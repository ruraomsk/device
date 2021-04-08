package devs

import (
	"bufio"
	"bytes"
	"github.com/JanFant/aura"
	"github.com/ruraomsk/TLServer/logger"
	"net"
	"strconv"
	"strings"
	"time"
)

var clients map[string]Client
var perfect Controller

type Client struct {
	Data    Controller
	login   bool
	echo    echo
	now     time.Time
	reading chan string
	writing chan string
	socket  net.Conn
	stop    chan interface{}
	deleted bool
}

func (c Client) SendInfo() {
	list := aura.Marshal(c.Data.Info)

	for _, l := range list {
		if isNeedInsertTime(l) {
			l = insertTime(l)
		}
		c.writing <- l
	}
}

type echo struct {
	modem  bool
	gps    bool
	pspd   bool
	server bool
	pult   bool
	can    bool
	modbus bool
	def    bool
}

func (e *echo) set(code string) {
	c, _ := strconv.Atoi(code)
	if c&1 != 0 {
		e.def = true
	} else {
		e.def = false
	}
	if c&2 != 0 {
		e.modem = true
	} else {
		e.modem = false
	}
	if c&4 != 0 {
		e.gps = true
	} else {
		e.gps = false
	}
	if c&8 != 0 {
		e.pspd = true
	} else {
		e.pspd = false
	}
	if c&16 != 0 {
		e.server = true
	} else {
		e.server = false
	}
	if c&32 != 0 {
		e.pult = true
	} else {
		e.pult = false
	}
	if c&64 != 0 {
		e.can = true
	} else {
		e.can = false
	}
	if c&128 != 0 {
		e.modbus = true
	} else {
		e.modbus = false
	}
}
func (e *echo) sayProtocol() []string {
	var result []string
	if e.modem {
		result = append(result, "Modem message")
	}
	if e.gps {
		result = append(result, "GPS message")
	}
	if e.server {
		result = append(result, "Server message")
	}
	if e.pult {
		result = append(result, "Pult message")
	}
	if e.can {
		result = append(result, "CAN message")
	}
	if e.modbus {
		result = append(result, "Modbus message")
	}
	return result
}

func (c *Client) init(soc net.Conn) {
	logger.Info.Printf("Новый клиент сервера устройства %s", soc.RemoteAddr().String())
	c.login = false
	c.socket = soc
	c.stop = make(chan interface{})
	c.reading = make(chan string, 1000)
	c.writing = make(chan string, 1000)
	c.now = time.Now()
	c.deleted = false
	c.Data = perfect
	go c.DeviceReader()
	go c.DeviceWriter()
	go c.DeviceWorker()
}
func (c *Client) StopDevice() {
	if c.deleted {
		return
	}
	c.deleted = true
	logger.Info.Printf("Остановлена работа с %s ", c.getId())
	time.Sleep(time.Second)
	c.stop <- "stop"
	c.stop <- "stop"
	c.socket.Close()
	delete(clients, c.getId())
}

func (c *Client) DeviceReader() {
	defer c.StopDevice()
	reader := bufio.NewReader(c.socket)
	for {
		bytel, _, err := reader.ReadLine()
		if c.deleted {
			return
		}
		if err != nil {
			logger.Error.Printf("Ошибка при чтении %s %s", c.getId(), err.Error())
			return
		}
		if bytes.HasSuffix(bytel, []byte{0x0D}) {
			bytel = bytel[0 : len(bytel)-1]
		}
		c.reading <- string(bytel)
	}
}
func (c *Client) getId() string {
	return c.socket.RemoteAddr().String()
}

func (c *Client) DeviceWriter() {
	writer := bufio.NewWriter(c.socket)
	for {
		select {
		case <-c.stop:
			return
		case str := <-c.writing:
			writer.WriteString(str)
			writer.WriteString("\r\n")
			err := writer.Flush()
			if err != nil {
				logger.Error.Printf("Ошибка при записи в %s %s", c.getId(), err.Error())
				c.StopDevice()
				return
			}
		}
	}
}

func (c *Client) DeviceWorker() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			if c.echo.def {
				c.SendInfo()
			}
			list := c.echo.sayProtocol()
			for _, l := range list {
				c.writing <- l
			}
		case command := <-c.reading:
			c.makeCommand(command)
		}
	}
}

func (c *Client) makeCommand(com string) {
	if strings.HasPrefix(com, "login") {
		c.login = true
		c.writing <- "OK"
		return
	}
	if !c.login {
		c.writing <- "need login"
		return
	}
	if strings.HasPrefix(com, "exit") {
		c.StopDevice()
		return
	}
	if isEmpty(com) {
		return
	}
	if isSave(com) {
		saveDevice(c.Data)
		c.writing <- "OK"
		return
	}
	if isSetTime(com) {
		logger.Info.Println("Установить время команда")
		return
	}
	com = strings.ReplaceAll(com, "=", ":")
	cc := new(Command)
	var ls []string
	ls = append(ls, com)
	err := aura.UnMarshal(ls, cc)
	if err != nil {
		c.writing <- "ERROR - " + com
		return
	}
	//fmt.Printf("%v\n", cc)
	if len(cc.Mode) != 0 {
		if strings.Compare(cc.Mode, "000") == 0 {
			c.echo.set(cc.Mode)
			c.SendInfo()
			return
		} else {
			c.echo.set(cc.Mode)
			return
		}
	}
	if len(cc.Gate) != 0 {
		c.Data.Info.Gate = strings.ReplaceAll(cc.Gate, "\"", "")
		return
	}
	if len(cc.Ipcontroller) != 0 {
		//
		ls := strings.Split(cc.Ipcontroller, ",")
		c.Data.Info.Ipcontroller = strings.ReplaceAll(ls[0], "\"", "")
		return
	}
	if len(cc.IpserverGPRS) != 0 {
		ls := strings.Split(cc.IpserverGPRS, ",")
		c.Data.Info.IpserverGPRS = strings.ReplaceAll(ls[0], "\"", "")
		return
	}
	if len(cc.IpserverLAN) != 0 {
		c.Data.Info.IpserverLAN = strings.ReplaceAll(cc.IpserverLAN, "\"", "")
		return
	}
	if len(cc.Mask) != 0 {
		c.Data.Info.Mask = strings.ReplaceAll(cc.Mask, "\"", "")
		return
	}
	c.writing <- "not define " + com
}
func Listen() {
	ln, err := net.Listen("tcp", ":8888")
	clients = make(map[string]Client)
	loadDevice()
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
		c := new(Client)
		c.init(socket)
		clients[c.getId()] = *c
	}
}
