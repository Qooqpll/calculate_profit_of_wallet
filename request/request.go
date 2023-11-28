package request

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	url    string = "https://deep-index.moralis.io/api/v2.2/"
	isMock bool   = false
)

type QueryParams struct {
	Address   string
	Chain     string
	FromBlock int
	ToBlock   int
}

type Request struct {
	Url string
}

func newQueryParams(address, chain string, fromBlock, toBlock int) *QueryParams {
	return &QueryParams{
		Address:   address,
		Chain:     chain,
		FromBlock: fromBlock,
		ToBlock:   toBlock,
	}
}

func (r *Request) sendRequestByUrl(payloadData, method string) []byte {
	payload := strings.NewReader(payloadData)
	req, _ := http.NewRequest(method, r.Url, payload)

	req.Header.Add("Accept", "application/json")
	if method == "POST" {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("X-API-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6IjBhMTBiN2UwLTM1YTEtNDgzYi1iNGRhLTJiMDI2YWExZmQwYiIsIm9yZ0lkIjoiMzY1NDI1IiwidXNlcklkIjoiMzc1NTYzIiwidHlwZUlkIjoiYjdiMDMwY2MtN2UwYi00ZjE3LWFhMjEtMzA1ZWFmY2VhNzFkIiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE3MDA2NjkyOTksImV4cCI6NDg1NjQyOTI5OX0.TJWpx0e1HSDhaepuJl6LWp_1wKunjxuOZg9N20QPqlU")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func (q *QueryParams) getUrlForNativeTransactionsByWallet() *Request {
	s := url + q.Address + "/verbose?chain=" + q.Chain
	return &Request{
		Url: s,
	}
}

func (q *QueryParams) getUrlForTokenTransfers() *Request {
	s := url + q.Address + "/erc20/transfers?chain=" + q.Chain
	return &Request{
		Url: s,
	}
}

func (q *QueryParams) getUrlForGetTokenPrices() *Request {
	s := url + "erc20/prices?chain=" + q.Chain
	return &Request{
		Url: s,
	}
}
