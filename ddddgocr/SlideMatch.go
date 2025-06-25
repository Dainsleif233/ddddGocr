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
	"math"
)

// SlideBBox 滑块边界框结构
type SlideBBox struct {
	TargetY int `json:"target_y"`
	X1      int `json:"x1"`
	Y1      int `json:"y1"`
	X2      int `json:"x2"`
	Y2      int `json:"y2"`
}

// SlideComparisonResult 坑位匹配结果
type SlideComparisonResult struct {
	X uint32 `json:"x"`
	Y uint32 `json:"y"`
}

// 滑块匹配主函数
func SlideMatch(targetImageData, backgroundImageData []byte) (*SlideBBox, error) {
	// 解码图像
	targetImg, _, err := image.Decode(bytes.NewReader(targetImageData))
	if err != nil {
		return nil, fmt.Errorf("解码目标图像失败: %v", err)
	}

	backgroundImg, _, err := image.Decode(bytes.NewReader(backgroundImageData))
	if err != nil {
		return nil, fmt.Errorf("解码背景图像失败: %v", err)
	}

	// 检查图像尺寸
	if backgroundImg.Bounds().Dx() < targetImg.Bounds().Dx() {
		return nil, errors.New("背景图片的宽度必须大于等于目标图片的宽度")
	}

	if backgroundImg.Bounds().Dy() < targetImg.Bounds().Dy() {
		return nil, errors.New("背景图片的高度必须大于等于目标图片的高度")
	}

	// 转换为RGBA格式
	targetRGBA := toRGBA(targetImg)

	// 裁剪透明区域
	croppedTarget, startY, _ := cropTransparent(targetRGBA)

	// 转换为灰度图
	targetGray := rgbaToGrayScale(croppedTarget)
	backgroundGray := toGrayScale(backgroundImg)

	// 边缘检测
	targetEdges := cannyEdgeDetection(targetGray, 100.0, 200.0)
	backgroundEdges := cannyEdgeDetection(backgroundGray, 100.0, 200.0)

	// 模板匹配
	matchResult := matchTemplate(backgroundEdges, targetEdges)
	if matchResult == nil {
		return nil, errors.New("模板匹配失败")
	}

	// 找到最佳匹配位置
	maxVal, maxX, maxY, _, _, _ := findExtremes(matchResult)

	if maxVal < 0.3 { // 设置一个阈值来判断匹配质量
		return nil, errors.New("匹配质量过低")
	}

	return &SlideBBox{
		TargetY: startY,
		X1:      maxX,
		Y1:      maxY,
		X2:      maxX + targetEdges.Bounds().Dx(),
		Y2:      maxY + targetEdges.Bounds().Dy(),
	}, nil
}

// 简单滑块匹配（无透明区域裁剪）
func SimpleSlideMatch(targetImageData, backgroundImageData []byte) (*SlideBBox, error) {
	// 解码图像
	targetImg, _, err := image.Decode(bytes.NewReader(targetImageData))
	if err != nil {
		return nil, fmt.Errorf("解码目标图像失败: %v", err)
	}

	backgroundImg, _, err := image.Decode(bytes.NewReader(backgroundImageData))
	if err != nil {
		return nil, fmt.Errorf("解码背景图像失败: %v", err)
	}

	// 检查图像尺寸
	if backgroundImg.Bounds().Dx() < targetImg.Bounds().Dx() {
		return nil, errors.New("背景图片的宽度必须大于等于目标图标的宽度")
	}

	if backgroundImg.Bounds().Dy() < targetImg.Bounds().Dy() {
		return nil, errors.New("背景图片的高度必须大于等于目标图标的高度")
	}

	// 转换为灰度图
	targetGray := toGrayScale(targetImg)
	backgroundGray := toGrayScale(backgroundImg)

	// 边缘检测
	targetEdges := cannyEdgeDetection(targetGray, 100.0, 200.0)
	backgroundEdges := cannyEdgeDetection(backgroundGray, 100.0, 200.0)

	// 模板匹配
	matchResult := matchTemplate(backgroundEdges, targetEdges)
	if matchResult == nil {
		return nil, errors.New("模板匹配失败")
	}

	// 找到最佳匹配位置
	maxVal, maxX, maxY, _, _, _ := findExtremes(matchResult)

	if maxVal < 0.3 { // 设置一个阈值来判断匹配质量
		return nil, errors.New("匹配质量过低")
	}

	return &SlideBBox{
		TargetY: 0,
		X1:      maxX,
		Y1:      maxY,
		X2:      maxX + targetEdges.Bounds().Dx(),
		Y2:      maxY + targetEdges.Bounds().Dy(),
	}, nil
}

