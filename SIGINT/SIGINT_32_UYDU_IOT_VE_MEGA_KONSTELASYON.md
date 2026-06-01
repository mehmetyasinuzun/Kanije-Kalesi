# SIGINT EL KİTABI — BÖLÜM 32: UYDU-IoT VE MEGA-KONSTELASYONLAR

## Starlink, Iridium ve Doğrudan-Cihaz (D2D): LEO Devrimi, Pasif Gözlem ve Savunma

> Amaç: Bölüm 11 uydu haberleşmesinin genel iskeletini kurdu — yörünge sınıfları, Doppler, transponder mimarisi, link bütçesi, DVB ailesi, GEO yayın uydularının spektrum imzası. Bölüm 22 mega-konstelasyon olgusuna ilk girişi yaptı. Bu bölüm o ikisinin üstüne **ileri alarak** inşa edilir: tekrar etmez, derinleşir. Konumuz, son on yılda uzay haberleşmesini kökten değiştiren iki akımdır. Birincisi **LEO mega-konstelasyonlar** (Starlink, OneWeb, Amazon Leo): binlerce alçak yörünge uydusu, faz-dizili kullanıcı terminalleri, uydular-arası lazer bağları ve karasal fiberle yarışan gecikme. İkincisi **uydu-IoT ve doğrudan-cihaz (Direct-to-Device, D2D)**: artık sıradan bir IoT sensörünün ya da değiştirilmemiş bir cep telefonunun doğrudan uyduya bağlanması. Hedefimiz, Bölüm 11'deki gibi operatör reçetesi değil mühendislik sezgisidir: bir spektrumda Starlink kullanıcı downlink'ini ya da Iridium ağ patlamasını gördüğünde arkasındaki yörünge dinamiğini, Doppler eğrisini, handover ritmini ve saldırı yüzeyini zihninde canlandırabilmen.

> Yasal çerçeve: Bu bölüm de serinin geri kalanı gibi anlama, savunma ve spektrum okuryazarlığı amaçlıdır. Anlatılan her şey **pasif alımdır (RX)**: kullanıcı downlink'ini spektrumda tanıma, açık ağ katmanı sinyallerini gözleme, bir LEO uydusunun geçişini Doppler kaymasıyla izleme. Hiçbir yetkisiz uplink, terminal aktivasyonu, abonelik/şifreli içeriğin çözülmesi veya kullanıcı trafiğine erişim önerilmez ya da tarif edilmez. Iridium/Starlink gibi sistemlerin **kullanıcı içeriği şifrelidir** ve bu kitap onu açmaz; yalnızca açık ağ katmanının (senkronizasyon, ağ duyuruları, sistem işaretçileri) varlığını ve fiziksel imzasını ele alır. Uydu erişim reçetesi verilmez. Kendi ülkenin ve sürümünün mevzuatını teyit et.

> Güncellik uyarısı: Bu alan çok hızlı değişiyor. Uydu sayıları, frekans tahsisleri, D2D hizmet durumları ve şirket isimleri (örneğin Project Kuiper → Amazon Leo, 2025 sonu) aylar içinde değişebilir. Aşağıdaki değerler 2025–2026 dönemine aittir; kritik bir karar için **operatör/ITU/FCC kaynağından teyit edilmeli** notunu ciddiye al. Fiziksel ilkeler (Doppler, link bütçesi, handover geometrisi) değişmez; sayılar değişir.

---

## İÇİNDEKİLER

