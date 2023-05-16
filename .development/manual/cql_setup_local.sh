podman network exists mincepierank-cass-network || podman network create mincepierank-cass-network

podman start mincepierank-cass || podman run --name mincepierank-cass \
    --network mincepierank-cass-network \
    -e CASSANDRA_CLUSTER_NAME="MincePieRank Test Cluster" \
    -e MAX_HEAP_SIZE="512M" \
    -e HEAP_NEWSIZE="100M" \
    --memory 1g \
    -p 9143:9042 \
    -d \
    docker.io/cassandra