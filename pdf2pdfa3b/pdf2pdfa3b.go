package pdf2pdfa3b

import (
	"context"
	"fmt"
	"strings"

	"github.com/karelbilek/ghostscriptwasm"
)

type ErrPdfaReverted struct {
	StdErr string
}

func (e *ErrPdfaReverted) Error() string {
	return fmt.Sprintf("PDF/A reverted to PDF. Stderr:\n%s", e.StdErr)
}

// PDF2PDFA converts PDF to PDF/A 3b.
func PDF2PDFA3b(ctx context.Context, gs *ghostscriptwasm.GS, policy int, paintFont bool, pdf []byte) ([]byte, error) {

	opts := []string{
		`-dNOSAFER`,
		`-dBATCH`,
		`-dNOPAUSE`,

		`-sDEVICE=pdfwrite`,
		`-dPDFA=3`,
		`-sColorConversionStrategy=RGB`,
		fmt.Sprintf(`-dPDFACompatibilityPolicy=%d`, policy),
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
			fmt.Sprintf(`-dPDFACompatibilityPolicy=%d`, policy),
			`-sOutputFile=outfile`,
			`-dNoOutputFonts`,
			`/gs_profiles/pdfa_def.ps`,
			`infile`,
		}
	}
	bs, stderr, err := gs.BasicRun(ctx, pdf, opts)
	if err != nil {
		return nil, err
	}
	// this is how ghostscript marks reversion to from PDF/A to PDF
	// with policy=1
	if strings.Contains(stderr, "reverting to") {
		return nil, &ErrPdfaReverted{
			StdErr: stderr,
		}
	}
	return bs, err
}
