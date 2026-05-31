package telegram

import "strings"

// catItem is one command's presentation. A blank `menu` hides it from the
// Telegram "/" autocomplete menu (the command still works and still appears in
// /yardim) — used for rare or dangerous commands that shouldn't clutter the menu.
type catItem struct {
	cmd  string // canonical command, e.g. "/status"
	help string // one-line description shown in /yardim
	menu string // short label for the "/" menu ("" = hidden from menu)
}

// commandCatalog is the SINGLE SOURCE OF TRUTH for how commands are presented.
// Both /yardim (buildHelp) and the "/" menu (menuCommands) are generated from it,
// so the two can never drift out of sync again. Order here = display order.
//
// This catalog is presentation-only; authorization still lives in commandCaps and
// dispatch still routes in the bot's switch. Add a new command in BOTH places.
var commandCatalog = []struct {
	header string
	items  []catItem
}{
	{"📊 İzleme & Durum", []catItem{
		{"/status", "Sistem durumu — CPU, RAM, disk", "📊 Sistem durumu"},
		{"/pil", "Pil yüzdesi, şarj, kalan süre", "🔋 Pil durumu"},
		{"/olaylar", "Son olaylar (tip/sayı ile filtrele)", "📋 Son olaylar"},
		{"/ozet", "Son 7 günün olay özeti", "📈 7 gün özeti"},
		{"/defender", "Microsoft Defender durumu + taramalar", "🛡️ Defender durumu"},
		{"/dogrula", "Olay günlüğü bütünlüğünü doğrula", ""},
		{"/ping", "Bağlantı kontrolü", ""},
	}},
	{"📷 Anlık Komutlar", []catItem{
		{"/foto", "Kameradan anlık fotoğraf", "📷 Anlık fotoğraf"},
		{"/ekran", "Ekran görüntüsü", "🖥️ Ekran görüntüsü"},
		{"/seskayit", "Anlık ses kaydı — örn. <code>/seskayit 30</code>", "🎤 Anlık ses (X sn)"},
		{"/pano", "Panodaki metni getir", "📋 Pano içeriği"},
		{"/panik", "Tek komutta kanıt: foto + ekran + ses + IP", "🆘 Panik — kanıt topla"},
		{"/dinle", "Şimdi canlı dinle — örn. <code>/dinle 10</code>", "🎧 Canlı dinle (şimdi)"},
	}},
	{"🎬 Tetik Modları — sen yokken sürekli izler", []catItem{
		{"/tetikkamera", "Hareket görünce foto serisi — <code>/tetikkamera ac</code>", "🎥 Hareket-tetikli kamera"},
		{"/tetikses", "Sesi otomatik yakalar (sessizken durur) — <code>/tetikses ac</code>", "🔊 Ses-tetikli kayıt"},
	}},
	{"🛡️ Güvenlik & Koruma", []catItem{
		{"/koruma", "Fiziksel tehdit: dead-man · USB · yanlış-giriş", "🛡️ Koruma motoru"},
		{"/tuzak", "Honeypot — tuzağa dokunan yakalanır", "🍯 Tuzak (honeypot)"},
		{"/kilit", "Ekranı kilitle · <code>/kilit tam</code> = lockdown", "🔒 Ekranı kilitle"},
		{"/erisim", "Dosyana kim erişti — <code>/erisim kur &lt;yol&gt;</code>", "👁️ Dosya erişim izi"},
	}},
	{"⚙️ Sistem & Dosya", []catItem{
		{"/dosya", "Dosya gez / indir — <code>/dosya al &lt;yol&gt;</code>", "📁 Dosya gez/indir"},
		{"/zamanla", "Komut zamanla — <code>/zamanla 30dk /foto</code>", "⏰ Komut zamanla"},
		{"/terminal", "Uzak komut çalıştır", "💻 Uzak komut"},
		{"/terminalix", "Yönetici (admin) terminali", ""},
		{"/yeniden", "Sistemi yeniden başlat", "🔄 Yeniden başlat"},
		{"/kapat", "Sistemi kapat", "⏻ Kapat"},
		{"/guncelle", "Yeni sürümü kontrol et ve kur", "⬆️ Güncelle"},
		{"/iptal", "Bekleyen / devam eden işlemi iptal et", ""},
	}},
	{"⚠️ Sahip — Geri Dönüşsüz", []catItem{
		{"/kaldir", "Kanije'yi izsiz kaldır (görev + dosya + eski sürümler)", ""},
		{"/aktar", "Botu yeni sahibe devret — <code>/aktar &lt;chat_id&gt;</code>", ""},
		{"/imha", "Kanije verisi + Müzik klasörü içini güvenli sil (OS silinmez)", ""},
	}},
	{"🔧 Ayarlar & Yardım", []catItem{
		{"/rehber", "Kullanım rehberi — nasıl çalışır, hangi modlar açık kalsın", "📖 Kullanım rehberi"},
		{"/kurulum", "Etkileşimli kurulum menüsü", "⚙️ Kurulum menüsü"},
		{"/ayarlar", "Mevcut yapılandırmayı gör", "🔧 Ayarlar"},
		{"/yardim", "Bu komut listesi", "📜 Komut listesi"},
	}},
	{"👥 Fleet & Kullanıcılar", []catItem{
		{"/cihazlar", "Tüm cihazları listele", "🛰️ Cihazlar"},
		{"/ekle", "Kişi ekle — <code>/ekle &lt;chat_id&gt; [isim]</code>", "➕ Kişi ekle"},
		{"/yonetim", "Eklediklerini gör / düzenle / çıkar", "👥 Kişi yönetimi"},
		{"/loglar", "Son işlemler (kim ne yaptı)", "🧾 İşlem günlüğü"},
	}},
}

// FormatHelp renders the /yardim command list from the catalog.
func FormatHelp() string {
	var b strings.Builder
	b.WriteString("🏰 <b>Kanije Kalesi — Komutlar</b>\n")
	b.WriteString("📖 <i>İlk kez mi, ya da hangi modlar açık kalsın? → /rehber</i>\n")

	for _, cat := range commandCatalog {
		b.WriteString("\n<b>")
		b.WriteString(cat.header)
		b.WriteString("</b>\n")
		for _, it := range cat.items {
			b.WriteString(it.cmd)
			b.WriteString(" — ")
			b.WriteString(it.help)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n<b>🔔 Otomatik Bildirimler</b> <i>(komut gerekmez)</i>\n")
	b.WriteString("Giriş/çıkış · ekran kilit-açılış · USB ve <b>her aygıt</b> (fare/klavye/telefon) tak-çıkar · kurcalama — hepsi kendiliğinden bildirilir.\n")
	b.WriteString("\n<i>💡 Grupta belirli cihaza: <code>/foto dizustu</code> · Adım adım kullanım → /rehber · Her komut yetkine bağlıdır.</i>")
	return b.String()
}

// menuCommands builds the Telegram "/" autocomplete menu from the catalog
// (only items with a non-empty menu label).
func menuCommands() []BotCommand {
	cmds := make([]BotCommand, 0, len(commandCatalog)*4)
	for _, cat := range commandCatalog {
		for _, it := range cat.items {
			if it.menu == "" {
				continue
			}
			cmds = append(cmds, BotCommand{strings.TrimPrefix(it.cmd, "/"), it.menu})
		}
	}
	return cmds
}
