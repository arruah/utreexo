package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/btcsuite/btcd/chaincfg"

	bridge "github.com/mit-dci/utreexo/bridgenode"
	"github.com/mit-dci/utreexo/csn"
)

var msg = `
Usage: utreexo COMMAND [OPTION]
A dynamic hash based accumulator designed for the Bitcoin UTXO set

Commands:
  ibdsim         simulates an initial block download with ttl.testnet.txos as an input
  genproofs      generates proofs from the ttl.testnet.txos file
OPTIONS:
  -net=testnet   configure whether to use testnet. Optional.
  -net=regtest   configure whether to use regtest. Optional.
`

// bit of a hack. Standard flag lib doesn't allow flag.Parse(os.Args[2]). You need a subcommand to do so.
var optionCmd = flag.NewFlagSet("", flag.ExitOnError)
var netCmd = optionCmd.String("net", "testnet",
	"Target network. (testnet, regtest, mainnet) Usage: '-net=regtest'")

func main() {
	// check if enough arguments were given
	if len(os.Args) < 2 {
		fmt.Println(msg)
		os.Exit(1)
	}
	optionCmd.Parse(os.Args[2:])
	var param chaincfg.Params // wire.BitcoinNet
	if *netCmd == "testnet" {
		param = chaincfg.TestNet3Params
	} else if *netCmd == "regtest" {
		param = chaincfg.RegressionNetParams
	} else if *netCmd == "mainnet" {
		param = chaincfg.MainNetParams
	} else {
		fmt.Println("Invalid net flag given.")
		fmt.Println(msg)
		os.Exit(1)
	}
	//listen for SIGINT, SIGTERM, or SIGQUIT from the os
	sig := make(chan bool, 1)
	handleIntSig(sig)

	switch os.Args[1] {
	case "ibdsim":
		err := csn.RunIBD(&param, sig)
		if err != nil {
			panic(err)
		}
	case "genproofs":
		err := bridge.BuildProofs(param, sig)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println(msg)
		os.Exit(0)
	}
}

func handleIntSig(sig chan bool) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-s
		sig <- true
	}()
}
