podman network exists mincepierank-cass-network || podman network create mincepierank-cass-network

podman start mincepierank-cass || podman run --name mincepierank-cass \
    --network mincepierank-cass-network \
    -e CASSANDRA_CLUSTER_NAME="MincePieRank Test Cluster" \
    -p 9143:9042 \
    -d \
    docker.io/cassandra