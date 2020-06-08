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
 * The sample smart contract for documentation topicasfsdfasdf:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"log"
	"math/big"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract-Define the Smart Contract structure
type SmartContract struct {
}

// Header-Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Header struct {
	Index         string `json:"index"`
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
	X         string  `json:"x"`
	Y         string  `json:"y"`
}

//
func Hash(tag string, data []byte) string {
	h := hmac.New(sha512.New512_256, []byte(tag))
	h.Write(data)
	return toString(h.Sum(nil))
}

// NewSigningKey generates a random P-256 ECDSA private key.
func NewSigningKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return key, err
}

//sign a message
func Sign(data []byte, privkey *ecdsa.PrivateKey) ([]byte, error) {
	// hash message
	digest := sha256.Sum256(data)

	// sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privkey, digest[:])
	if err != nil {
		return nil, err
	}

	// encode the signature {R, S}
	// big.Int.Bytes() will need padding in the case of leading zero bytes
	params := privkey.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8
	rBytes, sBytes := r.Bytes(), s.Bytes()
	signature := make([]byte, curveOrderByteSize*2)
	copy(signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(signature[curveOrderByteSize*2-len(sBytes):], sBytes)

	return signature, nil
}

// Verify checks a raw ECDSA signature.
// Returns true if it's valid and false if not.
func Verify(data, signature []byte, pubkey *ecdsa.PublicKey) bool {
	// hash message
	digest := sha256.Sum256(data)

	curveOrderByteSize := pubkey.Curve.Params().P.BitLen() / 8

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:curveOrderByteSize])
	s.SetBytes(signature[curveOrderByteSize:])

	return ecdsa.Verify(pubkey, digest[:], r, s)
}

func toString(text []byte) string {
	return base64.StdEncoding.EncodeToString(text)
}

func toByte(text string) []byte {
	decodeBytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Fatalln(err)
	}
	return decodeBytes
}

//sha256
func SHA256(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return string(h.Sum(nil))
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
	} else if function == "queryAllMsg" {
		return s.queryAllMsg(APIstub)
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

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}
	header := Header{
		Index:         args[6],
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
	APIstub.PutState(header.Index, msgAsBytes)
	message.Header.TransactionID = APIstub.GetTxID()

	//encrypt message here
	passphrase := SHA256("password")
	tag := "hashing message for lookup key"
	megsHash := Hash(tag, []byte(message.Content))

	encMessage := Message{
		Header: Header{
			Index:         toString(Encrypt([]byte(args[6]), passphrase)),
			TransactionID: toString(Encrypt([]byte(APIstub.GetTxID()), passphrase)),
			SenderID:      toString(Encrypt([]byte(args[0]), passphrase)),
			ReceiverID:    toString(Encrypt([]byte(args[1]), passphrase)),
			SourceID:      toString(Encrypt([]byte(args[2]), passphrase)),
			DestinationID: toString(Encrypt([]byte(args[3]), passphrase)),
		},
		Content: toString(Encrypt([]byte(args[4]), passphrase)),
	}
	//encrypt message end

	//Sign the message
	msgAsBytes, _ = json.Marshal(encMessage)
	key, err := NewSigningKey()
	if err != nil {
		panic(err)
	}
	signature, err := Sign(msgAsBytes, key)
	if err != nil {
		panic(err)
	}

	X := key.PublicKey.X.String()
	Y := key.PublicKey.Y.String()

	//############################
	header.TransactionID = APIstub.GetTxID()
	event := Event{
		Header:    header,
		Message:   encMessage,
		Hash:      megsHash,
		Signature: toString(signature),
		File:      args[5],
		X:         X,
		Y:         Y,
	}
	//	################
	//eventPayload := args[0] + ";" + args[1] + ";" + args[2] + ";" + args[3] + ";" + args[4] + ";" + args[5]
	payloadAsBytes, _ := json.Marshal(event)
	// eventBytes := []byte(string(text1) + "!!!!!!!!!!!!!!" + text2)
	eventError := APIstub.SetEvent("msgEvent", payloadAsBytes)
	if eventError != nil {
		return shim.Error("Failed to set event!")
	}
	return shim.Success(nil)
}

