package config

import (
	"flag"
	"fmt"
	"os"

	legacyconfig "github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
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

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		_ = v.Unmarshal(&global.GVA_CONFIG)
	})
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %w", err))
	}

	return v, global.GVA_CONFIG
}

func Current() Config {
	return global.GVA_CONFIG
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
	}
	fmt.Printf("gin mode %s, config: %s\n", gin.Mode(), config)

	if _, err := os.Stat(config); err != nil || os.IsNotExist(err) {
		config = ConfigDefaultFile
		fmt.Printf("config not found, using default: %s\n", config)
	}
	return config
}
