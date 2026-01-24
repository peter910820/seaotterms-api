package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	seaottermsdb "seaotterms-db"

	"seaotterms-api/teach"
)

var (
	// store
	blogStore *session.Store
	// management database connect
	dbm = make(map[string]*seaottermsdb.DBModel)
)

func init() {
	// init logrus settings
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	// init env file
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file load error: %v", err)
	}

	// init store
	blogStore = session.New(session.Config{
		Expiration:     14 * 24 * time.Hour,
		CookieDomain:   os.Getenv("SESSION_DOMAIN"),
		CookieSameSite: "None",
		KeyLookup:      "cookie:blog-userinfo-session",
		// CookieHTTPOnly: true,
	}) // blog user session

	dbName := strings.Split(os.Getenv("DATABASE_NAME"), ",")
	dbPort, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		logrus.Fatalf("db port parse error: %v", err)
	}

	// init db
	for _, n := range dbName {
		dbModel, err := seaottermsdb.InitDsn(seaottermsdb.ConnectDBConfig{
			Owner:    os.Getenv("DATABASE_OWNER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			DBName:   strings.TrimSpace(n),
			Port:     dbPort,
		})
		if err != nil {
			logrus.Fatal(err)
		}

		dbm[n] = dbModel
		// init migration
		seaottermsdb.Migration(dbModel)
	}
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_URL"),
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept",
		AllowCredentials: true,
	}))
	// api route group
	apiGroup := app.Group("/api") // main api route group

	// site route group
	// 根據資料庫註冊指定路由(因為目前設計是站台以及資料庫是一對一，一個站台一組根路由)
	// 多租戶模式確保不同站台不會衝突
	for _, m := range dbm {
		switch m.GetDBModel() {
		// case seaottermsdb.BlogModel:
		// 	blogrouter.BlogRouter(apiGroup, m)
		// case seaottermsdb.DiscordBotModel:
		// 	teachrouter.TeachRouter(apiGroup, m)
		// case seaottermsdb.AuthModel:
		// 	teachrouter.TeachRouter(apiGroup, m)
		case seaottermsdb.TeachModel:
			teach.TeachRouter(apiGroup, m)
		default:
		}
	}

	logrus.Fatal(app.Listen(fmt.Sprintf("127.0.0.1:%s", os.Getenv("PRODUCTION_PORT"))))
}
