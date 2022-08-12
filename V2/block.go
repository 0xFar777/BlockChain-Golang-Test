package main

import "time"

//定义区块结构
type Block struct {
	//当前区块链协议版本
	Version uint64
	//前区块哈希
	PrevHash []byte
	//默克尔树根哈希
	MerkleRoot []byte
	//时间戳
	TimeStamp uint64
	//挖矿难度值
	Difficulty uint64
	//出块随机数
	Nonce uint64
	//当前区块哈希(Bitcoin中没有该字段,这里是为了实现方便)
	Hash []byte
	//区块数据
	Data []byte
}

//定义创建区块的方法
func NewBlock(data string, prevHash []byte) *Block {
	block := Block{
		Version:    0,
		PrevHash:   prevHash,
		MerkleRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}
	//创建一个Pow对象
	pow := NewProofOfWork(&block)
	//挖矿过程,成功便接收Hash和随机数
	Hash, Nonce := pow.Run()
	//对block的数据重新赋值
	block.Hash = Hash
	block.Nonce = Nonce
	return &block
}
