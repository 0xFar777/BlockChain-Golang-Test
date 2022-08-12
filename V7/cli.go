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
	printChain        	"正向打印区块链"
	getBalance ADDRESS  "获取账户余额"
	transfer FROM TO AMOUNT MINER DATA  "发起转账"
	NewWallet           "创建新的钱包"
	listAddress         "列举所有的钱包地址"
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
	//如果第二个参数是printChain,意味着用户想要读取区块链信息
	case "printChain":
		if len(args) == 2 {
			fmt.Println("正在执行打印区块命令")
			cli.PrintBlockChain()
		} else {
			fmt.Println("无效命令,请重新输入")
			fmt.Print(Usage)
		}
	//如果第二个参数是getBalance,意味着用户想要获取自己或其他人账户的余额
	case "getBalance":
		//命令行第三个参数即为用户想要获取余额的账户地址
		if len(args) == 3 {
			fmt.Println("正在执行获取账户余额命令")
			address := args[2]
			cli.GetBalance(address)
		} else {
			fmt.Println("无效命令,请重新输入")
			fmt.Print(Usage)
		}
	//如果第二个参数是transfer,意味着用户想要发起一笔交易
	case "transfer":
		if len(args) == 7 {
			fmt.Println("正在执行转账命令")
			FROM := args[2]   //转账发起人
			TO := args[3]     //收款人
			AMOUNT := args[4] //转账金额
			MINER := args[5]  //矿工
			DATA := args[6]   //区块信息
			//注意,该版本尚未实现区块链网络层,所以暂时没办法完成矿工监听交易信息的操作
			cli.Transfer(FROM, TO, AMOUNT, MINER, DATA)
		} else {
			fmt.Println("无效命令,请重新输入")
			fmt.Print(Usage)
		}
	//如果第二个参数是NewWallet,意味着用户想要创建一个新的钱包地址
	case "NewWallet":
		if len(args) == 2 {
			fmt.Println("正在创建新的钱包地址")
			cli.CommandNewWallet()
		} else {
			fmt.Println("无效命令,请重新输入")
			fmt.Print(Usage)
		}
	//如果第二个参数是listAddress,意味着用户想要获取当前所有的钱包地址
	case "listAddress":
		if len(args) == 2 {
			fmt.Println("正在列举所有钱包地址")
			cli.CommandListAddress()
		} else {
			fmt.Println("无效命令,请重新输入")
			fmt.Print(Usage)
		}
	//如果上述都不满足,说明用户输入了无效命令
	default:
		fmt.Println("无效的命令,请重新输入")
		fmt.Print(Usage)
	}
}
