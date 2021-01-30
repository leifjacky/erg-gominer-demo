## Ergo(ERG) cpuminer written in Golang

Test only.

### usage

```bash
go get -v github.com/leifjacky/erg-gominer-demo
cd $GOPATH/src/github.com/leifjacky/erg-gominer-demo
go run *.go
```

## Ergo(ERG) stratum protocol

### mining.subscribe

```json
{"id":1,"method":"mining.subscribe","params":["ergminer-v1.0.0"]}

{"id":"mining.subscribe","result":[[["mining.set_difficulty","0.000031"],["mining.notify","000400001f164ca2"]],"d7eb83fe",4],"error":null}		// nonce1, nonce2 size
```



### mining.authorize

```json
{"id":2,"method":"mining.authorize","params":["3WyMz9efAu7XtQXKtjcDBrXBhQgak8AceQsnGtXg1JcK9ceCiw9x.worker1","x"]}	// params: [username, password]

{"id":2,"result":true,"error":null}
{"id":null,"method":"mining.set_target","params":["00007fffffffffffffffffffffffffffffffffffffffffffffffffffffff8000"]}
```



### mining.notify

```json
{
	"id": null,
	"method": "mining.notify",
	"params": ["47388824","b15c08abeb302e990330fda83cbf7835a2172bbfbd8542796fb4255c1aaf4bde",4738,true]
} // params: [jobId, msg, height, cleanJob]
```



### mining.submit

```json
{"id":102,"method":"mining.submit","params":["3WxsCRcd9e2fGDHgxk1ptn3stiXPG4p2fi6im4VzQgfQJuS9YSN8.worker1","47388824","00239534"]}  // params: [username, jobId, nonce2]

{"id":100,"result":true,"error":null}    // accepted share response
{"id":100,"result":false,"error":[21,"low difficulty",null]}  // rejected share response
```



```bash
MINER LOG:
solving nonce2: 00239534
msgWithNonce: b15c08abeb302e990330fda83cbf7835a2172bbfbd8542796fb4255c1aaf4bded7eb83fe00239534
msgHash: 576b18b99605a586654026614e50f690295580ab14a0dce645cb076480841eb5
i: 00841eb5
h: 00001282
e: 2c6b53fe68958f676183583d49040f95512992c576a100c4a32c716abcc874
J: [03ea5f7c 025f7cf6 037cf65e 00f65e81 025e81a4 0281a487 01a4878e 00878ecc 038ecc18 02cc1849 001849d4 0049d448 01d44880 004880d2 0080d2b8 00d2b833 02b83394 00339478 039478c7 0078c77c 00c77c47 037c4778 00477829 0378296d 00296ddb 016ddbf5 01dbf5b3 03f5b37a 01b37a21 037a2167 022167ea 0167ea5f]
f: 116ee0186a8bf4c0cfa57f0c2aa629e170fb0171ea526bb6d04ba05a6773fa6f
fh: 000079d2cfe14668445d1d78eb4fedf6676ae5c762260b0f1132752d84236b27
share found: 00239534 - 000079d2cfe14668445d1d78eb4fedf6676ae5c762260b0f1132752d84236b27
```

