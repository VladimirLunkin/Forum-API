package main

import (
	"database/sql"
	"fmt"
	"github.com/fasthttp/router"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"log"
	"tech-db-forum/internal/app/forum"
	"tech-db-forum/internal/app/service"
	"tech-db-forum/internal/app/user"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

type Config struct {
	Host string
	Port string
	DB   DBConfig
}

func (c Config) Addr() string {
	return c.Host + ":" + c.Port
}

func getPostgres(config DBConfig) *sql.DB {
	dsn := fmt.Sprintf(`user=%s dbname=%s password=%s host=%s port=%s sslmode=disable`,
		config.Username, config.DBName, config.Password, config.Host, config.Port)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10000)
	return db
}

func main() {
	viper.AddConfigPath("./config/")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	r := router.New()

	user.SetUserRouting(r, &user.Handlers{
		Logger:   logger,
		UserRepo: user.NewPgxRepository(getPostgres(config.DB)),
	})

	forum.SetForumRouting(r, &forum.Handlers{
		Logger:    logger,
		ForumRepo: forum.NewPgxRepository(getPostgres(config.DB)),
		UserRepo:  user.NewPgxRepository(getPostgres(config.DB)),
	})

	service.SetServiceRouting(r, &service.Handlers{
		Logger:      logger,
		ServiceRepo: service.NewPgxRepository(getPostgres(config.DB)),
	})

	log.Fatal(fasthttp.ListenAndServe(config.Addr(), r.Handler))
}
