package devs

import (
	"strings"
	"time"
)

var NotCommands []string = []string{"#SYS.OPTS", "#FWBL", "#DBG.SWVER", "#SYS.PBASE", "SYS.SSTAT"}

var SaveController []string = []string{"#DBG.RESET", "#DBG.STCON"}

var SetTime []string = []string{"#SYS.SETTM"}

type Controller struct {
	Time time.Time //Время
	Info Info      //Информация об устройстве
}

//
type Info struct {
	Nomer        string `aura:"#BRD.IDNUM"` //
	Hardware     string `aura:"#BRD.HWVER"` //
	Software     string `aura:"#BRD.SWVER"` //
	IpserverGPRS string `aura:"#SRV.IPADR"` // команда и ответ
	IpserverLAN  string `aura:"#LAN.SRVIP"` // команда и ответ
	Ipcontroller string `aura:"#LAN.IPADR"` // команда и ответ
	Mask         string `aura:"#LAN.MASK"`  //  команда и ответ
	Gate         string `aura:"#LAN.GATE"`  //  команда и ответ
	Lanstate     string `aura:"#LAN.STATE"` //
	Lanerror     string `aura:"#LAN.ERROR"` //
	Chanel       string `aura:"#SYS.CHAN"`  //
	Mcode        string `aura:"#SYS.MCODE"` //
	Gpserror     string `aura:"#GPS.ERROR"` //
	Error485     string `aura:"#485.ERROR"` //
	Pwrstate     string `aura:"#PWR.STATE"` //
	Memerror     string `aura:"#MEM.ERROR"` //
	Syssync      string `aura:"#SYS.SYNC"`  //
	Gpssat       string `aura:"#GPS.SAT"`   //
	Syssntm      string `aura:"#SYS.SNTM"`  //
	Systime      string `aura:"#SYS.TIME"`  //
	GpsTime      string `aura:"#GPS.TIME"`  //
	Modtype      string `aura:"#MOD.TYPE"`  //
	Modfwver     string `aura:"#MOD.FWVER"` //
	Gsmopera     string `aura:"#GSM.OPERA"` //
	Gsmlevel     string `aura:"#GSM.LEVEL"` //
	Sysdelay     string `aura:"#SYS.DELAY"` //
	Gsmsrvip     string `aura:"#GSM.SRVIP"` //
	Gsmstate     string `aura:"#GSM.STATE"` //
	Gsmerror     string `aura:"#GSM.ERROR"`
}
type Command struct {
	IpserverGPRS string `aura:"#SRV.IPADR"` // команда и ответ
	IpserverLAN  string `aura:"#LAN.SRVIP"` // команда и ответ
	Ipcontroller string `aura:"#LAN.IPADR"` // команда и ответ
	Mask         string `aura:"#LAN.MASK"`  //  команда и ответ
	Gate         string `aura:"#LAN.GATE"`  //  команда и ответ
	Mode         string `aura:"#DBG.MODE"`
}

//

func isEmpty(command string) bool {
	for _, cc := range NotCommands {
		if strings.Compare(command, cc) == 0 {
			return true
		}
	}
	return false
}
func isSave(command string) bool {
	for _, cc := range SaveController {
		if strings.Compare(command, cc) == 0 {
			return true
		}
	}
	return false
}
func isSetTime(command string) bool {
	for _, cc := range SetTime {
		if strings.Compare(command, cc) == 0 {
			return true
		}
	}
	return false
}
