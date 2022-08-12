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
	//交易需要签名进行验证
	//注意:不是对整个Transcation进行签名,而是对钱的来源进行签名(即另一笔交易的Output)
	Sig string
}

//定义交易的输出结构
type TXOutput struct {
	//交易需要知道转给对方多少钱
	Value float64
	//输出要存储公钥哈希,否则无法被以后的一笔交易当作钱的来源
	// PubKeyHash []byte
	//但是该版本还未实现签名,因此先用地址进行替代
	Address string
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

//实现Coinbase交易(即铸币)函数
func NewCoinbase(address string, data string) *Transaction {
	input := TXInput{
		//Coinbase交易的输入TXID为空,同样的Index也为空
		TXID:  []byte{},
		Index: -1,
		//Coinbase交易无需在Sig中填自己私钥的签名,因此可以填任何信息
		Sig: data,
	}
	output := TXOutput{
		//对于Coinbase交易,Value是出块奖励,Address是自己接收出块奖励的地址
		Value:   reward,
		Address: address,
	}
	//封装输入输出到整笔交易中
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
	}
	//对Coinbase交易取哈希
	tx.SetHash()
	return &tx
}

//实现普通转账
func NewTransfer(from, to string, amount float64, bc *BlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	//获取最为合理的UTXO集合
	utxo, resValue := bc.FindNeedUTXOs(from, amount)
	//当返回的UTXO集合中币的总数量小于想要转账的数量时,说明账户余额不足
	if resValue < amount {
		extraValue := amount - resValue
		fmt.Printf("余额不足,转账失败,还差%f才可转账成功\n", extraValue)
		return nil
	}
	//构造[]TxInput
	//遍历返回的UTXO,找到所有待消耗的Output所在的交易ID和Index
	for TxId, TxIndex := range utxo {
		//一笔交易可能消耗同个交易ID下的多个Output
		for _, i := range TxIndex {
			input := TXInput{
				TXID:  []byte(TxId),
				Index: int64(i),
				Sig:   from,
			}
			inputs = append(inputs, input)
		}
	}
	//构造[]TxOutput
	output := TXOutput{
		Value:   amount,
		Address: to,
	}
	outputs = append(outputs, output)
	//如果返回的UTXO集合币的总数量大于想要转账的数量,则需要找零
	//找零是构造一个自己把多余的钱转给自己的output
	if resValue > amount {
		reback := resValue - amount
		output = TXOutput{reback, from}
		outputs = append(outputs, output)
	}
	//构造交易结构
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  inputs,
		TXOutputs: outputs,
	}
	//设置交易哈希(Id)
	tx.SetHash()
	fmt.Printf("交易成功,%s向%s转了%f\n", from, to, amount)
	return &tx
}
