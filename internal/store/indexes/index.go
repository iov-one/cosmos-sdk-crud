package indexes

import (
	"encoding/binary"
	"fmt"
	"github.com/iov-one/cosmos-sdk-crud/internal/store/types"
	"math"
)

// maxKeyLength defines the index key maximum length in bytes
const maxKeyLength = math.MaxUint16

// lengthInBytes defines the bytes required to express the length of a key
const lengthInBytes = 2

// encodeIndexKey takes a types.SecondaryKEy and encodes it
// the way in which it's encoded is the following
// key = <[1]byte=secondaryKey.ID><[2]byte=littleEndian(len(secondaryKey.Value))><[]byte(SecondaryKey.Value)>
// in this way we have keys, that when iterated, do not go over domains of longer keys which contain
// the key, example:
// keyA = <0x1,0x2,0x3>
// keyB = <0x1,0x2,0x3,0x4>
// if we wanted to iterate over index keyA we would end up in keyB domain too
// as keyB has, as prefix, keyA. In order to avoid this we encode keys using
// the strategy highlighted above.
// Error types are of types.ErrBadArgument, and happen when
// the representation of the length of the index key takes more than 2 bytes (uint16).
func encodeIndexKey(sk types.SecondaryKey) ([]byte, error) {
	length := len(sk.Value)
	if length > maxKeyLength {
		return nil, fmt.Errorf("%w: index keys bigger than %d bytes are not allowed, got: %d", types.ErrBadArgument, maxKeyLength, length)
	}
	encodedLength := make([]byte, 2)
	binary.LittleEndian.PutUint16(encodedLength, uint16(length))
	finalKey := append([]byte{sk.ID}, encodedLength...)
	return append(finalKey, sk.Value...), nil
}

// decodeIndexKey takes a key and tries to turn it into a secondary key
// error returned here are of types.ErrInternal, this is because
// since, before doing a decode, you must have done an encode
// and they should be respectively their reverse functions
// it means that either state was corrupted or there is no backwards compatibility
// between the two anymore.
func decodeIndexKey(key []byte) (sk types.SecondaryKey, err error) {
	// minimumKeyLength defines the minimum length a key has to have
	// to be converted into a secondary key
	const minimumKeyLength = lengthInBytes + 1 + 1
	// sanity checks
	length := len(key)
	if length < minimumKeyLength {
		return sk, fmt.Errorf("%w: minimum length not reached, got: %d, want: %d", types.ErrInternal, len(key), minimumKeyLength)
	}
	decodedLength := binary.LittleEndian.Uint16(key[1:3])
	valueLength := length - minimumKeyLength + 1
	if int(decodedLength) != valueLength {
		return sk, fmt.Errorf("%w, mismatch in length, decoded: %d, got: %d", types.ErrInternal, decodedLength, valueLength)
	}
	// create secondary key
	value := make([]byte, length-minimumKeyLength+1)
	copy(value, key[minimumKeyLength-1:])
	return types.SecondaryKey{
		ID:    key[0],
		Value: value,
	}, nil
}
