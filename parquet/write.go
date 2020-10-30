package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

func main() {
	start := time.Now()
	filename := os.Args[1]

	fw, err := local.NewLocalFileWriter(filename)
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}

	pw, err := writer.NewParquetWriter(fw, new(Sample), 4)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}

	row_bytes := binary.Size(new(Sample))
	str_row_bytes := fmt.Sprintf("%d", row_bytes)
	kv := &parquet.KeyValue{"row_bytes", &str_row_bytes}
	pw.Footer.KeyValueMetadata = append(pw.Footer.KeyValueMetadata, kv)

	//pw.RowGroupSize = 128 * 1024 * 1024 //128M
	//pw.PageSize = 8 * 1024              //8K
	//pw.CompressionType = parquet.CompressionCodec_UNCOMPRESSED
	//pw.CompressionType = parquet.CompressionCodec_SNAPPY
	//pw.CompressionType = parquet.CompressionCodec_GZIP
	pw.CompressionType = parquet.CompressionCodec_ZSTD
	num := 1000000
	for i := 0; i < num; i++ {
		s := generateSample(i)
		if i < 5 || i >= num-5 {
			fmt.Printf("%d %v\n", i, s)
		} else if i == 5 {
			fmt.Println("...")
		}
		if err = pw.Write(s); err != nil {
			log.Println("Write error", err)
		}
	}
	if err = pw.WriteStop(); err != nil {
		log.Println("WriteStop error", err)
		return
	}
	fw.Close()
	elapsed := time.Now().Sub(start)

	// performance stats
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	filesz := float64(stat.Size())
	memsz := float64(num * row_bytes)
	fmt.Println()
	fmt.Printf("write speed = %g MB/s\n", memsz/float64(elapsed.Microseconds()))
	fmt.Printf("compression factor = %g\n", memsz/filesz)
}
