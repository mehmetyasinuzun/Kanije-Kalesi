# SIGINT EL KİTABI — BÖLÜM 11: UYDU HABERLEŞMESİ

## SATCOM Nasıl Çalışır, Nasıl Dinlenir (Pasif/Yasal), Nasıl Korunur

> Amaç: Önceki bölümler sinyalin fiziğini (Bölüm 1), SDR donanımını (Bölüm 2), anteni (Bölüm 3), yazılım zincirini (Bölüm 4), protokol çözümlemeyi (Bölüm 5), savunmayı (Bölüm 6), yoğun spektrumda ayıklamayı (Bölüm 7) ve frekans tahsisini (Bölüm 8) verdi. Bu bölüm, yerden 550 km ile 36.000 km arasında dönen bir vericiyi konu alır. Uydu haberleşmesi (SATCOM), karasal radyonun aynı fiziği üzerine kurulur ama üç şeyi uca taşır: mesafe (yol kaybı astronomik büyür), hız (alçak yörünge uydusu gökyüzünü dakikalar içinde geçer) ve görüş alanı (tek bir GEO uydusu bir kıtayı görür). Hedefimiz operatör reçetesi değil mühendislik sezgisidir: bir spektrumda uydu taşıyıcısını gördüğünde arkasındaki yörüngeyi, link bütçesini ve transponderi zihninde canlandırabilmen; meteoroloji uydusundan kendi elinle görüntü alabilmen; ve neden bazı eski SATCOM sistemlerinin savunmasız olduğunu, bunun nasıl tespit edildiğini anlayabilmen.

> Yasal çerçeve: Bu bölüm de serinin geri kalanı gibi anlama, savunma ve spektrum okuryazarlığı amaçlıdır. Anlatılan her şey pasif alımdır (RX): meteoroloji uyduları, açık çerçevede gözlemlenen ağ sinyalleri, amatör uydular ve SatNOGS gözlemi. Hiçbir yetkisiz uplink, transponder kullanımı veya şifreli içerik çözme önerilmez. Uydu downlink'i çoğu açık yayın için dinlenebilir olsa da, abonelik/şifreli içeriğin çözülmesi ve **her türlü uplink/iletişim** lisans ve yetki ister; bu kitap bunun adımlarını vermez. Bölüm boyunca "yetkisiz erişim nasıl yapılır" değil, "mekanizma nasıl çalışır, neden savunmasız olabilir, nasıl tespit edilir" sorularına yanıt verilir. Kendi ülkenin ve sürümünün mevzuatını teyit et.

---

## İÇİNDEKİLER

