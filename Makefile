kws :=KWS_Server

.PHONY: build run clean

build: main.go
	go build -o $(kws) .

run:	build
	chmod +x ./$(kws) 
	./$(kws)



clean: 
	rm -rf $(kws)
