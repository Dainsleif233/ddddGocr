//go:build !opencv
// +build !opencv

package ddddGocr

import (
	"fmt"

	"github.com/Dainsleif233/ddddGocr/ddddgocr"
)

func slideMatchWithOpenCV(_, _ []byte, _ SlideMatchType) (*ddddgocr.SlideBBox, error) {
	return nil, fmt.Errorf("OpenCV 支持未启用，请使用 -tags opencv 编译")
}
