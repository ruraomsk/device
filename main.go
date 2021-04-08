package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/device/devs"
	"github.com/ruraomsk/device/terminal"
	"net"
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
		host:     "192.168.115.85",
		port:     "1100",
		name:     "Отладочный",
		login:    "login",
		password: "password",
	}
	//list = append(list, d)
	rows, err := dbb.Query("select state from public.cross;")
	if err != nil {
		logger.Error.Printf("При чтении списка устройств %s", err.Error())
		return
	}
	for rows.Next() {
		var c pudge.Cross
		var jc []byte
		_ = rows.Scan(&jc)
		err = json.Unmarshal(jc, &c)
		if len(c.WiFi) == 0 {
			continue
		}
		d.host = c.WiFi
		d.port = "1100"
		d.name = c.Name
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
