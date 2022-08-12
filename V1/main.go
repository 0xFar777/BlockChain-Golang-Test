package main

import (
	"fmt"
)

func main() {
	createBlockChain := NewBlockChain()
	createBlockChain.AvarageBlock("ABCDEFG")
	createBlockChain.AvarageBlock("HIJKLMN")
	createBlockChain.AvarageBlock("OPQRSTU")
	createBlockChain.AvarageBlock("VWXYZ+-")
	for i, block := range createBlockChain.Block {
		fmt.Printf("====== 当前区块高度为:%d ======\n", i)
		fmt.Printf("前区块哈希是:%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希是:%x\n", block.Hash)
		fmt.Printf("区块数据为:%s\n", block.Data)
	}
}
