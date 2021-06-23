/*
This is the blockchain implementation
with the structs hashing and cecking if the block
is valid

if you contributed here please list your
github userbame below
========== Contributers
aryel - arydevy
*/

package blockchain

import(
	"crypto/sha256"
	"encoding/hex"
	"bytes"
	"encoding/gob"
	"time"

	"log"
)

const(
	Version = "1"
)

func Handle(e error){
		log.Fatal(e)
		//print(e)
		//panic(e)
}


//TODO:Add RingSignature

//Block - a item in the blockchain
type Block struct{
	Version string //Block verion
	Index int //Block index
	Timestamp string //time of the block creation
	Data string //Data
	Hash string //block hash
	PrevHash string //last block's hash
	Validator string //the validator of the block
}


type Blockchain struct {
	Blocks []Block
}

func (b *Blockchain) Init(){
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{
		Version:Version,
		Index:0,
		Timestamp:t.String(),
		Data:Hash(t.String()),
		Hash:HashBlock(genesisBlock),
		PrevHash:"",
		Validator:""}

	b.Blocks = append(b.Blocks, genesisBlock)

	Wdb(genesisBlock)
}

//Usefull Hashing function
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


func BlockNew(version,data,address string,oldBlock Block) Block{
	var newBlock Block

	t:=time.Now()

	newBlock.Version=version
	newBlock.Index=oldBlock.Index+1
	newBlock.Timestamp=t.String()
	newBlock.Data=data
	newBlock.PrevHash=oldBlock.Hash
	newBlock.Hash=HashBlock(newBlock)
	newBlock.Validator=address

	return newBlock
}

func (b Block) Valid(oldBlock Block) bool{
	if oldBlock.Index+1 != b.Index{
		return false
	}
	if oldBlock.Hash != b.PrevHash{
		return false
	}
	if HashBlock(b) != b.Hash{
		return false
	}

	return true
}


func (b Block) Encode() []byte{
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Decode(data []byte) *Block{
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}