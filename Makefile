
build:
	@echo "Building Paxos Server"	
	go build -o GoPaxos main.go

test:
	@echo "Testing GoPaxos"	
	go test -v --cover ./...

all:
	@echo "Build From Scratch"
	make build && make provision

requests:
	@echo "Run client/client.go"
	/bin/bash -c "go run ${PWD}/client/ -port=8000& go run ${PWD}/client/ -port=8001& go run ${PWD}/client/ -port=8002"
	
logs:
	@echo "Run mytest/correctness.go"
	go run ${PWD}/mytest/

cmp:
	@echo "compare all logs"
	@echo "---------------------  logs ---------------------------"
	@tree ${HOME}/paxos_logs
	@echo "-------------- node 0  && node 1------------------------"
	@diff ${HOME}/paxos_logs/n0 ${HOME}/paxos_logs/n1
	@echo "-------------- node 1  && node 2------------------------"
	@diff ${HOME}/paxos_logs/n1 ${HOME}/paxos_logs/n2

diff:
	@echo "use diy comparator"
	@go run ${PWD}/scripts

crash:
	@echo "test loop crash"
	@bash ${PWD}/scripts/loopCrash.sh

provision:
	@echo "Provisioning Paxos Cluster"	
	bash scripts/provision.sh

docker:
	@echo "Building Paxos Docker Image"	
	docker build -t paxos -f Dockerfile .

paxos-run:
	@echo "Running Single Paxos Docker Container"
	docker run -p 8080:8080 -d paxos

info:
	echo "Paxos Cluster Nodes"
	docker ps | grep 'paxos'
	docker network ls | grep paxos_network

stop:
	@echo "Stop Paxos Cluster"
	@docker ps -a | awk '$$2 ~ /paxos/ {print $$1}' | xargs -I {} docker stop  {} 

clean: 
	@echo "Cleaning Paxos Cluster"
	@docker ps -a | awk '$$2 ~ /paxos/ {print $$1}' | xargs -I {} docker rm -f {} 
	@docker network ls | grep paxos_network | xargs docker network rm paxos_network 2>/dev/null

	