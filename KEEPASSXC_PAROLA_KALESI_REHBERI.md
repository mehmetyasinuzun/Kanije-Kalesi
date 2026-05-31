# 🗝️ KEEPASSXC — TAM USTALIK REHBERİ
## Çevrimdışı Parola Kasasından Üçlü Donanım Kilidine, Püf Noktalarıyla Uçtan Uca

> **Amaç:** KeePassXC'yi "indir-parola-koy" seviyesinden çıkarıp, bir veri ihlali, bir hırsızlık ya da hedefli bir saldırgan karşısında **gerçekten ayakta kalacak** şekilde kullanmayı öğretmek. Bu rehber yalnızca *nasıl*'ı değil, **neden** ve **hangi durumda işe yaramaz**'ı da anlatır. Forum cevaplarında bulamayacağın Argon2 ayar matematiği, üçlü kilit (parola + key file + YubiKey), "açık veritabanı RAM'de çözülü" gerçeği, pano sızıntısı ve adli bilişim karşıtı detaylar burada.

> ⚠️ **Önce bunu oku:** Parola yöneticisi, *yanlış kullanıldığında* sana **tek noktadan çöküş (single point of failure)** yaratır — tüm yumurtalar tek sepette. Doğru kullanıldığında ise hayatındaki en güçlü güvenlik kararıdır. Özellikle **Tehdit Modeli**, **Üçlü Kilit** ve **Yedekleme** bölümlerini atlama. Tek bir key file kaybı ya da tek bir kopya yedek, yıllarca biriken her şeyi **bir anda** yok edebilir.

---

## 📑 İÇİNDEKİLER

