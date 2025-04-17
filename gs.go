package pstopdfa_gs_wazero

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"

	"github.com/tetratelabs/wazero/experimental/sys"
	"github.com/tetratelabs/wazero/experimental/sysfs"

	"github.com/karelbilek/wazero-fs-tools/memfs"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/emscripten"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed build/out/gs.wasm
var gsWasm []byte

//go:embed build/out/ghostscript
var sharedFiles embed.FS

type GS struct {
	module  wazero.CompiledModule
	runtime wazero.Runtime
}

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

	sub, err := fs.Sub(sharedFiles, "build/out/ghostscript")
	if err != nil {
		panic(err)
	}

	fsConfig := wazero.NewFSConfig()
	sharedFS := &sysfs.AdaptFS{FS: sub}

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

	fsConfig = fsConfig.(sysfs.FSConfig).WithSysFSMount(sharedFS, "/ghostscript")
	fsConfig = fsConfig.(sysfs.FSConfig).WithSysFSMount(tracked, "/")

	gsArgs := []string{"gs"}
	gsArgs = append(gsArgs, args...)

	moduleConfig := wazero.NewModuleConfig().
		WithStartFunctions("_start").
		WithStdout(stdOut).
		WithStderr(stdErr).
		WithFSConfig(fsConfig).
		WithName("").
		WithArgs(gsArgs...)

	_, err = gs.runtime.InstantiateModule(ctx, gs.module, moduleConfig)
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

func NewGS() *GS {
	ctx := context.Background()
	runtimeConfig := wazero.NewRuntimeConfig()
	wazeroRuntime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, wazeroRuntime); err != nil {
		log.Fatal(err)
	}

	compiledModule, err := wazeroRuntime.CompileModule(ctx, gsWasm)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := emscripten.InstantiateForModule(ctx, wazeroRuntime, compiledModule); err != nil {
		log.Fatal(err)
	}

	return &GS{
		module:  compiledModule,
		runtime: wazeroRuntime,
	}
}
