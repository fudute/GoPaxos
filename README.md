# GoPaxos

Paxos Made Simple, Implemented On Docker Containers

## Introduction

Paxos is a consensus algorithm used to establish consensus among several nodes in a distributed system. Here, we use Docker to containerize a Paxos cluster node, and Golang to implement the Paxos made simple protocol among every node. 
Using Paxos the cluster forms a very simple distributed Key-Value Store enabling the user to write and read data across any node in the cluster.

## Steps
### Requirenments
This program requires you to install:
> docker golang (necessary)
> lsof tree (for convenience)

### Start the Cluster
After cloning the repo. To provision the cluster:

```
$ make all
```    

This cleans all old bin files, containers, logs. And build everything from scatch and then run it in docker.

This creates a 3 node Paxos cluster established in their own docker network called `paxos_network` by default, and maps the container's `/home/logs` to your host path `${HOME}/paxos_logs/`, you can find logs inside it.

## Convenient Instructions
Here are some other instructions may be helpful to you:

`make docker` only build docker images

`make build` only build GoPaxos bin file (the paxos node)

`make provison` build everything but GoPaxos bin file from scratch 

`make info` inspect the cluster

`make request` After you run the cluster, this can let you send some random request to the cluster

`make logs` This let the cluster print out the logs

`make diff` This checks out the consistency of all the logs

`make cmp` simply use diff to check the consistency, may report unexisted errors

`make crash` let you could shutdown some nodes periodly

`make stop` stop the cluster from running

## Test the Cluster
The recommanded way to to test the cluster is:
```
$ make all
# and wait 3 seconds

$ make request
$ make logs
$ make diff # check logs
```

## Client Request Interfaces
The deafult rest api for client to send request is as follows:
```
url                     Method         Params
/store/set              Post           {"key":"key1","value":"val1"} //json
/store/get/key1         Get            
/store/nop              Get
/crash                  Get
/log/print/logFileName  Get
```
The default port starts from 8000, we may require root priviledge to kill the prosess that occupied the port.

Now we can send requests to Set and Get Key-Values to any peer node using its port allocated.

```
$ curl -i localhost:<peer-port>/store/set/<key>/<value>
$ curl -i localhost:<peer-port>/store/get/<key>
```

## Crash And Recover
The crashed node would restart automatically, and restore the state mechine using the stored logs in the internal level-db. Then it re-connect to its peers and require them to re-connect to itself. Finnally it starts to listen and serve the outworld client.

The proposal of nop operation could help it catch up others. You can send it by visit `/nop ` or start a loop using goroutine inside the node to achive this.


This is not certain to clean up all the locally created docker images at times. You can do a docker rmi to delete them.

## BenchMark
just run
```
go test -bench=./mytest
```


##  Paxos

The Paxos consensus algorithm is implemented using Golang running as a Paxos server in each node. Paxos consists of 3 phases:

- **Prepare Phase**:  This is the start of the Paxos phase enabled when a client would like to write data to the cluster. Here, the Prepare process generates a round ID of its own and propagates it to all the nodes in the cluster. Once a majority of nodes accept the prepare message it then moves to the accept phase.


- **Accept Phase**:  Here the same leader node that transmitted the prepare message sends an accept request to all the nodes again to accept the given value to be chosen, thus achieving consensus among a majority of nodes and at times all the nodes.


- **Learn Phase**: Once the above two phases are complete the leader then sends a learn request which enables all the nodes to persist the agreed-upon value to its store.

## Docker

Docker enables each Paxos node to be isolated and can run anywhere. Docker network here establishes a network across all the nodes so that each node can communicate with each other and ingress/egress with the host machine.

## Simple Paxos vs Multi-Paxos

The current implementation of Paxos here is Paxos Made Simple protocol, which in a real-world production environment would fare much better. Future improvements to GoPaxos would look at upgrading the protocol to Multi-Paxos. Multi-Paxos works by running multiple Paxos rounds across the nodes, auto leader election, log replication to handle failure scenarios, and several other improvements.

## References

 - [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf) [Leslie Lamport]
 - [Paxos lecture (Raft user study)](https://youtu.be/JEpsBg0AO6o) [Diego Ongaro & John Ousterhout]