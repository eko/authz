package configs

type Database struct {
	Host     string `config:"database_host"`
	Port     string `config:"database_port"`
	SSLMode  string `config:"database_ssl"`
	User     string `config:"database_user"`
	Password string `config:"database_password"`
	Dbname   string `config:"database_name"`
	Timezone string `config:"database_timezone"`
}

func (d Database) DSN() string {
	return "host=" + d.Host +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.Dbname +
		" port=" + d.Port +
		" sslmode=" + d.SSLMode +
		" TimeZone=" + d.Timezone
}

func newDatabase() *Database {
	return &Database{
		Host:     "localhost",
		User:     "root",
		Password: "toor",
		Dbname:   "root",
		Port:     "5432",
		SSLMode:  "disable",
		Timezone: "Europe/Paris",
	}
}
