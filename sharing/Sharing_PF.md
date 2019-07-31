# ユースケース図

## ユースケース一覧

```plantuml

@startuml
skinparam componentStyle uml2

actor 権利発行者
actor 権利利用者

rectangle 権利流通PF  {
  usecase 権利発行
  usecase 権利利用者変更
  usecase 権利分割
  usecase 権利削除
  usecase 利用者登録
  usecase 利用者診断
  usecase 価値評価
  usecase 発行者登録
  usecase トラブル解決
  usecase 自動決済
}

権利発行 <-- 権利発行者
権利削除 <-- 権利発行者
利用者診断 -> 権利発行者
権利利用者 --> 権利利用者変更 : 権利利用登録,権利譲渡
権利発行者 --> 権利分割
権利利用者 --> 利用者登録
利用者登録 <.. 利用者診断 : include
権利分割 <.. 権利発行 : include
権利分割 <.. 権利削除 : include
価値評価 .> 権利発行 : include
価値評価 -> 権利発行者
発行者登録 <- 権利発行者
トラブル解決 <-- 権利利用者
トラブル解決 <-- 権利発行者
権利利用者 --> 自動決済
@enduml

```

| ユースケース名 | 概要 |
|---|---|
| 権利発行 | 権利を所有するユーザ又はサービス事業者が権利情報を権利流通PFに登録することにより、<br>権利利用者が権利の利用登録をすることにより権利の利用を可能とする。 |
| 権利利用者変更 | 権利利用者が権利を利用可能とするために、権利の利用権の所有者を変更する。 |
| 権利分割 | 権利発行者が権利の分割利用を可能とするため、権利を分割する。<br>ex. 駐車場の1日単位貸し→駐車場15分単位へ変更 |
| 権利削除 | 権利流通PFに登録していた権利情報を削除する。 |
| 利用者登録 | 権利利用者の情報を権利流通PFに登録する。 |
| 利用者診断 | 権利利用者の情報から権利利用者がどのような性格傾向があるかを診断する。 |
| 価値評価 | 権利情報から権利利用権の妥当な基準価格を判断し提供する。 |
| 発行者登録 | 権利を発行するユーザ情報をを登録する。 |
| トラブル解決 | 権利行使時に発生したトラブルを自動的に解決する仕組み。 |
| 自動決済 | 権利の行使を確認後自動的に権利利用料金の決済を行う。 |

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

rectangle メモ #white [
  権利流通PFは権利管理サーバに該当する。
]

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

## 権利流通PF

```plantuml
@startuml
skinparam componentStyle uml2

component 権利情報管理機能
component 価格判定機能
component 利用者性格分析機能
component 権利利用者認証機能
component 権利情報バリデーション機能
component 権利発行者認証機能

@enduml

```

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
## 権利発行(UIアプリ経由)
- 処理概要
権利所有者がUIアプリを使用して権利を発行する。

```plantuml
@startuml
skinparam componentStyle uml2
activate 権利発行者
権利発行者 -> UIアプリケーション : 権利情報登録
  activate UIアプリケーション
    UIアプリケーション -> TSMプロキシ : 権利情報登録
    activate TSMプロキシ
      TSMプロキシ -> SP_TSM
      activate SP_TSM
        rnote over SP_TSM #white
          適切なアプリケーション
          からのアクセスか認証
        end rnote
        alt 認証NG
          TSMプロキシ <-- SP_TSM : 権利登録NG
        end
        SP_TSM -> 権利管理機能
        activate 権利管理機能
          rnote over 権利管理機能 #white
            通知された権利情報の
            内容をチェック
          end rnote
          alt チェックNG
            SP_TSM <-- 権利管理機能 : 登録NG
          end
          rnote over 権利管理機能 #white
            通知された権利情報を
            権利情報Databaseへ登録
          end rnote
        SP_TSM <-- 権利管理機能 : 登録OK
        deactivate 権利管理機能
        alt 権利登録NG
          TSMプロキシ <-- SP_TSM : 権利登録NG
        end
      TSMプロキシ <-- SP_TSM : 権利登録OK
      deactivate SP_TSM
      alt 権利登録NG
        UIアプリケーション <-- TSMプロキシ : 権利登録NG
      end
    UIアプリケーション <-- TSMプロキシ : 権利登録OK
    deactivate TSMプロキシ
    alt 権利登録NG
      権利発行者 <-- UIアプリケーション : 権利登録NG
    end
  権利発行者 <-- UIアプリケーション : 権利登録OK
  deactivate UIアプリケーション
deactivate 権利発行者
@enduml
```
## 権利行使

