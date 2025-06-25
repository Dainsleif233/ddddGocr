package main

import (
	"github.com/Dainsleif233/ddddGocr/ddddgocr"
	"github.com/Dainsleif233/ddddGocr/ddddgocr/withopencv"
)

func SlideMatch(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return ddddgocr.SlideMatch(targetImageData, backgroundImageData)
}

func SimpleSlideMatch(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return ddddgocr.SimpleSlideMatch(targetImageData, backgroundImageData)
}

func EnhancedSlideMatch(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return ddddgocr.EnhancedSlideMatch(targetImageData, backgroundImageData)
}

func SlideComparison(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideComparisonResult, error) {
	return ddddgocr.SlideComparison(targetImageData, backgroundImageData)
}

func SlideMatchUseOpenCV(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return withopencv.SlideMatch(targetImageData, backgroundImageData)
}

func SimpleSlideMatchUseOpenCV(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return withopencv.SimpleSlideMatch(targetImageData, backgroundImageData)
}

func EnhancedSlideMatchUseOpenCV(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
	return withopencv.EnhancedSlideMatch(targetImageData, backgroundImageData)
}

func SlideComparisonUseOpenCV(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideComparisonResult, error) {
	return withopencv.SlideComparison(targetImageData, backgroundImageData)
}

func main() {

	// 增强版滑块匹配的文件路径版本
	// func EnhancedSlideMatchWithPath(targetImagePath, backgroundImagePath string) (*SlideBBox, error) {
	// 	targetData, err := os.ReadFile(targetImagePath)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("读取目标图像文件失败: %v", err)
	// 	}

	// 	backgroundData, err := os.ReadFile(backgroundImagePath)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("读取背景图像文件失败: %v", err)
	// 	}

	// 	return EnhancedSlideMatch(targetData, backgroundData)
	// }

	// 使用示例
	// result, err := withopencv.EnhancedSlideMatchWithPath("test/tgt1.png", "test/bgd1.jpg")
	// result, err := EnhancedSlideMatchWithPathUseOpenCV("test/tgt2.png", "test/bgd2.png")
	// if err != nil {
	// 	print("滑块匹配失败: " + err.Error() + "\n")
	// 	return
	// }

	// print(result.X1)
}
