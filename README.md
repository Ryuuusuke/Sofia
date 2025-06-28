# Sofia IRC Bot

**Sofia** adalah bot IRC ringan yang ditulis dalam bahasa Go. Bot ini mendukung:
- Autentikasi SASL
- Fetch judul dari link (termasuk YouTube)
- Loop RSS (modular)
- Kirim pesan langsung dari terminal (stdin)
- Konfigurasi via file `.ini`

---

## 🔧 Konfigurasi

Buat file bernama `config.ini` di direktori utama:

```ini
[sasl]
sasl = true
user = your-sasl-username
password = your-sasl-password

[irc]
server = irc.libera.chat:6697
nickname = sofiaaa
username = SofiaPertama
realname = Ratu Sofia
channel = `##sofia`
````

> Pastikan kamu menggunakan backtick (`) untuk nilai `channel`agar karakter`#\` tidak dianggap komentar.

---

## 🚀 Cara Menjalankan

```bash
go run main.go
```

Atau build dulu:

```bash
go build -o sofia .
./sofia
```

---

## 🖥️ Kirim Pesan dari Terminal

Ketik langsung di terminal tempat kamu menjalankan bot untuk mengirim pesan ke channel IRC yang sudah dikonfigurasi.

---

## 🌐 Fitur Link Preview

* Bila seseorang kirim tautan di channel, bot akan mencoba mengambil **judul halaman** secara otomatis.
* Link YouTube akan difetch menggunakan [YouTube oEmbed API](https://www.youtube.com/oembed).
* Untuk halaman biasa, bot menggunakan `chromedp` (headless Chrome via Go) untuk ambil judul `<title>`.

---

## 📰 RSS Feed

Kamu bisa menambahkan modul RSS kamu sendiri di folder `rss/`. Bot sudah modular dan mendukung fungsi loop RSS yang bisa dimodifikasi sesuai kebutuhan.

---

## 🧱 Struktur Direktori

```
.
├── main.go          # Entry point
├── config.ini       # Config file (user provided)
├── stdin/           # Modul pembaca dari stdin
├── rss/             # Modul RSS handler
└── go.mod           # Module file
```

---

## 📦 Dependencies

* [`go-ini/ini`](https://github.com/go-ini/ini) - Untuk parsing config `.ini`
* [`chromedp`](https://github.com/chromedp/chromedp) - Untuk mengambil title dari halaman web
* Standard Go libraries (`net`, `bufio`, `tls`, `regexp`, dll)

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah **BSD 3-Clause License**.  
Silakan lihat file [`LICENSE`](./LICENSE) untuk detail lengkapnya.


---
