# 🏰 LİNUX — "KALE" MİMARİSİ: TAM SPEKTRUM SERTLEŞTİRME REHBERİ
## Kali Linux · Ubuntu Öncelikli

> **Hedef:** BIOS donanım kilidinden uygulama süreç yalıtımına kadar her katmanı kilitleyen, fiziksel çalınma dahil tüm senaryolara dayanıklı Linux sistemi inşa etmek.

> [!NOTE]
> 🐉 **Kali Linux** saldırı aracı olarak tasarlanmıştır; varsayılan yapılandırması savunmaya uygun değildir. Bu rehber Kali'yi defensif bir makineye dönüştürür.
> 🟠 **Ubuntu** daha güvenli varsayılan gelir (AppArmor, UFW); bu rehber onu daha da sertleştirir.

---

## 📐 MİMARİ HARITA

```
┌─────────────────────────────────────────────────────────┐
│  KATMAN 0 — Donanım / BIOS/UEFI                        │
│  Admin Şifresi · Secure Boot · TPM · IOMMU             │
├─────────────────────────────────────────────────────────┤
│  KATMAN 1 — Önyükleme (Boot)                           │
│  GRUB Şifresi · LUKS2 Tam Disk Şifreleme · Secure Boot │
├─────────────────────────────────────────────────────────┤
│  KATMAN 2 — Çekirdek (Kernel)                          │
│  sysctl · kernel.modules · Secure Boot imzası          │
├─────────────────────────────────────────────────────────┤
│  KATMAN 3 — Bellek (RAM)                               │
│  Swap Şifreleme · Hibernate · DMA Koruması             │
├─────────────────────────────────────────────────────────┤
│  KATMAN 4 — Süreç / Uygulama                           │
│  AppArmor · SELinux · seccomp · Capabilities           │
├─────────────────────────────────────────────────────────┤
│  KATMAN 5 — Ağ                                         │
│  UFW/iptables · SSH · DNS-over-HTTPS · fail2ban        │
├─────────────────────────────────────────────────────────┤
│  KATMAN 6 — Kimlik / Hesap                             │
│  PAM · sudo · şifre politikası · SSH key               │
├─────────────────────────────────────────────────────────┤
│  KATMAN 7 — Denetim / İzleme                           │
│  auditd · AIDE · syslog · rootkit tarama               │
└─────────────────────────────────────────────────────────┘
```

---

## KATMAN 0 — Donanım ve BIOS/UEFI Kilitleme

Windows kılavuzuyla aynı adımlar — Linux kurulumundan önce yapılmalı:

```
BIOS → Security → Supervisor Password → Güçlü şifre
BIOS → Security → Boot Password → Güçlü şifre
BIOS → Secure Boot → Enabled (Standard mod)
BIOS → TPM → Enabled (Firmware TPM)
BIOS → VT-d / AMD IOMMU → Enabled  ← DMA koruması için ZORunlu
BIOS → Boot Order → USB en alta al
```

> [!CAUTION]
> Kali kurulumunda Secure Boot ön tanımlı olarak **devre dışı** bırakmak gerekebilir; bazı Kali kernel modülleri ve araçları imzasız gelir. Ubuntu ise `shim` ile Secure Boot destekler ve imzalı kernel sunar.

---

## KATMAN 1 — Önyükleme Kilidi

### 1.1 — LUKS2 Tam Disk Şifreleme (Kurulum Sırasında)

Linux'ta BitLocker'ın karşılığı **LUKS2** (Linux Unified Key Setup ver.2).

#### Ubuntu Kurulumunda:
```
Ubuntu Kurulum Sihirbazı → "Diske Yükle" aşaması
  → "Gelişmiş Özellikler" → ✅ LVM ile şifreleme kullan
  → Bir şifreleme parolası oluştur (Pre-Boot — GRUB sonrası sorulur)
  → Yükle
```

#### Kali Linux Kurulumunda:
```
Kali Installer → Disk bölümlendirme
  → "Kılavuzlu - tüm diski kullan ve şifreli LVM kur"
  → Şifreleme parolası gir
```

#### Manuel LUKS2 Kurulumu (Gelişmiş Kontrol):
```bash
# Disk seç (dikkat: tüm veri silinir)
DISK="/dev/nvme0n1"

# GPT tablosu oluştur
parted -s $DISK mklabel gpt

# EFI bölümü (512 MB)
parted -s $DISK mkpart ESP fat32 1MiB 513MiB
parted -s $DISK set 1 esp on

# LUKS bölümü (kalan tüm alan)
parted -s $DISK mkpart primary 513MiB 100%

# LUKS2 ile şifrele — AES-256-XTS, Argon2id anahtar türetme
cryptsetup luksFormat \
    --type luks2 \
    --cipher aes-xts-plain64 \
    --key-size 512 \
    --hash sha512 \
    --pbkdf argon2id \
    --pbkdf-memory 1048576 \
    --pbkdf-parallel 4 \
    --iter-time 5000 \
    ${DISK}p2

# Şifreli bölümü aç
cryptsetup open ${DISK}p2 cryptroot

# LVM kur (isteğe bağlı ama önerilir)
pvcreate /dev/mapper/cryptroot
vgcreate vg0 /dev/mapper/cryptroot
lvcreate -L 50G vg0 -n root
lvcreate -L 8G  vg0 -n swap
lvcreate -l 100%FREE vg0 -n home

# Dosya sistemleri
mkfs.ext4 /dev/vg0/root
mkfs.ext4 /dev/vg0/home
mkswap /dev/vg0/swap
mkfs.fat -F32 ${DISK}p1
```

