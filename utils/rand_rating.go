package utils

import (
	"crypto/rand"
	"math"
	"math/big"
)

// secureRandFloat64 使用 crypto/rand 生成安全的随机浮点数 [0.0, 1.0)
func secureRandFloat64() float64 {
	maxInt := big.NewInt(1 << 53)
	n, _ := rand.Int(rand.Reader, maxInt)
	return float64(n.Int64()) / float64(maxInt.Int64())
}

// GenerateRating 生成随机应用评分 (1.0 - 5.0)
// 返回值保留一位小数
func GenerateRating() float64 {
	// 生成 1.0 到 5.0 之间的随机评分
	rating := 1.0 + secureRandFloat64()*4.0
	// 保留一位小数
	return math.Round(rating*10) / 10
}

// GenerateRatingWithDistribution 生成符合正态分布的随机评分
// 大多数评分会集中在 3.5 - 4.5 之间，更符合实际应用评分情况
func GenerateRatingWithDistribution() float64 {
	// 使用 Box-Muller 变换生成正态分布
	// 均值为 4.0，标准差为 0.5
	mean := 4.0
	stdDev := 0.5

	u1 := secureRandFloat64()
	u2 := secureRandFloat64()

	// Box-Muller 变换
	z := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)
	z = mean + stdDev*z

	// 限制在 1.0 到 5.0 之间
	if z < 1.0 {
		z = 1.0
	}
	if z > 5.0 {
		z = 5.0
	}

	// 保留一位小数
	return math.Round(z*10) / 10
}

// GenerateRatingWithWeight 根据权重生成评分
// 可以设置不同星级的概率
func GenerateRatingWithWeight() float64 {
	// 定义权重：5星(40%), 4星(30%), 3星(20%), 2星(7%), 1星(3%)
	weights := []float64{0.03, 0.07, 0.20, 0.30, 0.40}
	cumWeights := make([]float64, len(weights))

	// 计算累积权重
	cumWeights[0] = weights[0]
	for i := 1; i < len(weights); i++ {
		cumWeights[i] = cumWeights[i-1] + weights[i]
	}

	// 生成随机数
	r := secureRandFloat64()

	// 确定星级
	var starBase int
	for i, w := range cumWeights {
		if r <= w {
			starBase = i + 1
			break
		}
	}

	// 在该星级范围内生成随机小数
	rating := float64(starBase) + secureRandFloat64()*0.9

	// 保留一位小数
	return math.Round(rating*10) / 10
}
