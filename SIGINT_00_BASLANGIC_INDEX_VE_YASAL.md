# 📡 SİNYAL İSTİHBARATI (SIGINT) EL KİTABI
## Sıfırdan Uzmana — RF Fiziğinden SDR'a, Protokol Çözümlemeden Telekom Güvenliğine

> **Bu bir kitaptır, bir komut listesi değil.** Amacı seni bir *operatör* ya da *scriptçi* yapmak değil — o kodları yazacak, o cihazı kuracak, o sinyali anlayacak **altyapıyı ve mühendislik sezgisini** sana kazandırmaktır. Balık vermez; **balık tutmayı** öğretir. Hem hiç radyo bilmeyen biri buradan başlayıp ilerleyebilir, hem de konuyu daha yoğun öğrenmek isteyen ileri seviye okuyucu derinlik bulur.

> Bu el kitabı, Kanije Kalesi **CTI rehber kütüphanesinin** parçasıdır. Diğer rehberler (VeraCrypt, Tails, MITRE ATT&CK, Wireshark, Malware Analizi…) "veri ve ağ" katmanını ele alır; bu kitap **fiziksel/radyo katmanını** — tüm haberleşmenin altında yatan elektromanyetik gerçekliği — ele alır.

---

## 🚨 YASAL & ETİK MANİFESTO — HER ŞEYDEN ÖNCE OKU

Sinyal istihbaratı **meşru, akademik ve son derece öğretici** bir alandır: amatör radyo, uydu görüntüleme, havacılık/denizcilik takibi, spektrum araştırması, güvenlik testi ve savunma. **Ama aynı bilgi, yanlış elde suç işler.** Bu kitabın çizgisi nettir:

| ✅ Genelde SERBEST (eğitim/hobi) | 🚫 Genelde SUÇ (yetki/lisans gerekir) |
|---|---|
| **Dinlemek/almak** (RX) — açık yayınlar: ADS-B, AIS, NOAA uydu, amatör radyo, FM | **Yetkisiz YAYIN (TX)** — lisanssız vericilik |
| Spektrumu gözlemlemek, waterfall izlemek | **Jamming (karıştırma)** — her yerde ağır suç, can güvenliği tehdidi |
| Kendi cihazını/sinyalini analiz etmek | **Başkasının haberleşmesini dinlemek/çözmek** (telefon, özel telsiz) |
| rtl_433 ile **kendi** sensörlerini okumak | **Spoofing** — sahte GPS, sahte baz istasyonu (IMSI catcher) |
| Lisanslı amatör radyo (çağrı işareti ile) | Şifreli kurumsal/askeri/polis trafiğini çözmek |

**Türkiye'de** ilgili çerçeve: TCK 132–140 (haberleşmenin gizliliği), BTK/ICTA telsiz mevzuatı, amatör telsiz için BTK sınav/lisans. **ABD:** FCC + Wiretap Act. **AB:** ePrivacy. **Kuralın kendisini, kendi ülkenden ve güncel mevzuattan teyit et — bu kitap hukuki danışmanlık değildir.**

> ⚖️ **Altın kural:** *"Alıcı çoğu yerde serbesttir; verici her yerde sorumluluktur."* Şüphedeysen **yapma** veya önce lisans/izin al. Bu kitaptaki tüm alıştırmalar bilinçli olarak **kendi cihazların, açık/yayın sinyaller ve ev ortamı** üzerine kuruludur.

---

## 👥 Bu Kitap Kimin İçin?

- **Yeni başlayan** — "SDR nedir, ne alırım, nasıl başlarım?" diyen; sıfır elektronik bilgisiyle gelip temelden kurabilen.
- **İleri seviye** — modülasyon matematiği, link budget, anten tasarımı, protokol tersine mühendisliği ve telekom güvenliği derinliği arayan.
- **CTI / güvenlik araştırmacısı** — RF saldırı yüzeyini, telekom zafiyetlerini ve savunmayı anlamak isteyen.
- **Hobici / maker** — kendi antenini, LNA'sını, PCB'sini kurmak isteyen.

---

## 📚 KİTABIN YAPISI — 6 Bölüm

| # | Bölüm | Dosya | Ne öğrenirsin |
|---|---|---|---|
| 1 | **Temeller: RF Fiziği & Modülasyon** | [`SIGINT_01_TEMELLER_RF_VE_MODULASYON.md`](SIGINT_01_TEMELLER_RF_VE_MODULASYON.md) | EM spektrum, dB matematiği, AM/FM/PSK/QAM/OFDM, IQ, Nyquist, formüller (f=c/λ, FSPL, Shannon) |
| 2 | **SDR Cihazları Derinlemesine** | [`SIGINT_02_SDR_CIHAZLARI_DERINLEMESINE.md`](SIGINT_02_SDR_CIHAZLARI_DERINLEMESINE.md) | RTL-SDR V3/V4, HackRF, Airspy, SDRplay, LimeSDR, USRP, Pluto, KrakenSDR, Flipper Zero — karşılaştırma + "hangisi" |
| 3 | **Antenler, Donanım & Devre** | [`SIGINT_03_ANTEN_DONANIM_VE_DEVRE_TASARIMI.md`](SIGINT_03_ANTEN_DONANIM_VE_DEVRE_TASARIMI.md) | Anten türleri/boyu, LNA, filtre, watt/güç mantığı, kendi PCB'ni çizme, "nereden biliyorlar" |
| 4 | **Yazılım, OS & Kurulum** | [`SIGINT_04_YAZILIM_OS_VE_KURULUM.md`](SIGINT_04_YAZILIM_OS_VE_KURULUM.md) | Linux/DragonOS/Windows, sürücüler, GQRX/SDR#/GNU Radio/URH, komut araçları, sorun→çözüm |
| 5 | **Protokoller & Sinyal Çözümleme** | [`SIGINT_05_PROTOKOLLER_VE_SINYAL_COZUMLEME.md`](SIGINT_05_PROTOKOLLER_VE_SINYAL_COZUMLEME.md) | ADS-B, AIS, NOAA uydu, rtl_433, APRS, FT8 (yasal); 2G/3G/4G/5G mimari & prensip |
| 6 | **Güvenlik, Açıklar & Savunma** | [`SIGINT_06_GUVENLIK_ACIKLAR_VE_SAVUNMA.md`](SIGINT_06_GUVENLIK_ACIKLAR_VE_SAVUNMA.md) | Manipüle edilebilir/edilemez sinyaller, replay/spoofing/jamming prensibi, SS7/IMSI-catcher, savunma, açık takip |

