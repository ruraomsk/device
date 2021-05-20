#!/bin/bash
echo "Start to Windows deploy"
GOOS=windows GOARCH=amd64 go build
FILE=/home/rura/mnt/ASDU/device/device.exe
if [ -f "$FILE" ]; then
    echo "Mounted the server drive"
else
    echo "Mounting the server drive"
    sudo mount -t cifs -o username=asdu,password=162747 \\\\192.168.115.115\\d /home/rura/mnt/ASDU
fi
sudo cp device.exe /home/rura/mnt/ASDU/device
sudo cp config.toml /home/rura/mnt/ASDU/device