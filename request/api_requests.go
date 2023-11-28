package request

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GetTokenTransferByWallet() []byte {
	if isMock {
		jsonFile, err := os.Open("request/mock/transactions.json")
		if err != nil {
			fmt.Println(err)
		}

		defer jsonFile.Close()

		body, _ := ioutil.ReadAll(jsonFile)
		return body
	}

	queryParams := newQueryParams("0xfd20b05a83f7e956bda2dbab58d2062eb914bc12", "eth", 0, 0)
	url := queryParams.getUrlForTokenTransfers()
	return url.sendRequestByUrl("", "GET")

}

func GetTokenPrices(data []TokenData) []byte {
	if isMock {
		jsonFile, err := os.Open("request/mock/tokens.json")
		if err != nil {
			fmt.Println(err)
		}

		defer jsonFile.Close()

		body, _ := ioutil.ReadAll(jsonFile)
		return body
	}
	payloadData := "{\"tokens\":["
	for i, RequestToken := range data {
		if i > 0 {
			payloadData += ","
		}
		payloadData += "{\"token_address\":\"" + RequestToken.TokenAddress + "\",\"to_block\":\"" + RequestToken.ToBlock + "\"}"
	}
	payloadData += "]}"

	queryParams := newQueryParams("0xfd20b05a83f7e956bda2dbab58d2062eb914bc12", "eth", 0, 0)
	url := queryParams.getUrlForGetTokenPrices()

	return url.sendRequestByUrl(payloadData, "POST")
}
