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
	pstopdfa_gs_wazero.DoInit()
	re, err := pstopdfa_gs_wazero.PDF2PDFA3b(context.TODO(), false, data)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(os.Args[2], re, 0o666)
	if err != nil {
		panic(err)
	}
}
