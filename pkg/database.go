package pkg

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	db     *sql.DB
	dbfile string = "weibo.db"
	mutex  sync.Mutex
)

func init() {
	_, err := os.Stat(dbfile)
	if os.IsNotExist(err) {
		if err := createDatabase(); err != nil {
			log.Fatal("创建数据库失败：", err)
		}
	}

	db, err = sql.Open("sqlite", dbfile)
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
	mutex.Lock()
	defer mutex.Unlock()
	var count int
	if err := db.QueryRow("SELECT COUNT(id) AS counts FROM weibo WHERE link = ?", url).Scan(&count); err != nil {
		log.Println(err)
		return false
	}
	return count > 0
}

func InsertDB(title, url string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	results, err := db.Exec("INSERT INTO weibo(summary, link) VALUES(?, ?)", title, url)
	if err != nil {
		log.Println("Insert Err: ", err)
		return false
	}
	result, _ := results.RowsAffected()
	return result > 0
}
