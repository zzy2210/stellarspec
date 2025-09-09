package config

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

type BaseConfig struct {
	APIServer string
	Model     string
	Key       string
	Language  string
}

func LoadFile(path string) (*BaseConfig, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("load config file failed: err= %v", err)
	}
	config := &BaseConfig{}
	// 读取配置值
	config.APIServer = cfg.Section("").Key("APIServer").String()
	config.Model = cfg.Section("").Key("Model").String()
	config.Key = cfg.Section("").Key("Key").String()
	
	// 读取语言配置，默认为 zh
	config.Language = cfg.Section("").Key("Language").String()
	if config.Language == "" {
		config.Language = "zh"
	}

	return config, nil

}

func ensureConfigFile(path string) error {
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件不存在，创建一个空的 INI 文件
		cfg := ini.Empty()
		if err := cfg.SaveTo(path); err != nil {
			return fmt.Errorf("create config file failed: err= %v", err)
		}
	}
	return nil
}

func SaveAPIServer(apiserver string, path string) error {
	if err := ensureConfigFile(path); err != nil {
		return fmt.Errorf("ensure config path failed: err= %v", err)
	}
	cfg, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("load config file failed: err= %v", err)
	}

	cfg.Section("").Key("APIServer").SetValue(apiserver)

	if err := cfg.SaveTo(path); err != nil {
		return err
	}

	return nil
}

func SaveKey(key string, path string) error {
	if err := ensureConfigFile(path); err != nil {
		return fmt.Errorf("ensure config path failed: err= %v", err)
	}
	cfg, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("load config file failed: err= %v", err)
	}

	cfg.Section("").Key("Key").SetValue(key)

	if err := cfg.SaveTo(path); err != nil {
		return err
	}

	return nil
}

func SaveModel(model string, path string) error {
	if err := ensureConfigFile(path); err != nil {
		return fmt.Errorf("ensure config path failed: err= %v", err)
	}
	cfg, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("load config file failed: err= %v", err)
	}

	cfg.Section("").Key("Model").SetValue(model)

	if err := cfg.SaveTo(path); err != nil {
		return err
	}

	return nil
}

func SaveLanguage(language string, path string) error {
	if err := ensureConfigFile(path); err != nil {
		return fmt.Errorf("ensure config path failed: err= %v", err)
	}
	cfg, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("load config file failed: err= %v", err)
	}

	cfg.Section("").Key("Language").SetValue(language)

	if err := cfg.SaveTo(path); err != nil {
		return err
	}

	return nil
}
