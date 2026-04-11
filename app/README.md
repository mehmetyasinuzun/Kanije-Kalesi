# 🏰 Kanije Kalesi — Kurulum ve Kullanım Rehberi

> Bilgisayarında ne olduğunu her zaman bil. Telegram'dan öğren.

---

## 📋 Ne Yapar?

Kanije Kalesi, bilgisayarında belirlediğin olayları **anlık olarak Telegram'a bildirir**:

| Olay | Bildirim |
|------|---------|
| ✅ Başarılı giriş | Kimin, ne zaman oturum açtığı |
| 🚨 Yanlış şifre | Kamera fotoğrafı + IP adresi |
| 🔒 Ekran kilitlenme/açılma | Saat bilgisiyle |
| 🖥️ Bilgisayar açıldı/kapandı | Çalışma süresiyle |
| 😴 Uyku/Hibernate | Saat damgasıyla |
| ☀️ Uykudan uyanma | Otomatik mi, manuel mi |
| 🔌 USB takıldı/çıkarıldı | Sürücü adı ve boyutu |
| 📡 İnternet kesildi/bağlandı | Hangi ağdan, ne zaman |
| 🔀 WiFi/Kablo değişimi | Eski → Yeni ağ adı |

---

## ⚡ Hızlı Kurulum (5 Adım)

### 1. Python Kur

[python.org/downloads](https://python.org/downloads) → Python 3.11+ indir ve kur.
> ✅ Kurulum sırasında **"Add Python to PATH"** kutusunu işaretle!

### 2. Bağımlılıkları Yükle

Admin PowerShell'de:
```powershell
cd c:\...\Kanije-Kalesi\app
pip install -r requirements.txt
```

### 3. Telegram Bot Oluştur

1. Telegram'da **@BotFather**'a yaz → `/newbot`
2. Bot adı ver → Token al (örnek: `8750193308:AAEt8Q...`)
3. **@userinfobot**'a `/start` yaz → Chat ID'ni öğren (örnek: `987654321`)

### 4. Kimlik Bilgilerini Gir

İki seçenek:

**A) Otomatik (önerilen):**
```powershell
python kanije.py start
```
Token yoksa **kurulum penceresi otomatik açılır** → Token ve Chat ID'yi gir → Kaydet.

**B) Elle:**
`app/.env` dosyası oluştur:
```
KANIJE_BOT_TOKEN=8750193308:AAEt8Qg...
KANIJE_CHAT_ID=987654321
```

### 5. Başlat

```powershell
# Yönetici PowerShell'de (sağ tık → Yönetici olarak çalıştır):
cd c:\...\Kanije-Kalesi\app
python kanije.py start
```

---

## 🖥️ System Tray (Sistem Tepsisi) Menüsü

Uygulama çalışırken sağ alt köşede **🏰 ikonu** görünür.

Sağ tıkla:

| Menü | Açıklama |
|------|---------|
| ⚙️ Telegram Ayarları | Token ve Chat ID'yi değiştir |
| 📋 Log Dosyasını Aç | `kanije.log` → Notepad'de açılır |
| ❌ Çıkış | Uygulamayı kapat ve Telegram'a bildir |

---

## 📱 Telegram Komutları

Botuna şu komutları yazabilirsin:

| Komut | Ne Yapar |
|-------|---------|
| `/status` | CPU, RAM, Disk, uptime |
| `/events` | Son 10 olay |
| `/config` | Hangi bildirimler açık/kapalı |
| `/photo` | Anlık kamera fotoğrafı |
| `/screenshot` | Anlık ekran görüntüsü |
| `/lock` | Ekranı kilitle |
| `/ping` | Bota ulaşılıyor mu? |
| `/restart` | Bilgisayarı yeniden başlat (15sn onay) |
| `/shutdown` | Bilgisayarı kapat (15sn onay) |
| `/cancel` | Bekleyen işlemi iptal et |
| `/help` | Tüm komutlar |

---

## 🔐 Güvenlik

- `.env` dosyası **git'e yüklenmez** (`.gitignore` ile korumalı)
- Kimseyle paylaşma — tokenın oturumlarını `@BotFather → /mybots → Revoke` ile yenileyebilirsin
- Sadece senin Chat ID'nden gelen komutlar çalışır
- Gönderilen fotoğraflar yerel olarak silinir (varsayılan)

---

## 🔄 Otomatik Başlatma (Bilgisayar açılınca)

Admin PowerShell'de:
```powershell
.\install.ps1
```
Bu komut, Windows Görev Zamanlayıcısı'na oturum açıldığında otomatik başlatma ekler.

---

## ❓ Sorun Giderme

| Sorun | Çözüm |
|-------|-------|
| Bildirim gelmiyor | Yönetici PowerShell ile çalıştır |
| Login olayları yok | Yönetici yetkisi gerekli |
| `pip install` hatası | `pip install -r requirements.txt --user` dene |
| Python bulunamadı | PATH'e Python ekle veya tam yolu kullan |
| Token hatalı | `python kanije.py test` ile test et |
