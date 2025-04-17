package main

import (
	"context"
	"os"

	"github.com/karelbilek/ghostscript_wazero"
)

func main() {
	fn := os.Args[1]
	data, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	gs := ghostscript_wazero.NewGS()
	re, err := gs.PS2PDFA3B(context.TODO(), data)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(os.Args[2], re, 0o666)
	if err != nil {
		panic(err)
	}
}
