package main

import (
	"log"

	"github.com/boltdb/bolt"
)

/*该文件实现两个功能
1.遍历区块链所有信息
2.遍历指定区块号信息 //该版本尚未实现
*/

//把迭代器操作要用到的内容装进结构体,方便后面对其进行操作
type BlockIterator struct {
	db               *bolt.DB
	CurrentHashPoint []byte
}

//对于一个区块链构造迭代器
func (bc *BlockChain) Iterator() *BlockIterator {
	return &BlockIterator{
		db:               bc.db,
		CurrentHashPoint: bc.tail,
	}
}

//迭代遍历数据
func (it *BlockIterator) Next() *Block {
	db := it.db
	//从Bolt数据库取到的当前区块会被赋值到currentBlock
	var currentBlock Block
	//打开bucket
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BlockBucket))
		if bucket == nil {
			log.Panic("不存在此bucket,请检查")
		}
		//读取当前区块的信息
		resultTmp := bucket.Get(it.CurrentHashPoint)
		//将读取到数据反序列化
		currentBlock = DeSerialize(resultTmp)
		//更新迭代器指针
		it.CurrentHashPoint = currentBlock.PrevHash
		return nil
	})
	return &currentBlock
}
