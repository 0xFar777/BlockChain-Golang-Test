rm ./main.exe
rm ./blockChain.db
rm ./blockChain.db.lock

go build \
main.go block.go blockchain.go auxiliaryFunc.go \
proofOfWork.go iterator.go cli.go commandLine.go transcation.go \
getBalance.go wallet.go wallets.go sign.go verify.go

./main.exe getBalance 1Kwjqnee9o55JCt6SVPuAoTtDm9MMwVoPk
./main.exe getBalance 1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird
./main.exe getBalance 1GWtvNnn5f4vgjwPnEfsG3P56jBdhmkNpq

./main.exe transfer 1Kwjqnee9o55JCt6SVPuAoTtDm9MMwVoPk \
1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird 5.08 1GWtvNnn5f4vgjwPnEfsG3P56jBdhmkNpq "Trans success002"

./main.exe getBalance 1Kwjqnee9o55JCt6SVPuAoTtDm9MMwVoPk
./main.exe getBalance 1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird
./main.exe getBalance 1GWtvNnn5f4vgjwPnEfsG3P56jBdhmkNpq

./main.exe transfer 1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird \
1GWtvNnn5f4vgjwPnEfsG3P56jBdhmkNpq 1.25 1Kwjqnee9o55JCt6SVPuAoTtDm9MMwVoPk "Trans success003"

./main.exe getBalance 1Kwjqnee9o55JCt6SVPuAoTtDm9MMwVoPk
./main.exe getBalance 1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird
./main.exe getBalance 1GWtvNnn5f4vgjwPnEfsG3P56jBdhmkNpq

./main.exe transfer 1GWtvNnn5f4vgjwPnEfsG3P56jBdhmkNpq \
1Kwjqnee9o55JCt6SVPuAoTtDm9MMwVoPk 2.99 1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird "Trans success004"

./main.exe getBalance 1Kwjqnee9o55JCt6SVPuAoTtDm9MMwVoPk
./main.exe getBalance 1Am9Nu5VJ4aDMtQ4Mc89WKuPmpn5qSJird
./main.exe getBalance 1GWtvNnn5f4vgjwPnEfsG3P56jBdhmkNpq

#创建5个钱包
for i in {1..5}
do
    ./main.exe NewWallet
    let i++
done

./main.exe printChain

# ./main.exe listAddress

read -p "press enter end"