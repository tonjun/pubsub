package pubsub

import (
	"encoding/json"
)

type Message struct {
	OP     string      `json:"op"`
	ID     string      `json:"id"`
	Topics []string    `json:"topics,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Sender interface{} `json:"data,omitempty"`
}

const (
	OPSubscribe = "subscribe"
)

func ToBytes(m *Message) []byte {
	b, err := json.Marshal(m)
	if err != nil {
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
  "op": "connect-response",
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
  "op": "publish-response",
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
