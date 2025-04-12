kws :=KWS_Core

.PHONY: build run clean conf

conf:	
	./build/go.sh
	./build/download.sh
	./build/autoconfig.sh
	build

build: main.go
	go build -o $(kws) .

run:	build
	chmod +x ./$(kws) 
	./$(kws)



clean: 
	rm -rf $(kws)
