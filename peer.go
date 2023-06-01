/**
  @author:BOEN
  @data:2023/6/1
  @note:Peer数据结构定义及其函数
**/
package main

import (
	"crypto/rand"
	"fmt"
)

// Peer数据结构
type Peer struct {
	ID [IDLength]byte
	//每个peer维护一个DHT
	DHT *DHT
}

// 创建一个Peer
func NewPeer(dht *DHT) *Peer {
	var id [IDLength]byte
	//随机生成id
	rand.Read(id[:])
	return &Peer{
		ID:  id,
		DHT: dht,
	}
}

// 在节点维护的路由表中查找节点
func (p *Peer) FindNode(nodeID [IDLength]byte) bool {
	for _, bucket := range p.DHT.Buckets {
		//如果节点id在范围内 => 返回true
		if containsID(bucket.RangeStart, bucket.RangeEnd, nodeID) {
			bucket.insertNode(p, Node{ID: nodeID})
			return true
		}
	}
	return false
}

// peer接口层 =>插入节点
func (p *Peer) insertNode(node Node) {
	//如果节点存在，直接返回
	for i := range p.DHT.Buckets {
		if containsID(p.DHT.Buckets[i].RangeStart, p.DHT.Buckets[i].RangeEnd, node.ID) {
			p.DHT.Buckets[i].insertNode(p, node)
			return
		}
	}

	// 否则为节点creat一个
	kb := KBucket{
		RangeStart: node.ID,
		RangeEnd:   node.ID,
		Nodes:      []Node{node},
	}
	p.DHT.Buckets = append(p.DHT.Buckets, kb)
}

// setvalue
func (p *Peer) SetValue(key, value [IDLength]byte, peers []*Peer) bool {
	// 1.判断hash(key)是否等于hash(value)
	if compareID(key, value) != 0 {
		return false
	}

	// 2.判断peer是否保存了这个键值对
	for _, nodes := range p.DHT.Buckets {
		//再循环bucket中的node
		for _, node := range nodes.Nodes {
			if compareID(node.ID, value) != 0 {
				// 否则保存这个键值对 => 插入该ID
				p.insertNode(Node{ID: value})
			}
		}
	}

	// 3.计算key与Peer.ID的距离
	distance := xorDistance(key, p.ID)

	// 3.计算这个节点对应的桶
	bucketIndex := getBucketIndex(distance, *p.DHT)
	if bucketIndex == -1 {
		fmt.Println("没有找到这个桶")
		return false
	}

	// 3.找到最近的两个节点
	closestNodes := p.DHT.Buckets[bucketIndex].Nodes
	selectedNodes := selectClosestNodes(closestNodes, key, 2)

	// 重复setvalue
	for _, node := range selectedNodes {
		//找到ID对应的peer
		peer := findPeerByID(node.ID, peers)
		if peer != nil {
			if peer.SetValue(key, value, peers) {
				return true
			}
		}
	}

	return false
}

// getvalue
func (p *Peer) GetValue(value [IDLength]byte, peers []*Peer) [IDLength]byte {
	// 1.判断这个value是否已经存储
	// 先循环bucket
	for _, nodes := range p.DHT.Buckets {
		//再循环bucket中的node
		for _, node := range nodes.Nodes {
			if compareID(node.ID, value) == 0 {
				return node.ID
			}
		}
	}

	// 计算距离
	distance := xorDistance(value, p.ID)

	// 计算这个节点对应的桶
	bucketIndex := getBucketIndex(distance, *p.DHT)

	// 2.找到最近的两个节点
	closestNodes := p.DHT.Buckets[bucketIndex].Nodes
	selectedNodes := selectClosestNodes(closestNodes, value, 2)

	// Invoke GetValue on the selected nodes
	for _, node := range selectedNodes {
		peer := findPeerByID(node.ID, peers)
		if peer != nil {
			if value := peer.GetValue(value, peers); value != [IDLength]byte{0} {
				return value
			}
		}
	}

	return [IDLength]byte{}
}
