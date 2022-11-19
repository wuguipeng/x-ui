package firewall

import (
	"fmt"
	"os/exec"
	"x-ui/logger"
)

func release() string {
	shell_script :=
		`
		# check os
		if [[ -f /etc/redhat-release ]]; then
			release="centos"
		elif cat /etc/issue | grep -Eqi "debian"; then
			release="debian"
		elif cat /etc/issue | grep -Eqi "ubuntu"; then
			release="ubuntu"
		elif cat /etc/issue | grep -Eqi "centos|red hat|redhat"; then
			release="centos"
		elif cat /proc/version | grep -Eqi "debian"; then
			release="debian"
		elif cat /proc/version | grep -Eqi "ubuntu"; then
			release="ubuntu"
		elif cat /proc/version | grep -Eqi "centos|red hat|redhat"; then
			release="centos"
		else
			echo -e "${red}未检测到系统版本，请联系脚本作者！${plain}\n" && exit 1
		fi
		&& echo $release
	`
	output, err := exec.Command(shell_script, "-l").Output()
	if err != nil {
		logger.Errorf("")
		return ""
	}
	return string(output)
}

func Open(port int) {
	ubuntu_open_firewall_script := fmt.Sprintf(`iptables -I INPUT -s 0.0.0.0/0 -p tcp --dport %s -j ACCEPT`, port)
	centos_open_firewall_script := fmt.Sprintf(`firewall-cmd --zone=public --add-port=%s/tcp --permanent`, port)
	release := release()
	if release == "ubuntu" {
		err := exec.Command(ubuntu_open_firewall_script).Run()
		if err != nil {
			logger.Error(fmt.Sprintf("端口%s开放开放失败"), port)
		}
	}
	if release == "centos" {
		err := exec.Command(centos_open_firewall_script).Run()
		if err != nil {
			logger.Error(fmt.Sprintf("端口%s开放开放失败"), port)
		}
	}
}

func Close(port int) {
	ubuntu_close_firewall_script := fmt.Sprintf(`iptables -I INPUT -p tcp --dport %s -j DROP`, port)
	centos_close_firewall_script := fmt.Sprintf(`firewall-cmd --zone=public --remove-port=%s/tcp --permanent`, port)
	release := release()
	if release == "ubuntu" {
		err := exec.Command(ubuntu_close_firewall_script).Run()
		if err != nil {
			logger.Error(fmt.Sprintf("端口%s开放开放失败"), port)
		}
	}
	if release == "centos" {
		err := exec.Command(centos_close_firewall_script).Run()
		if err != nil {
			logger.Error(fmt.Sprintf("端口%s开放开放失败"), port)
		}
	}
}
