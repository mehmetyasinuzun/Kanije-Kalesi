# SIGINT EL KİTABI — BÖLÜM 25: AÇIK KAYNAK ARAÇ VE KAYNAK DİZİNİ — GITHUB VE ÖTESİ

## Meşru, Açık-Kaynak SDR/SIGINT/Kablosuz-Güvenlik Araçlarının Kategorize Başvuru Dizini

> Amaç: Bu bölüm bir okuma metni değil bir BAŞVURU DİZİNİdir. Önceki bölümler her aracın ne işe yaradığını, hangi iş akışına oturduğunu ve fiilen nasıl çalıştırıldığını anlattı (özellikle Bölüm 04, 12, 15, 16). Burada o araçları tek çatı altında, kaynağına kadar indirgenmiş biçimde topluyoruz: her aracın ne yaptığı, projenin GitHub/depo adresi (`org/repo` biçiminde), olgunluk ve bakım durumu, hangi görevde ne ölçüde kullanıldığı ve lisans notu. Hedef, "şu işi yapan açık kaynak araç hangisiydi, resmî deposu neredeydi, ne kadar güvenilir" sorusuna tek sayfadan bakıp cevap verebilmen. Bu bölüm sık sık başvurulmak, üstünden tarayarak okunmak için tasarlandı; baştan sona düz okunmak için değil.

> Yasal çerçeve: Bu dizindeki her araç meşru, açık-kaynak, kamuya açık bir projedir. Bir aracın "açık kaynak ve indirilebilir" olması, onunla yapacağın her şeyin yasal olduğu anlamına GELMEZ. Alıcı (RX) araçları geniştir; ama yakaladığın içeriğin kaydı/çözümü/paylaşımı, her türlü yayın (TX), hücresel test hücresi kurma, GPS sinyal üretme, yetkisiz ağa pentest — bunların yasallığı tamamen ülkene, izin durumuna ve hedefe bağlıdır. Hücresel, GNSS ve WiFi/kablosuz-güvenlik başlıklarında sınır özellikle nettir ve her tabloda tekrar edilir: yalnızca kendi cihazın, kendi izole/Faraday test ortamın, yazılı yetki verilmiş kapsam veya açık/kamuya yayınlanan sinyaller. Bu kitap hukuki danışmanlık değildir; kendi ülkenin güncel mevzuatını teyit et.

> Depo adresleri hakkında kritik uyarı: Aşağıdaki `org/repo` yolları, yazım anında bilinen resmî depolardır. Açık kaynak dünyasında depolar taşınır, organizasyon adı değişir, ana geliştirme çatallanır (fork) ve kötü amaçlı taklit çatallar türeyebilir. Bir depo yolunu körlemesine `git clone`'lamadan önce, projenin resmî sitesinden veya bilinen topluluk kaynağından teyit et. Adresinden tam emin olmadığım yerlerde "GitHub'da '<isim>' aratıp resmî depoyu teyit et" notu bıraktım; uydurma/yanlış yol vermektense teyit notu vermeyi tercih ettim. Bir GitHub projesinin güvenilirliğini değerlendirmenin yöntemi §16'da; kötü amaçlı çatal tuzakları §17'de.

---

## İÇİNDEKİLER

