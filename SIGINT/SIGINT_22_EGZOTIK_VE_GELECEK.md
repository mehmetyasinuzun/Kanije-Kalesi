# SIGINT EL KİTABI — BÖLÜM 22: EGZOTİK YAYILIM, SIRA DIŞI SİNYALLER VE GELECEĞİN SIGINT'İ

## Atmosferin, Ayın ve Meteorların Radyo Aynaları; Mega-Konstelasyonlar, Bilişsel Radyo ve Kuantum Çağına Doğru Spektrum

> Amaç: Önceki bölümler spektrumun gündelik halini anlattı: fiziği (Bölüm 1), donanımı (Bölüm 2), anteni (Bölüm 3), demodülasyonu (Bölüm 4-5), savunmayı (Bölüm 6), ayıklamayı (Bölüm 7) ve frekans planını (Bölüm 8). Bu bölüm, o sıradan çerçevenin kenarlarına gider: sinyalin atmosferin üst katmanlarından, bir meteorun iyonize izinden, hatta Ay'ın yüzeyinden sektiği egzotik yayılım modlarına; denizaltıyı bulan devasa ELF antenlerinden yıldırımın doğal radyosuna; gürültünün metrelerce altında konuşan zayıf-sinyal sanatına; ve oradan da geleceğe — gökyüzünü saran mega-konstelasyonlara, kendi spektrumunu pazarlayan bilişsel radyolara, OFDM sonrası dalga formlarına, post-kuantum kriptografiye ve yapay zekâ-tanımlı radyoya. Kullanıcı bilinçli olarak en bilinmeyen, en egzotik ve en güncel/gelecek konuları istedi; bu bölüm o isteğin karşılığıdır.

> Yasal çerçeve: Seri boyunca olduğu gibi bu bölüm de anlama, savunma ve spektrum okuryazarlığı amaçlıdır. Anlatılan her şey ya pasif gözlemdir (dinleme, izleme, hesap) ya da kavramsal-fiziksel açıklamadır; hiçbir yerde yetkisiz iletim, karıştırma veya yetkisiz içerik çözme önerilmez. Alıştırmalar yalnızca alıcı (RX) tarafıyla, doğal sinyallerle ve açık/yasal yayınlarla sınırlıdır. Bir uydunun spektrumunu izlemek genelde serbestken belirli bandların kaydı/yayılması çoğu ülkede düzenlemeye tabidir; kendi ülkenin ve sürümünün mevzuatını teyit et. Fizik evrenseldir, ama yasa yereldir.

> Fizik dürüstlüğü notu: Bu bölüm çok geniş bir yelpazeyi kapsar ve bir kısmı (kuantum algılama, QKD'nin saha olgunluğu, THz haberleşme, OTFS'in kitlesel konuşlanması) hâlâ araştırma veya erken-konuşlanma aşamasındadır. Bu tür yerlerde "araştırma aşamasında" veya "teyit edilmeli" notunu açıkça düşeceğim. Yerleşik fizik (iyonosfer, EME yol kaybı, Doppler) ile spekülatif/gelişmekte olanı karıştırmamak, bu bölümün en önemli disiplinidir.

---

## İÇİNDEKİLER

