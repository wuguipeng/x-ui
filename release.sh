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
echo $release