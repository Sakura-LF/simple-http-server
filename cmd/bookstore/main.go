package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	_ "simple-http-server/internal/store"
	"simple-http-server/server"
	"simple-http-server/store/factory"
	"syscall"
	"time"
)

func main() {
	// 创建存储对象
	s, err := factory.New("mem")

	if err != nil {
		panic(err)
	}

	// 创建BookStoreServer
	srv := server.NewBookStoreServer(":8080", s)

	errchan, err := srv.ListenAndServe()
	if err != nil {
		log.Println("Web Server Start Failed", err)
		return
	}
	log.Print("Web Server Start Success")

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-errchan:
		log.Println("bookstore server error:", err)
		return
	case <-c:
		log.Println("bookstore program is exiting...")
		ctx, cf := context.WithTimeout(context.Background(), time.Second)
		defer cf()
		err = srv.Shutdown(ctx) // 优雅关闭http服务实例
	}

	if err != nil {
		log.Println("bookstore program exit error:", err)
		return
	}
	log.Println("bookstore program exit ok")
}
