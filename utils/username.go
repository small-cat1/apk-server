package utils

import (
	"fmt"
	"math/rand"
)

// 网名生成器接口
type UsernameGenerator interface {
	Generate() string
	GetStyle() string
}

// 1. 萌系/可爱风格
type CuteGenerator struct{}

func (g *CuteGenerator) Generate() string {
	prefixes := []string{
		"软萌", "呆萌", "甜甜", "软软", "奶茶", "草莓",
		"樱花", "柠檬", "蜜桃", "棉花糖", "布丁", "果冻",
		"奶油", "糖果", "蛋糕", "冰淇淋", "巧克力", "曲奇",
	}

	nouns := []string{
		"小仙女", "小公主", "宝宝", "酱", "喵", "兔",
		"熊", "鸭", "猫猫", "团子", "丸子", "精灵",
		"天使", "萝莉", "少女", "妹妹", "姐姐", "宝贝",
	}

	return prefixes[rand.Intn(len(prefixes))] + nouns[rand.Intn(len(nouns))]
}

func (g *CuteGenerator) GetStyle() string {
	return "萌系可爱"
}

// 2. 古风/文艺风格
type AncientGenerator struct{}

func (g *AncientGenerator) Generate() string {
	templates := []func() string{
		// 模板1: 地点+景物
		func() string {
			places := []string{"陌上", "南巷", "北城", "西窗", "东篱", "江畔", "湖边", "山间"}
			scenes := []string{"烟雨", "青苔", "落花", "残雪", "孤影", "流年", "浮生", "微凉"}
			return places[rand.Intn(len(places))] + scenes[rand.Intn(len(scenes))]
		},
		// 模板2: 动词+名词+意境
		func() string {
			verbs := []string{"听", "看", "醉", "忆", "寻", "等", "念", "望"}
			nouns := []string{"风", "雨", "雪", "月", "星", "云", "花", "叶"}
			moods := []string{"吟", "语", "歌", "诗", "梦", "影", "声", "韵"}
			return verbs[rand.Intn(len(verbs))] + nouns[rand.Intn(len(nouns))] + moods[rand.Intn(len(moods))]
		},
		// 模板3: 颜色+物品
		func() string {
			colors := []string{"青", "白", "红", "紫", "碧", "墨", "素", "朱"}
			items := []string{"衫", "衣", "尘", "陌", "落", "笺", "墨", "笔"}
			suffixes := []string{"客", "人", "者", "书生", "过客", "如故"}
			return colors[rand.Intn(len(colors))] + items[rand.Intn(len(items))] + suffixes[rand.Intn(len(suffixes))]
		},
	}

	return templates[rand.Intn(len(templates))]()
}

func (g *AncientGenerator) GetStyle() string {
	return "古风文艺"
}

// 3. 霸气/游戏风格
type CoolGenerator struct{}

func (g *CoolGenerator) Generate() string {
	prefixes := []string{
		"至尊", "霸王", "狂", "傲", "战神", "魔", "帝",
		"王者", "无敌", "绝世", "冷血", "暗影", "嗜血",
	}

	nouns := []string{
		"杀手", "剑客", "刀客", "枪神", "法师", "龙", "虎",
		"狼", "鹰", "狮", "战士", "猎人", "刺客",
	}

	suffixes := []string{
		"", "无双", "归来", "降临", "觉醒", "重生",
		"天下", "称霸", "独尊", "纵横",
	}

	prefix := prefixes[rand.Intn(len(prefixes))]
	noun := nouns[rand.Intn(len(nouns))]
	suffix := suffixes[rand.Intn(len(suffixes))]

	return prefix + noun + suffix
}

func (g *CoolGenerator) GetStyle() string {
	return "霸气游戏"
}

// 4. 搞笑/沙雕风格
type FunnyGenerator struct{}

func (g *FunnyGenerator) Generate() string {
	templates := []func() string{
		// 模板1: 形容词+小+动物
		func() string {
			adjectives := []string{"可爱", "迷人", "性感", "帅气", "优雅", "高冷", "傲娇", "温柔"}
			animals := []string{"猪猪", "狗狗", "鸭鸭", "熊熊", "兔兔", "喵喵", "羊羊"}
			return adjectives[rand.Intn(len(adjectives))] + "小" + animals[rand.Intn(len(animals))]
		},
		// 模板2: 网络热词
		func() string {
			phrases := []string{
				"在逃公主", "隐藏巨星", "人间清醒", "打工人",
				"社恐患者", "熬夜冠军", "网瘾少女", "肥宅快乐",
				"划水选手", "摆烂专家", "吃瓜群众", "吸猫患者",
				"秃头少女", "佛系青年", "养生朋克", "熬夜冠军",
			}
			return phrases[rand.Intn(len(phrases))]
		},
		// 模板3: 动作+身份
		func() string {
			actions := []string{"爱学习的", "爱睡觉的", "爱吃饭的", "爱玩的", "想躺平的", "在摸鱼的"}
			identities := []string{"小废物", "小可爱", "小天才", "小机灵", "小迷糊", "小懒虫"}
			return actions[rand.Intn(len(actions))] + identities[rand.Intn(len(identities))]
		},
	}

	return templates[rand.Intn(len(templates))]()
}

func (g *FunnyGenerator) GetStyle() string {
	return "搞笑沙雕"
}

// 5. 混搭风格（英文+中文）
type MixedGenerator struct{}

