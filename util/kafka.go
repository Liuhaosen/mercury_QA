package util

import (
	"encoding/json"
	logger "modtest/gostudy/lesson1/log"

	"github.com/Shopify/sarama"
)

var (
	producer sarama.SyncProducer
)

//初始化kafka
func InitKafka(addr string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //赋值为-1：这意味着producer在follower副本确认接收到数据后才算一次发送完成。
	config.Producer.Partitioner = sarama.NewRandomPartitioner //写到随机分区中，默认设置8个分区
	config.Producer.Return.Successes = true
	// msg := &sarama.ProducerMessage{}
	// msg.Topic = "mercury_topic"
	// msg.Value = sarama.StringEncoder("this is a good test")
	producer, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		logger.Error("producer close err: %v", err)
		return
	}
	logger.Debug("kafka启动成功")
	// defer producer.Close()
	return
}

//将生成的question发送到kafka队列
func SendtoKafka(topic string, value interface{}) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("json marshal failed, err:%v", err)
		return
	}
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	pid, offset, err := producer.SendMessage(msg)
	if err != nil {
		logger.Error("send message failed, err: %v", err)
		return
	}
	logger.Debug("pid:%v, offset:%v, data:%v", pid, offset, string(data))
	return
}
