package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	_defaultMaxIdleConns = 10
	_defaultMaxOpenConns = 20
)

type Mysql struct {
	maxIdleConns int
	maxOpenConns int
	Engine       *xorm.Engine
}

func New(url string, opts ...Option) (*Mysql, error) {

	mysql := &Mysql{
		maxIdleConns: _defaultMaxIdleConns,
		maxOpenConns: _defaultMaxOpenConns,
	}

	for _, opt := range opts {
		opt(mysql)
	}

	var err error

	mysql.Engine, err = xorm.NewEngine("mysql", url)

	if err != nil {
		return nil, fmt.Errorf("mysql - NewMySQL -NewEngine: %w", err)
	}

	mysql.Engine.SetMaxIdleConns(mysql.maxIdleConns)
	mysql.Engine.SetMaxOpenConns(mysql.maxOpenConns)

	if err = mysql.Engine.DB().Ping(); err != nil {
		return nil, fmt.Errorf("mysql - NewMySQL - Ping == 0: %w", err)
	}
	return mysql, nil
}

func (m *Mysql) Close() {
	m.Engine.Close()
}
