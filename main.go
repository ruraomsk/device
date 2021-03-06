package main

import (
	"fmt"
	"github.com/ruraomsk/device/dataBase"
	"github.com/ruraomsk/device/dumper"
	"github.com/ruraomsk/device/logsystem"
	"github.com/ruraomsk/device/memDB"
	"github.com/ruraomsk/device/sender"
	"github.com/ruraomsk/device/tests"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/ag-server/extcon"
	"github.com/ruraomsk/device/setup"
	//pprof init

	_ "net/http/pprof"
)

var err error

//Секция инициализации программы
func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFile("config.toml", &setup.Set); err != nil {
		fmt.Println("Can't load config file - ", err.Error())
	}
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	path := setup.Set.Home
	err = logger.Init(path + "/log")
	if err != nil {
		fmt.Println("Error opening logger subsystem ", err.Error())
		return
	}
	logger.Info.Println("Start devices work...")
	fmt.Println("Start devices work...")
	err = dataBase.Init()
	if err != nil {
		logger.Error.Println(err.Error())
		fmt.Println(err.Error())
		return
	}
	fmt.Println("БД адаптирована...")

	stop := make(chan interface{})
	extcon.BackgroundInit()
	ready := make(chan interface{})
	go memDB.Start(ready)
	<-ready
	go dumper.Start()
	go logsystem.Start(ready)
	<-ready
	go sender.ListenPhone()
	if err = tests.PhoneTest(); err != nil {
		return
	}
	extcon.BackgroundWork(1*time.Second, stop)
	logger.Info.Println("Exit devices working...")
	fmt.Println("\nExit devices working...")
}
