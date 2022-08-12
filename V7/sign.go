package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
)

/*签名详细步骤：

首先明白对谁签名：
新生成的交易中的input有一个字段叫sig,是用来存储签名生成的r和s字节流的,因此一笔交易中,有几个input
就要签名几次,签名不是对单笔input进行签名,而是对当前交易的所有input+所有output进行签名,每次签名变化
的只有一个参数：当前input的PubKey字段

然后来看一下签名函数：
func ecdsa.Sign(rand io.Reader, priv *ecdsa.PrivateKey, hash []byte)
(r *big.Int, s *big.Int, err error)
可见,签名需要三个参数,一个随机数,一个钱包私钥,还有每次签名需要用到的数据,这个数据
指的是每次新生成的Transaction的TXID

形象点理解就是这样：
        注：一笔交易中input的数量一定与所引用的output的数量相等,
		    但是与当前新生成交易中有多少output没有关系
		第一次签名需要用到下面的东西：
		input1:  TXID Index nil PubKey    output1:   value PubKeyHash
		input2:  TXID Index nil nil       output2:   value PubKeyHash
		input3:  TXID Index nil nil
		对上面的数据取哈希(即得到TXID),这个TXID就是ecdsa.Sign()的第三个参数
		第二次签名需要用到下面的东西：
		input1:  TXID Index nil nil       output1:   value PubKeyHash
		input2:  TXID Index nil PubKey    output2:   value PubKeyHash
		input3:  TXID Index nil nil
		对上面的数据取哈希(即得到TXID),这个TXID就是ecdsa.Sign()的第三个参数
		第三次签名需要用到下面的东西：
		input1:  TXID Index nil nil       output1:   value PubKeyHash
		input2:  TXID Index nil nil       output2:   value PubKeyHash
		input3:  TXID Index nil PubKey
		对上面的数据取哈希(即得到TXID),这个TXID就是ecdsa.Sign()的第三个参数

因此,签名的详细实现步骤为:
1.先找到待签名交易中input所引用的所有output的所在交易(即找prevTXS)
2.对待签名交易进行值拷贝操作,生成txcopy,txcopy中的inputs中的Sig和PubKey字段设为nil
3.将prevTXS中的output的PubKeyHash的值赋值给待签名交易的input中的PubKey字段
4.顺次改变唯一变量:input.PubKey,有多少个input就签多少次名

*/

//新创建的交易需要进行数字签名,SignTransaction()是签名的入口函数
func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	//这里要先收集所有input所引用的output的集合,因此先定义
	prevTXS := make(map[string]Transaction)
	//遍历待签名交易的input
	for _, input := range tx.TXInputs {
		//通过TXID找到所有引用的output的所在交易
		prtx, err := bc.FindTransactionByTXid(input.TXID)
		if err != nil {
			log.Panicln(err)
		}
		//收集
		prevTXS[string(input.TXID)] = prtx
	}
	//收集完后进行签名(签名详情见Sign函数)
	tx.Sign(privateKey, prevTXS)
}

//这是一个辅助函数,通过TXID找到所有所引用的output所在的交易
func (bc *BlockChain) FindTransactionByTXid(id []byte) (Transaction, error) {
	//生成一个区块迭代器
	it := bc.Iterator()
	for {
		//从后到前遍历区块
		block := it.Next()
		for _, tx := range block.Transactions {
			//比较待获取output所在的交易的TXID是否等于当前ID
			if bytes.Equal(id, tx.TXID) {
				return *tx, nil
			}
		}
		//遍历结束,如果还没有找到,证明交易不合法
		if len(block.PrevHash) == 0 {
			fmt.Println("区块链已遍历结束,但是没有找到该笔交易")
			break
		}
	}
	return Transaction{}, errors.New("无效的交易,请检查")
}

//签名函数
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXS map[string]Transaction) {
	//先拷贝(值拷贝)一份当前的待签名交易,目的是方便改变签名要变的唯一变量input.PubKey
	txcopy := tx.TrimmedCopy()
	for i, input := range txcopy.TXInputs {
		prevTX := prevTXS[string(input.TXID)]
		if len(prevTX.TXID) == 0 {
			//在FindTransactionByTXid函数中,如果找不到TXID,只返回error但是没有退出程序
			//因此在这里要Panic退出
			log.Panicln("引用的交易无效,请检查")
		}
		//将input所引用的output的PubKeyHash赋值给副本transaction的input.PubKey字段
		//这里其实就说明了为何要进行拷贝操作了,如果在原交易进行操作,会把input.PubKey给覆盖掉
		txcopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		//赋值完后直接对交易副本取哈希(这是生成临时交易ID,方便签名用的,不是整笔交易的TXID)
		//整笔交易的哈希(即交易的TXID)在签名动作发生前就已经生成好了
		txcopy.SetHash()
		//记得对input.PubKey进行置nil操作,以免影响当前交易的下一次签名
		txcopy.TXInputs[i].PubKey = nil
		//临时生成的交易哈希(TXID)就是当次签名需要的参数之一
		signDataHash := txcopy.TXID
		//用ecdsa.Sign函数签名,生成r和s
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		if err != nil {
			log.Panicln(err)
		}
		//r和s都是*big.Int类型,需要先转换为[]byte字节流
		//转换为字节流是为了在P2P网络进行传输,将签名所产生的r和s数据传输到验证端
		var signature []byte
		signature = append(signature, r.Bytes()...)
		signature = append(signature, s.Bytes()...)
		//r和s字段拼接起来就是input.Sig字段的值
		tx.TXInputs[i].Sig = signature
	}
	fmt.Println("交易签名成功")
}

//值拷贝当前待签名的交易
func (tx *Transaction) TrimmedCopy() Transaction {
	//拷贝出来的Transaction要改变两个东西,input.Sig和input.PubKey字段都需要置为nil
	//其他的不变
	var inputs []TXInput
	for _, input := range tx.TXInputs {
		inputs = append(inputs, TXInput{
			TXID:   input.TXID,
			Index:  input.Index,
			Sig:    nil,
			PubKey: nil,
		})
	}
	//output原封不动拷贝即可
	return Transaction{tx.TXID, inputs, tx.TXOutputs}
}
