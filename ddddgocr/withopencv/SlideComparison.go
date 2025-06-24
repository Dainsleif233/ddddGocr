package withopencv

import (
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/Dainsleif233/ddddGocr/ddddgocr"
	"gocv.io/x/gocv"
)

// SlideComparison 坑位匹配
// 通过比较两张相同尺寸的图片，找出差异区域来定位坑位位置
func SlideComparison(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideComparisonResult, error) {
	// 从字节数据解码为Mat
	targetMat, err := gocv.IMDecode(targetImageData, gocv.IMReadColor)
	if err != nil {
		return nil, fmt.Errorf("解码目标图像失败: %v", err)
	}
	defer targetMat.Close()

	backgroundMat, err := gocv.IMDecode(backgroundImageData, gocv.IMReadColor)
	if err != nil {
		return nil, fmt.Errorf("解码背景图像失败: %v", err)
	}
	defer backgroundMat.Close()

	// 检查图像尺寸是否相等
	if targetMat.Cols() != backgroundMat.Cols() || targetMat.Rows() != backgroundMat.Rows() {
		return nil, errors.New("图片尺寸不相等")
	}

	// 计算差异图像
	diffMat := gocv.NewMat()
	defer diffMat.Close()

	gocv.AbsDiff(targetMat, backgroundMat, &diffMat)

	// 转换为灰度图
	grayMat := gocv.NewMat()
	defer grayMat.Close()
	gocv.CvtColor(diffMat, &grayMat, gocv.ColorBGRToGray)

	// 二值化处理
	binaryMat := gocv.NewMat()
	defer binaryMat.Close()
	gocv.Threshold(grayMat, &binaryMat, 80, 255, gocv.ThresholdBinary)

	var startX, startY uint32 = 0, 0

	// 按列扫描寻找差异区域
	width := binaryMat.Cols()
	height := binaryMat.Rows()

	for x := 0; x < width; x++ {
		count := 0

		for y := 0; y < height; y++ {
			pixelValue := binaryMat.GetUCharAt(y, x)

			// 如果像素不是黑色（即存在差异）
			if pixelValue != 0 {
				count++
			}

			// 如果连续发现5个差异像素且还未设置startY
			if count >= 5 && startY == 0 {
				if y >= 5 {
					startY = uint32(y - 5)
				} else {
					startY = 0
				}
			}
		}

		// 如果该列有足够的差异像素，说明找到了坑位的起始位置
		if count >= 5 {
			startX = uint32(x + 2) // 稍微向右偏移2个像素
			break
		}
	}

	return &ddddgocr.SlideComparisonResult{
		X: startX,
		Y: startY,
	}, nil
}

// SlideComparisonWithPath 从文件路径读取图像进行坑位匹配
func SlideComparisonWithPath(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideComparisonResult, error) {
	targetData, err := os.ReadFile(targetImagePath)
	if err != nil {
		return nil, fmt.Errorf("读取目标图像文件失败: %v", err)
	}

	backgroundData, err := os.ReadFile(backgroundImagePath)
	if err != nil {
		return nil, fmt.Errorf("读取背景图像文件失败: %v", err)
	}

	return SlideComparison(targetData, backgroundData)
}
