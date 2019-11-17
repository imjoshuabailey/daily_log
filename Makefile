bin:
	mkdir -p bin

daily_log: bin
	cd cmd/daily_log && go build -o ../../bin/daily_log

daily_log-windows: bin
	cd cmd/daily_log && GOOS=windows GOARCH=386 go build -o ../../bin/daily_log.exe

migrate: bin
	cd cmd/migrate && go build -o ../../bin/migrate