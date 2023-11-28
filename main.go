package main

import (
	"CalculateProfitLose/database"
	"CalculateProfitLose/request"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strconv"
	"time"
)

type Transactions struct {
	PageSize int        `json:"page_size"`
	Page     int        `json:"page"`
	Cursor   string     `json:"cursor"`
	Result   []Transfer `json:"result"`
}

type Transfer struct {
	TokenName       string    `json:"token_name"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	TokenAddress    string    `json:"address"`
	BlockTimestamp  time.Time `json:"block_timestamp"`
	Value           string    `json:"value"`
	BlockNumber     string    `json:"block_number"`
	TransactionHash string    `json:"transaction_hash"`
}

const (
	Buy            = "buy"
	Sell           = "sell"
	CurrentAddress = "0xfd20b05a83f7e956bda2dbab58d2062eb914bc12"
)

// TODO разработать скрипт который будет менять ключи по истечению Compute units в сервисе moralis
func main() {
	database.Connect()

	tokens := []request.TokenData{}

	body := request.GetTokenTransferByWallet()

	Transactions := Transactions{}
	err := json.Unmarshal(body, &Transactions)
	if err != nil {
		fmt.Println(err)
	}

	transfers := removeTokensDuplicate(Transactions.Result)

	for _, transaction := range transfers {

		dataForTokenRequest := request.TokenData{TokenAddress: transaction.TokenAddress, ToBlock: transaction.BlockNumber}
		tokens = append(tokens, dataForTokenRequest)

		fmt.Println("action: " + getAction(transaction.ToAddress, transaction.FromAddress))
		fmt.Println("tokenName: " + transaction.TokenName)
		fmt.Println("TokenAddress: " + transaction.TokenAddress)
		fmt.Println("from: " + transaction.FromAddress)
		fmt.Println("to: " + transaction.ToAddress)
		fmt.Println("time: ", transaction.BlockTimestamp.Unix())
		fmt.Println("BlockNumber: " + transaction.BlockNumber)
		fmt.Println("Hash: " + transaction.TransactionHash)
		fmt.Println("Value: " + transaction.Value)
		fmt.Println("-------------")

	}

	sort.Slice(Transactions.Result, func(i, j int) bool {
		return Transactions.Result[i].BlockNumber < Transactions.Result[j].BlockNumber
	})

	var tokenPrices []byte
	var n0, n1 int
	length := len(tokens)

	for i := 0; i < length; i += 25 {
		n1 = n0 + 25
		if n1 > length {
			n1 = length
		}

		req := request.GetTokenPrices(tokens[n0:n1])
		if i > 0 {
			// Add a comma between JSON arrays except for the first one
			tokenPrices = append(tokenPrices, ',')
		}
		tokenPrices = append(tokenPrices, req...)
		n0 = n1
	}

	// Wrap the combined JSON arrays with square brackets to form a valid JSON array
	tokenPrices = append([]byte{'['}, append(tokenPrices, ']')...)

	fmt.Println(string(tokenPrices))

	var infoTokens [][]request.Token
	err = json.Unmarshal(tokenPrices, &infoTokens)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	} else {
		fmt.Println("Number of tokens:", len(infoTokens))
	}
	fmt.Println("==================")
	fmt.Println(infoTokens)
	fmt.Println("==================")
	var result = make(map[string]float64)

	for _, transaction := range transfers {
		for i, token := range infoTokens {
			if transaction.BlockNumber == token[i].ToBlock && transaction.TokenAddress == token[i].TokenAddress {
				if transaction.ToAddress == CurrentAddress {
					// Buy
					if _, ok := result[transaction.TokenAddress]; !ok {
						result[transaction.TokenAddress] = float64(0)
					}
					result[transaction.TokenAddress] -= transaction.CalculateTokenPrice(token[i].UsdPriceFormatted, token[i].TokenDecimals)
				} else if transaction.FromAddress == CurrentAddress {
					// Sell
					if _, ok := result[transaction.TokenAddress]; !ok {
						result[transaction.TokenAddress] = float64(0)
					}
					result[transaction.TokenAddress] += transaction.CalculateTokenPrice(token[i].UsdPriceFormatted, token[i].TokenDecimals)
				}
			}
		}
	}
	fmt.Println(result)
	fmt.Println(calculateProfitLoss(result))
}

func (t *Transfer) CalculateTokenPrice(tokenPrice, tokenDecimals string) float64 {
	price := new(big.Float)
	price.SetString(tokenPrice)
	value := new(big.Float)
	value.SetString(t.Value)
	res := new(big.Float).Mul(price, value)
	decimals, err := strconv.Atoi(tokenDecimals)
	if err != nil {
		fmt.Println(err)
	}
	result := new(big.Float).Quo(res, big.NewFloat(math.Pow10(decimals)))
	r, _ := result.Float64()
	return r
}

func getAction(toAddress, fromAddress string) string {
	if toAddress == CurrentAddress {
		// Buy
		return "buy"
	} else if fromAddress == CurrentAddress {
		// Sell
		return "sell"
	}
	return ""
}

func removeTokensDuplicate(transfers []Transfer) []Transfer {
	// Шаг 1: Суммируем значения для одинаковых TokenAddress и BlockNumber
	sumMap := make(map[string]map[int]float64)
	for _, transfer := range transfers {
		if sumMap[transfer.TokenAddress] == nil {
			sumMap[transfer.TokenAddress] = make(map[int]float64)
		}
		// Преобразуем BlockNumber в int перед использованием в качестве индекса
		blockNumber, err := strconv.Atoi(transfer.BlockNumber)
		if err != nil {
			// Обработка ошибки, например, вывод сообщения или возврат 0
			fmt.Println("Ошибка конвертации BlockNumber в int:", err)
			return nil
		}
		sumMap[transfer.TokenAddress][blockNumber] += parseFloat(transfer.Value)
	}

	// Шаг 2: Создаем уникальный массив учитывая суммы
	result := []Transfer{}
	for _, transfer := range transfers {
		// Получаем сумму для данного TokenAddress и BlockNumber
		blockNumber, err := strconv.Atoi(transfer.BlockNumber)
		if err != nil {
			// Обработка ошибки, например, вывод сообщения или возврат 0
			fmt.Println("Ошибка конвертации BlockNumber в int:", err)
			return nil
		}
		sum := sumMap[transfer.TokenAddress][blockNumber]
		// Конвертируем сумму обратно в строку
		sumString := fmt.Sprintf("%.6f", sum)
		// Проверяем, был ли этот элемент добавлен ранее
		if !containsTransfer(result, transfer) {
			// Если нет, то добавляем его с учетом суммы
			transfer.Value = sumString
			result = append(result, transfer)
		}
	}

	return result
}

// Функция для проверки наличия элемента в массиве
func containsTransfer(transfers []Transfer, t Transfer) bool {
	for _, transfer := range transfers {
		if transfer.TokenAddress == t.TokenAddress && transfer.BlockNumber == t.BlockNumber {
			return true
		}
	}
	return false
}

// Функция для конвертации строки во float64
func parseFloat(s string) float64 {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// Обработка ошибки, например, вывод сообщения или возврат 0
		fmt.Println("Ошибка конвертации значения в float64:", err)
		return 0
	}
	return value
}

func calculateProfitLoss(transfers map[string]float64) float64 {
	var profit float64
	for _, transfer := range transfers {
		profit += transfer
	}

	return profit
}
