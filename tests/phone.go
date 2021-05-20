package tests

import (
	"bufio"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/device/dataBase"
	"github.com/ruraomsk/device/setup"
	"net"
	"strconv"
	"strings"
)

type changer struct {
	writer *bufio.Writer
	reader *bufio.Reader
	ph     *dataBase.Phone
}

func (c *changer) readDB() bool {
	fmt.Print("Тест БД")
	c.writer.WriteString(c.ph.Login + ":" + c.ph.Password + ":" + "getList\n")
	c.writer.Flush()
	for {
		s, err := c.reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("read DB %s", err.Error())
			return false
		}
		if strings.HasPrefix(s, "end") {
			return true
		}
		if strings.HasPrefix(s, "BAD") {
			return false
		}
	}
}

func (c *changer) setConnect(cr dataBase.Cross) bool {
	fmt.Print("Тест соединения")
	key := fmt.Sprintf("%d:%d", cr.Area, cr.ID)
	c.writer.WriteString(c.ph.Login + ":" + c.ph.Password + ":setConnect:" + key + "\n")
	c.writer.Flush()
	for {
		s, err := c.reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("setConnect %s", err.Error())
			return false
		}
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
	c.writer.WriteString(c.ph.Login + ":" + c.ph.Password + ":outConnect:" + key + "\n")
	c.writer.Flush()
	for {
		s, err := c.reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("outConnect %s", err.Error())
			return false
		}
		logger.Info.Println(s)
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
	c.writer.WriteString(c.ph.Login + ":" + c.ph.Password + ":setRU:" + key + "\n")
	c.writer.Flush()
	for {
		s, err := c.reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("setRU %s", err.Error())
			return false
		}
		logger.Info.Println(s)
		if strings.HasPrefix(s, "Ok") {
			return true
		}
		if strings.HasPrefix(s, "BAD") {
			return false
		}
	}

}

var ch changer

func PhoneTest() {
	soc, err := net.Dial("tcp", "localhost:"+strconv.Itoa(setup.Set.Port))
	if err != nil {
		logger.Error.Printf("tcp dial %s", err.Error())
		return
	}
	defer soc.Close()
	ph := new(dataBase.Phone)
	ph.Login = "rura"
	ph.Password = "162747"
	ph.Areas = make([]int, 0)
	ph.Status = dataBase.Status{Connection: false}
	ch.ph = ph
	ch.ph.GetPhone()
	ch.ph.SetPhone()

	ch.writer = bufio.NewWriter(soc)
	ch.reader = bufio.NewReader(soc)

	fmt.Println("Старт тесты")
	if !ch.readDB() {
		fmt.Println(" не прошел")
		return
	} else {
		fmt.Println(" ok")
	}
	crosses := dataBase.GetCrosses(ch.ph.Areas)
	for _, c := range crosses {
		if !ch.setConnect(c) {
			fmt.Println(" не прошел")
			return
		} else {
			fmt.Println(" ok")
		}
		if !ch.setPhase(c, 5, 1) {
			fmt.Println(" не прошел")
			return
		} else {
			fmt.Println(" ok")
		}
		if !ch.setPhase(c, 5, 5) {
			fmt.Println(" не прошел")
			return
		} else {
			fmt.Println(" ok")
		}
		if !ch.disConnect(c) {
			fmt.Println(" не прошел")
			return
		} else {
			fmt.Println(" ok")
		}

	}

	fmt.Println("Все тесты ok")

}
