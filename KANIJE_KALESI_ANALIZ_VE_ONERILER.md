# 🔍 Kanije Kalesi — Kapsamlı Analiz & İyileştirme Önerileri

> 9 dosyanın tümü incelendi. Mevcut güçlü yönler, eksikler, iyileştirme önerileri ve yeni dosya fikirleri aşağıda.

---

## 📊 MEVCUT DOSYA DURUM DEĞERLENDİRMESİ

| Dosya | Kalite | Tamamlanma | Kritik Eksik |
|-------|--------|------------|--------------|
| `WINDOWS11_HARDENING_KALE.md` | ★★★★★ | %95 | Küçük detaylar |
| `WINDOWS10_HARDENING_KALE.md` | ★★★★★ | %95 | Küçük detaylar |
| `LINUX_HARDENING_KALE.md` | ★★★★☆ | %90 | Fedora/Arch desteği zayıf |
| `DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md` | ★★★★★ | %95 | Neredeyse eksiksiz |
| `SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md` | ★★★★★ | %98 | Çok kapsamlı |
| `HARICI_USB_SSD_BOOT_REHBERI.md` | ★★★★★ | %95 | Neredeyse eksiksiz |
| `VENTOY_WTG_MULTIBOOT_REHBERI.md` | ★★★★☆ | %90 | Troubleshooting genişletilebilir |
| `GHOSTGUARD_MIMARI_VE_PLAN.md` | ★★★★★ | %95 | Plan hazır, kod yok |
| `README.md` | ★★★☆☆ | %70 | Contributing, lisans, badge eksik |

---

## 🔧 BÖLÜM 1 — MEVCUT DOSYALARA EKLENMESİ GEREKENLER

### 1.1 — `WINDOWS11_HARDENING_KALE.md` İyileştirmeler

| # | Konu | Detay | Önem |
|---|------|-------|------|
| 1 | **Windows Sandbox** | Win 11 Pro'da gelen Sandbox özelliğini güvenli uygulama çalıştırma için ekle. Güvenilmeyen dosyaları izole ortamda açma. | ★★★★☆ |
| 2 | **Microsoft Defender Exclusion Attack** | Saldırganların Defender istisna listesine zararlı yollar ekleyebildiği senaryoyu anlat ve koruma yöntemi ekle | ★★★★☆ |
| 3 | **Recall Özelliği (Win 11 24H2)** | AI tabanlı ekran kaydı özelliğinin güvenlik riskleri ve nasıl kapatılacağı | ★★★★★ |
| 4 | **Otomatik Cihaz Şifreleme vs BitLocker** | Win 11 Home'da varsayılan "Cihaz Şifreleme" ile Pro BitLocker farkı | ★★★☆☆ |
| 5 | **Hardening Doğrulama Script'i** | Windows 11 için tüm katmanları tek seferde kontrol eden `kale_durum_win11.ps1` scripti eksik (Win 10 belgesinde var ama ayrı dosya değil) | ★★★★★ |

### 1.2 — `WINDOWS10_HARDENING_KALE.md` İyileştirmeler

| # | Konu | Detay | Önem |
|---|------|-------|------|
| 1 | **EOL (End of Life) Uyarısı** | Windows 10, Ekim 2025'te destek sona eriyor. Güvenlik yamaları bitmeden önce geçiş planı veya ESU (Extended Security Updates) satın alma uyarısı | ★★★★★ |
| 2 | **Windows 10 → 11 Karşılaştırma Tablosu** | "Bu rehberdeki hangi özellik Win 11'de farklı?" diye net bir karşılaştırma tablosu | ★★★★☆ |
| 3 | **AppLocker** | WDAC'a ek olarak AppLocker kuralları — bazı Enterprise ortamlarında daha kolay yönetilir | ★★★☆☆ |

### 1.3 — `LINUX_HARDENING_KALE.md` İyileştirmeler

| # | Konu | Detay | Önem |
|---|------|-------|------|
| 1 | **Fedora/RHEL Desteği** | SELinux yapılandırması Fedora özelinde eksik. `semanage`, `restorecon`, `audit2allow` komutları eklenmeli | ★★★★☆ |
| 2 | **ClamAV/chkrootkit** | Antivirüs ve rootkit tarama araçları kurulumu ve otomatik tarama cron job'ları | ★★★★☆ |
| 3 | **Unattended Upgrades** | Ubuntu/Debian için otomatik güvenlik güncellemeleri yapılandırması eksik | ★★★★★ |
| 4 | **USB Guard** | `usbguard` ile USB cihaz whitelist/blacklist politikası — fiziksel saldırıya karşı çok güçlü | ★★★★☆ |
| 5 | **Lynis Güvenlik Taraması** | `lynis audit system` ile kapsamlı güvenlik denetimi ve skor yükseltme önerileri | ★★★★★ |
| 6 | **Docker/Container Güvenliği** | Kali özellikle container ile kullanılıyor — Docker daemon güvenliği, rootless containers | ★★★☆☆ |

