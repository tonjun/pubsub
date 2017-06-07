### Architecture

```

     +------------+

     |            |-+

     | (a) Client | |

     |            | |

     +------------+ |

       +------------+

           ^

           | Load balancer

           v

     +------------+

     |            |-+

     | (b) Server | |
      
     | (pub-sub)  | |

     +------------+ |

       +------------+

           ^

           |

           v

     +------------+

     |            |

     | (c) cfgsrv |

     |            |

     +------------+


```


## Connect

```json
{
  "op": "connect",
  "type": "request",
  "id": "request_id_1",
  "auth": {
    "provider": "facebook",
    "access_token": "xxxxxx",
  },
}
```

#### server response

```json
{
  "op": "connect",
  "type": "success",
  "id": "request_id_1",
  "session": "xxxxxxxyyyyy",
  "user": {
    "id": "000000001",
    "name": "Bob",
  },
}
```
#### sample failed response

```json
{
  "op": "connect",
  "type": "error",
  "code": 401,
  "message": "Unauthorized",
}
```

## Subscribe

```json
{
  "op": "sub",
  "type": "request",
  "id": "req-123",
  "topics": [ 
    "com.myapp.chats/me", 
    "com.myapp.room1", 
    "com.myapp.room1/presence", 
  ],
}
```

```json
{
  "op": "sub",
  "type": "success",
  "id": "req-123",
  "topics": [ 
    "com.myapp.chats/me", 
    "com.myapp.room1", 
    "com.myapp.room1/presence", 
  ],
  "user": {
    "id": "000000001",
    "name": "Bob",
  },
}

```

## Unsubscribe

```json
{
  "op": "unsub",
  "type": "request",
  "id": "unsub-123",
  "topics": [ "com.myapp.room1" ],
}
```

```json
{
  "op": "unsub",
  "type": "success",
  "id": "unsub-123",
}
```


## Publish

```json
{
  "op": "pub",
  "type": "request",
  "id": "pub-456",
  "topics": [ "com.myapp.room1" ],
  "data": {
    "type": "group_chat",
    "payload": "Hi All",
  },
}
```

```json
{
  "op": "pub",
  "type": "success",
  "id": "pub-456",
}
```

## Push Message

```json
{
  "op": "mesg",
  "type": "push",
  "id": "push-abc",
  "topics": [ "com.myapp.room1" ],
  "data": {
    "type": "group_chat",
    "payload": "Hi All",
  },
  "sender": {
    "id": "000000001",
    "name": "Bob",
  },
}
```

## Presence

#### sample presence join (when someone subscribes to a topic)

```json
{
  "op": "pres",
  "type": "join",
  "id": "pres-123",
  "topics": [ "com.myapp.room1" ],
  "sender": {
    "id": "000000001",
    "name": "Bob",
  },
}
```

#### sample presence leave (when someone unsubscribes or gets disconnected)

```json
{
  "op": "pres",
  "type": "leave",
  "id": "pres-123",
  "topics": [ "com.myapp.room1" ],
  "sender": {
    "id": "000000001",
    "name": "Bob",
  },
}
```

## Resume connection

```json
{
  "op": "resume",
  "id": "abc123",
  "session": "xxxxxxyyyy",
}
```
