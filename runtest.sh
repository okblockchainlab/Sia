
go test -v ./okwallet/okwallet
go test -race -v -tags='debug testing vlong  netgo' -timeout=5000s ./node/api -run="TestWalletCommitTransactions"
