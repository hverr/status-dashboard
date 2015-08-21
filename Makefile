.PHONY: all
all: dashboard-server dashboard-client

dashboard-server:
	gulp build
	go build -o dashboard-server server/main/*.go
	rm -f assets; zip -r assets dist
	cat assets.zip >> dashboard-server
	rm assets.zip

dashboard-client:
	go build -o dashboard-client client/main/*.go

.PHONY: clean
clean:
	rm -f dashboard-server
	rm -f dashboard-client
