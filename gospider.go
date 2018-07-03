package gospider

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nange/gospider/common"
	"github.com/nange/gospider/web"
	"github.com/pkg/errors"
)

const (
	Name    = "gospider"
	Version = "1.0.0"
)

type GoSpider struct {
	backend string
	mysql   common.MySQLConf
	web     *web.Server
}

func New(opts ...func(*GoSpider)) *GoSpider {
	gs := &GoSpider{}
	gs.init()

	for _, f := range opts {
		f(gs)
	}

	gs.parseSettingsFromEnv()

	return gs
}

// 启动gosipder
func (gs *GoSpider) Run() error {
	gs.print()
	db, err := common.NewGormDB(gs.mysql)
	if err != nil {
		return errors.Wrap(err, "new gorm db failed")
	}
	gs.web.SetDB(db)

	return errors.Wrap(gs.web.Run(), "web run failed")
}

func BackendMySQL() func(*GoSpider) {
	return func(gs *GoSpider) {
		gs.backend = "mysql"
	}
}

func BackendSQLite() func(*GoSpider) {
	return func(gs *GoSpider) {
		gs.backend = "sqlite"
	}
}

func MySQLHost(host string) func(*GoSpider) {
	return func(gs *GoSpider) {
		gs.mysql.Host = host
	}
}

func MySQLPort(port int) func(*GoSpider) {
	return func(gs *GoSpider) {
		gs.mysql.Port = port
	}
}

func MySQLUser(user string) func(*GoSpider) {
	return func(gs *GoSpider) {
		gs.mysql.User = user
	}
}

func MySQLPassword(password string) func(*GoSpider) {
	return func(gs *GoSpider) {
		gs.mysql.Password = password
	}
}

func MySQLDBName(dbname string) func(*GoSpider) {
	return func(gs *GoSpider) {
		gs.mysql.DBName = dbname
	}
}

func WebIP(ip string) func(*GoSpider) {
	return func(gs *GoSpider) {
		gs.web.IP = ip
	}
}

func WebPort(port int) func(*GoSpider) {
	return func(gs *GoSpider) {
		gs.web.Port = port
	}
}

func (gs *GoSpider) init() {
	gs.backend = "mysql"
	gs.mysql.Host = "127.0.0.1"
	gs.mysql.Port = 3306
	gs.mysql.User = "root"
	gs.mysql.MaxIdleConns = 3
	gs.mysql.MaxOpenConns = 10
	gs.mysql.DBName = "gospider"

	gs.web = &web.Server{Port: 8080}
}

func (gs *GoSpider) print() {
	log.Println(Name, Version)
	log.Printf("gospider backend conf:%+v\n", gs.mysql)
}

var envMap = map[string]func(*GoSpider, string){
	"DB_HOST": func(gs *GoSpider, val string) {
		gs.mysql.Host = val
	},
	"DB_PORT": func(gs *GoSpider, val string) {
		port, err := strconv.Atoi(val)
		if err == nil {
			gs.mysql.Port = port
		}
	},
	"DB_USER": func(gs *GoSpider, val string) {
		gs.mysql.User = val
	},
	"DB_PASSWORD": func(gs *GoSpider, val string) {
		gs.mysql.Password = val
	},
	"DB_NAME": func(gs *GoSpider, val string) {
		gs.mysql.DBName = val
	},
	"WEB_IP": func(gs *GoSpider, val string) {
		gs.web.IP = val
	},
	"WEB_PORT": func(gs *GoSpider, val string) {
		port, err := strconv.Atoi(val)
		if err == nil {
			gs.web.Port = port
		}
	},
}

func (gs *GoSpider) parseSettingsFromEnv() {
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "GOSPIDER_") {
			continue
		}
		pair := strings.SplitN(e[9:], "=", 2)
		if f, ok := envMap[pair[0]]; ok {
			f(gs, pair[1])
		} else {
			log.Println("Unknown env variable:", pair[0])
		}
	}
}
