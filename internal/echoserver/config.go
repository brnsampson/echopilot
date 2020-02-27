package echoserver

import (
	"encoding/json"
	"github.com/brnsampson/echopilot/pkg/logger"
	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
	"os"
)

// Application specific config
type config struct {
	ConfigFile    string `json:"ConfigFile" env:"ECHO_CONFIG_FILE"`
	GrpcAddress   string `json:"grpcAddress" env:"ECHO_GRPC_ADDR"`
	RestAddress   string `json:"restAddress" env:"ECHO_REST_ADDR"`
	TlsCert       string `json:"tlsCert" env:"ECHO_TLS_CERT"`
	TlsKey        string `json:"tlsKey" env:"ECHO_TLS_KEY"`
	TlsSkipVerify bool   `json:"tlsSkipVerify" env:"ECHO_GATEWAY_TLS_SKIP_VERIFY"`
}

func (conf config) withMerge(second *config) *config {
	if second.ConfigFile != "" {
		conf.ConfigFile = second.ConfigFile
	}

	if second.GrpcAddress != "" {
		conf.GrpcAddress = second.GrpcAddress
	}

	if second.RestAddress != "" {
		conf.RestAddress = second.RestAddress
	}

	if second.TlsCert != "" {
		conf.TlsCert = second.TlsCert
	}

	if second.TlsKey != "" {
		conf.TlsKey = second.TlsKey
	}

	if second.TlsSkipVerify != false {
		conf.TlsSkipVerify = second.TlsSkipVerify
	}
	return &conf
}

func (conf config) withDefaults(logger logger.Logger) *config {
	if conf.GrpcAddress == "" {
		conf.GrpcAddress = "127.0.0.1:8080"
		logger.Info("GrpcAddress is empty. Defaulting to 127.0.0.1:8080")
	}

	if conf.RestAddress == "" {
		conf.RestAddress = "127.0.0.1:3000"
		logger.Info("RestAddress is empty. Defaulting to 127.0.0.1:3000")
	}

	if conf.TlsCert == "" {
		conf.TlsCert = "/etc/echopilot/cert.pem"
		logger.Info("TlsCert path is empty when loading config. Defaulting to /etc/echopilot/cert.pem")
	}

	if conf.TlsKey == "" {
		conf.TlsKey = "/etc/echopilot/key.pem"
		logger.Info("TlsKey path is empty when loading config. Defaulting to /etc/echopilot/key.pem")
	}
	return &conf
}

func NewFullConfig(logger logger.Logger, flags *pflag.FlagSet) (*config, error) {
	conf, err := NewConfigFromFlags(logger, flags)
	if err != nil {
		logger.Error("Error: could not load config from flags!")
		return conf, err
	}

	if conf.ConfigFile != "" {
		fileConf, err := NewConfigFromFile(logger, conf.ConfigFile)
		if err != nil {
			logger.Error("Error: could not load config from file!")
			return conf, err
		}
		conf = conf.withMerge(fileConf)
	}

	envConf, err := NewConfigFromEnv(logger)
	if err != nil {
		logger.Error("Error: could not load config from environment!")
		return conf, err
	}

	conf = conf.withMerge(envConf)
	conf = conf.withDefaults(logger)

	logger.Debugf("Loaded combines config from all sources: %+v", conf)

	return conf, nil
}

func NewConfigFromFlags(logger logger.Logger, flags *pflag.FlagSet) (*config, error) {
	var c config
	var err error

	c.ConfigFile, err = flags.GetString("config")
	if err != nil {
		logger.Debug("Failed to load config file path from flags")
	}

	c.GrpcAddress, err = flags.GetString("GrpcAddress")
	if err != nil {
		logger.Debug("Failed to load GrpcAddress from flags")
	}

	c.RestAddress, err = flags.GetString("RestAddress")
	if err != nil {
		logger.Debug("Failed to load RestAddress from flags")
	}

	c.TlsCert, err = flags.GetString("TlsCert")
	if err != nil {
		logger.Debug("Failed to load TlsCert file path from flags")
	}

	c.TlsKey, err = flags.GetString("TlsKey")
	if err != nil {
		logger.Debug("Failed to load TlsKey file path from flags")
	}

	c.TlsSkipVerify, err = flags.GetBool("TlsSkipVerify")
	if err != nil {
		logger.Debug("Failed to load TlsSkipVerify from flags")
	}

	logger.Infof("Loaded config from flags: %+v", c)

	return &c, nil
}

func NewConfigFromFile(logger logger.Logger, ConfigFile string) (*config, error) {
	var c config

	file, err := os.Open(ConfigFile)
	if err != nil {
		return &c, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)

	if err != nil {
		return &c, err
	}

	logger.Infof("Loaded config from file %s: %+v", ConfigFile, c)

	return &c, nil
}

func NewConfigFromEnv(logger logger.Logger) (*config, error) {
	c := config{}

	if err := env.Parse(&c); err != nil {
		logger.Error("Failed to load config from env variables!")
		return &c, err
	}

	logger.Debugf("Loaded config from env variables: %+v", c)

	return &c, nil
}
