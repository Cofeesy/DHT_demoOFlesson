/**
  @author:BOEN
  @data:2023/5/31
  @note:程序入口
**/

package main

import "math/rand"

func main() {

	/*<!--第一次实验测试-->
	--------------------
	//5个基础节点基础
	initialPeerCount := 5
	newPeerCount := 200

	peers := make([]*Peer, initialPeerCount)
	//新建5个基础peer
	for i := 0; i < initialPeerCount; i++ {
		peers[i] = NewPeer()
	}

	fmt.Println("Initial Peer IDs:")
	for _, peer := range peers {
		fmt.Printf("%x\n", peer.ID)
	}

	// 产生200个peer并加入网络
	for i := 0; i < newPeerCount; i++ {
		newPeer := NewPeer()
		randomPeer := peers[rand.Intn(len(peers))]
		randomPeer.Bucket.insertNode(Node{ID: newPeer.ID})
		peers = append(peers, newPeer)
	}

	// 打印205个节点桶中信息
	fmt.Println("Bucket Contents:")
	for _, peer := range peers {
		fmt.Printf("Peer ID: %x\n", peer.ID)
		printBucketContents("Bucket", peer.Bucket.Nodes)
		fmt.Println("--------------------")
	}*/

	//<!--第二次实验测试-->
	//------------------
	// Step 1: Initialize 100 nodes and assign them to appropriate buckets
	// Create a new DHT
	// Create a new DHT
	dht := &DHT{}

	// 初始化100个桶
	for i := 0; i < 100; i++ {
		rangeStart := [IDLength]byte{}
		rangeEnd := [IDLength]byte{}
		rand.Read(rangeStart[:])
		//直接复制
		copy(rangeEnd[:], rangeStart[:])

		bucket := KBucket{
			RangeStart: rangeStart,
			RangeEnd:   rangeEnd,
		}

		//新增桶
		dht.AddBucket(bucket)
	}

	// creat一个100个peer的数组
	peerList := make([]*Peer, 100)
	for i := 0; i < 100; i++ {
		peer := NewPeer(dht)
		peerList[i] = peer
	}

	// Step 2: Randomly generate 200 keys and select a random node to perform SetValue operation
	keys := [][IDLength]byte{}
	for i := 0; i < 200; i++ {
		var key [IDLength]byte
		rand.Read(key[:])
		keys = append(keys, key)

		// Randomly select a node from the peerList to perform SetValue operation
		nodeIndex := rand.Intn(len(peerList))
		node := peerList[nodeIndex]
		node.SetValue(key, key, peerList)
	}

	// Step 3: Randomly select 100 keys and perform GetValue operation on a random node
	for i := 0; i < 100; i++ {
		// Randomly select a key from the keys array
		keyIndex := rand.Intn(len(keys))
		key := keys[keyIndex]

		// Randomly select a node from the peerList to perform GetValue operation
		nodeIndex := rand.Intn(len(peerList))
		node := peerList[nodeIndex]
		node.GetValue(key, peerList)
	}
}
