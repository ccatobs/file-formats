// +build test1

package main

import (
	"math"
)

type Sample struct {
	Day       int32   `parquet:"name=day, type=INT32"`
	TimeOfDay float64 `parquet:"name=time_of_day, type=DOUBLE"`
	Azimuth   float64 `parquet:"name=azimuth, type=DOUBLE"`
	Elevation float64 `parquet:"name=elevation, type=DOUBLE"`
}

func generateSample(i int) Sample {
	t := 12345.6789 + 0.005*float64(i)
	az := 2 * math.Sin(t)
	el := 45 + 2*math.Cos(t)

	return Sample{
		Day:       310,
		TimeOfDay: t,
		Azimuth:   az,
		Elevation: el,
	}
}
