const config = require('./config')
const axios = require('axios')
const url = 'http://127.0.0.1:8081/api/order/noti';
const pg = require('pg')
const pool = new pg.Pool(config.db)

const emptyDb = async () => {
  const client = await pool.connect()
  let now = Math.round(new Date().getTime()/1000)
  let sql = `DELETE FROM jp_orders`
  // console.log('---------------------------------')
  // console.log(sql)
  let ret
  try {
    ret = await client.query(sql)
    if (ret.rowCount === 1) return true
    else return false
  } catch (err) {
    console.log(err)
  } finally {
    client.release()
  }
}

const fireDeposit = async (status) => {
  console.log(status, ' - ',new Date().getTime())
  const priKey = Buffer.from('vxKZb+6qKXe2yg0zoOi9LM/EhExvin5tFcCZ+NpKJVw=', 'base64')
  const crypto = require('@jadepool/crypto')
  const moment = require('moment')
  const axios = require('axios')
  let msg = {
      "id":"3000",
      "state": status,
      "bizType":"DEPOSIT",
      "coinType":"ETH",
      "to":"0xEcce5fDF42da3b3E77833E8902fC8112ef78437f",
      "value":"0.001",
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
      url,
      data: data,
      proxy: false
  })
  console.log('result', status, result.data)
}


const sleepForSeconds = async (s) => {
  return new Promise(resolve => {
    setTimeout(resolve, s * 1000)
  })
}

const insertFirst = async () => {
  let data = {
    
  }
}

const test = async () => {
  console.log('test start')
  await emptyDb()
  await fireDeposit('pending')
  await sleepForSeconds(3)
  let t1 = fireDeposit('pending')
  let t2 = fireDeposit('done')
  let t3 = fireDeposit('pending')
  let t4 = fireDeposit('pending')
  let t5 = fireDeposit('pending')
  let t6 = fireDeposit('pending')
  let t7 = fireDeposit('pending')
  let t8 = fireDeposit('pending')
  let t9 = fireDeposit('pending')
  let t10 = fireDeposit('pending')
  await Promise.all([t1, t2, t3])
  // await Promise.all([t1, t2, t3, t4, t5, t6, t7, t8, t9, t10])
  console.log('test finished')
  process.exit(0)
}

test()
