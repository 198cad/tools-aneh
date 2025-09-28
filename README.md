# Tools CLI Aneh 🚀

Jadi gini guys, ini tuh CLI tools yang literally bisa ngebantu kalian manage PostgreSQL, RabbitMQ, sama MinIO. No cap, ini tools emang se-useful itu sih!

## ⚠️ DISCLAIMER PENTING BANGET ⚠️

**BACA DULU SEBELUM PAKE!**

Nih ya, gw kasih tau dari awal: **INI PROYEK PRIBADI GW**. Literally gw bikin buat kebutuhan gw sendiri, sesuai workflow gw, sesuai selera gw. Tools-tools ini bakal gw update kalo gw butuh, gw hapus kalo gw bosen, gw ubah seenak jidat kalo gw mau.

**Lo cocok? Gas pake aja.**  
**Lo ga cocok? Ya udah pergi aja, bikin sendiri sana.**

Gw ga nerima komplain, request fitur, atau drama apapun. This is my playground, deal with it. 💅

## Cara Install 

```bash
go build -o tools.exe .
```

Gampang kan? Tinggal build aja udah jadi deh 💯

## Konfigurasi

Nah ini dia yang seru, kalian bisa setup pake environment variables. Literally tinggal bikin file `.env` di folder yang sama:

```bash
cp .env.example .env
# Edit .env sesuai kebutuhan kalian ya bestie
```

### Environment Variables yang Bisa Dipake

#### PostgreSQL
- `PGHOST` - Host database kalian (default: localhost)
- `PGPORT` - Port-nya berapa (default: 5432)
- `PGUSER` - Username database (default: postgres)
- `PGPASSWORD` - Password-nya dong
- `PGDATABASE` - Database default (default: postgres)

#### RabbitMQ
- `RABBITMQ_HOST` - Host RabbitMQ (default: localhost)
- `RABBITMQ_MANAGEMENT_PORT` - Port buat management API (default: 15672)
- `RABBITMQ_DEFAULT_USER` - Username-nya (default: guest)
- `RABBITMQ_DEFAULT_PASS` - Password dong (default: guest)
- `RABBITMQ_DEFAULT_VHOST` - Virtual host gitu (default: /)

#### MinIO
- `MINIO_ENDPOINT` - Endpoint MinIO kalian (default: localhost:9000)
- `MINIO_ACCESS_KEY` atau `MINIO_ROOT_USER` - Access key nya
- `MINIO_SECRET_KEY` atau `MINIO_ROOT_PASSWORD` - Secret key jangan sampe bocor ya
- `MINIO_USE_SSL` - Pake SSL apa nggak (default: false)

## Cara Pake

### Command Database yang Kece

```bash
# Liat semua database yang ada
./tools.exe db list

# Bikin database baru
./tools.exe db create mydatabase

# Hapus database (hati-hati ya!)
./tools.exe db drop mydatabase

# Backup database biar aman
./tools.exe db backup mydatabase backup.sql

# Restore database kalo kenapa-kenapa
./tools.exe db restore mydatabase backup.sql


```

### Command RabbitMQ yang Mantap

```bash
# Liat semua queue
./tools.exe rabbit queues

# Bikin queue baru
./tools.exe rabbit create-queue myqueue

# Delete queue
./tools.exe rabbit delete-queue myqueue

# Bersihin message di queue
./tools.exe rabbit purge myqueue

# Liat exchanges yang ada
./tools.exe rabbit exchanges

# Bikin exchange baru
./tools.exe rabbit create-exchange myexchange --type topic

# Kirim message
./tools.exe rabbit publish exchange routing-key "isi pesan kalian"

# Liat statistik
./tools.exe rabbit stats
```

### Command MinIO yang Gokil