---

## 🎓 ÖĞRENME YOLU — Zero to Hero

Bölümler numara sırasıyla da okunur ama **en hızlı öğrenme yolu** şudur:

```
   BAŞLA
     │
     ▼
[1] Temeller ──────► "Sinyal nedir, dB ne demek, modülasyon nasıl bilgi taşır?"
     │                 (Fizik olmadan gerisi ezber olur — atlama.)
     ▼
[2] Cihaz Seç ─────► "Bütçeme ve amacıma göre hangi SDR?" → muhtemelen RTL-SDR Blog V4 ile başla
     │
     ▼
[4] Yazılım Kur ───► "GQRX/SDR# kur, ilk kez bir sinyal gör." (Hızlı başarı = motivasyon)
     │
     ▼
[5] İlk Sinyaller ─► ADS-B (uçaklar), NOAA (uydudan hava görüntüsü), rtl_433 (kendi sensörlerin)
     │                 ← Burada "vay be" anını yaşarsın. Hepsi YASAL.
     ▼
[3] Donanımı Geliştir ► "Daha iyi anten/LNA yap, neden işe yaradığını anla, kendi PCB'ne uzan."
     │
     ▼
[6] Güvenlik ──────► "Sinyaller nasıl saldırıya uğrar, nasıl savunulur?" (En son — temel oturunca)
     │
     ▼
   USTA  →  Artık yeni bir sinyal gördüğünde "bu ne, nasıl çözerim?" diyebilirsin.
```

> 💡 **Pedagoji:** Her bölümün sonunda **ev/güvenli ortam alıştırmaları** var. Bilgi okuyarak değil, **pratik yaparak** oturur. RTL-SDR Blog V4 + basit bir anten (~40$) ile bu kitaptaki yasal alıştırmaların neredeyse tamamını yapabilirsin.

---

## 🛒 HIZLI CİHAZ SEÇİMİ (detay → Bölüm 2)

| Amacın | Önerilen başlangıç | Neden |
|---|---|---|
| **İlk SDR / dinleme / en ucuz** | **RTL-SDR Blog V4** (~30-40$) | Geniş kapsama, mükemmel topluluk, tüm yasal alıştırmalar |
| **Verici/deney (TX)** | **HackRF One** | 1 MHz–6 GHz, TX yeteneği *(yasal sınıra dikkat)* |
| **Temiz HF / kısa dalga** | **SDRplay RSPdx** / **Airspy HF+** | 14-bit, düşük gürültü |
| **Yön bulma / pasif radar** | **KrakenSDR** | 5 kanal faz-tutarlı |
| **Taşınabilir RF çoklu-araç** | **Flipper Zero** | Sub-GHz/NFC/RFID/IR *(SDR değil, sınırlı)* |
| **Araştırma / full-duplex** | **USRP B2xx / LimeSDR** | Geniş bant, MIMO |

---

## 🧪 GÜVENLİ PRATİK ALANLARI (hepsi yasal, evden)

1. ✈️ **ADS-B** — üzerinden geçen uçakları haritada izle (1090 MHz, dump1090)
2. 🛰️ **NOAA / Meteor uydu** — gökyüzünden geçen uydudan **kendi** hava durumu görüntünü al (137 MHz)
3. 🌡️ **rtl_433** — evindeki kablosuz termometre/priz/istasyonu oku (kendi cihazların)
4. 📻 **FT8 / amatör radyo** — dünya çapında zayıf sinyalleri (sadece dinleyerek) gözle
5. 🎛️ **Kendi uzaktan kumandan** — garaj/araç kumandanın (KENDİ malın) sinyalini URH ile analiz et, sabit-kod mu rolling-code mu öğren

---

## 🔗 İlgili Kanije CTI Kütüphanesi

Bu el kitabı, repo'daki güvenlik rehberi koleksiyonunun RF/fiziksel-katman ayağıdır. Tamamlayıcılar:
`MITRE_ATTACK_USTALIK_REHBERI.md` · `TTP_AVCILIGI_USTALIK_REHBERI.md` · `MISP_THREAT_INTEL_USTALIK_REHBERI.md` · `WIRESHARK_AG_ANALIZ_USTALIK_REHBERI.md` · `OSINT_ARAC_SETI_USTALIK_REHBERI.md` · `MALWARE_ANALIZ_USTALIK_REHBERI.md` · `VERACRYPT_USTALIK_REHBERI.md`

---

> 🛰️ **Kapanış:** Elektromanyetik spektrum görünmezdir ama her yerdedir — telefonun, WiFi'ın, uçaklar, uydular, garaj kapın. Bu kitap, o görünmez dünyayı **görünür** kılar. Onu **anlamak** güçtür; onu **sorumlu kullanmak** olgunluktur. İkisini birlikte taşı.
>
> *Sıradaki adım: [Bölüm 1 — RF Fiziği & Modülasyon](SIGINT_01_TEMELLER_RF_VE_MODULASYON.md) ile başla.*
