package db

import (
	"encoding/binary"

	"github.com/hoanganh2070/pqdb/config"
)

type Meta struct {
	FreeListPage uint64
}

func NewEmptyMeta() *Meta {
	return &Meta{}
}

func (m *Meta) Serialize(buf []byte) {
	pos := 0
	binary.LittleEndian.PutUint64(buf[pos:], uint64(m.FreeListPage))
	pos += config.PageNumSize
}

func (m *Meta) Deserialize(buf []byte) {
	pos := 0
	m.FreeListPage = uint64(binary.LittleEndian.Uint64(buf[pos:]))
	pos += config.PageNumSize
}