1. [Yörünge Sınıfları: LEO, MEO, GEO, HEO](#1)
2. [Yörünge Geometrisi: Neden GEO Sabit, Neden LEO Hızlı Geçer](#2)
3. [Doppler Kayması: Alçak Yörüngenin İmzası](#3)
4. [Uydu Link Mimarisi: Uplink, Downlink, Transponder](#4)
5. [Bent-Pipe ve İşlemeli (Regenerative) Transponder](#5)
6. [Çapraz-Link, Yer İstasyonu ve TT&C](#6)
7. [Frekans Bantları: L/S/C/X/Ku/Ka ve Yağmur Sönümlemesi](#7)
8. [Uyduda Link Bütçesi: EIRP, Yol Kaybı, G/T, C/N0](#8)
9. [Modülasyon ve Kodlama: DVB-S/S2/S2X, QPSK/8PSK/APSK, FEC](#9)
10. [Spektrumda Bir Uydu Taşıyıcısını Tanımak](#10)
11. [Pasif/Yasal Dinleme — Meteoroloji Uyduları (NOAA APT, Meteor LRPT, GOES)](#11)
12. [Inmarsat, Iridium ve Amatör Uydular (Açık Çerçevede Gözlem)](#12)
13. [SatNOGS: Dağıtık Yer İstasyonu Ağı](#13)
14. [Donanım: Anten, LNB, Takip ve SDR Zinciri](#14)
15. [Güvenlik (Savunma Perspektifi): Neden Bazı SATCOM Zayıftır](#15)
16. ["Satellite Piracy" Olgusu: Mekanizma, Yasadışılık, Tespit](#16)
17. [VSAT/Terminal Saldırı Yüzeyi ve TT&C Güvenliği](#17)
18. [Alıştırmalar (Yasal)](#18)
19. [Hızlı Referans ve Diğer Bölümler](#19)

---

<a id="1"></a>
## 1. Yörünge Sınıfları: LEO, MEO, GEO, HEO

Bir uyduyu anlamanın başlangıcı, onun nereye yerleştirildiğidir. Yörünge yüksekliği tek bir sayı değil; gecikmeyi, kapsama alanını, geçiş süresini, Doppler kaymasını, anten gereksinimini ve link bütçesini aynı anda belirleyen bir karardır. Yörünge sınıfını bildiğinde, sinyalin nasıl davranacağını büyük ölçüde tahmin edebilirsin.

Uyduların yüksekliği genelde dört ana kuşağa ayrılır. Sınırlar doktrine göre biraz oynar; aşağıdaki değerler yaygın mühendislik kabulleridir.

| Sınıf | Açılım | Tipik yükseklik | Yörünge periyodu | Tek geçiş görünürlüğü | Tek başına kapsama |
|---|---|---|---|---|---|
| **LEO** | Low Earth Orbit | ~300–2.000 km | ~90–120 dakika | birkaç–~15 dakika | dar (yerel) |
| **MEO** | Medium Earth Orbit | ~2.000–35.786 km arası (tipik ~20.000 km) | ~2–12 saat | saatler | bölgesel-geniş |
| **GEO** | Geostationary Orbit | 35.786 km (ekvator, dairesel) | ~23h 56m (sidereal gün) | **sürekli (sabit görünür)** | kıtasal (~1/3 yerküre) |
| **HEO** | Highly Elliptical Orbit | apoje ~40.000 km'nin üstü, perije birkaç yüz km | ~12 saat (örn. Molniya) | apoje yakınında saatlerce | yüksek enlemler |

Buradan çıkan ilk pratik sonuç: gecikme. Sinyal ışık hızıyla gider; yükseklik arttıkça gidiş-dönüş gecikmesi (round-trip) artar.

```
  Tek yön gecikme (yaklaşık, dik geçiş):
  LEO  (~700 km)   :  700 km  / 300.000 km/s ≈ 2,3 ms
  MEO  (~20.000 km):  ~67 ms
  GEO  (~36.000 km):  ~120 ms (tek yön)  →  uplink+downlink ≈ 240 ms
  GEO çift atlama  :  ≈ 480–600 ms (yer-uydu-yer-uydu-yer)
```

Bu yüzden GEO üzerinden geçen sesli görüşmede o tanıdık "boşluk" duyulur: yaklaşık çeyrek saniyelik tek yön, yarım saniyeye yakın gecikme. Telekonferans veya gerçek zamanlı oyun için GEO kötüdür; bir kıtayı tek uyduyla yayınlamak (TV) için mükemmeldir. LEO ise gecikmeyi karasal fiberle yarışacak düzeye indirir (birkaç ms), ama bunun bedeli aşağıda göreceğimiz gibi sürekli hareket ve çok sayıda uydu ihtiyacıdır.

İkinci pratik sonuç: kapsama ve takım (constellation) büyüklüğü. Bir GEO uydusu ekvator üzerinde sabit durduğu için tek bir uydu, görüş açısındaki her yere sürekli hizmet verir; üç GEO uydusu (yaklaşık 120° aralıkla) kutuplar hariç neredeyse tüm yerküreyi kaplar. LEO'da ise her uydu yalnızca dakikalarca görünür; kesintisiz hizmet için onlarca, yüzlerce, bazen binlerce uydudan oluşan bir takım gerekir (modern geniş bant LEO takımları bunun örneğidir). Mühendislik dengesi nettir: **GEO az uydu + büyük gecikme + sabit anten; LEO çok uydu + küçük gecikme + hareketli/izleyen anten.**

HEO, bu dengeyi özel bir coğrafya için kırar. Yüksek enlemler (kutba yakın bölgeler) GEO'dan zor görülür, çünkü GEO ekvator düzlemindedir ve kuzeyden bakınca ufka çok yakın kalır. Molniya tipi yüksek eliptik yörünge, apojesini (en uzak nokta) kuzey yarımküre üzerinde tutarak uydunun günün büyük bölümünü o bölgenin "tepesinde yavaşça asılı" geçirmesini sağlar (Kepler'in ikinci yasası: apojede yavaş, perijede hızlı hareket). Böylece birkaç HEO uydusu, GEO'nun ulaşamadığı yüksek enlemlere GEO benzeri süreklilik sunar.

Not: GEO terimi sıkça GSO (geosenkron) ile karıştırılır. Geosenkron, periyodu bir sidereal güne eşit her yörüngedir; geostationary ise bunun özel hali olup ekvatorda, dairesel ve eğimsiz (inclination ≈ 0) olandır ve gökyüzünde gerçekten sabit görünür. Eğimli bir geosenkron uydu, gökyüzünde gün boyunca bir "8" (analemma) çizer; tam sabit durmaz. Bu ayrım, çanak hizalamada önemlidir.

---

<a id="2"></a>
## 2. Yörünge Geometrisi: Neden GEO Sabit, Neden LEO Hızlı Geçer

Bir uydunun gökyüzünde neden sabit ya da hızlı göründüğünü anlamak için tek bir denkleme bakmak yeterlidir: yörünge periyodu yalnızca yörünge yarıçapına bağlıdır. Kepler'in üçüncü yasası, dairesel yörünge için şu sade biçimi alır:

```
        ┌─────────────────────────┐
        │   T = 2π · √( a³ / μ )   │
        └─────────────────────────┘

  T = yörünge periyodu (s)
  a = yörünge yarıçapı = R_dünya + yükseklik  (m)   [R_dünya ≈ 6.371 km]
  μ = G·M_dünya ≈ 3,986 × 10^14 m³/s²  (Dünya'nın standart yerçekim parametresi)
```

Bu denklem, periyodun uydunun kütlesinden tamamen bağımsız olduğunu söyler — bir uydu da bir vida da aynı yükseklikte aynı sürede döner. Periyot yalnızca `a`'ya (dolayısıyla yüksekliğe) bağlıdır. Şimdi GEO'nun sihrini çözelim: GEO yüksekliği (35.786 km) tam olarak öyle seçilmiştir ki periyot bir **sidereal güne** (23 saat 56 dakika 4 saniye — Dünya'nın yıldızlara göre tam bir tur süresi) eşit olsun. Uydu Dünya ile aynı açısal hızda döndüğü için, yerden bakan bir gözlemciye göre gökyüzünde aynı noktada **asılı kalır**. Çanağını bir kez o noktaya kilitlersin, bir daha oynatmazsın.

```
   GEO — uydu Dünya ile birlikte döner (göreli olarak sabit):

              ☉ uydu (sabit görünür)
               \
                \  36.000 km
                 \
        ┌─────────●─────────┐   ← yer istasyonu (çanak sabit kilitli)
        │      DÜNYA         │
        │   (döner, ama uydu │
        │    aynı hızda)     │
        └───────────────────┘

   LEO — uydu Dünya'dan çok daha hızlı döner (gökyüzünü tarar):

      ufuk        zenit        ufuk
   (doğuş) ──────► (tepe) ──────► (batış)
        \           |           /
         \          |          /        ~700 km
          \         |         /
   ════════●════════●════════●════════  ← gözlemci
        görünür süre: birkaç–15 dakika
        çanak/anten uyduyu TAKİP etmeli (azimut+yükseklik)
```

LEO uydusu (~700 km) için aynı denklem ~98 dakikalık periyot verir; uydu Dünya çevresinde Dünya'nın kendi dönüşünden çok daha hızlı döner. Yerden bakınca ufuktan doğar, gökyüzünü kabaca dik bir yay çizerek geçer ve karşı ufuktan batar. Tepe noktasına yakın bir geçişte toplam görünürlük birkaç ile on beş dakika arasındadır; ufka teğet geçişlerde çok daha kısadır. İşte bu yüzden LEO alımında ya elle çevirdiğin yönlü bir anten ya da motorlu (azimut+yükseklik) bir takip sistemi ya da geniş açıyı kapsayan yön bağımsız bir anten (QFH, turnike) gerekir — bunları Bölüm 14'te ele alacağız.

Yükseklik arttıkça periyodun nasıl uzadığını birkaç örnekle görelim (yaklaşık):

| Yükseklik | Yörünge yarıçapı a | Periyot T | Davranış |
|---|---|---|---|
| 400 km (ISS) | ~6.771 km | ~92,7 dk | hızlı LEO; gökyüzünü ~5–10 dk'da geçer |
| 700 km (NOAA/Meteor) | ~7.071 km | ~98,8 dk | tipik kutupsal meteoroloji LEO |
| ~20.200 km (GNSS) | ~26.571 km | ~11h 58m | MEO; yavaş sürüklenir |
| 35.786 km (GEO) | ~42.157 km | ~23h 56m | sabit görünür |

Kutupsal yörünge (polar / güneş-senkron) bir ek incelik taşır: yörünge düzlemi neredeyse kutuplardan geçecek şekilde eğimlidir (~98°). Böylece Dünya altında dönerken uydu, her turda dünyanın farklı bir boylamını tarar ve gün içinde tüm yerküreyi şeritler halinde kapsar. Meteoroloji uydularının kutupsal olmasının nedeni budur: tek bir uydu, art arda turlarla tüm gezegeni görüntüler. Güneş-senkron varyantı, uydunun her geçişte aynı yerel güneş saatinde geçmesini sağlar (gölgeler tutarlı, görüntüler karşılaştırılabilir).

Pratikte: Bir uydudan görüntü almayı planlarken iki şeyi bilmen yeter — uydunun ne zaman senin gökyüzünden geçeceği (geçiş tahmini, "pass prediction") ve geçişin maksimum yükseklik açısı (elevation). Yüksek açılı (zenite yakın) geçişler en uzun ve en temiz sinyali verir; ufka yakın geçişler kısa ve gürültülüdür (atmosferden uzun yol + yer engelleri). Geçiş tahmini için kullanılan veri, Bölüm 13'te değineceğimiz TLE (Two-Line Element) yörünge parametreleridir.

---

<a id="3"></a>
## 3. Doppler Kayması: Alçak Yörüngenin İmzası

Hızlı hareket eden bir vericinin sinyali, yaklaşırken yukarı, uzaklaşırken aşağı frekansa kayar. Bu, ambulans sireninin perdesinin geçerken düşmesinin radyo karşılığıdır ve LEO uydu alımının en belirgin imzasıdır. Karasal radyoda çoğu zaman ihmal edilen Doppler, LEO'da merkez frekansı bantın bir ucundan diğerine sürükleyecek kadar büyür.

Doppler frekans kayması yaklaşık olarak:

```
        ┌──────────────────────────────┐
        │   Δf ≈ f₀ · (v_radyal / c)    │
        └──────────────────────────────┘

  Δf       = frekans kayması (Hz)
  f₀       = nominal taşıyıcı frekansı (Hz)
  v_radyal = uydunun gözlemciye doğru/uzak BİLEŞEN hızı (m/s)  [yaklaşırken +, uzaklaşırken −]
  c        = ışık hızı ≈ 3 × 10⁸ m/s
```

Kritik nokta `v_radyal`: önemli olan uydunun yörünge hızının tamamı (~7,5 km/s LEO için) değil, gözlemciye doğru olan **bileşeni**dir. Uydu tam tependen (zenitten) geçerken o an sana ne yaklaşıyor ne uzaklaşıyordur; radyal hız sıfırdır ve Doppler kayması o anda sıfırdan geçer. Doğuşta sana hızla yaklaşır (maksimum pozitif kayma), tepede sıfırlanır, batışta hızla uzaklaşır (maksimum negatif kayma). Bu yüzden tipik bir LEO geçişinde merkez frekans bir "S" eğrisi çizer:

```
   Frekans kayması (Δf) — tipik LEO geçişi:

   +Δf ┤●                            (doğuş: yaklaşıyor, yukarı kaymış)
       │ ●
       │   ●
     0 ┤─────●──────────────────────  (zenit: radyal hız 0, kayma yok)
       │        ●
       │          ●
   −Δf ┤            ●                 (batış: uzaklaşıyor, aşağı kaymış)
       └────────────────────────────► zaman (geçiş boyunca ~10 dk)
```

Büyüklük hissi için: 137 MHz'de (NOAA APT) tipik LEO geçişinde toplam kayma ±3 kHz mertebesindedir — APT'nin geniş FM bandı (~34 kHz) içinde yutulur, yazılım çoğunlukla bunu görmezden gelebilir. Ama 145/435 MHz amatör uydularda dar bantlı SSB/CW için kayma birkaç kHz'i bulur ve sinyali dar filtreden kaçırır; orada Doppler düzeltmesi şarttır. Iridium gibi 1.6 GHz'de çalışan LEO sistemlerde kayma onlarca kHz'e ulaşır. Genel kural: **frekans yükseldikçe ve bant daraldıkça Doppler daha çok canını sıkar.**

Doppler ile başa çıkmanın üç yolu vardır. Birincisi yazılımla otomatik düzeltme: program, uydunun TLE'sinden anlık radyal hızı hesaplar ve SDR'ın ayar frekansını geçiş boyunca canlı kaydırarak sinyali merkezde tutar (gpredict + rigctld zinciri bunu yapar). İkincisi geniş bant alıp sonradan işleme: APT gibi sinyalleri yeterince geniş örnekleyip kaymayı çözümlemede telafi etmek. Üçüncüsü, dar bantlı amatör çalışmada elle ayar — pratik ama yorucu. Mühendislik açısından Doppler bir sorun olduğu kadar bir **bilgi kaynağıdır** da: kayma eğrisinin şekli ve sıfır geçiş anı, uydunun ne zaman tepende olduğunu ve hangi yöne gittiğini söyler; SatNOGS gibi ağlar bu eğrileri otomatik analize katar.

Not: GEO uydularında Doppler pratikte ihmal edilebilir (uydu sabit görünür, radyal hız ~0). Hafif yörünge tutma manevraları ve eğim salınımı nedeniyle çok küçük bir artık Doppler kalabilir ama bu, tüketici alımında önemsizdir. Doppler esas olarak LEO/MEO problemidir.

---

<a id="4"></a>
## 4. Uydu Link Mimarisi: Uplink, Downlink, Transponder

Bir uydu haberleşme bağı, en sade haliyle üç parçadır: yerden uyduya giden yol (uplink), uydunun içindeki işlem zinciri (transponder), uydudan yere dönen yol (downlink). Bu üçlü, ister bir TV yayını ister bir VSAT internet bağı isterse bir uydu telefonu olsun, aynı iskelet üzerinde kurulur.

```
                         ┌───────────────────┐
                         │      UYDU          │
                         │  ┌─────────────┐   │
              UPLINK     │  │ TRANSPONDER │   │   DOWNLINK
           (yüksek f) ──▶│  │ LNA→çevrim→ │   │──▶ (düşük f)
           örn. 14 GHz   │  │  HPA        │   │   örn. 11 GHz
                         │  └─────────────┘   │
                         └───────────────────┘
                ▲                                      │
                │                                      ▼
        ┌───────────────┐                      ┌───────────────┐
        │ YER İSTASYONU │                      │  ALICI(LAR)   │
        │  (uplink)     │                      │  VSAT / TV /  │
        │  büyük çanak  │                      │  el terminali │
        │  + HPA        │                      │  + LNB        │
        └───────────────┘                      └───────────────┘
```

Neden uplink ve downlink farklı frekanslardadır? Çünkü uydu aynı anda hem dinleyip hem yayın yapar; aynı frekansta olsalardı kendi güçlü vericisi kendi hassas alıcısını sağır ederdi. Bu yüzden her transponder bir uplink bandını alır, sabit bir miktar **frekans çevrimi** (frequency translation) uygular ve daha alçak bir downlink bandında geri yayınlar. Örneğin tipik bir Ku-bandı transponderinde uplink ~14 GHz, downlink ~11 GHz civarındadır; arada sabit bir yerel osilatör farkı vardır. Bu çevrim, transponderin temel kimliğidir.

Bir uydu tek bir transponder taşımaz; onlarca transponder, uydunun toplam bant genişliğini dilimlere böler. Klasik bir yayın uydusunda her transponder ~36 MHz (veya 54/72 MHz) genişliğinde bir "boru" gibidir ve içinden bir veya birkaç taşıyıcı geçer. Operatör, müşterilere bu boruları (ya da bir borunun bir kısmını, MHz cinsinden kapasite) kiralar. Spektrumda bir GEO uydusuna baktığında gördüğün şey, yan yana dizilmiş bu transponder bloklarıdır.

Uplink ve downlink kazanç/güç dengesi de asimetriktir. Uplink tarafında yer istasyonu büyük çanak ve yüksek güçlü yükselteç (HPA — High Power Amplifier) kullanabilir; yerde yer, soğutma ve enerji boldur. Downlink tarafında ise uydunun gücü kısıtlıdır (güneş paneli + sınırlı HPA) ve bu güç bütün kapsama alanına yayılır; bu yüzden yerdeki alıcının işi zordur ve düşük gürültülü ön kademe (LNB/LNA) ile büyük anten gerektirir. Link bütçesi bölümünde (Bölüm 8) bu dengesizliğin sayısal nedenini göreceğiz.

Çift yönlü (interaktif) sistemlerde — VSAT internet gibi — terminal hem alır hem de yayınlar; yani küçük yer istasyonu rolünü üstlenir ve kendi uplink'ini yapar. Tek yönlü sistemlerde (TV yayını) ise yerdeki kullanıcı yalnızca alır; uplink'i yalnızca yayıncının merkez istasyonu yapar. Bu ayrım, ileride güvenlik tartışmasında kritik olacaktır: yalnızca alan bir terminal pasiftir; yayın yapan bir terminal spektrumda iz bırakır ve yetkisiz kullanımı tespit edilebilir kılar.

---

<a id="5"></a>
## 5. Bent-Pipe ve İşlemeli (Regenerative) Transponder

Transponderin "boru" benzetmesi tesadüf değildir; en yaygın uydu transponderi gerçekten de bir borudur. İki temel mimari vardır ve aralarındaki fark, hem performansı hem de güvenliği doğrudan belirler.

**Bent-pipe (şeffaf / transparent) transponder** sinyali içeriğine hiç bakmadan işler. Yaptığı tek şey şudur: uplink sinyalini al, düşük gürültüyle yükselt (LNA), frekansını çevir (downlink bandına indir), güçlü yükselt (HPA) ve geri yayınla. Uydu, sinyalin ne taşıdığını, hangi modülasyonu kullandığını, şifreli olup olmadığını **bilmez ve umursamaz**. Ne girdiyse onu, sadece frekansı kaydırılmış ve güçlendirilmiş halde geri verir. "Bent pipe" (bükülmüş boru) adı buradan gelir: sinyal yerden çıkar, uydunun tepesinde 90 derece "bükülür" ve yere geri döner.

```
   BENT-PIPE (şeffaf) transponder — içeriğe bakmaz:

   uplink ──▶ [ LNA ] ──▶ [ mikser/LO: f çevrimi ] ──▶ [ HPA ] ──▶ downlink
   (14 GHz)                 (sadece frekans kaydırır)            (11 GHz)

   ► Demodülasyon YOK, hata düzeltme YOK, yeniden kodlama YOK.
   ► Uydu sinyali "anlamaz"; sadece taşır.
   ► Gürültü de yükselir (uplink gürültüsü downlink'e taşınır).


   REGENERATIVE (işlemeli / on-board processing) transponder:

   uplink ──▶ [ LNA ] ──▶ [ DEMOD ] ──▶ [ baseband işleme: ]
                                         │ - bit'leri geri kazan
                                         │ - FEC ile hata düzelt
                                         │ - yönlendir/anahtarla
                                         ▼
                          [ yeniden MOD ] ──▶ [ HPA ] ──▶ downlink

   ► Uydu sinyali bit düzeyinde geri kazanır (regenerate eder).
   ► Uplink gürültüsü TEMİZLENİR; downlink taze sinyaldir.
   ► Uydu içinde anahtarlama/yönlendirme mümkün (paket switch).
```

Bent-pipe'ın avantajı sadelik ve esnekliktir: uydu hangi modülasyon/standart kullanıldığını bilmediği için, yıllar sonra yerdeki ekipman yeni bir standarda (örneğin DVB-S'ten DVB-S2'ye) geçtiğinde uyduya dokunmak gerekmez — uydu zaten içeriğe kayıtsızdır. Dezavantajı, uplink'te eklenen gürültünün downlink'e aynen taşınmasıdır (uydu sinyali temizlemez), ve aşağıda göreceğimiz güvenlik açığıdır.

**İşlemeli (regenerative / on-board processing, OBP) transponder** ise sinyali uydu üzerinde gerçekten **demodüle eder**: uplink'i baseband'e indirir, bit'leri geri kazanır, hata düzeltmesini (FEC) uygular, gerekirse paketleri uydu içinde anahtarlar/yönlendirir, sonra taze bir taşıyıcı olarak yeniden modüle edip yayınlar. Avantajı, uplink ve downlink gürültüsünün birbirinden ayrılmasıdır (uydu temiz bit aldıktan sonra yeniden ürettiği için downlink gürültüsüzdür) ve uydu içinde akıllı yönlendirme yapılabilmesidir. Dezavantajı karmaşıklık, güç tüketimi ve standarda bağımlılıktır — uydu belirli bir dalga biçimini "bildiği" için, o standardı sonradan değiştirmek genelde mümkün değildir.

Güvenlik açısından bu ayrım belirleyicidir ve Bölüm 15'te derinleşeceğiz; ama tohumu şimdi atalım: **bent-pipe transponder içeriğe kayıtsız olduğu için, uyduya yetkili uplink erişimi olan herhangi bir taşıyıcıyı geri yayınlar — uydu "bu yetkili mi" diye soramaz.** Erişim kontrolü uydunun kendisinde değil, yer tarafındaki idari/operasyonel düzeydedir (kim hangi transponderi kiralamış, hangi yer istasyonu yetkili). Bu, eski/klasik şeffaf transponderlerin neden yetkisiz kullanıma ("satellite piracy") açık olabileceğinin teknik temelidir. İşlemeli transponderler, kimlik doğrulama ve protokol farkındalığını uyduya taşıyabildiği için bu açığı daraltma potansiyeli sunar — ama bunun garanti değil, tasarım tercihi olduğunu vurgulamak gerekir.

---

<a id="6"></a>
## 6. Çapraz-Link, Yer İstasyonu ve TT&C

Uydu mimarisinin geri kalan üç parçası — çapraz-link, yer istasyonu ve TT&C — bir takımı tek tek uydulardan bir **ağa** dönüştürür ve uydunun yörüngede sağlıklı kalmasını sağlar.

**Çapraz-link (inter-satellite link, ISL):** Klasik mimaride her sinyal yere iner ve yerden tekrar yükselir (uydu-yer-uydu). Modern LEO takımlarında ise uydular birbirleriyle doğrudan, uzayda haberleşir — genelde lazer (optik) ya da yüksek frekanslı RF bağlarıyla. Çapraz-link, bir kullanıcı paketini yere hiç indirmeden uydudan uyduya zıplatarak hedefe en yakın uyduya taşır.

```
   Yer-merkezli (klasik):           Çapraz-linkli (modern LEO):

   kullanıcı                         kullanıcı
      │ uplink                          │
      ▼                          ☉──ISL──☉──ISL──☉
   ☉ uydu                        │  (uzayda zıplama) │
      │ downlink                 ▼                   ▼
      ▼                       yakın hiç inmeden    hedef
   GATEWAY (yer)              gateway gerekmeden    kullanıcı
      │                       paket taşınır
      ▼
   internet / hedef
```

Çapraz-linkin avantajı, yeryüzünde gateway (geçit yer istasyonu) bulunmayan bölgeler (okyanus ortası, kutuplar) üzerinde dahi hizmet verebilmek ve gecikmeyi azaltmaktır (ışık boşlukta camdan daha hızlıdır; uzun mesafede uzaydan zıplama, yer fiberinden hızlı olabilir). Pasif gözlem açısından önemli nokta: çapraz-link sinyalleri yere inmez, bu yüzden yerden alınamaz — bunlar uydular arası, dar hüzmeli ve genelde yere yönelmemiş bağlardır.

**Yer istasyonu / gateway:** Uydunun yerle buluştuğu büyük altyapıdır. Bir gateway, büyük çanaklar, yüksek güçlü uplink yükselteçleri, hassas alıcılar ve karasal omurgaya (internet/şebeke) bağlantı taşır. Yayın uydularında bu rol "uplink merkezi"dir (içeriği uyduya basan tek yetkili nokta). VSAT ağlarında "hub" denir ve yüzlerce uzak terminali tek merkezden yönetir (yıldız topoloji: bütün terminaller hub üzerinden konuşur).

**TT&C (Telemetry, Tracking and Command — Telemetri, Takip ve Komuta):** Bu, uydunun haberleşme yükünden (payload) ayrı, uydunun **kendisini yöneten** hayati bağdır. Üç işlevi vardır:

- **Telemetri (T):** Uydu yerle kendi sağlık durumunu paylaşır — batarya gerilimi, sıcaklık, yönelim (attitude), güneş paneli akımı, yakıt durumu. Yerdeki operasyon merkezi uydunun "nabzını" buradan dinler.
- **Takip (Tracking):** Uydunun tam konumu ve yörüngesi izlenir (mesafe ölçümü/ranging, açı takibi). Yörünge tahmini ve manevra planlaması buna dayanır.
- **Komuta (C):** Yerden uyduya gönderilen komutlar — yörünge düzeltme manevrası ateşle, antenleri yönlendir, transponderi aç/kapat, mod değiştir. Uydu bu komutlarla yönetilir.

```
   TT&C bağı — uydunun "yönetim konsolu":

        ┌────────────────────────────┐
        │           UYDU             │
        │  payload (haberleşme yükü) │  ← kullanıcı trafiği (ayrı)
        │  ─────────────────────────  │
        │  BUS (uydunun gövdesi):     │
        │   batarya, ADCS, itki,      │
        │   bilgisayar  ◀── KOMUTA    │
        │               ──▶ TELEMETRİ │
        └──────────┬─────────────────┘
                   │  güvenli (olması gereken) TT&C bağı
                   ▼
        ┌────────────────────────────┐
        │   UYDU KONTROL MERKEZİ      │
        │   (operatör, sadece yetkili)│
        └────────────────────────────┘
```

TT&C'nin güvenliği, tüm sistemin güvenliğinin temelidir ve Bölüm 17'de ayrıca ele alacağız. Mantığı basittir: payload'a (kullanıcı trafiğine) yetkisiz erişim kötüdür ama TT&C'ye yetkisiz erişim **felakettir** — çünkü komuta bağını ele geçiren biri uydunun yönelimini, transponderlerini, hatta yörüngesini etkileyebilir; yani uydunun kendisini kaybettirebilir. Bu yüzden TT&C bağları, payload'dan çok daha sıkı korunmalı (kimlik doğrulama, şifreleme, komut yetkilendirme) ve genelde ayrı, sınırlı erişimli frekanslarda çalışır. Pratikte eski uyduların bir kısmında TT&C güvenliği bugünkü standartların gerisindedir; bu, savunma perspektifinin ana endişe alanlarından biridir.

---

<a id="7"></a>
## 7. Frekans Bantları: L/S/C/X/Ku/Ka ve Yağmur Sönümlemesi

Uydu haberleşmesi, mikrodalga spektrumunun belirli "harf" bantlarında yoğunlaşır. Her bant, bir denge noktasıdır: düşük frekans daha az yol kaybı ve yağmura dayanıklılık ama daha az bant genişliği ve daha büyük anten; yüksek frekans bol bant genişliği ve küçük anten ama yağmura duyarlılık ve daha çok yol kaybı sunar. Bir bandın "kime/neye" ait olduğunu bilmek, spektrumda uydu sinyalini tanımanın yarısıdır. (Karasal bant adlandırması ve genel tahsis için bkz. Bölüm 8; burada uyduya özgü kullanım vurgulanır.)

| Bant | Yaklaşık frekans | Tipik uydu kullanımı | Anten boyu | Yağmur etkisi |
|---|---|---|---|---|
| **L** | ~1–2 GHz | Mobil uydu (Inmarsat, Iridium, Thuraya), GNSS (GPS/Galileo), meteoroloji (137 MHz APT teknik olarak VHF, ama düşük-band karakteri benzer) | küçük/yama, heliks | ihmal edilebilir |
| **S** | ~2–4 GHz | TT&C, bazı mobil uydu, uzay araştırma, meteoroloji veri bağı | orta | çok az |
| **C** | ~4–8 GHz | Klasik TV/veri yayını, uluslararası bağlar; yağmura dayanıklı olduğu için tropik bölgede tercih | büyük çanak (1,8–3 m+) | düşük |
| **X** | ~8–12 GHz | Askeri SATCOM, devlet/uzay araştırma (sivil tahsis kısıtlı) | orta-büyük | düşük-orta |
| **Ku** | ~12–18 GHz | DTH TV (uydu çanağı yayını), VSAT, kurumsal veri | küçük çanak (60–120 cm) | orta (şiddetli yağmurda kesinti) |
| **Ka** | ~26,5–40 GHz | Yüksek hızlı geniş bant internet (modern VSAT/HTS), nokta-hüzme | çok küçük çanak | yüksek (yağmura çok duyarlı) |

Bandın seçimi neredeyse tamamen iki gerilim arasındadır: **bant genişliği iştahı** (yukarı çıkar) ile **yağmur sönümlemesi** (yukarı çıktıkça kötüleşir). Bunu somutlaştıralım.

**Yağmur sönümlemesi (rain fade):** Atmosferdeki yağmur damlaları, dalga boyu damla çapına yaklaştıkça sinyali daha çok soğurur ve saçar. Damla çapı milimetre mertebesindedir; frekans yükselip dalga boyu milimetreye indikçe (Ku, özellikle Ka) yağmur ciddi bir engel olur. C-bandında (dalga boyu birkaç santimetre) yağmur büyük ölçüde "görünmez"; bu yüzden tropik/muson bölgelerde ve yüksek güvenilirlik gereken yayınlarda C-bandı tercih edilir. Ka-bandında ise şiddetli bir sağanak, linki tamamen kesebilir (birkaç dB değil, onlarca dB sönümleme).

```
   Yağmur sönümlemesi — frekansla artar:

   sönümleme (dB)
      │                                          ╱ Ka (çok duyarlı)
      │                                    ╱╱╱
      │                            ╱╱╱  Ku
      │                  ╱╱╱╱
      │        ╱╱╱  C, X
      │ ───  L, S  (neredeyse düz, düşük)
      └──────────────────────────────────────────► frekans
        1GHz    4GHz    12GHz    20GHz    30GHz+
```

Bu nedenle Ka/Ku sistemleri **uplink power control** (yağmur algılanınca yer istasyonu gücünü artırma), site diversity (coğrafi olarak ayrık ikinci gateway; biri yağmurdayken diğeri açık) ve adaptif kodlama (DVB-S2'nin ACM özelliği: yağmurda daha güçlü FEC'e düşme, açıkta daha hızlı moda çıkma) gibi tekniklerle savunulur. Mühendislik sezgisi: **bir uydu sisteminin bandını duyduğunda, yağmurla başının nasıl dertte olduğunu da öğrenmiş olursun.**

L-bandının özel yeri: el terminalleri (uydu telefonu) ve GNSS L-bandındadır çünkü bu frekanslarda hem yol kaybı görece düşük hem de küçük, yön bağımsız antenle (yama, kısa heliks) çalışılabilir — kullanıcı çanağı tepeye doğrultmak zorunda kalmaz. Bunun bedeli düşük bant genişliğidir; bu yüzden uydu telefonu sesi dardır, veri hızı düşüktür. Meteoroloji uydularının VHF/L benzeri düşük bant kullanması da aynı mantıkladır: basit, yön bağımsız antenle herkesin alabilmesi istenir.

Not: Bant sınırları (özellikle C/X/Ku/Ka kesişimleri) kaynaktan kaynağa birkaç yüz MHz oynayabilir; yukarıdaki değerler yaygın mühendislik aralıklarıdır, kesin tahsis için ITU/ulusal tabloya bakılmalı (Bölüm 8). "X-bandı askeri" genellemesi de mutlak değildir; sivil/araştırma kullanımları vardır.

---

<a id="8"></a>
## 8. Uyduda Link Bütçesi: EIRP, Yol Kaybı, G/T, C/N0

Bir uydu bağının çalışıp çalışmayacağı tek bir hesaba bağlıdır: link bütçesi. Bu, vericiden çıkan gücün, yol boyunca yitirilen her şeyin ve alıcıda toplanan her şeyin dB cinsinden toplamıdır. Karasal radyoda link bütçesini Bölüm 1'de gördük (FSPL); uyduda aynı denklem geçerlidir ama mesafe terimi devasa olduğu için her parça acımasızca önemlidir. 36.000 km, neden büyük çanak ve düşük gürültülü LNB gerektiğini tek başına açıklar.

Önce serbest uzay yol kaybı (FSPL — Free Space Path Loss), yani sinyalin sadece "yayılarak" zayıflaması:

```
        ┌──────────────────────────────────────────┐
        │  FSPL(dB) = 20·log₁₀(d) + 20·log₁₀(f) + 92,45 │
        │     d: km cinsinden,  f: GHz cinsinden        │
        └──────────────────────────────────────────┘

  GEO downlink örneği:  d = 36.000 km,  f = 12 GHz (Ku)
  FSPL = 20·log₁₀(36000) + 20·log₁₀(12) + 92,45
       = 20·(4,556) + 20·(1,079) + 92,45
       ≈ 91,1 + 21,6 + 92,45
       ≈ 205 dB
```

**205 dB.** Bu sayının büyüklüğünü hissetmek gerekir: 205 dB, gücün 10^20,5'e (yaklaşık 300 milyar kere milyar) bölünmesi demektir. Uydudan çıkan birkaç on watt'lık sinyal, yere ulaştığında pikowatt'ın altına iner. İşte bu yüzden GEO alımında ne yapılırsa yapılsın büyük bir anten kazancı (çanak) ve son derece düşük gürültülü bir ön kademe (LNB) zorunludur — başka türlü sinyal gürültünün altında kalır.

Tam link denklemini parçalarıyla yazalım. Alıcıda taşıyıcı/gürültü-yoğunluğu oranı (C/N0), sistemin nihai kalite ölçütüdür:

```
   Downlink C/N0 (dB-Hz) =  EIRP_uydu
                          − FSPL
                          − (diğer kayıplar: atmosfer, yağmur, anten yanlış hizalama)
                          + G/T_alıcı
                          − 10·log₁₀(k)        [k = Boltzmann sabiti; −228,6 dBW/K/Hz]

   Terimler:
   EIRP   = uydunun etkin izotropik yayılan gücü (dBW) = verici gücü + anten kazancı
   FSPL   = serbest uzay yol kaybı (dB)  ← yukarıda hesaplandı (~205 dB GEO/Ku)
   G/T    = alıcının "iyilik faktörü" (dB/K) = anten kazancı − sistem gürültü sıcaklığı
   k      = Boltzmann sabiti (gürültü tabanını belirler)
```

Üç terim, üç farklı mühendisin sorumluluğudur:

- **EIRP (uydu tarafı):** Uydunun yere ne kadar güçlü "bağırdığı". EIRP = HPA gücü + uydu anten kazancı. Modern HTS (High Throughput Satellite) uyduları, gücü geniş bir alana yaymak yerine dar **nokta hüzmelerine** (spot beam) odaklayarak EIRP'i yükseltir; aynı güç küçük bir alana düştüğünde o alandaki EIRP artar. Operatör, kapsama haritasında EIRP'i dBW cinsinden footprint çizgileriyle yayımlar; bir bölgede EIRP ne kadar yüksekse, orada o kadar küçük çanak yeter.

- **FSPL + atmosfer (yol tarafı):** Mühendisin kontrol edemediği, mesafenin dayattığı kayıp. GEO için ~205 dB sabittir; yağmur bunun üstüne Ku/Ka'da onlarca dB ekleyebilir. Bu terim "doğa vergisi"dir.

- **G/T (alıcı tarafı):** Yer alıcısının kalite mührü. G (anten kazancı) çanağı büyüttükçe artar; T (sistem gürültü sıcaklığı) LNB ne kadar düşük gürültülüyse o kadar düşer. G/T'yi yükseltmenin iki yolu: daha büyük çanak veya daha iyi (soğuk) LNB. Bir alıcının "ne kadar iyi" olduğu tek bir sayıda, G/T'de özetlenir.

Sayısal sezgi için: tipik bir DTH Ku alımında uydu EIRP'i ~50 dBW mertebesinde, FSPL ~205 dB, ev çanağı G/T ~12 dB/K civarı olabilir; bu değerler 60–90 cm çanakla yeterli C/N0 verir. Aynı uyduyu daha düşük EIRP footprint kenarında almak istersen ya daha büyük çanak (G artır) ya daha iyi LNB (T düşür) gerekir. **Link bütçesi, "neden bu çanak bu kadar büyük" sorusunun matematiksel cevabıdır.**

Uplink bütçesi de simetrik biçimde hesaplanır ama orada işaretler tersine döner: yerin EIRP'i (büyük çanak + HPA ile yüksek tutulabilir) ile uydunun G/T'si arasında. Toplam bağ kalitesi, uplink ve downlink C/N0'larının birleşiminden (paralel dirençler gibi, daha zayıf olan baskın) çıkar — bent-pipe'ta her iki gürültü toplanır, bu da bent-pipe'ın kalite dezavantajının sayısal nedenidir (Bölüm 5).

Not: Link bütçesi terimleri (özellikle margin, implementation loss, pointing loss) sistemden sisteme farklı ele alınır; yukarıdaki çatı doğru ama proje düzeyinde tam değerler operatör spesifikasyonundan teyit edilmeli.

---

<a id="9"></a>
## 9. Modülasyon ve Kodlama: DVB-S/S2/S2X, QPSK/8PSK/APSK, FEC

Uydudan inen bir taşıyıcının "içinde" ne olduğu, modülasyon ve kanal kodlamasıyla belirlenir. Uydu yayın dünyasının fiili standardı DVB ailesidir; onu tanımak, spektrumdaki bir GEO taşıyıcısının türünü okumanın anahtarıdır. (Modülasyon temelleri — QPSK/8PSK/QAM, sembol/bit ilişkisi — için bkz. Bölüm 1; burada uyduya özgü pratiğe odaklanırız.)

**Neden PSK/APSK, neden çok-seviyeli QAM değil?** Uydu HPA'ları (özellikle uydu üzerindeki TWT/SSPA yükselteçler) yüksek verim için doygunluğa yakın çalışır; doygun bir yükselteç sinyalin **genliğini** sıkıştırır (nonlineerlik). Genlikte bilgi taşıyan modülasyonlar (yüksek-mertebe QAM) bu sıkışmadan çok zarar görür. Bu yüzden uydu, bilgiyi ağırlıkla **fazda** taşıyan modülasyonları sever: QPSK (2 bit/sembol), 8PSK (3 bit/sembol). Daha yüksek verim için APSK (Amplitude and Phase Shift Keying — 16APSK, 32APSK) kullanılır; APSK, genliği yalnızca birkaç halkaya böler (örn. 16APSK'da 4+12 noktalı iki/üç halka), böylece nonlineer HPA ile karasal QAM'den daha uyumludur. Mühendislik sezgisi: **uydu modülasyonu nonlineer yükselteçle barışık olmak zorundadır; bu yüzden "halka" yapılı PSK/APSK egemendir, yoğun dikdörtgen QAM nadirdir.**

```
   Uydu modülasyon ailesi (artan verim, azalan dayanıklılık):

   QPSK        8PSK            16APSK              32APSK
   (2 bit)     (3 bit)         (4 bit)             (5 bit)

     ● ●         ●              ● ● ● ●           halkalı, yoğun
              ●   ●           ●  ● ● ●  ●         daha çok bit
     ● ●     ●     ●          ● ●●●●●● ●          ama daha çok
              ●   ●           ●  ● ● ●  ●         SNR ister
                ●              ● ● ● ●
   en dayanıklı ──────────────────────────► en verimli
   (düşük SNR'da çalışır)        (yüksek SNR + iyi link ister)
```

**FEC (Forward Error Correction — İleri Hata Düzeltme):** Uydu kanalı gürültülüdür; ham bit'ler yolda bozulur. FEC, gönderilen veriye fazladan "yedek" bit ekleyerek alıcının hataları (geri iletim istemeden) düzeltmesini sağlar. FEC oranı (code rate), faydalı bit'in toplam bit'e oranıdır: 1/2 oranı (her faydalı bit için bir yedek; çok dayanıklı ama yavaş) ile 9/10 (çok az yedek; hızlı ama zayıf linkte kırılır) arasında değişir. Link kötüleştikçe (yağmur, footprint kenarı) daha düşük (güçlü) FEC'e inilir.

DVB ailesinin gelişimi, bu iki kolun (modülasyon + FEC) sürekli iyileştirilmesidir:

| Standart | Yıl (yaklaşık) | Modülasyon | FEC | Karakteristik |
|---|---|---|---|---|
| **DVB-S** | 1990'lar | QPSK | Viterbi + Reed-Solomon (concatenated) | İlk nesil; basit, dayanıklı, verimi düşük |
| **DVB-S2** | 2000'ler | QPSK/8PSK/16APSK/32APSK | LDPC + BCH | Shannon sınırına yakın; ACM (adaptif kodlama/modülasyon); ~%30 daha verimli |
| **DVB-S2X** | 2010'lar | DVB-S2 + daha ince adımlar, 64/128/256APSK'a kadar | gelişmiş LDPC | Daha ince MODCOD basamakları, daha düşük roll-off; HTS için ek verim |

DVB-S2'nin getirdiği iki kavram pratikte çok önemlidir. **MODCOD**, modülasyon ve FEC oranının birleşik bir "vites" seçimidir (örn. "8PSK 3/4" bir MODCOD'dur). **ACM (Adaptive Coding and Modulation)**, her alıcı için anlık link kalitesine göre en uygun MODCOD'u seçer: açık havada yüksek verimli vites (32APSK 9/10), yağmurda dayanıklı vitese düşüş (QPSK 1/2). Böylece sistem, en kötü alıcıya göre sabitlenmek yerine, koşullara uyum sağlar. Bu, modern Ku/Ka VSAT'ın yağmurla başa çıkma yönteminin kalbidir (Bölüm 7 ile bağlantılı).

**Sembol hızı (symbol rate):** Taşıyıcının saniyede kaç sembol gönderdiğidir ve taşıyıcının bant genişliğini doğrudan belirler (kabaca bant ≈ sembol hızı × (1 + roll-off)). Bir transponder içinde tek bir geniş taşıyıcı (yüksek sembol hızı, örn. 27,5 Msym/s) ya da çok sayıda dar taşıyıcı (her biri düşük sembol hızı) bulunabilir. Spektrumda taşıyıcının genişliği, doğrudan sembol hızını ele verir — bu, taşıyıcı tanımanın temel ipucudur (Bölüm 10).

---

<a id="10"></a>
## 10. Spektrumda Bir Uydu Taşıyıcısını Tanımak

Bir GEO çanağını alıcına bağlayıp spektrumu açtığında, karşına dünya kadar düz görünen bir manzara çıkar; ama eğitimli göz, oradaki blokları okur. Bu bölüm, yalnızca **gözlem** (RX) ile bir uydu taşıyıcısının türünü çıkarsamanın mühendislik yöntemidir — içerik çözme değil, kimlik okuma.

Bir DVB taşıyıcısının spektrumdaki imzası karakteristiktir: düz tepeli, dikdörtgenimsi, kenarları roll-off ile yuvarlatılmış bir "tepe platosu". Gürültü tabanından belirgin biçimde yükselen, içi düz (çünkü rastgele veri + FEC sinyali beyaz gürültü gibi düzleştirir), kenarları keskin bir blok görürsün:

```
   DVB taşıyıcısı — spektrum imzası:

   güç
    │        ┌──────────────────┐         ← düz tepe (rastgele veri = düz spektrum)
    │       ╱│                  │╲        ← roll-off kenarları (yumuşak iniş)
    │      ╱ │                  │ ╲
    │ ────╱  │                  │  ╲────  ← gürültü tabanı
    └────────┴──────────────────┴───────► frekans
             │◄── sembol hızı ──►│
             (bant genişliği ≈ Rs × (1+α))

   Bir transponder içinde birden çok taşıyıcı:

    │   ┌──┐  ┌────────┐   ┌──┐  ┌──┐
    │   │  │  │        │   │  │  │  │
    │ ──┘  └──┘        └───┘  └──┘  └──  ← farklı genişlikte = farklı sembol hızı
    └──────────────────────────────────► frekans
       dar    geniş        dar  dar
```

Taşıyıcıyı kimliklendirme adımları (hepsi pasif gözlem):

1. **Bant genişliğini ölç.** Tepe platosunun genişliği sembol hızını verir. ~36 MHz'lik bir blok klasik tek-taşıyıcı-per-transponder yayını; içinde birkaç dar blok ise çoklu taşıyıcı (SCPC — Single Channel Per Carrier) düzenini düşündürür.
2. **Şeklini oku.** Düz tepe + roll-off = DVB benzeri sayısal taşıyıcı. Sivri bir çizgi = dar bantlı taşıyıcı veya beacon (uydunun konum/kimlik işaretçisi; çok dar, sürekli). Gürültü gibi geniş ve düşük = yayılmış spektrum (spread spectrum) olabilir.
3. **Frekans/bandı not et.** ~11 GHz downlink + ~10,7–12,75 GHz aralığı Ku DTH'i; ~3,7–4,2 GHz C-bandı yayını; ~1,5 GHz L-bandı mobil uydu. Bant, kullanım türünü daraltır (Bölüm 7 tablosu).
4. **Polarizasyonu dene.** Uydu taşıyıcıları çoğunlukla iki ortogonal polarizasyonu (dikey/yatay ya da sağ/sol dairesel) ayrı kanal olarak kullanır; LNB polarizasyonunu çevirdiğinde bazı taşıyıcılar kaybolup başkaları görünüyorsa, frekans yeniden kullanımı (polarization reuse) görüyorsun demektir.
5. **Bir demodülatöre kilitle (yalnızca açık yayında).** SatDump/açık DVB araçları bir taşıyıcıya kilitlenip MODCOD, sembol hızı, FEC gibi katman-fiziksel parametreleri **şifre çözmeden** raporlayabilir — bu, taşıyıcının "kimlik kartı"dır; içeriği açmak ayrı (ve genelde yasal izin gerektiren) bir iştir.

Burada çizgi nettir: taşıyıcının **fiziksel parametrelerini** (genişlik, modülasyon, FEC, frekans, polarizasyon) okumak spektrum okuryazarlığıdır ve pasiftir. Abonelik/şifreli bir taşıyıcının **içeriğini** açmak (CA/CW çözme) tamamen farklı bir konudur, bu kitabın kapsamı dışındadır ve çoğu yerde yasadışıdır. Mühendislik amacımız "bu blok ne tür bir taşıyıcı" sorusudur, "bu blokta ne yazıyor" değil.

Pratikte: Açık (Free-To-Air, FTA) yayınlar şifresizdir ve bunların alımı/izlenmesi birçok ülkede serbesttir; uydu meraklılarının "FTA tarama" topluluğu bu açık yayınları kataloglar. Şifreli (abonelik) taşıyıcılar ise erişim kontrolü (Conditional Access) altındadır; onların çözümü yasal izin ister. Spektrumda ikisi de aynı düz-tepeli bloğa benzer — fark, içeriğin şifreli olup olmadığıdır, spektrum şeklinde değil.

---

<a id="11"></a>
## 11. Pasif/Yasal Dinleme — Meteoroloji Uyduları (NOAA APT, Meteor LRPT, GOES)

Bu bölüm, kitabın en somut ve en tatmin edici pratiğidir: tependen geçen bir uydudan, kendi antenin ve SDR'ınla, kendi gökyüzünün canlı hava görüntüsünü almak. Tamamen yasaldır (açık, şifresiz yayın; alım dünya çapında serbest), nispeten ucuzdur ve "vay be" anını garantiler. Meteoroloji uyduları bu iş için idealdir çünkü bilinçli olarak basit, yön bağımsız antenle alınabilecek şekilde tasarlanmıştır. (Bu uyduların protokol/çözümleme tarafı Bölüm 5'te de geçti; burada uydu-haberleşmesi perspektifinden uçtan uca ele alıyoruz.)

### 11.1 NOAA APT (137 MHz) — Klasik Başlangıç

NOAA'nın kutupsal meteoroloji uyduları (NOAA-15/18/19), APT (Automatic Picture Transmission) adlı analog benzeri bir formatta yaklaşık 137 MHz'de yayın yapar. APT, görüntüyü satır satır, 2400 Hz'lik bir alt-taşıyıcıyı genlik modüle ederek (sonra tümü FM ile) gönderen, onlarca yıllık dayanıklı bir formattır. Düşük çözünürlüklü ama almak çok kolaydır.

```
   NOAA APT alım zinciri:

   137 MHz LEO uydu (kutupsal)
        │  (geçiş ~10 dk, Doppler ±3 kHz — APT bandında yutulur)
        ▼
   [ 137 MHz anten: V-dipol / QFH / turnike ]   ← yön bağımsız, RHCP tercih
        │
        ▼
   [ (opsiyonel) 137 MHz LNA + bant filtre ]      ← zayıf sinyali kaldırır
        │
        ▼
   [ RTL-SDR (FM, ~40 kHz bant) ]                  ← Bölüm 2/4
        │  IQ veya demod ses
        ▼
   [ Yazılım: SatDump / WXtoImg-benzeri ]          ← APT çöz → satırları görüntüye diz
        │
        ▼
   Görüntü: görünür/IR kanal + telemetri çubuğu
```

İncelik: APT'nin gönderdiği iki görüntü kanalı (genelde bir görünür, bir kızılötesi) yan yana gelir; yazılım ayrıca yanlardaki telemetri/senkron çubuklarını kullanarak görüntüyü hizalar ve kalibre eder. Anten anahtar rolündedir: 137 MHz dalga boyu ~2,2 metredir, bu yüzden anten bir miktar boyludur; yön bağımsız (omnidirectional) ve dairesel polarize (uydu RHCP yayınlar, geçiş boyunca yönelim değişir) bir anten — V-dipol, turnike (turnstile) veya QFH (Quadrifilar Helix) — en iyi sonucu verir. Yagi de kullanılabilir ama uyduyu elle takip etmek gerekir.

Not: NOAA APT uydularının bir kısmı yaşlanmaktadır ve servis durumları zamanla değişir; hangi uydunun aktif olduğunu güncel kaynaktan teyit et. APT, dünya genelinde aşamalı olarak yerini sayısal formatlara bırakmaktadır.

### 11.2 Meteor-M LRPT (137 MHz) — Sayısal, Daha Net

Rusya'nın Meteor-M serisi, yine ~137 MHz'de ama **sayısal** LRPT (Low Rate Picture Transmission) formatında yayın yapar. LRPT, QPSK modülasyonlu, daha yüksek çözünürlüklü ve renkli (çok kanallı) görüntü verir — APT'nin sayısal halefi gibidir. Alım zinciri APT'ye çok benzer (aynı anten, aynı SDR), fark çözümleme yazılımındadır: LRPT, QPSK'i demodüle edip paketlerden görüntüyü yeniden kurar.

Pratik fark: LRPT sayısal olduğu için ya temiz alırsın (net görüntü) ya da paket kaybıyla bozulur (siyah satırlar) — APT'nin yumuşak, gürültülü ama "okunabilir" bozulmasının aksine sayısal bir uçurum vardır. Bu yüzden Meteor için iyi bir geçiş (yüksek elevation, temiz anten) daha kritiktir. Buna karşılık başarılı bir geçişte LRPT görüntüsü APT'den belirgin biçimde daha güzeldir.

### 11.3 GOES HRIT (GEO, ~1,7 GHz) — Sabit, Sürekli, İleri Seviye

GOES (ve benzeri GEO meteoroloji uyduları) kutupsal değil **geostationary**dir; gökyüzünde sabit durur ve sürekli, tam-disk Dünya görüntüleri yayınlar. Yayın formatı HRIT/LRIT (High/Low Rate Information Transmission) olup yaklaşık 1,7 GHz (L-bandı) civarındadır. GEO olduğu için geçiş beklemek yok — çanağını bir kez doğrultursun, sürekli görüntü akar.

```
   Meteoroloji uydusu — LEO (NOAA/Meteor) vs GEO (GOES):

   LEO (137 MHz):                    GEO (~1,7 GHz HRIT):
   - geçiş bekle (~10 dk)            - sürekli akar (sabit)
   - yön bağımsız anten              - küçük çanak/yönlü anten (sabit kilit)
   - şerit görüntü (uydu altı)       - tam-disk (yarımküre tek karede)
   - ucuz, kolay başlangıç           - daha çok donanım, ileri seviye
```

GOES alımı ileri seviyedir çünkü 1,7 GHz'de küçük bir yönlü anten (ızgara çanak veya yamalı anten dizisi) ve genelde bir LNA gerekir; ayrıca sürekli yüksek veri akışını çözmek için SatDump gibi yetkin bir yazılım şarttır. Buna karşılık ödülü büyüktür: evinden, gerçek zamanlıya yakın, profesyonel kalitede tam-disk Dünya görüntüleri. (Bazı GEO meteoroloji yayınları farklı bant/formatta olabilir; bölgesel uydu ve frekansı güncel kaynaktan teyit et.)

Üç sistemin karşılaştırması:

| Sistem | Yörünge | Frekans | Format | Anten | Zorluk |
|---|---|---|---|---|---|
| NOAA APT | LEO kutupsal | ~137 MHz | analog (APT) | V-dipol/QFH/turnike | en kolay |
| Meteor-M LRPT | LEO kutupsal | ~137 MHz | sayısal QPSK (LRPT) | V-dipol/QFH/turnike | kolay-orta |
| GOES HRIT | GEO | ~1,7 GHz (L) | sayısal (HRIT/LRIT) | küçük çanak + LNA | orta-ileri |

---

<a id="12"></a>
## 12. Inmarsat, Iridium ve Amatör Uydular (Açık Çerçevede Gözlem)

Meteoroloji uydularının ötesinde, pasif gözlemle çalışılabilen birkaç sistem daha vardır. Burada çizgi dikkatle çizilir: **ağ sinyallerinin varlığını, yapısını ve davranışını gözlemlemek** (spektrum okuryazarlığı, sistem mühendisliği) ile **özel haberleşme içeriğini çözmek** (çoğu yerde yasadışı) tamamen ayrıdır. Bu bölüm yalnızca birincisini ele alır.

### 12.1 Inmarsat (GEO, L-bandı ~1,5 GHz)

Inmarsat, GEO uydular üzerinden denizcilik, havacılık ve kara mobil hizmet veren klasik bir operatördür. L-bandında (~1,5 GHz) çalıştığı ve GEO (sabit) olduğu için, sabit küçük bir yamalı/heliks antenle downlink'i gözlemlemek görece kolaydır. Açık çerçevede gözlemlenebilen birkaç sistem bileşeni vardır:

- **STD-C / EGC (Enhanced Group Call):** Denizcilik güvenlik yayınları (NAVTEX benzeri), SafetyNET deniz uyarıları ve grup çağrıları. Bunların bir kısmı tasarımı gereği **kamuya açık güvenlik yayını** niteliğindedir (denizcilik emniyet bilgisi). Bu açık yayın bileşenlerini almak, denizcilik güvenlik bilgisinin nasıl dağıtıldığını anlamak için öğreticidir.
- **AERO:** Havacılık uydu haberleşmesi. Sistemin **varlığını ve çerçeve yapısını** spektrumda gözlemlemek mümkündür.

Önemli sınır: Inmarsat downlink'inde taşınan trafiğin bir kısmı kamuya açık güvenlik yayını olsa da, **özel/üçüncü taraf haberleşmesinin** içeriğini çözmek haberleşmenin gizliliğini ihlal eder ve yasaktır (Türkiye'de TCK 132–140; bkz. Bölüm 0 ve 6). Bu kitap STD-C/AERO'yu yalnızca **sistem mimarisi ve açık güvenlik yayını** düzeyinde, "bu ağ spektrumda nasıl görünür, nasıl çerçevelenir" sorusuyla ele alır; özel mesaj içeriğine erişim yöntemi vermez.

### 12.2 Iridium (LEO, ~1,6 GHz)

Iridium, alçak yörüngede (~780 km), çok sayıda uydudan oluşan, çapraz-linkli (Bölüm 6) küresel bir takımdır. ~1,6 GHz'de çalışır. LEO olduğu için uydular hızla geçer ve sinyalde belirgin Doppler vardır. Iridium'un ilgi çekici yanı, sürekli yayınladığı **ağ yönetim sinyalleridir** (ring alert, çağrı kanalı çerçeveleri): bunlar belirli bir kullanıcıya ait içerik değil, ağın kendini yöneten sistem mesajlarıdır ve açık çerçevede gözlemlenebilir.

Mühendislik açısından Iridium, bir LEO takımının **ağ davranışını** incelemek için iyi bir laboratuvardır: uydu geçişleri, hüzme anahtarlamaları, Doppler eğrileri, çerçeve yapısı — hepsi pasif olarak görünür. Yine sınır aynıdır: ağın **sistem/yönetim sinyallerini** gözlemlemek (sistem mühendisliği) ile **kullanıcı haberleşmesinin içeriğini** çözmek (yasadışı) birbirinden ayrıdır; bu kitap yalnızca birincisini, "takım nasıl çalışır" perspektifiyle ele alır.

### 12.3 Amatör Uydular (AMSAT, SO-50, ISS)

Amatör radyo uyduları, öğrenme için en temiz ve en meşru alandır çünkü tasarımları gereği **açıktır** ve amatör topluluğu için yayın yaparlar. Önemli ayrım: bu uydularda **dinlemek (RX) herkese açıktır**; ancak **yayın yapmak (TX)** geçerli bir amatör telsiz lisansı ister (Türkiye'de BTK sınav/lisans; bkz. Bölüm 0 ve 8).

- **AMSAT uyduları / SO-50:** Amatör "transponder" uyduları, yukarıdan (uplink) gelen amatör sinyalleri alıp aşağıya (downlink) geri yayınlayan, uzaydaki bent-pipe röleler gibidir. SO-50 popüler bir FM "papağan"dır: lisanslı amatörler 145 MHz'den girer, 436 MHz'den çıkar. **Sadece dinlemek** isteyen biri downlink'i (436 MHz) lisanssız da alabilir.
- **ISS (Uluslararası Uzay İstasyonu):** ISS, amatör bandlarda **APRS** (paket konum/mesaj, ~145,825 MHz) ve zaman zaman **SSTV** (Slow Scan TV — yavaş tarama görüntü) yayınları yapar. SSTV etkinliklerinde ISS, görüntüleri amatör FM üzerinden gönderir; **bunları almak (RX) lisans gerektirmez** ve harika bir alıştırmadır — tependen geçen uzay istasyonundan gelen bir görüntüyü çözmek. APRS dinlemek de aynı şekilde RX'tir.

```
   Amatör uydu — RX herkese açık, TX lisans ister:

         ☉ uydu/ISS (LEO, hızlı geçer, Doppler var)
        ╱ ╲
   uplink  downlink
   (TX:    (RX:
   LİSANS  HERKESE
   ister)  açık)
      │       │
      ▼       ▼
   [lisanslı  [herkes:
    amatör]    SDR + 145/436 anten]
              → SO-50 sesi, ISS SSTV görüntüsü, APRS paketi
```

İncelik (Doppler): Amatör uydular 145/436 MHz dar bantta çalıştığı için Doppler düzeltmesi pratik olarak gereklidir (Bölüm 3). Dinleme için bile, geçiş boyunca frekansı (elle veya gpredict ile otomatik) kaydırman gerekir; aksi halde dar FM/SSB sinyali bandından kaçar. Bu, Doppler'ı "kitapta okunan bir formül" olmaktan çıkarıp elinle deneyimlemenin en iyi yoludur.

---

<a id="13"></a>
## 13. SatNOGS: Dağıtık Yer İstasyonu Ağı

Tek bir yer istasyonunun temel kısıtı coğrafyadır: yalnızca senin gökyüzünden geçen uyduyu, yalnızca o geçiş süresince alabilirsin. SatNOGS (Satellite Networked Open Ground Station), bu kısıtı topluluk gücüyle aşan açık kaynaklı bir projedir: dünya genelinde gönüllülerin kurduğu yüzlerce yer istasyonunu tek bir ağda birleştirir ve gözlemleri herkese açık bir veritabanında paylaşır.

```
   SatNOGS — dağıtık alım ağı:

        uydu (LEO, herhangi bir bölgeden geçer)
          ╱    │    ╲
         ╱     │     ╲
   ┌────────┐ ┌────────┐ ┌────────┐
   │ist. TR │ │ist. EU │ │ist. US │   ← gönüllü yer istasyonları
   │(çanak/  │ │        │ │        │      (her biri SDR + anten + Pi)
   │ QFH+SDR)│ │        │ │        │
   └────┬───┘ └────┬───┘ └────┬───┘
        └──────────┼──────────┘
                   ▼
        ┌──────────────────────┐
        │   SatNOGS NETWORK     │   ← merkezi planlama + açık veritabanı
        │ - geçiş planlama (TLE)│      (herkes gözlem zamanı rezerve edebilir,
        │ - otomatik gözlem     │       sonuçlar herkese açık)
        │ - waterfall + demod   │
        └──────────────────────┘
```

SatNOGS'un işleyişi öğreticidir çünkü bu bölümün tüm parçalarını bir araya getirir:

1. **TLE ile planlama:** Ağ, her uydunun güncel TLE'sini (Two-Line Element — yörünge parametreleri) kullanarak hangi istasyonun ne zaman hangi uyduyu göreceğini hesaplar (Bölüm 2'deki geçiş geometrisi). Bir kullanıcı, herhangi bir istasyonda gelecekteki bir geçiş için gözlem rezerve edebilir.
2. **Otomatik alım:** Rezerve edilen geçişte istasyon, antenini (sabit veya motorlu) uyduya yönlendirir, SDR ile kaydı yapar, Doppler düzeltmesini uygular (Bölüm 3).
3. **Açık sonuç:** Gözlemin waterfall'ı, demodüle verisi ve (mümkünse) çözümlenmiş telemetrisi açık veritabanına yüklenir; herkes inceleyebilir.

SatNOGS'un değeri çok yönlüdür. Donanımın yoksa bile, başka bir istasyonun gözlemini izleyerek öğrenebilirsin (yörünge geçişi, Doppler eğrisi, sinyal yapısı gerçek verilerde nasıl görünür). Kendi istasyonunu kurarsan (bir Raspberry Pi + SDR + uygun anten ile), küresel bir bilim/gözlem ağına katkı yaparsın; özellikle üniversite CubeSat'ları ve amatör uydular için SatNOGS, telemetri toplamanın fiili açık altyapısıdır. Hepsi tasarımı gereği **pasif alım** ve açık veriyle çalışır; bu yüzden hem yasal hem öğreticidir.

Pratikte: SatNOGS, bu bölümün "yörünge → geçiş → Doppler → anten → SDR → demodülasyon → veri" zincirini uçtan uca, gerçek ve paylaşılan veriyle deneyimlemenin en erişilebilir yoludur. Donanım kurmadan önce birkaç gözlemi incelemek, kendi alım denemelerin için sağlam bir sezgi verir.

---

<a id="14"></a>
## 14. Donanım: Anten, LNB, Takip ve SDR Zinciri

Uydu alımının donanım zinciri, karasal zincirin (Bölüm 2–4) uyduya uyarlanmış halidir; iki ek zorluk vardır: çok zayıf sinyal (büyük yol kaybı → yüksek kazanç + düşük gürültü gerekir) ve hareket (LEO için takip ya da geniş açı gerekir). Zinciri uçtan uca görelim.

```
   Genel uydu alım zinciri (RX):

   [ ANTEN ]──[ LNB/LNA ]──[ (downconverter) ]──[ SDR ]──[ YAZILIM ]
      │            │              │                 │          │
   yörünge/    düşük gürültü   yüksek bandı      IQ örnek   demod +
   banda göre  + kazanç        SDR'ın görebil-  (Bölüm 2)  görüntü/veri
   seçilir     (ilk kademe)    diği aralığa indir          (SatDump/GNU Radio)
```

**Anten seçimi — yörünge belirler:**

| Senaryo | Uygun anten | Neden |
|---|---|---|
| LEO 137 MHz (NOAA/Meteor) | V-dipol, turnike, QFH | yön bağımsız, dairesel polar; uydu gökyüzünü tarar, takip istemez |
| LEO 145/436 MHz (amatör) | el Yagi (takip) veya turnike | dar bant + kazanç; Yagi ile elle takip ya da turnike ile sabit |
| GEO Ku/C (TV/VSAT) | parabolik çanak (sabit) | yüksek kazanç + sabit hedef; bir kez kilitle |
| GEO L (~1,5–1,7 GHz) | yama/helix veya küçük ızgara çanak | orta kazanç; sabit GEO hedefi |
| GOES HRIT (~1,7 GHz) | küçük çanak / yama dizisi + LNA | zayıf GEO sinyali; yönlü kazanç + ön yükselteç |

Anten seçiminin iki ekseni: **yön bağımsızlık vs kazanç** (LEO geniş açı ister ama kazanç düşer; GEO yüksek kazanç ister ama yönlü olmalı) ve **polarizasyon** (uydular çoğunlukla dairesel — RHCP/LHCP — polarize yayınlar; doğru dairesel polarize anten 3 dB ve üzeri kazandırır). QFH (Quadrifilar Helix), 137 MHz LEO için ideal kabul edilir çünkü hem yarımküresel geniş açı hem dairesel polarizasyon sunar.

**LNB ve LNA — ilk kademe her şeyi belirler:** Link bütçesinde (Bölüm 8) gördüğümüz G/T'nin "T" (gürültü) kısmını, alıcının **ilk** aktif kademesi domine eder (Friis gürültü formülü: ilk kademenin gürültüsü tüm zincire baskındır). Bu yüzden antene mümkün olduğunca yakın, düşük gürültülü bir ön kademe konur:

- **LNB (Low Noise Block downconverter):** Ku/C çanak uydu alımında standart. Hem düşük gürültüyle yükseltir (LNA) hem de yüksek uydu bandını (örn. 11 GHz) SDR'ın/uydu alıcısının görebileceği ara frekansa (örn. ~1 GHz L-band IF) **indirir** (downconvert). Çanağın odağındaki o "boynuz", LNB'dir.
- **LNA (Low Noise Amplifier):** 137 MHz / 1,7 GHz gibi SDR'ın doğrudan görebildiği bandlarda ayrı bir downconverter gerekmez; sadece düşük gürültülü bir yükselteç (ve bant filtresi) zayıf sinyali kaldırmaya yeter. Antene yakın monte edilmeli (kablo kaybından önce).

**Downconverter:** Eğer uydu bandı SDR'ın frekans aralığının üstündeyse (örn. Ku 11 GHz, RTL-SDR ~1,7 GHz'e kadar görür), araya bir downconverter (veya LNB'nin kendisi) girerek bandı SDR'ın görebileceği aralığa taşır. GOES ve L-band gibi durumlarda SDR doğrudan görebildiği için downconverter gerekmeyebilir.

**Takip: motorlu vs sabit:**
- **Sabit + yön bağımsız (LEO 137 MHz):** QFH/turnike geniş açıyı kapsar; uyduyu izlemeye gerek yok, geçiş boyunca otomatik alır. En basit, en yaygın başlangıç.
- **Motorlu takip (LEO dar bant / yüksek kazanç):** Yagi/çanak gibi yönlü anten, geçiş boyunca uyduyu azimut+yükseklikte izlemek için motora (rotator) bağlanır; gpredict gibi yazılım TLE'den hesaplayıp rotatoru sürer. Daha çok kazanç ama daha çok karmaşıklık.
- **Sabit kilit (GEO):** GEO sabit olduğu için çanak bir kez doğrultulur, motor gerekmez. En kararlı senaryo.

**SDR ve yazılım:** SDR tarafı Bölüm 2'de detaylı; uyduya özgü tek not, yeterli bant genişliği (LRPT/HRIT için birkaç MHz) ve kararlı saat (Doppler düzeltmesi ve sayısal demod için TCXO'lu RTL-SDR tercih) gerektiğidir. Yazılım tarafında modern fiili standart **SatDump** (NOAA/Meteor/GOES ve birçok uydu için uçtan uca: alım → demod → görüntü) ve esnek akışlar için **GNU Radio**'dur (Bölüm 4). gpredict + rigctld/rotctld ikilisi, Doppler ve rotator otomasyonunu sağlar.

---

<a id="15"></a>
## 15. Güvenlik (Savunma Perspektifi): Neden Bazı SATCOM Zayıftır

Bu bölüm ve devamı, bilinçli olarak **savunma ve tespit** perspektifindedir. Amaç, bir sistemin neden savunmasız olabileceğini ve bunun nasıl **tespit edilip korunduğunu** anlamaktır — yetkisiz erişimin adımlarını vermek değil. Bölüm 6'daki genel RF güvenlik mantığı (manipüle edilebilir/edilemez sinyaller, kimlik doğrulamanın rolü) burada uydu bağlamına taşınır.

Uydu güvenliğinin temel sorunu, mimarinin doğasında yatar: uydu sinyali, tanımı gereği **çok geniş bir alana yayılır** (GEO bir kıtayı, LEO bir bölgeyi aydınlatır) ve bu alandaki herkes downlink'i fiziksel olarak alabilir. Karasal radyoda bir vericiye yaklaşmak gerekir; uyduda verici gökyüzündedir ve herkese eşit "görünür". Bu, gizliliğin **fiziksel mesafeyle değil, yalnızca şifrelemeyle** sağlanabileceği anlamına gelir. Şifreleme yoksa, downlink açıktır.

Tarihsel olarak birçok SATCOM sistemi şu nedenlerle zayıf tasarlanmıştır:

**1. Şifresiz bent-pipe transponder (Bölüm 5).** Klasik şeffaf transponder içeriğe kayıtsızdır; üzerinden geçen veri şifreliyse şifreli, açıksa açık iner. Eski yayın ve veri sistemlerinin önemli bir kısmı, downlink içeriğini şifrelemeden gönderir (çünkü maliyet, eski standart, ya da "zaten yörüngede, kim dinler" varsayımı). Sonuç: o downlink'i alan herkes, içeriği ham haliyle görür. Bu, gizliliğin mimaride değil "kimse bakmıyor" varsayımında aranmasının (güvenlik-belirsizlikle, security through obscurity) klasik örneğidir ve zayıftır.

**2. Kimlik doğrulaması olmayan uplink (eski sistemler).** Bazı eski sistemlerde uydu, kendisine gelen bir uplink taşıyıcısının **yetkili olup olmadığını kriptografik olarak doğrulayamaz**; bent-pipe ise zaten doğrulamaz (içeriğe bakmaz). Erişim kontrolü idari düzeydedir (kim hangi transponderi kiralamış). Bu, bir sonraki bölümdeki "satellite piracy" olgusunun teknik kökenidir.

**3. Eski/zayıf veya hiç şifreleme (TT&C dahil).** Bazı eski uyduların kontrol bağı (TT&C, Bölüm 6) modern kimlik doğrulama/şifreleme standartlarının gerisindedir. Payload zayıflığı içeriği açar; TT&C zayıflığı çok daha ciddidir (Bölüm 17).

**4. İnternete açık yer terminalleri.** Uydunun kendisi değil, yerdeki VSAT/SATCOM terminalleri ve yönetim arayüzleri, karasal ağ güvenliği zafiyetleri taşıyabilir (Bölüm 17).

Mühendislik sezgisi: **Uydu sinyalinin geniş alana yayılması bir "özellik" (kapsama) ama aynı zamanda bir "zafiyet" (herkes alır) olduğu için, gizlilik tek bir şeye indirgenir: uçtan uca güçlü şifreleme.** Bunu sağlamayan her sistem, yayıldığı alandaki herkese karşı açıktır. Savunma da buradan çıkar: downlink içeriğini ve özellikle TT&C'yi modern, kimlik doğrulamalı şifrelemeyle korumak; "kimse dinlemez" varsayımına asla güvenmemek.

```
   Neden uydu downlink'i "herkese açık":

        ☉ uydu (downlink TÜM footprint'e iner)
       ╱│╲
      ╱ │ ╲   geniş hüzme
     ╱  │  ╲
   ──┴──┴──┴──  ← bu alandaki HERKES sinyali alır
   yetkili  meraklı  rakip

   ► Fiziksel mesafe gizlilik sağlamaz.
   ► Tek savunma: içeriği şifrele (uçtan uca).
   ► Şifresiz downlink = footprint'teki herkese açık.
```

Not: "Eski sistemler zayıf" genellemesi modern sistemleri kapsamaz; günümüz operatörleri downlink ve TT&C şifrelemesini ciddiye alır. Buradaki amaç, **neden** bir sistemin zayıf olabileceğinin mekanizmasını göstermek ve modern savunmanın bu mekanizmaları nasıl kapattığını anlamaktır.

---

<a id="16"></a>
## 16. "Satellite Piracy" Olgusu: Mekanizma, Yasadışılık, Tespit

"Satellite piracy" (uydu korsanlığı), bir uydunun transponder kapasitesinin **yetkisiz** kullanımını anlatan bir olgudur. Bu bölüm olguyu **kavramsal olarak** açıklar: neden teknik olarak mümkün olabildiği, neden kesinlikle yasadışı olduğu ve operatörlerin bunu nasıl **tespit edip durdurduğu**. Burada hiçbir yöntem, adım veya reçete yoktur; amaç savunma ve tespit tarafını anlamaktır.

### 16.1 Neden Teknik Olarak Mümkün Olabilir

Mekanizma, Bölüm 5'teki bent-pipe transponderin doğasından çıkar. Şeffaf transponder içeriğe ve kaynağa kayıtsızdır: kendisine ulaşan, doğru frekansta ve yeterli güçte herhangi bir uplink taşıyıcısını alır, çevirir ve geri yayınlar. Uydu "bu taşıyıcı yetkili mi" diye **soramaz** çünkü kriptografik kimlik doğrulaması yapacak işlem yükü (regenerative değilse) yoktur. Erişim kontrolü uydunun kendisinde değil, **yer tarafındaki idari düzendedir** (kim hangi transponderi/kapasiteyi kiralamış).

```
   Bent-pipe'ın yetkisiz kullanıma açıklığı (kavramsal):

   transponder = içeriğe/kaynağa kayıtsız boru
        │
        ├─ yetkili uplink ──▶ geri yayınlanır ✓ (kiralanmış)
        │
        └─ yetkisiz uplink ─▶ uydu AYIRT EDEMEZ ──▶ yine geri yayınlanır ✗
              (kriptografik kimlik doğrulama yoksa)

   ► Açık, uydunun "ayırt edememesi"dir.
   ► Engel idaridir (kiralama/lisans), teknik kimlik doğrulama değil.
   ► Bu yüzden tespit + hukuk, savunmanın merkezidir.
```

Bu, transponderin bir tasarım sadeleştirmesinin (içeriğe kayıtsızlık, esneklik) yan etkisidir. Önemle: **bunun "mümkün olması" onu meşru kılmaz.** Aşağıda göreceğimiz gibi, hem fiilin kendisi hem de gerektirdiği lisanssız uplink ayrı ayrı suçtur.

### 16.2 Neden Kesinlikle Yasadışı

Yetkisiz transponder kullanımı en az iki bağımsız hukuki ihlali bir araya getirir:

1. **Lisanssız uplink (yetkisiz yayın/TX).** Uyduya sinyal basmak, bir uplink vericisiyle **yayın yapmak** demektir. Yetkisiz/lisanssız yayın, hemen her ülkede telsiz mevzuatının ağır ihlalidir (Türkiye'de BTK/telsiz mevzuatı; bkz. Bölüm 0 ve 8). Bölüm 0'ın altın kuralı: "Alıcı çoğu yerde serbesttir; verici her yerde sorumluluktur." Uplink, tanımı gereği vericidir.

2. **Hırsızlık / yetkisiz erişim.** Kiralanmış, sahibi olan bir kapasiteyi izinsiz kullanmak, mülkiyet/hizmet hırsızlığıdır — başkasının ödediği ve sahip olduğu bir kaynağı çalmaktır. Buna ek olarak, o kapasiteyi meşru kullanan tarafın yayınına **girişim** (interference) yaratır; bu da ayrı bir zarar ve suçtur.

Yani "satellite piracy" hem yayın suçu hem hırsızlık hem de girişim üretir. Tek bir fiilde birden çok ağır ihlal birleşir. Bu kitabın çizgisi gereği, bunun **nasıl yapılacağına** dair hiçbir bilgi verilmez; konu yalnızca olgunun varlığı, mekanizması ve karşı tarafın (operatör) savunması düzeyindedir.

### 16.3 Nasıl Tespit Edilir (Operatörün Savunması)

İşte savunma perspektifinin merkezi: operatörler yetkisiz taşıyıcıları **tespit eder ve yerini bulur**. Başlıca yöntemler:

**Carrier monitoring (taşıyıcı izleme).** Operatör, kendi transponderlerini sürekli izler ve "spektral envanter" tutar: hangi frekansta, hangi sembol hızında, hangi modülasyonda, hangi güçte yetkili taşıyıcı olması gerektiğini bilir. Beklenmeyen bir taşıyıcı (kiralanmamış bir frekansta beliren, beklenmeyen parametrelerde bir blok) hemen "anomali" olarak işaretlenir. Bu, Bölüm 10'da öğrendiğimiz taşıyıcı tanıma yeteneğinin operatör tarafındaki sürekli, otomatik halidir.

```
   Carrier monitoring — beklenen vs gerçek:

   BEKLENEN (kiralama kaydı):
   │ ┌──┐    ┌────┐      ┌──┐
   │ │A │    │ B  │      │C │            ← yetkili taşıyıcılar
   └─┴──┴────┴────┴──────┴──┴──► f

   GERÇEK (canlı spektrum):
   │ ┌──┐    ┌────┐ ┌─┐  ┌──┐
   │ │A │    │ B  │ │?│  │C │            ← "?" beklenmeyen = anomali
   └─┴──┴────┴────┴─┴─┴──┴──┴──► f
                       ▲
                  yetkisiz taşıyıcı tespit edildi
```

**Geolocation (kaçak taşıyıcının yerini bulma).** Bir anomali taşıyıcısı tespit edilince, operatör onun **yeryüzündeki kaynağını** bulmak ister. Bunun klasik yöntemi iki uydu üzerinden konumlandırmadır (dual-satellite geolocation): yetkisiz uplink, hedef uyduya ek olarak yakındaki ikinci (komşu) bir uyduya da bir miktar "sızar" (anten yan-lobları nedeniyle). Aynı taşıyıcının iki uyduya varış zamanı farkı (TDOA — Time Difference of Arrival) ve frekans farkı (FDOA — Frequency Difference of Arrival) ölçülür; bu iki ölçüm, yeryüzünde kesişen eğriler çizer ve kaynağı coğrafi olarak daraltır.

```
   Dual-satellite geolocation (TDOA/FDOA) — kavramsal:

        ☉ hedef uydu      ☉ komşu uydu
         \                /
          \  yetkisiz    /  (yan-lob sızıntısı)
           \ uplink     /
            \          /
             ●────────●  ← yer (kaçak verici)
                ▲
   iki uyduya VARIŞ ZAMANI farkı (TDOA) → bir eğri
   iki uyduya FREKANS farkı (FDOA)      → ikinci eğri
   kesişim → vericinin coğrafi konumu (daraltılır)
```

TDOA bir hiperbolik eğri, FDOA başka bir eğri verir; ikisinin kesişimi kaynağı bir bölgeye indirir. Operatör böylece kaçak vericinin nereden yayın yaptığını bulur, ilgili otoriteye/regülatöre bildirir ve hukuki süreç başlar. Bu, Bölüm 6 ve 7'deki konumlandırma (DF/TDOA) mantığının uydu ölçeğine taşınmış halidir.

Mühendislik sezgisi: **Yetkisiz uplink, tanımı gereği spektrumda iz bırakır (yayın yapar) ve bu iz, hem anomali olarak görülür hem de iki uydu üzerinden konumlandırılır.** "Korsan" pasif değildir; aktif olmak zorundadır (yayın yapmadan transponder kullanılamaz) ve aktiflik tespiti kaçınılmaz kılar. Savunmanın gücü buradadır: operatörün sürekli izleme + geolocation yeteneği, yetkisiz kullanımı hem görünür hem cezalandırılabilir kılar. Modern karşı önlemler ayrıca uplink kimlik doğrulaması, taşıyıcı kimliklendirme (carrier ID — DVB-S2'de taşıyıcıya gömülü kimlik bilgisi) ve regeneratif/akıllı payload ile bu açığı daha en baştan kapatma yönündedir.

Pratikte: Bu bölüm, "korsanlığın nasıl yapıldığını" değil, **operatörün savunmasının nasıl işlediğini** öğretir. Bir savunmacı/araştırmacı için değerli olan, tespit ve konumlandırma zincirini anlamaktır; saldırı tarafının adımları bu kitapta bilinçli olarak yoktur.

---

<a id="17"></a>
## 17. VSAT/Terminal Saldırı Yüzeyi ve TT&C Güvenliği

Uydunun kendisi yörüngede korunaklı olsa bile, sistemin zayıf halkası çoğu zaman **yerdedir**: VSAT/SATCOM terminalleri, yönetim arayüzleri ve TT&C altyapısı. Bu bölüm, araştırma topluluğunun yıllar içinde kavramsal olarak gösterdiği zafiyet yüzeylerini ve **savunmasını** özetler — yine reçete değil, farkındalık ve korunma düzeyinde.

### 17.1 VSAT/SATCOM Terminallerinin Saldırı Yüzeyi (Kavramsal)

VSAT terminalleri ve uydu modemleri, özünde birer ağ cihazıdır: bir uydu modemi + yönetim yazılımı + (çoğu zaman) bir web/yönetim arayüzü. Bu, onları karasal ağ cihazlarıyla aynı zafiyet sınıflarına açar:

- **İnternete açık yönetim arayüzleri.** Bazı terminallerin yönetim/telemetri arayüzleri, yanlış yapılandırma sonucu doğrudan internetten erişilebilir hale gelebilir. İnternet tarama servisleri (Bölüm 6'da değinilen türden cihaz arama motorları), bu açık arayüzlerin kavramsal varlığını ortaya koymuştur.
- **Zayıf/varsayılan kimlik bilgileri.** Saha cihazlarında varsayılan parolaların değiştirilmemesi klasik bir zafiyettir; uydu terminalleri de bundan muaf değildir.
- **Eski/yamasız yazılım (firmware).** Uzak/zor erişilen sahalardaki (gemi, petrol platformu, uzak şube) terminaller yıllarca güncellenmeyebilir; bilinen zafiyetler açık kalır.
- **Zayıf protokol güvenliği.** Bazı eski terminal yönetim protokolleri kimlik doğrulamasız veya şifrelemesizdir.

Önemle: Bu başlık **kavramsaldır** ve araştırma bulgularının özeti niteliğindedir; belirli cihaz, belirli açık veya istismar adımı **verilmez**. Amaç, bir savunmacının "uydu terminali de bir ağ cihazıdır, karasal güvenlik hijyeni gerektirir" sezgisini kazanmasıdır.

### 17.2 Savunma — Yer Terminali Tarafı

Yer terminali güvenliği, büyük ölçüde iyi bilinen ağ güvenliği hijyeninin uyguludur (Bölüm 6 ile uyumlu):

| Zafiyet | Savunma |
|---|---|
| İnternete açık yönetim arayüzü | Yönetimi internetten ayır (out-of-band/VPN), maruz yüzeyi kapat |
| Varsayılan/zayıf parola | Güçlü, benzersiz kimlik bilgileri; mümkünse çok faktörlü |
| Eski firmware | Düzenli güncelleme/yama; uzak saha için merkezi yönetim |
| Şifresiz yönetim protokolü | Şifreli/kimlik doğrulamalı yönetim kanalı |
| Şifresiz kullanıcı trafiği | Uçtan uca şifreleme (uydu linkine güvenme, üstüne VPN) |
| İzlenmeyen anomali | Terminal davranış izleme + merkezi log/SIEM |

Mühendislik sezgisi: **Uydu linki güvenli sayılmamalıdır; üstüne her zaman uçtan uca şifreleme konmalıdır.** Link şifreli olsa bile, terminalin kendisi (yönetim, firmware) karasal bir cihaz gibi sertleştirilmelidir. "Uzak ve erişilmesi zor" olması, güvenli olduğu anlamına gelmez; tersine, ihmal edilmeye en açık cihazlardır.

### 17.3 TT&C Güvenliğinin Kritikliği

Tüm sistemin en kritik güvenlik sınırı, Bölüm 6'da tanıttığımız TT&C (Telemetri, Takip, Komuta) bağıdır. Mantık katmanlıdır:

```
   Güvenlik etkisi — payload vs TT&C:

   PAYLOAD'a yetkisiz erişim (kullanıcı trafiği):
   ► içerik açığa çıkar / kapasite çalınır
   ► ciddi ama uydunun KENDİSİ sağlam kalır

   TT&C'ye yetkisiz erişim (komuta bağı):
   ► uydunun YÖNELİMİ, transponderleri, MOD'u, hatta YÖRÜNGESİ etkilenebilir
   ► uydunun kendisi kaybedilebilir  ← FELAKET
   ►►► bu yüzden TT&C, payload'dan KAT KAT daha sıkı korunmalı
```

TT&C'yi ele geçirmek, uydunun "yönetim konsolunu" ele geçirmektir: yönelimi bozarak antenleri hedeften kaçırma, transponderleri kapatma, yörünge düzeltme yakıtını boşa harcatma, hatta uyduyu kullanılamaz kılma potansiyeli. Bu yüzden TT&C güvenliği bir uydu programının **en üst** güvenlik önceliğidir ve şunları gerektirir:

- **Komut kimlik doğrulaması:** Uydu, yalnızca kriptografik olarak doğrulanmış (yetkili kontrol merkezinden geldiği kanıtlanmış) komutları yürütmeli; sahte/tekrar (replay) komutlarını reddetmeli.
- **TT&C şifreleme:** Telemetri ve komut bağı şifreli olmalı (içerik ve komut yapısı gizli).
- **Replay/jamming dayanıklılığı:** Komut bağı, yeniden oynatma saldırılarına ve karıştırmaya karşı dirençli tasarlanmalı (sayaç/nonce, zaman damgası, spread-spectrum vb.).
- **Erişim sınırlama:** TT&C, payload'dan ayrı, sınırlı erişimli kanallarda ve ayrı altyapıda çalışmalı.

Tarihsel sorun şudur: bazı eski uydular, TT&C'yi bugünkü standartların gerisinde (zayıf veya hiç şifreleme/kimlik doğrulama) tasarladı; çünkü tasarlandıkları dönemde tehdit modeli farklıydı ("kim uyduya komut gönderebilir ki" varsayımı). Modern uydu güvenliği bu varsayımı reddeder: TT&C, sistemin can damarı olarak en güçlü kriptografiyle korunmalıdır. Savunma perspektifinin özeti: **payload zafiyeti veri kaybettirir; TT&C zafiyeti uydu kaybettirir — bu yüzden güvenlik bütçesi orantısız biçimde TT&C'ye ayrılmalıdır.**

Not: Modern uydu güvenliği aktif bir mühendislik ve araştırma alanıdır; standartlar ve uygulamalar gelişmektedir. Buradaki çerçeve kavramsaldır; belirli bir uydunun/operatörün uygulaması ayrıca değerlendirilmelidir.

---

<a id="18"></a>
## 18. Alıştırmalar (Yasal)

Aşağıdaki alıştırmaların tamamı **pasif alımdır (RX)**, açık/yayın sinyaller veya kendi cihazların üzerinedir ve birçok ülkede serbesttir; yine de kendi mevzuatını teyit et (Bölüm 0). Hiçbiri uplink/yayın (TX) içermez. Sırasıyla artan zorlukta düzenlenmiştir.

**Alıştırma 1 — İlk meteoroloji görüntüsü (NOAA APT).** En kolay başlangıç. Bir RTL-SDR + basit bir 137 MHz anten (V-dipol yeter; QFH idealdir) ile bir NOAA APT geçişini al. Geçiş tahmini için bir uydu geçiş takip aracı kullan (yüksek elevation'lı bir geçiş seç). SatDump ile sinyali çöz ve görüntüyü üret. Hedef: kendi gökyüzünün canlı hava görüntüsünü almak ve görünür/IR kanalları ayırt etmek.

**Alıştırma 2 — Doppler'ı gözlemle.** NOAA/Meteor geçişi boyunca sinyalin merkez frekansının nasıl kaydığını (doğuşta yukarı, zenitte sıfır, batışta aşağı; Bölüm 3) waterfall'da izle. Mümkünse gpredict'in Doppler tahminini açıp gerçek kayma ile karşılaştır. Hedef: Doppler S-eğrisini formülden değil gözle deneyimlemek.

**Alıştırma 3 — Meteor-M LRPT (sayısal).** Aynı 137 MHz anten/SDR ile bu kez sayısal LRPT al. APT'nin yumuşak bozulması ile LRPT'nin sayısal "ya net ya bozuk" davranışını karşılaştır. Hedef: analog vs sayısal uydu yayını farkını uçtan uca görmek.

**Alıştırma 4 — Bir GEO TV taşıyıcısını TANI (yalnızca RX, içerik çözme yok).** Bir uydu çanağı + LNB + uygun bir alıcı/SDR ile bir Ku-band downlink'i spektrumda incele. Bir DVB taşıyıcısının düz-tepeli imzasını bul (Bölüm 10), bant genişliğinden sembol hızını kestir, polarizasyon değiştirince hangi taşıyıcıların değiştiğini gözlemle. **Yalnızca açık (FTA) yayında**, bir DVB aracıyla taşıyıcının MODCOD/sembol hızı/FEC parametrelerini oku. Hedef: taşıyıcı kimliklendirme (içerik değil, fiziksel parametre okuryazarlığı). Şifreli içeriği çözmeye **çalışma**.

**Alıştırma 5 — ISS'ten SSTV/APRS al.** Bir ISS SSTV etkinliği sırasında (~145,8 MHz, FM) geçişi izle ve SSTV görüntüsünü çöz; ya da ISS APRS paketini (~145,825 MHz) yakala. Sadece dinleme (RX) lisans gerektirmez. Hedef: tependen geçen uzay istasyonundan gelen bir görüntü/paket almak ve LEO Doppler'ı dar bantta (düzeltme gerektirir) deneyimlemek.

**Alıştırma 6 — SatNOGS gözlemi planla/incele.** Donanımın olsun olmasın, SatNOGS ağında bir uydu için gözlem incele (waterfall, demod, telemetri) ve mümkünse bir geçiş için gözlem rezerve et. Hedef: yörünge → geçiş → Doppler → demod zincirini gerçek, paylaşılan veride görmek; bir dağıtık alım ağının nasıl çalıştığını anlamak.

**Alıştırma 7 (ileri) — GOES HRIT tam-disk.** Küçük bir çanak/yama anten + LNA + SatDump ile bir GEO meteoroloji uydusunun (GOES tipi) HRIT yayınını al ve tam-disk Dünya görüntüsü üret. GEO sabit olduğu için çanağı bir kez doğrult. Hedef: LEO (geçiş) ile GEO (sürekli) alım farkını uçtan uca yaşamak ve profesyonel kalitede görüntü elde etmek.

**Alıştırma 8 (gözlem/savunma) — Carrier monitoring sezgisi.** Bir GEO transponderini spektrumda incele ve "beklenen taşıyıcılar" ile "gerçek spektrum" zihinsel karşılaştırmasını yap (Bölüm 16). Operatörün anomali tespitini nasıl yapacağını kavramsal olarak modelle. Hedef: yalnızca gözlemle, savunma tarafının (carrier monitoring) mantığını içselleştirmek. (Bu bir gözlem alıştırmasıdır; hiçbir uplink/müdahale içermez.)

Pratikte: 1, 2, 5 alıştırmaları ~40$'lık bir RTL-SDR + basit antenle yapılabilir ve bu bölümün özünü (yörünge, geçiş, Doppler, demodülasyon, görüntü) eksiksiz deneyimletir. 4, 7 ileri donanım ister; 6, 8 donanımsız da öğreticidir.

---

<a id="19"></a>
## 19. Hızlı Referans ve Diğer Bölümler

### Yörünge sınıfları — özet

| Sınıf | Yükseklik | Gecikme (tek yön) | Görünürlük | Anten | Tipik kullanım |
|---|---|---|---|---|---|
| LEO | ~300–2.000 km | ~2–10 ms | dakikalar (geçer) | yön bağımsız/takip | meteoroloji, geniş bant takım, ISS |
| MEO | ~20.000 km | ~50–70 ms | saatler | orta | GNSS, bazı mobil |
| GEO | 35.786 km | ~120 ms | sürekli (sabit) | çanak (sabit) | TV, VSAT, GEO meteoroloji |
| HEO | apoje ~40.000 km+ | değişken | apojede uzun | takip | yüksek enlem süreklilik (Molniya) |

### Frekans bantları — özet

| Bant | ~Frekans | Kullanım | Yağmur etkisi |
|---|---|---|---|
| L | 1–2 GHz | mobil uydu, GNSS, GEO meteoroloji (HRIT) | ihmal |
| S | 2–4 GHz | TT&C, mobil, uzay araştırma | çok az |
| C | 4–8 GHz | klasik TV/veri, tropik | düşük |
| X | 8–12 GHz | askeri/devlet SATCOM | düşük-orta |
| Ku | 12–18 GHz | DTH TV, VSAT | orta |
| Ka | 26,5–40 GHz | yüksek hız geniş bant (HTS) | yüksek |

### Anahtar formüller — özet

```
  Yörünge periyodu:   T = 2π·√(a³/μ)        μ ≈ 3,986e14 m³/s²,  a = R_dünya + h
  Doppler kayması:    Δf ≈ f₀·(v_radyal/c)  (LEO'da büyük, GEO'da ~0)
  Serbest uzay kaybı: FSPL(dB) = 20log₁₀(d_km) + 20log₁₀(f_GHz) + 92,45
                      → GEO/Ku ≈ 205 dB
  Downlink kalitesi:  C/N0 = EIRP − FSPL − kayıplar + G/T − 10log₁₀(k)
  Taşıyıcı bant gen.: BW ≈ sembol_hızı × (1 + roll-off)
```

### Pasif/yasal alım — hızlı eşleştirme

| Ne almak istiyorsun | Frekans | Anten | Yazılım | Zorluk |
|---|---|---|---|---|
| NOAA hava görüntüsü | ~137 MHz | V-dipol/QFH | SatDump | kolay |
| Meteor sayısal görüntü | ~137 MHz | V-dipol/QFH | SatDump | kolay-orta |
| ISS SSTV/APRS | ~145,8 MHz | turnike/Yagi | SatDump/SSTV/APRS dec. | kolay |
| Amatör uydu (SO-50, RX) | 436 MHz | Yagi (takip) | gpredict + SDR | orta |
| GEO TV taşıyıcı tanıma | Ku ~11 GHz | çanak + LNB | DVB/SatDump | orta |
| GOES tam-disk | ~1,7 GHz | çanak + LNA | SatDump | ileri |

### Güvenlik (savunma) — özet

| Konu | Zayıflık nedeni | Savunma / Tespit |
|---|---|---|
| Şifresiz bent-pipe downlink | içeriğe kayıtsız, footprint'teki herkes alır | uçtan uca şifreleme |
| Yetkisiz uplink (piracy) | bent-pipe kimlik doğrulamaz; idari erişim kontrolü | carrier monitoring + dual-sat geolocation (TDOA/FDOA) + carrier ID |
| VSAT terminal yüzeyi | internete açık/varsayılan parola/eski firmware | ağ hijyeni, yama, out-of-band yönetim, uçtan uca şifreleme |
| TT&C zayıflığı | eski sistemlerde zayıf/yok şifreleme | komut kimlik doğrulama + TT&C şifreleme + replay/jam dayanıklılık + erişim ayrımı |

> **Bölüm özeti.** Uydu haberleşmesi, karasal radyonun aynı fiziğini mesafe (yol kaybı ~205 dB GEO), hız (LEO Doppler ve hızlı geçiş) ve görüş alanı (tek GEO bir kıta) eksenlerinde uca taşır; yörünge sınıfı (LEO/MEO/GEO/HEO) gecikmeyi, kapsamayı ve anteni birlikte belirler, transponder mimarisi (bent-pipe vs regenerative) hem performansı hem güvenliği şekillendirir. Pasif ve yasal pratik zengindir: NOAA APT/Meteor LRPT/GOES HRIT'ten kendi elinle hava görüntüsü, ISS SSTV/APRS, açık çerçevede ağ gözlemi ve SatNOGS — hepsi RX'tir ve yörünge → geçiş → Doppler → demod → görüntü zincirini somutlaştırır. Güvenlik tarafı savunma ve tespit perspektifindedir: şifresiz bent-pipe'ın neden açık olduğu, "satellite piracy"nin neden mümkün ama kesinlikle yasadışı (lisanssız uplink + hırsızlık + girişim) olduğu ve operatörün bunu carrier monitoring + dual-satellite geolocation (TDOA/FDOA) ile nasıl tespit ettiği; en kritik sınır ise TT&C güvenliğidir (payload zafiyeti veri, TT&C zafiyeti uydu kaybettirir) — bu kitap yetkisiz erişim reçetesi vermez, yalnızca mekanizma + savunma + tespit öğretir.

---

Bu bölüm, Kanije Kalesi SIGINT El Kitabı'nın parçasıdır. Tüm bölümler ve önerilen okuma sırası için indekse bakın: [SIGINT_00 — Başlangıç ve İndeks](SIGINT_00_BASLANGIC_INDEX_VE_YASAL.md).

Doğrudan ilgili bölümler:
- [SIGINT_01 — RF Fiziği ve Modülasyon](SIGINT_01_TEMELLER_RF_VE_MODULASYON.md): link bütçesi, FSPL ve DVB modülasyonunun fizik temeli.
- [SIGINT_03 — Antenler, Donanım ve Devre Tasarımı](SIGINT_03_ANTEN_DONANIM_VE_DEVRE_TASARIMI.md): QFH/turnike/çanak, LNB, polarizasyon.
- [SIGINT_10 — GNSS/GPS Sistemleri](SIGINT_10_GNSS_GPS_SISTEMLERI.md): uydu zamanlama ve düşük güçlü uydu sinyali.
- [SIGINT_32 — Uydu-IoT ve Mega-Konstelasyonlar](SIGINT_32_UYDU_IOT_VE_MEGA_KONSTELASYON.md): bu bölümün LEO/mega-konstelasyon ileri uzantısı.
- [SIGINT_31 — SCADA, Endüstriyel Kontrol ve Telemetri RF'i](SIGINT_31_SCADA_ENDUSTRIYEL_RF.md): uydu SCADA/VSAT bağlarının güvenliği.
