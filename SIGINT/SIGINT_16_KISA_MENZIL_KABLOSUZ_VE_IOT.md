# SIGINT EL KİTABI — BÖLÜM 16: KISA MENZİLLİ KABLOSUZ EKOSİSTEM VE IoT GÜVENLİĞİ

## Bluetooth, RFID/NFC, Sub-GHz, Zigbee, LoRa ve Keyless — Yakındaki Spektrumu Anlamak ve Savunmak

> Amaç: Önceki bölümler uzak ve "büyük" sinyalleri ele aldı — uydular (Bölüm 11), GNSS (Bölüm 10), hücresel ağlar, radar ve ELINT (Bölüm 7). Bu bölüm tam tersine, kolunuzdaki saatten kapı kilidine, araç anahtarından akıllı prize kadar **metrelerle ölçülen menzilde** çalışan kablosuz ekosistemi ele alır. Bu cihazlar her gün etrafınızdadır, çoğu zaman görünmezdir ve toplu olarak modern saldırı yüzeyinin en geniş, en az denetlenen katmanını oluşturur. Hedefimiz operatör reçetesi değil; hangi teknolojinin hangi frekansta konuştuğunu, hangi aracın neyi yakaladığını, zafiyetin kavramsal kökenini ve **somut savunmayı** anlamaktır. Bir cihazı eline aldığında "bu nasıl konuşuyor, nerede zayıf, nasıl korunur?" diye düşünebilmen.

