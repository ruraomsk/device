package memDB

import (
	"encoding/json"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/device/dataBase"
)

var CrossesTable *Tx

func initCrosses() {
	CrossesTable = Create()
	CrossesTable.name = "crosses"
	CrossesTable.ReadAll = func() map[string]interface{} {
		res := make(map[string]interface{})
		w := fmt.Sprintf("select  key,crosst from crosses;")
		rows, err := db.Query(w)
		if err != nil {
			logger.Error.Printf("запрос %s %s", w, err.Error())
			return res
		}
		var key string
		var buffer []byte
		var cross dataBase.Cross
		for rows.Next() {
			_ = rows.Scan(&key, &buffer)
			_ = json.Unmarshal(buffer, &cross)
			res[key] = cross
		}
		return res
	}
	CrossesTable.UpdateFn = func(key string, value interface{}) string {
		val, _ := json.Marshal(value)
		return fmt.Sprintf("update crosses set crosst='%s' where key='%s';", string(val), key)
	}
	CrossesTable.AddFn = func(key string, value interface{}) string {
		val, _ := json.Marshal(value)
		return fmt.Sprintf("insert into crosses (key,crosst) values ('%s','%s');", key, string(val))
	}
	CrossesTable.DeleteFn = func(key string) string {
		return fmt.Sprintf("delete from crosses where key='%s';", key)
	}
	return
}
func GetCross(key string) (dataBase.Cross, error) {
	value, err := CrossesTable.Get(key)
	if err != nil {
		return dataBase.Cross{}, err
	}
	return value.(dataBase.Cross), err
}

func ListCrosses(areas []int) []dataBase.Cross {
	result := make([]dataBase.Cross, 0)
	var cross dataBase.Cross
	keys := CrossesTable.GetAllKeys()
	for _, key := range keys {
		cross, _ = GetCross(key)
		if len(areas) == 0 {
			result = append(result, cross)
		} else {
			for _, area := range areas {
				if area == cross.Area {
					result = append(result, cross)
					break
				}
			}
		}
	}
	return result
}
