package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

//定义区块结构
type Block struct {
	//当前区块链协议版本
	Version uint64
	//前区块哈希
	PrevHash []byte
	//默克尔树根哈希
	MerkleRoot []byte
	//时间戳
	TimeStamp string
	//挖矿难度值
	Difficulty uint64
	//出块随机数
	Nonce uint64
	//当前区块哈希(Bitcoin中没有该字段,这里是为了实现方便)
	Hash []byte
	//区块交易数据
	Transactions []*Transaction
}

//定义创建区块的方法
//该函数要传prevHash,不同矿工传入的prevHash可能会不一样,
//因为不同的矿工眼里最长合法链可能不一致(概率较小)
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	block := Block{
		Version:    0,
		PrevHash:   prevHash,
		MerkleRoot: []byte{},
		//这里将时间戳转换为格式化时间,方便阅读
		TimeStamp:    TimeStampToStandardTime(int64(time.Now().Unix())),
		Difficulty:   0,
		Nonce:        0,
		Hash:         []byte{},
		Transactions: txs,
	}
	//更新默克尔树根
	block.MerkleRoot = block.MakeMerkleRoot()
	//创建一个Pow对象
	pow := NewProofOfWork(&block)
	//挖矿过程,成功便接收Hash和随机数
	Hash, Nonce := pow.Run()
	//对block的数据重新赋值
	block.Hash = Hash
	block.Nonce = Nonce
	return &block
}

//编码(序列化)函数
func (block *Block) Serialize() []byte {
	//将编码后的数据放在buffer中
	var buffer bytes.Buffer
	//func gob.NewEncoder(w io.Writer) *gob.Encoder
	//NewEncoder函数返回一个Encoder类型的指针,该指针将编码后的数据写入"w"中
	encoder := gob.NewEncoder(&buffer)
	//Encode函数将block进行序列化,得到的字节流存放到buffer中
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("编码失败\n")
	}
	return buffer.Bytes()
}

//解码(反序列化)函数
func DeSerialize(data []byte) Block {
	//将解码后的数据以Block结构体的形式返回
	var block Block
	//func gob.NewDecoder(r io.Reader) *gob.Decoder
	//NewDecoder函数返回一个Decoder类型的指针,意味着对"r"进行解码
	//首先要将buffer里的字节流数据转化为可读的形式,因此需要使用bytes.NewReader
	decoder := gob.NewDecoder(bytes.NewReader(data))
	//Decode函数将buffer输入流进行反序列化,得到的结果存储再block中
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码失败\n")
	}
	return block
}

//将区块体中的所有交易取一个默克尔根哈希放到区块头
//因为最后挖矿需要的仅仅是区块头的信息,不包含区块体
func (block *Block) MakeMerkleRoot() []byte {
	//获取当前区块所有交易信息
	txs := block.Transactions
	var tmp []byte
	//正常比特币版本是通过二叉树的形式来实现默克尔根哈希的,但这里采用哈希拼接
	for _, tx := range txs {
		tmp = append(tmp, tx.TXID...)
	}
	//生成哈希
	hash := sha256.Sum256(tmp)
	return hash[:]
}
