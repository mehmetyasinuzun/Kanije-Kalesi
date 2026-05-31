# 🦈 WIRESHARK & AĞ TRAFİK ANALİZİ — TAM CTI USTALIK REHBERİ
## Paket Yakalamadan C2 Avına, JA3'ten Beacon Tespitine, Püf Noktalarıyla Uçtan Uca

> **Amaç:** Wireshark'ı "aç, başlat, kaydır" seviyesinden çıkarıp, bir **siber tehdit istihbaratı (CTI)** analisti gibi ağ trafiğinde **C2 kanalı, veri sızması (exfiltration), DNS tüneli ve lateral movement** avlayacak şekilde kullanmayı öğretmek. Bu rehber yalnızca *hangi filtre*'yi değil, **neden o filtre**, **paket nerede yakalanır** ve **hangi durumda hiçbir şey göremezsin**'i de anlatır. Forum cevaplarında bulamayacağın capture-vs-display filtre farkı, yakalama noktası körlüğü, beacon matematiği, JA3 parmak izi ve TLS'in seni nasıl kör ettiği — ama metadata'nın hâlâ konuştuğu — burada.

> ⚠️ **Önce bunu oku — HUKUK:** Başkasının ağ trafiğini **yetkisiz dinlemek/yakalamak çoğu yargı bölgesinde SUÇTUR** (Türkiye TCK 243/244 bilişim suçları, ABD Wiretap Act, AB ePrivacy). Bu rehber **yalnızca kendi ağında, yetkili bir laboratuvar/SOC ortamında, ya da yazılı izinli pentest kapsamında** kullanılmak içindir. "Ortak Wi-Fi'de komşumun trafiğine bakayım" **suçtur** — teknik olarak mümkün olması yasal olduğu anlamına gelmez. Bölüm 13'ü atlama.

> ⚠️ **İkinci uyarı — yanlış güvenlik hissi:** Bir pcap'i açıp "bağlantı görüyorum, demek ki güvendeyim" demek, *hiç bakmamaktan* daha tehlikelidir. Switch'li bir ağda **trafiğin çoğunu zaten göremezsin** (Bölüm 3); TLS trafiğin **içeriğini okuyamazsın** (Bölüm 6 & 8). Neyi göremediğini bilmeden yapılan analiz, sahte bir "temiz" raporu üretir. **Yakalama noktası körlüğü** ve **TLS körlüğü** bölümlerini atlama.

---

## 📑 İÇİNDEKİLER

