# SIGINT EL KİTABI — BÖLÜM 21: TARİH, AKTÖRLER VE ELEKTRONİK HARP

## Bağlam ve Kültür — Room 40'tan Bilişsel Elektronik Harbe, Sayı İstasyonlarından Amatör Radyo Etiğine

> Amaç: Önceki bölümler sinyalin fiziğini, donanımını, ayıklamasını ve tehdit manzarasını verdi. Bu bölüm, tüm o tekniğin içine doğduğu zemini ele alır: SIGINT bir laboratuvar disiplini değil, yüz yılı aşkın bir tarihin, kurumların, doktrinlerin ve kültürlerin ürünüdür. Bir mühendis darbe tekrar aralığını ölçebilir ama o ölçümün neden önemli olduğunu, hangi tarihsel dersten doğduğunu ve hangi etik sınırın içinde durduğunu bilmiyorsa, yalnızca bir alet kullanıcısıdır. Burada hedef, sinyal istihbaratının kültürel ve tarihsel okuryazarlığını kazandırmaktır: bir sinyali duyduğunda arkasındaki yüz yıllık birikimi, bir kurum adını gördüğünde onun açık-kaynak düzeydeki bağlamını ve bir gözetim tartışmasını duyduğunda meselenin iki yakasını tanımak.

> Yasal ve epistemik çerçeve: Bu bölüm tamamen açık kaynak, tarihsel ve kavramsal bilgiye dayanır. Anlatılan her teşkilat, operasyon ve olay, kamuya mal olmuş, kitaplara geçmiş, çoğu zaman onyıllar önce sınıflandırması kaldırılmış bilgidir. Hiçbir gizli kaynak, hiçbir doğrulanmamış komplo iddiası, hiçbir operasyonel reçete içermez. Casus haberleşmesi ve gizli vericiler tarihsel bir olgu olarak anlatılır; kurma talimatı verilmez. Sayı istasyonu ve HF dinleme alıştırmaları yalnızca pasif alımdır (RX) ve dünyanın çoğu yerinde serbesttir; yine de kayıt, çözme ve yayma sınırları için kendi ülkenin mevzuatını teyit et. Tarihsel doğruluk bu bölümün en önemli kalite ölçütüdür: emin olunamayan her tarih, sayı veya iddia "teyit edilmeli" notuyla işaretlenmiştir. Teşkilat ve operasyon iddialarında bilinçli olarak kamuya-mal-olmuş düzeyde kalınmıştır.

---

## İÇİNDEKİLER

