# Sinyal İstihbaratı (SIGINT) El Kitabı
## Sıfırdan Uzmana — RF Fiziğinden Yapay Zekaya, Tam Kapsamlı Başvuru Kaynağı

Bu bir başvuru kitabıdır; bir komut listesi değildir. Amacı okuyucuyu bir operatör ya da scriptçi yapmak değil, o cihazı kuracak, o sinyali çözecek ve o tehdidi savunacak mühendislik altyapısını ve sezgisini kazandırmaktır. Hiç radyo bilmeyen biri buradan başlayıp ilerleyebilir; konuyu derinlemesine bilen biri ise her bölümde ileri seviye derinlik bulur. Yirmi üç bölüm, RF fiziğinin temelinden modern yapay zeka tabanlı sinyal istihbaratına ve geleceğin dalga formlarına kadar uzanır.

Bu el kitabı, Kanije Kalesi güvenlik dokümantasyon kütüphanesinin radyo/fiziksel katman ayağıdır. Diğer rehberler (VeraCrypt, Tails, MITRE ATT&CK, Wireshark, OSINT, Malware Analizi vb.) veri ve ağ katmanını ele alır; bu kitap tüm haberleşmenin altında yatan elektromanyetik gerçekliği inceler.

---

## Yasal ve Etik Çerçeve — Her Şeyden Önce

Sinyal istihbaratı meşru, akademik ve son derece öğretici bir mühendislik alanıdır: amatör radyo, uydu görüntüleme, havacılık ve denizcilik takibi, spektrum araştırması, güvenlik testi ve savunma bunun parçasıdır. Aynı bilgi yanlış kullanıldığında ise suç işler. Bu kitabın çizgisi nettir ve istisnasız korunur:

| Genelde serbest (eğitim ve hobi) | Genelde suç (yetki veya lisans gerekir) |
|---|---|
| Dinlemek ve almak (RX): ADS-B, AIS, NOAA uydu, amatör radyo, FM | Yetkisiz yayın (TX): lisanssız vericilik |
| Spektrumu gözlemlemek, waterfall izlemek | Jamming (karıştırma): her yerde ağır suç, can güvenliği tehdidi |
| Kendi cihazını ve sinyalini analiz etmek | Başkasının haberleşmesini dinlemek veya çözmek |
| Kendi ağında yetkili penetrasyon testi | Başkasının ağına, kartına, cihazına izinsiz erişim |

Türkiye'de ilgili çerçeve: TCK 132-140 (haberleşmenin gizliliği) ve TCK 243-244 (bilişim sistemlerine yetkisiz erişim), BTK telsiz mevzuatı, amatör telsiz için BTK sınav ve lisansı. Amerika Birleşik Devletleri: FCC ve Wiretap Act. Avrupa Birliği: ePrivacy. Kuralın kendisini kendi ülkenden ve güncel mevzuattan teyit et; bu kitap hukuki danışmanlık değildir.

Yol gösterici ilke: alıcı çoğu yerde serbesttir, verici her yerde sorumluluktur. Şüphedeysen yapma veya önce lisans ve izin al. Bu kitaptaki her alıştırma bilinçli olarak kendi cihazların, açık yayınlar ve ev ortamı üzerine kuruludur. Jamming, yetkisiz dinleme, spoofing ve sistemlere izinsiz müdahale bu kitabın kapsamı dışındadır; bu konular yalnızca nasıl çalıştıkları ve nasıl savunulacakları yönüyle, savunma perspektifiyle ele alınır.

---

## Bu Kitap Kimin İçin

- Yeni başlayan: "SDR nedir, ne alırım, nasıl başlarım" diyen, sıfır elektronik bilgisiyle gelip temelden kurabilen okuyucu.
- İleri seviye: modülasyon matematiği, DSP iç mimarisi, anten tasarımı, protokol tersine mühendisliği, telekom güvenliği ve yapay zeka tabanlı sinyal analizi derinliği arayan okuyucu.
- Güvenlik araştırmacısı: RF saldırı yüzeyini, telekom ve kablosuz zafiyetlerini ve savunmayı anlamak isteyen profesyonel.
- Hobici ve maker: kendi antenini, LNA'sını, devresini ve PCB'sini kurmak isteyen meraklı.

---

## Kitabın Yapısı — 23 Bölüm

### Temel Katman

