kws :=KWS_Core

.PHONY: build run clean conf test test-v

conf:	
	./build/go.sh
	./build/download.sh
	./build/autoconfig.sh
	./build/apparmor.sh
	build

build: main.go
	go build -o $(kws) .

run:	build
	chmod +x ./$(kws) 
	./$(kws)



test:
	go test ./...

test-v:
	go test -v ./...

clean:
	rm -rf $(kws)
