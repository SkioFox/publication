package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

type result struct {
	value string
}

func first(servers ...*httptest.Server) (result, error) {
	c := make(chan result, len(servers))
	queryFunc := func(server *httptest.Server) {
		url := server.URL
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("http get error: %s\n", err)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		c <- result{
			value: string(body),
		}
	}
	for _, serv := range servers {
		go queryFunc(serv)
	}
	return <-c, nil
}

func fakeWeatherServer(name string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s receive a http request\n", name)
		time.Sleep(1 * time.Second)
		w.Write([]byte(name + ":ok"))
	}))
}

/*
*
超时(timeout)与取消(cancel)模式
我们经常会使用 Go 编写向服务发起请求并获取应答结果的客户端应用。这里我们就来看一个这样的例子：我们要编写一个从气象数据服务中心获取气象信息的客户端。该客户端每次会并发向从三个气象数据服务中心发起数据查询请求，并以返回最快的那个响应信息作为此次请求的应答返回值。
*/
func main() {
	result, err := first(fakeWeatherServer("open-weather-1"),
		fakeWeatherServer("open-weather-2"),
		fakeWeatherServer("open-weather-3"))
	if err != nil {
		log.Println("invoke first error:", err)
		return
	}

	log.Println(result)
}
