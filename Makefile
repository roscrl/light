#########################
##      Constants      ##
#########################

APP_NAME=todos
DB_NAME=app
DB_FOLDER=./db

BUILD_TAGS=fts5
SQLITE_PATH_SCHEMA=$(DB_FOLDER)/schema.sql
SQLITE_PATH_DB=$(DB_FOLDER)/$(DB_NAME).db
SQLITE_PATH_DB_SHM=$(DB_FOLDER)/$(DB_NAME).db-shm
SQLITE_PATH_DB_WAL=$(DB_FOLDER)/$(DB_NAME).db-wal

VPS_USER=root
VPS_APP_SERVICE_NAME=$(APP_NAME).service
VPS_APP_FOLDER=~/$(APP_NAME)
VPS_APP_CADDY_PATH=$(APP_NAME).caddy

CLOUDFLARE_ZONE_ID=CHANGE_ME

#########################
##     Development     ##
#########################

lint:
	golangci-lint run --config config/.golangci.yml && \
	govulncheck ./...                               && \
	cd $(DB_FOLDER) 								&& \
	sqlc diff
#	sqlc vet # this is panicking, might be a bug in sqlc with SQLite

format:
	gofumpt -extra -l -w . 										&& \
	gci write -s standard -s default -s "prefix(github.com/)" .

sqlc:
	cd $(DB_FOLDER) && sqlc generate

generate: sqlc
	./bin/tailwindcss -i ./core/views/assets/main.css -o ./core/views/assets/dist/main.css --config ./config/tailwind.config.js
	./bin/esbuild core/views/assets/dist/js/vendor/stimulus-3.2.1/stimulus.js           --minify --outfile=core/views/assets/dist/js/vendor/stimulus-3.2.1/stimulus.min.js
	./bin/esbuild core/views/assets/dist/js/vendor/turbo-7.3.0/dist/turbo.es2017-esm.js --minify --outfile=core/views/assets/dist/js/vendor/turbo-7.3.0/dist/turbo.es2017-esm.min.js

test:
	go test --tags "$(BUILD_TAGS)" -v -race -cover ./... -coverprofile=bin/test-coverage.out
	go tool cover -html=bin/test-coverage.out -o bin/test-coverage-report.html

test-browser-slow:
	go test --tags "$(BUILD_TAGS)" -v -race ./tests/browser -rod=show,trace,slow=0.4s

run:
	go run --tags "$(BUILD_TAGS)" . --config ./config/.local & make tailwind-watch

run-mock:
	go run --tags "$(BUILD_TAGS)" . --config ./config/.local.mock & make tailwind-watch

output-schema:
	sqlite3 $(SQLITE_PATH_DB) .schema > $(SQLITE_PATH_SCHEMA)

vacuum:
	sqlite3 $(SQLITE_PATH_DB) "VACUUM;"

tailwind-watch:
	./bin/tailwindcss -i ./core/views/assets/main.css -o ./core/views/assets/dist/main.css --watch --config ./config/tailwind.config.js

bench:
	go test --tags "$(BUILD_TAGS)" -run=XXX -bench=. ./... | tee bin/bench.txt

bench-individual-cpu:
	go test --tags "$(BUILD_TAGS)" -run=XXX -bench=BenchmarkHome -cpuprofile ./bin/BenchmarkHome.out ./tests/bench

pprof:
	go tool pprof -http=:8090 bin/BenchmarkHome.out

#########################
##        Builds       ##
#########################

pre-build: lint format generate test

build-amd64-linux: pre-build
	CC="zig cc -target x86_64-linux-musl" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --tags "$(BUILD_TAGS)" -o bin/app .

build-arm64-linux: pre-build
	CC="zig cc -target aarch64-linux-musl" CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build --tags "$(BUILD_TAGS)" -o bin/app .

build: pre-build
	go build --tags "$(BUILD_TAGS)" -o bin/app .

build-quick:
	go build --tags "$(BUILD_TAGS)" -o bin/app .

#########################
##         VPS         ##
#########################

ssh:
	ssh $(VPS_USER)@$(VPS_IP)

