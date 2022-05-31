package segment

import (
	"bytes"
	"log"
	"testing"
)

func TestMem(t *testing.T) {
	s := NewSegmentInMemory()

	s.Add(3, []Term{{"user_id", "aabbccdd"}})
	s.Add(4, []Term{{"user_id", "ffbbccdd"}})
	s.Add(5, []Term{{"user_id", "ccbbccdd"}})
	s.Add(6, []Term{{"user_id", "aabbccdd"}})
	s.Add(7, []Term{{"user_id", "aabbccdd"}})
	encoded := s.Encode()
	log.Printf(string(encoded.EncodedPointers))

	pointers, _ := PointersFromBytes(encoded.EncodedPointers)

	reader := bytes.NewReader(encoded.EncodedPostings)

	log.Printf("%v", pointers.PostingsFromBytes(encoded.EncodedPostings, Term{"user_id", "aabbccdd"}))
	p, _ := pointers.PostingsFromReader(reader, Term{"user_id", "aabbccdd"})
	log.Printf("%v", p)
}
