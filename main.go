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

func main() {
	db, err := sql.Open("sqlite3", "./tweets.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tweets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL
	);`)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	// CORS (so React frontend can talk to Go backend)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// List tweets
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

	// Post a tweet
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

	r.Run(":8080")
}