**LUKS Header Yedekle (48 Haneli Anahtar Gibi):**
```bash
# Header yedekle — ayrı bir USB'ye koy, güvenli sakla
cryptsetup luksHeaderBackup /dev/nvme0n1p2 \
    --header-backup-file /mnt/usb/luks_header_$(hostname)_$(date +%Y%m%d).img

# Header durumunu görüntüle
cryptsetup luksDump /dev/nvme0n1p2
```

> [!CAUTION]
> LUKS header bozulursa disk **tamamen erişilemez** olur. Header yedeği, BitLocker kurtarma anahtarı kadar kritik — fiziksel kasada sakla.

### 1.2 — GRUB Önyükleyici Şifresi

GRUB şifresi olmadan herhangi biri `e` tuşuna basıp kernel parametrelerini değiştirebilir (`init=/bin/bash` → root shell).

```bash
# GRUB şifre hash'i oluştur
grub-mkpasswd-pbkdf2
# İstediğin şifreyi gir → hash çıktısı kopyala

# GRUB yapılandırmasına ekle
sudo nano /etc/grub.d/40_custom
```

```bash
# /etc/grub.d/40_custom içine ekle:
set superusers="grubadmin"
password_pbkdf2 grubadmin grub.pbkdf2.sha512.10000.HASH_BURAYA...

# GRUB'u güncelle
sudo update-grub

# Doğrulama: Yeniden başlatınca "e" tuşu şifre isteyecek
```

### 1.3 — Secure Boot (Ubuntu)

Ubuntu, `shim` ile Secure Boot desteği sağlar. Kali'de ekstra adım gerekir.

```bash
# Ubuntu — Secure Boot durumu
mokutil --sb-state
# SecureBoot enabled → Doğru

# Kali — İmzalı kernel kurulumu (isteğe bağlı)
# Kali birçok modülü imzasız kullandığından Secure Boot kısıtlı çalışır
# Pentesting araçlarını kırmaması için "Custom Key" yöntemi:
sudo apt install -y sbsigntool mokutil shim-signed
```

---

## KATMAN 2 — Kernel Sertleştirme

### 2.1 — sysctl Güvenlik Parametreleri

```bash
# /etc/sysctl.d/99-hardening.conf dosyasını oluştur
sudo tee /etc/sysctl.d/99-hardening.conf <<'EOF'
# ──── Çekirdek ve İşaretçi Koruması ────
kernel.kptr_restrict = 2              # Kernel pointer'larını gizle (ayrıcalıksız kullanıcıdan)
kernel.dmesg_restrict = 1             # dmesg'i root'a kısıtla
kernel.perf_event_paranoid = 3        # perf event erişimini kısıtla
kernel.unprivileged_bpf_disabled = 1  # eBPF'i ayrıcalıksız kullanıcıdan kısıtla
net.core.bpf_jit_harden = 2          # BPF JIT sertleşmesi

# ──── Bellek / Stack Koruması ────
kernel.randomize_va_space = 2         # ASLR tam (2 = Stack + Heap + Mmap)
vm.mmap_rnd_bits = 32                 # Mmap randomizasyon biti sayısı
vm.mmap_min_addr = 65536              # NULL pointer dereference koruması

# ──── Çekirdek Modülü Kilitleme ────
kernel.modules_disabled = 0           # Kurulumda 0, sonra 1 yapılır (aşağıya bak)

# ──── Kullanıcı Alanı İzolasyonu ────
kernel.yama.ptrace_scope = 2          # ptrace yalnızca root (2) veya parent (1)
                                      # Mimikatz eşdeğeri hafıza okuma engeli

# ──── Ağ Sertleştirme — Spoofing ve MITM ────
net.ipv4.conf.all.rp_filter = 1              # Reverse path filtering
net.ipv4.conf.default.rp_filter = 1
net.ipv4.conf.all.accept_redirects = 0       # ICMP redirect kabul etme
net.ipv4.conf.default.accept_redirects = 0
net.ipv6.conf.all.accept_redirects = 0
net.ipv4.conf.all.send_redirects = 0         # ICMP redirect gönderme
net.ipv4.conf.all.accept_source_route = 0    # Kaynak yönlendirme reddet
net.ipv6.conf.all.accept_source_route = 0
net.ipv4.conf.all.log_martians = 1           # Sahte/geçersiz paketleri logla
net.ipv4.icmp_echo_ignore_broadcasts = 1     # Smurf saldırısı koruması
net.ipv4.tcp_syncookies = 1                  # SYN flood koruması
net.ipv4.tcp_rfc1337 = 1                     # TIME_WAIT saldırısı koruması
net.ipv4.conf.all.secure_redirects = 1       # Yalnızca ağ geçidinden redirect
net.ipv6.conf.all.accept_ra = 0              # IPv6 Router Advertisement reddet (MITM)
net.ipv6.conf.default.accept_ra = 0

# ──── IPv6 (Kullanmıyorsan kapat) ────
net.ipv6.conf.all.disable_ipv6 = 0           # Kullanıyorsan 0 bırak, kullanmıyorsan 1
EOF

# Uygula
sudo sysctl --system

# Doğrulama:
sudo sysctl kernel.randomize_va_space   # 2 olmalı
sudo sysctl kernel.kptr_restrict        # 2 olmalı
sudo sysctl net.ipv4.tcp_syncookies     # 1 olmalı
```

