package main

import (
	"encoding/json"
	"fmt"

	"github.com/syzoj/syzoj-ng-go/lib/divine"
)

type items struct {
	Items []*divine.Item `json:"items"`
}

var config = `
{
  "items": [
    {
      "title": "刷题",
      "detail": [
        "一遍过样例",
        null
      ]
    },
    {
      "title": "装弱",
      "detail": [
        "我好菜啊",
        "你太强了"
      ]
    },
    {
      "title": {
        "boy": "搞基",
        "girl": "搞姬" 
      },
      "detail": [
        "爱上学习",
        "会被掰弯"
      ]
    },
    {
      "title": "直播写代码",
      "detail": [
        "月入百万",
        "CE, RE and T，身败名裂"
      ]
    },
    {
      "title": "学数论",
      "detail": [
        "思维敏捷",
        "咋看都不会"
      ]
    },
    {
      "title": "参加模拟赛",
      "detail": [
        "AK 虐场",
        "爆零"
      ]
    }
  ]
}
`

func main() {
	var val items
	if err := json.Unmarshal([]byte(config), &val); err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", val.Items[0])
}
