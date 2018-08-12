### 编译
```shell
export GOPATH=$GOPATH:/your/go/path/directory  #设置GOPATH路径
cd /your/go/path/directory
git clone https://github.com/okblockchainlab/Sia.git ./gitlab.com/NebulousLabs/Sia
cd ./gitlab.com/NebulousLabs/Sia
./build.sh #run this script only if you first time build the project
./runbuild.sh
./runtest.sh
```

### 新增API
服务进程新增加了两个http形式的API接口，以方便实现冷热钱包分离的功能。调用代码示例可以参见node/api/wallet_test.go文件中的TestWalletCommitTransactions函数。

##### /wallet/checkoutput [GET]
此命令用来检查某个transaction的output是否可以花费。在调用CreateRawTransaction时，调用者必须确保所给出的所有output是可花费的，此命令可以帮助调用者进行筛选。  
输入：  
```
//要检查的transaction的base64的数据  
:transaction  
```
输出示例：
```json
{
		"Spendable":[1, 2],
		"Unspendable": [0]
}
```

##### /wallet/committransactions [POST]
提交已签名的transactions。可以一次提交多个。  
输入：
```
//transaction数组
:transactions
```
输出：
```json
{
  "transactionid" : [
      "1234567890abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  ]
}
```
