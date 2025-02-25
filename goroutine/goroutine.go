package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {

	// 現在の作業ディレクトリを取得
	rootdir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}
	fmt.Println("Current working directory (rootdir):", rootdir)

	// 実行
	// マルチスレッド・ゴルーチン
	errMul := RemoveMulti(rootdir)
	if err != nil {
		fmt.Println("Error:", errMul)
	} else {
		fmt.Println("マルチスレッド・ゴルーチンが正常に終了しました")
	}
	// シングルスレッド・ゴルーチン
	errSin := RemoveSingle(rootdir)
	if err != nil {
		fmt.Println("Error:", errSin)
	} else {
		fmt.Println("シングルスレッド・ゴルーチンが正常に終了しました")
	}

}

// マルチスレッド・ゴルーチン
// 最終更新日時が24時間より過去という条件下でファイルを削除する
func RemoveMulti(rootdir string) error {

	// 存在しない場合は処理を終了
	if _, err := os.Stat(rootdir); os.IsNotExist(err) {
		logrus.Tracef("Directory not found. rootdir: %s", rootdir)
		return nil
	}

	// 排他制御ファイルをチェック
	lockFile := rootdir + "/.lock"
	fi, flerr := os.Stat(lockFile)
	if flerr == nil {
		if fi.ModTime().Before(time.Now().Add(-24 * time.Hour)) {
			// 排他制御ファイルが24時間以上前であれば削除
			if err := os.Remove(lockFile); err != nil {
				logrus.WithError(err).Error("Failed to remove lock file.")
				return err
			}
		} else {
			// ロックファイルが存在し24時間以内であれば処理を終了
			logrus.Warningf("Lock file remains within the last 24 hours. lockFile: %s", lockFile)
			return nil
		}
	} else if !os.IsNotExist(flerr) {
		// ロックファイルのステータス確認に失敗した場合
		logrus.Errorf("Failed to check lock file. lockFile: %s", lockFile)
		return flerr
	}

	// 排他制御ファイルを作成
	lockFileHandle, err := os.Create(lockFile)
	if err != nil {
		logrus.Errorf("Failed to create lock file. lockFile: %s, err: %v", lockFile, err)
		return err
	}
	lockFileHandle.Close()

	// 関数終了時にロックファイルを削除
	defer func() {
		if err := os.Remove(lockFile); err != nil {
			logrus.Errorf("Failed to remove lock file. lockFile: %s, err: %v", lockFile, err)
		} else {
			logrus.Tracef("Executed remove lock file. lockFile: %s", lockFile)
		}
	}()

	// ファイル削除を並列処理で実行するためのチャネル
	type fileTask struct {
		path    string
		modTime time.Time
	}
	taskChan := make(chan fileTask, 300)

	// エラーチャネル
	errorChan := make(chan error, 1)

	// ワーカーを起動
	var wg sync.WaitGroup
	const numWorkers = 2
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				// MPDファイルが24時間より前のものかどうかをチェック
				if task.modTime.Before(time.Now().Add(-24 * time.Hour)) {
					// ファイル削除実行
					if err := os.Remove(task.path); err != nil {
						logrus.Warningf("Failed to remove file. path: %s, err: %v", task.path, err)
						errorChan <- err
					} else {
						logrus.Tracef("Successfully removed file. path: %s", task.path)
					}
				}
			}
		}()
	}

	// ファイルを探索しタスクをチャネルに送る
	go func() {
		defer close(taskChan)
		err := filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// ファイル名がパターンに一致するかどうかをチェック
			if matched, _ := filepath.Match("sess_*", info.Name()); matched {
				taskChan <- fileTask{path: path, modTime: info.ModTime()}
			}
			return nil
		})
		if err != nil {
			logrus.Errorf("Failed to walk the path. mpddir: %s, err: %v", rootdir, err)
			errorChan <- err
		}
	}()

	// ワーカーが全て終了するのを待つ
	wg.Wait()
	close(errorChan)

	// エラーが発生していた場合最初のエラーを返す
	if err := <-errorChan; err != nil {
		return err
	}

	return nil
}

////////////////////
////////////////////

// シングルスレッド・ゴルーチン
// (最終更新日時が24時間より過去)を削除する
func RemoveSingle(rootdir string) error {

	// ディレクトリが存在しない場合は処理を終了
	if _, err := os.Stat(rootdir); os.IsNotExist(err) {
		logrus.Tracef("Directory not found. rootdir: %s", rootdir)
		return nil
	}

	// 排他制御ファイルをチェック
	lockFile := rootdir + "/.lock"
	fi, flerr := os.Stat(lockFile)
	if flerr == nil {
		// 排他制御ファイルが24時間以上前であれば削除
		if fi.ModTime().Before(time.Now().Add(-24 * time.Hour)) {
			if err := os.Remove(lockFile); err != nil {
				logrus.Errorf("Failed to remove lock file. lockFile: %s, err: %v", lockFile, err)
				return err
			}
			logrus.Tracef("Removed Lock file. lockFile: %s", lockFile)
		} else {
			// ロックファイルが存在し24時間以内であれば処理を終了
			logrus.Tracef("Lock file remains within the last 24 hours. lockFile: %s", lockFile)
			return nil
		}
	} else if !os.IsNotExist(flerr) {
		// ロックファイルのステータス確認に失敗した場合
		logrus.Errorf("Failed to check lock file. lockFile: %s, flerr: %v", lockFile, flerr)
		return flerr
	}

	// 排他制御ファイルを作成
	lockFileHandle, err := os.Create(lockFile)
	if err != nil {
		logrus.Errorf("Failed to create lock file. lockFile: %s, err: %v", lockFile, err)
		return err
	}
	lockFileHandle.Close()

	// ゴルーチン内でファイルを探索し削除する
	go func() {
		logrus.Tracef("Started file cleanup in goroutine. rootdir: %s", rootdir)
		// ロックファイルを削除
		defer func() {
			if err := os.Remove(lockFile); err != nil {
				logrus.Errorf("Failed to remove lock file. lockFile: %s, err: %v", lockFile, err)
			} else {
				logrus.Tracef("Executed remove lock file. lockFile: %s", lockFile)
			}
		}()
		// ファイルを探索し削除する
		walkErr := filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logrus.Warningf("Failed to process file or directory. rootdir: %s, err: %v", rootdir, err)
				return err
			}

			// ファイル名がパターンに一致するかどうかをチェック
			if matched, _ := filepath.Match("sess_*", info.Name()); matched {
				// ファイルの最終変更日時が24時間より前かどうかをチェック
				if info.ModTime().Before(time.Now().Add(-24 * time.Hour)) {
					// ファイル削除実行
					if err := os.Remove(path); err != nil {
						logrus.Warningf("Failed to remove file. path: %s, err: %v", path, err)
					} else {
						logrus.Debugf("Successfully removed file. path: %s", path)
					}
				}
			}
			return nil
		})
		if walkErr != nil {
			logrus.Warningf("Failed to walk the path. rootdir: %s, walkErr: %v", rootdir, walkErr)
		}
	}()

	return nil
}

////////////////////
////////////////////
