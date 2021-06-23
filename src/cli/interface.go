package cli

import (
	"flag"
	"fmt"
	"github.com/arielcoin/arielcoin/blockchain"
	"github.com/davecgh/go-spew/spew"
	"os"
)

func Run(b blockchain.Blockchain) {
	//subcomand definition
	help := flag.NewFlagSet("help", flag.ExitOnError)
	listBlocks := flag.NewFlagSet("listBlocks", flag.ExitOnError)
	createBlock := flag.NewFlagSet("createBlock", flag.ExitOnError)
	stake := flag.NewFlagSet("stake", flag.ExitOnError)

	//args
	data := createBlock.String("data", "default", "block data")
	blockHash := stake.String("blockHash", "", "block hash")

	if len(os.Args) < 2 {
		fmt.Println("uknown values")
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help":
		help.Parse(os.Args[2:])
		printHelp()
	case "listBlocks":
		listBlocks.Parse(os.Args[2:])
		spew.Dump(b.Blocks[0])
		fmt.Println("Blocks...")
	case "createBlock":
		createBlock.Parse(os.Args[2:])
		ccreateBlock(*data, &b)
	case "stake":
		stake.Parse(os.Args[2:])
		stakeHash(*blockHash, b)
	default:
		fmt.Println("uknown values")
		printHelp()
		return
	}

}

func printHelp() {
	fmt.Println("ArielCoin Interface")
	fmt.Println()
	fmt.Println("help - This message")
	fmt.Println("listBlocks - Prints the Blockchain")
	fmt.Println("createBlock - Create a block")
	fmt.Println("stake - Stake a block")
	fmt.Println()
}

func ccreateBlock(data string, b *blockchain.Blockchain) {

	//get old block
	oldBlock := b.Blocks[len(b.Blocks)-1]
	block := blockchain.BlockNew(blockchain.Version, data, "", oldBlock)

	if !block.Valid(oldBlock) {
		fmt.Println("Unvalid Block")
	}

	b.Blocks = append(b.Blocks, block)

	fmt.Println("Block Created!")
}

func stakeHash(hash string, b blockchain.Blockchain) {
	fmt.Println(hash)
	fmt.Println(b.Blocks)
	fmt.Println("Staked!")
}
