.PHONY: all clean
.PHONY: dashboard-server dashboard-client
.PHONY: docker docker-build docker-release
.PHONY: docker-release-component

all: dashboard-server dashboard-client

GOOS=
GOARCH=

DOCKER_BUILD_PARAMS=\
	--rm\
	-v "${PWD}":/go/src/github.com/hverr/status-dashboard \
	-w /go/src/github.com/hverr/status-dashboard \
	status-dashboard

dashboard-server:
	npm install
	gulp build
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o dashboard-server server/main/*.go
	rm -f assets; zip -r assets dist
	cat assets.zip >> dashboard-server
	rm assets.zip

dashboard-client:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o dashboard-client client/main/*.go

docker:
	docker build -t status-dashboard .

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

clean:
	rm -f dashboard-server
	rm -f dashboard-client
