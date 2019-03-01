# API
本文介绍后端提供的 API 接口，方便前端开发。

在 [server.go](app/api/server.go) 中的 `setupRoutes` 函数中包含了所有开放的 API 接口。

在 [model](app/model) 中的 `*.proto` 文件包含了要用到的数据的格式。其中 `api.proto` 包含了各个 API 的请求及响应的格式，`primitive.proto` 包含了 ObjectID 的格式，`model.proto` 仅表示数据库的格式，在前端开发时不会用到。

后端所有 API 接口会接受 JSON 格式的请求，返回 JSON 格式的响应，字段名称和 proto 文件中的名称对应。注意虽然 ObjectID 是一个 message，但在 JSON 中会自动转换成一个字符串而不是一个 object。
