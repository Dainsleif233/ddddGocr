package main

import "github.com/Dainsleif233/ddddGocr"

func main() {
	// 使用示例
	result, err := ddddGocr.SlideMatch("test/bgd1.jpg", "test/bg1.png", "comparison", "opencv")
	if err != nil {
		print("滑块匹配失败: " + err.Error() + "\n")
		return
	}

	print(result.X1)
}
