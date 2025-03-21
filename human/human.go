package human

import (
	"fmt"
	"math"
)

const (
	B   = 1
	KB  = 1e3
	MB  = 1e6
	GB  = 1e9
	TB  = 1e12
	PB  = 1e15
	KiB = int64(1) << 10
	MiB = int64(1) << 20
	GiB = int64(1) << 30
	TiB = int64(1) << 40
	PiB = int64(1) << 50
)

func Dehumanize(value string) (result int64, err error) {
	var s float64
	var suffix string

	if _, err = fmt.Sscanf(value, "%f%s", &s, &suffix); err != nil {
		return
	}

	var factor int64
	switch suffix {
	case "B":
		factor = B
	case "kB":
		factor = KB
	case "MB":
		factor = MB
	case "GB":
		factor = GB
	case "TB":
		factor = TB
	case "PB":
		factor = PB
	case "KiB", "kb":
		factor = KiB
	case "MiB", "m", "mb":
		factor = MiB
	case "GiB", "g", "gb":
		factor = GiB
	case "TiB":
		factor = TiB
	case "PiB":
		factor = PiB
	default:
		factor = 1
	}

	result = int64(s * float64(factor))
	return
}

func Humanize(s float64, iec bool) (string, string) {
	sizes := []string{" B", " kB", " MB", " GB", " TB", " PB", " EB"}
	base := 1000.0

	if iec {
		sizes = []string{" B", " KiB", " MiB", " GiB", " TiB", " PiB", " EiB"}
		base = 1024.0
	}

	if s < 10 {
		return fmt.Sprintf("%2.0f", s), sizes[0]
	}
	e := math.Floor(logN(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f, val), suffix
}

func logN(n, base float64) float64 {
	return math.Log(n) / math.Log(base)
}
