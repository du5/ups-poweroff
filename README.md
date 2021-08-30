# ups-poweroff
用于 macOS 网络 UPS 工具, 理论上适用于全部 `pmset -g batt` 可识别信息的 UPS。

## make

```bash
go build .
# linux/macOS
sudo mkdir -p /var/ups-poweroff
sudo mv ups-poweroff /var/ups-poweroff/
# nas
mv ups-poweroff /share/homes/
```

## edit config

```bash
# 客户端
# cp .client.yaml /var/ups-poweroff/.ups-config.yaml
# 服务端
# linux/macOS
sudo cp .service.yaml /var/ups-poweroff/.ups-config.yaml
sudo vim /var/ups-poweroff/.ups-config.yaml
# nas
sudo cp .client.yaml /share/homes/.ups-config.yaml
```
> 服务端仅支持在 macOS 设备上运行

# add service

```bash
# linux
sudo mv ups-poweroff.service /lib/systemd/system/
# 设置自启动
systemctl enable --now ups-poweroff.service
# 查看状态
systemctl status ups-poweroff.service

```

```bash
# macOS
mv ups-poweroff.plist ~/Library/LaunchAgents
# 重启后会自动启动
launchctl list | grep ups
```

```bash
# nas
mv autorun.sh /tmp/config/
# 重启后会自动启动
```