1. [Bu Dizin Nasıl Okunur: Sütunlar, Olgunluk Ölçeği, Lisans Notasyonu](#1)
2. [SDR Çerçeveleri ve Altyapı](#2)
3. [Sürücüler ve Cihaz Kütüphaneleri](#3)
4. [Alıcı ve GUI Uygulamaları](#4)
5. [Çözücü ve Uygulama Araçları](#5)
6. [Analiz ve Tersine Mühendislik](#6)
7. [Hücresel Araştırma (Yasal/Akademik Sınır)](#7)
8. [GNSS / GPS (Sınır Net)](#8)
9. [Yön Bulma ve Konum Tespiti](#9)
10. [WiFi / WLAN Güvenliği (Yetkili Pentest)](#10)
11. [BLE / RFID / IoT Donanım Araçları](#11)
12. [Veri Kümeleri ve Makine Öğrenmesi](#12)
13. [Veritabanı, Topluluk ve Çevrimiçi Kaynaklar](#13)
14. [Kurulum Dağıtımları (Hazır OS İmajları)](#14)
15. [Donanım Açık-Tasarım Depoları](#15)
16. [Bir GitHub Projesi Nasıl Değerlendirilir](#16)
17. [Çatal Tuzakları, Derleme/Bağımlılık ve Güvenlik](#17)
18. [Alıştırmalar (Yasal)](#18)
19. [Hızlı Referans ve Diğer Bölümler](#19)

---

<a id="1"></a>
## 1. Bu Dizin Nasıl Okunur: Sütunlar, Olgunluk Ölçeği, Lisans Notasyonu

Bu bölümdeki tablolar tutarlı bir sütun şemasına oturur. Şemayı bir kez öğrenince tüm dizin hızla taranabilir hale gelir.

### Sütun şeması

| Sütun | Ne anlatır |
|---|---|
| Araç | Projenin yaygın adı |
| Depo (`org/repo`) | Resmî GitHub/proje yolu; emin olunmayan yerlerde teyit notu |
| Görev / kapsam | Aracın hangi işte ne ölçüde kullanıldığı; ilgili bölüm referansı |
| Olgunluk | Bakım ve üretime hazırlık durumu (aşağıdaki ölçek) |
| Lisans | Yaygın bilinen lisans ailesi (teyit edilmeli; lisans zamanla değişebilir) |

### Olgunluk ölçeği

Olgunluk, "ne kadar yıldız aldığı" değil, *üretimde/sahada bugün ne kadar güvenle dayanılabileceği*dir. Beş kademeli bir ölçek kullanıyorum:

```
   OLGUNLUK KADEMELERİ
   ─────────────────────────────────────────────────────────────────
   Referans     → fiili standart; ekosistem bunun üstüne kurulu,
                  geniş bakım, kararlı API (örn GNU Radio, aircrack-ng).
   Olgun        → uzun süredir bakımlı, üretime uygun, geniş kullanıcı.
   Aktif        → düzenli geliştirilen, canlı topluluk, hızlı değişebilir.
   Niş/uzman    → belirli bir işte çok iyi ama dar kapsam veya az bakım.
   Arşiv/durağan → tarihsel değeri var ama bakım yavaş/durmuş;
                  kullanırken çatal durumunu teyit et.
```

> Not: Olgunluk durağan değildir. Bir proje aylar içinde "Aktif"ten "Arşiv"e geçebilir (geliştirici çekilir) veya tersi (yeni çatal canlanır). Bu yüzden her tabloda olgunluk *yazım anı tahminidir* ve §16'daki canlılık kontrolüyle (son commit, açık issue, sürüm sıklığı) doğrulanmalıdır.

### Lisans notasyonu

Lisans sütunundaki kısaltmalar yaygın açık-kaynak lisans ailelerine işaret eder; tam metin ve sürüm her zaman deponun `LICENSE` dosyasından teyit edilmelidir:

| Kısaltma | Aile | Pratik anlamı (kabaca; hukuki tavsiye değildir) |
|---|---|---|
| GPL-2 / GPL-3 | GNU GPL | Güçlü copyleft; türev işi de aynı lisansla açman beklenir |
| LGPL | Lesser GPL | Kütüphane olarak bağlamada daha esnek copyleft |
| AGPL-3 | Affero GPL | Ağ servisi olarak sunmak da kaynak açmayı tetikler |
| MIT / BSD | İzin verici | Atıf yeterli; kapalı türeve izin |
| Apache-2 | İzin verici | MIT benzeri + patent hükmü |
| Karışık | Birden çok | Alt bileşenler farklı lisanslarda; dikkatle oku |

> Kritik: Lisans, "bedava mı" değil "ne yapabilirim" sorusudur. Bir aracı kişisel/eğitim amaçlı çalıştırmak çoğu lisansta serbesttir; ama onu bir ürüne gömmek, dağıtmak veya servis olarak sunmak GPL/AGPL ailesinde yükümlülük doğurur. Üretim/dağıtım düşünüyorsan lisansı teyit et. Bu kitaptaki lisans notları yalnızca yönlendiricidir, hukuki görüş değildir.

### Çapraz referans haritası

Bu dizin tek başına durmaz; her aracın *fiili kullanımı* serinin başka bölümlerindedir. Dizinden bir araç bulduğunda, derinlik için ilgili bölüme git:

| Konu | Dizin kategorisi (bu bölüm) | Derinlemesine bölüm |
|---|---|---|
| Kurulum, OS, sürücü, VM/USB | §3, §14 | Bölüm 04 |
| Araçların fiili komut/iş akışı | §2-§9 | Bölüm 12 |
| WiFi/WLAN güvenliği | §10 | Bölüm 15 |
| BLE/RFID/IoT kısa menzil | §11 | Bölüm 16 |
| Protokol/çözümleme prensibi | §5 | Bölüm 05 |
| GNSS/GPS sistemleri | §8 | Bölüm 10 |
| Yön bulma/takip | §9 | Bölüm 09 |
| İleri hücresel 4G/5G | §7 | Bölüm 20 |
| Yapay zeka / ML SIGINT | §12 | Bölüm 19 |
| Güvenlik/açık/savunma sınırı | tüm tablolar | Bölüm 06, Bölüm 00 |

---

<a id="2"></a>
## 2. SDR Çerçeveleri ve Altyapı

Bu katman, üstündeki tüm araçların oturduğu temeldir: sinyal akışı çatıları, donanım soyutlama katmanları ve blok/akış paradigması. Çoğu üst-seviye araç (alıcı, çözücü, analizör) doğrudan veya dolaylı olarak bu projelerden birine bağımlıdır. Fiili kullanım ve akış kurma Bölüm 12 §4'te.

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| GNU Radio | `gnuradio/gnuradio` | Blok-akış DSP çatısı; ekosistemin omurgası. Kendi alıcı/çözücünü bloklardan kurmak (GRC), `gr-*` eklenti tabanı. Bölüm 12 §4. | Referans | GPL-3 |
| SoapySDR | `pothosware/SoapySDR` | Donanım soyutlama katmanı; farklı SDR'ları tek API altında toplar. GQRX/SatDump/SDRangel "ne varsa kullan" der. Bölüm 04, Bölüm 12 §2. | Referans | Boost (izin verici) |
| gr-osmosdr | `osmocom/gr-osmosdr` | GNU Radio için çok-donanımlı kaynak/sink bloğu (osmocom Source). RTL/HackRF/Airspy/USRP'yi GRC akışına bağlar. Not: bazı dağıtımlarda `gnuradio/gr-osmosdr` çatalı kullanılır — teyit et. | Olgun | GPL-3 |
| Pothos | `pothosware/PothosCore` | GNU Radio'ya alternatif, düğüm-tabanlı veri-akış çatısı ve GUI (Pothos Flow). Niş ama tutarlı. | Niş/uzman | Boost |
| LuaRadio | `vsergeev/luaradio` | Hafif, bağımlılığı az, Lua tabanlı DSP çatısı; gömülü/öğrenme için. GNU Radio'ya göre minimal. | Niş/uzman | MIT |
| VOLK | `gnuradio/volk` | SIMD-hızlandırılmış vektör çekirdek kütüphanesi; GNU Radio performansının altındaki kas. Genelde GNU Radio ile gelir. | Referans | GPL-3 (teyit: LGPL'e geçiş tartışmaları olmuştur) |

> Pratikte: Bu katmanda neredeyse her şey GNU Radio ve SoapySDR ikilisine indirgenir. Bir üst-seviye aracın "neden çalışmadığını" kazarken sorun çoğu zaman bu katmandaki sürüm uyumudur (yanlış `gr-*` sürümü, eksik SoapySDR modülü). Bölüm 04, bu bağımlılık cehennemini ve DragonOS'un onu nasıl paketlediğini ayrıntılandırır.

---

<a id="3"></a>
## 3. Sürücüler ve Cihaz Kütüphaneleri

SoapySDR bir soyutlama katmanıdır; altında her donanım için gerçek sürücü/kütüphane çalışır. Bir cihazın "görünmemesi" sorununun kökü neredeyse her zaman bu katmandır (yanlış sürücü sürümü, eksik udev kuralı, çakışan çekirdek modülü). Kurulum derinliği ve Zadig/blacklist tuzakları Bölüm 04'te.

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| rtl-sdr (Osmocom) | `osmocom/rtl-sdr` | RTL2832U tabanlı ucuz DVB-T çubuklarının temel sürücüsü ve `rtl_test`/`rtl_fm`/`rtl_power`/`rtl_sdr` araçları. Tarihsel ana hat. | Referans | GPL-2 |
| rtl-sdr (Blog çatalı) | `rtlsdrblog/rtl-sdr-blog` | RTL-SDR Blog V3/V4 donanımına özel düzeltmeler (özellikle V4 için güncel `librtlsdr`). V4 görünmüyorsa bu çatal gerekir. Bölüm 04 V4 notu. | Aktif | GPL-2 |
| HackRF | `greatscottgadgets/hackrf` | HackRF One sürücüsü, firmware ve `hackrf_info`/`hackrf_transfer` araçları. Geniş-bant TX/RX donanımının resmî deposu. | Referans | GPL-2 (firmware ayrı lisanslar) |
| Airspy | `airspy/airspyone_host` | Airspy R2/Mini host kütüphanesi ve `airspy_info`/`airspy_rx`. HF+ için ayrı `airspy/airspyhf` deposu. | Olgun | BSD |
| LimeSuite | `myriadrf/LimeSuite` | LimeSDR (USB/Mini) sürücüsü, kalibrasyon ve `LimeUtil`/`LimeQuickTest`. Not: daha yeni `LimeSuiteNG` geçişi var — güncel hat için teyit et. | Olgun | Apache-2 |
| UHD | `EttusResearch/uhd` | USRP ailesinin resmî sürücü çatısı; `uhd_find_devices`/`uhd_usrp_probe` ve firmware imajları. Profesyonel/araştırma sınıfı. | Referans | GPL-3 |
| SoapyRemote | `pothosware/SoapyRemote` | Bir SDR'ı ağ üzerinden uzaktan paylaşma SoapySDR modülü; sunucu-istemci. Uzak istasyon/baş düğüm mimarisi. | Olgun | Boost |
| SoapySDR modülleri (genel) | `pothosware/SoapyRTLSDR`, `pothosware/SoapyHackRF`, `pothosware/SoapyAirspy`, `pothosware/SoapyUHD` | Her donanımı SoapySDR'a bağlayan köprü modülleri. Donanımın SoapySDR altında görünmesi için ilgili modül kurulu olmalı. | Olgun | Boost/çeşitli |
| SDRplay API | (kapalı kaynak ikili) | RSP serisi resmî API servisi. AÇIK KAYNAK DEĞİL; üreticiden ikili indirilir. Burada bütünlük için anılır; dizinin geri kalanı açık kaynaktır. | — | Tescilli |

> Not: SDRplay satırı dışında bu katmanın tamamı açık kaynaktır. SDRplay'i listeye ekosistemde çok kullanıldığı için kattım; ama açık-kaynak kuralının istisnası olduğunu açıkça işaretledim. Donanım takıp `SoapySDRUtil --find` cihazı göstermiyorsa, çözüm bu tablodaki ilgili sürücü/modülü güncellemek veya udev iznini düzeltmektir (Bölüm 04, sorun→çözüm).

---

<a id="4"></a>
## 4. Alıcı ve GUI Uygulamaları

Genel alıcılar, spektrumu görsel olarak gezip bir sinyali demodüle edip dinlemeni sağlayan radyo ön yüzleridir. Felsefe farkları ve "hangisi ne zaman" matrisi Bölüm 12 §3'te ayrıntılı.

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| GQRX | `gqrx-sdr/gqrx` | Linux/macOS'ta en yaygın "ilk radyo" alıcısı; GNU Radio + SoapySDR üzerine. Tek-kanal dinleme, IQ kaydı, bookmark. Bölüm 12 §3. | Olgun | GPL-3 |
| SDRangel | `f4exb/sdrangel` | "İsviçre çakısı": çoklu-cihaz/çoklu-kanal, geniş demod/eklenti yelpazesi (ADS-B, DAB, DVB-S, APT, uydu izleyici). Öğrenme eğrisi dik. | Aktif | GPL-3 |
| SDR++ | `AlexandreRouma/SDRPlusPlus` | Modern, çapraz platform (Win/Linux/macOS/Android), modüler, yüksek performanslı alıcı; band-plan şeridi, çok-VFO. | Aktif | GPL-3 |
| CubicSDR | `cjcliffe/CubicSDR` | SoapySDR üzerine kurulu, görsel olarak temiz, çapraz platform tek-kanal alıcı. Sade gezinme için iyi. Bakım yavaşlamış olabilir — teyit et. | Olgun (bakım yavaş) | GPL-2 |
| OpenWebRX | `jketterl/openwebrx` | Tarayıcıdan erişilen çok-kullanıcılı web SDR alıcısı; bir SDR'ı ağda paylaşmak. Not: orijinal `simonyiszk/openwebrx` artık `jketterl` çatalında sürüyor — güncel hat için teyit et. | Aktif | Karışık (AGPL-3 + bileşenler) |

> Pratikte: Tek istasyon dinleme için GQRX/SDR++; aynı anda çok şey izleme veya çok-cihaz için SDRangel; bir alıcıyı ağda/web'de paylaşmak için OpenWebRX. Hepsi DragonOS'ta hazır gelir (Bölüm 12 §1). Web-paylaşım araçlarında (OpenWebRX) erişimi açık internete koyuyorsan, *kimin dinlediğinin* ve *neyi yayınladığının* yasal sorumluluğu sende — RX paylaşımı bile bazı ülkelerde düzenlenir.

---

<a id="5"></a>
## 5. Çözücü ve Uygulama Araçları

Çözücüler, ham RF'i anlamlı yapısal veriye (JSON, uçak listesi, görüntü, mesaj) çevirir. Protokol prensipleri Bölüm 05'te; fiili komut/boru-hattı desenleri Bölüm 12 §5-§7'de.

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| rtl_433 | `merbanan/rtl_433` | 433/868/915 MHz ISM bandında yüzlerce IoT/sensör protokolü çözücüsü (hava istasyonu, sıcaklık/nem, TPMS, kapı zili). JSON çıktısıyla boru-hattı temeli. Kendi cihazların. Bölüm 12 §5, Bölüm 16. | Olgun | GPL-2 |
| dump1090 (FlightAware) | `flightaware/dump1090` | 1090 MHz ADS-B uçak telemetrisi çözücüsü + web harita; en yaygın FA çatalı. Açık yayın. Bölüm 05. | Olgun | BSD |
| readsb | `wiedehopf/readsb` | dump1090'ın modern, yüksek performanslı çatalı; çok-besleme, daha iyi 1090/978 işleme. `wiedehopf` hattı en aktif bakımlı sürümdür. | Aktif | Karışık (GPL/BSD bileşenler) |
| dump978 | `flightaware/dump978` | 978 MHz UAT (ABD'ye özgü ADS-B alt sistemi) çözücüsü. Bölgesel — ABD dışında genelde gereksiz. | Olgun | BSD |
| multimon-ng | `EliasOenal/multimon-ng` | Demodüle ses akışından sayısal modları çözer (POCSAG, FLEX, AFSK, DTMF). `rtl_fm \| multimon-ng` borusu. Yasal sınır: içerik çözme çoğu yerde suç; mekanik gösterim/kendi test sinyalin. Bölüm 12 §5. | Olgun | GPL-2 |
| SatDump | `SatDump/SatDump` | Uydu görüntüleme için hepsi-bir-arada: NOAA APT, Meteor LRPT, GOES HRIT; demod→senkron→hata düzeltme→görüntü→coğrafi referans. Canlı ve kayıttan. Açık yayın. Bölüm 12 §6, Bölüm 11. | Aktif | GPL-3 |
| WSJT-X | (resmî depo SourceForge'da) | FT8/FT4/JT65 zayıf-sinyal sayısal modları (RX/TX). GitHub'da çok ayna var; resmî kaynak SourceForge `wsjt`/`wsjtx`. "GitHub'da 'WSJT-X' aratıp resmî/birincil kaynağı teyit et" — yanlış aynaya güvenme. TX=amatör lisans. Bölüm 12 §7. | Olgun | GPL-3 |
| Dire Wolf | `wb2osz/direwolf` | Yazılım TNC; APRS/AX.25 paket çözücü. `rtl_fm \| direwolf` borusu, KISS/AGW portu. Açık amatör. Bölüm 12 §7. | Olgun | GPL-2 |
| Gpredict | `csete/gpredict` | Uydu geçiş takibi/yörünge tahmini (NOAA/Meteor/amatör uydu pasları). SatDump kayıt zamanlamasını planlamak için. | Olgun (bakım yavaş) | GPL-2 |
| acarsdec | `TLeconte/acarsdec` | VHF ACARS (uçak veri bağlantısı) çözücüsü. Niş ama olgun; içerik dinleme yasal sınırına dikkat (bölgesel). | Niş/uzman | GPL-2 |
| AIS çözücüler (örnek) | `dgiardini/rtl-ais` | ~162 MHz AIS gemi telemetrisi çözücüsü (RTL tabanlı). Açık yayın. Alternatifler için "GitHub'da 'AIS decoder SDR' arat". Bölüm 05. | Olgun | GPL-2 |

> Pratikte tuzaklar: (1) WSJT-X gibi araçlarda GitHub'da çok sayıda ayna/çatal görürsün; resmî birincil kaynağı teyit etmeden indirme. (2) ADS-B/AIS'te az veri görüyorsan sorun genelde anten/konum, yazılım değil (Bölüm 03). (3) multimon-ng ve ACARS gibi araçlar teknik olarak içerik çözebilir; bunun yasallığı içeriğin türüne ve ülkene bağlıdır — POCSAG/özel haberleşme içeriği çözmek çoğu yerde suçtur (Bölüm 00/06).

---

<a id="6"></a>
## 6. Analiz ve Tersine Mühendislik

Bu araçlar bilinen protokolleri değil, *bilinmeyen* bir sinyali sıfırdan anlamak içindir: kaydet, görselleştir, modülasyon/sembol hızı çöz, bitlere indir, kodlamayı kır, alanları etiketle. Yasal çizgi katıdır: yalnızca kendi cihazın. Fiili iş akışı Bölüm 12 §8'de.

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| Universal Radio Hacker (URH) | `jopohl/urh` | Uçtan uca tersine mühendislik tek pencerede: kayıt→demod→bit hizalama→kod çözme→alan etiketleme. Yeni başlayan için en bütünleşik. YALNIZCA kendi cihazın. Bölüm 12 §8. | Olgun | GPL-3 |
| Inspectrum | `miek/inspectrum` | Kaydedilmiş IQ'nun spektrogramını inceleme; imleçle sembol süresi/baud ölçme, sinyal yapısını gözle çözme. URH'nin tamamlayıcısı. Kendi kaydın. | Olgun (bakım yavaş) | GPL-3 |
| baudline | (kapalı kaynak, ücretsiz ikili) | Zaman-frekans spektral analiz aracı. AÇIK KAYNAK DEĞİL ama ekosistemde geleneksel olarak anılır; ücretsiz ikili. Açık-kaynak alternatif: Inspectrum. | — | Tescilli (ücretsiz) |
| gr-* eklentileri (genel) | `gnuradio/...` ve topluluk | GNU Radio analiz/uygulama eklenti ailesi (örn gr-satellites, gr-paint, gr-fosphor). İhtiyaca göre derlenir. Resmî/topluluk ayrımını teyit et. | Çeşitli | Genelde GPL-3 |
| gr-satellites | `daniestevez/gr-satellites` | Çok sayıda küçük/amatör uydunun telemetri çözücü GNU Radio paketi. Niş ama çok bakımlı; SatDump'ı tamamlar. | Aktif | GPL-3 |
| sigrok / PulseView | `sigrokproject/...` | Mantık analizörü/osiloskop sinyal çözme (RF değil ama yan-kanal/dijital protokol tersine müh. için sık kullanılır). Bölüm 17 ile bağlantılı. | Olgun | GPL-3 |

> Pratikte: URH + Inspectrum ikilisi, kendi 433 MHz kumandan/IoT sensörün gibi *kendi* cihazlarını çözmenin standart takımıdır. baudline gibi kapalı ama ücretsiz araçları bütünlük için andım; bu bölümün ruhu açık kaynak olduğundan, açık-kaynak muadili (Inspectrum) tercih edilmelidir. Tersine mühendislik sonucu bulduğun bir sinyali *yeniden üretmek/yayınlamak* (replay) kendi cihazın dışında her hedefte suçtur (Bölüm 06).

---

<a id="7"></a>
## 7. Hücresel Araştırma (Yasal/Akademik Sınır)

Bu kategori en hassas yasal sınıra sahiptir. Buradaki araçlar *mimariyi anlamak* (pasif gözlem), *kendi izole laboratuvarını kurmak* (Faraday/kablolu test hücresi) ve *kendini savunmak* (IMSI-catcher tespiti) için meşrudur. Kullanıcı içeriği yakalamak/çözmek, sahte baz istasyonu kurmak/yayın yapmak, kimlik/abone takibi — kapsam dışıdır ve ağır suçtur. İleri 4G/5G derinliği Bölüm 20'de; yasal çerçeve Bölüm 06/00'da.

```
   ┌───────────────────────────────────────────────────────────────────────┐
   │  HÜCRESEL ARAÇLAR — MEŞRU SINIR                                         │
   │                                                                         │
   │  ✓ MEŞRU:                                                               │
   │    - Downlink kontrol/yayın kanalı yapısını PASİF gözlemlemek (eğitim). │
   │    - PPM kalibrasyonu (downlink taşıyıcısını frekans referansı almak).  │
   │    - KENDİ izole test hücreni (srsRAN/Open5GS) Faraday/kabloda kurmak.  │
   │    - Savunma: IMSI-catcher TESPİT araçları çalıştırmak.                 │
   │                                                                         │
   │  ✗ SUÇ (anlatılmaz, yapılmaz):                                          │
   │    - Kullanıcı SES/SMS/veri trafiğini yakalamak/çözmek.                 │
   │    - Sahte baz istasyonu (IMSI-catcher) KURMAK / açık antenle yayın.    │
   │    - Başkasının cihazına/aboneliğine yönelik her işlem; abone takibi.   │
   └───────────────────────────────────────────────────────────────────────┘
```

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| gr-gsm | `ptrkrysik/gr-gsm` | GSM downlink kontrol/yayın kanalı gözlemi (`grgsm_livemon`) ve kalibrasyon. YALNIZCA downlink kontrol; içerik değil. Bölüm 12 §9. | Olgun (bakım yavaş) | GPL-3 |
| kalibrate-rtl (kal) | `steve-m/kalibrate-rtl` | GSM downlink taşıyıcısını referans alıp RTL-SDR ppm ofsetini ölçer. Pasif frekans ölçümü; içerik çözmez. Bölüm 12 §2/§9. | Olgun | BSD |
| srsRAN Project (5G) | `srsran/srsRAN_Project` | Açık kaynak 5G gNB/CU/DU yığını. YALNIZCA kendi izole test hücren (Faraday/kablolu), kendi test SIM/UE. Açık antenle yayın = suç. | Aktif | AGPL-3 |
| srsRAN 4G | `srsran/srsRAN_4G` | LTE eNB + EPC + UE yazılım yığını (eski ad srsLTE). Aynı izolasyon şartı. Bölüm 12 §9. | Olgun | AGPL-3 |
| Open5GS | `open5gs/open5gs` | Açık kaynak 4G EPC / 5G çekirdeği (core network). Test hücresinin çekirdek tarafı; srsRAN/UERANSIM ile eşlenir. İzole laboratuvar. | Aktif | AGPL-3 |
| OpenAirInterface (OAI) | (Eurecom GitLab'da barınır) | Araştırma sınıfı 4G/5G RAN+core. Birincil kaynak Eurecom GitLab'dır, GitHub değil — "OpenAirInterface gitlab eurecom" diye resmî kaynağı teyit et. İzole laboratuvar. | Aktif | OAI Public License (Apache benzeri, özel) |
| UERANSIM | `aligungr/UERANSIM` | 5G UE ve gNB *simülatörü* (RF'siz, çekirdek test için). Open5GS ile yazılım test hücresi kurmanın RF-sız yolu; en güvenli yasal seçenek. | Aktif | GPL-3 |
| YateBTS | (Yate/Legba kaynakları) | GSM baz istasyonu yazılımı. Sahte-BTS riski en yüksek araçlardandır; YALNIZCA tam izole/yetkili araştırma. Resmî kaynağı "YateBTS resmi repo" diye teyit et; rastgele çatala güvenme. | Niş/uzman | AGPL/çeşitli |

> Yasal uyarı (mutlak): Bu tablodaki ağ-tarafı araçların (srsRAN, Open5GS, OAI, YateBTS) bir baz istasyonu *yayını* başlatması, açık antenle yapıldığında hem ruhsatlı operatör spektrumuna girişimdir hem de sahte baz istasyonu niteliğine girer — her ikisi de ağır suçtur. Bu araçlar yalnızca kablolu bağlantı veya Faraday kafes içinde, kendi test SIM/UE'nle çalıştırılır. RF'e açmadan önce mutlaka izolasyonu doğrula. En güvenli öğrenme yolu UERANSIM gibi tamamen RF-sız simülasyondur. gr-gsm tarafında ise bazı ülkelerde downlink kontrol kanalı gözlemi bile gri alandadır — şüphedeysen yapma.

---

<a id="8"></a>
## 8. GNSS / GPS (Sınır Net)

GNSS araçlarında sınır kategorik olarak nettir: *alıcı/çözücü* tarafı (kendi konumunu öğrenmek, sinyal yapısını incelemek) tamamen meşrudur; *sinyal üretme/yayın* tarafı (gps-sdr-sim ile IQ üretip antenle yayınlamak) her yerde ağır suçtur (GPS spoofing). Sistem derinliği Bölüm 10'da; spoofing yasal-teknik sınırı Bölüm 06'da.

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| GNSS-SDR | `gnss-sdr/gnss-sdr` | Açık kaynak yazılım GNSS *alıcısı*: acquisition→tracking→PVT. GPS/Galileo/GLONASS/BeiDou çözer. Tamamen meşru RX/öğrenme. Bölüm 10. | Olgun | GPL-3 |
| gps-sdr-sim | `osqzss/gps-sdr-sim` | GPS L1 C/A *sinyal üreteci* (IQ dosyası üretir). YALNIZCA kablolu/Faraday izole test; ANTENLE YAYIN YASAK (spoofing = ağır suç). Bölüm 06/10. | Olgun | MIT |
| GPS-SDR-SIM benzeri Galileo/GLONASS üreteçler | (çeşitli çatallar) | Diğer takımyıldızlar için sinyal üreteçleri. Aynı mutlak kısıt: yalnızca izole test, yayın yok. Resmî/güvenilir çatalı teyit et. | Niş/uzman | Genelde MIT |
| RTKLIB | `tomojitakasu/RTKLIB` | GNSS hassas konumlama (RTK/PPP) işleme kütüphanesi/araçları. Ölçüm sonrası işleme; yayınla ilgisi yok, tamamen meşru. | Olgun (bakım yavaş; aktif çatallar var) | BSD-2 |

> Yasal uyarı (mutlak): GNSS-SDR ve RTKLIB tamamen meşrudur — biri sinyal alıp çözer, diğeri ölçüm işler; ikisi de yayın yapmaz. Buna karşılık gps-sdr-sim ve benzeri *üreteçler* bir IQ dosyası üretir; bu dosyayı bir SDR ile antenden yayınlamak GPS spoofing'dir ve her yerde ağır suçtur (uçak/gemi/altyapı navigasyonunu tehlikeye atar). Bu araçların meşru tek kullanımı, *kendi* GNSS alıcını kablolu/Faraday izole ortamda test etmektir. Antene asla bağlama.

---

<a id="9"></a>
## 9. Yön Bulma ve Konum Tespiti

Yön bulma (DF), bir sinyalin *nereden geldiğini* faz/zaman farkından kestirir; çok-kanallı tutarlı (coherent) alıcı gerektirir. Pasif gözlemdir (yayın yapmaz). Donanım Bölüm 02'de, anten/faz Bölüm 03'te, takip kuramı Bölüm 09'da.

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| KrakenSDR DOA | `krakenrf/krakensdr_doa` | KrakenSDR (5-kanal tutarlı RTL) için varış-yönü (DOA) kestirim yazılımı; harita/triangülasyon. Açık kaynak/açık donanım DF'in fiili referansı. Bölüm 09, Bölüm 12 §11. | Aktif | GPL-3 |
| KrakenSDR firmware/araçlar | `krakenrf/krakensdr_docs` ve ilgili | Belgeler, firmware ve yardımcı araçlar. Donanım+yazılım birlikte; depo ailesini "GitHub'da 'krakenrf' organizasyonu" altında teyit et. | Aktif | Karışık |
| KerberosSDR (öncül) | (krakenrf öncesi proje) | Kraken'in öncülü; 4-kanal tutarlı DF. Büyük ölçüde KrakenSDR'a evrildi. Yeni kurulum için Kraken hattını tercih et; Kerberos arşiv niteliğinde. | Arşiv/durağan | GPL (teyit et) |

> Pratikte: Açık kaynak DF dünyası fiilen KrakenSDR ekosistemine (krakenrf organizasyonu) toplanmıştır. KerberosSDR adını eski belgelerde görürsün; yeni iş için Kraken hattı geçerlidir. DF pasif bir ölçümdür — yayın yapmaz, içerik çözmez; yalnızca sinyalin geliş yönünü kestirir. Çok-kanallı tutarlı yakalama için bare-metal kurulum önerilir (Bölüm 12 §1, VM-USB tuzağı).

---

<a id="10"></a>
## 10. WiFi / WLAN Güvenliği (Yetkili Pentest)

Bu kategori, *yalnızca yazılı yetki verilmiş kapsamda* veya *kendi ağında* meşrudur. Başkasının ağına izinsiz erişim, deauth ile servis kesme, el sıkışma yakalayıp kırma — yetkisiz yapıldığında suçtur. Bu araçlar savunma testi, kendi ağını sertleştirme ve yetkili kırmızı-takım işleri içindir. WiFi/WLAN güvenliği derinliği Bölüm 15'te; yasal çerçeve Bölüm 00/06'da.

```
   ┌───────────────────────────────────────────────────────────────────────┐
   │  WIFI/KABLOSUZ GÜVENLİK ARAÇLARI — MEŞRU SINIR                          │
   │                                                                         │
   │  ✓ MEŞRU:  Kendi ağın; YAZILI yetki verilmiş pentest kapsamı;          │
   │            izole laboratuvar/test SSID; savunma izleme (Kismet).        │
   │  ✗ SUÇ:    Başkasının ağına izinsiz erişim; yetkisiz deauth/DoS;        │
   │            izinsiz el-sıkışma yakalama/kırma; komşu ağ takibi.          │
   └───────────────────────────────────────────────────────────────────────┘
```

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| Aircrack-ng | `aircrack-ng/aircrack-ng` | WiFi denetim paketi: yakalama (airodump-ng), enjeksiyon, WPA/WPA2 el-sıkışma kırma. Kablosuz pentest'in fiili referansı. Yetkili kapsam. Bölüm 15. | Referans | GPL-2 |
| hcxdumptool | `ZerBea/hcxdumptool` | WPA/WPA2 PMKID/el-sıkışma yakalama (modern, verimli). Yalnızca yetkili/kendi ağ. hcxtools ile eşlenir. | Aktif | MIT |
| hcxtools | `ZerBea/hcxtools` | hcxdumptool çıktısını hashcat formatına dönüştürme/işleme araçları. Yakalama sonrası analiz. | Aktif | MIT |
| Bettercap | `bettercap/bettercap` | Ağ/kablosuz saldırı-çerçevesi (MITM, WiFi/BLE/HID modülleri). Çok yetenekli → yetkisiz kullanımı ağır suç. Yalnızca yetkili pentest/laboratuvar. | Aktif | GPL-3 |
| Kismet | `kismetwireless/kismet` | Kablosuz keşif/izleme/IDS (WiFi/BLE/diğer); pasif algılama. Savunma izleme için en meşru kullanım. Bölüm 15. | Aktif | GPL-2 |
| Nexmon | `seemoo-lab/nexmon` | Broadcom WiFi yonga firmware yamalama çatası (monitör mod/enjeksiyon, örn Pi/telefon yongaları). Araştırma sınıfı; donanıma özel. | Aktif | Karışık (GPL + tescilli kısımlar) |
| Pwnagotchi | `evilsocket/pwnagotchi` | Öğrenen, otonom WiFi el-sıkışma/PMKID toplayıcı (Pi tabanlı). Yetkisiz toplama suçtur — yalnızca kendi laboratuvar/yetkili kapsam. | Aktif (topluluk çatalları da var) | GPL-3 |
| esp8266_deauther | `SpacehuhnTech/esp8266_deauther` | ESP8266 ile deauth/test cihazı; eğitim/farkındalık amaçlı. Deauth bir DoS'tur — yetkisiz kullanımı suç. Yalnızca kendi/izole test. | Aktif | MIT (teyit et) |
| Wifite | `kimocoder/wifite2` | Aircrack-ng/hcx araçlarını otomatize eden sarmalayıcı. Birincil hattı "GitHub'da 'wifite2' aratıp güncel bakımlı çatalı" teyit et. Yetkili kapsam. | Aktif | GPL-2 |

> Yasal uyarı (mutlak): Bu tablodaki araçların neredeyse tamamı *aktif saldırı* yapabilir (deauth/DoS, el-sıkışma yakalama, MITM). Yetkisiz hedefte kullanıldığında bunlar yasal "test" değil, bilişim suçudur. Tek meşru zemin: (a) kendi ağın/donanımın, (b) yazılı yetki ve tanımlı kapsamı olan bir pentest, (c) RF-izole laboratuvar. Kismet gibi *pasif izleme* araçları en geniş meşru alana sahiptir (savunma). Şüphedeysen, yazılı yetki olmadan paket enjekte etme veya deauth gönderme.

---

<a id="11"></a>
## 11. BLE / RFID / IoT Donanım Araçları

Bu kategori kısa-menzil kablosuz (Bluetooth/BLE, RFID/NFC, Zigbee, alt-GHz IoT) güvenlik ve araştırma araçlarını kapsar. Çoğu kendi cihazını/etiketsini incelemek için meşrudur; başkasının kartını/cihazını klonlamak/yetkisiz erişim suçtur. Kısa-menzil derinliği Bölüm 16'da.

| Araç | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| Sniffle | `nccgroup/Sniffle` | BLE5 sniffer (TI uyumlu donanımla); bağlantı takibi dahil güçlü BLE yakalama. Araştırma/kendi cihazın. Bölüm 16. | Aktif | GPL-3 (teyit et) |
| crackle | `mikeryan/crackle` | BLE eski eşleşme (legacy pairing) şifre çözme aracı (yakalanmış el-sıkışmadan). Yalnızca kendi cihazların/araştırma. | Olgun (bakım yavaş) | BSD-2 |
| Proxmark3 | `RfidResearchGroup/proxmark3` | RFID/NFC araştırma platformunun fiili referans firmware/istemcisi (RRG/Iceman çatalı). Kendi etiketlerin/erişim kartların. Yetkisiz klonlama suç. Bölüm 16. | Aktif | GPL-3 |
| KillerBee | `riverloopsec/killerbee` | Zigbee/IEEE 802.15.4 değerlendirme çatası. Niş; donanıma bağımlı (örn ApiMote). Araştırma. | Niş/uzman | GPL (teyit et) |
| Flipper Zero firmware | `flipperdevices/flipperzero-firmware` | Çok-protokollü cep aracı resmî firmware (alt-GHz, NFC, RFID, IR, GPIO). Kendi cihazların. Topluluk çatalları (örn Unleashed/RogueMaster) yetenekleri genişletir ama yasal sınırı genişletmez. | Aktif | GPL-3 |
| Mirage | (resmî depo teyit edilmeli) | BLE/Zigbee/diğer için modüler kablosuz saldırı/araştırma çatası (akademik kökenli). "GitHub/GitLab'da 'Mirage wireless framework' aratıp resmî kaynağı teyit et." | Niş/uzman | Genelde GPL |
| ubertooth | `greatscottgadgets/ubertooth` | Bluetooth/BLE araştırma donanımı host araçları (Ubertooth One). BLE keşif/yakalama; açık donanım+yazılım. | Olgun | GPL-2 |

> Pratikte: Bu araçların meşru zemini "kendi cihazın/kartın/etiketin"dir. Proxmark3 ile kendi erişim kartını okumak/anlamak eğitseldir; başkasının kartını klonlamak (fiziksel erişim suistimali) suçtur. Flipper firmware'inin topluluk çatalları daha fazla *teknik* yetenek açar; ama *yasal* sınırı değiştirmez — yetkisiz kullanım yine suçtur (Bölüm 16, Bölüm 00).

---

<a id="12"></a>
## 12. Veri Kümeleri ve Makine Öğrenmesi

RF/SIGINT'e makine öğrenmesi uygulamak (modülasyon sınıflandırma, sinyal tespiti, anormallik) için veri kümeleri ve ML kütüphaneleri. Derinlik Bölüm 19'da.

| Kaynak | Depo / yer | Görev / kapsam | Olgunluk | Lisans/erişim |
|---|---|---|---|---|
| torchsig | `TorchDSP/torchsig` | PyTorch tabanlı SIGINT/RF makine öğrenmesi kütüphanesi: sentetik veri üretimi, modülasyon sınıflandırma veri kümeleri, modeller. Bölüm 19. | Aktif | MIT (teyit et) |
| RadioML / DeepSig veri kümeleri | (DeepSig sitesi) | Modülasyon sınıflandırma için klasik RadioML 2016/2018 veri kümeleri. GitHub'da değil, DeepSig dağıtım sayfasından indirilir — "DeepSig RadioML dataset" diye resmî kaynağı bul. Akademik kullanım koşullarına dikkat. | Referans (veri) | DeepSig koşulları (teyit et) |
| gr-inference | (resmî depo teyit edilmeli) | GNU Radio içinde sinir-ağı çıkarımı (inference) çalıştırma eklentisi. "GitHub'da 'gr-inference' aratıp resmî/güncel depoyu teyit et." | Niş/uzman | Genelde GPL-3 |
| SatNOGS veri/gözlem | `satnogs/...` | Topluluk uydu yer-istasyonu ağı; gözlem verisi ve yazılım. ML için gerçek uydu telemetri/IQ kaynağı olabilir. §13'te de anılır. | Aktif | AGPL/karışık |

> Pratikte: RF-ML için en hızlı başlangıç, sentetik veri (torchsig) + bir referans veri kümesidir (RadioML). Gerçek-dünya verisi için kendi IQ kayıtların (Bölüm 12 §13) veya SatNOGS/topluluk gözlemleri kullanılır. Veri kümelerinin *lisans/kullanım koşulları* GitHub kodundan farklı olabilir (özellikle akademik-yalnız); indirmeden önce koşulu oku. RadioML gibi kümeler GitHub'da değil üreticinin sitesinde dağıtılır.

---

<a id="13"></a>
## 13. Veritabanı, Topluluk ve Çevrimiçi Kaynaklar

Bunlar yazılım deposu değil, *bilgi ve veri* kaynaklarıdır: bir sinyali tanımlamak, bir frekansın kime ait olduğunu bulmak, topluluk gözlemlerine erişmek için. İstihbarat-kaynakları/takip yöntemi Bölüm 14'te; band planı Bölüm 08'de.

| Kaynak | Yer / erişim | Görev / kapsam | Olgunluk | Not |
|---|---|---|---|---|
| Signal Identification Guide (sigidwiki) | `sigidwiki.com` | Bilinmeyen bir sinyalin "sesine/şekline" bakıp ne olduğunu tanımlama wiki'si (ses örnekleri + waterfall). Sinyal ayıklamanın (Bölüm 07) sözlüğü. | Topluluk referansı | Web; içerik lisansını teyit et |
| RadioReference | `radioreference.com` | Frekans tahsisi/veritabanı (özellikle ABD; bölgesel kapsam değişir). "Bu frekans kimin?" sorusu. Bölüm 08. | Olgun (veritabanı) | Bazı içerik üyelik gerektirir |
| SatNOGS | `satnogs/...` ve `satnogs.org` | Açık kaynak uydu yer-istasyonu ağı: gözlem zamanlama, telemetri/IQ verisi, açık donanım istasyon tasarımları. | Aktif | AGPL/karışık; açık veri |
| ADS-B Exchange | `adsbexchange.com` | Filtrelenmemiş topluluk ADS-B uçak verisi (besleme/erişim). Kendi dump1090/readsb beslemeni buraya verebilirsin. | Olgun (topluluk) | Erişim/kullanım koşulu değişti — teyit et |
| GPSJam | `gpsjam.org` | Günlük GPS girişim/parazitlenme haritası (uçak verisinden türetilmiş). GNSS girişimi farkındalığı (Bölüm 06/10). | Aktif (görselleştirme) | Web; veri kaynağı türev |
| "awesome-sdr" derlemeleri | (GitHub'da "awesome sdr" arat) | Topluluk-derlemesi bağlantı listeleri (araç/öğrenme/donanım). Tek bir resmî liste yoktur; birden çok "awesome-*" deposu var — yıldız/güncellik bak. | Çeşitli | Genelde CC/MIT |

> Pratikte: Bilmediğin bir sinyali bulduğunda iş akışı: önce sigidwiki ile *tanımla* (ne bu?), sonra RadioReference/band planı ile *kime ait* (Bölüm 08), gerekirse SatNOGS/topluluk ile *gözlem* paylaş. "awesome-sdr" türü derlemeler iyi bir başlangıç haritasıdır ama tek bir resmî liste yoktur; birden fazlasını tarayıp güncel/bakımlı olanı seç (§16 kriterleri). Bu kaynakların çoğu web sitesidir; veri/içerik kullanım koşulları GitHub lisansından ayrıdır.

---

<a id="14"></a>
## 14. Kurulum Dağıtımları (Hazır OS İmajları)

Bu projeler tek bir araç değil, *önceden derlenmiş araç ekosistemini* paketleyen Linux dağıtımlarıdır. Değerleri, bağımlılık cehennemini (Bölüm 04) çözüp "indir-aç-çalıştır" mesafesi sunmalarıdır. DragonOS'un felsefesi/sınırları Bölüm 12 §1'de.

| Dağıtım | Yer / kaynak | Görev / kapsam | Olgunluk | Not |
|---|---|---|---|---|
| DragonOS | (resmî dağıtım/SourceForge + cemaxecuter kanalları) | "Her şey kurulu" SDR dağıtımı (FocalX masaüstü, Pi/ARM hatları). GNU Radio, GQRX/SDR++/SDRangel, rtl_433, SatDump, URH, KrakenSDR vb. hazır. İmaj kaynağını resmî sayfadan teyit et. Bölüm 12 §1. | Aktif | Karışık (Ubuntu + araçların kendi lisansları) |
| Skywave Linux | (resmî dağıtım sayfası) | SDR/kısa-dalga odaklı canlı Linux dağıtımı; alıcı/çözücü araçları hazır. DragonOS'a alternatif. "Skywave Linux resmi" diye kaynağı teyit et. | Niş/uzman | Karışık |
| Kali Linux (ilgili paketler) | `kali.org` / paket depoları | Genel pentest dağıtımı; aircrack-ng/kismet/bettercap gibi kablosuz-güvenlik araçları paketli (ama SDR-uydu odaklı değil). WiFi pentest tarafı için pratik. | Referans (güvenlik dağıtımı) | Karışık |

> Pratikte: SDR/uydu/sinyal işi için DragonOS, kablosuz-güvenlik/pentest işi için Kali tarafı pratik başlangıçlardır; ikisi farklı amaca hizmet eder. Hazır imajların *sınırı*, içlerindeki araç sürümlerinin imaj çıktığı andaki sürümler olmasıdır — yeni özellik için ya imajı güncellersin ya aracı elle. Üretim istasyonu çoğu zaman kendi Ubuntu/Debian kurulumuna geçer (Bölüm 04 yol haritası). İmajları daima *resmî/bilinen* kaynaktan indir; rastgele ayna güvenlik riskidir.

---

<a id="15"></a>
## 15. Donanım Açık-Tasarım Depoları

Açık kaynak yalnızca yazılım değildir; bazı SDR/güvenlik donanımları şematik, PCB ve firmware'iyle birlikte açıktır. Bu, donanımı anlamak, üretmek/onarmak ve güvenmek (kapalı kutu olmaması) açısından değerlidir. Donanım derinliği Bölüm 02/03'te.

| Donanım | Depo (`org/repo`) | Görev / kapsam | Olgunluk | Lisans |
|---|---|---|---|---|
| HackRF One | `greatscottgadgets/hackrf` | Şematik/PCB/firmware/host yazılımı tek depoda; geniş-bant açık SDR'ın referans açık tasarımı. | Referans | Donanım açık + GPL yazılım |
| LimeSDR | `myriadrf/LimeSDR-USB` (ve `LimeSDR-Mini`) | LimeSDR donanım tasarım dosyaları (şematik/PCB). MyriadRF organizasyonu altında; tam depo adını teyit et. | Olgun | Açık donanım (CERN OHL vb.; teyit et) |
| KrakenSDR | `krakenrf/...` | 5-kanal tutarlı RTL DF donanımı; belge/tasarım/yazılım depoları krakenrf organizasyonunda. Donanımın açıklık derecesini depoda teyit et. | Aktif | Karışık |
| Osmocom projeleri (genel) | `osmocom/...` | Açık telekom/SDR ekosistemi (rtl-sdr, OsmoTRX, çeşitli GSM/telecom altyapı). Organizasyon geniştir; ilgili alt-proje deposunu teyit et. | Olgun | Genelde GPL |
| Ubertooth One | `greatscottgadgets/ubertooth` | Bluetooth/BLE araştırma donanımı; açık şematik/firmware/host. | Olgun | Açık donanım + GPL |

> Pratikte: Açık donanım depoları, bir cihazın *nasıl çalıştığını* öğrenmenin ve kapalı-kutu olmadığını doğrulamanın en saf yoludur. Great Scott Gadgets (HackRF/Ubertooth) ve MyriadRF (LimeSDR) bu felsefenin güçlü örnekleridir. Donanımı kendin üretmesen bile, açık tasarımı okumak Bölüm 02/03'teki kavramları (RF ön-uç, kazanç katmanları, saat) somutlaştırır. Bir deponun "açık donanım" iddiasını LICENSE/hardware klasöründen teyit et; "açık kaynak yazılım" ile "açık donanım" farklı lisanslardır.

---

<a id="16"></a>
## 16. Bir GitHub Projesi Nasıl Değerlendirilir

Bu dizindeki "olgunluk" sütunu bir başlangıçtır; ama açık kaynak hızla değişir. Bir aracı kullanmaya karar vermeden önce, deposunun *bugünkü* canlılığını kendin ölçmelisin. Aşağıdaki sinyaller, "bu projeye dayanabilir miyim" sorusunu yanıtlar.

### Canlılık ve bakım sinyalleri

```
   PROJE CANLILIK KONTROL LİSTESİ
   ─────────────────────────────────────────────────────────────────
   1) Son commit tarihi    → Aylar/yıllar önce mi? Durağan olabilir.
                             (Ama "tamamlanmış" araçlar da seyrek
                              commit alır; tek başına kanıt değil.)
   2) Sürüm/release sıklığı → Etiketli sürümler var mı, düzenli mi?
   3) Açık vs kapalı issue  → Yüzlerce açık, hiç kapanmayan issue =
                              terk edilme işareti. Sağlıklı oran ve
                              geliştirici yanıtı önemli.
   4) Katkıda bulunan sayısı → Tek kişi mi (otobüs riski) yoksa ekip mi?
   5) Yıldız/fork           → Popülerlik ipucu; ama YILDIZ KALİTE DEĞİL.
                              Çok yıldızlı ama terk edilmiş projeler var.
   6) Belge/README kalitesi → Kurulum, bağımlılık, örnek var mı?
   7) CI/test durumu        → Yeşil CI rozeti, test = bakım disiplini.
   8) Resmî kaynak teyidi   → Proje sitesi/wiki, deponun "resmî" olduğunu
                              doğruluyor mu? (Çatal tuzağı için §17.)
```

> Pratikte: Hiçbir sinyal tek başına yeterli değildir. "Son commit eski" bir aracın *çalışmadığı* anlamına gelmez — olgun, tamamlanmış bir kütüphane (örn kararlı bir çözücü) yıllarca commit almadan mükemmel çalışabilir. Tersine, çok yıldızlı bir proje terk edilmiş ve derlenmez olabilir. Sinyalleri *birlikte* oku: aktif issue yanıtı + son sürüm + makul katkıda bulunan = güvenli; yüzlerce yanıtsız issue + tek geliştirici + iki yıl sessizlik = dikkat.

### Yıldız sayısı neden yanıltıcıdır

Yıldız (star), bir GitHub kullanıcısının "ilginç/yer imi" işaretidir; *kullanıyorum* veya *bakımlı* demek değildir. Bir araç bir konferans sunumuyla bir hafta içinde on bin yıldız toplayıp sonra terk edilebilir. Tersine, kritik bir kütüphane (örn bir SoapySDR modülü) az yıldızla ama sağlam bakımla yıllarca taşınabilir. Yıldızı "dikkat çekti" olarak oku, "güvenilir" olarak değil. Asıl ölçüt bakım disiplinidir (issue yanıtı, sürüm, CI).

### Olgunluk-görev eşlemesi

Bir aracın olgunluğu, *görevin kritikliğiyle* eşleşmelidir:

| Görevin kritikliği | Kabul edilebilir olgunluk |
|---|---|
| Üretim/saha istasyonu (güvenilmeli) | Referans/Olgun |
| Düzenli kullanım, tolere edilebilir aksama | Aktif |
| Tek seferlik deneme/öğrenme | Niş/uzman dahil her şey |
| Tarihsel/karşılaştırma | Arşiv bile olur (ama çatal teyidiyle) |

> Sonuç: Kritik bir görev için "Aktif ama tek-geliştirici, hızlı değişen" bir aracı temel alma — yarın API'si değişebilir veya terk edilebilir. Öğrenme/deneme için ise niş/arşiv araçlar bile değerlidir. Olgunluğu göreve göre seç.

---

<a id="17"></a>
## 17. Çatal Tuzakları, Derleme/Bağımlılık ve Güvenlik

Açık kaynak indirmenin en gerçek riskleri burada toplanır: yanlış/kötü amaçlı çatalı klonlamak, bağımlılık cehenneminde boğulmak ve derleme sırasında güvenliği ihlal etmek.

### Çatal (fork) tuzakları

GitHub'da bir aracı arattığında çoğu zaman *birden çok* depo görürsün: orijinal, resmî çatal, kişisel çatallar ve bazen kötü amaçlı taklitler. Yanlışını klonlamak en az zaman kaybı, en kötü güvenlik ihlalidir.

```
   ÇATAL AYIRT ETME
   ─────────────────────────────────────────────────────────────────
   - RESMÎ KAYNAK NEREDE?   Proje sitesi/wiki hangi depoyu gösteriyor?
                            (Bu dizin de teyit notlarıyla yardımcı olur.)
   - HANGİ ÇATAL AKTİF?     Bazı araçlarda geliştirme orijinalden bir
                            çatala TAŞINMIŞTIR (örn ADS-B'de readsb
                            wiedehopf hattı; rtl-sdr'da blog çatalı;
                            Proxmark3'te RRG/Iceman). Resmî/aktif hattı
                            kullan, eski orijinali değil.
   - TAKLİT/KÖTÜ AMAÇLI?    İsmi neredeyse aynı ama sahibi tanınmayan,
                            ani oluşturulmuş, tek commit'lik depolar
                            şüphelidir. İndirme betiği/binary ekleyen
                            "kolaylık" çatallarına dikkat.
   - SÜRÜM ETİKETİ?         Resmî projeler genelde imzalı/etiketli
                            sürümler yayınlar; rastgele çatalda yoktur.
```

> Kritik güvenlik uyarısı: Bu dizinde bazı araçlar için bilerek "resmî çatal" ile "orijinal" ayrımı yaptım (örn rtl-sdr osmocom vs blog, Proxmark3 RRG/Iceman). Sebep tam budur: bazı ekosistemlerde *aktif geliştirme orijinal depodan bir çatala taşınmıştır* ve doğru çatalı kullanmak hem işlevsel hem güvenlik açısından önemlidir. Bir depo yolunu körlemesine kullanmak yerine, projenin resmî sitesinden/bilinen topluluğundan doğru hattı teyit et. Adından emin olmadığım her yerde "GitHub'da '<isim>' aratıp resmî depoyu teyit et" notu bıraktım — bu notu ciddiye al.

### Derleme ve bağımlılık tuzakları

SDR ekosisteminin acı gerçeği, araçların kendisi değil onları kurmaktır (Bölüm 04, Bölüm 12 §1). Kaynaktan derlerken tipik tuzaklar:

```
   DERLEME/BAĞIMLILIK TUZAKLARI
   ─────────────────────────────────────────────────────────────────
   - SÜRÜM BAĞIMLILIĞI   gr-* eklentisi belirli bir GNU Radio sürümüyle
                         derlenmiş olmalı; uyumsuzluk = derlenmez/çöker.
   - SİSTEM PAKETLERİ    libusb, libvolk, boost, fftw, qt — eksik geliştirme
                         (-dev/-devel) paketleri cmake'i durdurur.
   - PYTHON ORTAMI       Sistem Python'ı kirletmemek için sanal ortam;
                         yanlış ortam = "modül bulunamadı".
   - PARALEL KURULUM     Çok aracı uyumlu kurmak en zor kısım → bu yüzden
                         DragonOS gibi dağıtımlar var (önceden çözülmüş).
   - DOĞRU ÇATAL         Yanlış çatalı derlemek, doğru bağımlılıkla bile
                         beklenenden farklı davranabilir.
```

### Kaynaktan derlemenin güvenlik hijyeni

Bir depoyu derlemek, onun kodunu (ve build betiklerini) makinende çalıştırmaktır. Açık kaynak "görülebilir" demektir, "güvenli" değil. Minimum hijyen:

```
   - KAYNAĞI TEYİT ET    Resmî/aktif depo mu? (§16/§17 üstü.)
   - BUILD BETİĞİNE BAK  Kurulum betiği uzak bir yerden binary indirip
                         çalıştırıyor mu? "curl | sudo bash" desenine dikkat.
   - İZOLE DENE          Güvenmediğin bir aracı önce VM/konteynerde dene.
   - İMZA/CHECKSUM       Yayınlanmış sürümlerde sağlanıyorsa doğrula.
   - EN AZ AYRICALIK     root gerekmeyen şeyi root'la çalıştırma; udev
                         kuralı gibi şeyleri bilerek ve sınırlı ekle.
```

> Sonuç: Açık kaynak güveni "kodu görebiliyorum" ile başlar ama orada bitmez. Doğru çatalı seçmek (§17 üstü), bağımlılığı yönetmek (mümkünse hazır dağıtım kullanmak) ve build hijyeni (betiğe bakmak, izole denemek) üçlüsü, "indirdiğim araç beni güvenlik açığına sokmasın" güvencesinin pratik halidir. Şüpheli/aniden-oluşmuş bir çatalı, cazip görünse de klonlama; resmî hattı bekle.

---

<a id="18"></a>
## 18. Alıştırmalar (Yasal)

Aşağıdaki üç alıştırma, bu dizini *kullanma* becerisini pekiştirir: kendi görevine göre bir araç matrisi kurmak, bir aracı kaynaktan derlemek ve bir topluluk derlemesini eleştirel taramak. Hepsi yasaldır ve yalnızca açık/kendi kaynaklarınla yapılır.

### Alıştırma 1 — Kendi görev listene göre araç matrisi kur

```
HEDEF: Kişisel/eğitim hedeflerine göre, bu dizinden kendi araç matrisini çıkarmak.
ADIMLAR:
  1) Yapmak istediğin 4-5 YASAL görevi yaz (örn: "kendi IoT sensörlerimi oku",
     "NOAA uydu görüntüsü al", "kendi 433 MHz kumandamı çöz", "ADS-B uçak takip",
     "kendi WiFi ağımı sertleştirme testi").
  2) Her görev için bu dizinden birincil + alternatif aracı seç (ilgili tabloya bak).
  3) Her seçim için: depo yolunu not et, olgunluğu §16 kontrol listesiyle DOĞRULA
     (son commit/sürüm/issue), lisans notunu yaz, yasal sınırı işaretle.
  4) Sonucu Bölüm 12 §14'teki "görev→araç" matrisiyle karşılaştır; örtüşüyor mu?
ÖĞRENME: Dizini "okumak" değil "kullanmak"; olgunluk-görev eşlemesi (§16).
DOĞRULAMA: 4-5 satırlık, depo+olgunluk+lisans+yasal-not içeren kişisel matris.
SINIR: Yalnızca yasal görevler; her satırın yasal-not sütununu doldurmadan bitirme.
```

### Alıştırma 2 — Bir aracı kaynaktan derle

```
HEDEF: Hazır paketten değil, kaynaktan derleyerek bağımlılık zincirini görmek.
ADIMLAR:
  1) Görece sade, açık-kaynak bir araç seç (örn rtl_433: merbanan/rtl_433).
  2) §16/§17 ile RESMÎ depoyu teyit et; doğru çatal olduğundan emin ol.
  3) README'deki bağımlılıkları kur (örn libusb, librtlsdr geliştirme paketleri).
  4) cmake/make ile derle; eksik bağımlılık hatası çıkarsa hangi -dev paketinin
     eksik olduğunu çöz ve kur (Bölüm 12 §1 bağımlılık mantığı).
  5) Derlenen aracı kendi cihazınla çalıştırıp doğrula (örn kendi sensörünü oku).
ÖĞRENME: "Araçların kendisi değil kurmaktır zor" gerçeği; -dev paketi/sürüm uyumu.
DOĞRULAMA: Kaynaktan derlenmiş ikili, kendi cihazınla çalışan bir çıktı üretir.
GÜVENLİK: Build betiğine §17 hijyeniyle bak; "curl | sudo bash" desenine takılma.
```

### Alıştırma 3 — Bir "awesome-sdr" listesini eleştirel tara

```
HEDEF: Topluluk derlemesini körlemesine değil, kalite süzgeciyle okumak.
ADIMLAR:
  1) GitHub'da "awesome sdr" arat; birden çok derleme çıkacak — en güncel/bakımlı
     görüneni seç (§16: son commit, yıldız değil bakım).
  2) Listeden bu dizinde OLMAYAN 3 araç seç.
  3) Her biri için §16 canlılık kontrolünü uygula: resmî mi, aktif mi, lisans ne,
     hangi görevde ne kadar olgun?
  4) Üçünden hangisinin gerçekten "güvenilebilir", hangisinin "dikkatli dene",
     hangisinin "terk edilmiş/atla" olduğunu gerekçesiyle sınıfla.
ÖĞRENME: "awesome" listeleri başlangıç haritasıdır, onay mührü değil; eleştirel
         değerlendirme (§16) ve çatal teyidi (§17) her zaman senin işin.
DOĞRULAMA: 3 araç için olgunluk-sınıfı + gerekçe + yasal-not içeren kısa değerlendirme.
SINIR: Listede yasadışı/exploit-odaklı bir araç görürsen onu DEĞERLENDİRME DIŞI bırak;
       bu kitap yalnızca meşru, açık-kaynak araçları kapsar.
```

> Pedagoji: Üç alıştırma, dizinin üç farklı becerisini kapsar — seçme (matris), inşa etme (derleme) ve eleştirel okuma (awesome tarama). Hepsi yasaldır ve hepsi senin değerlendirme kasını çalıştırır. Asıl ders şudur: bu dizin bir başlangıç noktasıdır, son söz değil; her aracın *bugünkü* durumunu (§16) ve doğru hattını (§17) doğrulamak daima senin sorumluluğundur.

---

<a id="19"></a>
## 19. Hızlı Referans ve Diğer Bölümler

### Kategori → bölüm haritası

| Dizin kategorisi (bu bölüm) | Temel araçlar | Derinlemesine bölüm |
|---|---|---|
| SDR çerçeveleri (§2) | GNU Radio, SoapySDR, gr-osmosdr | Bölüm 12 §4, Bölüm 04 |
| Sürücüler (§3) | rtl-sdr, hackrf, LimeSuite, UHD | Bölüm 04, Bölüm 02 |
| Alıcı/GUI (§4) | GQRX, SDRangel, SDR++, OpenWebRX | Bölüm 12 §3 |
| Çözücüler (§5) | rtl_433, dump1090/readsb, SatDump, Dire Wolf | Bölüm 12 §5-§7, Bölüm 05 |
| Tersine müh. (§6) | URH, Inspectrum, gr-satellites | Bölüm 12 §8 |
| Hücresel (§7) | gr-gsm, srsRAN, Open5GS, UERANSIM | Bölüm 20, Bölüm 06 |
| GNSS (§8) | GNSS-SDR, gps-sdr-sim (sınır), RTKLIB | Bölüm 10, Bölüm 06 |
| Yön bulma (§9) | KrakenSDR DOA | Bölüm 09, Bölüm 12 §11 |
| WiFi (§10) | aircrack-ng, hcxdumptool, Kismet, bettercap | Bölüm 15, Bölüm 06 |
| BLE/RFID/IoT (§11) | Sniffle, Proxmark3, Flipper fw, ubertooth | Bölüm 16 |
| Veri/ML (§12) | torchsig, RadioML, SatNOGS | Bölüm 19 |
| Topluluk kaynak (§13) | sigidwiki, RadioReference, SatNOGS, awesome-sdr | Bölüm 14, Bölüm 08 |
| Dağıtımlar (§14) | DragonOS, Skywave, Kali | Bölüm 12 §1, Bölüm 04 |
| Açık donanım (§15) | HackRF, LimeSDR, KrakenSDR, osmocom | Bölüm 02, Bölüm 03 |

### Değerlendirme hızlı referansı

```
PROJE SEÇERKEN:    son commit + sürüm + issue yanıtı + katkıda bulunan
                   (YILDIZ kalite değildir) → §16
ÇATAL SEÇERKEN:    resmî kaynağı teyit et; aktif hat orijinalden TAŞINMIŞ olabilir
                   (readsb/wiedehopf, rtl-sdr/blog, proxmark3/RRG) → §17
DERLERKEN:         -dev paketleri + GNU Radio sürüm uyumu + sanal ortam;
                   zorsa hazır dağıtım (DragonOS) → §17, Bölüm 04
GÜVENLİK:          build betiğine bak ("curl|sudo bash" reddet), izole dene,
                   checksum/imza doğrula → §17
LİSANS:            kişisel kullanım ≠ dağıtım/servis; GPL/AGPL yükümlülük → §1
YASAL:             açık kaynak ≠ yasal kullanım; RX geniş, TX/içerik/pentest
                   ülkeye-yetkiye-hedefe bağlı → Bölüm 00/06
```

### Diğer bölümler

| # | Bölüm | İlişki |
|---|---|---|
| 0 | Başlangıç, Index, Yasal | Tüm araçların yasal çizgisi burada tanımlı |
| 4 | Yazılım, OS & Kurulum | Bağımlılık/sürücü/derleme; bu dizinin kurulum tarafı |
| 5 | Protokoller & Çözümleme | §5 çözücülerinin protokol prensibi |
| 6 | Güvenlik, Açıklar & Savunma | Hücresel/GNSS/WiFi araçlarının yasal-teknik sınırı |
| 8 | Frekans Tahsisi & Band Planı | §13 kaynaklarının "ne nereye ait" tarafı |
| 9 | Yer Tespiti & Yön Bulma | §9 DF araçlarının kuramı |
| 10 | GNSS/GPS Sistemleri | §8 GNSS araçlarının sistem tarafı |
| 12 | DragonOS & Araç Ekosistemi | §2-§9 araçlarının FİİLİ kullanımı (komut/iş akışı) |
| 14 | İstihbarat Kaynakları & Takip | §13 topluluk/veritabanı kaynaklarının yöntemi |
| 15 | WiFi/WLAN Güvenliği | §10 WiFi araçlarının derinliği |
| 16 | Kısa Menzil Kablosuz & IoT | §11 BLE/RFID/IoT araçlarının derinliği |
| 19 | Yapay Zeka & ML SIGINT | §12 veri/ML kaynaklarının uygulaması |
| 20 | İleri Hücresel 4G/5G | §7 hücresel araçlarının derinliği |

> Kapanış: Bu bölüm bir *dizindi* — hangi açık-kaynak araç, hangi resmî depo, ne kadar olgun, hangi lisans, hangi görevde ne ölçüde. Ama dizin yalnızca bir haritadır; her aracın *fiili* kullanımı (komut, akış, tuzak) Bölüm 12'de, *fiziği* Bölüm 01-03'te, *yasal sınırı* Bölüm 00/06'dadır. Bir araç ararken buraya bak, doğru hattı (§16/§17) teyit et, sonra ilgili derinlik bölümüne geç.

> Son yasal hatırlatma: Bu dizindeki her araç meşru ve açık-kaynaktır; ama "açık kaynak ve indirilebilir olmak" yasal kullanım garantisi DEĞİLDİR. Alıcı/çözücü araçları geniştir ama içerik çözme/kayıt/paylaşma; her türlü TX (GPS üretme, hücresel test hücresi yayını, deauth); yetkisiz pentest — ülkene, yazılı yetkine ve hedefe göre suç olabilir ve bazıları (GPS spoofing, sahte baz istasyonu, yetkisiz ağ saldırısı) her yerde ağır suçtur. Şüphedeysen yapma; yalnızca kendi cihazların, izole/Faraday test ortamın, yazılı yetkili kapsamın veya açık/kamuya yayınlanan sinyaller üzerinde çalış. Bu kitap hukuki danışmanlık değildir; kendi ülkenin güncel mevzuatını teyit et.
