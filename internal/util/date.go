package util

import (
	"os"
	"regexp"
	"time"
)

func MakeDateFilter(expr string, format string) (DateFilter, error) {
	if format == "" {
		format = "20060102"
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return func(file os.FileInfo) (time.Time, error) {
		matches := re.FindStringSubmatch(file.Name())
		return time.Parse(format, matches[1])
	}, nil
}
