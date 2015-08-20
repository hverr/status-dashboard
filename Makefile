.PHONY: all
all: dashboard-server

dashboard-server:
	gulp build
	go build -o dashboard-server server/main/*.go
	rm -f assets; zip -r assets dist
	cat assets.zip >> dashboard-server
	rm assets.zip

.PHONY: clean
clean:
	rm -f dashboard-server
	rm -f dashboard-client