vps-new:
	ssh $(VPS_USER)@$(VPS_IP) "mkdir -p $(VPS_APP_FOLDER)"
	make vps-dependencies
	make caddy-root-config
	make caddy-cert
	make caddy-service-reload
	make caddy-reload
	make db-copy-local-to-prod
	make app-service-reload
	make deploy

vps-dependencies:
	ssh $(VPS_USER)@$(VPS_IP) "sudo apt-get update && sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https && curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg && curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list && sudo apt -y update && sudo apt -y install caddy lnav unattended-upgrades apt-listchanges"

caddy-root-config:
	scp -r ./config/Caddyfile $(VPS_USER)@$(VPS_IP):/etc/caddy/Caddyfile

caddy-cert:
	scp -r ./config/public.pem $(VPS_USER)@$(VPS_IP):/etc/ssl/certs/$(APP_NAME).pem
	ssh $(VPS_USER)@$(VPS_IP) "mkdir -p /etc/ssl/private"
	scp -r ./config/private.pem $(VPS_USER)@$(VPS_IP):/etc/ssl/private/$(APP_NAME).pem

caddy-service-reload:
	scp -r ./config/caddy.service $(VPS_USER)@$(VPS_IP):/lib/systemd/system/caddy.service
	ssh $(VPS_USER)@$(VPS_IP) "systemctl daemon-reload"
	ssh $(VPS_USER)@$(VPS_IP) "systemctl restart caddy"

caddy-reload:
	scp -r ./config/$(VPS_APP_CADDY_PATH) $(VPS_USER)@$(VPS_IP):/etc/caddy/$(VPS_APP_CADDY_PATH)
	ssh $(VPS_USER)@$(VPS_IP) "systemctl reload caddy"

db-copy-local-to-prod:
	rsync -avz --ignore-existing $(SQLITE_PATH_DB) $(VPS_USER)@$(VPS_IP):$(VPS_APP_FOLDER)/db/

db-copy-prod-to-local:
	rsync -avz --ignore-existing $(VPS_USER)@$(VPS_IP):$(VPS_APP_FOLDER)/db/ $(SQLITE_PATH_DB).prod

db-copy-local-to-prod-force:
	ssh $(VPS_USER)@$(VPS_IP) "mkdir -p $(VPS_APP_FOLDER)/db/archive"
	ssh $(VPS_USER)@$(VPS_IP) "if [ -f $(VPS_APP_FOLDER)/db/$(DB_NAME).db ];     then mv $(VPS_APP_FOLDER)/db/$(DB_NAME).db     $(VPS_APP_FOLDER)/db/archive/$(DB_NAME)_$$(date +"%Y%m%d%H%M%S").db;     fi"
	ssh $(VPS_USER)@$(VPS_IP) "if [ -f $(VPS_APP_FOLDER)/db/$(DB_NAME).db-shm ]; then mv $(VPS_APP_FOLDER)/db/$(DB_NAME).db-shm $(VPS_APP_FOLDER)/db/archive/$(DB_NAME)_$$(date +"%Y%m%d%H%M%S").db-shm; fi"
	ssh $(VPS_USER)@$(VPS_IP) "if [ -f $(VPS_APP_FOLDER)/db/$(DB_NAME).db-wal ]; then mv $(VPS_APP_FOLDER)/db/$(DB_NAME).db-wal $(VPS_APP_FOLDER)/db/archive/$(DB_NAME)_$$(date +"%Y%m%d%H%M%S").db-wal; fi"
	rsync -avz $(SQLITE_PATH_DB) $(VPS_USER)@$(VPS_IP):$(VPS_APP_FOLDER)/db/
	rsync -avz $(SQLITE_PATH_DB_SHM) $(VPS_USER)@$(VPS_IP):$(VPS_APP_FOLDER)/db/
	rsync -avz $(SQLITE_PATH_DB_WAL) $(VPS_USER)@$(VPS_IP):$(VPS_APP_FOLDER)/db/

