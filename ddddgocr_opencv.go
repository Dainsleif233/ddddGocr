//go:build opencv
// +build opencv

package ddddGocr

import (
	"fmt"

	"github.com/Dainsleif233/ddddGocr/ddddgocr"
	"github.com/Dainsleif233/ddddGocr/ddddgocr/withopencv"
)

func slideMatchWithOpenCV(targetData, backgroundData []byte, matchType SlideMatchType) (*ddddgocr.SlideBBox, error) {
	switch matchType {
	case Simple:
		return withopencv.SimpleSlideMatch(targetData, backgroundData)
	case Standard:
		return withopencv.SlideMatch(targetData, backgroundData)
	case Enhanced:
		return withopencv.EnhancedSlideMatch(targetData, backgroundData)
	case Comparison:
		return withopencv.SlideComparison(targetData, backgroundData)
	default:
		return nil, fmt.Errorf("匹配类型错误")
	}
}