### 2.2 — Kernel Modülü Kısıtlaması

Kurulum tamamlandıktan sonra yeni modül yüklenmesini engelle:

```bash
# Kullanılmayan tehlikeli modülleri kara listeye al
sudo tee /etc/modprobe.d/blacklist-hardening.conf <<'EOF'
# Thunderbolt DMA
blacklist thunderbolt

# Eski güvensiz protokoller
blacklist dccp
blacklist sctp
blacklist rds
blacklist tipc

# USB uzak erişim protokolleri
blacklist usb-storage       # Sadece USB bellek kullanmıyorsan
# install usb-storage /bin/false   # Tamamen engelle

# Rare filesystems — saldırı yüzeyi
blacklist cramfs
blacklist freevxfs
blacklist jffs2
blacklist hfs
blacklist hfsplus
blacklist squashfs
blacklist udf

# FireWire — DMA saldırı vektörü
blacklist firewire-core
blacklist firewire-ohci

# Bluetooth (kullanmıyorsan)
# blacklist bluetooth
# blacklist btusb
EOF

# Güncelle
sudo update-initramfs -u

# Kurulum tamamsa yeni modül yüklemeyi kapat (geri dönüşü zor — dikkatli)
# echo "kernel.modules_disabled=1" | sudo tee -a /etc/sysctl.d/99-hardening.conf
# sudo sysctl -w kernel.modules_disabled=1
```

### 2.3 — Kernel Komut Satırı Sertleştirmesi

```bash
sudo nano /etc/default/grub
```

```bash
# GRUB_CMDLINE_LINUX satırını düzenle:
GRUB_CMDLINE_LINUX="quiet splash \
    apparmor=1 security=apparmor \
    init_on_alloc=1 init_on_free=1 \
    page_alloc.shuffle=1 \
    randomize_kstack_offset=on \
    vsyscall=none \
    debugfs=off \
    oops=panic \
    module.sig_enforce=1 \
    lockdown=confidentiality \
    iommu=force \
    intel_iommu=on \
    amd_iommu=on \
    pti=on \
    spectre_v2=on \
    spec_store_bypass_disable=on \
    l1tf=full,force \
    mds=full,nosmt"
```

```bash
# Uygula
sudo update-grub

# Parametreler:
# init_on_alloc/free=1  → Bellek tahsisi/serbest bırakma sırasında sıfırla (Cold Boot)
# vsyscall=none         → vsyscall saldırı yüzeyi kapat
# lockdown=confidentiality → Kernel lockdown (root bile bazı işlemleri yapamaz)
# pti=on               → Meltdown koruması (Page Table Isolation)
# iommu=force          → DMA koruması
```

---

## KATMAN 3 — RAM ve Swap Güvenliği

### 3.1 — Şifreli Swap

LUKS disk şifrelemesi aktifse swap otomatik şifreli gelir (LVM üzerindeyse). Değilse:

```bash
# Mevcut swap'ı kontrol et
cat /proc/swaps
swapon --show

# Swap dosyası yerine şifreli swap bölümü kullan
# /etc/crypttab'a ekle (her açılışta rastgele anahtar — en güvenli)
sudo tee -a /etc/crypttab <<'EOF'
cryptswap /dev/sda2 /dev/urandom swap,cipher=aes-xts-plain64,size=512,hash=sha512
EOF

# /etc/fstab'da swap satırını güncelle
# /dev/sda2 yerine /dev/mapper/cryptswap kullan
echo "/dev/mapper/cryptswap none swap sw 0 0" | sudo tee -a /etc/fstab
```

### 3.2 — Hibernate (Hazırda Bekletme) Kapatma

```bash
# Hibernate RAM içeriğini diske yazar — DİSK ŞİFRELİ olsa bile risk:
sudo systemctl mask sleep.target suspend.target hibernate.target hybrid-sleep.target

# Doğrulama:
systemctl status sleep.target
# Loaded: masked → Doğru

# Logind'e de uygula
sudo sed -i 's/#HandleSuspendKey=suspend/HandleSuspendKey=ignore/' /etc/systemd/logind.conf
sudo sed -i 's/#HandleHibernateKey=hibernate/HandleHibernateKey=ignore/' /etc/systemd/logind.conf
sudo sed -i 's/#HandleLidSwitch=suspend/HandleLidSwitch=lock/' /etc/systemd/logind.conf
sudo systemctl restart systemd-logind
```

### 3.3 — Bellek Temizleme (Kapatma Sırasında)

```bash
# systemd-shutdown hook — kapatma sırasında swap sıfırla
sudo tee /etc/systemd/system/clear-swap.service <<'EOF'
[Unit]
Description=Clear swap at shutdown
DefaultDependencies=no
Before=shutdown.target reboot.target halt.target
RequiresMountsFor=/

[Service]
Type=oneshot
ExecStart=/sbin/swapoff -a
ExecStart=/bin/sync

[Install]
WantedBy=halt.target reboot.target shutdown.target
EOF

sudo systemctl enable clear-swap.service
```