app-service-reload:
	scp -r ./config/.prod $(VPS_USER)@$(VPS_IP):$(VPS_APP_FOLDER)
	scp -r ./config/$(VPS_APP_SERVICE_NAME) $(VPS_USER)@$(VPS_IP):/lib/systemd/system/$(VPS_APP_SERVICE_NAME)
	ssh $(VPS_USER)@$(VPS_IP) "systemctl daemon-reload"
	ssh $(VPS_USER)@$(VPS_IP) "systemctl restart $(VPS_APP_SERVICE_NAME)"

upload: build-amd64-linux
	scp -r ./config/.prod $(VPS_USER)@$(VPS_IP):$(VPS_APP_FOLDER)
	ssh $(VPS_USER)@$(VPS_IP) "mkdir -p $(VPS_APP_FOLDER)/new"
	scp -r bin/app $(VPS_USER)@$(VPS_IP):$(VPS_APP_FOLDER)/new/app

deploy: upload
	ssh $(VPS_USER)@$(VPS_IP) "mkdir -p $(VPS_APP_FOLDER)/archive"
	ssh $(VPS_USER)@$(VPS_IP) "if [ -f $(VPS_APP_FOLDER)/app ]; then mv $(VPS_APP_FOLDER)/app $(VPS_APP_FOLDER)/archive/app_$$(date +"%Y%m%d%H%M%S"); fi"
	ssh $(VPS_USER)@$(VPS_IP) "mv $(VPS_APP_FOLDER)/new/app $(VPS_APP_FOLDER)/app"
	ssh $(VPS_USER)@$(VPS_IP) "systemctl restart $(VPS_APP_SERVICE_NAME)"
	make purge-cache-prod

purge-cache-prod:
	curl -X POST https://api.cloudflare.com/client/v4/zones/$(CLOUDFLARE_ZONE_ID)/purge_cache \
		-H "X-Auth-Email: $(CLOUDFLARE_EMAIL)" \
		-H "X-Auth-Key: $(CLOUDFLARE_KEY)" \
		-H "Content-Type: application/json" \
		--data '{"purge_everything":true}'

logs-prod:
	echo "make ssh then run 'journalctl -u $(VPS_APP_SERVICE_NAME) | lnav'"

logs-prod-tail:
	ssh $(VPS_USER)@$(VPS_IP) "journalctl -u $(VPS_APP_SERVICE_NAME) -f"

logs-caddy-prod:
	ssh $(VPS_USER)@$(VPS_IP) "journalctl -u caddy -f"

#########################
##    Tools Install    ##
#########################

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	go install mvdan.cc/gofumpt@v0.5.0
	go install github.com/daixiang0/gci@v0.11.0
	go install golang.org/x/vuln/cmd/govulncheck@v1.0.1
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.20.0
	mkdir -p ./bin/
	make tools-esbuild
	make tools-tailwind
	echo "Remember to install Zig for the built-in C cross-compiler to Linux (or any C compiler for the 'make build' targets)"

tools-esbuild:
	curl -fsSL https://esbuild.github.io/dl/v0.17.17 | sh
	mv esbuild ./bin/

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

tools-tailwind:
ifeq ($(UNAME_S), Darwin)
ifeq ($(UNAME_M), arm64)
	# MacOS ARM specific
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.3/tailwindcss-macos-arm64
	chmod +x tailwindcss-macos-arm64
	mv tailwindcss-macos-arm64 tailwindcss
	mv tailwindcss ./bin/
else
	# MacOS x64 specific (assuming this URL exists, update accordingly)
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.3/tailwindcss-macos-x64
	chmod +x tailwindcss-macos-x64
	mv tailwindcss-macos-x64 tailwindcss
	mv tailwindcss ./bin/
endif
else ifeq ($(OS), Windows_NT)
	# Windows specific
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.3/tailwindcss-windows-x64.exe
	mv tailwindcss-windows-x64.exe tailwindcss.exe
	mv tailwindcss.exe ./bin/
else
	# Assuming Linux x64 by default, update with the correct URL
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.3/tailwindcss-linux-x64
	chmod +x tailwindcss-linux-x64
	mv tailwindcss-linux-x64 tailwindcss
	mv tailwindcss ./bin/
endif