### 1.4 — `GHOSTGUARD_MIMARI_VE_PLAN.md` İyileştirmeler

| # | Konu | Detay | Önem |
|---|------|-------|------|
| 1 | **RDP/SSH Uzak Erişim Algılama** | Event ID 4624 Logon Type 10 (RDP) ve 3 (Network) ayrımı — uzak giriş uyarısı | ★★★★★ |
| 2 | **Disk Full Koruması** | Log dosyaları diski doldurursa ne olacak — disk usage kontrolü ve eski logları otomatik silme | ★★★★☆ |
| 3 | **Self-Update Mekanizması** | GitHub'dan versiyon kontrolü ve otomatik güncelleme bildirim | ★★★☆☆ |
| 4 | **Telegram Botu Komut Desteği** | `/status`, `/screenshot`, `/photo` komutları ile uzaktan bilgi alma — iki yönlü iletişim | ★★★★★ |
| 5 | **WiFi Ağ Değişikliği Algılama** | Bilinen ağdan bilinmeyen ağa geçişte uyarı — MitM koruması | ★★★★☆ |
| 6 | **Çoklu Bildirim Kanalı** | Telegram yanında e-posta (SMTP), Discord webhook, Pushover desteği | ★★★☆☆ |

### 1.5 — `VENTOY_WTG_MULTIBOOT_REHBERI.md` İyileştirmeler

| # | Konu | Detay | Önem |
|---|------|-------|------|
| 1 | **Ventoy Hash Doğrulama** | ISO dosyalarının SHA-256 hash kontrolü — indirilen ISO'nun bozuk/değiştirilmiş olmadığından emin olma | ★★★★★ |
| 2 | **Auto-install Desteği** | Ventoy'un `autoinstall` özelliği ile Ubuntu/Windows unattended kurulum | ★★★☆☆ |
| 3 | **VMware/VirtualBox'ta VHD Test Etme** | VHD'yi fiziksel boot'tan önce VM'de test etme rehberi | ★★★★☆ |

### 1.6 — `README.md` İyileştirmeler

| # | Konu | Detay | Önem |
|---|------|-------|------|
| 1 | **Badges** | Rehber sayısı, son güncelleme tarihi, lisans badge'leri | ★★★☆☆ |
| 2 | **Contributing** | Katkıda bulunma kuralları | ★★☆☆☆ |
| 3 | **Yıldız Geçmişi** | İlham ve kaynak linkleri | ★★★☆☆ |
| 4 | **Hızlı Başlangıç Rehberi** | "Nereden başlamalıyım?" sorusuna cevap veren bir akış şeması | ★★★★★ |

---

## 🆕 BÖLÜM 2 — YENİ DOSYA ÖNERİLERİ

### 2.1 — 📋 `AG_GUVENLIGI_VE_VPN.md` — Ağ Güvenliği Rehberi (ÖNCELİKLİ)

Mevcut rehberlerde ağ güvenliği her dosyada dağınık. Özel bir ağ güvenliği belgesi:

```
İçerik:
├── WiFi Güvenliği (WPA3, Evil Twin koruması, MAC filtreleme gerçekliği)
├── VPN Seçimi ve Yapılandırması (WireGuard, OpenVPN, IPSec)
├── Kendi VPN Sunucun (VPS + WireGuard kurulumu)
├── DNS Sızıntı Testi (dnsleaktest.com, WebRTC leak)
├── Tor Ağı ve Tails OS Kullanımı
├── Man-in-the-Middle (MitM) Koruması
├── Port Tarama ve Tespit (Nmap ile kendi ağını tara)
├── Router Sertleştirme (admin şifresi, UPnP kapatma, firmware güncelleme)
└── Ağ İzleme Araçları (Wireshark, tcpdump, ntopng)
```

**Neden Gerekli:** Tüm OS sertleştirmesi yapılsa bile ağ katmanı açıksa saldırgan trafiği okuyabilir.

---

