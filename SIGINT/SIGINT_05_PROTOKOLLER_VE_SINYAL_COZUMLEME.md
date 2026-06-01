# SIGINT EL KİTABI — BÖLÜM 5
## Protokoller ve Sinyal Çözümleme — Yakaladığın Sinyali Anlamak

> **Bu bölüm neyi öğretir:** Önceki bölümlerde SDR'ını kurdun, anteni taktın, waterfall'da renkli şeritler gördün. Şimdi asıl soru: **"Bu çizgi ne anlatıyor?"** Bir sinyali yakalamak ile onu **çözmek** (decode) arasında dağlar kadar fark vardır. Bu bölüm, bilinmeyen bir sinyali tanımlama metodolojisini, ardından dünyanın her yerinde meşru olarak dinlenebilen **açık protokolleri** uçtan uca (frekans → araç → komut → örnek çıktı) öğretir. Sonunda hücresel haberleşmenin (2G/3G/4G/5G) **nasıl çalıştığını** prensip düzeyinde açıklar — çünkü bir savunmacı, savunduğu şeyin mimarisini bilmek zorundadır.

> **ÖNCE BUNU OKU — YASAL ÇERÇEVE (atlamak yok):**
> - **Açık/meşru yayınlar serbesttir.** ADS-B (uçak), AIS (gemi), NOAA hava durumu uyduları, amatör radyo (APRS/FT8), rtl_433 (kendi ev sensörlerin) — bunlar **kamuya açık, şifresiz, herkesin dinlemesi için tasarlanmış** yayınlardır. Almak, çözmek, haritalamak yasaldır ve harika bir mühendislik alıştırmasıdır.
> - **ÖZEL/hücresel haberleşmeyi yetkisiz dinlemek/çözmek SUÇTUR.** Başkasının telefon görüşmesi, SMS'i, mobil verisi, şifreli telsizi, çağrı cihazı içeriği — bunları izinsiz dinlemek/çözmek/kaydetmek Türkiye'de TCK ve ilgili mevzuat, dünyada hemen her ülkede ağır cezalı suçtur. Yalnızca yakalamak bile suç sayılabilir.
> - **Bu bölümün hücresel kısmı bir saldırı reçetesi DEĞİLDİR.** Sana "şu adımları yap, komşunun telefonunu dinle" demiyoruz ve demeyeceğiz. Mimariyi, neden bazı nesillerin savunmasız bazılarının sağlam olduğunu, ve **kendini nasıl koruyacağını** anlatıyoruz. İşin **yasal yolu** = test/araştırma ortamı, **kendi** baz istasyonun (lisanslı veya Faraday kafesinde), operatör/üniversite izni.
> - Şüphedeyken: **alıcı (RX) olarak meşru yayınları dinle, dur.** Verici (TX) ve özel trafik = ayrı dünya.

---

## İÇİNDEKİLER

