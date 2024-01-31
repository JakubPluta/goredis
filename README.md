# goredis
In-memory redis database implementation in Go


### Goals:
- Write parsers to handle RESP protocol messages
- Handly multiple connections with goroutines
- Use append only files to persist data


```bash
SET user jp # The key > string,  value serialized 
GET user  
    >> "jp"
```

```bash
SET user jp
*3\r\n$3\r\nset\r\n$4\r\nuser\r\n$2\r\njp
---------
*3
$3
set
$4
user
$2
jp
```

```
*3 indicates that we have an array with a size of 3. 
We will read 6 lines. 
Each pair of lines represents the type and size of the object, 
Second line contains the value of that object.

$ indicates that it is a string with a length of 5. So the next line will contain exactly 5 characters.

Similarly, when we say `GET user`, it returns the same object structure with different values
$2\r\jp\r\n

```









Ref:
RESP: https://redis.io/docs/reference/protocol-spec/
Redis serialization protocol (RESP) is the wire protocol that clients implement

To communicate with the Redis server, Redis clients use a protocol called REdis Serialization Protocol (RESP). While the protocol was designed specifically for Redis, you can use it for other client-server software projects.

RESP can serialize different data types including integers, strings, and arrays. It also features an error-specific type. A client sends a request to the Redis server as an array of strings. The array's contents are the command and its arguments that the server should execute. The server's reply type is command-specific.

RESP is binary-safe and uses prefixed length to transfer bulk data so it does not require processing bulk data transferred from one process to another.


A client connects to a Redis server by creating a TCP connection to its port (the default is 6379)
The following table summarizes the RESP data types that Redis supports:

```

| RESP data type    | Minimal protocol version | Category | First byte |
|-------------------|--------------------------|----------|------------|
| Simple strings    | RESP2                    | Simple   | +          |
| Simple Errors     | RESP2                    | Simple   | -          |
| Integers          | RESP2                    | Simple   | :          |
| Bulk strings      | RESP2                    | Aggregate| $          |
| Arrays            | RESP2                    | Aggregate| *          |
| Nulls             | RESP3                    | Simple   | _          |
| Booleans          | RESP3                    | Simple   | #          |
| Doubles           | RESP3                    | Simple   | ,          |
| Big numbers       | RESP3                    | Simple   | (          |
| Bulk errors       | RESP3                    | Aggregate| !          |
| Verbatim strings  | RESP3                    | Aggregate| =          |
| Maps              | RESP3                    | Aggregate| %          |
| Sets              | RESP3                    | Aggregate| ~          |
| Pushes            | RESP3                    | Aggregate| >          |

```