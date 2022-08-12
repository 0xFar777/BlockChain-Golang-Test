package main

func main() {
	//初始化区块链
	bc := NewBlockChain("Satoshi")
	//初始化命令行结构体 -- 使用命令行完成  下面的操作
	cli := CommandLine{
		bc: bc,
	}
	//运行cli的读取命令行函数
	cli.Run()
}
