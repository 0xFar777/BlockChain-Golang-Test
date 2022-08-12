package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//定义结构体
type Info struct {
	Name string
	Age  uint
}

func main() {
	//声明student1
	student1 := Info{
		Name: "小明",
		Age:  20,
	}

	//1. 编码(序列化)
	//将编码后的数据放在buffer中
	var buffer bytes.Buffer
	//func gob.NewEncoder(w io.Writer) *gob.Encoder
	//NewEncoder函数返回一个Encoder类型的指针,该指针将编码后的数据写入"w"中
	encoder := gob.NewEncoder(&buffer)
	//Encode函数将student1进行序列化,得到的字节流存放到buffer中
	err := encoder.Encode(&student1)
	if err != nil {
		log.Panic("编码失败")
	}
	fmt.Printf("编码后的student1: %v\n", buffer.Bytes())

	//2. 解码(反序列化)
	//将解码后的数据以Info结构体的形式返回
	var student2 Info
	//func gob.NewDecoder(r io.Reader) *gob.Decoder
	//NewDecoder函数返回一个Decoder类型的指针,意味着对"r"进行解码
	//首先要将buffer里的字节流数据转化为可读的形式,因此需要使用bytes.NewReader
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	//Decode函数将buffer输入流进行反序列化,得到的结果存储再student2中
	err = decoder.Decode(&student2)
	if err != nil {
		log.Panic("解码失败")
	}
	fmt.Printf("解码后的student2: %v\n", student2)
}
