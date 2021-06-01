package memDB

import (
	"database/sql"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/ag-server/extcon"
	"github.com/ruraomsk/device/dataBase"
	"github.com/ruraomsk/device/setup"
	"time"
)

var memoryDB []*Tx
var db *sql.DB
var id int

func Start(ready chan interface{}) {
	memoryDB = make([]*Tx, 0)
	db, id = dataBase.GetDB()
	initPhone()
	initCrosses()
	memoryDB = append(memoryDB, PhonesTable)
	memoryDB = append(memoryDB, CrossesTable)
	for _, t := range memoryDB {
		t.mdb.data = t.ReadAll()
	}
	ready <- 0
	logger.Info.Println("memDB start work")
	updateTablesTicker := time.NewTicker(time.Duration(setup.Set.DataBase.UpdateInterval) * time.Second)
	p, _ := extcon.NewContext("memDB")
	for true {
		select {
		case <-updateTablesTicker.C:
			for _, t := range memoryDB {
				if t.Save() != nil {
					return
				}
			}
		case <-p.Done():
			for _, t := range memoryDB {
				_ = t.Save()
			}
			dataBase.FreeDB(id)
			logger.Info.Println("memDB end work")
			return

		}
	}
}
