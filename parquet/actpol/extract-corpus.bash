#!/bin/bash
set -eu -o noclobber -o pipefail

if [[ -f corpus.zip ]]; then
    echo "corpus.zip already exists"
    exit 0
fi

# run on niagara 2022-08-10

datadir="/project/r/rbond/actpol/data"

# list all files in zirfiles
census_list="census.txt.gz"
(
    cd "$datadir"
    for filename in season?/merlin/*/*.ar?.zip; do
        unzip -l "$filename" | tail -n+4 | head -n-2 | gawk "{ print \$1, \"$filename/\"\$4 }"
    done
) | gzip > "$census_list"

# extract just the TES channels
tesdata_list="tesdata.txt.gz"
zcat census.txt.gz | grep tesdata | grep slm | grep -v corrupt | gzip > "$tesdata_list"

# create random sample list
corpus_list="$PWD/corpus-list.txt"
zcat "$tesdata_list" | shuf -n 10000 > "$corpus_list"

# extract random sample from zirfiles
tmpdir=$(mktemp -d)
pushd "$tmpdir"
mkdir corpus
pushd corpus
i=0
for line in $(gawk '{print $2}' "$corpus_list"); do
    zipfile="$datadir/${line%/*}"
    tesdata_slm="${line##*/}"
    tesdata="${tesdata_slm%.slm}"

    # unzip & get size
    unzip "$zipfile" "$tesdata_slm"
    slim_sz=$(stat --printf='%s' "$tesdata_slm")

    # unslim & get size
    unslim "$tesdata_slm"
    uncompressed_sz=$(stat --printf='%s' "$tesdata")

    # rename
    newname=$(printf "%04d.i32" $i)
    mv "$tesdata" "$newname"

    echo "$newname $uncompressed_sz $slim_sz" >> manifest.txt

    i=$((i + 1))
done
popd
zip -rmT corpus.zip corpus
popd
mv "$tmpdir/corpus.zip" .

echo ok
