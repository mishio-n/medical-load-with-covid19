# medical-load-with-covid19

コロナ期間の医療提供体制データを元にした、地域の医療体制の強度を可視化するサービス

## 使用したオープンデータ

https://corona.go.jp/dashboard より、**全国医療機関の医療提供体制の状況 オープンデータ**を使用。

- API
  - https://opendata.corona.go.jp/api/covid19DailySurvey
- レスポンスに含まれるデータ
  - データ提供元の医療機関情報
  - 医療機関の営業体制
    - 通常・制限・停止・未回答・設置なし

## アイデアとターゲット

- コロナウイルスの蔓延期間に稼働できている日数が多い＝パンデミックに耐えられる医療機関と言えるのでは？
  - 回答日数のうちの営業状況で指数化できそう
- 範囲指定して、エリア内の医療機関の指数があると、家探しの観点で需要あるかも
- ターゲットは引越し先を探している人と、不動産屋
  - 家探し用のアプリや、不動産屋の Web アプリに情報提供して使ってもらう想定をしています

## 技術的にやりたいこと

- 位置情報を用いた、物件の範囲検索
  - ElasticSearch の geo_distance クエリのようなイメージ
- 更新データ取り込みの日次バッチ作成
  - オープンデータの API レスポンスを DB に格納
- 蓄積データを分析するバッチ作成
  - DB を読み込み、分析して別テーブルへ格納

## テーブル設計

![](./prisma/ERD.png)

- 病院データテーブル
- API で取得できる日々の提出データ格納テーブル
- 提出データの分析結果格納テーブル

## 作成したプログラム

- バッチ
  - 日次で API から前日データを収集するバッチ
  - 日次で分析処理を実行するバッチ
  - 手動で API から過去データを収集するバッチ
- API
  1. 任意の地点（緯度経度）と距離（半径）を指定して、対象エリアの評価指数と描画用データを取得できる API
  2. 都道府県単位で、医療機関の分析データを取得できる API

## API 仕様書

1. [範囲検索 API](./docs/area-api.md)
2. [医療機関の分析データ取得 API](./docs/facility-api.md)

## 範囲検索 API の GeoJson データの使い方

GeoJson フォーマットに対応している地図情報サービスやライブラリで使用することで、  
地図上に検索結果を描画することができます。

### イメージ

[こちら](./demo.geojson)  
※GitHub の GeoJson プレビュー機能では詳細な表現はできないためイメージとなります。

### 簡単な使い方

1. 中心となる緯度経度と、半径になる距離(km)を決める
2. jq をインストールする

```bash
brew install jq
```

3. パラメータをクエリで指定してコマンドで API を実行する

```
curl 'https://5ii04kds12.execute-api.ap-northeast-1.amazonaws.com/Prod/rate?lon=139.7711&lat=35.6916&distance=2' | jq '.geoJson' | pbcopy
```

4. https://geojson.io/ にアクセスし、クリップボードの内容を画面右側のエディタ部分へ貼り付ける
