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
    "type": "facebook",
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
  "topics": [ "appone_chats:self", "appone_topic1", "appone_topic2" ],
}
```

```json
{
  "op": "sub",
  "type": "success",
  "id": "req-123",
  "topics": [ "appone_chats:self", "appone_topic1", "appone_topic2" ],
  "user": {
    "id": "000000001",
    "name": "Bob",
  },
}

```

## Publish

```json
{
  "op": "pub",
  "type": "request",
  "id": "pub-456",
  "topics": [ "appone_topic1" ],
  "data": {
    "type": "event",
    "payload": "profile_update",
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
  "topics": [ "appone_topic1" ],
  "data": {
    "type": "event",
    "payload": "profile_update",
  },
  "sender": {
    "id": "000000001",
    "name": "Bob",
  },
}
```
