/**
  @author:BOEN
  @data:2023/6/1
  @note:
**/
package main

//路由表
type DHT struct {
	Buckets []KBucket
}

//在DHT路由表中添加Bucket
func (dht *DHT) AddBucket(bucket KBucket) {
	dht.Buckets = append(dht.Buckets, bucket)
}
