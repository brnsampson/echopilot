package appserver

// Implements a custom server for echoing requests back to the user.
//
// There are a few key components to this:
//
// EchoConnectServer: Implements the connect protocol RPC functions
// we defined in our api (protobuf) files. This should be mostly
// boilerplate middleware and request handling code and all domain
// logic should be implemented in an EchoService pkg.
//
// EchoServer: an encapsulating struct that creates a new server
// using the EchoConnectServer and generic SignaledServer pkg.

import (
    "os"
	"github.com/brnsampson/echopilot/rpc/echo"
	"github.com/brnsampson/echopilot/features/memory"
	"github.com/brnsampson/echopilot/pkg/server"
	"github.com/brnsampson/echopilot/pkg/config"
	"github.com/spf13/pflag"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/charmbracelet/log"
)

// Struct used as primary entrypoint for an RPC based interface.
func NewAppServer(flags *pflag.FlagSet) (*AppServer, error) {
    logger := log.NewWithOptions(os.Stderr, log.Options{
        ReportCaller: true,
        ReportTimestamp: true,
        Prefix: "echopilot",
        Level: log.DebugLevel,
    })

    conf, err := config.NewServerConfig(flags)
    if err != nil {
        return nil, err
    }

    echoService := echo.NewService(logger)
    memoryFeature := memory.NewFeature(logger)
    //router := router.NewRouter()
    router := chi.NewRouter()
    router.Use(middleware.Logger)
    router.Use(middleware.Recoverer)

    router.Route("/", routeRoot)
    router.Mount(memoryFeature.GetHandler())
    router.Mount(echoService.GetHandler())
    //router.AddHandlerFunc("/", serveEchoComponents)

	return &AppServer{
		router,
		server.NewServer(logger),
		logger,
		conf,
	}, nil
}

type AppServer struct {
    router *chi.Mux
	server *server.Server
	logger *log.Logger
	config *config.ServerConfig
}

func (es *AppServer) Run() {
	es.server.Run(es.router, es.config)
}

func (es *AppServer) BlockingRun() int {
	return es.server.BlockingRun(es.router, es.config)
}
