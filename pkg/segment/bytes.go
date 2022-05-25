package segment

import "encoding/binary"

func ByteArrayToIntA(data []byte) []int32 {
	postings := make([]int32, len(data)/4)
	for i := 0; i < len(postings); i++ {
		from := i * 4
		postings[i] = int32(binary.LittleEndian.Uint32(data[from : from+4]))
	}
	return postings
}
func IntArrayToByteA(data []int32) []byte {
	b := make([]byte, 4*len(data))
	for i, did := range data {
		binary.LittleEndian.PutUint32(b[i*4:], uint32(did))
	}
	return b
}
