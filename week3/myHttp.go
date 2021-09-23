package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出
func main() {
	g, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	srv := &http.Server{Addr: ":8080"}

	g.Go(func() error {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
			fmt.Fprintln(w, "hello")
		})
		return srv.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()
		fmt.Println("shutdown http...")
		return srv.Shutdown(ctx)
	})

	g.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case <- ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			cancel()
			return fmt.Errorf("get signal:%v", sig)
		}
	})

	err := g.Wait()
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	}
}


