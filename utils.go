/**
  @author:BOEN
  @data:2023/5/31
  @note:工具类函数
**/
package main

import (
	"crypto/sha1"
	"fmt"
)

//返回节点是否包含在桶的范围内
func containsID(rangeStart, rangeEnd [IDLength]byte, id [IDLength]byte) bool {
	return compareID(id, rangeStart) >= 0 && compareID(id, rangeEnd) <= 0
}

//由于采用的是随机生成节点ID 所以计算距离换成了以比较节点ID是否在桶的范围内来进行下一步的逻辑处理
func compareID(id1, id2 [IDLength]byte) int {
	for i := 0; i < IDLength; i++ {
		if id1[i] < id2[i] {
			return -1
		} else if id1[i] > id2[i] {
			return 1
		}
	}
	return 0
}

//打印每个桶中存在的NodeID
func printBucketContents(bucketName string, nodes []Node) {
	fmt.Println(bucketName, "Nodes:")
	for _, node := range nodes {
		fmt.Printf("%x\n", node.ID)
	}
}

//对字符串进行hash
func hash(data []byte) [IDLength]byte {
	// Placeholder hash function
	hash := sha1.Sum(data)
	return hash
}

//计算距离 => 异或
func xorDistance(a, b [IDLength]byte) [IDLength]byte {
	var distance [IDLength]byte
	for i := 0; i < IDLength; i++ {
		distance[i] = a[i] ^ b[i]
	}
	return distance
}

//获取桶
func getBucketIndex(distance [IDLength]byte, dht DHT) int {
	for i, bucket := range dht.Buckets {
		if containsID(bucket.RangeStart, bucket.RangeEnd, distance) {
			return i
		}
	}
	return -1 // 表示没有找到匹配的桶
}

//选择最近的两个节点
func selectClosestNodes(nodes []Node, key [IDLength]byte, count int) []Node {
	// 返回这个桶中的前两个节点
	return nodes[:count]
}

//通过ID找到peer
func findPeerByID(id [IDLength]byte, peers []*Peer) *Peer {
	for _, peer := range peers {
		if compareID(peer.ID, id) == 0 {
			return peer
		}
	}
	return nil
}

//通过bucket找到对应的DHT
func findDHTByBucket(peer *Peer, bucket *KBucket) **DHT {
	if peer.DHT.Buckets == nil {
		return nil
	}

	for _, b := range peer.DHT.Buckets {
		if &b == bucket {
			return &peer.DHT
		}
	}

	return nil
}
