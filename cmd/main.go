package main

import (
	"os"

	"github.com/hoanganh2070/pqdb/internal"
)

func main() {
	dal, _ := internal.NewDataAccessLayer("test.db", os.Getpagesize())

	p := dal.AllocateEmptyPage()
	p.Num = dal.GetNextPage()
	copy(p.Data[:], "hello")

	_ = dal.WritePage(p)

}
