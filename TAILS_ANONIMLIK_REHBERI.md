# 🧅 TAILS OS — TAM ANONİMLİK & AMNEZİ REHBERİ
## Live USB'den Tor-Zorunlu Anonimliğe, Püf Noktalarıyla Uçtan Uca

> **Amaç:** Tails'i "USB'ye yaz, aç, gez" seviyesinden çıkarıp, bir gözetim aktörü karşısında **gerçekten anonim** kalacak şekilde kullanmayı öğretmek. Bu rehber yalnızca *nasıl*'ı değil, **neden** ve **hangi durumda anonimliğin çöker**'i de anlatır. Forum cevaplarında bulamayacağın korelasyon saldırıları, exit node tehlikesi, Persistent Storage'ın inkâr edilebilirliği nasıl bozduğu, Evil Maid ve amnezi felsefesinin sınırları burada.

> ⚠️ **Önce bunu oku:** Anonimlik aracı, *yanlış kullanıldığında* sana **yanlış güvenlik hissi** verir — bu, hiç anonimlik aracı kullanmamaktan daha tehlikelidir, çünkü riskli işler yaparsın. Bir bölümü anlamadan uygulama. Özellikle **Tehdit Modeli**, **Persistent Storage** ve **Korelasyon Saldırıları** bölümlerini atlama.

---

## 📑 İÇİNDEKİLER

