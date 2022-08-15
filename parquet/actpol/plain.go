//go:build plain

package main

const encoding = "plain"

type Sample struct {
	Sample int32 `parquet:"name=sample, type=INT32, encoding=PLAIN"`
}
