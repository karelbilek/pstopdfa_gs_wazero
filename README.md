pstopdfa_gs_wazero
===

Convert PS to PDF-a/3b.

This includes an entire copy of patched GhostScript 9.06 (the last one GPL).

The patched ghostscript is compiled with patched emscripten to wasm, which is then run with experimental memory FS (because old ghostscript is VERY unsafe).

This is mostly just making sense as a part of PDF-to-PDFA conversion, which I will publish soon.

License
GPLv3

Copyright
(C) 2023 Karel Bilek, Jeroen Bobbeldijk (https://github.com/jerbob92)
For ghostscript license see https://github.com/karelbilek/ghostscript-9.06/blob/main/LICENSE