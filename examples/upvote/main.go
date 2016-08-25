package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-steem/rpc"
	"github.com/go-steem/rpc/apis/database"
	"github.com/go-steem/rpc/transports/websocket"

	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func run() (err error) {
	// Process flags.
	flagAddress := flag.String("rpc_endpoint", "ws://localhost:8090", "steemd RPC endpoint address")
	flag.Parse()

	url := *flagAddress

	// Process args.
	args := flag.Args()
	if len(args) != 3 {
		return errors.New("3 arguments required")
	}
	voter, author, permlink := args[0], args[1], args[2]

	// Start catching signals.
	var interrupted bool
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Drop the error in case it is a request being interrupted.
	defer func() {
		if err == websocket.ErrClosing && interrupted {
			err = nil
		}
	}()

	// Instantiate the WebSocket transport.
	t, err := websocket.NewTransport(url)
	if err != nil {
		return err
	}

	// Use the transport to get an RPC client.
	client := rpc.NewClient(t)
	defer func() {
		if !interrupted {
			client.Close()
		}
	}()

	// Start processing signals.
	go func() {
		<-signalCh
		fmt.Println()
		log.Println("Signal received, exiting...")
		signal.Stop(signalCh)
		interrupted = true
		client.Close()
	}()

	// Get the props for the transaction.
	props, err := client.Database.GetDynamicGlobalProperties()
	if err != nil {
		return err
	}

	// Prepare the transaction.
	tx := transactions.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(props.HeadBlockNumber),
		RefBlockPrefix: transactions.RefBlockPrefix(props.HeadBlockID),
		Expiration:     &expiration,
	})

	tx.PushOperation(&operations.Vote{
		Voter:    voter,
		Author:   author,
		Permlink: permlink,
		Weight:   10000,
	})

	tx.PushWIF(wifKey)

	if err := tx.Sign(); err != nil {
		return err
	}

	// Broadcast.
	raw, err := client.NetworkBroadcast.BroadcastTransactionSynchronousRaw(tx.Transaction)
	if err != nil {
		return err
	}
	fmt.Println("RESPONSE:", string(*raw))

	// Success!
	return nil
}

func readWIF() (string, error) {
	content, err := ioutil.ReadFile("wif.txt")
	if err != nil {
		return "", errors.Wrap(err, "failed to read wif.txt")
	}

	return strings.TrimSpace(string(content)), nil
}
