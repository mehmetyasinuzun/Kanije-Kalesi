# 📡 MISP — TEHDİT İSTİHBARATI (CTI) TAM USTALIK REHBERİ
## IOC'den Korelasyona, Galaxy'den PyMISP Otomasyonuna, Püf Noktalarıyla Uçtan Uca

> **Amaç:** MISP'i "kur, bir IOC gir, kapat" seviyesinden çıkarıp, bir **siber tehdit istihbaratı (CTI)** analisti gibi **yapılandırılmış istihbarat üretmeyi, korele etmeyi, paylaşmayı ve otomatikleştirmeyi** öğretmek. Bu rehber yalnızca *hangi alanı doldur*'u değil, **neden o veri modeli**, **bir IOC ne zaman değerli ne zaman gürültü**, **neyi paylaşırsın neyi ASLA paylaşmazsın** sorularını da yanıtlar. Forum cevaplarında bulamayacağın TLP disiplini, korelasyon gürültüsü yönetimi, warning-list ile false-positive bastırma, feed zehirlenmesi ve attribution riski — burada. Hiç bilmeyen biri bu rehberin sonunda kendi MISP olayını (event) sıfırdan kurabilir, PyMISP ile otomatikleştirebilir ve bir CTI raporunu MISP'e bağlayabilir.

> ⚠️ **Önce bunu oku — TLP yanlış paylaşımı bir felakettir:** MISP teknik bir araçtır ama içine koyduğun şey **ham istihbarattır** — yanlış kişiyle, yanlış sınıflandırmayla paylaşılırsa bir operasyonu yakar, bir mağduru ifşa eder, bir atıf hatasını yayar. Bir IOC'yi yanlış **TLP** ile işaretlemek, hiç paylaşmamaktan **daha tehlikelidir**, çünkü güvenilir görünür ve zincirleme yayılır. **Bölüm 9 (TLP) ve Bölüm 16 (OPSEC) atlanmaz.**

> ⚠️ **İkinci uyarı — yanlış güvenlik hissi:** "MISP'e 10.000 IOC besledim, artık korunuyorum" demek, *hiç beslememekten* daha tehlikelidir. Bunların çoğu **tazeliğini yitirmiş (expired)**, bir kısmı **false-positive** (Google'ın IP'si!), bir kısmı **zehirli feed**'den gelmiştir. Yapılandırılmamış, warning-list ile temizlenmemiş, sighting ile doğrulanmamış bir IOC denizi SIEM'ini alarm yağmuruna boğar ve **gerçek tehdidi gürültüde gizler**. Bölüm 12 (Warning Lists) ve Bölüm 8 (Korelasyon Gürültüsü) atlanmaz.

---

## 📑 İÇİNDEKİLER

