package main

//定义区块链结构
type BlockChain struct {
	//区块链 == N个区块相连,故为[]*Block
	Block []*Block
}

//创建创世区块
func GenesisBlock() *Block {
	return NewBlock("FaTeacher的公链", []byte{})
}

//创建区块链
func NewBlockChain() *BlockChain {
	//生成创世区块,并作为第一个区块添加到区块链中
	genesisBlock := GenesisBlock()
	bc := BlockChain{
		Block: []*Block{genesisBlock},
	}
	return &bc
}

//创建普通区块(非创世区块)
func (bc *BlockChain) AvarageBlock(data string) {
	//获取上一个区块
	lastBlock := bc.Block[len(bc.Block)-1]
	//取得上一个区块的哈希
	prevHash := lastBlock.Hash
	//创建新的区块
	newBlock := NewBlock(data, prevHash)
	//将新的区块添加到区块链中
	bc.Block = append(bc.Block, newBlock)
}
