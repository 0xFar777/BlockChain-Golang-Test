package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
)

//辅助函数 -- 将uint64转换为byte类型
//--  因为在求当前区块的哈希时，需要将区块的所有信息进行拼接在取哈希，转换成相同的类型便于拼接
func Uint64ToByte(num uint64) []byte {
	// func Write(w io.Writer, order ByteOrder, data interface{}) error
	//将data的binary编码格式写入w，data必须是定长值、定长值的切片、定长值的指针。
	//order指定写入数据的字节序，写入结构体时，名字中有'_'的字段会置为0。
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panicln(err)
	}
	return buffer.Bytes()
}

//定义生成当前区块哈希的方法
func (block *Block) SetHash() {
	//拼接数据
	//方法一:(bytes.Join()方法)
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkleRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}
	//将二维的切片数组拼接起来,拼成一个一维的切片
	blockInfo := bytes.Join(tmp, []byte{})

	//方法二:(append方法：将block里的内容依次拼接)
	// var blockInfo []byte
	// blockInfo = append(blockInfo, byte(block.Version))
	// blockInfo = append(blockInfo, block.PrevHash...)
	// blockInfo = append(blockInfo, block.MerkleRoot...)
	// blockInfo = append(blockInfo, byte(block.TimeStamp))
	// blockInfo = append(blockInfo, byte(block.Difficulty))
	// blockInfo = append(blockInfo, byte(block.Nonce))
	// blockInfo = append(blockInfo, block.Data...)

	//使用sha256生成哈希
	hash := sha256.Sum256(blockInfo)
	//将生成的哈希赋值给block.Hash
	block.Hash = hash[:]
}
