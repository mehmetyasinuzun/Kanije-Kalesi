# Güvenlik Politikası

Kanije Kalesi bir **güvenlik izleme** aracıdır; bu nedenle kendi güvenliğini de
ciddiye alır. Bu belge, aracın güvenlik modelini ve bir açık bulduğunuzda ne
yapmanız gerektiğini açıklar.

## Güvenlik Modeli

- **Yetkilendirme.** Bot yalnızca yapılandırılmış `chat_id` (ve isteğe bağlı
  `allowed_chat_ids`) listesindeki sohbetlerden gelen komutları işler. Diğer tüm
  mesajlar sessizce reddedilir ve loglanır.
- **Hassas veri.** Araç kamera fotoğrafı ve ekran görüntüsü çekebilir. Bu veriler
  yalnızca Telegram'a, yapılandırılmış sohbete gönderilir; varsayılan olarak yerel
  diske yazılmaz (`security.delete_captures_after_send = true`).
- **Bot token gizliliği.** Token tam yetki demektir. Config dosyası `0600`,
  config dizini `0700` izinleriyle yazılır. Token'ı tercihen `KANIJE_BOT_TOKEN`
  ortam değişkeniyle verin; repoya, loglara veya ekran görüntülerine sokmayın.
  `/ayarlar` çıktısında token maskelenir.
- **Mesaj güvenliği.** Telegram'a giden tüm dinamik alanlar (kullanıcı adı, cihaz
  etiketi, SSID, hata metni) HTML olarak kaçışlanır; enjeksiyon veya bozuk mesaj
  oluşmaz.
- **Tek örnek.** Aynı anda iki kopya çalışmasını engellemek için işletim sistemi
  düzeyinde kilit kullanılır.

## Sorumlu Kullanım

Bu araç **kendi cihazlarınızı** izlemek içindir. Başkalarının cihazlarına rızası
olmadan kurmak yasa dışı olabilir ve etik dışıdır. Kamera/ekran yakalama
özelliklerini yürürlükteki gizlilik yasalarına uygun kullanın.

## Açık Bildirimi

Bir güvenlik açığı bulursanız:

1. **Herkese açık issue açmayın.**
2. Depo sahibine özel olarak ulaşın (GitHub profili üzerinden) ve şunları paylaşın:
   - Etkilenen sürüm / commit
   - Yeniden üretme adımları
   - Olası etki
3. Düzeltme yayınlanana kadar detayları gizli tutun.

Makul bir sürede yanıt vermeye ve doğrulanan açıkları gidermeye çalışırız.

## Kapsam Dışı

- Kök/yönetici yetkisi gerektiren yerel saldırılar (araç zaten yüksek yetkiyle çalışır).
- Telegram'ın kendi altyapısına yönelik sorunlar.
- `app/` altındaki **kullanımdan kalkmış** Python sürümü (bkz. [app/DEPRECATED.md](app/DEPRECATED.md)).
