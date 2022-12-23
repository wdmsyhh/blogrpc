package bson

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"os"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type M = primitive.M

type E = primitive.E

type A = primitive.A

type D = primitive.D

type RegEx = primitive.Regex
type Raw = bson.RawValue
type ObjectId string

var NilObjectId = ObjectIdHex("000000000000000000000000")

var ErrNotFound = qmgo.ErrNoSuchDocuments

// ObjectIdHex returns an ObjectId from the provided hex representation.
// Calling this function with an invalid hex representation will
// cause a runtime panic. See the IsObjectIdHex function.
func ObjectIdHex(s string) ObjectId {
	d, err := hex.DecodeString(s)
	if err != nil || len(d) != 12 {
		panic(fmt.Sprintf("invalid input to ObjectIdHex: %q", s))
	}
	return ObjectId(d)
}

// IsObjectIdHex returns whether s is a valid hex representation of
// an ObjectId. See the ObjectIdHex function.
func IsObjectIdHex(s string) bool {
	if len(s) != 24 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

// objectIdCounter is atomically incremented when generating a new ObjectId
// using NewObjectId() function. It's used as a counter part of an id.
var objectIdCounter = readRandomUint32()

// readRandomUint32 returns a random objectIdCounter.
func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot read random object id: %v", err))
	}
	return uint32((uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24))
}

// machineId stores machine id generated once and used in subsequent calls
// to NewObjectId function.
var machineId = readMachineId()
var processId = os.Getpid()

// readMachineId generates and returns a machine id.
// If this function fails to get the hostname it will cause a runtime error.
func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

// NewObjectId returns a new unique ObjectId.
func NewObjectId() ObjectId {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	b[7] = byte(processId >> 8)
	b[8] = byte(processId)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return ObjectId(b[:])
}

// NewObjectIdWithTime returns a dummy ObjectId with the timestamp part filled
// with the provided number of seconds from epoch UTC, and all other parts
// filled with zeroes. It's not safe to insert a document with an id generated
// by this method, it is useful only for queries to find documents with ids
// generated before or after the specified timestamp.
func NewObjectIdWithTime(t time.Time) ObjectId {
	var b [12]byte
	binary.BigEndian.PutUint32(b[:4], uint32(t.Unix()))
	return ObjectId(string(b[:]))
}

// String returns a hex string representation of the id.
// Example: ObjectIdHex("4d88e15b60f486e428412dc9").
func (id ObjectId) String() string {
	return fmt.Sprintf(`ObjectIdHex("%x")`, string(id))
}

// Hex returns a hex representation of the ObjectId.
func (id ObjectId) Hex() string {
	return hex.EncodeToString([]byte(id))
}

// MarshalJSON turns a bson.ObjectId into a json.Marshaller.
func (id ObjectId) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%x"`, string(id))), nil
}

// MarshalText turns bson.ObjectId into an encoding.TextMarshaler.
func (id ObjectId) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%x", string(id))), nil
}

// UnmarshalText turns *bson.ObjectId into an encoding.TextUnmarshaler.
func (id *ObjectId) UnmarshalText(data []byte) error {
	if len(data) == 1 && data[0] == ' ' || len(data) == 0 {
		*id = ""
		return nil
	}
	if len(data) != 24 {
		return fmt.Errorf("invalid ObjectId: %s", data)
	}
	var buf [12]byte
	_, err := hex.Decode(buf[:], data[:])
	if err != nil {
		return fmt.Errorf("invalid ObjectId: %s (%s)", data, err)
	}
	*id = ObjectId(string(buf[:]))
	return nil
}

// Valid returns true if id is valid. A valid id must contain exactly 12 bytes.
func (id ObjectId) Valid() bool {
	return len(id) == 12 && id != NilObjectId
}

func (id ObjectId) IsZero() bool {
	return !id.Valid()
}

// byteSlice returns byte slice of id from start to end.
// Calling this function with an invalid id will cause a runtime panic.
func (id ObjectId) byteSlice(start, end int) []byte {
	if len(id) != 12 {
		panic(fmt.Sprintf("invalid ObjectId: %q", string(id)))
	}
	return []byte(string(id)[start:end])
}

// Time returns the timestamp part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ObjectId) Time() time.Time {
	// First 4 bytes of ObjectId is 32-bit big-endian seconds from epoch.
	secs := int64(binary.BigEndian.Uint32(id.byteSlice(0, 4)))
	return time.Unix(secs, 0)
}

func (id ObjectId) Timestamp() time.Time {
	return id.Time()
}

// Machine returns the 3-byte machine id part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ObjectId) Machine() []byte {
	return id.byteSlice(4, 7)
}

// Pid returns the process id part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ObjectId) Pid() uint16 {
	return binary.BigEndian.Uint16(id.byteSlice(7, 9))
}

// Counter returns the incrementing value part of the id.
// It's a runtime error to call this method with an invalid id.
func (id ObjectId) Counter() int32 {
	b := id.byteSlice(9, 12)
	// Counter is stored as big-endian 3-byte value
	return int32(uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]))
}

func Marshal(in interface{}) (out []byte, err error) {
	return bson.MarshalWithRegistry(DefaultRegistry, in)
}

func MarshalJSON(value interface{}) ([]byte, error) {
	return bson.MarshalExtJSONWithRegistry(DefaultRegistry, value, false, false)
}

func MarshalValue(value interface{}) (bsontype.Type, []byte, error) {
	return bson.MarshalValueWithRegistry(DefaultRegistry, value)
}

func Unmarshal(in []byte, out interface{}) (err error) {
	return bson.UnmarshalWithRegistry(DefaultRegistry, in, out)
}

func UnmarshalJSON(data []byte, value interface{}) error {
	return bson.UnmarshalExtJSONWithRegistry(DefaultRegistry, data, false, value)
}

func GetBulkUpsertIds(result *mongo.BulkWriteResult) []ObjectId {
	if result == nil {
		return nil
	}
	ids := make([]ObjectId, len(result.UpsertedIDs))
	index := 0
	for _, v := range result.UpsertedIDs {
		if oId, ok := v.(primitive.ObjectID); ok {
			ids[index] = ObjectIdHex(oId.Hex())
		}
		index++
	}
	return ids
}

func GetInsertManyIds(result *qmgo.InsertManyResult) []ObjectId {
	if result == nil {
		return []ObjectId{}
	}
	ids := make([]ObjectId, len(result.InsertedIDs))
	index := 0
	for _, v := range result.InsertedIDs {
		ids[index] = v.(ObjectId)
		index++
	}
	return ids
}

func GetInsertOneId(result *qmgo.InsertOneResult) ObjectId {
	return result.InsertedID.(ObjectId)
}

func IsDup(bulkWriteErr *mongo.BulkWriteException) bool {
	if bulkWriteErr != nil {
		if len(bulkWriteErr.WriteErrors) == 0 {
			return false
		}
		for _, writeErr := range bulkWriteErr.WriteErrors {
			if qmgo.IsDup(writeErr.WriteError) {
				return true
			}
		}
	}
	return false
}
