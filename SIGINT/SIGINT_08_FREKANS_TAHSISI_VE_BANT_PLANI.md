# SIGINT El Kitabı — Bölüm 8: Frekans Tahsisi ve Bant Planı

## Kim, Nerede, Hangi Frekansı Kullanır

Bu bölüm, radyo spektrumunun nasıl bölüştürüldüğünü ve hangi kullanıcı grubunun (havacılık, denizcilik, askeri, kamu güvenliği, amatör, ticari, uydu) hangi frekans aralığında bulunduğunu, açık ve kamuya açık tahsis bilgisi düzeyinde ele alır. Amaç, bir frekans değeri ya da bir sinyal duyulduğunda "bu kime ait, ne işe yarar" sorusunu sistematik biçimde yanıtlayabilmektir.

Bu doküman bir **bant tahsisi farkındalığı** kaynağıdır. İçindeki tablolar, hangi bandın hangi hizmete *ayrıldığını* gösterir; spesifik operasyonel kanalları, şifreli içeriği ya da gizli kullanım listelerini değil. Askeri bölümlerde verilen aralıklar, uluslararası ve ulusal tahsis tablolarında açıkça yayımlanmış **bant sınırlarıdır**; o bantlar içindeki belirli operasyonel frekanslar, çağrı planları ve şifreleme açık kaynak değildir ve burada yer almaz.

