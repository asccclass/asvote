# 「研究成果影像展」投票網站（正式區）

# # 網路投票系統
主要使用FB、Google 認證進行投票。

## 安裝
* 將根目錄下的 vote.db 複製到 data 目錄下
* sqlite 資料庫可以使用 [sqliteadmin](https://sqliteadmin.orbmu2k.de/)工具進行管理後，將vote.db上傳至伺服器

## Client 端

## Server 端提供的 API

# 系統安裝
* 先於 /etc/caddy/Caddyfile 增加下列 route 後，重啟 caddy 服務。

```
// 適用於：Caddy 0.11.1 (non-commercial use only)

proxy /vote localhost:10028 {
   without /vote
   websocket
   transparent
}
```

* 首次執行

```
make run
```

* 更版
請先更新 makefile 中的版本號，然後執行：

```
make re
```
