package internal

import (
	"fmt"
	"os"
)

type Page struct {
	Num  uint64
	Data []byte
}

type DataAccessLayer struct {
	file     *os.File
	pageSize int

	*FreeList
}

func NewDataAccessLayer(path string, pageSize int) (*DataAccessLayer, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	dal := &DataAccessLayer{
		file:     file,
		pageSize: pageSize,
		FreeList: NewFreeList(),
	}
	return dal, nil

}

func (d *DataAccessLayer) Close() error {
	if d.file != nil {
		err := d.file.Close()
		if err != nil {
			return fmt.Errorf("cound not close file: %s", err)
		}
		d.file = nil
	}
	return nil
}

func (d *DataAccessLayer) AllocateEmptyPage() *Page {
	return &Page{
		Data: make([]byte, d.pageSize),
	}
}

func (d *DataAccessLayer) ReadPage(pageNum uint64) (*Page, error) {
	p := d.AllocateEmptyPage()

	offset := int(pageNum) * d.pageSize

	_, err := d.file.ReadAt(p.Data, int64(offset))

	if err != nil {
		return nil, err
	}

	return p, err
}

func (d *DataAccessLayer) WritePage(p *Page) error {
	offset := int64(p.Num) * int64(d.pageSize)
	_, err := d.file.WriteAt(p.Data, offset)
	return err
}
