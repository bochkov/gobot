package util

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type DbParams struct {
	host   string
	port   int
	dbname string
	user   string
	passwd string
}

func (p *DbParams) ConnectString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		p.user, p.passwd, p.host, p.port, p.dbname)
}

func (p *DbParams) isOk() bool {
	if p.host == "" || p.port == 0 || p.dbname == "" || p.user == "" || p.passwd == "" {
		return false
	}
	return true
}

func obtainFromFlag(p *DbParams) {
	flag.StringVar(&p.host, "host", "", "database host")
	flag.IntVar(&p.port, "port", 0, "database port")
	flag.StringVar(&p.dbname, "db", "", "database name")
	flag.StringVar(&p.user, "user", "", "database user login")
	flag.StringVar(&p.passwd, "password", "", "database user password")
	flag.Parse()
}

func obtainFromEnv(p *DbParams) {
	p.host = os.Getenv("DB_HOST")
	p.port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	p.user = os.Getenv("DB_USER")
	p.passwd = os.Getenv("DB_PASSWORD")
	p.dbname = os.Getenv("DB_NAME")
}

func ParseDbParameters() (*DbParams, error) {
	var param DbParams

	obtainFromFlag(&param)
	if param.isOk() {
		return &param, nil
	}

	obtainFromEnv(&param)
	if param.isOk() {
		return &param, nil
	}

	return nil, errors.New("cannot get database parameters")

}
