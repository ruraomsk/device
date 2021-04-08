package devs

import (
	"encoding/json"
	"github.com/ruraomsk/TLServer/logger"
	"io/ioutil"
	"os"
)

func loadDevice() bool {
	file, err := os.Open("device.json")
	if err != nil {
		logger.Info.Printf("Попытка открыть файл %s", err.Error())
		return false
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("Попытка читать файл " + err.Error())
	}

	err = json.Unmarshal(data, &perfect)
	if err != nil {
		panic("Попытка раскодировать " + err.Error())
	}
	return true
}
func saveDevice(ctrl Controller) {
	file, err := os.Create("device.json")
	if err != nil {
		panic("Попытка создать файл " + err.Error())
	}
	defer file.Close()
	data, err := json.Marshal(ctrl)
	if err != nil {
		panic("Попытка закодировать " + err.Error())
	}
	_, err = file.Write(data)
	if err != nil {
		panic("Попытка записать файл " + err.Error())
	}
}
