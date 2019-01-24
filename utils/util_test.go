package utils

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestBuildJPNotiMsg(t *testing.T) {
	str := `{
	    "id":"403",
	    "state":"pending",
	    "bizType":"DEPOSIT",
	    "coinType":"BTC",
	    "to":"1CvVvwwtVMaxvA4dLWHvrf47bkYJXCeV1j",
	    "value":"0.01000000",
	    "confirmations":3,
	    "create_at":1520325892149,
	    "update_at":1520326180664,
	    "fee":"0.00009619",
	    "data":{
	        "type":"Bitcoin",
	        "hash":"cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b",
	        "fee":0.00009619,
	        "blockNumber":512241,
	        "blockHash":"00000000000000000005675cd684528cb310de8ece0c22befb198d97a12366fa",
	        "confirmations":3,
	        "from":[
	            {
	                "address":"3QQDiUoKwNUVVnRY5Cyt5gKDhcocL7w5YP",
	                "value":"17.32975394",
	                "txid":"2a941eb498fd6235408cc2ac39456d80c33a018a66f3eb69214fc3cbf2310623",
	                "n":0
	            }
	        ],
	        "to":[
	            {
	                "address":"1CvVvwwtVMaxvA4dLWHvrf47bkYJXCeV1j",
	                "value":"0.01000000",
	                "txid":"",
	                "n":0
	            },
	            {
	                "address":"3CtstmqVNNgW2Jdj88QVtnwZnFUdXsqH8J",
	                "value":"17.31965775",
	                "txid":"",
	                "n":1
	            }
	        ]
		},
	    "hash":"cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b",
		"extraData":"",
		"sendAgain": false,
		"memo": "",
		"from": "",
		"timestamp": 1547713400684
	}`

	obj := map[string]interface{}{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(str)))
	decoder.UseNumber()
	err := decoder.Decode(&obj)
	if err != nil {
		t.Errorf("error: %v", err)
		return
	}
	msg := BuildMsg(obj)
	nodeMsg := "bizTypeDEPOSITcoinTypeBTCconfirmations3create_at1520325892149datablockHash00000000000000000005675cd684528cb310de8ece0c22befb198d97a12366fablockNumber512241confirmations3fee0.00009619from0address3QQDiUoKwNUVVnRY5Cyt5gKDhcocL7w5YPn0txid2a941eb498fd6235408cc2ac39456d80c33a018a66f3eb69214fc3cbf2310623value17.32975394hashcb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574bto0address1CvVvwwtVMaxvA4dLWHvrf47bkYJXCeV1jn0txidvalue0.010000001address3CtstmqVNNgW2Jdj88QVtnwZnFUdXsqH8Jn1txidvalue17.31965775typeBitcoinextraDatafee0.00009619fromhashcb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574bid403memosendAgainfalsestatependingtimestamp1547713400684to1CvVvwwtVMaxvA4dLWHvrf47bkYJXCeV1jupdate_at1520326180664value0.01000000"
	if msg != nodeMsg {
		t.Fail()
	}
}

func TestPriToPub(t *testing.T) {
	pri := "bf12996feeaa2977b6ca0d33a0e8bd2ccfc4844c6f8a7e6d15c099f8da4a255c"
	pub := PriToPub(pri)
	if pub != "03ace32532c90652e1bae916248e427a7ab10aeeea1067949669a3f4da10965ef9" {
		t.Fail()
	}
}
