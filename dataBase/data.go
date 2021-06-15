package dataBase

import (
	"github.com/ruraomsk/device/setup"
	"strings"
	"time"
)

type Cross struct {
	Area     int    `json:"area"`
	ID       int    `json:"id"`
	IDevice  int    `json:"idevice"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	SSID     string `json:"ssid"`     //SSID wifi с контроллером
	PassID   string `json:"passid"`   //password для соединения с локальной сетью
	Login    string `json:"login"`    //login для соединения с контроллером
	Password string `json:"password"` //password для соединения с контроллером
	Fazes    []Faza `json:"fazes"`    //доступные фазы на контроллере
}
type Faza struct {
	Number int    `json:"number"` //Номер фазы
	Name   string `json:"name"`   //Краткое наименование
	Png    []byte `json:"png"`    //Рисунок фазы
}
type Phone struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Areas    []int  `json:"areas"`
	Status   Status `json:"status"`
}
type Status struct {
	TimeUpdateDB time.Time `json:"dateDB"`   //Время последнего обновления БД на телефоне
	LastTime     time.Time `json:"ltime"`    //Время последней операции
	LastOps      string    `json:"last_ops"` //Последняя операция
	Device       string    `json:"device"`   //К какому устройству подключен пусто если нет
	Connection   bool      `json:"connect"`  //Истина если произведено текущее подключение
	NeedFaze     int       `json:"nfaze"`    //Заказанная фаза на устройстве
	CurrentFaze  int       `json:"cfaze"`    //Текущая фаза на устройстве
}

func (c *Cross) Default(ipHost string) {
	if ipHost != "" {
		ips := strings.Split(ipHost, ":")
		if len(ips) != 0 {
			c.Host = ips[0]
		} else {
			c.Host = setup.Set.Default.Host
		}
	} else {
		c.Host = setup.Set.Default.Host
	}
	c.Port = setup.Set.Default.Port
	c.SSID = setup.Set.Default.SSID
	c.PassID = setup.Set.Default.PassID
	c.Login = setup.Set.Default.Login
	c.Password = setup.Set.Default.Password
	c.Fazes = make([]Faza, 0)
}
