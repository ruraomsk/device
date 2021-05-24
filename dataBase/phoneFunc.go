package dataBase

import (
	"encoding/json"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"golang.org/x/crypto/bcrypt"
)

func IsPhoneCorrect(login string, password string) (bool, []int) {
	db, id := GetDB()
	defer FreeDB(id)
	result := make([]int, 0)
	w := fmt.Sprintf("select phone from phones where login='%s';", login)
	rows, err := db.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		return false, result
	}
	var user Phone
	var us []byte
	for rows.Next() {
		rows.Scan(&us)
		json.Unmarshal(us, &user)
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			return false, result
		}
		return true, user.Areas
	}
	return false, result
}

func (p *Phone) GetPhone() bool {
	db, id := GetDB()
	defer FreeDB(id)
	w := fmt.Sprintf("select phone from phones where login='%s';", p.Login)
	rows, err := db.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		return false
	}
	var us []byte
	for rows.Next() {
		rows.Scan(&us)
		json.Unmarshal(us, &p)
		return true
	}
	p.Areas = make([]int, 0)
	return false
}
func (p *Phone) SetPhone() {
	db, id := GetDB()
	defer FreeDB(id)
	w := fmt.Sprintf("select count(*) from phones where login='%s';", p.Login)
	rows, err := db.Query(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		return
	}
	var count int
	var us []byte
	for rows.Next() {
		rows.Scan(&count)
		us, _ = json.Marshal(&p)
		if count == 0 {
			w = fmt.Sprintf("insert into phones (login,phone) values ('%s','%s');", p.Login, string(us))
		} else {
			w = fmt.Sprintf("update phones set phone='%s' where login='%s';", string(us), p.Login)
		}
	}
	rows.Close()
	_, err = db.Exec(w)
	if err != nil {
		logger.Error.Printf("запрос %s %s", w, err.Error())
		return
	}
}
