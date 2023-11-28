package request

import (
	"io/ioutil"
	"net/http"
)

type Tokens struct {
	Tokens []Token
}

type Token struct {
	TokenName         string      `json:"tokenName"`
	TokenDecimals     string      `json:"tokenDecimals"`
	UsdPriceFormatted string      `json:"usdPriceFormatted"`
	UsdPrice          float64     `json:"usdPrice"`
	ToBlock           string      `json:"toBlock"`
	TokenAddress      string      `json:"tokenAddress"`
	NativePrice       NativePrice `json:"nativePrice"`
}

type NativePrice struct {
	Decimals int    `json:"decimals"`
	Value    string `json:"value"`
}

type TokenData struct {
	TokenAddress string `json:"token_address"`
	Exchange     string `json:"exchange"`
	ToBlock      string `json:"to_block"`
}

func RequestGetTokenPrice(tokenAddress, blockNumber string) []byte {
	url := "https://deep-index.moralis.io/api/v2.2/erc20/" + tokenAddress + "/price?chain=eth&to_block=" + blockNumber

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6IjE0OGY0Mjc5LWQ1NjItNDkwNi1iM2JkLTg2MGYzNjRhYmYxNyIsIm9yZ0lkIjoiMzYzNDEwIiwidXNlcklkIjoiMzczNDkwIiwidHlwZUlkIjoiYjdjODQwYzUtMTlhMS00ODJmLWFmYmEtNGZhNDc0ZmY0OGQxIiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE2OTkxODM4MTMsImV4cCI6NDg1NDk0MzgxM30.GeW9nWUiKba2Nn_0KjtrdKpOk_p00MQbmanXtCvPdBA")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}