### 3.4 — Çekirdek DMA Koruması

IOMMU kernel parametresi (`iommu=force`) ile Katman 2'de etkinleştirildi. Doğrulama:

```bash
# IOMMU aktif mi?
dmesg | grep -i iommu | head -20
# "AMD-Vi: IOMMU enabled" veya "Intel-IOMMU: enabled" görünmeli

# DMAR tablosu
sudo cat /sys/kernel/security/lockdown
# confidentiality → Doğru (kernel lockdown aktif)
```

---

## KATMAN 4 — Süreç ve Uygulama Yalıtımı

### 4.1 — AppArmor (Ubuntu Varsayılan, Kali'ye Kur)

```bash
# Ubuntu — AppArmor durumu
sudo aa-status

# Kali — AppArmor kur ve etkinleştir
sudo apt install -y apparmor apparmor-utils apparmor-profiles apparmor-profiles-extra
sudo systemctl enable --now apparmor

# Tüm profilleri enforce moduna al
sudo aa-enforce /etc/apparmor.d/*

# Yeni uygulama için profil oluştur (örn: firefox)
sudo aa-genprof firefox
# Uygulamayı çalıştır, AppArmor erişimleri öğrenir → kaydet

# Mevcut profil durumu
sudo aa-status | grep "profiles are in enforce mode"
```

**Kritik Uygulamalar İçin AppArmor Profil Zorlaması:**
```bash
# Mevcut profilleri enforce et
sudo aa-enforce /etc/apparmor.d/usr.sbin.cupsd
sudo aa-enforce /etc/apparmor.d/usr.bin.firefox
sudo aa-enforce /etc/apparmor.d/usr.sbin.sshd
sudo aa-enforce /etc/apparmor.d/usr.sbin.apache2

# Tüm complain modundakileri enforce'a al
for prof in $(sudo aa-status | grep "profiles are in complain mode" -A 100 | grep "^   " | awk '{print $1}'); do
    sudo aa-enforce "$prof" 2>/dev/null
done
```

### 4.2 — seccomp (Sistem Çağrısı Filtresi)

```bash
# Docker veya direkt süreçler için seccomp profili:
# Örnek: Tehlikeli syscall'ların engellenmesi

# Mevcut süreçlerin seccomp durumu
cat /proc/$$/status | grep Seccomp
# Seccomp: 2 → filtre aktif, 0 → aktif değil

# systemd servislerinde seccomp filtresi (örn özel servis):
sudo tee /etc/systemd/system/guvenli-uygulama.service <<'EOF'
[Service]
ExecStart=/usr/bin/uygulama
# Tehlikeli sistem çağrılarını engelle
SystemCallFilter=~@clock @cpu-emulation @debug @keyring @module @mount @obsolete @privileged @raw-io @reboot @swap @sync
NoNewPrivileges=yes
PrivateTmp=yes
PrivateDevices=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/var/lib/uygulama
EOF
```

### 4.3 — Capabilities Kısıtlaması

```bash
# Tüm binary'lerin gereksiz capabilities'ini kaldır
# Önce mevcut durumu göster
sudo getcap -r / 2>/dev/null

# Gereksiz capabilities'i kaldır (örn ping için SETUID yerine cap kullanabilirsin)
# Hassas: net_raw capability
sudo setcap -r /usr/bin/ping 2>/dev/null || true

# SUID/SGID bit taraması — gereksiz olanları belirle
sudo find / -type f \( -perm -4000 -o -perm -2000 \) 2>/dev/null | sort

# Gereksiz SUID bitlerini kaldır (dikkatli — bazıları gerekli)
# Örnek — su yerine sudo kullanıyorsan:
# sudo chmod u-s /bin/su
```

### 4.4 — Firejail ile Uygulama Sandbox

```bash
# Kur
sudo apt install -y firejail firejail-profiles

# Kullanım
firejail firefox                 # Firefox'u sandboxla
firejail --private chromium      # Geçici ev diziniyle çalıştır
firejail --net=none vlc          # Ağsız çalıştır

# Varsayılan uygulama olarak ayarla — her çalışmada sandbox
sudo ln -sf /usr/bin/firejail /usr/local/bin/firefox

# Profil listesi
ls /etc/firejail/*.profile | head -20
```

---

## KATMAN 5 — Ağ Güvenliği

### 5.1 — UFW (Ubuntu) / iptables Güvenlik Duvarı

```bash
# UFW — Ubuntu
sudo apt install -y ufw

# Varsayılan politikalar
sudo ufw default deny incoming
sudo ufw default deny outgoing    # Whitelist yaklaşımı
sudo ufw default deny forward

# İzin verilecekler (beyaz liste)
sudo ufw allow out 53/udp   comment "DNS"
sudo ufw allow out 53/tcp   comment "DNS TCP"
sudo ufw allow out 443/tcp  comment "HTTPS"
sudo ufw allow out 80/tcp   comment "HTTP (gerekirse)"
sudo ufw allow out 123/udp  comment "NTP"

# SSH yalnızca belirli IP'den (gerekliyse)
sudo ufw allow from 192.168.1.0/24 to any port 22 proto tcp

sudo ufw enable
sudo ufw status verbose
```

