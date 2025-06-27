package ddddGocr

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/Dainsleif233/ddddGocr/ddddgocr"
)

type MatchEngine string

const (
	Default MatchEngine = "default"
	OpenCV  MatchEngine = "opencv"
)

// 滑块匹配类型
type SlideMatchType string

const (
	Simple     SlideMatchType = "simple"
	Standard   SlideMatchType = "standard"
	Enhanced   SlideMatchType = "enhanced"
	Comparison SlideMatchType = "comparison"
)

// 目标图片路径、背景图片路径/Base64编码、匹配方式、匹配引擎，
// 比较模式的背景图为完整图片
func SlideMatch(targetStr, backgroundStr string, matchType SlideMatchType, matchEngine MatchEngine) (*ddddgocr.SlideBBox, error) {
	var targetData, backgroundData []byte
	_, err := os.Stat(targetStr)
	if err == nil {
		targetData, err = os.ReadFile(targetStr)
		if err != nil {
			return nil, fmt.Errorf("读取目标图像文件失败: %v", err)
		}
	} else {
		targetData, err = base64.StdEncoding.DecodeString(targetStr)
		if err != nil {
			return nil, fmt.Errorf("解析目标图像失败: %v", err)
		}
	}
	_, err = os.Stat(backgroundStr)
	if err == nil {
		backgroundData, err = os.ReadFile(backgroundStr)
		if err != nil {
			return nil, fmt.Errorf("读取背景图像文件失败: %v", err)
		}
	} else {
		backgroundData, err = base64.StdEncoding.DecodeString(backgroundStr)
		if err != nil {
			return nil, fmt.Errorf("解析背景图像失败: %v", err)
		}
	}

	return SlideMatchWithByte(targetData, backgroundData, matchType, matchEngine)
}

// 目标图片、背景图片、匹配方式、匹配引擎，
// 比较模式的背景图为完整图片
func SlideMatchWithByte(targetData, backgroundData []byte, matchType SlideMatchType, matchEngine MatchEngine) (*ddddgocr.SlideBBox, error) {
	if matchEngine == OpenCV {
		return slideMatchWithOpenCV(targetData, backgroundData, matchType)
	} else {
		switch matchType {
		case Simple:
			return ddddgocr.SimpleSlideMatch(targetData, backgroundData)
		case Standard:
			return ddddgocr.SlideMatch(targetData, backgroundData)
		case Enhanced:
			return ddddgocr.EnhancedSlideMatch(targetData, backgroundData)
		case Comparison:
			return ddddgocr.SlideComparison(targetData, backgroundData)
		default:
			return nil, fmt.Errorf("匹配类型错误")
		}
	}
}
