package pan

import (
	"fmt"
	"regexp"
	"strconv"
)

// 预编译正则（批量转存循环中频繁调用，避免每次重复编译）
var (
	reSurlParam  = regexp.MustCompile(`surl=([a-zA-Z0-9_-]+)`)
	reSurlPath   = regexp.MustCompile(`/s/([a-zA-Z0-9_-]+)`)
	reShareID    = regexp.MustCompile(`"shareid":(\d+)[,}]`)
	// share_uk 历史是字符串 "share_uk":"123"，新版页面偶为数字 "share_uk":123，
	// 故引号设为可选。捕获组始终是数字本身。
	reUserID     = regexp.MustCompile(`"share_uk":"?(\d+)"?[,}]`)
	reFsID       = regexp.MustCompile(`"fs_id":(\d+)[,}]`)
	reFilename   = regexp.MustCompile(`"server_filename":"([^"]+)"[,}]`)
	reIsDir      = regexp.MustCompile(`"isdir":(\d+)[,}]`)
)

// baiduShareFile 百度分享页解析出的文件信息（包内可见）
type baiduShareFile struct {
	FsID       int64
	ShareID    int64
	UK         int64
	ServerName string
	IsDir      bool
}

// ErrorCodeMap 百度网盘 errno → 可读中文消息
var ErrorCodeMap = map[int]string{
	-1:     "链接错误，链接失效或缺少提取码",
	-3:     "分享失败，文件不存在或无法分享",
	-4:     "转存失败，无效登录。请退出账号在其他地方的登录",
	-6:     "转存失败，请用浏览器无痕模式获取 Cookie 后再试",
	-7:     "转存失败，转存文件夹名有非法字符，不能包含 < > | * ? \\ :",
	-8:     "转存失败，目录中已有同名文件或文件夹存在",
	-9:     "链接错误，提取码错误",
	-10:    "转存失败，容量不足",
	-12:    "链接错误，提取码错误",
	-62:    "转存失败，链接访问次数过多，请稍后再试",
	0:      "转存成功",
	2:      "转存失败，目标目录不存在",
	4:      "转存失败，目录中存在同名文件",
	12:     "转存失败，转存文件数超过限制",
	20:     "转存失败，容量不足",
	105:    "链接错误，所访问的页面不存在",
	115:    "分享链接已失效（文件禁止分享）",
	145:    "分享链接已失效",
	-65:    "触发频率限制",
	132:    "删除文件需要二次身份验证（手机短信），请在浏览器登录账号完成验证后再试",
	200025: "提取码输入错误，请检查提取码",
}

// ErrnoMessage 查 errno 对应消息，未知码返回兜底文案
func ErrnoMessage(errno int) string {
	if msg, ok := ErrorCodeMap[errno]; ok {
		return msg
	}
	return fmt.Sprintf("未知错误(errno=%d)", errno)
}

// ExtractSurl 从百度分享链接提取 surl（去掉 /s/ 前导的 "1"）
func ExtractSurl(shareURL string) string {
	if m := reSurlParam.FindStringSubmatch(shareURL); len(m) > 1 {
		return m[1]
	}
	m := reSurlPath.FindStringSubmatch(shareURL)
	if len(m) < 2 {
		return ""
	}
	surl := m[1]
	// 百度 /s/ 短码的展示形式以伪字符 "1" 开头（如 https://pan.baidu.com/s/1xxxx），
	// 真实 surl 不含该前导 1，需剥离。而 surl= 参数形式已是原始值，不剥离。
	if len(surl) > 0 && surl[0] == '1' {
		surl = surl[1:]
	}
	return surl
}

// ParseSharePageHTML 解析百度分享页 HTML/JSON 片段，提取 shareid/share_uk/fs_id/server_filename/isdir
func ParseSharePageHTML(response string) ([]baiduShareFile, error) {
	shareIDs := reShareID.FindAllStringSubmatch(response, -1)
	userIDs := reUserID.FindAllStringSubmatch(response, -1)
	fsIDs := reFsID.FindAllStringSubmatch(response, -1)
	filenames := reFilename.FindAllStringSubmatch(response, -1)
	isDirs := reIsDir.FindAllStringSubmatch(response, -1)

	if len(shareIDs) == 0 || len(userIDs) == 0 || len(fsIDs) == 0 {
		// 截取响应前 300 字符放进错误信息，方便定位是 Cookie 失效、被风控、还是字段格式变化
		snippet := response
		if len(snippet) > 300 {
			snippet = snippet[:300]
		}
		return nil, fmt.Errorf("解析分享链接响应失败, 可能是提取码错误或链接失效 (shareid=%d share_uk=%d fs_id=%d, 响应长度=%d, 前缀: %s)",
			len(shareIDs), len(userIDs), len(fsIDs), len(response), snippet)
	}

	shareID, _ := strconv.ParseInt(shareIDs[0][1], 10, 64)
	uk, _ := strconv.ParseInt(userIDs[0][1], 10, 64)

	var files []baiduShareFile
	for i, m := range fsIDs {
		fsID, _ := strconv.ParseInt(m[1], 10, 64)
		f := baiduShareFile{
			FsID:    fsID,
			ShareID: shareID,
			UK:      uk,
		}
		if i < len(filenames) {
			f.ServerName = filenames[i][1]
		}
		if i < len(isDirs) {
			f.IsDir = isDirs[i][1] == "1"
		}
		files = append(files, f)
	}
	return files, nil
}
