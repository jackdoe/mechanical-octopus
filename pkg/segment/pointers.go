package segment

import "io"

//go:generate msgp
type Pointer struct {
	Len int
	Off int
}

func (p *Pointer) PostingsFromBytes(data []byte) []int32 {
	return ByteArrayToIntA(data[p.Off : p.Off+p.Len])
}

func (p *Pointer) PostingsFromReader(r io.ReaderAt) ([]int32, error) {
	data := make([]byte, p.Len)
	_, err := r.ReadAt(data, int64(p.Off))
	return ByteArrayToIntA(data[p.Off : p.Off+p.Len]), err
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

func (p *Pointers) PostingsFromReader(r io.ReaderAt, t Term) ([]int32, error) {
	if f, ok := p.Data[t.Field]; ok {
		if v, ok := f[t.Value]; ok {
			return v.PostingsFromReader(r)
		}
	}
	return []int32{}, nil
}

func (p *Pointers) Encode() []byte {
	data, err := p.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}
	return data
}

func PointersFromBytes(data []byte) (*Pointers, error) {
	p := &Pointers{}
	_, err := p.UnmarshalMsg(data)
	return p, err
}
