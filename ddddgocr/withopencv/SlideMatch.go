package withopencv

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/Dainsleif233/ddddGocr/ddddgocr"
	"gocv.io/x/gocv"
)

// SlideMatch 滑块匹配主函数
func SlideMatch(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
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

	// 检查图像尺寸
	if backgroundMat.Cols() < targetMat.Cols() {
		return nil, errors.New("背景图片的宽度必须大于等于目标图片的宽度")
	}

	if backgroundMat.Rows() < targetMat.Rows() {
		return nil, errors.New("背景图片的高度必须大于等于目标图片的高度")
	}

	// 处理透明区域（如果目标图像有透明通道）
	var processedTarget gocv.Mat
	var startY int

	if targetMat.Channels() == 4 {
		// 有透明通道，需要处理
		processedTarget, startY = cropTransparentOpenCV(targetMat)
		defer processedTarget.Close()
	} else {
		processedTarget = targetMat.Clone()
		defer processedTarget.Close()
		startY = 0
	}

	// 转换为灰度图
	targetGray := gocv.NewMat()
	defer targetGray.Close()
	gocv.CvtColor(processedTarget, &targetGray, gocv.ColorBGRToGray)

	backgroundGray := gocv.NewMat()
	defer backgroundGray.Close()
	gocv.CvtColor(backgroundMat, &backgroundGray, gocv.ColorBGRToGray)

	// Canny边缘检测
	targetEdges := gocv.NewMat()
	defer targetEdges.Close()
	gocv.Canny(targetGray, &targetEdges, 100, 200)

	backgroundEdges := gocv.NewMat()
	defer backgroundEdges.Close()
	gocv.Canny(backgroundGray, &backgroundEdges, 100, 200)

	// 模板匹配
	matchResult := gocv.NewMat()
	defer matchResult.Close()
	gocv.MatchTemplate(backgroundEdges, targetEdges, &matchResult, gocv.TmCcoeffNormed, gocv.NewMat())

	// 找到最佳匹配位置
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(matchResult)

	if maxVal < 0.3 {
		return nil, errors.New("匹配质量过低")
	}

	return &ddddgocr.SlideBBox{
		TargetY: startY,
		X1:      maxLoc.X,
		Y1:      maxLoc.Y,
		X2:      maxLoc.X + targetEdges.Cols(),
		Y2:      maxLoc.Y + targetEdges.Rows(),
	}, nil
}

// SlideMatchWithPath 从文件路径读取图像进行滑块匹配
func SlideMatchWithPath(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideBBox, error) {
	targetData, err := os.ReadFile(targetImagePath)
	if err != nil {
		return nil, fmt.Errorf("读取目标图像文件失败: %v", err)
	}

	backgroundData, err := os.ReadFile(backgroundImagePath)
	if err != nil {
		return nil, fmt.Errorf("读取背景图像文件失败: %v", err)
	}

	return SlideMatch(targetData, backgroundData)
}

// SimpleSlideMatch 简单滑块匹配（无透明区域裁剪）
func SimpleSlideMatch(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
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

	// 检查图像尺寸
	if backgroundMat.Cols() < targetMat.Cols() {
		return nil, errors.New("背景图片的宽度必须大于等于目标图片的宽度")
	}

	if backgroundMat.Rows() < targetMat.Rows() {
		return nil, errors.New("背景图片的高度必须大于等于目标图片的高度")
	}

	// 转换为灰度图
	targetGray := gocv.NewMat()
	defer targetGray.Close()
	gocv.CvtColor(targetMat, &targetGray, gocv.ColorBGRToGray)

	backgroundGray := gocv.NewMat()
	defer backgroundGray.Close()
	gocv.CvtColor(backgroundMat, &backgroundGray, gocv.ColorBGRToGray)

	// Canny边缘检测
	targetEdges := gocv.NewMat()
	defer targetEdges.Close()
	gocv.Canny(targetGray, &targetEdges, 100, 200)

	backgroundEdges := gocv.NewMat()
	defer backgroundEdges.Close()
	gocv.Canny(backgroundGray, &backgroundEdges, 100, 200)

	// 模板匹配
	matchResult := gocv.NewMat()
	defer matchResult.Close()
	gocv.MatchTemplate(backgroundEdges, targetEdges, &matchResult, gocv.TmCcoeffNormed, gocv.NewMat())

	// 找到最佳匹配位置
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(matchResult)

	if maxVal < 0.3 {
		return nil, errors.New("匹配质量过低")
	}

	return &ddddgocr.SlideBBox{
		TargetY: 0,
		X1:      maxLoc.X,
		Y1:      maxLoc.Y,
		X2:      maxLoc.X + targetEdges.Cols(),
		Y2:      maxLoc.Y + targetEdges.Rows(),
	}, nil
}

