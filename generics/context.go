

//context パッケージは、ゴルーチンのキャンセルやタイムアウト管理、リクエスト単位の情報共有などを行うために使われる非常に重要な仕組み

//### Contextの主な特徴

1. **キャンセル可能**  
   親のContextがキャンセル（もしくはタイムアウトなど）されると、その子孫のContextも全てキャンセルされるという伝播の仕組みがある。

2. **値の受け渡しが可能**  
   `context.Context` には、リクエストスコープなどで必要となる値を保存して、関数間で受け渡す仕組みがある。  
   - ただし、値の受け渡しは最小限にし、設定する値はリクエストIDや認証情報など必要最小限に抑えるのが推奨。

3. **標準ライブラリやGoのベストプラクティス**  
   Goの標準ライブラリの多くの関数は、第一引数に `ctx context.Context` を取るようになってきている。  
   - 例えば `net/http` やデータベースの操作など、多くのI/O系の処理で `Context` が利用できる形になっている。


## まとめ

- `context` パッケージは、**キャンセル可能なゴルーチン管理とリクエストスコープの値の伝達を一元的に扱える**ツール。  
- `Background()` や `TODO()` から始まり、必要に応じて `WithCancel`, `WithDeadline`, `WithTimeout` で派生。  
- `ctx.Done()` を使ってキャンセル通知を受け取るのが基本パターン。  
- 値の保存は**最小限かつ適切なキー管理**のもとで使う。  
- Goの並行処理でよくある「どこで止めたらいいか分からない」「タイムアウト機構をどう作るか」という悩みを解決する。

これらを踏まえると、複数のゴルーチンをまたいで処理するときや、
リクエストごとに明確なタイムアウトを設定したいときに非常に便利です。
Webサーバーはもちろん、CLIツールなどでも、大きな並列処理を行う場合には 
`context` を活用することで**一貫性のあるエラーハンドリングとリソースの管理**が可能になります。

------

Go言語において `context` パッケージは、**ゴルーチンのキャンセルやタイムアウト管理、リクエスト単位の情報共有**などを行うために使われる非常に重要な仕組みです。以下では、主な使い方とメリット・注意点などをまとめます。

---

## contextパッケージの概要

`context` パッケージは、元々サーバーサイドなどで「リクエストがキャンセルされたらその処理に関連するすべてのゴルーチンを止めたい」「タイムアウトを設定したい」といったニーズを満たすために導入されました。  
基本的には**キャンセルやタイムアウトの合図を伝播させたり、リクエスト単位のスコープで値を管理したり**するために使われます。



---

## Contextの生成方法

### 1. `context.Background()`
```go
ctx := context.Background()
```
- 空のContextを返す。  
- アプリケーションのルートやメイン関数での起点として使われることが多い。

### 2. `context.TODO()`
```go
ctx := context.TODO()
```
- まだどんなContextが必要なのか決まっていない場合や、将来的にContextが必要になるがとりあえずplaceholderとして使いたいときに利用。  
- 本番コードでは最終的に `Background()` に置き換えることが推奨されることが多い。

---

## Contextの派生（With系関数）

### 1. `context.WithCancel(parentContext)`
```go
ctx, cancel := context.WithCancel(parentCtx)
defer cancel()
```
- 親のContext（`parentCtx`）から新しいContextを生成し、キャンセル関数(`cancel`)を取得する。  
- `cancel()` を呼び出すと、このContextおよびその子孫のContextがキャンセルされる。

### 2. `context.WithTimeout(parentContext, timeout)`
```go
ctx, cancel := context.WithTimeout(parentCtx, 3*time.Second)
defer cancel()
```
- 指定したタイムアウトが経過すると自動的にキャンセルされるContextを生成する。  
- 生成されたContextがキャンセルされるか、タイムアウトまでの間だけ有効。  
- 処理が終わったら `cancel()` を呼んでリソースを解放するのがベストプラクティス。

### 3. `context.WithDeadline(parentContext, deadline)`
```go
deadline := time.Now().Add(3 * time.Second)
ctx, cancel := context.WithDeadline(parentCtx, deadline)
defer cancel()
```
- タイムアウトではなく、指定した**時刻**（Deadline）を過ぎるとキャンセルされるContextを生成する。  
- 中身は `WithTimeout` とほぼ同じような使い方。時刻指定でのキャンセルが必要な場合に利用。

---

## Contextのキャンセルを受け取る側の処理

Contextを受け取る側では、しばしば `select` 文や `Done()` チャネルを利用して、キャンセル（もしくはタイムアウト）が通知されたかどうかを判定します。

### シンプルな例

```go
func doSomething(ctx context.Context) error {
    select {
    case <-time.After(2 * time.Second):
        // 時間のかかる処理の代わり
        return nil
    case <-ctx.Done():
        // 親のキャンセルやタイムアウトが発生
        return ctx.Err()
    }
}
```

- `ctx.Done()` はContextがキャンセルされたタイミングで閉じられるチャネルを返す。  
- 受け取るとキャンセルを検知できる。

---

## Contextで値をやり取りする

`context` には値を格納できる関数 `context.WithValue(parentContext, key, val)` が用意されています。ただし公式ドキュメントにもあるように、**値の格納には最小限に留める** ことが推奨されています。

```go
type keyType string

func doSomething(ctx context.Context) {
    userID := ctx.Value(keyType("user_id"))
    if userID != nil {
        fmt.Println("UserID:", userID)
    }
}

func main() {
    ctx := context.Background()
    ctx = context.WithValue(ctx, keyType("user_id"), "12345")
    
    doSomething(ctx)
}
```

- 利用する際には、キー用の型を定義したり、グローバルな変数を使わないなど、衝突を避けるための工夫が必要。
- **大量のデータやコンテキストと無関係な情報を詰め込むのは避ける**べき。

---

## メリット

1. **ゴルーチンのリークを防ぐ**  
   キャンセルやタイムアウトを伝播させることで、不要になったゴルーチンを早期に停止させる。これによりリソース浪費を防ぐ。

2. **統一されたキャンセル/タイムアウト管理**  
   `select <-ctx.Done()` というシンプルな構造で、複数の関数・複数のゴルーチンに一貫したキャンセルの仕組みを適用できる。

3. **リクエストスコープの値受け渡し**  
   Webアプリケーションなどで1リクエストに対応する値（認証情報、トレーシングIDなど）を、呼び出しチェーンの奥深くまで簡単に運ぶことができる。

4. **標準ライブラリとの親和性**  
   Goの標準ライブラリ ( `net/http`, `database/sql`, `os/exec` など ) でも `Context` を使ったインターフェースが増えており、対応することで書き方が整う。

---

## 注意点

1. **Contextは値の格納に乱用しない**  
   - Contextにデータを詰め込みすぎると、可読性が下がったり、意図しないデータ競合が起こり得る。

2. **関数の引数は必ず `ctx context.Context` を先頭に**  
   - Goの慣習的に、Contextを受け取る関数は第一引数に `ctx` を持たせる。

3. **キャンセル関数は必ず呼ぶ**  
   - `WithCancel`, `WithTimeout`, `WithDeadline` を使った場合は、処理後の `cancel()` 呼び出しでリソース解放するのが良い（特にタイムアウトやデッドラインの場合は重要）。

4. **nil の Context は使わない**  
   - `context.Background()` や `context.TODO()` を使うのが安全。

5. **Contextはスレッドセーフだが、書き換えはしない**  
   - 一度生成したContextを変更するのではなく、必ず新しいContextを生成して受け渡す。

---