1. [Egzotik Yayılımın Haritası: Sıradan Çizgi-Görüş Ötesi](#1)
2. [HF ve İyonosfer: Sky-Wave, NVIS, Gri-Çizgi, Güneş Döngüsü](#2)
3. [Troposfer Modları: Sapma, Ducting ve Sporadic-E](#3)
4. [Meteor Scatter: Bir Meteorun İzinden Saçılma (MSK144)](#4)
5. [Moonbounce/EME: Ay'ı Ayna Yapmak ve Aşırı Zayıf Sinyal](#5)
6. [Aurora ve Diğer Doğal Yansıtıcılar](#6)
7. [Spektrumun Dibi: VLF/ELF/SLF ve Denizaltı Haberleşmesi](#7)
8. [Navigasyonun Tarihi ve Yeniden Doğuşu: Omega, LORAN, eLORAN](#8)
9. [Doğal Radyo: Whistler, Sferic, Schumann Rezonansı](#9)
10. [Zayıf-Sinyal Sanatı: WSJT-X, Kodlama Kazancı ve Shannon](#10)
11. [Mega-Konstelasyonlar: Starlink, OneWeb, Kuiper, Iridium NEXT](#11)
12. [Uydu IoT ve Doğrudan-Cihaza Uzay Bağlantısı](#12)
13. [Bilişsel Radyo ve Dinamik Spektrum: White Space, CBRS](#13)
14. [Yeni Dalga Formları: OTFS, FBMC, Massive MIMO, mmWave/THz](#14)
15. [Görünür Işık Haberleşmesi (Li-Fi) ve Optik/Lazer Bağlar](#15)
16. [Güvenliğin Geleceği: Post-Kuantum, QKD ve Kuantum Algılama](#16)
17. [Yapay Zekâ-Tanımlı Radyo ve Kendini Ayarlayan Spektrum](#17)
18. [Sıra Dışı ve Merak: Pirat Radyo, Gizli Kanallar, Astronomik Kaynaklar](#18)
19. [Bilinmeyen Sinyalle Karşılaşma Metodolojisi (Kapanış Sentezi)](#19)
20. [Alıştırmalar (Yasal, Yalnızca Alıcı ve Gözlem)](#20)
21. [Hızlı Referans ve Diğer Bölümler](#21)

---

<a id="1"></a>
## 1. Egzotik Yayılımın Haritası: Sıradan Çizgi-Görüş Ötesi

Bölüm 1'de bir sinyalin temel yayılımını verdik: serbest uzay yol kaybı (FSPL), çizgi-görüş (line-of-sight, LOS) ve frekansa göre değişen davranış. O çerçeve, gündelik VHF/UHF haberleşmesinin %95'ini açıklar. Ama spektrumun en ilginç olayları tam da bu çerçevenin dışında, "sinyal nasıl buraya geldi?" diye sorduğumuz yerlerde yaşanır. Egzotik yayılım, sinyalin doğal bir ortam tarafından (iyonosfer, troposfer, bir meteorun izi, Ay'ın yüzeyi) yansıtıldığı, kırıldığı veya saçıldığı; böylece çizgi-görüşün çok ötesine — bazen yarım dünyaya, bazen Ay'a kadar — ulaştığı modların ortak adıdır.

Bir RF/yayılım uzmanı için bu modların ortak teması şudur: ortam, sinyalin yol açısını değiştiren bir aynaya, merceğe veya saçıcıya dönüşür. Hangi ortamın hangi frekansta ayna gibi davrandığı, frekans ile ortamın fiziksel ölçeği arasındaki orana bağlıdır. İyonosfer HF'i (3-30 MHz) yansıtır ama VHF'i delip geçer; troposfer VHF/UHF'i bükebilir; bir meteorun iyonize izi VHF'i kısa süreliğine saçar; Ay yüzeyi neredeyse her frekansı zayıfça yansıtır. Bu eşleşmeleri tek bir tabloda toplamak, bölümün geri kalanının iskeletini verir.

| Yayılım modu | Yansıtıcı/ortam | Tipik frekans | Tipik mesafe | Gerekli koşul | Karakter |
|---|---|---|---|---|---|
| Sky-wave (HF) | İyonosfer (F katmanı) | 3-30 MHz | 1000-12000+ km | Uygun katman yoğunluğu, gece/gündüz | Kıtalararası, değişken |
| NVIS | İyonosfer (dik geliş) | 2-10 MHz | 0-400 km | Yüksek açılı yayın | Bölgesel, gölge-doldurucu |
| Tropo ducting | Troposfer (sıcaklık tersinmesi) | 50 MHz-10 GHz | 100-2000 km | İnversiyon katmanı, deniz | Anlık, hava-bağımlı |
| Sporadic-E (Es) | E katmanında yoğun bulutlar | 30-220 MHz | 500-2300 km | Es bulutu (yaz zirvesi) | Ani, kısa, çok güçlü |
| Meteor scatter | Meteor iyonize izi | 30-150 MHz | 600-2300 km | Meteor girişi (saniyeler) | Burst (kısa patlama) |
| Moonbounce (EME) | Ay yüzeyi | 50 MHz-10 GHz | ~770000 km (gidiş-dönüş) | Yüksek güç + büyük anten | Aşırı zayıf, gecikmeli |
| Aurora | Auroral iyonizasyon | 30-450 MHz | 500-2000 km | Jeomanyetik fırtına | Çarpık/gürültülü ses |
| Ground-wave | Yer yüzeyi (yüzey dalgası) | < 3 MHz (en iyi) | 10-1000 km | Düşük frekans, iletken yer | Kararlı, yavaş zayıflar |
| VLF/ELF | İyonosfer-yer dalga kılavuzu | 3 Hz-30 kHz | Küresel | Devasa anten/güç | Su altına işler, çok yavaş |

Bu tablo bir referans çapasıdır; aşağıdaki başlıklar her satırı fiziğiyle açar. Dikkat edilecek temel sezgi: egzotik modlar "her zaman açık" değildir. İyonosfer güneşin keyfine, ducting havaya, meteor scatter gökyüzüne bir cismin girmesine, EME ise Ay'ın ufkun üstünde olmasına bağlıdır. Bu kesintililik, egzotik yayılımı hem büyüleyici hem de operasyonel olarak güvenilmez kılar; bu yüzden zayıf-sinyal modları (Bölüm 10) bu kesintili pencerelerden en çok bilgiyi sıkıştırmak için doğmuştur.

```
 Yayılım modlarının yükseklik haritası (kabaca, ölçeksiz):

  yükseklik
    ▲
 ~400 km │   ░░░░░░░░ F katmanı (HF sky-wave, NVIS) ░░░░░░░░
         │
 ~110 km │   ▒▒▒▒ E katmanı / Sporadic-E ▒▒▒▒   (meteor izleri ~85-120 km)
         │   ──── D katmanı (gündüz, HF'i yutar) ────
  ~12 km │   ≈≈≈≈ Troposfer: ducting, tropo-sapma ≈≈≈≈
         │
    0 km │ ███ Yer: ground-wave, yüzey dalgası ███
         └──────────────────────────────────────────────► (yatay mesafe)

         Ay'a giden EME yolu bu şemanın çok ötesinde, ~384000 km uzakta.
```

---

<a id="2"></a>
## 2. HF ve İyonosfer: Sky-Wave, NVIS, Gri-Çizgi, Güneş Döngüsü

Egzotik yayılımın en eski ve en zengin örneği HF sky-wave'idir. Kısa dalga radyonun yarım dünyaya ulaşması, bir mucize değil, iyonosferin sistematik bir davranışıdır. İyonosfer, atmosferin yaklaşık 60-1000 km arasındaki, güneş ışınımının (özellikle morötesi ve X-ışını) gaz moleküllerini iyonize ettiği katmandır. İyonize, yani serbest elektron içeren bu bölge, belirli frekanslardaki radyo dalgalarını kıran (ve yeterli açıda gelende geri yansıtan) bir plazma gibi davranır.

### İyonosferin katmanları

İyonosfer tek bir tabaka değil, yoğunluğu ve davranışı farklı alt katmanlardan oluşur. Geleneksel olarak D, E ve F (gündüz F1/F2'ye ayrılır) diye adlandırılırlar.

| Katman | Yükseklik | Gündüz/gece | Radyoya etkisi |
|---|---|---|---|
| D | ~60-90 km | Yalnız gündüz | Düşük HF'i (özellikle < 10 MHz) yutar; gece kaybolur |
| E | ~90-150 km | Gündüz baskın | Orta HF'i yansıtır; Sporadic-E ayrı bir olgu (Bölüm 3) |
| F1 | ~150-220 km | Yalnız gündüz | Gündüz F2'den ayrılır, ek kırılma |
| F2 | ~220-400+ km | Gece de sürer | HF sky-wave'in ana yansıtıcısı; en yüksek frekansları döndürür |

Sky-wave'in mantığı şudur: yere göre belirli bir açıyla yukarı yayılan HF dalgası, iyonosferde yeterince kırılırsa yere geri döner; orada yer yüzeyinden tekrar yukarı sekebilir (multi-hop). Her "hop" yüzlerce ila birkaç bin kilometre kazandırır; birkaç hop ile sinyal kıtalararası gider.

```
 Sky-wave tek ve çok-hop:

   iyonosfer (F2)
   ════════════════════════════════════════════════════
      ╲          ╱╲          ╱╲          ╱
       ╲        ╱  ╲        ╱  ╲        ╱
        ╲      ╱    ╲      ╱    ╲      ╱
   ──────╲────╱──────╲────╱──────╲────╱──────────────► yer
        TX          hop1         hop2        RX
   (yüksek açı → kısa skip; alçak açı → uzun skip)

   Skip zone: vericiye yakın, ground-wave bittiği ama
   ilk sky-wave hop'unun düşmediği "sağır" halka.
```

### Kritik frekans, MUF ve LUF

İyonosferin yansıtma yeteneği frekansa bağlıdır. Dik yukarı gönderilen bir sinyalin geri dönebileceği en yüksek frekansa kritik frekans (foF2, F2 katmanı için) denir; bu, katmandaki elektron yoğunluğunun doğrudan ölçüsüdür. Eğik geliş için yansıtılabilen en yüksek frekans, kullanılabilir azami frekanstır (MUF, Maximum Usable Frequency) ve geliş açısına göre kritik frekanstan yüksektir (eğik geliş, daha yüksek frekansları yansıtır). Alt sınırda ise, D katmanının yutması nedeniyle iletişimin kesildiği en düşük kullanılabilir frekans (LUF, Lowest Usable Frequency) bulunur. Çalışılabilir pencere LUF ile MUF arasıdır.

> Mühendislik sezgisi: HF planlaması, "frekansı MUF'un hemen altına, LUF'un üstüne yerleştir" sanatıdır. MUF'a yakın frekanslar daha az yutulur ve daha uzun atlar; ama MUF güneşin durumuna göre saatten saate kayar. Bu yüzden HF operatörü gün boyunca bant değiştirir: gündüz yüksek bandlar (örneğin 14-28 MHz) açık, gece düşük bandlar (3-7 MHz) açıktır. Bu ritmin sebebi tamamen iyonosferin gündüz-gece elektron yoğunluğu değişimidir.

### Güneş döngüsü etkisi

İyonosferin yoğunluğu, en büyük ölçekte 11 yıllık güneş lekesi döngüsüne bağlıdır. Güneş aktivitesinin yüksek olduğu yıllarda (solar maximum) iyonizasyon artar, foF2 ve MUF yükselir; en yüksek HF bandları (örneğin 28 MHz) dünya çapında açılır. Düşük aktivitede (solar minimum) yüksek bandlar büyük ölçüde ölür, iletişim düşük frekanslara çekilir. Güneş patlamaları ise ani iyonosferik bozulmalara (SID, Sudden Ionospheric Disturbance) yol açar; bir X-sınıfı patlamadan gelen yoğun X-ışını D katmanını şişirir ve gündüz HF'ini dakikalar-saatler boyunca yutabilir (radio blackout). Jeomanyetik fırtınalar (güneş rüzgârı/CME kaynaklı) F2'yi bozar ve yüksek enlemlerde iletişimi çökertirken aurora yayılımını (Bölüm 6) açar.

> Güncel bağlam: 25. güneş döngüsünün (Cycle 25) zirvesi 2024-2025 dolayında beklenmekteydi; tam zirve tarihi ve şiddeti gözleme dayalı olarak kurumlarca güncellenir ve kesin değer kaynaktan teyit edilmeli. Pratik sonuç: bu yıllarda yüksek HF bandları tarihsel olarak daha sık ve daha uzun açılır; bir alıcı için bu, dünya çapında zayıf-sinyal yakalama (Bölüm 10) için elverişli bir dönemdir.

### Gri-çizgi (gray-line) yayılımı

İyonosferin en zarif olaylarından biri gri-çizgi yayılımıdır. Dünyada gündüzü geceden ayıran alacakaranlık kuşağına (terminator) gri-çizgi denir. Bu hat boyunca D katmanı (yutucu) ya henüz oluşmamış ya da yeni dağılmıştır; ama F katmanı hâlâ iyonizedir. Sonuç, çok düşük yutmayla uzun mesafe iletişim için kısa bir altın penceredir. Gri-çizgi üzerindeki iki nokta, gün doğumu/batımı anlarında normalde imkânsız olan mesafelere düşük güçle ulaşabilir. Gri-çizgi her gün dünya etrafında dolaşan, yaklaşık iki kez (yerel gün doğumu ve batımı) açılan bir yayılım koridorudur.

```
 Gri-çizgi (terminator) yayılımı:

      GÜNDÜZ        │  gri-çizgi  │        GECE
   (D katmanı       │ (D zayıf,   │   (D yok, F sürer)
    yutuyor)        │  F güçlü)   │
                    │             │
     ◄ yutulma      │  ◄ düşük    │   ◄ düşük
       yüksek       │    yutma →  │     yutma →
                    │             │
   İki gri-çizgi noktası arası: kısa, çok verimli pencere.
   Bu hat, gün boyunca batıdan doğuya dünyayı dolaşır.
```

### NVIS: dik geliş, bölgesel kapsama

Sky-wave uzun mesafeye odaklanırken, aynı iyonosferin bambaşka bir kullanımı NVIS'tir (Near Vertical Incidence Skywave — Dik Gelişe Yakın Sky-wave). Burada sinyal kasten neredeyse dik yukarı (yüksek açıyla, ~70-90 derece) gönderilir; iyonosfer onu adeta bir şemsiye gibi geniş bir bölgeye geri yağdırır. NVIS'in amacı uzak değil, yakın-orta menzildir (0-400 km), özellikle ground-wave'in dağlar/engebe yüzünden yetmediği ve normal sky-wave'in "skip zone" yüzünden boşluk bıraktığı bölgesel kapsama. Tipik olarak düşük HF (2-10 MHz, gündüz/gece farklı) ve alçak, yatay antenler kullanılır.

```
 NVIS geometrisi:

         iyonosfer
   ══════════════════════════════════
       ╲│╱     ╲│╱     ╲│╱     ╲│╱
        │       │       │       │       (dik çıkıp geniş
   ─────┴───────┴───────┴───────┴────►   alana geri "yağar")
       TX  ◄──── 0-400 km bölgesel ────►
   Dağ/engebe ground-wave'i kesse bile yukarıdan doldurur.
```

NVIS, afet haberleşmesi ve engebeli arazide bölgesel bağlantı için klasik bir çözümdür; sky-wave ve ground-wave arasındaki boşluğu doldurur. Bir savunma/dayanıklılık perspektifinden NVIS, altyapı (baz istasyonu, uydu) çöktüğünde bölgesel haberleşmeyi ayakta tutabilen, dış altyapıya bağımlı olmayan bir moddur.

---

<a id="3"></a>
## 3. Troposfer Modları: Sapma, Ducting ve Sporadic-E

İyonosfer HF'in dünyasıyken, VHF/UHF (ve mikrodalga) çoğunlukla iyonosferi delip geçer ve normalde çizgi-görüşle sınırlıdır. Ancak alt atmosfer (troposfer) ve E katmanının özel bir hali, bu frekanslara da çizgi-görüş ötesi menzil kazandıran egzotik kapılar açar.

### Troposferik sapma ve ducting

Troposferde sıcaklık, nem ve basınç yükseklikle değişir; bu da havanın kırılma indisini değiştirir. Normalde dalga hafifçe aşağı kırılır (dünyanın eğriliğini bir miktar takip eder; bu yüzden radyo ufku optik ufuktan biraz uzaktır). Ama bir sıcaklık tersinmesi (inversion) olduğunda — örneğin sıcak hava soğuk havanın üstüne oturduğunda, ki bu deniz üzerinde ve durgun yüksek-basınç havasında sık görülür — kırılma indisinde keskin bir gradyan oluşur. Bu gradyan, VHF/UHF/mikrodalga dalgasını bir dalga kılavuzu (duct) içine hapsedebilir; dalga, yer ile inversiyon katmanı arasında düşük kayıpla yüzlerce, bazen 1000+ km ilerler. Buna troposferik ducting denir.

| Tropo modu | Tetikleyici | Frekans | Menzil | Karakter |
|---|---|---|---|---|
| Standart kırılma | Normal atmosfer | VHF+ | Radyo ufku (~ optik + %15) | Daima var, küçük katkı |
| Tropo-sapma (scatter) | Türbülanslı saçılma | VHF/UHF | 100-500 km | Zayıf ama sürekli mümkün |
| Ducting | Sıcaklık tersinmesi | VHF-10 GHz | 100-2000 km | Güçlü, hava-bağımlı, anlık |

Ducting, özellikle deniz üzerinde (deniz ducti) ve büyük yüksek-basınç sistemlerinin altında dramatiktir; VHF/UHF TV, radyo ve hatta hücresel sinyaller normalde imkânsız mesafelerden alınabilir. Bir SIGINT/spektrum perspektifinden ducting iki yönlü bir sürprizdir: hem uzaktaki bir yayıcıyı beklenmedik biçimde alınır kılar, hem de kendi yayınının çok daha uzaktan duyulmasına yol açar (OPSEC açısından dikkat; Bölüm 6).

```
 Troposferik ducting (kanal içine hapsolma):

   inversiyon katmanı (sıcak/soğuk sınır)
   ──────────────────────────────────────────────
   ↘   ↗   ↘   ↗   ↘   ↗   ↘   ↗   ↘   ↗   ↘
     ↘ ↗     ↘ ↗     ↘ ↗     ↘ ↗     ↘ ↗
   ═════════════════════════════════════════════ yer (özellikle deniz)
   TX ◄────── duct içinde düşük kayıpla 100-2000 km ──────► RX
```

### Sporadic-E (Es)

Sporadic-E, adından da anlaşılacağı gibi düzensiz, ani ortaya çıkan bir olgudur. E katmanında (yaklaşık 90-120 km) zaman zaman çok yoğun, ince, yamalı iyonizasyon bulutları oluşur. Bu bulutlar normalde iyonosferi delip geçecek VHF frekanslarını (özellikle 30-150 MHz, bazen 220 MHz'e kadar) güçlü biçimde yansıtır. Sonuç, 6 metre (50 MHz) ve FM yayın bandı dahil VHF'te ani, çok güçlü, 500-2300 km menzilli açılımlardır.

Es'in kesin oluşum mekanizması tam olarak yerleşmemiştir; rüzgâr makaslaması (wind shear) ve metalik iyonların (meteor kaynaklı) yoğunlaşması başlıca açıklamalardır ve bu mekanizmaların ayrıntısı araştırma konusudur, kesin model kaynaktan teyit edilmeli. Gözlemsel olarak çok nettir: kuzey yarımkürede yaz ayları (özellikle Mayıs-Ağustos) zirvedir, açılımlar dakikalar-saatler sürer ve coğrafi olarak yamalıdır. Es açılımının imzası, normalde ölü olan bir VHF bandında aniden uzak istasyonların çok güçlü (sky-wave gibi zayıf değil, neredeyse yerel güçte) belirmesidir.

> Pratik gözlem: Sporadic-E, alıcı sahibinin en kolay yakalayabileceği egzotik açılımdır çünkü VHF FM/TV bandlarında gündelik dinlemeyle bile fark edilir: 88-108 MHz FM bandında yaz öğleden sonrası aniden uzak şehirlerden istasyonlar girip çıkıyorsa, büyük olasılıkla Es açılımına tanık oluyorsundur. Bu, Bölüm 20'deki alıştırmalardan birinin temelidir.

---

<a id="4"></a>
## 4. Meteor Scatter: Bir Meteorun İzinden Saçılma (MSK144)

Egzotik yayılımın en şiirsel hali, sinyalin bir meteorun ardında bıraktığı iyonize izden sektiği meteor scatter'dır (meteor burst communication, MBC). Atmosfere giren her meteor (çoğu kum tanesi büyüklüğünde) sürtünmeyle yanar ve arkasında, birkaç saniyeliğine var olan, çok ince ve yoğun bir iyonize gaz sütunu (meteor trail) bırakır; bu sütun ~85-120 km yükseklikte oluşur. Bu kısa ömürlü iz, VHF dalgalarını (tipik olarak 30-150 MHz, özellikle 50 ve 144 MHz) yansıtabilen geçici bir aynadır.

Dünya atmosferine her gün milyarlarca meteor girer; bu yüzden gökyüzü aslında sürekli, ama anlık ve rastgele meteor izleriyle pırıldar. Bir meteor scatter bağlantısı, bu rastgele izlerden birinin tam olarak verici ve alıcıyı geometrik olarak birbirine yansıtacak konumda oluşmasını bekler. İz birkaç saniye (kısa izler), nadiren on saniyeler (uzun, yoğun izler) yaşar; bu yüzden meteor scatter haberleşmesi sürekli değil, kısa patlamalar (burst) halindedir: sessizlik, sonra bir-iki saniyelik güçlü açılım, sonra yine sessizlik.

```
 Meteor scatter geometrisi:

         ☄ meteor (atmosfere giriyor, ~100 km)
          ╲   iyonize iz (birkaç saniye yaşar)
   ════════╲═══════════════
       ╱    ╲    ╲          (iz, TX ve RX'i geometrik
      ╱      ╲    ╲          olarak birbirine yansıtacak
   ──╱────────╲────╲──────►  konumdaysa kısa bir "burst" olur)
    TX         600-2300 km   RX
```

### Neden özel bir dalga formu? WSJT ve MSK144

Bir meteor izi yalnızca bir-iki saniye dayandığından, klasik bir telsiz görüşmesi imkânsızdır; bütün bilgiyi o kısa patlamaya sıkıştırmak gerekir. İşte bu, zayıf-sinyal sanatının (Bölüm 10) doğduğu yerlerden biridir. WSJT-X yazılım ailesinin MSK144 modu tam bu iş için tasarlanmıştır: çok kısa süreli (saniyenin altında tekrarlanan), gürültüye dayanıklı, hızlı kodlanmış mesaj çerçeveleri gönderir; alıcı yazılım, bir burst boyunca bu çerçevelerden birini yakalayıp çözmeye çalışır. MSK144, ismindeki MSK (Minimum Shift Keying) ile verimli, sabit-zarflı bir modülasyon kullanır ve güçlü ileri hata düzeltme (FEC) ile bir saniyelik açılımdan tam bir mesaj (örneğin iki çağrı işareti ve bir rapor) çıkarabilir.

> Mühendislik sezgisi: Meteor scatter, "kanal nadiren ve kısa açılıyor" probleminin ders kitabı çözümüdür. Çözüm, mesajı küçültmek (yalnızca temel bilgi), hızlı tekrarlamak (bir burst içinde birden çok deneme) ve ağır kodlamak (her tekrar gürültüde bile çözülebilsin) üçlüsüdür. Bu üçlü, Bölüm 10'daki diğer WSJT modlarının da ortak felsefesidir; fark, her modun farklı bir kanal zorluğuna (kısa açılım, çok düşük SNR, Doppler) ayarlanmış olmasıdır.

Meteor scatter'ın tarihsel bir uygulaması, bazı askeri/uzak-bölge ağlarının (örneğin yüksek enlemlerde, iyonosferin güvenilmez olduğu yerlerde) düşük-veri-hızlı ama dayanıklı bir kanal olarak meteor burst sistemlerini kullanmasıydı; bu sistemler altyapısız, uzun mesafe ve nispeten gizli (kısa, yönlü, anlık) iletişim sağlardı. Meteor yağmurları (örneğin Perseid'ler Ağustos'ta, Geminid'ler Aralık'ta) sırasında iz oranı belirgin artar ve açılımlar sıklaşır; bu, gözlem için en verimli pencerelerdir.

---

<a id="5"></a>
## 5. Moonbounce/EME: Ay'ı Ayna Yapmak ve Aşırı Zayıf Sinyal

Egzotik yayılımın zirvesi, en uzak ve en zorlu yansıtıcıyı kullanır: Ay. Earth-Moon-Earth (EME), halk dilindeki adıyla moonbounce, bir sinyali Ay'a gönderip yüzeyinden yansıyan kısmını yeryüzünde geri almaktır. Ay, yaklaşık 384000 km uzakta, pürüzlü ve radyoyu zayıf yansıtan bir yüzeydir; bu yüzden EME, amatör radyonun teknik olarak en zorlu, sinyal-bütçesi en acımasız modudur.

### Yol kaybı: neden bu kadar zor?

EME'nin zorluğu tamamen mesafede ve Ay'ın zayıf yansıtıcılığındadır. Gidiş-dönüş toplam yol ~768000 km'dir ve serbest uzay yol kaybı mesafenin karesiyle (her yönde) büyür. Buna Ay yüzeyinin düşük yansıtma katsayısı (yaklaşık %7, yani enerjinin çoğu emilir/saçılır) eklenince, geri dönen sinyal astronomik ölçüde zayıftır. Bölüm 1'deki FSPL formülünü hatırlayalım:

```
 FSPL(dB) = 20·log10(d) + 20·log10(f) + 32.44     (d: km, f: MHz)
```

Tek yön Dünya-Ay için d ≈ 384000 km; bu tek başına çok büyük bir kayıptır ve gidiş-dönüş bunu ikiye katlar (yansıma kaybı da ayrıca eklenir). Toplam yol kaybı, frekansa bağlı olarak tipik VHF/UHF EME'de yaklaşık 240-260 dB mertebesindedir (kesin değer frekans ve Ay mesafesinin anlık değerine göre değişir; sayısal değer kaynaktan teyit edilmeli). Bu, dünyadaki herhangi bir karasal yol kaybının çok ötesindedir.

```
 EME yol bütçesi (kavramsal):

   TX ──(yüksek güç, büyük yönlü anten)──► [serbest uzay ~384000 km] ──►
                                                                    ☾ Ay
   ☾ Ay yüzeyi (zayıf, ~%7 yansıma, dağınık) ──► [geri ~384000 km] ──►
   RX (büyük anten + çok düşük gürültülü ön-yükselteç + zayıf-sinyal modu)

   Toplam ~240-260 dB kayıp → sinyal gürültü tabanının altında →
   ancak kodlama kazancı (WSJT) ile çözülür.
```

### EME'yi mümkün kılan üç sütun

EME, üç şeyin bir araya gelmesiyle mümkün olur. Birincisi büyük anten kazancı: yüksek kazançlı, yönlü diziler (VHF'te uzun Yagi dizileri, UHF/mikrodalgada büyük çanak antenler) hem iletimde gücü Ay'a odaklar hem de alımda zayıf yansımayı toplar. İkincisi düşük-gürültü alıcı: alıcının kendi gürültüsü (gürültü tabanı), gelen sinyalden daha zayıf olmalıdır; bu yüzden çok düşük gürültü figürlü (low-noise figure) ön-yükselteçler şarttır. Üçüncüsü ise — ve bu, geleceğe köprüdür — kodlama kazancı: WSJT-X'in JT65 ve özellikle Q65 gibi modları, sinyali gürültü tabanının çok altında (negatif SNR'da) bile çözebilen ağır hata düzeltme kullanır. Bu üçüncü sütun, donanım gereksinimini dramatik biçimde düşürmüştür; eskiden devasa istasyonlar gerektiren EME, bugün mütevazı (ama yine de ciddi) kurulumlarla yapılabilir.

### EME'nin kendine özgü bozulmaları

EME sinyali yalnızca zayıf değil, aynı zamanda kendine özgü biçimlerde bozulur:

- Yol gecikmesi: Işık hızında gidiş-dönüş yaklaşık 2.5 saniyedir. Kendi yayınının yankısını 2.5 saniye sonra duymak (kendi yankını dinlemek) EME'nin meşhur deneyidir.
- Doppler kayması: Ay ve Dünya'nın göreli hareketi (Ay'ın yörüngesi + Dünya'nın dönüşü) taşıyıcı frekansını kaydırır; bu kayma gün boyunca değişir ve hassas modlarda telafi gerekir.
- Libration fading: Ay'ın görünür yüzeyinin hafif salınımı (libration) ve pürüzlü yüzeyin çok-yollu yansıması, sinyalde dakikalar içinde derinlemesine değişen, "ışıldayan" bir sönümleme (fading) yaratır.
- Faraday dönmesi: Sinyal iyonosferden iki kez geçerken polarizasyon düzlemi döner; bu, sabit-polarizasyonlu antenlerde periyodik sinyal kayıplarına yol açar.

> Mühendislik sezgisi: EME, "kanal mümkün olan en kötüsü" sınır durumudur ve bu yüzden zayıf-sinyal kodlamasının (Bölüm 10) en saf vitrinidir. Burada öğrenilen ders evrenseldir: yeterince güçlü ileri hata düzeltme ve uyumlu modülasyonla, sinyali gürültünün altına gömüp yine de çıkarmak mümkündür. Bu ilke, EME'den derin-uzay haberleşmesine, oradan da modern hücresel ağların gürültü-sınırlı kıyılarına kadar uzanır.

---

<a id="6"></a>
## 6. Aurora ve Diğer Doğal Yansıtıcılar

Egzotik yayılım modlarının haritasını tamamlamak için birkaç daha az yaygın ama öğretici yansıtıcıyı toparlayalım.

### Auroral yayılım

Güçlü jeomanyetik fırtınalar sırasında, kutup bölgelerindeki yüksek atmosfer (auroral oval) yoğun biçimde iyonize olur; bu iyonizasyon, VHF (özellikle 50 ve 144 MHz) ve hatta UHF'in alt ucundaki sinyalleri saçabilir. Auroral yayılımın imzası çok özeldir: yansıma, hareketli ve düzensiz iyonizasyondan geldiği için sinyal güçlü bir Doppler yayılması ve "boğuk", hışırtılı bir distorsiyon kazanır; sesli yayınlar anlaşılmaz ölçüde bozulur, bu yüzden auroral çalışmada da dar-bant, gürültüye dayanıklı modlar tercih edilir. Aurora, antenleri coğrafi hedefe değil, kuzeydeki (kuzey yarımkürede) auroral perdeye doğrultmayı gerektirir — yansıtıcı gökyüzünde, hedefin yönünde değildir.

### Field-aligned irregularities ve diğer modlar

Daha egzotik ve daha çok araştırma konusu olan modlar da vardır: manyetik alan çizgilerine hizalı düzensizliklerden (field-aligned irregularities, FAI) saçılma, ekvatoral yayılı-F (equatorial spread-F) olayları ve trans-ekvatoral yayılım (TEP, manyetik ekvatorun iki yanındaki istasyonlar arası VHF açılımı). Bunların fizik ayrıntıları ve öngörülebilirliği iyonosferik fizik araştırmasının aktif alanlarıdır; kesin mekanizmaları ve koşulları kaynaktan teyit edilmeli. Ortak tema yine aynıdır: güneş-jeomanyetik aktivite, üst atmosferde geçici yansıtıcı/saçıcı yapılar yaratır ve bunlar VHF'e çizgi-görüş ötesi, ama düzensiz ve kısa-ömürlü kapılar açar.

| Doğal yansıtıcı | Kaynak | Frekans | Karakteristik imza |
|---|---|---|---|
| Aurora | Jeomanyetik fırtına iyonizasyonu | 50-450 MHz | Doppler-yayılmış, boğuk, kuzeye yönelik |
| TEP | Ekvatoral iyonosfer | VHF (50-220 MHz) | Manyetik ekvator ötesi uzun açılım |
| Spread-F / FAI | İyonosferik düzensizlikler | HF-VHF | Saçılma, sinyalde yayılma/titreşim |

---

<a id="7"></a>
## 7. Spektrumun Dibi: VLF/ELF/SLF ve Denizaltı Haberleşmesi

Şimdiye kadar yukarı (iyonosfer, Ay) baktık; şimdi spektrumun en dibine, frekansın neredeyse durduğu yere inelim. Çok düşük frekans (VLF, 3-30 kHz), aşırı düşük frekans (ELF, 3-30 Hz) ve süper düşük frekans (SLF, 30-300 Hz) bantları, dalga boylarının onlarca ila on binlerce kilometre olduğu, fiziğin tamamen değiştiği bir dünyadır.

| Bant | Açılım | Frekans | Dalga boyu | Tipik kullanım |
|---|---|---|---|---|
| SLF | Super Low Frequency | 30-300 Hz | 1000-10000 km | Denizaltı (bazı sistemler) |
| ELF | Extremely Low Frequency | 3-30 Hz | 10000-100000 km | Derin denizaltı haberleşmesi |
| VLF | Very Low Frequency | 3-30 kHz | 10-100 km | Denizaltı, navigasyon, saat sinyali |
| LF | Low Frequency | 30-300 kHz | 1-10 km | Uzun dalga yayın, navigasyon (eLORAN) |

### Neden ELF? Deniz suyu ve cilt derinliği

Bu bandların başlıca itici gücü tek bir fiziksel gerçektir: deniz suyu iletkendir ve iletken ortam yüksek frekansları hızla yutar. Bir elektromanyetik dalganın bir iletkene ne kadar nüfuz edebildiği cilt derinliği (skin depth) ile ölçülür; cilt derinliği frekans arttıkça azalır. Deniz suyunda VHF/UHF birkaç santimetrede yok olur, ama ELF (örneğin onlarca Hz) deniz suyuna onlarca metre nüfuz edebilir. Bu yüzden derin dalmış bir denizaltıyla yüzeye çıkmadan haberleşmenin neredeyse tek yolu ELF/SLF'tir; frekans ne kadar düşükse, su o kadar derine işlenir.

```
 Cilt derinliği ve denizaltı (kavramsal):

   deniz yüzeyi ═══════════════════════════════
        │ VHF/UHF: birkaç cm'de yutulur (X)
        │ VLF: birkaç metre
        │ SLF: onlarca metre
        ▼ ELF: en derine işler → derin denizaltı
   ────────────────────────────────────────────
   (frekans düştükçe cilt derinliği artar → daha derine ulaşır)
```

### Bedeli: devasa antenler ve çok düşük veri hızı

ELF'in faturası ağırdır. Birincisi anten: verimli bir anten dalga boyunun anlamlı bir kesri kadar büyük olmalıdır; ELF'te dalga boyu binlerce kilometre olduğundan, "normal" bir anten imkânsızdır. Tarihsel ELF sistemleri, antenin bir bacağı olarak kilometrelerce uzunlukta gömülü tel hatları ve hatta yerkabuğunun kendisini kullanan devasa tesislerdi; verim yine de son derece düşüktü ve yüksek güç gerekiyordu. İkincisi veri hızı: bant genişliği frekansla orantılı olduğundan, birkaç Hz'lik bir bant yalnızca saniyede bitin kesirleri mertebesinde veri taşıyabilir. Pratikte ELF, bir denizaltıya tam mesaj göndermez; çoğunlukla yalnızca "yüzeye yakın bir derinliğe çık ve daha yüksek frekanslı (VLF/uydu) bir kanaldan tam mesajı al" anlamına gelen çok kısa bir uyandırma/çağrı (bell ringer) gönderir. Bir karakteri iletmek dakikalar alabilir.

> Mühendislik sezgisi: ELF, Shannon kapasitesinin (Bölüm 1) ekstrem bir köşesidir. Kapasite bant genişliğiyle orantılıdır; bant genişliği neredeyse sıfıra giderse, kapasite de öyle. ELF, bu bedeli bilerek öder çünkü tek alternatifi (denizaltının yüzeye çıkması) operasyonel olarak çok daha maliyetlidir. Burada ders şudur: bir kanalı seçerken frekans yalnızca menzil değil, aynı anda nüfuz (cilt derinliği), anten boyutu ve veri kapasitesi arasında zorlu bir takastır.

VLF (3-30 kHz) bu ölçeğin biraz daha yumuşak halidir: hâlâ devasa antenler ve yüksek güç ister ama yüzeye yakın denizaltılara ve dünya çapında (iyonosfer-yer dalga kılavuzu içinde, çok düşük kayıpla yayılarak) ulaşır. VLF ayrıca çok kararlı olduğu için hassas zaman/frekans referansı ve eski navigasyon sistemlerinde kullanılmıştır.

---

<a id="8"></a>
## 8. Navigasyonun Tarihi ve Yeniden Doğuşu: Omega, LORAN, eLORAN

Düşük frekansların klasik bir kullanımı seyrüsefer (navigasyon) idi ve bu alan, GPS çağında "öldü" sanılırken bugün dayanıklılık (resilience) gerekçesiyle yeniden gündemde. Bu, "eski bir teknolojinin neden gelecekte geri dönebileceğine" dair öğretici bir vakadır.

### Eski sistemler: Omega ve LORAN

Omega, VLF bandında çalışan, yalnızca sekiz dünya çapında istasyonla küresel kapsama sağlayan bir navigasyon sistemiydi; alıcı, istasyonlardan gelen sinyallerin faz farkından konum çıkarırdı. Düşük frekansın küresel yayılımı (iyonosfer-yer kılavuzu) sayesinde az istasyonla dünyayı kapsıyordu ama doğruluğu kabaydı. GPS yaygınlaşınca kullanımdan kaldırıldı.

LORAN (Long Range Navigation), LF bandında (klasik LORAN-C ~100 kHz) çalışan, darbe-zamanlamasına dayalı bir hiperbolik navigasyon sistemiydi. Birden çok istasyondan gelen darbelerin varış zamanı farkından (TDOA — Bölüm 9'daki yer tespiti mantığının aynısı) konum bulunurdu. LORAN-C uzun menzilli, güvenilir ve denizcilikte yaygındı; o da GPS sonrası büyük ölçüde kapatıldı.

### eLORAN: PNT yedeği olarak güncel dönüş

Burada güncel ve geleceğe dönük kısım başlar. GPS ve diğer GNSS sistemleri (Bölüm 10) son derece kullanışlıdır ama kritik bir zaafları vardır: sinyalleri çok zayıftır (gürültü tabanının altında) ve bu yüzden karıştırmaya (jamming) ve aldatmaya (spoofing) açıktırlar (Bölüm 10 ve 13). Modern toplum yalnızca konum için değil, hassas zaman (PNT — Positioning, Navigation, Timing) için de GNSS'e bağımlı hale geldi: elektrik şebekeleri, finans zaman damgaları, telekom senkronizasyonu hep GNSS saatine yaslanır. Tek bir noktaya bu kadar bağımlılık, stratejik bir kırılganlıktır.

eLORAN (enhanced LORAN), bu kırılganlığa yanıt olarak klasik LORAN'ın modernize edilmiş halidir: daha yüksek güçlü, daha doğru, ek veri kanalı (zaman dağıtımı) içeren bir LF sistemi. Mantığı şudur: eLORAN, GNSS'ten tamamen bağımsız fiziksel ilkelerle (LF darbe zamanlaması, karasal istasyonlar) çalıştığı için, GNSS karıştırıldığında veya çöktüğünde bir yedek PNT kaynağı sağlar. LF sinyalleri yüksek güçlü ve karıştırması GNSS'e göre çok daha zordur (gürültü tabanının altında değil, güçlü karasal sinyallerdir). Birkaç ülke, kritik altyapı için eLORAN'ı bir GNSS-yedeği olarak değerlendirmekte/konuşlandırmaktadır.

> Güncel not: eLORAN'ın konuşlanma durumu ülkeden ülkeye değişir ve bazı programlar pilot/değerlendirme aşamasındadır; güncel operasyonel durum kaynaktan teyit edilmeli. Kavramsal ders kalıcıdır: tek bir konum/zaman kaynağına (GNSS) toplumsal bağımlılık, bağımsız bir yedek (farklı fizik, farklı bant, farklı altyapı) gerektirir. Bu, geleceğin PNT mimarisinin merkezî bir tasarım ilkesidir.

| Sistem | Bant | İlke | Durum | Rolü |
|---|---|---|---|---|
| Omega | VLF | Faz farkı, küresel | Tarihsel (kapalı) | İlk küresel navigasyon |
| LORAN-C | LF (~100 kHz) | Darbe TDOA, hiperbolik | Büyük ölçüde kapalı | Uzun menzil denizcilik |
| eLORAN | LF | Modernize TDOA + veri | Değerlendirme/pilot/kısmi | GNSS-bağımsız PNT yedeği |
| GNSS (GPS vb.) | L bandı (~1.2-1.6 GHz) | Uydu zamanlama | Yaygın (Bölüm 10) | Birincil PNT (ama kırılgan) |

---

<a id="9"></a>
## 9. Doğal Radyo: Whistler, Sferic, Schumann Rezonansı

Şimdiye kadar insan yapımı sinyallerin egzotik yayılımına baktık. Ama doğa da radyo yayar; üstelik VLF/ELF bandında, herhangi bir verici olmadan, sadece bir basit alıcıyla dinlenebilen olağanüstü sesler üretir. Bu doğal radyo, hem büyüleyici hem de atmosferik fiziğin doğrudan kanıtıdır.

### Sferics: yıldırımın radyo patlaması

Her yıldırım deşarjı, çok geniş bir frekans aralığında (özellikle VLF'te güçlü) anlık bir elektromanyetik darbe yayar. Bu darbelere sferic (atmospheric'in kısaltması) denir. Dünyada her saniye onlarca yıldırım çakar; bu yüzden VLF alıcısı sürekli bir "çıtırtı/tıklama" gürültüsü duyar — her tıklama, dünyanın bir yerindeki bir yıldırımdır. VLF yıldırım izleme ağları (sferic'lerin varış zamanı farkından — yine TDOA) yıldırımları küresel ölçekte gerçek zamanlı haritalandırır.

### Whistler: manyetosferde yolculuk eden ses

Doğal radyonun en çarpıcı olayı whistler'dır. Bir yıldırımın ürettiği VLF enerjisinin bir kısmı, dünyanın manyetik alan çizgileri boyunca manyetosfere (binlerce km yukarı) tırmanır, diğer yarımküreye geçer ve geri döner. Bu uzun yolculukta, plazmanın dağıtıcı (dispersive) özelliği yüzünden yüksek frekanslar düşük frekanslardan daha hızlı ilerler. Sonuç: kulakta yüksek perdeden alçak perdeye kayan, birkaç saniye süren bir ıslık sesi — whistler. Sesteki perde düşüşünün hızı, sinyalin manyetosferde ne kadar yol aldığının ve plazma yoğunluğunun ölçüsüdür; bu yüzden whistler'lar uzay fiziği için bir tanı aracıdır.

```
 Whistler oluşumu (kavramsal):

   yıldırım (VLF darbe)
        │
        ▼  manyetik alan çizgisi boyunca
   ╱────────────────╲
  │  manyetosfer     │   (yüksek frekans önce gelir,
  │  (plazma,        │    düşük frekans sonra → kayan ıslık)
   ╲────────────────╱
        ▲
   diğer yarımkürede (veya geri) alıcı:
   "fiiiiuuuuuw" — yüksekten alçağa kayan ıslık (~1-3 sn)
```

Whistler dinlemek tamamen pasiftir ve özel bir VLF alıcısı (genellikle basit bir uzun tel anten + ses-bandı yükselteci, çünkü VLF frekansları doğrudan ses kartı bandına düşer) ile yapılır. En iyi sonuç elektriksel gürültüden (şehir şebekesi, elektronik) uzakta, kırsalda alınır. Bu, Bölüm 20'deki doğal radyo alıştırmasının konusudur.

### Schumann rezonansı

En düşük frekansta, dünyanın kendisi bir rezonatördür. Dünya yüzeyi ile iyonosfer arasındaki boşluk, küresel bir kapalı dalga kılavuzu/kavite oluşturur. Dünya çapındaki sürekli yıldırım aktivitesi bu kaviteyi uyarır ve belirli rezonans frekanslarında (temel mod yaklaşık 7.83 Hz, ardından ~14, ~20, ~26 Hz harmonikleri) çok zayıf, küresel bir ayakta-duran dalga oluşur. Bunlara Schumann rezonansları denir. Bu, ELF bandının doğal "zili"dir: dünya-iyonosfer kavitesinin doğal frekansı. Ölçülmesi son derece hassas ekipman ve çok düşük gürültülü ortam gerektirir; temel mod frekansındaki küçük kaymalar bile küresel yıldırım aktivitesi ve iyonosfer durumuyla ilişkilendirilir ve bilimsel ilgi konusudur.

> Mühendislik sezgisi: Doğal radyo, "verici olmadan da spektrum doludur" gerçeğinin somut hatırlatıcısıdır. Bir SIGINT/spektrum analisti için bunun iki dersi var. Birincisi, gürültü tabanı boş değildir; doğal kaynaklar (yıldırım, güneş, galaktik gürültü) ona katkı verir ve bunları tanımak, insan yapımı sinyali ayıklamada (Bölüm 7) önemlidir. İkincisi, en basit alıcıyla bile (bir tel ve ses kartı) gözlemlenebilen zengin bir fizik vardır; egzotik spektrumun kapısı pahalı donanım değil, doğru bant ve sessiz bir ortamdır.

---

<a id="10"></a>
## 10. Zayıf-Sinyal Sanatı: WSJT-X, Kodlama Kazancı ve Shannon

Egzotik yayılım modlarının ortak teması — EME'nin aşırı kaybı, meteor scatter'ın kısa açılımı, gri-çizginin değişken penceresi — hep aynı zorluğa çıkar: sinyal çoğu zaman gürültünün içinde veya altındadır. Bu zorluğun modern çözümü, son yıllarda amatör ve profesyonel haberleşmeyi dönüştüren zayıf-sinyal modlarıdır; en bilineni WSJT-X yazılım ailesidir (FT8, FT4, JT65, Q65, MSK144, WSPR ve diğerleri).

### Neden mümkün? Shannon'a geri dönüş

Bölüm 1'de Shannon-Hartley kapasitesini vermiştik:

```
 C = B · log2(1 + S/N)        (C: bit/s, B: bant genişliği Hz, S/N: sinyal/gürültü)
```

Bu formülün gizli dersi şudur: kapasite, SNR çok düşük (negatif dB) olsa bile sıfır değildir — yalnızca çok küçüktür. Yani gürültünün altındaki bir sinyalden de bilgi çıkarmak teorik olarak mümkündür; sadece çok yavaş ve çok az bit ile. Zayıf-sinyal modları tam bu kapıdan girer: veri hızını acımasızca düşürür, mesajı küçültür ve çok güçlü ileri hata düzeltme (FEC) ile her biti gürültüde defalarca "yedekler". Bu yaklaşımla, kulakla hiçbir şey duyulamayan (sinyalin SNR'ı -20 dB'nin altında) bir kanaldan bile yazılım tam bir mesaj çözebilir. Buna kodlama kazancı (coding gain) denir: kanalın ham SNR'ını, kodlama sayesinde etkin olarak yükseltmek.

> Mühendislik sezgisi: Zayıf-sinyal sanatının özü bir takastır — bilgiyi yavaşlatıp küçülterek dayanıklılık satın almak. Bir FT8 mesajı yalnızca ~13 karakter taşır ve gönderimi ~13 saniye sürer; bu çok yavaştır, ama tam da bu yavaşlık ve ağır kodlama, mesajın dünyanın öbür ucundan, gürültünün altından, birkaç watt'la çözülmesini sağlar. Hız ile dayanıklılık arasındaki bu takas, Shannon'ın formülünün pratiğe dökülmüş halidir.

### WSJT-X mod ailesi: her kanal için bir araç

WSJT-X'in farklı modları, farklı kanal zorluklarına ayarlanmıştır. Birini diğerinden ayıran, taşıdığı veri miktarı, gönderim süresi ve hangi bozulmaya (çok düşük SNR, kısa açılım, Doppler) dayandığıdır.

| Mod | Tasarım hedefi | Mesaj/süre | Tipik kullanım | Karakter |
|---|---|---|---|---|
| FT8 | Genel zayıf-sinyal QSO | ~13 karakter / ~15 s | HF DX, günlük zayıf-sinyal | En popüler; çok düşük SNR |
| FT4 | Hızlı kontes/QSO | Kısa / ~6 s | Yoğun bant, hız öncelikli | FT8'den hızlı, biraz az hassas |
| JT65 | Çok zayıf sinyal (EME, HF) | Kısa / ~60 s | EME, derin DX | Çok yavaş, çok hassas |
| Q65 | EME/troposcatter/zor kanal | Kısa / değişken | EME, meteor, ionoscatter | JT65 ardılı; Doppler'e dayanıklı |
| MSK144 | Meteor scatter (kısa burst) | Çok kısa / alt-saniye tekrar | VHF meteor scatter | Hızlı tekrar, burst yakalama |
| WSPR | Yayılım sondalama (beacon) | Çağrı+konum+güç / ~2 dk | Zayıf-sinyal yayılım haritası | Sadece "duyuldum mu" raporu |

### WSPR: yayılımın dünya çapında haritası

Bu ailenin SIGINT/yayılım perspektifinden en ilginç üyesi WSPR'dır (Weak Signal Propagation Reporter — "whisper" diye okunur). WSPR, mesaj alışverişi için değil, yalnızca yayılımı sondalamak için tasarlanmıştır. Bir WSPR vericisi periyodik olarak çok düşük güçle (bazen miliwatt'lar) çağrı işaretini, konumunu (grid locator) ve güç seviyesini yayar. Dünya çapındaki binlerce WSPR alıcısı, duyduğu her WSPR yayınını merkezi bir veritabanına raporlar. Sonuç, gerçek zamanlı, küresel bir "hangi frekansta, hangi noktadan hangi noktaya, hangi güçle yayılım açık?" haritasıdır.

Bir alıcı sahibi için WSPR, kendi alıcısının dünya çapında nereleri duyabildiğini (yani hangi yayılım yollarının kendisine açık olduğunu) tamamen pasif olarak, hiçbir yayın yapmadan görmesini sağlar: sadece WSPR alıcısı kurulur, çözülen yayınlar raporlanır ve harita üzerinde alıcının "duyma menzili" canlanır. Bu, Bölüm 20'deki başlıca alıştırmadır ve egzotik yayılımı (gri-çizgi, gece/gündüz, güneş döngüsü) doğrudan gözlemlemenin en somut, en yasal yoludur.

> Not: WSPR alıcısı çalıştırmak ve sonuçları raporlamak pasiftir ve genellikle serbesttir; WSPR yayını yapmak ise iletim olduğundan amatör telsiz lisansı ve ilgili bandlarda yetki gerektirir. Bu bölümün alıştırmaları yalnızca alıcı (RX) tarafıyla sınırlıdır.

---

<a id="11"></a>
## 11. Mega-Konstelasyonlar: Starlink, OneWeb, Kuiper, Iridium NEXT

Egzotik yayılımdan geleceğe geçişin en görünür köprüsü, gökyüzünü saran mega-konstelasyonlardır. Bölüm 11 uydu haberleşmesinin temellerini verir; burada özellikle son yılların devrimini — alçak yer yörüngesindeki (LEO) binlerce uydudan oluşan internet konstelasyonlarını — spektrum ve SIGINT perspektifinden ele alıyoruz.

### LEO devrimi: neden binlerce uydu?

Geleneksel uydu internetinin (Bölüm 11) çoğu, yer-eşzamanlı yörüngedeki (GEO, ~35786 km) tek bir uyduya dayanırdı. GEO'nun avantajı uydunun gökyüzünde sabit görünmesidir; dezavantajı ise mesafedir: 35786 km gidiş-dönüş, ışık hızında bile ~240 ms'nin üzerinde gecikme (latency) demektir; bu, interaktif uygulamalar (oyun, video görüşme) için rahatsız edicidir. LEO (~340-1200 km) uyduları bu gecikmeyi dramatik biçimde düşürür (tek yön ~milisaniyeler). Ama LEO'nun bedeli vardır: bu kadar alçaktan bir uydu yalnızca küçük bir alanı görür ve gökyüzünde hızla hareket eder (birkaç dakikada ufuktan ufka geçer). Sürekli kapsama için, her an her noktanın üstünde bir uydu olacak şekilde binlerce uydudan oluşan bir konstelasyon gerekir. Mega-konstelasyon kavramı budur.

| Konstelasyon | Yörünge | Yaklaşık ölçek | Bant | Öne çıkan özellik |
|---|---|---|---|---|
| Starlink | LEO (~550 km) | Binlerce (büyüyen) | Ku/Ka (+ uydu-arası lazer) | En büyük; faz dizili kullanıcı terminali |
| OneWeb | LEO (~1200 km) | Yüzlerce-binler | Ku/Ka | Daha yüksek yörünge, kurumsal odak |
| Kuiper | LEO (planlanan) | Binlerce (planlanan) | Ka | Konuşlanma aşamasında; rakip sistem |
| Iridium NEXT | LEO (~780 km) | 66 (+ yedek) | L bandı (+ Ka besleme) | Küresel ses/veri/IoT, çapraz bağlı |

> Ölçek notu: Mega-konstelasyonların uydu sayıları hızla değişir (sürekli fırlatma, planlanan genişlemeler); buradaki sayılar mertebe vermek içindir, güncel kesin rakam kaynaktan teyit edilmeli.

### Faz dizili anten ve hüzme yönlendirme

LEO uyduları hızla hareket ettiğinden, kullanıcı terminali (örneğin Starlink "dishy") uyduyu mekanik olarak takip edemez — çok hızlı olurdu. Bunun yerine elektronik hüzme yönlendirme (beam steering) yapan bir faz dizili anten (phased array) kullanır: yüzlerce küçük anten elemanının fazını ayarlayarak hüzmeyi mekanik hareket olmadan, elektronik olarak gökyüzünde gezdiren bir dizi. Bu, Bölüm 18'de (Massive MIMO ve beamforming) derinlemesine işlenen ilkenin tüketici ölçeğinde somutlaşmış halidir. Terminal, bir uydu ufka inerken hüzmesini bir sonraki yükselen uyduya elektronik olarak "atlatır" (handover).

```
 LEO geçişi ve handover (kullanıcı terminali bakışı):

   gökyüzü:   ☉sat-A →→→        ☉sat-B →→→        ☉sat-C →→→
              (alçalıyor)       (zirvede)         (yükseliyor)

   terminal hüzmesi:    ╲                │                ╱
   (faz dizili,          ╲               │               ╱
    elektronik           ╲              │              ╱
    yönlendirme)          ╲             │             ╱
   ═══════════════════════════════ terminal ═══════════════════
   Her birkaç dakikada hüzme bir uydudan diğerine "handover".
```

### Doppler ve handover: SIGINT/gözlem imzası

LEO uydularının hızlı hareketi, iki belirgin spektral imza yaratır. Birincisi Doppler kayması: uydu yaklaşırken alınan frekans yükselir, uzaklaşırken düşer; bir geçiş boyunca taşıyıcı frekans belirgin bir S-eğrisi çizer. İkincisi handover ritmi: birkaç dakikada bir uydu/hüzme değişimi, trafikte periyodik bir yapı oluşturur. Bir spektrum gözlemcisi için bu imzalar, bir LEO konstelasyon sinyalini diğer (sabit) yayınlardan ayırt etmenin doğal yoludur: hızlı Doppler tarama + periyodik handover = LEO konstelasyon parmak izi.

### Spektrum ve SIGINT etkisi

Mega-konstelasyonların spektrum/SIGINT açısından etkileri çok yönlüdür. Birincisi spektrum yoğunluğu: binlerce uydu, Ku/Ka bandlarını yoğun biçimde kullanır; bu, hem spektrum yönetimi (girişim, koordinasyon) hem de bu bandlardaki diğer servisler için yeni bir ortam yaratır. İkincisi gözlemlenebilirlik: bu uydular hem RF'te (downlink yayınları) hem optik olarak (özellikle Starlink'in "tren" halinde geçen parlak uyduları, fırlatmadan hemen sonra) gözlenebilir. Optik astronomi camiası, mega-konstelasyonların gökyüzü parlaklığına ve teleskop görüntülerine etkisi konusunda ciddi kaygı dile getirmiştir; bu, "gelecek spektrum"un yalnızca RF değil, optik domeni de etkilediğinin somut örneğidir. Üçüncüsü dayanıklılık ve bağımlılık: bu sistemler iletişim altyapısının yeni bir katmanı haline geldikçe, hem stratejik bir varlık hem de yeni bir bağımlılık/hedef yüzeyi oluştururlar.

> Gözlem fırsatı: Bir Starlink "uydu treni"ni (özellikle yeni bir fırlatmadan sonraki ilk günlerde) çıplak gözle veya basit bir teleskople izlemek tamamen yasal ve etkileyici bir gözlemdir; geçiş zamanları çeşitli açık takip araçlarıyla önceden hesaplanabilir. RF tarafında, uygun bandda bir SDR ile LEO downlink'lerinin Doppler-taranan imzasını spektrumda gözlemek de pasif bir alıştırmadır (Bölüm 20).

---

<a id="12"></a>
## 12. Uydu IoT ve Doğrudan-Cihaza Uzay Bağlantısı

Mega-konstelasyonların yanında, geleceğin uzay spektrumunun ikinci büyük dalgası uydu-IoT ve doğrudan-cihaza (direct-to-device, D2D) bağlantıdır. Buradaki fikir, internet hızı sunmak değil, dünyanın her yerindeki düşük-veri cihazlarına (sensörler, takip cihazları, hatta sıradan cep telefonları) uzaydan ince bir bağlantı vermektir.

İki belirgin yön var. Birincisi uydu-IoT: tarım, lojistik, boru hattı, deniz ve uzak bölge sensörleri için tasarlanmış, küçük ve ucuz, çok düşük veri hızlı (zayıf-sinyal modlarına yakın felsefe) LEO konstelasyonları. Bu sistemler, karasal hücresel kapsama olmayan yerlerde (okyanus, çöl, dağ) küçük veri paketlerini (konum, ölçüm) toplar. Çoğu, kısa mesaj ve uzun pil ömrü için optimize edilmiştir; bu da onları zayıf-sinyal/dar-bant tasarım ilkelerine (Bölüm 10) yaklaştırır.

İkincisi — ve daha çarpıcı olanı — doğrudan-cep-telefonuna uydu bağlantısıdır. Son yıllarda hem uydu operatörleri hem hücresel operatörler, değiştirilmemiş sıradan akıllı telefonların doğrudan uyduyla (en azından acil mesaj, SMS, konum paylaşımı düzeyinde) haberleşmesini sağlayan sistemler geliştirmektedir. Buradaki mühendislik meydan okuması devasadır: bir cep telefonunun küçük anteni ve düşük gücü, yüzlerce km uzaktaki bir uyduya ulaşmak için tasarlanmamıştır. Çözüm, uydu tarafında çok büyük antenler (yüksek kazanç), uyduda hassas alıcılar ve yine zayıf-sinyal/ağır-kodlama felsefesidir — yani EME'de ve WSJT'de gördüğümüz aynı ilkelerin (çok büyük anten + çok düşük SNR'da çözebilen kodlama) ticari ölçekte uygulanması.

> Mühendislik sezgisi ve güncellik: Doğrudan-cihaza uydu, bu kitabın yazımı sırasında hızla gelişen, kısmen konuşlanmış ama tam olgunlaşmamış bir alandır; kapsama, veri hızı ve hangi cihazların desteklendiği operatöre ve bölgeye göre değişir ve güncel durum kaynaktan teyit edilmeli. Kavramsal ders kalıcıdır: egzotik yayılımın (EME, meteor scatter) zorladığı "çok zayıf sinyali çözme" sanatı, geleceğin kitlesel uydu bağlantısının temel olanağıdır. Spektrum/SIGINT açısından bu, hücresel bandların bir kısmının uzaya da uzanması ve dünyanın her noktasının (en ıssız yerler dahil) bir RF imza üretebilir hale gelmesi anlamına gelir.

| Yön | Hedef | Veri hızı | Kilit zorluk | Tasarım ilkesi |
|---|---|---|---|---|
| Uydu-IoT | Uzak sensörler/takip | Çok düşük | Pil ömrü, kapsama | Dar-bant, kısa paket, ağır kodlama |
| Direct-to-device | Sıradan telefon | Düşük (mesaj düzeyi) | Telefonun zayıf anteni/gücü | Uyduda dev anten + zayıf-sinyal kodlama |

---

<a id="13"></a>
## 13. Bilişsel Radyo ve Dinamik Spektrum: White Space, CBRS

Geleceğin spektrumu yalnızca yeni yerler (uzay) değil, aynı zamanda mevcut spektrumu daha akıllı kullanmaktır. Klasik spektrum yönetimi statiktir: her band bir servise sabit tahsis edilir (Bölüm 8). Bu, basit ama israflıdır — birçok band çoğu zaman ve çoğu yerde boştur, ama başkası kullanamaz. Bilişsel radyo (cognitive radio) ve dinamik spektrum erişimi (DSA, Dynamic Spectrum Access), bu israfı gidermek için spektrumu statik tahsisten akıllı, koşullu paylaşıma taşıma fikridir.

### Bilişsel radyonun çekirdeği

Bilişsel radyo, çevresindeki spektrumu algılayan (spectrum sensing), boş kanalları bulan ve kendi yayınını birincil kullanıcıya zarar vermeden o boşluklara yerleştiren, gerektiğinde geri çekilen bir radyodur. Çekirdek yetenek spektrum algılamadır: radyonun, bir kanalın gerçekten boş mu yoksa zayıf bir birincil kullanıcı tarafından mı kullanıldığını güvenilir biçimde tespit etmesi. Bu, Bölüm 7'deki enerji tespiti ve çevrimsel durağanlık (cyclostationarity) tekniklerinin tam da uygulama alanıdır: zayıf bir birincil sinyali kaçırmamak (yanlış "boş" kararı zararlıdır) için yapı-arayan duyarlı algılama gerekir.

### TV white space

İlk büyük pratik uygulama TV white space (TVWS) oldu. Dijital TV'ye geçişle, TV bandlarında coğrafi olarak değişen boşluklar (bir bölgede kullanılmayan TV kanalları) oluştu. Bu "beyaz alanlar", elverişli yayılım özellikleri olan UHF frekanslarıdır (iyi nüfuz, iyi menzil). Fikir, ikincil cihazların bu boş TV kanallarını (örneğin kırsal geniş bant için) kullanmasıdır — ama yalnızca o konumda gerçekten boşsa. Boşluğun tespiti iki yolla yapılır: spektrum algılama ve/veya bir coğrafi konum veritabanı sorgusu (cihaz konumunu bildirir, veritabanı o konumda hangi kanalların serbest olduğunu söyler). Veritabanı yaklaşımı pratikte daha güvenilir bulunmuş ve yaygınlaşmıştır.

### CBRS: katmanlı paylaşım

Daha gelişmiş bir model CBRS'tir (Citizens Broadband Radio Service, ~3.5 GHz bandı). CBRS, spektrumu üç öncelik katmanına böler ve merkezi bir koordinatör (Spectrum Access System, SAS) gerçek zamanlı olarak kimin hangi kanalı kullanacağını yönetir:

| Katman | Öncelik | Kim | Davranış |
|---|---|---|---|
| Incumbent | En yüksek | Mevcut kullanıcılar (örneğin deniz radarı) | Daima korunur; varsa diğerleri çekilir |
| PAL | Orta | Lisanslı öncelikli erişim | Açık artırmayla alınan öncelikli haklar |
| GAA | En düşük | Genel yetkili erişim (lisanssız) | Boş kapasiteyi fırsatçı kullanır |

CBRS'in zarafeti, statik tahsis ile lisanssız kaos arasında bir orta yol kurmasıdır: bir koordinatör (SAS), incumbent'ı algılayıp (örneğin kıyıdaki radar aktifse) o bölgedeki ikincil kullanıcıları gerçek zamanlı kanaldan çıkarır. Bu, spektrumu "tahsis edilmiş mülk" yerine "koşullu, paylaşılan, koordine edilen kaynak" olarak gören geleceğin yönetim modelinin öncüsüdür.

> Mühendislik sezgisi: Dinamik spektrum, SIGINT'in bir yeteneğini (spektrum algılama, sinyal varlığı tespiti) bir altyapı hizmetine dönüştürür. "Bu kanalda biri var mı?" sorusu, hem bir dinleyicinin hem bir bilişsel radyonun temel sorusudur; fark, dinleyicinin istihbarat için, bilişsel radyonun nezaket (girişim önleme) için sormasıdır. Geleceğin radyosu, sürekli kendi çevresini dinleyen ve buna göre davranan bir radyodur; bu da onu doğası gereği bir spektrum-farkındalık (ve dolayısıyla SIGINT-komşusu) cihazı yapar.

---

<a id="14"></a>
## 14. Yeni Dalga Formları: OTFS, FBMC, Massive MIMO, mmWave/THz

Modülasyon ve dalga formu da duruyor değil. Bugünün baskın dalga formu OFDM'dir (Orthogonal Frequency Division Multiplexing — 4G/5G, WiFi, dijital TV'nin temeli); ama OFDM'in zayıf yanları (yüksek hareketli/Doppler kanallarda bozulma, yan-lob sızıntısı, tepe/ortalama güç oranı) yeni dalga formu araştırmasını besliyor.

### OTFS: gecikme-Doppler alanı

OFDM, sinyali zaman-frekans ızgarasına yerleştirir; bu, sabit veya yavaş kanallarda mükemmeldir ama hızlı hareket eden (yüksek Doppler) kanallarda — yüksek hızlı tren, uçak, LEO uydu — taşıyıcılar arası girişim yüzünden bozulur. OTFS (Orthogonal Time Frequency Space), bilgiyi zaman-frekans yerine gecikme-Doppler (delay-Doppler) alanına yerleştiren bir dalga formu kavramıdır. Fikir şudur: hızlı bir kanalın etkisi (çok-yol gecikmesi ve Doppler) gecikme-Doppler alanında çok daha durağan ve seyrek görünür; sinyali bu alanda taşımak, yüksek-Doppler kanallara doğal bir dayanıklılık verir. OTFS, özellikle yüksek hareketlilik (high-mobility) senaryoları — hızlı araçlar, uydu — için umut verici bir aday olarak araştırılmaktadır.

> Olgunluk notu: OTFS aktif bir araştırma ve standartlaşma-öncesi tartışma konusudur; kitlesel ticari konuşlanması henüz yerleşmemiştir ve performans iddiaları senaryoya bağlıdır, kaynaktan teyit edilmeli. Kavramsal değeri nettir: hareketlilik arttıkça (uydu, yüksek hızlı ulaşım), dalga formunu Doppler'e dayanıklı bir alana taşımak mantıklı bir yöndür.

### FBMC ve diğer çok-taşıyıcı adaylar

FBMC (Filter Bank Multi-Carrier), OFDM'in yan-lob sızıntısı (her alt-taşıyıcının komşu bandlara sızması) sorununu, her alt-taşıyıcıyı keskin bir filtreyle şekillendirerek azaltmayı amaçlayan bir alternatiftir. Daha temiz spektral şekil, dinamik spektrum (Bölüm 13) gibi parçalı, boşluk-doldurma senaryolarında değerlidir çünkü komşu birincil kullanıcıya daha az sızar. FBMC ve akrabaları (GFDM, UFMC) çeşitli avantaj/dezavantaj takasları sunar ve hepsi araştırma/erken-değerlendirme aşamasındadır; hiçbiri OFDM'i kitlesel olarak henüz tahtından etmemiştir.

| Dalga formu | Ana fikir | Güçlü yanı | Durum |
|---|---|---|---|
| OFDM | Zaman-frekans ortogonal taşıyıcılar | Olgun, basit denkleştirme | Baskın (4G/5G/WiFi) |
| OTFS | Gecikme-Doppler alanı | Yüksek Doppler dayanıklılığı | Araştırma/aday (uydu, hız) |
| FBMC | Filtrelenmiş alt-taşıyıcılar | Düşük yan-lob sızıntısı | Araştırma/değerlendirme |

### Massive MIMO ve beamforming (derin bağ)

5G ve ötesinin performans sıçramasının büyük kısmı yeni modülasyondan değil, uzamsal işlemeden gelir: Massive MIMO (çok sayıda anten elemanı) ve beamforming (hüzme şekillendirme). Çok sayıda antenle, baz istasyonu enerjiyi belirli kullanıcılara yönlü hüzmelerle gönderebilir (spatial multiplexing — aynı frekansı farklı yönlerdeki kullanıcılara aynı anda kullandırma) ve girişimi bastırabilir. Bu, Bölüm 18'de derinlemesine işlenen konudur; burada vurgu, geleceğin kapasitesinin büyük ölçüde uzamsal boyuttan (anten dizisi + akıllı hüzme) geleceğidir. Bölüm 11'deki LEO faz dizili terminaller ve Bölüm 13'teki bilişsel radyo, hep bu uzamsal-işleme devriminin farklı yüzleridir.

### mmWave ve THz: spektrumun üst sınırına tırmanmak

Veri talebi arttıkça, daha fazla bant genişliği için spektrumun üst uçlarına tırmanılıyor. Milimetre dalga (mmWave, kabaca 24-100 GHz) 5G'nin yüksek-kapasite katmanında zaten kullanımda; çok geniş bantlar (dolayısıyla yüksek veri hızı) sunar ama menzili kısa ve engellere (duvar, yağmur, hatta el) çok duyarlıdır. Bunun ötesinde terahertz (THz, kabaca 0.1-10 THz) haberleşme, 6G ve sonrası için araştırılan bir sınırdır: muazzam bant genişliği vaadi, ama çok kısa menzil, atmosferik soğurma (özellikle su buharı emilim hatları) ve donanım zorlukları. THz, optik ile RF arasındaki sınır bölgesidir ve hâlâ büyük ölçüde laboratuvar/araştırma aşamasındadır.

```
 Frekans-bant genişliği-menzil takası (kabaca):

   düşük frekans ◄──────────────────────────────► yüksek frekans
   (HF/VHF/UHF)        (mikrodalga)      (mmWave)      (THz)
   az bant genişliği   orta              çok geniş     muazzam
   uzun menzil         orta menzil       kısa menzil   çok kısa
   engel-geçer         orta              engele duyarlı  her şey engel
   olgun               olgun             konuşlanıyor   araştırma
```

> Mühendislik sezgisi: Spektrumun üstüne tırmanmak (mmWave → THz) bant genişliği kazandırır ama fiziği acımasızlaştırır: menzil kısalır, engeller öldürür, donanım zorlaşır. Bu yüzden gelecek, "tek bir sihirli frekans" değil, katmanlı bir mimaridir — alçak frekanslar kapsama/dayanıklılık için, yüksek frekanslar yoğun-kapasite noktaları için, ve hepsini akıllı beamforming + dinamik spektrum yönetimi birbirine bağlar. THz'nin haberleşmedeki rolü ve algılama (sensing) ile birleşmesi (joint communication and sensing) araştırma aşamasındadır ve kaynaktan teyit edilmeli.

---

<a id="15"></a>
## 15. Görünür Işık Haberleşmesi (Li-Fi) ve Optik/Lazer Bağlar

Spektrumun üst sınırına tırmanmanın doğal devamı, radyoyu tamamen aşıp ışığa geçmektir. Görünür ışık, kızılötesi ve lazer, elektromanyetik spektrumun RF'in çok ötesindeki bölgeleridir ve haberleşme için giderek daha fazla kullanılıyor.

### Li-Fi: ışıkla veri

Li-Fi (Light Fidelity), görünür ışık haberleşmesinin (VLC, Visible Light Communication) bir markasıdır: bir LED'i insan gözünün algılayamayacağı kadar hızlı yakıp söndürerek (ya da yoğunluğunu modüle ederek) veri taşımak. Aydınlatma ve haberleşme aynı LED'den gelir; ışık zaten ortamı aydınlatırken, üzerine binen yüksek hızlı modülasyon veriyi taşır. Avantajları ilginçtir: görünür ışık duvarları geçmez, dolayısıyla sinyal bir odanın içinde kalır — bu hem bir güvenlik avantajı (oda dışından dinlenemez; fiziksel sınırlama) hem RF spektrum sıkışıklığından kaçış sağlar. Dezavantajları da bariz: ışık engelleri geçemez (bir nesne araya girince bağlantı kopar), gün ışığı/ortam ışığı girişimi yapar ve genellikle tek yönlü ya da kısa menzillidir.

> Güvenlik sezgisi: Li-Fi, "yayılımın fiziksel sınırlanması" fikrinin güzel bir örneğidir. RF her yöne sızarken (Bölüm 6'daki TEMPEST/OPSEC kaygısı), görünür ışık duvarda durur. Bu, bazı senaryolarda (kapalı, kontrollü oda içi bağlantı) doğal bir gizlilik katmanı sunar — ama "ışık geçmez" varsayımı pencere, kapı aralığı ve yansımalarla zayıflar; mutlak değildir ve gerçek güvenlik için yine şifreleme gerekir.

### Serbest-uzay optik ve lazer bağlar

Daha uzun menzil için serbest-uzay optik haberleşme (FSO, Free-Space Optics) ve lazer bağlar kullanılır: iki nokta arasında dar bir lazer hüzmesiyle, fiber çekmeden, çok yüksek hızlı bir bağ kurmak. Avantajı muazzam bant genişliği ve dar hüzme (yönlü, dinlemesi/karıştırması zor); dezavantajı görüş hattı zorunluluğu ve atmosfere (sis, yağmur, türbülans, hizalama) aşırı duyarlılık. FSO, binalar arası yüksek hızlı bağ ve özellikle uzayda öne çıkar.

### Uzay-içi lazer ve derin-uzay optik haberleşme

Geleceğin en heyecan verici optik uygulaması uzaydadır. Uydu-arası lazer bağlar (örneğin Starlink'in uydular arası lazer ağı, Bölüm 11) uyduları RF yerine ışıkla birbirine bağlar; bu, çok yüksek hız sağlar ve yer istasyonu ihtiyacını azaltır (veri uzayda uydudan uyuya optik olarak taşınır). Daha da ileride, derin-uzay optik haberleşme (deep-space optical communication) — uzak uzay araçlarından lazerle veri göndermek — geleneksel RF derin-uzay bağlarına göre çok daha yüksek veri hızı vaat eder ve test/gösterim aşamasındadır. Burada zorluklar EME'yi (Bölüm 5) andırır ama optik domende: aşırı mesafe, çok dar hüzmeyi milyonlarca km öteden bir alıcıya isabet ettirme (hassas yönlendirme), ve çok zayıf foton-seviyesi sinyali algılama.

> Olgunluk notu: Uydu-arası lazer ağları kısmen konuşlanmıştır; derin-uzay optik haberleşme aktif test/gösterim aşamasındadır ve operasyonel olgunluğu kaynaktan teyit edilmeli. Kavramsal ders: spektrum stratejisi artık yalnızca "hangi RF bandı" değil, "RF mı yoksa optik mi" sorusunu da içeriyor; optik, dar-hüzme/yüksek-hız/fiziksel-gizlilik istenen yerlerde RF'e güçlü bir tamamlayıcı (ve bazen alternatif) oluyor.

| Optik mod | Menzil | Bant genişliği | Kilit zorluk | Gizlilik özelliği |
|---|---|---|---|---|
| Li-Fi / VLC | Oda içi | Yüksek | Engel/ortam ışığı | Duvarı geçmez (oda-içi gizli) |
| FSO (karasal) | Bina-bina (km) | Çok yüksek | Sis/yağmur/hizalama | Dar hüzme, dinlemesi zor |
| Uydu-arası lazer | Uzay (binlerce km) | Çok yüksek | Hassas yönlendirme | Uzayda, erişimi zor |
| Derin-uzay optik | Milyonlarca km | Yüksek (RF'e göre) | Foton-seviyesi algılama | Aşırı dar hüzme |

---

<a id="16"></a>
## 16. Güvenliğin Geleceği: Post-Kuantum, QKD ve Kuantum Algılama

Geleceğin SIGINT'i yalnızca yeni sinyaller değil, yeni güvenlik paradigmalarıdır. Kuantum hesaplamanın yükselişi, hem bir tehdit (mevcut şifrelemeyi kırma potansiyeli) hem de yeni savunma araçları (kuantum anahtar dağıtımı, kuantum algılama) doğuruyor. Bu başlık RF/spektrum bağlamında bu üç ekseni ele alır.

### Post-kuantum kriptografi (PQC): RF bağlamı

Bugünün açık-anahtar kriptografisinin (örneğin RSA, eliptik eğri) güvenliği, belirli matematik problemlerinin (büyük sayı çarpanlara ayırma, ayrık logaritma) klasik bilgisayarlarla pratik olarak çözülemez olmasına dayanır. Yeterince büyük ve kararlı bir kuantum bilgisayar, bilinen kuantum algoritmalarıyla bu problemleri çok daha hızlı çözebilir; bu da bugünkü açık-anahtar şemalarını teorik olarak kırılabilir kılar. Böyle bir kuantum bilgisayarın ne zaman (ya da pratikte hiç) ortaya çıkacağı belirsizdir ve aktif tartışma konusudur; ama tehdit yeterince ciddiye alınmaktadır.

Yanıt post-kuantum kriptografidir (PQC): kuantum bilgisayarların da (bilindiği kadarıyla) verimli çözemediği matematik problemlerine dayanan yeni algoritmalar (örneğin kafes-tabanlı, kod-tabanlı, hash-tabanlı şemalar). Standartlaştırma süreci olgunlaşmış ve ilk PQC algoritmaları standart olarak yayımlanmış durumdadır (kesin algoritma listesi ve durum kaynaktan teyit edilmeli). RF/haberleşme bağlamında PQC'nin önemi şudur: kablosuz cihazlar (IoT, uydu, hücresel) uzun ömürlüdür ve bugün gönderilen (ve belki kaydedilen) şifreli trafik, gelecekte kuantum yetenekle çözülebilir ("harvest now, decrypt later" — şimdi topla, sonra çöz tehdidi). Bu yüzden kritik kablosuz sistemlerin anahtar değişimini PQC'ye geçirmesi, gelecekteki bir tehdide karşı bugünden alınan bir önlemdir.

> Savunma sezgisi: "Şimdi topla, sonra çöz" tehdidi, Bölüm 11'deki trafik analizi dersinin kuantum-çağı uzantısıdır. Bir saldırgan içeriği bugün çözemese bile, şifreli trafiği kaydedip gelecekte çözmeyi umabilir. Bu, uzun-ömürlü/yüksek-değerli kablosuz trafiğin neden bugün PQC'ye geçmesi gerektiğini açıklar: tehdit gelecekte, ama korunması gereken veri bugün havada.

### Kuantum anahtar dağıtımı (QKD): kavram

QKD (Quantum Key Distribution), şifreleme anahtarını klasik matematiğe değil, kuantum fiziğinin temel yasalarına dayanarak güvenli paylaşmayı amaçlar. Temel fikir: kuantum durumlarını (örneğin tek fotonların polarizasyonunu) ölçmek, onları kaçınılmaz olarak bozar (ölçüm geri-etkisi, no-cloning teoremi). Bu yüzden bir dinleyici (eavesdropper) kanalı dinlemeye çalışırsa, iletilen kuantum durumlarını bozar ve bu bozulma taraflarca istatistiksel olarak tespit edilebilir. Yani QKD'nin güvenlik vaadi, "dinleme fiziksel olarak iz bırakır ve fark edilir" ilkesine dayanır.

QKD bugün çoğunlukla optik fiber üzerinden (sınırlı mesafede) ve deneysel olarak uydu-yer (serbest-uzay optik) bağlarında gösterilmiştir; uydu-tabanlı QKD, fiberin mesafe sınırını aşmak için araştırılan bir yöndür. Ancak QKD pratikte önemli kısıtlara sahiptir: özel donanım gerektirir, mesafe/hız sınırlıdır, ve yalnızca anahtar dağıtımını çözer (kimlik doğrulama hâlâ ayrı gerekir). Bu yüzden QKD'nin gerçek-dünya rolü ve PQC'ye kıyasla pratikliği aktif tartışma konusudur; bazı uzmanlar yakın gelecekte PQC'yi daha pratik bulur.

> Olgunluk notu: QKD laboratuvar ve sınırlı saha gösterimlerinde çalışır; geniş, pratik konuşlanması ve maliyet/fayda dengesi araştırma ve tartışma aşamasındadır, kaynaktan teyit edilmeli. RF/SIGINT bağlamında doğrudan etkisi (çoğunlukla optik olduğu için) sınırlıdır, ama "dinlemenin fiziksel olarak tespit edilebildiği bir kanal" kavramı, klasik RF güvenliğinin (dinleme genelde sessiz ve tespit edilemez) tam tersi olduğu için kavramsal olarak öğreticidir.

### Kuantum algılama (quantum sensing): kavram

Kuantum fiziğinin SIGINT'e belki daha doğrudan dokunan yüzü kuantum algılamadır. Kuantum sistemlerin (örneğin atomik enerji seviyeleri, NV-merkezleri, Rydberg atomları) çevresel alanlara olağanüstü hassasiyeti, çok zayıf elektromanyetik alanları ölçen yeni nesil sensörler vaat eder. Özellikle araştırılan bir yön Rydberg-atom tabanlı RF algılamadır: yüksek uyarılmış atomların elektrik alanına aşırı duyarlılığını kullanarak, geleneksel antenlerden farklı ilkelerle ve potansiyel olarak çok geniş bantta RF alan ölçmek. Eğer olgunlaşırsa, böyle sensörler alıcı tasarımında (özellikle çok zayıf sinyal ve geniş-bant algılama) yeni olanaklar açabilir.

> Olgunluk notu: Kuantum RF algılama (Rydberg vb.) aktif araştırma aşamasındadır; pratik, alan-konuşlanmış SIGINT alıcısı olarak olgunluğu henüz kanıtlanmamıştır ve iddialar kaynaktan teyit edilmeli. Kavramsal ilgi: bu teknolojiler, "antenin ne olduğu" tanımını bile değiştirebilir — bir kristal/atom bulutu, geleneksel bir metal anten yerine alan algılayıcısı olabilir. Bu, zayıf-sinyal sanatının (Bölüm 10) donanım tarafındaki olası bir geleceğidir.

| Kuantum ekseni | Ne yapar | RF/SIGINT bağlamı | Olgunluk |
|---|---|---|---|
| PQC | Kuantum-dirençli şifreleme | Kablosuz trafiği gelecekteki çözmeye karşı korur | Standartlar yayımlandı, geçiş sürüyor |
| QKD | Fiziğe dayalı anahtar dağıtımı | Dinleme tespit edilebilir kanal (çoğunlukla optik) | Saha gösterimi; pratiklik tartışmalı |
| Kuantum algılama | Aşırı hassas alan ölçümü | Yeni nesil RF alıcı/anten kavramı | Araştırma aşamasında |

---

<a id="17"></a>
## 17. Yapay Zekâ-Tanımlı Radyo ve Kendini Ayarlayan Spektrum

Geleceğin radyosunun belki en kapsayıcı temması, yazılım-tanımlı radyonun (SDR, Bölüm 2) bir adım ötesi: yapay zekâ-tanımlı radyo (AI-defined / cognitive radio'nun öğrenen hali). Bölüm 7'de otomatik modülasyon tanıma (AMC) ve makine öğrenmesinin sinyal sınıflandırmadaki rolünü işledik; burada bunu bir bütüne, kendi davranışını öğrenen ve uyarlayan radyoya genişletiyoruz. (Bu konunun derinlemesine ele alınışı Bölüm 19'da planlanmaktadır; burada gelecek-spektrum bağlamında özetlenir.)

Klasik bir radyo sabit parametrelerle çalışır: belirli bir modülasyon, belirli bir bant, belirli bir güç. Yazılım-tanımlı radyo bunları yazılımdan değiştirilebilir kılar. Yapay zekâ-tanımlı radyo bir adım daha gider: radyo, çevresini (spektrum durumu, girişim, kanal kalitesi) sürekli algılar ve en iyi parametreleri (frekans, dalga formu, güç, kodlama) gerçek zamanlı kendisi seçer/öğrenir. Bu, üç yetenek katmanının birleşimidir:

| Katman | Yetenek | İlgili bölüm |
|---|---|---|
| Algılama | Spektrum/sinyal durumunu tanıma (AMC, enerji/yapı tespiti) | Bölüm 7, 13 |
| Karar | En iyi dalga formu/frekans/gücü seçme (öğrenme, optimizasyon) | Bu bölüm, Bölüm 19 |
| Uyum | Donanımı (SDR) yeni parametrelere anında ayarlama | Bölüm 2 |

Bunun SIGINT için iki yüzü vardır. Savunma/dost tarafta: kendi kendine uyum sağlayan radyo, karıştırmadan kaçabilir (jamming algılayınca frekans/dalga formu değiştirir), girişimi en aza indirir ve zorlu kanallarda dayanıklılığı artırır — Bölüm 13'teki dinamik spektrumun öğrenen, otonom hali. Karşı/analiz tarafta ise zorluk büyür: eğer hedef radyo dalga formunu ve frekansını sürekli, akıllıca değiştiriyorsa, onu ayıklamak ve sınıflandırmak (Bölüm 7) çok daha zordur; LPI/LPD felsefesi (Bölüm 7, sonraki bölümlerde de işlenir) yapay zekâ ile güçlenir. Yani yapay zekâ hem yayıcıyı daha kaçışkan yapar hem de analizciye daha güçlü ayıklama araçları (öğrenen sınıflandırıcılar) verir; bu, sürekli ilerleyen bir karşılıklı tırmanıştır.

> Mühendislik sezgisi: Geleceğin radyosu "kör ve sabit" değil, "gören ve uyarlanan" bir cihazdır. Bu, spektrumu statik bir tahsis tablosundan (Bölüm 8) yaşayan, sürekli yeniden müzakere edilen bir ekosisteme dönüştürür. Hem dost radyo hem analizci aynı temel yeteneğe — spektrumu algılayıp anlamlandırma — yaslanır; bu yüzden gelecekte radyo mühendisliği ile sinyal istihbaratı, ortak bir algılama-öğrenme çekirdeği etrafında giderek daha çok iç içe geçer.

---

<a id="18"></a>
## 18. Sıra Dışı ve Merak: Pirat Radyo, Gizli Kanallar, Astronomik Kaynaklar

Egzotik spektrumun bir de "merak köşesi" vardır: ana akım dışı, sıra dışı ama öğretici sinyal kaynakları. Bunları tanımak, hem spektrum okuryazarlığını zenginleştirir hem de bilinmeyen bir sinyalle karşılaşma (Bölüm 19) sezgisini besler.

### Pirat radyo ve kaçak vericiler

Pirat radyo, yetkisiz (lisanssız) yayın yapan vericilere verilen addır; tarihsel olarak kıyı ötesi gemilerden FM müzik yayını yapan istasyonlardan, kısa dalgada belirli zamanlarda beliren kaçak yayıncılara kadar uzanır. Bir spektrum gözlemcisi açısından bunlar, "bant planında (Bölüm 8) olmaması gereken yerde beliren sinyal" örnekleridir ve tam da bu yüzden anormallik-tespiti pratiği için ilginçtir. Pirat/kaçak vericilerin tespiti, Bölüm 9'daki yön bulma (DF) tekniklerinin klasik bir sivil uygulamasıdır (yetkili kurumlar yetkisiz vericileri DF ile bulur). Burada vurgu yalnızca tanımadır; bu kitap hiçbir yetkisiz yayını teşvik etmez — aksine, "spektrumda olmaması gerekeni fark etme" savunma refleksini örnekler.

### Gizli kanallar (covert channels)

Daha incelikli bir konu gizli kanallardır: bilginin, beklenmedik bir taşıyıcıya gizlice gömülerek iletilmesi. RF bağlamında bu, meşru bir sinyalin içine fark edilmeyecek biçimde ek bilgi sıkıştırmak (steganografik modülasyon), ya da normalde haberleşme için tasarlanmamış bir emisyonu (örneğin bir cihazın istem dışı RF sızıntısı — Bölüm 6'daki TEMPEST'in tersi) bilgi taşımak için kullanmak olabilir. Gizli kanallar hem bir tehdit (veri sızdırma yolu) hem bir araştırma konusudur. Savunma perspektifinden önemi, "beklenen sinyalin içinde beklenmeyen bir yapı" aramayı gerektirmesidir — ki bu, çevrimsel durağanlık ve ince istatistiksel analiz (Bölüm 7) ile yapılır. Bu kitap gizli kanal kurmayı değil, varlığını fark etme ve savunma bilincini hedefler.

### Astronomik ve doğal radyo kaynakları

Spektrumun en sıra dışı yayıcıları dünyada değil, gökyüzündedir. Radyo astronomi, evrenin doğal radyo kaynaklarını inceler ve bunların bir kısmı SIGINT/spektrum dünyasıyla kesişir:

| Astronomik kaynak | Doğası | Radyo imzası | SIGINT/spektrum kesişimi |
|---|---|---|---|
| Güneş | En yakın yıldız | Güneş radyo patlamaları (geniş bant gürültü) | Patlamalar HF-mikrodalgada girişim/karıştırma yapar |
| Pulsarlar | Dönen nötron yıldızları | Olağanüstü düzenli darbeler (saniye-altı PRI gibi) | Doğal "darbe treni"; zamanlama referansı kavramı |
| Galaktik gürültü | Samanyolu arka planı | Geniş-bant gürültü tabanı (özellikle düşük frekans) | Alıcının gürültü tabanına katkı (Bölüm 1) |
| Jüpiter | Manyetosferik emisyon | Dekametrik (HF) patlamalar | HF'te doğal güçlü emisyon kaynağı |

Pulsarlar özellikle ilginçtir: dönen bir nötron yıldızının yaydığı radyo hüzmesi, dünyaya olağanüstü düzenli aralıklarla ulaşır — adeta evrenin en kararlı darbe treni (Bölüm 5'teki PRI kavramının doğal, kozmik hali). Pulsar zamanlamasının kararlılığı o kadar yüksektir ki, kavramsal olarak bir zaman/navigasyon referansı (pulsar tabanlı navigasyon) olarak araştırılmaktadır; olgunluğu araştırma aşamasındadır ve kaynaktan teyit edilmeli. Güneş radyo patlamaları ise pratik bir yan etkiye sahiptir: güçlü bir patlama, GNSS dahil (Bölüm 10) geniş bir bandda doğal girişim/gürültü yükselmesine yol açabilir — yani gökyüzü zaman zaman istemeden bir "karıştırıcı" olur.

### SETI: en uç merak

En uç sıra dışı sinyal arayışı SETI'dir (Search for Extraterrestrial Intelligence — Dünya Dışı Zekâ Arayışı): doğal kaynaklarla açıklanamayan, yapay (teknolojik) kökenli olabilecek dar-bant veya yapılı radyo sinyallerini gökyüzünde aramak. SETI, kavramsal olarak bu kitabın çekirdek problemiyle aynıdır: gürültü ve doğal kaynaklar arasından yapay/anlamlı bir sinyali ayırt etmek (Bölüm 7'deki ayıklama ve sınıflandırmanın kozmik ölçekte hali). SETI'nin temel sezgisi, "doğal süreçler genellikle geniş-bant ve gürültü-benzeridir; çok dar-bant veya açıkça yapılandırılmış bir sinyal yapaylığa işaret edebilir" varsayımıdır — ki bu, modülasyon tanımanın (Bölüm 7) en felsefi formudur. SETI bugüne kadar doğrulanmış bir tespit yapmamıştır ve bilimsel olarak açık bir arayıştır.

> Mühendislik sezgisi: Bu merak köşesinin ortak dersi, spektrumun her ölçekte — yetkisiz bir yerel verici, gizli bir kanal, bir pulsarın darbesi, hatta varsayımsal bir uzaylı işareti — "beklenen ile beklenmeyeni ayırma" problemi olmasıdır. Bir analizcinin temel refleksi her durumda aynıdır: önce neyin normal/beklenen olduğunu bil (bant planı, doğal gürültü tabanı, bilinen kaynaklar), sonra ondan sapanı fark et. Bu refleks, bir sonraki ve son teknik başlığın — bilinmeyen sinyalle karşılaşma metodolojisinin — çekirdeğidir.

---

<a id="19"></a>
## 19. Bilinmeyen Sinyalle Karşılaşma Metodolojisi (Kapanış Sentezi)

Bu bölüm, egzotikten geleceğe geniş bir yelpaze gezdi. Şimdi hepsini tek bir pratik soruda toplayalım, çünkü bir RF araştırmacısının gerçek işi budur: spektrumda hiç tanımadığın bir sinyalle karşılaştın. Ne yaparsın? Bu, Bölüm 5'teki (radar/darbe analizi) ve Bölüm 7'deki (ayıklama, AMC, kütüphane eşleştirme) yöntemlerin pratik bir sentezidir ve serinin bu noktaya kadar verdiği her şeyi bir akışta birleştirir.

Bilinmeyen bir sinyali çözmek doğrusal bir dedektiflik sürecidir: gözlemden hipoteze, hipotezden teyide. Aşağıdaki akış, bir araştırmacının zihninde (çoğu zaman farkında olmadan) işleyen sıradır.

```
 Bilinmeyen sinyalle karşılaşma akışı:

 [1] NEREDE?  → Frekans ve bant planı (Bölüm 8)
       │         "Bu bandda ne bulunması beklenir?"
       ▼
 [2] NE ZAMAN/NASIL? → Zaman davranışı (sürekli/darbeli/burst/atlamalı)
       │              Sky-wave mı? Es açılımı mı? LEO Doppler mi? (Bölüm 2-6,11)
       ▼
 [3] NE BİÇİMDE? → Modülasyon tanıma (AMC: neyin sabit/değişken) (Bölüm 7)
       │           Genlik/faz/frekans hangisi taşıyor?
       ▼
 [4] HANGİ RİTİMDE? → PRI/sembol hızı/periyodik yapı (Bölüm 5,7)
       │              Çevrimsel durağanlık ile saklı imza
       ▼
 [5] NEREDEN? → Yön bulma / geliş açısı (DF/AOA) (Bölüm 9)
       │         Sabit konum mu, hareketli (uydu) mi?
       ▼
 [6] EŞLEŞTİR → Açık veritabanı (sigidwiki vb.) + kütüphane (Bölüm 7,10)
       │          "Bu imza bilinen bir kayda uyuyor mu?"
       ▼
 [7] HİPOTEZ + TEYİT → Bir tahmin kur, bağımsız kanıtla doğrula/çürüt
                        Uymuyorsa: yeni/bilinmeyen → kaydet, geri besle
```

Bu akışın her adımı serinin bir bölümüne bağlanır ve her biri bir soruyu eler:

| Adım | Soru | Araç/Bölüm | Eler |
|---|---|---|---|
| Nerede | Hangi frekans/band? | Bant planı (Bölüm 8) | Beklenen servis sınıfı |
| Ne zaman | Zaman deseni ne? | Waterfall gözlemi (Bölüm 2), yayılım (Bu bölüm) | Yayılım modu, sürekli/burst |
| Ne biçimde | Hangi modülasyon? | AMC, "sabit/değişken" (Bölüm 7) | Modülasyon ailesi |
| Hangi ritimde | Periyodik yapı? | PRI/cyclostationarity (Bölüm 5,7) | Sistem tipi/sembol hızı |
| Nereden | Geliş açısı/konum? | DF/AOA (Bölüm 9) | Sabit/hareketli kaynak |
| Eşleştir | Bilinen mi? | Veritabanı/kütüphane (Bölüm 7,10) | Bilinen vs yeni |

> Mühendislik sezgisi (sentez): Bilinmeyen sinyal çözme, tek bir sihirli ölçüm değil, birbirini daraltan kanıtların kümülatif birikimidir. Hiçbir adım tek başına kesin değildir — ama frekans + zaman deseni + modülasyon + ritim + geliş açısı birlikte, olasılık uzayını tek bir hipoteze kadar daraltır. Bu, Bölüm 7'deki "ayırt edici güç tek parametrede değil, parametrelerin birlikte oluşturduğu örüntüdedir" ilkesinin nihai ifadesidir. Ve en önemlisi: hipotezini her zaman bağımsız bir kanıtla teyit et; "böyle görünüyor" ile "böyle olduğunu doğruladım" arasındaki fark, meraklı ile araştırmacı arasındaki farktır.

Egzotik yayılım bu metodolojiye bir boyut daha ekler. Bilinmeyen bir sinyal "olmaması gereken" bir yerden geliyorsa, ilk soru içerik değil yayılımdır: Bu sinyal gerçekten yerel mi, yoksa bir egzotik mod (sky-wave, Es, ducting, meteor scatter, LEO geçişi) tarafından çok uzaktan mı getirildi? Bir VHF bandında aniden beliren uzak bir istasyon büyük olasılıkla yeni bir verici değil, bir Sporadic-E açılımıdır (Bölüm 3); hızlı Doppler-taranan bir sinyal yeni bir karasal kaynak değil, bir LEO uydusudur (Bölüm 11). Yayılımı anlamak, "bu ne?" sorusundan önce "bu nasıl buraya geldi?" sorusunu cevaplamayı sağlar ve çoğu zaman gizemi açıklar.

---

<a id="20"></a>
## 20. Alıştırmalar (Yasal, Yalnızca Alıcı ve Gözlem)

> Bu alıştırmalar yalnızca alıcı (RX), pasif gözlem ve hesap içindir. Hiçbiri iletim gerektirmez. WSPR/zayıf-sinyal alıcısı çalıştırmak ve raporlamak pasiftir; bu modlarda yayın yapmak ise lisans ve yetki gerektirir ve bu alıştırmaların kapsamı dışındadır. Doğal radyo ve uydu gözlemi de tamamen pasiftir. Şüphedeysen yapma; bandını ve ülkenin mevzuatını teyit et.

### A) WSPR ile kendi alıcının dünya çapındaki "duyma menzilini" haritalamak (sadece RX)

Bu, egzotik yayılımı doğrudan gözlemlemenin en somut yoludur. Bir SDR (Bölüm 2) ve WSPR çözücü yazılım kurarak, bir HF WSPR frekansını (örneğin 14 MHz dolayındaki standart WSPR alt-bandı; güncel frekansları kaynaktan teyit et) dinle ve çözdüğün yayınları açık WSPR raporlama ağına ilet (yalnızca alıcı; hiçbir yayın yapmıyorsun).

| Gözlem penceresi | Hangi noktalar duyuldu? | Mesafe (kabaca) | Yayılım yorumu |
|---|---|---|---|
| Gündüz (öğle) | ? | ? | Yüksek band açık mı? D-yutması? |
| Gün batımı (gri-çizgi) | ? | ? | Gri-çizgi pencere açıldı mı? |
| Gece | ? | ? | Düşük band, uzun mesafe? |

Birkaç gün boyunca farklı saatlerde dinle ve "alıcımın duyma menzili" haritasının nasıl değiştiğini izle. Amaç: gece/gündüz farkını (Bölüm 2 iyonosfer), gri-çizgi penceresini (Bölüm 2) ve mevcut güneş döngüsü koşulunu (Bölüm 2) kendi alıcının verisiyle, tamamen pasif olarak gözlemlemek. İleri adım: aynı uzak istasyonun sinyal seviyesinin (raporlanan SNR) gün içinde nasıl değiştiğini not et; bu, yayılımın canlı nabzıdır.

### B) Bir Sporadic-E / VHF açılımını gözlemlemek (pasif dinleme)

Yaz aylarında (kuzey yarımkürede Mayıs-Ağustos zirvesi), normalde yalnızca yerel istasyonların duyulduğu FM yayın bandını (88-108 MHz) ya da bir VHF amatör bandını periyodik dinle. Bir Es açılımının imzasını ara: normalde ölü olan bantta aniden uzak şehirlerden güçlü istasyonların girip çıkması.

```
 Gözlem kaydı:
 Tarih/saat: ______   Band: ______
 Normalde duyulanlar: ______ (yerel)
 Açılımda beliren uzak istasyonlar: ______
 Süre (ne kadar açık kaldı): ______
 Karakter (ani/güçlü mü, yavaş/zayıf mı): ______
```

Amaç: Bölüm 3'teki Es'in "ani, kısa, çok güçlü" karakterini bizzat tanımak ve sky-wave (zayıf, kademeli) ile Es (güçlü, ani) arasındaki farkı kulağınla ayırt etmek. Not: Bu tamamen pasif dinlemedir; kayıt/yayma konusunda yerel mevzuatı teyit et.

### C) Bir Starlink (LEO) geçişini gözlemlemek (optik ve/veya spektrum)

İki katmanlı, tamamen pasif bir gözlem.

Optik: Açık bir takip aracıyla bulunduğun konum için bir Starlink (özellikle yeni fırlatma sonrası "tren") veya parlak LEO geçişinin zamanını hesapla; belirtilen zamanda çıplak gözle veya basit dürbünle gökyüzünü izle. Ardışık parlak noktaların düzgün bir hat halinde geçişi (LEO treni) etkileyici ve tartışmasız yasaldır.

Spektrum (varsa uygun donanım): Uygun bir bandda bir SDR ile, bir LEO geçişi sırasında downlink sinyalinin Doppler-taranan imzasını (taşıyıcının geçiş boyunca yükselip alçalan S-eğrisi) waterfall'da gözle. Bu, Bölüm 11'deki Doppler/handover imzasının canlı halidir.

| Gözlem | Optik | Spektrum |
|---|---|---|
| Ne ararsın | Hat halinde geçen parlak noktalar | Doppler S-eğrisi, periyodik yapı |
| Yasal durum | Tamamen serbest (gökyüzü gözlemi) | Pasif RX; band mevzuatını teyit et |
| Bağlandığı kavram | LEO yörünge geometrisi (Bölüm 11) | Doppler/handover (Bölüm 11) |

Amaç: "Gelecek spektrum"un (mega-konstelasyon) hem optik hem RF imzasını somut deneyimlemek ve LEO'yu sabit yayıncılardan ayıran imzayı (hızlı hareket → Doppler + optik tren) tanımak.

### D) VLF doğal radyo: whistler ve sferic dinlemek (pasif, doğal kaynak)

En düşük frekansın doğal sesini dinle. Bir VLF alıcısı kur — en basit haliyle uzun bir tel anten + ses-bandı yükselteci + bilgisayar ses kartı (VLF frekansları doğrudan ses bandına düştüğü için özel RF donanımı gerekmeyebilir; basit "doğal radyo alıcısı" şemaları açık kaynaklarda mevcut, kaynaktan teyit et). Elektriksel gürültüden olabildiğince uzakta (kırsal, şebeke gürültüsünden uzak) dinle.

| Ararsın | Ses karakteri | Kaynak (Bölüm 9) |
|---|---|---|
| Sferic | Kısa "tık/çıtırtı" | Uzak yıldırım darbeleri |
| Whistler | Yüksekten alçağa kayan ıslık (~1-3 sn) | Manyetosferde yol alan yıldırım enerjisi |
| Dawn chorus | Cıvıltı benzeri (varsa) | Manyetosferik plazma dalgaları |

Amaç: Bölüm 9'daki doğal radyoyu bizzat duymak; "verici olmadan da spektrum doludur" gerçeğini ve en basit alıcıyla bile gözlemlenebilen atmosferik fiziği deneyimlemek. En iyi sonuç, elektriksel olarak sessiz bir ortamda ve şebeke gürültüsünden (50/60 Hz harmonikleri) uzakta alınır.

### E) Bilinmeyen sinyal metodolojisini bir gerçek gözleme uygulamak (sentez)

Yukarıdaki gözlemlerden birinde (ya da herhangi bir bantta) tanımadığın bir sinyalle karşılaştığında, Bölüm 19'daki akışı bir gözlem kâğıdında baştan sona uygula:

1. Nerede (frekans/band — Bölüm 8'e göre ne beklenir)?
2. Ne zaman/nasıl (sürekli mi, burst mu, Doppler mı — hangi yayılım modu olabilir)?
3. Ne biçimde (modülasyon: neyin sabit/değişken — Bölüm 7)?
4. Hangi ritimde (periyodik yapı var mı)?
5. Nereden (yönü kestirebiliyor musun — Bölüm 9)?
6. Eşleşiyor mu (açık veritabanı/sigidwiki ile karşılaştır)?
7. Hipotezin ne ve onu nasıl teyit edebilirsin?

Amaç: Serinin tüm araçlarını tek bir gerçek (yasal, pasif) gözlemde birleştirmek ve "bu nasıl buraya geldi?" (yayılım) sorusunu "bu ne?" (içerik/tip) sorusundan önce sormayı içselleştirmek.

---

<a id="21"></a>
## 21. Hızlı Referans ve Diğer Bölümler

### Kavram kartı

| Kavram | Bir cümlelik öz |
|---|---|
| Sky-wave / NVIS | İyonosferin uzun-mesafe yansıması / dik-geliş bölgesel doldurma |
| Gri-çizgi | Terminator boyunca düşük-yutmalı kısa altın yayılım penceresi |
| Güneş döngüsü | 11 yıllık aktivite; yüksekte HF üst bandları açılır, patlamalar HF'i yutar |
| Ducting / Sporadic-E | Troposferik kanal hapsi / E-katmanı ani güçlü VHF açılımı |
| Meteor scatter | Meteor izinden saniyelik VHF burst; MSK144 ile çözülür |
| EME (moonbounce) | Ay'dan yansıma; ~240-260 dB kayıp; kodlama kazancıyla çözülür |
| VLF/ELF | Spektrum dibi; cilt derinliği yüzünden denizaltıya işler; devasa anten, çok yavaş |
| eLORAN | GNSS-bağımsız LF PNT yedeği (karıştırmaya dayanıklı) |
| Doğal radyo | Whistler/sferic/Schumann — vericisiz, en basit alıcıyla dinlenir |
| Zayıf-sinyal (WSJT) | Hızı düşürüp ağır kodlayarak gürültü altında çözme; WSPR yayılım haritası |
| Mega-konstelasyon | Binlerce LEO uydu; faz dizili anten, Doppler + handover imzası |
| Direct-to-device | Sıradan telefona uydu bağı; zayıf-sinyal felsefesinin ticari hali |
| Bilişsel radyo / DSA | Spektrumu algılayıp boşlukları paylaşan radyo (TVWS, CBRS) |
| OTFS / FBMC | Gecikme-Doppler / filtreli çok-taşıyıcı; OFDM ardılı adaylar (araştırma) |
| mmWave / THz | Üst frekans = çok bant genişliği ama çok kısa menzil; THz araştırmada |
| Li-Fi / optik | Işıkla veri; duvar geçmez (oda-içi gizli); FSO/uzay-lazer/derin-uzay optik |
| PQC | Kuantum-dirençli şifreleme; "şimdi topla sonra çöz" tehdidine karşı |
| QKD / kuantum algılama | Fiziğe dayalı anahtar / aşırı hassas alan ölçümü (çoğu araştırma aşamasında) |
| AI-tanımlı radyo | Çevresini algılayıp parametrelerini kendi seçen/öğrenen radyo |
| Bilinmeyen sinyal akışı | Nerede→ne zaman→ne biçimde→hangi ritimde→nereden→eşleştir→teyit |

### Ezber sezgiler

- Egzotik modlar "her zaman açık" değildir; iyonosfer güneşe, ducting havaya, meteor scatter gökyüzüne, EME Ay'a bağlıdır.
- Frekans yalnızca menzil değildir; aynı anda nüfuz (cilt derinliği), anten boyutu ve kapasite arasında bir takastır (ELF ekstrem örnek).
- Zayıf-sinyal sanatının özü: hızı yavaşlatıp mesajı küçülterek dayanıklılık satın almak (Shannon'ın pratiği).
- Gürültünün altındaki sinyal de bilgi taşır; yeterli kodlama kazancı onu çıkarır (EME'den uyduya, WSJT'den hücresele).
- LEO'yu sabit yayıncıdan ayıran imza: hızlı Doppler tarama + periyodik handover.
- Geleceğin spektrumu statik tahsis değil, algılanan/müzakere edilen paylaşımdır (bilişsel radyo, CBRS).
- Spektrumun üstüne tırmanmak (mmWave/THz) bant genişliği kazandırır, fiziği acımasızlaştırır; gelecek katmanlıdır.
- "Şimdi topla, sonra çöz" tehdidi, uzun-ömürlü kablosuz trafiğin bugün PQC'ye geçişini gerektirir.
- Bilinmeyen sinyalde önce "bu nasıl buraya geldi?" (yayılım), sonra "bu ne?" (içerik) sorulur.
- Ayırt edici güç tek ölçümde değil, kanıtların birlikte daralttığı örüntüdedir; hipotezi daima bağımsız teyit et.

### Spekülatif/olgunlaşmamış konuların dürüstlük listesi

Bu bölümün bazı konuları yerleşik fizik değil, araştırma/erken-konuşlanma aşamasındadır. Karıştırmamak için açıkça işaretleniyor:

| Konu | Durum |
|---|---|
| İyonosfer, EME yol kaybı, Doppler, cilt derinliği | Yerleşik fizik (kesin) |
| Sporadic-E oluşum mekanizması | Gözlem net, kesin mekanizma araştırma aşamasında |
| eLORAN konuşlanması | Ülkeye göre pilot/kısmi; güncel durum teyit edilmeli |
| Direct-to-device uydu | Kısmen konuşlanmış, olgunlaşmamış; teyit edilmeli |
| OTFS / FBMC | Araştırma/standart-öncesi; kitlesel konuşlanma yok |
| THz haberleşme | Büyük ölçüde laboratuvar/araştırma |
| QKD pratik konuşlanması | Saha gösterimi var; geniş pratiklik tartışmalı |
| Kuantum RF algılama (Rydberg) | Araştırma aşamasında |
| Pulsar tabanlı navigasyon | Araştırma/kavram aşamasında |
| Derin-uzay optik haberleşme | Test/gösterim aşamasında |

### Serinin diğer bölümleri (çapraz referans)

- SIGINT_01 — RF Fiziği, Spektrum ve Modülasyon: FSPL, dB/dBm, Shannon kapasitesi, IQ. (Bu bölümdeki EME yol kaybı, zayıf-sinyal kodlama kazancı ve ELF kapasite-takasının fiziksel temeli oradadır.)
- SIGINT_02 — SDR Donanımları ve Yazılım Ekosistemi: RTL-SDR/HackRF/USRP, GNU Radio. (Alıştırmalardaki WSPR/VLF/LEO gözlemlerini burada üretirsin.)
- SIGINT_03 — Anten Teorisi, Kapsama ve Yön Bulma (DF): dizi antenler, DF temelleri. (EME'nin büyük yönlü dizileri ve faz dizili anten ilkesinin temeli orada.)
- SIGINT_05 — Hücresel/WiFi/BT ve IoT Spektrumu: protokol imzaları. (Direct-to-device ve uydu-IoT'nin karasal komşusu orada.)
- SIGINT_06 — TEMPEST, RF Sızıntısı ve Savunma: istem dışı yayın, OPSEC. (Ducting'in "uzaktan duyulma" riski ve Li-Fi'nin "duvar geçmeyen" gizliliği oranın kardeşidir.)
- SIGINT_07 — Disiplinler ve Sinyal Ayıklama: AMC, deinterleaving, kütüphane eşleştirme, cyclostationarity. (Bilinmeyen-sinyal metodolojisinin ayıklama/sınıflandırma çekirdeği orada.)
- SIGINT_08 — Frekans Tahsisi ve Bant Planı: "bu bandda ne beklenir" referansı. (Bilinmeyen sinyal akışının ilk adımı oraya dayanır.)
- SIGINT_09 — Yer Tespiti, Yön Bulma ve Takip: TDOA, AOA, çoklu-alıcı konumlama. (eLORAN'ın TDOA mantığı ve metodolojinin "nereden" adımı orada.)
- SIGINT_10 — GNSS/GPS Sistemleri: uydu zamanlama, zayıf sinyal, jamming/spoofing. (eLORAN'ın yedeklediği birincil PNT ve "GNSS kırılganlığı" oradadır.)
- SIGINT_11 — Uydu Haberleşmesi: GEO/LEO temelleri, link bütçesi. (Mega-konstelasyon ve Doppler/handover ayrıntısının temeli oradadır.)
- SIGINT_13 — RF Tehdit ve Karşı Önlemler: jamming, spoofing, dayanıklılık. (eLORAN yedekleme ve AI-radyonun karıştırmadan kaçma gerekçesi orada.)
- SIGINT_18 — Massive MIMO ve Beamforming (planlanan): uzamsal işleme derinlemesine. (Faz dizili anten ve hüzme yönlendirmenin tam matematiği orada.)
- SIGINT_19 — Yapay Zekâ-Tanımlı Radyo (planlanan): öğrenen radyo derinlemesine. (Bu bölümdeki AI-radyo özetinin tam işlenişi orada.)
- SIGINT_20 — İleri Hücresel 4G/5G Güvenlik: modern hücresel mimari. (OFDM, beamforming ve yeni dalga formlarının hücresel bağlamı orada.)

> Kapanış: Sıradan spektrum, çizgi-görüş ve sabit tahsisten ibaret görünür; ama kenarlarında atmosfer bir aynaya, bir meteorun izi geçici bir yansıtıcıya, Ay uzak bir hedefe dönüşür — ve dibinde, yıldırımın doğal radyosu hiç susmaz. Geleceğe baktığımızda ise spektrum giderek daha kalabalık (mega-konstelasyonlar), daha akıllı (bilişsel/AI radyo), daha yüksek (mmWave/THz) ve daha güvenli-ama-yeni-tehditli (post-kuantum) bir ekosisteme dönüşüyor. Bu bölümün tek bir sezgisi varsa o da şudur: spektrum statik bir tablo değil, fiziği güneşe ve havaya, geleceği insanın mühendisliğine bağlı, sürekli yaşayan bir ortamdır. Bir araştırmacının ustalığı, tanımadığı bir sinyalle karşılaştığında paniğe değil metoda — nerede, ne zaman, nasıl, hangi ritimde, nereden, neye benziyor, nasıl teyit ederim — başvurabilmesidir. Bandını, ülkeni ve sürümünü teyit et; bu kitap anlama, savunma ve merak içindir.
>
> Bu doküman Kanije Kalesi güvenlik/teknik rehberleri koleksiyonunun SIGINT serisinin 22. bölümüdür. İlgili: SIGINT_01–20, `VERACRYPT_USTALIK_REHBERI.md`, `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`.
