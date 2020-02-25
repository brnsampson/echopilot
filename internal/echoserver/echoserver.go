package echoserver

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	pb "github.com/brnsampson/echopilot/api/echo"
	"github.com/brnsampson/echopilot/pkg/echo"
	"github.com/brnsampson/echopilot/pkg/logger"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"strconv"
	"time"
)

// Application specific functionas and behavior.

// Begin grpc server for echo logic
type echoServer struct {
	pb.UnimplementedEchoServer
	logger logger.Logger
}

func NewEchoServer(logger logger.Logger) *echoServer {
	return &echoServer{logger: logger}
}

func (es *echoServer) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoReply, error) {
	es.logger.Debug("Recieved Echo request: ", ctx)

	begin := time.Now()
	result, err := echo.EchoString(req.Content)
	if err != nil {
		return &pb.EchoReply{Content: ""}, err
	}
	es.logger.Infof("Echo request handled in %+v", time.Since(begin))
	es.logger.Debugf("Replying to Echo request with: %s", result)
	return &pb.EchoReply{Content: result}, nil
}

// End grpc server for echo logic

// Befin File config struct.

type config struct {
	config        string
	grpcAddress   string
	restAddress   string
	tlsCert       string
	tlsKey        string
	tlsSkipVerify bool
}

func (conf config) withMerge(second *config) *config {
	if second.config != "" {
		conf.config = second.config
	}

	if second.grpcAddress != "" {
		conf.grpcAddress = second.grpcAddress
	}

	if second.restAddress != "" {
		conf.restAddress = second.restAddress
	}

	if second.tlsCert != "" {
		conf.tlsCert = second.tlsCert
	}

	if second.tlsKey != "" {
		conf.tlsKey = second.tlsKey
	}

	if second.tlsSkipVerify != false {
		conf.tlsSkipVerify = second.tlsSkipVerify
	}
	return &conf
}

func (conf config) withDefaults(logger logger.Logger) *config {
	if conf.grpcAddress == "" {
		conf.grpcAddress = "127.0.0.1:8080"
		logger.Info("grpcAddress is empty. Defaulting to 127.0.0.1:8080")
	}

	if conf.restAddress == "" {
		conf.restAddress = "127.0.0.1:3000"
		logger.Info("restAddress is empty. Defaulting to 127.0.0.1:3000")
	}

	if conf.tlsCert == "" {
		conf.tlsCert = "/etc/echopilot/cert.pem"
		logger.Info("tlsCert path is empty when loading config. Defaulting to /etc/echopilot/cert.pem")
	}

	if conf.tlsKey == "" {
		conf.tlsKey = "/etc/echopilot/key.pem"
		logger.Info("tlsKey path is empty when loading config. Defaulting to /etc/echopilot/key.pem")
	}
	return &conf
}