// SimpleSlideMatchWithPath 从文件路径读取图像进行简单滑块匹配
func SimpleSlideMatchWithPath(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideBBox, error) {
	targetData, err := os.ReadFile(targetImagePath)
	if err != nil {
		return nil, fmt.Errorf("读取目标图像文件失败: %v", err)
	}

	backgroundData, err := os.ReadFile(backgroundImagePath)
	if err != nil {
		return nil, fmt.Errorf("读取背景图像文件失败: %v", err)
	}

	return SimpleSlideMatch(targetData, backgroundData)
}

// EnhancedSlideMatch 增强版滑块匹配
func EnhancedSlideMatch(targetImageData, backgroundImageData []byte) (*ddddgocr.SlideBBox, error) {
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

	// 检查图像尺寸
	if backgroundMat.Cols() < targetMat.Cols() {
		return nil, errors.New("背景图片的宽度必须大于等于目标图片的宽度")
	}

	if backgroundMat.Rows() < targetMat.Rows() {
		return nil, errors.New("背景图片的高度必须大于等于目标图片的高度")
	}

	// 处理透明区域（如果有的话）
	var processedTarget gocv.Mat
	var startY int

	if targetMat.Channels() == 4 {
		processedTarget, startY = cropTransparentOpenCV(targetMat)
		defer processedTarget.Close()
	} else {
		processedTarget = targetMat.Clone()
		defer processedTarget.Close()
		startY = 0
	}

	// 转换为灰度图
	targetGray := gocv.NewMat()
	defer targetGray.Close()
	gocv.CvtColor(processedTarget, &targetGray, gocv.ColorBGRToGray)

	backgroundGray := gocv.NewMat()
	defer backgroundGray.Close()
	gocv.CvtColor(backgroundMat, &backgroundGray, gocv.ColorBGRToGray)

	results := make([]*ddddgocr.SlideBBox, 0)

	// 策略1: 直接灰度模板匹配
	matchResult1 := gocv.NewMat()
	defer matchResult1.Close()
	gocv.MatchTemplate(backgroundGray, targetGray, &matchResult1, gocv.TmCcoeffNormed, gocv.NewMat())
	_, maxVal1, _, maxLoc1 := gocv.MinMaxLoc(matchResult1)

	if maxVal1 > 0.6 {
		results = append(results, &ddddgocr.SlideBBox{
			TargetY: startY,
			X1:      maxLoc1.X,
			Y1:      maxLoc1.Y,
			X2:      maxLoc1.X + targetGray.Cols(),
			Y2:      maxLoc1.Y + targetGray.Rows(),
		})
	}

	// 策略2: 低阈值边缘检测匹配
	targetEdges1 := gocv.NewMat()
	defer targetEdges1.Close()
	gocv.Canny(targetGray, &targetEdges1, 30, 80)

	backgroundEdges1 := gocv.NewMat()
	defer backgroundEdges1.Close()
	gocv.Canny(backgroundGray, &backgroundEdges1, 30, 80)

	matchResult2 := gocv.NewMat()
	defer matchResult2.Close()
	gocv.MatchTemplate(backgroundEdges1, targetEdges1, &matchResult2, gocv.TmCcoeffNormed, gocv.NewMat())
	_, maxVal2, _, maxLoc2 := gocv.MinMaxLoc(matchResult2)

	if maxVal2 > 0.3 {
		results = append(results, &ddddgocr.SlideBBox{
			TargetY: startY,
			X1:      maxLoc2.X,
			Y1:      maxLoc2.Y,
			X2:      maxLoc2.X + targetEdges1.Cols(),
			Y2:      maxLoc2.Y + targetEdges1.Rows(),
		})
	}

	// 策略3: 中等阈值边缘检测匹配
	targetEdges2 := gocv.NewMat()
	defer targetEdges2.Close()
	gocv.Canny(targetGray, &targetEdges2, 50, 150)

	backgroundEdges2 := gocv.NewMat()
	defer backgroundEdges2.Close()
	gocv.Canny(backgroundGray, &backgroundEdges2, 50, 150)

	matchResult3 := gocv.NewMat()
	defer matchResult3.Close()
	gocv.MatchTemplate(backgroundEdges2, targetEdges2, &matchResult3, gocv.TmCcoeffNormed, gocv.NewMat())
	_, maxVal3, _, maxLoc3 := gocv.MinMaxLoc(matchResult3)

	if maxVal3 > 0.2 {
		results = append(results, &ddddgocr.SlideBBox{
			TargetY: startY,
			X1:      maxLoc3.X,
			Y1:      maxLoc3.Y,
			X2:      maxLoc3.X + targetEdges2.Cols(),
			Y2:      maxLoc3.Y + targetEdges2.Rows(),
		})
	}

	// 策略4: 使用SIFT特征匹配（可选）
	if len(results) == 0 {
		siftResult := siftFeatureMatch(backgroundGray, targetGray)
		if siftResult != nil {
			results = append(results, siftResult)
		}
	}

	if len(results) == 0 {
		return nil, errors.New("所有匹配策略都失败了")
	}

	// 选择最可信的结果
	var bestResult *ddddgocr.SlideBBox
	for _, result := range results {
		if result.X1 > 0 {
			if bestResult == nil || result.X1 < bestResult.X1 {
				bestResult = result
			}
		}
	}

	if bestResult == nil {
		bestResult = results[0]
	}

	return bestResult, nil
}

