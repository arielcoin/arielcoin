package blockchain

import(
	badger "github.com/dgraph-io/badger/v3"
)

const (
	dbPath = "/tmp/aum"
)


//write a block to the db
func Wdb(b Block){
	block := b.Encode()
	hash := []byte(b.Hash) //get the last hash

	//Updating the db

	// Open the Badger database.
  	// It will be created if it doesn't exist.
  	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
  	
  	Handle(err)
  	// after the function finishes close the db
  	defer db.Close()
  	
  	err = db.Update(func(txn *badger.Txn) error {
  		
  		err  = txn.Set(hash,block) // add the block to the db
  		Handle(err)
  		err  = txn.Set([]byte("ls"),hash) //now change the lh(last hash)

  		Handle(err)

  		return err
	})

	Handle(err)


}
