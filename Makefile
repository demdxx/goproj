
install:
	go build -o /usr/local/bin/goproj cmd/main.go

uninstall:
	rm /usr/local/bin/goproj
