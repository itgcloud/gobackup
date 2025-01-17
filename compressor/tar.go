package compressor

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/itgcloud/gobackup/helper"
	"github.com/itgcloud/gobackup/logger"
)

type Tar struct {
	*Base
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

	if err := tar.checkIncludes(); err != nil {
		return nil, err
	}

	opts = append(opts, "-cP")
	opts = append(opts, tar.additionalArgs()...)
	opts = append(opts, tar.excludesArgs()...)

	opts = append(opts, "-f")
	opts = append(opts, archiveFilePath)
	opts = append(opts, tar.includesArgs()...)

	return opts, nil
}

func (tar *Tar) additionalArgs() []string {
	if tar.model.Archive == nil {
		return []string{}
	}

	return tar.model.Archive.GetStringSlice("additional_arguments")
}

func (tar *Tar) excludesArgs() []string {
	if tar.model.Archive == nil {
		return []string{}
	}

	excludes := tar.model.Archive.GetStringSlice("excludes")
	excludes = cleanPaths(excludes)

	var args []string
	for _, exclude := range excludes {
		args = append(args, "--exclude="+filepath.Clean(exclude))
	}

	return args
}

func (tar *Tar) checkIncludes() error {
	if tar.model.Databases == nil {
		if len(tar.model.Archive.GetStringSlice("includes")) == 0 {
			return fmt.Errorf("archive.includes have no config")
		}
	}

	return nil
}

func (tar *Tar) includesArgs() []string {
	logger := logger.Tag("Compressor")
	var includes []string

	if tar.model.Archive == nil && tar.model.Databases != nil {
		includes = []string{tar.model.DumpPath}
		includes = cleanPaths(includes)

		return includes
	}

	includes = tar.model.Archive.GetStringSlice("includes")
	includes = cleanPaths(includes)

	logger.Info("=> includes", len(includes), "rules")

	return includes
}

func cleanPaths(paths []string) (results []string) {
	for _, p := range paths {
		results = append(results, filepath.Clean(p))
	}
	return
}
