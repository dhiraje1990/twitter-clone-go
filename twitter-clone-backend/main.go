package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Tweet struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tweets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL
	);`)
	return db, err
}

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/tweets", func(c *gin.Context) {
		rows, _ := db.Query("SELECT id, text FROM tweets ORDER BY id DESC")
		var tweets []Tweet
		for rows.Next() {
			var t Tweet
			rows.Scan(&t.ID, &t.Text)
			tweets = append(tweets, t)
		}
		c.JSON(http.StatusOK, tweets)
	})

	r.POST("/tweets", func(c *gin.Context) {
		var t Tweet
		if err := c.BindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		_, err := db.Exec("INSERT INTO tweets (text) VALUES (?)", t.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed"})
			return
		}
		c.Status(http.StatusCreated)
	})

	return r
}

func main() {
	db, err := InitDB("./tweets.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := SetupRouter(db)
	r.Run(":8080")
}
