# 🧊 QUBES OS — TAM BÖLMELEME & İZOLASYON USTALIK REHBERİ
## "Security by Compartmentalization" Felsefesinden Split-GPG'ye, Püf Noktalarıyla Uçtan Uca

> **Amaç:** Qubes OS'u "kur-tıkla-kullan" seviyesinden çıkarıp, bir tehdit aktörü karşısında **gerçekten ayakta kalacak** şekilde kullanmayı öğretmek. Qubes'in vaadi tek cümlede özetlenir: **"Bir şey ele geçirilecekse, her şey değil sadece o ele geçirilsin."** Bu rehber yalnızca *nasıl*'ı değil, **neden** ve **hangi durumda işe yaramaz**'ı da anlatır. Forum cevaplarında bulamayacağın dom0 hijyeni, pano saldırı yüzeyi, TemplateVM disiplini, Split-GPG mantığı ve kimlik-bazlı qube ayrımı burada.

> ⚠️ **Önce bunu oku:** Qubes, *yanlış kullanıldığında* sana **yanlış güvenlik hissi** verir — her şeyi tek bir AppVM'de toplayan bir Qubes kurulumu, sıradan bir Linux'tan daha güvenli değildir; sadece daha yavaştır. İzolasyon, **senin onu doğru kurduğun kadar** vardır. Özellikle **dom0 hijyeni**, **TemplateVM mantığı** ve **pano/dosya transferi saldırı yüzeyi** bölümlerini atlama.

> 🧪 **Teyit notu:** Qubes hızlı gelişen bir projedir; sürümler arası varsayılan qube adları, komut isimleri ve davranışlar değişebilir. Aşağıdaki komutlar Qubes 4.1/4.2 mantığına dayanır. **Kendi sürümünde `qubes-os.org` resmi dokümanından teyit et** — yanlış güvenlik tavsiyesi tehlikelidir, bu yüzden emin olmadığım yerleri "teyit et" notuyla işaretledim.

---

## 📑 İÇİNDEKİLER

