# 🎯 MITRE ATT&CK — TAM USTALIK REHBERİ (CTI Perspektifi)
## Matris Anatomisinden Navigator'a, Aktör Profillemeden Rapor Yazımına — Sıfırdan Uzmana

> **Amaç:** MITRE ATT&CK'i "bir technique ID'si ezberlemek" seviyesinden çıkarıp, bir **siber tehdit istihbaratı (CTI)** analisti gibi — **gözlemi technique'e mapleyebilen, aktör türünü TTP imzasından çıkarabilen, atfetme tuzaklarını bilen ve doğru raporda doğru ATT&CK ağırlığını kullanan** — ustaca kullanmayı öğretmek. Bu rehber yalnızca *ATT&CK nedir*'i değil, **bir saldırı anlatısını nasıl technique zincirine çevirirsin**, **gözlemlenen TTP'lerden hangi tür aktörle karşı karşıya olduğunu nasıl çıkarırsın** ve **hangi raporda ne kadar ATT&CK kullanmalısın**'ı da anlatır. Hiç bilmeyen biri bu tek dosyayı okuyarak ATT&CK'i operasyonel olarak kullanan bir analiste dönüşebilmeli.

> ⚠️ **Önce bunu oku — ATT&CK'in iki ölümcül yanılgısı:**
> 1. **"Technique ID saymak = istihbarat."** Yanlış. Bir rapora 40 tane T-numarası serpiştirmek seni analist yapmaz; **bağlamsız ID gürültüdür.** ATT&CK'in değeri, gözlemleri *ortak bir dile* çevirip aktör davranışını **karşılaştırılabilir** kılmasındadır — sayıda değil, **ilişkide.**
> 2. **"TTP'ler eşleşti, demek ki kesin X aktörü."** Yanlış. TTP'ler **çakışır** (birçok grup aynı PowerShell/phishing'i kullanır), **kopyalanır** (copy-cat), **kasten taklit edilir** (false flag). ATT&CK'i atfetme için kullanmak güçlüdür ama **acımasızca tuzaklıdır** — bu yüzden Bölüm 9 (Aktör Profilleme) ve atfetme tehlikeleri bu rehberin en uzun, en dikkatli bölümleridir.
>
> Bu iki yanılgı, neden bu rehberin "saldırı türü bulma antrenmanı" (Bölüm 10) ve "aktör türü belirleme" (Bölüm 9) bölümlerinin en çok emek verilen kısımlar olduğunu açıklar. Onları atlama.

> 🧭 **Kapsam ayrımı (oku):** Bu rehber **ATT&CK ARACINI/FRAMEWORK'ünü** nasıl kullanacağına odaklanır. **TTP kavramının derin teorisi** (Indicator vs TTP farkı, Pyramid of Pain'in tam mantığı, davranış-temelli avcılık felsefesi) **ayrı bir dosyada** işlenir: `TTP_AVCILIGI_USTALIK_REHBERI.md`. Burada teoriyi tekrarlamıyoruz; ATT&CK'i **pratikte** kullanıyoruz. Teorik derinlik için sürekli o dosyaya çapraz referans vereceğiz.

> 🔢 **ID doğruluğu sözü:** Bu rehberdeki technique ID'leri (T-numaraları), tactic isimleri ve grup numaraları (G-numaraları) **dikkatle** verilmiştir. Yine de ATT&CK **canlı bir framework**'tür — sürümle (v14, v15, v16…) technique'ler bölünür, birleşir, yeniden numaralanır (örn. eski "deprecated" ID'ler). **Bir ID'yi rapora koymadan önce her zaman `attack.mitre.org`'dan teyit et.** Emin olmadığım yerleri açıkça *"attack.mitre.org'dan teyit et"* notuyla işaretledim — **yanlış ID, tüm rehberi çürütür.**

---

## 📑 İÇİNDEKİLER

