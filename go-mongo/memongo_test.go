package main

import (
	"fmt"
	"testing"

	"github.com/benweissmann/memongo"
)

func Test(t *testing.T) {
	mongoServer, err := memongo.Start("4.0.5")
	if err != nil {
		t.Fatal(err)
	}
	defer mongoServer.Stop()

	connectAndDoStuff(mongoServer.URI(), memongo.RandomDatabase())
}

func connectAndDoStuff(uri string, database string) {
	fmt.Println(uri)
	fmt.Println(database)
}
