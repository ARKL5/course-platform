package template

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

// GetTemplateFunctions 返回所有模板函数的映射
func GetTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"iterate":          iterateFunction,
		"formatNumber":     formatNumberFunction,
		"add":              addFunction,
		"sub":              subFunction,
		"formatFileSize":   formatFileSizeFunction,
		"getFileType":      getFileTypeFunction,
		"getVideoDuration": getVideoDurationFunction,
	}
}

// iterateFunction 生成指定数量的整数切片
func iterateFunction(count int) []int {
	result := make([]int, count)
	for i := 0; i < count; i++ {
		result[i] = i
	}
	return result
}

// formatNumberFunction 格式化数字显示（如 1000 -> 1.0k）
func formatNumberFunction(num int) string {
	if num >= 1000 {
		return fmt.Sprintf("%.1fk", float64(num)/1000)
	}
	return fmt.Sprintf("%d", num)
}

// addFunction 两个数字相加
func addFunction(a, b int) int {
	return a + b
}

// subFunction 两个数字相减
func subFunction(a, b int) int {
	return a - b
}

// formatFileSizeFunction 格式化文件大小
func formatFileSizeFunction(size interface{}) string {
	// 如果size是字符串，直接返回
	if sizeStr, ok := size.(string); ok {
		return sizeStr
	}

	// 如果是数字，格式化为合适的单位
	if sizeNum, ok := size.(int64); ok {
		return formatBytesSize(sizeNum)
	}

	return "未知大小"
}

// formatBytesSize 将字节数格式化为人类可读的格式
func formatBytesSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// getFileTypeFunction 根据文件名判断文件类型
func getFileTypeFunction(filename string) string {
	if len(filename) == 0 {
		return "文件"
	}

	ext := strings.ToLower(filepath.Ext(filename))
	return getFileTypeByExtension(ext)
}

// getFileTypeByExtension 根据扩展名返回文件类型
func getFileTypeByExtension(ext string) string {
	videoExts := map[string]bool{
		".mp4": true, ".avi": true, ".mov": true, ".mkv": true,
		".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
	}

	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".bmp": true, ".svg": true, ".webp": true, ".ico": true,
	}

	documentExts := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".txt": true,
		".rtf": true, ".odt": true, ".pages": true,
	}

	presentationExts := map[string]bool{
		".ppt": true, ".pptx": true, ".odp": true, ".key": true,
	}

	audioExts := map[string]bool{
		".mp3": true, ".wav": true, ".flac": true, ".aac": true,
		".ogg": true, ".wma": true, ".m4a": true,
	}

	switch {
	case videoExts[ext]:
		return "视频"
	case imageExts[ext]:
		return "图片"
	case documentExts[ext]:
		return "文档"
	case presentationExts[ext]:
		return "演示文稿"
	case audioExts[ext]:
		return "音频"
	default:
		return "文件"
	}
}

// getVideoDurationFunction 模拟获取视频时长
func getVideoDurationFunction(filename string) string {
	// 模拟不同课程的视频时长
	durations := []string{
		"8:45", "12:30", "15:20", "9:15", "18:40",
		"11:25", "7:30", "14:55", "22:10", "6:35",
	}

	// 基于文件名生成一致的哈希值
	hash := calculateSimpleHash(filename)
	return durations[hash%len(durations)]
}

// calculateSimpleHash 计算字符串的简单哈希值
func calculateSimpleHash(s string) int {
	hash := 0
	for _, char := range s {
		hash = hash*31 + int(char)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}
