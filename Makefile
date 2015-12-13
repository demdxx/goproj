
install:
	go build -o /usr/local/bin/goproj cmd/goproj/main.go

uninstall:
	rm /usr/local/bin/goproj
