package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"course-platform/internal/shared/pb/contentpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ContentGRPCClientService 内容服务gRPC客户端
type ContentGRPCClientService struct {
	client contentpb.ContentServiceClient
	conn   *grpc.ClientConn
}

// NewContentGRPCClientService 创建内容服务gRPC客户端
func NewContentGRPCClientService(address string) (*ContentGRPCClientService, error) {
	// 建立gRPC连接
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("连接内容服务失败: %w", err)
	}

	client := contentpb.NewContentServiceClient(conn)

	log.Printf("✅ 内容服务gRPC客户端已连接: %s", address)
	return &ContentGRPCClientService{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭连接
func (s *ContentGRPCClientService) Close() error {
	return s.conn.Close()
}

// UploadFile 上传文件
func (s *ContentGRPCClientService) UploadFile(ctx context.Context, req *contentpb.UploadFileRequest) (*contentpb.UploadFileResponse, error) {
	// 设置请求超时
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 调用gRPC服务
	resp, err := s.client.UploadFile(ctx, req)
	if err != nil {
		log.Printf("❌ 调用内容服务上传文件失败: %v", err)
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	log.Printf("📁 调用内容服务上传文件成功: %s", req.FileName)
	return resp, nil
}

// GetFiles 获取文件列表
func (s *ContentGRPCClientService) GetFiles(ctx context.Context, req *contentpb.GetFilesRequest) (*contentpb.GetFilesResponse, error) {
	// 设置请求超时
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 调用gRPC服务
	resp, err := s.client.GetFiles(ctx, req)
	if err != nil {
		log.Printf("❌ 调用内容服务获取文件列表失败: %v", err)
		return nil, fmt.Errorf("获取文件列表失败: %w", err)
	}

	log.Printf("📂 调用内容服务获取文件列表成功，共 %d 条记录", len(resp.Files))
	return resp, nil
}

// DeleteFile 删除文件
func (s *ContentGRPCClientService) DeleteFile(ctx context.Context, req *contentpb.DeleteFileRequest) (*contentpb.DeleteFileResponse, error) {
	// 设置请求超时
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 调用gRPC服务
	resp, err := s.client.DeleteFile(ctx, req)
	if err != nil {
		log.Printf("❌ 调用内容服务删除文件失败: %v", err)
		return nil, fmt.Errorf("删除文件失败: %w", err)
	}

	log.Printf("🗑️ 调用内容服务删除文件成功: 文件ID=%d", req.FileId)
	return resp, nil
}

