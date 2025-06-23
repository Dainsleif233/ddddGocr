package main

import (
	"github.com/Dainsleif233/ddddGocr/ddddgocr"
)

func SlideMatch(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return ddddgocr.SlideMatch(targetImageData, backgroundImageData)
}

func SlideMatchWithPath(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideBBox, error) {
	return ddddgocr.SlideMatchWithPath(targetImagePath, backgroundImagePath)
}

func SimpleSlideMatch(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return ddddgocr.SimpleSlideMatch(targetImageData, backgroundImageData)
}

func SimpleSlideMatchWithPath(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideBBox, error) {
	return ddddgocr.SimpleSlideMatchWithPath(targetImagePath, backgroundImagePath)
}

func EnhancedSlideMatch(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return ddddgocr.EnhancedSlideMatch(targetImageData, backgroundImageData)
}

func EnhancedSlideMatchWithPath(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideBBox, error) {
	return ddddgocr.EnhancedSlideMatchWithPath(targetImagePath, backgroundImagePath)
}

func main() {
	// 使用示例
	// result, err := EnhancedSlideMatchWithPath("test/tgt1.png", "test/bgd1.png")
	// if err != nil {
	// 	fmt.Printf("滑块匹配失败: %v\n", err)
	// 	return
	// }

	// fmt.Printf("匹配结果: X1=%d, Y1=%d, X2=%d, Y2=%d, TargetY=%d\n",
	// 	result.X1, result.Y1, result.X2, result.Y2, result.TargetY)
}
