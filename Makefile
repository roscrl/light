#########################
##     Development     ##
#########################

DB_NAME=$(APP_NAME)
DB_FOLDER=./core/db
LOCAL_SQLITE_SCHEMA=$(DB_FOLDER)/schema.sql
LOCAL_SQLITE_DB_PATH=$(DB_FOLDER)/$(DB_NAME).db
LOCAL_SQLITE_SHM_DB_PATH=$(DB_FOLDER)/$(DB_NAME).db-shm
LOCAL_SQLITE_WAL_DB_PATH=$(DB_FOLDER)/$(DB_NAME).db-wal

lint:
	golangci-lint run --config config/.golangci.yml && \
	govulncheck ./...                               && \
	cd $(DB_FOLDER) 								&& \
	sqlc vet 										&& \
	sqlc diff

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
	go test -v -race ./...

test-browser-slow:
	go test -v -race ./... -rod=show,slow=3s,trace

run:
	go run . --config ./config/.local

run-mock:
	go run . --config ./config/.local.mock

output-schema:
	sqlite3 $(LOCAL_SQLITE_DB_PATH) .schema > $(LOCAL_SQLITE_SCHEMA)

tailwind-watch:
	./bin/tailwindcss -i ./core/views/assets/main.css -o ./core/views/assets/dist/main.css --watch --config ./config/tailwind.config.js

bench:
	go test -run=^$ -bench=. ./...

pprof:
	go tool pprof -http=:8090 bin/profile

#########################
##        Builds       ##
#########################

pre-build: lint format generate test

build-amd64-linux: pre-build
	CC="zig cc -target x86_64-linux-musl" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/app .

build-arm64-linux: pre-build
	CC="zig cc -target aarch64-linux-musl" CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o bin/app .

build: pre-build
	go build -o bin/app .

build-quick:
	go build -o bin/app .

#########################
##         VPS         ##
#########################

USER=root

APP_NAME=todos
APP_FOLDER=~/$(APP_NAME)

APP_CADDY_PATH=$(APP_NAME).caddy
SERVICE_NAME=$(APP_NAME).service

CLOUDFLARE_ZONE_ID=CHANGE_ME

ssh:
	ssh $(USER)@$(VPS_IP)

vps-new:
	ssh $(USER)@$(VPS_IP) "mkdir -p $(APP_FOLDER)"
	make vps-dependencies
	make caddy-root-config
	make caddy-cert
	make caddy-service-reload
	make caddy-reload
	make db-copy-to-prod
	make app-service-reload
	make deploy

vps-dependencies:
	ssh $(USER)@$(VPS_IP) "sudo apt-get update && sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https && curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg && curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list && sudo apt -y update && sudo apt -y install caddy && sudo apt -y install lnav"

caddy-root-config:
	scp -r ./config/Caddyfile $(USER)@$(VPS_IP):/etc/caddy/Caddyfile

caddy-cert:
	scp -r ./config/public.pem $(USER)@$(VPS_IP):/etc/ssl/certs/$(APP_NAME).pem
	ssh $(USER)@$(VPS_IP) "mkdir -p /etc/ssl/private"
	scp -r ./config/private.pem $(USER)@$(VPS_IP):/etc/ssl/private/$(APP_NAME).pem

caddy-service-reload:
	scp -r ./config/caddy.service $(USER)@$(VPS_IP):/lib/systemd/system/caddy.service
	ssh $(USER)@$(VPS_IP) "systemctl daemon-reload"
	ssh $(USER)@$(VPS_IP) "systemctl restart caddy"

caddy-reload:
	scp -r ./config/$(APP_CADDY_PATH) $(USER)@$(VPS_IP):/etc/caddy/$(APP_CADDY_PATH)
	ssh $(USER)@$(VPS_IP) "systemctl reload caddy"

db-copy-local-to-prod:
	rsync -avz --ignore-existing $(LOCAL_SQLITE_DB_PATH) $(USER)@$(VPS_IP):$(APP_FOLDER)/db/

db-copy-prod-to-local:
	rsync -avz --ignore-existing $(USER)@$(VPS_IP):$(APP_FOLDER)/db/ $(LOCAL_SQLITE_DB_PATH).prod

