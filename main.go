package main

import (
	"io/ioutil"
	"log"

	"github.com/decred/dcrlnd/lnrpc"
	"github.com/decred/dcrlnd/macaroons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	macaroon "gopkg.in/macaroon.v2"
)

func main() {
	// Create the TransportCredentials and DialOptions
	// for gRPC client connection.
	creds, err := credentials.NewClientTLSFromFile("tls.cert", "")
	if err != nil {
		log.Printf("Unable to read cert file: %v", err)
		return
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	// Load the macaroon file.
	macBytes, err := ioutil.ReadFile("admin.macaroon")
	if err != nil {
		log.Println(err)
		return
	}
	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macBytes); err != nil {
		log.Println(err)
		return
	}

	// Append macaroon credentials to the dial options.
	opts = append(
		opts,
		grpc.WithPerRPCCredentials(macaroons.NewMacaroonCredential(mac)),
	)

	// Dial to the dcrlnd.
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Printf("Unable to dial to dcrlnd's gRPC server: %v", err)
		return
	}

	// Connect and create a client for dcrlnd.
	dcrlnd := lnrpc.NewLightningClient(conn)

	log.Println(dcrlnd)
	return
}
