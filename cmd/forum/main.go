package main

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"log"
	"tech-db-forum/internal/app/forum"
	"tech-db-forum/internal/app/service"
	"tech-db-forum/internal/app/user"
	fasthttpprom "tech-db-forum/internal/pkg/Metrics"
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
	return ":" + c.Port
}

func getPostgres(config DBConfig) *pgx.ConnPool {
	dsn := fmt.Sprintf(`user=%s dbname=%s password=%s host=%s port=%s sslmode=disable`,
		config.Username, config.DBName, config.Password, config.Host, config.Port)

	conn, err := pgx.ParseConnectionString(dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: 1000,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Fatalf("Error %s occurred during connection to database", err)
	}

	return pool
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

	// metrics
	p := fasthttpprom.NewPrometheus("")
	p.Use(r)

	// connect db
	db := getPostgres(config.DB)

	user.SetUserRouting(r, &user.Handlers{
		Logger:   logger,
		UserRepo: user.NewPgxRepository(db),
	})

	forum.SetForumRouting(r, &forum.Handlers{
		Logger:    logger,
		ForumRepo: forum.NewPgxRepository(db),
		UserRepo:  user.NewPgxRepository(db),
	})

	service.SetServiceRouting(r, &service.Handlers{
		Logger:      logger,
		ServiceRepo: service.NewPgxRepository(db),
	})

	fmt.Printf("Start server on port %s\n", config.Port)

	log.Fatal(fasthttp.ListenAndServe(config.Addr(), p.Handler))
}
