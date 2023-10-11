package config

import (
	"crypto/tls"
    "strings"
    "strconv"

	"github.com/spf13/pflag"
    "github.com/charmbracelet/log"
)

func NewServerConfig(flags *pflag.FlagSet) (*ServerConfig, error) {
    tlsConf := tls.Config{
		Certificates:             []tls.Certificate{},
		MinVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	conf := ServerConfig{
		flags:  flags,
        tlsConf: &tlsConf,
	}
    if err := conf.update(); err != nil {
        return nil, err
    }
	return &conf, nil
}

type ServerConfig struct {
	flags   *pflag.FlagSet
	config  *StaticConfig
	tlsConf *tls.Config
}

func (c *ServerConfig) update() error {
	conf, err := NewFullReloadableConfig(c.flags)
	if err != nil {
		log.Error("Could not update echo server config due to error loading", "error", err)
		return err
	}

    staticConf := conf.Finalize()
	log.Info("Updating echo server config from merged config", "config", staticConf)
	c.config = &staticConf

	if staticConf.TlsEnabled {
        cert, err := tls.LoadX509KeyPair(staticConf.TlsCert, staticConf.TlsKey)
		if err != nil {
			log.Error("Updating echo server TLS Certificate failed", "error", err)
			return err
		}

		c.tlsConf = &tls.Config{
			Certificates:             []tls.Certificate{cert},
			MinVersion:               tls.VersionTLS13,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
	}

	return nil
}

func (c *ServerConfig) GetAddr(update bool) (string, error) {
	if update {
		if err := c.update(); err != nil {
			return "", err
		}
	}

    addr := strings.Join([]string{c.config.IP, strconv.Itoa(c.config.Port)}, ":")
	return addr, nil
}

func (c *ServerConfig) GetHost(update bool) (string, error) {
	if update {
		if err := c.update(); err != nil {
			return "", err
		}
	}

	return c.config.Host, nil
}

func (c *ServerConfig) GetTlsConfig(update bool) (*tls.Config, error) {
	if update {
		if err := c.update(); err != nil {
			return c.tlsConf, err
		}
	}

	return c.tlsConf, nil
}

func (esc *ServerConfig) GetTlsEnabled(update bool) (bool, error) {
	if update {
		if err := esc.update(); err != nil {
			return esc.config.TlsEnabled, err
		}
	}

	return esc.config.TlsEnabled, nil
}

