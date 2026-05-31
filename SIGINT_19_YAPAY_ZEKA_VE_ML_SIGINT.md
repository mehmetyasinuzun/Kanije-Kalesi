# SIGINT EL KİTABI — BÖLÜM 19: YAPAY ZEKA VE MAKİNE ÖĞRENMESİ İLE SİNYAL İSTİHBARATI

## Ham IQ'dan Karara: Öğrenen Sistemlerin RF Dünyasındaki Güncel Hâli

> Amaç: Önceki bölümler sinyalin fiziğini (Bölüm 1), donanımı (Bölüm 2), anteni (Bölüm 3), demodülasyonu (Bölüm 4-5), savunmayı (Bölüm 6), ayıklama ve sınıflandırmanın klasik (uzman-kuralı) mühendisliğini (Bölüm 7) ve yön bulmayı (Bölüm 9) verdi. Bu bölüm, o klasik zincirin üzerine binen yeni bir katmanı ele alır: makine öğrenmesi ve derin öğrenme, sinyal istihbaratının hangi adımlarını dönüştürüyor, nerede gerçekten kazandırıyor, nerede tehlikeli biçimde yanıltıyor? Hedefimiz bir kütüphane çağrısı reçetesi değil, mühendislik sezgisidir: bir derin ağ "şu modülasyon BPSK" dediğinde arkasındaki mekanizmayı, güven sınırını ve kırılma noktasını zihninde canlandırabilmen. Bu alanın bir özelliği, çok hızlı değişmesidir; bu yüzden somut sayılara değil, kalıcı ilkelere ve doğrulama reflekslerine ağırlık veriyoruz.

> Yasal ve etik çerçeve: Bu bölüm de serinin geri kalanı gibi anlama, savunma ve spektrum okuryazarlığı amaçlıdır. Anlatılan tüm ML teknikleri tasarım gereği pasif analiz, savunma ve kendi sistemini sertleştirme bağlamında işlenir. Üretken ve karşıt (adversarial) RF başlıkları bilinçli olarak savunma perspektifinden, yani "bir sınıflandırıcı nasıl kandırılır ki ona karşı nasıl savunulur" sorusuyla ele alınır; saldırı reçetesi vermez. Alıştırmalar yalnızca açık veri kümeleri (RadioML gibi) ve kendi cihazlarınla, yasal/açık sinyallerle sınırlıdır. Hiçbir yerde yetkisiz dinleme, içerik çözme veya iletim önerilmez. Bant kullanımı ve kayıt konusunda kendi ülkenin mevzuatını teyit et.

> Doğruluk notu: Bu alan, yayın hızının ve abartının yüksek olduğu bir alandır. Veri kümesi adları, araç adları ve özellikle performans (doğruluk) iddiaları sürümden sürüme, koşuldan koşula değişir. Metinde somut bir ad veya sayı geçtiğinde, kritik kullanımdan önce "teyit edilmeli" uyarısını ciddiye al; verilen rakamlar büyüklük mertebesi sezgisi içindir, kesin değer değildir.

---

## İÇİNDEKİLER

