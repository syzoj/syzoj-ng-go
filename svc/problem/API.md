# problem
管理题目文档的服务。

## GET /api/problem/id/:id
获取标号为 `id` 的题目。输出包含字段 `id`、`title`、`statement`、`tags`、`name`。

### 输入格式
```
GET /api/problem/id/m8V64e5OaF2aP1YW
```

### 输出格式
```
{
	"id": "m8V64e5OaF2aP1YW",
	"title": "标题",
	"statement": "markdown 题面",
	"tags": ["标签"],
	"name": "1"
}
```

## POST /api/problem
创建一道题目，支持字段 `title`、`statement`、`tags`、`name`。输出题目，仅支持字段 id。

### 输入格式
```
POST /api/problem
{
	"title": "标题",
	"statement": "markdown 题面",
	"tags": ["标签"],
	"name": "1" // 可选
}
```

### 输出格式
```
{
	"id": "m8V64e5OaF2aP1YW"
}
```

## PUT /api/problem/id/:id
修改一道题目，支持字段 `title`、`statement`、`tags`。没有输出。id 必须有效，否则状态码为 404。

### 输入格式
```
PUT /api/problem/id/:id
{
	"title": "标题",
	"statement": "markdown 题面",
	"tags": ["标签"]
}
```

## DELETE /api/problem/id/:id
删除一道题目。没有输出。id 必须有效，否则状态码为 404。

### 输入格式
```
DELETE /api/problem/id/:id
```

## POST /api/problem/search
搜索题目。支持字段 `id`、`title`、`tags`。

### 输入格式
```
POST /api/problem/search
{
	"query": "query string" // Lucene 查询
}
```

### 输出格式
```
{
	"hits": [
		{
			"id": "m8V64e5OaF2aP1YW",
			"title": "标题",
			"statement": "",
			"tags": ["标签"]
		}
	]
}
```
