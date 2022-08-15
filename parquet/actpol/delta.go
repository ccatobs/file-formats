//go:build delta

package main

const encoding = "delta"

type Sample struct {
	Sample int32 `parquet:"name=sample, type=INT32, encoding=DELTA_BINARY_PACKED"`
}
