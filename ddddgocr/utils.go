package ddddgocr

import (
	"image"
	"image/color"
	"math"
)

// Canny边缘检测算法
func cannyEdgeDetection(img *image.Gray, lowThreshold, highThreshold float64) *image.Gray {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 高斯模糊
	blurred := gaussianBlur(img)

	// 计算梯度
	gradX := make([][]float64, height)
	gradY := make([][]float64, height)
	magnitude := make([][]float64, height)
	direction := make([][]float64, height)

	for i := range gradX {
		gradX[i] = make([]float64, width)
		gradY[i] = make([]float64, width)
		magnitude[i] = make([]float64, width)
		direction[i] = make([]float64, width)
	}

	// Sobel算子计算梯度
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// Sobel X
			gx := -1*getGrayValue(blurred, x-1, y-1) + 1*getGrayValue(blurred, x+1, y-1) +
				-2*getGrayValue(blurred, x-1, y) + 2*getGrayValue(blurred, x+1, y) +
				-1*getGrayValue(blurred, x-1, y+1) + 1*getGrayValue(blurred, x+1, y+1)

			// Sobel Y
			gy := -1*getGrayValue(blurred, x-1, y-1) + -2*getGrayValue(blurred, x, y-1) + -1*getGrayValue(blurred, x+1, y-1) +
				1*getGrayValue(blurred, x-1, y+1) + 2*getGrayValue(blurred, x, y+1) + 1*getGrayValue(blurred, x+1, y+1)

			gradX[y][x] = gx
			gradY[y][x] = gy
			magnitude[y][x] = math.Sqrt(gx*gx + gy*gy)
			direction[y][x] = math.Atan2(gy, gx)
		}
	}

	// 非最大抑制
	suppressed := nonMaximumSuppression(magnitude, direction, width, height)

	// 双阈值检测
	result := doubleThreshold(suppressed, lowThreshold, highThreshold, width, height)

	return result
}

// 双阈值检测
func doubleThreshold(suppressed [][]float64, lowThreshold, highThreshold float64, width, height int) *image.Gray {
	result := image.NewGray(image.Rect(0, 0, width, height))

	for y := range height {
		for x := range width {
			if suppressed[y][x] >= highThreshold {
				result.Set(x, y, color.Gray{Y: 255})
			} else if suppressed[y][x] >= lowThreshold {
				result.Set(x, y, color.Gray{Y: 128})
			} else {
				result.Set(x, y, color.Gray{Y: 0})
			}
		}
	}

	// 边缘连接（简化版）
	edgeTracking(result, width, height)

	return result
}

// 边缘跟踪
func edgeTracking(img *image.Gray, width, height int) {
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			if img.GrayAt(x, y).Y == 128 {
				// 检查8邻域是否有强边缘
				hasStrongEdge := false
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if img.GrayAt(x+dx, y+dy).Y == 255 {
							hasStrongEdge = true
							break
						}
					}
					if hasStrongEdge {
						break
					}
				}

				if hasStrongEdge {
					img.Set(x, y, color.Gray{Y: 255})
				} else {
					img.Set(x, y, color.Gray{Y: 0})
				}
			}
		}
	}
}

// 非最大抑制
func nonMaximumSuppression(magnitude [][]float64, direction [][]float64, width, height int) [][]float64 {
	result := make([][]float64, height)
	for i := range result {
		result[i] = make([]float64, width)
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			angle := direction[y][x] * 180 / math.Pi
			if angle < 0 {
				angle += 180
			}

			var q, r float64

			// 根据梯度方向确定邻近像素
			if (0 <= angle && angle < 22.5) || (157.5 <= angle && angle <= 180) {
				q = magnitude[y][x+1]
				r = magnitude[y][x-1]
			} else if 22.5 <= angle && angle < 67.5 {
				q = magnitude[y+1][x-1]
				r = magnitude[y-1][x+1]
			} else if 67.5 <= angle && angle < 112.5 {
				q = magnitude[y+1][x]
				r = magnitude[y-1][x]
			} else if 112.5 <= angle && angle < 157.5 {
				q = magnitude[y-1][x-1]
				r = magnitude[y+1][x+1]
			}

			if magnitude[y][x] >= q && magnitude[y][x] >= r {
				result[y][x] = magnitude[y][x]
			} else {
				result[y][x] = 0
			}
		}
	}

	return result
}

// 获取灰度值
func getGrayValue(img *image.Gray, x, y int) float64 {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return 0
	}
	return float64(img.GrayAt(x, y).Y)
}

