    const pubKey = Buffer.from('03ace32532c90652e1bae916248e427a7ab10aeeea1067949669a3f4da10965ef9', 'hex')
    const crypto = require('@jadepool/crypto')

    let bodyMsg = JSON.parse('{"crypto":"ecc","hash":"sha3","encode":"base64","appid":"app","timestamp":1548411683000,"sig":{"r":"W8ehp5kPnH1c+yn1SkRgLC5NTydQ7CHp3C9e0kJLO24=","s":"P3bqh3FOmdHi3Vf8oIg/R0IZJ/BOHeG7+/p8+/UUjsY=","v":0},"data":{"type":"","value":"","to":"","timestamp":1548411683000,"callback":"","extraData":""}}')
    const msgCrypto = bodyMsg.crypto || 'ecc'
    const msgHash = bodyMsg.hash || 'sha3'
    const msgEncode = bodyMsg.encode || 'base64'
    const appid = bodyMsg.appid || 'app'
    console.log(bodyMsg.sig)
    let verifyOk = crypto.ecc.verify(bodyMsg.data, bodyMsg.sig, pubKey, {hash: msgHash, crypto: msgCrypto, encode: msgEncode})
    console.log("verify result:", verifyOk)

    bodyMsg = JSON.parse('{"crypto":"ecc","hash":"sha3","encode":"base64","appid":"app","timestamp":1548661151000,"sig":{"r":"D5W5WYkKGZiYSiT5/SxYTmTtv+zCAeb88yBCSK69jQg=","s":"Ifc8aEid0Mbh33Z6i5lluSk25KiWEXMQrNwDRyUJe2M=","v":0},"data":{"type":"","timestamp":1548661151000,"callback":""}}')
    console.log(bodyMsg.sig)
    verifyOk = crypto.ecc.verify(bodyMsg.data, bodyMsg.sig, pubKey, {hash: msgHash, crypto: msgCrypto, encode: msgEncode})
    console.log("verify result:", verifyOk)