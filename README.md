pstopdfa_gs_wazero
===

Convert PS to PDF-a/3b.

This includes an entire ghostscript 10.5.0, unpatched.

Ghostscript is compiled with patched emscripten to wasm, which is then run with memory FS.

This is mostly just making sense as a part of my PDF-to-PDFA conversion tool, which I will publish soon. But the whole thing can be use as ghostscript in the backend, from go, with cgo-less build.

It DOES work on Mac, Linux and Windows.

NOTE THAT THE WHOLE THING IS AFFERO GPL - that means, if you use it in a backend as a library, you need to provide sources of the entire backend! For GPLv3 version, you can use GS version 9.06 (it's safe as it is contained with no access to FS) - see here https://github.com/karelbilek/ghostscript-9.06 - and compile it yourself, as done in build/ (it should work, but I no longer want to maintain the build). I used to use the GPLv3 version, but it's just too buggy and old; it can work for your needs.

License
Affero GPL

Copyright
(C) 2023 Karel Bilek, Jeroen Bobbeldijk (https://github.com/jerbob92)
For ghostscript license see https://github.com/karelbilek/ghostscript-9.06/blob/main/LICENSE