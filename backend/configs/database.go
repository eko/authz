package configs

import "fmt"

const (
	DriverPostgres = "postgres"
	DriverMysql    = "mysql"
	DriverSqlite   = "sqlite"
)

type Database struct {
	Driver   string `config:"database_driver"`
	Host     string `config:"database_host"`
	Port     string `config:"database_port"`
	SSLMode  string `config:"database_ssl"`
	User     string `config:"database_user"`
	Password string `config:"database_password"`
	Dbname   string `config:"database_name"`
	Timezone string `config:"database_timezone"`
}

func (d Database) MysqlDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Dbname,
	)
}

func (d Database) PostgresDSN() string {
	return "host=" + d.Host +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.Dbname +
		" port=" + d.Port +
		" sslmode=" + d.SSLMode +
		" TimeZone=" + d.Timezone
}

func (d Database) SqliteDSN() string {
	return fmt.Sprintf("file:%s?cache=shared&mode=rwc&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)", d.Dbname)
}

func newDatabase() *Database {
	return &Database{
		Driver:   "postgres",
		Host:     "localhost",
		User:     "root",
		Password: "toor",
		Dbname:   "root",
		Port:     "5432",
		SSLMode:  "disable",
		Timezone: "UTC",
	}
}
