# 💎 HAYAT KURTARAN YAZILIMLAR — Format Sonrası Köşetaşları
## Premium Deneyim · Açık Kaynak · Ücretsiz · Reklamsız · Bloatware'siz

> **Felsefe:** Her kategoride **tek doğru araç** seçildi. Alternatif cehennemine düşmeden, yıllardır topluluk tarafından kanıtlanmış, güncellenmeye devam eden, reklamsız ve güvenilir yazılımlar.

> [!TIP]
> 📜 Bu listenin bir sonraki evrimi: **tek PowerShell scripti ile format sonrası hepsini otomatik kuran** `kale_format_sonrasi.ps1` olacak. (`winget` + `choco` tabanlı)

---

## 📐 KURULUM ÖNCELİK SIRASI

```
Format atıldı → Windows kuruldu → İnternet bağlandı → Sırasıyla:

AŞAMA 1: Temel Altyapı (ilk 10 dakika)
  → Tarayıcı, şifre yöneticisi, arşivleme, runtime'lar

AŞAMA 2: Sistem Araçları (ilk 30 dakika)
  → Disk analiz, temizlik, donanım izleme, terminal

AŞAMA 3: Üretkenlik (ilk 1 saat)
  → Metin editörü, ofis, PDF, not alma, ekran görüntüsü

AŞAMA 4: Medya ve İçerik (ilk 2 saat)
  → Video oynatıcı, codec, görüntü düzenleme, ses

AŞAMA 5: Geliştirici Araçları (gerektiğinde)
  → Git, IDE, Python, terminal, SSH

AŞAMA 6: Güvenlik ve Gizlilik (her zaman)
  → VPN, şifreleme, güvenlik duvarı, sandbox
```

---

## AŞAMA 1 — TEMEL ALTYAPI

### 🌐 Tarayıcı

