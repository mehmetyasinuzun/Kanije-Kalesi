# 🕵️ OSINT ARAÇ SETİ — TAM USTALIK REHBERİ (CTI Perspektifi)
## theHarvester'dan Shodan'a, SpiderFoot/Maltego/Recon-ng'e ve OPSEC Püf Noktalarıyla Uçtan Uca

> **Amaç:** OSINT'i (Açık Kaynak İstihbaratı) "bir araca komut yaz, çıktıya bak" seviyesinden çıkarıp, bir **siber tehdit istihbaratı (CTI)** analisti gibi — **disiplinli, izlenemez ve doğrulanmış** — kullanmayı öğretmek. Bu rehber yalnızca *hangi aracı nasıl çalıştırırsın*'ı değil, **neden bu sırayla**, **araştırırken kendini nasıl ele vermezsin** ve **hangi durumda topladığın veri yalan**'ı da anlatır. Forum cevaplarında bulamayacağın OPSEC altyapısı, sock puppet yaşatma, canary token tuzakları, Shodan dork ustalığı, Maltego transform sızıntısı ve attribution kaçınma detayları burada.

> ⚠️ **Önce bunu oku — OSINT'in iki ölümcül yanılgısı:**
> 1. **"Pasif olduğum için izlenemem."** Yanlış. Yanlış kurulmuş bir OSINT oturumu, hedefe *senin onu araştırdığını* sızdırır (LinkedIn profil görüntüleme bildirimi, DNS sorgusu, web bug, doğrudan bağlantı kurma) — ve bir CTI operasyonunda **araştırıldığını fark eden hedef, kanıtı yok eder ya da sana karşı oyun kurar.**
> 2. **"Açık kaynakta bulduysam doğrudur."** Yanlış. Açık kaynak, **dezenformasyonun**, **bayat verinin** ve **yanlış-pozitifin** kaynağıdır. Tek kaynağa dayanan istihbarat, istihbarat değil **dedikodudur.**
>
> Bu iki yanılgı, bu rehberin **OPSEC** (Bölüm 3) ve **Doğrulama** (Bölüm 14) bölümlerinin neden en uzun bölümler olduğunu açıklar. Onları atlama.

> 🚫 **Yasal kapsam (peşinen):** Bu rehber **yalnızca meşru, açık kaynak istihbaratı** anlatır — herkese açık veriyi toplama, ilişkilendirme ve analiz. **İzinsiz erişim, parola kırma, kimlik avı, sömürü (exploit), yetkisiz tarama ya da bir sistemde "içeri girme" bu rehberin kapsamı DIŞINDADIR ve çoğu yargı bölgesinde suçtur.** OSINT'in tüm gücü ve tüm yasallığı, *kapıyı zorlamadan, açık pencereden görüneni* okumaktan gelir. Bu çizgiyi (Bölüm 18) hiç geçme.

---

## 📑 İÇİNDEKİLER

