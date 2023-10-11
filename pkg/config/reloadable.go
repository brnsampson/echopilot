package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/brnsampson/echopilot/pkg/option"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
    "github.com/charmbracelet/log"
)

const DEFAULT_CONFIG_FILE = "/etc/echopilot/echopilot.json"
const DEFAULT_HOST = "localhost"
const DEFAULT_IP = "127.0.0.1"
const DEFAULT_PORT = 3000
const DEFAULT_TLS_ENABLED = true
const DEFAULT_TLS_SKIP_VERIFY = false
const DEFAULT_TLS_CERT = "/etc/echopilot/tls/cert.pem"
const DEFAULT_TLS_KEY = "/etc/echopilot/tls/key.pem"

type StaticConfig struct {
	ConfigFile         string
	Host               string
	IP                 string
	Port               int
	TlsCert            string
	TlsKey             string
	TlsEnabled         bool
	TlsSkipVerify      bool
}


// Generic server configuration which can be reloaded on demand.
type ReloadableConfig struct {
	ConfigFile         option.Option[string] `json:"configFile" env:"ECHOPILOT_CONFIG_FILE"`
	Host               option.Option[string] `json:"serverHost" env:"ECHOPILOT_HOST"`
	IP                 option.Option[string] `json:"bindHost" env:"ECHOPILOT_BIND_IP"`
	Port               option.Option[int]    `json:"serverPort" env:"ECHOPILOT_PORT"`
	TlsCert            option.Option[string] `json:"tlsCert" env:"ECHOPILOT_TLS_CERT"`
	TlsKey             option.Option[string] `json:"tlsKey" env:"ECHOPILOT_TLS_KEY"`
	TlsEnabled         option.Option[bool]   `json:"tlsEnabled" env:"ECHOPILOT_TLS_ENABLED"`
	TlsSkipVerify      option.Option[bool]   `json:"tlsSkipVerify" env:"ECHOPILOT_TLS_SKIP_VERIFY"`
}

func emptyReloadableConfig() ReloadableConfig {
    return ReloadableConfig {
        ConfigFile: option.None[string](),
        Host: option.None[string](),
        IP: option.None[string](),
        Port: option.None[int](),
        TlsCert: option.None[string](),
        TlsKey: option.None[string](),
        TlsEnabled: option.None[bool](),
        TlsSkipVerify: option.None[bool](),
    }
}

func (r *ReloadableConfig) Finalize() StaticConfig {
    f := func() string {
	    // We only want to default to using a config file if the default file exists
	    if _, err := os.Stat(DEFAULT_CONFIG_FILE); err == nil {
	    	// The file exists, so we default to using it.
	    	return DEFAULT_CONFIG_FILE
	    } else {
            return ""
        }
    }
    configFile := r.ConfigFile.UnwrapOrElse(f)
    host := r.Host.UnwrapOrDefault(DEFAULT_HOST)
    ip := r.IP.UnwrapOrDefault(DEFAULT_IP)
    port := r.Port.UnwrapOrDefault(DEFAULT_PORT)
    tlsCert := r.TlsCert.UnwrapOrDefault(DEFAULT_TLS_CERT)
    tlsKey := r.TlsKey.UnwrapOrDefault(DEFAULT_TLS_KEY)
    tlsEnabled := r.TlsEnabled.UnwrapOrDefault(DEFAULT_TLS_ENABLED)
    tlsSkipVerify := r.TlsSkipVerify.UnwrapOrDefault(DEFAULT_TLS_SKIP_VERIFY)

    conf := StaticConfig {
        ConfigFile: configFile,
        Host: host,
        IP: ip,
        Port: port,
        TlsCert: tlsCert,
        TlsKey: tlsKey,
        TlsEnabled: tlsEnabled,
        TlsSkipVerify: tlsSkipVerify,
    }

	return conf
}

func (conf ReloadableConfig) withMerge(second ReloadableConfig) ReloadableConfig {
	if second.ConfigFile.IsSome() {
		conf.ConfigFile = second.ConfigFile
	}

	if second.Host.IsSome() {
		conf.Host = second.Host
	}

	if second.IP.IsSome() {
		conf.IP = second.IP
	}

	if second.Port.IsSome() {
		conf.Port = second.Port
	}

	if second.TlsCert.IsSome() {
		conf.TlsCert = second.TlsCert
	}

	if second.TlsKey.IsSome() {
		conf.TlsKey = second.TlsKey
	}

	if second.TlsEnabled.IsSome() {
		conf.TlsEnabled = second.TlsEnabled
	}

	if second.TlsSkipVerify.IsSome() {
		conf.TlsSkipVerify = second.TlsSkipVerify
	}

	return conf
}

