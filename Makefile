kws :=KWS_Server

build:
	go build  -o $(kws) .

run: build
	chmod +x ./$(kws) 
	./$(kws)



clean: 
	rm -rf $(kws)
