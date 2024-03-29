package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"x-ui/database"
	"x-ui/database/model"
	"x-ui/util/firewall"
	"x-ui/util/http"

	"github.com/tidwall/gjson"
)

type SubscribeService struct {
	inboundService InboundService
	xrayService    XrayService
}

const vmess_type = "vmess"

const subconverterUrl = "http://oracle.wocc.cf/sub?target=clash&new_name=true&url="

const scanPortApi = "https://duankou.wlphp.com/api.php"

func (s *SubscribeService) Publish() string {
	db := database.GetDB()
	var inbounds []*model.Inbound
	db.Model(model.Inbound{}).Where("user_id = 1 and enable = 1").Find(&inbounds)

	// 创建订阅
	text := ""
	for _, inbound2 := range inbounds {
		if inbound2.Protocol != vmess_type {
			continue
		}

		// 端口扫描
		b := scanPort(gjson.Get(inbound2.StreamSettings, "tlsSettings.serverName").Str, inbound2.Port)
		if !b {
			// 更新端口
			updatePort(inbound2)
		}
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
		data, err := json.MarshalIndent(vmess, "", "\t")
		if err != nil {
			fmt.Println("序列化出错,错误原因: ", err)
			return ""
		}
		if text != "" {
			text += "\n"
		}
		sEnc := vmess_type + "://" + base64.URLEncoding.EncodeToString(data)
		text = text + sEnc
	}
	return text
}

func (s *SubscribeService) Clash() string {
	text := s.Publish()
	return vmessToClash(text)
}

func updatePort(inbound *model.Inbound) (int, int) {
	oldPort := inbound.Port
	inbound.Port = inbound.Port + 1
	service := InboundService{}
	err := service.UpdateInbound(inbound)
	if err == nil {
		xrayService := XrayService{}
		xrayService.SetToNeedRestart()
	}
	return firewall.Open(inbound.Port), firewall.Close(oldPort)
}

// vmess 转clash订阅
func vmessToClash(url string) string {
	body, err := http.GetHttp(subconverterUrl + url + "%3D%3D&insert=false&config=https%3A%2F%2Fraw.githubusercontent.com%2Fwuguipeng%2Fclash-rule%2Fmaster%2Fconfig%2FACL4SSR_Online.ini")
	if err != nil {
		fmt.Println(err)
		fmt.Println("请求错误")
	}
	return string(body)
}

// 端口扫描
func scanPort(ip string, port int) bool {
	resp, err := http.GetHttp(scanPortApi + fmt.Sprintf("?i=%s&p=%d", ip, port))
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