func (s *SmartContract) receiveMsg(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var encEvent Event
	var decMsg Message
	tag := "hashing message for lookup key"
	passphrase := SHA256("password")
	err := json.Unmarshal([]byte(args[0]), &encEvent)
	if err != nil {
		shim.Error("Failed to json!")
	}

	X, Y := encEvent.X, encEvent.Y

	key, err := NewSigningKey()
	if err != nil {
		panic(err)
	}

	key.PublicKey.X, _ = new(big.Int).SetString(X, 10)
	key.PublicKey.Y, _ = new(big.Int).SetString(Y, 10)
	msgAsBytes, _ := json.Marshal(encEvent.Message)

	eventBytes := msgAsBytes
	if !Verify(msgAsBytes, toByte(encEvent.Signature), &key.PublicKey) {
		eventBytes, _ = json.Marshal("signature")
		eventError := APIstub.SetEvent("receiveEvent", eventBytes)
		if eventError != nil {
			return shim.Error("Failed to set event!")
		}
	} else {
		decMsg.Header.Index = string(Decrypt(toByte(encEvent.Message.Header.Index), passphrase))
		decMsg.Header.TransactionID = string(Decrypt(toByte(encEvent.Message.Header.TransactionID), passphrase))
		decMsg.Header.SenderID = string(Decrypt(toByte(encEvent.Message.Header.SenderID), passphrase))
		decMsg.Header.ReceiverID = string(Decrypt(toByte(encEvent.Message.Header.ReceiverID), passphrase))
		decMsg.Header.SourceID = string(Decrypt(toByte(encEvent.Message.Header.SourceID), passphrase))
		decMsg.Header.DestinationID = string(Decrypt(toByte(encEvent.Message.Header.DestinationID), passphrase))
		decMsg.Content = string(Decrypt(toByte(encEvent.Message.Content), passphrase))
		eventBytes, _ = json.Marshal(decMsg)
		result, _ := APIstub.GetState(encEvent.Header.Index)
		if result != nil {
			eventBytes, _ = json.Marshal("redundant")
			eventError := APIstub.SetEvent("receiveEvent", eventBytes)
			if eventError != nil {
				return shim.Error("Failed to set event!")
			}
		} else {
			if encEvent.Hash != Hash(tag, []byte(decMsg.Content)) {
				eventBytes, _ = json.Marshal("Hash")
			} else if decMsg.Content == "" {
				eventBytes, _ = json.Marshal("Content")
			} else if encEvent.Header.TransactionID != decMsg.Header.TransactionID || decMsg.Header.TransactionID == "" {
				eventBytes, _ = json.Marshal("TransactionID")
			} else if encEvent.Header.SenderID != decMsg.Header.SenderID || decMsg.Header.SenderID == "" {
				eventBytes, _ = json.Marshal("SenderID")
			} else if encEvent.Header.ReceiverID != decMsg.Header.ReceiverID || decMsg.Header.ReceiverID == "" {
				eventBytes, _ = json.Marshal("ReceiverID")
			} else if encEvent.Header.SourceID != decMsg.Header.SourceID || decMsg.Header.SourceID == "" {
				eventBytes, _ = json.Marshal("SourceID")
			} else if encEvent.Header.DestinationID != decMsg.Header.DestinationID || decMsg.Header.DestinationID == "" {
				eventBytes, _ = json.Marshal("DestinationID")
			} else {
				err2 := APIstub.PutState(encEvent.Header.Index, eventBytes)
				if err2 != nil {
					shim.Error("Failed")
				}
			}
			eventError := APIstub.SetEvent("receiveEvent", eventBytes)
			if eventError != nil {
				return shim.Error("Failed to set event!")
			}
		}
	}

	return shim.Success(nil)
}

func (s *SmartContract) queryAllMsg(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

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

	return shim.Success(buffer.Bytes())
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Encrypt(plaintext []byte, key string) (ciphertext []byte) {
	byteKey := []byte(key)
	block, err := aes.NewCipher(byteKey[:])
	if err != nil {
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil
	}

	return gcm.Seal(nonce, nonce, plaintext, nil)
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key string) (plaintext []byte) {
	byteKey := []byte(key)
	block, err := aes.NewCipher(byteKey[:])
	if err != nil {
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil
	}

	plaintext, _ = gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
	return plaintext
}

// func (s *SmartContract) encrypt(byteArray []byte, key []byte) []byte {

// 	// Create the AES cipher
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Empty array of 16 + byteArray length
// 	// Include the IV at the beginning
// 	ciphertext := make([]byte, aes.BlockSize+len(byteArray))

// 	// Slice of first 16 bytes
// 	iv := ciphertext[:aes.BlockSize]

// 	// Write 16 rand bytes to fill iv
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		panic(err)
// 	}

// 	// Return an encrypted stream
// 	stream := cipher.NewCFBEncrypter(block, iv)

// 	// Encrypt bytes from byteArray to ciphertext
// 	stream.XORKeyStream(ciphertext[aes.BlockSize:], byteArray)

// 	return ciphertext
// }

// func (s *SmartContract) decrypt(ciphertext []byte, key []byte) []byte {

// 	// Create the AES cipher
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Before even testing the decryption,
// 	// if the text is too small, then it is incorrect
// 	if len(ciphertext) < aes.BlockSize {
// 		panic("Text is too short")
// 	}

// 	// Get the 16 byte IV
// 	iv := ciphertext[:aes.BlockSize]

// 	// Remove the IV from the ciphertext
// 	ciphertext = ciphertext[aes.BlockSize:]

// 	// Return a decrypted stream
// 	stream := cipher.NewCFBDecrypter(block, iv)

// 	// Decrypt bytes from ciphertext
// 	stream.XORKeyStream(ciphertext, ciphertext)

// 	return ciphertext
// }

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
