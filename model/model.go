package model

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/itgcloud/gobackup/archive"
	"github.com/itgcloud/gobackup/compressor"
	"github.com/itgcloud/gobackup/config"
	"github.com/itgcloud/gobackup/database"
	"github.com/itgcloud/gobackup/encryptor"
	"github.com/itgcloud/gobackup/helper"
	"github.com/itgcloud/gobackup/logger"
	"github.com/itgcloud/gobackup/notifier"
	"github.com/itgcloud/gobackup/splitter"
	"github.com/itgcloud/gobackup/storage"
)

// Model class
type Model struct {
	Config config.ModelConfig
}

// Perform model
func (m Model) Perform() (err error) {
	logger := logger.Tag(fmt.Sprintf("Model: %s", m.Config.Name))

	m.before()

	defer func() {
		if err != nil {
			logger.Error(err)
			notifier.Failure(m.Config, err.Error())
		} else {
			notifier.Success(m.Config)
		}
	}()

	logger.Info("WorkDir:", m.Config.DumpPath)

	defer func() {
		if r := recover(); r != nil {
			m.after()
			logger.Fatalf("PANIC: %v", r)
		}

		m.after()
	}()

	if err = database.Run(m.Config); err != nil {
		return
	}

	if err = archive.Run(m.Config); err != nil {
		return
	}

	// It always to use compressor, default use tar, even not enable compress.
	archivePath, err := compressor.Run(m.Config)
	if err != nil {
		return
	}

	archivePath, err = encryptor.Run(archivePath, m.Config)
	if err != nil {
		return
	}

	archivePath, err = splitter.Run(archivePath, m.Config)
	if err != nil {
		return
	}

	err = storage.Run(m.Config, archivePath)
	if err != nil {
		return
	}

	return nil
}

func (m Model) before() {
	// Execute before_script
	if len(m.Config.BeforeScript) == 0 {
		return
	}

	logger.Info("Executing before_script...")
	if _, err := helper.ExecWithStdio(m.Config.BeforeScript, true); err != nil {
		logger.Error(err)
	}

	return
}

// Cleanup model temp files
func (m Model) after() {
	logger := logger.Tag("Model")

	tempDir := m.Config.TempPath
	if viper.GetBool("useTempWorkDir") {
		tempDir = viper.GetString("workdir")
	}
	logger.Infof("Cleanup temp: %s/", tempDir)
	if err := os.RemoveAll(tempDir); err != nil {
		logger.Errorf("Cleanup temp dir %s error: %v", tempDir, err)
	}

	if len(m.Config.AfterScript) == 0 {
		return
	}

	logger.Info("Executing after_script...")
	if _, err := helper.ExecWithStdio(m.Config.AfterScript, true); err != nil {
		logger.Error(err)
	}

	return
}

// GetModelByName get model by name
func GetModelByName(name string) *Model {
	modelConfig := config.GetModelConfigByName(name)
	if modelConfig == nil {
		return nil
	}

	return &Model{
		Config: *modelConfig,
	}
}

// GetModels get models
func GetModels() []*Model {
	models := make([]*Model, 0, len(config.Models))

	for _, modelConfig := range config.Models {
		m := Model{
			Config: modelConfig,
		}
		models = append(models, &m)
	}

	return models
}
