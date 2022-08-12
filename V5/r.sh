rm ./main.exe
rm ./blockChain.db
rm ./blockChain.db.lock
go build main.go block.go blockchain.go auxiliaryFunc.go proofOfWork.go iterator.go cli.go commandLine.go transcation.go getBalance.go
./main.exe transfer Satoshi xiaoming 5.08 Alice "巴黎"
./main.exe transfer xiaoming Alice 1.25 Satoshi "伦敦"
./main.exe transfer Alice Satoshi 2.99 xiaoming "纽约"
./main.exe getBalance Satoshi
./main.exe getBalance xiaoming
./main.exe getBalance Alice
./main.exe printChain
read -p "press enter end"