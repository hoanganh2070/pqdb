package internal

const metaPage = 0

type FreeList struct {
	maxPage       uint64
	releasedPages []uint64
}

func NewFreeList() *FreeList {
	return &FreeList{
		maxPage:       metaPage,
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
