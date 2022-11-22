package firewall

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"x-ui/logger"
)

func release() string {
	out, err := exec.Command("sh", "-c", "./release.sh").Output()
	if err != nil {
		logger.Errorf("未检测到系统版本")
		return ""
	}
	return strings.ReplaceAll(string(out), "\n", "")
}

// 执行shell命令
func ec(cmd string) (*bytes.Buffer, error) {
	command := exec.Command("sh")
	in := bytes.NewBuffer(nil)
	out := bytes.NewBuffer(nil)
	command.Stdin = in
	command.Stdout = out
	command.Stderr = os.Stderr
	in.WriteString(cmd)
	//    in.WriteString("exit\n")
	return out, command.Run()
}

func Open(port int) int {
	ubuntu_open_firewall_script := fmt.Sprintf(`iptables -I INPUT -s 0.0.0.0/0 -p tcp --dport %d -j ACCEPT`, port)
	centos_open_firewall_script := fmt.Sprintf(`firewall-cmd --zone=public --add-port=%d/tcp --permanent`, port)
	release := release()
	if release == "ubuntu" {
		_, err := ec(ubuntu_open_firewall_script)
		if err != nil {
			logger.Error(fmt.Sprintf("端口%d开放失败", port))
			return 1
		}
	}
	if release == "centos" {
		_, err := ec(centos_open_firewall_script)
		if err != nil {
			logger.Error(fmt.Sprintf("端口%d开放失败", port))
			return 1
		}
	}
	return port
}

func Close(port int) int {
	ubuntu_close_firewall_script := fmt.Sprintf(`iptables -I INPUT -p tcp --dport %d -j DROP`, port)
	centos_close_firewall_script := fmt.Sprintf(`firewall-cmd --zone=public --remove-port=%d/tcp --permanent`, port)
	release := release()
	if release == "ubuntu" {
		_, err := ec(ubuntu_close_firewall_script)
		if err != nil {
			logger.Error(fmt.Sprintf("端口%d关闭失败", port))
			return 0

		}
	}
	if release == "centos" {
		_, err := ec(centos_close_firewall_script)
		if err != nil {
			logger.Error(fmt.Sprintf("端口%d关闭失败", port))
			return 0
		}
	}
	return port
}
