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
)

type AppServer struct {
	*http.Server
}

func NewServer(addr string, handler http.Handler) *AppServer {
	return &AppServer{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
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

		shutdown <- as.Shutdown(ctx)
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
