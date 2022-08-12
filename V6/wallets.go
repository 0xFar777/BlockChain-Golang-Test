package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet.data"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

//加载已有的钱包(从"wallet.data"文件中读取出所有的地址,并赋值到map映射中)
func NewWallets() *Wallets {
	//创建一个从钱包地址到钱包公私钥的映射
	var wallets Wallets
	wallets.WalletsMap = make(map[string]*Wallet)
	//这个LoadFile函数提供了将"wallet.data"中保存的所有地址取出并赋值给map映射的功能
	wallets.LoadFile()
	return &wallets
}

//创建一个新地址,并将其保存到map映射和"wallet.data"文件中
//map映射是临时存储,把可执行文件(.exe)关了就不复存在了,而"wallet.data"文件是持久化存储
func (ws *Wallets) CreateWallet() string {
	//创建一对公私钥
	wallet := NewWallet()
	//从公钥生成地址
	address := wallet.CreateAddress()
	//保存到本地map映射
	ws.WalletsMap[address] = wallet
	//保存到"wallet.data"中以实现持久化存储
	ws.SaveToFile()
	return address
}

//将新创建的钱包地址保存到本地"wallet.data"文件中
func (ws *Wallets) SaveToFile() {
	//存储文件前，需要将待存储内容编码成字节流
	var buffer bytes.Buffer
	//由于存在interface类型,因此需要向gob进行注册。
	//panic: gob: type not registered for interface: elliptic.p256Curve
	gob.Register(elliptic.P256())
	//定义gob编码器
	encoder := gob.NewEncoder(&buffer)
	//对内容进行编码(转成字节流)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panicln("文件编码出错")
		return
	}
	//将内容(字节流形式)写入到"wallet.data"文件中
	err = ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
	if err != nil {
		log.Panicln("写入钱包文件出错")
		return
	}
}

//LoadFile函数提供了将"wallet.data"中保存的所有地址取出并赋值给map映射的功能
func (ws *Wallets) LoadFile() {
	//判断文件状态,看是否是空文件
	_, err := os.Stat(walletFile)
	//如果是空文件,直接返回即可
	//注意：不能用log.Panic抛出异常,因为第一次创建文件时就是空的,里面没有存储
	//     任何地址,直接返回即可
	if os.IsNotExist(err) {
		return
	}
	//读取"wallet.data"内容(此时还是字节流形式,需要转换类型)
	contentByte, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panicln("读取文件出错")
		return
	}
	//同理,存在interface类型,需要向gob进行注册
	gob.Register(elliptic.P256())
	//定义gob解码器
	decoder := gob.NewDecoder(bytes.NewReader(contentByte))
	//定义一个Wallets类型的结构体,字节流要转化成该形式
	walletInfo := Wallets{}
	//对字节流进行解码
	decoder.Decode(&walletInfo)
	if err != nil {
		log.Panicln("解码钱包文件内容失败")
		return
	}
	//将"wallet.data"中保存的所有地址取出并赋值给map映射
	ws.WalletsMap = walletInfo.WalletsMap
}

//将map映射中的所有地址全部读取出来
func (ws *Wallets) ListAddress() []string {
	var addresses []string
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}
	return addresses
}
