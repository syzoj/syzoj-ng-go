package divine

import (
	"time"
	"hash/crc32"
	"math/rand"
)

type Divine struct {
	Fortune string `json:"fortune"`
	Good []*Item
	Bad []*Item
}

type Item struct {
	Title string `json:"title"`
	Detail string `json:"detail"`
}

type item struct {
	Title string
	Detail [2]string
}

var items = []*item{
	&item{Title: "刷题", Detail: [2]string{"一遍过样例", ""}},
	&item{Title: "装弱", Detail: [2]string{"我好菜啊", "你太强了"}},
	&item{Title: "搞x", Detail: [2]string{"爱上学习", "会被掰弯"}},
	&item{Title: "直播写代码", Detail: [2]string{"月入百万", "CE, RE and T, 身败名裂"}},
	&item{Title: "学数论", Detail: [2]string{"思维敏捷", "咋看都不会"}},
	&item{Title: "参加模拟赛", Detail: [2]string{"AK 虐场", "爆零"}},
}

func DoDivine(name string, sex int) *Divine {
	res := &Divine{Good: []*Item{}, Bad: []*Item{}}
	seed := crc32.ChecksumIEEE([]byte(name + time.Now().Format("20060102")))
	random := rand.New(rand.NewSource(int64(seed)))

	f := random.Float32()
	switch {
	case f <= 0.25:
		res.Fortune = "大吉"
	case f <= 0.5:
		res.Fortune = "大凶"
	case f <= 0.6:
		res.Fortune = "中平"
	case f <= 0.7:
		res.Fortune = "小吉"
	case f <= 0.8:
		res.Fortune = "小凶"
	case f <= 0.9:
		res.Fortune = "吉"
	default:
		res.Fortune = "凶"
	}

	pitems := make([]*item, len(items))
	copy(pitems, items)
	makeItem := func(typ int) *Item {
		for {
			id := random.Intn(len(pitems))
			item := pitems[id]
			if item == nil {
				continue
			}
			citem := &Item{Title: item.Title}
			if citem.Title == "搞x" {
				switch sex {
				case 0:
					citem.Title = "搞基"
				case 1:
					citem.Title = "搞姬"
				}
			}
			if item.Detail[typ] == "" {
				continue
			} else {
				citem.Detail = item.Detail[typ]
			}
			pitems[id] = nil
			return citem
		}
	}
	if res.Fortune != "大凶" {
		res.Good = []*Item{makeItem(0), makeItem(0)}
	}
	if res.Fortune != "大吉" {
		res.Bad = []*Item{makeItem(1), makeItem(1)}
	}
	return res
}