// 高斯模糊
func gaussianBlur(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 简化的高斯模糊，使用3x3核
	kernel := [][]float64{
		{1, 2, 1},
		{2, 4, 2},
		{1, 2, 1},
	}

	// 归一化核
	sum := 16.0
	for i := range kernel {
		for j := range kernel[i] {
			kernel[i][j] /= sum
		}
	}

	result := image.NewGray(bounds)

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var value float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					value += getGrayValue(img, x+kx, y+ky) * kernel[ky+1][kx+1]
				}
			}
			result.Set(x, y, color.Gray{Y: uint8(math.Max(0, math.Min(255, value)))})
		}
	}

	return result
}

// 查找极值
func findExtremes(matrix [][]float64) (maxVal float64, maxX, maxY int, minVal float64, minX, minY int) {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return
	}

	maxVal = matrix[0][0]
	minVal = matrix[0][0]
	maxX, maxY = 0, 0
	minX, minY = 0, 0

	for y := range matrix {
		for x := range matrix[y] {
			if matrix[y][x] > maxVal {
				maxVal = matrix[y][x]
				maxX, maxY = x, y
			}
			if matrix[y][x] < minVal {
				minVal = matrix[y][x]
				minX, minY = x, y
			}
		}
	}

	return
}

// 模板匹配 - 标准化交叉相关
func matchTemplate(background, template *image.Gray) [][]float64 {
	bgBounds := background.Bounds()
	tplBounds := template.Bounds()

	bgWidth := bgBounds.Dx()
	bgHeight := bgBounds.Dy()
	tplWidth := tplBounds.Dx()
	tplHeight := tplBounds.Dy()

	resultWidth := bgWidth - tplWidth + 1
	resultHeight := bgHeight - tplHeight + 1

	if resultWidth <= 0 || resultHeight <= 0 {
		return nil
	}

	result := make([][]float64, resultHeight)
	for i := range result {
		result[i] = make([]float64, resultWidth)
	}

	// 计算模板的均值
	var templateSum float64
	templatePixels := tplWidth * tplHeight
	for y := range tplHeight {
		for x := range tplWidth {
			templateSum += getGrayValue(template, x, y)
		}
	}
	templateMean := templateSum / float64(templatePixels)

	// 计算模板的标准差
	var templateSumSq float64
	for y := range tplHeight {
		for x := range tplWidth {
			diff := getGrayValue(template, x, y) - templateMean
			templateSumSq += diff * diff
		}
	}
	templateStd := math.Sqrt(templateSumSq)

	// 对每个可能的位置进行匹配
	for y := range resultHeight {
		for x := range resultWidth {
			// 计算当前窗口的均值
			var windowSum float64
			for wy := range tplHeight {
				for wx := range tplWidth {
					windowSum += getGrayValue(background, x+wx, y+wy)
				}
			}
			windowMean := windowSum / float64(templatePixels)

			// 计算当前窗口的标准差和相关性
			var windowSumSq, correlation float64
			for wy := range tplHeight {
				for wx := range tplWidth {
					bgVal := getGrayValue(background, x+wx, y+wy)
					tplVal := getGrayValue(template, wx, wy)

					bgDiff := bgVal - windowMean
					tplDiff := tplVal - templateMean

					windowSumSq += bgDiff * bgDiff
					correlation += bgDiff * tplDiff
				}
			}

			windowStd := math.Sqrt(windowSumSq)

			// 标准化交叉相关
			if windowStd > 0 && templateStd > 0 {
				result[y][x] = correlation / (windowStd * templateStd)
			} else {
				result[y][x] = 0
			}
		}
	}

	return result
}

// 将图像转换为RGBA格式
func toRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	return rgba
}

// 裁剪透明区域
func cropTransparent(img *image.RGBA) (*image.RGBA, int, int) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	startX, startY := width, height
	endX, endY := 0, 0

	// 找到不透明像素的边界
	for x := range width {
		for y := range height {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0 {
				if x < startX {
					startX = x
				}
				if y < startY {
					startY = y
				}
				if x > endX {
					endX = x
				}
				if y > endY {
					endY = y
				}
			}
		}
	}

	// 如果没有不透明像素，返回原图
	if startX > endX || startY > endY {
		return img, startY, startX
	}

	// 裁剪图像
	cropWidth := endX - startX + 1
	cropHeight := endY - startY + 1
	cropped := image.NewRGBA(image.Rect(0, 0, cropWidth, cropHeight))

	for y := range cropHeight {
		for x := range cropWidth {
			cropped.Set(x, y, img.At(startX+x, startY+y))
		}
	}

	return cropped, startY, startX
}

// 将彩色图像转换为灰度图像
func toGrayScale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			gray.Set(x, y, grayColor)
		}
	}

	return gray
}

// 将RGBA图像转换为灰度图像
func rgbaToGrayScale(img *image.RGBA) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// 使用标准灰度转换公式
			grayValue := uint8((299*r + 587*g + 114*b) / 1000 / 256)
			gray.Set(x, y, color.Gray{Y: grayValue})
		}
	}

	return gray
}

// absDiff 计算两个uint8值的绝对差值
func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}
