package main

import (
	"fmt"
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"os"
	"math/rand"
)

type RefEntity struct {
	RefEntityId   string  `json:"RefEntityId"`
	RefEntityName string  `json:"RefEntityName"`
	RefEntityCur  string  `json:"RefEntityCur"`
}

type Cds struct {
	TradeDate   string  `json:"TradeDate"`
	EffectiveDate string  `json:"EffectiveDate"`
	ProtectionSeller  string  `json:"ProtectionSeller"`
	ProtectionBuyer  string  `json:"ProtectionBuyer"`
	ReferenceEntityId  string  `json:"ReferenceEntityId"`
	CalculationAmount  string  `json:"CalculationAmount"`
	CalculationCurrency  string  `json:"CalculationCurrency"`
	MasterAgreementType  string  `json:"MasterAgreementType"`
	FixedRate  string  `json:"FixedRate"`
}

type UnMarRefEntity struct {
	RefEntityId   string
	RefEntityName string
	RefEntityCur  string
}

type ValidationResult struct {
	RefEntId  string `json:"RefEntId"`
	RefEntCur string `json:"RefEntCur"`
	RefResult string `json:"RefResult"`
}

type CDSTransactionResult struct {
	CDSTransactionDetails  string `json:"CDSTransactionDetails"`
	CDSTransResult string `json:"CDSTransResult"`
}

type CDSTransEvent struct {
	CDSTransRefEntityId string `json:"CDSTransRefEntityId"`
	CDSTransId string `json:"CDSTransId"`
}

type SimpleChaincode struct {

}
//Variable for Random String Generation
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//Logging at the Shim Level
var infoLevel,_ = shim.LogLevel("INFO")
//var debugLevel,_= shim.LogLevel("DEBUG")
//var errorLevel,_=shim.LogLevel("ERROR")
//var crticalLevel,_=shim.LogLevel("CRITICAL")

//Creating Logger Instance
var myLogger = shim.NewLogger("CreditDefaultSawpLogger")


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	myLogger.SetLevel(infoLevel)
	fmt.Println(myLogger.IsEnabledFor(infoLevel))
	myLogger.Info("Init firing. Function will be ignored: " + function)
	fmt.Println("Init firing. Function will be ignored: " + function)
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	err := stub.PutState("LoadEntities", blankBytes)

	if err != nil {
		myLogger.Error("Failed to initialize Loading Entities collection")
		fmt.Println("Failed to initialize Loading Entities collection")
	}
	fmt.Println("Initialization Complete")

	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//fmt.Println(myLogger.IsEnabledFor(infoLevel))
	//var refentid, jsonResp string
	myLogger.Info("Query firing -- " + function)
	fmt.Println("Query firing -- " + function)

	//need two args --> reference entity id and currency
	if len(args) < 1 {
		return nil, errors.New("Missing Reference Entity Arguments to query. Expecting a Refernce Entity id and currency......")
	}

	fmt.Println("Inside Query -- " + function)

	if function == "validate_RefIdAndCur" {
		fmt.Println("Inside Query -- " + function)

		return t.validate_RefIdAndCur(stub, args)

	}else if function == "retrieve_CdsTransactionDetails"{
		fmt.Println("Inside Query -- " + function)
		myLogger.Info("Inside Query -- " + function)
		return t.retrieve_CdsTransactionDetails(stub, args)
	}else {
		fmt.Println("Generic Query call")
		bytes, err := stub.GetState(args[0])
		if err != nil {
			fmt.Println("Generic Error Occured. There is nothing to query")
			return nil, errors.New("Generic Error Occured. There is nothing to query")
		}

		fmt.Println("All success, returning from generic")
		return bytes, nil
	}
}

func (t *SimpleChaincode) retrieve_CdsTransactionDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	queryCdTransId := args[0]
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments in getCdsTransactionDetails()")
	}

	cdsTransactionBytes, err := stub.GetState(queryCdTransId)
	myLogger.Info("Querying Transaction ID -- " + queryCdTransId)
	if err != nil {
		fmt.Println("Error retrieving Cds Transaction Details for " + queryCdTransId)

		result,err := json.Marshal(CDSTransactionResult{CDSTransactionDetails: queryCdTransId, CDSTransResult: "ERROR - No Trnsaction ID Found"})
		if err != nil {
			return nil, errors.New("Error Marshalling the Query Data1")
		}
		//result := string(vr)
		return result, nil
	}
	return cdsTransactionBytes,nil
}



