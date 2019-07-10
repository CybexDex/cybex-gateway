const config = require('./config')
const pg = require('pg')
const pool = new pg.Pool(config.db)

const insertWithDraw = async () => {
  const client = await pool.connect()
  let now = Math.round(new Date().getTime()/1000)
  let sql = `INSERT INTO jp_orders(created_at,updated_at,asset,block_chain,cyb_user,out_addr,total_amount,amount,fee,status,type,current,current_state) VALUES(to_timestamp(${now}),to_timestamp(${now}),'NASH','ETH','yinnan-test1','0x3ae306d3fe3584ec90a765db587815b6d990ce4a',0.2,0.15,0.05,'jporder','WITHDRAW','jp','INIT');`
  console.log('---------------------------------')
  console.log(sql)
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

const test = async () => {
  let loop = 2
  for (let i = 0; i < loop; i++) {
    await insertWithDraw()
  }
  process.exit(0)
}

test()
