package router

import (
    "net/http"
)

func All[T any](ts []T, pred func(T) bool) bool {
    for _, t := range ts {
        if !pred(t) {
            return false
        }
    }
    return true
}

func DoAll[T any](ts []T, pred func(T) error) error {
    var err error
    for _, t := range ts {
        if e := pred(t); e != nil {
            err = e
        }
    }
    return err
}

// Common server implementation boilerplate
type HandlerGenerator interface {
    // If a HandlerGenerator does not run a goroutine, then Run() and Halt() are no-ops and IsRunning and IsHalted() always return true
    Run() error // Run() should be idempotent
    Halt() error
    IsRunning() bool
    IsHalted() bool
	GetHttpHandlers() []*HandlerPair
}

type HandlerPair struct {
    path string
    handler http.Handler
}

func NewHandlerPair(path string, handler http.Handler) *HandlerPair {
    return &HandlerPair {path, handler}
}

type HandlerFuncPair struct {
    path string
    handler http.HandlerFunc
}

func NewRouter() *Router {
    handlers := make([]HandlerPair, 0)
    handlerFuncs := make([]HandlerFuncPair, 0)
    handlerGens := make([]HandlerGenerator, 0)
    return &Router { handlers, handlerFuncs, handlerGens }
}

type Router struct {
    handlers []HandlerPair
    handlerFuncs []HandlerFuncPair
    handlerGens []HandlerGenerator
}

func (r *Router) AddHandler(path string, handler http.Handler) {
    r.handlers = append(r.handlers, HandlerPair{path, handler}) }

func (r *Router) AddHandlerFunc(path string, handler http.HandlerFunc) {
    r.handlerFuncs = append(r.handlerFuncs, HandlerFuncPair{path, handler})
}

func (r *Router) AddHandlerGen(handleGen HandlerGenerator) {
    r.handlerGens = append(r.handlerGens, handleGen)
}

// implement HandlerGenerator for Router
func (r *Router) Run() error{
    return DoAll(r.handlerGens, func(h HandlerGenerator) error { return h.Run() })
}

func (r *Router) Halt() error {
    return DoAll(r.handlerGens, func(h HandlerGenerator) error { return h.Halt() })
}

func (r *Router) IsRunning() bool {
    return All(r.handlerGens, func(h HandlerGenerator) bool { return h.IsRunning() })
}

func (r *Router) IsHalted() bool {
    return All(r.handlerGens, func(h HandlerGenerator) bool { return h.IsHalted() })
}

func (r *Router) GetHttpHandler() (path string, handler http.Handler) {
    mux := http.NewServeMux()
    for _, handlerPair := range r.handlers {
        mux.Handle(handlerPair.path, handlerPair.handler)
    }

    for _, handlerFuncPair := range r.handlerFuncs {
        mux.Handle(handlerFuncPair.path, handlerFuncPair.handler)
    }

    for _, handlerGen := range r.handlerGens {
        for _, handlerpair := range handlerGen.GetHttpHandlers() {
            mux.Handle(handlerpair.path, handlerpair.handler)
        }
    }
    return "/", mux
}
// End implementation of HandlerGenerator for Router
