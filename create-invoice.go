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

	// Create and add an invoice.
	invoiceReq := &lnrpc.Invoice{
		Memo:   "Create my first invoice", // Description
		Value:  1000,                      // Amount in atoms
		Expiry: 3600,                      // Seconds to expiry invoice
	}

	invoice, err := dcrlnd.AddInvoice(context.Background(), invoiceReq)
	if err != nil {
		log.Printf("Error on generate an invoice: %v", err)
		return
	}

	spew.Dump(invoice)
	fmt.Println("\nPayment request:", invoice.PaymentRequest)
	return
}
