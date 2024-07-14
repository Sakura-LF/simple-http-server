package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"simple-http-server/server/middleware"
	"simple-http-server/store"
	"time"
)

// BookStoreServer 定义一个专门用于处理 BookStore 的 Server
type BookStoreServer struct {
	s   store.Store
	srv *http.Server
}

func NewBookStoreServer(addr string, s store.Store) *BookStoreServer {
	srv := &BookStoreServer{
		s: s,
		srv: &http.Server{
			Addr: addr,
		},
	}
	mux := http.NewServeMux()

	// RESTFul 风格的API 标准接口
	// 创建书接口
	mux.HandleFunc("POST /book", srv.createBookHandler)
	// 更新书接口
	mux.HandleFunc("PATCH /book/{id}", srv.updateBookHandler)
	// 获取单本书接口
	mux.HandleFunc("GET /book/{id}", srv.getBookHandler)
	// 获取全部书接口
	mux.HandleFunc("GET /book", srv.getAllBookHandler)
	// 删除书接口
	mux.HandleFunc("DELETE /book/{id}", srv.deleteBookHandler)

	// 添加中间件，按顺序添加
	// 1. 打印请求信息中间件
	// 2. 参数校验中间件
	srv.srv.Handler = middleware.Logging(middleware.Validating(mux))

	return srv
}

func (bs *BookStoreServer) ListenAndServe() (<-chan error, error) {
	var err error
	errChan := make(chan error)
	go func() {
		err = bs.srv.ListenAndServe()
		errChan <- err
	}()

	select {
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second):
		return errChan, nil
	}
}

func (bs *BookStoreServer) Shutdown(ctx context.Context) error {
	return bs.srv.Shutdown(ctx)
}

func (bs *BookStoreServer) createBookHandler(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello")
	dec := json.NewDecoder(req.Body)
	var book store.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	if err := bs.s.Create(&book); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bs *BookStoreServer) getBookHandler(resp http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")

	book, err := bs.s.Get(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
	response(resp, book)
}

func (bs *BookStoreServer) updateBookHandler(resp http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")

	dec := json.NewDecoder(req.Body)
	var book store.Book
	if err := dec.Decode(&book); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	book.Id = id
	if err := bs.s.Update(&book); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bs *BookStoreServer) deleteBookHandler(resp http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")

	err := bs.s.Delete(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}
}

func (bs *BookStoreServer) getAllBookHandler(resp http.ResponseWriter, req *http.Request) {
	books, err := bs.s.GetAll()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	response(resp, books)
}

func response(resp http.ResponseWriter, req interface{}) {
	data, err := json.Marshal(req)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}