1. [Wireshark Nedir, CTI'de Ağ Analizinin Rolü](#1)
2. [Yakalama Temelleri — Arayüz, Promiscuous, dumpcap, Ring Buffer](#2)
3. [🔥 Capture Filter (BPF) vs Display Filter — En Sık Karıştırılan Konu](#3)
4. [🔥 Yakalama NOKTASI — TAP vs SPAN vs Hub & Switch Körlüğü](#4)
5. [Display Filter Ustalığı — Söz Dizimi & Hazır Kütüphane](#5)
6. [Follow Stream, Export Objects, Statistics, IO Graphs](#6)
7. [Protokol-Protokol Analiz — DNS / HTTP / TLS / ICMP / SMB](#7)
8. [🔥 Malware Trafiği Avı — Beaconing, C2, DNS Tüneli, Exfil](#8)
9. [TLS Şifre Çözme — SSLKEYLOGFILE, RSA Anahtar, Forward Secrecy Sınırı](#9)
10. [tshark CLI — Otomasyon, Toplu IOC Çıkarma, Dilimleme](#10)
11. [🔥 PÜF NOKTALARI — Piyasada Bulamayacakların](#11)
12. [Yaygın Ölümcül Hatalar](#12)
13. [🏰 Kanije Kalesi ile Birlikte Kullanım & Öz-Farkındalık](#13)
14. [Hızlı Referans — En Yararlı Filtreler & Analiz Kontrol Listesi](#14)
15. [⚖️ Hukuki Sınır — Yetkisiz Dinleme SUÇTUR](#15)

---

<a id="1"></a>
## 1. 🧭 Wireshark Nedir, CTI'de Ağ Analizinin Rolü

**Wireshark**, dünyanın en yaygın açık kaynaklı **ağ protokolü analizörüdür** (eski adı *Ethereal*, 2006'da Wireshark oldu; Gerald Combs ve geniş bir topluluk sürdürür). Ağ kartından geçen ham paketleri yakalar, **2000'den fazla protokolü** ayrıştırır (dissect eder) ve sana her paketi katman katman (Ethernet → IP → TCP → uygulama) gösterir.

Wireshark **pasif** bir araçtır: trafiği *üretmez*, yalnızca *gözlemler*. Bu onu CTI ve adli ağ analizi (Network Forensics) için ideal kılar — "telde gerçekte ne aktığını" yalan söyleyemeyecek tek kaynak budur. Bir uç noktadaki kötücül yazılım kendini gizleyebilir (rootkit), ama ağda iletişim kurması gerekiyorsa **o paketler telden geçmek zorundadır**, ve orada onu yakalarsın.

### CTI'de ağ analizinin beş temel rolü

| Rol | Ne aranır | Tipik imza |
|---|---|---|
| **C2 (Command & Control) tespiti** | Kötücül yazılımın "patronuyla" konuşması | Düzenli aralıklı (beacon) giden bağlantılar, bilinen kötü IP/domain, anormal JA3 |
| **Veri sızması (Exfiltration)** | Verinin dışarı kaçırılması | Anormal büyük giden (upload) hacim, DNS/ICMP içine gömülü veri, "low and slow" sızıntı |
| **IOC çıkarma** | İhlal Göstergesi (Indicator of Compromise) toplama | IP, domain, URL, JA3/JA3S hash, sertifika seri no, User-Agent, dosya hash'i |
| **Lateral Movement (yanal hareket)** | Saldırganın ağ içinde yayılması | SMB/RPC, PsExec, WMI, RDP, Kerberos anomalileri, iç→iç tarama |
| **Malware trafiği sınıflandırma** | Ailenin/kampanyanın belirlenmesi | Beacon deseni, URI yapısı, TLS parmak izi, payload kalıbı |

> 🧠 **CTI ilkesi — "ağ yalan söylemez":** Bir saldırgan diskteki dosyaları silebilir, log'ları temizleyebilir, EDR'yi atlatabilir. Ama **dışarıyla konuşmak zorundaysa**, o trafik bir yakalama noktasından geçer. Ağ analizi, uç-nokta telemetrisinin (EDR) **kör noktalarını kapatan** bağımsız bir kanıt katmanıdır. İkisi birlikte kullanılır; biri diğerinin yerine geçmez.

### Wireshark vs alternatifler (dürüst kıyas)

| Araç | Güçlü yanı | Ne zaman |
|---|---|---|
| **Wireshark (GUI)** | Derin, görsel, etkileşimli analiz; "neden" sorusu | Tek bir olayı/oturumu derinlemesine kazmak |
| **tshark (CLI)** | Wireshark'ın komut satırı kardeşi; scriptlenebilir | Otomasyon, toplu IOC çıkarma, dev pcap'leri dilimleme |
| **dumpcap** | Wireshark'ın *yakalama motoru* (en hafif, en az kayıp) | Uzun süreli/yüksek hacimli yakalama (ring buffer) |
| **tcpdump** | Sunucularda her yerde var, hafif | Uzaktan/headless yakalama → sonra Wireshark'la aç |
| **Zeek (eski Bro)** | Paketi *bağlam/log*'a çevirir (akış-bazlı) | Sürekli SOC izleme, devasa hacim, davranışsal log |
| **Suricata** | İmza tabanlı IDS/IPS + akış log | Bilinen tehdit imzalarını gerçek-zamanlı yakalamak |

> 💡 **Katmanlı yaklaşım:** Pratikte SOC akışı şudur: **Suricata/Zeek** sürekli izler ve alarm/log üretir → bir alarm tetiklendiğinde ilgili **pcap** çıkarılır → **Wireshark/tshark** ile insan analist derinlemesine kazar. Wireshark "her şeyi 7/24 dinleyen" araç değil; **"şüpheyi kanıta çeviren mikroskop"**tur.

---

<a id="2"></a>
## 2. 🎛️ Yakalama Temelleri — Arayüz, Promiscuous, dumpcap, Ring Buffer

### Arayüz seçimi
Wireshark'ı açınca **yakalanacak arayüzü** seçersin. Her arayüzün yanındaki **sparkline** (mini grafik) o anda hangisinde trafik aktığını gösterir — doğru arayüzü seçmenin en hızlı yolu budur.

- **Wi-Fi (`wlan0` / `Wi-Fi`):** Kablosuz; promiscuous/monitor modunda özel davranır (aşağıda).
- **Ethernet (`eth0` / `Ethernet`):** Kablolu; en temiz yakalama genelde burada.
- **Loopback (`lo` / `Npcap Loopback Adapter`):** `127.0.0.1` üzerindeki yerel-içi trafik (bir uygulamanın localhost'taki başka servise konuşması).
- **`any` (yalnızca Linux):** Tüm arayüzleri aynı anda yakalar (sözde-arayüz). Pratik ama bazı L2 detayını kaybeder.
- **USBPcap / Bluetooth / sanal adaptör:** İlgili donanım trafiği.

> ⚠️ **Windows notu:** Yakalama için **Npcap** sürücüsü gerekir (Wireshark kurulumuyla gelir; "WinPcap uyumluluk modu" ve "loopback desteği" kutularını işaretle). Linux/maclog'da `libpcap` kullanılır ve yakalama için ya **root** ya da `wireshark` grubu + `dumpcap` setcap yetkisi gerekir.

### Promiscuous mode (karışık mod) — ve büyük yanılgı
**Promiscuous mode**, ağ kartına "yalnızca bana adreslenmiş çerçeveleri değil, **gördüğün TÜM çerçeveleri** işletim sistemine ver" der. Wireshark bunu varsayılan olarak açar.

> 🔥 **EN SIK YANILGI:** "Promiscuous mode açtım, artık ağdaki herkesin trafiğini görürüm." **YANLIŞ.** Promiscuous mode yalnızca **kartına fiziksel olarak ULAŞAN** çerçeveleri görmeni sağlar. **Modern switch'li bir ağda, switch trafiği yalnızca hedef porta yollar** — başkasının trafiği senin portuna hiç gelmez, dolayısıyla promiscuous mode açık olsa da **görecek bir şeyin yoktur.** (Eski **hub**'larda her şey her porta giderdi; switch bunu bitirdi.) Başkasının trafiğini görmek istiyorsan **yakalama NOKTASINI** değiştirmen gerekir — Bölüm 4'ün tüm konusu budur. Promiscuous mode gerekli ama **yeterli değildir.**

> **Monitor mode (Wi-Fi'ye özel):** Promiscuous'tan farklıdır. **Monitor mode**, kablosuz kartı *hiçbir ağa bağlı olmadan* havadaki tüm 802.11 çerçevelerini (yönetim, kontrol, veri) dinler. Windows'ta Npcap çoğu kartta monitor mode'u kısıtlı destekler; Linux'ta `airmon-ng` ile etkinleşir. Wi-Fi'de "başka cihazların trafiği" için monitor mode gerekir — ama veri çerçeveleri WPA2/WPA3 ile şifrelidir (anahtarı bilmeden içeriği çözemezsin).

### Büyük pcap yönetimi & ring buffer
Uzun süreli yakalamada (saatler/günler) tek bir dev `.pcap` yönetilemez hale gelir (RAM dolar, Wireshark donar). Çözüm **ring buffer**:

- **Dosya boyutu/süreyle böl:** Capture Options → Output → "Create a new file automatically" → her **N MB** ya da **N saniye/dakika**'da yeni dosya.
- **Ring buffer:** "Use a ring buffer with N files" → yalnızca son N dosyayı tutar, eskiyi siler → disk dolmaz, **sonsuza dek** yakalayabilirsin (son N dosya penceresi).
- **dumpcap ile yakala (önerilen):** GUI Wireshark yakalarken *aynı anda* arayüzü çizmeye/ayrıştırmaya çalışır → yüksek hacimde **paket düşürür (drop)**. **dumpcap** sadece yazar, ayrıştırmaz → en az kayıp. Uzun yakalamada her zaman dumpcap kullan, sonra dosyaları Wireshark'la aç.

```bash
# dumpcap ile ring buffer: her 100 MB'da yeni dosya, en fazla 50 dosya tut (~5 GB pencere)
dumpcap -i eth0 -b filesize:100000 -b files:50 -w /capture/sensor.pcapng

# Zaman bazlı: her 60 sn'de yeni dosya, son 1440 dosya (24 saat)
dumpcap -i eth0 -b duration:60 -b files:1440 -w /capture/rolling.pcapng

# Yalnızca ilgili trafiği yakala (capture filter ile gürültüyü baştan kes — Bölüm 3)
dumpcap -i eth0 -f "not port 22 and not arp" -b filesize:100000 -b files:50 -w /capture/clean.pcapng
```

> 🔑 **Püf — pcapng vs pcap:** Modern format **pcapng**'dir (çoklu arayüz, paket yorumları/anotasyon, nanosaniye zaman damgası, arayüz metadata'sı destekler). Eski `.pcap` daha taşınabilir ama metadata fakiri. Analiz/anotasyon için **pcapng**; başka araçlara devrederken bazen `.pcap`'e çevirmen gerekebilir (`editcap -F libpcap in.pcapng out.pcap`).

---

<a id="3"></a>
## 3. 🔥 Capture Filter (BPF) vs Display Filter — En Sık Karıştırılan Konu

Bu, Wireshark öğrenenlerin **en çok karıştırdığı ve en kritik** ayrımdır. İkisi tamamen farklı dil, farklı an, farklı amaç.

```
┌──────────────────────────────────────────────────────────────────┐
│                         YAKALAMA AKIŞI                            │
│                                                                  │
│   Ağ kartı                                                       │
│      │                                                           │
│      ▼                                                           │
│  ┌────────────────────┐                                         │
│  │  CAPTURE FILTER    │  ← BPF dili (tcpdump). Çekirdek/sürücü   │
│  │  (BPF)             │    seviyesinde, YAKALAMADAN ÖNCE eler.   │
│  │  "port 53"         │    Elenen paket ASLA kaydedilmez/görülmez│
│  └────────┬───────────┘    → geri alınamaz, ama performans+disk  │
│           │                  kazandırır.                          │
│           ▼                                                       │
│      pcap dosyası (yalnızca geçenler)                            │
│           │                                                       │
│           ▼                                                       │
│  ┌────────────────────┐                                         │
│  │  DISPLAY FILTER    │  ← Wireshark dili. Yakalanmış paketleri  │
│  │  "dns"             │    EKRANDA filtreler. Paket diskte durur,│
│  └────────┬───────────┘    sadece görünmez → istediğinde geri al.│
│           ▼                                                       │
│      Ekranda görünen paketler                                   │
└──────────────────────────────────────────────────────────────────┘
```

### Capture Filter (BPF — Berkeley Packet Filter)
- **Ne zaman:** Yakalama **başlamadan önce**, çekirdek/sürücü seviyesinde.
- **Dil:** **BPF** (tcpdump söz dizimi) — Wireshark display diliyle **alakasız**.
- **Etki:** Filtreye uymayan paket **hiç yakalanmaz** → diske yazılmaz, RAM'e girmez, **sonradan geri getirilemez.**
- **Amaç:** Gürültüyü baştan kesmek, **paket düşürmeyi önlemek** (yüksek hacimde kritik), disk/RAM tasarrufu.
- **Nerede girilir:** Capture → Options → "Capture filter for selected interfaces" kutusu (yakalama başlamadan).

```bash
# BPF örnekleri (capture filter)
host 192.168.1.50                  # yalnızca bu IP'nin trafiği
net 10.0.0.0/8                     # bu alt ağ
port 53                            # yalnızca DNS (TCP+UDP 53)
tcp port 443                       # yalnızca HTTPS
not port 22 and not port 3389      # SSH ve RDP gürültüsünü at
src host 10.0.0.5 and dst port 80  # bu kaynaktan giden HTTP
udp and port 53                    # DNS sorgu/yanıt
host 8.8.8.8 or host 1.1.1.1       # iki DNS sunucusu
ether host aa:bb:cc:dd:ee:ff       # belirli MAC
vlan 100                           # belirli VLAN
```

> ⚠️ **BPF'in tehlikesi — geri dönüşsüzlük:** Capture filter ile bir şeyi *elersen*, o paket **yoktur.** Yakalama bitince "keşke onu da tutsaydım" diyemezsin. Bu yüzden BPF'i **emin olduğun gürültü** için kullan (örn. kendi SSH yönetim oturumun, bilinen yedekleme trafiği). Soruşturmada *belki lazım olur* dediğin hiçbir şeyi BPF ile atma.

### Display Filter (Wireshark görüntü filtresi)
- **Ne zaman:** Yakalama **sonrasında/sırasında**, paketler zaten diskteyken.
- **Dil:** **Wireshark display filter** dili (`dns`, `ip.addr`, `tcp.flags.syn`…) — BPF'ten **tamamen farklı**.
- **Etki:** Paket diskte **durur**, sadece **ekranda gizlenir.** Filtreyi temizleyince hepsi geri gelir. **Kayıp yok**, istediğin kadar dene-yanıl yap.
- **Amaç:** Yakalanmış devasa pcap içinde **avlamak**, daraltmak, kanıt bulmak.
- **Nerede girilir:** Ana penceredeki üst filtre çubuğu (yeşil = geçerli, kırmızı = sözdizimi hatası, sarı = riskli/uyarı).

```
# Display filter örnekleri (Wireshark dili)
dns                                # tüm DNS
ip.addr == 192.168.1.50            # bu IP kaynak VEYA hedef
tcp.port == 443                    # HTTPS (her iki yön)
http.request.method == "POST"      # yalnızca HTTP POST
tls.handshake.type == 1            # TLS Client Hello
tcp.flags.syn == 1 && tcp.flags.ack == 0   # bağlantı kurma denemeleri
frame contains "password"          # ham bayt içinde "password" geçen
```

> 🔥 **KARIŞIKLIĞIN KAYNAĞI — aynı şeyi iki farklı dilde yazmak:** "Sadece DNS yakala" demek:
> - **Capture filter** olarak: `port 53` (BPF)
> - **Display filter** olarak: `dns` (Wireshark)
>
> İkisini **karıştırırsan hata alırsın.** Capture filter kutusuna `dns` yazarsan **hata** (BPF `dns` bilmez). Display filter kutusuna `port 53` yazarsan **kırmızı/hata** (Wireshark `port` anahtar kelimesini öyle kullanmaz — `tcp.port == 53 || udp.port == 53` ister). **Hatırlatma:** *Yakalamadan önce = BPF (port 53). Yakaladıktan sonra = Wireshark (dns).*

### Hangisini ne zaman?

| Durum | Kullan |
|---|---|
| Yüksek hacimli ağ, paket düşürme riski | **Capture filter (BPF)** — gürültüyü baştan kes |
| Disk/RAM kısıtlı, uzun yakalama | **Capture filter (BPF)** |
| "Belki lazım olur" — emin değilsin | **Hiç filtreleme**, hepsini yakala, sonra **display filter** ile ara |
| Soruşturma/adli — hiçbir kanıtı kaybetme | **Hepsini yakala** + **display filter** |
| Yakalanmış pcap içinde avlama | **Display filter** (zaten capture geçti) |

> 🧠 **Altın kural:** **Capture filter ile yakalarken kaybettiğin paketi geri alamazsın; display filter ile gizlediğin paket her zaman oradadır.** Şüphedeysen daha fazla yakala, sonra ekranda filtrele. BPF'i yalnızca *performans/disk zorunluluğu* ve *kesin gürültü* için kullan.

---

<a id="4"></a>
## 4. 🔥 Yakalama NOKTASI — TAP vs SPAN vs Hub & Switch Körlüğü

Filtreden bile önce gelen soru: **"Paketi NEREDE yakalıyorsun?"** Yanlış noktada, en iyi filtre bile boş pcap üretir. Bu, sahada en çok zaman kaybettiren ve en az anlaşılan konudur.

### Neden switch'te "her şeyi göremezsin"
Eski **hub**, gelen her çerçeveyi **tüm portlara** kopyalardı → herkes herkesi görürdü (yakalama kolaydı ama ağ yavaştı). Modern **switch** ise MAC tablosu tutar ve çerçeveyi **yalnızca hedef portuna** yollar. Sonuç:

```
        ESKİ HUB (her şey her yere)          MODERN SWITCH (yalnızca hedefe)
        ┌──────────────────┐                 ┌──────────────────┐
   A ───┤                  ├─── B            A ───┤   MAC tablo   ├─── B
        │   tüm portlara   │                      │   A↔port1     │
   C ───┤   YAYINLAR       ├─── Sen          C ───┤   B↔port2     ├─── Sen
        └──────────────────┘                      │   ...         │   (port4)
   Sen A↔B trafiğini GÖRÜRSÜN.               └──────────────────┘
                                             A→B trafiği port4'e GELMEZ.
                                             Sen YALNIZCA kendi trafiğini görürsün.
```

> 🔥 **Switch körlüğü:** Sıradan bir switch portuna takılı dizüstünde Wireshark açarsan, promiscuous mode açık olsa bile **yalnızca kendine gelen/giden + broadcast/multicast** trafiği görürsün. Komşu makinelerin birbiriyle ya da gateway ile konuşması **sana hiç gelmez.** "Ağdaki tehdidi göremiyorum" şikâyetinin %1 numaralı sebebi budur. Çözüm: doğru yakalama noktası seç.

### Yakalama noktası seçenekleri

| Yöntem | Nasıl çalışır | Artı | Eksi |
|---|---|---|---|
| **Network TAP** | Hatta seri bağlanan donanım; trafiği **fiziksel olarak kopyalar** | **En güvenilir**, paket düşürmez, gizli (fark edilmez), tam-dupleks | Donanım gerekir, hatta erişim/kesinti |
| **Port Mirroring (SPAN)** | Switch'i "şu portları/VLAN'ı şu monitör portuna kopyala" diye yapılandırırsın | Donanımsız (switch destekliyorsa), esnek | **Yüksek yükte paket düşürür**, switch CPU'su sınırlı, yapılandırma hatası |
| **Hub (eski)** | Her şeyi her porta yayınlar | Ucuz, basit | Bulması zor, yarı-dupleks, çarpışma, sadece 10/100 Mbps |
| **Inline (bridge/PC iki kart)** | Trafiği bir karttan alıp diğerinden ver, arada yakala | Tam görüş | Hattı sen taşırsın (arıza = kesinti), gecikme |
| **Gateway/sunucuda yerel** | İlgilendiğin makinenin *kendi* arayüzünde yakala | İzinli, kolay, o makinenin tüm trafiği | Yalnızca o makineyi görürsün, makine ele geçtiyse güvenilmez |
| **Sanal switch (vSwitch) port mirror** | Hypervisor'da promiscuous/mirror portu | VM ağı için ideal | Hypervisor yapılandırması, kaçak trafik |

### Stratejik konumlandırma
- **Tüm dış (internet) trafiğini görmek istiyorsan:** Yakalama noktasını **internet gateway/firewall ile switch arasına** koy (TAP) ya da firewall'un WAN/LAN portunu **SPAN**'le. C2/exfil avı için **en değerli nokta budur** — her şey buradan geçer.
- **Belirli bir şüpheli makineyi izlemek:** O makinenin switch portunu mirror'la ya da **makinenin kendisinde** yakala (izin varsa).
- **Lateral movement (iç→iç) görmek:** Tek bir uplink noktası yetmez — iç trafik switch içinde dağıtık akar. **Birden çok SPAN** ya da çekirdek switch'te VLAN mirror gerekir; çoğu kuruluş bunu **Zeek/NDR sensörleriyle** dağıtık yapar.

### ⚠️ ARP spoofing — etik ve yasal SINIR
İnternet "ARP spoof ile switch'te başkasının trafiğini Wireshark'a çekersin" der. Teknik olarak mümkündür (kendini gateway gibi tanıtıp ortadaki-adam olursun, MITM), **AMA:**

> 🔥 **ARP spoofing AKTİF bir saldırıdır, pasif yakalama değildir.** Trafiği **manipüle eder** (yeniden yönlendirir), ağı bozabilir, ve **yetkisiz yapıldığında açıkça suçtur** (TCK 243/244, "bilişim sistemine müdahale"). Wireshark'ın pasif/yasal ruhuyla çelişir. **Yalnızca:** (a) tamamen kendi laboratuvar ağında, (b) yazılı izinli bir pentest kapsamında, (c) eğitim amaçlı izole ortamda yapılır. Üretim ağında ya da başkasının trafiğini görmek için **asla.** Doğru yol her zaman **TAP/SPAN** (pasif, izinli) ya da **kendi makinende** yakalamaktır.

> 🧠 **Altın kural — körlüğünü bil:** Bir pcap'i analiz etmeden önce **"Bu nereden yakalandı? Neyi göremiyorum?"** diye sor. Sadece bir uç noktada yakaladıysan, ağın geri kalanı kördür. Sadece gateway'de yakaladıysan, iç→iç (lateral) trafiği kaçırırsın. **Gördüğün, yakalama noktanın gösterdiği kadarıdır** — bütünü değil.

---

<a id="5"></a>
## 5. 🎯 Display Filter Ustalığı — Söz Dizimi & Hazır Kütüphane

Display filter, Wireshark'ın asıl gücüdür. Mantığı: **`protokol.alan operatör değer`**.

### Operatörler

| Operatör | Anlam | Alternatif | Örnek |
|---|---|---|---|
| `==` | eşit | `eq` | `ip.addr == 10.0.0.1` |
| `!=` | eşit değil | `ne` | `tcp.port != 443` |
| `>` `<` `>=` `<=` | karşılaştırma | `gt lt ge le` | `tcp.len > 1000` |
| `&&` | VE | `and` | `ip.addr==10.0.0.1 && tcp.port==80` |
| `\|\|` | VEYA | `or` | `tcp.port==80 \|\| tcp.port==443` |
| `!` | DEĞİL | `not` | `!arp` |
| `contains` | bayt/metin içerir | — | `frame contains "MZ"` |
| `matches` | regex (PCRE) | `~` | `http.host matches "\\.ru$"` |
| `in` | küme üyeliği | — | `tcp.port in {80 443 8080}` |

### Temel alan filtreleri (en sık kullanılanlar)

```
# --- IP / adres ---
ip.addr == 192.168.1.50            # kaynak VEYA hedef bu IP (en çok kullanılan)
ip.src == 192.168.1.50             # yalnızca kaynak
ip.dst == 8.8.8.8                  # yalnızca hedef
ip.addr == 10.0.0.0/24             # alt ağ
!(ip.addr == 192.168.1.0/24)       # yerel ağ DIŞI (dışarıyla konuşan = ilginç)

# --- TCP / UDP / port ---
tcp.port == 443                    # bu port (her iki yön)
tcp.dstport == 3389                # hedef RDP
udp.port == 53                     # DNS
tcp.flags.syn == 1 && tcp.flags.ack == 0    # SYN (bağlantı kurma)
tcp.flags.reset == 1               # RST (reddedilen/kesilen bağlantı)
tcp.analysis.retransmission        # yeniden iletim (ağ sorunu/kayıp)
tcp.analysis.flags                 # Wireshark'ın işaretlediği TCP anomalileri
tcp.stream == 12                   # belirli bir TCP oturumu (Follow ile çıkar)

# --- HTTP ---
http                               # tüm HTTP
http.request                       # yalnızca istekler
http.request.method == "POST"      # POST (sıklıkla exfil/komut)
http.host == "evil.example.com"    # belirli host başlığı
http.user_agent contains "curl"    # User-Agent anomalisi
http.response.code == 200          # başarılı yanıtlar
http.request.uri contains "/admin" # şüpheli yol

# --- DNS ---
dns                                # tüm DNS
dns.flags.response == 0            # yalnızca sorgular
dns.qry.name contains "dropbox"    # ada göre
dns.flags.rcode == 3               # NXDOMAIN (domain yok) — DGA göstergesi
dns.qry.type == 16                 # TXT kaydı (tünelleme sık kullanır)
dns.count.answers == 0             # yanıtsız sorgular

# --- TLS / SSL ---
tls                                # tüm TLS
tls.handshake.type == 1            # Client Hello (JA3 burada)
tls.handshake.type == 2            # Server Hello (JA3S burada)
tls.handshake.extensions_server_name           # SNI alanı var
tls.handshake.extensions_server_name == "x.com"   # belirli SNI
tls.record.version == 0x0301       # TLS sürümü
x509ce.dNSName                     # sertifikadaki domain adları

# --- Çerçeve / ham içerik ---
frame contains "password"          # ham baytlarda metin
frame.len > 1400                   # büyük çerçeveler
frame.time >= "2026-05-31 14:00:00"   # zamana göre
icmp                               # ICMP (ping / tünel)
arp                                # ARP (spoofing/keşif)
eth.addr == aa:bb:cc:dd:ee:ff      # MAC
```

### Filtre kombinasyonu — gerçek av örnekleri

```
# Yerel ağdan DIŞARIYA giden HTTP POST (potansiyel exfil/C2)
http.request.method == "POST" && ip.src == 192.168.0.0/16 && !(ip.dst == 192.168.0.0/16)

# Bilinen kötü IP listesine giden HER ŞEY
ip.addr == 185.220.101.1 || ip.addr == 45.155.205.99

# NXDOMAIN patlaması (DGA — Domain Generation Algorithm göstergesi)
dns.flags.rcode == 3

# Şifresiz parola sızıntısı (HTTP Basic Auth ya da form)
http.authorization || (http.request.method == "POST" && frame contains "passwd")

# Standart-dışı portta TLS (gizlenmiş C2 sık 443 dışı port kullanır)
tls.handshake.type == 1 && !(tcp.port == 443)

# Belirli bir hosttan giden tüm dış TCP bağlantı kurma denemeleri
ip.src == 10.0.0.55 && tcp.flags.syn == 1 && tcp.flags.ack == 0 && !(ip.dst == 10.0.0.0/8)
```

> 🔥 **Püf — `ip.addr ==` vs `ip.src ==`:** `ip.addr == X` "kaynak **veya** hedef X" demektir (çift yönlü, en kullanışlı). `!(ip.addr == X)` ise "ne kaynak ne hedef X" — `ip.addr != X` ile **aynı DEĞİLDİR!** `ip.addr != X` mantıken "en az bir uç X değil" olur ve neredeyse her paketi geçirir (tuzak). Bir IP'yi **tamamen dışlamak** için her zaman `!(ip.addr == X)` yaz.

> 💡 **Püf — Sağ tık ile filtre üret:** Bir paketin herhangi bir alanına sağ tıkla → **"Apply as Filter" → "Selected"** → Wireshark senin için doğru display filter'ı yazar. Söz dizimini ezberlemeden, alan adlarını bu yolla öğrenirsin. **"Prepare as Filter"** ise filtreyi çubuğa yazar ama uygulamaz (üzerine ekleme yapmak için).

> 💡 **Püf — Coloring Rules (renklendirme):** View → Coloring Rules. TCP RST'leri kırmızı, retransmission'ları sarı, kendi tanımladığın "şüpheli IP" kuralını mor yaparsın → pcap'i kaydırırken anomaliler **göze çarpar**, filtre yazmadan. Hazır kural seti zaten gelir (Bad TCP = siyah/kırmızı). Av sırasında "Bad TCP" renkleri seni doğru pakete götürür.

---

<a id="6"></a>
## 6. 🔬 Follow Stream, Export Objects, Statistics, IO Graphs

Wireshark'ın "tek tek paket" görünümünden **oturum/istatistik** görünümüne çıkaran araçları — asıl analizin yapıldığı yer.

### Follow Stream (akışı izle) — kanıtın olduğu yer
Bir TCP/UDP oturumundaki **tüm paketleri yeniden birleştirip** insan-okunur tek bir konuşma olarak gösterir.

- Bir pakete sağ tık → **Follow → TCP Stream** (ya da UDP/HTTP/TLS/QUIC Stream).
- İstemci verisi (gönderilen) ve sunucu verisi (alınan) **farklı renkte** gösterilir.
- Otomatik olarak `tcp.stream == N` display filter'ı uygulanır → o oturuma kilitlenirsin.
- **CTI değeri:** HTTP isteğinin tam gövdesi, çalınan kimlik bilgisi, C2 komutu, sızdırılan dosyanın baytları **burada** görünür (şifresizse). "Show data as" → ASCII / Hex Dump / Raw / C Arrays seçilebilir.

> 🔥 **Püf — "follow stream"de kanıt:** Bir C2 ya da exfil şüphesinde, ilgili pakete sağ tıkla → Follow TCP Stream → konuşmanın **tamamını** bir ekranda gör. Tek tek paketlere bakmak yerine "bu oturumda ne konuşuldu" sorusunun cevabı budur. **Ama TLS ise** (Bölüm 9 olmadan) yalnızca şifreli baytları görürsün — içerik kapalıdır, sadece *kiminle/ne kadar* konuştuğu açıktır.

### Export Objects — dosya çıkarma (carving)
Trafikten **transfer edilen dosyaları** ayıkla:
- **File → Export Objects → HTTP** (ya da SMB, FTP-DATA, TFTP, IMF, DICOM).
- Listeden dosyayı seç → **Save** → diske yaz.
- **CTI değeri:** Bir HTTP/SMB oturumunda indirilen **malware payload**'ını, sızdırılan belgeyi, ya da C2'nin gönderdiği aracı **doğrudan diske çıkarıp** hash'ini alır, sandbox'a atarsın. (TLS şifreliyse çıkmaz — yine Bölüm 9 gerekir.)

```bash
# tshark ile HTTP nesnelerini topluca dışa aktar (otomasyon)
tshark -r capture.pcapng --export-objects http,/cikti/http_objeler/
tshark -r capture.pcapng --export-objects smb,/cikti/smb_objeler/
```

### Statistics — büyük resmi gör
Menü: **Statistics**. CTI'de en kritik üçü:

| Araç | Ne gösterir | CTI kullanımı |
|---|---|---|
| **Protocol Hierarchy** | Trafiğin protokol dağılımı (% bayt/paket) | "Bu pcap'te ne var?" — anormal protokol (beklenmedik IRC, çok DNS) hemen göze çarpar |
| **Conversations** | Tüm konuşmalar (IP↔IP, port↔port), bayt/paket/süre | **En çok konuşan çiftler**; "şu iç IP şu dış IP ile niye bu kadar konuşmuş?" |
| **Endpoints** | Tüm uç noktalar (IP/MAC/port), trafik hacmi, GeoIP | Dışarıdaki **en aktif hedefler**; coğrafya |
| **DNS** | DNS sorgu/yanıt istatistiği | Anormal sorgu hacmi, NXDOMAIN oranı |
| **HTTP → Requests** | Tüm istenen host/URI'ler | İstenen tüm domainler tek listede → IOC çıkarma |
| **IO Graph** | Zaman içinde trafik (paket/bayt/bps) | **Beacon tespiti** (aşağıda) — düzenli sivri uçlar |

> 🔥 **Püf — Conversations'ta sırala & filtrele:** Statistics → Conversations → **TCP** sekmesi → "Bytes" sütununa göre sırala → en hacimli oturum en üstte. **Büyük giden (Tx) hacim = potansiyel exfil.** "Limit to display filter" kutusuyla yalnızca filtrelediğin trafiği say. Bir satıra sağ tıkla → **Apply as Filter** → o konuşmaya odaklan. Bu, "devasa pcap'te nereden başlasam" sorusunun en hızlı cevabıdır.

### IO Graph — zaman ekseninde trafik & beacon avı
**Statistics → I/O Graph.** X ekseni zaman, Y ekseni trafik (paket/sn, bayt/sn, ya da bir filtrenin değeri). Birden çok grafik üst üste, her birine ayrı display filter + renk verirsin.

- **Aralık (interval):** 1 sn, 100 ms, 1 dk seçilebilir — beacon periyodunu yakalamak için ayarla.
- **Y ekseni:** "Packets", "Bytes", "Bits", ya da `AVG/MAX/MIN(alan)`.
- **Beacon deseni:** Bir C2 her N saniyede bir "patronu arar" → IO Graph'ta **düzenli, eşit aralıklı sivri uçlar** (tarak/testere dişi deseni) belirir. İnsan trafiği düzensizdir; **mekanik düzenlilik = otomasyon = beacon şüphesi.** (Bölüm 8'de derin.)

```
IO Graph — beacon deseni (her 60 sn'de bir bağlantı):

paket/sn
  │
 4│   █          █          █          █          █
 3│   █          █          █          █          █
 2│   █          █          █          █          █
 1│   █    .  .  █   .   .  █  .    .  █   .    .  █
  └───┴──────────┴──────────┴──────────┴──────────┴──── zaman
     0s         60s        120s       180s       240s
     ▲ Eşit aralık + benzer hacim = MEKANİK = beacon (insan değil)
```

---

<a id="7"></a>
## 7. 🧬 Protokol-Protokol Analiz — DNS / HTTP / TLS / ICMP / SMB

Her protokolün kendi "kötüye kullanım imzaları" vardır. CTI analisti bunları tanır.

### 7.1 DNS — en çok kötüye kullanılan protokol
DNS neredeyse hiçbir yerde engellenmez (ağ çalışsın diye açık) → saldırganın en sevdiği kanal.

- **DNS tünelleme:** Veri, DNS sorgularının **subdomain alanına** (örn. `<base32-veri>.tunnel.evil.com`) ya da yanıtların **TXT kaydına** gömülür. İmza: anormal **uzun/rastgele görünümlü subdomain'ler**, çok sayıda **TXT/NULL** sorgusu, tek bir domain'e **olağanüstü yüksek sorgu hacmi**.
  ```
  dns.qry.name.len > 50                          # anormal uzun isim
  dns.qry.type == 16                             # TXT (tünel sık kullanır)
  dns.qry.type == 10                             # NULL kaydı
  dns && dns.qry.name matches "[a-f0-9]{20,}"    # hex/rastgele subdomain
  ```
- **DGA (Domain Generation Algorithm):** Malware, C2 domain'ini algoritmayla üretir (`kq3v9zx1p.com`, `xn7gh2la.net`…). Çoğu kayıtlı değildir → **NXDOMAIN patlaması.** İmza: kısa sürede çok sayıda `dns.flags.rcode == 3` (NXDOMAIN), rastgele görünümlü isimler.
  ```
  dns.flags.rcode == 3                           # NXDOMAIN — DGA göstergesi
  ```
- **NXDOMAIN patlaması:** Tek bir host saniyeler içinde onlarca başarısız çözümleme yapıyorsa → DGA/malware C2 arayışı.

### 7.2 HTTP — açık metin altın madeni (giderek azalıyor)
Şifresiz HTTP'de her şey okunur. Ama dünya HTTPS'e geçtikçe HTTP trafiği "eski/şüpheli" hale geliyor.

- **User-Agent anomalisi:** Meşru tarayıcılar standart UA gönderir. Malware sık sık **garip/boş/tekil** UA kullanır (`Mozilla/4.0`, `curl/7.x`, rastgele dize, ya da hiç). İmza:
  ```
  http.user_agent contains "curl" || http.user_agent contains "python" || http.user_agent contains "powershell"
  http.request && !http.user_agent             # UA olmayan istekler (şüpheli)
  ```
- **Beacon (HTTP C2):** Düzenli aralıklı GET/POST, aynı URI kalıbına (örn. `/api/v1/beacon`, `/submit.php`, `/jquery.min.js` taklidi). Sabit boyutlu yanıtlar, "200 OK" ama içerik tuhaf.
- **Anormal POST:** Büyük POST gövdeleri = veri sızması ya da çalınan kimlik bilgisi yüklemesi.
- **URI kalıbı:** Base64/hex kodlu uzun URI'ler, `.php?id=<uzun>` → komut/veri taşıyor olabilir.

### 7.3 TLS — kimliği gizler, metadata'yı gizlemez
Trafiğin çoğu artık TLS (HTTPS). İçeriği **göremezsin** (Bölüm 9 olmadan), ama **el sıkışma (handshake) açık metindir** ve çok şey söyler.

- **SNI (Server Name Indication):** İstemci, hangi domain'e bağlandığını el sıkışmada **açık** söyler (ESNI/ECH yoksa). Bu, "şifreli trafik nereye gidiyor" sorusunun cevabıdır:
  ```
  tls.handshake.extensions_server_name                     # SNI alanı
  tls.handshake.extensions_server_name contains ".onion"   # (teorik)
  ```
- **Sertifika:** Server Hello'da sunucu sertifikası (genelde) açık görünür → **issuer, subject, geçerlilik, seri no** IOC olur. Kendinden imzalı (self-signed) ya da tuhaf CN'li sertifikalar şüphelidir.
  ```
  x509ce.dNSName                  # sertifikadaki domain'ler
  x509sat.printableString         # subject/issuer alanları
  ```
- **🔥 JA3 / JA3S parmak izi:** TLS Client Hello'nun parametreleri (sürüm, şifre listesi, uzantılar, eliptik eğriler) **istemci yazılımına özgü bir parmak izi** oluşturur. Bu parametrelerin MD5'i **JA3 hash**'idir. Sunucunun Server Hello'sundan üretileni **JA3S**'tir. İkisi birlikte (**JA3+JA3S**) bir istemci-sunucu çiftini tanımlar.
  - **Neden güçlü:** TLS içeriğini çözemesen bile, **belirli bir malware'in TLS imzası** (örn. belirli bir Cobalt Strike/Metasploit/Trickbot JA3'ü) tehdit istihbaratı listelerinde bilinir. **Şifreli C2'yi içeriğini açmadan tanırsın** — sadece el sıkışma parmak izinden.
  - **Sınır/teyit:** Wireshark sürümüne göre JA3 doğrudan bir alan olarak gelmeyebilir; çoğu ortamda **Zeek/Suricata/JA3 betikleri** ya da bir Lua eklentisi ile üretilir. *JA3'ün senin Wireshark sürümünde yerleşik alan olup olmadığını teyit et; yoksa Zeek `ssl.log` ya da `ja3` eklentisi kullan.* Ayrıca JA3, TLS kütüphanesi güncellemeleriyle ve **JA3 randomizasyonu/GREASE** ile değişebilir; tek başına kesin kanıt değil, **güçlü bir gösterge**dir.
- **Standart-dışı port + TLS:** 443 dışı bir portta TLS el sıkışması → gizlenmeye çalışan C2 olabilir.

### 7.4 ICMP & DNS — kovert kanal / exfiltration
- **ICMP tünelleme:** `ping` paketlerinin **veri (payload) alanı** keyfi veri taşıyabilir. Normal ping'in payload'ı sabit/öngörülebilirdir; **anormal büyük ya da değişken ICMP payload'ı** → ICMP tüneli/exfil.
  ```
  icmp && data.len > 64           # büyük ICMP payload'ı (tünel şüphesi)
  icmp.type == 8 || icmp.type == 0   # echo request/reply (ping)
  ```
- **DNS exfiltration:** (7.1) Veri DNS sorgularına gömülür — firewall'ı atlatmanın klasik yolu.

### 7.5 SMB & lateral movement
İç ağda yanal hareketin ana protokolleri:
- **SMB (445):** Dosya paylaşımı + PsExec/uzak yürütme. Anormal **iç→iç SMB**, `ADMIN$`/`C$` paylaşımına erişim, `.exe` yazma → lateral movement.
  ```
  smb2 || smb                     # SMB trafiği
  smb2.cmd == 5                   # SMB2 Create (dosya/pipe açma)
  smb2.filename contains ".exe"   # uzak .exe (PsExec imzası)
  ```
- **Kerberos (88):** Anormal bilet istekleri, **Kerberoasting** (çok sayıda servis bileti), Golden/Silver Ticket göstergeleri.
- **RDP (3389), WMI/RPC (135 + dinamik):** Uzak yönetim → lateral.
- **İmza:** İç ağda bir iş istasyonunun **başka iş istasyonlarına** SMB/RDP/WMI ile bağlanması (normalde sunuculara bağlanır, akranlarına değil) → yanal hareket alarmı.

> 🧠 **Protokol sezgisi:** Her protokol için "normal nasıl görünür?"ü öğren, sonra **sapmayı** ara. DNS'te normal = kısa isimler, az TXT. HTTP'de normal = tarayıcı UA'sı. TLS'te normal = bilinen JA3, geçerli sertifika, 443. SMB'de normal = iş istasyonu→sunucu. **Anomali = sapmadır**, ve CTI anomalinin peşindedir.

---

<a id="8"></a>
## 8. 🔥 Malware Trafiği Avı — Beaconing, C2, DNS Tüneli, Exfil

Bu bölüm rehberin kalbidir: gerçek bir tehdidi pcap'te **nasıl avlarsın.**

### 8.1 Beaconing tespiti — C2'nin nabzı
Çoğu C2, ele geçmiş makineye "uyu, sonra beni ara" der. Makine **düzenli aralıklarla** (her 30 sn / 5 dk / 1 saat) C2'yi yoklar = **beacon (işaret atışı).** Bu **düzenlilik** en güçlü tespit imzasıdır.

```
NORMAL İNSAN TRAFİĞİ (düzensiz, kümeli):
  █ ██   █        ███ █    █  ██        █     ███
  └─────────────────────────────────────────────► zaman
  (göz atma, mola, yine göz atma — öngörülemez)

BEACON (mekanik, eşit aralıklı):
  █         █         █         █         █
  └─────────────────────────────────────────────► zaman
  ◄──60s──► ◄──60s──► ◄──60s──► ◄──60s──►
  (her 60 sn ± küçük jitter — bir SAAT gibi)
```

**Wireshark'ta beacon avı:**
1. **Conversations** ile şüpheli bir iç↔dış IP çiftini bul (sürekli az miktarda trafik, uzun süre).
2. O çifte **Apply as Filter**.
3. **IO Graph** aç, aralığı beacon periyoduna göre ayarla (1 sn → 1 dk dene).
4. **Eşit aralıklı sivri uçlar** = beacon. İnsan trafiği bu deseni üretmez.
5. Paketler arası zamanı incele: bir paket seç → "Time since previous frame" → **sabit ya da dar aralıkta mı?**

> 🔥 **Püf — jitter'ı hesaba kat:** Gelişmiş C2'ler (Cobalt Strike vb.) tespitten kaçmak için **jitter** ekler: "her 60 sn" yerine "60 sn ± %30 rastgele" → IO Graph'ta desen biraz dağılır ama **hâlâ bir merkezi periyot** etrafındadır. Mükemmel düzenlilik yoksa beacon yok sanma; **ortalama aralık + dar varyans** hâlâ otomasyona işaret eder. İstatistiksel olarak: paketler-arası süreyi çıkar (tshark `-e frame.time_delta`), histogramını al → insan trafiği geniş/çok-tepeli, beacon dar/tek-tepeli.

> 🔥 **Püf — "low and slow":** Akıllı saldırgan beacon'ı **çok seyrek** (saatte bir, günde bir) ve **küçük** yapar → hacim alarmı tetiklemez, kısa bir pcap'te hiç görünmez. Tespit için **uzun süreli yakalama** (Bölüm 2 ring buffer) + uzun pencerede IO Graph gerekir. Kısa pcap "temiz" görünebilir ama low-and-slow beacon'ı kaçırmış olabilirsin — **süre bilincini** koru.

### 8.2 C2 kanal türleri & imzaları

| Kanal | Nasıl gizlenir | Wireshark imzası |
|---|---|---|
| **HTTP(S) C2** | Meşru web trafiği taklidi | Beacon deseni, tuhaf URI/UA, bilinen JA3, sabit yanıt boyutu |
| **DNS C2/tünel** | DNS hiç engellenmez | Uzun/rastgele subdomain, çok TXT, tek domain'e yüksek hacim |
| **TLS C2** | Şifreli → içerik gizli | JA3 imzası, kendinden imzalı sertifika, 443-dışı port, SNI yok/tuhaf |
| **ICMP C2** | Ping görünür | Büyük/değişken ICMP payload, sürekli ping |
| **Domain fronting** | CDN arkasına saklanır (SNI bir, Host başka) | SNI ile gerçek hedef uyuşmazlığı (TLS çözülürse) |

### 8.3 DNS tünelleme tespiti (adım adım)
1. `dns` filtrele, **Statistics → DNS** ile sorgu hacmine bak.
2. Tek bir domain'e **olağanüstü çok sorgu** mu? → şüphe.
3. `dns.qry.name.len > 50` → anormal uzun isimler.
4. `dns.qry.type == 16` (TXT) ya da `== 10` (NULL) yoğunluğu.
5. Subdomain'ler **base32/base64/hex** mi görünüyor (rastgele harf-rakam)? → veri taşıyor.
6. Follow ile birkaç sorguyu incele: subdomain'i decode et (base32/64) → gömülü veri çıkarsa **kanıt.**

### 8.4 Veri sızması (exfiltration) imzaları
- **Hacim anomalisi:** Conversations'ta **giden (Tx) bayt ≫ gelen (Rx)** olan dış bağlantılar. Normalde indirme (Rx) baskındır; **büyük upload** ters ve şüphelidir.
- **Beklenmedik hedef:** Veri tanımadık bir bulut/IP'ye, tuhaf saatte (gece 3) gidiyorsa.
- **Yavaş sızıntı:** Küçük parçalar halinde, uzun zamana yayılmış (low-and-slow) → toplam hacim büyük ama anlık küçük.
- **Protokol uyumsuzluğu:** 443 portunda ama TLS olmayan trafik; DNS/ICMP içinde veri.
- **Sıkıştırma/şifreleme imzası:** Exfil edilen veri sık sık şifreli/sıkıştırılmış (yüksek entropi) → Follow stream'de "rastgele bayt çorbası" görünür.

### 8.5 Kovert (gizli) kanallar
- **Timing channel (zamanlama kanalı):** Veri, paketlerin **gönderim zamanlamasına** kodlanır (1 = hızlı, 0 = gecikmeli) → içerik tamamen normal ama **aralıklar** mesaj taşır. Tespiti çok zordur; istatistiksel zamanlama analizi gerekir.
- **Padding/boyut kanalı:** Veri, paket **boyutlarına** kodlanır.
- **Protokol başlık alanları:** Kullanılmayan TCP/IP başlık bitlerine (IP ID, TCP seq başlangıcı, TTL) veri gizlenir.
- **Bu kanallar Wireshark'ta "normal" görünür** — tespit için derin istatistiksel/davranışsal analiz (Zeek, özel betik) gerekir. Bilinçli ol: *gördüğün normal, gizli kanal olmadığını kanıtlamaz.*

> 🧠 **Av zihniyeti:** Malware avı "imza eşleştirme" değil, **"normalden sapma"** avıdır. Önce ağının **temel çizgisini (baseline)** öğren: kim kiminle, ne sıklıkta, ne hacimde konuşur? Sonra **sapmayı** ara — yeni hedef, yeni periyot, ters hacim, tuhaf protokol. CTI'de "bilinen kötü"yü (IOC listesi) ararsın **ama** asıl ustalık "bilinmeyen kötü"yü davranıştan yakalamaktır.

---

<a id="9"></a>
## 9. 🔓 TLS Şifre Çözme — Ne Zaman Mümkün, Ne Zaman Değil

Trafiğin çoğu TLS. Wireshark **doğru anahtarla** TLS'i çözebilir — ama bu her zaman mümkün değildir. Bu bölüm "neyi çözebilirsin, neyi asla"yı netleştirir.

### Yöntem 1 — SSLKEYLOGFILE (en pratik, en güçlü)
Tarayıcılar (Chrome, Firefox) ve birçok uygulama, bir ortam değişkeni ayarlıysa **TLS oturum anahtarlarını** bir dosyaya yazar. Wireshark bu dosyayı okuyup trafiği çözer — **forward secrecy olsa bile** (çünkü oturum-başına anahtarın *kendisi* dosyadadır).

```bash
# Windows (PowerShell) — tarayıcıyı başlatmadan ÖNCE ayarla
$env:SSLKEYLOGFILE = "C:\Users\Yasin\tls_keys.log"

# Linux / macOS
export SSLKEYLOGFILE=~/tls_keys.log

# Sonra tarayıcıyı bu ortamdan aç, trafiği yakala.
```

Wireshark'ta: **Edit → Preferences → Protocols → TLS → "(Pre)-Master-Secret log filename"** → bu dosyayı göster. Artık o oturumların HTTP/2 içeriği, gönderilen veriler, her şey **açık** görünür (Follow → HTTP Stream çözülmüş gelir).

> 🔥 **Püf — SSLKEYLOGFILE neden işe yarar (forward secrecy'ye rağmen):** Modern TLS **ephemeral (geçici) anahtar değişimi** (ECDHE) kullanır → her oturumun anahtarı benzersiz, sunucunun uzun-ömürlü özel anahtarı bile geçmiş oturumları çözmeye **yetmez** (forward secrecy). AMA `SSLKEYLOGFILE`, oturum anahtarını **doğrudan istemciden** alır → forward secrecy'yi atlar, çünkü "anahtarı kırmıyoruz, istemci bize veriyor." **Koşul:** Yakalamayı yapan sen, istemci tarafına (kendi makinen/lab) erişebilmelisin. Başkasının trafiğini bu yöntemle çözemezsin (anahtar log'una erişimin yok).

### Yöntem 2 — RSA özel anahtarı (sınırlı, eskiyen yöntem)
Sunucunun **RSA özel anahtarına** sahipsen (örn. kendi sunucun) ve oturum **RSA anahtar değişimi** kullanıyorsa:
- **Edit → Preferences → Protocols → TLS → RSA Keys List** → sunucu IP/port + özel anahtar (`.pem`/`.key`).

> ⚠️ **Kritik sınır — RSA yöntemi forward secrecy'de ÇALIŞMAZ:** Modern TLS 1.2/1.3 neredeyse her zaman **ECDHE** (ephemeral) kullanır → RSA özel anahtarı oturum anahtarını **çözmeye yetmez.** RSA özel anahtarıyla çözme yalnızca **eski, RSA-key-exchange** oturumlarda işe yarar (giderek yok oluyor). TLS 1.3 zaten RSA anahtar değişimini **kaldırdı.** Yani: kendi sunucunun bile özel anahtarına sahipsen, **forward secrecy varsa (ki vardır) RSA yöntemi başarısızdır** — `SSLKEYLOGFILE`'a ya da sunucu-tarafı log'a (ya da bir TLS-terminating proxy) ihtiyacın olur.

### Ne zaman ASLA çözemezsin (dürüst sınır)

| Durum | Çözülür mü? | Neden |
|---|---|---|
| Kendi tarayıcın, SSLKEYLOGFILE ayarlı | ✅ | Oturum anahtarı elinde |
| Kendi sunucun, RSA key-exchange (eski) | ✅ | RSA özel anahtarı yeter |
| Kendi sunucun, ECDHE/forward secrecy | ❌ (RSA ile) | Ephemeral anahtar; SSLKEYLOGFILE/proxy gerek |
| **Başkasının TLS trafiği, anahtar yok** | ❌ **ASLA** | Anahtara erişimin yok — şifreleme işini yapıyor |
| TLS 1.3, anahtar log yok | ❌ | Forward secrecy + RSA yok |

> 🧠 **Altın kural — TLS çoğu zaman seni kör eder, ama metadata hâlâ konuşur:** Pratikte sahada yakaladığın TLS trafiğinin **içeriğini genelde çözemezsin** (anahtarın yok). Bu yüzden TLS analizi **metadata analizidir:** SNI (nereye), sertifika (kim), JA3/JA3S (hangi yazılım), zamanlama/hacim (beacon mı, exfil mi), port (standart mı). **"İçeriği göremiyorum" ≠ "hiçbir şey göremiyorum."** Bir analistin TLS karşısındaki gücü, **el sıkışma metadata'sını ve davranışı** okumaktır. (Bölüm 8 beacon/JA3 tam bu yüzden değerli.)

---

<a id="10"></a>
## 10. ⌨️ tshark CLI — Otomasyon, Toplu IOC Çıkarma, Dilimleme

**tshark**, Wireshark'ın komut satırı motorudur. GUI'nin yapamadığı şey: **otomasyon, scripting, devasa pcap'leri programatik işleme, toplu IOC çıkarma.** Senior bir analistin günlük ekmeği.

### Temel kullanım
```bash
# Bir pcap'i oku ve özetle
tshark -r capture.pcapng

# Display filter uygula (Wireshark dili — -Y)
tshark -r capture.pcapng -Y "dns.flags.rcode == 3"

# Capture filter ile canlı yakala (BPF — -f)
tshark -i eth0 -f "port 53" -w dns_only.pcapng

# İstatistik: protokol hiyerarşisi
tshark -r capture.pcapng -q -z io,phs

# İstatistik: konuşmalar (en hacimliyi bul)
tshark -r capture.pcapng -q -z conv,tcp
```

### 🔥 Alan (field) çıkarma — IOC otomasyonu
`-T fields -e <alan>` ile istediğin alanları CSV/satır olarak çıkar → grep/sort/uniq ile IOC üret:

```bash
# Tüm benzersiz DNS sorgularını çıkar (domain IOC listesi)
tshark -r capture.pcapng -Y "dns.flags.response == 0" \
  -T fields -e dns.qry.name | Sort-Object -Unique

# Tüm HTTP host + URI (web IOC)
tshark -r capture.pcapng -Y "http.request" \
  -T fields -e http.host -e http.request.uri -E separator=" "

# Tüm TLS SNI değerleri (şifreli trafiğin gittiği domainler)
tshark -r capture.pcapng -Y "tls.handshake.extensions_server_name" \
  -T fields -e tls.handshake.extensions_server_name | Sort-Object -Unique

# Tüm dış hedef IP'ler + port (bağlantı IOC)
tshark -r capture.pcapng -Y "ip.dst != 192.168.0.0/16 && tcp.flags.syn==1" \
  -T fields -e ip.dst -e tcp.dstport | Sort-Object -Unique

# Tüm User-Agent değerleri (anomali avı)
tshark -r capture.pcapng -Y "http.user_agent" \
  -T fields -e http.user_agent | Sort-Object -Unique

# Paketler-arası süre (beacon analizi için — bir IP çiftinde)
tshark -r capture.pcapng -Y "ip.addr==10.0.0.5 && ip.addr==185.1.2.3" \
  -T fields -e frame.time_delta_displayed
```

> 💡 **Not (Windows/PowerShell):** Yukarıda `Sort-Object -Unique` PowerShell'dir. Linux/macOS'ta `| sort -u` ya da `| sort | uniq -c | sort -rn` (sayıp sıralar) kullan. tshark'ın kendisi her iki platformda da aynı çalışır.

### 🔥 Büyük pcap'i dilimleme & ön işleme (editcap / mergecap / capinfos)
GUI Wireshark devasa bir pcap'i açmaya çalışırken donar. **Önce CLI ile dilimle, sonra parçayı GUI'de aç.**

```bash
# pcap hakkında özet bilgi (boyut, paket sayısı, süre)
capinfos capture.pcapng

# Zamana göre dilimle: yalnızca belirli zaman aralığı
editcap -A "2026-05-31 14:00:00" -B "2026-05-31 15:00:00" buyuk.pcapng dilim.pcapng

# Paket sayısına göre böl: her 100.000 pakette bir dosya
editcap -c 100000 buyuk.pcapng parca.pcapng

# Süreye göre böl: her 600 saniyede bir dosya
editcap -i 600 buyuk.pcapng parca.pcapng

# Birden çok pcap'i birleştir (zaman sırasına göre)
mergecap -w birlesik.pcapng dosya1.pcapng dosya2.pcapng dosya3.pcapng

# Belirli paket aralığını çıkar (örn. 5000–6000)
editcap -r buyuk.pcapng kesit.pcapng 5000-6000

# Format çevir (pcapng → eski pcap, başka araç için)
editcap -F libpcap in.pcapng out.pcap
```

> 🔥 **Püf — CLI ile filtrele, GUI ile incele:** En verimli akış: `tshark -r dev.pcapng -Y "<dar filtre>" -w kucuk.pcapng` ile devasa pcap'ten **yalnızca ilgili paketleri** yeni küçük bir dosyaya yaz, sonra `kucuk.pcapng`'i Wireshark GUI'de aç. 10 GB'lık pcap'i GUI'de açıp kasmaktansa, 50 MB'lık ilgili kesiti incele. tshark filtreleme + GUI görselleştirme = en güçlü kombinasyon.

### Otomatik IOC çıkarma scripti (kavram)
Bir pcap geldiğinde tek komutla tüm IOC'leri dökmek:
```bash
# Tek script: domain, SNI, dış IP, UA, dosya hash → ayrı dosyalara
tshark -r $PCAP -Y "dns.qry.name" -T fields -e dns.qry.name | sort -u > iocs/domains.txt
tshark -r $PCAP -Y "tls.handshake.extensions_server_name" -T fields -e tls.handshake.extensions_server_name | sort -u > iocs/sni.txt
tshark -r $PCAP -Y "ip.dst and not ip.dst in {10.0.0.0/8 192.168.0.0/16}" -T fields -e ip.dst | sort -u > iocs/ext_ips.txt
tshark -r $PCAP -Y "http.user_agent" -T fields -e http.user_agent | sort -u > iocs/uas.txt
tshark -r $PCAP --export-objects http,iocs/files/    # dosyaları çıkar → hash'le
```

---

<a id="11"></a>
## 11. 🔥 PÜF NOKTALARI — Piyasada Bulamayacakların

Bu bölüm, çoğu rehberin atladığı ve gerçek analizi **çökerten ya da kurtaran** detaylardır. **En önemli bölüm budur.**

### 11.1 Capture filter ile display filter'ı karıştırmak (en sık hata)
- **Belirti:** Capture filter kutusuna `dns` yazdın → hata. Ya da display çubuğuna `port 53` yazdın → kırmızı.
- **Gerçek:** İki ayrı dil. **Yakalamadan önce = BPF (`port 53`). Yakaladıktan sonra = Wireshark (`dns`).** (Bölüm 3.)
- **Etki:** Yeni başlayanların %90'ı buradan saat kaybeder. BPF'in geri-dönüşsüzlüğü ile birleşince, "yanlış capture filter koydum, kanıtı elemişim" felaketine döner.
- **Çözüm:** Şüphede **hiç capture filter koyma**, hepsini yakala, **display filter** ile oyna. BPF'i yalnızca emin gürültü + performans zorunluluğu için kullan.

### 11.2 BPF ile gürültü azaltma & performans (paket düşürme)
- **Sorun:** Yüksek hacimli ağda (1 Gbps+) Wireshark GUI ayrıştırırken **paket düşürür** (Statistics → Capture → "dropped" sayısına bak — sıfır olmalı). Düşen paket = kayıp kanıt.
- **Çözüm:** (a) **dumpcap** kullan (ayrıştırmaz, yazar). (b) **Capture filter (BPF)** ile gürültüyü çekirdek seviyesinde at (örn. yedekleme/replikasyon trafiğini, kendi SSH oturununu). (c) Ring buffer ile parçala. (d) Snap length kısıtla (`-s 96` sadece başlıkları al — payload gerekmiyorsa).
- **Püf:** "Dropped packets" sıfırdan büyükse analizine **güvenme** — eksik veri yanlış sonuç verir.

### 11.3 Beacon tespiti — IO Graph + jitter okuma
- (Bölüm 8.1 derinliği.) İnsan trafiği düzensiz, beacon mekaniktir. IO Graph'ta **eşit aralıklı sivri uçlar** ara. Jitter eklenmişse mükemmel düzen bozulur ama **dar varyanslı bir merkezi periyot** kalır.
- **Püf:** `frame.time_delta` histogramı al; insan = geniş/çok-tepeli, beacon = dar/tek-tepeli. Mükemmel düzen yok diye beacon yok deme.

### 11.4 JA3/JA3S ile şifreli C2'yi içeriğini açmadan tanımak
- İçeriği çözemediğin TLS C2'sini, **TLS el sıkışma parmak izinden** (JA3) tanıyabilirsin. Bilinen malware JA3'leri tehdit istihbaratında listelidir.
- **Sınır:** Wireshark'ta JA3 yerleşik alan olmayabilir (sürüme bağlı — **teyit et**); çoğu ortamda Zeek/Suricata/eklenti üretir. GREASE/randomizasyon JA3'ü oynatabilir → kesin değil, güçlü gösterge.

### 11.5 DNS tüneli tespiti
- Uzun/rastgele subdomain, çok TXT/NULL sorgusu, tek domain'e yüksek hacim. (Bölüm 8.3.) DNS hiç engellenmediği için en sinsi exfil/C2 kanalı.
- **Püf:** `dns.qry.name.len` üzerinden sırala; insan domain'leri kısa, tünel domain'leri anormal uzun. Subdomain'i base32/64 decode et → gömülü veri çıkarsa kanıt.

### 11.6 Promiscuous mode yanılgısı (yine — çünkü en kritik)
- Promiscuous mode "ağdaki herkesi görürüm" **demek değildir.** Switch'te yalnızca sana ulaşan çerçeveleri görürsün. (Bölüm 3 & 4.)
- **Püf:** Promiscuous açık ama boş pcap → sorun **yakalama noktası**, mod değil. Wi-Fi'de başka cihazları görmek için **monitor mode** (ayrı kavram) gerekir, o da şifreli veriyi açmaz.

### 11.7 Yakalama noktası körlüğü (sahada 1 numaralı zaman kaybı)
- "Tehdidi göremiyorum" → muhtemelen yanlış noktadasın. Tek uç noktada yakaladıysan ağ kördür; gateway'de yakaladıysan lateral (iç↔iç) kördür. (Bölüm 4.)
- **Püf:** Analizden önce **"Bu nereden yakalandı, neyi göremiyorum?"** diye sor. Gördüğün, noktanın gösterdiği kadarıdır. C2/exfil için **gateway/firewall TAP**; lateral için dağıtık SPAN/Zeek.

### 11.8 TLS'in çoğu trafiği gizlemesi — ama metadata hâlâ konuşur
- Sahadaki TLS'in içeriğini genelde çözemezsin (anahtar yok). **Ama** SNI (nereye), sertifika (kim), JA3 (hangi yazılım), zamanlama/hacim (beacon/exfil), port (standart mı) **açıktır.** (Bölüm 9.)
- **Püf:** "Şifreli, demek ki bakacak bir şey yok" deme — **metadata analizi** TLS karşısında asıl silahtır. Beacon ve JA3 tam bu yüzden değerli.

### 11.9 Büyük pcap'i tshark/editcap ile dilimleme
- GUI devasa pcap'te donar/çöker. **Önce CLI ile dilimle** (zaman/sayı/filtre), sonra küçük kesiti GUI'de incele. (Bölüm 10.)
- **Püf:** `tshark -r big -Y "<filtre>" -w small` → ilgili paketleri ayır → GUI'de aç. 10 GB'ı açmaya çalışma; 50 MB'ı incele.

### 11.10 Coloring rules — anomaliyi göze çarpan yap
- View → Coloring Rules. Bad TCP kırmızı, retransmission sarı, "şüpheli IP" kendi kuralın mor. Pcap'i kaydırırken **renk** seni anomaliye götürür, filtre yazmadan. (Bölüm 5.)
- **Püf:** Soruşturma başında bilinen-kötü IP'lere parlak bir renk kuralı ata → o trafik tüm pcap boyunca **yanıp söner**.

### 11.11 "Follow stream"de kanıtı bulmak
- Tek tek pakete değil, **oturumun bütününe** bak: sağ tık → Follow → TCP/HTTP Stream → konuşmanın tamamı tek ekranda. Çalınan kimlik, C2 komutu, exfil edilen veri **burada.** (Bölüm 6.)
- **Püf:** Şifresizse altın; TLS ise yalnızca şifreli bayt görürsün (Bölüm 9 gerekir). "Show as" ile ASCII/Hex geçiş yap.

### 11.12 Zaman damgası & saat senkronizasyonu (korelasyonun temeli)
- Birden çok kaynaktan (pcap + firewall log + EDR) olayları **eşleştirmek** için saatlerin senkron olması (NTP, tercihen **UTC**) şart. Sensörün saati 3 dk kaymışsa, olay zaman çizgin yanlış olur ve korelasyon çöker.
- **Püf:** Wireshark'ta **View → Time Display Format → UTC** kullan (yerel saat dilimi karışıklık yaratır). **Time Reference** (Ctrl+T) ile bir paketi "sıfır" yap → ondan sonraki göreli süreleri ölç (beacon aralığı, gecikme analizi için ideal). Farklı pcap'leri birleştirirken (`mergecap`) zaman sırasına dikkat.

### 11.13 Gizli veri kanalı (timing/padding) — "normal" yanıltır
- Veri, paket **zamanlamasına** ya da **boyutuna** kodlanabilir (kovert kanal) → içerik tamamen normal görünür ama mesaj taşır. (Bölüm 8.5.) Wireshark'ta standart filtrelerle görünmez.
- **Püf:** *Gördüğün "normal", gizli kanal olmadığını kanıtlamaz.* Yüksek-güven ortamında zamanlama/boyut istatistiği (Zeek, özel analiz) gerekir. Bu, "analiz temiz" demenin sınırını hatırlatır.

### 11.14 GeoIP yanıltması (coğrafya kesin kimlik değil)
- Wireshark **MaxMind GeoIP** veritabanıyla IP'leri ülkeye eşleyebilir (Edit → Preferences → Name Resolution → MaxMind database directory; Statistics → Endpoints'te ülke sütunu).
- **Püf — yanıltma:** GeoIP **kesin değildir.** Saldırgan **VPN/proxy/bulut/ele geçmiş ara sunucu** kullanır → "trafik Hollanda'ya gidiyor" aslında C2'nin gerçek yeri olmayabilir; sadece bir sıçrama noktası. Bulut IP'leri (AWS/Azure) coğrafyayı tamamen anlamsızlaştırır. GeoIP **bir ipucu**, kimlik atfı (attribution) **değildir.** Ülkeye bakıp "Rus saldırısı" deme — altyapı kiralık/ele geçmiş olabilir.

### 11.15 Yakalananın hassas veri içermesi (KENDİ OPSEC'in)
- Bir pcap **her şeyi** içerir: şifresiz parolalar, çerezler/oturum token'ları, kişisel veri, dahili IP/topoloji, e-posta içerikleri, API anahtarları. **Bir pcap'i paylaşmak, o trafikteki tüm sırları paylaşmaktır.**
- **Püf — pcap hijyeni:** (a) pcap'leri **şifreli** sakla (VeraCrypt — çapraz referans). (b) Paylaşmadan önce **anonimleştir/sanitize et** (`tracewrangler`, `pcapng` düzenleme — IP/MAC/payload maskeleme). (c) İhtiyacın bittiğinde **güvenli sil** (`/imha`). (d) Yakalama yaparken **kendi kimlik bilgilerin de** pcap'e düşer — yetkili ortamda bile kendi parolanı HTTP'den girme. **Analist olarak topladığın kanıt, sızdırılırsa silahtır.**

### 11.16 Name resolution'ın yakalamayı kirletmesi
- Wireshark, IP'leri isimlere çevirmek için **DNS sorgusu yapabilir** (Edit → Preferences → Name Resolution → "Resolve network addresses"). **Bu, yakaladığın ağa senin ürettiğin ek DNS trafiği ekler** — hem pcap'ini kirletir hem gizliliğini bozar (sen orada olduğunu ağa duyurursun).
- **Püf:** Pasif/gizli analizde **"Resolve network (IP) addresses"** kapat (canlı DNS sorgusu yapmasın). MAC üreticisi çözümü (yerel, OUI tablosundan) güvenlidir; ağ adresi çözümü ağa paket gönderir → kapalı tut.

### 11.17 Expert Info — Wireshark sana ipucu verir
- **Analyze → Expert Information**: Wireshark'ın otomatik tespit ettiği anomaliler (retransmission, duplicate ACK, RST, malformed paket, "connection refused") renk-kodlu (Error/Warn/Note/Chat) listelenir.
- **Püf:** Bir pcap'i açar açmaz Expert Info'ya bak → ağ sorunları, taranan kapalı portlar (çok RST = port tarama), bozuk paketler hemen görünür. Lateral movement keşfi (port tarama) sık sık **bir sürü RST/SYN** olarak Expert Info'da belirir.

### 11.18 İsim çözümleme ile MAC/OUI → cihaz türü ipucu
- MAC adresinin ilk yarısı (**OUI**) üreticiyi gösterir (Wireshark yerel tabloyla çözer, ağa paket göndermeden). Bir cihazın üreticisi ("Raspberry Pi", "Hikvision kamera", "Apple") ne tür cihaz olduğuna dair **bedava ipucudur.**
- **Püf:** Ağda beklenmedik bir üreticinin cihazı (örn. tanımadığın bir IoT/kamera OUI'si) → yetkisiz/rogue cihaz olabilir. Statistics → Endpoints → Ethernet sekmesinde üreticileri tara.

---

<a id="12"></a>
## 12. ☠️ Yaygın Ölümcül Hatalar

1. **Capture filter (BPF) ile display filter'ı karıştırmak** → kutuda hata, ya da yanlış BPF ile kanıtı geri-dönüşsüz elemek. (En temel, en sık.)
2. **Yanlış yakalama noktası** → switch körlüğü; promiscuous açık ama boş/eksik pcap; "tehdidi göremiyorum" → bütünü değil, noktanın gösterdiğini görürsün.
3. **Promiscuous mode'u sihir sanmak** → "açtım, herkesi görürüm" yanılgısı; switch'te işe yaramaz.
4. **Paket düşürmeyi (drop) görmezden gelmek** → eksik veriyle "temiz" raporu üretmek; dumpcap/BPF kullanmamak.
5. **`ip.addr != X` ile bir IP'yi dışladığını sanmak** → mantık hatası; her paketi geçirir. Doğrusu `!(ip.addr == X)`.
6. **TLS'i "kapalı kutu" sayıp metadata'yı atlamak** → SNI/sertifika/JA3/beacon'ı kaçırmak; "şifreli, bakacak bir şey yok" demek.
7. **Kısa pcap'e bakıp "temiz" demek** → low-and-slow beacon/exfil'i kaçırmak; süre bilincinin olmaması.
8. **GeoIP'yi kimlik atfı sanmak** → VPN/proxy/bulut yanıltmasını görmezden gelip yanlış ülkeye/aktöre işaret etmek.
9. **pcap'i sanitize etmeden paylaşmak** → içindeki tüm parola/çerez/kişisel veri/topolojiyi sızdırmak (kendi OPSEC ihlali).
10. **Name resolution açıkken pasif analiz** → Wireshark'ın canlı DNS sorgusuyla pcap'i kirletmesi + ağa varlığını duyurması.
11. **Devasa pcap'i GUI'de açmaya zorlamak** → donma/çökme; önce tshark/editcap ile dilimlememek.
12. **Saat senkronu olmayan sensörlerle korelasyon** → olay zaman çizgisinin kayması; UTC kullanmamak.
13. **ARP spoof ile "yakalama" yapmak** → pasif analizi aktif (ve çoğu yerde yasa dışı) saldırıya çevirmek; TAP/SPAN dururken.
14. **Tek bir IOC'ye (IP/domain) kilitlenip davranışı kaçırmak** → "bilinen kötü" listesi boş çıkınca "temiz" sanmak; bilinmeyen kötüyü davranıştan aramamak.

---

<a id="13"></a>
## 13. 🏰 Kanije Kalesi ile Birlikte Kullanım & Öz-Farkındalık

Bu repo (Kanije Kalesi), fiziksel tehdit anında cihazı koruyan bir muhafızdır ve **Telegram Bot API üzerinden** komut alıp bildirim gönderir. Wireshark/ağ analizi burada iki yönden devreye girer: (1) **Kanije'nin kendi trafiğini anlamak/doğrulamak**, (2) **bir analistin gözünden Kanije'nin nasıl göründüğünü bilmek (öz-farkındalık)** ve (3) **şüpheli giden bağlantıları avlamak.**

### 13.1 Kanije'nin ağ ayak izi (ne üretir)
Kanije'nin ağ-ilgili davranışı, kaynak koddan teyitli olarak şudur:

| Davranış | Trafik deseni | Display filter |
|---|---|---|
| **Bağlantı kontrolü (heartbeat)** | Her **~5 sn**'de `api.telegram.org:443`'e TCP bağlantı denemesi (açıp kapatır) | `ip.dst == <telegram_ip> && tcp.flags.syn==1` |
| **Komut bekleme (long-poll)** | `getUpdates`'e ~**30 sn** uzun süren HTTPS isteği (sürekli açık bağlantı) | `tls && ip.addr == <telegram_ip>` |
| **Genel IP sorgusu** | Ara sıra (≤5 dk'da bir) `api.ipify.org`'a HTTPS | `tls.handshake.extensions_server_name contains "ipify"` |
| **Bildirim/medya gönderimi** | Olayda `api.telegram.org`'a HTTPS POST (foto/ses/dosya = büyük upload) | `tls && ip.dst == <telegram_ip>` |
| **Yerel ağ izleme** | `netsh`/`iwgetid`/`ip route` (yerel komutlar, **ağa paket üretmez**) | — (yerel, telde görünmez) |

> 🔥 **ÖZ-FARKINDALIK — Kanije bir analiste nasıl görünür:** Kanije'nin trafik deseni, **bir C2 (Command & Control) beacon'ına tehlikeli derecede benzer:**
> - **Düzenli aralıklı** giden bağlantılar (5 sn heartbeat + 30 sn long-poll) → **beacon deseni** (Bölüm 8.1). Bir IO Graph'ta Kanije, eşit aralıklı sivri uçlar üretir — tıpkı bir C2 gibi.
> - **Sabit bir dış uç noktaya** (`api.telegram.org`) sürekli konuşma → "ele geçmiş makine patronunu arıyor" gibi okunabilir.
> - **TLS içinde gizli komut/veri** (Telegram Bot API şifreli) → analist içeriği göremez (Bölüm 9), yalnızca "şu makine Telegram'a beacon atıyor + ara sıra büyük upload (foto/ses) yapıyor" görür.
> - **Olayda büyük giden upload** (kamera fotoğrafı, mikrofon kaydı, `/dosya al` ile dosya) → **exfiltration imzasına** birebir benzer (Tx ≫ Rx).
>
> Yani Kanije, **savunma amaçlı** olmasına rağmen, ağ trafiği imzası bakımından bir **Telegram-tabanlı RAT/C2'den ayırt edilmesi zordur.** Bu, hem bir tasarım gerçeği hem bir uyarıdır: Kanije'yi yetkili olmadığın bir kurumsal ağda çalıştırırsan, ağ savunması (NDR/SOC) onu **kötücül beacon** olarak işaretleyebilir. (Telegram-tabanlı C2'ler — örn. çeşitli "ToxicEye"/Telegram-RAT aileleri — gerçek tehditlerdir, bu yüzden SOC'lar `api.telegram.org` beacon'ına duyarlıdır.)

### 13.2 Kanije'yi doğrulamak için Wireshark kullanımı
Kendi makinende Kanije çalışırken, gerçekten **yalnızca beklenen hedeflerle** konuştuğunu doğrulayabilirsin:
```
# Kanije'nin TÜM dış bağlantılarını gör (yerel ağ hariç)
ip.dst != 192.168.0.0/16 && ip.dst != 10.0.0.0/8 && tcp.flags.syn == 1

# Beklenen: yalnızca api.telegram.org + api.ipify.org IP'leri.
# BAŞKA bir hedef görünüyorsa → ya yanlış yapılandırma ya tehlike → araştır.
```
> 💡 **Güven ama doğrula:** Kanije CGo-free, üçüncü-parti SDK kullanmayan, yalnızca stdlib `net/http` ile Telegram'a konuşan bir araçtır. Wireshark ile **bunu kendin teyit edebilirsin**: SNI'lerine bak (`tls.handshake.extensions_server_name`), yalnızca `api.telegram.org` ve `api.ipify.org` görmelisin. Beklenmedik bir SNI/IP = inceleme konusu.

### 13.3 Şüpheli giden bağlantı avı (Kanije'nin koruduğu makinede)
Kanije'nin koruduğu makinenin **başka** süreçlerinin kötücül olup olmadığını ağdan avlamak:
1. Makinenin tüm dış SYN'lerini çıkar (`ip.dst not in {yerel} && tcp.flags.syn==1`).
2. Hedefleri **Conversations** ile sırala; tanımadığın dış IP/domain'leri **tehdit istihbaratıyla** (VirusTotal, AbuseIPDB) karşılaştır.
3. **Beacon** ara (IO Graph) — Kanije'nin kendi beacon'ı dışında **ikinci bir düzenli beacon** varsa → başka bir şey "ev arıyor" → şüphe.
4. **DNS** anomalisi (NXDOMAIN patlaması, uzun subdomain) → DGA/tünel.
5. **Büyük upload** (Tx ≫ Rx) → Kanije'nin medya gönderimi dışında bir exfil varsa → kanıt.

### Karşılaştırma — kim neyi yapar

| Senaryo | Wireshark/Ağ analizi rolü | Kanije Kalesi rolü |
|---|---|---|
| Makinede beacon/C2 şüphesi | Beacon'ı IO Graph ile tespit, IOC çıkar | — (Kanije savunma muhafızı) |
| Kanije'nin trafiğini doğrulama | SNI/IP teyidi, "yalnızca Telegram mı?" | Kendi davranışını şeffaf üretir |
| Veri sızması avı | Tx≫Rx, DNS/ICMP tünel tespiti | `/imha` ile veriyi yok et (sızmadan önce) |
| Fiziksel olay anı | (pasif gözlem) | `/koruma`, `/panik`, `/kilit tam` |
| Adli sonrası analiz | pcap'ten olay zaman çizgisi çıkar | Olay log'u + bildirim geçmişi |

> 🧠 **Felsefe örtüşmesi & dürüstlük:** Wireshark "ağda gerçekte ne aktığını" gösteren tarafsız bir mikroskoptur — Kanije dahil. Bu rehber, Kanije'nin kendi trafiğinin bir C2 beacon'ına benzediğini **saklamaz, açıkça söyler**: çünkü iyi bir analist (ve iyi bir savunma aracı kullanıcısı) **kendi araçlarının ağda nasıl göründüğünü bilmek zorundadır.** "Benim aracım iyi niyetli" demek, ağ savunmasının onu masum sayacağı anlamına gelmez. Kanije'yi **yalnızca kendi cihazında/yetkili olduğun ortamda** çalıştır — tıpkı Wireshark'ı yalnızca kendi ağında çalıştırman gibi.

---

<a id="14"></a>
## 14. ✅ Hızlı Referans — En Yararlı Filtreler & Analiz Kontrol Listesi

### En yararlı display filter'lar (Wireshark dili)

| Amaç | Display Filter |
|---|---|
| Belirli IP (her yön) | `ip.addr == 10.0.0.5` |
| Bir IP'yi tamamen dışla | `!(ip.addr == 10.0.0.5)` |
| Yerel ağ DIŞINA giden | `ip.src == 192.168.0.0/16 && !(ip.dst == 192.168.0.0/16)` |
| Bağlantı kurma (SYN) | `tcp.flags.syn == 1 && tcp.flags.ack == 0` |
| Reddedilen/kesilen (RST) | `tcp.flags.reset == 1` |
| TCP anomalileri | `tcp.analysis.flags` |
| Tüm DNS sorguları | `dns.flags.response == 0` |
| NXDOMAIN (DGA göstergesi) | `dns.flags.rcode == 3` |
| DNS TXT (tünel) | `dns.qry.type == 16` |
| Uzun DNS adı (tünel) | `dns.qry.name.len > 50` |
| HTTP POST (exfil/C2) | `http.request.method == "POST"` |
| HTTP host | `http.host contains "evil"` |
| Şüpheli User-Agent | `http.user_agent contains "curl" \|\| http.user_agent contains "python"` |
| TLS Client Hello (JA3) | `tls.handshake.type == 1` |
| TLS SNI | `tls.handshake.extensions_server_name` |
| 443-dışı TLS (gizli C2) | `tls.handshake.type == 1 && !(tcp.port == 443)` |
| Büyük ICMP (tünel) | `icmp && data.len > 64` |
| SMB (lateral) | `smb2 \|\| smb` |
| Uzak .exe (PsExec) | `smb2.filename contains ".exe"` |
| Ham metin arama | `frame contains "password"` |
| Belirli TCP oturumu | `tcp.stream == N` |

### En yararlı BPF capture filter'lar (yakalamadan önce)

| Amaç | BPF |
|---|---|
| Tek host | `host 10.0.0.5` |
| Alt ağ | `net 192.168.1.0/24` |
| DNS | `port 53` |
| HTTPS | `tcp port 443` |
| Gürültü at (SSH+RDP) | `not port 22 and not port 3389` |
| Bir kaynaktan dışarı | `src host 10.0.0.5 and not dst net 10.0.0.0/8` |

### Analiz kontrol listesi (bir pcap geldiğinde)

**Bağlam & körlük**
- [ ] **Bu nereden yakalandı?** (uç nokta / gateway / SPAN / TAP) → **neyi göremiyorum?**
- [ ] Yakalama noktası körlüğü değerlendirildi (lateral mi, sadece dış mı görünüyor)
- [ ] **Paket düşürme (dropped) = 0** mı? (Statistics → Capture)
- [ ] Zaman dilimi **UTC**'ye alındı, sensör saati senkron mu?

**Büyük resim**
- [ ] **Protocol Hierarchy** → anormal protokol var mı?
- [ ] **Conversations** → en hacimli / en uzun süren çiftler kim?
- [ ] **Endpoints** → en aktif dış hedefler + (varsa) GeoIP
- [ ] **Expert Info** → retransmission/RST patlaması/malformed?

**Tehdit avı**
- [ ] Yerel ağ **DIŞINA** giden bağlantılar tarandı
- [ ] **Beacon** arandı (IO Graph, düzenli aralık + jitter)
- [ ] **DNS** anomalisi (NXDOMAIN patlaması, uzun/rastgele subdomain, çok TXT)
- [ ] **HTTP** anomalisi (tuhaf UA, anormal POST, şüpheli URI)
- [ ] **TLS** metadata (SNI, sertifika, JA3 — şifreli C2 göstergesi)
- [ ] **Exfil** (Tx ≫ Rx, ICMP/DNS tünel, low-and-slow)
- [ ] Şüpheli **dosyalar** Export Objects ile çıkarıldı + hash'lendi
- [ ] **IOC**'ler (IP/domain/SNI/UA/JA3/hash) çıkarıldı (tshark) → TI ile karşılaştırıldı

**OPSEC (kendi güvenliğin)**
- [ ] **Name resolution** (canlı DNS) kapalı (pasif/gizli analizde)
- [ ] pcap **şifreli** saklanıyor (VeraCrypt)
- [ ] Paylaşılacaksa **sanitize** edildi (IP/MAC/payload maskeleme)
- [ ] İş bitince **güvenli silme** planı var
- [ ] Yakalama **yalnızca yetkili olduğun ağda** yapıldı (Bölüm 15)

---

<a id="15"></a>
## 15. ⚖️ Hukuki Sınır — Yetkisiz Dinleme SUÇTUR

> 🔴 **Bu bölüm tavsiye değil, sınırdır. Atlama.**

- **Yetkisiz trafik yakalama = suç.** Başkasının ağ trafiğini izni olmadan yakalamak/dinlemek **çoğu yargı bölgesinde açıkça suçtur:**
  - **Türkiye:** TCK m.243 (bilişim sistemine hukuka aykırı erişim), m.244 (sistemi engelleme/bozma/verileri ele geçirme), ayrıca **haberleşmenin gizliliğini ihlal** (m.132) ve **kişisel verileri hukuka aykırı ele geçirme** (m.136).
  - **ABD:** Wiretap Act (18 U.S.C. § 2511), Computer Fraud and Abuse Act (CFAA).
  - **AB:** ePrivacy Direktifi + GDPR (kişisel veri içeren trafik).
- **"Teknik olarak mümkün" ≠ "yasal".** Ortak Wi-Fi'de, kafede, iş yerinde başkasının trafiğine "merak ettim" diye bakmak **suçtur** — Wireshark'ın bunu kolaylaştırması seni aklamaz.
- **ARP spoofing/MITM** (Bölüm 4) **aktif müdahaledir** ve yetkisiz yapıldığında ek suç oluşturur (sistemi/iletişimi bozma).
- **Nerede yasaldır:**
  - **Kendi ağın, kendi cihazların** (senin trafiğin, senin ağın).
  - **Yazılı izinli** pentest/kırmızı takım (kapsam ve yetki belgesiyle, **scope dışına çıkma**).
  - **Yetkili SOC/IR** görevi (kurumun kendi ağında, iş sözleşmesi/politikası kapsamında — yine de **çalışan mahremiyeti** ve **iç politika** sınırlarına uy).
  - **İzole eğitim laboratuvarı** (kendi VM'lerin, kendi ürettiğin trafik).
- **pcap = kişisel veri taşıyıcısı.** Yakaladığın trafik parola, mesaj, kişisel veri içerir. **Yetkili olsan bile** bunu saklama/işleme **veri koruma yükümlülüğüne** tabidir: şifreli sakla, gereğinden uzun tutma, paylaşmadan sanitize et, görev bitince güvenli sil.
- **Atıf (attribution) dikkatli yapılır.** GeoIP/IP, kesin kimlik değildir (Bölüm 11.14). "Şu ülke/aktör yaptı" gibi iddialar hukuki/itibari sonuç doğurur — **kanıt zinciri ve teyit** olmadan atıf yapma.
- **Yasal sorumluluk:** Bu rehber **meşru savunma, olay müdahalesi (IR), tehdit istihbaratı, ağ sorun giderme ve eğitim** içindir. Bulunduğun yargı bölgesinin yasalarına ve içinde bulunduğun kurumun politikasına uy. **Şüphedeysen, yakalama.**

---

> 🏰 **Kapanış:** Ağ analizi bir buton değil, bir **disiplindir.** En güçlü filtre bile, yanlış yakalama noktasında boş pcap üretir; en keskin göz bile, TLS'in kapattığı içeriği göremez — ama **metadata'yı, beacon'ı, JA3'ü, hacmi** okuyabilen analist karanlıkta bile avlar. Wireshark sana telin üzerindeki **tartışılmaz gerçeği** verir; **doğru yerden bakmak, neyi göremediğini bilmek ve yasanın içinde kalmak senin işin.** Bir C2 beacon'ı ile Kanije'nin heartbeat'ini ayırt etmek, "şifreli, demek temiz" tuzağına düşmemek, ve topladığın pcap'in kendisinin bir sır kasası olduğunu unutmamak — ustalık budur. Kanije Kalesi de tam burada: kendi trafiğini bile şeffafça gösteren, ama yalnızca **senin** cihazında nöbet tutan muhafız.
>
> *Bu doküman Kanije Kalesi güvenlik rehberleri koleksiyonunun parçasıdır. İlgili: `TAILS_ANONIMLIK_REHBERI.md` (Tor/exit node/trafik korelasyonu), `VERACRYPT_USTALIK_REHBERI.md` (pcap'i şifreli sakla), `QUBES_OS_BOLMELEME_REHBERI.md` (ağ izolasyonu/sys-net), `GNUPG_GPG_USTALIK_REHBERI.md` (metadata sızıntısı), `LINUX_HARDENING_KALE.md`, `WINDOWS11_HARDENING_KALE.md`.*
