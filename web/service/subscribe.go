package service

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "x-ui/database"
    "x-ui/database/model"
    "x-ui/logger"

    "github.com/tidwall/gjson"
)

type SubscribeService struct {
}

func  (s *SubscribeService) Publish() (string){
    userId := "1"
    enable := "1"
    db := database.GetDB()
    var inbounds []*model.Inbound
    db.Model(model.Inbound{}).Where("user_id = ? and enable = ?", userId, enable).Find(&inbounds)
    var text string
    for _, inbound := range inbounds {
        vmess := model.Vmess{
            V:    "2",
            Ps:   inbound.Remark,
            //            Add:  gjson.Get(inbound.StreamSettings, "tlsSettings.serverName").Str,
            Add: "oracle.xyxdbp.xyz",
            Port: inbound.Port,
            Id:   gjson.Get(inbound.Settings, "clients.0.id").Str,
            Aid:  int(gjson.Get(inbound.Settings, "clients.0.alterId").Int()),
            Net:  gjson.Get(inbound.StreamSettings, "network").Str,
            Type: gjson.Get(inbound.StreamSettings, "tcpSettings.header.type").Str,
            Host: "",
            Path: "",
            Tls:  gjson.Get(inbound.StreamSettings, "security").Str,
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
    logger.Info(text)
    return text
}
