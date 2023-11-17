package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"

	"42tokyo-road-to-dojo-go/pkg/repositories"
	"42tokyo-road-to-dojo-go/pkg/server"
)

var (
	// Listenするアドレス+ポート
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "tcp host:port to connect")
	flag.Parse()
	log.Default().SetFlags(log.LstdFlags | log.Llongfile)
}

// connectDB DBに接続し、Pingで接続を確認する。
// いずれかが失敗したら、ログを出してプログラムを終了する。
func connectDB() *sql.DB {

	dsn := "root:ca-tech-dojo@tcp(localhost:3306)/CA_Tech_Dojo"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping to the database: %v", err)
	}
	return db
}

func newRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return rdb
}

func cacheMasterData(repos *repositories.Repositories) {
	if err := repos.ItemRepository.CacheItems(); err != nil {
		log.Fatalf("Failed to cache items: %v", err)
	}
	if err := repos.GameSettingsRepository.CacheActiveGameSettings(); err != nil {
		log.Fatalf("Failed to cache game settings: %v", err)
	}
}

func main() {
	db := connectDB()
	defer db.Close()

	rdb := newRedisClient()
	defer rdb.Close()

	repos := repositories.NewRepositories(db, rdb)
	cacheMasterData(repos)
	server.Serve(addr, repos)
}
