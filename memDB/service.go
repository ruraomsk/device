package memDB

import (
	"github.com/ruraomsk/TLServer/logger"
)

func (tx *Tx) Save() error {
	tx.Lock()
	defer tx.Unlock()
	dbTx, err := db.Begin()
	if err != nil {
		logger.Error.Printf("запрос begin %s", err.Error())
		return err
	}
	for key := range tx.deleted {
		w := tx.DeleteFn(key)
		_, err = dbTx.Exec(w)
		if err != nil {
			_ = dbTx.Rollback()
			logger.Error.Printf("запрос %s %s", w, err.Error())
			return err
		}
	}
	tx.deleted = make(map[string]bool)
	for key := range tx.added {
		w := tx.AddFn(key, tx.mdb.data[key])
		_, err = dbTx.Exec(w)
		if err != nil {
			_ = dbTx.Rollback()
			logger.Error.Printf("запрос %s %s", w, err.Error())
			return err
		}
	}
	tx.added = make(map[string]bool)
	for key := range tx.updated {
		w := tx.UpdateFn(key, tx.mdb.data[key])
		_, err = dbTx.Exec(w)
		if err != nil {
			_ = dbTx.Rollback()
			logger.Error.Printf("запрос %s %s", w, err.Error())
			return err
		}
	}
	tx.updated = make(map[string]bool)
	_ = dbTx.Commit()
	tx.mdb.data = tx.ReadAll()
	return nil
}
