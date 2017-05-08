package pubsub

import (
	"encoding/json"
	"log"
)

type Message struct {
	OP     string      `json:"op"`
	ID     string      `json:"id"`
	Topics []string    `json:"topics,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Sender interface{} `json:"sender,omitempty"`
}

const (
	OPSubscribe         = "subscribe"
	OPSubscribeResponse = "subscribe_response"
	OPPublish           = "publish"
	OPPublishResponse   = "publish_response"
)

func ToBytes(m *Message) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("ToBytes: error: %s", err.Error())
		return []byte("")
	}
	return b
}

/*
{
  "op": "connect",
  "id": "req123",
}
{
  "op": "connect_response",
  "id": "req123",
  "data": {
    "type": "success",
	"connection_id": "123abc",
  },
}
*/

/*
{
  "op": "subscribe",
  "id": "reqid1",
  "topics": [ "topic1", "topic2" ],
}
{
  "op": "subscribe_response",
  "id": "reqid1",
  "data": {
    "type": "success",
  },
}

*/

/*
{
  "op": "publish",
  "id": "reqid1",
  "topics": [ "topic1" ],
  "data": "abcdef",
  "sender": {
    "name": "Bob"
  },
}
{
  "op": "publish_response",
  "id": "reqid1",
  "data": {
    "type": "success",
  },
}

*/

/*
{
  "op": "join",
  "id": "xxxx",
  "topics": [ "topic1" ],
  "sender": {
    "id": "bob123",
    "name": "Bob",
  }
}
*/

/*
{
  "op": "leave",
  "id": "xxxx",
  "topics": [ "topic1" ],
  "sender": {
    "id": "bob123",
    "name": "Bob",
  }
}
*/
