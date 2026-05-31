# 🛡️ SIGINT EL KİTABI — BÖLÜM 23: KABLOSUZ SALDIRI VEKTÖRLERİ TAKSONOMİSİ VE HİZMET-DIŞI BIRAKMA

## Erişim, Gizlilik, Bütünlük, Erişilebilirlik ve Kimlik — Kaynak/Paket Tüketme Saldırılarının Anatomisi, Tespiti ve Savunması

> Amaç: Önceki bölümler tek tek teknolojileri (WiFi — Bölüm 15, kısa menzilli/IoT — Bölüm 16, hücresel — Bölüm 20), aktif RF tehditlerini (karıştırma/spoofing — Bölüm 13) ve genel sinyal güvenliğini (Bölüm 6) ayrı ayrı işledi. Bu bölüm üzerlerine binen birleştirici bir soruyu ele alır: bütün bu saldırılar nasıl bir çatıya oturur? Bir analist, "şu an gördüğüm anomali hangi sınıfa girer, hangi güvenlik özelliğini hedefliyor, neresinden savunulur?" sorusunu ancak bir taksonomi varsa tutarlı yanıtlar. Kullanıcının somut sorusu — "başkasının internet paketini/kaynağını kullanan ve tüketen saldırılar" — bu çatının erişilebilirlik (DoS/kaynak tüketme) ve erişim (yetkisiz bağlanma) kollarının kesişiminde durur ve burada prensip, tespit ve savunma ekseninde işlenir. Hedef, bir saldırganın icra listesi değil; bir savunmacının, kritik altyapı analistinin ve ağ güvenliği mühendisinin tehdidi sınıflandırıp ölçebilmesidir.

> Yasal çerçeve (kritik): Bu bölüm anlama, tespit ve savunma amaçlıdır. Anlatılan her vektör için bilinçli olarak "nasıl çalışır + nasıl tespit edilir + nasıl savunulur" üçlüsü verilir; adım-adım operasyonel exploit, parola kırma reçetesi veya "şu komutla şu sistemin kaynağını tüket" tarzı icra talimatı **verilmez**. Başkasının ağına izinsiz bağlanmak, kaynağını/veri kotasını tüketmek, hizmetini kasten kesmek veya trafiğini araya girip değiştirmek Türkiye'de TCK 243 (bilişim sistemine yetkisiz erişim) ve TCK 244 (sistemi engelleme, bozma, verileri yok etme/değiştirme) kapsamında suçtur; ayrıca elektronik haberleşmenin gizliliğini ihlal ve karıştırma ayrı yaptırımlara tabidir. Bu bölümdeki tüm alıştırmalar **yalnızca kendi sahip olduğun ağ ve cihazlarda ya da yazılı yetki verilmiş bir sızma testi (pentest) kapsamında** yapılmalıdır. Belirli savunma yapılandırmalarının kendi sürümünüzdeki davranışı ve yasal sınırların güncel hali kaynaktan teyit edilmelidir.

---

## 📑 İÇİNDEKİLER

