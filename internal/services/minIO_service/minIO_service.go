package minioservice

import (
	"context"
	"ecom_promotion_v2/internal"
	"ecom_promotion_v2/internal/utils_call"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

type MergeImageRequest struct {
	BaseObjectName   string `json:"base_object"`
	FrameObjectName  string `json:"frame_object"`
	OutputObjectName string `json:"output_object"`
	Bucket           string `json:"bucket"`
}

type MergeImageService interface {
	MergeAndUploadImage(input *MergeImageRequest) (string, *internal.SystemStatus)
}

type mergeImageService struct{}

func NewMergeImageService() MergeImageService {
	return &mergeImageService{}
}

func (s *mergeImageService) MergeAndUploadImage(input *MergeImageRequest) (string, *internal.SystemStatus) {
	funcName := "MergeAndUploadImage"

	internal.Log.Info("📷 Bắt đầu ghép ảnh",
		zap.String("funcName", funcName),
		zap.Any("input", input),
	)

	ctx := context.Background()

	// Khởi tạo MinIO
	client, err := utils_call.InitMinioClient()
	if err != nil {
		internal.Log.Error("❌ Lỗi kết nối MinIO",
			zap.String("funcName", funcName),
			zap.Error(err),
		)
		return "", internal.SysStatus.SystemError
	}

	// Lấy ảnh gốc từ MinIO
	baseObject, err := client.GetObject(ctx, input.Bucket, input.BaseObjectName, minio.GetObjectOptions{})
	if err != nil {
		return handleMinIOError("BaseObject", input.BaseObjectName, funcName, err)
	}
	defer baseObject.Close()

	// Lấy ảnh khung từ MinIO
	frameObject, err := client.GetObject(ctx, input.Bucket, input.FrameObjectName, minio.GetObjectOptions{})
	if err != nil {
		return handleMinIOError("FrameObject", input.FrameObjectName, funcName, err)
	}
	defer frameObject.Close()

	// Decode ảnh
	baseImage, err := jpeg.Decode(baseObject)
	if err != nil {
		return handleDecodeError("JPEG Base Image", funcName, err)
	}

	frameImage, err := png.Decode(frameObject)
	if err != nil {
		return handleDecodeError("PNG Frame Image", funcName, err)
	}

	// Resize và ghép ảnh
	baseResized := resize.Resize(
		uint(frameImage.Bounds().Dx()),
		uint(frameImage.Bounds().Dy()),
		baseImage,
		resize.Lanczos3,
	)
	output := image.NewRGBA(frameImage.Bounds())
	draw.Draw(output, output.Bounds(), baseResized, image.Point{}, draw.Src)
	draw.Draw(output, output.Bounds(), frameImage, image.Point{}, draw.Over)

	outputObjectName := input.OutputObjectName
	if outputObjectName == "" {
		outputObjectName = fmt.Sprintf("output/merged_%d.jpg", time.Now().Unix())
	}

	// Định dạng file và Content-Type
	ext := filepath.Ext(outputObjectName)
	if ext == "" {
		ext = ".jpg"
		outputObjectName += ext
	}
	contentType := "image/jpeg"
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("merged_%d%s", time.Now().Unix(), ext))
	outFile, err := os.Create(tmpFile)
	if err != nil {
		return handleSystemError("Tạo file tạm thất bại", funcName, err)
	}
	defer outFile.Close()

	// Ghi file
	switch ext {
	case ".png":
		contentType = "image/png"
		err = png.Encode(outFile, output)
	default:
		err = jpeg.Encode(outFile, output, &jpeg.Options{Quality: 70})
	}
	if err != nil {
		return handleSystemError("Ghi file ảnh thất bại", funcName, err)
	}

	// Upload lên MinIO
	_, err = client.FPutObject(ctx, input.Bucket, outputObjectName, tmpFile, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return handleSystemError("Upload lên MinIO thất bại", funcName, err)
	}

	link := fmt.Sprintf("http://localhost:9000/%s/%s", input.Bucket, outputObjectName)
	internal.Log.Info("✅ Ghép và upload ảnh thành công", zap.String("link", link))

	return link, nil
}

func handleMinIOError(objectType, objectName, funcName string, err error) (string, *internal.SystemStatus) {
	internal.Log.Error(fmt.Sprintf("❌ Không thể đọc %s từ MinIO", objectType),
		zap.String("funcName", funcName),
		zap.String("objectName", objectName),
		zap.Error(err),
	)
	// go utils_call.SendStatusMessage(fmt.Sprintf("❌ Lỗi lấy object %s: %v", objectName, err), 3)

	return "", &internal.SystemStatus{
		Status: internal.CODE_SYSTEM_ERROR,
		Msg:    fmt.Sprintf("Cannot get %s", objectType),
		Detail: err.Error(),
	}
}

func handleDecodeError(info, funcName string, err error) (string, *internal.SystemStatus) {
	internal.Log.Error(fmt.Sprintf("❌ Decode lỗi: %s", info),
		zap.String("funcName", funcName),
		zap.Error(err),
	)
	return "", &internal.SystemStatus{
		Status: internal.CODE_SYSTEM_ERROR,
		Msg:    fmt.Sprintf("Decode error: %s", info),
		Detail: err.Error(),
	}
}

func handleSystemError(msg, funcName string, err error) (string, *internal.SystemStatus) {
	internal.Log.Error(msg,
		zap.String("funcName", funcName),
		zap.Error(err),
	)
	// go utils_call.SendStatusMessage(fmt.Sprintf("❌ %s: %v", msg, err), 3)

	return "", &internal.SystemStatus{
		Status: internal.CODE_SYSTEM_ERROR,
		Msg:    msg,
		Detail: err.Error(),
	}
}
