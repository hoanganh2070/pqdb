package main

import (
	"fmt"
	"os"

	"github.com/hoanganh2070/pqdb/internal"
)

func main() {
	dal, err := internal.NewDataAccessLayer("test.db", os.Getpagesize())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing DataAccessLayer: %v\n", err)
		return
	}

	p := dal.AllocateEmptyPage()
	p.Num = dal.GetNextPage()
	copy(p.Data[:], "hello")

	_ = dal.WritePage(p)
	dal.WriteFreeList()

	dal.Close()

	dal, _ = internal.NewDataAccessLayer("test.db", os.Getpagesize())

	p = dal.AllocateEmptyPage()
	p.Num = dal.GetNextPage()
	copy(p.Data[:], "world")

	_ = dal.WritePage(p)
	dal.WriteFreeList()

	dal.Close()

}
