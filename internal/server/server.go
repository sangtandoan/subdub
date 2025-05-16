package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sangtandoan/subscription_tracker/internal/chrono"
)

type AppServer struct {
	*http.Server
	background *chrono.Background
}

func NewServer(addr string, handler http.Handler, background *chrono.Background) *AppServer {
	return &AppServer{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		background: background,
	}
}

func (as *AppServer) Run() {
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := as.Shutdown(ctx)
		if err != nil {
			shutdown <- err
		}

		// Call Wai() to block until our WaitGroup counter is zero --- essentially
		// blocking until all background tasks are done. Then we return nil on
		// the shutdown channel to indicate that the server has shut down gracefully.
		as.background.Wg.Wait()
		shutdown <- nil
	}()

	fmt.Println("")
	fmt.Println("server started at port 8080!")
	fmt.Println("")
	err := as.ListenAndServe()
	if err != nil {
		/// This error is expected when the server is gratefull shutdown
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server error: %v\n", err)
			return
		}
	}

	// ctx expires or error while closing listeners
	err = <-shutdown
	if err != nil {
		fmt.Printf("server forced to shutdown: %v\n", err)
		return
	}

	fmt.Println("server shutdown gratefully")
}
