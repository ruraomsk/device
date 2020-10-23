package terminal

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"rura/device/devs"
	"strings"
	"time"
)

var socket net.Conn
var err error
var prompt string

func readFromDevice() {
	reader := bufio.NewReader(socket)
	for {
		bytel, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("stop terminal " + err.Error())
			return
		}
		if bytes.HasSuffix(bytel, []byte{0x0D}) {
			bytel = bytel[0 : len(bytel)-1]
		}
		fmt.Printf("\n%s\n%s", string(bytel), prompt)
	}
}
func Terminal() {
	time.Sleep(5 * time.Second)
	prompt = ">"
	for {
		socket, err = net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			panic("Нет связи! " + err.Error())
		}
		go readFromDevice()
		reader := bufio.NewReader(os.Stdin)
		writer := bufio.NewWriter(socket)
		_, _ = writer.WriteString("login=login password=password")
		_, _ = writer.WriteString("\r\n")
		err = writer.Flush()
		if err != nil {
			panic("Ошибка вывода " + err.Error())
		}
		fmt.Println("login=login password=password")
		for true {
			cmd, _, err := reader.ReadLine()
			if err != nil {
				panic("Беда с клавиатурой))")
			}
			command := string(cmd)
			if strings.HasSuffix(command, "save") {
				command = devs.SaveController[0]
			}
			if strings.HasSuffix(command, "main") {
				command = "#DBG.MODE:000"
			}
			if strings.HasSuffix(command, "all") {
				command = "#DBG.MODE:255"
			}
			_, _ = writer.WriteString(command)
			_, _ = writer.WriteString("\r\n")
			err = writer.Flush()
			if err != nil {
				panic("Ошибка вывода " + err.Error())
			}
			if strings.HasSuffix(command, "exit") {
				_ = socket.Close()
				break
			}
		}
	}
}
