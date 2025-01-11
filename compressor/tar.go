package compressor

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/itgcloud/gobackup/helper"
	"github.com/itgcloud/gobackup/logger"
)

type Tar struct {
	Base
}

func (tar *Tar) perform() (string, error) {
	filePath := tar.archiveFilePath(tar.ext)

	opts, err := tar.options(filePath)
	if err != nil {
		return "", err
	}

	_, err = helper.Exec("tar", opts...)

	return filePath, err
}

func (tar *Tar) options(archiveFilePath string) ([]string, error) {
	var opts []string

	if helper.IsGnuTar {
		opts = append(opts, "--ignore-failed-read")
	}

	compressArgs := []string{"-a"}
	if len(tar.parallelProgram) > 0 {
		if path, err := exec.LookPath(tar.parallelProgram); err == nil {
			compressArgs = []string{"--use-compress-program", path}
		}
	}
	opts = append(opts, compressArgs...)

	includes := tar.model.Archive.GetStringSlice("includes")
	includes = cleanPaths(includes)

	if len(includes) == 0 {
		return nil, fmt.Errorf("archive.includes have no config")
	}

	logger.Info("=> includes", len(includes), "rules")

	excludes := tar.model.Archive.GetStringSlice("excludes")
	excludes = cleanPaths(excludes)

	additionalArguments := tar.model.Archive.GetStringSlice("additional_arguments")

	opts = append(opts, "-cP")
	opts = append(opts, additionalArguments...)

	for _, exclude := range excludes {
		opts = append(opts, "--exclude="+filepath.Clean(exclude))
	}

	opts = append(opts, "-f")
	opts = append(opts, archiveFilePath)
	opts = append(opts, includes...)

	return opts, nil
}

func cleanPaths(paths []string) (results []string) {
	for _, p := range paths {
		results = append(results, filepath.Clean(p))
	}
	return
}
