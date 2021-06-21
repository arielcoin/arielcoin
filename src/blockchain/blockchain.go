package blockchain

//Block - a item in the blockchain
type Block struct{
	Index int
	Timestamp string
	Data string
	Hash string
	PrevHash string
	Validator string
}


type Blockchain struct {
	[]Block
}

//Usefull Hashing function
func Hash(s string) string{
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}