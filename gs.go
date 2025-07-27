package ghostscriptwazero

import (
	"context"
	"crypto/rand"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/tetratelabs/wazero/experimental/sys"
	"github.com/tetratelabs/wazero/experimental/sysfs"

	"github.com/karelbilek/wazero-fs-tools/memfs"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/emscripten"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed build/out/gs.wasm
var gsWasm []byte

//go:embed build/out/ghostscript_lib
var sharedFiles embed.FS

//go:embed gs_profiles
var gsProfiles embed.FS

type File struct {
	Path    string
	Content []byte
}

type memFSTrackFiles struct {
	*memfs.MemFS
	files []string
}

func (m *memFSTrackFiles) OpenFile(path string, flag sys.Oflag, perm fs.FileMode) (sys.File, sys.Errno) {
	f, e := m.MemFS.OpenFile(path, flag, perm)
	if e != 0 {
		return f, e
	}

	// just dumbly track all successfully created files, do not track unlinks
	if flag&sys.O_CREAT != 0 {
		m.files = append(m.files, path)
	}
	return f, e
}

var (
	ErrCannotWriteInitialFile   = errors.New("cannot write initial file")
	ErrRunningGhostscriptModule = errors.New("error running ghostscript module")
	ErrCannotReadResultFile     = errors.New("cannot read result file")
)

func (gs *GS) Run(ctx context.Context, stdOut, stdErr io.Writer, args []string, files []File) (map[string]File, error) {
	fsConfig := wazero.NewFSConfig()

	sharedFS := &sysfs.AdaptFS{FS: sharedFilesCut}
	gsProfilesFS := &sysfs.AdaptFS{FS: gsProfilesCut}

	memFS := memfs.New()

	tracked := &memFSTrackFiles{MemFS: memFS}

	errno := memFS.Mkdir("tmp", 0)
	if errno != 0 {
		return nil, fmt.Errorf("%w %q: %w", ErrCannotWriteInitialFile, "tmp", errno)
	}
	for _, f := range files {
		// this will not track the file, as it just goes to the embedded memfs
		errno := memFS.WriteFile(f.Path, f.Content)
		if errno != 0 {
			return nil, fmt.Errorf("%w %q: %w", ErrCannotWriteInitialFile, f.Path, errno)
		}
	}

	// memFStmp := memfs.New()
	fsConfig = fsConfig.(sysfs.FSConfig).WithSysFSMount(sharedFS, "/ghostscript/share/ghostscript/10.06.0/lib")
	// fsConfig = fsConfig.(sysfs.FSConfig).WithSysFSMount(memFStmp, "/tmp")
	fsConfig = fsConfig.(sysfs.FSConfig).WithSysFSMount(gsProfilesFS, "/gs_profiles")

	fsConfig = fsConfig.(sysfs.FSConfig).WithSysFSMount(tracked, "/")

	gsArgs := []string{"gs"}
	gsArgs = append(gsArgs, args...)

	moduleConfig := wazero.NewModuleConfig().
		WithStartFunctions("_start").
		WithStdout(stdOut).
		WithStderr(stdErr).
		WithFSConfig(fsConfig).
		WithRandSource(rand.Reader).
		WithSysWalltime().
		WithSysNanotime().
		WithSysNanosleep().
		WithName("").
		WithArgs(gsArgs...)

	_, err := gs.wruntime.InstantiateModule(ctx, gs.compiled, moduleConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRunningGhostscriptModule, err)
	}

	resFiles := map[string]File{}

	// now we know that we run successfully, let's read the files
	for _, f := range tracked.files {
		cont, errno := tracked.ReadFile(f)
		if errno != 0 {
			if errno == sys.ENOENT {
				// skipping removed tmp files
				continue
			}
			return nil, fmt.Errorf("%w %q: %w", ErrCannotReadResultFile, f, errno)
		}
		resFiles[f] = File{
			Path:    f,
			Content: cont,
		}
	}

	return resFiles, nil
}

type ErrRunningGhostscript struct {
	Err    error
	StdOut string
	StdErr string
}

func (e *ErrRunningGhostscript) Error() string {
	return fmt.Sprintf("error running ghostscript: %v\n\nStdout:\n%s\n\nStderr:\n%s", e.Err, e.StdOut, e.StdErr)
}

func (e *ErrRunningGhostscript) Unwrap() error {
	return e.Err
}

type ErrMissingFile struct {
	ResultFileName string
	StdOut         string
	StdErr         string
}

func (e *ErrMissingFile) Error() string {
	return fmt.Sprintf("missing ghostscript result file %s.\n\nStdout:\n%s\n\nStderr:\n%s", e.ResultFileName, e.StdOut, e.StdErr)
}

func (gs *GS) BasicRun(ctx context.Context, in []byte, opts []string) ([]byte, string, error) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	files, err := gs.Run(ctx, stdout, stderr, opts, []File{{
		Path:    `infile`,
		Content: in,
	}})
	if err != nil {
		return nil, stderr.String(), &ErrRunningGhostscript{
			Err:    err,
			StdOut: stdout.String(),
			StdErr: stderr.String(),
		}
	}
	stderrS := stderr.String()

	if strings.Contains(stderrS, "The following errors") || strings.Contains(stderrS, "error occurred") {
		return nil, stderrS, &ErrRunningGhostscript{
			Err:    fmt.Errorf("error on conversion"),
			StdOut: stdout.String(),
			StdErr: stderrS,
		}
	}

	for k, v := range files {
		if k == "outfile" {
			return v.Content, stderrS, nil
		}
	}
	return nil, stderrS, &ErrMissingFile{
		ResultFileName: "outfile",
		StdOut:         stdout.String(),
		StdErr:         stderr.String(),
	}
}

type GS struct {
	compiled wazero.CompiledModule
	wruntime wazero.Runtime
}

var gsProfilesCut fs.FS
var sharedFilesCut fs.FS

func init() {
	var err error
	gsProfilesCut, err = fs.Sub(gsProfiles, "gs_profiles")
	if err != nil {
		panic(err)
	}

	sharedFilesCut, err = fs.Sub(sharedFiles, "build/out/ghostscript_lib")
	if err != nil {
		panic(err)
	}
}

func NewGS(ctx context.Context) (*GS, error) {
	// ctx := experimental.WithFunctionListenerFactory(context.Background(), logging.NewHostLoggingListenerFactory(os.Stdout, logging.LogScopeAll))
	runtimeConfig := wazero.NewRuntimeConfig()
	wazeroRuntime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, wazeroRuntime); err != nil {
		return nil, err
	}

	compiledModule, err := wazeroRuntime.CompileModule(ctx, gsWasm)
	if err != nil {
		return nil, err
	}
	if _, err := emscripten.InstantiateForModule(ctx, wazeroRuntime, compiledModule); err != nil {
		return nil, err
	}

	return &GS{
		compiled: compiledModule,
		wruntime: wazeroRuntime,
	}, nil
}