| # | Bölüm | Dosya |
|---|---|---|
| 00 | Başlangıç, İndeks ve Yasal Çerçeve | bu dosya |
| 01 | Temeller: RF Fiziği ve Modülasyon | [SIGINT_01_TEMELLER_RF_VE_MODULASYON.md](SIGINT_01_TEMELLER_RF_VE_MODULASYON.md) |
| 02 | SDR Cihazları Derinlemesine | [SIGINT_02_SDR_CIHAZLARI_DERINLEMESINE.md](SIGINT_02_SDR_CIHAZLARI_DERINLEMESINE.md) |
| 03 | Antenler, Donanım ve Devre Tasarımı | [SIGINT_03_ANTEN_DONANIM_VE_DEVRE_TASARIMI.md](SIGINT_03_ANTEN_DONANIM_VE_DEVRE_TASARIMI.md) |
| 04 | Yazılım, İşletim Sistemi ve Kurulum | [SIGINT_04_YAZILIM_OS_VE_KURULUM.md](SIGINT_04_YAZILIM_OS_VE_KURULUM.md) |
| 05 | Protokoller ve Sinyal Çözümleme | [SIGINT_05_PROTOKOLLER_VE_SINYAL_COZUMLEME.md](SIGINT_05_PROTOKOLLER_VE_SINYAL_COZUMLEME.md) |
| 06 | Güvenlik, Açıklar ve Savunma | [SIGINT_06_GUVENLIK_ACIKLAR_VE_SAVUNMA.md](SIGINT_06_GUVENLIK_ACIKLAR_VE_SAVUNMA.md) |

### İstihbarat ve Konum Katmanı

| # | Bölüm | Dosya |
|---|---|---|
| 07 | SIGINT Disiplinleri ve Sinyal Ayıklama | [SIGINT_07_DISIPLINLER_VE_SINYAL_AYIKLAMA.md](SIGINT_07_DISIPLINLER_VE_SINYAL_AYIKLAMA.md) |
| 08 | Frekans Tahsisi ve Bant Planı (askeri/sivil/havacılık/denizci) | [SIGINT_08_FREKANS_TAHSISI_VE_BANT_PLANI.md](SIGINT_08_FREKANS_TAHSISI_VE_BANT_PLANI.md) |
| 09 | Yer Tespiti, Yön Bulma ve Takip | [SIGINT_09_YER_TESPITI_YON_BULMA_VE_TAKIP.md](SIGINT_09_YER_TESPITI_YON_BULMA_VE_TAKIP.md) |
| 10 | GNSS/GPS Sistemleri | [SIGINT_10_GNSS_GPS_SISTEMLERI.md](SIGINT_10_GNSS_GPS_SISTEMLERI.md) |
| 11 | Uydu Haberleşmesi | [SIGINT_11_UYDU_HABERLESMESI.md](SIGINT_11_UYDU_HABERLESMESI.md) |
| 12 | DragonOS ve Araç Ekosistemi | [SIGINT_12_DRAGONOS_VE_ARAC_EKOSISTEMI.md](SIGINT_12_DRAGONOS_VE_ARAC_EKOSISTEMI.md) |
| 13 | RF Tehdit Manzarası ve Karşı-Önlemler | [SIGINT_13_RF_TEHDIT_VE_KARSI_ONLEMLER.md](SIGINT_13_RF_TEHDIT_VE_KARSI_ONLEMLER.md) |
| 14 | İstihbarat Kaynakları ve Güncel Takip | [SIGINT_14_ISTIHBARAT_KAYNAKLARI_VE_TAKIP.md](SIGINT_14_ISTIHBARAT_KAYNAKLARI_VE_TAKIP.md) |

### İleri ve Güncel Katman

| # | Bölüm | Dosya |
|---|---|---|
| 15 | WiFi/WLAN Güvenliği (el-sıkışması yakalama dahil) | [SIGINT_15_WIFI_WLAN_GUVENLIGI.md](SIGINT_15_WIFI_WLAN_GUVENLIGI.md) |
| 16 | Kısa Menzilli Kablosuz ve IoT (BLE, RFID/NFC, Zigbee, LoRa) | [SIGINT_16_KISA_MENZIL_KABLOSUZ_VE_IOT.md](SIGINT_16_KISA_MENZIL_KABLOSUZ_VE_IOT.md) |
| 17 | TEMPEST, Emanasyon ve Yan-Kanal | [SIGINT_17_TEMPEST_EMANASYON_VE_YAN_KANAL.md](SIGINT_17_TEMPEST_EMANASYON_VE_YAN_KANAL.md) |
| 18 | Sayısal Sinyal İşleme ve SDR İç Mimarisi | [SIGINT_18_DSP_VE_SDR_IC_MIMARI.md](SIGINT_18_DSP_VE_SDR_IC_MIMARI.md) |
| 19 | Yapay Zeka ve Makine Öğrenmesi ile SIGINT | [SIGINT_19_YAPAY_ZEKA_VE_ML_SIGINT.md](SIGINT_19_YAPAY_ZEKA_VE_ML_SIGINT.md) |
| 20 | İleri Hücresel: 4G/5G Güvenlik | [SIGINT_20_ILERI_HUCRESEL_4G_5G_GUVENLIK.md](SIGINT_20_ILERI_HUCRESEL_4G_5G_GUVENLIK.md) |
| 21 | SIGINT Tarihi, Aktörler ve Elektronik Harp | [SIGINT_21_TARIH_AKTORLER_VE_ELEKTRONIK_HARP.md](SIGINT_21_TARIH_AKTORLER_VE_ELEKTRONIK_HARP.md) |
| 22 | Egzotik Yayılım ve Geleceğin SIGINT'i | [SIGINT_22_EGZOTIK_VE_GELECEK.md](SIGINT_22_EGZOTIK_VE_GELECEK.md) |

