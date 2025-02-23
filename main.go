package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"math/rand"
	"time"

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

	// Endpoint root menampilkan informasi JSON statis
	router.GET("/", apiInfo)
	
	// Endpoint untuk mengambil semua kata dengan limit
	router.GET("/words", getAllWords)

	router.GET("/randomwords", getRandomWords)

	// Endpoint untuk mencari kata berdasarkan query parameter
	router.GET("/search", searchWord)

	// Jalankan server
	router.Run(":8080")
}

// Handler untuk menampilkan informasi API
func apiInfo(c *gin.Context) {
	info := gin.H{
		"api": gin.H{
			"name":        "API KBBI IV",
			"description": "API KBBI (Kamus Besar Bahasa Indonesia) versi IV ini digunakan untuk mencari arti kata dalam bahasa Indonesia.",
			"version":     "1.0.0",
			"endpoint": []gin.H{
				{
					"url":         "/",
					"description": "Menampilkan informasi tentang API KBBI IV.",
					"method":      "GET",
					"params":      []gin.H{},
				},
				{
					"url":         "/search/:kata",
					"description": "Mencari arti kata dalam bahasa Indonesia.",
					"method":      "GET",
					"params": []gin.H{
						{
							"name":        "kata",
							"description": "Kata yang ingin dicari artinya.",
							"type":        "string",
							"required":    true,
						},
					},
				},
				{
					"url":         "/words?limit=10",
					"description": "Menampilkan daftar kata yang tersedia dalam API KBBI IV.",
					"method":      "GET",
					"params": []gin.H{
						{
							"name":        "limit",
							"description": "Batas jumlah kata yang ingin ditampilkan.",
							"type":        "number",
							"required":    false,
						},
					},
				},
				{
					"url":         "/randomwords?limit=10",	
					"description": "Menampilkan daftar kata yang tersedia dalam API KBBI IV secara acak.",
					"method":      "GET",
					"params": []gin.H{
						{
							"name":        "limit",
							"description": "Batas jumlah kata yang ingin ditampilkan.",
							"type":        "number",
							"required":    false,
						},
					},
				},
				{
					"url":         "/search?word=kata",
					"description": "Mencari arti kata dalam bahasa Indonesia.",
					"method":      "GET",
					"params": []gin.H{
						{
							"name":        "word",
							"description": "Kata yang ingin dicari artinya.",
							"type":        "string",
							"required":    true,
						},
					},
				},
			},
		},
		"author": gin.H{
			"name":   "Kang Cahya",
			"blog":   "https://kang-cahya.com",
			"github": "https://github.com/dyazincahya",
		},
	}
	c.JSON(http.StatusOK, info)
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

func getRandomWords(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 1000 // Default limit jika tidak diatur oleh user

	// Parsing limit jika diberikan
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'limit' harus berupa angka positif"})
			return
		}
		if parsedLimit > 1000 {
			parsedLimit = 1000 // Batasi limit maksimal menjadi 1000
		}
		limit = parsedLimit
	}

	// Ambil semua kata dari database
	rows, err := db.Query("SELECT word, arti, type FROM api_kbbi_IV")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	defer rows.Close()

	// Simpan semua kata dalam slice
	var words []KBBI
	for rows.Next() {
		var entry KBBI
		if err := rows.Scan(&entry.Word, &entry.Definition, &entry.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data"})
			return
		}
		words = append(words, entry)
	}

	// Acak daftar kata
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })

	// Batasi jumlah kata sesuai limit
	if len(words) > limit {
		words = words[:limit]
	}

	// Berikan response JSON
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
	err := db.QueryRow("SELECT word, arti, type FROM api_kbbi_IV WHERE word LIKE ?", "%"+word+"%").
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