func (g *MixedGenerator) Generate() string {
	englishWords := []string{
		"Sunny", "Moon", "Star", "Sky", "Cloud", "Dream",
		"Angel", "Devil", "King", "Queen", "Ice", "Fire",
		"Dark", "Light", "Shadow", "Night", "Dawn",
	}

	chineseWords := []string{
		"小姐姐", "小哥哥", "宝贝", "酱", "君", "子",
		"仙", "神", "大魔王", "小可爱", "殿下", "公主",
	}

	eng := englishWords[rand.Intn(len(englishWords))]
	chn := chineseWords[rand.Intn(len(chineseWords))]

	// 随机决定顺序
	if rand.Float32() < 0.5 {
		return eng + chn
	}
	return chn + eng
}

func (g *MixedGenerator) GetStyle() string {
	return "中英混搭"
}

// 6. 符号装饰风格
type DecoratedGenerator struct{}

func (g *DecoratedGenerator) Generate() string {
	names := []string{
		"星辰", "月光", "晨曦", "暮色", "流年", "梦境",
		"幻想", "迷雾", "清风", "落雪", "烟火", "时光",
	}

	prefixSymbols := []string{
		"°", "·", "『", "【", "「", "〔", "꧁", "ꦿ", "࿐",
	}

	suffixSymbols := []string{
		"°", "·", "』", "】", "」", "〕", "꧂", "ꦿ", "࿐",
	}

	name := names[rand.Intn(len(names))]

	// 70%概率添加前缀符号
	if rand.Float32() < 0.7 {
		name = prefixSymbols[rand.Intn(len(prefixSymbols))] + name
	}

	// 70%概率添加后缀符号
	if rand.Float32() < 0.7 {
		name = name + suffixSymbols[rand.Intn(len(suffixSymbols))]
	}

	return name
}

func (g *DecoratedGenerator) GetStyle() string {
	return "符号装饰"
}

// 7. 诗意/意境风格
type PoeticGenerator struct{}

func (g *PoeticGenerator) Generate() string {
	templates := []func() string{
		// 四字成语风
		func() string {
			words1 := []string{"春", "夏", "秋", "冬", "晨", "暮", "朝", "夕"}
			words2 := []string{"风", "雨", "雪", "霜", "露", "雾", "烟", "云"}
			words3 := []string{"初", "细", "微", "轻", "薄", "淡", "浅", "深"}
			words4 := []string{"醉", "梦", "吟", "思", "念", "忆", "望", "叹"}
			return words1[rand.Intn(len(words1))] + words2[rand.Intn(len(words2))] +
				words3[rand.Intn(len(words3))] + words4[rand.Intn(len(words4))]
		},
		// 自然景物
		func() string {
			scenes := []string{
				"山间明月", "水上清风", "竹外疏梅", "柳下笛声",
				"花间一壶酒", "雨打芭蕉", "雪落无痕", "云淡风轻",
			}
			return scenes[rand.Intn(len(scenes))]
		},
	}

	return templates[rand.Intn(len(templates))]()
}

func (g *PoeticGenerator) GetStyle() string {
	return "诗意意境"
}

// 8. 二次元风格
type AnimeGenerator struct{}

func (g *AnimeGenerator) Generate() string {
	prefixes := []string{
		"次元", "二次元", "动漫", "ACG", "宅",
		"萌二", "元气", "中二", "病娇", "傲娇",
	}

	nouns := []string{
		"少女", "少年", "魔法师", "勇者", "骑士",
		"公主", "王子", "萝莉", "御姐", "正太",
	}

	suffixes := []string{
		"", "酱", "桑", "君", "殿", "sama", "chan",
	}

	return prefixes[rand.Intn(len(prefixes))] +
		nouns[rand.Intn(len(nouns))] +
		suffixes[rand.Intn(len(suffixes))]
}

func (g *AnimeGenerator) GetStyle() string {
	return "二次元"
}

// 网名生成器管理器
type UsernameGeneratorManager struct {
	generators []UsernameGenerator
}

func NewUsernameGeneratorManager() *UsernameGeneratorManager {
	return &UsernameGeneratorManager{
		generators: []UsernameGenerator{
			&CuteGenerator{},
			&AncientGenerator{},
			&CoolGenerator{},
			&FunnyGenerator{},
			&MixedGenerator{},
			&DecoratedGenerator{},
			&PoeticGenerator{},
			&AnimeGenerator{},
		},
	}
}

// 生成指定风格的网名
func (m *UsernameGeneratorManager) GenerateByStyle(styleIndex int) (string, string) {
	if styleIndex < 0 || styleIndex >= len(m.generators) {
		styleIndex = rand.Intn(len(m.generators))
	}
	generator := m.generators[styleIndex]
	return generator.Generate(), generator.GetStyle()
}

// 随机生成任意风格网名
func (m *UsernameGeneratorManager) GenerateRandom() (string, string) {
	generator := m.generators[rand.Intn(len(m.generators))]
	return generator.Generate(), generator.GetStyle()
}

// 批量生成网名
func (m *UsernameGeneratorManager) GenerateBatch(count int, styleIndex int) []string {
	usernames := make([]string, count)
	for i := 0; i < count; i++ {
		if styleIndex == -1 {
			username, _ := m.GenerateRandom()
			usernames[i] = username
		} else {
			username, _ := m.GenerateByStyle(styleIndex)
			usernames[i] = username
		}
	}
	return usernames
}

// 显示所有可用风格
func (m *UsernameGeneratorManager) ShowAllStyles() {
	fmt.Println("可用的网名风格：")
	for i, generator := range m.generators {
		fmt.Printf("%d. %s\n", i, generator.GetStyle())
	}
}
