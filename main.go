package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)


//Block - a item in the blockchain
type Block struct{
	Index int
	Timestamp string
	Data string
	Hash string
	PrevHash string
	Validator string
}


//Blockchain 
var Blockchain []Block
var tempBlocks []Block

var candidateBlocks = make(chan Block)

//PeerToPeer stuff
var anouncements = make(chan string)
var mutex = &sync.Mutex{}

var validators = make(map[string]int)

//Hashing
func Hash(s string) string{
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//HashBlock
func HashBlock(block Block) string{
	record := string(block.Index)+block.Timestamp+block.Data+block.PrevHash
	return Hash(record)
}

//MakeBlock creates a block
func MakeBlock(lastBlock Block,data string,address string) Block{
	var newBlock Block

	t:=time.Now()

	newBlock.Index=lastBlock.Index+1
	newBlock.Timestamp=t.String()
	newBlock.Data=data
	newBlock.PrevHash=lastBlock.Hash
	newBlock.Hash=HashBlock(newBlock)
	newBlock.Validator=address

	return newBlock
}

//CeckValidity makes sure the block is valid
func CeckValidity(newBlock,oldBlock Block) bool{
	if oldBlock.Index+1 != newBlock.Index{
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash{
		return false
	}
	if HashBlock(newBlock) != newBlock.Hash{
		return false
	}

	return true
}

//TCP implementation
func handelConn(conn net.Conn){
	
	//Close the connection before the function returns
	defer conn.Close()

	go func() {
		for {
			msg:=<-anouncements
			io.WriteString(conn,msg)

		}
	}()

	//Node Address
	var address string

	//allow user to set the number of tokens
	//to stake, The greather the number the greather the
	//chance to win
	io.WriteString(conn,"Enter Token Balance:")
	scanBalance:= bufio.NewScanner(conn)

	for scanBalance.Scan(){
		balance,err := strconv.Atoi(scanBalance.Text())
		if err != nil{
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}

		t:=time.Now()
		address = Hash(t.String())
		validators[address]=balance
		fmt.Println(validators)
		break
	}


	io.WriteString(conn,"\nEnter Data:")

	scanData:= bufio.NewScanner(conn)

	go func() {
		for scanData.Scan(){
			data:=scanData.Text()

			mutex.Lock()
			oldLastIndex :=Blockchain[len(Blockchain)-1]
			mutex.Unlock()

			//Create a block to be forged
			newBlock := MakeBlock(oldLastIndex,data,address)

			if CeckValidity(newBlock,oldLastIndex){
				candidateBlocks <- newBlock
			}

			io.WriteString(conn,"\nEnter data:")

		}
	}()

	//Simuate reciving broadcast
	for {
		
		time.Sleep(time.Minute)
		mutex.Lock()
		output,err := json.Marshal(Blockchain)
		mutex.Unlock()

		if err!=nil{
			log.Fatal(err )
		}

		io.WriteString(conn,string(output)+"\n")
	}
}

//PickWinner gathers nodes and then choses a winner
//by randomly sellecting
func PickWinner(){
	time.Sleep(30*time.Second)
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	pool:= []string{}
	if len(temp)>0{
		// slightly modified traditional proof of stake algorithm
		// from all validators who submitted a block, weight them by the number of staked tokens
		// in traditional proof of stake, validators can participate without submitting a block to be forged

		OUTER:
			for _,block:=range temp{
				//if in the pool skip
				for _,node:=range pool{
					if block.Validator == node {
						continue OUTER
					}
				}

				//Lock to prevent data race
				mutex.Lock()
				setValidators := validators
				mutex.Unlock()


				k,ok := setValidators[block.Validator]
				if ok{
					for i := 0;i<k;i++{
						pool = append(pool,block.Validator)
					}
				}
			}


			//Randomly pick a winner
			s := rand.NewSource(time.Now().Unix())
			r := rand.New(s)
			winner:= pool[r.Intn(len(pool))]

			//add the winner's block to the blockchain
			//and let the other nodes know
			for _,block:=range temp{
				if block.Validator == winner{
					mutex.Lock()
					Blockchain = append(Blockchain,block)
					mutex.Unlock()
					for _ = range validators{
						anouncements <- "\nWinner is:"+winner+"\n"
					}
					break
				}
			}
	}

	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}


func main(){
	err := godotenv.Load()
	if err!= nil{
		log.Fatal(err)
	}

	//Create genisis Block
	t:=time.Now()
	genisisBlock := Block{}
	genisisBlock = Block{0,t.String(),"genisisBlock",HashBlock(genisisBlock),"",""}
	spew.Dump(genisisBlock)
	Blockchain = append(Blockchain,genisisBlock)


	//start TCP
	server,err := net.Listen("tcp",":"+os.Getenv("ADDR"))
	if err != nil{
		log.Fatal()
	}
	defer server.Close()

	go func() {
		for candidate:=range candidateBlocks{
			mutex.Lock()
			tempBlocks = append(tempBlocks,candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		for candidate := range candidateBlocks{
			mutex.Lock()
			tempBlocks=append(tempBlocks,candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		PickWinner()
	}()

	for {
		conn,err:=server.Accept()
		if err !=nil{
			log.Fatal(err)
		}
		go handelConn(conn)
	}
}