// EnhancedSlideMatchWithPath 增强版滑块匹配的文件路径版本
func EnhancedSlideMatchWithPath(targetImagePath, backgroundImagePath string) (*ddddgocr.SlideBBox, error) {
	targetData, err := os.ReadFile(targetImagePath)
	if err != nil {
		return nil, fmt.Errorf("读取目标图像文件失败: %v", err)
	}

	backgroundData, err := os.ReadFile(backgroundImagePath)
	if err != nil {
		return nil, fmt.Errorf("读取背景图像文件失败: %v", err)
	}

	return EnhancedSlideMatch(targetData, backgroundData)
}

// cropTransparentOpenCV 使用OpenCV裁剪透明区域
func cropTransparentOpenCV(img gocv.Mat) (gocv.Mat, int) {
	// 分离通道
	channels := gocv.Split(img)
	defer func() {
		for _, ch := range channels {
			ch.Close()
		}
	}()

	if len(channels) < 4 {
		// 没有透明通道，返回原图
		return img.Clone(), 0
	}

	alphaMat := channels[3] // 透明通道

	// 找到非透明区域的边界
	nonZeroPoints := gocv.NewMat()
	defer nonZeroPoints.Close()
	gocv.FindNonZero(alphaMat, &nonZeroPoints)

	if nonZeroPoints.Rows() == 0 {
		// 没有非透明像素，返回原图
		return img.Clone(), 0
	}

	// 计算边界框
	pointVec := gocv.NewPointVectorFromMat(nonZeroPoints)
	defer pointVec.Close()
	boundingRect := gocv.BoundingRect(pointVec)

	// 裁剪图像
	croppedImg := img.Region(boundingRect)

	return croppedImg, boundingRect.Min.Y
}

// siftFeatureMatch 使用SIFT特征进行匹配
func siftFeatureMatch(background, target gocv.Mat) *ddddgocr.SlideBBox {
	// 创建SIFT检测器
	sift := gocv.NewSIFT()
	defer sift.Close()

	// 检测关键点和描述符
	kp1, desc1 := sift.DetectAndCompute(target, gocv.NewMat())
	defer desc1.Close()

	kp2, desc2 := sift.DetectAndCompute(background, gocv.NewMat())
	defer desc2.Close()

	if desc1.Rows() == 0 || desc2.Rows() == 0 {
		return nil
	}

	// 创建匹配器
	matcher := gocv.NewBFMatcher()
	defer matcher.Close()

	// 进行匹配
	matches := matcher.Match(desc1, desc2)

	if len(matches) < 4 {
		return nil
	}

	// 提取匹配点
	srcPoints := make([]image.Point, 0)
	dstPoints := make([]image.Point, 0)

	for _, match := range matches {
		if match.Distance < 0.7 { // 过滤好的匹配
			srcPoint := image.Point{X: int(kp1[match.QueryIdx].X), Y: int(kp1[match.QueryIdx].Y)}
			dstPoint := image.Point{X: int(kp2[match.TrainIdx].X), Y: int(kp2[match.TrainIdx].Y)}

			srcPoints = append(srcPoints, srcPoint)
			dstPoints = append(dstPoints, dstPoint)
		}
	}

	if len(srcPoints) < 4 {
		return nil
	}

	// 使用RANSAC找到最佳匹配位置
	// 这里简化处理，取平均偏移量
	var sumX, sumY int
	for i := range srcPoints {
		sumX += dstPoints[i].X - srcPoints[i].X
		sumY += dstPoints[i].Y - srcPoints[i].Y
	}

	avgX := sumX / len(srcPoints)
	avgY := sumY / len(srcPoints)

	return &ddddgocr.SlideBBox{
		TargetY: 0,
		X1:      avgX,
		Y1:      avgY,
		X2:      avgX + target.Cols(),
		Y2:      avgY + target.Rows(),
	}
}