1. [MISP Nedir, Neden? CTI Ekosistemindeki Yeri](#1)
2. [Tehdit İstihbaratı Paylaşımının Mantığı — Tek Başına vs Topluluk](#2)
3. [Kurulum — VM / Docker / Üretim & İlk Yapılandırma](#3)
4. [🔥 VERİ MODELİ (Zero-to-Hero) — Event, Attribute, Object, Tag, Galaxy, Taxonomy, Sighting](#4)
5. [🔥 Bir Event'i Sıfırdan Oluşturma — Phishing Kampanyası (Adım Adım)](#5)
6. [🔥 GALAXY — ATT&CK, Threat Actor, Ransomware Cluster'ları Bağlama](#6)
7. [TAXONOMY & Machine-Tag Mantığı](#7)
8. [🔥 KORELASYON — Otomatik IOC İlişkilendirme & Gürültü Yönetimi](#8)
9. [🔥 TLP — Traffic Light Protocol Disiplini](#9)
10. [FEEDS — CIRCL OSINT, Açık Kaynak Feed'ler & Güvenilirlik](#10)
11. [IMPORT / EXPORT — STIX, OpenIOC, CSV & SIEM/IDS Beslemesi](#11)
12. [🔥 WARNING LISTS & False-Positive Bastırma](#12)
13. [SHARING — Sharing Groups, Organizasyonlar & Sync Server](#13)
14. [🔥 PyMISP / API — Otomasyon (Kod Örnekleriyle)](#14)
15. [SIGHTINGS — "Ben de Gördüm" & IOC Tazeliği/Yaygınlığı](#15)
16. [🔥 PÜF NOKTALARI — Piyasada Bulamayacakların](#16)
17. [🔥 CTI RAPORU İÇİN MISP — Rapordan Event'e, Event'ten Rapora](#17)
18. [Yaygın Ölümcül Hatalar](#18)
19. [🏰 Kanije Kalesi ile Birlikte Kullanım — Forensik Çıktıyı IOC'ye Çevirme](#19)
20. [Hızlı Referans & Operasyonel Kontrol Listesi](#20)
21. [⚖️ Hukuki & Etik Sınır + Çapraz Referanslar](#21)

---

<a id="1"></a>
## 1. 🧭 MISP Nedir, Neden? CTI Ekosistemindeki Yeri

**MISP** (Malware Information Sharing Platform & Threat Sharing), tehdit istihbaratını **yapılandırılmış, makine-okunur ve paylaşılabilir** biçimde toplamak, korele etmek, depolamak ve dağıtmak için kullanılan **açık kaynaklı bir platformdur** (GPL lisanslı; başlangıçta Belçika Ordusu CERT'i için geliştirildi, bugün **CIRCL — Computer Incident Response Center Luxembourg** liderliğinde geniş bir toplulukça sürdürülür).

Tek cümleyle: MISP, bir analistin kafasındaki ya da bir Excel dosyasındaki dağınık "şu IP kötü, şu hash zararlı, şu aktör bu tekniği kullanır" bilgisini **standart bir veri modeline** oturtur, otomatik olarak **diğer olaylarla ilişkilendirir (korelasyon)** ve güvendiğin kurumlarla **kontrollü biçimde paylaşır**.

### MISP üç temel problemi çözer

| Problem | MISP'siz hâli | MISP ile |
|---|---|---|
| **Yapı (structure)** | IOC'ler e-postada, Excel'de, PDF'te dağınık; her analist kendi formatında | Event/Attribute/Object şeması — makine-okunur, tutarlı |
| **Korelasyon** | "Bu IP'yi daha önce görmüş müydük?" → kimse bilmiyor | Otomatik: aynı IOC tüm olaylarda anında ilişkilendirilir |
| **Paylaşım** | Güven varsa e-postayla, yoksa hiç; TLP elle takip | Sharing group + TLP + sync — kontrollü, ölçeklenebilir, denetlenebilir |

### CTI istihbarat döngüsünde MISP'in yeri

Klasik istihbarat döngüsü (intelligence cycle) altı aşamadır; MISP esas olarak **işleme, analiz ve yayma** aşamalarının omurgasıdır:

```
   ┌─────────────┐
   │ 1. Yönlendirme│  (Direction: neyi istihbar edeceğiz? — örn. bizi hedefleyen ransomware)
   └──────┬──────┘
          ▼
   ┌─────────────┐
   │ 2. Toplama   │  (Collection: feed, OSINT, IR çıktısı, sandbox, honeypot)
   └──────┬──────┘
          ▼
   ┌─────────────┐    ┌──────────────────────────────────────┐
   │ 3. İşleme    │◄──►│  M I S P                              │
   │ 4. Analiz    │    │  Event ▸ Attribute/Object ▸ Tag/Galaxy │
   │ 5. Yayma     │    │  Korelasyon ▸ Sighting ▸ Sharing       │
   └──────┬──────┘    └──────────────────────────────────────┘
          ▼
   ┌─────────────┐
   │ 6. Geri besl.│  (Feedback: sighting, false-positive, IOC güncelleme → döngü tekrar)
   └─────────────┘
```

> 🧠 **CTI ilkesi — MISP bir "veritabanı" değil, bir "ilişki motoru"dur.** Asıl değeri tek tek IOC'lerde değil, **IOC'ler arasındaki ve olaylar arasındaki bağlarda** yatar. "Bu hash" değil, "bu hash → bu C2 domain'i → bu aktör → bu kampanya → bizi geçen ay vuran olay" zinciri istihbarattır. Tek başına bir IOC veridir; bağlamlandırılmış IOC istihbarattır.

### Kimler kullanır?

| Aktör | Kullanım amacı |
|---|---|
| **CERT / CSIRT** (ulusal/kurumsal) | Olay müdahalesinde IOC toplama, üye kurumlara dağıtma (CIRCL, US-CERT, ulusal CERT'ler) |
| **SOC** (Security Operations Center) | SIEM/IDS/EDR'a IOC besleme, alarm zenginleştirme (enrichment), avlama |
| **ISAC / ISAO** (sektörel paylaşım toplulukları) | Aynı sektördeki kurumlar arası tehdit paylaşımı (finans FS-ISAC, sağlık, enerji) |
| **CTI ekipleri / tehdit avcıları** | Aktör/kampanya takibi, TTP eşleme, rapor üretimi |
| **MSSP / tehdit istihbaratı satıcıları** | Müşterilere kürate edilmiş feed sağlama |

### MISP vs alternatifler (dürüst kıyas)

| Araç | Güçlü yanı | Ne zaman |
|---|---|---|
| **MISP** | Açık kaynak, zengin veri modeli, korelasyon, galaxy/taxonomy, geniş topluluk feed'i | İstihbarat **paylaşımı + korelasyon** odaklıysa; ücretsiz, kendi altyapın |
| **OpenCTI** | Modern grafik-tabanlı (knowledge graph), STIX 2.1 yerli, görsel ilişki | Bilgi grafiği/aktör-merkezli derin modelleme; MISP ile **birlikte** (konnektörle) |
| **ThreatConnect / Anomali / Recorded Future** | Ticari, hazır kürate feed, scoring, entegrasyon | Bütçe varsa, "anahtar teslim" istihbarat isteniyorsa |
| **TheHive + Cortex** | Olay yönetimi (case management) + analiz orkestrasyonu | IR/olay yönetimi odaklı; MISP ile entegre (TheHive↔MISP) |
| **STIX/TAXII (ham)** | Standart format + taşıma protokolü (araç değil, dil/protokol) | Araçlar arası **birlikte çalışabilirlik**; MISP bunları konuşur |

> 💡 **Katmanlı yaklaşım:** Pratikte modern bir CTI yığını şudur: **MISP** (IOC ambarı + korelasyon + paylaşım) ↔ **OpenCTI** (knowledge graph + STIX modelleme) ↔ **TheHive** (olay/case yönetimi) ↔ **SIEM/EDR** (tüketim). MISP tek başına bir silah değil; **istihbaratın kalbi** — beslenir, korele eder, dağıtır. ATT&CK ile **derinlemesine TTP modellemesi** için ayrı dosyaya bak (Bölüm 21 çapraz referans).

---

<a id="2"></a>
## 2. 🤝 Tehdit İstihbaratı Paylaşımının Mantığı — Tek Başına vs Topluluk

CTI'nin temel önermesi şudur: **bir saldırgan aynı altyapıyı, aynı araçları, aynı TTP'leri birden çok hedefe karşı tekrar kullanır.** Dolayısıyla **bir kurbanın bugün gördüğü IOC, başka bir kurbanın yarınki erken uyarısıdır.**

```
   TEK BAŞINA SAVUNMA                    TOPLULUK SAVUNMASI (MISP)
   ─────────────────                    ─────────────────────────
   Kurum A vuruldu                      Kurum A vuruldu
        │                                    │
        ▼                                    ▼  IOC'leri MISP'e koyar (TLP:AMBER)
   IOC'leri kendine saklar               ┌──────────────┐
        │                                │  MISP / ISAC │
        ▼                                └──┬───┬───┬───┘
   Kurum B aynı aktörce                     ▼   ▼   ▼
   sıfırdan vurulur ──► kayıp           B   C   D kurumları
   (A'nın acısından                     IOC'leri ÖNCEDEN bloklar
    kimse öğrenmedi)                     ──► saldırı ENGELLENDİ
```

### Neden paylaşmak rasyonel?

- **Asimetri:** Saldırgan bir kez geliştirir, çok kez kullanır. Savunma da bir kez tespit eder, çok kez paylaşırsa asimetri savunma lehine döner.
- **Erken uyarı:** "Patient zero" (ilk kurban) acı çeker; geri kalan topluluk **onun bedeliyle** korunur. Sıra sana geldiğinde sen de korunmuş olursun.
- **Doğrulama (sighting):** Bir IOC'yi 5 farklı kurum bağımsız gördüyse, o IOC'nin gerçekliği ve aktifliği güçlenir.
- **Bağlam zenginleşmesi:** Sen sadece bir hash görmüş olabilirsin; başkası o hash'in C2'sini, aktörünü, kampanyasını eklemiştir.

> 🧠 **Ama paylaşım pasif değil disiplinli olmalı:** "Her şeyi herkesle paylaş" felaketle sonuçlanır (Bölüm 9, 16). Paylaşım her zaman **TLP sınıflandırması**, **OPSEC değerlendirmesi** ve **atıf riski** süzgecinden geçer. "Ne paylaşılır" kadar "ne ASLA paylaşılmaz" da CTI ustalığının çekirdeğidir.

---

<a id="3"></a>
## 3. 📥 Kurulum — VM / Docker / Üretim & İlk Yapılandırma

MISP üç ana yolla kurulur. Öğrenme için VM, hızlı deneme için Docker, ciddi kullanım için üretim kurulumu.

### Seçenek 1 — Resmi VM (öğrenmeye en hızlı başlangıç)

CIRCL, hazır yapılandırılmış bir **VirtualBox/VMware OVA** yayınlar (`misp-vm`). İndir, içe aktar, başlat — dakikalar içinde çalışan bir MISP.

- **İndirme:** Yalnızca resmi kaynak — `https://www.misp-project.org/` → Download → Training VM (OVA).
- **Varsayılan kimlik:** Web arayüzü `admin@admin.test` / `admin` (parolayı **ilk açılışta değiştir** — değiştirmeden hiçbir şey yapma). SSH ve MySQL parolaları belge sayfasında listelidir; hepsini değiştir.
- **Ağ:** Öğrenme için NAT/host-only yeterli. **İnternete açık bırakma** (test VM'i sertleştirilmemiştir).

> ⚠️ **Resmi VM bir öğrenme/eğitim aracıdır, üretim değildir.** Varsayılan parolalar herkesçe bilinir; içine gerçek/hassas istihbarat koyma.

### Seçenek 2 — Docker (`misp-docker`)

Resmi topluluk Docker derlemesi (`MISP/misp-docker` deposu) hızlı, tekrarlanabilir kurulum sağlar.

```bash
git clone https://github.com/MISP/misp-docker.git
cd misp-docker
cp template.env .env
# .env içinde EN AZ şunları düzenle:
#   BASE_URL        → https://misp.kurumun.local  (sertifika/erişim için doğru olmalı)
#   ADMIN_EMAIL     → ilk admin kullanıcı
#   ADMIN_PASSWORD  → güçlü parola (varsayılanı ASLA bırakma)
#   MYSQL_PASSWORD / MYSQL_ROOT_PASSWORD → güçlü
docker compose pull
docker compose up -d
# İlk başlatmada şema kurulur (birkaç dakika); logları izle:
docker compose logs -f misp-core
```

> 🔑 **Püf:** `BASE_URL` yanlışsa arayüz çalışır ama **link/STIX export/sync URL'leri bozulur** ve PyMISP bağlanamaz. Kurulumda ilk doğrulayacağın şey budur.

### Seçenek 3 — Üretim kurulumu (kısa)

Ciddi kullanım için bare-metal/VM üzerine resmi **`INSTALL.sh`** (Ubuntu LTS hedefli) ya da olgun bir Docker derlemesi. Üretim, kurulumdan çok **işletim disiplinidir** (aşağıdaki güvenlik notları ve Bölüm 16).

### İlk yapılandırma — sırasıyla yapılacaklar

1. **Admin parolasını değiştir** + (mümkünse) **çok faktörlü doğrulama** etkinleştir.
2. **`MISP.baseurl`** ve **`MISP.external_baseurl`** ayarlarını doğru gir (Administration → Server Settings → MISP).
3. **Kendi organizasyonunu** tanımla (host org). Administration → Organisations. UUID'sini not et — paylaşımda kimliğindir.
4. **Diagnostics**'i çalıştır (Administration → Server Settings → Diagnostics) — kırmızıları temizle (worker, GnuPG, redis, dosya izinleri).
5. **Background worker'lar** çalışıyor mu doğrula (correlation, default, email, cache jobs). Korelasyon ve feed çekme bunlara bağlıdır — durmuşlarsa MISP "sessizce" işlemez.
6. **GnuPG/e-posta** yapılandır (bildirim ve şifreli e-posta için; opsiyonel ama önerilir).
7. **TLS** zorunlu kıl — MISP arayüzü ve API **asla** düz HTTP'de çalışmamalı.

### Kullanıcı & rol modeli

MISP **rol-tabanlı** erişim kullanır. Temel roller:

| Rol | Yetki |
|---|---|
| **admin** (site admin) | Tüm sistem: organizasyonlar, sunucular, ayarlar, tüm veri |
| **org admin** | Yalnızca kendi organizasyonunun kullanıcıları/verisi |
| **user** | Olay oluşturma/düzenleme (kendi org'unda), yayımlama izni role bağlı |
| **publisher** | Olayları **yayımlama** (publish) yetkisi — paylaşımı tetikler |
| **read only / sync user** | Yalnızca okuma; sync user yalnızca senkronizasyon API'si için |

> 🔑 **Püf — yayımlama (publish) ayrı bir yetkidir.** Bir olayı *oluşturmak* onu paylaşmaz. **Publish** edilene kadar olay sadece senin org'unda taslaktır. Yayımlama, sync ve e-posta bildirimini tetikler → **publish düğmesi, paylaşım kararının verildiği andır.** Yanlış TLP'li bir olayı publish etmek geri alınması zor bir hatadır.

---

<a id="4"></a>
## 4. 🔥 VERİ MODELİ (Zero-to-Hero)

Bu bölüm rehberin kalbidir. MISP'i anlamak = veri modelini anlamaktır. Hiyerarşiyi önce görsel olarak oturt:

```
ORGANISATION (senin kurumun / kaynak kurum)
   │
   └── EVENT  (bir "olay" = bir kampanya/olay/rapor kabı; tüm bağlamı taşır)
        │   ├─ info: "APT-X Bordro Phishing Kampanyası — Mayıs 2026"
        │   ├─ date, threat_level, analysis, distribution (paylaşım), uuid
        │   │
        │   ├── ATTRIBUTE  (tekil IOC: bir veri parçası)
        │   │      ├─ category: "Network activity"
        │   │      ├─ type: "ip-dst"
        │   │      ├─ value: "203.0.113.45"
        │   │      ├─ to_ids: true  (IDS/SIEM'e ihraç edilsin mi?)
        │   │      └─ comment, distribution, tags...
        │   │
        │   ├── OBJECT  (yapılandırılmış IOC grubu = ilişkili attribute demeti)
        │   │      └─ "file" objesi: { filename, md5, sha1, sha256, size, entropy }
        │   │         └─ ilişki (relationship): "downloads" → başka obje/attribute
        │   │
        │   ├── TAG  (etiket — taxonomy/galaxy'den ya da serbest)
        │   │      ├─ "tlp:amber"          (taxonomy machine-tag)
        │   │      └─ "misp-galaxy:mitre-attack-pattern=..." (galaxy cluster)
        │   │
        │   └── (event seviyesinde de) GALAXY cluster'ları, SIGHTING'ler
        │
        └── KORELASYON (otomatik): bu event'in attribute'ları,
            BAŞKA event'lerdeki aynı değerlerle otomatik bağlanır
```

### 4.1 Event (Olay) — en üst kap

Bir **Event**, tek bir tutarlı bağlamı temsil eder: bir saldırı kampanyası, bir IR vakası, bir tehdit raporu, bir malware ailesi analizi. **IOC'ler tek başına durmaz; her zaman bir event'e ait olurlar.** Event'in kritik alanları:

| Alan | Ne işe yarar |
|---|---|
| **info** | İnsan-okunur başlık/özet. Net, aranabilir olmalı ("APT28 SMTP C2 — 2026-05" gibi) |
| **date** | Olayın tarihi (saldırı/gözlem tarihi, oluşturma değil) |
| **threat_level** | High / Medium / Low / Undefined |
| **analysis** | Initial / Ongoing / Completed (analiz olgunluğu) |
| **distribution** | Paylaşım kapsamı (aşağıda) — **en kritik alan** |
| **uuid** | Evrensel benzersiz kimlik — senkronizasyon ve atıfta bu kullanılır |

**Distribution (dağıtım) seviyeleri — paylaşımın can damarı:**

| Seviye | Anlamı |
|---|---|
| **0 — Your organisation only** | Sadece senin org'un görür. Hassas/taslak veri burada başlar |
| **1 — This community only** | Senin MISP örneğindeki tüm org'lar |
| **2 — Connected communities** | Senin + doğrudan bağlı (sync) topluluklar |
| **3 — All communities** | Sınırsız yayılır (tüm sync zinciri) — **en dikkatli kullanılan** |
| **4 — Sharing group** | Yalnızca seçtiğin **belirli** org/topluluk kümesi (en granüler) |

> 🧠 **Distribution, attribute seviyesinde event'i geçemez ama daha da kısıtlayabilir.** Yani event "community only" ise içindeki bir attribute "your org only" olabilir (daha kapalı), ama "all communities" olamaz (daha açık olamaz). Hassas tek bir IOC'yi olayın geri kalanından daha kapalı tutmak için bu kullanılır.

### 4.2 Attribute (IOC) — atomik istihbarat birimi

Bir **Attribute**, tek bir veri parçasıdır ve iki eksende tanımlanır: **category** (bağlam — "bu IOC ne işe yarıyor?") ve **type** (veri türü — "bu IOC neyin verisi?"). Aynı `type` farklı `category`'lerde olabilir (örn. bir IP "Network activity" da olabilir "Payload delivery" de).

**En sık kullanılan attribute türleri (type) ve kategorileri:**

| type | Açıklama | Tipik category | to_ids varsayılan |
|---|---|---|---|
| `ip-src` | Kaynak IP (saldırgan IP'si) | Network activity | evet |
| `ip-dst` | Hedef IP (C2/exfil sunucusu) | Network activity | evet |
| `hostname` | Sunucu adı | Network activity | evet |
| `domain` | Alan adı (C2/phishing) | Network activity | evet |
| `url` | Tam URL (indirme/phishing linki) | Payload delivery / Network | evet |
| `uri` | URI yolu | Network activity | evet |
| `md5` | Dosya MD5 hash'i | Payload delivery | evet |
| `sha1` | Dosya SHA-1 hash'i | Payload delivery | evet |
| `sha256` | Dosya SHA-256 hash'i (tercih edilen) | Payload delivery | evet |
| `filename` | Dosya adı | Payload delivery / Artifacts | hayır (genelde) |
| `filename\|sha256` | Ad + hash birleşik | Payload delivery | evet |
| `email-src` | Gönderen e-posta (phishing) | Payload delivery | evet |
| `email-subject` | E-posta konusu | Payload delivery | hayır |
| `email-dst` | Alıcı (genelde hedef — **PII riski!**) | Payload delivery | hayır |
| `ssdeep` | Fuzzy hash (varyant gruplama) | Payload delivery | evet |
| `imphash` | Import hash (PE benzerliği) | Payload delivery | evet |
| `ja3-fingerprint-md5` | TLS istemci parmak izi | Network activity | evet |
| `btc` | Bitcoin adresi (fidye) | Financial fraud | evet |
| `mutex` | Malware mutex adı | Artifacts dropped | evet |
| `regkey` | Registry anahtarı (persistence) | Persistence mechanism | evet |
| `vulnerability` | CVE numarası | External analysis | hayır |
| `yara` | YARA kuralı | Payload installation | hayır |
| `pattern-in-traffic` | Trafikteki desen | Network activity | evet |
| `windows-service-name` | Servis adı (persistence) | Persistence mechanism | hayır |

> 🔥 **`to_ids` bayrağı — en yanlış anlaşılan alan.** Bu bayrak "bu attribute makinesel **tespite (IDS/SIEM/EDR)** ihraç edilmeli mi?" demektir. `to_ids=true` → bu IOC, IDS/Suricata/Sigma export'una girer ve gerçek alarm üretir. **Bağlamsal ama tespit edilebilir olmayan** veriler (örn. bir e-posta konusu, bir genel filename, bir CVE) `to_ids=false` bırakılır — yoksa SIEM'i çöple doldurur. **Kural: yüksek-güvenilir, ayırt edici, bloklanabilir IOC → to_ids=true; bağlam/gürültü → to_ids=false.** Bu disiplin, MISP'in SIEM beslemesini kullanışlı kılan tek şeydir.

### 4.3 Object (Nesne) — yapılandırılmış IOC grupları

Tek tek attribute'lar bir dosyayı tam tanımlamak için dağınık kalır. **Object**, ilişkili attribute'ları **şablonlu (template)** bir demet hâlinde gruplar. Örnek: bir `file` objesi içinde `filename`, `md5`, `sha1`, `sha256`, `size-in-bytes`, `entropy`, `mimetype` birlikte durur — hepsi **aynı dosyaya** ait olduğu açıktır.

Yaygın object şablonları: `file`, `pe` (+`pe-section`), `url`, `domain-ip`, `email`, `network-connection`, `registry-key`, `x509` (sertifika), `vulnerability`, `ja3`, `whois`, `credential`.

**Object'ler birbirine ve attribute'lara `relationship` ile bağlanır:**

```
[email objesi] ──"contains"──► [url objesi] ──"downloads"──► [file objesi]
   from: a@evil                  http://evil/x.doc            sha256: ab12...
   subject: "Bordro"                                          (downloads)
                                                                   │
                                                            "drops"│
                                                                   ▼
                                                          [file objesi: payload.exe]
                                                          sha256: cd34...
```

> 🧠 **Attribute mı Object mi?** Tek, bağımsız bir IOC (yalın bir C2 IP'si) → **attribute** yeter. Çok yönlü, birbirine bağlı bir varlık (bir dosya: ad+hash'ler+boyut+imza) → **object** kullan. Object'ler korelasyonu da zenginleştirir: bir dosyanın SHA-256'sı başka olayda görülürse, **tüm dosya bağlamı** ilişkilenir.

### 4.4 Tag, Galaxy, Taxonomy — bağlam katmanı

- **Tag (etiket):** Bir event ya da attribute'a yapıştırılan işaret. Serbest metin olabilir ama **gerçek güç machine-tag'lerdedir** (taxonomy/galaxy'den gelen yapılandırılmış etiketler).
- **Taxonomy:** Standart, üzerinde anlaşılmış **etiket sözlükleri** (TLP, PAP, kill-chain, confidence, vb.). Makine-okunur namespace:predicate=value formatında. → Bölüm 7.
- **Galaxy:** Zengin, **bilgi-kümeleri (cluster)** — aktörler, teknikler, malware, sektörler. Bir galaxy cluster'ı bir event'e yapıştırınca o event'e koca bir bilgi gövdesi (aktörün takma adları, ülkesi, hedef sektörleri vb.) bağlanır. → Bölüm 6.

### 4.5 Sighting — "bu IOC'yi ben de gördüm"

Bir **Sighting**, bir attribute'a iliştirilen "ben bu IOC'yi *şu tarihte* gördüm" kaydıdır. Üç tipi vardır: **gözlem (true positive)**, **false-positive**, **expiration (artık geçerli değil)**. Sighting'ler IOC'nin **tazeliğini, yaygınlığını ve gerçekliğini** ölçer. → Bölüm 15.

---

<a id="5"></a>
## 5. 🔥 Bir Event'i Sıfırdan Oluşturma — Phishing Kampanyası

Teoriyi pratikle birleştirelim. Senaryo: SOC'una düşen bir **bordro temalı phishing** olayı. Gönderen `hr@evil-payroll.tld`, içinde `http://evil-payroll.tld/bordro.html` linki, link bir `bordro.xlsm` (makro) indiriyor, makro `203.0.113.45` adresindeki C2'ye bağlanıp `update.exe` (SHA-256 `cd34...`) çekiyor. Aktör izleri **APT-benzeri**, teknik **spearphishing link + macro execution**.

### Adım 1 — Event oluştur (kabı kur)

`Add Event` (Event Actions → Add Event):

| Alan | Değer |
|---|---|
| **Date** | 2026-05-28 (gözlem tarihi) |
| **Distribution** | `Your organisation only` (başlangıçta KAPALI — olgunlaşınca açarsın) |
| **Threat Level** | High |
| **Analysis** | Ongoing |
| **Info** | `Bordro Temalı Spearphishing — Macro Loader → C2 (203.0.113.45) — 2026-05` |

> 🔑 **Püf — daima kapalı başla.** Yeni event'i **`Your org only`** ile aç, doğrula/temizle/etiketle, sonra `publish` ederken bilinçli olarak distribution'ı yükselt. "Önce community'ye aç sonra düşün" → düzeltemeyeceğin sızıntı.

### Adım 2 — IOC'leri attribute olarak ekle

`Add Attribute` ile tek tek, ya da `Populate from → Freetext Import` ile toplu (MISP metni parse edip türleri otomatik önerir):

| category | type | value | to_ids | comment |
|---|---|---|---|---|
| Payload delivery | `email-src` | `hr@evil-payroll.tld` | true | Phishing gönderen |
| Payload delivery | `email-subject` | `2026 Mayıs Bordronuz` | false | Bağlam (tespit değil) |
| Network activity | `url` | `http://evil-payroll.tld/bordro.html` | true | Phishing landing |
| Network activity | `domain` | `evil-payroll.tld` | true | Kötücül alan |
| Network activity | `ip-dst` | `203.0.113.45` | true | C2 sunucusu |

> 💡 **Freetext Import gücü:** IR'dan gelen ham notu (IP/hash/URL karışık) yapıştır → MISP türleri otomatik tanır, sen onaylarsın. En hızlı toplu giriş yolu. Ama **her satırı gözden geçir** — yanlış tür ataması korelasyonu bozar.

### Adım 3 — Dosyaları Object olarak ekle (yapılandırılmış)

İki dosya var; her biri `file` objesi (`Add Object → file`):

**Object 1 — bordro.xlsm (ilk aşama):**
```
file objesi:
  filename     = bordro.xlsm
  md5          = 11aa...
  sha1         = 22bb...
  sha256       = ab12...
  size-in-bytes= 48213
  mimetype     = application/vnd.openxmlformats... (xlsm)
```

**Object 2 — update.exe (payload):**
```
file objesi:
  filename     = update.exe
  sha256       = cd34...
  imphash      = 99ff...   ← varyant avlamada altın
  size-in-bytes= 215040
```

**İlişkileri bağla (Add Reference):**
```
[url: .../bordro.html] ──"downloads"──► [file: bordro.xlsm]
[file: bordro.xlsm]    ──"drops"─────► [file: update.exe]
[file: update.exe]     ──"connects-to"► [ip-dst: 203.0.113.45]
```

Artık event sadece "5 IOC" değil; **bir saldırı zincirinin makine-okunur haritası**.

### Adım 4 — Etiketle (TLP + Galaxy + Taxonomy)

- **TLP:** `tlp:amber` (topluluk içi paylaşılabilir ama dışarı değil) — event seviyesinde.
- **Kill-chain / kategori:** `kill-chain:Delivery`, `kill-chain:Command and Control`.
- **Confidence:** `admiralty-scale:source-reliability="B"` veya `misp-confidence` taxonomy.
- **Galaxy — teknik:** ATT&CK pattern cluster'ları ekle (Bölüm 6): **T1566.002 (Spearphishing Link)**, **T1204.002 (Malicious File / User Execution)**, **T1071.001 (Application Layer Protocol: Web)**.
- **Galaxy — aktör (atıf varsa):** İlişkili threat-actor cluster'ı — **ama atıf kesin değilse ekleme** (Bölüm 16, OPSEC).

### Adım 5 — to_ids disiplinini son kez gözden geçir

- `email-subject` → **to_ids=false** (genel ifade, SIEM'i şişirir).
- `domain`, `ip-dst`, `url`, dosya `sha256`'ları → **to_ids=true** (bloklanabilir, ayırt edici).
- Genel `filename` (`update.exe` — çok yaygın ad!) → **to_ids=false** (tek başına false-positive üretir; hash zaten yakalar).

### Adım 6 — Warning list kontrolü & Publish

- **Enrichment/Warninglist** çalıştır: IOC'lerin known-good listede (Alexa/Umbrella/bogon) olup olmadığını kontrol et (Bölüm 12). `203.0.113.45` bir TEST-NET adresi → gerçek olayda bu bir uyarı verirdi; gerçek IOC'de bogon/özel IP'leri **temizle**.
- **Publish:** Distribution'ı bilinçli yükselt (örn. `Sharing group: Finans ISAC`), TLP'yi son kez doğrula, **Publish**. Bu an paylaşım kararının verildiği andır.

> 🧠 **Sonuç:** Artık bu event, başka bir kurum `203.0.113.45` ya da `ab12...` hash'ini gördüğünde **otomatik korele olur** (Bölüm 8) ve `to_ids` IOC'leri Suricata/Sigma export'una düşerek **gerçek bloklama** üretir (Bölüm 11).

---

<a id="6"></a>
## 6. 🔥 GALAXY — ATT&CK, Threat Actor, Ransomware Cluster'larını Bağlama

**Galaxy**, MISP'in IOC'yi **bağlama oturtma** mekanizmasıdır. Bir galaxy, **cluster**'lardan (bilgi-kümeleri) oluşur; her cluster zengin meta-veri taşır (takma adlar, açıklama, kaynak, ilişkiler). Bir cluster'ı event'e/attribute'a yapıştırınca, o koca bilgi gövdesi olayına bağlanır.

### En kritik galaxy'ler

| Galaxy | İçerik | Kullanım |
|---|---|---|
| **MITRE ATT&CK** (attack-pattern, course-of-action, software, group) | Teknik/taktikler (T-numaraları), yazılım, gruplar | Event'e **TTP** etiketleme — "bu olay hangi tekniği kullandı?" |
| **Threat Actor** | Aktör profilleri (takma adlar, ülke, hedef sektör, motivasyon) | Atıf — "bu kampanya hangi aktöre ait?" |
| **Ransomware** | Fidye aileleri (şifreleme, fidye notu, ödeme TTP'leri) | Olayı bir ransomware ailesine bağlama |
| **Tool / mitre-tool** | Saldırgan/ikili-kullanım araçları (Cobalt Strike, Mimikatz) | Kullanılan aracı işaretleme |
| **Sector** | Hedef sektörler | Hedefleme analizi |
| **Country** | Coğrafi bağlam | Hedef/kaynak coğrafyası |
| **Malpedia / Microsoft activity group** | Malware aileleri / MS aktör adları | Çapraz-adlandırma eşleme |

### Bir event'e ATT&CK technique cluster ekleme

Event görünümünde **Galaxies → Add → MITRE ATT&CK → Attack Pattern** → T-numarasını ara (örn. `T1566.002`) → ekle. Artık event'in galaxy bölümünde:

```
🌌 Galaxies
   ├─ ATT&CK Pattern: T1566.002 — Phishing: Spearphishing Link
   ├─ ATT&CK Pattern: T1204.002 — User Execution: Malicious File
   ├─ ATT&CK Pattern: T1071.001 — Application Layer Protocol: Web Protocols
   └─ Threat Actor: (yalnızca atıf güvenliyse)
```

Bu cluster'lar makine-tag olarak da görünür: `misp-galaxy:mitre-attack-pattern="Spearphishing Link - T1566.002"`. Bu sayede **MISP'ten doğrudan bir ATT&CK Navigator layer'ı üretebilirsin** (Bölüm 17).

### Threat Actor cluster — atıf

Threat-actor galaxy'sine bir aktör eklemek olayı o aktörün tüm bilinen profiline (takma adları: APT28=Fancy Bear=Sofacy; hedefleri; geçmiş kampanyaları) bağlar.

> 🔥 **ATIF UYARISI:** Bir threat-actor cluster'ı yapıştırmak **iddiadır**. Yanlış aktör ataması, raporunu çürütür ve topluluğu yanıltır. Aktör cluster'ını **yalnızca yüksek güvenle** (çoklu bağımsız kanıt, kendi telemetrin + güvenilir kaynak teyidi) ekle. Şüphedeysen aktör değil **yalnızca TTP** (ATT&CK) etiketle — teknikler gözlemdir, atıf yorumdur.

> 📌 **ATT&CK DERİNLİĞİ BURADA DEĞİL:** MISP'in galaxy ile ATT&CK'i nasıl *bağladığını* burada gördün. Tekniklerin/taktiklerin **derinlemesine anlamı, alt-teknik mantığı, kill-chain eşlemesi, data source/detection** konuları ve **TTP avcılığı teorisi** ayrı dosyalarda işlenir:
> - **`MITRE_ATTACK_USTALIK_REHBERI.md`** — ATT&CK matrisinin derin kullanımı, Navigator, technique seçimi.
> - **`TTP_AVCILIGI_USTALIK_REHBERI.md`** — TTP teorisi, davranışsal avlama, Pyramid of Pain.
> Burada tekrarlamıyoruz; MISP tarafında galaxy = **bağlama oturtma** aracıdır, o kadar.

---

<a id="7"></a>
## 7. 🏷️ TAXONOMY & Machine-Tag Mantığı

**Taxonomy**, standart **etiket sözlükleridir**. Galaxy "kim/ne/hangi teknik" sorularını (zengin cluster) yanıtlarken, taxonomy **sınıflandırma** sorularını (TLP, güven, PAP, kill-chain aşaması) yanıtlar. Hepsi **machine-tag** formatındadır:

```
namespace : predicate = value
   tlp    :   amber
   tlp    :   red
admiralty-scale : source-reliability = "B"
   PAP   :   amber
kill-chain :  Exploitation
estimative-language : likelihood-probability = "likely"
```

Machine-tag'in gücü: makine ayrıştırabilir, filtreleyebilir, otomasyon kuralları yazabilirsin ("tlp:red olan hiçbir şeyi şu org'a sync etme").

### Faydalı taxonomy'ler

| Taxonomy | Ne için |
|---|---|
| **tlp** | Paylaşım sınırı (Traffic Light Protocol) — Bölüm 9 |
| **PAP** (Permissible Actions Protocol) | IOC ile **ne yapabilirsin** (sadece pasif gözlem mi, aktif bloklama mı?) — atıf yakmamak için |
| **admiralty-scale** | Kaynak güvenilirliği (A-F) + bilgi doğruluğu (1-6) |
| **estimative-language** | Belirsizlik ifadesi (likely/probable) — ICD 203 uyumlu |
| **kill-chain** | Lockheed Martin kill-chain aşaması |
| **confidence-level / misp-confidence** | IOC güven seviyesi |
| **course-of-action** | Önerilen savunma aksiyonu |
| **workflow** | İç süreç (to-review, complete, rejected) |

> 🔑 **PAP, TLP'nin az bilinen kardeşidir.** TLP "kiminle paylaşabilirim?" der; **PAP** "bu IOC ile **ne yapabilirim**?" der. `PAP:RED` → "sadece pasif, hedefi uyarma" (bir C2'yi bloklarsan aktör altyapısının yandığını anlar). `PAP:GREEN` → "serbestçe blokla/aksiyona geç". OPSEC-hassas operasyonlarda PAP, TLP kadar önemlidir.

> ⚠️ **Over-tagging tuzağı:** Her olaya 30 taxonomy etiketi yapıştırmak gürültüdür ve filtrelemeyi bozar. **Az ama doğru** etiketle: TLP (zorunlu), confidence, kill-chain, gerekirse PAP. Geri kalanı yalnızca anlamlıysa. Bkz. Bölüm 16.

---

<a id="8"></a>
## 8. 🔥 KORELASYON — Otomatik IOC İlişkilendirme & Gürültü Yönetimi

Korelasyon, MISP'in **en değerli ve en yanlış anlaşılan** özelliğidir. MISP, bir attribute eklendiğinde otomatik olarak **aynı değere sahip diğer tüm event'lerdeki attribute'larla** ilişki kurar. Manuel iş yok — sen `203.0.113.45`'i eklediğin an, MISP onu daha önce bu IP'yi içeren her event'e bağlar.

```
        EVENT A (senin phishing olayın)
         ip-dst: 203.0.113.45 ●───────┐
                                       │  KORELASYON (otomatik)
        EVENT B (3 ay önceki, başka org)│
         ip-dst: 203.0.113.45 ●────────┤
                                       │
        EVENT C (ransomware kampanyası) │
         ip-dst: 203.0.113.45 ●────────┘
                  │
                  ▼
   "Bu C2 IP'si 3 ayrı kampanyada görüldü → muhtemelen
    paylaşılan altyapı / aynı aktör kümesi" = İSTİHBARAT
```

### "Bu IOC başka hangi olaylarda görüldü?"

Her attribute'un yanında **korelasyon göstergesi** belirir; tıkladığında o IOC'nin geçtiği tüm event'leri görürsün. Event görünümünde **Correlation Graph** ile ilişkileri görsel olarak gezersin:

```
   [Event A]───203.0.113.45───[Event B]
       │                          │
   evil-payroll.tld          cd34...(sha256)
       │                          │
   [Event D]                  [Event C]───[Threat Actor: APT-X]
```

Bu graf, "tekil bir IOC"yi "kampanya ağına" çeviren şeydir. Pyramid of Pain üst basamakları (TTP, aktör) işte bu korelasyonla ortaya çıkar.

### 🔥 Korelasyon GÜRÜLTÜSÜ — ustalığın ayrıldığı yer

Korelasyon kör çalışırsa **felaket** olur. Düşün: `8.8.8.8` (Google DNS) yanlışlıkla bir IOC olarak girilmiş. MISP onu içeren **her olayı** birbirine bağlar → binlerce sahte ilişki, anlamsız bir spagetti grafiği, ve gerçek ilişkilerin gizlenmesi. Gürültüyü yönetmenin araçları:

| Araç | Ne yapar |
|---|---|
| **Warning lists** (Bölüm 12) | Known-good (Google/Cloudflare/Alexa/bogon) değerlerin korelasyonunu işaretler/bastırır |
| **disable_correlation bayrağı** | Belirli bir attribute için korelasyonu kapat (örn. çok yaygın bir değer) |
| **Correlation exclusion list** | Hiç korele edilmeyecek değerler (kurumsal IP'ler, bilinen sinkhole'lar) |
| **over-correlating value eşiği** | Bir değer çok fazla event'le korele oluyorsa MISP onu otomatik "over-correlating" işaretler ve grafikten dışlar |
| **Sinkhole/parking IP farkındalığı** | Ele geçirilmiş domain'ler genelde aynı sinkhole IP'ye işaret eder → sahte korelasyon kaynağı |

> 🔥 **Altın kural — korelasyon kalite ister, miktar değil.** "Çok korelasyon = iyi" YANLIŞ. Anlamlı korelasyon, **ayırt edici (yüksek-entropili)** IOC'lerden gelir: bir SHA-256, bir benzersiz C2 domain'i, bir imphash. **Düşük-entropili** değerler (yaygın IP, genel filename, paylaşılan hosting IP'si) korelasyonu **zehirler**. Warning-list ve correlation-exclusion **olmazsa olmazdır** — onlarsız MISP birkaç ay içinde kullanılamaz bir gürültü grafiğine döner.

> 🧠 **Sinkhole tuzağı:** Bir kötücül domain "takedown" edildiğinde genelde bir **sinkhole IP**'sine yönlendirilir. Onlarca alakasız kampanyanın domain'i aynı sinkhole IP'sine bakar → MISP hepsini birbirine bağlar. Bu **sahte korelasyondur**; sinkhole IP'leri exclusion list'e ekle ya da warning-list (sinkhole) kullan.

---

<a id="9"></a>
## 9. 🔥 TLP — Traffic Light Protocol Disiplini

**TLP (Traffic Light Protocol)**, istihbaratın **kiminle paylaşılabileceğini** belirleyen evrensel sınıflandırma standardıdır (FIRST tarafından sürdürülür). MISP'te `tlp` taxonomy'si olarak machine-tag'lenir ve **paylaşım kararının çekirdeğidir**.

> 📌 **Sürüm notu — TLP 2.0:** Güncel standart **TLP 2.0**'dır ve `TLP:CLEAR` ile `TLP:AMBER+STRICT` getirmiştir. Eski **TLP 1.0**'daki `TLP:WHITE` artık **`TLP:CLEAR`** olarak adlandırılır. MISP her iki etiketi de taşır; kurumunun ve paylaşım topluluğunun hangi sürümü kullandığını **teyit et** (resmi: first.org/tlp).

| TLP (2.0) | Renk | Kiminle paylaşılır | Tipik kullanım |
|---|---|---|---|
| **TLP:RED** | 🔴 | **Yalnızca toplantıdaki/doğrudan belirtilen kişiler.** Daha geniş paylaşılamaz | Çok hassas, atıf-riskli, aktif operasyon; isimli alıcılar |
| **TLP:AMBER+STRICT** | 🟠 | Yalnızca **kendi organizasyonu** içinde (müşteri/3. tarafa bile değil) | Org-içi hassas; harici paylaşım yasak |
| **TLP:AMBER** | 🟠 | Kendi org + **bilmesi gereken müşteriler/üyeler** | ISAC içi, sınırlı topluluk |
| **TLP:GREEN** | 🟢 | **Topluluk** geneli (sektör/güven çemberi), ama **kamuya/internete değil** | Sektörel paylaşım, geniş ama kapalı |
| **TLP:CLEAR** (eski WHITE) | ⚪ | **Sınırsız** — kamuya açık paylaşılabilir | Açık OSINT, yayın izni olan bilgi |

### TLP ↔ MISP distribution eşlemesi (kritik)

TLP bir **politika etiketi**; MISP distribution bir **teknik kontrol**. İkisini **tutarlı** tutmak zorundasın, yoksa etiket "AMBER" der ama sistem onu "all communities"e sync eder → **felaket**.

| TLP | Önerilen MISP distribution |
|---|---|
| TLP:RED | `Your organisation only` (0) — hatta sharing group'ta isimli alıcılar |
| TLP:AMBER / AMBER+STRICT | `Your org only` (0) ya da dar bir **Sharing group** (4) |
| TLP:GREEN | `This community` (1) ya da `Connected communities` (2) |
| TLP:CLEAR | `All communities` (3) — yalnızca gerçekten kamuya açıksa |

> 🔥 **YANLIŞ PAYLAŞIM FELAKETİ — somut senaryo:** Bir IR analisti, müşterinin ele geçirildiğini gösteren bir event'i (içinde müşteri iç IP'leri, çalışan e-postaları) **TLP:AMBER** olması gerekirken yanlışlıkla **distribution: All communities** ile publish eder. Event saniyeler içinde sync zinciriyle onlarca kuruma yayılır → **müşterinin ihlali, iç altyapısı ve çalışan PII'si geri alınamaz biçimde ifşa olur.** Sözleşme ihlali, yasal sorumluluk, itibar kaybı. **Bu yüzden: publish öncesi TLP+distribution iki kez doğrulanır; otomasyon kuralı tlp:red/amber'i asla geniş sync etmeyecek şekilde kurulur.**

> 🧠 **TLP'ye saygı iki yönlüdür:** Sana TLP:AMBER gelen bir IOC'yi sen **daha geniş** paylaşamazsın (downgrade edemezsin). Kaynak ne dediyse o sınır geçerlidir. TLP'yi **yükseltmek** (daha kapalı yapmak) serbest; **düşürmek** güven ihlalidir.

---

<a id="10"></a>
## 10. 🌐 FEEDS — CIRCL OSINT, Açık Kaynak Feed'ler & Güvenilirlik

**Feed**, dış kaynaklardan otomatik IOC çeken bir akıştır. MISP iki tip feed konuşur: **MISP-format feed** (başka MISP'ten/JSON event) ve **freetext/CSV feed** (ham IOC listeleri).

### Önemli açık feed'ler

| Feed | Kaynak | İçerik |
|---|---|---|
| **CIRCL OSINT Feed** | CIRCL (MISP-format) | Kürate edilmiş, yüksek kaliteli — **varsayılan başlangıç** |
| **Botvrij.eu** | OSINT (MISP-format) | Çeşitli OSINT IOC'leri |
| **abuse.ch** (URLhaus, Feodo Tracker, SSLBL, MalwareBazaar, ThreatFox) | abuse.ch | Aktif malware URL/C2/SSL/hash — yüksek değer |
| **Tor exit nodes** | (CSV) | Tor çıkış IP'leri (bağlam için, IOC değil) |
| **CINS / blocklist.de / Emerging Threats** | Çeşitli | IP reputation/blocklist |

### Feed ekleme/etkinleştirme

Sync Actions → Feeds → Add Feed → URL + format (MISP/freetext/CSV) + distribution + (varsa) tag.

```
1. Feed ekle (URL, format, hedef distribution)
2. "Fetch and store all feed data" YA DA cache'le (önce önizleme için cache önerilir)
3. Cache'lenen feed'i ÖNİZLE → IOC kalitesini gör → sonra import et
4. Periyodik fetch (cron/worker) ile güncel tut
```

> 🔑 **Püf — önce cache, sonra import.** Bir feed'i körlemesine "fetch & store" etme. **Cache** edip **önizle**; içinde ne tür IOC var, ne kadar gürültülü, warning-list'e takılan çöp var mı gör. Beğenirsen import et. Aksi halde MISP'ine binlerce şüpheli IOC'yi geri-dönüşsüz boca edersin.

### 🔥 Feed güvenilirliği & zehirlenme

| Risk | Açıklama | Önlem |
|---|---|---|
| **Düşük kalite** | Feed yanlış-pozitif dolu (meşru IP/domain'ler IOC diye) | Warning list + güvenilen feed seç + cache-önizle |
| **Zehirleme (poisoning)** | Kötü niyetli/ele geçirilmiş feed kasıtlı yanlış IOC enjekte eder (örn. Google'ı "C2" yapar → bloklama → DoS) | Sadece **güvenilir, köklü** feed; feed'i to_ids=false ile getir; auto-block etme |
| **Stale (bayat)** | Feed güncellenmiyor, IOC'ler ölü | Tazelik kontrolü, expiration |
| **Atıf yanlışı** | Feed yanlış aktör/aile etiketi taşır | Etiketleri körü körüne devralma; teyit et |

> 🔥 **Feed zehirlenmesi gerçek bir tehdittir.** Düşün: bir saldırgan bir OSINT feed'ine `1.1.1.1` (Cloudflare DNS) ya da kurumsal banka IP'lerini "C2" olarak enjekte etmeyi başarır. Sen otomatik bloklama (auto-block) kuruyorsan → kendi kullanıcılarını meşru servislerden **kesersin** (öz-DoS). **Bu yüzden: dış feed'lerden gelen IOC'leri varsayılan `to_ids=false` getir, kritik bloklamayı sadece kürate/teyit edilmiş IOC'lerle yap, ve warning-list'i her zaman aktif tut.** Güven, kaynağın köklülüğüyle (CIRCL, abuse.ch gibi) orantılıdır.

---

<a id="11"></a>
## 11. 🔄 IMPORT / EXPORT — STIX, OpenIOC, CSV & SIEM/IDS Beslemesi

MISP bir **çeviri merkezi**dir: içeri çeşitli formatlardan alır, dışarı çeşitli formatlara verir.

### Desteklenen formatlar

| Format | Yön | Not |
|---|---|---|
| **MISP JSON** | ↔ | Yerli format; MISP↔MISP en sadık aktarım |
| **STIX 1.x (XML)** | ↔ | Eski ama hâlâ yaygın (TAXII 1.x) |
| **STIX 2.0 / 2.1 (JSON)** | ↔ | Modern standart (TAXII 2.x); OpenCTI/diğer araçlarla birlikte çalışma |
| **OpenIOC** | ↔ | Mandiant formatı (eski raporlar) |
| **CSV** | ↔ | Basit IOC listeleri; sütun eşleme |
| **Suricata / Snort** | → (export) | IDS imza kuralları (to_ids=true attribute'lardan) |
| **Sigma** | → | SIEM-agnostik tespit kuralı |
| **Bro/Zeek intel** | → | Zeek intel framework formatı |
| **STIX bundle** | → | Rapor paylaşımı |
| **(text/CSV/JSON API)** | → | SIEM/EDR ham besleme |

> ⚠️ **STIX uyumluluğu pürüzsüz değildir.** MISP veri modeli ile STIX 2.1 nesne modeli **birebir örtüşmez**; çevirilerde bazı alanlar (özel galaxy'ler, MISP-özel objeler) kaybolabilir ya da custom-property'ye düşer. STIX'e çevirip başka araca verirken **kayıp olup olmadığını doğrula**. STIX 1↔2 dönüşümü ayrıca kayıplıdır. Kritik aktarımda her iki uçta da **örnek bir event ile test et** ve alan kaybını teyit et (resmi STIX mapping dokümanından kontrol et).

### SIEM/IDS'e besleme (en yaygın tüketim)

MISP'in asıl iş değeri: **to_ids=true IOC'leri savunma araçlarına otomatik beslemek.**

```
   MISP (to_ids=true attribute'lar)
        │
        ├──► Suricata/Snort rules ──► IDS/IPS (ağ tespiti)
        ├──► Sigma rules ──────────► SIEM (Splunk/Elastic/Sentinel)
        ├──► Zeek intel ───────────► Zeek (ağ izleme)
        ├──► STIX/TAXII ───────────► EDR / diğer platform
        └──► PyMISP (API) ─────────► özel entegrasyon (Bölüm 14)
```

**Tipik akış:** Bir SIEM/EDR, MISP'in **export API** ucundan periyodik olarak `to_ids=true` IOC'leri çeker (örn. son 24 saat, belirli tag/TLP filtresiyle) ve kendi tespit listesine yükler. Suricata için MISP doğrudan `.rules` üretir.

> 🔑 **Püf — filtreli export.** SIEM'e **her şeyi** verme. `to_ids:true` + güncel + uygun TLP + yüksek-confidence filtresiyle çek. Yoksa SIEM'in alarm yağmuruna boğulur (ve gerçek alarm gürültüde kaybolur — Bölüm 8'in aynı dersi tüketim tarafında).

---

<a id="12"></a>
## 12. 🔥 WARNING LISTS & False-Positive Bastırma

**Warning list**, bilinen **iyi (known-good)** ya da problemli değerlerin listeleridir. MISP, bir IOC bir warning-list'le eşleştiğinde onu **işaretler** (otomatik silmez — sana "dikkat, bu meşru olabilir" der). Bu, false-positive'i bastıran ve korelasyon gürültüsünü kesen birincil savunmadır.

### Önemli warning list'ler

| Warning list | İçerik | Niçin |
|---|---|---|
| **Alexa / Cisco Umbrella / Majestic top 1M** | En popüler domain'ler | `google.com` IOC olarak girmesin |
| **Cisco Umbrella / Tranco top sites** | Popüler siteler | Aynı |
| **Bogon / RFC1918 / reserved IP** | Özel/ayrılmış/yönlendirilemez IP'ler | `192.168.x`, `10.x`, `203.0.113.x` (TEST-NET) IOC olamaz |
| **Public DNS resolvers** | Google/Cloudflare/Quad9 (8.8.8.8, 1.1.1.1) | C2 sanılmasın |
| **CDN/cloud ranges (AWS/Azure/GCP/Akamai/Cloudflare)** | Paylaşımlı hosting IP'leri | Sahte korelasyon kaynağı; bloklamak meşru servisi keser |
| **Microsoft/Google/Apple known domains** | Telemetri/güncelleme uçları | False-positive |
| **TLD list / sinkhole list** | Geçerli TLD'ler / bilinen sinkhole'lar | Geçersiz domain / sahte korelasyon |
| **Common hashes (boş dosya, sık DLL)** | `d41d8cd9...` (boş MD5) vb. | Anlamsız hash IOC'leri |

> 🔥 **Warning-list olmadan MISP kullanılamaz hale gelir.** Bir analist `8.8.8.8`'i (malware DNS sorgusu yaptı diye) C2 sanıp IOC girer → warning-list yoksa bu SIEM'e gider, Google DNS bloklanır, **tüm kurum DNS'siz kalır.** Ya da CDN IP'si bloklanır, yarısı çalışmayan bir internet. Warning-list, bu felaketleri **publish ANINDA** "bu IOC bir public resolver/CDN/bogon — emin misin?" diye uyararak önler. **İlk kurulumda tüm temel warning-list'leri etkinleştir.**

> 🧠 **Warning-list ≠ silme.** MISP IOC'yi **silmez**, sadece bayrak takar — çünkü bazen meşru altyapı gerçekten kötüye kullanılır (ele geçirilmiş bir popüler site). Karar senindir: uyarıyı görüp **bilinçli** karar verirsin (genelde to_ids=false yaparsın ya da çıkarırsın).

---

<a id="13"></a>
## 13. 🤝 SHARING — Sharing Groups, Organizasyonlar & Sync Server

### Organisations

MISP'te her veri bir **organisation**'a aittir. İki tür org vardır: **local** (senin örneğindeki gerçek kullanıcılı org'lar) ve **external** (sadece paylaşımda atıf için tanınan, kullanıcısı olmayan org'lar). Her org'un bir **UUID**'si vardır — paylaşımda kimliğidir.

### Sharing Groups — granüler paylaşım

`distribution: All communities` çok geniş, `community only` çok dar olduğunda **Sharing Group** devreye girer: **tam olarak hangi org'larla** paylaşacağını sen seçersin.

```
   Sharing Group: "Finans-ISAC Çekirdek"
   ┌──────────────────────────────────┐
   │ ✓ Banka A   (extend: hayır)       │  ← extend=hayır: yeniden paylaşamaz
   │ ✓ Banka B   (extend: hayır)       │
   │ ✓ CERT-Fin  (extend: EVET)        │  ← extend=evet: kendi çemberine yayabilir
   └──────────────────────────────────┘
   Bu event SADECE bu 3 org'a gider. Başkası göremez.
```

**Extend bayrağı** kritik: bir org'a "extend" verirsen, o org event'i **kendi** paylaşım çemberine de yayabilir. Vermezsen, paylaşım orada durur.

### Sync Server — organizasyonlar arası senkronizasyon

İki MISP örneği **sync** ile bağlanır: biri diğerinden event çeker (pull) ve/veya iter (push). Sync, **sync user** kimlik bilgisi (API key) ve **dikkatli distribution/TLP filtresi** ile kurulur.

```
   MISP-A (sizin) ──push (yalnızca tlp:green+)──► MISP-B (ISAC merkez)
   MISP-A (sizin) ◄──pull (community feed)────── MISP-B (ISAC merkez)
```

> 🔥 **Sync, TLP felaketinin en hızlı yayıldığı kanaldır.** Bir sync push kuralını "tüm event'leri it" diye kurarsan, **TLP:RED bir event saniyeler içinde karşı tarafa gider.** Sync kurallarını **her zaman distribution/TLP filtresiyle** kur: "yalnızca distribution≥community ve tlp:red DEĞİL olanları push et". Push/pull rule'larını ve "push edilen org listesini" kurulumda **iki kez** doğrula. Sync API key'i bir admin yetkisidir — sızarsa tüm paylaşım kontrolün gider (Bölüm 16).

> 🔑 **Push vs Pull tercihi:** **Pull** (sen merkezi çekersin) daha güvenli başlangıçtır — kontrolü sende tutar, yanlışlıkla bir şey *itmezsin*. **Push** (sen iletirsin) güçlüdür ama her yanlış filtre bir sızıntıdır. Yeni başlayan: **pull-only** ile başla.

---

<a id="14"></a>
## 14. 🔥 PyMISP / API — Otomasyon (Kod Örnekleriyle)

**PyMISP**, MISP'in resmi Python kütüphanesidir. Manuel arayüz öğrenme içindir; gerçek operasyon **otomasyondur** — feed işleme, IR çıktısı→event, IOC sorgulama, SIEM beslemesi hep API ile.

### Kurulum & bağlantı

```python
# pip install pymisp
from pymisp import PyMISP, MISPEvent, MISPAttribute, MISPObject

MISP_URL = "https://misp.kurumun.local"
MISP_KEY = "API_ANAHTARIN"           # Profil → Auth keys (ASLA koda gömme — env/secret)
VERIFY_TLS = True                     # üretimde DAİMA True

misp = PyMISP(MISP_URL, MISP_KEY, ssl=VERIFY_TLS)
```

> 🔑 **API key = parola gücünde sır.** Asla kaynak koda/repoya gömme. Ortam değişkeni ya da secret manager kullan. Her entegrasyon için **ayrı, dar yetkili** key üret (sadece-okuma entegrasyonu için read-only key). Key sızarsa **derhal iptal et** (Profil → Auth keys → revoke). Bkz. Bölüm 16.

### Programatik event oluşturma (Bölüm 5'in kod karşılığı)

```python
event = MISPEvent()
event.info = "Bordro Spearphishing — Macro Loader → C2 — 2026-05"
event.distribution = 0          # Your org only (kapalı başla)
event.threat_level_id = 1       # 1=High, 2=Medium, 3=Low, 4=Undefined
event.analysis = 1              # 0=Initial, 1=Ongoing, 2=Completed
event.add_tag("tlp:amber")
event.add_tag('misp-galaxy:mitre-attack-pattern="Spearphishing Link - T1566.002"')

# Attribute'lar (to_ids disiplinine dikkat)
event.add_attribute("ip-dst", "203.0.113.45", category="Network activity", to_ids=True,
                    comment="C2 sunucusu")
event.add_attribute("domain", "evil-payroll.tld", category="Network activity", to_ids=True)
event.add_attribute("email-src", "hr@evil-payroll.tld", category="Payload delivery", to_ids=True)
event.add_attribute("email-subject", "2026 Mayis Bordronuz", category="Payload delivery",
                    to_ids=False)   # bağlam → tespit değil

# File object (yapılandırılmış)
fobj = MISPObject("file")
fobj.add_attribute("filename", "update.exe", to_ids=False)  # genel ad → to_ids kapalı
fobj.add_attribute("sha256", "cd34" + "0"*60, to_ids=True)
fobj.add_attribute("imphash", "99ff" + "0"*28, to_ids=True)
event.add_object(fobj)

created = misp.add_event(event, pythonify=True)
print("Event ID:", created.id, "UUID:", created.uuid)
```

### Arama / IOC çekme

```python
# 1) Belirli bir IOC kurumumuzda görülmüş mü? (korelasyon sorgusu)
res = misp.search(controller="attributes", value="203.0.113.45", pythonify=True)
for attr in res:
    print(attr.event_id, attr.type, attr.value, attr.to_ids)

# 2) Son 24 saatte eklenen, IDS'e gidecek (to_ids) IOC'leri SIEM için çek
iocs = misp.search(controller="attributes",
                   to_ids=True,
                   timestamp="24h",
                   type_attribute=["ip-dst", "domain", "url", "sha256"],
                   pythonify=True)
ip_block_list = [a.value for a in iocs if a.type == "ip-dst"]

# 3) Belirli tag/TLP ile event ara (rapor için)
events = misp.search(controller="events",
                     tags=["tlp:green", 'misp-galaxy:ransomware%'],
                     pythonify=True)

# 4) Tag'e göre ATT&CK technique'leri topla (layer üretimi için — Bölüm 17)
attack_events = misp.search(controller="events",
                            tags='misp-galaxy:mitre-attack-pattern%',
                            pythonify=True)
```

### Sighting ekleme (API ile "ben de gördüm")

```python
from pymisp import MISPSighting
s = MISPSighting()
s.value = "203.0.113.45"
s.source = "SOC-EDR"
s.type = "0"    # 0=gözlem(true positive), 1=false-positive, 2=expiration
misp.add_sighting(s)
```

### Export / feed otomasyonu

```python
# Suricata kurallarını al (to_ids IOC'lerden) — IDS'e besle
suricata_rules = misp.search(controller="attributes", to_ids=True,
                             return_format="suricata")

# STIX 2.x bundle olarak bir event'i dışa ver
stix = misp.search(controller="events", eventid=created.id, return_format="stix2")
```

> 🧠 **Otomasyon mimarisi:** Tipik üretim akışı — (1) **IR/sandbox/honeypot** → PyMISP ile event oluştur; (2) **cron job** her saat `to_ids` IOC'leri çekip SIEM/EDR'a basar; (3) **SIEM eşleşmesi** → PyMISP ile **sighting** geri yazar (Bölüm 15); (4) **expiration job** bayat IOC'leri to_ids=false yapar. Döngü kapanır: MISP hem üretir hem geri-besleme alır.

---

<a id="15"></a>
## 15. 👁️ SIGHTINGS — "Ben de Gördüm" & IOC Tazeliği/Yaygınlığı

**Sighting**, bir attribute'a "bu IOC'yi şu tarihte/kaynakta gördüm" diyen bir kayıttır. Üç tipi:

| Tip | Anlamı | Etkisi |
|---|---|---|
| **0 — Sighting (true positive)** | "Bu IOC'yi gerçekten gördüm/yakaladım" | IOC'nin **aktifliğini ve yaygınlığını** artırır |
| **1 — False positive** | "Bu IOC bende yanlış alarm üretti / meşru" | Güveni düşürür, gözden geçirmeyi tetikler |
| **2 — Expiration** | "Bu IOC artık geçerli değil (C2 kapandı)" | Tazeliği düşürür → to_ids kapatma adayı |

### Neden sighting değerlidir?

- **Doğrulama:** 5 farklı kurum bağımsız "gördüm" derse, IOC'nin gerçekliği matematiksel olarak güçlenir. 1 kişinin tek girdisi ≠ 50 kurumun teyidi.
- **Tazelik (freshness):** Son sighting'i 6 ay önce olan bir C2 IP'si muhtemelen ölüdür → bloklamaya değmez, gürültü. Yeni sighting'li IOC önceliklidir.
- **Yaygınlık (prevalence):** Bir IOC çok yerde görülüyorsa ya çok aktif bir kampanyadır ya da... false-positive (yaygınlık ikisini de gösterebilir — bağlamla yorumla).
- **Önceliklendirme:** SOC, "taze + çok-sighting'li + yüksek-confidence" IOC'leri öne alır; bayat/tek-kaynak olanları arka plana atar.

> 🔥 **Sighting, IOC tazeliği problemini çözen tek mekanizmadır.** Bir C2 IP'si bugün aktiftir, 3 hafta sonra ölü. Sighting olmadan MISP bunu bilemez → ölü IOC'yi sonsuza dek "aktif" sanıp SIEM'e besler (gürültü). **Sighting + expiration disiplini** ile IOC'ler **yaşayan** veridir: aktif olanlar parlar, ölüler söner. SOC'unu sighting geri-yazacak şekilde otomatikleştir (eşleşme oldu mu → sighting at).

> 🧠 **Expiration/decay:** MISP, taxonomy ile **IOC decay** (zamanla güven azaltma) modelleyebilir — bir IOC'nin "skoru" tipine ve son sighting'ine göre zamanla düşürülür, eşik altına inince to_ids otomatik kapatılabilir. (Decay modeli yapılandırması için resmi MISP decay dokümanından teyit et.)

---

<a id="16"></a>
## 16. 🔥 PÜF NOKTALARI — Piyasada Bulamayacakların

Aşağıdakiler bir CTI analistini "MISP'e veri giren kişi"den "istihbarat üreten kişi"ye ayıran inceliklerdir.

### 16.1 — TLP disiplini ve yanlış paylaşım felaketi
Her zaman önce TLP düşün, sonra distribution ayarla, sonra publish. TLP↔distribution tutarsızlığı en pahalı hatadır (Bölüm 9). Publish öncesi **iki kez** doğrula. Otomasyon kuralın tlp:red/amber'i geniş sync **etmesin**.

### 16.2 — Warning list ile false-positive bastır
İlk kurulumda **tüm temel warning-list'leri etkinleştir** (top-1M, bogon, public resolver, CDN, sinkhole). Bunlar olmadan MISP birkaç ayda gürültüye boğulur, Google DNS'ini bloklarsın (Bölüm 12).

### 16.3 — Feed güvenilirliği / zehirlenmesi
Sadece köklü feed (CIRCL, abuse.ch). Dış feed IOC'lerini `to_ids=false` getir, kritik bloklamayı teyitliyle yap. Önce cache+önizle, sonra import. Feed zehirlenmesi öz-DoS üretebilir (Bölüm 10).

### 16.4 — Korelasyon gürültüsü
Korelasyon kalite ister, miktar değil. Düşük-entropili değerler (yaygın IP, genel filename, CDN) korelasyonu zehirler. Exclusion list + warning-list + over-correlation eşiği kullan. Sinkhole IP'lerine dikkat (Bölüm 8).

### 16.5 — IOC tazeliği (expire/decay)
IOC'ler yaşar ve ölür. Sighting + expiration + decay modeli ile bayat IOC'leri söndür. Ölü C2'yi sonsuza dek bloklamak gürültüdür (Bölüm 15).

### 16.6 — OPSEC: ne paylaşılır, ne paylaşılmaz
Bu, CTI ustalığının çekirdeğidir. **ASLA paylaşılmayacaklar:**
- **Mağdur kimliği / iç altyapı:** Müşteri iç IP'leri, çalışan e-postaları/isimleri, iç hostname'ler — bunlar PII ve ifşa riskidir. Paylaşımdan **önce sanitize et** (maskele/çıkar).
- **Operasyonel gizlilik:** Bir aktörü aktif izliyorsan, IOC'sini PAP:RED'siz paylaşmak **aktörü uyarır** (altyapısını yakar, kaçar). Hedefi uyaracak veriyi PAP/TLP ile koru.
- **Atıf riski:** Kesin olmayan aktör ataması yapma. Yanlış atıf hukuki/itibari sonuç doğurur. Şüphede → yalnızca TTP etiketle, aktör değil.
- **Ham hassas örnek:** Hedefli bir malware örneğini public sandbox'a/feed'e atmak mağduru ifşa eder (bkz. MALWARE rehberi).

### 16.7 — Galaxy ile bağlamlandır
Çıplak IOC veridir; galaxy ile (TTP/aktör/aile) bağlamlanmış IOC istihbarattır. Ama atıf cluster'larını yalnızca güvenle ekle (Bölüm 6). Bağlam Pyramid of Pain'de yukarı çıkarır.

### 16.8 — Event'i over-tagging yapma
30 etiket = gürültü. Az ama doğru: TLP (zorunlu), confidence, kill-chain, gerekirse PAP/aktör. Gereksiz etiket filtrelemeyi ve makinesel işlemeyi bozar (Bölüm 7).

### 16.9 — Üretim kurulum güvenliği
TLS zorunlu, varsayılan parolaları değiştir, MFA aç, worker'ları izle, Diagnostics'i temiz tut. MISP içi istihbarat hassastır → sunucunun kendisi sertleştirilmeli (bkz. LINUX_HARDENING_KALE). İnternete açıksa WAF/VPN arkasına al.

### 16.10 — API key yönetimi
Her entegrasyona ayrı, dar yetkili key. Koda gömme (env/secret). Sızarsa derhal revoke. Read-only iş için read-only key. Sync key admin gücündedir — özellikle koru (Bölüm 14).

### 16.11 — Sighting'in değeri
Sighting'i ihmal etme — IOC'yi yaşayan veri yapan tek şey budur. SOC eşleşmelerini sighting'e çevir; bu hem tazelik hem topluluk doğrulaması sağlar (Bölüm 15).

### 16.12 — STIX uyumluluğu
STIX export/import kayıplı olabilir. Kritik aktarımda test event ile alan kaybını doğrula. MISP↔MISP için yerli JSON en sadık (Bölüm 11).

### 16.13 — Organizasyon hijyeni
External org'ları doğru tanımla (atıf kimliği). UUID'leri tutarlı tut. Aynı kurumu iki org olarak yaratma → korelasyon ve atıf karışır.

### 16.14 — Yedekleme
MISP veritabanı (MySQL) + dosya eklentileri + config düzenli yedeklenmeli. Korelasyon tablosu yeniden üretilebilir ama event verisi değil. Yedeği **şifreli** sakla (VeraCrypt — içinde TLP:RED veri var). Geri-yükleme prosedürünü **test et**.

### 16.15 — Bonus: Confidence ve kaynak izlenebilirliği
Her IOC'nin **nereden geldiği** (kaynak feed/IR/OSINT) ve **ne kadar güvenilir** olduğu (admiralty-scale) izlenebilir olmalı. "Bu IOC'ye neden güveniyoruz?" sorusunun cevabı event'te olmalı — yoksa körlemesine bloklama/atıf yaparsın.

---

<a id="17"></a>
## 17. 🔥 CTI RAPORU İÇİN MISP — Rapordan Event'e, Event'ten Rapora

MISP bir rapor aracı değildir ama bir CTI raporunun **omurgasıdır**: ham gözlemleri yapılandırır, sonra rapor onları anlatır. İki yön de işler.

### Yön 1 — Rapordan MISP'e (raporu event'e bağlama)

Bir tehdit raporu (kendi IR'ından ya da bir vendor raporundan) okudun. İçindeki IOC'leri ve TTP'leri MISP'e **yapılandırılmış** koy:

```
1. Yeni event aç → info: rapor başlığı, date: rapor/olay tarihi
2. Raporun IOC'lerini Freetext Import ile yapıştır → türleri onayla
3. Dosyaları file object'e çevir (hash'ler + boyut)
4. Raporun TTP'lerini ATT&CK galaxy cluster'ı olarak ekle
5. Atıf varsa (ve güvenliyse) threat-actor cluster ekle
6. Rapor PDF'ini event'e EKLE (attachment attribute: 'attachment' type)
   → raporun kendisi ile yapılandırılmış IOC'ler aynı event'te bağlı
7. TLP'yi rapordan devral (rapor TLP:AMBER ise event de AMBER)
8. external-analysis kategorisine rapor linkini (URL) ekle
```

Artık rapor "ölü bir PDF" değil; **korele olabilen, sorgulanabilen, SIEM'e beslenebilen** canlı istihbarat.

### Yön 2 — MISP'ten rapora (IOC tablosu + ATT&CK layer çıkarma)

Tersine: MISP'teki bir event'ten **rapor malzemesi** üret:

- **IOC tablosu çıkar:** Event → Download → CSV (ya da PyMISP `return_format="csv"`). Bunu raporun "Göstergeler" ekine koy. to_ids/category/comment sütunlarıyla temiz bir tablo.

```python
csv_iocs = misp.search(controller="attributes", eventid=EVENT_ID,
                       to_ids=True, return_format="csv")
# rapora gidecek temiz IOC tablosu (ip/domain/url/hash + comment)
```

- **ATT&CK Navigator layer üret:** Event'in ATT&CK galaxy cluster'larından bir **Navigator JSON layer** çıkar (MISP bunu doğrudan export edebilir ya da galaxy tag'lerini toplayıp layer JSON kurarsın). Bu layer, raporun "kullanılan teknikler" görselini verir.

```
MISP event galaxy tag'leri:
   T1566.002, T1204.002, T1071.001
        │
        ▼ (Navigator layer JSON: techniques[] = bu T-numaraları, renkli)
   ATT&CK Navigator katmanı ──► raporun matris görseli
```

> 📌 **ATT&CK Navigator layer'ının DERİN üretimi/yorumu** (renk skorlama, çoklu layer birleştirme, gap analizi) **`MITRE_ATTACK_USTALIK_REHBERI.md`**'de işlenir. Burada sadece "MISP event'inden layer'a köprü" gösterilir.

- **Raporu event'e geri-bağla:** Rapor yayımlandığında, raporun final linkini/sürümünü event'e `link`/`external-analysis` olarak ekle → ileride biri o IOC'yi korele ettiğinde **raporun tamamına** ulaşır. Döngü kapanır.

> 🧠 **CTI raporu = anlatı + MISP = yapı.** Rapor *neden/kim/nasıl/ne yapmalı* anlatır (insan için); MISP *hangi IOC/TTP, nasıl korele, nasıl tespit* taşır (makine için). İkisi birbirine **bağlı** olmalı: rapordaki her IOC MISP'te, MISP'teki event raporun linkinde. Bağlı değilse altı ay sonra ne sen ne SOC'un o IOC'nin hikayesini hatırlar.

---

<a id="18"></a>
## 18. ❌ Yaygın Ölümcül Hatalar

| Hata | Sonuç | Doğrusu |
|---|---|---|
| Yanlış TLP/distribution ile publish | Hassas veri/PII geri-dönüşsüz sızar | Publish öncesi TLP↔distribution iki kez doğrula; kapalı başla |
| Warning-list'i etkinleştirmemek | Google DNS/CDN bloklanır, öz-DoS; korelasyon çöplüğü | İlk kurulumda tüm temel warning-list'leri aç |
| Her şeye to_ids=true | SIEM alarm yağmuru, gerçek tehdit gizlenir | Sadece ayırt edici/bloklanabilir IOC → to_ids=true |
| Dış feed'i körlemesine auto-block | Feed zehirlenmesi → meşru servis kesilir | to_ids=false getir, cache+önizle, köklü feed seç |
| Kesin olmayan aktör ataması | Yanlış atıf, hukuki/itibari hasar | Şüphede yalnızca TTP etiketle, aktör değil |
| Mağdur PII/iç IP paylaşmak | İfşa, sözleşme/yasa ihlali | Paylaşımdan önce sanitize et (maskele/çıkar) |
| Sighting/expiration ihmali | Ölü IOC'ler sonsuza dek "aktif", gürültü | SOC eşleşmesini sighting'e çevir; decay uygula |
| API key'i koda gömmek | Tüm istihbarat çalınabilir | env/secret; dar yetki; sızınca revoke |
| Over-tagging (30 etiket) | Filtreleme/işleme bozulur | Az ama doğru etiket (TLP+confidence+kill-chain) |
| STIX export'u doğrulamamak | Sessiz alan kaybı, eksik istihbarat aktarımı | Test event ile alan kaybını teyit et |
| Aynı kurumu iki org yapmak | Korelasyon/atıf karışır | Org hijyeni; UUID tutarlılığı |
| MISP'i internete sertleştirmeden açmak | İstihbarat ambarı ele geçirilir | TLS+MFA+VPN/WAF; sunucu hardening |
| Yedeksiz çalışmak | Event verisi (üretilemez) kaybolur | Şifreli, test-edilmiş düzenli yedek |
| Korelasyonu yorumsuz "kanıt" saymak | Sahte ilişki (sinkhole/CDN) yanlış sonuç | Düşük-entropili korelasyonu sorgula |

---

<a id="19"></a>
## 19. 🏰 Kanije Kalesi ile Birlikte Kullanım — Forensik Çıktıyı IOC'ye Çevirme

Kanije Kalesi, **kendi cihazını** koruyan ve anti-hırsızlık/forensik veriler üreten bir muhafızdır. Bu veriler doğru ele alındığında **birinci-elden, yüksek-güvenilir IOC kaynağıdır** — çünkü kendi telemetrinden gelir (feed'den değil). Fikir: Kanije'nin ürettiği olayları yapılandırılmış IOC'ye çevirip (uygun TLP ile) MISP'e beslemek.

### Kanije forensik çıktısı → MISP attribute eşlemesi

| Kanije komutu / çıktısı | Üretilen veri | MISP attribute (type) | Not |
|---|---|---|---|
| **`/defender`** (AV durumu + tehdit) | Tehdit adı, imza sürümü, karantina yolu | `text` (threat name), `filename`, dosya `sha256` (varsa) | AV tetiklenen dosyanın hash'i altın IOC |
| **`/erisim`** (erişim/forensik) | Yetkisiz erişim izleri, bağlanan cihazlar, zaman damgaları | `comment`/bağlam; cihaz seri/ID → `text` | Bağlamı event'e taşı |
| **USB algılama** (HID/ses/ağ/telefon) | Bağlanan cihaz türü, ID, zaman | `text` (device id), bağlam | "Rubber-ducky/Bad-USB" şüphesi |
| **`/tuzak`** (honeypot/canary tetiği) | Tuzağa dokunan süreç/işlem | `filename`, `regkey`, süreç adı | **Yüksek değer** — tuzak yalnızca kötücül erişimde tetiklenir |
| **`/koruma`** (lockdown tetiği) | Tetikleyen olay/koşul | Bağlam, zaman | Olay zaman çizelgesi |
| **Kamera/ses hareket tetiği** | Olay zamanı, medya | Bağlam (medya **paylaşılmaz** — PII!) | Sadece olay metadatası, içerik değil |

### Önerilen akış (PyMISP ile)

```python
# Kanije /tuzak (honeypot) tetiklendi: tuzağa dokunan süreç + dosya yakalandı.
# Bunu yapılandırılmış, KAPALI bir MISP event'ine çevir (kendi org'una).
event = MISPEvent()
event.info = "Kanije Honeypot Tetigi — yetkisiz erisim — 2026-05-31"
event.distribution = 0            # SADECE kendi org — bu HASSAS, kendi cihazın
event.threat_level_id = 1
event.add_tag("tlp:red")          # kendi cihaz/forensik → en kapalı başla
event.add_tag('misp-galaxy:mitre-attack-pattern="Valid Accounts - T1078"')  # örnek

# Forensik IOC'ler (kendi telemetrin = yüksek güven, ama atıf YOK)
event.add_attribute("sha256", "<tuzaga dokunan dosya hash>", to_ids=True)
event.add_attribute("text", "USB-HID device id: VID_xxxx PID_yyyy", category="Artifacts dropped",
                    to_ids=False)
event.add_attribute("comment", "/tuzak tetigi; /erisim: yerel oturum, 03:14 TSI", to_ids=False)
misp.add_event(event)
```

> 🔥 **KRİTİK OPSEC — Kanije verisi en hassas veridir:**
> - Kanije çıktısı **senin cihazına/hayatına** dairdir (konum, kamera, oturum, kişisel dosya). Bunu MISP'e koyarken **TLP:RED + Your-org-only** ile başla. **ASLA** kamera/ses/konum içeriğini ya da kişisel veriyi topluluğa paylaşma — yalnızca **soyutlanmış, sanitize IOC** (hash, cihaz türü, teknik gösterge) topluluğa **bilinçli** açılabilir.
> - **Atıf yapma:** "Beni X kişi hackledi" gibi attribution Kanije verisinden **çıkarılamaz** (Bölüm 16.6). Yalnızca teknik gözlem (TTP/IOC) kaydet.
> - Kanije forensik çıktısının ham hali (loglar, medya) **VeraCrypt ile şifreli** saklanmalı (kendi içinde kanıt/PII taşır).

> 💡 **Değer:** Bu entegrasyon, Kanije'yi izole bir "alarm cihazı"ndan **kişisel CTI sensörüne** çevirir: tetiklenen her gerçek olay yapılandırılmış IOC'ye dönüşür, MISP korelasyonu "bu hash/cihaz başka nerede görüldü?" der, ve (sanitize edilirse) topluluğa katkı olur. Ama **sınır nettir:** kendi cihaz verisi varsayılan **kapalı + şifreli + atıfsız**.

---

<a id="20"></a>
## 20. ⚡ Hızlı Referans & Operasyonel Kontrol Listesi

### En sık attribute türleri (hızlı hatırlatma)

| İhtiyaç | type |
|---|---|
| C2/exfil IP | `ip-dst` |
| Saldırgan IP | `ip-src` |
| Kötücül domain | `domain` |
| Phishing URL | `url` |
| Dosya kimliği (tercih) | `sha256` |
| Varyant gruplama | `imphash`, `ssdeep` |
| Phishing gönderen | `email-src` |
| TLS parmak izi | `ja3-fingerprint-md5` |
| Fidye adresi | `btc` |
| Persistence | `regkey`, `windows-service-name` |
| Mutex | `mutex` |
| Zafiyet | `vulnerability` (CVE) |

### TLP hızlı tablo (2.0)

| TLP | Kim | MISP distribution |
|---|---|---|
| RED | İsimli alıcılar | Your org only (0) |
| AMBER(+STRICT) | Org / dar topluluk | 0 ya da Sharing group (4) |
| GREEN | Topluluk (kamuya değil) | Community (1) / Connected (2) |
| CLEAR | Kamuya açık | All communities (3) |

### Event oluşturma kontrol listesi
- [ ] Event **kapalı** başladı (`Your org only`)
- [ ] `info` net, aranabilir; `date` = gözlem tarihi
- [ ] IOC'ler doğru **type + category** ile girildi (Freetext gözden geçirildi)
- [ ] Dosyalar **file object** + hash'ler + ilişkiler (relationship) bağlandı
- [ ] **to_ids disiplini** uygulandı (bloklanabilir→true, bağlam→false)
- [ ] **TLP** etiketlendi; distribution TLP ile **tutarlı**
- [ ] **Galaxy** (ATT&CK TTP) eklendi; aktör **yalnızca güvenliyse**
- [ ] Az ama doğru **taxonomy** (confidence/kill-chain/PAP) — over-tagging yok
- [ ] **Warning-list** kontrolü yapıldı (bogon/CDN/resolver temizlendi)
- [ ] **PII/mağdur/iç altyapı sanitize** edildi
- [ ] Publish öncesi TLP+distribution **iki kez** doğrulandı → **Publish**

### Operasyonel hijyen kontrol listesi
- [ ] **Warning-list'ler** etkin (top-1M, bogon, resolver, CDN, sinkhole)
- [ ] **Feed'ler** köklü kaynaktan; dış feed `to_ids=false`; cache+önizle akışı
- [ ] **Korelasyon** exclusion list + over-correlation eşiği ayarlı
- [ ] **Sighting** otomasyonu kurulu (SOC eşleşmesi → sighting)
- [ ] **Expiration/decay** ile bayat IOC söndürülüyor
- [ ] **API key'ler** dar yetkili, env'de, sızınca revoke prosedürü var
- [ ] **Sync** kuralları TLP/distribution filtreli (tlp:red asla geniş push)
- [ ] **Background worker'lar** çalışıyor (correlation/feed/email)
- [ ] **TLS + MFA + hardening**; internete açıksa VPN/WAF arkasında
- [ ] **Yedek** şifreli + test-edilmiş (event verisi üretilemez)
- [ ] **STIX export** kritik aktarımda doğrulandı

---

<a id="21"></a>
## 21. ⚖️ Hukuki & Etik Sınır + Çapraz Referanslar

> 📛 **Bu bölüm tavsiye değil, sınırdır. Atlama.**

- **Paylaşılan veri = kişisel veri taşıyabilir.** IOC'ler (özellikle `email-dst`, iç IP, hostname) **PII** içerebilir → GDPR/KVKK kapsamındadır. Paylaşmadan **sanitize et**, gereğinden uzun tutma, şifreli sakla.
- **TLP'ye uy — yasal/sözleşmesel yükümlülük.** Sana belirli TLP ile gelen istihbaratı **daha geniş paylaşmak** güven ihlali ve çoğu zaman sözleşme/NDA ihlalidir. TLP:RED/AMBER'i yayma.
- **Atıf (attribution) sorumluluktur.** Bir aktör/ülke ataması **hukuki ve itibari sonuç** doğurur. Kesin kanıt zinciri olmadan atıf yapma — yanlış atıf iftira/itibar davası riskidir. Şüphede TTP gözlemiyle yetin.
- **Mağduru ifşa etmek etik ihlal + olası suç.** Bir ihlal kurbanının kimliğini/verisini rızasız paylaşmak hem etik dışıdır hem veri-koruma ihlalidir.
- **Yetki sınırı.** MISP'i ve içindeki istihbaratı **yalnızca yetkili olduğun kapsamda** (kendi kurumun, üyesi olduğun ISAC, yetkili IR görevi) işle ve paylaş. Kurumunun paylaşım politikasına uy.
- **Malware örneği paylaşımı ayrı kurallara tabidir** (bkz. MALWARE rehberi) — hedefli örneği public'e atmak mağduru yakar ve birçok yargı bölgesinde malware **yaymak** suçtur.
- **Şüphedeysen, paylaşma / daha kapalı TLP seç.** Geri alınamayan tek şey, yayılmış veridir.

---

> 🏰 **Kapanış:** MISP bir veritabanı değil, bir **disiplindir.** En zengin veri modeli, en güçlü korelasyon motoru bile; yanlış TLP ile publish ettiğin, warning-list'i kapalı bıraktığın ya da bir feed'i körlemesine bloklattığın bir akşam **felakettir**. MISP sana istihbaratı **yapılandıracak, koreleyecek ve paylaşacak** omurgayı verir — ama "ne paylaşılır, ne ASLA paylaşılmaz" çizgisini, IOC'nin ne zaman değerli ne zaman gürültü olduğunu, ve bir korelasyonun gerçek mi sinkhole-sahtesi mi olduğunu ayırt etmek **senin işindir.** İki altın kuralı asla unutma: **(1) çıplak IOC veridir — galaxy/sighting/korelasyon ile bağlamlanınca istihbarat olur;** **(2) paylaşım disiplindir — TLP+OPSEC süzgecinden geçmeyen hiçbir şey publish edilmez.** Kanije Kalesi de tam burada: kendi cihazından birinci-elden forensik üretir, ama o veri en hassas olandır — varsayılan kapalı, şifreli, atıfsız. İstihbarat üreten analist karanlıkta avlar; veri biriktiren ise gürültüde boğulur — fark, **disiplindir.**
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. **CTI çekirdek çapraz referansları:** `MITRE_ATTACK_USTALIK_REHBERI.md` (ATT&CK matrisinin DERİN kullanımı, Navigator layer, technique seçimi — MISP galaxy buraya bağlanır), `TTP_AVCILIGI_USTALIK_REHBERI.md` (TTP teorisi, davranışsal avlama, Pyramid of Pain). **İlgili teknik rehberler:** `MALWARE_ANALIZ_USTALIK_REHBERI.md` (IOC/TTP üretimi, örnek paylaşım OPSEC'i, infected-zip), `WIRESHARK_AG_ANALIZ_USTALIK_REHBERI.md` (ağdan IOC çıkarma: IP/domain/JA3), `OSINT_ARAC_SETI_USTALIK_REHBERI.md` (IOC zenginleştirme/dış kaynak), `VERACRYPT_USTALIK_REHBERI.md` (TLP:RED veri ve yedeği şifreli sakla), `GNUPG_GPG_USTALIK_REHBERI.md` (paylaşımda imza/şifreleme, metadata sızıntısı), `LINUX_HARDENING_KALE.md` / `QUBES_OS_BOLMELEME_REHBERI.md` (MISP sunucusunu sertleştirme/izole etme).*