```bash
# Liat semua bucket
./tools.exe minio buckets

# Bikin bucket baru
./tools.exe minio create-bucket mybucket

# Hapus bucket
./tools.exe minio delete-bucket mybucket
./tools.exe minio delete-bucket mybucket --force  # Hapus sekalian sama isinya

# Liat isi bucket
./tools.exe minio list mybucket
./tools.exe minio list mybucket --recursive

# Upload file
./tools.exe minio upload mybucket localfile.txt
./tools.exe minio upload mybucket localfile.txt remote-name.txt

# Download file
./tools.exe minio download mybucket remote-file.txt
./tools.exe minio download mybucket remote-file.txt local-name.txt

# Hapus object
./tools.exe minio delete mybucket object-name.txt

# Copy antar bucket
./tools.exe minio copy source-bucket file.txt dest-bucket new-file.txt

# Cek info object/bucket
./tools.exe minio stat mybucket
./tools.exe minio stat mybucket object.txt

# Mirror folder lokal ke bucket
./tools.exe minio mirror ./local-dir mybucket
```

## Flag Command-Line Biar Makin Flexible

Semua command bisa di-override pake flag, jadi ga stuck sama env variables doang:

```bash
# Database pake koneksi custom
./tools.exe db list --host 192.168.1.10 --port 5433 --user admin --password secret

# RabbitMQ pake host lain
./tools.exe rabbit queues --host rabbitmq.example.com --port 15673

# MinIO pake endpoint sendiri
./tools.exe minio buckets --endpoint s3.example.com:9000 --access-key mykey --secret-key mysecret
```

## Buat yang Mau Develop

### Yang Harus Ada Dulu
- Go 1.21 atau yang lebih baru
- Akses ke PostgreSQL/MySQL (buat command database)
- Akses ke RabbitMQ Management API (buat command RabbitMQ)
- Akses ke MinIO server (buat command MinIO)

### Build dari Source Code

```bash
git clone git@github.com:198cad/tools-aneh.git
cd tools-aneh
go mod download
go build -o tools.exe .
```

### Jalanin Pake Air (Auto Reload Gitu)

```bash
air
```

## Update & Maintenance

FYI aja ya, koleksi tools ini bakal gw update sesuai kebutuhan pribadi gw. Jadi kalo tiba-tiba ada fitur baru, fitur ilang, atau behavior berubah, ya wajar aja. Gw ga punya roadmap, ga punya schedule release, ga punya planning apapun. 

**Update-nya kapan?** Kalo gw lagi butuh atau lagi gabut.  
**Fix bug-nya kapan?** Kalo bug-nya ganggu gw.  
**Tambah fitur baru?** Kalo gw perlu.

Simple as that. No drama needed 💁‍♂️

## Lisensi

**Apache License 2.0** 🔥

Kenapa Apache? Karena gw mau dong:
- **Lo boleh pake** buat project komersil, personal, apapun terserah
- **Lo boleh modif** sesuka hati lo, tapi inget kasih credit ya
- **Lo ga perlu share** source code modifikasi lo (beda sama GPL yang ribet)
- **Ada perlindungan patent** jadi lo aman dari drama patent infringement
- **Lo bisa sublicense** pake lisensi yang lebih ketat kalo mau

Yang **WAJIB** lo lakuin:
- Kasih credit ke gw (attribution)
- Sertain copy license Apache 2.0
- Kasih tau kalo ada perubahan major

Yang **GA BOLEH**:
- Pake trademark/nama gw sembarangan
- Salahin gw kalo ada bug atau masalah (no warranty bestie)
- Sue gw soal patent (kalo lo sue, hak lo dicabut otomatis)

Intinya: **Pake aja sesuka hati, tapi jangan lupa kasih credit dan jangan bawa-bawa gw kalo ada masalah.** 

Full license text ada di LICENSE file ya. Kalo males baca, ya udah skip aja, tapi tanggung sendiri resikonya 🤷‍♂️

---

**Copyright 2024 - tools-aneh project**  
*Licensed under the Apache License, Version 2.0*