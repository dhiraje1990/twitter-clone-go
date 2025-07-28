package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Tweet struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Step 1: Ensure tweets table exists (initial table)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tweets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL
	);`)
	if err != nil {
		return nil, err
	}

	// Step 2: Auto-migrate: add missing columns
	rows, err := db.Query(`PRAGMA table_info(tweets);`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make(map[string]bool)

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt *string
		rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		columns[name] = true
	}

	if !columns["username"] {
		_, err = db.Exec(`ALTER TABLE tweets ADD COLUMN username TEXT DEFAULT 'anonymous';`)
		if err != nil {
			return nil, err
		}
	}

	if !columns["created_at"] {
		_, err = db.Exec(`ALTER TABLE tweets ADD COLUMN created_at DATETIME DEFAULT CURRENT_TIMESTAMP;`)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/tweets", func(c *gin.Context) {
		rows, _ := db.Query("SELECT id, username, text, created_at FROM tweets ORDER BY id DESC")
		var tweets []Tweet
		for rows.Next() {
			var t Tweet
			rows.Scan(&t.ID, &t.Username, &t.Text, &t.CreatedAt)
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
		_, err := db.Exec("INSERT INTO tweets (username, text) VALUES (?, ?)", t.Username, t.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed"})
			return
		}
		c.Status(http.StatusCreated)
	})

	r.DELETE("/tweets/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM tweets WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
			return
		}
		c.Status(http.StatusOK)
	})

	r.PUT("/tweets/:id", func(c *gin.Context) {
		id := c.Param("id")
		var t Tweet
		if err := c.BindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		_, err := db.Exec("UPDATE tweets SET text = ? WHERE id = ?", t.Text, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
			return
		}
		c.Status(http.StatusOK)
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
