(async () => {
    const priKey = Buffer.from('vxKZb+6qKXe2yg0zoOi9LM/EhExvin5tFcCZ+NpKJVw=', 'base64')
    const crypto = require('@jadepool/crypto')
    const moment = require('moment')
    const axios = require('axios')
    let msg = {
        "id":"3",
        "state":"pending",
        "bizType":"DEPOSIT",
        "coinType":"ETH",
        "to":"0x9014f690a2e6ae544c69ff2c3f12fa46f67e956f",
        "value":"0.01",
        "confirmations":3,
        "create_at":1520229729535,
        "update_at":1520229783025,
        "fee":"0.00021",
        "data":{
            "type":"Ethereum",
            "hash":"0x7442fd4bb80566d106671ac80461f2de96e8fb02134829532b596005505bdcde",
            "blockNumber":5199257,
            "blockHash":"0x91e1ce6a55340754deefd869439c626be64113968b08916a09c802ec69c53273",
            "fee":0.00021,
            "confirmations":3,
            "from":[
                {
                    "address":"0x0029d396902D034b3afe2A1D81D7CB9706d7D694",
                    "value":"0.01021"
                }
            ],
            "to":[
                {
                    "address":"0x9014F690a2E6aE544C69FF2C3F12fA46F67E956F",
                    "value":"0.01"
                }
            ],
            "ethereum":{
                "input":"0x",
                "gasUsed":21000,
                "gasPrice":"10000000000",
                "nonce":144
            }
        },
        "hash":"0x7442fd4bb80566d106671ac80461f2de96e8fb02134829532b596005505bdcde",
        "extraData":"",
        "sendAgain": false,
        "memo": "",
        "from": "0x0029d396902D034b3afe2A1D81D7CB9706d7D694",
        "timestamp": moment.now()
    }
    
    const sigObj = crypto.ecc.sign(msg, priKey, { hash: 'sha3', accept: 'object' })
    let data = {}
    data.crypto = 'ecc'
    data.timestamp = sigObj.timestamp
    data.sig = sigObj.signature
    data.result = msg

    let result = null
    result = await axios({
        method: 'POST',
        url: 'http://127.0.0.1:8081/api/order/noti',
        data: data,
        proxy: false
    })

    console.log('result', result.data)
})()