func (t *SimpleChaincode) validate_RefIdAndCur(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var refid, refcur string
	var refEntityObj UnMarRefEntity
	var err error
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments in ValidateRefAndCur()")
	}

	refid = args[0]
	refcur = args[1]
	refBytes, err := stub.GetState(refid)
	if err != nil {
		fmt.Println("Error retrieving RefEntity " + refid)
		result,err := json.Marshal(ValidationResult{RefEntId: refid, RefEntCur: refcur, RefResult: "ERROR - No RefEntid"})
		if err != nil {
			return nil, errors.New("Error Marshalling the Query Data1")
		}
		//result := string(vr)
		return result, nil
	}
	err = json.Unmarshal(refBytes, &refEntityObj)
	fmt.Println("Passed Unmarshalling")
	if err != nil {
		fmt.Println("Error unmarshalling refbytes " + refid)
		result,err := json.Marshal(ValidationResult{RefEntId: refid, RefEntCur: refcur, RefResult: "ERROR - Failed Unmarshlling Result"})
		if err != nil {
			return nil, errors.New("Error Marshalling the Query Data2")
		}
		//result := string(vr)
		return result, nil
	}

	currencyVal := refEntityObj.RefEntityCur
	if currencyVal == refcur {
		//fmt.Println("Inside Currency Validation - Success")
		myLogger.Info("Inside Currency Validation - Success")
		result,err := json.Marshal(ValidationResult{RefEntId: refid, RefEntCur: refcur, RefResult: "SUCCESS"})
		if err != nil {
			return nil, errors.New("Error Marshalling the Query Data3")
		}
		//result := string(vr)
		return result, nil
	} else {
		fmt.Println("Inside Currency Validation - Failed")
		result,err := json.Marshal(ValidationResult{RefEntId: refid, RefEntCur: refcur, RefResult: "ERROR - Currency Code Not Matching"})
		if err != nil {
			return nil, errors.New("Error Marshalling the Query Data4")
		}
		//result := string(vr)
		return result, nil
	}

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Invoke firing -- " + function)
	if (function) == "init" {
		return t.Init(stub, "init", args)
	} else if function == "load_entities" {
		fmt.Println("Loading Entities --")
		return t.load_entities(stub, args)
	} else if function == "create_cds" {
		fmt.Println("Creating CDS --")
		return t.create_cds(stub, args)
	}

	return nil, nil
}

func (t *SimpleChaincode) create_cds(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Inside Perfrom CDS -- ")
	n := 10
	var cdsObj Cds

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments in Create_Cds")
	}

	cdsInputData := args[0]
	fmt.Println("CDS Input Data is --" + cdsInputData)
	json.Unmarshal([]byte(cdsInputData), &cdsObj)
	//tradedt = cdsObj.TradeDate
	//effectivedt = cdsObj.EffectiveDate
	//protectionseller = cdsObj.ProtectionSeller
	//protectionbuyer = cdsObj.ProtectionBuyer
	referenceid := cdsObj.ReferenceEntityId
	fmt.Println("Refernce Entity Id is" + referenceid)
	myLogger.Info("Refernce Entity Id is: " + referenceid)
	//calculationamt = cdsObj.CalculationAmount
	//calculationcur = cdsObj.CalculationCurrency
	//masteragtype = cdsObj.MasterAgreementType
	//fixedrt = cdsObj.FixedRate

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	fmt.Println("Random String is " + string(b))

	cdstransactionid :=  referenceid + "-" + string(b)
	fmt.Println("CDS Transaction Id is" + cdstransactionid)
	myLogger.Info("Refernce Transaction Id is: " + cdstransactionid)
	err := stub.PutState(cdstransactionid, []byte(cdsInputData))
	if err != nil {
		fmt.Println("Could not create CDS")
		return nil, errors.New("Not Able to Create Transaction")
	}
	var cdsevent = CDSTransEvent{referenceid,cdstransactionid}
	cdsEventBytes,err := json.Marshal(&cdsevent)
	if err != nil {
		fmt.Println("Error Marshalling CDS Event")
		return nil, errors.New("Error Markshalling CDS Event")
	}
	err = stub.SetEvent("cdsEventSender", cdsEventBytes)
	if err != nil {
		fmt.Println("Error Creating CDS Event")
		return nil, errors.New("Error Creating CDS Event")
	}
	return nil, nil

}

func (t *SimpleChaincode) load_entities(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	fmt.Println("Start loading entities --")
	/*raw, err := ioutil.ReadFile("./loadentities.json")
	if err != nil {
		fmt.Println("Erorr Loading Entities from JSON")
		return nil, errors.New("Loading Entities file should be in Parent Directory")
		os.Exit(1)
	}*/

	entitiesload := []byte(`[
				  {
				    "RefEntityId": "002BB2",
				    "RefEntityName": "Abbey National PLC",
				    "RefEntityCur": "EUR"
				  },
				  {
				    "RefEntityId": "8G836J",
				    "RefEntityName": "Tenet Healthcare Corporation",
				    "RefEntityCur": "USD"
				  },
				  {
				    "RefEntityId": "4AB951",
				    "RefEntityName": "Republic of Italy",
				    "RefEntityCur": "EUR"
				  },
				  {
				    "RefEntityId": "008FAQ",
				    "RefEntityName": "Aiful Corporation",
				    "RefEntityCur": "JPY"
				  }

				]`)

	var c []RefEntity
	json.Unmarshal(entitiesload, &c)
	entities := c
	for _, refent := range entities {
		fmt.Println("Inside Loading Entities After Unmarshall")
		bytes, err := json.Marshal(refent)
		if err != nil {
			fmt.Println("Error Loading Entity--" + refent.RefEntityId + "--EntityName--" + refent.RefEntityName)
			os.Exit(1)
		}
		err = stub.PutState(refent.RefEntityId, bytes)
		fmt.Println("Loaded Entity" + refent.RefEntityName)
	}

	fmt.Println("Finished loading entities --")
	return nil, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting Simple chaincode: %s", err)
	}
}


