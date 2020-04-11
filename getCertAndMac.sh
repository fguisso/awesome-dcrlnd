# This script is useful if you are running the
# docker-compose from dcrlnd repository.
#
# Usage: ./getCertAndMac <container_name>
echo "get all things from: $1"
rm admin.macaroon tls.cert
docker cp $1:/root/.dcrlnd/tls.cert .
docker cp $1:/root/.dcrlnd/data/chain/decred/simnet/admin.macaroon .