```bash
# nftables (modern — Ubuntu 22.04+ ve Kali)
sudo tee /etc/nftables.conf <<'EOF'
#!/usr/sbin/nft -f
flush ruleset

table inet filter {
    chain input {
        type filter hook input priority 0; policy drop;
        ct state established,related accept
        iif "lo" accept
        ct state invalid drop
        tcp flags syn limit rate 20/second burst 5 packets accept
        # SSH (sadece LAN'dan)
        ip saddr 192.168.1.0/24 tcp dport 22 accept
        drop
    }
    chain output {
        type filter hook output priority 0; policy drop;
        ct state established,related accept
        iif "lo" accept
        udp dport {53, 123} accept
        tcp dport {80, 443, 53} accept
        drop
    }
    chain forward {
        type filter hook forward priority 0; policy drop;
    }
}
EOF

sudo systemctl enable --now nftables
```

### 5.2 — SSH Sertleştirme

SSH, Linux'un en büyük saldırı vektörü. Kapsamlı sertleştirme:

```bash
sudo tee /etc/ssh/sshd_config.d/99-hardening.conf <<'EOF'
# ──── Kimlik Doğrulama ────
PasswordAuthentication no           # Parola girişini KAPAT — yalnızca anahtar
PubkeyAuthentication yes            # Anahtar tabanlı giriş
PermitRootLogin no                  # Root ile doğrudan SSH girişi yasak
AuthorizedKeysFile .ssh/authorized_keys
MaxAuthTries 3                      # 3 başarısız denemede bağlantıyı kes
LoginGraceTime 20                   # 20 saniye içinde giriş yapılmazsa kes
PermitEmptyPasswords no             # Boş şifreli hesaplara giriş yasak

# ──── Bağlantı Güvenliği ────
Protocol 2                          # Yalnızca SSH protokol 2
Port 2222                           # Varsayılan portu değiştir (opsiyonel)
AddressFamily inet                  # Yalnızca IPv4 (IPv6 kullanmıyorsan)
ListenAddress 0.0.0.0

# ──── Şifreleme Algoritmaları ────
KexAlgorithms curve25519-sha256,curve25519-sha256@libssh.org,diffie-hellman-group16-sha512,diffie-hellman-group18-sha512
Ciphers chacha20-poly1305@openssh.com,aes256-gcm@openssh.com,aes128-gcm@openssh.com
MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com

# ──── Oturum Güvenliği ────
ClientAliveInterval 300             # 5 dakika inaktivite → bağlantıyı kes
ClientAliveCountMax 2
MaxSessions 3
X11Forwarding no                    # X11 forwarding kapat (saldırı vektörü)
AllowAgentForwarding no
AllowTcpForwarding no               # TCP tünelleme kapat (proxychains engeli)
PermitTunnel no
GatewayPorts no
Banner /etc/ssh/banner.txt          # Uyarı banneri

# ──── Loglama ────
LogLevel VERBOSE
SyslogFacility AUTH
EOF

# SSH uyarı banneri
sudo tee /etc/ssh/banner.txt <<'EOF'
╔════════════════════════════════════════╗
║  YETKİSİZ ERİŞİM YASAKTIR             ║
║  Tüm bağlantılar kaydedilmektedir.    ║
╚════════════════════════════════════════╝
EOF

# SSH anahtarı oluştur (Ed25519 — en güvenli)
ssh-keygen -t ed25519 -a 100 -f ~/.ssh/id_ed25519 -C "$(whoami)@$(hostname)-$(date -I)"

# Servisi yeniden başlat
sudo systemctl restart sshd
sudo sshd -t  # Yapılandırma hatası var mı?
```

### 5.3 — fail2ban (Brute-Force Koruması)

```bash
sudo apt install -y fail2ban

sudo tee /etc/fail2ban/jail.local <<'EOF'
[DEFAULT]
bantime  = 3600          # 1 saat ban
findtime = 600           # 10 dakika içinde
maxretry = 3             # 3 başarısız deneme → ban
banaction = ufw          # UFW ile entegre
ignoreip = 127.0.0.1/8 192.168.1.0/24

[sshd]
enabled  = true
port     = 2222
filter   = sshd
logpath  = /var/log/auth.log
maxretry = 3

[nginx-http-auth]
enabled = false

[apache-auth]
enabled = false
EOF

sudo systemctl enable --now fail2ban

# Durum
sudo fail2ban-client status sshd
sudo fail2ban-client status
```

### 5.4 — DNS over HTTPS / TLS

```bash
# systemd-resolved ile DoT (DNS over TLS)
sudo tee -a /etc/systemd/resolved.conf <<'EOF'
[Resolve]
DNS=9.9.9.9#dns.quad9.net 1.1.1.1#cloudflare-dns.com
FallbackDNS=149.112.112.112#dns.quad9.net
DNSSEC=yes
DNSOverTLS=yes
EOF

sudo systemctl restart systemd-resolved

# Doğrulama
resolvectl status
resolvectl query example.com
```

### 5.5 — Kullanılmayan Servisleri Kapat