Uyarı (yasal): Çoğu açık bandın (havacılık ATC, denizcilik VHF, amatör, yayın) **dinlenmesi** birçok ülkede serbesttir. Buna karşılık bazı bantların içeriğini çözmek (kamu telsizinin şifreli trafiği, ticari/kurumsal sayısal sistemler, askeri şifreli kanallar) ya da dinlediğini üçüncü kişilere aktarmak/kayıt altına alıp yaymak yasak olabilir. Yayın yapmak (TX) ise neredeyse her bantta lisans/yetki ister. Kendi ülkenin telekomünikasyon düzenleyicisinin (Türkiye'de BTK) kurallarını teyit etmeden hareket etme. Bu bölüm eğitim, spektrum farkındalığı ve savunma amaçlıdır.

---

## İçindekiler

1. [Spektrum Yönetiminin Mantığı](#1)
   - 1.1 [ITU ve Üç Bölge (Region 1/2/3)](#1-1)
   - 1.2 [Ulusal Regülatörler — BTK, FCC ve Diğerleri](#1-2)
   - 1.3 [Tahsis, Tahsisat, Tahis — Terim Ayrımı](#1-3)
   - 1.4 [Birincil ve İkincil Kullanıcı](#1-4)
   - 1.5 [Bant Adlandırma (VLF–EHF) ve Yayılım Hatırlatması](#1-5)
2. [Havacılık — Pilot ve Hava Trafik (Sivil)](#2)
3. [Askeri Havacılık — UHF Hava Bandı](#3)
4. [Denizcilik — Gemi ve Kıyı](#4)
5. [Kara Askeri Taktik Haberleşme](#5)
6. [Kamu Güvenliği — Polis, İtfaiye, Ambulans](#6)
7. [Amatör Radyo — Tam Bant Planı](#7)
8. [ISM / SRD — Lisanssız Kısa Menzil](#8)
9. [Uydu, Seyrüsefer ve Yayın](#9)
10. [GSM / LTE / 5G — Hücresel Bantlar](#10)
11. [Bir Frekansın Sahibini Belirleme Metodolojisi](#11)
12. [Tahsis Veritabanları ve Kaynaklar](#12)
13. [Alıştırmalar (Yasal)](#13)
14. [Çapraz Referans ve Sonraki Bölümler](#14)

---

<a id="1"></a>
## 1. Spektrum Yönetiminin Mantığı

Radyo spektrumu, sınırlı ve paylaşılan bir doğal kaynaktır. Aynı frekansta aynı bölgede iki güçlü vericinin koordinasyonsuz çalışması, her ikisini de kullanılamaz hale getirir (girişim). Bu nedenle spektrum, idari olarak hizmetlere ve kullanıcılara bölüştürülür. Bu bölüştürmenin üç katmanı vardır: uluslararası (ITU), bölgesel ve ulusal (her ülkenin kendi regülatörü). Bir frekansın "kime ait olduğunu" anlamanın temeli, bu katmanlı tahsis mantığını bilmektir.

<a id="1-1"></a>
### 1.1 ITU ve Üç Bölge (Region 1/2/3)

Uluslararası Telekomünikasyon Birliği (ITU — International Telecommunication Union), Birleşmiş Milletler'in radyo spektrumu ve uydu yörüngelerini düzenleyen ihtisas kuruluşudur. ITU'nun temel ürünü, periyodik Dünya Radyokomünikasyon Konferansları'nda (WRC) güncellenen **Radyo Tüzüğü** (Radio Regulations) ve onun içindeki **Frekans Tahsis Tablosu**dur (Table of Frequency Allocations). Bu tablo, hangi frekans aralığının hangi *radyokomünikasyon hizmetine* (sabit, mobil, havacılık mobil, deniz mobil, radyo seyrüsefer, radyo astronomi, amatör, yayın, vb.) ayrıldığını belirler.

ITU, yayılım koşulları ve tarihsel kullanım farkları nedeniyle dünyayı üç bölgeye ayırır:

| Bölge | Kapsadığı coğrafya | Karakteristik |
|---|---|---|
| **Region 1** | Avrupa, Afrika, Orta Doğu (Basra Körfezi'nin batısı dahil), eski Sovyet coğrafyası, Moğolistan | Türkiye buradadır. Amatör ve yayın bantlarında Region 2/3'ten farklar var. |
| **Region 2** | Kuzey ve Güney Amerika, Grönland, doğu Pasifik adaları | ABD/Kanada düzeni; 60 m amatör, bazı 40/80 m segment farkları. |
| **Region 3** | Asya'nın büyük kısmı, Avustralya, Yeni Zelanda, batı Pasifik, İran | Region 1'e birçok noktada benzer ama bağımsız farklar içerir. |

Pratikte: Bir amatör bant ya da yayın frekansı için "şu segment SSB, şu segment CW" gibi sınırlar **bölgeye göre değişebilir.** Türkiye Region 1 kurallarına tabidir; bir ABD kaynağındaki (Region 2) segment sınırlarını Türkiye'ye doğrudan uygulamak hataya yol açar. Havacılık ve denizcilik gibi küresel mobil hizmetlerde ise frekanslar büyük ölçüde dünya genelinde aynıdır — uçak ve gemi sınır tanımadan hareket ettiği için bu hizmetler ITU düzeyinde küresel standarda bağlanmıştır.

<a id="1-2"></a>
### 1.2 Ulusal Regülatörler — BTK, FCC ve Diğerleri

ITU tablosu bir üst çerçevedir; o çerçeve içinde her ülke kendi **ulusal frekans tahsis tablosunu** yayımlar. Ülke, ITU'nun hizmet ayrımına genel olarak uyar ama bant içindeki kanal planını, güç sınırlarını, lisans rejimini ve belirli alt segmentlerin kullanımını kendi belirler.

| Ülke / bölge | Düzenleyici | Yayımladığı temel belge |
|---|---|---|
| **Türkiye** | BTK (Bilgi Teknolojileri ve İletişim Kurumu) | Milli Frekans Planı; ilgili tebliğ ve yönetmelikler |
| **ABD** | FCC (sivil) ve NTIA (federal/devlet) | FCC Online Table of Frequency Allocations; CFR Title 47 |
| **Birleşik Krallık** | Ofcom | UK Frequency Allocation Table |
| **Almanya** | Bundesnetzagentur | Frequenzplan |
| **Genel Avrupa uyumu** | CEPT / ECC | ECA Table (European Common Allocation) |

Not (ABD'ye özgü ikilik): ABD'de spektrum, sivil/ticari kullanım (FCC) ile federal hükümet kullanımı (NTIA) arasında bölünmüştür. Bu yüzden ABD kaynaklarında bir bandın "government" ya da "non-government" olduğu ayrıca belirtilir. Türkiye'de tek merkezi sivil otorite BTK'dır; askeri tahsisler ise ayrıca koordine edilir ve kamuya kapalıdır.

Kritik kural: Bu bölümdeki ülke/bölge bağımlı her değer (özellikle kamu telsizi standartları, GSM/LTE operatör bant dağılımı, bazı amatör segment sınırları) **kesin karar için ilgili ulusal tahsis tablosundan teyit edilmelidir.** Küresel mobil hizmetler (havacılık 118–137 MHz, denizcilik 156–162 MHz, GPS 1575.42 MHz gibi) ise ülkeye göre değişmez.

<a id="1-3"></a>
### 1.3 Tahsis, Tahsisat, Tahis — Terim Ayrımı

ITU terminolojisinde üç kademeli bir hak verme zinciri vardır; Türkçe regülasyon dilinde de karşılıkları kullanılır:

- **Tahsis (allocation):** Bir frekans bandının, tablo düzeyinde bir ya da birden çok **radyokomünikasyon hizmetine** ayrılmasıdır. Örnek: "108–137 MHz havacılık mobil hizmetine tahsis edilmiştir." En üst, en soyut kademe budur.
- **Tahsisat / paylaştırma (allotment):** Tahsis edilmiş bir bandın, belirli bir coğrafi alan ya da ülke grubu için belirli kullanıma planlanmasıdır. Örnek: Bir yayın bandının ülkeler arasında kanal planına bölünmesi.
- **Tahis / atama (assignment):** Belirli bir istasyona, belirli bir frekansın (ya da kanalın) fiilen verilmesidir — lisanslama. Örnek: Bir havaalanı kulesine 121.300 MHz'in atanması.

Pratikte: Bu bölümdeki tablolar büyük ölçüde **tahsis** ve kısmen **tahsisat** düzeyindedir (bant → hizmet → tipik kanal planı). Belirli bir istasyonun hangi frekansı kullandığı (atama) ise tahsis tablolarında değil, lisans veritabanlarında ya da havaalanı/liman yayın bilgilerinde bulunur. Askeri **atamalar** kamuya açık değildir; bu bölüm yalnızca askeri **tahsis** (bant) düzeyinde kalır.

<a id="1-4"></a>
### 1.4 Birincil ve İkincil Kullanıcı

Çoğu bant tek bir hizmete değil, birden çok hizmete aynı anda tahsis edilir. Bu durumda hizmetler arasında bir öncelik hiyerarşisi konur:

- **Birincil (primary) kullanıcı:** Bandın korunan asıl sahibidir. Girişimden korunma hakkına sahiptir ve diğer hizmetlere girişim yapma konusunda öncelikli sayılır. Tahsis tablolarında genellikle **BÜYÜK HARFLE** yazılır.
- **İkincil (secondary) kullanıcı:** Bandı kullanabilir ama birincil kullanıcıya girişim yapmamak zorundadır ve birincil kullanıcıdan korunma talep edemez. Tahsis tablolarında küçük harfle yazılır.

Örnek: Birçok ülkede bazı amatör bantlar (özellikle bazı mikrodalga segmentler) amatör hizmete **ikincil** statüde tahsis edilir; aynı bandın birincil sahibi radyolokasyon ya da sabit hizmet olabilir. Bu, amatörün o bandı kullanabileceği ama radar ya da birincil hizmet aktif olduğunda ona katlanmak ve girişim yapmamak zorunda olduğu anlamına gelir.

Not: Birincil/ikincil ayrımı, "bir frekansta beklenmedik bir sinyal neden var?" sorusunun yanıtı olabilir. Bir amatör segmentte zaman zaman duyulan güçlü, amatör olmayan sinyal, çoğu zaman o bandın birincil sahibinin (örneğin bir radyolokasyon sistemi) meşru kullanımıdır.

<a id="1-5"></a>
### 1.5 Bant Adlandırma (VLF–EHF) ve Yayılım Hatırlatması

Frekans aralıkları, on katlı dilimler halinde standart adlarla anılır. Bir frekansın hangi dilimde olduğu, onun **nasıl yayıldığını** (dolayısıyla kim tarafından, hangi menzil için kullanıldığını) büyük ölçüde belirler.

| Kısaltma | Açılım | Frekans | Dalga boyu | Tipik yayılım ve kullanım |
|---|---|---|---|---|
| **VLF** | Very Low Frequency | 3–30 kHz | 100–10 km | Yer dalgası, denizaltı haberleşmesi, zaman/standart işaret |
| **LF** | Low Frequency | 30–300 kHz | 10–1 km | Yer dalgası, uzun dalga (LW) yayın, seyrüsefer (eski) |
| **MF** | Medium Frequency | 300 kHz – 3 MHz | 1000–100 m | Orta dalga (MW) AM yayın, deniz MF, NAVTEX |
| **HF** | High Frequency | 3–30 MHz | 100–10 m | İyonosferden yansıma (skywave), kıtalararası; kısa dalga, havacılık/deniz HF, amatör HF, taktik HF |
| **VHF** | Very High Frequency | 30–300 MHz | 10–1 m | Görüş hattı (line-of-sight); FM yayın, havacılık VHF, deniz VHF, amatör 2 m, kara taktik alt-VHF |
| **UHF** | Ultra High Frequency | 300 MHz – 3 GHz | 1 m – 10 cm | Görüş hattı; TV, GSM/LTE, askeri hava UHF, GPS, amatör 70 cm, ISM 433/868/915, WiFi 2.4 GHz |
| **SHF** | Super High Frequency | 3–30 GHz | 10–1 cm | Mikrodalga link, uydu, radar, WiFi 5 GHz |
| **EHF** | Extremely High Frequency | 30–300 GHz | 10–1 mm | Milimetre dalga, gelişmiş radar, 5G mmWave, uydu üst bantları |

Pratikte: HF ve altı, iyonosfer/yer dalgası sayesinde ufuk ötesine gider — bu yüzden kıtalararası (okyanusüstü havacılık, denizaşırı deniz, uzun menzilli taktik, kısa dalga yayın) hizmetler buradadır. VHF ve üstü esasen görüş hattıyla sınırlıdır — bu yüzden yerel ve hava/deniz yakın menzil hizmetleri (havaalanı kulesi, liman, şehir telsizi, cep telefonu) buradadır. Bir frekansın hangi dilimde olduğunu bilmek, kullanıcı grubunu daraltmanın ilk adımıdır. Yayılım mekanizmalarının ayrıntısı için bkz. Bölüm 1.

---

<a id="2"></a>
## 2. Havacılık — Pilot ve Hava Trafik (Sivil)

Sivil havacılık, tüm spektrum kullanıcıları içinde en standartlaşmış ve en çok küresel uyumlu olanıdır; bir uçak dünyanın herhangi bir yerinde aynı bantları aynı biçimde kullanır. Havacılık, sesli hava trafik kontrolü (ATC), seyrüsefer yardımcıları ve gözetim (surveillance) olmak üzere üç ana işlev için spektrum kullanır. Sesli ATC ve airband, dinlemesi birçok ülkede serbest olan, başlangıç için en öğretici bantlardandır.

Not (modülasyon): Havacılık sesli haberleşmesi, neredeyse istisnasız **genlik modülasyonu (AM)** kullanır; FM değil. Bunun klasik gerekçesi, AM'in "yakalama etkisi" (capture effect) göstermemesidir: iki uçak aynı anda konuşursa, FM'de güçlü olan zayıfı tamamen bastırırken, AM'de her iki sinyal de duyulur (heterodyne uğultusuyla birlikte) ve kulenin en azından çakışmayı fark etmesi sağlanır. Bu, can güvenliği açısından tercih edilmiş bilinçli bir mühendislik kararıdır.

### 2.1 Sivil havacılık bant tablosu

| Frekans aralığı | Kullanım | Modülasyon | Not |
|---|---|---|---|
| **108.000 – 111.975 MHz** | ILS Localizer + VOR (seyrüsefer) | AM / özel taşıyıcılar | Çift kullanım: tek MHz onda biri VOR, çift olanlar ILS yatay kılavuz |
| **111.975 – 117.975 MHz** | VOR (VHF Omnidirectional Range) | AM taşıyıcı + alt taşıyıcı | Seyrüsefer yön referansı; yer istasyonu yayını |
| **108 – 117.975 MHz (genel)** | Seyrüsefer yardımcıları (NAV) | — | Bu blok haberleşme değil seyrüsefer; pilot konuşmaz, alıcı çözer |
| **118.000 – 136.975 MHz** | ATC sesli haberleşme (COM/airband) | **AM (A3E)** | Kule, yaklaşma, saha, ACC; kanal aralığı 25 kHz ya da **8.33 kHz** |
| **121.500 MHz** | Acil durum / imdat (VHF guard) | AM | Uluslararası havacılık imdat frekansı; sürekli izlenir |
| **123.100 MHz** | Arama-kurtarma (SAR) sahne | AM | Olay yeri koordinasyonu |
| **121.700 – 121.900 MHz (tipik)** | Yer kontrol (ground) | AM | Taksi/yer trafiği; havaalanına göre değişir |
| **131.x MHz (ör. 131.550, 131.725)** | ACARS (veri linki) | AM üzerinden MSK/veri | Uçak–yer kısa metin/telemetri; çözülebilir veri |
| **136.900 – 137.000 MHz civarı** | VDL Mode 2 (sayısal veri linki) | D8PSK | ACARS'ın sayısal evrimi |

Not (8.33 kHz kanal aralığı): Avrupa'da (Türkiye dahil) hava trafiğinin yoğunluğu, klasik 25 kHz kanal aralığını yetersiz bıraktığı için birçok ATC kanalı **8.33 kHz** aralığa geçmiştir; bu, aynı bant içinde üç kat daha fazla kanal demektir. Bir alıcıda airband dinlerken hem 25 kHz hem 8.33 kHz adımlamayı desteklemek gerekir. Hangi sahanın hangi adımlamayı kullandığı ulusal havacılık yayınlarından (AIP) öğrenilir.

### 2.2 HF havacılık (okyanusüstü / uzun menzil)

VHF airband görüş hattıyla sınırlı olduğundan (yaklaşık ufka kadar), okyanus üstü ve uzak bölge uçuşlarında kıtalararası haberleşme için HF kullanılır. HF, iyonosferden yansıyarak binlerce kilometre gider.

| Frekans aralığı | Kullanım | Modülasyon | Not |
|---|---|---|---|
| **2.850 – 22.000 MHz (havacılık HF tahsisleri)** | Okyanusüstü/uzun menzil ATC | **SSB (USB)** | Belirli aile frekanslarına (MWARA, family) bölünmüş; gündüz/gece farklı bant |
| Örnek aile bantları (3, 5, 6, 8, 10, 13, 17 MHz dilimleri) | Bölgesel okyanus rotaları | USB | Hangi ailenin hangi rotaya hizmet ettiği AIP/HF plandan teyit edilir |
| **HFDL (HF Data Link)** | Sayısal uçak–yer veri | PSK | HF üzerinden konum/telemetri; çözülebilir |

Not: Havacılık HF, gün/gece ve mevsime göre yayılımı değiştiği için "aynı uçak gündüz 8 MHz, gece 5 MHz ailesini kullanır" gibi bir esneklikle çalışır. Bu, HF yayılımının (Bölüm 1) doğrudan operasyonel sonucudur.

### 2.3 Gözetim (surveillance) ve transponder

Hava trafik kontrolünün uçağı "görmesi", radar ve transponder sistemlerine dayanır. Bunlar sesli değil, sayısal darbe sistemleridir.

| Frekans | Kullanım | Modülasyon | Not |
|---|---|---|---|
| **1030 MHz** | İkincil gözetim radarı (SSR) sorgu — yerden uçağa | Darbe (pulse) | Kule sorar |
| **1090 MHz** | Transponder yanıtı + **ADS-B** (uçağın kendi konum yayını) | Darbe konum modülasyonu (PPM) | En popüler dinleme hedefi; `dump1090` ile çözülür |
| **978 MHz (UAT)** | Universal Access Transceiver (genel havacılık, ABD) | — | 1090 ADS-B'ye alternatif, düşük irtifa/genel havacılık; esasen ABD |

Pratikte: 1090 MHz ADS-B, ucuz bir alıcıyla çözülebilen, uçağın kimliğini, konumunu, irtifasını ve hızını açıkça yayınlayan bir sistemdir; bu yüzden havacılık SIGINT'inde en sık başlanan noktadır. Cihaz seçimi için bkz. Bölüm 2 (SDR cihazları); kod çözme için bkz. Bölüm 5.

---

<a id="3"></a>
## 3. Askeri Havacılık — UHF Hava Bandı

Askeri havacılık, sivil VHF airband'den ayrı, kendine ait bir bant kullanır: **225–400 MHz UHF askeri hava bandı.** Bu, NATO ve birçok ülke için fiili standarttır. Bu bölümde verilen tek değer, bu bandın **dış sınırlarıdır**; bant içindeki belirli operasyonel frekanslar, çağrı işaretleri, atlama planları ve şifreli içerik açık kaynak değildir ve burada yer almaz.

| Frekans aralığı | Kullanım | Modülasyon | Not |
|---|---|---|---|
| **225 – 400 MHz** | Askeri havacılık sesli (UHF airband) | Esasen **AM**; sayısal/şifreli modlar da kullanılır | Bant tahsisi açıktır; içindeki belirli kanallar/atamalar açık değildir |
| **243.000 MHz** | Askeri imdat (UHF guard) | AM | Sivil 121.5 MHz'in askeri karşılığı; tam iki katı olması tasarımdır |

Not (modülasyon mantığı): UHF askeri hava bandı da, sivil airband gibi tarihsel olarak AM tabanlıdır — aynı capture-effect gerekçesiyle. Modern platformlar bunun üstüne sayısal sesli ve şifreli modlar ekler.

### Frekans atlamalı haberleşme kavramı (Have Quick)

Sabit bir frekansta süren askeri haberleşme, hem dinlemeye hem karıştırmaya (jamming) açıktır. Bunu azaltmak için askeri UHF hava haberleşmesinde **frekans atlama** (frequency hopping) yaklaşımı geliştirilmiştir; bu ailenin yaygın bilinen adı **Have Quick**'tir. Temel fikir şudur: verici ve alıcı, ortak ve gizli bir desene (ve hassas bir ortak zaman referansına) göre, saniyede birçok kez farklı frekanslara birlikte atlar. Dışarıdan bakan biri için sinyal, bant boyunca rastgele dağılmış kısa parçalar gibi görünür; deseni ve zaman senkronizasyonunu bilmeyen bir dinleyici, anlamlı ve kesintisiz bir kanal yakalayamaz.

Bu, **kavramsal** bir açıklamadır. Atlama deseni, anahtarlar, zaman senkronizasyon yöntemi ve içerik açık değildir; burada anlatılan yalnızca "neden sabit frekansta dinleyince bir şey bulamazsın" sorusunun teknik gerekçesidir. Frekans atlamanın daha geniş anlatımı için bkz. yayılı spektrum (spread spectrum) konusu, Bölüm 1.

Uyarı: Askeri hava bandı tahsisini bilmek (bu bandın 225–400 MHz olduğunu bilmek) açık bilgidir ve bir tarayıcıda bu aralıkta aktivite görmek mümkündür; ancak bu bandın içeriğini çözmeye, şifresini açmaya ya da kayıtlarını yaymaya çalışmak birçok ülkede yasaktır. Bu bölüm yalnızca bandın *varlığı* ve *sınırları* hakkında farkındalık sağlar.

---

<a id="4"></a>
## 4. Denizcilik — Gemi ve Kıyı

Denizcilik haberleşmesi, havacılık gibi küresel olarak standartlaşmıştır ve büyük ölçüde açıktır; deniz VHF dinlemek çoğu yerde serbesttir. Denizcilik, yakın menzil (VHF), orta/uzun menzil (MF/HF), sayısal çağrı (DSC), otomatik kimlik (AIS) ve acil/güvenlik (EPIRB, NAVTEX) için ayrı sistemler kullanır.

Not (modülasyon): Havacılığın aksine, deniz VHF sesli haberleşmesi **FM (daha doğrusu dar bant FM, NBFM)** kullanır. Bu, havacılık (AM) ile denizcilik (FM) arasındaki en akılda kalıcı ayrımlardandır ve bir sinyalin türünü teyit ederken işe yarar.

### 4.1 Deniz VHF bandı (156–162 MHz)

Deniz VHF, kanal numaralarıyla anılan, 25 kHz aralıklı sabit bir kanal planı kullanır. Kanal numaraları uluslararasıdır.

| Frekans / kanal | Kullanım | Modülasyon | Not |
|---|---|---|---|
| **156.800 MHz (Kanal 16)** | Uluslararası imdat, güvenlik ve çağrı | FM (NBFM) | Sürekli dinlenir; ilk temas ve acil çağrı kanalı |
| **156.525 MHz (Kanal 70)** | DSC (Digital Selective Calling) | FM üzerinden sayısal | Sesli değil; sayısal acil/çağrı sinyali, GMDSS'in parçası |
| **156.300 MHz (Kanal 06)** | Gemiden gemiye güvenlik | FM | Köprüden köprüye koordinasyon |
| **156.650 MHz (Kanal 13)** | Köprü-köprü seyir güvenliği | FM | Çatışmayı önleme, manevra koordinasyonu |
| **157.x / 161.x MHz (dubleks kanallar)** | Kıyı telsiz, liman operasyonu, public correspondence | FM | Gemi ve kıyı farklı frekanslarda (dubleks) konuşur |
| **156 – 162 MHz (genel deniz VHF)** | Gemi/kıyı sesli + sayısal | FM | Tam kanal listesi ITU Appendix 18'den teyit edilir |

Not: Kanal 16 (156.800 MHz), denizciliğin "ortak buluşma noktasıdır"; havacılıktaki 121.5 MHz'in dengidir. Acil bir durumda ya da ilk temasta buradan çağrı yapılır, sonra çalışma kanalına geçilir. Bu kanalı dinlemek, deniz haberleşmesinin nabzını tutmanın en doğrudan yoludur.

### 4.2 AIS — Otomatik Kimlik Sistemi

AIS (Automatic Identification System), gemilerin kimliğini, konumunu, rotasını ve hızını otomatik olarak yayınladığı sayısal bir sistemdir; havacılıktaki ADS-B'nin deniz karşılığıdır. Açık ve çözülebilir.

| Frekans | Kullanım | Modülasyon | Not |
|---|---|---|---|
| **161.975 MHz (AIS 1 / Kanal 87B)** | Gemi konum/kimlik yayını | GMSK | Çift kanaldan biri |
| **162.025 MHz (AIS 2 / Kanal 88B)** | Gemi konum/kimlik yayını | GMSK | İki kanal arasında zaman paylaşımlı (TDMA) |

Pratikte: AIS, ucuz bir VHF alıcısıyla çözülerek yakındaki gemilerin canlı haritası çıkarılabilir; ADS-B uçaklar için neyse, AIS gemiler için odur. Kod çözme için bkz. Bölüm 5.

### 4.3 Deniz MF/HF ve acil/güvenlik sistemleri

| Frekans | Kullanım | Modülasyon | Not |
|---|---|---|---|
| **2182 kHz (2.182 MHz)** | MF imdat ve çağrı (sesli) | SSB (J3E) | Klasik MF imdat frekansı |
| **2187.5 kHz** | MF DSC | Sayısal | DSC acil çağrı (MF) |
| **4 / 6 / 8 / 12 / 16 MHz deniz HF bantları** | Uzun menzil gemi-kıyı | SSB (USB) | Okyanus ötesi; ITU kanal planına bağlı |
| **518 kHz** | NAVTEX (uluslararası) | FSK (sayısal teleks) | Seyir/hava uyarıları otomatik metin yayını; çözülebilir |
| **490 kHz / 4209.5 kHz** | NAVTEX (ulusal / tropik) | FSK | 518 kHz'in ulusal ve HF tamamlayıcıları |
| **406.0 – 406.1 MHz** | EPIRB / COSPAS-SARSAT acil konum verici | Sayısal burst | Gemi/can yeleği acil işareti; uydu üzerinden arama-kurtarmaya gider |
| **121.500 MHz** | Eski EPIRB homing (yardımcı) | AM | Artık ana acil değil; yön bulma yardımcısı olarak kullanılabilir |

Uyarı: 406 MHz EPIRB ve havacılıktaki 121.5 MHz gibi acil/imdat frekansları, can güvenliğine ayrılmıştır. Bu frekanslarda **yayın yapmak** (yanlışlıkla dahi) arama-kurtarma kaynaklarını seferber eder ve ağır yaptırıma tabidir. Dinleme ayrı bir konudur; bu frekanslarda asla test yayını yapılmaz.

---

<a id="5"></a>
## 5. Kara Askeri Taktik Haberleşme

Bu bölüm, kara taktik haberleşmenin hangi **bant aralıklarına** tahsis edildiğini, açık ve genel düzeyde gösterir. Burada verilen değerler, uluslararası tahsis tablolarında ilgili hizmetlere (kara mobil, sabit, taktik) ayrılmış **geniş bant sınırlarıdır.** Bu tablo bir **bant farkındalığı** aracıdır.

Kritik sınır: Spesifik operasyonel frekanslar, çağrı planları, atlama desenleri, şifreleme anahtarları ve içerik **açık kaynak değildir** ve bu bölümde kesinlikle yer almaz. "Şu birlik şu frekansı kullanır" türünde bir bilgi ne burada vardır ne de açık kaynaktan edinilebilir; böyle bir içeriğe ulaşmaya/çözmeye çalışmak yasadışıdır. Burada anlatılan, yalnızca taktik haberleşmenin hangi spektrum dilimlerinde *yaşadığı* ve hangi kavramlara (atlama, ağ) dayandığıdır.

| Bant aralığı (genel/açık) | Tipik taktik kullanım | Karakteristik | Not |
|---|---|---|---|
| **1.6 – 30 MHz (HF)** | Uzun menzil / ufuk ötesi taktik ve stratejik | İyonosferden yansıma; uydu gerekmeden kıtalararası | ALE (otomatik link kurulumu) ve sayısal modlar yaygın |
| **30 – 88 MHz (alt VHF)** | Kısa-orta menzil tabur/manga taktik telsiz | Görüş hattı; frekans atlamalı sistemler bu banttadır | SINCGARS ailesi bu aralıkta tasarlanmıştır (kavram aşağıda) |
| **108 – 174 MHz / 225 – 400 MHz (VHF/UHF)** | Taktik sesli, hava-kara, geri bağlantı | Görüş hattı; UHF kısmı askeri hava ile örtüşür | UHF üst-taktik ve uydu yer terminalleri |
| **225 – 400 MHz** | Askeri UHF (hava + kara köprü) | Bkz. Bölüm 3 | Hava ve kara için ortak askeri UHF bandı |
| **969 – 1206 MHz** | **Link-16** taktik veri ağı (JTIDS/MIDS) | TDMA + frekans atlama; sayısal veri | Aşağıda ayrı açıklanır |

### 5.1 SINCGARS ve VHF taktik frekans atlama (kavram)

SINCGARS (Single Channel Ground and Airborne Radio System), 30–88 MHz alt-VHF aralığında çalışan, NATO ülkelerinde yaygın bir taktik telsiz ailesidir. İki temel modu vardır: tek kanal (sabit frekans) ve **frekans atlama.** Atlama modunda, bir ağdaki tüm telsizler ortak bir atlama setine ve ortak bir ağ zamanına göre saniyede çok sayıda frekansa birlikte sıçrar. Dışarıdan bakıldığında, 30–88 MHz bandı boyunca kısa, dağınık enerji parçaları görünür; deseni ve ağ anahtarını bilmeyen biri kesintisiz bir kanal yakalayamaz.

Bu, Have Quick (Bölüm 3) ile aynı temel mantıktır, sadece farklı bant ve sistemde. Her ikisi de yayılı spektrumun (Bölüm 1) askeri uygulamasıdır: dinlemeye ve karıştırmaya karşı dayanıklılık.

### 5.2 Link-16 (kavram)

Link-16, modern taktik kuvvetlerin durumsal farkındalığı paylaştığı bir **sayısal veri ağıdır** (sesli değil, esasen veri). 969–1206 MHz aralığında çalışır ve iki tekniği birleştirir: zaman bölmeli çoklu erişim (TDMA — her terminale ağda zaman dilimleri ayrılır) ve frekans atlama. Bu birleşim, çok sayıda terminalin aynı ağda, birbirini karıştırmadan ve dış müdahaleye karşı dirençli biçimde konum, hedef ve durum verisi paylaşmasını sağlar.

Not: Link-16'nın bant aralığı (969–1206 MHz) açık bilgidir ve bu aralık aynı zamanda sivil havacılık seyrüsefer (DME/TACAN) ve gözetim sistemleriyle dikkatli koordine edilir. İçeriği şifrelidir ve açık değildir; burada anlatılan yalnızca sistemin hangi bantta ve hangi erişim mantığıyla çalıştığıdır.

Tekrar uyarı: Bu bölümün tamamı bant farkındalığı içindir. Hiçbir satır operasyonel kullanım, çözme ya da müdahale rehberi değildir. Taktik askeri trafiği dinlemeye/çözmeye/karıştırmaya yönelik her girişim yasadışıdır ve bu dokümanın kapsamı dışındadır.

---

<a id="6"></a>
## 6. Kamu Güvenliği — Polis, İtfaiye, Ambulans

Kamu güvenliği telsizleri (polis, itfaiye, ambulans, afet, belediye) tarihsel olarak VHF ve UHF bantlarında çalışmıştır. Son yirmi yılda bu hizmetler büyük ölçüde **sayısal ve şifreli** sistemlere geçmiştir; bu, dinleme açısından kritik bir ayrım yaratır.

Önemli yasal ayrım: Kamu güvenliği telsizinin *bant tahsisini* bilmek açık bilgidir. Ancak birçok ülkede (ve giderek artan biçimde) kamu telsizi trafiği **şifrelidir**; şifreli trafiği çözmek yasaktır. Şifresiz (analog ya da açık sayısal) kamu telsizi dinlemenin yasallığı ülkeye göre değişir — bazı ülkelerde serbest, bazılarında (dinlediğini aktarmak/kullanmak özellikle) yasaktır. Türkiye dahil birçok ülkede kamu güvenliği telsizine müdahale ve içeriğini ifşa ağır suçtur. Aşağıdaki tablo standart/bant farkındalığı içindir.

### 6.1 Sayısal kamu telsizi standartları

| Standart | Tipik bant | Bölge / yaygınlık | Karakteristik |
|---|---|---|---|
| **TETRA** | **380 – 400 MHz** (kamu güvenliği tahsisi); ayrıca 410–430 vb. ticari | Avrupa (ve birçok ülke) standardı | Trunked, sayısal, genellikle şifreli; grup çağrı, doğrudan mod |
| **P25 (APCO-25)** | VHF/UHF (136–174, 380–470, 700/800 MHz) | Esasen Kuzey Amerika | Sayısal kamu güvenliği; şifresiz ya da şifreli olabilir |
| **DMR** | VHF/UHF | Küresel, ticari + bazı kamu | Sayısal, iki zaman dilimli TDMA; ticari telsizde yaygın |
| **TETRAPOL** | UHF | Bazı Avrupa ülkeleri (alternatif) | TETRA'ya alternatif sayısal kamu sistemi |

### 6.2 Klasik (analog) kamu/ticari telsiz bantları

Sayısal sistemlere geçilmeden önce ve hâlâ bazı yerlerde, kamu ve ticari telsiz şu genel bantlarda bulunur (kesin atamalar ulusal tahsis tablosundan teyit edilmeli):

| Bant aralığı | Tipik kullanım | Modülasyon | Not |
|---|---|---|---|
| **VHF düşük (30–50 MHz)** | Eski kırsal/kamu telsiz | FM | Giderek terk edildi |
| **VHF yüksek (136–174 MHz)** | Kamu/ticari telsiz, bazı ambulans/itfaiye | FM (NBFM) | Klasik "telsiz bandı"nın çekirdeği |
| **UHF (380–470 MHz)** | Kamu güvenliği (sayısal dahil), ticari | FM / sayısal | TETRA kamu 380–400 burada |
| **PMR446 (446.0–446.2 MHz)** | Lisanssız el telsizi (Avrupa) | FM / dPMR | Halka açık, düşük güç; dinlemesi serbest |

Türkiye'de durum (genel): Türkiye'de kamu güvenliği için sayısal telsiz altyapısı kullanılmaktadır; bu sistemlerin trafiği genellikle şifrelidir ve içeriği kamuya açık değildir. Bant tahsisleri BTK Milli Frekans Planı'nda yer alır; spesifik kanal ve şifreleme bilgisi açık değildir ve dinlenmesi/çözülmesi yasaktır. Halka açık ve dinlenebilir tek kategori, lisanssız el telsizi (PMR446 benzeri) ve amatör bantlardır.

Pratikte: Bir UHF tarayıcıda 380–400 MHz civarında TETRA'nın karakteristik sayısal sesini (insan kulağına anlamsız, ritmik sayısal vınlama) duyabilirsin; bu, sistemin *varlığını* teyit eder ama içeriği şifreliyse çözülemez ve çözmeye çalışmak yasadışıdır. Sayısal sesli kod çözme (yalnızca şifresiz ve yasal sistemlerde) için bkz. Bölüm 5.

---

<a id="7"></a>
## 7. Amatör Radyo — Tam Bant Planı

Amatör (radyo amatörlüğü, "ham radio"), lisanslı bireylerin deney, haberleşme ve acil yardım için kullandığı, ITU'nun amatör hizmete tahsis ettiği bantlardır. Amatör bantlar, dinleme açısından tamamen açıktır ve öğrenmek için en zengin alanlardandır; bant içi mod düzeni (CW/SSB/sayısal/FM segmentleri) bir "bant planı" ile düzenlenir.

Not: Aşağıdaki sınırlar Region 1 (Türkiye dahil) için tipik değerlerdir. Bant kenarları ve segment sınırları (özellikle 40 m, 80 m, 60 m ve bazı üst bantlar) bölgeye ve ulusal düzenlemeye göre değişebilir; **kesin sınırlar için ulusal amatör bant planı (Türkiye'de BTK / TRAC) teyit edilmelidir.** "Tipik segment" sütunu, küresel amatör geleneğinin genel düzenini gösterir.

### 7.1 HF amatör bantları (Region 1 tipik)

| Bant adı | Frekans aralığı | Tipik mod segmentleri | Not |
|---|---|---|---|
| **160 m** | 1.810 – 2.000 MHz | Alt uç CW, üst SSB | Gece bandı; düşük frekans, geniş anten ister |
| **80 m** | 3.500 – 3.800 MHz (R1) | 3.500–3.580 CW, 3.580–3.600 sayısal, 3.600–3.800 SSB | Bölgesel/gece; R2'de üst kenar farklı (4.000) |
| **60 m** | 5.351.5 – 5.366.5 kHz (R1, dar) | Çoğunlukla USB / sınırlı | Yeni ve dar tahsis; ülkeye göre kanal/segment değişir |
| **40 m** | 7.000 – 7.200 MHz (R1) | 7.000–7.040 CW, 7.040–7.060 sayısal, 7.060–7.200 SSB | R2/R3'te 7.300'e kadar; yayın bandı çakışması tarihsel sorun |
| **30 m** | 10.100 – 10.150 MHz | CW ve sayısal (SSB yok) | WARC bandı; dar, sadece CW/data |
| **20 m** | 14.000 – 14.350 MHz | 14.000–14.070 CW, 14.070–14.099 sayısal, 14.101–14.350 SSB | En popüler DX bandı; gündüz uzak mesafe |
| **17 m** | 18.068 – 18.168 MHz | CW/sayısal alt, SSB üst | WARC bandı |
| **15 m** | 21.000 – 21.450 MHz | 21.000–21.070 CW, sayısal, üst SSB | Güneş aktivitesine duyarlı; iyi koşulda uzun DX |
| **12 m** | 24.890 – 24.990 MHz | CW/sayısal/SSB | WARC bandı |
| **10 m** | 28.000 – 29.700 MHz | 28.000–28.070 CW, 28.070–28.190 sayısal, FM üst (29.x repeater) | Güneş maksimumunda olağanüstü açılır; FM ve uydu da burada |

Not (CW ve SSB taban yan bandı geleneği): Amatör HF'te 10 MHz'in altında genellikle alt yan bant (LSB), üstünde üst yan bant (USB) kullanılır. Bu bir kural değil yerleşik gelenektir ama bir SSB sinyalini doğru çözmek için bilmek gerekir.

### 7.2 VHF / UHF amatör bantları (Region 1 tipik)

| Bant adı | Frekans aralığı | Tipik kullanım | Not |
|---|---|---|---|
| **6 m** | 50.000 – 52.000 MHz (R1) | CW/SSB alt uç, FM üst | "Sihirli bant"; sporadik-E ile beklenmedik uzun mesafe |
| **4 m** | 70.000 – 70.500 MHz | Bazı R1 ülkelerinde | Türkiye dahil her ülkede tahsisli değil; teyit gerekir |
| **2 m** | 144.000 – 146.000 MHz (R1) | 144.000–144.150 CW/SSB, üst FM/repeater, uydu | En yaygın VHF amatör bandı; yerel FM ve uydu |
| **70 cm** | 430.000 – 440.000 MHz (R1) | SSB/CW alt, FM/repeater, sayısal, uydu | İkincil statü olabilir (radyolokasyonla paylaşımlı) |
| **23 cm** | 1240 – 1300 MHz | Mikrodalga giriş; ATV, sayısal | İkincil; seyrüseferle paylaşımlı, dikkat |
| Üst mikrodalga (13 cm, 9 cm, 6 cm, 3 cm...) | 2.3 / 3.4 / 5.6 / 10 GHz dilimleri | Deneysel, mikrodalga DX | Çoğu ikincil; ISM ile örtüşen bölümler var |

### 7.3 Amatör uydu segmentleri

Amatör radyo, kendi uydularını (OSCAR serisi ve çok sayıda kübsat) işletir. Uydu haberleşmesi, bant içinde ayrılmış belirli segmentlerde yapılır ve genellikle bir bantta yukarı (uplink), başka bir bantta aşağı (downlink) bağlantı kullanır.

| Segment | Bant | Kullanım | Not |
|---|---|---|---|
| **2 m uydu segmenti** | 145.800 – 146.000 MHz civarı | Uydu uplink/downlink | ISS dahil; APRS ve sesli röle |
| **70 cm uydu segmenti** | 435 – 438 MHz civarı | Uydu uplink/downlink | V/U ve U/V transponder kombinasyonları |
| **ISS / amatör** | 145.800 (downlink) civarı | Uzay istasyonu amatör trafiği | Dinlemesi serbest; geçiş zamanları hesaplanabilir |

Pratikte: Amatör bantlar, dinleme yasası açısından en serbest, içerik açısından en açık ve deney açısından en zengin alandır. Bir bant planını "okumayı" (segment → mod → ne duyulur) öğrenmenin en iyi yolu, kendi bölgendeki 2 m ya da 40 m bandını izleyip duyduğun sinyali plandaki segmentle eşleştirmektir.

---

<a id="8"></a>
## 8. ISM / SRD — Lisanssız Kısa Menzil

ISM (Industrial, Scientific and Medical) ve SRD (Short Range Devices) bantları, lisans gerektirmeden, düşük güçle, herkesin belirli kurallar dahilinde kullanabildiği bantlardır. Günlük hayattaki kablosuz cihazların (uzaktan kumanda, IoT sensör, WiFi, Bluetooth, kablosuz kulaklık) büyük kısmı buradadır. Dinlemesi serbesttir; bu bantlar `rtl_433` gibi araçlarla en çok deney yapılan alanlardandır.

Not: ISM bantları lisanssız ve "paylaşımlı" olduğu için, aynı bantta çok sayıda farklı cihaz bir arada bulunur. Bir frekansın "kime ait" olduğu burada tek bir cihaza değil, bir cihaz **sınıfına** karşılık gelir. Ülkeye göre tam sınırlar ve izinli güç değişir; kesin değer için ulusal SRD düzenlemesi (Türkiye'de BTK kısa mesafe cihaz tebliği) teyit edilmelidir.

| Frekans aralığı | Bölge | Tipik kullanım | Not |
|---|---|---|---|
| **13.553 – 13.567 MHz (13.56 MHz)** | Küresel | RFID/NFC, endüstriyel | Temassız kart, etiket |
| **26.957 – 27.283 MHz (27 MHz)** | Küresel | Eski oyuncak kumandaları, CB yakını | Düşük güç model/kumanda |
| **40.66 – 40.70 MHz** | Küresel | Telemetri, kumanda | Dar ISM dilimi |
| **433.05 – 434.79 MHz (433 MHz)** | Region 1 (Avrupa dahil) | Uzaktan kumanda, garaj, hava istasyonu, TPMS, sensör | Avrupa'nın en yoğun SRD bandı; `rtl_433` cenneti |
| **863 – 870 MHz (868 MHz)** | Region 1 (Avrupa) | IoT (LoRa, sensör), akıllı sayaç, alarm | Avrupa kısa menzil/IoT; LoRaWAN burada |
| **902 – 928 MHz (915 MHz)** | Region 2 (Amerika) | IoT, kumanda, bazı LoRa | Amerika muadili; Avrupa'da bu blok farklı tahsisli |
| **2400 – 2483.5 MHz (2.4 GHz)** | Küresel | WiFi (2.4G), Bluetooth/BLE, ZigBee, mikrodalga fırın | En kalabalık ISM bandı; her şey burada |
| **5725 – 5875 MHz (5.8 GHz)** | Küresel (bölgesel farklar) | WiFi (5G kısmı), bazı kumandalar, FPV video | Daha geniş, daha az kalabalık |

Uyarı (433 ve 868/915 farkı): Avrupa (Region 1) ve Amerika (Region 2) arasında en kritik SRD farkı, 868 MHz (Avrupa) ile 915 MHz (Amerika) ayrımıdır. Türkiye Region 1'dedir; bu yüzden Türkiye'de IoT/LoRa cihazları 868 MHz blokunda olmalıdır, 915 MHz Amerika'ya yönelik bir cihaz Türkiye'de yasal olmayabilir. Bir 915 MHz sinyali Türkiye'de duyulması beklenmeyen bir durumdur ve genellikle yanlış bölge için üretilmiş bir cihaza işaret eder.

Pratikte: 433 MHz bandını bir ucuz alıcı ve `rtl_433` ile dinlemek, çevredeki kablosuz hava istasyonlarını, araç lastik basıncı (TPMS) sensörlerini, kapı zillerini ve bazı uzaktan kumandaları pasifçe çözüp, bant tahsisi ile gerçek dünyadaki cihazları eşleştirmenin en somut yoludur. Bu, hem bir alıştırma hem de IoT güvenlik farkındalığı dersidir.

---

<a id="9"></a>
## 9. Uydu, Seyrüsefer ve Yayın

Bu bölüm, üç büyük "yukarıdan gelen" kullanıcı grubunu toplar: küresel seyrüsefer uyduları (GNSS), haberleşme/hava durumu uyduları ve karasal yayın (radyo/TV). Bunların çoğu açık, çözülebilir ve dinlemesi serbest sistemlerdir.

### 9.1 Küresel seyrüsefer uyduları (GNSS)

GNSS (Global Navigation Satellite System), konum/zaman sağlayan uydu takımlarıdır. Hepsi L bandında (1–2 GHz) çalışır; sinyaller çok zayıftır (gürültü tabanının altında, korelasyonla çözülür).

| Sistem / sinyal | Frekans | Sahip / takım | Not |
|---|---|---|---|
| **GPS L1 (C/A)** | **1575.42 MHz** | ABD | Sivil temel sinyal; en yaygın çözülen |
| **GPS L2** | 1227.60 MHz | ABD | İkinci frekans (çoğunlukla askeri/çift frekans alıcı) |
| **GPS L5** | 1176.45 MHz | ABD | Yeni sivil güvenlik-kritik sinyal (havacılık) |
| **GLONASS L1** | ~1598–1606 MHz (1602 nominal) | Rusya | FDMA tabanlı; kanal frekansları kaydırmalı |
| **GLONASS L2** | ~1242–1248 MHz | Rusya | İkinci bant |
| **Galileo E1** | 1575.42 MHz (GPS L1 ile örtüşür) | Avrupa (AB) | E1/E5/E6; L1 üzerinde GPS ile birlikte yaşar |
| **Galileo E5a/E5b** | 1176.45 / 1207.14 MHz | Avrupa | Geniş bant, yüksek doğruluk |
| **BeiDou B1** | ~1561.098 / 1575.42 MHz | Çin | B1/B2/B3; küresel kapsama |

Uyarı: GNSS dinlemek (alıcıyla konum çözmek) tamamen serbesttir. Ancak GNSS bandında **yayın yapmak** — sahte sinyal (spoofing) ya da karıştırma (jamming) — neredeyse her ülkede ağır suçtur; uçak, gemi ve kara seyrüseferini doğrudan tehlikeye atar. Bu bölüm yalnızca bantların *alımı* hakkındadır.

### 9.2 Haberleşme ve mobil uydu sistemleri

| Sistem | Frekans (kullanıcı linki) | Kullanım | Not |
|---|---|---|---|
| **Inmarsat** | ~1.525–1.559 GHz (downlink), ~1.626–1.660 GHz (uplink) | Deniz/hava/kara mobil uydu haberleşme | L bandı; bazı yayınları (STD-C) çözülebilir |
| **Iridium** | **1616 – 1626.5 MHz** | Küresel uydu telefon/veri (LEO takım) | Alçak yörünge; karakteristik "ring alert" burst'leri |
| **Thuraya** | ~1.5/1.6 GHz | Bölgesel mobil uydu telefon | Orta Doğu/Asya/Afrika kapsama |
| **AMSAT / amatör uydu** | 145/435 MHz (bkz. Bölüm 7.3) | Amatör uydu | Açık; dinlemesi serbest |

### 9.3 Hava durumu ve görüntü uyduları

| Sistem | Frekans | Kullanım | Not |
|---|---|---|---|
| **NOAA APT** | ~137 MHz (137.1 / 137.62 / 137.9125 vb.) | Düşük yörünge hava uydusu görüntü (analog APT) | Ucuz alıcıyla canlı görüntü indirilir; klasik proje |
| **Meteor-M LRPT** | ~137 MHz | Rus hava uydusu sayısal görüntü | APT'nin sayısal/daha keskin muadili |
| **GOES (HRIT/EMWIN)** | ~1.69 GHz | Sabit yörünge (jeostasyoner) hava uydusu | Sürekli yarıküre görüntüsü; daha büyük anten/LNA ister |
| **GNSS hava/yardımcı yayınlar** | çeşitli | SBAS (EGNOS/WAAS) seyrüsefer düzeltme | GPS L1 yakını; havacılık hassasiyet artırma |

Pratikte: 137 MHz NOAA APT, ucuz bir alıcı ve basit bir antenle gökyüzünden canlı uydu görüntüsü almanın en erişilebilir yoludur ve bir uydunun geçişini, frekansını ve modülasyonunu (bant planından) eşleştirmenin somut bir alıştırmasıdır.

### 9.4 Karasal yayın (radyo / TV)

| Bant / aralık | Hizmet | Modülasyon | Not |
|---|---|---|---|
| **148.5 – 283.5 kHz (LF)** | Uzun dalga (LW) yayın | AM | Region 1'e özgü; Amerika'da yok |
| **526.5 – 1606.5 kHz (MF)** | Orta dalga (MW) AM yayın | AM | Klasik AM radyo; gece çok uzaktan gelir |
| **HF yayın bantları (49 m, 41 m, 31 m, 25 m, 19 m, 16 m, 13 m)** | Kısa dalga (SW) yayın | AM (ve bazı DRM) | Uluslararası yayın; iyonosferle kıtalararası |
| **87.5 – 108.0 MHz (VHF Band II)** | FM radyo yayını | FM (geniş bant) | Yerel FM; RDS veri alt taşıyıcısı çözülebilir |
| **174 – 240 MHz (VHF Band III)** | DAB / DAB+ sayısal radyo, bazı TV | OFDM (DAB) | Ülkeye göre DAB ya da TV |
| **470 – 694 MHz (UHF)** | DVB-T / DVB-T2 sayısal TV | OFDM | Karasal sayısal TV; ülkeye göre kanal planı |

Not: FM yayın bandı (87.5–108 MHz), bir alıcının en kolay duyduğu güçlü sinyallerdir ve aynı zamanda komşu bantlarda (havacılık 108 MHz hemen üstte, deniz/VHF) zayıf sinyalleri "boğan" başlıca girişim kaynağıdır. Bu yüzden 8-bit ucuz alıcılarda FM çentik (notch) filtresi sık kullanılır (bkz. Bölüm 2 ve 3).

---

<a id="10"></a>
## 10. GSM / LTE / 5G — Hücresel Bantlar

Hücresel (cep telefonu) haberleşme, kendine ayrılmış sayısal bantlarda çalışır ve **şifrelidir.** Bu bölüm bant tahsisi farkındalığı içindir; hücresel trafiğin içeriği şifreli olduğundan çözülemez ve çözmeye çalışmak yasadışıdır. Burada anlatılan, yalnızca hangi bant numarasının hangi frekans aralığına karşılık geldiği ve uplink/downlink mantığıdır.

Temel kavram (FDD uplink/downlink): Çoğu hücresel bant **FDD** (Frequency Division Duplex) çalışır: telefondan kuleye (uplink) ve kuleden telefona (downlink) ayrı frekans bloklarında, eşzamanlı. Bir bant numarası (örn. Band 1, Band 3) belirli bir uplink ve belirli bir downlink aralığını birlikte tanımlar. Bazı yeni bantlar **TDD** (zaman bölmeli) çalışır: aynı frekansta sırayla yön değiştirir.

### 10.1 Yaygın hücresel bantlar (3GPP bant numaraları)

| Bant | Uplink (telefon→kule) | Downlink (kule→telefon) | Dupleks | Yaygın ad / not |
|---|---|---|---|---|
| **Band 1** | 1920 – 1980 MHz | 2110 – 2170 MHz | FDD | 2100 MHz; 3G/4G yaygın |
| **Band 3** | 1710 – 1785 MHz | 1805 – 1880 MHz | FDD | 1800 MHz; LTE'nin en yaygın bandı (Avrupa/Türkiye) |
| **Band 7** | 2500 – 2570 MHz | 2620 – 2690 MHz | FDD | 2600 MHz; yüksek kapasite LTE |
| **Band 8** | 880 – 915 MHz | 925 – 960 MHz | FDD | 900 MHz; GSM mirası, geniş kapsama |
| **Band 20** | 832 – 862 MHz | 791 – 821 MHz | FDD | 800 MHz; "dijital temettü", kırsal LTE kapsama |
| **Band 28** | 703 – 748 MHz | 758 – 803 MHz | FDD | 700 MHz; geniş kapsama LTE/5G |
| **Band 38 / n78** | 2570 – 2620 MHz / 3300 – 3800 MHz | (TDD, aynı blok) | TDD | 2600 TDD ve 3.5 GHz 5G ana bandı |

Not (GSM tarihsel bantları): Klasik GSM, 900 MHz (GSM-900, kabaca Band 8 aralığı) ve 1800 MHz (GSM-1800/DCS, kabaca Band 3 aralığı) bantlarında kurulmuştu. LTE ve 5G büyük ölçüde aynı fiziksel bantları yeniden kullanır (refarming); bu yüzden bir frekans bloğu zamanla GSM'den LTE'ye, oradan 5G'ye dönüşebilir.

### 10.2 Türkiye'de operatör bantları (genel)

Türkiye'de mobil hizmet, BTK tarafından operatörlere ihale/tahsis edilen bloklarla sunulur. Genel olarak kullanılan bantlar 800 (Band 20), 900 (Band 8), 1800 (Band 3), 2100 (Band 1), 2600 (Band 7) ve 5G için 3.5 GHz (n78) çevresindedir. Hangi operatörün hangi bloğa (hangi frekans aralığına) sahip olduğu BTK'nın yayımladığı tahsis kararlarıyla belirlenir ve zaman içinde değişir.

Kritik: Hangi operatörün hangi alt bloğu kullandığı **ulusal tahsis kararından (BTK) teyit edilmelidir**; bu dağılım ülkeye ve döneme özgüdür ve burada operatör-blok eşlemesi verilmemiştir. Bu bölüm yalnızca bant numarası → frekans aralığı düzeyinde farkındalık sağlar.

Uyarı: Hücresel trafik şifrelidir. Bir spektrum analizöründe bu bloklarda yoğun, geniş bant sayısal enerji görmek mümkündür (sistemin varlığını teyit eder) ama içerik çözülemez ve çözmeye/araya girmeye (IMSI catcher vb.) çalışmak ağır suçtur ve bu dokümanın kapsamı dışındadır.

---

<a id="11"></a>
## 11. Bir Frekansın Sahibini Belirleme Metodolojisi

Bir tarayıcıda ya da şelale (waterfall) görünümünde bilinmeyen bir sinyal yakaladığında, "bu kime ait, ne işe yarar" sorusunu tahmine değil, sistematik bir akıl yürütmeye dayandırabilirsin. Aşağıdaki adımlar, bilinmeyenden kimliğe giden bir daraltma sürecidir.

### Adım 1 — Frekansı oku ve bant dilimine yerleştir (kaba daraltma)

İlk iş, frekansın hangi büyük dilimde (HF/VHF/UHF/SHF) olduğunu ve bu bölümdeki hangi tabloya düştüğünü belirlemektir. Bu tek adım, kullanıcı grubunu büyük ölçüde daraltır:

- 118–137 MHz arası, AM → büyük olasılıkla sivil havacılık ATC.
- 156–162 MHz arası, FM → büyük olasılıkla deniz VHF.
- 380–400 MHz arası, sayısal "vınlama" → muhtemelen TETRA kamu/ticari.
- 433 MHz civarı, kısa burst'ler → ISM/SRD cihazları (sensör, kumanda).
- 1090 MHz, kısa darbeler → ADS-B/transponder.
- 1575.42 MHz → GPS L1 (çok zayıf, korelasyonsuz duyulmaz).

### Adım 2 — Tahsisi doğrula (tablo/veritabanı)

Bant diliminden çıkan hipotezi, tahsis tablosuyla doğrula. Frekans, ilgili hizmete tahsisli mi? Birincil mi ikincil mi? Bu adımda bu dokümanın tabloları ya da ulusal tahsis tablosu (BTK) ve çevrimiçi tahsis veritabanları (Bölüm 12) kullanılır. Tahsis, "kim olabilir" listesini netleştirir.

### Adım 3 — Olası kullanımı daralt (kanal planı / yerel bilgi)

Tahsis edilmiş bir bant içinde, kanal planı ve yerel bağlam kullanımı daha da daraltır. Örnek: 118–137 MHz havacılık ise, kanal aralığı (8.33 ya da 25 kHz) ve yerel havaalanı yayınları (AIP) hangi tür istasyon olduğunu (kule, yaklaşma, saha) söyler. 156–162 MHz deniz ise, kanal numarası (16, 70, vb.) işlevi verir.

### Adım 4 — Modülasyonu teyit et (kimlik doğrulama)

Son adım, modülasyon tipini dinleyerek/ölçerek hipotezi doğrulamaktır. Modülasyon, kullanıcı grubunun karakteristik imzasıdır:

| Duyduğun / gördüğün | Olası kimlik |
|---|---|
| AM ses (uğultulu, capture yok) | Havacılık ATC/airband, askeri UHF hava |
| FM ses (temiz, dar bant) | Deniz VHF, amatör FM, kamu/ticari analog telsiz |
| SSB ses (kısık, "ördek" sesi, doğru BFO ile netleşir) | HF amatör, deniz/havacılık HF, taktik HF |
| Sayısal ritmik vınlama (anlamsız ses) | TETRA/P25/DMR, sayısal sistemler (çoğu şifreli) |
| Kısa periyodik burst (sensör) | ISM 433/868, telemetri |
| Sürekli geniş bant sayısal "duvar" | GSM/LTE/5G, sayısal yayın (DVB-T/DAB) |
| Kısa darbe çiftleri, ~1090 MHz | ADS-B/transponder |
| Otomatik teleks/FSK metin | NAVTEX (518 kHz), POCSAG/FLEX çağrı |

Pratikte: Bu dört adım (dilim → tahsis → kanal/bağlam → modülasyon) birlikte uygulandığında, bilinmeyen sinyallerin büyük çoğunluğu birkaç dakikada makul biçimde kimliklendirilebilir. Şifreli sayısal sistemlerde "kimlik" tespit edilebilir (örn. "bu TETRA") ama içerik çözülemez ve çözmeye çalışmak yasadışıdır. Modülasyon ve demodülasyon ayrıntıları için bkz. Bölüm 5; spektrum/şelale okuma için bkz. Bölüm 1.

---

<a id="12"></a>
## 12. Tahsis Veritabanları ve Kaynaklar

Frekans kimliklendirme, hafızadaki tablolardan çok güncel ve resmi veritabanlarına dayanmalıdır. Aşağıdaki kaynaklar, açık tahsis ve kullanım bilgisinin başlıca dayanaklarıdır.

| Kaynak | Tür | Ne sağlar |
|---|---|---|
| **ITU Radio Regulations / Table of Frequency Allocations** | Uluslararası resmi | Hizmet bazlı küresel tahsis çerçevesi; Region 1/2/3 farkları |
| **BTK Milli Frekans Planı (Türkiye)** | Ulusal resmi | Türkiye'nin bant tahsisleri, kanal planları, SRD/amatör kuralları |
| **FCC Online Table of Frequency Allocations (ABD)** | Ulusal resmi | ABD sivil/federal tahsis; karşılaştırma için yararlı |
| **CEPT / ECC ECA Table (Avrupa)** | Bölgesel uyum | Avrupa ortak tahsis tablosu; Region 1 uyumu |
| **RadioReference.com** | Topluluk veritabanı | Yerel/ülke bazlı bilinen frekanslar, kamu telsizi (yasal olduğu yerde), tarayıcı bilgisi |
| **ITU Appendix 18 (deniz VHF)** | Uluslararası resmi | Deniz VHF kanal numarası → frekans tam listesi |
| **Ulusal AIP (havacılık yayını)** | Resmi havacılık | Havaalanı/saha ATC frekansları, HF aile planı |
| **TRAC / ulusal amatör kuruluş bant planı** | Amatör otorite | Region 1/ulusal amatör bant ve segment sınırları |
| **3GPP bant tanımları** | Standart | Hücresel bant numarası → uplink/downlink eşlemesi |

Not: Topluluk veritabanları (RadioReference vb.) son derece yararlıdır ama resmi tahsis değildir; bir frekansın yasal statüsü ya da güncel kanal planı için her zaman **resmi ulusal kaynak (BTK)** esas alınmalıdır. Topluluk verisi "ne duyulmuş" gösterir; resmi tablo "ne tahsis edilmiş" söyler — ikisi farklı sorulardır.

---

<a id="13"></a>
## 13. Alıştırmalar (Yasal)

Amaç: Bant planını ezberlemekten çok, bir frekansı kimliklendirme refleksini ve resmi kaynak kullanma alışkanlığını pekiştirmek. Aşağıdaki alıştırmaların tamamı, dinlemesi açık/serbest olan bantlar ve kamuya açık veritabanları üzerinedir. Önce kendin çöz, sonra cevaba bak.

### Alıştırma 1 — Kendi bölgenin airband ATC'sini dinle (yasal)

En yakın havaalanının kule (tower) frekansını ulusal AIP'den ya da bir havacılık veritabanından bul, bir alıcıda **AM** modunda ve 8.33/25 kHz uygun adımlamayla dinle. Konuşmaların düzenini (uçak çağrı işareti, kalkış/iniş izni) gözlemle.

<details><summary>Beklenen kazanım</summary>

118–137 MHz havacılık bandının AM olduğunu, kanal aralığının (Avrupa'da çoğunlukla 8.33 kHz) önemini ve sesli ATC'nin nasıl yapılandığını doğrudan deneyimlersin. Bu, "frekans → tahsis → kanal → modülasyon" zincirinin (Bölüm 11) gerçek bir uygulamasıdır. Not: Dinleme serbesttir; duyduğunu yaymak ya da operasyonel kullanmak ayrı bir konudur ve ülkeye göre kısıtlı olabilir.
</details>

### Alıştırma 2 — Deniz VHF Kanal 16'yı kimliklendir

Kıyıya yakın bir yerdeysen, 156.800 MHz'i (Kanal 16) **FM** modunda izle. Havacılıktan farkını (FM vs AM) ve bu kanalın işlevini (imdat/çağrı) gözlemle. Mümkünse bir AIS alıcısıyla 161.975/162.025 MHz'i çözüp yakındaki gemilerin kimliğini gör.

<details><summary>Beklenen kazanım</summary>

Deniz VHF'in FM olduğunu (havacılığın AM'inin tersi), Kanal 16'nın denizciliğin ortak acil/çağrı noktası olduğunu (havacılıktaki 121.5 MHz'in dengi) ve AIS'in gemiler için ADS-B'nin karşılığı olduğunu kavrarsın. Modülasyon farkı (AM/FM), iki hizmeti ayırmanın en hızlı yoludur.
</details>

### Alıştırma 3 — Bilinmeyen bir frekansı bant planından kimliklendir

Sana üç frekans verildi: (a) 162.025 MHz, GMSK sayısal burst; (b) 145.500 MHz, FM ses; (c) 1090 MHz, kısa darbe çiftleri. Yalnızca bu bölümün tablolarını kullanarak her birinin kullanıcı grubunu ve işlevini belirle.

<details><summary>Çözüm</summary>

- **(a) 162.025 MHz, GMSK** → Deniz VHF bandı içinde, AIS 2 (Kanal 88B) frekansı. GMSK ve 162.025 değeri birlikte AIS'i işaret eder; gemi konum/kimlik yayını (Bölüm 4.2).
- **(b) 145.500 MHz, FM** → 2 m amatör bandı (144–146 MHz) içinde, FM çağrı/simpleks bölgesi. FM ses + bu frekans = amatör radyo (Bölüm 7.2). Deniz değildir (deniz 156+ MHz), havacılık değildir (havacılık AM ve 118–137).
- **(c) 1090 MHz, kısa darbe çiftleri** → Transponder/ADS-B (Bölüm 2.3). Bu frekans ve darbe modülasyonu havacılık gözetimine özgüdür; sesli değildir, `dump1090` ile çözülür.

Ders: Frekans + modülasyon birlikte, çoğu sinyali tek tabloya indirger. Tek başına frekans bazen iki hizmete uyabilir (komşu bantlar); modülasyon ayrımı bağı koparır.
</details>

### Alıştırma 4 — Ulusal tahsis tablosunu (BTK) incele

BTK Milli Frekans Planı'nı (ya da bulunduğun ülkenin tahsis tablosunu) aç ve şu üç bandın nasıl tahsis edildiğini bul: (1) 433 MHz SRD, (2) 868 MHz IoT, (3) 2 m amatör. Her biri için birincil/ikincil statüyü ve varsa güç sınırını not et.

<details><summary>Beklenen kazanım</summary>

Resmi bir tahsis tablosunu okumayı, bu dokümandaki "tipik" değerlerle ülkenin gerçek tahsisini karşılaştırmayı ve ülke/bölge bağımlı değerlerin neden "ulusal tablodan teyit edilmeli" notuyla verildiğini somut olarak görürsün. Ayrıca Türkiye'nin Region 1'de olmasının (868 MHz var, 915 MHz Amerika'ya ait) pratik sonucunu doğrularsın.
</details>

---

<a id="14"></a>
## 14. Çapraz Referans ve Sonraki Bölümler

Bu bölüm, SIGINT El Kitabı'nın frekans tahsisi ve bant planı parçasıdır. Bir frekansın "kime ait" olduğunu belirledikten sonra, o sinyali yakalamak, çözmek ve yasal çerçevede ele almak diğer bölümlerin konusudur.

Kapanış: "Bir frekans duydum, kime ait?" sorusunun cevabı sihir değil, sistemdir: önce hangi dilimde olduğuna bak, sonra tahsis tablosuyla hizmeti daralt, kanal planı ve modülasyonla kimliği doğrula. Bu bölümün tabloları o sistemin haritasıdır; ancak küresel hizmetler (havacılık, denizcilik, GNSS) dışında ülke/bölge bağımlı her değer, kesin karar için ulusal tahsis tablosundan (Türkiye'de BTK) teyit edilmelidir. Ve sınır nettir: tahsisi bilmek ve açık bantları dinlemek farkındalıktır; şifreli içeriği çözmek, acil/seyrüsefer/askeri bantlarda yayın yapmak ya da araya girmek suçtur. Bu rehber haritayı verir; haritanın hangi yollarının yasal olduğunu kendi ülkenin kuralları belirler.

---

Bu bölüm, Kanije Kalesi SIGINT El Kitabı'nın parçasıdır. Tüm bölümler ve önerilen okuma sırası için indekse bakın: [SIGINT_00 — Başlangıç ve İndeks](SIGINT_00_BASLANGIC_INDEX_VE_YASAL.md).

Doğrudan ilgili bölümler:
- [SIGINT_01 — RF Fiziği ve Modülasyon](SIGINT_01_TEMELLER_RF_VE_MODULASYON.md): bant mantığının ve yayılımın fiziksel temeli.
- [SIGINT_05 — Protokoller ve Sinyal Çözümleme](SIGINT_05_PROTOKOLLER_VE_SINYAL_COZUMLEME.md): bu bölümde tanımlanan sinyallerin içeriğine erişme.
- [SIGINT_09 — Yer Tespiti, Yön Bulma ve Takip](SIGINT_09_YER_TESPITI_YON_BULMA_VE_TAKIP.md): bir frekansın nereden geldiğini bulma.
- [SIGINT_10 — GNSS/GPS Sistemleri](SIGINT_10_GNSS_GPS_SISTEMLERI.md): GNSS L bandı tahsislerinin sistem tarafı.
- [SIGINT_06 — Güvenlik, Açıklar ve Savunma](SIGINT_06_GUVENLIK_ACIKLAR_VE_SAVUNMA.md): hangi içeriği çözmenin/yaymanın yasal sınırı.
