package ipfs

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

const (
	Pid      = ""
	Pse      = ""
	CID      = ""
	TmpFile  = ""
	TmpTab   = ""
	TestFile = ""
	IsTable  = true
)

func Init() {
	InitInfura(Pid, Pse)
}

func Test_Add(t *testing.T) {
	Init()

	file, err := os.Open(TestFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	cid, err := Inf.Shell.Add(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Data successfully stored in IPFS: %v\n", cid)
}

func Test_Get(t *testing.T) {
	Init()
	err := Inf.Shell.Get(CID, TmpFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Get", CID, "to", TmpFile, "OK")
}

func Test_Cat(t *testing.T) {
	Init()

	content, err := Inf.Shell.Cat(CID)
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	fmt.Println("content", buf.String())
}

func Test_AddTable(t *testing.T) {
	Init()

	cid, err := Inf.Shell.AddDir(TmpTab)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Data successfully stored in IPFS: %v\n", cid)
}
