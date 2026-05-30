# Katkı Rehberi

Kanije Kalesi'ne katkıda bulunmak isteyenler için. Aktif geliştirme
[`go/`](go/) dizinindeki Go uygulamasında yapılır. (`app/` Python sürümü
[kullanımdan kalkmıştır](app/DEPRECATED.md).)

## Geliştirme Ortamı

- **Go 1.21+**
- **ffmpeg** (kamera/cihaz testleri için, opsiyonel)
- Linux'ta `-race` testleri için bir C derleyicisi (gcc/clang)

```bash
cd go
go mod download
```

## Çalıştırma

```bash
make run                # KANIJE_LOG_LEVEL=debug ile başlatır
# veya
go run ./cmd/kanije/ version
```

## Kalite Kapıları (CI bunları zorunlu kılar)

PR açmadan önce yereldeki tüm kapıların geçtiğinden emin olun:

```bash
cd go
gofmt -l .              # çıktı BOŞ olmalı
go vet ./...
go build ./...                                   # mevcut platform
GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build ./...
GOOS=linux  GOARCH=arm64 CGO_ENABLED=0 go build ./...   # Raspberry Pi 4/5
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ./...
go test ./...
go test -race ./...     # CGo (C derleyici) gerektirir
```

> **Yerelde `-race`:** Race dedektörü CGo ister. Windows'ta bir C derleyici yoksa
> kuramazsınız — en kolayı [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) veya MSYS2.
> Kurulduktan sonra `CGO_ENABLED=1 go test -race ./...` çalışır. (Yoksa CI'ın Linux
> job'ı zaten `-race` koşar.)

Linter:

```bash
golangci-lint run        # yapılandırma: go/.golangci.yml
# veya yalnızca staticcheck:
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

> **Önemli:** Bu uygulama **CGo'suz** derlenir (saf-Go SQLite ve ffmpeg
> alt-süreci sayesinde). Cross-compile'ı bozan CGo bağımlılıkları eklemeyin.

## Kod Stili

- **Biçim:** `gofmt`/`goimports` zorunlu.
- **İsimlendirme:** Go konvansiyonları (initializmler büyük: `XML`, `GUID`, `ID`).
- **Yorumlar ve tanımlayıcılar:** İngilizce.
- **Kullanıcıya görünen metinler** (loglar, Telegram mesajları, CLI çıktısı): Türkçe,
  tam diakritikle (`ı ş ğ ç ö ü İ`).
- **Eşzamanlılık:** Paylaşılan durum korunmalı. `go test -race` temiz olmalı.
  Yeni paylaşılan alanlar için mutex/atomic kullanın; `config` paketinde alanları
  doğrudan okumak yerine kilitli getter'ları tercih edin.
- **Hata yönetimi:** Hataları yutmayın; ya işleyin ya da en azından loglayın.

## Mimari (özet)

```
cmd/kanije        → CLI giriş noktası, platforma özel listener enjeksiyonu
internal/app      → orchestrator (errgroup ile yaşam döngüsü, graceful shutdown)
internal/event    → merkezi event bus (rate limit + dedup)
internal/listener → platforma özel olay kaynakları (windows/, linux/, stub/)
internal/notifier → Telegram istemcisi, bot, formatter, kurulum sihirbazı
internal/storage  → SQLite kalıcılık (saf-Go sürücü)
internal/capture  → kamera (ffmpeg) + ekran görüntüsü
internal/network  → bağlantı/SSID izleme
internal/config   → TOML + ortam değişkeni yapılandırması, çalışma zamanı mutasyonu
```

Yeni bir olay kaynağı eklemek için: `listener.Listener` arayüzünü uygulayın ve
ilgili platformun `cmd/kanije/listeners_*.go` dosyasında kaydedin.

## Commit / PR

- Küçük, odaklı commit'ler.
- Açıklayıcı başlık (Türkçe veya İngilizce).
- PR açıklamasında: ne değişti, neden, nasıl test edildi.