func NewFullConfig(logger logger.Logger, flags *pflag.FlagSet) (*config, error) {
	conf, err := NewConfigFromFlags(logger, flags)
	if err != nil {
		logger.Error("Error: could not load config from flags!")
		return conf, err
	}

	if conf.config != "" {
		fileConf, err := NewConfigFromFile(logger, conf.config)
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

	c.config, err = flags.GetString("config")
	if err != nil {
		logger.Debug("Failed to load config file path from flags")
	}

	c.grpcAddress, err = flags.GetString("grpcAddress")
	if err != nil {
		logger.Debug("Failed to load grpcAddress from flags")
	}

	c.restAddress, err = flags.GetString("restAddress")
	if err != nil {
		logger.Debug("Failed to load restAddress from flags")
	}

	c.tlsCert, err = flags.GetString("tlsCert")
	if err != nil {
		logger.Debug("Failed to load tlsCert file path from flags")
	}

	c.tlsKey, err = flags.GetString("tlsKey")
	if err != nil {
		logger.Debug("Failed to load tlsKey file path from flags")
	}

	c.tlsSkipVerify, err = flags.GetBool("tlsSkipVerify")
	if err != nil {
		logger.Debug("Failed to load tlsSkipVerify from flags")
	}

	logger.Infof("Loaded config from flags: %+v", c)

	return &c, nil
}

func NewConfigFromFile(logger logger.Logger, configFile string) (*config, error) {
	var c config

	file, err := os.Open(configFile)
	if err != nil {
		return &c, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)

	if err != nil {
		return &c, err
	}

	logger.Infof("Loaded config from file %s: %+v", configFile, c)

	return &c, nil
}

func NewConfigFromEnv(logger logger.Logger) (*config, error) {
	var c config

	// grpc Addr
	grpcAddr := os.Getenv("ECHO_GRPC_ADDR")
	c.grpcAddress = grpcAddr

	// restAddr
	restAddr := os.Getenv("ECHO_REST_ADDR")
	c.restAddress = restAddr

	// tlsCert and tlsKey
	cert := os.Getenv("ECHO_SERVER_CERT")
	c.tlsCert = cert

	key := os.Getenv("ECHO_SERVER_KEY")
	c.tlsKey = key

	// tlsSkipVerify
	var tlsSkipVerify bool
	skipVerify := os.Getenv("ECHO_GATEWAY_SKIP_VERIFY")
	if skipVerify == "" {
		logger.Info("ECHO_GATEWAY_SKIP_VERIFY is empty. Defaulting to false")
		tlsSkipVerify = false
	} else {
		var err error
		tlsSkipVerify, err = strconv.ParseBool(skipVerify)
		if err != nil {
			logger.Info("ECHO_GATEWAY_SKIP_VERIFY could not be parsed. Try true, false, or empty.")
		}
	}
	c.tlsSkipVerify = tlsSkipVerify

	logger.Debugf("Loaded config from env variables: %+v", c)

	return &c, nil
}

// Generic wrapper for grpc server to be called in server.go
func NewServerRegistrant(logger logger.Logger) (*echoServer, error) {
	return NewEchoServer(logger), nil
}

func (es *echoServer) RegisterWithServer(grpcServer *grpc.Server) error {
	pb.RegisterEchoServer(grpcServer, es)
	return nil
}

//Implement ServerOptions interface
type ServerOpts struct {
	grpcAddr string
	restAddr string
	srvOpts  []grpc.ServerOption
	dialOpts []grpc.DialOption
}

func (so *ServerOpts) GetGrpcAddr() string {
	return so.grpcAddr
}

func (so *ServerOpts) GetRestAddr() string {
	return so.restAddr
}

func (so *ServerOpts) GetSrvOpts() []grpc.ServerOption {
	return so.srvOpts
}

func (so *ServerOpts) GetDialOpts() []grpc.DialOption {
	return so.dialOpts
}

func NewServerOpts(logger logger.Logger, flags *pflag.FlagSet) (*ServerOpts, error) {
	var opts ServerOpts

	conf, err := NewFullConfig(logger, flags)
	if err != nil {
		logger.Error("Error: could not load config!")
		return &opts, err
	}

	logger.Infof("Creating server opts from merged config: %+v", conf)

	opts.grpcAddr = conf.grpcAddress

	opts.restAddr = conf.restAddress

	if conf.tlsCert == "" {
		logger.Error("Error: tlsCert is empty when loading config!")
		return &opts, errors.New("NewServerOptionsFromEnv: tls certificate location must be set")
	}

	if conf.tlsKey == "" {
		logger.Error("Error: tlsKey is empty when loading config!")
		return &opts, errors.New("NewServerOptionsFromEnv: tls key location must be set")
	}

	creds, err := credentials.NewServerTLSFromFile(conf.tlsCert, conf.tlsKey)
	if err != nil {
		logger.Error("Error: Failed to create credentials from cert and key when loading config! Are the tls cert and key paths set properly?")
		return &opts, err
	}
	grpcCreds := grpc.Creds(creds)

	var srvOpts []grpc.ServerOption
	srvOpts = append(srvOpts, grpcCreds)

	opts.srvOpts = srvOpts

	var dialOpts []grpc.DialOption
	if conf.tlsSkipVerify {
		tlsSkipVerifyConfig := &tls.Config{
			InsecureSkipVerify: conf.tlsSkipVerify,
		}
		tlsSkipVerifyOpt := grpc.WithTransportCredentials(credentials.NewTLS(tlsSkipVerifyConfig))
		dialOpts = append(dialOpts, tlsSkipVerifyOpt)
	} else {
		creds, err := credentials.NewClientTLSFromFile(conf.tlsCert, "")
		if err != nil {
			logger.Error("Error: Failed to create client credentials from cert when loading config! Is the tls certificate path set properly?")
			return &opts, err
		}
		clientCreds := grpc.WithTransportCredentials(creds)
		dialOpts = append(dialOpts, clientCreds)
	}

	opts.dialOpts = dialOpts

	return &opts, nil
}
