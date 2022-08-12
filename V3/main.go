package main

import "fmt"

func main() {
	createBlockChain := NewBlockChain()
	createBlockChain.AvarageBlock("ABCDEFG")
	createBlockChain.AvarageBlock("HIJKLMN")
	createBlockChain.AvarageBlock("OPQRSTU")
	createBlockChain.AvarageBlock("VWXYZ+-")

	//为区块链定义一个迭代器
	createIterator := createBlockChain.Iterator()
	for {
		currentBlock := createIterator.Next()
		fmt.Println("=============================================")
		fmt.Printf("当前区块的出块随机数为:%d\n", currentBlock.Nonce)
		fmt.Printf("前区块哈希是:%x\n", currentBlock.PrevHash)
		fmt.Printf("当前区块哈希是:%x\n", currentBlock.Hash)
		fmt.Printf("区块数据为:%s\n", currentBlock.Data)
		//迭代器如果读到创世区块,则停止读取
		if len(createIterator.CurrentHashPoint) == 0 {
			break
		}
	}
}
