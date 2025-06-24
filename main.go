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

func SlideComparison(targetImageData, backgroundImageData []byte) (*struct {
	First  uint32
	Second uint32
}, error) {
	return ddddgocr.SlideComparison(targetImageData, backgroundImageData)
}

func SlideComparisonWithPath(targetImagePath, backgroundImagePath string) (*struct {
	First  uint32
	Second uint32
}, error) {
	return ddddgocr.SlideComparisonWithPath(targetImagePath, backgroundImagePath)
}

func main() {
	// 使用示例
	// result, err := SlideComparisonWithPath("test/bgd1.jpg", "test/bg1.png")
	// if err != nil {
	// 	fmt.Printf("滑块匹配失败: %v\n", err)
	// 	return
	// }

	// fmt.Printf("匹配结果: X1=%d\n",
	// 	result.First)
}
