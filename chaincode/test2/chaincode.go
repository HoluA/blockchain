package test2

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	// database import
	"github.com/zemirco/couchdb"

)

type SimpleChaincode struct {
}

type Document struct {
	couchdb.Document
	Id string 'json:"default_id"'
	Name string 'json:"default_name"'
	Data string 'json:"default_data"'

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("test1 Init")
	_, args := stub.GetFunctionAndParameters()
	var id, name string    // Entities
	var data string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// Initialize the chaincode
 	id = args[0]
	name = args[1]
	data = args[2]
	fmt.Printf("id : %s, name : %s\n", id, name)

	// Write the state to the ledger
	err = stub.PutState(id, []byte(data))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(name, []byte(data))
	if err != nil {
		return shim.Error(err.Error())
	}

	// insert data to couchdb
	db := client.Use("test")
	doc := &Document{
		Id: id,
		Name : name,
		Data : data,
	}
	result, err := db.Post(doc)
	if err != nil {
		panic(err)
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("test1 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var id, name string    // Entities
	var info string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	id = args[0]
	name = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	idValbytes, err := stub.GetState(id)
	nameValbytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if idValbytes == nil {
		return shim.Error("Entity not found")
	}

	//info = string(idValbytes)

	// Perform the execution
	info = args[2]

	// Write the state back to the ledger
	err = stub.PutState(id, []byte(info))
	if err != nil {
		return shim.Error(err.Error())
	}

	// modify data to couchdb
	db := client.Use("test")
	doc := &Document{
                Id: id,
                Name : name,
                Data : info,
        }
        result, err := db.Put(doc)
        if err != nil {
                panic(err)
        }

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	id := args[0]
	name := args[1]
	data := stub.GetState(id)
	fmt.Printf("id : %s, name : %s\n", id, name)

	// Delete the key from the state in ledger
	err := stub.DelState(id)
	if err != nil {
		return shim.Error("Failed to delete state")
	}
	err := stub.DelState(name)
	if err != nil {
		return shim.Error("Faild to delete state")
	}

	// delete data on couchdb
	db := client.Use("test")
        doc := &Document{
                Id: id,
                Name : name,
                Data : data,
        }
	if _, err = db.Delete(doc); err != nil {
		panic(err)
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var id string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	id = args[0]

	// Get the state from the ledger
	idValbytes, err := stub.GetState(id)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	}

	if idValbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + id + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + id + "\",\"Amount\":\"" + string(idValbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(idValbytes)
}
