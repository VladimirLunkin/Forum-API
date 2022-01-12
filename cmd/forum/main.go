package main

import (
	"database/sql"
	"github.com/fasthttp/router"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"log"
	"tech-db-forum/internal/app/forum"
	"tech-db-forum/internal/app/user"
)

func getPostgres() *sql.DB {
	dsn := "user=postgres dbname=forum1 password=251099 host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(1000)
	return db
}

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	r := router.New()

	user.SetUserRouting(r, &user.Handlers{
		Logger:   logger,
		UserRepo: user.NewPgxRepository(getPostgres()),
	})

	forum.SetForumRouting(r, &forum.Handlers{
		Logger:    logger,
		ForumRepo: forum.NewPgxRepository(getPostgres()),
		UserRepo: user.NewPgxRepository(getPostgres()),
	})

	log.Fatal(fasthttp.ListenAndServe(":5000", r.Handler))
}
