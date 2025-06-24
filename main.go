package main

import (
	"github.com/Dainsleif233/ddddGocr/ddddgocr"
	"github.com/Dainsleif233/ddddGocr/ddddgocr/withopencv"
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

func SlideComparison(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideComparisonResult, error) {
	return ddddgocr.SlideComparison(targetImageData, backgroundImageData)
}

func SlideComparisonWithPath(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideComparisonResult, error) {
	return ddddgocr.SlideComparisonWithPath(targetImagePath, backgroundImagePath)
}

func SlideMatchUseOpenCV(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return withopencv.SlideMatch(targetImageData, backgroundImageData)
}

func SlideMatchWithPathUseOpenCV(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideBBox, error) {
	return withopencv.SlideMatchWithPath(targetImagePath, backgroundImagePath)
}

func SimpleSlideMatchUseOpenCV(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return withopencv.SimpleSlideMatch(targetImageData, backgroundImageData)
}

func SimpleSlideMatchWithPathUseOpenCV(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideBBox, error) {
	return withopencv.SimpleSlideMatchWithPath(targetImagePath, backgroundImagePath)
}

func EnhancedSlideMatchUseOpenCV(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return withopencv.EnhancedSlideMatch(targetImageData, backgroundImageData)
}

func EnhancedSlideMatchWithPathUseOpenCV(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideBBox, error) {
	return withopencv.EnhancedSlideMatchWithPath(targetImagePath, backgroundImagePath)
}

func SlideComparisonUseOpenCV(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideComparisonResult, error) {
	return withopencv.SlideComparison(targetImageData, backgroundImageData)
}

func SlideComparisonWithPathUseOpenCV(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideComparisonResult, error) {
	return withopencv.SlideComparisonWithPath(targetImagePath, backgroundImagePath)
}

func main() {
	// 使用示例
	// result, err := withopencv.EnhancedSlideMatchWithPath("test/tgt1.png", "test/bgd1.jpg")
	// result, err := EnhancedSlideMatchWithPathUseOpenCV("test/tgt2.png", "test/bgd2.png")
	// if err != nil {
	// 	print("滑块匹配失败: " + err.Error() + "\n")
	// 	return
	// }

	// print(result.X1)
}
