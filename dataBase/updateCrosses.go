package dataBase

import (
	"encoding/json"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/device/setup"
	"strings"
)

var cross map[string]pudge.Cross

func updateCrosses() {
	if !setup.Set.External.Make {
		return
	}
	if !work {
		return
	}
	db, idb := GetDB()
	defer FreeDB(idb)
	cross = make(map[string]pudge.Cross)
	w := fmt.Sprintf("select area,id,state from public.\"cross\" where region=%d;", setup.Set.Region)
	rows, err := conExternal.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		work = false
		return
	}
	var (
		area  int
		id    int
		state string
	)
	var (
		st  pudge.Cross
		buf []byte
		cr  Cross
	)
	update := make([]string, 0)
	for rows.Next() {
		_ = rows.Scan(&area, &id, &state)
		_ = json.Unmarshal([]byte(state), &st)
		key := fmt.Sprintf("%d:%d", area, id)
		cross[key] = st
		ww := fmt.Sprintf("select crossT from crosses where key='%s';", key)
		rs, err := db.Query(ww)
		if err != nil {
			logger.Error.Printf("запрос %s %s", ww, err.Error())
			work = false
			return
		}
		found := false
		for rs.Next() {
			found = true
		}
		_ = rs.Close()
		if !found {
			cr.Area = st.Area
			cr.ID = st.ID
			cr.Name = st.Name
			cr.IDevice = st.IDevice
			cr.Default()
			buf, _ = json.Marshal(&cr)
			www := fmt.Sprintf("insert into crosses (key,crossT) values ('%s','%s')", key, string(buf))
			update = append(update, www)

		}
	}
	for _, w := range update {
		_, err := db.Exec(w)
		if err != nil {
			logger.Error.Printf("запрос %s %s", w, err.Error())
			work = false
			return
		}
	}
	_ = rows.Close()
	w = "select key,crossT from crosses;"
	rows, err = db.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		work = false
		return
	}
	update = make([]string, 0)
	var key string
	for rows.Next() {
		_ = rows.Scan(&key, &buf)
		_ = json.Unmarshal([]byte(buf), &cr)
		st, is := cross[key]
		if !is {
			www := fmt.Sprintf("delete from crosses where key='%s';", key)
			update = append(update, www)
		} else {
			if strings.Compare(st.Name, cr.Name) != 0 || st.IDevice != cr.IDevice {
				cr.Name = st.Name
				cr.IDevice = st.IDevice
				buf, _ = json.Marshal(&cr)
				www := fmt.Sprintf("update crosses set crossT='%s' where key='%s';", string(buf), key)
				update = append(update, www)
			}
		}
	}
	_ = rows.Close()
	for _, w := range update {
		_, err = db.Exec(w)
		if err != nil {
			logger.Error.Printf("запрос %s %s", w, err.Error())
			work = false
			return
		}
	}
}
