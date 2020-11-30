package util

import (
	"github.com/tunarider/rotator/internal/config"
	"io/ioutil"
	"os"
	"time"
)

type TargetChecker func(os.FileInfo) bool

type DateFilter func(os.FileInfo) (time.Time, error)

func GetTargetFiles(checkTarget TargetChecker, target config.Target) ([]os.FileInfo, error) {
	contents, err := ioutil.ReadDir(target.Path)
	if err != nil {
		return nil, err
	}
	return filterTargetFiles(checkTarget, contents), nil
}

func filterTargetFiles(checkTarget TargetChecker, files []os.FileInfo) []os.FileInfo {
	var filtered []os.FileInfo
	for _, file := range files {
		isTarget := checkTarget(file)
		if isTarget {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func MakeTargetChecker(retention int, filterDate DateFilter) (TargetChecker, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	criteria := today.AddDate(0, 0, -1*retention)
	return func(file os.FileInfo) bool {
		if file.IsDir() {
			return false
		}
		fileDate, err := filterDate(file)
		if err != nil {
			return false
		}
		if criteria.After(fileDate) {
			return true
		}
		return false
	}, nil
}
