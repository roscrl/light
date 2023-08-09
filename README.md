# light

fullstack go web app template, for those without [imaginary problems](https://cerebralab.com/Imaginary_Problems_Are_the_Root_of_Bad_Software).

## Setup

```bash
go install golang.org/x/tools/cmd/gonew@latest
```

```bash
gonew github.com/roscrl/light
```

## Dependencies

`tailwindcss` styling

`@hotwired/turbo` SPA like navigation

`@hotwired/stimulus` lightweight JS functionality

`go-sqlite3` database driver, requires CGO to build. Prefer `zig cc` over `gcc` for easier cross compilation

`sqlc` generate Go code from [SQL queries](core/db/query.sql)

`is` test assertions

`rod` browser testing

`fsnotify` watch Go template changes in dev mode without recompiling

## Deploy

Prefer to deploy on a [VPS](https://specbranch.com/posts/one-big-server/)

#### VPS Setup

- Ensure `config/private.pem` exists (cloudflare origin certificate private key)
- Ensure `config/.prod` exists (app config)
- Set `VPS_IP` environment variable
- Set `CLOUDFLARE_EMAIL` environment variable
- Set `CLOUDFLARE_KEY` environment variable
- Run `make vps-new`

### Cloudflare

- Set SSL `Full (strict)`
- Add an A record in the DNS settings pointing to VPS IP
- Create Origin Certificate and place in `config/public.pem` & `config/private.pem`
- Enable Rate Limiting
  - `(http.request.uri.path contains "/")` 50 requests per 10s
- Enable [Bot Fight Mode](https://developers.cloudflare.com/bots/get-started/free/)
- Enable Page Rules Caching to respect `Cache-Control` headers returned
    - playlistvote.com/* Cache Level: Cache Everything
- Always Use HTTPS, Enable Brotli

### Hetzner

- Set firewall to allow only [Cloudflare IPs](https://www.cloudflare.com/en-gb/ips/) on port 443
- Set firewall to allow only personal IP on port 22