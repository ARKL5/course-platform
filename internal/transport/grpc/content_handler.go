package grpc

import (
	"context"
	"log"
	"strconv"

	"course-platform/internal/domain/content/model"
	"course-platform/internal/domain/content/service"
	"course-platform/internal/shared/pb/contentpb"
)

// ContentHandler 内容gRPC处理器
type ContentHandler struct {
	contentpb.UnimplementedContentServiceServer
	contentService service.ContentService
}

// NewContentHandler 创建内容处理器实例
func NewContentHandler(contentService service.ContentService) *ContentHandler {
	return &ContentHandler{
		contentService: contentService,
	}
}

// UploadFile 上传文件
func (h *ContentHandler) UploadFile(ctx context.Context, req *contentpb.UploadFileRequest) (*contentpb.UploadFileResponse, error) {
	log.Printf("📁 收到文件上传请求: %s", req.FileName)

	// 转换请求
	uploadReq := &service.UploadFileRequest{
		FileData:   req.FileData,
		FileName:   req.FileName,
		FileType:   req.FileType,
		CourseID:   uint(req.CourseId),
		UploaderID: uint(req.UploaderId),
	}

	// 调用服务层上传文件
	file, err := h.contentService.UploadFile(ctx, uploadReq)
	if err != nil {
		log.Printf("❌ 文件上传失败: %v", err)
		return &contentpb.UploadFileResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 转换文件信息 - 修复类型转换错误
	fileInfo := &contentpb.FileInfo{
		FileId:     strconv.FormatUint(uint64(file.ID), 10), // uint转换为string
		FileName:   file.FileName,
		FileUrl:    file.FileURL,
		FileType:   file.FileType,
		FileSize:   file.FileSize,
		CourseId:   uint32(file.CourseID),
		UploaderId: uint32(file.UploaderID),
		CreatedAt:  file.UploadTime.Format("2006-01-02 15:04:05"),
		UpdatedAt:  file.UploadTime.Format("2006-01-02 15:04:05"),
	}

	log.Printf("✅ 文件上传成功: %s", file.FileName)
	return &contentpb.UploadFileResponse{
		Code:     200,
		Message:  "文件上传成功",
		FileInfo: fileInfo,
	}, nil
}

// GetFiles 获取文件列表
func (h *ContentHandler) GetFiles(ctx context.Context, req *contentpb.GetFilesRequest) (*contentpb.GetFilesResponse, error) {
	log.Printf("📂 收到获取文件列表请求: 课程ID=%d", req.CourseId)

	// 设置默认分页参数
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 构建过滤器
	filter := &model.FileFilter{
		CourseID: uint(req.CourseId),
		FileType: req.FileType,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用服务层获取文件列表
	files, total, err := h.contentService.GetFiles(ctx, filter)
	if err != nil {
		log.Printf("❌ 获取文件列表失败: %v", err)
		return &contentpb.GetFilesResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	// 转换文件列表 - 修复类型转换错误
	pbFiles := make([]*contentpb.FileInfo, len(files))
	for i, file := range files {
		pbFiles[i] = &contentpb.FileInfo{
			FileId:     strconv.FormatUint(uint64(file.ID), 10), // uint转换为string
			FileName:   file.FileName,
			FileUrl:    file.FileURL,
			FileType:   file.FileType,
			FileSize:   file.FileSize,
			CourseId:   uint32(file.CourseID),
			UploaderId: uint32(file.UploaderID),
			CreatedAt:  file.UploadTime.Format("2006-01-02 15:04:05"),
			UpdatedAt:  file.UploadTime.Format("2006-01-02 15:04:05"),
		}
	}

	log.Printf("✅ 获取文件列表成功，共 %d 条记录", len(files))
	return &contentpb.GetFilesResponse{
		Code:    200,
		Message: "获取文件列表成功",
		Files:   pbFiles,
		Total:   uint32(total),
	}, nil
}

// DeleteFile 删除文件
func (h *ContentHandler) DeleteFile(ctx context.Context, req *contentpb.DeleteFileRequest) (*contentpb.DeleteFileResponse, error) {
	log.Printf("🗑️ 收到删除文件请求: 文件ID=%s, 用户ID=%d", req.FileId, req.UserId)

	// 修复string到uint的转换
	fileID, err := strconv.ParseUint(req.FileId, 10, 64)
	if err != nil {
		log.Printf("❌ 文件ID格式错误: %v", err)
		return &contentpb.DeleteFileResponse{
			Code:    400,
			Message: "文件ID格式错误",
		}, nil
	}

	// 调用服务层删除文件
	err = h.contentService.DeleteFile(ctx, uint(fileID), uint(req.UserId))
	if err != nil {
		log.Printf("❌ 删除文件失败: %v", err)
		return &contentpb.DeleteFileResponse{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	log.Printf("✅ 删除文件成功: 文件ID=%s", req.FileId)
	return &contentpb.DeleteFileResponse{
		Code:    200,
		Message: "文件删除成功",
	}, nil
}
