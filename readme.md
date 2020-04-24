#golang製　書庫サーバー

使用ライブラリー
https://github.com/mattn/go-sqlite3

必要パッケージ

```
sudo apt install g++-arm-linux-gnueabihf
sudo apt install poppler-utils
```

コンパイル方法

armの場合
```
CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ \
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 \
go build -o bookserver
```

以下のコマンドを実行をすれば上記と同じコマンド
```
source arm_build
```

以下のフォルダが必要なもの
* 実行ファイル
* html
* html_tmp