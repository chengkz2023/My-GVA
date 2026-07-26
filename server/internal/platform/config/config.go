package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	legacyconfig "github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	ConfigEnv         = "GVA_CONFIG"
	ConfigDefaultFile = "config.yaml"
	ConfigTestFile    = "config.test.yaml"
	ConfigDebugFile   = "config.debug.yaml"
	ConfigReleaseFile = "config.release.yaml"
)

type Config = legacyconfig.Server

func Load() (*viper.Viper, Config) {
	configPath := getConfigPath()

	var cfg Config
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		_ = v.Unmarshal(&cfg)
	})
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %w", err))
	}

	return v, cfg
}

func getConfigPath() string {
	var config string
	flag.StringVar(&config, "c", "", "choose config file.")
	flag.Parse()
	if config != "" {
		fmt.Printf("using command-line config: %s\n", config)
		return config
	}
	if env := os.Getenv(ConfigEnv); env != "" {
		fmt.Printf("using env %s config: %s\n", ConfigEnv, env)
		return env
	}

	switch gin.Mode() {
	case gin.DebugMode:
		config = ConfigDebugFile
	case gin.ReleaseMode:
		config = ConfigReleaseFile
	case gin.TestMode:
		config = ConfigTestFile
	default:
		config = ConfigDefaultFile
	}
	fmt.Printf("gin mode %s, config: %s\n", gin.Mode(), config)

	// Check configs/ directory first, then root directory
	configsPath := filepath.Join("configs", config)
	if _, err := os.Stat(configsPath); err == nil {
		fmt.Printf("using configs/ directory: %s\n", configsPath)
		return configsPath
	}
	if _, err := os.Stat(config); err == nil {
		return config
	}

	// Fallback: try configs/config.yaml, then config.yaml
	configsDefault := filepath.Join("configs", ConfigDefaultFile)
	if _, err := os.Stat(configsDefault); err == nil {
		fmt.Printf("fallback to configs/ default: %s\n", configsDefault)
		return configsDefault
	}
	fmt.Printf("using default config: %s\n", ConfigDefaultFile)
	return ConfigDefaultFile
}
