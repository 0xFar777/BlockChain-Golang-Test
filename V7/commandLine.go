package main

import (
	"fmt"
	"log"
	"strconv"
)

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
		fmt.Printf("当前区块默克尔树根哈希为:%x\n", currentBlock.MerkleRoot)
		fmt.Printf("区块信息为:%v\n", currentBlock.Transactions[0].TXInputs[0].SetBytesDataToStringData())
		fmt.Printf("出块矿工为:%v\n", currentBlock.Transactions[0].TXOutputs[0].SetAddress())
		fmt.Printf("出块奖励为:%v\n", currentBlock.Transactions[0].TXOutputs[0].Value)
		fmt.Printf("出块时间戳为:%s\n", currentBlock.TimeStamp)
		fmt.Println("=============================================")
		//迭代器如果读到创世区块,则停止读取
		if len(createIterator.CurrentHashPoint) == 0 {
			break
		}
	}
}

//如果用户命令行输入的第二个参数是getBalance,则调用此函数获取账户的余额
func (cli *CommandLine) GetBalance(address string) {
	//获取余额前务必先检查地址的合法性
	if !IsVaildAddress(address) {
		fmt.Println("无效的地址,请检查")
		return
	}
	//通过FindUTXOByUser函数,来获取与该地址有关的所有未被消费(即未被当成另一笔交易的input)的Output
	UTXOS := cli.bc.FindUTXOByUser(address)
	var total float64
	//遍历所有的Output,并将每个Output内的Value累加即为用户余额
	for _, utxo := range UTXOS {
		total += utxo.Value
	}
	fmt.Printf("%s的余额为:%f\n", address, total)
}

//如果用户命令行输入的第二个参数是transfer,则调用此函数进行转账
//该版本尚未实现区块链网络层,因此矿工无法监听到其他节点产生的交易信息
//故这里设定为:发起一笔交易即可产生一个新区块
func (cli *CommandLine) Transfer(from, to, amount, miner, data string) {
	//务必先校验三个地址的合法性：
	var num = 0
	//检查发送人地址是否合法
	if !IsVaildAddress(from) {
		fmt.Println("发送人地址错误,请检查")
		num++
	}
	//检查收款人地址是否合法
	if !IsVaildAddress(to) {
		fmt.Println("接收人地址错误,请检查")
		num++
	}
	//检查矿工地址是否合法
	if !IsVaildAddress(miner) {
		fmt.Println("矿工地址错误,请检查")
		num++
	}
	//只要三个地址有一个输入错误,就停止转账
	if num > 0 {
		fmt.Println("无效交易,请检查")
		return
	}
	//先把从终端获取到的转账数额转成浮点型
	floatAmount, _ := strconv.ParseFloat(amount, 64)
	if floatAmount < 0 {
		fmt.Println("无效交易,转账金额务必大于等于0,请检查")
		return
	}
	//创建Coinbase交易
	coinbase := NewCoinbase(miner, data)
	//创建普通交易
	transfer := NewTransfer(from, to, floatAmount, cli.bc)
	if transfer == nil {
		log.Panicln("无效交易,请检查")
	}
	//创建新区块
	cli.bc.AvarageBlock([]*Transaction{coinbase, transfer})
}

//如果用户命令行输入的第二个参数是NewWallet,则调用此函数创建新钱包
func (cli *CommandLine) CommandNewWallet() {
	// //生成钱包(即生成一对公钥和私钥)
	// wallet := NewWallet()
	// //根据公钥生成地址
	// address := wallet.CreateAddress()
	// fmt.Printf("公钥: %v\n", wallet.PublicKey)
	// fmt.Printf("私钥: %v\n", wallet.PrivateKey)
	// fmt.Printf("地址: %s\n", address)

	//注意：以上代码为非公开代码,因为不能让其他人知道我的私钥,对外可见的一般只有地址
	//(一般公钥也是不对外可见的,只有需要验证的时候才会向验证者公开公钥)

	//以下为公开代码
	wallets := NewWallets()
	address := wallets.CreateWallet()
	fmt.Println("恭喜您,新地址创建成功")
	fmt.Printf("新地址为: %s\n", address)
}

func (cli *CommandLine) CommandListAddress() {
	wallets := NewWallets()
	addresses := wallets.ListAddress()
	//对map映射中所有的地址进行遍历
	for _, address := range addresses {
		fmt.Printf("地址：%s\n", address)
	}
}
