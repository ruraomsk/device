package tests

import (
	"bufio"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/device/dataBase"
	"github.com/ruraomsk/device/memDB"
	"github.com/ruraomsk/device/setup"
	"net"
	"strconv"
	"strings"
	"time"
)

type changer struct {
	writer   *bufio.Writer
	reader   *bufio.Reader
	ph       dataBase.Phone
	password string
}

func (c *changer) readDB() bool {
	fmt.Print("Тест БД")
	c.writer.WriteString(dataBase.CodeString(c.ph.Login+":"+c.password+":"+"getList") + "\n")
	_ = c.writer.Flush()
	for {
		//soc.SetReadDeadline(time.Now().Add(5*time.Second))
		s, err := c.reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("read DB %s", err.Error())
			return false
		}
		s = dataBase.DecodeString(s)
		if strings.HasPrefix(s, "end") {
			return true
		}
		if strings.HasPrefix(s, "BAD") {
			return false
		}
	}
}
func (c *changer) logout() {
	c.writer.WriteString(dataBase.CodeString(c.ph.Login+":"+c.password+":"+"logout") + "\n")
	_ = c.writer.Flush()
	time.Sleep(time.Second)
}

func (c *changer) setConnect(cr dataBase.Cross) bool {
	fmt.Print("Тест соединения")
	key := fmt.Sprintf("%d:%d", cr.Area, cr.ID)
	c.writer.WriteString(dataBase.CodeString(c.ph.Login+":"+c.password+":setConnect:"+key) + "\n")
	_ = c.writer.Flush()
	for {
		//soc.SetReadDeadline(time.Now().Add(5*time.Second))
		s, err := c.reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("setConnect %s", err.Error())
			return false
		}
		s = dataBase.DecodeString(s)
		if strings.HasPrefix(s, "Ok") {
			return true
		}
		if strings.HasPrefix(s, "BAD") {
			return false
		}
	}

}
func (c *changer) disConnect(cr dataBase.Cross) bool {
	fmt.Print("Тест разрыва")
	key := fmt.Sprintf("%d:%d", cr.Area, cr.ID)
	c.writer.WriteString(dataBase.CodeString(c.ph.Login+":"+c.password+":outConnect:"+key) + "\n")
	_ = c.writer.Flush()
	for {
		//soc.SetReadDeadline(time.Now().Add(5*time.Second))
		s, err := c.reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("outConnect %s", err.Error())
			return false
		}
		s = dataBase.DecodeString(s)
		if strings.HasPrefix(s, "Ok") {
			return true
		}
		if strings.HasPrefix(s, "BAD") {
			return false
		}
	}

}

func (c *changer) setPhase(cr dataBase.Cross, set int, now int) bool {
	fmt.Printf("Ставим на %d %d фазу %d текущая %d Тест", cr.Area, cr.ID, set, now)
	key := fmt.Sprintf("%d:%d:%d:%d", cr.Area, cr.ID, set, now)
	c.writer.WriteString(dataBase.CodeString(c.ph.Login+":"+c.password+":setRU:"+key) + "\n")
	_ = c.writer.Flush()
	for {
		//soc.SetReadDeadline(time.Now().Add(5*time.Second))
		s, err := c.reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("setRU %s", err.Error())
			return false
		}
		s = dataBase.DecodeString(s)
		if strings.HasPrefix(s, "Ok") {
			return true
		}
		if strings.HasPrefix(s, "BAD") {
			return false
		}
	}

}

var ch changer
var soc net.Conn
var err error

func PhoneTest() error {
	//l := "testTestTest"
	//c := dataBase.CodeString(l)
	//d := dataBase.DecodeString(c)
	//if strings.Compare(l, d) != 0 {
	//	fmt.Printf("<%s> != <%s>\n", l, d)
	//	panic("exit")
	//}
	soc, err = net.Dial("tcp", "localhost:"+strconv.Itoa(setup.Set.Port))
	if err != nil {
		logger.Error.Printf("tcp dial %s", err.Error())
		return err
	}
	defer soc.Close()
	var stoped = fmt.Errorf("тесты не пройдены")
	ph, err := memDB.GetPhone("newrura")
	if err != nil {
		ph = dataBase.Phone{Login: "newrura", Password: dataBase.GetHasPassword("162747"), Name: "Тестовый пользователь", Areas: make([]int, 0)}
		ph.Areas = append(ph.Areas, 1, 2, 3)
		ph.Status = dataBase.Status{Connection: false}
		memDB.SetPhone(ph)
	}
	memDB.SetPhone(ph)
	ch.ph = ph
	ch.password = "162747"

	ch.writer = bufio.NewWriter(soc)
	ch.reader = bufio.NewReader(soc)

	fmt.Println("Старт тесты")
	if !ch.readDB() {
		fmt.Println(" не прошел")
		return stoped
	} else {
		fmt.Println(" ok")
	}
	crosses := memDB.ListCrosses(ch.ph.Areas)
	for _, c := range crosses {
		if !ch.setConnect(c) {
			fmt.Println(" не прошел")
			return stoped
		} else {
			fmt.Println(" ok")
		}
		if !ch.setPhase(c, 5, 1) {
			fmt.Println(" не прошел")
			return stoped
		} else {
			fmt.Println(" ok")
		}
		if !ch.setPhase(c, 5, 5) {
			fmt.Println(" не прошел")
			return stoped
		} else {
			fmt.Println(" ok")
		}
		if !ch.disConnect(c) {
			fmt.Println(" не прошел")
			return stoped
		} else {
			fmt.Println(" ok")
		}

	}
	ch.logout()
	fmt.Println("Все тесты ok")
	return nil

}