1. [Neden LEO? Mega-Konstelasyon Devriminin Fiziği](#1)
2. [Yörünge Kabukları, İnter-Satellite Link ve Yer Ağı Geçitleri](#2)
3. [LEO Doppler ve Geçiş Süresi: Sayısal Derinlik](#3)
4. [LEO Link Bütçesi: Bölüm 11'den LEO'ya Taşıma](#4)
5. [Faz-Dizili Kullanıcı Terminali vs Mekanik İzleme (Bölüm 27 Bağı)](#5)
6. [Handover: LEO'da Sürekli Uydu Değiştirme](#6)
7. [STARLINK Derin: Mimari, Frekanslar, Sinyal Yapısı](#7)
8. [Starlink Gözlemi (Pasif) ve PNT Potansiyeli](#8)
9. [IRIDIUM Derin: NEXT Konstelasyonu, L-Bant, TDMA/FDMA, SBD](#9)
10. [Iridium Ağ Sinyallerini Gözlemek (gr-iridium / iridium-toolkit Kavram)](#10)
11. [Diğer Konstelasyonlar: OneWeb, Amazon Leo, Telesat, Globalstar, Orbcomm, Inmarsat](#11)
12. [Uydu-IoT: Küçük-Uydu Mesajlaşma (Swarm, Astrocast, Lacuna, Kineis/Argos)](#12)
13. [Doğrudan-Cihaz (D2D / Direct-to-Cell) ve 3GPP NTN](#13)
14. [Güvenlik (Savunma Perspektifi): Saldırı Yüzeyi, Kimlik Doğrulama, Gizlilik](#14)
15. [Alıştırmalar (Yasal, Pasif/RX)](#15)
16. [Hızlı Referans ve Diğer Bölümler](#16)

---

<a id="1"></a>
## 1. Neden LEO? Mega-Konstelasyon Devriminin Fiziği

Bölüm 11, GEO ile LEO arasındaki temel dengeyi kurdu: GEO az uydu + büyük gecikme + sabit anten; LEO çok uydu + küçük gecikme + hareketli anten. Mega-konstelasyon devrimi, bu denklemin LEO tarafını ekonomik ve teknik olarak mümkün kılan üç gelişmenin kesişiminden doğdu. Bunları anlamadan Starlink'in neden bir kırılma olduğu kavranamaz.

**Birinci itki: gecikme.** Geri dönelim. GEO için tek yön gecikme ~120 ms, gidiş-dönüş (RTT) ~240 ms, çift atlamalı bağlarda ~480–600 ms idi. Bu, interaktif uygulamalar (video konferans, çevrimiçi oyun, finansal işlem, VPN) için kabul edilemez. LEO ise uyduyu 36.000 km yerine 550 km'ye indirir; mesafe ~65 kat azalır. Tek yön gecikme birkaç milisaniyeye düşer ve uçtan uca RTT, ağ mimarisi iyiyse karasal fiberle yarışacak ~20–50 ms mertebesine iner. İşte mega-konstelasyonların ana satış argümanı budur: **düşük gecikme, GEO'nun yapısal olarak veremediği şey.**

**İkinci itki: kapasite ve frekans yeniden kullanımı.** Bir GEO uydusu, kapsadığı geniş alandaki tüm kullanıcılarla aynı spektrumu paylaşmak zorundadır; toplam kapasite, görece sınırlı bir spektrumu büyük bir nüfusa bölmekle kısıtlanır. LEO mega-konstelasyonu ise gökyüzünü onlarca dar nokta-hüzme (spot beam) ile döşer; her hüzme küçük bir coğrafi hücreyi aydınlatır ve birbirinden yeterince uzak hücreler aynı frekansı yeniden kullanabilir (tıpkı karasal hücresel ağdaki frekans yeniden kullanımı gibi). Binlerce uydu × her uyduda onlarca hüzme = aynı spektrumun yüzlerce-binlerce kez yeniden kullanımı. Toplam sistem kapasitesi, GEO'nun ulaşamayacağı terabit/saniye mertebesine çıkar.

**Üçüncü itki: fırlatma ekonomisi ve seri üretim.** Yeniden kullanılabilir roketler ve uyduların otomotiv benzeri seri üretimi, "binlerce uydu" fikrini bütçeye sığdırdı. Bir GEO uydusu el yapımı, pahalı ve uzun ömürlü (15+ yıl) iken, bir LEO mega-konstelasyon uydusu ucuz, çok sayıda ve kısa ömürlüdür (~5 yıl; sürekli yenilenir). Bu, devrimin görünmeyen ama belirleyici ayağıdır.

```
   GECİKME karşılaştırması (tek yön, dik geçiş — Bölüm 11'den ileri):

   GEO  (35.786 km) ████████████████████████████████  ~120 ms
   MEO  (~8.000 km) ██████████                          ~27 ms
   LEO  (~550 km)   █                                   ~1,8 ms
                    └──────────────────────────────────► gecikme
   Not: uçtan uca RTT'ye ağ/gateway/yönlendirme gecikmesi eklenir.
   LEO'nun avantajı "düşük taban gecikme"dir; ağ tasarımı bunu korumalı.
```

Bedeli nedir? LEO mega-konstelasyonu üç yeni problem getirir; bu bölümün teknik omurgası bunları çözmektir:

1. **Sürekli hareket** — hiçbir uydu sabit değil; her biri dakikalar içinde geçer (Bölüm 11'deki LEO geçiş süresi). Kullanıcı terminali sürekli yeni uyduya geçmeli (handover, Bölüm 6).
2. **Büyük Doppler** — Ku/Ka gibi yüksek frekanslarda LEO Doppler'i yüzlerce kHz'e ulaşır (Bölüm 3). Sürekli düzeltme şart.
3. **İzleme** — hızlı geçen uyduyu izlemek için ya motorlu anten ya da **elektronik tarama yapan faz-dizili anten** gerekir (Bölüm 5, Bölüm 27 bağı). Mega-konstelasyonların tercihi ikincisidir.

Mühendislik sezgisi: **LEO mega-konstelasyonu, GEO'nun "tek büyük uydu" modelini "binlerce küçük, hareketli, sürekli el değiştiren hücre" modeline çevirir. Karasal hücresel ağı gökyüzüne taşımak gibidir — ama hücreler saatte ~27.000 km hızla uçar.**

---

<a id="2"></a>
## 2. Yörünge Kabukları, İnter-Satellite Link ve Yer Ağı Geçitleri

Mega-konstelasyon "rastgele binlerce uydu" değildir; dikkatle tasarlanmış **yörünge kabukları (orbital shells)** sisteminden oluşur. Bir kabuk, aynı yükseklik ve eğimde, belirli sayıda yörünge düzlemine (orbital plane) ve her düzlemde belirli sayıda uyduya bölünmüş bir alt-konstelasyondur. Walker konstelasyon deseni denen bu geometri, gökyüzünü tekdüze kaplamayı amaçlar.

```
   YÖRÜNGE KABUĞU yapısı (kavramsal):

   Kabuk = { yükseklik h, eğim i, P düzlem, her düzlemde S uydu }

        düzlem 1   düzlem 2   düzlem 3  ...  düzlem P
          ☉          ☉          ☉              ☉
          ☉          ☉          ☉              ☉      her düzlem
          ☉          ☉          ☉              ☉      Dünya çevresinde
          ☉          ☉          ☉              ☉      bir halka (S uydu)
          ☉          ☉          ☉              ☉
        \────────── eğim i ile Dünya'yı sarar ──────────/

   Toplam uydu = P × S.   Farklı eğimli birden çok kabuk,
   farklı enlemleri (kutuplar dahil) kapsamak için üst üste konur.
```

Eğim seçimi, kapsamanın nereye yoğunlaşacağını belirler. Düşük eğim (örneğin ~53°) nüfusun yoğun olduğu orta enlemlere odaklanır; yüksek eğim (~70–98°, kutupsal yakın) kutup bölgelerini de kapsar. Bir mega-konstelasyon tipik olarak farklı eğimli birden çok kabuk birleştirir: ana kabuk orta enlemleri yoğun kaplar, ek kutupsal kabuklar yüksek enlemleri tamamlar. Bölüm 11'deki kutupsal/güneş-senkron yörünge mantığının çok-kabuklu ölçeklenmiş halidir bu.

**İnter-Satellite Link (ISL — uydular-arası bağ).** Bölüm 11 çapraz-linki tanıttı; mega-konstelasyonda bu, sistemin can damarıdır. Klasik mimaride her paket yere iner ve yerden tekrar yükselir (uydu-gateway-internet). Ama okyanus ortasında, kutupta ya da gateway'in olmadığı bir bölgede uçan uydunun yere inecek bir geçidi yoktur. ISL bunu çözer: uydu, paketi yere hiç indirmeden komşu uyduya **lazer (optik)** bağ ile aktarır; paket uzayda uydudan uyduya zıplayarak en yakın gateway'i olan uyduya ya da doğrudan hedef bölgenin uydusuna ulaşır.

```
   ISL'siz (yer-merkezli) vs ISL'li (lazer mesh):

   ISL YOK:                          ISL VAR (lazer mesh):
   kullanıcı                         kullanıcı
      │                                 │
      ▼                          ☉≈≈lazer≈≈☉≈≈lazer≈≈☉
   ☉ uydu                        │   (uzayda zıplama)    │
      │ (gateway şart)           ▼                       ▼
      ▼                       gateway              hedef kullanıcı
   GATEWAY (yerde olmalı)     gerekmeden          (okyanus/kutup
      │                       paket taşınır        üstünde bile)
      ▼
   internet
```

Lazer ISL'nin iki büyük getirisi vardır. Birincisi **kapsama bağımsızlığı**: gateway olmayan bölgelerde hizmet. İkincisi **gecikme avantajı**: ışık boşlukta cam fiberden ~%47 daha hızlı ilerler (fiberde kırılma indisi ışığı yavaşlatır). Uzun mesafelerde (kıtalararası) uzaydan lazer zıplama, yer fiberinden düşük gecikmeli olabilir — bu yüzden bazı düşük gecikme kritik uygulamalar (finansal arbitraj gibi) için mega-konstelasyon ISL'si cazip görülür. Pasif gözlem açısından kritik not: **ISL sinyalleri yere inmez ve optiktir; yerden RF ile alınamaz.** Bunlar uydular arası, son derece dar hüzmeli optik bağlardır.

**Yer ağı geçitleri (gateway / ground station).** Mega-konstelasyon, internet omurgasına gateway'ler üzerinden bağlanır. Gateway, büyük çanaklar, yüksek frekanslı (genelde Ka-bant) bağ ekipmanı ve karasal fiber bağlantısı taşır. Bir gateway, üzerinden geçen uydularla yüksek kapasiteli besleme bağı (feeder link) kurar ve kullanıcı trafiğini internete taşır. ISL'li bir sistemde gateway sayısı azaltılabilir (uydular trafiği gateway'i olan bölgeye taşır); ISL'siz bir sistemde her kapsama bölgesinin görüş alanında bir gateway şarttır.

```
   Mega-konstelasyon uçtan uca veri yolu:

   [kullanıcı terminali] ──Ku──▶ [uydu A] ──lazer ISL──▶ [uydu B]
                                                              │ Ka feeder
                                                              ▼
                                                        [GATEWAY] ──fiber──▶ [internet POP]
                                                              ▲
                                                        [PoP / veri merkezi]
```

Bu üç katman — kabuklar (geometri), ISL (uzay omurgası), gateway (yer kapısı) — birlikte mega-konstelasyonun "ağ" kimliğini oluşturur. Bölüm 11'in tek-uydu link mimarisinden farkı budur: burada uydu bir düğüm, konstelasyon bir yönlendirilmiş ağdır.

---

<a id="3"></a>
## 3. LEO Doppler ve Geçiş Süresi: Sayısal Derinlik

Bölüm 11 Doppler'in fiziğini ve "S eğrisi" şeklini verdi. Mega-konstelasyon ve uydu-IoT, yüksek frekanslarda çalıştığı için Doppler'i Bölüm 11'in NOAA APT örneğinin çok ötesine taşır; burada sayısal derinleştiriyoruz.

LEO uydusunun yörünge hızı, dairesel yörünge için:

```
        ┌─────────────────────────┐
        │   v_orbit = √( μ / a )   │
        └─────────────────────────┘

  μ = 3,986 × 10^14 m³/s²   (Dünya standart yerçekim parametresi)
  a = R_dünya + h           (yörünge yarıçapı, m)

  Örnek: h = 550 km → a ≈ 6.921 km = 6,921 × 10^6 m
  v_orbit = √(3,986e14 / 6,921e6) ≈ √(5,759e7) ≈ 7.589 m/s ≈ 7,59 km/s
```

Yani 550 km'lik bir Starlink kabuğunda uydu ~7,6 km/s gider. Şimdi Doppler. Bölüm 11'deki `Δf ≈ f₀ · (v_radyal / c)` formülünde kritik olan radyal bileşendi. Maksimum radyal hız, uydu ufukta doğarken/batarken oluşur ve yörünge hızının bir kısmıdır (tam tepe geçişinde radyal hız sıfırdır). Kabaca, alçak yükseklik açılarında radyal hız yörünge hızının önemli bir bölümüne ulaşabilir.

```
   MAKSİMUM Doppler büyüklüğü — frekansa göre (h≈550 km LEO, kaba):

   Taşıyıcı f₀     | Δf_max mertebesi (yaklaşık)  | Yorum
   ----------------|------------------------------|---------------------------
   137 MHz (APT)   | ± ~3 kHz                      | Bölüm 11; bant içinde yutulur
   1,6 GHz (Iridium)| ± ~37 kHz                    | dar TDMA kanalı için ciddi
   11 GHz (Starlink DL)| ± ~250 kHz                | büyük; sürekli düzeltme şart
   14 GHz (Starlink UL)| ± ~320 kHz                | uplink'te terminal ön-düzeltir
   30 GHz (Ka gateway)| ± ~680 kHz                 | gateway izleme + düzeltme

   Kaba kural: Δf ölçekler f₀ ile. v_radyal,max ~ birkaç km/s alınırsa
   her GHz başına ~birkaç-on kHz kayma çıkar. Tam değer geometri (geçiş
   açısı, yükseklik) ile değişir; bu mertebe sezgisidir, kesin değer değil.
```

Starlink'in 11 GHz downlink'inde ~250 kHz mertebesindeki Doppler, kanalın 240 MHz genişliğine (Bölüm 7) göre küçük bir oran olsa da, OFDM alt-taşıyıcı aralığına göre büyüktür ve mutlaka telafi edilir. Iridium'un 1,6 GHz'inde ~37 kHz, dar TDMA patlamaları için sinyali kaçıracak büyüklüktedir; sistem ve gözlemci yazılımı bunu düzeltmek zorundadır.

**Doppler değişim hızı (Doppler rate).** Sadece kayma değil, kaymanın ne kadar hızlı değiştiği de önemlidir. Tam tepe geçişinde radyal hız sıfırdan geçerken Doppler en hızlı değişir (eğrinin en dik noktası). Bu, alıcının frekans takip döngüsünün (PLL/FLL) ne kadar çevik olması gerektiğini belirler. Yüksek yükseklik açılı (zenite yakın) geçişler en güçlü sinyali verir ama en hızlı Doppler değişimini de getirir — sinyal kalitesi ile takip zorluğu arasında bir denge.

**Geçiş süresi (pass duration).** Bir LEO uydusu gökyüzünde ne kadar görünür kalır? Geometriye bağlıdır: yükseklik arttıkça ufuk çizgisinin altındaki görüş alanı (footprint) genişler, geçiş uzar; geçişin maksimum yükseklik açısı arttıkça (tepeye yakın geçiş) süre uzar.

```
   Görüş yarım-açısı (uydu, minimum yükseklik açısı ε ile görülürken):

   Maksimum görünürlük süresi yaklaşık (tepe geçişi, ε minimum açı):

        T_pass ≈ (2 / ω_rel) · arccos[ (R_dünya/a) · cos(ε) ... ]
                 (basitleştirilmiş; tam ifade küresel geometri ister)

   Pratik mertebe (h ≈ 550 km, ε = 25°):
     yüksek geçiş (tepeye yakın):  ~4–6 dakika görünür
     alçak geçiş (ufka yakın):     ~1–2 dakika görünür

   Karşılaştır: h ≈ 781 km (Iridium):  tepe geçişi ~7–10 dakika
                h ≈ 1200 km (OneWeb):  tepe geçişi daha da uzun
```

Yükseklik arttıkça (Iridium 781 km, OneWeb 1200 km) geçiş süresi uzar ama gecikme de artar — yine o LEO dengesi. Düşük 550 km kabuğu en düşük gecikmeyi verir ama en kısa geçiş ve en sık handover'ı gerektirir. Mühendislik sezgisi: **kabuk yüksekliği seçimi, gecikme ile geçiş süresi/handover sıklığı arasında bir takastır. Starlink düşük gecikme için alçak gider, bedelini sık handover ile öder.**

---

<a id="4"></a>
## 4. LEO Link Bütçesi: Bölüm 11'den LEO'ya Taşıma

Bölüm 8 (Bölüm 11 içinde) link bütçesini GEO için kurdu: FSPL ~205 dB (36.000 km, 12 GHz), EIRP/G/T/C/N0 çatısı. LEO için aynı denklem geçerlidir ama mesafe terimi dramatik biçimde küçüldüğü için tablo tersine döner — ve bu, mega-konstelasyonun neden küçük kullanıcı terminaliyle çalışabildiğinin matematiksel sebebidir.

FSPL'yi LEO mesafesi için hesaplayalım (Bölüm 8'deki formülle):

```
   FSPL(dB) = 20·log₁₀(d_km) + 20·log₁₀(f_GHz) + 92,45

   GEO downlink:  d = 36.000 km, f = 12 GHz → ~205 dB   (Bölüm 8)
   LEO downlink:  d = 550 km,    f = 11 GHz

   FSPL = 20·log₁₀(550) + 20·log₁₀(11) + 92,45
        = 20·(2,740) + 20·(1,041) + 92,45
        ≈ 54,8 + 20,8 + 92,45
        ≈ 168 dB

   FARK: 205 − 168 ≈ 37 dB daha AZ kayıp (LEO lehine)
```

37 dB. Bu devasa bir farktır — gücün ~5.000 katı. Bölüm 8'de gördük ki GEO'nun 205 dB'lik kaybı büyük çanak + soğuk LNB'yi zorunlu kılıyordu. LEO'da 37 dB'lik kazanç, bu bütçeyi öyle gevşetir ki:

- **Daha küçük kullanıcı anteni** yeterli olur (büyük çanak yerine düz faz-dizili panel, Bölüm 5).
- **Daha düşük uydu EIRP'i** ile aynı C/N0 sağlanır (ya da aynı EIRP ile çok daha yüksek veri hızı).
- En kritiği: bu bütçe gevşekliği, sıradan bir **cep telefonunun** zayıf anteniyle bile (D2D, Bölüm 13) uyduya bağlanmasını teorik olarak mümkün kılar — GEO'da bu fiziksel olarak imkânsızdır.

Ancak LEO bedava öğle yemeği değildir. Mesafe sürekli değişir (uydu doğup batarken d 550 km'den ~2.000 km'ye çıkar, ufka yakın geçişte); dolayısıyla FSPL ve dolayısıyla C/N0 geçiş boyunca dalgalanır. Sistem bunu **adaptif kodlama/modülasyon** (Bölüm 9'daki ACM mantığı) ve hüzme yönetimi ile telafi eder: uydu zenitteyken yüksek verimli MODCOD, ufka yakınken dayanıklı MODCOD.

```
   LEO link bütçesi geçiş boyunca DEĞİŞİR (GEO'da sabitti):

   C/N0
    │        ╭────────╮          ← zenit: en yakın (d~550km), en yüksek C/N0
    │      ╭─╯        ╰─╮            → yüksek verimli MODCOD (32APSK gibi)
    │    ╭─╯            ╰─╮
    │  ╭─╯                ╰─╮      ← ufka yakın: d~2000km, FSPL ~+11 dB
    │ ─╯                    ╰─     → dayanıklı MODCOD (QPSK gibi) + handover
    └──────────────────────────► geçiş süresi (dakikalar)
       doğuş      zenit      batış
```

Bu yüzden mega-konstelasyon terminali GEO çanağı gibi "kur ve unut" değildir: sürekli en iyi uyduyu seçer, link bütçesini her uydu için yeniden değerlendirir ve handover yapar. Mühendislik sezgisi: **LEO'nun 37 dB'lik mesafe armağanı, küçük terminali ve D2D'yi mümkün kılar; bedeli, sabit GEO bütçesinin yerini geçiş boyunca dalgalanan dinamik bir bütçeye bırakmasıdır.**

---

<a id="5"></a>
## 5. Faz-Dizili Kullanıcı Terminali vs Mekanik İzleme (Bölüm 27 Bağı)

Bölüm 11, LEO uydusunu izlemek için ya elle çevrilen yönlü anten, ya motorlu azimut+yükseklik takip sistemi, ya da yön bağımsız geniş açılı anten (QFH/turnike) gerektiğini söylemişti. Mega-konstelasyon, hiçbirini kullanmaz; dördüncü ve en güçlü yolu seçer: **elektronik tarama yapan faz-dizili anten (phased array).** Bu konunun fiziği Bölüm 27'de (Anten Dizileri ve Beamforming) derinlemesine işlenir; burada mega-konstelasyona özgü uygulamayı ele alıyoruz.

Mekanik izleme ile elektronik tarama arasındaki temel fark, hüzmenin nasıl yönlendirildiğidir. Mekanik antende çanak fiziksel olarak döner; ağır, yavaş, aşınan parçalar içerir ve aynı anda yalnızca bir yöne bakar. Faz-dizili antende ise sabit, düz bir panel üzerinde yüzlerce-binlerce küçük anten elemanı vardır; her elemana giden sinyalin fazı elektronik olarak ayarlanarak ana hüzme istenen yöne "eğilir" — hiçbir hareketli parça olmadan, mikrosaniyeler içinde.

```
   FAZ-DİZİLİ hüzme yönlendirme (Bölüm 27 fiziği, LEO uygulaması):

   Elemanlar:   E1   E2   E3   E4   E5  ...  (düz panel, sabit)
   Faz gecikme: 0    Δφ   2Δφ  3Δφ  4Δφ      (elektronik, ayarlanabilir)
                │    │    │    │    │
                ▼    ▼    ▼    ▼    ▼
   Eş-faz cepheleri belirli bir açıda birleşir:

        hüzme yönü θ  ↗
                    ╱╱╱        sin(θ) = (λ · Δφ) / (2π · d_eleman)
                  ╱╱╱
         ════════════════════  panel (yatay, gökyüzüne bakar)
         E1  E2  E3  E4  E5

   Δφ'yi elektronik değiştir → hüzme anında yeni uyduya kayar.
   Hareketli parça YOK. Bir uydudan diğerine geçiş mikrosaniyeler.
```

Mega-konstelasyon için bu neden zorunlu? Çünkü:

1. **Hız:** Uydu gökyüzünü dakikalarda geçer; hüzme onu pürüzsüz izlemeli. Elektronik tarama, mekanik motordan çok daha hızlı ve titreşimsiz izler.
2. **Handover:** Bir uydu batarken diğeri doğar; faz-dizili anten, neredeyse anında yeni uyduya hüzme atabilir (bazı tasarımlar aynı anda iki hüzme tutup "kırılmadan" geçiş yapabilir, Bölüm 6).
3. **Dayanıklılık:** Hareketli parça yok; çatıya monte edilen, kar/buz/rüzgâra dayanan, bakım istemeyen düz panel.
4. **Maliyet ve seri üretim:** Tüketici fiyatına faz-dizili anten üretebilmek (silikon faz kaydırıcılar, baskı devre antenler), mega-konstelasyonun tüketiciye ulaşmasının anahtarıdır.

Starlink'in kullanıcı terminali ("Dishy" olarak bilinen düz panel), tam olarak budur: yüzlerce elemanlı bir Ku-bant faz-dizili anten. İlk kurulumda mekanik bir motorla kabaca gökyüzüne yönelir (en açık gökyüzü dilimini bulmak için), ama uyduyu **izleme** tamamen elektroniktir — panel sabit kaldığı halde hüzme uydudan uyduya elektronik olarak atlar.

Faz-dizili antenin LEO'ya özgü bir zorluğu, **tarama açısıyla kazanç düşüşüdür** (scan loss). Panel yatay durduğunda zenite (tam tepe) en yüksek kazançla bakar; hüzme ufka doğru eğildikçe etkin açıklık (panelin görünen izdüşümü) küçülür ve kazanç düşer (kabaca cos θ ile). Bu yüzden faz-dizili terminal, ufka çok yakın uyduları zayıf görür; sistem bunu, mümkünse yüksek açılı uyduları tercih ederek ve panel eğimini optimize ederek yönetir. Bu, Bölüm 27'de işlenen genel faz-dizili sınırlamasının LEO bağlamındaki yansımasıdır.

Mühendislik sezgisi: **GEO çanağı bir kez kilitlenir ve unutulur (Bölüm 11); LEO mega-konstelasyon terminali, hareketsiz görünmesine rağmen içinde sürekli elektronik hüzme dansı yapan bir faz-dizili sistemdir. Görünürdeki sadelik, Bölüm 27'deki beamforming karmaşıklığını gizler.**

---

<a id="6"></a>
## 6. Handover: LEO'da Sürekli Uydu Değiştirme

Karasal hücresel ağda telefon, bir baz istasyonundan diğerine geçerken handover yapar — ama baz istasyonları sabittir, hareket eden sensin. LEO mega-konstelasyonunda durum tersine döner ve şiddetlenir: **kullanıcı sabit, hücreler (uydular) saatte ~27.000 km hızla uçar.** Bu yüzden kullanıcı tamamen hareketsiz dursa bile, birkaç dakikada bir hizmet veren uydu değişir ve handover kaçınılmazdır.

```
   LEO HANDOVER zaman çizgisi (sabit kullanıcı, uçan uydular):

   t=0          t=3dk        t=6dk        t=9dk
   ☉ Uydu A     ☉ A (batıyor) 
        ╲            ╲       ☉ Uydu B      ☉ B (batıyor)
         ╲ hizmet     ╲       ╱ hizmet          ╲    ☉ Uydu C
          ╲            ╲ HO  ╱                    ╲    ╱ hizmet
           ▼            ▼ ▼ ▼                      ▼  ╱
        [kullanıcı terminali — yerinde sabit, hüzmesi A→B→C atlar]

   HO = handover anı. Her ~birkaç dakikada bir tekrar eder.
   "Make-before-break": B'ye bağlan, sonra A'yı bırak (kesintisiz).
```

Handover stratejileri, sistemin gecikme ve kesinti hedeflerine göre değişir:

- **Make-before-break (kır-madan-bağlan):** Terminal, eski uyduyla bağı kesmeden önce yeni uyduyla bağ kurar; faz-dizili anten aynı anda iki hüzme tutabiliyorsa kesinti neredeyse sıfırdır. Tercih edilen yöntem budur.
- **Zamanlanmış handover:** Konstelasyon deterministik olduğu için (uydu konumları yörünge mekaniğinden tam bilinir), handover anları önceden hesaplanabilir; terminal hangi uyduya ne zaman geçeceğini bilir. Bu, "ölç ve karar ver" yerine "planla ve uygula" yaklaşımıdır.
- **Hüzme-içi vs uydu-arası handover:** Bazen aynı uydu içinde bir hüzmeden diğerine geçilir (uydu hâlâ görünür ama kullanıcı farklı bir spot beam hücresine giriyor); bazen tümüyle yeni bir uyduya geçilir. İkisi farklı maliyet ve zamanlama taşır.

Handover'ın spektrum imzası, pasif gözlemci için ilginçtir (Bölüm 8). Bir LEO uydusu downlink'ini spektrumda izlerken, uydu battıkça o taşıyıcının zayıflayıp kaybolduğunu ve yeni doğan uydunun taşıyıcısının güçlenerek belirdiğini gözlemleyebilirsin — bu, handover ritminin RF'teki görünür izidir. Tekil bir kullanıcının trafiğini görmezsin (içerik şifreli ve hüzme dar), ama konstelasyonun "nefes alıp verme" ritmini, geçiş periyotlarını ve uydu yoğunluğunu spektrumdaki bu geliş-gidişlerden çıkarabilirsin.

```
   Pasif gözlemde handover'ın izi (Starlink DL bandında, kavramsal):

   güç
    │  ▂▄▆█▆▄▂              ▂▄▆█▆▄▂              ← her tümsek bir uydu geçişi
    │ ╱       ╲           ╱       ╲                taşıyıcı belirir → tepe → kaybolur
    │╱  uydu A ╲ ───── ╱  uydu B  ╲ ─────       (handover periyodu ~birkaç dk)
    └──────────────────────────────────────► zaman
    Not: bu downlink VARLIĞININ izidir; içerik DEĞİL. Pasif/yasal.
```

Mühendislik sezgisi: **LEO handover, karasal handover'ın tersidir — ağ hareket eder, kullanıcı durur. Konstelasyon deterministik olduğu için handover ölçüm değil, takvim işidir; ve bu takvimin ritmi, pasif spektrum gözleminde uydu geçişleri olarak görünür hale gelir.**

---

<a id="7"></a>
## 7. STARLINK Derin: Mimari, Frekanslar, Sinyal Yapısı

Starlink, mega-konstelasyon olgusunun en büyük ve en çok incelenen örneğidir. Burada mimarisini, frekanslarını ve — açık akademik literatürden çıkarılmış — fiziksel sinyal yapısını ele alıyoruz. Tüm bu bilgi pasif gözlem ve açık yayınlanmış araştırmaya dayanır; kullanıcı içeriği şifrelidir ve burada çözülmez.

> Teyit notu: Starlink hızla evrilen bir sistemdir (uydu nesilleri V1.0, V1.5, V2 mini, V2; frekans tahsisleri ve hüzme planları değişir). Aşağıdaki değerler 2025–2026 dönemine ve açık FCC/akademik kaynaklara aittir. Operasyonel detaylar SpaceX tarafından sık güncellenir; kritik bir karar için güncel kaynak teyit edilmeli.

**Mimari.** Starlink binlerce uydudan oluşan, çok-kabuklu bir LEO konstelasyonudur. Ana kabuklar ~540–570 km civarındadır (düşük gecikme için alçak); farklı eğimlerde kabuklar orta ve yüksek enlemleri kapsar. Yeni nesil uydular lazer ISL taşır (Bölüm 2), böylece okyanus ve kutup kapsaması ile gateway-bağımsız yönlendirme sağlanır. Her uydu, kullanıcı bağı için **çok sayıda Ku-bant faz-dizili anten** ve gateway/besleme bağı için **Ka-bant** (ve bazı raporlara göre E-bant) anten taşır.

**Kullanıcı terminali ("Dishy").** Bölüm 5'te ele alındığı gibi, yüzlerce elemanlı bir Ku-bant faz-dizili düz paneldir. Elektronik tarama ile uyduyu izler; ilk kurulumda mekanik motorla kabaca yönelir, izleme elektroniktir. Bu, Bölüm 27'deki beamforming ilkelerinin tüketici ölçeğinde uygulanmasıdır.

**Frekanslar (teyit edilmeli — değişebilir).** Açık FCC tahsisleri ve teknik kaynaklara göre:

| Bağ yönü | Bant | Frekans aralığı (yaklaşık) | Not |
|---|---|---|---|
| Kullanıcı **downlink** (uydu→Dishy) | Ku | ~10,7–12,7 GHz | FSS Ku; gözlemlenebilir |
| Kullanıcı **uplink** (Dishy→uydu) | Ku | ~14,0–14,5 GHz | terminal yayını |
| Gateway **uplink** (yer→uydu) | Ka | ~27,5–30 GHz | besleme bağı |
| Gateway **downlink** (uydu→yer) | Ka | ~17,8–20,2 GHz | besleme bağı |
| TT&C | (ayrı tahsis) | — | uydu yönetimi (Bölüm 11/6) |
| ISL (uydular arası) | optik (lazer) | — | RF değil; yerden alınamaz |

Bölüm 11'deki bant tablosuyla tutarlı: Ku kullanıcı için (küçük anten + makul yağmur dayanımı), Ka yüksek kapasiteli gateway için (bol bant genişliği, yağmura daha duyarlı ama gateway'lerde site diversity ve uplink power control ile yönetilir — Bölüm 7). Doppler bu frekanslarda büyüktür (Bölüm 3: Ku DL'de ~250 kHz mertebesinde) ve sürekli düzeltilir.

**Sinyal yapısı (açık akademik literatürden).** Texas Üniversitesi Radyonavigasyon Laboratuvarı'nın (Humphreys ve ekibi) açık yayınlanmış "kör sinyal tanımlama" çalışması, Starlink Ku-bant downlink'inin fiziksel yapısını içerik çözmeden karakterize etmiştir. Bu, pasif gözlem ve PNT araştırması için önemli bir referanstır (Bölüm 8):

```
   STARLINK Ku-bant downlink — açık literatürden fiziksel yapı (kavramsal):

   ┌─────────────────────────────────────────────────────────┐
   │ Kanal genişliği:        ~240 MHz                          │
   │ Dalga biçimi:           OFDM (dik frekans bölmeli çoğullama)│
   │ Alt-taşıyıcı sayısı:    ~1024                              │
   │ Çevrimsel önek (CP):    ~0,133 µs koruma aralığı          │
   │ Sembol süresi (~):      4,4 µs aralıklar                  │
   │ Çerçeve (frame) periyodu: 1/750 s (~1,33 ms)              │
   │ Çerçeve başı senkr.:    PSS → SSS (4QAM OFDM sembolü)     │
   └─────────────────────────────────────────────────────────┘

   OFDM kanalı (240 MHz) — alt-taşıyıcılara bölünmüş:

   güç
    │ ▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏  ← ~1024 dik alt-taşıyıcı
    │ ▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏▏     (her biri dar bant taşır)
    └────────────────────────────────────► frekans
      │◄──────────── ~240 MHz ──────────►│
```

OFDM seçimi anlamlıdır: karasal 4G/5G ve modern geniş bant sistemlerinin tercih ettiği dalga biçimidir; yüksek veri hızı, frekans seçici sönümlemeye dayanıklılık ve esnek alt-taşıyıcı tahsisi sağlar (Bölüm 1/5'teki modülasyon temelleri). Beam-hopping (hüzme atlama): uydu, sınırlı sayıda hüzme oluşturucuyu zaman içinde farklı hücreler arasında hızla paylaştırarak geniş bir alanı "sırayla" aydınlatabilir — sabit sürekli hüzme yerine zaman-paylaşımlı hüzme. Bu, kapasiteyi talebe göre dağıtmanın yoludur. PSS/SSS senkronizasyon sembolleri (karasal hücresel ağdaki senkronizasyon sinyallerine kavramsal benzerlik), alıcının çerçeveye kilitlenmesini ve zamanlama/frekans senkronizasyonunu sağlar; bu sembollerin yapısı PNT için kritiktir (Bölüm 8).

Mühendislik sezgisi: **Starlink downlink'i, GEO yayın uydularının (Bölüm 9) tek-taşıyıcı DVB-S2 imzasından tamamen farklıdır — geniş (240 MHz), OFDM-tabanlı, çerçeveli ve senkronizasyon sembolleriyle yapılandırılmış, karasal hücresel/geniş bant dünyasına yakın bir dalga biçimidir. Spektrumda bu genişlik ve yapı, onu klasik SATCOM taşıyıcısından ayırır.**

---

<a id="8"></a>
## 8. Starlink Gözlemi (Pasif) ve PNT Potansiyeli

Bu bölüm tamamen **pasif gözlem** (RX) üzerinedir: Starlink kullanıcı downlink'ini spektrumda tanımak, varlığını ve geçiş ritmini gözlemlemek. Kullanıcı içeriği şifrelidir ve burada çözülmez; amaç fiziksel imza okuryazarlığıdır (Bölüm 10'daki "taşıyıcı kimliklendirme" mantığının LEO'ya uzantısı).

**Downlink'i spektrumda tanıma (kavramsal).** Starlink kullanıcı downlink'i Ku-bandındadır (~10,7–12,7 GHz). Bu, sıradan bir SDR'ın doğrudan ulaşamayacağı kadar yüksek bir frekanstır; tıpkı GEO Ku alımında olduğu gibi (Bölüm 14), önce uygun bir **LNB** (Low Noise Block downconverter) ile sinyal ara frekansa (IF) indirilmeli, sonra SDR ile bu IF gözlemlenmelidir. Spektrumda aranan imza, Bölüm 7'de tarif edilen geniş (240 MHz), düz-tepeli OFDM bloklarıdır; ve bunların Bölüm 6'da tarif edilen handover ritmiyle belirip kaybolmasıdır.

```
   PASİF Starlink downlink gözlemi — sinyal zinciri (kavramsal):

   [Ku LNB] ──IF──▶ [geniş bant SDR] ──▶ [spektrum görüntüleme]
   (10.7-12.7 GHz       (≥ kanal           Aranan: ~240 MHz geniş
    → IF'e indir)        genişliğini         OFDM blokları, belirip
                         örnekleyebilen)      kaybolan (handover ritmi)

   Gözlenen (yasal):   downlink VARLIĞI, genişliği, geçiş zamanlaması
   Gözlenmeyen (yapılmaz): şifreli kullanıcı içeriği
```

**Açık-kaynak Starlink gözlem projeleri (kavram).** Akademik ve hobi topluluğu, Starlink downlink'ini araştırmak için açık araçlar geliştirmiştir. Bunlar içerik çözmez; fiziksel sinyali karakterize eder, senkronizasyon sembollerini tespit eder ve PNT için kullanılabilirliğini inceler. Texas Radyonavigasyon Laboratuvarı'nın yayınladığı sinyal yapısı (Bölüm 7) ve buna dayalı sinyal simülatörleri bu çalışmaların temelidir. Bu projelerin amacı, ticari LEO sinyallerinin GNSS'e alternatif/tamamlayıcı bir konumlama kaynağı olup olamayacağını araştırmaktır.

**PNT potansiyeli (Positioning, Navigation, Timing).** En ilginç gelişmelerden biri budur. GNSS (GPS/Galileo, ~20.200 km MEO — Bölüm 11) zayıf sinyalli ve jamming/spoofing'e açıktır (Bölüm 14, Bölüm 20). LEO mega-konstelasyon sinyalleri ise:

- **Çok daha güçlüdür** (LEO ~550 km, GNSS ~20.200 km; FSPL farkı Bölüm 4'teki mantıkla onlarca dB — LEO sinyali yerde çok daha kuvvetli).
- **Hızlı hareket eder** (büyük Doppler — Bölüm 3) ki bu, konum çözümü için zengin geometri çeşitliliği sağlar.
- **Bol sayıdadır** (binlerce uydu, gökyüzünde her an çok sayıda görünür).

Bu üç özellik, Starlink (ve genel olarak LEO) sinyallerini **GNSS-bağımsız PNT** için cazip kılar. Fikir, sinyalin senkronizasyon sembollerinden (PSS/SSS, Bölüm 7) sözde-mesafe (pseudorange) ölçmek ve birden çok uydudan bu ölçümlerle konum çözmektir — tıpkı GNSS'in yaptığı gibi, ama daha güçlü ve daha çeşitli geometriyle. Bu henüz operasyonel bir hizmet değil, aktif bir araştırma alanıdır (opportunistic navigation / signals of opportunity). Savunma açısından önemi: GNSS jamming'in yaygınlaştığı bir ortamda (Bölüm 20, Bölüm 23), LEO sinyalleri yedek bir PNT kaynağı sunabilir.

```
   GNSS vs LEO sinyali — PNT için karşılaştırma:

                    GNSS (GPS/Galileo)      LEO (Starlink vb.)
   Yükseklik        ~20.200 km (MEO)        ~550 km (LEO)
   Sinyal gücü      zayıf (yerde -130dBm)   çok daha güçlü
   Doppler          küçük                   büyük (geometri zengin)
   Jamming direnci  düşük (zayıf sinyal)    yüksek (güçlü sinyal)
   PNT durumu       operasyonel             araştırma (opportunistic)
   Amaç             konum sağlamak          fırsatçı/yedek PNT
```

Mühendislik sezgisi: **Starlink gözlemi, içerik değil yapı okumaktır — downlink'in varlığı, genişliği, senkronizasyon imzası ve handover ritmi. Bu yapının en heyecan verici yan ürünü, mega-konstelasyon sinyallerinin GNSS'e bağımsız bir PNT kaynağı olma potansiyelidir; güçlü ve bol LEO sinyali, jamming'e açık zayıf GNSS'in yedeği olabilir.**

---

<a id="9"></a>
## 9. IRIDIUM Derin: NEXT Konstelasyonu, L-Bant, TDMA/FDMA, SBD

Iridium, mega-konstelasyon kavramının öncüsüdür — Starlink'ten çok önce, çapraz-linkli bir LEO ağı kuran ilk sistemdir. Bugünkü hali **Iridium NEXT** konstelasyonudur. Geniş bant değil, küresel ses/dar-veri/IoT odaklıdır ve mimarisi mega-konstelasyon mantığının erken ve zarif bir örneğidir.

**NEXT konstelasyonu.** Iridium, **66 aktif uydudan** oluşur (artı yörüngede dönen yedekler ve yerde bekleyen yedekler). Uydular **~781 km** yükseklikte, **~86,4° eğimle** (neredeyse kutupsal) altı yörünge düzlemine dağılmıştır. Bu kutupsal yakın geometri, Iridium'a benzersiz bir özellik kazandırır: **gerçek küresel kapsama** — kutuplar dahil yeryüzünün her noktası. Starlink dahil çoğu mega-konstelasyon kutupları zayıf kaplarken (orta enlem odaklı kabuklar), Iridium tasarımı gereği her yeri görür. Bu yüzden kutup araştırması, denizcilik, havacılık ve uzak bölge IoT'sinde tercih edilir.

```
   IRIDIUM NEXT geometrisi — 6 düzlem, kutupsal yakın:

         düzlem 1  2  3  4  5  6
            │   │  │  │  │  │      her düzlemde 11 aktif uydu
        ────┼───┼──┼──┼──┼──┼────  (6 × 11 = 66)
            │   │  │  │  │  │      eğim ~86,4° (neredeyse kutupsal)
        ╲   │   │  │  │  │  │   ╱  yükseklik ~781 km
         ╲  Kutup üstünden geçer ╱  → KÜRESEL kapsama (kutuplar dahil)
          ╲────────────────────╱
                  DÜNYA

   Çapraz-link (ISL): her uydu 4 Ka-bant bağ tutar —
   2 düzlem-içi (öndeki/arkadaki komşu) + 2 düzlem-arası (yan komşular)
   → dinamik halka (ring) topolojisi, ~25 Mbps her bağ.
```

**Çapraz-link (ISL).** Iridium'un imza özelliği, Starlink'ten yıllar önce uyguladığı uydular-arası bağlardır. Her uydu **dört aktif Ka-bant ISL** tutar: ikisi aynı yörünge düzlemindeki önündeki ve arkasındaki uyduya (düzlem-içi), ikisi komşu düzlemlerdeki uydulara (düzlem-arası). Bu, gökyüzünde dinamik bir **halka (ring) topolojisi** oluşturur; bir çağrı, yere hiç inmeden uydudan uyduya zıplayarak dünyanın öbür ucundaki gateway'e taşınabilir. Bu, Iridium'un az sayıda gateway ile küresel hizmet verebilmesinin sırrıdır (okyanus ortasındaki bir telefon, çağrıyı ISL üzerinden karadaki gateway'e ulaştırır).

**L-bant (kullanıcı bağı).** Iridium kullanıcı bağı **L-bandındadır: ~1616–1626,5 MHz.** Bölüm 11'deki L-bant mantığı tam olarak geçerli: bu frekansta hem yol kaybı görece düşük, hem de küçük, yön bağımsız antenle (yama, kısa heliks) çalışılabilir — kullanıcı çanağı tepeye doğrultmak zorunda değildir (uydu telefonu cebe sığar). Bedeli düşük bant genişliğidir; bu yüzden Iridium sesi dar, veri hızı düşüktür. Doppler bu frekansta ~37 kHz mertebesindedir (Bölüm 3) ve sistem bunu düzeltir.

**TDMA/FDMA çoklu erişim.** Iridium, sınırlı L-bant spektrumunu çok sayıda kullanıcıya paylaştırmak için hibrit bir çoklu erişim kullanır: **FDMA** (Frequency Division Multiple Access — bant, frekans kanallarına bölünür) üzerine **TDMA** (Time Division Multiple Access — her kanal, zaman dilimlerine/slot'lara bölünür). Yani spektrum hem frekansta hem zamanda dilimlenir; bir kullanıcıya belirli bir frekans kanalında belirli bir zaman slot'u atanır. Bu, Bölüm 5'teki modülasyon ve çoklu erişim temellerinin uydu bağlamındaki uygulamasıdır.

```
   IRIDIUM çoklu erişim — FDMA × TDMA (kavramsal):

   frekans
    ▲   ┌──┬──┬──┬──┐  ┌──┬──┬──┬──┐   her hücre: FDMA kanalları (frekans)
    │   │s1│s2│s3│s4│  │s1│s2│s3│s4│   × TDMA slotları (zaman)
    │   ├──┼──┼──┼──┤  ├──┼──┼──┼──┤   
    │   │s1│s2│s3│s4│  │s1│s2│s3│s4│   bir kullanıcı = (kanal, slot) çifti
    │   └──┴──┴──┴──┘  └──┴──┴──┴──┘   TDD: gidiş/dönüş zamanda ayrık
    └────────────────────────────────► zaman (TDMA çerçeveleri)
        patlamalı (burst) iletim — sürekli değil
```

İletim **patlamalıdır (burst)**: kullanıcı sürekli değil, kendi zaman slot'unda kısa patlamalar gönderir. Bu, Iridium'un spektrum imzasının ayırt edici özelliğidir — gözlemde (Bölüm 10) sürekli bir taşıyıcı değil, gelip giden kısa patlamalar görürsün.

**SBD (Short Burst Data) ve hizmetler.** Iridium, sesin yanında dar-veri hizmetleri sunar. **SBD (Short Burst Data)**, çok küçük veri paketlerini (onlarca-yüzlerce bayt) iletmek için tasarlanmıştır ve uydu-IoT/M2M'in (makineden makineye) temelidir: uzak bir sensör, bir varlık takip cihazı, bir denizcilik mesajlaşma terminali küçük durum mesajlarını SBD ile gönderir. Iridium ağında milyonlarca aktif IoT cihazı SBD kullanır. Daha yüksek hız gerektiğinde **Iridium Certus** (geniş bant hizmeti) L-bant terminallerle yüzlerce kbps'e kadar çıkar. Bu, Bölüm 12'deki uydu-IoT olgusunun en olgun ve yerleşik örneğidir.

Mühendislik sezgisi: **Iridium, mega-konstelasyonun atasıdır — kutupsal 66 uydu, Ka-bant çapraz-linkli halka, L-bant kullanıcı bağı, FDMA×TDMA patlamalı erişim. Starlink'in geniş bant OFDM dünyasından farklı olarak, Iridium dar-bant patlamalı bir küresel mesajlaşma/ses ağıdır; ve SBD ile uydu-IoT'nin temelini atmıştır.**

---

<a id="10"></a>
## 10. Iridium Ağ Sinyallerini Gözlemek (gr-iridium / iridium-toolkit Kavram)

Iridium, açık-kaynak topluluğunun en çok çalıştığı uydu sistemlerinden biridir; çünkü L-bandı (~1,6 GHz) sıradan bir SDR'ın ulaşabileceği bir frekanstır ve ağ katmanının bir kısmı açık çerçevede yayılır. Bu bölüm, **pasif gözlemin** ne olduğunu ve sınırını net çizer.

> Yasal/etik sınır: Aşağıda anlatılan, Iridium ağının **açık ağ katmanını** (senkronizasyon, ağ duyuruları, sistem işaretçileri gibi şifresiz çerçeveler) gözlemektir. Kullanıcı ses/veri içeriği ve özel mesajlar bu kapsamda DEĞİLDİR; başkalarının iletişim içeriğini ele geçirmek/çözmek çoğu yerde yasadışıdır ve bu kitap onu tarif etmez. Amaç, fiziksel sinyal ve açık ağ sinyalizasyonu okuryazarlığıdır — başkalarının haberleşmesini dinlemek değil.

**Donanım/yazılım zinciri (kavram).** Iridium L-bandı ~1616–1626,5 MHz'dedir. Bu frekansa ulaşabilen bir SDR (Bölüm 2) ve uygun bir L-bant anten (yama ya da heliks, gökyüzünü geniş açıyla gören — Bölüm 3) ile patlamalar alınabilir. Açık-kaynak işleme zinciri kavramsal olarak iki parçadır:

```
   IRIDIUM açık ağ gözlemi — işleme zinciri (kavramsal):

   [L-bant anten] ──▶ [SDR @ ~1,6 GHz] ──▶ [gr-iridium]
   (geniş açı,           (yeterli bant,        patlama (burst)
    gökyüzü)              ~1626 MHz'i kapsar)   tespit + demodülasyon
                                                      │
                                                      ▼
                                              ham patlama bit'leri
                                                      │
                                                      ▼
                                              [iridium-toolkit]
                                              çerçeve türü ayrıştırma:
                                              - IRA (Ring Alert) ← açık ağ duyurusu
                                              - sistem/zamanlama çerçeveleri
                                              - (içerik DEĞİL — ayrıştırılmaz/çözülmez)
```

- **gr-iridium:** Iridium patlamalarını tespit edip demodüle eden açık-kaynak araçtır (GNU Radio tabanlı). Spektrumdaki kısa patlamaları yakalar ve ham bit'lere çevirir. (Alternatif bağımsız uygulamalar da vardır; iridium-sniffer gibi C tabanlı, gömülü/düşük güçlü sistemler için.)
- **iridium-toolkit:** gr-iridium çıktısını alıp çerçeve türlerine ayrıştıran araçtır. Tanıyabildiği açık çerçeveler arasında **IRA (Iridium Ring Alert)** gibi ağ duyuruları vardır — bunlar, uydunun varlığını ve ağ zamanlamasını taşıyan, içerik olmayan sistem sinyalleridir.

**Ne gözlenir, ne gözlenmez?** Sınır nettir:

| Gözlenebilir (açık ağ katmanı, pasif/yasal odak) | Gözlenmez / yapılmaz (kapsam dışı) |
|---|---|
| Patlamaların varlığı, zamanlaması, frekansı | Kullanıcı ses içeriği |
| IRA (Ring Alert) tipi ağ duyuru çerçeveleri | Özel kullanıcı mesaj içeriği |
| Sistem/zamanlama sinyalizasyonu (açık) | Şifreli kullanıcı verisi |
| Doppler eğrisi, geçiş ritmi (uydu izleme) | Üçüncü kişilerin haberleşmesini çözme |

Iridium'un açık ağ katmanının (özellikle Ring Alert gibi düz yayılan sistem çerçevelerinin) gözlemlenebilir olması, akademik güvenlik araştırmalarının da konusu olmuştur (Bölüm 14): sistem sinyallerinin açık olması, ağ yapısının ve bazı meta-verilerin pasif olarak çıkarılabilmesi anlamına gelir, bu da gizlilik/savunma tartışmasının parçasıdır. Önemli olan, gözlemcinin **açık sistem katmanı** ile **şifreli kullanıcı içeriği** arasındaki çizgiyi bilmesi ve ikincisine geçmemesidir.

**Gözlemin değeri.** Pasif Iridium gözlemi, mega-konstelasyon dinamiğini elle tutulur kılar: kendi gökyüzünden geçen Iridium uydularının patlamalarını yakalayıp, Doppler kaymasından geçişlerini izleyebilir (Bölüm 3'teki S eğrisi), ağ duyurularının ritmini görebilir ve TDMA patlamalı yapının (Bölüm 9) RF'teki gerçek görünümünü deneyimleyebilirsin. Bu, "uydu-IoT ağı nasıl bir şeydir" sorusunun teorik değil, gözlemsel cevabıdır.

Mühendislik sezgisi: **Iridium gözlemi, mega-konstelasyon pasif gözleminin en erişilebilir laboratuvarıdır — L-bandı SDR menzilinde, açık ağ katmanı gözlemlenebilir, patlamalı TDMA yapısı ve Doppler eğrisi doğrudan görülebilir. Çizgi nettir: açık sistem sinyalizasyonunu gözlemek serbest, kullanıcı içeriği değil.**

---

<a id="11"></a>
## 11. Diğer Konstelasyonlar: OneWeb, Amazon Leo, Telesat, Globalstar, Orbcomm, Inmarsat

Starlink ve Iridium tek oyuncu değildir. LEO/MEO/GEO yelpazesinde birbirinden farklı tasarım felsefeleriyle çok sayıda konstelasyon vardır. Aşağıdaki tablo bunları yörünge × frekans × kullanım ekseninde özetler (değerler 2025–2026; **teyit edilmeli — hızla değişir**).

| Sistem | Yörünge | Frekans (kullanıcı / besleme) | Kullanım | Durum/Not |
|---|---|---|---|---|
| **Starlink** | LEO ~550 km | Ku (kull.) / Ka (gateway) / lazer ISL | Geniş bant internet, D2D | En büyük; binlerce uydu |
| **Iridium NEXT** | LEO ~781 km, kutupsal | L (kull.) / Ka (ISL+feeder) | Ses, dar-veri, IoT (SBD), küresel | 66 aktif; kutuplar dahil |
| **OneWeb (Eutelsat)** | LEO ~1200 km | Ku (kull.) / Ka (gateway) | Geniş bant (kurumsal/B2B) | ~648 uydu; Eutelsat ile birleşti |
| **Amazon Leo** (eski Project Kuiper) | LEO (alçak) | Ka (kull. ve gateway) | Geniş bant internet | İnşa halinde; ~3.232 planlı; 2025 sonu yeniden adlandırıldı |
| **Telesat Lightspeed** | LEO | Ka | Geniş bant (kurumsal) | İnşa/planlama aşamasında |
| **Globalstar** | LEO ~1.400 km | L/S | Ses, dar-veri, IoT; Apple SOS arka planı | C-3 nesli yenileniyor |
| **Orbcomm** | LEO | VHF (~137–150 MHz) | M2M/IoT (varlık takip, telematik) | Olgun, dar-veri IoT |
| **Inmarsat** | GEO (Bölüm 11) | L (kull.) / C-Ka (feeder) | Denizcilik/havacılık ses-veri, güvenlik | GEO; LEO değil — Bölüm 11 |

Birkaç önemli ayrım:

**OneWeb (artık Eutelsat OneWeb).** ~1200 km'lik daha yüksek bir LEO kabuğunda ~648 uydu işletir; kullanıcıya Ku, gateway'e Ka kullanır (Starlink'e benzer bant düzeni ama daha yüksek yörünge → daha uzun geçiş, biraz daha yüksek gecikme). Tüketiciye doğrudan satıştan çok kurumsal/B2B, telekom backhaul ve devlet pazarına odaklıdır. GEO operatörü Eutelsat ile birleşerek GEO+LEO melez bir portföy oluşturmuştur.

**Amazon Leo (eski Project Kuiper).** Amazon'un Ka-bant LEO geniş bant konstelasyonudur; 2025 sonunda "Amazon Leo" olarak yeniden adlandırılmıştır. Binlerce uydu (planlanan ~3.232) hedeflenir; düzenleyici takvim, konstelasyonun belirli oranının belirli tarihlere kadar yörüngede olmasını şart koşar. Starlink'in en ciddi geniş bant rakibi olarak konumlanır. (Teyit edilmeli — aktif konuşlandırma aşamasında, sayılar sık güncellenir.)

**Telesat Lightspeed.** Ka-bant LEO, kurumsal/telekom odaklı bir başka geniş bant konstelasyonu; planlama/erken konuşlandırma aşamasında.

**Globalstar.** LEO ses/dar-veri/IoT sistemi; teknik açıdan en dikkat çekici yanı, **Apple'ın acil durum uydu hizmetinin (Emergency SOS)** altyapı sağlayıcısı olmasıdır (Bölüm 13). Konstelasyonunu yeni nesil uydularla yeniler.

**Orbcomm.** VHF (~137–150 MHz) bandında çalışan, olgun bir M2M/IoT konstelasyonu; varlık takip, ağır makine telematiği, denizcilik (AIS dahil — Bölüm 11/20 ile bağlantılı) gibi dar-veri uygulamalarına odaklıdır. Düşük frekans → küçük anten, düşük veri hızı (Bölüm 11 L/VHF mantığı).

**Inmarsat.** Burada bir hatırlatma: Inmarsat bir **GEO** sistemidir (LEO/mega-konstelasyon değil), Bölüm 11'de ele alınmıştır. Denizcilik ve havacılıkta L-bant ses/veri ve güvenlik hizmetleri sunar. Mega-konstelasyon bölümünde anılması, karşılaştırma içindir: GEO'nun süreklilik avantajı (tek uydu sürekli görünür) vs LEO'nun düşük gecikme avantajı. Inmarsat operatörü de LEO ortaklıklarıyla melez stratejilere yönelmektedir.

Mühendislik sezgisi: **Mega-konstelasyon manzarası tek tip değil — geniş bant (Starlink, OneWeb, Amazon Leo, Telesat: Ku/Ka, yüksek hız) ile dar-veri/IoT (Iridium, Globalstar, Orbcomm: L/VHF/S, düşük hız, küçük anten) iki ayrı dünyadır. Frekans bandı (Bölüm 11 tablosu) sistemin hangi dünyaya ait olduğunu çoğu zaman ele verir: Ku/Ka geniş bant, L/VHF dar-veri/IoT.**

---

<a id="12"></a>
## 12. Uydu-IoT: Küçük-Uydu Mesajlaşma (Swarm, Astrocast, Lacuna, Kineis/Argos)

Mega-konstelasyon devriminin sessiz ama yaygın kolu, **uydu-IoT'dir**: küçük, ucuz uydularla (genelde CubeSat sınıfı) düşük güçlü, düşük veri hacimli, gecikmeye toleranslı mesajlaşma. Hedef, insanın değil **makinenin** bağlanmasıdır: okyanustaki bir şamandıra, çöldeki bir boru hattı sensörü, tundradaki bir hayvan izleme tasması, tarladaki bir toprak nemi ölçeri. Bu cihazlar küçük durum mesajlarını (konum, sıcaklık, seviye) aralıklı olarak gönderir; sürekli geniş bant gerekmez.

**Uydu-IoT'nin tasarım felsefesi.** Geniş bant mega-konstelasyondan (Starlink) temelde farklıdır:

```
   GENİŞ BANT LEO  vs  UYDU-IoT (küçük uydu):

                  Starlink/OneWeb        Swarm/Astrocast/Kineis
   Uydu boyutu    büyük, pahalı          CubeSat, çok küçük, ucuz
   Veri hacmi     yüksek (Mbps-Gbps)     çok düşük (bayt-kilobayt/mesaj)
   Güç (cihaz)    yüksek (terminal+güç)  çok düşük (pil yıllarca)
   Gecikme        düşük, gerçek zamanlı  toleranslı (store-and-forward)
   Kapsama anı    sürekli (binlerce uydu) aralıklı (az uydu, geçişte iletir)
   Anten          faz-dizili panel       küçük yama/tel, yön bağımsız
   Amaç           insan/internet         makine/sensör (M2M)
```

Kilit kavram **store-and-forward (depola-ve-ilet)**: az sayıda uyduyla küresel anlık kapsama mümkün olmadığından, cihaz mesajını gökyüzünden bir uydu geçtiği anda gönderir; uydu mesajı depolar ve bir gateway'in üzerinden geçtiğinde yere indirir. Bu, gerçek zamanlı değildir (dakikalar-saatler gecikme olabilir) ama sensör verisi için genelde yeterlidir ve az uyduyla küresel erişim sağlar. Bölüm 3'teki geçiş tahmini (gpredict) burada doğrudan işe yarar: cihaz, uydunun ne zaman geçeceğini bilirse o anda iletir.

**Başlıca uydu-IoT oyuncuları (2025–2026; teyit edilmeli):**

- **Swarm Technologies** (artık SpaceX bünyesinde): Çok küçük uydularla (sandviç boyutu) ucuz, düşük güçlü IoT bağlantısı sunmuştur. Çok küçük modemler ve düşük maliyetli mesajlaşma ile tanınır.
- **Astrocast:** İsviçre merkezli; iki yönlü (gidiş-dönüş) küçük-uydu IoT konstelasyonu. Cihaz hem gönderir hem komut alır.
- **Lacuna Space:** Karasal IoT'nin yaygın açık standardı olan **LoRaWAN'ı** uzaya taşır; mevcut LoRa cihazlarının uyduya bağlanabilmesini hedefler (Omnispace gibi ortaklıklarla). Açık standart yaklaşımı ayırt edicidir.
- **Kineis (Fransa):** Köklü **Argos** sisteminin (onlarca yıldır çevre/yaban hayatı izleme, deniz şamandıraları, kutup araştırması için kullanılan uydu telemetri sistemi) modern ticari devamı. Argos, özellikle hayvan göçü takibi ve çevresel veri toplamada bir standarttır; Kineis bunu yeni nesil küçük uydularla genişletir. **Argos-4** gibi yeni nesil cihazlar/yükler bu evrimin parçasıdır.
- (Ek olarak Kepler Communications, Myriota, Hiber gibi oyuncular da bu alanda yer alır.)

**Ground station-as-a-service.** Uydu-IoT ekosistemini mümkün kılan bir altyapı katmanı, **yer istasyonu hizmeti** (ground station-as-a-service) modelidir: küçük operatörler kendi yer istasyon ağlarını kurmak yerine, bulut sağlayıcıların ya da uzmanlaşmış şirketlerin dünya çapına dağılmış yer istasyonlarını kiralar. Bu, SatNOGS'un (Bölüm 11/13) ticari/dağıtık mantığının bir uzantısıdır — uydunuz geçerken en yakın kiralık yer istasyonu veriyi alır ve buluta teslim eder. Küçük uydu operatörünün giriş engelini dramatik biçimde düşürür.

**Protokoller.** Uydu-IoT, dar-bant ve düşük güç için optimize protokoller kullanır: bazıları tescilli (Swarm, Astrocast kendi dalga biçimleri), bazıları açık standart (Lacuna'nın LoRaWAN'ı). Ortak özellik, çok kısa paketler, güçlü FEC (Bölüm 9 — zayıf linkte hata düzeltme), düşük sembol hızı ve düşük güç tüketimidir. Bunların çoğu, Iridium SBD'nin (Bölüm 9) açtığı yolun farklı teknik çözümleridir.

Mühendislik sezgisi: **Uydu-IoT, mega-konstelasyonun "az veri, çok cihaz, düşük güç, gecikmeye toleranslı" kardeşidir. Store-and-forward ile az uyduyla küresel erişim sağlar; geçiş tahmini (Bölüm 3) ve yer istasyonu hizmeti onu mümkün kılar; Iridium SBD'den Kineis/Argos'a, tescilliden LoRaWAN'a uzanan bir teknik çeşitlilik gösterir.**

---

<a id="13"></a>
## 13. Doğrudan-Cihaz (D2D / Direct-to-Cell) ve 3GPP NTN

Bu, alanın en güncel ve en çok ilgi çeken kırılmasıdır: **değiştirilmemiş, sıradan bir cep telefonunun doğrudan uyduya bağlanması.** Geleneksel uydu telefonu (Iridium, Thuraya — Bölüm 11) özel, büyük antenli, pahalı bir cihazdır. Doğrudan-cihaz (Direct-to-Device, D2D; cep telefonu özelinde Direct-to-Cell) ise cebindeki standart telefonun, baz istasyonu menzili dışında doğrudan bir uyduya mesaj/ses/veri göndermesini hedefler. Fizik açısından bu, Bölüm 4'teki LEO link bütçesi armağanı (~37 dB) sayesinde sınırda mümkün hale gelir — ama yine de sınırdadır.

**Neden zor, neden yeni mümkün?** Sıradan bir telefonun anteni küçük ve verimsizdir (avuç içine sığar, her yöne yayar), güç çıkışı düşüktür (pil ve SAR/sağlık sınırları). GEO'ya (36.000 km) bu telefonla bağlanmak imkânsızdır (FSPL ~205 dB — Bölüm 8). LEO'ya (~550 km, FSPL ~168 dB — Bölüm 4) bağlanmak ise hâlâ zordur ama mümkündür; çünkü:

1. **Mesafe çok daha az** (~37 dB kazanç, Bölüm 4).
2. **Uydu tarafı çok güçlüdür**: D2D uyduları **devasa antenler** taşır (yüzlerce metrekareye varan açıklıklar), böylece zayıf telefon sinyalini toplayabilir ve telefona güçlü hüzme gönderebilir. Telefonun zayıflığını uydunun gücü telafi eder.
3. **Karasal hücresel spektrum** kullanılır: telefon zaten bu frekansları (örneğin LTE bantları) destekler; "uyduyu havada uçan bir baz istasyonu" gibi görür.

```
   D2D link asimetrisi — telefonun zayıflığını uydu telafi eder:

   ZAYIF taraf (telefon):        GÜÇLÜ taraf (uydu):
   - küçük, verimsiz anten        - DEVASA anten (yüzlerce m²)
   - düşük güç (pil/SAR)          - yüksek güç, hassas alıcı
   - her yöne yayar               - dar, güçlü hüzme
        │                              │
        └──── LEO ~550 km (FSPL ~168 dB) ────┘
        Bağ, ancak uydunun büyük anteni + LEO yakınlığı sayesinde kapanır.
   GEO'da (FSPL ~205 dB, +37 dB) bu bağ sıradan telefonla KAPANMAZ.
```

**Yaklaşımlar ve oyuncular (2025–2026; teyit edilmeli — durum hızla değişir):**

| Yaklaşım | Örnek | Bant/yöntem | Yetenek (yaklaşık) |
|---|---|---|---|
| Özel uydu + mevcut telefon, acil mesaj | **Apple Emergency SOS** (Globalstar altyapısı) | L/S (Globalstar) | Acil durum mesajlaşma (sınırlı) |
| Karasal MNO spektrumu + LEO, "uzaydan baz istasyonu" | **Starlink Direct-to-Cell** (T-Mobile / T-Satellite) | PCS ~1900 MHz (LTE Band 25) | Önce metin (2025 ticari), sonra veri |
| Geniş bant D2D, büyük anten | **AST SpaceMobile** (BlueBird) | MNO spektrumu, çok büyük dizi | Yüksek hıza (broadband) hedefli |
| Aralıklı mesajlaşma | **Lynk Global** | MNO spektrumu | Aralıklı metin (ada/gelişmekte olan pazarlar) |

- **Apple Emergency SOS (Globalstar):** İlk yaygın D2D örneği. Telefon, Globalstar uydularına acil durum mesajı gönderir; kullanıcı telefonu uyduya doğru elle yönlendirir (uydu nerede, ekran gösterir). Yalnızca acil durum — sınırlı ama hayat kurtaran bir hizmet.
- **Starlink Direct-to-Cell (T-Mobile / "T-Satellite"):** Starlink uyduları, T-Mobile'ın **PCS 1900 MHz** (LTE Band 25) karasal spektrumunu kullanarak telefona "uzayda uçan bir kule" gibi hizmet verir; telefonda özel donanım/yazılım gerekmez (çoğu modern telefon Band 25'i destekler). 2025'te ticari olarak önce metin, sonra veri yeteneğiyle başlamıştır; yüzlerce özel D2D uydusu konuşlandırılmıştır. (Teyit edilmeli — kapsama ve yetenek sürekli genişliyor; FCC, belirli spektrum segmentlerinde yüksek güç için düzenleyici onaylar vermiştir.)
- **AST SpaceMobile (BlueBird):** Çok büyük antenli uydularla doğrudan-cihaz **geniş bant** hedefler (yalnızca mesaj değil, daha yüksek hız). Büyük operatör ortaklıkları ile ilerler.
- **Lynk Global:** Daha küçük ölçekte, aralıklı mesajlaşmaya odaklı D2D.

**3GPP NTN (Non-Terrestrial Networks).** Bu işin standart temeli, 3GPP'nin (karasal hücresel standartları belirleyen kuruluş — 4G/5G) **NTN (Non-Terrestrial Networks — karasal-olmayan ağlar)** çalışmasıdır. NTN, hücresel standartları uyduya/yüksek-irtifa platformuna uyarlar; iki ana kol vardır:

```
   3GPP NTN iki kol:

   ┌─────────────────────────────────────────────────────────┐
   │ NB-IoT-NTN   : dar-bant IoT'nin uydu uyarlaması           │
   │                (düşük veri, düşük güç sensörler — Bölüm 12)│
   │ NR-NTN       : 5G NR'nin uydu uyarlaması                  │
   │                (telefon/geniş bant doğrudan-cihaz)        │
   └─────────────────────────────────────────────────────────┘

   NTN'in çözmesi gereken uydu-özgü problemler:
   - Büyük Doppler (Bölüm 3)        → protokolde ön-düzeltme/raporlama
   - Uzun yayılım gecikmesi (LEO)   → zamanlama/HARQ uyarlaması
   - Hızlı handover (Bölüm 6)       → uydu hareketine uyarlı prosedürler
```

NTN'in önemi, **standartlaşmadır**: D2D'nin tescilli, operatöre/uyduya özgü çözümlerden çıkıp, telefon yonga setlerine (Qualcomm, MediaTek gibi) "standart özellik" olarak girmesini sağlar. Böylece gelecekte NTN yeteneği, tıpkı Wi-Fi/Bluetooth gibi, telefonun varsayılan bir özelliği olabilir — operatöre özel anlaşma gerekmeden. 3GPP'nin ardışık sürümleri (Release 17 ilk NTN, sonraki sürümler Ka-bant NTN sinyalizasyonu vb.) bu yolu döşemektedir. (Teyit edilmeli — sürüm içerikleri ve yonga seti desteği gelişiyor.)

**Sınırlar.** D2D sihir değildir; fizik sınırları açıktır:

- **Düşük kapasite:** Tek uydu geniş bir alanı kaplar; o alandaki tüm telefonlar uydunun sınırlı kapasitesini paylaşır. Karasal hücresel ağın yoğunluğu yoktur. Bu yüzden ilk hizmetler metin/acil durum, sonra sınırlı veri — karasal geniş bandın yerini tutmaz.
- **Görüş hattı gerekir:** Uydu sinyali bina içine, derin orman altına, kanyon dibine zor girer (L/S/PCS frekansları engelden geçmede sınırlı). Açık gökyüzü gerekir.
- **Gecikme ve handover:** LEO geçişi ve handover (Bölüm 6) sürekli yönetilmeli; ses için bu, dikkatli mühendislik ister.

Mühendislik sezgisi: **D2D, LEO'nun ~37 dB mesafe armağanını (Bölüm 4) + uydunun devasa antenini birleştirerek sıradan telefonu sınırda uyduya bağlar. Apple SOS'tan Starlink Direct-to-Cell ve AST SpaceMobile'a uzanan oyuncular farklı bant/yetenek noktalarında durur; 3GPP NTN bunu tescilliden standarda taşır. Ama fizik, D2D'yi "her yerde geniş bant" değil "kapsama boşluğunda metin/sınırlı veri/acil durum" olarak konumlandırır.**

---

<a id="14"></a>
## 14. Güvenlik (Savunma Perspektifi): Saldırı Yüzeyi, Kimlik Doğrulama, Gizlilik

Bu bölüm, mega-konstelasyon ve uydu-IoT'yi **savunma gözüyle** inceler: nereler kırılgan, nasıl korunmalı. Bölüm 11'in SATCOM güvenlik tartışmasını (bent-pipe açığı, TT&C, satellite piracy) mega-konstelasyon bağlamına taşır. Amaç saldırı reçetesi değil, savunma için zafiyet haritasıdır (Bölüm 6, Bölüm 20, Bölüm 23 ile bağlantılı).

**Saldırı yüzeyi katmanları.** Bir mega-konstelasyon, GEO yayın uydusundan çok daha geniş bir saldırı yüzeyi sunar; çünkü çok katmanlı bir ağdır:

```
   MEGA-KONSTELASYON saldırı yüzeyi katmanları:

   ┌─────────────────────────────────────────────────────────┐
   │ 1. KULLANICI TERMİNALİ (en erişilebilir)                  │
   │    - faz-dizili modem, ürün yazılımı, yan-kanal           │
   │    - fiziksel ele geçirme, donanım tersine mühendislik    │
   ├─────────────────────────────────────────────────────────┤
   │ 2. RF BAĞI (uplink/downlink)                              │
   │    - jamming (Bölüm 20/23), spoofing, ele geçirme         │
   ├─────────────────────────────────────────────────────────┤
   │ 3. AĞ / GATEWAY                                           │
   │    - yer altyapısı, yönlendirme, internet kenarı          │
   ├─────────────────────────────────────────────────────────┤
   │ 4. UYDU / TT&C (en kritik — Bölüm 11/6)                   │
   │    - komuta bağı ele geçirme = uydu kaybı                 │
   ├─────────────────────────────────────────────────────────┤
   │ 5. ISL / UZAY OMURGASI (en korunaklı — optik, dar)        │
   └─────────────────────────────────────────────────────────┘
```

**Kullanıcı terminali (en geniş yüzey).** Mega-konstelasyonda milyonlarca kullanıcı terminali sahada, fiziksel olarak erişilebilir konumdadır. Her terminal bir modem, ürün yazılımı ve kriptografik kimlik taşır. Saldırı vektörleri: ürün yazılımı çıkarma/tersine mühendislik, yan-kanal analizi (Bölüm 23'teki donanım saldırı mantığı), klonlama girişimleri, fiziksel kurcalama. Savunma: güvenli önyükleme (secure boot), donanım kök güveni, ürün yazılımı imzalama, terminal kimlik doğrulama. Pasif gözlem açısından terminal, uplink yaptığında (Bölüm 11) spektrumda iz bırakır — bu, hem yetkisiz kullanımı tespit edilebilir kılar hem de bir gizlilik boyutu yaratır (aşağıda).

**Jamming ve spoofing direnci/riski.** RF bağı, en klasik saldırı yüzeyidir (Bölüm 20, Bölüm 23):

- **Jamming (sıkışmacılık):** Güçlü bir gürültü vericisiyle kullanıcı bağını boğmak. LEO'nun avantajı, uydunun yakın ve sinyalin görece güçlü olması (Bölüm 4) — bu, jamming'i GEO'ya göre kısmen zorlaştırır; ayrıca faz-dizili anten (Bölüm 5, Bölüm 27) yön ayrımıyla (null steering — jammer yönüne sıfır oluşturma) dirence katkı sağlayabilir. Dezavantaj: yüksek frekans (Ku/Ka) ve geniş bant, hedefli jamming'e açıktır.
- **Spoofing (sahtecilik):** Sahte uydu/baz sinyaliyle terminali kandırmak. Özellikle D2D (Bölüm 13) bağlamında, telefon "havada uçan baz istasyonu" gördüğü için, sahte bir karasal/uydu sinyali kimlik doğrulama zayıfsa risk oluşturur. PNT bağlamında (Bölüm 8) LEO sinyali spoofing'i de bir araştırma/savunma konusudur.

**Uydu-IoT kimlik doğrulama.** Uydu-IoT cihazları (Bölüm 12) küçük, ucuz ve düşük güçlüdür; bu, kriptografik kimlik doğrulamayı zorlaştırır (sınırlı işlem gücü, pil). Zayıf veya eksik kimlik doğrulama, bir saldırganın **sahte sensör verisi enjekte etmesine** (örneğin sahte konum, sahte ölçüm) ya da cihazı taklit etmesine yol açabilir. Savunma: hafif ama güçlü kriptografi (cihaz kimliği imzalama), uçtan uca mesaj bütünlüğü, tekrar saldırısına (replay) karşı zaman damgası/sayaç. Bu, IoT güvenliğinin (Bölüm 23 ile bağlantılı) uydu bağlamındaki yansımasıdır.

**TT&C — en kritik halka (Bölüm 11/6).** Bölüm 11'de vurgulandığı gibi: payload'a (kullanıcı trafiği) yetkisiz erişim kötü, ama **TT&C'ye (komuta bağı) yetkisiz erişim felakettir** — uydunun yönelimini, hüzmelerini, hatta yörüngesini etkileyebilir. Mega-konstelasyonda binlerce uydu olduğu için TT&C güvenliği ölçeklenmeli; tek bir zafiyet binlerce uyduyu etkileyebilir. Savunma: güçlü kimlik doğrulamalı/şifreli komut bağı, komut yetkilendirme, ayrı korunaklı frekanslar. Akademik güvenlik araştırmaları, bazı uydu sistemlerinin (özellikle eski/küçük) TT&C ve bağ güvenliğinde zafiyetler bulmuştur — bu, savunma perspektifinin ana endişesidir.

**GNSS-bağımsız PNT (savunma fırsatı).** Bölüm 8'de işlendiği gibi, LEO mega-konstelasyon sinyalleri **GNSS'e bağımsız bir PNT kaynağı** olma potansiyeli taşır. Bu, bir güvenlik fırsatıdır: GNSS jamming/spoofing'in (Bölüm 20) yaygınlaştığı bir ortamda, güçlü ve bol LEO sinyalleri yedek konum/zaman kaynağı sunabilir. Savunma mimarisi, tek bir PNT kaynağına (GNSS) bağımlılığı azaltmak için LEO sinyallerini tamamlayıcı olarak değerlendirebilir.

**Gözetim ve gizlilik boyutu.** Mega-konstelasyon ve D2D, ciddi bir gizlilik boyutu taşır:

- **Terminal/cihaz konumu:** Bir kullanıcı terminali ya da D2D telefonu uplink yaptığında, varlığı ve kabaca konumu (hangi hüzme/uydu üzerinden) ağ tarafından bilinir. Geniş ölçekte bu, konum meta-verisi üretir.
- **Açık ağ katmanı:** Bazı sistemlerin (Bölüm 10'daki Iridium örneği) açık ağ sinyalizasyonu, pasif gözlemle meta-veri çıkarımına açıktır.
- **Merkezi altyapı:** Mega-konstelasyon, az sayıda operatörün elinde merkezileşmiş, küresel kapsamalı bir altyapıdır; bu, hem dayanıklılık (tek nokta) hem de gözetim/erişim politikası açısından stratejik bir konudur.

Savunma/gizlilik dengesi: kullanıcı içeriğinin uçtan uca şifrelenmesi (içerik gizliliği), meta-verinin minimize edilmesi ve açık sistem katmanının gerektiğinden fazla bilgi sızdırmaması önemlidir. Bu, Bölüm 6'daki genel savunma ve Bölüm 20'deki spektrum-tehdit perspektifinin uydu katmanına taşınmasıdır.

Mühendislik sezgisi: **Mega-konstelasyon, GEO yayın uydusundan çok daha katmanlı bir saldırı yüzeyi sunar — en erişilebilir uçta milyonlarca kullanıcı terminali, en kritik uçta binlerce uydunun TT&C'si. LEO yakınlığı jamming'e karşı kısmi avantaj ve GNSS-bağımsız PNT fırsatı verirken; ölçek, kimlik doğrulama (özellikle düşük güçlü IoT) ve gözetim/gizlilik yeni savunma zorlukları getirir.**

---

<a id="15"></a>
## 15. Alıştırmalar (Yasal, Pasif/RX)

Aşağıdaki alıştırmaların tamamı **pasif alımdır (RX)** ve açık/gözlemlenebilir sinyal katmanına odaklanır. Hiçbiri uplink, terminal aktivasyonu, abonelik/şifreli içerik çözme ya da kullanıcı trafiğine erişim içermez. Her birini yapmadan önce kendi ülkenin mevzuatını teyit et (Bölüm 11 yasal çerçevesi geçerli).

**Alıştırma 1 — Bir LEO uydusunun geçişini Doppler kaymasıyla izle.**
Amaç: Bölüm 3'teki S eğrisini gerçek bir geçişte görmek.
- Sıradan SDR menzilinde, bilinen bir LEO uydusunun (örneğin bir meteoroloji uydusu — Bölüm 11, ya da bir amatör uydu) bilinen taşıyıcı frekansını seç.
- gpredict ile (aşağıda) geçiş zamanını öğren; geçiş sırasında SDR spektrumunu izle.
- Merkez frekansın doğuşta yukarı kaymış, tepe geçişinde nominale döner, batışta aşağı kaymış olduğunu gözlemle. Bu, Bölüm 3'teki Doppler eğrisinin canlı halidir.
- İleri: kayma büyüklüğünü ölçüp Bölüm 3 formülüyle (Δf ≈ f₀·v_radyal/c) radyal hız tahmini yap.

**Alıştırma 2 — Uydu geçiş tahmini (gpredict).**
Amaç: Konstelasyon dinamiğini takvimle kavramak.
- gpredict kur; güncel TLE verisini (Bölüm 11/13) indir.
- İlgilendiğin uyduları (LEO meteoroloji, Iridium, Starlink uyduları) ekle.
- Kendi konumun için önümüzdeki geçişleri, maksimum yükseklik açılarını ve geçiş sürelerini listele. Yüksek açılı geçişlerin neden daha uzun/temiz olduğunu (Bölüm 3) gözlemle.
- İleri: bir Iridium ya da Starlink uydusunun gökyüzünü ne kadar hızlı geçtiğini ölç; Bölüm 3'teki geçiş süresi mertebeleriyle karşılaştır.

**Alıştırma 3 — Iridium ağ sinyallerini gözle (açık ağ katmanı).**
Amaç: Patlamalı TDMA yapısını ve açık ağ duyurularını RF'te görmek (Bölüm 9, Bölüm 10).
- L-bant (~1626 MHz) alabilen bir SDR ve uygun anten (yama/heliks, gökyüzü görüşlü) hazırla.
- gr-iridium ile patlamaları tespit et; spektrumda sürekli taşıyıcı değil, gelip giden kısa patlamalar gördüğünü doğrula (TDMA imzası, Bölüm 9).
- iridium-toolkit ile **açık ağ çerçevelerini** (Ring Alert / IRA tipi sistem duyuruları) ayrıştır. Çizgiyi koru: yalnızca açık ağ/sistem sinyalizasyonu — kullanıcı içeriği DEĞİL (Bölüm 10 yasal/etik sınırı).
- Gözlem: ağ duyurularının ritmi ve uydu geçişlerinin patlama yoğunluğuna etkisi.

**Alıştırma 4 — Starlink kullanıcı downlink'ini spektrumda tanı (uygun LNB/SDR ile).**
Amaç: Mega-konstelasyon geniş bant imzasını ve handover ritmini gözlemek (Bölüm 7, Bölüm 8).
- Ku-bant (~10,7–12,7 GHz) için uygun bir LNB ile sinyali IF'e indir; geniş bant bir SDR ile IF'i gözlemle (Bölüm 8 zinciri).
- Spektrumda ~240 MHz genişliğinde, düz-tepeli OFDM blokları ara (Bölüm 7 imzası). GEO DVB-S2 tek-taşıyıcısından (Bölüm 9) ne kadar farklı göründüğünü karşılaştır.
- Zaman içinde izle: blokların belirip kaybolmasını (uydu geçişi/handover ritmi — Bölüm 6) gözlemle. Yalnızca downlink **varlığını** ve yapısını gözlüyorsun; şifreli içeriği değil.

**Alıştırma 5 — Konstelasyon imza karşılaştırması (masa başı + gözlem).**
Amaç: Frekans bandının sistem türünü ele verdiğini içselleştirmek (Bölüm 11 tablosu).
- Bölüm 11'in alıştırmalarındaki GEO DVB-S2 taşıyıcısı, bu bölümdeki Starlink OFDM bloğu ve Iridium L-bant patlamasını yan yana koy.
- Her birinin bant genişliği, şekli (düz tepe / patlama / sürekli) ve frekans bandını not et; bunlardan sistem türünü (geniş bant LEO / dar-veri LEO / GEO yayın) çıkarsama pratiği yap (Bölüm 10 kimliklendirme mantığı).

---

<a id="16"></a>
## 16. Hızlı Referans ve Diğer Bölümler

**Bu bölümün özü (mega-konstelasyon ve D2D):**

- **Neden LEO:** düşük gecikme (~ms, GEO ~120 ms değil) + yüksek kapasite (frekans yeniden kullanımı, binlerce hüzme) + ucuz seri üretim. Bedeli: sürekli hareket, büyük Doppler, sık handover.
- **LEO link armağanı:** FSPL GEO'ya göre ~37 dB az (550 km vs 36.000 km) → küçük terminal ve D2D mümkün; ama bütçe geçiş boyunca dalgalanır (ACM ile telafi).
- **Faz-dizili terminal (Bölüm 27):** hareketsiz panel, elektronik hüzme izleme; mekanik motor yok; handover için ideal.
- **Starlink:** Ku kullanıcı (DL ~10,7–12,7 GHz / UL ~14–14,5 GHz), Ka gateway (~27,5–30 / 17,8–20,2 GHz), lazer ISL; downlink ~240 MHz OFDM, ~1024 alt-taşıyıcı, çerçeve 1/750 s, PSS/SSS senkronizasyon (açık literatür). PNT potansiyeli.
- **Iridium NEXT:** 66 aktif uydu, ~781 km kutupsal, küresel kapsama; L-bant (~1616–1626,5 MHz) kullanıcı, Ka çapraz-link (4 bağ, halka); FDMA×TDMA patlamalı; SBD ile IoT temeli.
- **Uydu-IoT:** küçük uydu, düşük güç/veri, store-and-forward; Swarm, Astrocast, Lacuna (LoRaWAN), Kineis/Argos; yer istasyonu hizmeti.
- **D2D / Direct-to-Cell:** sıradan telefon → uydu; Apple SOS (Globalstar), Starlink Direct-to-Cell (T-Mobile PCS 1900), AST SpaceMobile, Lynk; 3GPP NTN (NB-IoT-NTN, NR-NTN) standart temeli. Sınır: düşük kapasite, görüş hattı, metin/sınırlı veri.
- **Güvenlik:** katmanlı saldırı yüzeyi (terminal → RF → gateway → TT&C → ISL); jamming/spoofing direnci-riski; IoT kimlik doğrulama zorluğu; GNSS-bağımsız PNT fırsatı; gözetim/gizlilik boyutu.

**Pasif gözlemin değişmez çizgisi:** Downlink **varlığını**, **yapısını** (genişlik, modülasyon, frekans, bant), **geçiş ritmini** ve **açık ağ/sistem sinyalizasyonunu** gözlemek spektrum okuryazarlığıdır ve pasiftir. Şifreli **kullanıcı içeriğini** çözmek, terminali aktive etmek ya da uplink yapmak farklı, izin gerektiren ve çoğu yerde yasadışı bir iştir; bu kitap onu tarif etmez.

**Teyit edilmeli (hızla değişen) konular:** Starlink frekans tahsisleri ve uydu nesilleri; D2D hizmet durumu/kapsama/yetenek; mega-konstelasyon uydu sayıları (Amazon Leo, Telesat, Kuiper→Leo); 3GPP NTN sürüm içerikleri ve yonga seti desteği; konstelasyon isim/sahiplik değişiklikleri. Sayılar değişir, fizik (Doppler, link bütçesi, handover geometrisi) değişmez.

**Diğer bölümlere köprüler:**

- **Bölüm 1** — sinyal/modülasyon temelleri (OFDM, PSK/QAM, FSPL); bu bölümdeki dalga biçimlerinin alt yapısı.
- **Bölüm 2** — SDR donanımı; L-bant/IF gözlemi için alıcı seçimi.
- **Bölüm 3** — anten/devre; L-bant yama/heliks, Ku LNB ön kademe.
- **Bölüm 10** — protokol/sinyal çözümleme yöntemleri; patlama/çerçeve ayrıştırma mantığı.
- **Bölüm 11** — uydu haberleşmesi genel (yörünge sınıfları, Doppler, transponder, link bütçesi, DVB, SatNOGS, GEO/Inmarsat); bu bölüm onun LEO/mega-konstelasyon ileri uzantısıdır.
- **Bölüm 20** — spektrum tehditleri/jamming-spoofing; uydu bağına uygulanışı.
- **Bölüm 22** — mega-konstelasyon giriş; bu bölüm onun derinleştirilmiş devamı.
- **Bölüm 23** — saldırı vektörleri/DoS; terminal ve IoT saldırı yüzeyi, donanım/yan-kanal.
- **Bölüm 27** — anten dizileri ve beamforming; faz-dizili kullanıcı terminalinin (Dishy) ve uydu hüzme oluşturmanın fiziği.

---

> Kapanış: Mega-konstelasyon ve doğrudan-cihaz, Bölüm 11'in tek-GEO-uydu dünyasını alt üst etti. Artık gökyüzü, saatte ~27.000 km hızla uçan binlerce hücreyle döşeli; sıradan bir telefon, kapsama boşluğunda doğrudan uzaya mesaj atabiliyor; bir CubeSat, okyanustaki bir şamandıranın verisini topluyor. Bu bölüm, o devrimin fiziğini (Doppler, link bütçesi, handover), mimarisini (kabuklar, ISL, gateway, faz-dizili terminal), iki amiral gemisini (Starlink, Iridium) ve savunma boyutunu, pasif gözlem çizgisinden çıkmadan ele aldı. Sayılar yarın değişecek; ama bir LEO geçişinin Doppler eğrisini bir kez kendi gözünle gördüğünde, mega-konstelasyonun nasıl nefes aldığını artık bilirsin. Bir sonraki sinyal, gökyüzünden geçerken seni hazır bulsun.
