package utils

import (
	"context"
	"fmt"
	"github.com/hodukihugi/winglets-api/models"
	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"mime/multipart"
	"sync"
)

func UploadProfileImageAsync(
	wg *sync.WaitGroup,
	ctx context.Context,
	ch chan models.ImageUploadResult,
	ik *imagekit.ImageKit,
	src multipart.File,
	imageSlotId int,
	imageName string,
	folderName string,
) {
	defer wg.Done()
	base64Image, err := ConvertImageToBase64(src)
	if err != nil {
		return
	}

	// base64Image, err = ResizeBase64ImageTo50x50(base64Image)
	fmt.Printf("Start upload image: %s\n", imageName)

	uploadResponse, err := ik.Uploader.Upload(ctx, base64Image, uploader.UploadParam{
		FileName: fmt.Sprintf("%s.jpg", imageName),
		Folder:   folderName,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Done upload image: %s - Image Url: %s\n", imageName, uploadResponse.Data.Url)
	ch <- models.ImageUploadResult{SlotId: imageSlotId, FileId: uploadResponse.Data.FileId}
}
