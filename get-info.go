package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/davecgh/go-spew/spew"
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
	conn, err := grpc.Dial("localhost:10009", opts...)
	if err != nil {
		log.Printf("Unable to dial to dcrlnd's gRPC server: %v", err)
		return
	}

	// Connect and create a client for dcrlnd.
	dcrlnd := lnrpc.NewLightningClient(conn)

	// Get info from dcrlnd's node
	infoReq := &lnrpc.GetInfoRequest{}
	nodeInfo, err := dcrlnd.GetInfo(context.Background(), infoReq)
	if err != nil {
		log.Printf("Unable to get info: %v", err)
		return
	}

	spew.Dump(nodeInfo)

	fmt.Println("\nIdentity Pubkey:", nodeInfo.IdentityPubkey)
	fmt.Println("Alias:", nodeInfo.Alias)
	fmt.Println("Pending Channels:", nodeInfo.NumPendingChannels)
	fmt.Println("Active Channels:", nodeInfo.NumActiveChannels)
	fmt.Println("Peers:", nodeInfo.NumPeers)
	fmt.Println("Block Height:", nodeInfo.BlockHeight)
	fmt.Println("Block Hash:", nodeInfo.BlockHash)
	fmt.Println("Synced:", nodeInfo.SyncedToChain)
	fmt.Println("Is testnet:", nodeInfo.Testnet)
	fmt.Println("URIs:", nodeInfo.Uris)
	fmt.Println("Best header timestamp:", nodeInfo.BestHeaderTimestamp)
	fmt.Println("Version:", nodeInfo.Version)
	fmt.Println("Inactive Channels:", nodeInfo.NumInactiveChannels)
	fmt.Println("Chain:", nodeInfo.Chains[0].Chain) // the lnd started for use more than on chain, but at this moment just support one
	fmt.Println("Network:", nodeInfo.Chains[0].Network)
	fmt.Println("Color:", nodeInfo.Color)
	fmt.Println("Synced to graph:", nodeInfo.SyncedToGraph)
	return
}