1. [Neden Taksonomi: Saldırıyı Sınıflandırmadan Savunulamaz](#1)
2. [CIA Üçlüsü ve Kablosuz Saldırı Yüzeyine Eşleme](#2)
3. [Saldırı Vektörü Taksonomisi: Beş Ana Kol](#3)
4. [Saldırı Zinciri: Keşif → Erişim → Kalıcılık → Amaç](#4)
5. [Kaynak ve Paket Tüketme Saldırıları (Kullanıcının Sorusu)](#5)
   - 5.1 [Yetkisiz WiFi Erişimi ile Başkasının Kotasını Kullanma](#5-1)
   - 5.2 [Deauth/Disassoc Seli ile Bağlantı Kesme](#5-2)
   - 5.3 [Ağ Katmanı Kaynak Tüketme: DHCP, ARP, Bağlantı Tablosu](#5-3)
   - 5.4 [RF Katmanı: Karıştırma ve Spektrum Tüketme](#5-4)
   - 5.5 [Enerji/Pil Tüketme: Uyku Yoksunluğu (IoT/BLE)](#5-5)
   - 5.6 [Bulut/Veri Kotası Tüketme: Botnet ve Proxy Kötüye Kullanımı](#5-6)
6. [Ortadaki-Adam (MitM) Sınıfı: Rogue AP, Evil Twin, SSL Stripping](#6)
7. [Enjeksiyon, Replay ve Spoofing Sınıfı](#7)
8. [İstemci-Taraflı Saldırılar: Probe İfşası, MAC İzleme, Captive Portal](#8)
9. [Tedarik Zinciri ve Firmware: Varsayılan Kimlik Bilgileri](#9)
10. [Araçlar ve Ne Ölçüde: Yetkili Test Bağlamı](#10)
11. [Tespit ve Savunma Mimarisi: WIDS/WIPS, PMF, WPA3, Sıfır-Güven](#11)
12. [Alıştırmalar (Yalnızca Kendi Ağında / Yetkili)](#12)
13. [Hızlı Referans ve Diğer Bölümler](#13)

---

<a id="1"></a>
## 1. Neden Taksonomi: Saldırıyı Sınıflandırmadan Savunulamaz

Güvenlik dünyasında en yaygın hata, saldırıları birbirinden kopuk "numaralar" olarak görmektir: deauth ayrı bir hile, evil twin ayrı bir hile, ARP zehirleme bambaşka bir şey. Bu parçalı bakış savunmacıyı zayıf düşürür çünkü her yeni saldırı varyantı tamamen yabancı görünür ve panik yaratır. Oysa saldırıların sayısı pratikte sonsuzdur; arkalarındaki **sınıf** sayısı azdır. Bir savunmacı sınıfı tanırsa, daha önce hiç görmediği bir varyantı bile doğru rafa koyabilir ve doğru kontrolü uygulayabilir.

Taksonomi, bu sınıflandırmanın disiplinidir. İyi bir taksonomi üç soruyu aynı anda yanıtlamamızı sağlar:

- **Bu saldırı neyi hedefliyor?** (Hangi güvenlik özelliği bozulmaya çalışılıyor — erişim mi, gizlilik mi, erişilebilirlik mi?)
- **Saldırı zincirinin hangi aşamasında?** (Keşif mi, ilk erişim mi, amaç icrası mı?)
- **Hangi katmanda çalışıyor?** (RF/fiziksel mi, bağlantı/MAC mi, ağ mı, uygulama mı?)

Bu üç eksen birbirine diktir. Aynı saldırı her üç eksende de bir koordinata sahiptir. Örneğin bir deauth seli: hedefi *erişilebilirlik* (DoS), zincirdeki yeri çoğu zaman *erişim için ön hazırlık* (istemciyi düşürüp yeniden el-sıkışmaya zorlamak) veya doğrudan *amaç* (sadece hizmeti kesmek), katmanı ise *bağlantı/MAC* (802.11 yönetim çerçevesi). Bu üç koordinatı bilince savunma kendiliğinden belirir: MAC katmanı yönetim çerçevesi koruması (802.11w/PMF) hedefin tam üstüne oturur.

```
                  ÜÇ EKSENLİ SINIFLANDIRMA UZAYI
                  (her saldırı bir koordinattır)

      Katman                    Hedef (CIA + ek)
   (nerede çalışır)          (neyi bozar)
        ▲                          ▲
   Uygulama │                      │ Kimlik (Authenticity)
     Ağ     │                      │ Erişilebilirlik (Availability)
   Bağlantı │      ● deauth        │ Bütünlük (Integrity)
    RF/Fiz  │     (Bağlantı,       │ Gizlilik (Confidentiality)
        │   │      Erişilebilirlik,│ Erişim (Authorization)
        │   │      Erişim-öncesi)  │
        └───┴──────────────────────┴──────────► Zincir aşaması
            Keşif  Erişim  Kalıcılık  Amaç
```

Not: Bu bölüm boyunca her vektörü bu üç eksene yerleştireceğiz. Amaç ezber değil; bir anomali gördüğünde onu refleks olarak bu uzaya oturtabilmen. Bu beceri, henüz isimlendirilmemiş yeni saldırılar karşısında bile seni savunulabilir kılar.

### Taksonominin savunmaya pratik faydası

Savunma bütçesi sonludur. Her saldırı varyantına ayrı bir kontrol almak hem pahalı hem imkânsızdır. Taksonomi, kontrolleri **sınıf düzeyinde** seçmeyi sağlar: tek bir 802.11w/PMF kontrolü, bütün yönetim-çerçevesi tabanlı DoS sınıfını (deauth, disassoc, sahte beacon temelli bazı saldırılar) aynı anda zayıflatır. Tek bir WPA3-SAE geçişi, çevrimdışı parola tahmini sınıfını büyük ölçüde kapatır. Mühendislik için kritik sezgi şudur: **savunmayı saldırı sayısıyla değil, sınıf sayısıyla ölçeklendir.**

---

<a id="2"></a>
## 2. CIA Üçlüsü ve Kablosuz Saldırı Yüzeyine Eşleme

Bilgi güvenliğinin klasik temeli CIA üçlüsüdür: Confidentiality (Gizlilik), Integrity (Bütünlük), Availability (Erişilebilirlik). Kablosuz/RF bağlamında bu üçlü iki ek eksenle genişler: erişim/yetkilendirme (Authorization — kimin ağı kullanmaya hakkı var) ve kimlik/özgünlük (Authenticity — karşı tarafın iddia ettiği varlık olup olmadığı). Beş eksen birlikte, kablosuz saldırı yüzeyinin tamamını kaplar.

| Güvenlik özelliği | Soru | Kablosuz saldırı örneği | Hedeflenen değer |
|---|---|---|---|
| Erişim (Authorization) | Kim bağlanabilir? | Yetkisiz WiFi bağlanma, parola zayıflığı istismarı, açık ağa sızma | Bant genişliği, veri kotası, iç ağ erişimi |
| Gizlilik (Confidentiality) | Veriyi kim okuyabilir? | Pasif dinleme, el-sıkışma yakalama, MitM, SSL stripping | Mahremiyet, kimlik bilgileri, içerik |
| Bütünlük (Integrity) | Veri değiştirildi mi? | Paket enjeksiyonu, replay, içerik değiştirme | Verinin doğruluğu, işlem güvenliği |
| Erişilebilirlik (Availability) | Hizmet ayakta mı? | Deauth seli, karıştırma, DHCP tükenmesi, kaynak tüketme | Hizmet sürekliliği, kullanılabilirlik |
| Kimlik (Authenticity) | Karşı taraf gerçek mi? | Evil twin, rogue AP, MAC spoofing, sahte baz istasyonu | Güven ilişkisi, kimlik |

Kullanıcının sorduğu "başkasının paketini/kaynağını tüketen saldırılar" iki eksende oturur: bir başkasının bağlantısına izinsiz bağlanıp **kotasını harcamak** erişim ihlalidir; hizmeti kasten kesmek veya altyapıyı boğmak ise **erişilebilirlik** saldırısıdır. Bu ikisi sık birlikte kullanılır: erişilebilirliği geçici olarak bozmak (deauth), çoğu zaman erişim elde etmenin (el-sıkışmayı yeniden tetikleyip yakalamak) ön adımıdır. Bölüm 5 tam olarak bu kesişimi işler.

```
        CIA + KABLOSUZ EK EKSENLER — KAPSAMA HARİTASI

   ┌──────────────────────────────────────────────────────────┐
   │  ERİŞİM           GİZLİLİK         BÜTÜNLÜK               │
   │  (kim girer)      (kim okur)       (kim değiştirir)       │
   │     │                │                  │                 │
   │  yetkisiz         dinleme           enjeksiyon            │
   │  bağlanma         MitM              replay                │
   │  parola zaafı     el-sıkışma yak.   içerik değiştirme     │
   │                                                            │
   │  ERİŞİLEBİLİRLİK                    KİMLİK                 │
   │  (hizmet ayakta)                    (karşı taraf gerçek)  │
   │     │                                  │                  │
   │  deauth seli                       evil twin              │
   │  karıştırma                        rogue AP               │
   │  kaynak tüketme                    sahte baz istasyonu    │
   │  (DHCP/ARP/bağlantı tablosu)       MAC spoofing           │
   └──────────────────────────────────────────────────────────┘
        ▲ kullanıcının sorusu burada (Erişim + Erişilebilirlik)
```

Not: CIA üçlüsü bir başlangıç çerçevesidir, evrensel doktrin değil. Bazı modeller "inkâr edilemezlik" (non-repudiation) ve "hesap verebilirlik" (accountability) eksenlerini de ekler. Kablosuz bağlamda bu beş eksen pratikte yeterlidir; ancak kuruma özel risk modellemesinde ek eksenler kaynaktan teyit edilmeli.

---

<a id="3"></a>
## 3. Saldırı Vektörü Taksonomisi: Beş Ana Kol

Şimdi tüm kablosuz saldırı vektörlerini beş ana kola ayıralım. Her kol bir güvenlik özelliğine karşılık gelir; her kolun altında somut vektörler, mekanizmaları, tespit ipuçları ve savunmaları yer alır. Aşağıdaki büyük tablo bu bölümün omurgasıdır; sonraki başlıklar tablonun satırlarını açar.

| Kol | Vektör | Mekanizma (prensip) | Tespit ipucu | Birincil savunma |
|---|---|---|---|---|
| **Erişim** | Yetkisiz WiFi bağlanma | Zayıf/sızdırılmış parola, açık ağ; istemci olarak ağa katılma | Bilinmeyen MAC/cihaz, kota anomalisi | WPA3-SAE, güçlü parola, istemci envanteri, kota izleme |
| **Erişim** | Evil twin ile istemci çekme | Aynı SSID'li sahte AP, istemciyi kendine bağlama | İki özdeş SSID, BSSID/kanal anomalisi | 802.1X/EAP karşılıklı doğrulama, PMF, istemci profili |
| **Gizlilik** | Pasif dinleme | Şifresiz/zayıf şifreli trafiği havadan okuma | (Pasif — doğrudan tespiti zor) | Güçlü uçtan-uca şifreleme (WPA3, TLS, VPN) |
| **Gizlilik** | El-sıkışma yakalama | Kimlik doğrulama mesajlarını kaydedip çevrimdışı analiz | Tekrarlı deauth + yeniden bağlanma örüntüsü | WPA3-SAE (çevrimdışı tahmine dirençli), uzun parola |
| **Gizlilik/Kimlik** | MitM (rogue AP, karma) | Araya girip trafiği aktarma; sahte ağ kimliği | Beklenmeyen gateway/DNS, sertifika uyarısı | VPN, sertifika sabitleme, HSTS, PMF, 802.1X |
| **Bütünlük** | Paket enjeksiyonu | Sahte/değiştirilmiş çerçeve gönderme | Sıra/numara anomalisi, beklenmeyen kaynak | Mesaj bütünlük kontrolü (MIC), şifreleme, PMF |
| **Bütünlük** | Replay (tekrar) | Yakalanmış geçerli mesajı yeniden gönderme | Aynı paketin/nonce'un tekrarı | Nonce/sayaç, zaman penceresi, taze anahtar |
| **Erişilebilirlik** | Deauth/disassoc seli | Yönetim çerçevesiyle istemciyi düşürme | Yoğun deauth, hızlı bağlantı kesilme döngüsü | 802.11w/PMF (yönetim çerçevesi koruması) |
| **Erişilebilirlik** | Ağ katmanı tüketme | DHCP havuzu/ARP/bağlantı tablosu doldurma | Havuz bitmesi, ARP tablo taşması, oturum patlaması | DHCP snooping, port güvenliği, oran sınırlama |
| **Erişilebilirlik** | RF karıştırma | Bandı gürültüyle doldurup spektrumu tüketme | Yüksek gürültü tabanı, SNR çöküşü, geniş kirlilik | Spektrum izleme, yön bulma, lisanslı/korumalı band (yasa dışı saldırı) |
| **Erişilebilirlik** | Enerji/pil tüketme | IoT/BLE cihazını sürekli uyandırıp pilini bitirme | Anormal uyanma sıklığı, pil hızlı tükenme | Bağlantı/oran sınırı, kimlik doğrulamalı uyanma |
| **Erişilebilirlik** | Kota/bulut tüketme | Botnet/proxy ile veri kotası/hizmet kaynağı yakma | Trafik/maliyet ani artışı, coğrafi anomali | Oran sınırlama, anomali tespiti, kota uyarısı |
| **Kimlik** | MAC/SSID spoofing | Başka cihaz/ağ kimliğine bürünme | Aynı MAC iki yerde, OUI tutarsızlığı | Karşılıklı kimlik doğrulama, sertifika, izleme |

Bu tablonun her satırı bir sınıf, her sütun bir savunma boyutudur. Pratikte savunmacı önce hangi sütunun zayıf olduğunu sorar (tespitim var mı? savunmam var mı?), sonra o sütunu güçlendirir. Sıradaki başlık, satırların kullanıcının asıl sorusuna (erişilebilirlik/erişim) düşen kısmını derinleştirir.

---

<a id="4"></a>
## 4. Saldırı Zinciri: Keşif → Erişim → Kalıcılık → Amaç

Taksonominin ikinci ekseni zaman/aşamadır. Gerçek bir olayda saldırı tek bir vektör değil, bir zincir halinde ilerler. Bu zinciri tanımak savunmacıya iki avantaj sağlar: birincisi, zincirin **erken** halkasında kesilen saldırı sonraki halkalara ulaşamaz (savunmayı solda kır); ikincisi, gördüğün bir vektör zincirin neresinde olduğunu söyler ve bir sonraki adımı **öngörmeni** sağlar.

```
   KABLOSUZ SALDIRI ZİNCİRİ (RF/WLAN bağlamı)

   ┌─────────┐   ┌──────────┐   ┌────────────┐   ┌───────────┐
   │  KEŞİF  │──▶│  ERİŞİM  │──▶│  KALICILIK │──▶│   AMAÇ    │
   │ (recon) │   │ (access) │   │(persistence)│  │ (impact)  │
   └─────────┘   └──────────┘   └────────────┘   └───────────┘
        │             │               │                │
   pasif tarama   yetkisiz        sahte AP'yi      gizlilik ihlali
   SSID/istemci   bağlanma        kalıcı kılma     (dinleme),
   keşfi          (parola/evil    rogue cihaz      erişilebilirlik
   probe izleme   twin), MitM     bırakma          (DoS/tüketme),
   kanal/güç      konumlanması                     bütünlük (enjeksiyon)
   haritası
        │                                                │
        │   ◀── savunma "solda" kesmeyi hedefler ──      │
        └────────────────────────────────────────────────┘
            (keşfi tespit et → erişim hiç olmasın)
```

| Aşama | Saldırgan ne yapar (prensip) | Savunmacı ne arar | Erken kesme kontrolü |
|---|---|---|---|
| Keşif | Pasif olarak SSID, istemci, kanal, sinyal gücü haritalar; probe isteklerini dinler | Anormal tarama yoğunluğu, bilinmeyen monitör-mod cihaz imzaları | SSID gizleme sınırlı fayda; asıl olan sonraki aşamayı kapatmak (PMF, WPA3) |
| Erişim | Parola zaafı/evil twin/MitM ile ağa veya istemciye girer | Yeni/bilinmeyen istemci, ikiz SSID, sertifika uyarısı, gateway değişimi | WPA3-SAE, 802.1X, PMF, istemci envanteri |
| Kalıcılık | Rogue AP/cihaz bırakır, kimlik bilgisi saklar, yeniden bağlanmayı kolaylaştırır | Yetkisiz AP/cihaz, yapılandırma değişikliği, beklenmeyen DHCP kiralaması | Düzenli WIDS taraması, port güvenliği, varlık envanteri |
| Amaç | Dinler (gizlilik), keser/tüketir (erişilebilirlik) veya değiştirir (bütünlük) | Trafik/kota anomalisi, hizmet kesintisi, içerik tutarsızlığı | Segmentasyon, oran sınırlama, sıfır-güven |

Mühendislik sezgisi: Çoğu kablosuz olayda **keşif sessizdir ve büyük ölçüde tespit edilemez** (pasif dinleme havada iz bırakmaz). Bu yüzden savunma stratejisi keşfi engellemeye değil, **erişim ve sonrası aşamaları** olabildiğince zorlaştırmaya odaklanmalıdır. "Düşman zaten içeriyi dinliyor olabilir" varsayımıyla tasarlamak (sıfır-güven) bu gerçeğin doğrudan sonucudur.

Çapraz referans: Bu zincir modeli MITRE ATT&CK gibi çerçevelerle uyumludur; MITRE'nin taktik kategorileri (Reconnaissance, Initial Access, Persistence, Impact) bu dört aşamanın daha ince taneli halidir. MITRE rehberine bakınız.

---

<a id="5"></a>
## 5. Kaynak ve Paket Tüketme Saldırıları (Kullanıcının Sorusu)

Bu başlık, kullanıcının asıl merak ettiği soruyu — "başkasının internet paketini/kaynağını kullanan ve tüketen saldırılar" — taksonominin içine oturtarak, her biri için prensip + tespit + savunma üçlüsüyle işler. Hatırlatma: hiçbir alt başlıkta icra reçetesi yoktur; amaç bir savunmacının bu sınıfı tanıyıp ölçebilmesidir.

Kaynak tüketme saldırılarının ortak mantığı tek cümleyle özetlenir: **sonlu bir kaynağı, sahibinin meşru kullanımı dışında harcayarak ya tükenmesine ya da maliyetinin sahibine yıkılmasına yol açmak.** Burada "kaynak" farklı katmanlarda farklı şeydir: RF katmanında spektrum/zaman, bağlantı katmanında bağlantı süreklilik, ağ katmanında IP/adres havuzu ve durum tabloları, uygulama/bulut katmanında veri kotası ve işlem bütçesi, cihaz katmanında ise **enerji (pil)**. Aşağıdaki tablo bu kaynak haritasını verir.

| Katman | Tüketilen kaynak | Tipik saldırı sınıfı | Sahibe yansıyan zarar |
|---|---|---|---|
| RF / Fiziksel | Spektrum, zaman, SNR | Karıştırma (jamming) | Hizmet kesilir; kimse iletişemez |
| Bağlantı / MAC | Bağlantı sürekliliği | Deauth/disassoc seli | İstemciler sürekli düşer |
| Ağ | Adres havuzu, durum tablosu | DHCP tükenmesi, bağlantı tablosu doldurma | Yeni istemci adres/oturum alamaz |
| Uygulama / Bulut | Veri kotası, işlem bütçesi | Yetkisiz erişimle kota yakma, proxy kötüye kullanımı | Fatura artışı, kota dolması, IP itibar kaybı |
| Cihaz | Enerji (pil) | Uyku yoksunluğu (sleep deprivation) | Pil hızla biter, cihaz devre dışı kalır |

Not: Bu saldırıların çoğu "asimetriktir" — saldırganın harcadığı kaynak (birkaç paket), savunmacıya yıkılan zarardan (saatlerce kesinti, dolu havuz, biten pil) çok daha küçüktür. Asimetri, bu sınıfı cazip ve tehlikeli kılan şeydir; savunma da bu asimetriyi tersine çevirmeyi (saldırının maliyetini artırma, etkisini sınırlama) hedefler.

<a id="5-1"></a>
### 5.1 Yetkisiz WiFi Erişimi ile Başkasının Kotasını Kullanma

**Prensip.** Kullanıcının sorusunun çekirdeği budur: bir saldırgan, sahibi olmadığı bir kablosuz ağa bağlanarak o ağın internet bağlantısını, bant genişliğini ve (mobil/kotalı bağlantılarda) **veri kotasını** harcar. Bu mümkün olur çünkü ağa katılma yetkisi, çoğu ev/küçük işletme ağında tek bir paylaşılan sırla (WPA2-PSK parolası) korunur. Üç temel zayıflık bu erişimi açar:

1. **Açık ağlar.** Şifresiz misafir ağları veya yanlış yapılandırılmış erişim noktaları, herhangi bir istemcinin kimlik doğrulamasız katılmasına izin verir. Burada "erişim" zaten serbesttir; kaynak tüketme yalnızca bir bağlanmadan ibarettir.

2. **Parola zayıflığı.** WPA2-PSK'da parola, çevrimdışı tahmine açıktır: el-sıkışma bir kez yakalanırsa (yakalama mekanizması Bölüm 15'te ayrıntılı), zayıf/sözlükte bulunan bir parola hava bağlantısı olmadan, ayrı bir makinede tahmin yoluyla bulunabilir. Parola bulununca saldırgan meşru bir istemci gibi bağlanır ve kotayı harcar. Bu bölümde kırma süreci anlatılmaz; yalnızca **neden zayıf parolanın bu riski açtığı** vurgulanır.

3. **Evil twin ile istemci çekme.** Saldırgan, hedef ağla aynı SSID'yi yayan sahte bir AP kurarsa, istemciler (özellikle "otomatik bağlan" açıkken) sahte ağa geçebilir; bu bir erişim ihlalinden çok kimlik/MitM saldırısıdır (Bölüm 6'da) ama dolaylı olarak meşru kullanıcının trafiğini saldırgan üzerinden akıtarak kaynak/gizlilik etkisi yaratır.

```
   YETKİSİZ ERİŞİMLE KOTA TÜKETME — KAVRAMSAL AKIŞ

   Meşru ağ (kotalı bağlantı)
   ┌──────────────┐        ┌─────────────────┐
   │  Ev/İş AP    │◀──────▶│  İnternet (kota) │
   └──────┬───────┘        └─────────────────┘
          │ paylaşılan sır (PSK) / açık ağ
          │
   ┌──────┴───────┐      ┌──────────────────────┐
   │ Meşru istemci│      │ Yetkisiz istemci      │
   │ (sahip)      │      │ (zayıf parola/açık ağ │
   └──────────────┘      │  üzerinden katılır)   │
                         └──────────────────────┘
                                  │
                         kota/bant genişliği harcanır
                         (sahibe fatura/yavaşlama yansır)
```

**Tespit.** Bu saldırının imzası ağ tarafında görünür: istemci envanterinde **bilinmeyen bir MAC/cihaz**, beklenmeyen saatlerde **veri kullanım sıçraması**, mobil bağlantılarda **kota anomalisi**. Birçok yönlendirici "bağlı cihazlar" listesi sunar; düzenli denetim, tanımadığın bir cihazı erken yakalar. OUI (MAC üreticisi öneki) tutarsızlığı veya rastgele/değişen MAC adresleri de ipucu olabilir.

**Savunma.**

| Kontrol | Ne yapar | Etki |
|---|---|---|
| WPA3-SAE | Parolayı çevrimdışı tahmine dirençli kılar (Dragonfly/SAE) | El-sıkışma yakalansa bile çevrimdışı tahmin pratik dışı |
| Güçlü/uzun parola | Sözlük ve kaba-kuvvet maliyetini astronomik yapar | WPA2'de bile zayıflığı büyük ölçüde kapatır |
| İstemci envanteri/izleme | Bilinmeyen cihazı tespit eder | Yetkisiz erişimi erken yakalama |
| Veri kotası uyarısı | Anormal kullanımı raporlar | Kota yakma saldırısını görünür kılar |
| Misafir ağ yalıtımı | Misafiri iç ağdan ayırır, bant sınırlar | Erişim olsa bile zarar sınırlanır |
| İstemci yalıtımı (client isolation) | İstemcilerin birbirini görmesini engeller | Yanal hareket ve istemci-istemci saldırıyı keser |
| MAC filtreleme (zayıf) | Yalnızca tanınan MAC'lere izin | Tek başına yetersiz (MAC spoofing'le aşılır); ek katman |

Pratikte: Tek en güçlü adım WPA3'e geçiş ve güçlü paroladır; bu ikisi "kotamı başkası harcıyor" senaryosunun en yaygın kökünü kapatır. MAC filtreleme yalnızca tamamlayıcı bir gürültü engelidir, esas savunma sayılmamalıdır (Bölüm 15).

Çapraz referans: Parola/el-sıkışma zayıflığının ayrıntılı anatomisi ve WPA2→WPA3 geçiş mantığı Bölüm 15'tedir; bu bölüm yalnızca kaynak-tüketme açısından özetler.

<a id="5-2"></a>
### 5.2 Deauth/Disassoc Seli ile Bağlantı Kesme

**Prensip.** 802.11'in (eski sürümlerinde) yönetim çerçeveleri — özellikle deauthentication ve disassociation — **kimlik doğrulaması ve şifrelemesi olmadan** gönderilebiliyordu. Bu, tasarımın bir boşluğudur: bir cihaz, başka bir istemci veya AP adına "bağlantıyı sonlandır" çerçevesi yayabilir ve hedef istemci ağdan düşer. Sürekli (sel halinde) gönderildiğinde istemci bir türlü bağlı kalamaz; bu klasik bir erişilebilirlik (DoS) saldırısıdır.

Bu saldırının iki amacı olabilir: (a) saf hizmet kesme — kurban interneti kullanamaz; (b) erişim için ön hazırlık — istemciyi düşürmek, onu yeniden el-sıkışmaya zorlar ve bu el-sıkışma yakalanabilir (Bölüm 5.1 ve Bölüm 15 ile bağlantı). Bu yüzden deauth, hem bir DoS sınıfı hem de bir erişim-öncesi araçtır.

```
   DEAUTH SELİ — KAVRAMSAL (yönetim çerçevesi sahtekarlığı)

   ┌─────────┐    "bağlantıyı kes" (sahte, AP/istemci adına)
   │ Saldırgan│ ─────────────────────────────────────────┐
   └─────────┘                                            ▼
                                              ┌───────────────────┐
   ┌─────────┐  ◀── bağlantı kopar ──         │   Meşru istemci    │
   │   AP    │      (sürekli tekrar)          │ (sürekli düşürülür)│
   └─────────┘                                └───────────────────┘

   PMF (802.11w) yoksa: çerçeve kabul edilir → istemci düşer.
   PMF (802.11w) varsa : yönetim çerçevesi bütünlük korumalı →
                         sahte deauth reddedilir → saldırı etkisiz.
```

**Tespit.** İmza nettir ve WIDS tarafından kolay yakalanır: kısa sürede **yoğun deauth/disassoc çerçevesi**, aynı istemcide **hızlı bağlan-kopar döngüsü**, çoğu zaman çerçevelerin **aynı kaynaktan** gelmesi. Wireshark'ta deauth çerçeveleri (`wlan.fc.type_subtype == 0x0c`) tek bir kanal üzerinde sayıca patlar; Kismet gibi WIDS bunu otomatik alarma çevirir.

**Savunma.** Bu saldırının tek ve kesin savunması **802.11w / Protected Management Frames (PMF)**'dir. PMF, yönetim çerçevelerine bütünlük koruması ekler; sahte (anahtara sahip olmayan) bir deauth çerçevesi reddedilir. WPA3 PMF'i **zorunlu** kılar; WPA2'de PMF isteğe bağlı olarak (WPA2-PMF) etkinleştirilebilir. PMF açıkken klasik deauth seli pratikte etkisizleşir.

| Kontrol | Ne yapar | Sınır |
|---|---|---|
| 802.11w/PMF | Yönetim çerçevesini bütünlük korumasına alır | Saldırının çekirdeğini kapatır |
| WPA3 (PMF zorunlu) | PMF'i varsayılan/zorunlu yapar | En temiz çözüm |
| WIDS alarmı (Kismet) | Deauth selini tespit/raporlar | Kesmez ama görünür kılar, yön bulmaya yarar |
| Kanal çeşitliliği / 6 GHz | Bazı eski saldırı araçlarının kapsamını sınırlar | Kısmi; PMF'in yerini tutmaz |

Pratikte: Alıştırma 12'de kendi ağında PMF'i açıp deauth'un etkisizleştiğini bizzat gözlemleyeceksin — bu, savunmanın somut olarak işlediğini görmenin en öğretici yoludur.

<a id="5-3"></a>
### 5.3 Ağ Katmanı Kaynak Tüketme: DHCP, ARP, Bağlantı Tablosu

RF ve bağlantı katmanından bir adım yukarı çıkınca, kaynak tüketme **ağ katmanının sonlu yapılarını** hedef alır. Buradaki ortak fikir, bir tablonun/havuzun sahte taleplerle doldurulup meşru talebe yer kalmamasıdır. Üç klasik biçim:

**DHCP havuzu tükenmesi (DHCP starvation).** Bir DHCP sunucusu, dağıtabileceği IP adreslerini sonlu bir havuzdan verir. Saldırgan çok sayıda sahte istemci kimliğiyle adres talep ederse havuz tükenir; gerçek bir istemci ağa katılmak istediğinde **adres alamaz** ve bağlanamaz. Bu bir erişilebilirlik saldırısıdır. Bazı varyantlarda saldırgan havuzu tükettikten sonra **sahte bir DHCP sunucusu** (rogue DHCP) devreye sokar ve istemcilere kendi (kötü niyetli) ağ geçidi/DNS bilgisini dağıtır — bu noktada saldırı MitM'e (Bölüm 6) dönüşür.

```
   DHCP TÜKENMESİ — KAVRAMSAL

   Meşru DHCP havuzu: [.10 .11 .12 ... .250]  (sonlu)

   Saldırgan (çok sahte kimlik) ──▶ DHCP sunucu: "bana adres ver" x N
                                          │
   Havuz dolar: [DOLU DOLU DOLU ... DOLU]  (X)
                                          │
   Gerçek istemci ──▶ "bana adres" ──▶ (havuz boş) ──▶ (X) bağlanamaz
```

**Savunma:** DHCP snooping (anahtar üzerinde güvenilir/güvenilmez port ayrımı), port başına MAC/oran sınırlama, port güvenliği. DHCP snooping ayrıca rogue DHCP sunucusunu da engeller (yalnızca güvenilir porttan DHCP yanıtı kabul edilir).

**ARP zehirleme/seli.** ARP, yerel ağda IP↔MAC eşlemesini sağlar ve klasik biçiminde kimlik doğrulamasızdır. Saldırgan sahte ARP yanıtlarıyla ya bir ağ geçidinin MAC'ini kendine yönlendirir (MitM, Bölüm 6) ya da ARP tablosunu sahte girişlerle doldurur (kaynak tüketme/karışıklık). 

**Savunma:** Dinamik ARP İnceleme (Dynamic ARP Inspection — DAI, DHCP snooping'e dayanır), statik ARP girişleri (küçük ölçek), istemci yalıtımı.

**Bağlantı/durum tablosu doldurma.** Yönlendirici/güvenlik duvarı, her aktif bağlantı için bir durum kaydı tutar (NAT/oturum tablosu). Saldırgan çok sayıda yarı-açık veya kısa ömürlü bağlantı üreterek bu tabloyu doldurursa, cihaz yeni meşru bağlantı kuramaz. Bu, kablolu DoS'un klasik bir biçimidir ama kablosuz istemci üzerinden de tetiklenebilir.

**Savunma:** Bağlantı oran sınırlama, oturum zaman aşımı ayarı, SYN cookie benzeri durumsuz koruma, kaynak başına bağlantı sınırı.

| Vektör | Tüketilen yapı | Birincil savunma | İkincil savunma |
|---|---|---|---|
| DHCP tükenmesi | IP adres havuzu | DHCP snooping | Port güvenliği, oran sınırı |
| Rogue DHCP | İstemci güveni | DHCP snooping (güvenilir port) | İstemci yalıtımı |
| ARP zehirleme/sel | ARP tablosu, güven | Dinamik ARP İnceleme (DAI) | Statik ARP, izleme |
| Durum tablosu doldurma | Oturum/NAT tablosu | Oran sınırlama, zaman aşımı | Kaynak başına sınır |

Mühendislik sezgisi: Bu saldırıların ortak savunması **"güvenilir kaynak" kavramını ağa kazandırmaktır.** DHCP snooping ve DAI, aslında "bu bilgi yalnızca güvenilir bir kaynaktan gelebilir" kuralını uygular; oran sınırlama ise "tek bir kaynak sonsuz talep üretemez" kuralını. Sıfır-güven mimarisi bu kuralların genelleştirilmiş halidir.

<a id="5-4"></a>
### 5.4 RF Katmanı: Karıştırma ve Spektrum Tüketme

**Prensip.** En alt katmanda, kaynak doğrudan **spektrumun kendisidir**. Karıştırma (jamming), hedef bandı gürültü veya rakip sinyalle doldurarak meşru iletişimin sinyal-gürültü oranını (SNR) çökertir; alıcı, kendi sinyalini gürültüden ayıramaz hale gelir ve haberleşme durur. Bu, en kaba ama en zor savunulan erişilebilirlik saldırısıdır çünkü fiziksel katmanı hedef alır ve şifreleme/kimlik doğrulama gibi üst katman kontrolleri burada işe yaramaz.

Bu bölüm karıştırmanın **yapılışını anlatmaz**; mekanizma, türleri (geniş bant, dar bant, reaktif, protokol-bilinçli), tespiti (spektrum izleme, gürültü tabanı artışı, yön bulma) ve yasal çerçevesi Bölüm 13'te ayrıntılıdır. Burada yalnızca taksonomideki yerini sabitliyoruz: karıştırma, **RF katmanında spektrum/zaman kaynağını tüketen** bir DoS sınıfıdır.

```
   KARIŞTIRMA — SNR ÇÖKÜŞÜ (kavramsal)

   Normal:    sinyal ▔▔▔▔  ╱gürültü tabanı  →  SNR yüksek, çözülür
              ────────────────────────────────────────

   Karıştırma: sinyal ▔▔▔▔
              gürültü ███████████████████████  →  SNR çöker, çözülemez
              ────────────────────────────────────────
```

**Tespit.** Spektrum analizörü/SDR ile gürültü tabanında ani ve sürekli artış, geniş bantta enerji kirliliği, belirli bir bandın "ölmesi". Yön bulma (Bölüm 9) ile kaynağın yaklaşık konumu kestirilebilir; bu, müdahale için kritiktir.

**Savunma.** Önemli uyarı: karıştırma çoğu ülkede **yasa dışıdır** ve karşı-karıştırma genellikle düzenleyici/yetkili kurumların işidir. Sivil savunma daha çok **tespit, raporlama ve dayanıklılık** üzerinedir: frekans atlamalı/yayılı spektrum sistemler bazı karıştırma türlerine daha dayanıklıdır; kritik bağlar için yedekli/çeşitli kanallar; spektrum izleme ile erken uyarı; ve olayın **yetkili otoriteye bildirilmesi**. Aktif karşı-yayın (jamming'e jamming) hem teknik hem yasal olarak doğru yol değildir.

Çapraz referans: Karıştırma türleri, dayanıklılık teknikleri ve yasal çerçeve Bölüm 13'te; spektrum izleme/yön bulma Bölüm 9'da ayrıntılıdır.

<a id="5-5"></a>
### 5.5 Enerji/Pil Tüketme: Uyku Yoksunluğu (IoT/BLE)

**Prensip.** Pille çalışan IoT ve BLE cihazları, pil ömrünü uzatmak için çoğu zaman **derin uyku** durumunda bekler ve yalnızca gerektiğinde kısa süre uyanır. Uyku yoksunluğu (sleep deprivation) saldırısı bu tasarımı tersine çevirir: cihazı sürekli uyanık tutarak (sahte bağlantı istekleri, sürekli sorgular, geçersiz ama işlenmesi gereken paketler) enerjisini hızla tüketir. Sonuç, pilin normalde aylar dayanacakken günler/saatler içinde bitmesi ve cihazın **erişilemez hale gelmesidir** — bir erişilebilirlik saldırısı, ama hedefi spektrum veya tablo değil, **enerji bütçesi**.

Bu sınıf özellikle düşük güçlü geniş alan ağlarında (LPWAN), BLE çevre birimlerinde ve sensör düğümlerinde önemlidir; çünkü bu cihazlar tasarımları gereği enerji-kısıtlıdır ve bir kez pil bitince fiziksel müdahale (pil değişimi) gerekebilir. Saldırının asimetrisi burada en uçtadır: saldırgan birkaç paketle, sahanın derinindeki bir sensörü kalıcı olarak susturabilir.

```
   UYKU YOKSUNLUĞU — ENERJİ TÜKETME (kavramsal)

   Normal döngü:  [UYKU ────────────] [uyan] [UYKU ───────────]
                  düşük tüketim, pil aylarca dayanır

   Saldırı:       [uyan][uyan][uyan][uyan][uyan][uyan][uyan]...
                  sürekli uyandırma → pil hızla biter → cihaz ölür
```

**Tespit.** Cihaz/ağ tarafında **anormal uyanma sıklığı**, beklenmeyen bağlantı isteği yoğunluğu, telemetri varsa **pil seviyesinde hızlı düşüş**. Bağlantı denemelerinin kaynağındaki tutarsızlık (tanınmayan eşler) ek ipucudur.

**Savunma.**

| Kontrol | Ne yapar |
|---|---|
| Kimlik doğrulamalı uyanma | Cihaz yalnızca doğrulanmış isteklerde tam uyanır; sahte istek hızlı reddedilir |
| Bağlantı/oran sınırı | Birim zamanda kabul edilen istek sayısını sınırlar |
| Hızlı reddetme yolu | Geçersiz paketi minimum enerjiyle elemek (derin işleme öncesi) |
| Bağlantı beyaz listesi | Yalnızca bilinen eşlerle bağlantıya izin |
| Gözlem/telemetri | Pil ve uyanma anomalisini erken raporlama |

Mühendislik sezgisi: Enerji-kısıtlı cihazlarda güvenlik tasarımının altın kuralı, **"pahalı işlemden önce ucuz doğrulama"** koymaktır. Cihaz bir isteği işlemeye (enerji harcamaya) başlamadan önce, o isteğin meşruluğunu mümkün olan en düşük maliyetle elemelidir. Bu ilke uyku yoksunluğu saldırısının asimetrisini büyük ölçüde tersine çevirir.

Çapraz referans: BLE/IoT eşleşme, bağlantı modeli ve cihaz güvenliği Bölüm 16'da; düşük güçlü protokollerin (LoRa, Zigbee) enerji profili de orada ele alınır.

<a id="5-6"></a>
### 5.6 Bulut/Veri Kotası Tüketme: Botnet ve Proxy Kötüye Kullanımı

**Prensip.** En üst katmanda kaynak tüketme, fiziksel spektrumdan tamamen kopar ve **veri kotası, işlem bütçesi ve hizmet kapasitesi** gibi mantıksal/ekonomik kaynakları hedef alır. İki kavramsal biçim öne çıkar:

1. **Botnet kaynaklı tüketme.** Ele geçirilmiş çok sayıda cihaz (sıklıkla zayıf/varsayılan kimlik bilgili IoT cihazları — Bölüm 9 ile bağlantı) bir komuta altında toplanır ve hedefe ya dağıtık hizmet engelleme (DDoS) trafiği yağdırır ya da hedefin kaynaklarını (bant genişliği, sunucu işlem süresi, bulut faturası) tüketir. Burada saldırgan kendi kaynağını değil, **kurban cihaz sahiplerinin** kaynağını kullanır; bu, kullanıcının sorusunun "başkasının kaynağını kullanma" temasının en geniş ölçekli halidir.

2. **Proxy/kota kötüye kullanımı.** Saldırgan, ele geçirdiği veya kötü niyetli yazılım yüklediği cihazları bir **proxy çıkış noktası** olarak kullanır; kurbanın internet bağlantısı ve veri kotası, saldırganın trafiğini taşımak için harcanır. Mobil/kotalı bağlantılarda bu doğrudan fatura/kota zararı doğurur; ayrıca kurbanın IP adresi saldırgan trafiğiyle ilişkilendiği için **itibar kaybı** (kara listeye düşme) yaşanır.

```
   BOTNET / PROXY KÖTÜYE KULLANIMI — KAVRAMSAL

   ┌──────────┐  komuta
   │ Saldırgan │────────┐
   └──────────┘         ▼
              ┌───────────────────────────────┐
              │  Ele geçirilmiş cihazlar       │
              │  (zayıf/varsayılan kimlik)     │
              │   ● ● ● ● ● ● ● ● ● ●          │
              └───────────────┬───────────────┘
                              │ kurbanların kotası/bant genişliği
                              ▼
                  ┌──────────────────────┐
                  │ Hedef / İnternet      │
                  │ (DDoS, kota yakma,    │
                  │  proxy trafiği)       │
                  └──────────────────────┘
```

**Tespit.** Kurban tarafında **trafik ve maliyetin ani, açıklanamayan artışı**, beklenmeyen saatlerde sürekli yüklü bağlantı, coğrafi/uygulama anomalisi (cihaz hiç gitmediği yerlere bağlanıyor), IP itibar uyarıları. Ağ geçidi düzeyinde giden trafik profili çıkarmak (baseline) bu anomaliyi görünür kılar.

**Savunma.**

| Katman | Kontrol | Etki |
|---|---|---|
| Cihaz | Varsayılan parolayı değiştirme, firmware güncelleme | Botnet'e katılımı baştan engeller (Bölüm 9) |
| Ağ | Giden oran sınırlama, anomali tespiti | Kötüye kullanımı sınırlar/görünür kılar |
| Bulut/operatör | Kota uyarısı, DDoS koruması | Maliyet patlamasını erken durdurur |
| Segmentasyon | IoT'yi ayrı ağa alma | Ele geçirilen cihazın erişimini daraltır |
| İzleme | Giden trafik baseline + alarm | Proxy/botnet davranışını yakalar |

Mühendislik sezgisi: Botnet ekonomisinin yakıtı **zayıf ve güncellenmemiş cihazlardır.** Bu yüzden bireysel savunmanın en yüksek getirili adımı, kendi IoT cihazlarının varsayılan kimlik bilgilerini değiştirmek ve firmware'ini güncel tutmaktır (Alıştırma 12'de denetleyeceksin). "Benim küçük kameramdan kime ne zarar gelir?" yanılgısı, tam olarak botnet'lerin büyüdüğü boşluktur.

Not: DDoS'un ağ/uygulama katmanı ayrıntıları ve hacim/protokol/uygulama-katmanı alt türleri bu RF-odaklı serinin kapsamı dışındadır; konunun derini için ağ güvenliği/DDoS kaynakları teyit edilmeli. Burada amaç, RF/IoT katmanının botnet besleme yüzeyini ve savunmasını göstermektir.

---

<a id="6"></a>
## 6. Ortadaki-Adam (MitM) Sınıfı: Rogue AP, Evil Twin, SSL Stripping

**Prensip.** Ortadaki-adam saldırıları, iletişimin iki ucu arasına girerek trafiği aktaran (ve isteğe bağlı okuyan/değiştiren) bir konum elde etmeyi amaçlar. Kablosuz bağlamda bu konum çoğu zaman **sahte bir erişim noktası** kurarak elde edilir. Sınıfın başlıca biçimleri:

| Biçim | Mekanizma | Hedeflediği güven |
|---|---|---|
| Rogue AP | Yetkisiz, kötü niyetli bir AP'yi ağa/ortama yerleştirme | "Bu AP güvenilir altyapının parçası" varsayımı |
| Evil twin | Meşru ağla aynı SSID'yi yayan sahte AP; istemciyi çekme | İstemcinin "tanıdığı SSID'ye otomatik bağlan" davranışı |
| Karma saldırısı | İstemcinin probe ettiği SSID'lere "evet, o benim" yanıtı | İstemcinin tercih edilen ağ listesi (PNL) ifşası |
| SSL stripping | HTTPS'e yükseltmeyi engelleyip düz HTTP'de tutma | Kullanıcının "kilit simgesini kontrol etmeme" alışkanlığı |
| Sahte baz istasyonu | Hücresel istemciyi sahte hücreye çekme (IMSI catcher mantığı) | Telefonun en güçlü hücreye bağlanma davranışı |

Evil twin, kullanıcının sorusuyla doğrudan bağlantılıdır: istemci sahte AP'ye geçtiğinde hem trafiği saldırgan üzerinden akar (gizlilik) hem de meşru kullanıcının bağlantı kaynağı dolaylı olarak saldırganın denetimine girer. Karma, evil twin'in "hedef SSID'yi önceden bilmeye gerek bırakmayan" genellemesidir.

```
   EVİL TWİN / MitM TOPOLOJİSİ (kavramsal)

   Meşru:   İstemci ───── Meşru AP ───── İnternet

   Evil twin / MitM:
            İstemci ──┐
                      ▼
              ┌───────────────┐        ┌──────────┐
              │  Sahte AP     │───────▶│ İnternet │
              │ (aynı SSID,   │  (trafik buradan akar:
              │  saldırgan)   │   okunabilir/değiştirilebilir)
              └───────────────┘
              ▲
              │ savunma: karşılıklı kimlik doğrulama (802.1X/EAP),
              │ sertifika doğrulama, VPN, HSTS, PMF
```

**Tespit.** İki özdeş SSID'nin farklı BSSID/kanalda görünmesi, beklenen AP'nin sinyal/konum profilinin değişmesi, sertifika/uyarı diyalogları, gateway veya DNS değişimi. WIDS (Kismet) yeni/yetkisiz BSSID'leri ve SSID çakışmalarını alarma çevirebilir. Kullanıcı tarafında ani "bu sertifika güvenilmez" uyarısı klasik bir kırmızı bayraktır.

**Savunma.**

| Kontrol | Neyi sağlar |
|---|---|
| 802.1X / EAP (karşılıklı doğrulama) | İstemci AP'yi de doğrular; sahte AP kimliği kanıtlayamaz |
| Sertifika doğrulama/sabitleme | TLS'te sunucunun gerçekliğini garantiler; SSL stripping ve sahte sunucuyu zorlaştırır |
| VPN | Uçtan-uca tünel; aradaki AP içeriği okuyamaz | 
| HSTS | Tarayıcıyı HTTPS'e zorlar; SSL stripping'i engeller |
| 802.11w/PMF | Yönetim çerçevesi manipülasyonunu (deauth ile istemci çekme) zorlaştırır |
| WIDS/WIPS | Rogue/evil twin AP'yi tespit (ve WIPS ile sınırlı müdahale) |

Pratikte: MitM sınıfına karşı en sağlam tek ilke **karşılıklı kimlik doğrulamadır** — istemci yalnızca kendini kanıtlamaz, AP/sunucu da kendini kanıtlar. 802.1X (kurumsal) ve sertifika sabitleme (uygulama) bunu sağlar. Halka açık/açık WiFi'de ise VPN, "ortamı güvenilmez kabul et" yaklaşımının pratik karşılığıdır.

Çapraz referans: Sahte baz istasyonu (IMSI catcher) mekanizması ve hücresel MitM Bölüm 20 (ileri hücresel) ve genel sinyal güvenliği Bölüm 6'da; rogue AP/evil twin'in WiFi'ye özgü ayrıntısı Bölüm 15'tedir.

---

<a id="7"></a>
## 7. Enjeksiyon, Replay ve Spoofing Sınıfı

Bu üç vektör **bütünlük** ve **kimlik** eksenlerinde oturur ve sık birlikte anılır çünkü hepsi "sahte ama geçerli görünen veri" üretmeye dayanır.

**Enjeksiyon (injection).** Saldırgan, iletişime ait olmayan bir çerçeveyi/paketi ortama sokar. Bu çerçeve geçerli formatta olduğu için alıcı tarafından işlenebilir; sonuç, sahte komut, sahte yanıt veya durum bozulmasıdır. Şifreleme ve mesaj bütünlük kontrolü (MIC), enjekte edilen sahte çerçevenin doğru bütünlük etiketi taşıyamayacağı için onu reddeder; bu yüzden güçlü şifreleme enjeksiyonun en doğal savunmasıdır.

**Replay (tekrar).** Saldırgan, daha önce yakaladığı **geçerli** bir mesajı yeniden gönderir. Mesaj gerçekten meşrudur (bütünlük etiketi doğrudur), ama bağlamı sahtedir — örneğin "kapıyı aç" komutunun kaydedilip sonradan tekrar oynatılması. Savunma, her mesajı **bir kez geçerli** kılan mekanizmalardır: nonce (tek kullanımlık sayı), artan sayaç, zaman damgası/pencere, taze oturum anahtarı. Bu mekanizmalar olmadan, şifreleme tek başına replay'i durdurmaz.

**Spoofing (kimlik taklidi).** Saldırgan, başka bir varlığın kimliğine bürünür: MAC spoofing (başka istemci/AP gibi görünme), SSID spoofing (sahte ağ adı), IP/ARP spoofing. Spoofing genellikle başka bir saldırının (MitM, DoS, erişim) yapı taşıdır. Savunma, kimliğin **kanıtlanmasını** zorunlu kılmaktır: kriptografik kimlik doğrulama, sertifika, karşılıklı doğrulama.

| Vektör | Hedef eksen | Çekirdek savunma | Neden şifreleme tek başına yetmez |
|---|---|---|---|
| Enjeksiyon | Bütünlük | MIC + şifreleme | (Genelde yeter; sahte çerçeve MIC'i tutturamaz) |
| Replay | Bütünlük/kimlik | Nonce/sayaç/zaman penceresi | Yakalanan mesaj gerçekten geçerlidir; tazelik gerekir |
| Spoofing | Kimlik | Kriptografik kimlik doğrulama | Kimlik alanı sahte; içerik şifreli olsa da kaynak yalan olabilir |

Mühendislik sezgisi: Bu üçlüyü ayıran soru "mesaj geçerli mi?" değil, **"mesaj bu bağlamda, bu kaynaktan, ilk kez mi geliyor?"** sorusudur. Şifreleme "veri okunamaz" der; bütünlük "değiştirilmemiş" der; tazelik (nonce/sayaç) "tekrar değil" der; kimlik doğrulama "doğru kaynaktan" der. Sağlam bir protokol dördünü birden sağlar.

Çapraz referans: Replay/spoofing'in araç anahtarı, keyless ve RF bağlamındaki somut biçimleri Bölüm 13 (sinyal manipülasyonu) ve Bölüm 16'da (keyless/IoT); genel sinyal manipülasyonu prensibi Bölüm 6'dadır.

---

<a id="8"></a>
## 8. İstemci-Taraflı Saldırılar: Probe İfşası, MAC İzleme, Captive Portal

Saldırı yüzeyi yalnızca AP/ağ tarafı değildir; **istemci cihazın kendisi** de bilgi sızdırır ve hedeflenebilir. Bu sınıf çoğu zaman gizlilik/mahremiyet eksenindedir.

**Probe isteği ifşası.** Birçok istemci, daha önce bağlandığı ağları aramak için aktif probe istekleri yayar ve bu istekler **tercih edilen ağ listesini (PNL)** ifşa edebilir. Bir gözlemci, bir cihazın hangi ağlara bağlandığını (ev, iş, kafe) pasif olarak öğrenip cihazı/kullanıcıyı profilleyebilir; karma saldırısı (Bölüm 6) tam olarak bu ifşayı istismar eder.

**Savunma:** Modern işletim sistemleri pasif tarama ve PNL gizleme yönünde ilerledi; kullanıcı tarafında "otomatik bağlan"ı kapatmak ve eski/kullanılmayan ağ profillerini silmek ifşayı azaltır.

**MAC izleme.** Sabit bir MAC adresi, bir cihazı (ve dolaylı olarak kişiyi) farklı zaman/konumlarda izlemeye olanak tanır (örneğin perakende analitiği veya kötü niyetli takip). 

**Savunma:** MAC rastgeleleştirme (MAC randomization) — modern cihazlar tarama ve bazı bağlantılarda rastgele MAC kullanır; bu, pasif izlemeyi büyük ölçüde zorlaştırır. Kullanıcı, cihazında "özel/rastgele MAC" seçeneğini etkin tutmalıdır.

**Captive portal kötüye kullanımı.** Açık misafir ağlarındaki "giriş için tıklayın" portalı, kullanıcıyı tanıdık bir akışa alıştırır; saldırgan sahte bir captive portal ile kimlik bilgisi veya ödeme bilgisi toplamayı deneyebilir (kimlik avı). 

**Savunma:** Portal üzerinden hassas bilgi girmemek, sertifika/alan adını doğrulamak, mümkünse açık ağda VPN kullanmak; kurumsal tarafta portalın resmî alan adı ve HTTPS ile sunulması.

| Vektör | Sızan/hedeflenen | İstemci savunması | Ağ/kurum savunması |
|---|---|---|---|
| Probe ifşası (PNL) | Bağlandığı ağlar | Otomatik bağlanı kapat, eski profilleri sil | (Doğrudan ağ kontrolü sınırlı) |
| MAC izleme | Cihaz/kişi kimliği | MAC rastgeleleştirme | İzleme amaçlı toplamaktan kaçınma |
| Captive portal kötüye kullanımı | Kimlik/ödeme bilgisi | Hassas bilgi girme, VPN, alan doğrula | Resmî alan + HTTPS, net markalama |
| Çerez/oturum çalma | Oturum jetonu | HTTPS, güvenli çerez, oturum süresi | HSTS, güvenli/HttpOnly çerez |

Mühendislik sezgisi: İstemci-taraflı saldırıların çoğu, cihazın **"yardımcı olma" davranışlarını** (otomatik bağlanma, bilinen ağı arama, portala güvenme) istismar eder. Savunmanın özü, bu otomatik güveni azaltmak ve "ortamı kanıtlanana kadar güvenilmez say" ilkesini istemciye taşımaktır.

---

<a id="9"></a>
## 9. Tedarik Zinciri ve Firmware: Varsayılan Kimlik Bilgileri

**Prensip.** Kablosuz/IoT ekosisteminin en sessiz ama en yaygın zafiyeti, cihazların **varsayılan kimlik bilgileri** ve **güncellenmeyen firmware**'idir. Birçok IoT cihazı (kamera, yönlendirici, akıllı priz, DVR) fabrikadan sabit veya tahmin edilebilir bir yönetici parolasıyla çıkar; kullanıcı bunu değiştirmezse, cihaz internete açıldığı anda kütlesel taramaların kolay hedefi olur. Bu, Bölüm 5.6'daki botnet ekonomisinin ana besleyicisidir: ele geçirilen cihaz, sahibinin kaynağını başkalarına karşı kullanan bir araca dönüşür.

Firmware tarafında iki ayrı risk vardır: (a) **bilinen güvenlik açıkları** — yayımlanmış ama yamalanmamış zafiyetler; (b) **güncelleme mekanizmasının kendisinin zayıflığı** — imzasız/doğrulanmamış firmware güncellemesi, saldırganın cihaza kötü niyetli yazılım yüklemesine kapı aralayabilir. Tedarik zinciri boyutu ise bileşenin/üreticinin güvenilirliğiyle ilgilidir ve değerlendirilmesi çoğu zaman son kullanıcının ötesindedir.

| Risk | Mekanizma | Savunma |
|---|---|---|
| Varsayılan parola | Sabit/tahmin edilebilir yönetici kimliği | İlk kurulumda güçlü, benzersiz parola |
| Yamalanmamış zafiyet | Yayımlanmış açığın kapatılmaması | Düzenli firmware güncelleme, otomatik güncelleme |
| İmzasız firmware | Doğrulanmamış güncelleme yüklenmesi | İmzalı/doğrulanmış güncelleme zorunlu (üretici özelliği) |
| Gereksiz açık servis | İnternete dönük yönetim arayüzü | Servisleri kapatma, yerel ağa hapsetme, segmentasyon |
| Tedarik zinciri | Güvenilmez bileşen/üretici | İtibarlı üretici, güncel destek taahhüdü |

Pratikte: Bireysel kullanıcı için en yüksek getirili üç adım — varsayılan parolayı değiştir, firmware'i güncel tut, IoT'yi ayrı/yalıtılmış bir ağa al. Bu üçü, hem cihazın ele geçirilmesini hem de ele geçirilse bile zararının yayılmasını sınırlar (Alıştırma 12).

Çapraz referans: IoT cihaz sınıfları, protokolleri ve cihaza özgü zafiyet kökenleri Bölüm 16'da; varsayılan-kimlik/firmware temelli ele geçirmenin botnet'e dönüşmesi Bölüm 5.6'da işlenir.

---

<a id="10"></a>
## 10. Araçlar ve Ne Ölçüde: Yetkili Test Bağlamı

Savunmacının bu araçları tanıması şarttır — çünkü bir WIDS'in alarma çevirdiği imzayı anlamak, aracın ne ürettiğini bilmeyi gerektirir; ve kendi ağının dayanıklılığını yalnızca **yetkili, kontrollü** bir testle ölçebilirsin. Aşağıdaki tablo her aracı "ne için, yetkili test bağlamında ne ölçüde" çerçevesinde verir. Bu bir kullanım reçetesi değil, bir **rol haritasıdır**; hiçbir satır "şununla şu saldırıyı yap" talimatı içermez.

| Araç | Birincil rol | Yetkili test bağlamında kullanım | Sınır/uyarı |
|---|---|---|---|
| Wireshark | Paket/çerçeve analizi (pasif) | Kendi ağında trafiği gözlemleme, deauth/anomali tespiti, savunma doğrulama | Pasif; içerik çözme yalnızca kendi anahtarınla |
| Kismet | Kablosuz keşif + WIDS | Kendi ortamında AP/istemci envanteri, rogue/deauth alarmı, anomali görme | İzleme amaçlı; başkalarının trafiğini toplamak yasal değil |
| aircrack-ng paketi | WiFi güvenlik denetimi | Kendi ağında el-sıkışma yakalama, PMF/WPA3 dayanıklılığını doğrulama | Yalnızca sahip olduğun ağ; parola kırma reçetesi bu kitapta yok |
| bettercap | Ağ saldırı/savunma test çatısı | Yetkili pentest'te MitM/ARP senaryosunu **kontrollü** doğrulama | Yüksek güçlü; izinsiz kullanım açıkça suç |
| mdk4 | 802.11 DoS/stres testi | Kendi ağında deauth/beacon stres testiyle PMF'in etkisini ölçme | DoS üretir; yalnızca izole/sahip olunan ortam |
| Scapy | Paket üretimi/ayrıştırma (kod) | Savunma testinde özel çerçeve üretip tespit kuralını sınama | Esnek; sorumluluk tamamen kullanıcıda |
| Spektrum analizörü / SDR | RF gözlem | Gürültü tabanı, karıştırma tespiti, kanal kullanımı izleme | Pasif gözlem; yayın yapmak ayrı yetki/lisans ister |

Kritik çerçeve: Bu araçların hepsinin meşru, savunma-odaklı bir kullanımı vardır (kendi ağını denetleme, bir WIDS kuralını doğrulama, bir savunmanın gerçekten işlediğini görme). Aynı araçların izinsiz bir ağa/cihaza yöneltilmesi TCK 243/244 kapsamına girer. Bu kitap yalnızca **savunma doğrulama ve farkındalık** bağlamını destekler; saldırı icrası için adım listesi vermez.

Not: bettercap ve mdk4 gibi araçlar gerçek DoS/MitM üretebildiği için, yetkili bir testte bile **izole bir laboratuvar ağı** (kendi cihazların, ayrı SSID, başka kullanıcıyı etkilemeyen ortam) tercih edilmelidir. Yan etkilerin başkalarına ulaşmaması savunmacının sorumluluğudur.

---

<a id="11"></a>
## 11. Tespit ve Savunma Mimarisi: WIDS/WIPS, PMF, WPA3, Sıfır-Güven

Tek tek vektörlerin savunmasını ilgili başlıklarda verdik. Şimdi bunları **bütünleşik bir savunma mimarisine** bağlayalım; çünkü gerçek koruma, dağınık kontrollerin değil, katmanlı ve birbirini tamamlayan bir mimarinin sonucudur.

```
   KATMANLI KABLOSUZ SAVUNMA MİMARİSİ

   ┌────────────────────────────────────────────────────────────┐
   │ İZLEME / TESPİT KATMANI                                     │
   │  WIDS/WIPS (Kismet) · spektrum izleme · anomali tespiti    │
   │  → rogue AP, deauth seli, karıştırma, yeni cihaz alarmı     │
   ├────────────────────────────────────────────────────────────┤
   │ KİMLİK / ŞİFRELEME KATMANI                                  │
   │  WPA3-SAE · 802.11w/PMF · 802.1X/EAP · sertifika           │
   │  → erişim, gizlilik, kimlik, yönetim çerçevesi koruması     │
   ├────────────────────────────────────────────────────────────┤
   │ AĞ KONTROL KATMANI                                          │
   │  DHCP snooping · DAI · port güvenliği · oran sınırlama      │
   │  · segmentasyon (IoT/misafir ayrı) · istemci yalıtımı      │
   │  → kaynak tüketme, ARP/DHCP, yanal hareket                  │
   ├────────────────────────────────────────────────────────────┤
   │ SIFIR-GÜVEN / UYGULAMA KATMANI                             │
   │  VPN · TLS/sertifika sabitleme · HSTS · en az ayrıcalık    │
   │  → "ağı güvenilmez say", uçtan-uca güven                    │
   └────────────────────────────────────────────────────────────┘
        her katman bir öncekinin boşluğunu kapatır
```

### WIDS ve WIPS arasındaki fark

| Özellik | WIDS (tespit) | WIPS (önleme) |
|---|---|---|
| İşlev | Anomaliyi algılar ve alarma çevirir | Algılar **ve** otomatik müdahale eder |
| Müdahale | Pasif (raporlar) | Aktif (rogue AP'yi izole etme vb.) |
| Risk | Düşük (yalnızca gözlem) | Müdahale yan etkili olabilir; dikkatli politika gerekir |
| Tipik araç | Kismet (açık kaynak) | Kurumsal WIPS çözümleri |

Önemli: WIPS'in aktif müdahale yetenekleri (örneğin bir cihazı ağdan düşürme) yanlış yapılandırıldığında meşru cihazları etkileyebilir ve hukuki sorumluluk doğurabilir; bu yüzden müdahale politikaları yalnızca **sahip olunan/yönetilen** altyapıyla sınırlı tanımlanmalıdır.

### Savunmanın çekirdek ilkeleri

| İlke | Açıklama | Hangi sınıfı kapatır |
|---|---|---|
| Karşılıklı kimlik doğrulama | İki taraf da birbirini kanıtlar | Evil twin, rogue AP, MitM, spoofing |
| Yönetim çerçevesi koruması (PMF) | Yönetim çerçevelerine bütünlük | Deauth/disassoc DoS |
| Çevrimdışı tahmine direnç (SAE) | Parola yakalansa bile tahmin pahalı | Yetkisiz erişim, el-sıkışma istismarı |
| Tazelik (nonce/sayaç) | Her mesaj bir kez geçerli | Replay |
| Mesaj bütünlüğü (MIC) | Değiştirilmiş çerçeve reddedilir | Enjeksiyon |
| Segmentasyon | Ağı bölmelere ayır | Yanal hareket, kaynak tüketme yayılması |
| Oran sınırlama | Tek kaynak sonsuz talep üretemez | DHCP/ARP/bağlantı tablosu, kota tüketme |
| En az ayrıcalık / sıfır-güven | Hiçbir şeye otomatik güvenme | Genel; tüm zinciri zorlaştırır |
| Gözlem/anomali tespiti | Normalden sapmayı gör | Hepsinin görünürlüğü |

Mühendislik sezgisi: Bu mimarinin gücü **derinlikte (defense in depth)** yatar. Hiçbir tek kontrol her saldırıyı durdurmaz; ama katmanlar üst üste binince, saldırganın her halkayı ayrı ayrı aşması gerekir ve asimetri savunmacı lehine döner. Kullanıcının sorusu açısından: WPA3 (erişim), PMF (deauth DoS), segmentasyon+oran sınırlama (kaynak tüketme) ve izleme (görünürlük) birlikte, "başkası kaynağımı tüketiyor" senaryosunun büyük çoğunluğunu kapatır.

---

<a id="12"></a>
## 12. Alıştırmalar (Yalnızca Kendi Ağında / Yetkili)

Aşağıdaki alıştırmalar tasarım gereği **yalnızca kendi sahip olduğun ağ ve cihazlarda** ya da yazılı yetki verilmiş bir ortamda yapılır; hiçbiri başkasının kaynağına/ağına dokunmayı içermez. Hepsi savunmayı **gözlemleme ve doğrulama** amaçlıdır.

**Alıştırma 1 — PMF ile deauth'un etkisizleştiğini gözle.** Kendi yönlendiricinde önce PMF kapalıyken, sonra PMF (802.11w) etkinken kendi test istemcinin bağlantı kararlılığını gözlemle. Wireshark'ta yönetim çerçevelerini izleyerek, PMF açıkken sahte deauth çerçevelerinin neden artık istemciyi düşürmediğini kendi gözünle doğrula. Amaç: yönetim çerçevesi korumasının somut etkisini görmek.

**Alıştırma 2 — Kendi ağında bir WIDS (Kismet) kur ve anomali gör.** Kendi ortamında Kismet'i çalıştır; kendi AP'lerini, istemcilerini ve kanal kullanımını envantere al. Sonra (yine kendi ağında, kontrollü) bir bağlantı kesme/yeniden bağlanma örüntüsü oluştuğunda Kismet'in bunu nasıl alarma çevirdiğini gözlemle. Amaç: bir WIDS'in hangi imzaları yakaladığını ve normal-anomali ayrımını içselleştirmek.

**Alıştırma 3 — Kendi açık misafir ağında istemci yalıtımını test et.** Kendi yönlendiricinde bir misafir ağı oluştur; istemci yalıtımı (client isolation) kapalı ve açıkken iki test cihazının birbirini görüp göremediğini kontrol et. Amaç: yalıtımın istemci-istemci görünürlüğünü ve dolayısıyla yanal saldırı yüzeyini nasıl kapattığını anlamak.

**Alıştırma 4 — Kendi IoT cihazlarının varsayılan kimlik/firmware durumunu denetle.** Sahip olduğun IoT cihazlarını (kamera, priz, yönlendirici) tek tek geç: yönetici parolası hâlâ varsayılan mı? Firmware güncel mi, üretici yama yayımlamış mı? İnternete dönük gereksiz bir yönetim servisi açık mı? Bulgularını bir envantere yaz ve gerekenleri düzelt. Amaç: botnet besleme yüzeyini kendi cihazlarında kapatmak (Bölüm 5.6 ve 9 ile bağlantı).

**Alıştırma 5 — Kendi ağının kota/trafik baseline'ını çıkar.** Yönlendiricinin bağlı cihaz listesini ve (varsa) trafik istatistiklerini bir hafta boyunca düzenli kontrol et; normal kullanım profilini öğren. Amaç: "anormal" olanı tanıyabilmek için önce "normal"i tanımlamak — kota/kaynak tüketme tespitinin temeli budur.

**Alıştırma 6 — WPA2 ve WPA3 davranış farkını gözle.** Yönlendiricin destekliyorsa, kendi test ağını WPA2-PSK ve WPA3-SAE/transition modlarında ayrı ayrı kur; istemcilerin bağlanma davranışını ve (Wireshark ile) el-sıkışma akışının nasıl farklılaştığını gözlemle. Amaç: SAE'nin çevrimdışı tahmine direncinin neden erişim sınıfını kapattığını kavramsal olarak pekiştirmek.

Not: Bu alıştırmaların hiçbiri başkasının ağına/cihazına yönlendirilmemelidir. mdk4/bettercap gibi aktif araçlar kullanılacaksa, yan etkileri başkalarına ulaşmayacak **izole bir laboratuvar** zorunludur; aksi hâlde yasal sorumluluk doğar.

---

<a id="13"></a>
## 13. Hızlı Referans ve Diğer Bölümler

### Vektör → Hedef → Birincil Savunma (cep kartı)

| Vektör | Hedef eksen | Birincil savunma |
|---|---|---|
| Yetkisiz WiFi erişimi (kota tüketme) | Erişim | WPA3-SAE + güçlü parola + izleme/kota |
| Deauth/disassoc seli | Erişilebilirlik | 802.11w/PMF |
| DHCP tükenmesi / rogue DHCP | Erişilebilirlik | DHCP snooping |
| ARP zehirleme/sel | Erişilebilirlik/bütünlük | Dinamik ARP İnceleme (DAI) |
| Bağlantı tablosu doldurma | Erişilebilirlik | Oran sınırlama, zaman aşımı |
| Karıştırma (jamming) | Erişilebilirlik | Tespit + yön bulma + raporlama (Bölüm 13) |
| Uyku yoksunluğu (pil tüketme) | Erişilebilirlik | Kimlik doğrulamalı uyanma + oran sınırı |
| Botnet/proxy (kota/bulut) | Erişilebilirlik | Cihaz sertleştirme + giden oran sınırı + segmentasyon |
| Evil twin / rogue AP / MitM | Kimlik/gizlilik | Karşılıklı doğrulama (802.1X) + VPN + PMF |
| Enjeksiyon | Bütünlük | MIC + şifreleme |
| Replay | Bütünlük/kimlik | Nonce/sayaç/zaman penceresi |
| Spoofing (MAC/SSID/ARP) | Kimlik | Kriptografik kimlik doğrulama |
| Probe ifşası / MAC izleme | Gizlilik | MAC rastgeleleştirme + otomatik bağlanı kapat |
| Varsayılan kimlik / firmware | Erişim/kalıcılık | Parola değiştir + güncelle + segmentasyon |

### Üç-eksen hatırlatması

- **Hedef ekseni (ne bozulur):** Erişim · Gizlilik · Bütünlük · Erişilebilirlik · Kimlik
- **Zincir ekseni (ne zaman):** Keşif → Erişim → Kalıcılık → Amaç
- **Katman ekseni (nerede):** RF/Fiziksel · Bağlantı/MAC · Ağ · Uygulama/Bulut

### İlgili bölümler

| Konu | Bölüm |
|---|---|
| WiFi el-sıkışma, WPA2→WPA3, evil twin ayrıntısı | Bölüm 15 — WiFi/WLAN Güvenliği |
| BLE/IoT, keyless, LoRa/Zigbee, cihaz güvenliği | Bölüm 16 — Kısa Menzilli Kablosuz ve IoT |
| Karıştırma türleri, drone RF, replay/spoofing (RF) | Bölüm 13 — RF Tehdit ve Karşı-Önlemler |
| Genel sinyal manipülasyonu, açık/savunma çerçevesi | Bölüm 6 — Güvenlik, Açıklar ve Savunma |
| Sahte baz istasyonu, hücresel MitM | Bölüm 20 — İleri Hücresel 4G/5G Güvenlik |
| Yön bulma (karıştırma kaynağı tespiti) | Bölüm 9 — Yer Tespiti ve Yön Bulma |
| Ayıklama/sınıflandırma, trafik analizi | Bölüm 7 — Disiplinler ve Sinyal Ayıklama |
| Paket analizi pratiği | Wireshark rehberi |
| Saldırı zinciri taktik çerçevesi | MITRE ATT&CK rehberi |
| Botnet/zararlı yazılım ekosistemi | Malware rehberi |

### Anahtar çıkarımlar

1. **Saldırı sayısı sonsuz, sınıf sayısı azdır.** Taksonomi (üç eksen: hedef, zincir, katman) yeni varyantları bile doğru rafa koymanı sağlar.
2. **Kaynak tüketme asimetriktir.** Saldırganın küçük maliyeti, savunmacıya büyük zarar yıkar; savunma bu asimetriyi tersine çevirmeyi (maliyet artırma, etki sınırlama) hedefler.
3. **"Başkasının kaynağını tüketme" iki kolun kesişimidir:** yetkisiz erişim (kota/bant) + erişilebilirlik (DoS). WPA3 + PMF + segmentasyon + izleme birlikte bunun çoğunu kapatır.
4. **Yönetim çerçevesi koruması (PMF) ve SAE (WPA3),** kablosuzun iki en yüksek getirili savunma anahtarıdır; biri DoS sınıfını, diğeri erişim sınıfını büyük ölçüde kapatır.
5. **Botnet'in yakıtı zayıf IoT cihazlarıdır;** kendi cihazını sertleştirmek (varsayılan parola + firmware + segmentasyon) hem kendini hem başkalarını korur.
6. **Karşılıklı kimlik doğrulama ve sıfır-güven,** MitM/spoofing/evil twin sınıfının ortak panzehiridir; "ağı kanıtlanana kadar güvenilmez say".
7. **Savunma katmanlıdır (defense in depth);** hiçbir tek kontrol yeterli değildir, ama katmanlar asimetriyi savunmacı lehine çevirir.
8. **Yasal sınır mutlaktır:** anlatılanların tamamı kendi ağında/yetkili pentest'te geçerlidir; izinsiz erişim ve kaynak tüketme TCK 243/244 kapsamında suçtur.

> Kapanış notu: Bu bölüm bir saldırı el kitabı değil, bir **sınıflandırma ve savunma haritasıdır**. Bir anomaliyle karşılaştığında refleksin "bu hangi eksende, hangi sınıfta, hangi katmanda?" olmalı; çünkü doğru sınıflandırma, doğru savunmanın yarısıdır. Saldırı tekniklerini öğrenmenin tek meşru amacı, onları daha iyi savunmaktır — ve bu savunma her zaman kendi sistemlerinde ve yasanın çizdiği sınır içinde kalır. Belirli yapılandırmaların güncel davranışı ve mevzuatın son hali için resmî kaynaklar teyit edilmelidir.