```plantuml
@startuml
skinparam componentStyle uml2

opt 権利利用開始
  activate RW端末
  RW端末 -> 権利利用者端末 : 利用権利認証処理開始
    activate 権利利用者端末
    権利利用者端末 -> RW端末 : 利用権利情報
    activate RW端末
      rnote over RW端末 #white
        利用権利情報を認証
      end rnote
      alt 認証NG
        権利利用者端末 <-- RW端末 : 権利認証NG
      end
      RW端末 -> 権利管理サーバ : 権利状態変更
      activate 権利管理サーバ
        rnote over 権利管理サーバ #white
          権利状態を権利行使中に変更
        end rnote
        RW端末 <-- 権利管理サーバ
      deactivate 権利管理サーバ
      権利利用者端末 <-- RW端末 : 権利利用OK
    deactivate RW端末
    RW端末 <-- 権利利用者端末
  deactivate 権利利用者端末
end

opt 権利利用終了
  activate RW端末
  RW端末 -> 権利利用者端末 : 利用権利認証処理開始
    activate 権利利用者端末
    権利利用者端末 -> RW端末 : 利用権利情報
    activate RW端末
      rnote over RW端末 #white
        利用権利情報を認証
      end rnote
      alt 認証NG
        権利利用者端末 <-- RW端末 : 権利認証NG
      end
      RW端末 -> 権利管理サーバ : 権利状態変更
      activate 権利管理サーバ
        rnote over 権利管理サーバ #white
          権利の状態を未行使に変更
        end rnote
        RW端末 <-- 権利管理サーバ
      deactivate 権利管理サーバ
      権利利用者端末 <-- RW端末 : 権利利用OK
    deactivate RW端末
    RW端末 <-- 権利利用者端末
  deactivate 権利利用者端末
end

@enduml
```

## 検討事項

1. 権利情報のマスターをUICCに持つかそれとも権利管理サーバに持つか。

| 方式 | 機能性 | 信頼性 | 使用性 | 効率性 | 保守性 | 移植性 |
|---|---|---|---|---|---|---|
| UICCに権利情報保持 | ○ | △ 端末故障、電源断時に<br>機能が満たせない  | ○ | △ 端末主導の処理は<br>処理性能が端末依存になる。<br>使えるリソースも少ない。 | × 権利管理PFのソフトウェアに<br>変更が入った場合、すべての端末で<br>UPDATEが必要。  | ○ |
| 権利管理サーバに権利情報保持 | ○ | ○ | ○ | ○ | ○ | ○ |

- 機能性
ソフトウェアを指定された条件のもとで動作するとき、要求されている仕様を満たす能力のこと。
- 信頼性
ソフトウェアを指定された条件のもとで動作するとき、達成水準を維持し続ける能力のこと。誤作動時の復旧や、障害に対する許容性をあらわす場合もある。
- 使用性
ソフトウェアを指定された条件のもとで動作するとき、利用者が理解、習得、利用がスムーズにおこなえる能力こと。いわゆる「使い勝手」や「使いやすさ」、「操作性」のこと。
- 効率性
与えられたリソースに対して、適切な性能を発揮する能力のこと。たとえば、決められた処理時間の中でいかに早く、数多くの処理ができるか、などがあります。
- 保守性
できたソフトウェアの修正のしやすさの能力のこと。作った本人しか理解できないプログラムでは、改修が発生した際に多くのコストがかかってしまいます。これは利用者には直接は関係しない特性のように見えますが、最終的なサービスリリースまでにかかるコストの軽減は、利用者へのメリットにつながることが多くあります。
- 移植性
別な環境へ移すことになった際に、容易に移せる能力のこと。サーバーの移行や、使うフレームワークが変更になった場合などに重要になってくる。
