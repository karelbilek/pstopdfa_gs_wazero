package pstopdfa_gs_wazero

import (
	"context"
	"fmt"
	"strings"
)

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

func (gs *GS) basicRun(ctx context.Context, in []byte, opts []string) ([]byte, error) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	files, err := gs.Run(ctx, stdout, stderr, opts, []File{{
		Path:    `infile`,
		Content: in,
	}})
	if err != nil {
		return nil, &ErrRunningGhostscript{
			Err:    err,
			StdOut: stdout.String(),
			StdErr: stderr.String(),
		}
	}
	for k, v := range files {
		if k == "outfile" {
			return v.Content, nil
		}
	}
	return nil, &ErrMissingFile{
		ResultFileName: "outfile",
		StdOut:         stdout.String(),
		StdErr:         stderr.String(),
	}
}

// PDF2PDFA converts PDF to PDF/A.
func (gs *GS) PS2PDFA3B(ctx context.Context, pdf []byte) ([]byte, error) {
	return gs.basicRun(ctx, pdf, []string{
		`-dPDFA=2`, // hack... we act like it's 3... we hard-code 3 with our patched GhostScript...
		`-dBATCH`,
		`-dNOPAUSE`,
		`-dUseCIEColor`,
		`-sProcessColorModel=DeviceCMYK`,
		`-dUseCIEColor`,
		`-sColorConversionStrategy=/UseDeviceIndependentColor`,
		`-sDEVICE=pdfwrite`,
		`-dPDFACompatibilityPolicy=2`,
		`-sOutputFile=outfile`,
		`infile`,
	})
}