```bash
# Çalışan servisleri listele
systemctl list-units --type=service --state=running

# Kali'de genellikle kapatılabilecekler:
sudo systemctl disable --now avahi-daemon    # mDNS — ağ keşfi vektörü
sudo systemctl disable --now cups            # Yazıcı servisi — kullanmıyorsan
sudo systemctl disable --now bluetooth       # Bluetooth — kullanmıyorsan
sudo systemctl disable --now rpcbind         # NFS RPC — kullanmıyorsan
sudo systemctl disable --now nfs-server      # NFS server
sudo systemctl disable --now telnet          # Telnet — hiçbir zaman kullanma

# Ağ üzerinde dinleyen portları listele
sudo ss -tlnp
sudo netstat -tlnp 2>/dev/null || sudo ss -tlnp
```

---

## KATMAN 6 — Kimlik ve Hesap Güvenliği

### 6.1 — Root Hesabı Kilitleme

```bash
# Root'un doğrudan giriş yapmasını engelle
sudo passwd -l root

# Root'un shell'ini kaldır
sudo usermod -s /sbin/nologin root

# Doğrulama
sudo passwd -S root
# root L (Locked)
```

### 6.2 — Kali'ye Özgü: Default Root'tan Standart Kullanıcıya Geçiş

> [!IMPORTANT]
> Kali Linux 2020'den itibaren varsayılan olarak standart kullanıcıyla gelir. Eski Kali sürümlerinde root olarak çalışmak son derece tehlikelidir.

```bash
# Standart kullanıcı oluştur (Kali)
sudo adduser kullanici_adi
sudo usermod -aG sudo kullanici_adi

# Kali araçlarına erişim için gruba ekle
sudo usermod -aG kali kullanici_adi

# Root ile otomatik girişi kapat
sudo sed -i 's/^autologin-user=root/autologin-user=/' /etc/lightdm/lightdm.conf 2>/dev/null || true
```

### 6.3 — sudo Kısıtlaması

```bash
# sudo log tüm komutları kaydetsin
sudo tee /etc/sudoers.d/hardening <<'EOF'
# Tüm sudo kullanımını logla
Defaults logfile=/var/log/sudo.log
Defaults log_input, log_output
Defaults!NOPASSWD: ALL
Defaults passwd_timeout=1
Defaults timestamp_timeout=5    # 5 dakika sonra şifreyi tekrar sor

# Root shell açan komutları kısıtla
Defaults !visiblepw
Defaults env_reset
Defaults secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
EOF

sudo chmod 440 /etc/sudoers.d/hardening
sudo visudo -c  # Sözdizimi doğrulama
```

### 6.4 — PAM Şifre Politikası

```bash
sudo apt install -y libpam-pwquality

sudo tee /etc/security/pwquality.conf <<'EOF'
minlen = 16          # Minimum 16 karakter
minclass = 4         # Büyük, küçük, rakam, özel karakter zorunlu
maxrepeat = 2        # Aynı karakteri en fazla 2 kez tekrarla
gecoscheck = 1       # Kullanıcı adını şifreden engelle
difok = 8            # Eski şifreyle en az 8 karakter fark
EOF

# PAM hesap kilitleme — 5 başarısız deneme
sudo tee /etc/security/faillock.conf <<'EOF'
deny = 5
unlock_time = 900
fail_interval = 900
audit
EOF
```

### 6.5 — Parola Son Kullanma Tarihi

```bash
# Tüm mevcut kullanıcılar için
sudo chage -M 90 -m 1 -W 14 kullanici_adi
# -M 90: 90 günde bir değiştir
# -m 1:  en erken 1 gün sonra değiştirebilir
# -W 14: 14 gün önce uyar

# Varsayılan politika (yeni kullanıcılar için)
sudo sed -i 's/^PASS_MAX_DAYS.*/PASS_MAX_DAYS 90/'  /etc/login.defs
sudo sed -i 's/^PASS_MIN_DAYS.*/PASS_MIN_DAYS 1/'   /etc/login.defs
sudo sed -i 's/^PASS_WARN_AGE.*/PASS_WARN_AGE 14/'  /etc/login.defs
```

---

## KATMAN 7 — Denetim, İzleme ve Bütünlük

### 7.1 — auditd (Linux Denetim Çerçevesi)

```bash
sudo apt install -y auditd audispd-plugins

# Kural seti — önemli olayları kaydet
sudo tee /etc/audit/rules.d/99-hardening.rules <<'EOF'
# Mevcut kuralları temizle
-D

# Arabellek boyutu
-b 8192

# Başarısız sistem çağrıları
-a always,exit -F arch=b64 -S all -F success=0 -k failures

# Kimlik değişiklikleri
-w /etc/passwd    -p wa -k identity
-w /etc/shadow    -p wa -k identity
-w /etc/group     -p wa -k identity
-w /etc/sudoers   -p wa -k actions
-w /etc/sudoers.d -p wa -k actions

# Kullanıcı giriş/çıkış
-w /var/log/faillog -p wa -k logins
-w /var/log/lastlog -p wa -k logins
-w /var/run/faillock -p wa -k logins

# Oturum izleme
-w /var/run/utmp -p wa -k session
-w /var/log/wtmp -p wa -k session
-w /var/log/btmp -p wa -k session

# Dosya silme
-a always,exit -F arch=b64 -S unlinkat -S rename -k delete

# sudo kullanımı
-w /usr/bin/sudo -p x -k priv_esc

# SSH
-w /etc/ssh/sshd_config -p wa -k sshd

# Kernel modülü
-w /sbin/insmod  -p x -k modules
-w /sbin/rmmod   -p x -k modules
-w /sbin/modprobe -p x -k modules
-a always,exit -F arch=b64 -S init_module -S finit_module -k modules

# Ağ yapılandırması
-a always,exit -F arch=b64 -S sethostname -S setdomainname -k network

# Mount
-a always,exit -F arch=b64 -S mount -k mounts

# Crontab değişiklikleri
-w /etc/cron.d    -p wa -k cron
-w /etc/crontab   -p wa -k cron
-w /var/spool/cron -p wa -k cron

# Kural setini kilitle — yeniden başlatmadan değiştirilemez
-e 2
EOF

sudo systemctl enable --now auditd
sudo auditctl -l   # Kuralları listele

# Log izleme
sudo ausearch -k identity | aureport -f -i
sudo ausearch -k priv_esc
```

