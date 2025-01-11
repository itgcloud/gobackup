package archive

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/itgcloud/gobackup/config"
	"github.com/itgcloud/gobackup/helper"
	"github.com/itgcloud/gobackup/logger"
)

// Run archive
func Run(model config.ModelConfig) error {
	logger := logger.Tag("Archive")

	if model.Archive == nil {
		return nil
	}

	if err := helper.MkdirP(model.DumpPath); err != nil {
		logger.Errorf("Failed to mkdir dump path %s: %v", model.DumpPath, err)

		return err
	}

	// Archive + compress with tar in one step if compression is enabled and databases are not empty
	if model.CompressWith.Type != "" && len(model.Databases) == 0 {
		return nil
	}

	opts, err := options(model)
	if err != nil {
		return err
	}

	if _, err = helper.Exec("tar", opts...); err != nil {
		return err
	}

	return nil
}

func options(model config.ModelConfig) ([]string, error) {
	var opts []string

	includes := model.Archive.GetStringSlice("includes")
	includes = cleanPaths(includes)

	if len(includes) == 0 {
		return nil, fmt.Errorf("archive.includes have no config")
	}

	logger.Info("=> includes", len(includes), "rules")

	excludes := model.Archive.GetStringSlice("excludes")
	excludes = cleanPaths(excludes)

	tarPath := path.Join(model.DumpPath, "archive.tar")
	if helper.IsGnuTar {
		opts = append(opts, "--ignore-failed-read")
	}

	additionalArguments := model.Archive.GetStringSlice("additional_arguments")
	opts = append(opts, "-cP")
	opts = append(opts, additionalArguments...)
	opts = append(opts, "-f")
	opts = append(opts, tarPath)

	for _, exclude := range excludes {
		opts = append(opts, "--exclude="+filepath.Clean(exclude))
	}

	opts = append(opts, includes...)

	return opts, nil
}

func cleanPaths(paths []string) []string {
	var results []string

	for _, p := range paths {
		results = append(results, filepath.Clean(p))
	}

	return results
}
