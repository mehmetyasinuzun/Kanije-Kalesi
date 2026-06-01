# SIGINT EL KİTABI — BÖLÜM 9: YER TESPİTİ, YÖN BULMA VE TAKİP

## Bir Vericinin Konumunu Yalnızca Sinyalinden Kestirmek — Geolocation, DF ve İz Sürme

> Amaç: Önceki bölümler sinyalin fiziğini (Bölüm 1), donanımı (Bölüm 2), anteni ve yön bulmanın temellerini (Bölüm 3), demodülasyonu (Bölüm 4), SIGINT'e karşı savunmayı (Bölüm 6) ve yoğun spektrumda yayıcı ayıklamayı (Bölüm 7) verdi. Bu bölüm tek bir soruya odaklanır: bir verici "ben buradayım" demeden, yalnızca yaydığı elektromanyetik dalgadan onun coğrafi konumunu nasıl kestiririz? Bu, SIGINT zincirinin uzamsal halkasıdır; Bölüm 7'de bir yayıcıyı parametre uzayında bir bulut olarak ayırdık, burada o bulutu haritada bir noktaya (ya da bir hata elipsine) oturtuyoruz. Hedef operatör reçetesi değil, mühendislik sezgisidir: bir DF ekranındaki açı çizgisini, bir multilateration çözümünün hata elipsini ya da bir Kalman izini gördüğünde arkasındaki geometriyi ve matematiği zihninde canlandırabilmen.

> Yasal çerçeve: Bu bölüm de serinin geri kalanı gibi anlama, savunma ve spektrum okuryazarlığı amaçlıdır. Anlatılan yön bulma ve konumlandırma teknikleri tasarım gereği pasif yöntemlerdir; hiçbir iletim, karıştırma ya da yetkisiz içerik çözme önermez. Pasif yön bulmayı yalnızca kendi cihazlarının, açık/yasal sinyallerin (örneğin bir FM yayın istasyonu) ya da lisanslı kendi test vericinin üzerinde uygulamak serbest olabilir; başkasının haberleşmesini izlemek, konumunu kestirmek veya takip etmek çoğu ülkede yasal sınırlara tabidir ve bu bölüm açıkça savunma ve kendi-test perspektifindedir. Bandını, ülkeni ve sürümünü teyit et; Türkiye'de düzenleyici BTK'dır.

---

## İÇİNDEKİLER

