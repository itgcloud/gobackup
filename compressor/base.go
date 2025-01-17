package compressor

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"

	"github.com/itgcloud/gobackup/config"
	"github.com/itgcloud/gobackup/logger"
)

// Base compressor
type Base struct {
	name            string
	ext             string
	parallelProgram string
	model           config.ModelConfig
	viper           *viper.Viper
}

// Compressor
type Compressor interface {
	perform() (archivePath string, err error)
}

func (c *Base) archiveFilePath(ext string) string {
	return filepath.Join(c.model.TempPath, c.model.Name+"-"+time.Now().Format("2006-01-02-15-04-05")+ext)
}

// Run compressor, return archive path
func Run(model config.ModelConfig) (string, error) {
	base, err := newBase(model)
	if err != nil {
		return "", err
	}

	c := &Tar{base}

	logger := logger.Tag("Compressor")
	logger.Info("=> Compress | " + model.CompressWith.Type)

	if err := changeWorkDir(model); err != nil {
		return "", err
	}

	archivePath, err := c.perform()
	if err != nil {
		return "", err
	}

	logger.Info("->", archivePath)

	return archivePath, nil
}

func newBase(model config.ModelConfig) (*Base, error) {
	base := &Base{
		name:  model.Name,
		model: model,
		viper: model.CompressWith.Viper,
	}

	var ext, parallelProgram string
	switch model.CompressWith.Type {
	case "gz", "tgz", "taz", "tar.gz":
		ext = ".tar.gz"
		parallelProgram = "pigz"
	case "Z", "taZ", "tar.Z":
		ext = ".tar.Z"
	case "bz2", "tbz", "tbz2", "tar.bz2":
		ext = ".tar.bz2"
		parallelProgram = "pbzip2"
	case "lz", "tar.lz":
		ext = ".tar.lz"
	case "lzma", "tlz", "tar.lzma":
		ext = ".tar.lzma"
	case "lzo", "tar.lzo":
		ext = ".tar.lzo"
	case "xz", "txz", "tar.xz":
		ext = ".tar.xz"
		parallelProgram = "pixz"
	case "zst", "tzst", "tar.zst":
		ext = ".tar.zst"
	case "tar":
		ext = ".tar"
	case "":
		ext = ".tar"
		model.CompressWith.Type = "tar"
	default:
		return nil, fmt.Errorf("Unsupported compress type: %s", model.CompressWith.Type)
	}

	// save Extension
	model.Viper.Set("Ext", ext)

	base.ext = ext
	base.parallelProgram = parallelProgram

	return base, nil
}

func changeWorkDir(model config.ModelConfig) error {
	if len(model.Databases) == 0 {
		return nil
	}

	if err := os.Chdir(filepath.Join(model.DumpPath, "../")); err != nil {
		return fmt.Errorf("chdir to dump path: %s: %w", model.DumpPath, err)
	}

	return nil
}
