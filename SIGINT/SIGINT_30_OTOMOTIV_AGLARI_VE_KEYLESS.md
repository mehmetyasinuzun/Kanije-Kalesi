# SIGINT EL KİTABI — BÖLÜM 30: OTOMOTİV-İÇİ AĞLAR VE ANAHTARSIZ SİSTEMLER — DERİN

## CAN/LIN/FlexRay/Otomotiv Ethernet, ECU Mimarisi, RKE/PKES/Immobilizer, Relay Saldırısı ve Katmanlı Savunma

> Amaç: Bölüm 16, kısa menzilli kablosuz ekosistemi (BLE, RFID, Sub-GHz, Zigbee) genel bir çerçevede ele aldı ve keyless/relay'i kavramsal düzeyde tanıttı; Bölüm 13 aynı relay/replay/spoof üçlüsünü RF tehdit doktrini açısından işledi. Bu bölüm tek bir hedefe iner: **otomobil.** Modern bir araç, tekerlekli bir bilgisayar ağıdır — onlarca ECU, birkaç farklı veriyolu, lisanssız RF bantlarında konuşan anahtarlar ve giderek artan bir uzaktan-bağlanabilirlik yüzeyi. Bu bölüm, aracın içindeki ağların (CAN bus başta) nasıl çalıştığını ve **neden güvensiz tasarlandığını**, anahtarsız giriş/çalıştırma sistemlerinin (RKE, PKES, immobilizer) mimarisini, bunlara yönelik saldırı sınıflarını **yalnızca prensip ve "neden mümkün" düzeyinde** ve her birinin somut savunmasını ele alır. Hedef, bir aracı çalmak/açmak için reçete değil; bir otomotiv güvenlik araştırmacısının, tasarımcının ve bilinçli sürücünün "bu sistem nerede zayıf, nasıl korunur?" sorusunu yanıtlayabilmesidir.

> Yasal çerçeve (kritik, mutlaka oku): Bu bölüm tasarım gereği **otomotiv güvenliği eğitimi ve savunma** içindir. Anlatılan tüm pratik gözlem ve analiz teknikleri **yalnızca kendi aracında veya yazılı olarak yetkilendirilmiş bir test ortamında** uygulanır. Başkasının aracına yetkisiz erişim, kapısını açma, çalıştırma, anahtar sinyalini relay/röle ile uzatma, CAN veriyoluna izinsiz müdahale veya immobilizer atlatma; Türkiye'de hırsızlık (TCK 141–142), bilişim sistemine yetkisiz erişim ve sistemi engelleme/bozma (TCK 243–245) ve haberleşmenin gizliliğini ihlal (TCK 132–140) kapsamında **ağır suçtur**; ülkene göre teyit et. Bu metin hiçbir araca yetkisiz erişim, kilit açma veya çalıştırma için adım-adım operasyonel reçete vermez. Saldırı sınıfları yalnızca **kavramsal köken + neden mümkün + savunma** ekseninde verilir. Emin olunmayan teknik ayrıntılar (özellikle KeeLoq/Hitag2/Rolling-PWN/RollJam'in belirli ürünlerdeki etkisi, ISO 15118 ve V2X güvenlik detayları) **kaynaktan teyit edilmelidir** ve bölümde böyle işaretlenir.

---

## İÇİNDEKİLER

