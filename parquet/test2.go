// +build test2

package main

import (
	"math"
)

type Sample struct {
	Day       int32 `parquet:"name=day, type=INT32, encoding=RLE_DICTIONARY"`
	TimeOfDay int64 `parquet:"name=time_of_day, type=INT64, encoding=DELTA_BINARY_PACKED"`
	Azimuth   int64 `parquet:"name=azimuth, type=INT64, encoding=DELTA_BINARY_PACKED"`
	Elevation int64 `parquet:"name=elevation, type=INT64, encoding=DELTA_BINARY_PACKED"`
}

func generateSample(i int) Sample {
	t := 12345.6789 + 0.005*float64(i)
	az := 2 * math.Sin(t)
	el := 45 + 2*math.Cos(t)

	it := int64(t * 1e9)
	iaz := int64(az * 1e9)
	iel := int64(el * 1e9)
	return Sample{
		Day:       310,
		TimeOfDay: it,
		Azimuth:   iaz,
		Elevation: iel,
	}
}
