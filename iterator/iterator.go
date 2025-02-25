package main

import (
	"fmt"
	"iter"
)

func main() {
	// numbers内の偶数だけが2倍になって取得
	for i := range double(even(numbers())) {
		fmt.Println(i)
	}
}

// 0~9の整数をPushするイテレータを返す関数
func numbers() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range 10 {
			yield(i)
		}
	}
}

// イテレータを受け取って偶数の場合だけPushするイテレータを返す関数
func even(seq iter.Seq[int]) iter.Seq[int] {
	// 偶数の場合だけPushするイテレータ
	return func(yield func(int) bool) {
		// 関数外から渡されたイテレータをrangeに渡して値をもらう
		for i := range seq {
			if i%2 == 0 {
				yield(i)
			}
		}
	}
}

// イテレータを受け取って要素を2倍にしてPushするイテレータを返す関数
func double(seq iter.Seq[int]) iter.Seq[int] {
	// 各要素を2倍にしてPushするイテレータ
	return func(yield func(int) bool) {
		// 関数外から渡されたイテレータをrangeに渡して値をもらう
		for i := range seq {
			yield(i * 2)
		}
	}
}