1. [Problem: Sinyalden Konuma — Pasif Geolocation'ın Doğası](#1)
2. [Ölçüm Türleri ve Birbirleriyle İlişkisi (AOA, TDOA, FDOA, RSS)](#2)
3. [AOA — Geliş Açısı, Triangulation ve Hata Elipsleri](#3)
4. [TDOA — Varış Zamanı Farkı, Hiperbol ve Multilateration](#4)
5. [FDOA — Doppler Farkı ve Hareketli Platform Geolocation](#5)
6. [RSS — Güç Tabanlı Konumlandırma ve Neden En Zayıf Olduğu](#6)
7. [DF Teknikleri: Watson-Watt'tan Alt-Uzay Algoritmalarına](#7)
8. [Faz-Tutarlı Çok-Kanallı Alıcılar: KrakenSDR ile Pratik DOA](#8)
9. [Pasif Radar (PCL): Aydınlatıcı Fırsatçılığı ve Çapraz-Belirsizlik](#9)
10. [Hareketli Hedef Takibi: Kalman Filtresi ile İz Sürme](#10)
11. [Modern Konumlandırma: WiFi/BLE, Hücresel, EPIRB/Argos](#11)
12. [Hata Kaynakları: Çok-Yol, Geometri (GDOP), Senkronizasyon, Kalibrasyon](#12)
13. [Alıştırmalar (Yasal, Kendi Cihazların ve Kâğıt Üzerinde)](#13)
14. [Hızlı Referans ve Diğer Bölümler](#14)

---

<a id="1"></a>
## 1. Problem: Sinyalden Konuma — Pasif Geolocation'ın Doğası

Bir verici uzaya bir elektromanyetik dalga yayar. Bu dalga, kaynağının konumu hakkında bilgi taşır; ama bu bilgi sinyalin içeriğinde yazılı değildir, sinyalin geometrisinde gizlidir. Pasif geolocation, hiçbir şey yaymadan (yalnızca dinleyerek) bu geometrik ipuçlarını ölçüp kaynağın konumunu geriye doğru çözme problemidir.

Aktif bir sistem (örneğin radar) hedefe enerji gönderip yankısını ölçer; mesafeyi doğrudan elde eder. Pasif bir sistem bu lükse sahip değildir. Yalnızca hedefin kendi yaydığı sinyali görür ve şu dört fiziksel büyüklükten birini (ya da birkaçını) ölçebilir:

- Dalganın hangi yönden geldiği (açı),
- Dalganın farklı alıcılara hangi zaman farklarıyla ulaştığı (zaman),
- Hareket varsa farklı alıcılarda görülen frekans kaymalarının farkı (Doppler),
- Alınan sinyalin gücü (zayıflama).

Her ölçüm türü, hedefin konumu üzerinde bir kısıt (constraint) üretir. Tek bir ölçüm konumu tam belirlemez; yalnızca konumu bir eğri ya da yüzey üzerine sınırlar. İki bağımsız kısıt kesiştiğinde bir konum kestirimi doğar.

```
 Pasif geolocation'ın temel fikri: her ölçüm bir KISIT EĞRİSİ üretir,
 eğrilerin kesişimi konumu verir.

   AOA  → kaynaktan alıcıya bir DOĞRU (yarım-doğru) boyunca kısıt
   TDOA → bir HİPERBOL (iki alıcıdan eşit zaman-farkı yeri) boyunca kısıt
   FDOA → bir eş-Doppler EĞRİSİ boyunca kısıt
   RSS  → bir DAİRE (sabit mesafe = sabit güç) boyunca kısıt

         AOA doğrusu                    TDOA hiperbolü
            ╲                              ╲      ╱
             ╲                              ╲    ╱
              ╲ ● ← kesişim = kaynak         ╲  ╱
               ╲   (iki bağımsız kısıt)       ●  ← kesişim
                ╲                            ╱  ╲
```

Pasif olmanın iki büyük sonucu vardır. Birincisi gizlilik: dinleyen taraf hiçbir enerji yaymadığından kendini ele vermez (Bölüm 6'daki emisyon güvenliği felsefesinin tersten okunuşu). İkincisi belirsizlik: doğrudan mesafe ölçemediğin için tek bir alıcıyla (tek açı dışında) konum bulamazsın; çoğu yöntem birden çok mekânsal olarak ayrık alıcı ya da hareketli bir platform gerektirir. Bu coğrafi çeşitliliğin (baseline) kalitesi, çözümün ne kadar iyi olacağını birinci derecede belirler; bu konuya geometri ve GDOP başlığında (Bölüm 12) döneceğiz.

> Mühendislik sezgisi: Geolocation, "doğrudan ölçemediğin bir büyüklüğü (konum), ölçebildiğin büyüklüklerin (açı/zaman/Doppler/güç) kesişiminden çıkarmaktır". Her ölçüm bir kısıt yüzeyi koyar; problem, yeterince bağımsız yüzeyi kesiştirip belirsizliği bir noktaya (gerçekte küçük bir bölgeye) sıkıştırmaktır. Çözüm asla bir nokta değil, bir olasılık bulutu/hata elipsidir; iyi mühendis konumu değil, konumun belirsizliğini de raporlar.

---

<a id="2"></a>
## 2. Ölçüm Türleri ve Birbirleriyle İlişkisi (AOA, TDOA, FDOA, RSS)

Geolocation'ın dört temel ölçüm türünü ve ürettikleri geometrik kısıtı bir arada görmek, hangisinin ne zaman kullanılacağını anlamanın anahtarıdır. Aşağıdaki tablo bu bölümün haritasıdır; sonraki başlıklar her satırı açar.

| Ölçüm | Açılım | Fiziksel büyüklük | Geometrik kısıt | Tek alıcıyla? | Senkronizasyon ihtiyacı |
|---|---|---|---|---|---|
| AOA / DOA | Angle / Direction of Arrival | Geliş açısı | Kaynağa doğru bir yarım-doğru | Açı evet, konum hayır | Düşük (anten dizisi içi faz) |
| TDOA | Time Difference of Arrival | Varış zamanı farkı | Hiperbol (3B'de hiperboloid) | Hayır (≥2 alıcı çifti) | Çok yüksek (ns düzeyi ortak saat) |
| FDOA | Frequency Difference of Arrival | Doppler farkı | Eş-Doppler eğrisi | Hayır (hareket gerekir) | Yüksek (ortak frekans referansı) |
| RSS | Received Signal Strength | Alınan güç | Daire (sabit mesafe) | Mesafe (kötü), konum hayır | Düşük (kalibrasyon gerekir) |

Bu dört türün önemli bir hiyerarşisi vardır. AOA ve TDOA "iskelet" yöntemlerdir; profesyonel pasif geolocation çoğunlukla bu ikisine ya da kombinasyonuna dayanır. FDOA, hareketli platform varsa TDOA'yı güçlendiren tamamlayıcı bir boyuttur (uydu SIGINT'inde TDOA+FDOA klasik ikilidir). RSS ise en kolay ölçülen ama en güvenilmez büyüklüktür; çok-yol ve ortam değişkenliği güç-mesafe ilişkisini kırdığından, ciddi geolocation'da tek başına nadiren kullanılır, daha çok kapalı-alan/IoT konumlandırmasında (Bölüm 11) yardımcı olur.

```
 Kısıtların boyutu (2B düzlemde):

 AOA:  bir DOĞRU       (1 serbestlik derecesini siler)
       ────────●────────────────→
                  kaynak doğru üstünde bir yerde

 TDOA: bir HİPERBOL    (1 serbestlik derecesini siler, eğrisel)
            ╲          ╱
             ╲        ╱     kaynak hiperbol üstünde bir yerde
              ╲______╱

 RSS:  bir DAİRE       (1 serbestlik derecesini siler, ama yarıçap gürültülü)
              ____
            ╱      ╲
           │   ●R   │   kaynak daire üstünde bir yerde (R çok belirsiz)
            ╲______╱
```

İki boyutlu bir konum bulmak için (enlem, boylam) iki bağımsız kısıt gerekir; örneğin iki AOA doğrusu (triangulation), ya da iki TDOA hiperbolü (üç alıcı), ya da bir AOA + bir TDOA. Üç boyutta (yükseklik dahil) bir kısıt daha gerekir. Pratikte gürültüye karşı dayanıklılık için gereğinden fazla ölçüm toplanır ve problem en küçük kareler (least squares) ile aşırı-belirlenmiş (overdetermined) biçimde çözülür; bu hem konumu hem de hata elipsini verir.

> Not: Hangi yöntemin "en iyi" olduğu evrensel değildir; banda, hedefin sürekliliğine (anlık darbe mi sürekli yayın mı), alıcı sayısına, senkronizasyon imkânına ve geometriye bağlıdır. İyi mühendislik, tek yönteme tapmak yerine eldeki kısıtlara göre yöntem seçmek (ya da füzyonlamak) demektir.

---

<a id="3"></a>
## 3. AOA — Geliş Açısı, Triangulation ve Hata Elipsleri

Geliş açısı (AOA, Angle of Arrival; yön bulmada DOA, Direction of Arrival da denir) ölçümü, dalganın alıcıya hangi yönden geldiğini kestirir. Bölüm 3 yön bulmanın anten/donanım temelini verdi; burada bu açıdan konuma nasıl geçildiğini ve hatanın nasıl modellendiğini ele alıyoruz.

### Faz farkından açıya: dizi geometrisinin matematiği

İki anten elemanı `d` kadar aralıklı olsun ve uzaktaki bir kaynaktan gelen düzlem dalga, dizi normaline `θ` açısıyla gelsin. Dalga, iki elemana farklı yol uzunluklarıyla ulaşır; bu yol farkı:

```
 Δr = d · sin(θ)

        gelen dalga cepheleri
         ╲   ╲   ╲   ╲
          ╲   ╲   ╲   ╲      θ = geliş açısı (normalden)
           ╲   ╲   ╲   ╲
   ─────────●────────●─────────  anten dizisi
            │◄── d ──►│
            │        ╱
            │ Δr=d·sinθ  (ikinci elemana ekstra yol)
```

Bu yol farkı, taşıyıcı dalga boyu `λ` üzerinden bir faz farkına (`Δφ`) dönüşür:

```
 Δφ = (2π / λ) · d · sin(θ)
```

Buradan geliş açısı çözülür:

```
 θ = arcsin( (λ · Δφ) / (2π · d) )
```

Bu, faz-interferometrik yön bulmanın çekirdek denklemidir: ölç ettiğin faz farkı `Δφ`'den açı `θ` çıkar. Önemli bir kısıt: faz `±π` aralığında sarmalanır (wrapping); eğer `d > λ/2` ise aynı faz farkı birden çok açıya karşılık gelebilir ve belirsizlik (ambiguity) doğar. Bunu çözmek için ya elemanlar `λ/2`'den yakın tutulur (tek belirsizlik bölgesi, ama düşük açısal çözünürlük) ya da çok-aralıklı diziler / fazladan elemanlarla belirsizlik giderilir (Bölüm 8'deki KrakenSDR dizilim tartışması bu dengenin pratiğidir).

### Triangulation: iki açı bir konum verir

Tek bir alıcı yalnızca bir açı (bir yarım-doğru) verir; kaynak bu doğru üzerinde bir yerdedir ama mesafe bilinmez. İki coğrafi olarak ayrık alıcı, her biri bir açı ölçerse, iki yarım-doğrunun kesişimi konumu belirler. Buna triangulation denir.

```
 Triangulation (iki DF istasyonu):

   İstasyon 1 ●─────────────────╲ (açı θ1)
   (konumu      ╲                 ╲
    bilinen)      ╲                 ╲
                    ╲                 ● ← KAYNAK (iki açı doğrusunun kesişimi)
                      ╲             ╱
   İstasyon 2 ●────────╲─────────╱ (açı θ2)
   (konumu             ╲       ╱
    bilinen)             ╲   ╱
```

İstasyonların konumu ve ölçtükleri açılar biliniyorsa, kesişim noktası basit trigonometriyle bulunur. İki istasyon arasındaki uzaklığa baseline denir; baseline ne kadar uygun (kaynağa göre geniş açı görecek şekilde) ise kesişim o kadar keskin olur.

### Hata elipsi: neden nokta değil, elips?

Açı ölçümü hatasızsa kesişim bir nokta olurdu. Gerçekte her açıda bir belirsizlik (`±σ_θ`) vardır; bu, her doğruyu bir "ışın demetine" (koni) dönüştürür. İki demetin kesişimi bir nokta değil, bir paralelkenarımsı bölge — istatistiksel olarak bir hata elipsi (error ellipse / uncertainty ellipse) olur.

```
 Açı belirsizliği → hata elipsi:

   İst.1 demeti (θ1 ± σ)           İki demetin kesişim alanı:
        ╲╲╲                              ____
         ╲╲╲                           ╱╱    ╲╲   ← hata elipsi
          ╲╲╲                         (  ●     )    (kaynağın bulunma
   ────────╲╲╲──────────              ╲╲ ____ ╱╱     bölgesi; uzun
   İst.2    ╲╲╲ demeti               '         '     eksen = en zayıf
   (θ2 ± σ)  ╲╲╲                                     yön)
```

Hata elipsinin biçimi geometriye çok bağlıdır. İki açı doğrusu birbirine dik kesişiyorsa elips küçük ve dairesele yakındır (iyi geometri). Doğrular birbirine yakın açıyla (sığ kesişim) kesişiyorsa elips uzar; mesafe yönündeki belirsizlik patlar. Bu, GDOP'un (Bölüm 12) AOA karşılığıdır: "kötü geometri = uzun, ince elips". Belirsizliğin niceliği için sezgisel yaklaşım: kaynak istasyona `R` uzaklıktaysa ve açı hatası `σ_θ` radyan ise, açıya dik yöndeki konum hatası kabaca

```
 σ_yanal ≈ R · σ_θ      (radyan cinsinden)
```

olur. Yani aynı açı hatası, uzaktaki bir kaynakta çok daha büyük konum hatası demektir (`R` ile doğrusal büyür). Bu, AOA'nın yakın hedeflerde güçlü, çok uzak hedeflerde zayıflayan bir yöntem olduğunu açıklar.

| AOA özelliği | Sonuç |
|---|---|
| Tek istasyon bir açı verir | Konum için ≥2 istasyon (triangulation) gerekir |
| Hata `R · σ_θ` ile büyür | Yakın hedefte iyi, uzakta zayıf |
| Geometri elipsi belirler | Dik kesişim iyi, sığ kesişim kötü |
| Kararlı boyut (konum değişmez) | Bölüm 7'de "en kararlı ayraç" denmesinin nedeni |

---

<a id="4"></a>
## 4. TDOA — Varış Zamanı Farkı, Hiperbol ve Multilateration

Varış zamanı farkı (TDOA, Time Difference of Arrival), pasif geolocation'ın belki en güçlü ve en yaygın yöntemidir. Fikir şudur: aynı sinyal, mekânsal olarak ayrık iki alıcıya farklı zamanlarda ulaşır, çünkü kaynaktan alıcılara olan yol uzunlukları farklıdır. Bu zaman farkı, kaynağın konumu üzerinde bir kısıt üretir.

### Mutlak zaman değil, zaman FARKI

Kritik incelik: kaynağın yayın anını bilmiyoruz (pasifiz). Dolayısıyla tek bir alıcıdaki mutlak varış zamanı işe yaramaz; çünkü onu kaynağın bilinmeyen yayın anıyla karşılaştıracak referansımız yok. Ama iki alıcıdaki varış zamanlarının farkını alırsak, kaynağın yayın anı (her ikisi için ortak) sadeleşir; geriye yalnızca yol uzunluğu farkının yarattığı zaman farkı kalır.

İki alıcı `A` ve `B`, kaynak `S` olsun. TDOA:

```
 Δt_AB = t_A − t_B = (|S−A| − |S−B|) / c
```

burada `|S−A|` kaynaktan A'ya mesafe, `c` ışık hızıdır. Bunu mesafe farkı olarak yazarsak:

```
 Δd_AB = c · Δt_AB = |S−A| − |S−B|  =  sabit
```

### Sabit mesafe farkı = hiperbol

Geometrinin güzelliği burada: "iki sabit noktaya (odak) olan mesafelerin FARKI sabit" olan noktaların yeri, tanımı gereği bir hiperboldür. A ve B alıcıları hiperbolün iki odağıdır; ölçülen `Δd_AB` hangi hiperbol kolunu seçeceğini belirler. Kaynak bu hiperbol üzerinde bir yerdedir.

```
 İki alıcı (A, B) = hiperbolün odakları. Sabit Δd = bir hiperbol kolu.

              ╲                       ╱
               ╲  hiperbol kolu      ╱
                ╲ (|SA|−|SB| = Δd)  ╱
                 ╲                 ╱
       A ●········╲···············╱········● B
        (odak)     ╲             ╱      (odak)
                    ╲           ╱
                     ╲         ╱
                      ╲       ╱
        Δd = 0 ise: A-B'nin tam ortasından geçen dikey doğru (dejenere hiperbol)
        Δd ≠ 0 ise: A ya da B'ye yakın tarafa bükülmüş hiperbol kolu
```

### Multilateration: hiperbolleri kesiştirmek

Tek bir alıcı çifti bir hiperbol verir; bu, konumu bir eğriye sınırlar ama tek noktaya indirmez. Üç alıcı (`A`, `B`, `C`) ile iki bağımsız çift (örneğin A-B ve A-C) iki hiperbol üretir; bunların kesişimi konumu belirler. Bu yönteme multilateration (ya da hiperbolik konumlandırma) denir.

```
 Üç alıcı, iki hiperbol → kesişim = konum:

        A ●                    ● B
            ╲ hiperbol A-B    ╱
             ╲               ╱
              ╲      ● ←─────╱── KAYNAK (iki hiperbolün kesişimi)
               ╲    ╱ ╲     ╱
   hiperbol A-C ╲  ╱   ╲   ╱
                 ╲╱     ╲ ╱
                  ●      ╳
                  C    (kesişim)
```

İki hiperbol genelde iki noktada kesişebilir; bu ikilik (ambiguity), üçüncü bir hiperbol (dördüncü alıcı, ya da A-B, A-C, B-C üçlüsünden fazlalık) ya da kaba bir ön bilgi (örneğin AOA, ya da "kaynak şu bölgede") ile giderilir. Dört veya daha çok alıcıyla problem aşırı-belirlenir ve en küçük karelerle hem konum hem hata elipsi çözülür.

### Senkronizasyon: TDOA'nın Aşil topuğu

TDOA'nın bütün gücü, zaman farkını ne kadar hassas ölçebildiğine bağlıdır. Işık hızında, `1 ns` zaman hatası `c · 1ns ≈ 0,3 m` konum hatasıdır. Yani metre düzeyinde konum istiyorsan nanosaniye düzeyinde zaman senkronizasyonu gerekir. Bu, ayrı alıcıların ortak ve son derece kararlı bir saat referansını paylaşmasını zorunlu kılar.

```
 Zaman hatası → konum hatası (TDOA):
 Δt hatası        konum hatası ≈ c·Δt
 ────────────────────────────────────
   1 ns           ≈ 0,30 m
  10 ns           ≈ 3,0 m
 100 ns           ≈ 30 m
   1 µs           ≈ 300 m
```

Pratik çözüm GPS-disiplinli osilatör (GPSDO) kullanmaktır: her alıcı, GPS'in çok kararlı zaman işaretini (1 PPS, saniyede bir darbe ve disipline edilmiş 10 MHz referans) alarak yerel saatini ortak bir mutlak zamana kilitler. Böylece coğrafi olarak ayrık alıcılar nanosaniye düzeyinde aynı saati paylaşır. Sinyalin iki alıcıdaki kopyaları arasındaki zaman farkı ise pratikte çapraz-korelasyon (cross-correlation) ile bulunur: iki kayıt birbirine göre kaydırılır, korelasyonun tepe yaptığı kayma miktarı `Δt`'dir.

```
 Çapraz-korelasyon ile Δt bulma:
 Alıcı A kaydı:  ─────╱╲────────────────
 Alıcı B kaydı:  ──────────╱╲───────────   (aynı sinyal, geç gelmiş)
 Korelasyon R(τ):              ▲ tepe burada → τ_tepe = Δt
                  ───────────┘ └──────────
```

> Uyarı: TDOA'nın doğruluğu sinyalin bant genişliğiyle de yakından ilişkilidir; geniş bantlı (keskin korelasyon tepesi veren) sinyaller zaman farkını dar bantlılardan çok daha hassas verir. Dar bantlı sürekli bir taşıyıcının korelasyon tepesi geniş ve belirsizdir. Bu yüzden TDOA, darbeli ya da geniş bantlı sinyallerde parlar; saf bir CW taşıyıcıda zorlanır. Kesin doğruluk bütçesi banda, SNR'a ve geometriye bağlıdır ve uygulamaya göre teyit edilmelidir.

---

<a id="5"></a>
## 5. FDOA — Doppler Farkı ve Hareketli Platform Geolocation

Varış frekansı farkı (FDOA, Frequency Difference of Arrival; literatürde DFO, Differential Frequency Offset de denir), kaynak ile alıcı(lar) arasında göreli hareket olduğunda ortaya çıkan Doppler kaymalarının farkını kullanır. TDOA'nın "zaman" boyutunun frekans karşılığıdır ve özellikle hareketli platformlarda (uçak, uydu) TDOA'yı tamamlayan ikinci bir bağımsız kısıt verir.

### Doppler hatırlatması

Bir alıcı, kaynağa göre `v_r` bağıl radyal hızıyla (yaklaşma/uzaklaşma hız bileşeni) hareket ediyorsa, gözlediği frekans kayar:

```
 f_gözlenen = f_0 · ( 1 + v_r / c )     →    Δf_Doppler = f_0 · (v_r / c)
```

burada `f_0` gerçek taşıyıcı frekansı, `v_r` radyal hız (yaklaşırken pozitif), `c` ışık hızıdır. Radyal hız, platform hız vektörünün kaynağa doğru izdüşümüdür: `v_r = v · cos(α)`, `α` hız vektörü ile kaynak yönü arasındaki açıdır.

### Frekans farkı = eş-Doppler eğrisi

Pasif olduğumuz için `f_0`'ı (kaynağın gerçek frekansı) tam bilmeyiz; tıpkı TDOA'da yayın anını bilmediğimiz gibi. Çözüm yine FARK almaktır: iki alıcı (ya da iki anda aynı hareketli alıcı) farklı radyal hızlara sahipse, gözledikleri Doppler kaymalarının farkı `f_0`'ın bilinmezliğini büyük ölçüde sadeleştirir ve konum üzerinde bir kısıt verir. Sabit FDOA değeri, düzlemde bir eş-Doppler eğrisi (iso-Doppler contour) çizer; TDOA hiperbolünden farklı şekilli bir eğridir.

```
 İki alıcı çiftinin kısıtları FARKLI şekilli eğriler verir → kesişim daha keskin:

      TDOA hiperbolü                 FDOA eş-Doppler eğrisi
          ╲      ╱                        ~~~~~~~
           ╲    ╱                       ~~       ~~
            ╲  ╱                      ~~           ~~
             ●  ← TDOA+FDOA kesişimi: iki farklı geometrili eğri
            ╱  ╲    aynı noktada kesişince konum çok daha iyi belirlenir
           ~~   ~~
```

### Neden TDOA+FDOA ikilisi güçlüdür

Uydu/hava SIGINT'inin klasik reçetesi TDOA+FDOA birlikteliğidir. Nedeni şudur: TDOA bir hiperbol, FDOA bir eş-Doppler eğrisi verir ve bu iki eğri farklı geometriye sahiptir; aynı noktada kesiştiklerinde, tek tip iki kısıttan (örneğin iki hiperbol, sığ açıyla kesişen) elde edilenden çok daha keskin bir konum çözümü doğar. İki uydu (ya da hareketli iki platform) tek bir kaynaktan hem zaman farkını hem Doppler farkını ölçerse, tek bir geçişte bile konum kestirilebilir. FDOA, hareket olmadan (statik alıcılar, statik kaynak) sıfırdır ve bilgi vermez; bu yüzden FDOA platform hareketinin armağanıdır.

> Mühendislik sezgisi: TDOA "yol uzunluğu farkını" (geometri), FDOA "yol uzunluğunun değişim hızı farkını" (geometrinin türevi) ölçer. Biri konumu, diğeri konumun harekete göre nasıl değiştiğini kısıtlar. İkisi bağımsız bilgi taşıdığından, birlikte tek bir platform-çiftiyle bile sağlam çözüm verirler. FDOA'nın bedeli, alıcılar arasında çok kararlı bir ortak frekans referansı (yine GPSDO sınıfı) gerektirmesidir; Doppler farkları çoğu zaman Hz'in altındadır ve yerel osilatör kararsızlığı bu ince farkı yutar.

---

<a id="6"></a>
## 6. RSS — Güç Tabanlı Konumlandırma ve Neden En Zayıf Olduğu

Alınan sinyal gücü (RSS, Received Signal Strength), ölçmesi en kolay büyüklüktür; her alıcı zaten sinyal seviyesini bilir. Fikir basittir: sinyal kaynaktan uzaklaştıkça zayıflar; alınan güçten mesafeyi, birden çok alıcının mesafesinden konumu kestirebiliriz. Ne yazık ki bu yöntem, basitliğine karşın pasif geolocation'ın en güvenilmez üyesidir.

### Güç-mesafe ilişkisi (log-distance modeli)

Boş uzayda güç, mesafenin karesiyle azalır (FSPL, Bölüm 1). Gerçek ortamda zayıflama daha hızlıdır ve log-distance path loss modeliyle ifade edilir:

```
 P(d) [dBm] = P(d_0) − 10 · n · log10( d / d_0 )
```

burada `P(d)` mesafe `d`'deki güç, `d_0` referans mesafe, `n` yol-kaybı üssüdür (path loss exponent). Boş uzayda `n = 2`; bina içinde, şehirde, ormanda `n` 2,5 ile 5+ arasında değişir. Mesafeyi çözersek:

```
 d = d_0 · 10^( (P(d_0) − P(d)) / (10·n) )
```

### Neden en zayıf yöntem

Sorun bu formülün her teriminin kırılgan olmasıdır:

```
 RSS daire kısıtı (her alıcı bir daire) — ama yarıçaplar çok belirsiz:

   Alıcı 1 ◯ (d1 ± büyük belirsizlik)
        ╲   ___
         ╲╱     ╲      Daireler kalın "halka"lar gibi; kesişim
        ( ● ?    )     geniş bir bulanık bölge → konum çok kaba
         ╲ ___ ╱
   Alıcı 2 ◯ (d2 ± büyük belirsizlik)
```

| Sorun | Etkisi |
|---|---|
| `n` ortamla değişir ve önceden bilinmez | Mesafe kestirimi sistematik olarak kayar |
| Çok-yol (fading) gücü saniyeler içinde ±20 dB oynatabilir | Aynı mesafe, çok farklı güç → mesafe gürültülü |
| Verici çıkış gücü/anten kazancı bilinmeyebilir | `P(d_0)` referansı belirsiz → tüm ölçek kayar |
| Gölgeleme (binalar, arazi) | Tutarsız, yönsel zayıflama |

Güç bir mesafeye, mesafe ise mesafenin karesi/üssü üzerinden değiştiğinden, güçteki küçük bir hata mesafede büyük bir hataya patlar. Çok-yollu bir ortamda alınan güç, kaynağın mesafesinden çok o anki sönümlenme (fading) durumunu yansıtabilir. Bu yüzden ciddi geolocation'da RSS tek başına neredeyse hiç kullanılmaz; ancak kontrollü, kalibre edilmiş ortamlarda (kapalı alan WiFi/BLE konumlandırması, Bölüm 11) ve genelde bir parmak izi (fingerprinting) veritabanıyla birlikte işe yarar; orada bile RTT/AOA gibi başka boyutlarla desteklenir.

> Pratikte: RSS'i "kaba bir yakınlık göstergesi" olarak düşün, "mesafe ölçer" olarak değil. "Sinyal güçleniyor, kaynağa yaklaşıyoruz" demek genelde doğrudur (foxhunt/yön avı sezgisi); ama "güç −60 dBm, demek ki tam 120 metre" demek çok-yollu bir ortamda yanıltıcıdır. Güç, yön/zaman/Doppler kadar saygın bir geometrik kısıt değildir.

---

<a id="7"></a>
## 7. DF Teknikleri: Watson-Watt'tan Alt-Uzay Algoritmalarına

Yön bulma (DF, Direction Finding), AOA ölçümünün pratik mühendisliğidir. Bölüm 3 anten/donanım temelini verdi; burada DF tekniklerini, en eski analog yöntemlerden modern yüksek-çözünürlüklü alt-uzay algoritmalarına kadar bir ölçek üzerinde sıralıyoruz. Hepsinin ortak amacı aynıdır: gelen dalganın yönünü kestirmek; fark, açısal çözünürlükte ve çok-kaynak ayırmadaki güçte.

### Watson-Watt / Adcock

Klasik ve olgun bir tekniktir. Adcock anten dizilimi (genelde dik yönlerde yerleştirilmiş dipoller, örneğin kuzey-güney ve doğu-batı çiftleri) gelen dalganın iki dik bileşenini (sinüs ve kosinüs benzeri) algılar. Watson-Watt yöntemi bu iki bileşeni karşılaştırarak (genelde bir görüntüleme/işleme ile) geliş açısını çıkarır. Avantajı sadeliği ve geniş bant uyumudur; tek anlık ölçümle (instantaneous) açı verir. Dezavantajı düşük çözünürlüğü ve çok-yollu/çok-kaynaklı ortamda zayıflamasıdır; iki kaynak aynı anda gelirse açıyı karıştırır.

### Doppler / Pseudo-Doppler DF

Bir anteni kaynağın etrafında fiziksel olarak döndürdüğünü düşün; antenin kaynağa yaklaşıp uzaklaşması, alınan sinyalde periyodik bir Doppler kayması yaratır ve bu kaymanın fazı geliş yönünü kodlar. Fiziksel döndürme pratik olmadığından, pseudo-Doppler bir daire üzerine dizilmiş birden çok anteni hızlıca elektronik olarak sırayla seçerek (komütasyon) dönmeyi taklit eder; oluşan yapay Doppler modülasyonunun fazından açı çıkar. Sade ve ucuzdur, amatör DF setlerinde yaygındır; çözünürlüğü orta düzeydedir ve gürültüye/çok-yola duyarlıdır.

```
 Pseudo-Doppler: daire üzerinde antenleri sırayla seçmek = dönen anteni taklit

         A1
       ╱    ╲
     A4  ◌   A2      antenler sırayla A1→A2→A3→A4→A1... seçilir;
       ╲    ╱        gelen dalga yönüne göre oluşan yapay Doppler'in
         A3          FAZI, geliş açısını verir
```

### Faz interferometri

Bölüm 3'teki temel: iki (veya daha çok) anten arasındaki faz farkını doğrudan ölçüp `θ = arcsin(λΔφ / 2πd)` ile açıya çevirmek. Faz-tutarlı çok-kanallı alıcı gerektirir (her kanal aynı saat/LO ile çalışmalı). Yüksek doğruluk verir; belirsizlik (`d > λ/2`) sorununu çok-aralıklı dizilerle çözer. Modern sayısal DF'in (KrakenSDR dahil) temel mekanizmalarından biridir.

### Alt-uzay algoritmaları: MUSIC ve ESPRIT

Buraya kadarki yöntemler bir tür "tek açı" mantığıyla çalışır ve çok kaynağı ayırmada zorlanır. Alt-uzay (subspace) algoritmaları niteliksel bir sıçramadır: aynı bantta, aynı anda gelen birden çok kaynağı ayrı ayrı çözebilir ve süper-çözünürlük (super-resolution) sunar, yani anten açıklığının klasik (Rayleigh) çözünürlük sınırından daha keskin açı ayırır.

Temel fikir (kavramsal): Çok-kanallı alıcıdan gelen örneklerin kovaryans matrisi (kanallar arası korelasyonu özetleyen matris) hesaplanır. Bu matrisin özdeğer ayrışımı (eigen-decomposition), uzayı iki dik alt-uzaya böler: sinyal alt-uzayı (gelen kaynakların yaydığı, güçlü özdeğerlere karşılık gelen) ve gürültü alt-uzayı (geri kalan, zayıf özdeğerler). MUSIC (MUltiple SIgnal Classification) şu gözlemi kullanır: gerçek bir kaynağın yön-vektörü (steering vector — o açıdan gelen dalganın dizide yaratacağı faz deseni) sinyal alt-uzayında yatar, dolayısıyla gürültü alt-uzayına diktir. MUSIC, tüm olası açılar için yön-vektörünün gürültü alt-uzayına ne kadar dik olduğunu ölçen bir fonksiyon (MUSIC spektrumu) hesaplar; bu fonksiyon gerçek geliş açılarında keskin tepeler verir.

```
 MUSIC spektrumu (kavramsal): gerçek açılarda keskin tepeler

 P_MUSIC ▲        ▌                    ▌
         │        ▌                    ▌      ← iki kaynak, iki keskin tepe
         │        ▌                    ▌        (klasik yöntem bunları tek
         │     ___▌___              ___▌___      geniş tümsek görürdü)
         └────┴───────┴────────────┴───────┴──► geliş açısı θ
                 θ1                    θ2
```

MUSIC'in yüksek çözünürlüğünün nedeni, açıyı bir genlik/güç tepesi olarak değil, "gürültü alt-uzayına diklik" keskinliği olarak aramasıdır; diklik koşulu çok keskin (paydası sıfıra giden) bir tepe üretir. ESPRIT (Estimation of Signal Parameters via Rotational Invariance Techniques) aynı alt-uzay fikrini kullanır ama açıyı tüm açıları taramadan (spektrum hesaplamadan), dizinin özel yapısından (öteleme değişmezliği) doğrudan cebirsel olarak çözer; daha hesap-verimlidir.

| DF tekniği | Çözünürlük | Çok-kaynak ayırma | Karmaşıklık | Tipik kullanım |
|---|---|---|---|---|
| Watson-Watt / Adcock | Düşük | Zayıf | Düşük | Klasik, geniş bant, hızlı |
| Pseudo-Doppler | Orta | Zayıf | Düşük | Amatör DF, foxhunt |
| Faz interferometri | Yüksek | Orta | Orta | Sayısal çok-kanallı DF |
| MUSIC / ESPRIT | Çok yüksek (süper-çöz.) | Güçlü | Yüksek | Yoğun/çok-yollu ortam |

> Uyarı: MUSIC/ESPRIT'in performansı, kaynak sayısının doğru kestirilmesine, kanalların iyi kalibre edilmesine ve yeterli örnek toplanmasına bağlıdır; korele (örneğin çok-yoldan gelen kendiyle korele) kaynaklar temel MUSIC'i bozar ve uzamsal yumuşatma (spatial smoothing) gibi ek işlemler gerektirir. Buradaki anlatım kavramsaldır; kesin uygulama koşulları ve sınırlar kaynaktan teyit edilmelidir.

---

<a id="8"></a>
## 8. Faz-Tutarlı Çok-Kanallı Alıcılar: KrakenSDR ile Pratik DOA

Yukarıdaki faz-tabanlı ve alt-uzay teknikleri, kâğıtta güzeldir; pratikte ise faz-tutarlı (phase-coherent) çok-kanallı bir alıcı gerektirir. "Faz-tutarlı" demek, tüm alıcı kanallarının aynı yerel osilatör (LO) ve aynı saatle çalışması, böylece kanallar arasındaki faz ilişkisinin yalnızca antene gelen dalganın geometrisinden kaynaklanması (alıcının kendi rastgele fazından değil) demektir. Bu koşul sağlanmazsa faz farkı ölçümü anlamsızdır.

KrakenSDR, bu işi erişilebilir kılan, beş adet RTL-SDR tabanlı kanalı ortak bir saat ve LO ile faz-tutarlı çalıştıran bir cihazdır. Beş kanal, beş anten elemanından gelen sinyali eşzamanlı örnekler; aralarındaki faz farklarından (yukarıdaki faz interferometri / MUSIC matematiğiyle) geliş açısı (DOA) hesaplanır.

### Kalibrasyon: neden zorunlu

Beş kanalın donanım yolları (kablolar, RTL parçaları) tıpatıp aynı değildir; her kanal sinyale küçük, bilinmeyen bir faz ve kazanç kayması ekler. Bu kanal-içi farklar, antenden gelen gerçek geometrik faz farkına karışır ve açıyı bozar. Bu yüzden ölçümden önce kalibrasyon şarttır: bilinen bir referans sinyali (KrakenSDR'de dahili bir gürültü kaynağı / koherens referansı) tüm kanallara aynı anda verilir; kanallar arasındaki ölçülen faz/kazanç farkları, "antenden gelmeyen, donanımdan gelen" kısım olarak kaydedilir ve sonraki ölçümlerden çıkarılır (kalibre edilir). Kalibrasyon olmadan DOA çıktısı sistematik olarak yanlıştır.

### Anten dizilimi: UCA mı ULA mı

Beş anteni nasıl yerleştirdiğin, çözebileceğin açı aralığını ve belirsizliği belirler.

```
 ULA (Uniform Linear Array — düz dizi):   UCA (Uniform Circular Array — dairesel):

   ●   ●   ●   ●   ●                              ●
   │←d→│                                       ╱     ╲
   açısal kapsama ~180° (ön/arka              ●   ◌   ●     açısal kapsama 360°
   belirsizliği var: soldan mı                 ╲     ╱      (her yönü ayırır,
   sağdan mı geldiği karışabilir)                 ●         ön/arka belirsizliği yok)
```

Düz dizi (ULA) daha yüksek tek-eksen çözünürlük verir ama ön-arka belirsizliği taşır (dizinin önünden mi arkasından mı geldiği ayrılamaz; ~180° kapsama). Dairesel dizi (UCA) 360° kapsama sunar (her yönü ayırır) ve ön-arka belirsizliği yoktur; bu yüzden bilinmeyen yönden gelen sinyali aramada (örneğin mobil DF) UCA tercih edilir. Eleman aralığı `d` yine `λ/2` dengesindedir: yakın aralık belirsizliği azaltır ama çözünürlüğü düşürür, geniş aralık tersini yapar. Hedef frekansa göre anten boyu ve aralığı ayarlanmalıdır (Bölüm 3'teki rezonans/boy ilişkisi).

### Tek istasyon yön verir, konum için hareket/çoklu istasyon gerekir

Tek bir KrakenSDR yalnızca bir açı (bir yarım-doğru) verir; konumu belirlemek için ya iki ayrı sabit istasyonun açılarını kesiştirmek (triangulation) ya da tek istasyonu hareket ettirip (örneğin araçla) farklı konumlardan alınan açıları birleştirmek gerekir. Hareketli tek-istasyon DF'inde, art arda alınan açı çizgilerinin kesişimi (ya da bir filtre ile birleştirilmesi) konuma yakınsar; bu, foxhunt/yön avı pratiğinin sayısallaşmış halidir.

| KrakenSDR pratik konu | Dikkat noktası |
|---|---|
| Faz-tutarlılık | 5 kanal ortak saat/LO; şart |
| Kalibrasyon | Her ölçüm öncesi/periyodik; gürültü referansıyla |
| Dizilim | UCA (360°, ön-arka belirsizliği yok) çoğu mobil senaryoda tercih |
| Eleman aralığı | `λ/2` civarı; belirsizlik–çözünürlük dengesi |
| Konum | Tek istasyon açı verir; triangulation ya da hareket gerekir |

> Pratikte: KrakenSDR ile yasal bir FM istasyonunun yönünü bulmak, bütün bu teori için mükemmel bir alıştırmadır (Bölüm 13). FM yayını güçlü, sürekli ve yasal olarak dinlenebilir; DOA çıktısının istasyonun gerçek yönüne oturmasını gözlemlemek, kalibrasyonun ve dizilimin etkisini somutlaştırır. İçeriği çözmüyor, yalnızca yönü ölçüyorsun; bu pasif ve eğitseldir.

---

<a id="9"></a>
## 9. Pasif Radar (PCL): Aydınlatıcı Fırsatçılığı ve Çapraz-Belirsizlik

Şimdiye kadar hedefin kendi yaydığı sinyali kullandık. Pasif radar (PCL, Passive Coherent Location; ya da pasif bistatik radar), bir adım farklı ve zarif bir fikirdir: hedef hiçbir şey yaymasa bile, çevredeki mevcut yayınların (FM radyo, DVB-T sayısal televizyon, baz istasyonu sinyalleri) hedeften yansımasını dinleyerek hedefi tespit ve konumlandırır. Kendi vericin yok; başkasının zaten havada olan sinyalini "aydınlatıcı" (illuminator of opportunity) olarak ödünç alırsın.

### Geometri: iki yol, bir fark

Sistemde iki alım kanalı vardır:

- Referans kanalı: doğrudan aydınlatıcıdan (örneğin FM vericisinden) gelen sinyal.
- Gözlem (surveillance) kanalı: hedeften yansıyıp gelen, gecikmiş ve Doppler-kaymış sinyal.

```
                       HEDEF (uçak vb.)
                        ╱  ╲
       yansıyan ╱╱╱╱╱╱╱╱    ╲╲╲╲╲ gözlem kanalına
       (gecikmiş+Doppler)        ╲╲╲╲
   AYDINLATICI ●─────doğrudan─────────► ● ALICI
   (FM/DVB-T vericisi)  (referans kanalı)   (2 kanal: referans + gözlem)
```

Hedeften yansıyan sinyal, doğrudan sinyale göre fazladan yol katettiği için gecikmiştir (bistatik menzil) ve hedef hareketliyse Doppler kaymıştır (bistatik hız). Bu gecikme ve Doppler, hedefin konumu ve hızı hakkında bilgi taşır.

### Çapraz-belirsizlik fonksiyonu (CAF)

Hedefin gecikmesini ve Doppler'ini birlikte çıkarmak için referans ile gözlem kanalı arasında çapraz-belirsizlik fonksiyonu (CAF, Cross-Ambiguity Function) hesaplanır. Kavramsal olarak CAF, "referans sinyalini her olası gecikme `τ` kadar kaydır ve her olası Doppler `f_d` kadar frekans-ötele, sonra gözlem sinyaliyle korelasyonuna bak" demektir:

```
 CAF(τ, f_d) = ∫ s_gözlem(t) · s*_referans(t − τ) · e^(−j2π f_d t) dt
```

Bu fonksiyon, gerçek hedefin `(τ, f_d)` değerinde bir tepe verir. CAF'i bir gecikme-Doppler düzleminde bir yüzey olarak düşün; tepeler hedefleri, tepenin koordinatları o hedefin bistatik menzilini ve hızını gösterir.

```
 CAF yüzeyi (gecikme-Doppler düzlemi):

 Doppler f_d ▲          ▲ ← hedef 1 tepesi (τ1, fd1)
             │
             │                      ▲ ← hedef 2 tepesi (τ2, fd2)
             │   (taban: doğrudan sinyal ve gürültü)
             └──────────────────────────────► gecikme τ
```

Pratikte en zorlu adım, çok güçlü doğrudan sinyalin (referans) gözlem kanalına sızıntısını ve sabit yansımaları (clutter) bastırmaktır; zayıf hedef yansıması, güçlü doğrudan sinyalin yanında kolayca kaybolur. Bu yüzden PCL sistemlerinin kalbi, doğrudan-sinyal ve clutter bastırma (adaptif iptal) algoritmalarıdır. Aydınlatıcı seçimi de önemlidir: DVB-T gibi geniş bantlı sayısal sinyaller, FM'e göre daha keskin (daha iyi menzil çözünürlüklü) bir CAF tepesi verir, çünkü bant genişliği menzil çözünürlüğünü belirler.

| PCL bileşeni | Rolü |
|---|---|
| Aydınlatıcı (FM/DVB-T/baz istasyonu) | Bedava, hep havada olan enerji kaynağı |
| Referans kanalı | Doğrudan sinyal (karşılaştırma şablonu) |
| Gözlem kanalı | Hedeften yansıyan gecikmiş/Doppler sinyal |
| CAF | Gecikme+Doppler'i tepe olarak çıkarır |
| Clutter/doğrudan-sinyal bastırma | Zayıf hedefi güçlü gürültüden kurtarır |

> Mühendislik sezgisi: Pasif radar, "kendi enerjini harcamadan başkasının enerjisiyle görmek"tir. Gizlidir (yayın yapmaz), ucuzdur (verici gerektirmez) ve mevcut altyapıyı kullanır; bedeli, fırsatçı aydınlatıcının dalga biçimini kontrol edememen ve zayıf yansımayı güçlü doğrudan sinyalden ayıracak ağır işaret işlemesidir. CAF, bu bölümdeki TDOA (gecikme) ve FDOA (Doppler) fikirlerinin tek bir iki-boyutlu yüzeyde birleşmiş halidir.

---

<a id="10"></a>
## 10. Hareketli Hedef Takibi: Kalman Filtresi ile İz Sürme

Tek bir konum kestirimi anlık bir fotoğraftır. Hedef hareketliyse (uçak, gemi, araç), art arda alınan gürültülü konum ölçümlerini birleştirip pürüzsüz, sürekli bir iz (track) oluşturmak ve gelecekteki konumu öngörmek gerekir. Bunun standart aracı Kalman filtresidir.

### Temel fikir: öngör ve düzelt

Kalman filtresi, hedefin durumunu (state) tutar; tipik olarak konum ve hız (`[x, y, vx, vy]`). İki adımı döngüsel tekrarlar:

```
 ┌─────────────┐   yeni ölçüm gelir    ┌──────────────┐
 │  ÖNGÖR       │ ──────────────────▶  │  DÜZELT       │
 │ (predict)    │                       │ (update)      │
 │ hareket      │ ◀──────────────────  │ ölçüm ile     │
 │ modeliyle    │   düzeltilmiş durum   │ düzelt        │
 │ ileri taşı   │   bir sonraki adıma   │ (Kalman kazancı)
 └─────────────┘                       └──────────────┘
```

Öngör adımı: bir hareket modeli (örneğin "sabit hızla gidiyor") kullanarak mevcut durumu bir sonraki ana taşır; "ölçüm gelmeseydi hedef nerede olurdu" sorusunu yanıtlar. Bu sırada belirsizlik (kovaryans) büyür, çünkü öngörü her geçen anda biraz daha güvenilmezleşir.

Düzelt adımı: yeni (gürültülü) bir konum ölçümü geldiğinde, öngörü ile ölçümü ağırlıklı ortalar. Ağırlığı belirleyen Kalman kazancı (Kalman gain), "öngörüye mi yoksa yeni ölçüme mi daha çok güveneyim?" sorusunu, ikisinin belirsizliklerine bakarak çözer. Ölçüm çok gürültülüyse öngörüye, model güvenilmezse ölçüme yaslanır.

```
 Öngörü vs ölçüm füzyonu:
   öngörülen konum  ●----------●  ölçülen konum (gürültülü)
                       ▲
                       │ Kalman kazancı, ikisi arasında
                       ● ← düzeltilmiş (filtrelenmiş) konum
                         (belirsizliği ikisinden de küçük)
```

Kalman filtresinin zarafeti, çıktının belirsizliğinin (kovaryansın) hem öngörüden hem ölçümden daha küçük olmasıdır: iki bilgi kaynağını birleştirmek, her birinden ayrı ayrı daha iyi bir kestirim verir. Ayrıca filtre, ölçüm gelmediği anlarda bile öngörü adımıyla hedefin nerede olacağını tahmin eder; bu, geçici sinyal kayıplarını (hedef bir an susarsa) köprülemeyi sağlar.

Not: Hareket doğrusal değilse (manevra yapan hedef, açı ölçümünün konuma doğrusal-olmayan bağı) temel (doğrusal) Kalman yetmez; genişletilmiş Kalman filtresi (EKF) ya da unscented Kalman filtresi (UKF) gibi türevler kullanılır. Bunlar doğrusal-olmayan model/ölçüm fonksiyonlarını yerel olarak doğrusallaştırarak (EKF) ya da örnek noktalarla (UKF) aynı öngör-düzelt iskeletini sürdürür.

### İz birleştirme (track association / fusion)

Çoklu hedef ve çoklu sensör varsa iki ek problem doğar. Birincisi veri ilişkilendirme (data association): bu yeni ölçüm hangi mevcut ize ait? Yanlış eşleştirme izleri karıştırır (Bölüm 7'deki deinterleaving'in uzamsal-zamansal kardeşi: "hangi gözlem hangi nesneye ait?"). İkincisi iz birleştirme (track fusion): farklı sensörlerin (örneğin bir DF istasyonu + bir TDOA ağı) ürettiği izleri tek bir tutarlı dünya resmine kaynaştırmak. Bu, TCPED zincirinin (Bölüm 7) Exploitation/füzyon adımının uzamsal yüzüdür.

| Takip kavramı | Özü |
|---|---|
| Durum (state) | Genelde konum + hız `[x,y,vx,vy]` |
| Öngör (predict) | Hareket modeliyle ileri taşı, belirsizlik büyür |
| Düzelt (update) | Ölçümle füzyon; Kalman kazancı ağırlıklar |
| EKF / UKF | Doğrusal-olmayan hareket/ölçüm için türevler |
| Veri ilişkilendirme | "Bu ölçüm hangi ize ait?" |
| İz birleştirme | Çoklu sensör izlerini tek resme kaynaştırma |

> Mühendislik sezgisi: Tek bir geolocation çözümü "hedef şu an muhtemelen burada" der; Kalman filtresi bunu "hedef nereden geldi, şimdi nerede, birazdan nerede olacak" hikâyesine dönüştürür. Geolocation noktayı, takip ise noktanın zaman içindeki yörüngesini ve onun belirsizliğini verir. İyi takip, gürültülü tekil ölçümleri, fizik (hareket modeli) süzgecinden geçirip tutarlı bir ize dönüştürmektir.

---

<a id="11"></a>
## 11. Modern Konumlandırma: WiFi/BLE, Hücresel, EPIRB/Argos

Bu bölümün yöntemleri (AOA, TDOA, FDOA, RSS) yalnızca askeri/ELINT bağlamında değil, gündelik teknolojinin içinde de yaşar. Aynı fizik, farklı isimlerle her yerdedir; bunları tanımak hem savunma (kendi cihazının nasıl konumlandığını/sızdığını bilmek) hem de sezgi açısından öğreticidir.

### WiFi / BLE: RSSI + RTT

İç mekân konumlandırması (mağaza, havaalanı, ofis) sıkça WiFi ve Bluetooth Low Energy (BLE) kullanır.

- RSSI tabanlı: Cihazın çeşitli erişim noktalarından/işaretçilerden (beacon) aldığı sinyal gücü (RSSI), bir parmak izi (fingerprint) haritasıyla eşleştirilerek konum kestirilir. Bu, Bölüm 6'daki RSS'in kapalı-alan uygulamasıdır; kalibre edilmiş bir radyo haritası gerektirir ve çok-yola karşı kırılgandır.
- RTT (Round-Trip Time) tabanlı: Modern WiFi (özellikle 802.11mc Fine Timing Measurement, FTM) sinyalin gidiş-dönüş süresini ölçerek mesafeyi doğrudan kestirir; bu, RSS'ten çok daha güvenilir bir mesafe verir (zaman tabanlı, güç tabanlı değil). BLE tarafında ise yön bulma (AoA/AoD) özellikleri (çok-antenli işaretçiler) açı bilgisi sunar.

Önemli savunma notu: BLE/WiFi yayını yapan her cihaz (telefon, kulaklık, etiket) bu yöntemlerle konumlandırılabilir; rastgele/dönüşümlü MAC adresi gibi önlemler tam da bu pasif izlemeyi zorlaştırmak içindir (Bölüm 5'teki cihaz sızıntı yüzeyiyle doğrudan ilişkili).

### Hücresel: OTDOA ve Hücre + TA

Hücresel ağlar konumlandırmayı standart olarak içerir.

- Hücre kimliği (Cell-ID) + TA: En kaba yöntem, cihazın bağlı olduğu baz istasyonunu (hücre) bilmektir; konumu o hücrenin kapsama alanına sınırlar. Zaman hizalama (TA, Timing Advance) değeri, cihazın baz istasyonuna olan kaba mesafesini (zaman tabanlı bir halka) ekler; hücre + TA, kabaca bir yay/halka bölgesi verir.
- OTDOA (Observed Time Difference of Arrival): LTE'nin hassas yöntemidir ve tam olarak bu bölümün TDOA'sıdır: cihaz, birden çok baz istasyonundan gelen referans sinyallerin (PRS) varış zaman farklarını ölçer; bu farklar hiperboller üretir, kesişimleri konumu verir. 5G'de bu aile genişler (çok-antenli AOA/AOD dahil çeşitli yöntemler).

```
 Hücresel konumlandırma kabalıktan hassasa:
 Cell-ID         → hangi hücre (geniş bölge)
 Cell-ID + TA    → hücre + mesafe halkası (yay bölge)
 OTDOA (TDOA)    → çok baz istasyonu hiperbol kesişimi (hassas)
```

### EPIRB / Argos: Doppler tabanlı kurtarma lokasyonu

Tarihsel olarak en zarif geolocation örneklerinden biri, uydu tabanlı acil konum belirlemedir. EPIRB (denizcilik), ELT (havacılık) ve PLB (kişisel) acil işaretçileri ile Argos sistemi, alçak-yörüngeli (LEO) bir uydunun üzerinden geçerken işaretçinin sabit frekanslı sinyalinde gözlediği Doppler kaymasının zamanla değişiminden (Doppler eğrisi) işaretçinin konumunu çözer. Uydu yaklaşırken frekans yukarı, en yakın noktada (TCA) gerçek frekans, uzaklaşırken aşağı kayar; bu eğrinin tam şekli, işaretçinin uydu izine göre konumunu belirler. Bu, tek bir hareketli alıcının (uydu) Doppler ölçümüyle (FDOA fikrinin saf bir uygulaması) konum bulmasıdır ve gerçekten can kurtarır.

```
 Doppler eğrisi ile konum (LEO uydu geçişi):
 frekans ▲
 kayması │ ●●●
         │     ●●●          (uydu yaklaşırken + ; TCA'da 0 ; uzaklaşırken −)
   0 ─────────────●●●──────────── zaman
         │            ●●●         eğrinin eğimi ve sıfır-geçiş anı
         │               ●●●      → işaretçinin konumunu kodlar
```

Not: Modern kurtarma işaretçileri sıkça GNSS konumunu sinyalin içine gömer (kendi konumunu söyler); bu durumda Doppler çözümü ya yedek ya da doğrulama olur. Yine de Doppler-tabanlı lokasyon, altyapının temel ve tarihsel mekanizmasıdır.

| Modern sistem | Kullandığı çekirdek yöntem | Bu bölümdeki karşılığı |
|---|---|---|
| WiFi/BLE RSSI fingerprint | RSS + parmak izi | Bölüm 6 (RSS) |
| WiFi FTM (RTT) | Gidiş-dönüş zamanı (mesafe) | Zaman tabanlı menzil |
| BLE AoA/AoD | Açı ölçümü | Bölüm 3 (AOA) |
| Hücresel Cell-ID+TA | Hücre + zaman halkası | Kaba menzil |
| Hücresel OTDOA / 5G | TDOA (ve AOA) | Bölüm 4 (TDOA) |
| EPIRB/Argos | Doppler eğrisi | Bölüm 5 (FDOA) |

---

<a id="12"></a>
## 12. Hata Kaynakları: Çok-Yol, Geometri (GDOP), Senkronizasyon, Kalibrasyon

Geolocation'ın matematiği temiz, dünyası ise kirlidir. Gerçek bir çözümün doğruluğu, dört büyük hata kaynağının toplam etkisiyle belirlenir. İyi mühendislik, bu hataları tanımak ve ölçüm geometrisini/donanımını ona göre tasarlamaktır.

### Çok-yol (multipath)

Sinyal kaynaktan alıcıya yalnızca doğrudan yoldan değil, binalardan/araziden yansıyarak da gelir. Bu yansımalar:

- AOA'da: yanlış açılar üretir (yansımanın geldiği yön, kaynağın yönü değildir); DF "hayalet" yönler gösterir.
- TDOA'da: korelasyon tepesini bozar ya da çoğaltır; yansıyan kopya, doğrudan kopyadan biraz gecikmiş ikinci bir tepe yaratır ve gerçek `Δt`'yi karıştırır.
- RSS'te: gücü dramatik biçimde (fading) oynatır; mesafe kestirimini çökertir.

Çok-yol, kentsel/iç-mekân ortamların baş belasıdır ve alt-uzay yöntemlerinde (korele kaynak sorunu) bile özel önlem gerektirir.

### Geometri ve GDOP

Belki en az sezgisel ama en belirleyici faktör geometridir. Aynı ölçüm hatasıyla bile, alıcıların hedefe göre yerleşimi çözümün doğruluğunu uçurum farkıyla değiştirir. Bu etkiyi niceleyen kavram GDOP (Geometric Dilution of Precision — Geometrik Hassasiyet Seyrelmesi):

```
 konum hatası ≈ GDOP × ölçüm hatası
```

GDOP, "ölçüm hatasının konum hatasına ne oranda büyütüldüğü"dür; birimsiz bir çarpandır. İyi geometride GDOP küçüktür (≈1-2); kötü geometride büyür (10+), yani aynı ölçüm hatası konumda kat kat büyük hataya dönüşür.

```
 GDOP geometriye nasıl bağlı (TDOA/triangulation sezgisi):

 İYİ geometri (alıcılar hedefi geniş açıyla sarıyor):
        A ●         ● B          kısıt eğrileri DİK kesişir
            ╲      ╱              → küçük, dairesel hata elipsi
             ╲ ●  ╱               → GDOP küçük
              ╲│ ╱
        C ●────┼────              (hedef alıcıların ortasında)

 KÖTÜ geometri (alıcılar hep aynı tarafta, hizada):
   A ● B ● C ●                    kısıt eğrileri SIĞ açıyla kesişir
        ╲ ╲ ╲                     → uzun, ince hata elipsi
         ╲ ╲ ╲                    → GDOP büyük (mesafe yönünde patlar)
            ● hedef (uzakta, aynı yönde)
```

Temel kural: alıcılar hedefi ne kadar farklı yönlerden "sararsa" geometri o kadar iyi, GDOP o kadar küçük olur. Tüm alıcılar hedefe göre hemen hemen aynı yöndeyse (hizalıysa), kısıt eğrileri sığ açıyla kesişir, hata elipsi mesafe yönünde uzar ve çözüm o yönde neredeyse belirsizleşir. Bu, GPS'te de aynıdır (uyduların gökyüzüne iyi dağılması = düşük GDOP); pasif geolocation'da alıcı yerleşimini tasarlarken birinci düşünülmesi gereken budur. Sıkça karşılaşılan kötü durum collinear (eş-doğrusal) yerleşimdir: alıcılar bir doğru üzerinde dizilirse o doğruya dik yöndeki belirsizlik patlar.

### Senkronizasyon hatası

TDOA ve FDOA için ele alındı (Bölüm 4-5): zaman senkronizasyonundaki `Δt` hatası doğrudan `c·Δt` konum hatasına (1 ns ≈ 0,3 m), frekans referansındaki kararsızlık FDOA'da Doppler farkını yutmaya dönüşür. GPSDO sınıfı ortak referans bu hataları minimuma indirir; onsuz TDOA/FDOA pratikte çalışmaz.

### Kalibrasyon hatası

AOA/DF için ele alındı (Bölüm 8): çok-kanallı alıcıda kanallar arası bilinmeyen faz/kazanç farkları açıyı sistematik olarak kaydırır. Kalibrasyon bu donanım-kaynaklı farkı ölçüp çıkarır; kalibre edilmemiş bir dizi, matematik mükemmel olsa bile yanlış açı verir.

| Hata kaynağı | Hangi yöntemi vurur | Azaltma yolu |
|---|---|---|
| Çok-yol (multipath) | Hepsi (AOA, TDOA, RSS) | Yüksek bant genişliği, alt-uzay yöntemleri, anten yerleşimi, mekânsal yumuşatma |
| Geometri / GDOP | Hepsi | Alıcıları hedefi saracak biçimde dağıtmak; collinear'dan kaçınmak |
| Senkronizasyon | TDOA, FDOA | GPSDO / ortak zaman+frekans referansı |
| Kalibrasyon | AOA / faz-tabanlı DF | Referans sinyaliyle kanal faz/kazanç kalibrasyonu |
| SNR / bant genişliği | TDOA (korelasyon keskinliği), hepsi | İntegrasyon süresi, anten kazancı, LNA (Bölüm 3) |

> Mühendislik sezgisi: Geolocation hata bütçesinde matematik genelde en kolay kısımdır; gerçek hatayı çok-yol ve geometri belirler. "Çözümüm neden kötü?" sorusunun cevabı çoğu zaman algoritmada değil, ya yansımalarda ya da alıcıların kötü (hizalı) yerleşimindedir. Önce geometriyi düzelt (GDOP), sonra senkronizasyon/kalibrasyonu sıkılaştır, sonra algoritmayı zenginleştir.

---

<a id="13"></a>
## 13. Alıştırmalar (Yasal, Kendi Cihazların ve Kâğıt Üzerinde)

> Bu alıştırmalar yalnızca kendi cihazların, açık/yasal sinyaller ve gözlem/hesap içindir. İletim, karıştırma ve yetkisiz içerik çözme yoktur. Pasif yön bulma yaparken yalnızca yasal olarak dinlenebilir sinyalleri (örneğin bir FM yayın istasyonu) ya da KENDİ lisanslı/serbest-bant test vericini hedef al; başkasının haberleşmesini konumlandırma. Kâğıt-kalem alıştırmaları hiçbir yayın gerektirmez; hepsi gözlem ve hesaptır. Şüphedeysen yapma.

### A) Bir FM istasyonunun yönünü bulmak (DOA refleksi)

KrakenSDR'in varsa, faz-tutarlı diziyi (tercihen UCA, Bölüm 8) kur, kalibrasyonu yap ve güçlü bir yerel FM istasyonunun (yasal, sürekli yayın) geliş açısını ölç. DOA çıktısını, istasyonun haritadaki gerçek yönüyle karşılaştır. KrakenSDR yoksa, yönlü bir anten (Yagi, Bölüm 3) ile sinyalin en güçlü olduğu yönü el ile tarayarak kaba bir DF yap (foxhunt mantığı: maksimum güç = kaynak yönü). Şu soruları yanıtla: ölçtüğün açı gerçek yöne ne kadar yakın? Kalibrasyonu atlarsan/bozarsan açı nasıl kayıyor? Anteni döndürünce ön-arka belirsizliği görüyor musun (ULA ise)?

| Gözlem | Not |
|---|---|
| Ölçülen DOA | ? derece |
| Gerçek yön (haritadan) | ? derece |
| Fark (hata) | ? |
| Kalibrasyon etkisi | ? |

### B) Kendi düşük-güç ISM/lisanslı test vericinin yönü (kontrollü DF)

Yasal sınırlar içinde KENDİ bir test kaynağın varsa (örneğin serbest ISM bandında çalışan kendi cihazın ya da lisanslıysan kendi amatör verici test sinyalin), onu bilinen bir konuma koy ve A'daki yöntemle yönünü bul. Avantaj: gerçek yönü tam bildiğin için hatayı kesin ölçebilirsin. Kaynağı farklı uzaklıklara koyarak `σ_yanal ≈ R · σ_θ` ilişkisini gözlemle: uzaklaştıkça aynı açı hatası daha büyük konum hatası mı veriyor? (Yalnızca yasal güç/bant sınırında yayın; lisanssız TX yok.)

### C) TDOA geometrisini kâğıtta hesaplamak (hiperbol sezgisi)

Hiç yayın yapmadan, tamamen kâğıt-kalem: Bir düzlemde iki alıcı koy, örneğin `A = (0, 0)` ve `B = (10, 0)` km. Bir kaynak `S = (3, 4)` km olsun.

1. `|S−A|` ve `|S−B|` mesafelerini hesapla (Pisagor): `|S−A| = √(3²+4²) = 5` km; `|S−B| = √((3−10)²+4²) = √(49+16) = √65 ≈ 8,06` km.
2. Mesafe farkı: `Δd = |S−A| − |S−B| = 5 − 8,06 = −3,06` km.
3. Zaman farkı: `Δt = Δd / c = −3,06 km / (3×10⁵ km/s) ≈ −10,2 µs`.
4. Bu `Δd = −3,06` km değeri, A ve B odaklı bir hiperbol kolunu tanımlar; kaynak bu kolun üstündedir. Üçüncü bir alıcı `C` ekleyip A-C için aynı hesabı yaparsan, ikinci hiperbol kaynakta birinciyle kesişir.

Şunu doğrula: `Δt`'yi `±100 ns` değiştirirsen (senkronizasyon hatasını taklit), hiperbol ne kadar kayar, dolayısıyla konum kestirimi kabaca ne kadar bozulur? (İpucu: `c·100ns ≈ 30 m`.)

### D) GDOP'un istasyon yerleşimine bağlılığını göstermek (geometri sezgisi)

Kâğıt üzerinde iki senaryo çiz:

- Senaryo 1 (iyi): Üç alıcıyı bir kaynağı çevreleyecek biçimde, kabaca eşkenar üçgenin köşelerine koy; kaynak ortada.
- Senaryo 2 (kötü): Aynı üç alıcıyı neredeyse bir doğru üzerinde (collinear), kaynak da o doğrunun uzağında ve aynı yönde olacak şekilde koy.

Her senaryoda her alıcıdan kaynağa giden kısıt eğrilerinin (TDOA için hiperbol normalleri, ya da AOA için açı doğruları) hangi açıyla kesiştiğini kabaca çiz. Senaryo 1'de eğriler dik kesişir (küçük, dairesel hata elipsi, düşük GDOP); Senaryo 2'de sığ açıyla kesişir (uzun, ince elips, yüksek GDOP). Aynı `±σ` ölçüm hatasını her iki çizime ekleyip hata elipsinin nasıl uzadığını gözlemle. Sonuç: ölçüm cihazın aynı olsa bile, alıcıları nereye koyduğun doğruluğu belirler.

### E) Doppler eğrisinden konum çıkarımı (FDOA/Argos sezgisi)

Kâğıt-kalem düşünce deneyi: Sabit frekanslı bir kaynağın üzerinden geçen bir LEO uyduyu düşün. Uydu yaklaşırken gözlenen frekans `f_0`'ın üstünde, en yakın geçiş anında (TCA) tam `f_0`, uzaklaşırken altındadır. `Δf = f_0 · (v_r/c)` formülünü kullanarak şunu akıl yürüt: Doppler eğrisinin sıfır-geçiş anı (TCA) uydu izinin kaynağa en yakın olduğu anı, eğrinin eğimi (frekansın ne kadar hızlı değiştiği) ise kaynağın uydu izine olan dik mesafesini kodlar. Neden eğri daha "dik" ise kaynak uydu izine daha yakındır? (İpucu: yakın geçişte radyal hız bileşeni daha hızlı işaret değiştirir.) Bu, EPIRB/Argos'un (Bölüm 11) tek bir hareketli alıcıyla nasıl konum bulduğunu somutlaştırır.

---

<a id="14"></a>
## 14. Hızlı Referans ve Diğer Bölümler

### Kavram kartı

| Kavram | Bir cümlelik öz |
|---|---|
| Pasif geolocation | Yaymadan, sinyalin geometrisinden kaynağın konumunu kestirme |
| AOA / DOA | Geliş açısı; tek istasyon doğru, iki istasyon (triangulation) konum verir |
| Faz-açı bağı | `θ = arcsin(λ·Δφ / 2π·d)`; `d>λ/2` belirsizlik üretir |
| Hata elipsi | Açı/ölçüm belirsizliğinin kesişimi; geometri biçimini belirler |
| TDOA | Varış zaman farkı → hiperbol; 3 alıcı = multilateration |
| TDOA senkronizasyon | 1 ns ≈ 0,3 m; GPSDO/ortak saat şart; çapraz-korelasyonla `Δt` |
| FDOA | Doppler farkı → eş-Doppler eğrisi; hareketli platformda TDOA'yı tamamlar |
| RSS | Güç→mesafe; çok-yol/`n` belirsizliği nedeniyle en zayıf yöntem |
| DF teknikleri | Watson-Watt → pseudo-Doppler → faz interferometri → MUSIC/ESPRIT |
| MUSIC | Gürültü alt-uzayına diklikten süper-çözünürlüklü açı tepeleri |
| KrakenSDR | 5 kanal faz-tutarlı DOA; kalibrasyon ve dizilim (UCA/ULA) kritik |
| Pasif radar (PCL) | Fırsat aydınlatıcısı (FM/DVB-T) yansımasını CAF ile çözme |
| CAF | Gecikme-Doppler düzleminde tepe = hedefin bistatik menzil+hızı |
| Kalman filtresi | Öngör+düzelt ile gürültülü ölçümlerden pürüzsüz iz |
| GDOP | konum hatası ≈ GDOP × ölçüm hatası; geometri her şeydir |

### Ezber sezgiler

- Her ölçüm bir kısıt eğrisi (AOA doğru, TDOA hiperbol, RSS daire) koyar; konum kesişimdir.
- Pasifte mutlak değil FARK ölçülür: TDOA zaman farkı, FDOA Doppler farkı (kaynağın bilinmeyeni sadeleşir).
- AOA hatası mesafeyle büyür (`σ_yanal ≈ R·σ_θ`): yakında güçlü, uzakta zayıf.
- TDOA nanosaniye senkronizasyon ister (1 ns ≈ 0,3 m); geniş bant keskin korelasyon = iyi `Δt`.
- TDOA+FDOA birlikte güçlüdür çünkü farklı şekilli eğriler keskin kesişir.
- RSS bir mesafe ölçer değil, kaba yakınlık göstergesidir; çok-yol onu çökertir.
- MUSIC açıyı güç tepesi değil "diklik keskinliği" olarak arar → süper-çözünürlük.
- Pasif radar başkasının enerjisiyle görür; CAF = TDOA(gecikme)+FDOA(Doppler) tek yüzeyde.
- Kalman geolocation noktasını yörüngeye çevirir; öngör+düzelt, belirsizliği ikisinden de küçültür.
- Doğruluğu matematik değil çoğu zaman geometri (GDOP) ve çok-yol belirler; önce alıcı yerleşimini düzelt.

### Ve daima: yasal sınır ve perspektif

Bu bölümdeki tüm teknikler tasarım gereği pasif yöntemlerdir; hiçbiri iletim, karıştırma ya da yetkisiz içerik çözme önermez. Pasif yön bulmayı yalnızca kendi cihazların, açık/yasal sinyaller (FM yayını gibi) ya da lisanslı kendi test vericin üzerinde uygula; başkasının haberleşmesini konumlandırmak/izlemek/takip etmek yasal sınırlara tabidir ve bu bölüm açıkça savunma ve kendi-test perspektifindedir. Hedef operatörlük değil, mühendislik sezgisidir: bir DF ekranındaki açı çizgisini, bir multilateration hata elipsini ya da bir Kalman izini gördüğünde arkasındaki geometriyi ve matematiği tanımak. Bandını, ülkeni ve sürümünü teyit et (Türkiye'de BTK); bu kitap anlama ve savunma içindir.

---

### Serinin diğer bölümleri (çapraz referans)

- SIGINT_01 — RF Fiziği, Spektrum ve Modülasyon: `f = c/λ`, dB/dBm, FSPL, Doppler ve faz kavramlarının fiziksel temeli. (Bu bölümdeki `Δφ`, Doppler ve güç-mesafe formüllerinin kökeni oradadır.)
- SIGINT_02 — SDR Donanımları ve Yazılım Ekosistemi: RTL-SDR/HackRF/USRP, faz-tutarlı çoklu-alıcı altyapısı. (KrakenSDR ve TDOA için gereken çok-kanallı/senkron donanımın bağlamı orada.)
- SIGINT_03 — Antenler, RF Donanımı ve Devre Tasarımı (DF): dipol/Yagi/dizi, yön bulmanın anten temelleri, rezonans ve `λ/2` aralığı. (AOA ölçümünün, yani bu bölümün anten/faz altyapısı orada.)
- SIGINT_04 — Sayısal Demodülasyon ve Protokol Kod Çözme: konumlandırılan sinyalin içeriğine inen kod çözme adımı. (Geolocation "nerede"yi, demodülasyon "ne"yi verir.)
- SIGINT_05 — Hücresel, WiFi/BT ve IoT Spektrumu (Savunma Bakışı): GSM/LTE/5G, WiFi/BLE imzaları, cihaz sızıntı yüzeyi. (Bu bölümdeki modern konumlandırmanın — OTDOA, RTT, BLE AoA — protokol bağlamı orada.)
- SIGINT_06 — TEMPEST, RF Sızıntısı ve SIGINT'e Karşı Savunma: emisyon güvenliği, OPSEC. (Pasif geolocation'a karşı savunma — yaymamak, izlenebilir imzayı azaltmak — oranın felsefesidir.)
- SIGINT_07 — Disiplinler ve Sinyal Ayıklama: COMINT/ELINT/FISINT, PDW, deinterleaving, AMC, trafik analizi. (Bu bölüm, Bölüm 7'de ayıklanan yayıcıyı haritaya oturtur; AOA orada "en kararlı ayraç", burada bir konum kısıtıdır.)
- SIGINT_08 — Frekans Tahsisi ve Bant Planı: hangi frekans kime ait. (Konumlandıracağın sinyalin hangi servise ait olduğunu, dolayısıyla yasal statüsünü oradan okursun.)

> Kapanış: Bir verici konumunu asla doğrudan söylemez; ama yaydığı dalga, geldiği açıda, alıcılara varış zamanlarının farkında, hareketle değişen Doppler'inde ve zayıflayan gücünde konumunu istemeden ele verir. Yer tespiti bu geometrik ipuçlarını kısıt eğrilerine çevirmek, yön bulma açıyı ölçmek, multilateration eğrileri kesiştirmek, takip ise noktayı zaman içinde bir yörüngeye bağlamaktır. Bu zinciri kavradığında bir geolocation ekranındaki hata elipsi artık soyut bir leke değil, fiziğini ve geometrisini tanıdığın bir belirsizlik bütçesi olur. Ve her zaman olduğu gibi: bu sezgiyi yalnızca kendi cihazlarının ve yasal sinyallerin üzerinde, pasif ve savunma amaçlı sına.
>
> Bu doküman Kanije Kalesi güvenlik/teknik rehberleri koleksiyonunun SIGINT serisinin 9. bölümüdür. İlgili: SIGINT_01–08, `VERACRYPT_USTALIK_REHBERI.md`, `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`.