1. [Neden ML: Klasik Sinyal İşlemenin Duvarları](#1)
2. [RF Verisinin Doğası: IQ, Görüntü ve Öğrenilebilir Temsiller](#2)
3. [Otomatik Modülasyon Sınıflandırma (AMC): Klasik Öznitelikten Derin Öğrenmeye](#3)
4. [AMC Mimarileri: CNN, ResNet, Özyineli ve Hibrit Ağlar](#4)
5. [SNR'a Karşı Doğruluk: Performans Eğrisinin Anatomisi](#5)
6. [RF Parmak İzi ve Spesifik Yayıcı Tanıma ML ile (Bölüm 7'ye Köprü)](#6)
7. [Spektrum Algılama ve Anomali Tespiti: Bilişsel Radyo ve Kaçak Yayıcı](#7)
8. [Sinyal Tespiti ve Kaynak Ayrıştırma: Gürültüde Zayıf Sinyal](#8)
9. [Üretken Modeller: Sentez, Veri Artırımı ve Karşıt RF](#9)
10. [Kör Sinyal Analizi ve Otomatik Protokol Çıkarımı](#10)
11. [Konumlandırma ve DF için ML (Bölüm 9'a Köprü)](#11)
12. [Uçta (Edge) ML: Gömülü SDR + Sinir Ağı, Gerçek-Zaman Çıkarım](#12)
13. [Veri Kümeleri, Araçlar ve Eğitim Hattı](#13)
14. [Sınırlar ve Tuzaklar: Dağıtım Kayması, Açıklanabilirlik, Kırılganlık](#14)
15. [Güncel Araştırma Yönelimleri: Transformer, Kendi-Gözetim, Temel Modeller](#15)
16. [Alıştırmalar (Yasal, Lab, Açık Veri)](#16)
17. [Hızlı Referans ve Diğer Bölümler](#17)

---

<a id="1"></a>
## 1. Neden ML: Klasik Sinyal İşlemenin Duvarları

Önceki bölümlerde anlattığımız sinyal işleme zinciri (filtre, FFT, demodülatör, PRI histogramı, kümülant özniteliği) son derece güçlüdür ve onlarca yıldır SIGINT'in bel kemiğidir. O hâlde neden makine öğrenmesine ihtiyaç duyulsun? Cevap, klasik yaklaşımın üç somut duvarına çarpmasındadır. Bu duvarları net görmek, ML'i bir moda olarak değil, belirli bir problem sınıfına verilen mühendislik cevabı olarak konumlandırmak için şarttır.

Birinci duvar, uzman-kuralının ölçeklenmemesidir. Klasik bir modülasyon tanıyıcı veya darbe sınıflandırıcı, bir uzmanın elle tasarladığı eşiklere ve karar ağaçlarına dayanır: "anlık genliğin varyansı şu eşiği aşıyorsa ASK, anlık frekansın momenti buysa FSK." Bu kurallar bilinen, iyi tanımlı dalga formları için mükemmel çalışır. Ama spektrumdaki dalga formu çeşidi sürekli artarken (yeni haberleşme standartları, özel/tescilli protokoller, çevik radarlar), her yeni sınıf için elle yeni kural yazmak insan emeğiyle ölçeklenmez. ML'in vaadi tam buradadır: kuralı elle yazmak yerine, etiketli örneklerden kuralı öğrenmek.

İkinci duvar, bilinmeyen ve gürültülü dalga formlarıdır. Klasik öznitelikler, sinyalin yapısı hakkında varsayımlar yapar (örneğin "taşıyıcı sabittir", "sembol hızı şu aralıktadır"). Gerçek ortamda sinyal çok-yol sönümlemesinden geçmiş, donanım kusurlarıyla bozulmuş, düşük SNR'da gürültüye gömülmüş ve belki de daha önce hiç görülmemiş bir yapıda olabilir. Elle tasarlanmış öznitelik, varsayımları ihlal eden bu koşullarda kırılır. Öğrenen bir model, yeterli ve temsil edici veriyle eğitilirse, bu bozulmaların bir kısmını veriden öğrenip tolere edebilir.

Üçüncü duvar, yoğunluk ve hızdır. Modern bir geniş-bant alıcı saniyede çok büyük hacimde IQ üretir; bu akışta yüzlerce eşzamanlı yayıcıyı insan analist gözüyle takip etmek imkânsızdır. ML, bu akışı otomatik ön-eleme yaparak (neyin ilginç, neyin rutin, neyin anomali olduğunu işaretleyerek) insan dikkatini ölçeklendirir. Burada ML insanı değiştirmez; insanın bakacağı yeri daraltan bir dikkat hunisi kurar.

```
 Klasik zincir (uzman-kuralı)              Öğrenen zincir (veri-güdümlü)
 ───────────────────────────              ───────────────────────────
 IQ → elle öznitelik → elle eşik          IQ → (öğrenilen öznitelik) → öğrenilen karar
       (uzman tasarlar)                          (model veriden çıkarır)
       │                                          │
   bilinen dalga formunda                    görülmüş dağılımda
   kusursuz, yorumlanır                      güçlü, ölçeklenir
       │                                          │
   yeni/bilinmeyen/çok-yol →               dağıtım kayması →
   her sınıf için yeni kural               sessizce yanılır (tehlike)
```

Kritik denge şudur: ML bu üç duvarı aşarken kendi duvarlarını getirir (Bölüm 14). Klasik yöntem yorumlanabilir ve öngörülebilir biçimde başarısız olurken (varsayım ihlal edilince mühendis bunu görür), ML sessizce ve kendinden emin biçimde yanılabilir. Bu yüzden olgun bir SIGINT sistemi ikisini karşıt değil tamamlayıcı kullanır: ML ön-eleme ve zor sınıfları yapar, klasik işleme doğrulama ve fiziksel tutarlılık kontrolü sağlar. Bu bölümün alttan akan tezi budur; "ML her şeyi çözer" değil, "ML belirli darboğazları açar, karşılığında yeni mühendislik sorumlulukları getirir."

> Mühendislik sezgisi: ML'i, klasik sinyal işlemenin yerine geçen bir sihir olarak değil, "elle öznitelik tasarımının ölçeklenmediği yerde, veriden öznitelik ve karar çıkaran bir otomasyon" olarak düşün. Kazanç, esneklik ve ölçek; bedel, açıklanabilirlik kaybı ve dağıtım kaymasına kırılganlıktır. Bu takas her tasarım kararında karşına çıkacak.

---

<a id="2"></a>
## 2. RF Verisinin Doğası: IQ, Görüntü ve Öğrenilebilir Temsiller

Bir derin ağa ne besliyoruz? Bu soru, RF'te makine öğrenmesinin görüntü veya metinden ayrıldığı ilk ve en önemli noktadır. Bilgisayarlı görüde girdi nettir (piksel ızgarası); RF'te ise aynı fiziksel sinyali birbirinden çok farklı temsil eden birçok seçenek vardır ve seçim, modelin neyi öğrenebileceğini doğrudan belirler.

Temel ham temsil, IQ örnekleridir (Bölüm 1 ve 2). Karmaşık taban-bant sinyal, eşfazlı (I) ve dikfazlı (Q) iki gerçek bileşenle temsil edilir. Bir derin ağ açısından bu, iki kanallı bir zaman serisidir: her zaman adımında bir (I, Q) çifti. Ham IQ üzerinde çalışmanın cazibesi, hiçbir bilgi atılmamasıdır; genlik, faz, anlık frekans, hepsi I ve Q içinde gizlidir. Bedeli, ağın bu yapıyı sıfırdan öğrenmesi gerekmesidir.

```
 Aynı sinyalin farklı yüzleri (aynı bilgi, farklı temsil):

 (a) Ham IQ (2 kanallı zaman serisi)
     I: ─╲╱╲╱╲────╲╱╲────   ← ağ önişlemsiz, "her şey burada"
     Q: ╱╲╱╲──╲╱╲╱╲╱──╲

 (b) Genlik/Faz (kutupsal)        |IQ| ve ∠IQ olarak ayrıştırma
     A: ▁▃▅▇▅▃▁▃▅▇   φ: ╱╱╱╲╲╲     bazı kusurlar fazda daha görünür

 (c) Spektrogram / zaman-frekans (görüntü)
     f ▲ ░░▓▓░░          STFT ile 2B görüntü → CNN'e "resim" gibi
       │ ▓▓░░▓▓          zaman-frekans yapısı göze çıkar
       └──────► t

 (d) Çevrimsel / SCF imzası (Bölüm 7)
     gizli periyodiklikler → sembol hızı/taşıyıcı tepe olarak belirir
```

İkinci büyük temsil ailesi, zaman-frekans görüntüleridir. Kısa-zaman Fourier dönüşümü (STFT) ile üretilen spektrogram, sinyali bir 2B görüntüye çevirir; böylece bilgisayarlı görünün olgun evrişimli ağ (CNN) cephaneliği doğrudan uygulanabilir. Bu yaklaşımın gücü, zaman-frekans yapısının (örneğin frekans atlamalı bir sinyalin merdiven deseni, bir chirp'in eğik çizgisi) görsel olarak ayırt edilebilir hâle gelmesidir. Bedeli, STFT'nin zaman-frekans çözünürlük takası (pencere boyu seçimi) ve bazı ince faz bilgisinin görüntüye geçerken kaybolmasıdır.

Üçüncü aile, türetilmiş/uzman temsillerdir: genlik-faz ayrışımı, anlık frekans serisi, çevrimsel öznitelikler (SCF — Bölüm 7), yüksek dereceli kümülantlar. Burada klasik sinyal işleme bilgisini ön-işleme olarak modele enjekte ederiz; ağ sıfırdan öğrenmek yerine bilgiyle ısıtılmış bir girdiyle başlar. Bu, az-veri rejiminde ve fiziksel yorumlanabilirlik istendiğinde değerlidir.

| Temsil | Ne saklar | Güçlü olduğu yer | Zayıf yanı |
|---|---|---|---|
| Ham IQ (I/Q kanalları) | Her şey (genlik+faz) | Maks. bilgi; uçtan-uca öğrenme | Ağ yapıyı sıfırdan öğrenmeli; hizalamaya duyarlı |
| Genlik/Faz (kutupsal) | Aynı bilgi, ayrışık | Bazı kusurlar fazda belirgin | Faz sarması (wrapping) sorunları |
| Spektrogram (STFT) | Zaman-frekans yapısı | CNN doğrudan uygulanır; görsel desen | Çözünürlük takası; ince faz kaybı |
| Çevrimsel/SCF | Gizli periyodiklik | Düşük SNR, sembol hızı kestirimi | Yüksek hesap maliyeti |
| Kümülant/uzman öznitelik | İstatistiksel imza | Az-veri, yorumlanabilir | Öznitelik tasarımı uzmanlık ister |

Bu bölümün ileri başlıklarını okurken bu tabloyu zihninde tut: bir AMC veya parmak izi modelinin başarısı kadar başarısızlığı da çoğu zaman temsil seçimiyle açıklanır. "Model kötü" demeden önce sorulacak ilk soru, "modele neyi besledik ve o temsil aradığımız ayrımı taşıyor mu?" sorusudur.

> Mühendislik sezgisi: Temsil seçimi, RF-ML'in en az konuşulan ama en belirleyici kararıdır. Ham IQ "hiçbir şey atma ama her şeyi öğren" der; spektrogram "yapıyı görselleştir ama bir miktar bilgi feda et" der; uzman öznitelik "fizik bilgini enjekte et ama varsayımlarına bağlan" der. Doğru seçim probleme, veri miktarına ve gecikme/hesap bütçesine bağlıdır; tek doğru yoktur.

---

<a id="3"></a>
## 3. Otomatik Modülasyon Sınıflandırma (AMC): Klasik Öznitelikten Derin Öğrenmeye

Otomatik modülasyon sınıflandırma, RF'te derin öğrenmenin "merhaba dünya"sıdır ve alanın en olgun, en çok kıyaslanmış problemidir. Bölüm 7'de AMC'nin klasik (öznitelik ve istatistik tabanlı) yüzünü ele almıştık: yüksek dereceli kümülantlar, spektral momentler, çevrimsel durağanlık. Burada o temelin üstüne derin öğrenmenin nasıl bindiğini, neyi değiştirdiğini ve neyi değiştirmediğini işliyoruz. Bölüm 7'deki kümülant ve çevrimsel öznitelik anlatımı bu bölümün ön koşuludur; oradaki "neyin sabit, neyin değişken" sezgisi burada da geçerlidir.

Problemi netleştirelim: elimizde kanaldan geçmiş, etiketsiz bir IQ parçası var ve sorumuz "bu hangi modülasyon ailesinden?" (BPSK mi, QPSK mi, 16-QAM mı, GFSK mi, AM mi...). Klasik çözüm, IQ'dan elle öznitelik çıkarıp (örneğin dördüncü derece kümülant) bir sınıflandırıcıya vermekti. Derin öğrenme çözümü, bu elle öznitelik çıkarma adımını ortadan kaldırıp ham IQ'yu (veya spektrogramı) doğrudan ağa verir ve ağın hem özniteliği hem kararı birlikte öğrenmesini bekler. Buna uçtan-uca (end-to-end) öğrenme denir.

```
 İki paradigma yan yana:

 KLASİK (öznitelik mühendisliği)
 IQ ──► [kümülant, spektral moment, SCF] ──► [SVM / rastgele orman] ──► sınıf
        └── uzman tasarlar (sabit) ──┘        └── eğitilir ──┘

 DERİN (uçtan-uca)
 IQ ──► [evrişim katmanları: öznitelik öğrenir] ──► [tam-bağlı: karar] ──► sınıf
        └────────── hepsi birlikte geri-yayılımla eğitilir ──────────┘
```

İki paradigmanın somut karşılaştırması:

| Boyut | Klasik öznitelik + ML | Derin uçtan-uca |
|---|---|---|
| Öznitelik kaynağı | Uzman tasarımı (kümülant, SCF) | Ağ veriden öğrenir |
| Veri ihtiyacı | Görece az | Çok (etiketli) |
| Yorumlanabilirlik | Yüksek (öznitelik fiziksel) | Düşük (kara kutu) |
| Görülmemiş bozulma | Varsayım ihlalinde kırılır | Eğitimde gördüyse tolere eder |
| Düşük SNR davranışı | Çevrimsel öznitelik güçlü | Veriye/mimariye bağlı, değişken |
| Yeni sınıf ekleme | Yeni öznitelik gerekebilir | Yeniden eğitim (veri varsa) |

Burada sık yapılan bir hatayı düzeltelim: derin öğrenme her zaman klasik özniteliği yenmez. Bol ve temsil edici etiketli veri olduğunda ve test dağılımı eğitime benzediğinde, derin ağlar genellikle üstündür. Ama az-veri rejiminde, ağır dağıtım kaymasında veya yorumlanabilirliğin zorunlu olduğu durumlarda, iyi tasarlanmış klasik öznitelikler (özellikle çevrimsel durağanlık tabanlılar, çünkü fiziksel olarak gürültüye dayanıklıdırlar) hâlâ son derece rekabetçidir. Olgun yaklaşım çoğu zaman hibrittir: ağa ham IQ'nun yanında uzman öznitelikleri de besleyip, ağın ikisini birleştirmesine izin vermek.

DeepSig firmasının yayımladığı RadioML veri kümeleri (örneğin sıkça anılan RML2016.10a ve daha büyük/gerçekçi RML2018.01a sürümleri — kesin ad ve içerik teyit edilmeli), bu problemin ortak kıyaslama zeminini kurdu. Bu veri kümeleri, çok sayıda modülasyon türünü geniş bir SNR aralığında sentetik (ve bazı sürümlerde donanımdan geçmiş) IQ parçaları olarak sunar. AMC literatüründeki neredeyse her sayı, "hangi RadioML sürümünde, hangi SNR'da" bağlamıyla okunmalıdır; bu bağlam olmadan doğruluk rakamı anlamsızdır.

> Mühendislik sezgisi: AMC'de derin öğrenmenin asıl kazancı, "elle kümülant tasarlamaktan kurtulmak" değil, görülmemiş bozulmaları (kanal, kusur) veriden öğrenip tolere edebilme potansiyelidir. Ama bu potansiyel ancak eğitim verisi o bozulmaları temsil ediyorsa gerçekleşir. Sentetik temiz veriyle eğitilmiş bir ağ, sahada çakılır; bu, Bölüm 14'teki dağıtım kaymasının en somut örneğidir.

---

<a id="4"></a>
## 4. AMC Mimarileri: CNN, ResNet, Özyineli ve Hibrit Ağlar

AMC için kullanılan sinir ağı mimarilerini tanımak, hangi mimarinin hangi yapıyı yakaladığını anlamak demektir. Burada amaç katman katman bir reçete vermek değil, her mimari ailesinin "RF sinyalinin hangi yönünü görmek için" tasarlandığını kavramaktır. Bölüm 2'deki örnekleme ve Bölüm 1'deki modülasyon bilgisiyle birlikte okunduğunda bu mimariler somutlaşır.

### Evrişimli ağlar (CNN) — yerel desen avcısı

CNN, ham IQ üzerinde 1B evrişim (zaman ekseninde kayan filtreler) veya spektrogram üzerinde 2B evrişim (görüntüde kayan filtreler) uygular. Sezgi şudur: evrişim filtreleri, sinyaldeki yerel desenleri (bir sembol geçişinin şekli, bir frekans sıçraması) öğrenir; ardışık katmanlar bu yerel desenleri giderek daha soyut örüntülere birleştirir. İlk katman "bu küçük dalgacık nasıl görünüyor", son katman "bu örüntü topluluğu hangi modülasyona ait" sorusuna karşılık gelir.

```
 1B-CNN, ham IQ üzerinde (kavramsal katman akışı):

 Girdi: 2 × N   (I ve Q, N örnek)
   │
   ▼  [Conv1D: K filtre, çekirdek boyu w] ── yerel desen (sembol geçişi vb.)
   ▼  [ReLU]                              ── doğrusalsızlık
   ▼  [Conv1D: daha çok filtre]           ── desenlerin bileşimi
   ▼  [Pooling/stride]                    ── boyut indirgeme, kayma toleransı
   ▼  ... (birkaç blok daha) ...
   ▼  [Global pooling / Flatten]
   ▼  [Tam-bağlı katman(lar)]             ── karar
   ▼  [Softmax: M sınıf]                  ── her modülasyon için olasılık
 Çıktı: sınıf olasılıkları
```

### Artık ağlar (ResNet) — derinliği mümkün kılan

Ağ derinleştikçe (daha çok katman) eğitim zorlaşır (kaybolan gradyan sorunu). Artık bağlantılar (residual/skip connections), bir bloğun girdisini çıktısına doğrudan ekleyerek gradyanın derin ağda akmasını sağlar. RF'te ResNet tarzı mimariler, daha derin ve dolayısıyla daha karmaşık öznitelikleri öğrenebildikleri için bazı AMC kıyaslarında düz CNN'i geçer. Sezgi: artık bağlantı, ağa "bu bloğu istersen atla" esnekliği verir; bu, çok katmanlı ağı eğitilebilir kılar.

### Özyineli ağlar (RNN/LSTM/GRU) — zamansal bağımlılık

IQ bir zaman serisidir ve sembollerin sırası anlam taşır. Özyineli ağlar (özellikle uzun-kısa süreli bellek, LSTM ve kapılı özyineli birim, GRU) ardışık örnekler arasındaki zamansal bağımlılığı modeller. Sezgi: RNN, sinyali soldan sağa "okurken" bir bellek taşır; uzun menzilli zamansal yapı (örneğin bir çerçeveleme deseni) bu bellekte yakalanır. Pratikte saf RNN, çok uzun IQ dizilerinde yavaş ve zor eğitilir; bu yüzden sık görülen çözüm hibrittir.

### Hibrit ve dönüştürücü (transformer) yaklaşımlar

En yaygın güçlü tasarımlardan biri CNN+RNN hibritidir (bazı çalışmalarda CLDNN olarak anılır — convolutional, LSTM, deep neural network birleşimi): CNN yerel desenleri çıkarır, RNN bu desenlerin zamansal dizilişini modeller, tam-bağlı katmanlar karar verir. Daha yeni bir akım, dönüştürücü (transformer) tabanlı mimarilerin RF'e taşınmasıdır; öz-dikkat (self-attention) mekanizması, dizideki uzak konumlar arası ilişkileri RNN'den daha iyi yakalayabilir (Bölüm 15). 

| Mimari | Neyi yakalar | Tipik güç | Tipik zayıflık |
|---|---|---|---|
| 1B/2B CNN | Yerel zaman/zaman-frekans deseni | Hızlı, sağlam temel | Uzun menzilli bağımlılık zayıf |
| ResNet (artık) | Daha derin/karmaşık öznitelik | Derinlik → doğruluk | Daha çok hesap/veri |
| RNN (LSTM/GRU) | Zamansal sıra bağımlılığı | Sıralı yapı | Uzun dizide yavaş/zor eğitim |
| CNN+RNN hibrit (CLDNN) | Yerel + zamansal birlikte | Güçlü genel performans | Karmaşık, ayarı zor |
| Transformer | Uzak ilişkiler, ölçeklenir | Büyük veri/model rejiminde güçlü | Veri/hesap açlığı, RF'te genç |

> Mühendislik sezgisi: Mimariyi sinyalin yapısına eşle. Modülasyon ayrımı çoğunlukla yerel istatistikte (takımyıldız geometrisi) saklıysa CNN yeter; ayrım uzun menzilli zamansal/çerçeveleme yapısındaysa RNN veya transformer katkısı belirginleşir. "Daha büyük ağ daima daha iyi" bir yanılgıdır; doğru ağ, aradığın ayrımın saklı olduğu ölçeğe uyan ağdır. Yanlış ölçekteki dev bir ağ, hem aşırı uydurur hem uçta sığmaz (Bölüm 12).

---

<a id="5"></a>
## 5. SNR'a Karşı Doğruluk: Performans Eğrisinin Anatomisi

RF-ML'de tek bir doğruluk sayısı neredeyse hiçbir şey ifade etmez; anlamlı olan, doğruluğun sinyal-gürültü oranına (SNR) karşı çizdiği eğridir. Bu eğrinin şeklini okumak, bir AMC (veya herhangi bir RF sınıflandırıcı) sonucunu değerlendirmenin temel becerisidir. Bölüm 1'deki SNR kavramı ve Bölüm 5'teki demodülasyon eşiği burada doğrudan işe yarar.

Tipik bir AMC doğruluk-SNR eğrisi karakteristik bir S şekli (sigmoid benzeri) çizer: çok düşük SNR'da (örneğin negatif on dB ve altı) sinyal gürültüye gömülüdür, sınıflandırıcı rastgele tahmine yakındır (M sınıf için ~1/M); SNR arttıkça belirli bir eşik bölgesinde doğruluk hızla yükselir; yüksek SNR'da bir tavana oturur. Bu tavan çoğu zaman yüzde yüz değildir, çünkü bazı modülasyon çiftleri yüksek SNR'da bile özünde karışır.

```
 Doğruluk
   ▲
1.0│                          ┌──────────────  ← yüksek-SNR tavanı
   │                        ╱                    (bazı çiftler hâlâ karışır,
   │                      ╱                        tavan < %100 olabilir)
   │                    ╱   ← geçiş bölgesi
0.5│                  ╱        (SNR'a en duyarlı,
   │               ╱           küçük SNR farkı büyük doğruluk farkı)
   │          ╱╱╱
1/M│ ─────────  ← gürültü tabanı (rastgele tahmine yakın)
   └──┴────────┴──────────┴──────────► SNR (dB)
    -20      -10         0        +10
```

Bu eğriden çıkarılacak mühendislik dersleri:

| Eğri bölgesi | Ne anlatır | Tasarım çıkarımı |
|---|---|---|
| Düşük-SNR tabanı | Bilgi yok denecek kadar az | Buradan doğruluk beklemek hayal; ön-işleme/integrasyon gerekir |
| Geçiş bölgesi | Modelin "ayırt etme eşiği" | Sistemin çalışma SNR'ı buraya düşmemeli; en kırılgan bölge |
| Yüksek-SNR tavanı | Modelin teorik tavanı | Tavan < %100 ise sınıflar özünde örtüşüyor; mimari değil problem |

İki modülasyonun karışması, sıkça yüksek SNR'da bile sürer. Klasik örnek, yüksek dereceli QAM aileleri içinde (örneğin 16-QAM ile 64-QAM) veya benzer takımyıldız geometrileri arasındaki karışmadır; bu çiftler az sayıda örnekle veya hafif kanal bozulmasıyla birbirine dönüşebilir. Karışıklık matrisi (confusion matrix), tam olarak hangi çiftlerin karıştığını gösterir ve bir AMC sonucunu teşhis etmenin en bilgilendirici aracıdır: köşegen güçlüyse model iyi, köşegen dışı belirli hücreler doluysa o sınıf çifti özünde zor demektir.

Bir uyarı daha: yayımlanmış doğruluk-SNR eğrileri, neredeyse her zaman belirli bir veri kümesinin (örneğin bir RadioML sürümünün) ürettiği SNR tanımına ve kanal modeline bağlıdır. Farklı veri kümesi, farklı SNR tanımı, farklı sonuç verir. Bir sayı gördüğünde refleksin "hangi veri kümesi, hangi SNR aralığı, hangi modülasyon kümesi?" sorusu olmalı; bağlamsız doğruluk iddiaları teyit edilmeli (bu, alanın en sık tuzaklarından biridir).

> Mühendislik sezgisi: Tek bir doğruluk yüzdesi reklamdır; doğruluk-SNR eğrisi ve karışıklık matrisi mühendisliktir. Bir RF sınıflandırıcıyı değerlendirirken daima "hangi SNR'da ve hangi sınıf çiftinde başarısız oluyor?" diye sor. Sistemin gerçek çalışma SNR'ı geçiş bölgesine düşüyorsa, o sistem kâğıt üzerinde iyi görünüp sahada güvenilmez olur.

---

<a id="6"></a>
## 6. RF Parmak İzi ve Spesifik Yayıcı Tanıma ML ile (Bölüm 7'ye Köprü)

Bölüm 7'de spesifik yayıcı tanımayı (SEI) ve RF parmak izini klasik perspektiften ele aldık: üretim toleranslarının (osilatör kayması, güç yükselteci doğrusalsızlığı, açılış geçici tepkisi) cihaza özgü, istem dışı izler bıraktığını ve bunların ayırt edici olduğunu gördük. Bu bölümde aynı problemin makine öğrenmesiyle nasıl ele alındığına bakıyoruz. Bölüm 7'deki "geçici vs kararlı-durum öznitelik" ayrımı burada da omurgadır; ML, o özniteliği elle çıkarmak yerine veriden öğrenir.

Klasik SEI, uzmanın seçtiği özniteliği (örneğin açılış geçici tepkisinin belirli istatistikleri, faz gürültüsü spektrumunun şekli) ölçer ve eşler. ML tabanlı SEI ise, ham IQ'yu (veya spektrogramı) doğrudan bir derin ağa verir ve ağdan "bu sinyal hangi fiziksel cihazdan?" sorusuna cevap vermesini bekler. Ağ, hangi minik kusurun ayırt edici olduğunu kendi keşfeder. Sezgi: modülasyon sınıflandırmada ağ "ne söyleniyor"u (modülasyon ailesini) öğrenir; parmak izinde aynı ağ mimarisi "kim söylüyor"u (cihaz kimliğini) öğrenir, ama bu kez sınıflar modülasyon türleri değil, tek tek fiziksel cihazlardır.

```
 ML tabanlı RF parmak izi (kavramsal):

 Cihaz 1 ──┐
 Cihaz 2 ──┤  her cihazdan IQ kaydı (aynı protokol/mesaj olsa bile
 Cihaz 3 ──┤  donanım kusuru farklı) ──► [derin ağ] ──► cihaz kimliği
   ...    ──┘                                  │
                              ağ, kusur imzasını (osilatör/PA/transient)
                              veriden öğrenir; sınıflar = cihazlar
```

Bu yaklaşımın güçlü ve kırılgan yanları aynı kökten gelir. Güçlü yan: ağ, bir uzmanın gözden kaçırabileceği ince, çok boyutlu kusur kombinasyonlarını yakalayabilir ve protokol-katmanı kimliğin (adres, sertifika) altında fiziksel-katman bir doğrulama sağlar; bir saldırgan protokol kimliğini taklit etse bile donanım parmak izini taklit etmesi çok zordur. Kırılgan yan: ağ, cihazın kusurunu öğrendiğini sanırken aslında ortamın bir özelliğini (o gün, o kanal, o sıcaklık) öğrenmiş olabilir. Bu, RF parmak izinin en sinsi tuzağıdır.

| Boyut | Klasik SEI (Bölüm 7) | ML tabanlı SEI |
|---|---|---|
| Öznitelik | Uzman seçer (transient, faz gürültüsü) | Ağ öğrenir |
| Yeni cihaz | Yeni öznitelik gerekebilir | Yeniden eğitim/kayıt |
| Sinsi hata | Az (öznitelik fiziksel) | Kanal/ortamı kusur sanma |
| Ölçek | Sınırlı | Çok cihaza ölçeklenir |
| Doğrulama gücü | Fiziksel-katman | Fiziksel-katman (ama kararlılığı kanıtlanmalı) |

ML tabanlı SEI'nin geçerliliği için kritik test, eğitim ve test verisinin farklı oturumlarda/günlerde/koşullarda toplanmış olmasıdır. Aynı oturumda toplanıp rastgele bölünmüş veride neredeyse mükemmel doğruluk elde etmek kolaydır ama yanıltıcıdır; model muhtemelen cihazın kalıcı kusurunu değil, o oturumun geçici imzasını öğrenmiştir. Sağlam bir parmak izi modeli, "farklı gün, farklı kanal, aynı cihaz" testinden geçmelidir. Bu, Bölüm 14'teki dağıtım kaymasının parmak izine özgü, çok somut hâlidir ve sahadaki başarısızlıkların başlıca nedenidir.

Savunma değeri açısından, ML tabanlı RF parmak izi özellikle IoT ve kablosuz güvenlikte klonlanmış/sahte cihaz tespiti için araştırılan bir katmandır (Bölüm 7'deki sahtecilik-tespiti tartışmasının ML uzantısı). Ancak öznitelik kararlılığı sorunu (parmak izinin sıcaklık, yaşlanma ve kanal ile kayması) operasyonel güvenilirliği sınırlar; somut güvenilirlik iddiaları ortama özgüdür ve teyit edilmeli.

> Mühendislik sezgisi: ML tabanlı parmak izinde başarının ölçütü, "aynı veride yüksek doğruluk" değil, "farklı koşulda aynı cihazı hâlâ tanıma"dır. Bir parmak izi modeli gördüğünde ilk soru "eğitim ve test farklı oturumlardan mı?" olmalı. Değilse, raporlanan doğruluk büyük olasılıkla oturum imzasını ezberlemenin ürünüdür; gerçek parmak izi değil.

---

<a id="7"></a>
## 7. Spektrum Algılama ve Anomali Tespiti: Bilişsel Radyo ve Kaçak Yayıcı

Şimdiye kadar "bu sinyal nedir/kimdir" sorusuna odaklandık. Bir adım geriye gidip daha temel bir soru var: "burada bir sinyal var mı, yoksa sadece gürültü mü?" Bu, spektrum algılama (spectrum sensing) problemidir ve hem bilişsel radyonun (cognitive radio) hem de savunma amaçlı kaçak yayıcı tespitinin temelidir. Bölüm 7'deki enerji-tespiti (energy detection) bunun klasik halidir; ML, bu temelin üstüne anomali ve örüntü tabanlı bir katman ekler.

Bilişsel radyo bağlamında spektrum algılama, bir bandın belirli bir an ve yerde boş (kullanılmıyor) mu yoksa işgal mi edildiğini tespit etmektir; amaç, boş "spektrum boşluklarını" (spectrum holes / white spaces) ikincil kullanıcıların fırsatçı biçimde kullanabilmesidir. Klasik yöntemler enerji tespiti (eşik üstü enerji varsa işgal var) ve eşleşmiş süzgeç (matched filter, sinyal yapısı biliniyorsa) gibi tekniklerdir. ML'in katkısı, özellikle düşük SNR'da ve sinyal yapısı tam bilinmediğinde, işgal/boş ayrımını veriden öğrenilmiş bir karar sınırıyla yapabilmesidir; enerji tabanlı eşiğin gürültü dalgalanmasına yenildiği rejimde öğrenen model, sinyalin ince yapısal izlerini (örneğin çevrimsel imza) kullanarak daha iyi ayırt edebilir.

```
 Spektrum işgal haritası (zaman-frekans, ML ile sınıflandırılmış):

 f ▲ ░░░░██████░░░░░░░░░░░░    ██ = işgal (bilinen servis)
   │ ░░░░██████░░░░▓▓▓▓░░░░    ▓▓ = ANOMALİ (beklenmeyen/kaçak yayıcı)
   │ ░░░░░░░░░░░░░░▓▓▓▓░░░░    ░░ = boş (spektrum boşluğu)
   │ ██████░░░░░░░░░░░░██████
   └────────────────────────► t
        ▲ buradaki ▓▓: bu bantta bu zamanda olmaması gereken yayıcı
```

Savunma perspektifinde asıl ilginç olan anomali tespitidir: "bu ortamda normalde ne bulunur" öğrenilir, sonra normalden sapan her şey (yeni, beklenmeyen, kaçak bir yayıcı) işaretlenir. Bunun gücü, anomalinin önceden tanımlanmasının gerekmemesidir; model neyin normal olduğunu öğrenir ve geri kalanı şüpheli sayar. Bu, gözetimsiz veya yarı-gözetimli bir problemdir (etiketli "anomali" örneği genelde yoktur). Tipik yaklaşımlar:

| Yaklaşım | Mantık | RF'te kullanımı |
|---|---|---|
| Otokodlayıcı (autoencoder) | Normali sıkıştırıp geri kurar; anomalide kurma hatası yüksek | Spektrogramı yeniden kuramama → anomali |
| Tek-sınıf modeller (one-class) | Sadece "normal" sınıfının sınırını öğrenir | Sınır dışı → kaçak yayıcı |
| İstatistiksel temel-çizgi | Normal spektrumun istatistiğini öğren, sapmayı ölç | Beklenmeyen güç/bant aktivitesi |
| Tahmin-hatası | Normalin sonraki halini tahmin et, hata büyükse anomali | Zaman serisi spektrum izleme |

Otokodlayıcı sezgisi öğreticidir: ağ, normal spektrumu düşük boyutlu bir koda sıkıştırıp geri kurmayı öğrenir. Normal girdileri iyi kurar (düşük hata); ama daha önce hiç görmediği bir yayıcı geldiğinde onu iyi kuramaz (yüksek kurma hatası). Kurma hatasındaki bu sıçrama, anomali alarmıdır. Bu mekanizma, Kanije Kalesi gibi bir savunma sisteminin "çevremde alışılmadık bir RF kaynağı belirdi mi?" sorusuna ML cevabının iskeletidir.

Uyarı: anomali tespiti, yüksek yanlış-alarm oranıyla maluldür. RF ortamı doğal olarak değişkendir (geçici girişimler, hava koşulları, yeni meşru cihazlar); "anomali" sandığın şey çoğu zaman zararsız bir yeni meşru yayıcıdır. Pratik bir anomali tespit sistemi, ham alarm sayısını değil, alarmları bağlamla (konum, zaman, tekrar) zenginleştirip insana sunmayı hedefler. Tek başına ham anomali skoru, operasyonel olarak gürültülüdür.

> Mühendislik sezgisi: Spektrum algılamada ML'in kazancı, "bilinen sinyali tespit" değil (bunu eşleşmiş süzgeç zaten iyi yapar), "neyin normal olduğunu öğrenip sapmayı yakalama"dır. Ama bu güç, yanlış-alarm bedeliyle gelir. Anomali tespitini bir kesin alarm değil, bir dikkat-yönlendirici olarak tasarla; nihai kararı bağlam ve insan versin.

---

<a id="8"></a>
## 8. Sinyal Tespiti ve Kaynak Ayrıştırma: Gürültüde Zayıf Sinyal

Spektrumda "bir şey var mı" sorusunun ötesinde iki zor problem var: gürültüye gömülü zayıf bir sinyali tespit etmek ve üst üste binmiş birden çok sinyali birbirinden ayırmak (kaynak ayrıştırma). İkisi de klasik sinyal işlemenin (Bölüm 5, 7) zorlandığı yerlerdir ve ML'in giderek daha çok uygulandığı alanlardır.

### Gürültüde zayıf sinyal tespiti

Çok düşük SNR'da, sinyal enerjisi gürültü tabanının altında veya yakınındadır; klasik enerji tespiti çuvallar. Klasik çözümler, sinyalin yapısını kullanan tekniklerdir: eşleşmiş süzgeç (yapı tam biliniyorsa), çevrimsel tespit (gizli periyodikliği kullanır, Bölüm 7) ve uzun süre integrasyonu (sinyal tutarlı, gürültü değil; biriktirince sinyal öne çıkar). ML'in katkısı, sinyalin tam yapısı bilinmese bile, eğitim verisinden o sinyal sınıfının "imzasını" öğrenip gürültüden ayırabilmesidir. Bir derin ağ, zayıf ama yapılı bir sinyali, saf gürültüden, insan gözünün ve basit eşiğin ayıramayacağı SNR'larda ayırt edebilir; çünkü ağ, gürültünün yapısızlığı ile sinyalin yapısı arasındaki ince farkı öğrenmiştir.

```
 Gürültüde zayıf sinyal (kavramsal):

 Ham:   ╲╱╲╱gürültü╲╱╲ +zayıf sinyal+ ╱╲gürültü╲╱   ← gözle ayırt edilemez
                          │
                   [öğrenen tespitçi]
                          │
 Karar: ────────────────[SİNYAL VAR]──────────────   ← yapı imzasından
        ağ, gürültünün yapısızlığı ↔ sinyalin yapısı farkını yakalar
```

### Kaynak ayrıştırma: üst üste binen sinyalleri ayırmak

Bölüm 7'deki deinterleaving, darbeli sinyalleri zaman örüntüsünden ayırıyordu. Haberleşme tarafında benzer bir problem var: aynı bantta, aynı anda yayan birden çok sinyal alıcıda toplanır; bunları geri ayırmak kaynak ayrıştırma (source separation) problemidir. Klasik kör kaynak ayrıştırma (blind source separation) teknikleri (örneğin bağımsız bileşen analizi, ICA) belirli varsayımlar altında çalışır. ML, özellikle derin öğrenme tabanlı ayrıştırma, daha esnek koşullarda (örneğin tek antenli, güçlü örtüşme) sinyalleri ayırmayı öğrenebilir; ses işlemedeki "kokteyl partisi" ayrıştırmasının RF karşılığı gibi düşünülebilir.

| Problem | Klasik araç | ML katkısı | Zorluk |
|---|---|---|---|
| Zayıf sinyal tespiti | Eşleşmiş süzgeç, çevrimsel, integrasyon | Yapı imzasını öğrenip düşük SNR'da ayırma | Yapı bilinmiyorsa eğitim verisi gerek |
| Kaynak ayrıştırma | ICA, dizi işleme | Tek anten/güçlü örtüşmede öğrenilmiş ayrıştırma | Etiketli karışım verisi üretmek zor |
| Karışım + tanıma | Ayır sonra sınıflandır (ardışık) | Uçtan-uca birlikte öğrenme | Hata birikmesi vs ortak optimizasyon |

### Gözetimli mi, gözetimsiz mi?

Bu problemlerde etiket bulmak zordur (gerçek sahada "doğru cevap" yoktur). Bu yüzden üç rejim iç içedir: gözetimli (etiketli sentetik veriyle eğit, sahada uygula — dağıtım kayması riski), gözetimsiz (yapı/küme keşfet, etiket yok) ve yarı-gözetimli (az etiket + çok etiketsiz). Pratikte sık görülen strateji, sentetik karışımlar üretip (bilinen iki sinyali toplayıp "doğru cevabı" elde tutarak) gözetimli eğitmek, sonra gerçek veride sınamaktır; üretilen karışımın gerçekçiliği başarının anahtarıdır.

> Mühendislik sezgisi: Zayıf sinyal tespitinde ML'in sırrı, "gürültünün yapısızlığını" tanımasıdır; gürültü öngörülemezken sinyal (modülasyonun saat tıkırtısı, çevrimsel imza) yapı taşır. Kaynak ayrıştırmada ise asıl darboğaz model değil veridir: gerçekçi etiketli karışım üretmek zordur, ve sentetik karışımla eğitilip gerçek örtüşmede sınanan modeller dağıtım kaymasından en çok etkilenen ailedir.

---

<a id="9"></a>
## 9. Üretken Modeller: Sentez, Veri Artırımı ve Karşıt RF

Şimdiye kadarki başlıklar ayırt edici (discriminative) modellerdi: girdiyi alıp bir etiket veriyorlardı. Üretken (generative) modeller ters yönde çalışır: bir dağılımı öğrenip ondan yeni örnekler üretir. RF'te üretken modellerin üç önemli kullanımı var ve üçü de bu serinin savunma çerçevesiyle dikkatle ele alınmalıdır.

### Sentez ve veri artırımı

RF-ML'in en büyük darboğazı etiketli gerçek veri kıtlığıdır (Bölüm 14). Üretken modeller, özellikle üretken çekişmeli ağlar (GAN — generative adversarial networks), öğrendikleri sinyal dağılımından gerçekçi sentetik IQ örnekleri üreterek eğitim kümesini zenginleştirebilir (veri artırımı / augmentation). Sezgi: bir GAN'da üretici ağ sahte sinyal üretir, ayırt edici ağ "gerçek mi sahte mi" der; ikisi yarışırken üretici giderek daha gerçekçi sinyaller üretmeyi öğrenir. Bu sentetik veri, az-veri sınıflarını dengelemek veya nadir koşulları (belirli SNR, belirli kanal) çoğaltmak için kullanılır.

```
 GAN ile RF veri artırımı (kavramsal):

   gürültü ──► [ÜRETİCİ] ──► sahte IQ ─┐
                                        ├──► [AYIRT EDİCİ] ──► gerçek/sahte?
   gerçek IQ ──────────────────────────┘            │
        ▲                                            │
        └──── üretici, ayırt ediciyi kandıracak ◄────┘
              kadar gerçekçi sinyal üretmeyi öğrenir
              → çıktı: eğitim için ek sentetik veri
```

Uyarı: sentetik veri çift taraflı keskindir. Gerçekçiyse darboğazı açar; gerçeklik dağılımını tam yakalayamazsa modele yanlış bir dünya öğretir ve dağıtım kaymasını gizlice büyütür. Sentetik veriyle eğitilen modelin gerçek veride mutlaka sınanması, üzerinde pazarlık edilemez bir kuraldır.

### Karşıt (adversarial) RF — savunma perspektifi

Burada serinin savunma duruşu kritik. Karşıt örnekler, bir sınıflandırıcıyı kasten yanıltmak için tasarlanmış, insana neredeyse normal görünen ama modele yanlış sınıf dedirten girdilerdir. Görüntü dünyasında iyi bilinen bu olgu RF'e de taşınır: bir sinyale eklenen, enerjisi düşük ve dikkatle hesaplanmış bir bozulma (perturbation), bir AMC veya parmak izi modelini şaşırtabilir. Bunu burada öğrenme amacımız saldırı yapmak değil; kendi savunma sınıflandırıcımızın bu tür manipülasyona ne kadar kırılgan olduğunu anlamak ve sertleştirmektir.

| Kavram | Ne demek | Savunma anlamı |
|---|---|---|
| Karşıt örnek | Modeli yanıltan, az bozulmuş girdi | Sınıflandırıcının kırılganlık yüzeyi |
| Karşıt eğitim | Karşıt örnekleri eğitime katıp dayanıklılık | Modeli sertleştirme yöntemi |
| Aktarılabilirlik | Bir modeli kandıran örnek başkasını da kandırır | Savunmanın zor yanı |
| Sağlamlık değerlendirmesi | Modeli kasten zorlayıp sınır ölçme | Savunma testinin parçası |

Mühendislik dersi şudur: RF güvenliğinde ML kullanan bir savunma katmanı (örneğin parmak izi tabanlı cihaz doğrulama), saf doğruluğun yanında karşıt-sağlamlık açısından da değerlendirilmelidir. Bir saldırgan, modelin karar sınırını bilir veya tahmin ederse, onu yanıltmaya çalışabilir; bu yüzden ML tabanlı savunma, tek başına değil, fiziksel-katman tutarlılık kontrolleri (Bölüm 7'deki çoklu-parametre eşleşmesi, Bölüm 9'daki konum tutarlılığı) ile katmanlı kullanılmalıdır. Tek bir ML kararına dayanan savunma, kırılgan savunmadır.

### Saldırı simülasyonu — kırmızı takım

Üretken ve karşıt teknikler, savunma testinde meşru ve değerli bir rol oynar: kendi sistemine karşı sentetik saldırı senaryoları üretip ("kırmızı takım"), savunmanın nerede çuvalladığını görmek. Burada üretken model, sahada karşılaşılması zor durumları lab ortamında güvenli biçimde çoğaltma aracıdır. Bu, saldırı geliştirmek değil, savunmayı stres-testine sokmaktır ve sorumlu güvenlik mühendisliğinin parçasıdır.

> Mühendislik sezgisi: Üretken modeller RF'te iki yüzlüdür: veri kıtlığını çözen dost (sentez/artırım) ve sınıflandırıcı kırılganlığını ifşa eden uyarıcı (karşıt örnek). Savunma açısından ikisini de tanımak gerekir: birincisini darboğazı açmak için kullan, ikincisini kendi savunmanı sertleştirmek için. Hiçbir ML savunma kararına tek başına güvenme; karşıt-sağlamlık ve fiziksel tutarlılık ile katmanla.

---

<a id="10"></a>
## 10. Kör Sinyal Analizi ve Otomatik Protokol Çıkarımı

Bölüm 5, bilinen protokollerin nasıl çözüleceğini ele aldı. Peki ya protokol bilinmiyorsa? Hiçbir referans dokümanı, hiçbir bilinen yapı olmadan, ham bir sinyalden onun parametrelerini ve yapısını çıkarmaya kör sinyal analizi (blind signal analysis) denir. ML, bu zorlu alanın bazı adımlarını otomatikleştirmeye başlamıştır; ama burada özellikle dürüst olmak gerekir, çünkü bu en çok abartılan başlıklardan biridir.

Kör analiz tipik olarak bir katmanlar zinciridir; her katman bir önceki çıktıyı girdi alır:

```
 Kör sinyal analizi katmanları (kabadan inceye):

 Ham IQ
   │
   ▼ (1) Tespit + parametre: merkez frekans, bant genişliği, başlangıç/bitiş
   ▼ (2) Modülasyon tanıma (AMC, bu bölüm 3-4): hangi modülasyon ailesi
   ▼ (3) Sembol hızı / zamanlama kestirimi (çevrimsel öznitelik, Bölüm 7)
   ▼ (4) Demodülasyon → sembol/bit akışı
   ▼ (5) Çerçeveleme/yapı çıkarımı: tekrar eden desenler, senkron sözcükleri
   ▼ (6) Üst-katman yapı (protokol alanları) — en zor, en spekülatif
   │
   ▼ (kısmi) yapı tahmini
```

ML'in bu zincire katkısı katman katman değişir. İlk katmanlarda (tespit, modülasyon tanıma, sembol hızı) ML olgun ve gerçekten yardımcıdır; bu adımlar bu bölümün önceki başlıklarının (AMC, çevrimsel öznitelik) doğrudan uygulamasıdır. Orta katmanda (sembol/bit akışına inme), ML demodülasyonu kanal koşullarına uyarlamada katkı sağlar. Üst katmanlarda (çerçeveleme, protokol alanı çıkarımı) işler hızla zorlaşır ve büyük ölçüde örüntü-madenciliği, tekrar tespiti ve istatistiksel çıkarımın işi olur; "bir derin ağ ham IQ'dan protokol spesifikasyonu çıkarır" iddiası bugün gerçekçi değildir, en azından genel ve güvenilir biçimde değil (teyit edilmeli).

| Kör analiz adımı | ML olgunluğu | Gerçekçi beklenti |
|---|---|---|
| Tespit + parametre kestirimi | Olgun | Otomatik, güvenilir |
| Modülasyon tanıma | Olgun (bu bölüm) | İyi, SNR'a bağlı |
| Sembol hızı / zamanlama | Olgun (çevrimsel) | İyi |
| Demodülasyon | Gelişmekte | Kanala uyarlamada katkı |
| Çerçeveleme/desen | Kısmi | Tekrar/senkron tespiti; yarı-otomatik |
| Protokol alanı semantiği | Olgunlaşmamış | İnsan-destekli; tam otomasyon gerçekçi değil |

Otomatik protokol çıkarımı (bit akışından yapı keşfi) kendi başına bir araştırma alanıdır ve büyük ölçüde dizilerdeki tekrarları, sabit alanları (senkron sözcükleri), değişen alanları (sayaçlar, yük) istatistiksel olarak ayırt etmeye dayanır. ML burada örüntü keşfini hızlandırabilir ama "anlam"ı (bu alan ne işe yarar) çıkarmak hâlâ büyük ölçüde insan uzmanlığı ve bağlam gerektirir. Bu alandaki dürüst duruş: ML, kör analizin alt katmanlarını ciddi biçimde otomatikleştirir; üst katmanlarda insanın hızlandırıcısıdır, yerine geçeni değil.

> Mühendislik sezgisi: Kör sinyal analizinde ML'in katkısı katmanla ters orantılıdır: alt katmanlarda (tespit, modülasyon, sembol hızı) güçlü ve olgun; üst katmanlarda (protokol semantiği) zayıf ve insan-bağımlı. "Yapay zeka bilinmeyen protokolü otomatik çözer" cümlesini duyduğunda, hangi katmandan bahsedildiğini sor; alt katman gerçek, üst katman çoğunlukla abartıdır (teyit edilmeli).

---

<a id="11"></a>
## 11. Konumlandırma ve DF için ML (Bölüm 9'a Köprü)

Bölüm 9, yön bulmayı (DF) ve yer tespitini klasik perspektiften ele aldı: faz farkından açı kestirimi, çoklu alıcıdan üçgenleme, varış zamanı farkı (TDOA), varış açısı (AOA). Bu bölümde aynı problemlere makine öğrenmesinin nasıl katkı sağladığına bakıyoruz. Bölüm 9'daki çok-yol ve geometri kavramları burada da ön koşuldur.

ML'in konumlandırmaya en belirgin katkısı parmak-izi tabanlı konumlandırmadır (fingerprinting-based localization). Klasik geometrik yöntemler (üçgenleme, TDOA) açık alanda ve net görüş hattında iyi çalışır ama çok-yollu, yansımalı ortamlarda (şehir içi, bina içi) ciddi biçimde bozulur; çünkü sinyal birden çok yoldan gelir ve geometri varsayımları çöker. Parmak-izi yaklaşımı bu sorunu tersine çevirir: çok-yolu bir hata kaynağı olarak değil, konumun imzası olarak kullanır.

```
 Parmak-izi tabanlı konumlandırma mantığı:

 EĞİTİM (offline): bilinen konumlarda RF imzası topla
   konum A ─► imza_A  (RSSI/kanal yanıtı/çok-yol deseni)
   konum B ─► imza_B
   konum C ─► imza_C        → bir "radyo haritası" (radio map) kur

 ÇALIŞMA (online): bilinmeyen konumda imza ölç ─► imza_X
   imza_X'i haritadaki imzalarla eşle ─► en yakın → tahmini konum
   (çok-yol deseni burada DOST: her konumun parmak izini benzersiz kılar)
```

Sezgi şudur: her konumun çevresindeki yansıtıcılar (duvarlar, binalar) o konuma özgü bir çok-yol deseni yaratır; bu desen, alınan sinyal gücü (RSSI), kanal frekans yanıtı veya kanal durum bilgisi (CSI) olarak ölçülebilir. Önce bilinen konumlarda bu imzalardan bir "radyo haritası" toplanır (eğitim); sonra bilinmeyen bir konumun imzası bu haritayla eşleştirilerek konum kestirilir. ML, bu eşleştirmeyi (yüksek boyutlu, doğrusal olmayan imza→konum eşlemesi) öğrenmede güçlüdür.

| Yöntem | Çok-yola tepki | Güçlü olduğu yer | Zayıf yanı |
|---|---|---|---|
| Geometrik (TDOA/AOA, Bölüm 9) | Çok-yol bozar | Açık alan, görüş hattı | Şehir/iç mekânda kötü |
| Parmak-izi (ML) | Çok-yolu imza yapar | İç mekân, yoğun yansıma | Önceden harita gerek; ortam değişince bayatlar |
| Hibrit | Geometri + imza | İkisinin güçlüsü | Karmaşık; veri+model maliyeti |

ML'in ayrıca klasik DF'in kendisini iyileştirme rolü var: açı kestirimini (örneğin dizinin faz örüntüsünden açıya eşleme) bir derin ağa öğretmek, belirli koşullarda klasik öz-uzay yöntemlerinden (Bölüm 9'da değinilen türden) daha gürbüz olabilir, özellikle düşük anlık görüntü sayısı veya bozuk dizi koşullarında. Yine de geometrik yöntemlerin fiziksel garantileri (ML'in vermediği) vardır; olgun sistemler ikisini birleştirir.

Parmak-izi yaklaşımının kritik zayıflığı, radyo haritasının bayatlamasıdır. Ortam değişirse (mobilya, inşaat, hatta yoğunluk farkı), imzalar kayar ve harita geçerliliğini yitirir; bu, Bölüm 14'teki dağıtım kaymasının konumlandırmaya özgü hâlidir. Pratik sistemler haritayı periyodik güncellemek veya değişime dayanıklı imzalar seçmek zorundadır.

> Mühendislik sezgisi: ML konumlandırmada paradigmayı tersine çevirir: klasik yöntem için düşman olan çok-yol, parmak-izi için dosttur (her konumu benzersiz imzalar). Ama bu güç, "ortam sabit kalır" varsayımına bağlıdır; ortam değişince radyo haritası bayatlar. Geometrik yöntemin fiziksel garantisi ile parmak-izinin çok-yol toleransını katmanlamak, en sağlam tasarımdır (Bölüm 9 ile birlikte oku).

---

<a id="12"></a>
## 12. Uçta (Edge) ML: Gömülü SDR + Sinir Ağı, Gerçek-Zaman Çıkarım

Şimdiye kadarki her şey, bol hesap gücü olan bir masaüstü/sunucu varsaydı. Gerçek bir SIGINT veya RF-savunma cihazı (Kanije Kalesi'nin çalıştığı türden gömülü bir platform dâhil) çoğu zaman sınırlı işlemci, sınırlı bellek ve sıkı güç bütçesiyle, üstelik gerçek zamanlı karar vermek zorundadır. Bu, uçta makine öğrenmesi (edge ML) problemidir ve RF'te kendine özgü zorlukları vardır. Bölüm 2'deki SDR donanımı ve Bölüm 4'teki sistem kaynakları bu başlığın zeminidir.

Temel gerilim şudur: en doğru modeller genellikle en büyük ve en yavaş modellerdir, ama uçta hem boyut (bellek/flash) hem gecikme (gerçek-zaman) hem güç (pil/ısı) kısıtı vardır. RF'te bu gerilim daha da serttir, çünkü IQ veri hızı çok yüksektir; saniyede milyonlarca örnek üzerinde gerçek-zamanlı çıkarım yapmak, görüntüde kareler arası çıkarımdan çok daha sıkı bir bütçe demektir.

```
 Uçta ML takas üçgeni (RF bağlamı):

            DOĞRULUK
               ▲
              ╱ ╲
             ╱   ╲      hedef: üçgenin içinde,
            ╱     ╲     uygulamanın kabul ettiği bölgede kalmak
           ╱       ╲    (üçünü birden maksimize edemezsin)
          ╱─────────╲
   GECİKME           BOYUT/GÜÇ
   (gerçek-zaman)    (bellek/pil/ısı)

   IQ hızı yüksek → gecikme kısıtı RF'te özellikle sert
```

Bu gerilimi yönetmek için kullanılan başlıca teknikler:

| Teknik | Ne yapar | RF'te kazanç | Bedel |
|---|---|---|---|
| Niceleme (quantization) | Ağırlıkları düşük bit'e indir (32→8 bit vb.) | Bellek↓, hız↑, güç↓ | Hafif doğruluk kaybı |
| Budama (pruning) | Önemsiz bağlantıları sil | Model küçülür, hızlanır | Aşırı budama doğruluğu bozar |
| Bilgi damıtma (distillation) | Büyük "öğretmen"den küçük "öğrenci" eğit | Küçük model, yakın doğruluk | Eğitim karmaşıklığı |
| Mimari seçimi | Baştan küçük/verimli ağ tasarla | Uca uygun temel | Doğruluk tavanı düşebilir |
| Donanım hızlandırma | FPGA/NPU/DSP üzerinde çıkarım | Düşük gecikme, verim | Donanım+geliştirme maliyeti |

Niceleme sezgisi öğreticidir: bir sinir ağının ağırlıkları genelde 32-bit kayan noktada saklanır; bunları 8-bit tamsayıya (veya daha aza) indirgemek, belleği ve hesap maliyetini büyük ölçüde düşürür ve çoğu zaman doğruluğu yalnızca biraz feda eder. Gömülü RF cihazında bu, "model cihaza sığar mı ve gerçek zamanda yetişir mi" sorusunun çoğu zaman olumlu cevabıdır.

RF'e özgü bir mimari soru, çıkarımı nerede yapacağındır: tüm derin ağı uçta mı koşturmalı, yoksa uçta hafif bir ön-eleme (örneğin "ilginç bir şey var mı") yapıp ağır analizi merkeze mi göndermeli? Bu, klasik bir uç-bulut iş bölümüdür: uçta ucuz, hızlı, kaba bir tetikleyici (sürekli koşar); merkezde pahalı, yavaş, ince bir analiz (yalnızca tetiklendiğinde). Bant genişliği, gecikme ve gizlilik kısıtları bu bölünmeyi şekillendirir. Bir savunma cihazı için "her şeyi uçta yap" gizlilik ve bağımsızlık verir ama doğruluk tavanını düşürür; denge uygulamaya bağlıdır.

Gömülü RF-ML için olgunlaşan araç zinciri (kavramsal): eğitimi sunucuda yap (PyTorch/TensorFlow), modeli nicele ve dışa aktar (örneğin TensorFlow Lite, ONNX Runtime, ya da donanıma özgü derleyiciler — kesin araç ve uyumluluk teyit edilmeli), uçta hafif bir çalışma zamanıyla koştur. SDR tarafında, GNU Radio akış grafiğine bir çıkarım bloğu eklemek (örneğin gr-inference türü bir köprü — teyit edilmeli) yaygın bir entegrasyon kalıbıdır; böylece canlı IQ akışı doğrudan modele beslenir.

> Mühendislik sezgisi: Uçta RF-ML, "en iyi modeli" değil, "kısıt içinde yeterli modeli" arar. IQ hızının yüksekliği gecikme kısıtını RF'te özellikle acımasız kılar. Doğru tasarım çoğu zaman katmanlıdır: uçta ucuz/hızlı tetikleyici, merkezde pahalı/ince analiz. "Cihaza sığmayan mükemmel model", sığan iyi modelden işe yaramazdır.

---

<a id="13"></a>
## 13. Veri Kümeleri, Araçlar ve Eğitim Hattı

RF-ML'de model mimarisi genellikle en kolay kısımdır; asıl iş veri ve hattadır. Bu başlık, bir RF-ML projesinin somut bileşenlerini ve uçtan uca akışını verir. Tüm ad ve sürümler hızla değiştiği için kesin güncel durum teyit edilmeli; burada kalıcı olan, hattanın iskeletidir.

### Veri kümeleri

| Veri kümesi (kavramsal) | İçerik | Tipik kullanım | Not |
|---|---|---|---|
| RadioML (DeepSig; örn. RML2016.10a, 2018.01a) | Çok modülasyon × geniş SNR, sentetik/donanımdan | AMC kıyaslaması | Ad/sürüm/lisans teyit edilmeli |
| Kendi kayıtların (SDR ile) | Gerçek, yerel, etiketli (kontrollü) | Dağıtım kayması testi, gerçekçilik | En değerli "gerçek" sınama |
| Sentetik (GNU Radio üretimi) | İstenen modülasyon/kanal kontrollü | Eğitim artırımı, nadir koşul | Gerçekçilik dağılımı kritik |
| Topluluk/kurumsal setler | Değişken | Probleme özgü | Lisans ve kalite teyit edilmeli |

RadioML, AMC için fiili ortak kıyaslama zeminidir ve alana giriş için en yaygın başlangıç noktasıdır; ama yalnızca onunla eğitip değerlendiren bir model, sentetik dünyaya aşırı uyarlanmış olabilir. Olgun pratik, açık veri kümesiyle başlayıp kendi kaydettiğin gerçek sinyallerle sınamaktır (Bölüm 16 alıştırması tam bunu yaptırır).

### Araçlar

| Katman | Tipik araç (kavramsal) | Rol |
|---|---|---|
| RF yakalama/akış | GNU Radio, SoapySDR + SDR donanımı (Bölüm 2) | IQ üretimi, ön-işleme |
| RF-ML veri/dönüşüm | torchsig türü kütüphaneler (teyit edilmeli) | Sentez, dönüşüm, veri yükleyici |
| Model eğitimi | PyTorch / TensorFlow | Ağ tanımı, eğitim, değerlendirme |
| Uç çıkarım | TF Lite / ONNX Runtime / donanım derleyici | Niceleme, gömülü koşum (Bölüm 12) |
| SDR↔ML köprüsü | gr-inference türü bloklar (teyit edilmeli) | Canlı IQ → model akışı |

torchsig (PyTorch tabanlı RF sinyal kütüphanesi — kesin kapsam ve sürdürülürlük teyit edilmeli) ve GNU Radio + derin öğrenme çatısı kombinasyonu, son dönemde RF-ML hattının yaygın iskeleti olmuştur. gr-inference benzeri köprüler, eğitilmiş bir modeli canlı SDR akışına bağlamayı hedefler. Bu adların güncel durumu, sürdürülürlüğü ve tam yetenekleri projeye başlamadan doğrulanmalıdır; ekosistem hızlı değişir.

### Uçtan uca eğitim hattı

```
 RF-ML eğitim hattı (iskelet):

 (1) VERİ TOPLAMA ───────► açık set (RadioML) + kendi SDR kayıtların
          │
 (2) ETİKETLEME ─────────► kontrollü koşulda kaydet → etiket "bedava"
          │                (gerçek sahada etiket en pahalı kaynak)
          │
 (3) TEMSİL/ÖN-İŞLEME ───► IQ mı, spektrogram mı, öznitelik mi (Bölüm 2)
          │
 (4) BÖLME ──────────────► eğitim/doğrulama/TEST — sızıntısız ayır!
          │                (oturum/cihaz bazında ayır; rastgele değil)
          │
 (5) EĞİTİM + AYAR ──────► mimari, hiperparametre, düzenlileştirme
          │
 (6) DEĞERLENDİRME ──────► doğruluk-SNR eğrisi + karışıklık matrisi (Bölüm 5)
          │                + farklı koşulda test (dağıtım kayması)
          │
 (7) NİCELEME/DIŞA AKTAR ► uca uygunlaştır (Bölüm 12)
          │
 (8) SAHA SINAMASI ──────► gerçek IQ'da doğrula → geri besle (1'e dön)
```

Hattanın en sık ihmal edilen ve en çok zarar veren adımı (4) veri bölmedir. RF'te rastgele bölme sızıntı (data leakage) yaratır: aynı kaydın komşu parçaları hem eğitime hem teste düşerse, model ezberler ve test doğruluğu sahte yükselir. Doğru bölme oturum, cihaz veya zaman bazında yapılır; "eğitimde gördüğüm cihazın/oturumun farklı bir parçası" testte olmamalıdır. Bu kural, Bölüm 6'daki (parmak izi) ve Bölüm 14'teki (dağıtım kayması) tüm uyarıların pratik temelidir.

> Mühendislik sezgisi: RF-ML'de zorluk modelde değil veridedir: etiketli gerçek veri toplamak pahalı, sızıntısız bölmek inceliklidir. Bir RF-ML sonucu duyduğunda ilk iki soru "veri nereden geldi?" ve "eğitim/test nasıl bölündü?" olmalı. Mimari adı etkileyici ama veri hattası kararları başarıyı (ve çoğu sahte başarıyı) belirler.

---

<a id="14"></a>
## 14. Sınırlar ve Tuzaklar: Dağıtım Kayması, Açıklanabilirlik, Kırılganlık

Bu başlık, bölümün en önemlisidir. RF-ML'in vaadini abartmamak ve sahadaki başarısızlıkların nedenlerini önceden tanımak, bir mühendisi heveskârdan ayırır. Buradaki her tuzak, bir sürü etkileyici makale doğruluğunun gerçek dünyada neden buharlaştığını açıklar.

### Dağıtım kayması (distribution shift) — bir numaralı katil

RF-ML'in en yaygın ve en sinsi başarısızlığı budur. Model, eğitim verisinin dağılımında öğrenir; gerçek dünya farklı bir dağılım sunarsa (farklı donanım, farklı kanal, farklı SNR, farklı gürültü, farklı zaman/yer), model sessizce çöker. RF'te bu kayma her yerdedir: sentetik veriyle eğitip gerçeğe uygulamak, bir SDR'la eğitip başkasıyla test etmek, bir günün koşulunda eğitip başka gün uygulamak — hepsi kaymadır. Bu bölüm boyunca tekrar tekrar uyardığımız parmak izi (Bölüm 6), kaynak ayrıştırma (Bölüm 8), konumlandırma (Bölüm 11) başarısızlıklarının ortak kökü budur.

```
 Dağıtım kayması (kavramsal):

 EĞİTİM dağılımı        GERÇEK dağılım
   ╱▔▔▔╲                    ╱▔▔▔╲
  ╱     ╲                  ╱     ╲      model burada öğrendi (sol),
 ╱       ╲    ≠           ╱       ╲     burada uygulanıyor (sağ);
 sentetik/                gerçek/        örtüşmedikleri ölçüde
 tek SDR/                 çok SDR/        sessizce yanılır
 tek gün                  her koşul
```

### Aşırı uydurma (overfitting) ve sahte ilişki (spurious correlation)

Model, sınıfı ayıran gerçek nedeni (modülasyonun yapısı, cihazın kusuru) değil, veri kümesinin bir tesadüfünü öğrenebilir. RF'te klasik örnek: bir cihazın sinyalleri hep aynı gün/kanalda toplanmışsa, model cihazı değil o günün kanal imzasını öğrenir (Bölüm 6). Veya bir modülasyon sınıfı hep belirli bir SNR'da kaydedilmişse, model modülasyonu değil gürültü seviyesini ayırt etmeyi öğrenir. Bu sahte ilişkiler eğitim/testte (aynı dağılımdan) görünmez, sahada patlar.

### Açıklanabilirlik (explainability) eksikliği

Derin ağ "BPSK" der ama nedenini söylemez. SIGINT gibi yüksek-bahisli bir alanda (yanlış atıf ciddi sonuç doğurabilir), "model böyle dedi" yeterli gerekçe değildir. Açıklanabilir ML teknikleri (modelin hangi girdi bölgesine baktığını gösteren yöntemler) yardımcı olur ama RF'te görüntüdeki kadar olgun ve sezgisel değildir. Pratik savunma: ML kararını klasik, yorumlanabilir bir kontrolle (fiziksel parametre tutarlılığı, Bölüm 7) çapraz-doğrulamak.

### Karşıt kırılganlık (adversarial vulnerability)

Bölüm 9'da gördüğümüz gibi, ML sınıflandırıcılar kasten tasarlanmış küçük bozulmalarla yanıltılabilir. Savunma amaçlı bir ML sistemi, bu kırılganlık açısından da değerlendirilmeli; saf doğruluk, sağlamlık garantisi değildir.

### Veri etiketleme zorluğu

Gerçek RF verisini etiketlemek pahalıdır: gerçek sahada "doğru cevabı" kim biliyor? Bu yüzden alan, sentetik veriye veya kontrollü kayıtlara bağımlıdır ki bu da dağıtım kaymasını besler. Etiket kıtlığı, RF-ML'in yapısal darboğazıdır (Bölüm 13).

| Tuzak | Belirti | Önlem |
|---|---|---|
| Dağıtım kayması | Lab harika, saha kötü | Farklı koşulda test; saha verisi; uyarlama |
| Aşırı uydurma/sahte ilişki | Test iyi, gerçek dünya kötü | Sızıntısız bölme; çeşitli koşulda veri |
| Açıklanabilirlik eksikliği | "Neden?" cevapsız | Klasik kontrolle çapraz-doğrula |
| Karşıt kırılganlık | Küçük bozulma → yanlış sınıf | Karşıt eğitim; katmanlı savunma |
| Etiket kıtlığı | Yeterli gerçek etiket yok | Sentez + kontrollü kayıt; yarı-gözetim |

> Mühendislik sezgisi: RF-ML'in en tehlikeli özelliği, sessizce ve kendinden emin biçimde yanılabilmesidir. Klasik yöntem varsayımı ihlal edilince görünür biçimde bozulur; ML, dağıtım kaydığında yüksek güvenle yanlış cevap verir. Bu yüzden altın kural: ML kararını asla tek başına ve sorgusuz kabul etme; farklı koşulda sına, fiziksel tutarlılıkla çapraz-doğrula, ve sahte başarının kaynağı olan veri sızıntısına karşı acımasız ol.

---

<a id="15"></a>
## 15. Güncel Araştırma Yönelimleri: Transformer, Kendi-Gözetim, Temel Modeller

Bu başlık, alanın gittiği yönü kavramsal düzeyde verir. Burada somut iddialardan çok eğilimler ve bunların RF'e neden taşındığı önemlidir; alan hızlı değiştiği için kesin durum ve sonuçlar teyit edilmeli. Amaç, kullanıcının "en güncel" merakını, abartıya kaymadan, sağlam bir çerçeveyle karşılamaktır.

### Dönüştürücüler (transformers) RF'te

Dil ve görüntüde paradigma değiştiren dönüştürücü mimarisi (öz-dikkat tabanlı) RF'e de taşınmaktadır. Çekiciliği, öz-dikkatin bir dizideki uzak konumlar arası ilişkileri (RNN'in zorlandığı uzun menzilli bağımlılıkları) doğrudan modelleyebilmesidir; bu, uzun IQ dizilerinde veya karmaşık çerçeveleme yapılarında avantaj vaat eder (Bölüm 4). Ayrıca dönüştürücüler büyük veri ve büyük model rejiminde iyi ölçeklenir. Karşılığında veri ve hesap açlığı yüksektir ve RF'teki olgunlukları görüntü/dile göre gençtir; "transformer her RF problemini çözer" beklentisi henüz temellenmemiştir (teyit edilmeli).

### Kendi-gözetimli öğrenme (self-supervised learning)

RF-ML'in bir numaralı darboğazı etiketli veri kıtlığıydı (Bölüm 13-14). Kendi-gözetimli öğrenme tam bu darboğaza saldırır: etiketsiz bol IQ verisinden, etiket gerektirmeyen bir ön-görev (pretext task) ile yararlı temsiller öğrenmek. Sezgi: modele "bu sinyalin maskelenmiş kısmını tahmin et" veya "bu iki parça aynı yayından mı" gibi etiketsiz görevler verilir; bunları çözerken model, sonradan az etiketli gerçek görevde (AMC, parmak izi) işe yarayan genel RF temsilleri öğrenir. Bu, etiket pahalı ama ham veri bol olan RF için doğal bir uyumdur ve aktif bir araştırma yönüdür.

### Temel modeller (foundation models) RF için

Dil ve görüntüdeki temel modellerin (büyük, genel-amaçlı, çok veriyle ön-eğitilmiş, sonra göreve uyarlanabilen modeller) RF karşılığı, son dönemin en iddialı ve en spekülatif yönelimidir. Fikir: çok büyük miktarda çeşitli RF verisiyle ön-eğitilmiş genel bir "RF temel modeli" kurmak ve onu küçük veriyle çeşitli alt görevlere (modülasyon, parmak izi, anomali) uyarlamak (ince ayar / fine-tuning). Vaat, her görev için sıfırdan eğitmek yerine bir kez genel öğrenip çok göreve transfer etmektir.

| Yönelim | Çözmeye çalıştığı sorun | RF'te durum (kavramsal) | Dürüst uyarı |
|---|---|---|---|
| Transformer | Uzun menzilli bağımlılık, ölçek | Taşınıyor, gelişmekte | Veri/hesap açlığı; RF'te genç |
| Kendi-gözetim | Etiket kıtlığı | Aktif, umut verici | Hangi ön-görev RF'te en iyi, açık değil |
| Temel modeller | Görev-başı sıfırdan eğitim | İddialı, erken | Çok spekülatif; sonuçlar teyit edilmeli |
| Karşıt-sağlamlık | ML savunmanın kırılganlığı | Olgunlaşıyor | Genel garanti hâlâ zor |

Bu yönelimleri okurken sağlam duruş şudur: hepsi gerçek bir darboğaza (etiket kıtlığı, dağıtım kayması, ölçek) cevap olarak doğmuştur ve kavramsal olarak mantıklıdır; ama RF'teki olgunlukları dil/görüntünün gerisindedir ve özellikle "temel model" söylemi henüz büyük ölçüde vaat aşamasındadır. Bir mühendis olarak bu yönelimleri izle, kavramsal mantığını anla, ama somut performans iddialarını kaynaktan ve kendi koşulunda doğrulamadan kabul etme.

> Mühendislik sezgisi: Güncel yönelimlerin ortak teması, RF-ML'in iki kronik derdine (etiket kıtlığı ve dağıtım kayması) saldırmaktır: kendi-gözetim etiket kıtlığına, temel modeller transfer yoluyla az-veri rejimine, transformer ölçeğe. Bunlar gerçek umutlar ama dil/görüntüdeki olgunluğa erişmiş değiller. "En güncel" demek "en kanıtlanmış" demek değildir; yönelimi öğren, iddiayı doğrula.

---

<a id="16"></a>
## 16. Alıştırmalar (Yasal, Lab, Açık Veri)

Aşağıdaki alıştırmalar yalnızca açık veri kümeleri ve kendi cihazlarınla, yasal/açık sinyallerle yapılır. Hiçbiri yetkisiz dinleme, içerik çözme veya iletim içermez. Amaç, bu bölümün sezgilerini (özellikle dağıtım kayması ve SNR-doğruluk eğrisi) elle deneyimleyerek pekiştirmektir.

### Alıştırma 1 — Açık veriyle basit bir AMC sınıflandırıcı eğit

Hedef: RadioML türü açık bir veri kümesiyle temel bir modülasyon sınıflandırıcı eğitmek ve doğruluk-SNR eğrisini çizmek.

1. Açık bir AMC veri kümesi edin (RadioML; lisans ve indirme koşullarını teyit et).
2. Küçük bir model kur: önce klasik öznitelik + basit sınıflandırıcı (örneğin birkaç kümülant + bir ağaç/SVM), sonra küçük bir 1B-CNN (Bölüm 4).
3. Veriyi sızıntısız böl (eğitim/doğrulama/test; aynı kaydın parçaları aynı kümede kalsın, Bölüm 13).
4. Test setinde her SNR seviyesi için doğruluğu ölç ve doğruluk-SNR eğrisini çiz (Bölüm 5).
5. Beklenen gözlem: S şekilli eğri; düşük SNR'da rastgeleye yakın, geçiş bölgesi, yüksek SNR tavanı. Karışıklık matrisini çıkar; hangi modülasyon çiftleri karışıyor?

### Alıştırma 2 — Kendi kaydınla dağıtım kaymasını gözle (en öğretici)

Hedef: Açık veriyle eğitilmiş modelin, kendi kaydettiğin gerçek sinyalde nasıl bozulduğunu görmek.

1. Alıştırma 1'in modelini hazır tut.
2. Kendi SDR'ınla (Bölüm 2), yasal/açık bir kaynaktan (örneğin kendi ürettiğin bir test sinyali ya da açıkça serbest bir yayın) aynı modülasyon ailesinden gerçek IQ kaydet.
3. Bu gerçek kaydı, açık veriyle eğitilmiş modele ver ve doğruluğu ölç.
4. Beklenen gözlem: doğruluk büyük olasılıkla açık-veri test doğruluğunun belirgin altında. Bu, dağıtım kaymasının (Bölüm 14) bizzat yaşanmış hâlidir. Neden? Donanımın, kanalın, gürültünün eğitim dağılımından farkı.
5. Düşün: Modeli kendi kayıtlarınla bir miktar yeniden eğitince (veya artırınca) fark kapanıyor mu?

### Alıştırma 3 — Bir sınıflandırıcının düşük SNR'da nasıl bozulduğunu gözle

Hedef: Gürültü ekleyerek SNR'ı kontrollü düşürmek ve modelin kırılma eşiğini görmek.

1. Temiz (yüksek SNR) bir test kaydı al (kendi ürettiğin ya da açık veriden).
2. Üzerine kontrollü, artan miktarda gürültü ekleyerek bir SNR merdiveni oluştur (yüksek → düşük).
3. Her SNR seviyesinde modelin doğruluğunu ölç ve eğriyi çiz.
4. Beklenen gözlem: belirli bir eşiğin altında doğruluk hızla çöker (geçiş bölgesi, Bölüm 5). Bu eşik, modelin pratik çalışma sınırıdır.
5. Karşılaştır: klasik çevrimsel öznitelik (Bölüm 7) tabanlı basit bir tespit, aynı düşük SNR'da derin modele göre nasıl davranıyor? (Çevrimsel öznitelik düşük SNR'da sıkça daha dayanıklıdır.)

### Alıştırma 4 — Anomali tespiti sezgisi (gözetimsiz)

Hedef: "Normali öğren, sapmayı yakala" mantığını küçük ölçekte denemek.

1. Belirli bir bandın "normal" spektrogramından bir küme topla (yalnızca rutin, beklenen aktivite).
2. Basit bir otokodlayıcı eğit (normali sıkıştırıp geri kursun, Bölüm 7).
3. Sonra normalde olmayan (kendi ürettiğin yasal bir test yayını) bir kaynağı ortama kat; kurma hatasının sıçradığını gözle.
4. Beklenen gözlem: normal girdilerde düşük kurma hatası, beklenmeyen kaynakta yüksek hata → anomali alarmı. Yanlış-alarm oranını da gözle (Bölüm 7 uyarısı).

> Güvenlik ve yasa hatırlatması: Tüm alıştırmalarda yalnızca kendi ürettiğin sinyalleri ve açıkça serbest yayınları kullan; bant kullanımı, kayıt ve (varsa) test yayını için kendi ülkenin mevzuatını teyit et. Bu alıştırmaların değeri sezgidedir: dağıtım kayması ve SNR eşiği, bir kez elle görülünce bir daha unutulmaz.

---

<a id="17"></a>
## 17. Hızlı Referans ve Diğer Bölümler

### Bu bölümün özünü tek tabloda

| Konu | Çekirdek sezgi | Bağlı bölüm |
|---|---|---|
| Neden ML | Elle öznitelik ölçeklenmediğinde veriden öğren; bedel açıklanabilirlik + kayma | Bölüm 7 |
| Temsil seçimi | Ham IQ / spektrogram / uzman öznitelik — en belirleyici karar | Bölüm 1, 2 |
| AMC | Klasik kümülant/çevrimsel ↔ derin uçtan-uca; hibrit çoğu zaman en iyi | Bölüm 7 |
| Mimari | CNN yerel, RNN zamansal, transformer uzak ilişki; ölçeğe eşle | Bölüm 4 |
| SNR-doğruluk | Tek sayı değil eğri + karışıklık matrisi | Bölüm 5 |
| RF parmak izi (ML) | "Kim söylüyor"u öğren; farklı koşulda test şart | Bölüm 7 |
| Spektrum/anomali | Normali öğren, sapmayı yakala; yanlış-alarm bedeli | Bölüm 7 |
| Zayıf sinyal/ayrıştırma | Gürültünün yapısızlığını tanı; veri darboğaz | Bölüm 5, 7 |
| Üretken/karşıt | Sentez dost, karşıt örnek uyarıcı; savunmayı sertleştir | Bölüm 6 |
| Kör analiz | Alt katman olgun, protokol semantiği insan-bağımlı | Bölüm 5 |
| ML konumlandırma | Çok-yol düşmandan dosta döner; harita bayatlar | Bölüm 9 |
| Edge ML | Doğruluk/gecikme/boyut üçgeni; katmanlı uç-bulut | Bölüm 2, 4 |
| Veri/hatta | Zorluk modelde değil veride; sızıntısız böl | Bölüm 13 |
| Sınırlar | Dağıtım kayması bir numaralı katil; sessizce yanılır | Bölüm 14 |
| Güncel yönelim | Kendi-gözetim/temel model/transformer; yönelimi öğren, iddiayı doğrula | Bölüm 15 |

### Tekrarlayan altın kurallar

- Tek bir doğruluk yüzdesine güvenme; doğruluk-SNR eğrisi ve karışıklık matrisi iste.
- Bir RF-ML sonucunda ilk iki soru: "veri nereden geldi?" ve "eğitim/test nasıl bölündü?"
- ML kararını tek başına kabul etme; klasik/fiziksel tutarlılıkla (Bölüm 7, 9) çapraz-doğrula.
- Parmak izi/konumlandırma modelini farklı koşulda (gün/cihaz/ortam) sına; aynı-koşul doğruluğu yanıltır.
- "En güncel teknik" ile "en kanıtlanmış teknik" aynı şey değildir; somut iddiayı kaynaktan teyit et.

### Diğer bölümlerle ilişki

- Bölüm 1 — Temeller (RF ve modülasyon): IQ, SNR, modülasyon türleri; bu bölümün fiziksel zemini.
- Bölüm 2 — SDR cihazları: IQ üretimi, örnekleme, uç donanım; edge ML'in (Bölüm 12) temeli.
- Bölüm 4-5 — Yazılım/OS ve protokol çözümleme: demodülasyon ve kör analizin (Bölüm 10) klasik tarafı.
- Bölüm 6 — Güvenlik ve savunma: karşıt RF ve OPSEC bağlamı (Bölüm 9'un savunma kökü).
- Bölüm 7 — Disiplinler ve sinyal ayıklama: AMC, SEI/parmak izi ve çevrimsel özniteliğin klasik temeli; bu bölümün doğrudan ön koşulu.
- Bölüm 9 — Yer tespiti, yön bulma ve takip: DF ve konumlandırmanın klasik tarafı; ML konumlandırmanın (Bölüm 11) temeli.

> Kapanış sezgisi: Yapay zeka, SIGINT'i klasik sinyal işlemenin yerine geçerek değil, onun ölçeklenmediği darboğazları açarak dönüştürüyor. Kazanç gerçek (ölçek, esneklik, bilinmeyene tolerans); ama bedel de gerçek (açıklanabilirlik kaybı, dağıtım kaymasına sessiz kırılganlık, etiket kıtlığı). Olgun bir RF mühendisi ML'i bir sihir değil, dikkatle doğrulanması gereken güçlü bir araç olarak kullanır: veriye şüpheyle, sonuca çapraz-kontrolle, "en güncel" iddiaya teyit refleksiyle yaklaşır. Bu bölümün tek cümlelik özü: öğrenen sistemler RF'te kapı açar, ama kapının ardında ne olduğunu hâlâ mühendislik disiplini doğrular.
