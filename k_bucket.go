/**
  @author:BOEN
  @data:2023/5/30
  @note:
**/
package main

import (
	"fmt"
)

// K_bucket数据结构
type KBucket struct {
	//桶的初始范围 => 起始ID => 每个桶负责维护一定范围的节点
	RangeStart [IDLength]byte
	//桶的结束范围 => 结束ID
	RangeEnd [IDLength]byte
	Nodes    []Node
}

//插入节点
func (kb *KBucket) insertNode(peer *Peer, node Node) {
	//如果桶未满 => 插入
	if len(kb.Nodes) < K {
		kb.Nodes = append(kb.Nodes, node)
	} else {
		//否则 = >判断节点是否再桶的范围内
		//如果在 => 分裂
		if containsID(kb.RangeStart, kb.RangeEnd, node.ID) {
			kb.splitAndRedistribute(peer, node)
		} else {
			// 否则插入该节点
			fmt.Println("Node discarded:", node.ID)
		}
	}
}

// 分裂模拟
func (kb *KBucket) splitAndRedistribute(peer *Peer, node Node) {
	newRangeEnd := kb.RangeEnd
	copy(newRangeEnd[:], node.ID[:])

	// Create a new bucket with the updated range
	newBucket := KBucket{
		RangeStart: newRangeEnd,
		RangeEnd:   kb.RangeEnd,
	}

	// Update the current bucket's range
	kb.RangeEnd = newRangeEnd

	// Reallocate nodes between the two buckets
	for i := len(kb.Nodes) - 1; i >= 0; i-- {
		if containsID(kb.RangeStart, kb.RangeEnd, kb.Nodes[i].ID) {
			newBucket.Nodes = append(newBucket.Nodes, kb.Nodes[i])
			kb.Nodes = append(kb.Nodes[:i], kb.Nodes[i+1:]...)
		}
	}

	// Add the new node to the appropriate bucket
	if containsID(kb.RangeStart, kb.RangeEnd, node.ID) {
		kb.Nodes = append(kb.Nodes, node)
	} else {
		newBucket.Nodes = append(newBucket.Nodes, node)
	}

	// Add the new bucket to the DHT
	dht := findDHTByBucket(peer, kb)
	if dht != nil {
		(*dht).AddBucket(newBucket)
	}

	// Print the new bucket ranges and node allocation
	fmt.Println("Bucket 1 Range:", kb.RangeStart, "-", kb.RangeEnd)
	fmt.Println("Bucket 2 Range:", newBucket.RangeStart, "-", newBucket.RangeEnd)
	fmt.Println("Node Allocation:")
	printBucketContents("Bucket 1", kb.Nodes)
	printBucketContents("Bucket 2", newBucket.Nodes)
}
