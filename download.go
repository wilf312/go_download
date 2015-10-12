package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"
)


var (
    Leng int      = 150 // 連番 150 〜 001までダウンロード
    threadNum int = 5   // 同時ダウンロード数
    launchSec time.Duration = time.Second * 600 // 起動時間
    dir string    = "_download/" // 保存先
    dateStr string    = createTime( 0 ) // 20060102 今日から+引数日をフォーマットして返す
    formatStr string    = "http://hoge/%s_%03d.html"
)
func downloadFromUrl(url string) {
    // ファイル名の取得
    tokens := strings.Split(url, "/")
    fileName := tokens[len(tokens)-1]
    fmt.Println("Downloading", url, "to", fileName)


    output, err := os.Create(dir + fileName)
    if err != nil {
        fmt.Println("Error while creating", fileName, "-", err)
        return
    }
    defer output.Close()

    response, err := http.Get(url)
    if err != nil {
        fmt.Println("Error while downloading", url, "-", err)
        return
    } else if response.StatusCode != 200 {
        fmt.Println("Status Code not ok =>", response.StatusCode)
        return
    }
    defer response.Body.Close()


    n, err := io.Copy(output, response.Body)
    if err != nil {
        fmt.Println("Error while downloading", url, "-", err)
        return
    }

    fmt.Println(fileName, n / 1000000, "MB downloaded.")



}
func main() {

    fmt.Println(dateStr)

    // 保存先フォルダがなければディレクトリ作成
    if !isExist(dir) {
        if err := os.Mkdir(dir, 0777); err != nil {
            fmt.Println(err)
            return
        }
    }

    // ダウンロード実行goroutineの生成
    for cnt :=0; cnt<threadNum; cnt++ {
        go download(cnt);
    }

    // 起動時間
    time.Sleep(launchSec)
    fmt.Println("end")
}


func download(num int) {
    fmt.Println("Thread => ", num)
    Leng--
    // 0以下になれば終了
    if Leng <= 0 {
        return
    }

    // URLの設定
    url := fmt.Sprintf(formatStr, dateStr, Leng)
    downloadFromUrl(url)

    go download(num);
}

func createTime(num int) string {

    day := time.Now()
    day = day.AddDate(0, 0, num) // 1日前
    const layout = "20060102"
    return day.Format(layout)
}

func isExist( aDirName string ) bool {
    if _, err := os.Stat(aDirName); err != nil {
        return false
    } else {

        return true
    }
}
