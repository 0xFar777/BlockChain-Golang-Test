package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const BlockChainDb = "blockChain.db"
const BlockBucket = "blockChain"

//定义BlockChain结构
/*区块链所有的数据存储在Bolt数据库中,我们说区块链是分布式的,只要有一个节点在工作,
  那么这个系统就是持续运行的,单个节点想要加入挖矿系统,需要从Bolt数据库中读取最后
  一个区块的哈希然后才能开始挖新的区块,但作为矿工是不知道现在区块链运行到哪一个区
  块高度的(即矿工并不知道需要从Bolt数据库中读取哪一个key),因此我们需要在Bolt数据
  库中加入一个字段,这个字段存储当前系统最后一个区块的哈希值,单个节点开始挖矿,只需
  要读取该字段,便可以知道现在区块链运行到哪个高度了
*/

//因此需要重定义BlockChain的数据结构
type BlockChain struct {
	db *bolt.DB
	//Bolt数据库中用tail存储最后一个区块的哈希,
	//节点加入后读取该字段即可知道当前区块链运行到哪一个高度了
	tail []byte
}

/*以下是对应Bolt数据库中的结构
      key             value
  block1.Hash ---> block1的内容
  block2.Hash ---> block2的内容
  block3.Hash ---> block3的内容
  ......
     lastHash ---> blockN.Hash
*/

//创建创世区块
func GenesisBlock(address string, data string) *Block {
	tx := NewCoinbase(address, data)
	return NewBlock([]*Transaction{tx}, []byte{})
}

//以下这个函数有两个用意：
/*如果该区块链不存在,意味着你是本区块链的创建者,需要你创建一个db文件和开辟一个bucket,
  并添加创世区块;如果该区块链存在,意味着你可能是该区块链系统新加入(或退出重进)的矿工,
  需要读取该区块链当前的运行高度,得到最后一个区块的哈希值再开始挖矿
*/
func NewBlockChain(address string) *BlockChain {
	//用lastHash存储区块链当前最后一个区块的哈希
	var lastHash []byte
	//打开bolt数据库
	//打开"blockChain.db"文件,如果文件不存在则自动创建
	//0600表示给文件操作者赋予读写权限
	db, err := bolt.Open(BlockChainDb, 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败\n")
	}
	//该数据库db文件不能关闭！！！！！！

	db.Update(func(tx *bolt.Tx) error {
		//打开"blockChain.db"文件中的"blockChain"抽屉
		bucket := tx.Bucket([]byte(BlockBucket))
		//如果该抽屉不存在,意味着该区块链不存在,则需要创建区块链
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(BlockBucket))
			//创建区块链可能不成功
			if err != nil {
				log.Panic("创建bucket(" + BlockBucket + ")失败\n")
			}
			//如果该区块链存在或该区块链刚被创建成功,则往该区块链写入创世区块
			//生成创世区块,并作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock(address, "FaTeacher")
			//更新Bolt数据库
			//将创世区块的哈希最为key,创世区块的内容的字节流作为value
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			//更新lastHash这个key,其value为创世区块的哈希
			bucket.Put([]byte("lastHash"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		} else {
			//如果该抽屉存在,意味着该区块链存在
			//直接读取lastHash字段得到最后一个区块的哈希即可
			lastHash = bucket.Get([]byte("lastHash"))
		}
		return nil
	})
	//最后返回区块链结构即可
	return &BlockChain{
		db:   db,
		tail: lastHash,
	}
}

//创建普通区块(非创世区块)
func (bc *BlockChain) AvarageBlock(txs []*Transaction) {
	db := bc.db
	//取得上一个区块的哈希
	lastHash := bc.tail
	//创建新的区块
	newBlock := NewBlock(txs, lastHash)

	db.Update(func(tx *bolt.Tx) error {
		//打开"blockChain.db"文件中的"blockChain"抽屉
		bucket := tx.Bucket([]byte(BlockBucket))
		if bucket == nil {
			log.Panic("打开bucket失败\n")
		}
		//打开成功
		//1. Bolt数据库中添加新的区块数据
		bucket.Put(newBlock.Hash, newBlock.Serialize())
		//2. Bolt数据库中更新lastHash
		bucket.Put([]byte("lastHash"), newBlock.Hash)
		//3. 更新内存中的lastHash
		bc.tail = newBlock.Hash
		return nil
	})
}

func (bc *BlockChain) FindNeedUTXOs(from string,
	amount float64) (map[string][]uint64, float64) {
	var spentOutputs = make(map[string][]uint64)
	var UTXOS = make(map[string][]uint64)
	var total float64
	blockChain := bc.Iterator()
	for {
		//遍历所有区块
		block := blockChain.Next()
		//遍历区块的每个交易
		for _, transaction := range block.Transactions {
		OUTPUT:
			//遍历每个区块的Output
			for i, output := range transaction.TXOutputs {
				//剔除有关消耗
				if spentOutputs[string(transaction.TXID)] != nil {
					for _, spentOutputsIndex := range spentOutputs[string(transaction.TXID)] {
						if spentOutputsIndex == uint64(i) {
							continue OUTPUT
							//跳转到for i, output := range transaction.TXOutputs语句,但是i+1
						}
					}
				}
				//剔除无关消耗和无关未消耗
				//只有有关未消耗才append进UTXOS中
				if output.Address == from {
					if total < amount {
						UTXOS[string(transaction.TXID)] = append(UTXOS[string(transaction.TXID)], uint64(i))
						total += output.Value
						if total >= amount {
							return UTXOS, total
						}
					}
				}
			}
			//遍历每个区块的Input
			for _, input := range transaction.TXInputs {
				//筛选有关消耗
				if input.Sig == from {
					spentOutputs[string(input.TXID)] = append(spentOutputs[string(input.TXID)], uint64(input.Index))
				}
			}
		}
		//前区块哈希为空,终止遍历
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return UTXOS, total
}
