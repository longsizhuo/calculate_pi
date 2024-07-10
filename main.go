package main

import (
	"fmt"
	"math/big"
	"os"
	"sync"
)

// 文件路径
const FILENAME = "pi_digits.txt"

// 初始化精度
const initialPrec = 100

// Chudnovsky算法计算圆周率
func chudnovskyAlgorithm(prec int) *big.Float {
	// 设置高精度
	pi := new(big.Float).SetPrec(uint(prec))

	// 初始化常数
	C := new(big.Float).SetPrec(uint(prec)).Mul(big.NewFloat(426880), new(big.Float).Sqrt(big.NewFloat(10005)))

	K := new(big.Int).SetInt64(6)
	M := new(big.Int).SetInt64(1)
	X := new(big.Int).SetInt64(1)
	L := new(big.Int).SetInt64(13591409)
	S := new(big.Float).SetPrec(uint(prec)).SetInt(L)

	temp := new(big.Float).SetPrec(uint(prec))
	for k := int64(1); k <= int64(prec)/14+1; k++ {
		M.Mul(M, new(big.Int).Div(new(big.Int).Sub(new(big.Int).Mul(K, new(big.Int).Mul(K, K)), new(big.Int).Mul(big.NewInt(16), K)), new(big.Int).Mul(new(big.Int).SetInt64(k), new(big.Int).Mul(new(big.Int).SetInt64(k), new(big.Int).SetInt64(k)))))
		L.Add(L, big.NewInt(545140134))
		X.Mul(X, big.NewInt(-262537412640768000))
		temp.SetInt(M)
		temp.Mul(temp, new(big.Float).SetPrec(uint(prec)).SetInt(L))
		temp.Quo(temp, new(big.Float).SetPrec(uint(prec)).SetInt(X))
		S.Add(S, temp)
		K.Add(K, big.NewInt(12))
	}

	pi.Quo(C, S)
	return pi
}

// 读取已有的圆周率位数
func readExistingDigits() string {
	if _, err := os.Stat(FILENAME); err == nil {
		content, err := os.ReadFile(FILENAME)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return ""
		}
		return string(content)
	}
	return ""
}

// 写入新的圆周率位数
func writeDigits(digits string) {
	file, err := os.OpenFile(FILENAME, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(digits)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

// 计算圆周率函数
func calculatePi() {
	existingDigits := readExistingDigits()
	currentPosition := len(existingDigits)

	currentPrec := initialPrec
	if currentPosition > 0 {
		currentPrec = currentPosition + 100 // 根据已有位数动态设置初始精度
	}

	resultChan := make(chan *big.Float)
	var wg sync.WaitGroup

	go func() {
		for pi := range resultChan {
			piStr := pi.Text('f', currentPrec)[2:] // 去除"3."
			for currentPosition < len(piStr) {
				newDigits := piStr[currentPosition:min(currentPosition+50, len(piStr))]
				writeDigits(newDigits) // 写入新计算的位数
				currentPosition += 50
			}
			wg.Done()
		}
	}()

	for {
		currentPrec += 50 // 增加精度
		wg.Add(1)
		go func(prec int) {
			resultChan <- chudnovskyAlgorithm(prec)
		}(currentPrec)

		wg.Wait()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	calculatePi()
}
