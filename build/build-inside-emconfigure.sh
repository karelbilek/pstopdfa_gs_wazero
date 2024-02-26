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

# There is a bug in this version of Ghostscript that prevents passing in gcc to compile the build tools, replace the var manually.
sed -i "s/CCAUX=@CC@/CCAUX=gcc/g" base/Makefile.in

emconfigure ./autogen.sh \
  CFLAGSAUX= CPPFLAGSAUX= \
  --host="wasm32-unknown-linux" \
  --prefix="$OUT_DIR" \
  --disable-cups \
  --disable-dbus \
  --disable-gtk \
  --with-system-libtiff
