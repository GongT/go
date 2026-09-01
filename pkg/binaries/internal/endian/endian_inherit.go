package endian

func (s ByteOrder) Uint16(b []byte) uint16       { return s.impl.Uint16(b) }
func (s ByteOrder) Uint32(b []byte) uint32       { return s.impl.Uint32(b) }
func (s ByteOrder) Uint64(b []byte) uint64       { return s.impl.Uint64(b) }
func (s ByteOrder) PutUint16(b []byte, v uint16) { s.impl.PutUint16(b, v) }
func (s ByteOrder) PutUint32(b []byte, v uint32) { s.impl.PutUint32(b, v) }
func (s ByteOrder) PutUint64(b []byte, v uint64) { s.impl.PutUint64(b, v) }
func (s ByteOrder) String() string               { return s.impl.String() }

func (s ByteOrder) AppendUint16(b []byte, v uint16) []byte {
	return s.impl.AppendUint16(b, v)
}
func (s ByteOrder) AppendUint32(b []byte, v uint32) []byte {
	return s.impl.AppendUint32(b, v)
}
func (s ByteOrder) AppendUint64(b []byte, v uint64) []byte {
	return s.impl.AppendUint64(b, v)
}
