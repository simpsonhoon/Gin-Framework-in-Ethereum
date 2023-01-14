package service

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func TransferCtxSSH(from string, to string, _value int64) (receipt string, err error) {
	// 블록체인 네트워크와 연결할 클라이언트를 생성하기 위한 rpc url 연결
	client, err := ethclient.Dial("https://api.test.wemix.com")
	if err != nil {
		fmt.Println("client error")
	}

	// 본인이 배포한 토큰 컨트랙트 어드레스
	tokenAddress := common.HexToAddress("0xb58E525a38bb9Dc9Fe4fb3C2b957f7A9863093bF")

	privateKey, err := crypto.HexToECDSA(from) // privateKey(개인키)
	if err != nil {
		fmt.Println(err)
	}

	//privateKey -> publicKey 변환 -> 자신의 address 변환
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 현재 계정의 nonce를 가져옴. 다음 트랜잭션에서 사용할 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		fmt.Println(err)
	}

	// 전송할 양, gasLimit, gasPrice 설정. 추천되는 gasPrice를 가져옴
	value := big.NewInt(_value)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	// 보낼 주소
	toAddress := common.HexToAddress(to)

	// 컨트랙트 전송시 사용할 함수명
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID))

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	paddedAmount := common.LeftPadBytes(value.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000
	zvalue := big.NewInt(0)

	//컨트랙트 전송 정보 입력
	var pdata []byte
	pdata = append(pdata, methodID...)
	pdata = append(pdata, paddedAddress...)
	pdata = append(pdata, paddedAmount...)

	gasLimit := uint64(200000)
	fmt.Println(gasLimit)

	// 트랜잭션 생성
	tx := types.NewTransaction(nonce, tokenAddress, zvalue, gasLimit, gasPrice, pdata)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	// 트랜잭션 서명
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println(err)
	}

	// 트랜잭션 전송
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		fmt.Println(err)
	}

	return signedTx.Hash().Hex(), nil
}
