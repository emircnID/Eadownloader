# EaDownloader

Telegram icin hafif bir sosyal medya indirme botu.

Desteklenen platformlar:

- Facebook
- Instagram
- TikTok
- Twitter (X)

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
git clone https://github.com/KULLANICI_ADIN/EaDownloader.git
cd EaDownloader
cp .env.example .env
```

`.env` icinde en az sunlari degistir:

```env
BOT_TOKEN=BotFather_tokenin
DB_PASSWORD=guclu-bir-sifre
ADMINS=telegram_user_id
```

Sonra baslat:

```bash
docker compose up -d --build
docker compose logs -f bot
```

Guncelleme icin:

```bash
git pull
docker compose up -d --build
```

## Cookie Dosyalari

Twitter/X ve Facebook bazen giris cookie'si olmadan medya vermez.

Cookie dosyalarini Netscape formatinda su yollara koy:

```text
private/cookies/twitter.txt
private/cookies/facebook.txt
```

## Notlar

- Bot polling ile calisir; sunucuda webhook veya domain zorunlu degil.
- `METRICS_PORT` ve `PROFILER_PORT` Docker tarafinda sadece `127.0.0.1` uzerinden acilir.
- `MAX_FILE_SIZE` MB cinsindendir.
