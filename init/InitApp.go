package initapp

import (
	"hoho-framework-v2/adapters/repository"
	"hoho-framework-v2/infrastructure/cache"
	"hoho-framework-v2/infrastructure/database/connection"
	"hoho-framework-v2/infrastructure/router"
	"hoho-framework-v2/registry"
	"log"
	"os"

	sLog "hoho-framework-v2/log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func InitApp(envPath string) (*echo.Echo, *repository.SymperOrm, cache.RedisClient) {
	godotenv.Load(envPath)
	db, err := connection.NewPostgresCon().Conn()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.CloseDB()
	rdb := cache.NewRedisClient()
	sLog.NewLogger()
	sLog.Info("Server listen at http://localhost"+":"+os.Getenv("SERVER_PORT"), map[string]interface{}{"line": sLog.Trace()})

	r := registry.NewRegistry(db, rdb)

	e := echo.New()
	e = router.NewRouter(e, r.NewAppController())
	if err := e.Start(":" + os.Getenv("SERVER_PORT")); err != nil {
		log.Fatalln(err)
	}

	return e, db, rdb
}
