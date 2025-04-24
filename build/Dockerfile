FROM emscripten/emsdk:4.0.7

# build and install old autoconf
RUN apt update
RUN apt-get install -y m4 libtool


WORKDIR /
RUN curl -O http://ftp.gnu.org/gnu/autoconf/autoconf-2.69.tar.gz
RUN tar zxvf autoconf-2.69.tar.gz
WORKDIR /autoconf-2.69
RUN ./configure
RUN make
RUN sudo make install

# other tools later needed for ghostscript
RUN apt-get update && apt-get install -y autotools-dev automake libtool

# build emscripten freetype
# RUN embuilder build freetype

# # copy GNU ghostscript code
COPY ghostscript-10.05.0 /ghostscript-src

# # replace gs with emscripten freetype
# RUN rm -Rf /ghostscript-src/freetype
# RUN mkdir /ghostscript-src/freetype
# RUN cp -R /emsdk/upstream/emscripten/cache/ports/freetype/FreeType-version_1/* /ghostscript-src/freetype

RUN rm -Rf /emsdk/upstream/emscripten/cache

# check in https://github.com/emscripten-core/emscripten/pull/18285
COPY patched_emscripten/system/lib/standalone/standalone.c /emsdk/upstream/emscripten/system/lib/standalone/standalone.c
COPY patched_emscripten/system/lib/libc/musl/src/stdio/freopen.c /emsdk/upstream/emscripten/system/lib/libc/musl/src/stdio/freopen.c

# we run configure in docker build, make then as a separate docker run as it's faster
# (also I kept hitting weird bugs in docker buildkit)
RUN mkdir /ghostscript-build

COPY build-inside-emconfigure.sh /ghostscript-build/build-inside-emconfigure.sh

WORKDIR /ghostscript-build
RUN ./build-inside-emconfigure.sh

# just copy the emmake script
COPY build-inside-emmake.sh /ghostscript-build/build-inside-emmake.sh