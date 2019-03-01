syzoj-ng
---

## 许可证
本软件以 [GNU Affero General Public License v3.0](https://choosealicense.com/licenses/agpl-3.0/) 发布，见 [LICENSE](LICENSE).

## 安装
syzoj-ng 采用前后端分离架构，前端在 <https://github.com/syzoj/syzoj-ng-app> 安装，后端在本仓库，两者互相独立，通过 nginx 使两者配合工作。

安装方法：
1. 安装 MongoDB、nodejs 和 nginx，具体操作取决于操作系统。
2. 下载 <https://github.com/syzoj/syzoj-ng-app>。依次执行 `npm install`，`npm install -g react-scripts` `npm run build` 编译前端。
3. 安装 [Go](https://golang.org) 语言，并设置 `GOPATH`，将 bin 目录加入到 `PATH`。具体步骤为将以下内容加入到 `~/.profile`：
```bash
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH
```
4. 执行 `go get github.com/syzoj/syzoj-ng-go`。此时 syzoj-ng-go 命令应该安装完成。
5. 创建一个目录保存相关文件。把 `$GOPATH/src/github.com/syzoj/syzoj-ng-go/config-example.json` 复制到当前目录并命名为 `config.json`，按需要调整配置。
6. 配置 nginx，根目录指向编译好的前端，`/api` 目录指向后端。此时网站已经可以访问。评测端还未完成，因此暂时无法评测。
