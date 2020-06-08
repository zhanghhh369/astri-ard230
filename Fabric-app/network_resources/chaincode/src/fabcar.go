/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract-Define the Smart Contract structure
type SmartContract struct {
}

// Header-Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Header struct {
	TransactionID string `json:"txID"`
	SenderID      string `json:"senderID"`
	ReceiverID    string `json:"receiverID"`
	SourceID      string `json:"sourceID"`
	DestinationID string `json:"destinationID"`
}

//Message.
type Message struct {
	Header  Header `json:"header"`
	Content string `json:"content"`
}

//Event-Define the msg event
type Event struct {
	Header    Header  `json:"header"`
	Message   Message `json:"message"`
	Hash      string  `json:"hash"`
	Signature string  `json:"signature"`
	File      string  `json:"file"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCar" {
		return s.queryCar(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createMsg" {
		return s.createMsg(APIstub, args)
	} else if function == "receiveMsg" {
		return s.receiveMsg(APIstub, args)
	} else if function == "queryAllCars" {
		return s.queryAllCars(APIstub)
	} else if function == "changeCarOwner" {
		return s.changeCarOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(carAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	msg := Message{
		Header: Header{
			SenderID:      "0",
			ReceiverID:    "0",
			SourceID:      "0",
			DestinationID: "0",
		},
		Content: "0",
	}

	msgAsBytes, _ := json.Marshal(msg)
	APIstub.PutState("0", msgAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) createMsg(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}
	header := Header{
		SenderID:      args[0],
		ReceiverID:    args[1],
		SourceID:      args[2],
		DestinationID: args[3],
	}

	message := Message{
		Header:  header,
		Content: args[4],
	}
	msgAsBytes, _ := json.Marshal(message)
	APIstub.PutState(args[0], msgAsBytes)

	//agent := a.CreateAgent(agentId, agentName, agentAddress, stub)
	// ==== Agent saved. Set Event ===
	message.Header.TransactionID = APIstub.GetTxID()
	header.TransactionID = APIstub.GetTxID()
	event := Event{
		Header:    header,
		Message:   message,
		Hash:      "XXXX",
		Signature: "XXXX",
		File:      args[5]}
	//eventPayload := args[0] + ";" + args[1] + ";" + args[2] + ";" + args[3] + ";" + args[4] + ";" + args[5]
	//payloadAsBytes := []byte(eventPayload)
	eventBytes, _ := json.Marshal(event)
	eventError := APIstub.SetEvent("msgEvent", eventBytes)
	//eventError := APIstub.SetEvent("MsgCreatedEvent", payloadAsBytes)
	if eventError != nil {
		return shim.Error("Failed to set event!")
	}
	return shim.Success(nil)
}
func (s *SmartContract) receiveMsg(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var encMsg Message
	err := json.Unmarshal([]byte(args[0]), &encMsg)
	if err != nil {
		shim.Error("Failed to json!")
	}
	msgBytes, _ := json.Marshal(encMsg)
	APIstub.PutState(encMsg.Header.TransactionID, msgBytes)

	//agent := a.CreateAgent(agentId, agentName, agentAddress, stub)
	// ==== Agent saved. Set Event ===

	event := "Transaction " + encMsg.Header.TransactionID + " has been accepted!"
	eventBytes, _ := json.Marshal(event)
	eventError := APIstub.SetEvent("receiveEvent", eventBytes)
	//eventError := APIstub.SetEvent("MsgCreatedEvent", payloadAsBytes)
	if eventError != nil {
		return shim.Error("Failed to set event!")
	}
	return shim.Success(nil)
}

func (s *SmartContract) queryAllCars(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CAR0"
	endKey := "CAR999"

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

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeCarOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// if len(args) != 2 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 2")
	// }

	// carAsBytes, _ := APIstub.GetState(args[0])
	// car := Car{}

	// json.Unmarshal(carAsBytes, &car)
	// car.Owner = args[1]

	// carAsBytes, _ = json.Marshal(car)
	// APIstub.PutState(args[0], carAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