> Yasal çerçeve: Bu bölüm de serinin geri kalanı gibi **anlama, savunma ve spektrum okuryazarlığı** amaçlıdır. Anlatılan tüm gözlem teknikleri tasarım gereği **kendi cihazların ve yetkili test** içindir. Başkasının erişim kartını okumak/klonlamak, başkasının aracına/kapısına yetkisiz erişim, başkasının kablosuz trafiğini çözmek veya cihazını taklit etmek **çoğu ülkede suçtur** (Türkiye'de TCK 132–140, 243–245; ülkene göre teyit et). Bu metin hiçbir cihaza yetkisiz erişim için adım adım reçete vermez; saldırı sınıflarını yalnızca **prensip düzeyinde** ve **savunmayı çıkarabilmek için** açıklar. Klonlama yalnızca kendi kartında/cihazında veya yazılı yetkiyle yapılır.

---

## İÇİNDEKİLER

1. [Kısa Menzilli Spektruma Genel Bakış: Bant, Mesafe, Saldırı Yüzeyi](#1)
2. [Bluetooth Klasik ve BLE: Yığın, Reklam, Eşleşme, Adres Randomizasyonu](#2)
3. [BLE Sniffing ve Saldırı Yüzeyi Kavramı + Savunma](#3)
4. [RFID/NFC Temelleri: LF ve HF, Kart Mimarisi](#4)
5. [RFID/NFC Araçları: Proxmark3, Flipper Zero, ChameleonMini](#5)
6. [MIFARE Classic ve Crypto-1: Neden Eski Kartlar Zayıftır + Savunma](#6)
7. [Sub-GHz IoT: 433/868/915 MHz, OOK/FSK, Sabit vs Rolling-Code](#7)
8. [Zigbee ve Z-Wave (802.15.4 Ailesi): Mimari, Anahtar, Sniffing + Savunma](#8)
9. [LoRa ve LoRaWAN: Mimari, Sınıflar, OTAA/ABP, Bilinen Zayıflıklar + Savunma](#9)
10. [Otomotiv Keyless: RKE, PKES, Relay/Röle Saldırısı Prensibi + Savunma](#10)
11. [Diğer: Kablosuz Klavye/Fare (Mousejack), Drone, Tıbbi/Endüstriyel Telemetri](#11)
12. [Donanım Envanteri: Hangi Cihaz Neyi Yakalar (Tablo)](#12)
13. [Bütünleşik Savunma Doktrini: Şifreleme, Kimlik, Mesafe, Segmentasyon, İzleme](#13)
14. [Alıştırmalar (Yalnızca Kendi Cihazların / Yetkili)](#14)
15. [Hızlı Referans ve Diğer Bölümler](#15)

---

<a id="1"></a>
## 1. Kısa Menzilli Spektruma Genel Bakış: Bant, Mesafe, Saldırı Yüzeyi

Uzun menzilli sistemler (uydu, hücresel) güçlü vericiler, lisanslı bantlar ve görece az sayıda profesyonelce tasarlanmış protokol kullanır. Kısa menzilli ekosistem bunun tam karşıtıdır: düşük güçlü, çoğu lisanssız ISM (Industrial, Scientific, Medical) bantlarında çalışan, milyarlarca adet üretilmiş, **maliyet baskısı altında tasarlanmış** cihazlar. Bu üç özellik — düşük güç, açık bant, maliyet baskısı — saldırı yüzeyinin neden bu kadar geniş olduğunu doğrudan açıklar.

Düşük güç ve kısa menzil bir savunma sanılır ("zaten birkaç metre gidiyor"), ama bu yanıltıcıdır: yönlü anten ve düşük gürültülü yükselteç (LNA, Bölüm 3) ile alıcı menzili tasarım menzilinin kat kat üstüne çıkarılabilir. Bir BLE cihazının "10 metre" menzili, iyi bir antenle onlarca metreden dinlenebilir. Dolayısıyla "yakın olmak gerekiyor" varsayımına dayanan her güvenlik kararı baştan zayıftır; bu, bu bölümün tekrar tekrar döneceği bir temadır (özellikle keyless relay, Bölüm 10).

Aşağıdaki tablo, bu bölümde ele alınan teknolojileri tek bakışta konumlandırır. Frekanslar yaygın değerlerdir; bölgesel tahsis farkları için Bölüm 8'deki (frekans tahsisi) ve ulusal plana bakılmalı.

| Teknoloji | Tipik bant | Menzil (nominal) | Birincil kullanım | Şifreleme durumu (tipik) |
|---|---|---|---|---|
| Bluetooth Klasik | 2,4 GHz ISM | ~10–100 m | Ses, dosya, çevre birim | Eşleşmeye bağlı (zayıftan güçlüye) |
| BLE | 2,4 GHz ISM | ~10–50 m | Sensör, giyilebilir, beacon | LE Secure Connections varsa güçlü |
| RFID LF | 125 / 134 kHz | ~10 cm | Erişim kartı, hayvan çipi | Genelde yok (sadece UID) |
| RFID/NFC HF | 13,56 MHz | ~5–10 cm | Erişim, ödeme, ulaşım | Karta göre (yok → güçlü) |
| Sub-GHz IoT | 433/868/915 MHz | ~10–500 m | Kumanda, sensör, alarm | Sabit-kod (yok) → rolling-code |
| Zigbee | 2,4 GHz (+sub-GHz) | ~10–100 m | Ev otomasyonu, mesh | AES-128 (anahtar yönetimine bağlı) |
| Z-Wave | ~868/908 MHz | ~30–100 m | Ev otomasyonu, mesh | S2 ile güçlü |
| LoRaWAN | 433/868/915 MHz | ~2–15 km | Geniş alan IoT (LPWAN) | AES-128 (OTAA/ABP'ye bağlı) |
| Keyless (RKE/PKES) | ~315/433/868 MHz (+UWB/LF) | ~1–100 m | Araç giriş/çalıştırma | Rolling-code; relay'e açık |
| Kablosuz klavye/fare | 2,4 GHz (özel) | ~10 m | Giriş cihazları | Çoğu zaman zayıf/yok |

> Mühendislik sezgisi: Kısa menzilli güvenlik üç soruya indirgenir. **Şifreli mi?** (içeriği gizlemek) — bir saldırgan trafiği okuyabiliyor mu? **Kimlik doğrulanıyor mu?** (sahteciliği önlemek) — saldırgan meşru cihaz gibi davranabilir mi? **Tazelik/mesafe doğrulanıyor mu?** (tekrar ve relay'i önlemek) — eski bir mesaj yeniden oynatılabilir, ya da uzaktaki bir cihaz "yakın" gibi gösterilebilir mi? Bu üç soru, Bölüm 6 ve 13'teki "auth + nonce + mesafe" denetiminin kısa-menzil karşılığıdır ve bu bölümün her başlığında karşımıza çıkar.

---

<a id="2"></a>
## 2. Bluetooth Klasik ve BLE: Yığın, Reklam, Eşleşme, Adres Randomizasyonu

Bluetooth, "Bluetooth" tek bir şey değildir; aslında iki ayrı radyo teknolojisidir. **Bluetooth Klasik (BR/EDR)** ses ve sürekli yüksek hızlı veri için tasarlanmıştır (kulaklık, hoparlör, dosya). **Bluetooth Low Energy (BLE)**, 4.0 ile gelen, düşük güç ve kısa, aralıklı veri için tasarlanmış tamamen ayrı bir yığındır (sensör, giyilebilir, beacon, akıllı kilit). İkisi aynı 2,4 GHz ISM bandını paylaşır ama kanal planı, modülasyon ve protokol akışı farklıdır. Bu bölümde ağırlık BLE'dedir, çünkü IoT'nin baskın kısa-menzil protokolü odur.

### BLE kanal planı ve reklam (advertising)

BLE, 2,4 GHz bandını 2 MHz aralıklı 40 kanala böler. Bunların **üçü reklam kanalıdır** (kanal 37, 38, 39 — frekansları bilerek bandın WiFi'dan en az etkilenen yerlerine serpiştirilmiştir), kalan 37'si veri kanalıdır. Bir cihaz bağlanmadan önce, varlığını duyurmak için reklam paketlerini bu üç kanalda sırayla yayınlar. Bir telefon (tarayıcı/central) bu kanalları dinleyerek çevredeki cihazları keşfeder. Bu reklam mekanizması, BLE'nin hem en kullanışlı hem de gizlilik açısından en sızdıran yanıdır: cihaz, kimse bağlanmasa bile sürekli "buradayım" diye yayın yapar.

```
 BLE bandı (2,4 GHz, 40 kanal × 2 MHz):

 2402   2426            2480  MHz
  │      │               │
  ▼      ▼               ▼
 [37]  ...[38]...        [39]      ← 3 reklam kanalı (bandın kenarlarına+ortasına)
  │  data data data ...   │           serpilmiş; WiFi 1/6/11 ile çakışmayı azaltır
  └── 37 veri kanalı (bağlantı sırasında atlamalı kullanılır) ──┘

 Reklam akışı (zaman →):
  cihaz:  [ADV ch37][ADV ch38][ADV ch39] ...bekle (adv interval)... tekrar
           "buradayım, adım X, servislerim Y"
  tarayıcı: 37/38/39'u dinler → cihazı keşfeder
```

### Yığın: GAP ve GATT

BLE yazılım yığınının iki temel katmanını bilmek, hem geliştirme hem güvenlik için şarttır.

**GAP (Generic Access Profile)**, cihazın dünyaya nasıl göründüğünü ve nasıl bağlantı kurulduğunu yönetir: reklam yapma, tarama, rol (peripheral/central), bağlantı kurma. "Cihaz keşfi ve bağlantı" katmanıdır.

**GATT (Generic Attribute Profile)**, bağlantı kurulduktan sonra **veriyi nasıl alışveriş ettiğini** tanımlar. Veri, bir ağaç yapısında düzenlenir: **servisler** (örneğin "Kalp Atışı Servisi") içinde **karakteristikler** (örneğin "kalp atış değeri") bulunur; her karakteristiğin okunabilir/yazılabilir/bildirim özellikleri ve bir tanımlayıcısı (UUID) vardır. Bir BLE cihazını "anlamak", onun GATT ağacını dökmek (servis/karakteristik listesini çıkarmak) demektir.

```
 BLE GATT ağacı (kavramsal):

 Cihaz (peripheral)
  └─ Servis: Cihaz Bilgisi (UUID 0x180A)
  │    ├─ Karakteristik: Üretici Adı   [read]
  │    └─ Karakteristik: Firmware Sürümü [read]
  └─ Servis: Akıllı Kilit (özel UUID)
       ├─ Karakteristik: Kilit Durumu   [read, notify]
       └─ Karakteristik: Kilit Komutu   [write]   ← güvenlik açısından kritik nokta
```

Güvenlik refleksi şudur: yazılabilir bir komut karakteristiği (örneğin "kilit komutu") **kimlik doğrulama gerektirmiyorsa**, bağlanabilen herkes onu tetikleyebilir. Pek çok ucuz IoT cihazının zafiyeti tam buradadır — şifrelemeyi protokol seviyesinde değil, "uygulama bunu kimse fark etmez sanır" varsayımıyla atlamış olmalarıdır.

### Eşleşme (pairing) ve bağlanma (bonding)

Eşleşme, iki cihazın güvenli iletişim için anahtar üzerinde anlaşma sürecidir. BLE'nin tarihsel olarak en zayıf noktası buydu. **Legacy pairing** (4.0/4.1), anahtar değişimini zayıf bir yöntemle yapardı ve pasif bir dinleyici eşleşmeyi yakalarsa oturum anahtarını türetebilirdi. **LE Secure Connections** (4.2 ile gelen), eşleşmeyi **eliptik eğri Diffie-Hellman (ECDH)** üzerine kurar; pasif dinleyici artık anahtarı türetemez. Bu, BLE güvenliğindeki en önemli kırılma noktasıdır: bir cihazın LE Secure Connections kullanıp kullanmadığı, pasif sniffing'e karşı dayanıklılığını büyük ölçüde belirler.

Eşleşme sırasında "association model" (eşleşme doğrulama yöntemi) seçilir ve bu, ortadaki-adam (MITM) direncini belirler:

| Eşleşme yöntemi | Nasıl çalışır | MITM direnci | Tipik cihaz |
|---|---|---|---|
| Just Works | Kullanıcı doğrulaması yok, otomatik | **Yok** (en zayıf) | Ekransız ucuz cihazlar |
| Passkey Entry | Bir cihazda 6 hane gösterilir, diğerine girilir | Var | Ekran+tuş olan cihazlar |
| Numeric Comparison | İki ekranda aynı sayı gösterilir, onaylanır | Var | İki ekranlı cihazlar (LE SC) |
| Out of Band (OOB) | Anahtar başka kanaldan (NFC vb.) taşınır | Kanala bağlı güçlü | NFC destekli eşleşme |

"Just Works" yaygındır çünkü ucuzdur ve kullanıcı dostudur — ama tanım gereği MITM'e açıktır. Bonding ise eşleşmede türetilen anahtarların **kalıcı saklanmasıdır**; böylece cihazlar her seferinde yeniden eşleşmez. Bonding güvenliği, ilk eşleşmenin güvenliği kadardır.

### Adres randomizasyonu (gizlilik)

BLE cihazları sabit bir donanım adresi (public address) yayınlarsa, bu adres bir takip kimliğine dönüşür: cihaz nereye giderse adresi de gider, ve bir gözlemci reklam paketlerini dinleyerek o cihazı (ve sahibini) izleyebilir. Bunu önlemek için BLE, **rastgele özel adres (Resolvable Private Address, RPA)** mekanizması tanımlar: cihaz, periyodik olarak değişen rastgele bir adres yayınlar; yalnızca paylaşılan bir kimlik çözme anahtarına (IRK) sahip eşleşmiş cihaz, bu rastgele adresin arkasındaki gerçek kimliği çözebilir.

> Savunma sezgisi: Adres randomizasyonu doğru uygulanırsa güçlü bir gizlilik korumasıdır, ama **uygulama hataları** onu sık sık deler. Bazı cihazlar adresi yeterince sık değiştirmez, ya da reklam payload'ında (cihaz adı, üretici verisi) **adresten bağımsız sabit bir tanımlayıcı** sızdırır — adres rastgele olsa bile, payload'daki bu sabit alan cihazı yeniden tanınabilir kılar. Yani gizlilik, yalnızca adres katmanında değil, **tüm reklam içeriğinde** sağlanmalıdır. Bu, Bölüm 11'deki "şifreleme içeriği korur ama meta-veri sızar" dersinin kısa-menzil versiyonudur.

---

<a id="3"></a>
## 3. BLE Sniffing ve Saldırı Yüzeyi Kavramı + Savunma

Bir BLE cihazını **kendi cihazını** anlamanın yolu, reklamlarını ve (varsa) bağlantı trafiğini gözlemlemektir. Bu, tasarım gereği pasif bir gözlem etkinliğidir; aşağıdaki araçlar bunun için vardır.

### Sniffing araçları

| Araç | Ne yapar | Güçlü yanı | Sınırı |
|---|---|---|---|
| nRF52840 + Sniffle | Açık kaynak BLE sniffer firmware'i (Nordic çipi üzerinde) | Bağlantı takibi, ucuz, güncel BLE sürümleri | Tek anten; çoklu bağlantıda sınır |
| nRF Sniffer (Nordic) | Nordic'in resmî Wireshark eklentili sniffer'ı | Wireshark entegrasyonu, kararlı | Reklam+tek bağlantı odaklı |
| Ubertooth One | Açık kaynak 2,4 GHz donanımı (BT Klasik + BLE) | Klasik BT'ye de erişim, esnek | BLE'de modern sürümlerde Sniffle daha güçlü |
| Adafruit Bluefruit LE Sniffer | Hazır küçük donanım (nRF tabanlı) | Ucuz, başlangıç dostu | Temel düzey yetenek |
| Wireshark (BLE eklentisi) | Yakalanan paketleri çözümleme/görselleştirme | Standart analiz arayüzü | Kendi başına yakalamaz; sniffer'a bağlı |

Tipik akış: sniffer reklam kanallarını (37/38/39) dinler, kendi cihazının reklamlarını yakalar; eğer cihaz bağlanırsa, sniffer bağlantı parametrelerini (kanal atlama dizisi, zamanlama) takip ederek veri paketlerini de izler. Yakalanan paketler Wireshark'ta GATT seviyesine kadar çözülür.

```
 BLE sniffing topolojisi (kendi cihazın):

   [Kendi BLE cihazın]  ⇠⇠⇠ reklam (37/38/39) ⇢⇢⇢  [Sniffer: nRF52840+Sniffle]
        │ │                                                │
        │ └── (bağlanırsan) veri kanalı atlamalı ──────────┤
        │                                                  ▼
   [Telefon/uygulama] ⇠⇠ bağlantı ⇢⇢                  [Wireshark: GATT çözümleme]
                                                       "hangi servis, hangi karakteristik,
                                                        şifreli mi, auth var mı?"
```

### Saldırı yüzeyi — yalnızca kavram, savunma için

Bu sınıfları **savunmayı görebilmek için** kavramsal olarak sıralıyoruz; hiçbiri reçete değildir ve yetkisiz hedefe uygulanması suçtur.

**Pasif sniffing.** Reklam ve (legacy pairing'de) eşleşme trafiği dinlenir. Şifreleme zayıfsa içerik açığa çıkar. **Savunma:** LE Secure Connections (ECDH) ile eşleşme; pasif dinleyici anahtarı türetemez.

**Eşleşme zayıflıkları (MITM).** "Just Works" gibi doğrulamasız eşleşme, araya giren bir cihazın iki tarafı kandırmasına açıktır. **Savunma:** Passkey Entry veya Numeric Comparison gibi MITM-dirençli yöntem; kritik cihazlarda OOB.

**BLE replay.** Şifrelenmemiş veya tazelik (nonce/sayaç) içermeyen bir komut paketi kaydedilip yeniden gönderilebilir (örneğin "kilidi aç" komutu sabitse). **Savunma:** Uygulama katmanında nonce/sayaç + MAC; komutun her seferinde farklı ve doğrulanabilir olması (Bölüm 13'teki "tekrar edilebilen güvenli değildir").

**BLE relay.** Tıpkı keyless araçlarda olduğu gibi (Bölüm 10), iki nokta arasında BLE bağlantısı köprülenerek uzaktaki telefon "yakın" gösterilebilir; akıllı kilitlerde kritik bir risktir. **Savunma:** Mesafe/zaman doğrulama (RSSI tek başına zayıftır; gecikme/zaman ölçümü veya UWB daha güçlü), kullanıcı onayı.

**KNOB ve BLESA (kavramsal).** Literatürde raporlanmış iki sınıf: **KNOB** (Key Negotiation of Bluetooth) — eşleşmede anahtar entropisinin (uzunluğunun) pazarlığını zorlayıp zayıf anahtara düşürme fikri; **BLESA** (BLE Spoofing Attack) — yeniden bağlanma sırasında doğrulama eksiğinden yararlanıp sahtecilik fikri. Bunların ayrıntıları belirli akademik çalışmalara dayanır ve sürüm/yamaya göre değişir; **kesin etki ve etkilenen sürümler kaynaktan teyit edilmeli.** **Savunma:** güncel Bluetooth sürümü ve firmware, minimum anahtar uzunluğunun zorlanması, yeniden bağlanmada da kimlik doğrulama.

> Mühendislik sezgisi: BLE güvenliğinin tamamı üç soruya iner — eşleşme pasif dinleyiciye dayanıklı mı (LE SC)? Eşleşme MITM'e dayanıklı mı (doğrulama yöntemi)? Komutlar tekrar/relay'e dayanıklı mı (nonce + mesafe)? Bir akıllı kilit ya da giyilebilir cihazı değerlendirirken bu üç sorunun cevabı, güvenliğin neredeyse tamamını belirler. "BLE şifreli" demek yetmez; **hangi eşleşme, hangi doğrulama, hangi tazelik** sorulmalıdır.

---

<a id="4"></a>
## 4. RFID/NFC Temelleri: LF ve HF, Kart Mimarisi

RFID (Radio Frequency Identification), bir okuyucunun yaydığı alanla beslenen (çoğu zaman pilsiz/pasif) bir etiketin/kartın, kendi kimliğini veya verisini geri yansıtmasıdır. NFC (Near Field Communication), 13,56 MHz'de çalışan, RFID'nin bir alt kümesi sayılabilecek, iki yönlü ve kısa mesafeli (~birkaç cm) bir standarttır. Erişim kartları, ödeme, ulaşım, pasaport ve "telefonu okut" senaryolarının altında bu aile yatar.

Temel ayrım frekanstadır ve bu ayrım hem fiziği hem güvenliği belirler:

| Özellik | LF (Low Frequency) | HF (High Frequency) |
|---|---|---|
| Frekans | 125 kHz / 134 kHz | 13,56 MHz |
| Tipik menzil | ~10 cm (çok kısa) | ~5–10 cm |
| Eşleşme (coupling) | Manyetik (yakın alan) | Manyetik (yakın alan) |
| Tipik kartlar | EM4100, HID Prox | MIFARE Classic/DESFire, NTAG, NFC |
| Veri/güvenlik | Genelde sadece UID, **şifreleme yok** | Yoktan (UID) güçlüye (DESFire) geniş yelpaze |
| Tipik kullanım | Eski erişim, hayvan çipi, otopark | Modern erişim, ödeme, ulaşım, pasaport |

```
 RFID okuma fiziği (yakın alan, pasif kart):

   [OKUYUCU]  ~~~ manyetik alan (taşıyıcı) ~~~>  [KART]
       │  ← geri yansıma (load modulation) ────────┘   kart, alanı modüle
       │                                                ederek UID/veri yollar
       ▼
   "UID = 0x... , (HF ise) blok verisi / kimlik doğrulama"

   Menzil çok kısadır çünkü yakın-alan manyetik eşleşme mesafeyle
   çok hızlı zayıflar — ama "kısa menzil = güvenli" sanmak yanlıştır:
   güçlü okuyucu/anten ile mesafe bir miktar uzatılabilir.
```

### LF kartların mimarisi: çoğu zaman sadece bir kimlik

LF kartların büyük kısmı (EM4100 gibi) yalnızca sabit bir kimlik numarası (UID) yayınlar; **şifreleme, kimlik doğrulama, sayaç yoktur**. Okuyucu kartı alanla besler, kart UID'sini geri yansıtır, kapı bu UID'yi bir listede arar. HID Prox gibi bazıları tesis kodu + kart numarası taşır ama yine kriptografik koruma içermez. Bu mimarinin güvenlik sonucu açıktır: kartın tek "sırrı" UID'sidir ve UID yayındadır.

### HF kartların mimarisi: bloklar, sektörler, kimlik doğrulama

HF kartlar çok daha çeşitlidir. Basit NDEF/NTAG etiketleri (NFC link/kartvizit) yalnızca veri taşır. **MIFARE Classic**, hafızayı sektörlere ve bloklara böler; her sektörün erişimi anahtarlarla (Crypto-1 ile) korunur — ama bu koruma tarihsel olarak kırılmıştır (Bölüm 6). **MIFARE DESFire** (EV1/EV2/EV3) ise standart, güçlü kriptografi (AES, 3DES) ve gerçek karşılıklı kimlik doğrulama kullanır; modern güvenli kartların temelidir.

Bir kartı "anlamak", önce **tipini** belirlemek (UID'nin yapısından, ATQA/SAK yanıtından, yongadan), sonra mimarisini (sektör/blok düzeni, kullanılan kimlik doğrulama) çıkarmaktır. Savunma değerlendirmesi tam da burada başlar: "bu kart kriptografik kimlik doğrulama yapıyor mu, yoksa sadece bir UID mi?"

---

<a id="5"></a>
## 5. RFID/NFC Araçları: Proxmark3, Flipper Zero, ChameleonMini

Bu araçlar **kendi kartını** okumak, türünü/zafiyetini belirlemek ve (yalnızca kendi kartında veya yazılı yetkiyle) öykünmek/klonlamak için kullanılır. Başkasının kartını okumak/kopyalamak suçtur ve bu metin bunun adımlarını vermez.

| Araç | LF | HF | Öne çıkan yetenek | Konum |
|---|---|---|---|---|
| Proxmark3 (RDV4/Easy) | ✔ | ✔ | Düşük seviye RFID analiz/oku/öykün; en güçlü açık platform | RFID araştırmasının standardı |
| Flipper Zero | ✔ (125 kHz) | ✔ (13,56 MHz) | Taşınabilir çoklu-araç; kolay kullanım, sınırlı derinlik | Hobici/saha; SDR değil |
| ChameleonMini / Tiny | — | ✔ | HF kart **öykünme** (emulation) odaklı | Kart davranışı taklidi/test |
| ACR122U vb. NFC okuyucu | — | ✔ | Masaüstü NFC okuma/yazma | Geliştirme/test |

**Proxmark3**, RFID'nin "osiloskobu" gibidir: ham alan seviyesinde sinyali görür, kartın tipini otomatik tanır, sektör/blok yapısını döker, ve desteklenen kartlarda kimlik doğrulama denemelerini yönetir. En düşük seviyeye erişim sağladığı için hem en güçlü analiz hem en bilgili savunma değerlendirmesini mümkün kılar.

**Flipper Zero**, çok daha erişilebilir bir taşınabilir cihazdır; LF (125 kHz) ve HF (13,56 MHz) kartları okuyabilir, bazı tipleri kaydedip öykünebilir. Sınırı, Proxmark3'ün düşük-seviye esnekliğine ve hesap gücüne sahip olmamasıdır; bilinen/desteklenen kart tipleriyle pratiktir, derin analiz aracı değildir. Aynı cihaz sub-GHz ve IR de yapar (Bölüm 7).

**ChameleonMini/Tiny**, asıl olarak bir **kart öykünücüdür**: bir HF kartın davranışını taklit ederek okuyucu/sistem testini mümkün kılar — kendi sisteminin bir kartı nasıl karşıladığını yetkili biçimde sınamak için.

```
 RFID iş akışı (kendi kartın / yetkili):

  1) ALGILA   → kartın tipi ne? (LF mi HF mi; EM4100 / Prox / MIFARE Classic / DESFire / NTAG)
       │          (ATQA/SAK, UID yapısı, frekans yanıtı)
       ▼
  2) ANALİZ   → mimari: sadece UID mi? sektör/blok + kimlik doğrulama mı? hangi kripto?
       │
       ▼
  3) DEĞERLENDİR → zafiyet: şifreleme yok mu? Crypto-1 (zayıf) mı? AES/DESFire (güçlü) mü?
       │
       ▼
  4) (yalnızca kendi kartında/yetkiyle) öykün/klonla → "bu kart kopyalanabilir mi?" dersini gör
```

> Yasal sınır: 4. adım — öykünme/klonlama — yalnızca **kendi kartında** veya **yazılı yetkiyle** yapılır. Başkasının erişim kartını okumak, kaydetmek veya kopyalamak (kartı bir an için ele geçirip okumak dahil) yetkisiz erişim ve bilişim suçudur. Amaç "kendi eski kartım kopyalanabiliyor mu, değiştirmeli miyim?" sorusunu yetkili biçimde yanıtlamaktır.

---

<a id="6"></a>
## 6. MIFARE Classic ve Crypto-1: Neden Eski Kartlar Zayıftır + Savunma

MIFARE Classic, dünyada milyarlarca adet üretilmiş, erişim ve ulaşım sistemlerinin uzun süre belkemiği olmuş bir HF kart ailesidir. Güvenliğini **Crypto-1** adlı tescilli (proprietary) bir akış şifresine dayandırır. Bu bölüm, neden eski MIFARE Classic kartların artık güvenli sayılmadığını **prensip düzeyinde** açıklar; amaç bir kırma reçetesi değil, "neden eski kartı değiştirmeliyim" mühendislik dersidir.

### Crypto-1 neden zayıf — kavramsal

Crypto-1, 48-bit bir doğrusal geri beslemeli kaydırma yazmacına (LFSR) dayanan, kapalı tasarlanmış bir akış şifresidir. İki tarihsel gerçek onu çökertti:

Birincisi, **gizlilikle güvenlik (security through obscurity)** başarısızlığı. Algoritma gizli tutularak korunmaya çalışıldı; ama 2000'lerin sonunda akademik araştırmacılar (donanım tersine mühendisliği ve kriptanaliz ile) algoritmayı açığa çıkardı. Açığa çıkınca, tasarımdaki yapısal zayıflıklar görünür oldu.

İkincisi, **kriptografik zayıflıklar**. 48-bitlik durum uzayı modern hesap gücü için küçüktür; dahası, anahtar akışı üretimindeki ve kimlik doğrulamadaki yapısal kusurlar (zayıf rastgelelik, durum sızıntısı), anahtarların kaba kuvvetten çok daha hızlı türetilmesine olanak verdi. Akademik literatürde, bilinen kart/okuyucu etkileşimlerinden anahtarların pratik sürelerde elde edilebildiği gösterilmiştir (yöntem ayrıntıları akademik kaynaklarda; burada **prensip** veriliyor, **adım değil**).

Sonuç şudur: bir MIFARE Classic kartın sektör anahtarları kırılabilirse, kartın tüm içeriği okunabilir ve **birebir bir kopyası** (klon) üretilebilir. Erişim sistemi yalnızca kartın UID'sine veya içeriğine bakıyorsa, klon orijinal gibi davranır.

```
 Neden eski kart zayıf (kavramsal zincir):

 [Tescilli Crypto-1]  ──açığa çıktı──►  [48-bit küçük durum + yapısal kusur]
        │                                          │
   "gizlilik = güvenlik" yanılgısı          kriptanaliz pratik hale geldi
        │                                          │
        └──────────────► [sektör anahtarları elde edilebilir] ──► [kart klonlanabilir]
                                                                        │
                                          erişim sistemi sadece UID/içeriğe bakıyorsa:
                                                          klon = orijinal gibi geçer
```

### Klonlama: prensip + yasal sınır

Klonlamanın prensibi basittir: kartın okunabilen tüm bilgisini (UID — bazı "magic" kartlarda yazılabilir — ve blok içeriği) yeni bir karta yazmak. Bir sistem yalnızca bu okunabilir bilgiye güveniyorsa, klon ayırt edilemez. **Bu, tam olarak eski kartların neden değiştirilmesi gerektiğinin kanıtıdır.**

> Yasal sınır — vurgulu: Klonlama yalnızca **kendi kartında** ya da **yazılı yetkiyle** yapılır; amaç kendi sisteminin/kartının zafiyetini görmektir. **Başkasının erişim kartını klonlamak suçtur** (yetkisiz erişim, bilişim sistemine girme, hırsızlık teşebbüsü). Bir kişinin kartını izinsiz okumak için bir an yakınına getirmek bile (skimming) yetkisiz erişimdir. Bu metin klonlama adımlarını vermez; yalnızca "neden zayıf" prensibini ve savunmayı verir.

### Savunma: güçlü karta geçmek

| Önlem | Ne sağlar | Not |
|---|---|---|
| MIFARE DESFire EV2/EV3'e geçiş | AES tabanlı, gerçek karşılıklı kimlik doğrulama | Modern erişimin temeli; Crypto-1 sorununu kökten çözer |
| Kriptografik karşılıklı kimlik doğrulama | Kart ve okuyucu birbirini kanıtlar; salt-UID yetmez | Klon, anahtarı bilmeden geçemez |
| Çeşitlenen/dinamik veri (CMAC, sayaç) | Her okuma farklı/doğrulanabilir; replay zorlaşır | Statik içeriğe güvenmeme |
| Çok faktör (kart + PIN/biyometri) | Tek başına kart yeterli olmaz | Kritik kapılarda |
| Anti-passback, izleme, kart yaşam döngüsü | Çalıntı/klon kullanımını tespit/kısıtla | Operasyonel katman |

> Mühendislik sezgisi: MIFARE Classic dersi, kriptografinin altın kuralını somutlaştırır — **güvenlik gizlilikten değil, açık ve sağlam kriptografiden gelir.** Tescilli, gizli, kısa-anahtarlı bir şema er ya da geç açığa çıkar ve çöker. Bir erişim kartını değerlendirirken tek soru şudur: "kart, okuyucuyu ve okuyucu kartı, **paylaşılan bir sırrı açığa vurmadan** kanıtlıyor mu (gerçek karşılıklı kimlik doğrulama), yoksa sadece bir kimlik mi gösteriyor?" İkincisiyse, kopyalanabilir.

---

<a id="7"></a>
## 7. Sub-GHz IoT: 433/868/915 MHz, OOK/FSK, Sabit vs Rolling-Code

1 GHz altındaki ISM bantları (bölgeye göre 433 MHz, 868 MHz Avrupa, 915 MHz Amerika) garaj kumandaları, kapı zilleri, alarm sensörleri, hava istasyonları, akıllı prizler ve sayısız basit IoT cihazıyla doludur. Bu bantların cazibesi: düşük frekansta daha iyi menzil ve duvar geçişi, ucuz radyo. Modülasyon çoğunlukla en basitlerden ikisidir — **OOK** (On-Off Keying: taşıyıcıyı aç/kapa) ve **FSK** (iki/çok frekans arası). Bölüm 1 bu modülasyonların fiziğini, Bölüm 4/5 ise demodülasyonu verir; burada güvenlik açısından kritik ayrım **kod tipindedir.**

### Sabit-kod vs rolling-code — güvenliğin kalbi

| Özellik | Sabit-kod (fixed code) | Rolling-code (hopping code) |
|---|---|---|
| Davranış | Her basışta **aynı** kod gönderilir | Her basışta **farklı**, sırayla artan kod |
| Temel | Sabit bir bit deseni | Paylaşılan gizli anahtar + sayaç (örn. KeeLoq) |
| Replay'e karşı | **Savunmasız** — bir kez yakalanan kod tekrar oynatılabilir | Dirençli — yakalanan kod bir sonraki basışta geçersiz |
| Tipik cihaz | Eski garaj kumandaları, ucuz uzaktan kumandalar | Modern garaj kumandaları, araç anahtarları |
| Ders | Değiştirilmeli | Doğru uygulanmışsa kabul edilebilir |

**Sabit-kod** cihazlar, "kapıyı aç" için her zaman aynı sinyali gönderir. Bir gözlemci bu sinyali bir kez kaydedip yeniden gönderirse kapı açılır — klasik **replay** (Bölüm 13). Bu yüzden sabit-kod cihazlar (kendi cihazında gözlemleyebileceğin gibi) bugün güvensiz kabul edilir.

**Rolling-code** (KeeLoq bunun bilinen bir örneğidir), her basışta paylaşılan gizli anahtar ve artan bir sayaçtan türetilen **farklı** bir kod gönderir. Alıcı, beklenen sayaç penceresindeki kodları kabul eder; kullanılan/eskimiş kod reddedilir. Böylece basit replay etkisizdir. KeeLoq'un kavramsal ayrıntısı ve tarihsel zafiyetleri Bölüm 6 ve 13'te ele alınır; oradaki tasarım dersi şudur: koruma "rolling" etiketinde değil, **kriptografik kalitede ve doğru sayaç/nonce yönetimindedir** — zayıf rastgelelik veya kötü sayaç yönetimi olan bazı eski uygulamalarda tarihsel zafiyetler raporlanmıştır.

```
 Sabit-kod (zayıf):           Rolling-code (güçlü):
  basış 1: 1011001              basış 1: f(anahtar, sayaç=42) = A7F3...
  basış 2: 1011001  (aynı!)     basış 2: f(anahtar, sayaç=43) = 9C12...
  basış 3: 1011001  (aynı!)     basış 3: f(anahtar, sayaç=44) = 5E80...
       │                              │
  yakalanan kod tekrar             yakalanan kod bir sonraki
  oynatılabilir → replay           basışta zaten geçersiz
```

### Flipper Zero sub-GHz: yeteneği ve sınırı

Flipper Zero, sub-GHz bandında (yaygın olarak ~300–928 MHz aralığında, bölgesel yasal sınırlar dahilinde) OOK/FSK sinyalleri **alabilir, çözebilir ve kendi cihazların için kaydedebilir.** Pek çok sabit-kod protokolünü tanır. **Sınırı** ve **etik tasarımı** önemlidir: cihaz, rolling-code sistemleri taklit ederek araç/kapı açmak için tasarlanmamıştır; güncel firmware bu tür kötüye-kullanımı kısıtlar, ve zaten rolling-code'un doğası gereği yakalanan kod tekrar işe yaramaz. Flipper, **kendi sabit-kod cihazını analiz etmek** (sabit mi rolling mi ayrımını görmek) için pratik bir araçtır; bir saldırı aracı olarak değil.

### rtl_433 ile kendi cihazların

`rtl_433`, bir RTL-SDR ile sub-GHz bandını dinleyip **bilinen yüzlerce cihaz protokolünü** (hava istasyonları, sıcaklık/nem sensörleri, lastik basıncı, bazı kumandalar) otomatik çözen açık kaynak bir araçtır. Bu, serinin baştan beri vurguladığı **yasal, kendi cihazların** alıştırmasının kalbidir: evdeki kendi termometre/sensör/istasyonunu okuyup verisini görmek (Bölüm 5'te de geçer).

```
 rtl_433 ile kendi sensörünü okumak (kavramsal):

  [Kendi hava istasyonu sensörün] ──433/868 MHz OOK/FSK──> [RTL-SDR] ──> rtl_433
                                                                            │
                                                                            ▼
                                       "Acurite-XYZ, kanal 2, sıcaklık 21.4°C, nem %48"
                                       (kendi cihazının verisini çözer; yasal, pasif)
```

> Savunma sezgisi: Sub-GHz'de güvenlik tek soruya iner — **kod her seferinde değişiyor ve doğrulanıyor mu?** Sabit-kod (replay'e açık) ile rolling-code (replay'e dirençli) arasındaki fark, bütün savunmayı belirler. Kendi eski bir garaj kumandanı gözlemleyip "her basışta aynı mı?" diye bakmak, bu dersi elle deneyimletir (Bölüm 14 ve Bölüm 13 alıştırması). Sabit-kodsa: değiştir. Rolling-code'sa: uygulamanın kriptografik kalitesine güven, ama bilinen zafiyetli ürünlerden kaçın.

---

<a id="8"></a>
## 8. Zigbee ve Z-Wave (802.15.4 Ailesi): Mimari, Anahtar, Sniffing + Savunma

Akıllı ev ve bina otomasyonunun (ışık, sensör, kilit, termostat) baskın iki ağ teknolojisi Zigbee ve Z-Wave'dir. İkisi de düşük güçlü, **mesh** (her düğüm başkasının mesajını iletebilir) ağlardır; bu mesh yapısı menzili genişletir ama saldırı yüzeyini de dağıtır.

### Zigbee mimarisi (IEEE 802.15.4 üzerinde)

Zigbee, fiziksel/MAC katmanı için IEEE 802.15.4 standardını kullanır (yaygın olarak 2,4 GHz, ayrıca sub-GHz bantlar). Ağda üç rol vardır:

| Rol | İşlev | Güç |
|---|---|---|
| Coordinator | Ağı kuran, anahtarları yöneten tek düğüm | Sürekli besli |
| Router | Mesajları ileten, mesh'i genişleten düğüm | Genelde besli |
| End device | Yaprak cihaz (sensör, anahtar); uyur | Pille çalışır |

Güvenliğin kalbi **anahtar yönetimidir**. Zigbee, AES-128 ile şifreleme yapar ama gerçek güvenlik, anahtarların ağa **nasıl dağıtıldığına** bağlıdır. Tarihsel zafiyet noktası şudur: cihaz ağa katılırken (joining), ağ anahtarı bazı eski/uyumluluk modlarında **iyi bilinen bir varsayılan "transport key" ile** veya zayıf korumayla taşınabiliyordu; bu kısa katılım anını dinleyen biri ağ anahtarını ele geçirebilir ve sonrasında tüm trafiği çözebilirdi. Modern Zigbee (3.0) bunu **install code** (cihaza özgü, kutudan gelen benzersiz başlangıç anahtarı) ile büyük ölçüde kapatır.

```
 Zigbee mesh + anahtar dağıtımı (kritik an):

        [Coordinator] ── ağ anahtarı yönetir
         /     |     \
    [Router][Router][Router]      ← mesh; menzili genişletir
      /  \      |       \
 [End][End]  [End]    [End]       ← yaprak cihazlar (uyur, pille)

 KRİTİK AN: yeni cihaz katılırken (joining) ağ anahtarı taşınır.
   Eski/zayıf modda: tahmin edilebilir transport key → katılımı dinleyen
                     ağ anahtarını alır → tüm mesh çözülür.
   Zigbee 3.0: install code (cihaza özgü) → katılım anı korunur.
```

### Sniffing — kavram ve araçlar (kendi ağın)

Kendi Zigbee ağını anlamak/savunma değerlendirmek için sniffing araçları vardır; bunlar pasif gözlem içindir.

| Araç/çerçeve | Ne yapar | Not |
|---|---|---|
| CC2531 / CC2652 USB dongle | 802.15.4 paketlerini yakalar (ucuz, yaygın) | Wireshark'a beslenir |
| KillerBee (çerçeve) | 802.15.4/Zigbee güvenlik analiz araç seti | Yakalama+çözümleme kavramsal iskelet |
| ApiMote | KillerBee uyumlu özel donanım | Araştırma odaklı |
| Ubertooth | Esas 2,4 GHz BT; bazı 802.15.4 gözlemi | Birincil Zigbee aracı değil |
| Wireshark (802.15.4/Zigbee) | Paket çözümleme; anahtar verilirse şifre çözme | Kendi ağının anahtarıyla kendi trafiğini incele |

Kendi ağında tipik akış: bir 802.15.4 dongle'ı kanalını ağınla eşitler, paketleri yakalar; kendi ağ anahtarını Wireshark'a girerek **kendi** şifreli trafiğini çözüp "katılım nasıl korunuyor, hangi anahtar modu" sorusunu yanıtlarsın.

### Z-Wave — benzer mantık

Z-Wave, ~868 MHz (Avrupa) / ~908 MHz (Amerika) sub-GHz bandında çalışan, yine mesh tabanlı bir otomasyon teknolojisidir. Güvenlik açısından tarih benzerdir: eski **S0** güvenlik çerçevesi, katılım sırasında zayıf bir anahtar değişimi nedeniyle eleştirilmişti; yerini alan **S2**, eşleşmede güçlü (ECDH tabanlı) anahtar kurulumu ve cihaza özgü PIN/DSK ile bu sorunu büyük ölçüde kapatır. Ders aynıdır: **şifreleme var olması yetmez; anahtarın ağa nasıl ve ne kadar güvenli girdiği** belirleyicidir.

> Savunma sezgisi (Zigbee/Z-Wave ortak): Mesh otomasyonda güvenliğin tek kritik anı **cihazın ağa katıldığı andır.** Şifreleme (AES-128) zaten vardır; sorun, ağ anahtarının o ilk anda nasıl korunduğudur. Savunma: install code / S2 destekleyen (Zigbee 3.0, Z-Wave S2) cihazlar kullan; cihaz katılımını (pairing) **güvenli, kısa, kontrollü** bir anda yap; mümkünse IoT ağını ana ağdan segmente et (Bölüm 13). Bilinen-varsayılan-anahtar modlarına izin veren eski cihazlardan kaçın.

---

<a id="9"></a>
## 9. LoRa ve LoRaWAN: Mimari, Sınıflar, OTAA/ABP, Bilinen Zayıflıklar + Savunma

LoRa, **uzun menzilli (kilometrelerce), çok düşük güçlü** bir radyo modülasyonudur (chirp spread spectrum — chirp yayılı spektrum); sub-GHz ISM bantlarında (433/868/915 MHz) çalışır. **LoRaWAN**, LoRa fiziği üzerine kurulu bir ağ protokolüdür ve LPWAN (Low Power Wide Area Network) dünyasının önde gelenidir: az veriyi çok uzağa, pille yıllarca taşıyan sensör ağları. Önceki kısa-menzil teknolojilerinden menzil olarak ayrılır ama IoT güvenlik mantığı aynıdır.

### Mimari ve sınıflar

LoRaWAN dört temel bileşenden oluşur: **uç cihaz** (sensör) → **gateway** (anteni geniş; LoRa'yı IP'ye köprüler) → **ağ sunucusu** (network server; trafiği yönetir) → **uygulama sunucusu** (application server; veriyi tüketir). Uç cihazlar üç sınıfa ayrılır:

| Sınıf | Davranış | Gecikme/güç dengesi |
|---|---|---|
| A | Sadece kendi gönderiminden sonra kısa alım penceresi | En düşük güç, en yüksek aşağı-yön gecikmesi |
| B | A + zamanlanmış periyodik alım pencereleri (beacon) | Orta |
| C | Sürekli alım (gönderim hariç hep dinler) | En düşük gecikme, en yüksek güç |

### Güvenlik: iki anahtar, iki katman

LoRaWAN güvenliği AES-128 üzerine kuruludur ve (1.0.x'te) iki anahtar etrafında döner: **NwkSKey** (ağ oturum anahtarı — bütünlük/MIC, ağ katmanı) ve **AppSKey** (uygulama oturum anahtarı — yük şifreleme, uçtan uca). Bu oturum anahtarları iki yöntemle elde edilir:

**OTAA (Over-The-Air Activation).** Cihaz, kök anahtar(lar) (1.0.x'te **AppKey**; 1.1'de **NwkKey** + **AppKey**) ile bir katılım (join) prosedürü yürütür; oturum anahtarları her katılımda **taze nonce'larla türetilir.** Bu tercih edilen, daha güvenli yöntemdir çünkü oturum anahtarları dinamiktir.

**ABP (Activation By Personalization).** Oturum anahtarları cihaza **sabit olarak gömülür**; katılım prosedürü yoktur. Daha basittir ama anahtarlar statik olduğundan ve kötü uygulamalarda **sayaç (frame counter) yönetimi** hatalıdır — bu, replay riskine ve anahtar yeniden-kullanımına kapı açar.

```
 LoRaWAN mimari + aktivasyon:

  [Uç cihaz] ⇢ LoRa (km'ler) ⇢ [Gateway] ⇢ IP ⇢ [Ağ Sunucusu] ⇢ [Uygulama Sunucusu]
      │                                              │                    │
   AppSKey (yük şifreleme, uçtan uca) ───────────────┼────────────────────┘
   NwkSKey (MIC/bütünlük, ağ katmanı) ───────────────┘

  OTAA (tercih):  kök anahtar + join → taze nonce → dinamik oturum anahtarları
  ABP (basit):    oturum anahtarları sabit gömülü → sayaç yönetimi kritik (zayıf nokta)
```

### Bilinen zayıflıklar — kavramsal

LoRaWAN'ın raporlanmış zayıflık temaları (ayrıntılar akademik/üretici kaynaklarda; **teyit edilmeli**):

- **ABP nonce/sayaç yönetimi:** frame counter sıfırlanır veya kötü yönetilirse, eski paketler yeniden oynatılabilir (replay) ya da anahtar akışı yeniden kullanılır. OTAA bunu taze join-nonce ile azaltır.
- **Anahtar yönetimi/tedarik:** kök anahtarlar (AppKey) cihaza güvensiz gömülürse veya tedarik zincirinde sızarsa, tüm güvenlik çöker. Donanımdan anahtar çıkarımı (fiziksel erişimle) bir tehdittir.
- **1.0.x'te tek kök anahtar:** 1.0.x'te ağ ve uygulama güvenliği tek AppKey'e dayanırken, **1.1** kök anahtarları (NwkKey/AppKey) ayırarak ağ ile uygulama güvenini ayrıştırır ve roaming senaryolarını güçlendirir.

> Savunma sezgisi: LoRaWAN'da güvenlik üç karara iner. **OTAA kullan** (ABP yerine) — dinamik oturum anahtarları replay'i ve anahtar yeniden-kullanımını azaltır. **Kök anahtarları güvenli sakla** — cihazda güvenli element/secure storage, tedarik zincirinde sızdırma. **Frame counter'ı doğru yönet** — sıfırlama/atlama izleyip reddet. Mümkünse LoRaWAN **1.1**'i tercih et (ayrık kök anahtarlar). IoT'nin "uzun menzil" kolu olsa da, ders kısa-menzille aynıdır: şifreleme var; belirleyici olan anahtar yönetimi ve tazeliktir.

---

<a id="10"></a>
## 10. Otomotiv Keyless: RKE, PKES, Relay/Röle Saldırısı Prensibi + Savunma

Modern araç anahtarları kısa-menzilli kablosuz güvenliğin en görünür ve en yüksek-bahisli örneğidir. İki ayrı işlevi karıştırmamak gerekir:

**RKE (Remote Keyless Entry) — uzaktan kilit.** Anahtardaki butona basınca, araca tek yönlü bir RF komutu (kilit/aç) gider (yaygın olarak ~315/433/868 MHz). Bu bir **rolling-code** komutudur (Bölüm 7): butona her basışta sayaç+anahtardan türetilen farklı bir kod. Replay'e karşı tasarlanmıştır.

**PKES (Passive Keyless Entry and Start) — pasif giriş.** Burada buton yoktur: anahtar cebinizdeyken araca yaklaşıp kapıyı çekmeniz veya çalıştırma düğmesine basmanız yeter. Araç, anahtarın **yakında** olduğunu, kısa menzilli bir sorgu-yanıt (challenge-response) ile doğrular: araç bir LF (~125 kHz) "uyandırma/sorgu" alanı yayar, yakındaki anahtar bunu duyup bir RF yanıtla cevaplar. İşte bu "yakınlık varsayımı" relay saldırısının hedefidir.

### Relay (röle / menzil uzatma) saldırısı — prensip

Relay saldırısı **kodu kırmaya çalışmaz**; meşru sinyali olduğu gibi **uzatır.** PKES'te araç ile anahtar yalnızca birbirine yakınken konuşabilmelidir; relay, bu konuşmayı iki nokta arasında köprüleyerek aracın anahtarı "yakında sanmasını" sağlar. Kod hiç kırılmaz; yalnızca **mesafe yapay olarak uzatılır.**

```
 PKES relay (röle) saldırısı — prensip (yalnızca kavram):

   [ARAÇ] ── LF sorgu ──► (cihaz A) ⇠⇠⇠⇠ uzun köprü ⇢⇢⇢⇢ (cihaz B) ── LF ──► [ANAHTAR]
      ▲                                                                          │
      │ ◄────────── RF yanıt, köprü ile geri taşınır ◄────────────────────────────┘
      │
   Araç "anahtar 1 metrede" sanır; oysa anahtar evde, araç sokakta.
   KOD KIRILMAZ — sadece "yakınlık" yalanı söylenir.
```

Bu, RKE'nin rolling-code'unu **atlar**, çünkü rolling-code yalnızca komutun tazeliğini garanti eder; **mesafeyi** garanti etmez. PKES'in zaafı kriptografik değil, **fizikseldir**: "yakın görünmek" ile "fiziksel olarak yakın olmak" arasındaki farkı ölçmemesi.

### Savunma — mesafeyi gerçekten ölçmek

| Önlem | Ne sağlar | Not |
|---|---|---|
| **UWB (Ultra-Wideband) time-of-flight** | Sinyalin gidiş-dönüş süresinden **gerçek mesafe** ölçer; relay köprüsünün eklediği gecikme "yakın" iddiasını çürütür | Modern güvenli keyless'in yöneldiği çözüm |
| **Hareket sensörü (anahtarda)** | Anahtar bir süre hareketsizse (gece, masada) yaymayı durdurur | Relay penceresini kapatır |
| **Faraday kesesi/kutusu** | Anahtarı saran ekranlama, RF'i tamamen keser | Ucuz, etkili kullanıcı önlemi |
| **PKES'i kapatma seçeneği** | Pasif girişi devre dışı bırakıp butona dönmek | Bazı araçlarda mümkün |
| **Daha kısa zamanlama pencereleri** | Yanıt için izin verilen gecikmeyi sıkma | Relay'i zorlaştırır (tek başına yetersiz) |

Kritik savunma fikri tek cümledir: **relay'i yenmenin yapısal yolu, "yakın görünmek" ile "fiziksel olarak yakın olmak" arasındaki farkı ölçmektir.** UWB tabanlı mesafe ölçümü tam da bunu yapar (gidiş-dönüş süresinden gerçek mesafe); köprü ek gecikme eklediği için "yakın" iddiası fiziksel olarak çürür. Bu, Bölüm 13'te detaylandırılan aynı doktrindir.

### Rolling-code ve TPMS gizliliği — iki ek not

**Rolling-code** (RKE tarafı): yukarıda ve Bölüm 7/13'te işlendiği gibi, doğru uygulanmış rolling-code basit replay'i çözer; ders, korumanın "rolling" etiketinde değil kriptografik kalitede olduğudur.

**TPMS (Tire Pressure Monitoring System) gizliliği:** her lastikteki basınç sensörü, ~315/433 MHz'de periyodik olarak benzersiz bir sensör kimliği yayar. Bu, bir güvenlik açığından çok bir **gizlilik/takip** sorunudur: sabit TPMS kimlikleri, bir aracı (dolayısıyla sahibini) yol kenarındaki bir alıcıyla **tanıma/izleme** imkânı verebilir. Savunma kullanıcı elinde sınırlıdır (üretici tasarımına bağlı); farkındalık düzeyinde bilinmesi, "araç bile sürekli kimlik yayar" gerçeğini gösterir — Bölüm 2'deki BLE adres-randomizasyonu ve Bölüm 11'deki meta-veri sızıntısı temasının otomotiv karşılığıdır.

> Yasal sınır — vurgulu: **Kendi aracının** anahtar/kumanda RF imzasını yasal olarak gözlemleyip "rolling mi, hangi bant" diye analiz edebilirsin. Başkasının aracına relay/röle ile girmek, kod yakalamak veya çalıştırmak **hırsızlık ve bilişim suçudur**; bu metin hiçbir relay kurulum adımı vermez. Burada verilen tek şey **prensip + savunma**dır: aracını UWB destekli sistemle tercih et, anahtarını Faraday kesesinde/hareketsizken-yaymayan modda tut.

> Mühendislik sezgisi (Bölüm 13 ile ortak): "Tekrar edilebilen ve doğrulanmayan şey güvenli değildir" — sabit kod tekrar edilebilir (replay → çözüm: rolling/nonce); kimlik doğrulanmayan komut sahtelenebilir (spoof → çözüm: kripto kimlik); **ölçülmeyen mesafe yanıltılabilir (relay → çözüm: gerçek mesafe ölçümü, UWB).** Keyless, bu üçlünün en somut sahnesidir.

---

<a id="11"></a>
## 11. Diğer: Kablosuz Klavye/Fare (Mousejack), Drone, Tıbbi/Endüstriyel Telemetri

Kısa-menzilli ekosistem yalnızca yukarıdakilerle bitmez; saldırı yüzeyi günlük çevre birimlerine ve kritik telemetriye kadar uzanır. Bu başlık, kapsamı tamamlayan üç alana **farkındalık** düzeyinde değinir.

### Kablosuz klavye/fare — Mousejack kavramı

Pek çok kablosuz klavye/fare (Bluetooth değil, üreticinin **özel 2,4 GHz** dongle protokolüyle çalışanlar) tarihsel olarak güvenliği zayıf tasarlamıştır. **Mousejack** olarak adlandırılan zafiyet ailesi, bazı dongle'ların **fare hareketlerini şifrelemeden** ve **gönderenin kimliğini yeterince doğrulamadan** kabul etmesine dayanır; sonuçta saldırgan, kurbanın dongle'ına sahte tuş-vuruşu paketleri **enjekte ederek** (kavramsal olarak) bilgisayara komut yazdırabilir. Kritik nokta: klavye trafiği bazen şifreliyken, fare trafiği şifresiz bırakılmış ve aynı alıcı sahte klavye paketlerini de kabul edebilmiştir.

| Boyut | Zayıf tasarım | Savunma |
|---|---|---|
| Şifreleme | Fare (bazen klavye) trafiği şifresiz | Şifreli giriş protokolü kullanan cihaz seç |
| Kimlik doğrulama | Dongle, kaynağı doğrulamadan paket kabul eder | Eşleştirmeyi doğrulayan/şifreleyen dongle |
| Yama | Eski firmware açık kalır | Üretici firmware güncellemesi (varsa) |
| Alternatif | — | **Bluetooth (LE SC)** veya **kablolu** giriş cihazı |

> Savunma sezgisi: Bir giriş cihazı (klavye/fare) bilgisayara komut **yazabildiği** için yüksek değerli bir hedeftir. Çözüm basittir: trafiği **şifreleyen ve kaynağı doğrulayan** cihazlar (modern Bluetooth LE Secure Connections destekli, ya da güvenli özel protokol) tercih et; mümkünse kritik makinelerde kablolu kullan; üretici güvenlik yamalarını uygula. Bu, "şifreleme + kimlik doğrulama" temel ikilisinin (Bölüm 13) en sade uygulamasıdır.

### Drone RF bağları — Bölüm 13'e referans

Tüketici dronlarının kumanda (uplink) ve video (downlink) bağları (yaygın olarak 2,4/5,8 GHz, kimi zaman şifresiz veya zayıf korumalı) ayrı bir saldırı yüzeyidir; jamming, link kopması ve spoofing açısından **Bölüm 13 (Drone/İHA RF güvenliği ve Counter-UAS)** detaylı işler. Buradaki bağ: drone da kısa-menzilli kablosuz ekosistemin bir üyesidir ve aynı "şifreleme/kimlik/dayanıklılık" sorularına tabidir.

### Tıbbi ve endüstriyel telemetri — yüksek bahis, eski tasarım

İki alan özel farkındalık gerektirir, çünkü bahis **can güvenliği** ve **kritik altyapıdır**:

**Tıbbi cihazlar:** kablosuz insülin pompaları, kalp pilleri/defibrilatörler, hasta telemetrisi tarihsel olarak zayıf/şifresiz kablosuz bağlarla raporlanmıştır. Akademik çalışmalar, bazı eski cihazlarda kimlik doğrulama eksikliğinin tehlikeli senaryolara kapı açtığını göstermiştir. Bu, üretici ve düzenleyici sorumluluğudur; kullanıcı tarafında **farkındalık** ve güncel/yamalı cihaz tercihi anlamlıdır.

**Endüstriyel telemetri (SCADA/ICS kablosuz):** saha sensörleri, vana/aktüatör kontrolü ve uzak terminal birimleri (RTU) bazen şifresiz veya zayıf korumalı sub-GHz/özel radyolarla konuşur. Bölüm 13 (Kritik Altyapı RF Riski) bunu işler; ders şudur: eski endüstriyel kablosuz, **kimlik doğrulama ve şifreleme** olmadan tasarlandığında manipülasyona açıktır ve burada bahis fizikseldir (üretim, enerji, su).

> Mühendislik sezgisi: Bu üç alanın ortak dersi şudur — **bir kablosuz bağın bahsi yükseldikçe (giriş cihazı → drone → tıbbi/endüstriyel), zayıf tasarımın bedeli de yükselir.** Aynı üç soru (şifreli mi, kimlik doğrulanıyor mu, tazelik/dayanıklılık var mı) her katmanda geçerlidir; yalnızca yanlış cevabın sonucu, bir tuş enjeksiyonundan can güvenliğine kadar ağırlaşır.

---

<a id="12"></a>
## 12. Donanım Envanteri: Hangi Cihaz Neyi Yakalar (Tablo)

Bu bölümdeki tüm gözlem alıştırmaları belirli donanım gerektirir. Aşağıdaki envanter, "neyi anlamak/yetkili test etmek için hangi cihaz" sorusunu tek tabloda yanıtlar. Tüm kullanım **kendi cihazların/yetkili** çerçevesindedir.

| Cihaz | Hangi teknoloji | Bant | Birincil yeteneği | Sınırı / not |
|---|---|---|---|---|
| **RTL-SDR (Blog V4)** | Sub-GHz IoT, genel RX | ~24 MHz–1,7 GHz | rtl_433 ile kendi sensörlerini çöz; sub-GHz gözlem | Sadece RX; 2,4 GHz BLE'yi doğrudan iyi yakalamaz |
| **HackRF One** | Geniş RX/TX | 1 MHz–6 GHz | Geniş bant gözlem; TX (yasal sınır!) | BLE/Zigbee için özel araçlar daha iyi |
| **nRF52840 + Sniffle** | BLE | 2,4 GHz | BLE reklam + bağlantı takibi (en güncel) | Tek anten; çok-bağlantıda sınır |
| **nRF Sniffer (Nordic)** | BLE | 2,4 GHz | Resmî BLE sniffer + Wireshark | Reklam/tek bağlantı odaklı |
| **Ubertooth One** | BT Klasik + BLE | 2,4 GHz | Klasik BT gözlemi, esnek 2,4 GHz | BLE'de Sniffle çoğu işte daha güçlü |
| **Adafruit Bluefruit Sniffer** | BLE | 2,4 GHz | Ucuz başlangıç BLE sniffer | Temel düzey |
| **Proxmark3 (RDV4)** | RFID/NFC | 125 kHz + 13,56 MHz | LF+HF oku/analiz/öykün; en derin RFID aracı | Öğrenme eğrisi; öykünme yalnızca kendi kartın |
| **Flipper Zero** | RFID/NFC + Sub-GHz + IR | 125 kHz / 13,56 MHz / ~300–928 MHz | Taşınabilir çoklu-araç; kendi cihaz analizi | SDR değil; rolling-code/derin analizde sınırlı |
| **ChameleonMini/Tiny** | NFC (HF) | 13,56 MHz | HF kart öykünme (test) | Esas öykünücü; okuma/analizde sınırlı |
| **CC2531 / CC2652 dongle** | Zigbee/802.15.4 | 2,4 GHz | Zigbee paket yakalama (Wireshark) | Kendi ağ anahtarınla kendi trafiğin |
| **ApiMote + KillerBee** | Zigbee/802.15.4 | 2,4 GHz | 802.15.4 güvenlik analiz çerçevesi | Araştırma odaklı |
| **LoRa geliştirme kartı/gateway** | LoRa/LoRaWAN | 433/868/915 MHz | Kendi LoRaWAN cihaz/ağını test | Bölgesel bant uyumu gerek |

> Pratik not: Tek bir cihaz her şeyi yapmaz. Kabaca: **2,4 GHz BLE** için nRF52840+Sniffle; **RFID/NFC** için Proxmark3 (derin) veya Flipper (pratik); **sub-GHz kendi sensörlerin** için RTL-SDR+rtl_433; **Zigbee** için bir 802.15.4 dongle. "Hangi cihaz" sorusunun cevabı her zaman "hangi bandı ve hangi protokolü" sorusunun cevabıdır (Bölüm 1-2).

---

<a id="13"></a>
## 13. Bütünleşik Savunma Doktrini: Şifreleme, Kimlik, Mesafe, Segmentasyon, İzleme

Bölüm boyunca her teknolojinin kendi savunması verildi; burada hepsini **tek bir doktrine** indiriyoruz. Kısa-menzilli kablosuz/IoT güvenliği, beş katmanlı bir denetim listesidir. Bu liste, Bölüm 6 (savunma) ve Bölüm 13'teki (RF tehdit) doktrinin IoT karşılığıdır.

```
 KISA-MENZİL/IoT SAVUNMA YIĞINI (yukarıdan aşağı sorulur):

 ┌─────────────────────────────────────────────────────┐
 │ 1. ŞİFRELEME      → içerik gizli mi? (pasif sniff'e) │ ← LE SC, AES-128, DESFire
 ├─────────────────────────────────────────────────────┤
 │ 2. KİMLİK DOĞR.   → karşı taraf kanıtlı mı? (MITM/   │ ← karşılıklı auth, MAC,
 │                     spoof'a)                          │   eşleşme doğrulama
 ├─────────────────────────────────────────────────────┤
 │ 3. TAZELİK        → mesaj taze mi? (replay'e)         │ ← nonce/sayaç + reddetme
 ├─────────────────────────────────────────────────────┤
 │ 4. MESAFE         → gerçekten yakın mı? (relay'e)     │ ← UWB ToF, hareket sensörü
 ├─────────────────────────────────────────────────────┤
 │ 5. SEGMENT+İZLEME → kuşatılmış ve görünür mü?         │ ← IoT VLAN, anomali izleme
 └─────────────────────────────────────────────────────┘
```

### Beş katman, somut karşılıklar

| Katman | Hangi tehdide karşı | Bu bölümdeki somut karşılığı |
|---|---|---|
| **Şifreleme** | Pasif sniffing, içerik okuma | BLE LE Secure Connections; Zigbee/LoRa AES-128; DESFire kart |
| **Kimlik doğrulama** | MITM, spoofing, sahte cihaz | BLE Passkey/Numeric; kartta karşılıklı auth; dongle doğrulaması |
| **Tazelik (nonce/sayaç)** | Replay (tekrar oynatma) | Rolling-code; LoRaWAN frame counter; BLE komut nonce |
| **Mesafe doğrulama** | Relay/röle (menzil uzatma) | Keyless UWB ToF; Faraday kese; BLE relay'e karşı zaman ölçümü |
| **Segmentasyon + izleme** | Yanal hareket, görünürlük kaybı | IoT'yi ayrı VLAN/ağ; pairing'i kontrollü anda; anomali tespiti |

### Kesişen ilkeler

**Güncel firmware ve protokol sürümü.** Bu bölümdeki zafiyetlerin neredeyse tamamı (legacy BLE pairing, Zigbee varsayılan anahtar, Z-Wave S0, MIFARE Classic, Mousejack) **eski sürümlerde** yaşar ve yeni sürümlerde (LE SC, Zigbee 3.0 install code, Z-Wave S2, DESFire EV2/3, LoRaWAN 1.1) büyük ölçüde kapanır. **"Eski cihazı değiştir/güncelle"**, tek başına en yüksek getirili savunmadır.

**Mesafe-bağlı kimlik doğrulama.** "Yakınlık" bir güvenlik kararının temeli olacaksa, mesafe **ölçülmeli** (UWB ToF), **varsayılmamalıdır**. Bu, relay ailesinin (keyless, BLE kilit) yapısal çözümüdür.

**Segmentasyon.** IoT cihazları (tanımı gereği yama-zayıf, uzun ömürlü, güveni düşük) ana ağdan ve hassas sistemlerden **ayrılmalıdır** (ayrı VLAN/SSID). Bir IoT cihazı ele geçse bile yanal hareket sınırlanır.

**İzleme.** Beklenmeyen pairing, ani trafik değişimi, bilinmeyen cihazın belirmesi tespit edilebilmelidir. Bu, Bölüm 11'deki meta-veri/anomali mantığının savunma tarafıdır.

> Mühendislik sezgisi: Bir kısa-menzilli/IoT cihazı değerlendirirken sırayla sor — **şifreli mi? kimlik doğrulanıyor mu? tazelik var mı? mesafe gerçekten ölçülüyor mu? kuşatılmış ve izleniyor mu?** Bu beş soru, bu bölümün bütün teknolojilerini (BLE, RFID, sub-GHz, Zigbee, LoRa, keyless) tek bir denetim listesinde birleştirir. Yanıtların çoğu "hayır"sa, cihaz değiştirilmeli veya kuşatılmalıdır.

---

<a id="14"></a>
## 14. Alıştırmalar (Yalnızca Kendi Cihazların / Yetkili)

> Bu alıştırmalar yalnızca **kendi cihazların** ve **yetkili test** içindir. Hiçbiri yetkisiz erişim, başkasının kartını/aracını/ağını hedef alma veya klonlama içermez. Gözlem ve analiz odaklıdır; iletim/karıştırma yoktur. Klonlama/öykünme yalnızca kendi kartında ya da yazılı yetkiyle yapılır. Şüphedeysen yapma; bandını ve ülkeni teyit et.

### A) Kendi BLE cihazının reklamlarını gözlemlemek (GAP/GATT refleksi)

Bir nRF52840+Sniffle (veya nRF Sniffer / Ubertooth) ile reklam kanallarını (37/38/39) dinle ve kendi bir BLE cihazını (giyilebilir, akıllı priz, kendi telefonunun bir uygulaması) yakala. Şu satırları doldur:

| Gözlem | Değer / not |
|---|---|
| Reklam aralığı (adv interval) | ? ms |
| Yayınlanan adres türü (public / random) | ? |
| Reklam payload'ında sabit tanımlayıcı var mı? | ? (cihaz adı, üretici verisi sabit mi?) |
| Bağlanınca GATT: yazılabilir komut karakteristiği var mı? | ? |
| Eşleşme türü (Just Works / Passkey / Numeric)? | ? |

Amaç: Cihazın **gizlilik** (adres rastgele mi, payload sabit kimlik sızdırıyor mu — Bölüm 2) ve **güvenlik** (yazılabilir komut auth gerektiriyor mu, eşleşme MITM-dirençli mi — Bölüm 3) duruşunu kendi gözünle değerlendirmek.

### B) Kendi erişim kartının türünü ve zayıflığını belirlemek (RFID değerlendirme refleksi)

Bir Proxmark3 veya Flipper Zero ile **kendi** bir erişim/ulaşım kartını (kendi ofis kartın, kendi otopark kartın — yalnızca senin olan) oku ve sınıflandır:

```
 Doldur:
  Frekans:   LF (125 kHz) mi  /  HF (13,56 MHz) mi ?
  Tip:       EM4100 / HID Prox / MIFARE Classic / DESFire / NTAG / ? 
  Mimari:    sadece UID mi  /  sektör+blok+kimlik doğrulama mı ?
  Kripto:    yok  /  Crypto-1 (zayıf)  /  AES/DESFire (güçlü) ?
  Sonuç:     "bu kart klonlanabilir mi? değiştirmeli miyim?"
```

Amaç: Bölüm 4-6'nın dersini somutlaştırmak — kartın yalnızca bir UID mi gösterdiğini yoksa **gerçek karşılıklı kimlik doğrulama** mı yaptığını belirleyip, eski/zayıf bir kartsa değiştirme kararını gerekçelendirmek. **Klonlama yalnızca kendi kartında ve dersi görmek için**; başkasının kartı kesinlikle hariç.

### C) Kendi sub-GHz kumandanı analiz etmek (sabit vs rolling refleksi)

Bir RTL-SDR + URH (Universal Radio Hacker) veya Flipper Zero ile **kendi** bir sub-GHz kumandanı (kendi garaj/araç kumandan, kendi kablosuz zilin) gözlemle. Aynı butona birkaç kez bas ve her basışın sinyalini karşılaştır:

| Basış | Sinyal görünümü (kabaca) | Aynı mı, farklı mı? |
|---|---|---|
| 1 | ? | — |
| 2 | ? | (1 ile aynı / farklı) |
| 3 | ? | (öncekiyle aynı / farklı) |

Yorumla: Her basışta görünüş **aynıysa** sabit-kod ailesine (replay'e açık → değiştir), **değişiyorsa** rolling-code ailesine işaret eder (Bölüm 7). Bu, "neden eski sabit-kod cihazlarımı değiştirmeliyim" dersini elle deneyimletir. **Çok önemli:** yalnızca **gözlem**; kodları kaydedip tekrar oynatma (replay) amacı taşımaz — kendi cihazında bile amaç sabit/rolling ayrımını görmektir, kod yakalamak değil. Başkasının cihazı hariç.

### D) Kendi Zigbee/IoT ağının katılım güvenliğini incelemek (anahtar yönetimi refleksi)

Bir CC2531/CC2652 dongle + Wireshark ile **kendi** Zigbee ağını (kendi akıllı ev kurulumun) dinle. Kendi ağ anahtarını Wireshark'a girerek kendi trafiğini çöz ve şunları yanıtla:

1. Bir cihaz ağa katılırken (pairing) anahtar nasıl korunuyor — install code mu, eski/varsayılan mod mu (Bölüm 8)?
2. Ağ AES-128 ile mi şifreli; çözülen trafikte komutlar açık mı görünüyor (kendi anahtarınla)?
3. IoT ağın ana ağdan/hassas sistemlerden segmente mi (Bölüm 13)?

Amaç: Şifrelemenin (AES-128) **var olmasının** yetmediğini, asıl belirleyicinin **katılım anındaki anahtar yönetimi** ve **segmentasyon** olduğunu kendi ağında görmek.

### E) Kendi sisteminin beş-soru denetimini uygulamak (savunma doktrini refleksi)

İletim/yakalama olmadan, kavramsal bir egzersiz: Evindeki/işindeki bir kısa-menzilli cihaz seç (akıllı kilit, keyless araç, akıllı priz, kablosuz klavye) ve Bölüm 13'ün beş sorusunu kâğıda uygula:

1. **Şifreli mi?** (içerik pasif sniffing'e karşı korunuyor mu)
2. **Kimlik doğrulanıyor mu?** (karşı taraf kanıtlı mı, MITM/spoof'a karşı)
3. **Tazelik var mı?** (nonce/sayaç; replay'e karşı)
4. **Mesafe gerçekten ölçülüyor mu?** (relay'e karşı — özellikle keyless/akıllı kilit)
5. **Segmente ve izleniyor mu?** (yanal hareket sınırlı mı, anomali görünür mü)

Amaç: Bu bölümün bütün teknolojilerini tek bir savunma denetim listesinde birleştirmek; her "hayır" cevabını somut bir iyileştirmeye (cihaz değiştir, firmware güncelle, Faraday kese, VLAN'a al) bağlamak.

---

<a id="15"></a>
## 15. Hızlı Referans ve Diğer Bölümler

### Kavram kartı

| Kavram | Bir cümlelik öz |
|---|---|
| Bluetooth Klasik vs BLE | Klasik = ses/sürekli; BLE = düşük güç/aralıklı sensör; ikisi ayrı yığın, aynı 2,4 GHz |
| BLE reklam (37/38/39) | Bağlanmadan önce varlık duyurma; en kullanışlı ama en sızdıran katman |
| GAP / GATT | GAP = keşif/bağlantı; GATT = servis/karakteristik veri ağacı |
| LE Secure Connections | ECDH tabanlı eşleşme; pasif dinleyici anahtarı türetemez (kritik kırılma) |
| BLE eşleşme yöntemi | Just Works (MITM'e açık) ↔ Passkey/Numeric/OOB (dirençli) |
| Adres randomizasyonu (RPA) | Takibi önler; ama payload sabit kimlik sızdırırsa delinir |
| RFID LF vs HF | LF (125 kHz, çoğu sadece UID, şifresiz) ↔ HF (13,56 MHz, UID'den DESFire'a) |
| Crypto-1 / MIFARE Classic | Tescilli, 48-bit, açığa çıkmış zayıf şifre; eski kart klonlanabilir → değiştir |
| DESFire EV2/EV3 | AES tabanlı gerçek karşılıklı kimlik doğrulama; modern güvenli kart |
| Sabit-kod vs rolling-code | Sabit = replay'e açık (değiştir); rolling = farklı/sayaçlı kod (dirençli) |
| Zigbee/Z-Wave kritik anı | Cihazın ağa katılımı; install code (Zigbee 3.0) / S2 (Z-Wave) korur |
| LoRaWAN OTAA vs ABP | OTAA = dinamik oturum anahtarı (tercih); ABP = statik, sayaç yönetimi kritik |
| Keyless RKE vs PKES | RKE = butonlu rolling-code; PKES = pasif yakınlık → relay'e açık |
| Relay (röle) saldırısı | Kodu kırmaz, mesafeyi yapay uzatır; çözüm UWB time-of-flight |
| Mousejack | Zayıf kablosuz fare/klavye dongle'ına tuş enjeksiyonu; çözüm şifreli+auth giriş |
| Beş-soru denetimi | Şifreleme? kimlik? tazelik? mesafe? segment+izleme? |

### Ezber sezgiler

- Kısa menzil güvenliği üç soruya iner: **şifreli mi, kimlik doğrulanıyor mu, tazelik/mesafe doğrulanıyor mu?**
- "Yakın olmak gerekiyor" bir güvenlik değildir; yönlü anten+LNA menzili açar, relay mesafeyi yalanlar.
- BLE'de belirleyici: **LE Secure Connections** (pasif sniff'e) + **eşleşme doğrulama** (MITM'e) + **komut nonce'u** (replay'e).
- RFID'de tek soru: kart **gerçek karşılıklı kimlik doğrulama** mı yapıyor, yoksa sadece bir **UID** mi gösteriyor?
- MIFARE Classic dersi: **güvenlik gizlilikten değil, açık ve sağlam kriptografiden gelir** — Crypto-1 açığa çıkınca çöktü.
- Sub-GHz'de tek ayrım: kod her basışta **aynı mı (sabit → replay)** yoksa **farklı mı (rolling → dirençli)**?
- Zigbee/Z-Wave/LoRaWAN: AES var olması yetmez; belirleyici **anahtarın ağa/cihaza nasıl girdiğidir** (katılım/aktivasyon).
- Keyless: rolling-code tazeliği verir ama **mesafeyi vermez**; relay'in çözümü gerçek mesafe ölçümüdür (UWB).
- Eski sürüm = zafiyet, yeni sürüm = kapanma; **"güncelle/değiştir"** en yüksek getirili savunmadır.
- IoT'yi **segmente et ve izle**; yama-zayıf, uzun ömürlü, güveni düşük cihazları ana ağdan ayır.

### Yasal sınır ve perspektif (daima)

Bu bölümdeki tüm teknikler tasarım gereği **kendi cihazların ve yetkili test** içindir. Başkasının kartını okumak/klonlamak, aracına/kapısına relay ile girmek, ağını/cihazını izinsiz dinlemek veya taklit etmek çoğu ülkede suçtur (Türkiye'de TCK 132–140 haberleşme gizliliği, 243–245 bilişim suçları; ülkene göre teyit et). Klonlama/öykünme yalnızca kendi kartında ya da yazılı yetkiyle yapılır. Hedef operatörlük değil, **mühendislik sezgisi ve savunma**dır: bir kablosuz cihazı eline aldığında onun nasıl konuştuğunu, nerede zayıf olduğunu ve nasıl korunacağını görebilmek. Emin olmadığın teknik ayrıntıları (özellikle KNOB/BLESA, Crypto-1 yöntemi, LoRaWAN nonce, KeeLoq) **kaynaktan teyit et**; bu metin anlama ve savunma içindir.

---

### Serinin diğer bölümleri (çapraz referans)

- **SIGINT_01** — RF Fiziği, Spektrum ve Modülasyon: ISM bantları, OOK/FSK/chirp, dB/menzil mantığı. (Bu bölümdeki tüm modülasyon ve bant kavramlarının fiziksel temeli orada.)
- **SIGINT_02** — SDR Cihazları Derinlemesine: RTL-SDR, HackRF, Flipper Zero karşılaştırması. (Bu bölümün donanım envanterinin ayrıntılı hali orada.)
- **SIGINT_03** — Antenler, Donanım & Devre (LNA, yön bulma): "kısa menzil neden uzatılabilir" — yönlü anten + LNA mantığı orada.
- **SIGINT_05** — Protokoller & Sinyal Çözümleme: rtl_433 ile kendi sensörlerin, WiFi/BT/Zigbee imzaları, IoT spektrumu. (Bu bölümün çözümleme/protokol tarafı orada.)
- **SIGINT_06** — Güvenlik, Açıklar & Savunma: replay/spoofing prensibi, savunma doktrini. (Bu bölümün savunma çerçevesinin kökü orada.)
- **SIGINT_07** — Disiplinler ve Sinyal Ayıklama: SEI/RF parmak izi (cihazları fiziksel kusurdan ayırma), trafik analizi. (BLE adres-randomizasyonu ve meta-veri sızıntısının ayıklama bağlamı orada.)
- **SIGINT_08** — Frekans Tahsisi ve Bant Planı: ISM bantlarının bölgesel tahsisi (433/868/915 farkı). (Bu bölümdeki bant tablolarının resmî dayanağı orada.)
- **SIGINT_13** — RF Tehdit ve Karşı-Önlemler: replay/spoofing/relay derinlemesine, KeeLoq, keyless UWB savunması, drone RF, kritik altyapı. (Bu bölümün keyless/relay ve drone başlıklarının ana referansı orada.)

> Kapanış: Kısa menzilli kablosuz ekosistem, modern hayatın görünmez sinir sistemidir — kapı kilidinden araç anahtarına, akıllı prizden tıbbi pompaya. Toplu olarak en geniş ve en az denetlenen saldırı yüzeyini oluşturur, çünkü çoğu cihaz maliyet baskısı altında, "kısa menzil zaten güvenlidir" yanılgısıyla tasarlanmıştır. Ama her teknolojinin güvenliği aynı beş soruya iner: şifreli mi, kimlik doğrulanıyor mu, taze mi, gerçekten yakın mı, kuşatılmış ve izlenir mi? Bu denetim listesini içselleştirdiğinde, eline aldığın her kablosuz cihaz artık bir kara kutu değil; nasıl konuştuğunu, nerede zayıf olduğunu ve nasıl savunulacağını gördüğün tanıdık bir sistemdir. Bir sonraki adım, bu sezgiyi yalnızca **kendi cihazlarının yasal gözlemi** üzerinde sınamaktır.
>
> Bu doküman Kanije Kalesi güvenlik/teknik rehberleri koleksiyonunun SIGINT serisinin 16. bölümüdür. İlgili: SIGINT_01–13, `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`, `VERACRYPT_USTALIK_REHBERI.md`.