1. [KeePassXC Nedir, Neden? (Bitwarden/Bulut Kıyası)](#1)
2. [Tehdit Modeli — Neyi Korur, Neyi KORUMAZ](#2)
3. [Kurulum + İmza/Hash Doğrulama (Atlama!)](#3)
4. [Temel Kavramlar — .kdbx / Master Key / KDF](#4)
5. [Veritabanı Oluşturma & 🔥 ÜÇLÜ Kilit](#5)
6. [KDF Derinlemesine — Argon2id vs AES-KDF](#6)
7. [Şifreleme Algoritması — AES-256 vs ChaCha20](#7)
8. [Key File Stratejisi (Derinlemesine)](#8)
9. [YubiKey / Challenge-Response Derinlemesine](#9)
10. [TOTP (2FA) Saklama — Tek Sepet Tartışması](#10)
11. [Auto-Type & Pano Güvenliği](#11)
12. [Tarayıcı Entegrasyonu (KeePassXC-Browser)](#12)
13. [SSH Agent Entegrasyonu](#13)
14. [Parola/Passphrase Üretici (Diceware & Entropi)](#14)
15. [Yedekleme & Senkronizasyon (Syncthing/Bulut)](#15)
16. [🔥 PÜF NOKTALARI — Piyasada Bulamayacakların](#16)
17. [Yaygın Ölümcül Hatalar](#17)
18. [🏰 Kanije Kalesi ile Birlikte Kullanım](#18)
19. [Hızlı Referans & Operasyonel Kontrol Listesi](#19)

---

<a id="1"></a>
## 1. 🧭 KeePassXC Nedir, Neden?

KeePassXC, klasik **KeePass**'in (Windows odaklı, .NET) topluluk tarafından sürdürülen, çapraz platform (Windows/macOS/Linux), modern C++/Qt fork'udur. **Açık kaynaktır**, aktif olarak geliştirilir ve denetlenir. Verini **tek bir şifreli dosyada** (`.kdbx`) saklar — bu dosya **senin diskinde** durur, **buluta çıkmaz**, bir şirketin sunucusunda barınmaz.

Ne yapar: Tüm parolalarını, TOTP kodlarını, SSH anahtarlarını, notlarını **tek bir ana parola** (master password) ile kilitlenen yerel bir kasada toplar. Veri diskte **hiçbir zaman** çözülmüş halde durmaz; yalnızca veritabanı **açıkken** ve **yalnızca RAM'de** çözülür.

### Bitwarden / 1Password / LastPass (Bulut) ile Dürüst Kıyas

| Özellik | KeePassXC | Bitwarden | 1Password | LastPass |
|---|---|---|---|---|
| Açık kaynak | ✅ | ✅ (sunucu kısmen) | ❌ | ❌ |
| **Veri nerede** | **Senin diskinde** | Bulut (şirket sunucusu) | Bulut | Bulut |
| Buluta zorunlu çıkar mı | ❌ (sen seçersen) | ✅ | ✅ | ✅ |
| Çevrimdışı tam çalışır | ✅ | Kısmen (cache) | Kısmen | Kısmen |
| Şirkete güven gerekir | ❌ | ✅ | ✅ | ✅ |
| Sunucu ihlalinde risk | **Yok** (senin diskin) | Şifreli blob sızar | Şifreli blob sızar | **Sızdı (2022)** |
| Senkronizasyon kolaylığı | Manuel (sen kurarsın) | Otomatik | Otomatik | Otomatik |
| Donanım anahtarı (kasa kilidi) | ✅ (YubiKey CR) | ✅ (2FA, kasayı şifrelemez*) | ✅ | ✅ |
| Ücret | Ücretsiz | Ücretsiz/Premium | Ücretli | Ücretli |

> ℹ️ \*Çoğu bulut yöneticisinde YubiKey **giriş 2FA'sıdır** (hesaba erişimi engeller), kasanın **şifreleme anahtarının parçası değildir**. KeePassXC'de Challenge-Response **doğrudan veritabanını şifreleyen anahtara karışır** — bu kavramsal olarak daha güçlüdür (Bölüm 9).

**CTI / hassas veri perspektifi:** Bulut yöneticileri **kullanım kolaylığı** sunar ama **güven modelini** dışarı verir: verin (şifreli de olsa) bir şirketin altyapısında durur. LastPass 2022 ihlalinde **şifreli kasa yedekleri çalındı**; zayıf ana parolası olan kullanıcılar offline brute-force'a maruz kaldı. KeePassXC'de **böyle bir saldırı yüzeyi yoktur** — `.kdbx` dosyan senin elinde, bir sunucu ihlalinde sızacak merkezî bir havuz **yok**. Bunun bedeli: senkronizasyon, yedek ve cihazlar arası erişimi **sen** yönetirsin (Bölüm 15). Hassas araştırma, gazetecilik ve kaynak koruma için **çevrimdışı KeePassXC, bulut yöneticisine göre daha sağlam bir tehdit duruşudur** — özellikle "kimseye güvenmek istemiyorum" senaryosunda.

> 💡 **Katmanlı yaklaşım:** İkisi birlikte de kullanılabilir — günlük/düşük riskli parolalar Bitwarden'da (kolaylık), kritik sırlar (BitLocker/BIOS/kurtarma/kripto seed) çevrimdışı KeePassXC'de (maksimum güven). Bölüm 16'da "çoklu veritabanı" desenine bak.

---

<a id="2"></a>
## 2. 🎯 Tehdit Modeli — Neyi Korur, Neyi KORUMAZ

Kasanı kurmadan önce **kime karşı** korunduğunu netleştir. Aksi halde yanlış güvenlik hissi yaşarsın. Parola yöneticisinin en tehlikeli yanı, "her şey güvende" sanıp **tek sepetin** zafiyetlerini görmemektir.

### ✅ KeePassXC'nin KORUDUĞU senaryolar
- **Veri ihlali / sızıntı:** Bir web sitesi hacklenip parolaları sızsa bile, her site için **ayrı, rastgele** parola kullandığından zincirleme ele geçirme (credential stuffing) olmaz.
- **Bulut güveni yok:** Hiçbir şirket senin kasanı tutmaz → sunucu ihlali seni etkilemez.
- **Disk/dosya çalınması (veritabanı kapalı):** `.kdbx` başka makineye kopyalansa → güçlü ana parola + KDF varsa offline brute-force pratikte imkânsızdır.
- **Çevrimdışı çalışma:** İnternet olmadan, hava boşluğu (air-gap) ortamında bile tam işlevseldir.
- **Zayıf insan hafızası:** Tek güçlü ana parolayı ezberle; yüzlerce karmaşık parolayı KeePassXC hatırlasın → parola tekrarını (en yaygın zafiyet) bitirir.

### ❌ KeePassXC'nin KORUMADIĞI senaryolar
- **Keylogger / tuş kaydedici:** Ana parolanı yazarken yakalar → tüm kasa düşer. Şifreleme burada **çaresizdir.** (Auto-Type kısmi yardım eder — Bölüm 11.)
- **Pano (clipboard) sızıntısı:** Parolayı "Kopyala" deyince panoya düşer; pano izleyen kötücül yazılım / başka uygulama / pano geçmişi (Win+V) okur.
- **Açık veritabanı = RAM'de çözülü:** Kasa açıkken parolalar **belleğe çözülmüş** durur. O anda makineye erişen kötücül yazılım/uzak erişim okuyabilir; cold-boot ile RAM çekilebilir (Bölüm 16).
- **Uç nokta uzlaşması (endpoint compromise):** Makinen RAT/rootkit ile ele geçirildiyse, ekran görüntüsü/bellek/disk her şeye erişilir → parola yöneticisi de çaresizdir.
- **Ekran görüntüsü / omuz sörfü:** Parolayı ekranda "göster" deyip açtığında çevrendeki gözler/kayıt yakalar.
- **Zayıf ana parola:** KDF ne kadar güçlü olursa olsun, ana parola "123456" ise kasa kâğıttan kaledir.
- **Phishing (kullanıcının kendisi):** KeePassXC-Browser doğru alan adını eşler (phishing'e kısmi koruma) ama parolayı elle kopyalarsan sahte siteye yapıştırabilirsin.

> 🧠 **Altın kural:** KeePassXC **"durağan veriyi" (data-at-rest)** ve **parola hijyenini** korur. **"Kullanımdaki veriyi" (data-in-use — kasa açıkken)** ve **uç noktanı** korumak senin operasyonel disiplinine bağlıdır (auto-lock, pano temizleme, temiz/sertleştirilmiş makine, anti-keylogger). Kasa **kullanmadığın her an KİLİTLİ** olmalı.

---

<a id="3"></a>
## 3. 📥 Kurulum + İmza/Hash Doğrulama

**Asla** rastgele bir "indir.com" aynasından indirme. Yalnızca resmi **keepassxc.org** ya da resmi paket depoları (Microsoft Store, Flatpak `org.keepassxc.KeePassXC`, dağıtım deposu). Değiştirilmiş bir yükleyici, *kurduğun anda* tüm kasanı sızdırır.

### İmza/Hash doğrulama (ATLAMA — kurulumun en kritik adımı)

KeePassXC sürümleri **PGP ile imzalanır** (anahtar parmak izini resmi siteden teyit et). İndirdiğin dosyanın yanında bir `.sig` ve `.DIGEST` dosyası bulunur.

```bash
# 1. KeePassXC imzalama anahtarını içe aktar (parmak izini keepassxc.org/verifying-signatures üzerinden TEYİT ET)
gpg --recv-keys CFB4C2166397D0D2DD8C46D9D7B62EBD1F71E4A1   # teyit et: sürümle değişebilir

# 2. PGP imzasını doğrula (Windows yükleyicisi örneği)
gpg --verify "KeePassXC-2.7.x-Win64.msi.sig" "KeePassXC-2.7.x-Win64.msi"
#   → "Good signature from ... KeePassXC Release ..." ve doğru parmak izi görmelisin

# 3. SHA-256 hash karşılaştır
#   Windows (PowerShell):
Get-FileHash .\KeePassXC-2.7.x-Win64.msi -Algorithm SHA256
#   Çıkan değeri sitedeki .DIGEST içeriğiyle karakter karakter karşılaştır
```

- **Windows:** Yükleyici **Authenticode imzasını** da kontrol et (Sağ tık → Özellikler → Dijital İmzalar → "KeePassXC" / "DroidMonkey Apps"). > teyit et: imzalayan kuruluş adı sürümle değişebilir.
- **macOS:** Notarize edilmiş `.dmg`; Gatekeeper doğrular. Yine de PGP/hash kontrolü en güçlüsüdür.
- **Linux:** Flatpak/dağıtım deposu imzayı kendi doğrular; AppImage indirdiysen **mutlaka** PGP/DIGEST kontrolü yap.

> 🔑 **Püf:** İmza doğrulamayı *her sürüm güncellemesinde* tekrarla. Tedarik zinciri saldırıları en çok "güncelleme" anında işler. KeePassXC'nin kendi otomatik güncelleme indirmesi yoktur (yalnızca "yeni sürüm var" bildirir) — bu **kasıtlı** ve güvenlik açısından iyidir; indirmeyi sen kontrol edersin.

---

<a id="4"></a>
## 4. 🧱 Temel Kavramlar

### `.kdbx` dosyası
Tüm kasan **tek bir şifreli dosyadır**: `Kasam.kdbx`. KeePassXC bugün **KDBX 4** formatını kullanır (Argon2 ve modern şifreleri destekler; eski KDBX 3.1 yalnızca AES-KDF + AES destekler — yeni veritabanı için **KDBX 4 seç**). Bu dosyanın içinde her şey (parolalar, notlar, dosya ekleri, TOTP sırları) şifreli durur.

### Anahtarın yolculuğu (basitleştirilmiş)
```
Ana Parola  +  (Key File)  +  (YubiKey Challenge-Response)
        │
        ▼  birleştirilir → "composite key"
        │
        ▼  KDF (Argon2id — yüksek bellek+iterasyon)   ◄── brute-force'u burada yavaşlatırsın
   Şifreleme Anahtarı
        │
        ▼  .kdbx gövdesini çözer (AES-256 / ChaCha20)
   Çözülmüş veritabanı ──► YALNIZCA RAM'de, yalnızca kasa AÇIKKEN
```

**Kritik sonuçlar:**
1. Composite key'in **her bileşeni gereklidir.** Key file kullanıyorsan, parola doğru olsa bile key file olmadan **asla** açılmaz.
2. KDF (Bölüm 6), saldırganın her parola denemesini **çok pahalı** hale getirir — zayıf parolayı bile bir miktar korur, güçlü parolayı pratikte kırılmaz yapar.
3. Veritabanı **açıkken** içerik RAM'de **çözülüdür** → kilitleme disiplini hayatidir (Bölüm 16).

### Master key bileşenleri (üç faktör)
- **Ana Parola** (bildiğin bir şey)
- **Key File** (sahip olduğun bir dosya) — Bölüm 8
- **Hardware Key / YubiKey Challenge-Response** (sahip olduğun bir cihaz) — Bölüm 9

Üçünü birden kullanmak = **gerçek çok faktörlü** kasa (Bölüm 5).

---

<a id="5"></a>
## 5. 🔐 Veritabanı Oluşturma & 🔥 ÜÇLÜ Kilit

İlk kasanı **doğru** kur — sonradan ekleme yapılabilir ama baştan sağlam temel at.

### Adım adım
1. **Database → New Database.**
2. **General:** Veritabanı adı + açıklama (bunlar **şifrelenmez**, görünür kalabilir — masum bir ad seç; Bölüm 16.7).
3. **Encryption Settings (kritik):**
   - **Encryption Algorithm:** AES-256 (varsayılan, donanım hızlandırmalı) veya ChaCha20 (Bölüm 7).
   - **Key Derivation Function:** **Argon2id** (Bölüm 6 — donanımına göre ayarla, "Benchmark 1.0 s" yetmez, elle yükselt).
   - **Format:** KDBX 4 (Argon2 için zorunlu).
4. **Master Key — işte ÜÇLÜ kilit:**

```
╔═══════════════════════════════════════════════════════════╗
║  COMPOSITE KEY  =  3 bağımsız faktör (hepsi gerekir)       ║
╠═══════════════════════════════════════════════════════════╣
║  ① Ana Parola      → ezberlediğin uzun passphrase         ║
║       (bildiğin)      "doğru-at-pil-zımba-volkan-78"       ║
║                                                            ║
║  ② Key File        → ayrı USB'de duran rastgele dosya      ║
║       (sahip olduğun) Kasam.keyx  (asla değişmez)          ║
║                                                            ║
║  ③ YubiKey (CR)    → fiziksel donanım, dokunmadan açılmaz ║
║       (sahip olduğun) HMAC-SHA1 Challenge-Response          ║
╚═══════════════════════════════════════════════════════════╝
       ▼ üçü birleşir ▼
   Bu üçünü AYNI ANDA ele geçirmek pratikte imkânsızdır.
```

   - **Add Password:** Uzun passphrase gir (Bölüm 14 — diceware, 5-7 kelime).
   - **Add Key File / Add Additional Protection → Key File:** Var olanı seç ya da **Generate** ile yeni rastgele key file üret (Bölüm 8). Bunu **ayrı bir USB'ye** koy.
   - **Add Additional Protection → Add Challenge Response:** YubiKey'ini HMAC-SHA1 modunda seç (Bölüm 9).
5. **Kaydet:** `.kdbx` dosyasını **masum bir yere ve adla** kaydet (Bölüm 16.7).

> 🔥 **Üçlü kilidin gücü:** Tek başına parola → keylogger riski. Parola + key file → key file çalınırsa zayıflar. Parola + key file + YubiKey → saldırganın **üçünü birden** (ezberindeki passphrase + ayrı USB'deki dosya + cebindeki fiziksel anahtar) ele geçirmesi gerekir. Bu, kasanı **devlet-seviyesi** rakibe karşı bile ciddi bir kaleye dönüştürür.

> ⚠️ **Üçlü kilidin bedeli:** Her ek faktör bir **kayıp riskidir.** Key file'ı kaybedersen VEYA YubiKey'in bozulur/kaybolursa **kasa açılmaz** — yedeği şart (Bölüm 8, 9). Faktör eklemek güvenliği artırır ama **kurtarma planı olmadan** felakete davetiyedir.

---

<a id="6"></a>
## 6. ⚙️ KDF Derinlemesine — Argon2id vs AES-KDF

KDF (Key Derivation Function), ana parolanı şifreleme anahtarına dönüştürürken işi **kasıtlı olarak yavaşlatır.** Amaç: saldırganın her parola denemesini pahalı kılmak. Anlamadan kullanan çok; ustaca ayarlayan az. **Bu, brute-force direncinin kalbidir.**

### İki seçenek
| KDF | Açıklama | Ne zaman |
|---|---|---|
| **Argon2id** | 2015 Password Hashing Competition galibi. **Bellek-zor (memory-hard):** GPU/ASIC ile paralel saldırıyı pahalı kılar. id varyantı yan-kanal + GPU direncini dengeler. | **VARSAYILAN, ÖNERİ** — modern, en güçlüsü |
| **AES-KDF** | Sadece iterasyonlu AES (eski KeePass uyumu). Bellek-zor **değil** → GPU/ASIC ile çok daha hızlı saldırılabilir. | Yalnızca eski KDBX 3.1 uyumu gerekiyorsa |

### Argon2id'in üç parametresi (kritik — donanımına göre AYARLA)
KeePassXC üç değer ister:

| Parametre | Ne yapar | Pratik öneri |
|---|---|---|
| **Iterations (t)** | Kaç tur dönecek (zaman maliyeti) | 2–10 (bellek yüksekse düşük tut) |
| **Memory (m)** | Her denemede kaç MiB RAM tüketecek | **YÜKSEK = GPU saldırısını öldürür.** 256 MiB–1 GiB |
| **Parallelism (p)** | Kaç paralel iş parçacığı | Mantıksal çekirdek sayın (örn. 4–8) |

> 🔥 **Püf — "Benchmark butonu yetmez":** KeePassXC'nin **varsayılanı düşüktür** (uyumluluk için, ~64 MiB civarı). "1.0 second" benchmark, *senin güçlü makinende* hızlı çözer — **ama saldırganın 64-GPU'lu rig'i için de hızlıdır.** Bellek-zorluğun tüm anlamı, **bellek miktarını yükseltmektir.** Hedef: **mount süren ~1 sn olacak şekilde ama bellek olabildiğince yüksek** (örn. 256 MiB–1 GiB) + iterasyonu makul tut. Yüksek bellek, GPU/ASIC çiftliklerini ekonomik olarak çökertir — bu, parametre ayarının asıl gücüdür.

> 🧠 **CTI ayar reçetesi:**
> - **Masaüstü/güçlü laptop (kritik kasa):** Memory **512 MiB–1 GiB**, Iterations 4–10, Parallelism = çekirdek sayısı. Açılış 1–2 sn sürer; kabul et — sen günde birkaç kez açarsın, saldırgan **milyarlarca** kez denemek zorunda.
> - **Zayıf cihazda da açılacaksa (telefon/eski makine):** Belleği o cihazın RAM'inin yarısını aşmayacak şekilde tut (yoksa o cihazda **hiç açılmaz**). Tüm cihazların en zayıfına göre tavanı belirle.
> - **Önemli:** Parametreleri **veritabanına yazarsın**; her cihaz onu okur. En düşük donanımlı cihazın belleği, üst sınırını belirler. Telefonda 1 GiB Argon2 çökebilir.

> ⚠️ **AES-KDF tuzağı:** Eğer eski KeePass/KeePassDX uyumu için AES-KDF kullanıyorsan, iterasyon sayısını **yüksek** tut (milyonlar mertebesi) — ama yine de Argon2id bellek-zorluğu kadar GPU-dirençli **değildir.** Mümkünse her şeyi KDBX 4 + Argon2id'e taşı.

---

<a id="7"></a>
## 7. 🔐 Şifreleme Algoritması — AES-256 vs ChaCha20

KDF anahtarı türettikten **sonra**, gerçek veriyi şu simetrik şifrelerden biri korur:

| Seçim | Açıklama | Ne zaman |
|---|---|---|
| **AES-256** | NIST standardı, **AES-NI** ile donanım hızlandırmalı (çoğu modern CPU). Çok hızlı, en yaygın, en çok denetlenmiş. | **Varsayılan, güvenli tercih** |
| **ChaCha20** | Modern akış şifresi (Bernstein). Donanım AES'i olmayan ortamlarda (eski/gömülü cihaz) yazılımda **AES'ten hızlı**, zamanlama saldırılarına dirençli. | AES-NI yoksa, modern tercih isteyen |
| **Twofish** | (Bazı sürümlerde plugin) AES finalisti, sağlam. | Niş alternatif |

> 🔥 **Püf — gerçek fark küçük:** Bir parola yöneticisinde şifrelenen veri **megabaytlar** mertebesindedir (gigabaytlık disk değil), bu yüzden AES vs ChaCha20 **hız farkı pratikte önemsizdir.** İkisi de bugün **kırılmadı.** AES-NI olan modern bir PC'de **AES-256** mükemmeldir. AES'in donanım hızlandırması olmayan **çok eski/gömülü** bir cihazda da kasayı açacaksan **ChaCha20** seç. **Asıl güvenlik şifrede değil, KDF'de (Bölüm 6) ve ana parolanın gücündedir (Bölüm 14).** Şifre seçimine takılıp KDF'yi zayıf bırakmak en yaygın hatadır.

---

<a id="8"></a>
## 8. 🗝️ Key File Stratejisi (Derinlemesine)

Key file, "bildiğin bir şey" (parola) yanına "sahip olduğun bir şey" (dosya) ekler → ikinci faktör. KeePassXC composite key'e karıştırır: **parola doğru olsa bile, key file olmadan kasa açılmaz.**

### Nasıl çalışır
KeePassXC iki tür key file destekler:
- **Üretilmiş key file (`.keyx` — XML v2):** KeePassXC'nin kendi ürettiği, içinde rastgele anahtar **ve** SHA-256 bütünlük hash'i bulunan dosya. **Önerilen budur** — bozulma/değişme tespiti yapar.
- **Rastgele bir dosya (eski yöntem):** Herhangi bir dosyanın içeriğinin hash'i anahtar olur. **Riskli** — dosya bir bayt bile değişirse kasa açılmaz.

> ℹ️ Key file tipi **içeriğe göre** belirlenir; **uzantının önemi yoktur** (`.keyx`, `.key`, ya da uzantısız fark etmez — KeePassXC içeriğe bakar). Bu yüzden key file'ı masum bir adla saklayabilirsin (Bölüm 16.7).

### 🔥 Püf noktaları
- **Üreterek oluştur:** Yeni veritabanı kurarken **Generate** ile `.keyx` üret — tahmin/yeniden üretim **imkânsız**, bütünlük korumalı.
- **Değişen dosyayı ASLA key file yapma:** Bir Word belgesi, log, ya da senkronize bir fotoğraf → bir gün değişir → kasa **sonsuza dek** açılmaz. KeePassXC bu yüzden **amaca özel (purpose-built) key file** önerir, rastgele büyük dosya değil.
- **Nerede tut — AYRI ortamda (kritik):** Key file'ı veritabanından **AYRI fiziksel ortamda** sakla: küçük bir USB, bir donanım anahtarı, hatta ayrı bir cihaz. Mantık: `.kdbx` çalınsa bile key file yanında olmaz → **veri açılamaz.** "Diski/dosyayı çaldım ama açamıyorum" senaryosunu bu yaratır. Key file'ı `.kdbx` ile **aynı klasörde tutmak**, ikinci faktörü **anlamsızlaştırır** (en yaygın hata).
- **Gizleme:** Key file'ı binlerce dosya arasına masum bir adla serpiştir (`IMG_2847.dat`, `font_cache.bin`) — saldırgan hangisinin key file olduğunu bilemez.
- **🚨 Ölümcül risk — KAYIP = VERİ KAYBI:** Key file'ı kaybedersen ya da bozulursa, ana parolan doğru olsa bile **kasa SONSUZA dek açılmaz.** Bu yüzden key file'ın **güvenli yedeği** şart: en az iki ayrı güvenli USB'de tut (ama her kopya da güvenli/şifreli olmalı). Key file yedeğini ana parola yedeğinden **ayrı yerde** tut.

> 🧠 **En güçlü kombinasyon:** Uzun passphrase **+** üretilmiş key file (ayrı USB'de, yedekli) **+** YubiKey CR (Bölüm 9). Üç faktörü birden ele geçirmek pratikte imkânsızdır — ama her birinin **yedeği** de planlı olmalı.

---

<a id="9"></a>
## 9. 🔑 YubiKey / Challenge-Response Derinlemesine

KeePassXC, YubiKey (ve OnlyKey gibi uyumlu cihazlar) ile **HMAC-SHA1 Challenge-Response** destekler. Bu, kasanı bir **fiziksel donanım anahtarına** bağlar.

### Nasıl çalışır (önemli — diğer 2FA'dan farklı)
- YubiKey'in bir slotunda gizli bir **HMAC-SHA1 sırrı** saklanır (programlanır).
- KeePassXC, veritabanının **master seed**'ini (her kayıtta değişen rastgele bayt dizisi) **challenge** olarak YubiKey'e gönderir.
- YubiKey, sır + challenge ile bir **response** üretir ve geri verir.
- Bu response, **şifreleme anahtarının türetilmesine karışır.** Yani YubiKey **giriş kapısı değil**, doğrudan **şifreleme anahtarının parçasıdır** — bulut yöneticilerinin "giriş 2FA"sından **kavramsal olarak daha güçlü.**

```
KeePassXC                          YubiKey
   │  master seed (challenge)  ──────►  │
   │                                    │ HMAC-SHA1(sır, challenge)
   │  ◄──────  response ────────────────│
   ▼
 response → composite key → şifreleme anahtarı
 (YubiKey takılı + dokunulmuş değilse açılmaz)
```

### Slot ayarı
- YubiKey Manager ile bir slotu **"HMAC-SHA1 Challenge-Response"** moduna programla.
- **"Require touch"** (dokunma şartı) aç → kasa açılırken YubiKey'e **fiziksel dokunmadan** response üretilmez. Bu, bir kötücül yazılımın YubiKey takılıyken **senin haberin olmadan** kasayı açmasını engeller (uç nokta uzlaşmasına karşı ciddi bir kalkan).

### 🚨 Yedekleme — KAYIP = VERİ KAYBI (en kritik nokta)
- **Sır, YubiKey'in içindedir ve dışa okunamaz** (kasıtlı). YubiKey'in bozulur/kaybolursa, **aynı sırla programlanmış ikinci bir YubiKey** olmadan kasa **sonsuza dek açılmaz.**
- **Çözüm:** İlk programlamada kullandığın **HMAC sırrını** (genelde elle girdiğin/ürettiğin hex değer) **güvenli sakla** (kâğıda yaz + güvenli kasaya koy ya da ayrı şifreli ortamda tut) **VE** **ikinci bir YubiKey'i aynı sırla programla.** İkisini coğrafi olarak ayır.
- Kayıp durumunda kurtarma için: **ana parola + key file + yedek HMAC sırrı** üçü birden gerekir. Birini bile kaybedersen kasa gider.

### ⚠️ Bulut senkronizasyonu uyarısı (kritik kısıtlama)
- KeePassXC, challenge olarak **master seed** kullanır ve **master seed her kaydetmede değişir** → response da değişir → **her kaydetmede YubiKey'e dokunman gerekir.**
- Bu, **bulut senkronizasyonunu sorunlu kılar:** Syncthing/bulut ile sürekli yazılan bir veritabanında her yazma YubiKey ister; ayrıca farklı cihazlardan eşzamanlı düzenleme çakışma yaratabilir.
- **Çözüm deseni (Bölüm 16.6 / 10):** **Çoklu veritabanı** kullan:
  - **Kritik kasa** (BIOS/BitLocker/seed/kurtarma) → parola + key file + **YubiKey CR**, **senkronize ETME** (yerel + manuel yedek).
  - **Günlük kasa** (web siteleri) → parola (+key file), **senkronize ET** (Bölüm 15).

---

<a id="10"></a>
## 10. 🔢 TOTP (2FA) Saklama — Tek Sepet Tartışması

KeePassXC, her girdi için **TOTP (zaman-bazlı tek kullanımlık kod — Google Authenticator'ın ürettiği 6 haneli kodlar)** saklayabilir ve üretebilir. Girdiye sağ tık → **TOTP → Set up TOTP** ile QR/secret eklersin; sonra **TOTP → Show / Copy** ile anlık kod alırsın.

### 🔥 Tek yerde toplama — risk mi, avantaj mı? (dürüst tartışma)

| Görüş | Argüman |
|---|---|
| **Avantaj** | Tek yerde yönetim, kolay yedek, telefon kaybında 2FA kaybolmaz, KeePassXC açık kaynak ve çevrimdışı. Auto-Type ile parola+TOTP otomatik girilir. |
| **Risk** | TOTP'yi **parolayla aynı kasada** tutmak, "iki faktörü tek sepete" koyar: kasa düşerse **hem parola hem 2FA** düşer → 2FA'nın anlamı (ayrı faktör olması) zayıflar. |

> 🧠 **CTI dengeli öneri:**
> - **Kasan zaten çok güçlüyse** (üçlü kilit + güçlü KDF + temiz uç nokta), TOTP'yi içinde tutmak **çoğu birey için makul** — kolaylık güvenlikten önemli ölçüde feda ettirmez ve telefon-kaybı senaryosunu çözer.
> - **Maksimum güvenlik / yüksek tehdit** istiyorsan: **kritik hesapların TOTP'sini AYRI bir araçta** tut (ayrı bir KeePassXC veritabanı **ya da** ayrı bir cihazdaki authenticator **ya da** fiziksel YubiKey'in kendi OATH/FIDO2'su). Böylece parola kasası düşse bile 2FA ayrı kalır.
> - **En kritik hesaplarda** (e-posta, banka, kripto) TOTP yerine **FIDO2/WebAuthn donanım anahtarı** kullan — TOTP'den (paylaşılan sır) üstündür ve phishing'e dirençlidir.

> 💡 **Pragmatik orta yol:** Günlük/düşük-riskli hesapların TOTP'sini KeePassXC'de tut (kolaylık); e-posta + banka + kripto gibi "kale anahtarı" hesapları için 2FA'yı **ayrı** tut.

---

<a id="11"></a>
## 11. ⌨️ Auto-Type & Pano Güvenliği

Parolayı uygulamaya/siteye aktarmanın iki yolu var: **panoya kopyala** (riskli) ya da **Auto-Type** (kısmen daha güvenli). İkisinin de tuzakları var.

### Pano (Clipboard) — sessiz sızıntı kaynağı
"Copy Password" deyince parola **panoya** düşer. Pano, sistemdeki **tüm uygulamaların okuyabildiği** ortak bir alandır:
- Pano izleyen kötücül yazılım anında okur.
- **Windows pano geçmişi (Win+V)** parolayı kaydedebilir → sonradan görünür.
- Bazı uygulamalar/senkron araçları panoyu buluta taşır (örn. telefon-PC pano paylaşımı).

> 🔥 **Püf — pano otomatik temizleme:** **Tools → Settings → Security → "Clear clipboard after N seconds"** ayarını **10–15 sn** yap (varsayılan açık ama süreyi teyit et). KeePassXC parolayı bu süre sonunda panodan siler. Yine de bu süre içinde başka bir şey kopyalarsan ya da pano geçmişi açıksa risk sürer → mümkünse **panoyu hiç kullanma, Auto-Type tercih et.** Ek olarak: Windows'ta **pano geçmişini (Win+V) kapat** (Ayarlar → Sistem → Pano → Geçmişi kapat) ya da en azından KeePassXC'nin "exclude from clipboard history" desteğine güven (sürüme göre değişir — teyit et).

### Auto-Type — keylogger'a karşı kısmi avantaj
Auto-Type (varsayılan `Ctrl+Alt+A` ya da global kısayol), KeePassXC'nin kullanıcı adı + Tab + parola + Enter'ı **klavye tuşlarını taklit ederek** aktif pencereye "yazmasıdır." Panoyu hiç kullanmaz.

- **Avantaj (panoya göre):** Parola panoda durmaz → pano sızıntısı yok.
- **Avantaj (keylogger'a karşı KISMİ):** Bazı basit/uygulama-içi keylogger'lar yalnızca fiziksel klavyeyi dinler; sahte tuş olaylarını kaçırabilir. **Ama bu garanti DEĞİLDİR** — çekirdek-seviyesi (kernel) keylogger sahte olayları da yakalar. Auto-Type keylogger'a karşı **kısmi/zayıf** bir avantajdır, **çözüm değil.**
- **Risk — yanlış pencere:** Auto-Type **odaktaki pencereye** yazar. Yanlış pencere öndeyse parolanı **düz metin** olarak oraya yazar (sohbet kutusu, adres çubuğu...). KeePassXC pencere başlığını eşleştirir (Auto-Type association) — bunu doğru kur, "her pencereye yaz" deme.
- **Risk — phishing:** Auto-Type pencere başlığına bakar, **gerçek alan adına bakmaz** (tarayıcıda) → sahte sayfa doğru başlığı taşıyorsa yanlış yere yazabilir. Tarayıcıda **KeePassXC-Browser** (Bölüm 12) alan adını eşlediğinden phishing'e karşı Auto-Type'tan **daha güvenlidir.**

> 🧠 **Öncelik sırası (en güvenliden):** ① KeePassXC-Browser (alan adı eşlemeli, pano yok) → ② Auto-Type (pano yok, ama pencere riski) → ③ Pano kopyala (en riskli, yalnızca kısa süre + temizleme ile). Keylogger'lı bir makinede **hiçbiri** seni kurtarmaz — temiz uç nokta şarttır.

---

<a id="12"></a>
## 12. 🌐 Tarayıcı Entegrasyonu (KeePassXC-Browser)

**KeePassXC-Browser** eklentisi (Firefox/Chrome/Edge/Brave), tarayıcı ile KeePassXC arasında **şifreli yerel bir kanal** kurar ve parolaları doğru sitede otomatik doldurur.

### Kurulum
1. KeePassXC: **Settings → Browser Integration → Enable** + kullandığın tarayıcıyı işaretle.
2. Tarayıcıya **resmi** "KeePassXC-Browser" eklentisini kur (resmi mağazadan; sahte kopyalara dikkat).
3. Eklenti ↔ KeePassXC ilk bağlantıda **eşleşir** (bir anahtar değişimi yapar, onayla).
4. Bir sitede ilk girişte "bu girdiyi bu siteyle ilişkilendir" → bundan sonra otomatik doldurur.

### Nasıl güvenlik sağlar (phishing'e karşı)
- Eklenti, parolayı **yalnızca KeePassXC girdisindeki URL ile EŞLEŞEN alan adında** önerir. Sahte `paypa1.com` sitesine `paypal.com` parolasını **önermez** → phishing'e karşı güçlü, otomatik bir kalkan.
- Parola panoya **düşmeden** doğrudan form alanına gider.

### ⚠️ Riskleri
- **Saldırı yüzeyi büyür:** Tarayıcı = en çok saldırılan yazılım. Eklenti, kasanı tarayıcıya açar (yalnızca kasa açıkken). Tarayıcı/eklenti bir zafiyet taşırsa, açık kasaya köprü oluşur.
- **Yerel kanal hedefi:** Eklenti ile KeePassXC arasındaki yerel mesajlaşma (native messaging) kötücül başka bir eklenti tarafından taklit edilmeye çalışılabilir → KeePassXC bağlantı onayı + alan adı eşleme bunu sınırlar ama **sıfır risk değildir.**
- **"Her şeyi doldur" tembelliği:** Eklentiyi "sormadan doldur"a alma; her dolduruşta site doğruluğunu görsel teyit et.

> 🔥 **Püf — ayrı kasa + en az yetki:** Tarayıcı entegrasyonunu yalnızca **günlük/web kasanda** kullan; **kritik kasanı** (BIOS/seed/kurtarma) tarayıcıya **hiç açma** (Browser Integration o kasa için kapalı kalsın). KeePassXC'de **"Only return best matching credentials"** ve **HTTP siteler için uyarı** seçeneklerini aç. Eklenti izinlerini minimumda tut.

---

<a id="13"></a>
## 13. 🔐 SSH Agent Entegrasyonu

KeePassXC, SSH özel anahtarlarını kasada şifreli tutup, kasa açıkken bunları **SSH agent'a** sunabilir → her bağlantıda parola/passphrase girmeden, ama anahtar **diskte düz değil, kasada şifreli** durur.

### Kurulum (özet)
1. KeePassXC: **Settings → SSH Agent → Enable.** (Windows'ta OpenSSH agent ya da Pageant; Linux/macOS'ta `ssh-agent` / `SSH_AUTH_SOCK`.)
2. SSH özel anahtarını bir KeePassXC girdisine **Advanced → Attachment** olarak ekle (örn. `id_ed25519`).
3. Girdinin **SSH Agent** sekmesinde: "Add key to agent when database is opened/unlocked" + isteğe bağlı "Remove key when database is locked."
4. Artık kasa açıkken `ssh` komutu anahtarı agent'tan alır; kasa kilitlenince anahtar agent'tan **düşer** (ayarladıysan).

### 🔥 Püf / Riskler
- **Avantaj:** Özel anahtar diskte **düz/passphrase'siz dosya olarak durmaz** → disk çalınsa anahtar şifreli kasada. Kasa kilidi = SSH erişimi kapanır.
- **"Remove key when database is locked" AÇ:** Yoksa kasayı kilitlesen bile anahtar agent'ta (RAM'de) açık kalır → kilitlemenin amacı zayıflar.
- **Agent forwarding tehlikesi:** Uzak sunucuya `ssh -A` (agent forwarding) ile bağlanma alışkanlığı, ele geçirilmiş sunucunun agent'ına erişip anahtarını kullanmasına yol açabilir → forwarding'i yalnızca güvendiğin hop'larda, mümkünse hiç kullanma.
- **Anahtarın kendisine passphrase koy:** KeePassXC kasası + SSH anahtarının kendi passphrase'i = iki katman.

---

<a id="14"></a>
## 14. 🎲 Parola/Passphrase Üretici (Diceware & Entropi)

KeePassXC'nin yerleşik üreticisi (girdi oluştururken "zar" ikonu ya da **Tools → Password Generator**) iki mod sunar: **Password** (rastgele karakter) ve **Passphrase** (rastgele kelime — diceware).

### İki mod
| Mod | Örnek | Ne zaman |
|---|---|---|
| **Password** | `k7$Wm@2pX!qL9zR#` | Site parolaları (KeePassXC hatırlar, sen ezberlemezsin) → **uzun + tüm karakter sınıfları**, 20+ karakter |
| **Passphrase (diceware)** | `volkan-zımba-pil-atış-78-koru` | **Ana parola** (ezberlemen gereken) → kelime tabanlı, hatırlaması kolay, entropisi yüksek |

### Entropi — gerçek güç ölçüsü
**Entropi (bit)**, parolanın tahmin-direncidir. Her bit, brute-force maliyetini **iki katına** çıkarır.
- Diceware: her rastgele kelime ~**12.9 bit** (7776 kelimelik listeden). **5 kelime ≈ 64 bit**, **7 kelime ≈ 90 bit.**
- **Ana parola için hedef: en az 5-7 diceware kelimesi (≥ 64-77 bit).** Güçlü KDF (Bölüm 6) ile birleşince offline brute-force pratikte imkânsız olur.
- Site parolaları için: 20+ karakter rastgele → KeePassXC üretsin, sen asla görmesen/ezberlemesen de olur.

> 🔥 **Püf — "parola gücü ölçer yanılgısı":** Web sitelerindeki ve bazı araçlardaki renkli "güçlü/zayıf" çubuğu **entropiyi gerçekten ölçmez** — çoğu yalnızca uzunluk + karakter çeşidine bakar ve **kalıpları** (klavye dizisi, kelime+sayı, l33t değişimi) yakalayamaz. `P@ssw0rd123!` çoğu ölçere "güçlü" görünür ama saldırgan sözlüğünde **ilk denemelerdedir.** **Tek güvenilir ölçü: gerçek rastgelelikten gelen entropi.** Bu yüzden ana parolanı **kendin uydurma** (insan beyni rastgele değildir) — **KeePassXC'nin diceware üreticisiyle ürettiğini ezberle.** "Akıllıca" uydurduğun parola, üreticinin ürettiğinden neredeyse her zaman zayıftır.

> 🧠 **Ana parola reçetesi:** Password Generator → Passphrase → **6-7 kelime**, kelimeler arası ayraç (`-` ya da boşluk), istersen bir kelimeyi büyük harfle/sayıyla zenginleştir. Üret → **birkaç gün boyunca yazarak ezberle** → kâğıt yedeğini güvenli kasaya koy → dijital izini sil.

---

<a id="15"></a>
## 15. 💾 Yedekleme & Senkronizasyon

`.kdbx` **şifreli** olduğundan, onu yedeklemek/senkronlamak güvenlidir — dosya çalınsa bile içerik (güçlü parola + KDF varsa) açılamaz. Ama **yedek = hayat**; tek kopya **ölümcül hatadır.**

### Yedekleme (3-2-1 kuralı)
- **3 kopya, 2 farklı ortam, 1 saha dışı (offsite).** Örn: yerel disk + harici USB + uzak/şifreli konum.
- KeePassXC **kaydetmeden önce otomatik yedek** alabilir: **Settings → General → "Backup database file before saving"** → her kayıtta `.kdbx.old` (ya da zaman damgalı) yedek tutar. **Aç.**
- **Key file ve YubiKey HMAC sırrını AYRI yedekle** (Bölüm 8, 9) — `.kdbx` yedeği tek başına işe yaramaz, faktörler de gerekir.
- **Test et:** Yedeği **düzenli olarak başka makinede aç** → açılmayan yedek, yedek değildir.

### Senkronizasyon — cihazlar arası
`.kdbx` tek dosya olduğundan herhangi bir senkron aracıyla taşınır:

| Yöntem | Güven | Not |
|---|---|---|
| **Syncthing** | **En iyi** — uçtan uca, sunucusuz, açık kaynak | Cihazların arasında doğrudan; bulut yok. **Önerilen.** |
| **Şifreli bulut** (Cryptomator + Drive/Dropbox) | İyi | `.kdbx` zaten şifreli; ek katman olarak Cryptomator |
| **Düz bulut** (Drive/Dropbox/OneDrive) | Kabul edilebilir | `.kdbx` şifreli olduğundan içerik korunur; ama metadata/erişim şirkette |
| **VeraCrypt birimi içinde + manuel** | En sağlam (kritik kasa) | Bölüm 16.10 |

> 🔥 **Püf — çakışma (conflict) yönetimi:** İki cihazda **aynı anda** düzenleyip ikisini de kaydedersen, senkron araçları **çakışan kopya** (`Kasam.kdbx.sync-conflict-...`) üretir. KeePassXC bunu otomatik **birleştirmez** (dosya seviyesi senkronda). Çözüm:
> - KeePassXC'nin **Database → Merge from Database** özelliğini kullan → iki `.kdbx`'i girdi seviyesinde birleştir (her iki taraftaki değişiklikleri korur).
> - Daha iyisi: **bir cihazda düzenle, kapat, senkronun bitmesini bekle, sonra diğerinde aç** → çakışmayı baştan önle.
> - YubiKey CR'li kasayı **senkronlama** (her kayıt master seed'i değiştirir, çakışma + her cihazda YubiKey gerekir — Bölüm 9).

> ⚠️ **Senkron + şifreleme uyumu:** `.kdbx`'i senkronlarken **kayıt sırasında** senkron aracının yarı-yazılmış dosyayı kopyalamamasına dikkat et (nadiren bozulma). KeePassXC atomik yazma yapar ama yine de "backup before saving" açık olsun.

---

<a id="16"></a>
## 16. 🔥 PÜF NOKTALARI — Piyasada Bulamayacakların

Bu bölüm, çoğu rehberin atladığı ve gerçek dünyada güvenliği **çökerten ya da kurtaran** detaylardır.

### 16.1 Argon2 Parametrelerini Donanımına Göre Ayarla
Varsayılan KDF ayarı **uyumluluk için düşüktür** ve güçlü makinende anlamsız hızlı çözer → saldırganın GPU'su için de hızlı. **Belleği elle yükselt** (256 MiB–1 GiB), iterasyonu makul tut, açılış ~1-2 sn olsun. Bellek-zorluk, GPU/ASIC saldırı çiftliklerini ekonomik olarak öldürür. (Bölüm 6 — bu, en yüksek getirili tek ayardır.)

### 16.2 Key File'ı Veritabanından AYRI Ortamda Tut
`.kdbx` ve key file **aynı klasörde/diskte** ise, ikinci faktör **anlamsızdır** — ikisini birden çalan her şeyi açar. Key file'ı **ayrı USB'de** tut → "dosyayı çaldım ama açamıyorum" senaryosunu yarat. (Bölüm 8.)

### 16.3 Veritabanı Kilitleme — Açık Kasa = RAM'de Çözülü
**En kritik operasyonel gerçek:** Veritabanı **açıkken** tüm parolalar **bellekte çözülü** durur. KeePassXC'yi şu durumlarda **otomatik kilitle** (Settings → Security):
- **"Lock databases after inactivity of N seconds"** → boşta kalınca (örn. 60-300 sn) kilitle.
- **"Lock databases when session is locked or lid is closed"** → `Win+L` / kapak kapanınca kilitle.
- **"Lock databases when minimizing the main window"** → küçültünce kilitle.
- **"Lock databases after switching the screen lock"** + ekran koruyucu tetikleyicisi.

Kilitlenince KeePassXC, çözülmüş veriyi **RAM'den temizler** (composite key'i unutur). **Masadan kalkarken kasa kilitli olmalı.**

### 16.4 Pano Temizleme Süresi
"Copy Password" parolayı panoya koyar → **Settings → Security → "Clear clipboard after N seconds"** = 10-15 sn. Ek olarak Windows pano geçmişini (Win+V) kapat. Mümkünse panoyu hiç kullanma → Auto-Type/Browser. (Bölüm 11.)

### 16.5 "Açık Veritabanı = RAM'de" Gerçeği & Cold-Boot
- Kasa açıkken composite key ve çözülmüş içerik **RAM'dedir.** Cold-boot saldırısında (makine açık/uyku/hazırda iken RAM dondurulup çekilir) bu sızabilir.
- **Çözüm:** İş biter bitmez **kilitle** (RAM temizlenir); makineyi bırakırken **tamamen kapat** (uyku/hazırda değil) → RAM boşalır. Hassas senaryoda **hibernate'i kapat** (`powercfg -h off`) — yoksa açık kasa hazırda bekletmede `hiberfil.sys`'e (şifresiz disk) düşebilir.
- Kanije Kalesi entegrasyonu (Bölüm 18): fiziksel tehditte otomatik kilit/kapatma.

### 16.6 Çoklu Veritabanı — Günlük / Kritik Ayrımı
Her şeyi tek kasaya koyma. **En az iki `.kdbx`:**
- **Günlük kasa:** Web siteleri, forumlar, düşük-riskli → parola (+key file), tarayıcı entegrasyonu açık, **senkronize.**
- **Kritik kasa:** BitLocker/BIOS/kurtarma anahtarları, kripto seed, VeraCrypt parolaları, banka → **üçlü kilit (parola+key file+YubiKey)**, tarayıcı entegrasyonu **kapalı**, **senkronize edilmez**, yalnızca gerektiğinde açılır.

Bu, "günlük kasa bir şekilde düşse bile kale anahtarları ayrı ve daha korunaklı" demektir.

### 16.7 Veritabanı Dosyasını Masum Adla Gizleme
`.kdbx` uzantısı "burada bir parola kasası var" der. KeePassXC dosyayı **içeriğe göre** açar, ada/uzantıya bakmaz → kasayı `vergi_2019.dat`, `proje_yedek.bin`, `aile_videolari.dat` gibi **masum bir adla** sakla. Aynı şey key file için de geçerli (Bölüm 8). Tam gizlilik için `.kdbx`'i bir **VeraCrypt birimi içine** koy (16.10) → dosyanın **varlığı** bile gizlenir.

### 16.8 Geçmiş (History) & Silinen Girdiler Kalır
- KeePassXC her girdi için **değişiklik geçmişi** tutar (eski parolalar dahil) → bir parolayı değiştirsen bile **eskisi geçmişte durur.** İyi yan: yanlışlıkla silineni kurtarırsın. Riskli yan: **sızdırılmış eski parola hâlâ kasada.** Girdi → **History** sekmesinden eski kayıtları **temizle** (özellikle bir parola gerçekten ele geçirildiyse).
- **Silinen girdiler "Recycle Bin" grubuna gider** (kalıcı silinmez) → çöp kutusunu **boşalt** (sağ tık → Empty Recycle Bin). Yoksa "sildim" sandığın hassas girdi kasada kalır.
- **Ayar:** İstersen geçmiş tutma limitini düşür (Database Settings → History) ama kurtarma yeteneğini de feda edersin — denge.

### 16.9 Entropi / Parola Gücü Ölçer Yanılgısı
KeePassXC'nin (ve sitelerin) gösterdiği güç çubuğu **kalıpları tam yakalayamaz.** Tek güvenilir güç ölçüsü **gerçek entropidir** → ana parolanı **kendin uydurma**, **diceware üreticisiyle üret** (Bölüm 14). "Güçlü" görünen insan-yapımı parola, üreticinin ürettiğinden genelde zayıftır.

### 16.10 `.kdbx`'i VeraCrypt Birimi İçinde Tutmak — İkinci Katman
Kritik kasayı bir **VeraCrypt (tercihen hidden) birimi** içinde sakla:
- Birim **kapalıyken** `.kdbx`'in **varlığı bile görünmez** (yalnızca şifreli blob).
- Birim açıkken KeePassXC'yi içinden çalıştır; iş bitince **dismount** → çift kilit (VeraCrypt parolası + KeePassXC üçlü kilit).
- Zorlanma (rubber-hose) senaryosunda VeraCrypt **hidden volume inkâr edilebilirliği** sağlar: dış birimde masum bir kasa göster, kritik kasa gizli birimde görünmez kalsın.
- Bu, **ikinci bağımsız şifreleme katmanı** + **inkâr edilebilirlik** demektir. (Bkz. `VERACRYPT_USTALIK_REHBERI.md`.)

### 16.11 Quick Unlock / Biyometri — Kolaylık-Güvenlik Dengesi
KeePassXC (sürüme göre) **Quick Unlock** (Windows Hello / Touch ID / kısa süreli) sunabilir. Kolaylık sağlar ama: biyometri **ana parolanı RAM'de tutar** ve fiziksel zorlamaya (parmağını/yüzünü kullanmaya zorlama) açıktır. **Yüksek tehditte kapalı tut**; sıradan kullanımda kabul edilebilir. > teyit et: Quick Unlock davranışı/varlığı sürüm ve platforma göre değişir.

### 16.12 Ekran Görüntüsü & "Parolayı Göster" Riski
Parolayı ekranda "göster" (göz ikonu) ile açtığında, ekran görüntüsü alan kötücül yazılım/omuz sörfü yakalayabilir. KeePassXC ana penceresinin ekran görüntüsünü/ekran paylaşımını engellemeye çalışır (platforma göre) ama garanti değildir → hassas parolayı gereksiz yere ekranda **gösterme.**

---

<a id="17"></a>
## 17. ☠️ Yaygın Ölümcül Hatalar

1. **Tek kopya `.kdbx`** → disk ölünce/bozulunca **tüm parolalar bir anda gider.** (En sık felaket — 3-2-1 yedek şart.)
2. **Key file kaybı/değişmesi** → parola doğru olsa bile kasa sonsuza dek açılmaz. (Yedeklenmemiş key file = saatli bomba.)
3. **YubiKey'i yedeksiz kullanmak** → tek YubiKey bozulunca/kaybolunca kasa gider. (İkinci anahtar + HMAC sırrı yedeği şart.)
4. **Zayıf ana parola** → KDF ne olursa olsun "123456"/`P@ssw0rd1` ile kasa kâğıttan kale. (Diceware kullan.)
5. **Kasayı açık bırakmak** → RAM'de çözülü; biri makineye eriştiğinde her şey okunur. (Auto-lock + Win+L.)
6. **Key file'ı `.kdbx` ile aynı yerde tutmak** → ikinci faktör anlamsız.
7. **Parolayı/PIM'i/HMAC sırrını aynı yerde saklamak** → faktör ayrımını yok etmek.
8. **Argon2'yi varsayılan düşük ayarda bırakmak** → GPU brute-force'a kapı aralamak.
9. **Pano geçmişi açıkken parola kopyalamak** → parola Win+V geçmişinde kalır.
10. **İmza doğrulamadan kurmak** → truva atlı yükleyiciyle tüm kasayı sızdırmak.
11. **Keylogger'lı/uzlaşmış makinede ana parolayı yazmak** → kasa tasarımı çaresiz; temiz uç nokta şart.
12. **Eski/sızmış parolaları History + Recycle Bin'de bırakmak** → "değiştirdim/sildim" sandığın sır kasada durur.
13. **Tek kasaya her şeyi (günlük + kritik + TOTP) tıkmak** → tek nokta çöküşü; çoklu kasa ayır.

---

<a id="18"></a>
## 18. 🏰 Kanije Kalesi ile Birlikte Kullanım

Bu repo (Kanije Kalesi), **fiziksel tehdit anında** cihazı uzaktan/otomatik koruyan bir muhafızdır. KeePassXC **durağan parolaları** korur; Kanije **olay anını** yönetir. İkisi katmanlı bir savunma oluşturur:

| Senaryo | KeePassXC rolü | Kanije Kalesi rolü |
|---|---|---|
| Cihaz çalındı (kasa kapalı) | Parolalar şifreli, erişilemez | — |
| Cihaz açık, kasa açık, sen uzaktayken biri yaklaştı | Kasa RAM'de çözülü (**risk**) | `/koruma` dead-man / USB / yanlış-giriş tetikleyici → **kilitle + alarm + foto** |
| Acil durum | — | `/panik`, `/kilit tam` (lockdown) |
| Kasayı kalıcı yok etme | `.kdbx`'i hedef al | `/imha` ile `.kdbx`'i + key file'ı güvenli sil |
| Adli inceleme öncesi | `.kdbx` VeraCrypt hidden volume'da → inkâr edilebilirlik | RAM-only mod, iz temizleme |

### 🔥 Önerilen entegrasyon deseni
1. **Kritik parolaları KeePassXC'de tut** — BitLocker kurtarma anahtarı, BIOS/UEFI parolası, VeraCrypt birim parolaları, kripto seed, banka. Bunları başka yerde tutmak (kâğıt/dosya) daha risklidir; tek güçlü kasada toplayıp **üçlü kilitle** koru.
2. **`.kdbx`'i VeraCrypt hidden volume'da sakla** (Bölüm 16.10) → varlığı gizli + ikinci şifreleme katmanı + inkâr edilebilirlik.
3. **KeePassXC auto-lock'u** "screensaver/idle/lid/minimize"e bağla → Kanije kilitlerken/ekran kilitlenirken kasa da kilitlenir, RAM'deki composite key uçar.
4. **Kanije `/koruma`'yı aç** — dead-man switch + USB dead-man. Sen X süre check-in yapmazsan ya da key file USB'si çıkarılırsa Kanije devreye girer (kilit + alarm + delil).
5. **Hızlı yok etme (crypto-shredding):** En kritik senaryoda, Kanije'nin `/imha` hedeflerine **`.kdbx` + key file**'ı ekle → kasayı saniyeler içinde güvenli sil. Key file ayrı USB'deyse, o USB'yi fiziksel olarak yok etmek de kasayı erişilemez kılar (key file olmadan açılmaz).
6. **RAM-only çalışma:** Kanije'nin RAM-only modu + KeePassXC'nin auto-lock disiplini = güç kesildiğinde çözülmüş kasa **RAM'den uçar, iz kalmaz.**

> 🧠 **Felsefe örtüşmesi:** KeePassXC sana **çevrimdışı, kimseye güvenmeyen** bir parola kalesi verir; Kanije sen unuttuğunda/tehlikede olduğunda **kapıyı kilitleyen nöbetçidir.** Kritik kasanın anahtarlarını (key file, YubiKey) fiziksel olarak yanında taşırsın; tehdit anında Kanije kilitler, gerekirse `/imha` ile kasa crypto-shredding ile yok edilir. VeraCrypt katmanı (Bölüm 16.10) inkâr edilebilirliği tamamlar.

---

<a id="19"></a>
## 19. ✅ Hızlı Referans & Operasyonel Kontrol Listesi

### Kasa oluştururken
- [ ] Yükleyici **PGP imzası + SHA-256 hash doğrulandı**
- [ ] **KDBX 4** formatı + **Argon2id** KDF
- [ ] Argon2 **belleği elle yükseltildi** (256 MiB–1 GiB, ~1-2 sn açılış)
- [ ] Şifreleme: **AES-256** (ya da AES-NI yoksa ChaCha20)
- [ ] **Ana parola = diceware 6-7 kelime** (üreticiyle üretildi, ezberlendi)
- [ ] **Key File** üretildi (`.keyx`) ve **ayrı USB'ye** kondu
- [ ] **YubiKey Challenge-Response** eklendi + **"require touch"** açık
- [ ] Veritabanı **masum adla** kaydedildi

### Yedek & faktör güvenliği
- [ ] **3-2-1 yedek** (3 kopya, 2 ortam, 1 saha dışı), "backup before saving" açık
- [ ] **Key file yedeği** ayrı ve güvenli (≥2 USB)
- [ ] **YubiKey HMAC sırrı yedeği** + **ikinci YubiKey** aynı sırla programlı
- [ ] Parola / key file / HMAC sırrı **AYRI yerlerde**
- [ ] Yedek **başka makinede açılarak test edildi**

### Güvenlik ayarları (Settings)
- [ ] **Auto-lock:** inactivity + session lock + lid + minimize
- [ ] **Clear clipboard after 10-15 s**
- [ ] Windows **pano geçmişi (Win+V) kapalı**
- [ ] **Browser integration:** yalnızca günlük kasada açık, kritik kasada kapalı
- [ ] **SSH agent:** "remove key when locked" açık (kullanıyorsan)

### Her oturum
- [ ] Temiz / sertleştirilmiş uç nokta (keylogger/RAT yok)
- [ ] Pano yerine **Auto-Type / Browser** tercih
- [ ] Parolayı ekranda gereksiz **gösterme**
- [ ] İş bitince **kasayı KİLİTLE**; masadan kalkarken **Win+L**
- [ ] Makineyi bırakırken **tamamen kapat** (uyku/hazırda değil — kritik kasa açıksa)

### Periyodik
- [ ] **Çoklu kasa** ayrımı korunuyor (günlük / kritik)
- [ ] Sızan/değişen parolaların **History + Recycle Bin temizlendi**
- [ ] Senkron **çakışmaları** Merge ile çözüldü
- [ ] Yedekler güncel ve **test edildi**
- [ ] Yeni sürümde **imza yeniden doğrulandı**

---

> 🏰 **Kapanış:** Parola yöneticisi bir ürün değil, bir **disiplindir.** En güçlü Argon2id ve üçlü kilit bile, kasayı açık bıraktığın bir akşam ya da keylogger'lı bir makinede ana parolanı yazdığın an çaresizdir. KeePassXC sana **çevrimdışı, kimseye güvenmeyen, matematiksel bir parola kalesi** verir — bulutun aksine verin senin elinde, bir şirketin sunucu ihlalinde sızacak merkez yok. Ama bu kalenin tek kapısı senin operasyonel disiplinindir: **kilitli tut, ayrı sakla, yedekle, doğrula.** Kanije Kalesi de tam burada — sen unuttuğunda kapıyı kilitleyen, tehlikede crypto-shredding ile kasayı yok eden nöbetçi — devreye girer.
>
> ⚠️ **Doğruluk notu:** Argon2 varsayılanları, Quick Unlock davranışı, imzalama anahtarı parmak izi ve menü adları **sürüme göre değişebilir** — kendi KeePassXC sürümünden ve resmi `keepassxc.org/docs` adresinden **teyit et**. Yanlış bir güvenlik varsayımı (örn. "key file yedeğim var sandım") tüm kaleyi çökertir.
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. İlgili: `VERACRYPT_USTALIK_REHBERI.md`, `HAYAT_KURTARAN_YAZILIMLAR.md`, `SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md`, `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`.*
