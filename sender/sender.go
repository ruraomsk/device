package sender

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/device/dataBase"
	"github.com/ruraomsk/device/logsystem"
	"github.com/ruraomsk/device/memDB"
	"github.com/ruraomsk/device/setup"
	"net"
	"strconv"
	"strings"
	"time"
)

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
	var phone dataBase.Phone
	var login string
	var gdata logsystem.LoggData
	logger.Info.Printf("Новый клиент списка устройств %s", socket.RemoteAddr().String())
	reader := bufio.NewReader(socket)
	writer := bufio.NewWriter(socket)
	var message string
	defer func() {
		socket.Close()
	}()
	for {
		//logger.Debug.Println("ready ReadString")
		message, err = reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("Чтение от %s ошибка %s", socket.RemoteAddr().String(), err.Error())
			return
		}
		message = strings.ReplaceAll(message, "\n", "")
		message = strings.ReplaceAll(message, "\r", "")

		//fmt.Println(message)

		message = dataBase.DecodeString(message)
		ms := strings.Split(message, ":")
		if len(ms) < 3 {
			logger.Error.Printf("Пришло от %s строка %s", socket.RemoteAddr().String(), message)
			_, _ = writer.WriteString(dataBase.CodeString("BAD") + "\n")
			_ = writer.Flush()
			continue
		}
		is, areas := memDB.IsPhoneCorrect(ms[0], ms[1])
		if !is {
			logger.Info.Printf("Пользователь %s не зарегистрирован или неверный пароль", ms[0])
			_, _ = writer.WriteString(dataBase.CodeString("BAD") + "\n")
			_ = writer.Flush()
			continue
		}
		login = ms[0]
		gdata.Login = ms[0]
		gdata.External = false
		if strings.Compare(ms[2], "logout") == 0 {
			logger.Info.Printf("Пользователь %s завершает сеанс", login)
			return
		}

		if strings.Compare(ms[2], "getList") == 0 {
			crosses := memDB.ListCrosses(areas)
			for _, cr := range crosses {
				buf, _ := json.Marshal(&cr)
				_, _ = writer.WriteString(dataBase.CodeString(string(buf)))
				_, _ = writer.WriteString("\n")
				_ = writer.Flush()
			}
			_, _ = writer.WriteString(dataBase.CodeString("end") + "\n")
			_ = writer.Flush()
			//logger.Info.Println("getList send done...")

			memDB.PhonesTable.Lock()
			phone, err = memDB.GetPhone(login)
			if err != nil {
				logger.Error.Printf("%s %s", login, err.Error())
				return
			}
			//logger.Info.Println("getList reading phone...")

			phone.Status.LastOps = "Загрузка БД"
			phone.Status.TimeUpdateDB = time.Now()
			memDB.SetPhone(phone)
			memDB.PhonesTable.Unlock()
			//logger.Info.Println("getList done...")
			continue
		}
		gdata.Key = ms[3] + ":" + ms[4]
		gdata.External = false

		if _, err = memDB.GetCross(gdata.Key); err != nil {
			logger.Info.Printf("Пользователь %s нет такого перекрестка %s", phone.Login, gdata.Key)
			_, _ = writer.WriteString(dataBase.CodeString("BAD") + "\n")
			_ = writer.Flush()
			continue
		}
		if len(ms) == 5 && strings.Compare(ms[2], "setConnect") == 0 {
			response := "BAD"
			memDB.PhonesTable.Lock()
			phone, _ = memDB.GetPhone(login)
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
				//logger.Debug.Println("setConnect found...")
				phone.Status.LastTime = time.Now()
				phone.Status.LastOps = "Подключение к " + gdata.Key
				phone.Status.Device = gdata.Key
				phone.Status.Connection = true
				phone.Status.CurrentFaze = 0
				phone.Status.NeedFaze = 0
				memDB.SetPhone(phone)
				//logger.Debug.Println("setConnect write phone...")
				response = "Ok"
				gdata.Txt = fmt.Sprintf("Подключился к устройству")
				logsystem.LoggChan <- gdata
			}
			memDB.PhonesTable.Unlock()
			_, _ = writer.WriteString(dataBase.CodeString(response))
			_, _ = writer.WriteString("\n")
			_ = writer.Flush()
			//logger.Debug.Println("setConnect done...")
			continue
		}
		if len(ms) == 5 && strings.Compare(ms[2], "outConnect") == 0 {
			memDB.PhonesTable.Lock()
			phone, err = memDB.GetPhone(login)
			if err != nil {
				logger.Error.Printf("%s %s", login, err.Error())
				memDB.PhonesTable.Unlock()
				_, _ = writer.WriteString(dataBase.CodeString("BAD") + "\n")
				_ = writer.Flush()
				//logger.Debug.Println("outConnect abdone...")
				continue
			}
			phone.Status.LastTime = time.Now()
			phone.Status.LastOps = "Отмена РУ на " + gdata.Key
			phone.Status.Connection = false
			phone.Status.CurrentFaze = 0
			phone.Status.NeedFaze = 0
			memDB.SetPhone(phone)
			memDB.PhonesTable.Unlock()
			gdata.Txt = fmt.Sprintf("Отмена РУ")
			gdata.External = true
			logsystem.LoggChan <- gdata
			_, _ = writer.WriteString(dataBase.CodeString("Ok") + "\n")
			_ = writer.Flush()
			//logger.Debug.Println("outConnect done...")
			continue
		}
		if len(ms) == 7 && strings.Compare(ms[2], "setRU") == 0 {
			setPhase, _ := strconv.Atoi(ms[5])
			nowPhase, _ := strconv.Atoi(ms[6])
			memDB.PhonesTable.Lock()
			phone, _ = memDB.GetPhone(login)
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
			memDB.SetPhone(phone)
			memDB.PhonesTable.Unlock()
			logsystem.LoggChan <- gdata
			_, _ = writer.WriteString(dataBase.CodeString("Ok") + "\n")
			_ = writer.Flush()
			continue
		}

	}

}
