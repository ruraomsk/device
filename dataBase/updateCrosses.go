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
	//w := fmt.Sprintf("select area,id,state from public.\"cross\" where region=%d;", setup.Set.Region)
	w := fmt.Sprintf("SELECT device,state FROM public.devices LEFT JOIN public.\"cross\" ON public.devices.id = public.\"cross\".idevice where region=%d;", setup.Set.Region)

	rows, err := conExternal.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		work = false
		return
	}
	var (
		state  string
		device string
	)
	var (
		st  pudge.Cross
		dv  pudge.Controller
		buf []byte
		cr  Cross
	)
	update := make([]string, 0)
	for rows.Next() {
		_ = rows.Scan(&device, &state)
		_ = json.Unmarshal([]byte(state), &st)
		_ = json.Unmarshal([]byte(device), &dv)
		key := fmt.Sprintf("%d:%d", st.Area, st.ID)
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
			cr.Default(dv.IPHost)
			buf, _ = json.Marshal(&cr)
			www := fmt.Sprintf("insert into crosses (key,crossT) values ('%s','%s')", key, string(buf))
			update = append(update, www)

		}
	}
	for _, w := range update {
		tx, err := db.Begin()
		_, err = tx.Exec(w)
		if err != nil {
			logger.Error.Printf("запрос %s %s", w, err.Error())
			work = false
			return
		}
		_ = tx.Commit()
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
		tx, err := db.Begin()
		_, err = tx.Exec(w)
		if err != nil {
			logger.Error.Printf("запрос %s %s", w, err.Error())
			work = false
			return
		}
		_ = tx.Commit()
	}
}
