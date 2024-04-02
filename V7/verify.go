package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"log"
	"math/big"
)

/*
区块交易验证过程：
首先验证交易的只能是矿工(全节点),轻节点不行
何时验证？  ------ 矿工要新发布一个区块的时候,需要验证打包来的交易是否合法
//注：由于目前只实现了应用层的代码,较底层的P2P网络还暂未实现(即暂未实现分布式),因此无法实现交易
      在网络中传输,只能在本地进行存储,因此缺少了某一矿工成功发布区块后其他矿工同步区块的功能

然后需要清楚的是:并不是每一笔交易只验证一次,而是一笔交易有多少个input就验证几次(这点跟签名一样)

每一次验证都需要三个参数:
1.当前input的Sig字段(但不是Sig本身,而是将Sig平分拆成的两段,即"r"和"s");
   //判断交易是不是由自己的地址所发起的
2.当前input的PubKey字段(但不是PubKey本身,而是PubKey.X和PubKey.Y)
   //判断交易所引用的是否是自己的UTXO
3.第三个参数与签名时生成的signDataHash一模一样(一样的生成步骤,一样的结果)
   //知道对哪一个input进行验证

*/

//区块交易验证的入口函数
//在新增非创世区块的时候:即AvarageBlock()函数中已循环遍历当前待验证区块的所有交易,
//所以这里传入的是单笔交易,而不是整个区块的交易
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	//先判断这笔交易是否为Coinbase交易,如果是就无需验证
	if tx.IsCoinbase() {
		fmt.Println("这笔交易是coinbase交易,无需验证")
		//记得返回
		return true
	}
	//新建一个映射,用来存储所引用的output的集合
	prevTXS := make(map[string]Transaction)
	for _, input := range tx.TXInputs {
		//遍历当前交易的input,调用bc.FindTransactionByTXid函数进行查找
		prtx, err := bc.FindTransactionByTXid(input.TXID)
		if err != nil {
			log.Panicln("无效交易，验证失败")
		}
		prevTXS[string(input.TXID)] = prtx
	}
	//至此,所引用的Output的所在交易已备好,可以开始验证
	return tx.Verify(prevTXS)
}

//验证函数
func (tx *Transaction) Verify(prevTXS map[string]Transaction) bool {
	var num int = 0
	//值拷贝一份副本,跟签名一样,目的是比较用同样方法获取到得VerifyData和SigDataHash一不一致,
	//如果不一致,说明存在input是违法的
	txcopy := tx.TrimmedCopy()
	//注意,这里与Sign只有一处不同,Sign是对副本的Inputs进行遍历,而Verify是对原交易的Inputs进行遍历
	//因为Verify要获取的Sig和PubKey在副本中没有
	for i, input := range tx.TXInputs {
		prevTx := prevTXS[string(input.TXID)]
		if len(prevTx.TXID) == 0 {
			log.Panicln("引用的交易无效,请检查")
		}
		txcopy.TXInputs[i].PubKey = prevTx.TXOutputs[input.Index].PubKeyHash
		txcopy.SetHash()
		//与Sign一样,每一次都要对input.PubKey置nil,避免对当前交易的下一个input验证时造成影响
		txcopy.TXInputs[i].PubKey = nil
		verifyData := txcopy.TXID
		PubKey := input.PubKey
		//将input.Sig字段平均拆分再通过SetBytes方法得到r和s
		r := big.Int{}
		s := big.Int{}
		r.SetBytes(input.Sig[:len(input.Sig)/2])
		s.SetBytes(input.Sig[len(input.Sig)/2:])
		//将input.PubKey字段平均拆分再通过SetBytes方法得到x和y
		x := big.Int{}
		y := big.Int{}
		x.SetBytes(PubKey[:len(PubKey)/2])
		y.SetBytes(PubKey[len(PubKey)/2:])
		//构造原生的公钥
		PublicKeyOrigin := ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     &x,
			Y:     &y,
		}
		//正式进行验证,传入参数
		res := ecdsa.Verify(&PublicKeyOrigin, verifyData, &r, &s)
		num++
		if !res {
			fmt.Printf("交易ID为%x的第%d次验证失败\n", tx.TXID, num)
			return false
		}
		fmt.Printf("交易ID为%x的第%d次验证通过\n", tx.TXID, num)
	}
	return true
}

//辅助函数: 判断待验证的交易是否为Coinbase交易
func (tx *Transaction) IsCoinbase() bool {
	//Coinbase交易的特点:input只有一个,并且没有引用的output
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXID) == 0 {
		return true
	}
	return false
}
