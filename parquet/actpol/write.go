package main

import (
	"archive/zip"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

func fileSize(filename string) (int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

type Result struct {
	Filename          string
	Encoding          string
	Compression       string
	PageSize          int64
	CompressedSize    int64
	CompressionFactor float64
	WriteSpeedMBps    float64
	WriteSpeedMHz     float64
}

func main() {
	var compression string
	flag.StringVar(&compression, "compression", "none", "compression kind")
	var pageSize int64
	flag.Int64Var(&pageSize, "pageSize", 8*1024, "page size")
	flag.Parse()

	// Open a zip archive for reading.
	z, err := zip.OpenReader(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer z.Close()

	var results []Result

	// iterate through the files in the corpus
	for _, f := range z.File {
		ifilename := f.Name
		if !strings.HasSuffix(ifilename, ".i32") {
			continue
		}
		//fmt.Printf("Contents of %s:\n", ifilename)
		fr, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		num := int(f.UncompressedSize64 / 4)
		data := make([]int32, num)
		err = binary.Read(fr, binary.LittleEndian, &data)
		if err != nil {
			log.Fatal(err)
		}
		if err := fr.Close(); err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		ofilename := path.Base(ifilename) + ".parquet"
		fw, err := local.NewLocalFileWriter(ofilename)
		if err != nil {
			log.Fatal("Can't create local file", err)
		}

		pw, err := writer.NewParquetWriter(fw, new(Sample), 2)
		if err != nil {
			log.Fatal("Can't create parquet writer", err)
		}

		switch compression {
		case "gzip":
			pw.CompressionType = parquet.CompressionCodec_GZIP
		case "lz4":
			pw.CompressionType = parquet.CompressionCodec_LZ4
		case "none":
			pw.CompressionType = parquet.CompressionCodec_UNCOMPRESSED
		case "snappy":
			pw.CompressionType = parquet.CompressionCodec_SNAPPY
		case "zstd":
			pw.CompressionType = parquet.CompressionCodec_ZSTD
		default:
			log.Fatal("unknown compression: ", compression)
		}

		pw.PageSize = pageSize

		for i, x := range data {
			s := &Sample{x}
			/*
				if i < 5 || i >= num-5 {
					fmt.Printf("%d %v\n", i, s)
				} else if i == 5 {
					fmt.Println("...")
				}
			*/
			if err := pw.Write(s); err != nil {
				log.Fatalf("Write error at %d %v: %v\n", i, s, err)
			}
		}
		if err := pw.WriteStop(); err != nil {
			log.Fatal("WriteStop error", err)
		}
		if err := fw.Close(); err != nil {
			log.Fatal("Close error", err)
		}

		elapsed_secs := 1e-9 * float64(time.Now().Sub(start).Nanoseconds())

		row_bytes := binary.Size(new(Sample))
		memsz := float64(num * row_bytes)
		filesz, err := fileSize(ofilename)
		if err != nil {
			log.Fatal("Can't get file size", err)
		}

		if err := os.Remove(ofilename); err != nil {
			log.Fatal(err)
		}

		results = append(results, Result{ifilename, encoding, compression, pageSize, filesz, memsz / float64(filesz), 1e-6 * memsz / elapsed_secs, 1e-6 * float64(num) / elapsed_secs})

	}

	fmt.Printf("input,encoding,compression,page_size,compressed_file_size,compression_factor,write_speed_MBps,write_speed_MHz\n")
	for _, result := range results {
		fmt.Printf("%s,%s,%s,%d,%d,%g,%g,%g\n",
			result.Filename, result.Encoding, result.Compression, result.PageSize, result.CompressedSize, result.CompressionFactor, result.WriteSpeedMBps, result.WriteSpeedMHz)
	}
}