### 2.2 — 🔑 `SIFRE_YONETIMI_VE_2FA.md` — Parola ve Kimlik Güvenliği

```
İçerik:
├── Şifre Yöneticisi Karşılaştırma (Bitwarden, KeePassXC, 1Password)
├── KeePassXC Tam Kurulum ve Kullanım (cross-platform)
├── TOTP (2FA) Nasıl Çalışır (algoritma düzeyinde açıklama)
├── YubiKey / FIDO2 Donanım Anahtarı Kurulumu
│   ├── Windows Hello entegrasyonu
│   ├── Linux PAM entegrasyonu
│   ├── SSH key olarak kullanma
│   └── Google/GitHub/Microsoft hesap koruması
├── Parolasız (Passkey) Gelecek — FIDO2 WebAuthn standardı
├── SIM Swap Saldırısı ve Korunma
├── Sosyal Mühendislik Savunması
└── Güçlü Parola Oluşturma Stratejileri (diceware, xkcd yöntemi)
```

**Neden Gerekli:** Tüm teknik önlemler, "123456" şifresi kullanılıyorsa boşa gider.

---

### 2.3 — 🛡️ `YEDEKLEME_VE_FELAKET_KURTARMA.md` — Backup & Disaster Recovery

```
İçerik:
├── 3-2-1 Yedekleme Kuralı (3 kopya, 2 farklı medya, 1 offsite)
├── Windows Yedekleme
│   ├── wbAdmin ile sistem görüntüsü
│   ├── Dosya Geçmişi (File History) yapılandırması
│   ├── VSS (Volume Shadow Copy) ayarları
│   └── Bare-metal recovery
├── Linux Yedekleme
│   ├── rsync + cron otomatik yedekleme
│   ├── Timeshift sistem snapshot'ı
│   ├── borgbackup (şifreli, deduplicated backup)
│   └── LUKS header yedeklemesi (zaten kısmen var ama derinek)
├── Bulut Yedekleme Güvenliği
│   ├── Client-side şifreleme (rclone + crypt)
│   ├── Cryptomator ile bulut şifreleme
│   └── Zero-knowledge yedekleme servisleri
├── Ransomware Koruması
│   ├── Offline yedeklerin önemi
│   ├── Immutable backup stratejisi
│   └── Air-gapped yedekleme
└── Felaket Kurtarma Senaryoları
    ├── Disk tamamen bozuldu → Ne yapılır?
    ├── Ransomware → Adım adım kurtarma
    └── BitLocker kurtarma anahtarı kayboldu → Alternatifler
```

**Neden Gerekli:** En sık gözden kaçan alandır. Disk şifreli ama yedek yoksa — donanım arızası = veri kaybı.

---

### 2.4 — 🔬 `FORENSIC_VE_INCIDENT_RESPONSE.md` — Olay Müdahalesi

```
İçerik:
├── Olay Müdahale Süreci (NIST 800-61 çerçevesi)
│   ├── Hazırlık
│   ├── Tespit ve Analiz
│   ├── Sınırlandırma
│   ├── Temizleme ve Kurtarma
│   └── Sonrası Analiz (Post-Incident)
├── Canlı Sistem Analizi
│   ├── Çalışan süreçleri analiz etme (procmon, Process Explorer)
│   ├── Ağ bağlantılarını analiz etme (netstat, TCPView)
│   ├── Zamanlanmış görevleri inceleme
│   └── Registry persistence noktaları
├── Bellek Analizi (RAM Forensics)
│   ├── Volatility Framework kullanımı
│   ├── RAM dump alma (WinPmem, DumpIt)
│   └── Process injection tespit etme
├── Disk Forensics
│   ├── Autopsy / FTK Imager
│   ├── Silinmiş dosya kurtarma
│   └── Zaman çizelgesi analizi
├── Log Analizi
│   ├── Windows Event Log analizi (Event ID korelasyonu)
│   ├── Linux auth.log / journal analizi
│   └── Sysmon log korelasyonu
└── IOC (Indicators of Compromise) Tarama
    ├── YARA kuralları
    ├── Sigma kuralları
    └── VirusTotal API kullanımı
```

**Neden Gerekli:** Sertleştirme "önleme"dir. Ama saldırı gerçekleştiğinde ne yapılacağını bilmek en az sertleştirme kadar önemli.

---

### 2.5 — 📱 `MOBIL_GUVENLIK.md` — Telefon ve Tablet Güvenliği

