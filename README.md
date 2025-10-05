# RSS Aggregator

RSS feed'lerini takip edip yönetmenizi sağlayan bir Go (Golang) REST API projesi.

## 🚀 Özellikler

- ✅ Kullanıcı oluşturma ve API key ile authentication
- ✅ RSS feed ekleme ve listeleme
- ✅ Feed takip sistemi (follow/unfollow)
- ✅ PostgreSQL database
- ✅ RESTful API yapısı

## 📁 Proje Yapısı

```
rss-aggregator/
├── cmd/
│   └── api/              # API server'ın main dosyası
│       └── main.go
├── internal/             # Private package'lar (sadece bu proje içinde kullanılır)
│   ├── auth/            # Authentication yardımcı fonksiyonları
│   ├── database/        # SQLC tarafından generate edilen database kodları
│   ├── handlers/        # HTTP route handler'ları
│   ├── middleware/      # HTTP middleware'ler (auth, logging, vb.)
│   └── models/          # API response modelleri ve dönüşüm fonksiyonları
├── sql/
│   ├── queries/         # SQLC için SQL query'leri
│   └── schema/          # Database migration dosyaları
├── .env                 # Environment variables (git'e eklenMEZ!)
├── .env.example         # Environment variables örneği
├── go.mod               # Go modül tanımı ve bağımlılıklar
└── README.md            # Bu dosya
```

## 🛠️ Teknolojiler

- **Go 1.25+** - Programlama dili
- **Chi Router** - HTTP router ve middleware
- **PostgreSQL** - Veritabanı
- **SQLC** - Type-safe SQL kod generator
- **Goose** - Database migration tool

## 📋 Gereksinimler

- Go 1.25 veya üstü
- PostgreSQL
- SQLC (opsiyonel, kod generate için)
- Goose (opsiyonel, migration için)

## ⚙️ Kurulum

### 1. Repository'yi klonlayın

```bash
git clone https://github.com/mehmettalhairmak/rss-aggregator-go.git
cd rss-aggregator-go
```

### 2. Bağımlılıkları yükleyin

```bash
go mod download
```

### 3. PostgreSQL database oluşturun

```bash
createdb rss_aggregator
```

### 4. Environment variables'ları ayarlayın

`.env.example` dosyasını `.env` olarak kopyalayın:

```bash
cp .env.example .env
```

`.env` dosyasını düzenleyip kendi değerlerinizi girin:

```
PORT=8080
DB_URL=postgres://username:password@localhost:5432/rss_aggregator?sslmode=disable
```

### 5. Database migration'ları çalıştırın

```bash
cd sql/schema
goose postgres "$DB_URL" up
```

### 6. Uygulamayı çalıştırın

```bash
go run cmd/api/main.go
```

Server `http://localhost:8080` adresinde başlayacak.

## 📚 API Endpoints

### Health Check

```
GET /v1/ready          # Server hazır mı kontrolü
GET /v1/error          # Test error endpoint
```

### Users

```
POST /v1/users         # Yeni kullanıcı oluştur
GET  /v1/users         # Kullanıcı bilgilerini getir (Auth gerekli)
```

**Örnek: Kullanıcı oluşturma**

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Mehmet"}'
```

### Feeds

```
POST /v1/feed          # Yeni feed ekle (Auth gerekli)
GET  /v1/feed          # Tüm feed'leri listele
```

**Örnek: Feed ekleme**

```bash
curl -X POST http://localhost:8080/v1/feed \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey YOUR_API_KEY" \
  -d '{"name": "Blog", "url": "https://example.com/feed.xml"}'
```

### Feed Follows

```
POST   /v1/feed_follows              # Feed takip et (Auth gerekli)
GET    /v1/feed_follows              # Takip edilen feed'leri listele (Auth gerekli)
DELETE /v1/feed_follows/{id}         # Feed takibini bırak (Auth gerekli)
```

## 🔐 Authentication

API'nin korunan endpoint'leri API Key ile authentication gerektirir.

API Key'inizi HTTP header'da gönderin:

```
Authorization: ApiKey YOUR_API_KEY_HERE
```

API Key, kullanıcı oluştururken response'da döner.

## 🗄️ Database Schema

### users

- id (UUID, Primary Key)
- created_at (Timestamp)
- updated_at (Timestamp)
- name (Text)
- api_key (Text, Unique)

### feeds

- id (UUID, Primary Key)
- created_at (Timestamp)
- updated_at (Timestamp)
- name (Text)
- url (Text, Unique)
- user_id (UUID, Foreign Key → users)

### feed_follows

- id (UUID, Primary Key)
- created_at (Timestamp)
- updated_at (Timestamp)
- user_id (UUID, Foreign Key → users)
- feed_id (UUID, Foreign Key → feeds)

## 🧪 Development

### SQLC ile kod generate etme

SQL query'lerini düzenledikten sonra:

```bash
sqlc generate
```

### Yeni migration oluşturma

```bash
cd sql/schema
goose create migration_name sql
```

## 📝 TODO

- [ ] RSS feed'leri otomatik olarak fetch etme (scraper/worker)
- [ ] Posts tablosu ve endpoint'leri
- [ ] Pagination desteği
- [ ] Rate limiting
- [ ] Logging middleware
- [ ] Unit testler
- [ ] Docker desteği

## 📄 License

MIT

## 👨‍💻 Geliştirici

Mehmet Talha Irmak - [@mehmettalhairmak](https://github.com/mehmettalhairmak)
