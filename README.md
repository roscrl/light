# light

fullstack go web app template, for those who enjoy [radical simplicty](https://www.radicalsimpli.city).

## Setup


<p align="center">
  <img width="681" alt="image" src="https://github.com/roscrl/light/assets/13072760/58030551-0b2f-43e3-898b-b3d388b4b85f">
</p>

```bash
git clone https://github.com/roscrl/light.git &&
cd light                                      &&
make tools                                    &&
make run
```

## Go Template

```bash
go install golang.org/x/tools/cmd/gonew@latest
```

```bash
gonew github.com/roscrl/light
```

## Dependencies

#### Backend

`go-sqlite3` database driver, requires CGO to build. Prefer `zig cc` over `gcc`/`clang` for easier cross compilation


#### Frontend

`tailwindcss` styling

`@hotwired/turbo` SPA like navigation

`@hotwired/stimulus` lightweight JS functionality

`esbuild` bundling

#### Development

`sqlc` generate Go code from [SQL queries](db/query.sql)

`is` test assertions

`rod` browser testing

`fsnotify` watch Go template changes in dev mode without recompiling

`golangci-lint` linting

`gofumpt` formatting

## VPS Deploy Checklist

### Cloudflare

- Set SSL `Full (strict)`
- Add an A record in the DNS settings pointing to VPS IP
- Create a 15 year 'Origin Certificate' and place in `config/public.pem` & `config/private.pem`
- Enable Rate Limiting
  - `(http.request.uri.path contains "/")` 50 requests per 10s
- Enable [Bot Fight Mode](https://developers.cloudflare.com/bots/get-started/free/)
- Enable Page Rules Caching to respect `Cache-Control` headers returned
  - playlistvote.com/* Cache Level: Cache Everything
- Always Use HTTPS, Enable Brotli

### Hetzner

- Set firewall to allow only [Cloudflare IPs](https://www.cloudflare.com/en-gb/ips/) on port 443
- Set firewall to allow only personal IP on port 22

### VPS Setup

- Change `Makefile` constant `APP_NAME` to your own
- Change `Makefile` constant `CLOUDFLARE_ZONE_ID` to your own


- Change filename `config/todos.caddy` to `<APP_NAME>.caddy`
- Change `config/todos.caddy` `tls /etc/ssl/certs/todos.pem /etc/ssl/private/todos.pem` to `tls /etc/ssl/certs/<APP_NAME>.pem /etc/ssl/private/<APP_NAME>.pem`


- Change filename `config/todos.service` to `<APP_NAME>.service`
- Change `config/todos.service` `EnvironmentFile=/root/todos/.prod` to `EnvironmentFile=/root/<APP_NAME>/.prod`
- Change `config/todos.service` `WorkingDirectory=/root/todos` to `WorkingDirectory=/root/<APP_NAME>`
- Change `config/todos.service` `ExecStart=/root/todos/app` to `ExecStart=/root/<APP_NAME>/app`


- Create `config/.prod` using `config/.prod.template` as a template


- Ensure `config/private.pem` exists (cloudflare origin certificate private key from cloudflare setup)
- Set `VPS_IP` environment variable to your VPS IP
- Set `CLOUDFLARE_ZONE_ID` environment variable to your cloudflare zone id
- Set `CLOUDFLARE_EMAIL` environment variable to your cloudflare email
- Set `CLOUDFLARE_KEY` environment variable to your cloudflare key
- Run `make vps-new`


### The classic abandoned TODO section

- On PR, create a ephemeral preview environment with usage of Makefile VPS creation commands
- On PR, performance testing should execute and attach a graph with ability to see performance changes overtime. `rod` for e2e? but `tests/bench` for individual endpoints
- Code coverage + quality + unit test reports integrated to PR
- `core/jobs` package for background jobs
- `core/notify` package for sending notifications (email, sms, push, etc)