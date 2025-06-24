package ddddgocr

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"  // 导入GIF格式支持
	_ "image/jpeg" // 导入JPEG格式支持
	_ "image/png"  // 导入PNG格式支持
	"os"
)

// SlideComparisonResult 坑位匹配结果
type SlideComparisonResult struct {
	X uint32 `json:"x"`
	Y uint32 `json:"y"`
}

// SlideComparison 坑位匹配
// 通过比较两张相同尺寸的图片，找出差异区域来定位坑位位置
func SlideComparison(targetImageData, backgroundImageData []byte) (*SlideComparisonResult, error) {
	// 解码图像
	targetImg, _, err := image.Decode(bytes.NewReader(targetImageData))
	if err != nil {
		return nil, fmt.Errorf("解码目标图像失败: %v", err)
	}

	backgroundImg, _, err := image.Decode(bytes.NewReader(backgroundImageData))
	if err != nil {
		return nil, fmt.Errorf("解码背景图像失败: %v", err)
	}

	// 检查图像尺寸是否相等
	if targetImg.Bounds().Dx() != backgroundImg.Bounds().Dx() ||
		targetImg.Bounds().Dy() != backgroundImg.Bounds().Dy() {
		return nil, errors.New("图片尺寸不相等")
	}

	width := targetImg.Bounds().Dx()
	height := targetImg.Bounds().Dy()

	// 创建差异图像
	diffImage := image.NewGray(image.Rect(0, 0, width, height))

	// 计算像素差异
	for y := range height {
		for x := range width {
			// 获取目标图像和背景图像的像素值
			targetR, targetG, targetB, _ := targetImg.At(x, y).RGBA()
			bgR, bgG, bgB, _ := backgroundImg.At(x, y).RGBA()

			// 转换为8位值
			tR, tG, tB := uint8(targetR>>8), uint8(targetG>>8), uint8(targetB>>8)
			bR, bG, bB := uint8(bgR>>8), uint8(bgG>>8), uint8(bgB>>8)

			// 计算RGB差异的平均值
			diffR := absDiff(tR, bR)
			diffG := absDiff(tG, bG)
			diffB := absDiff(tB, bB)
			avgDiff := (uint16(diffR) + uint16(diffG) + uint16(diffB)) / 3

			// 如果差异大于80，设为白色(255)，否则为黑色(0)
			if avgDiff > 80 {
				diffImage.Set(x, y, color.Gray{Y: 255})
			} else {
				diffImage.Set(x, y, color.Gray{Y: 0})
			}
		}
	}

	var startX, startY uint32 = 0, 0

	// 按列扫描寻找差异区域
	for x := range width {
		count := 0

		for y := range height {
			pixel := diffImage.GrayAt(x, y)

			// 如果像素不是黑色（即存在差异）
			if pixel.Y != 0 {
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

	return &SlideComparisonResult{
		X: startX,
		Y: startY,
	}, nil
}

// SlideComparisonWithPath 从文件路径读取图像进行坑位匹配
func SlideComparisonWithPath(targetImagePath, backgroundImagePath string) (*SlideComparisonResult, error) {
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
