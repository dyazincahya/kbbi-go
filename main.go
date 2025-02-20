package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Struktur untuk menyimpan data
type KBBI struct {
	Word       string `json:"word"`
	Definition string `json:"arti"`
	Type       int    `json:"type"`
}

var db *sql.DB

func main() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal memuat file .env:", err)
	}

	// Ambil konfigurasi dari file .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Buat koneksi string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}
	defer db.Close()

	// Cek koneksi
	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak dapat dijangkau:", err)
	}

	fmt.Println("Berhasil terhubung ke database!")

	// Inisialisasi router
	router := gin.Default()

	// Endpoint untuk mengambil semua kata dengan limit
	router.GET("/", getAllWords)

	// Endpoint untuk mencari kata berdasarkan query parameter
	router.GET("/search", searchWord)

	// Jalankan server
	router.Run(":8080")
}

// Handler untuk mengambil semua kata dengan limit
func getAllWords(c *gin.Context) {
	limitStr := c.Query("limit")

	// Default limit
	limit := 100

	// Jika parameter limit diberikan, parsing ke integer
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'limit' harus berupa angka positif"})
			return
		}
		// Batasi limit maksimal menjadi 1000
		if parsedLimit > 1000 {
			parsedLimit = 1000
		}
		limit = parsedLimit
	}

	// Ambil data dari database dengan batas limit
	rows, err := db.Query("SELECT word, arti, type FROM api_kbbi_IV LIMIT ?", limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	// Buat slice untuk menyimpan hasil
	var words []KBBI
	for rows.Next() {
		var entry KBBI
		if err := rows.Scan(&entry.Word, &entry.Definition, &entry.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data"})
			return
		}
		words = append(words, entry)
	}

	// Berikan respon JSON
	c.JSON(http.StatusOK, words)
}

// Handler untuk mencari kata berdasarkan query parameter
func searchWord(c *gin.Context) {
	word := c.Query("word")

	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'word' harus diisi"})
		return
	}

	var entry KBBI
	err := db.QueryRow("SELECT word, arti, type FROM api_kbbi_IV WHERE word = ?", word).
		Scan(&entry.Word, &entry.Definition, &entry.Type)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kata tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		}
		return
	}

	// Berikan respon JSON
	c.JSON(http.StatusOK, entry)
}