| Yazılım | Açıklama | Neden Bu? | Alternatif |
|---------|----------|-----------|------------|
| **[Firefox](https://www.mozilla.org/firefox/)** | Açık kaynak tarayıcı | Google'a bağımlı değil, `about:config` ile tam kontrol, uBlock Origin desteği, container tabs | — |
| **[Brave](https://brave.com/)** | Chromium tabanlı gizlilik odaklı | Reklam engelleyici dahili, Tor sekmesi, fingerprint koruması | Ungoogled Chromium |

> [!NOTE]
> **İkisini de kur.** Firefox günlük kullanım + gizlilik, Brave Chromium gerektiren siteler için yedek.

**Firefox Olmazsa Olmaz Eklentiler:**

| Eklenti | Ne Yapar |
|---------|----------|
| **uBlock Origin** | Reklam + tracker engelleyici (tüm zamanların en iyisi) |
| **Bitwarden** | Şifre yöneticisi tarayıcı entegrasyonu |
| **Dark Reader** | Her siteyi karanlık moda çevir (göz sağlığı) |
| **ClearURLs** | URL'lerdeki tracking parametrelerini temizle |
| **Privacy Badger** | EFF'nin tracker öğrenen engelleyicisi |
| **Multi-Account Containers** | Farklı hesapları izole sekmelerde kullan |

---

### 🔑 Şifre Yöneticisi

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[Bitwarden](https://bitwarden.com/)** | Bulut senkronizasyonlu şifre yöneticisi | Açık kaynak, ücretsiz, tüm platformlar, TOTP desteği, self-host edilebilir |
| **[KeePassXC](https://keepassxc.org/)** | Offline şifre yöneticisi | Buluta çıkmaz, yerel `.kdbx` dosyası, TOTP, SSH agent, tarayıcı entegrasyonu |

```
Strateji:
  Bitwarden → günlük web şifreleri (otomatik doldurma)
  KeePassXC → kritik şifreler (BitLocker, BIOS, kurtarma anahtarları)
              USB'de taşınabilir versiyonu da var
```

---

### 📦 Arşivleme

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[7-Zip](https://7-zip.org/)** | Arşiv açma/oluşturma | Açık kaynak, reklamsız, tüm formatlar (zip, rar, 7z, tar, gz, xz), AES-256 şifreleme. WinRAR'ın ücretsiz ve daha güçlü versiyonu |

```
winget install -e --id 7zip.7zip
```

> [!IMPORTANT]
> **WinRAR kullanma.** 7-Zip her şeyi yapar, ücretsizdir, reklamsızdır ve güvenlik açısından daha sık güncellenir.

---

### ⚙️ Runtime'lar ve Temel Bileşenler

| Yazılım | Neden Gerekli |
|---------|---------------|
| **[Visual C++ Redistributable (All-in-One)](https://github.com/abbodi1406/vcredist)** | Neredeyse tüm yazılımlar için gerekli DLL'ler |
| **[.NET Desktop Runtime](https://dotnet.microsoft.com/download)** | .NET uygulamaları için |
| **[Java (Adoptium/Temurin)](https://adoptium.net/)** | JDownloader, Minecraft, Java uygulamaları |
| **[Python 3.12+](https://python.org/)** | Script'ler, GhostGuard, otomasyon |
| **[Node.js LTS](https://nodejs.org/)** | Web geliştirme, npm araçları |

```powershell
# Hepsini tek seferde:
winget install -e --id Microsoft.VCRedist.2015+.x64
winget install -e --id Microsoft.DotNet.DesktopRuntime.8
winget install -e --id EclipseAdoptium.Temurin.21.JDK
winget install -e --id Python.Python.3.12
winget install -e --id OpenJS.NodeJS.LTS
```

---

## AŞAMA 2 — SİSTEM ARAÇLARI

### 💽 Disk Analiz ve Temizlik

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[WinDirStat](https://windirstat.net/)** | Disk kullanımını görsel haritayla göster | Hangi dosya/klasör ne kadar yer kaplıyor? Treemap ile anında görüyorsun. Format öncesi temizlik için vazgeçilmez |
| **[WizTree](https://wiztreefree.com/)** | WinDirStat'ın çok daha hızlı versiyonu | NTFS MFT'yi okuduğu için saniyeler içinde tarar. WinDirStat dakikalar alır. **Hız istiyorsan bu** |
| **[BleachBit](https://www.bleachbit.org/)** | Sistem temizleme (CCleaner alternatifi) | Açık kaynak, reklamsız, geçici dosyalar + tarayıcı cache + log temizleme. CCleaner artık Avast bloatware |
| **[Bulk Crap Uninstaller (BCUninstaller)](https://www.bcuninstaller.com/)** | Gelişmiş program kaldırma | Kalıntı dosya ve registry girdilerini de temizler. Windows'un kendi "Programı Kaldır"ı yetersiz |

```powershell
winget install -e --id AntibodySoftware.WizTree
winget install -e --id BleachBit.BleachBit
winget install -e --id Klocman.BulkCrapUninstaller
```

---

### 🔧 Donanım İzleme ve Sistem Bilgisi

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[HWiNFO64](https://www.hwinfo.com/)** | Donanım bilgisi ve sensör izleme | CPU/GPU sıcaklığı, fan hızları, voltaj, SSD ömrü — hepsi gerçek zamanlı. Overclock yapanlar için şart |
| **[CrystalDiskInfo](https://crystalmark.info/en/software/crystaldiskinfo/)** | Disk sağlık durumu (S.M.A.R.T.) | SSD/HDD'nin ömrünü, sıcaklığını, hata sayısını gösterir. Disk ölmeden önce uyarır |
| **[CrystalDiskMark](https://crystalmark.info/en/software/crystaldiskmark/)** | Disk hız testi | Okuma/yazma hızını ölç — USB SSD kutusunun gerçek performansını test et |
| **[CPU-Z](https://www.cpuid.com/softwares/cpu-z.html)** | CPU, RAM, anakart detayları | İşlemci modeli, RAM hızı/tipi, anakart chipset bilgisi |
| **[GPU-Z](https://www.techpowerup.com/gpuz/)** | Ekran kartı detayları | GPU sıcaklığı, VRAM kullanımı, sürücü versiyonu |

```powershell
winget install -e --id REALiX.HWiNFO
winget install -e --id CrystalDewWorld.CrystalDiskInfo
winget install -e --id CrystalDewWorld.CrystalDiskMark
winget install -e --id CPUID.CPU-Z
winget install -e --id TechPowerUp.GPU-Z
```

---

### 🖥️ Terminal ve Dosya Yönetimi

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[Windows Terminal](https://aka.ms/terminal)** | Modern terminal emülatörü | Sekmeli, GPU hızlandırmalı, PowerShell + CMD + WSL tek yerde. Windows'un varsayılan cmd.exe'si karanlık çağdan kalma |
| **[PowerToys](https://github.com/microsoft/PowerToys)** | Microsoft'un güç kullanıcı araç seti | FancyZones (pencere yönetimi), Color Picker, PowerRename, File Locksmith, Paste as Plain Text, Always on Top ve daha fazlası |
| **[Everything](https://www.voidtools.com/)** | Anlık dosya arama | Windows Search saatler alırken Everything milisaniyede bulur. NTFS indeksleme ile çalışır — **oyun değiştirici** |
| **[Files](https://files.community/)** | Modern dosya yöneticisi | Sekmeli, çift panel, modern arayüz. Windows Explorer'ın olması gereken hali |
| **[Total Commander](https://www.ghisler.com/)** | Çift panelli dosya yöneticisi (Power User) | FTP, arşiv, toplu yeniden adlandırma, dosya karşılaştırma — 30 yıllık efsane |

```powershell
winget install -e --id Microsoft.WindowsTerminal
winget install -e --id Microsoft.PowerToys
winget install -e --id voidtools.Everything
```

---

### 🔄 Yedekleme ve Senkronizasyon

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[FreeFileSync](https://freefilesync.org/)** | Dosya senkronizasyon ve yedekleme | İki klasörü karşılaştır ve senkronize et. Harici diske yedekleme için mükemmel |
| **[Macrium Reflect Free](https://www.macrium.com/reflectfree)** | Disk imaj yedekleme | Tam disk görüntüsü al → felaket anında geri yükle |
| **[Syncthing](https://syncthing.net/)** | P2P dosya senkronizasyonu | Bulut yok, sunucu yok — cihazlar arası doğrudan şifreli senkronizasyon |

---

## AŞAMA 3 — ÜRETKENLİK

### 📝 Metin Editörü ve Kod

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[VS Code](https://code.visualstudio.com/)** | Kod editörü | Eklenti ekosistemi devasa, Git entegrasyonu, terminal, hata ayıklama — endüstri standardı |
| **[Notepad++](https://notepad-plus-plus.org/)** | Hızlı metin editörü | Basit dosya düzenleme, log okuma, regex arama — Windows Notepad'in olması gereken hali |
| **[Obsidian](https://obsidian.md/)** | Markdown not alma | Bağlantılı notlar (Zettelkasten), grafik görünüm, eklenti ekosistemi, yerel dosyalar (bulut yok). **İkinci beyin** |

```powershell
winget install -e --id Microsoft.VisualStudioCode
winget install -e --id Notepad++.Notepad++
winget install -e --id Obsidian.Obsidian
```

---

### 📄 Ofis ve PDF

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[LibreOffice](https://www.libreoffice.org/)** | Ofis paketi | Word/Excel/PowerPoint ücretsiz alternatifi. Microsoft Office lisansı yoksa tek seçenek. `.docx`, `.xlsx` tam uyumlu |
| **[Sumatra PDF](https://www.sumatrapdfreader.org/)** | PDF okuyucu | 3 MB, anında açılır, reklamsız. Adobe Reader 300 MB + güvenlik açıkları. **Karşılaştırma bile yok** |
| **[Okular](https://okular.kde.org/)** | Gelişmiş PDF okuyucu | İmzalama, form doldurma, açıklama ekleme gerekiyorsa |
| **[PDF24](https://www.pdf24.org/)** | PDF araç kutusu | Birleştir, böl, sıkıştır, OCR, dönüştür — her şey ücretsiz ve yerelde çalışır |

```powershell
winget install -e --id TheDocumentFoundation.LibreOffice
winget install -e --id SumatraPDF.SumatraPDF
winget install -e --id geek.pdf24
```

---

### 📸 Ekran Görüntüsü ve Kayıt

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[ShareX](https://getsharex.com/)** | Ekran görüntüsü + kayıt + GIF | Bölge seçimi, otomatik upload, OCR, renk seçici, blur, ok çizme — **tek araç her şeyi yapar**. Snipping Tool'un 100 katı |
| **[OBS Studio](https://obsproject.com/)** | Ekran ve kamera kaydı / canlı yayın | Açık kaynak, profesyonel seviye, sahne yönetimi, oyun kaydı, streaming |
| **[ScreenToGif](https://www.screentogif.com/)** | GIF kaydedici ve editör | Ekranı kaydet → GIF olarak düzenle ve kaydet |

```powershell
winget install -e --id ShareX.ShareX
winget install -e --id OBSProject.OBSStudio
winget install -e --id NickeManarin.ScreenToGif
```

---

### 📋 Pano Yönetimi ve Otomasyon

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[Ditto](https://ditto-cp.sourceforge.io/)** | Pano geçmişi yöneticisi | Son 1000 kopyaladığın şeyi saklar. Ctrl+C → istediğin zaman geri getir |
| **[AutoHotkey v2](https://www.autohotkey.com/)** | Klavye/fare otomasyon | Kısayol tuşları, makrolar, pencere yönetimi — tekrarlayan işleri otomatikleştir |

---

## AŞAMA 4 — MEDYA ve İÇERİK

### 🎬 Video

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[VLC Media Player](https://www.videolan.org/)** | Video/ses oynatıcı | Her formatı açar, codec gerektirmez, reklamsız, altyazı desteği. **Efsane** |
| **[mpv](https://mpv.io/)** | Minimalist video oynatıcı | VLC'den daha hafif, komut satırı dostu, GPU hızlandırmalı, config dosyasıyla tam kontrol |
| **[HandBrake](https://handbrake.fr/)** | Video dönüştürücü ve sıkıştırıcı | MKV→MP4, 4K→1080p, codec değiştirme, preset'ler. **Ücretsiz dünyanın en iyi video encoder'ı** |
| **[MKVToolNix](https://mkvtoolnix.download/)** | MKV dosya düzenleyici | Altyazı, ses ekleme/çıkarma, bölüm yönetimi — yeniden encode etmeden |
| **[Subtitle Edit](https://github.com/SubtitleEdit/subtitleedit)** | Altyazı düzenleyici | Zamanlama düzeltme, format dönüştürme, OCR, senkronizasyon |

```powershell
winget install -e --id VideoLAN.VLC
winget install -e --id HandBrake.HandBrake
winget install -e --id MoritzBunkus.MKVToolNix
```

---

### 🎵 Ses

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[Audacity](https://www.audacityteam.org/)** | Ses düzenleme | Kayıt, kesme, efekt, gürültü temizleme — podcast ve ses işleme için standart |
| **[foobar2000](https://www.foobar2000.org/)** | Müzik oynatıcı | Hafif, özelleştirilebilir, FLAC/ALAC desteği, ReplayGain |
| **[EarTrumpet](https://eartrumpet.app/)** | Uygulama bazlı ses kontrolü | Her uygulamanın sesini ayrı ayrı kontrol et — Windows mixer'ının çok daha iyisi |

```powershell
winget install -e --id Audacity.Audacity
winget install -e --id PeterPawlowski.foobar2000
winget install -e --id File-New-Project.EarTrumpet
```

---

### 🖼️ Görüntü

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[GIMP](https://www.gimp.org/)** | Görüntü düzenleme (Photoshop alternatifi) | Katmanlar, maskeler, filtreler — ücretsiz ve açık kaynak |
| **[Paint.NET](https://www.getpaint.net/)** | Hafif görüntü düzenleme | GIMP karmaşık geliyorsa bu. Katmanlar, efektler, şeffaflık — Paint'in olması gereken hali |
| **[IrfanView](https://www.irfanview.com/)** | Hızlı görüntü görüntüleyici | 2 MB, her formatı açar, toplu dönüştürme, basit düzenleme. **30 yıllık efsane** |
| **[ImageMagick](https://imagemagick.org/)** | Komut satırı görüntü işleme | Toplu resize, format dönüştürme, watermark — script'lerle otomasyon |

```powershell
winget install -e --id GIMP.GIMP
winget install -e --id dotPDN.PaintDotNet
winget install -e --id IrfanSkiljan.IrfanView
```

---

## AŞAMA 5 — İNDİRME ve DOSYA AKTARIMI

### ⬇️ İndirme Yöneticileri

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[JDownloader 2](https://jdownloader.org/)** | İndirme yöneticisi | Çoklu bağlantı, otomatik CAPTCHA çözme, paket tanıma, hız sınırı yönetimi, tarayıcı entegrasyonu. **İndirme konusunda rakipsiz** |
| **[Free Download Manager (FDM)](https://www.freedownloadmanager.org/)** | İndirme yöneticisi (alternatif) | Torrent desteği dahil, tarayıcı entegrasyonu, zamanlama |
| **[qBittorrent](https://www.qbittorrent.org/)** | Torrent istemcisi | Açık kaynak, reklamsız, RSS, IP filtreleme. uTorrent artık reklam ve kripto madenciliği dolu — **KULLANMA** |
| **[yt-dlp](https://github.com/yt-dlp/yt-dlp)** | Video indirici (komut satırı) | YouTube, Twitter, Instagram ve 1000+ siteden video/ses indir. youtube-dl'nin aktif fork'u |

```powershell
winget install -e --id AppWork.JDownloader
winget install -e --id qBittorrent.qBittorrent
winget install -e --id yt-dlp.yt-dlp
```

> [!WARNING]
> **uTorrent kullanma.** Arka planda kripto madenciliği yaptığı, reklam gösterdiği ve güvenlik açıkları olduğu tespit edildi. qBittorrent her açıdan üstün.

---

### 📡 Dosya Aktarımı

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[WinSCP](https://winscp.net/)** | SFTP/SCP/FTP istemcisi | Sunucuya dosya aktarımı, anahtar tabanlı kimlik doğrulama, script desteği |
| **[FileZilla](https://filezilla-project.org/)** | FTP/SFTP istemcisi | Daha kullanıcı dostu arayüz, çoklu bağlantı, site yöneticisi |
| **[LocalSend](https://localsend.org/)** | Cihazlar arası dosya aktarımı | AirDrop alternatifi — Windows, Android, iOS, Linux, Mac hepsi arasında WiFi üzerinden dosya gönder. Sunucu yok, internet gerekmiyor |

```powershell
winget install -e --id WinSCP.WinSCP
winget install -e --id LocalSend.LocalSend
```

---

## AŞAMA 6 — GELİŞTİRİCİ ARAÇLARI

### 💻 Geliştirme

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[Git](https://git-scm.com/)** | Versiyon kontrol sistemi | Endüstri standardı. `git bash` ile Linux komutları Windows'ta |
| **[GitHub CLI (gh)](https://cli.github.com/)** | GitHub komut satırı | Repo oluşturma, PR, issue yönetimi terminalden |
| **[WSL2 (Ubuntu)](https://learn.microsoft.com/wsl/)** | Windows üzerinde Linux | Gerçek Linux kernel'i, Docker desteği, tam geliştirme ortamı |
| **[Docker Desktop](https://www.docker.com/products/docker-desktop/)** | Container platformu | İzole geliştirme ortamları, veritabanı, servis yönetimi |
| **[Postman](https://www.postman.com/)** | API test aracı | REST/GraphQL API testi, koleksiyon yönetimi, ortam değişkenleri |
| **[DBeaver](https://dbeaver.io/)** | Veritabanı yöneticisi | MySQL, PostgreSQL, SQLite, MongoDB — hepsini tek arayüzden yönet |

```powershell
winget install -e --id Git.Git
winget install -e --id GitHub.cli
winget install -e --id Docker.DockerDesktop
winget install -e --id dbeaver.dbeaver
```

---

### 🔌 SSH ve Uzak Erişim

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[PuTTY](https://www.putty.org/)** | SSH istemcisi (klasik) | Hafif, taşınabilir, anahtar yönetimi (Pageant) |
| **[MobaXterm](https://mobaxterm.mobatek.net/)** | Gelişmiş SSH istemcisi | X11 forwarding, SFTP, sekmeli SSH, yerleşik araçlar — **PuTTY'nin evrimi** |
| **[Tabby](https://tabby.sh/)** | Modern terminal + SSH | GPU hızlandırmalı, sekmeler, SSH yapılandırma yöneticisi, güzel arayüz |
| **[RustDesk](https://rustdesk.com/)** | Uzak masaüstü (TeamViewer alternatifi) | Açık kaynak, self-host edilebilir, ücretsiz, TeamViewer kısıtlaması yok |

```powershell
winget install -e --id Eugeny.Tabby
winget install -e --id RustDesk.RustDesk
```

---

## AŞAMA 7 — GÜVENLİK ve GİZLİLİK

### 🛡️ Güvenlik Araçları

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[VeraCrypt](https://www.veracrypt.fr/)** | Disk/dosya şifreleme | TrueCrypt'in güvenli devamı. Konteyner, tam disk, gizli volüm |
| **[Cryptomator](https://cryptomator.org/)** | Bulut depolama şifreleme | Google Drive / OneDrive / Dropbox'taki dosyalarını yerel olarak şifreler, buluta şifreli gönderir |
| **[Portmaster](https://safing.io/portmaster/)** | Uygulama güvenlik duvarı | Hangi uygulama nereye bağlanıyor? DNS filtreleme, tracker engelleme — GlassWire alternatifi ama açık kaynak |
| **[Simplewall](https://github.com/henrypp/simplewall)** | Hafif güvenlik duvarı | Windows Filtering Platform tabanlı, basit arayüz, uygulama bazlı kurallar |
| **[Malwarebytes](https://www.malwarebytes.com/)** | Anti-malware tarayıcı | İkinci görüş taraması (Defender'a ek). Ücretsiz sürüm yeterli — manuel tarama |
| **[Sysinternals Suite](https://learn.microsoft.com/sysinternals/)** | Sistem analiz araçları | Process Explorer, Autoruns, Process Monitor, TCPView — **Windows yöneticisinin İsviçre çakısı** |

```powershell
winget install -e --id IDRIX.VeraCrypt
winget install -e --id Cryptomator.Cryptomator
winget install -e --id Safing.Portmaster
winget install -e --id Microsoft.Sysinternals.ProcessExplorer
winget install -e --id Microsoft.Sysinternals.Autoruns
```

---

### 🌐 VPN

| Yazılım | Açıklama | Neden Bu? |
|---------|----------|-----------|
| **[Mullvad VPN](https://mullvad.net/)** | Gizlilik odaklı VPN | Hesap yok, e-posta yok, sadece numara. WireGuard desteği, denetlenmiş, İsveç yasaları |
| **[ProtonVPN](https://protonvpn.com/)** | Ücretsiz VPN (güvenilir) | ProtonMail ekosistemi, ücretsiz sunucular, no-log kanıtlanmış |
| **[WireGuard](https://www.wireguard.com/)** | VPN protokolü | Kendi VPN sunucunu kur (VPS + WireGuard). En hızlı, en modern VPN protokolü |

---

## AŞAMA 8 — YAŞAM KALİTESİ BONUSLARI

### ✨ Günlük Hayatı Premium Yapan Şeyler

| Yazılım | Açıklama | Neden Hayat Kurtarır? |
|---------|----------|----------------------|
| **[f.lux](https://justgetflux.com/)** | Ekran renk sıcaklığı | Gece mavi ışığı azaltır — göz sağlığı ve uyku kalitesi. Windows Night Light'tan çok daha iyi |
| **[Flameshot](https://flameshot.org/)** | Linux'ta ekran görüntüsü | ShareX'in Linux karşılığı. Ok, bulanıklık, metin ekleme |
| **[SpaceSniffer](https://www.intactus.hu/)** | Disk görselleştirme (alternatif) | WinDirStat/WizTree'ye ek — gerçek zamanlı animasyonlu treemap |
| **[Rufus](https://rufus.ie/)** | USB bootable medya oluşturma | Windows/Linux ISO → bootable USB. Ventoy'un tek seferlik alternatifi |
| **[Ventoy](https://ventoy.net/)** | Çoklu USB boot | Bir kez kur, istediğin kadar ISO at — hepsi menüde çıkar |
| **[Stardock Fences](https://www.stardock.com/products/fences/)** | Masaüstü düzenleyici (ücretli) | Masaüstü ikonlarını gruplara ayır, çift tıkla gizle |
| **[Flow Launcher](https://www.flowlauncher.com/)** | Uygulama başlatıcı (Spotlight) | Alt+Space → yaz → aç. Dosya ara, hesap yap, web'de ara — **macOS Spotlight Windows'ta** |
| **[Barrier](https://github.com/debauchee/barrier)** | Yazılımsal KVM switch | Tek klavye-fare ile birden fazla bilgisayarı kontrol et |
| **[PowerShell 7](https://github.com/PowerShell/PowerShell)** | Modern PowerShell | Windows PowerShell 5.1 eski. PSH 7 cross-platform, daha hızlı, daha güçlü |
| **[Winget-UI (UniGetUI)](https://github.com/marticliment/UniGetUI)** | Paket yöneticisi GUI | winget, choco, scoop, pip, npm paketlerini güzel arayüzden yönet. **App Store ama düzgünü** |

```powershell
winget install -e --id flux.flux
winget install -e --id Rufus.Rufus
winget install -e --id Flow-Launcher.Flow-Launcher
winget install -e --id Microsoft.PowerShell
winget install -e --id MartiCliment.UniGetUI
```

---

## 🚫 KULLANMA — BUNLARDAN UZAK DUR

| ❌ Kullanma | Neden | ✅ Yerine Kullan |
|------------|-------|-----------------|
| **WinRAR** | Ücretli, güncellenmiyor, gereksiz | **7-Zip** |
| **CCleaner** | Avast satın aldı, bloatware, reklam | **BleachBit** |
| **uTorrent** | Kripto madenciliği, reklam, güvenlik açığı | **qBittorrent** |
| **Adobe Reader** | 300 MB, yavaş, güvenlik açığı, telemetri | **Sumatra PDF** |
| **TeamViewer** | Kişisel kullanımda kısıtlama, telemetri | **RustDesk** |
| **Internet Download Manager** | Ücretli, crack gerektirir | **JDownloader 2** |
| **McAfee / Norton** | Bloatware, PC yavaşlatır, gereksiz | **Windows Defender + Malwarebytes** |
| **Avast / AVG** | Kullanıcı verisini satıyor (kanıtlanmış) | **Windows Defender** |
| **Opera / Opera GX** | Çin'e satıldı, kredi servisi, güvenilmez | **Firefox / Brave** |
| **Windows Media Player** | Codec desteği kötü, yavaş | **VLC / mpv** |
| **Notepad (Windows)** | Fonksiyon yok | **Notepad++** |

---

## 📜 FORMAT SONRASI OTOMATİK KURULUM — Ön Hazırlık

> Bu bölüm, bir sonraki aşamada hazırlanacak `kale_format_sonrasi.ps1` scriptinin taslağıdır.

```powershell
# ═══════════════════════════════════════════════════════
# 🏰 KANIJE KALESİ — Format Sonrası Kurulum Scripti
# ═══════════════════════════════════════════════════════
# Çalıştırma: PowerShell (Yönetici) → .\kale_format_sonrasi.ps1
# Ön koşul: İnternet bağlantısı + Windows 10/11

# ── winget Doğrulama ──
if (!(Get-Command winget -ErrorAction SilentlyContinue)) {
    Write-Host "winget bulunamadı! Microsoft Store'dan App Installer'ı güncelle." -ForegroundColor Red
    exit 1
}

$apps = @(
    # Aşama 1: Temel
    "Mozilla.Firefox",
    "7zip.7zip",
    "Bitwarden.Bitwarden",
    "KeePassXCTeam.KeePassXC",

    # Aşama 2: Sistem
    "AntibodySoftware.WizTree",
    "BleachBit.BleachBit",
    "Klocman.BulkCrapUninstaller",
    "REALiX.HWiNFO",
    "CrystalDewWorld.CrystalDiskInfo",
    "Microsoft.WindowsTerminal",
    "Microsoft.PowerToys",
    "voidtools.Everything",

    # Aşama 3: Üretkenlik
    "Microsoft.VisualStudioCode",
    "Notepad++.Notepad++",
    "Obsidian.Obsidian",
    "TheDocumentFoundation.LibreOffice",
    "SumatraPDF.SumatraPDF",
    "ShareX.ShareX",

    # Aşama 4: Medya
    "VideoLAN.VLC",
    "HandBrake.HandBrake",
    "Audacity.Audacity",
    "GIMP.GIMP",

    # Aşama 5: İndirme
    "AppWork.JDownloader",
    "qBittorrent.qBittorrent",
    "yt-dlp.yt-dlp",

    # Aşama 6: Geliştirici
    "Git.Git",
    "GitHub.cli",
    "Python.Python.3.12",

    # Aşama 7: Güvenlik
    "IDRIX.VeraCrypt",
    "Cryptomator.Cryptomator",

    # Aşama 8: Yaşam Kalitesi
    "flux.flux",
    "Rufus.Rufus",
    "Flow-Launcher.Flow-Launcher",
    "Microsoft.PowerShell",
    "MartiCliment.UniGetUI"
)

$total = $apps.Count
$current = 0

foreach ($app in $apps) {
    $current++
    $pct = [math]::Round(($current / $total) * 100)
    Write-Host "[$pct%] Kuruluyor: $app" -ForegroundColor Cyan
    winget install -e --id $app --accept-source-agreements --accept-package-agreements -h
}

Write-Host "`n✅ Tüm yazılımlar kuruldu!" -ForegroundColor Green
Write-Host "🏰 Kanije Kalesi hazır." -ForegroundColor Yellow
```

> [!IMPORTANT]
> **Son versiyon script** ayrı bir `kale_format_sonrasi.ps1` dosyasında tutulacak. Yukarıdaki taslak bir başlangıç noktasıdır. İhtiyacına göre listeden ekleme/çıkarma yapabilirsin.

---

## 📊 HIZLI REFERANS — KATEGORİ BAŞINA EN İYİ TEK SEÇİM

Karar veremiyorsan, her kategoride **sadece bir tane** seç:

| Kategori | Tek Seçim | Neden |
|----------|----------|-------|
| Tarayıcı | **Firefox** | Gizlilik + uBlock + Container |
| Şifre | **Bitwarden** | Cross-platform sync |
| Arşiv | **7-Zip** | Tartışmasız |
| Disk Analiz | **WizTree** | WinDirStat'tan 100x hızlı |
| Temizlik | **BleachBit** | CCleaner'ın temiz versiyonu |
| Dosya Arama | **Everything** | Milisaniye |
| Metin Editörü | **VS Code** | Editör + terminal + Git = tek araç |
| Not Alma | **Obsidian** | İkinci beyin |
| PDF | **Sumatra PDF** | 3 MB, anlık açılış |
| Ekran Görüntüsü | **ShareX** | Her şeyi yapar |
| Video Oynatıcı | **VLC** | Her format, her codec |
| Video Encode | **HandBrake** | Endüstri standardı |
| İndirme | **JDownloader 2** | Rakipsiz |
| Torrent | **qBittorrent** | Temiz, reklamsız |
| Uzak Erişim | **RustDesk** | Açık kaynak TeamViewer |
| Şifreleme | **VeraCrypt** | Askeri standart |
| Power User | **PowerToys** | Microsoft'tan, her şey dahil |

---

*Son güncelleme: 2026-03-22 · Windows 10/11 için hazırlanmıştır.*
*Bir sonraki adım: `kale_format_sonrasi.ps1` — tek script ile hepsini kur.*
