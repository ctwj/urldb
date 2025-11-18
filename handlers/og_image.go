package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ctwj/urldb/utils"
	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/entity"
	"github.com/gin-gonic/gin"
	"github.com/fogleman/gg"
	"image/color"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

// OGImageHandler 处理OG图片生成请求
type OGImageHandler struct{}

// NewOGImageHandler 创建新的OG图片处理器
func NewOGImageHandler() *OGImageHandler {
	return &OGImageHandler{}
}

// Resource 简化的资源结构体
type Resource struct {
	Title       string
	Description string
	Cover       string
	Key         string
}

// getResourceByKey 通过key获取资源信息
func (h *OGImageHandler) getResourceByKey(key string) (*Resource, error) {
	// 这里简化处理，实际应该从数据库查询
	// 为了演示，我们先返回一个模拟的资源
	// 在实际应用中，您需要连接数据库并查询

	// 模拟数据库查询 - 实际应用中请替换为真实的数据库查询
	dbInstance := db.DB
	if dbInstance == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var resource entity.Resource
	result := dbInstance.Where("key = ?", key).First(&resource)
	if result.Error != nil {
		return nil, result.Error
	}

	return &Resource{
		Title:       resource.Title,
		Description: resource.Description,
		Cover:       resource.Cover,
		Key:         resource.Key,
	}, nil
}

