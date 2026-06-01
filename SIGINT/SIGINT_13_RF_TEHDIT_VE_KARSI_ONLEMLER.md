# SIGINT EL KİTABI — BÖLÜM 13: RF TEHDİT MANZARASI VE KARŞI-ÖNLEMLER

## Karıştırma, Drone/İHA RF Saldırıları ve Sinyal Manipülasyonu — Tehdidi Tanıma, Tespit, Savunma ve Yasal Çerçeve

> Amaç: Önceki bölümler sinyalin fiziğini (Bölüm 1), donanımı (Bölüm 2), anteni ve yön bulmayı (Bölüm 3), demodülasyonu (Bölüm 4), hücresel/IoT yüzeyini (Bölüm 5), TEMPEST ve sızıntı savunmasını (Bölüm 6), ayıklama ve sınıflandırmayı (Bölüm 7) verdi. Bu bölüm, üzerlerine binen bir savunma sorusunu ele alır: RF spektrumunu hedef alan aktif tehditler — karıştırma (jamming), drone/İHA radyo bağlarına yönelik saldırılar, replay/spoofing/relay manipülasyonları — fiziksel olarak nasıl çalışır, neden belirli sinyaller savunmasızdır, bu tehditler nasıl TESPİT edilir, nasıl SAVUNULUR ve yasal çerçeve nedir? Hedefimiz bir saldırganın icra reçetesi değil; bir savunmacının, kritik altyapı koruma analistinin ve siber tehdit istihbaratı (CTI) uzmanının tehdidi tanıyıp ölçebilmesidir.

> Yasal çerçeve (bu bölümün tamamı için bağlayıcı): Bu bölüm anlama, tespit ve savunma amaçlıdır. Karıştırma (jamming), bir drone/uydu/araç/sisteme yetkisiz müdahale, GPS/GNSS sahteleme (spoofing), keyless araç relay saldırısı ve benzeri aktif RF müdahaleleri dünyanın hemen her ülkesinde ağır suçtur; ayrıca havacılık, denizcilik, acil çağrı ve tıbbi sistemleri etkilediğinde doğrudan can güvenliği tehdididir. Bu metin hiçbir çalışan karıştırıcı (jammer) yapımını, devresini, komponent reçetesini; hiçbir drone/uydu/sistem ele geçirme (hijack) adımını; hiçbir aktif saldırı operasyonel reçetesini vermez ve veremez. Anlatılan saldırı sınıfları yalnızca kavramsal taksonomi düzeyindedir; amaç, savunmacının tehdidi spektrumda tanıması ve karşısına doğru kontrolü koymasıdır. Her aktif teknikten önce yasallık sınırını kendi ülkenin ve sürümünün mevzuatından teyit et.

---

## İÇİNDEKİLER

