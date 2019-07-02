# コンポーネント図

## コンポーネント一覧

```plantuml

@startuml
skinparam componentStyle uml2

component 権利利用者発行者
component サービス事業者
component 権利管理サーバ {
  database 権利情報
  component 権利管理機能
}
component SPTSM
component MNOTSM
component ユーザ端末 {
  component TSMプロキシ
  component UIアプリケーション
  component Applet
}

権利利用者発行者-サービス事業者
サービス事業者--権利管理機能
権利管理機能--SPTSM
SPTSM--TSMプロキシ
MNOTSM--TSMプロキシ
TSMプロキシ-UIアプリケーション
権利利用者発行者-UIアプリケーション
TSMプロキシ--Applet

@enduml

```

| コンポーネント名 | 機能概要 |
|---|---|
| 権利利用者/発行者 | 権利を利用するユーザ、権利を発行するユーザ |
| サービス事業者 | 権利発行、検索、利用予約のIFを権利利用者、発行者に提供する。 |
| 権利管理サーバ | 権利情報を一元管理する。 |
| SP-TSM |  |
| MNO-TSM |  |
| サービス事業者 |  |
| サービス事業者 |  |

## 権利発行

権利発行者がサービス事業者のAPIを経由して権利情報を登録する場合。

```plantuml
@startuml
skinparam componentStyle uml2

component 権利利用者発行者

component サービス事業者
サービス事業者 - () 001



component 権利管理サーバ {
  database 権利情報
  component 権利管理機能
  権利管理機能 - () 002
}

権利利用者発行者 ..> 001 : 権利情報登録
サービス事業者 ..> 002 : 権利情報登録
権利管理機能 --> 権利情報 : 権利情報登録

@enduml

```

権利発行者が自端末のUIアプリケーションのAPIを経由して権利情報を登録する場合。

```plantuml
@startuml
skinparam componentStyle uml2

component 権利利用者発行者
component ユーザ端末 {
  component TSMプロキシ
  TSMプロキシ - () 004
  component UIアプリケーション
  UIアプリケーション - () 003
}
component SPTSM
SPTSM - () 005
component 権利管理サーバ {
  database 権利情報
  component 権利管理機能
  権利管理機能 - () 002
}
権利利用者発行者 ..> 003 : 権利情報登録
UIアプリケーション ..> 004 : 権利情報登録
TSMプロキシ ..> 005 : 権利情報登録
SPTSM ..> 002 : 権利情報登録
権利管理機能 --> 権利情報 : 権利情報登録
@enduml

```

## IF一覧

| IF通番 | IF保有CP | IF利用CP | プロトコル | 用途 |
|---|---|---|---|---|
| 001 | サービス事業者 | 権利利用者発行者 | HTML(REST) | 権利発行,変更,削除,検索 |
| 002 | 権利管理サーバ | サービス事業者<br>SP TSM | HTML(REST) |  権利発行,変更,削除,検索 |
| 003 | UIアプリケーション | 権利利用者発行者 | GUI |  権利発行,検索,利用登録 |
| 004 | TSMプロキシ | UIアプリケーション | HTML(REST) |  HTML(REST)プロキシ |
| 005 | SP TSM | TSMプロキシ | HTML(REST) | 権利発行,変更,削除,検索  |

# シーケンス図
## 権利発行(サービス事業者経由)
- 処理概要
権利所有者がサービス事業者を経由して権利を発行する。

```plantuml
@startuml
skinparam componentStyle uml2
activate 権利発行者
権利発行者->サービス事業者 : 権利情報登録
  activate サービス事業者
  サービス事業者-> 権利管理機能 : 権利情報登録
    activate 権利管理機能
    rnote over 権利管理機能 #white
      通知された権利情報の
      内容をチェック
    end rnote
    alt チェックNG
      サービス事業者 <-- 権利管理機能 : 登録NG
    end
    rnote over 権利管理機能 #white
      通知された権利情報を
      権利情報Databaseへ登録
    end rnote
    サービス事業者 <-- 権利管理機能 : 登録OK
    deactivate 権利管理機能
    alt 登録NG
      権利発行者<--サービス事業者 : 登録NG
    end
  権利発行者<--サービス事業者 : 登録OK
  deactivate サービス事業者
deactivate 権利発行者
@enduml
```
