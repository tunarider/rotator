package cmd

import (
	"errors"
	"fmt"
	"github.com/mholt/archiver/v3"
	"github.com/mitchellh/go-homedir"
	"github.com/tunarider/rotator/internal/config"
	"github.com/tunarider/rotator/internal/util"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

func Rotate(conf *config.Config, dry bool) error {
	for _, target := range conf.Targets {
		var err error
		target.Path, err = homedir.Expand(target.Path)
		if err != nil {
			return cli.Exit(err, 1)
		}
		if target.Action == config.ActionArchive {
			err := archive(target, dry)
			if err != nil {
				return cli.Exit(err, 2)
			}
		}
	}
	return nil
}

func archive(target config.Target, dry bool) error {
	var targetPaths []string
	filterDate, err := util.MakeDateFilter(target.Regexp, target.DateFormat)
	if err != nil {
		return errors.New("failed to init date filter")
	}
	checkTarget, err := util.MakeTargetChecker(target.Retention, filterDate)
	if err != nil {
		return errors.New("failed to init target checker")
	}
	files, err := util.GetTargetFiles(checkTarget, target)
	if err != nil {
		return errors.New("failed to filter target files")
	}
	for _, file := range files {
		targetPaths = append(targetPaths, filepath.Join(target.Path, file.Name()))
	}
	if target.GroupBy == config.GroupByDate {
		var groups []util.FileGroup
		groups, err = util.GroupByDate(target.Name, filterDate, files)
		if err != nil {
			return errors.New("failed to group target files")
		}
		for _, group := range groups {
			targetPaths := util.GetFilePathsFromGroup(target.Path, group)
			err = archiveByGroup(target, targetPaths, filepath.Join(target.Path, group.Name)+".tar.gz", dry)
			if err != nil {
				return err
			}
		}
	} else if target.GroupBy == config.GroupByName {
		return archiveByGroup(target, targetPaths, filepath.Join(target.Path, target.Name)+".tar.gz", dry)
	} else if target.GroupBy == config.GroupByFile {
		for _, path := range targetPaths {
			err = archiveFile(target, path, dry)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func archiveByGroup(target config.Target, targetPaths []string, output string, dry bool) error {
	if dry {
		fmt.Printf("archive:\n  group: %s\n", output)
		for _, p := range targetPaths {
			fmt.Printf("    file: %s\n", p)
		}
	} else {
		err := archiver.Archive(targetPaths, output)
		if err != nil {
			return err
		}
		if target.Remove {
			for _, p := range targetPaths {
				err = os.Remove(p)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func archiveFile(target config.Target, path string, dry bool) error {
	if dry {
		fmt.Printf("archive: %s to %s\n", path, path+".gz")
	} else {
		err := archiver.CompressFile(path, path+".gz")
		if err != nil {
			return err
		}
		if target.Remove {
			return os.Remove(path)
		}
	}
	return nil
}
