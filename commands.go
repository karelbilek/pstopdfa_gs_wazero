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

func basicRun(ctx context.Context, in []byte, opts []string) ([]byte, error) {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	files, err := Run(ctx, stdout, stderr, opts, []File{{
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
	stderrS := stderr.String()

	if strings.Contains(stderrS, "The following errors") || strings.Contains(stderrS, "error occurred") {
		return nil, &ErrRunningGhostscript{
			Err:    fmt.Errorf("error on conversion"),
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

// PDF2PDFA converts PDF to PDF/A 3b.
func PDF2PDFA3b(ctx context.Context, paintFont bool, pdf []byte) ([]byte, error) {
	opts := []string{
		`-dNOSAFER`,
		`-dBATCH`,
		`-dNOPAUSE`,

		`-sDEVICE=pdfwrite`,
		`-dPDFA=3`,
		`-sColorConversionStrategy=RGB`,
		`-dPDFACompatibilityPolicy=2`,
		`-sOutputFile=outfile`,
		`/gs_profiles/pdfa_def.ps`,
		`infile`,
	}
	if paintFont {
		opts = []string{
			`-dNOSAFER`,
			`-dBATCH`,
			`-dNOPAUSE`,

			`-sDEVICE=pdfwrite`,
			`-dPDFA=3`,
			`-sColorConversionStrategy=RGB`,
			`-dPDFACompatibilityPolicy=2`,
			`-sOutputFile=outfile`,
			`-dNoOutputFonts`,
			`/gs_profiles/pdfa_def.ps`,
			`infile`,
		}
	}
	return basicRun(ctx, pdf, opts)
}