1. [OSINT Nedir, CTI İstihbarat Döngüsündeki Yeri](#1)
2. [Pasif vs Aktif Keşif — Çizgiyi Bilmek](#2)
3. [🔥 OPSEC — EN KRİTİK: Araştırırken Kendini Ele Vermemek](#3)
4. [Sock Puppet / Araştırma Kimlikleri (Oluşturma & Yaşatma)](#4)
5. [Ayrı Altyapı — Tarayıcı, VM, Ağ, E-posta, Telefon](#5)
6. [theHarvester — E-posta / Subdomain / Host Toplama](#6)
7. [Shodan — İnternete Bakan Cihazların Arama Motoru](#7)
8. [SpiderFoot — Otomatik OSINT & Korelasyon](#8)
9. [Maltego — Graph / Entity / Transform Mantığı (ve OPSEC)](#9)
10. [Recon-ng — Modüler Keşif Çatısı](#10)
11. [Tamamlayıcı Teknikler (Dorking, Wayback, crt.sh, WHOIS, EXIF, Sosyal)](#11)
12. [Certificate Transparency & Passive DNS — Pasif Altın Madeni](#12)
13. [Metadata & EXIF — Dosyaların İtirafları](#13)
14. [🔥 Veri Doğrulama / Yanlış-Pozitif / Dezenformasyon Ayıklama](#14)
15. [🔥 PÜF NOKTALARI — Piyasada Bulamayacakların](#15)
16. [☠️ Yaygın Ölümcül Hatalar](#16)
17. [🏰 Kanije Kalesi ile Birlikte Kullanım](#17)
18. [⚖️ Hukuki & Etik Sınır — Yalnızca Açık Kaynak](#18)
19. [✅ Hızlı Referans & OPSEC Kontrol Listesi](#19)

---

<a id="1"></a>
## 1. 🧭 OSINT Nedir, CTI İstihbarat Döngüsündeki Yeri

**OSINT (Open Source Intelligence — Açık Kaynak İstihbaratı)**, *herkese açık ve yasal olarak erişilebilir* kaynaklardan (web siteleri, sosyal medya, sızıntı veritabanları, DNS kayıtları, sertifika logları, devlet sicilleri, haritalar, forumlar, kod depoları) veri toplayıp bunu **istihbarata** dönüştürme disiplinidir. Kritik nokta şudur: ham veri OSINT değildir. OSINT, ham verinin **işlenmiş, doğrulanmış ve bağlama oturtulmuş** hâlidir. "Bu IP adresi şu sertifikayı sunuyor" bir **gözlem**; "Bu altyapı, geçmişte X fidye yazılımı kampanyasında kullanılan C2 sunucularıyla aynı TLS parmak izini ve aynı barındırıcıyı paylaşıyor, dolayısıyla aynı tehdit aktörüne ait olma olasılığı yüksek" bir **istihbarat ürünüdür.**

**CTI (Cyber Threat Intelligence — Siber Tehdit İstihbaratı)**, bu disiplinin güvenlik tarafına uygulanışıdır: tehdit aktörlerini profillemek, saldırı altyapısını haritalamak, IOC'leri (Indicators of Compromise — uzlaşma göstergeleri: IP, alan adı, hash, e-posta) zenginleştirmek ve **kendi kurumunun saldırgan gözünden nasıl göründüğünü** denetlemek (attack surface / saldırı yüzeyi keşfi). OSINT, CTI'nın **toplama (collection)** ayağının en büyük ve en ucuz kaynağıdır.

### İstihbarat Döngüsü — OSINT'in oturduğu çerçeve

OSINT rastgele "Google'da arama" değildir; klasik **istihbarat döngüsünün** içinde yaşar. Her ciddi CTI operasyonu bu beş aşamadan geçer ve döngü **kapalıdır** (yayım, yeni soruları doğurur → tekrar planlama):

```
        ┌──────────────────────────────────────────────────────┐
        │                                                      │
        ▼                                                      │
  ┌───────────┐    ┌───────────┐    ┌───────────┐              │
  │ 1.PLANLAMA│───►│ 2.TOPLAMA │───►│ 3.İŞLEME  │──┐           │
  │  & YÖNLENDİ│    │(COLLECTION)│   │(PROCESSING)│ │          │
  │  RME      │    │  theHarv. │    │ normalize │  │           │
  │  "Neyi    │    │  Shodan   │    │  ayıkla   │  ▼           │
  │  öğrenmem │    │  SpiderFoot│   │  dedup    │ ┌───────────┐ │
  │  lazım?"  │    │  Maltego  │    └───────────┘ │ 4.ANALİZ  │ │
  │  (PIR/EEI)│    │  Recon-ng │                  │(ANALYSIS) │ │
  └───────────┘    └───────────┘                  │ korelasyon│ │
        ▲                                         │ doğrulama │ │
        │                                         │ hipotez   │ │
        │          ┌───────────┐                  └─────┬─────┘ │
        └──────────│ 5.YAYIM   │◄───────────────────────┘       │
   yeni sorular    │(DISSEMIN.)│                                │
                   │ rapor+IOC │────────────────────────────────┘
                   │ imzalı    │   geri besleme
                   └───────────┘
```

| Aşama | Ne yapılır | OSINT karşılığı |
|---|---|---|
| **1. Planlama & Yönlendirme** | İstihbarat gereksinimlerini (PIR — Priority Intelligence Requirement) tanımla. "Tam olarak neyi öğrenmek istiyorum?" | Hedef kapsamı, anahtar kelimeler, sınırlar (scope), OPSEC planı |
| **2. Toplama** | Kaynaklardan ham veri çek | theHarvester, Shodan, SpiderFoot, Maltego, Recon-ng, dorking |
| **3. İşleme** | Ham veriyi normalize et, dedup'la, çevir, ayıkla | CSV/JSON birleştirme, format düzeltme, gürültü atma |
| **4. Analiz** | Veriyi ilişkilendir, doğrula, hipotez kur, sonuç çıkar | Korelasyon, çapraz doğrulama, yanlış-pozitif ayıklama |
| **5. Yayım** | Karar vericiye **işe yarar** istihbaratı ulaştır | İmzalı CTI raporu + yapılandırılmış IOC (STIX/MISP) |

> 🧠 **Altın kural:** En yaygın amatör hatası **doğrudan 2. aşamadan başlamaktır** ("hadi tarayalım"). Profesyonel, **1. aşamada** durur: *Neyi öğrenmem gerekiyor? Hangi soruya cevap arıyorum? Bu araştırmanın OPSEC bütçesi ne?* Planlamasız toplama, gigabaytlık gürültü ve bir de hedefe sızdırılmış izler üretir. Toplama ucuzdur; **analiz ve OPSEC pahalıdır** — orada usta olunur.

> 💡 **PIR vs EEI:** PIR, üst düzey soru ("X tehdit grubu bizim sektörü hedef alıyor mu?"). EEI (Essential Element of Information), onu cevaplayan somut veri parçaları ("X'in bilinen C2 alan adları neler?", "Hangi TLS sertifikalarını kullanıyorlar?"). Her toplama adımını bir EEI'ye bağla — bağlanmıyorsa, muhtemelen gereksiz gürültü ve gereksiz OPSEC riski topluyorsundur.

---

<a id="2"></a>
## 2. 🎯 Pasif vs Aktif Keşif — Çizgiyi Bilmek

OSINT'in OPSEC kalbi **bu ayrımdadır.** "Pasif" ve "aktif" arasındaki çizgiyi yanlış çizmek, hem seni ifşa eder hem de (aktif tarafta) yasal sınırı geçirebilir.

### Tanım

- **Pasif keşif:** Veriyi **hedefin altyapısına hiç dokunmadan**, üçüncü taraf kaynaklardan toplarsın. Hedefin log'unda **senin izin yoktur.** Örnek: crt.sh'tan sertifika geçmişi çekmek, Shodan'ın *önceden topladığı* banner'ı okumak, Wayback Machine'den eski sayfaları görmek, WHOIS sorgusu (bu da bir "ortada" sorgudur ama hedef sunucuya değil tescil sağlayıcısına gider), Google'da arama.
- **Aktif keşif:** Hedefin sistemine **doğrudan bir paket/istek gönderirsin** → hedefin log'unda **senin (ya da kullandığın altyapının) izi kalır.** Örnek: hedef alan adına `nmap` taraması, hedef web sitesini *senin* tarayıcınla ziyaret etmek, subdomain'leri *canlı* DNS ile brute-force etmek, hedef e-postaya doğrulama (verification) isteği atmak.

```
   SEN ──► [3. taraf kaynak: Shodan/crt.sh/Wayback/Google] ──► veri      ✅ PASİF
           (hedef bunu GÖRMEZ — iz üçüncü tarafta kalır)

   SEN ──────────────────────────────────────────────► [HEDEF sunucu]   ⚠️ AKTİF
           (hedefin erişim log'unda IP'in/parmak izin var)
```

### Neden hayati?

| Senaryo | Pasif | Aktif |
|---|---|---|
| Hedef seni fark eder mi? | **Hayır** (genelde) | **Evet** — log'unda IP/UA/zamanın var |
| Gerçek zamanlı mı? | Hayır (veri "bayat" olabilir, Bölüm 15) | Evet (anlık gerçek durum) |
| Yasal risk | Düşük (açık kaynak) | Yetkisizse **yüksek** (izinsiz tarama suç olabilir) |
| Erken aşamada tercih | ✅ **Önce hep pasif** | Yalnızca yetki + gerekçe varsa |

> 🔥 **Disiplin kuralı:** **Tüm pasif toplama bitmeden aktif keşfe geçme.** Çoğu hedef hakkında bilmen gerekenin %80'ini pasif yollarla, hedef hiç fark etmeden öğrenebilirsin. Aktife geçmek bir karardır — gerekçesi, yetkisi ve OPSEC altyapısı (Bölüm 5) olmalı. Bir CTI angajmanında "yanlışlıkla hedefin sitesini ziyaret etmek" bile, dikkatli bir saldırgan altyapısında **senin varlığını ele verir** ve operasyonu yakar.

> ⚠️ **Sinsi tuzak — "yarı pasif" araçlar:** Bazı araçlar pasif sanılır ama bir modülü **sessizce aktif** çalışır. SpiderFoot'un birçok modülü pasiftir ama bazıları hedef sunucuya bağlanır (Bölüm 8). theHarvester'ın çoğu kaynağı pasiftir ama DNS brute-force modülü **aktiftir.** Aracın her modülünün hedefe dokunup dokunmadığını **bilerek** kullan — "pasif araç" diye bir şey yoktur, yalnızca "pasif kullanılan modül" vardır.

---

<a id="3"></a>
## 3. 🔥 OPSEC — EN KRİTİK: Araştırırken Kendini Ele Vermemek

> Bu, rehberin **en önemli bölümüdür** ve bilerek en başa, en uzun şekilde kondu. Çünkü bir CTI/OSINT operasyonunda **topladığın veri ikincildir; o veriyi toplarken kendini ele verip vermediğin birincildir.** Hedefe "seni araştırıyorum" sinyali sızdıran bir analist, hem operasyonu yakar hem de (bir tehdit aktörü söz konusuysa) kendini hedef hâline getirir.

**OPSEC (Operations Security)**, askeri kökenli bir disiplindir: *düşmanın senin niyetini/yeteneğini/kimliğini çıkarmasına yarayacak ipuçlarını sistematik olarak yok etmek.* OSINT'te düşman = araştırdığın hedef (ya da onu izleyen taraf). İpuçları = senin IP'in, tarayıcı parmak izin, zaman dilimin, dil ayarın, hesap izlerin, davranış kalıbın.

### OPSEC'in altın kuralı: **Attribution kaçınma (ilişkilendirilememe)**

Amaç, topladığın hiçbir izin **gerçek kimliğine, kurumuna ya da niyetine** geri bağlanamamasıdır. Bu üç katmanda kırılır:

1. **Ağ katmanı:** IP adresin seni ve coğrafyanı ele verir.
2. **Kimlik katmanı:** Kullandığın hesaplar (giriş yapmış LinkedIn/Twitter/Google) seni doğrudan adlandırır.
3. **Davranış katmanı:** Zaman dilimin (mesai saatlerin), dilin, yazım kalıbın, arama deseninin tutarlılığı seni **dolaylı** ele verir (örtülü parmak izi).

### 🔥 "Hedefe araştırdığını sızdırmama" — somut sızıntı kanalları

Bunlar, amatörlerin sürekli düştüğü ve operasyon yakan **gerçek** kanallardır:

#### 3.1 LinkedIn "profil görüntüleme" bildirimi
LinkedIn, biri profiline baktığında karşı tarafa **"profilinizi kim görüntüledi"** bildirir. Bir kurumun çalışanlarını araştırırken **kişisel (ya da kuruma bağlanabilen) hesabınla** profillere bakarsan, hedefin İK'sı/güvenlik ekibi *kimin onları incelediğini* görür → operasyon yanar.
- **Çözüm:** LinkedIn'i ya **anonim/özel mod**'a al (Settings → Visibility → Profile viewing options → "Private mode") — ama bu özelliklerin **çoğu özel mod kullanırken karşı tarafın kim olduğunu görme yeteneğini de kapatır ve premium gerektirebilir; resmi LinkedIn ayarlarından güncel davranışı teyit et.** Ya da hiç giriş yapma (out-of-network görünüm sınırlıdır). En temizi: araştırmayı LinkedIn'in *kendi* arama/bildirim mekanizmasına dokunmadan, **Google cache / üçüncü taraf** üzerinden yap.

#### 3.2 Doğrudan bağlantı / etkileşim kurma
Hedefe arkadaşlık isteği atmak, gönderisini beğenmek, takip etmek, mesaj atmak — hepsi **sock puppet'ini bile** ifşa eder ve hedefe "biri benimle ilgileniyor" sinyali verir. **Asla hedefle etkileşime girme** (yalnızca pasif gözlem). Etkileşim, sosyal mühendisliktir ve bu rehberin (yasal OSINT) kapsamı dışındadır.

#### 3.3 DNS / referer / CDN sızıntısı
- **DNS sızıntısı:** VPN/Tor kullansan bile, sistemin DNS sorguları **VPN dışından** (ISP'nin DNS'i) çıkarsa, ISP — ve dolaylı olarak izleyen taraf — hangi alan adlarını çözdüğünü görür. Üstelik *senin DNS çözücün* hedefin authoritative DNS'ine ulaşırsa iz bırakır.
  - **Çözüm:** DNS leak testi yap; VPN'in DNS'ini zorla; Tor üzerinden çalış (Tor DNS'i devre üzerinden çözer).
- **Referer sızıntısı:** Bir sayfadan bir linke tıklarsan, tarayıcın `Referer` başlığında **nereden geldiğini** hedef sunucuya söyleyebilir → iç bir araç/portal URL'in ya da arama sorgun sızar.
- **CDN/3. taraf isteği:** Bir hedef sayfada gömülü kaynaklar (analytics, font, beacon) senin tarayıcından **hedefin kontrol ettiği** sunuculara istek attırabilir → seni hedefin log'una sokar (bu zaten "aktif" tarafa kayar).

#### 3.4 Web bug / canary token / tracking pixel tuzakları
Dikkatli bir hedef (ya da seni bekleyen bir tehdit aktörü), **kasten yem belge/bağlantı** koyar:
- **Web bug / tracking pixel:** Bir belgeye/sayfaya gömülü 1x1 görünmez görsel; sen açınca *onların* sunucusuna istek gider → **IP'ini, zamanını, kullanıcı-ajanını** verir.
- **Canary token:** Özel hazırlanmış bir dosya/link/QR/AWS anahtarı; biri ona dokunduğunda token sahibine **anında alarm** gider ("biri yemi yuttu"). Tehdit avcıları bunu savunma için kullanır; ama hedefin de senin için kurmuş olabilir.
- **Çözüm:** Şüpheli belgeleri **asla kendi makinende, kendi IP'inle açma.** İzole VM + Tor/VPN + (gerekirse) belgeyi internetsiz açan bir görüntüleyici kullan. Office/PDF makrolarını ve uzak içerik yüklemeyi **kapat.** (Bkz. Bölüm 5, `QUBES_OS_BOLMELEME_REHBERI.md` disposable VM mantığı.)

> 🍯 **Simetriyi gör:** Kanije Kalesi'nin `/tuzak` (honeypot) komutu **tam olarak bu silahı sana verir** — *senin* sistemine yaklaşan saldırganı yakalamak için. Yani sen hedefe canary kurabildiğin gibi, **hedef de sana kurabilir.** OSINT yaparken her zaman "bu açtığım şey bir tuzak mı?" diye düşün.

#### 3.5 Zaman dilimi & dil parmak izi (örtülü attribution)
- **Zaman deseni:** Sock puppet hesabın hep **senin yerel mesai saatlerinde** aktifse, gözlemci zaman dilimini → coğrafyanı → muhtemel kimliğini daraltır. Aktör profili "Türkiye saatiyle 09:00-18:00 arası aktif" diye not düşülür.
- **Dil/yerel ayar:** Tarayıcı `Accept-Language: tr-TR` başlığı, işletim sistemi dili, sock puppet hesabının yazım dili, klavye düzeni izleri — hepsi coğrafi/kültürel ipucu sızdırır.
- **Çözüm:** Sock puppet'in zaman dilimini ve dilini **kapağına** (cover story) uygun ayarla, kendi gerçek desenine değil. Aktiviteyi farklı saatlere yay. Tarayıcı dilini/zaman dilimini kapakla tutarlı yap (ama parmak izi tutarlılığına dikkat — Bölüm 5).

#### 3.6 Tarayıcı parmak izi (browser fingerprinting)
IP'ini gizlesen bile, tarayıcının **benzersiz parmak izi** (User-Agent, ekran çözünürlüğü, yüklü fontlar, canvas/WebGL hash'i, eklentiler, saat dilimi) seni oturumlar arası **tekilleştirebilir.** Aynı parmak izi hem sock puppet'inde hem (bir gün dikkatsizce) kişisel oturumunda görünürse → ikisi ilişkilendirilir.
- **Çözüm:** Her kimliğe **ayrı tarayıcı profili/VM** (Bölüm 5). Anti-detect tarayıcılar (Mullvad Browser, Tor Browser) parmak izini *standartlaştırarak* kalabalığa karıştırır. Eklenti bağımlı, kişiselleştirilmiş tarayıcı **kötü** OPSEC'tir (benzersizleşir).

### OPSEC altyapısı — katmanlı diyagram

```
┌─────────────────────────────────────────────────────────────────┐
│                       SEN (gerçek kimlik)                       │
│                            │                                    │
│   ASLA karışmaz ───────────┼──────────── ASLA karışmaz          │
│                            ▼                                    │
│              ┌──────────────────────────┐                       │
│              │   İZOLE ARAŞTIRMA ORTAMI │  (ayrı VM / ayrı OS)   │
│              │  ┌────────────────────┐  │                       │
│              │  │ Ayrı tarayıcı prof.│  │  ← parmak izi ayrı     │
│              │  │ Sock puppet hesap. │  │  ← kişiselle karışmaz  │
│              │  └────────────────────┘  │                       │
│              └────────────┬─────────────┘                       │
│                           ▼                                    │
│              ┌──────────────────────────┐                       │
│              │   ANONİMLEŞTİRME KATMANI │                       │
│              │   VPN ya da Tor          │  ← gerçek IP gizli     │
│              │   + DNS leak koruması    │                       │
│              └────────────┬─────────────┘                       │
│                           ▼                                    │
│   ════════════════════════════════════════════════════════     │
│                    AÇIK İNTERNET / HEDEF KAYNAKLAR              │
│   (yalnızca PASİF; aktif yalnızca yetki + gerekçe + bu altyapı)│
└─────────────────────────────────────────────────────────────────┘
```

> 🧠 **OPSEC felsefesi:** OPSEC bir araç değil, bir **alışkanlıktır.** En iyi VPN bile, yanlışlıkla kişisel Google hesabınla giriş yaptığın bir sekmede çöker. Kural basit ve acımasızdır: **araştırma kimliği ile kişisel kimlik asla, hiçbir koşulda, tek bir sekmede/uygulamada/cihazda buluşmaz.** Buna "duvar" (the wall) denir ve bir kez delindiğinde geri alınamaz.

---

<a id="4"></a>
## 4. 🎭 Sock Puppet / Araştırma Kimlikleri (Oluşturma & Yaşatma)

**Sock puppet** ("kukla hesap"), OSINT araştırmacısının kullandığı, gerçek kimliğine bağlanamayan **sahte ama inandırıcı** bir çevrimiçi kimliktir. Sosyal medya OSINT'inde, kapalı grupları görmek, profil incelemek ve kişisel hesabını korumak için zorunludur.

> ⚖️ **Etik/yasal sınır önce:** Sock puppet, **gözlem ve kişisel kimliği koruma** içindir — **sosyal mühendislik, kandırma yoluyla bilgi sızdırma, sahte kimlikle dolandırıcılık DEĞİLDİR.** Çoğu platformun şartlarına aykırıdır (hesap kapatılabilir) ve hedefle *etkileşime* girerek bilgi koparmak yasal/etik çizgiyi aşar. Bu rehberde sock puppet = **pasif gözlem için bir paravan**, aktif bir aldatma aracı değil.

### Oluşturma — inandırıcı bir kimlik kurmak

1. **Kapak hikayesi (cover/legend) önce:** Önce *kim olduğuna* karar ver — yaş, şehir, meslek, ilgi alanları. Tutarlı ve **sıkıcı/sıradan** olsun (dikkat çekmesin). "İstanbul'da yaşayan, kedi seven, futbol tutan bir muhasebeci" — kontrol edilebilir ama ilgi çekici olmayan.
2. **Yüz/avatar:** Gerçek bir kişinin fotoğrafını **asla** kullanma (ters görsel aramayla ifşa olur + kişiyi mağdur eder). **GAN ile üretilmiş yapay yüz** (ör. "thispersondoesnotexist" tarzı üreticiler) kullan — ama bunların da **bilinen artefaktları** var (kulak/gözlük/arka plan tutarsızlıkları) ve bazı tespit araçları yakalayabilir; bir yüzü **birden çok kimlikte kullanma.**
3. **İsim:** Bölgeye uygun, yaygın, sıradan bir ad. Gerçek ünlü/var olan biriyle çakışmasın.
4. **E-posta + telefon:** Kişiselinden **tamamen ayrı** (Bölüm 5). Hesap kurtarma e-postası/telefonu da ayrı olmalı — yoksa kurtarma zinciri seni gerçek kimliğine bağlar.
5. **Yavaş büyüt:** Yeni, boş, hemen "araştırma yapan" bir hesap **şüphelidir** ve platform anti-spam'ine takılır.

### 🔥 Sock puppet'i YAŞATMA sanatı (en çok burada ölür)

Hesap *oluşturmak* kolay; **yaşatmak** zordur. Çoğu sock puppet, oluşturulduktan sonra "soğuk" kaldığı, tutarsız davrandığı ya da kişiselle karıştığı için yanar.

- **Yaşlandır (aging):** Bir kimliği kuracaksan **kullanmadan haftalar/aylar önce kur** ve organik şekilde "yaşat" — birkaç gönderi, birkaç takip (alâkasız, masum hesaplar), profili doldur. Olgun, geçmişi olan bir hesap inandırıcıdır.
- **Tutarlı kişilik:** Yazım tarzı, ilgi alanları, paylaşım saatleri hep **kapakla** uyumlu olsun (kendi gerçek desenine değil — Bölüm 3.5). Bir gün "muhasebeci kedi sever", ertesi gün "sızma testi uzmanı" gibi davranma.
- **Karşılıklı ayrım (compartmentalization):** **Her sock puppet ayrı tarayıcı profili/VM** ve **mümkünse ayrı kimliğe bağlı IP**. İki puppet'i aynı oturumda kullanırsan, platform onları (ve seni) **aynı parmak izi/çerez/IP** üzerinden ilişkilendirir → "bu iki hesap aynı kişi" → ikisi birden yanar.
- **Asla kişiselle çapraz kullanma:** Sock puppet sekmesinde kişisel Gmail'e girme, kişisel telefonla SMS doğrulama yapma, kişisel ödeme kartı kullanma. **Tek çapraz temas, tüm puppet ağını gerçek kimliğine bağlar.**
- **Para/abonelik izi:** Bir puppet'e premium (LinkedIn/Twitter) gerekiyorsa, ödemeyi kişisel kartınla yapma → ön ödemeli/sanal kart veya hiç ödeme. Ödeme izi en güçlü deanonimleştiricilerden biridir.
- **Telefon doğrulama tuzağı:** Çoğu platform telefon ister. Kişisel numaranı **asla** verme. Ayrı, kimliğe bağlanmayan bir numara gerekir (Bölüm 5) — ama "sanal numara" servislerinin birçoğu platformlarca **kara listededir** ve hesabı baştan yakar; bu, sock puppet kurmanın en kırılgan adımıdır.

> 🧠 **Pratik:** İhtiyacın olan kadar puppet kur, fazlasını değil — her biri **bakım** ister (yaşatmak, tutarlılık, ayrı altyapı). Bakamayacağın 10 ölü puppet'ten, iyi yaşatılmış 2 puppet daha değerlidir. Ve unutma: bir puppet bir kez yanarsa (platform kapatırsa ya da ifşa olursa) onu **kurtarmaya çalışma** — gömüp yenisini kur, yoksa yanan kimlik diğerlerini de ele verir.

---

<a id="5"></a>
## 5. 🧱 Ayrı Altyapı — Tarayıcı, VM, Ağ, E-posta, Telefon

Sock puppet ve OPSEC, ancak **altyapı gerçekten ayrıksa** işe yarar. "Duvarı" fiziksel/teknik olarak kuran katman budur.

### 5.1 Ayrı tarayıcı profili / VM / OS

| Seviye | Ne | İzolasyon gücü |
|---|---|---|
| **Ayrı tarayıcı profili** | Aynı tarayıcıda farklı "profile" | Zayıf — çerez ayrı ama parmak izi/IP aynı, sızıntı kolay |
| **Ayrı tarayıcı (Tor/Mullvad)** | Farklı uygulama + standart parmak izi | Orta — parmak izi kalabalığa karışır |
| **Ayrı VM** | Araştırmaya özel sanal makine | İyi — OS düzeyinde izole, snapshot ile temizlenir |
| **Ayrı OS / Tails / Qubes** | Amnezik ya da bölmeli işletim sistemi | **En güçlü** — donanım/ağ/kimlik tam ayrı |

- **VM kullan:** Araştırmayı **disposable (tek kullanımlık) VM**'de yap; iş bitince snapshot'a dön ya da yok et → canary/web bug bulaşsa bile temizlenir. (Bkz. `QUBES_OS_BOLMELEME_REHBERI.md` — her kimliğe ayrı "qube".)
- **Tails:** Amnezik; kapatınca iz bırakmaz, tüm trafik Tor'dan çıkar → yüksek riskli pasif OSINT için ideal. (Bkz. `TAILS_ANONIMLIK_REHBERI.md`.)

### 5.2 Ağ — VPN vs Tor

| | VPN | Tor |
|---|---|---|
| Gizlediği | Gerçek IP (sağlayıcıya güven gerekir) | Gerçek IP (dağıtık, güven gerektirmez) |
| Hız | Hızlı | Yavaş |
| Çıkış IP itibarı | Temiz (datacenter ise bazen bloklanır) | Çıkış düğümleri sıkça **bloklu/CAPTCHA'lı** |
| Coğrafya seçimi | Kolay (ülke seç) | Sınırlı (çıkış düğümü ülkesi) |
| Log riski | Sağlayıcı log tutabilir → **no-log + yargı bölgesi** önemli | Tek nokta log tutamaz |
| **OSINT'te** | Çoğu pasif iş için pratik; kapakla uyumlu ülke seç | En hassas, attribution-kritik işlerde |

- **DNS leak koruması şart** (Bölüm 3.3): VPN/Tor kullansan bile DNS sızarsa altyapı çöker. Bağlandıktan sonra **mutlaka DNS leak testi** yap.
- **Çıkış coğrafyası ↔ kapak tutarlılığı:** Sock puppet'in "İstanbul'da muhasebeci" ise ama IP'in Hollanda datacenter'ından çıkıyorsa tutarsızlık doğar. Kapakla uyumlu çıkış konumu seç (ama datacenter IP'leri "VPN" diye işaretlenebilir — residential proxy ihtiyacı doğabilir; bu da maliyet/etik dengesi gerektirir).

### 5.3 Ayrı e-posta & telefon

- **E-posta:** Her kimlik için ayrı, gizlilik dostu sağlayıcı; kurtarma e-postası da **ayrı** olmalı. Kurtarma zinciri (recovery chain) en sık gözden kaçan deanonimleştiricidir — kişisel bir e-posta/telefonu kurtarma olarak eklersen, tüm puppet'i ona bağlamış olursun.
- **Telefon:** Kişisel numara **asla.** Ayrı, kimliğe bağlanmayan numara gerekir; ama (Bölüm 4) sanal numaralar platformlarca sıkça reddedilir — bu adım kırılgandır, planla.
- **Ödeme:** Premium gereken yerde kişisel kart kullanma → sanal/ön ödemeli kart ya da hiç. **Ödeme = en güçlü kimlik bağı.**

> 🔥 **Püf — "temiz oda" prensibi:** Araştırma VM'ini bir **temiz oda** gibi düşün: içeri hiçbir kişisel şey girmez (kişisel hesap, dosya, kart, çerez), dışarı hiçbir araştırma izi sızmaz (gerçek IP, DNS, parmak izi). VM ile ana sistem arasında **kes-yapıştır/dosya paylaşımı bile** bir köprüdür — kapat. Snapshot disiplini: temiz bir baseline'dan başla, iş bitince **baseline'a geri dön** (web bug/canary/çerez kalıntısı silinir).

---

<a id="6"></a>
## 6. 🌾 theHarvester — E-posta / Subdomain / Host Toplama

**theHarvester**, bir alan adı/kurum hakkında **e-posta adresleri, subdomain'ler, host'lar, çalışan adları, IP'ler ve açık portları** *çok sayıda açık kaynaktan* toplayan klasik bir keşif aracıdır. Bir CTI angajmanının **erken pasif toplama** ayağının iş atıdır.

### Ne işe yarar
Bir hedefin "saldırı yüzeyini" hızlıca çıkarır: hangi e-posta adresleri herkese açık dolaşımda (kimlik avı/sızıntı riski göstergesi), hangi subdomain'ler var (unutulmuş/test sunucuları), hangi host'lar internete bakıyor. CTI tarafında: hem **kendi** kurumunu denetlemek (saldırgan ne görüyor?) hem **tehdit aktörü** altyapısını haritalamak için.

### Kurulum
```bash
# Kali/çoğu pentest dağıtımında kurulu gelir
theHarvester -h

# Manuel (Python):
pip install theHarvester
# ya da kaynaktan:
git clone https://github.com/laramies/theHarvester
cd theHarvester
python3 -m pip install -r requirements.txt
```
> ⚠️ **Sürüm farkı:** Bayrak adları sürümle değişebilir (`-b` kaynak adları, modül listesi). `theHarvester -h` ile **kendi sürümünün** desteklediği kaynakları ve bayrakları teyit et.

### Temel kullanım
```bash
# hedef alan adı için, belirli bir kaynaktan
theHarvester -d ornek.com -b bing

# tüm (varsayılan) kaynaklardan, sonuç limiti ile
theHarvester -d ornek.com -l 500 -b all

# çıktıyı dosyaya (HTML/JSON/XML — sürüme göre)
theHarvester -d ornek.com -b all -f cikti_ornek
```

| Bayrak | Anlamı |
|---|---|
| `-d` | Hedef alan adı/kurum |
| `-b` | Kaynak (source): `bing`, `duckduckgo`, `crtsh`, `hackertarget`, `otx`, `shodan` vb. ya da `all` |
| `-l` | Sonuç limiti |
| `-f` | Çıktı dosyası (rapor) |
| `-s` | Başlangıç ofseti (sayfalama) |

### İleri kullanım — kaynak modülleri & API anahtarları
theHarvester'ın gücü **kaynak çeşitliliğindedir.** Bazı kaynaklar **anahtarsız** çalışır (bing, duckduckgo, crtsh, hackertarget); ama **en zengin** kaynaklar (Shodan, SecurityTrails, Hunter, Censys, GitHub, VirusTotal/OTX vb.) **API anahtarı** ister. Anahtarları yapılandırma dosyasına eklersin:
```bash
# Sürüme göre konum değişir — TEYİT ET:
#   ~/.theHarvester/api-keys.yaml   (yaygın yeni konum)
#   ya da kurulum dizinindeki api-keys.yaml
# Örnek (yaml):
#   shodan: { key: "ANAHTARIN" }
#   securityTrails: { key: "ANAHTARIN" }
```

> 🔥 **theHarvester'a özel püf — kaynak çakıştırma & DNS brute-force ayrımı:**
> - **Çok kaynak çalıştır, sonra birleştir:** Tek kaynak eksik döner; `bing` farklı, `crtsh` farklı, `otx` farklı subdomain bulur. `-b all` çalıştır ama sonra **çıktıyı dedup'la** (aynı subdomain birçok kaynaktan gelir → çakışma = doğrulama, ama listeyi şişirir).
> - **DNS brute-force MODÜLÜ AKTİFTİR:** theHarvester'ın subdomain *brute-force* / DNS *resolution* özelliği (sürüme göre `-c`, `-n` gibi bayraklar) hedefin DNS'ine **canlı sorgu** gönderir → **pasif değil, aktiftir** (Bölüm 2). Saf pasif kalmak istiyorsan **yalnızca pasif kaynakları** seç (bing/crtsh/otx) ve brute-force/resolve adımlarını **kapalı tut.** İz bırakmadan keşif yapıyorsan bu ayrım kritik.
> - **crt.sh kaynağı = bedava altın:** `-b crtsh` (Certificate Transparency, Bölüm 12) tamamen pasiftir ve genelde en zengin subdomain listesini verir. Anahtar bile gerekmez.

---

<a id="7"></a>
## 7. 🛰️ Shodan — İnternete Bakan Cihazların Arama Motoru

**Shodan**, "internete bağlı her şeyin Google'ı"dır. Web sayfalarını değil, **servis banner'larını** indeksler: bir IP'nin hangi portlarında ne çalışıyor, hangi yazılım/sürüm, hangi TLS sertifikası, hangi başlık (HTTP title), hangi ülke/kurum. Sunucular, kameralar, yönlendiriciler, ICS/SCADA, veritabanları, IoT — internete bakan ne varsa.

> 🧠 **Neden pasif?** Shodan **kendi tarayıcılarıyla** interneti sürekli tarar ve sonucu depolar. **Sen Shodan'da arama yaptığında hedefe HİÇBİR paket gitmez** — Shodan'ın *önceden topladığı* veriyi okursun. Bu yüzden Shodan, OSINT'in en güçlü **pasif** araçlarından biridir (ama veri **bayat** olabilir — Bölüm 15).

### Banner grabbing nedir
Bir servise bağlanıp ondan dönen "tanıtım" metnini (banner) okumak: `Server: nginx/1.18.0`, SSH sürümü, FTP karşılama mesajı, sertifika CN'i. Shodan bunu senin yerine *zaten* yapmıştır.

### Kurulum & CLI
```bash
# CLI (Python)
pip install shodan
shodan init <API_ANAHTARIN>      # hesabındaki anahtarla kur

# Bir IP hakkında her şey
shodan host 1.2.3.4

# Arama (web arayüzündeki sorgunun aynısı)
shodan search 'apache country:"TR"'

# Sonuç sayısı (kota harcamadan keşif)
shodan count 'product:"MongoDB"'

# Tarama kredisi durumu
shodan info
```

### 🔥 Arama filtreleri / dork'ları — Shodan'ın asıl gücü
Shodan ustalığı **filtre dilindedir.** Doğru filtreyle, milyonlarca cihaz arasından tam aradığını bulursun:

| Filtre | Ne yapar | Örnek |
|---|---|---|
| `net:` | IP/CIDR bloğu | `net:198.51.100.0/24` |
| `org:` | Kuruluş adı (WHOIS) | `org:"Example Corp"` |
| `hostname:` | Host adı içeren | `hostname:.ornek.com` |
| `ssl:` | TLS sertifikasında geçen | `ssl:"ornek.com"` |
| `ssl.cert.subject.cn:` | Sertifika CN'i | `ssl.cert.subject.cn:"*.ornek.com"` |
| `http.title:` | Sayfa başlığı | `http.title:"Admin Login"` |
| `http.html:` | HTML gövdesinde geçen | `http.html:"powered by"` |
| `port:` | Açık port | `port:3389` (RDP) |
| `product:` | Ürün/servis | `product:"MySQL"` |
| `country:` / `city:` | Coğrafya | `country:"TR" city:"Istanbul"` |
| `asn:` | Otonom sistem no | `asn:AS12345` |
| `os:` | İşletim sistemi | `os:"Windows"` |
| `vuln:` | Bilinen zafiyet (CVE) — **akademik/kurumsal plan gerekir** | `vuln:CVE-2021-44228` |

```bash
# Bir kuruma ait, RDP açık, Türkiye'deki sistemler (kendi yüzeyini denetleme)
shodan search 'org:"Example Corp" port:3389 country:"TR"'

# Bir alan adının TÜM sertifika-eşleşen sistemleri (subdomain/altyapı keşfi)
shodan search 'ssl:"ornek.com"'

# Aynı self-signed sertifikayı paylaşan altyapı (C2 korelasyonu — CTI)
shodan search 'ssl.cert.serial:"<seri_no>"'
```

> 🔥 **Shodan dork ustalığı (en güçlü püf):**
> - **`ssl:` ile altyapı pivotu:** Bir tehdit aktörünün bir sunucusunu bulduysan, o sunucunun **TLS sertifikasının** belirgin bir alanını (CN, organizasyon, **serial**, ya da **JARM/JA3S** parmak izi) al → `ssl.cert...` veya `ssl.jarm:` ile aynı parmak izini paylaşan **diğer sunucularını** bul. Saldırganlar altyapıyı kopyala-yapıştır kurar; sertifika/JARM **aynı kalır** → tek sunucudan tüm kampanya altyapısı **pasif** olarak ortaya çıkar. Bu, modern CTI'nın en güçlü pivot tekniğidir.
> - **`http.favicon.hash:` ile teknoloji eşleştirme:** Belirli bir panel/C2 framework'ünün favicon hash'ini bilirsen, internetteki **tüm örneklerini** tek sorguyla bulursun. (Favicon hash pivotu, Cobalt Strike vb. araçların altyapısını avlamada meşhurdur.)
> - **`count` ile kota koru:** Keşif yaparken `shodan search` yerine `shodan count` kullan — kaç sonuç olduğunu kredi harcamadan görürsün; sorguyu daraltıp daraltmadığını ucuza test edersin.
> - **`vuln:` herkese açık değil:** Zafiyet filtresi genelde **akademik/iş planı** gerektirir; ücretsiz/temel planda çalışmaz. Resmi plan sayfandan teyit et.

### Monitoring (sürekli izleme)
Shodan, IP/aralığını **sürekli izleyip** yeni açık port/servis/zafiyet çıkınca **alarm** verebilir (Shodan Monitor). Kendi kurumunun saldırı yüzeyini "yanlışlıkla açılan port" için 7/24 izlemenin pratik yoludur.

### Shodan vs Censys (hangisi?)
| | Shodan | Censys |
|---|---|---|
| Odak | Geniş cihaz/IoT/banner kapsamı, olgun dork dili | Derin sertifika/host veri modeli, akademik köken |
| Sertifika analizi | İyi | **Çok güçlü** (CT entegre, zengin sorgu) |
| Arama dili | Filtre/dork tabanlı | Yapılandırılmış sorgu dili (alan bazlı) |
| CTI kullanımı | Hızlı pivot, favicon/JARM | Derin sertifika korelasyonu, geçmiş |
> 💡 **Pratik:** İkisini **birlikte** kullan — Shodan'da hızlı bul, Censys'te sertifika/host ilişkisini **derinleştir ve çapraz doğrula.** Tek kaynağa güvenme (Bölüm 14).

---

<a id="8"></a>
## 8. 🕷️ SpiderFoot — Otomatik OSINT & Korelasyon

**SpiderFoot**, bir hedef (alan adı, IP, e-posta, isim, telefon, Bitcoin adresi, ad-soyad) verdiğinde **onlarca açık kaynağı otomatik tarayan, sonuçları birbirine bağlayan (korelasyon) ve bir veriyi diğerine pivotlayan** bir OSINT otomasyon motorudur. Elle 50 araç çalıştırmak yerine, SpiderFoot bunu **modüller** üzerinden orkestrasyonla yapar.

### Ne işe yarar
Tek bir veri parçasından (örn. bir e-posta) başlayıp; o e-postanın geçtiği sızıntıları, ilişkili alan adlarını, sosyal hesapları, IP'leri, subdomain'leri, sertifikaları **otomatik** keşfeder ve hepsini bir ilişki grafiğine/raporuna döker. CTI'da bir IOC'yi **zenginleştirmenin** (enrichment) hızlı yoludur.

### Kurulum & çalıştırma
```bash
# Python
git clone https://github.com/smicallef/spiderfoot
cd spiderfoot
pip3 install -r requirements.txt

# Web arayüzünü başlat (yerel)
python3 ./sf.py -l 127.0.0.1:5001
# Tarayıcıdan http://127.0.0.1:5001

# Komut satırından tek tarama
python3 ./sf.py -s ornek.com -t DOMAIN_NAME
```
> ⚠️ Arayüzü **yalnızca localhost'a** bağla (`127.0.0.1`) — açık IP'ye bağlarsan tüm OSINT sonuçların ve API anahtarların ağa açılır.

### Modüller & API entegrasyonları
SpiderFoot'un gücü **modüllerindedir** (200+). Çoğu anahtarsız çalışır; en zenginleri (Shodan, VirusTotal, Hunter, SecurityTrails, HaveIBeenPwned, Censys vb.) **API anahtarı** ister → Settings'ten ekle. Tarama başlatırken bir **use case** seçersin:
- **Footprint** — hedefin internet ayak izini geniş çıkar.
- **Investigate** — kötü amaçlılık göstergeleri (kara liste, malware ilişkisi).
- **Passive** — **hedefe hiç dokunmadan** yalnızca pasif kaynaklar.

### 🔥 SpiderFoot'a özel püf — PASİF MOD ve "hedefe dokunma" tehlikesi
> Bu, bu aracın **en kritik** OPSEC noktasıdır (ve en çok atlananı):
> - **SpiderFoot her modülü pasif DEĞİLDİR.** Bazı modüller hedefin web sitesini **canlı tarar** (spidering), DNS'ini **çözer**, portlarını yoklar → bu **aktif keşiftir** (Bölüm 2) ve hedefin log'una **senin altyapının izini** bırakır. "Otomatik OSINT" rahatlığı, farkında olmadan hedefe dokunmana yol açar.
> - **Çözüm:** Hedefe iz bırakmaması gereken bir CTI işinde **mutlaka "Passive" use case** seç ya da tarama ayarlarında **yalnızca pasif modülleri** etkinleştir; web spider / DNS brute / port-scan modüllerini **kapat.** Tarama öncesi modül listesini **gözden geçir** — "bu modül hedefe paket gönderir mi?" diye her birini değerlendir.
> - **Korelasyon = asıl değer:** SpiderFoot'un raporundaki tekil bulgular ham veridir; asıl değer **korelasyon kurallarında** (aynı IP'yi paylaşan varlıklar, bir e-postanın birden çok yerde geçmesi). Raporu "liste" gibi değil "graf" gibi oku — hangi düğümler birbirine bağlanıyor?
> - **API kotası yönetimi:** "Footprint/Investigate" tüm modülleri ateşler → API anahtarlarının kotasını **bir taramada** tüketebilir (özellikle Shodan/VT). Önce dar kapsamlı/pasif tara, kotayı bilinçli harca (Bölüm 15).

---

<a id="9"></a>
## 9. 🕸️ Maltego — Graph / Entity / Transform Mantığı (ve OPSEC)

**Maltego**, OSINT'in **görsel ilişki analizi** (link analysis) standardıdır. Veriyi tablo değil **graf** olarak düşünür: düğümler (entity) ve onları birbirine bağlayan ilişkiler. Bir tehdit aktörünü, bir altyapıyı, bir kişiyi araştırırken "kim neyle bağlantılı" sorusunu **görsel** olarak çözer.

### Üç temel kavram
1. **Entity (varlık):** Bir veri parçası — alan adı, IP, e-posta, kişi, telefon, şirket, sosyal hesap, hash, BTC adresi. Graf'taki bir düğüm.
2. **Transform (dönüşüm):** Bir entity'yi alıp **ilişkili yeni entity'ler üreten** bir işlem. Örn. "alan adı" entity'sine "DNS'ten IP'ye çöz" transform'u uygula → bağlı IP entity'leri belirir. "E-posta" → "geçtiği sızıntılar" transform'u → ilişkili veriler. Transform'lar zincirlenir → graf büyür.
3. **Graph (graf):** Tüm entity'ler + ilişkiler. Pivot yaparak (bir düğümden komşularına geçerek) bir veri parçasından koca bir ağa ulaşırsın.

```
   [ornek.com]──(DNS çöz)──►[1.2.3.4]──(reverse DNS)──►[mx.ornek.com]
        │                       │
   (WHOIS)                  (Shodan)
        ▼                       ▼
   [admin@ornek.com]        [port 443 / sertifika]──(ssl pivot)──►[diger-altyapi.com]
        │
   (sızıntı kontrolü)
        ▼
   [parola sızıntısı kaydı]
```

### CE vs ticari & Hub
- **Maltego CE (Community Edition):** Ücretsiz; transform sayısı ve sonuç limiti **kısıtlı**, ticari kullanım sınırlı. Öğrenmek ve küçük araştırmalar için yeterli.
- **Ticari (Pro/Enterprise):** Sınırsıza yakın, kurumsal veri kaynakları, daha büyük graf.
- **Transform Hub:** Üçüncü taraf veri sağlayıcıların (Shodan, VirusTotal, HIBP, vb.) transform'larını eklediğin pazar yeri. Çoğu **kendi API anahtarını** ister.

### 🔥 Maltego'ya özel püf — TRANSFORM OPSEC SIZINTISI (en kritik)
> Bu, Maltego'nun **en tehlikeli ve en az bilinen** OPSEC tuzağıdır:
> - **Transform'lar genelde KİMİN sunucusunda çalışır?** Çoğu Maltego transform'u **senin makinende değil, sağlayıcının (Paterva/Maltego ya da üçüncü taraf) sunucusunda** çalışır. Yani bir entity'ye transform uyguladığında, **araştırdığın hedefin adını/IP'sini/e-postasını o üçüncü taraf sunucuya GÖNDERİRSİN.** Sonuç: araştırma niyetin ve hedefin, **senin kontrolünde olmayan** bir sunucuya sızar. Hassas bir CTI operasyonunda bu, "kim kimi araştırıyor" bilgisini dış bir tarafa teslim etmektir.
> - **"Local transform" mı "remote transform" mı?** Bazı transform'lar yerel (senin makinende) çalışır → hedef verisi dışarı gitmez. Hangi transform'un **local** hangisinin **remote** olduğunu **bilerek** seç. Çok hassas işlerde yalnızca local transform / kendi barındırdığın TDS (Transform Distribution Server) kullan.
> - **API anahtarın = kimliğin:** Transform Hub'a eklediğin Shodan/VT anahtarın **senin hesabına bağlıdır** → o sağlayıcı, "bu anahtar sahibi şu hedefi araştırdı" kaydını tutar. Anahtar hijyeni (Bölüm 15) burada da geçerli: araştırma için ayrı hesap/anahtar.
> - **Otomatik pivot graf'ı patlatır:** "Tüm transform'ları çalıştır" baştan çekici ama (a) API kotanı tüketir, (b) hedef verisini onlarca dış sunucuya saçar, (c) gürültüyle graf okunmaz olur. **Tek tek, gerekçeli** transform uygula — her transform bir EEI'ye (Bölüm 1) hizmet etsin.

> 🧠 **Özet OPSEC:** Maltego'da "bir butona basıp graf büyütmek" eğlenceli görünür; ama **her transform, hedef verisini bir yere gönderen bir ağ işlemidir.** Pasif sandığın bir tıklama, hedef verisini üçüncü tarafa (ve dolaylı olarak belki hedefe) sızdırabilir. **Önce transform'un nerede çalıştığını öğren, sonra tıkla.**

---

<a id="10"></a>
## 10. 🧰 Recon-ng — Modüler Keşif Çatısı

**Recon-ng**, "web tabanlı keşfin Metasploit'i" olarak tasarlanmış, **modüler, komut satırı tabanlı** bir OSINT çatısıdır (framework). theHarvester gibi tek iş yapmaz; **workspace** (çalışma alanı), **veritabanı**, **modül marketplace** ve **raporlama** ile uçtan uca bir keşif iş akışı sunar.

### Neden çatı?
theHarvester bir *araç*, Recon-ng bir *ortamdır*: bulguları (host, contact, credential-leak göstergesi, port) merkezi bir **SQLite veritabanında** biriktirir; bir modülün çıktısı diğerinin girdisi olur (bir alan adından subdomain'leri çıkar → onları başka modüle besle → IP'leri çöz → raporla). Tekrarlanabilir, ölçeklenebilir araştırma için tasarlanmıştır.

### Kurulum & temel kullanım
```bash
# Kali'de kurulu; manuel:
git clone https://github.com/lanmaster53/recon-ng
cd recon-ng
pip install -r REQUIREMENTS
./recon-ng

# Konsol içinde:
[recon-ng] > workspaces create ornek_arastirma   # izole çalışma alanı
[recon-ng] > marketplace search                   # mevcut modüller
[recon-ng] > marketplace install all              # (ya da seçili modülleri)
[recon-ng] > modules load recon/domains-hosts/hackertarget
[recon-ng] > options set SOURCE ornek.com
[recon-ng] > run
[recon-ng] > show hosts                            # toplanan host'lar
```

### Workspaces (çalışma alanları)
Her hedef/angajman için **ayrı workspace** aç → veri karışmaz, raporlar ayrı, OPSEC olarak da hijyenik (bir müşterinin verisi diğerine sızmaz). Bu, Recon-ng'in en güçlü organizasyon özelliğidir.

### Marketplace & modüller
- `marketplace search` / `marketplace info <modül>` / `marketplace install <modül>`.
- Modül kategorileri: `recon/` (toplama), `discovery/`, `import/`, `report/`.
- Bazı modüller **API anahtarı** ister (`keys add shodan_api <anahtar>`, `keys list`).

### API key yönetimi
```bash
[recon-ng] > keys add shodan_api <ANAHTARIN>
[recon-ng] > keys add virustotal_api <ANAHTARIN>
[recon-ng] > keys list
```
Anahtarlar workspace'ten bağımsız, merkezi tutulur → araştırma anahtarlarını **kişisel anahtarlarından ayrı** bir kurulumda/kullanıcıda tut (Bölüm 15).

### Raporlama
```bash
[recon-ng] > modules load reporting/html
[recon-ng] > options set FILENAME /yol/rapor.html
[recon-ng] > run
# CSV / JSON / XML / Markdown / list de mevcut (modüle göre)
```

> 🔥 **Recon-ng'e özel püf:**
> - **Aktif modül ayrımı:** `recon/` modüllerinin çoğu pasif (3. taraf API) ama bazıları (DNS brute-force, `discovery/` altındakiler, web içerik tarama) **hedefe dokunur → aktiftir** (Bölüm 2). Modülü yüklemeden `marketplace info` / `info` ile **ne yaptığını** oku; pasif kalman gerekiyorsa aktif modülleri çalıştırma.
> - **Workspace = OPSEC sınırı:** Her angajmana ayrı workspace; bittiğinde **veritabanını güvenli sakla/sil** (hedef hakkında topladığın hassas veri orada birikir — bir saldırı yüzeyidir; Bölüm 17 → `/imha`).
> - **`db` ile veriyi sorgula:** Topladığın her şey SQLite'ta; `db query SELECT ...` ile (sürüme göre) doğrudan sorgulayıp dışa aktarabilirsin → analiz aşamasında (Bölüm 1) işe yarar.
> - **Marketplace'i güncel tut ama körlemesine `install all` yapma:** Bazı modüller ölü API'lara bağlıdır (çalışmaz/hata verir). İhtiyacın olanı kur; her modül bir bağımlılık ve bazen bir API anahtarı demektir.

---

<a id="11"></a>
## 11. 🔎 Tamamlayıcı Teknikler (Dorking, Wayback, crt.sh, WHOIS, EXIF, Sosyal)

Araç setinin yanında, hiçbir kuruluma ihtiyaç duymayan ama **en çok iş gören** elle teknikler:

### 11.1 Google / GitHub Dorking
**Dork**, arama motorunun gelişmiş operatörleriyle **hassas/gizli kalmış** içeriği bulmaktır.
```text
# Google dork örnekleri (tamamen pasif — Google'ın index'i)
site:ornek.com -www                       # subdomain/sayfa keşfi
site:ornek.com filetype:pdf               # açık PDF'ler (metadata için!)
site:ornek.com intitle:"index of"         # açık dizin listeleri
site:ornek.com inurl:admin                # yönetim panelleri
"ornek.com" filetype:xls (password | parola)   # sızmış tablolar
```
```text
# GitHub dork — kod/secret sızıntısı (kendi kurumunu denetle!)
"ornek.com" password
org:ExampleCorp filename:.env
"api_key" "ornek.com"
AWS_SECRET_ACCESS_KEY  (bir kurum adıyla birlikte)
```
> 🔥 **Püf:** GitHub dorking, **kendi geliştiricilerinin yanlışlıkla commit'lediği sırları** (API anahtarı, `.env`, parola) bulmanın en hızlı yoludur — saldırgandan **önce sen bul.** Hedef bir tehdit aktörüyse, onun açık deposundaki altyapı ipuçlarını ortaya çıkarır. GitHub'ın **kendi secret-scanning** uyarılarını da kullan. (Dorking sonucu bulduğun *başkasının* sırrını **kullanmak** ≠ OSINT — Bölüm 18.)

### 11.2 Wayback Machine (Internet Archive) & arşivler
`web.archive.org`, web sayfalarının **geçmiş hâllerini** saklar. Bir hedefin **kaldırdığı** sayfa, eski e-posta/telefon, eski teknoloji yığını, silinmiş çalışan listesi orada **hâlâ durur.**
```bash
# Bir alan adının arşivlenmiş TÜM URL'lerini çek (CDX API — pasif)
curl "http://web.archive.org/cdx/search/cdx?url=ornek.com*&output=text&fl=original&collapse=urlkey"
```
> 🔥 **Püf:** Bir hedef "izini temizlemiş" olabilir ama Wayback (ve `archive.today`, Google cache) **eski hâli tutar.** Silinmiş bir LinkedIn gönderisi, kaldırılmış bir `robots.txt` (gizli dizinleri ele verir!), eski bir `/about` sayfası → araştırmanın altın madeni. **Pasiftir** (Archive'ın kopyasını okursun, hedefe gitmezsin).

### 11.3 crt.sh & Certificate Transparency
Bölüm 12'de derinleşiyor. Özet: `crt.sh`, TLS sertifikası loglarını sorgular → bir alan adının **tüm subdomain'lerini** (sertifika alındığı için ifşa olur) **pasif** olarak verir.
```bash
# Bir alan adının sertifikalarındaki tüm isimleri çek (JSON)
curl "https://crt.sh/?q=%25.ornek.com&output=json"
```

### 11.4 WHOIS & Passive DNS
- **WHOIS:** Alan adı tescil bilgisi (kayıt tarihi, tescil eden, name server'lar, bazen iletişim — ama çoğu artık **GDPR/privacy proxy** ile gizli).
```bash
whois ornek.com
```
- **Passive DNS (pDNS):** Bir alan adının/IP'nin **geçmişteki** DNS çözümlemelerini gösterir (üçüncü taraf toplayıcılardan — pasiftir). "Bu IP geçmişte hangi alan adlarını barındırdı?" → altyapı pivotunun temeli (Bölüm 12).

### 11.5 EXIF / Metadata
Bölüm 13'te derinleşiyor. Özet: Fotoğraf/belge dosyaları, **GPS konumu, cihaz, yazar, yazılım, zaman damgası** gibi gömülü metadata taşır.

### 11.6 Sosyal Medya OSINT (SOCMINT)
- Kişi/kurum profilleri, gönderiler, takipçi ağı, etiketlenen konumlar, paylaşım saatleri (yaşam deseni), bağlantılar.
- **Araçlar:** Platforma özel arama operatörleri, profil arşivleyiciler, kullanıcı-adı çapraz arama (aynı handle'ı birçok platformda arama).
> ⚠️ **OPSEC kritik:** SOCMINT, OPSEC'in **en çok sınandığı** yerdir — sürekli giriş yapmış hesap, profil görüntüleme bildirimleri (Bölüm 3.1), yanlışlıkla beğeni/takip. **Daima sock puppet + ayrı VM** (Bölüm 4-5). **Asla kişisel hesabınla** hedef profillere dalma.

---

<a id="12"></a>
## 12. 📜 Certificate Transparency & Passive DNS — Pasif Altın Madeni

Bu iki teknik, **en güçlü pasif keşif** kaynaklarıdır ve ayrı bir bölümü hak eder — çünkü hedefe **tek paket göndermeden** koca bir altyapıyı haritalarlar.

### Certificate Transparency (CT) — sertifikaların herkese açık günlüğü
Modern TLS ekosistemi, her verilen sertifikanın **herkese açık, değiştirilemez günlüklere** (CT logs) yazılmasını zorunlu kılar (tarayıcılar bunu denetler). Sonuç: bir kurum `vpn.ornek.com`, `test.ornek.com`, `internal-mail.ornek.com` için sertifika aldığında, **bu adlar CT loglarında ifşa olur** — kurum o subdomain'i hiç yayınlamasa bile.

```bash
# crt.sh — tüm sertifika kayıtlarındaki isimler (pasif)
curl -s "https://crt.sh/?q=%25.ornek.com&output=json" | \
  jq -r '.[].name_value' | sort -u
```
> 🔥 **Püf — CT ile "gizli" subdomain avı:** DNS brute-force (aktif, gürültülü, eksik) yerine **CT logları** (pasif, sessiz, kapsamlı) çoğu zaman daha fazla subdomain verir — üstelik hedef **hiç fark etmez.** `test.`, `dev.`, `staging.`, `vpn.`, `jira.`, `git.` gibi *unutulmuş/iç* sistemler genelde CT'de görünür. Bir CTI angajmanında **ilk** bakacağın yer burasıdır. **Wildcard sertifika** (`*.ornek.com`) tek tek subdomain'leri gizleyebilir → o zaman CT yetmez, pDNS + pasif DNS dataset'leri devreye girer.

### Passive DNS (pDNS) — DNS'in geçmişi
Normal DNS sorgusu **şu anki** cevabı verir (ve hedefin DNS'ine dokunur = aktif). **Passive DNS**, üçüncü taraf sensörlerin **zamanla topladığı** DNS cevaplarının arşividir → **tamamen pasif** ve **tarihsel.**

Ne sorarsın:
- "Bu IP geçmişte hangi alan adlarını barındırdı?" (reverse — IP'den alan adlarına)
- "Bu alan adı geçmişte hangi IP'lere çözüldü?" (forward — altyapı göçü/geçmiş)
- "Bu name server'ı kullanan başka hangi alan adları var?"

> 🔥 **Püf — pDNS ile altyapı pivotu (CTI'nın kalbi):** Bir tehdit aktörünün bir alan adını bulduysan, pDNS ile o alan adının **çözüldüğü IP'leri** bul → o IP'lere pDNS reverse uygula → **aynı IP'yi paylaşan diğer kötü amaçlı alan adlarını** ortaya çıkar. Saldırganlar altyapıyı paylaşır/yeniden kullanır; pDNS bu paylaşımı **geçmişe dönük** görünür kılar. Shodan'ın `ssl:` pivotu (Bölüm 7) + crt.sh (CT) + pDNS = **üç bacaklı pasif altyapı haritalama**, bir aktöre tek paket göndermeden tüm ağını çıkarır.

> ⚠️ **Doğrulama notu:** pDNS verisi **bayat** olabilir (eski kayıt) ve **paylaşımlı barındırma** (shared hosting/CDN) yanlış-pozitif üretir — bir IP'de yüzlerce alakasız site olabilir. pDNS pivotunu **her zaman** ikinci bir kaynakla çapraz doğrula (Bölüm 14): CDN/shared IP mi, dedicated mı?

---

<a id="13"></a>
## 13. 🖼️ Metadata & EXIF — Dosyaların İtirafları

Dosyalar konuşur. Bir kurumun yayınladığı PDF, Word, Excel ya da fotoğraf, **gömülü metadata** taşır ve bu metadata sık sık **istemeden** sızar.

### Ne sızar
| Dosya türü | Olası metadata |
|---|---|
| **Fotoğraf (JPEG)** | **GPS koordinatı**, kamera marka/model, çekim tarihi-saati, yazılım |
| **PDF** | Yazar adı, oluşturan yazılım/sürüm, şablon yolu (iç kullanıcı adı!), oluşturma/değiştirme tarihi |
| **Office (docx/xlsx)** | Yazar, son düzenleyen, kurum, şablon, **iç dosya yolları**, revizyon geçmişi |
| **Tümü** | İç ağ yolları (`\\sunucu\paylaşım\...`), kullanıcı adları, yazılım envanteri |

### Araçlar
```bash
# ExifTool — metadata çıkarmanın standardı
exiftool fotograf.jpg
exiftool -gps:all -a fotograf.jpg          # yalnızca GPS
exiftool *.pdf | grep -i author            # toplu yazar çıkar

# Toplu site metadata toplama (PDF/doc indir + metadata çıkar)
# (örn. metagoofil / FOCA tarzı araçlar — kurumun açık belgelerini tarar)
```

> 🔥 **Püf — kurumsal metadata madenciliği:** Bir hedefin sitesinden **tüm açık belgeleri** (Google dork: `site:ornek.com filetype:pdf OR filetype:docx`) indir → `exiftool` ile **yazar adlarını, iç kullanıcı adlarını, yazılım sürümlerini, iç ağ yollarını** çıkar. Bu, kurumun **iç kullanıcı adı şemasını** (örn. `ad.soyad`), **yazılım envanterini** (eski/zafiyetli sürümler) ve hatta **iç sunucu adlarını** ifşa eder — hepsi **pasif** (belge zaten herkese açık). Tehdit aktörü tarafında: aktörün yayınladığı bir belgenin metadata'sı (yazar, dil, yazılım) **attribution** ipucudur (Bölüm 14).

> 🛡️ **Savunma yüzü (kendini denetle):** Aynı tekniği **kendi kurumuna** uygula — yayınladığın PDF/Office dosyaları iç kullanıcı adı/ağ yolu sızdırıyor mu? Yayınlamadan önce metadata **temizle** (`exiftool -all= dosya`). Bu, Kanije Kalesi felsefesindeki "saldırgan ne görüyor?" denetiminin (Bölüm 17) bir parçasıdır.

> ⚠️ **Tuzak hatırlatması:** İndirdiğin belge bir **canary token** (Bölüm 3.4) olabilir — açtığında hedefe haber gider. Belgeleri **izole VM'de, internetsiz, makro kapalı** aç. `exiftool` dosyayı *çalıştırmadan* okur (görece güvenli) ama yine de izole ortamda çalış.

---

<a id="14"></a>
## 14. 🔥 Veri Doğrulama / Yanlış-Pozitif / Dezenformasyon Ayıklama

> OSINT'in **ikinci en kritik** bölümü (OPSEC'ten sonra). Çünkü topladığın veri, **doğrulanana kadar istihbarat değildir** — sadece iddiadır. Yanlış doğrulanmış OSINT, yanlış kararlar (ve bir CTI bağlamında yanlış attribution → yanlış aktöre suçlama → diplomatik/hukuki felaket) üretir.

### Neden açık kaynak özellikle kirli?
- **Dezenformasyon:** Tehdit aktörleri **kasten** sahte iz bırakır (false flag): başka bir ülkenin diline/araçlarına ait artefakt yerleştirir, sahte "sızıntı" yayar, yanlış altyapıya yönlendirir.
- **Bayat veri:** Shodan banner'ı 3 ay önceki olabilir; pDNS kaydı eski olabilir; arşiv silinmiş bir gerçeği gösterebilir (Bölüm 15).
- **Yanlış-pozitif:** Shared hosting/CDN bir IP'de yüzlerce alakasız siteyi buluşturur; aynı favicon binlerce sitede olabilir; bir e-posta deseni tesadüfen eşleşebilir.
- **Bağlam kaybı:** Doğru veri, yanlış bağlamda yanlış sonuç verir.

### Doğrulama disiplini — somut kurallar

#### 14.1 Çok kaynak çapraz doğrulama (en temel kural)
**Tek kaynak = sıfır kaynak.** Bir bulguyu en az **iki bağımsız** kaynaktan teyit et. "Shodan diyor ki port 22 açık" → Censys de aynı şeyi diyor mu? Bağımsız demek: *aynı verinin iki kopyası değil* (Shodan ve Shodan'ı kaynak alan bir araç bağımsız değildir), gerçekten **farklı toplama yöntemi.**

#### 14.2 Birincil kaynağa in
Bir iddia bir blogdan geliyorsa, o blog nereden almış? **Zincirin ucundaki birincil kaynağa** kadar in. "Herkes böyle diyor" bir şeyi doğru yapmaz — herkes aynı tek (belki yanlış) kaynaktan kopyalamış olabilir (circular reporting / döngüsel raporlama tuzağı).

#### 14.3 Yanlış-pozitif elemesi
- **IP pivotu yaptıysan:** O IP **shared/CDN mi**? (PTR kaydı, ASN, aynı IP'deki site sayısı). Shared ise "aynı IP = aynı sahip" **yanlıştır.**
- **Sertifika/favicon pivotu:** Bu sertifika/favicon **jenerik mi** (varsayılan panel, popüler CMS)? Jenerikse eşleşme tesadüfidir.
- **İsim/e-posta eşleşmesi:** Yaygın bir ad mı? Tesadüf olabilir.

#### 14.4 Dezenformasyon / false flag farkındalığı
- Bir "kanıt" **fazla mı temiz/fazla mı uygun**? (Aktörün dilini, saat dilimini, aracını *tam da beklediğin gibi* gösteren artefakt → şüphelen; yerleştirilmiş olabilir.)
- Metadata (Bölüm 13) **çelişiyor mu**? (Belge "Rusça" ama derleme zaman dilimi başka, dil paketi başka → kasıtlı yanıltma olabilir.)
- **Atıf güveni dereceli olmalı:** Asla "kesin X aktörü" deme; **"orta/yüksek güvenle X ile tutarlı"** de. CTI'da güven seviyesi (low/medium/high confidence) ve **alternatif hipotez** (rakip hipotez analizi — ACH) standarttır.

#### 14.5 Zamansal doğrulama (verinin tazeliği)
Her bulgunun **zaman damgasını** kaydet: "Bu banner ne zaman toplandı?" "Bu DNS kaydı hangi tarihe ait?" Eski veriyle güncel sonuç çıkarma. (Bölüm 15 → "bayatlama".)

> 🧠 **Analist zihniyeti:** İyi bir OSINT analisti **doğrulamak için değil, çürütmek için** uğraşır (falsification). "Bu hipotez doğruysa, başka ne görmem gerekir? Görüyor muyum? Bu hipotezi *yanlışlayacak* ne ararım?" Onaylama yanlılığı (confirmation bias), OSINT'in bir numaralı analiz katilidir — beklediğini görürsün, görmek istediğini doğrularsın. **Hipotezini öldürmeye çalış; öldüremezsen, o zaman güçlüdür.**

### Arşivleme & zaman damgası (kanıt bütünlüğü)
Doğruladığın her bulguyu **kanıt olarak sabitle:**
- **Arşivle:** Sayfayı `archive.today`'e gönder / `web.archive.org`'a kaydettir / yerel olarak tam sayfa + ekran görüntüsü al. (Hedef sayfayı silerse kanıtın kalır.)
- **Zaman damgala:** Bulguyu **ne zaman** ve **nereden** topladığını kaydet. Kanıt değerini artırmak için topladığın dosyaları **GPG ile imzala / hash'le** → "bu IOC dosyası şu tarihte bizdeydi ve değişmedi" (Bkz. `GNUPG_GPG_USTALIK_REHBERI.md` — detached signature; RFC 3161 TSA ile bağımsız zaman damgası).
- **Zincir bütünlüğü:** Kaynak → toplama yöntemi → tarih → analist → sonuç zincirini belgele. Mahkemeye/karar vericiye giden istihbaratın **izlenebilir** olması şart.

---

<a id="15"></a>
## 15. 🔥 PÜF NOKTALARI — Piyasada Bulamayacakların

Çoğu rehberin atladığı, gerçek dünyada bir OSINT/CTI operasyonunun başarısını **ve araştırmacının güvenliğini** belirleyen detaylar.

### 15.1 API key hijyeni — anahtarların kimliğindir
Shodan, VirusTotal, SecurityTrails, Hunter, Censys anahtarların **senin hesabına bağlıdır.** Bir araca (theHarvester, SpiderFoot, Recon-ng, Maltego) eklediğin anahtar, o sağlayıcıya **"bu hesap şu hedefleri araştırdı"** kaydı bıraktırır. Üstelik anahtarın config dosyasında **düz metin** durur → makinen ele geçerse anahtarların (ve dolaylı olarak araştırma geçmişin) sızar. **Kurallar:** (a) araştırma için **kişiselden ayrı** hesaplar/anahtarlar kullan; (b) anahtarları **kod deposuna asla commit'leme** (en sık sızıntı — GitHub dork'la kendi anahtarını arat!); (c) anahtarları şifreli sakla (Bkz. `KEEPASSXC_PAROLA_KALESI_REHBERI.md`); (d) sızma şüphesinde **derhal döndür (rotate).**

### 15.2 Rate-limit yönetimi — sessiz kal, kotayı koru
Her API'nin **istek sınırı** (rate limit) ve **kotası** vardır. İki sebepten kritik: (a) sınırı aşarsan **ban/throttle** yersin ya da kotanı bir taramada tüketirsin (özellikle SpiderFoot/Maltego "tüm modülleri çalıştır" derken); (b) bazı **aktif** kaynaklarda hızlı/agresif sorgu **hedefin de dikkatini çeker** (anomali tespiti). **Kurallar:** `shodan count` ile keşfet (kredi harcamadan), sorguyu **daralt** sonra çalıştır; otomasyonlarda gecikme/throttle ayarla; toplu işleri **gece/dağıtık** değil **kapakla tutarlı** saatlere yay (Bölüm 3.5); ücretsiz kotayı bilinçli harca.

### 15.3 Sock puppet'i yaşatma sanatı (operasyonel özet)
Bölüm 4'ün özü, püf olarak: **kimlik kurmak kolay, yaşatmak zordur.** Bir puppet'i **kullanmadan önce yaşlandır** (haftalar/aylar), **tutarlı kişilik** taşı (kapakla uyumlu, kendi deseninle değil), **her birine ayrı altyapı** ver, ve **asla kişiselle çapraz kullanma.** Bir puppet yanarsa **kurtarma**, **göm.** Bakabileceğinden fazla puppet kurma — her biri bakım borcudur.

### 15.4 Pasif → aktif geçmeme disiplini
Bölüm 2'nin özü: **tüm pasif toplama bitmeden aktife geçme.** Aktife geçmek bir karardır (yetki + gerekçe + altyapı). "Yarı pasif" araçların (theHarvester DNS brute, SpiderFoot spider, Recon-ng discovery, hedef sitesini tarayıcıyla ziyaret) sessizce hedefe dokunduğunu **bil.** Bir CTI işinde hedefe dokunmak, dikkatli bir altyapıda **seni ele verir** ve operasyonu yakar.

### 15.5 Shodan dork ustalığı (operasyonel özet)
Bölüm 7'nin özü: Shodan'ın gücü **filtre/pivot dilindedir.** `ssl:` / `ssl.cert.serial:` / `ssl.jarm:` / `http.favicon.hash:` ile **tek sunucudan tüm altyapıya** pasif pivot yap. `count` ile kotayı koru. `vuln:` ücretli plandadır — varsayma, **teyit et.** Censys ile çapraz doğrula. **Banner bayat olabilir** (15.10) — Shodan'ın "last seen" zamanına bak.

### 15.6 Maltego transform OPSEC sızıntısı (operasyonel özet)
Bölüm 9'un özü: **çoğu transform hedef verisini ÜÇÜNCÜ TARAF sunucuya gönderir** (remote transform). Pasif sandığın bir tıklama, "kim kimi araştırıyor" bilgisini dışarı sızdırır. **Local vs remote** transform ayrımını bilerek seç; çok hassas işte yalnızca local/kendi TDS'in. "Tüm transform'ları çalıştır"ma — gerekçeli, tek tek uygula.

### 15.7 Otomasyonun (SpiderFoot/Recon-ng aktif modülleri) hedefe dokunması
Bölüm 8 & 10'un özü: **"otomatik OSINT" rahatlığı en büyük OPSEC tuzağıdır.** Bir butona basıp "her şeyi tara" dediğinde, araç senin adına **hedefe paket gönderen** modülleri de ateşler (web spider, DNS brute, port-scan). Pasif kalman gerekiyorsa **modül listesini önceden denetle** ve aktif olanları **kapat.** Otomasyonu körlemesine güvenme — her modülün hedefe dokunup dokunmadığını bil.

### 15.8 Canary token / honeypot tuzakları (savunma ve farkındalık)
Bölüm 3.4'ün özü: hedef (ya da seni bekleyen aktör) **kasten yem** koymuş olabilir — açtığında IP'ini/zamanını veren web bug, dokununca alarm veren canary token, balküpü dosya. **Şüpheli her şeyi izole VM + Tor/VPN + makro kapalı** aç; uzak içerik yüklemeyi kapat. **Simetri:** Kanije Kalesi'nin `/tuzak`'ı *senin* için bu silahı sana verir — yani sen de yakalanabilirsin. "Bu açtığım şey bir tuzak mı?" sorusunu **her zaman** sor.

### 15.9 Attribution kaçınma (operasyonel özet)
Bölüm 3'ün kalbi: topladığın hiçbir iz **gerçek kimliğine/kurumuna/niyetine** bağlanamamalı. Üç katman — **ağ** (IP → VPN/Tor + DNS leak koruması), **kimlik** (hesap → sock puppet, asla kişisel), **davranış** (zaman/dil/parmak izi → kapakla tutarlı, kişiselle değil). **Tek çapraz temas** (bir sekmede kişisel Gmail) tüm zinciri çözer. "Duvar" bir kez delinir.

### 15.10 Verinin "bayatlaması" (data staleness)
OSINT verisi **zaman içinde çürür.** Shodan banner'ı haftalar önceki olabilir (servis artık kapalı/değişmiş); pDNS kaydı eski IP'yi gösterir; arşiv silinmiş bir gerçeği yansıtır; bir e-posta artık geçerli değildir. **Her bulgunun "ne zaman doğru olduğunu" kaydet.** Eski veriyle güncel karar verme; kritik bulguyu **taze bir kaynakla** teyit et (ama "taze teyit" çoğu zaman aktife geçmek demektir — dengeyi OPSEC'le kur). Bayat veriye dayanan "gerçek zamanlı" iddia, klasik OSINT hatasıdır.

### 15.11 Çok kaynak çapraz doğrulama (operasyonel özet)
Bölüm 14'ün kalbi: **tek kaynak = sıfır kaynak.** Her bulguyu **bağımsız** ikinci kaynakla teyit et (Shodan + Censys, crt.sh + pDNS). "Bağımsız" = farklı toplama yöntemi, aynı verinin iki kopyası değil. **Döngüsel raporlamadan** kaçın (herkes aynı tek kaynaktan kopyalamış olabilir). Birincil kaynağa in.

### 15.12 Arşivleme / zaman damgası (kanıt bütünlüğü — operasyonel özet)
Bölüm 14'ün kalbi: doğruladığın her bulguyu **arşivle + zaman damgala + (gerekirse) imzala.** Hedef silse bile kanıtın kalır (`archive.today`, tam sayfa kaydı, ekran görüntüsü). IOC dosyalarını **GPG detached-sign** et (`GNUPG_GPG_USTALIK_REHBERI.md`) → "bu istihbarat şu tarihte bizdeydi, değişmedi." Kaynak→yöntem→tarih→analist zincirini belgele.

### 15.13 Kapsam (scope) disiplini & "tavşan deliği" tuzağı
OSINT **sonsuz** bir labirenttir; her pivot yeni pivot doğurur. Net bir **kapsam** (Bölüm 1, PIR/EEI) olmadan saatlerce ilgisiz veride kaybolursun ("rabbit hole"). **Her pivottan önce sor:** "Bu, cevapladığım soruya hizmet ediyor mu?" Etmiyorsa **dur.** Ayrıca kapsam **yasal/etik bir sınırdır** — bir CTI angajmanında müşterinin verdiği kapsam dışına çıkmak (üçüncü taraf altyapıyı taramak) hem yasal risk hem itibar riskidir.

### 15.14 Ham veri ≠ istihbarat (analiz katmanını atlama)
En sık amatör hatası: **araç çıktısını "istihbarat" sanmak.** theHarvester'ın 500 satırlık subdomain listesi, Shodan'ın 10.000 sonucu, SpiderFoot'un dev raporu — bunlar **ham veridir** (Bölüm 1, 3. aşama). İstihbarat, bunun **işlenmiş, doğrulanmış, bağlama oturmuş ve bir soruyu cevaplayan** hâlidir. Karar vericiye "işte 10.000 sonuç" verme; **"bu üç bulgu şu anlama geliyor ve şu eylemi gerektiriyor"** ver. Analiz ve doğrulama olmadan, topladığın her şey gürültüdür.

### 15.15 Operasyonel günlük (logbook) tut
Profesyonel OSINT, **kendi adımlarını belgeler:** hangi sorguyu ne zaman çalıştırdın, hangi kaynaktan ne buldun, hangi hipotezi neden eledin. Neden? (a) **Tekrarlanabilirlik** (başkası/sen 6 ay sonra aynı sonuca ulaşabilmeli); (b) **kanıt zinciri** (Bölüm 14); (c) **OPSEC denetimi** (yanlışlıkla aktife geçtin mi, kişiselle karıştın mı — geriye dönük yakalanır). Ama dikkat: günlüğün **kendisi** hassas bir saldırı yüzeyidir → şifreli sakla, gerektiğinde `/imha` (Bölüm 17).

### 15.16 Yasal/etik sınır = teknik bir kısıt gibi davran (operasyonel özet)
Bölüm 18'in özü, püf olarak: **"yapabilmek" ≠ "yapmaya hakkın olması."** Bir parolayı dork'la bulabilmen, onu **kullanma** hakkı vermez. Bir sisteme erişebilmen, erişmen gerektiği anlamına gelmez. OSINT'in tüm yasallığı **açık kaynakta gözlem** çizgisindedir; o çizgiyi geçtiğin an (izinsiz erişim, kimlik avı, sömürü) OSINT bitti, suç başladı. Bu sınırı **teknik bir kural gibi** koy — istisna yok.

### 15.17 İzole, tek kullanımlık ortam disiplini
Bölüm 5'in özü, püf olarak: araştırmayı **disposable VM/Tails/Qubes**'te yap, **baseline'a geri dön.** Web bug/canary/çerez/parmak izi kalıntısı her oturumda temizlensin. Ana sisteminle araştırma VM'i arasında **kes-yapıştır/dosya köprüsü bile** kurma. "Temiz oda" prensibi: içeri kişisel girmez, dışarı araştırma sızmaz.

### 15.18 Hedefe değil, hedefin "çevresine" bak (dolaylı keşif)
Çoğu zaman hedefin kendisi sağlam korunmuştur ama **çevresi** (tedarikçiler, eski çalışanlar, alt yükleniciler, kişisel sosyal hesaplar, eski sürümleri Wayback'te, çalışanların GitHub'ı) sızıntı doludur. Bir kurumu doğrudan "taramak" yerine, **etrafındaki halkayı** pasif tara — iç kullanıcı adı şeması bir PDF metadata'sından, teknoloji yığını bir iş ilanından, altyapı bir eski çalışanın GitHub commit'inden gelir. Bu hem daha verimli hem daha sessizdir (hedefe hiç dokunmazsın).

---

<a id="16"></a>
## 16. ☠️ Yaygın Ölümcül Hatalar

1. **Kişisel hesapla/IP ile araştırmak** → hedefe (LinkedIn bildirimi, log'daki IP) *kim olduğunu* sızdırmak; operasyonu ve kendini yakmak. (En sık felaket.)
2. **Pasif sanıp aktife geçmek** → theHarvester DNS brute / SpiderFoot spider / hedef sitesini ziyaret ile hedefin log'una iz bırakmak.
3. **Tek kaynağa güvenmek** → doğrulanmamış/bayat/dezenforme veriyi "istihbarat" sanıp yanlış karar (ve yanlış attribution) üretmek.
4. **Maltego'da "tüm transform'ları çalıştır"** → hedef verisini onlarca üçüncü taraf sunucuya saçmak (kim-kimi-araştırıyor sızıntısı) + API kotasını tüketmek.
5. **Sock puppet'i kişiselle çapraz kullanmak** → tek çapraz temasla (kişisel Gmail/telefon/kart) tüm puppet ağını gerçek kimliğine bağlamak.
6. **Canary/web bug'lı belgeyi kendi makinende açmak** → IP/zaman/kullanıcı-ajanını hedefe vermek; tuzağa düşmek.
7. **API anahtarını kod deposuna commit'lemek** → kendi anahtarının (ve araştırma geçmişinin) GitHub dork'la sızması.
8. **Rate-limit/kotayı patlatmak** → ban/throttle yemek, anahtarı bir taramada tüketmek, agresif sorguyla hedefin dikkatini çekmek.
9. **DNS leak'i görmezden gelmek** → VPN/Tor kullanırken DNS'in dışarıdan çıkıp gerçek niyetini/coğrafyanı sızdırması.
10. **Ham veriyi istihbarat sanmak** → 10.000 satır çıktıyı "rapor" diye teslim etmek; analiz/doğrulama katmanını atlamak.
11. **Zaman damgasız/arşivsiz çalışmak** → hedef veriyi silince kanıtın yok olması; verinin tazeliğini bilememek.
12. **Onaylama yanlılığı** → hipotezini çürütmeye değil doğrulamaya çalışmak; görmek istediğini görmek.
13. **Kapsamsız "tavşan deliğine" düşmek** → saatlerce ilgisiz pivotta kaybolmak; bazen yasal kapsam dışına taşmak.
14. **Yasal çizgiyi geçmek** → dork'la bulduğun parolayı/sırrı *kullanmak*, izinsiz erişim, kimlik avı → OSINT'ten suça geçmek.
15. **Sock puppet'i yaşlandırmadan/tutarsız kullanmak** → platform anti-spam'ine takılmak, hesabın yanması, hatta tutarsızlığın seni ele vermesi.

---

<a id="17"></a>
## 17. 🏰 Kanije Kalesi ile Birlikte Kullanım

Bu repo (Kanije Kalesi), CTI/güvenlik odaklı bir **muhafız ve araç setidir.** OSINT, onun **istihbarat ve öz-denetim** katmanını besler — felsefe birebir örtüşür: *"saldırganı tanı, kendini saldırgan gözüyle gör, kanıtı koru."* OSINT araç seti **dışarıyı** keşfeder; Kanije Kalesi **olay anını ve kendi savunmanı** yönetir. İkisi birlikte tam bir CTI döngüsü kurar.

| İhtiyaç | OSINT araç setinin rolü | Kanije Kalesi ile bağ |
|---|---|---|
| **Tehdit aktörü profilleme** | Aktörün altyapısını (Shodan `ssl:`/JARM pivot, crt.sh, pDNS), TTP'lerini, IOC'lerini pasif haritala | Toplanan IOC'lerle Kanije'nin tehdit/tespit bağlamını zenginleştir; `/defender` tespitlerini bilinen aktör IOC'leriyle eşleştir |
| **Kendi dijital ayak izini denetleme** | theHarvester/Shodan/dork ile **kurumunun saldırgan gözünden** ne göründüğünü çıkar (açık port, sızmış e-posta, GitHub secret, belge metadata'sı) | Kanije'nin koruduğu cihazın saldırı yüzeyini önceden kapat; `/erisim` ve `/tuzak` ile gözetlenen/sızdırılan dosyaları yakala |
| **IOC zenginleştirme** | Bir IP/alan adı/hash'i SpiderFoot/Recon-ng/Maltego ile zenginleştir (ilişkili altyapı, kötü amaçlılık skoru, geçmiş) | Zenginleştirilmiş IOC → Kanije'nin alarm/karar bağlamı; "bu IP bilinen bir C2" → daha sert aksiyon |
| **Canary/tuzak simetrisi** | Hedefin kurmuş olabileceği canary token/web bug'ları tanı ve kaçın (Bölüm 3.4, 15.8) | `/tuzak` ile **kendi** sistemine honeypot kur → sana yaklaşan saldırganı yakala (aynı silahın savunma yüzü) |
| **Kanıt bütünlüğü** | Topladığın OSINT bulgularını arşivle + zaman damgala + imzala (Bölüm 14) | `/imha` ile araştırma izini/hassas OSINT veritabanını gerektiğinde secure-wipe (Recon-ng/SpiderFoot DB bir saldırı yüzeyidir) |

### 🔥 Önerilen entegrasyon deseni
1. **Öz-denetim önce (kendini OSINT'le tara):** Düzenli olarak **kendi kurumunu/cihazını** pasif OSINT'ten geçir — saldırgan ne görüyor? Açık RDP (Shodan `org:` + `port:3389`), sızmış e-posta (theHarvester), GitHub secret (dork), belge metadata sızıntısı (exiftool). Bulduğun her açığı **saldırgandan önce** kapat. Kanije Kalesi'nin koruduğu uç noktanın **saldırı yüzeyini sıfıra** yaklaştır.
2. **Tehdit aktörü → IOC → Kanije bağlamı:** Bir aktörü pasif profillediğinde (altyapı pivotu, Bölüm 7/12) çıkardığın **IOC'leri** Kanije'nin tespit/karar mantığını zenginleştirmek için kullan — `/defender` bir tespit gösterdiğinde, onu bilinen aktör IOC'leriyle eşleştirip **daha isabetli** aksiyon al.
3. **Tuzak simetrisini kur:** OSINT yaparken *hedefin* canary'lerinden kaçınırken (Bölüm 15.8), aynı mantıkla **kendi** tarafında `/tuzak kur` ile honeypot dosyaları yerleştir → seni araştıran/sızan saldırgan **senin** tuzağına düşsün. Avcı ve av aynı silahı kullanır; sen iki tarafı da bil.
4. **Hassas OSINT verisini koru ve gerektiğinde imha et:** Recon-ng workspace DB'si, SpiderFoot tarama sonuçları, sock puppet kimlik bilgileri, API anahtarları → hepsi **hassas birer saldırı yüzeyidir.** Şifreli sakla (`KEEPASSXC_PAROLA_KALESI_REHBERI.md` / `VERACRYPT_USTALIK_REHBERI.md`); fiziksel tehditte Kanije `/koruma` + `/imha` ile araştırma izini güvenli sil.
5. **Bulguları imzala (kanıt zinciri):** CTI raporunu ve IOC listesini **GPG detached-sign** et (`GNUPG_GPG_USTALIK_REHBERI.md`) → "bu istihbarat bizden ve değişmedi" matematiksel olarak kanıtlanır. Kanije'nin ürettiği olay kayıtlarıyla birleşince **savunulabilir** bir kanıt zinciri çıkar.

> 🧠 **Felsefe örtüşmesi:** OSINT **dışarıyı görür** (saldırgan + saldırı yüzeyi), Kanije Kalesi **olay anını ve cihazı savunur** (kilit, tuzak, imha, izleme). Birlikte: *önce saldırganı ve kendi açığını tanı (OSINT), sonra açığı kapat ve olay anına hazır ol (Kanije).* İstihbarat olmadan savunma kördür; savunma olmadan istihbarat eylemsizdir.

---

<a id="18"></a>
## 18. ⚖️ Hukuki & Etik Sınır — Yalnızca Açık Kaynak

> Bu bölüm "ek bilgi" değil, rehberin **anayasasıdır.** OSINT'in tüm gücü ve tüm meşruiyeti, tek bir çizgide durmaktan gelir: **açık kaynakta gözlem.** O çizgiyi geçtiğin an, yaptığın şey OSINT değil **suçtur.**

### Çizginin neresinde durmalı

| ✅ Meşru OSINT (açık kaynak gözlem) | ❌ Çizgiyi geçmek (suç/etik ihlali) |
|---|---|
| crt.sh / pDNS / Shodan'ın *önceden topladığı* veriyi okumak | Hedef sisteme izinsiz **aktif tarama**, port-scan (yetkisiz) |
| Herkese açık LinkedIn/Twitter profilini görüntülemek | Hesaba **giriş denemesi**, parola tahmini/kırma |
| Açık bir PDF'in metadata'sını okumak | Korumalı/özel bir sisteme **erişmek** (yetkisiz) |
| Dork'la açık bir dizini bulmak | Bulunan **sırrı/parolayı kullanmak** (yetkisiz erişim) |
| Sock puppet'le **pasif gözlem** | Sahte kimlikle **sosyal mühendislik** / kandırarak bilgi koparma |
| Wayback'ten silinmiş sayfayı okumak | Sızmış veriyi **yaymak/satmak/şantaj** |
| Kendi kurumunu denetlemek | Üçüncü tarafı **yetkisiz** kapsam dışı taramak |

### Temel ilkeler
- **"Yapabilmek" ≠ "yapmaya hakkın olması."** Bir parolayı/sırrı dork'la *bulabilmen*, onu *kullanma* hakkı vermez. Açık bir kapı, içeri girme davetiyesi değildir. (Bölüm 15.16)
- **Yetki & kapsam (scope):** Bir CTI angajmanında **yazılı yetki** ve **net kapsam** olmadan aktif keşif yapma. Kapsam dışına çıkmak (müşterinin sahibi olmadığı altyapıyı taramak) yasal sorumluluk doğurur.
- **Yetkisiz erişim yasaları:** Çoğu ülkede "bilişim sistemine yetkisiz erişim" suçtur (örn. ABD'de CFAA, Türkiye'de TCK md. 243-245 bilişim suçları, BK'de Computer Misuse Act). **Açık kaynağı okumak** bunun dışındadır; **bir sisteme girmek/erişmek** içindedir.
- **Kişisel veri (KVKK/GDPR):** Kişiler hakkında OSINT yaparken kişisel veri koruma yasaları geçerlidir — toplama amacı meşru, ölçülü ve hukuka uygun olmalı; gözetleme/taciz/profilleme amaçlı kötüye kullanım ayrı bir suçtur.
- **Sock puppet & platform şartları:** Sahte hesaplar çoğu platformun şartlarına aykırıdır; bu rehber bunları **kişisel kimliği koruma + pasif gözlem** için anlatır, aldatma/dolandırıcılık için değil.

> 🚫 **Kırmızı çizgi (tek cümle):** *OSINT, açık pencereden görüneni okumaktır; kapıyı zorlamak, içeri girmek, bulduğunu kullanmak ya da birini kandırmak değildir.* Bu çizgiyi teknik bir kural gibi koy — istisna yok. Şüphedeysen, **durmak** doğru cevaptır.

> ⚖️ **Yasal sorumluluk:** Bu rehber **meşru CTI, savunma, gazetecilik, akademik araştırma ve öz-denetim** içindir. Bulunduğun yargı bölgesinin yasalarına uy; kurumsal bir angajmandaysan **yazılı yetki ve kapsam** olmadan tek paket gönderme. Doğru komut/filtre verilmeye çalışıldı; emin olmadığın her ayrıntıyı (Shodan plan kapsamı, araç bayrak adları, API konumları) **resmi dokümandan teyit et.**

---

<a id="19"></a>
## 19. ✅ Hızlı Referans & OPSEC Kontrol Listesi

### Araç → kullanım tablosu

| Araç | Ne için | Tipik komut / kullanım | Pasif mi? |
|---|---|---|---|
| **theHarvester** | E-posta/subdomain/host toplama | `theHarvester -d ornek.com -b all -l 500 -f rapor` | Pasif (DNS brute modülü **aktif**) |
| **Shodan (CLI)** | İnternete bakan cihaz/banner | `shodan search 'ssl:"ornek.com"'` · `shodan host 1.2.3.4` · `shodan count '...'` | **Pasif** (Shodan'ın arşivi) |
| **SpiderFoot** | Otomatik OSINT + korelasyon | `python3 sf.py -s ornek.com -t DOMAIN_NAME` · **"Passive" use case** | Modül bazlı (bazıları **aktif**) |
| **Maltego** | Görsel ilişki/graf analizi | Entity → transform → graf; **local vs remote** transform seç | Transform bazlı (remote = **veri dışarı**) |
| **Recon-ng** | Modüler keşif çatısı | `workspaces create` · `modules load recon/...` · `run` | Modül bazlı (discovery/brute **aktif**) |
| **crt.sh / CT** | Subdomain (sertifikadan) | `curl "https://crt.sh/?q=%25.ornek.com&output=json"` | **Pasif** (saf) |
| **Passive DNS** | Tarihsel DNS / altyapı pivotu | 3. taraf pDNS sorgusu (IP↔alan adı geçmişi) | **Pasif** (saf) |
| **Wayback** | Silinmiş/eski içerik | `curl "http://web.archive.org/cdx/search/cdx?url=ornek.com*..."` | **Pasif** (saf) |
| **WHOIS** | Tescil bilgisi | `whois ornek.com` | Pasif (tescile sorar) |
| **ExifTool** | Dosya metadata/EXIF | `exiftool -gps:all -a foto.jpg` · `exiftool *.pdf` | **Pasif** (yerel dosya) |
| **Google/GitHub dork** | Açık hassas içerik | `site:ornek.com filetype:pdf` · `org:X filename:.env` | **Pasif** (arama index'i) |
| **Censys** | Sertifika/host derin sorgu | Sertifika/host korelasyonu (Shodan çapraz doğrulama) | **Pasif** (arşiv) |

### 🔥 OPSEC kontrol listesi

**Operasyon öncesi (kurulum)**
- [ ] **PIR/EEI tanımlı** — neyi öğreneceğim, hangi soruya cevap arıyorum
- [ ] **Ayrı VM / Tails / Qubes** hazır (disposable, baseline snapshot)
- [ ] **VPN ya da Tor** aktif + **DNS leak testi geçti**
- [ ] **Sock puppet'ler** yaşlandırılmış, kapakla tutarlı, **ayrı altyapıda**
- [ ] **API anahtarları** araştırmaya özel hesaptan, **şifreli**, kod deposunda **değil**
- [ ] Tarayıcı **parmak izi** standart (Tor/Mullvad) ya da ayrı profil
- [ ] Zaman dilimi/dil **kapakla uyumlu** (kişiselle değil)

**Toplama sırasında**
- [ ] **Önce pasif** — tüm pasif kaynaklar tükenmeden aktife geçme
- [ ] Aracın her modülü için "**hedefe dokunuyor mu?**" sorusu cevaplandı
- [ ] SpiderFoot **"Passive"**, Recon-ng aktif modülleri kapalı, theHarvester brute kapalı
- [ ] Maltego **remote transform** bilinçli (hedef verisi nereye gidiyor?)
- [ ] **Rate-limit/kota** korunuyor (`shodan count` ile keşif, throttle)
- [ ] **Kişisel hesap/IP ile asla** — LinkedIn vb. profil görüntüleme bildirimi riski yok
- [ ] Şüpheli belge → **izole VM, makro kapalı, uzak içerik kapalı** (canary/web bug)

**Analiz & doğrulama**
- [ ] Her bulgu **≥2 bağımsız kaynaktan** çapraz doğrulandı
- [ ] **Yanlış-pozitif** elendi (shared/CDN IP? jenerik sertifika/favicon?)
- [ ] **Dezenformasyon/false flag** sorgulandı (fazla mı temiz/uygun?)
- [ ] **Verinin tazeliği** (zaman damgası) kayıtlı; bayat veriyle güncel iddia yok
- [ ] Hipotez **çürütülmeye** çalışıldı (onaylama yanlılığına karşı), alternatif hipotez var
- [ ] Güven seviyesi **dereceli** ("yüksek güvenle X ile tutarlı", "kesin" değil)

**Yayım & kapanış**
- [ ] Bulgular **arşivlendi + zaman damgalandı** (archive.today / tam sayfa / ekran görüntüsü)
- [ ] IOC/rapor **GPG ile imzalandı** (kanıt bütünlüğü)
- [ ] Ham veri → **işlenmiş istihbarat**a dönüştü (analiz katmanı atlanmadı)
- [ ] **Operasyonel günlük** (logbook) tam — tekrarlanabilir + denetlenebilir
- [ ] Hassas veri (DB, puppet kimlikleri, anahtarlar) **şifreli saklandı / gerekirse `/imha`**
- [ ] **Yasal/etik çizgi** hiç geçilmedi (yalnızca açık kaynak, yetki/kapsam içinde)

---

> 🏰 **Kapanış:** OSINT bir araç koleksiyonu değil, bir **disiplindir.** En güçlü Shodan dork'u bile, kişisel hesabınla baktığın bir LinkedIn profilinde, DNS'in sızdığı bir VPN oturumunda ya da tek kaynağa güvenip doğrulamadığın bir "bulguda" çaresizdir — hatta tehlikelidir. OSINT sana **dünyanın açık penceresinden bakma** gücünü verir; **o pencereden bakarken görünmez kalmak, gördüğünü doğrulamak ve çizgiyi geçmemek senin işin.** İstihbarat döngüsünün kalbi araçlar değil, **planlama, OPSEC ve analizdir** — araçlar yalnızca toplama ayağıdır. Kanije Kalesi de tam burada: sen dışarıyı keşfederken, **kendi kapını saldırgan gözüyle görmeni** ve olay anında **kilitlemeni** sağlayan nöbetçi olarak devreye girer.
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. İlgili: `GNUPG_GPG_USTALIK_REHBERI.md` (kanıt imzalama/bütünlük), `VERACRYPT_USTALIK_REHBERI.md` (araştırma verisini şifreli saklama), `KEEPASSXC_PAROLA_KALESI_REHBERI.md` (sock puppet kimlik/anahtar yönetimi), `TAILS_ANONIMLIK_REHBERI.md` (amnezik araştırma ortamı), `QUBES_OS_BOLMELEME_REHBERI.md` (kimlik başına izole bölmeleme), `WINDOWS11_HARDENING_KALE.md` / `LINUX_HARDENING_KALE.md` (araştırma uç noktasını sertleştirme).*
