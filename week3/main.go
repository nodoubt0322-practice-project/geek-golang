/*
問題描述：
基於 errgroup 實現一個 http server 的啓動和關閉 ，以及 linux signal 信號的註冊和處理，要保證能夠一個退出，全部註銷退出。
根據描述信息，可以簡單匯總成3塊內容：
1.實現HTTP server的啓動和關閉
2.監聽linux signal信號，使用chan實現對linux signal中斷的註冊和處理 按ctrl+c之類退出程序
3.errgroup實現多個goroutine的級聯退出，通過errgroup+context的形式，對1、2中的goroutine進行級聯註銷
*/

package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)



func main() {

	g, ctx := errgroup.WithContext(context.Background())

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("world"))
	})


	serverOut := make(chan struct{})
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		serverOut <- struct{}{} 
	})

	server := http.Server{
		Addr: ":8080",
	}


	g.Go(func() error {
		return server.ListenAndServe()
	})


	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("errgroup exit...")
		case <-serverOut:
			log.Println("server will out...")
		}

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		log.Println("shutting down server...")
		return server.Shutdown(timeoutCtx) 
	})

	g.Go(func() error {
		c := make(chan os.Signal, 1) 
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-c:
			return errors.Errorf("get os signal: %v", sig)
		}
	})

	fmt.Printf("errgroup exiting: %+v\n", g.Wait())
}
