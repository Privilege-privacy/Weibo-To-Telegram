package internal

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func init() {
	_, file := os.Stat("./weibo.db")
	if os.IsNotExist(file) {
		db, err := sql.Open("sqlite", "./weibo.db")
		if err != nil {
			log.Println("创建数据库失败", err)
		}
		db.Exec("CREATE TABLE IF NOT EXISTS \"weibo\" (\"id\" INTEGER PRIMARY KEY, \"summary\" TEXT, \"link\" TEXT);")
		defer db.Close()
	}
}

func Check(url string) int {
	db, err := sql.Open("sqlite", "./weibo.db")
	if err != nil {
		log.Println(err)
	}
	var counts int
	err = db.QueryRow("SELECT COUNT(id) AS counts FROM weibo WHERE link = ?", url).Scan(&counts)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	return counts
}

func Insert(title, link string) int {
	db, err := sql.Open("sqlite", "./weibo.db")
	if err != nil {
		log.Println(err)
	}
	results, err := db.Exec("INSERT INTO weibo(summary, link) VALUES(?, ?)", title, link)
	if err != nil {
		log.Printf("insert failed: %s\n", err)
	}
	result, _ := results.RowsAffected()
	defer db.Close()
	return int(result)
}
