# light

fullstack go web app template, for those who enjoy [radical simplicty](https://www.radicalsimpli.city).

## Setup

```bash
go install golang.org/x/tools/cmd/gonew@latest
```

```bash
gonew github.com/roscrl/light
```
<img width="681" alt="image" src="https://github.com/roscrl/light/assets/13072760/58030551-0b2f-43e3-898b-b3d388b4b85f">



## Dependencies

#### Backend

`go-sqlite3` database driver, requires CGO to build. Prefer `zig cc` over `gcc`/`clang` for easier cross compilation


#### Frontend

`tailwindcss` styling

`@hotwired/turbo` SPA like navigation

`@hotwired/stimulus` lightweight JS functionality

`esbuild` bundling

#### Development

`sqlc` generate Go code from [SQL queries](core/db/query.sql)

`is` test assertions

`rod` browser testing

`fsnotify` watch Go template changes in dev mode without recompiling

`golangci-lint` linting

`gofumpt` formatting

## Deploy

Prefer to deploy on a [VPS](https://specbranch.com/posts/one-big-server/)
- Search for `CHANGE_ME` in the codebase and replace with your own values

#### VPS Setup

- Ensure `config/private.pem` exists (cloudflare origin certificate private key)
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
