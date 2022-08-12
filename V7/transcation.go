package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward float64 = 12.5

//定义交易的结构
type Transaction struct {
	//一笔交易的ID(是一个哈希,由交易的字节流取哈希生成)
	TXID []byte
	//一笔交易的输入和输出,都可能有多个,所以用数组
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

//定义交易的输入结构
type TXInput struct {
	//交易要说明钱的来源,钱来源于另外一笔交易的输出
	//因此要知道钱来源的那笔交易的ID及其在Output中的索引
	TXID  []byte
	Index int64
	//一笔Input需要提供私钥的签名供矿工验证
	//注意:不是对整个Transcation进行签名,而是对钱的来源进行签名(即另一笔交易的Output)
	Sig []byte
	//一笔Input需要提供公钥供矿工验证
	PubKey []byte
}

//定义交易的输出结构
type TXOutput struct {
	//交易需要知道转给对方多少钱
	Value float64
	//Output要存储公钥哈希,否则无法被以后的一笔Input当作钱的来源
	//矿工验证某Input的合法性时需要其引用的Output的公钥哈希
	PubKeyHash []byte
}

//交易ID(哈希)生成函数,是属于结构体Transcation的方法
func (tx *Transaction) SetHash() {
	//将编码后的数据放在buffer中
	var buffer bytes.Buffer
	//func gob.NewEncoder(w io.Writer) *gob.Encoder
	//NewEncoder函数返回一个Encoder类型的指针,该指针将编码后的数据写入"w"中
	encoder := gob.NewEncoder(&buffer)
	//Encode函数将tx进行序列化,得到的字节流存放到buffer中
	err := encoder.Encode(&tx)
	if err != nil {
		log.Panic("设置交易哈希失败\n")
	}
	//对buffer中字节流取哈希
	Hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = Hash[:]
}

func (txop *TXOutput) Lock(address string) {
	//将地址转换位160位公钥哈希并赋值给TXOutput结构体的PubKeyHash
	txop.PubKeyHash = GetPubkeyHashFromAdderss(address)
}

//设置Output,较之前版本,Output中的地址已经改成了地址对应的公钥的哈希
func SetTXOutput(value float64, address string) *TXOutput {
	//构造一个output
	output := TXOutput{
		Value: value,
	}
	//将地址转换为其对应的160位公钥哈希(不是256位公钥哈希)
	output.Lock(address)
	return &output
}

//实现Coinbase交易(即铸币)函数
func NewCoinbase(address string, data string) *Transaction {
	input := TXInput{
		//Coinbase交易的输入TXID为空,同样的Index也为空
		TXID:  []byte{},
		Index: -1,
		//Coinbase交易无需在Sig中填自己私钥的签名,因此直接填nil
		Sig: nil,
		//Coinbase交易无需在PubKey中填自己私钥的签名,因此可以填任何信息
		PubKey: []byte(data),
	}
	//对于Coinbase交易,Value是出块奖励
	output := SetTXOutput(reward, address)
	//封装输入输出到整笔交易中
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{*output},
	}
	//对Coinbase交易取哈希
	tx.SetHash()
	return &tx
}

//实现普通转账
func NewTransfer(from, to string, amount float64, bc *BlockChain) *Transaction {
	//创建交易后要进行数字签名--->签名需要交易发起者的私钥--->因此要打开钱包获取私钥
	//创建交易后要验证所有的TXInput是否花的都是交易创建者自己的钱,因此还需要交易发起者的公钥
	//拿到交易发起者的公钥之后,对公钥取160位哈希,看是否与Input里引用的Output的PubKeyHash相同
	//如果相同,证明该笔交易确实是花的自己的钱
	ws := NewWallets()
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Printf("没有找到地址为：%s的钱包,创建交易失败\n", from)
		//创建失败就要返回nil了,否则会继续执行后续代码
		return nil
	}
	PublicKey := wallet.PublicKey
	PrivateKey := wallet.PrivateKey

	PubKeyHash160 := SetPubkeyHash(PublicKey)

	//获取最为合理的UTXO集合
	utxo, resValue := bc.FindNeedUTXOs(PubKeyHash160, amount)
	//当返回的UTXO集合中币的总数量小于想要转账的数量时,说明账户余额不足
	if resValue < amount {
		extraValue := amount - resValue
		fmt.Printf("余额不足,转账失败,还差%f才可转账成功\n", extraValue)
		return nil
	}
	//开始构建该笔交易的所有input和output:
	var inputs []TXInput
	var outputs []TXOutput
	//构造[]TxInput
	//遍历返回的UTXO,找到所有待消耗的Output所在的交易ID和Index
	for TxId, TxIndex := range utxo {
		//一笔交易可能消耗同个交易ID下的多个Output
		for _, i := range TxIndex {
			input := TXInput{
				TXID:   []byte(TxId),
				Index:  int64(i),
				Sig:    nil,
				PubKey: PublicKey,
			}
			inputs = append(inputs, input)
		}
	}
	//构造[]TxOutput
	output := SetTXOutput(amount, to)
	outputs = append(outputs, *output)
	//如果返回的UTXO集合币的总数量大于想要转账的数量,则需要找零
	//找零是构造一个自己把多余的钱转给自己的output
	if resValue > amount {
		reback := resValue - amount
		output = SetTXOutput(reback, from)
		outputs = append(outputs, *output)
	}
	//构造交易结构
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  inputs,
		TXOutputs: outputs,
	}
	//设置交易哈希(Id)
	tx.SetHash()
	//对该笔交易进行数字签名,签名的详细过程见sign.go文件
	bc.SignTransaction(&tx, PrivateKey)
	fmt.Printf("交易成功,%s向%s转了%f\n", from, to, amount)
	return &tx
}