1. [Elektronik Harp Taksonomisi: EA, EP, ES ve Sivil/Güvenlik Karşılığı](#1)
2. [Karıştırma Fiziği: Sinyal-Karıştırma Oranı (J/S) ve Neden Zayıf Sinyaller Kolay Bastırılır](#2)
3. [Karıştırıcı Türleri (Kavramsal Sınıflandırma — Yapım Yok)](#3)
4. [Karıştırmanın Etkileri: GPS, GSM, Havacılık, Acil — Can Güvenliği Boyutu](#4)
5. [Karıştırma Tespiti: Spektrum İzleme, Gürültü Tabanı, Anomali ve Yön Bulma](#5)
6. [Karıştırmaya Karşı Savunma: Yayılı Spektrum, Atlama, Null-Steering, Güç Kontrolü, Yedeklilik](#6)
7. [Drone/İHA RF Güvenliği: Tipik Bağlar ve Neden Savunmasız](#7)
8. [Drone Saldırı Sınıfları (Yalnızca Kavram) ve Counter-UAS Savunması](#8)
9. [Replay / Spoofing / Relay Tehditleri: Prensip ve Savunma](#9)
10. [Sinyal Dayanıklılığı: Hangi Sinyaller Manipülasyona Açık/Kapalı — Tasarım Dersleri](#10)
11. [TEMPEST ve Yan-Kanal Yayılım: Savunma Farkındalığı](#11)
12. [Kritik Altyapı RF Riski: GPS/Zaman, SCADA Telemetri ve Savunma](#12)
13. [Tehdit → Tespit → Savunma Matrisi (Birleşik Referans)](#13)
14. [Alıştırmalar (Yalnızca Yasal/Savunma)](#14)
15. [Kapanış: Savunmacı Zihniyeti, Etik ve Yasal Sınır + Çapraz Referans](#15)

---

<a id="1"></a>
## 1. Elektronik Harp Taksonomisi: EA, EP, ES ve Sivil/Güvenlik Karşılığı

RF tehdidini doğru konumlandırmak için önce kavramsal çerçeveyi kurmak gerekir. Elektronik harp (EW, Electronic Warfare) doktrini, elektromanyetik spektrumu kullanan faaliyetleri üç ana fonksiyona ayırır. Bu ayrım askeri kökenlidir ama sivil güvenlik, kritik altyapı koruma ve siber tehdit istihbaratı dünyasında birebir karşılığı vardır. Savunmacı için değeri şudur: bir olayla karşılaştığında onun hangi fonksiyona ait olduğunu adlandırabilmek, doğru kontrolü seçmenin ilk adımıdır.

| EW fonksiyonu | Açılım | Özü | Sivil/güvenlik karşılığı |
|---|---|---|---|
| EA | Electronic Attack (Elektronik Saldırı) | Spektrumu düşmanın aleyhine kullanma: karıştırma, aldatma, sahteleme | Yasadışı jammer kullanımı, GPS spoofing olayı, keyless relay hırsızlığı — savunmacı için "tehdit" |
| EP | Electronic Protection (Elektronik Koruma) | Kendi sistemini EA ve istem dışı parazitten koruma | FHSS/DSSS, kimlik doğrulamalı protokol, anten null-steering, sağlam GNSS alıcısı — savunmacının inşa ettiği dayanıklılık |
| ES | Electronic Support (Elektronik Destek) | Spektrumu izleyip tehdidi tespit/tanı/yer tespiti yapma | Spektrum izleme, gürültü tabanı takibi, RF sensör ağı, drone RF tespiti — savunmacının "gözü" |

Not: Eski doktrinde aynı üçlü ECM (Electronic Countermeasures), ECCM (Electronic Counter-Countermeasures) ve ESM (Electronic Support Measures) olarak adlandırılırdı. Modern terimlerle EA≈ECM, EP≈ECCM, ES≈ESM örtüşür. Literatürde her iki sözlüğe de rastlanır.

```
                   ELEKTROMANYETİK SPEKTRUM
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
       EA                  EP                  ES
  (Elektronik Saldırı)  (Elektronik Koruma) (Elektronik Destek)
        │                   │                   │
   karıştırma          dayanıklı tasarım    izleme/tespit
   aldatma/sahteleme    FHSS/auth/null      yön bulma/tanıma
        │                   │                   │
  ── SAVUNMACI İÇİN ──  ── SAVUNMACININ ──   ── SAVUNMACININ ──
     TEHDİT (bunu          KALKANI (bunu        GÖZÜ (bunu
     tanı ve tespit et)    inşa et)             işlet)
```

Bu kitabın savunma duruşu şu hizalamayla özetlenir: Bu bölüm EA'yı yalnızca tanımak için anlatır (icra için değil); asıl yatırımını EP (dayanıklılık inşası) ve ES (tespit/izleme) üzerine yapar. Savunmacının üçgeni budur — saldırıyı tanı, kalkanı kur, gözü aç.

### Savunma açısından üç fonksiyonun birbirini beslemesi

ES olmadan EP kördür: neye karşı koruduğunu bilmeden dayanıklılık tasarlayamazsın; spektrumu izlemeden bir karıştırma olayının başladığını fark edemezsin. EP olmadan ES yetersizdir: tehdidi görürsün ama dayanağın yoksa savrulursun. Bu yüzden olgun bir savunma mimarisi her zaman ikisini birlikte kurar: bir tarafta sürekli spektrum izleme (ES), diğer tarafta yayılı spektrum + kimlik doğrulama + yedeklilik (EP). EA tarafı ise savunmacı için bir yetenek değil, bir tehdit kataloğudur — onu yalnızca modellemek ve tespit etmek için inceleriz.

> Mühendislik sezgisi: "Tehdidi tanı, kalkanı kur, gözü aç." Bu üçlü, bölümün geri kalanının iskeletidir. Her tehdit başlığında önce fizik (nasıl çalışır), sonra tespit (ES tarafı), sonra savunma (EP tarafı), en sonunda yasal sınır gelir.

Çapraz referans: ES'in ölçüm tabanı Bölüm 3'teki yön bulma (DF) ve Bölüm 7'deki ayıklama/sınıflandırmadır; EP'nin yayılı spektrum temeli Bölüm 1 (modülasyon) ve Bölüm 7'nin (LPI/LPD) konusudur.

---

<a id="2"></a>
## 2. Karıştırma Fiziği: Sinyal-Karıştırma Oranı (J/S) ve Neden Zayıf Sinyaller Kolay Bastırılır

Karıştırmanın (jamming) ne olduğunu ve neden işe yaradığını anlamak için tek bir büyüklüğü kavramak yeterlidir: alıcının girişinde istenen sinyalin gücüne karşı karıştırma gücünün oranı. Bu oran sinyal-karıştırma oranı (J/S, jamming-to-signal ratio) olarak adlandırılır ve karıştırmanın tüm fiziği bu orana indirgenebilir. Uyarı: Burada amaç bir karıştırıcının nasıl yapıldığını değil, neden çalıştığını ve neden bazı sinyallerin kırılgan olduğunu anlamaktır; bu kavrayış doğrudan savunma tasarımına döner.

Bir alıcı, bir sinyali ancak istenen sinyalin gücü gürültü ve girişim toplamına göre yeterince yüksekse çözebilir. Karıştırma, alıcının frekansına yapay enerji basarak bu dengeyi bozmaktan ibarettir. İstenen sinyal gücü S, karıştırma gücü J ise, alıcıdaki etkin oran:

```
        J        P_j · G_jr · (λ / 4π R_j)^2 · L_j
   ───────  =  ─────────────────────────────────────
        S        P_s · G_sr · (λ / 4π R_s)^2 · L_s
```

Buradaki büyüklükler kavramsaldır (formül, savunma için "neyin önemli olduğunu" gösterir, bir cihaz tasarımı değildir):

| Sembol | Anlamı | Savunma için çıkarım |
|---|---|---|
| P_j , P_s | Karıştırıcı ve istenen vericinin yayın gücü | Saldırgan gücü artırabilir; savunan istenen sinyalin gücünü/kazancını yönetir |
| G_jr , G_sr | Alıcı anteninin karıştırıcıya ve istenen kaynağa yönündeki kazancı | Yönlü anten + null-steering, karıştırıcı yönündeki kazancı düşürerek J'yi azaltır |
| R_j , R_s | Karıştırıcı–alıcı ve kaynak–alıcı mesafesi | Karıştırıcı uzaktaysa serbest uzay kaybı (FSPL, Bölüm 1) lehine işler |
| λ | Dalga boyu | Frekansa bağımlılık; bant seçimi savunmanın parçası |
| L_j , L_s | Yol/işleme kayıpları, yayma kazancı (processing gain) | Yayılı spektrumun "işleme kazancı" J/S'i etkin olarak düşürür — EP'nin kalbi |

Kritik sezgi şudur: Bir sinyali karıştırmanın zorluğu, o sinyalin alıcıdaki gücüyle doğru orantılıdır. Alıcıya çok güçlü ulaşan bir sinyali bastırmak için orantılı biçimde yüksek karıştırma gücü gerekir; buna karşılık alıcıya zayıf ulaşan bir sinyal, küçük bir karıştırma gücüyle bile gürültüye gömülebilir. Bu yüzden uzaydan gelen, alıcıya son derece zayıf ulaşan uydu sinyalleri (özellikle GNSS) yapısal olarak en kırılgan sinyallerdir.

```
   Sinyal gücü ekseni (alıcı girişinde)

   güçlü ◄───────────────────────────────────────► zayıf
     │                                                 │
  yer üstü güçlü verici                         uydu (GNSS) sinyali
  (yakın baz istasyonu)                         (~ gürültü tabanı altı)
     │                                                 │
  karıştırmak için                              karıştırmak için
  YÜKSEK J gerekir                              ÇOK DÜŞÜK J yeter
     │                                                 │
  görece dayanıklı                              yapısal olarak kırılgan
```

### GNSS'in kırılganlığının sayısal sezgisi

GNSS sinyalleri yeryüzüne yaklaşık gürültü tabanının altında bir güç yoğunluğuyla ulaşır; alıcı, bunu ancak yayma kodunun (DSSS) işleme kazancı sayesinde geri kazanır (Bölüm 7'deki LPI/LPD mantığının tersten okunması). Bu, iki yönlü bir gerçektir:

- Savunma yönü: İşleme kazancı, alıcıya bir miktar karıştırma bağışıklığı verir; yayma ne kadar genişse, dar bantlı bir karıştırmaya karşı kazanç o kadar yüksektir.
- Kırılganlık yönü: Sinyal başlangıçta zaten gürültü tabanının altında olduğundan, alıcıya ulaşan görece küçük bir yapay enerji bile işleme kazancının sağladığı marjı tüketmeye yetebilir.

Bu denge, neden GNSS karıştırma/sahteleme olaylarının (genellikle havalimanı çevrelerinde, çatışma bölgelerinde) bu kadar tekrarlandığını ve neden kritik altyapının asla tek başına GNSS'e bağlı olmaması gerektiğini açıklar (bkz. Bölüm 12).

> Mühendislik sezgisi: Karıştırmayı yenmenin üç fiziksel kolu vardır ve hepsi J/S formülünde gizlidir — (1) istenen sinyalin alıcıdaki gücünü/işleme kazancını artır (yayma, güç kontrolü), (2) karıştırıcı yönündeki anten kazancını düşür (null-steering, yönlülük), (3) karıştırıcı–alıcı mesafesini/engelini lehine kullan (konumlandırma, ekranlama). Bu üç kol, Bölüm 6'daki savunma katmanlarının tam karşılığıdır.

Yasal uyarı: J/S oranını "yükseltmek" için yapay enerji yaymak, yani karıştırma, lisanssız yayın ve kasıtlı girişim olarak hemen her ülkede ağır suçtur ve can güvenliği sistemlerini etkilediğinde ayrıca ağırlaştırıcı niteliklidir. Bu başlık yalnızca savunma tasarımı ve tespit için fiziği açıklar.

---

<a id="3"></a>
## 3. Karıştırıcı Türleri (Kavramsal Sınıflandırma — Yapım Yok)

Savunmacının bir karıştırma olayını spektrumda tanıması için, karıştırmanın spektral ve zamansal "imza ailelerini" bilmesi gerekir; çünkü farklı karıştırma stratejileri waterfall'da farklı görünür ve farklı tespit/savunma yaklaşımı gerektirir. Aşağıdaki sınıflandırma yalnızca kavramsaldır — hiçbir devre, komponent, güç kademesi veya yapım adımı içermez; amaç, savunmacının "ekranda ne görüyorum" sorusuna isim verebilmesidir.

| Karıştırma sınıfı (kavram) | Spektral/zamansal imzası | Savunmacı için tespit ipucu | En etkili savunma yönü |
|---|---|---|---|
| Barrage (geniş bantlı set) | Geniş bir bandı sürekli gürültüyle doldurur; waterfall'da geniş, kalıcı yükselen tabaka | Geniş bantta eşzamanlı gürültü tabanı sıçraması | Yayılı spektrum + yedek bant/sistem |
| Spot (nokta) | Tek bir dar frekansa yoğunlaşmış enerji | Tek kanalda keskin, kalıcı güç adası | Frekans atlama (FHSS) — saldırgan tek noktayı tutamaz |
| Sweep (taramalı) | Enerjiyi banttan banda kaydırır; waterfall'da diyagonal/süpüren çizgi | Periyodik olarak gezen dar tepe; zaman-frekansta eğik iz | FHSS + dar bant filtreleme; tarama periyoduyla senkron kaçınma (sistem tasarımı) |
| Protokol-bilinçli | Hedef protokolün belirli alanlarını/zamanlamasını hedefler; enerji düşük ama "doğru anda" | Yalnızca belirli protokol olaylarında bozulma; düşük genel gürültü ama yüksek hata oranı | Kimlik doğrulama + bütünlük + zamanlama çeşitlemesi |
| Reaktif (dinle-sonra-vur) | Yalnızca hedef yayın algılanınca tetiklenir; aksi halde sessiz | Yayın yokken temiz, yayın anında ani bozulma; "sebepli" korelasyon | Kısa burst + atlama + düşük tespit profili (LPD); öngörülemez zamanlama |

Not: Bu tablo bir "menü" değildir; savunmacının teşhis cetvelidir. Pratikte gözlemlenen birçok olay melez veya kalitesizdir (örneğin ucuz bir yasadışı cihaz hem geniş bandı kirletip hem belirli bir banda yoğunlaşabilir). Önemli olan, waterfall'daki imzadan hangi aileye yakın olduğunu çıkarıp doğru savunmayı seçebilmektir.

```
   Waterfall'da kavramsal imza karşılaştırması (zaman ↓, frekans →)

   BARRAGE                SPOT                 SWEEP
   f→                     f→                   f→
   ████████████  ↓t       ░░░██░░░░░░  ↓t       █░░░░░░░░░░  ↓t
   ████████████           ░░░██░░░░░░           ░█░░░░░░░░░
   ████████████           ░░░██░░░░░░           ░░█░░░░░░░░
   ████████████           ░░░██░░░░░░           ░░░█░░░░░░░
   (geniş, dolu)          (dar, kalıcı)        (gezen diyagonal)

   REAKTİF (zaman ekseni kritik)
   t↓   olay yok → temiz ......
        hedef yayını başladı → ██ (ani)
        yayın bitti → temiz ......
```

Savunmacı için pratik ayrım: Barrage ve spot kalıcıdır ve gürültü tabanı izlemeyle (Bölüm 5) kolay yakalanır; sweep, zaman-frekans (waterfall) izini gerektirir; reaktif ve protokol-bilinçli olanlar en sinsidir çünkü ortalama güç düşüktür ve yalnızca olayla korelasyonda görünür — bunları yakalamak, ham gürültü tabanı değil, "hizmet bozulmasıyla RF anomalisinin zaman korelasyonu" gerektirir.

Yasal uyarı: Bu sınıflandırma yalnızca tehdidi tanıma içindir. Bu sınıfların hiçbirini üretmek/yaymak yasal değildir; metin kasıtlı olarak hiçbir yapım, parametre veya cihaz ayrıntısı vermez.

---

<a id="4"></a>
## 4. Karıştırmanın Etkileri: GPS, GSM, Havacılık, Acil — Can Güvenliği Boyutu

Karıştırmanın neden yalnızca bir "teknik arıza" değil, doğrudan bir can güvenliği ve kritik altyapı tehdidi olduğunu anlamak, savunma önceliklendirmesi için zorunludur. Aşağıdaki tablo, karıştırmanın hangi sistemde hangi sonucu doğurabileceğini ve neden bunun ölümcül olabileceğini gösterir; bu, savunmacının "neyi önce korumalıyım" sorusunun cevabıdır.

| Hedef sistem | Karıştırmanın doğrudan etkisi | Neden can/altyapı kritik |
|---|---|---|
| GPS/GNSS (konum) | Konum/kilit kaybı, hatalı konum | Havacılık seyrüsefer, denizcilik, ambulans yönlendirme, hassas zaman |
| GNSS zaman | Zaman senkron kaybı | Telekom, enerji şebekesi, finans damgalama, SCADA (Bölüm 12) |
| GSM/LTE/5G (hücresel) | Çağrı/veri kesintisi, acil çağrı (112/911) engeli | Acil yardım çağrısı yapılamaması — doğrudan can riski |
| Havacılık telsizi/seyrüsefer | ATC sesi, VOR/ILS, ADS-B bozulması | Uçuş güvenliği; yaklaşma/iniş kritik fazları |
| Denizcilik (AIS/GMDSS) | Gemi konumu/acil çağrı kaybı | Çarpışma önleme, arama-kurtarma |
| Kamu güvenliği telsizi (TETRA/P25) | İlk müdahale ekiplerinin haberleşme kaybı | Olay yerinde koordinasyon çöküşü |
| Tıbbi telemetri / alarm | Hasta izleme/alarm sinyali kaybı | Hastane/ev tipi cihazlarda fark edilmeyen kriz |

Bu tablonun savunmacı için verdiği ders nettir: Karıştırma, "sinyal kalitesi" meselesi değil, insan güvenliği meselesidir. Bu yüzden hemen her ülke karıştırmayı yalnızca telekomünikasyon ihlali değil, kamu güvenliğine yönelik ağır suç sayar; havacılık veya acil çağrı etkilendiğinde cezalar ağırlaşır.

> Pratikte: Bir kurumda "açıklanamayan GPS sürüklenmesi", "belirli bir saatte tekrarlayan hücresel kesinti" veya "yalnızca belirli bir araç/cihaz çevresinde bağlantı kaybı" gibi örüntüler, bir karıştırma olayının erken işaretleri olabilir. Savunmacı bu örüntüleri arıza değil, potansiyel EA olayı olarak da değerlendirmeli ve spektrum tarafında doğrulamalıdır (Bölüm 5).

Yasal uyarı: Bu etkiler tam olarak karıştırmanın neden suç olduğunu açıklar. Kişisel "sinyal engelleyici/jammer" cihazları (araç içi GPS engelleyici, sınıf/sınav jammer'ı, kişisel hücresel engelleyici dahil) çoğu ülkede satışı, bulundurması ve kullanımı yasaktır; can güvenliği sistemlerini etkilediğinde sonuçlar ağırdır. Savunma her zaman pasif (tespit + dayanıklılık + yedeklilik) yoldan kurulur.

---

<a id="5"></a>
## 5. Karıştırma Tespiti: Spektrum İzleme, Gürültü Tabanı, Anomali ve Yön Bulma

Karıştırmaya karşı ilk savunma, onu görmektir (ES fonksiyonu). Tespit, dört tamamlayıcı katmandan oluşur ve hepsi pasiftir — yalnızca dinler, hiçbir şey yaymaz. Bu, savunmacının yasal olarak ve sürekli işletebileceği bir yetenektir.

### 5.1 Gürültü tabanı (noise floor) izleme

En temel ve en güçlü tespit yöntemi, ilgili bantlardaki gürültü tabanını sürekli izlemektir. Normal koşullarda bir bandın gürültü tabanı belirli bir aralıkta dalgalanır; bir karıştırma olayı bu tabanı belirgin ve kalıcı biçimde yükseltir.

```
   Güç (dBm)
     ▲
     │                       ┌──────────────┐  ← karıştırma başladı
     │                       │   yükselmiş   │     (gürültü tabanı
     │                       │ gürültü tabanı│      ani ve kalıcı
     │   ~~~~~~~~~~~~~~~~~~~~~┘              └~~~~  yükseldi)
     │   normal taban (dalgalı ama sınırlı)
     └───────────────────────────────────────────► zaman
                            ▲
                      eşik aşımı = alarm
```

Yöntem: Bir referans (baseline) profil oluştur — bandın saatlik/günlük normal gürültü tabanı dağılımı. Sonra canlı ölçümü bu profile karşı kıyasla; istatistiksel olarak anlamlı (örneğin baseline + birkaç standart sapma) ve süreklilik gösteren bir yükseliş, karıştırma adayıdır. Tek seferlik kısa sıçramalar (geçici girişim) ile kalıcı yükselişi ayırmak için zaman penceresi (örneğin "N saniyeden uzun süren yükseliş") kullanılır.

### 5.2 Spektral anomali ve şekil analizi

Gürültü tabanının yalnızca seviyesi değil, şekli de bilgi taşır (Bölüm 3'ten farklı olarak, yön bulmadan).

| Anomali türü | Spektral görünüm | Olası yorum |
|---|---|---|
| Geniş bantlı kalkma | Tüm bant eşzamanlı yükselir | Barrage benzeri |
| Dar bant kalıcı tepe (beklenmedik yerde) | Tahsiste sinyal olmaması gereken yerde sürekli enerji | Spot benzeri |
| Gezen tepe | Tepe zamanla frekansta kayıyor | Sweep benzeri |
| Olayla korelasyonlu bozulma | RF anomalisi tam da hizmet bozulması anında | Reaktif/protokol-bilinçli |

### 5.3 Hizmet/sinyal kalitesi metrikleriyle çapraz doğrulama

Saf RF ölçümü her zaman tek başına yeterli değildir; karıştırmayı kesinleştirmek için sistem-tarafı metriklerle çapraz doğrulama güçlüdür: GNSS alıcısında C/N0 (taşıyıcı/gürültü yoğunluğu) düşüşü, izlenen uydu sayısının ani azalması, hücresel modemde sinyal kalitesi (RSRP/RSRQ/SINR) çöküşü, hata oranı (BER/PER) artışı. RF tarafındaki gürültü tabanı yükselişi ile sistem tarafındaki kalite çöküşünün zaman içinde örtüşmesi, karıştırma teşhisini güçlü biçimde destekler.

### 5.4 Yön bulma (DF) ile kaynağı yer tespiti

Karıştırma doğrulandıktan sonra kaynağı bulmak, yetkili makamların müdahalesi (yasal) için kritiktir. Bu, Bölüm 3'teki yön bulma tekniklerinin doğrudan uygulamasıdır: birden çok konumdan kestirim açısı (bearing) alınır ve açıların kesişimiyle kaynak konumu üçgenlenir.

```
   Yön bulma ile karıştırıcı yerini üçgenleme (kavram)

   Sensör 1 ──────►        ◄────── Sensör 2
            \  açı θ1      açı θ2  /
             \                   /
              \                 /
               \      ✕        /     ✕ = açıların kesiştiği nokta
                \   (kaynak)  /          = karıştırıcı tahmini konumu
                 \           /
                  Sensör 3 ──► açı θ3 (doğrulama)
```

Önemli: Yön bulma pasif bir ölçümdür (yalnızca dinler) ve yasaldır; ancak bulunan kaynağa müdahale (cihazı etkisiz kılma, karşı yayın) yetki ister ve genellikle yalnızca yetkili kurumların görevidir. Savunmacının rolü tespit + raporlamadır; icra değil.

> Mühendislik sezgisi: Karıştırma tespiti üç katmanın füzyonudur — (1) gürültü tabanı yükseldi mi (seviye), (2) imza hangi aileye benziyor (şekil/zaman), (3) sistem kalitesi eşzamanlı çöktü mü (çapraz doğrulama). Üçü birden işaret ediyorsa teşhis güçlüdür; yalnızca biri varsa şüphe düzeyindedir. Kaynağı bulmak için dördüncü katman yön bulmadır.

Çapraz referans: Gürültü tabanı ve waterfall okuma becerisi Bölüm 2 (SDR) ve Bölüm 7 (ayıklama); yön bulma matematiği Bölüm 3.

---

<a id="6"></a>
## 6. Karıştırmaya Karşı Savunma: Yayılı Spektrum, Atlama, Null-Steering, Güç Kontrolü, Yedeklilik

Tespit (ES) tehdidi görür; dayanıklılık (EP) onu yener. Karıştırmaya karşı savunma katmanlıdır ve her katman J/S formülünün (Bölüm 2) farklı bir terimine saldırır. Hiçbir tek katman yeterli değildir; güç, katmanların birlikteliğindedir.

```
   KARIŞTIRMAYA KARŞI SAVUNMA KATMANLARI (dıştan içe)

   ┌──────────────────────────────────────────────────┐
   │ 5. YEDEKLİLİK / DEĞİŞKEN MOD                       │
   │    (çoklu bant, alternatif sistem, GNSS+INS)       │
   │  ┌────────────────────────────────────────────┐   │
   │  │ 4. ANTEN / UZAMSAL (null-steering, yönlülük)│   │
   │  │  ┌──────────────────────────────────────┐   │   │
   │  │  │ 3. GÜÇ KONTROLÜ / LİNK BÜTÇESİ         │   │   │
   │  │  │  ┌────────────────────────────────┐   │   │   │
   │  │  │  │ 2. FREKANS ATLAMA (FHSS)        │   │   │   │
   │  │  │  │  ┌──────────────────────────┐   │   │   │   │
   │  │  │  │  │ 1. YAYILI SPEKTRUM/işleme │   │   │   │   │
   │  │  │  │  │    kazancı (DSSS)         │   │   │   │   │
   │  │  │  │  └──────────────────────────┘   │   │   │   │
   │  │  │  └────────────────────────────────┘   │   │   │
   │  │  └──────────────────────────────────────┘   │   │
   │  └────────────────────────────────────────────┘   │
   └──────────────────────────────────────────────────┘
```

| Katman | Mekanizma | J/S formülünde saldırdığı terim | Pratik karşılığı |
|---|---|---|---|
| Yayılı spektrum (DSSS) | Enerjiyi geniş banda yayma; alıcıda işleme kazancı | L (işleme kazancı) — etkin J/S'i düşürür | GNSS, bazı askeri/IoT bağları |
| Frekans atlama (FHSS) | Taşıyıcıyı sözde rastgele dizide değiştirme | Saldırganın tek bandı tutamaması | Bluetooth (uyarlanır), bazı telsiz/dron bağları |
| Güç kontrolü / link bütçesi | İstenen sinyalin alıcıdaki gücünü/marjını yönetme | S'i yükseltme, R_s'i kısaltma | Hücresel güç kontrolü, kısa link mesafesi |
| Anten / uzamsal | Null-steering, yönlü anten, dizi işleme | G_jr'yi (karıştırıcı yönü kazancı) düşürme | CRPA tipi anti-jam GNSS antenleri (kavram), yönlü baz anten |
| Yedeklilik / değişken mod | Çoklu bant, alternatif sistem, eylemsizlik (INS) yedeği | Sistemin tek bir RF bağına bağımlılığını kırma | GNSS+INS, çoklu GNSS takımyıldızı, kablolu yedek |

### Katmanların açıklaması (savunma okuryazarlığı)

Yayılı spektrum ve işleme kazancı: Enerjiyi geniş banda yaymak, dar bantlı bir karıştırıcının etkisini alıcıdaki "despreading" işlemiyle bastırır; bu, J/S'i etkin olarak işleme kazancı kadar düşürür. Bölüm 7'de LPI/LPD'yi "tespiti zorlaştırma" olarak okumuştuk; burada aynı fizik "karıştırmaya direnç" olarak çalışır.

Frekans atlama: Saldırgan tek bir frekansa yoğunlaşsa bile, sistem o frekansta yalnızca kısa süre kalır; atlama dizisi gizli ve hızlıysa, saldırganın "doğru anda doğru yerde" olması zorlaşır. Reaktif karıştırmaya karşı en etkili yapısal savunmalardan biridir.

Anten/uzamsal (null-steering): Çok elemanlı bir anten dizisi, karıştırıcının geldiği yöne bir "boşluk" (null) yerleştirerek o yöndeki kazancı (G_jr) bastırır; istenen sinyali korurken karıştırıcıyı uzamsal olarak söndürür. GNSS dünyasında bu yaklaşımın anti-jam anten ailesi (kavramsal olarak CRPA) vardır.

Yedeklilik: En sağlam savunma, sisteme tek bir RF bağına bağımlı olmamayı öğretmektir. GNSS karıştırılırsa eylemsizlik seyrüseferine (INS) geçen bir sistem, konum çözümünü bir süre koruyabilir; çoklu takımyıldız (GPS+Galileo+GLONASS+BeiDou) ve çoklu bant, tek noktayı vurmayı zorlaştırır. Kritik altyapıda zaman için GNSS+yerel atomik saat (holdover) bu felsefenin örneğidir (Bölüm 12).

> Mühendislik sezgisi: Karıştırmaya karşı savunma, J/S formülünün her terimine ayrı bir kontrol koymaktır. Tek bir "sihirli kalkan" yoktur; yayma (L), atlama (zaman boyutu), güç (S), anten (G_jr) ve yedeklilik (bağımlılığı kır) birlikte, saldırganın işini katlanarak zorlaştırır.

Yasal/etik not: Bu savunmaların hepsi savunan tarafın kendi sistemine uyguladığı pasif/tasarımsal önlemlerdir. Hiçbiri karşı-yayın (karşı-karıştırma) değildir; karşı yayın da bir yayındır ve yetki ister. Savunmacının alanı dayanıklı tasarım + tespittir.

---

<a id="7"></a>
## 7. Drone/İHA RF Güvenliği: Tipik Bağlar ve Neden Savunmasız

İnsansız hava araçları (İHA, drone), savunmacı için iki ayrı şapkayla incelenir: (a) korunması gereken bir varlık (kendi/kurumun dronu hedef olabilir) ve (b) bir tehdit (yetkisiz bir drone, mahremiyet/güvenlik/kritik altyapı riski oluşturabilir). Her iki şapka da aynı RF gerçeğine dayanır: bir drone, birden çok radyo bağına bağımlıdır ve her bağ bir saldırı yüzeyidir. Bu başlık, bu bağların nasıl çalıştığını ve neden bazılarının zayıf olduğunu — saldırı icrası vermeden — anlatır.

### Tipik drone RF bağları

```
                    ┌──────────────────────┐
   GNSS uyduları    │       DRONE          │
        ░░░ ──────► │  (konum: GPS/GNSS)   │
   (zayıf, alttan)  │                      │
                    │  ┌────────────────┐  │
   Kumanda (RC)     │  │ uçuş kontrol    │  │   Video downlink
   2.4 / 5.8 GHz ◄──┼─►│ + telemetri     │──┼──► operatör/gözlük
   (kontrol bağı)   │  └────────────────┘  │   (2.4 / 5.8 GHz, analog/dijital)
                    │                      │
                    │  Remote ID yayını ───┼──► çevreye kimlik/konum
                    └──────────────────────┘       (yasal şeffaflık)
```

| Bağ | Tipik frekans | İşlevi | Zayıflık kaynağı (bazı modellerde) |
|---|---|---|---|
| Kumanda (RC kontrol) | 2.4 GHz / 5.8 GHz (ISM) | Pilot komutları | Bazı eski/ucuz modellerde şifresiz veya zayıf kimlik doğrulama; bağ kaybında failsafe davranışı öngörülebilir |
| Video downlink | 2.4 / 5.8 GHz | Kameradan operatöre görüntü | Analog video çoğunlukla şifresiz; dinlenmeye/bozulmaya açık |
| GNSS | L bandı (GPS/GNSS) | Konum/dönüş-eve (RTH) | Bölüm 2/4/12'deki GNSS kırılganlığı: zayıf sinyal, karıştırma/sahtelemeye açık |
| Telemetri | Değişken (modeme bağlı) | Durum/konum geri bildirimi | Şifreleme/kimlik doğrulama modele göre değişir |
| Remote ID | Genellikle yerel kablosuz (Wi‑Fi/BT yayını) | Yasal kimlik/konum şeffaflığı | Tasarımı gereği açık (savunmacı için tespit fırsatı, Bölüm 8) |

### Neden bazı drone bağları savunmasız

Drone RF güvenliğinin temel dersi şudur: güvenlik, modelden modele dramatik biçimde değişir. Üst seviye ticari ve kurumsal sistemler şifreli, kimlik doğrulamalı, atlamalı bağlar ve sağlam GNSS alıcıları kullanırken; ucuz/eski/hobi sınıfı bazı cihazlar şifresiz kontrol, açık analog video ve korumasız GNSS bağımlılığı taşıyabilir. Savunmasızlığın üç yapısal kaynağı:

1. GNSS bağımlılığı: Otonom uçuş, dönüş-eve ve konum tutma GNSS'e dayanır; GNSS yapısal olarak zayıf sinyalli olduğundan (Bölüm 2), bu, dronun en kırılgan duyusudur.
2. Kimlik doğrulama zayıflığı (bazı modeller): Kontrol bağı kriptografik olarak kimliklendirilmiş değilse, "kimden geldiği doğrulanmamış" komutlara ilkesel olarak açıktır.
3. Açık downlink: Şifresiz video/telemetri, bir gözlemcinin içeriği görmesine ve meta-veriden çıkarım yapmasına izin verir (Bölüm 7'deki trafik analizi mantığı).

> Mühendislik sezgisi: Bir drone "uçan, çok bağlı bir IoT cihazıdır". Her radyo bağı (kontrol, video, GNSS, telemetri, Remote ID) ayrı bir yüzeydir ve en zayıf bağ, tüm sistemin güvenlik seviyesini belirler. Savunmacı, kendi dronunu seçerken/işletirken bu bağların her birinin kriptografik durumunu sorar; tehdit olarak değerlendirirken ise bu bağların tespit edilebilir imzalarını arar (Bölüm 8).

Yasal uyarı: Bu başlık yalnızca bağların nasıl çalıştığını ve neden zayıf olabileceğini açıklar. Bir drona (kendininki dışında) yetkisiz erişim, komut enjeksiyonu, ele geçirme veya GNSS sahteleme suçtur ve hava güvenliğini tehdit eder. Kendi dronunu yasal olarak analiz edebilirsin (kendi kontrol bağının imzasını gözlemlemek gibi); başkasının sistemine müdahale edemezsin.

---

<a id="8"></a>
## 8. Drone Saldırı Sınıfları (Yalnızca Kavram) ve Counter-UAS Savunması

Bu başlık iki yarıdan oluşur: önce savunmacının tanıması gereken saldırı sınıflarının kavramsal taksonomisi (icra adımı YOK), sonra bu tehditlere karşı meşru karşı-İHA (Counter-UAS, C‑UAS) savunması ve — en önemlisi savunmacı için — drone TESPİTİ.

### 8.1 Saldırı sınıfları — yalnızca kavram (icra reçetesi yok)

| Saldırı sınıfı (kavram) | Temel fikir (yüksek seviye) | Neden mümkün olabilir | Bunun karşısındaki dayanıklılık |
|---|---|---|---|
| Karıştırma → failsafe tetikleme | Kontrol/GNSS bağını bozarak dronu failsafe moduna (iniş/RTH/havada kalma) zorlama | Bağ kaybında davranış öngörülebilirse | Şifreli+atlamalı bağ, sağlam GNSS, öngörülemez failsafe değil—güvenli failsafe |
| GNSS sahteleme → konum yanıltma | Sahte GNSS ile dronun "yanlış yerde olduğunu sanması" | GNSS kimliklendirilmemiş ve zayıf sinyalli | Çok-takımyıldız, INS füzyonu, sahtecilik tespiti (tutarlılık kontrolü) |
| Protokol replay | Daha önce gözlemlenen kontrol mesajlarını tekrar oynatma | Mesajlar nonce/sayaç/zaman damgasıyla korunmuyorsa | Rolling/nonce tabanlı kimlik doğrulama (bkz. Bölüm 9) |

Çok önemli sınır: Yukarıdaki üç satır yalnızca "bu tehdit kategorileri vardır" demek içindir. Bu metin bu saldırıların hiçbirinin nasıl yapılacağını (frekans, zamanlama, mesaj formatı, araç, adım) vermez ve vermeyecektir. Amaç, savunmacının bir olayı gördüğünde "bu bir GNSS-sahteleme örüntüsüne benziyor" diyebilmesi ve doğru dayanıklılığı (sağ sütun) seçmesidir.

### 8.2 Drone TESPİTİ (savunmacının asıl işi — pasif ve yasal)

Bir yetkisiz drone tehdidine karşı savunmanın ilk ve en yasal adımı tespittir. Tespit çoğunlukla pasiftir (yalnızca dinler) ve birden çok yöntemi birleştirir:

| Tespit yöntemi | Nasıl çalışır | Güçlü yanı | Sınırı |
|---|---|---|---|
| RF tespiti / imza | Drone kontrol/video bağının spektral imzasını tanıma | Pasif, menzilli, yön bulmaya açık | Şifreli/atlamalı/sessiz bağlarda zorlaşır |
| RF yön bulma (DF) | Birden çok sensörle açı alıp drone+operatör konumunu üçgenleme | Hem dronu hem pilotu konumlandırabilir | Çok sensör + temiz açı gerektirir (Bölüm 3) |
| Remote ID dinleme | Yasal kimlik/konum yayınını okuma | Standart, yasal, kimlik+konum verir | Yalnızca uyumlu/dürüst cihazlar yayınlar |
| Radar | Aktif radar ile küçük hedef tespiti | RF-sessiz dronu da görebilir | Maliyet; küçük hedef/kuş ayrımı; yetki |
| Akustik | Pervane sesini tanıma | Pasif, kısa menzil | Gürültülü ortamda zayıf |
| Elektro-optik/IR | Kamera/termal ile görsel tespit | Görsel doğrulama | Görüş hattı, hava koşulu |

Remote ID, savunmacı için özellikle değerlidir: birçok yargı bölgesinde dronların kimlik ve konum bilgisini açıkça yayınlaması yasal zorunluluktur; bu yayın pasif olarak (yasal biçimde) dinlenip uyumlu dronların kimliği ve konumu okunabilir. Bu, "şeffaflıkla güvenlik" yaklaşımının somut bir örneğidir.

```
   Katmanlı drone tespiti (sensör füzyonu)

   RF imza ──┐
   RF DF ────┤
   Remote ID ┼──► FÜZYON ──► "drone var mı? nerede? kimliği? pilotu nerede?"
   Radar ────┤        │
   Akustik ──┤        └──► doğrulanmış uyarı (tek sensör değil, korelasyon)
   EO/IR ────┘
```

### 8.3 Counter-UAS (etkisiz kılma) — yetki ve yasallık

Tespit yasal ve geniş biçimde uygulanabilirken, etkisiz kılma (mitigation) — yani dronu karıştırma, ele geçirme, fiziksel olarak durdurma — tamamen farklı bir hukuki alandır. Bu eylemler:

- Genellikle yalnızca yetkili kamu kurumlarına (kolluk, ordu, sivil havacılık otoritesi) ve belirli yasal yetkilerle tanınır.
- Karıştırma içeren C-UAS, sıradan karıştırma gibi spektrum hukukuna tabidir ve çevredeki diğer sistemleri (havacılık, hücresel) etkileyebileceğinden ek risk taşır.
- "Ele geçirme/komut alma" içeren C-UAS, bilgisayar/iletişim sistemlerine yetkisiz erişim hukukuna da girebilir.

Savunmacının (kurumsal güvenlik, kritik altyapı operatörü) rolü neredeyse her zaman tespit + raporlama + fiziksel/idari önlem (erişim kontrolü, geofence talebi, yetkili makamı bilgilendirme) ile sınırlıdır; etkisiz kılma yetkili makama bırakılır.

### 8.4 Geofence ve Remote ID — tasarımsal savunma

Drone ekosisteminde iki tasarımsal kontrol, savunmaya doğrudan katkı sağlar: Geofence (üretici tarafından hassas bölgelerde — havalimanı, hapishane, kritik tesis çevresinde — uçuşu yazılımsal kısıtlama) ve Remote ID (kimlik/konum şeffaflığı). Bunlar saldırıyı durdurmaz ama tehdidi azaltır ve tespiti kolaylaştırır; savunmacının yetkililerden talep edebileceği/uygulanmasını destekleyebileceği yapısal önlemlerdir.

> Mühendislik sezgisi: Drone tehdidine karşı savunmacının üçgeni: (1) TESPİT et (RF+DF+Remote ID füzyonu — pasif, yasal), (2) RAPORLA ve fiziksel/idari önlem al, (3) ETKİSİZ KILMAYI yetkiliye bırak. Saldırı sınıflarını yalnızca tanımak için öğren; icra etmek için değil.

Yasal uyarı: Bir drona yetkisiz müdahale (karıştırma, sahteleme, ele geçirme) ve genel olarak yetkisiz C-UAS faaliyeti suçtur ve hava güvenliğini tehlikeye atar. Tespit (özellikle Remote ID dinleme ve pasif RF gözlem) genellikle yasaldır; etkisiz kılma yetki ister. Sınırı kendi ülkenden teyit et.

---

<a id="9"></a>
## 9. Replay / Spoofing / Relay Tehditleri: Prensip ve Savunma

Bu başlık, gündelik hayata en yakın RF manipülasyon ailesini ele alır: uzaktan kumandalar, keyless araç sistemleri ve genel olarak kimlik doğrulaması zayıf RF kontrol sinyalleri. Burada da prensip + savunma verilir; kullanıcının yalnızca kendi cihazını analiz etmesi yasaldır.

### 9.1 Sabit-kod (fixed-code) replay

En basit ve en eski zayıflık: bir uzaktan kumanda her seferinde aynı sabit kodu gönderiyorsa, bu kod bir kez gözlemlendiğinde tekrar oynatılarak (replay) aynı etki üretilebilir. Eski garaj kapıları, bazı basit RF prizler ve oyuncaklar bu sınıftadır.

```
   SABİT KOD (kırılgan)            ROLLING KOD (korumalı)
   her basışta:                   her basışta:
   ┌──────────┐                   ┌──────────┐
   │ KOD: A7F3│  ──► aynı         │ KOD: f(K, sayaç) │ ──► her seferinde
   │ KOD: A7F3│      her          │ sayaç++          │      farklı
   │ KOD: A7F3│      seferinde    │ KOD: değişir     │      (replay işe
   └──────────┘                   └──────────┘            yaramaz)
   replay = aynı kodu tekrarla    eski kod "geçti" sayılır, reddedilir
```

### 9.2 Rolling-code (atlamalı kod) koruması

Modern uzaktan kumandalar (araç anahtarları, çağdaş garaj kapıları) rolling-code (hopping code) kullanır: her basışta, paylaşılan bir gizli anahtar ve artan bir sayaçtan türetilen farklı bir kod gönderilir. Alıcı, beklenen sayaç penceresindeki kodları kabul eder; bir kez kullanılan/eskimiş kod reddedilir. Bu, basit replay'i etkisiz kılar çünkü gözlemlenen kod bir sonraki basışta zaten geçersizdir.

Not: Rolling-code kavramsal olarak replay'i çözer; ancak gerçek dünyada bazı eski/zayıf uygulamaların tarihsel zafiyetleri raporlanmıştır (zayıf rastgelelik, sayaç yönetimi hataları, vb.). Buradan çıkan tasarım dersi (Bölüm 10): koruma "rolling" etiketinde değil, kriptografik kalitede ve doğru sayaç/nonce yönetimindedir.

### 9.3 RF relay (menzil uzatma) — keyless araç örneği

Relay saldırısı kavramı, kodu kırmaya çalışmaz; bunun yerine, meşru sinyali olduğu gibi uzatır. Keyless giriş sistemlerinde araç ve anahtar birbirine yakınken konuşur; relay kavramı, bu konuşmayı iki nokta arasında köprüleyerek aracın anahtarı "yakında sanmasını" sağlar — kod hiç kırılmaz, yalnızca menzil yapay olarak uzatılır.

```
   RELAY KAVRAMI (yüksek seviye — icra yok)

   [ARAÇ] ◄···· yakınlık beklentisi ····► [ANAHTAR gerçekte uzakta]
       │                                        │
       │   ── köprü/relay ile sinyal taşınır ──  │
       └────────────────────────────────────────┘
   Araç anahtarı "yakında" sanır; kod kırılmaz, mesafe yanıltılır.

   SAVUNMA: mesafeyi GERÇEKTEN ölç (UWB time-of-flight) →
   "yakın gibi görünmek" yetmez, fiziksel olarak yakın olmak gerekir.
```

### 9.4 Savunma matrisi

| Tehdit | Neden mümkün | Savunma (tasarım) | Savunma (kullanıcı) |
|---|---|---|---|
| Sabit-kod replay | Kod her seferinde aynı | Rolling-code / nonce / sayaç | Eski sabit-kod cihazları değiştir |
| Zayıf rolling-code | Zayıf rastgelelik/sayaç yönetimi | Güçlü kripto + doğru sayaç penceresi | Üretici güncellemeleri, bilinen zafiyetli ürünlerden kaçınma |
| RF relay (keyless) | Yakınlık varsayımı, mesafe ölçülmüyor | UWB time-of-flight mesafe ölçümü, hareket sensörü | Anahtarı Faraday kesesinde/kutuda tutma, "hareketsizken yayma" özelliğini açma |
| Genel spoofing | Kimlik doğrulanmıyor | Kriptografik kimlik doğrulama + bütünlük | Kimliklendirilmiş sistemleri tercih etme |

Kritik savunma fikri (relay'e karşı): Relay'i yenmenin yapısal yolu, "yakın görünmek" ile "fiziksel olarak yakın olmak" arasındaki farkı ölçmektir. UWB (Ultra-Wideband) tabanlı mesafe ölçümü (time-of-flight), sinyalin gidiş-dönüş süresinden gerçek mesafeyi hesaplar; relay köprüsü ek gecikme eklediği için "yakın" iddiası fiziksel olarak çürür. Bu yüzden modern güvenli keyless sistemleri UWB mesafe doğrulamasına yönelmiştir.

> Mühendislik sezgisi: Bu üç tehdidin ortak dersi tek cümlede toplanır — "tekrar edilebilen ve doğrulanmayan şey güvenli değildir". Sabit kod tekrar edilebilir (replay); kimlik doğrulanmayan komut sahtelenebilir (spoof); ölçülmeyen mesafe yanıltılabilir (relay). Savunma sırasıyla: değişen kod (nonce/sayaç), kriptografik kimlik, ve gerçek mesafe ölçümü.

Yasal uyarı: Kendi cihazının (kendi garaj kumandanın, kendi aracının anahtarının) RF imzasını yasal olarak gözlemleyip "sabit mi rolling mi" diye analiz edebilirsin (bkz. Bölüm 14 alıştırması). Başkasının aracına/kapısına yetkisiz erişim, kod yakalama veya relay ile giriş hırsızlık ve bilişim suçudur. Bu metin hiçbir cihaza yetkisiz erişim adımı vermez.

---

<a id="10"></a>
## 10. Sinyal Dayanıklılığı: Hangi Sinyaller Manipülasyona Açık/Kapalı — Tasarım Dersleri

Önceki başlıkların hepsi tek bir tasarım sorusuna yakınsar: bir sinyali manipülasyona (replay, spoof, jamming-aldatma) açık ya da kapalı yapan nedir? Aşağıdaki matris, savunmacının ve sistem tasarımcısının "bu sinyal güvenilir mi" sorusuna yapısal cevabıdır. Bu, bölümün belki de en kalıcı dersidir.

| Özellik | Yoksa (kırılgan) | Varsa (dayanıklı) | Hangi tehdidi engeller |
|---|---|---|---|
| Kimlik doğrulama (authentication) | Komut kimden geldiği belirsiz | Kriptografik olarak kimliklendirilmiş | Spoofing, sahte komut |
| Tazelik (nonce / sayaç / zaman damgası) | Eski mesaj tekrar geçerli | Her mesaj benzersiz, eskisi reddedilir | Replay |
| Bütünlük (MAC/imza) | Mesaj fark edilmeden değiştirilebilir | Değişiklik tespit edilir | Tampering, manipülasyon |
| Şifreleme (confidentiality) | İçerik okunabilir | İçerik gizli | Dinleme, içerik sızıntısı |
| Yayılı spektrum / atlama (FHSS/DSSS) | Dar bantta kolay bastırılır/yakalanır | İşleme kazancı + zaman çeşitlemesi | Karıştırma, tespit/ayıklama |
| Mesafe doğrulama (UWB ToF) | "Yakın görünmek" yeterli | Gerçek mesafe ölçülür | Relay |
| Kaynak çeşitliliği / yedeklilik | Tek bağa bağımlı | Çoklu bant/sistem/sensör | Tek noktadan karıştırma |

```
   DAYANIKLILIK MERDİVENİ (alttan yukarı güçlenir)

   ┌─────────────────────────────────────────────┐
   │ YEDEKLİLİK (çoklu bağ, INS, çoklu takımyıldız)│ ← tek nokta vurmayı kır
   ├─────────────────────────────────────────────┤
   │ YAYMA/ATLAMA (FHSS/DSSS)                      │ ← karıştırmaya direnç
   ├─────────────────────────────────────────────┤
   │ MESAFE DOĞRULAMA (UWB ToF)                    │ ← relay'i kır
   ├─────────────────────────────────────────────┤
   │ TAZELİK (nonce/sayaç/zaman)                   │ ← replay'i kır
   ├─────────────────────────────────────────────┤
   │ BÜTÜNLÜK (MAC/imza)                           │ ← manipülasyonu yakala
   ├─────────────────────────────────────────────┤
   │ KİMLİK DOĞRULAMA (auth)                       │ ← spoofing'i kır
   └─────────────────────────────────────────────┘
        (en altta yoksa, üsttekiler tek başına yetmez)
```

Tasarım dersleri (savunmacı/tasarımcı için):

1. Şifreleme tek başına yetmez: Şifreleme içeriği gizler ama kimlik doğrulama olmadan sahte komutu, tazelik olmadan replay'i engellemez. (Bölüm 7'deki "şifreleme içeriği korur, meta-veri sızar" dersinin kontrol tarafındaki kardeşi.)
2. "Rolling" etiketi koruma garantisi değildir: Koruma, kriptografik kalitede ve doğru sayaç/nonce yönetimindedir; etikette değil.
3. Yapısal kırılganlık seçimle değil tasarımla çözülür: GNSS'in zayıf sinyali bir hata değil, fiziğin sonucudur; çözüm sinyali "güçlendirmek" değil, yedeklilik ve tutarlılık kontrolü eklemektir.
4. En zayıf bağ seviyeyi belirler: Çok bağlı bir sistemde (drone gibi) güvenlik, en zayıf radyo bağı kadardır; her bağ ayrı denetlenmelidir.

> Mühendislik sezgisi: Bir sinyalin güvenliğini değerlendirirken sırayla sor — kimliklendirilmiş mi (auth)? tazelik var mı (nonce)? bütünlük korunuyor mu (MAC)? mesafe gerçekten ölçülüyor mu (relay'e karşı)? tek bağa mı bağımlı (yedeklilik)? Bu beş soru, neredeyse tüm RF manipülasyon ailesini kapsayan bir denetim listesidir.

---

<a id="11"></a>
## 11. TEMPEST ve Yan-Kanal Yayılım: Savunma Farkındalığı

RF tehdidi yalnızca kasıtlı bir vericiden gelmez; cihazların istem dışı yaydığı elektromanyetik enerji de bir sızıntı kanalıdır. Bu, TEMPEST ve yan-kanal (side-channel) yayılım başlığıdır ve Bölüm 6'da derinlemesine işlenmiştir; burada RF tehdit manzarası bütünlüğü için kısa bir savunma farkındalığı olarak ele alınır.

Temel fikir: Çalışan elektronik (ekran, kablo, işlemci, klavye) işini yaparken, taşıdığı bilgiyle ilişkili istem dışı elektromanyetik yayınlar üretebilir. Yeterince hassas ve yakın bir alıcı, bazı durumlarda bu yayınlardan bilgi (ekran görüntüsü izleri, tuş vuruşu zamanlaması gibi) çıkarabilir. Bu pasif bir tehdittir — hedefe hiçbir şey gönderilmez, yalnızca onun sızdırdığı dinlenir; bu da onu tespit etmeyi zorlaştırır.

| Yan-kanal türü | Sızıntı kaynağı | Savunma yönü |
|---|---|---|
| Yayılan emisyon (radiated) | Ekran/kablo/devreden hava yoluyla EM yayılım | Ekranlama (shielding), TEMPEST sertifikalı ekipman, mesafe |
| İletilen emisyon (conducted) | Güç/veri hatları üzerinden sızıntı | Filtreleme, izolasyon, ayrı besleme |
| Zamanlama yan-kanalı | İşlem süresinin veriyle korelasyonu | Sabit-zaman uygulamalar (kripto), gürültü ekleme |

Savunma farkındalığı (bu bölüm bağlamında):

- En hassas işlemleri (anahtar üretimi, parola girişi) ekranlanmış/izole ortamda yapmak, fiziksel mesafe ve ekranlamanın yayılan emisyonu hızla zayıflattığı gerçeğine dayanır (mesafeyle güç düşüşü — Bölüm 1 FSPL).
- Kritik kabloların ekranlanması ve düzenlenmesi, iletilen ve yayılan sızıntıyı azaltır.
- Bu, "saldırgan sana bir şey göndermeden seni dinleyebilir" gerçeğinin savunma karşılığıdır ve neden bazı kurumların fiziksel/emisyon güvenliğine yatırım yaptığını açıklar.

> Mühendislik sezgisi: Aktif RF tehdidi (karıştırma/sahteleme) "sana enerji gönderir"; yan-kanal tehdidi "senin sızdırdığın enerjiyi dinler". İlkine karşı dayanıklılık + tespit; ikincisine karşı ekranlama + mesafe + sabit-zaman tasarım. İkisi RF tehdit manzarasının iki yüzüdür.

Çapraz referans: Tam ayrıntı Bölüm 6 (TEMPEST, RF sızıntısı, emisyon güvenliği, OPSEC). Burada yalnızca tehdit manzarası bütünlüğü için anılmıştır.

---

<a id="12"></a>
## 12. Kritik Altyapı RF Riski: GPS/Zaman, SCADA Telemetri ve Savunma

RF tehdidinin en yüksek bahisli alanı kritik altyapıdır: enerji şebekesi, telekom, finans, su, ulaşım ve endüstriyel kontrol (ICS/SCADA). Bu sistemler genellikle RF'e iki kritik noktadan bağımlıdır — konum/zaman (GNSS) ve telemetri/kontrol bağları — ve her ikisi de bu bölümde anlatılan tehditlere açıktır. Savunmacı için buradaki ders, "tek bir RF bağına asla varoluşsal biçimde bağımlı olma" ilkesidir.

### 12.1 GNSS zaman bağımlılığı (sessiz kritiklik)

Çoğu insan GNSS'i konum için bilir; oysa kritik altyapının en sessiz ve en yaygın GNSS bağımlılığı zamandır. Telekom şebekeleri, enerji şebekesi senkronizasyonu, finans işlem damgalama ve SCADA olay sıralaması, GNSS'ten gelen hassas zaman/frekans referansına dayanabilir. GNSS karıştırma veya sahteleme, yalnızca "konumu" değil, bu zaman omurgasını da tehdit eder.

| Altyapı | GNSS bağımlılığı | Karıştırma/sahteleme etkisi | Savunma |
|---|---|---|---|
| Telekom | Baz istasyonu zaman/frekans senkronu | Senkron kaybı, kapasite/kalite düşüşü | Yerel atomik saat holdover, PTP/ağ zamanı yedeği |
| Enerji şebekesi | PMU/koruma zaman damgası | Yanlış olay sıralaması, koruma hatası | Çoklu zaman kaynağı, holdover osilatör |
| Finans | İşlem zaman damgası (mevzuat) | Damgalama tutarsızlığı | Yedek zaman dağıtımı, denetim izleri |
| Ulaşım/lojistik | Konum/zaman | Yanlış konum, gecikme | Çoklu GNSS + INS + harita eşleme |

### 12.2 SCADA/ICS telemetri bağları

Endüstriyel kontrol sistemleri, uzak sahalardan (boru hattı, trafo, pompa istasyonu) veri toplamak için sıklıkla RF telemetri (lisanslı telsiz, hücresel modem, uydu) kullanır. Bu bağlar:

- Karıştırmaya açıktır (telemetri kesintisi → körlük veya yanlış "son bilinen değer" ile çalışma).
- Bazı eski kurulumlarda zayıf veya hiç kimlik doğrulamasız olabilir (komut spoofing riski — Bölüm 9/10 dersleri).

```
   KRİTİK ALTYAPI RF BAĞIMLILIĞI ve YEDEKLİLİK

                 GNSS (konum + ZAMAN)
                      ░░░  ◄── kırılgan (zayıf sinyal)
                       │
        ┌──────────────┼──────────────┐
        │              │              │
     Telekom        Enerji          Finans
     senkron        PMU/zaman       damgalama
        │              │              │
        └──────► YEDEKLİLİK ◄─────────┘
        (yerel atomik saat / holdover osilatör,
         ağ zamanı PTP, çoklu GNSS, INS, kablolu yedek)

   SCADA telemetri ── RF bağı ── kontrol merkezi
        │                              │
        └── savunma: kimlik doğrulama + bütünlük + bağ-kaybı güvenli durumu
```

### 12.3 Savunma ilkeleri (kritik altyapı)

1. Zaman için GNSS'i tek kaynak yapma: Yerel holdover (atomik/OCXO osilatör), ağ tabanlı zaman (PTP), ve çoklu GNSS takımyıldızı birlikte, GNSS kesintisinde sistemin bir süre doğru zamanı koruyabilmesini sağlar.
2. Telemetri bağlarını kimliklendir ve bütünlüğünü koru: Komutların kimden geldiği doğrulanmalı (auth), değiştirilmediği garanti edilmeli (MAC) — Bölüm 10 merdiveni.
3. Bağ kaybında güvenli duruma geç: Telemetri/GNSS kaybında sistem "son değeri körü körüne kullanmak" yerine tanımlı bir güvenli duruma (fail-safe) geçmeli.
4. Sürekli spektrum izleme (ES): Kritik tesis çevresinde gürültü tabanı/anomali izleme (Bölüm 5), bir karıştırma olayını erken yakalar.
5. Yetkili tespit + raporlama zinciri: Karıştırma şüphesinde kaynağı yön bulmayla konumlandırıp (Bölüm 5) yetkili telekom/güvenlik makamına raporlama.

> Mühendislik sezgisi: Kritik altyapı savunmasının özü tek cümledir — "hiçbir varoluşsal işlev tek bir RF bağına bağlı olmamalı". GNSS karıştırılabilir; o halde zaman için yedek saat olmalı. Telemetri kesilebilir; o halde güvenli duruma geçiş olmalı. Tek nokta bağımlılığı, RF tehdidinin en sevdiği hedeftir.

Yasal uyarı: Kritik altyapıya yönelik karıştırma/sahteleme/yetkisiz telemetri müdahalesi en ağır suç kategorilerindendir ve toplum güvenliğini doğrudan tehdit eder. Savunmacının alanı dayanıklılık tasarımı + pasif tespit + yetkiliye raporlamadır.

---

<a id="13"></a>
## 13. Tehdit → Tespit → Savunma Matrisi (Birleşik Referans)

Aşağıdaki birleşik matris, bölümün tüm tehditlerini tek bir savunmacı referansında toplar. Her satır: tehdit (kavram), nasıl tespit edilir (ES, pasif), nasıl savunulur (EP, tasarım), ve yasal not. Bu tablo, bir olayla karşılaştığında hızlı yönlendirme içindir.

| Tehdit (kavram) | Fiziksel temel | Tespit (pasif/ES) | Savunma (EP/tasarım) | Yasal |
|---|---|---|---|---|
| Barrage karıştırma | Geniş banda gürültü | Geniş bant gürültü tabanı sıçraması | Yayma + yedek bant/sistem | Karıştırma = suç |
| Spot karıştırma | Tek frekansa enerji | Beklenmedik kalıcı dar tepe | FHSS (tek noktayı tutamaz) | Karıştırma = suç |
| Sweep karıştırma | Banttan banda kayan enerji | Waterfall'da gezen diyagonal | FHSS + dar bant filtre | Karıştırma = suç |
| Reaktif karıştırma | Yayın algılayınca vur | Olayla zaman korelasyonu | Kısa burst + atlama + LPD | Karıştırma = suç |
| GNSS karıştırma | Zayıf GNSS'i bastırma | C/N0 düşüşü + uydu kaybı + gürültü tabanı | Anti-jam anten (CRPA), çoklu takımyıldız, INS | Suç + can güvenliği |
| GNSS sahteleme | Sahte GNSS sinyali | Tutarsızlık (INS/zaman/konum sıçraması) | Tutarlılık kontrolü, çok kaynak, sahtecilik tespiti | Suç + can güvenliği |
| Drone (tehdit) | Çok bağlı uçan RF | RF imza + DF + Remote ID füzyonu | Tespit + raporlama + yetkili C-UAS | Tespit yasal; müdahale yetki ister |
| Sabit-kod replay | Aynı kod tekrarı | Kendi cihazında tekrar deseni gözlemi | Rolling-code/nonce | Başkasına = suç |
| RF relay (keyless) | Mesafe yanıltma | (Tasarım tarafı) mesafe tutarsızlığı | UWB ToF, Faraday kesesi | Başkasına = suç |
| Genel spoofing | Kimlik doğrulanmaması | Anomali/tutarsızlık | Auth + nonce + MAC | Yetkisiz = suç |
| Yan-kanal/TEMPEST | İstem dışı EM sızıntı | (Pasif, tespiti zor) | Ekranlama + mesafe + sabit-zaman | Dinleme bağlama göre suç |
| SCADA telemetri saldırı | Zayıf/karıştırılan bağ | Bağ kaybı + anomali izleme | Auth+MAC+güvenli durum+yedek | Yetkisiz = ağır suç |

### Manipüle edilebilir / edilemez hızlı ayrım

| Sinyal/sistem türü | Manipülasyona açıklık | Belirleyici özellik |
|---|---|---|
| Sabit-kod RC, açık analog video, korumasız GNSS bağımlılığı | Yüksek | Kimlik/tazelik/yayma yok |
| Rolling-code (zayıf uygulama) | Orta | Tazelik var ama kripto/sayaç zayıf |
| Auth + nonce + MAC + UWB + yedeklilik | Düşük | Dayanıklılık merdiveninin tamamı |
| FHSS/DSSS + kriptografik kimlik doğrulamalı bağ | Düşük | Yayma + kimlik birlikte |

> Mühendislik sezgisi: Bu matris bir "ne yapayım" cetvelidir. Soldan sağa oku: tehdidi tanı → pasif olarak tespit et → tasarımsal dayanıklılıkla savun → yasal sınırı hatırla (icra değil, savunma). Hiçbir satırda savunmacının doğru cevabı "karşı saldırı" değildir; her zaman tespit + dayanıklılık + yetkiliye raporlamadır.

---

<a id="14"></a>
## 14. Alıştırmalar (Yalnızca Yasal/Savunma)

Aşağıdaki alıştırmalar tamamen yasal, pasif ve savunma odaklıdır; yalnızca kendi cihazların ve açık/yasal sinyallerle sınırlıdır. Hiçbiri iletim, karıştırma, yetkisiz erişim veya içerik manipülasyonu içermez. Her birinde yasal sınır hatırlatılmıştır.

### A) Kendi cihazının RF imzasını ve güç desenini ölçmek (pasif)

Amaç: Bir vericinin (kendi Wi-Fi yönlendiricin, kendi Bluetooth cihazın, kendi uzaktan kumandan — yalnızca kendine ait olanlar) spektral imzasını ve güç desenini tanımak; "normal" baseline'ı kurmak.

Adımlar (kavramsal):
1. SDR ile (Bölüm 2) ilgili ISM bandını izle; cihazını çalıştır ve durdur, waterfall'daki imzayı gözlemle.
2. Cihaza yaklaşıp uzaklaşırken alıcıdaki gücün nasıl değiştiğini not et (mesafeyle düşüş — Bölüm 1 FSPL sezgisi).
3. Bu imzayı "baseline" olarak kaydet; bu, ileride anomaliyi (örneğin beklenmedik bir kaynak) tanımanın referansıdır.

Yasal sınır: Yalnızca kendi cihazın. Başkasının cihazını hedeflemek/içeriğini çözmek bu alıştırmanın dışındadır.

### B) Bir karıştırma olayını spektrumda TESPİT etmeyi simüle etmek (gürültü tabanı izleme — yayın yok)

Amaç: Karıştırma tespitinin (Bölüm 5) çekirdek becerisini, hiçbir karıştırıcı kullanmadan, yalnızca gürültü tabanı izleyerek pratik etmek.

Adımlar (kavramsal, yayın yok):
1. Bir bandın gürültü tabanını uzun süre kaydet; saatlik normal dağılımını (baseline) çıkar.
2. Normal varyasyonun istatistiğini (ortalama, standart sapma) hesapla; bir "alarm eşiği" tanımla (örn. baseline + birkaç σ, N saniye süreklilik).
3. Banttaki doğal/yasal yoğun bir kullanım anında (örneğin meşru bir cihazın yoğun yayını) tabanın nasıl yükseldiğini gözlemle ve eşiğinin bunu nasıl yakalayacağını test et. Böylece gerçek bir karıştırmanın imzasını (kalıcı, geniş yükseliş) zihninde modellersin.

Önemli: Bu alıştırma hiçbir karıştırıcı yapmaz/kullanmaz. Yalnızca pasif izleme ve eşik mantığı pratiğidir. Gerçek karıştırma üretmek suçtur.

### C) Kendi eski kumandanın sabit-kod mu rolling-code mu olduğunu analiz etmek

Amaç: Bölüm 9'daki replay dersini kendi cihazında somutlaştırmak.

Adımlar (kavramsal):
1. Kendi eski garaj/kapı/araba kumandanın (yalnızca kendine ait) yayın frekansını (cihaz etiketinden veya yasal referanstan) öğren ve SDR ile gözlemle.
2. Aynı tuşa arka arkaya bas ve her basıştaki yayının waterfall/zaman görünümünün aynı mı yoksa değişen mi olduğunu gözlemle.
3. Yorumla: Her seferinde görünüş aynıysa sabit-kod ailesine, değişiyorsa rolling-code ailesine işaret eder. Bu, "neden eski sabit-kod cihazları değiştirmeliyim" dersini elle deneyimletir.

Çok önemli: Yalnızca kendi cihazın; yalnızca gözlem. Kodları kaydedip tekrar oynatmak (replay), kendi cihazında bile başkasının erişimini etkileyebilecek bir senaryoya dönüşmemeli; amaç sabit/rolling ayrımını görmektir, kod yakalama değil. Başkasının cihazı kesinlikle hariç.

### D) Drone Remote ID yayınını (yasal/açık) gözlemlemek

Amaç: Counter-UAS savunmasının (Bölüm 8) en yasal tespit katmanını — Remote ID — tanımak.

Adımlar (kavramsal):
1. Bölgendeki Remote ID düzenlemesini ve hangi açık kablosuz yayının (genellikle Wi-Fi/BT tabanlı) kullanıldığını öğren.
2. Yasal/açık bir ortamda (örneğin kendi uyumlu dronunu uçururken ya da kamuya açık bir gösteri alanında) Remote ID yayınını pasif olarak gözlemle; yayınlanan kimlik/konum alanlarının yapısını anla.
3. Yorumla: Bu, "şeffaflıkla tespit" yaklaşımının nasıl çalıştığını ve savunmacının uyumlu dronları neden kolayca tanıyabildiğini gösterir.

Yasal sınır: Yalnızca açık/yasal Remote ID yayınının pasif gözlemi. Başka bir dronun kontrol bağına müdahale, içerik çözme veya konum sahteleme bu alıştırmanın dışındadır ve suçtur.

### E) Dayanıklılık denetimi (kâğıt üstü düşünce deneyi)

Amaç: Bölüm 10'daki dayanıklılık merdivenini kendi kullandığın bir RF sistemine uygulamak.

Adımlar:
1. Kullandığın bir RF sistemini seç (akıllı kapı zili, garaj kumandası, IoT sensör, araç anahtarı).
2. Bölüm 10'un beş sorusunu uygula: kimlik doğrulama var mı? tazelik (nonce/sayaç)? bütünlük (MAC)? mesafe ölçümü (relay'e karşı)? yedeklilik?
3. Eksik kalan basamakları işaretle ve "hangi tehdide açık" sütununu doldur. Bu, gerçek bir tehdit modelleme refleksi kazandırır.

Amaç: Savunmanın yalnızca "şifreleme var mı" değil, katmanlı bir denetim listesi olduğunu içselleştirmek; bir tüketici/operatör olarak hangi soruları sorman gerektiğini öğrenmek.

---

<a id="15"></a>
## 15. Kapanış: Savunmacı Zihniyeti, Etik ve Yasal Sınır + Çapraz Referans

### Kavram kartı

| Kavram | Bir cümlelik öz |
|---|---|
| EA / EP / ES | Saldırı (tehdit) / koruma (kalkan) / destek (göz); savunmacı EP+ES'e yatırım yapar |
| J/S oranı | Karıştırmanın tüm fiziği; zayıf sinyaller (GNSS) yapısal olarak kırılgan |
| Karıştırıcı türleri | Barrage/spot/sweep/protokol-bilinçli/reaktif — yalnızca teşhis imzaları, yapım yok |
| Karıştırma etkisi | Teknik arıza değil, can/altyapı güvenliği meselesi |
| Karıştırma tespiti | Gürültü tabanı + imza + sistem kalitesi çapraz doğrulama + yön bulma |
| Karıştırma savunması | Yayma + atlama + güç + null-steering + yedeklilik (J/S'in her terimi) |
| Drone RF | Çok bağlı uçan IoT; en zayıf bağ seviyeyi belirler; GNSS en kırılgan duyusu |
| Counter-UAS | TESPİT yasal (RF+DF+Remote ID); ETKİSİZ KILMA yetki ister |
| Replay/spoof/relay | Tekrar/doğrulanmama/ölçülmeyen mesafe; çözüm nonce/auth/UWB |
| Dayanıklılık merdiveni | Auth → tazelik → bütünlük → mesafe → yayma → yedeklilik |
| Yan-kanal/TEMPEST | "Senin sızdırdığını dinler"; savunma ekranlama+mesafe+sabit-zaman |
| Kritik altyapı RF | Hiçbir varoluşsal işlev tek RF bağına bağlı olmamalı (özellikle GNSS zaman) |

### Ezber sezgiler

- Tehdidi tanı, kalkanı kur, gözü aç (EA-tanı, EP-kur, ES-aç).
- Karıştırmanın tüm fiziği J/S'tir; alıcıya zayıf ulaşan sinyal kolay bastırılır — bu yüzden GNSS yapısal olarak kırılgandır.
- Karıştırmayı yenmenin üç fiziksel kolu: işleme kazancını artır, karıştırıcı yönündeki anten kazancını düşür, mesafeyi/yedekliliği lehine kullan.
- Karıştırma tespiti üç katmanın füzyonudur: gürültü tabanı seviyesi + imza şekli/zamanı + sistem kalitesi çöküşü; kaynak için dördüncü katman yön bulma.
- Bir drone uçan, çok bağlı bir IoT'dur; en zayıf radyo bağı tüm güvenliği belirler.
- Drone savunmasının doğru sırası: tespit (pasif, yasal) → raporla → etkisiz kılmayı yetkiliye bırak.
- "Tekrar edilebilen ve doğrulanmayan şey güvenli değildir": replay'i nonce, spoof'u auth, relay'i gerçek mesafe ölçümü çözer.
- "Rolling" bir etikettir; koruma kriptografik kalitede ve doğru sayaç yönetimindedir.
- Şifreleme içeriği korur; kimlik doğrulama ve tazelik olmadan sahte komut ve replay durmaz.
- Kritik altyapıda tek nokta bağımlılığı (özellikle GNSS zaman) en büyük risktir; yedeklilik şarttır.
- Savunmacının doğru cevabı asla "karşı saldırı/karşı yayın" değildir; her zaman tespit + dayanıklılık + yetkiliye raporlamadır.

### Etik ve yasal sınır (bölümün özü)

Bu bölüm baştan sona tek bir ayrımı korudu: tehdidi tanımak ile onu icra etmek arasındaki çizgi. Bir savunmacı, bir karıştırmanın fiziğini, bir GNSS sahtelemesinin neden mümkün olduğunu, bir relay saldırısının mantığını bilmek zorundadır — aksi halde karşısına doğru kontrolü koyamaz. Ama bu bilgi, yalnızca tanıma, tespit ve savunma içindir. Bu metin bilinçli olarak hiçbir karıştırıcı yapımı, hiçbir drone/uydu/sistem ele geçirme adımı, hiçbir aktif saldırı reçetesi vermedi ve vermez.

Yasal gerçek nettir ve tekrar edilmeyi hak eder: Karıştırma, yetkisiz drone müdahalesi, GNSS sahteleme, keyless relay hırsızlığı ve kritik altyapıya RF müdahale dünyanın hemen her yerinde ağır suçtur; çoğu zaman can güvenliğini doğrudan tehdit eder ve cezaları buna göre ağırdır. Pasif tespit, dayanıklı tasarım ve yetkiliye raporlama yasaldır ve teşvik edilir; aktif müdahale (özellikle herhangi bir yayın) yetki ister ve genellikle yalnızca yetkili kurumların alanıdır.

> Kapanış: RF spektrumu görünmez ama savunmasız değildir. Her karıştırma bir gürültü tabanı yükseltir, her sahteleme bir tutarsızlık bırakır, her zayıf protokol bir denetim sorusuna takılır. Savunmacının işi bu izleri tanımak, pasif olarak tespit etmek, sistemi dayanıklı tasarlamak ve sınırı geçmeden — yetkiliyle birlikte — yanıt vermektir. Tehdidi en iyi, onu icra eden değil, onu tanıyıp karşısına doğru kalkanı koyan kişi yener.
>
> Bu doküman Kanije Kalesi güvenlik/teknik rehberleri koleksiyonunun SIGINT serisinin 13. bölümüdür. Tehdidi tanıma + savunma + yasal sınır üçlüsü, serinin omurgasıdır.

---

### Serinin diğer bölümleri (çapraz referans)

- SIGINT_01 — RF Fiziği, Spektrum ve Modülasyon: FSPL, dB/dBm, SNR, yayılı spektrum temeli. (J/S oranının ve mesafeyle güç düşüşünün fiziksel temeli; yayma/işleme kazancının kökü orada.)
- SIGINT_02 — SDR Donanımları ve Yazılım Ekosistemi: RTL-SDR/HackRF/USRP, waterfall okuma. (Karıştırma tespiti ve RF imza alıştırmalarının pasif ölçüm tabanı orada.)
- SIGINT_03 — Anten Teorisi, Kapsama ve Yön Bulma (DF): yön bulma, dizi, null-steering temeli. (Karıştırıcı ve drone kaynak tespitinin, null-steering savunmasının matematiği orada.)
- SIGINT_04 — Sayısal Demodülasyon ve Protokol Kod Çözme: ADS-B, AIS vb. (Drone/telemetri bağlarının protokol bağlamı orada.)
- SIGINT_05 — Hücresel, WiFi/BT ve IoT Spektrumu (Savunma Bakışı): GSM/LTE/5G, ISM imzaları. (Karıştırmanın hücresel etkisi ve IoT/keyless bağların yüzeyi orada.)
- SIGINT_06 — TEMPEST, RF Sızıntısı ve SIGINT'e Karşı Savunma: emisyon güvenliği, ekranlama, OPSEC. (Bölüm 11'in yan-kanal/TEMPEST farkındalığının tam ayrıntısı orada.)
- SIGINT_07 — Disiplinler ve Sinyal Ayıklama: PDW, kümeleme, deinterleaving, LPI/LPD, trafik analizi. (Karıştırıcı imza ayıklama, LPI/LPD'nin karıştırma-direnci olarak okunması ve drone tespit füzyonu oranın yöntemlerine dayanır.)
- SIGINT_08 — Frekans Tahsisi ve Bant Planı: tahsis tabloları, hizmet eşleştirme. (Beklenmedik bantta enerji = anomali yorumunun referans tablosu orada.)

> Bu bölüm, serinin aktif-tehdit yüzünü tamamen savunmacı/CTI perspektifinden kapatır: jamming, drone RF, replay/spoof/relay ve kritik altyapı riski — hepsi prensip + tespit + savunma + yasal çerçevede, hiçbir saldırı reçetesi olmadan. İlgili kale rehberleri: `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`, `VERACRYPT_USTALIK_REHBERI.md`.