1. [Tekerlekli Bilgisayar: Otomotiv Saldırı Yüzeyine Genel Bakış](#1)
2. [ECU Mimarisi ve Araç-İçi Ağ Topolojisi](#2)
3. [CAN Bus Derinlemesine: Çerçeve, Arbitrasyon, Bit Rate](#3)
4. [CAN Neden Güvensiz Tasarlandı: Broadcast ve Kimlik-Doğrulamasızlık](#4)
5. [CAN-FD, LIN, FlexRay, MOST ve Otomotiv Ethernet](#5)
6. [OBD-II Portu ve CAN Erişimi (Kendi Araç / Yetkili)](#6)
7. [SocketCAN ve can-utils: Pasif Okuma ve Analiz Kavramı](#7)
8. [CAN Savunması: IDS, Mesaj Kimlik Doğrulama, Segmentasyon, Gateway](#8)
9. [Anahtarsız Sistemler I — RKE: Buton, UHF, Rolling Code](#9)
10. [Anahtarsız Sistemler II — PKES: LF Uyandırma, UHF Yanıt, Yakınlık](#10)
11. [Immobilizer ve Transponder: RFID 125 kHz / 13,56 MHz](#11)
12. [Saldırı Sınıfları I — Sabit-Kod Replay ve Rolling-Code Zayıflıkları](#12)
13. [Saldırı Sınıfları II — Relay (Röle) Saldırısı: Prensip ve Savunma](#13)
14. [Saldırı Sınıfları III — CAN Enjeksiyon ve Immobilizer Atlatma Kavramı](#14)
15. [TPMS: Lastik Basıncı İzleme ve Araç-İzleme Gizliliği](#15)
16. [Modern Saldırı Yüzeyi: Telematik, EV Şarj, V2X](#16)
17. [Savunma Derinliği: Katmanlı Mimari, ISO/SAE 21434, Sorumlu Açıklama](#17)
18. [Alıştırmalar (Yalnızca Kendi Araç / Yasal)](#18)
19. [Hızlı Referans ve Diğer Bölümler](#19)

---

<a id="1"></a>
## 1. Tekerlekli Bilgisayar: Otomotiv Saldırı Yüzeyine Genel Bakış

Modern bir otomobili anlamanın en yararlı yolu, onu mekanik bir nesne değil, **dağıtık bir gömülü bilgisayar ağı** olarak görmektir. Orta sınıf bir araçta 30–70, üst segmentte 100'ü aşan sayıda elektronik kontrol ünitesi (ECU) motoru, frenleri, direksiyonu, gövde elektroniğini, eğlence sistemini ve sürüş yardımcılarını yönetir; bunlar birkaç farklı veriyolu üzerinden konuşur ve toplam yazılım hacmi onlarca milyon satır koda ulaşır.

Bu mimari güvenlik açısından üç temel gerilim taşır. Birincisi **tarihsel miras**: araç-içi ağ protokollerinin çoğu 1980'ler–1990'larda, ağın **fiziksel olarak kapalı ve güvenilir** olduğu varsayımıyla tasarlandı; o dönemde tehdit modeli kötü niyetli bir düğüm değil, elektriksel gürültü/arızaydı. İkincisi **bağlanabilirlik patlaması**: bir zamanlar izole bu ağ bugün hücresel modem (telematik), Bluetooth, WiFi, USB, anahtar RF'i ve şarj arabirimi üzerinden dış dünyaya açıldı. Üçüncüsü **uzun ömür ve yama zorluğu**: araç 15–20 yıl yolda kalır; güvenlik yamasının (özellikle OTA olmayan araçlarda) sahaya ulaşması yıllar alır ya da hiç ulaşmaz.

Saldırı yüzeyini iki büyük eksende düşünmek faydalıdır. **Fiziksel/yerel erişim ekseni**: OBD-II portu, far/tampon arkasındaki kablo demetleri, USB, doğrudan ECU pinleri. Bu eksen, araca fiziksel temas (veya çok yakınlık) gerektirir. **Uzaktan erişim ekseni**: anahtar RF'i (RKE/PKES), TPMS, telematik (hücresel), Bluetooth/WiFi, V2X, EV şarj iletişimi. Bu eksen, fiziksel temas olmadan menzil dahilinden ulaşılabilir ve bu yüzden hırsızlık ve büyük ölçekli tehdit açısından daha kritiktir.

```
 Otomotiv saldırı yüzeyi (kavramsal harita):

   UZAKTAN ERİŞİM EKSENİ                  YEREL/FİZİKSEL ERİŞİM EKSENİ
   Anahtar RF (RKE 315/433/868)           OBD-II portu (CAN'e doğrudan)
   PKES (LF 125 kHz + UHF), TPMS          Far/tampon arkası kablo demeti
   Telematik (hücresel/eSIM)              USB / medya portu
   Bluetooth/WiFi, V2X, EV şarj           ECU debug/JTAG; veriyolu kesme
        └──────────────┬───────────────────────┘
                       ▼
        [ARAÇ-İÇİ AĞLAR: CAN / CAN-FD / LIN / FlexRay / Ethernet]
        [ECU'lar: motor, fren, direksiyon, gövde, ADAS, infotainment]
```

Bölümün omurgası: önce **iç ağı** (özellikle CAN) ve yapısal güvensizliğini anlamak; sonra **dış kapıları** (anahtarsız sistemler, TPMS, telematik) ve bunların iç ağa nasıl köprülendiğini görmek; en sonunda her katman için **savunmayı** çıkarmak. Bütün saldırı sınıfları, Bölüm 13 ve 16'da yerleşen tek bir denetim listesine indirgenir: **şifreli mi, kimlik doğrulanıyor mu, taze mi (replay'e karşı), mesafe gerçekten ölçülüyor mu (relay'e karşı), ağ segmente ve izleniyor mu?**

---

<a id="2"></a>
## 2. ECU Mimarisi ve Araç-İçi Ağ Topolojisi

Bir ECU (Electronic Control Unit), bir mikrodenetleyici (MCU), bir/birden çok veriyolu alıcı-vericisi (transceiver), giriş/çıkış arabirimleri ve gömülü yazılım (firmware) içeren özelleşmiş bir gömülü bilgisayardır. İşlevine göre adlandırılır: ECM/PCM (motor/güç aktarımı), BCM (gövde kontrol modülü — kapılar, ışıklar, kornalar), ABS/ESP (fren/stabilite), EPS (elektrikli direksiyon), TCU (telematik kontrol ünitesi), IVI (infotainment), gateway (ağ geçidi).

Araç-içi ağ tek bir düz veriyolu değil, **birden çok ayrı veriyolunun bir gateway etrafında segmentlenmiş** halidir; farklı kritiklikteki sistemler ayrı veriyollarına yerleştirilir, aralarındaki trafik bir **merkezi gateway** üzerinden filtrelenerek geçer:

```
 Modern segmentli araç-içi ağ (kavramsal):

                       ┌─────────────────────┐
                       │   MERKEZİ GATEWAY    │  ← veriyolları arası
                       │ (yönlendirme+filtre) │    köprü ve filtre noktası
                       └──────────┬──────────┘
            ┌────────────┬────────┼────────────┬──────────────┐
            ▼            ▼        ▼            ▼              ▼
      [Powertrain]  [Şasi/Fren] [Gövde]   [Infotainment]  [Teşhis/OBD]
       CAN (yüksek   CAN/FlexRay  CAN/LIN    Ethernet/MOST   CAN
        hız)         (kritik)     (konfor)   (yüksek bant)
       ECM, TCM      ABS, EPS     BCM, kapı  IVI, amfi       OBD-II portu
                                  modülleri
```

Bu segmentasyon kritik bir kavramdır: amaç, düşük-güven bir alandan (infotainment, OBD portu) kaynaklanan trafiğin doğrudan yüksek-kritiklik alanına (fren, direksiyon, motor) ulaşmasını engellemektir. Gateway, hangi mesajın hangi veriyolundan hangisine geçebileceğine dair politika uygular. **Bu segmentasyonun kalitesi** modern araç güvenliğinin belkemiğidir — kötü tasarlanmış veya filtresiz bir gateway, "iç ağ kapalı ve güvenilir" varsayımının çöktüğü noktada tek savunma hattıdır ve sıklıkla yetersiz kalır.

Eski (ve bazı maliyet-odaklı modern) araçlarda bu segmentasyon zayıftır veya yoktur: tüm ECU'lar tek CAN veriyolunu paylaşabilir, OBD portu doğrudan kritik veriyoluna bağlı olabilir. Bu, CAN enjeksiyon riskini (Bölüm 14) doğrudan büyütür.

---

<a id="3"></a>
## 3. CAN Bus Derinlemesine: Çerçeve, Arbitrasyon, Bit Rate

CAN (Controller Area Network), 1980'lerde Bosch tarafından geliştirilen ve otomotivin baskın iç-ağ standardı olan bir seri veriyolu protokolüdür (ISO 11898). Anlaşılması güvenliğin temelidir, çünkü güvensizliği rastgele değil **tasarım hedeflerinin doğrudan sonucudur.** CAN, gürültülü ortamda **güvenilir, gerçek-zamanlı, çok-uçlu (multi-master)** iletişim için tasarlandı; gizlilik veya kimlik doğrulama hiçbir zaman hedef değildi.

### Fiziksel katman

Klasik yüksek-hızlı CAN, iki kablolu **diferansiyel** bir veriyoludur: CAN-High (CANH) ve CAN-Low (CANL); bilgi iki hat arasındaki gerilim farkıyla taşınır (ortak-mod gürültüye dayanıklılık), iki ucunda 120 ohm sonlandırma vardır. İki mantıksal durum: **dominant bit (mantıksal 0)** — hatlar aktif ayrılır; **recessive bit (mantıksal 1)** — hat serbest bırakılır. Kritik elektriksel kural: bir düğüm dominant (0) sürerken başkası recessive (1) sürse bile **veriyolu dominant olur** — 0, 1'i "bastırır". Bu, arbitrasyonun çalışma mekanizmasıdır.

### CAN çerçeve (frame) yapısı

CAN'de mesajlar **çerçeveler** halinde gönderilir. En önemli özellik: **çerçevelerde gönderen adresi yoktur.** Bunun yerine her çerçevenin bir **tanımlayıcısı (identifier, ID)** vardır; bu ID, mesajın *kim tarafından gönderildiğini* değil, *neyle ilgili olduğunu/önceliğini* belirtir. Klasik standart çerçevede ID 11 bit (2048 olası değer), genişletilmiş çerçevede 29 bittir.

```
 Klasik CAN veri çerçevesi (standart format, 11-bit ID) — alanlar:

 ┌─────┬──────────────┬─────┬──────┬──────────┬───────┬─────┬─────┬─────┐
 │ SOF │ Identifier   │ RTR │ IDE  │ DLC      │ DATA  │ CRC │ ACK │ EOF │
 │ 1b  │ 11 bit       │ 1b  │ r0   │ 4 bit    │ 0-8   │ 15b │ 2b  │ 7b  │
 │     │ (öncelik +   │     │      │ (veri    │ bayt  │ +1  │     │     │
 │     │  mesaj türü) │     │      │  uzunluğu)│       │ del │     │     │
 └─────┴──────────────┴─────┴──────┴──────────┴───────┴─────┴─────┴─────┘
   │         │                          │         │       │     │
   │         │                          │         │       │     └ alındı onayı
   │         │                          │         │       └ döngüsel artıklık (bütünlük)
   │         │                          │         └ asıl yük: en fazla 8 bayt
   │         │                          └ kaç bayt veri var (0-8)
   │         └ mesaj tanımlayıcı: düşük sayısal değer = yüksek öncelik
   └ çerçeve başlangıcı (start of frame)
```

Güvenlik açısından önemli alanlar: **Identifier** (öncelik ve içerik etiketi — kimlik değil), **DLC** (kaç bayt veri), **DATA** (asıl yük, klasik CAN'de en fazla **8 bayt**), **CRC** (15 bitlik döngüsel artıklık — *iletim hatalarına* karşı bütünlük, ama **kriptografik değil**; saldırgan geçerli CRC üretebilir), **ACK** (herhangi bir düğümün "hatasız aldım" onayı). Çerçevede ne gönderen kimliği, ne alıcı adresi, ne de kimlik-doğrulama alanı (MAC/imza) vardır.

### Arbitrasyon: çakışmaların öncelikle çözülmesi

CAN bir **broadcast** veriyoludur: bir çerçeveyi **bütün düğümler görür**; her düğüm ID'ye bakarak mesajın kendisini ilgilendirip ilgilendirmediğine karar verir (kabul filtresi). İki düğüm aynı anda gönderirse çakışma **arbitrasyon** ile tahrip-etmeden (non-destructive) çözülür: gönderenler ID'yi bit bit yayarken veriyolunu dinler; recessive (1) yaymak isterken dominant (0) gören düğüm kaybettiğini anlar ve susar. Düşük sayısal ID her zaman kazanır — yani **düşük ID = yüksek öncelik.** Bu, gerçek-zaman için zarif bir çözümdür (acil fren mesajı konforun önüne geçer) ama bir zafiyet kaynağıdır: sürekli düşük-ID çerçeve yayan bir düğüm veriyolunu işgal edebilir (bus monopolizasyonu — bir DoS biçimi; bkz. Bölüm 23).

```
 CAN arbitrasyon (kavram): üç düğüm aynı anda gönderiyor

   bit zamanı →   t1  t2  t3  t4  ...
   Düğüm A (ID 0x100): 0   0   1   ...   ← t3'te recessive(1) yaymak isterken
   Düğüm B (ID 0x010): 0   0   0   ...     veriyolunda dominant(0) görür → A susar
   Düğüm C (ID 0x301): 0   1  ...          ← t2'de kaybeder, susar
   Veriyolu (sonuç):   0   0   0   ...   → Düğüm B (en düşük ID) kazanır
```

### Bit rate ve zamanlama

Klasik yüksek-hızlı CAN tipik olarak **500 kbit/s** (powertrain) veya **250 kbit/s**; düşük-hızlı/konfor CAN 125 kbit/s veya altına iner. Bit rate veriyolu uzunluğuyla ters orantılıdır (yayılım gecikmesi arbitrasyonu sınırlar): 1 Mbit/s ancak ~40 m'ye kadar mümkündür. Bir CAN veriyolunu **dinlemek için adaptörün bit rate'i veriyolununkiyle aynı ayarlanmalıdır** — yanlış bit rate anlamsız çerçeve verir; bu, kendi aracını dinlerken (Bölüm 18) karşılaşılacak ilk pratik ayrıntıdır.

---

<a id="4"></a>
## 4. CAN Neden Güvensiz Tasarlandı: Broadcast ve Kimlik-Doğrulamasızlık

CAN'in güvenlik zaafları "hata" değil, 1980'lerin tehdit modelinin sonucudur: veriyolu **fiziksel olarak kapalı** kabul edildi (kötü niyetli düğüm bağlanamaz), birincil tehdit elektriksel arıza/gürültüydü, hedefler **güvenilirlik/gerçek-zaman/düşük maliyetti**; gizlilik ve kimlik doğrulama listede yoktu. Sonuçlar üç başlıkta:

**1. Broadcast — gizlilik yok.** Her çerçeveyi tüm düğümler görür; veriyoluna erişen bir düğüm **tüm trafiği pasif okuyabilir** (şifreleme yok). Fiziksel erişim sağlayan bir saldırgan (OBD portu veya kesilmiş kablo demeti üzerinden) aracın iç durumunu (hız, devir, kapı) izleyebilir.

**2. Kaynak kimlik doğrulaması yok — spoofing serbest.** Çerçevede gönderen kimliği ve MAC bulunmadığından **herhangi bir düğüm herhangi bir ID ile çerçeve yayabilir** ve alıcı bunu meşru sanır; bir aktüatör ECU'su komutu kimin gönderdiğini doğrulayamaz, yalnızca ID'ye/veriye bakar. CAN enjeksiyonunun (Bölüm 14) kavramsal temeli budur.

**3. Tazelik/replay koruması yok.** Çerçevelerde nonce/zaman damgası/sıra-sayacı standart yoktur; geçerli çerçeve kaydedilip **aynen yeniden yayınlanabilir** (replay). Bazı üst-katman protokolleri sayaç ekler, ama bu CAN tabanında yoktur.

| Zaaf | Kök neden (1980'ler) | Modern savunma yönü |
|---|---|---|
| Gizlilik yok (broadcast) | "ağ fiziksel kapalı" | Segmentasyon + fiziksel erişim kontrolü |
| Auth yok (spoof serbest) | maliyet + gerçek-zaman | Mesaj kimlik doğrulama (CAN MAC / SecOC) |
| Tazelik yok (replay) | basitlik | Sayaç/nonce + MAC, IDS davranış analizi |

Buradaki ders serinin tekrar eden temasıdır (Bölüm 6, 13): **bir protokolün güvenliği tasarım hedeflerinin içinde yoksa, sonradan sağlam eklenmesi zordur.** CAN'e güvenlik bugün üst katmanlardan (SecOC, gateway segmentasyonu, IDS) ekleniyor — ama bu, tabandan güvenli tasarlanmış bir protokolün yerini tam tutmaz ve geriye-uyumluluk yükü taşır.

---

<a id="5"></a>
## 5. CAN-FD, LIN, FlexRay, MOST ve Otomotiv Ethernet

CAN tek başına değildir; modern araçta birden çok veriyolu, farklı kritiklik/bant-genişliği ihtiyacına göre bir arada kullanılır. Güvenlik açısından önemli olan, her birinin kimlik doğrulama ve gizlilikte **tıpkı CAN gibi tasarım-gereği zayıf** olabileceğini görmektir.

| Veriyolu | Tipik hız | Kullanım | Topoloji | Kimlik doğrulama (taban) | Güvenlik notu |
|---|---|---|---|---|---|
| Klasik CAN | 125–500 kbit/s | Genel kontrol | Çok-uçlu bus | Yok | Broadcast, spoof/replay'e açık |
| CAN-FD | ~2–8 Mbit/s veri | Daha çok veri gereken kontrol | Çok-uçlu bus | Yok (taban) | 64 bayta kadar yük; SecOC eklenebilir |
| LIN | ~20 kbit/s | Düşük maliyet konfor (ayna, cam) | Tek master + slave'ler | Yok | Çok ucuz, yavaş; spoof'a açık |
| FlexRay | ~10 Mbit/s | Zaman-kritik (şasi, x-by-wire) | Bus/yıldız, zaman bölümlü | Yok (taban) | Deterministik zamanlama; kripto değil |
| MOST | 25–150 Mbit/s | Multimedya (eski) | Halka (ring) | Yok | Eğlence/ses; yerini Ethernet alıyor |
| Otomotiv Ethernet | 100 Mbit/s–10 Gbit/s | ADAS, kamera, IVI, omurga | Noktadan noktaya/switch | Var olabilir (MACsec, TLS) | Tek-çift kablo (100BASE-T1); modern, kriptoya uygun |

**CAN-FD** (Flexible Data-rate), CAN'in genişletilmiş halidir: veri alanı 8 bayttan **64 bayta** çıkar, veri fazında hız yükselir. Önemli güvenlik yan etkisi **MAC için yer açmasıdır** — 8 baytlık klasik CAN'de hem veri hem MAC sığdırmak zorken, geniş yük SecOC (Secure Onboard Communication) gibi şemaları pratik kılar.

**LIN** (Local Interconnect Network), düşük-hızlı konfor işlevleri (aynalar, koltuk motorları) için tek-master, çok ucuz bir alt-ağdır; güvenliği esasen yoktur ama kritik fonksiyon taşımaz, bir CAN ECU'sunun arkasında durur.

**FlexRay**, zaman-tetiklemeli, yüksek-hızlı, **deterministik** bir veriyoludur (fren-by-wire, aktif süspansiyon); determinizm güvenilirlik içindir, kripto için değil.

**Otomotiv Ethernet** (tek-çift-kablolu 100BASE-T1/1000BASE-T1), ADAS kameraları ve araç-içi omurga için benimseniyor. Önemi: **olgun IT güvenlik mekanizmalarını** (VLAN, MACsec katman-2 şifreleme, IPsec/TLS, firewall) otomotive taşır; geleceğin mimarileri kritik trafiği şifreli/kimlik-doğrulamalı Ethernet omurgaya taşıma eğilimindedir — ama bu da IT saldırı yüzeyini (IP yığını zafiyetleri) araca getirir.

---

<a id="6"></a>
## 6. OBD-II Portu ve CAN Erişimi (Kendi Araç / Yetkili)

OBD-II (On-Board Diagnostics, ikinci nesil), 1996 sonrası (ABD; AB'de EOBD ~2001 kademeli) araçlarda **zorunlu** standart bir teşhis arabirimidir: 16 pinli bir konektör (genellikle direksiyon altında) emisyon/arıza teşhisi için ECU'lara erişim sağlar ve modern araçların çoğunda fiziksel olarak bir/birden çok CAN veriyoluna (ISO 15765 / "CAN üzerinde diagnostik") bağlıdır.

Güvenlik açısından önemi büyüktür: bu port **iç ağa bilerek açılmış bir kapıdır.** Standart bir adaptörle (ELM327 tabanlı ucuz Bluetooth/USB dongle veya profesyonel cihaz) araç-içi CAN trafiğine erişilir. Meşru bakım/teşhis için tasarlanmıştır; ama aynı port fiziksel erişimi olan biri için trafiği okuma ve (kötü segmente araçlarda) enjeksiyon noktasıdır. Bu yüzden OBD'ye takılı kalan üçüncü-parti dongle'lar (sigorta/filo takip) başlı başına bir saldırı yüzeyidir: dongle'ın kablosuz arabirimi (Bluetooth/hücresel) zayıfsa dışarıdan iç ağa köprü olabilir — geçmişte araştırmalara konu olmuştur, **dongle güvenliği kaynaktan teyit edilmelidir.**

OBD-II'de iki ana iletişim türü vardır: (1) **ham CAN trafiği** — ECU'ların normal yaydığı periyodik çerçeveler (hız, devir) pasif gözlemlenebilir; (2) **diagnostik protokoller** — standart OBD-II PID'leri (emisyon verileri) ve üreticiye özel UDS (Unified Diagnostic Services, ISO 14229). UDS; okuma/yazma, rutin çalıştırma ve ECU yeniden programlama gibi güçlü işlevler sunar; kritik olanlar **"security access"** (seed-key) ile korunur — ama gücü üreticiye göre değişir, zayıf seed-key şemaları geçmişte eleştirilmiştir.

```
 OBD-II → CAN erişim zinciri (kendi araç / yetkili):

   [Kendi aracın]──16-pin OBD-II──[Adaptör: ELM327/SocketCAN uyumlu]──USB/BT──[Bilgisayar]
        │                                                                          │
        │  pin 6 = CANH, pin 14 = CANL (yüksek hızlı CAN, tipik)                    ▼
        │                                                            can-utils / Wireshark
        └ pasif: periyodik çerçeveleri gözle (hız/devir/kapı)        "hangi ID neyle ilişkili?"
          aktif (yetkili): standart OBD PID sorgusu                  (sadece kendi aracında öğren)
```

**Yasal ve etik sınır (vurgulu):** Kendi aracında OBD-II'den CAN trafiğini okumak ve hangi ID'nin neyle ilişkili olduğunu öğrenmek meşrudur ve bölümün alıştırmasıdır (Bölüm 18). Başkasının aracının OBD portuna izinsiz erişmek, trafiğini okumak veya komut yazmak suçtur. Bu bölüm, hiçbir aracın bir fonksiyonunu (kilit, çalıştırma) tetikleyecek **somut CAN mesajını/ID'sini** vermez — bu araç-özgü tersine mühendislikle bulunan, kötüye kullanıma açık bilgidir. Öğretilen **yöntem ve savunmadır**, hedef-özgü saldırı verisi değil.

---

<a id="7"></a>
## 7. SocketCAN ve can-utils: Pasif Okuma ve Analiz Kavramı

Linux çekirdeği, CAN veriyolunu bir ağ arabirimi gibi ele alan **SocketCAN** altyapısını içerir: CAN arabirimleri (`can0`, `vcan0`) standart soket API'siyle açılır. Üzerine kurulu **can-utils** açık kaynak araç seti, otomotiv güvenlik araştırmasının temel araçlarıdır. Amaç **kavramsal**: ne işe yaradıklarını ve hangi öğrenme akışını sağladıklarını anlamak — bir aracı manipüle etme reçetesi değil.

| Araç (can-utils) | İşlevi (kavramsal) | Tipik öğrenme kullanımı (kendi araç) |
|---|---|---|
| `candump` | Veriyolundaki çerçeveleri **pasif olarak** ekrana döker | Aracın hangi ID'leri ne sıklıkla yaydığını gözlemleme |
| `cansniffer` | Değişen baytları vurgulayarak canlı izler | Bir işlemi yaparken (kapı açma) hangi baytın değiştiğini görme |
| `cansend` | Tek bir çerçeve gönderir | **Yalnızca kendi test düzeneğinde/yetkili**; canlı araçta tehlikeli |
| `canplayer` | Kaydedilmiş trafiği yeniden oynatır | Lab analizinde kayıttan çalışma |
| `cangen` | Trafik üretir (test) | Sanal veriyolunda (vcan) öğrenme |

**Pasif okuma (`candump`)** risksiz ve öğreticidir: aracın normal yaydığı çerçeveler gözlemlenir. **Aktif gönderme (`cansend`/`canplayer`)** ise **tamamen farklı bir kategoridir:** çalışan bir araçta keyfi CAN çerçevesi göndermek motor/fren/direksiyon davranışını etkileyebilir ve **hayati tehlike** doğurabilir; ayrıca başkasının aracında suçtur. Bu yüzden öğrenme ya **sanal CAN (vcan)**, ya **hareketsiz/güvenli lab düzeneğinde** (banktaki ECU), ya da yalnızca **pasif gözlemle** yapılmalıdır. Bu metin canlı araçta aktif gönderme için hiçbir adım vermez.

```
 Güvenli öğrenme katmanları (en güvenliden riskliye):

   1. vcan (sanal CAN)        → tamamen yazılımsal; hiç araç yok    ← başla burada
   2. Kendi araç, PASİF       → sadece candump; hiçbir şey gönderme
   3. Bank/lab ECU düzeneği   → izole, hareketsiz, yetkili
   4. Canlı araç, AKTİF       → TEHLİKELİ + çoğu durumda yasadışı   ← yapma
```

Bir CAN veriyolunu "anlamak", bir işlemi (kendi aracının camını indirmek) yaparken `candump`/`cansniffer` ile hangi çerçevelerin değiştiğini gözlemleyip korelasyon kurmaktır; bu, otomotiv tersine mühendisliğinin temelidir ve **kendi aracında pasif olarak** meşrudur. Bulunan korelasyonları başkasının aracında kullanmak sınırın dışıdır.

---

<a id="8"></a>
## 8. CAN Savunması: IDS, Mesaj Kimlik Doğrulama, Segmentasyon, Gateway

CAN'in yapısal zaaflarını (Bölüm 4) tabandan değiştirmek geriye-uyumluluk nedeniyle zordur; savunma **katmanlı** olarak, çoğunlukla üst katmanlardan ve mimari düzeyden gelir. Dört ana sütun vardır.

**1. Ağ segmentasyonu ve güvenli gateway.** En temel ve etkili mimari savunma: kritik veriyollarını (powertrain, fren, direksiyon) düşük-güven veriyollarından (infotainment, OBD, telematik) ayırmak ve aralarındaki trafiği **filtreleyen, politika uygulayan bir gateway** üzerinden geçirmek ("infotainment'tan fren veriyoluna komut geçemez", "OBD'den kritik UDS yalnızca kimlik doğrulamayla erişilir"). Bu, dış kapıdan (telematik/OBD) giren bir saldırganın kritik sistemlere ulaşmasını engelleyen tek hat olabilir.

**2. Mesaj kimlik doğrulama (CAN MAC / SecOC).** CAN'in "kim gönderdi?" cevabı olmadığından, modern yaklaşım her kritik çerçeveye bir **MAC** ve bir **tazelik değeri (sayaç)** eklemektir. AUTOSAR'ın **SecOC (Secure Onboard Communication)** spesifikasyonu tam bunu tanımlar: paylaşılan simetrik anahtarla hesaplanan kısaltılmış MAC + sayaç çerçeveye eklenir, alıcı doğrular; spoof'u (auth) ve replay'i (sayaç) aynı anda hedefler. Klasik 8-baytlık CAN'de yer dar olduğundan kısaltılmış MAC kullanılır; **CAN-FD'nin geniş yükü** bu şemayı pratikleştirir.

**3. Saldırı tespit sistemi (CAN IDS).** CAN trafiği çok düzenlidir: her ID belirli periyot ve değer aralıklarında yayılır. IDS bu **normal davranış profilini** öğrenir ve sapmaları işaretler: bir ID'nin beklenenden çok sık görünmesi (enjeksiyon işareti), iki ECU'nun aynı ID'yi yayması, fiziksel-katman parmak izi (her ECU'nun gerilim/zamanlama imzası — voltaj tabanlı transmitter parmak izi) tutarsızlığı. IDS, kimlik doğrulamanın olmadığı yerde **anomaliyi** yakalar.

| Savunma katmanı | Neyi hedefler | Mekanizma | Sınırı |
|---|---|---|---|
| Segmentasyon + gateway | Yanal hareket (lateral) | Veriyolları arası filtreleme/politika | Gateway zayıfsa delinir |
| Mesaj auth (SecOC) | Spoof + replay | Çerçeveye MAC + sayaç | Yer (klasik CAN), anahtar yönetimi |
| CAN IDS | Enjeksiyon/anomali | Davranış/parmak-izi profili | Yanlış pozitif, sofistike enjeksiyon |
| Fiziksel erişim kontrolü | OBD/kablo demeti erişimi | Mekanik koruma, port denetimi | Belirlenmiş saldırgana sınırlı |

**4. Fiziksel erişim kontrolü ve güvenli teşhis.** OBD portunun ve kablo demetlerinin fiziksel korunması, UDS "security access" şemalarının güçlü (zayıf seed-key değil) olması ve teşhis oturumlarının kimlik doğrulamalı olması temel hijyendir.

> Savunma sezgisi: CAN güvenliği tek bir sihirli çözümle değil, **derinlemesine savunma** ile gelir: SecOC "kim gönderdi?" cevabı verir, segmentasyon "nereden nereye geçebilir?" sınırlar, IDS "bu davranış normal mi?" izler, fiziksel kontrol kapıyı zorlaştırır — hiçbiri tek başına yetmez. Bu, Bölüm 13'teki katmanlı savunma doktrininin (auth + nonce + segment + izleme) otomotiv iç-ağ karşılığıdır.

---

<a id="9"></a>
## 9. Anahtarsız Sistemler I — RKE: Buton, UHF, Rolling Code

Anahtarsız sistemlerde iki ayrı işlevi kesinlikle ayırmak gerekir; farklı fizik, frekans ve tehdit modeli taşırlar: **RKE** (bu başlık) ve **PKES** (Bölüm 10).

**RKE (Remote Keyless Entry)**, anahtarlıktaki bir **butona basıldığında** araca tek-yönlü bir RF komutu (kilitle/aç/bagaj) gönderir — geleneksel "uzaktan kumandalı anahtar". Tipik olarak **UHF bandında** çalışır: bölgeye göre **315 MHz** (Kuzey Amerika, Japonya), **433,92 MHz** veya **868 MHz** (Avrupa); modülasyon genellikle OOK (ASK) veya FSK'dır (Bölüm 1, 16). İletim tek yönlüdür; araç yanıt vermez.

RKE güvenliği **rolling code** üzerine kuruludur. Sabit-kod sistemlerin (her basışta aynı kodu yollayan) replay'e açık olduğu erken anlaşıldığından, modern RKE her basışta **farklı** bir kod üretir:

```
 RKE rolling-code akışı (kavramsal):

   [Anahtar butonu] basılır
        │
        ▼
   sayaç++ (her basışta artar)         araç tarafı: aynı anahtarı + sayaç penceresini bilir
        │                                   │
        ▼                                   ▼
   kod = f(gizli_anahtar, sayaç, komut)   gelen kodu çözer → sayaç penceresinde mi?
        │                                   │  evet → komutu uygula, pencereyi ilerlet
        ▼                                   │  hayır → reddet
   UHF (315/433/868) ── OOK/FSK ──►─────────┘
```

Mekanizma: anahtar ve araç bir **paylaşılan gizli anahtar** ve bir **senkron sayaç** paylaşır. Her basışta sayaç artar; kod gizli anahtar + sayaç + komuttan türetilir; araç gelen kodu çözüp sayacın **kabul penceresi** içinde (ve son kullanılandan ileride) olup olmadığına bakar, geçerliyse komutu uygular ve penceresini ilerletir. Eski bir kod tekrar gönderilirse pencerenin gerisinde kaldığı için reddedilir — basit replay böyle engellenir. Pencere, sürücünün araçtan uzakta birkaç kez basması (senkron kayması) durumuna tolerans sağlar.

Rolling-code'un **kriptografik kalitesi** üretici/algoritmaya göre büyük ölçüde değişir. Tarihsel olarak yaygın algoritmalar:

| Algoritma | Tip | Notlar (güvenlik) |
|---|---|---|
| KeeLoq | Blok şifre tabanlı rolling code | Çok yaygın; akademik çalışmalarda çeşitli zayıflıklar raporlandı (anahtar çıkarımı, üretici-anahtarı senaryoları). Belirli ürün/etki **kaynaktan teyit edilmeli** |
| Hitag2 | Akış şifre tabanlı (transponder + RKE) | Eski; kriptanalizi yayımlanmış, zayıf kabul edilir |
| AES tabanlı (modern) | 128-bit AES rolling code | Doğru uygulanırsa güçlü; modern araçların yöneldiği taban |

> Savunma sezgisi (RKE): Rolling-code **tazelik** (her basış farklı kod → replay'e direnç) sağlar; bu gereklidir ama iki şeyi sağlamaz: (1) zayıf algoritma (eski Hitag2, bazı KeeLoq kullanımları) kriptografik kırılabilir → savunma: AES tabanlı güncel sistem; (2) "yakala-engelle-tekrar" (RollJam, Bölüm 12) gibi senkron-penceresini kötüye kullanan saldırılara doğal koruma yoktur → savunma uygulama detayı ve zaman pencerelerindedir. Ders: koruma "rolling" etiketinde değil, **kriptografik kalitede ve protokol detayındadır.**

---

<a id="10"></a>
## 10. Anahtarsız Sistemler II — PKES: LF Uyandırma, UHF Yanıt, Yakınlık

**PKES (Passive Keyless Entry and Start)** — pasif anahtarsız giriş/çalıştırma — RKE'den temelde farklıdır: **buton yoktur.** Anahtar cebinizdeyken araca yaklaşıp kapı kolunu çeker veya çalıştırma düğmesine basarsınız; araç, anahtarın **yakında** olduğunu otomatik bir sorgu-yanıt (challenge-response) ile doğrular. Kolaylık yüksektir; tehdit modeli de tam burada değişir.

PKES iki frekans bandı kullanır ve bu ikilik hem işleyişin hem relay zaafının (Bölüm 13) anahtarıdır:

- **LF (Low Frequency) uyandırma/sorgu — ~125 kHz:** Araç, kapı kollarına ve kabine gömülü LF antenlerden **çok kısa menzilli** (~1–2 m) bir alan/sorgu yayar. Bu kısa menzil **bilinçli tasarımdır**: araç, "anahtar fiziksel olarak çok yakın olmalı" varsayımını LF'in doğal menzil sınırına dayandırır. LF ayrıca anahtarı "uyandırır" (normalde pil tasarrufu için uykudadır).

- **UHF yanıt — ~315/433/868 MHz:** Uyanan anahtar, LF sorgusuna kriptografik yanıtı **UHF** bandında geri gönderir (çift yönlü diyalog burada tamamlanır).

```
 PKES normal çalışma (kavramsal):

   [ARAÇ]  ── LF sorgu (~125 kHz, ~1-2 m) ──►  [ANAHTAR (cepte)]
      │          "yakında mısın? işte challenge"        │
      │                                                 │ uyanır, kriptografik
      │  ◄── UHF yanıt (315/433/868) ◄──────────────────┘ response hesaplar
      ▼
   challenge-response geçerli + anahtar "yakın" sanılıyor → kapı açılır / çalıştırılır
```

Sistem sorgu-yanıtı kriptografik doğrular (anahtarın gizli anahtarı bilmesi gerekir — spoofing'i zorlaştırır). Ama kritik tasarım varsayımı şudur: **"LF sorgusunu duyabilen anahtar fiziksel olarak yakındır."** Sistem mesafeyi doğrudan ölçmez; LF'in kısa doğal menzilini bir **mesafe vekili (proxy)** olarak kullanır. Relay saldırısının (Bölüm 13) saldırdığı nokta tam budur: LF sorgusu ve UHF yanıtı yapay olarak uzağa taşınırsa, kriptografi geçerli ve mesafe ölçülmemiş olduğundan araç hâlâ "anahtar yakın" sanır.

PKES'in zaafı **kriptografik değil, fiziksel/mimaridir:** "yakın görünmek" (LF'i duyabilmek) ile "fiziksel yakın olmak" arasındaki farkı ölçmemesi. Bu ayrım, en yaygın araç hırsızlığı yönteminin (relay) neden işe yaradığını ve savunmanın (gerçek mesafe ölçümü, UWB) neden tam burayı hedeflediğini açıklar.

---

<a id="11"></a>
## 11. Immobilizer ve Transponder: RFID 125 kHz / 13,56 MHz

RKE ve PKES "girişle" (kapı) ilgilenir; **immobilizer** ise "çalıştırmayla" (motor) ilgilenen ayrı bir katmandır ve sıklıkla karıştırılır. İşlevi: doğru **transponder** (anahtar içindeki RFID çipi) yokken motorun çalışmasını **engellemek.** Çoğu ülkede 1990'ların sonundan zorunlu olması, "anahtarı zorla/kopyala, motoru çalıştır" tipi hırsızlığı büyük ölçüde azaltmış ve hırsızlığı kablosuz/elektronik düzleme kaydırmıştır.

Mekanizma: anahtarda pasif bir **transponder** (pilsiz, okuyucunun alanından beslenen RFID çipi) bulunur; marş/kontak bölgesinde bir **anten halkası** vardır. Anahtar kontağa sokulunca (veya PKES'te düğmeye basılınca) okuyucu transponder'ı besler ve sorgu-yanıt yürütür; doğru yanıt gelmezse motor ECU'su yakıt/ateşlemeyi serbest bırakmaz.

| Boyut | Tipik değer | Not |
|---|---|---|
| Frekans (eski/yaygın) | 125 kHz (LF RFID) | Hitag, Megamos vb. transponder aileleri |
| Frekans (yeni/bazı) | 13,56 MHz (HF RFID) | Bazı modern transponder uygulamaları |
| Çip örnekleri | Hitag2, Megamos Crypto, DST, AES tabanlı modern | Eski olanların bir kısmının kriptanalizi yayımlandı |
| Bağımsızlık | RKE/PKES'ten **ayrı** katman | Giriş açılsa bile motor çalışmaz |

Tarihsel ders nettir ve Bölüm 16'daki MIFARE/Crypto-1 dersiyle birebir örtüşür: bazı eski transponder algoritmaları (Hitag2, Megamos Crypto) **tescilli/gizli (security through obscurity)** tasarlandı ve açığa çıkıp kriptanalize uğradığında zayıf oldukları görüldü. Belirli araç modellerindeki etki ve mevcut durum **kaynaktan teyit edilmelidir**; genel ders, güvenliğin gizlilikten değil **açık ve sağlam kriptografiden** (modern AES tabanlı transponder) gelmesidir.

> Savunma sezgisi (immobilizer): Immobilizer, girişten bağımsız ikinci bir kilittir ve mimari değerlidir (giriş aşılsa bile çalıştırmayı engeller). Ama gücü tamamen **transponder kriptosunun kalitesine** bağlıdır: eski/kırılmış çipler (Hitag2/Megamos eski uygulamaları) klonlanabilir, modern AES tabanlı sistemler çok daha dirençlidir. Ders: immobilizer'ın "ayrı katman" değeri gerçektir, ama "kripto kalitesi belirleyicidir".

---

<a id="12"></a>
## 12. Saldırı Sınıfları I — Sabit-Kod Replay ve Rolling-Code Zayıflıkları

Bu ve sonraki iki başlık, saldırı sınıflarını **yalnızca kavramsal köken, "neden mümkün" ve savunma** ekseninde ele alır. Hiçbiri operasyonel reçete değildir; bilerek saldırı verisi, frekans-özgü zamanlama veya araç-özgü mesaj içermez. Amaç, savunmanın neye karşı tasarlandığını görmek.

### Sabit-kod replay (çoğunlukla eski araçlar)

En eski ve basit sınıf. Sabit-kod RKE (her basışta aynı kodu yayan eski sistemler) veya sabit-kod garaj kumandaları gönderdikleri kodu **hiç değiştirmez**; bir gözlemci kodu bir kez yakalarsa aynen yeniden yayarak (replay) komutu tekrarlayabilir. **Neden mümkün:** tazelik yok — "eski" ile "yeni" ayırt edilemez (Bölüm 4'teki CAN replay ile aynı kök). **Kimde:** çok eski araçlar, basit garaj kapıları, ucuz alarm kumandaları. **Savunma:** rolling-code veya kriptografik sorgu-yanıt; sabit-kod cihazı değiştirmek. Modern araçların ezici çoğunluğu bunu rolling-code ile çözmüştür.

### Rolling-code zayıflıkları — kavramsal

Rolling-code basit replay'i çözer ama "her rolling-code güvenlidir" demek yanlıştır; literatürde iki kavramsal sınıf tartışılır:

**RollJam (yakala-engelle-tekrar) kavramı.** Rolling-code'un *senkron penceresini* kötüye kullanma fikri: saldırgan kullanıcının ürettiği geçerli kodu **yakalarken aynı anda hafifçe karıştırarak (jamming) aracın o kodu duymasını engeller**; kullanıcı tekrar bastığında ikinci kodu da yakalar, bu kez ilk (henüz "harcanmamış") kodu araca iletip ikinciyi saklı tutar — elinde aracın henüz görmediği geçerli bir kod kalır. **Neden mümkün:** rolling-code "yakalanmış ama araca ulaşmamış" bir kodun geçerliliğini korur ve jamming araç-tarafı alımı engelleyebilir. **Belirli koşullara/uygulamaya bağlıdır; uygulanabilirliği ve modern yamaların etkisi kaynaktan teyit edilmelidir.** **Savunma:** kodlara kısa geçerlilik penceresi (eski yakalanmış kod zamanla geçersizleşir), çift-yönlü doğrulama, jamming tespiti, mümkünse mesafe-doğrulamalı (UWB) sistem.

**Rolling-PWN / sayaç-resenkronizasyon kavramı.** Belirli üreticilere yönelik raporlanan bir sınıf, rolling-code **sayaç senkronizasyonunun** kötüye kullanılabildiği fikrine dayanır: belirli koşullarda ardışık kodların yeniden gönderimi sayaç penceresini istenmeyen biçimde resenkronize edip eski bir kodu yeniden geçerli kılabilir. **Belirli araştırma/modellere özgü bir iddiadır; kesin etki, etkilenen modeller ve yama durumu mutlaka kaynaktan teyit edilmelidir** — bu metin doğrulanmış genel gerçek olarak sunmaz. **Savunma:** sağlam sayaç-pencere yönetimi (gereksiz geniş pencere açmama, eski koda dönüşü engelleme), üretici yaması.

```
 Replay ailesinin tek-cümlelik özü:

   Sabit kod      → her seferinde aynı → yakala+tekrar (replay)         [çözüm: rolling/nonce]
   RollJam (kavr.)→ geçerli kodu yakala+jam+sakla+sonra kullan         [çözüm: zaman penceresi,
                                                                         jam tespiti, UWB]
   Sayaç resenkr. → sayaç penceresini kötüye kullanıp eski kodu diril  [çözüm: sağlam pencere
   (kavr., teyit) →                                                      yönetimi, yama]
```

> Savunma sezgisi: Rolling-code "tazelik" verir ama tek başına yeterli garanti değildir; saldırı yüzeyi **senkron-pencere yönetimine ve jamming'e** kayar — tazelik mekanizmasının kendisi de (pencere genişliği, eski kodun geçerlilik süresi, jamming tespiti) bir saldırı yüzeyidir. En güçlü yapısal çözüm, anahtarsız sistemde **mesafeyi gerçekten ölçmeye** (UWB, Bölüm 13) geçmektir — bu, hem relay'i hem "yakalanmış kodu sonra kullan" senaryolarını kökten zorlaştırır.

---

<a id="13"></a>
## 13. Saldırı Sınıfları II — Relay (Röle) Saldırısı: Prensip ve Savunma

Relay (röle / menzil uzatma) saldırısı, **günümüzde anahtarsız araçlara yönelik en yaygın hırsızlık yöntemidir** ve PKES'in (Bölüm 10) fiziksel varsayımını hedef alır. Bu başlık prensip ve savunmaya odaklanır; bilinçli olarak hiçbir kurulum adımı, donanım reçetesi veya zamanlama parametresi vermez.

### Prensip — kodu kırmaz, mesafeyi uzatır

Relay'in zarif (ve tehlikeli) yanı **hiçbir kriptografiyi kırmaya çalışmamasıdır.** PKES'te araç ile anahtar yalnızca yakınken konuşabilmelidir (LF'in kısa menzili sayesinde); relay bu konuşmayı iki nokta arasında **köprüleyerek** aracın LF sorgusunu uzaktaki anahtara, anahtarın UHF yanıtını geri araca taşır. Araç meşru kriptografik yanıtı aldığı ve mesafeyi ölçmediği için anahtarı "yakında" sanır — kod kırılmaz, yalnızca **mesafe yapay uzatılır.**

```
 PKES relay (röle) saldırısı — yalnızca prensip:

   [ARAÇ (sokakta)]                                      [ANAHTAR (evde, masada)]
        │ LF sorgu                                              ▲ LF (yapay)
        ▼                                                       │
   (Cihaz A: aracın yanında)  ⇠⇠⇠ uzun köprü ⇢⇢⇢  (Cihaz B: anahtarın yakınında)
        ▲                                                       │
        │ UHF yanıt köprüyle geri taşınır                       ▼ UHF yanıt
        └───────────────────────────────────────────────── anahtar normal cevap verir

   Araç "anahtar ~1 m'de" sanır; oysa anahtar onlarca metre/duvar arkasında.
   KRİPTO KIRILMAZ — yalnızca "yakınlık" yalanı söylenir.
```

**Neden mümkün:** PKES mesafeyi ölçmez, LF'in kısa menzilini bir vekil olarak kullanır (Bölüm 10); köprü bu vekili geçersiz kılar — anahtar uzakken bile "LF duyulabilir" hale gelir. Rolling-code/kriptografi koruma sağlamaz, çünkü iletilen yanıt **gerçek ve geçerlidir**; sorun tazelikte değil, **mesafededir.**

### Savunma — mesafeyi gerçekten ölçmek

| Önlem | Ne sağlar | Not / sınır |
|---|---|---|
| **UWB (Ultra-Wideband) mesafe ölçümü (time-of-flight)** | Sinyalin gidiş-dönüş süresinden **gerçek fiziksel mesafe** ölçer; köprünün eklediği gecikme "yakın" iddiasını fiziksel olarak çürütür | Modern güvenli keyless'in yapısal çözümü; yaygınlaşıyor |
| **Hareket sensörlü fob** | Anahtar bir süre hareketsizse (gece, masada) LF'e yanıt vermeyi/yaymayı durdurur | Relay penceresini kapatır; ucuz, etkili |
| **Faraday kesesi / kutusu (kullanıcı önlemi)** | Anahtarı saran iletken ekranlama RF'i tamamen keser → anahtar LF sorgusunu duyamaz | En ucuz, en erişilebilir kullanıcı önlemi; evde/işte anahtarı içine koymak |
| **PKES'i kapatma seçeneği** | Pasif girişi devre dışı bırakıp butona (RKE) dönmek | Bazı araçlarda menüden/fob'dan mümkün |
| **OEM güvenlik yamaları / sıkı zamanlama** | Yanıt için izin verilen gecikmeyi daraltma (köprü gecikmesini zorlaştırma) | Tek başına yetersiz; UWB kadar yapısal değil |

Kritik savunma fikri tek cümledir ve serinin omurga dersidir (Bölüm 13, 16): **relay'i yenmenin yapısal yolu, "yakın görünmek" ile "fiziksel yakın olmak" arasındaki farkı ölçmektir.** UWB time-of-flight bunu yapar: gidiş-dönüş süresi mesafeyle orantılıdır, köprü kaçınılmaz gecikme ekler, ölçülen "mesafe" gerçeğinden büyük çıkar ve "yakın" iddiası çürür.

> Mühendislik sezgisi: Relay, "tekrar edilebilen ve doğrulanmayan şey güvenli değildir" ilkesinin mesafe boyutudur. Sabit kod tekrar edilebilir (replay → nonce/rolling); doğrulanmayan komut sahtelenebilir (spoof → kripto kimlik); **ölçülmeyen mesafe yanıltılabilir (relay → gerçek mesafe ölçümü, UWB).** Anahtarsız araçlar bu üçlünün en somut ve yüksek-bahisli sahnesidir. Kullanıcı için en erişilebilir önlem: anahtarı evde Faraday kesesinde tutmak ve hareket-sensörlü/UWB'li sistem tercih etmek.

---

<a id="14"></a>
## 14. Saldırı Sınıfları III — CAN Enjeksiyon ve Immobilizer Atlatma Kavramı

Bu başlık, son yıllarda öne çıkan bir hırsızlık sınıfını ve immobilizer atlatma fikrini **yalnızca kavramsal** ele alır. Hiçbir araca uygulanacak somut mesaj, ID, konektör pini veya adım verilmez; bu bilgi araç-özgüdür ve doğrudan kötüye kullanıma açıktır.

### CAN enjeksiyon (fiziksel/yerel erişim üzerinden) — kavram

CAN enjeksiyonu, Bölüm 4'teki yapısal zaafın (auth yok; herhangi bir düğüm herhangi bir ID ile çerçeve yayabilir) doğrudan istismarıdır. Kavramsal mantık: saldırgan araç-içi CAN veriyoluna **fiziksel erişir** ve gövde/giriş ECU'larının beklediği belirli kontrol çerçevelerini **enjekte ederek** (örneğin "kapı kilidi aç" veya "geçerli anahtar var" durum mesajı) aracı kandırır. Erişim noktası genellikle dış-erişimli bir kablo demetidir — far ünitesi, tampon sensörleri gibi, gövde paneli sökülerek ulaşılan ve CAN'e bağlı bir nokta ("far-içi/headlight" erişim kavramı medyada bu sınıf için kullanılmıştır).

**Neden mümkün:** (1) CAN'de kaynak kimlik doğrulaması yok → enjekte çerçeve meşru sanılır; (2) ağ **yetersiz segmente** ise dış düğümden (far) kritik veriyoluna doğrudan ulaşılabilir; (3) bazı araçlarda giriş/çalıştırma yetkisi tek bir CAN durum-mesajına güvenir. **Savunma** doğrudan Bölüm 8'den gelir: **SecOC** enjekte çerçeveyi reddeder (MAC tutmaz); **segmentasyon + güvenli gateway** geçişi engeller; **CAN IDS** anomaliyi (beklenmeyen ID/sıklık, çift gönderen) yakalar; **fiziksel erişim zorlaştırma** kapıyı daraltır. Bu, "iç ağ kapalı ve güvenilir" varsayımının neden artık geçerli olmadığının en somut kanıtıdır.

```
 CAN enjeksiyon — kavramsal zincir ve savunma kesişimleri:

   dış-erişimli kablo demeti (far/tampon)  ── fiziksel erişim ──┐
                                                                ▼
        [yetersiz segmentasyon?] ──hayır filtre──►  kritik CAN veriyolu
                                                                │
        enjekte çerçeve (sahte durum/komut)  ──auth yok──►  ECU "meşru" sanır
                                                                │
   SAVUNMA: ① SecOC (MAC tutmaz→ret)  ② gateway segment (geçişi kes)
            ③ IDS (anomali yakala)    ④ fiziksel koruma (erişimi zorlaştır)
```

### Immobilizer atlatma — kavram

Immobilizer atlatma fikri (Bölüm 11) iki kavramsal yola dayanır: (1) **zayıf transponder kriptosu** — eski/kırılmış çiplerde (Hitag2/Megamos eski uygulamaları) transponder yanıtının taklit/klonlanabilmesi; (2) **mimari atlatma** — immobilizer kontrolünün CAN üzerinden taşınan bir durum-mesajına güvendiği ve bu mesajın enjeksiyonla taklit edilebildiği durumlar (CAN enjeksiyonuyla kesişir). **Neden mümkün:** gizliliğe dayalı zayıf kripto ve/veya kimlik-doğrulamasız iç-ağ. **Savunma:** modern AES tabanlı transponder; immobilizer kararının kimlik-doğrulamalı (SecOC) ve segmente bir kanaldan gelmesi; sırra değil sağlam algoritmaya dayanma. Belirli araç/çip etkileri **kaynaktan teyit edilmelidir.**

> Yasal sınır (vurgulu, tekrar): Bu başlık CAN enjeksiyon ve immobilizer atlatma için **hiçbir somut mesaj, ID, pin, donanım veya adım vermez.** Verilen tek şey, savunmanın neye karşı tasarlandığını anlamak için **kavramsal köken + neden mümkün + savunma**dır. Başkasının aracında bunları denemek hırsızlık ve bilişim suçudur (TCK 141–142, 243–245). Araştırmacı bu bilgiyi yalnızca **savunma tasarımı** ve **kendi aracının/yetkili düzeneğin yasal analizi** için kullanır.

---

<a id="15"></a>
## 15. TPMS: Lastik Basıncı İzleme ve Araç-İzleme Gizliliği

TPMS (Tire Pressure Monitoring System — lastik basıncı izleme), birçok ülkede zorunludur ve her lastikteki sensörün basınç/sıcaklık verisini araca kablosuz iletir. **Doğrudan TPMS bir kapı/çalıştırma açığından çok bir gizlilik/izleme meselesidir** ve bu yönüyle Bölüm 16'daki BLE adres-randomizasyonu ve Bölüm 11'deki meta-veri sızıntısı temasının otomotiv karşılığıdır.

Teknik: çoğu doğrudan TPMS sensörü **315 MHz** (Kuzey Amerika) veya **433 MHz** (Avrupa) bandında periyodik (veya hareket/basınç değişiminde) bir paket yayar; pakette tipik olarak **benzersiz ve çoğunlukla sabit bir sensör kimliği (ID)** ile basınç/sıcaklık verisi bulunur. Şifreleme genellikle yok veya zayıftır (sensör pilini korumak için minimal tasarım). Gizlilik riski iki yönlüdür:

**Araç izleme/tanıma:** Sabit TPMS ID'leri, bir aracı (dolayısıyla sahibini) yol kenarındaki alıcılarla **tanıma ve hareketini izleme** imkânı verebilir — plaka okumadan, pasif RF dinlemeyle; kitlesel takip altyapısı açısından kavramsal bir endişedir.

**Sahte veri kavramı:** TPMS kimlik-doğrulamasız olduğundan sahte basınç paketleri yayınlayıp panelde yanlış uyarı tetiklemek kavramsal olarak mümkündür (güvenlik/rahatsızlık sorunu; fren sistemine doğrudan etki tipik olarak yoktur ama **araç-özgü davranış teyit edilmeli**).

| Boyut | Durum (tipik) | Savunma / not |
|---|---|---|
| Frekans | 315 MHz (KA) / 433 MHz (AB) | rtl_433 ile kendi sensörün gözlemlenebilir (Bölüm 18) |
| Kimlik | Sabit, benzersiz sensör ID | İzleme riskinin kaynağı |
| Şifreleme | Genellikle yok/zayıf | Sahte-veri ve okuma mümkün |
| Kullanıcı savunması | Sınırlı (üretici tasarımına bağlı) | Esas çözüm üretici: dönüşümlü/şifreli ID |
| Üretici savunması | Dönüşümlü ID, kimlik doğrulama | Henüz yaygın değil; teyit edilmeli |

> Savunma sezgisi (TPMS): TPMS, "araç bile sürekli benzersiz kimlik yayar" gerçeğinin somut örneğidir — gizlilik yalnızca telefon/BLE katmanında değil, aracın her bileşeninde düşünülmelidir. Kullanıcı tarafında doğrudan önlem sınırlıdır; asıl çözüm üretici tarafında **dönüşümlü/şifreli sensör kimliği ve kimlik doğrulamadır** — henüz yaygın değildir, teyit edilmelidir. Ders: meta-veri (sabit ID), içerik şifreli olmasa bile başlı başına bir izleme kanalıdır (Bölüm 11).

---

<a id="16"></a>
## 16. Modern Saldırı Yüzeyi: Telematik, EV Şarj, V2X

Araç giderek daha "bağlı" hale geldikçe, saldırı yüzeyi anahtar RF'i ve OBD portunun çok ötesine, **uzaktan erişilebilir** yeni kanallara genişler. Bu başlık üç önemli alanı kavramsal olarak tanıtır; ayrıntıların bir kısmı hızla gelişmektedir ve **kaynaktan teyit edilmelidir.**

### Telematik ve uzaktan erişim (en kritik uzak yüzey)

Modern araçlar bir **telematik kontrol ünitesi (TCU)** içerir: gömülü hücresel modem (giderek **eSIM**), bazen WiFi, üreticinin bulutu ve mobil uygulamasıyla konuşan bir bağlantı. İşlevleri: acil çağrı (eCall), uzaktan kilit/aç, uzaktan çalıştırma/iklimlendirme, konum, OTA. Güvenlik açısından **en kritik uzak yüzey budur**, çünkü:

- **Fiziksel temas gerektirmez:** hücresel kapsamadaki bir saldırgan (veya buluttaki bir zafiyet) potansiyel olarak dünyanın öbür ucundan erişebilir — relay/RF saldırılarının menzil sınırını tamamen aşar.
- **İç ağa köprü olabilir:** TCU kötü segmente ise, uzaktan ele geçirilen bir telematik birimi CAN'e komut köprüleyebilir. Geçmişte yüksek-profilli araştırmalar, zayıf segmentasyon + telematik zafiyeti birleşiminin kritik fonksiyonlara uzaktan erişime yol açabildiğini göstermiştir (detaylar **kaynaktan teyit edilmeli**).
- **Bulut/uygulama yüzeyi:** Araç sağlam olsa bile üreticinin bulut API'si/uygulaması zayıfsa (kimlik doğrulama/yetkilendirme hataları), saldırgan hesap üzerinden uzaktan kilit/konum/çalıştırma elde edebilir — klasik web/API güvenliği problemi otomotive taşınmıştır.

**Savunma:** TCU'nun kritik veriyollarından **sıkı segmentasyonu**, güvenli önyükleme ve imzalı/doğrulanmış OTA, bulut/API tarafında güçlü kimlik doğrulama/yetkilendirme, TCU arabirimlerinin (hücresel/WiFi/Bluetooth) sağlamlaştırılması. Modern tehdit modeli aracı bir **uzaktan-erişilebilir gömülü sistem** olarak ele almak zorundadır.

### EV şarj iletişimi (ISO 15118)

Elektrikli araçlar şarj istasyonuyla yalnızca güç değil **veri** de alışverişi yapar (şarj pazarlığı, kimlik, ödeme). **ISO 15118** bu iletişimi (özellikle "Plug & Charge" — kablo takılınca otomatik kimlik/ödeme) tanımlar ve **PLC (Power Line Communication)** üzerinden taşır. Saldırı yüzeyi kavramsaldır: şarj iletişiminin dinlenmesi/araya girilmesi kimlik/ödeme verisi açısından risk taşıyabilir; ayrıca şarj noktası altyapısı (yönetim ağları) bir hedeftir. ISO 15118 bu yüzden **TLS tabanlı güvenlik ve PKI** öngörür; ayrıntılar/olgunluk gelişmektedir, **kaynaktan teyit edilmelidir.** **Savunma:** kripto/PKI özelliklerinin doğru uygulanması, şarj altyapısının ayrı güvenlik yönetimi.

### V2X (DSRC / C-V2X) — kavram

V2X (Vehicle-to-Everything), aracın diğer araçlarla (V2V), altyapıyla (V2I), yayalarla (V2P) iletişimini kapsar (çarpışma uyarısı, trafik). İki teknoloji yarışmıştır: **DSRC** (802.11p) ve **C-V2X** (hücresel). Kavramsal endişeler: sahte mesaj enjeksiyonu (var olmayan tehlike/araç bildirme), gizlilik (araç hareketi izleme), mesaj bütünlüğü. Bu yüzden V2X, **imzalı mesajlar ve büyük ölçekli bir PKI** (dönüşümlü takma-ad sertifikalarıyla gizlilik) üzerine tasarlanır. Alan gelişmektedir, dağıtımı bölgeseldir; ayrıntılar **kaynaktan teyit edilmelidir.**

```
 Bağlı-araç uzak saldırı yüzeyi — menzil ekseni:

   yakın ────────────────────────────────────────────────────► uzak
   RKE/PKES   TPMS     Bluetooth/   V2X        EV şarj      Telematik/
   relay      izleme   WiFi (IVI)   (yerel)    (istasyon)   bulut/uygulama
   (~m-10sm)  (~10m)   (~10-100m)   (~yüzlerce  (kablolu)    (sınırsız:
                                     m)                       hücresel/internet)
   └─ fiziksel/RF yakınlık gerekir ─┘            └─ fiziksel temas GEREKMEZ ─┘
                                                  (en yüksek ölçeklenme riski)
```

> Mühendislik sezgisi: Otomotiv saldırı yüzeyinin ağırlık merkezi, RF-yakınlık (anahtar/TPMS) gerektiren saldırılardan **fiziksel temas gerektirmeyen uzaktan-bağlanabilirlik** (telematik/bulut) saldırılarına kayıyor. Neden ölçeklenebilirlik: relay tek aracı, tek seferde, fiziksel yakınlıkla hedefler; bir bulut/telematik zafiyeti potansiyel olarak **bütün bir filoyu uzaktan** hedefleyebilir. Bu yüzden modern otomotiv güvenliği aracı bir RF cihazı değil, **bağlı gömülü sistem + bulut servisi** olarak savunmak zorundadır (Bölüm 24 bu kaymayı yansıtır).

---

<a id="17"></a>
## 17. Savunma Derinliği: Katmanlı Mimari, ISO/SAE 21434, Sorumlu Açıklama

Bütün önceki başlıkların savunmaları tek bir doktrinde birleşir: **derinlemesine savunma (defense in depth).** Hiçbir tek mekanizma yeterli değildir; güvenlik, üst üste binen katmanlardan ve bir **süreçten** doğar.

### Katmanlı savunma — otomotiv özeti

```
 Otomotiv derinlemesine savunma katmanları (dıştan içe):

   ① Dış kapı sertleştirme   : RKE AES + UWB mesafe; TPMS dönüşümlü ID;
                               telematik/bulut auth; şarj/V2X PKI
   ② Segmentasyon + gateway  : kritik veriyolları ayrı; dış düğümden geçişi filtrele
   ③ Mesaj kimlik doğrulama  : SecOC (MAC + sayaç) → spoof + replay'i kes
   ④ İzleme (IDS)            : CAN anomali + parmak-izi → enjeksiyonu yakala
   ⑤ Güvenli önyükleme + OTA : imzalı firmware; yamayı sahaya güvenli ulaştır
   ⑥ Fiziksel erişim kontrolü: OBD/kablo demeti koruması; güvenli teşhis
   ⑦ Süreç (21434)           : tehdit analizi, güvenli geliştirme, izleme, yanıt
```

Her katman farklı bir saldırı sınıfını hedefler ve biri delinse bile diğeri devrede kalır: relay'i ① (UWB), CAN enjeksiyonu ②③④, uzaktan ele geçirmeyi ①②⑤, fiziksel saldırıyı ⑥ ele alır. Önceki bölümlerde tek tek gördüğümüz savunmalar, aslında bu tablonun satırlarıdır.

### ISO/SAE 21434 — güvenliği sürece bağlamak

Otomotiv siber güvenliği artık yalnızca teknik değil, **düzenleyici ve süreçseldir.** **ISO/SAE 21434** (Road vehicles — Cybersecurity engineering), aracın tüm yaşam döngüsü boyunca (kavram → geliştirme → üretim → işletim → bakım → hizmetten çıkarma) siber güvenlik mühendisliğini tanımlar: tehdit analizi/risk değerlendirmesi (TARA), güvenli geliştirme, tedarik zinciri güvenliği, sahada **izleme/olaya yanıt.** Eşlik eden **UNECE WP.29 (R155/R156)**, üreticilerden bir siber güvenlik yönetim sistemi (CSMS) ve yazılım güncelleme yönetimi (SUMS) ister. Anlamı: güvenlik "bir kez bitirilen" değil, **sürekli yönetilen** bir disiplindir. (Kapsam/yürürlük ayrıntıları **kaynaktan teyit edilmelidir.**)

### Sorumlu açıklama (responsible disclosure)

Otomotiv güvenlik araştırmasının etik omurgası **sorumlu açıklamadır:** bir araştırmacı zafiyet bulduğunda doğru yol, bulguyu önce **üreticiye/CERT'e özel bildirmek**, makul düzeltme süresi tanımak, etkilenenleri riske atmadan koordineli açıklamaktır. Yapılmaması gereken: tam exploit/araç-özgü saldırı verisini kamuya saçmak veya kötüye kullanmak. Seri boyunca (Bölüm 6, 13, 24) işlenen ilke geçerlidir: **bulgunun amacı savunmayı güçlendirmektir, zarar vermek değil.** Bu bölümün baştan sona "prensip + savunma, reçete yok" duruşu bu çerçevenin yansımasıdır.

> Kapanış sezgisi: Otomotiv güvenliği "aracı kırılamaz yapmak" değil — bu imkânsızdır — **saldırıyı yeterince zorlaştırmak, tespit etmek ve yanıt verebilmektir.** Her katman saldırganın işini büyütür; süreç (21434) bunu sürdürülebilir kılar; sorumlu açıklama ekosistemi güçlendirir. Araştırmacının görevi, bu katmanların nerede ince olduğunu **savunmayı iyileştirmek için** görmektir.

---

<a id="18"></a>
## 18. Alıştırmalar (Yalnızca Kendi Araç / Yasal)

Aşağıdaki alıştırmalar tasarım gereği **yalnızca kendi aracında, kendi anahtarınla ve pasif/güvenli biçimde** yapılır; hiçbiri başkasının aracına, kilit açmaya veya çalıştırmaya yönelik değildir, hepsi **gözlem ve anlama** amaçlıdır. Canlı araçta CAN'e çerçeve **göndermek** (pasif okumak değil) hayati tehlike ve çoğu durumda yasal sorun doğurur — bu alıştırmalar bunu içermez.

### Alıştırma 1 — Kendi aracının CAN trafiğini pasif okumak (candump)

**Amaç:** CAN'in broadcast yapısını canlı görmek. **Düzenek:** kendi aracın + SocketCAN uyumlu OBD-II adaptörü (ELM327; bazı analizlerde daha yetenekli arabirim gerekebilir, teyit et). **Adımlar:** arabirimi doğru bit rate ile aç (yanlış hız anlamsız veri verir), `candump` ile çerçeveleri pasif izle; kontak açıkken (**araç güvenli/hareketsiz**) hangi ID'lerin ne sıklıkla aktığını gözle. **Gözlem:** "kaynak adresi yok, ID neyle ilgili olduğunu söylüyor" gerçeğini ve trafiğin yoğunluğunu/periyodikliğini somut gör. **Sınır:** yalnızca `candump` (pasif); hiçbir şey gönderme.

### Alıştırma 2 — Bir işlemi CAN'de korele etmek (cansniffer, pasif)

**Amaç:** "hangi çerçeve neyle ilişkili" sorusunu deneyimlemek. **Adımlar:** `cansniffer` ile (değişen baytları vurgular) izlerken, **araç güvenli ve hareketsizken** zararsız bir işlem yap (kapını kilitle/aç, cam aç) ve hangi baytların değiştiğini gözle. **Gözlem:** korelasyon mantığını (tersine mühendisliğin temeli) kendi aracında pasif gör. **Sınır:** sadece gözlem; bulunan korelasyonu **göndermeye çalışma**, başka araçta kullanma.

### Alıştırma 3 — Kendi anahtar fob'unun RF sinyalini dinlemek (sadece RX)

**Amaç:** RKE'nin bandını/modülasyonunu somut görmek. **Düzenek:** RTL-SDR + spektrum görüntüleyici (Bölüm 2, 4). **Adımlar:** alıcıyı ~315/433/868 MHz civarına ayarla; kendi fob'unun butonuna bastığında spektrumda kısa bir patlama (burst) ve (mümkünse) modülasyon tipini (OOK/FSK) gözle. **Gözlem:** "her basışta sinyal var, bandı şu" — rolling-code'u çözmeye **çalışma**, yalnızca varlığını/bandını gözle. **Sınır:** yalnızca dinleme (RX); **yayınlama (TX) yok**; başkasının fob'u yok.

### Alıştırma 4 — Faraday kesesinin relay'i nasıl kestiğini test etmek

**Amaç:** relay savunmasını (Bölüm 13) elle deneyimlemek. **Adımlar:** kendi anahtarını bir Faraday kesesine/kutusuna koy; PKES'li aracında anahtar kese içindeyken aracın anahtarı "görmediğini" (kapı açılmaz) gözle, keseden çıkarınca normale döndüğünü doğrula. **Gözlem:** ekranlamanın LF sorgusunu nasıl kestiğini, dolayısıyla relay köprüsünün de neden "duyacak anahtar bulamayacağını" somut gör. **Sınır:** kendi aracın/anahtarın; bu bir savunma testidir, saldırı değil.

### Alıştırma 5 — Kendi TPMS sensör ID'lerini rtl_433 ile görmek

**Amaç:** "araç bile kimlik yayar" gerçeğini (Bölüm 15) somutlaştırmak. **Düzenek:** RTL-SDR + `rtl_433`. **Adımlar:** alıcıyı uygun banda (315/433 MHz) ayarla; `rtl_433` ile kendi aracının lastik sensörlerinin paketlerini (sensör ID + basınç/sıcaklık) gözle (araç hareket ederken daha aktif). **Gözlem:** sabit benzersiz ID'leri ve şifrelemenin yokluğunu gör; bunun bir izleme kanalı olduğunu kavra. **Sınır:** yalnızca kendi aracın; pasif dinleme; sahte paket **yayınlama**.

> Alıştırma çerçevesi (daima): Hepsi **kendi araç + kendi anahtar + pasif/savunma-testi** sınırındadır. Kural: **CAN'e canlı araçta çerçeve gönderme; başka araca dokunma; TX yapma (yalnızca RX/gözlem).** Emin olmadığın her ayrıntıyı (bit rate, adaptör yeteneği, bant, yasal durum) **kaynaktan/üretici belgesinden teyit et.** Amaç, eline aldığın sistemin nasıl konuştuğunu ve nerede zayıf olduğunu **savunma gözüyle** görmektir.

---

<a id="19"></a>
## 19. Hızlı Referans ve Diğer Bölümler

### Kavram kartı

| Kavram | Bir cümlelik öz |
|---|---|
| ECU | Aracın işlevlerini yöneten özelleşmiş gömülü bilgisayar; modern araçta onlarca-yüzlerce |
| CAN bus | Otomotivin baskın iç veriyolu; broadcast, çok-uçlu, ID = öncelik/içerik (kimlik değil) |
| CAN çerçevesi | SOF+ID+DLC+veri(0-8B)+CRC+ACK; gönderen adresi/auth/nonce **yok** |
| Arbitrasyon | Düşük ID = yüksek öncelik; dominant(0) recessive(1)'i bastırır; tahrip etmez |
| CAN güvensizliği | Gizlilik yok (broadcast) + auth yok (spoof) + tazelik yok (replay) — 1980'ler tehdit modeli |
| CAN-FD | 64 bayt yük + yüksek hız; MAC (SecOC) için yer açar |
| LIN/FlexRay/Ethernet | Ucuz-yavaş / zaman-kritik / yüksek-bant; Ethernet kriptoya (MACsec/TLS) uygun |
| OBD-II | İç ağa açılmış standart teşhis kapısı; takılı dongle'lar saldırı yüzeyi |
| SocketCAN / can-utils | Linux CAN soketi + araçlar; `candump`=pasif (güvenli), `cansend`=aktif (tehlikeli) |
| SecOC | Çerçeveye MAC + sayaç → spoof (auth) ve replay (tazelik) birlikte çözülür |
| CAN IDS | Periyot/değer/parmak-izi profilinden anomali (enjeksiyon) yakalar |
| RKE | Butonlu, tek-yönlü UHF (315/433/868), rolling-code; replay'e karşı tasarlı |
| PKES | Butonsuz; LF(125kHz) uyandırma + UHF yanıt; "LF duyuluyor = yakın" varsayımı |
| Immobilizer | Doğru transponder yoksa motoru engelleyen ayrı katman; gücü kripto kalitesinde |
| Transponder kripto | Hitag2/Megamos (eski, kırık) ↔ AES (modern, güçlü) |
| Sabit-kod replay | Aynı kod → yakala+tekrar; çözüm rolling/nonce |
| RollJam (kavram) | Geçerli kodu yakala+jam+sakla+sonra kullan; çözüm zaman penceresi/jam tespiti/UWB |
| Rolling-PWN (kavram, teyit) | Sayaç penceresini kötüye kullanıp eski kodu diriltme iddiası; **kaynaktan teyit** |
| Relay (röle) | Kodu kırmaz, LF/UHF köprüleyip mesafeyi yalanlar; **en yaygın keyless hırsızlığı** |
| Relay savunması | UWB time-of-flight (gerçek mesafe), hareket-sensörlü fob, Faraday kesesi |
| CAN enjeksiyon (kavram) | Far/kablo demetinden CAN'e sahte çerçeve; çözüm SecOC+segment+IDS+fiziksel |
| TPMS | 315/433 MHz, sabit ID → araç-izleme gizlilik riski; çözüm dönüşümlü/şifreli ID |
| Telematik | Hücresel/bulut/uygulama → fiziksel temassız uzak yüzey; **en yüksek ölçek riski** |
| ISO 15118 / V2X | EV şarj / araç-her şey iletişimi; TLS/PKI tabanlı; detay **teyit edilmeli** |
| ISO/SAE 21434 | Aracın ömür boyu siber güvenlik mühendisliği süreci; WP.29 R155/R156 |

### Ezber sezgiler

- Modern araç = tekerlekli, bağlı bir gömülü ağ; güvenliği iç-ağdan (CAN) dış-kapılara (keyless/telematik) iki eksende düşün.
- CAN'in üç günahı **tasarım gereğidir**: broadcast (gizlilik yok), auth yok (spoof), tazelik yok (replay) — 1980'lerde "ağ kapalı" sanıldı.
- CAN'e güvenlik **üst katmandan** eklenir: SecOC (MAC+sayaç), gateway segmentasyonu, IDS — hiçbiri tek başına yetmez.
- `candump` (pasif) güvenli ve öğretici; `cansend` (aktif, canlı araç) tehlikeli ve çoğunlukla yasadışı — **vcan/lab'da öğren**.
- RKE ≠ PKES: RKE butonlu UHF rolling-code (replay'e dirençli); PKES butonsuz LF+UHF yakınlık (relay'e açık).
- Rolling-code **tazelik** verir, **mesafe vermez**; relay tam bu boşluğu kullanır — çözüm **gerçek mesafe ölçümü (UWB)**.
- Relay kriptoyu kırmaz, **mesafeyi yalanlar**; en erişilebilir kullanıcı savunması anahtarı **Faraday kesesinde** tutmaktır.
- Immobilizer "ayrı bir kilittir" (giriş aşılsa bile çalıştırmayı engeller); gücü tamamen **transponder kriptosunun kalitesinde**.
- Saldırı yüzeyi yakınlıktan (RF/relay) **uzaktan-bağlanabilirliğe (telematik/bulut)** kayıyor — ölçek riski en büyük orada.
- Eski/gizli kripto (Hitag2/Megamos, bazı KeeLoq) açığa çıkınca çöker; **açık ve sağlam kripto (AES)** belirleyicidir (Crypto-1 dersinin otomotiv hali).
- Güvenlik bir ürün değil **süreçtir** (ISO/SAE 21434); araştırmanın etiği **sorumlu açıklamadır**.

### Yasal sınır ve perspektif (daima)

Bu bölümdeki tüm teknikler tasarım gereği **kendi aracında ve yetkili test** içindir. Başkasının aracına yetkisiz erişim, kapısını açma, çalıştırma, anahtar sinyalini relay ile uzatma, CAN'e izinsiz müdahale veya immobilizer atlatma; Türkiye'de hırsızlık (TCK 141–142) ve bilişim suçları (TCK 243–245) ile haberleşme gizliliği ihlali (TCK 132–140) kapsamında **ağır suçtur**; ülkene göre teyit et. Bu metin hiçbir araç için somut kilit-açma/çalıştırma reçetesi, mesaj/ID, donanım kurulumu veya zamanlama parametresi vermez; yalnızca **prensip + neden mümkün + savunma** sunar. Hedef operatörlük değil, **mühendislik sezgisi ve savunma**dır: bir otomotiv sistemini eline aldığında nasıl konuştuğunu, nerede zayıf olduğunu ve nasıl korunacağını görebilmek. Emin olmadığın teknik ayrıntıları (özellikle KeeLoq/Hitag2/Megamos etkisi, RollJam/Rolling-PWN'in belirli modeldeki durumu, ISO 15118/V2X ve 21434 ayrıntıları) **kaynaktan teyit et**; bu metin anlama ve savunma içindir.

---

> Kapanış: Modern otomobil, mühendislik tarihinin en görünür "kapalı sanılan ağın açılması" örneğidir. CAN ve kardeş veriyolları, kötü niyetli bir düğümün asla bağlanmayacağı varsayımıyla tasarlandı; anahtarsız sistemler "LF'i duyan yakındır" varsayımına; immobilizer'lar gizli kriptoya güvendi. Her varsayım, bağlanabilirlik ve araştırma ilerledikçe sınandı ve kısmen çöktü — ve her çöküş, aynı beş soruyla onarılıyor: şifreli mi, kimlik doğrulanıyor mu, taze mi, gerçekten yakın mı, kuşatılmış ve izlenir mi? Bir otomotiv güvenlik araştırmacısının işi, bu soruları bir aracın her katmanında — iç veriyolundan bulut servisine — sormak ve cevapları **savunmayı güçlendirmek** için kullanmaktır; asla başkasının aracına zarar vermek için değil. Bir sonraki adım, bu sezgiyi yalnızca **kendi aracının ve anahtarının yasal, pasif gözlemi** üzerinde sınamaktır.

---

Bu bölüm, Kanije Kalesi SIGINT El Kitabı'nın parçasıdır. Tüm bölümler ve önerilen okuma sırası için indekse bakın: [SIGINT_00 — Başlangıç ve İndeks](SIGINT_00_BASLANGIC_INDEX_VE_YASAL.md).

Doğrudan ilgili bölümler:
- [SIGINT_01 — RF Fiziği ve Modülasyon](SIGINT_01_TEMELLER_RF_VE_MODULASYON.md): UHF/LF bantları, OOK/FSK ve menzil mantığının temeli.
- [SIGINT_16 — Kısa Menzilli Kablosuz ve IoT](SIGINT_16_KISA_MENZIL_KABLOSUZ_VE_IOT.md): BLE/RFID/sub-GHz/keyless çerçevesi; MIFARE Crypto-1 dersi.
- [SIGINT_13 — RF Tehdit Manzarası ve Karşı-Önlemler](SIGINT_13_RF_TEHDIT_VE_KARSI_ONLEMLER.md): keyless relay/UWB savunması, KeeLoq, drone RF.
- [SIGINT_23 — Kablosuz Saldırı Vektörleri ve DoS](SIGINT_23_SALDIRI_VEKTORLERI_VE_DOS.md): bus monopolizasyonu/jamming gibi kaynak-tüketme sınıfları.
- [SIGINT_24 — Güncel Zafiyet Manzarası](SIGINT_24_GUNCEL_ZAFIYET_MANZARASI.md): otomotiv saldırı yüzeyinin telematik/buluta kayışı.

İlgili kale rehberleri: `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`, `VERACRYPT_USTALIK_REHBERI.md`.
