# Awesome `dcrlnd` (Decred Lightning Network)

Everything that you wnat to know about Decred Lightninh Network(Layer 2), or almost.(WIP)

## TLS and Macaroons

Every program needs the tls and macarron files, to simplify this proccess you can use
the `dcrlnd` docker-compose to start a simnet environment and the script `getCertAndMac.sh`
to copy data from dcrlnd nodes container.

## Create a gRPC connection with your node

Use the `main.go` and edit the config:
```
credentials.NewClientTLSFromFile("path_to_tls.cert", "")

ioutil.ReadFile("path_to_admin.macaroon")

grpc.Dial("node_host:node_port", opts...)
```

Run separate file:
`go run main.go`
