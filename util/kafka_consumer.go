package util

import (
	logger "modtest/gostudy/lesson1/log"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

var (
	wg sync.WaitGroup
)

//初始化kafka消费者
func InitKafkaConsumer(addr, topic string, consume func(message *sarama.ConsumerMessage)) (err error) {
	consumer, err := sarama.NewConsumer(strings.Split(addr, ","), nil)
	if err != nil {
		logger.Error("Failed to start consumer: %v", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) //获得该topic所有的分区
	if err != nil {
		logger.Error("Failed to get the list of partition:, %v", err)
		return
	}

	for partition := range partitionList {
		pc, errRet := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if errRet != nil {
			logger.Error("Failed to start consumer for partition %d: %s\n", partition, errRet)
			return
		}
		wg.Add(1)
		go func(sarama.PartitionConsumer) { //为每个分区开一个go协程去取值
			for msg := range pc.Messages() { //阻塞直到有值发送过来，然后再继续等待
				logger.Debug("Partition:%d, Offset:%d, key:%s, value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				consume(msg)
			}
			// defer pc.AsyncClose()
			wg.Done()
		}(pc)
	}
	// wg.Wait()
	// consumer.Close()
	return
}
