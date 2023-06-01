/**
  @author:BOEN
  @data:2023/5/31
  @note:node数据结构及其函数
**/
package main

const (
	//每个桶三个节点
	K = 3
	//ID长度
	IDLength = 20
)

//简化的node;没有IP和Port
type Node struct {
	ID [IDLength]byte
}
