package pkg

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

var db *sql.DB

func init() {
	_, err := os.Stat("./weibo.db")
	if os.IsNotExist(err) {
		if err := createDatabase(); err != nil {
			log.Fatal("创建数据库失败：", err)
		}
	}

	db, err = sql.Open("sqlite", "./weibo.db")
	if err != nil {
		log.Fatal("连接数据库失败：", err)
	}
}

func createDatabase() error {
	conn, err := sql.Open("sqlite", "./weibo.db")
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

func Check(url string) (result int) {
	if err := db.QueryRow("SELECT COUNT(id) AS counts FROM weibo WHERE link = ?", url).Scan(&result); err != nil {
		log.Println(err)
	}
	return result
}

func Insert(title, url string) int {
	results, err := db.Exec("INSERT INTO weibo(summary, link) VALUES(?, ?)", title, url)
	if err != nil {
		log.Println("Insert Err: ", err)
	}
	result, _ := results.RowsAffected()
	return int(result)
}
