package util

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Flags struct {
	dbHost   string
	dbPort   int
	dbName   string
	dbUser   string
	dbPasswd string
	port     int
	dev      bool
}

func (f *Flags) DbConnectString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		f.dbUser, f.dbPasswd, f.dbHost, f.dbPort, f.dbName)
}

func (f *Flags) ServeAddr() string {
	return fmt.Sprintf(":%d", f.port)
}

func (f *Flags) Port() int {
	return f.port
}

func (f *Flags) InvokeForTesting() bool {
	return f.dev
}

func (f *Flags) isOk() bool {
	if f.dbHost == "" || f.dbPort == 0 || f.dbName == "" || f.dbUser == "" || f.dbPasswd == "" {
		return false
	}
	return true
}

func obtainFromFlag(f *Flags) {
	flag.StringVar(&f.dbHost, "dbhost", "", "database host")
	flag.IntVar(&f.dbPort, "dbport", 0, "database port")
	flag.StringVar(&f.dbName, "dbname", "", "database name")
	flag.StringVar(&f.dbUser, "dbuser", "", "database user login")
	flag.StringVar(&f.dbPasswd, "dbpassword", "", "database user password")
	flag.IntVar(&f.port, "port", 5000, "server port")
	flag.BoolVar(&f.dev, "dev", false, "enable dev endpoints")
	flag.Parse()
}

func obtainFromEnv(f *Flags) {
	f.dbHost = os.Getenv("DB_HOST")
	f.dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	f.dbUser = os.Getenv("DB_USER")
	f.dbPasswd = os.Getenv("DB_PASSWORD")
	f.dbName = os.Getenv("DB_NAME")
}

func ParseParameters() (*Flags, error) {
	var param Flags

	obtainFromFlag(&param)
	if param.isOk() {
		return &param, nil
	}

	obtainFromEnv(&param)
	if param.isOk() {
		return &param, nil
	}

	return nil, errors.New("cannot parse parameters")

}
