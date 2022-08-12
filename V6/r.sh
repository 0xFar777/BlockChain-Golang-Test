rm ./main.exe
rm ./blockChain.db
rm ./blockChain.db.lock

go build main.go block.go blockchain.go auxiliaryFunc.go \
proofOfWork.go iterator.go cli.go commandLine.go transcation.go \
getBalance.go wallet.go wallets.go

./main.exe transfer Satoshi xiaoming 5.08 Alice "巴黎"
./main.exe transfer xiaoming Alice 1.25 Satoshi "伦敦"
./main.exe transfer Alice Satoshi 2.99 xiaoming "纽约"

./main.exe getBalance Satoshi
./main.exe getBalance xiaoming
./main.exe getBalance Alice

./main.exe printChain

# 创建20个钱包
for i in {1..20}
do
    ./main.exe NewWallet
    let i++
done

./main.exe listAddress

read -p "press enter end"