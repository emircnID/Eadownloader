# EaDownloader

Telegram icin hafif bir sosyal medya indirme botu.

Desteklenen platformlar:

- Facebook
- Instagram
- TikTok
- Twitter (X)
- YouTube

Bot public calisacak sekilde hazirlandi. Whitelist varsayilan olarak kapali, diller sadece English ve Turkce.

## Sunucuda Calistirma

Oracle Free Tier Ubuntu uzerinde en kolay yol Docker Compose ile calistirmaktir.

```bash
sudo apt update
sudo apt install -y git docker.io docker-compose-plugin
sudo usermod -aG docker "$USER"
newgrp docker
```

Repo'yu sunucuya cek:

```bash
git clone https://github.com/emircnID/EaDownloader.git
cd EaDownloader
nano .env
```

`.env` icinde en az sunlari degistir:

```env
BOT_TOKEN=BotFather_tokenin
DB_PASSWORD=guclu-bir-sifre
ADMINS=telegram_user_id
```

Sonra baslat:

```bash
docker compose pull
docker compose up -d
docker compose logs -f bot
```

Guncelleme icin:

```bash
git pull
docker compose pull
docker compose up -d
```

## Cookie Dosyalari

Twitter/X, Facebook ve YouTube public iceriklerde cookiesiz denenir. Private,
yas kisitli veya giris isteyen icerikler yine cookie gerektirebilir.

Cookie dosyalarini Netscape formatinda su yollara koy:

```text
private/cookies/twitter.txt
private/cookies/facebook.txt
private/cookies/youtube.txt
```

## Notlar

- Bot polling ile calisir; sunucuda webhook veya domain zorunlu degil.
- `METRICS_PORT` ve `PROFILER_PORT` Docker tarafinda sadece `127.0.0.1` uzerinden acilir.
- `MAX_FILE_SIZE` MB cinsindendir.
