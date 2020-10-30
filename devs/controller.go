package devs

import (
	"strings"
	"time"
)

var NotCommands = []string{"#SYS.OPTS", "#FWBL", "#DBG.SWVER", "#SYS.PBASE", "SYS.SSTAT"}

var SaveController = []string{"#DBG.RESET", "#DBG.STCON"}

var SetTime = []string{"#SYS.SETTM"}
var InsertTime = []string{"#TCH.TCSTA", "#TCH.DKSTA", "#TCH.HARDW", "#TCH.PHASE", "#TCH.PANEL",
	"#TCH.TVPST", "#TCH.SENST", "#TCH.SSTAT", "#TCH.INPST", "#TCH.MGRST"}

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

	TchTcsta string `aura:"#TCH.TCSTA"` // Потом идет ТехРежим,ПК,СК,НК,00
	TchDksta string `aura:"#TCH.DKSTA"` // Потом идет Устр,Режим,Состояние
	TchHardw string `aura:"#TCH.HARDW"` // Потом идет Неисправность
	TchPhase string `aura:"#TCH.PHASE"` // Потом идет PH, TU, TTU, TS, TTS, POS, NXT, ST
	//PH – текущая фаза по плану
	//TU – фактическая фаза ТУ
	//TTU – время фазы ТУ
	//TS – фактическая фаза ТС
	//TTS – время фазы ТС
	//POS – позиция в плане (секунда от начала)
	//NXT – позиция следующей смены фазы
	//SW – номер переключения от начала плана
	TchPanel string `aura:"#TCH.PANEL"` // NN – Входы (NN – HEX-число, каждый бит отвечает за
	//активность входа)
	//Бит 0 – Тумблер ЖМ включен (ЖМ)
	//Бит 1 – Тумблер ОС включен (ОС)
	//Бит 2 – Автомат &quot;Светофоры&quot; отключен (НАГРУЗ)
	//Бит 3 – Автомат &quot;Сеть&quot; выключен (ВВОД)
	//Бит 4 – Тумблер БП выключен (СЕТЬ)
	//Бит 5 – Дверь открыта (ДВЕРЬ)
	//Бит 6 – ТВП 1 замкнут (ТВП1)
	//Бит 7 – ТВП 2 замкнут (ТВП2)
	TchTvpst string `aura:"#TCH.TVPST"` // NN – ТВП
	//Бит 0 – Состояние ТВП1
	//Бит 1 – Состояние ТВП1
	//Бит 2 – Состояние ТВП2
	//Бит 3 – Состояние ТВП2
	//
	//Бит 4 – ТВП1 лампа &quot;Ждите&quot;
	//Бит 5 – ТВП2 лампа &quot;Ждите&quot;
	//Бит 6 – ТВП 1 неисправность входа
	//Бит 7 – ТВП 2 неисправность входа
	//Состояние ТВП1,2
	//00 – Готов (готовность к приему заявки)
	//01 – Вызов (отрабатывается принятая заявка)
	//02 – Ждите (ожидание отработки принятой заявки)
	TchSenst string `aura:"#TCH.SENST"` // D1, D2, D3, D4, D5, D6, D7, D8, D9, D10, D11, D12, D13,
	//	D14, D15, D16, S1, S2, S3, S4, S5, S6, S7, S8, S9, S10, S11, S12, S13, S14, S15, S16
	//	D1-D15 – данные датчиков за текущую секунду (количество ТС)
	//	S1-S15 – текущая накопленная статистика с начала интервала усреднения
	//(Количество ТС)
	TchSstat string `aura:"#TCH.SSTAT"` // S01, S02, S03, S04, S05, S06, S07, S08, S09, S10, S11,
	//S12, S13, S14, S15
	//S01 – S15 – переданная на сервер статистика (количество ТС по истечении
	//интервала или по команде)
	TchInpst string `aura:"#TCH.INPST"` // INT, EXT, TVP, VPU, MGR1, MGR2, FLT
	//INT – интервал усреднения статистики (минуты)
	//EXT – привязанные входы (HEX побитно)
	//TVP – входы в режиме ТВП (HEX побитно)
	//VPU – входы в режиме ВПУ (HEX побитно)
	//MGR1 – входы в режиме МГР ДК1 (HEX побитно)
	//MGR2 – входы в режиме МГР ДК2 (HEX побитно – не используются, но есть в
	//привязке)
	//FLT – неисправные входы (HEX побитно)
	TchMgrst string `aura:"#TCH.MGRST"` // T01, T02, T03, T04, T05, T06, T07, T08, XXXX
	//T01-T08 – таймеры продления фаз при МГР
	//XXXX – признаки неисправности входов датчиков по битам
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
func isNeedInsertTime(result string) bool {
	ls := strings.Split(result, ":")
	for _, cc := range InsertTime {
		if strings.Compare(ls[0], cc) == 0 {
			return true
		}
	}
	return false
}
func insertTime(result string) string {
	ls := strings.Split(result, ":")
	result = ls[0] + ":" + time.Now().Format("15:04:05 ")
	for i := 1; i < len(ls); i++ {
		result += ls[i]
	}
	return result
}
