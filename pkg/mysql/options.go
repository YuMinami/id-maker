package mysql

type Option func(mysql *Mysql)

func MaxIdleConns(size int) Option {
	return func(c *Mysql) {
		c.maxIdleConns = size
	}
}

func MaxOpenConns(size int) Option {
	return func(mysql *Mysql) {
		mysql.maxOpenConns = size
	}
}
