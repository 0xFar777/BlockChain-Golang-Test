package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

//定义钱包结构体
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

//创建钱包
func NewWallet() *Wallet {
	//创建椭圆曲线
	curve := elliptic.P256()
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panicln("生成私钥出错")
	}
	//生成公钥
	publicKeyOrigin := privateKey.PublicKey
	//将公钥转化成[]byte形式
	var publicKey []byte
	publicKey = append(publicKey, publicKeyOrigin.X.Bytes()...)
	publicKey = append(publicKey, publicKeyOrigin.Y.Bytes()...)
	return &Wallet{PrivateKey: privateKey, PublicKey: publicKey}
}

//生成地址
/*
地址生成过程：
对公钥取哈希(256位) ---> 生成ripemd160编码器 ---> 将公钥哈希传入ripemd160编码器 --->
获取160位的公钥哈希 ---> 生成21位payload(由version和160位公钥哈希组成) --->
对21位payload取二次哈希 ---> 生成25位payload(version + 160公钥哈希 + 21位payloade的二次哈希) --->
对25位payload进行base58编码 ---> 生成地址
*/
func (w *Wallet) CreateAddress() string {
	publicKey := w.PublicKey
	//获取160位的公钥哈希
	PublicKeyHash160 := SetPubkeyHash(publicKey)
	//定义区块链版本号
	version := byte(00)
	//生成21位payload(由version和160位公钥哈希组成)
	payload := append([]byte{version}, PublicKeyHash160...)
	//对21位payload取二次哈希
	tmp := sha256.Sum256(payload)
	payloadHash := sha256.Sum256(tmp[:])
	checkCode := payloadHash[:4]
	//生成25位payload(version + 160公钥哈希 + 21位payloade的二次哈希)
	payload = append(payload, checkCode...)
	//对25位payload进行base58编码
	address := base58.Encode(payload)
	return address
}

func SetPubkeyHash(Pubkey []byte) []byte {
	//对公钥取哈希(256位)
	PublicKeyHash256 := sha256.Sum256(Pubkey)
	//生成ripemd160编码器
	rip160hasher := ripemd160.New()
	//将公钥哈希传入ripemd160编码器
	_, err := rip160hasher.Write(PublicKeyHash256[:])
	if err != nil {
		log.Panicln(err)
	}
	//获取160位的公钥哈希
	PublicKeyHash160 := rip160hasher.Sum(nil)
	return PublicKeyHash160
}

//此函数提供了将公钥哈希转化为地址的功能,方便遍历区块时可以返回出块矿工的地址
func (output *TXOutput) SetAddress() string {
	PubKeyHash := output.PubKeyHash
	//定义区块链版本号
	version := byte(00)
	//生成21位payload(由version和160位公钥哈希组成)
	payload := append([]byte{version}, PubKeyHash...)
	//对21位payload取二次哈希
	tmp := sha256.Sum256(payload)
	payloadHash := sha256.Sum256(tmp[:])
	checkCode := payloadHash[:4]
	//生成25位payload(version + 160公钥哈希 + 21位payloade的二次哈希)
	payload = append(payload, checkCode...)
	//对25位payload进行base58编码
	address := base58.Encode(payload)
	return address
}