1. [ATT&CK Nedir, Neden Var — Saldırgan Davranışının Ortak Dili](#1)
2. [Tarihçe & Cyber Kill Chain'den Farkı](#2)
3. [🔥 MATRİS ANATOMİSİ — 14 Tactic Tek Tek](#3)
4. [Technique vs Sub-technique vs Procedure — Granülarite](#4)
5. [Enterprise / Mobile / ICS Matrisleri](#5)
6. [🔥 Bir Technique Sayfasını OKUMA (T1059 baştan sona)](#6)
7. [Groups (G-numaraları) & Software (S-numaraları) Sayfaları](#7)
8. [🔥 attack.mitre.org'da Arama & İlişkilerde Gezinme](#8)
9. [🔥 ATT&CK Navigator — Sıfırdan Uzmana](#9)
10. [🔥🔥 MANUEL TEHDİT AKTÖRÜ TÜRÜ BELİRLEME (Profilleme Metodolojisi)](#10)
11. [🔥🔥 SALDIRI TÜRÜ BULMA ANTRENMANI (Çözümlü Egzersizler)](#11)
12. [🔥🔥 CTI RAPOR YAZIMINDA ATT&CK — Hangi Raporda Ne Kadar](#12)
13. [ATT&CK'in Kill Chain · Diamond · Pyramid of Pain ile İlişkisi](#13)
14. [Gerçek Aktörler & İmza TTP'leri](#14)
15. [☠️ Yaygın Ölümcül Hatalar](#15)
16. [🏰 Kanije Kalesi ile Birlikte Kullanım](#16)
17. [🗺️ Sıfırdan Uzmana Yol Haritası (Hafta Hafta)](#17)
18. [✅ Hızlı Referans (14 Tactic Tablosu) & Kontrol Listesi](#18)
19. [📚 Kaynaklar & Çapraz Referanslar](#19)

---

<a id="1"></a>
## 1. 🧭 ATT&CK Nedir, Neden Var — Saldırgan Davranışının Ortak Dili

**MITRE ATT&CK** (Adversarial Tactics, Techniques, and Common Knowledge — "Düşman Taktikleri, Teknikleri ve Ortak Bilgisi"), gerçek dünyadaki tehdit aktörlerinin **gözlemlenmiş davranışlarını** sistematik olarak belgeleyen, herkese açık ve sürekli güncellenen bir **bilgi tabanı (knowledge base)** ve **çerçevedir (framework)**. Kritik nokta şudur: ATT&CK bir "araç" değil, bir **ortak dildir.** Bir aktörün "ne yaptığını" (davranışını), dünyanın her yerindeki analistlerin **aynı şekilde adlandırıp karşılaştırabileceği** standart bir sözlüğe çevirir.

### Neden bir "ortak dil"e ihtiyaç var?

ATT&CK'ten önce, iki farklı güvenlik ekibi aynı saldırıyı tamamen farklı kelimelerle anlatırdı. Biri "saldırgan zararlı bir Office belgesi gönderdi" der, diğeri "makro tabanlı initial access vektörü" derdi. Aynı olay, **karşılaştırılamaz** iki cümleydi. Sonuç: istihbarat paylaşılamaz, savunma önceliklendirilemez, "bu aktör de aynı şeyi mi yapıyor?" sorusu cevaplanamazdı.

ATT&CK bunu çözer: o davranış artık **T1566.001 (Spearphishing Attachment)** + **T1204.002 (Malicious File)** + **T1059.001 (PowerShell)** zinciridir. Artık iki ekip, iki ülke, iki vendor **tam olarak aynı şeyden** bahsettiğini bilir. Bir aktörün TTP setini başka bir aktörünkiyle **kelimesi kelimesine** kıyaslayabilirsin.

> 🧠 **Altın kural:** ATT&CK'in tek bir cümlelik özü: *"Saldırgan davranışını, herkesin aynı anladığı standart bir kataloğa indeksle."* IOC'ler (IP, hash, domain) saldırganın **değiştirebildiği** şeylerdir; ATT&CK ise saldırganın **nasıl davrandığını** — değiştirmesi çok daha pahalı olan kısmı — yakalar. (Bunun derin teorisi: Pyramid of Pain → `TTP_AVCILIGI_USTALIK_REHBERI.md`.)

### ATT&CK ne DEĞİLDİR (yaygın kafa karışıklığı)

- **Bir saldırı sırası/playbook değildir.** Tactic'lerin matristeki **soldan sağa sırası** zorunlu bir akış değildir — saldırgan tactic'ler arasında atlar, geri döner, bazılarını hiç kullanmaz. (Bu, Kill Chain'den en büyük farkıdır — Bölüm 2.)
- **Bir zafiyet/CVE listesi değildir.** ATT&CK *davranışı* anlatır ("credential dumping yaptı"), *hangi açığı kullandığını* değil. Zafiyet tarafı için CVE/CWE ayrı dünyalardır.
- **Bir savunma ürünü değildir.** ATT&CK bir **referans**'tır; onu savunmaya, tespite, kırmızı takıma, istihbarata *sen* uygularsın.
- **Bir araç değildir** (Navigator hariç — o da yalnızca matrisi görselleştiren bir yardımcı, Bölüm 9). ATT&CK'in kendisi bir **web sitesi + JSON veri kümesidir** (`attack.mitre.org`).

### CTI'da ATT&CK'in dört ana kullanım alanı

| Kullanım | Ne işe yarar | Kim kullanır |
|---|---|---|
| **Tehdit istihbaratı** | Aktörleri/kampanyaları TTP setiyle profille, karşılaştır, izle | CTI analisti |
| **Tespit & avcılık (detection/hunting)** | "Bu technique'i hangi data source ile yakalarım?" → tespit geliştir | SOC, threat hunter |
| **Kırmızı takım / emülasyon** | Bir aktörü taklit ederek savunmanı test et (Atomic Red Team, Bölüm 11) | Red team, purple team |
| **Savunma boşluk analizi** | "Hangi technique'lere karşı kör noktam var?" → kapsama haritası | Blue team, savunma mimarı |

Bu rehber **CTI/tehdit istihbaratı** perspektifine odaklanır (profilleme, raporlama, atfetme) ama dördü de iç içedir.

---

<a id="2"></a>
## 2. 🕰️ Tarihçe & Cyber Kill Chain'den Farkı

### Kısa tarihçe

- **MITRE**, ABD merkezli, kâr amacı gütmeyen bir araştırma kuruluşudur (devlet için Ar-Ge yapar; CVE listesini de o yönetir).
- ATT&CK, **2013'te** MITRE içinde bir araştırma projesi olarak doğdu: "Gerçek bir kurumsal ağda saldırganların *gözlemlenen* post-compromise davranışlarını kataloglayabilir miyiz?" sorusuyla. İlk odak: **Windows kurumsal ortamı.**
- **2015'te** kamuya açıldı. Zamanla genişledi: önce Windows, sonra macOS/Linux, ardından **Mobile** (iOS/Android), **ICS** (endüstriyel kontrol sistemleri), **Cloud** (Enterprise içinde alt-platformlar olarak).
- **PRE-ATT&CK** adında ayrı bir "saldırı öncesi" matrisi vardı; sonradan Enterprise'a **Reconnaissance** ve **Resource Development** tactic'leri olarak **entegre edildi** (artık ayrı PRE matrisi yok — bunu bilen az, sık karıştırılır).
- ATT&CK **versiyonlanır** (v1, v2, … v14, v15, v16…). Her sürüm technique ekler/böler/günceller. Rapor yazarken **hangi sürümü** baz aldığını belirtmek profesyonelliktir.

### 🔥 Cyber Kill Chain'den Farkı (en kritik kavramsal ayrım)

**Lockheed Martin Cyber Kill Chain** (2011), bir saldırıyı **7 doğrusal aşamaya** böler: Reconnaissance → Weaponization → Delivery → Exploitation → Installation → Command & Control → Actions on Objectives. Yıkıcı derecede etkili ama **yüksek seviyeli ve doğrusaldır** — "saldırı şu sırayla ilerler" der, ama *her aşamada saldırgan SOMUT olarak ne yapar*'ı detaylandırmaz.

ATT&CK, Kill Chain'in **"ne yapıldığını" boyutunu derinleştirir:** her tactic (= bir saldırı amacı/aşama) altında **onlarca somut technique** (= o amacı gerçekleştirmenin gözlenmiş yolları) listeler.

```
   CYBER KILL CHAIN (7 doğrusal aşama — "saldırı sırası")
   ┌──────────┬──────────────┬──────────┬─────────────┬──────────────┬──────┬───────────────────┐
   │ Recon    │ Weaponization│ Delivery │ Exploitation│ Installation │ C2   │ Actions on Object.│
   └────┬─────┴──────┬───────┴────┬─────┴──────┬──────┴──────┬───────┴──┬───┴─────────┬─────────┘
        │            │            │            │             │          │             │
        ▼            ▼            ▼            ▼             ▼          ▼             ▼
   MITRE ATT&CK (14 tactic — her biri "saldırganın amacı" + altında somut technique'ler)
   ┌──────────────────────────────────────────────────────────────────────────────────────────┐
   │ Recon │Resource │Initial │Execution│Persist.│PrivEsc│Defense│Cred. │Discov.│Lateral│Coll.│ │
   │       │ Dev.    │ Access │         │        │       │Evasion│Access│       │ Move  │     │ │
   │       │         │        │         │        │       │       │      │       │       │ C2  │ │
   │       │         │        │         │        │       │       │      │       │       │Exfil│ │
   │       │         │        │         │        │       │       │      │       │       │Impact│ │
   └──────────────────────────────────────────────────────────────────────────────────────────┘
        ▲ ATT&CK doğrusal DEĞİL: saldırgan tactic'ler arasında atlar, döner, bazılarını atlar.
```

| Boyut | Cyber Kill Chain | MITRE ATT&CK |
|---|---|---|
| **Yapı** | 7 aşama, **doğrusal** (sıralı akış) | 14 tactic, **doğrusal DEĞİL** (atlamalı, döngüsel) |
| **Granülarite** | Yüksek seviye (aşama adı) | Derin (tactic → technique → sub-technique → procedure) |
| **"Nasıl" detayı** | Yok/az | **Yüzlerce somut technique** |
| **Odak** | Saldırının **akışı/önlenmesi** | Saldırganın **somut davranışı/tespiti** |
| **Kapsam** | Tüm saldırı yaşam döngüsü (üst düzey) | Çoğunlukla **post-/peri-compromise** davranış, derin |
| **CTI kullanımı** | Aşama-bazlı genel konumlandırma | TTP profilleme, karşılaştırma, tespit eşleme |

> 🧠 **İkisi rakip değil, katmanlıdır:** Kill Chain "saldırının hangi *aşamasında*yız?" sorusuna üst düzey cevap verir (yönetici diline yakın); ATT&CK "o aşamada saldırgan *tam olarak ne yaptı* ve bunu *nasıl yakalarım*?" sorusuna derin cevap verir (analist diline yakın). Olgun bir CTI ekibi ikisini birlikte kullanır: Kill Chain'le hikâyeyi *çerçevele*, ATT&CK'le içini *doldur*. (Diamond Model ile ilişki: Bölüm 13.)

> 💡 **Pratik köprü:** "Initial Access" tactic'i ≈ Kill Chain'in "Delivery + Exploitation"ı; "Command and Control" tactic'i ≈ Kill Chain'in "C2"si; "Exfiltration + Impact" ≈ "Actions on Objectives". Birebir değil ama zihinsel eşleme kurmak (Bölüm 13 ASCII'sine bak) ATT&CK'i Kill Chain bilen birine hızla öğretir.

---

<a id="3"></a>
## 3. 🔥 MATRİS ANATOMİSİ — 14 Enterprise Tactic Tek Tek

Matrisi okumak ustalığın temelidir. **Sütunlar = tactic'ler** (saldırganın *amaçları*, "neden bunu yapıyor"), **her sütunun altındaki kutular = technique'ler** (o amacı gerçekleştirmenin *gözlenmiş yolları*, "nasıl yapıyor").

> 🧠 **Tactic = "NEDEN", Technique = "NASIL".** Bu cümleyi ezberle. "Saldırgan parolaları çaldı" → *neden* = **Credential Access** (tactic), *nasıl* = **OS Credential Dumping / T1003** (technique). Her gözlemde önce "saldırganın amacı neydi?" (tactic) diye sor, sonra "bunu hangi yöntemle yaptı?" (technique) diye sor.

Enterprise matrisinde **14 tactic** vardır. İşte her birinin **saldırganın amacı** mantığı + örnek technique'leri (ID'leri attack.mitre.org'dan teyit edilebilir):

### TA0043 — Reconnaissance (Keşif)
> **Saldırganın amacı:** Hedef hakkında, ona *dokunarak veya dokunmadan*, gelecekteki operasyon için bilgi toplamak. (Saldırı *öncesi* aşama; eski PRE-ATT&CK'ten gelir.)
- **T1595** — Active Scanning (hedefin altyapısını tarama)
- **T1589** — Gather Victim Identity Information (e-posta, kimlik bilgisi toplama)
- **T1598** — Phishing for Information (bilgi sızdırmaya yönelik kimlik avı)

> 🔗 Bu, **OSINT** dünyasıyla doğrudan örtüşür — saldırganın keşif tarafını derinlemesine `OSINT_ARAC_SETI_USTALIK_REHBERI.md` işler (theHarvester, Shodan, crt.sh, pasif/aktif keşif ayrımı). ATT&CK tarafında: Reconnaissance technique'lerini gördüğünde, savunma açısından "saldırgan beni keşfediyor" demektir.

### TA0042 — Resource Development (Kaynak Geliştirme)
> **Saldırganın amacı:** Operasyon için altyapı/araç/kimlik *edinmek, kiralamak, satın almak veya üretmek* (saldırıdan önce silahını hazırlamak).
- **T1583** — Acquire Infrastructure (domain, sunucu, VPS edinme)
- **T1587** — Develop Capabilities (kendi malware/sertifikasını geliştirme)
- **T1586** — Compromise Accounts (mevcut hesapları ele geçirip kullanma)

> 💡 CTI değeri büyük: Bir aktörün **altyapı edinme deseni** (hangi registrar, hangi hosting, kendi mi yazıyor satın mı alıyor) güçlü bir parmak izidir. Altyapı pivotu (Shodan `ssl:`/JARM, pDNS) tam burada CTI'ya bağlanır → `OSINT_ARAC_SETI_USTALIK_REHBERI.md` Bölüm 7/12.

### TA0001 — Initial Access (İlk Erişim)
> **Saldırganın amacı:** Hedef ağa/sisteme **ilk ayağı sokmak** (kapıdan içeri girmek).
- **T1566** — Phishing (kimlik avı — en yaygın initial access vektörü)
- **T1190** — Exploit Public-Facing Application (internete bakan uygulamadaki açığı sömürme)
- **T1078** — Valid Accounts (çalınmış/meşru kimlik bilgileriyle giriş)
- **T1133** — External Remote Services (VPN/RDP gibi uzak servisler üzerinden)

> 🧠 "Saldırgan içeri *nasıl* girdi?" sorusunun cevabı bu tactic'tedir. Bir incident'te **ilk** bakacağın tactic budur (kök neden / patient zero).

### TA0002 — Execution (Yürütme)
> **Saldırganın amacı:** Hedef sistemde **kötü amaçlı kodu çalıştırmak.**
- **T1059** — Command and Scripting Interpreter (PowerShell, Bash, cmd, Python… — bkz. Bölüm 6 derin inceleme)
- **T1204** — User Execution (kullanıcıyı bir şeyi çalıştırmaya kandırma — makro, link)
- **T1053** — Scheduled Task/Job (zamanlanmış görevle çalıştırma — ayrıca persistence!)

> 💡 Birçok technique **birden çok tactic'e** hizmet eder (T1053 hem Execution hem Persistence hem Privilege Escalation). Bu yüzden bir technique'i "hangi amaçla kullanıldığı"yla (tactic context) okumak şart.

### TA0003 — Persistence (Kalıcılık)
> **Saldırganın amacı:** Sistem yeniden başlasa, kullanıcı çıkış yapsa, kimlik bilgisi değişse bile **erişimini korumak** (içeride kalmak).
- **T1547** — Boot or Logon Autostart Execution (başlangıçta otomatik çalışma — Registry Run anahtarları vb.)
- **T1053** — Scheduled Task/Job (zamanlanmış görev kalıcılığı)
- **T1136** — Create Account (kalıcı erişim için yeni hesap açma)
- **T1505** — Server Software Component (web shell vb.)

> 🧠 **Persistence, APT'lerin imza tactic'idir.** Uzun süre, sessizce içeride kalmak isteyen devlet-destekli aktörler burada *yaratıcıdır* (Bölüm 10 — aktör türü ayrımı). Fidye grubu genelde persistence'a az yatırım yapar (vur-kaç).

### TA0004 — Privilege Escalation (Yetki Yükseltme)
> **Saldırganın amacı:** Daha **yüksek izinler** elde etmek (kullanıcı → admin → SYSTEM/root).
- **T1548** — Abuse Elevation Control Mechanism (UAC bypass, sudo kötüye kullanımı)
- **T1068** — Exploitation for Privilege Escalation (yerel açık sömürerek yükselme)
- **T1055** — Process Injection (başka bir işleme kod enjekte ederek — ayrıca defense evasion!)

> 💡 Execution, Persistence, Privilege Escalation sık **iç içe**dir: aynı technique (Scheduled Task, Valid Accounts) üç amaca birden hizmet eder. ATT&CK bunu "bir technique birden çok tactic altında listelenir" diye gösterir.

### TA0005 — Defense Evasion (Savunma Atlatma)
> **Saldırganın amacı:** Tespit edilmekten **kaçınmak**, güvenlik araçlarını **körleştirmek**, izleri **gizlemek.** (ATT&CK'in **en kalabalık** tactic'i — onlarca technique.)
- **T1070** — Indicator Removal (log silme, dosya/artefakt temizleme, **timestomp**)
- **T1036** — Masquerading (meşru dosya/işlem gibi görünme — isim/yol taklidi)
- **T1027** — Obfuscated Files or Information (kod gizleme, base64, packing)
- **T1562** — Impair Defenses (AV/güvenlik aracını devre dışı bırakma)
- **T1055** — Process Injection (meşru işlem içinde gizlenme)

> 🔗 **Kanije ilgisi:** Kanije Kalesi'nin kendi kodu, savunma atlatma *karşıtı* bir farkındalıkla yazılmıştır — `/imha` mantığı **T1561 (Disk Wipe)** AV sinyalini *kasten vermez* (Bölüm 16). T1036 (Masquerading) ve T1070 (Indicator Removal) bir muhafızın "ne yapMAMASI gerektiğini" bilmesi açısından kritiktir.

### TA0006 — Credential Access (Kimlik Bilgisi Erişimi)
> **Saldırganın amacı:** Parola, hash, token, anahtar gibi **kimlik bilgilerini çalmak** (başka sistemlere geçişin/yükselmenin yakıtı).
- **T1003** — OS Credential Dumping (LSASS dump, SAM, /etc/shadow — Mimikatz sınıfı)
- **T1110** — Brute Force (parola deneme/spray)
- **T1555** — Credentials from Password Stores (tarayıcı/parola yöneticisi kasalarından)
- **T1056** — Input Capture (keylogging)

> 🔗 Savunma tarafında parola hijyeni doğrudan buna karşı koyar → `KEEPASSXC_PAROLA_KALESI_REHBERI.md`. Bir incident'te T1003 görmek **ciddiyet** işaretidir (yanal harekete zemin hazırlanıyor).

### TA0007 — Discovery (Keşif/İç Gözlem)
> **Saldırganın amacı:** Ele geçirilen ortamı **içeriden tanımak** — hangi makineler, kullanıcılar, ağ, savunma var? (Recon *dışarıdan*, Discovery *içeriden*.)
- **T1057** — Process Discovery (çalışan işlemleri listeleme)
- **T1082** — System Information Discovery (OS, donanım bilgisi)
- **T1018** — Remote System Discovery (ağdaki diğer makineleri bulma)
- **T1083** — File and Directory Discovery
- **T1033** — System Owner/User Discovery

> 🧠 Discovery technique'leri çoğu zaman **meşru sistem komutlarıyla** yapılır (`whoami`, `net view`, `tasklist`) → "living off the land". Bu yüzden tespiti zordur ve davranışsal analiz gerektirir (→ `TTP_AVCILIGI_USTALIK_REHBERI.md`).

### TA0008 — Lateral Movement (Yanal Hareket)
> **Saldırganın amacı:** İlk eriştiği makineden **ağdaki diğer sistemlere yayılmak.**
- **T1021** — Remote Services (RDP, SMB, SSH, WinRM ile diğer makinelere geçiş)
- **T1550** — Use Alternate Authentication Material (pass-the-hash, pass-the-ticket)
- **T1570** — Lateral Tool Transfer (araçları makineler arası taşıma)

> 💡 Lateral movement, **bir incident'in "tek makine mi, tüm ağ mı?" sorusunu** belirler. Görüldüğünde olay büyür: artık containment (izolasyon) kritiktir.

### TA0009 — Collection (Toplama)
> **Saldırganın amacı:** Sızdırmadan **önce** ilgilendiği veriyi **toplamak/paketlemek.**
- **T1113** — Screen Capture (ekran görüntüsü)
- **T1056** — Input Capture (klavye/girdi yakalama — Credential Access ile örtüşür)
- **T1005** — Data from Local System
- **T1560** — Archive Collected Data (sızdırmadan önce sıkıştırıp şifreleme — bir .rar/.zip hazırlama)

> 🧠 Collection + Exfiltration birlikte, aktörün **asıl niyetini** ele verir: ekran/belge toplayıp dışarı atan = casusluk; veriyi şifreleyip kilitleyen = fidye (Impact). Niyet okuması Bölüm 10'un kalbidir.

### TA0011 — Command and Control (C2 — Komuta Kontrol)
> **Saldırganın amacı:** Ele geçirilen sistemle **iletişim kurmak**, komut göndermek, veri almak (uzaktan kumanda kanalı).
- **T1071** — Application Layer Protocol (HTTP/HTTPS/DNS üzerinden C2)
- **T1573** — Encrypted Channel (şifreli C2 trafiği)
- **T1572** — Protocol Tunneling
- **T1090** — Proxy (trafiği gizleme/yönlendirme)

> 🔗 C2 trafiğini ağ seviyesinde yakalamak → `WIRESHARK_AG_ANALIZ_USTALIK_REHBERI.md`. C2 altyapısını pasif haritalama (JARM, sertifika pivotu) → `OSINT_ARAC_SETI_USTALIK_REHBERI.md`. C2 IOC'lerini yaşatma/paylaşma → MISP (`MISP_*` rehberi, Bölüm 19).

### TA0010 — Exfiltration (Sızdırma)
> **Saldırganın amacı:** Toplanan veriyi **ağ dışına çıkarmak** (çalmak).
- **T1041** — Exfiltration Over C2 Channel (veriyi C2 kanalından sızdırma)
- **T1567** — Exfiltration Over Web Service (cloud depolama/legit servis üzerinden)
- **T1048** — Exfiltration Over Alternative Protocol (DNS, FTP gibi alternatif kanal)

> 🧠 Exfiltration görmek = **veri kaybı gerçekleşti/gerçekleşiyor** demektir; en yüksek aciliyet. Fidye gruplarının "double extortion" (önce çal, sonra şifrele) modelinde Exfiltration + Impact birlikte görülür.

### TA0040 — Impact (Etki)
> **Saldırganın amacı:** Veriyi/sistemi **bozmak, yok etmek, erişilemez kılmak veya manipüle etmek** (asıl yıkıcı eylem).
- **T1486** — Data Encrypted for Impact (**fidye yazılımının imza technique'i** — veriyi şifreleyip kilitleme)
- **T1485** — Data Destruction (veriyi geri dönüşsüz silme — wiper)
- **T1490** — Inhibit System Recovery (yedek/shadow copy silme — kurtarmayı engelleme)
- **T1489** — Service Stop
- **T1498/T1499** — Network/Endpoint Denial of Service (DDoS)
- **T1491** — Defacement (web sitesi tahrifi — **hacktivist imzası**)

> 🔥 **Aktör türü ipucu burada en güçlüdür:** **T1486 (Data Encrypted for Impact)** → fidye; **T1485/T1561 (Data Destruction/Disk Wipe)** → yıkıcı/sabotaj (bazı devlet operasyonları); **T1491 (Defacement)** + **T1498 (DDoS)** → hacktivist. Impact tactic'i, "bu saldırının *amacı* neydi?" sorusunun en doğrudan cevabıdır (Bölüm 10).

---

<a id="4"></a>
## 4. 🧩 Technique vs Sub-technique vs Procedure — Granülarite

ATT&CK'in gücü, davranışı **üç granülarite seviyesinde** ifade edebilmesidir. Bu hiyerarşiyi karıştırmak en yaygın acemi hatasıdır.

```
   TACTIC (amaç)            ──►  Execution (TA0002)
        │
   TECHNIQUE (genel yöntem) ──►  T1059  Command and Scripting Interpreter
        │
   SUB-TECHNIQUE (özel yol)  ──►  T1059.001  PowerShell
        │                        T1059.003  Windows Command Shell
        │                        T1059.004  Unix Shell  ...
        │
   PROCEDURE (somut uygulama)──►  "APT29, base64-encode edilmiş bir PowerShell
                                   indirici çalıştırdı"  (gözlemlenen GERÇEK olay)
```

| Seviye | Tanım | Örnek | Numaralandırma |
|---|---|---|---|
| **Tactic** | Saldırganın amacı (NEDEN) | Execution | `TA00xx` |
| **Technique** | Amacı gerçekleştirmenin genel yolu (NASIL) | Command and Scripting Interpreter | `Txxxx` (4 hane) |
| **Sub-technique** | Technique'in daha özel bir varyantı | PowerShell | `Txxxx.0yy` (.001, .002…) |
| **Procedure** | Belirli bir aktörün **somut, gözlenmiş** uygulaması | "Grup X şu komutu şöyle çalıştırdı" | (serbest metin / "Procedure Examples") |

### Ayrımı netleştir

- **Technique (T1566 Phishing):** "Kimlik avı yaptı" — genel.
- **Sub-technique (T1566.001 Spearphishing Attachment):** "Kimlik avını *ekli zararlı dosya* yöntemiyle yaptı" — özel varyant. (Diğerleri: `.002` Spearphishing Link, `.003` Spearphishing via Service.)
- **Procedure:** "Aktör, bir CV.docx eki gönderdi; ek açıldığında makro bir HTA indirdi" — **gerçekte ne olduğu.** Procedure'lar technique sayfasındaki **"Procedure Examples"** bölümünde, hangi grubun/yazılımın o technique'i nasıl kullandığıyla listelenir (Bölüm 6).

> 🧠 **CTI için neden önemli?** **Granülarite seviyesi, raporun türünü belirler** (Bölüm 12):
> - Stratejik rapor: çoğunlukla **tactic** seviyesi ("saldırgan verinizi şifreliyor") — ID neredeyse yok.
> - Operasyonel rapor: **technique** seviyesi (T1486, T1059).
> - Taktik rapor: **sub-technique + procedure** seviyesi (T1059.001 + "şu komutu çalıştırdı").
>
> Yanlış seviyede konuşmak (CISO'ya sub-technique procedure'u anlatmak ya da avcıya sadece "kötü amaçlı kod çalıştırdı" demek) iletişimi öldürür.

> 💡 **Pratik kural:** Mümkün olan **en spesifik** seviyede mapleyebildiğin kadar maple, ama emin değilsen **bir üst seviyede kal.** "PowerShell olduğundan eminim" → T1059.001. "Script çalıştırdı ama hangisi belli değil" → sadece T1059 (sub-technique uydurma!). Yanlış spesifiklik, yanlış ID kadar zararlıdır.

---

<a id="5"></a>
## 5. 🗂️ Enterprise / Mobile / ICS Matrisleri

ATT&CK tek bir matris değil, **teknoloji alanına göre ayrı matrislerdir.** Doğru matrisi seçmek, doğru technique'i bulmanın ön koşuludur.

| Matris | Kapsam | Tactic sayısı (yaklaşık) | Ne zaman kullan |
|---|---|---|---|
| **Enterprise** | Windows, macOS, Linux, **Cloud** (IaaS/SaaS/Office 365/Azure AD/Google Workspace), Network, Containers | **14** | Kurumsal IT, sunucu, uç nokta, bulut — en yaygın, bu rehberin odağı |
| **Mobile** | iOS, Android | ~12 (farklı set) | Mobil cihaz tehditleri (mobil casus yazılım, kötü amaçlı uygulama) |
| **ICS** | Endüstriyel Kontrol Sistemleri (SCADA, PLC, OT) | ~12 (kendine özgü) | Enerji, üretim, kritik altyapı, OT ağları |

### Önemli ayrıntılar

- **Enterprise içinde alt-platformlar var:** Bir technique birden çok platforma (Windows + Linux + macOS + Cloud) uygulanabilir; technique sayfası **"Platforms"** alanında bunu belirtir. Cloud, ayrı bir matris değil, Enterprise'ın bir **alt-platformudur** (bunu sık karıştırırlar).
- **Mobile matrisi farklı tactic'ler içerir** (örn. mobil özelinde farklı initial access yolları); Enterprise technique'lerini körü körüne mobile'a taşıma.
- **ICS matrisi tamamen ayrı bir dünya** (technique'leri fiziksel süreç manipülasyonuna kadar gider — "Damage to Property", "Loss of Safety"). Bir OT incident'ini Enterprise matrisine zorlamak yanlış sonuç verir.

> 🧠 **Pratik:** %90 işin **Enterprise** matrisindedir. Mobil bir tehdit (örn. devlet-destekli mobil casus yazılım — Pegasus sınıfı) ya da bir OT/SCADA olayı analiz ediyorsan, **bilinçli olarak** Mobile/ICS matrisine geç. Yanlış matris = yanlış technique ID = çürük rapor.

---

<a id="6"></a>
## 6. 🔥 Bir Technique Sayfasını OKUMA (T1059 Baştan Sona)

Bir technique sayfasını profesyonelce okumak, ATT&CK ustalığının kalbidir. Örnek olarak **T1059 — Command and Scripting Interpreter** sayfasını baştan sona gezelim (her technique sayfası aynı şablonu kullanır).

> Adres deseni: `attack.mitre.org/techniques/T1059/` (sub-technique: `/T1059/001/`).

### Bir technique sayfasının alanları

#### 1. Başlık & ID
```
T1059  —  Command and Scripting Interpreter
Tactic: Execution (TA0002)
```
- **ID (T1059):** Evrensel referans. Raporda bu ID'yi kullanırsın.
- **Tactic bağlamı:** Bu technique hangi amaca hizmet eder (burada Execution). Bir technique **birden çok tactic** altında listelenebilir.

#### 2. Açıklama (Description)
Technique'in *ne olduğunu* düz dille anlatır: "Saldırganlar, komut/script yorumlayıcılarını (PowerShell, cmd, Bash, Python, JavaScript…) kötü amaçlı komut çalıştırmak için kötüye kullanır." Buradan technique'in **kapsamını** (neyin bu technique'e girip neyin girmediğini) öğrenirsin.

#### 3. Sub-techniques (Alt teknikler)
T1059'un alt teknikleri (numaralar attack.mitre.org'dan teyit edilebilir):
```
T1059.001  PowerShell
T1059.002  AppleScript
T1059.003  Windows Command Shell (cmd)
T1059.004  Unix Shell (bash/sh)
T1059.005  Visual Basic
T1059.006  Python
T1059.007  JavaScript
T1059.008  Network Device CLI
...
```
> 🧠 Mümkünse **sub-technique** seviyesinde maple. "PowerShell gördüm" → T1059.001, sadece T1059 değil. Spesifiklik, raporun tespit/savunma değerini artırır.

#### 4. 🔥 Procedure Examples (Prosedür Örnekleri) — CTI'nın altın madeni
Bu bölüm, **hangi grubun/yazılımın bu technique'i SOMUT olarak nasıl kullandığını** listeler. Örnek satırlar (illüstratif):
```
G0007  APT28         "...PowerShell ile indirici çalıştırdı..."
G0016  APT29         "...base64-encode edilmiş PowerShell komutları kullandı..."
S0002  Mimikatz      "...PowerShell üzerinden çalıştırıldı..."
```
> 💡 **CTI'da nasıl kullanılır?** Procedure Examples, "**bu technique'i kimler kullanıyor?**" sorusunun cevabıdır → aktör profillemenin (Bölüm 10) temel girdisi. Tersi de mümkün: bir grup sayfasından (Bölüm 7) o grubun technique listesine, oradan procedure'lara inersin. İlişkiler **çift yönlüdür.**

#### 5. Mitigations (Önlemler)
Bu technique'i **önlemenin/zorlaştırmanın** yolları (her birinin Mxxxx ID'si vardır):
```
M1042  Disable or Remove Feature or Program  (gereksiz yorumlayıcıları kaldır)
M1038  Execution Prevention                  (uygulama beyaz listesi/AppLocker)
M1026  Privileged Account Management
M1049  Antivirus/Antimalware
...
```
> 🔗 **Rapor değeri:** Raporunun "Öneriler/Mitigation" bölümünü doğrudan buradan beslersin (Bölüm 12). "T1059.001'e karşı M1038 Execution Prevention uygulayın" — somut, ATT&CK-bağlantılı öneri.

#### 6. Detection (Tespit)
Technique'i **yakalamanın** yolları — neye, hangi veride bakmalısın. Modern ATT&CK bunu **Data Sources / Data Components** ile yapılandırır:
```
Data Source: Command          → Command Execution
Data Source: Process          → Process Creation
Data Source: Script           → Script Execution
Data Source: Module           → Module Load
```
> 💡 **Avcı/SOC değeri:** "T1059.001'i yakalamak için **PowerShell Script Block Logging** (Event ID 4104) ve **Process Creation** (4688) loglarına bak." Detection alanı, tespit mühendisliğinin başlangıç noktasıdır (→ derin avcılık `TTP_AVCILIGI_USTALIK_REHBERI.md`).

#### 7. Data Sources (Veri Kaynakları)
Technique'i tespit etmek için **hangi telemetriyi toplaman** gerektiği. Bu, savunma **boşluk analizinin** temelidir: "T1059'u tespit için Process Creation logu lazım; bende bu var mı?" Yoksa, o technique'e karşı **körsün** (Navigator heatmap'inde bunu görselleştirirsin — Bölüm 9).

#### 8. İlişkili alanlar
- **Platforms:** Windows, Linux, macOS, Network… (hangi platformda geçerli).
- **Permissions Required / Effective Permissions:** Gereken yetki seviyesi.
- **References:** Bu technique'in belgelendiği gerçek tehdit raporları (birincil kaynaklar — atfetme için değerli).

### 🔥 Bir technique'i okurken analist soru listesi

> Her technique sayfasında kendine sor:
> 1. **Bu technique tam olarak neyi kapsıyor?** (Description — sınırları)
> 2. **Hangi sub-technique benim gözlemimle eşleşiyor?** (en spesifik doğru seviye)
> 3. **Kimler kullanıyor?** (Procedure Examples → aktör hipotezi)
> 4. **Nasıl tespit ederim?** (Detection + Data Sources → bende o veri var mı?)
> 5. **Nasıl önlerim/azaltırım?** (Mitigations → rapor önerisi)
> 6. **Bu technique kaç tactic'e hizmet ediyor?** (bağlam — neden kullanılmış?)

> 🧠 **T1566 (Phishing) için aynı egzersizi yap:** Sub-technique'leri (`.001` Attachment, `.002` Link, `.003` via Service), Procedure Examples'ı (hangi APT hangi phishing'i kullanıyor), Mitigations'ı (M1017 User Training, M1031 Network Intrusion Prevention) ve Detection'ı (e-posta gateway, ek analizi) oku. İki farklı technique'i baştan sona gezdiğinde **şablonu içselleştirmiş** olursun — gerisi tekrardır.

---

<a id="7"></a>
## 7. 👥 Groups (G-numaraları) & Software (S-numaraları) Sayfaları

ATT&CK yalnızca technique'leri değil, **kim** (Groups) ve **neyle** (Software) sorularını da kataloglar. Bunlar CTI'nın "aktör" boyutudur.

### Groups (Tehdit Grupları) — G-numaraları

Bir **Group** sayfası, MITRE'nin **birden çok kamuya açık rapora** dayanarak izlediği bir tehdit aktörünü (genelde bir APT veya organize suç grubu) belgeler. Her grubun bir **Gxxxx** ID'si vardır.

> 🧠 **Kritik nüans — isimler kaos, ID'ler düzen:** Aynı aktör, farklı vendor'larda farklı isimlerle anılır (örn. APT28 = Fancy Bear = Sofacy = STRONTIUM/Forest Blizzard = Pawn Storm…). MITRE bunu **tek bir G-numarasında** birleştirir ve **"Associated Groups"** alanında tüm takma adları listeler. **Raporda aktöre G-numarasıyla atıfta bulunmak**, isim kargaşasını çözer ("APT28 (G0007)").

#### Bir grup sayfası nasıl okunur — APT28 (G0007) örneği

> APT28 = Fancy Bear, Rusya GRU ile ilişkilendirilen, casusluk odaklı devlet-destekli bir aktör. (Atıf iddiaları kamuya açık raporlardandır; G-numarasını ve güncel technique listesini attack.mitre.org/groups/G0007'den teyit et.)

Bir grup sayfasında şu alanlar olur:
```
G0007  APT28
  Associated Groups: Fancy Bear, Sofacy, STRONTIUM/Forest Blizzard, Pawn Storm, Sednit ...
  Description:       Kim oldukları, atıf, hedef sektörler, motivasyon (casusluk)
  Techniques Used:   Kullandıkları technique'ler — HER BİRİ bir procedure açıklamasıyla
                       T1566.001  Spearphishing Attachment  "...şu kampanyada..."
                       T1059.001  PowerShell                 "...şöyle kullandı..."
                       T1003      OS Credential Dumping      "..."
                       ...
  Software Used:     Kullandıkları araçlar/malware (S-numaralarına link)
                       S0002 Mimikatz, S0xxx X-Agent, ...
  References:        Bu profili oluşturan birincil tehdit raporları
```

> 🔥 **CTI'da grup sayfasının kullanımı:**
> - **Aktör → TTP seti:** Bir grubun **"Techniques Used"** listesi, o aktörün **TTP imzasıdır** → Navigator'da layer'a dökülür (Bölüm 9), aktör karşılaştırmasında kullanılır (Bölüm 10).
> - **Aktör → araçlar:** "Software Used", aktörün araç setini verir → IOC zenginleştirme.
> - **Çift yönlü gezinme:** Technique sayfasından (Procedure Examples) gruba, grup sayfasından technique'e — ilişki ağında **pivot** yaparsın (Bölüm 8).
> - **References = birincil kaynak:** Atfetme yaparken MITRE'nin özetine değil, **References'taki orijinal raporlara** in (doğrulama disiplini — Bölüm 10, atfetme tehlikeleri).

### Software (Yazılım/Araçlar) — S-numaraları

Bir **Software** sayfası, tehdit aktörlerinin kullandığı bir **malware** veya **meşru-ama-kötüye-kullanılan araç**ı (tool) belgeler. Her birinin **Sxxxx** ID'si vardır.

- İki tip: **Malware** (saldırgana özel kötü amaçlı yazılım) ve **Tool** (meşru ama kötüye kullanılan — örn. PsExec, Mimikatz, Cobalt Strike).
- Bir software sayfası şunları içerir: hangi **technique'leri** uyguladığı (yani o aracı kullanmak hangi TTP'leri ima eder), hangi **gruplar** tarafından kullanıldığı, platformları.

```
S0002  Mimikatz  (Tool)
  Techniques:  T1003.001 LSASS Memory, T1550.002 Pass the Hash, T1558 Kerberos ...
  Groups:      G0007 APT28, G00xx ..., (Mimikatz'i kullanan tüm gruplar)
```

> 🧠 **Neden değerli?** Bir incident'te **Mimikatz (S0002)** gördüysen, bu otomatik olarak bir **technique kümesini** (T1003 credential dumping vb.) ve **olası aktör havuzunu** (Mimikatz'i kullanan gruplar) ima eder. Software, technique ile grubu birbirine bağlayan köprüdür. **Ama dikkat:** Mimikatz/Cobalt Strike gibi **herkesin kullandığı** araçlar **zayıf atıf sinyalidir** (Bölüm 10 — TTP overlap tuzağı). Aracın varlığı technique'i kanıtlar, aktörü kanıtlamaz.

### Campaigns (Kampanyalar) — C-numaraları (bonus)

Yeni ATT&CK sürümleri **Campaigns** (Cxxxx) de ekledi: belirli bir zaman aralığında, belirli bir hedefe yönelik, ilişkilendirilmiş bir saldırı etkinliği kümesi. Bir kampanya sayfası, o kampanyada kullanılan technique/software/group ilişkilerini bir arada verir → tarihsel kampanya analizinde değerli. (Varlığını ve içeriğini attack.mitre.org'dan teyit et.)

---

<a id="8"></a>
## 8. 🔎 attack.mitre.org'da Arama & İlişkilerde Gezinme

ATT&CK web sitesini akıcı kullanmak, bir analistin günlük refleksidir. Site bir **ilişki grafiğidir** (technique ↔ group ↔ software ↔ mitigation ↔ data source) ve asıl güç **gezinmededir.**

### Arama yöntemleri

| Aradığın | Nasıl bulursun |
|---|---|
| **Belirli ID** | URL'e doğrudan git: `attack.mitre.org/techniques/T1059/` · `/groups/G0007/` · `/software/S0002/` · `/mitigations/M1038/` |
| **Technique (isimle)** | Üstteki arama kutusuna "PowerShell" / "phishing" / "credential dumping" yaz |
| **Group (isim/takma ad)** | "Fancy Bear" ara → G0007'ye yönlendirir (Associated Groups eşlemesi) |
| **Software** | "Mimikatz" / "Cobalt Strike" ara |
| **Mitigation** | Mitigations sekmesinden M-numaraları |
| **Data Source** | Data Sources sekmesi → hangi technique'leri hangi veri yakalar |
| **Tactic** | Matriste sütun başlığına tıkla → o tactic'in tüm technique'leri |

### 🔥 İlişkilerde gezinme (pivot) — asıl ustalık

ATT&CK'i "arama" değil **"gezinme"** aracı olarak kullan. Tipik pivot zincirleri:

```
   Bir technique buldun (T1003)
        │
        ├──► "Procedure Examples" → hangi GRUPLAR kullanıyor? (G0007, G00xx...)
        │         │
        │         └──► O grubun sayfası → başka hangi technique'leri var? (aktör TTP seti)
        │
        ├──► "Mitigations" → nasıl önlerim? (M1026, M1027...)
        │
        ├──► "Detection / Data Sources" → nasıl yakalarım? (Process, Command...)
        │
        └──► İlgili "Software" → bu technique'i hangi araçlar uyguluyor? (S0002...)
                  │
                  └──► O software'i hangi gruplar kullanıyor? → yeni aktör havuzu
```

> 🔥 **Pivot örneği (uçtan uca):** Bir incident'te **LSASS dump** gördün → T1003.001 sayfasına git → Procedure Examples'ta bunu kullanan grupları gör → birinin sayfasına gir → o grubun *diğer* technique'lerini al → bunları kendi gözlemlerinle karşılaştır ("benim gördüğüm 5 technique'in 4'ü bu grupta var mı?") → aktör hipotezi kur (Bölüm 10). **Tek bir gözlemden, ilişki grafiğinde gezinerek bir aktör profiline ulaşmak** — ATT&CK'in asıl gücü budur.

> 💡 **STIX / programatik erişim:** ATT&CK verisinin tamamı **STIX 2.x JSON** formatında GitHub'da (`mitre-attack/attack-stix-data`) ve bir **TAXII** sunucusu üzerinden makine-okunur olarak yayınlanır. Büyük ölçekli analiz, otomasyon veya kendi araçlarına ATT&CK gömme için bu veri kümesini çek. `python-attackcti` / `mitreattack-python` kütüphaneleri bunu kolaylaştırır. (MISP, ATT&CK galaksisini bu veriden besler → Bölüm 19.)

---

<a id="9"></a>
## 9. 🔥 ATT&CK Navigator — Sıfırdan Uzmana

**ATT&CK Navigator**, matrisi **görselleştiren, renklendiren, skorlayan ve katmanlayan** ücretsiz, açık kaynak web aracıdır. CTI'da Navigator = "ATT&CK'i bir resme dönüştürme" aracıdır: aktör TTP'lerini boyamak, savunma kapsamını haritalamak, iki aktörü karşılaştırmak, heatmap üretmek için kullanılır.

> Erişim: `mitre-attack.github.io/attack-navigator/` (tarayıcıdan, kurulumsuz) ya da kendi sunucunda barındırabilirsin (Docker). Veri yerelde işlenir.

### Temel kavram: Layer (Katman)

Navigator'da her şey bir **layer**'dır: matrisin üstüne bindirilmiş, technique'lere **renk + skor + yorum** atayan bir kaplama. Bir layer **JSON** olarak export/import edilir → paylaşılabilir, versiyonlanabilir, birleştirilebilir.

### 🔥 Sıfırdan adım adım: ilk layer'ını oluştur

1. **Aç:** `mitre-attack.github.io/attack-navigator/` → **"Create New Layer"** → **Enterprise** (matris seç).
2. **Technique seç:** Bir technique kutusuna tıkla → seçilir (vurgulanır). Çoklu seçim: Ctrl+tık ya da arama ile filtrele.
3. **Renklendir (color):** Seçili technique'lere sağ panelden/araç çubuğundan **arka plan rengi** ata. (Örn. "bu aktörün kullandığı technique'ler" = kırmızı.)
4. **Skor ver (scoring):** Her technique'e bir **sayısal skor** ata (örn. gözlem sıklığı, güven, önem). Navigator skorları **renk gradyanına** (heatmap) çevirir.
5. **Yorum (comment):** Bir technique'e not ekle ("bu kampanyada şu tarihte görüldü, kaynak: X raporu").
6. **Meta:** Layer'a isim/açıklama ver (örn. "APT28 — gözlemlenen TTP'ler v1").
7. **Export:** Araç çubuğundan **JSON indir** (paylaşmak/saklamak için) veya **SVG/Excel** (rapora gömmek için).

### 🔥 Kullanım 1 — Bir aktörün TTP'lerini boyamak

Bir grubun TTP imzasını görselleştir:
- ATT&CK sitesinden grubun (örn. G0007 APT28) **"Techniques Used"** listesini al.
- Navigator'da yeni layer aç, bu technique'leri seç, **tek renge** boya, isim ver.
- Sonuç: APT28'in matristeki "ayak izi" — hangi tactic'lerde yoğun, nerede zayıf, görsel olarak okunur.

> 💡 **Kısayol:** ATT&CK grup sayfalarının çoğunun altında **"ATT&CK Navigator Layers"** indirme bağlantısı vardır — o grubun TTP'leri **hazır bir layer JSON'u** olarak. İndir → Navigator'da "Open Existing Layer" → import et. Sıfırdan boyamana gerek kalmaz (ama kendi gözlemini eklemek için elle düzenlersin).

### 🔥🔥 Kullanım 2 — İki aktörü KARŞILAŞTIRMA (layer kesiştirme/birleştirme)

Bu, Navigator'ın en güçlü CTI özelliğidir. İki aktörün TTP'lerini matematiksel olarak karşılaştırırsın:

1. **İki layer hazırla:** Layer A = Aktör 1 TTP'leri (her technique skoru = 1), Layer B = Aktör 2 TTP'leri (skoru = 1).
2. **Yeni "score expression" layer'ı oluştur:** Navigator'da **"Create Layer from other layers"** seç.
3. **Matematiksel ifade gir:** Layer'lara bir harf atanır (a, b). Sonra bir **skor ifadesi** yazarsın:
   - **`a + b`** → her iki aktörde de olan technique'ler **yüksek skor** alır (kesişim vurgusu — *ortak* TTP'ler).
   - **`a - b`** → yalnızca Aktör 1'de olanlar pozitif, ortak olanlar 0 (fark vurgusu — *ayırt edici* TTP'ler).
   - **`a * b`** → yalnızca **her ikisinde de** olanlar 1, diğerleri 0 (saf **kesişim**).
4. **Gradyan ata** → ortak technique'ler bir renkte, ayrışanlar başka renkte parlar.

```
   Layer A (APT-X TTP'leri)        Layer B (APT-Y TTP'leri)
        │ skor=1 her technique          │ skor=1 her technique
        └───────────────┬───────────────┘
                        ▼
            "Create from other layers"
            ifade:  a + b   (ya da a*b kesişim, a-b fark)
                        ▼
        ┌──────────────────────────────────────┐
        │  Birleşik heatmap:                    │
        │   skor 2 = İKİSİNDE de var (ortak TTP)│ ← koyu kırmızı
        │   skor 1 = yalnız birinde (ayırt edici)│ ← açık
        │   skor 0 = hiçbirinde                 │ ← boş
        └──────────────────────────────────────┘
```

> 🔥 **CTI değeri:** İki şüpheli kampanyanın **aynı aktör mü** olduğunu sorgularken: TTP'lerini `a*b` ile kesiştir → **çok ortak ve özellikle NADİR technique'lerde örtüşüyorlarsa**, aynı aktör hipotezi güçlenir. `a-b` ile her birinin **ayırt edici** technique'lerini bul. Bu, atfetmenin görsel altyapısıdır (ama tek başına kanıt değil — Bölüm 10, atfetme tehlikeleri).

### 🔥🔥 Kullanım 3 — Savunma KAPSAMA HARİTASI (detection coverage)

Navigator'ın blue-team için en değerli kullanımı:
1. **Layer 1 — Tehdit:** Sektörünü hedefleyen aktörlerin TTP'lerini boya (örn. seni hedefleyen 3 aktörün birleşimi, `a+b+c`). Bu, **"bana karşı kullanılan technique'ler"**dir.
2. **Layer 2 — Savunma:** Mevcut tespit yeteneğin olan technique'leri boya (hangi technique için log/kural/EDR kapsaman var). Bu, **"yakalayabildiklerim"**dir.
3. **Çıkar:** `tehdit - savunma` → **pozitif skorlu technique'ler = KÖR NOKTALARIN** (sana karşı kullanılıyor ama yakalayamıyorsun). Heatmap'te kırmızı parlayanlar, **savunma yatırımının önceliğidir.**

> 🧠 **Bu, "ATT&CK ile risk önceliklendirme"nin özüdür:** Sınırlı savunma bütçeni, "modaya uygun" technique'lere değil, **sana karşı gerçekten kullanılan + senin kör olduğun** technique'lere harca. Navigator bu boşluğu bir resimde gösterir. (Operasyonel raporun (Bölüm 12) en güçlü eki budur.)

### Heatmap üretme & skorlama desenleri

- **Sıklık heatmap'i:** Bir kampanyada her technique kaç kez gözlendi → skor = sıklık → "en çok kullanılan TTP'ler" ısı haritası.
- **Güven heatmap'i:** Her technique için atıf/gözlem güveni (1-100) → düşük güvenli gözlemler soluk, yüksek güvenli parlak.
- **Çoklu kampanya:** Her kampanya bir layer, hepsini topla (`a+b+c+...`) → "bu aktör zaman içinde hangi TTP'lerde *tutarlı*" (en stabil = imza).

### Export / Import / Paylaşım

- **JSON export:** Layer'ı paylaş, versiyonla, rapora ek yap, başka analistle birleştir. ATT&CK layer JSON'u standart bir formattır (sürüm alanı içerir — `versions`).
- **SVG/Excel export:** Rapora **görsel** gömmek için (Bölüm 12 — Appendix'te Navigator layer'ı).
- **Import:** Başkasının (ya da ATT&CK sitesinin hazır) layer'ını yükle, üstüne çalış.

> 💡 **Pratik püf:** Her CTI raporuna, kampanyanın TTP'lerini gösteren **bir Navigator layer'ı (JSON + SVG)** ekle. Okuyucu (özellikle savunma ekibi) bunu kendi kapsama layer'ıyla kesiştirip **kendi kör noktalarını** anında çıkarır. Bu, raporunu "okunan" değil "kullanılan" yapar.

---

<a id="10"></a>
## 10. 🔥🔥 MANUEL TEHDİT AKTÖRÜ TÜRÜ BELİRLEME (Profilleme Metodolojisi)

> Bu, rehberin **en kritik ve en çok talep edilen** bölümüdür. Soru şu: *"Elimde gözlemlenmiş davranışlar (TTP'ler) var — bunlardan, karşımdaki aktörün KİM ve özellikle hangi TÜR olduğunu nasıl çıkarırım?"* ATT&CK bu işin omurgasıdır, ama **acımasız tuzaklarla** doludur. Önce metodoloji, sonra aktör türleri, sonra tehlikeler.

> ⚠️ **Peşinen uyarı — atfetme zordur:** Profesyonel atıf (özellikle "hangi spesifik grup/ülke") devasa kaynak, çoklu istihbarat kaynağı ve aylar ister; çoğu zaman ulus-devlet/üst düzey vendor işidir. Bir analist olarak senin gerçekçi ve değerli hedefin çoğunlukla **aktör TÜRÜNÜ** (devlet-destekli mi? fidye mi? hacktivist mi? insider mı?) ve **muhtemel aktör ailesini** belirlemektir — "kesin Rusya" demek değil. Bu ayrımı baştan kabul et; aşırı-atıf en yaygın ve en zararlı hatadır.

### 🔥 Adım adım profilleme metodolojisi

```
   1. GÖZLEM           Logları/IOC'leri/olay anlatısını topla (ham davranış)
        │                "kullanıcıya makro Word geldi", "LSASS dump", "veri şifrelendi"
        ▼
   2. TECHNIQUE MAP    Her gözlemi ATT&CK technique'ine çevir (Bölüm 11 antrenmanı)
        │                → T1566.001, T1003.001, T1486 ...
        ▼
   3. TTP SETİ ÇIKAR   Gözlemlenen technique'lerin bütününü bir SET olarak yaz
        │                {T1566.001, T1059.001, T1053.005, T1003.001, T1486, ...}
        ▼
   4. AKTÖR KARŞILAŞTIR ATT&CK Groups'u bu sete göre tara (Navigator a*b kesişim)
        │                "hangi grupların TTP seti benimkiyle örtüşüyor?"
        ▼
   5. TÜR AYRIMI       TTP imzasından aktör TÜRÜNÜ çıkar (aşağıdaki tablo)
        │                Impact=T1486 → fidye? Uzun persistence+stealth → APT?
        ▼
   6. HİPOTEZ + GÜVEN  Rakip hipotezlerle (ACH) tart, dereceli güvenle ifade et
                         "orta güvenle fidye operasyonu, X ailesiyle tutarlı"
```

#### Adım 1-3: Gözlem → Technique → TTP seti
Olayı **çıplak davranışlara** ayır, her birini technique'e mapele (bunun pratiği Bölüm 11'de), gözlemlenen technique'leri bir **küme** olarak yaz. Bu küme, aktörün **TTP imzasıdır.**

#### Adım 4: ATT&CK Groups'u TTP setine göre karşılaştırma
- Her gözlemlenen technique'in **"Procedure Examples"**ına bak → o technique'i kullanan grupları topla.
- Birden çok technique'inde **tekrar tekrar aynı grup** çıkıyorsa, o grup bir hipotez adayıdır.
- **Navigator'da yap (Bölüm 9):** Senin gözlem layer'ın × aday grup layer'ı → kesişim (`a*b`) → ne kadar örtüştüğünü gör. Özellikle **NADİR technique'lerde** örtüşme değerlidir (herkesin yaptığı PowerShell örtüşmesi zayıf sinyal; nadir bir custom persistence örtüşmesi güçlü sinyal).

#### Adım 5: Aktör TÜRÜNÜ TTP imzasından çıkarma (en pratik kısım)

Spesifik grubu bulamasan bile, TTP imzası genelde **aktörün türünü** ele verir. Her türün karakteristik bir davranış profili vardır:

| Aktör türü | Motivasyon | Karakteristik TTP imzası | Tipik "tonu" |
|---|---|---|---|
| **Devlet-destekli APT** | Casusluk, istihbarat, uzun vadeli erişim | **Uzun persistence** (T1547, T1053, custom backdoor), **derin stealth** (T1070 Indicator Removal, T1027 Obfuscation, T1036 Masquerading, living-off-the-land), hedefli **Collection** (T1113, T1005) + **Exfiltration** (T1041), düşük gürültü, **özel araçlar** (custom malware S-numaraları), hedefe özel spearphishing (T1566.001) | Sabırlı, sessiz, hedefe özel, "iz bırakmadan kal ve dinle" |
| **Siber suç / Fidye (ransomware)** | Para (fidye, double extortion) | **Impact: T1486 (Data Encrypted for Impact)** + **T1490 (Inhibit System Recovery — shadow copy silme)** + **T1489 (Service Stop)**, hızlı **lateral movement** (T1021), credential dumping (T1003), sık **Exfiltration önce** (double extortion), **commodity araçlar** (Cobalt Strike, Mimikatz, PsExec), affiliate modeli | Hızlı, gürültülü, "şifrele-tehdit-tahsil et", vur-kaç |
| **Hacktivist** | İdeoloji, mesaj, itibar zedeleme | **T1491 (Defacement)**, **T1498/T1499 (DoS/DDoS)**, veri sızdırıp **kamuya açıklama** (utandırma), düşük teknik sofistikasyon (genelde), açık üstlenme/propaganda | Gösterişçi, kamuya dönük, "mesaj ver ve duyur" |
| **Insider (içeriden tehdit)** | İntikam, para, casusluk | **Valid Accounts (T1078)** ile *meşru* erişim (initial access yok — zaten içeride), **Collection** + **Exfiltration**, az/hiç malware, normal araçlarla anormal davranış, yetki kötüye kullanımı | Düşük teknik iz, "zaten anahtarı var", anomali davranışta |
| **Script kiddie / fırsatçı** | Merak, kolay hedef, düşük beceri | Hazır exploit/araç, gürültülü tarama (T1595), tutarsız/dağınık TTP, otomatik kitler, hedef ayrımı yok | Dağınık, otomatik, "ne tutarsa" |

> 🔥 **En keskin tek ipucu — Impact tactic'i:** Aktörün **niyeti** çoğu zaman en net Impact (TA0040) tactic'inde okunur:
> - **T1486 Data Encrypted for Impact** → **fidye** (neredeyse kesin işaret).
> - **T1485 Data Destruction / T1561 Disk Wipe** → **yıkıcı/sabotaj** (bazı devlet operasyonları, "wiper" — örn. NotPetya sınıfı).
> - **T1491 Defacement** + **T1498 DDoS** → **hacktivist**.
> - **Hiç Impact yok, sadece sessiz Collection/Exfiltration** → **casusluk (APT)**.
>
> Impact yoksa "niyet = bilgi çalmak/kalmak"; Impact varsa, *hangi* Impact technique'i aktör türünü neredeyse söyler.

> 🧠 **Diğer ayırt edici sinyaller:**
> - **Persistence derinliği:** Çok katmanlı, özel persistence → APT. Persistence umursamayan, hızlı şifreleyen → fidye.
> - **Stealth yatırımı:** Yoğun T1070/T1027/T1036 → tespitten kaçınma önemli → APT/sofistike. Gürültülü, umursamaz → fidye/kiddie.
> - **Araç tipi:** Custom malware (özel S-numaraları) → kaynaklı/APT. Sadece commodity (Cobalt Strike, Mimikatz, açık kaynak) → fidye/affiliate/kiddie (ama APT'ler de bilinçli commodity kullanıp **karışmaya** çalışır — Bölüm aşağıda false flag).
> - **Hedefleme:** Çok hedefli spearphishing (belirli kişiler) → APT. Geniş, ayrımsız → fidye/kiddie/fırsatçı.
> - **Initial Access yokluğu:** Hiç initial access technique'i yok, sadece Valid Accounts → **insider** ihtimali.

#### Adım 6: Hipotez + dereceli güven

Asla "kesin X" deme. CTI'da atıf **dereceli güven** ve **rakip hipotez analizi (ACH — Analysis of Competing Hypotheses)** ile ifade edilir.

### 🎲 Diamond Model — atfetmenin dört köşesi

Aktör profillemede ATT&CK'i tamamlayan klasik çerçeve **Diamond Model**'dir. Her saldırı olayını dört köşeli bir elmasla modeller:

```
                    ADVERSARY (Aktör)
                   "kim?" (atıf, motivasyon)
                          ▲
                          │
       INFRASTRUCTURE ◄───┼───► CAPABILITY
       "neyle bağlandı?"  │     "hangi araç/TTP?"
       (C2, domain, IP)   │     (malware, ATT&CK TTP'leri)
                          │
                          ▼
                      VICTIM (Kurban)
                  "kime?" (hedef, sektör)
```

| Köşe | Soru | ATT&CK / OSINT bağı |
|---|---|---|
| **Adversary** | Kim? Motivasyon? | Groups (G-numaraları), atıf — bu bölüm |
| **Capability** | Hangi yetenek/araç/TTP? | **ATT&CK technique'leri + Software (S-numaraları)** |
| **Infrastructure** | Hangi altyapı? | C2/domain/IP — OSINT pivotu (`OSINT_ARAC_SETI...`), pDNS, JARM |
| **Victim** | Kim hedef alındı? | Sektör, coğrafya, kişi — hedefleme deseni (tür ipucu) |

> 🧠 **ATT&CK ↔ Diamond ilişkisi:** ATT&CK, Diamond'ın **Capability** (ve kısmen Adversary) köşesini doldurur — "aktörün ne yapabildiği/yaptığı". Infrastructure köşesi OSINT'le, Victim köşesi hedefleme analizi ile dolar. **Dört köşeyi birlikte** kullandığında atıf güçlenir; tek köşeye (sadece TTP'ye) dayanan atıf **kırılgandır.** (Diamond Model'in derin teorisi → `TTP_AVCILIGI_USTALIK_REHBERI.md`.)

### ⚠️🔥 ATFETME TEHLİKELERİ — bu kısmı atlama

Atfetme, CTI'nın en tehlikeli sularıdır. ATT&CK TTP'leriyle atıf yaparken bu tuzaklar **sürekli** seni bekler:

1. **🚩 False flag (sahte bayrak):** Sofistike aktörler **kasten yanıltır** — başka bir aktörün araçlarını, dilini, çalışma saatlerini, hatta bilinen bir grubun imza TTP'lerini **taklit eder.** "Çok temiz, çok uygun" bir atıf delili → **şüphelen** (yerleştirilmiş olabilir). Tarihte gerçek false flag operasyonları yaşandı (bir saldırının başka bir ülkeye/gruba yıkılmaya çalışılması).
2. **🔁 TTP overlap (örtüşme):** Birçok aktör **aynı yaygın technique'leri** kullanır — herkes PowerShell (T1059.001), herkes phishing (T1566), herkes Mimikatz (S0002). **Yaygın TTP'lerde örtüşme atıf KANITI değildir.** Atıf değeri yalnızca **nadir, özel, ayırt edici** TTP'lerin örtüşmesindedir.
3. **🐑 Copy-cat / araç paylaşımı:** Bir aktörün araçları sızdığında (örn. sızdırılmış Cobalt Strike, açık kaynak olmuş malware), **başka aktörler aynısını kullanır.** Araç eşleşmesi → "aynı aktör" demek değil. Affiliate/RaaS (ransomware-as-a-service) modelinde **aynı malware'i onlarca farklı aktör** kullanır.
4. **📋 Circular reporting (döngüsel raporlama):** Bir vendor "bu APT28" der, herkes onu kopyalar, sonra "5 kaynak da APT28 diyor" sanılır — ama hepsi **tek (belki yanlış) kaynaktan** gelir. Birincil kaynağa in (References, Bölüm 7).
5. **🎯 Confirmation bias (onaylama yanlılığı):** "Bu APT28 olsun istiyorum" → APT28'i destekleyen kanıtları görür, çürüten kanıtları görmezsin. **Hipotezini çürütmeye çalış**, doğrulamaya değil (ACH disiplini).
6. **⏰ Bayat/eksik veri:** Eksik loglarla yapılan atıf, resmin yarısını görmektir. Gördüğün TTP'ler aktörün **tüm** TTP'leri değil, sadece **yakalayabildiklerin** olabilir.

### Atfetmenin güven dili (dereceli ifade)

Atıfı **asla mutlak** ifade etme. Standart güven seviyeleri:

| Güven seviyesi | Ne zaman | Örnek dil |
|---|---|---|
| **Düşük güven (low confidence)** | Az/dolaylı kanıt, alternatifler güçlü | "Gözlemlenen TTP'ler fidye operasyonuyla **tutarlı olabilir**, ancak veri sınırlı." |
| **Orta güven (medium confidence)** | Birden çok kaynak/TTP destekliyor ama kesin değil | "**Orta güvenle**, bu kampanya X ailesiyle (RaaS) tutarlıdır." |
| **Yüksek güven (high confidence)** | Çoklu bağımsız kaynak, nadir TTP örtüşmesi, güçlü altyapı bağı | "**Yüksek güvenle**, gözlemlenen TTP, altyapı ve hedefleme G00xx ile ilişkilidir." |

> 🧠 **Altın cümle kalıbı:** *"[Düşük/orta/yüksek] güvenle, gözlemlenen TTP'ler [aktör türü / aktör ailesi] ile tutarlıdır; [şu ayırt edici technique'ler] bu değerlendirmeyi destekler, [şu sınırlama] güveni kısıtlar."* Bu kalıp seni hem **savunulabilir** hem **dürüst** tutar. "Kesin X yaptı" cümlesi, bir CTI analistinin yapabileceği **en tehlikeli** beyandır (yanlış atıf → yanlış aksiyon, diplomatik/hukuki sonuç).

---

<a id="11"></a>
## 11. 🔥🔥 SALDIRI TÜRÜ BULMA ANTRENMANI (Çözümlü Egzersizler)

> Teori bitti — şimdi **pratik.** ATT&CK'te ustalık, "bir saldırı anlatısını/log'u okuyup technique zincirine mapleyebilmek"tir. Bu beceri ancak **tekrarlı egzersizle** gelir. Aşağıda zorluğu artan **çözümlü senaryolar** var. Her egzersizi önce **kendin çöz** (technique'leri yaz), sonra çözümü aç.

### 🏋️ Kendi başına nasıl pratik yaparsın (kaynaklar)

- **MITRE'nin kendi eğitimi:** MITRE ATT&CK ücretsiz **"ATT&CK for Cyber Threat Intelligence"** eğitim modülleri (attack.mitre.org → Resources → Training) — tam da "rapor oku, technique'e maple" pratiği yaptırır. **MAD (MITRE ATT&CK Defender)** sertifika programı daha derin (bazı modüller ücretsiz).
- **Atomic Red Team:** Her ATT&CK technique için **çalıştırılabilir küçük test'ler** (`atomics/`). Bir technique'in *gerçekte* nasıl göründüğünü güvenli bir lab'da çalıştırıp logunu görerek öğrenirsin — "T1059.001 çalışınca hangi log oluşur?" sorusunun en iyi cevabı. (Yalnızca **kendi izole lab ortamında** çalıştır.)
- **CTI rapor okuma alıştırması:** Yayınlanmış gerçek tehdit raporlarını (vendor blog'ları, CISA advisory'leri) al, **her cümleyi technique'e maplemeye** çalış, sonra raporun kendi ATT&CK eşlemesiyle (çoğu modern rapor verir) karşılaştır. En gerçekçi antrenman budur.
- **Kendi loglarınla:** Bir lab'da Atomic Red Team çalıştır → Sysmon/EDR loglarını topla → "bu log hangi technique?" diye tersine maple.

> 🧠 **Antrenman disiplini:** Her senaryoda 6 adımı uygula (Bölüm 10): gözlem → technique map → TTP seti → tactic sıralaması → aktör türü → güven. Sadece "T-numarası bul"ma; **tactic zincirini** (saldırının amaç akışını) ve **aktör türü hipotezini** de çıkar.

---

### 📦 EGZERSİZ 1 (Kolay) — Klasik makro-phishing zinciri

> **Senaryo:** Bir kullanıcıya, ekinde makro içeren bir Word belgesi olan bir e-posta geldi. Kullanıcı belgeyi açtı ve makroyu etkinleştirdi. Makro, base64 ile kodlanmış bir PowerShell komutu çalıştırdı. Bu komut bir zamanlanmış görev kurdu. Ardından sistem, dışarıdaki bir sunucuya 443 portu (HTTPS) üzerinden bağlandı.
>
> **Görev:** Olayı tactic sırasına göre technique'lere maple. Aktör türü hakkında ne söyleyebilirsin?

<details>
<summary><b>✅ ÇÖZÜM — Egzersiz 1</b></summary>

| # | Gözlem | Tactic | Technique |
|---|---|---|---|
| 1 | Ekli makrolu Word ile e-posta | Initial Access | **T1566.001** Spearphishing Attachment |
| 2 | Kullanıcı belgeyi açtı / makroyu etkinleştirdi | Execution | **T1204.002** Malicious File (User Execution) |
| 3 | Base64 PowerShell çalıştı | Execution | **T1059.001** PowerShell (+ **T1027** Obfuscated Files/Info — base64 kodlama) |
| 4 | Zamanlanmış görev kuruldu | Persistence | **T1053.005** Scheduled Task |
| 5 | 443'ten dış sunucuya bağlantı | Command and Control | **T1071.001** Web Protocols (HTTPS) (+ muhtemel **T1573** Encrypted Channel) |

**Tactic zinciri:** Initial Access → Execution → Persistence → C2.

**Aktör türü değerlendirmesi:** Bu **çok yaygın** bir başlangıç zinciridir — neredeyse her aktör türü kullanır (commodity malware, fidye öncesi, APT öncesi). **Tek başına aktör türünü ayırt etmez** (Impact/Collection görmeden niyet belirsiz). Düşük güvenle: "tipik bir initial access + persistence + C2 kurulumu; niyet (fidye mi casusluk mu) sonraki aşamalara bağlı." **Bu egzersizin dersi:** Erken-aşama TTP'ler aktör türü için **zayıf** sinyaldir; niyet Impact/Collection'da okunur.

*(Tüm ID'leri attack.mitre.org'dan teyit et.)*
</details>

---

### 📦 EGZERSİZ 2 (Orta) — Fidye operasyonu

> **Senaryo:** İnternete bakan bir VPN cihazındaki bilinen bir açık sömürülerek ağa girildi. Saldırgan, LSASS belleğinden parola hash'lerini boşalttı. Bu kimlik bilgileriyle RDP üzerinden 12 sunucuya yayıldı. Yayıldığı her sunucuda Windows Defender'ı devre dışı bıraktı. Ardından Volume Shadow Copy'leri sildi ve tüm dosyaları şifreleyerek her klasöre bir fidye notu bıraktı. Şifrelemeden önce 40 GB veri bir bulut depolama servisine yüklenmişti.
>
> **Görev:** Technique zincirini çıkar. Aktör türü ve niyeti nedir? Hangi technique bunu en net ele veriyor?

<details>
<summary><b>✅ ÇÖZÜM — Egzersiz 2</b></summary>

| # | Gözlem | Tactic | Technique |
|---|---|---|---|
| 1 | VPN açığı sömürüldü | Initial Access | **T1190** Exploit Public-Facing Application |
| 2 | LSASS dump | Credential Access | **T1003.001** LSASS Memory |
| 3 | Çalınan kimlikle RDP yayılımı | Lateral Movement | **T1021.001** Remote Desktop Protocol (+ **T1078** Valid Accounts) |
| 4 | Defender devre dışı | Defense Evasion | **T1562.001** Disable or Modify Tools |
| 5 | Shadow Copy silindi | Impact | **T1490** Inhibit System Recovery |
| 6 | Bulut depolamaya 40 GB yükleme | Exfiltration | **T1567.002** Exfiltration to Cloud Storage |
| 7 | Dosyalar şifrelendi + fidye notu | Impact | **T1486** Data Encrypted for Impact |

**Tactic zinciri:** Initial Access → Credential Access → Lateral Movement → Defense Evasion → Exfiltration → Impact.

**Aktör türü:** **Fidye (ransomware) operasyonu — yüksek güven.** En net işaret: **T1486 (Data Encrypted for Impact)**. Ayrıca **T1490 (Inhibit System Recovery)** ve fidye notu fidyeyi doğrular. **Exfiltration + şifreleme birlikte = "double extortion"** (önce çal, sonra şifrele, hem fidye hem ifşa tehdidi) — modern fidye gruplarının imzası. Hızlı lateral movement + commodity teknikler (LSASS dump, RDP) bu profili destekler.

**Ders:** Impact tactic'i (T1486) niyeti **tek başına** ele verir; double extortion (Exfil + Encrypt) çağdaş fidyeyi işaret eder.

*(ID'leri teyit et.)*
</details>

---

### 📦 EGZERSİZ 3 (Zor) — Sessiz casusluk (APT)

> **Senaryo:** Belirli üç üst düzey yöneticiye, kişiselleştirilmiş, kurumsal bir konuyla ilgili bir e-posta geldi; ekteki bir bağlantı, geçerli bir koddan yararlanan özel bir implant indirdi. İmplant, meşru bir Windows hizmeti adıyla (`svchost`-benzeri) kendini gizledi ve WMI olay aboneliğiyle kalıcılık kurdu. Saldırgan haftalarca yalnızca standart Windows araçlarıyla (`net`, `nltest`, `tasklist`) ağı keşfetti, hiç dosya indirmedi. Belirli e-posta kutularından ve belge sunucularından seçili dosyaları toplayıp, normal HTTPS trafiğine karışan şifreli bir kanaldan, küçük parçalar hâlinde dışarı sızdırdı. Hiçbir dosya şifrelenmedi, hiçbir şey bozulmadı. Olay, 4 ay sonra fark edildi.
>
> **Görev:** Technique zincirini çıkar. Aktör türü nedir ve hangi davranışsal işaretler bunu gösteriyor? Bu, Egzersiz 2'den nasıl ayrışıyor?

<details>
<summary><b>✅ ÇÖZÜM — Egzersiz 3</b></summary>

| # | Gözlem | Tactic | Technique |
|---|---|---|---|
| 1 | Üç yöneticiye kişiselleştirilmiş hedefli e-posta + zararlı link | Initial Access | **T1566.002** Spearphishing Link (çok hedefli) |
| 2 | Kullanıcı linke tıkladı / implant indi | Execution | **T1204.001** Malicious Link |
| 3 | Meşru hizmet adıyla gizlenme | Defense Evasion | **T1036.004/.005** Masquerading (Masquerade Task or Service / Match Legitimate Name) |
| 4 | WMI olay aboneliğiyle kalıcılık | Persistence | **T1546.003** WMI Event Subscription |
| 5 | Yalnızca yerleşik araçlarla keşif | Discovery | **T1087** Account Discovery, **T1018** Remote System Discovery, **T1057** Process Discovery (living-off-the-land — **T1059** yerleşik araçlar) |
| 6 | Seçili e-posta/belge toplama | Collection | **T1114** Email Collection, **T1005** Data from Local System |
| 7 | HTTPS'e karışan şifreli C2 | Command and Control | **T1071.001** Web Protocols + **T1573** Encrypted Channel |
| 8 | Küçük parçalar hâlinde sızdırma | Exfiltration | **T1041** Exfil Over C2 Channel (+ muhtemel **T1030** Data Transfer Size Limits) |
| (yok) | Hiçbir şey şifrelenmedi/bozulmadı | Impact | **— (Impact YOK)** |

**Tactic zinciri:** Initial Access → Execution → Defense Evasion → Persistence → Discovery → Collection → C2 → Exfiltration. **Impact yok.**

**Aktör türü:** **Devlet-destekli APT / casusluk operasyonu — yüksek güven.** Davranışsal işaretler:
- **Impact'in YOKLUĞU** + sessiz Collection/Exfiltration → niyet = **bilgi çalmak ve kalmak**, yıkmak değil (casusluğun imzası).
- **Uzun süre + 4 ay tespit edilememe** → derin stealth, sabır.
- **Living-off-the-land** (yalnızca yerleşik araçlar, hiç dosya indirmeme) → tespitten kaçınma önceliği.
- **Çok hedefli spearphishing** (3 belirli yönetici) → istihbarat hedefli, ayrımsız değil.
- **Özel implant + WMI persistence + meşru isim taklidi** → kaynaklı, sofistike aktör.

**Egzersiz 2'den ayrım:** Egzersiz 2 (fidye) = **gürültülü, hızlı, Impact=T1486, para**. Egzersiz 3 (APT) = **sessiz, yavaş, Impact YOK, casusluk**. Aynı erken-aşama benzese de (ikisi de phishing'le başlıyor), **Impact ve stealth/persistence derinliği** iki türü kesin ayırır. **Bu, Bölüm 10'un tür-ayrım tablosunun canlı uygulamasıdır.**

*(ID'leri teyit et — özellikle T1546.003, T1036 alt teknikleri, T1114.)*
</details>

---

### 📦 EGZERSİZ 4 (Uzman) — Belirsiz/karışık senaryo + atfetme tuzağı

> **Senaryo:** Bir kurum saldırıya uğradı. Saldırgan, sızdırılmış meşru bir hesapla VPN'den girdi (initial access için exploit YOK). İçeride **Cobalt Strike** ve **Mimikatz** kullandı, PowerShell ile hareket etti. Bulgular arasında, Rusça dize içeren bir araç ve Moskova saat diliminde (UTC+3) derlenmiş bir binary var. Veri sızdırıldı, sonra sistemlerin bir kısmı şifrelendi (T1486) — ama fidye notu **alışılmadık biçimde** belirli bir ulus-devlet APT grubunun bilinen sloganına benziyordu. Bir vendor blog'u "kesinlikle [ünlü Rus APT grubu]" dedi.
>
> **Görev:** Technique'leri maple. Aktör türü ne? Vendor'ın "kesinlikle [Rus APT]" atfına nasıl yaklaşırsın? Hangi tuzaklar devrede?

<details>
<summary><b>✅ ÇÖZÜM — Egzersiz 4</b></summary>

**Technique eşlemesi (özet):**
| Gözlem | Tactic | Technique |
|---|---|---|
| Sızdırılmış hesapla VPN girişi (exploit yok) | Initial Access | **T1078** Valid Accounts (+ **T1133** External Remote Services) |
| Cobalt Strike (commodity C2/araç) | Execution/C2/çok | **S0154** Cobalt Strike (→ T1059, T1055, T1071…) |
| Mimikatz | Credential Access | **S0002** Mimikatz → **T1003** OS Credential Dumping |
| PowerShell hareketi | Execution | **T1059.001** PowerShell |
| Veri sızdırma | Exfiltration | **T1041 / T1567** |
| Kısmi şifreleme | Impact | **T1486** Data Encrypted for Impact |

**Aktör türü değerlendirmesi:** Yüzeyde **fidye** (T1486) gibi ama **karışık sinyaller** var. Burada asıl ders **atfetme tuzaklarıdır:**

1. **🔁 TTP overlap:** Cobalt Strike (S0154), Mimikatz (S0002), PowerShell — **herkesin kullandığı** commodity araçlar. Bunlar **atıf KANITI değildir**; onlarca aktör (fidye affiliate'leri, APT'ler, kiddie'ler) kullanır. Bu örtüşme **sıfır** ayırt edici değer taşır.
2. **🚩 False flag şüphesi:** "Rusça dize" + "Moskova saat dilimi derlemesi" + "ünlü Rus APT'nin sloganına benzeyen fidye notu" → **fazla temiz, fazla uygun.** Bu artefaktlar **kolayca yerleştirilir** (dil paketi değiştirilir, build zamanı manipüle edilir, slogan kopyalanır). Gerçek bir sofistike aktör imzasını bu kadar bariz bırakmaz → **kasıtlı yanıltma (false flag) ihtimali yüksek.**
3. **🐑 Copy-cat:** Fidye notunun bir APT slogan'ına benzemesi, **bir fidye grubunun o APT'yi taklit ederek** ya korku salmak ya atfı saptırmak istemesi olabilir.
4. **📋 Circular reporting + confirmation bias:** Vendor'ın "kesinlikle [Rus APT]" demesi → tek kaynak, üstelik **mutlak dil** (CTI'da kırmızı bayrak). Diğerleri bunu kopyalarsa döngüsel raporlama doğar.

**Doğru analist tavrı:** *"Gözlemlenen Impact (T1486) bir fidye/yıkıcı amaca işaret ediyor. Ancak 'Rus APT' atfını destekleyen artefaktlar (dil, derleme zamanı, slogan) **kolayca taklit edilebilir niteliktedir ve false flag olasılığını taşır**; kullanılan araçlar (Cobalt Strike, Mimikatz) ise ayrımsız commodity olduğundan atıf değeri taşımaz. **Düşük güvenle**, bu bir fidye operasyonuyla tutarlıdır; **belirli bir devlet APT'sine atıf, mevcut kanıtla desteklenmez ve mutlak ifadeden kaçınılmalıdır.** Atıf için bağımsız altyapı analizi (Diamond — Infrastructure) ve nadir/ayırt edici TTP'ler gerekir."*

**Ders:** Bu egzersiz, Bölüm 10'un **tüm atfetme tehlikelerini** tek senaryoda toplar. "Kesin X" diyen vendor'a değil, **kanıtın niteliğine** bak: araç commodity mi (zayıf)? Artefakt taklit edilebilir mi (false flag)? Tek kaynak ve mutlak dil var mı (kırmızı bayrak)? **Aşırı-atıf, en zararlı CTI hatasıdır.**

*(ID'leri teyit et.)*
</details>

> 🧠 **Antrenmanı sürdürme:** Bu 4 egzersiz şablonu sana yeter — şimdi **gerçek raporlarla** çalış (CISA advisory, vendor APT raporu). Her raporu okurken: önce kendin maple, sonra raporun ATT&CK tablosuyla karşılaştır, **uyuşmadığın yerleri** sorgula (sen mi kaçırdın, rapor mu fazla iddialı?). Atomic Red Team ile bir lab'da technique'leri **çalıştırıp** logunu görmek, "kağıt üstü technique"i "gerçek davranış"a bağlar — usta-seviye köprü budur.

---

<a id="12"></a>
## 12. 🔥🔥 CTI RAPOR YAZIMINDA ATT&CK — Hangi Raporda Ne Kadar

> Bu, kullanıcının açık isteğidir ve sahada en çok yanlış yapılan şeydir: **her rapora aynı miktarda ATT&CK koymak.** Doğru kullanım, **okuyucuya göre ATT&CK dozunu ayarlamaktır.** Bir CISO'ya 40 technique ID'si saçmak da, bir avcıya "kötü amaçlı kod çalıştırdı" demek de başarısızlıktır.

### 🎯 Temel ilke: ATT&CK dozu = okuyucunun teknik seviyesi

```
   STRATEJİK (CISO/yönetici)   ████░░░░░░  ATT&CK MİNİMAL  — iş etkisi dili, ID yok/çok az
   OPERASYONEL (SOC/savunma)   ███████░░░  ATT&CK ORTA-YÜKSEK — kampanya TTP haritası, Navigator
   TAKTİK (analist/avcı)       ██████████  ATT&CK MAKSİMUM — her gözlem ID+procedure+detection
```

### Rapor türleri ve ATT&CK ağırlığı

#### 1. Stratejik Rapor (Strategic) — yönetici/CISO/yönetim kurulu
> **ATT&CK ağırlığı: MİNİMAL.** Bu okuyucu **iş riski** diliyle düşünür: para, itibar, operasyon kesintisi, yasal sorumluluk. Teknik ID **anlamsızdır, hatta zararlıdır** (gözünü korkutur, mesajı boğar).

- **Ne yaz:** "Bu aktör, kurumunuzun verilerini şifreleyerek operasyonu durdurma ve fidye talep etme eğilimindedir; benzer kurumlarda ortalama X gün kesinti ve Y maliyet görülmüştür." → **iş etkisi**, technique değil.
- **ATT&CK kullanımı:** Ya hiç ID yok, ya da en fazla **bir-iki üst-seviye trend** ("aktör giderek daha çok *meşru kimlik bilgisi çalma* yöntemine yöneliyor" — tactic adı, ID yok).
- **Yasak:** "T1486, T1490, T1003.001 gözlemlendi" → CISO için **gürültü.** Asla.

> 💡 Stratejik raporda ATT&CK'i **arka planda** kullan: analizini ATT&CK'le yaptın ama **çıktıyı iş diline çevir.** "Verinizi şifreliyor" cümlesinin arkasında T1486 var, ama okuyucu onu görmez.

#### 2. Operasyonel Rapor (Operational) — SOC, savunma ekibi, IR ekibi, threat hunting lead
> **ATT&CK ağırlığı: ORTA-YÜKSEK.** Bu okuyucu **savunmayı önceliklendirmek** ister: "bu aktöre karşı neyi izleyeyim, neyi sertleştireyim, kapsamamda nerede boşluk var?"

- **Ne yaz:** Kampanyanın **TTP haritası** (technique seviyesi), aktörün **hangi tactic'lerde aktif** olduğu, **savunma önceliklendirme**.
- **ATT&CK kullanımı:** **Navigator layer'ı** (kampanya TTP'leri — Bölüm 9), tactic-bazlı özet, "şu technique'lere karşı tespit önerisi", kör nokta analizi (`tehdit - savunma`).
- **Görsel:** Operasyonel raporun kalbi bir **Navigator heatmap'idir** (SVG) — savunma ekibi bunu kendi kapsamasıyla kesiştirir.

> 💡 Operasyonel rapor, technique ID kullanır ama **eyleme dönük**: "T1003.001 (LSASS dump) gözlendi → **M1043 Credential Access Protection** (LSA korumasını etkinleştir) ve Event ID 4688/Sysmon 10 ile izle." ID + öneri + tespit birlikte.

#### 3. Taktik Rapor (Tactical) — analist, threat hunter, detection engineer, IR analisti
> **ATT&CK ağırlığı: MAKSİMUM.** Bu okuyucu **derinlik** ister: her gözlem, en spesifik technique/sub-technique ile, procedure detayıyla, hangi data source'la yakalandığıyla.

- **Ne yaz:** **Her gözlem** technique ID'siyle, **procedure detayı** ("şu komutu şöyle çalıştırdı"), **detection/data source eşlemesi** ("bu, Sysmon Event 1 + 4104'te görünür"), IOC'ler.
- **ATT&CK kullanımı:** Tam **sub-technique** seviyesi, procedure açıklamaları, her technique için detection logic, data source haritası. Bu, başka bir analistin **tespit kuralı yazabileceği** detaydır.

> 💡 Taktik rapor, ATT&CK'i **eksiksiz** kullanır — ama yine de **bağlamla.** "T1059.001" değil, "T1059.001 (PowerShell): saldırgan `IEX (New-Object Net.WebClient).DownloadString(...)` deseniyle indirici çalıştırdı; PowerShell Script Block Logging (4104) ve Process Creation (Sysmon 1) ile tespit edilebilir."

### 📑 Rapor BÖLÜMLERİNDE ATT&CK kullanımı

Aynı rapor içinde bile, **bölüme göre** ATT&CK dozu değişir:

| Bölüm | ATT&CK kullanımı | Technique ID? |
|---|---|---|
| **Executive Summary** | İş etkisi, üst düzey niyet ("veriyi şifreliyor", "casusluk") | **Hayır** (ya da çok az, tactic adı) |
| **Technical Analysis** | **Tam TTP zinciri**, tactic sırasıyla, sub-technique + procedure | **Evet, maksimum** |
| **Detection & Mitigation** | ATT&CK **Mitigations (M-no)** + **Detections/Data Sources** | **Evet** (M-no + data source) |
| **Appendix / IOC** | **Navigator layer** (JSON+SVG) + **technique tablosu** (ID, ad, tactic, gözlem) | **Evet, tablo halinde** |

> 🧠 **Bölüm mantığı:** Executive Summary'yi **CISO** okur (ID yok); Technical Analysis'i **analist** okur (tam ID); Detection'ı **SOC** kullanır (M-no + data source); Appendix'i **herkes referans** alır (Navigator + tablo). **Tek rapor, çok okuyucu** — her bölüm doğru dozda.

### Technique ID'lerini metinde referanslama (yazım kuralı)

Profesyonel CTI yazımında technique ID'leri **ad + ID** birlikte, parantezle verilir:

> ✅ Doğru: *"Saldırgan, hedefli bir kimlik avı eki (**T1566.001 — Spearphishing Attachment**) ile ilk erişimi sağladı; ardından kodlanmış bir PowerShell betiği (**T1059.001 — PowerShell**) çalıştırdı."*
>
> ❌ Yanlış: *"T1566.001 yapıldı, T1059.001 çalıştı, T1053.005 kuruldu."* (ID yığını, okunmaz, bağlamsız.)

- **İlk geçişte** tam ver: "T1059.001 (PowerShell)". Sonra "PowerShell (T1059.001)" ya da sadece "PowerShell" kullanabilirsin.
- **ID'yi cümlenin içine göm**, listeye değil. Metin **anlatı**, ID **etiket**.
- **Sürüm belirt:** Raporun başında "ATT&CK Enterprise v[x] baz alınmıştır" (technique'ler sürümle değişir).

### 🔥 ÖRNEK RAPOR İSKELETİ (kullanılabilir şablon)

```
╔══════════════════════════════════════════════════════════════╗
║  CTI RAPORU — [Kampanya/Olay Adı]                            ║
║  Sınıflandırma: [TLP:AMBER vb.]  ·  Tarih  ·  ATT&CK v[x]    ║
╠══════════════════════════════════════════════════════════════╣

1. EXECUTIVE SUMMARY                          [ATT&CK: yok/minimal]
   - Ne oldu (iş diliyle): "Bir fidye aktörü, X gün içinde
     verileri şifreledi ve sızdırdı; tahmini etki: ..."
   - Aktör türü + güven: "Orta güvenle çift-gasplı fidye operasyonu"
   - En kritik 3 öneri (iş eylemi)

2. KEY FINDINGS                               [ATT&CK: tactic düzeyi]
   - Aktör hangi amaçlarla aktif (tactic özeti)
   - Niyet: casusluk / fidye / sabotaj? (Impact okuması)

3. TECHNICAL ANALYSIS                         [ATT&CK: MAKSİMUM]
   - Tactic sırasıyla tam TTP zinciri:
     3.1 Initial Access   — T1190 (...) : prosedür detayı
     3.2 Credential Access— T1003.001  : "LSASS'tan ... ile dump"
     3.3 Lateral Movement — T1021.001  : "RDP ile 12 sunucu..."
     3.4 ... (her gözlem: ID + sub-technique + procedure)
     3.5 Impact           — T1486      : "şifreleme + fidye notu"
   - Saldırı zinciri diyagramı (ASCII / görsel)

4. ATTRIBUTION ASSESSMENT                     [ATT&CK + Diamond]
   - TTP imzası karşılaştırması (Navigator kesişim)
   - Aktör türü + olası aile, DERECELI GÜVEN
   - Alternatif hipotezler (ACH) + atfetme sınırlamaları/false flag notu

5. DETECTION & MITIGATION                     [ATT&CK: M-no + data source]
   - Her kilit technique için: Mitigation (M-no) + Detection (data source/log)
   - Örn: T1003.001 → M1043 (LSA Protection) + Event 4688/Sysmon 10

6. APPENDIX                                   [ATT&CK: tablo + Navigator]
   A. IOC listesi (IP/domain/hash) — (GPG imzalı, bkz. GNUPG rehberi)
   B. ATT&CK Navigator layer (JSON + SVG heatmap)
   C. Technique tablosu: | ID | Ad | Tactic | Gözlem/Procedure |
╚══════════════════════════════════════════════════════════════╝
```

> 💡 **İskeletin mantığı:** Yukarıdan aşağı **teknik yoğunluk artar** — Executive (iş dili) → Technical (maksimum ID) → Appendix (yapısal referans). Bir okuyucu kendi seviyesinde durur: CISO bölüm 1'de, SOC 5'te, analist 3 ve 6'da. **ATT&CK her bölümde var ama dozu okuyucuya göre.**

> 🔗 **Rapor + kanıt bütünlüğü:** Appendix'teki IOC listesini ve raporu **GPG ile imzala** (`GNUPG_GPG_USTALIK_REHBERI.md`) → "bu istihbarat bizden ve değişmedi." IOC'leri **MISP**'e yapılandırılmış olarak gönder (ATT&CK galaksisiyle etiketli — `MISP_*` rehberi, Bölüm 19) → makine-okunur, paylaşılabilir istihbarat.

---

<a id="13"></a>
## 13. 🔗 ATT&CK'in Kill Chain · Diamond · Pyramid of Pain ile İlişkisi

ATT&CK tek başına yaşamaz; CTI'nın diğer çerçeveleriyle **katmanlı** çalışır. Burada **kısa** ilişkiyi veriyoruz — derin teori `TTP_AVCILIGI_USTALIK_REHBERI.md`'de.

### Üçlü ilişki — özet

| Çerçeve | Soruya cevap verir | ATT&CK ile ilişki |
|---|---|---|
| **Cyber Kill Chain** | Saldırı **hangi aşamada**? (üst düzey akış) | ATT&CK tactic'leri Kill Chain aşamalarını **derinleştirir** (Bölüm 2) |
| **Diamond Model** | **Kim, neyle, nereden, kime**? (olay yapısı) | ATT&CK, Diamond'ın **Capability** köşesini doldurur (Bölüm 10) |
| **Pyramid of Pain** | Bir IOC'yi engellemek saldırgana **ne kadar acı** verir? | ATT&CK **TTP** seviyesindedir = piramidin **tepesi** (saldırgana en çok acı) |

### Pyramid of Pain ↔ ATT&CK (en önemli ilişki)

```
        ▲  Acı (saldırgana maliyet)
        │
   ████ │  TTPs            ◄── ATT&CK BURADA (davranış — değiştirmesi en pahalı)
   ███  │  Tools
   ██   │  Network/Host Artifacts
   █    │  Domain Names
   ▌    │  IP Addresses
   ▌    │  Hash Values     ◄── değiştirmesi en kolay (saldırgan saniyede değiştirir)
        └──────────────────►
```

> 🧠 **Neden ATT&CK piramidin tepesinde?** Bir saldırgan hash'ini (bir bit değiştir), IP'sini (yeni sunucu), domain'ini (yeni kayıt) **anında** değiştirebilir — bunları engellemek ona **az acı** verir. Ama **davranışını (TTP)** değiştirmek — operasyon tarzını, araç kullanımını, persistence yöntemini değiştirmek — **çok pahalıdır.** ATT&CK, savunmayı bu **acı veren** seviyeye taşır: "bu IP'yi engelle" yerine "bu *davranışı* tespit et." (Bu felsefenin tamamı → `TTP_AVCILIGI_USTALIK_REHBERI.md`.)

### Birleşik kullanım (tek cümle)

> Kill Chain ile saldırının **aşamasını** çerçevele → ATT&CK ile o aşamada **somut TTP'yi** belgele → Diamond ile olayı **dört köşeye** (aktör/yetenek/altyapı/kurban) otur → Pyramid of Pain ile savunmanı **en acı veren** (TTP) seviyeye yükselt. Dördü birlikte, eksiksiz bir CTI analiz iskeletidir.

---

<a id="14"></a>
## 14. 🌍 Gerçek Aktörler & İmza TTP'leri

Profillemeyi somutlaştırmak için gerçek aktörlerin **kamuya açık raporlardan bilinen** karakteristik TTP'lerine bakalım. (Atıf iddiaları kamuya açık kaynaklardandır; **G-numaralarını ve güncel technique listelerini attack.mitre.org'dan teyit et** — aktör profilleri sürümle güncellenir.)

| Aktör (ATT&CK G-no) | Tür | Bilinen imza/karakteristik TTP'ler | Not |
|---|---|---|---|
| **APT28** (G0007) — Fancy Bear, Sofacy, Forest Blizzard | Devlet-destekli (Rusya/GRU iddiası) | Hedefli spearphishing (T1566.001/.002), credential harvesting (T1003, T1110), özel implantlar, casusluk odaklı Collection/Exfiltration | Uzun soluklu casusluk; teyit: G0007 |
| **APT29** (G0016) — Cozy Bear, Midnight Blizzard, The Dukes | Devlet-destekli (Rusya/SVR iddiası) | İleri stealth, living-off-the-land, bulut/kimlik odaklı (token, OAuth), tedarik zinciri (SolarWinds kampanyası ile anılır), sabırlı persistence | Yüksek sofistikasyon; teyit: G0016 |
| **Lazarus Group** (G0032) — Hidden Cobra | Devlet-destekli (K. Kore iddiası) | Hem casusluk hem **finansal** (banka/kripto soygunu — nadir karma), özel malware, yıkıcı bileşenler, geniş araç seti | Casusluk + para birlikte (atipik); teyit: G0032 |
| **APT1** (G0006) — Comment Crew | Devlet-destekli (Çin/PLA iddiası) | Tarihsel; geniş çaplı fikri mülkiyet casusluğu, standartlaşmış araçlar | İlk büyük kamuya açık APT raporlarından (Mandiant 2013) |
| **FIN7** (G0046) | Siber suç (finansal) | POS/ödeme kartı hedefli, sofistike spearphishing, özel araçlar — **suç ama APT-seviyesi titizlik** | "Finansal motivasyonlu ama sofistike" örneği; teyit: G0046 |
| **Fidye grupları** (çeşitli — örn. LockBit, Conti vb.) | Siber suç / fidye (RaaS) | **T1486 (Data Encrypted for Impact)** + T1490 (Inhibit Recovery) + double extortion (Exfil önce), commodity araçlar (Cobalt Strike, Mimikatz), affiliate modeli | **Aktör türü ipucu için en net Impact imzası.** Spesifik grup/G-no için attack.mitre.org'dan **mutlaka** teyit et (bu alan çok hızlı değişir) |

> ⚠️ **Kritik uyarı:** Yukarıdaki **takma adlar ve atıf iddiaları kamuya açık tehdit raporlarından** alıntıdır; kesin atıf devlet-seviyesi bir iddiadır ve bu rehber bunu **eğitim amaçlı** aktarır. **G-numaralarını, güncel takma adları ve technique listelerini her zaman attack.mitre.org'dan doğrula** — özellikle **fidye gruplarının** spesifik G-numaraları ve isimleri sık değişir/eklenir (uydurma G-no verme; emin değilsen "attack.mitre.org Groups'tan teyit et" notu koy).

> 🧠 **Aktör türü tablosuyla bağ (Bölüm 10):** Yukarıdaki gerçek örnekler tür imzalarını doğrular: APT28/29/Lazarus/APT1 → **uzun persistence + stealth + casusluk** (APT imzası); fidye grupları → **T1486 + double extortion** (fidye imzası); FIN7 → **sofistike ama finansal** (tür sınırlarının bulanıklaşabildiğini gösterir — her aktör tabloya tam oturmaz). Profillerken tabloyu **rehber** olarak kullan, **dogma** olarak değil.

---

<a id="15"></a>
## 15. ☠️ Yaygın Ölümcül Hatalar

1. **Technique sayma yarışı** → rapora 40 ID serpiştirip "kapsamlı analiz" sanmak. ATT&CK **ilişkidir**, sayı değil. Bağlamsız ID gürültüdür.
2. **Bağlamsız ID** → "T1059 gözlendi" deyip *ne yaptığını, hangi tactic amacıyla, nasıl tespit edileceğini* yazmamak. ID bir etikettir; anlatı şart.
3. **Over-attribution (aşırı atıf)** → "kesinlikle [ünlü APT]" demek. Atıf **dereceli güven** ister; mutlak dil CTI'da kırmızı bayraktır. (Bölüm 10.)
4. **TTP overlap'i kanıt sanmak** → herkesin kullandığı PowerShell/Mimikatz/Cobalt Strike örtüşmesini "aynı aktör" diye okumak. Yalnızca **nadir/ayırt edici** TTP örtüşmesi atıf değeri taşır.
5. **False flag'e kanmak** → "çok temiz/uygun" artefaktları (dil, saat dilimi, taklit slogan) sorgusuz kabul etmek. Fazla uygunsa, yerleştirilmiş olabilir.
6. **Yanlış granülarite** → emin olmadan sub-technique uydurmak (T1059 yeterken T1059.007 demek) ya da fazla genel kalmak (T1059.001 belliyken sadece "script çalıştı" demek).
7. **Yanlış matris** → bir mobil/ICS olayını Enterprise matrisine zorlamak (ya da tersi). Doğru matris = doğru technique.
8. **Tactic sırasını "playbook" sanmak** → matrisin soldan sağ sırasını zorunlu saldırı akışı zannetmek. ATT&CK doğrusal değildir.
9. **Yanlış seviyede rapor** → CISO'ya sub-technique procedure'u anlatmak ya da avcıya "kötü amaçlı kod çalıştırıldı" demek. Doz = okuyucu seviyesi. (Bölüm 12.)
10. **Eski/yanlış ID** → deprecated bir technique'i ya da hatalı numarayı kullanmak. ATT&CK canlıdır; **her ID'yi attack.mitre.org'dan teyit et** ve **sürüm belirt.**
11. **Impact'i atlamak** → niyeti (fidye mi casusluk mu sabotaj mı) belirlerken Impact tactic'ine bakmamak. Niyet en net Impact'te okunur.
12. **Tek kaynak/circular reporting** → bir vendor'ın atfını birincil kaynağa inmeden kopyalamak; "herkes diyor" tuzağı.
13. **Confirmation bias** → istediği aktörü destekleyen kanıtı görüp çürüteni görmezden gelmek. Hipotezini **çürütmeye** çalış (ACH).
14. **Sadece IOC'ye takılmak** → ATT&CK varken hash/IP seviyesinde kalmak (Pyramid of Pain'in dibinde). Saldırgana acı veren **TTP** seviyesine çık.
15. **Veriyi tüm resim sanmak** → gördüğün TTP'leri "aktörün tüm davranışı" zannetmek. Gördüklerin sadece **yakalayabildiklerin**; eksik telemetri eksik resim demektir.

---

<a id="16"></a>
## 16. 🏰 Kanije Kalesi ile Birlikte Kullanım

Bu repo (Kanije Kalesi), CTI/güvenlik odaklı bir **muhafız ve araç setidir.** ATT&CK, onun olaylarını **standart bir dile** çevirir ve — ilginç biçimde — Kanije'nin **kendi kodu** ATT&CK farkındalığıyla yazılmıştır. İki yön var: (a) **Kanije'nin kendi davranışının ATT&CK eşlemesi**, (b) **Kanije'nin yakaladığı olayları ATT&CK'e mapleme.**

### 16.1 Kanije'nin kendi kodunun MITRE eşlemesi

Kanije bir **savunma** aracıdır ama fiziksel-tehdit anında bazı "agresif" aksiyonlar alır (kilitleme, secure-wipe, honeypot). Bu aksiyonların **ATT&CK'te karşılıkları** vardır ve repo, bunları **bilinçli olarak** ele alır — özellikle **AV/EDR'yi tetiklememek** için:

| Kanije davranışı | ATT&CK technique | Kanije'nin yaklaşımı (repo notu) |
|---|---|---|
| `/imha` — hassas dosyaları (bot token, config, DB, log) overwrite-then-delete ile güvenli silme | **T1485** Data Destruction / **T1070.004** File Deletion | `destroy.go`/`destroy_windows.go`: hedefli secure-wipe, **factory reset YOK** → kasten **T1561 (Disk Wipe)** AV sinyalini *vermez* ("keeps the OS intact and avoids the T1561 AV signal"). Yıkıcı değil, **veri koruma** amaçlı. |
| **T1561** Disk Wipe (kaçınılan) | **T1561** Disk Wipe | Kanije bunu **bilinçli kullanmaz** — tüm disk silme/format AV/EDR tarafından yıkıcı malware sayılır; Kanije yalnızca **kendi hassas izini** siler |
| `/koruma` — dead-man switch, USB çıkarma, yanlış-giriş tetikleyici → kilitle + alarm + foto | (savunma — Defense/Deception tarafı) | `deadman.go`, `protect.go`: tehlikede otomatik aksiyon. Saldırgan-TTP değil, **savunma tetikleyicisi** |
| `/tuzak` (honeypot) — masaüstünde decoy dosyalar, SACL ile erişim izleme | **Deception** (savunma) / kavramsal **T1056-karşıtı** | `canary.go`, `canary_commands.go`: saldırganın **verisini KOPYALAMAZ**, yalnızca tuzağa dokunanı yakalar. Bu, ATT&CK'in **savunma** (deception) yüzü |
| `/kilit tam` (lockdown), `/panik` | (acil savunma) | Olay anı yönetimi — saldırgan davranışı değil, **savunan tarafın** refleksi |

> 🔥 **Kritik felsefe — "T1561'den kaçınmak":** Kanije'nin kodundaki en öğretici ATT&CK kararı şudur: bir veri-koruma aracı *teknik olarak* tüm diski silebilir (T1561 Disk Wipe), ama bu davranış **her AV/EDR tarafından yıkıcı malware imzası olarak** yakalanır (çünkü gerçek wiper'lar — NotPetya sınıfı — bunu yapar). Kanije bu yüzden **factory reset/disk wipe yapmaz**; yalnızca **kendi hassas dosyalarını** (token, DB, log) overwrite-then-delete ile siler. Sonuç: veri-koruma amacına ulaşır ama **savunma yazılımının önünde "iyi huylu"** kalır. **Bu, ATT&CK'i savunma tasarımında kullanmanın güzel bir örneğidir:** "yapabildiğim ama AV-sinyali yüzünden yapMAdığım technique."

> 🧠 **Masquerading (T1036) & Indicator Removal (T1070) farkındalığı:** Bir muhafızın bu technique'leri **bilmesi** kritiktir — çünkü (a) Kanije *kendi* meşruiyetini korumak için bunlardan **kaçınır** (gizlenmeye çalışan bir "güvenlik aracı" şüphelidir), (b) Kanije'nin yakaladığı bir saldırgan **bunları kullanırsa** (örn. log silme = T1070), bunu tanıyıp doğru technique'e maplemek gerekir. Savunma aracını yazarken "hangi davranışım hangi ATT&CK technique'ine benzer ve bu beni şüpheli yapar mı?" sorusu, **tasarım** sorusudur.

### 16.2 Kanije olaylarını ATT&CK'e mapleme

Kanije bir olay yakaladığında (USB takıldı, yanlış giriş, honeypot dokunuldu, ağ anomalisi), bunu **ATT&CK diline** çevirmek, olayı daha büyük bir CTI bağlamına oturtur:

| Kanije olayı | Olası ATT&CK eşlemesi (saldırgan tarafından) | CTI değeri |
|---|---|---|
| `/tuzak` honeypot dosyasına erişim | **T1083** File and Directory Discovery / **T1005** Data from Local System (saldırgan keşfi) | Birinin/bir şeyin **keşif/collection** yaptığının kanıtı → olay tetikle |
| USB cihaz takılması (yetkisiz) | **T1091** Replication Through Removable Media / **T1200** Hardware Additions | Fiziksel erişim/yayılma vektörü işareti |
| Tekrarlı yanlış giriş | **T1110** Brute Force | Kimlik bilgisi saldırısı göstergesi |
| Ağ anomalisi / beklenmedik bağlantı | **T1071** Application Layer Protocol (olası C2) | Şüpheli C2/exfil işareti → `WIRESHARK_*` ile derinleş |

> 💡 **Önerilen entegrasyon deseni:**
> 1. **Kanije olayını technique'e maple:** Bir `/defender` tespiti ya da `/erisim`/`/tuzak` olayı geldiğinde, onu ATT&CK technique'ine çevir → olayı **standart dile** sok, başka tespitlerle karşılaştırılabilir yap.
> 2. **IOC zenginleştir + aktör bağlamı:** Olaydan çıkan IOC'leri (IP, dosya) OSINT'le zenginleştir (`OSINT_ARAC_SETI_USTALIK_REHBERI.md`), ATT&CK Groups'ta o technique'i kullanan aktörlere bak → "bu davranış bilinen bir aktör türüyle tutarlı mı?"
> 3. **Navigator layer'ı tut:** Kanije'nin yakaladığı olaylardan zamanla bir **gözlemlenen-TTP layer'ı** oluştur (Bölüm 9) → "bu cihaza/kuruma karşı en çok hangi technique'ler deneniyor?" heatmap'i.
> 4. **Kanıt zinciri:** Olayları + ATT&CK eşlemesini + IOC'leri **GPG imzala** (`GNUPG_GPG_USTALIK_REHBERI.md`) ve gerekirse **MISP**'e gönder (ATT&CK galaksisiyle etiketli) → savunulabilir, paylaşılabilir istihbarat.
> 5. **Savunma boşluğu:** Kanije'nin koruduğu uç noktanın savunma kapsamasını ATT&CK'e mapleyip (`tehdit - savunma` heatmap'i) kör noktaları kapat — Kanije bir katman, ama tek katman değil.

> 🧠 **Felsefe örtüşmesi:** ATT&CK **saldırgan davranışını standartlaştırır**; Kanije Kalesi **olay anını yakalar ve savunur.** Birlikte: Kanije'nin yakaladığı her olayı ATT&CK diline çevir → olay artık *karşılaştırılabilir, aktör-bağlamına oturtulabilir, raporlanabilir* bir TTP gözlemidir. Kanije gözünü açar; ATT&CK gördüğüne **ortak bir ad** verir.

---

<a id="17"></a>
## 17. 🗺️ Sıfırdan Uzmana Yol Haritası (Hafta Hafta)

ATT&CK ustalığı bir gecede gelmez. İşte ölçülü, kaynaklı bir program:

### Hafta 1 — Temel kavramlar & matris okuma
- **Oku:** Bu rehberin Bölüm 1-5'i + `attack.mitre.org` ana sayfası, "Getting Started".
- **Yap:** Enterprise matrisini aç, **14 tactic'i** ezberle (bu rehberin Bölüm 3 + Bölüm 18 tablosu). Her tactic'in "amacını" kendi cümlenle yaz.
- **Hedef:** "Tactic = neden, technique = nasıl" refleksi otursun.

### Hafta 2 — Technique sayfalarını okuma
- **Oku:** Bölüm 6. `attack.mitre.org`'da **5 technique sayfasını** baştan sona gez (T1059, T1566, T1003, T1486, T1071 öner).
- **Yap:** Her birinin sub-technique'lerini, Procedure Examples'ını, Mitigations ve Detection'ını çıkar. Bir technique için "nasıl tespit ederim?" sorusunu cevapla.
- **Hedef:** Herhangi bir technique sayfasını akıcı okuyabilmek.

### Hafta 3 — Groups & Software + ilişkilerde gezinme
- **Oku:** Bölüm 7-8. Bir **grup sayfasını** (G0007 APT28) ve bir **software sayfasını** (S0002 Mimikatz) derinlemesine incele.
- **Yap:** Technique → Procedure → Group → diğer technique'ler pivot zincirini elle yap (Bölüm 8 diyagramı). İsim↔G-numarası eşlemesini (Fancy Bear=G0007) pratik et.
- **Hedef:** ATT&CK'i "arama" değil "gezinme" aracı olarak kullanmak.

### Hafta 4 — Navigator
- **Oku:** Bölüm 9. `mitre-attack.github.io/attack-navigator/` aç.
- **Yap:** İlk layer'ını oluştur; bir grubun hazır layer'ını import et; **iki aktörü kesiştir** (`a*b`); bir **savunma kapsama heatmap'i** (`tehdit - savunma`) üret; JSON export/import et.
- **Hedef:** Navigator'da aktör karşılaştırması ve kapsama haritası üretebilmek.

### Hafta 5-6 — Mapleme antrenmanı (en kritik)
- **Oku:** Bölüm 10-11.
- **Yap:** Bu rehberdeki **4 egzersizi** çöz. Sonra **gerçek CTI raporları** al (CISA advisory'leri, vendor APT blog'ları), her cümleyi technique'e maple, raporun kendi eşlemesiyle karşılaştır. **Atomic Red Team**'i bir izole lab'da kur, birkaç technique çalıştır, logunu gör.
- **Hedef:** Bir saldırı anlatısını/log'u akıcı şekilde technique zincirine çevirmek + aktör türü çıkarmak.

### Hafta 7 — Aktör profilleme & atfetme disiplini
- **Oku:** Bölüm 10 (derinlemesine) + Bölüm 14.
- **Yap:** Bir kampanya seç, TTP setini çıkar, Groups'la karşılaştır, **aktör türünü** belirle, **dereceli güvenle** bir hipotez yaz. Atfetme tuzaklarını (false flag, overlap) kendi analizinde **bilinçli ara.**
- **Hedef:** Aşırı-atıftan kaçınan, ACH disiplinli, savunulabilir profilleme.

### Hafta 8 — Rapor yazımı + ekosistem
- **Oku:** Bölüm 12-13 + `TTP_AVCILIGI_USTALIK_REHBERI.md` (TTP teorisi) + MISP rehberi (Bölüm 19).
- **Yap:** Bir olay için **üç versiyon** rapor taslağı (stratejik/operasyonel/taktik) yaz — her birinde doğru ATT&CK dozuyla. Navigator layer'ını ek yap. IOC'leri MISP formatında düşün.
- **Hedef:** Okuyucuya göre doz ayarlayan, ATT&CK'i diğer çerçevelerle (Kill Chain/Diamond/Pyramid) bütünleştiren profesyonel CTI çıktısı.

### İleri seviye (devam eden)
- **MITRE ATT&CK Defender (MAD)** sertifikasyonu (resmi eğitim).
- **MITRE ATT&CK for CTI** eğitim modülleri (ücretsiz, attack.mitre.org → Training).
- ATT&CK STIX verisini programatik kullanma (`mitreattack-python`).
- Düzenli olarak **yeni ATT&CK sürümlerini** takip et (release notes — neler değişti, hangi technique bölündü/eklendi).

> 🧠 **Yol haritası felsefesi:** İlk 4 hafta **araç/framework** (matris, sayfa, gezinme, Navigator). 5-7 **beceri** (mapleme, profilleme — asıl ustalık burada). Hafta 8 **çıktı** (rapor). Çoğu kişi 1-4'te takılır ("ID ezberler") ve 5-7'yi (gerçek beceri) atlar. **ATT&CK ustalığı, sayfaları ezberlemek değil, gözlemi davranışa ve davranışı aktör hipotezine çevirebilmektir.**

---

<a id="18"></a>
## 18. ✅ Hızlı Referans (14 Tactic Tablosu) & Kontrol Listesi

### 🔥 14 Enterprise Tactic — hızlı referans tablosu

| # | Tactic | ID | Saldırganın amacı (NEDEN) | Örnek technique |
|---|---|---|---|---|
| 1 | **Reconnaissance** | TA0043 | Hedef hakkında bilgi topla (saldırı öncesi) | T1595 Active Scanning, T1589 Gather Identity Info |
| 2 | **Resource Development** | TA0042 | Altyapı/araç/kimlik edin, hazırla | T1583 Acquire Infrastructure, T1587 Develop Capabilities |
| 3 | **Initial Access** | TA0001 | Ağa ilk ayağı sok | T1566 Phishing, T1190 Exploit Public-Facing App, T1078 Valid Accounts |
| 4 | **Execution** | TA0002 | Kötü amaçlı kodu çalıştır | T1059 Command/Scripting Interpreter, T1204 User Execution |
| 5 | **Persistence** | TA0003 | Erişimi koru (içeride kal) | T1547 Autostart, T1053 Scheduled Task, T1136 Create Account |
| 6 | **Privilege Escalation** | TA0004 | Daha yüksek izin al | T1548 Abuse Elevation, T1068 Exploitation for PrivEsc |
| 7 | **Defense Evasion** | TA0005 | Tespitten kaç, izi gizle (en kalabalık) | T1070 Indicator Removal, T1036 Masquerading, T1027 Obfuscation |
| 8 | **Credential Access** | TA0006 | Kimlik bilgisi çal | T1003 OS Credential Dumping, T1110 Brute Force, T1056 Input Capture |
| 9 | **Discovery** | TA0007 | Ortamı içeriden tanı | T1057 Process Disc., T1082 System Info, T1018 Remote System Disc. |
| 10 | **Lateral Movement** | TA0008 | Ağda yayıl | T1021 Remote Services, T1550 Alt. Auth Material (pass-the-hash) |
| 11 | **Collection** | TA0009 | Sızdırmadan önce veri topla | T1113 Screen Capture, T1005 Data from Local Sys, T1560 Archive |
| 12 | **Command and Control** | TA0011 | Uzaktan kumanda kanalı kur | T1071 App Layer Protocol, T1573 Encrypted Channel |
| 13 | **Exfiltration** | TA0010 | Veriyi dışarı çıkar | T1041 Exfil Over C2, T1567 Over Web Service, T1048 Alt. Protocol |
| 14 | **Impact** | TA0040 | Boz/yok et/erişilemez kıl/manipüle et | **T1486 Data Encrypted** (fidye), T1485 Destruction, T1490 Inhibit Recovery |

> ⚠️ Tüm ID'ler **attack.mitre.org'dan teyit edilmelidir** (sürümle değişebilir). Tactic ID'leri (TA00xx) stabildir; technique ID'leri (Txxxx) bölünebilir/güncellenebilir.

### 🔍 Aktör türü hızlı ipucu (Bölüm 10 özeti)

| Gördüğün | Muhtemel aktör türü |
|---|---|
| **T1486** (Data Encrypted) + T1490 + double extortion | **Fidye / siber suç** |
| Uzun persistence + derin stealth (T1070/T1027/T1036) + Impact YOK + casusluk Collection | **Devlet-destekli APT** |
| **T1491** (Defacement) + T1498 (DDoS) + kamuya açıklama | **Hacktivist** |
| Sadece T1078 (Valid Accounts), initial access yok, içeriden anormal | **Insider** |
| **T1485/T1561** (Destruction/Wipe), yıkıcı, geri dönüşsüz | **Sabotaj / yıkıcı (bazı devlet op.)** |
| Dağınık, otomatik, ayrımsız, hazır kit | **Script kiddie / fırsatçı** |

### ✅ Technique mapleme kontrol listesi (her olayda)
- [ ] Olayı **çıplak davranışlara** ayırdım (her gözlem ayrı)
- [ ] Her gözlem için **tactic** (amaç) belirledim — NEDEN?
- [ ] Her gözleme **en spesifik doğru technique/sub-technique** verdim (uydurma yok)
- [ ] Emin olmadığım sub-technique'te **bir üst seviyede** kaldım
- [ ] **Doğru matris** (Enterprise/Mobile/ICS) kullandım
- [ ] Her ID'yi **attack.mitre.org'dan teyit ettim** + **sürüm** not ettim
- [ ] **Impact** tactic'ine bakıp **niyeti** (fidye/casusluk/sabotaj) okudum

### ✅ Aktör profilleme & atfetme kontrol listesi
- [ ] Gözlemlerden **TTP setini** çıkardım
- [ ] Groups'u TTP setine göre karşılaştırdım (Navigator kesişim)
- [ ] **Aktör TÜRÜNÜ** TTP imzasından çıkardım (tür tablosu)
- [ ] Atıfı **dereceli güvenle** ifade ettim ("kesin" demedim)
- [ ] **False flag / TTP overlap / copy-cat / circular reporting** tuzaklarını sorguladım
- [ ] **Alternatif hipotez** (ACH) düşündüm; hipotezimi **çürütmeye** çalıştım
- [ ] Yaygın TTP örtüşmesini **kanıt saymadım** (yalnızca nadir TTP)
- [ ] **Diamond** dört köşesini (özellikle Infrastructure'ı OSINT'le) dahil ettim

### ✅ Rapor yazımı kontrol listesi (Bölüm 12)
- [ ] Okuyucuyu belirledim → **ATT&CK dozunu** ayarladım (stratejik/operasyonel/taktik)
- [ ] **Executive Summary'de ID YOK** (iş etkisi dili)
- [ ] **Technical Analysis'te tam TTP zinciri** (sub-technique + procedure)
- [ ] **Detection/Mitigation'da M-no + data source** verdim
- [ ] **Appendix'te Navigator layer (JSON+SVG) + technique tablosu** ekledim
- [ ] ID'leri metne **gömdüm** (ad + ID), liste yığını yapmadım
- [ ] **ATT&CK sürümünü** belirttim
- [ ] IOC/raporu **GPG imzaladım** (kanıt bütünlüğü)

---

<a id="19"></a>
## 19. 📚 Kaynaklar & Çapraz Referanslar

### Resmi MITRE ATT&CK kaynakları
- **attack.mitre.org** — ana bilgi tabanı (technique, group, software, mitigation, data source, campaign).
- **ATT&CK Navigator** — `mitre-attack.github.io/attack-navigator/` (layer, heatmap, karşılaştırma).
- **ATT&CK Training** — attack.mitre.org → Resources → Training ("ATT&CK for CTI" modülleri, ücretsiz).
- **MITRE ATT&CK Defender (MAD)** — resmi sertifika programı (bazı modüller ücretsiz).
- **ATT&CK STIX data** — `github.com/mitre-attack/attack-stix-data` + TAXII (programatik erişim).
- **ATT&CK release notes** — sürüm değişiklikleri (hangi technique eklendi/bölündü).
- **Atomic Red Team** — `github.com/redcanaryco/atomic-red-team` (technique başına çalıştırılabilir test — **yalnızca izole lab**).

### 🔗 Bu koleksiyondaki ilgili rehberler (çapraz referans)
- **`TTP_AVCILIGI_USTALIK_REHBERI.md`** — **TTP kavramının derin teorisi**: Indicator vs TTP, Pyramid of Pain tam mantığı, Diamond Model derinliği, davranış-temelli avcılık felsefesi. *(Bu rehber ATT&CK aracını kullanır; o rehber TTP'nin "neden"ini işler — ikisi ikizdir.)*
- **`MISP_*_REHBERI.md`** — IOC ve TTP istihbaratını **yapılandırılmış paylaşma/yaşatma**: MISP'in ATT&CK galaksisi, STIX export, olay paylaşımı. *(ATT&CK ile mapelediğin TTP'leri MISP'te kataloglar, paylaşırsın.)*
- **`OSINT_ARAC_SETI_USTALIK_REHBERI.md`** — Reconnaissance/Resource Development tarafı + **altyapı pivotu** (Shodan `ssl:`/JARM, crt.sh, pDNS) = Diamond'ın Infrastructure köşesi. Aktör altyapısını pasif haritalama.
- **`WIRESHARK_AG_ANALIZ_USTALIK_REHBERI.md`** — C2 (T1071), Exfiltration (T1041/T1048) trafiğini ağ seviyesinde yakalama; technique'leri paket kanıtına bağlama.
- **`MALWARE_ANALIZ_USTALIK_REHBERI.md`** — Software (S-numaraları) analizinin derinliği; bir malware'in hangi technique'leri uyguladığını tersine mühendislikle çıkarma.
- **`GNUPG_GPG_USTALIK_REHBERI.md`** — CTI raporu ve IOC listesini **imzalama** (kanıt bütünlüğü/zinciri).
- **`KEEPASSXC_PAROLA_KALESI_REHBERI.md`** — Credential Access (T1003/T1110/T1555) tarafına **savunma**; analist anahtar/kimlik hijyeni.
- **`VERACRYPT_USTALIK_REHBERI.md`** — Hassas CTI verisini (analiz, IOC, rapor taslakları) şifreli saklama.

---

> 🏰 **Kapanış:** MITRE ATT&CK bir technique listesi değil, saldırgan davranışının **ortak dilidir.** En uzun T-numarası tablosunu ezberlemek bile, o ID'leri **gözleme bağlayamıyorsan, aktör türüne çeviremiyorsan ve doğru okuyucuya doğru dozda anlatamıyorsan** çaresizdir — hatta (aşırı-atıfla) tehlikelidir. ATT&CK sana, dünyanın her analistiyle **aynı dili konuşma** gücünü verir; **o dili gözlemi davranışa, davranışı aktör hipotezine ve hipotezi savunulabilir bir rapora çevirmek için kullanmak senin işin.** ATT&CK'in kalbi matriste değil, **mapleme, profilleme ve dürüst (dereceli) atıfta**dır — matris yalnızca sözlüktür. Kanije Kalesi de tam burada: yakaladığı her olayı bu ortak dile çevirip, savunmanı saldırganın **en çok acı duyduğu** (TTP) seviyeye taşıyan nöbetçi olarak devreye girer.
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. İlgili: `TTP_AVCILIGI_USTALIK_REHBERI.md` (TTP teorisi — bu rehberin ikizi), `MISP_*_REHBERI.md` (TTP/IOC paylaşımı), `OSINT_ARAC_SETI_USTALIK_REHBERI.md` (altyapı/Recon), `WIRESHARK_AG_ANALIZ_USTALIK_REHBERI.md` (C2/Exfil trafiği), `MALWARE_ANALIZ_USTALIK_REHBERI.md` (Software/S-no analizi), `GNUPG_GPG_USTALIK_REHBERI.md` (kanıt imzalama).*
