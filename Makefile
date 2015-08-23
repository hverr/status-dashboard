.PHONY: all clean
.PHONY: dashboard-server dashboard-client
.PHONY: run-server

.PHONY: docker docker-build docker-release
.PHONY: docker-run-server
.PHONY: docker-release-component

.PHONY: nginx nginx-run

all: dashboard-server dashboard-client

GOOS=
GOARCH=

dashboard-server:
	npm install
	gulp build
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
	go get -v github.com/mattn/go-isatty
	rm -rf release && mkdir release
	make docker-release-component GOOS=linux GOARCH=386
	make docker-release-component GOOS=linux GOARCH=amd64
	make docker-release-component GOOS=darwin GOARCH=386
	make docker-release-component GOOS=darwin GOARCH=amd64
	make docker-release-component GOOS=windows GOARCH=386
	make docker-release-component GOOS=windows GOARCH=amd64

docker-release-component:
	make docker-build GOOS=$(GOOS) GOARCH=$(GOARCH)
	cp dashboard-server release/dashboard-server-$(GOOS)-$(GOARCH)
	cp dashboard-client release/dashboard-client-$(GOOS)-$(GOARCH)

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
