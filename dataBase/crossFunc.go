package dataBase

import (
	"encoding/json"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
)

func GetOneCross(key string) Cross {
	var cross Cross
	db, id := GetDB()
	defer FreeDB(id)
	w := fmt.Sprintf("select crossT from crosses where key='%s';", key)
	rows, err := db.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		return cross
	}
	var buf []byte
	for rows.Next() {
		rows.Scan(&buf)
		_ = json.Unmarshal(buf, &cross)
		return cross

	}
	return cross
}

func GetCrosses(areas []int) []Cross {
	result := make([]Cross, 0)
	var cross Cross
	db, id := GetDB()
	defer FreeDB(id)
	w := fmt.Sprintf("select crossT from crosses;")
	rows, err := db.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		return result
	}
	var buf []byte
	for rows.Next() {
		rows.Scan(&buf)
		_ = json.Unmarshal(buf, &cross)
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
func IsCross(key string) bool {
	db, id := GetDB()
	defer FreeDB(id)
	w := fmt.Sprintf("select count(*) from crosses where key='%s';", key)
	rows, err := db.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		return false
	}
	var count int
	for rows.Next() {
		rows.Scan(&count)
		if count == 0 {
			return false
		}
		return true
	}
	return false
}
