package internal

import (
	"errors"
	"fmt"
	"os"

	"github.com/hoanganh2070/pqdb/config"
)

type Page struct {
	Num  uint64
	Data []byte
}

type DataAccessLayer struct {
	file     *os.File
	pageSize int

	*FreeList
	*Meta
}

func NewDataAccessLayer(path string, pageSize int) (*DataAccessLayer, error) {

	dal := &DataAccessLayer{
		Meta:     NewEmptyMeta(),
		pageSize: pageSize,
	}

	_, err := os.Stat(path)
	if err == nil {
		dal.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			_ = dal.Close()
			return nil, err
		}

		meta, err := dal.ReadMeta()
		if err != nil {
			return nil, err
		}

		dal.Meta = meta

		freeList, err := dal.ReadFreelist()
		if err != nil {
			return nil, err
		}
		dal.FreeList = freeList
	} else if errors.Is(err, os.ErrNotExist) {
		dal.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			_ = dal.Close()
			return nil, err
		}

		dal.FreeList = NewFreeList()
		dal.FreeListPage = dal.GetNextPage()
		_, err := dal.WriteFreeList()
		if err != nil {
			return nil, err
		}

		dal.WriteMeta(dal.Meta)
	} else {
		return nil, err
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
		fmt.Println(err)
		return nil, err
	}

	return p, err
}

func (d *DataAccessLayer) WritePage(p *Page) error {
	offset := int64(p.Num) * int64(d.pageSize)
	_, err := d.file.WriteAt(p.Data, offset)
	return err
}

func (d *DataAccessLayer) WriteMeta(meta *Meta) (*Page, error) {
	p := d.AllocateEmptyPage()
	p.Num = config.MetaPageNum
	meta.Serialize(p.Data)

	err := d.WritePage(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *DataAccessLayer) ReadMeta() (*Meta, error) {
	p, err := d.ReadPage(config.MetaPageNum)
	if err != nil {
		return nil, err
	}

	meta := NewEmptyMeta()
	meta.Deserialize(p.Data)
	return meta, nil
}

func (d *DataAccessLayer) ReadFreelist() (*FreeList, error) {
	p, err := d.ReadPage(d.FreeListPage)
	if err != nil {
		return nil, err
	}

	freelist := NewFreeList()
	freelist.Deserialize(p.Data)
	return freelist, nil
}

func (d *DataAccessLayer) WriteFreeList() (*Page, error) {
	p := d.AllocateEmptyPage()
	p.Num = d.FreeListPage
	d.FreeList.Serialize(p.Data)

	err := d.WritePage(p)
	if err != nil {
		return nil, err
	}
	d.FreeListPage = p.Num

	return p, nil
}
