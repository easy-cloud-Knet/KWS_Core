kws :=KWS_Core
NETWORK:=ovn


.PHONY: build run clean conf test test-v

conf:	
	./build/go.sh
	./build/download.sh
	./build/autoconfig.sh
	./build/apparmor.sh
	build

build: main.go
	go build -ldflags "-X 'github.com/easy-cloud-Knet/KWS_Core/vm/parsor.NetworkMode=$(NETWORK)'" -o $(kws) .

run:	build
	chmod +x ./$(kws) 
	./$(kws)



test:
	go test ./...

test-v:
	go test -v ./...

clean:
	rm -rf $(kws)
