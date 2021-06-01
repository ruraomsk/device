package dataBase

import (
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/device/setup"
)

func updateArea() {
	if !setup.Set.External.Make {
		return
	}
	if !work {
		return
	}
	db, idb := GetDB()
	defer FreeDB(idb)
	tx, err := db.Begin()
	_, err = tx.Exec("delete from area;")
	if err != nil {
		logger.Error.Printf("delete from area; %s", err.Error())
		work = false
		return
	}
	_ = tx.Commit()
	w := fmt.Sprintf("select area,namearea from region where region=%d; ", setup.Set.Region)
	rows, err := conExternal.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		work = false
		return
	}
	var (
		area     int
		namearea string
	)
	for rows.Next() {
		_ = rows.Scan(&area, &namearea)
		ww := fmt.Sprintf("insert into area (area, namearea) values (%d, '%s');", area, namearea)
		tx, err := db.Begin()
		_, err = tx.Exec(ww)
		if err != nil {
			logger.Error.Printf("запрос %s %s", w, err.Error())
			work = false
			return
		}
		_ = tx.Commit()
	}
	_ = rows.Close()
}
