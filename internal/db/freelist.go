package db

import (
	"encoding/binary"

	"github.com/hoanganh2070/pqdb/config"
)

type FreeList struct {
	maxPage       uint64
	releasedPages []uint64
}

func NewFreeList() *FreeList {
	return &FreeList{
		maxPage:       0,
		releasedPages: []uint64{},
	}
}

func (fr *FreeList) GetNextPage() uint64 {
	if len(fr.releasedPages) != 0 {
		pageID := fr.releasedPages[len(fr.releasedPages)-1]
		fr.releasedPages = fr.releasedPages[:len(fr.releasedPages)-1]
		return pageID

	}
	fr.maxPage += 1
	return fr.maxPage
}

func (fr *FreeList) ReleasePage(page uint64) {
	fr.releasedPages = append(fr.releasedPages, page)
}

func (fr *FreeList) Serialize(buf []byte) []byte {
	pos := 0
	binary.LittleEndian.PutUint16(buf[pos:], uint16(fr.maxPage))

	pos += 2
	binary.LittleEndian.PutUint16(buf[pos:], uint16(len(fr.releasedPages)))
	for _, page := range fr.releasedPages {
		binary.LittleEndian.PutUint64(buf[pos:], uint64(page))
		pos += config.PageNumSize
	}

	return buf

}

func (fr *FreeList) Deserialize(buf []byte) {
	pos := 0
	fr.maxPage = uint64(binary.LittleEndian.Uint16(buf[pos:]))
	pos += 2

	releasedPagesCount := int(binary.LittleEndian.Uint16(buf[pos:]))
	pos += 2
	for range releasedPagesCount {
		fr.releasedPages = append(fr.releasedPages, binary.LittleEndian.Uint64(buf[pos:]))
		pos += config.PageNumSize
	}
}
