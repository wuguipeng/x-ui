package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"time"
	"x-ui/database"
	"x-ui/database/model"
)

type SubscribeService struct {
	inboundService InboundService
}

func (s *SubscribeService) Publish() string {
	db := database.GetDB()
	var inbounds []*model.Inbound
	db.Model(model.Inbound{}).Where("user_id = 1 and enable = 1").Find(&inbounds)

	// 扫描端口
	for _, inbound1 := range inbounds {
		b := scanPort(gjson.Get(inbound1.StreamSettings, "tlsSettings.serverName").Str, inbound1.Port)
		if !b {
			updatePort(inbound1)
		}
	}
	// 创建订阅
	var text string
	for _, inbound2 := range inbounds {
		vmess := model.Vmess{
			V:    "2",
			Ps:   inbound2.Remark,
			Add:  gjson.Get(inbound2.StreamSettings, "tlsSettings.serverName").Str,
			Port: inbound2.Port,
			Id:   gjson.Get(inbound2.Settings, "clients.0.id").Str,
			Aid:  int(gjson.Get(inbound2.Settings, "clients.0.alterId").Int()),
			Net:  gjson.Get(inbound2.StreamSettings, "network").Str,
			Type: gjson.Get(inbound2.StreamSettings, "tcpSettings.header.type").Str,
			Host: "",
			Path: "",
			Tls:  gjson.Get(inbound2.StreamSettings, "security").Str,
		}
		data, err := json.Marshal(&vmess)
		if err != nil {
			fmt.Println("序列化出错,错误原因: ", err)
			return ""
		}

		sEnc := "vmess://" + base64.StdEncoding.EncodeToString([]byte(string(data)))
		text = text + sEnc + "\n"
	}
	text = base64.StdEncoding.EncodeToString([]byte(text))
	return text
}

func updatePort(inbound *model.Inbound) {
	inbound.Port = inbound.Port + 1
	service := InboundService{}
	err := service.UpdateInbound(inbound)
	if err != nil {
		return
	}
}

//func scanPort(network string, ip string, port int) bool {
//	conn, _ := net.DialTimeout(network, fmt.Sprintf("%s:%d", ip, port), time.Millisecond*time.Duration(500))
//	if conn != nil {
//		err := conn.Close()
//		if err != nil {
//			panic(err)
//		}
//		return true
//	}
//	return false
//}
func scanPort(ip string, port int) bool {
	resp, err := GetHttp(fmt.Sprintf("https://duankou.wlphp.com/api.php?i=%s&p=%d", ip, port))
	if err != nil {
		fmt.Println(err)
	}
	status := gjson.Get(string(resp), "msg.status")
	if status.Str == "Openning" {
		return true
	} else {
		return false
	}
}

func GetHttp(url string) (body []byte, err error) {

	// 创建 client 和 resp 对象
	var client http.Client
	var resp *http.Response

	// 这里博主设置了10秒钟的超时
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
