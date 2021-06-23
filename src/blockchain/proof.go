/*
Here is the Satking of blocks implementation

if you contributed here please list your
github userbame below
========== Contributers
aryel - arydevy
*/

package blockchain

import(
	"fmt"
	"time"
	"math/rand"
)

const(
	MinTokens = 25
	MinValidators = 1 //this shoud be changed
)

type Validator struct{
	Tokens int
	Address string
	ProposedBlock Block

}

type ValidatorPool struct{
	Validators []Validator
}

//Propose adds a block to the winner validator
//returns true - success false - invalid block
func (v *Validator) Propose(block ,oldBlock Block) bool{
	valid := block.Valid(oldBlock)
	if !valid{
		fmt.Printf("[ERR]:Proposed invalid Block")
		return false
	}
	v.ProposedBlock = block
	return true
}


func Validate(pool *ValidatorPool,Blockchain *Blockchain){
	if len(pool.Validators) >= MinValidators{
		//Randomly pick a winner
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		winner:= pool.Validators[r.Intn(len(pool.Validators))]

		//TODO:Here you shoud listen to the non winner
		//validators to say if they trust or not the 
		//proposed block if not delete the winners
		//stack coins

		//get old block
		oldBlock:=Blockchain.Blocks[len(Blockchain.Blocks)-1]

		valid := winner.ProposedBlock.Valid(oldBlock)
		if !valid{
			fmt.Printf("[ERR]:Proposed invalid Block")
			//pool = delete(pool,winner) //TODO:fix the delete
			//TODO:Delete Staked coins
		}

		//Add Block
		Blockchain.Blocks = append(Blockchain.Blocks,winner.ProposedBlock)

		//TODO:Brodcast The added Block
	}

}