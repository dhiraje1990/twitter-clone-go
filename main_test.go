package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var testRouter *gin.Engine

func TestMain(m *testing.M) {
	// Setup
	os.Remove("./test.db")
	db, err := InitDB("./test.db")
	if err != nil {
		panic(err)
	}
	testRouter = SetupRouter(db)

	// Run tests
	code := m.Run()

	// Teardown
	db.Close()
	os.Remove("./test.db")

	os.Exit(code)
}

func TestPostTweet(t *testing.T) {
	body := []byte(`{"text":"Test tweet from automated test"}`)
	req := httptest.NewRequest("POST", "/tweets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d", w.Code)
	}
}

func TestGetTweets(t *testing.T) {
	// Insert a tweet
	req := httptest.NewRequest("POST", "/tweets", bytes.NewBuffer([]byte(`{"text":"Another test"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	// Now test GET
	req = httptest.NewRequest("GET", "/tweets", nil)
	w = httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	var tweets []Tweet
	json.Unmarshal(body, &tweets)

	if len(tweets) == 0 {
		t.Errorf("Expected at least 1 tweet, got 0")
	}
}