func EnhancedSlideMatch(targetImageData, backgroundImageData []byte) (*SlideBBox, error) {
	// 解码图像
	targetImg, _, err := image.Decode(bytes.NewReader(targetImageData))
	if err != nil {
		return nil, fmt.Errorf("解码目标图像失败: %v", err)
	}

	backgroundImg, _, err := image.Decode(bytes.NewReader(backgroundImageData))
	if err != nil {
		return nil, fmt.Errorf("解码背景图像失败: %v", err)
	}

	// 检查图像尺寸
	if backgroundImg.Bounds().Dx() < targetImg.Bounds().Dx() {
		return nil, errors.New("背景图片的宽度必须大于等于目标图片的宽度")
	}

	if backgroundImg.Bounds().Dy() < targetImg.Bounds().Dy() {
		return nil, errors.New("背景图片的高度必须大于等于目标图片的高度")
	}

	// 策略1: 直接灰度匹配（不进行边缘检测）
	targetGray := toGrayScale(targetImg)
	backgroundGray := toGrayScale(backgroundImg)

	// 如果是RGBA图像，先处理透明区域
	var croppedTarget *image.Gray
	var startY int
	if _, ok := targetImg.(*image.RGBA); ok || hasTransparency(targetImg) {
		targetRGBA := toRGBA(targetImg)
		cropped, sy, _ := cropTransparent(targetRGBA)
		croppedTarget = rgbaToGrayScale(cropped)
		startY = sy
	} else {
		croppedTarget = targetGray
		startY = 0
	}

	results := make([]*SlideBBox, 0)

	// 策略1: 直接灰度模板匹配
	matchResult1 := matchTemplate(backgroundGray, croppedTarget)
	if matchResult1 != nil {
		maxVal, maxX, maxY, _, _, _ := findExtremes(matchResult1)
		// fmt.Printf("策略1 - 灰度匹配: 最大值=%.4f, 位置=(%d, %d)\n", maxVal, maxX, maxY)
		if maxVal > 0.6 {
			results = append(results, &SlideBBox{
				TargetY: startY,
				X1:      maxX,
				Y1:      maxY,
				X2:      maxX + croppedTarget.Bounds().Dx(),
				Y2:      maxY + croppedTarget.Bounds().Dy(),
			})
		}
	}

	// 策略2: 边缘检测匹配（低阈值）
	targetEdges1 := cannyEdgeDetection(croppedTarget, 30.0, 80.0)
	backgroundEdges1 := cannyEdgeDetection(backgroundGray, 30.0, 80.0)
	matchResult2 := matchTemplate(backgroundEdges1, targetEdges1)
	if matchResult2 != nil {
		maxVal, maxX, maxY, _, _, _ := findExtremes(matchResult2)
		// fmt.Printf("策略2 - 低阈值边缘: 最大值=%.4f, 位置=(%d, %d)\n", maxVal, maxX, maxY)
		if maxVal > 0.3 {
			results = append(results, &SlideBBox{
				TargetY: startY,
				X1:      maxX,
				Y1:      maxY,
				X2:      maxX + targetEdges1.Bounds().Dx(),
				Y2:      maxY + targetEdges1.Bounds().Dy(),
			})
		}
	}

	// 策略3: 边缘检测匹配（中等阈值）
	targetEdges2 := cannyEdgeDetection(croppedTarget, 50.0, 150.0)
	backgroundEdges2 := cannyEdgeDetection(backgroundGray, 50.0, 150.0)
	matchResult3 := matchTemplate(backgroundEdges2, targetEdges2)
	if matchResult3 != nil {
		maxVal, maxX, maxY, _, _, _ := findExtremes(matchResult3)
		// fmt.Printf("策略3 - 中阈值边缘: 最大值=%.4f, 位置=(%d, %d)\n", maxVal, maxX, maxY)
		if maxVal > 0.2 {
			results = append(results, &SlideBBox{
				TargetY: startY,
				X1:      maxX,
				Y1:      maxY,
				X2:      maxX + targetEdges2.Bounds().Dx(),
				Y2:      maxY + targetEdges2.Bounds().Dy(),
			})
		}
	}

	// 策略4: 差分匹配（寻找缺口）
	diffResult := findSlotByDifference(backgroundGray, croppedTarget)
	if diffResult != nil {
		// fmt.Printf("策略4 - 差分匹配: 位置=(%d, %d)\n", diffResult.X1, diffResult.Y1)
		results = append(results, diffResult)
	}

	if len(results) == 0 {
		return nil, errors.New("所有匹配策略都失败了")
	}

	// 选择最可信的结果（优先选择X1 > 0的结果）
	var bestResult *SlideBBox
	for _, result := range results {
		if result.X1 > 0 {
			if bestResult == nil || result.X1 < bestResult.X1 {
				bestResult = result
			}
		}
	}

	// 如果没有X1 > 0的结果，选择第一个
	if bestResult == nil {
		bestResult = results[0]
	}

	// fmt.Printf("最终选择结果: X1=%d, Y1=%d\n", bestResult.X1, bestResult.Y1)
	return bestResult, nil
}

