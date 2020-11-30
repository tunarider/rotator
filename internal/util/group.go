package util

import (
	"os"
	"path"
	"strings"
)

type FileGroup struct {
	Name string
	Files []os.FileInfo
}

func GroupByDate(base string, filter DateFilter, files []os.FileInfo) ([]FileGroup, error) {
	var groups []FileGroup
	for _, file := range files {
		date, err := filter(file)
		if err != nil {
			return nil, err
		}
		groupName := strings.Join([]string{base, date.Format("20060102")}, "-")
		isExists := false
		for i, g := range groups {
			if g.Name == groupName {
				groups[i].Files = append(groups[i].Files, file)
				isExists = true
			}
		}
		if !isExists {
			group := FileGroup{
				Name: groupName,
				Files: []os.FileInfo{file},
			}
			groups = append(groups, group)
		}
	}
	return groups, nil
}

func GetFilePathsFromGroup(root string, group FileGroup) []string {
	var paths []string
	for _, file := range group.Files {
		paths = append(paths, path.Join(root, file.Name()))
	}
	return paths
}