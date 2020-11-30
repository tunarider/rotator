package cmd

import (
	"errors"
	"github.com/mholt/archiver/v3"
	"github.com/mitchellh/go-homedir"
	"github.com/tunarider/rotator/internal/config"
	"github.com/tunarider/rotator/internal/util"
	"github.com/urfave/cli/v2"
	"path/filepath"
)

func Rotate(conf *config.Config) error {
	for _, target := range conf.Targets {
		var err error
		target.Path, err = homedir.Expand(target.Path)
		if err != nil {
			return cli.Exit(err, 1)
		}
		if target.Action == config.ActionArchive {
			err := archive(target)
			if err != nil {
				return cli.Exit(err, 2)
			}
		}
	}
	return nil
}

func archive(target config.Target) error {
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
	if target.Group {
		var groups []util.FileGroup
		if target.GroupDate {
			groups, err = util.GroupByDate(target.Name, filterDate, files)
			if err != nil {
				return errors.New("failed to group target files")
			}
			for _, group := range groups {
				err = archiver.Archive(
					util.GetFilePathsFromGroup(target.Path, group),
					filepath.Join(target.Path, group.Name)+".tar.gz",
				)
				if err != nil {
					return err
				}
			}
		} else {
			err = archiver.Archive(targetPaths, filepath.Join(target.Path, target.Name)+".tar.gz")
			if err != nil {
				return err
			}
		}
	} else {
		for _, path := range targetPaths {
			err = archiver.CompressFile(path, path+".gz")
		}
	}
	return nil
}
