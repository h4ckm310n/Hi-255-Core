package utils

type MessageItem struct {
	MessageType int32
	Timestamp   int64
	DeviceID    string
	Content     string
}

var MessageQueue []*MessageItem

func MessageEnqueue(messageType int32, timestamp int64, deviceID string, content string) {
	MessageQueue = append(MessageQueue, &MessageItem{
		MessageType: messageType,
		Timestamp:   timestamp,
		DeviceID:    deviceID,
		Content:     content,
	})
}

func MessageDequeue() *MessageItem {
	if len(MessageQueue) == 0 {
		return nil
	}
	item := MessageQueue[0]
	MessageQueue = MessageQueue[1:]
	return item
}
