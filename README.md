# seaotterms-api

* 因為近期專案都準備前後端完全分離，所以開了一個api專案去集中API管理，減少伺服器需要開的PORT(專案)數量。
* 目前旗下有管理三個專案的API:
    * blog/seaotterms.com: 主站(上線中)
    * gal/gal.seaotterms.com: galgame文章資源分享站(開發中)
    * teach/teach.seaotterms.com: 教學文章站(暫時關閉中)

## 專案架構

### api模組
* 存放API方法，用**站台別**分子模組
### dto模組
* 存放DTO結構，用**站台別**分子模組
### middleware模組
* 存放middleware，用**站台別**分子模組
### model模組
* 存放資料表結構，用**資料庫別**分子模組  
### router模組
* 存放站台路由，用**站台別**分子模組
### utils模組
* 其他工具程式，用**站台別**分子模組

## 專案注意事項

* 因為後端使用Session處理，因為跨域無法設定Cookie問題，所以要確保前端跟後端框架要跑在同一個Domain下(可以不同Sub Domain)