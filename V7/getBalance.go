package main

import (
	"bytes"
)

/*FindUTXOByUser函数做了以下处理:
遍历了区块链的所有区块,又对每个区块的所有交易进行了遍历,又对每个交易的所有Input进行了遍历,
如果该input的address等于用户想要获取余额的address,则创建该input所引用的output所在交易的key,
value值为该input所引用的output在其交易中的索引(筛选有关消耗)。同时遍历每个交易的所有Output,
如果某个output被某个与自己有关的input消耗过,则被过滤掉(剔除有关消耗),如果某个output的address
本身与自己无关,也被过滤掉(筛选有关未消耗,同时也排除了无关消耗和无关未消耗)。
(output中一共有四种类型:有关未消耗,有关消耗,无关未消耗,无关消耗)
*/
func (bc *BlockChain) FindUTXOByUser(address string) []TXOutput {
	publicHash160 := GetPubkeyHashFromAdderss(address)
	var UTXO []TXOutput
	var spentOutputs = make(map[string][]uint64)
	blockChain := bc.Iterator()
	for {
		//遍历所有区块
		block := blockChain.Next()
		//遍历区块的每个交易
		for _, transaction := range block.Transactions {
		OUTPUT:
			//遍历每个区块的Output
			for i, output := range transaction.TXOutputs {
				//剔除有关消耗
				if spentOutputs[string(transaction.TXID)] != nil {
					for _, spentOutputsIndex := range spentOutputs[string(transaction.TXID)] {
						if spentOutputsIndex == uint64(i) {
							continue OUTPUT
							//跳转到for i, output := range transaction.TXOutputs语句,但是i+1
						}
					}
				}
				//剔除无关消耗和无关未消耗
				//只有有关未消耗才append进UTXO中
				//注意:两个byte并不能直接用"=="来进行比较,需要用到bytes包中的Equal函数
				if bytes.Equal(output.PubKeyHash, publicHash160) {
					UTXO = append(UTXO, output)
				}
			}
			//遍历每个区块的Input
			for _, input := range transaction.TXInputs {
				//筛选有关消耗
				//注意input里面存储的是公钥,需要先将其转换为哈希才可比较
				if bytes.Equal(SetPubkeyHash(input.PubKey), publicHash160) {
					spentOutputs[string(input.TXID)] = append(spentOutputs[string(input.TXID)], uint64(input.Index))
				}
			}
		}
		//前区块哈希为空,终止遍历
		if len(block.PrevHash) == 0 {
			break
		}
	}
	//返回所有有关未消耗
	return UTXO
}