```
İçerik:
├── Android Güvenlik Sertleştirme
│   ├── Telefon şifreleme (FDE vs FBE)
│   ├── Geliştirici seçenekleri güvenliği
│   ├── USB debugging kapatma
│   ├── Uygulama izinleri yönetimi
│   ├── DNS over HTTPS (Private DNS)
│   └── GrapheneOS / CalyxOS alternatifleri
├── iOS Güvenlik
│   ├── Lockdown Mode
│   ├── Advanced Data Protection
│   └── iCloud Private Relay
├── Mobil Tehditler
│   ├── Pegasus/özellikli spyware
│   ├── Evil Twin WiFi saldırısı
│   ├── ADB saldırıları
│   └── Sahte uygulama marketleri
└── İki Cihaz Arasında Güvenli Veri Aktarımı
```

---

### 2.6 — 🌐 `TARAYICI_VE_WEB_GUVENLIGI.md` — Web Tarama Güvenliği

```
İçerik:
├── Tarayıcı Sertleştirme
│   ├── Firefox — about:config güvenlik ayarları
│   ├── Chromium/Brave güvenlik ayarları
│   ├── Eklenti güvenliği (uBlock Origin, Privacy Badger, NoScript)
│   └── JavaScript fingerprinting koruması
├── HTTPS Everywhere ve HSTS
├── Certificate Pinning ve CA güvenliği
├── Phishing Algılama ve Korunma
├── Cookie ve Tracker Yönetimi
├── Sandboxed Tarayıcı Kullanımı (Firejail ile Firefox)
└── Tor Browser güvenli kullanımı
```

---

### 2.7 — ⚡ `HIZLI_SERTLESTIRME_SCRIPTLERI.md` — Otomatik Sertleştirme

```
İçerik:
├── Windows Sertleştirme Scripti (tek PowerShell scripti — tüm katmanlar)
│   ├── kale_harden_win11.ps1
│   ├── kale_harden_win10.ps1
│   └── kale_durum_raporu.ps1 (mevcut durumu kontrol et)
├── Linux Sertleştirme Scripti
│   ├── kale_harden_ubuntu.sh
│   ├── kale_harden_kali.sh
│   └── kale_durum_linux.sh (Lynis entegreli)
├── Parametre Açıklamaları
│   └── Her script satırı NEDEN orada — acemi dostu
└── Geri Alma (Rollback) Scriptleri
    ├── kale_rollback_win.ps1
    └── kale_rollback_linux.sh
```

**Neden Gerekli:** 700+ satırlık rehberleri elle uygulamak zor. Tek script ile tüm sertleştirme uygulanabilmeli.

---

### 2.8 — 🔐 `SIFRELEME_DERINKOKU.md` — Kriptografi ve Şifreleme Temelleri

```
İçerik:
├── Simetrik vs Asimetrik Şifreleme (AES, RSA, ECC)
├── Hash Fonksiyonları (SHA-256, SHA-512, Argon2)
├── Key Derivation (PBKDF2, scrypt, Argon2id)
├── BitLocker İç Yapısı (FVEK, VMK, SRK zinciri detaylı)
├── LUKS2 İç Yapısı (header, keyslot, digest)
├── TLS/SSL Nasıl Çalışır (handshake adım adım)
├── GPG ile E-posta/Dosya Şifreleme
│   ├── GPG anahtar oluşturma
│   ├── Dosya şifreleme/imzalama
│   └── Web of Trust kavramı
├── SSH Anahtar Çeşitleri (RSA vs Ed25519 vs ECDSA)
└── Quantum Computing ve Şifreleme Geleceği (Post-Quantum)
```

---

## 📋 BÖLÜM 3 — GENEL İYİLEŞTİRME ÖNERİLERİ

### 3.1 — Tüm Dosyalara Uygulanması Gereken Değişiklikler

| # | Değişiklik | Açıklama |
|---|-----------|----------|
| 1 | **Versiyon/Changelog** | Her dosyaya versiyon numarası ve son değişiklik özeti ekle (`v1.2 — USB Guard bölümü eklendi`) |
| 2 | **Zorluk Seviyesi Etiketi** | Her bölüme `👶 Başlangıç / 🧑‍💻 Orta / 🔬 İleri` etiketi koy — okuyucu hangi seviyede anlasın |
| 3 | **"Bunu Neden Yapıyorum?" Kutusu** | Her teknik adımın hemen altına `> ❓ Bu adımı atlarsam ne olur?` şeklinde kısa risk açıklaması |
| 4 | **Cross-reference Linkleri** | Belgeler arası çapraz referans: "Bu konunun detayı için → [DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md](./DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md#bölüm-2)" |
| 5 | **Doğrulama Komutları Vurgusu** | Her adımın sonundaki "doğrulama" komutlarını `✅ DOĞRULAMA:` ile başlat — göz taramasında kolay bulsun |