func NewFullReloadableConfig(flags *pflag.FlagSet) (*ReloadableConfig, error) {
	conf, err := NewReloadableConfigFromFlags(flags)
	if err != nil {
		log.Error("Error: could not load config from flags!")
		return &conf, err
	}

	envConf, err := NewReloadableConfigFromEnv()
	if err != nil {
		log.Error("Error: could not load config from environment!")
		return &conf, err
	}

	conf = conf.withMerge(envConf)

	if conf.ConfigFile.IsSome() {
        tmp := conf.ConfigFile.Clone()
        file := (&tmp).UnsafeUnwrap()
		fileConf, err := NewReloadableConfigFromFile(file)
		if err != nil {
			log.Error("Error loadng config from file", "filename", conf.ConfigFile, "error", err)
		} else {
			conf = conf.withMerge(fileConf)
		}
	}

	log.Debug("Loaded combines config from all sources", "config", conf)

	return &conf, nil
}

func NewReloadableConfigFromFlags(flags *pflag.FlagSet) (ReloadableConfig, error) {
    var configFile option.Option[string]
    tmp, err := flags.GetString("config")
	if err != nil || tmp == "" {
        configFile = option.None[string]()
		log.Debug("Failed to load config file path from flags")
	} else {
        log.Debug("Found config file to load config file path from flags", "filename", tmp)
        configFile = option.NewOption(tmp)
    }

    var host option.Option[string]
    tmp, err = flags.GetString("host")
	if err != nil || tmp == "" {
        host = option.None[string]()
		log.Debug("Failed to load host from flags")
	} else {
        host = option.NewOption(tmp)
    }

    var ip option.Option[string]
    tmp, err = flags.GetString("ip")
	if err != nil || tmp == "" {
        ip = option.None[string]()
		log.Debug("Failed to load ip from flags")
	} else {
        ip = option.NewOption(tmp)
    }

    var port option.Option[int]
    tmpint, err := flags.GetInt("port")
	if err != nil {
	}
	if err != nil {
        port = option.None[int]()
		log.Debug("Failed to load port from flags")
	} else {
        port = option.NewOption(tmpint)
    }

    var tlsCert option.Option[string]
    tmp, err = flags.GetString("tlsCert")
	if err != nil || tmp == "" {
        tlsCert = option.None[string]()
		log.Debug("Failed to load TlsCert file path from flags")
	} else {
        tlsCert = option.NewOption(tmp)
    }

    var tlsKey option.Option[string]
    tmp, err = flags.GetString("tlsKey")
	if err != nil || tmp == "" {
        tlsKey = option.None[string]()
		log.Debug("Failed to load TlsKey file path from flags")
	} else {
        tlsKey = option.NewOption(tmp)
    }

    var tlsEnabled option.Option[bool]
    tmpbool, err := flags.GetBool("tlsEnabled")
	if err != nil {
        tlsEnabled = option.None[bool]()
		log.Debug("Failed to load TlsEnabled from flags")
	} else {
        tlsEnabled = option.NewOption(tmpbool)
    }

    var tlsSkipVerify option.Option[bool]
    tmpbool, err = flags.GetBool("tlsSkipVerify")
	if err != nil {
        tlsSkipVerify = option.None[bool]()
		log.Debug("Failed to load TlsSkipVerify from flags")
        tlsSkipVerify = option.NewOption(tmpbool)
	} else {
    }

    c := ReloadableConfig{
        ConfigFile: configFile,
        Host: host,
        IP: ip,
        Port: port,
        TlsCert: tlsCert,
        TlsKey: tlsKey,
        TlsEnabled: tlsEnabled,
        TlsSkipVerify: tlsSkipVerify,
    }

	log.Info("Loaded config from flags", "config", c)

	return c, nil
}

func NewReloadableConfigFromFile(ConfigFile string) (ReloadableConfig, error) {
    c := emptyReloadableConfig()

	if _, err := os.Stat(ConfigFile); errors.Is(os.ErrNotExist, err) {
		return c, err
	}

	file, err := os.Open(ConfigFile)
	if err != nil {
		return c, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)

	if err != nil {
		return c, err
	}

	log.Info("Loaded config from file", "filename", ConfigFile, "config", c)

	return c, nil
}

func NewReloadableConfigFromEnv() (ReloadableConfig, error) {
	c := emptyReloadableConfig()

	if err := env.Parse(&c); err != nil {
		log.Error("Failed to load config from env variables!")
		return c, err
	}

	log.Debug("Loaded config from env variables", "config", c)

	return c, nil
}
