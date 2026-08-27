# seaotterms-api

* 因為近期專案都準備前後端完全分離，所以開了一個api專案去集中API管理，減少伺服器需要開的PORT(專案)數量。
* 目前旗下有管理兩個專案的API:
    * blog: 主站(上線中)
    * teach: 教學文章站(暫時關閉中)

> [!NOTICE]
> `gal` 已下線，相關 API 不再維護。

## 專案架構

### blog / teach 套件
* 依**站台別**拆分，各自包含該站台的路由、handler、middleware 等邏輯
* 根路徑分別為 `/api/blog`、`/api/teach`，並依連線的資料庫自動註冊對應路由

### seaotterms-db 依賴
* 資料表結構、migration、DB 初始化已抽離至獨立套件 [`seaotterms-db`](https://github.com/peter910820/seaotterms-db)
* 本專案透過 `go.mod` 的 `replace` 指向本地路徑（預設 `../seaotterms-db`）

## 專案注意事項

* 因為後端使用Session處理，因為跨域無法設定Cookie問題，所以要確保前端跟後端框架要跑在同一個Domain下(可以不同Sub Domain)
