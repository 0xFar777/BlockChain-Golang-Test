package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	//打开bolt数据库
	//打开testBolt.db文件,如果文件不存在则自动创建
	//0600表示给文件操作者赋予读写权限
	db, err := bolt.Open("testBolt.db", 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败")
	}
	//所有操作完成后关闭该文件
	defer db.Close()

	//Bolt数据库写入操作
	//func (*bolt.DB).Update(fn func(*bolt.Tx) error) error
	db.Update(func(tx *bolt.Tx) error {
		//打开数据库中的一个抽屉(Bucket)
		bucket := tx.Bucket([]byte("b1"))
		//如果该抽屉不存在,则创建该抽屉
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte("b1"))
			//创建抽屉可能不成功
			if err != nil {
				log.Panic("创建bucket(b1)失败")
			}
		}
		//如果该抽屉存在或该抽屉刚被创建成功,则往该抽屉写入数据
		//Bolt数据库采用key-value写入数据
		bucket.Put([]byte("1"), []byte("HelloWorld"))
		bucket.Put([]byte("2"), []byte("BlockChain"))
		return nil
	})

	//Bolt数据库读取操作
	//func (*bolt.DB).View(fn func(*bolt.Tx) error) error
	db.View(func(tx *bolt.Tx) error {
		//打开数据库中的一个抽屉(Bucket)
		bucket := tx.Bucket([]byte("b1"))
		//如果该抽屉不存在,则报错
		if bucket == nil {
			log.Panic("该抽屉不存在")
		}
		//抽屉存在,则根据key读取value
		value1 := bucket.Get([]byte("1"))
		value2 := bucket.Get([]byte("2"))
		fmt.Printf("value1: %s\n", value1)
		fmt.Printf("value2: %s\n", value2)
		return nil
	})
}
