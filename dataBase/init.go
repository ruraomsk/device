package dataBase

import (
	"database/sql"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/device/setup"
	"sync"
	"time"
)

var conExternal *sql.DB

type usedDb struct {
	db   *sql.DB
	used bool
}

var dbPool []usedDb
var mutex sync.Mutex
var err error
var work bool
var info string

func Init() error {
	// Проверяем наличие всех таблиц для работы
	if setup.Set.External.Make {
		// Есть внешняя БД проверим ее
		info = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			setup.Set.External.Host, setup.Set.External.User,
			setup.Set.External.Password, setup.Set.External.DBname)
		conExternal, err = sql.Open("postgres", info)
		if err != nil {
			return fmt.Errorf("запрос на открытие %s %s", info, err.Error())
		}
	}
	info = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		setup.Set.DataBase.Host, setup.Set.DataBase.User,
		setup.Set.DataBase.Password, setup.Set.DataBase.DBname)
	con, err := sql.Open("postgres", info)
	if err != nil {
		return fmt.Errorf("запрос на открытие %s %s", info, err.Error())
	}

	defer con.Close()

	_, err = con.Exec(loggerTableCreate)
	if err != nil {
		return fmt.Errorf("запрос  %s %s", loggerTableCreate, err.Error())
	}
	_, err = con.Exec(areaTableCreate)
	if err != nil {
		return fmt.Errorf("запрос  %s %s", areaTableCreate, err.Error())
	}
	_, err = con.Exec(phoneTableCreate)
	if err != nil {
		return fmt.Errorf("запрос  %s %s", phoneTableCreate, err.Error())
	}
	_, err = con.Exec(crossesTableCreate)
	if err != nil {
		return fmt.Errorf("запрос  %s %s", crossesTableCreate, err.Error())
	}
	dbPool = make([]usedDb, 0)
	for i := 0; i < setup.Set.DataBase.MaxOpenConn; i++ {
		d, err := sql.Open("postgres", info)
		if err != nil {
			logger.Error.Printf("dbase ConnectDB %s", err.Error())
			return err
		}
		dbPool = append(dbPool, usedDb{db: d, used: false})
	}
	work = true
	updateCrosses()
	updateArea()

	if work {
		return nil
	}
	return fmt.Errorf("во время адаптации БД произли ошбки проверьте лог системы")
}

//GetDB обращение к БД
func GetDB() (db *sql.DB, id int) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, d := range dbPool {
		if !d.used {
			dbPool[i].used = true
			dbPool[i].db, _ = sql.Open("postgres", info)
			dbPool[i].db.SetMaxIdleConns(0)
			dbPool[i].db.SetMaxOpenConns(1)
			dbPool[i].db.SetConnMaxIdleTime(1 * time.Second)
			dbPool[i].db.SetConnMaxLifetime(5 * time.Second)
			return dbPool[i].db, i
		}
	}
	logger.Error.Printf("dbase закончился пул соединений")
	return nil, 0
}
func FreeDB(id int) {
	mutex.Lock()
	defer mutex.Unlock()
	if id < 0 || id >= len(dbPool) {
		logger.Error.Printf("dbase freeDb неверный индекс %d", id)
		return
	}
	dbPool[id].used = false
	dbPool[id].db.Close()
}
