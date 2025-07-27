FROM emscripten/emsdk:4.0.7

# other tools later needed for ghostscript
RUN apt-get update && apt-get install -y autotools-dev automake libtool

# copy ghostscript code
COPY ghostpdl /ghostscript-src

RUN rm -Rf /emsdk/upstream/emscripten/cache

COPY patched_emscripten/system/lib/standalone/standalone.c /emsdk/upstream/emscripten/system/lib/standalone/standalone.c
COPY patched_emscripten/system/lib/libc/musl/src/stdio/freopen.c /emsdk/upstream/emscripten/system/lib/libc/musl/src/stdio/freopen.c

# we run `configure` in docker build, `make` in docker run
RUN mkdir /ghostscript-build

COPY build-inside-emconfigure.sh /ghostscript-build/build-inside-emconfigure.sh

WORKDIR /ghostscript-build
RUN ./build-inside-emconfigure.sh

# just copy the emmake script
COPY build-inside-emmake.sh /ghostscript-build/build-inside-emmake.sh