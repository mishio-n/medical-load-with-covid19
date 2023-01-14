# 医療機関データ取得 API

```
GET https://5ii04kds12.execute-api.ap-northeast-1.amazonaws.com/Prod/facilities?prefecture=:都道府県名
```

For example.

```
GET https://5ii04kds12.execute-api.ap-northeast-1.amazonaws.com/Prod/facilities?prefecture=%E6%9D%B1%E4%BA%AC%E9%83%BD
```

## 概要

- クエリパラメータで指定した条件に合致する集計データを取得できます

## API 仕様

### エンドポイント

| エンドポイント | HTTP メソッド |
| -------------- | ------------- |
| /facilities    | GET           |

### クエリパラメータ

以下のクエリパラメータを指定してデータを取得可能です

| パラメータ | 説明                                   | 型                                            | 必須 | デフォルト値 |
| ---------- | -------------------------------------- | --------------------------------------------- | ---- | ------------ |
| prefecture | 都道府県名（**現状は「東京都」のみ**） | string                                        | ◯    |              |
| cityCode   | 全国地方公共団体コード                 | string                                        |      |              |
| type       | フィルタしたいん医療機関種別           | string(`HOSPITAL`, `EMERGENCY`, `OUTPATIENT`) |      |              |

### レスポンス

| ステータスコード | 説明                                       |
| ---------------- | ------------------------------------------ |
| 200              | JSON 形式でデータが返却される              |
| 400              | パラメータエラー（必須パラメーターの不足） |
| 404              | 条件に合致するデータが存在しない           |

#### 正常パターンのレスポンス形式

※ルート要素は配列となります

| フィールド   | 型     | 説明                                                           |
| ------------ | ------ | -------------------------------------------------------------- |
| id           | string | 医療機関 ID                                                    |
| name         | string | 医療機関名                                                     |
| prefecture   | string | 都道府県                                                       |
| address      | string | 住所                                                           |
| latitude     | float  | 緯度                                                           |
| longitude    | float  | 経度                                                           |
| city         | string | 市区町村名                                                     |
| cityCode     | string | 地方公共団体コード                                             |
| validDays    | int    | 有効回答が得られた日数                                         |
| normalDays   | int    | 有効回答のうち、営業状況が「通常」の日数                       |
| limittedDays | int    | 有効回答のうち、営業状況が「制限」の日数                       |
| stoppedDays  | int    | 有効回答のうち、営業状況が「停止」の日数                       |
| rate         | float  | 営業状況の日数を元にした評価値                                 |
| facilityType | string | 医療機関種別（HOSPITAL:入院, EMERGENCY:救急, OUTPATIENT:外来） |

For example.

```
[
  {
    "id": "XXXXXXX",
    "name": "XXXXXXXXXXXXXXXXXXXXXXXXX",
    "prefecture": "東京都",
    "address": "XXXXXXXXXXXXXXXXXX",
    "latitude": xx.xxxxxx,
    "longitude": xxx.xxxxxx,
    "city": "千代田区",
    "cityCode": "131016",
    "validDays": 388,
    "normalDays": 0,
    "limittedDays": 388,
    "stoppedDays": 0,
    "rate": 0.3,
    "facilityType": "HOSPITAL"
  }
]
```