// GenerateOGImage 生成OG图片
func (h *OGImageHandler) GenerateOGImage(c *gin.Context) {
	// 获取请求参数
	key := strings.TrimSpace(c.Query("key"))
	title := strings.TrimSpace(c.Query("title"))
	description := strings.TrimSpace(c.Query("description"))
	siteName := strings.TrimSpace(c.Query("site_name"))
	theme := strings.TrimSpace(c.Query("theme"))
	coverUrl := strings.TrimSpace(c.Query("cover"))

	width, _ := strconv.Atoi(c.Query("width"))
	height, _ := strconv.Atoi(c.Query("height"))

	// 如果提供了key，从数据库获取资源信息
	if key != "" {
		resource, err := h.getResourceByKey(key)
		if err == nil && resource != nil {
			if title == "" {
				title = resource.Title
			}
			if description == "" {
				description = resource.Description
			}
			if coverUrl == "" && resource.Cover != "" {
				coverUrl = resource.Cover
			}
		}
	}

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

	// 获取当前请求的域名
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	domain := scheme + "://" + host

	// 生成图片
	imageBuffer, err := createOGImage(title, description, siteName, theme, width, height, coverUrl, key, domain)
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
func createOGImage(title, description, siteName, theme string, width, height int, coverUrl, key, domain string) (*bytes.Buffer, error) {
	dc := gg.NewContext(width, height)

	// 设置圆角裁剪区域
	cornerRadius := 20.0
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)

	// 设置背景色
	bgColor := getBackgroundColor(theme)
	dc.SetColor(bgColor)
	dc.Fill()

	// 绘制渐变效果
	gradient := gg.NewLinearGradient(0, 0, float64(width), float64(height))
	gradient.AddColorStop(0, getGradientStartColor(theme))
	gradient.AddColorStop(1, getGradientEndColor(theme))
	dc.SetFillStyle(gradient)
	dc.Fill()

	// 定义布局区域
	imageAreaWidth := width / 3  // 左侧1/3用于图片
	textAreaWidth := width * 2 / 3  // 右侧2/3用于文案
	textAreaX := imageAreaWidth  // 文案区域起始X坐标

	// 统一的字体加载函数，确保中文显示正常
	loadChineseFont := func(fontSize float64) bool {
		// 优先使用项目字体
		if err := dc.LoadFontFace("font/SourceHanSansSC-Regular.otf", fontSize); err == nil {
			return true
		}

		// Windows系统常见字体，按优先级顺序尝试
		commonFonts := []string{
			"C:/Windows/Fonts/msyh.ttc",     // 微软雅黑
			"C:/Windows/Fonts/simhei.ttf",   // 黑体
			"C:/Windows/Fonts/simsun.ttc",   // 宋体
		}

		for _, fontPath := range commonFonts {
			if err := dc.LoadFontFace(fontPath, fontSize); err == nil {
				return true
			}
		}

		// 如果都失败了，尝试使用粗体版本
		boldFonts := []string{
			"C:/Windows/Fonts/msyhbd.ttc",   // 微软雅黑粗体
			"C:/Windows/Fonts/simhei.ttf",   // 黑体
		}

		for _, fontPath := range boldFonts {
			if err := dc.LoadFontFace(fontPath, fontSize); err == nil {
				return true
			}
		}

		return false
	}

	// 加载基础字体（24px）
	fontLoaded := loadChineseFont(24)
	dc.SetHexColor("#ffffff")

	// 绘制封面图片（如果存在）
	if coverUrl != "" {
		if err := drawCoverImageInLeftArea(dc, coverUrl, width, height, imageAreaWidth); err != nil {
			utils.Error("绘制封面图片失败: %v", err)
		}
	}

	// 设置站点标识
	dc.DrawStringAnchored(siteName, float64(textAreaX)+60, 50, 0, 0.5)

	// 绘制标题
	dc.SetHexColor("#ffffff")

	// 标题在右侧区域显示，考虑文案宽度限制
	maxTitleWidth := float64(textAreaWidth - 120)  // 右侧区域减去左右边距

	// 动态调整字体大小以适应文案区域，使用统一的字体加载逻辑
	fontSize := 48.0
	titleFontLoaded := false
	for fontSize > 24 {  // 最小字体24
		// 优先使用项目粗体字体
		if err := dc.LoadFontFace("font/SourceHanSansSC-Bold.otf", fontSize); err == nil {
			titleWidth, _ := dc.MeasureString(title)
			if titleWidth <= maxTitleWidth {
				titleFontLoaded = true
				break
			}
		} else {
			// 尝试系统粗体字体
			boldFonts := []string{
				"C:/Windows/Fonts/msyhbd.ttc",   // 微软雅黑粗体
				"C:/Windows/Fonts/simhei.ttf",   // 黑体
			}
			for _, fontPath := range boldFonts {
				if err := dc.LoadFontFace(fontPath, fontSize); err == nil {
					titleWidth, _ := dc.MeasureString(title)
					if titleWidth <= maxTitleWidth {
						titleFontLoaded = true
						break
					}
					break  // 找到可用字体就跳出内层循环
				}
			}
			if titleFontLoaded {
				break
			}
		}
		fontSize -= 4
	}

	// 如果粗体字体都失败了，使用常规字体
	if !titleFontLoaded {
		loadChineseFont(36) // 使用稍大的常规字体
	}

	// 标题左对齐显示在右侧区域
	titleX := float64(textAreaX) + 60
	titleY := float64(height)/2 - 80
	dc.DrawString(title, titleX, titleY)

	// 绘制描述
	if description != "" {
		dc.SetHexColor("#e5e7eb")
		// 使用统一的字体加载逻辑
		loadChineseFont(28)

		// 自动换行处理，适配右侧区域宽度
		wrappedDesc := wrapText(dc, description, float64(textAreaWidth-120))
		descY := titleY + 60  // 标题下方

		for i, line := range wrappedDesc {
			y := descY + float64(i)*30  // 行高30像素
			dc.DrawString(line, titleX, y)
		}
	}

	// 添加装饰性元素
	drawDecorativeElements(dc, width, height, theme)

	// 绘制底部URL访问地址
	if key != "" && domain != "" {
		resourceURL := domain + "/r/" + key
		dc.SetHexColor("#d1d5db") // 浅灰色

		// 使用统一的字体加载逻辑
		loadChineseFont(20)

		// URL位置：底部居中，距离底部边缘40像素，给更多空间
		urlY := float64(height) - 40

		dc.DrawStringAnchored(resourceURL, float64(width)/2, urlY, 0.5, 0.5)
	}

	// 添加调试信息（仅在开发环境）
	if title == "DEBUG" {
		dc.SetHexColor("#ff0000")
		dc.DrawString("Font loaded: "+strconv.FormatBool(fontLoaded), 50, float64(height)-80)
	}

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

// drawCoverImageInLeftArea 在左侧1/3区域绘制封面图片
func drawCoverImageInLeftArea(dc *gg.Context, coverUrl string, width, height int, imageAreaWidth int) error {
	// 下载封面图片
	resp, err := http.Get(coverUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取图片数据
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return err
	}

	// 获取图片尺寸和宽高比
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	aspectRatio := float64(imgWidth) / float64(imgHeight)

	// 计算图片区域的可显示尺寸，留出边距
	padding := 40
	maxImageWidth := imageAreaWidth - padding*2
	maxImageHeight := height - padding*2

	var scaledImg image.Image
	var drawWidth, drawHeight, drawX, drawY int

	// 判断是竖图还是横图，采用不同的缩放策略
	if aspectRatio < 1.0 {
		// 竖图：充满整个左侧区域（去掉边距）
		drawHeight = height - padding*2  // 留上下边距
		drawWidth = int(float64(drawHeight) * aspectRatio)

		// 如果宽度超出左侧区域，则以宽度为准充满整个区域宽度
		if drawWidth > imageAreaWidth - padding*2 {
			drawWidth = imageAreaWidth - padding*2
			drawHeight = int(float64(drawWidth) / aspectRatio)
		}

		// 缩放图片
		scaledImg = scaleImage(img, drawWidth, drawHeight)

		// 垂直居中，水平居左
		drawX = padding
		drawY = (height - drawHeight) / 2
	} else {
		// 横图：优先占满宽度
		drawWidth = maxImageWidth
		drawHeight = int(float64(drawWidth) / aspectRatio)

		// 如果高度超出限制，则以高度为准
		if drawHeight > maxImageHeight {
			drawHeight = maxImageHeight
			drawWidth = int(float64(drawHeight) * aspectRatio)
		}

		// 缩放图片
		scaledImg = scaleImage(img, drawWidth, drawHeight)

		// 水平居中，垂直居中
		drawX = (imageAreaWidth - drawWidth) / 2
		drawY = (height - drawHeight) / 2
	}

	// 绘制图片
	dc.DrawImage(scaledImg, drawX, drawY)

	// 添加半透明遮罩效果，让文字更清晰（仅在有图片时添加）
	maskColor := color.RGBA{0, 0, 0, 80} // 半透明黑色，透明度稍低
	dc.SetColor(maskColor)
	dc.DrawRectangle(float64(drawX), float64(drawY), float64(drawWidth), float64(drawHeight))
	dc.Fill()

	return nil
}

// scaleImage 图片缩放函数
func scaleImage(img image.Image, width, height int) image.Image {
	// 使用 gg 库的 Scale 变换来实现缩放
	srcWidth := img.Bounds().Dx()
	srcHeight := img.Bounds().Dy()

	// 创建目标尺寸的画布
	dc := gg.NewContext(width, height)

	// 计算缩放比例
	scaleX := float64(width) / float64(srcWidth)
	scaleY := float64(height) / float64(srcHeight)

	// 应用缩放变换并绘制图片
	dc.Scale(scaleX, scaleY)
	dc.DrawImage(img, 0, 0)

	return dc.Image()
}

// drawCoverImage 绘制封面图片（保留原函数作为备用）
func drawCoverImage(dc *gg.Context, coverUrl string, width, height int) error {
	// 下载封面图片
	resp, err := http.Get(coverUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取图片数据
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return err
	}

	// 计算封面图片的位置和大小，放置在左侧
	coverWidth := 200  // 封面图宽度
	coverHeight := 280 // 封面图高度
	coverX := 50
	coverY := (height - coverHeight) / 2

	// 绘制封面图片（按比例缩放）
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// 计算缩放比例，保持宽高比
	scaleX := float64(coverWidth) / float64(imgWidth)
	scaleY := float64(coverHeight) / float64(imgHeight)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// 计算缩放后的尺寸
	newWidth := int(float64(imgWidth) * scale)
	newHeight := int(float64(imgHeight) * scale)

	// 居中绘制
	offsetX := coverX + (coverWidth-newWidth)/2
	offsetY := coverY + (coverHeight-newHeight)/2

	dc.DrawImage(img, offsetX, offsetY)

	// 添加半透明遮罩效果，让文字更清晰
	maskColor := color.RGBA{0, 0, 0, 120} // 半透明黑色
	dc.SetColor(maskColor)
	dc.DrawRectangle(float64(coverX), float64(coverY), float64(coverWidth), float64(coverHeight))
	dc.Fill()

	return nil
}