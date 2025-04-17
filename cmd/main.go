package main

import (
	"context"
	"os"

	"github.com/karelbilek/pstopdfa_gs_wazero"
)

func main() {
	fn := os.Args[1]
	data, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	gs := pstopdfa_gs_wazero.NewGS()
	re, err := gs.PS2PDFA3B(context.TODO(), data)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(os.Args[2], re, 0o666)
	if err != nil {
		panic(err)
	}
}