1. [Neden Tarih ve Kültür? Tekniğin İçine Doğduğu Zemin](#1)
2. [SIGINT Tarihçesi I: Telgraf, Telsiz ve I. Dünya Savaşı (Room 40, Zimmermann)](#2)
3. [SIGINT Tarihçesi II: II. Dünya Savaşı (Enigma, Bletchley Park, Trafik Analizi, Pasifik)](#3)
4. [SIGINT Tarihçesi III: Soğuk Savaş ve Modern Dönem](#4)
5. [İstihbarat Teşkilatları ve SIGINT (Açık-Kaynak Düzeyde): NSA, GCHQ, Five Eyes](#5)
6. [ECHELON Tartışması ve Snowden İfşaatlarının Kamuya Mal Olmuş Bağlamı](#6)
7. [Sayı İstasyonları (Numbers Stations): HF'de Gizem ve Tek-Yönlü Ajan Bağı](#7)
8. [One-Time Pad ve Casus Radyo Tarihçesi: Burst Transmission, The Thing](#8)
9. [Elektronik Harp Doktrini ve Sistemleri: EA/EP/ES, RWR, Chaff/Flare, SEAD](#9)
10. [ELINT ve Radar İstihbaratı Tarihçesi: Radar Parmak İzi ve Savaş Düzeni](#10)
11. [Sivil ve Ticari SIGINT: Spektrum İzleme, Kaçak Verici Avı, Geolocation](#11)
12. [Türkiye ve Bölge Bağlamı: Spektrum Yönetimi ve Amatör Radyo Kültürü](#12)
13. [Etik ve Mahremiyet: Kitlesel Gözetim, Hukuk ve Pratisyenin Sorumluluğu](#13)
14. [Kültür: Amatör Radyo Topluluğu, DXing, Contesting, ARES, Maker RF](#14)
15. [Alıştırmalar (Yasal: Dinle, Araştır, Belgele)](#15)
16. [Hızlı Referans ve Diğer Bölümler](#16)

---

<a id="1"></a>
## 1. Neden Tarih ve Kültür? Tekniğin İçine Doğduğu Zemin

Bu serinin önceki bölümleri bir soruya cevap verdi: "Sinyal nasıl yakalanır, çözülür ve anlamlandırılır?" Bu bölüm farklı bir soruya cevap verir: "Bu disiplin nereden geldi, kimler şekillendirdi, hangi derslerle bugünkü halini aldı ve hangi etik sınırların içinde durur?" İkisi birbirini tamamlar; çünkü teknik bilgi, tarihsel bağlamından koparıldığında ezbere, etik çerçevesinden koparıldığında tehlikeye dönüşür.

Birkaç somut neden, tarihin neden mühendislik kadar önemli olduğunu gösterir.

Birincisi, bugün kullandığın her kavramın bir doğum hikâyesi vardır. Trafik analizinin (Bölüm 7) neden "içeriği çözmeden istihbarat üretebildiğini" gerçekten anlamak için, onun Enigma kırılamadan önce bile İngiliz analistlere savaş düzeni çıkarttığı II. Dünya Savaşı sahnesini bilmek gerekir. Frekans atlamanın (Bölüm 13) neden hem askeri hem ticari haberleşmenin temeli olduğunu anlamak için, fikrin bir film yıldızı ile bir besteci tarafından bir torpido kontrol patenti olarak ortaya atıldığı tarihsel ironiyi bilmek aydınlatıcıdır.

İkincisi, kurumların adları her gün haberlerde, belgesellerde ve teknik tartışmalarda geçer; bir SIGINT okuryazarının bunları açık-kaynak düzeyde, abartısız ve komplosuz biçimde konumlandırabilmesi gerekir. "NSA şunu yapıyor" türü iddiaların hangi kısmının kamuya mal olmuş belge, hangi kısmının spekülasyon olduğunu ayırt etmek, bu alanın temel hijyenidir.

Üçüncüsü, etik. Sinyal istihbaratı, doğası gereği başkalarının haberleşmesine dokunma potansiyeli taşır. Tarih, bu gücün hem meşru savunma hem de kitlesel gözetim için kullanıldığını gösterir. Bir pratisyenin kendi sorumluluğunu tanımlayabilmesi için, geçmişteki tartışmaların (ECHELON, Snowden ifşaatları, kitlesel veri toplama) iki yakasını da bilmesi gerekir. Bu bölüm bir taraf tutmaz; meseleyi olgusal ve dengeli sunar.

Dördüncüsü, kültür. Sinyal dünyasının en canlı, en açık ve en öğretici yüzü, devlet teşkilatları değil amatör radyo topluluğudur. Bir asır boyunca standartları, kurtarma haberleşmesini, teknik yeniliği ve etiği taşıyan bu gönüllü kültür, bir SIGINT öğrencisinin gerçek "okulu"dur. Bu bölüm, o kapıyı da aralar.

> Mühendislik ve insan: Bir analiz ekranındaki küme, soyut bir renk değildir; arkasında Room 40'ta el yazısıyla telgraf çözen bir dilbilimci, Bletchley'de bombe çeviren bir matematikçi, soğuk bir gecede HF'de sayı sayan bir istasyon ve bir kurtarma operasyonunda morse ile hayat kurtaran bir radyo amatörü vardır. Tekniği bu insanlık tarihiyle birlikte taşımak, onu hem daha derin hem daha sorumlu kılar.

---

<a id="2"></a>
## 2. SIGINT Tarihçesi I: Telgraf, Telsiz ve I. Dünya Savaşı (Room 40, Zimmermann)

### Telgraftan telsize: dinlenebilir hale gelen haberleşme

Sinyal istihbaratının kökü, haberleşmenin fiziksel olarak ele geçirilebilir hale gelmesiyle başlar. Telgraf çağında (19. yüzyıl ortası) mesajlar kablolar üzerinden gitti; bunları okumak için fiziksel erişim (hatta dokunma, kablo kesme) gerekiyordu. Asıl dönüm noktası, 19. yüzyılın sonunda telsiz telgrafın (Marconi ve çağdaşları) ortaya çıkışıdır. Bilgi artık serbest uzaya yayılıyordu; yani bir vericinin menzilindeki herkes, prensipte onu dinleyebilirdi. Haberleşmenin havaya verilmesi, dinlemeyi (interception) ilk kez ölçeklenebilir kıldı. Modern SIGINT'in doğum koşulu budur: yayılan sinyal, herkesin sinyalidir.

Bu, donanma için özellikle kritikti, çünkü gemiler birbirleriyle ve karayla telsizle haberleşmek zorundaydı; başka iletişim yolu yoktu. I. Dünya Savaşı'na girilirken büyük güçler, düşman telsiz trafiğini dinlemenin değerini hızla kavradı.

### Room 40: modern kripto-istihbaratın doğuşu

I. Dünya Savaşı'nın en bilinen ve iyi belgelenmiş SIGINT kurumu, İngiliz Amiralliği'nin kod çözme birimi olan ve oda numarasından adını alan Room 40'tır (Oda 40). 1914'te kurulan bu birim, Alman donanma ve diplomatik telsiz trafiğini dinleyip çözmek üzere dilbilimciler, matematikçiler ve meraklılardan oluşturuldu.

Room 40'ın bir avantajı, savaşın erken döneminde ele geçirilen Alman kod kitaplarıydı (örneğin batırılan/karaya oturan gemilerden ve diğer kaynaklardan elde edilen şifre materyalleri). Kod kitabı, çözme işini muazzam kolaylaştırır; çünkü şifreyi sıfırdan kırmak yerine, ele geçen anahtarla doğrudan okumayı mümkün kılar. Bu, SIGINT'in kalıcı bir dersini erken gösterdi: en güçlü "kriptanaliz", çoğu zaman matematik değil, anahtarın fiziksel olarak ele geçirilmesidir.

### Zimmermann Telgrafı: SIGINT'in tarihi değiştirdiği an

Room 40'ın en ünlü başarısı, 1917 tarihli Zimmermann Telgrafı'dır. Alman Dışişleri Bakanı Arthur Zimmermann'ın Meksika'ya gönderdiği, ABD savaşa girerse Meksika'ya ittifak ve toprak vaat eden şifreli diplomatik mesaj, Room 40 tarafından çözüldü. İngiltere, bu istihbaratı kaynağını (yani Alman kodlarını okuyabildiğini) ele vermeden ABD'ye ulaştırma sorununu ustaca yönetti ve telgrafın içeriği kamuoyuna yansıdığında ABD'nin savaşa girme kararında önemli bir etken oldu.

Zimmermann Telgrafı, SIGINT tarihinin kanonik örneğidir; çünkü iki kalıcı dersi aynı anda taşır. Birincisi, çözülmüş bir tek mesajın stratejik tarihi değiştirebileceği. İkincisi, kaynak korumanın (yani "bunu nasıl bildiğini" gizlemenin) en az çözmenin kendisi kadar önemli olduğu. İstihbaratı kullanırken kaynağı ifşa edersen, düşman kodunu değiştirir ve geleceğin tüm istihbaratını kaybedersin. Bu ikilem (istihbaratı kullanmak ile kaynağı korumak arasındaki gerilim) bütün SIGINT tarihinin sabit temasıdır.

```
 I. Dünya Savaşı SIGINT zinciri (Zimmermann örneği):

 Alman diplomatik   Telsiz/kablo      Room 40           Çözülen          Stratejik
 şifreli mesaj  ──► dinleme       ──► kriptanaliz   ──► içerik        ──► etki (ABD'nin
 (Zimmermann)       (interception)    (+kod kitabı)     + kaynak koruma     savaşa girişine
                                                                            katkı)
```

> Not: Zimmermann Telgrafı'nın çözülmesinin ve ABD'ye ulaştırılmasının ayrıntıları (hangi kopyanın nasıl elde edildiği, kaynağın nasıl gizlendiği) tarihçiler tarafından geniş biçimde belgelenmiştir ancak bazı operasyonel ayrıntılarda kaynaklar farklılaşabilir; ana hatlar kamuya mal olmuş tarihtir, ince ayrıntılar için akademik kaynaklardan teyit edilmeli.

### Trafik analizinin embriyosu ve yön bulma

I. Dünya Savaşı, içerik çözmenin yanında iki tekniği daha olgunlaştırdı. Birincisi, telsiz yön bulma (direction finding, DF): düşman vericisinin yönünü ölçerek konumunu kestirmek. Donanma muharebelerinde gemi telsizlerinin yönünü bulmak, filo hareketlerini izlemenin bir yoluydu. İkincisi, trafik analizinin ilk halleri: mesaj içeriği çözülemese bile, hangi istasyonun ne zaman, ne sıklıkla ve kiminle haberleştiğinden örüntü çıkarmak. Bu iki teknik (Bölüm 3'teki yön bulma ve Bölüm 7'deki trafik analizi), bir asır önce burada kök saldı ve doğası gereği bugün de aynı mantıkla çalışır.

---

<a id="3"></a>
## 3. SIGINT Tarihçesi II: II. Dünya Savaşı (Enigma, Bletchley Park, Trafik Analizi, Pasifik)

II. Dünya Savaşı, sinyal istihbaratını el sanatından endüstriyel ve giderek hesaplamalı bir disipline dönüştürdü. Bu dönemin merkezinde, modern bilgisayarın da doğum sancılarını taşıyan bir hikâye vardır.

### Enigma ve Bletchley Park

Almanya, askeri haberleşmesini Enigma adlı elektromekanik bir şifre makinesiyle koruyordu. Enigma, dönen rotorlar ve bir kablo bağlantı panosu (plugboard) aracılığıyla her harfi değişen bir şekilde şifreliyordu; olası anahtar uzayı astronomik büyüklükteydi ve Almanlar makineyi pratikte kırılamaz sayıyordu.

Enigma'nın kırılması, çok-uluslu ve çok-aşamalı bir başarıdır ve tarihsel adaleti için her halkasını anmak gerekir:

- Temel kırılma çalışması savaştan önce Polonya'da başladı. Polonyalı matematikçiler (en çok anılan isim Marian Rejewski ve ekibi) 1930'larda Enigma'nın yapısını matematiksel olarak çözdü ve "bomba" (bomba kryptologiczna) adı verilen elektromekanik bir yardımcı cihaz geliştirdi. Savaş çıkmadan önce bu bilgiyi İngiliz ve Fransız müttefiklerle paylaştılar; bu paylaşım, sonraki başarının zeminini kurdu.

- Çalışma İngiltere'de, Londra yakınındaki Bletchley Park'ta (Government Code and Cypher School, GC&CS) endüstriyel ölçeğe taşındı. Burada Alan Turing ve Gordon Welchman gibi isimlerin katkısıyla geliştirilen ve Polonya bombasından ilham alıp onu aşan İngiliz bombe makinesi, Enigma anahtarlarını günlük olarak aramak için kullanıldı. Bletchley'in çözdüğü istihbarata Ultra kod adı verildi.

- Daha sonra, Almanların daha karmaşık Lorenz şifre makinesiyle korunan yüksek komuta trafiğini çözmek için, Tommy Flowers liderliğinde Colossus geliştirildi. Colossus, programlanabilirliği sınırlı olsa da elektronik (vakum tüplü) sayısal bir hesaplama makinesiydi ve bilgisayar tarihinin kilometre taşlarından biri sayılır.

```
 Enigma/Ultra zinciri (çok-uluslu):

 Polonya (1930'lar)        İngiltere — Bletchley Park            Ürün
 ──────────────────        ───────────────────────────          ──────────
 Rejewski & ekibi          Turing, Welchman, Flowers...
 Enigma'yı matematiksel    bombe (Enigma) + Colossus (Lorenz)   "Ultra"
 olarak çözer, "bomba"  ─► endüstriyel ölçekte günlük       ─► istihbaratı
 cihazı; müttefiklerle        anahtar arama / çözme               (sahaya etki)
 paylaşır
```

Enigma hikâyesinin kalıcı dersleri: (1) "Kırılamaz" sayılan bir sistem, operasyonel kusurlar (tahmin edilebilir mesaj kalıpları, anahtar yönetim hataları, tekrar eden selamlamalar — kribler) ve devasa hesaplama gücü birleşince kırılabilir. (2) Başarı bir dehanın değil, bir sistemin ürünüdür: binlerce insan, makine, lojistik ve disiplinli kaynak koruma. (3) Yine kaynak koruma: müttefikler Ultra'yı koruyabilmek için, çözdükleri bilgiye dayanan eylemleri çoğu zaman başka bir keşif kaynağı varmış gibi gizlemek (ör. önce bir keşif uçağı göndermek) zorunda kaldılar.

> Not: Enigma'nın kırılmasına ilişkin ana hatlar (Polonya'nın öncülüğü, Bletchley'in rolü, Turing'in katkısı, bombe ve Colossus) çok iyi belgelenmiş tarihtir. Belirli savaşların seyrine Ultra'nın tam katkı payı tarihçiler arasında tartışılır ve nicel iddialar (ör. "savaşı şu kadar kısalttı") tahminîdir; bu tür sayılar için akademik kaynaktan teyit edilmeli.

### Trafik analizi: içerik çözülmeden istihbarat

II. Dünya Savaşı, trafik analizinin (Bölüm 7) altın çağıydı. Enigma'nın çözülemediği dönemlerde veya çözmenin yetişemediği anlarda bile, İngiliz analistler düşman telsiz ağının yapısından muazzam bilgi çıkardı: hangi birim hangi frekansta, hangi çağrı işaretiyle, ne zaman ve kiminle konuşuyor. Çağrı işaretlerinin ve frekansların değişim örüntüsü, yeni birliklerin konuşlanmasını, bir harekâtın hazırlığını ve ağ topolojisini ele veriyordu. Bir radyo sessizliği, yaklaşan bir harekâtın habercisi olabiliyordu. Bu, "şifreleme içeriği korur ama varlık, ritim ve örüntü sızar" ilkesinin (Bölüm 7) tarihsel kanıtıdır.

### Pasifik tiyatrosu ve donanma SIGINT'i

Pasifik'te ABD donanma kriptanalizi, Japon donanma kodları üzerinde çalıştı. Bu çabanın en bilinen sonucu, 1942'deki Midway Muharebesi öncesinde Japon hedefine dair kritik istihbarattır; çözülen trafik, Amerikan donanmasına büyük bir taktik avantaj sağladı. Pasifik aynı zamanda yön bulma ağlarının (HF/DF) düşman gemi ve denizaltı konumlarını izlemekte ne kadar değerli olduğunu gösterdi.

> Not: Midway öncesi kriptanalizin rolü kamuya mal olmuş tarihtir, ancak çözmenin tam kapsamı ve karar süreçlerindeki ağırlığı tarihçilerce nüanslı biçimde ele alınır; "tek başına savaşı kazandırdı" gibi indirgemeci ifadelerden kaçınmak ve kaynaktan teyit etmek doğru olur.

### Atlantik ve denizaltı savaşı

Atlantik'te Alman denizaltılarına (U-boot) karşı verilen mücadelede, hem Enigma'nın deniz versiyonunun çözülmesi hem de yön bulma ağları belirleyici oldu. Denizaltıların telsizle konuşma zorunluluğu (komuta ile koordinasyon için) onları yön bulma ile izlenebilir kıldı. Bu, modern bir OPSEC dersinin kökenidir: haberleşmek zorunda olmak, bir zayıflıktır; her yayın bir konum ve örüntü sızıntısıdır.

---

<a id="4"></a>
## 4. SIGINT Tarihçesi III: Soğuk Savaş ve Modern Dönem

### Soğuk Savaş: ölçeğin ve ELINT'in çağı

II. Dünya Savaşı'ndan sonra sinyal istihbaratı kurumsallaştı ve devasa bir ölçeğe ulaştı. Soğuk Savaş'ın iki bloğu, birbirinin haberleşmesini, radarlarını ve füze testlerini sürekli izledi. Bu dönemin getirdiği yenilikler:

ELINT'in yükselişi: Radarın savaşta yaygınlaşmasıyla (Bölüm 10), düşman radar ve hava savunma sistemlerinin parametrelerini çıkarmak (TechELINT) ve konuşlanmalarını izlemek (OpELINT) stratejik öncelik oldu. Uçaklar, gemiler ve sonradan uydular, "elektronik savaş düzeni" (electronic order of battle) çıkarmak için radar yayınlarını topladı.

Hava ve uzay platformları: Yüksek irtifa keşif uçakları ve giderek istihbarat uyduları, SIGINT toplamayı coğrafi sınırların ötesine taşıdı. Uzaydan sinyal toplama, erişilemeyen bölgelerin haberleşme ve radar yayınlarını dinlemenin yolu oldu.

Sınır ve risk: Soğuk Savaş, SIGINT toplamanın tehlikelerini de gösterdi; keşif uçaklarının düşürülmesi, dinleme gemilerine yönelik olaylar gibi vakalar, bu faaliyetin politik ve insani bedellerini hatırlattı. (Belirli olayların ayrıntıları için tarihsel kaynaklardan teyit edilmeli.)

### Modern dönem: dijital, paketlenmiş ve şifreli dünya

Soğuk Savaş'ın sonundan bugüne, haberleşmenin doğası kökten değişti ve SIGINT'i yeni bir zemine taşıdı:

| Dönüşüm | Eski dünya | Yeni dünya | SIGINT'e etkisi |
|---|---|---|---|
| Taşıyıcı | Analog telsiz, devre anahtarlamalı | Paket anahtarlamalı, IP, fiber | Hedef artık "frekans" değil, akış/oturum |
| İçerik | Açık veya zayıf şifreli | Yaygın güçlü şifreleme (TLS, uçtan uca) | İçerik çoğu zaman okunamaz; meta-veri öne çıkar |
| Hacim | Sınırlı kanal sayısı | Devasa veri akışları | Toplama değil, filtreleme/seçim sorunu |
| Taşınabilirlik | Sabit istasyonlar | Mobil cihaz, her yerde sinyal | Yer tespiti ve cihaz ilişkilendirme kritik |

Modern SIGINT'in büyük teması, "interception" sorununun yer değiştirmesidir. Geçmişte zorluk sinyali yakalamaktı; bugün sinyal her yerde ve şifreli, zorluk ise devasa akış içinden ilgiliyi seçmek (Bölüm 7'deki TCPED'in Tasking ve Processing adımları) ve şifreleme içeriği kapattığında meta-veriden (trafik analizi) anlam çıkarmaktır. Bu yüzden modern dönem, trafik analizini ve yer tespitini, içerik çözmenin önüne geçecek kadar değerli kıldı.

```
 SIGINT'in ağırlık merkezinin kayması (kavramsal):

 I. DS / II. DS        Soğuk Savaş            Modern
 ───────────────       ───────────────        ───────────────
 İçerik çözme          ELINT + ölçek          Meta-veri + seçim
 (kriptanaliz)         (radar/uzay)           (şifreli çağ)
 "yakalayabilir        "her şeyi              "akıştan
  miyim?"               topla"                 ilgiliyi seç + 
                                               meta-veriden anla"
```

> Not: Modern dönem teşkilat faaliyetlerine ilişkin somut iddialar büyük ölçüde sınıflandırılmıştır; bu başlık, kamuya açık teknoloji eğilimleri (paketleşme, yaygın şifreleme, mobilite) üzerinden kavramsal bir çerçeve verir, belirli operasyon iddiaları içermez.

---

<a id="5"></a>
## 5. İstihbarat Teşkilatları ve SIGINT (Açık-Kaynak Düzeyde): NSA, GCHQ, Five Eyes

Bu başlık bilinçli olarak açık-kaynak düzeyinde kalır: yalnızca kamuya mal olmuş, kurumların kendi web sitelerinde veya yaygın referans kaynaklarında bulunan bilgiyi aktarır. Amaç, bir SIGINT okuryazarının bu kurum adlarını duyduğunda onları doğru kategoriye koyabilmesidir; operasyonel iddia veya spekülasyon değil.

### Başlıca SIGINT kurumları (kamuya açık tanım düzeyinde)

| Kurum | Ülke | Kamuya açık tanım | Tarihsel kök |
|---|---|---|---|
| NSA | ABD | National Security Agency; ABD'nin sinyal istihbaratı ve haberleşme güvenliği (bilgi güvencesi) kurumu | 1952'de kuruldu; II. DS dönemi kriptoloji birimlerinin halefi |
| GCHQ | Birleşik Krallık | Government Communications Headquarters; İngiltere'nin sinyal istihbaratı ve siber güvenlik kurumu | Bletchley Park'ın halefi GC&CS'in devamı |
| CSE | Kanada | Communications Security Establishment | — |
| ASD | Avustralya | Australian Signals Directorate | — |
| GCSB | Yeni Zelanda | Government Communications Security Bureau | — |

Bu kurumların ikili (dual) bir misyonu olduğu kamuya açık biçimde belirtilir: hem yabancı sinyal istihbaratı toplamak hem de kendi ulusal haberleşmesini ve sistemlerini korumak (bilgi güvencesi / siber güvenlik). İkinci misyon, bu serinin savunma duruşuyla doğrudan örtüşür: aynı bilgi, hem dinlemeyi hem savunmayı besler.

### Five Eyes (Beş Göz) ittifakı

Five Eyes, II. Dünya Savaşı sonrası imzalanan ve kökleri UKUSA Anlaşması'na (1946 civarı, ABD-İngiltere sinyal istihbaratı işbirliği) dayanan, beş İngilizce konuşan ülke (ABD, Birleşik Krallık, Kanada, Avustralya, Yeni Zelanda) arasındaki istihbarat paylaşım ittifakının yaygın adıdır. UKUSA Anlaşması'nın varlığı uzun süre gizliydi; ilgili belgelerin bir kısmı sonradan kamuya açıklandı (deklasifiye edildi). İttifakın temel mantığı, coğrafi sorumluluk paylaşımı ve toplanan sinyal istihbaratının üye ülkeler arasında paylaşılmasıdır.

```
 Five Eyes (kamuya açık çerçeve):

        ┌─────────── UKUSA kökü (≈1946) ───────────┐
        │                                          │
      ABD ── BK ── Kanada ── Avustralya ── Yeni Zelanda
       │     │       │           │              │
       └─────┴───────┴───────────┴──────────────┘
              ortak sinyal istihbaratı işbirliği
              (coğrafi paylaşım + veri paylaşımı)

 Genişletilmiş halkalar (kamuoyunda anılır): "Nine Eyes",
 "Fourteen Eyes" gibi daha geniş paylaşım düzenlemeleri —
 üyelik ve kapsam iddiaları için kaynaktan teyit edilmeli.
```

> Not: Five Eyes ve UKUSA'nın varlığı ile beş üye ülke kamuya açık ve belgelenmiş bilgidir. İttifakın işleyişine, kapsamına ve "Nine/Fourteen Eyes" gibi genişletilmiş düzenlemelere dair ayrıntılı iddialar büyük ölçüde gazetecilik ve sızıntı kaynaklıdır; bunların tam kapsamı resmî olarak doğrulanmamıştır ve kaynaktan teyit edilmelidir. Burada yalnızca ittifakın varlığı ve temel mantığı, kamuya-mal-olmuş düzeyde aktarılmıştır.

### Bu kitabın duruşu

Bu seri, teşkilatların operasyonel faaliyetlerini ne över ne de spekülasyonla suçlar. Bir SIGINT öğrencisi için doğru tutum şudur: bu kurumların var olduğunu, ikili (toplama + koruma) misyon taşıdıklarını ve sinyal istihbaratının devlet düzeyinde nasıl kurumsallaştığını bilmek; ama belirli operasyon iddialarını kamuya açık, doğrulanmış kaynaklarla sınırlı tutmak ve komplo anlatılarından uzak durmak. Bir sonraki başlık (ECHELON ve Snowden), tam da bu "kamuya mal olmuş bağlam" çizgisinde işlenecektir.

---

<a id="6"></a>
## 6. ECHELON Tartışması ve Snowden İfşaatlarının Kamuya Mal Olmuş Bağlamı

Bu başlık, kamuoyunda en çok tartışılan iki SIGINT olgusunu, kasıtlı bir disiplinle ele alır: yalnızca kamuya mal olmuş, gazetecilik ve resmî soruşturmalarla yaygınlaşmış çerçeveyi aktarır; doğrulanmamış ayrıntılara veya operasyonel reçeteye girmez. Vurgu, "bu araç ve kabiliyetlerin var olduğu" bilgisinin kamusallaşmasının kendisindedir; nasıl kullanıldığına dair gizli ayrıntıda değil.

### ECHELON: bir kavramın kamusallaşması

ECHELON, 1990'larda ve 2000'lerin başında kamuoyu ve parlamento gündemine gelen, Five Eyes ittifakına atfedilen küresel bir sinyal istihbaratı toplama ağı iddiasının yaygın adıdır. Tartışma, özellikle Avrupa Parlamentosu'nun konuyu ele alan bir raporuyla (2001 dönemi) kurumsal bir zemine taşındı; rapor, böyle bir kapasitenin var olabileceği ve mahremiyet ile ekonomik casusluk açısından kaygı doğurduğu çerçevesinde tartışma yürüttü.

ECHELON tartışmasının SIGINT okuryazarlığı açısından önemi, teknik ayrıntısından çok şu noktadadır: kitlesel haberleşme toplama kapasitesinin var olabileceği fikri, ilk kez geniş bir kamuoyu ve yasama tartışmasının konusu oldu. Bu, "gözetim ve özgürlük dengesi" tartışmasının (Bölüm 13) modern başlangıç noktalarından biridir.

> Not: ECHELON'a atfedilen sistemin tam mimarisi, kapsamı ve yetenekleri resmî olarak teyit edilmemiştir; konu büyük ölçüde gazetecilik, araştırmacı iddialar ve parlamento raporları üzerinden bilinir. Burada yalnızca "böyle bir tartışmanın kamusal olarak yürütüldüğü" olgusu aktarılmıştır; sistemin somut iddiaları için kaynaktan teyit edilmeli ve doğrulanmamış iddialar olarak ele alınmalıdır.

### Snowden ifşaatları: kavramsal bağlam

2013'te, bir eski ABD istihbarat sözleşmelisi olan Edward Snowden'ın gazetecilere sağladığı belgeler, çok sayıda sinyal istihbaratı programı ve kabiliyeti iddiasını kamuoyuna taşıdı. Bu ifşaatların ayrıntılı operasyonel içeriği bu kitabın konusu değildir ve burada aktarılmayacaktır. Ancak ifşaatların kamuya mal olmuş, kavramsal düzeydeki etkisi, bir SIGINT okuryazarının bilmesi gereken bağlamdır:

- İfşaatlar, devlet düzeyinde sinyal istihbaratının ölçeği ve kapsamı üzerine küresel bir kamu tartışması başlattı; mahremiyet, hukuki gözetim ve yasal sınırlar konusunu ana akım siyasetin merkezine taşıdı.
- Teknoloji şirketleri ve standart kuruluşları açısından, uçtan uca şifrelemenin yaygınlaşması ve haberleşme güvenliğine verilen önemin artması gibi gözlemlenebilir sonuçlar doğurdu. (Bu, modern dönemin "yaygın güçlü şifreleme" eğiliminin hızlanmasıyla örtüşür; Bölüm 4.)
- Hukuki ve düzenleyici tarafta, çeşitli ülkelerde istihbarat gözetimi ve denetim mekanizmaları üzerine reform tartışmalarını tetikledi.

### Sızdırılan ANT katalog fenomeni: araçların var olduğu bilgisi

Snowden bağlamında kamuoyuna yansıyan olgulardan biri, çeşitli RF ve donanım istihbarat araçlarını listelediği iddia edilen bir kataloğun (yaygın olarak "ANT kataloğu" adıyla anılan) kamusallaşmasıdır. Bu kitabın bu olguya yaklaşımı kesin ve sınırlıdır: önemli olan, bu tür araçların (örneğin yan-kanal/RF yayılım yoluyla bilgi sızdırma, donanıma yerleştirilen pasif/aktif implantlar gibi kategorilerin) var olabildiği bilgisinin kamusallaşmasıdır; bu araçların operasyonel kullanımı, kurulumu veya yeniden üretimi değil.

Bir savunmacı için bu fenomenin tek doğru kullanımı, savunma farkındalığıdır: "Demek ki donanım/RF düzeyinde bu tür yan-kanal ve implant kategorileri teorik ve pratik olarak mümkün; öyleyse kendi tehdit modelimde fiziksel güvenliği, tedarik zinciri bütünlüğünü ve TEMPEST/yan-kanal savunmasını (Bölüm 13) ciddiye almalıyım." Bu, ANT katalog fenomeninin bu kitaptaki tek meşru dersidir. Hiçbir aracın tasarımı, çalışma reçetesi veya kullanımı burada anlatılmaz; yalnızca "kategori olarak mümkün" bilgisi, savunma tarafını besler.

> Net sınır: Bu bölüm, ifşaatların operasyonel detayını, araç tasarımını veya kullanım yöntemini aktarmaz ve aktaramaz. Yalnızca kamuya mal olmuş kavramsal etkiyi (kamu tartışması, şifreleme eğilimi, yasal reform, savunma farkındalığı) işler. İfşaatlara dair belirli teknik iddialar gazetecilik kaynaklıdır ve bağımsız teyit gerektirir.

### Dengeli okuma

ECHELON ve Snowden bağlamı, bu kitapta ne bir suçlama ne bir savunma olarak sunulur. İki olgu da gerçektir: kitlesel sinyal istihbaratı kapasitesinin var olabildiği kamuoyuna mal olmuştur ve bu, meşru bir gözetim-özgürlük tartışması doğurmuştur. Bu tartışmanın iki yakası (ulusal güvenlik ihtiyacı ile bireysel mahremiyet hakkı) Bölüm 13'te dengeli biçimde ele alınır. Bir SIGINT pratisyeninin görevi taraf tutmak değil, meselenin gerçekliğini ve karmaşıklığını tanımaktır.

---

<a id="7"></a>
## 7. Sayı İstasyonları (Numbers Stations): HF'de Gizem ve Tek-Yönlü Ajan Bağı

Kısa dalga (HF) bandını tarayan herkes, er ya da geç tuhaf bir yayınla karşılaşır: monoton bir ses (çoğu zaman sentetik veya kayıt bir kadın/erkek/çocuk sesi) sıralı sayılar veya harfler okur; bazen bir melodi veya elektronik bir motifle başlar; sonra biter. Bunlara sayı istasyonları (numbers stations) denir ve bunlar, açık spektrumun en gizemli ama aynı zamanda en iyi anlaşılmış olgularından biridir.

### Ne işe yarar? Tek-yönlü ses bağı (one-way voice link)

Sayı istasyonlarının en yaygın ve kabul gören açıklaması, tek-yönlü ajan haberleşmesi (one-way voice link, OWVL) olmalarıdır. Mantık şudur: bir istihbarat servisi, sahadaki bir ajanına mesaj göndermek ister; ama ajanın yakalanmasını riske atacak şekilde ona bir verici taşıtmak istemez. Çözüm dahiyane biçimde basittir: ajan yalnızca alır. Sahadaki kişinin tek ihtiyacı, dünyanın her yerinde satılan, üzerinde hiçbir suç unsuru olmayan sıradan bir kısa dalga radyosudur. HF, iyonosfer yansımasıyla kıtalararası menzile ulaştığından (Bölüm 1), merkez bir vericiyle binlerce kilometre öteye yayın yapmak mümkündür.

Bu modelin neden bu kadar dayanıklı olduğunu anlamak önemlidir: tek-yönlü bağ, ajan tarafında hiçbir iz bırakmaz. Sadece dinleyen birinin elinde sıradan bir radyo vardır; ne bir verici, ne bir bağlantı kaydı, ne de bir ağ izi. Bu yüzden sayı istasyonları, dijital ve izlenebilir modern haberleşmenin tersine, izlenemezliğin (untraceability) klasik bir biçimidir.

### Neden hâlâ var? Düşük teknolojinin dayanıklılığı

Bir soru kaçınılmazdır: uydu, internet ve şifreli mesajlaşma çağında neden hâlâ HF'de sayı okuyan istasyonlar var? Açık-kaynak düzeydeki yaygın değerlendirme, tam da basitliğin getirdiği güvenlik avantajıdır:

| Özellik | Sayı istasyonu (HF OWVL) | Modern dijital bağ |
|---|---|---|
| Ajan tarafında iz | Yok (sadece sıradan radyo) | Cihaz, hesap, ağ izi olabilir |
| Konum ele verme | Alıcı yayın yapmaz → izlenemez | Bağlantı meta-verisi sızabilir |
| Altyapı bağımlılığı | İyonosfer + sıradan radyo | İnternet/uydu altyapısı |
| İnkar edilebilirlik | Yüksek (içerik tek-kullanımlık şifreli) | Daha düşük |

Yani sayı istasyonları "geri kalmış" değil, tam tersine belirli bir tehdit modeli için fazlasıyla rasyonel bir çözümdür: minimum iz, maksimum inkar edilebilirlik. Tek-kullanımlık şifre (one-time pad, bir sonraki başlık) ile birleştiğinde, içeriği teorik olarak kırılamaz hale gelir.

### Conet Project: olgunun belgelenmesi

Sayı istasyonları olgusunun kamuoyunda tanınmasında önemli bir kilometre taşı, The Conet Project'tir: 1990'larda yayımlanan, dünyanın dört bir yanından kaydedilmiş sayı istasyonu yayınlarını derleyen bir ses arşividir. Bu derleme, daha önce yalnızca kısa dalga dinleyicilerinin (DXer) bildiği bu olguyu geniş bir kültürel ve araştırmacı kitleye taşıdı ve sayı istasyonlarının akademik/gazetecilik düzeyinde belgelenmesine katkı sundu.

### Nasıl dinlenir? (Yalnızca pasif, yasal RX)

Sayı istasyonları, açık HF bandında yayın yapar ve dinlenmesi (yalnızca alım) dünyanın çoğu yerinde serbesttir; içerik şifreli ve sana yönelik olmadığından çözmen zaten mümkün değildir ve amaç da değildir. Pasif gözlem için pratik notlar:

- Bir HF kapsamlı alıcı (uygun bir SDR ve HF anteni; Bölüm 2-3) veya bir çevrimiçi WebSDR kullanılabilir. WebSDR'ler, dünyanın çeşitli yerlerindeki alıcılara tarayıcıdan erişim sağlar ve HF olgularını gözlemlemenin tamamen yasal, donanımsız bir yoludur.
- Sayı istasyonları çoğunlukla belirli HF dilimlerinde, çoğu zaman düzenli zaman çizelgeleriyle yayın yapar; topluluklar (aşağıda) bu çizelgeleri açık biçimde takip eder ve istasyonlara takma adlar (ör. ünlü bir örnek olarak anılan "Lincolnshire Poacher" gibi, melodisiyle tanınan tarihsel istasyonlar) verir.
- Gözlemlerken yalnızca dinle ve notla; kaydetme ve yeniden yayma konusunda kendi ülkenin mevzuatını teyit et. Amaç olguyu tanımak, sesini ve örüntüsünü gözlemlemektir; içeriği çözmek ne mümkün ne hedeftir.

> Mühendislik ve tarih sezgisi: Sayı istasyonu, modern izlenebilir haberleşmenin tersini öğretir. Bütün bu kitap "her yayın bir iz bırakır" der; sayı istasyonu ise bu dersi ajan tarafından silme zanaatıdır — yayını merkeze, izi sıfıra indirir. Bu, OPSEC'in (Bölüm 6) en saf tarihsel örneğidir: en güvenli verici, hiç yayın yapmayan alıcıdır.

> Not: Belirli sayı istasyonlarının kime ait olduğu, hangi servisin işlettiği ve içeriğin tam doğası çoğunlukla resmî olarak doğrulanmamıştır; bu konudaki atıflar araştırmacı/topluluk kaynaklıdır ve kesin atıf iddiaları teyit edilmelidir. "Tek-yönlü ajan bağı" açıklaması yaygın ve makul kabul edilen çerçevedir ama her istasyon için kanıtlanmış değildir.

---

<a id="8"></a>
## 8. One-Time Pad ve Casus Radyo Tarihçesi: Burst Transmission, The Thing

Sayı istasyonları neden kırılamaz? Cevap, bir kriptografi klasiğindedir: tek-kullanımlık şifre defteri (one-time pad). Ve casus haberleşmesinin RF tarihi, bu kavramın etrafında, izi en aza indirmeye çalışan bir dizi zekice çözüm üretti. Bu başlık bu tarihi kavramsal olarak ele alır; hiçbir kurma talimatı içermez.

### One-time pad: teorik mükemmellik

Tek-kullanımlık şifre, basit ama matematiksel olarak kanıtlanmış biçimde kırılamaz bir şifreleme yöntemidir. Mantığı şudur: mesaj kadar uzun, tamamen rastgele bir anahtar (pad) üretilir; mesaj bu anahtarla birleştirilir (klasik haliyle harf/sayı toplama, dijital haliyle XOR); anahtar yalnızca bir kez kullanılır ve sonra imha edilir. Anahtar gerçekten rastgele, mesaj kadar uzun, gizli ve tek-kullanımlık ise, şifreli metin hiçbir bilgi sızdırmaz; her olası açık metin eşit olasılıkla mümkündür. Bu, Claude Shannon'un biçimlendirdiği "mükemmel gizlilik" (perfect secrecy) kavramının pratik karşılığıdır.

One-time pad'in zayıflığı matematikte değil, pratiktedir: anahtar dağıtımı ve yönetimi. Anahtar mesaj kadar uzun olmalı, güvenli biçimde her iki tarafa ulaşmalı, gerçekten rastgele olmalı ve asla tekrar kullanılmamalıdır. Tarihte one-time pad'in çöktüğü durumlar, neredeyse her zaman bu pratik kurallardan birinin ihlalidir (örneğin anahtarın yeniden kullanılması). Sayı istasyonları + one-time pad birleşimi tam olarak şu yüzden güçlüdür: ajan istasyondan sayıları alır, kâğıt pad'iyle çözer, pad'i imha eder; ele geçse bile geçmiş mesajlar kurtarılamaz.

```
 One-time pad mantığı (kavramsal):

 Açık metin:   H  E  L  L  O      (sayıya çevrili)
 Pad (rastgele,
  tek kullanım): 23 11 04 19 07
 Şifreli:      (açık + pad) mod N  ──► yayınlanır (ör. sayı istasyonu)

 Alıcı: şifreli − pad = açık metin,  sonra pad'i İMHA et.
 Anahtar tek-kullanımlık ve rastgele ise → matematiksel olarak kırılamaz.
 Tek zaafı: pad'in üretimi, dağıtımı ve "asla tekrar kullanma" kuralı.
```

### Casus radyo tarihçesi ve burst transmission

Soğuk Savaş'ın casus telsizleri (yaygın olarak "spy radio" denen taşınabilir verici/alıcı setleri), izlenme riskini azaltmak için tasarım yapıldı. En önemli kavramlardan biri burst transmission (patlama/sıkıştırılmış iletim) idi.

Mantık, yön bulmaya (Bölüm 3) karşı bir savunmadır. Bir ajan elle morse gönderirse, mesaj dakikalarca sürer; bu süre, karşı-istihbaratın yön bulma ekiplerine vericiyi konumlandırmak için bolca zaman tanır. Burst transmission, mesajı önceden kaydedip çok yüksek hızda, saniyenin altında bir sürede tek bir "patlama" halinde gönderir. Yayın o kadar kısadır ki, yön bulma ekipleri ölçümü tamamlayamadan biter. Bu, "havada geçirilen süreyi en aza indir" OPSEC ilkesinin (Bölüm 6) klasik donanım uygulamasıdır: ne kadar kısa yayın, o kadar zor tespit.

| Casus radyo kavramı | Amaç | Karşı koyduğu tehdit |
|---|---|---|
| Tek-yönlü alım (sayı istasyonu) | Ajanda verici/iz olmaması | Verici tespiti, ele geçirme |
| Burst transmission | Yayın süresini saniye altına indirme | Yön bulma (DF) ile konumlandırma |
| One-time pad | İçeriği matematiksel olarak korumak | Kriptanaliz |
| Gizli/kamufle cihaz | Fiziksel ele geçmede inkar | Fiziksel arama |

Bu kavramların hepsinin ortak teması, yine bu kitabın merkezi dersidir: haberleşmek bir zorunluluk ama aynı zamanda bir risktir; casus radyo zanaatı, bu riski (iz, süre, içerik, fiziksel kanıt) her boyutta en aza indirme sanatıdır.

### The Thing (Great Seal bug): pasif rezonatörün dehası

RF casusluk tarihinin en zarif ve en çok anlatılan örneği, yaygın olarak "The Thing" veya "Great Seal bug" diye anılan pasif dinleme cihazıdır. Tarihsel anlatıya göre, bir ABD büyükelçilik binasındaki ahşap bir devlet arması (Great Seal) içine yerleştirilmiş bu cihaz, uzun süre tespit edilemedi; çünkü alışılmış böcek arama yöntemlerinin aradığı şeyi içermiyordu.

Dehası, çalışma prensibindeydi. The Thing, içinde batarya, aktif elektronik veya kendi vericisi olmayan pasif bir rezonatördü. Kendi başına hiçbir sinyal yaymıyordu; dolayısıyla yayın arayan bir dedektör onu bulamıyordu. Çalışması için dışarıdan, güçlü bir radyo dalgasıyla "aydınlatılması" gerekiyordu. Bu harici dalga geldiğinde, cihazın içindeki boşluğun (kavite) rezonansı, odadaki ses titreşimleriyle hafifçe değişiyor ve yansıyan dalgayı ses bilgisiyle modüle ediyordu. Yani cihaz, kendi enerjisini değil, dışarıdan verilen enerjiyi ses bilgisiyle "renklendirerek" geri yansıtıyordu.

```
 The Thing — pasif rezonatör prensibi (kavramsal):

   Harici verici  ──── güçlü RF dalgası ────►  ┌─────────────┐
   (aydınlatıcı)                                │  Pasif      │
                                                │  rezonatör  │◄── oda sesi
   Harici alıcı  ◄─── ses ile modüle ───────── │  (bataryasız,│    (zar/kavite
   (dinleyici)        edilmiş yansıma           │   vericisiz) │     titreşir)
                                                └─────────────┘
   Cihaz kendi enerji yaymaz → yayın arayan dedektör KÖR kalır.
   Ancak harici aydınlatma + ses modülasyonu = dinleme.
```

The Thing'in tarihsel önemi iki yönlüdür. Birincisi, casusluk dehası açısından: tespit edilebilirliği en aza indirmenin en uç biçimi, hiç yayın yapmamaktır; cihaz yalnızca aydınlatıldığında "canlanır". İkincisi, bu fikir modern teknolojinin atasıdır: pasif, harici enerjiyle çalışan ve yansımayı modüle eden bu prensip, bugün her gün kullandığımız pasif RFID etiketlerinin (Bölüm 5) temel çalışma mantığının erken ve ünlü bir örneğidir. Bir RFID etiketi de kendi bataryası olmadan, okuyucunun verdiği enerjiyle "uyanır" ve cevabını yansıtır.

> Savunma dersi (Bölüm 13 ile bağ): The Thing, neden "yayın aramanın" tek başına yeterli bir TEMPEST/böcek savunması olmadığını tarihsel olarak kanıtlar. Pasif, harici-enerjili tehditler, klasik yayın dedektörünü atlatabilir; bu yüzden modern fiziksel güvenlik, yalnızca aktif yayın taraması değil, harici aydınlatma altında non-lineer yanıt arama (non-linear junction detection) gibi yöntemleri de içerir. Tarih, savunmanın kapsamını genişletmemiz gerektiğini öğretir.

> Not: The Thing'in tarihsel ayrıntıları (yerleştirilme biçimi, ne kadar süre tespit edilemediği, teknik tasarım detayı) yaygın biçimde anlatılır ve genel hatlarıyla kabul görür; bazı ayrıntılar kaynaklar arasında farklılaşabilir ve kesin teknik iddialar için tarihsel kaynaklardan teyit edilmeli. Burada amaç kurma reçetesi değil, pasif rezonatör prensibinin kavranmasıdır.

---

<a id="9"></a>
## 9. Elektronik Harp Doktrini ve Sistemleri: EA/EP/ES, RWR, Chaff/Flare, SEAD

Sinyal istihbaratı, elektronik harbin (EW, Electronic Warfare) bir komşusu ve çoğu zaman ön koşuludur. Bölüm 13, EW taksonomisini savunma odağıyla (EA'yı yalnızca tanımak, EP ve ES'i inşa etmek) ele aldı. Burada konuyu tarihsel ve kavramsal/kültürel çerçevede tamamlıyoruz: EW sistemleri nereden çıktı, hangi klasik araçlar bu doktrini şekillendirdi ve modern eğilim nereye gidiyor. Bütün anlatım açık-kaynak ve kavramsaldır; hiçbir icra reçetesi içermez.

### EA / EP / ES: doktrinin üç ayağı (özet ve tarihsel bağ)

EW'nin üç ana fonksiyonu (modern terimlerle) şöyledir; ayrıntılı işlenişi Bölüm 13'tedir:

| Fonksiyon | Açılım | Özü | Eski adı (tarihsel) |
|---|---|---|---|
| ES | Electronic Support | Spektrumu izleyip tehdidi tespit/tanı/yer tespiti yapmak | ESM (Electronic Support Measures) |
| EA | Electronic Attack | Spektrumu rakip aleyhine kullanmak: karıştırma, aldatma, sahteleme | ECM (Electronic Countermeasures) |
| EP | Electronic Protection | Kendi sistemini EA ve parazitten korumak | ECCM (Electronic Counter-Countermeasures) |

Tarihsel olarak bu üçlü, II. Dünya Savaşı'nda radar etrafında doğdu. Radar bir "elektronik göz" olarak savaşa girince (Bölüm 10), neredeyse hemen onu kör etme (karıştırma, EA), karıştırmaya direnme (EP) ve düşman radarını dinleyip tanıma (ES) ihtiyaçları doğdu. SIGINT'in ES ile ilişkisi yakındır: ES, bir EW kavramı olarak, SIGINT'in (özellikle ELINT'in) savaş alanı gerçek-zamanlı kardeşidir. Bir RWR (aşağıda) aslında küçük, gerçek-zamanlı bir ELINT alıcısıdır.

### Chaff ve flare: en eski karşı önlemler

EW'nin en eski ve en sezgisel araçları, II. Dünya Savaşı'na uzanır:

- Chaff (saman/radar yansıtıcı): Radar dalga boyuna uyumlu boyutta kesilmiş, havaya saçılan ince metalik şeritler. Havada bir bulut oluşturur ve radar ekranında devasa sahte yankılar yaratarak gerçek hedefi gizler veya radarı şaşırtır. Tarihsel olarak bu teknik, II. Dünya Savaşı hava harekâtlarında düşman radarlarını boğmak için kullanıldı (İngilizlerin "Window" kod adıyla andığı uygulama bunun erken örneğidir). Chaff, bir pasif EA biçimidir: enerji yaymaz, yalnızca yansıtır.

- Flare (sıcak parazit/yem): Kızılötesi güdümlü füzelere karşı, uçaktan atılan yüksek sıcaklıklı parlayıcılar. Füzenin kilitlendiği sıcak hedefi (motor egzozu) taklit ederek füzeyi kendine çeker. Bu, RF değil kızılötesi alanda bir karşı önlemdir ama EW kültürünün ayrılmaz parçasıdır.

```
 Chaff prensibi (kavramsal, radar ekranı):

 Gerçek hedef tek yankı:        Chaff sonrası ekran:
   ·                              ▓▓▓▓▓▓▓
   (uçak)                         ▓▓·▓▓▓▓   ← gerçek hedef sahte
                                  ▓▓▓▓▓▓▓     yankı bulutunda kaybolur
                                  ▓▓▓▓▓▓
   Chaff = dalga boyuna kesilmiş metalik şeritler, havada bulut.
   Radar için "her yer hedef" → gerçeği ayırt edemez.
```

### RWR: radar uyarı alıcısı

Radar Warning Receiver (RWR, radar uyarı alıcısı), bir platformu (uçak, gemi) aydınlatan düşman radarlarını tespit edip pilota/operatöre uyaran sistemdir. İşleyişi, Bölüm 7'deki ELINT mantığının gerçek-zamanlı, gömülü halidir: gelen darbeleri yakalar, parametrelerini (RF, PRI, PW, tarama) ölçer, bir tehdit kütüphanesiyle eşleştirir ve "seni şu tip radar aydınlatıyor, şu yönden, ve muhtemelen izleme/kilit modunda" bilgisini üretir. RWR'nin kalbi, tam olarak Bölüm 7'de anlatılan PDW üretimi, deinterleaving ve kütüphane eşleştirmesidir; sadece milisaniyeler içinde ve hayatta kalma amacıyla.

RWR bir EP/ES aracıdır: kendisi saldırmaz, ama tehdidi tanıyarak pilotun doğru kaçınma/karşı önlem kararını (chaff atma, manevra, jammer devreye alma) almasını sağlar. Sivil dünyadaki karşılığı, bir kritik altyapı tesisinin "beni bir radar/güçlü yayıcı mı aydınlatıyor?" diye spektrumu izleyen ES sensörüdür (Bölüm 13).

### Kendini koruma jammer'ları

Kendini koruma karıştırıcısı (self-protection jammer), bir platformun, kendisine kilitlenen radar veya güdüm sistemini karıştırarak (gürültü veya aldatıcı sahte hedefler üreterek) füzenin/atışın isabetini bozmaya çalışan EP/EA sistemidir. Bunlar tipik olarak pod (harici kapsül) veya gömülü sistem biçimindedir. Bu kitap işleyiş reçetesi vermez; kavramsal olarak önemli olan, bunun EA'nın savunma amaçlı (kendini koruma) kullanımı olması ve modern hava platformlarında RWR + chaff/flare + jammer'ın bir "kendini koruma paketi" olarak birlikte çalışmasıdır.

### Anti-radyasyon füze ve SEAD kavramı

EW kültürünün bir diğer klasik kavramı, düşman hava savunmasının bastırılmasıdır (SEAD, Suppression of Enemy Air Defenses). Bu doktrinin bir aracı, anti-radyasyon füzesidir (ARM, Anti-Radiation Missile): düşman radarının yaydığı sinyalin kendisine yönelerek (sinyali "ev" olarak kullanarak) radarı imha eden füze. Buradaki kavramsal güzellik, SIGINT/ELINT ile silahın doğrudan birleşmesidir: radar yayın yaptığı sürece, kendi yayınıyla hedef olur. Bu, "yayın yapmak risktir" ilkesinin (Bölüm 6) en uç askeri sonucudur ve düşman radar operatörlerini, hayatta kalmak için yayınlarını kısa tutmaya veya kapatmaya (emission control) zorlar.

```
 SEAD/ARM kavramsal döngüsü:

   Düşman radarı yayın yapar ──► ELINT/RWR yayını tespit eder
            ▲                              │
            │                              ▼
   "yayını kapat/kısalt"          ARM, radar yayınına yönelir
   (emission control = OPSEC)             │
            │                              ▼
            └──── radar susmaya/saklanmaya ── radar tehdit altında
                  zorlanır (bu da bir EW etkisidir)
```

> Not: SEAD ve ARM kavramları açık askeri doktrin literatüründe geniş yer tutar; belirli sistemlerin yetenekleri ve performans iddiaları sınıflandırılmış olabilir. Burada yalnızca doktrinin mantığı (yayının kendisinin hedef olması ve emission control baskısı) kavramsal olarak verilmiştir.

### Modern eğilim: bilişsel elektronik harp

EW'nin modern yönelimi, bilişsel elektronik harptir (cognitive EW). Klasik EW sistemleri, önceden programlanmış bir tehdit kütüphanesine (Bölüm 7, Bölüm 10) dayanır: gördükleri sinyali bilinen kayıtlarla eşleştirir ve önceden tanımlı bir tepki uygular. Sorun şudur: düşman sistemleri artık yazılım tanımlı (SDR tabanlı) ve çevik; parametrelerini (frekans, PRI, dalga biçimi) hızla değiştirip kütüphanedeki kaydı geçersiz kılabilirler.

Bilişsel EW, bu meydan okumaya makine öğrenmesi ve uyarlanabilirlikle yanıt verir: sistem, daha önce görmediği bir yayıcıyı gerçek zamanlı analiz eder, davranışını öğrenir ve en uygun tepkiyi (hangi karşı önlem, hangi parametre) kütüphaneye bel bağlamadan kendi üretir. Bu, Bölüm 7'deki "öğrenme tabanlı ayıklama/sınıflandırma" eğiliminin EW'deki karşılığıdır. Kavramsal olarak EW, sabit kütüphane eşleştirmeden, gerçek-zamanlı öğrenen bir döngüye doğru kayıyor.

| EW kuşağı | Yaklaşım | Sınırı |
|---|---|---|
| Klasik | Sabit tehdit kütüphanesi + önceden tanımlı tepki | Bilinmeyen/çevik tehdide kör |
| Uyarlanabilir | Parametrik esneklik, sınırlı seçenek havuzu | Hâlâ önceden tasarlanmış senaryolara bağlı |
| Bilişsel (modern) | Gerçek-zamanlı öğrenme, kütüphanesiz tepki üretimi | Veri/güven, açıklanabilirlik, doğrulama zorlukları |

> Mühendislik ve doktrin sezgisi: EW ile SIGINT aynı fiziği paylaşır ama amaçları farklıdır. SIGINT spektrumu anlamak için dinler (pasif, bilgi üretir); EW spektrumda etki yaratmak (EA) veya etkiye direnmek (EP) için davranır. ES, ikisinin köprüsüdür: hem SIGINT'in savaş alanı kardeşi hem EW'nin "gözü". Bu kitap, savunma duruşu gereği ağırlığını ES (tespit/izleme) ve EP (dayanıklılık) tarafına koyar; EA'yı yalnızca tanımak için anlatır (Bölüm 13).

---

<a id="10"></a>
## 10. ELINT ve Radar İstihbaratı Tarihçesi: Radar Parmak İzi ve Savaş Düzeni

ELINT (Electronic Intelligence, Bölüm 7), haberleşme olmayan elektronik yayınların — başta radar — istihbaratıdır. Tarihsel kökü, radarın savaşa girişiyle eşzamanlıdır ve bugünkü tehdit kütüphanelerinin, RWR'lerin ve elektronik savaş düzeninin temelini atan bir öyküdür.

### Radarın doğuşu ve "radar savaşı"

Radar (RAdio Detection And Ranging), II. Dünya Savaşı öncesinde ve sırasında, hava savunmasının belkemiği olarak olgunlaştı. Bir radar, hedeflere enerji darbeleri gönderir ve yankıları dinleyerek menzil, yön ve hız çıkarır. Radar savaş alanına girdiği an, onu hedef alan bir "karşı-disiplin" de doğdu: düşman radarını dinlemek, parametrelerini çıkarmak, konumunu bulmak ve gerektiğinde karıştırmak. Bu sürekli hamle-karşı hamle (radar geliştir → karıştır → karıştırmaya dirençli radar geliştir → yeni karıştırma...) "radar savaşı" veya "dalga boyları savaşı" diye anılan tarihsel dinamiği yarattı ve EW doktrinini (Bölüm 9) şekillendirdi.

### Radar parmak izi ve ELINT kütüphanesinin doğuşu

ELINT'in temel keşfi şuydu: her radar, parametrelerinin birleşimiyle (RF, PRI ve türü, PW, tarama tipi/periyodu, polarizasyon, darbe içi modülasyon — Bölüm 5-6) bir parmak izi taşır. Aynı modelin radarları benzer parametre demetine sahipken, farklı modeller ayırt edilebilir biçimde farklıdır. Bu, "radar parmak izi" kavramını ve onun kurumsal ürününü, ELINT parametre kütüphanesini (Bölüm 7, Bölüm 10) doğurdu: bilinen her yayıcı tipinin ölçülmüş parametre aralıklarının kataloğu.

Bir adım derinde, aynı modelin tek tek cihazlarını bile ayırt etme fikri (kasıtsız modülasyon ve üretim kusurlarından) ortaya çıktı; bu, Bölüm 7'de işlenen SEI (Specific Emitter Identification) ve ELINT-MASINT kesişiminin tarihsel köküdür. Yani "hangi model radar" sorusundan "hangi tekil cihaz" sorusuna geçiş, ELINT'in olgunlaşma yörüngesidir.

### Elektronik savaş düzeni (EOB)

ELINT'in operasyonel ürünü, elektronik savaş düzenidir (EOB, Electronic Order of Battle): bir coğrafyadaki düşman elektronik yayıcılarının (radarlar, hava savunma sistemleri, haberleşme düğümleri) türlerinin, konumlarının ve faaliyet durumlarının haritası. EOB, OpELINT'in (Bölüm 7) çıktısıdır ve şu soruya cevap verir: "Karşımda hangi sistemler var, neredeler, hangileri aktif ve neyi koruyorlar?"

```
 Elektronik savaş düzeni (EOB) — kavramsal:

   ELINT toplama (TechELINT)        ELINT izleme (OpELINT)
   ──────────────────────          ──────────────────────
   "Bu yayıcı hangi tip?"          "Bu yayıcı nerede,
   (parametre → kütüphane)          ne zaman, aktif mi?"
           │                                │
           └────────────┬───────────────────┘
                        ▼
              ELEKTRONİK SAVAŞ DÜZENİ (EOB)
              tip + konum + durum haritası
              → tehdit resmi, kaçınma/SEAD planı,
                durumsal farkındalık
```

EOB'nin sivil/savunma karşılığı, bir kritik altyapı korumacısının çevresindeki RF emitör resmini (hangi vericiler var, hangileri normal, hangisi anomali) çıkarmasıdır (Bölüm 13). Mantık birebir aynıdır: önce çevreni tanı (baseline), sonra sapmayı yakala.

### Tarihsel ders: yayın yapan, tanınır

ELINT tarihinin kalıcı dersi, bu kitabın merkezî temasını bir kez daha doğrular: bir radar işini yapmak için yayın yapmak zorundadır (hedefleri aydınlatmalıdır), ama bu yayın aynı zamanda onu tanınabilir, konumlandırılabilir ve hedeflenebilir kılar (SEAD/ARM, Bölüm 9). Bu yüzden modern radarlar LPI tekniklerine (Bölüm 7, düşük tespit olasılığı), emission control'e ve çevikliğe yönelmiştir. "Görmek için yayın yapmak zorundasın, ama yayın yaptığın an görünürsün" gerilimi, ELINT ile radar tasarımı arasındaki ebedi yarışın motorudur.

> Not: Belirli radar sistemleri, ülkeler ve EOB içerikleri sınıflandırılmıştır; bu başlık yalnızca kamuya açık ELINT kavramlarını (radar parmak izi, kütüphane, EOB, SEAD baskısı) tarihsel çerçevede verir. Somut sistem iddiaları içermez.

---

<a id="11"></a>
## 11. Sivil ve Ticari SIGINT: Spektrum İzleme, Kaçak Verici Avı, Geolocation

SIGINT yalnızca askeri ve istihbarat dünyasının işi değildir. Aynı teknik temel (yakala, ölç, sınıflandır, yer tespit et), tamamen sivil, meşru ve günlük bir alanda da kullanılır: spektrum yönetimi ve düzenleme. Bu, bir SIGINT okuryazarının görmesi gereken "öteki yüz"dür; çünkü buradaki faaliyetler açık, yasal ve çoğu zaman kamuya hizmet edicidir.

### Spektrum izleme kurumları (regülatör monitoring)

Her ülkenin bir telekomünikasyon düzenleyicisi, radyo spektrumunu yönetir: kim hangi frekansta, hangi güçle, hangi amaçla yayın yapabilir (Bölüm 8'deki frekans tahsisi). Bu düzeni korumak için düzenleyiciler, spektrum izleme (spectrum monitoring) faaliyeti yürütür: sabit ve mobil istasyonlarla spektrumu dinler, kim yayın yapıyor diye ölçer ve tahsis planına uymayanları tespit eder. Uluslararası düzeyde ITU (Uluslararası Telekomünikasyon Birliği), spektrum izleme için teknik el kitapları ve standartlar yayımlar; ulusal düzenleyiciler bu çerçeveye uyumlu çalışır.

Bu faaliyetin işlevleri:

| Sivil spektrum izleme işlevi | Amaç | SIGINT karşılığı |
|---|---|---|
| Tahsis uyum denetimi | Yetkili yayıncıların kurallara uyması | Yayıcı tanıma + parametre ölçümü |
| Parazit/girişim avı | Zararlı girişim kaynağını bulmak | Yön bulma + yer tespiti (Bölüm 3, 9) |
| Kaçak/lisanssız verici tespiti | İzinsiz yayıncıyı yakalamak | Anomali tespiti + geolocation |
| Spektrum doluluk ölçümü | Bandların kullanımını planlamak | Spektrum tarama + istatistik |

### Kaçak verici avı (interference hunting)

Sivil spektrum izlemenin en "saha" yanı, girişim/parazit avıdır (interference hunting): bir bantta beklenmedik, zararlı bir sinyal ortaya çıktığında (örneğin havacılık veya acil haberleşme bandını bozan bir kaynak) onun konumunu bulup susturmak. Teknik tamamen bu kitabın yön bulma (Bölüm 3) ve yer tespiti (Bölüm 9) bölümlerinin pratiğidir: yönlü anten, sinyal gücü gradyanı, üçgenleme ve giderek mobil/araç üstü ölçüm ile kaynağa doğru "yaklaşma". Bu, tamamen yasal ve kamu yararına bir SIGINT pratiğidir; çoğu zaman arızalı bir cihazı, yanlış kurulmuş bir tekrarlayıcıyı veya lisanssız bir yayıncıyı bulup düzeltmekle sonuçlanır.

Radyo amatörleri (Bölüm 14) bu beceriyi sportif bir biçimde de yaşatır: "fox hunting" veya ARDF (Amateur Radio Direction Finding) denen, gizlenmiş düşük güçlü bir vericiyi yön bulma teknikleriyle bulma yarışları, kaçak verici avının eğitsel ve eğlenceli kardeşidir. Bu, yön bulmayı öğrenmenin tamamen yasal ve topluluk destekli yoludur.

### Ticari geolocation servisleri

Sivil tarafın bir başka boyutu, ticari sinyal coğrafi konumlama (geolocation) servisleridir. Açık ve yasal örnekler:

- Uçak ve gemi takip servisleri: ADS-B (uçak) ve AIS (gemi) yayınları açık olduğundan (Bölüm 5), dünya çapında gönüllü ve ticari alıcı ağları bu sinyalleri toplayıp gerçek-zamanlı harita servislerine dönüştürür. Bu, dağıtık, kitle-kaynaklı bir SIGINT toplama ağının tamamen yasal ve kamuya açık örneğidir; herkes bir alıcı kurup bu ağlara katkı verebilir.
- Kablosuz konumlama veritabanları: WiFi erişim noktalarının ve baz istasyonlarının konumlarını haritalayan ticari/topluluk veritabanları, cihazların GPS'siz konum kestirimi için kullanılır. Bunlar, "RF imzasından konum" fikrinin (Bölüm 9) sivil, ölçekli uygulamasıdır.

```
 Sivil SIGINT ekosistemi (yasal, açık):

   Düzenleyici (regülatör)        Topluluk / ticari
   ──────────────────────         ──────────────────
   spektrum izleme                ADS-B / AIS toplama ağları
   girişim avı (DF)               WiFi/hücre konum veritabanı
   kaçak verici tespiti           fox hunting / ARDF (amatör)
        │                              │
        └────────── ortak teknik ──────┘
        yakala · ölç · sınıflandır · yer tespit et
        (askeri SIGINT ile AYNI fizik, FARKLI amaç ve yasa)
```

> Mühendislik ve toplum sezgisi: Sivil spektrum izleme, SIGINT tekniğinin "ışıkta" yapılan halidir. Aynı yön bulma, aynı sınıflandırma, aynı yer tespiti; ama amaç gözetim değil, ortak bir kaynağın (spektrum) düzenini ve güvenliğini korumak. Bir SIGINT öğrencisi için bu alan, becerilerini tamamen yasal ve kamuya yararlı biçimde kullanmanın somut yoludur — ve çoğu ülkede gönüllü olarak katkı vermek mümkündür.

---

<a id="12"></a>
## 12. Türkiye ve Bölge Bağlamı: Spektrum Yönetimi ve Amatör Radyo Kültürü

Bu başlık, genel ve açık bilgi düzeyinde, Türkiye ve benzer bölge ülkeleri bağlamında spektrum yönetimi ve amatör radyo kültürünü ele alır. Amaç, okuyucuyu kendi yerel çerçevesine yönlendirmektir; ayrıntılı ve güncel kurallar için her zaman resmî kaynaklara başvurulmalıdır.

### Spektrum yönetimi: yerel düzenleyici çerçeve

Çoğu ülkede olduğu gibi, spektrum kamusal ve sınırlı bir kaynaktır ve bir devlet kurumu tarafından yönetilir. Türkiye'de telekomünikasyon ve telsiz alanının düzenleyicisi Bilgi Teknolojileri ve İletişim Kurumu'dur (BTK). BTK, frekans tahsisi (Bölüm 8), telsiz cihazlarının yetkilendirilmesi, amatör telsiz sınav ve belgelendirmesi ve spektrum izleme gibi işlevleri yürütür. Genel ilkeler bu kitabın geri kalanıyla uyumludur: dinleme/alım (RX) açık yayınlar için büyük ölçüde serbestken, yayın (TX) yetki ve genellikle lisans gerektirir; karıştırma her durumda yasaktır ve ağır yaptırımı vardır.

Yasal çerçeveye ilişkin genel hatlar (Bölüm 0'da da özetlendiği gibi): haberleşmenin gizliliği ceza hukukunda korunur (Türkiye'de Türk Ceza Kanunu'nun ilgili maddeleri), telsiz/spektrum kullanımını düzenleyici mevzuat ve uluslararası düzeyde ITU çerçevesi belirler. Bu metin hukuki danışmanlık değildir; somut bir faaliyetin (özellikle herhangi bir yayın, kayıt veya çözme) yasallığı için güncel resmî mevzuattan ve gerekirse hukukçudan teyit alınmalıdır.

> Not: Belirli kurum adları, kanun maddeleri ve düzenlemeler zamanla değişebilir; burada verilen çerçeve genel ve tanıtıcıdır. Güncel ve bağlayıcı bilgi için doğrudan resmî düzenleyici kaynaklara başvurulmalı ve "teyit edilmeli" notu ciddiye alınmalıdır.

### Amatör radyo kültürü ve lisans

Amatör radyo (ham radio), spektrumun belirli bantlarının, lisanslı meraklıların teknik deney, kendini geliştirme ve acil haberleşme amacıyla kullanımına ayrıldığı, dünya çapında düzenlenmiş bir hobidir. Türkiye'de de amatör telsiz lisansı, ilgili düzenleyicinin (BTK) sınavıyla alınır; lisans sahibine bir çağrı işareti (call sign) verilir ve belirli amatör bantlarda, belirli kurallar dahilinde yayın hakkı tanır. Amatör radyo, bu kitabın tamamında vurgulanan "yayın bir sorumluluktur" ilkesinin meşru ve düzenlenmiş kapısıdır: lisansla, çağrı işaretiyle ve kurallara uyarak yayın yapma yetkisi.

Türkiye'de ve bölgede amatör radyo, dernekler ve topluluklar etrafında örgütlenir; bu topluluklar sınav hazırlığı, teknik eğitim, contest (yarışma) katılımı ve özellikle afet/acil durum haberleşmesi konularında aktiftir. Deprem gibi afetlere açık bir coğrafyada, amatör radyonun acil haberleşme yedekliliği sağlama rolü (altyapı çöktüğünde dahi çalışabilen, merkezi olmayan haberleşme) özellikle değerlidir; bu, Bölüm 14'teki ARES/acil haberleşme kültürünün yerel karşılığıdır.

> Yönlendirme: Kendi bölgendeki amatör radyo topluluğunu, yerel dernekleri, sınav takvimini ve lisans gerekliliklerini araştırmak, bu kitabın alıştırmalarından biridir (Bölüm 15). Yerel bir kulübe katılmak, hem yasal çerçeveyi doğru öğrenmenin hem de pratiği güvenli/mentörlü biçimde edinmenin en iyi yoludur.

---

<a id="13"></a>
## 13. Etik ve Mahremiyet: Kitlesel Gözetim, Hukuk ve Pratisyenin Sorumluluğu

Bu bölümün belki en önemli başlığı budur. Sinyal istihbaratı, doğası gereği güçlü ve aynı oranda hassas bir alandır: başkalarının haberleşmesine ve hareketine dokunabilir. Bu güç, hem meşru savunma ve kamu yararı için hem de bireysel özgürlüğü tehdit eden kitlesel gözetim için kullanılabilir. Bir pratisyenin teknik yetkinliği kadar etik çerçevesi de olmalıdır.

### Kitlesel gözetim tartışması: iki yaka

ECHELON ve Snowden bağlamının (Bölüm 6) gösterdiği gibi, devlet düzeyinde kitlesel sinyal istihbaratı kapasitesi, modern dünyanın bir gerçeğidir. Bu, derin ve çözülmemiş bir tartışma doğurur. Bu kitap taraf tutmaz; ama dengeli bir okuryazarın her iki yakayı da tanıması gerekir:

| Eksen | Güvenlik/toplum argümanı | Mahremiyet/özgürlük argümanı |
|---|---|---|
| Amaç | Terör, ağır suç, devlet güvenliği tehditlerini önlemek | Bireyin özel hayatı temel bir haktır |
| Kapsam | Hedefli, denetimli toplama meşru olabilir | Ayrım gözetmeyen kitlesel toplama orantısızdır |
| Denetim | Yargısal/yasal gözetim kötüye kullanımı sınırlar | Gizli programlar denetlenemezse güç suistimal edilir |
| Soğutucu etki | Güvenlik özgürlüğün ön koşuludur | Gözetildiğini bilmek ifade ve örgütlenmeyi kısıtlar |
| Veri kalıcılığı | Veri ileride suç çözmede değerli olabilir | Toplanan veri ileride amacı dışında kullanılabilir |

Sağlıklı bir tutum, bu eksenlerin hiçbirini yok saymaz. Mesele "gözetim mi özgürlük mü" gibi ikili bir seçim değil, bunlar arasında orantılılık, denetlenebilirlik ve hukuka uygunluk ilkeleriyle kurulan bir dengedir. Demokratik hukuk düzenlerinde bu denge; yargısal izin, parlamenter denetim, bağımsız gözetim kurumları ve şeffaflık mekanizmalarıyla kurulmaya çalışılır. Bu mekanizmaların yeterliliği, sürekli ve meşru bir kamusal tartışmanın konusudur.

### Hukuki çerçeve: gözetim ve özgürlük dengesi

Çoğu hukuk düzeninde haberleşmenin gizliliği anayasal/temel bir haktır ve buna devlet müdahalesi (yasal dinleme dahil) genellikle şu koşullara bağlanır: kanuni dayanak, meşru amaç, orantılılık ve bağımsız (çoğu zaman yargısal) denetim. Yani "yapılabilir olması" ile "yapılmasına izin verilmesi" farklı şeylerdir; teknik kapasite, hukuki yetkiyi otomatik olarak doğurmaz. Bir pratisyen için bu ayrım hayatidir: bir şeyi teknik olarak yapabiliyor olmak, onu yapma hakkın olduğu anlamına gelmez.

```
 Gözetim-özgürlük dengesi (kavramsal terazi):

      GÜVENLİK / KAMU YARARI          MAHREMİYET / ÖZGÜRLÜK
              │                              │
              ▼                              ▼
         ┌─────────────────────────────────────┐
         │      DENGE: hukuk + orantılılık      │
         │      + bağımsız denetim + şeffaflık  │
         └─────────────────────────────────────┘
                          │
              "Yapılabilir" ≠ "Yapmaya izinli"
              Teknik kapasite, hukuki yetkiyi doğurmaz.
```

### Bir SIGINT pratisyeninin etik sorumluluğu

Bu kitabın okuru için somut etik çerçeve, baştan beri savunulan duruşun bir özetidir:

1. Yasallık önce gelir. Her faaliyetten önce kendi ülkenin ve sürümünün mevzuatını teyit et. Şüphedeysen yapma. "Alıcı çoğu yerde serbesttir; verici her yerde sorumluluktur" (Bölüm 0).
2. Pasiflik ve savunma. Bu serinin tüm teknikleri anlama, savunma ve spektrum okuryazarlığı içindir; başkasının haberleşmesini dinlemeyi, çözmeyi, karıştırmayı veya sahtelemeyi değil. Alıştırmalar yalnızca kendi cihazların ve açık/yasal sinyallerle sınırlıdır.
3. Mahremiyete saygı. Teknik olarak erişebileceğin bilgi, erişme hakkın olduğu anlamına gelmez. Başkalarının özel haberleşmesi, içeriği teknik olarak çözülebilir olsa bile, dokunulmazdır.
4. Bilginin sorumlu kullanımı. Bu kitaptaki bilgi (ve genel olarak SIGINT bilgisi) çift kullanımlıdır: hem savunma hem saldırı için kullanılabilir. Onu savunma, eğitim ve kamu yararı yönünde kullanmak bir tercih değil, bir sorumluluktur.
5. Şeffaflık ve hesap verebilirlik. Meşru bir kurumsal bağlamda (örneğin bir güvenlik ekibinde) çalışıyorsan, faaliyetin denetlenebilir, yetkilendirilmiş ve belgelenmiş olmalıdır. Gizli, denetimsiz ve keyfî kullanım, bu alanın en büyük etik tehlikesidir.

> İnsan ve teknik: Bu alanın tarihi, gücün hem hayat kurtardığını (Bletchley'in savaşa katkısı, kurtarma haberleşmesi) hem de kötüye kullanılabildiğini (denetimsiz kitlesel gözetim kaygıları) gösterir. Aradaki fark teknikte değil, etikte ve hukuktadır. Bir SIGINT pratisyenini sorumlu kılan, ne kadar şey yapabildiği değil, neyi yapmamayı seçtiğidir.

---

<a id="14"></a>
## 14. Kültür: Amatör Radyo Topluluğu, DXing, Contesting, ARES, Maker RF

Sinyal dünyasının en açık, en canlı ve bir SIGINT öğrencisi için en erişilebilir yüzü, devlet teşkilatları değil; bir asırlık gönüllü amatör radyo kültürüdür. Bu kültür, hem teknik ustalığı hem etiği hem de kamuya hizmeti taşıyan, dünya çapında bir topluluktur ve bu kitabın değerlerinin yaşayan örneğidir.

### Amatör radyo (ham radio): gönüllü teknik kültür

Amatör radyo, lisanslı meraklıların teknik deney, kendini geliştirme ve kamuya hizmet amacıyla, kâr gütmeden yaptığı bir uğraştır. Yüzyılı aşkın geçmişinde radyo amatörleri, modülasyon tekniklerinden anten tasarımına, zayıf-sinyal haberleşmesinden dijital protokollere kadar pek çok yeniliğe öncülük etti veya katkı verdi. Amatör radyonun kültürel kodları (çağrı işareti kimliği, kurallara titiz uyum, "ragchew" denen sohbetten contest'e uzanan etkinlikler, ve birbirine yardım/mentörlük geleneği) bir SIGINT öğrencisinin de benimsemesi gereken değerlerdir: teknik merak + kurallara saygı + topluluk.

### DXing: uzak ve zayıfın peşinde

DXing, mümkün olan en uzak veya en zayıf istasyonları yakalama/iletişim kurma uğraşıdır ("DX" telgraf kısaltmasından "uzak mesafe" anlamına gelir). HF'de iyonosfer koşullarını (Bölüm 1) ustaca kullanarak kıtalararası temas kurmak, ya da yalnızca dinleyerek (SWL, shortwave listening) dünyanın uzak köşelerinden yayınları yakalamak, bir sabır ve beceri sanatıdır. DXing'in SIGINT'e doğrudan eğitsel değeri vardır: yayılım koşullarını, anten performansını, zayıf-sinyal alımını ve bant planlarını (Bölüm 8) gerçek dünyada öğretir. Sayı istasyonu ve HF utility dinlemenin (Bölüm 7) kültürel evi de bu DXing/SWL geleneğidir.

### Contesting: sportif yoğunluk

Contesting (radyo yarışması), belirli bir süre içinde olabildiğince çok istasyonla, belirli kurallarla temas kurma yarışmasıdır. Yoğun, hızlı ve disiplinli bir etkinliktir; operatörlük becerisini, ekipman hazırlığını ve dayanıklılığı sınar. Contesting, amatör radyonun rekabetçi/sportif yüzüdür ve teknik beceriyi baskı altında uygulamayı öğretir.

### ARES ve acil/kurtarma haberleşmesi

Amatör radyonun belki en değerli toplumsal yüzü, afet ve acil durum haberleşmesidir. ARES (Amateur Radio Emergency Service) gibi örgütlü yapılar ve gönüllü amatörler, deprem, sel, fırtına gibi felaketlerde altyapı (telefon, internet, hücresel) çöktüğünde devreye girer. Amatör radyonun merkezi-olmayan, altyapıdan bağımsız çalışabilen doğası, bu durumlarda hayat kurtaran bir yedeklilik sağlar: bir telsiz, bir batarya/güç kaynağı ve bir anten ile, başka hiçbir şey çalışmazken haberleşme kurulabilir.

```
 Amatör radyonun değer eksenleri:

   TEKNİK          SPORTİF           HİZMET
   ──────          ───────           ──────
   DXing           contesting        ARES / acil haberleşme
   anten/cihaz     operatör          afet yedekliliği
   deneyi          becerisi          (altyapı çökünce çalışır)
       │               │                 │
       └───────────────┴─────────────────┘
            ortak zemin: lisans + çağrı işareti +
            kurallara saygı + topluluk dayanışması
```

Bu, Bölüm 13'teki kritik altyapı dayanıklılığı ve yedeklilik fikrinin (örneğin GNSS/haberleşme çöktüğünde alternatif) gönüllü, insani karşılığıdır. Afete açık bir coğrafyada (Bölüm 12), bu kültürün değeri teorik değil, somuttur.

### Maker ve hackerspace RF kültürü

Modern dönemde amatör radyo kültürüne yeni bir damar eklendi: maker/hackerspace dünyası ve uygun fiyatlı SDR'ların (Bölüm 2) demokratikleştirdiği RF deneyciliği. Açık kaynak yazılım (GNU Radio, SDR uygulamaları), düşük maliyetli donanım (RTL-SDR ve benzeri) ve paylaşımcı topluluk kültürü, bir zamanlar pahalı ve seçkin olan sinyal analizini herkese açtı. Konferanslar, çevrimiçi topluluklar, sigidwiki gibi açık veritabanları (Bölüm 7) ve paylaşılan projeler, bu kültürün ürünleridir. Bu kitabın kendisi de bu demokratikleşmenin bir çocuğudur: bir asır önce devlet tekelinde olan bilgi, bugün sorumlu biçimde herkesin öğrenebileceği bir okuryazarlıktır.

> Kültür ve sorumluluk: Amatör radyo ve maker RF kültürü, bu kitabın savunduğu dengeyi yaşayan bir biçimde gösterir: derin teknik merak, ama her zaman lisans, kural ve etik içinde. Bir SIGINT öğrencisi için bu topluluklar yalnızca bilgi kaynağı değil, doğru değerleri (açıklık, paylaşım, kamuya hizmet, kurallara saygı) edineceği bir okuldur. Tekniği öğrenmek kadar bu kültüre katılmak da bir SIGINT okuryazarını olgunlaştırır.

---

<a id="15"></a>
## 15. Alıştırmalar (Yasal: Dinle, Araştır, Belgele)

> Bu alıştırmalar yalnızca pasif alım (RX), açık/yasal sinyaller, açık-kaynak araştırma ve kâğıt-kalem çalışmasıdır. Hiçbiri yayın (TX), karıştırma veya yetkisiz içerik çözme içermez. Dinleme yapacaksan yalnızca açık yayınları ve serbest bantları dinle; kayıt ve yeniden yayma konusunda kendi ülkenin mevzuatını teyit et. Şüphedeysen yapma. Bu bölüm tarih ve kültür odaklı olduğundan alıştırmaların çoğu araştırma ve gözlemdir.

### A) Bir sayı istasyonu veya HF utility dinle (yalnızca RX, pasif)

Bir WebSDR (tarayıcı tabanlı, donanımsız) veya kendi HF kapsamlı alıcın (Bölüm 2-3) ile kısa dalga bandını gözlemle. Amaç bir sayı istasyonu (Bölüm 7) veya bir HF "utility" yayını (örneğin hava durumu faks, deniz/havacılık HF haberleşmesi gibi açık olgular) duymaya çalışmak ve gözlemini belgelemektir.

| Gözlem alanı | Not |
|---|---|
| Frekans (yaklaşık) | ? |
| Zaman (UTC) | ? |
| Ne duydun (ses/ton/sayı/melodi/dijital) | ? |
| Yayılım nasıldı (kuvvetli/zayıf/sönümlü) | ? |
| Tahmini olgu tipi | ? |

Yalnızca dinle ve notla; içeriği çözmeye çalışma (sayı istasyonlarında zaten mümkün değildir ve amaç değildir). Topluluk takvimlerine bakarak bir sayı istasyonunun yayın zamanını önceden öğrenip o anda dinlemeyi deneyebilirsin. Amaç: Bölüm 7'deki olguyu kendi kulağınla tanımak ve HF yayılımının (Bölüm 1) gerçekliğini deneyimlemek.

### B) Bir tarihsel SIGINT olayını araştır ve özetle

Aşağıdakilerden birini (veya kendi seçtiğin başka bir kamuya-mal-olmuş olayı) açık kaynaklardan araştır ve kısa bir özet yaz:

- Zimmermann Telgrafı: nasıl çözüldü, kaynak nasıl korundu, savaşa etkisi ne oldu? (Bölüm 2)
- Enigma/Ultra: Polonya'nın öncü rolü, Bletchley'in katkısı, kaynak koruma örnekleri. (Bölüm 3)
- The Thing (Great Seal bug): pasif rezonatör prensibi neden tespiti zorlaştırdı, RFID ile bağı nedir? (Bölüm 8)
- ECHELON parlamento tartışması: hangi kaygılar dile getirildi, sonuç ne oldu? (Bölüm 6)

Özetinde şu üç soruyu yanıtla: (1) Olayda hangi SIGINT prensibi işliyor (içerik çözme / trafik analizi / kaynak koruma / pasif dinleme)? (2) Bu olay bugünün hangi kavramına bağlanıyor (bu kitaptaki hangi bölüm)? (3) Olaydan çıkan kalıcı ders nedir? Kaynak değerlendirirken kamuya-mal-olmuş, çok kaynakça doğrulanmış bilgiye dayan; tek kaynaklı veya komplo niteliğindeki iddiaları "doğrulanmamış" diye işaretle. Bu, Bölüm 14'teki (istihbarat kaynakları/değerlendirme) kaynak hijyeninin tarihsel pratiğidir.

### C) Kendi bölgendeki amatör radyo topluluğunu ve lisansı araştır

Tamamen masa başı, yasal bir araştırma: Kendi bölgendeki amatör radyo çerçevesini öğren ve bir sayfalık not çıkar.

1. Yerel düzenleyici kim ve amatör telsiz lisansı nasıl alınır (sınav, sınıflar, çağrı işareti)? (Bölüm 12)
2. Bölgendeki amatör radyo dernekleri/kulüpleri hangileri, hangi etkinlikleri yapıyorlar (contest, fox hunting, eğitim)?
3. Bölgende afet/acil haberleşme (ARES benzeri) yapılanması var mı, nasıl çalışıyor? (Bölüm 14)
4. Yerel olarak hangi bantlar/etkinlikler bir başlangıç için uygun ve yasal?

Amaç: Tekniği yalnızca kitaptan değil, yaşayan ve yasal bir topluluktan öğrenmenin yolunu bulmak. Mümkünse bir yerel kulübün açık etkinliğine katılmak, hem mevzuatı doğru öğrenmenin hem de mentörlü pratik edinmenin en sağlam yoludur.

### D) Kendi etik çerçeveni yaz (düşünce egzersizi)

Kâğıt üzerinde, kendi SIGINT etik çerçeveni (Bölüm 13) maddele. Şu soruları kendine sor ve yanıtla: Hangi faaliyetler benim için kesinlikle "serbest" (kendi cihazım, açık yayın, pasif)? Hangileri "asla" (başkasının haberleşmesi, yayın/karıştırma, yetkisiz çözme)? Bir şeyi "teknik olarak yapabiliyor olmak" ile "yapma hakkım olmak" arasındaki sınırı kendi kelimelerimle nasıl çizerim? Amaç: Bölüm 13'teki etik ilkeleri soyut bilgiden kişisel bir pusulaya dönüştürmek.

### E) EW taksonomisini sivil bir senaryoya uygula (kavramsal eşleme)

Kâğıt üzerinde: Bir kritik altyapı tesisinin RF güvenliğini düşün (Bölüm 13). EA / EP / ES (Bölüm 9) fonksiyonlarını bu sivil senaryoya eşle: Bu tesis için "tehdit" (EA karşılığı) ne olabilir? "Kalkan" (EP karşılığı) ne olmalı? "Göz" (ES karşılığı) nasıl kurulur? Bir tarihsel EW dersini (örneğin "yayın yapan tanınır", ELINT/SEAD, Bölüm 10) bu sivil bağlama nasıl taşırsın? Amaç: Askeri kökenli EW doktrinini sivil savunma okuryazarlığına çevirmek ve Bölüm 9-10-13 arasındaki köprüyü içselleştirmek.

---

<a id="16"></a>
## 16. Hızlı Referans ve Diğer Bölümler

### Kavram kartı

| Kavram | Bir cümlelik öz |
|---|---|
| Room 40 | I. DS İngiliz kod çözme birimi; Zimmermann Telgrafı'nı çözdü; kaynak koruma dersi |
| Zimmermann Telgrafı | Çözülmüş tek mesajın stratejik tarihi değiştirebileceğinin kanıtı |
| Enigma / Ultra | II. DS Alman şifre makinesi ve onu kıran çok-uluslu çaba (Polonya→Bletchley) |
| Bletchley Park | GC&CS; bombe ve Colossus; modern hesaplamanın ve kurumsal SIGINT'in beşiği |
| Trafik analizi (tarihsel) | Enigma çözülmeden bile savaş düzeni çıkardı; "şifre içeriği korur, örüntü sızar" |
| NSA / GCHQ | ABD/İngiltere sinyal istihbaratı ve haberleşme güvenliği kurumları (ikili misyon) |
| Five Eyes / UKUSA | II. DS sonrası beş ülke sinyal istihbaratı paylaşım ittifakı (kök ≈1946) |
| ECHELON | Kitlesel toplama kapasitesinin kamuya mal olduğu tartışma (iddialar teyit gerektirir) |
| Snowden bağlamı | Kamu tartışması + şifreleme eğilimi + yasal reform + savunma farkındalığı (kavramsal) |
| ANT katalog fenomeni | "Bu araç kategorileri var" bilgisinin savunma dersi; reçete değil |
| Sayı istasyonu (OWVL) | Tek-yönlü ajan bağı; ajanda iz sıfır; "en güvenli verici, yayın yapmayan alıcı" |
| One-time pad | Tek-kullanımlık, rastgele anahtar → matematiksel olarak kırılamaz; zaafı anahtar yönetimi |
| Burst transmission | Yayını saniye altına sıkıştırıp yön bulmayı engelleme (OPSEC donanımı) |
| The Thing | Pasif rezonatör casus cihaz; yayın yapmaz, aydınlatılınca canlanır; RFID'nin atası |
| EA / EP / ES | Elektronik saldırı / koruma / destek (eski: ECM/ECCM/ESM); ES ≈ SIGINT'in EW kardeşi |
| Chaff / flare | En eski karşı önlemler; sahte yankı / sahte ısı hedefi |
| RWR | Radar uyarı alıcısı = gerçek-zamanlı gömülü ELINT |
| SEAD / ARM | Düşman radarını bastırma; radar yayınıyla hedef olur → emission control baskısı |
| Bilişsel EW | Sabit kütüphaneden gerçek-zamanlı öğrenen tepkiye geçiş (Bölüm 7 ML'nin EW karşılığı) |
| EOB | Elektronik savaş düzeni: tip+konum+durum haritası (OpELINT ürünü) |
| Sivil SIGINT | Spektrum izleme, girişim avı (DF), kaçak verici tespiti, ADS-B/AIS toplama (yasal) |
| Gözetim-özgürlük dengesi | "Yapılabilir" ≠ "yapmaya izinli"; orantılılık + denetim + hukuk |
| Amatör radyo kültürü | DXing/contesting/ARES/maker; lisans + çağrı işareti + etik = SIGINT okulu |

### Ezber sezgiler

- Tekniği tarihinden ve etiğinden koparma; biri ezbere, diğeri tehlikeye götürür.
- Kaynak koruma en az çözme kadar önemlidir (Zimmermann, Ultra dersi): istihbaratı kullanırken kaynağı yakma.
- "Şifre içeriği korur ama varlık, ritim ve örüntü sızar" — trafik analizinin tarihsel kanıtı II. DS'dir.
- Teşkilat iddialarında kamuya-mal-olmuş düzeyde kal; komplo ve doğrulanmamış iddiadan uzak dur.
- En güvenli verici, hiç yayın yapmayan alıcıdır (sayı istasyonu); en güvenli içerik, tek-kullanımlık pad'dir.
- Pasif/harici-enerjili tehditler (The Thing) klasik yayın taramasını atlatır → savunmanın kapsamını genişlet.
- "Yayın yapan tanınır" (ELINT/SEAD): görmek için yaymak zorundasın, ama yaydığın an görünürsün.
- EW modern eğilimi sabit kütüphaneden bilişsel/öğrenen tepkiye kayar (Bölüm 7 ML'nin kardeşi).
- "Yapılabilir olması", "yapmaya iznin olduğu" anlamına gelmez — etik ve hukuk, tekniğin önündedir.
- Sinyal kültürünün gerçek okulu amatör radyodur: teknik merak + kurallara saygı + kamuya hizmet.

### Ve daima: tarihsel doğruluk, açık kaynak ve perspektif

Bu bölümdeki her olay, kurum ve kavram, kamuya mal olmuş ve açık kaynak düzeyindedir; emin olunamayan tarih, sayı ve iddialar bilinçli olarak "teyit edilmeli" notuyla işaretlenmiştir. Teşkilat ve operasyon iddialarında komplo ve doğrulanmamış spekülasyondan kaçınılmış, dengeli ve olgusal bir çerçeve tutulmuştur. Hedef, bir SIGINT okuryazarına tarihsel ve kültürel bağlamı kazandırmaktır: tekniği, içine doğduğu yüz yıllık birikim ve etik sorumlulukla birlikte taşımak. Kendi ülkenin ve sürümünün mevzuatını teyit et; bu kitap anlama, savunma ve okuryazarlık içindir.

---

### Serinin diğer bölümleri (çapraz referans)

- SIGINT_00 — Başlangıç, İndeks ve Yasal Manifesto: serinin yasal/etik çerçevesi ve öğrenme yolu. (Bu bölümün etik duruşunun kökü oradadır.)
- SIGINT_01 — RF Fiziği, Spektrum ve Modülasyon: HF/iyonosfer yayılımı, dalga boyu, modülasyon. (Sayı istasyonu ve DXing'in HF yayılım fiziği oradadır.)
- SIGINT_03 — Anten Teorisi, Kapsama ve Yön Bulma (DF): yön bulma temelleri. (Burst transmission'a karşı DF ve kaçak verici avının tekniği orada.)
- SIGINT_05 — Hücresel, WiFi/BT ve IoT Spektrumu: RFID ve protokol bağlamı. (The Thing'in modern akrabası pasif RFID oradadır.)
- SIGINT_06 — TEMPEST, RF Sızıntısı ve Savunma: yan-kanal, emisyon güvenliği, OPSEC. (ANT katalog fenomeninin savunma dersi ve The Thing savunması oraya bağlanır.)
- SIGINT_07 — Disiplinler ve Sinyal Ayıklama: COMINT/ELINT/FISINT, trafik analizi, kütüphane eşleştirme, SEI. (RWR ve ELINT kütüphanesinin teknik mekaniği oradadır.)
- SIGINT_08 — Frekans Tahsisi ve Bant Planı: spektrum tahsisi ve düzenleme. (Sivil spektrum izleme ve Türkiye/BTK çerçevesi oraya bağlanır.)
- SIGINT_09 — Yer Tespiti, Yön Bulma ve Takip: geolocation ve takip. (Ticari ADS-B/AIS geolocation ve fox hunting oraya bağlanır.)
- SIGINT_13 — RF Tehdit Manzarası ve Karşı-Önlemler: EW taksonomisi (EA/EP/ES), karıştırma, savunma. (Bu bölümün EW doktrini başlığı, oranın savunma odaklı işlenişinin tarihsel/kültürel tamamlayıcısıdır.)
- SIGINT_14 — İstihbarat Kaynakları ve Takip: kaynak değerlendirme ve doğrulama. (Tarihsel olay araştırmasının kaynak hijyeni oraya bağlanır.)

> Kapanış: Sinyal istihbaratı, bir asır önce havaya verilen ilk telsiz mesajıyla başlayan, Room 40'ın el yazısı çözümlerinden Bletchley'in bombe'larına, soğuk bir gecenin sayı istasyonundan modern bilişsel elektronik harbe uzanan yaşayan bir tarihtir. Bu tarih iki şeyi birden öğretir: tekniğin gücünü ve o gücün sorumluluğunu. Bir analiz ekranındaki küme, bir histogram ya da HF'de duyulan bir sayı dizisi, yalnızca fizik değil; yüz yıllık bir insanlık birikiminin, kurumların, doktrinlerin ve etik tartışmaların izidir. Bu izi tanımak, tekniği hem daha derin hem daha sorumlu kılar.
>
> Bu doküman Kanije Kalesi güvenlik/teknik rehberleri koleksiyonunun SIGINT serisinin 21. bölümüdür. İlgili: SIGINT_00–14, `VERACRYPT_USTALIK_REHBERI.md`, `WINDOWS11_HARDENING_KALE.md`, `LINUX_HARDENING_KALE.md`.