### 7.2 — AIDE (Dosya Bütünlük İzleme)

AIDE, kritik sistem dosyalarının hash değerlerini kaydeder ve değişiklik olursa uyarır.

```bash
sudo apt install -y aide aide-common

# AIDE yapılandırması
sudo tee -a /etc/aide/aide.conf <<'EOF'
# İzlenecek kritik dizinler
/etc    NORMAL
/bin    NORMAL
/sbin   NORMAL
/usr/bin NORMAL
/usr/sbin NORMAL
/boot   NORMAL
/lib    NORMAL
/lib64   NORMAL

# Hariç tut
!/var/log
!/tmp
!/run
EOF

# İlk veritabanını oluştur (referans noktası)
sudo aideinit
sudo cp /var/lib/aide/aide.db.new /var/lib/aide/aide.db

# Bütünlük kontrolü çalıştır
sudo aide --check

# Haftalık otomatik kontrol — crontab
echo "0 3 * * 0 root /usr/bin/aide --check | mail -s 'AIDE Raporu' root" | \
    sudo tee /etc/cron.d/aide-check
```

### 7.3 — Rootkit Tarama

```bash
# rkhunter — Rootkit Hunter
sudo apt install -y rkhunter chkrootkit

# rkhunter veritabanını güncelle ve tara
sudo rkhunter --update
sudo rkhunter --propupd   # Referans veritabanı oluştur
sudo rkhunter --check --rwo   # Tarama (uyarıları göster)

# chkrootkit
sudo chkrootkit

# Haftalık otomatik
echo "0 2 * * 0 root /usr/bin/rkhunter --check --rwo >> /var/log/rkhunter_weekly.log 2>&1" | \
    sudo tee /etc/cron.d/rkhunter-weekly
```

### 7.4 — Merkezi Loglama Sertleştirme

```bash
# rsyslog uzak sunucuya yedek (isteğe bağlı ama önerilir)
# Log silinse bile uzak kopya kalır

sudo tee -a /etc/rsyslog.conf <<'EOF'
# TLS ile uzak log (SIEM / log sunucusu varsa)
# *.* @@logs.ornekdomain.com:6514  # TCP TLS
*.* action(type="omfile" file="/var/log/tum_olaylar.log")
EOF

# Log dosyası izinlerini kısıtla
sudo chmod 640 /var/log/auth.log
sudo chmod 640 /var/log/syslog
sudo chown root:adm /var/log/auth.log

# Log rotasyonu
sudo tee /etc/logrotate.d/secure-logs <<'EOF'
/var/log/auth.log {
    weekly
    rotate 52
    compress
    delaycompress
    missingok
    notifempty
    create 640 root adm
}
EOF
```

---

## KALI VE UBUNTU ARASINDAKİ TEMEL FARKLAR

