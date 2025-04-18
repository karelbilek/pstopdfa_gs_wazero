#!/bin/bash
set -euo pipefail

OUT_DIR="$PWD/out"
ROOT="$PWD"
EMCC_FLAGS_DEBUG="-g"
EMCC_FLAGS_RELEASE="-O2"

export CPPFLAGS="-I$OUT_DIR/include"
export LDFLAGS="-L$OUT_DIR/lib"
export PKG_CONFIG_PATH="$OUT_DIR/lib/pkgconfig"
export EM_PKG_CONFIG_PATH="$PKG_CONFIG_PATH"
export CFLAGS="$EMCC_FLAGS_DEBUG"
export CXXFLAGS="$CFLAGS"
export TARGET_ARCH_FILE="/ghostscript-src/arch/wasm.h"
#export EMCC_DEBUG=1

mkdir -p "$OUT_DIR"

cd "/ghostscript-src"

export GS_LDFLAGS="\
-s ALLOW_MEMORY_GROWTH=1 \
-s WASM=1 \
-s ALLOW_MEMORY_GROWTH=1 \
-s STANDALONE_WASM=1 \
-sERROR_ON_UNDEFINED_SYMBOLS=0 \
-s USE_ZLIB=1 \
-s WASM_BIGINT=1 \
-g \
--profile"

mkdir -p /ghostscript
sudo chmod 777 /ghostscript

nproc | xargs -I % emmake make \
  LDFLAGS="$LDFLAGS $GS_LDFLAGS" \
  prefix="/ghostscript" \
  -j% install

rm -r /out/*
mkdir -p /out/ghostscript/share/ghostscript/9.06
cp -r /ghostscript/share/ghostscript/9.06/lib /out/ghostscript/share/ghostscript/9.06/lib
cp bin/gs.wasm /out/gs.wasm