package main

import(
	"github.com/arielcoin/arielcoin/blockchain"
	"github.com/arielcoin/arielcoin/cli"
)

func main(){
	//init Blockchain
	var blockchain Blockchain
	blockchain.Init()

	//Run cli interface
	cli.Run(blockchain)
}