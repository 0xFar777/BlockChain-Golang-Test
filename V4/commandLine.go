package main

import "fmt"

//如果用户命令行输入的第二个参数是addBlock,则调用此函数添加区块
func (cli *CommandLine) AddBlocks(data string) error {
	cli.bc.AvarageBlock(data)
	fmt.Println("恭喜你,添加新区块成功")
	return nil
}

//如果用户命令行输入的第二个参数是printChain,则调用此函数获取区块链信息
func (cli *CommandLine) PrintBlockChain() {
	bc := cli.bc
	//为区块链定义一个迭代器
	createIterator := bc.Iterator()
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