1. [Qubes OS Nedir, Neden? (VM vs Konteyner vs Tek Sistem)](#1)
2. [Tehdit Modeli — Neyi Korur, Neyi KORUMAZ](#2)
3. [Mimari — dom0, sys-net, sys-firewall, sys-usb, AppVM, TemplateVM, DisposableVM](#3)
4. [Kurulum + Donanım Gereksinimleri (VT-x/VT-d/IOMMU) + İmza Doğrulama](#4)
5. [🔥 TemplateVM Mantığı — Kök Şablon → Türetilen AppVM'ler](#5)
6. [Renk Kodlu Pencere Kenarlıkları — Görsel Güven Ayrımı](#6)
7. [DisposableVM — Tek Kullanımlık İzolasyon](#7)
8. [Qube'lar Arası Güvenli Transfer (qvm-copy / Pano) & Riskleri](#8)
9. [sys-usb ile USB İzolasyonu (BadUSB Savunması)](#9)
10. [Whonix Entegrasyonu — Tor ile Ağ Seviyesinde Anonimlik](#10)
11. [Split-GPG — Özel Anahtarı Asla Dışarı Çıkarmamak](#11)
12. [Split-SSH — Aynı Mantık, SSH İçin](#12)
13. [USB Qube, Kamera & Mikrofon Kontrolü](#13)
14. [🔥 PÜF NOKTALARI — Piyasada Bulamayacakların](#14)
15. [Yedekleme & Kurtarma (qvm-backup şifreli)](#15)
16. [Yaygın Ölümcül Hatalar](#16)
17. [🏰 Kanije Kalesi ile Birlikte Kullanım](#17)
18. [Hızlı Referans (qvm-* komutları) & Operasyonel Kontrol Listesi](#18)
19. [Hukuki Sınır & Operasyonel Notlar](#19)

---

<a id="1"></a>
## 1. 🧭 Qubes OS Nedir, Neden?

Qubes OS, **"security by compartmentalization (bölmeleme yoluyla güvenlik)"** ilkesi üzerine kurulu, ücretsiz ve açık kaynaklı bir işletim sistemidir. Joanna Rutkowska ve Rafał Wojtczuk tarafından başlatılmıştır (Invisible Things Lab). Temel fikir şudur: **Tek bir devasa, güvenilir işletim sistemi yerine, her işi kendi izole sanal makinesinde (qube) çalıştır.** Bir qube ele geçirilse bile saldırgan diğerlerine — ve en önemlisi yönetim katmanına — geçemez.

Qubes bir Linux dağıtımı **değildir**; bir **Xen hipervizörü** üzerine kurulu bir *meta-işletim sistemi*'dir. Sen "Qubes kullanırken" aslında onlarca küçük VM'i (Fedora, Debian, Whonix tabanlı) tek bir bütünleşik masaüstünde yönetirsin. Pencereler farklı qube'lardan gelir ama tek bir ekranda, renk kodlu kenarlıklarla birlikte yaşar.

### Felsefe: "Reasonably Secure"

Qubes kendini **"reasonably secure (makul ölçüde güvenli)"** olarak tanımlar — "kırılamaz" değil. Bu dürüstlük önemlidir:

- **Mutlak güvenlik yoktur.** Qubes, *hata yapacağını varsayar* — bir uygulama, bir tarayıcı sekmesi, bir PDF eki bir gün seni ele geçirecek. Soru "ele geçirilir miyim" değil, **"ele geçirildiğimde hasar ne kadar yayılır"**dır.
- Qubes bu hasarı **tek bir qube ile sınırlar.** Bankacılık ayrı qube'da, gündelik gezinti ayrı qube'da, iş ayrı qube'da → birinin çökmesi diğerini etkilemez.

### Neden VM, neden konteyner değil? (Dürüst Kıyas)

| Yaklaşım | İzolasyon sınırı | Saldırı yüzeyi | Qubes'in tercihi |
|---|---|---|---|
| **Tek sistem (normal OS)** | Yok — her şey aynı çekirdek/kullanıcı | Devasa | ❌ Bir bulaşma = her şey |
| **Konteyner (Docker/LXC)** | Paylaşılan çekirdek (namespace) | Çekirdek 0-day ile aşılır | ❌ Çekirdek tek hata noktası |
| **VM (Xen hipervizör)** | Donanım destekli (VT-x/VT-d) | Çok daha küçük (ince hipervizör) | ✅ Qubes'in temeli |
| **Ayrı fiziksel makineler** | Mükemmel ama pratik değil | — | Qubes bunu *tek makinede* taklit eder |

**CTI / hassas-veri perspektifi:** Konteynerler *aynı çekirdeği* paylaşır; bir çekirdek açığı tüm konteynerleri düşürür. Qubes, izolasyonu **donanım sanallaştırması (Xen + Intel VT-x/VT-d)** ile sağlar — saldırganın bir qube'dan kaçması için **hipervizör veya IOMMU'yu aşması** gerekir, bu çok daha yüksek bir çıtadır. CTI analisti için bu şu demektir: **şüpheli bir örneği (malware) bir qube'da patlatabilir, ağ qube'unu izole tutabilir, iş qube'una asla bulaştırmadan analiz edebilirsin.**

> 💡 **Zihinsel model:** Qubes'i "tek makinede çalışan bir ağ" gibi düşün. Her qube ayrı bir bilgisayar; aralarında **varsayılan olarak hiçbir bağlantı yok**; ne paylaşacaklarını **sen açıkça** belirlersin (dosya kopyalama, pano, ağ). VeraCrypt diski şifreler; Qubes **çalışan sistemi bölmelere ayırır** — farklı katmanlar, birlikte güçlü.

---

<a id="2"></a>
## 2. 🎯 Tehdit Modeli — Neyi Korur, Neyi KORUMAZ

Qubes kurmadan önce **kime karşı** korunduğunu netleştir. İzolasyonun gücü, neyi izole ettiğini doğru anlamandan gelir.

### ✅ Qubes'in KORUDUĞU senaryolar

- **Tek qube'un ele geçirilmesi:** Gündelik tarayıcında bir 0-day exploit çalışsa bile, saldırgan o qube'da kalır — bankacılık qube'una, iş belgelerine, GPG anahtarına **erişemez.** Bu, Qubes'in en güçlü olduğu senaryodur.
- **Malware analizi:** Şüpheli bir dosyayı DisposableVM'de açarsın; iş biter, qube yok edilir, bulaşma onunla gider. Ana sisteme **hiç dokunmaz.**
- **Kimlik ayrımı:** İş kimliğin, kişisel kimliğin, anonim kimliğin ve bankacılık kimliğin ayrı qube'larda → biri sızsa diğerleri açığa çıkmaz; aralarında *korelasyon* kurulamaz.
- **Kötücül USB (BadUSB):** USB cihazı `sys-usb` qube'una takılır; kötücül bir USB denetleyici dom0'ı değil, yalnızca izole `sys-usb`'yi hedefler.
- **Ağ saldırı yüzeyi izolasyonu:** Ağ donanımı (Wi-Fi/Ethernet) `sys-net`'te izole; bir ağ sürücüsü açığı ele geçse bile `sys-net`'te kalır, asıl verine ulaşamaz.

### ❌ Qubes'in KORUMADIĞI senaryolar

- **Hipervizör 0-day (VM escape):** Xen'de kritik bir açık varsa, bir qube'dan kaçıp dom0'a ya da başka qube'a geçilebilir. Nadirdir ama **mümkündür** — Qubes'in tek en büyük teorik zaafı budur.
- **Donanım / firmware / Intel ME:** Qubes işletim sistemi katmanında çalışır. **Intel Management Engine (ME)**, BIOS/UEFI implantı, kötücül firmware ya da donanım arka kapısı, Qubes'in *altında* çalışır → Qubes bunu göremez/durduramaz. (Teyit: Coreboot + ME devre dışı bırakma kısmi azaltma sağlar; bkz. Bölüm 14.)
- **Yan kanal saldırıları (side-channel):** Spectre/Meltdown sınıfı CPU açıkları, izolasyon sınırlarını teorik olarak delebilir. Qubes mikrokod/yamalarla azaltır ama CPU donanım kusurlarını tam kapatamaz.
- **Kullanıcı hatası — YANLIŞ QUBE'A YAPIŞTIRMA:** En sık ve en ölümcül zaaf. Bir parolayı/dosyayı **yanlış qube'un panosuna** yapıştırırsan, izolasyonu *kendi elinle* delersin. Qubes seni bundan koruyamaz — bu senin disiplinine bağlıdır.
- **dom0'ı kirletmek:** dom0'da internet açmak, dosya indirmek, yabancı yazılım kurmak → tüm güvenlik modelini çökertir (Bölüm 3, 14).
- **Fiziksel erişim (kapalı değilken):** Açık/kilitli bir makineye fiziksel erişen biri cold-boot, DMA (IOMMU yoksa) ya da Evil Maid uygulayabilir.

> 🧠 **Altın kural:** Qubes, **"bir hatanın yayılmasını"** engeller — hatanın *olmasını* değil. İzolasyon, hipervizörün sağlamlığı + senin operasyonel disiplinin kadar güçlüdür. En zayıf halka neredeyse her zaman **dom0 hijyeni** ve **yanlış qube'a yapıştırma**dır.

---

<a id="3"></a>
## 3. 🧱 Mimari — Qube Türleri ve Sistem İskeleti

Qubes'in tüm gücü mimarisindedir. Şu şemayı zihnine kazı:

```
                        ┌───────────────────────────────────────┐
                        │              FİZİKSEL DONANIM          │
                        │   CPU (VT-x/VT-d) · RAM · Disk · NIC   │
                        └───────────────────┬───────────────────┘
                                            │
                        ┌───────────────────▼───────────────────┐
                        │           XEN HİPERVİZÖR (ince)        │
                        └───────────────────┬───────────────────┘
                                            │
        ┌───────────────────────────────────┼───────────────────────────────────┐
        │                                   │                                   │
┌───────▼────────┐                 ┌────────▼─────────┐               ┌─────────▼────────┐
│     dom0       │                 │  SİSTEM QUBE'LARI │               │  KULLANICI QUBE'LARI│
│  (YÖNETİM)     │                 │ ┌──────────────┐  │               │ ┌────────────────┐ │
│  ⚫ İNTERNET YOK │◄── GUI olayları │ │   sys-net    │◄─┼── donanım NIC │ │  iş (AppVM)    │ │
│  ⚫ DOSYA AÇMA  │   (sadece çizim) │ │ (ağ donanımı)│  │               │ │  kişisel(AppVM)│ │
│     YOK        │                 │ └──────┬───────┘  │               │ │  banka (AppVM) │ │
│  Pencere yön.  │                 │ ┌──────▼───────┐  │               │ │  anon (Whonix) │ │
│  qvm-* araçları│                 │ │ sys-firewall │◄─┼───────────────┼─┤  (ağ buradan)  │ │
└────────────────┘                 │ │  (ağ kuralı) │  │               │ └────────────────┘ │
                                   │ └──────────────┘  │               └────────────────────┘
                                   │ ┌──────────────┐  │               ┌────────────────────┐
                                   │ │   sys-usb    │◄─┼── USB denetleyici│ │  TemplateVM(ler)  │
                                   │ │ (USB izole)  │  │               │ │ fedora · debian   │
                                   │ └──────────────┘  │               │ │ whonix-gw/ws      │
                                   └───────────────────┘               │ │ (kök şablonlar)   │
                                                                       │ └────────────────────┘
                                                                       │ ┌────────────────────┐
                                                                       │ │   DisposableVM     │
                                                                       │ │ (tek kullanımlık)  │
                                                                       │ └────────────────────┘
                                                                       └────────────────────────
```

### dom0 — Yönetim Alanı (Kalenin En İç Burcu)

dom0, Xen'in ayrıcalıklı yönetim domaini; tüm grafik arayüzü (pencere yöneticisi) ve `qvm-*` yönetim araçları buradadır. **Mutlak kural:**

- **dom0'da ASLA internet yok.** dom0'ın ağ erişimi *yoktur ve olmamalıdır.* İnternete bağlı bir dom0 = tüm sistemin tek hata noktası.
- **dom0'da ASLA dosya açma.** İndirdiğin hiçbir dosyayı, PDF'i, belgeyi dom0'da açma. dom0 ele geçerse **her qube biter** — çünkü dom0 hepsini yönetir.
- **dom0 minimal kalsın.** Yabancı yazılım kurma, depo ekleme, script indirme. dom0 yalnızca pencereleri çizer ve qube'ları yönetir.

> ⚠️ **dom0 = krallık.** Diğer her qube feda edilebilir; dom0 edilemez. dom0'a bulaşan bir şey, *tüm bölmeleme modelini* anlamsız kılar. Qubes'in birinci emri: **dom0'ı kutsal tut.**

### Sistem Qube'ları (sys-*)

| Qube | Görevi | Neden ayrı? |
|---|---|---|
| **sys-net** | Fiziksel ağ donanımını (Wi-Fi/Ethernet) tutar; internete *doğrudan* bağlanan tek qube | Ağ donanımı/sürücüleri saldırı yüzeyidir; bir açık ele geçse `sys-net`'te hapsolur, verine ulaşamaz |
| **sys-firewall** | `sys-net` ile diğer qube'lar arasında güvenlik duvarı; hangi qube nereye bağlanabilir kuralları | Ağ politikasını tek noktadan yönetir; `sys-net` çökse bile firewall katmanı durur |
| **sys-usb** | Tüm USB denetleyicilerini tutar; USB cihazları önce buraya takılır | BadUSB/kötücül USB denetleyici dom0 yerine izole `sys-usb`'yi hedefler |

**Akış:** `sys-net` (donanım) → `sys-firewall` (kural) → senin AppVM'lerin. Her AppVM internete **doğrudan değil**, bu zincir üzerinden çıkar. Bir AppVM'in "NetVM"ini `none` yaparsan o qube **tamamen çevrimdışı** olur (örn. soğuk anahtar deposu, hava boşluğu taklidi).

### AppVM (Uygulama Qube'u)

Senin günlük işlerini yaptığın qube'lar. Her biri bir **TemplateVM**'den türetilir (kök sistem ondan gelir) ama **kendi `/home`'una sahiptir** (kalıcı kişisel veri). "iş", "kişisel", "banka", "anon" gibi kimliklere göre ayır (Bölüm 14).

### TemplateVM (Şablon Qube'u)

AppVM'lerin kök dosya sistemini sağlayan **salt-okunur kök**. Yazılım buraya kurulur, güncelleme buradan yapılır, türetilen tüm AppVM'ler bundan beslenir (Bölüm 5). Doğrudan TemplateVM'de çalışmazsın — onu yalnızca *bakım* için açarsın.

### DisposableVM (Tek Kullanımlık Qube)

Açılır, iş görülür, **kapanınca tüm izleriyle yok olur.** Şüpheli ek/link açmanın, malware patlatmanın en güvenli yolu (Bölüm 7).

---

<a id="4"></a>
## 4. 📥 Kurulum + Donanım Gereksinimleri + İmza Doğrulama

### Donanım gereksinimleri (atlanırsa Qubes ya çalışmaz ya da güvensiz çalışır)

| Gereksinim | Neden ZORUNLU | Teyit yöntemi |
|---|---|---|
| **Intel VT-x / AMD-V** | Donanım sanallaştırma — Xen'in qube'ları çalıştırması için | BIOS/UEFI'de etkin olmalı |
| **Intel VT-d / AMD-Vi (IOMMU)** | Donanım izolasyonu + DMA koruması; `sys-net`/`sys-usb`'yi gerçekten izole eder | **Kritik!** Olmazsa `sys-net` DMA ile dom0'ı vurabilir |
| **64-bit CPU** | Xen 64-bit | — |
| **RAM** | **Minimum 6 GB, gerçekçi 16 GB+** | Her qube RAM yer; çok qube açacaksan 16-32 GB hedefle |
| **Disk** | **Minimum 32 GB, gerçekçi 256 GB+ SSD** | Şablonlar + qube'lar şişer; SSD performans için şart |
| **UEFI/Legacy uyumu** | Önyükleme | HCL'den kontrol et |

> 🔥 **HCL (Hardware Compatibility List) — ATLAMA:** Qubes her donanımda sorunsuz çalışmaz (özellikle Wi-Fi, GPU, suspend). **Satın almadan/kurmadan önce `qubes-os.org/hcl` listesinden modelini kontrol et.** "Çalışıyor" raporu olan bir makine, saatlerce uğraşı önler. IOMMU (VT-d) **olmayan** bir makinede Qubes kurma — izolasyonun büyük kısmı kağıt üzerinde kalır.

### IOMMU/VT-d gerçeği (en çok atlanan teknik nokta)

VT-x qube'ları *çalıştırır*, ama asıl izolasyonu **VT-d (IOMMU)** sağlar. IOMMU, `sys-net` gibi donanıma dokunan qube'ların **DMA (Direct Memory Access) ile dom0'ın belleğini okumasını/yazmasını engeller.** IOMMU yoksa: ele geçirilmiş bir `sys-net`, doğrudan dom0 belleğine DMA yapıp tüm izolasyonu delebilir. **VT-d'siz Qubes, "izole görünen ama aslında delik" bir sistemdir.**

### İmza doğrulama (kurulumun en kritik adımı — ATLAMA)

İndirdiğin Qubes ISO'su değiştirilmişse, *kurduğun ilk andan itibaren* arka kapılı bir sistem kullanırsın. Bu yüzden imzayı doğrula:

```bash
# 1. Qubes Master Signing Key'i içe aktar ve parmak izini RESMİ siteden teyit et
gpg --import qubes-master-signing-key.asc
# Parmak izini qubes-os.org üzerinden el ile karşılaştır (kritik!)

# 2. ISO'yu imzalayan release key'i içe aktar / doğrula
gpg --import qubes-release-X-signing-key.asc

# 3. ISO'nun imza dosyasını doğrula
gpg --verify Qubes-RX.X-x86_64.iso.asc Qubes-RX.X-x86_64.iso
# "Good signature" + doğru release key parmak izi görmelisin
```

> 🔑 **Püf:** Master Signing Key parmak izini **birden çok bağımsız kaynaktan** (resmi site + arşiv + güvendiğin bir ayna) teyit et. Tedarik zinciri saldırısı tam burada, "anahtarı da sahteleyerek" işler. ISO'yu **güvenli, temiz bir makinede** doğrula. Kuruluş diskini yazarken (Rufus/dd) imzası doğrulanmış imajı kullan.

---

<a id="5"></a>
## 5. 🔥 TemplateVM Mantığı — Kök Şablon → Türetilen AppVM'ler

Bu, Qubes'in en zarif ve en çok yanlış anlaşılan mekanizmasıdır. Anladığında her şey yerine oturur.

### Üç parçalı dosya sistemi modeli

Her AppVM'in dosya sistemi **üç katmandan** oluşur:

```
   ┌──────────────────────────────────────────────────────────┐
   │  TemplateVM'den gelen KÖK (/usr, /bin, kurulu yazılım)    │  ← SALT-OKUNUR (çalışırken)
   │  → tüm türetilen AppVM'ler bu KÖK'ü PAYLAŞIR              │
   ├──────────────────────────────────────────────────────────┤
   │  AppVM'in KENDİ /home + /rw  (kişisel veri, ayarlar)      │  ← KALICI (yalnızca bu AppVM'e ait)
   ├──────────────────────────────────────────────────────────┤
   │  Geçici /root değişiklikleri (çalışırken)                 │  ← KAPANINCA SİLİNİR
   └──────────────────────────────────────────────────────────┘
```

**Kritik sonuç — kalıcılık kuralı:**
- TemplateVM'e yazılım kurarsın → **o şablondan türeyen TÜM AppVM'lerde** o yazılım belirir.
- AppVM'in kök sistemine (örn. `/usr`) çalışırken yaptığın değişiklik → **AppVM kapanınca KAYBOLUR.** Kalıcı olması için ya TemplateVM'e kur, ya `/home`'da tut.
- AppVM'in `/home` ve `/rw` dizini → **kalıcıdır**, yalnızca o AppVM'e aittir.

### Neden bu kadar güçlü?

| Avantaj | Açıklama |
|---|---|
| **Tek noktadan güncelleme** | 20 AppVM'in hepsini ayrı ayrı güncellemezsin; **TemplateVM'i bir kez güncellersin**, hepsi yararlanır (yeniden başlattığında) |
| **Saldırı yüzeyi küçülür** | AppVM'in kökü salt-okunur; bir malware `/usr/bin`'i kalıcı kirletemez → AppVM yeniden başlayınca **temiz kök geri gelir** |
| **Disk tasarrufu** | 20 AppVM, tek bir şablonun kökünü paylaşır; her biri ayrı 10 GB OS taşımaz |
| **Sıfırlanabilirlik** | Şüpheli bir AppVM'i yeniden başlatmak, kök seviyesindeki çoğu bulaşmayı temizler |

### Şablon türleri ve disiplin

- **fedora-XX**, **debian-XX** → varsayılan resmi şablonlar.
- **whonix-gateway / whonix-workstation** → Tor için özel şablonlar (Bölüm 10).
- **minimal şablonlar** (fedora-minimal, debian-minimal) → çok küçük, yalnızca ihtiyacın olanı kurarsın → **en küçük saldırı yüzeyi** (Bölüm 14 püf noktası).

> 🔥 **TemplateVM disiplini (piyasada az anlatılan):**
> 1. **TemplateVM'i yalnızca yazılım kurmak/güncellemek için aç**, içinde gündelik iş yapma. Şablon ne kadar "temiz" kalırsa, ondan türeyen her AppVM o kadar güvenli başlar.
> 2. **Farklı güven seviyeleri için farklı şablon kullan.** Örn: "hassas iş" AppVM'leri minimal/sıkı bir şablondan, "gündelik" AppVM'ler genel bir şablondan türesin. Böylece gündelik şablonun şişkinliği hassas işlere bulaşmaz.
> 3. **TemplateVM'in ağı varsayılan kapalı olmalı**, yalnızca güncelleme sırasında (Qubes'in güncelleme proxy'si üzerinden) açılır → şablona internetten doğrudan bulaşma yüzeyi minimize.
> 4. **Şablonu güncelle → türetilen AppVM'leri yeniden başlat.** Güncelleme, AppVM'i *yeniden başlatana kadar* o AppVM'e yansımaz (çalışan AppVM eski kökü RAM'de tutar).

---

<a id="6"></a>
## 6. 🎨 Renk Kodlu Pencere Kenarlıkları — Görsel Güven Ayrımı

Qubes'in en akıllı UX güvenlik özelliklerinden biri: **her qube'a bir renk atanır**, ve o qube'dan gelen *her pencere* o renkli kenarlıkla çizilir. Pencere başlığı da `[qube-adı] Uygulama` formatındadır.

```
  ╔═══════════════════════════════╗   ← KIRMIZI kenarlık (güvenilmez/gündelik)
  ║ [gündelik] Firefox            ║
  ╚═══════════════════════════════╝

  ╔═══════════════════════════════╗   ← YEŞİL kenarlık (güvenilir/banka)
  ║ [banka] Firefox               ║
  ╚═══════════════════════════════╝

  ╔═══════════════════════════════╗   ← MOR/SİYAH (dom0 — yönetim, en hassas)
  ║ Qube Manager                  ║
  ╚═══════════════════════════════╝
```

### Yaygın renk konvansiyonu (zorunlu değil — sen belirlersin)

| Renk | Tipik kullanım |
|---|---|
| 🔴 **Kırmızı** | En güvenilmez — gündelik gezinti, bilinmeyen kaynaklar, malware analizi |
| 🟠 **Turuncu / Sarı** | Orta güven — genel iş, sosyal medya |
| 🟢 **Yeşil** | Güvenilir — bankacılık, hassas iş |
| 🔵 **Mavi** | Kişisel |
| ⚫ **Siyah/Gri** | dom0 (yönetim — en hassas, ona göre ayrı renk) |

> 🔥 **Püf — renk koduna GÜVEN ama aşırı güvenme:** Renkli kenarlık, "şu an hangi qube'dayım, neye yapıştırıyorum" sorusunu görsel olarak yanıtlar → yanlış-qube'a-yapıştırma hatasını **azaltır.** Ama:
> - **Renk salt görseldir;** parolanı yazmadan önce **pencere başlığındaki qube adını da oku** (renk körlüğü, benzer tonlar yanıltır).
> - dom0 penceresine bir AppVM penceresi **kendini benzetemez** — çünkü kenarlığı dom0 çizer, AppVM değil. Bu, sahte-dom0-penceresi (UI spoofing) saldırısına karşı kasıtlı bir savunmadır.
> - Renkleri **anlamlı ve tutarlı** ata; "banka=yeşil" dediysen başka bir yeşili gündelik işe verme.

---

<a id="7"></a>
## 7. 🗑️ DisposableVM — Tek Kullanımlık İzolasyon

DisposableVM (disp-VM), açılıp iş görüldükten sonra **tüm değişiklikleriyle birlikte yok edilen** geçici bir qube'dur. Bir DisposableVM Template'ten (dvm-şablonu) anlık türetilir; kapanınca diski silinir.

### Ne işe yarar

- **Şüpheli ek/link açma:** Bilinmeyen birinden gelen PDF/Office eki → DisposableVM'de aç. Malware patlasa bile qube kapanınca **iz bırakmadan** gider.
- **Tek seferlik gezinti:** Riskli bir siteye girmen gerekiyor → disp-VM'de aç, kapat, bitti.
- **Güvenli inceleme:** Şüpheli bir USB içeriğine bakmak, bir betiği test etmek.

```
   Ek geldi → [DisposableVM aç] → PDF'i içinde aç → oku → [Kapat]
                                                              │
                                                              ▼
                                              Qube + tüm bulaşma YOK OLUR
```

### Kullanım yolları

- **Dosyaya sağ tık → "Open in DisposableVM"** (dosya yöneticisinden).
- **Uygulama menüsünden** disp-VM başlat.
- **Varsayılan disp-VM şablonu** ayarlanır; tüm tek-kullanımlıklar ondan türer.

> 🔥 **Püf — DisposableVM'i varsayılan YAPMA tuzağı:** DisposableVM güçlüdür ama **her şeyi** disp-VM'de yapmak (örn. kalıcı kimlik gerektiren işleri) yanlıştır — her açılışta sıfırlanır, çerez/oturum/ayar kaybolur, *kalıcı kimlik izolasyonu* (Bölüm 14) sağlamaz. Doğru kullanım: **şüpheli/tek-seferlik** işler disp-VM'de; **kalıcı kimlikler** ayrı AppVM'lerde. İkisini karıştırma.
>
> 💡 **Whonix + DisposableVM:** Anonimlik gereken şüpheli linkleri **Whonix tabanlı bir DisposableVM**'de aç → hem Tor üzerinden anonim, hem tek-kullanımlık. CTI/OSINT için altın kombinasyon.

---

<a id="8"></a>
## 8. 📋 Qube'lar Arası Güvenli Transfer & Riskleri

Qube'lar **varsayılan olarak izoledir** — aralarında dosya/pano paylaşımı yoktur. Paylaşmak istediğinde Qubes iki güvenli mekanizma sunar; ama **her transfer bir saldırı yüzeyidir.**

### Dosya kopyalama — qvm-copy

```bash
# AppVM içinden (kaynak qube'da):
qvm-copy dosya.pdf
# → Qubes bir diyalog açar: HANGİ hedef qube'a? (sen seçersin)
# → Hedef qube'un ~/QubesIncoming/<kaynak-qube>/ dizinine düşer
```

- **Yön kontrolü sende:** Hangi qube'a gideceğini **sen onaylarsın** (dom0 diyaloğu). Kötücül bir qube, senin onayın olmadan başka qube'a dosya itemez.
- **qvm-move:** Kopyalamak yerine taşır (kaynaktan siler).

### Pano (Clipboard) — Güvenli Kopyalama Akışı

Qubes'te pano qube'lar arası **otomatik paylaşılmaz.** Özel bir akış vardır:

```
1. Kaynak qube'da normal kopyala (Ctrl+C)
2. Ctrl+Shift+C → içeriği "Qubes global panosuna" yükle
3. Hedef qube'a geç, Ctrl+Shift+V → global panodan o qube'a indir
4. Hedef qube'da normal yapıştır (Ctrl+V)
   → Global pano YENİDEN BOŞALIR (tek seferlik aktarım)
```

Bu üç-adımlı akış kasıtlıdır: pano içeriği **yalnızca senin açık komutunla** bir qube'dan çıkar; arka planda sessiz sızıntı olmaz.

### 🔥 Transfer saldırı yüzeyi (piyasada nadiren konuşulan)

| Risk | Açıklama | Azaltma |
|---|---|---|
| **Yanlış qube'a yapıştırma** | Parolanı banka qube'una değil, gündelik qube'a yapıştırırsan → izolasyon delinir | Yapıştırmadan önce **renk + başlık** kontrol et; parola yöneticisini doğrudan hedef qube'da kullan |
| **Pano içeriğinin hedefte kötüye kullanımı** | Yapıştırdığın metin, hedef qube'da çalışan kötücül bir uygulama tarafından okunabilir | Hassas veriyi yalnızca *güvendiğin* qube'a yapıştır |
| **Dosya formatı exploit'i** | `qvm-copy` ile taşınan bir PDF/dosya, hedef qube'da açılırken o qube'u ele geçirebilir | Şüpheli dosyaları **DisposableVM'de** aç, kalıcı qube'da değil |
| **Pano "yapışkan kalması"** | Global panoyu yükleyip kullanmazsan, bir sonraki Ctrl+Shift+V yanlış yere taşıyabilir | Aktarım sonrası panoyu temiz tut; alışkanlık haline getirme |

> 🧠 **Altın kural:** Her qube-arası transfer, izolasyon duvarında **bilinçli bir kapı** açar. Kapıyı *gerektiğinde* aç, *işin bitince* kapat (yapıştırınca pano boşalır). En tehlikeli transfer, **dosya formatı exploit'i** (kötücül PDF kalıcı qube'a kopyalanıp açılınca) ve **insan hatası** (yanlış qube'a yapıştırma)dır.

---

<a id="9"></a>
## 9. 🔌 sys-usb ile USB İzolasyonu (BadUSB Savunması)

USB, modern bilgisayarın en tehlikeli saldırı yüzeylerinden biridir. Bir USB cihazı, kendini başka bir cihaz gibi tanıtabilir (BadUSB: klavye taklidi yapıp komut yazar) ya da USB denetleyici sürücüsündeki açığı sömürebilir.

### sys-usb ne yapar

`sys-usb`, makinenin **tüm USB denetleyicilerini** üstlenen ayrı bir qube'dur. Bir USB cihazı taktığında, cihaz **doğrudan dom0'a değil, izole `sys-usb`'ye** bağlanır. Kötücül bir USB:
- BadUSB ile komut çalıştırmaya çalışsa → yalnızca `sys-usb`'de çalışır, dom0'a/verine ulaşamaz.
- USB denetleyici açığını sömürse → `sys-usb`'de hapsolur.

```
   USB takıldı
       │
       ▼
  ┌──────────┐   IOMMU ile izole   ┌──────────┐
  │ sys-usb  │ ──── X ──────────►  │   dom0   │   ← USB, dom0'a ASLA doğrudan ulaşamaz
  │ (USB)    │                     │ (krallık)│
  └────┬─────┘                     └──────────┘
       │ qvm-usb attach (SEN onaylarsın)
       ▼
  ┌──────────────┐
  │ hedef AppVM  │  ← cihazı yalnızca açık komutla bir AppVM'e bağlarsın
  └──────────────┘
```

### USB cihazını bir AppVM'e bağlama

```bash
# dom0'dan: bağlı USB cihazlarını listele
qvm-usb list

# Belirli bir USB cihazını hedef AppVM'e bağla
qvm-usb attach <hedef-AppVM> sys-usb:<cihaz-id>

# İş bitince ayır
qvm-usb detach <hedef-AppVM> sys-usb:<cihaz-id>
```

### 🔥 Püf noktaları

- **USB klavye/fare tuzağı:** Eğer USB klavyen/faren varsa, `sys-usb` onları yönetir → BadUSB bir klavye taktığında `sys-usb`'de kalır. **Ama** kurulumda dikkat: tek USB klavyen varsa ve `sys-usb` yapılandırması hatalıysa, girişten kilitlenebilirsin (teyit: Qubes kurulumu USB giriş cihazları için özel onay/politika sorar).
- **USB depolama → DisposableVM'de aç:** Bilinmeyen bir USB belleğin içeriğini **DisposableVM'e** bağlayıp orada incele; kalıcı qube'a değil.
- **qvm-block ile blok cihaz:** Salt depolama için `qvm-block` ile yalnızca blok cihazını (USB'nin tamamını değil) bir qube'a verebilirsin → daha dar yüzey.
- **USB qube'u ağdan ayrı tut:** `sys-usb`'nin ağ erişimi olmamalı (NetVM = none) → ele geçen bir USB qube'u dışarı veri sızdıramaz.

> ⚠️ **Sınır:** `sys-usb` kötücül USB *denetleyici/firmware*'ini izole eder ama **USB üzerinden gelen verinin içeriğini** (örn. exploit'li bir dosya) sen başka qube'da açarsan o qube'u vurabilir. USB izolasyonu donanım katmanını korur; içerik güvenliği yine **DisposableVM disiplinine** bağlıdır.

---

<a id="10"></a>
## 10. 🧅 Whonix Entegrasyonu — Tor ile Ağ Seviyesinde Anonimlik

Qubes, **Whonix**'i (Tor tabanlı anonimlik OS'u) birinci sınıf entegre eder. Whonix iki parçadan oluşur ve Qubes bunları iki qube'a böler:

```
   ┌─────────────────────┐         ┌──────────────────────┐
   │   sys-whonix        │         │   anon-whonix        │
   │   (Whonix Gateway)  │◄────────│   (Whonix Workstation)│
   │   TÜM trafiği Tor'a │  trafik  │   senin çalıştığın    │
   │   ZORLAR            │  buradan │   anonim qube         │
   │   IP'yi gizler      │  geçer   │   (IP'yi GÖREMEZ)     │
   └─────────┬───────────┘         └──────────────────────┘
             │
             ▼  Tor ağı
        İnternet (çıkış IP = Tor exit node)
```

### Mimari neden güçlü

- **anon-whonix (Workstation)**, çalıştığın anonim qube'dur ama **gerçek IP'ni asla göremez.** Çünkü ağı yalnızca `sys-whonix` üzerinden alır.
- **sys-whonix (Gateway)**, *tüm* trafiği Tor'a zorlar. anon-whonix'teki bir uygulama Tor'u baypas etmeye çalışsa bile, ağ fiziksel olarak yalnızca Gateway'den çıkar → **IP sızıntısı (leak) ağ seviyesinde engellenir.**
- Bu, "uygulama düzeyinde Tor"dan (örn. sadece Tor Browser) çok daha güçlüdür: **sistem düzeyinde, sızdıramaz bir Tor zorlaması.**

### Kullanım desenleri

| Amaç | Qube | Not |
|---|---|---|
| Anonim gezinti/araştırma | anon-whonix | Tor Browser dahili |
| Anonim + tek kullanımlık | Whonix DisposableVM | Şüpheli link + anonim + iz bırakmaz |
| Başka qube'u Tor'a sokmak | NetVM = sys-whonix | Herhangi bir AppVM'i sys-whonix'e bağla → o qube anonimleşir |

> 🔥 **Püf noktaları (OSINT/CTI perspektifi):**
> - **Kimlik korelasyonuna dikkat:** Aynı anon-whonix qube'unda hem gerçek kimliğinle ilişkili bir hesaba, hem anonim bir işe girersen → **korelasyon** oluşur. Anonim iş için **ayrı, hiç kişisel iz taşımayan** bir Whonix qube kullan; en güvenlisi her oturum için **Whonix DisposableVM.**
> - **Tor ≠ sihir:** Whonix IP'ni gizler ama **davranışsal parmak izi** (yazım tarzı, oturum saatleri, hesap ilişkileri) seni ele verebilir. Ağ anonimliği, operasyonel anonimliğin yalnızca bir katmanıdır.
> - **Stream isolation:** Whonix, farklı uygulamaları farklı Tor devrelerine ayırır → bir uygulamanın çıkış düğümü diğerini ele vermez (teyit: Whonix varsayılan davranışı).
> - **Clearnet kimliğini Whonix'e karıştırma:** Bankacılık/gerçek-isim işlerini **asla** Whonix qube'unda yapma — Tor exit node'ları izlenebilir ve bazı servisler Tor'u engeller/işaretler.

---

<a id="11"></a>
## 11. 🔑 Split-GPG — Özel Anahtarı Asla Dışarı Çıkarmamak

Qubes'in en zarif ileri-seviye özelliklerinden biri. Sorun şu: GPG özel anahtarın, onu kullanan uygulamayla **aynı qube'da** durursa, o qube ele geçince **anahtarın da gider.** E-posta okuduğun, dosya açtığın qube en çok saldırıya açık olandır — ve GPG anahtarın orada olmamalı.

### Split-GPG mantığı

```
   ┌────────────────────────┐                  ┌──────────────────────────┐
   │   AppVM (e-posta)      │                  │   gpg-qube (OFFLINE)     │
   │   GÜVENİLMEZ           │  "şunu imzala/   │   ÖZEL ANAHTAR BURADA    │
   │   ┌─────────────────┐  │   çöz" isteği    │   ┌────────────────────┐ │
   │   │ Thunderbird     │──┼─────────────────►│   │ GPG özel anahtarı  │ │
   │   │ (anahtara SAHİP │  │                  │   │ (ASLA bu qube'dan  │ │
   │   │  DEĞİL)         │◄─┼──────────────────┤   │  ÇIKMAZ)           │ │
   │   └─────────────────┘  │  imzalı/çözülmüş │   └────────────────────┘ │
   │                        │  SONUÇ döner     │   ⚠️ onay diyaloğu      │
   └────────────────────────┘                  └──────────────────────────┘
```

- Özel anahtar, **ayrı, ağsız (NetVM=none) bir `gpg-qube`**'da durur.
- E-posta qube'u bir şeyi imzalamak/çözmek istediğinde, **ham veriyi gpg-qube'a gönderir**; gpg-qube işlemi yapar ve **yalnızca sonucu** geri verir.
- **Özel anahtar, gpg-qube'un sınırını asla geçmez.** E-posta qube'u ele geçse bile saldırgan anahtarı çalamaz — yalnızca, sen onay verdiğin sürece, imzalama/çözme *isteğinde* bulunabilir.

### Kurulum (kavramsal — komutlar sürümle değişir, teyit et)

```bash
# 1. gpg-qube oluştur (ağsız), özel anahtarı YALNIZCA orada tut
qvm-create gpg-qube --label black
# NetVM = none yap (qube ayarlarından ya da):
qvm-prefs gpg-qube netvm none

# 2. İstemci AppVM'de Split-GPG'yi kur ve gpg-qube'u hedef göster
#    (qubes-gpg-client paketi + QUBES_GPG_DOMAIN değişkeni)
#    Resmi Split-GPG dokümanını adım adım izle.
```

### 🔥 Püf noktaları

- **Onay diyaloğu = son savunma:** Her GPG isteğinde Qubes (genelde) bir onay sorabilir → e-posta qube'u ele geçse bile saldırgan **senin tıklaman olmadan** toplu imza/çözme yapamaz. Onayları otomatikleştirme cazibesine kapılma; manuel onay, kötücül toplu kullanımı durdurur.
- **gpg-qube ağsız olmalı:** `NetVM = none`. Anahtarın bulunduğu qube internete bağlıysa, izolasyonun yarısını kaybedersin.
- **Anahtar üretimini gpg-qube'da yap:** Anahtarı başka yerde üretip taşıma; *doğduğu andan itibaren* o ağsız qube'da kalsın.
- **Yedek:** gpg-qube'u `qvm-backup` ile şifreli yedekle (Bölüm 15) — anahtar tek qube'da olduğu için o qube'un yedeği = anahtarın yedeği.

> 🧠 **Felsefe:** Split-GPG, "anahtarı kullanan kod" ile "anahtarın kendisi"ni **fiziksel olarak ayırır.** Bu, donanım güvenlik anahtarının (YubiKey) yazılımsal taklididir — ve VeraCrypt'in keyfile'ı ayrı USB'de tutma mantığıyla aynı ruhu taşır: **en değerli sırrı, onu kullanan en açık yüzeyden uzak tut.**

---

<a id="12"></a>
## 12. 🔐 Split-SSH — Aynı Mantık, SSH İçin

Split-GPG ile birebir aynı felsefe, bu kez SSH özel anahtarın için. SSH ile sunuculara bağlandığın qube en çok saldırıya açık olandır — özel anahtarın orada durmamalı.

```
   ┌─────────────────────┐               ┌──────────────────────┐
   │  AppVM (geliştirme) │  "şu sunucuya │  ssh-qube (OFFLINE)  │
   │  GÜVENİLMEZ         │   bağlanmak    │  SSH ÖZEL ANAHTARI   │
   │  ssh client ────────┼───────────────►│  (ssh-agent burada)  │
   │  (anahtar YOK)      │  için imza iste│  anahtar ASLA çıkmaz │
   └─────────────────────┘               └──────────────────────┘
```

- SSH özel anahtarın **ayrı, ağsız `ssh-qube`**'da bir `ssh-agent` içinde durur.
- Geliştirme qube'u SSH bağlantısı kurarken, kimlik doğrulama imzasını `ssh-qube`'dan ister (`SSH_AUTH_SOCK` qube-arası proxy üzerinden).
- **Özel anahtar ssh-qube'dan çıkmaz**; geliştirme qube'u ele geçse bile saldırgan anahtarı **çalamaz** (yalnızca, sen aktifken, imza isteğinde bulunabilir).

> 🔥 **Püf:** Split-SSH'i **onay/zaman-aşımı** ile birleştir → her bağlantıda (ya da belirli aralıkla) onay iste. Böylece ele geçen geliştirme qube'u, sen makinedeyken bile **sınırsız** sunucu erişimi yapamaz. Komut detayları için resmi "Split SSH" Qubes topluluk dokümanını izle (sürümle değişir — **teyit et**).

> 💡 **Genel desen:** Split-GPG ve Split-SSH, **"hassas anahtar izolasyonu"** desenin iki örneğidir. Aynı mantığı parola yöneticisi, kripto cüzdanı özel anahtarı gibi her kritik sır için kurabilirsin: **sır → ağsız özel qube; kullanan uygulama → ayrı qube; arada onaylı proxy.**

---

<a id="13"></a>
## 13. 🎥 USB Qube, Kamera & Mikrofon Kontrolü

Qubes, donanım çevre birimlerine erişimi de qube'lara böler — bir qube, sen *açıkça vermedikçe* kameraya/mikrofona/USB'ye erişemez.

### Varsayılan: erişim YOK

- Bir AppVM, **varsayılan olarak** kameraya, mikrofona ya da USB cihazlarına erişemez. Bunlar `sys-usb` ya da ilgili denetleyicide tutulur.
- Bir cihazı bir qube'a vermek için **açık `qvm-usb attach` / aygıt atama** gerekir → kazara/kötücül erişim yok.

### Kamera/mikrofon disiplini

| Tehdit | Qubes savunması |
|---|---|
| Kötücül uygulama gizlice kamerayı açar | Uygulama qube'unda kamera **bağlı değilse**, fiziksel olarak erişemez |
| Mikrofon dinlemesi | Mikrofonu yalnızca **ihtiyaç anında** ilgili qube'a bağla, iş bitince ayır |
| Görüntülü görüşme | Kamera/mikrofonu yalnızca o görüşme qube'una, yalnızca o süre için ver |

```bash
# Mikrofonu bir qube'a ver (örnek — cihaz adı sistemine göre değişir, teyit et)
qvm-device mic attach <hedef-qube> dom0:mic

# İş bitince ayır
qvm-device mic detach <hedef-qube> dom0:mic
```

> 🔥 **Püf noktaları:**
> - **"İhtiyaç anında bağla, hemen ayır"** disiplini → bir qube ele geçse bile, kamera/mikrofon ona bağlı değilse dinleyemez/izleyemez.
> - **Kamera LED'i baypası Qubes'in kapsamı dışıdır:** Qubes *erişimi* engeller (cihaz bağlı değilse açılamaz), ama bir qube'a kamerayı *verdiysen*, o qube içindeki yazılım LED davranışını yönetir — bu donanım/sürücü katmanıdır. En güvenli kamera, **bağlı olmayan kameradır** (ya da fiziksel kapatıcı).
> - **Hassas görüşmeler için ayrı qube:** Görüntülü görüşmeyi gündelik/iş qube'unda değil, **özel bir görüşme qube'unda** yap; bitince cihazları ayır, qube'u kapat.

---

<a id="14"></a>
## 14. 🔥 PÜF NOKTALARI — Piyasada Bulamayacakların

Bu bölüm, çoğu rehberin atladığı ve gerçek dünyada izolasyonu **çökerten** ya da **ustalaştıran** detaylardır.

### 14.1 dom0'ı ASLA Kirletme (Birinci Emir)

- dom0, tüm qube'ları yönetir → dom0 ele geçerse **her şey biter.** İzolasyonun *anlamı* dom0'ın temizliğine bağlıdır.
- **Yapma:** dom0'da internet açma, dosya indirme/açma, yabancı yazılım/depo ekleme, rastgele script çalıştırma, USB içeriğine dom0'dan bakma.
- **Yap:** dom0'ı yalnızca `qubes-os.org` güncelleme kanalından güncelle; minimal tut; dom0'ı yalnızca pencere yönetimi ve `qvm-*` araçları için kullan.

### 14.2 Kimlik-Bazlı Qube Ayrımı (Mimarinin Kalbi)

Qubes'in asıl gücü, qube'ları **kimliklerine göre** ayırmandır — uygulamaya göre değil, *kim olduğuna* göre:

```
   ┌──────────┐  ┌───────────┐  ┌──────────┐  ┌───────────┐  ┌──────────┐
   │   iş     │  │  kişisel  │  │  banka   │  │  anonim   │  │ gündelik │
   │ (sarı)   │  │  (mavi)   │  │ (yeşil)  │  │ (Whonix)  │  │ (kırmızı)│
   │ iş maili │  │ aile      │  │ sadece   │  │ OSINT/    │  │ rastgele │
   │ iş dosya │  │ foto      │  │ bankacılık│ │ araştırma │  │ gezinti  │
   └──────────┘  └───────────┘  └──────────┘  └───────────┘  └──────────┘
        ▲              ▲              ▲             ▲              ▲
        └──────────────┴──── Aralarında VARSAYILAN bağlantı YOK ───┴──────┘
```

- **Banka qube'u:** *yalnızca* bankacılık. Başka site açma, başka dosya koyma → en küçük saldırı yüzeyi, en yüksek güven.
- **İş ≠ kişisel ≠ anonim:** Biri sızsa diğer kimliklerin açığa çıkmaz; aralarında korelasyon kurulamaz.
- **CTI/OSINT için:** Araştırma kimliğini (anonim), iş kimliğini (kurumsal) ve kişisel kimliğini **kesinlikle ayrı** tut → bir hedef seni "geri izleyemez."

### 14.3 DisposableVM'i Varsayılan YAPMA Tuzağı

- DisposableVM şüpheli/tek-seferlik işler içindir. **Kalıcı kimlik** gereken işleri (çerez, oturum, ayar) disp-VM'de yapma — her açılışta sıfırlanır.
- **Doğru:** şüpheli ek/link → disp-VM; kalıcı iş → AppVM. (Bölüm 7.)

### 14.4 TemplateVM'i Minimal Tut

- Şablon ne kadar şişerse (gereksiz paketler), ondan türeyen *her* AppVM o kadar geniş saldırı yüzeyiyle başlar.
- **fedora-minimal / debian-minimal** şablonlardan, *yalnızca ihtiyacın olanı* kurarak özel şablonlar üret → minimum yüzey. (Teyit: minimal şablonlar elle paket kurulumu ister.)
- Hassas işler için **ayrı, sıkı bir şablon**; gündelik işler için ayrı bir şablon → şişkinlik hassasa bulaşmaz.

### 14.5 sys-net'i Ayrı ve Disposable Tut

- `sys-net`, internete doğrudan dokunan en açık qube'dur → ağ sürücüsü açıkları burada patlar.
- `sys-net`'i **DisposableVM tabanlı** yapabilirsin (her başlatmada temiz) → kalıcı bulaşma tutamaz (teyit: Qubes 4.x bunu destekler).
- `sys-net`'in *yalnızca* ağ donanımını taşıdığından emin ol; başka iş yükleme.

### 14.6 Pano & Dosya Transferi = Saldırı Yüzeyi

- Her qube-arası transfer izolasyon duvarında kapı açar (Bölüm 8). En sık hata: **yanlış qube'a yapıştırma.**
- **Azaltma:** Parola yöneticisini **doğrudan hedef qube'da** kullan (panoyla taşıma); transferden önce **renk + başlık** doğrula; şüpheli dosyaları disp-VM'de aç.

### 14.7 Renk Kodlarına "Tek Başına" Güvenme

- Renk yardımcıdır ama **kesin değildir** (renk körlüğü, benzer tonlar). Kritik işlemden (parola, banka) önce **pencere başlığındaki qube adını oku.** dom0 penceresi taklit edilemez (kenarlığı dom0 çizer) — ama AppVM'ler arası benzer renkler insanı yanıltabilir.

### 14.8 Bellek & Disk Gereksinimleri Gerçeği

- Qubes RAM açlığı çeker: her açık qube RAM yer. **6 GB ile zar zor, 16 GB rahat, 32 GB konforlu.** Az RAM'de qube'lar açılmaz/donar → insanlar "her şeyi tek qube'da toplama" hatasına kayar (izolasyonu kaybeder).
- Disk: şablonlar + qube'lar hızla şişer. **256 GB+ SSD** hedefle; HDD'de Qubes acı verir.
- **Gerçekçi ol:** Yetersiz donanımda Qubes'i zorlamak, izolasyondan ödün vermeye iter — doğru donanımla başla.

### 14.9 Güncelleme Disiplini

- **dom0'ı, TemplateVM'leri ve Whonix'i ayrı ayrı güncelle.** TemplateVM güncellemesi, türetilen AppVM'lere ancak **yeniden başlatınca** yansır.
- Güncellemeleri Qubes'in **güncelleme proxy'si** üzerinden yap (şablonlar doğrudan internete çıkmaz) → tedarik zinciri yüzeyi daralır.

### 14.10 Donanım/Firmware Katmanı (Qubes'in Altı)

- Qubes, **Intel ME / BIOS / firmware** katmanını koruyamaz (Bölüm 2). İleri seviye azaltma: **Coreboot/Heads** firmware + **Intel ME devre dışı/neutralize** edilmiş donanım (örn. belirli ThinkPad'ler) → Evil Maid ve firmware implantına karşı kısmi savunma. (Teyit: donanım-spesifik, dikkatli araştır.)
- **Anti-Evil-Maid (AEM):** Qubes, önyükleme bütünlüğünü doğrulayan AEM'i destekler (TPM ile) → bootloader'ın kurcalanıp kurcalanmadığını anlarsın. Hassas tehdit modelinde değerlidir (teyit: kurulumu donanıma bağlı).

### 14.11 Yedekleme Şart (qvm-backup şifreli)

- Qube'lar diskte; disk ölürse her şey gider. **`qvm-backup` ile düzenli, şifreli yedek al** (Bölüm 15). Yedeği farklı, güvenli ortamda tut; **geri yükleyerek test et.**

> 🧠 **Püf noktalarının özü:** Qubes'in güvenliği üç ayağa dayanır — **(1) dom0 temizliği, (2) doğru qube ayrımı (kimlik-bazlı + minimal şablon), (3) transfer/yapıştırma disiplini.** Üçünden biri çökerse, en pahalı donanım bile seni kurtarmaz.

---

<a id="15"></a>
## 15. 💾 Yedekleme & Kurtarma (qvm-backup Şifreli)

Qubes'te qube'lar diskte yaşar — disk ölürse, dom0 bozulursa ya da kurulum gidersé her şey gider. Qubes yerleşik, **şifreli** bir yedekleme aracı sunar.

### qvm-backup mantığı

```bash
# GUI: dom0 → "Qubes Backup" aracı → yedeklenecek qube'ları seç →
#      güçlü parola → hedef qube/dizin (örn. harici şifreli disk) → başlat

# CLI (dom0):
qvm-backup --dest-vm <hedef-qube> /yedek/yolu  qube1 qube2 gpg-qube ...
```

- **Yedekler şifrelidir** — yedek parolası, geri yüklemenin tek anahtarıdır. **Unutursan yedek ölür.**
- Yedek bir **hedef qube/dizine** yazılır → oradan harici şifreli diske/buluta taşı.
- **Geri yükleme:** "Qubes Restore" aracı (ya da `qvm-backup-restore`) ile seçili qube'ları geri al.

### 🔥 Püf noktaları

- **gpg-qube / ssh-qube'u mutlaka yedekle:** En değerli sırların (özel anahtarlar) o qube'larda. Onların yedeği = anahtarlarının yedeği. (Bölüm 11-12.)
- **Yedek parolası ≠ qube parolaları:** Ayrı, güçlü bir parola seç; güvenli sakla (ezber + ayrı yedek).
- **3-2-1 kuralı:** 3 kopya, 2 farklı ortam, 1 saha dışı. Tek yedek = yedek değildir.
- **Geri yüklemeyi TEST ET:** Yedeği farklı/temiz bir Qubes kurulumuna geri yükleyerek doğrula. Açılmayan yedek, yedek değildir.
- **Şifreli harici diske yaz:** Yedeği VeraCrypt/LUKS şifreli bir diske koy → **katmanlı koruma** (Qubes şifreler + disk şifreler).

> 💡 **VeraCrypt köprüsü:** Qubes yedeğini bir **VeraCrypt hidden volume**'a yaz → hem şifreli, hem inkâr edilebilir. Zorlanırsan dış birimi verirsin, Qubes yedeğini barındıran gizli birim görünmez kalır. (Bkz. `VERACRYPT_USTALIK_REHBERI.md`.)

---

<a id="16"></a>
## 16. ☠️ Yaygın Ölümcül Hatalar

1. **dom0'da iş yapmak** → internet açma, dosya indirme/açma, yazılım kurma. **En ölümcül hata** — dom0 ele geçerse tüm bölmeleme anlamsızlaşır.
2. **Yanlış qube'a yapıştırma** → parolayı/dosyayı yanlış güven seviyesindeki qube'a taşımak; izolasyonu *kendi elinle* delmek.
3. **Her şeyi tek qube'da toplamak** → "tek büyük AppVM"de bankacılık + gündelik + iş → sıradan bir Linux'tan farksız (sadece daha yavaş); Qubes'in tüm anlamını kaybetmek.
4. **TemplateVM'de gündelik iş yapmak** → şablonu kirletmek; ondan türeyen *her* AppVM'i baştan riskli başlatmak.
5. **VT-d/IOMMU'suz donanımda Qubes kurmak** → izolasyonun büyük kısmı kağıt üzerinde kalır; `sys-net` DMA ile dom0'ı vurabilir.
6. **DisposableVM'i kalıcı kimlik için kullanmak** → her açılışta sıfırlanır; oturum/çerez/kimlik kaybı + yanlış güvenlik hissi.
7. **gpg-qube/ssh-qube'a internet vermek** → anahtar izolasyonunun yarısını kaybetmek (Split-GPG/SSH ağsız qube ister).
8. **Yedek almamak / geri yüklemeyi test etmemek** → disk ölünce her qube ve her sır gider.
9. **Whonix qube'una gerçek kimlik karıştırmak** → anonimlik ile clearnet kimliğini aynı qube'da buluşturup korelasyon yaratmak.
10. **Renk koduna kör güven** → yalnızca renge bakıp başlığı okumamak; benzer tonlarda yanlış qube'a hassas veri vermek.
11. **İmza doğrulamadan ISO kurmak** → değiştirilmiş bir Qubes, *kurulumdan itibaren* arka kapılıdır.
12. **Şablonu güncelleyip AppVM'leri yeniden başlatmamak** → güncellemenin (güvenlik yamasının) AppVM'lere yansımadığını sanmak.

---

<a id="17"></a>
## 17. 🏰 Kanije Kalesi ile Birlikte Kullanım

Bu repo (Kanije Kalesi), **fiziksel tehdit anında** cihazı uzaktan/otomatik koruyan bir muhafızdır. Qubes **çalışan sistemi bölmelere ayırır**; Kanije **olay anını** yönetir. Üçü (Qubes + VeraCrypt + Kanije) katmanlı bir savunma oluşturur:

| Senaryo | Qubes rolü | VeraCrypt rolü | Kanije Kalesi rolü |
|---|---|---|---|
| Bir uygulama ele geçirildi | Bulaşma **tek qube'da** kalır | — | — |
| Şüpheli dosya/malware analizi | **DisposableVM**'de patlat, iz bırakmadan yok et | — | — |
| Cihaz açık, biri yaklaştı | Açık qube'lar risk altında | Mount'lu birim açık (risk) | `/koruma` dead-man / USB / yanlış-giriş → **kilitle + alarm + foto** |
| Acil durum | — | — | `/panik`, `/kilit tam` (lockdown) |
| Hassas iş izolasyonu | Kimlik-bazlı qube'da yalıtık çalış | VeraCrypt hidden volume'da veri | — |
| Verinin yok edilmesi | Qube diskini sil | Birim header'ını sil (crypto-shred) | `/imha` ile hedef dosyaları güvenli sil |
| Adli inceleme öncesi | DisposableVM'ler iz bırakmaz | Hidden volume inkâr edilebilirlik | RAM-only mod, iz temizleme |

### 🔥 Önerilen entegrasyon deseni

1. **Hassas işi izole qube'larda yap:** Kimlik-bazlı ayrım (iş/kişisel/anon/banka). En kritik veriyi, bir AppVM içinde **VeraCrypt hidden volume**'da tut → Qubes izolasyonu + VeraCrypt şifreleme **iç içe iki kale.**
2. **Malware analizini DisposableVM'de yap:** Şüpheli örneği ağsız (NetVM=none) ya da Whonix-disp bir qube'da patlat → ana sisteme hiç dokunmaz, kapanınca yok olur. CTI analisti için ideal kum havuzu.
3. **Fiziksel tehdit duruşu:** Kanije `/koruma`'yı (dead-man + USB dead-man) aç → sen X süre check-in yapmazsan ya da USB çıkarılırsa Kanije **kilitler + alarm + foto** çeker. Qubes'in açık qube'larındaki "kullanımdaki veri", fiziksel tehditte böyle korunur.
4. **VeraCrypt + Qubes katmanı:** Disk seviyesinde tüm Qubes kurulumu zaten şifrelidir (LUKS); kritik sırrı ek olarak bir qube içindeki VeraCrypt hidden volume'a koy → **çift şifreleme + inkâr edilebilirlik.**
5. **Crypto-shredding ile hızlı yok etme:** En kritik senaryoda, VeraCrypt birim **header'ını** silmek (64 KB) tüm veriyi anında erişilemez kılar — Kanije'nin `/imha` hedeflerine bu konumu eklemek, "veriyi sil"i saniyeler içinde yapar (gigabaytları silmeye gerek yok; **anahtarı yok et, veri çöp olur**).

> 🧠 **Felsefe örtüşmesi:** Üç araç da aynı ilkeyi paylaşır — **"en kötüyü varsay, hasarı sınırla."** Qubes hasarı tek qube'a hapseder; VeraCrypt veriyi matematiksel kaleye kilitler; Kanije sen unuttuğunda/tehdit anında kapıyı kapatır. Tek başına hiçbiri yetmez; **birlikte derinlemesine savunma (defense-in-depth)** kurarlar.

---

<a id="18"></a>
## 18. ✅ Hızlı Referans (qvm-* Komutları) & Operasyonel Kontrol Listesi

### Sık qvm-* komutları (dom0'da çalışır — teyit et, sürümle değişebilir)

```bash
# --- Qube yaşam döngüsü ---
qvm-ls                       # Tüm qube'ları ve durumlarını listele
qvm-start <qube>             # Qube başlat
qvm-shutdown <qube>          # Qube kapat
qvm-kill <qube>              # Yanıt vermeyen qube'u zorla kapat
qvm-create <ad> --label red  # Yeni qube oluştur (renk ata)
qvm-remove <qube>            # Qube sil

# --- Yapılandırma ---
qvm-prefs <qube>             # Qube ayarlarını göster
qvm-prefs <qube> netvm sys-whonix   # Qube'u Whonix'e (Tor'a) bağla
qvm-prefs <qube> netvm none         # Qube'u tamamen çevrimdışı yap
qvm-prefs <qube> template <şablon>  # Qube'un şablonunu değiştir
qvm-prefs <qube> memory 4000        # RAM ayarı

# --- Transfer ---
qvm-copy <dosya>             # Dosyayı başka qube'a kopyala (hedefi onaylarsın)
qvm-move <dosya>             # Taşı (kaynaktan sil)
# Pano: Ctrl+Shift+C (yükle) → hedefte Ctrl+Shift+V (indir)

# --- USB / aygıt ---
qvm-usb list                 # Bağlı USB cihazları
qvm-usb attach <qube> sys-usb:<id>   # USB'yi qube'a bağla
qvm-usb detach <qube> sys-usb:<id>   # Ayır
qvm-block / qvm-device       # Blok cihaz / genel aygıt yönetimi

# --- Şablon / güncelleme ---
qubes-update-gui             # GUI güncelleme aracı (dom0, şablonlar, Whonix)
qubes-dom0-update            # dom0'ı güncelle (yalnızca resmi kanal)

# --- Yedek ---
qvm-backup                   # Şifreli yedek (GUI: "Qubes Backup")
qvm-backup-restore           # Geri yükle
```

### Kurulum / ilk yapılandırma kontrol listesi
- [ ] Donanım **HCL**'de doğrulandı (Wi-Fi/GPU/suspend çalışıyor)
- [ ] **VT-x VE VT-d (IOMMU)** BIOS'ta etkin
- [ ] ISO **imzası doğrulandı** (Master Signing Key parmak izi teyit edildi)
- [ ] RAM **16 GB+**, disk **256 GB+ SSD**
- [ ] **sys-usb** kuruldu (USB izolasyonu aktif)
- [ ] **sys-whonix + anon-whonix** kuruldu (anonim ihtiyaç varsa)

### Mimari / qube ayrımı
- [ ] **Kimlik-bazlı** qube'lar (iş / kişisel / banka / anon / gündelik) ayrı
- [ ] **Banka qube'u** yalnızca bankacılık (en dar yüzey)
- [ ] Anonim iş **Whonix** üzerinden, gerçek kimlik karışmıyor
- [ ] Hassas anahtarlar **Split-GPG / Split-SSH** ile ağsız qube'da
- [ ] **TemplateVM minimal**, gündelik iş şablonda yapılmıyor
- [ ] **sys-net** ayrı (ideal: disposable)

### Her oturum / operasyonel
- [ ] **dom0'da hiçbir şey açma/indirme** (internet yok, dosya yok)
- [ ] Yapıştırmadan önce **renk + pencere başlığı** kontrol
- [ ] Şüpheli ek/link **DisposableVM**'de
- [ ] USB/kamera/mikrofon yalnızca **ihtiyaç anında bağla, hemen ayır**
- [ ] İş bitince şüpheli qube'ları **kapat** (disp-VM yok olur)

### Periyodik
- [ ] **dom0 + TemplateVM + Whonix** güncel (şablon sonrası AppVM yeniden başlat)
- [ ] **qvm-backup** şifreli yedek alındı, **farklı makinede test** edildi
- [ ] gpg-qube / ssh-qube yedeği güvende
- [ ] Kullanılmayan qube'lar temizlendi (saldırı yüzeyi azalt)

---

<a id="19"></a>
## 19. ⚖️ Hukuki Sınır & Operasyonel Notlar

- **Meşru amaç:** Bu rehber **gizlilik, gazetecilik, güvenlik araştırması, CTI/malware analizi ve kişisel veri koruma** içindir. Qubes, hassas işle uğraşan herkesin (aktivist, araştırmacı, gazeteci, analist) meşru aracıdır.
- **Anonimlik ≠ dokunulmazlık:** Whonix/Tor ağ izini gizler ama **operasyonel hatalar** (kimlik karıştırma, davranışsal parmak izi, yanlış qube) seni ele verebilir. Teknoloji, disiplinin yerini tutmaz.
- **Sınır geçişi:** Tüm disk şifreli olsa bile, bazı yargı bölgelerinde parola talep edilebilir/reddetmek suç olabilir. Qubes'in LUKS şifrelemesi + (gerekirse) içteki VeraCrypt hidden volume katmanı, inkâr edilebilirlik için VeraCrypt rehberindeki notlarla birlikte değerlendir.
- **Malware analizi sorumluluğu:** Şüpheli örnekleri **yalnızca izole/ağsız DisposableVM**'de patlat; kötücül trafiği kasıtsız yaymamak için ağ qube'unu kontrollü tut. Analiz ettiğin örneklerin yasal statüsünü ve kurumsal politikanı gözet.
- **Yasal uyum:** Bulunduğun yargı bölgesinin yasalarına uy. Bu rehber savunma ve gizlilik içindir, kötüye kullanım için değil.

---

> 🏰 **Kapanış:** Qubes OS bir ürün değil, bir **disiplindir.** En ince Xen hipervizörü bile, dom0'da açtığın bir dosyayla ya da yanlış qube'a yapıştırdığın bir parolayla çaresizdir. Qubes sana **bölmeleme yoluyla matematiksel-yapısal bir kale** verir — onlarca izole burç, tek surun içinde. Ama **hangi kapının açık, hangi burçta kim olduğunu bilmek senin işin.** En zayıf halka her zaman donanımın değil, **disiplinindir.**
>
> İzolasyon, ele geçirilmeyi *önlemez* — ele geçirildiğinde **hasarı tek burca hapseder.** Bunu içselleştirdiğinde Qubes'i ustaca kullanıyorsun demektir.
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. İlgili: `VERACRYPT_USTALIK_REHBERI.md`, `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`, `SIFRE_KRONOLOJISI_VE_USB_SIFRELEME.md`, `DUAL_BOOT_VE_DEPOLAMA_GUVENLIGI.md`.*
