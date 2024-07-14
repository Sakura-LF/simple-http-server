package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"simple-http-server/internal/store2"
	"simple-http-server/server"
	"simple-http-server/store"
	"simple-http-server/store/factory"
	"syscall"
	"time"
)

func main() {
	// 手动创建factory
	storeFactory := factory.NewStoreFactory()
	// 注册store
	storeFactory.Register("mem", &store2.MemBookStore{Books: make(map[string]*store.Book)})

	// 获取store
	// 创建存储对象
	s, err := storeFactory.New("mem")
	if err != nil {
		panic(err)
	}

	// 创建BookStoreServer
	srv := server.NewBookStoreServer(":8080", s)

	// 启动服务器
	errChan, err := srv.ListenAndServe()
	if err != nil {
		log.Println("Web Server Start Failed", err)
		return
	}
	log.Print("Web Server Start Success")

	// 监听系统信号
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	// 发送错误的信号
	case err = <-errChan:
		log.Println("bookstore server error:", err)
		return
	// 正常推出的信号
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
