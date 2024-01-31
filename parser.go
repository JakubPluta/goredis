package main

import (
	"bufio"
	"io"
	"log"
	"strconv"
)

type redisMessageType byte

const (
	// Resp 2
	BulkString   = '$'
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	Array        = '*'
)

// struct for serialization and deserialization
type RedisMessage struct {
	// typ determines the data type.
	// str holds the value of a string.
	// num holds the value of an integer.
	// bulk stores a string.
	// array holds all the values from arrays.
	typ   string
	str   string
	num   int
	bulk  string
	array []RedisMessage
}

type Resp struct {
	// Reader contains methods for reading from the buffer and storing in the RedisMessage struct.
	reader *bufio.Reader
}

// NewResp creates a new Resp instance.
//
// It takes an io.Reader as a parameter and returns a pointer to Resp.
func NewResp(r io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(r),
	}
}

// This code defines a method readLine that reads a line from the Resp reader.
// It reads one byte at a time until it encounters the \r character, which signifies the end of the line.
// Then it returns the line without the last 2 bytes (\r\n) and the number of bytes in the line.
func (resp *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := resp.reader.ReadByte() // read one byte at a time
		if err != nil {
			return nil, 0, err
		}
		n += 1                                           // increment the number of bytes
		line = append(line, b)                           // append the byte to the line
		if len(line) >= 2 && line[len(line)-2] == '\r' { // check if the line ends with \r\n if yes break the loop and return the line
			break
		}
	}
	return line[:len(line)-2], n, nil // return the line without the last 2 bytes (\r\n)
}

// readInteger reads an integer from the Resp object.
// It returns the integer read, the number of bytes read, and an error if any.
func (resp *Resp) readInteger() (int, int, error) {
	line, n, err := resp.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(i64), n, nil
}

// ReadMessage reads a message from the Resp object.
func (r *Resp) ReadMessage() (RedisMessage, error) {
	dataType, err := r.reader.ReadByte() // read first byte to determine the data type
	if err != nil {
		return RedisMessage{}, err
	}
	switch dataType { // switch on the data type
	case BulkString: // if the data type is a bulk string
		return r.readBulkString()
	case byte(Array): // if the data type is an array
		return r.readArray()
	default:
		log.Println("Invalid type: ", dataType)
		return RedisMessage{}, nil
	}
}

func (r *Resp) readArray() (RedisMessage, error) {
	// skip first byte as it's already read in ReadMessage method
	v := RedisMessage{}
	v.typ = "array"

	length, _, err := r.readInteger() // Read number of elements in the array.
	if err != nil {
		log.Println(err)
		return v, err
	}
	v.array = make([]RedisMessage, 0) // for every line in the array parse and read it

	// Iterate over the array and for each line,
	// call the ReadMessage method to parse the type according to the character at the beginning of the line.

	for i := 0; i < length; i++ {

		msg, err := r.ReadMessage()
		if err != nil {
			log.Println(err)
			return v, err
		}
		v.array = append(v.array, msg)
	}
	return v, nil
}

func (r *Resp) readBulkString() (RedisMessage, error) {
	// skip first byte as it's already read in ReadMessage method

	v := RedisMessage{}
	v.typ = "bulk"
	length, _, err := r.readInteger() // Read integer that represents number of bytes in the bulk string.
	if err != nil {
		log.Println(err)
		return v, err
	}
	// Read the bulk string, followed by the ‘\r\n’ that indicates the end of the bulk string.
	bulk := make([]byte, length)
	r.reader.Read(bulk) // Read bulk string.
	v.bulk = string(bulk)
	r.readLine() // Read \r\n

	return v, nil
}
