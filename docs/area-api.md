# 範囲検索 API

```
GET https://5ii04kds12.execute-api.ap-northeast-1.amazonaws.com/Prod/rate?lat=:latitude&lon:=longitude
```

For example.

```
GET https://5ii04kds12.execute-api.ap-northeast-1.amazonaws.com/Prod/rate?lon=139.7730717&lat=35.698383
```

## 概要

- クエリパラメータで指定したエリアの評価値と含まれる医療機関データを取得できます
- 有効回答数が 1 年（365 日）以上のデータのみが対象となります

## API 仕様

### エンドポイント

| エンドポイント | HTTP メソッド |
| -------------- | ------------- |
| /rate          | GET           |

### クエリパラメータ

| パラメータ | 説明            | 型    | 必須 | デフォルト値 |
| ---------- | --------------- | ----- | ---- | ------------ |
| lat        | 緯度(latitude)  | float | ◯    |              |
| lon        | 経度(longitude) | float | ◯    |              |
| distance   | 探索距離(km)    | float |      | 1.2          |

※1.2km=山手線の平均駅間距離

### レスポンス

| ステータスコード | 説明                                                                      |
| ---------------- | ------------------------------------------------------------------------- |
| 200              | JSON 形式でデータが返却される                                             |
| 400              | パラメータエラー（必須パラメーターの不足）                                |
| 404              | 指定範囲以内にデータが存在しないもしくは、存在するが有効回答数が 1 年未満 |

#### 正常パターンのレスポンス形式

| フィールド |              | 型     | 説明                                                           |
| ---------- | ------------ | ------ | -------------------------------------------------------------- |
| areaRate   |              | float  | 指定エリア全体の評価値                                         |
| facilities |              | Object | 指定エリアに含まれる医療機関の情報と集計データ                 |
|            | id           | string | 医療機関 ID                                                    |
|            | name         | string | 医療機関名                                                     |
|            | prefecture   | string | 都道府県                                                       |
|            | address      | string | 住所                                                           |
|            | latitude     | float  | 緯度                                                           |
|            | longitude    | float  | 経度                                                           |
|            | city         | string | 市区町村名                                                     |
|            | cityCode     | string | 地方公共団体コード                                             |
|            | validDays    | int    | 有効回答が得られた日数                                         |
|            | normalDays   | int    | 有効回答のうち、営業状況が「通常」の日数                       |
|            | limittedDays | int    | 有効回答のうち、営業状況が「制限」の日数                       |
|            | stoppedDays  | int    | 有効回答のうち、営業状況が「停止」の日数                       |
|            | rate         | float  | 営業状況の日数を元にした評価値                                 |
|            | facilityType | string | 医療機関種別（HOSPITAL:入院, EMERGENCY:救急, OUTPATIENT:外来） |
|            | distance     | float  | 指定エリアの中心からの距離（km）                               |
| geoJson    |              | Object | 地図情報サービス等で描画するための GeoJson データ              |

geoJson フィールドで取得できるデータの使い方については README を参照してください。

For example.

```
{
  "areaRate": 0.98,
  "facilities": [
    {
      "id": "xxxxxxx",
      "name": "XXXXXXXXXXXXXXX",
      "prefecture": "東京都",
      "address": "xxxxxxxxxxxxx",
      "latitude": xx.xxxxxx,
      "longitude": xxx.xxxxxx,
      "city": "千代田区",
      "cityCode": "131016",
      "validDays": 608,
      "normalDays": 585,
      "limittedDays": 23,
      "stoppedDays": 0,
      "rate": 0.98,
      "facilityType": "HOSPITAL",
      "distance": 0.396
    }
  ],
  "geoJson": {}
}
```
