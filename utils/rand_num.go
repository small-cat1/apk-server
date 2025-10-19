package utils

import (
	"crypto/rand"
	"log"
	"math/big"
)

func RandNumber(minNum, maxNum int64) int64 {
	// 计算范围
	n := maxNum - minNum + 1
	// 生成一个安全的随机数
	// rand.Int() 返回一个 *big.Int
	// rand.Int(rand.Reader, n) 生成 [0, n) 范围内的随机数
	// 在这里，我们生成 [0, n) 范围内的随机数
	bigRand, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		log.Fatal(err)
	}
	// 将随机数转换为 int64，并加上最小值
	// 结果是 [min, max] 范围内的随机数
	randomNum := bigRand.Int64() + minNum
	return randomNum
}
