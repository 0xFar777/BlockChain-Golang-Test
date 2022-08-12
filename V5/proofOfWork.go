package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//定义工作量证明结构
type ProofOfWork struct {
	//当前区块数据
	Block *Block
	//当前区块的难度目标值
	//big.Int是一个非常大的数,有很丰富的方法，比如比较、赋值等
	target *big.Int
}

//生成目标域值
func NewProofOfWork(block *Block) *ProofOfWork {
	//先将区块结构赋值给ProofOfWork结构体
	pow := ProofOfWork{
		Block: block,
	}
	//定义当前区块难度的目标值,当前版本先定义一个固定值
	target := "0000100000000000000000000000000000000000000000000000000000000000"
	//定义一个big.Int类型的临时变量
	tmp := big.Int{}
	//将String类型的target转换为big.Int(16进制)类型并存储到tmp中
	tmp.SetString(target, 16)
	pow.target = &tmp
	return &pow
}

//提供不断进行Hash的函数
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var Nonce uint64
	var Hash [32]byte
	for {
		tmp1 := [][]byte{
			Uint64ToByte(pow.Block.Version),
			pow.Block.PrevHash,
			pow.Block.MerkleRoot,
			Uint64ToByte(pow.Block.TimeStamp),
			Uint64ToByte(pow.Block.Difficulty),
			Uint64ToByte(Nonce),
		}
		//将二维的切片数组拼接起来,拼成一个一维的切片
		blockInfo := bytes.Join(tmp1, []byte{})
		//使用sha256生成Hash
		Hash = sha256.Sum256(blockInfo)
		//我们的目的是将所求Hash与target进行比较,但是现在Hash还是[32]byte类型
		//而target是big.Int类型,因此需要把Hash也转换成big.Int类型
		tmp2 := big.Int{}
		tmp2.SetBytes(Hash[:])
		//将Hash与target进行比较,小于target意味着找到了符合要求的随机数,退出循环
		/*func (*big.Int).Cmp(y *big.Int) (r int)
		Cmp compares x and y and returns:
		-1 if x < y    0 if x == y    +1 if x > y
		*/
		if tmp2.Cmp(pow.target) == -1 {
			fmt.Println("======= 挖矿成功 =======")
			fmt.Printf("Hash: %x,Nonce: %d\n", Hash[:], Nonce)
			break
		} else {
			Nonce++
		}
	}
	return Hash[:], Nonce
}
