package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/chaincode/test2"
	"net/url"
	"github.com/zemirco/couchdb"
)

func main() {
	err := shim.Start(new(test2.SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	u := url.Parse("http://localhost:5984/")
	client, err := couchdb.NewClient("rokmc", "8972", u)
	if err != nil {
		panic(err)
	}
}
