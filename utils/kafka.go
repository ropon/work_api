package utils

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ropon/work_api/conf"
)

func KafkaSet(topic, message string) (pid int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Key = sarama.StringEncoder("testkey")
	msg.Value = sarama.StringEncoder(message)
	pid, offset, err = conf.KafkaProducer.SendMessage(msg)
	return
}

func KafkaGet(topic string) error {
	partitionList, err := conf.KafkaConsumer.Partitions(topic)
	if err != nil {
		return err
	}
	for partition := range partitionList {
		pc, err := conf.KafkaConsumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			return err
		}
		defer pc.AsyncClose()
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
	return nil
}
