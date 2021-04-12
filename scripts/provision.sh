#!/bin/bash

# Initialize ENV Variables peers
# A list of all the ports
# of the paxos peers
declare -a peers=()
declare -a peer_id_list=()

# Name of the Paxos cluster network
# connecting the peers
network="paxos_network"

# number of peers to be provisioned
# Default 3 peers are provisioned
peers_count=$1

# Err check number of peers
# If no peers count is given defaul to 3
if [[ $peers_count -eq "" ]]; then
    peers_count=3
fi

echo "Number of peers: $peers_count"

# Exit when there are less than 3 peers
if [[ $peers_count -le 2 ]]; then
    echo "Number of peers cannot be less than 3"
    exit 255
fi

# Exit when there are more than 1000 peers
if [[ $peers_count -ge 1000 ]]; then
    echo "Number of peers cannot be more than 1000"
    exit 255
fi

# Exit when there are an even
# number of peers provided
if [[ $(($peers_count % 2)) -eq 0 ]]; then
    echo "Number of peers cannot be an even value"
    exit 255
fi

# Check if port is available and then
# append to peers starting from 8000
available_port=8000
provisioned_ports_count=0

echo "Cleaning previous stale peers"
docker ps -a | awk '$2 ~ /paxos/ {print $1}' | xargs -I {} docker rm -f {} 2>/dev/null
docker network ls | grep paxos_network
if [[ $? -eq 0 ]]; then
    docker network rm "$network" 2>/dev/null
fi

echo "Reserving ports for peers"

for port in {8000..9000}; do
    if [[ provisioned_ports_count -eq peers_count ]]; then
        break
    fi
    sudo lsof -i:${port} | awk '$3 ~ /root/ {print $2} ' |xargs -i sudo kill -9 {} 2>/dev/null
    peers+=($port)
    ((provisioned_ports_count++))
    
done

if [[ provisioned_ports_count -ne peers_count ]]; then
    echo "Unable to reserve ports for peers"
    exit 255
fi

echo "Reserved ports:" ${peers[*]}
comma_separated_peers=$(
    IFS=,
    echo "${peers[*]}"
)

# Docker create peers from peer list
# and pass PORT = peers[[i]]
echo "Provisioning Paxos Docker Cluster"

echo "Building Paxos Docker Image"
docker build -t paxos -f Dockerfile .

if [[ $? -ne 0 ]]; then
    echo "Unable To Build Paxos Docker Image"
    exit 255
fi
paxos_data_path="paxos_logs"
echo "clear Paxos Volumes"
#docker volume ls |grep "${paxos_volume_prefix}"|awk '{print $2}'|xargs -i docker volume rm {}
rm -rf ${HOME}/${paxos_data_path} 2>/dev/null
mkdir -p ${HOME}/${paxos_data_path}
paxos_volume_prefix="n"

echo "Building Paxos Cluster Network"
docker network create "$network"

for ((id = 0; id < $peers_count; ++id)); do
    peer_id_list+=(peer-$id)
done

comma_separated_peer_id_list=$(
    IFS=,
    echo "${peer_id_list[*]}"
)

for peer_index in "${!peers[@]}"; do
    volume_path="${HOME}/${paxos_data_path}/${paxos_volume_prefix}${peer_index}" ;
    mkdir -p ${volume_path}
    docker run -p "${peers[$peer_index]}":8080 --restart=on-failure --net $network --privileged=true -v ${volume_path}":/home/log" -e "PEERS="$comma_separated_peer_id_list"" -e "ME"=$peer_index"" -e "NETWORK="$network"" --name="peer-$peer_index" -d paxos
done

# Docker list peers on success
echo "Paxos Cluster Nodes"
docker ps | grep 'paxos'
docker network ls | grep "$network"
