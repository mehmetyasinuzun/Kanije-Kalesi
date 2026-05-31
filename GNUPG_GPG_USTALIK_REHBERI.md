# 🔑 GNUPG (GPG / OpenPGP) — TAM USTALIK REHBERİ
## Ana Anahtardan Air-Gapped Subkey'lere, YubiKey'e ve Püf Noktalarıyla Uçtan Uca

> **Amaç:** GnuPG'yi "şifrele-imzala" seviyesinden çıkarıp, bir tehdit aktörü karşısında **gerçekten ayakta kalacak** bir anahtar mimarisi kurmayı öğretmek. Bu rehber yalnızca *nasıl*'ı değil, **neden** ve **hangi durumda işe yaramaz**'ı da anlatır. Forum cevaplarında bulamayacağın ana anahtar/subkey ayrımı, air-gapped çevrimdışı master, revocation cert disiplini, paperkey ile kağıda yedek, gpg.conf sertleştirme, forward secrecy gerçeği ve metadata sızıntısı detayları burada.

> ⚠️ **Önce bunu oku:** GPG, *yanlış kullanıldığında* sana **yanlış güvenlik hissi** verir — bu, hiç kullanmamaktan daha tehlikelidir. Özel anahtarı yedeklemeden üretmek, revocation sertifikası olmadan anahtar yaymak, ya da expired bir anahtara güvenmek — hepsi gerçek dünyada insanları yaktı. Özellikle **Ana Anahtar / Subkey** ve **Revocation** bölümlerini atlama.

---

## 📑 İÇİNDEKİLER

