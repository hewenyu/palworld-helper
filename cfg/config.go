package cfg

import (
	"io"
	"os"

	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

type Config struct {
	RCONSettings       RCONSettings `yaml:"rcon_settings"`
	ArchiveTimeSeconds int32        `yaml:"archive_time_seconds"`
	Interval           string       `yaml:"interval"`
	Timeout            string       `yaml:"timeout"`
}

type RCONSettings struct {
	Endpoint string `yaml:"endpoint"`
	Password string `yaml:"password"`
}

type ConfigManager interface {
	ReadConfig() Config
}

type configStruct struct {
	configFile string
}

// ReadConfig reads the config file and returns a Config struct
func (cs *configStruct) ReadConfig() Config {
	file, err := os.Open(cs.configFile)
	if err != nil {
		slog.Error("Can't open file: ", "ERROR", err)
	}
	defer file.Close()

	reader := io.Reader(file)

	var config Config
	decoder := yaml.NewDecoder(reader)

	if err := decoder.Decode(&config); err != nil {
		slog.Error("Can't decode file: ", "ERROR", err)
		os.Exit(1)
	}

	return config
}

// NewConfigManager returns a ConfigManager interface
func NewConfigManager() ConfigManager {
	return &configStruct{
		configFile: configFile,
	}
}
