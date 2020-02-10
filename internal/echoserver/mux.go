package echoserver

import (
	"context"
	gw "github.com/brnsampson/echopilot/api/echo"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
)

// Begin REST gatwway mux implementation for echo logic.

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "dist/swagger-ui/swagger.json")
}

func serveUI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ui" {
		errorHandler(w, r, 404)
		return
	}
	http.ServeFile(w, r, "dist/index.html")
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		filePath := "dist/404.html"
		http.ServeFile(w, r, filePath)
	}
}

func NewGatewayMux(ctx context.Context, grpcAddr string, opts []grpc.DialOption) (http.Handler, error) {
	rmux := runtime.NewServeMux()
	err := gw.RegisterEchoHandlerFromEndpoint(ctx, rmux, grpcAddr, opts)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.HandleFunc("/swagger", serveSwagger)
	fs := http.FileServer(http.Dir("dist/swagger-ui"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))
	mux.HandleFunc("/ui", serveUI)

	return mux, nil
}

// End REST gateway mux implementation for echo logic.