### 3.2 — Repo Yapısı İyileştirmesi

```
Kanije-Kalesi/
├── README.md
├── CHANGELOG.md                    ← YENİ: değişiklik geçmişi
├── LICENSE                         ← YENİ: MIT veya CC BY-SA
│
├── os-hardening/                   ← Klasör: OS sertleştirme
│   ├── WINDOWS11_HARDENING_KALE.md
│   ├── WINDOWS10_HARDENING_KALE.md
│   └── LINUX_HARDENING_KALE.md
│
├── disk-security/                  ← Klasör: Disk ve şifreleme
│   ├── DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md
│   ├── SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md
│   ├── HARICI_USB_SSD_BOOT_REHBERI.md
│   └── VENTOY_WTG_MULTIBOOT_REHBERI.md
│
├── tools/                          ← Klasör: Yazılım ve araçlar
│   └── GHOSTGUARD_MIMARI_VE_PLAN.md
│
├── scripts/                        ← Klasör: Otomatik sertleştirme
│   ├── kale_harden_win11.ps1       ← YENİ
│   ├── kale_harden_win10.ps1       ← YENİ
│   ├── kale_harden_ubuntu.sh       ← YENİ
│   ├── kale_durum_win.ps1          ← YENİ
│   └── kale_durum_linux.sh         ← YENİ
│
└── network/                        ← Klasör: Ağ güvenliği
    └── AG_GUVENLIGI_VE_VPN.md      ← YENİ
```

---

## 🏆 BÖLÜM 4 — ÖNCELİK SIRASI (Ne İlk Yapılmalı?)

```
🔴 YÜKSEK ÖNCELİK (İlk Yapılacaklar):
  1. Sertleştirme scriptleri (scripts/) oluştur — en büyük pratik fayda
  2. Win10 EOL uyarısı ekle
  3. Ağ Güvenliği rehberi oluştur
  4. GhostGuard'ı kodlamaya başla
  5. Yedekleme ve Felaket Kurtarma rehberi

🟡 ORTA ÖNCELİK:
  6. Şifre Yönetimi ve 2FA rehberi
  7. Forensic ve Incident Response rehberi
  8. Linux rehberine Lynis ve USBGuard ekle
  9. Tarayıcı Güvenliği rehberi
  10. GhostGuard'a Telegram komut desteği

🟢 DÜŞÜK ÖNCELİK:
  11. Repo yapısını klasörlere ayır
  12. Mobil güvenlik rehberi
  13. Kriptografi rehberi
  14. README badge'leri ve contributing
  15. Windows Sandbox bölümü
```

---

## 💬 BÖLÜM 5 — SPESİFİK TEKNİK NOTLAR

### Gözlemlenen Güçlü Yönler
- **Mimari tutarlılık:** 7 katmanlı kale modeli tüm belgelerde tutarlı — harika
- **Doğrulama komutları:** Her adımda `Get-BitLockerVolume`, `aa-status` gibi kontrol komutları var
- **Tehdit modeli odaklı:** "Bu neden yapılıyor?" sorusuna çoğu yerde cevap verilmiş
- **Event ID referansları:** Windows Event Log ID'leri çok kapsamlı
- **Cross-platform:** Aynı güvenlik prensibi hem Windows hem Linux için uygulanmış

### Gözlemlenen Tekrarlar ve Çakışmalar
```
⚠️ BitLocker yapılandırması: 4 farklı dosyada tekrarlanıyor
   → SIFRE_KRONOLOJISI, DUAL_BOOT, WIN10_HARDENING, WIN11_HARDENING
   → Çözüm: Ana açıklama bir dosyada, diğerleri çapraz referans

⚠️ LUKS2 kurulumu: 2 dosyada tekrar
   → DUAL_BOOT, LINUX_HARDENING
   → Çözüm: Linux rehberine ana referans, dual boot'ta kısa özet + link

⚠️ BIOS ayarları: Neredeyse her dosyada
   → Çözüm: SIFRE_KRONOLOJISI'ni ana kaynak yap, diğerleri "bkz." desin
```

---

*Bu analiz 2026-03-22 tarihinde güncel dosya içerikleri üzerinden yapılmıştır.*