// 检查图像是否有透明度
func hasTransparency(img image.Image) bool {
	switch img.(type) {
	case *image.RGBA, *image.NRGBA, *image.RGBA64, *image.NRGBA64:
		return true
	default:
		return false
	}
}

// 通过差分方法寻找滑块缺口
func findSlotByDifference(background, target *image.Gray) *SlideBBox {
	bgBounds := background.Bounds()
	tgtBounds := target.Bounds()

	bgWidth := bgBounds.Dx()
	bgHeight := bgBounds.Dy()
	tgtWidth := tgtBounds.Dx()
	tgtHeight := tgtBounds.Dy()

	if bgWidth < tgtWidth || bgHeight < tgtHeight {
		return nil
	}

	// 计算每一列的垂直边缘强度
	columnEdges := make([]float64, bgWidth)

	for x := 1; x < bgWidth-1; x++ {
		var edgeStrength float64
		for y := 1; y < bgHeight-1; y++ {
			// 计算垂直梯度
			top := float64(background.GrayAt(x, y-1).Y)
			bottom := float64(background.GrayAt(x, y+1).Y)
			gradient := math.Abs(bottom - top)
			edgeStrength += gradient
		}
		columnEdges[x] = edgeStrength / float64(bgHeight-2)
	}

	// 寻找垂直边缘强度的峰值（可能是滑块缺口的左边缘）
	maxEdge := 0.0
	maxX := 0

	// 寻找合适区域内的最大边缘强度
	searchStart := tgtWidth / 2
	searchEnd := bgWidth - tgtWidth - tgtWidth/2

	for x := searchStart; x < searchEnd; x++ {
		if columnEdges[x] > maxEdge {
			maxEdge = columnEdges[x]
			maxX = x
		}
	}

	// fmt.Printf("差分匹配 - 最大边缘强度: %.2f, 位置: %d\n", maxEdge, maxX)

	if maxEdge > 10.0 { // 设置一个边缘强度阈值
		// 寻找最佳的Y位置
		bestY := findBestYPosition(background, target, maxX)

		return &SlideBBox{
			TargetY: 0,
			X1:      maxX,
			Y1:      bestY,
			X2:      maxX + tgtWidth,
			Y2:      bestY + tgtHeight,
		}
	}

	return nil
}

// 在给定X位置寻找最佳Y位置
func findBestYPosition(background, target *image.Gray, x int) int {
	bgBounds := background.Bounds()
	tgtBounds := target.Bounds()

	bgHeight := bgBounds.Dy()
	tgtHeight := tgtBounds.Dy()
	tgtWidth := tgtBounds.Dx()

	bestY := 0
	bestScore := -1.0

	for y := 0; y <= bgHeight-tgtHeight; y++ {
		// 计算在这个位置的匹配分数
		var score float64
		var count int

		for ty := range tgtHeight {
			for tx := range tgtWidth {
				if x+tx < bgBounds.Dx() && y+ty < bgBounds.Dy() {
					bgVal := float64(background.GrayAt(x+tx, y+ty).Y)
					tgtVal := float64(target.GrayAt(tx, ty).Y)

					// 使用归一化相关性
					score += bgVal * tgtVal
					count++
				}
			}
		}

		if count > 0 {
			score /= float64(count)
			if score > bestScore {
				bestScore = score
				bestY = y
			}
		}
	}

	return bestY
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
