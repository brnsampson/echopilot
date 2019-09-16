package echoserver

import (
	"net/http"
	"os/signal"
	"sync"
	"time"
)

type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Sync() error
}

func Echo(s string) string {
	return s
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	resp := Echo(body)
	io.WriteString(w, resp)
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request), logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		fn(w, r)
		logger.Debugf("Served request for %s in %d", r.URL.Path, time.Since(begin)/time.Second)
	}
}

func ServeWithReload(logger Logger, waitgroup *sync.WaitGroup, err chan<- error, done <-chan struct{}, hups <-chan os.Signal) {
	m := http.NewServerMux()

	for {
		m.HandleFunc("/", MakeHandler(EchoHandler))
		addr := "127.0.0.1:9090"
		s := &http.Server{Addr: addr, Handler: m}

		go func(<-chan error) {
			logger.Infof("Listending on %s", addr)
			if e := s.ListenAndServer(); e != nil {
				err <- e
			}
			return
		}(err)

		select {
		case <-hups:
			logger.Info("SIGHUP recieved. Reloading...")
			begin := time.Now()
			logger.Debug("Halting Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			s.Shutdown(ctx)
			logger.Debug("Server halted in %s ms", time.Since(begin)/time.Millisecond)
		case <-done:
			logger.Info("Server shutting down...")
			begin := time.Now()
			logger.Debug("Halting Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			s.Shutdown(ctx)
			logger.Debugf("Server halted in %s ms", time.Since(begin)/time.Millisecond)
			break
		}
	}
	waitgroup.Done()
}

func RunEchoServer(logger Logger) int {
	// First set up signal hanlding so that we can reload and stop.
	hups := make(chan os.Signal, 1)
	stop := make(chan os.Signal, 1)

	signal.Notify(hups, syscall.SIGHUP)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(hups)
	defer signal.Stop(stop)

	var waitgroup sync.WaitGroup

	done := make(chan struct{}, 1)
	err := make(chan struct{}, 1)
	defer close(done)
	defer close(err)

	waitgroup.Add(1)

	go ServeWithReload(logger, &waitgroup, err, done, hups)

	exitCode := 1
	select {
	case <-stop:
		logger.Info("Interrupt/kill recieved. Exiting...")
		exitCode = 0
		close(done)
	case <-err:
		logger.Errorf("Encoutered error %s. Exiting...", err)
		close(done)
	}
	waitgroup.Wait()
	logger.Info("All waits done. Execution complete.")
	return exitCode
}
