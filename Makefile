kws     := KWS_Core
NETWORK := ovn
CTL_IP  := 10.5.15.39
IP      :=
DNS     :=

.PHONY: build run clean conf test test-v ovn-install ovn-cluster ovn-worker

conf:
	./build/go.sh
	./build/download.sh
	./build/autoconfig.sh
ifeq ($(NETWORK),ovn)
	./build/apparmor.sh
endif
	$(MAKE) build

build: main.go
	go build -ldflags "-X 'github.com/easy-cloud-Knet/KWS_Core/vm/parsor.NetworkMode=$(NETWORK)'" -o $(kws) .

run: build
	chmod +x ./$(kws)
	./$(kws)

# OVN/OVS 소스 빌드 설치 (최초 1회)
ovn-install:
	./build/sdnConf.sh

# 컨트롤 노드: make ovn-cluster IP=10.5.15.39 DNS=10.5.15.1
ovn-cluster:
	@test -n "$(IP)" || { echo "usage: make ovn-cluster IP=x.x.x.x DNS=x.x.x.x"; exit 1; }
	./build/ovn-cluster.sh $(IP) $(DNS)

# 워커 노드: make ovn-worker IP=10.5.15.40 DNS=10.5.15.1 [CTL_IP=x.x.x.x]
ovn-worker:
	@test -n "$(IP)" || { echo "usage: make ovn-worker IP=x.x.x.x DNS=x.x.x.x [CTL_IP=x.x.x.x]"; exit 1; }
	./build/ovn-worker.sh $(IP) $(DNS) $(CTL_IP)

test:
	go test ./...

test-v:
	go test -v ./...

clean:
	rm -rf $(kws)
