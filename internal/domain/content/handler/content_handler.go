package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	service "course-platform/internal/infrastructure/grpc_client"
	"course-platform/internal/shared/pb/contentpb"

	"github.com/gin-gonic/gin"
)

// ContentHandler 内容处理器
type ContentHandler struct {
	contentClient *service.ContentGRPCClientService
}

// NewContentHandler 创建内容处理器实例
func NewContentHandler(contentClient *service.ContentGRPCClientService) *ContentHandler {
	return &ContentHandler{
		contentClient: contentClient,
	}
}

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传课程相关文件（需要认证）
// @Tags content
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "上传的文件"
// @Param course_id formData string true "课程ID"
// @Param file_type formData string true "文件类型 (image, video, document, audio, other)"
// @Success 200 {object} map[string]interface{} "上传成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Failure 500 {object} map[string]interface{} "内部错误"
// @Router /api/v1/content/upload [post]
func (h *ContentHandler) UploadFile(c *gin.Context) {
	log.Printf("📁 收到文件上传请求")

	// 获取用户ID（从JWT中间件设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "AUTH_REQUIRED",
			"message": "用户未认证",
		})
		return
	}

	// 获取表单参数
	courseIDStr := c.PostForm("course_id")
	fileType := c.PostForm("file_type")

	if courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "MISSING_COURSE_ID",
			"message": "课程ID不能为空",
		})
		return
	}

	if fileType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "MISSING_FILE_TYPE",
			"message": "文件类型不能为空",
		})
		return
	}

	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_COURSE_ID",
			"message": "课程ID格式错误",
		})
		return
	}

	// 特殊处理：头像上传时course_id为0
	if courseID == 0 && fileType == "image" {
		log.Printf("📸 检测到头像上传请求")
	}

	// 获取上传文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "MISSING_FILE",
			"message": "请选择要上传的文件",
		})
		return
	}

	// 验证文件大小（限制50MB）
	const maxFileSize = 50 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "FILE_TOO_LARGE",
			"message": "文件大小不能超过50MB",
		})
		return
	}

	// 验证文件类型
	allowedTypes := map[string][]string{
		"image":    {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"},
		"video":    {".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm"},
		"document": {".pdf", ".doc", ".docx", ".ppt", ".pptx", ".txt", ".md"},
		"audio":    {".mp3", ".wav", ".flac", ".aac", ".ogg"},
		"other":    {}, // 允许所有类型
	}

	if fileType != "other" {
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if validExts, exists := allowedTypes[fileType]; exists && len(validExts) > 0 {
			isValid := false
			for _, validExt := range validExts {
				if ext == validExt {
					isValid = true
					break
				}
			}
			if !isValid {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "INVALID_FILE_TYPE",
					"message": fmt.Sprintf("文件类型 %s 不支持 %s 格式", fileType, ext),
				})
				return
			}
		}
	}

	// 读取文件数据
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FILE_READ_ERROR",
			"message": "读取文件失败",
		})
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "FILE_READ_ERROR",
			"message": "读取文件数据失败",
		})
		return
	}

	// 调用内容服务上传文件
	req := &contentpb.UploadFileRequest{
		FileName:   fileHeader.Filename,
		FileData:   fileData,
		FileType:   fileType,
		CourseId:   uint32(courseID),
		UploaderId: uint32(userID.(uint)),
	}

	resp, err := h.contentClient.UploadFile(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ 调用内容服务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "UPLOAD_FAILED",
			"message": "文件上传失败",
			"error":   err.Error(),
		})
		return
	}

	if resp.Code != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "UPLOAD_FAILED",
			"message": resp.Message,
		})
		return
	}

	log.Printf("✅ 文件上传成功: %s", fileHeader.Filename)
	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "文件上传成功",
		"data":    resp.FileInfo,
	})
}

// GetFiles 获取文件列表
// @Summary 获取文件列表
// @Description 获取课程文件列表
// @Tags content
// @Accept json
// @Produce json
// @Param course_id query string false "课程ID"
// @Param file_type query string false "文件类型"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 500 {object} map[string]interface{} "内部错误"
// @Router /api/v1/content/files [get]
func (h *ContentHandler) GetFiles(c *gin.Context) {
	log.Printf("📂 收到获取文件列表请求")

	// 获取查询参数
	courseIDStr := c.Query("course_id")
	fileType := c.Query("file_type")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	// 转换参数
	var courseID uint32 = 0
	if courseIDStr != "" {
		id, err := strconv.ParseUint(courseIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "INVALID_COURSE_ID",
				"message": "课程ID格式错误",
			})
			return
		}
		courseID = uint32(id)
	}

	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil || page == 0 {
		page = 1
	}

	pageSize, err := strconv.ParseUint(pageSizeStr, 10, 32)
	if err != nil || pageSize == 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	// 调用内容服务
	req := &contentpb.GetFilesRequest{
		CourseId: courseID,
		FileType: fileType,
		Page:     uint32(page),
		PageSize: uint32(pageSize),
	}

	resp, err := h.contentClient.GetFiles(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ 调用内容服务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "GET_FILES_FAILED",
			"message": "获取文件列表失败",
			"error":   err.Error(),
		})
		return
	}

	if resp.Code != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "GET_FILES_FAILED",
			"message": resp.Message,
		})
		return
	}

	log.Printf("✅ 获取文件列表成功，共 %d 条记录", len(resp.Files))
	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取文件列表成功",
		"data": gin.H{
			"files":     resp.Files,
			"total":     resp.Total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除文件（需要认证，只能删除自己上传的文件）
// @Tags content
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file_id path string true "文件ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "认证失败"
// @Failure 403 {object} map[string]interface{} "权限不足"
// @Failure 500 {object} map[string]interface{} "内部错误"
// @Router /api/v1/content/files/{file_id} [delete]
func (h *ContentHandler) DeleteFile(c *gin.Context) {
	log.Printf("🗑️ 收到删除文件请求")

	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "AUTH_REQUIRED",
			"message": "用户未认证",
		})
		return
	}

	// 获取文件ID
	fileIDStr := c.Param("file_id")
	if fileIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "MISSING_FILE_ID",
			"message": "文件ID不能为空",
		})
		return
	}

	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_FILE_ID",
			"message": "文件ID格式错误",
		})
		return
	}

	// 调用内容服务删除文件
	req := &contentpb.DeleteFileRequest{
		FileId: fileIDStr,
		UserId: uint32(userID.(uint)),
	}

	resp, err := h.contentClient.DeleteFile(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ 调用内容服务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "DELETE_FAILED",
			"message": "删除文件失败",
			"error":   err.Error(),
		})
		return
	}

	if resp.Code != 200 {
		status := http.StatusInternalServerError
		if strings.Contains(resp.Message, "无权限") {
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{
			"code":    "DELETE_FAILED",
			"message": resp.Message,
		})
		return
	}

	log.Printf("✅ 文件删除成功: 文件ID=%d", fileID)
	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "文件删除成功",
	})
}
