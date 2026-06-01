# SIGINT EL KİTABI — BÖLÜM 24: GÜNCEL ZAFİYET MANZARASI

## Kablosuz ve Telekom Zafiyet Manzarası — Bilinen Açıklar, Prensipleri ve Savunması

> Amaç: Bu bölüm bir zafiyet farkındalık kataloğudur. Önceki bölümler sinyalin fiziğini, protokolünü ve savunmasını ayrı ayrı işledi; burada amaç, kablosuz ve telekom dünyasının bilinen ve hâlâ güncel açıklarını tek bir referans altında toplamak, her birinin neden mümkün olduğunu (hangi mekanizmanın zayıf olduğunu), etkisini, yamandı/sürüyor durumunu, savunmasını ve hangi açık-kaynak araçla yetkili biçimde test edileceğini netleştirmektir. Bu, bir CTI (Cyber Threat Intelligence) ya da savunma analisti için yazılmıştır: hedefi tehdit modellemesi, yama önceliklendirmesi ve risk değerlendirmesidir.

> Yasal ve etik çerçeve: Bu katalog "nasıl tetiklenir" reçetesi vermez. Her zafiyet, prensip + etki + savunma + yetkili test aracı üçgeninde ele alınır. Anlatılan açıklar kamuya açık akademik makalelerde, üretici advisory'lerinde ve CVE veritabanlarında yayımlanmış, sorumlu açıklama (CVD) sürecinden geçmiş kayıtlardır. Test araçları yalnızca kendi cihazlarında, kendi ağında ya da açık yazılı izinle (yetkili sızma testi) kullanılmak üzere anılır; başkasının trafiğini çözmek, klonlamak veya bozmak çoğu ülkede suçtur. Verilen her CVE numarası ve tarih elden geldiğince doğrulanmıştır; emin olunamayan kayıtlar "NVD/CVE'den teyit edilmeli" notuyla işaretlenmiştir — bu bölüm asla uydurma CVE vermez, şüpheyi açıkça yazar.

> Bu bölümün önceki bölümlerle ilişkisi: WiFi açıkları Bölüm 15'in, kısa menzilli/IoT açıkları Bölüm 16'nın, hücresel ileri konular Bölüm 20'nin, SS7/Diameter ve genel telekom çekirdeği Bölüm 6 ve 20'nin, GNSS olayları Bölüm 10'un, RF tehdit/karşı-önlem Bölüm 13'ün, güncel kaynak takibi Bölüm 14'ün derinlikli işlediği konulardır. Bu bölüm onların üstünde bir indeks ve zaman çizelgesi katmanı kurar; tekrar değil, birleştirme amaçlıdır.

---

## İÇİNDEKİLER

