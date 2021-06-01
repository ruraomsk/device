package memDB

import (
	"encoding/json"
	"fmt"
	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/device/dataBase"
	"golang.org/x/crypto/bcrypt"
)

var PhonesTable *Tx

func GetPhone(login string) (dataBase.Phone, error) {
	value, err := PhonesTable.Get(login)
	if err != nil {
		return dataBase.Phone{}, err
	}
	return value.(dataBase.Phone), nil
}
func IsPhoneCorrect(login string, password string) (bool, []int) {
	var user dataBase.Phone
	user, _ = GetPhone(login)
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return false, make([]int, 0)
	}
	return true, user.Areas
}
func SetPhone(phone dataBase.Phone) {
	PhonesTable.Set(phone.Login, phone)
}
func initPhone() {
	PhonesTable = Create()
	PhonesTable.name = "phones"
	PhonesTable.ReadAll = func() map[string]interface{} {
		res := make(map[string]interface{})
		w := fmt.Sprintf("select  login,phone from phones;")
		rows, err := db.Query(w)
		if err != nil {
			logger.Error.Printf("запрос %s %s", w, err.Error())
			return res
		}
		var key string
		var buffer []byte
		var phone dataBase.Phone
		for rows.Next() {
			_ = rows.Scan(&key, &buffer)
			_ = json.Unmarshal(buffer, &phone)
			res[key] = phone
		}
		return res
	}
	PhonesTable.UpdateFn = func(key string, value interface{}) string {
		val, _ := json.Marshal(value)
		return fmt.Sprintf("update phones set phone='%s' where login='%s';", string(val), key)
	}
	PhonesTable.AddFn = func(key string, value interface{}) string {
		val, _ := json.Marshal(value)
		return fmt.Sprintf("insert into phones (login,phone) values ('%s','%s');", key, string(val))
	}
	PhonesTable.DeleteFn = func(key string) string {
		return fmt.Sprintf("delete from phones where logon='%s';", key)
	}
	return
}