| Özellik | Ubuntu | Kali Linux |
|---------|--------|-----------|
| **Varsayılan Kullanıcı** | Standart | Standart (2020+) |
| **AppArmor** | ✅ Varsayılan aktif | ⚠️ Kurulu, kapalı gelebilir |
| **UFW** | ✅ Kurulu (kapalı) | ❌ Manuel kurulum |
| **Secure Boot** | ✅ shim ile destekli | ⚠️ İmzasız modüller sorun çıkarır |
| **Güncelleme Sıklığı** | LTS: 5 yıl destek | Rolling release (sürekli güncelleme) |
| **Kernel** | Stabil LTS | Kali özel (en son) |
| **Amaç** | Genel kullanım | Penetrasyon testi |
| **Sertleştirme Zorluğu** | Orta | Yüksek (araç-güvenlik çakışması) |
| **snap Paketleri** | ✅ (sandbox'lı) | ❌ Yok |
| **LUKS kurulum seçeneği** | ✅ GUI'den | ✅ Installer'dan |

---

## BÜTÜNLÜK DOĞRULAMA — Kalenin Sağlamlığını Test Et

```bash
#!/usr/bin/env bash
# linux_kale_durum.sh
# Yönetici olarak çalıştır: sudo bash linux_kale_durum.sh

GREEN='\033[0;32m'; RED='\033[0;31m'; YELLOW='\033[1;33m'; NC='\033[0m'
ok()   { echo -e "${GREEN}[✓]${NC} $1"; }
fail() { echo -e "${RED}[✗]${NC} $1"; }
warn() { echo -e "${YELLOW}[!]${NC} $1"; }

echo "========== LİNUX KALE DURUM RAPORU =========="

# LUKS Şifreleme
if cryptsetup status $(lsblk -rno NAME,TYPE | awk '$2=="crypt"{print $1}' | head -1) &>/dev/null; then
    ok "LUKS: Şifreli disk aktif"
else
    fail "LUKS: Şifreli disk bulunamadı!"
fi

# AppArmor
if systemctl is-active apparmor &>/dev/null; then
    ENFORCE=$(aa-status 2>/dev/null | grep "in enforce mode" | awk '{print $1}')
    ok "AppArmor: Aktif — $ENFORCE profil enforce modunda"
else
    fail "AppArmor: Kapalı!"
fi

# Swap Şifreleme
if grep -q crypt /proc/swaps 2>/dev/null; then
    ok "Swap: Şifreli"
else
    warn "Swap: Şifrelenmemiş veya yok"
fi

# Hibernate
if systemctl is-masked hibernate.target &>/dev/null; then
    ok "Hibernate: Devre dışı (masked)"
else
    warn "Hibernate: Aktif — Cold Boot riski"
fi

# SSH Root girişi
if sshd -T 2>/dev/null | grep -q "permitrootlogin no"; then
    ok "SSH: Root girişi devre dışı"
else
    fail "SSH: Root girişi AÇIK!"
fi

# SSH Parola kimlik doğrulama
if sshd -T 2>/dev/null | grep -q "passwordauthentication no"; then
    ok "SSH: Parola doğrulama kapalı (yalnızca anahtar)"
else
    fail "SSH: Parola doğrulama AÇIK!"
fi

# fail2ban
if systemctl is-active fail2ban &>/dev/null; then
    BANNED=$(fail2ban-client status sshd 2>/dev/null | grep "Total banned" | awk '{print $NF}')
    ok "fail2ban: Aktif — Toplam ban: $BANNED"
else
    fail "fail2ban: Kapalı!"
fi

# UFW/nftables
if ufw status 2>/dev/null | grep -q "Status: active"; then
    ok "UFW: Aktif"
elif systemctl is-active nftables &>/dev/null; then
    ok "nftables: Aktif"
else
    fail "Güvenlik Duvarı: Aktif değil!"
fi

# auditd
if systemctl is-active auditd &>/dev/null; then
    ok "auditd: Aktif"
else
    fail "auditd: Kapalı!"
fi

# Root hesabı kilidi
if passwd -S root 2>/dev/null | grep -q " L "; then
    ok "Root: Kilili"
else
    warn "Root: Kilili değil"
fi

# ASLR
ASLR=$(cat /proc/sys/kernel/randomize_va_space)
[ "$ASLR" -eq 2 ] && ok "ASLR: Tam aktif (2)" || fail "ASLR: $ASLR (2 olmalı)"

# kptr_restrict
KPTR=$(cat /proc/sys/kernel/kptr_restrict)
[ "$KPTR" -ge 1 ] && ok "kptr_restrict: $KPTR" || fail "kptr_restrict: $KPTR (1+ olmalı)"

# ptrace_scope
PTRACE=$(cat /proc/sys/kernel/yama/ptrace_scope 2>/dev/null)
[ "$PTRACE" -ge 1 ] && ok "ptrace_scope: $PTRACE" || fail "ptrace_scope: $PTRACE (1+ olmalı)"

# SYN Cookies
SYN=$(cat /proc/sys/net/ipv4/tcp_syncookies)
[ "$SYN" -eq 1 ] && ok "TCP SYN Cookies: Aktif" || fail "TCP SYN Cookies: Kapalı!"

# GRUB Şifresi
if grep -q "password_pbkdf2" /boot/grub/grub.cfg 2>/dev/null; then
    ok "GRUB: Şifreli"
else
    warn "GRUB: Şifre bulunamadı"
fi

echo "=============================================="
```

```bash
chmod +x linux_kale_durum.sh
sudo bash linux_kale_durum.sh
```

---

## ÖZET — Hangi Katman Neyi Koruyor?

| Tehdit | Katman | Tedbir |
|--------|--------|--------|
| Disk Çalındı | 1 | LUKS2 AES-256-XTS ile tam disk şifreleme |
| Cold Boot | 3 | Hibernate kapalı · Swap şifreleme · init_on_free |
| Thunderbolt DMA | 3 | IOMMU=force · VT-d/AMD IOMMU |
| Bootkit / Rootkit | 1,2 | GRUB şifresi · Secure Boot · rkhunter |
| Kernel Exploit | 2 | sysctl sertleştirme · lockdown=confidentiality |
| Mimikatz Eşdeğeri | 2 | ptrace_scope=2 · kptr_restrict=2 |
| SSH Brute-Force | 5,6 | fail2ban · MaxAuthTries · Yalnızca anahtar |
| Privilege Escalation | 4,6 | AppArmor · sudo kısıtlama · No root login |
| Uygulama Exploit | 4 | Firejail · seccomp · AppArmor profil |
| Ağ Dinleme | 5 | nftables/UFW · DNS-over-TLS · No redirect |
| Dosya Manipülasyonu | 7 | AIDE bütünlük kontrolü · auditd |
| Log Silme Saldırısı | 7 | auditd -e 2 (kurallar kilili) · Uzak log |

---

*Son güncelleme: 2026-03-22 · Ubuntu 22.04 LTS / Kali Linux 2024.x için hazırlanmıştır.*
