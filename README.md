ghostscriptwasm
===

Go library with Ghostscript (currently 10.5.1), built to WASM, run with Wazero with fake memory filesystem.

This is mostly just making sense as a part of my PDF-to-PDFA conversion tool. But the whole thing can be use as ghostscript in the backend, from go, without cgo, on all go platforms!

It works on all platforms where Go works; it was tested on Linux, Windows and Mac.

### How to use
See pdf2pdfa3b, which is a library + a cmd tool that tries to convert from PDF to PDF/A 3b.

You first need to init new GS runtime with `ghostscriptwasm.New()`, which creates a new runtime. Then you can use it with .Run() or .BasicRun(). Each run creates a new fake directory in memory. There is no read or write from disk; all files are either in memory or embeded in the binary.

### License

NOTE THAT THE WHOLE THING IS AFFERO GPL, because Ghostscript itself is Affero GPL. This can be a problem for your backend, consult your copyright lawyers.

If you want a GPLv3 version, you can use GS version 9.06 (it's safe as it is contained with no access to FS) - see here https://github.com/karelbilek/ghostscript-9.06 - but I don't want to maintain it.

The build of the wasm is a bit crazy, as it uses a custom fork of emscripten that produces
wasm that only works in wazero. But it works...?

License:

Affero GPL

Copyright

(C) 2023 Karel Bilek, Jeroen Bobbeldijk (https://github.com/jerbob92)

Ghostscript license details are in build/ghostpdl/ submodule. The current GS version is 10.05.1