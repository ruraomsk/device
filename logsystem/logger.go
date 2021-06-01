package logsystem

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/ag-server/extcon"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/device/dataBase"
	"github.com/ruraomsk/device/memDB"
	"github.com/ruraomsk/device/setup"
	"strconv"
	"strings"
	"time"
)

type LoggData struct {
	External bool   `json:"ext"`
	Login    string `json:"login"`
	Key      string `json:"key"`
	Txt      string `json:"txt"`
}

var LoggChan chan LoggData

func Start(ready chan interface{}) {
	db, id := dataBase.GetDB()
	defer dataBase.FreeDB(id)
	var con *sql.DB
	var err error
	if setup.Set.External.Make {
		// Есть внешняя БД проверим ее
		info := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			setup.Set.External.Host, setup.Set.External.User,
			setup.Set.External.Password, setup.Set.External.DBname)
		con, err = sql.Open("postgres", info)
		if err != nil {
			logger.Error.Printf("запрос на открытие %s %s", info, err.Error())
			return
		}
	}
	LoggChan = make(chan LoggData, 1000)
	ready <- 0
	p, _ := extcon.NewContext("logSystem")
	logger.Info.Println("logSystem start work")
	for {
		select {
		case l := <-LoggChan:
			{

				w := fmt.Sprintf("INSERT INTO logger(tm, login, key, txt) VALUES ('%s','%s' ,'%s' ,'%s');",
					string(pq.FormatTimestamp(time.Now())), l.Login, l.Key, l.Txt)
				//logger.Debug.Println(w)
				tx, err := db.Begin()
				_, err = tx.Exec(w)
				if err != nil {
					logger.Error.Printf("запрос %s %s", w, err.Error())
					return
				}
				err = tx.Commit()
				if err != nil {
					logger.Error.Printf("запрос %s %s", w, err.Error())
					return
				}
				if setup.Set.External.Make && l.External {
					// Есть внешняя БД запишем туда данные
					ks := strings.Split(l.Key, ":")
					id, _ := strconv.Atoi(ks[1])
					cross, _ := memDB.GetCross(l.Key)
					extLogg := pudge.JSONLog{Region: strconv.Itoa(setup.Set.Region), Area: ks[0], ID: id, Description: cross.Name, Type: 0}
					result, _ := json.Marshal(extLogg)
					w := fmt.Sprintf("insert into public.logdevice (id,tm,crossinfo,txt) values(%d,'%s','%s','%s');",
						cross.IDevice, string(pq.FormatTimestamp(time.Now())), result, l.Txt)
					tx, err := con.Begin()
					_, err = tx.Exec(w)
					if err != nil {
						logger.Error.Printf("запрос %s %s", w, err.Error())
						return
					}
					err = tx.Commit()
					if err != nil {
						logger.Error.Printf("запрос %s %s", w, err.Error())
						return
					}

				}

			}
		case <-p.Done():
			logger.Info.Println("logSystem end work")
			return
		}

	}
}
