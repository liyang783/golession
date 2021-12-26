package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, groupCtx := errgroup.WithContext(ctx)
	srv := &http.Server{
		Addr: ":8080",
	}

	g.Go(func() error {

		http.HandleFunc("/", indexHandle)
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		fmt.Println("starting in 8080 port")
		<-groupCtx.Done()
		return srv.Shutdown(groupCtx)
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	g.Go(func() error {
		select {
		case <-groupCtx.Done():
			return groupCtx.Err()
		case <-c:
			cancel()
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Println("errgroup wait err:", err)
	}
}

func indexHandle(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Hello world")
}
