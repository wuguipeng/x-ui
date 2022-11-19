package http

import (
	"io/ioutil"
	"net/http"
	"time"
)

func GetHttp(url string) (body []byte, err error) {

	// 创建 client 和 resp 对象
	var client http.Client
	var resp *http.Response

	// 设置10秒钟的超时
	client = http.Client{Timeout: 10 * time.Second}

	// 这里使用了 Get 方法，并判断异常
	resp, err = client.Get(url)
	if err != nil {
		return nil, err
	}
	// 释放对象
	defer resp.Body.Close()

	// 把获取到的页面作为返回值返回
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 释放对象
	defer client.CloseIdleConnections()

	return body, nil
}
