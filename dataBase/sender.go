package dataBase

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/device/setup"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mutexSender sync.Mutex

func ListenPhone() {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(setup.Set.Port))
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
		go workerPhone(socket)
	}
}
func workerPhone(socket net.Conn) {
	var err error
	var phone Phone
	var gdata LoggData
	logger.Info.Printf("Новый клиент списка устройств %s", socket.RemoteAddr().String())
	reader := bufio.NewReader(socket)
	writer := bufio.NewWriter(socket)
	var message string
	defer func() {
		socket.Close()
	}()
	for {
		message, err = reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("Чтение от %s ошибка %s", socket.RemoteAddr().String(), err.Error())
			return
		}
		message = strings.ReplaceAll(message, "\n", "")
		message = strings.ReplaceAll(message, "\r", "")
		//fmt.Println(message)
		ms := strings.Split(message, ":")
		if len(ms) < 3 {
			logger.Error.Printf("Пришло от %s строка %s", socket.RemoteAddr().String(), message)
			_, _ = writer.WriteString("BAD\n")
			writer.Flush()
			continue
		}
		is, areas := IsPhoneCorrect(ms[0], ms[1])
		if !is {
			logger.Info.Printf("Пользователь %s не зарегистрирован или неверный пароль", ms[0])
			_, _ = writer.WriteString("BAD\n")
			writer.Flush()
			continue
		}
		phone.Login = ms[0]
		gdata.Login = ms[0]
		gdata.External = false
		if strings.Compare(ms[2], "getList") == 0 {
			crosses := GetCrosses(areas)
			for _, cr := range crosses {
				buf, _ := json.Marshal(&cr)
				_, _ = writer.WriteString(string(buf))
				_, _ = writer.WriteString("\n")
				writer.Flush()
			}
			_, _ = writer.WriteString("end\n")
			writer.Flush()

			mutexSender.Lock()
			phone.GetPhone()
			phone.Status.LastOps = "Загрузка БД"
			phone.Status.TimeUpdateDB = time.Now()
			phone.SetPhone()
			mutexSender.Unlock()
			continue
		}
		gdata.Key = ms[3] + ":" + ms[4]
		gdata.External = false
		if !IsCross(gdata.Key) {
			logger.Info.Printf("Пользователь %s нет такого перекрестка %s", phone.Login, gdata.Key)
			_, _ = writer.WriteString("BAD\n")
			writer.Flush()
			continue
		}
		if len(ms) == 5 && strings.Compare(ms[2], "setConnect") == 0 {
			response := "BAD"
			mutexSender.Lock()
			phone.GetPhone()
			area, _ := strconv.Atoi(ms[3])
			found := false
			if len(phone.Areas) == 0 {
				found = true
			} else {
				for _, ar := range phone.Areas {
					if ar == area {
						found = true
						break
					}
				}
			}
			if found {
				phone.Status.LastTime = time.Now()
				phone.Status.LastOps = "Подключение к " + gdata.Key
				phone.Status.Device = gdata.Key
				phone.Status.Connection = true
				phone.Status.CurrentFaze = 0
				phone.Status.NeedFaze = 0
				phone.SetPhone()
				mutexSender.Unlock()
				response = "Ok"
				gdata.Txt = fmt.Sprintf("Подключился к устройству")
				LoggChan <- gdata
			} else {
				mutex.Unlock()
			}
			_, _ = writer.WriteString(response)
			_, _ = writer.WriteString("\n")
			writer.Flush()
			continue
		}
		if len(ms) == 5 && strings.Compare(ms[2], "outConnect") == 0 {
			mutexSender.Lock()
			phone.GetPhone()
			phone.Status.LastTime = time.Now()
			phone.Status.LastOps = "Отмена РУ на " + gdata.Key
			phone.Status.Connection = false
			phone.Status.CurrentFaze = 0
			phone.Status.NeedFaze = 0
			phone.SetPhone()
			mutexSender.Unlock()
			gdata.Txt = fmt.Sprintf("Отмена РУ")
			gdata.External = true
			LoggChan <- gdata
			_, _ = writer.WriteString("Ok")
			_, _ = writer.WriteString("\n")
			writer.Flush()
			continue
		}
		if len(ms) == 7 && strings.Compare(ms[2], "setRU") == 0 {
			setPhase, _ := strconv.Atoi(ms[5])
			nowPhase, _ := strconv.Atoi(ms[6])
			mutexSender.Lock()
			phone.GetPhone()
			gdata.External = true
			phone.Status.LastTime = time.Now()
			phone.Status.CurrentFaze = nowPhase
			phone.Status.NeedFaze = setPhase
			if setPhase != nowPhase {
				phone.Status.LastOps = fmt.Sprintf("Устройство %s Переход в РУ с фазой %d текущая %d", gdata.Key, setPhase, nowPhase)
				gdata.Txt = fmt.Sprintf("Переход в РУ с фазой %d текущая %d", setPhase, nowPhase)
			} else {
				phone.Status.LastOps = fmt.Sprintf("Устройство %s установлено РУ фаза %d", gdata.Key, setPhase)
				gdata.Txt = fmt.Sprintf("Установлен РУ с фазой %d ", setPhase)
			}
			phone.SetPhone()
			mutexSender.Unlock()
			LoggChan <- gdata
			_, _ = writer.WriteString("Ok")
			_, _ = writer.WriteString("\n")
			writer.Flush()
			continue
		}

	}

}
