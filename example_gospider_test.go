package gospider_test

import (
	"github.com/labstack/gommon/log"
	"github.com/nange/gospider"
)

// Example quitstart
func ExampleGoSpider_Run() {
	// if gospider.New() has no argments, will use env parameters
	// gs := gospider.New()

	gs := gospider.New(
		gospider.BackendMySQL(),
		gospider.MySQLHost("127.0.0.1"),
		gospider.MySQLPort(3306),
		gospider.MySQLDBName("test"),
		gospider.MySQLUser("root"),
		gospider.MySQLPassword(""),
		gospider.WebPort(8080),
	)
	log.Fatal(gs.Run())
}
