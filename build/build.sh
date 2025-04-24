#!/usr/bin/env bash

docker build --progress=plain -t gs-wazero-build . 2>&1 | tee build.log

docker run -it -v $(pwd)/out:/out:z gs-wazero-build /ghostscript-build/build-inside-emmake.sh