package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/JanFant/TLServer/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"net"
	"rura/device/devs"
	"rura/device/terminal"
	"strings"
	"time"
)

type dev struct {
	host     string
	port     string
	name     string
	login    string
	password string
}

func (d *dev) toString() string {
	d.name = strings.ReplaceAll(d.name, ":", "-")
	return d.host + ":" + d.port + ":" + d.name + ":" + d.login + ":" + d.password
}
func sender(soc net.Conn) {

	defer soc.Close()
	var err error
	logger.Info.Printf("Новый клиент списка устройств %s", soc.RemoteAddr().String())
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		"192.168.115.115", "postgres", "162747", "agserv")
	dbb, err := sql.Open("postgres", dbinfo)
	if err != nil {
		logger.Error.Printf("Запрос на открытие %s %s", dbinfo, err.Error())
		return
	}
	defer dbb.Close()
	if err = dbb.Ping(); err != nil {
		logger.Error.Printf("Ping %s", err.Error())
		return
	}
	list := make([]dev, 0)
	d := dev{
		host:     "192.168.115.159",
		port:     "8888",
		name:     "Отладочный",
		login:    "login",
		password: "password",
	}
	list = append(list, d)
	rows, err := dbb.Query("select id,device from public.devices;")
	if err != nil {
		logger.Error.Printf("При чтении списка устройств %s", err.Error())
		return
	}
	i := 10
	for rows.Next() {
		var c pudge.Controller
		var id int
		var jc []byte
		_ = rows.Scan(&id, &jc)
		err = json.Unmarshal(jc, &c)
		ips := strings.Split(c.IPHost, ":")
		if len(ips) == 2 {
			d.host = ips[0]
			d.port = ips[1]
			d.name = c.Name
		} else {
			d.host = fmt.Sprintf("192.168.1.%d", i)
			d.port = "8888"
			d.name = "Имитатор " + c.Name
			i++
		}
		list = append(list, d)
	}
	writer := bufio.NewWriter(soc)
	for _, d := range list {
		logger.Info.Printf(d.toString())
		_, _ = writer.WriteString(d.toString() + "\n")
		err = writer.Flush()
		if err != nil {
			logger.Error.Printf("При передаче списка устройств %s", err.Error())
			return
		}
	}
	writer.WriteString("end\n")
	err = writer.Flush()
	if err != nil {
		logger.Error.Printf("При передаче списка устройств %s", err.Error())
		return
	}
	time.Sleep(10 * time.Second)
}
func sendDevices() {
	ln, err := net.Listen("tcp", ":8088")
	if err != nil {
		logger.Error.Printf("Ошибка открытия порта %s", err.Error())
		return
	}
	defer ln.Close()
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Ошибка accept %s", err.Error())
			continue
		}
		go sender(socket)
	}

}

//Простой сервер для отладки adBox
//Слушаем порт если кто приперся то логин пароль и затем ждем вопросы и шлем ответы
func main() {
	_ = logger.Init(".")
	logger.Info.Println("Device start")
	fmt.Println("Device start")
	go devs.Listen()
	go sendDevices()
	terminal.Terminal()
}
