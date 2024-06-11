// -----------------------------------------------------
// REDIS-X
// © Utsav Virani
// Written by: (Utsav Virani)
// -----------------------------------------------------

package main

/*
The provided Go code implements a basic RESP (REdis Serialization Protocol) parser and writer.
RESP is a protocol used by Redis for communication between clients and servers.
*/

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+' // Constant for simple string type
	ERROR   = '-' // Constant for error type
	INTEGER = ':' // Constant for integer type
	BULK    = '$' // Constant for bulk string type
	ARRAY   = '*' // Constant for array type
)

/*
Value represents a Redis protocol value with various types.

typ   - used to determine the data type of the value.
str   - holds the value of the string received from simple strings.
num   - holds the value of the integer received from integers.
bulk  - stores the string received from bulk strings.
array - holds all the values received from arrays.
*/
type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

/*
Resp struct encapsulates a buffered reader to read data from an io.Reader.
*/
type Resp struct {
	reader *bufio.Reader
}

/*
NewResp creates a new Resp instance with a buffered reader.
*/
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

/*
readLine reads a line from the buffered reader until a CRLF (\r\n) is encountered.
*/
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

/*
readInteger reads an integer value from the buffered reader.
*/
func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

/*
Read reads a value based on the type prefix from the buffered reader.
*/
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

/*
readArray reads an array of values from the buffered reader.
*/
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// Read length of the array
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// Read each value in the array
	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array = append(v.array, val)
	}

	return v, nil
}

/*
readBulk reads a bulk string value from the buffered reader.
*/
func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	// Read length of the bulk string
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// Read the bulk string
	bulk := make([]byte, len)
	r.reader.Read(bulk)
	v.bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}

/*
Driver returns the byte representation of the value based on its type.
*/
func (v Value) Driver() []byte {
	switch v.typ {
	case "array":
		return v.arrayDriver()
	case "bulk":
		return v.bulkDriver()
	case "string":
		return v.stringDriver()
	case "null":
		return v.nullDriver()
	case "error":
		return v.errorDriver()
	default:
		return []byte{}
	}
}

/*
stringDriver returns the byte representation of a simple string value.
*/
func (v Value) stringDriver() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

/*
bulkDriver returns the byte representation of a bulk string value.
*/
func (v Value) bulkDriver() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

/*
arrayDriver returns the byte representation of an array value.
*/
func (v Value) arrayDriver() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Driver()...)
	}

	return bytes
}

/*
nullDriver returns the byte representation of a null value.
*/
func (v Value) nullDriver() []byte {
	return []byte("$-1\r\n")
}

/*
errorDriver returns the byte representation of an error value.
*/
func (v Value) errorDriver() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

/*
Writter struct encapsulates an io.Writer to write data.
*/
type Writter struct {
	writter io.Writer
}

/*
NewWritter creates a new Writter instance with an io.Writer.
*/
func NewWritter(w io.Writer) *Writter {
	return &Writter{writter: w}
}

/*
Write writes the byte representation of a value to the encapsulated writer.
*/
func (w *Writter) Write(v Value) error {
	var bytes = v.Driver()

	_, err := w.writter.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
