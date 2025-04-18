#!/usr/bin/env bash

docker build --progress=plain -t gs-wazero-build-3 . 2>&1 | tee build.log

docker run -it -v $(pwd)/out:/out:z gs-wazero-build-3 /ghostscript-build/build-inside-emmake.sh