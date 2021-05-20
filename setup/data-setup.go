package setup

//Set переменная для хранения текущих настроек
var Set *Setup

//Setup общая структура для настройки всей системы
type Setup struct {
	Home     string   `toml:"home"`
	Region   int      `toml:"region"`
	Location string   `toml:"location"` //Локация временной зоны
	DataBase DataBase `toml:"dataBase"`
	Port     int      `toml:"port"` //Стартовый номер порта на прием
	External External `toml:"external"`
	Dumper   Dumper   `toml:"dumper"`
	Default  Default  `toml:"default"`
}
type Default struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	SSID     string `toml:"ssid"`     //SSID wifi с контроллером
	PassID   string `toml:"passid"`   //password для соединения с локальной сетью
	Login    string `toml:"login"`    //login для соединения с контроллером
	Password string `toml:"password"` //password для соединения с контроллером
}
type External struct {
	Make     bool   `toml:"make"` //true если есть внешний сервер ASDU
	Step     int    `toml:"step"` //Интервал времени в секундах для расчетов
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBname   string `toml:"dbname"`
}

//DataBase настройки базы данных postresql
type DataBase struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	User        string `toml:"user"`
	Password    string `toml:"password"`
	DBname      string `toml:"dbname"`
	MaxOpenConn int    `toml:"maxopen"`
}

type Dumper struct {
	Make bool   `toml:"make"`
	Path string `toml:"path"`
	Time string `toml:"time"`
}
