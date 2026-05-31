# 🔒 VERACRYPT — TAM USTALIK REHBERİ
## Kapsayıcıdan Gizli İşletim Sistemine, Püf Noktalarıyla Uçtan Uca

> **Amaç:** VeraCrypt'i "kur-mount-et" seviyesinden çıkarıp, bir tehdit aktörü karşısında **gerçekten ayakta kalacak** şekilde kullanmayı öğretmek. Bu rehber yalnızca *nasıl*'ı değil, **neden** ve **hangi durumda işe yaramaz**'ı da anlatır. Forum cevaplarında bulamayacağın deniability tuzakları, SSD gerçeği, cold-boot, PIM matematiği ve adli bilişim karşıtı detaylar burada.

> ⚠️ **Önce bunu oku:** Şifreleme, *yanlış kullanıldığında* sana **yanlış güvenlik hissi** verir — bu, hiç şifrelememekten daha tehlikelidir. Bir bölümü anlamadan uygulama. Özellikle **Hidden Volume** ve **SSD** bölümlerini atlama.

---

## 📑 İÇİNDEKİLER

1. [VeraCrypt Nedir, Neden? (BitLocker/LUKS Kıyası)](#1)
2. [Tehdit Modeli — Neyi Korur, Neyi KORUMAZ](#2)
3. [Kurulum + İmza Doğrulama (Atlama!)](#3)
4. [Temel Kavramlar — Container / Partition / System / XTS](#4)
5. [Şifreleme & Hash Algoritması Seçimi](#5)
6. [PIM Derinlemesine — Gizli İkinci Faktör](#6)
7. [Keyfile Derinlemesine](#7)
8. [Standart Şifreli Kapsayıcı (Adım Adım)](#8)
9. [Gizli Birim (Hidden Volume) & İnkâr Edilebilirlik](#9)
10. [Sistem Şifreleme (Pre-Boot Authentication)](#10)
11. [Gizli İşletim Sistemi (Hidden OS) — En Üst Seviye](#11)
12. [🔥 PÜF NOKTALARI — Piyasada Bulamayacakların](#12)
13. [Yedekleme & Kurtarma (Header / Rescue Disk)](#13)
14. [Günlük Operasyon Güvenliği](#14)
15. [Komut Satırı & Otomasyon](#15)
16. [Kanije Kalesi ile Birlikte Kullanım](#16)
17. [Yaygın Ölümcül Hatalar](#17)
18. [Hızlı Referans & Operasyonel Kontrol Listesi](#18)
19. [Hukuki Sınır & Sınır Geçişi Notu](#19)

---

<a id="1"></a>
## 1. 🧭 VeraCrypt Nedir, Neden?

VeraCrypt, kapatılan **TrueCrypt**'in (2014) devamıdır. IDRIX (Mounir Idrassi) tarafından sürdürülür, **açık kaynaktır** ve 2016'da bağımsız bir güvenlik denetiminden (QuarksLab/OSTIF) geçmiştir. TrueCrypt'in iterasyon sayısını ciddi şekilde artırarak brute-force direncini güçlendirmiştir.

Ne yapar: Bir **sanal şifreli disk** (kapsayıcı dosya), bir **bölüm** ya da **tüm sistem diskini** "anında" (on-the-fly) şifreler. Veri diskte hiçbir zaman çözülmüş halde durmaz; yalnızca RAM'de, mount edildiği sürece çözülür.

### BitLocker / LUKS ile Dürüst Kıyas

| Özellik | VeraCrypt | BitLocker | LUKS (Linux) |
|---|---|---|---|
| Açık kaynak | ✅ (denetlenmiş) | ❌ (Microsoft) | ✅ |
| Platform | Win/macOS/Linux | Windows (+sınırlı) | Linux odaklı |
| **İnkâr edilebilirlik (hidden volume)** | ✅ **Tek o sunar** | ❌ | ❌ |
| Gizli işletim sistemi | ✅ | ❌ | ❌ |
| TPM ile şeffaf açılım | ❌ (kasıtlı) | ✅ | Sınırlı (Clevis) |
| Kaskad şifreler (AES-Serpent-Twofish) | ✅ | ❌ | ✅ (manuel) |
| Arka kapı şüphesi | Düşük (denetim) | Yüksek (kapalı, devlet erişimi söylentileri) | Düşük |
| Pre-boot PIM/keyfile | ✅ | ❌ | ✅ (sınırlı) |

**CTI / hassas veri perspektifi:** BitLocker, TPM'e güvenir; TPM anahtarı mühürler ve *senin parolan olmadan* (hatta sen uyurken) sistemi çözebilir → fiziksel el koymada DMA/cold-boot ile risklidir ve **kapalı kaynaktır**. VeraCrypt kasıtlı olarak TPM kullanmaz: "Parolayı bilen açar, başka yol yok" felsefesi. Hassas araştırma verisinde **VeraCrypt + parola/keyfile, BitLocker'a göre daha sağlam bir tehdit duruşu** sunar — özellikle inkâr edilebilirlik gerektiğinde rakipsizdir.

> 💡 **Katmanlı yaklaşım:** BitLocker'ı kapatmana gerek yok — ikisi birlikte kullanılabilir (BitLocker tüm disk + VeraCrypt hidden volume kritik veri için). Ama tek başına kritik veriyi korurken VeraCrypt hidden volume tercih et.

---

<a id="2"></a>
## 2. 🎯 Tehdit Modeli — Neyi Korur, Neyi KORUMAZ

Şifrelemeye başlamadan önce **kime karşı** korunduğunu netleştir. Aksi halde yanlış güvenlik hissi yaşarsın.

### ✅ VeraCrypt'in KORUDUĞU senaryolar
- **Disk/cihaz çalınması (kapalı/dismount halde):** Disk başka makineye takılır → veri okunamaz. Bu, en güçlü olduğu senaryodur.
- **El koyma (kapalı cihaz):** Adli inceleme diski klonlar → şifreli blob görür.
- **İnkâr edilebilirlik (hidden volume):** Parola vermeye zorlanırsan, dış birimin parolasını verirsin; gizli birim **matematiksel olarak kanıtlanamaz**.
- **Eski/atılan disk:** Format/silme yetersizken bile şifreli veri çöp kalır.

### ❌ VeraCrypt'in KORUMADIĞI senaryolar
- **Mount edilmiş haldeyken (sistem açık, birim bağlı):** Veri RAM'de çözülü; o anda erişen herkes (kötü amaçlı yazılım, omuz sörfü, uzak erişim) okur.
- **Keylogger / ekran kaydedici:** Parolanı yazarken yakalar. Şifreleme burada çaresizdir.
- **Cold-boot saldırısı (açık/uyku/hazırda):** RAM'deki anahtar fiziksel olarak çekilebilir (aşağıda detay).
- **Rubber-hose (zorla parola alma) — hidden volume yoksa:** Seni parola vermeye zorlarlarsa düz şifreleme çaresiz; **bu yüzden hidden volume var.**
- **Donanım implantı / firmware rootkit:** Pre-boot ortamından önce çalışan kötücül kod parolayı çalabilir (Evil Maid).
- **Sen mount ederken arkanda biri varken:** Sosyal/fiziksel zafiyet.

> 🧠 **Altın kural:** VeraCrypt **"durağan veriyi" (data-at-rest)** korur. **"Kullanımdaki veriyi" (data-in-use)** korumak senin operasyonel disiplinine bağlıdır (auto-dismount, ekran kilidi, temiz uç nokta, anti-keylogger).

---

<a id="3"></a>
## 3. 📥 Kurulum + İmza Doğrulama

**Asla** rastgele bir aynadan indirme. Yalnızca **veracrypt.eu / veracrypt.fr** (resmi) veya GitHub (`veracrypt/VeraCrypt`).

### İmza doğrulama (ATLAMA — kurulumun en kritik adımı)
İndirilen kurulum dosyası kötü amaçlıyla değiştirildiyse, *kurduğun anda* parolan çalınır. Bu yüzden PGP imzasını doğrula:

```bash
# 1. VeraCrypt PGP genel anahtarını içe aktar (resmi siteden parmak izini teyit et)
gpg --import VeraCrypt_PGP_public_key.asc

# 2. İmzayı doğrula
gpg --verify "VeraCrypt Setup 1.26.x.exe.sig" "VeraCrypt Setup 1.26.x.exe"

# "Good signature from IDRIX" ve doğru parmak izi görmelisin
```

- **Windows:** Yükleyici imzasını da kontrol et (Özellikler → Dijital İmzalar → IDRIX SARL).
- **SHA-256:** Resmi sitedeki hash ile karşılaştır.
- **Portable mod:** Kurulum yerine "Extract" seçeneğiyle taşınabilir kullan (sistemde iz bırakmaz; aşağıda "iz" bölümüne bak).

> 🔑 **Püf:** İmza doğrulamayı *her sürüm güncellemesinde* tekrarla. Tedarik zinciri saldırıları en çok "güncelleme" anında işler.

---

<a id="4"></a>
## 4. 🧱 Temel Kavramlar

### Üç tür hedef
1. **File Container (kapsayıcı dosya):** `gizli.hc` gibi tek bir dosya. Mount edilince sanal sürücü (örn. `X:`) olur. **En esnek, en taşınabilir** — yeni başlayan için ideal, hidden volume için temel.
2. **Non-system Partition/Drive:** Tüm bir USB/disk bölümü şifrelenir. Dosya değil, ham bölüm.
3. **System Encryption:** Windows'un kurulu olduğu sistem diskini şifreler → açılışta **pre-boot parola** sorar.

### Anahtar kavramı: XTS modu
VeraCrypt disk şifrelemede **XTS** kullanır (eski LRW/CBC değil). XTS, her sektörü konum-bağımlı şifreler → "watermark" ve sektör değiştirme saldırılarına dayanıklıdır. Sen bir şey yapmazsın; sadece "XTS doğru moddur, başkasını seçme" diye bil.

### Anahtarın yolculuğu (basitleştirilmiş)
```
Parola + (Keyfile) + PIM
        │
        ▼  PBKDF2 (yüz binlerce iterasyon, seçtiğin hash)
   Header Anahtarı
        │
        ▼  Volume Header'ı çözer (ilk 64 KB)
   Master Anahtar (XTS) ──► Veriyi anında çözer/şifreler (yalnızca RAM'de)
```
**Kritik sonuç:** Master anahtar **header içinde** şifreli durur. Header bozulursa/yanlış silinirse → **parolan doğru olsa bile veri SONSUZA dek erişilemez.** (Bkz. Bölüm 13 — Header Yedekleme.)

---

<a id="5"></a>
## 5. 🔐 Şifreleme & Hash Algoritması Seçimi

### Şifreleme algoritmaları
| Seçim | Açıklama | Ne zaman |
|---|---|---|
| **AES** | Donanım hızlandırmalı (AES-NI), çok hızlı, NIST standardı | Varsayılan, günlük kullanım |
| **Serpent** | AES finalisti, daha geniş güvenlik marjı, daha yavaş | AES'e güvenmeyen paranoyak |
| **Twofish** | Bruce Schneier; sağlam, orta hız | Alternatif |
| **Camellia / Kuznyechik** | Japon / Rus standartları | Niş |
| **Kaskadlar** (AES-Twofish, AES-Twofish-Serpent…) | Üç bağımsız şifre üst üste | Maksimum paranoya |

> 🔥 **Püf — kaskad gerçeği:** `AES(Twofish(Serpent()))` üçlü kaskadı, AES *tamamen* kırılsa bile diğer iki katmanın koruması demektir. **Ama** performans ~3 kat düşer ve AES-NI avantajını kaybedersin. Pratikte **AES tek başına bugün kırılmadı**; kaskad ancak "AES'te gizli kusur var" tehdidine inanıyorsan mantıklıdır. Devlet-seviyesi rakip + uzun ömürlü sır için: **Serpent veya AES-Serpent** iyi bir denge (Serpent'in güvenlik marjı AES'ten geniştir).

### Hash / KDF seçimi
| Hash | Not |
|---|---|
| **SHA-512** | Güçlü, yaygın, non-system için mükemmel varsayılan |
| **Whirlpool** | Sağlam alternatif |
| **BLAKE2s-256** | Hızlı, modern |
| **SHA-256** | Sistem şifrelemede kullanılır (boot ortamı kısıtlı) |
| **Streebog** | Rus GOST — NSA/FSB paranoyası olanlar genelde kaçınır |

> 🧠 **CTI önerisi:** Non-system için **Serpent + Whirlpool** (ya da AES + SHA-512); sistem şifreleme için boot uyumluluğu nedeniyle pratikte **AES + SHA-256**. Streebog/Kuznyechik'i, belirli bir devlet aktörüne güvenmiyorsan bilinçli seç/kaçın.

---

<a id="6"></a>
## 6. ⚙️ PIM Derinlemesine — Gizli İkinci Faktör

**PIM (Personal Iterations Multiplier)**, KDF'nin kaç kez döneceğini belirler. Anlamadan kullanan çok; ustaca kullanan az.

### Matematiği (sürümle değişebilir — kendi sürümünden teyit et)
- **Non-system, PIM = 0 (boş):** İterasyon = **500.000** (SHA-512/Whirlpool/BLAKE2s/Streebog).
- **Non-system, PIM belirtilmiş:** İterasyon = **15.000 + (PIM × 1.000)**.
- **System, PIM = 0:** İterasyon = **200.000** (SHA-256 ailesi).
- **System, PIM belirtilmiş:** İterasyon = **PIM × 2.048**.

### Neden önemli — iki yönlü kılıç
- **Yüksek PIM** → brute-force *katlanarak* yavaşlar (saldırgan her deneme için aynı yavaşlığı yaşar) **ama** senin mount sürenin de uzar (zayıf CPU'da 10-30 sn).
- **Düşük PIM (örn. system için <98)** → mount hızlı **ama** kısa parolayla birleşince zayıf.

### 🔥 PIM'i "gizli ikinci faktör" yap
PIM, parolayla birlikte **bilinmesi gereken ikinci sayıdır.** Saldırgan doğru parolayı bulsa bile **yanlış PIM ile asla açamaz** ve PIM'i bilmediği için her denemede *tüm makul PIM aralığını* taraması gerekir → arama uzayını devasa büyütür.

- Parolan güçlü değilse (uzun değilse), **özel ve akılda kalır bir PIM** seç (örn. doğum yılın değil — tahmin edilebilir; rastgele 3 haneli bir sayı).
- PIM'i **asla parolayla aynı yerde yazma.** Ayrı sakla ya da ezberle.
- **Sistem şifrelemede düşük PIM tehlikesi:** PIM < 98 ise iterasyon çok düşer; yalnızca **çok uzun** parola ile birleştir.

> 💡 **Operasyonel:** "Uzun passphrase + makul yüksek PIM" kombinasyonu, "orta parola + PIM yok"tan kat kat güçlüdür. PIM'i unutursan birim açılmaz — ezber + güvenli yedek şart.

---

<a id="7"></a>
## 7. 🗝️ Keyfile Derinlemesine

Keyfile, "bildiğin bir şey" (parola) yanına "sahip olduğun bir şey" (dosya) ekler → gerçek iki faktör.

### Nasıl çalışır
VeraCrypt, keyfile'ın **ilk 1.048.576 baytını (1 MB)** okur, CRC-32 ile işler ve parolayla karıştırır. Yani 1 MB'tan büyük dosyada sadece ilk 1 MB önemlidir.

### 🔥 Püf noktaları
- **Keyfile olarak ne seç:** Büyük, **asla değişmeyecek** bir dosya (bir fotoğraf, bir PDF). **Değişirse birim açılmaz!** Bu yüzden düzenlenen dosyaları (Word, log) ASLA keyfile yapma.
- **Daha iyisi:** VeraCrypt ile **rastgele keyfile üret** (Tools → Keyfile Generator) → tahmin/üretim imkânsız.
- **Nerede tut:** Keyfile'ı **ayrı bir fiziksel ortamda** (küçük USB, hatta NFC/akıllı kart) sakla → disk çalınsa keyfile yanında olmaz = veri açılamaz. Bu, "diski çaldım ama açamıyorum" senaryosunu yaratır.
- **Gizleme:** Keyfile'ı masum bir dosyanın içine göm (steganografi) ya da binlerce dosya arasına serpiştir — saldırgan hangisinin keyfile olduğunu bilemez.
- **Çoklu keyfile:** Birden çok keyfile zorunlu kılınabilir (hepsi gerekir) → "anahtarı 3 kişiye böl" senaryosu.
- **Ölümcül risk:** Keyfile'ı kaybedersen/değiştirirsen **veri gider.** Mutlaka güvenli bir kopyasını ayrı tut (ama kopyası da güvenli olmalı).

> 🧠 **En güçlü kombinasyon:** Uzun passphrase **+** rastgele keyfile (ayrı USB'de) **+** özel PIM. Üçünü birden ele geçirmek pratikte imkânsızdır.

---

<a id="8"></a>
## 8. 📦 Standart Şifreli Kapsayıcı (Adım Adım)

İlk birimini güvenle oluştur:

1. **Create Volume → Create an encrypted file container.**
2. **Standard VeraCrypt volume** (hidden'a sonra geçeceğiz).
3. **Volume Location:** Dosya adını masum seç. **Uzantı ipucu vermez** — `tatil_2019.mp4` gibi bir ad, kapsayıcıyı medya dosyası gibi gizler (boyut tutarlı olmalı).
4. **Encryption Options:** AES (veya Serpent/AES-Serpent) + SHA-512.
5. **Volume Size:** Gerçekçi seç. (Hidden volume kullanacaksan dış birimi yeterince büyük yap.)
6. **Volume Password:** **Uzun passphrase** (20+ karakter, rastgele kelime dizisi — "diceware" yöntemi mükemmeldir). İstersen keyfile + PIM ekle.
7. **Filesystem:** 
   - **exFAT** → platformlar arası (Win/mac/Linux), büyük dosya destekler.
   - **NTFS** → yalnızca Windows ağırlıklıysan.
   - **FAT** → küçük birimler, en uyumlu.
8. **🔥 Mouse Entropy (kritik):** Pencerede fareyi **rastgele ve uzun süre** gezdir → anahtar üretimi için entropi toplar. **Ne kadar uzun, o kadar güçlü rastgelelik.** Aceleyle geçme; bu, master anahtarın kalitesini doğrudan belirler.
9. **Format:** 
   - **Quick Format KULLANMA** (özellikle hidden volume planlıyorsan) → tüm alanı rastgele veriyle doldur ki boş alan "şifreli veri" mi yoksa "boşluk" mu ayırt edilemesin (deniability için şart).
   - **Dynamic (sparse) KULLANMA** → büyürken disk üzerinde "küçük dosya" gibi görünür, deniability'yi ve performansı bozar.
10. **Mount:** Select File → Slot seç → Mount → parola/keyfile/PIM gir → `X:` sürücüsü hazır.
11. **Dismount:** İş bitince **Dismount** (ya da Dismount All). Bilgisayardan kalkarken **mutlaka** dismount et.

---

<a id="9"></a>
## 9. 🥷 Gizli Birim (Hidden Volume) & İnkâr Edilebilirlik

VeraCrypt'in **kalbi** budur. Başka hiçbir yaygın araç bunu sunmaz.

### Mantık
Bir **dış (outer) birim** oluşturursun. Onun **boş alanının içine**, ikinci bir **gizli (hidden) birim** gömülür. İki ayrı parola:
- **Dış parola** → masum/yem veriyi gösterir.
- **Gizli parola** → asıl sırrı açar.

```
┌─────────────────────────────────────────────┐
│  DIŞ BİRİM (yem parola ile açılır)            │
│  ┌──────────────┐                             │
│  │ Yem veriler  │   ...boş görünen alan...    │
│  │ (eski foto,  │   ┌───────────────────────┐ │
│  │  önemsiz pdf)│   │ GİZLİ BİRİM           │ │
│  └──────────────┘   │ (asıl sır — gizli     │ │
│                     │  parola ile açılır)   │ │
│                     └───────────────────────┘ │
└─────────────────────────────────────────────┘
```

Gizli birimin verisi, dış birimin "boş alanındaki rastgele gürültüden" **matematiksel olarak ayırt edilemez.** Saldırgan dış parolayı alsa bile gizli birimin **varlığını kanıtlayamaz.**

### Oluşturma
Create Volume → **Hidden VeraCrypt volume** → "Normal mode" → önce dış, sonra gizli birim. Her ikisine **farklı, güçlü** parola ver.

### 🔥 İnkâr edilebilirliği AYAKTA tutan kurallar (en çok burada hata yapılır)
1. **Dış birime yem veri koy — ve inandırıcı olsun.** Tamamen boş dış birim şüphe çeker. Ama içine "korunmaya değer ama yıkıcı olmayan" veri koy (eski vergi belgeleri, sıradan fotoğraflar). Banka/kripto gibi "asıl değerli" görünen şeyi koyma — saldırgan "asıl sır bu değil" diye daha çok kazar.
2. **Dış birimi ara sıra güncelle.** Timestamp'ler 3 yıl öncesinden kalmışsa "bu sadece bir paravan" belli olur.
3. **DIŞ BİRİME YAZARKEN GİZLİYİ EZME!** Dış birime normal mount edip veri yazarsan, dosya sistemi gizli birimin oturduğu alana yazabilir → **gizli birim kalıcı bozulur.** Çözüm: Dış birimi mount ederken **Mount Options → "Protect hidden volume against damage"** işaretle ve **gizli parolayı da gir** → VeraCrypt gizli alanı korur. (Bu modda yazarken gizli birim yazma-korumalı olur.)
4. **Asla iki parolayı aynı yerde tutma.** İnkâr edilebilirliğin tüm anlamı, saldırganın gizli parolayı *bilmemesidir*.
5. **Uygulama izlerini temizle (Bölüm 12).** "Recent files"ta gizli birimin sürücü harfi/yolu görünürse deniability çöker.

> ⚖️ **Hukuki amaç:** Bazı ülkelerde parola vermeyi reddetmek suçtur (Bölüm 19). Hidden volume burada hayat kurtarır: **dış parolayı verirsin** ("işte parolam"), gizli birim görünmez kalır.

---

<a id="10"></a>
## 10. 💻 Sistem Şifreleme (Pre-Boot Authentication)

Tüm Windows sistem diskini şifreler. Açılışta, Windows yüklenmeden **önce** VeraCrypt bootloader parola sorar.

### Süreç
System → **Encrypt System Partition/Drive** → Normal → şifre seçimi → **Rescue Disk oluştur (zorunlu — ATLAMA)** → wipe mode → **pretest** (sistem önce parolayı test eder, çalışırsa şifrelemeye başlar) → arka planda şifreleme (kullanmaya devam edebilirsin).

### Kritik noktalar
- **Rescue Disk (kurtarma diski):** Bootloader/header bozulursa **tek kurtuluşun.** ISO'yu fiziksel CD'ye yaz ya da güvenli sakla. Kaybedersen ve header bozulursa → **sistem açılmaz, veri gider.**
- **Pre-boot hash = SHA-256:** Boot ortamı kısıtlı olduğundan tüm hash'ler desteklenmez.
- **Klavye düzeni:** Pre-boot ortamı genelde **ABD/QWERTY** varsayar. Parolanda Türkçe karakter veya düzen-bağımlı semboller varsa açılışta yazamayabilirsin → **parolayı yalnızca ASCII harfler/rakamlar/standart sembollerden** seç.
- **Çift kimlik doğrulama:** PIM ve/veya keyfile sistem şifrelemede de kullanılabilir (keyfile pre-boot için USB üzerinden).

---

<a id="11"></a>
## 11. 🕳️ Gizli İşletim Sistemi (Hidden OS) — En Üst Seviye

İnkâr edilebilirliğin zirvesi. **İki ayrı Windows kurulumu:**
- **Decoy OS (yem):** Normal görünür, masum kullanılır, parolası verilebilir.
- **Hidden OS (gizli):** Asıl çalıştığın sistem; varlığı kanıtlanamaz.

```
Disk:
┌────────────────────┬───────────────────────────────────┐
│ Sistem bölümü      │ İkinci bölüm                       │
│ (DECOY OS — yem)   │ ┌─────────────┬───────────────────┐│
│ açılış parolası 1  │ │ Dış birim   │ HIDDEN OS         ││
│                    │ │ (yem veri,  │ (gizli sistem,    ││
│                    │ │  parola 2)  │  parola 3)        ││
│                    │ └─────────────┴───────────────────┘│
└────────────────────┴───────────────────────────────────┘
```

Açılışta hangi parolayı girdiğine göre **decoy** ya da **hidden** OS açılır. Zorlanırsan decoy parolasını verirsin; gizli sistem görünmez.

> ⚠️ **Disiplin şart:** Decoy OS'i **gerçekten kullan** (gezinti geçmişi, belgeler olsun) ki inandırıcı olsun. Hassas işleri **yalnızca hidden OS'te** yap. Bu, en güçlü ama en çok operasyonel titizlik isteyen moddur. Kurulum karmaşıktır — resmi dokümanı adım adım izle.

---

<a id="12"></a>
## 12. 🔥 PÜF NOKTALARI — Piyasada Bulamayacakların

Bu bölüm, çoğu rehberin atladığı ve gerçek dünyada deniability'yi/güvenliği **çökerten** detaylardır.

### 12.1 SSD / TRIM / Wear-Leveling — Hidden Volume'un Baş Düşmanı
- SSD'ler **wear-leveling** yapar: veriyi fiziksel olarak gezdirir → "boş alan rastgele gürültü" varsayımı bozulabilir, gizli birimin varlığına dair **istatistiksel iz** kalabilir.
- **TRIM** komutu, "bu bloklar boş" bilgisini SSD'ye söyler → dış birimde boş görünen ama aslında gizli birim olan alan ele verilebilir. VeraCrypt **sistem şifrelemede TRIM'i deniability için kapatır** (performans/ömür pahasına).
- **Sonuç:** **Hidden volume için ideal ortam, SSD değil mekanik HDD'dir** (ya da TRIM'i bilinçli yönet). Kritik inkâr edilebilirlik gerekiyorsa bunu hesaba kat.

### 12.2 Hibernate / Hazırda Bekletme — Anahtarı Diske Sızdırır
- **Hazırda bekletme (hibernate)**, RAM'in tamamını `hiberfil.sys`'e yazar — **master anahtar dahil.** Non-system birim mount'ken hibernate edersen, anahtar şifresiz diske düşebilir.
- **Uyku (sleep)** ise RAM'i açık tutar → cold-boot riski.
- **Çözüm:** Hassas birim mount'ken **hibernate/sleep kullanma.** Ya tüm sistemi şifrele (hiberfil de şifrelenir) ya hibernate'i tamamen kapat:
  ```
  powercfg -h off
  ```

### 12.3 Cold-Boot Saldırısı — RAM'deki Anahtar
- RAM, güç kesildikten sonra saniyeler-dakikalar (soğutulursa daha uzun) içeriğini tutar. Saldırgan açık/uykudaki makineden RAM'i dondurup anahtarı çekebilir.
- **Çözüm:**
  - İş biter bitmez **Dismount** (anahtar RAM'den silinir).
  - VeraCrypt: **Settings → Preferences → "Dismount all when: user logs off / screensaver / power saving"** + **"Wipe cached passwords on dismount/exit."**
  - Bırakacaksan **tamamen kapat** (uyku/hazırda değil).
  - Kanije Kalesi entegrasyonu (Bölüm 16) burada devreye girer: fiziksel tehditte otomatik dismount.

### 12.4 Page File (Swap) Sızıntısı
- Windows page file (`pagefile.sys`), RAM'den taşan veriyi diske yazar — mount'ken hassas veri parçaları sızabilir.
- **Çözüm:** Tüm sistemi şifrele (pagefile de şifreli olur) ya da page file'ı şifreli birimde tut / minimize et.

### 12.5 Header Yedeği — "Parolam doğru ama açılmıyor" kâbusu
- İlk 64 KB'lık **volume header**, master anahtarı şifreli tutar. Tek bir bozulma (kötü sektör, yanlış yazma, format kazası) → **parola doğru olsa bile veri sonsuza dek gider.**
- **Çözüm:** **Tools → Backup Volume Header** ile her birimin header'ını yedekle, **ayrı ve güvenli** sakla. (Yedek header, eski parolayı da içerir — parola değiştirdiysen yedeği güncelle, eskisini güvenli imha et.)

### 12.6 Adli Bilişim Karşıtı — İz Temizleme
Şifrelemen mükemmel olsa bile **işletim sistemi seni ele verir:**
- **Son kullanılanlar / MRU:** VeraCrypt'in "recently mounted volumes" listesi → **Settings'ten temizle** ve **"Never save history"** işaretle.
- **Windows Recent / Jump Lists:** Mount edilen sürücüdeki dosya yolları kayıt altına alınır.
- **Shellbags, Prefetch, $MFT, USN Journal:** Dosya adları/erişim izleri bırakır.
- **Thumbnail cache:** Açtığın görsellerin küçük resimleri `thumbcache`'te kalır.
- **Uygulama izleri:** Office/PDF okuyucu "son açılan dosya" tutar.
- **Çözüm:** Hassas işleri **hidden OS** ya da **amnezi sistemi (Tails)** üzerinde yap; mümkünse şifreli birimdeki dosyaları dışarı kopyalama. VeraCrypt'i **portable** çalıştır (kurulum izi bırakmaz).

### 12.7 Dosya Sistemi Journaling & Defrag
- NTFS journaling, dış birime yazarken metadata bırakır; defrag dosyaları gezdirip gizli birime zarar verebilir.
- **Hidden volume'lu dış birimi defrag ETME.** Gerekiyorsa exFAT/FAT gibi journaling'siz dosya sistemi deniability için daha sessizdir (ama dayanıklılık daha az — denge).

### 12.8 Kapsayıcıyı Gizleme (Steganografik Adlandırma)
- Kapsayıcı dosyanın **uzantısı/adı** önemli ipucu verir. `gizli.hc` → herkes anlar. `2018_dugun.mp4` (gerçekçi boyutla) → medya dosyası sanılır.
- VeraCrypt kapsayıcıları "yüksek entropili" görünür (rastgele); gelişmiş analiz "bu dosya şifreli/sıkıştırılmış" diyebilir ama **içeriği** asla çözemez. Tam gizlilik için steganografi araçlarıyla başka katman ekleyebilirsin.

### 12.9 Birden Çok Parola Değişimi & "Anahtar Hijyeni"
- Parolayı değiştirmek master anahtarı **değiştirmez**, yalnızca header'ı yeni parolayla yeniden şifreler. Yani eski header yedeği hâlâ eski parolayla açar → **parola sızdıysa, eski header yedeklerini de imha et.**
- Gerçekten "anahtarı döndürmek" istiyorsan → **yeni birim oluştur, veriyi taşı, eskisini güvenli imha et.**

### 12.10 Mount Sırasında Yanlış Slot/Sürücü
- Aynı anda çok birim mount ediyorsan, yanlış sürücüye yazma riski. Her birime tanınır bir etiket ver; dismount disiplinini koru.

---

<a id="13"></a>
## 13. 💾 Yedekleme & Kurtarma

Şifreli veride **yedek = hayat.** Ama yedek de şifreli olmalı.

- **Header yedeği:** Her birim için (Bölüm 12.5). Header bozulmasının tek kurtarıcısı.
- **Rescue Disk (sistem şifreleme):** Fiziksel CD/USB; bootloader kurtarır.
- **Veri yedeği:** Şifreli birimi **bir bütün olarak** (kapsayıcı dosyayı) başka şifreli ortama kopyala. Asıl ile yedek arasında parola/keyfile aynıysa, ikisini de aynı anda kaybetmeyecek şekilde **coğrafi olarak ayır** (3-2-1 kuralı: 3 kopya, 2 farklı ortam, 1 saha dışı).
- **Hidden volume yedeği:** Dış birimi olduğu gibi kopyala → gizli birim de içinde taşınır (ayrıca yedeklemeye gerek yok, ama kopyalama "protect hidden" gerektirmez çünkü ham kopyadır).
- **Test et:** Yedeğini **düzenli olarak farklı makinede mount ederek doğrula.** Açılmayan yedek, yedek değildir.

---

<a id="14"></a>
## 14. 🛡️ Günlük Operasyon Güvenliği

- **Auto-dismount:** Settings → Preferences → idle/logoff/screensaver/power'da otomatik dismount.
- **Wipe cached passwords:** Çıkışta/dismount'ta parola önbelleğini temizle.
- **Hotkey:** "Dismount All" için global kısayol ata → tehlikede tek tuşla her şeyi kapat.
- **Ekran kilidi disiplini:** Masaüstünden kalkarken `Win+L`. Mount'ken kilitsiz bırakma.
- **Temiz uç nokta:** Keylogger/RAT bulaşmış makinede şifreleme çaresizdir → ayrı, sertleştirilmiş bir sistem kullan (bkz. repo'daki `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`).
- **Read-only mount:** Yalnızca okuyacaksan birimi salt-okunur bağla → yanlışlıkla yazma/deniability hasarı yok.

---

<a id="15"></a>
## 15. ⌨️ Komut Satırı & Otomasyon

VeraCrypt CLI (`veracrypt.exe` / `veracrypt`) script'lenebilir:

```bash
# Bağla (Windows)
"C:\Program Files\VeraCrypt\VeraCrypt.exe" /q /v "D:\gizli.hc" /l X /p "PAROLA" /pim 487 /keyfile "E:\anahtar.bin"

# Çöz (dismount)
"C:\Program Files\VeraCrypt\VeraCrypt.exe" /q /d X

# Tümünü çöz
"C:\Program Files\VeraCrypt\VeraCrypt.exe" /q /d
```

> 🔥 **Püf:** Parolayı **komut satırına düz yazma** — komut geçmişi (`PSReadLine`/bash history), process listesi ve ETW ile sızar. Script'lerde parolayı *etkileşimli* sor (parametresiz bırak, VeraCrypt güvenli pencerede ister) ya da otomasyonu yalnızca keyfile + güvenli prompt ile kur. CLI'ı en çok **"tehlikede tek komutla dismount"** için kullan.

---

<a id="16"></a>
## 16. 🏰 Kanije Kalesi ile Birlikte Kullanım

Bu repo (Kanije Kalesi), **fiziksel tehdit anında** cihazı uzaktan/otomatik koruyan bir muhafızdır. VeraCrypt **durağan veriyi** korur; Kanije **olay anını** yönetir. İkisi katmanlı bir savunma oluşturur:

| Senaryo | VeraCrypt rolü | Kanije Kalesi rolü |
|---|---|---|
| Cihaz çalındı (kapalı) | Veri şifreli, erişilemez | — |
| Cihaz açık, sen uzaktayken biri yaklaştı | Birim mount → veri açık (risk) | `/koruma` dead-man / USB / yanlış-giriş tetikleyici → **kilitle + alarm + foto** |
| Acil durum | — | `/panik`, `/kilit tam` (lockdown) |
| Verinin kalıcı yok edilmesi | Birim header'ı sil → veri matematiksel olarak gider | `/imha` ile hedef dosyaları güvenli sil |
| Adli inceleme öncesi | Hidden volume inkâr edilebilirlik | RAM-only mod, iz temizleme |

### 🔥 Önerilen entegrasyon deseni
1. **Hassas veriyi VeraCrypt hidden volume'da tut** (mekanik diskte, deniability için).
2. **Kanije `/koruma`'yı aç** — dead-man switch + USB dead-man. Sen X saat check-in yapmazsan ya da USB çıkarılırsa Kanije devreye girer.
3. **VeraCrypt auto-dismount'u** "screensaver/power/logoff"a bağla → Kanije kilitlerken birim de dismount olur, RAM'deki anahtar uçar.
4. **Hızlı yok etme:** En kritik senaryoda, VeraCrypt birim **header'ını** silmek (64 KB) tüm veriyi anında erişilemez kılar — Kanije'nin `/imha` hedeflerine birim header konumunu eklemek, "veriyi sil"i saniyeler içinde yapar (gigabaytları silmeye gerek yok; **anahtarı yok et, veri çöp olur**).
5. **RAM-only çalışma:** Kanije'nin RAM-only modu + VeraCrypt'in dismount disiplini = güç kesildiğinde **hiçbir iz kalmaz** (Tails/BusKill sınıfı duruş).

> 🧠 **Felsefe örtüşmesi:** Kanije'nin `/imha`'sı tüm dosyayı silmek yerine **anahtarı/header'ı yok etmeye** odaklanırsa (crypto-shredding), terabaytlık veri **bir anda** kurtarılamaz hale gelir. VeraCrypt bu yaklaşımın temelidir.

---

<a id="17"></a>
## 17. ☠️ Yaygın Ölümcül Hatalar

1. **Header yedeği almamak** → tek bozulmada tüm veri gider. (En sık felaket.)
2. **Quick Format / Dynamic kapsayıcı** kullanıp deniability'yi baştan öldürmek.
3. **Hidden volume'lu dış birime "protect hidden" olmadan yazıp gizliyi ezmek.**
4. **SSD'de hidden volume** kullanıp TRIM/wear-leveling ile varlığı ele vermek.
5. **Mount'ken hibernate/uyku** → anahtarı diske/RAM'e sızdırmak.
6. **Pre-boot parolada Türkçe/özel karakter** → açılışta yazamamak.
7. **Keyfile'ı değişen bir dosya yapmak** → bir gün dosya değişince birim açılmaz.
8. **Parolayı PIM/keyfile ile aynı yerde saklamak** → ikinci faktörü anlamsızlaştırmak.
9. **İmza doğrulamadan kurmak** → truva atı yükleyiciyle parolayı sızdırmak.
10. **Recent/MRU izlerini bırakmak** → şifre mükemmel ama işletim sistemi seni ele verir.
11. **Tek kopya** → disk ölünce hem veri hem şifre gider (yedek = şifreli olmalı).

---

<a id="18"></a>
## 18. ✅ Hızlı Referans & Operasyonel Kontrol Listesi

### Birim oluştururken
- [ ] Yükleyici **imzası doğrulandı**
- [ ] AES (veya Serpent/kaskad) + SHA-512
- [ ] **Uzun passphrase** (diceware) + **keyfile** (ayrı USB) + **özel PIM**
- [ ] Mouse entropi **bol bol** toplandı
- [ ] **Quick Format YOK**, **Dynamic YOK**
- [ ] Dosya adı **masum** (medya gibi)

### Hidden volume için ek
- [ ] Dış birimde **inandırıcı yem veri** var
- [ ] Dış birime yazarken **"Protect hidden volume"** açık
- [ ] **Mekanik disk** (SSD ise TRIM riskini kabul ettim)
- [ ] İki parola **ayrı yerlerde**

### Her oturum
- [ ] Temiz/sertleştirilmiş uç nokta
- [ ] İş bitince **Dismount**
- [ ] Kalkarken **ekran kilidi**
- [ ] **Hibernate/uyku yok** (mount'ken)

### Periyodik
- [ ] **Header yedeği** güncel ve güvende
- [ ] Yedek **farklı makinede test** edildi
- [ ] **MRU/Recent izleri** temiz
- [ ] Parola değiştiyse **eski header yedeği imha** edildi

---

<a id="19"></a>
## 19. ⚖️ Hukuki Sınır & Sınır Geçişi Notu

- **Parola zorunluluğu:** Bazı ülkelerde (ör. Birleşik Krallık RIPA Bölüm 49, Fransa) mahkeme/sınır görevlisi parola talep edebilir; reddetmek **suç** sayılabilir. Hidden volume tam bu yüzden vardır: **dış parolayı verirsin**, gizli birim görünmez kalır ve "başka bir şey yok" beyanın çürütülemez.
- **Sınır geçişi:** Cihazlar sınırda kopyalanabilir/incelenebilir. Hassas veriyi **fiziksel olarak yanında taşıma**; gizli/şifreli bir bulut üzerinden, geçtikten sonra indir. Ya da yalnızca **decoy OS** ile seyahat et.
- **İnkâr edilebilirlik mükemmel değildir:** Yeterince kaynaklı bir adli ekip, *davranışsal* ve *istatistiksel* ipuçları (boş alan dağılımı, kullanım izleri, SSD davranışı) arayabilir. Operasyonel disiplin (Bölüm 12) bu ipuçlarını yok etmek içindir.
- **Yasal sorumluluk:** Bu rehber **meşru gizlilik, gazetecilik, araştırma ve kişisel veri koruma** içindir. Bulunduğun yargı bölgesinin yasalarına uy.

---

> 🏰 **Kapanış:** Şifreleme bir ürün değil, bir **disiplindir.** En güçlü AES bile, dismount etmeyi unuttuğun bir akşam ya da keylogger'lı bir makinede çaresizdir. VeraCrypt sana matematiksel kaleyi verir; **kapıyı kilitli tutmak senin işin.** Kanije Kalesi de tam burada — sen unuttuğunda kapıyı kilitleyen nöbetçi — devreye girer.
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. İlgili: `SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md`, `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`, `DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md`.*
