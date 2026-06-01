# SIGINT EL KİTABI — BÖLÜM 27: ANTEN DİZİLERİ, HÜZME ŞEKİLLENDİRME VE MASSIVE MIMO

## Tek Antenden Uzamsal İşlemeye — Diziler, Beamforming ve Çok-Kanallı Matematik

> Amaç: Bölüm 3 tek antenin fiziğini (rezonans, kazanç, yönlülük, polarizasyon) verdi; o bölüm tek bir elemanın ne yaptığını anlattı. Bu bölüm sahneye ikinci, üçüncü ve nihayetinde yüzlerce elemanı ekler ve sorar: birden çok anteni uzayda dizip her birinden gelen örneği karmaşık bir ağırlıkla çarpıp toplarsak ne kazanırız? Cevap, modern radyonun en güçlü fikirlerinden biridir — uzamsal işleme (spatial processing). Bir dizi, frekansta veya zamanda yapılamayan bir şeyi yapar: gelen dalgaların geliş yönüne göre ayrım. Bu, kazanç (anten kazancından bağımsız bir dizi kazancı), yönlülük (elektronik olarak yönlendirilebilen dar bir hüzme), girişim bastırma (istenmeyen bir kaynağa "sıfır" çevirme) ve yön kestirim (DoA — Bölüm 9'un alt-uzay algoritmalarının dizi tarafı) anlamına gelir. Bu bölüm reçete değil türetme verir; bir faz dizisi (phased array), bir KrakenSDR çıktısı ya da bir 5G Massive MIMO baz istasyonu gördüğünde her ağırlığın, her hüzmenin, her sıfırın arkasındaki denklemi zihninde canlandırabilmen için yazıldı.

> Yasal çerçeve: Bu bölüm uzamsal işleme matematiğine ve dizi mimarisine odaklanır; doğası gereği pasif analiz ve eğitim içeriğidir. Anlatılan tüm algoritmalar (dizi faktörü, beamforming ağırlığı, MVDR, MUSIC) bir alıcının kendi çok-kanallı örnekleri üzerinde çalışır; hiçbir iletim, karıştırma ya da yönlü güç yayma önermez. "Null steering" burada bir alıcının istenmeyen bir kaynağı bastırması (alış tarafı uzamsal filtre) olarak ele alınır; verici tarafı yönlü karıştırma değildir ve önerilmez. Alıştırmalar yalnızca kendi cihazların, kendi ürettiğin sinyaller ve yasal/açık yayın bandlarıyla (örn. bir FM istasyonunun yönünü ölçmek) sınırlıdır. Belirli bandların kaydı veya yön kestirimi bazı yargı alanlarında düzenlenebilir; kendi ülkenin ve sürümünün mevzuatını teyit et.

---

## İÇİNDEKİLER

1. [Neden Dizi: Kazanç, Yönlülük, Uzamsal Çeşitlilik, Girişim Bastırma](#1)
2. [Dizi Temelleri: Dizi Faktörü, Eleman Aralığı ve Faz İlerlemesi](#2)
3. [Izgara Lobu (Grating Lobe) ve Neden d ≤ λ/2](#3)
4. [Dizi Geometrileri: ULA, UCA, Düzlemsel, Seyrek Diziler](#4)
5. [Yönlendirme Vektörü (Steering Vector) ve Dizi Veri Modeli](#5)
6. [Hüzme Şekillendirme: Gecikme-Topla, Ağırlık Vektörü, Analog/Sayısal/Hibrit](#6)
7. [Uyarlanabilir Diziler: MVDR/LCMV, Null Steering, Wiener/LMS](#7)
8. [Yön Kestirim (DoA): Gecikme-Topla, Capon, MUSIC, ESPRIT](#8)
9. [MIMO ve Massive MIMO: Uzamsal Çoğullama ve Kanal Kapasitesi](#9)
10. [Çeşitlilik-Çoğullama Ödünleşimi ve Hüzme-Uzayı (Beamspace)](#10)
11. [5G Massive MIMO ve FR2/mmWave Hüzme Yönetimi](#11)
12. [Faz Tutarlılığı ve Kalibrasyon: Çok-Kanallı Alıcı, Saat/LO Paylaşımı](#12)
13. [Karşılıklı Kuplaj, Kenar Etkileri ve Gerçek-Dünya Bozulmaları](#13)
14. [Uygulama: Pasif Radar Dizisi, Uydu İzleme, Girişim Haritalama](#14)
15. [Alıştırmalar (Yasal, Lab)](#15)
16. [Hızlı Referans ve Diğer Bölümler](#16)

---

<a id="1"></a>
## 1. Neden Dizi: Kazanç, Yönlülük, Uzamsal Çeşitlilik, Girişim Bastırma

Bölüm 3'te tek bir antenin kazancının, enerjiyi bir yöne toplamaktan geldiğini gördük: Yagi öne odaklar, çanak dar bir hüzme verir, dipol donut deseni saçar. Bu kazanç mekaniktir — antenin geometrisi sabittir, yöneldiği yön fiziksel olarak nereye baktığına bağlıdır. Bir anten dizisi (antenna array) aynı yönlülüğü tamamen farklı bir yoldan, elektronik olarak elde eder: birden çok özdeş elemanı uzayda dizip her birinin sinyalini ayrı bir karmaşık ağırlıkla işleyip toplayarak. Bu yaklaşımın dört temel kazanımı vardır.

Birincisi dizi kazancı (array gain). `N` elemanlı bir dizi, istenen yönden gelen sinyalleri eşfazlı topladığında, sinyal genliği `N` kat (güç `N²` kat) büyürken, elemanlar arası ilişkisiz gürültü yalnızca `√N` kat (güç `N` kat) büyür. Net SNR kazancı:

```
 G_dizi = 10·log10(N)   [dB]    (ideal, ilişkisiz gürültü)
```

Sekiz elemanlı bir dizi ~9 dB, yüz elemanlı bir dizi ~20 dB kazanç verir — üstelik her eleman düşük kazançlı, geniş açılı olsa bile. Bu, "çok sayıda küçük antenle bir büyük antenin yönlülüğünü kurmak" demektir.

İkincisi yönlülük ve elektronik tarama (electronic steering). Mekanik bir çanağı çevirmek için motor gerekir; bir faz dizisinde hüzmeyi bir yönden diğerine taşımak yalnızca ağırlıkların fazını değiştirmektir — mikrosaniyeler içinde, hareketli parça olmadan. Aynı diziyle aynı anda birden çok hüzme bile oluşturulabilir (çok-hüzme).

Üçüncüsü uzamsal çeşitlilik (spatial diversity). Birbirinden yeterince ayrık (tipik olarak yarım dalga boyundan fazla) elemanlar, çok-yollu sönümlemeyi (multipath fading) bağımsız yaşar: bir eleman sönümlü bir noktadayken komşusu güçlü olabilir. Bunları birleştirmek (Bölüm 18'deki eşleştirilmiş filtre mantığının uzamsal karşılığı) sönümlemeye karşı dayanıklılık verir. Bu, MIMO'nun çeşitlilik kanadıdır (Kısım 9-10).

Dördüncüsü ve SIGINT açısından en kıymetlisi uzamsal ayrım: girişim bastırma ve yön kestirim. Bir dizi, gelen dalgaları geliş yönüne göre ayırabilir. İstenmeyen bir kaynağa (girişim, jammer) hüzme deseninde bir sıfır (null) çevirip onu bastırırken istenen yönü açık tutabilir (Kısım 7; Bölüm 13'teki null-steering savunmasının matematiği). Aynı uzamsal ayrım, tersinden, gelen sinyalin yönünü kestirmeyi (DoA) sağlar (Kısım 8; Bölüm 9'un alt-uzay algoritmalarının dizi temeli). Frekans ve zaman bir dalganın "ne" ve "ne zaman" olduğunu söyler; dizi "nereden" sorusunu ekler. Bu, üçüncü bir bağımsız işleme eksenidir.

> Mühendislik sezgisi: Bir dizinin yaptığı şey, uzayda bir eşleştirilmiş filtredir. Bölüm 18'deki matched filter, sinyalin bilinen zaman şekline göre en iyi toplamayı yapıyordu; bir beamformer, sinyalin bilinen uzamsal şekline (geliş yönüne karşılık gelen faz desenine) göre en iyi toplamayı yapar. "Faz desenini eşle, eşfazlı topla" — hem zaman hem uzayda aynı fikir.

---

<a id="2"></a>
## 2. Dizi Temelleri: Dizi Faktörü, Eleman Aralığı ve Faz İlerlemesi

### Düzlem dalga ve elemanlar arası faz farkı

Bir dizinin matematiği tek bir fiziksel gözlemle başlar: uzak bir kaynaktan gelen dalga, dizi boyutlarına göre düzlem dalga (plane wave) gibi davranır — dalga cepheleri paralel düzlemlerdir. Bu cephe diziye eğik (`θ` açısıyla) gelirse, komşu elemanlara farklı zamanlarda ulaşır; aralarında bir yol farkı, dolayısıyla bir faz farkı oluşur.

```
 Gelen düzlem dalga, θ açısıyla (broadside = 0°):

        ╲   ╲   ╲   ╲        dalga cepheleri (θ eğimli)
         ╲   ╲   ╲   ╲
   ───●────●────●────●───►  dizi ekseni (eleman aralığı d)
      0    1    2    3       eleman indeksi n
      │←d→│
            ekstra yol = d·sin θ   (komşu elemanlar arası)
```

İki komşu eleman arasındaki ek yol `d·sin θ`'dir. Bu yolun faz karşılığı, dalga sayısı `k = 2π/λ` ile:

```
 ψ = k · d · sin θ = (2π/λ) · d · sin θ        (komşu elemanlar arası faz farkı)
```

`θ = 0` (broadside, diziye dik geliş) ise `ψ = 0`: dalga tüm elemanlara aynı anda ulaşır. `θ = 90°` (endfire, dizi ekseni boyunca) ise `ψ` maksimumdur. Bu `ψ`, dizi matematiğinin tek en önemli niceliğidir; geometriyi (`d`), dalga boyunu (`λ`) ve yönü (`θ`) tek skalerde toplar. Bölüm 9'daki `θ = arcsin(λΔφ/2πd)` ifadesi tam olarak bunun tersidir: ölçülen faz farkından açıyı geri çözmek.

### Dizi faktörü (array factor)

Şimdi `N` elemanlı düzgün doğrusal dizide (ULA) her elemanın çıktısını bir karmaşık ağırlık `w_n` ile çarpıp toplayalım. `n`. eleman, referans elemana göre `n·ψ` kadar faz kazandığından, dizinin yöne bağlı toplam yanıtı dizi faktörüdür (array factor, AF):

```
 AF(θ) = Σ_{n=0}^{N−1} w_n · e^{ j n ψ },     ψ = (2π/λ) d sin θ
```

Bu, ağırlıkların bir tür uzamsal Fourier dönüşümüdür: `w_n` katsayıları zaman örnekleri gibi, `ψ` ise frekans gibi davranır. Tüm hüzme şekillendirme bu denklemin içindedir — `w_n`'i nasıl seçtiğin hüzmenin nereye bakacağını ve nasıl şekilleneceğini belirler.

En basit hal, eşit ağırlıklı (`w_n = 1`) düzgün dizidir. Geometrik seri toplamı kapalı formdadır:

```
 AF(θ) = Σ_{n=0}^{N−1} e^{ j n ψ } = (1 − e^{ j N ψ}) / (1 − e^{ j ψ})
       = e^{ j(N−1)ψ/2} · sin(Nψ/2) / sin(ψ/2)
```

Genlik (normalize, faz terimini atarak) Dirichlet çekirdeği veya "periyodik sinc"tir:

```
 |AF(ψ)| = | sin(N ψ/2) / sin(ψ/2) |
```

Normalize biçim (tepe = 1):

```
 |AF_n(ψ)| = | sin(N ψ/2) / (N · sin(ψ/2)) |
```

Bu fonksiyon, Bölüm 18'deki CIC filtresinin genlik yanıtıyla (`sin(Nx)/sin(x)` ailesi) aynı matematiksel akrabadır — bir dizi, uzamsal tanım bölgesinde çalışan bir FIR filtresidir; eleman aralığı örnekleme aralığı, `ψ` ise normalize uzamsal frekanstır. Bu benzeşim derin ve kullanışlıdır: filtre tasarımının her aracı (pencereleme, Chebyshev, Taylor) dizi sentezine birebir taşınır.

### Ana hüzme, yan loblar ve hüzme genişliği

`|AF_n(ψ)|`'in özellikleri (eşit ağırlık, broadside):

```
 |AF_n(ψ)| ▲
  1.0 ────────╮ ana hüzme (main lobe)
              │╲
              │ ╲   ilk null (sin(Nψ/2)=0 → ψ = 2π/N)
              │  ╲ ╱╲      ╱╲
              │   ╳  ╲    ╱  ╲    yan loblar (side lobes)
  ────────────┴────────────────────────► ψ
              0   2π/N        4π/N
```

Ana hüzme `ψ = 0`'da tepe yapar. İlk sıfır (first null) `N ψ/2 = π`, yani `ψ = 2π/N`'de oluşur. Eleman sayısı `N` arttıkça ana hüzme daralır — daha fazla eleman = daha keskin yönlülük. İlk yan lobun seviyesi, eşit-ağırlıklı (uniform) dizide `N` ne olursa olsun yaklaşık `−13.26 dB`'dir; bu, dikdörtgen pencerenin yan-lob seviyesinin (Bölüm 18, spektral sızıntı) tam uzamsal karşılığıdır.

Yarı-güç (−3 dB) hüzme genişliği (half-power beamwidth, HPBW), broadside'a bakan, `d = λ/2` aralıklı, eşit-ağırlıklı ULA için yaklaşık:

```
 HPBW ≈ 0.886 · λ / (N d)   radyan   (broadside, geniş N yaklaşımı)
      ≈ 102° / N            ( d = λ/2 için pratik kestirim )
```

Sekiz elemanlı, yarım-dalga aralıklı bir dizi broadside'da ~13° hüzme verir; bu, tek bir dipolün geniş donut'undan çok daha keskindir ve elemanların hiçbiri tek başına yönlü değildir. Yönlülük tamamen uzamsal toplamadan doğar.

### Tam ışıma deseni: eleman deseni × dizi faktörü

Dizinin gerçek ışıma deseni, dizi faktörünün tek başına değil, eleman deseni (element pattern) ile çarpımıdır — buna desen çarpımı teoremi (pattern multiplication) denir:

```
 F_toplam(θ) = g_eleman(θ) · AF(θ)
```

`g_eleman(θ)` tek bir elemanın (örn. dipol) kendi yön deseni. ULA'da bu, dizi faktörünün keskin lob yapısını eleman deseninin geniş zarfıyla pencereler: endfire'a (`θ → 90°`) doğru tek eleman kazancı düştüğü için dizi de orada zayıflar. Bu çarpım, gerçek bir dizinin neden ideal `|AF|`'ten saptığını açıklar ve tasarımda eleman seçimini diziyle birlikte değerlendirmeyi gerektirir.

> Mühendislik sezgisi: Dizi faktörünü "ışık demetinin yarık deneyi" gibi düşün. Her eleman bir kaynak; belli bir yönde hepsi eşfazlı varırsa yapıcı girişim (ana hüzme), faz farkları birikip iptal olursa yıkıcı girişim (null). Hüzme şekillendirme, bu girişim desenini ağırlıklarla bilerek kurmaktır — istediğin yöne yapıcı, istemediğin yöne yıkıcı.

---

<a id="3"></a>
## 3. Izgara Lobu (Grating Lobe) ve Neden d ≤ λ/2

Dizi faktörü `ψ`'ye göre `2π` periyotludur: `|AF_n(ψ)| = |AF_n(ψ + 2πm)|`, çünkü `e^{jnψ}` periyodiktir. Ana hüzme `ψ = 0`'da oluşur ama tam aynı yükseklikte tepeler `ψ = ±2π, ±4π, ...`'de de oluşur. Bunlara ızgara lobları (grating lobes) denir — ana hüzmenin istenmeyen tam-boy kopyaları. Bir ızgara lobu, fiziksel bir `θ` açısına denk geliyorsa, dizi o yönden gelen sinyali ana hüzme kadar güçlü kabul eder ve yön belirsizliği (ya da TX'te boşa giden güç) yaratır.

`ψ`'nin fiziksel olarak ulaşabileceği aralık `θ ∈ [−90°, +90°]` için (`sin θ ∈ [−1, +1]`):

```
 ψ = (2π/λ) d sin θ    ⟹    ψ ∈ [ −2π d/λ , +2π d/λ ]
```

Bu görünür bölge (visible region) genişliği `4π d/λ`'dir. Bir ızgara lobunun (en yakını `ψ = ±2π`) görünür bölgeye girmemesi için:

```
 2π d/λ < 2π   ⟹   d < λ        (ızgara lobu fiziksel açıya düşmesin, broadside)
```

Daha genel olarak, hüzme `θ_0` yönüne yönlendirilirse (Kısım 6), görünür bölge kayar ve her iki kenarın da güvende olması için sıkı koşul:

```
 d / λ  ≤  1 / (1 + |sin θ_0|)
```

`θ_0 = 0` (broadside) için bu `d ≤ λ`; ama hüzmeyi `θ_0 = 90°`'ye (endfire) kadar tarayabilmek istersen payda 2 olur ve koşul `d ≤ λ/2`'ye sıkışır. İşte dizi tasarımının altın kuralı buradan gelir:

```
 d ≤ λ/2   ⟹   hiçbir yöne yönlendirmede ızgara lobu oluşmaz (tam tarama güvenli)
```

```
 d = λ/2:  görünür bölge tam [−π, +π]            (ızgara lobu sınırda, içeri girmez)
 d = λ:    görünür bölge [−2π, +2π]               (ψ=±2π ızgara lobu kenara girer)

 |AF| ▲ ana          ızgara                       ızgara
      │ ╱╲           ╱╲                            ╱╲
      │╱  ╲  ╱╲ ╱╲ ╱╱  ╲ ...                    ╱  ╲
   ───┴────────────────────────────────────────────────► ψ
     0   π/2  π        2π                          4π
       d=λ/2: yalnız [−π,π] görünür   d=λ: 2π'deki kopya da görünür → belirsizlik
```

Bu, Bölüm 18'deki Nyquist örneklemenin uzamsal kardeşidir. Eleman aralığı `d` bir uzamsal örnekleme aralığıdır; `d = λ/2`, `sin θ` ekseninde tam Nyquist hızıdır. `d > λ/2` uzamsal alt-örnekleme (undersampling) demektir ve ızgara lobu tam olarak uzamsal aliasing'dir — yüksek "uzamsal frekanslı" (büyük açılı) bir geliş, düşük açılı bir gelişin kılığında geri katlanır. Bölüm 9'daki faz interferometrisinde "`d > λ/2` belirsizlik yaratır" uyarısının kökü budur: ızgara lobu = açı belirsizliği.

Pratik sonuç: SIGINT ve yön bulma dizilerinde eleman aralığı neredeyse her zaman `λ/2` (ya da altında) seçilir. İlgilenilen bant geniş ise (örn. KrakenSDR ile geniş HF/VHF taraması), `λ` frekansla değişir; `d` sabit olduğundan dizi yalnızca belli bir frekans aralığında `d ≤ λ/2` koşulunu sağlar. Bu yüzden çok-bant DF sistemleri ya birden çok aralıklı eleman setine ya da frekansa göre eleman seçimine başvurur.

---

<a id="4"></a>
## 4. Dizi Geometrileri: ULA, UCA, Düzlemsel, Seyrek Diziler

Eleman yerleşimi, dizinin hangi açıları nasıl ayırabileceğini doğrudan belirler. Dört temel geometri.

### Düzgün Doğrusal Dizi (ULA — Uniform Linear Array)

Elemanlar tek bir hat üzerinde, eşit `d` aralıkla. En basit, en çok incelenen geometri (Kısım 2'nin tüm formülleri ULA içindir). Güçlü yön: matematiği kapalı form, MUSIC/ESPRIT en temiz burada. Zayıf yön: yalnızca tek bir açıyı (azimut, dizi eksenine göre `θ`) çözebilir ve önden-arkadan belirsizliği (front-back ambiguity) vardır — `+θ` ile `−θ`'yi (dizi ekseni etrafındaki koni) ayıramaz, çünkü her ikisi de aynı `sin θ`'yi verir. Konum vektörleri:

```
 p_n = (n·d, 0, 0),   n = 0 … N−1
```

### Düzgün Dairesel Dizi (UCA — Uniform Circular Array)

Elemanlar `R` yarıçaplı bir çember üzerinde eşit açısal aralıkla. KrakenSDR'ın DF için önerdiği klasik geometri budur (Bölüm 9, Kısım 8). En büyük avantajı 360° azimut kapsama ve simetri: her yöne benzer çözünürlük verir, ULA'nın önden-arkadan belirsizliğini doğal olarak çözer. `n`. elemanın açısal konumu `γ_n = 2πn/N`, konum:

```
 p_n = ( R cos γ_n , R sin γ_n , 0 ),   γ_n = 2πn/N
```

Bir `φ` azimut ve `θ` yükseliş açısından gelen düzlem dalga için `n`. elemandaki faz:

```
 Φ_n(φ,θ) = (2π/λ) · R · sin θ · cos(φ − γ_n)
```

UCA'da dizi faktörü Bessel fonksiyonlarıyla ifade edilir (faz argümanı kosinüs içinde olduğu için), bu yüzden ULA kadar basit kapalı form yoktur ama "faz modu" (phase mode) ayrıştırmasıyla ele alınır. Pratikte hesap sayısal yapılır; kavram aynıdır.

### Düzlemsel Dizi (Planar / URA — Uniform Rectangular Array)

Elemanlar bir ızgarada (`M × N`), iki eksende de aralıklı. Hem azimut hem yükseliş açısını (`φ, θ` ikisini birden) çözebilen ilk geometri — 2B yön bulma ve 3B hüzme yönlendirme bunu gerektirir. 5G Massive MIMO panelleri (Kısım 11) tipik olarak düzlemsel dizilerdir. Yönlendirme vektörü iki eksenin Kronecker çarpımıdır:

```
 a(φ,θ) = a_x(φ,θ) ⊗ a_y(φ,θ)       ( ⊗ = Kronecker çarpımı )
```

Bu çarpım yapısı, 2B hüzme yönetimini iki ayrı 1B problemine indirger ve hesabı büyük ölçüde basitleştirir (özellikle hibrit beamforming'de, Kısım 6).

### Seyrek / Rastgele Diziler (Sparse / Random Arrays)

Tüm `λ/2` ızgara dolu olmak zorunda değildir. Seyrek diziler, az sayıda elemanı geniş bir açıklığa (aperture) yayar. Aynı eleman sayısıyla daha geniş açıklık = daha dar ana hüzme = daha iyi açısal çözünürlük; ama bedeli daha yüksek yan loblar ve dikkatli tasarlanmazsa ızgara lobu riskidir. İki önemli alt-tür:

- Minimum-redundancy / minimum-hole diziler: eleman çiftleri arasındaki tüm farkları (lag) en az fazlalıkla kapsayacak şekilde yerleştirilir; korelasyon (kovaryans) tahmini için tüm gerekli faz farklarını az elemanla üretir.
- Eş-asal (coprime) ve iç içe (nested) diziler: iki alt-dizinin (aralıkları eş-asal) birleşimi, fiziksel eleman sayısından çok daha fazla "sanal" eleman (difference co-array) üretir; `N` fiziksel elemanla `O(N²)` serbestlik derecesine ulaşılır. Bu, az donanımla çok kaynak ayırmanın modern yoludur.

```
 ULA (dolu):    ● ● ● ● ● ● ● ●        N=8, açıklık=7d, basit
 Seyrek:        ●   ●     ●       ●     4 eleman, geniş açıklık, dar hüzme
 Eş-asal:       2 alt-dizi (aralık 2d ve 3d) → fark co-array çok daha yoğun
```

| Geometri | Kapsama | Çözebildiği açı | Tipik kullanım |
|---|---|---|---|
| ULA | 1B (önden-arkadan belirsiz) | Azimut (tek) | Klasik DF, en temiz matematik |
| UCA | 360° azimut | Azimut (+ kısıtlı yükseliş) | KrakenSDR, sabit DF istasyonu |
| Düzlemsel (URA) | 2B yarım-küre | Azimut + yükseliş | 5G Massive MIMO, 3B beamforming |
| Seyrek/eş-asal | Geniş açıklık, az eleman | Yüksek çözünürlük | Az donanımla çok kaynak ayırma |

---

<a id="5"></a>
## 5. Yönlendirme Vektörü (Steering Vector) ve Dizi Veri Modeli

Bütün uzamsal işlemeyi tek bir vektör-matris diline sıkıştıran kavram yönlendirme vektörüdür (steering vector, ya da array manifold vektörü). Tek bir `θ` yönünden gelen düzlem dalganın, her elemanda hangi karmaşık fazla göründüğünü toplayan vektör budur. ULA için (`d` aralık, referans eleman 0):

```
        ┌      1            ┐
        │  e^{ j ψ}         │
 a(θ) = │  e^{ j 2ψ}        │ ,   ψ = (2π/λ) d sin θ        (N×1 vektör)
        │     ⋮             │
        └  e^{ j (N−1)ψ}    ┘
```

`a(θ)`, "o yönden gelen birim genlikli bir dalganın dizi üzerindeki imzası"dır. Her olası `θ`'nin `a(θ)`'sının kümesine dizi manifoldu (array manifold) denir; uzamsal işleme, gözlenen veriyi bu manifolddaki vektörlerle eşleştirme problemidir.

### Veri modeli: çoklu kaynak + gürültü

`D` adet uzak kaynak (`s_1, …, s_D`) sırasıyla `θ_1, …, θ_D` yönlerinden geliyorsa, dizinin `t` anındaki `N×1` örnek vektörü `x(t)` lineer üst üste binmedir:

```
 x(t) = Σ_{i=1}^{D} a(θ_i) · s_i(t) + n(t) = A · s(t) + n(t)
```

Burada `A = [ a(θ_1) | a(θ_2) | … | a(θ_D) ]` `N×D` yönlendirme (manifold) matrisi, `s(t)` kaynak vektörü, `n(t)` ilişkisiz gürültü. Bu denklem bu bölümün ve Bölüm 9'un alt-uzay yöntemlerinin tüm temelidir — sade görünür ama içinde geliş açıları, kaynak sayısı ve uzamsal yapı saklıdır.

### Kovaryans matrisi: uzamsal işlemenin merkezi

Tek bir örnek vektörü gürültülüdür; istatistik gerekir. Diziler dünyasının merkezi nesnesi `N×N` uzamsal kovaryans matrisidir (spatial covariance / correlation matrix):

```
 R = E{ x(t) x^H(t) } = A · R_s · A^H + σ² I
```

`x^H` Hermitsel (eşlenik-devrik) transpoz, `R_s = E{s s^H}` kaynak kovaryansı, `σ²` gürültü gücü, `I` birim matris. Pratikte beklenen değer `K` örnek üzerinden örneklem kovaryansı ile kestirilir:

```
 R̂ = (1/K) Σ_{t=1}^{K} x(t) x^H(t)
```

Bütün beamforming ve DoA algoritmaları `R̂`'nin (ya da onun özayrışımının) bir fonksiyonudur. `R̂`'yi doğru kestirmek için yeterli örnek (`K ≫ N`) ve kaynakların yeterince hareketsiz/ayrık olması gerekir; bu, Bölüm 9'da MUSIC için "yeterli örnek topla" uyarısının sayısal nedenidir. `R`'nin yapısı şunu söyler: sinyal enerjisi `A` sütunlarının gerdiği `D`-boyutlu sinyal alt-uzayında yoğunlaşır; geri kalan `N−D` boyut yalnızca gürültüdür (gürültü alt-uzayı). MUSIC ve ESPRIT bu ayrımı sömürür (Kısım 8).

> Mühendislik sezgisi: `a(θ)`'yı "uzamsal bir karmaşık üstel" olarak gör. Bölüm 18'de bir frekansı seçmek için sinyali `e^{−j2πf₀n/f_s}` ile çarpıyorduk; burada bir yönü seçmek için örnek vektörünü `a(θ_0)`'ın eşleniğiyle çarpıp topluyoruz. Frekans seçimi zaman örneklerinin Fourier'i, yön seçimi uzay örneklerinin (elemanların) Fourier'i. Aynı matematik, farklı eksen.

---

<a id="6"></a>
## 6. Hüzme Şekillendirme: Gecikme-Topla, Ağırlık Vektörü, Analog/Sayısal/Hibrit

Hüzme şekillendirme (beamforming), dizi çıkışını tek bir skalere indirgeyen ağırlıklı toplamadır:

```
 y(t) = Σ_{n=0}^{N−1} w_n^* · x_n(t) = w^H · x(t)
```

`w = [w_0, …, w_{N−1}]^T` ağırlık vektörü (beamforming weights). Dizinin bir `θ` yönüne yanıtı, ağırlık vektörünün o yönün yönlendirme vektörüyle iç çarpımıdır:

```
 B(θ) = w^H · a(θ)        (hüzme deseni — beam pattern)
```

`B(θ)`, Kısım 2'deki dizi faktörünün genel halidir (`w_n = 1` koyarsan `AF`'e iner). Tüm soru: `w`'yi nasıl seçeriz?

### Gecikme-topla / faz-kaydır-topla (delay-and-sum / conventional beamformer)

En basit ve en sağlam beamformer, hüzmeyi istenen `θ_0` yönüne kilitleyen ağırlığı seçer: her elemanın o yöndeki fazını tam tersiyle düzelt, böylece `θ_0`'dan gelen dalga tüm elemanlarda eşfazlı toplanır.

```
 w = a(θ_0) / N        ⟹    B(θ_0) = a^H(θ_0) a(θ_0) / N = 1   (tepe θ_0'da)
```

Ağırlığın `n`. bileşeni `e^{−jnψ_0}/N` (`ψ_0 = (2π/λ)d sin θ_0`) — yani her elemana `θ_0`'ı broadside'a getiren bir faz kayması (phase shift) uygulanır. Bu yüzden buna faz dizisi (phased array) denir: hüzmeyi yönlendirmek = ağırlıkların fazını ayarlamak.

```
 Hüzme yönlendirme (steering):  her elemana ilerleyen faz ekle

   eleman:   0      1      2      3
   faz:      0    −ψ_0   −2ψ_0  −3ψ_0     (θ_0'dan geleni eşfazla)
             │      │      │      │
             ▼      ▼      ▼      ▼
            ───── Σ ───── → hüzme tepesi artık θ_0 yönünde
```

Yönlendirilmiş dizi faktörü, broadside desenin kaydırılmışıdır — `ψ` ekseninde `ψ → ψ − ψ_0`:

```
 |B(θ)| = | sin(N(ψ−ψ_0)/2) / (N sin((ψ−ψ_0)/2)) |
```

Geniş-bantlı sinyallerde dikkat: faz kaydırma yalnızca tek frekansta doğru gecikmeyi verir. Gerçek geniş-bant beamforming, faz değil gerçek zaman gecikmesi (true time delay) gerektirir; aksi halde hüzme yönü frekansla kayar (beam squint). Dar-bant kabul (sinyal bandı taşıyıcıya göre küçük) bu bölümün varsayımıdır; geniş-bantta gecikme tabanlı yapılar ya da alt-bant beamforming kullanılır.

### Pencereleme ile yan-lob kontrolü

Eşit ağırlık (`w_n = 1`) en dar ana hüzmeyi verir ama −13 dB yan loblar bırakır. Tıpkı FFT'de spektral sızıntı için pencere kullandığımız gibi (Bölüm 18, Kısım 13), dizide de ağırlıklara bir uzamsal pencere (taper) uygulayarak yan lobları bastırırız — bedeli ana hüzmenin genişlemesi:

```
 w_n = pencere(n) · e^{−j n ψ_0}      (örn. Hamming, Chebyshev, Taylor taper)
```

Dolph-Chebyshev taper, verilen ana-hüzme genişliği için tüm yan lobları eşit ve mümkün olan en düşük seviyeye indiren optimaldir (Parks-McClellan'ın uzamsal karşılığı). Pencere ↔ yan-lob ödünleşimi, Bölüm 18'deki spektral pencere ödünleşimiyle birebir aynıdır: bunu kavradıysan dizi sentezini de kavramışsındır.

### Analog vs sayısal vs hibrit beamforming

Ağırlıkları nerede uyguladığın bir mimari karardır ve maliyet/esneklik dengesini belirler.

```
 Analog BF:    RF fazlayıcılar, tek ADC
   ant ─[φ]─┐
   ant ─[φ]─┼─Σ─[karıştır]─[ADC]─ tek akış   (faz RF'te, 1 zincir)
   ant ─[φ]─┘

 Sayısal BF:   her eleman tam RF zinciri + ADC
   ant ─[RF+ADC]─┐
   ant ─[RF+ADC]─┼─(yazılımda w^H x)─ N akış   (tam esneklik, çok-hüzme)
   ant ─[RF+ADC]─┘

 Hibrit BF:    alt-dizi başına analog faz + az sayıda dijital zincir
   [alt-dizi: φ φ φ]─[ADC]─┐
   [alt-dizi: φ φ φ]─[ADC]─┼─ K≪N akış        (denge: maliyet ↔ esneklik)
   [alt-dizi: φ φ φ]─[ADC]─┘
```

| Mimari | RF zinciri sayısı | Esneklik | Maliyet/güç | Tipik yer |
|---|---|---|---|---|
| Analog | 1 (tüm dizi) | Tek hüzme, yavaş | Düşük | Klasik radar, basit faz dizisi |
| Sayısal | N (her eleman) | Çok-hüzme, anlık adaptif | Yüksek (N ADC) | KrakenSDR, SIGINT, MVDR/MUSIC |
| Hibrit | K (1 < K ≪ N) | Sınırlı çok-hüzme | Orta | 5G Massive MIMO, mmWave |

SIGINT ve yön bulma sayısal beamforming ister, çünkü MVDR/MUSIC gibi adaptif algoritmalar her elemanın ayrı karmaşık örneğine ihtiyaç duyar — bir analog dizide bu örnekler RF'te toplanıp kaybolmuştur. KrakenSDR'ın beş ayrı RTL-SDR kanalı tam da bunun içindir: tam sayısal, faz-tutarlı bir dizi (Kısım 12).

> Mühendislik sezgisi: Sayısal beamforming'in büyüsü, aynı `x(t)` veri setinden sonsuz sayıda farklı `w` ile sonsuz farklı hüzme oluşturabilmendir — kayıt yapıldıktan sonra bile. Diziye baktığın "yönü", veriyi topladıktan sonra yazılımda seçersin. Bir KrakenSDR kaydında her yöne ayrı ayrı "bakıp" hangi yönde sinyal olduğunu sonradan tarayabilirsin; analog dizide bu imkânsızdır çünkü yön donanımda dondurulmuştur.

---

<a id="7"></a>
## 7. Uyarlanabilir Diziler: MVDR/LCMV, Null Steering, Wiener/LMS

Gecikme-topla beamformer'ı sabittir: yan lobları geometrinin belirlediği yerdedir ve güçlü bir girişim tam bir yan loba denk gelirse onu bastıramaz. Uyarlanabilir diziler (adaptive arrays), ağırlıkları veriye (kovaryans `R`'ye) göre seçerek girişimi nerede olursa olsun bastırır. Bu, dizi işlemenin SIGINT ve savunma açısından en güçlü kısmıdır.

### MVDR / Capon beamformer (minimum varyans)

İstenen yön `θ_0`'dan gelen sinyali bozmadan geçirirken (kazanç tam 1) toplam çıkış gücünü en aza indir. Çıkış gücü `w^H R w` olduğundan, kısıtlı optimizasyon:

```
 min_w  w^H R w     öyle ki    w^H a(θ_0) = 1
```

Lagrange çarpanlarıyla çözüm Minimum Variance Distortionless Response (MVDR, ya da Capon) ağırlığıdır:

```
 w_MVDR = ( R^{−1} a(θ_0) ) / ( a^H(θ_0) R^{−1} a(θ_0) )
```

Sezgi derindir: `θ_0`'ı sabit tutup toplam gücü küçültmenin tek yolu, gücün geri kalanını (girişim + gürültü) bastırmaktır. Algoritma, `θ_0` dışındaki güçlü kaynaklara otomatik olarak null çevirir — onları nereye yerleştireceğini önceden bilmeden, sadece `R^{−1}` üzerinden. Bu, sabit beamformer'a göre devrimsel bir farktır: girişim hangi açıdan gelirse gelsin, MVDR onu öğrenir ve siler.

### LCMV: çoklu kısıt

MVDR tek bir kısıt (bir yönü koru) kullanır. Linearly Constrained Minimum Variance (LCMV) birden çok lineer kısıt uygular — örneğin "`θ_0`'ı koru ve aynı anda `θ_1`, `θ_2`'ye tam null koy". Kısıt matrisi `C` (`N×L`) ve istenen yanıt vektörü `f` (`L×1`) ile:

```
 min_w  w^H R w     öyle ki    C^H w = f

 w_LCMV = R^{−1} C ( C^H R^{−1} C )^{−1} f
```

`C`'ye istenen yönün steering vektörünü (yanıt 1) ve bastırılacak yönlerin steering vektörlerini (yanıt 0) koyarsan, dizi istenen yöne bakarken seçilen yönlere kesin sıfır çevirir.

### Null steering (sıfır çevirme)

Tersinden bakıldığında, hedefe bakmak değil belirli bir yöne sıfır koymak da bir kısıttır. `K` adet girişim yönü (`θ_1 … θ_K`) biliniyorsa, ağırlığı bu yönlerin yönlendirme vektörlerine dik olacak şekilde seçeriz:

```
 w^H a(θ_k) = 0,   k = 1 … K       (her girişim yönüne tam null)
 w^H a(θ_0) = 1                     (istenen yönü koru)
```

```
 İstenen sinyal θ_0 = +20°,  girişim θ_1 = −40°:

  |B(θ)| ▲
   1.0 ──────────╮ ana hüzme (θ_0 = +20°)
                 │╲
        ╱╲   ╱╲  │ ╲   ╱╲
   ─────╳────╳───┴──╳──────────► θ
        │         │
      −40°       +20°
        ▼
       NULL (B(θ_1)=0): girişim bastırıldı
```

Bölüm 13'teki anti-jam savunma katmanı olan "null-steering"in matematiği tam budur: alıcı dizisi, jammer yönüne hüzme deseninde bir sıfır koyarak istenen sinyali kurtarır. Kritik sınır: bir `N` elemanlı dizi en fazla `N−1` bağımsız null oluşturabilir (uzamsal serbestlik derecesi `N−1`). Beş kanallı bir KrakenSDR en çok dört girişimi aynı anda bastırabilir; daha fazlası için daha çok eleman gerekir. Bu, "kaç jammer'a karşı dayanıklıyım?" sorusunun donanımsal cevabıdır.

> Yasal not: Buradaki null steering tamamen alış tarafıdır — alıcı kendi RAM'indeki örneklerden istenmeyen bir kaynağı uzamsal olarak filtreler. Bu, pasif ve savunmacıdır. Verici tarafı yönlü enerji yayma (örneğin bir hedefe yönlü karıştırma) tamamen farklı ve genellikle yasa dışıdır; bu bölüm onu kapsamaz.

### Wiener çözümü ve LMS uyarlaması

Kovaryansı baştan bilmiyorsak, ağırlıkları çevrimiçi öğrenebiliriz. Bir referans/eğitim sinyali `d(t)` varsa (örn. bilinen bir pilot ya da preamble), optimal ağırlık Wiener çözümüdür:

```
 w_opt = R^{−1} · p,      p = E{ x(t) d^*(t) }      (çapraz korelasyon)
```

`R^{−1}`'i hesaplamak pahalı olduğundan, en dik iniş (gradient descent) ile yinelemeli yaklaşılır — LMS (Least Mean Squares) algoritması:

```
 e(t)   = d(t) − w^H(t) x(t)            (hata)
 w(t+1) = w(t) + μ · x(t) · e^*(t)      (ağırlık güncelle)
```

`μ` adım büyüklüğü (yakınsama hızı ↔ kararlılık ödünleşimi; çok büyük `μ` ıraksar). LMS, donanımda ucuzdur (matris tersi yok) ve girişim ile sinyal değiştikçe ağırlıkları kovalar — uyarlanabilir gürültü/girişim iptalinin (Bölüm 18'deki adaptif I/Q iptaline akraba) uzamsal halidir. Daha hızlı yakınsama için RLS (Recursive Least Squares) kullanılır ama hesabı ağırdır.

---

<a id="8"></a>
## 8. Yön Kestirim (DoA): Gecikme-Topla, Capon, MUSIC, ESPRIT

Beamforming "bilinen bir yöne bak" diyordu; yön kestirim (Direction of Arrival, DoA) ters problemi çözer: "sinyaller hangi yönlerden geliyor?" Bu, Bölüm 9'un yön bulma algoritmalarının dizi-matematik tarafıdır; burada her yöntemin `R`'den açıyı nasıl çıkardığını türetiyoruz.

### Gecikme-topla (Bartlett) spektrumu

En basit DoA, beamformer'ı tüm açılarda tarayıp çıkış gücünü çizmektir. Bartlett (conventional) spektrumu:

```
 P_Bartlett(θ) = a^H(θ) · R · a(θ)
```

Her `θ` için diziyi o yöne yönlendirip ne kadar güç geldiğine bakarsın; tepe yapan açılar kaynak yönleridir. Basit ve sağlam ama çözünürlüğü dizi açıklığıyla sınırlıdır (Rayleigh limiti): iki kaynak ana hüzme genişliğinden (`~λ/Nd`) yakınsa ayrılamaz. Bu, "uzamsal RBW" sınırıdır — Bölüm 18'deki frekans çözünürlüğünün (RBW) uzamsal kardeşi.

### Capon (MVDR) spektrumu

MVDR ağırlığını her açıya uygulayıp çıkış gücüne bakarsak süper-Bartlett bir spektrum elde ederiz:

```
 P_Capon(θ) = 1 / ( a^H(θ) · R^{−1} · a(θ) )
```

Capon, her yönde "o yön dışındaki her şeyi bastırarak" baktığı için Bartlett'ten çok daha keskin tepeler verir; yakın kaynakları ayırabilir. Bedeli, `R`'nin tersine ihtiyaç (gürültüye/az örneğe daha hassas) ve korele kaynaklarda bozulmadır.

### MUSIC (MUltiple SIgnal Classification)

MUSIC, çözünürlüğü açıklık limitinin ötesine taşıyan (süper-çözünürlük) alt-uzay yöntemidir. Temel: `R`'nin özayrışımında, en büyük `D` özdeğer sinyal alt-uzayını, kalan `N−D` özdeğer (hepsi `≈ σ²`) gürültü alt-uzayını gerer. Kilit gerçek: her gerçek yönün `a(θ_i)` vektörü sinyal alt-uzayındadır, dolayısıyla gürültü alt-uzayına diktir.

```
 R = U_s Λ_s U_s^H + U_n Λ_n U_n^H      (özayrışım: sinyal + gürültü alt-uzayı)
 U_n = [N×(N−D)] gürültü alt-uzayı özvektörleri

 P_MUSIC(θ) = 1 / ( a^H(θ) · U_n U_n^H · a(θ) )
```

Payda, `a(θ)`'nın gürültü alt-uzayına izdüşümünün normudur. `θ` bir gerçek kaynak yönüne eşitse `a(θ) ⊥ U_n` olur, payda `→ 0`, spektrum patlar (keskin tepe). MUSIC tepeyi "güç" olarak değil "gürültü alt-uzayına diklik keskinliği" olarak arar — Bölüm 9'daki "diklik keskinliği" ifadesinin tam matematiği budur.

```
 P_MUSIC ▲          ▌                    ▌
         │          ▌ keskin tepe        ▌ keskin tepe
         │          ▌ (θ_1)              ▌ (θ_2)
         │   ___────▌▌▌────________──────▌▌▌────___
   ──────┴──────────────────────────────────────────► θ
         (gerçek açılarda a(θ)⊥U_n → payda→0 → patlar)
```

MUSIC'in koşulları (Bölüm 9'daki uyarının nedeni): (1) kaynak sayısı `D` doğru kestirilmeli (özdeğer dağılımından, örn. MDL/AIC ile); (2) `N > D` (gürültü alt-uzayı boş olmamalı); (3) kanallar iyi kalibre, faz-tutarlı olmalı (Kısım 12); (4) korele kaynaklar (örn. çok-yoldan kendiyle korele sinyal) `R_s`'i tekil yapar ve MUSIC bozulur — çözüm uzamsal yumuşatma (spatial smoothing): diziyi örtüşen alt-dizilere bölüp kovaryansları ortalamak korelasyonu kırar (yalnız ULA/URA'da temiz çalışır).

### ESPRIT (rotational invariance)

ESPRIT (Estimation of Signal Parameters via Rotational Invariance Techniques), açıları bir spektrumu tüm `θ`'larda tarayarak değil, doğrudan bir özdeğer probleminden cebirsel olarak çözer — hesabı MUSIC'ten hafif, tarama gerektirmez. Dizinin kaydırma değişmezliği (shift invariance) yapısını kullanır: ULA'yı iki örtüşen alt-diziye (`x_1`, biri bir eleman kaydırılmış `x_2`) böl. İki alt-dizinin sinyal alt-uzayları, kaynak yönlerini taşıyan köşegen bir faz matrisi `Φ` ile ilişkilidir:

```
 U_{s2} = U_{s1} · Ψ ,        Ψ'nin özdeğerleri = e^{ j ψ_i }
 ψ_i = (2π/λ) d sin θ_i   ⟹   θ_i = arcsin( λ·angle(öz_i) / (2π d) )
```

Açılar doğrudan özdeğerlerin fazından okunur. ESPRIT'in avantajı: açı tarama (grid search) yok, kalibrasyona MUSIC kadar duyarlı değil (kaydırma yapısı kalibrasyonu kısmen massediyor); dezavantajı: kaydırma-değişmez geometri (ULA gibi düzgün yapı) şart.

| Yöntem | Çözünürlük | Korele kaynak | Hesap | Geometri | Tipik kullanım |
|---|---|---|---|---|---|
| Bartlett (gecikme-topla) | Düşük (açıklık limiti) | Dayanıklı | Çok düşük | Her | Hızlı kaba tarama |
| Capon (MVDR) | Orta-yüksek | Bozulur | Orta (R^{−1}) | Her | Adaptif, az kaynak |
| MUSIC | Çok yüksek (süper-çöz.) | Bozulur (smoothing gerekir) | Yüksek (özayrışım+tarama) | Her | Yoğun/çok-yollu ortam |
| ESPRIT | Çok yüksek | Bozulur | Orta (taramasız) | ULA/düzenli | Hızlı, kalibrasyona toleranslı |

> Mühendislik sezgisi: Üç yöntemi bir merdiven gibi düşün. Bartlett = "her yöne bak, en parlak nereye?" (sezgisel, kaba). Capon = "her yöne bak ama gerisini sustur" (keskin). MUSIC/ESPRIT = "güce hiç bakma; hangi yönün imzası gürültü uzayına dik?" (süper-çözünürlük, ama kalibrasyon ve kaynak-sayısı hassas). Bir KrakenSDR çıktısında gördüğün keskin DOA tepesi neredeyse her zaman MUSIC ailesindendir; o keskinlik güçten değil, dikliktendir.

---

<a id="9"></a>
## 9. MIMO ve Massive MIMO: Uzamsal Çoğullama ve Kanal Kapasitesi

Şimdiye kadar dizi tek tarafta (alıcı) idi ve tek bir akışa hizmet ediyordu. MIMO (Multiple-Input Multiple-Output), hem verici hem alıcı tarafa dizi koyup, açtığı çok-boyutlu uzamsal kanalı bir kapasite kazancına çevirir. Bu, Bölüm 1'deki Shannon kapasitesinin uzamsal eksene genişlemesidir ve modern hücresel/WiFi sistemlerinin (Bölüm 20) çekirdek teknolojisidir.

### MIMO kanal modeli

`N_t` verici ve `N_r` alıcı anten arasında, alınan vektör `y` (`N_r×1`), kanal matrisi `H` (`N_r×N_t`) üzerinden:

```
 y = H x + n
```

`H`'nin `(i,j)` elemanı, `j`. TX anteninden `i`. RX antenine olan karmaşık kanal kazancıdır. Tek-anten (SISO) durumunda `H` bir skalerdir; MIMO'da bir matristir ve içindeki bağımsız "yollar" kapasiteyi katlar.

### MIMO kapasite formülü

Kanal alıcıda bilinir ama vericide bilinmezse (toplam güç `P` antenlere eşit dağıtılır), kapasite:

```
 C = log2 · det( I_{N_r} + (ρ / N_t) · H H^H )      [bit/s/Hz]
```

`ρ = P/σ²` ortalama SNR, `det` determinant, `I` birim matris. Bu, Shannon `C = log2(1+SNR)`'nin matris genellemesidir. `H`'yi tekil değer ayrışımıyla (SVD) köşegenleştirelim: `H = U Σ V^H`, tekil değerler `σ_1 … σ_r` (`r = rank(H) ≤ min(N_t, N_r)`). Kapasite paralel alt-kanalların (eigen-kanallar) toplamına ayrışır:

```
 C = Σ_{i=1}^{r} log2( 1 + (ρ/N_t) · σ_i² )
```

Bu denklem MIMO'nun büyüsünü açıklar: `H` tam-ranklı ve tekil değerleri dengeliyse, kapasite yaklaşık `min(N_t, N_r)` katına çıkar — aynı bant, aynı güçle. Buna uzamsal çoğullama kazancı (spatial multiplexing gain) denir: zengin saçılmalı (rich scattering) bir ortam, bağımsız uzamsal yollar açar ve her yol ayrı bir veri akışı taşır.

```
 SISO:   tek kanal,  C = log2(1+ρ)
 MIMO:   r paralel eigen-kanal,  C ≈ Σ log2(1 + ρ σ_i²/N_t)

   TX akışları       H (saçılma)        RX
   s1 ─[ant]─╮  ╱─────────────╲  ╭─[ant]─ ŝ1
   s2 ─[ant]─┼─╳ çok-yollu kanal╳─┼─[ant]─ ŝ2   → r bağımsız akış
   s3 ─[ant]─╯  ╲─────────────╱  ╰─[ant]─ ŝ3      (kapasite ×r)
```

Kanal vericide de biliniyorsa (CSIT), güç eigen-kanallara su-doldurma (water-filling) ile en uygun dağıtılır — güçlü kanallara çok, zayıflara az güç:

```
 P_i = max( 0 , μ − N_t σ²/σ_i² ),      Σ P_i = P   (su seviyesi μ ile)
```

### Massive MIMO

Massive MIMO, antenlerin sayısını radikal artırır — baz istasyonunda yüzlerce eleman (`N_t ≫ kullanıcı sayısı K`). Üç sınır-davranış kazanımı vardır:

- Kanal sertleşmesi (channel hardening): `N` büyüdükçe rastgele kanalın etkisi ortalamada düzleşir, küçük-ölçek sönümleme yok olur; kanal neredeyse deterministik davranır.
- Uygun yayılma (favorable propagation): farklı kullanıcıların kanal vektörleri, `N → ∞` iken neredeyse dik (ortogonal) hale gelir; basit lineer işlem (MRC/ZF) kullanıcıları neredeyse mükemmel ayırır.
- Dizi kazancı: `N` antenli koherent toplama, kullanıcı başına SNR'a `~10·log10(N)` dB ekler — TX gücünü düşürüp enerji verimi sağlar.

Çok-kullanıcılı MIMO'da (MU-MIMO), baz istasyonu `K` kullanıcıya aynı zaman-frekans kaynağında ayrı uzamsal hüzmelerle hizmet eder. Lineer precoder örnekleri (downlink), kanal matrisi `H` (`K×N`) için:

```
 Eşleşik-filtre (MRT):  W = H^H                          (basit, kullanıcı içi)
 Sıfır-zorlama (ZF):    W = H^H (H H^H)^{−1}              (kullanıcılar arası girişim → 0)
 Regülarize (MMSE):     W = H^H (H H^H + (K/ρ) I)^{−1}    (gürültü-girişim dengesi)
```

ZF precoder, her kullanıcıya giderken diğer `K−1` kullanıcının yönüne null koyar (Kısım 7'nin TX karşılığı — ama bu meşru bir baz istasyonunun kendi kullanıcılarını ayırmasıdır). Massive MIMO'da `N ≫ K` olduğundan bu nulllar için bol uzamsal serbestlik vardır.

> SIGINT bağı: Bir analist için Massive MIMO iki şey demektir. Birincisi, baz istasyonunun enerjisi artık her yöne yayılmaz; belirli kullanıcılara dar hüzmelerle gider, dolayısıyla "yanlış" konumdan bakan bir pasif alıcı sinyali zayıf görür ya da hiç görmez (Bölüm 20'deki "beamforming'in dolaylı gizlilik etkisi"). İkincisi, downlink'teki çok-hüzmeli yapı, bir gözlemciye o anda kaç kullanıcının nerede olduğu hakkında uzamsal bilgi sızdırabilir — hüzme yönetimi sinyalleri (Kısım 11) bunun izidir.

---

<a id="10"></a>
## 10. Çeşitlilik-Çoğullama Ödünleşimi ve Hüzme-Uzayı (Beamspace)

MIMO'yu iki farklı amaçla kullanabilirsin ve ikisi aynı anda tam olarak elde edilemez — bu temel gerilim çeşitlilik-çoğullama ödünleşimidir (diversity-multiplexing tradeoff, DMT).

### İki kazanç türü

Çoğullama kazancı (multiplexing gain, `r`): bağımsız uzamsal akış sayısı; hızı artırır. Kapasite yüksek SNR'da `~r·log2(SNR)` gibi büyür.

Çeşitlilik kazancı (diversity gain, `g`): aynı bilgiyi birden çok bağımsız yoldan göndererek hata olasılığını düşürmek; güvenilirliği artırır. Hata olasılığı `~SNR^{−g}` gibi düşer. Maksimum çeşitlilik `g_max = N_t · N_r`.

Zheng-Tse ödünleşimi, bu ikisinin aynı anda ne kadar elde edilebileceğini bir doğru parçasıyla bağlar (`N_t = N_r = N`, yüksek SNR):

```
 g(r) = (N_t − r)(N_r − r),   r = 0, 1, …, min(N_t,N_r)
```

İki uç: `r=0` (hiç çoğullama, tüm antenler aynı bilgiyi taşır) → maksimum çeşitlilik `g = N_t N_r`; `r = min(N_t,N_r)` (tam çoğullama) → çeşitlilik `g = 0`. Aradaki her nokta bir denge: ya hız ya güvenilirlik. Bir sistem hangi noktada çalışacağını kanala ve uygulamaya göre seçer (örn. yüksek hareketlilikte/zayıf kanalda çeşitliliğe, sabit/güçlü kanalda çoğullamaya kayar).

```
 Çeşitlilik g ▲
   N·N ●╲                       r=0:  tüm güç güvenilirliğe
        │ ╲___
        │     ╲___              ara: hız ↔ güvenilirlik dengesi
        │         ╲___
      0 ┼─────────────●──────►   r=min(Nt,Nr): tüm kapasite hıza
        0            min(Nt,Nr)   Çoğullama r
```

### Hüzme-uzayı (beamspace) MIMO

Eleman-uzayında (element space) her antenin ayrı RF zinciri olması Massive MIMO'da (`N` yüzlerce) imkânsız pahalıdır. Hüzme-uzayı dönüşümü, sinyali sabit bir hüzme kümesine (örn. bir DFT hüzme bankası — Butler matrisi) yansıtır. Anahtar gözlem: mmWave kanalları seyrektir (sparse) — enerji yalnızca birkaç yönde (hüzmede) yoğunlaşır. Bu yüzden `N` elemanlı sinyal, yalnızca birkaç baskın hüzme ile temsil edilebilir:

```
 x_beam = U_DFT^H · x_element        (DFT hüzme bankası ile dönüşüm)
```

`x_beam`'in çoğu bileşeni sıfıra yakındır (seyrek); yalnızca enerjili birkaç hüzmeyi seçip (beam selection) onlara RF zinciri ayırmak, `N` zincir yerine `K ≪ N` zincirle çalışmayı sağlar. Bu, hibrit beamforming'in (Kısım 6) ve mmWave Massive MIMO'nun pratik anahtarıdır: kanalın uzamsal seyrekliğini sömürerek donanımı küçültmek. Polyphase kanallaştırıcının (Bölüm 18, Kısım 8) FFT'yle frekans alt-kanalları üretmesi gibi, hüzme-uzayı da DFT'yle uzamsal alt-kanallar (hüzmeler) üretir — aynı dönüşüm, frekans yerine açı ekseninde.

---

<a id="11"></a>
## 11. 5G Massive MIMO ve FR2/mmWave Hüzme Yönetimi

Bölüm 20, 5G NR'nin FR1 (sub-6 GHz) ve FR2 (mmWave, ~24-52+ GHz) ile beamforming/Massive MIMO kullandığını verdi. Burada o hüzmelerin dizi-matematik tarafını ve hüzme yönetiminin nasıl çalıştığını açıyoruz.

### Neden FR2 beamforming'i zorunlu kılar

mmWave'de iki fizik aynı yöne işaret eder. Birincisi yol kaybı: serbest-uzay yol kaybı (Bölüm 1) frekansla artar (`~f²` terim), bu yüzden 28 GHz, 3 GHz'e göre çok daha hızlı zayıflar. İkincisi (telafi edici) anten boyutu: sabit fiziksel alana sığan eleman sayısı `~(boyut/λ)²` ile ölçeklenir; `λ` küçük olduğundan (28 GHz'de ~1 cm) avuç içi büyüklüğünde bir panele yüzlerce eleman sığar. Sonuç: mmWave'de yüksek eleman sayısının ürettiği yüksek dizi kazancı (`~10log10(N)`), yüksek yol kaybını tam telafi eder — ama yalnızca enerji dar bir hüzmede toplanırsa. Bu yüzden FR2 yayını fiziksel olarak beamforming'e mecburdur; geniş, yönsüz yayın menzili olmaz.

### Hüzme yönetimi: tarama, ölçüm, izleme

Dar hüzmelerin bedeli, doğru yöne nişan alma zorunluluğudur. Baz istasyonu (gNB) ile kullanıcı (UE) birbirini "görmek" için bir hüzme yönetim protokolü (beam management) yürütür — kavramsal adımlar:

- Hüzme tarama (beam sweeping): gNB, bir senkronizasyon sinyali bloğunu (SSB) farklı hüzme yönlerinde sırayla yayar; her hüzme uzaydaki bir sektörü tarar. UE her hüzmeyi ölçer.
- Hüzme ölçüm ve raporlama: UE, en güçlü aldığı hüzmeyi (genelde bir RSRP ölçümüyle) belirler ve gNB'ye geri bildirir.
- Hüzme belirleme/iyileştirme (refinement): kaba SSB hüzmesinden daha dar CSI-RS hüzmelerine inilerek en iyi yön keskinleştirilir.
- Hüzme izleme ve kurtarma (tracking / recovery): UE hareket edince ya da hüzme engellenince (mmWave'de bir el bile bloke eder) bağlantı bozulur; hüzme kurtarma yeni bir uygun hüzme bulur.

```
 gNB hüzme tarama (SSB sweep):        UE en iyi hüzmeyi seçer
        ╱ hüzme 0
       ╱── hüzme 1                     UE  ◀── en güçlü: hüzme 2
   gNB ──── hüzme 2  ━━━▶ UE           rapor: "hüzme 2"
       ╲── hüzme 3                     gNB o yöne kilitlenir
        ╲ hüzme 4
```

Bu, Kısım 6'daki faz-dizisi yönlendirmenin protokol katmanıdır: gNB her SSB hüzmesi için farklı bir `w = a(θ_k)` ağırlığı uygular, uzayı tarar, UE'nin raporuna göre `θ_k`'yi seçer. Düzlemsel dizi (Kısım 4) sayesinde tarama hem azimut hem yükselişte yapılır.

> SIGINT bağı: Bu hüzme tarama yapısı, pasif bir gözlemciye bir miktar uzamsal-zamansal bilgi sızdırır: SSB hüzmelerinin periyodik taraması (beam sweep deseni) tespit edilebilir bir imzadır ve hangi hüzmenin ne zaman aktif olduğu, hangi yönde kullanıcı olduğuna dair ipucu verir. Ancak içerik yine şifrelidir; analist yalnızca uzamsal/zamansal deseni gözler. Ayrıntılı 3GPP NR hüzme yönetimi prosedürleri sürümlere göre tanımlıdır ve güncel şartnameden teyit edilmelidir; buradaki anlatım kavramsaldır.

---

<a id="12"></a>
## 12. Faz Tutarlılığı ve Kalibrasyon: Çok-Kanallı Alıcı, Saat/LO Paylaşımı

Bütün dizi matematiği tek bir varsayıma dayanır: kanallar arası faz farkları yalnızca geliş yönünden kaynaklanır, donanımdan değil. Gerçek donanım bu varsayımı doğal olarak sağlamaz; faz tutarlılığı (phase coherence) ve kalibrasyon, dizinin çalışmasının olmazsa olmaz koşuludur. Bölüm 9'daki "her kanal aynı saat/LO ile çalışmalı" şartının matematiksel ve donanımsal açılımı budur.

### Faz tutarlılığının iki koşulu

Bir çok-kanallı alıcıda her kanal kendi ADC'sine, kendi mikserine sahip olabilir. Sorun iki kaynaktan çıkar:

- Örnekleme saati (sample clock): kanallar farklı saatlerle örneklenirse, aralarında bilinmeyen, zamanla kayan bir faz/zaman farkı olur. Çözüm: tek bir ortak saat (shared clock) tüm ADC'leri sürmeli. KrakenSDR'ın beş RTL-SDR'ı tek bir saat kaynağından beslenir.
- Yerel osilatör (LO): her kanalın downconversion mikseri aynı LO'yu kullanmalıdır; aksi halde her kanalın taban-bandı farklı, bilinmeyen bir faz ofsetiyle döner. Çözüm: ortak LO dağıtımı. Ortak saat + ortak LO = faz-tutarlı dizi.

Bunlar sağlansa bile, kanallar açılışta rastgele bir başlangıç faz ofseti (LO'ların belirsiz başlangıç fazı, PLL kilitlenme anı) edinir. Bu yüzden her açılışta bir faz senkronizasyonu/kalibrasyonu gerekir.

### Kanallar arası kalibrasyon

Kalibrasyon, her kanalın bilinmeyen karmaşık kazancını (genlik `g_n` + faz `φ_n`) ölçüp düzeltmektir. Ölçülen örnek, ideal modele bir köşegen kusur matrisi `Γ` ile bağlanır:

```
 x_ölçülen(t) = Γ · ( A s(t) + n(t) ),    Γ = diag( g_1 e^{jφ_1}, …, g_N e^{jφ_N} )
```

`Γ` bilinirse `Γ^{−1}` ile düzeltilir ve ideal model geri gelir. `Γ`'yı kestirmenin yolları:

- Enjekte edilen referans ton (injected reference): bilinen bir ton tüm kanallara aynı anda (eşit faz/genlikle) dağıtılır; her kanalın çıkışındaki faz/genlik farkı doğrudan `g_n e^{jφ_n}`'i verir. KrakenSDR'ın dahili gürültü/referans kaynağı bu mantıkla periyodik kalibrasyon yapar.
- Bilinen yönden gelen kaynak: yönü bilinen güçlü bir kaynak (örn. konumu bilinen bir FM vericisi), beklenen `a(θ_bilinen)` ile karşılaştırılarak kalibrasyon sağlar.
- Öz-kalibrasyon (self-calibration): kaynak yönleri ve kanal kusurları birlikte, kovaryans yapısından yinelemeli kestirilir (daha karmaşık, daha az donanım).

```
 Kalibrasyon zinciri (KrakenSDR tipi):
   ortak saat ──┬─[RTL 0]─┐
   ortak LO  ───┼─[RTL 1]─┤
   ref. ton  ───┼─[RTL 2]─┼─► [faz/genlik ölç] ─► Γ^{−1} ─► kalibre IQ ─► DoA
                ├─[RTL 3]─┤
                └─[RTL 4]─┘
```

### Karşılıklılık (reciprocity)

TDD (zaman bölmeli dupleks) sistemlerde uplink ve downlink aynı frekansı kullandığından, kanal `H` her iki yönde aynıdır (karşılıklılık): `H_downlink = H_uplink^T`. Bu, Massive MIMO için kritik bir kolaylıktır — baz istasyonu, uplink pilotlarından kanalı öğrenip downlink precoding'i (Kısım 9) için kullanabilir; her kullanıcıdan ayrı downlink geri bildirimi gerekmez (`N` antenli sistemde bu geri bildirim aksi halde devasa olurdu). Ancak karşılıklılık yalnızca havadaki kanal için doğrudur; TX ve RX RF zincirleri (PA vs LNA) farklı olduğundan, donanım karşılıklı değildir ve ayrı bir karşılıklılık kalibrasyonu (reciprocity calibration) gerekir.

> Mühendislik sezgisi: Kalibre edilmemiş bir diziyle DoA yapmak, ayarsız bir teleskopla yıldız ölçmek gibidir — her kanalın eklediği bilinmeyen faz, gerçek açıya rastgele bir hata bindirir ve MUSIC tepesini kaydırır ya da çoğaltır. KrakenSDR pratiğinde "önce kalibrasyon" şart koşulmasının nedeni budur; faz-tutarlılık olmadan tüm dizi matematiği geçersizdir. Bir DoA çıktısı tutarsız/gezgin tepeler veriyorsa, ilk şüphe her zaman kalibrasyondur.

---

<a id="13"></a>
## 13. Karşılıklı Kuplaj, Kenar Etkileri ve Gerçek-Dünya Bozulmaları

İdeal dizi modeli, her elemanı bağımsız bir nokta-örnekleyici sayar. Gerçek antenler birbirini etkiler; bu bölüm modeli gerçeğe yaklaştıran bozulmaları sıralar.

### Karşılıklı kuplaj (mutual coupling)

Birbirine yakın (özellikle `d ≤ λ/2`) antenler elektromanyetik olarak etkileşir: bir elemana gelen enerji, komşu elemanların akımlarını da değiştirir. Bu, her elemanın etkin yönlendirme vektörünü ideal `a(θ)`'dan saptırır. Kuplaj, bir kuplaj matrisi `C` (`N×N`, köşegen-dışı terimler etkileşimi taşır) ile modellenir:

```
 a_gerçek(θ) = C · a_ideal(θ)
```

`C` bilinirse (ölçüm ya da elektromanyetik simülasyonla), `C^{−1}` ile düzeltilebilir — kuplaj telafisi (mutual coupling compensation). Düzeltilmezse DoA hatası ve hüzme deseni bozulması olur; özellikle MUSIC gibi süper-çözünürlüklü yöntemler kuplağa hassastır çünkü manifold modelinin doğruluğuna dayanırlar. Kuplaj, eleman aralığı küçüldükçe (yan-lob/çözünürlük için istenen) güçlenir — bir tasarım gerilimi: yakın elemanlar daha iyi uzamsal örnekleme ama daha çok kuplaj.

### Kenar etkileri (edge effects)

Sonlu bir dizide kenardaki elemanlar, ortadakilerden farklı bir elektromanyetik çevre görür (bir tarafında komşu yok). Bu yüzden kenar elemanlarının eleman deseni ve kuplajı ortadakilerden farklıdır; "tüm elemanlar özdeş" varsayımı kenarlarda en çok bozulur. Büyük dizilerde kenar etkisi oransal olarak küçülür; küçük dizilerde (örn. 5 elemanlı KrakenSDR) önemli olabilir. Bazı tasarımlar kenarlara aktif olmayan "kukla" (dummy/parasitic) elemanlar ekleyerek tüm aktif elemanların benzer çevre görmesini sağlar.

### Diğer gerçek-dünya bozulmaları

| Bozulma | Etkisi | Azaltma |
|---|---|---|
| Karşılıklı kuplaj | Manifold sapması, DoA hatası | Kuplaj matrisi `C^{−1}` telafisi |
| Kenar etkileri | Kenar elemanlarda desen farkı | Kukla elemanlar, model düzeltme |
| Kanal kusuru (genlik/faz) | Tüm matematik geçersiz | Kalibrasyon (Kısım 12) |
| Karşılıklı yön belirsizliği (ULA) | `±θ` ayrılamaz | UCA/2B geometri kullan |
| Eleman konum hatası | Manifold yanlış | Hassas montaj, konum kalibrasyonu |
| Geniş-bant beam squint | Hüzme yönü frekansla kayar | True-time-delay, alt-bant BF |
| Korele kaynak (çok-yol) | MUSIC bozulur | Uzamsal yumuşatma (smoothing) |
| Karşılıklı görüş yokluğu (NLOS) | DoA gerçek konuma değil yansımaya işaret eder | Çok-yol ayıklama, çoklu gözlem |

NLOS (non-line-of-sight) bozulması SIGINT'te özellikle sinsidir: bir bina yansımasından gelen güçlü çok-yol bileşeni, DoA'yı verici yönüne değil yansıtıcı yönüne işaret ettirir. Dizi "bir şeyin" yönünü doğru ölçer ama o şey gerçek verici değil hayaletidir. Bu yüzden tek bir DoA ölçümüne asla tam güvenilmez; Bölüm 9'daki gibi birden çok konum/zamanda gözlem birleştirilir.

---

<a id="14"></a>
## 14. Uygulama: Pasif Radar Dizisi, Uydu İzleme, Girişim Haritalama

Üç somut uygulama, bu bölümün matematiğinin sahaya nasıl indiğini gösterir.

### Pasif radar dizisi (PCL)

Bölüm 9, Kısım 9 pasif radarı (Passive Coherent Location) çapraz-belirsizlik fonksiyonuyla tanıttı: fırsatçı bir aydınlatıcının (FM, DVB-T) hedeften yansımasını referansla korele ederek hedef tespiti. Bir dizi bu sisteme uzamsal boyut ekler:

- Yön ayrımı: alıcı dizisi, hedeften yansıyan zayıf eko ile doğrudan gelen güçlü aydınlatıcı sinyalini geliş yönlerine göre ayırır. Doğrudan sinyale (genelde bilinen yönde) bir null çevirmek (Kısım 7), referans sızıntısını bastırıp eko tespitini dramatik kolaylaştırır — pasif radarın en büyük pratik sorunu olan "direct path interference"in uzamsal çözümü budur.
- Hedef DoA: tespit edilen eko üzerinde DoA (MUSIC), hedefin açısını verir; bu açı, gecikme (menzil) ve Doppler (hız) ile birleşince hedefin 3B durumunu (Bölüm 9, Kalman izleme) besler.

```
 Pasif radar dizisi:
   aydınlatıcı (FM) ─────güçlü doğrudan─────▶ [DİZİ]  ─null→ doğrudan bastırıldı
        │                                       │
        └──zayıf eko◀── hedef ──yansıma────────▶│      ─hüzme→ eko DoA ölçülür
                                                 referans + eko korelasyonu
```

### Uydu izleme dizisi

Bir LEO uydusu (Bölüm 11) gökyüzünde hızla hareket eder; mekanik bir çanağın onu takip etmesi motor ve hassas kontrol ister. Bir faz dizisi (Kısım 6) hüzmeyi elektronik olarak yönlendirerek uyduyu hareketli parça olmadan izler — ağırlık vektörü `a(θ(t))`, uydunun yörünge efemerisinden hesaplanan anlık `θ(t)`'ye göre güncellenir. Düzlemsel dizi (Kısım 4) hem azimut hem yükselişte izlemeyi sağlar. Avantajları: çok hızlı yeniden yönlendirme (uydudan uydoya anında geçiş), çok-hüzme (aynı anda birden çok uyduyu izleme), titreşim/aşınma yokluğu. Bu, modern uydu yer terminallerinin (özellikle hareketli platformlarda — gemi, uçak) faz dizilerine yönelmesinin nedenidir.

### Girişim haritalama (interference mapping)

Bir bölgede istenmeyen yayım kaynaklarının (girişim, yetkisiz verici) yönünü ve gücünü haritalamak için dizi tabanlı sistemler kullanılır. Yöntem: dizi sürekli DoA (Kısım 8) çalıştırır, tespit edilen her kaynağın açısını ve gücünü zaman içinde kaydeder. Birden çok dizi istasyonunun açıları kesiştirilirse (triangulation, Bölüm 9) kaynağın konumu çıkar. Bir dizi, aynı anda birden çok kaynağı (`< N` adet) ayrı ayrı izleyebildiği için, yoğun bir spektrumda hangi girişimin nereden geldiğini ayrıştırmak — tek antenli bir tarayıcının asla yapamayacağı bir şey — mümkün olur. Bu, spektrum düzenleyicilerin ve savunma izleme sistemlerinin temel aracıdır; pasiftir, yalnızca yön ve güç ölçülür, içerik çözülmez.

---

<a id="15"></a>
## 15. Alıştırmalar (Yasal, Lab)

> Tümü kendi cihazların, kendi ürettiğin sinyaller ya da yasal/açık yayın bandlarıyla sınırlıdır. İletim (TX) ve yönlü karıştırma içermez; hepsi pasif analiz, simülasyon ya da kendi ürettiğin sinyalin işlenmesidir.

### Alıştırma 1 — İki ve üç elemanlı dizide dizi faktörünü hesapla ve çiz

Amaç: dizi faktörünün eleman sayısı ve aralığıyla nasıl değiştiğini somut görmek (Kısım 2).

Adımlar:
1. `N=2` ve `N=3` için, `d=λ/2` ile eşit-ağırlıklı dizi faktörü `|AF(ψ)| = |sin(Nψ/2)/sin(ψ/2)|`'i `θ ∈ [−90°, +90°]` aralığında (`ψ = π sin θ`) hesapla ve dB ölçekte çiz.
2. Gözlemle: `N=2`'de ana hüzme geniş, tek null var; `N=3`'te ana hüzme daraldı, yan-lob belirdi. İlk-null açısını formülden (`ψ = 2π/N`) hesaplayıp grafikle karşılaştır.

```python
import numpy as np
theta = np.linspace(-90, 90, 1801)
psi = np.pi * np.sin(np.deg2rad(theta))      # d = lambda/2
for N in (2, 3, 8):
    AF = np.abs(np.sin(N*psi/2) / (N*np.sin(psi/2) + 1e-12))
    AF_dB = 20*np.log10(AF + 1e-9)
    # AF_dB'yi theta ekseninde çiz: N büyüdükçe ana hüzme daralır,
    # yan-lob seviyesi ~ -13.3 dB'e oturur
```

### Alıştırma 2 — λ/2 ile λ arasında ızgara lobunun ortaya çıkışı

Amaç: uzamsal aliasing'i (ızgara lobu) gözle görmek (Kısım 3).

Adımlar:
1. `N=4` sabit tutup eleman aralığını `d = 0.5λ`, `0.7λ`, `1.0λ` için dizi faktörünü çiz. Bu kez `ψ = 2π(d/λ)sin θ` kullan (aralığa bağlı).
2. Gözlemle: `d=0.5λ`'da yalnız tek ana hüzme; `d=λ`'da `θ` ekseninin kenarına doğru ikinci bir tam-boy tepe (ızgara lobu) belirir. Bunun fiziksel açı belirsizliği yarattığını (iki farklı `θ` aynı yanıtı verir) doğrula.
3. Hüzmeyi `θ_0 = 30°`'ye yönlendir (`w_n = e^{−jnψ_0}`) ve `d=0.7λ`'da ızgara lobunun görünür bölgeye nasıl kaydığını gör — yönlendirme ızgara lobu riskini artırır (`d ≤ λ/(1+|sin θ_0|)`).

### Alıştırma 3 — KrakenSDR ile uzamsal kalibrasyonu gözle (yasal FM ile)

Amaç: faz-tutarlılığın ve kalibrasyonun DoA'ya etkisini somutlaştırmak (Kısım 8, 12; Bölüm 9, Kısım 8).

Adımlar:
1. KrakenSDR'ı UCA (5 eleman, üretici önerisi aralık) olarak kur. Güçlü, yasal bir FM istasyonuna ayarla (içeriği çözmeden, yalnızca yön).
2. Önce kalibrasyon yapmadan DoA çalıştır: tepe gezgin/tutarsız olur (kanal faz ofsetleri manifoldü bozar). Sonra dahili referans/gürültü kaynağıyla kalibrasyon yap; aynı DoA'nın artık istasyonun gerçek yönüne oturup sabitlendiğini gözle.
3. MUSIC vs Bartlett karşılaştır: aynı veride MUSIC tepesinin Bartlett'ten çok daha keskin olduğunu (süper-çözünürlük) gör. Birden çok istasyon varsa MUSIC'in onları ayrı tepeler olarak çözmesini izle.
4. İleri: diziyi (ya da aracı) birkaç metre kaydırıp ikinci bir açı al; iki açı çizgisinin kesişiminin istasyon konumuna yakınsadığını (triangulation, Bölüm 9) doğrula.

### Alıştırma 4 — İki-eleman dizide bir FM istasyonuna null çevirme kavramı

Amaç: alış-tarafı null steering'in girişimi nasıl bastırdığını anlamak (Kısım 7; Bölüm 13). Tamamen pasif/analitik.

Adımlar (kendi kaydın üzerinde, çevrimdışı işleme):
1. Faz-tutarlı iki-kanallı kayıt al (KrakenSDR'ın iki kanalını kullan), ortamda güçlü bir yasal FM istasyonu olsun.
2. İstasyonun yönünü (`θ_jam`) DoA ile ölç. Sonra ağırlığı `w^H a(θ_jam) = 0` olacak şekilde seç (iki elemanda bu basit: `w = [1, −e^{−jψ_jam}]^T` türünden, `ψ_jam = π sin θ_jam`).
3. Birleştirilmiş çıkış `y = w^H x`'i hesapla ve o istasyonun gücünün düştüğünü (null'a düştüğünü), başka yöndeki bir sinyalin ise korunduğunu gözle. Bu, bir dizinin istenmeyen bir yönü uzamsal olarak "kapatabildiğini" kanıtlar.
4. Kavramsal bağ: `N` elemanlı bir dizinin en çok `N−1` null koyabildiğini hatırla (iki elemanda yalnız bir null). Bunun anti-jam dayanıklılığının donanımsal sınırı olduğunu (Bölüm 13) ilişkilendir.

> Bu dört alıştırma, bölümün dört kritik kavramını (dizi faktörü, ızgara lobu/uzamsal aliasing, kalibrasyon + süper-çözünürlüklü DoA, null steering) somut, yasal ve tamamen kendi sinyalin/simülasyonun üzerinde deneyimletir. Hepsi pasiftir; hiçbir yayın ya da karıştırma yapılmaz.

---

<a id="16"></a>
## 16. Hızlı Referans ve Diğer Bölümler

### Bu bölümün formül kartı

```
 Elemanlar arası faz:   ψ = (2π/λ) d sin θ
 Açıya geri çözme:      θ = arcsin( λ·Δφ / (2π d) )          (Bölüm 9 ile aynı)
 Dizi faktörü (genel):  AF(θ) = Σ w_n e^{ j n ψ }
 ULA eşit ağırlık:      |AF_n(ψ)| = | sin(Nψ/2) / (N sin(ψ/2)) |
 İlk null (ULA):        ψ = 2π/N
 İlk yan-lob (uniform): ≈ −13.26 dB  (N'den bağımsız)
 HPBW (d=λ/2):          ≈ 0.886 λ/(Nd) rad ≈ 102°/N
 Desen çarpımı:         F(θ) = g_eleman(θ) · AF(θ)
 Izgara lobu yok:       d ≤ λ/2  (tam tarama);  d ≤ λ/(1+|sin θ_0|) (θ_0'a)
 Steering vektörü:      a(θ) = [1, e^{jψ}, …, e^{j(N−1)ψ}]^T
 Veri modeli:           x(t) = A s(t) + n(t)
 Kovaryans:             R = A R_s A^H + σ² I;  R̂ = (1/K)Σ x x^H
 Beamformer çıkışı:     y = w^H x;   hüzme deseni B(θ)=w^H a(θ)
 Gecikme-topla ağırlık: w = a(θ_0)/N
 Dizi kazancı:          G = 10 log10(N) dB
 MVDR/Capon ağırlık:    w = R^{-1}a(θ_0) / (a^H(θ_0) R^{-1} a(θ_0))
 LCMV:                  w = R^{-1}C (C^H R^{-1} C)^{-1} f
 Maks. null sayısı:     N−1  (uzamsal serbestlik derecesi)
 LMS güncelleme:        w(t+1)=w(t)+μ x(t) e*(t);  e=d−w^H x
 Bartlett spektrum:     P(θ) = a^H(θ) R a(θ)
 Capon spektrum:        P(θ) = 1 / (a^H(θ) R^{-1} a(θ))
 MUSIC spektrum:        P(θ) = 1 / (a^H(θ) U_n U_n^H a(θ))
 ESPRIT:                θ_i = arcsin( λ·angle(öz_i)/(2π d) )
 MIMO model:            y = H x + n
 MIMO kapasite:         C = log2 det( I + (ρ/N_t) H H^H )
 Eigen-kanal toplamı:   C = Σ log2(1 + (ρ/N_t) σ_i²)
 Su-doldurma:           P_i = max(0, μ − N_t σ²/σ_i²)
 ZF precoder:           W = H^H (H H^H)^{-1}
 DMT (Zheng-Tse):       g(r) = (N_t−r)(N_r−r)
 Maks. çeşitlilik:      g_max = N_t · N_r
 Beamspace dönüşüm:     x_beam = U_DFT^H x_element
 Kalibrasyon kusuru:    x = Γ(A s + n),  Γ=diag(g_n e^{jφ_n});  düzelt: Γ^{-1}
 Kuplaj telafisi:       a_gerçek = C a_ideal;  düzelt: C^{-1}
 TDD karşılıklılık:     H_dl = H_ul^T  (hava kanalı; RF zinciri ayrı kalibre)
```

### Dizi işleme zincirinin tek-bakış haritası

```
 DİZİ (N eleman) → N×[RF+ADC] → faz-tutarlı IQ (ortak saat+LO)
   → KALİBRASYON (Γ^{-1}, kuplaj C^{-1})
   → kovaryans R̂ = (1/K)Σ x x^H
   → ┬─ BEAMFORMING: w (gecikme-topla / MVDR / LCMV / null-steer) → y=w^H x
     └─ DoA: Bartlett / Capon / MUSIC / ESPRIT → θ̂_i (kaynak yönleri)
   → (MIMO ise) H kestirimi → precoding/çoğullama → r paralel akış
   → uygulama: DF / pasif radar / uydu izleme / girişim haritalama
```

### Diğer bölümlerle bağ

Tüm bölümler ve önerilen okuma sırası için indekse bakın: [SIGINT_00 — Başlangıç ve İndeks](SIGINT_00_BASLANGIC_INDEX_VE_YASAL.md).

Doğrudan ilgili bölümler:
- [SIGINT_03 — Antenler, Donanım ve Devre Tasarımı](SIGINT_03_ANTEN_DONANIM_VE_DEVRE_TASARIMI.md): tek anten fiziği — bu bölümün eleman temeli.
- [SIGINT_09 — Yer Tespiti, Yön Bulma ve Takip](SIGINT_09_YER_TESPITI_YON_BULMA_VE_TAKIP.md): AOA, MUSIC/ESPRIT ve DoA'nın uygulama tarafı.
- [SIGINT_18 — Sayısal Sinyal İşleme ve SDR İç Mimarisi](SIGINT_18_DSP_VE_SDR_IC_MIMARI.md): uzamsal işlemenin DSP akrabası (NCO, FFT, channelizer).
- [SIGINT_20 — İleri Hücresel: 4G/5G Güvenlik](SIGINT_20_ILERI_HUCRESEL_4G_5G_GUVENLIK.md): Massive MIMO ve beamforming'in hücresel uygulaması.
- [SIGINT_13 — RF Tehdit Manzarası ve Karşı-Önlemler](SIGINT_13_RF_TEHDIT_VE_KARSI_ONLEMLER.md): null-steering'in anti-jam savunma bağlamı.

> Mühendislik özeti: Bu bölüm, tek antenden uzamsal işlemeye geçişi tamamladı. Tüm dizi matematiği tek bir niceliğe — elemanlar arası faz farkı `ψ = (2π/λ)d sin θ` — ve tek bir vektöre — yönlendirme vektörü `a(θ)` — indirgenir. Bundan sonrası bu vektörlerle lineer cebirdir: ağırlıklarla topla (beamforming), kovaryansı tersine çevir (MVDR), gürültü alt-uzayına dikliği ara (MUSIC), kanal matrisini ayrıştır (MIMO). Bir dizi, uzayda çalışan bir FIR filtresi, bir uzamsal Fourier dönüşümü, bir eşleştirilmiş filtredir — Bölüm 18'in tüm sezgileri buraya, frekans ekseninden açı eksenine taşınır. Frekans ve zaman bir sinyalin "ne" ve "ne zaman"ını verir; dizi "nereden"i ekler ve bu üçüncü eksen, modern SIGINT'in, 5G'nin ve pasif radarın temelidir.

> Yasal hatırlatma: Bu bölümün teknikleri pasif analiz ve eğitim içindir. Beamforming, DoA ve null-steering matematiği evrenseldir; alış tarafında (kendi örneklerini uzamsal işlemek) yapıldığında pasiftir. Verici tarafı yönlü enerji yayma — yönlü karıştırma dâhil — tamamen farklıdır ve çoğu yargı alanında yasa dışıdır; bu bölüm onu kapsamaz ve önermez. Kendi cihazların ve yasal/açık sinyallerle (örn. bir FM istasyonunun yönünü ölçmek) sınırlı kal; ülkenin ve sürümünün mevzuatını teyit et.
