.PHONY: all clean
.PHONY: dashboard-server dashboard-client release
.PHONY: run-server

.PHONY: docker docker-build docker-release
.PHONY: docker-run-server

.PHONY: nginx nginx-run

all: dashboard-server dashboard-client

GOOS=
GOARCH=

RELEASE_SERVER_OSARCH=\
	linux/amd64 \
	linux/386 \
	linux/arm \
	windows/386 \
	windows/amd64 \
	darwin/amd64

RELEASE_CLIENT_OSARCH=\
	linux/amd64 \
	linux/386 \
	linux/arm

dashboard-server:
	npm install
	gulp build
	bower install --allow-root
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o dashboard-server server/main/*.go
	rm -f assets; zip -r assets dist
	cat assets.zip >> dashboard-server
	rm assets.zip

dashboard-client:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o dashboard-client client/main/*.go

run-server:
	HTML_ROOT=app go run server/main/*.go -c server/main/dev_config.json

docker:
	docker build -t status-dashboard .

docker-run-server:
	docker run \
		--rm \
		-it \
		-v "${GOPATH}":/go \
		-w /go/src/github.com/hverr/status-dashboard \
		--name dashboard \
		status-dashboard \
		make run-server

docker-build:
	docker run \
		--rm \
		-v "${GOPATH}":/go \
		-w /go/src/github.com/hverr/status-dashboard \
		status-dashboard \
		make GOOS=$(GOOS) GOARCH=$(GOARCH)

docker-release:
	docker run \
		--rm \
		-v "${GOPATH}/src/github.com/hverr/status-dashboard":/go/src/github.com/hverr/status-dashboard \
		-w /go/src/github.com/hverr/status-dashboard \
		status-dashboard \
		make release

release:
	go get -v ./...
	go get -v github.com/mattn/go-isatty
	go get -v github.com/mitchellh/gox
	rm -rf release && mkdir release
	#
	npm install
	bower install --allow-root
	gulp build
	${GOPATH}/bin/gox \
		-osarch="$(RELEASE_SERVER_OSARCH)" \
		-output="release/dashboard-server_{{.OS}}_{{.Arch}}" \
		github.com/hverr/status-dashboard/server/main
	rm -f assets.zip; zip -r assets dist
	cat assets.zip | tee -a release/dashboard-server* >/dev/null
	rm assets.zip
	#
	${GOPATH}/bin/gox \
		-osarch="$(RELEASE_CLIENT_OSARCH)" \
		-output="release/dashboard-client_{{.OS}}_{{.Arch}}" \
		github.com/hverr/status-dashboard/client/main

nginx:
	docker build -t status-dashboard-nginx -f Dockerfile.nginx .

nginx-run:
	docker run \
		--rm \
		-it \
		--link dashboard:dashboard \
		-p 12443:443 \
		--name dashboard-nginx \
		status-dashboard-nginx \
		nginx -g 'daemon off;'

clean:
	rm -f dashboard-server
	rm -f dashboard-client
