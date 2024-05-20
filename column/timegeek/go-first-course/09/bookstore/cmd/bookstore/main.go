package main

import (
	_ "bookstore/internal/store"
	"bookstore/server"
	"bookstore/store/factory"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s, err := factory.New("mem")
	if err != nil {
		panic(err)
	}
	/** 这种函数原型的设计是 Go 语言的一种惯用设计方法，也就是接受一个接口类型参数，返回一个具体类型。返回的具体类型组合了传入的接口类型的能力。*/
	srv := server.NewBookStoreServer(":8080", s)

	// 测试byte字节格式
	s1 := []byte("你好，世界！")
	fmt.Println("byte:", s1) // 111111111 [228 189 160 229 165 189 239 188 140 228 184 150 231 149 140 239 188 129]  UTF-8 编码的字符串的字节序列
	str := string(s1)
	fmt.Println("string:", str) // "你好，世界！

	errChan, err := srv.ListenAndServe()
	if err != nil {
		log.Println("web server start failed:", err)
		return
	}
	log.Println("web server start ok")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-errChan:
		log.Println("web server run failed:", err)
		return
	case <-c:
		log.Println("bookstore program is exiting...")
		ctx, cf := context.WithTimeout(context.Background(), time.Second)
		defer cf()
		err = srv.Shutdown(ctx)
	}

	if err != nil {
		log.Println("bookstore program exit error:", err)
		return
	}
	log.Println("bookstore program exit ok")
}
