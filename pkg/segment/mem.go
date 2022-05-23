package segment

import (
	"encoding/binary"
	"encoding/json"
)

type Term struct {
	Field string
	Value string
}

type SegmentInMemory struct {
	Postings map[string]map[string][]int32
}

func NewSegmentInMemory() *SegmentInMemory {
	return &SegmentInMemory{
		Postings: map[string]map[string][]int32{},
	}
}

func (s *SegmentInMemory) Add(did int32, terms []Term) {
	for _, t := range terms {
		f, ok := s.Postings[t.Field]
		if !ok {
			f = map[string][]int32{}
			s.Postings[t.Field] = f
		}
		f[t.Value] = append(f[t.Value], did)
	}
}

func (s *SegmentInMemory) GetPostingsList(t Term) []int32 {
	if f, ok := s.Postings[t.Field]; ok {
		if v, ok := f[t.Value]; ok {
			return v
		}
	}
	return []int32{}

}

// PROOF OF CONCEPT
// TODO: implement some custom encoding that uses delta encoded runlength encoded group varints
func (s *SegmentInMemory) Encode() EncodedSegment {
	postingsData := []byte{}

	pointers := Pointers{
		Data: map[string]map[string]Pointer{},
	}

	for fname, fvalues := range s.Postings {
		perValuePointers, ok := pointers.Data[fname]
		if !ok {
			perValuePointers = map[string]Pointer{}
			pointers.Data[fname] = perValuePointers
		}
		for vname, postings := range fvalues {
			encoded := IntArrayToByteA(postings)

			perValuePointers[vname] = Pointer{Len: len(encoded), Off: len(postingsData)}

			postingsData = append(postingsData, encoded...)
		}
	}

	return EncodedSegment{EncodedPointers: pointers.Encode(), EncodedPostings: postingsData}
}

type Pointer struct {
	Len int
	Off int
}

func (p *Pointer) PostingsFromBytes(data []byte) []int32 {
	return PostingsFromBytes(data, p.Len, p.Off)
}

type Pointers struct {
	Data map[string]map[string]Pointer
}

func (p *Pointers) PostingsFromBytes(data []byte, t Term) []int32 {
	if f, ok := p.Data[t.Field]; ok {
		if v, ok := f[t.Value]; ok {
			return v.PostingsFromBytes(data)
		}
	}
	return []int32{}
}

func (p *Pointers) Encode() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return data
}

func PointersFromBytes(data []byte) (*Pointers, error) {
	p := &Pointers{}
	err := json.Unmarshal(data, p)
	return p, err
}

func PostingsFromBytes(data []byte, length, offset int) []int32 {
	return ByteArrayToIntA(data[offset : offset+length])
}

type EncodedSegment struct {
	EncodedPointers []byte
	EncodedPostings []byte
}

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
