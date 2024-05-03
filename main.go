package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	PUB_KEY = `0x5F95613eE4BA6933DD81333cB5a0e0aa0F83fF24`
)

type Addr struct {
	Router02 string `json:"router02"`
	Factory  string `json:"factory"`
}

func InitAddr() (Addr, error) {
	file, err := os.Open("addresses.json")
	if err != nil {
		return Addr{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Addr{}, err
	}

	var ret Addr
	if err := json.Unmarshal(data, &ret); err != nil {
		return Addr{}, err
	}

	return ret, nil
}

func main() {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/" + os.Getenv("INFURA_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("connected")

	addr, err := InitAddr()
	if err != nil {
		log.Fatal(err)
	}

	addresses := []common.Address{common.HexToAddress(addr.Factory), common.HexToAddress(addr.Router02)}

	query := ethereum.FilterQuery{
		Addresses: addresses,
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case log := <-logs:
				fmt.Println(log)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	abiFile, err := os.Open("./uniswap/IUniswapV2Factory.abi")
	if err != nil {
		log.Fatal(err)
	}
	defer abiFile.Close()

	data, err := io.ReadAll(abiFile)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}

	fmt.Println(contractAbi)
}
