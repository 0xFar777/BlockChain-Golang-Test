package main

import (
	"fmt"
	"os"
)

//把命令行操作要用到的内容装进结构体,方便后面对其进行操作
type CommandLine struct {
	bc *BlockChain
}

//客户端操作需要给用户一些指令提示
const Usage = `
	addBlock DATA     "添加区块"
	printChain        "正向打印区块链"

`

//用户输入命令处理函数
func (cli *CommandLine) Run() {
	//先获取用户在命令行输入的内容
	args := os.Args
	//如果用户在命令行中只输入了一个参数,则提示用户重新输入
	if len(args) < 2 {
		fmt.Print(Usage)
		//此处要退出函数,否则会继续执行该函数后面的内容
		return
	}

	//如果命令行参数大于一个
	cmd := args[1]
	//判断用户在命令行输入的第二个参数:
	switch cmd {
	//如果第二个参数是addBlock,意味着用户想要新增区块
	case "addBlock":
		//命令行的第三个参数即为新增Block的block.data
		if len(args) == 3 {
			data := args[2]
			cli.AddBlocks(data)
		} else {
			//如果在第二个参数为addBlock的情况下用户输入的参数不等于3,则用户输入出现了错误
			fmt.Println("无效命令,请重新输入")
			fmt.Print(Usage)
		}
		//如果第二个参数是addBlock,意味着用户想要读取区块链信息
	case "printChain":
		if len(args) == 2 {
			fmt.Println("执行打印区块命令")
			cli.PrintBlockChain()
		} else {
			fmt.Println("无效命令,请重新输入")
			fmt.Print(Usage)
		}
	default:
		fmt.Println("无效的命令,请重新输入")
		fmt.Print(Usage)
	}
}