1. [Tails Nedir, Neden? (Whonix/Qubes Kıyası)](#1)
2. [Tehdit Modeli — Neyi Korur, Neyi KORUMAZ](#2)
3. [Kurulum + İmza Doğrulama (Atlama!) & USB'ye Yazma](#3)
4. [İlk Açılış — Welcome Screen & Ağ Bağlantısı](#4)
5. [Tor Entegrasyonu — Devre Yönetimi & .onion](#5)
6. [MAC Adresi Spoofing — Varsayılan & Ne Zaman Kapatılmalı](#6)
7. [Sansür Aşma — Bridges, obfs4, meek](#7)
8. [Persistent Storage — Kalıcı Şifreli Depolama & RİSKLERİ](#8)
9. [Tails İçi Araçlar — KeePassXC / GnuPG / MAT2 / Electrum](#9)
10. [Güvenli Kapatma & RAM Temizleme (Cold-Boot Karşıtı)](#10)
11. [🔥 PÜF NOKTALARI — Piyasada Bulamayacakların](#11)
12. [Yaygın Ölümcül Hatalar](#12)
13. [🏰 Kanije Kalesi ile Birlikte Kullanım](#13)
14. [Hızlı Referans & Operasyonel Kontrol Listesi](#14)
15. [Hukuki Sınır & Sınır Geçişi Notu](#15)

---

<a id="1"></a>
## 1. 🧭 Tails Nedir, Neden?

**Tails** = **T**he **A**mnesic **I**ncognito **L**ive **S**ystem. Debian tabanlı, **USB bellekten çalışan canlı (live)** bir işletim sistemidir. İki temel söz verir:

1. **Amnezi (Amnesic):** Her kapanışta **her şeyi unutur.** Çalışırken tüm veriler yalnızca **RAM'de** tutulur; kapatınca RAM silinir → kullandığın bilgisayarda **hiçbir iz kalmaz.** Diske dokunmaz (sen açıkça istemedikçe).
2. **Anonimlik (Incognito):** **Tüm** internet trafiği **zorla Tor üzerinden** yönlendirilir. Tor dışına çıkmaya çalışan bağlantı (uygulama Tor bilmese bile) **engellenir** — sızıntı olmaz. IP adresin gittiğin siteye değil, Tor exit node'una görünür.

Açık kaynaktır, üreticisi kâr amacı gütmeyen **Tails Project**'tir (eskiden bağımsızdı, 2024'te **Tor Project** ile birleşti). Snowden'ın, gazetecilerin ve insan hakları savunucularının kullandığı referans araçtır.

### Neden "live USB + RAM"? — Felsefe
Normal bir bilgisayar sürekli diske yazar: geçici dosyalar, log, tarayıcı geçmişi, swap, thumbnail. Bu izler **adli incelemede** seni ele verir. Tails bu modeli ters çevirir: **disk yok, sadece RAM.** Güç kesilince kanıt da kesilir.

### Whonix / Qubes ile Dürüst Kıyas

| Özellik | Tails | Whonix | Qubes OS |
|---|---|---|---|
| Çalışma modeli | Live USB, **amnezik** | Kalıcı (VM içinde) | Kalıcı (bölümlenmiş OS) |
| **Yerel iz bırakmama** | ✅ **En güçlü o** | ❌ (host'a yazar) | Kısmî (disposable VM) |
| Tüm trafik Tor zorunlu | ✅ | ✅ (Gateway VM) | Hayır (Whonix entegre edilirse) |
| İzolasyon (compartment) | Düşük (tek sistem) | Orta | ✅ **En güçlü o** |
| Donanım izi (USB taşınabilir) | ✅ Cebe sığar | ❌ (kurulu makine) | ❌ (güçlü PC gerekir) |
| Kurulum kolaylığı | ✅ Çok kolay | Orta | Zor |
| Malware'in kalıcılığı | **Yok** (reboot temizler) | Olabilir | VM'e hapsolur |

**CTI / gazetecilik perspektifi:** **Tails**, "tek seferlik, iz bırakmayan, anonim oturum" için rakipsizdir — kaynakla görüşme, hassas dosya alma, anonim yayın. **Whonix/Qubes**, kalıcı bir anonim kimlik ve güçlü uygulama izolasyonu gerektiğinde (örn. uzun vadeli takma ad, malware analizi) daha uygundur. Hassas, **kısa ömürlü, yüksek-riskli** operasyonda Tails'in amnezisi en güçlü duruştur; ama bir saldırgan exit node'unu izliyorsa Tails de seni korelasyondan korumaz (Bölüm 2 & 11).

> 💡 **Katmanlı yaklaşım:** Tails'i tek başına kullanabilirsin, ama hassas dosyaları **VeraCrypt hidden volume** içinde tutup Tails'te açmak (Bölüm 13) hem amnezi hem inkâr edilebilirlik katar. Tails RAM-only, VeraCrypt durağan-veri — ikisi örtüşür.

---

<a id="2"></a>
## 2. 🎯 Tehdit Modeli — Neyi Korur, Neyi KORUMAZ

Tails'i kullanmaya başlamadan **kime karşı** korunduğunu netleştir. Tails **anonimlik** aracıdır, **uçtan uca dokunulmazlık** değil. Yanlış varsayım hayatına mal olabilir.

### ✅ Tails'in KORUDUĞU senaryolar
- **Yerel iz bırakmama:** Kullandığın bilgisayarın diskine dokunmaz. Kapanınca o makinede "ne yaptın" sorusunun adli cevabı **yoktur** (RAM silinir).
- **ISP / yerel ağdan gizlenme:** ISP'n yalnızca "Tor'a (ya da bir bridge'e) bağlandığını" görür — **hangi siteye gittiğini değil.**
- **IP/konum anonimliği (gittiğin siteye karşı):** Ziyaret ettiğin site senin gerçek IP'ni değil, Tor exit node'unun IP'sini görür.
- **Anonim kimlik:** Gerçek kimliğinle bağlantılı hiçbir veri (çerez, hesap, donanım izi) varsayılan olarak kalıcı değildir.
- **Bulaşan malware'in kalıcılığı:** Bir web zafiyeti seni o oturumda vursa bile, **reboot** ile temizlenir (diske yerleşemez — okuma-yalnız sistem).
- **El konulan bilgisayar (Tails kapalı, USB ayrı):** Makinede iz yok; USB başka yerde.

### ❌ Tails'in KORUMADIĞI senaryolar
- **Tor exit node gözetimi + şifresiz trafik:** Exit node, **HTTP (şifresiz)** trafiğini **okuyabilir/değiştirebilir.** Tor *kim olduğunu* gizler, *ne gönderdiğini* (HTTPS yoksa) gizlemez. **HTTPS şarttır.** (Bölüm 11.1)
- **Korelasyon / zamanlama saldırıları:** Hem girişini (sen → Tor) hem çıkışını (exit → hedef) izleyebilen **global pasif rakip**, trafik zamanlamasını eşleştirerek seni de-anonimize edebilir. Tor bu saldırıya karşı **tasarımı gereği** tam koruma vermez. (Bölüm 11.2)
- **Son nokta uzlaşması (endpoint compromise):** Hedef sunucu ele geçirilmişse ya da sen kendi kimliğini sızdırırsan (hesaba giriş, isim yazma) anonimlik biter. **Tor, senin yaptığın hatayı düzeltemez.**
- **Kötücül USB / firmware (Evil Maid):** Tails'i yazdığın USB ya da çalıştırdığın bilgisayarın **BIOS/UEFI/firmware**'i kötücül koda sahipse, Tails başlamadan önce çalışan kod seni izleyebilir. (Bölüm 11.4)
- **Global pasif rakip / devlet-seviyesi ağ gözlemi:** İnternet omurgasının büyük kısmını izleyen aktör, uçtan uca korelasyonla Tor kullanıcılarını hedefleyebilir.
- **Sen kendini ele verirsen:** Gerçek adınla hesap açmak, kişisel bilgi yazmak, aynı takma adı farklı bağlamlarda kullanmak → anonimlik **senin elinle** çöker.
- **"Tails kullandığın" gerçeğinin gizlenmesi:** ISP'n Tor'a bağlandığını görür (bridge kullanmazsan). Bazı bağlamlarda **Tor kullanıyor olman bile** dikkat çeker. (Bölüm 7 & 11.10)

> 🧠 **Altın kural:** Tails **"kim olduğunu" (kimlik/konum)** gizler. **"Ne söylediğini"** (içerik) HTTPS gizler. **"Tor kullandığını"** bridge gizler. Üçü ayrı problemdir; üçünü de ayrı yönetmen gerekir. Tek başına Tails hepsini çözmez.

---

<a id="3"></a>
## 3. 📥 Kurulum + İmza Doğrulama & USB'ye Yazma

**Asla** rastgele bir aynadan ya da torrent'ten indirme. Yalnızca **resmi site: `tails.net`** (eski adres `tails.boum.org`).

### İmza / bütünlük doğrulama (ATLAMA — kurulumun en kritik adımı)
İndirdiğin Tails imajı kötü amaçlıyla değiştirildiyse, **anonimlik aracın bizzat seni gözetleyen araç olur.** Bu yüzden imajı *kurmadan önce* doğrula. Tails iki yol sunar:

1. **Tarayıcı eklentisi / web doğrulama (en kolay):** `tails.net` indirme sayfasındaki resmi doğrulama aracı imajı tarayıcıda kontrol eder. Yeni başlayan için yeterli ve önerilen yoldur.
2. **OpenPGP imza doğrulama (en sağlam):** İmajın `.img` ile `.sig` dosyasını indir, Tails imzalama anahtarıyla doğrula:

```bash
# 1. Tails imzalama anahtarını içe aktar (parmak izini resmi siteden teyit et)
wget https://tails.net/tails-signing.key
gpg --import tails-signing.key

# 2. İmzayı doğrula
gpg --verify tails-amd64-*.img.sig tails-amd64-*.img

# "Good signature from Tails developers" görmeli ve
# anahtar parmak izini resmi sitedeki ile karşılaştırmalısın.
```

> 🔑 **Püf:** OpenPGP doğrulamada anahtarın **parmak izini** mutlaka resmi siteden (mümkünse ikinci bir cihaz/ağ üzerinden) teyit et. Saldırgan hem imajı hem "doğrulama anahtarını" değiştirebilir — anahtarın güvenini siteden bağımsız doğrulamazsan doğrulama anlamsızdır. *Anahtar parmak izinin güncel değerini resmi `tails.net` dokümanından teyit et.*

### USB'ye yazma
- **Tavsiye edilen araç:** Tails'in önerdiği yazıcı (**balenaEtcher** birçok platformda; ya da mevcut bir Tails'ten **Tails Cloner**). Resmi dokümandaki güncel yöntemi izle.
- **`dd` ile (Linux/macOS, ileri seviye):**
  ```bash
  # DİKKAT: of= hedefini ÜÇ KEZ kontrol et — yanlış disk = veri kaybı
  sudo dd if=tails-amd64-*.img of=/dev/sdX bs=16M oflag=sync status=progress
  ```
- **Minimum 8 GB** USB. **Kaliteli, yeni** bir USB kullan (ucuz/eski USB'ler bozulur → açılmaz).
- **İki USB stratejisi (önerilir):** Bir USB'den ikinci bir Tails kur (Tails Cloner). Biri bozulursa/kaybolursa yedeğin olur; ayrıca güncellemeleri test edebilirsin.

> ⚠️ **Donanım gerçeği:** Tails'i yazdığın USB'yi **güvendiğin, takip ettiğin** bir ortamdan al. Birinin sana verdiği "hazır Tails USB"sine **asla güvenme** — içine implant yazılmış olabilir (Bölüm 11.4 Evil Maid).

---

<a id="4"></a>
## 4. 🚀 İlk Açılış — Welcome Screen & Ağ Bağlantısı

Bilgisayarı USB'den başlatmak için **boot menüsü** (genelde F12 / F2 / Esc / Del — markaya göre) gerekir. Gerekirse **Secure Boot** ayarını UEFI'den yönet (Tails Secure Boot'u destekler, ama bazı makinelerde ayar gerekir).

### Welcome Screen (Tails Karşılama Ekranı)
Açılışta gelen bu ekran kritik seçimler sunar — **dikkatle geç:**

- **Dil / klavye / bölge:** Klavye düzenini doğru seç (parola yazarken hayati).
- **Additional Settings (Ek Ayarlar):**
  - **Administration password (yönetici parolası):** Varsayılan olarak **yoktur** (sudo kapalı = saldırı yüzeyi düşük). Yalnızca gerçekten gerekiyorsa (ör. donanım ayarı) bu oturum için bir parola belirle. **Gereksizse belirleme.**
  - **MAC address spoofing:** Varsayılan **açık** — bırak (Bölüm 6).
  - **Network connection:** "Doğrudan Tor" mu yoksa "bridge/sansür aşma" mı (Bölüm 7).
  - **Unsafe Browser:** Varsayılan kapalı; yalnızca captive portal (otel/havalimanı giriş sayfası) için gerekirse aç — **bu tarayıcı Tor KULLANMAZ**, dikkatli ol.

### Ağ bağlantısı — Tor Connection asistanı
Welcome Screen'den sonra **Tor Connection** asistanı açılır:
- **"Connect to Tor automatically"** → normal, sansürsüz ortam.
- **"Configure a Tor bridge"** → ISP/ülke Tor'u engelliyorsa ya da Tor kullandığını gizlemen gerekiyorsa (Bölüm 7).

Tor devresi kurulana kadar bekle. Bağlantı tamamlanınca **Tor Browser** ve diğer araçları güvenle kullanabilirsin.

> 🧠 **İz felsefesi:** Welcome Screen'de yaptığın seçimler (parola, persistence) **o oturumun saldırı yüzeyini** belirler. Minimum yetki ilkesi: ihtiyacın olmayan hiçbir şeyi açma.

---

<a id="5"></a>
## 5. 🧅 Tor Entegrasyonu — Devre Yönetimi & .onion

Tails'in tüm ağı **Tor** üzerinden akar. Tor, trafiğini **üç katmanlı** (soğan gibi) şifreleyip rastgele üç düğümden (relay) geçirir:

```
   SEN ──(şifreli)──► [Guard/Giriş] ──► [Orta] ──► [Exit/Çıkış] ──► HEDEF SİTE
   │                       │                            │
   │  Gerçek IP'ni         │ Kimseyi tam               │ Senin IP'ni BİLMEZ,
   │  yalnızca Guard       │ göremez (zincir           │ ama HEDEFE giden
   │  bilir (hedefi bilmez)│  kırık)                    │ trafiği görür (HTTP ise içeriği de!)
```

- **Guard (giriş düğümü):** Senin gerçek IP'ni bilir ama nereye gittiğini bilmez. Tor, aynı guard'ı bir süre sabit tutar (güvenlik için).
- **Exit (çıkış düğümü):** Hedefe giden trafiği taşır; senin kim olduğunu bilmez **ama trafiğin içeriğini (HTTPS yoksa) görür.** → Bölüm 11.1.
- **.onion (Onion Services):** `…​.onion` adresleri exit node'a hiç çıkmaz — trafik **uçtan uca Tor içinde** kalır. Exit node riski ortadan kalkar, hedefin de IP'si gizlidir. **Mümkünse .onion sürümünü kullan** (ör. büyük gazete/SecureDrop .onion adresleri).

### Devre (circuit) yönetimi
- **Yeni kimlik / yeni devre:** Tor Browser'da **"New Tor Circuit for this site"** ile o site için devreyi yenile; **"New Identity"** ile tüm oturumu sıfırla (sekmeler kapanır, yeni devreler kurulur). Kimlik ayrımı yapacaksan kullan.
- **Onion Circuits aracı:** Tails'te hangi devrelerin açık olduğunu ve exit node'un hangi ülkede olduğunu **Onion Circuits** uygulamasıyla görebilirsin.
- **Stream isolation:** Tails, farklı uygulamaları farklı Tor devrelerine ayırır → bir uygulamanın trafiği diğeriyle aynı exit'ten çıkıp ilişkilendirilmez.

> 🔥 **Püf — devre coğrafyası:** Exit node'un ülkesi önemlidir. Belirli bir ülkenin gözetiminden çekiniyorsan **Onion Circuits**'te exit'in nerede olduğunu kontrol et; uygunsuzsa devreyi yenile. Ama unutma: exit node'u sen seçemezsin (yalnızca yenileyebilirsin), Tor rastgele atar.

---

<a id="6"></a>
## 6. 🎭 MAC Adresi Spoofing — Varsayılan & Ne Zaman Kapatılmalı

**MAC adresi**, ağ kartının donanım kimliğidir. Bir ağa (Wi-Fi/Ethernet) bağlanırsan, o ağ senin MAC adresini **görür ve loglayabilir** → "şu cihaz, şu yerde, şu saatte bağlandı" izini bırakır.

Tails **varsayılan olarak MAC adresini rastgeleleştirir (spoofing)** → bağlandığın ağ senin **gerçek donanım kimliğini** görmez. Bu, özellikle **halka açık/izlenen ağlarda** (kafe, kütüphane, kampüs) kimliğini korur.

### Ne zaman AÇIK bırakmalı (varsayılan — çoğu durum)
- **Halka açık ağ** (kafe, otel, havalimanı): Gerçek MAC'ini gizle → fiziksel cihaz kimliğin loglanmasın.
- **Sana ait olmayan / başkasının** ağı: İz bırakma.

### Ne zaman KAPATMAYI düşünmeli
- **Kendi ev ağın, MAC filtresi varsa:** Yönlendiricin yalnızca bilinen MAC'lere izin veriyorsa, spoofing bağlantıyı **engeller** → kapatman gerekebilir.
- **Bağlantı sürekli kopuyorsa:** Bazı ağlar rastgele MAC'i reddeder → tanı için geçici kapat.
- **Cihazın zaten o ağda bilindiği** ve gizlemenin anlamı olmadığı durumda (ama "rastgele MAC ev ağında dikkat çeker mi?" diye düşün — genelde çekmez).

> ⚠️ **İnce nokta:** MAC spoofing **yerel ağı** korur, **internete giden anonimliği değil** (onu Tor yapar). Ayrıca bazı ortamda **rastgele MAC'in kendisi anomali** olarak işaretlenebilir (ör. tüm cihazların sabit MAC'li olduğu kurumsal ağda). Tails bunu Welcome Screen'de seçtirir — *resmi dokümandaki güncel MAC spoofing davranışını* riskli bir ağda kullanmadan önce teyit et.

---

<a id="7"></a>
## 7. 🌉 Sansür Aşma — Bridges, obfs4, meek

Bazı ISP/ülkeler **bilinen Tor giriş düğümlerini engeller.** Ayrıca normal Tor bağlantısı, ISP'nin **"bu kullanıcı Tor'a bağlanıyor"** demesine izin verir. İkisinin de çözümü **bridge (köprü)**.

### Bridge nedir
**Bridge**, herkese açık listede **olmayan**, gizli bir Tor giriş düğümüdür. Engellenen ülkelerde Tor'a giriş kapısı sağlar. Ayrıca **pluggable transport** ile trafiği gizler.

| Transport | Ne yapar | Ne zaman |
|---|---|---|
| **obfs4** | Tor trafiğini **rastgele gürültüye** benzetir → DPI (derin paket incelemesi) "bu Tor" diyemez | **Varsayılan seçim**, çoğu sansür |
| **meek** (meek-azure) | Trafiği **büyük bir buluta** (ör. Microsoft/CDN) gidiyormuş gibi gösterir (domain fronting) → engellemek "tüm bulutu engelle" demek olur | Ağır sansür, obfs4 de engelliyse |
| **snowflake** | Geçici, gönüllü tarayıcı proxy'leri üzerinden geçer | Alternatif, dinamik |

### Tails'te bridge kullanımı
- Tor Connection asistanında **"Configure a Tor bridge"** seç.
- **obfs4** köprü adreslerini gir. Köprüleri `bridges.torproject.org` adresinden ya da `bridges@torproject.org` e-postasından (Tor'un kendi kanalları) edinebilirsin.
- Persistent Storage'da **"Tor Bridge"** özelliğini açarsan köprüler oturumlar arası saklanır (ama Bölüm 8 risklerini oku).

> 🔥 **Püf — "Tor kullandığını gizlemek":** Bridge'in iki ayrı amacı var: (1) **engeli aşmak**, (2) **Tor kullandığını ISP'den gizlemek.** İkincisi senin için kritikse (Tor kullanmanın kendisi tehlikeliyse) **obfs4/meek bridge ŞART** — düz Tor bağlantısı ISP'ye "bu kişi Tor'da" der. Ama mükemmel değil: gelişmiş bir gözlemci yine de **trafik desenlerinden** şüphelenebilir (Bölüm 11.10). *Bridge edinme kanallarının güncel halini resmi Tor dokümanından teyit et.*

---

<a id="8"></a>
## 8. 💾 Persistent Storage — Kalıcı Şifreli Depolama & RİSKLERİ

Tails'in **amnezi** sözü, "hiçbir şey kalıcı olmasın" demektir. **Persistent Storage** bunu **kısmen bozar:** USB'nin boş alanında **şifreli (LUKS) bir bölme** oluşturur ve seçtiğin verileri (parolalar, dosyalar, ayarlar) oturumlar arası **saklar.**

### Kurulum
- **Applications → Tails → Persistent Storage** (eski sürümlerde "Configure persistent volume").
- Güçlü bir **passphrase** belirle (Bölüm 9 KeePassXC ile saklamayı düşün — ama o da nerede?).
- Hangi verilerin kalıcı olacağını **tek tek seç** (toggle): KeePassXC veritabanı, GnuPG anahtarları, Tor bridge'leri, ek yazılımlar, kişisel dosyalar klasörü vb.
- Açtıktan sonra **her açılışta** Welcome Screen'de passphrase girersin; girmezsen o oturum yine amnezik olur.

### 🔥 Persistent Storage'ın getirdiği RİSKLER (en kritik bölüm)
Persistent Storage **kolaylık** sağlar ama Tails'in **iki temel güvencesini zayıflatır:**

1. **Amneziyi delersin:** Artık USB'de **kalıcı, şifreli ama VAR OLAN** veri durur. El konulursa, "bu USB'de şifreli bölme var" **görünür ve kanıtlanabilir.** Amnezinin "hiçbir iz yok" gücü kaybolur — artık "şifreli iz var, ama açamıyorsun" durumundasın.

2. **İnkâr edilebilirliği (deniability) bozar:** Saf amnezik Tails'te "bu USB'de ne var?" sorusunun cevabı **"işletim sistemi, hepsi bu"** olabilir. Persistent Storage açıkça **bir LUKS şifreli bölme** olarak görünür → **inkâr edilemez.** Birisi seni passphrase vermeye zorlarsa (rubber-hose), "şifreli bir şey var ama açmıyorsun" konumundasın. VeraCrypt hidden volume'un aksine, **Tails Persistent Storage gizli/inkâr edilebilir DEĞİLDİR** — varlığı açıktır.

3. **Saldırı yüzeyi büyür:** Kalıcı veri = kalıcı hedef. Bir oturumda bulaşan kötücül kod **diske (persistent bölmeye) yazamaz** (sistem okuma-yalnız), ama persistent'taki *verilerini* okuyabilir ya da bozabilir.

4. **Parola kaybı = veri kaybı:** LUKS passphrase'ini unutursan, persistent veri **matematiksel olarak gider** (VeraCrypt header kaybı gibi). Yedek şart.

### Karar matrisi

| İhtiyaç | Persistent Storage? |
|---|---|
| Tek seferlik, maksimum-anonim, iz-sıfır oturum | ❌ **KAPALI** — saf amnezi |
| İnkâr edilebilirlik kritik (zorla parola riski) | ❌ Persistent yerine **VeraCrypt hidden volume** (Bölüm 13) |
| Sürekli kullanım, parola/anahtar yönetimi gerekiyor | ✅ Açabilirsin — ama riskleri kabul et |
| GnuPG anahtarlarını/KeePassXC'yi taşımak | ✅ ya da ayrı şifreli USB |

> 🧠 **Operasyonel ilke:** Persistent Storage'ı **gerçekten gerekmedikçe açma.** Açacaksan, **inkâr edilebilirlik isteyen hassas veriyi oraya değil, ayrı bir VeraCrypt hidden volume'a koy.** Tails'in amnezisini bilerek delmenin bedeli vardır; bunu bilinçli yap.

---

<a id="9"></a>
## 9. 🧰 Tails İçi Araçlar — KeePassXC / GnuPG / MAT2 / Electrum

Tails, anonim operasyon için seçilmiş araçlarla gelir. En kritikleri:

### KeePassXC — Parola kasası
- Tüm takma-ad parolalarını **tek bir şifreli veritabanında** (`.kdbx`) tut → her kimlik için **uzun, rastgele, tekrarsız** parola.
- Veritabanını **Persistent Storage** ya da ayrı şifreli USB'de sakla.
- **Püf:** Anonim kimliğin tüm hesap bilgilerini ezberlemeye çalışma; KeePassXC üret + sakla → parola tekrarından doğan **korelasyon** (aynı parolayı iki yerde kullanmak) riskini ortadan kaldırır.

### GnuPG / OpenPGP — Uçtan uca şifreleme & imza
- E-posta/dosya şifreleme ve **imza** için. Gazetecilikte kaynakla güvenli yazışmanın temeli.
- Anahtarlarını Persistent Storage'da sakla (ya da her oturum yeniden içe aktar).
- **Püf:** PGP **metadata sızdırır** (kim-kime, konu, zaman). İçeriği şifreler ama "kiminle konuştuğun" görünebilir. Hassas ilişki ağı için bunu hesaba kat; mümkünse .onion tabanlı kanal + PGP birlikte.

### MAT2 — Metadata Anonimleştirme (KRİTİK)
- **MAT2 (Metadata Anonymisation Toolkit)** dosyalardaki **gizli metadata'yı temizler:** EXIF (fotoğrafın GPS konumu, kamera modeli, çekim zamanı), Office/PDF yazar adı, yazılım sürümü, düzenleme geçmişi.
- **Neden hayati:** Anonim yayınladığın bir fotoğraf, EXIF'inde **GPS koordinatını** ya da **gerçek adını** taşıyabilir → tüm anonimliği tek dosya çökertir. Çok sayıda kaynak/aktivist **bu yüzden** ifşa olmuştur.
- **Kullanım:** Dosyaya sağ tık → **"Remove metadata"**, ya da terminalde:
  ```bash
  mat2 gizli_foto.jpg        # temizlenmiş kopya üretir
  mat2 --show gizli_foto.jpg # hangi metadata var, göster
  ```
- **Püf:** MAT2 her formatı %100 temizleyemez (bilinen formatlar iyi, egzotikler eksik). En güvenlisi: **hassas görseli yeniden çek/ekran görüntüsünden üret** ya da metadata tutmayan formata çevir + MAT2'den geçir + `--show` ile **doğrula.**

### Electrum — Bitcoin (dikkatli)
- Tails, **Electrum** Bitcoin cüzdanıyla gelir; Persistent Storage'da cüzdan saklanabilir.
- **Püf — anonimlik tuzağı:** Bitcoin **anonim değil, takip edilebilir** (zincir analizi). Tails üzerinden işlem yapsan bile, **fonların kaynağı/hedefi** kimliğinle ilişkilendirilebilir. Gerçek anonimlik için zincir hijyeni (coin ayrımı, gizlilik-odaklı yöntemler) ayrı bir disiplindir. Tails IP'ni gizler, **para izini gizlemez.**

### Diğer
- **Tor Browser** (güvenlik seviyesi ayarlanabilir: Standard / Safer / Safest — JS'i kısıtlar, saldırı yüzeyini düşürür).
- **Thunderbird** (şifreli e-posta), **OnionShare** (anonim dosya paylaşımı .onion üzerinden), **Metadata cleaner**, **KeePassXC**, **LibreOffice**.

> 🔥 **Güvenlik seviyesi püfü:** Yüksek-risk işte Tor Browser'ı **"Safest"** moduna al → JavaScript çoğu sitede kapanır, saldırı yüzeyi (tarayıcı exploit'leri) **dramatik düşer.** Site bozulursa bilinçli olarak yükselt, iş bitince düşür.

---

<a id="10"></a>
## 10. 🔌 Güvenli Kapatma & RAM Temizleme (Cold-Boot Karşıtı)

Tails'in amnezisi **RAM'in silinmesine** dayanır. Ama RAM, güç kesildikten sonra **saniyeler-dakikalar** (soğutulursa daha uzun) içeriğini tutabilir → **cold-boot saldırısı** ile anahtarlar/veriler RAM'den çekilebilir. Tails buna karşı önlem alır:

- **Kapatınca RAM'i aktif siler:** Tails, kapanış sırasında **RAM'i üzerine yazarak temizler** (sürüm geliştikçe yöntem güçlenir) → cold-boot ile anlamlı veri çekmek zorlaşır. *Güncel RAM-silme davranışını resmi dokümandan teyit et.*
- **USB'yi çıkarınca otomatik kapanma:** Tails USB'sini çalışırken **fiziksel olarak çıkarırsan**, Tails **acil kapanma** yapar ve RAM'i temizler. → Bu, **BusKill / dead-man switch** mantığıdır: USB'yi bedenine bir iple bağla; biri makineyi senden koparırsa USB çıkar → Tails anında kapanır.
- **Güvenli kapatma:** İş bitince **düzgün kapat** (Shutdown). Aceleyle güç düğmesine basmak yerine, mümkünse menüden kapat (RAM-silme rutini çalışsın).

### Cold-boot disiplini
- **İş bitince hemen kapat** — makineyi açık/anonim oturumda başıboş bırakma.
- **Fiziksel tehdit anında:** USB'yi çek (acil kapanma) ya da Shutdown. Bilgisayarı **uyku/hazırda** bırakma — RAM'de veri kalır.
- **Soğuk ortam riski:** Saldırgan RAM'i sprey ile dondurursa içerik daha uzun yaşar; bu yüzden **el konulması muhtemel** ortamda fiziksel olarak makineyi koru.

> 🧠 **Felsefe:** "Güç kesilince kanıt kesilir" Tails'in özüdür — ama **ancak RAM gerçekten silinirse.** Acil durumda **USB-çek refleksi** (BusKill) + düzgün kapatma, cold-boot penceresini kapatır. Kanije Kalesi entegrasyonu (Bölüm 13) bu refleksi otomatikleştirir.

---

<a id="11"></a>
## 11. 🔥 PÜF NOKTALARI — Piyasada Bulamayacakların

Bu bölüm, çoğu rehberin atladığı ve gerçek dünyada anonimliği **çökerten** detaylardır. **En önemli bölüm budur.**

### 11.1 Exit Node Tehlikesi & HTTPS Şartı
- Tor exit node'u, hedefe giden **son** düğümdür ve **şifresiz (HTTP)** trafiğini **okuyabilir, değiştirebilir, parola çalabilir.** Tor *kim olduğunu* gizler; *ne gönderdiğini* (HTTPS yoksa) gizlemez.
- Kötü niyetli kişiler **özellikle exit node işletir** (trafik dinlemek için).
- **Çözüm:**
  - **Daima HTTPS** kullan (URL'de `https://` ve kilit). Tor Browser HTTPS'i zorlamaya çalışır ama emin ol.
  - **Hassas giriş/parola** asla HTTP üzerinden gönderme.
  - **.onion sürümünü kullan** (varsa) → trafik exit'e hiç çıkmaz, uçtan uca Tor içinde kalır → exit riski **sıfırlanır.**

### 11.2 Korelasyon / Zamanlama Saldırıları (de-anonimizasyon)
- **En ciddi yapısal tehdit.** Hem **girişini** (sen → guard) hem **çıkışını** (exit → hedef) izleyebilen bir aktör, trafiğin **zamanlama ve hacim desenini** eşleştirerek "bu giriş = bu çıkış" diyebilir → seni de-anonimize eder.
- **Global pasif rakip** (internet omurgasının büyük kısmını izleyen devlet) bunu yapabilir. Tor, tasarımı gereği bu saldırıya **tam koruma vermez** (düşük-gecikmeli ağ olmanın bedeli).
- **Çözüm (tam değil, azaltma):**
  - **Trafik desenini öngörülemez kıl:** Sabit zamanlarda, sabit hacimde aynı işi yapma.
  - **Aynı oturumda kimlik karıştırma** → ayrı kimlikler için ayrı oturum (reboot).
  - .onion servisleri korelasyon yüzeyini azaltır (çıkış yok).
  - Bunun **çözümü yok**, yalnızca yönetimi var; yüksek-riskte bunu kabul et.

### 11.3 Persistent Storage İnkâr Edilebilirliği Bozar
- (Bölüm 8'in özeti — burada bir kez daha, çünkü en sık atlanan tuzak.) Persistent Storage **açıkça görünen** bir LUKS şifreli bölmedir. El konulduğunda "şifreli veri var ama açmıyorsun" konumuna düşersin — **inkâr edilemez.** Saf amnezik Tails'in "hiçbir kalıcı veri yok" gücünü kaybedersin.
- **Çözüm:** İnkâr gereken veri için Persistent yerine **VeraCrypt hidden volume** (matematiksel olarak gizli — Bölüm 13).

### 11.4 BIOS / UEFI / Firmware Tehdidi (Evil Maid)
- Tails okuma-yalnız bir USB'den çalışsa bile, **Tails başlamadan önce çalışan kod** (BIOS/UEFI, firmware, kötücül bootloader) seni izleyebilir, tuş kaydedebilir, Tails'i değiştirebilir. Buna **Evil Maid** denir: sen yokken biri (otel temizlikçisi metaforu) makineye/USB'ye fiziksel dokunur.
- Tails bu katmanın **altında** çalışamaz — firmware'i temizleyemez.
- **Çözüm:**
  - **Güvendiğin, fiziksel kontrolündeki** bir bilgisayar kullan; başıboş bırakma.
  - **Sana verilen "hazır Tails USB"sine güvenme.**
  - Firmware/Secure Boot bütünlüğüne dikkat et; mümkünse kendi temiz donanımın.
  - Yüksek-riskte: işten önce/sonra fiziksel kurcalanma izi (tamper-evidence) kontrol et.

### 11.5 Kamera / Mikrofon — Fiziksel Kapatma
- Tails yazılımsal koruma sağlar ama bir tarayıcı/uygulama exploit'i **kamera/mikrofona** erişebilir → yüzünü/sesini/ortamını sızdırır → anonimliği anında çökertir.
- **Çözüm:**
  - **Kamerayı fiziksel kapat** (kapak/bant).
  - **Mikrofon** için: mümkünse fiziksel kesme (bazı laptoplarda donanım anahtarı), ya da dahili mikrofonu olmayan/kapalı makine. Yazılım izni yetmez — **fiziksel** güven verir.

### 11.6 Saat / Zaman Dilimi Sızıntısı
- Sistemin **zaman dilimi ve saat sapması**, bulunduğun coğrafyaya dair ipucu verebilir. Tor, korelasyon için zamanlamaya duyarlıdır.
- Tails **saati Tor üzerinden senkronize eder ve UTC kullanma eğilimindedir** → yerel saat dilimi sızıntısını azaltır. *Güncel saat/zaman dilimi davranışını resmi dokümandan teyit et.*
- **Çözüm:** Welcome Screen'de gereksiz yere yerel zaman dilimi ayarlama; Tails'in varsayılan saat yönetimine güven. Uygulama seviyesinde "yerel saat" gösteren bir şeyle anonim kimliğini ifşa etme (ör. "benim saatimle 3'te" deme).

### 11.7 Dosya Metadata (EXIF) — Sessiz İfşa
- (Bölüm 9 MAT2'nin önemi.) Fotoğraf/PDF/Office dosyaları **GPS, gerçek ad, cihaz, zaman** taşır. Anonim yayınlanan bir dosya, metadata'sıyla seni ele verir.
- **Çözüm:** Paylaşmadan önce **HER dosyayı MAT2'den geçir** ve `mat2 --show` ile **temiz olduğunu doğrula.** Ekran görüntüsü genelde daha temizdir ama yine de kontrol et.

### 11.8 Ekran Görüntüsü & Pano (Clipboard) Sızıntısı
- **Ekran görüntüsü** aldığın anda, ekrandaki **her şey** (açık başka sekme, kimlik, isim) görüntüye girer → paylaşırsan sızar.
- **Pano (kopyala-yapıştır):** Bir kimlikten kopyaladığını yanlışlıkla başka kimliğin alanına yapıştırmak → korelasyon. Ayrıca bazı uygulamalar pano içeriğini loglayabilir.
- **Çözüm:** Ekran görüntüsü almadan önce **ekranı temizle** (tek pencere, kimlik gizli). Panoyu kimlikler arası taşıma; hassas veriyi kullandıktan sonra panoyu temizle.

### 11.9 Aynı Kimlikle Çoklu Oturum Korelasyonu
- Aynı takma adı/parolayı/yazı stilini/zamanlamayı **farklı bağlamlarda** kullanmak → gözlemci bunları birbirine bağlar → "bu hesaplar aynı kişi" → ardından gerçek kimliğe köprü.
- **Stylometry (yazı stili analizi):** Cümle yapın, kelime tercihlerin bile parmak izi olabilir.
- **Çözüm:**
  - **Her kimlik için ayrı oturum** (reboot — Tails amnezisi bunu kolaylaştırır), ayrı parola (KeePassXC), ayrı Tor New Identity.
  - Kimlikleri **asla** aynı oturumda yan yana kullanma.
  - Yazı stilini bilinçli değiştir; tanınır deyim/imza kullanma.

### 11.10 "Tails/Tor Kullandığın" Görünebilir
- Bridge kullanmazsan ISP'n **"bu kişi Tor'a bağlandı"** der. Bazı ortamda **Tor kullanıyor olman bile** seni şüpheli/hedef yapar (az sayıda Tor kullanıcısının olduğu bir ağda dikkat çekersin).
- **Çözüm:**
  - **obfs4 / meek bridge** kullan → trafik "Tor'a benzemez", ISP Tor dediğini söyleyemez (Bölüm 7).
  - Yine de **mükemmel değil:** gelişmiş DPI/desen analizi şüphelenebilir. "Tor kullandığım görünmemeli" tehdidi varsa bunu **operasyonun en başında** planla, sonradan değil.

### 11.11 Unsafe Browser & Captive Portal Tuzağı
- Tails'teki **Unsafe Browser**, captive portal (otel/havalimanı giriş sayfası) içindir ve **Tor KULLANMAZ** — gerçek IP'nle çalışır.
- **Tuzak:** Yanlışlıkla hassas bir işi Unsafe Browser'da yaparsan **anonimliğin tamamen çöker** (gerçek IP'n ifşa olur).
- **Çözüm:** Unsafe Browser'ı **yalnızca** captive portal'a giriş için aç, işi bitince **kapat.** Hassas hiçbir şeyi orada açma. Varsayılan kapalıdır — gereksizse hiç açma.

### 11.12 Güncellemeyi İhmal Etmek
- Tails ve Tor Browser düzenli güvenlik güncellemesi alır; eski sürümde **bilinen exploit'ler** seni de-anonimize edebilir.
- **Çözüm:** Tails'i güncel tut (otomatik yükseltici ya da yeni imaj + imza doğrulama). Güncellemeyi **imza doğrulayarak** yap — güncelleme anı tedarik-zinciri saldırısının en sevdiği andır (VeraCrypt rehberindeki aynı ilke).

---

<a id="12"></a>
## 12. ☠️ Yaygın Ölümcül Hatalar

1. **İmaj imzasını doğrulamadan kurmak** → truva atlı Tails ile bizzat gözetlenmek. (En temel felaket.)
2. **HTTP (şifresiz) üzerinden hassas veri/parola göndermek** → exit node okur. **HTTPS/.onion şart.**
3. **Persistent Storage'ı bilinçsiz açmak** → amneziyi delmek, inkâr edilebilirliği kaybetmek.
4. **Aynı kimliği/parolayı çoklu bağlamda kullanmak** → korelasyonla de-anonimize olmak.
5. **Dosya metadata'sını (EXIF) temizlememek** → fotoğrafın GPS'i ya da gerçek adın sızar (MAT2 atlanır).
6. **Unsafe Browser'da hassas iş yapmak** → gerçek IP ifşası.
7. **Gerçek kimlikle hesaba girmek / kişisel bilgi yazmak** → Tor'un anonimliğini kendi elinle çürütmek.
8. **Tor kullandığını gizlemen gerekirken bridge kullanmamak** → ISP'ye Tor kullandığını söylemek.
9. **Başkasının verdiği / başıboş bıraktığın makine/USB** → Evil Maid / firmware implantı.
10. **Kamera/mikrofonu fiziksel kapatmamak** → exploit'le yüz/ses sızması.
11. **İş bitince düzgün kapatmamak / uyku bırakmak** → cold-boot ile RAM'den veri çekilmesi.
12. **Bitcoin'i anonim sanmak** → zincir analiziyle para izinden kimliğe ulaşılması.
13. **Eski sürüm kullanmak** → bilinen exploit'lerle de-anonimizasyon.

---

<a id="13"></a>
## 13. 🏰 Kanije Kalesi ile Birlikte Kullanım

Bu repo (Kanije Kalesi), **fiziksel tehdit anında** cihazı uzaktan/otomatik koruyan bir muhafızdır. Tails **anonim ve amnezik oturumu** sağlar; Kanije **fiziksel olay anını** yönetir. İkisinin felsefesi örtüşür: **RAM-only, iz bırakma, güç kesilince kanıt kesilir.**

| Senaryo | Tails rolü | Kanije Kalesi rolü |
|---|---|---|
| Anonim oturum, iz-sıfır | Amnezi + Tor zorunlu | — |
| Sen uzaktayken biri makineye yaklaştı | Açık oturum → kimlik riski | `/koruma` dead-man / yanlış-giriş → **kilitle + alarm + foto** |
| USB'yi biri çekti / makine koparıldı | **Acil kapanma + RAM temizliği** | USB dead-man tetikleyici → alarm + bildirim |
| Acil durum | RAM-only zaten iz tutmaz | `/panik`, `/kilit tam` (lockdown) |
| Hassas dosyanın inkâr edilebilirliği | Persistent **inkâr edilemez** → kullanma | VeraCrypt hidden volume'u tut + Kanije ile koru |
| Adli inceleme öncesi | Kapanış = RAM silindi, iz yok | İz temizleme, RAM-only mod örtüşmesi |

### 🔥 Önerilen entegrasyon deseni
1. **Anonim/hassas işi Tails'te yap** — amnezi + Tor + (gerekirse obfs4 bridge).
2. **İnkâr gereken veriyi Persistent Storage'a DEĞİL, VeraCrypt hidden volume'a koy** → Tails'te bu volume'u mount et (VeraCrypt rehberi Bölüm 13'teki entegrasyon). Hidden volume matematiksel olarak gizli; Persistent Storage açıkça görünür → inkâr için VeraCrypt kazanır.
3. **USB dead-man / BusKill refleksi:** Tails USB'sini bedenine bağla. Biri makineyi koparırsa → Tails **acil kapanma + RAM temizliği** yapar. Kanije'nin USB dead-man'i ile birleştir → hem RAM uçar hem **alarm + foto + uzak bildirim** gider.
4. **Fiziksel muhafız:** Sen Tails oturumundayken kalkman gerekirse, Kanije `/koruma` (dead-man switch + yanlış-giriş) devrede olsun → sen yokken biri dokunursa kilitlenir, kanıt toplanır.
5. **RAM-only örtüşmesi:** Tails diske dokunmaz; Kanije'nin RAM-only modu da öyle. **Güç kesildiğinde her iki katmandan da hiçbir iz kalmaz** (Tails/BusKill sınıfı duruş).

> 🧠 **Felsefe örtüşmesi:** Tails "yazılım katmanında" amnezi ve anonimlik verir; Kanije "fiziksel katmanda" tehdidi algılayıp tepki verir. **Tails seni ağda görünmez kılar; Kanije seni masanın başında korumasız bırakmaz.** İkisi, "dijital + fiziksel" eksiksiz bir RAM-only kale kurar.

---

<a id="14"></a>
## 14. ✅ Hızlı Referans & Operasyonel Kontrol Listesi

### Kurulumda
- [ ] İmaj **yalnızca `tails.net`**'ten indirildi
- [ ] **İmza/bütünlük doğrulandı** (OpenPGP ya da web doğrulama)
- [ ] Anahtar **parmak izi bağımsız teyit** edildi
- [ ] **Kaliteli, güvendiğin** USB'ye yazıldı
- [ ] **Yedek USB** (ikinci Tails) hazır

### Her açılışta (Welcome Screen)
- [ ] Doğru **klavye/dil**
- [ ] **Admin parolası gereksizse belirlenmedi**
- [ ] **MAC spoofing açık** (halka açık ağda)
- [ ] Gerekiyorsa **bridge (obfs4/meek)** seçildi
- [ ] **Unsafe Browser kapalı** (gerekmedikçe)

### Her oturum
- [ ] **HTTPS / .onion** kullanılıyor (HTTP'den hassas veri yok)
- [ ] Tor Browser güvenlik seviyesi **uygun** (yüksek-riskte "Safest")
- [ ] **Kamera/mikrofon fiziksel kapalı**
- [ ] **Her kimlik için ayrı oturum** (reboot) / New Identity
- [ ] Paylaşılan **her dosya MAT2'den geçti** (`--show` ile doğrulandı)
- [ ] Ekran görüntüsü/pano **kimlik sızdırmıyor**
- [ ] Gerçek kimlikle hesaba **girilmedi**, kişisel bilgi yazılmadı

### Persistent Storage (kullanıyorsan)
- [ ] **Gerçekten gerekli** olduğu için açıldı
- [ ] **Güçlü passphrase** + **yedek**
- [ ] **İnkâr gereken veri burada DEĞİL** (VeraCrypt hidden volume'da)
- [ ] Risklerin (amnezi/deniability kaybı) **kabul edildi**

### Kapatırken
- [ ] İş bitince **düzgün Shutdown** (RAM silinsin)
- [ ] **Uyku/hazırda bırakılmadı**
- [ ] Fiziksel tehditte **USB-çek** (acil kapanma) refleksi hazır
- [ ] Tails **güncel** tutuluyor (imza doğrulayarak)

---

<a id="15"></a>
## 15. ⚖️ Hukuki Sınır & Sınır Geçişi Notu

- **Tor/Tails kullanmanın görünürlüğü:** Bazı yargı bölgelerinde/ağlarda **Tor kullanıyor olman bile** dikkat çeker ya da şüphe doğurur. Tor kullanmanın kendisini gizlemen gerekiyorsa **obfs4/meek bridge** kullan (Bölüm 7) — ama bunun da mükemmel olmadığını bil (Bölüm 11.10).
- **Sınır geçişi:** Cihazlar ve USB'ler sınırda incelenebilir/kopyalanabilir. **Persistent Storage'lı bir Tails USB**, "şifreli bir bölme var" diye **görünür ve inkâr edilemez** → bazı ülkelerde parola talep edilebilir, reddetmek suç sayılabilir (VeraCrypt rehberi Bölüm 19). Daha güvenli: **saf amnezik Tails** (kalıcı veri yok) + hassas veriyi geçtikten sonra **.onion/şifreli kanaldan indir.**
- **İnkâr edilebilirlik:** Tails Persistent Storage **gizli değildir** (varlığı bellidir). İnkâr gerekiyorsa **VeraCrypt hidden volume** tercih et (matematiksel gizlilik). İkisini karıştırma.
- **Amnezi ≠ dokunulmazlık:** Tails yerel iz bırakmaz, ama **ağ tarafı** (korelasyon, exit gözetimi, "Tor'a bağlandın" kaydı) ayrı bir gözetim yüzeyidir. Operasyonel disiplin (Bölüm 11) bu yüzeyi yönetmek içindir.
- **Yasal sorumluluk:** Bu rehber **meşru gizlilik, gazetecilik, insan hakları savunuculuğu, araştırma ve kişisel veri koruma** içindir. Bulunduğun yargı bölgesinin yasalarına uy.

---

> 🏰 **Kapanış:** Anonimlik bir ürün değil, bir **disiplindir.** En güçlü Tor devresi bile, gerçek adınla giriş yaptığın bir hesapta ya da metadata'sını temizlemediğin bir fotoğrafta çaresizdir. Tails sana **amnezik, Tor-zorunlu bir kale** verir; **kimliğini gizli tutmak senin işin.** Exit node'a şifresiz veri yollarsan, persistence'la amneziyi delersen ya da iki kimliği yan yana açarsan — kale değil, cam ev olur. Kanije Kalesi de tam burada — sen masadan kalktığında kapıyı kilitleyen, USB koparıldığında RAM'i temizleten nöbetçi — devreye girer.
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. İlgili: `VERACRYPT_USTALIK_REHBERI.md` (durağan veri + hidden volume), `LINUX_HARDENING_KALE.md`, `WINDOWS11_HARDENING_KALE.md`, `SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md`, `DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md`.*
