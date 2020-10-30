#!/bin/bash
set -eu -o pipefail

announce() {
    echo -e "\n### $(tput bold)$*$(tput sgr0) ###\n"
}

for test in test1 test2; do
    filename="$test.parquet"
    [[ ! -f $filename ]] || rm "$filename"

    announce "running Go write $test"
    go build -tags "$test" -o write_go
    ./write_go "$filename"

    announce "running Python read $test"
    ./read.py "$filename"
done

