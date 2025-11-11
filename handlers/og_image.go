package handlers

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/ctwj/urldb/utils"
	"github.com/gin-gonic/gin"
	"github.com/fogleman/gg"
	"image/color"
)

// OGImageHandler 处理OG图片生成请求
type OGImageHandler struct{}

// NewOGImageHandler 创建新的OG图片处理器
func NewOGImageHandler() *OGImageHandler {
	return &OGImageHandler{}
}

// GenerateOGImage 生成OG图片
func (h *OGImageHandler) GenerateOGImage(c *gin.Context) {
	// 获取请求参数
	title := strings.TrimSpace(c.Query("title"))
	description := strings.TrimSpace(c.Query("description"))
	siteName := strings.TrimSpace(c.Query("site_name"))
	theme := strings.TrimSpace(c.Query("theme"))

	width, _ := strconv.Atoi(c.Query("width"))
	height, _ := strconv.Atoi(c.Query("height"))

	// 设置默认值
	if title == "" {
		title = "老九网盘资源数据库"
	}
	if siteName == "" {
		siteName = "老九网盘"
	}
	if width <= 0 || width > 2000 {
		width = 1200
	}
	if height <= 0 || height > 2000 {
		height = 630
	}

	// 生成图片
	imageBuffer, err := createOGImage(title, description, siteName, theme, width, height)
	if err != nil {
		utils.Error("生成OG图片失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate image: " + err.Error(),
		})
		return
	}

	// 返回图片
	c.Data(http.StatusOK, "image/png", imageBuffer.Bytes())
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "public, max-age=3600")
}

// createOGImage 创建OG图片
func createOGImage(title, description, siteName, theme string, width, height int) (*bytes.Buffer, error) {
	dc := gg.NewContext(width, height)

	// 设置背景色
	bgColor := getBackgroundColor(theme)
	dc.SetColor(bgColor)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	// 绘制渐变效果
	gradient := gg.NewLinearGradient(0, 0, float64(width), float64(height))
	gradient.AddColorStop(0, getGradientStartColor(theme))
	gradient.AddColorStop(1, getGradientEndColor(theme))
	dc.SetFillStyle(gradient)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	// 设置站点标识
	dc.SetHexColor("#ffffff")
	// 尝试加载字体，如果失败则使用默认字体
	if err := dc.LoadFontFace("assets/fonts/SourceHanSansCN-Regular.ttc", 24); err != nil {
		// 使用默认字体设置
	}

	dc.DrawStringAnchored(siteName, 60, 50, 0, 0.5)

	// 绘制标题
	dc.SetHexColor("#ffffff")
	if err := dc.LoadFontFace("assets/fonts/SourceHanSansCN-Bold.ttc", 48); err != nil {
		// 使用默认字体设置
	}

	// 文字居中处理
	titleWidth, _ := dc.MeasureString(title)
	if titleWidth > float64(width-120) {
		// 如果标题过长，尝试加载较小字体
		if err := dc.LoadFontFace("assets/fonts/SourceHanSansCN-Bold.ttc", 42); err != nil {
			// 使用默认字体设置
		}
		titleWidth, _ = dc.MeasureString(title)
		if titleWidth > float64(width-120) {
			if err := dc.LoadFontFace("assets/fonts/SourceHanSansCN-Bold.ttc", 36); err != nil {
				// 使用默认字体设置
			}
		}
	}

	dc.DrawStringAnchored(title, float64(width)/2, float64(height)/2-30, 0.5, 0.5)

	// 绘制描述
	if description != "" {
		dc.SetHexColor("#e5e7eb")
		// 尝试加载较小字体
		if err := dc.LoadFontFace("assets/fonts/SourceHanSansCN-Regular.ttc", 28); err != nil {
			// 使用默认字体设置
		}

		// 自动换行处理
		wrappedDesc := wrapText(dc, description, float64(width-120))
		startY := float64(height)/2 + 40

		for i, line := range wrappedDesc {
			y := startY + float64(i)*35
			dc.DrawStringAnchored(line, float64(width)/2, y, 0.5, 0.5)
		}
	}

	// 添加装饰性元素
	drawDecorativeElements(dc, width, height, theme)

	// 生成图片
	buf := &bytes.Buffer{}
	err := dc.EncodePNG(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// getBackgroundColor 获取背景色
func getBackgroundColor(theme string) color.RGBA {
	switch theme {
	case "dark":
		return color.RGBA{31, 41, 55, 255} // slate-800
	case "blue":
		return color.RGBA{29, 78, 216, 255} // blue-700
	case "green":
		return color.RGBA{6, 95, 70, 255} // emerald-800
	case "purple":
		return color.RGBA{109, 40, 217, 255} // violet-700
	default:
		return color.RGBA{55, 65, 81, 255} // gray-800
	}
}

// getGradientStartColor 获取渐变起始色
func getGradientStartColor(theme string) color.Color {
	switch theme {
	case "dark":
		return color.RGBA{15, 23, 42, 255} // slate-900
	case "blue":
		return color.RGBA{30, 58, 138, 255} // blue-900
	case "green":
		return color.RGBA{6, 78, 59, 255} // emerald-900
	case "purple":
		return color.RGBA{91, 33, 182, 255} // violet-800
	default:
		return color.RGBA{31, 41, 55, 255} // gray-800
	}
}

// getGradientEndColor 获取渐变结束色
func getGradientEndColor(theme string) color.Color {
	switch theme {
	case "dark":
		return color.RGBA{55, 65, 81, 255} // slate-700
	case "blue":
		return color.RGBA{59, 130, 246, 255} // blue-500
	case "green":
		return color.RGBA{16, 185, 129, 255} // emerald-500
	case "purple":
		return color.RGBA{139, 92, 246, 255} // violet-500
	default:
		return color.RGBA{75, 85, 99, 255} // gray-600
	}
}

// wrapText 文本自动换行处理
func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	var lines []string
	words := []rune(text)

	currentLine := ""
	for _, word := range words {
		testLine := currentLine + string(word)
		width, _ := dc.MeasureString(testLine)

		if width > maxWidth && len(currentLine) > 0 {
			lines = append(lines, currentLine)
			currentLine = string(word)
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	// 最多显示3行
	if len(lines) > 3 {
		lines = lines[:3]
		// 在最后一行添加省略号
		if len(lines[2]) > 3 {
			lines[2] = lines[2][:len(lines[2])-3] + "..."
		}
	}

	return lines
}

// drawDecorativeElements 绘制装饰性元素
func drawDecorativeElements(dc *gg.Context, width, height int, theme string) {
	// 绘制装饰性圆点
	dc.SetHexColor("#ffffff")
	dc.SetLineWidth(2)

	for i := 0; i < 5; i++ {
		x := float64(100 + i*150)
		y := float64(100 + (i%2)*200)
		dc.DrawCircle(x, y, 8)
		dc.Stroke()
	}

	// 绘制底部装饰线
	dc.DrawLine(60, float64(height-80), float64(width-60), float64(height-80))
	dc.Stroke()
}