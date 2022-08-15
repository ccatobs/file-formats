# Parquet test: ACTpol data

Test using real world instrument data, specifically TES bolometer timestreams from [ACTpol](https://act.princeton.edu/). Can parquet compress the data as well as [slim](https://sourceforge.net/projects/slimdata/), ACT's custom encoding scheme?

## Corpus

From a set of 5 years of ACTpol TES bolometer data (72.6TB compressed, 270M files), we selected a random sample of 10,000 files (2.7GB compressed with slim, 8.9GB uncompressed) using the [extract-corpus.bash](extract-corpus.bash) script.

The files are a sequence of little-endian `int32` values.

## Parquet compression

The corpus was transcoded as parquet files using the [run.bash](run.bash) script and [write.go](write.go) program.

Go was used instead of python because as of the time of writing (August 2022) the Arrow C++ library did not support writing the `DELTA_BINARY_PACKED` encoding (see issue [PARQUET-490](https://issues.apache.org/jira/browse/PARQUET-490)).

## Results

Unfortunately, parquet files with delta encoding & zstd compression are about 25% larger than slim files. Changing the page size from the default 8kiB to 64kiB (same as slim) helps a little.

| compressed size [bytes] | compressed size [%] | compression factor | method |
| --- | --- | --- | --- |
| 8946113100 | 100% | 1.00x | raw |
| 3759819782 | 42% | 2.38x | zip |
| 3619588480 | 40% | 2.47x | delta/zstd/8k |
| 3404912818 | 38% | 2.63x | delta/zstd/64k |
| 2685028165 | 30% | 3.33x | slim |

## Why is slim better?

Slim's advantage comes from bit rotating the samples before delta encoding. Here's a few representative data points:

| sample | XOR with previous sample |
| --- | --- |
| `0x2311fc09` | |
| `0x23124e09` | `0x3b200` |
| `0x23127509` | `0x03b00` |
| `0x23124c89` | `0x03980` |
| `0x23125f09` | `0x01380` |
| `0x23124389` | `0x01c80` |
| `0x23120109` | `0x04280` |
| `0x2311bb09` | `0x3ba00` |

Note how the first 14 bits and the last 7 bits of each sample XOR'd with the previous sample are zero, i.e., they never change. The constant high-order bits are well compressed by delta encoding, but not the constant low-order bits. Before delta encoding, slim rotates each sample so that the constant low-order bits become constant high-order bits.

## Appendix

### Slim options

The ACTpol data was compressed with the following slim options:

```
slim --permit-bitrotation --deltas --compute-crc32 \
     --int --repeats 16384 <filename>
```

