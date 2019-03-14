bin:
	mkdir -p bin

hello: bin
	cd cmd/hello && go build -o ../../bin/hello

greet: bin
	cd cmd/greet && go build -o ../../bin/greet

greet-windows: bin
	cd cmd/greet && GOOS=windows GOARCH=386 go build -o ../../bin/greet.exe

ping_valet: bin
	cd cmd/ping_valet && go build -o ../../bin/ping_valet

cats: bin 
	cd cmd/cats && go build -o ../../bin/cats

dogs: bin
	cd cmd/dogs && go build -o ../../bin/dogs

weather: bin
	cd cmd/weather && go build -o ../../bin/weather