package utils

import (
	"context"
	"fmt"
	"github.com/imagekit-developer/imagekit-go"
	"sync"
)

func RemoveProfileImageAsync(
	wg *sync.WaitGroup,
	ctx context.Context,
	ik *imagekit.ImageKit,
	imageId string,
) {
	defer wg.Done()

	fmt.Printf("Start remove image: %s\n", imageId)

	_, err := ik.Media.DeleteFile(ctx, imageId)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Done remove image: %s", imageId)
}
