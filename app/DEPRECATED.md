# ⚠️ Bu Sürüm Kullanımdan Kalkmıştır (DEPRECATED)

Bu `app/` dizinindeki **Python sürümü artık geliştirilmiyor ve desteklenmiyor.**

Yerine, tamamen yeniden yazılmış ve aktif olarak bakımı yapılan **Go sürümünü**
kullanın:

➡️ **[`../go/`](../go/)** — daha hızlı, tek bağımsız binary, CGo'suz cross-compile
(Windows / Linux / Raspberry Pi), tüm yapılandırma Telegram üzerinden.

## Neden Go'ya geçildi?

| | Python (`app/`) | Go (`go/`) |
|-|-----------------|------------|
| Dağıtım | Python + bağımlılıklar gerekir | Tek bağımsız binary |
| Cross-compile | Zor | `make build-all` ile tek komut |
| Bellek / başlangıç | Ağır | Hafif |
| Bakım durumu | **Durduruldu** | **Aktif** |

Bu klasör yalnızca **arşiv / referans** amacıyla tutulmaktadır. Yeni kurulumlar ve
katkılar Go sürümüne yapılmalıdır. Güvenlik düzeltmeleri ve yeni özellikler
yalnızca [`../go/`](../go/) için yayınlanır.
