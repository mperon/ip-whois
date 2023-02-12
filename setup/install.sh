#!/usr/bin/env bash

sudo mkdir -p /opt/golang/{bin,share,include}
sudo chown -R devel /opt/golang
sudo chmod g+w /opt/golang/bin
rsync -avh ./bin/ip-whois-linux-amd64 devel@mperon:/opt/golang/bin/
sudo chmod +x /opt/golang/bin/ip-whois-linux-amd64


#/etc/systemd/system/ip-whois.service

systemctl start ip-whois
systemctl status ip-whois
