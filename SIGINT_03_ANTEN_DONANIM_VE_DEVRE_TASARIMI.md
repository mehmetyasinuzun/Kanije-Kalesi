# 📶 SIGINT EL KİTABI — BÖLÜM 3
## Antenler, RF Donanımı ve Devre Tasarımı — Sinyali Yakalamak ve Güçlendirmek

> **Bu bölüm neyi öğretir:** Bir önceki bölümde "sinyal nedir, SDR neyi nasıl örnekler" konusunu işledik. Şimdi **zincirin en başına** iniyoruz: antene. Çünkü dünyanın en pahalı alıcısı bile, ucundaki anten kötüyse **sağır**dır. Bu bölüm seni "rastgele tel takan" seviyesinden çıkarıp, **neden o boyda, neden o empedansta, neden o filtre** sorularını cevaplayan; gerektiğinde **kendi antenini ve hatta kendi RF kartını tasarlayacak** zihinsel modele taşır. Sonunda "insanlar bu bilgiyi nereden biliyor?" sorusunun da somut cevabını vereceğiz: kaynaklar, datasheet'ler, topluluk.

> ⚖️ **EN BAŞTA NET YASAL SINIR — ATLAMA:**
> Bu bölümün büyük kısmı **alıcı (RX) taraftır**: anten yapmak, dinlemek, LNA ile zayıf sinyali yükseltmek, filtre eklemek **çoğu yerde serbesttir** ve bu el kitabının ruhu budur — *gözlemlemek, anlamak*.
> **Verici (TX) taraf bambaşkadır.** "Sinyalin watt'ını güçlendirmek", yani **yayın yapmak / güç yükseltmek**, neredeyse her ülkede **lisans gerektirir**. Yetkisiz veya yüksek güçle yayın:
> - **Suçtur** (telsiz/telekomünikasyon kanunları; Türkiye'de BTK düzenler),
> - **Tehlikelidir** (RF yanığı, göz/doku hasarı, yangın),
> - **Zarar verir** (acil servis, havacılık, GPS, GSM bantlarına **girişim** — insan hayatını riske atar),
> - **Donanımı yakar** (yanlış empedans/SWR ile PA patlar).
>
> Bu bölümün TX/PA/EIRP/watt kısmı **PRENSİBİ ANLAMAK** içindir — *operatörlük talimatı değil*. "Nasıl daha çok yayarım" değil, "watt artırmanın mühendislik mantığı nedir, neden tehlikeli ve neden regüle edilir" sorusunu cevaplar. **Uygulamak istiyorsan: amatör telsiz sınavına gir, lisans al, yasal güç/bant sınırında çalış.** Lisanssız iş için yalnızca **RX** ve **kapalı kutu / dummy load** ile ölçüm yap.

---

## 📑 İÇİNDEKİLER

1. [Neden Anten Her Şeyin Başıdır](#1)
2. [Anten Temel Teorisi — Rezonans, Boy, Empedans](#2)
3. [VSWR, Kazanç, Işıma Deseni, Polarizasyon](#3)
4. [Anten Türleri — Hangisi, Ne Zaman, Nasıl Çalışır](#4)
5. [🔥 LNA — Düşük Gürültülü Yükselteç ve Friis Gürültü Formülü](#5)
6. [Filtreler — Bant-Dışı Gürültüyü Kesmek](#6)
7. [Koaksiyel Kablo ve Konnektörler](#7)
8. [Empedans Eşleme, Balun/Unun, Ölçüm (NanoVNA)](#8)
9. [Upconverter / Downconverter — Bandı Taşımak](#9)
10. [🔥 GÜÇ ve WATT — TX, PA, EIRP'nin Mantığı (Sınırla)](#10)
11. [🔥 GENİŞLETME MODÜLLERİ ve KENDİ PCB'Nİ ÇİZMEK](#11)
12. ["Nereden Biliyorlar?" — Bilginin Gerçek Kaynakları](#12)
13. [🧪 Alıştırmalar (Yasal, Ev Şartlarında)](#13)
14. [Özet, Kontrol Listesi ve Çapraz Referans](#14)

---

<a id="1"></a>
## 1. 🧭 Neden Anten Her Şeyin Başıdır

Bir SIGINT (sinyal istihbaratı / sinyal gözlemi) zinciri kabaca şöyledir:

```
 [ANTEN] → [Filtre] → [LNA] → [Filtre] → [SDR alıcı] → [Bilgisayar/DSP]
    ▲          ▲         ▲         ▲           ▲              ▲
 sinyali   istenmeyen  zayıf    görüntü     örnekler      anlam
 yakalar   bantları    sinyali  frekansı    (ADC)         çıkarır
           kes         yükselt  bastır
```

Bu zincirde değiştirilmesi en ucuz ama etkisi en büyük halka **antendir**. Sebep basit ve acımasız bir doğa kanunu: **kaybettiğin sinyali sonradan geri kazanamazsın.** Antende yakalayamadığın enerjiyi, sonraki hiçbir yükselteç "yaratamaz" — yükselteç sadece *zaten var olanı* (sinyal + gürültü) birlikte büyütür. Bu yüzden:

> 🧠 **Altın kural:** Önce anten + konum, sonra düşük gürültülü ön-uç (LNA + filtre), **en son** yazılım. Çoğu yeni başlayan tersini yapar; pahalı SDR alır, ucuna 10 cm'lik fabrika çubuğu takar ve "sinyal zayıf" diye şikâyet eder. Asıl sorun **antende ve konumdadır** (yükseklik, engel, parazit kaynaklarından uzaklık).

İki kavramı baştan ayıralım:
- **Kazanç (gain):** Antenin/devrenin sinyali *yönlendirme veya yükseltme* yeteneği. Pasif antende kazanç, enerjiyi bir yöne **toplamaktan** gelir (yoktan üretmez). Aktif devrede (LNA) gerçek elektriksel yükseltmedir.
- **Gürültü (noise):** Her şeyin üstüne binen istenmeyen enerji (termal, atmosferik, insan yapımı). Amaç sinyali değil, **sinyal/gürültü oranını (SNR)** iyileştirmektir.

Bütün bu bölüm aslında tek bir cümlenin açılımıdır: **"Doğru frekansta, doğru yönden, mümkün olduğunca çok sinyali, mümkün olduğunca az gürültüyle alıcıya taşı."**

---

<a id="2"></a>
## 2. 📐 Anten Temel Teorisi — Rezonans, Boy, Empedans

### 2.1 Neden boy dalga boyuna bağlı? (Rezonans)

Bir anten, içinde elektronların **ileri-geri salındığı** bir iletkendir. Tıpkı belirli uzunluktaki bir gitar teli yalnızca belirli notada en güçlü titreştiği gibi, belirli uzunluktaki bir iletken de belirli bir frekansta **rezonansa** girer — yani o frekanstaki dalgaya en verimli "kaplanır", en çok enerji alışverişi yapar.

Rezonans, iletken boyu **dalga boyunun** (λ, lambda) anlamlı bir kesriyle (λ/2, λ/4 gibi) eşleştiğinde oluşur. Çünkü o boyda, gelen dalganın elektrik alanı iletken üzerinde **duran dalga** kurar; akım ve gerilim düzgün, verimli bir desen oturtur.

**Dalga boyu formülü (ışık hızından):**

```
       c            300
λ (m) = ──    ≈    ─────        ( c ≈ 3×10⁸ m/s ;  f MHz cinsinden )
       f          f(MHz)
```

Örnek: 100 MHz (FM bandı ortası) → λ ≈ 300/100 = **3 metre**. 433 MHz → λ ≈ 0,69 m = 69 cm. 1090 MHz (uçak ADS-B) → λ ≈ 0,275 m = 27,5 cm.

### 2.2 Anten boyu pratik formülleri (🔥 doğru ezberle)

Teorik λ/2 ile **gerçekte kesilen tel** arasında küçük bir fark vardır. İletkenin çapı ve uç etkileri (end effect) yüzünden gerçek rezonans boyu, serbest uzaydaki λ/2'den **biraz kısadır** (tipik **velocity/kısalma faktörü ≈ 0,95**). Pratikte amatör dünyada şu formüller kullanılır:

**Yarım dalga dipol (toplam tel boyu):**

```
                143          ( metre cinsinden, ince tel için pratik;
L(yarım dalga) ≈ ─────         teorik 150/f'in ~0,95 katı = end-effect )
                f(MHz)

  (İnç-feet seven kaynaklarda:  L(ft) ≈ 468 / f(MHz)  — aynı şeyin emperyali)
```

**Çeyrek dalga monopol / dipolün tek bacağı:**

```
                71,5
L(çeyrek dalga) ≈ ──────      ( = yarım dalga / 2 )
                 f(MHz)
```

> ⚠️ **Teyit et:** `143/f` ve `468/f` yaygın *pratik* dipol formülleridir (end-effect dahil, ~0,95 faktörü). Saf teorik serbest-uzay yarım dalgası `150/f`'dir (`= (300/f)/2`). Kalın iletken, yalıtım veya yakın metal varsa gerçek boy değişir; bu yüzden son ayarı **NanoVNA ile rezonansı ölçüp tel ucundan kırparak** yaparsın (bkz. §8). Formül başlangıç noktasıdır, kutsal değil.

**Türetme (neden bunlar?):**
- Serbest uzayda yarım dalga: `λ/2 = (300/f)/2 = 150/f` metre.
- End-effect (uç kapasitansı) elektriksel boyu fiziksel boydan uzun gösterir → fiziksel teli **%5 kısaltırsın**: `150/f × 0,95 ≈ 143/f`.
- Çeyrek dalga onun yarısı: `≈ 71,5/f`.

**Hızlı tablo (pratik dipol / monopol boyları):**

| Frekans | λ (≈300/f) | Yarım dalga dipol (≈143/f) | Çeyrek dalga (≈71,5/f) | Tipik kullanım |
|---|---|---|---|---|
| 27 MHz | ~11,1 m | ~5,30 m | ~2,65 m | CB telsiz |
| 88–108 MHz | ~3,0 m | ~1,40 m | ~0,70 m | FM yayın (RX) |
| 137 MHz | ~2,19 m | ~1,04 m | ~0,52 m | NOAA/uydu (RX) |
| 433 MHz | ~0,69 m | ~33 cm | ~16,5 cm | ISM, uzaktan kumanda |
| 868/915 MHz | ~0,34 m | ~16,5 cm | ~8,2 cm | LoRa/ISM |
| 1090 MHz | ~0,275 m | ~13 cm | ~6,5 cm | ADS-B uçak (RX) |
| 1575 MHz | ~0,19 m | ~9,5 cm | ~4,8 cm | GPS L1 (RX) |

> 🔥 **Püf:** ADS-B (1090 MHz) için "tam boy" çeyrek dalga sadece **6,5 cm**'dir. Bir SMA konnektöre 6,5 cm bakır tel lehimleyip altına dört adet ~6,9 cm "radyal" eklersen, ev yapımı ground-plane anten **mağaza antenini sollar**. Maliyet: bir kahve parasından az. İşte SIGINT donanımının güzelliği bu — **bilgi, paradan değerlidir**.

### 2.3 Empedans — neden 50 Ω?

Her anten, beslendiği noktada bir **empedans** (Z, ohm cinsinden, AC direnci) gösterir. Bu, "anten kaynaktan ne kadar kolay enerji alır/verir" ölçüsüdür. Rezonanstaki ideal yarım dalga dipolün serbest uzay empedansı teorik olarak **≈ 73 Ω**, çeyrek dalga monopolünki (ideal toprak düzlemiyle) **≈ 36 Ω** civarıdır.

**Maksimum güç transferi**, kaynak empedansı ile yük empedansı **eşleştiğinde** olur (empedans uyumu). Sektör, RF için bir uzlaşı değeri seçmiştir: **50 Ω.**

Neden tam 50? Tarihsel ve fiziksel bir dengedir (koaksiyel kabloda):
- **En düşük kayıp** için optimum empedans (hava dielektrikli kabloda) ≈ **77 Ω**.
- **En yüksek güç taşıma** kapasitesi için optimum ≈ **30 Ω**.
- İkisinin **geometrik ortalaması** (√(77×30) ≈ 48) ≈ **50 Ω** → hem makul kayıp hem makul güç. Pratik bir altın orta.
- (Yayın/TV dünyasında **75 Ω** standardı vardır — düşük kayıp önceliği. Bu yüzden TV anten kablosu 75 Ω; RTL-SDR'ların çoğu 50 Ω ama bazıları 75 Ω F-konnektörlüdür. **Karıştırma**, küçük uyumsuzluk küçük kayıp demek.)

> 🧠 **Sezgi:** 50 Ω "kutsal" değil, **anlaşılmış ortak dil**. Herkes 50 Ω yaparsa anten, kablo, LNA, SDR sorunsuz birbirine takılır. Empedans uyumsuzluğu → güç geri yansır → kayıp + dalgalanma (bir sonraki başlık: VSWR).

---

<a id="3"></a>
## 3. 📊 VSWR, Kazanç, Işıma Deseni, Polarizasyon

### 3.1 VSWR / SWR — uyumun karnesi

**SWR (Standing Wave Ratio)** / **VSWR (Voltage SWR)**: Anten + kablo sistemiyle kaynak ne kadar iyi *eşleşmiş*? Uyumsuzlukta, gönderilen/alınan dalganın bir kısmı **geri yansır**; ileri ve geri dalga üst üste binip **duran dalga** kurar. VSWR, bu duran dalganın tepe/dip oranıdır.

**İlişkiler:**

```
Yansıma katsayısı:   Γ = (Z_yük − Z_hat) / (Z_yük + Z_hat)

                       1 + |Γ|
VSWR  =  ─────────────────────
                       1 − |Γ|

Geri dönen güç oranı = |Γ|²    (örn. |Γ|=0,2 → %4 güç yansır)
```

| VSWR | |Γ| | Yansıyan güç | Yorum |
|---|---|---|---|
| 1,0:1 | 0,00 | %0 | Kusursuz (ideal, gerçekte nadir) |
| 1,2:1 | 0,09 | ~%0,8 | Mükemmel |
| **1,5:1** | **0,20** | **~%4** | **İyi — pratik hedef sınırı** |
| 2,0:1 | 0,33 | ~%11 | Kabul edilebilir (çoğu durumda OK) |
| 3,0:1 | 0,50 | %25 | Kötü — ciddi kayıp, TX'te tehlikeli |
| ∞ (açık/kısa) | 1,0 | %100 | Tüm güç geri (anten yok/kopuk) |

> 🔥 **Neden 1,5:1 altı "iyi"?** 1,5'te yalnızca **~%4** güç yansır (RX için ihmal edilebilir). RX'te VSWR'a *aşırı* takılma — birkaç dB'lik uyumsuzluk gürültü tabanının altında kalabilir. Ama **TX'te VSWR ölümcüldür:** yansıyan güç PA'ya geri döner, ısı + gerilim tepeleri oluşturur, **PA'yı yakar.** İşte TX'in RX'ten neden çok daha hassas/tehlikeli olduğunun bir nedeni daha.

### 3.2 Anten kazancı — dBi, dBd

**Kazanç (gain)**: Antenin enerjiyi **belirli yönlere toplama** becerisi. Pasif anten enerji üretmez; sadece her yöne eşit saçmak yerine bir yöne **odaklar** (el feneri reflektörü gibi). Bir yöne daha çok → diğer yöne daha az.

- **dBi:** İzotropik (her yöne eşit saçan ideal nokta) kaynağa göre kazanç. Referans budur.
- **dBd:** Yarım dalga dipole göre kazanç. Dönüşüm: **dBi = dBd + 2,15.** (Dipolün izotropa göre kazancı ~2,15 dBi'dir.)

**dB sezgisi (🔥 ezberle):**

```
+3 dB  ≈ 2×  güç        −3 dB ≈ ½ güç
+6 dB  ≈ 4×  güç        +10 dB = 10× güç
+20 dB = 100× güç       +30 dB = 1000× güç

Formül:   dB = 10 · log₁₀( P_çıkış / P_giriş )      (GÜÇ oranı için)
          dB = 20 · log₁₀( V_çıkış / V_giriş )      (GERİLİM/genlik oranı için)
```

> ⚠️ **Bedava öğle yemeği yok:** "20 dBi anten" çok yüksek kazançtır ama **çok dar huzme** demektir — yalnızca tam o yöne baktığında işe yarar, gerisini görmez. Tarama/keşif için **düşük kazanç + geniş açı** (diskon, dipol); bilinen tek hedef için **yüksek kazanç + dar açı** (Yagi, çanak). Kazanç her zaman "daha iyi" değildir; **göreve bağlıdır.**

### 3.3 Işıma deseni (radiation pattern)

Antenin hangi yöne ne kadar duyarlı olduğunun "haritası". (Antenler **karşılıklıdır/reciprocal**: aynı anten, vericiyken hangi yöne yayarsa, alıcıyken o yönden o kadar iyi alır.)

```
 Dipol (yandan kesit):  "donut" / simit       Yagi (üstten):  ileriye sivri huzme
        ↑                                            ▲
     ___|___                                    ◄────█────►   arka: zayıf (front/back)
    /   |   \   yatay düzlemde                       │
   (    •    )  her yöne (omni)                  ileriye odaklı, yanlar/arka bastırılmış
    \___|___/                                        ▼
        ↓ (uçlardan ışımaz)
```

- **Omnidirectional (her yöne):** Dikey dipol/monopol → ufukta 360° dinler, dikeyde sınırlı. Tarama/keşif için ideal.
- **Directional (yönlü):** Yagi, log-periyodik, çanak → bir yönde yüksek kazanç, diğer yönler bastırılmış. Zayıf/uzak hedef veya yön bulma (DF) için.

### 3.4 Polarizasyon

Elektrik alanının salınım **düzlemi**. Verici ve alıcı anten polarizasyonu **uyuşmazsa** ciddi kayıp olur (dik polarizasyonlar arası teorik kayıp **çok büyüktür**, ~20 dB+).

| Polarizasyon | Tipik kullanım |
|---|---|
| **Dikey (vertical)** | Çoğu telsiz, mobil, ISM, yer-tabanlı omni — anten yere dik |
| **Yatay (horizontal)** | FM/TV yayını (bazı bölgelerde), bazı HF |
| **Dairesel (circular: RHCP/LHCP)** | **Uydu** (LEO geçişlerinde dönen geometriyle uyum) → QFH, turnike anten |

> 🔥 **Püf:** NOAA/Meteor uydularını dikey çubukla almaya çalışırsan sürekli "fading" (sönümlenme) yaşarsın çünkü uydu yörüngede döndükçe polarizasyon kayar. **Dairesel polarize QFH/turnike** anten bu sorunu çözer — uydu dinlemenin "neden bu garip anten?" cevabı budur.

---

<a id="4"></a>
## 4. 📡 Anten Türleri — Hangisi, Ne Zaman, Nasıl Çalışır

| Anten | Yönlülük | Bant genişliği | Tipik kazanç | En iyi olduğu iş |
|---|---|---|---|---|
| **Dipol (λ/2)** | Omni (donut) | Dar (rezonant) | ~2,15 dBi | Tek banda basit, referans |
| **Monopol (λ/4 + radyal)** | Omni | Dar | ~2–5 dBi | Mobil/taban omni, ADS-B |
| **Diskon** | Omni | **Çok geniş (10:1)** | ~0–2 dBi | **Geniş bant tarama/keşif** |
| **Yagi-Uda** | Yönlü | Dar–orta | 7–20 dBi | Uzak/zayıf hedef, DF, TV |
| **Log-periyodik (LPDA)** | Yönlü | **Geniş** | 6–11 dBi | Geniş bant + yön bir arada |
| **Helical (eksenel)** | Yönlü, dairesel | Orta | 10–15 dBi | Uydu/dairesel yer linki |
| **Biquad** | Yarı-yönlü | Orta | ~10–12 dBi | 2,4 GHz Wi-Fi/ISM, ucuz DIY |
| **PCB / yonga anten** | Değişken | Dar–orta | <0–2 dBi | Küçük cihaz, IoT, dahili |
| **Magnetik loop** | Sekiz (null'lı) | Çok dar (ayarlı) | düşük (verimli RX) | **HF, gürültülü ortam, küçük alan** |
| **QFH / turnike** | Yukarı yarım küre, dairesel | Orta | ~3–5 dBi | **LEO uydu (NOAA/Meteor)** |

### 4.1 Dipol — atası
İki adet λ/4 bacak (toplam ~λ/2), ortadan beslenir. Basit, öngörülebilir, **referans** anten. Donut deseni: yatayda omni, uçlardan kör.

```
   λ/4              λ/4
 ◄──────►   ╪   ◄──────►
 ════════ [besle] ════════
              │
           (koaks)
```

### 4.2 Monopol + ground plane (λ/4)
Dipolün bir bacağını "toprak düzlemi" (radyaller veya araba tavanı) ile değiştirir → yarısı kadar tel, toprak düzlemi diğer yarısı gibi davranır ("imge" prensibi). **ADS-B / mobil** için klasik.

```
        │ λ/4 dikey eleman
        │
   ─────┴─────  ← besleme noktası (SMA)
   ╱   ╱│╲   ╲
   radyaller (~λ/4, ufka doğru ~45° eğimli tipik)
```

### 4.3 Diskon — geniş bant tarama kralı 🔥
Bir **disk** + bir **koni** (disc + cone). Bu geometri, klasik rezonansa değil **frekanstan bağımsız** bir geçiş empedansına yaklaşır → **10:1 gibi devasa bant** (örn. 100–1000 MHz tek antenle). Kazancı mütevazıdır ama "ne çıkarsa bakalım" keşif/tarama için **en pratik tek anten**dir.

```
   ____ disk ____
  ╱──────────────╲     ← üstte düz disk
        │ besle
       ╱ ╲
      ╱   ╲   koni (ışınsal teller de olur)
     ╱     ╲
    ╱       ╲
```
> Diskonun "geniş bant" sırrı: rezonant tek boy yerine, **sürekli değişen kesit** her frekansta "yeterince uygun" bir empedans sunar. Mükemmel değil ama her yerde *çalışır*.

### 4.4 Yagi-Uda — yönlü, yüksek kazanç
Bir **sürülen eleman** (driven, dipol) + arkada bir **reflektör** + önde bir-çok **yönlendirici** (director). Parazitik elemanlar enerjiyi öne **odaklar**. Eleman sayısı ↑ → kazanç ↑, huzme ↓ (TV çatı anteni bunun ta kendisi).

```
 reflektör  sürülen   yönlendiriciler →→→
   │          ╪          │   │   │
   │     ════[besle]════  │   │   │      ileriye huzme ►►►
   │          │          │   │   │
 ──┴──────────┴──────────┴───┴───┴──  (boom)
```

### 4.5 Log-periyodik (LPDA) — geniş bant + yön
Boyları logaritmik oranla küçülen **bir dizi dipol**. Her frekansta o frekansa uyan elemanlar aktif → **geniş bantta** yönlü kazanç. "Geniş bant Yagi" gibi düşün (laboratuvar/EMC anteni klasiği).

### 4.6 Helical (eksenel mod)
Toprak düzlemi üstünde **yay/spiral** sarım. Eksenel modda **dairesel polarize**, ileriye iyi kazanç → **uydu** yer linki, dairesel polarizasyon gerektiren işler.

### 4.7 Biquad
İki kare ilmek + reflektör plaka. **2,4 GHz** (Wi-Fi/ISM) için ucuz, etkili, **DIY dostu** yarı-yönlü anten. Bakır tel + CD kutusu kapağıyla yapılır.

### 4.8 PCB / yonga (chip) anten
Telefonlar, IoT modülleri: anten **devre kartının bakırına** çizilir (meander/IFA deseni) ya da minik **seramik yonga** olarak lehimlenir. Yer kazandırır, kazancı düşüktür; tasarımı hassastır (üretici **datasheet'inde referans layout** verir — bkz. §11–12).

### 4.9 Magnetik loop (HF)
Küçük, ayarlı bir **iletken halka** + ayar kondansatörü. Fiziksel olarak küçük ama **HF**'de etkili; **elektrik alan gürültüsünü reddeder** → apartman/gürültülü şehir ortamında HF dinlemenin kurtarıcısı. Çok dar bantlıdır, kondansatörle frekansa **akort** edersin. Sekiz (figure-8) deseni → null'larıyla **parazit kaynağını bastırma**/yön bulma imkânı.

### 4.10 QFH / Turnike (uydu)
**QFH (Quadrifilar Helix)** ve **turnstile (turnike)**: yukarı yarım küreye **dairesel polarize** desen → ufuktan tepeye geçen **LEO uydularını** (NOAA APT, Meteor LRPT) dengeli alır, dönen polarizasyona dayanıklı. "Neden bu tuhaf burgu/artı şekilli anten?" sorusunun cevabı: **uydu geometrisi.**

---

<a id="5"></a>
## 5. 🔥 LNA — Düşük Gürültülü Yükselteç ve Friis Gürültü Formülü

### 5.1 LNA nedir, ne işe yarar?
**LNA (Low-Noise Amplifier)**, zayıf sinyali **kendi üstüne çok az gürültü ekleyerek** yükselten aktif bir devredir. Amaç sinyali "duyulur" yapmak **değil** sadece — asıl amaç, zincirin geri kalanının (özellikle kabloda kaybolan ve SDR'ın gürültülü ön-ucundan geçen) etkisini **maskelemek**.

### 5.2 🔥 Neden anten ucunda? (Friis gürültü formülü)

İşte bu bölümün en önemli mühendislik içgörüsü. Kademeli (kaskad) bir sistemin **toplam gürültü faktörü**, **Friis formülü** ile verilir:

```
                 F₂ − 1     F₃ − 1            Fₙ − 1
F_toplam = F₁ + ──────── + ──────── + … + ─────────────────
                  G₁         G₁·G₂        G₁·G₂·…·G_(n−1)
```

Burada:
- **F** = her katın **gürültü faktörü** (lineer oran; **NF = 10·log₁₀(F)** ile dB'ye çevrilir),
- **G** = her katın **kazancı** (lineer oran).

**Formülün anlattığı tek büyük gerçek:** İlk katın gürültüsü (F₁) **doğrudan**, bölünmeden toplama girer. Sonraki katların katkıları ise **önceki kazançlara bölünür**. Yani **ilk kat hem en düşük gürültülü olmalı hem yeterli kazanç vermeli** — ki kendinden sonraki her şeyin gürültüsünü "ezsin", önemsizleştirsin.

İşte bu yüzden LNA **antenin hemen ucuna** (kablodan ÖNCE) konur:
- Kablo, sinyali zayıflatan (kayıplı) pasif bir kattır; kaybı **L** ise gürültü faktörü de yaklaşık **L** kadardır.
- Eğer önce **kablo** (kayıp) sonra LNA gelirse: kablonun kaybı F₁ olarak doğrudan toplama girer → felaket.
- Eğer önce **LNA** (yüksek G, düşük F) sonra kablo gelirse: kablonun gürültüsü `G_LNA`'ya bölünür → **görünmez** olur.

```
KÖTÜ:  [ANTEN]→[uzun kablo, −6 dB]→[LNA]→[SDR]
        ilk kat 6 dB kayıp = sistem NF ≥ 6 dB'den başlar (mahvoldu)

İYİ:   [ANTEN]→[LNA, NF~1 dB, G~20 dB]→[uzun kablo, −6 dB]→[SDR]
        ilk kat NF~1 dB; kablo+SDR gürültüsü ÷100 (20 dB) → neredeyse yok olur
```

### 5.3 Sayısal örnek (kaskad — kâğıt üstünde çöz)

Diyelim:
- **LNA:** NF = 1 dB (F₁ ≈ 1,26), kazanç G₁ = 20 dB (= 100×).
- **Sonraki her şey (kablo + SDR birleşik):** NF = 10 dB (F₂ = 10).

**LNA antende (doğru sıra):**
```
F_top = F₁ + (F₂−1)/G₁ = 1,26 + (10−1)/100 = 1,26 + 0,09 = 1,35
NF_top = 10·log₁₀(1,35) ≈ 1,3 dB    ← muhteşem
```

**LNA yok (sadece kablo + SDR):**
```
F_top = 10  →  NF_top = 10 dB    ← 8,7 dB daha kötü
```

> 🔥 **Sonuç:** İyi yerleştirilmiş LNA, sistem gürültü figürünü **10 dB'den ~1,3 dB'ye** indirdi. Bu, gürültü tabanını ~8–9 dB düşürmek = daha önce **duyulamayan** zayıf sinyalleri **duyulur** yapmak demektir. **Antene en yakın aktif kat, sistemin kaderini belirler.**

### 5.4 ⚠️ Ne zaman LNA ZARAR verir? (aşırı sürme / IMD)
LNA her derde deva değildir — **yanlış kullanılırsa zincirini bozar:**
- **Aşırı kazanç / yakın güçlü sinyaller:** Tepenizdeki güçlü bir FM/TV/GSM vericisi, LNA'yı **doyurur (saturation)**. LNA lineerliğini kaybeder → **intermodülasyon (IMD)**: gerçekte olmayan **hayalet sinyaller** üretir, her yer "ızgara/birbiri içine geçmiş" görünür.
- **Gürültü tabanı zaten kabloda kaybolmuyorsa:** Kısa kablo + iyi SDR'da LNA gereksiz kazanç yığar, ADC'yi taşırır.
- **Çözüm:** LNA'dan **önce/sonra filtre** (bant-dışı güçlüleri kes — §6), kazancı abartma, gerekiyorsa **attenuator** (zayıflatıcı) ile dengeyi kur. "Daha çok kazanç = daha iyi" **yanlıştır**; doğru olan **doğru SNR ve lineerlik**.

> 🧠 **Mühendislik dengesi:** LNA = düşük gürültü için *yukarı*, ama lineerlik (büyük sinyale dayanım) için *aşağı* baskı. İkisi çakışır. Profesyonel ön-uç tasarımı bu iki uç arasında **filtreyle** denge kurmaktır.

---

<a id="6"></a>
## 6. 🔇 Filtreler — Bant-Dışı Gürültüyü Kesmek

Anten ve LNA "her şeyi" alır. Filtre, **istemediğin frekansları** alıcıya/LNA'ya ulaşmadan **bastırır**. Neden kritik?
- **Güçlü bant-dışı vericiler** (yerel FM 88–108 MHz, GSM, TV) ön-ucu doyurur, IMD yaratır (§5.4).
- **Görüntü frekansı (image):** Karıştırıcı (mixer) mimarisi, istenen frekansın "ayna" karşılığını da içeri alabilir → filtre bunu reddeder.
- **Geniş bant gürültü:** İlgilenmediğin spektrumdaki toplam enerji, ADC dinamik aralığını boşa harcar.

**Filtre türleri:**

| Tür | Ne yapar | Tipik SIGINT kullanımı |
|---|---|---|
| **LPF (alçak geçiren)** | Belli frekans **altını** geçirir | HF dinlerken üstteki VHF/FM'i kes |
| **HPF (yüksek geçiren)** | Belli frekans **üstünü** geçirir | Alttaki güçlü AM/FM'i kes |
| **BPF (bant geçiren)** | Sadece **belli bandı** geçirir | ADS-B (1090 MHz) only, GSM-band only |
| **Notch (çentik)** | Tek dar bandı **bastırır** | Yerel boğucu FM istasyonunu sil |
| **SAW** | Yüzey akustik dalga ile **keskin** BPF | Kompakt, dik etekli filtre (1090 SAW yaygın) |

```
BPF tepkisi (geçiş bandı + dik etekler):
 kazanç
   │      ┌──────────┐
   │     ╱            ╲          ← sadece bu bant geçer
   │    ╱   geçiş     ╲           dışı sert bastırılır
   │___╱     bandı     ╲_______
   └────────────────────────────► frekans
       f_alt        f_üst
```

> 🔥 **Klasik örnek 1 — ADS-B SAW filtre:** 1090 MHz uçak dinlerken, yerel GSM/LTE devasa güçtedir. Antenin hemen ardına **1090 MHz SAW BPF** koyarsan, LNA artık sadece uçak bandını "görür", doymaz → menzil ve görülen uçak sayısı katlanır.
> **Klasik örnek 2 — FM bant-stop (notch/HPF):** Şehirde 88–108 MHz FM o kadar güçlüdür ki her şeyi ezer. Geniş bant tararken bir **FM band-stop** filtre, bütün spektrumu "temizler".

**Sıralama mantığı:** Çoğu iyi ön-uç `ANTEN → BPF → LNA → BPF/SDR` şeklindedir: LNA'dan **önce** bir filtre onu büyük sinyallerden korur; sonrasında bir filtre LNA'nın ürettiği geniş gürültüyü/görüntüyü temizler.

---

<a id="7"></a>
## 7. 🔌 Koaksiyel Kablo ve Konnektörler

Kablo "sadece tel" değildir; **kayıplı bir bileşendir** ve kayıp **frekansla artar.** Yanlış/uzun/ucuz kablo, antenin tüm kazancını yer.

**Kablo kayıpları (yaklaşık, 100 m başına dB — teyit için üretici datasheet'i):**

| Kablo | Çap | ~100 MHz | ~1000 MHz | Not |
|---|---|---|---|---|
| **RG-174** | İnce (~2,8 mm) | ~8–12 dB | ~30–40 dB | Çok kayıplı; sadece **çok kısa** patch |
| **RG-58** | Orta (~5 mm) | ~5–6 dB | ~20+ dB | Yaygın ama UHF'de zayıf |
| **RG-6** | (75 Ω) | ~2–3 dB | ~7–9 dB | TV kablosu; uzun RX için fena değil |
| **LMR-240** | ~6 mm | ~2,5 dB | ~8 dB | İyi denge |
| **LMR-400** | Kalın (~10 mm) | ~1,2 dB | ~4 dB | Düşük kayıp, uzun hat / çatı anteni |

> 🔥 **Püf (frekansla kayıp):** Aynı RG-58, 100 MHz'de ~5 dB/100m iken 1 GHz'de ~20+ dB/100m olur — **frekans arttıkça kayıp artar.** Bu yüzden ADS-B/GHz işlerinde kablo **kısa ve kalın (LMR-400)** olmalı; ya da LNA'yı (§5) antene koyup kabloyu **LNA'dan sonra** çekersin (kayıp artık önemsiz). **Altın kural: ya kabloyu kısalt, ya LNA'yı yukarı taşı.**

**Konnektörler:**

| Konnektör | Tipik | Not |
|---|---|---|
| **SMA** | SDR donanımları (RTL-SDR, HackRF) | Küçük, vidalı; **RP-SMA** (Wi-Fi) ile karıştırma — pim ters! |
| **BNC** | Lab, osiloskop, eski telsiz | Bayonet (çeyrek tur), pratik |
| **N** | Yüksek güç / yüksek frekans / dış mekân | Sağlam, su geçirmez seçenekli, düşük kayıp |
| **F** | TV / 75 Ω | Ucuz, anten/TV; bazı RTL-SDR'larda |

> ⚠️ **SMA vs RP-SMA tuzağı:** Görünüş aynı ama iç pim ters yerleşimlidir; uyumsuz çift sinyali geçirmez. Wi-Fi anteni (RP-SMA) SDR'a (SMA) doğrudan oturmaz — adaptör gerekir. Çok yaygın bir DIY hatası.

**Genel ilke:** Kablo **kısa, kaliteli, doğru empedansta, sağlam konnektörlü** olsun. Her ek/adaptör küçük bir kayıp + yansıma noktasıdır. RF'te "az ek, kısa hat" daima kazanır.

---

<a id="8"></a>
## 8. 🎯 Empedans Eşleme, Balun/Unun, Ölçüm (NanoVNA)

### 8.1 Neden eşleme?
Anten (örn. 73 Ω dipol) ile hat/alıcı (50 Ω) **tam uyuşmaz** → küçük yansıma. Çoğu RX'te 73→50 ihmal edilir, ama daha büyük uyumsuzluklarda (loop, çok-bantlı anten, uzun besleme) **eşleme ağı** (matching network: L/C devreleri, stub'lar, transformatör) ile empedansı 50 Ω'a "çevirirsin" → minimum yansıma, maksimum transfer.

### 8.2 Balun ve Unun
- **Balun = BALanced↔UNbalanced.** Dipol gibi **dengeli (balanced)** bir anteni, koaks gibi **dengesiz (unbalanced)** bir hatta bağlarken kullanılır. Yoksa, koaksın dış örgüsü **anten gibi davranır** (akım dışta akar) → desen bozulur, gürültü/parazit girer. Basit bir **ferrit choke (1:1 balun)** çoğu dipolde bu sorunu çözer.
- **Unun = UNbalanced↔UNbalanced**, ayrıca **empedans dönüştürme** (örn. 9:1 unun, uzun-tel/EFHW antende ~450 Ω'u ~50 Ω'a yaklaştırır).

```
 Dipol (dengeli)            Koaks (dengesiz)
   ═══╪═══   ──[ 1:1 balun / ferrit choke ]── ║coax║──► alıcı
       │                                      (dış örgüde
   (besleme)                                   asalak akımı bastırır)
```

### 8.3 Ölçüm — SWR metre ve NanoVNA (kavram)
- **SWR metre:** Hatta ileri/geri gücü ölçer → VSWR'ı söyler. Antenin "uyumlu mu" sorusunu cevaplar.
- **NanoVNA (Vektör Network Analizörü):** Ucuz, taşınabilir mucize. Bir frekans aralığını **tarar**, her noktada **VSWR, empedans (R+jX), kazanç/kayıp (S11/S21)** verir. Anteni keserken/akort ederken **rezonansı gözle görürsün** (VSWR'ın dibe vurduğu frekans). "Tahmin etme, **ölç**" aracı.

> 🔥 **Pratik akış:** Formülle (§2.2) teli yaklaşık kes → NanoVNA bağla → VSWR eğrisine bak → rezonans **istediğin frekanstan yüksekteyse tel kısadır → uzat**; **alçaktaysa tel uzundur → kırp.** Birkaç iterasyonda 1,2:1 VSWR'a oturtursun. **NanoVNA olmadan** bunu yapmanın yolu: formüle güven, biraz **uzun kes** (kırpmak eklemekten kolay), mümkünse bir alıcıda S-metre/SNR'ı izleyerek en iyi boyu **deneysel** bul.

---

<a id="9"></a>
## 9. 🔁 Upconverter / Downconverter — Bandı Taşımak

Bazı alıcılar belirli bir frekans aralığını **göremez** (örn. tipik RTL-SDR ~24 MHz altını doğrudan alamaz → **HF kör**). Çözüm: ilgilenilen bandı, alıcının **gördüğü** bir banda **kaydırmak** — bir **karıştırıcı (mixer) + yerel osilatör (LO)** ile.

- **Upconverter (yukarı çevirici):** HF'yi (0–30 MHz) yukarı taşır. Örn. **+125 MHz LO** ile: 7 MHz'lik sinyal → 132 MHz'te görünür; SDR'ı 132 MHz'e ayarlar, aslında 7 MHz dinlersin. (Ham-It-Up gibi modüller bunu yapar.) **HF dinlemeyi RTL-SDR'a açan** tipik genişletme.
- **Downconverter (aşağı çevirici):** Çok yüksek bir bandı (örn. SHF/uydu LNB'si) alıcının erişebileceği daha düşük bir banda indirir. Uydu **LNB**'si tam olarak budur: 10–12 GHz'i ~1 GHz IF'e indirir.

```
 İlgilenilen bant         Mixer            SDR'ın gördüğü
   7 MHz  ───────►  ⊗  ◄── LO 125 MHz  ───►  132 MHz
                    (toplam/fark frekansı üretir)
```

> 🧠 **Kavram:** Mixer iki frekansı çarpar → **toplam ve fark** çıkar (f_LO ± f_RF). İstediğin tarafı filtreyle seçersin. Tüm süper-heterodin alıcıların kalbi budur — "frekansı taşıma" hilesi.

---

<a id="10"></a>
## 10. 🔥 GÜÇ ve WATT — TX, PA, EIRP'nin Mantığı (Sınırla)

> ⚖️ **TEKRAR — bu başlık PRENSİP içindir, operatörlük değil.** Yetkisiz/yüksek güç **yayın** suçtur, tehlikelidir, donanım yakar (bkz. en üstteki uyarı). Burada *neden* ve *nasıl çalıştığını* anlatıyoruz ki "watt'ı güçlendirmenin mantığı" sorusu cevaplansın ve **neden regüle edildiğini** kavrayasın. Uygulamak = lisans + yasal sınır.

### 10.1 TX gücü ne demek?
**Çıkış gücü (W)**, vericinin antene **verdiği** RF enerji oranıdır. Ama menzili belirleyen sadece bu değildir — **antenin kazancı** ve **kablo kaybı** da hesaba girer. Bunların birleşimi **EIRP/ERP** ile ölçülür.

### 10.2 🔥 EIRP / ERP — "etkili yayılan güç"

```
EIRP (dBm) = P_TX(dBm) − Kayıp_kablo/konnektör(dB) + Kazanç_anten(dBi)

ERP (dBd referanslı) = EIRP − 2,15 dB     (çünkü dBi = dBd + 2,15)

Watt↔dBm:   P(dBm) = 10·log₁₀( P(mW) )    →   30 dBm = 1 W ,  0 dBm = 1 mW
```

**Örnek (anlama amaçlı):**
- TX = 1 W = **30 dBm**, kablo kaybı **−2 dB**, anten kazancı **+12 dBi**:
  ```
  EIRP = 30 − 2 + 12 = 40 dBm = 10 W EIRP
  ```
  Yani 1 W'lık verici, 12 dBi yönlü antenle **o yönde sanki 10 W yayıyormuş gibi** etki yapar.

> 🔥 **Watt artırmanın MANTIĞI:** Menzili artırmanın iki yolu var:
> 1. **Vericinin gücünü (W) artırmak** → PA (güç yükselteci) ile. Pahalı, ısı/verim sorunlu, **regüle/sınırlı**, doğrusallık zorlaşır.
> 2. **Antenin kazancını/yüksekliğini artırmak** (pasif!) → çoğu zaman **çok daha akıllıca**: +3 dBi anten, gücü 2'ye katlamakla aynı etki — ama **bedava, ısısız, çoğu yerde serbest, tek yöne odaklı**. Bu yüzden profesyoneller **önce anteni ve yüksekliği** zorlar, watt'ı en son artırır.

### 10.3 PA (Power Amplifier) ne yapar, devre prensibi?
PA, düşük güçlü RF sinyalini (örn. mixer/modülatör çıkışı, birkaç mW) **kazanç katları** halinde yükselterek antene yetecek güce taşır. "Watt'ı büyütmek" = sinyali kontrollü şekilde **çoğaltmaktır** (yoktan üretmek değil; **DC besleme gücünü RF'e dönüştürmek**). Temel gerilim: PA, güç kaynağından çektiği DC enerjiyi, transistörle "şekillendirip" RF çıkışına aktarır; verim = (RF çıkış)/(DC giriş).

**Amplifikatör sınıfları — lineerlik ↔ verim dengesi:**

| Sınıf | İletim açısı | Lineerlik | Verim | Tipik kullanım |
|---|---|---|---|---|
| **A** | 360° (hep açık) | **En iyi** | **Düşük (~%25–30)** | LNA, hassas lineer RX/küçük sinyal |
| **B / AB** | ~180° / arası | İyi (AB pratik denge) | Orta (~%50–60) | Lineer modülasyon (SSB), ses |
| **C** | <180° | **Kötü (lineer değil)** | Yüksek (~%70+) | Sabit-zarflı (FM/CW) — modülasyonu bozmaz |
| **D / E / F** | Anahtarlamalı | (özel) | **Çok yüksek (>%80–90)** | Modern verimli RF/ses, anahtar-modlu |

> 🧠 **Anahtar fikir:** **Lineerlik ile verim çatışır.** Genliği bilgi taşıyan modülasyonlarda (QAM, SSB) **lineerlik şart** → A/AB sınıfı → düşük verim → çok ısı. Genliği sabit modülasyonlarda (FM, CW) **C/D/E/F** ile yüksek verim "bedava" gelir çünkü bozulma bilgiyi etkilemez. PA tasarımı = "hangi modülasyon, ne kadar doğrusallık, ne kadar ısıyı göze alırım" dengesidir.

> ⚠️ **Neden tehlikeli (donanım+insan):** Yanlış empedans/yüksek VSWR (§3.1) → yansıyan güç PA'da ısı + gerilim tepesi → **PA patlar**. Yetersiz soğutmada **yangın**. Yeterli güçte anten yakınında **RF yanığı / doku ısınması / göz hasarı**. Ve en önemlisi: yayın **başkalarının bandına girer** (GPS, havacılık, GSM, acil servis) → bu yüzden **devletler watt'ı ve bandı sıkı regüle eder** ve lisanssız yayın **suçtur**.

---

<a id="11"></a>
## 11. 🔥 GENİŞLETME MODÜLLERİ ve KENDİ PCB'Nİ ÇİZMEK

Kullanıcının asıl merakı: "İnsanlar **kendi RF kartlarını** nasıl tasarlıyor, bu kadar bilgiyi nereden alıyorlar?" Önce **modül kavramı**, sonra **kendi PCB akışı**, sonra (§12) **bilginin kaynağı.**

### 11.1 Genişletme modülleri — hazır yapı taşları
RF dünyası **modüler**dir. İnsanlar genelde sıfırdan çizmez; **hazır breakout/modül** zincirler:
- **LNA modülü** (örn. geniş bant ya da banda özel),
- **Filtre modülü** (SAW BPF, FM bant-stop),
- **Upconverter** (HF için),
- **Bias-tee** (koaks üstünden LNA'ya DC besleme: sinyal + güç aynı kabloda),
- **Attenuator / splitter / switch** modülleri,
- **SDR'ın kendisi** (RTL-SDR, HackRF, Airspy…) bir genişletme platformudur.

Bunlar **SMA ile uç uca** takılır → lehimsiz, denemesi kolay bir "RF lego". Çoğu ciddi sistem böyle **modül zincirinden** doğar; PCB çizmek **son adımdır** (modülleri tek karta birleştirmek ya da özel bir şey gerektiğinde).

### 11.2 Bias-tee — neden var?
LNA'yı **antenin ucunda** (çatıda) çalıştırmak istiyorsun (§5.2) ama oraya ayrı güç kablosu çekmek zor. **Bias-tee**, DC gücü **koaks kablonun içinden** LNA'ya gönderir (RF'i ayırıp DC'yi bindirerek). Pek çok RTL-SDR'ın "bias-tee" özelliği tam bunun için: **tek kabloyla** hem sinyali al hem uzaktaki LNA'yı besle.

### 11.3 🔥 Kendi RF PCB'ni çizmenin TEMELLERİ

RF'te PCB "sadece bağlantı" değil, **devrenin parçasıdır** — bakır izler belirli frekanslarda **hat/empedans** gibi davranır. Temel ilkeler:

1. **Kontrollü empedans iz (controlled impedance / 50 Ω trace):** RF taşıyan iz, **tam 50 Ω** olmalı. İzin genişliği, altındaki dielektrik (FR-4) kalınlığı ve bakır kalınlığı **birlikte** empedansı belirler. Bunu **mikroşerit (microstrip)** hesabıyla bulursun (üretici/araç "impedance calculator" verir; PCB üreticisi yığın bilgisini — stackup — söyler). Yanlış genişlik → yanlış empedans → yansıma/kayıp.

```
 Microstrip kesiti:
   ▭ W (iz genişliği)  ← bunu hesaplarsın
  ════════════════     üst bakır (sinyal izi)
  ░░░░░░░░░░░░░░░░  h  ← dielektrik (FR-4), kalınlık önemli
  ████████████████     alt: KESİNTİSİZ GROUND düzlemi
   Z₀, W/h ve εr (dielektrik sabiti) ile belirlenir
```

2. **Kesintisiz ground plane (toprak düzlemi):** RF izin **hemen altında** bütün, delik-deşik olmayan bir bakır toprak şart. Dönüş akımının yolu budur; bozulursa empedans bozulur, anten gibi ışır, gürültü girer.

3. **Via stitching (toprak dikişi):** Üst ve alt toprakları **çok sayıda via** ile birbirine "diker" → toprak gerçekten "tek potansiyel" olur, RF kaçaklarını/kavite rezonanslarını bastırır. RF kartlarındaki o "via tarlası" bunun içindir.

4. **RF/dijital ayrımı:** Gürültülü dijital (saat, USB, anahtarlamalı regülatör) bölümü, hassas RF bölümünden **fiziksel olarak ayır**; toprakları tek noktadan birleştir, besleme hatlarını **decoupling kondansatör + ferrit boncuk** ile filtrele. Yoksa dijital gürültü RF'e sızar.

5. **Komponent yerleşimi:** Sinyal zinciri **kısa ve düz** aksın (anten girişi → filtre → LNA → çıkış). Kıvrım, uzun iz, gereksiz via = kayıp + yansıma. Yüksek frekansta **her milimetre** sayar.

6. **Empedans eşleme alanları:** Çip (LNA/mixer) giriş-çıkışına **matching** için boş pad'ler (L/C için) bırakılır — datasheet'in önerdiği değerlerle doldurulur.

> 🔥 **Neden 50 Ω trace bu kadar önemli?** Çünkü tüm zincir (anten, kablo, konnektör, çip) 50 Ω. PCB izi de 50 Ω olmazsa, kartın **içinde** bir empedans sıçraması olur → yansıma, VSWR bozulması, kayıp, bazen osilasyon. PCB'yi "görünmez bir 50 Ω koaks" gibi tasarlarsın.

### 11.4 KiCad / EasyEDA ile akış (kavram)
- **Şematik** çiz (çiplerin datasheet'teki referans bağlantısı) → **footprint** ata → **PCB layout**'a geç.
- **Stackup/empedans:** Üreticinin (JLCPCB/PCBWay) stackup'ını gir, **impedance calculator** ile 50 Ω iz genişliğini bul, RF izi ona göre çiz.
- **Ground plane + via stitching** dök, RF/dijital ayır, kısa yolla.
- **DRC** (tasarım kuralı kontrolü) → **Gerber** dosyaları üret.
- **Üretim:** Gerber'ı **JLCPCB / PCBWay**'e yükle → birkaç dolara birkaç günde kart gelir (isteğe bağlı **dizgi/assembly** de yaptırılır). Bu erişilebilirlik, "evde RF kartı" çağını açtı.

> 🧠 **Gerçeklik:** İlk RF kartın muhtemelen çalışmaz ya da beklenenden kötü olur — bu **normaldir.** Profesyoneller bile **revizyon** (rev A, B, C) yapar: ölç (NanoVNA/spektrum), anla, düzelt, tekrar bas. RF tasarımı **iteratif** bir zanaattır; "tek seferde mükemmel" beklentisi yanlıştır.

---

<a id="12"></a>
## 12. 🧠 "Nereden Biliyorlar?" — Bilginin Gerçek Kaynakları

Kullanıcının çekirdek sorusu. Cevap **gizli/sezgisel değil**; herkese açık, izlenebilir kaynaklardan gelir. İnsanlar şuralardan öğreniyor:

1. **Komponent DATASHEET'leri (en temel kaynak):** Her RF çipinin (LNA, mixer, PA, SAW) üreticisi, ayrıntılı **veri sayfası** yayımlar: elektriksel değerler, gürültü figürü, kazanç, **önerilen uygulama devresi (typical application circuit)**, **referans PCB layout**, empedans eşleme değerleri. Mühendis çoğu zaman datasheet'teki **referans tasarımı kopyalar**. "Nereden biliyor?" → çoğu doğrudan **datasheet'ten**.

2. **Üretici APPLICATION NOTE'ları:** Texas Instruments, Analog Devices, NXP, Qorvo, Skyworks, Mini-Circuits gibi firmalar **uygulama notları** (AN-xxx) yayımlar: "Bu çiple 50 Ω eşleme nasıl yapılır", "LNA layout tuzakları", "PA termal tasarım". Bunlar **bedava, derin** ve pratik. **Referans tasarım kartları** (eval board) ve onların **şematik+gerber'ları** da paylaşılır → birebir incelersin.

3. **Amatör radyo (ham) topluluğu ve ARRL Handbook:** Onyıllardır biriken **açık bilgi**. **ARRL Handbook** ve **ARRL Antenna Book**, anten/RF tasarımının "kutsal kitapları" sayılır; formüller, tablolar, yapım projeleri içerir. Ham radyo kültürü "yap, ölç, paylaş" üzerine kuruludur.

4. **RF tasarım kitapları (akademik/pratik):** Pozar'ın *Microwave Engineering* (microstrip, eşleme, S-parametreleri), anten için Balanis/Kraus, pratik için ham yayınları. Üniversitelerin RF/mikrodalga dersleri bu temeli verir; ama **çoğu hobici kitabı + datasheet + deneyle** öğrenir.

5. **Açık donanım projeleri (GitHub vb.):** RTL-SDR uyumlu LNA'lar, filtreler, upconverter'lar, hatta tüm SDR'lar **açık kaynak** paylaşılır — **şematik + gerber + parça listesi**. İnsanlar bunları okuyup **çalışan tasarımdan** öğrenir, kendine uyarlar. (Açık kaynak donanım, RF eğitiminin en hızlı yoludur.)

6. **Topluluk / forumlar:** **RTL-SDR.com** (blog + rehber hazinesi), **Reddit r/RTLSDR**, **r/amateurradio**, **r/RFElectronics**, RF mühendislik forumları (ör. RFdesign tartışmaları), QRZ, eevblog. Burada gerçek sorunlar, ölçümler, hatalar paylaşılır → **kolektif akıl**.

> 🔥 **Demistifikasyon:** "Bu kadar bilgiyi nereden biliyorlar?" sorusunun dürüst cevabı: **kimse her şeyi kafasından bilmiyor.** Mühendis bile her tasarımda **datasheet + application note + referans tasarım + ölçüm + forum** karışımına başvurur. Fark, *bilgiyi nereden bulacağını ve nasıl doğrulayacağını* bilmektir. Bu bölüm sana o **haritayı** verdi: önce datasheet/app-note → referans tasarımı anla → KiCad'de uygula → JLCPCB'de bas → NanoVNA ile ölç → forumda doğrula → revize et. Süreç **tekrarlanabilir ve öğrenilebilir** — sihir değil, **disiplin + kaynak okuryazarlığı**.

---

<a id="13"></a>
## 13. 🧪 Alıştırmalar (Yasal, Ev Şartlarında)

> Hepsi **RX/hesap/anlama** odaklıdır — yayın yok, lisans gerekmez. Sadece kâğıt, tel, makas ve (varsa) bir RTL-SDR.

**Alıştırma 1 — Belirli frekans için dipol/monopol kes.**
Seçtiğin bir RX hedefi al (örn. **137,5 MHz** NOAA, ya da **1090 MHz** ADS-B).
- λ = 300/f, yarım dalga dipol ≈ 143/f, çeyrek dalga ≈ 71,5/f hesapla.
- 137,5 MHz için: λ ≈ 2,18 m; dipol ≈ **1,04 m** (her bacak ~52 cm); çeyrek dalga ≈ **52 cm**.
- 1090 MHz için: çeyrek dalga ≈ **6,56 cm** → 4 radyalli ground-plane yap.
- **Biraz uzun kes**, sonra (NanoVNA varsa) rezonansı ölçüp kırparak ayarla; yoksa SDR'da SNR'ı izleyerek en iyi boyu deneysel bul.

**Alıştırma 2 — Diskon neden geniş bant? Açıkla.**
Kendi cümlenle yaz: Rezonant tek-boy anten **dar** banttır çünkü yalnızca tek λ kesrinde uyar. Diskon, **sürekli değişen kesiti** (disk→koni geçişi) sayesinde geniş bir frekans aralığında "yeterince uygun" empedans sunar → 10:1 bant. Bedeli: mütevazı kazanç. (Tarama için neden ideal olduğunu bağla.)

**Alıştırma 3 — LNA'nın gürültü iyileştirmesini hesapla (kaskad/Friis).**
Verilenler: LNA NF = 0,8 dB (F₁ ≈ 1,20), G₁ = 18 dB (≈ 63×); arkadaki kablo+SDR birleşik NF = 9 dB (F₂ ≈ 7,94).
- **LNA antende:** F_top = 1,20 + (7,94−1)/63 = 1,20 + 0,110 = **1,31** → NF ≈ **1,17 dB**.
- **LNA yok:** NF = **9 dB**.
- Kazanım: ~**7,8 dB** daha düşük gürültü tabanı → kabaca menzil/duyarlık sıçraması. (Friis'in "ilk kat belirler" dersini kendi cümlenle özetle.)

**Alıştırma 4 — NanoVNA olmadan SWR'ı düşün.**
Cihazsız: VSWR'ın *ne anlattığını* yaz (yansıma oranı; 1,5:1 → ~%4 yansır). Sonra "rezonansı nasıl deneysel bulurum?" sorusuna cevap ver: teli **biraz uzun** kes, bir alıcıda S-metre/SNR'ı izle, **kırptıkça** sinyalin en güçlendiği boyu yakala; ya da bilinen güçlü bir işareti referans alıp en iyi alımı veren boyu bul. ("Ölçemiyorsan, **kontrollü deneyle yaklaş**" ilkesini içselleştir.)

**Alıştırma 5 (kavram) — Watt mı, anten mi?**
1 W TX + 3 dBi anten yerine **1 W TX + 9 dBi yönlü anten** kullanınca EIRP nasıl değişir? (+6 dBi = +4× → o yönde 1 W yerine ~4 W EIRP, **gücü hiç artırmadan**.) "Neden profesyoneller önce anten/yükseklik der?" sorusunu bu hesapla bağla — **ama yayının lisans/yasa gerektirdiğini** not düş.

---

<a id="14"></a>
## 14. ✅ Özet, Kontrol Listesi ve Çapraz Referans

### Özün özü
- **Anten = zincirin kaderi.** Önce anten + konum + yükseklik, sonra LNA/filtre, en son yazılım.
- **Boy dalga boyuna bağlı** (rezonans): yarım dalga dipol ≈ **143/f(MHz) m**, çeyrek dalga ≈ **71,5/f**. λ ≈ **300/f(MHz) m**.
- **50 Ω** = ortak dil; uyum **VSWR** ile ölçülür, **<1,5:1 iyi**.
- **LNA antene en yakın** (Friis): ilk kat sistem gürültüsünü belirler; ama **aşırı sürme = IMD/zarar.**
- **Filtre** bant-dışı güçlüleri keser (FM/GSM/görüntü frekansı); ön-ucu kurtarır.
- **Kablo kayıplıdır ve frekansla artar** → kısa+kalın, ya da LNA'yı yukarı taşı.
- **Watt artırmanın mantığı:** PA ile güç ya da (daha akıllıca) **anten kazancı/yükseklik**; EIRP = P − kablo + anten. **Ama yayın = lisans + yasa; yetkisiz/yüksek güç SUÇ + tehlike.**
- **Kendi PCB:** 50 Ω kontrollü iz + kesintisiz ground + via stitching + RF/dijital ayrımı; KiCad→Gerber→JLCPCB; **ölç-anla-revize** iteratif zanaat.
- **Bilgi nereden:** datasheet → application note → referans tasarım → ARRL/kitap → açık donanım → forum. **Sihir değil, kaynak okuryazarlığı + ölçüm + iterasyon.**

### Hızlı kontrol listesi (RX ön-uç kurarken)
- [ ] Anten **hedef banda** uygun mu (tarama→diskon/dipol; tek hedef→Yagi; uydu→QFH)?
- [ ] Boy formülle hesaplandı, mümkünse **NanoVNA ile akort** edildi (VSWR <1,5:1)?
- [ ] **LNA antene yakın** mı (kablodan önce), kazanç **abartılmadı** mı?
- [ ] Güçlü yerel kaynaklar için **filtre** (FM band-stop / banda özel BPF) var mı?
- [ ] Kablo **kısa/kaliteli/doğru empedans**, konnektör (SMA≠RP-SMA) doğru mu?
- [ ] Çatıdaki LNA için **bias-tee** ile tek kablodan besleme?
- [ ] (Uydu/dengeli anten) **balun/dairesel polarizasyon** doğru mu?

### ⚖️ Yasal hatırlatma (kapanış)
**Alıcı taraf** (anten, LNA, filtre, kablo, ölçüm, dinleme) bu el kitabının serbest ve teşvik edilen alanıdır. **Verici taraf** (güç yükseltme, yayın, PA, EIRP'yi sahada uygulamak) **lisans + yasal güç/bant sınırı** ister; yetkisiz veya yüksek güç yayın **suçtur, tehlikelidir (RF yanığı, yangın, donanım hasarı) ve başkalarının kritik bantlarına (havacılık/GPS/GSM/acil) girişimle insan hayatını riske atar.** Bu bölümün TX/PA/watt kısmı **yalnızca mühendislik prensibini anlamak** içindir.

---

> 📶 **Kapanış:** Anten görünmez bir el gibidir — uzaydaki dalgayı tutup tele indirir. Onu **boyuyla** akort eder, **empedansıyla** sisteme bağlar, **LNA'yla** ilk dokunuşu temiz tutar, **filtreyle** gürültüyü süzer, **kabloyla** kayıpsız taşırsın. "Watt'ı büyütmek" çoğu zaman yanlış sorudur; doğru soru **"sinyali nasıl daha temiz yakalarım"**dır — ve cevabın çoğu, watt değil, **geometri, empedans ve düşük gürültüdür.** Kendi kartını çizmek de erişilemez bir sihir değil; **datasheet okumak, 50 Ω'a saygı duymak ve ölçüp revize etmektir.**
>
> *Bu doküman Kanije Kalesi SIGINT El Kitabı serisinin **3. Bölümü**dür.*
> *← Önceki: `SIGINT_02_*.md` (Sinyal Teorisi, Modülasyon ve SDR ile Örnekleme) — burada işlenen "sinyal/SNR/spektrum/SDR" kavramlarının temeli.*
> *→ Sonraki: `SIGINT_04_*.md` (Yakalanan Sinyali Çözmek — Demodülasyon, Dekoderler ve Protokoller) — bu antenle topladığın RF'i **anlama** aşaması.*
> *İlgili repo rehberleri: `VERACRYPT_USTALIK_REHBERI.md` (yakalanan/üretilen veriyi koruma), Kanije Kalesi güvenlik komutları.*
