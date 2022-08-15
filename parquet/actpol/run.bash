#!/bin/bash
set -efu -o noclobber -o pipefail

./extract-corpus.bash

go fmt
for tag in plain delta; do
    go build -tags "$tag" -o "write-$tag"
done

for tag in delta; do
    for compression in zstd; do
        for pagesize in 8192 65536; do
            ./"write-$tag" -compression "$compression" corpus.zip > "${tag}-${compression}-pagesize${pagesize}.csv"
        done
    done
done