1. [GnuPG / OpenPGP Nedir, Neden? (PGP vs GPG)](#1)
2. [Tehdit Modeli — Neyi Korur, Neyi KORUMAZ](#2)
3. [Kurulum (Gpg4win / Linux / macOS)](#3)
4. [Temel Kavramlar — Asimetrik Şifreleme, İmza, Anahtar Halkası](#4)
5. [Algoritma Seçimi (RSA-4096 vs ED25519/Curve25519)](#5)
6. [Anahtar Çifti Üretme (Adım Adım)](#6)
7. [🔥 Ana Anahtar + Subkey Mimarisi — En Kritik Konu](#7)
8. [🔥 Ana Anahtarı Air-Gapped (Çevrimdışı) Tutma](#8)
9. [Revocation Certificate — Önceden Üret, Güvenli Sakla](#9)
10. [Web of Trust, Anahtar İmzalama, Keyserver, WKD](#10)
11. [Son Kullanma, Yenileme, İptal](#11)
12. [Pratik Kullanım — Şifrele / İmzala / Doğrula / Symmetric](#12)
13. [🔑 YubiKey / Akıllı Kart ile Donanımda Anahtar](#13)
14. [gpg-agent, SSH Auth, Parola Önbelleği](#14)
15. [🔥 PÜF NOKTALARI — Piyasada Bulamayacakların](#15)
16. [gpg.conf Sertleştirme (Hardening)](#16)
17. [Yedekleme & Kurtarma (paperkey dahil)](#17)
18. [🏰 Kanije Kalesi ile Birlikte Kullanım](#18)
19. [Yaygın Ölümcül Hatalar](#19)
20. [Hızlı Referans & Operasyonel Kontrol Listesi](#20)

---

<a id="1"></a>
## 1. 🧭 GnuPG / OpenPGP Nedir, Neden?

**OpenPGP** (RFC 4880, güncellemesi RFC 9580) bir **standarttır** — açık anahtarlı şifreleme, dijital imza ve kimlik için bir mesaj/anahtar formatı tanımlar. **PGP** ("Pretty Good Privacy"), Phil Zimmermann'ın 1991'de yazdığı orijinal (artık ticari, Symantec/Broadcom) uygulamadır. **GnuPG** (GPG, "GNU Privacy Guard"), OpenPGP standardının **özgür ve açık kaynak** uygulamasıdır — Werner Koch tarafından sürdürülür, fiilen evrensel standart hâline gelmiştir.

> 🧠 **Tek cümle:** *OpenPGP = standart, PGP = orijinal ticari yazılım, GPG = onun açık kaynak uygulaması.* Bugün "PGP ile imzaladım" diyen birinin %99'u aslında **GPG** kullanır.

### Neyi çözer — üç temel işlev

GPG **asimetrik (açık anahtarlı)** kriptografiye dayanır: herkesin bir **açık anahtarı** (paylaşılır) ve bir **özel anahtarı** (asla paylaşılmaz) vardır.

| İşlev | Nasıl | Ne sağlar |
|---|---|---|
| **Şifreleme** | Alıcının **açık** anahtarıyla şifrele → yalnızca alıcının **özel** anahtarı çözer | İçerik gizliliği (confidentiality) |
| **İmzalama** | Kendi **özel** anahtarınla imzala → herkes senin **açık** anahtarınla doğrular | Kaynak doğrulama + bütünlük + inkâr edilemezlik |
| **Kimlik doğrulama (auth)** | Özel anahtarla kimliğini kanıtla (örn. SSH) | Erişim kontrolü |

### Hibrit şifreleme (kapağın altında ne döner)

GPG aslında saf asimetrik değil **hibrittir**: her mesaj için rastgele bir **oturum anahtarı** (session key, simetrik — AES) üretir, veriyi onunla şifreler, sonra **yalnızca o küçük oturum anahtarını** alıcının açık anahtarıyla şifreler. Bu yüzden gigabaytlık dosyaları hızlıca şifreleyebilir (asimetrik yalnızca minik anahtara uygulanır).

### Uçtan uca (end-to-end) felsefesi

GPG **uçtan uca şifrelemedir**: veri, *senin makinende* şifrelenir, *alıcının makinesinde* çözülür. Arada hiçbir sunucu (e-posta sağlayıcı, keyserver) içeriği göremez. Sunucuya güven gerektirmez — matematiğe güvenir.

---

<a id="2"></a>
## 2. 🎯 Tehdit Modeli — Neyi Korur, Neyi KORUMAZ

Anahtar üretmeden önce **kime karşı** korunduğunu netleştir. Aksi halde yanlış güvenlik hissi yaşarsın.

### ✅ GPG'nin KORUDUĞU senaryolar
- **İçerik gizliliği (transit + at-rest):** E-posta/dosya geçişte yakalansa bile, özel anahtar olmadan içerik okunamaz. Sunucu (Gmail/ProtonMail dahi) içeriği göremez.
- **Kaynak doğrulama (authentication):** İmza, mesajın **gerçekten o anahtarın sahibinden** geldiğini kanıtlar. Sahte gönderici çürür.
- **Bütünlük (integrity):** İmza, mesajın **bir bit bile değişmediğini** kanıtlar. Ortadaki adam içeriği kurcalarsa imza geçersizleşir.
- **İnkâr edilemezlik (non-repudiation):** İmzalayan "ben imzalamadım" diyemez (özel anahtar yalnızca onda). *CTI raporu, sürüm dağıtımı, kanıt zinciri için kritik.*
- **Yazılım/sürüm bütünlüğü:** Bir exe/paket imzalanırsa, indiren kişi truva atı enjekte edilmediğini doğrular (VeraCrypt, Tor, Linux paketleri böyle çalışır).

### ❌ GPG'nin KORUMADIĞI senaryolar
- **Metadata (en büyük zafiyet):** GPG **kim kime, ne zaman, ne boyutta** mesaj gönderdiğini **gizlemez.** E-posta başlıkları (From/To/Subject/tarih) açıkta kalır. *Konu satırı şifrelenmez!* Trafik analizi seni ele verebilir.
- **Forward secrecy YOK:** GPG'nin **ileri gizliliği yoktur.** Özel anahtarın *bir gün* ele geçerse, geçmişte o anahtarla şifrelenmiş **TÜM eski mesajlar** geriye dönük çözülebilir. (Signal/TLS'in aksine — orada her oturum ayrı geçici anahtar kullanır.) Bu, GPG'nin yapısal en zayıf noktasıdır.
- **Uç nokta uzlaşması (endpoint compromise):** Makinende keylogger/RAT varsa, parolanı ve çözülmüş içeriği yakalar. Şifreleme burada çaresizdir.
- **Anahtar uzlaşması:** Özel anahtar dosyan + parolan çalınırsa, saldırgan **sensin.** Anahtarın güvenliği = her şeyin güvenliği.
- **Anahtar doğrulama hatası:** Yanlış kişinin "açık anahtarına" güvenirsen (MITM), ona şifrelersin → saldırgan okur. Bu yüzden **parmak izi doğrulama** (Web of Trust) şart.
- **Zorla parola alma (rubber-hose):** Seni parolanı söylemeye zorlarlarsa matematik çaresizdir.

> 🧠 **Altın kural:** GPG **"içeriği"** korur, **"iletişimin var olduğunu"** korumaz. Hassas, *anlık*, *inkâr edilebilir* iletişim için (forward secrecy + metadata minimizasyonu istiyorsan) **GPG değil Signal** kullan. GPG'nin güçlü olduğu yer: **at-rest dosyalar, imzalı kanıt/sürüm bütünlüğü, eşzamansız (asenkron) e-posta.**

---

<a id="3"></a>
## 3. 📥 Kurulum

| Platform | Paket | Not |
|---|---|---|
| **Windows** | **Gpg4win** (gpg4win.org) — GnuPG + Kleopatra GUI + GpgOL (Outlook) | Resmi siteden indir, **imzasını/SHA'sını doğrula** |
| **Linux** | `gnupg` (çoğu dağıtımda kurulu) | `sudo apt install gnupg` / `dnf install gnupg2` |
| **macOS** | **GPGTools** (gpgtools.org) ya da `brew install gnupg` | GPGTools, Mail.app eklentisi getirir |

```bash
# Sürümü ve desteklenen algoritmaları doğrula
gpg --version
```

> 🔑 **Püf:** Modern özellikler (ED25519, gelişmiş subkey yönetimi) için **GnuPG 2.2+**, tercihen **2.4+** kullan. `gpg --version` çıktısının başındaki sürümü teyit et. Çok eski sürümler (1.4) hâlâ dolaşımda — **kullanma**, modern algoritmaları ve güvenli varsayılanları desteklemez.

> ⚠️ **Tedarik zinciri:** Gpg4win/GPGTools'u **yalnızca resmi siteden** indir ve mümkünse mevcut bir GPG ile imzasını doğrula. Truva atlı bir GPG, ürettiğin tüm anahtarları baştan ele verir.

---

<a id="4"></a>
## 4. 🧱 Temel Kavramlar

### Anahtar çifti
- **Özel anahtar (secret/private key):** ASLA paylaşılmaz. Şifreyi çözer, imza atar. Parolayla (passphrase) korunur.
- **Açık anahtar (public key):** Herkese dağıtılır. Şifreler, imza doğrular.

### Anahtar halkaları (keyrings)
- **Açık anahtar halkası:** Başkalarının açık anahtarları (`pubring.kbx`).
- **Özel anahtar halkası:** Kendi özel anahtarların (modern GnuPG'de `private-keys-v1.d/` klasöründe, gpg-agent yönetir).

### Kimlik (UID)
Bir anahtara bağlı **isim + e-posta** ("Yasin Uzun <yasin@ornek.com>"). Bir anahtarın **birden çok UID'si** olabilir (iş + kişisel e-posta aynı anahtarda).

### Parmak izi (fingerprint) — kimliğin TEK gerçek kanıtı
Her anahtarın benzersiz bir **40 haneli parmak izi** vardır (anahtarın SHA-1/SHA-256 özeti). "Key ID" (son 8/16 hane) **çakıştırılabilir** (collision saldırısı) — **asla yalnızca Key ID'ye güvenme.** Bir kişinin anahtarını doğrularken **tam 40 haneli parmak izini** karşılaştır (telefon/yüz yüze/güvenli kanaldan).

```bash
# Tam parmak izini göster
gpg --fingerprint yasin@ornek.com
```

### Anahtarın yolculuğu (şifreleme — basitleştirilmiş)
```
Düz metin
   │
   ▼  Rastgele oturum anahtarı (AES-256) üretilir
Şifreli içerik (simetrik, hızlı)
   │
   ▼  Oturum anahtarı, ALICININ AÇIK anahtarıyla şifrelenir
Paket = [açık-anahtarla-şifreli oturum anahtarı] + [AES'le-şifreli veri]
   │
   ▼  Alıcı: ÖZEL anahtarıyla oturum anahtarını çözer → veriyi açar
Düz metin (yalnızca alıcının makinesinde)
```

---

<a id="5"></a>
## 5. 🔐 Algoritma Seçimi (RSA-4096 vs ED25519)

### Modern öneri: ED25519 + Curve25519

| Algoritma | İmza/Auth | Şifreleme | Notlar |
|---|---|---|---|
| **ED25519 / cv25519** ⭐ | ED25519 (EdDSA) | Curve25519 (ECDH) | **Modern tercih.** Kısa anahtar, çok hızlı, güçlü güvenlik. GnuPG 2.1.17+ |
| **RSA-4096** | RSA | RSA | Yaygın, çok uyumlu, **yavaş ve büyük.** Eski sistemler için |
| **RSA-2048** | RSA | RSA | Asgari kabul edilebilir; uzun ömürlü sır için **yetersiz**, tercih etme |
| ** NIST P-curve'ler** | ECDSA/ECDH | — | NIST eğrilerine (P-256 vb.) bazıları güvenmez (sabit kaynaklı şüphe); 25519 tercih edilir |

> 🔥 **Püf — hangisi?**
> - **Yeni anahtar üretiyorsan ve karşı taraflar modern GPG kullanıyorsa: ED25519 (imza/auth) + Curve25519 (şifreleme).** Daha küçük, daha hızlı, en az RSA-4096 kadar güçlü ve gelecek-uyumlu.
> - **Çok eski/heterojen sistemlerle (legacy mail gateway, eski cihaz) uyum şartsa: RSA-4096.** Daha yavaş ama her yerde çalışır.
> - **RSA-2048'i uzun ömürlü/devlet-seviyesi tehdit için seçme.** Asgari, ideal değil.

> 🧠 **Teyit et:** ED25519 desteği `gpg --version` çıktısında "Pubkey:" satırında "EDDSA, ECDH" görünmeli (GnuPG 2.1.17+). Görünmüyorsa GnuPG'yi güncelle.

### Anahtar geçerlilik süresi
- **Ana anahtar:** Genelde **uzun ömür** ya da süresiz (çünkü çevrimdışı saklanır, sızma riski düşük). Bazıları yine de 2-5 yıl koyar.
- **Subkey'ler:** **Sınırlı süre koy** (1-2 yıl). Süre dolunca yenilersin → bir subkey sızsa bile penceresi sınırlıdır.

> 💡 Süre dolması **felaket değil yenilenebilir**: ana anahtar elindeyse subkey son kullanma tarihini her zaman ileri atabilirsin. *Süresiz anahtardansa "süre dolan ama yenilenen" anahtar daha disiplinlidir* — ama son kullanmayı **takvimine kaydet**, yoksa bir gün imzaların geçersiz görünür.

---

<a id="6"></a>
## 6. 🛠️ Anahtar Çifti Üretme (Adım Adım)

### Yol A — Hızlı/rehberli (yeni başlayan)
```bash
# Tam menülü, kontrollü üretim (önerilen)
gpg --full-generate-key
```
- Algoritma sorulduğunda: modern için **"ECC (sign and encrypt)" → Curve 25519**; uyum için **"RSA and RSA" → 4096**.
- Geçerlilik: ana anahtar için makul bir süre (ör. 2y) ya da süresiz.
- İsim + e-posta gir.
- **Güçlü passphrase** belirle (aşağıya bak).

### Yol B — Uzman (expert modu, ince kontrol)
```bash
gpg --expert --full-generate-key
# "(11) ECC (set your own capabilities)" ile yetenekleri elle seç,
# ana anahtara yalnızca Certify (sertifikalama) yetkisi bırakmak için (Bkz. Bölüm 7)
```

### Passphrase (anahtarın parolası) — kritik
- **Uzun olsun:** 5-7 rastgele kelime (diceware) — kısa karmaşık paroladan iyidir.
- **Benzersiz:** Başka hiçbir yerde kullanma.
- **Bunu kaybedersen özel anahtar işe yaramaz** (yedek parola yok). Ezberle + güvenli parola yöneticisinde tut.

> 🔥 **Püf:** Passphrase, özel anahtar dosyasını *at-rest* korur. Dosya çalınsa bile parola olmadan açılamaz. Ama **keylogger** parolayı yakalar — passphrase, uç-nokta uzlaşmasına karşı korumaz (Bölüm 2). Bu yüzden **donanım anahtarı** (YubiKey, Bölüm 13) bir üst seviyedir: özel anahtar makineyi hiç terk etmez.

### Üretildikten sonra hemen yap
1. **Revocation certificate üret** (Bölüm 9) — *anahtarı yaymadan ÖNCE.*
2. **Özel anahtarı + revocation cert'i yedekle** (Bölüm 17).
3. **Subkey mimarisini kur** (Bölüm 7).

---

<a id="7"></a>
## 7. 🔥 Ana Anahtar + Subkey Mimarisi — En Kritik Konu

Bu, GPG ustalığının **kalbidir** ve çoğu rehberin atladığı yerdir. Doğru kurulduğunda, bir cihazın çalınması **kimliğini yok etmez.**

### Sorun
Eğer tek bir "her şeyi yapan" anahtarın varsa ve o cihaz çalınırsa:
- Saldırgan senin **kimliğin** olur (artık senin adına imzalar).
- Tüm Web of Trust ilişkilerin (insanların senin anahtarına attığı imzalar) **çöp olur** — yeni anahtar üretip baştan herkese imzalatman gerekir.

### Çözüm: Ana anahtar yalnızca SERTİFİKALAR, gerisini subkey'ler yapar

OpenPGP anahtarları **yeteneklere** ayrılır:
- **[C] Certify (Sertifikala):** Yalnızca ana anahtarda. *Diğer anahtarları imzalama* ve *kendi subkey'lerini oluşturma/iptal etme* yetkisi. **Kimliğinin köküdür.**
- **[S] Sign (İmzala):** Mesaj/dosya/yazılım imzalama.
- **[E] Encrypt (Şifrele):** Sana gelen veriyi çözme.
- **[A] Authenticate (Kimlik doğrula):** SSH gibi erişim.

**İdeal mimari:** Ana anahtar = **yalnızca [C]**. Günlük işler için ayrı [S], [E], [A] **subkey'leri**.

```
                  ┌─────────────────────────────────────┐
                  │   ANA ANAHTAR  [C] (Certify only)    │
                  │   ↳ KİMLİĞİNİN KÖKÜ                   │
                  │   ↳ ÇEVRİMDIŞI / AIR-GAPPED saklanır  │
                  │   ↳ Günlük makinede BULUNMAZ          │
                  └──────────────┬──────────────────────┘
                                 │ imzalar / oluşturur
              ┌──────────────────┼──────────────────────┐
              ▼                  ▼                      ▼
        ┌───────────┐     ┌───────────┐         ┌───────────┐
        │ Subkey [S]│     │ Subkey [E]│         │ Subkey [A]│
        │ İmzalama  │     │ Şifreleme │         │  SSH/Auth │
        └───────────┘     └───────────┘         └───────────┘
         (günlük makinede / YubiKey'de — bunlar çalınsa kimlik ölmez)
```

### Neden bu kadar güçlü?
- **Subkey çalınırsa:** Ana anahtar elinde → o subkey'i **iptal edersin**, yenisini üretirsin. **Kimliğin (ana anahtar) ve tüm WoT imzaların sağlam kalır.** Felaket → küçük olaya dönüşür.
- **Ana anahtar çevrimdışı** olduğu için günlük saldırı yüzeyinde **hiç bulunmaz.**

### Subkey'leri oluşturma
```bash
# Anahtarı düzenleme moduna gir
gpg --expert --edit-key yasin@ornek.com

# (edit-key içinde)
gpg> addkey          # yeni subkey ekle
# → ECC sign / ECC encrypt / set-your-own (auth için) seç
# → süre ver (örn. 1-2 yıl)
gpg> addkey          # gerekli her yetenek için tekrarla (S, E, A)
gpg> save
```

> 🔥 **Püf — ana anahtardan [S] ve [E] yeteneğini kaldırmak:** `gpg --full-generate-key` çoğu zaman ana anahtara hem C hem S verir. Saf "[C] only" istersen `--expert` ile özel yetenek (set your own capabilities) seçip yalnızca Certify bırak; imzalamayı ayrı bir [S] subkey'e yaptır. Bu, ileri seviye disiplinin işaretidir.

> 💡 **Pratik fayda:** Dışarıya verdiğin/keyserver'a koyduğun **açık anahtar**, ana anahtar + tüm subkey'leri içerir. İnsanlar sana **şifreleme subkey'inle** şifreler, sen imzaları **imza subkey'inle** atarsın — ama hepsinin **kökü** çevrimdışı duran ana anahtardır.

---

<a id="8"></a>
## 8. 🔥 Ana Anahtarı Air-Gapped (Çevrimdışı) Tutma

Subkey mimarisinin tamamlayıcısı: **ana anahtar özel parçasını günlük makinende TUTMA.** İnternete bağlı makinede yalnızca subkey'ler dursun.

### Hedef durum
- **Günlük makine:** Yalnızca **subkey'lerin özel parçaları** (S, E, A). Ana anahtarın özel parçası **YOK** (yerine bir "stub" / placeholder kalır).
- **Air-gapped ortam:** Ana anahtarın özel parçası — *yalnızca* subkey oluşturma/iptal, başka anahtar imzalama gerektiğinde kullanılır. Sızma yüzeyi ~sıfır.

### Air-gapped ortam ne olabilir?
- İnternete **hiç bağlanmamış** bir laptop ya da Raspberry Pi.
- **Tails** USB (amnezik; kapatınca iz bırakmaz).
- Şifreli bir USB/VeraCrypt birimi (bkz. `VERACRYPT_USTALIK_REHBERI.md`) içinde saklanan, yalnızca offline makinede açılan keyring.

### Subkey'leri dışa aktarıp ana anahtarı çıkarma (kavramsal akış)

> ⚠️ **Tam komut dizisi sürüme göre incelik gösterir — kendi GnuPG sürümünde resmî dokümandan teyit et.** Genel akış:

```bash
# 1) (AIR-GAPPED makinede) TÜM gizli anahtarları yedekle — felaket sigortası
gpg --export-secret-keys --armor yasin@ornek.com > master-full-backup.asc

# 2) (AIR-GAPPED) Yalnızca SUBKEY gizli parçalarını dışa aktar
gpg --export-secret-subkeys --armor yasin@ornek.com > subkeys-only.asc

# 3) Açık anahtarı dışa aktar (taşımak için)
gpg --export --armor yasin@ornek.com > public.asc

# --- subkeys-only.asc + public.asc'yi GÜNLÜK makineye taşı (USB ile) ---

# 4) (GÜNLÜK makinede) içe aktar
gpg --import public.asc
gpg --import subkeys-only.asc

# 5) (GÜNLÜK makinede) ana anahtarın GİZLİ parçasını SİL,
#    yalnızca subkey gizli parçaları + ana anahtar STUB'ı kalsın
#    (Bunu doğru keygrip ile yap; yanlış silme veri kaybıdır — TEYİT ET)
```

> 🔥 **Doğrulama püfü:** Günlük makinede `gpg --list-secret-keys` çalıştır. Ana anahtar satırında **`sec#`** görmelisin — `#` işareti "özel parça burada YOK, yalnızca stub" demektir. Subkey'lerde `ssb` (# yok) görünür = özel parça mevcut. **`sec#` görüyorsan kurulum doğru.** `sec` (# yok) görüyorsan ana anahtarın hâlâ bu makinede — air-gap kurulmamış.

```bash
# Stub kontrolü
gpg --list-secret-keys --keyid-format long
# sec#  ed25519/...   <-- # = ana anahtar gizli parçası burada DEĞİL (DOĞRU)
# ssb   cv25519/...   <-- subkey gizli parçası mevcut
```

> 💡 **Operasyonel rutin:** Yeni subkey gerekince ya da birinin anahtarını imzalaman gerekince → air-gapped makineyi aç, ana anahtarı kullan, çıktıyı USB ile taşı, offline makineyi kapat. Ana anahtar **internete bağlı bir ortamı asla görmez.** Bu, gazeteci/CTI analisti/yüksek-değerli hedef için en güçlü pratik duruştur.

---

<a id="9"></a>
## 9. 🔁 Revocation Certificate — Önceden Üret, Güvenli Sakla

**Revocation certificate (iptal sertifikası)**, anahtarını "artık geçersiz / ele geçti / kayboldu" diye **dünyaya duyurma** aracıdır. En kritik özelliği: **anahtar elinde olmasa/parolanı unutsan bile** anahtarı iptal edebilmen için **ÖNCEDEN** üretilmesi gerekir.

### Neden önceden?
Eğer özel anahtarını kaybedersen ya da parolanı unutursan, **o anda** iptal sertifikası **üretemezsin** (üretmek için özel anahtar + parola lazım). İnsanlar hâlâ senin (artık erişemediğin) anahtarına şifreler — sen okuyamazsın, onlar bilmez. Önceden üretilmiş bir revocation cert bunu çözer: yayınlarsın, herkes "bu anahtar ölü" görür.

```bash
# Anahtarı üretir üretmez, HEMEN revocation cert üret
gpg --output revoke-yasin.asc --gen-revoke yasin@ornek.com
# Sebep sorulur: 0=belirtilmemiş, 1=ele geçti, 2=yenisiyle değişti, 3=artık kullanılmıyor
```

### Nereye saklamalı?
- **Anahtardan AYRI**, güvenli bir yerde.
- 🔥 **Püf — basılı sakla:** Revocation cert'i **kağıda yazdır** (QR ya da düz metin/base64) ve fiziksel kasada tut. Dijital kopya bozulsa/silinsa bile basılı haldeki geri yazılabilir. *Air-gapped ana anahtar + basılı revocation cert = profesyonel duruş.*
- **Saklama tuzağı:** Revocation cert'i **kim ele geçirirse anahtarını iptal edebilir** (sabotaj — içeriği çözemez ama "anahtarını öldürebilir"). Bu yüzden gizli tut; ama yine de anahtarı *kaybetme* riskinden daha küçük bir risktir.

> ⚠️ **Modern GnuPG kolaylığı:** GnuPG 2.1+ kurulumda otomatik olarak `~/.gnupg/openpgp-revocs.d/<FINGERPRINT>.rev` altında bir revocation cert üretir. **Yine de** bunu güvenli/ayrı bir yere (ve basılı) kopyala — varsayılan konum makineyle birlikte kaybolur/ele geçer.

### İptali uygulama (gün geldiğinde)
```bash
# 1) İptal sertifikasını içe aktar
gpg --import revoke-yasin.asc
# 2) İptal edilmiş anahtarı keyserver'a/WKD'ye yay ki herkes görsün
gpg --keyserver hkps://keys.openpgp.org --send-keys <FINGERPRINT>
```

---

<a id="10"></a>
## 10. 🌐 Web of Trust, Anahtar İmzalama, Keyserver, WKD

GPG'nin **merkezi otoritesi yoktur** (TLS'teki CA'lar gibi). Güven, kullanıcıların birbirinin anahtarını imzalamasıyla yayılır: **Web of Trust (WoT).**

### Mantık
Bir anahtarın "gerçekten o kişiye ait" olduğunu **parmak izini doğrulayıp imzalayarak** onaylarsın. A, B'ye güveniyor ve B, C'nin anahtarını imzaladıysa, A da C'ye (dolaylı) güvenebilir.

### Bir anahtarı doğrulama + imzalama (kritik disiplin)
```bash
# 1) Anahtarı al
gpg --keyserver hkps://keys.openpgp.org --recv-keys <FINGERPRINT>

# 2) PARMAK İZİNİ GÜVENLİ KANALDAN doğrula (telefon, yüz yüze, güvenli mesaj)
gpg --fingerprint <FINGERPRINT>

# 3) Eşleşiyorsa imzala (sertifikala) — bu ana anahtarın [C] yeteneğini gerektirir
gpg --sign-key <FINGERPRINT>
```

> 🔥 **Püf — parmak izini ASLA aynı kanaldan alma.** Eğer birinin parmak izini sana e-postayla gönderdiği anahtarın *yanında* alıyorsan, MITM hem anahtarı hem parmak izini değiştirmiş olabilir. Parmak izini **bağımsız bir kanaldan** (yüz yüze görüşme, telefonda okutma, kartvizit, doğrulanmış web sitesi) teyit et.

### Keyserver'lar
| Keyserver | Not |
|---|---|
| **keys.openpgp.org** ⭐ | Modern, **e-posta doğrulamalı**, kimlik bilgisini onaysız yaymaz, gizlilik dostu. **Tercih edilen.** |
| **hkps://** (SKS havuzu — eski) | Eski SKS ağı doğrulamasızdı ve "anahtar zehirleme" saldırılarına uğradı; çoğu çöktü. Eski havuzlara güvenme. |

```bash
# Anahtarını yayınla
gpg --keyserver hkps://keys.openpgp.org --send-keys <FINGERPRINT>

# Birinin anahtarını çek
gpg --keyserver hkps://keys.openpgp.org --recv-keys <FINGERPRINT>
```

> ⚠️ **Gizlilik notu:** Anahtarını keyserver'a koymak **e-posta adresini + sosyal grafiğini** (kimler seni imzalamış) herkese açar — bu bir **metadata sızıntısıdır.** Yüksek gizlilik gereken kimliklerde keyserver yerine **WKD** ya da doğrudan elden dağıtım tercih et.

### WKD (Web Key Directory) — modern, temiz dağıtım
WKD, açık anahtarını **kendi alan adının HTTPS sunucusundan** otomatik bulunur kılar. Biri `yasin@ornek.com`'a mesaj yazınca, GPG anahtarı `ornek.com`'un WKD yolundan çeker. Keyserver'ın sosyal-grafik sızıntısı olmadan, *alan sahipliğine* dayalı güven verir. (Kurulum sunucu tarafı yapılandırma ister — alan adın varsa güçlü seçenektir.)

---

<a id="11"></a>
## 11. ⏳ Son Kullanma, Yenileme, İptal

### Son kullanma tarihini uzatma (subkey yenileme)
Subkey'in süresi doluyorsa, ana anahtarla (air-gapped makinede) ileri at:
```bash
gpg --edit-key yasin@ornek.com
gpg> key 1            # ilgili subkey'i seç (numara)
gpg> expire           # yeni süre ver
gpg> save
# Sonra güncel açık anahtarı yeniden yay (keyserver/WKD)
gpg --keyserver hkps://keys.openpgp.org --send-keys <FINGERPRINT>
```

### Tek bir subkey'i iptal etme (sızdıysa)
```bash
gpg --edit-key yasin@ornek.com
gpg> key 1
gpg> revkey           # bu subkey'i iptal et
gpg> save
gpg --keyserver hkps://keys.openpgp.org --send-keys <FINGERPRINT>
```

### Tüm anahtarı iptal etme
Bölüm 9'daki revocation certificate'i içe aktar ve yay.

> 💡 **Yenileme vs yeni anahtar:** Süre uzatmak (expire) anahtarı **değiştirmez**, yalnızca geçerliliği uzatır — WoT imzaların korunur. Ancak **anahtar sızdıysa** süre uzatma yetmez; iptal et ve **yeni anahtar** üret. Süre dolması ≠ güvenlik ihlali; sızma = güvenlik ihlali.

---

<a id="12"></a>
## 12. 🔧 Pratik Kullanım — Şifrele / İmzala / Doğrula / Symmetric

### Asimetrik şifreleme (alıcının açık anahtarıyla)
```bash
# Şifrele (yalnızca alıcı çözebilir) + kendi imzanı da ekle
gpg --encrypt --sign --recipient alici@ornek.com --armor dosya.txt
# çıktı: dosya.txt.asc  (--armor = ASCII; e-posta gövdesine yapışır)

# Çöz
gpg --decrypt dosya.txt.asc > dosya.txt
```

### İmzalama türleri (üçünü karıştırma)
```bash
# 1) Detached signature (AYRI imza dosyası) — exe/sürüm/büyük dosya için İDEAL
gpg --detach-sign --armor surum-1.0.exe
# çıktı: surum-1.0.exe.asc  → dosyaya dokunmaz, ayrı .asc imza
gpg --verify surum-1.0.exe.asc surum-1.0.exe

# 2) Clearsign — metin OKUNUR kalır, altına imza eklenir (e-posta/duyuru için)
gpg --clearsign duyuru.txt
# çıktı: -----BEGIN PGP SIGNED MESSAGE----- ... metin görünür ...

# 3) Gömülü imza (binary, sıkıştırılmış) — metin görünmez, tek dosya
gpg --sign rapor.txt
```

| İmza türü | Dosya bozulur mu? | Metin okunur mu? | En iyi kullanım |
|---|---|---|---|
| **Detached** (`--detach-sign`) | Hayır (ayrı .asc) | Evet | **Yazılım/sürüm/kanıt dosyası bütünlüğü** |
| **Clearsign** (`--clearsign`) | Hayır (sarmalanır) | Evet | E-posta, duyuru, README |
| **Gömülü** (`--sign`) | Evet (sarılır) | Hayır | Tek dosyada imza+veri taşıma |

> 🔥 **Püf:** Bir sürümü/raporu imzalayıp **dosyanın kendisini değiştirmek istemiyorsan → DAİMA detached signature.** Linux dağıtımları, Tor, VeraCrypt hep böyle yapar: `dosya` + `dosya.asc`. İndiren `gpg --verify dosya.asc dosya` ile bütünlüğü kanıtlar.

### Symmetric (parolayla, anahtarsız) şifreleme
Hiç anahtar çifti olmadan, yalnızca **parolayla** şifreleme. Kendi yedeğin ya da parolayı paylaşabildiğin biri için:
```bash
gpg --symmetric --cipher-algo AES256 gizli.zip
# çıktı: gizli.zip.gpg  → çözmek için yalnızca parola gerekir
gpg --decrypt gizli.zip.gpg > gizli.zip
```
> 💡 Symmetric mod, "anahtar altyapısı yok ama bu dosyayı parolayla kilitleyeyim" senaryosu için pratiktir. **Parolayı güvenli kanaldan paylaş** (asla aynı mesajda gönderme).

### Sıkça
```bash
gpg --list-keys                 # açık anahtarlar
gpg --list-secret-keys          # özel anahtarlar
gpg --export --armor KIMLIK     # açık anahtarı dışa aktar
gpg --import dosya.asc          # anahtar içe aktar
```

---

<a id="13"></a>
## 13. 🔑 YubiKey / Akıllı Kart ile Donanımda Anahtar

**En güçlü pratik duruş.** Özel anahtar bir **donanım belirtecine** (YubiKey, Nitrokey, OpenPGP smartcard) yazılır ve **bir daha asla** belirteci terk etmez. İmzalama/çözme **kartın içinde** olur; bilgisayar yalnızca "şunu imzala" der, anahtarı hiç görmez.

### Neden bu kadar güçlü?
- **Keylogger/RAT bile anahtarı çalamaz** — anahtar diskte/RAM'de yok, kartın güvenli yongasında.
- Her işlem için **fiziksel dokunuş** (touch-to-sign) zorunlu kılınabilir → kötücül yazılım sen yokken sessizce imza atamaz.
- Kartın **PIN denemesi sınırlıdır** (genelde 3) → brute-force imkânsız; karta el konsa bile PIN olmadan işe yaramaz.

### Akış (kavramsal)
```
Bilgisayar (güvensiz olabilir)          YubiKey (güvenli yonga)
   "Şu hash'i imzala" ───────────────►  [özel anahtar İÇERİDE]
                                         PIN + dokunuş gerekir
   imzalı sonuç      ◄───────────────  imzayı üret, anahtarı asla verme
```

### Kurulum (özet — kendi kartının dokümanından teyit et)
```bash
# Kartı tanı
gpg --card-status

# Kart yönetimi (PIN değiştir, sahip bilgisi, vb.)
gpg --card-edit
gpg/card> admin
gpg/card> passwd        # PIN ve Admin PIN'i DEĞİŞTİR (varsayılan 123456/12345678!)
```

İki yaygın yöntem:
1. **Kartta üret (en güvenli):** Anahtar **kartın içinde** üretilir, dışarı hiç çıkmaz → yedeklenemez de (kart bozulursa anahtar gider; bu yüzden **iki kart** ya da air-gapped yedek mantığı).
2. **keytocard (yaygın):** Air-gapped makinede ana anahtar + subkey üret, **subkey'leri `keytocard` ile karta taşı**, ana anahtarı/yedeği güvenli sakla. Kart bozulursa yedekten yeni karta yüklersin.

```bash
# Air-gapped makinede: subkey'i karta taşı (edit-key içinde)
gpg --edit-key yasin@ornek.com
gpg> key 1          # taşınacak subkey'i seç
gpg> keytocard      # karta yaz (gizli parça karttan sonra makineden silinebilir)
gpg> save
```

> 🔥 **Püf — ideal kombinasyon:** *Air-gapped ana anahtar (Bölüm 8) + subkey'ler YubiKey'de (keytocard) + basılı revocation cert (Bölüm 9) + paperkey yedeği (Bölüm 17).* Bu dörtlü, gazeteci/CTI analisti/yüksek-değerli hedef için sahada bilinen **en sağlam GPG duruşudur.** Günlük makinende özel anahtar **hiç bulunmaz** — yalnızca kart takılıyken, PIN+dokunuşla iş görür.

> ⚠️ **Yedek şart:** "Kartta üret" yöntemi en güvenli ama **yedeği yoktur** — kart kaybolursa/bozulursa anahtar gider. Kritik kimlikte ya **iki özdeş kart** hazırla ya **keytocard + air-gapped yedek** kullan.

---

<a id="14"></a>
## 14. 🧩 gpg-agent, SSH Auth, Parola Önbelleği

### gpg-agent
Modern GnuPG'de tüm özel anahtar işlemlerini **gpg-agent** yürütür. Parolanı bir kez sorar, yapılandırılmış süre boyunca **önbellekte** tutar (her işlemde tekrar sormaz).

```ini
# ~/.gnupg/gpg-agent.conf
default-cache-ttl 600        # önbellek 10 dk (saniye)
max-cache-ttl 7200           # en fazla 2 saat
# pinentry-program ...        # parola sorma arayüzü (platforma göre)
```
```bash
# Değişiklikten sonra agent'ı yeniden yükle
gpg-connect-agent reloadagent /bye
```

> 🔥 **Püf — önbellek tehdidi:** Uzun `cache-ttl`, makineden kalkınca birinin senin parolan olmadan imza/çözme yapabileceği bir **pencere** açar (parola RAM'de tutulu). Hassas ortamda **kısa TTL** kullan; ekran kilidi disiplinini koru. Tehlike anında önbelleği temizlemek için agent'ı yeniden başlat.

### SSH authentication GPG ile (auth subkey)
[A] yetenekli subkey'ini SSH anahtarın olarak kullanabilirsin — özellikle **YubiKey'de** ise, SSH özel anahtarın **donanımda** durur:
```ini
# ~/.gnupg/gpg-agent.conf
enable-ssh-support
```
```bash
# gpg-agent'ı SSH agent yap (kabuk profiline ekle)
export SSH_AUTH_SOCK=$(gpgconf --list-dirs agent-ssh-socket)
# Public SSH anahtarını al (sunucuya koymak için)
gpg --export-ssh-key yasin@ornek.com
```
> 💡 Böylece SSH ile sunucuya girerken özel anahtar diskte düz dosya olarak durmaz; YubiKey'deyse **dokunuşla** kimlik doğrularsın — sunucu erişimi için güçlü bir model.

---

<a id="15"></a>
## 15. 🔥 PÜF NOKTALARI — Piyasada Bulamayacakların

Çoğu rehberin atladığı, gerçek dünyada güvenliği belirleyen detaylar.

### 15.1 Air-gapped ana anahtar + günlük subkey (özet ve neden)
Bölüm 7-8'in özü: **kimliğin (ana anahtar) internete bağlı makineyi asla görmesin.** Bu tek karar, "cihaz çalındı → kimliğim öldü"yü "cihaz çalındı → bir subkey iptal ettim"e çevirir. **GPG ustalığının tek en önemli pratiğidir.**

### 15.2 Revocation cert'i önceden üret + BASILI sakla
Bölüm 9. Anahtarı yaymadan üret, *kağıda* yazdır, kasaya koy. Parolanı unutsan/anahtarı kaybetsen bile anahtarını "ölü" ilan edebilmenin tek yolu.

### 15.3 paperkey — anahtarı KAĞIDA yedekle
Dijital yedekler bozulur, silinir, şifrelenir (ransomware). **paperkey**, özel anahtarın *yalnızca gizli* kısmını minimal bir metne indirger → **yazdırırsın**, fiziksel kasada saklarsın. Mıknatıs/disk arızası/ransomware fiziksel kağıdı etkilemez.
```bash
# Özel anahtarı kağıda dökülecek metne çevir
gpg --export-secret-keys yasin@ornek.com | paperkey --output anahtar-kagit.txt
# (yazdır, kasaya koy)

# Geri yükleme: açık anahtar + kağıttaki gizli kısım birleştirilir
paperkey --pubring public.asc --secrets anahtar-kagit.txt | gpg --import
```
> 🔥 **Püf:** paperkey, açık anahtarı yeniden yazmaz (zaten elinde) — yalnızca **gizli kısmı** kağıda döker, böylece kağıt daha küçük olur. Geri yükleme için **açık anahtarın da bir kopyası** gerekir; ikisini birlikte planla.

### 15.4 gpg.conf sertleştirme — zayıf varsayılanları kapat
Eski sürümlerin varsayılanları zayıf algoritmalara izin verebilir. Tercih sıralamasını **SHA-512 / AES-256** lehine sabitle (tam içerik Bölüm 16). Karşı tarafın eski yazılımıyla *kırılmasın* diye yalnızca tercih sırası verirsin (zorlama değil); ama kendi imzaların güçlü digest kullanır.

### 15.5 Metadata sızıntısı — `--throw-keyids` ile alıcıyı gizle
Normalde şifreli bir paket, **kimin açık anahtarına** şifrelendiğini (alıcı Key ID'leri) açıkça taşır → "bu mesaj X kişisine" ifşa olur. `--hidden-recipient` / `--throw-keyids` bu Key ID'leri siler:
```bash
gpg --encrypt --hidden-recipient alici@ornek.com --throw-keyids dosya.txt
```
> ⚠️ **Bedeli:** Alıcı, hangi anahtarla çözeceğini bilemez → **tüm özel anahtarlarını sırayla dener** (yavaşlar). Yine de "kime şifrelendiği" metadata'sını gizlemek istediğin hassas senaryolarda değerlidir. **Konu/From/To e-posta başlıklarını GPG yine de gizlemez** (Bölüm 2) — bu yalnızca paket-içi alıcı kimliğini saklar.

### 15.6 "Encrypt to self" tuzağı
GnuPG genelde, şifrelediğin her şeyi **kendi anahtarına da** şifreler (`encrypt-to-self`) ki sonradan kendi gönderdiğini okuyabilesin.
- **Yararı:** Gönderdiğin şifreli e-postayı "Sent" klasöründe okuyabilirsin.
- **Tuzağı:** Bu, şifreli paketin **senin** anahtarınla da çözülebileceği anlamına gelir → senin anahtarın ele geçerse o mesajlar da açılır; ayrıca alıcı listesinde senin Key ID'in görünür (metadata). Yüksek-gizlilik tek-yönlü teslimde `--no-encrypt-to-self` düşün.

### 15.7 Tarih/saat güvenilmezliği
İmza üstündeki **zaman damgası, imzalayanın makine saatinden** gelir → kötü niyetli biri saatini değiştirip **geçmiş/gelecek tarihli** imza üretebilir. İmza zamanı *kanıt* değil, *iddiadır.* Kanıt zinciri/CTI raporu için bağımsız bir **zaman damgası otoritesi** (RFC 3161 TSA) ya da blok zinciri/şahit kullan. **GPG imza tarihine tek başına adli delil gibi güvenme.**

### 15.8 Sürüm uyumluluğu
ED25519, modern AEAD modları ve bazı paket formatları **eski GnuPG/PGP** ile çözülemeyebilir. Karşı tarafın eski yazılımı varsa: RSA-4096 + SHA-512 gibi geniş uyumlu seçim yap. Yeni bir özelliği kullanmadan **alıcının sürümünü** düşün; aksi halde "şifreledim ama açamıyor" yaşanır.

### 15.9 Clearsign vs detached — doğru aracı seç
Bölüm 12 tablosu: dosyayı **değiştirmeden** bütünlük → detached (`.asc` ayrı). Metnin **okunur** kalması + imza → clearsign. İkisini karıştırmak yaygın hatadır (örn. bir exe'yi clearsign'lamaya çalışmak anlamsızdır — binary'de detached şart).

### 15.10 Anahtar zehirleme (keyserver flooding) — eski SKS riski
Eski SKS keyserver'ları, bir anahtara **sınırsız sahte imza** eklenmesine izin veriyordu → anahtar o kadar şişiyordu ki içe aktaran GnuPG kilitleniyordu (DoS). **keys.openpgp.org** bunu çözer (imzaları doğrulamadan yaymaz). Eski havuzlardan körlemesine `--recv-keys` yapma.

### 15.11 Forward secrecy yok → hassas anlık iletişimde Signal
Tekrar vurgu (Bölüm 2): GPG'de tek bir uzun-ömürlü anahtar tüm geçmişi açar. Gerçek zamanlı, *yarın anahtarım çalınırsa bugünkü mesajım da gitsin istemiyorum* dediğin sohbetler için **Signal** (X3DH + Double Ratchet, her mesaja taze anahtar) doğru araçtır. GPG'yi **dosya/sürüm/e-posta/kanıt** için, Signal'i **anlık hassas konuşma** için kullan.

### 15.12 ASCII armor (`--armor`) ne zaman?
`--armor` çıktıyı base64 metne çevirir → e-posta gövdesine/sohbete yapışır, kopyalanır. Büyük binary dosyada gereksiz **%33 şişme** yapar → büyük dosyada armor'suz (binary `.gpg`) tut, yalnızca metin kanalında taşıyacaksan armor kullan.

### 15.13 Şifreleme ≠ imzalama — ikisini birlikte iste
Yalnızca `--encrypt` gizlilik verir ama **kaynağı kanıtlamaz** (saldırgan da alıcıya şifreleyebilir). Hem gizlilik hem kaynak için **`--encrypt --sign` birlikte.** "Şifreli geldi, demek ki ondan" yanılgısına düşme — imza yoksa gönderici doğrulanmaz.

---

<a id="16"></a>
## 16. 🛡️ gpg.conf Sertleştirme (Hardening)

`~/.gnupg/gpg.conf` dosyasına eklenecek, güçlü-varsayılan örneği. **Bunlar tercih sıralamasıdır** — karşı tarafın desteklediği en güçlüyü seçtirir, desteklemiyorsa kırılmaz.

```ini
# --- Kimlik gösterimi ---
keyid-format 0xlong            # Key ID'leri uzun ve 0x önekli göster
with-fingerprint               # parmak izini her zaman göster

# --- Algoritma tercih sıralaması (güçlüden zayıfa) ---
personal-cipher-preferences AES256 AES192 AES
personal-digest-preferences SHA512 SHA384 SHA256
personal-compress-preferences ZLIB BZIP2 ZIP Uncompressed
cert-digest-algo SHA512        # anahtar imzalarken (sertifika) SHA-512 kullan
default-preference-list SHA512 SHA384 SHA256 AES256 AES192 AES ZLIB BZIP2 ZIP Uncompressed

# --- Zayıf SHA-1'i kısıtla ---
weak-digest SHA1               # SHA-1'i zayıf işaretle (sürüm destekliyorsa)

# --- Genel sertleştirme ---
no-emit-version                # GnuPG sürümünü çıktıya yazma (metadata azalt)
no-comments                    # yorum satırı ekleme
require-cross-certification    # imza subkey'leri için çapraz sertifika zorunlu
no-symkey-cache                # simetrik parolayı agent'ta önbellekleme (paranoyak)
```

> ⚠️ **Teyit et:** Bazı seçenekler (örn. `weak-digest`, AEAD ile ilgili anahtarlar) **GnuPG sürümüne göre** isim/destek değiştirir. Kendi `gpg --version`'ında geçerliliğini doğrula; bilinmeyen seçenek hata verir. Yukarıdakiler 2.2/2.4 hattında yaygın olarak geçerlidir, ama körlemesine kopyalamadan önce her satırı sürümünde test et.

> 🔥 **Püf:** En kritik satır **`cert-digest-algo SHA512`** ve **cipher/digest tercih listeleridir** — bunlar senin ürettiğin imza ve şifrelerin güçlü algoritma kullanmasını sağlar. Eski varsayılan SHA-1 sertifika digest'i, modern güvenlikte kabul edilemez.

---

<a id="17"></a>
## 17. 💾 Yedekleme & Kurtarma

Özel anahtarını kaybedersen **şifreli her şey sonsuza dek gider** ve kimliğini yeniden kuramazsın. Yedek = hayat. Ama yedek de korunmalı.

### Neyi yedekle
1. **Tam gizli anahtar (ana + subkey'ler):**
   ```bash
   gpg --export-secret-keys --armor yasin@ornek.com > secret-full.asc
   ```
2. **Açık anahtar:**
   ```bash
   gpg --export --armor yasin@ornek.com > public.asc
   ```
3. **Owner trust (WoT güven atamaların):**
   ```bash
   gpg --export-ownertrust > ownertrust.txt
   ```
4. **Revocation certificate** (Bölüm 9) — ayrı + basılı.
5. **paperkey çıktısı** (Bölüm 15.3) — fiziksel kasada.

### Nereye / nasıl
- Gizli anahtar yedeğini **şifreli bir ortama** koy (VeraCrypt birimi — bkz. `VERACRYPT_USTALIK_REHBERI.md` — ya da donanım-şifreli USB).
- **3-2-1:** 3 kopya, 2 farklı ortam, 1 saha dışı. *Ama her kopya korunmalı — yedek de bir saldırı yüzeyidir.*
- **paperkey + basılı revocation** = dijital felakete (ransomware, disk ölümü) karşı son sigorta.

### Geri yükleme
```bash
gpg --import secret-full.asc          # gizli anahtarları geri al
gpg --import public.asc               # açık anahtarı geri al
gpg --import-ownertrust ownertrust.txt # güven atamalarını geri al
```

> 🔥 **Püf — yedeği TEST et:** Yedeğini **temiz/ayrı bir makinede** içe aktarıp bir test dosyasını çözmeyi dene. Açılmayan yedek, yedek değildir. Bunu *anahtarı kaybetmeden önce* yap.

---

<a id="18"></a>
## 18. 🏰 Kanije Kalesi ile Birlikte Kullanım

Bu repo (Kanije Kalesi), CTI/güvenlik odaklı bir muhafız ve araç setidir. GPG, onun **bütünlük ve gizli iletişim** katmanını sağlar — felsefe birebir örtüşür: *"matematiksel kanıt + kapıyı kilitli tutma disiplini."*

| İhtiyaç | GPG rolü | Kanije Kalesi ile bağ |
|---|---|---|
| **CTI raporu / kanıt bütünlüğü** | Raporu/IOC dosyasını **imzala** (detached `.asc`) → değiştirilmediği kanıtlanır, inkâr edilemez | Üretilen rapor çıktısına imza ekle → zincir kanıtı |
| **Kaynaklarla şifreli iletişim** | Muhbir/kaynak ile **uçtan uca** şifreli e-posta (açık anahtar değişimi) | Hassas kaynak yazışması; metadata için Signal'i tamamlayıcı kullan |
| **exe / sürüm imzalama** | Dağıttığın binary'yi **detached-sign** et → indiren truva atı olmadığını doğrular | Repo sürümlerini imzala; `VERACRYPT_USTALIK_REHBERI.md`'deki "imza doğrulama" felsefesinin aynısı |
| **Yapılandırma/log bütünlüğü** | Kritik config/log'u imzala → sonradan kurcalanma tespit edilir | Adli inceleme öncesi delil bütünlüğü |
| **Anahtar yedeği güvenliği** | Özel anahtar yedeğini VeraCrypt hidden volume'da tut | Kanije `/imha` ile gerekirse anahtar yedeğini secure-wipe |

### 🔥 Önerilen entegrasyon deseni
1. **Repo sürümlerini/raporlarını detached-sign et** → her yayında `dosya` + `dosya.asc`. Kullanıcı `gpg --verify` ile bütünlüğü kanıtlar. Bu, projenin "imza bütünlüğü" duruşunu tamamlar.
2. **Air-gapped ana anahtar + YubiKey subkey** ile imzala → imza atan makine ele geçse bile **kimliğin (ana anahtar) sağlam.**
3. **Kaynak iletişimi:** İçerik için GPG (e-posta, dosya), *anlık + forward-secrecy gereken* konuşma için Signal — ikisi katmanlı.
4. **Anahtar yedeği:** Özel anahtar + paperkey + revocation cert'i **VeraCrypt hidden volume**'da sakla; en kritik senaryoda Kanije'nin secure-wipe'ı ile yedeği yok et (kalan kopyalar air-gapped/basılı).
5. **Bütünlük zinciri:** CTI çıktılarını imzalayıp parmak izini bilinen bir yere sabitle → "bu rapor gerçekten bizden ve değişmedi" matematiksel olarak kanıtlanır.

> 🧠 **Felsefe örtüşmesi:** VeraCrypt **durağan veriyi** (gizlilik), GPG **bütünlüğü + kaynağı** (kim, değişti mi), Kanije Kalesi **olay anını** (fiziksel tehdit, imha) yönetir. Üçü birlikte: *gizli + doğrulanabilir + savunulabilir* bir duruş.

---

<a id="19"></a>
## 19. ☠️ Yaygın Ölümcül Hatalar

1. **Özel anahtarı yedeklememek** → disk ölünce kimlik + tüm şifreli veri sonsuza dek gider. (En sık felaket.)
2. **Revocation certificate'i önceden üretmemek** → anahtar/parola kaybında anahtarı "ölü" ilan edememek; insanlar erişemediğin anahtara şifrelemeye devam eder.
3. **Expired (süresi dolmuş) anahtara güvenmek / takip etmemek** → imzaların bir gün "geçersiz" görünür, şifre alamazsın; son kullanmayı takvime koy.
4. **Tek "her şeyi yapan" anahtar** (subkey mimarisi yok) → cihaz çalınınca kimlik + tüm WoT imzaları çöp.
5. **Ana anahtarı günlük/çevrimiçi makinede tutmak** → sızdığında kimliğin tamamen ele geçer (air-gapped tut).
6. **Parmak izini yalnızca Key ID'den ya da aynı kanaldan doğrulamak** → MITM ile yanlış anahtara şifrelemek.
7. **Yalnızca `--encrypt`, `--sign` olmadan** → gizlilik var ama kaynak doğrulanmaz; sahte gönderici mümkün.
8. **Passphrase'i unutmak / zayıf seçmek** → özel anahtar kullanılamaz hale gelir; zayıfsa dosya çalınınca açılır.
9. **Konu satırının/metadata'nın şifrelendiğini sanmak** → GPG From/To/Subject/tarihi gizlemez; hassas bilgiyi konuya yazma.
10. **Forward secrecy olduğunu sanmak** → anahtar bir gün sızarsa geçmiş tüm mesajlar açılır; hassas anlık iletişimde Signal kullan.
11. **İmza tarihini adli kanıt sanmak** → zaman damgası imzalayanın saatinden gelir, sahtelenebilir; TSA kullan.
12. **Eski SKS keyserver'dan körlemesine anahtar çekmek** → zehirlenmiş anahtar GnuPG'yi kilitleyebilir; keys.openpgp.org kullan.
13. **YubiKey'de "kartta üret" yapıp yedek almamak** → kart bozulunca anahtar gider; iki kart ya da keytocard+yedek.
14. **Yedeği test etmemek** → felaket anında "açılmayan yedek" çıkar.

---

<a id="20"></a>
## 20. ✅ Hızlı Referans & Operasyonel Kontrol Listesi

### En sık komutlar

| Amaç | Komut |
|---|---|
| Anahtar üret (rehberli) | `gpg --full-generate-key` |
| Anahtar üret (uzman) | `gpg --expert --full-generate-key` |
| Parmak izini göster | `gpg --fingerprint KIMLIK` |
| Açık anahtarı dışa aktar | `gpg --export --armor KIMLIK > public.asc` |
| Anahtar içe aktar | `gpg --import dosya.asc` |
| Şifrele + imzala | `gpg --encrypt --sign -r alici@x.com --armor dosya` |
| Çöz | `gpg --decrypt dosya.asc > dosya` |
| Detached imza (sürüm/dosya) | `gpg --detach-sign --armor dosya` |
| Detached doğrula | `gpg --verify dosya.asc dosya` |
| Clearsign (metin) | `gpg --clearsign dosya.txt` |
| Symmetric (parolayla) | `gpg --symmetric --cipher-algo AES256 dosya` |
| Revocation cert üret | `gpg --output revoke.asc --gen-revoke KIMLIK` |
| Subkey ekle/yönet | `gpg --expert --edit-key KIMLIK` → `addkey` / `keytocard` / `expire` / `revkey` |
| Kartı tanı (YubiKey) | `gpg --card-status` |
| Tam gizli yedek | `gpg --export-secret-keys --armor KIMLIK > secret.asc` |
| paperkey yedeği | `gpg --export-secret-keys KIMLIK \| paperkey --output kart.txt` |
| Anahtarı yay | `gpg --keyserver hkps://keys.openpgp.org --send-keys FINGERPRINT` |
| Anahtar çek | `gpg --keyserver hkps://keys.openpgp.org --recv-keys FINGERPRINT` |
| Gizli anahtar stub kontrolü | `gpg --list-secret-keys` (ana anahtarda `sec#` görmeli) |

### Anahtar üretirken
- [ ] **GnuPG 2.2+ (tercihen 2.4+)** ve imzası doğrulanmış kurulum
- [ ] Algoritma: **ED25519/Curve25519** (modern) ya da **RSA-4096** (uyum)
- [ ] **Güçlü passphrase** (diceware, benzersiz)
- [ ] **Ana anahtar = yalnızca [C]**; ayrı **[S]/[E]/[A] subkey'leri**
- [ ] Subkey'lere **son kullanma** (1-2 yıl), takvime kaydedildi
- [ ] **Revocation cert üretildi** + ayrı + **basılı** saklandı
- [ ] **Tam gizli anahtar yedeği** + **paperkey** + **ownertrust** alındı, **şifreli** saklandı
- [ ] Yedek **ayrı makinede test** edildi

### Mimari (ileri)
- [ ] **Ana anahtar air-gapped** (günlük makinede `sec#` stub görünüyor)
- [ ] Subkey'ler **YubiKey/akıllı kartta** (keytocard) ya da en azından subkey-only
- [ ] **gpg.conf sertleştirildi** (SHA-512 cert-digest, AES-256 tercihi)
- [ ] WKD ya da **keys.openpgp.org** ile dağıtım (eski SKS yok)

### Her kullanımda
- [ ] Şifrelerken **`--encrypt --sign` birlikte** (kaynak da kanıtlansın)
- [ ] Alıcının **parmak izi bağımsız kanaldan doğrulanmış**
- [ ] Sürüm/dosya bütünlüğü → **detached signature**
- [ ] Hassas metadata (kim-kime) gizlenmeli → `--hidden-recipient --throw-keyids` (bedelini bilerek)
- [ ] Anlık + forward-secrecy gereken konuşma → **GPG değil Signal**

### Periyodik
- [ ] Subkey **son kullanma** yaklaştı mı → air-gapped makinede uzat
- [ ] Anahtar sızma şüphesi → **iptal et + yeni anahtar**
- [ ] Yedek + revocation cert **hâlâ erişilebilir ve güvende** mi
- [ ] Parola değişti/anahtar döndü → **eski yedekleri** güvenli imha

---

> 🏰 **Kapanış:** GPG bir ürün değil, bir **disiplindir.** En güçlü ED25519 bile, ana anahtarını çevrimiçi bir makinede tuttuğun, revocation cert'i hiç üretmediğin ya da özel anahtarını yedeklemediğin gün çaresizdir. GPG sana **kimliğin matematiksel mührünü** ve **içeriğin kilidini** verir; **kökü (ana anahtarı) çevrimdışı, yedeği basılı, subkey'i donanımda tutmak senin işin.** Forward secrecy istediğin anlık iletişimde ise alet çantandan **Signal'i** çıkar — doğru aracı doğru tehdide.
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. İlgili: `VERACRYPT_USTALIK_REHBERI.md`, `SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md`, `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`, `DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md`.*
