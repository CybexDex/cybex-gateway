(async () => {
    const priKey = Buffer.from('vxKZb+6qKXe2yg0zoOi9LM/EhExvin5tFcCZ+NpKJVw=', 'base64')
    const crypto = require('@jadepool/crypto')
    const moment = require('moment')
    const axios = require('axios')
    let msg = {
        "id":"404",
        "state":"done",
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
                    "n":1
                },
                {
                    "address":"3CtstmqVNNgW2Jdj88QVtnwZnFUdXsqH8J",
                    "value":"17.31965775",
                    "txid":"",
                    "n":0
                }
            ]
      },
      "hash":"cb51b5174b1059549be8b54cd9a8710f510889a465da28fe590c43a38052574b",
      "extraData":"",
      "sendAgain": false,
      "memo": "",
      "from": "3QQDiUoKwNUVVnRY5Cyt5gKDhcocL7w5YP",
      "timestamp": moment.now()
    }

    const sigObj = crypto.ecc.sign(msg, priKey, { hash: 'sha3', accept: 'object' })
    let data = {}
    data.crypto = 'ecc'
    data.timestamp = sigObj.timestamp
    data.sig = sigObj.signature
    data.result = msg

    let result = null
    try {
        result = await axios({
            method: 'POST',
            url: 'http://127.0.0.1:8081/api/order/noti',
            data: data,
            proxy: false
        })
    } catch (err) {
        console.error(err)
        return
    }

    console.log('result', result.data)
})()
