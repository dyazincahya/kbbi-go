# KBBI GO
API KBBI IV menggunakan bahasa pemrograman Go Lang

## Spesifikasi
- GO v1.24.x
- MYSQL 5.7 atau lebih baru

## Basis data
Aplikasi ini menggunakan data KBBI dari [KBBI-SQL-database](https://github.com/dyazincahya/KBBI-SQL-database). Anda bisa mengimport data KBBI ke basis data di lokal komputer Anda.

## Instalasi
#### Env
Untuk instalasi bisa buat file `.env` yang dapat Anda salin kodenya dari file `.env.example`.
``` env
DB_USER=root
DB_PASS=
DB_HOST=localhost
DB_PORT=3306
DB_NAME=your_db_name
```
Sesuaikan konfigurasinya berdasarkan konfig yang ada di lokal komputer Anda.

#### Run
Untuk menjalankan aplikasi, Anda dapat menjalankan perintah
``` bash
go run main.go
```

#### Build
Untuk membangun aplikasi menjadi versi produksi, Anda dapat menjalankan perintah
``` bash
go build main.go
```

## Demo
Anda dapat mencoba demo API-nya disini : [https://services.x-labs.my.id/kbbi/](https://services.x-labs.my.id/kbbi/)

## Penulis
[Kang Cahya](https://www.kang-cahya.com)