---

## Önerilen Okuma Yolu

Bölümler numara sırasıyla okunabilir; ancak en verimli öğrenme yolu pratiğe erken geçmekten geçer.

```
Başla
  |
  v
[01] Temeller --------> Sinyal nedir, dB ne demek, modülasyon nasıl bilgi taşır.
  |                     Fizik olmadan gerisi ezber olur; bu bölümü atlama.
  v
[02] Cihaz Seç -------> Bütçene ve amacına göre hangi SDR. Genellikle RTL-SDR Blog V4.
  v
[04] Yazılım Kur -----> GQRX/SDR# kur, ilk kez bir sinyal gör. Hizli basari, motivasyon.
  v
[05] Ilk Sinyaller ---> ADS-B (ucaklar), NOAA (uydudan goruntu), rtl_433 (kendi sensorlerin).
  |                     Hepsi yasal; "calisiyor" anini burada yasarsin.
  v
[03] Donanim ---------> Daha iyi anten/LNA, neden ise yaradigini anla, kendi PCB'ne uzan.
  v
[07-14] Istihbarat ---> Disiplinler, frekans planlari, yer tespiti, GPS, uydu, DragonOS, savunma.
  v
[15-20] Ileri --------> WiFi, kisa menzil, TEMPEST, DSP ic mimari, yapay zeka, 5G.
  v
[18] DSP -------------> Motor kaputunu ac: FFT, filtreler, demodulasyon matematigi (en teknik bolum).
  v
[21-22] Baglam -------> Tarih, elektronik harp, egzotik yayilim, gelecek.
  v
Uzman -> Yeni bir sinyal gordugunde "bu ne, nasil cozerim, nasil savunurum" diyebilirsin.
```

Her bölümün sonunda kendi cihazların ve açık yayınlarla yapılabilecek alıştırmalar vardır. Bilgi okuyarak değil pratik yaparak oturur. RTL-SDR Blog V4 ve basit bir anten (yaklaşık 40 dolar) ile bu kitaptaki yasal alıştırmaların neredeyse tamamı yapılabilir.

---

## Hızlı Cihaz Seçimi (detay: Bölüm 2)

| Amacın | Önerilen başlangıç | Neden |
|---|---|---|
| İlk SDR, dinleme, en ucuz | RTL-SDR Blog V4 | Geniş kapsama, güçlü topluluk, tüm yasal alıştırmalar |
| Verici ve deney (TX) | HackRF One | 1 MHz - 6 GHz, TX yeteneği (yasal sınıra dikkat) |
| Temiz HF ve kısa dalga | SDRplay RSPdx, Airspy HF+ | 14 bit, düşük gürültü |
| Yön bulma, pasif radar | KrakenSDR | Beş kanal faz-tutarlı |
| Taşınabilir RF çoklu araç | Flipper Zero | Sub-GHz, NFC, RFID, IR (SDR değil, sınırlı) |
| Araştırma, full-duplex | USRP, LimeSDR | Geniş bant, MIMO |

---

## Güvenli Pratik Alanları (hepsi yasal, evden)

1. ADS-B: üzerinden geçen uçakları haritada izle (1090 MHz, dump1090).
2. NOAA ve Meteor uydusu: gökyüzünden geçen uydudan kendi hava durumu görüntünü al (137 MHz).
3. rtl_433: evindeki kablosuz termometre, priz ve istasyonu oku (kendi cihazların).
4. WSPR ve FT8: dünya çapında zayıf sinyalleri yalnızca dinleyerek gözle.
5. Kendi ağın ve cihazların: WiFi handshake, BLE reklamı, RFID kartı, uzaktan kumanda — yalnızca senin malın üzerinde, savunmayı öğrenmek için.

---

## İlgili Kanije Kütüphanesi

Bu el kitabı, repodaki güvenlik dokümantasyon koleksiyonunun RF ve fiziksel katman ayağıdır. Tamamlayıcılar: MITRE_ATTACK_USTALIK_REHBERI, TTP_AVCILIGI_USTALIK_REHBERI, MISP_THREAT_INTEL_USTALIK_REHBERI, WIRESHARK_AG_ANALIZ_USTALIK_REHBERI, OSINT_ARAC_SETI_USTALIK_REHBERI, MALWARE_ANALIZ_USTALIK_REHBERI, VERACRYPT_USTALIK_REHBERI.

---

Kapanış. Elektromanyetik spektrum görünmezdir ama her yerdedir: telefonun, kablosuz ağın, uçaklar, uydular, garaj kapın. Bu kitap o görünmez dünyayı görünür kılar. Onu anlamak güçtür; onu sorumlu kullanmak olgunluktur. İkisini birlikte taşı.

Sıradaki adım: [Bölüm 1 — RF Fiziği ve Modülasyon](SIGINT_01_TEMELLER_RF_VE_MODULASYON.md).
