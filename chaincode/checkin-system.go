package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	//"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
    sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Account struct {
	ID          string `json:"ID"`
	HardwareID  string `json:"HardwareID"`
	ProjectID   string `json:"ProjectID"`
	GroupID     string `json:"GroupID"`
	IDCard      string `json:"IDCard"`
	Value1      string `json:"Value1"`
	Value2      string `json:"Value2"`
	Worktime    int    `json:"Worktime"`
	Workload    int    `json"Workload"`
	Workprice   int    `json:"Workprice"`
	SalaryPer   int    `json:"SalaryPer"` 
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function,args := APIstub.GetFunctionAndParameters()
	if function == "query" {
		return s.query(APIstub, args)
	} else if function == "init" {
		return s.initAccount(APIstub)
	} else if function == "create" {
		return s.create(APIstub, args)
	} else if function == "list" {
		return s.list(APIstub)
	} else if function == "updateWork" {
		return s.updateWork(APIstub, args)
	} else if function == "updateSalary" {
		return s.updateSalary(APIstub, args)
	} 

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	dataAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(dataAsBytes)
}

func (s *SmartContract) initAccount(APIstub shim.ChaincodeStubInterface) sc.Response {
	Accounts := []Account {
		Account{
			ID: "0000",
			HardwareID: "0000",
			ProjectID: "0000",
			GroupID: "0000",
			IDCard: "000000000000000000",
			Value1: "0",
			Value2: "0",
			Worktime: 0,
			Workload: 0,
			Workprice: 0,
			SalaryPer: 0},
	}
	fmt.Println("This is the defalt Data.")
	accountAsBytes, _ := json.Marshal(Accounts[0])
	APIstub.PutState("ACCOUNT0", accountAsBytes)
	fmt.Println("Added", Accounts[0])

	return shim.Success(nil)
}

func (s *SmartContract) create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 12{
        return shim.Error("Incorrect number of arguments. Expecting 12")
    }

    int1, _ := strconv.Atoi(args[8])
	int2, _ := strconv.Atoi(args[9])
	int3, _ := strconv.Atoi(args[10])
	int4, _ := strconv.Atoi(args[11])

    var account = Account{
    	ID: args[1],
		HardwareID: args[2],
		ProjectID: args[3],
		GroupID: args[4],
		IDCard: args[5],
		Value1: args[6],
		Value2: args[7],
		Worktime: int1,
		Workload: int2,
		Workprice: int3,
		SalaryPer: int4}
    fmt.Println("New Added:", account)
    accountAsBytes, _ := json.Marshal(account)
    fmt.Println("New args[0]", args[0])
    fmt.Println("New accountAsBytes:",accountAsBytes)
    APIstub.PutState(args[0], accountAsBytes)

    return shim.Success(nil)
}

func (s *SmartContract) list(APIstub shim.ChaincodeStubInterface) sc.Response {

    startKey := "ACCOUNT0"
    endKey := "ACCOUNT999"

    resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    // buffer is a JSON array containing QueryResults
    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        // Add a comma before array members, suppress it for the first array member
        if bArrayMemberAlreadyWritten == true {
            buffer.WriteString(",")
        }
        buffer.WriteString("{\"Key\":")
        buffer.WriteString("\"")
        buffer.WriteString(queryResponse.Key)
        buffer.WriteString("\"")

        buffer.WriteString(", \"Record\":")
        // Record is a JSON object, so we write as-is
        buffer.WriteString(string(queryResponse.Value))
        buffer.WriteString("}")
        bArrayMemberAlreadyWritten = true
    }
    buffer.WriteString("]")

    fmt.Printf("- listAllAcount:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

func (s *SmartContract) updateWork(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    fmt.Println("Account update start")
    if len(args) != 4 {
        return shim.Error("Incorrect number of arguments. Expecting 4")
    }

    accountAsBytes, _ := APIstub.GetState(args[0])
    account := Account{}

    json.Unmarshal(accountAsBytes, &account)
    time, _ := strconv.Atoi(args[1])
    load, _ := strconv.Atoi(args[2])
    price, _ := strconv.Atoi(args[3])

    account.Worktime = time
    account.Workload = load
    account.Workprice = price
    fmt.Println("Account update Worktime:",args[1])
    fmt.Println("Account update Workload:",args[2])
    fmt.Println("Account update Workprice:",args[3])

    accountAsBytes, _ = json.Marshal(account)
    APIstub.PutState(args[0], accountAsBytes)

    fmt.Println("Account update end")
    return shim.Success(nil)
}

func (s *SmartContract) updateSalary(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    fmt.Println("Account update start")
    if len(args) != 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    accountAsBytes, _ := APIstub.GetState(args[0])
    account := Account{}

    json.Unmarshal(accountAsBytes, &account)
    salary, _ := strconv.Atoi(args[1])

    account.SalaryPer = salary
    fmt.Println("Account update SalaryPer:",args[1])

    accountAsBytes, _ = json.Marshal(account)
    APIstub.PutState(args[0], accountAsBytes)

    fmt.Println("Account update end")
    return shim.Success(nil)
}

func main() {

    err := shim.Start(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating new Smart Contract: %s", err)
    }
}