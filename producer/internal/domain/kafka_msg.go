package domain

type KafkaMessage struct {
	Session string
	Data    OBUData
}
