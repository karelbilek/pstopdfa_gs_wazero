package main

import (
	"context"
	"os"

	"github.com/karelbilek/ghostscriptwasm"
	"github.com/karelbilek/ghostscriptwasm/pdf2pdfa3b"
)

func main() {
	fn := os.Args[1]
	data, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	gs, err := ghostscriptwasm.NewGS(context.Background())
	if err != nil {
		panic(err)
	}

	re, err := pdf2pdfa3b.PDF2PDFA3b(context.TODO(), gs, 1, false, data)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(os.Args[2], re, 0o666)
	if err != nil {
		panic(err)
	}
}