db-copy-local-to-prod-force:
	ssh $(USER)@$(VPS_IP) "mkdir -p $(APP_FOLDER)/db/archive"
	ssh $(USER)@$(VPS_IP) "if [ -f $(APP_FOLDER)/db/$(DB_NAME).db ];     then mv $(APP_FOLDER)/db/$(DB_NAME).db     $(APP_FOLDER)/db/archive/$(DB_NAME)_$$(date +"%Y%m%d%H%M%S").db;     fi"
	ssh $(USER)@$(VPS_IP) "if [ -f $(APP_FOLDER)/db/$(DB_NAME).db-shm ]; then mv $(APP_FOLDER)/db/$(DB_NAME).db-shm $(APP_FOLDER)/db/archive/$(DB_NAME)_$$(date +"%Y%m%d%H%M%S").db-shm; fi"
	ssh $(USER)@$(VPS_IP) "if [ -f $(APP_FOLDER)/db/$(DB_NAME).db-wal ]; then mv $(APP_FOLDER)/db/$(DB_NAME).db-wal $(APP_FOLDER)/db/archive/$(DB_NAME)_$$(date +"%Y%m%d%H%M%S").db-wal; fi"
	rsync -avz $(LOCAL_SQLITE_DB_PATH) $(USER)@$(VPS_IP):$(APP_FOLDER)/db/
	rsync -avz $(LOCAL_SQLITE_SHM_DB_PATH) $(USER)@$(VPS_IP):$(APP_FOLDER)/db/
	rsync -avz $(LOCAL_SQLITE_WAL_DB_PATH) $(USER)@$(VPS_IP):$(APP_FOLDER)/db/

app-service-reload:
	scp -r ./config/$(SERVICE_NAME) $(USER)@$(VPS_IP):/lib/systemd/system/$(SERVICE_NAME)
	ssh $(USER)@$(VPS_IP) "systemctl daemon-reload"
	ssh $(USER)@$(VPS_IP) "systemctl restart $(SERVICE_NAME)"

upload: build-amd64-linux
	ssh $(USER)@$(VPS_IP) "mkdir -p $(APP_FOLDER)/new"
	scp -r bin/app $(USER)@$(VPS_IP):$(APP_FOLDER)/new/app

deploy: upload
	ssh $(USER)@$(VPS_IP) "mkdir -p $(APP_FOLDER)/archive"
	ssh $(USER)@$(VPS_IP) "if [ -f $(APP_FOLDER)/app ]; then mv $(APP_FOLDER)/app $(APP_FOLDER)/archive/app_$$(date +"%Y%m%d%H%M%S"); fi"
	ssh $(USER)@$(VPS_IP) "mv $(APP_FOLDER)/new/app $(APP_FOLDER)/app"
	ssh $(USER)@$(VPS_IP) "systemctl restart $(SERVICE_NAME)"
	make purge-cache-prod

purge-cache-prod:
	curl -X POST https://api.cloudflare.com/client/v4/zones/$(CLOUDFLARE_ZONE_ID)/purge_cache \
		-H "X-Auth-Email: $(CLOUDFLARE_EMAIL)" \
		-H "X-Auth-Key: $(CLOUDFLARE_KEY)" \
		-H "Content-Type: application/json" \
		--data '{"purge_everything":true}'

logs-prod:
	echo "make ssh then run 'journalctl -u $(SERVICE_NAME) | lnav'"

logs-prod-tail:
	ssh $(USER)@$(VPS_IP) "journalctl -u $(SERVICE_NAME) -f"

logs-caddy-prod:
	ssh $(USER)@$(VPS_IP) "journalctl -u caddy -f"

#########################
##    Tools Install    ##
#########################

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	go install mvdan.cc/gofumpt@v0.5.0
	go install github.com/daixiang0/gci@v0.11.0
	go install golang.org/x/vuln/cmd/govulncheck@v1.0.0
	go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.19.1
	mkdir -p ./bin/
	make tools-esbuild
	make tools-tailwind
	echo "Remember to install Zig for the built-in C cross-compiler to Linux (or any C compiler for the 'make build' targets)"

tools-esbuild:
	curl -fsSL https://esbuild.github.io/dl/v0.17.17 | sh
	mv esbuild ./bin/

tools-tailwind:
	# MacOS ARM specific
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.3/tailwindcss-macos-arm64
	chmod +x tailwindcss-macos-arm64
	mv tailwindcss-macos-arm64 tailwindcss
	mv tailwindcss ./bin/