1. [Zafiyet Manzarasını Okuma: Katalog Nasıl Kullanılır](#1)
2. [WiFi / 802.11 Zafiyetleri (WEP'ten FragAttacks'e)](#2)
3. [Bluetooth ve BLE Zafiyetleri (BlueBorne'dan BLESA'ya)](#3)
4. [Hücresel Zafiyetler (A5'ten 5G Downgrade'e)](#4)
5. [Telekom Çekirdeği: SS7 ve Diameter Sinyalleşme Zafiyetleri](#5)
6. [GNSS Zafiyetleri ve Gerçek-Dünya Kesintileri](#6)
7. [IoT / Sub-GHz / Otomotiv: Rolling-Code, KeeLoq, Crypto-1, Zigbee/Z-Wave, LoRaWAN](#7)
8. [RFID / NFC: Kart Klonlama ve Relay Zafiyetleri](#8)
9. [Donanım, Firmware ve SDR Yan-Kanal Zafiyet Kavramı](#9)
10. [Yıllara Göre Zaman Çizelgesi (Büyük Tablo)](#10)
11. [Açığın Yaşam Döngüsü: Keşiften KEV'e (ASCII)](#11)
12. [Nasıl Güncel Kalınır: Kaynak Disiplini](#12)
13. [Yetkili Test Araçları (Zafiyet Sınıfına Göre)](#13)
14. [Alıştırmalar (Yasal, Kendi Cihazların)](#14)
15. [Hızlı Referans ve Diğer Bölümler](#15)

---

<a id="1"></a>
## 1. Zafiyet Manzarasını Okuma: Katalog Nasıl Kullanılır

Bir zafiyet kataloğu, ezberlenecek bir liste değil, düşünülecek bir haritadır. Her kayıt için kendine dört soru sorman, kataloğu istihbarata çevirir. Birincisi mekanizma sorusu: bu açık hangi mekanizmanın zayıflığından doğuyor — bir kriptografik ilkellin kırılmasından mı (zayıf şifre, kısa anahtar), bir protokol durum makinesinin hatasından mı (yeniden gönderim, parçalama mantığı), bir uygulama (implementation) hatasından mı (bellek taşması, sınır denetimi yokluğu), yoksa bir tasarım kararından mı (downgrade'e izin veren geriye uyumluluk)? Bu ayrım, savunmanın da nereye konacağını belirler.

İkincisi katman sorusu: açık fiziksel katmanda mı (RF sinyalin kendisi), bağlantı/MAC katmanında mı (çerçeve, el-sıkışma), yoksa ağ/uygulama katmanında mı (sinyalleşme protokolü, çekirdek ağ) oturuyor? Üçüncüsü erişim sorusu: saldırganın menzilde olması yeterli mi (yakınlık), aynı ağda olması mı gerekiyor, yoksa operatör ağına erişimi mi şart (örneğin SS7)? Dördüncüsü durum sorusu: yamandı mı, kısmen mi azaltıldı, yoksa protokolün doğasında olduğu için kalıcı mı?

Bu dört eksen, kataloğun her satırını okurken zihninde işletmen gereken çerçevedir. Aşağıdaki tablo, zafiyet sınıflarını bu eksenlerde özetler ve kataloğun nasıl tarandığını gösterir.

| Mekanizma sınıfı | Tipik örnek | Katman | Savunma yönü |
|---|---|---|---|
| Zayıf/kırık kripto ilkel | WEP RC4, A5/1, Crypto-1 | Bağlantı/uygulama | İlkeli değiştir (AES, KASUMI üstü, AES kart) |
| Anahtar yönetimi/yeniden kullanım | KRACK (nonce reuse), KNOB (zayıf anahtar pazarlığı) | Bağlantı | Yama + güçlü anahtar zorlama |
| Protokol durum makinesi hatası | FragAttacks, BIAS, BLESA | Bağlantı/MAC | Yama + sıkı durum doğrulama |
| Uygulama (kod) hatası | BlueBorne, SweynTooth, BleedingTooth | Yığın/sürücü | Yama + bellek-güvenli kod |
| Yan kanal | Dragonblood (zamanlama/önbellek), Kr00k (sıfır anahtar) | Fiziksel/uygulama | Sabit-zaman kod + yama |
| Tasarım/geriye uyumluluk | 2G downgrade, IMSI ifşası, SS7 güveni | Sistem/ağ | Eski sürümü kapat, güven sınırı çiz |

> Sezgi: İyi yamanmış bir açık bile, yama dağıtılmadıkça (cihaz güncellenmedikçe) sömürülmeye devam eder. Bu yüzden katalogdaki "durum" sütunu iki şey söyler: protokol/üretici tarafında durum (yama yayımlandı mı) ve saha tarafında durum (cihazlar güncellendi mi). İkisi neredeyse hiçbir zaman aynı anda "tamam" olmaz. Eski açıkların hâlâ sömürülmesinin tek sebebi budur (Bölüm 11).

---

<a id="2"></a>
## 2. WiFi / 802.11 Zafiyetleri (WEP'ten FragAttacks'e)

WiFi, en çok incelenmiş kablosuz protokol ailesidir; dolayısıyla zafiyet geçmişi en zengin olanıdır. Aşağıdaki katalog tarihsel (artık yalnızca eski cihazlarda görülen) ve güncel açıkları birlikte verir. Derin protokol anlatımı Bölüm 15'tedir; burada her açığın katalog kaydı tutulur.

| Zafiyet / ad | CVE / referans | Sistem / protokol | Keşif | Prensip (neden savunmasız) | Etki | Durum | Savunma |
|---|---|---|---|---|---|---|---|
| WEP zayıflığı | (CVE öncesi; FMS 2001, sonraki PTW saldırıları) | 802.11 WEP, RC4 + zayıf IV | 2001 | RC4 anahtar akışı, 24-bit IV'lerin tekrarı ve zayıf-IV sızıntısıyla istatistiksel olarak çözülür; anahtar trafikten geri çıkarılabilir | Anahtar tamamen kurtarılır; ağ açılır | Tarihsel; WEP terk edildi | WEP'i tamamen kapat; WPA2/WPA3 kullan |
| WPA-TKIP zayıflığı | (Beck-Tews 2008 ve türevleri) | WPA (TKIP/MIC) | 2008-2009 | TKIP, WEP'in üstüne yamadır; MIC (Michael) ve QoS kanal yapısı sınırlı paket enjeksiyonu/çözümüne kapı aralar | Kısmi paket manipülasyonu; tam anahtar değil | Tarihsel; TKIP devre dışı bırakılmalı | TKIP'i kapat, yalnızca CCMP/AES; WPA3'e geç |
| KRACK | CVE-2017-13077 ... CVE-2017-13088 (10 CVE serisi) | WPA/WPA2 4-yönlü el-sıkışma (ve grup/FT) | 2017 | Anahtar Yeniden Kurulum (Key Reinstallation): el-sıkışmanın 3. mesajı yeniden gönderilince istemci nonce/replay sayacını sıfırlar; aynı anahtar+nonce ile şifreleme yeniden kullanılır | Paket çözme/yeniden oynatma; bazı yığınlarda tüm-sıfır anahtar | Yamandı (2017 istemci/AP yamaları) | İşletim sistemi/sürücü yamasını uygula; PMF; WPA3 |
| Kr00k | CVE-2019-15126 | Belirli WiFi yongaları (Broadcom/Cypress) | 2019 | İlişki kesilince (disassociation) yonga, tampondaki veriyi tamamı-sıfır oturum anahtarıyla şifreleyip gönderir; sıfır anahtar trivially çözülür | Birkaç çerçevelik şifreli verinin çözülmesi | Yamandı (firmware/sürücü) | Yonga firmware/sürücü güncellemesi; üst-katman şifreleme (TLS) |
| Dragonblood | CVE-2019-9494 (zamanlama yan kanal), CVE-2019-9495 (önbellek yan kanal), CVE-2019-9496/9497/9498/9499 (ilgili EAP-pwd dâhil) | WPA3-SAE (Dragonfly el-sıkışma) | 2019 | SAE'nin hash-to-curve/grup işlemleri sabit-zaman değilse, zamanlama ve önbellek yan kanalları parola hakkında bilgi sızdırır; ayrıca downgrade ve grup-indirme vektörleri | Parola bölümleme (partitioning) saldırısıyla parola tahmininin hızlanması; downgrade | Büyük ölçüde yamandı (hostapd/wpa_supplicant + standart iyileştirme, Hash-to-Element) | Güncel hostapd/wpa_supplicant; SAE-H2E; WPA3-only mod |
| FragAttacks | CVE-2020-24586, CVE-2020-24587, CVE-2020-24588 (tasarım) + uygulama CVE'leri CVE-2020-26139 ... CVE-2020-26147 | 802.11 çerçeve parçalama (fragmentation) ve toplama (aggregation) | 2021 (açıklama; bazıları 802.11'in doğasında) | Çerçeve toplama bayrağının doğrulanmaması, farklı anahtarlarla şifrelenmiş parçaların birleştirilebilmesi ve önbellekte kalan parçaların karışması; bazıları tasarım kusuru, bazıları uygulama hatası | Paket enjeksiyonu, sınırlı veri sızıntısı, istemciyi kötü sunucuya yönlendirme | Kısmen yamandı; bazı tasarım kusurları azaltma gerektirir | Üretici yamaları; uçtan-uca şifreleme (HTTPS/VPN); PMF |
| WPS PIN zayıflığı | CVE-2011-5053 (genel WPS) / Reaver sınıfı | WiFi Protected Setup, 8 haneli PIN | 2011 | PIN iki yarıya bölünüp ayrı doğrulandığından arama uzayı 10^8'den ~11.000'e düşer; çevrimdışı PIN (Pixie Dust) bazı yongalarda zayıf rastgelelikle anında çözülür | Router PIN'i ve WPA anahtarı kurtarılır | Sürüyor (eski/ucuz router'larda WPS açık) | WPS'i kapat; WPS yoksa açma; firmware güncelle |
| KARMA / PNL istismarı | (CVE'siz davranış sınıfı) | İstemci tercih edilen ağ listesi (PNL) + probe | 2004+ | İstemci geçmiş ağları açık probe ederse, sahte AP "evet o ağ benim" diye yanıtlayıp istemciyi kendine çeker | Evil Twin'e otomatik bağlanma, MITM | Kısmen azaltıldı (MAC randomizasyon, pasif tarama) | Otomatik bağlanmayı kapat; bilinen ağ doğrulaması; 802.1X |

Notlar ve doğruluk uyarıları: KRACK için on ayrı CVE atanmıştır (CVE-2017-13077'den CVE-2017-13088'e kadar, FT el-sıkışma ve grup anahtarı varyantları dâhil); bunların tam eşlemesi NVD'den teyit edilmelidir. FragAttacks'in üç "tasarım" CVE'si (CVE-2020-24586/24587/24588) ile çok sayıda "uygulama" CVE'si (CVE-2020-26139 ... CVE-2020-26147 aralığı) ayrı kayıtlardır; tam liste fragattacks.com ve NVD'den doğrulanmalıdır. WPS için tek kanonik CVE yerine bir davranış sınıfı söz konusudur; "Pixie Dust" belirli yonga/firmware'lere bağlı ayrı kayıtlardır.

> Savunma özeti (WiFi): Modern minimum çizgi — WPA3 (mümkünse WPA3-only, değilse WPA2/WPA3 geçiş modunda PMF zorunlu), TKIP ve WPS kapalı, sürücü/firmware güncel, üst-katmanda mutlaka uçtan-uca şifreleme (HTTPS/VPN). Bu beş madde, yukarıdaki kataloğun büyük kısmını saha tarafında etkisiz bırakır. Ayrıntı: Bölüm 15 §13.

---

<a id="3"></a>
## 3. Bluetooth ve BLE Zafiyetleri (BlueBorne'dan BLESA'ya)

Bluetooth (Klasik/BR-EDR) ve BLE, geniş bir saldırı yüzeyi sunar çünkü hem karmaşık bir yığın (stack) hem de çoklu eşleşme/şifreleme modu vardır. Aşağıdaki açıklar üç gruba düşer: yığın kod hataları (bellek), kripto/anahtar-pazarlık zayıflıkları ve protokol durum makinesi hataları. Protokol anlatımı Bölüm 16 §2-3'tedir.

| Zafiyet / ad | CVE / referans | Sistem / protokol | Keşif | Prensip (neden savunmasız) | Etki | Durum | Savunma |
|---|---|---|---|---|---|---|---|
| BlueBorne | CVE-2017-1000251 (Linux L2CAP), CVE-2017-1000250 (Linux SDP), CVE-2017-0781/0782/0783 (Android), CVE-2017-8628 (Windows), CVE-2017-14315 (iOS LEAP) | Bluetooth yığını (çoklu OS) | 2017 | Çeşitli yığınlarda eşleşme gerektirmeyen bellek bozulması/sınır hataları; saldırgan menzilde, kullanıcı etkileşimi olmadan kod çalıştırabilir veya MITM kurabilir | Uzaktan kod çalıştırma, bilgi sızıntısı, MITM | Yamandı (2017 OS yamaları) | OS/yığın yaması; Bluetooth'u gereksizse kapat |
| KNOB | CVE-2019-9506 | Bluetooth BR/EDR şifreleme anahtarı pazarlığı | 2019 | Eşleşmede şifreleme anahtarı entropisi (uzunluğu) pazarlığı korumasızdır; saldırgan anahtar boyunu 1 bayta indirgeyip kaba kuvvetle çözebilir | Şifreli bağlantının kırılması, dinleme/değiştirme | Yamandı (spec + yamalar: min. anahtar uzunluğu zorlama) | Minimum anahtar uzunluğu (7 bayt) zorlayan firmware; güncelleme |
| BIAS | CVE-2020-10135 | Bluetooth BR/EDR güvenli kimlik doğrulama | 2020 | Daha önce eşleşmiş bir cihazın kimliğine bürünme; rol değişimi ve eski kimlik doğrulama prosedürünün suistimaliyle karşılıklı doğrulama atlatılır | Cihaz taklidi, oturum ele geçirme, MITM | Yamandı (spec açıklama + yamalar) | Güvenli Bağlantılar (Secure Connections) zorla; firmware güncelle |
| SweynTooth | CVE-2019-16336, CVE-2019-17519, CVE-2019-17517, CVE-2019-19192, CVE-2019-18614 vb. (çoklu satıcı SoC) | BLE SoC/yığın uygulamaları (birden çok üretici) | 2020 | Belirli BLE denetleyici SoC firmware'lerinde paket-ayrıştırma hataları (taşma, beklenmeyen sıralama); kötü biçimli paketlerle çökme/kilitlenme, bazılarında atlama | Hizmet reddi (çökme/donma), bazı durumda güvenlik atlatma | Yamandı (etkilenen SoC üreticilerince) | Etkilenen SoC firmware güncellemesi; tıbbi/kritik cihazda öncelik |
| BleedingTooth | CVE-2020-12351, CVE-2020-12352, CVE-2020-24490 | Linux çekirdeği BlueZ (L2CAP/A2MP/HCI) | 2020 | Linux çekirdeğinin Bluetooth alt sisteminde bellek hataları; menzildeki saldırgan kötü biçimli paketlerle çekirdek bağlamında kod çalıştırabilir/sızdırabilir | Çekirdek düzeyinde RCE/bilgi sızıntısı (Linux) | Yamandı (çekirdek yamaları) | Linux çekirdeğini güncelle; BlueZ yaması |
| BLESA | (CVE'siz/satıcıya özel kayıtlar; akademik 2020) | BLE yeniden bağlanma (reconnection) kimlik doğrulama | 2020 | BLE'de yeniden bağlanmada kimlik doğrulama isteğe bağlı/zayıf uygulanırsa, sahte sunucu önceden eşleşmiş istemciye sahte veri besleyebilir; "reconnection spoofing" | Sahte veri enjeksiyonu, istemci aldatma | Kısmen; uygulama-bağımlı yamalar | Yeniden bağlanmada kimlik doğrulamayı zorunlu kıl; güncel yığın |

Notlar ve doğruluk uyarıları: BlueBorne sekiz ayrı zafiyetin şemsiye adıdır; her platform için ayrı CVE atanmıştır (yukarıdaki liste başlıcalarıdır, tamamı Armis açıklaması ve NVD'den teyit edilmeli). SweynTooth bir düzineyi aşkın CVE içeren bir paket adıdır; etkilenen SoC ve tam CVE listesi orijinal SUTD/açıklama belgesinden ve NVD'den doğrulanmalıdır — yukarıdaki CVE'ler temsilî örneklerdir. BLESA için tek bir kanonik CVE yaygın değildir; akademik bir bulgu olarak ve satıcıya özel kayıtlarla ele alınmalıdır (NVD'den teyit).

> Savunma özeti (Bluetooth/BLE): Yığın/firmware güncel tut (BlueBorne, SweynTooth, BleedingTooth doğrudan kod hatasıdır — yama tek gerçek çözümdür); eşleşmede Güvenli Bağlantılar + yeterli anahtar uzunluğunu zorla (KNOB, BIAS); BLE'de yeniden bağlanma kimlik doğrulamasını zorunlu kıl (BLESA); ihtiyaç dışı Bluetooth'u kapat ve keşfedilebilirliği sınırla. Ayrıntı: Bölüm 16 §13.

---

<a id="4"></a>
## 4. Hücresel Zafiyetler (A5'ten 5G Downgrade'e)

Hücresel açıklar üç kuşağa yayılır: 2G'nin kriptografik ve kimlik-doğrulama tasarım zayıflıkları, 4G/LTE'nin akademik olarak gösterilmiş katman-2 ve kimlik açıkları, 5G'nin büyük ölçüde kapatılmış ama tamamen bitmemiş artıkları. Derin mimari Bölüm 20'de, telekom çekirdeği §5'te (bu bölümde) ve Bölüm 6 §5'tedir.

| Zafiyet / ad | CVE / referans | Sistem / protokol | Keşif | Prensip (neden savunmasız) | Etki | Durum | Savunma |
|---|---|---|---|---|---|---|---|
| A5/1, A5/2 zayıflığı | (CVE öncesi; akademik 1990'lar-2000'ler) | GSM (2G) ses/veri şifreleme | 1994-2003 | A5/2 kasten zayıf; A5/1 kısa anahtar ve LFSR yapısı nedeniyle önceden hesaplanmış (rainbow) tablolarla pratik sürede çözülür | 2G çağrı/SMS içeriğinin çözülmesi | Tarihsel; A5/3 (KASUMI) ve 2G kapatma yaygınlaşıyor | 2G'yi kapat (mümkünse); A5/3; VoLTE/VoWiFi tercih |
| Tek-yönlü kimlik doğrulama (2G) | (tasarım; CVE'siz) | GSM ağ kimlik doğrulaması | tasarımsal | 2G'de yalnızca ağ aboneyi doğrular, abone ağı doğrulamaz; sahte baz istasyonu (IMSI catcher) bu boşluktan beslenir | Sahte hücreye bağlanma, izleme, downgrade'e zemin | Tasarımsal; 3G/4G/5G karşılıklı doğrulama ekledi | 2G'yi devre dışı bırak; modern ağ; downgrade reddi |
| IMSI catcher sınıfı (sahte baz istasyonu) | (cihaz/teknik sınıfı; CVE'siz) | 2G/3G/4G kimlik ifşası ve downgrade | uzun süredir bilinir | Cihaz, ağ kimliğini doğrulayamadığı (özellikle 2G) ya da ilk kayıtta IMSI açık iletildiği için sahte hücreye çekilir; downgrade ile zayıf kuşağa zorlanır | Konum izleme, IMSI toplama, downgrade ile dinleme | Sürüyor; 5G SUCI ifşayı azaltır ama eski kuşaklar açık | 5G SA + SUCI; 2G kapatma; IMSI-catcher tespit uygulamaları |
| aLTeR | (akademik 2019; CVE'siz) | LTE kullanıcı düzlemi (PDCP şifreli ama bütünlük yok) | 2019 | LTE kullanıcı düzlemi verisi şifrelidir ama bütünlük korumasızdır (AES-CTR esnek); aktif saldırgan şifreli trafiği değiştirip DNS yönlendirmesi yapabilir | Trafiği kötü sunucuya yönlendirme (DNS spoofing benzeri) | LTE'de doğası gereği; 5G kullanıcı-düzlemi bütünlüğü çözüm | 5G kullanıcı-düzlemi bütünlük koruması; DNSSEC/DoH; VPN |
| ToRPEDO | (akademik 2019; CVE'siz) | 4G/5G çağrı (paging) protokolü | 2019 | Paging mesajlarının sabit zamanlamasından ve TMSI değişim örüntüsünden yararlanarak hedefin orada olup olmadığı ve paging kimliği çıkarılır; IMSI-kırma ataklarına kapı | Varlık tespiti, paging takibi, IMSI'ye yaklaşma | Araştırma; standart/uygulama iyileştirmeleriyle azaltma | TMSI yenileme sıklığı; paging gizliliği iyileştirmeleri |
| LTE downgrade / sahte hücre | (sınıf; CVE'siz) | LTE → 2G/3G zorlama | bilinir | Saldırgan, RRC/yeniden seçim mekanizmalarını manipüle ederek ya da 2G sahte hücre sunarak cihazı zayıf kuşağa düşürür | Zayıf kripto kuşağında dinleme/izleme | Sürüyor; "2G kapat" en etkili azaltma | 2G/3G devre dışı; sadece-5G/4G modu; ağ tarafı downgrade tespiti |
| 5G kayıt/yan-kanal araştırmaları | (akademik; CVE'siz, sürekli güncellenen alan) | 5G NAS/RRC kayıt, SUCI uygulaması | 2019-sürüyor | SUCI doğru uygulanmazsa (ör. null-şema) IMSI yine sızar; bazı NAS prosedürlerinde downgrade/DoS vektörleri araştırılır | Kalan kimlik ifşası, sınırlı DoS | Araştırma; uygulama kalitesine bağlı | SUCI'de gerçek şifreleme şeması; null-şema yasağı; standart güncel |

Notlar ve doğruluk uyarıları: A5/1-A5/2, aLTeR, ToRPEDO, IMSI-kırma ve birçok LTE/5G katman-2 bulgusu akademik makale ve konferans (özellikle USENIX Security, NDSS, IEEE S&P) çıktılarıdır ve genellikle klasik anlamda bir üretici CVE'sine bağlanmaz; bunlar "saldırı sınıfı/araştırma bulgusu" olarak kataloglanır. aLTeR ve ToRPEDO için kanonik makaleler ve yıllar bağımsız teyit edilmelidir (NDSS 2019 civarı). 5G alanı hızla değiştiğinden, en güncel durum için GSMA CVD bildirimleri ve akademik literatür §12'deki kaynaklardan izlenmelidir.

> Savunma özeti (hücresel): Bireysel düzeyde en güçlü tek hamle 2G'yi (mümkünse 3G'yi) cihazda kapatıp downgrade yüzeyini daraltmaktır; bu, A5 kriptosunu, en kolay IMSI-catcher senaryolarını ve en kolay dinlemeyi saha dışına atar. Modern ağda 5G SA + doğru uygulanmış SUCI kimlik ifşasını azaltır. Uygulama düzeyinde uçtan-uca şifreleme (Signal/TLS/VPN) içerik gizliliğini ağdan bağımsızlaştırır. Ayrıntı: Bölüm 20 §16, Bölüm 6 §5-6.

---

<a id="5"></a>
## 5. Telekom Çekirdeği: SS7 ve Diameter Sinyalleşme Zafiyetleri

Hava arayüzünün ötesinde, operatörler arası sinyalleşme ağı (SS7 ve onun 4G karşılığı Diameter) tarihsel olarak "kapalı, güvenilir taraflar" varsayımıyla tasarlanmıştır. Bu güven varsayımı bugün geçerli değildir; ağa (yetkisiz ya da kötüye kullanılan operatör/aracı erişimiyle) ulaşan bir taraf, abone konumu sorgulama, çağrı/SMS yönlendirme ve kimlik doğrulama vektörü talebi gibi işlemleri kötüye kullanabilir. Derin işleyiş Bölüm 20 §12 ve Bölüm 6 §5'tedir.

| Zafiyet sınıfı | Referans | Sistem | Keşif/açıklama | Prensip (neden savunmasız) | Etki | Durum | Savunma |
|---|---|---|---|---|---|---|---|
| SS7 konum/sorgu istismarı | (sınıf; kamuya açık 2014 CCC sunumları) | SS7 (2G/3G çekirdek sinyalleşme) | kamuoyuna 2014 | Ağ, sorgu kaynağının meşruluğunu doğrulamaz; konum (örn. MAP anyTimeInterrogation/provideSubscriberLocation) ve abone bilgisi sorgulanabilir | Konum izleme, abone bilgisi ifşası | Sürüyor; SS7 firewall/filtre ile azaltılır | SS7 firewall, kategori bazlı filtreleme, GSMA FS.11 |
| SS7 çağrı/SMS ele geçirme | (sınıf) | SS7 (MAP) | 2014+ | Yönlendirme mesajlarıyla (örn. updateLocation/SRI-SM kötüye kullanımı) çağrı/SMS saldırgana yönlendirilebilir; SMS-OTP yakalama riski | SMS-tabanlı 2FA atlatma, dinleme | Sürüyor; filtre + SMS-OTP'den uzaklaşma | SS7 firewall; SMS-OTP yerine uygulama-tabanlı/donanım 2FA |
| Diameter (4G) istismarı | (sınıf; GSMA FS.19 rehberliği) | Diameter (LTE çekirdek sinyalleşme) | araştırma 2010'lar sonu | Diameter, SS7'nin sorunlarının çoğunu miras alır; yanlış yapılandırma/filtre eksikliğinde konum/abone sorguları ve DoS mümkün | Konum izleme, abone ifşası, DoS | Sürüyor; FS.19'a göre sertleştirme | Diameter Edge Agent (DEA) filtreleme, GSMA FS.19 |
| GTP / roaming arayüzü | (sınıf; GSMA FS.20) | GTP (kullanıcı/kontrol tüneli, roaming) | araştırma | Yanlış yapılandırılmış GTP arayüzleri abone bilgisi ifşası, sahtecilik ve DoS'a açabilir | Bilgi ifşası, sahtecilik, DoS | Sürüyor; FS.20 sertleştirme | GTP firewall; GSMA FS.20 rehberi |
| SEPP/HTTP-2 (5G roaming) | (5G; gelişen alan) | 5G SBA roaming arası (N32/SEPP) | 5G ile yeni | 5G, roaming arası güveni SEPP ile kriptografik olarak sınırlamayı hedefler; yanlış uygulama/yapılandırma yeni yüzey açabilir | Roaming arası ifşa/sahtecilik (yanlış yapılandırmada) | Yeni; standart koruma var, uygulama kalitesine bağlı | SEPP doğru dağıtımı; N32 koruması; GSMA güncel |

Notlar ve doğruluk uyarıları: SS7/Diameter/GTP istismarları neredeyse hiçbir zaman tekil "CVE" değil, GSMA dokümanlarında (FS.11 SS7, FS.19 Diameter, FS.20 GTP) ele alınan "kötüye kullanım sınıfları"dır; somut numara aramak yerine bu rehberlere ve GSMA CVD sürecine bakılmalıdır. Tarihsel kamuoyu farkındalığı 2014 Chaos Communication Congress sunumlarıyla (Tobias Engel, Karsten Nohl) artmıştır; tam atıf ve tarih bağımsız teyit edilebilir.

> Savunma özeti (çekirdek): Operatör tarafında sinyalleşme firewall'ı (SS7/Diameter/GTP filtreleme, GSMA FS serisi), abone tarafında ise en kritik pratik ders SMS-tabanlı 2FA'dan uzaklaşmaktır — SMS-OTP, SS7 yönlendirmesiyle yakalanabildiği için uygulama-tabanlı (TOTP) ya da donanım (FIDO2/güvenlik anahtarı) 2FA'ya geçmek bireysel en yüksek getirili savunmadır. Ayrıntı: Bölüm 20 §12-13.

---

<a id="6"></a>
## 6. GNSS Zafiyetleri ve Gerçek-Dünya Kesintileri

GNSS'in temel zafiyeti tasarımsaldır: sivil sinyaller son derece zayıf (yer seviyesinde gürültü tabanının altında, ~-130 dBm) ve çoğunlukla kimlik-doğrulamasızdır. Bu, hem bastırma (jamming) hem de aldatma (spoofing) için zemin oluşturur. Bu bölüm olayları ve prensibi kataloglar; mekanizma ve savunma Bölüm 10 §10-12'de, RF tehdit çerçevesi Bölüm 13'tedir.

| Zafiyet / olay sınıfı | Referans | Sistem | Dönem | Prensip (neden savunmasız) | Etki | Durum | Savunma |
|---|---|---|---|---|---|---|---|
| GNSS jamming (bastırma) | (sınıf; sayısız saha olayı) | GPS/GLONASS/Galileo/BeiDou L-bandı | süregelen | Sinyal çok zayıf olduğundan görece düşük güçlü bir gürültü kaynağı alıcıyı kilitlenmeden çıkarır; kişisel jammer'lar (yasadışı) bile yerel kesinti yapar | Konum/zaman kaybı; havacılık, denizcilik, zaman senkron etkisi | Sürüyor; bölgesel kesintiler raporlanıyor | Anten null-steering, çok-bantlı/çok-takım, atalet (INS) yedeği, izleme |
| GNSS spoofing (aldatma) | (sınıf; gösterimler ve saha olayları) | Sivil GPS C/A ve diğer açık kodlar | gösterimler 2010'lar, saha olayları sürüyor | Açık kod ve zamanlama bilinir; kimlik doğrulama yoktur; sahte sinyal gerçeği bastırıp alıcıyı yanlış konum/zamana çeker | Yanlış konum/zaman; navigasyon ve zaman-bağımlı sistem riski | Sürüyor; bölgesel/çatışma bölgesi olayları | Galileo OSNMA gibi kimlik doğrulama, çoklu-takım çapraz kontrol, INS, anomali tespiti |
| Zamanlama saldırısı (GNSS time) | (sınıf) | GNSS'e bağımlı zaman dağıtımı (telekom, enerji, finans) | süregelen | GNSS bir saat dağıtım sistemidir; sahte zaman, ağ senkronizasyonunu ve zaman damgalı sistemleri bozar | Telekom/enerji/finans senkron hatası | Sürüyor; kritik altyapı riski | Holdover osilatör, çoklu zaman kaynağı (PTP/eLoran), GNSS bütünlük izleme |

Notlar ve doğruluk uyarıları: GNSS jamming/spoofing olayları çoğunlukla CVE değil olay (incident) kategorisindedir; bölgesel kesinti raporları havacılık otoriteleri, denizcilik bültenleri ve açık izleme platformlarınca (Bölüm 10 ve 13) belgelenir. Galileo OSNMA (Open Service Navigation Message Authentication) sivil tarafta sinyal-mesajı kimlik doğrulaması getiren önemli bir savunma adımıdır; kapsama ve olgunluk durumu resmî Galileo/GSC kaynaklarından teyit edilmelidir.

> Savunma özeti (GNSS): Tek bir GNSS'e körü körüne güvenme — çoklu-takım + çoklu-bant, kimlik-doğrulamalı sinyal (OSNMA), atalet/holdover yedeği ve anomali tespiti (ani konum sıçraması, imkânsız hız, C/N0 anormalliği) katmanlı savunmayı oluşturur. Kritik zaman uygulamalarında GNSS'i tek kaynak yapma. Ayrıntı: Bölüm 10 §12, Bölüm 13 §12.

---

<a id="7"></a>
## 7. IoT / Sub-GHz / Otomotiv: Rolling-Code, KeeLoq, Crypto-1, Zigbee/Z-Wave, LoRaWAN

Bu kategori, ev/araç/endüstri çevresindeki kısa ve orta menzilli protokolleri kapsar. Ortak tema: maliyet ve enerji kısıtları yüzünden zayıf ya da hatalı uygulanmış kriptografi. Derin işleyiş Bölüm 16 §6-10'da, otomotiv Bölüm 16 §10 ve Bölüm 13 §9'da.

| Zafiyet / ad | CVE / referans | Sistem / protokol | Keşif | Prensip (neden savunmasız) | Etki | Durum | Savunma |
|---|---|---|---|---|---|---|---|
| KeeLoq zayıflığı | (akademik 2007-2008; CVE'siz) | KeeLoq blok şifre (garaj/araç uzaktan kumanda) | 2007-2008 | Yan-kanal ve kriptanaliz, üreticiye özgü ana anahtar çıkarımı ve bazı uygulamalarda klonlama imkânı verir | Kumanda klonlama, yetkisiz erişim | Tarihsel/uygulama-bağımlı; yeni sistemler daha güçlü | AES tabanlı modern rolling-code; üretici güncel ürün |
| RollJam (rolling-code yakalama-tutma) | (teknik sınıfı 2015; CVE'siz) | Sabit/zayıf rolling-code RKE (araç/garaj) | 2015 | Saldırgan geçerli kodu yakalayıp bloklayarak "biriktirir"; sonra kullanır. Zayıf rolling-code pencere yönetimi bunu mümkün kılar | Tek seferlik yetkisiz açma | Uygulama-bağımlı; zaman/sayaç sıkı yönetimle azalır | Zaman damgalı/challenge-response RKE; UWB mesafe ölçümü |
| Rolling-PWN | (2022 kamuoyu; üretici-bağımlı, CVE durumu teyit edilmeli) | Belirli araç RKE rolling-code uygulamaları | 2022 | İddiaya göre belirli araçlarda kod sayaç penceresi yeniden senkronize edilerek eski kodların yeniden geçerli kılınması mümkün | Tekrar oynatma ile açma (etkilenen modellerde) | Üretici-bağımlı; doğrulama gerektirir | Üretici firmware/güncelleme; sıkı sayaç penceresi |
| MIFARE Classic / Crypto-1 | (akademik 2008; CVE'siz, tasarım) | MIFARE Classic RFID (13.56 MHz HF) | 2008 | Tescilli Crypto-1 akış şifresi ve 48-bit anahtar tersine mühendislikle kırıldı; zayıf rastgelelik ve kısa anahtar pratik çözüm sağlar | Kart klonlama, erişim kontrolü atlatma | Tarihsel/sürüyor (hâlâ sahada çok kart var) | DESFire/AES kartlara geçiş; kart+arka uç doğrulama |
| Zigbee anahtar/eşleştirme zayıflıkları | (sınıf; çeşitli bulgular) | Zigbee (802.15.4) ağ/birleşme | araştırma | Bazı uygulamalarda birleşme (join) anında ağ anahtarı zayıf korunarak iletilebilir; varsayılan/iyi bilinen anahtarlar | Ağ anahtarı ele geçirme, cihaz enjeksiyonu | Uygulama-bağımlı; sürüm sertleştirme | Install-code tabanlı eşleştirme, güncel Zigbee 3.0, segmentasyon |
| Z-Wave (S0 sınıfı) zayıflığı | (sınıf; "S0 downgrade" tartışması) | Z-Wave eski S0 güvenlik sınıfı | araştırma 2018 civarı | Eski S0 anahtar değişimi zayıftır; bazı senaryolarda S2'den S0'a downgrade ile zafiyet | Anahtar ele geçirme (downgrade'te) | S2 ile büyük ölçüde çözüldü | S2 güvenlik sınıfını zorla; S0'ı sınırla |
| LoRaWAN nonce/anahtar zayıflıkları | (sınıf; 1.0.x'te bilinen) | LoRaWAN birleşme (join) ve oturum anahtarları | araştırma | 1.0.x'te DevNonce yeniden kullanımı/öngörülebilirliği ve ABP'de statik anahtarlar replay/oturum riskine açar | Replay, oturum ele geçirme (zayıf uygulamada) | 1.1 büyük ölçüde iyileştirdi | LoRaWAN 1.1, OTAA, nonce yönetimi, anahtar rotasyonu |

Notlar ve doğruluk uyarıları: Bu kategorideki bulguların çoğu akademik makale veya konferans gösterimidir, klasik CVE değil. "Rolling-PWN" kamuoyuna 2022'de duyurulmuş, üreticiye özel ve tartışmalı bir iddiadır; etkilenen modeller, CVE durumu ve geçerlilik bağımsız olarak (üretici advisory + NVD) teyit edilmelidir — bu katalog onu doğrulanmamış-iddia olarak işaretler. KeeLoq, Crypto-1 ve RollJam tarihsel olarak iyi belgelenmiştir ama belirli ürün etkisi üreticiye göre değişir.

> Savunma özeti (IoT/sub-GHz/otomotiv): Zayıf/tescilli kriptodan (Crypto-1, eski KeeLoq, S0) modern standart kriptoya (AES kart, S2, LoRaWAN 1.1) geç; rolling-code yerine challenge-response ve mümkünse UWB mesafe ölçümü tercih et (relay'e karşı); IoT cihazlarını ağ segmentasyonuyla izole et ve varsayılan anahtarları değiştir. Ayrıntı: Bölüm 16 §13.

---

<a id="8"></a>
## 8. RFID / NFC: Kart Klonlama ve Relay Zafiyetleri

RFID/NFC zafiyetleri iki köke iner: zayıf/eski kart kriptosu (klonlama) ve mesafe varsayımının kırılması (relay). Crypto-1 bir önceki bölümde geçti; burada NFC'ye özgü relay ve klonlama sınıfı toparlanır. Derin işleyiş Bölüm 16 §4-6'da.

| Zafiyet sınıfı | Referans | Sistem | Prensip (neden savunmasız) | Etki | Durum | Savunma |
|---|---|---|---|---|---|---|
| Kart klonlama (zayıf kart) | (Crypto-1 vb.) | MIFARE Classic ve benzeri eski kartlar | Zayıf/kırık kripto kartın tüm içeriğinin okunup birebir kopyalanmasına izin verir | Yetkisiz erişim, kimlik taklidi | Sürüyor (eski kart tabanı) | AES/DESFire kart; arka uç doğrulama; anti-clone |
| UID-only sistemler | (tasarım) | Yalnızca UID'ye güvenen erişim sistemleri | UID kimlik doğrulama için değildir; bazı kartlarda yazılabilir/taklit edilebilir | Trivial taklit | Sürüyor (kötü tasarımda) | UID'yi kimlik sanma; kriptografik karşılıklı doğrulama |
| NFC/temassız relay | (akademik sınıf) | Temassız ödeme/erişim (kısa menzil varsayımı) | Saldırgan, kart ile okuyucu arasındaki mesajları gerçek zamanlı röleler; fiziksel yakınlık varsayımı kırılır | Yetkisiz işlem/erişim mesafe varsayımına rağmen | Sürüyor; mesafe-bağlama (distance bounding) ile azaltılır | Distance bounding, zaman aşımı sıkılaştırma, kullanıcı onayı |

Notlar: RFID/NFC relay ve klonlama büyük ölçüde "saldırı sınıfı" olarak belgelenir; ürün-spesifik CVE'ler kart/okuyucu üreticisine göre değişir. Distance bounding'in pratik dağıtımı protokol ve donanıma bağlıdır; teyit Bölüm 16 ve üretici dokümanından.

> Savunma özeti (RFID/NFC): UID'yi asla kimlik doğrulama olarak kullanma; kriptografik karşılıklı doğrulamalı modern kartlara geç; relay'e karşı zaman aşımı/mesafe-bağlama ve yüksek değerli işlemlerde kullanıcı onayı uygula. Ayrıntı: Bölüm 16 §4-6.

---

<a id="9"></a>
## 9. Donanım, Firmware ve SDR Yan-Kanal Zafiyet Kavramı

Protokol ve kripto açıklarının altında bir katman daha vardır: donanımın ve firmware'in kendisi. Bu kategori tek tek CVE'lerden çok kavramsal bir farkındalık gerektirir, çünkü RF güvenliğinin sessiz ama kritik tarafıdır. Üç alt başlık öne çıkar.

Birincisi firmware zafiyetleri. WiFi/Bluetooth/baseband yongalarının firmware'i, kendi başına bir saldırı yüzeyidir. Baseband (hücresel modem) firmware'inde uzaktan tetiklenebilir bellek hataları, en ciddi mobil zafiyet sınıflarından biri olmuştur (hava arayüzünden kullanıcı etkileşimsiz RCE potansiyeli). Bu açıklar üreticiye özel CVE'lerle yamanır; tehlikesi, baseband'in ana işlemciden ayrı ve görece az denetlenen bir ortam olmasıdır. Savunma: güvenlik güncellemelerini geciktirmeden uygulamak ve güncellenmeyen (EOL) cihazları kademeli emekliye ayırmaktır.

İkincisi yan-kanal sızıntısı (TEMPEST/emanasyon). Bir cihazın istemeden yaydığı elektromanyetik emisyon, işlediği veri hakkında bilgi sızdırabilir (Bölüm 17'nin konusu). Bu bir "CVE" değil, fiziksel bir gerçektir; savunması ekranlama, filtreleme ve TEMPEST sertleştirmesidir. Kriptografik yan kanal (zamanlama, güç analizi) da bu aileye girer: Dragonblood'un (§2) yan-kanal vektörü, bunun kablosuz protokoldeki yansımasıdır.

Üçüncüsü SDR'nin iki yüzü. SDR (yazılım tanımlı radyo), hem savunma/araştırma aracıdır hem de düşük maliyetli saldırı yüzeyini büyütür: önceden pahalı/özel donanım gerektiren analiz ve test artık genel amaçlı SDR ile yapılabilir. Bu, savunmacı için "saldırganın eşiği düştü" anlamına gelir — tehdit modellemesinde maliyet varsayımlarını güncellemek gerekir. Aynı SDR, yetkili test için (§13) en değerli araçtır.

| Alt başlık | Doğa | Tipik örnek/kaynak | Savunma |
|---|---|---|---|
| Baseband/modem firmware | Uzaktan tetiklenir kod hatası | Üreticiye özel CVE'ler (Qualcomm/MediaTek/vb.) | Güvenlik yaması; EOL cihaz emekliliği |
| WiFi/BT yonga firmware | Kod/kripto hatası | Kr00k, SweynTooth (firmware düzeyi) | Sürücü/firmware güncel |
| TEMPEST/emanasyon | Fiziksel sızıntı | Bölüm 17 | Ekranlama, filtreleme, mesafe |
| Kripto yan kanal | Zamanlama/güç sızıntısı | Dragonblood yan-kanal | Sabit-zaman kod, maskeleme |
| SDR ile düşen saldırı eşiği | Tehdit modeli değişimi | Genel amaçlı SDR | Maliyet varsayımını güncelle; tespit |

> Sezgi: Donanım/firmware katmanı, "protokol güvenli ama cihaz değil" durumunun kaynağıdır. En iyi tasarlanmış protokol bile, altındaki firmware yamasız ise sömürülebilir. Bu yüzden yama yönetimi ve EOL politikası, kripto seçimi kadar önemli bir savunma kararıdır.

---

<a id="10"></a>
## 10. Yıllara Göre Zaman Çizelgesi (Büyük Tablo)

Aşağıdaki tablo, kablosuz/telekom zafiyet manzarasının kronolojik omurgasıdır. Amaç, "ne zaman ne kırıldı ve bugün ne durumda" sorusuna tek bakışta yanıt vermektir. Tarihler keşif/kamuoyuna açıklama yılıdır; CVE ve tam tarih için ilgili satırın §2-9'daki kaydına ve NVD'ye bakılmalıdır.

| Yıl | Zafiyet / olay | Sistem | Etki özeti | Bugünkü durum |
|---|---|---|---|---|
| 1994-2003 | A5/2, A5/1 kriptanalizi | GSM 2G | 2G şifreleme çözülür | Tarihsel; 2G kapatılıyor |
| 2001 | WEP (FMS) zayıflığı | 802.11 WEP | WEP anahtarı kurtarılır | Tarihsel; WEP terk |
| 2007-2008 | KeeLoq, Crypto-1 (MIFARE) | RKE / RFID | Klonlama/erişim atlatma | Tarihsel/sürüyor (eski taban) |
| 2008-2009 | WPA-TKIP (Beck-Tews) | WPA | Kısmi paket manipülasyonu | Tarihsel; TKIP kapat |
| 2011 | WPS PIN zayıflığı | WiFi WPS | Router PIN/anahtar kurtarma | Sürüyor (eski router) |
| 2014 | SS7 kötüye kullanımı kamuoyu | SS7 çekirdek | Konum/SMS/çağrı riski | Sürüyor; firewall ile azaltma |
| 2015 | RollJam (rolling-code) | Araç/garaj RKE | Yetkisiz açma | Uygulama-bağımlı |
| 2017 | KRACK | WPA/WPA2 | Anahtar yeniden kurulum, çözme | Yamandı |
| 2017 | BlueBorne | Bluetooth (çoklu OS) | Uzaktan kod/MITM | Yamandı |
| 2018 | Z-Wave S0 tartışması | Z-Wave | Downgrade ile anahtar riski | S2 ile çözüldü |
| 2019 | KNOB | Bluetooth BR/EDR | Anahtar entropisi düşürme | Yamandı |
| 2019 | Kr00k (CVE-2019-15126) | WiFi yongaları | Sıfır anahtarla çözme | Yamandı |
| 2019 | Dragonblood | WPA3-SAE | Yan kanal/downgrade, parola | Büyük ölçüde yamandı |
| 2019 | aLTeR, ToRPEDO | LTE/4G-5G | Yönlendirme, paging takibi | LTE'de doğası; 5G azaltır |
| 2020 | BIAS (CVE-2020-10135) | Bluetooth BR/EDR | Cihaz taklidi, MITM | Yamandı |
| 2020 | SweynTooth | BLE SoC'ler | DoS/atlatma | Yamandı |
| 2020 | BleedingTooth | Linux BlueZ | Çekirdek RCE/sızıntı | Yamandı |
| 2020 | BLESA | BLE yeniden bağlanma | Sahte veri enjeksiyonu | Kısmen; uygulama-bağımlı |
| 2021 | FragAttacks | 802.11 frag/aggr | Enjeksiyon/sızıntı | Kısmen yamandı |
| 2022 | Rolling-PWN (iddia) | Araç RKE (belirli) | Tekrar oynatma (iddia) | Üretici-bağımlı; teyit gerek |
| 2019-sürüyor | 5G SUCI/downgrade araştırmaları | 5G NR | Kalan ifşa/DoS | Araştırma; gelişiyor |
| süregelen | GNSS jamming/spoofing | GPS/GNSS | Konum/zaman kaybı/aldatma | Sürüyor; OSNMA vb. azaltır |

> Okuma notu: Tablonun üst yarısı (tarihsel) "kripto/tasarım kırıldı, sistem değişti" hikâyesidir; alt yarısı (2017+) "protokol/uygulama hatası bulundu, yamandı ama saha gecikiyor" hikâyesidir. İki desen farklı savunma gerektirir: birincisi mimari değişim (eski sistemi kapat), ikincisi yama disiplini.

---

<a id="11"></a>
## 11. Açığın Yaşam Döngüsü: Keşiften KEV'e (ASCII)

Bir zafiyetin doğuşundan kapanışına giden yol standart bir döngü izler. Bu döngüyü anlamak, hem "neden bir açık hâlâ tehlikeli" hem de "yamayı ne zaman öncelemeliyim" sorularını yanıtlar.

```
   KEŞİF                 SORUMLU AÇIKLAMA (CVD)        KOORDİNASYON / CVE
 ┌─────────┐            ┌──────────────────────┐     ┌──────────────────┐
 │araştırmacı│──bildirim─▶│ üretici/CERT bilgilen.│──▶│ CVE numarası atanır│
 │ açığı     │           │ embargo süresi        │    │ (MITRE/CNA)        │
 │ bulur     │           │ yama geliştirilir     │    │ NVD'ye işlenir     │
 └─────────┘            └──────────────────────┘     └──────────────────┘
       │                          │                           │
       │ (paralel risk)           ▼                           ▼
       │                  ┌──────────────┐            ┌──────────────────┐
       └─ sıfır-gün ─────▶│ YAMA YAYIMI  │───────────▶│ SAHA DAĞITIMI    │
          (yama yokken     │ (advisory)   │            │ (cihazlar güncel │
           riskli pencere) └──────────────┘            │  oldukça kapanır)│
                                  │                    └──────────────────┘
                                  ▼                            │
                          ┌──────────────────┐                 │
                          │ AKTİF SÖMÜRÜ?     │◄────────────────┘
                          │ → CISA KEV'e girer│  (yama varken hâlâ
                          │ (bilinen-sömürülen)│   sömürülüyorsa)
                          └──────────────────┘
                                  │
                                  ▼
                       n-gün penceresi: yama VAR ama
                       saha güncellenmediği için açık SÖMÜRÜLÜR
```

Döngünün kritik kavramları ve savunma anlamı:

| Aşama | Ne olur | Savunma açısından anlamı |
|---|---|---|
| Keşif | Araştırmacı açığı bulur | Henüz kimse bilmiyor olabilir; "sıfır-gün" buradan doğar |
| CVD (sorumlu açıklama) | Üreticiye gizlice bildirilir, embargo + yama | Etik norm; advisory bu süreçten çıkar |
| CVE atanması | Standart kimlik (CVE-YYYY-NNNN) verilir | Takip ve eşleştirme bu numarayla yapılır |
| Yama yayımı | Üretici düzeltmeyi yayımlar | Tehdit teknik olarak kapanır (üretici tarafında) |
| Saha dağıtımı | Cihazlar güncellenir | Asıl darboğaz burası; gecikme = risk |
| KEV (bilinen-sömürülen) | Aktif sömürülüyorsa CISA KEV listesine girer | En yüksek yama önceliği işareti |

Sıfır-gün (zero-day) ile n-gün (n-day) ayrımı bu döngünün kalbidir. Sıfır-gün, yama henüz yokken (üreticinin "sıfırıncı günü") sömürülen açıktır; savunması zordur, çünkü imza/yama yoktur — katmanlı savunma ve anomali tespitine bağımlısın. n-gün ise yama yayımlandıktan N gün sonra hâlâ sömürülen açıktır; teknik çözümü mevcuttur, sorun dağıtım gecikmesidir. Pratikte ihlallerin büyük kısmı sıfır-gün değil, n-gün'dür: yaması aylar/yıllar önce çıkmış ama uygulanmamış açıklar. Bu yüzden CISA KEV listesi, "yaması var, sömürülüyor, hemen yama" diyen en pragmatik önceliklendirme aracıdır.

> Neden eski açıklar hâlâ sömürülür: Yama gecikmesi (saha dağıtımı), EOL/güncellenemeyen cihazlar (eski router, IoT, eski telefon), geriye uyumluluk (2G/WPS gibi zayıf modların açık bırakılması) ve "çalışıyorsa dokunma" kültürü. Katalogdaki "sürüyor" etiketlerinin çoğu, protokol/üretici tarafında değil saha tarafında açıktır. Savunmanın gerçek cephesi burasıdır.

---

<a id="12"></a>
## 12. Nasıl Güncel Kalınır: Kaynak Disiplini

Bir zafiyet kataloğu, yazıldığı an eskimeye başlar. Bu yüzden katalogdan daha değerli olan, kataloğu güncel tutan kaynak disiplinidir. Bölüm 14 bu konuyu derinlemesine işler; burada RF/telekom odaklı çekirdek kaynaklar özetlenir.

| Kaynak türü | Örnek | Ne için | Güncellik |
|---|---|---|---|
| CVE/zafiyet veritabanı | NVD (NIST), MITRE CVE | Kanonik CVE kaydı, CVSS, etkilenen ürün | Sürekli |
| Bilinen-sömürülen liste | CISA KEV (Known Exploited Vulnerabilities) | "Hemen yama" önceliği | Sık güncellenir |
| Üretici advisory | Yonga/OS/router üreticisi güvenlik bültenleri | Yama varlığı ve sürüm | Ürüne bağlı |
| Telekom CVD | GSMA Coordinated Vulnerability Disclosure, FS serisi | Operatör/şebeke zafiyetleri | Periyodik |
| Akademik | USENIX Security, IEEE S&P, NDSS, ACM WiSec | Yeni saldırı sınıfları (çoğu CVE-öncesi) | Konferans takvimi |
| Konferans | DEF CON, Black Hat, Chaos Communication Congress (CCC) | Pratik gösterim ve araç | Yıllık |
| Tehdit istihbaratı | MITRE ATT&CK, MISP feed'leri, OSINT | TTP eşleme ve göstergeler | Sürekli |

İş akışı önerisi (savunma analisti için): (1) NVD/CVE'yi ürün envanterinle eşle (hangi CVE seni ilgilendiriyor); (2) CISA KEV'i öncelik filtresi olarak kullan (sömürülen + sana uyan = en üst); (3) üretici advisory'sinden yama varlığını teyit et; (4) telekom/şebeke tarafında GSMA FS rehberlerini izle; (5) yeni saldırı sınıfları için akademik/konferans çıktısını takip et (bunlar genelde CVE'den önce gelir); (6) MITRE ATT&CK ve MISP ile gözlemleri TTP'ye ve göstergeye bağla. Kritik kural: ikincil habere değil birincil kaynağa (CVE kaydı, advisory, makale) dayan; bir CVE numarasını her zaman NVD'den teyit et. Ayrıntı: Bölüm 14 §1, §3-5, §8.

---

<a id="13"></a>
## 13. Yetkili Test Araçları (Zafiyet Sınıfına Göre)

Aşağıdaki araçlar yalnızca yetkili bağlamda anılır: kendi cihazların, kendi laboratuvarın, kendi ağın ya da açık yazılı izinli (kapsamı tanımlı) sızma testi. Bu araçların çoğu çift kullanımlıdır (savunma ve saldırı); buradaki çerçeve savunma doğrulaması, yama teyidi ve güvenlik araştırmasıdır. Başkasının sistemine yönelik kullanım yasa dışıdır. Araçların yetenek/komut ayrıntısı bilinçli olarak verilmez; amaç "hangi sınıf için hangi araç ekosistemi" eşlemesidir.

| Zafiyet sınıfı | Açık-kaynak test aracı (yetkili) | Test bağlamı |
|---|---|---|
| WiFi (handshake/WEP/WPA) | aircrack-ng suiti | Kendi ağının el-sıkışma yakalama/çözüm direncini ölçme |
| WiFi (PMKID/hash toplama) | hcxtools / hcxdumptool | Kendi AP'nin PMKID maruziyetini denetleme |
| WiFi (parola direnci) | hashcat / John the Ripper | Kendi parolanın çevrimdışı kırma direncini ölçme |
| WiFi (KRACK/FragAttacks teyidi) | araştırmacıların yayımladığı PoC test betikleri (örn. FragAttacks test aracı) | Kendi cihazının yamalı olup olmadığını doğrulama |
| WPA3-SAE (Dragonblood) | Dragonblood test betikleri / güncel hostapd-wpa_supplicant test modları | Kendi AP/istemcinin SAE sertleşmesini teyit |
| Bluetooth/BLE genel | Mirage (BLE saldırı/analiz çerçevesi), bettercap (BLE modülü) | Kendi BLE cihazının savunmasını sınama |
| BLE SoC (SweynTooth) | SweynTooth PoC çerçevesi | Kendi BLE SoC firmware'inin yamalı olup olmadığını test |
| BLE sniffing/analiz | nRF Sniffer, Ubertooth | Kendi BLE trafiğini pasif gözlem (eğitim) |
| RFID/NFC (Crypto-1/klon) | Proxmark3 (yazılımı), libnfc | Kendi kartının kripto sınıfını/klonlanabilirliğini denetleme |
| Sub-GHz/rolling-code | RTL-SDR/HackRF + Universal Radio Hacker (URH), rtl_433 | Kendi uzaktan kumandanın sinyal yapısını analiz |
| Zigbee/802.15.4 | KillerBee, Zigbee sniffer donanımı | Kendi Zigbee ağının anahtar/birleşme güvenliğini test |
| LoRaWAN | ChirpStack (kendi sunucu), LoRa sniffer | Kendi LoRaWAN dağıtımının nonce/anahtar hijyenini denetleme |
| Hücresel (pasif/lab) | srsRAN, Open5GS, gr-gsm (pasif/eğitim) | Faraday/kendi test hücresinde pasif analiz (Bölüm 20 §15) |
| GNSS (lab) | GNSS-SDR (pasif alıcı/araştırma) | Kendi alıcının anti-spoof davranışını lab'da inceleme |
| Genel RF analiz | GNU Radio, Universal Radio Hacker, SDR# / GQRX | Sinyal yapısı/spektrum analizi (eğitim/araştırma) |

> Etik sınır: Bu tabloda hiçbir araç "saldırı reçetesi" olarak değil, savunma doğrulaması/yama teyidi/eğitim için listelenmiştir. Hücresel iletim (sahte hücre kurma), başkasının trafiğini çözme veya kart klonlayıp kullanma çoğu ülkede ağır suçtur; srsRAN/Open5GS gibi araçlar yalnızca yalıtılmış (Faraday) kendi test hücresinde, yasal çerçevede kullanılmalıdır. "Kendi cihazın, kendi ağın, açık izin" üçlüsü dışına çıkma. Ayrıntı: Bölüm 15 §14, Bölüm 16 §14, Bölüm 20 §15, 17.

---

<a id="14"></a>
## 14. Alıştırmalar (Yasal, Kendi Cihazların)

Aşağıdaki alıştırmalar bilgi ve savunma odaklıdır; hiçbiri yetkisiz iletim, çözme veya başkasının sistemine erişim içermez.

1. Bir CVE'yi baştan sona izle (KRACK örneği). KRACK'i seç; NVD'de CVE-2017-13077 ve serisini bul, etkilenen ürünleri ve CVSS'i not et. Ardından zinciri kur: orijinal akademik makale (keşif) → üretici/CERT advisory'leri (CVD/yama) → CVE kaydı → bugün hangi cihazların hâlâ yamasız olabileceği. Yaşam döngüsünün (§11) her aşamasını bu tek açıkta somutlaştır. Çıktı: bir sayfalık "açık biyografisi".

2. Kendi cihazlarının yama durumunu denetle. Telefonun, router'ın ve bir IoT cihazının güncel firmware/yazılım sürümünü çıkar. Her biri için: üreticinin son güvenlik bülteni tarihi nedir, cihaz EOL mı, bilinen bir RF açığı (örn. router'da WPS açık mı, telefonda 2G kapanabiliyor mu) var mı? Çıktı: bir "kişisel yama açığı" tablosu ve önceliklendirilmiş kapatma listesi.

3. CISA KEV'den RF-ilgili bir kayıt analiz et. CISA KEV listesini aç; kablosuz/telekom/IoT ile ilişkili bir kayıt bul (ör. bir router/WiFi/baseband zafiyeti). Kaydın CVE'sini NVD'de incele: mekanizma sınıfı (§1) nedir, hangi katman, yama mevcut mu, neden hâlâ aktif sömürülüyor (n-gün analizi)? Çıktı: kaydın §1 dört-eksen çerçevesinde sınıflandırması.

4. Bir saldırı sınıfını "yama tarafında" haritalandır. FragAttacks'i seç; üç tasarım CVE'si ile uygulama CVE'lerini ayır (§2 notu), her birinin tasarım mı uygulama mı olduğunu işaretle ve hangisinin "yama ile", hangisinin "azaltma ile" çözüldüğünü ayırt et. Bu, "her açık yamayla bitmez" sezgisini pekiştirir.

5. Kendi WiFi'ını savunma çizgisine göre denetle (yetkili). Yalnızca kendi ağında: WPA3/WPA2 modunu, PMF durumunu, TKIP/WPS açık mı olduğunu router arayüzünden teyit et; §2 savunma özetiyle karşılaştır. Gerekirse aircrack-ng ile yalnızca kendi el-sıkışmanın çözülme direncini ölç (parola gücü teyidi). Çıktı: "katalogdaki hangi WiFi satırı bende kapalı, hangisi açık" eşlemesi.

6. Telekom 2FA yüzeyini sertleştir. SMS-tabanlı 2FA kullanan kritik hesaplarını listele (§5'in pratik dersi); SS7/Diameter yönlendirme riskini hatırlayarak bunları uygulama-tabanlı (TOTP) veya donanım (FIDO2) 2FA'ya taşımak için bir plan çıkar. Bu, ağ-katmanı bir zafiyetin bireysel savunmasını gösterir.

---

<a id="15"></a>
## 15. Hızlı Referans ve Diğer Bölümler

Bu bölüm bir indeks/zaman çizelgesi katmanıdır; derin işleyiş için ilgili bölümlere geç.

| Konu | Bu bölümde | Derin bölüm |
|---|---|---|
| WiFi zafiyetleri (KRACK, Kr00k, Dragonblood, FragAttacks, WPS) | §2 | Bölüm 15 |
| Bluetooth/BLE (BlueBorne, KNOB, BIAS, SweynTooth, BleedingTooth, BLESA) | §3 | Bölüm 16 §2-3 |
| Hücresel (A5, IMSI catcher, aLTeR, ToRPEDO, downgrade, 5G) | §4 | Bölüm 20, Bölüm 6 §5-6 |
| SS7/Diameter/GTP çekirdek | §5 | Bölüm 20 §12-13, Bölüm 6 §5 |
| GNSS jamming/spoofing | §6 | Bölüm 10 §10-12, Bölüm 13 §12 |
| IoT/sub-GHz/otomotiv (KeeLoq, RollJam, Crypto-1, Zigbee, Z-Wave, LoRaWAN) | §7 | Bölüm 16 §6-10 |
| RFID/NFC klon/relay | §8 | Bölüm 16 §4-6 |
| Donanım/firmware/yan-kanal | §9 | Bölüm 17, Bölüm 13 §11 |
| Zaman çizelgesi | §10 | Bölüm 6 §7 (kronoloji) |
| Açık yaşam döngüsü / CVD / KEV | §11 | Bölüm 14 §5 |
| Güncel kalma / kaynaklar | §12 | Bölüm 14 §1-8 |
| Yetkili test araçları | §13 | Bölüm 12 (araç ekosistemi), 15/16/20 alıştırmaları |
| RF tehdit/karşı-önlem genel | — | Bölüm 13 |

Çapraz kavram bağları: Bu bölümdeki her "saldırı sınıfı", Bölüm 6'nın manipüle-edilebilir/edilemez çerçevesiyle (replay/spoofing'i ne engeller), Bölüm 13'ün tehdit→tespit→savunma matrisiyle ve MITRE ATT&CK / MISP / OSINT akışlarıyla (Bölüm 14) eşlenebilir. Bir zafiyeti yalnızca "CVE" olarak değil, bir TTP, bir tespit imzası ve bir yama önceliği olarak görmek, bu kataloğu istihbarata çeviren bakıştır.

> Kapanış: Zafiyet manzarası durağan değildir; bu katalog bir fotoğraftır, harita değil. Değeri, ezberinde değil, ürettiği reflekstedir: yeni bir açık duyduğunda onu §1'in dört ekseninde sınıflandırmak, §11'in yaşam döngüsünde konumlandırmak, §12'nin kaynaklarından teyit etmek ve §13'ün yetkili çerçevesinde — yalnızca kendi sistemlerinde — doğrulamak. CVE numaraları ve tarihler için daima birincil kaynağa (NVD, üretici advisory, akademik makale) dönülmeli; bu bölüm emin olmadığı her yerde bunu açıkça not etmiştir.
