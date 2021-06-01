package dumper

import (
	"github.com/jasonlvhit/gocron"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/device/setup"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func makeDump() {
	if runtime.GOOS == "linux" {
		logger.Error.Printf("Нет реализации для Linux!")
		return
	}
	file, err := os.Create("save.bat")
	if err != nil {
		logger.Error.Printf("Не могу создать файл save.bat %s", err.Error())
		return
	}
	_, _ = file.WriteString("SET PGPASSWORD=" + setup.Set.DataBase.Password + "\n")
	date := time.Now().Format(time.RFC3339)[0:10]
	path := setup.Set.Dumper.Path + "/dump" + date + ".sql"
	_, _ = file.WriteString("pg_dump -U" + setup.Set.DataBase.User + " -d" + setup.Set.DataBase.DBname +
		" -C -c --column-inserts --if-exists --no-comments -f" + path + "\n")
	_ = file.Close()
	time.Sleep(5 * time.Second)
	cmd := exec.Command("save.bat")
	err = cmd.Run()
	if err != nil {
		logger.Error.Printf("Не могу выполнить save.bat %s", err.Error())
		return
	}
	logger.Info.Printf("Dump writed..")

}
func Start() {
	if !setup.Set.Dumper.Make {
		logger.Info.Printf("Dumper dont start..")
		return
	}
	logger.Info.Printf("Dumper starting..")
	_ = gocron.Every(1).Day().At(setup.Set.Dumper.Time).Do(makeDump)
	<-gocron.Start()
	logger.Info.Printf("Dumper working..")
}