1. [Sinyal Tanımlama Metodolojisi — "Bu da ne?"](#1)
2. [Sinyalin Beş Parmak İzi (Merkez, Bant, Modülasyon, Sembol Hızı, Yapı)](#2)
3. [Referans Kütüphaneleri — sigidwiki ve Arkadaşları](#3)
4. [ ADS-B — Uçakları Gökyüzünden Çek (1090 MHz)](#4)
5. [ AIS — Gemileri Denizden Çek (162 MHz)](#5)
6. [ACARS — Uçakların Metin Mesajları (131 MHz)](#6)
7. [POCSAG / FLEX — Çağrı Cihazları (Pager) ](#7)
8. [APRS — Amatör Paket Radyo (144.800 MHz)](#8)
9. [ NOAA APT & Meteor-M LRPT — Uydudan Hava Durumu Görüntüsü](#9)
10. [ rtl_433 — Evdeki Her Şeyin Telsizi (433/868 MHz)](#10)
11. [LoRa — Uzun Menzilli IoT (868/915 MHz)](#11)
12. [Amatör Dijital Modlar — FT8 / PSK31 / RTTY](#12)
13. [Sinyal Tersine Mühendislik — URH ile KENDİ Kumandanı Çöz](#13)
14. [ HÜCRESEL — 2G/3G/4G/5G Mimarisi (Prensip + Yasal Sınır)](#14)
15. [Net Çerçeve — Gözlem mi, Dinleme mi? (Yasal/Yasadışı Ayrımı)](#15)
16. [ ALIŞTIRMALAR — Hepsi Evde, Hepsi Yasal](#16)
17. [Hızlı Referans Tablosu (Protokol × Frekans × Araç)](#17)
18. [Çapraz Referans & Sonraki Bölüm](#18)

---

<a id="1"></a>
## 1.  Sinyal Tanımlama Metodolojisi — "Bu da ne?"

Waterfall'da gezinirken er ya da geç şu an gelir: tanımadığın bir desen. Belki ritmik atan bir nokta dizisi, belki geniş bir "fırça darbesi", belki düzenli aralıklarla yanıp sönen iki çizgi. Acemi panikler ("hacklendim mi?"); usta ise **sistematik bir kimlik tespiti** akışına girer.

Sinyal çözümleme bir **eleme oyunudur**. Eline aldığın her ipucu (frekans, genişlik, ritim) olasılıkları daraltır. Hiç kimse bir sinyale bakıp "bu ADS-B" diye anında bilmez — **özelliklerini ölçer, referansla eşleştirir, doğru aracı dener.**

### "Tuhaf bir sinyal buldum" akışı (zihinsel kontrol listesi)

```
   ┌─────────────────────────────────────────────────────┐
   │  1. NEREDE?  Merkez frekansı tam olarak oku          │
   │     (örn. 1090.000 MHz, 162.025 MHz)                 │
   └───────────────────────┬─────────────────────────────┘
                           ▼
   ┌─────────────────────────────────────────────────────┐
   │  2. NE KADAR GENİŞ?  Bant genişliğini ölç            │
   │     (kHz mi, MHz mi? dar = az veri, geniş = çok)     │
   └───────────────────────┬─────────────────────────────┘
                           ▼
   ┌─────────────────────────────────────────────────────┐
   │  3. NASIL TİTREŞİYOR?  Modülasyonu tahmin et         │
   │     (AM/genlik mi? FM/frekans mı? dijital mi?)       │
   └───────────────────────┬─────────────────────────────┘
                           ▼
   ┌─────────────────────────────────────────────────────┐
   │  4. NE HIZDA?  Sembol/baud hızını oku                │
   │     (saniyede kaç sembol? darbe ne kadar dar?)       │
   └───────────────────────┬─────────────────────────────┘
                           ▼
   ┌─────────────────────────────────────────────────────┐
   │  5. NE RİTİMDE?  Zaman yapısı / periyot              │
   │     (sürekli mi? patlamalı (burst) mı? periyodik mi?)│
   └───────────────────────┬─────────────────────────────┘
                           ▼
   ┌─────────────────────────────────────────────────────┐
   │  6. EŞLEŞTİR:  sigidwiki + frekans tahsis tablosu    │
   │  7. ARACI DENE:  doğru decoder'ı çalıştır            │
   └─────────────────────────────────────────────────────┘
```

> **Altın kural:** Frekans **tek başına** çoğu zaman cevabı verir. Çünkü spektrum lisanslıdır — her bant belli bir kullanıma ayrılmıştır (uçak, denizcilik, amatör, ISM...). 1090 MHz'de bir sinyal görüyorsan, dünyada o frekans uçak transponderlarına ayrıldığı için %99 ADS-B'dir. **Önce "bu frekans kime ait?" diye sor.**

---

<a id="2"></a>
## 2.  Sinyalin Beş Parmak İzi

Bir sinyali tanımlamak, bir insanı tanımlamak gibidir: boy, kilo, ten rengi yerine **merkez frekans, bant genişliği, modülasyon, sembol hızı, zaman yapısı.** Bu beşlisini doğru ölçersen kimliği neredeyse kesinleşir.

### 2.1 Merkez Frekans (center frequency)
Sinyalin spektrumda oturduğu orta nokta. SDR yazılımında (SDR++, SDRangel, GQRX) imleci sinyalin tam ortasına koy ve oku. **En güçlü ipucu budur** — çünkü frekans tahsisleri yasal olarak sabittir.

### 2.2 Bant Genişliği (bandwidth)
Sinyalin spektrumda kapladığı **genişlik** (kHz/MHz). Genel sezgi:
- **Dar (< 25 kHz):** Az bilgi taşır — telsiz sesi, çağrı cihazı, basit telemetri.
- **Orta (25 kHz – birkaç yüz kHz):** APT görüntü, bazı dijital yayınlar.
- **Geniş (MHz'ler):** Yüksek veri — TV, Wi-Fi, hücresel taşıyıcılar, ADS-B darbeleri (anlık geniş).

Ölçmek için: waterfall'da sinyalin sol ve sağ kenarındaki frekansları oku, farkını al.

### 2.3 Modülasyon (modulation)
Bilgi taşıyıcıya **nasıl** bindiriliyor? Üç büyük aile:

| Tür | Ne değişiyor | Görsel/işitsel ipucu | Örnek |
|---|---|---|---|
| **Genlik (AM/ASK/OOK)** | Sinyalin **gücü** | Waterfall'da parlaklık titrer; OOK'ta nokta-boşluk | Havacılık sesi (AM), garaj kumandası (OOK) |
| **Frekans (FM/FSK/GFSK)** | Sinyalin **frekansı** | İki/çok paralel çizgi arası gidip gelir | FM radyo, POCSAG (FSK), rtl_433 sensörleri |
| **Faz (PSK/QPSK)** | Sinyalin **fazı** | Tek çizgi, kulağa düz "tıslama"; constellation gerekir | PSK31, uydu telemetri, modern dijital |

> **Püf:** **OOK (On-Off Keying)**, en ilkel ve en sık karşılaşacağın dijital moddur — sinyal ya **var** ya **yok** (Mors gibi). Garaj kumandaları, ucuz IoT, eski zil sistemleri OOK kullanır. Waterfall'da kesik kesik yanıp sönen yatay çizgi = büyük ihtimalle OOK.

### 2.4 Sembol / Baud Hızı (symbol rate)
Saniyede kaç **sembol** (bilgi birimi) gönderiliyor? Bu, çözücüye söylemen gereken kritik parametredir (örn. POCSAG 512/1200/2400 bps, AIS 9600 bps). Görsel SDR araçlarında bir "burst"ün en dar darbesini ölçerek tahmin edersin; URH gibi araçlar bunu **otomatik** tespit edebilir.

### 2.5 Zaman Yapısı (timing / burst pattern)
- **Sürekli (continuous):** Hiç durmayan yayın (FM radyo, bazı telemetri).
- **Patlamalı (burst):** Kısa paketler, aralarda sessizlik (ADS-B, AIS, sensörler — periyodik "tık").
- **Periyodik:** Düzenli aralıklarla tekrar (her saniye, her 30 sn). NOAA uydusu geçişi gibi.

> **Hepsini birleştir:** "1090 MHz, çok kısa patlamalar, genlik modülasyonlu, ~1 Mbps" → bu profil **yalnızca ADS-B**'ye uyar. "162 MHz, 9600 baud GMSK, düzenli kısa paketler" → **AIS**. Parmak izleri çakıştığında kimlik kesinleşir.

---

<a id="3"></a>
## 3.  Referans Kütüphaneleri — sigidwiki ve Arkadaşları

Sinyali ölçtün; şimdi **eşleştirme** zamanı. Tek tek ezberlemek imkânsız — referans veritabanları kullan:

- **Signal Identification Wiki (sigidwiki.com):** SIGINT dünyasının "tür rehberi". Her sinyalin **waterfall ekran görüntüsü**, **ses örneği (.wav)**, frekansı, modülasyonu ve çözüm aracını listeler. Gördüğün deseni oradaki görsellerle karşılaştır. *Bilinmeyen sinyal avının 1 numaralı durağı.* (teyit et: site içeriği topluluk katkısıdır, frekanslar bölgesel değişebilir)
- **Frekans tahsis tabloları:** Ülkenin telekom otoritesinin (Türkiye'de BTK) ya da ITU'nun band planı — "bu frekans kime ayrılmış?" sorusunun resmi cevabı.
- **RadioReference.com:** Özellikle bölgesel/yerel frekans veritabanı (ağırlıklı Kuzey Amerika ama metodoloji evrensel).
- **SDR topluluk forumları** (RTL-SDR.com blog, r/RTLSDR): "şu görüntü ne?" diye sorulan binlerce çözülmüş vaka.

> **İş akışı:** Beş parmak izini ölç → sigidwiki'de görsel/ses karşılaştır → aday protokolü bul → o protokolün adanmış çözücüsünü kur ve dene. Eşleşirse decode başlar; eşleşmezse parmak izlerinden birini yanlış ölçtün, geri dön.

---

<a id="4"></a>
## 4.  ADS-B — Uçakları Gökyüzünden Çek

**ADS-B (Automatic Dependent Surveillance–Broadcast)**, SIGINT'e yeni başlayan herkesin ilk "vay be!" anıdır. Uçaklar konum, irtifa, hız ve kimliklerini **şifresiz** olarak yayınlar; sen bir antenle bunu yakalayıp kendi canlı uçak haritanı kurarsın. Tamamen meşru ve müthiş tatmin edici.

| Özellik | Değer |
|---|---|
| **Frekans** | **1090 MHz** (1090ES, ana mod) — ayrıca 978 MHz (UAT, ağırlıklı ABD genel havacılık) |
| **Modülasyon** | PPM (Pulse Position Modulation), genlik temelli |
| **Veri hızı** | 1 Mbps |
| **Araç** | **dump1090** (mutlak klasik), **readsb** (modern fork), tar1090 (web arayüz) |
| **Anten** | 1090 MHz'e ayarlı dikey anten (kolaylıkla DIY); yükseklik menzili artırır |

### Adım adım akış (dump1090)

```bash
# 1. dump1090 ile 1090 MHz'i dinle, web arayüzünü aç
dump1090 --interactive --net

# 2. Tarayıcıda haritayı aç (yerel sunucu)
# http://localhost:8080  → uçaklar harita üzerinde belirir
```

`--interactive` terminalde canlı tablo, `--net` ise web haritası ve diğer araçlara veri akışı (Beast/SBS formatı, port 30003/30005) sağlar.

### Örnek çıktı (terminal, interactive)

```
 ICAO   Flight   Altitude  Speed  Hdg   Lat       Lon        Msgs  Seen
------- -------- --------  -----  ---   -------   --------    ----  ----
4BABCD  THY1924   37000    458    092   41.012    29.045      512   0
3C6DEF  DLH441    11200    312    268   40.998    28.812       88   1
```

- **ICAO:** Uçağın benzersiz 24-bit donanım adresi (hex). Dünyada her uçağa tahsisli.
- **Flight:** Çağrı işareti (THY1924 = Turkish Airlines seferi).
- Konum/irtifa/hız → bunları harita üzerinde gerçek zamanlı izlersin.

> **Püf:** Topladığın veriyi **FlightAware / FlightRadar24 / OpenSky Network** gibi platformlara besleyebilirsin (feeder olursun) → karşılığında premium erişim alırsın. Kendi anteninle küresel bir izleme ağına katkı = en sevilen yasal SDR projelerinden. **Menzil sırrı:** Antenini ne kadar yükseğe (çatı, pencere) ve engelsiz koyarsan, ufuk çizgisi o kadar uzar — düz arazide 250+ deniz mili görmek mümkün.

---

<a id="5"></a>
## 5.  AIS — Gemileri Denizden Çek

ADS-B'nin denizdeki kuzeni. **AIS (Automatic Identification System)**, gemilerin konum, hız, rota, isim ve MMSI kimliğini **şifresiz** yaydığı sistemdir. Sahil/liman yakınındaysan onlarca gemiyi haritalarsın.

| Özellik | Değer |
|---|---|
| **Frekans** | **161.975 MHz** (AIS 1 / kanal 87B) ve **162.025 MHz** (AIS 2 / kanal 88B) |
| **Modülasyon** | GMSK, 9600 bps |
| **Araç** | **AISdecoder**, **gnuais / gr-ais**, **rtl-ais** (RTL-SDR + decoder tek pakette) |
| **Anten** | Denizcilik VHF bandına (~162 MHz) ayarlı dikey anten |

### Adım adım akış (rtl-ais — en pratik)

```bash
# rtl-ais: iki AIS kanalını aynı anda dinler, NMEA cümleleri üretir,
# UDP ile haritalama yazılımına (OpenCPN vb.) yollar
rtl_ais -n -p 0
```

`-n` ile NMEA çıktısı UDP üzerinden yayınlanır; bunu **OpenCPN** gibi bir deniz haritası programına besleyip gemileri canlı görürsün.

### Örnek çıktı (ham NMEA AIS cümlesi)

```
!AIVDM,1,1,,A,15M: Dt0001G?tO`K>RA1wUbN0TKH,0*5C
        │                                    └ sağlama (checksum)
        └ payload (6-bit ASCII kodlu konum/MMSI/hız — decoder çözer)
```

Decoder bunu insan-okur alanlara dönüştürür: **MMSI** (gemi kimliği), enlem/boylam, SOG (hız), COG (rota), gemi adı/tipi.

> **Not:** Tıpkı ADS-B gibi, AIS verisini **MarineTraffic / AISHub** gibi platformlara besleyebilirsin. Sahile uzaksan gemi göremezsin — AIS VHF'tir, menzili ufukla sınırlıdır (~40-50 deniz mili).

---

<a id="6"></a>
## 6.  ACARS — Uçakların Metin Mesajları

**ACARS (Aircraft Communications Addressing and Reporting System)**, uçaklarla yer istasyonları arasındaki kısa **metin** mesajlaşma sistemidir (kalkış/varış saatleri, motor verisi, hava durumu istekleri, bakım mesajları). VHF üzerinden açık yayınlanan operasyonel mesajları çözebilirsin.

| Özellik | Değer |
|---|---|
| **Frekans** | **131.x MHz** bandı (bölgeye göre 131.525, 131.725, 131.825 MHz yaygın) |
| **Modülasyon** | MSK / AM-tabanlı, 2400 bps |
| **Araç** | **acarsdec** (çoklu kanal), **acarsdeco2** |

### Adım adım akış (acarsdec)

```bash
# Birden çok ACARS frekansını aynı anda dinle (RTL-SDR cihaz 0)
acarsdec -r 0 131.525 131.725 131.825
```

### Örnek çıktı

```
[#1 (F:131.725 L:-22) ----------------------------------------
Mode: 2  Label: H1  Id: 7
Aircraft reg: TC-JJK  Flight: TK0079
Message:
  POS N41.2 E028.8 FL370 ETA 1542 ...
]
```

> **Bağlam:** ACARS mesajları operasyoneldir, kişisel iletişim değildir — açık yayınlandıkları için dinlenmesi yaygın bir hobidir. Yine de **mesaj içeriğini kötüye kullanma** (gizliliğe saygı). Bir sonraki nesil sistem ise daha geniş bantlı VDL Mode 2'dir.

---

<a id="7"></a>
## 7.  POCSAG / FLEX — Çağrı Cihazları (Pager)

Çağrı cihazları (pager) hâlâ hastane, acil servis ve bazı kurumlarda kullanılır. **POCSAG** ve **FLEX**, pager mesajlarını taşıyan protokollerdir ve çoğunlukla **şifresizdir**.

| Özellik | Değer |
|---|---|
| **Frekans** | Bölgeye göre değişir — yaygın bantlar 138–175 MHz ve ~929/931 MHz (teyit et: yerel tahsis BTK band planına bağlı) |
| **Modülasyon** | FSK (2-FSK), POCSAG 512/1200/2400 bps; FLEX daha yüksek hızlı |
| **Araç** | **multimon-ng** (POCSAG/FLEX çözer), **PDW** (Windows) |

### Tipik zincir (rtl_fm → multimon-ng)

```bash
# RTL-SDR'dan ham FM akışını multimon-ng'ye boru ile ver
rtl_fm -f 153.350M -s 22050 - | multimon-ng -t raw -a POCSAG1200 -f alpha /dev/stdin
```

> **GİZLİLİK UYARISI — DİKKAT:** Pager içeriği teknik olarak şifresiz yayınlansa da, **kişisel/tıbbi/özel bilgi** taşıyabilir (hasta isimleri, acil çağrılar). Birçok ülkede bu içeriği **kaydetmek, paylaşmak, kullanmak yasadışıdır** — sinyalin "açık" olması içeriğini kullanmanı meşru kılmaz. **Önerimiz:** Bu protokolü **yalnızca modülasyon/çözümleme mekaniğini öğrenmek** için, içeriğe odaklanmadan, kaydetmeden incele. Şüphedeysen hiç dokunma; bunun yerine ADS-B/AIS/NOAA gibi içeriği gerçekten kamusal olan protokollerle pratik yap.

---

<a id="8"></a>
## 8.  APRS — Amatör Paket Radyo

**APRS (Automatic Packet Reporting System)**, amatör radyocuların konum, telemetri, kısa mesaj ve hava durumu paylaştığı bir paket sistemidir. Tamamen açık ve topluluk odaklıdır — amatör radyo dünyasının "canlı haritası".

| Özellik | Değer |
|---|---|
| **Frekans** | **144.800 MHz** (Avrupa/Türkiye) — Kuzey Amerika'da 144.390 MHz |
| **Modülasyon** | AFSK 1200 bps (Bell 202 tonları) FM taşıyıcı üzerinde, AX.25 çerçeve |
| **Araç** | **Direwolf** (yazılım TNC/modem — referans araç), multimon-ng (basit) |

### Adım adım akış (Direwolf)

```bash
# Direwolf'u RTL-SDR akışıyla besle (rtl_fm üzerinden) ve paketleri çöz
rtl_fm -f 144.800M -s 24000 - | direwolf -c sdr.conf -r 24000 -D 1 -
```

### Örnek çıktı (çözülmüş APRS paketi)

```
TA1ABC-9>APRS,WIDE1-1:!4101.23N/02858.45E>Mobil istasyon, Istanbul
       │                  │         │       └ durum/yorum metni
       │                  └ enlem    └ boylam
       └ gönderen çağrı işareti (amatör lisans)
```

> **Not:** APRS'i çözmek (RX) yasaldır. **Yayın yapmak (TX)** için **amatör telsiz lisansı** gerekir (Türkiye'de bu sınava girip lisans alabilirsin — SIGINT öğrenenler için harika bir resmi adım). Çözdüğün paketleri **aprs.fi** gibi haritalarda da görebilirsin.

---

<a id="9"></a>
## 9.  NOAA APT & Meteor-M LRPT — Uydudan Hava Durumu Görüntüsü

**SDR'ın taç mücevheri.** Başının üzerinden geçen bir **hava durumu uydusundan**, kendi antenınle, gerçek **görüntü** indirirsin. İlk başarılı APT görüntüsünü aldığında SDR'a tutkun olursun. %100 meşru, %100 büyüleyici.

### 9.1 NOAA APT (analog)

| Özellik | Değer |
|---|---|
| **Uydular** | NOAA-15 (137.620 MHz), NOAA-18 (137.9125 MHz), NOAA-19 (137.100 MHz) *(teyit et: aktiflik zamanla değişebilir)* |
| **Frekans** | **137 MHz** bandı |
| **Modülasyon** | APT — FM taşıyıcı üzerinde AM görüntü hattı (analog) |
| **Araç** | **WXtoImg** (klasik), **noaa-apt** (modern, açık kaynak) |
| **Anten** | **QFH** (Quadrifilar Helix) veya çapraz dipol — dairesel polarizasyon ideal |

#### Adım adım akış

```
1. GEÇİŞ ZAMANI:  Bir uydu takip uygulaması (Gpredict, "Look4Sat" mobil)
                  ile uydunun ne zaman tependen geçeceğini öğren.
                  (Uydu ufkun üzerindeyken ~10-15 dk pencere)
        │
        ▼
2. KAYIT:         Geçiş sırasında SDR ile 137.xxx MHz'i WFM modda
                  kaydet (.wav, 11025 veya 48000 Hz). Doppler kayması
                  olur — bazı araçlar otomatik düzeltir.
        │
        ▼
3. ÇÖZ:           noaa-apt ile .wav'ı görüntüye dönüştür:
                  noaa-apt kayit.wav -o sonuc.png
        │
        ▼
4. SONUÇ:         Dünyanın o anki uydu görüntüsü — bulutlar,
                  kıyı şeritleri, kızılötesi/görünür kanallar.
```

```bash
# Kaydettiğin WAV'ı görüntüye çevir (noaa-apt)
noaa-apt gecis_noaa19.wav -o izmir_gokyuzu.png
```

### 9.2 Meteor-M LRPT (dijital, daha keskin)

Rus **Meteor-M N2** serisi, NOAA'nın **dijital** ve **renkli/daha yüksek çözünürlüklü** muadilidir (LRPT — Low Rate Picture Transmission).

| Özellik | Değer |
|---|---|
| **Frekans** | ~**137.1 / 137.9 MHz** bandı *(uydu ve döneme göre değişir — teyit et)* |
| **Modülasyon** | QPSK, dijital LRPT |
| **Araç** | **SatDump** (modern, çok-uydu, önerilen), eski zincir: meteor_demod + LRPToffline |
| **Anten** | NOAA ile aynı (QFH / çapraz dipol, 137 MHz) |

> **Püf:** **SatDump** bugün hem NOAA APT hem Meteor LRPT hem de daha fazlasını tek çatı altında çözen modern araçtır — yeni başlıyorsan doğrudan SatDump öğren. **QFH anten** bu iş için altın standarttır çünkü uydular dairesel polarize yayar ve gökyüzünde hareket eder; yönlü anten gerektirmez. İlk denemende sabırlı ol: temiz görüntü için iyi anten + açık ufuk + doğru geçiş zamanı şart.

---

<a id="10"></a>
## 10.  rtl_433 — Evdeki Her Şeyin Telsizi

Etrafındaki onlarca cihaz **433.92 MHz** ve **868 MHz** (ISM bantları) üzerinden sürekli telsizle konuşuyor: meteoroloji istasyonları, lastik basıncı sensörleri (TPMS), akıllı prizler, kapı/pencere sensörleri, su/gaz sayaçları. **rtl_433**, bu yüzlerce protokolü **otomatik tanıyıp** çözen sihirli araçtır.

| Özellik | Değer |
|---|---|
| **Frekans** | **433.92 MHz** (Avrupa ISM ana), **868 MHz** (Avrupa), 315/915 MHz (Amerika) |
| **Modülasyon** | Çoğunlukla OOK/ASK ve FSK (cihaza göre) |
| **Araç** | **rtl_433** (200+ cihaz protokolünü hazır tanır) |

### Adım adım akış

```bash
# 433.92 MHz'i dinle, tanınan tüm cihazları JSON olarak dök
rtl_433 -f 433.92M -F json
```

`rtl_433` arka planda yüzlerce protokol şablonunu dener; bir eşleşme bulunca cihazı **otomatik** kimliklendirir.

### Örnek çıktı (JSON)

```json
{ "time":"2026-05-31 14:22:10", "model":"Nexus-TH", "id":42,
  "channel":1, "temperature_C":23.4, "humidity":56, "battery_ok":1 }
{ "time":"2026-05-31 14:22:18", "model":"Toyota-TPMS", "id":"a1b2c3d4",
  "pressure_kPa":230.0, "temperature_C":28 }
```

> **Püf — KENDİ sensörlerin:** En tatmin edici alıştırma, **kendi evindeki** kablosuz termometreni, hava istasyonunu ya da **kendi arabanın** lastik basıncı sensörünü (yanına park edip) rtl_433 ile görmektir. Burada **kendi cihazını** gözlemliyorsun = tamamen meşru. Sonuçları MQTT'ye basıp ev otomasyonuna (Home Assistant) bağlayabilirsin.  Komşunun sensörlerini de görebilirsin — bunlar açık yayınlar olsa da, kişisel verisini (örn. evde var/yok çıkarımı) kötüye kullanmak etik/yasal sınırı aşar; gözlemle, kullanma.

---

<a id="11"></a>
## 11.  LoRa — Uzun Menzilli IoT

**LoRa (Long Range)**, düşük güçle kilometrelerce menzil sağlayan bir IoT modülasyonudur (chirp spread spectrum — CSS). Akıllı şehir sensörleri, tarım telemetri, varlık takibi kullanır.

| Özellik | Değer |
|---|---|
| **Frekans** | **868 MHz** (Avrupa), **915 MHz** (Amerika), 433 MHz (bazı bölgeler) |
| **Modülasyon** | CSS (Chirp Spread Spectrum) — waterfall'da **karakteristik eğik "kayan" çizgiler** (chirp) |
| **Araç** | **gr-lora / gr-lora_sdr** (GNU Radio), bazı SDRangel eklentileri |

> **Tanıma ipucu:** LoRa, waterfall'da **yukarı/aşağı süpüren eğik rampa** desenleriyle hemen belli olur — başka hiçbir yaygın sinyal böyle "kayan" görünmez. **LoRaWAN** (üst katman ağ protokolü) genelde **şifrelidir**; yani fiziksel chirp'i görebilir/demodüle edebilirsin ama uygulama verisi anahtarsız okunamaz. Kendi LoRa geliştirme kartınla (Heltec/TTGO) **kendi** paketlerini gönderip çözmek mükemmel bir öğrenme projesidir.

---

<a id="12"></a>
## 12.  Amatör Dijital Modlar — FT8 / PSK31 / RTTY

Kısa dalga (HF) bantlarında amatör radyocular, dünyanın öbür ucuyla **dijital modlarla** haberleşir. Bunları **dinlemek (RX)** tamamen yasaldır ve "zayıf sinyalle dünya çapı iletişim" mucizesini gözlemlemenin en iyi yoludur. (HF için **upconverter'lı RTL-SDR** ya da doğrudan HF SDR — örn. ham bant kapsayan cihaz — gerekir.)

| Mod | Ne işe yarar | Araç | Tipik bantlar |
|---|---|---|---|
| **FT8** | Aşırı zayıf sinyalle (gürültü altında!) kısa "el sıkışma" mesajları — modern amatör fenomeni | **WSJT-X** | 14.074 MHz (20m) ve diğer ham bantlar |
| **PSK31** | Klavyeden-klavyeye sohbet, çok dar bant, verimli | **fldigi** | 14.070 MHz civarı |
| **RTTY** | Klasik radyo teleprinter (Baudot) | **fldigi** | Çeşitli HF ham bantları |

### FT8 akışı (WSJT-X)

```
1. SDR yazılımının sesini (USB modda, örn. 14.074 MHz) WSJT-X'e ver
   (sanal ses kablosu / VB-Cable / PulseAudio loopback).
2. WSJT-X bandı 15 saniyelik pencerelerde dinler ve gürültünün
   ALTINDAKİ sinyalleri bile çözer.
3. Ekranda dünyanın dört bir yanından çağrı işaretleri akar:
   "CQ JA1XYZ PM95"  → Japonya'dan bir istasyon çağrı yapıyor.
```

> **Püf:** **FT8'in büyüsü**, -20 dB SNR'da (yani sinyal gürültüden zayıfken!) bile çözebilmesidir — bu yüzden minik antenle, düşük güçle kıtalararası "görülmek" mümkündür. **Sadece RX yaparak** (lisans gerekmez) ekranında anlık olarak hangi ülkelerin "duyulduğunu" izlemek, propagasyonu (iyonosferik yayılım) öğrenmenin en somut yoludur. Aldığın spotları **PSK Reporter** haritasında dünya ölçeğinde görebilirsin.

---

<a id="13"></a>
## 13.  Sinyal Tersine Mühendislik — URH ile KENDİ Kumandanı Çöz

Şimdiye kadar **bilinen** protokolleri hazır araçlarla çözdük. Peki ya **bilinmeyen** bir sinyal — örneğin **kendi garaj kapı kumandan**, kendi kablosuz kapı zilin, kendi oyuncağının uzaktan kumandası? İşte burada **tersine mühendislik** devreye girer ve en güçlü araç **Universal Radio Hacker (URH)**'tır.

> **SINIR:** Bu teknik **yalnızca KENDİ cihazların** içindir. Kendi malın olan bir kumandanın sinyalini analiz etmek bir öğrenme egzersizidir. **Başkasının** garajını/aracını/kapısını hedef almak suçtur ve bu rehberin kapsamı dışındadır.

### URH ile bilinmeyen OOK/FSK sinyalini çözme (kavram + akış)

```
┌────────────────────────────────────────────────────────────┐
│ 1. KAYDET (Record):                                        │
│    URH ile, kumandaya basarken sinyali kaydet (örn.        │
│    433.92 MHz). Her tuş basışında bir "burst" görürsün.    │
├────────────────────────────────────────────────────────────┤
│ 2. ANALİZ (Analysis):                                      │
│    URH otomatik olarak modülasyonu (çoğu kumandada OOK),   │
│    sembol hızını ve sinyal/boşluk eşiğini tahmin eder.     │
├────────────────────────────────────────────────────────────┤
│ 3. DEMODÜLE → BİT DİZİSİ:                                  │
│    Burst'ler 1/0 dizisine dönüşür. Örn. her tuş için:      │
│    Tuş A:  1010 1100 1010 0101 ...                         │
│    Tuş B:  1010 1100 1010 1010 ...                         │
│    (Sabit bir "preamble" + değişen "komut" deseni görürsün)│
├────────────────────────────────────────────────────────────┤
│ 4. DESENİ ÇÖZ (Protocol):                                 │
│    Farklı tuşları kıyaslayarak hangi bitlerin "adres",     │
│    hangilerinin "komut" olduğunu çıkarırsın. URH bunları   │
│    etiketlemene ve alanları isimlendirmene izin verir.     │
└────────────────────────────────────────────────────────────┘
```

### Somut örnek (kavramsal)

Eski, basit bir garaj kumandası genelde **sabit kod (fixed code)** kullanır: her basışta **aynı** bit dizisini gönderir. URH'ta iki kez kaydedip karşılaştırırsan dizilerin birebir aynı olduğunu görürsün → bu, neden eski kumandaların güvensiz olduğunu öğretir (Bölüm 6'da "replay" zafiyeti). Modern kumandalar ise **rolling code** (her basışta değişen, senkron sayaç) kullanır → kaydedilen dizi tekrar işe yaramaz. **İşte tersine mühendilik, sana bir cihazın güvenli mi güvensiz mi tasarlandığını somut olarak gösterir** — bu da savunma mühendisliğinin özüdür.

> **Öğrenme değeri:** URH ile kendi cihazlarını çözmek, "modülasyon → sembol → bit → protokol alanı" zincirini **elinle** kurmanı sağlar. Hazır decoder'lar bu zinciri senin yerine yapar; URH ile bir kez kendin yaparsan **tüm dijital protokolleri** çok daha derin anlarsın.

---

<a id="14"></a>
## 14.  HÜCRESEL — 2G/3G/4G/5G Mimarisi (Prensip + Yasal Sınır)

> **EN ÖNEMLİ UYARI — TEKRAR:** Bu bölüm **mimariyi ve neden-savunmasız/savunmalı** sorusunu anlatır. **Hiçbir adım-adım yetkisiz dinleme reçetesi içermez ve içermeyecektir.** Başkasının hücresel haberleşmesini (görüşme/SMS/veri) izinsiz dinlemek/çözmek **ağır cezalı suçtur** — TCK ve uluslararası mevzuat. Burada amaç: bir savunmacının, saldırının nasıl mümkün olduğunu **anlayarak** kendini ve kurumunu koruyabilmesidir. Yasal pratik yolu Bölüm 15'te.

Hücresel ağlar, açık yayınların (ADS-B/AIS) tam tersidir: **kullanıcı trafiği şifrelidir** ve dinlenmesi hem teknik hem yasal olarak engellenmiştir. Ama her nesil aynı derecede sağlam değildir — işte özü.

### 14.1 GSM / 2G — Tarihsel Zafiyetin Dersi

GSM, 1990'larda tasarlandı; bugünün tehdit modeline göre **zayıf**tır ve bu *eğitici* bir örnektir.

**Mimari (temel kavramlar):**
- **MS (Mobile Station):** Senin telefonun. **BTS (Base Transceiver Station):** Baz istasyonu (hücre kulesi).
- **ARFCN (Absolute Radio Frequency Channel Number):** Hangi frekans kanalında çalışıldığını belirten numara (900/1800 MHz GSM bantları).
- **Erişim yöntemi:** TDMA + FDMA — yani hem **frekans** bölmeli hem **zaman** bölmeli. Bir taşıyıcı 8 zaman dilimine (timeslot) bölünür; her kullanıcı kendi diliminde konuşur.

```
GSM TDMA ÇERÇEVE YAPISI (kavramsal)
═══════════════════════════════════════════════════════════
Bir taşıyıcı (200 kHz), zamanda 8 dilime (TS0..TS7) bölünür:

│ TS0 │ TS1 │ TS2 │ TS3 │ TS4 │ TS5 │ TS6 │ TS7 │  → 1 çerçeve
  ▲                                                  (~4.6 ms)
  └─ Genelde KONTROL kanalları (BCCH, CCCH) burada

KONTROL KANALLARI (yayın/sinyalizasyon — açık niteliği farklı):
  • BCCH (Broadcast Control Channel):  hücre kimliği, komşu hücreler,
        ağ parametreleri → telefonların ağı bulması için YAYINLANIR
  • CCCH / PCH / RACH:  çağrı/erişim sinyalizasyonu
TRAFİK KANALLARI (TCH):  asıl ses/veri → ŞİFRELİ (A5 ailesi)
═══════════════════════════════════════════════════════════
```

**Neden 2G güvensiz sayılır (kavramsal):**
- **A5/1 şifrelemesi tarihsel olarak zayıftır:** Akademik çalışmalar A5/1'in kriptografik olarak kırılabildiğini gösterdi; **A5/2 daha da zayıf**, **A5/0 ise hiç şifreleme yok** demektir. Bu, GSM'in *tasarım çağının* sınırlarını yansıtır.
- **Tek yönlü kimlik doğrulama:** GSM'de **ağ telefonu doğrular ama telefon ağı doğrulamaz.** Bu asimetri, sahte baz istasyonu (**IMSI catcher / "stingray"**) kavramının temelidir: telefon, kendini en güçlü hücre sanan sahte bir BTS'e bağlanabilir. *(Bu yapının nasıl çalıştığını **bilmek** seni korur — örn. neden 2G'yi kapatmak savunma sağlar.)*
- **gr-gsm gibi araçlar** akademik olarak **downlink kontrol/yayın kanallarını** (BCCH gibi) gözlemlemek için kullanılır — yani *zaten yayınlanan* sinyalizasyonu incelemek. **Kullanıcı trafiğini çözmek bambaşka, yasadışı bir eylemdir** ve bu rehberde tarif edilmez.

> **Savunma dersi:** Telefonunu/kurumunu korumak için **2G'yi devre dışı bırak** (modern telefonlarda "yalnızca 4G/5G" zorlanabilir), **IMSI catcher tespiti** olan güvenlik uygulamalarını değerlendir, kritik iletişimde **uçtan uca şifreli** uygulamalar (Signal vb.) kullan — çünkü E2E şifreleme, taşıyıcı katmanı ele geçirilse bile içeriği korur.

### 14.2 3G / UMTS — Neden Dinlemesi Zorlaştı

3G (UMTS, W-CDMA), GSM'in zafiyetlerine doğrudan yanıttır:

- **W-CDMA (Wideband Code Division Multiple Access):** Kullanıcılar **frekans veya zamanla değil, kodla** ayrılır. Herkes aynı geniş bandı (5 MHz) paylaşır; her birinin sinyali kendine özgü bir **yayılım kodu** ile çarpılır. Doğru kodu bilmeyen için sinyal **gürültüye gömülü** görünür → **yayılı spektrum (spread spectrum)** doğal bir gizlilik katmanı ekler.
- **Karşılıklı (mutual) kimlik doğrulama:** GSM'in aksine, 3G'de **telefon de ağı doğrular** (USIM + AKA protokolü). Bu, sahte baz istasyonu saldırılarını GSM'e göre çok zorlaştırır.
- **Daha güçlü şifreleme** (KASUMI/SNOW tabanlı, A5/1'den ileri).

**Sonuç:** 3G'nin yayılı spektrumu + karşılıklı doğrulaması, pasif dinlemeyi pratikte çok zor kılar.

### 14.3 4G / LTE — Açık Yayın vs. Şifreli Trafik

LTE, "neyin açık, neyin kapalı" ayrımını net biçimde örnekler — SIGINT etik çizgisini öğrenmek için ideal.

- **OFDMA (Orthogonal Frequency-Division Multiple Access):** Bant, binlerce dar **alt-taşıyıcıya** bölünür; kullanıcılara zaman-frekans "kaynak blokları" (resource block) tahsis edilir. Çok verimli ama yapısı karmaşık.
- **EARFCN:** LTE'nin kanal/frekans numaralandırması (GSM'deki ARFCN'in LTE karşılığı).
- **Hücre arama ve YAYIN bilgileri (açık niteliği):** Telefon ağa bağlanmak için önce **senkronizasyon sinyallerini** (PSS/SSS) bulur, sonra **MIB (Master Information Block)** ve **SIB (System Information Block)** denen **yayın** mesajlarını okur. Bunlar **hücre kimliği, bant genişliği, komşu hücre listesi** gibi **herkese açık** ağ parametreleridir — tıpkı GSM'in BCCH'i gibi, telefonların ağı tanıması için **kasıtlı yayınlanır.**
- **AKADEMİK araçlar:** **srsRAN** (açık kaynak LTE/5G yığını, *araştırma ve kendi test ağını kurmak* için) ve **FALCON** gibi araçlar, LTE **downlink kontrol kanalı (PDCCH) analizi** ve hücre ölçümü için akademik olarak kullanılır — yani **yayın/sinyalizasyon katmanını** incelemek için.
- **Kullanıcı trafiği şifrelidir:** Asıl veri (PDSCH üzerindeki kullanıcı düzlemi), NAS/AS güvenlik bağlamıyla şifrelenir. **Bu yüzden LTE kullanıcı haberleşmesini çözmek hem teknik olarak engellenmiş hem yasadışıdır.**

> **Kritik ayrım (aklında kazınsın):** Bir LTE hücresinin **var olduğunu, hangi frekansta yayın yaptığını, MIB/SIB parametrelerini** ölçmek (akademik/araştırma, açık yayın) ile o hücredeki **bir kullanıcının verisini çözmek** (şifreli, özel, **SUÇ**) **tamamen farklı** iki şeydir. Birincisi radyo mühendisliği; ikincisi ağır suç.

### 14.4 5G NR — Kısaca

5G NR (New Radio), LTE'nin OFDMA temelini sürdürür ama **çok daha geniş bant**, **mmWave** (yüksek GHz) seçenekleri, **kapsamlı ışın yönlendirme (beamforming)** ve **daha güçlü gizlilik özellikleri** ekler:
- **SUCI (Subscription Concealed Identifier):** 5G, kalıcı abone kimliğini (SUPI/IMSI) açıktan göndermek yerine **şifreleyerek** (SUCI) iletir → IMSI catcher tarzı kimlik yakalama saldırılarını tasarımdan zorlaştırır.
- Beamforming ve geniş bant, pasif gözlemi fiziksel olarak da daha zor kılar.
- Yine: **yayın/senkronizasyon bilgisi** (SSB vb.) açıktır; **kullanıcı trafiği şifrelidir.**

> **Genel savunma çıkarımı:** Nesil ilerledikçe (2G→5G) hem **şifreleme** hem **kimlik gizliliği** hem **karşılıklı doğrulama** güçlenir. En zayıf halka neredeyse her zaman **en eski nesle (2G) düşürme (downgrade) saldırısıdır** — bu yüzden cihazını mümkünse 2G'den arındırmak somut bir koruma sağlar.

---

<a id="15"></a>
## 15.  Net Çerçeve — Gözlem mi, Dinleme mi?

Tüm bu bölümün ahlaki/hukuki pusulası tek bir ayrımda toplanır:

```
┌──────────────────────────────┬──────────────────────────────┐
│  GÖZLEM (çoğunlukla meşru)    │  DİNLEME (yetkisizse SUÇ)     │
├──────────────────────────────┼──────────────────────────────┤
│ • ADS-B / AIS / NOAA / APRS   │ • Başkasının telefon görüşmesi│
│   → herkese açık yayınlar     │ • Başkasının SMS / mobil verisi│
│ • Hücresel YAYIN/kontrol      │ • Şifreli özel telsiz içeriği │
│   kanalı varlığı/parametresi  │ • Pager içeriğinin kaydı/kullanımı│
│   (akademik/araştırma)        │ • Herhangi bir ŞİFRELİ özel   │
│ • KENDİ cihazlarının sinyali  │   trafiği çözme girişimi      │
│   (kendi kumandan, sensörün)  │ • Sahte baz istasyonu kurmak  │
└──────────────────────────────┴──────────────────────────────┘
```

### Hücresel/özel haberleşmeyi öğrenmenin YASAL yolu
Eğer bu teknolojiyi gerçekten ellerinle öğrenmek istiyorsan, **doğru ve yasal kapı şudur:**
1. **Kendi test ağın:** **srsRAN / Open5GS** gibi açık kaynak yığınlarla, **kendi SIM kartların ve kendi cihazlarınla** bir özel LTE/5G ağı kur. Kendi ağında, kendi trafiğini analiz etmek meşrudur.
2. **Faraday kafesi / RF-shielded oda:** Yayının dışarı sızmaması ve **başkalarının spektrumuna karışmaman** için izole ortam. Çoğu ülkede vericiyi (BTS tarafı) çalıştırmak **lisans** gerektirir; izolasyon bunu yönetmenin yoludur.
3. **Operatör / üniversite / kurum izni:** Akademik araştırma laboratuvarları, operatör test sahaları — yazılı izinle, kontrollü ortamda.
4. **Amatör radyo lisansı:** TX'e geçmenin (kendi sinyalini yayınlamanın) genel yasal yolu; sınava gir, çağrı işaretini al.

> **Özet ilke:** *Açığı dinle, özeli dinleme.* *Kendi cihazını çöz, başkasınınkini çözme.* *Yayın/kontrol katmanını gözlemlemek (araştırma) ile kullanıcı içeriğini ele geçirmek (suç) arasındaki çizgiyi asla bulanıklaştırma.* Şüphedeysen: RX-only + açık yayın = güvenli liman.

---

<a id="16"></a>
## 16.  ALIŞTIRMALAR — Hepsi Evde, Hepsi Yasal

Bu bölümü "okuyup geçme" — **yap.** Aşağıdakilerin hepsi yasaldır, ucuz bir RTL-SDR ile mümkündür ve seni acemiden ustaya taşır:

1. ** Uydudan görüntü indir (taç görev):** Bir NOAA (137 MHz) ya da Meteor-M geçişini Gpredict/Look4Sat ile yakala, **SatDump** veya **noaa-apt** ile görüntüye çevir. İlk bulut/kıyı görüntün — SDR'ın "vay be" anı.
2. ** Canlı uçak haritası:** **dump1090**'ı 1090 MHz'de çalıştır, tarayıcıda tependeki uçakların haritasını izle. Antenini yükselt, menzilin nasıl arttığını gör. İstersen OpenSky/FlightAware'e feeder ol.
3. ** Gemi haritası (sahildeysen):** **rtl-ais** + OpenCPN ile 162 MHz'den gemileri haritala.
4. ** Kendi sensörlerin:** **rtl_433 -f 433.92M -F json** ile evindeki kablosuz termometreyi/hava istasyonunu, kendi arabanın TPMS lastik sensörünü gör. JSON'u Home Assistant'a bağla.
5. ** Kendi kumandanı çöz:** **URH** ile kendi garaj/zil/oyuncak kumandanın OOK sinyalini kaydet, bit dizisini çıkar, "sabit kod mu, rolling code mu?" sorusunu yanıtla. (Yalnızca **kendi** malın.)
6. ** Dünyayı dinle (sadece RX):** HF-yetkin SDR + **WSJT-X** ile **FT8** (14.074 MHz) penceresini aç; gürültünün altından dünyanın dört bucağından gelen çağrı işaretlerini izle. **PSK Reporter**'da hangi ülkelerin "duyulduğunu" haritada gör. Lisans gerekmez (yalnızca dinliyorsun).
7. ** APRS gözlemi:** 144.800 MHz'i **Direwolf** ile çöz, yerel amatör istasyonların konum/telemetri paketlerini gör; **aprs.fi**'de karşılaştır.

> **Ustalık rozeti:** Bu yedi alıştırmayı bitirdiğinde, "sinyal yakalama → tanımlama → çözme" zincirini hem **açık protokollerde** hem **kendi cihazlarında** uçtan uca yapmış olursun. Hücresel kısmı ise **anlamış** ama yasal sınırın hangi tarafında durduğunu **bilen** biri olarak tamamlarsın — gerçek bir SIGINT savunmacısının duruşu budur.

---

<a id="17"></a>
## 17.  Hızlı Referans Tablosu (Protokol × Frekans × Araç)

| Protokol | Frekans (tipik) | Modülasyon | Ana Araç | Yasal durum |
|---|---|---|---|---|
| **ADS-B** | 1090 MHz (978 UAT) | PPM | dump1090 / readsb |  Açık yayın |
| **AIS** | 161.975 / 162.025 MHz | GMSK 9600 | rtl-ais / gnuais |  Açık yayın |
| **ACARS** | ~131 MHz bandı | MSK 2400 | acarsdec |  Açık (içeriğe saygı) |
| **POCSAG/FLEX** | 138–175 / ~929 MHz* | FSK | multimon-ng / PDW |  İçerik kişisel olabilir — kullanma |
| **APRS** | 144.800 MHz (EU/TR) | AFSK 1200 / AX.25 | Direwolf |  RX yasal (TX=lisans) |
| **NOAA APT** | 137 MHz | APT (FM+AM) | SatDump / noaa-apt |  Açık yayın |
| **Meteor-M LRPT** | ~137 MHz* | QPSK (LRPT) | SatDump |  Açık yayın |
| **rtl_433 (IoT)** | 433.92 / 868 MHz | OOK/ASK/FSK | rtl_433 |  Kendi cihazların |
| **LoRa** | 868 / 915 MHz | CSS (chirp) | gr-lora_sdr |  Kendi cihazların (LoRaWAN şifreli) |
| **FT8 / PSK31 / RTTY** | HF ham bantları (14.074 vb.) | FSK/PSK | WSJT-X / fldigi |  RX yasal |
| **GSM / 2G** | 900 / 1800 MHz (ARFCN) | GMSK | (yalnızca kavram) |  Kullanıcı trafiği = SUÇ |
| **3G / UMTS** | UMTS bantları | W-CDMA | (yalnızca kavram) |  Kullanıcı trafiği = SUÇ |
| **4G / LTE** | LTE bantları (EARFCN) | OFDMA | srsRAN (test ağı) |  Kullanıcı trafiği = SUÇ |
| **5G NR** | Sub-6 / mmWave | OFDMA + beamforming | Open5GS (test ağı) |  Kullanıcı trafiği = SUÇ |

> \* **Teyit et:** Yıldızlı frekanslar bölgesel tahsise (Türkiye'de BTK band planı) ve uydu/dönem değişimine bağlıdır; uygulamadan önce güncel band planını ve uydu durumunu kontrol et.

---

<a id="18"></a>
## 18.  Çapraz Referans & Sonraki Bölüm

Bu bölümde bir sinyali **tanıma** ve **çözme** sanatını öğrendin: beş parmak izi (frekans/bant/modülasyon/sembol/zaman), açık protokollerin uçtan uca akışı, kendi cihazlarını tersine mühendislik, ve hücresel mimarinin neden bazı nesillerde savunmasız bazılarında sağlam olduğu.

** Bir sonraki adım — Bölüm 6: Güvenlik, Açıklar ve Savunma.** Burada öğrendiğin "nasıl çalışır" bilgisini, "nasıl saldırıya uğrar ve nasıl savunulur" perspektifine taşıyacağız: replay/rolling-code zafiyetleri (bu bölümde URH ile gördüğün sabit-kod kumandasının neden tehlikeli olduğu), GPS spoofing/jamming farkındalığı, IMSI-catcher tespiti ve karşı önlemler, spektrum izleme ile anomali tespiti, ve Kanije Kalesi'nin fiziksel/RF tehdit duruşuyla bağlantısı.

> **Kapanış:** Bir sinyali yakalamak teknikse, onu **anlamak** sanattır; ama onu **hangi sınıra kadar** çözebileceğini bilmek **bilgeliktir.** Açık gökyüzü senindir — uçakları, gemileri, uyduları, kendi cihazlarını dinle ve öğren. Özel/şifreli haberleşmenin kapısında ise dur: orası mühendisliğin değil, hukukun ve etiğin alanıdır. Gerçek bir uzman, **yapabildiği** ile **yapması gerekeni** ayırt edendir.
>
> *Bu doküman Kanije Kalesi SIGINT El Kitabı serisinin **5. Bölümü**dür. İlgili: Bölüm 4 (SDR kurulumu/yakalama), **Bölüm 6 (Güvenlik, açıklar, savunma)**, ve repo'daki `WIRESHARK_AG_ANALIZ_USTALIK_REHBERI.md` (ağ katmanı analizi), `OSINT_ARAC_SETI_USTALIK_REHBERI.md` (açık kaynak istihbarat).*
