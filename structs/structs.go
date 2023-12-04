package structs

type Block struct {
	BlockId           string
	PrevHash          string
	SetOfTransactions []Transaction
}
