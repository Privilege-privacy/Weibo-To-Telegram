package pkg

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

const dbfile string = "weibo.db"

var db *sql.DB

func init() {
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		if err := createDatabase(); err != nil {
			log.Fatal("创建数据库失败：", err)
		}
	}

	db, err = sql.Open("sqlite", dbfile+"?cache=shared&mode=rwc&_journal_mode=WAL")
	db.SetMaxOpenConns(1)
	if err != nil {
		log.Fatal("连接数据库失败：", err)
	}
}

func createDatabase() error {
	conn, err := sql.Open("sqlite", dbfile)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS weibo (id INTEGER PRIMARY KEY, summary TEXT, link TEXT)")
	if err != nil {
		return err
	}

	return nil
}

func ExistsInDB(url string) bool {
	row := db.QueryRow("SELECT COUNT(id) AS counts FROM weibo WHERE link = ?", url)

	var count int
	if err := row.Scan(&count); err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "ExistsInDB Failed", slog.String("URL", url))
		return false
	}

	return count > 0
}

func InsertDB(summary, link string) bool {
	_, err := db.Exec("INSERT INTO weibo(summary, link) VALUES(?, ?)", summary, link)
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelWarn, "InsertDB Failed", slog.String("URL", link))
		return false
	}
	return true
}
