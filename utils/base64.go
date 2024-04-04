package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"regexp"
	"strings"

	"github.com/nfnt/resize"
)

// ResizeBase64ImageTo50x50 input format: "data:image/png;base64,iVBORw0KGg....="
func ResizeBase64ImageTo50x50(input string) (string, error) {
	seperatorIdx := strings.IndexByte(input, ',') + 1

	// prefix: "data:image/png;base64," , base64Str: "iVBORw0KGg....="
	prefix, base64Str := input[:seperatorIdx], input[seperatorIdx:]

	// Compile the regular expression
	regex := regexp.MustCompile(`data:image/([^;]+);base64,`)

	// Find the first match in the input string
	match := regex.FindStringSubmatch(prefix) // match: ["image/png", "png"]
	imageType := match[1]

	switch imageType {
	case "png":
		return resizeBase64PNGTo50x50(base64Str, prefix)
	case "jpeg":
		return resizeBase64JPEGTo50x50(base64Str, prefix)
	case "svg":
		return input, nil
	default:
		return "", errors.New(fmt.Sprintf(
			"Định dạng ảnh %s chưa được hỗ trợ. Vui lòng cung cấp ảnh có định dạng png, svg, jpg hoặc jpeg",
			imageType,
		))
	}
}

func resizeBase64PNGTo50x50(input, prefix string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(decoded)
	img, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}

	// Resize the image to 50x50
	newImg := resize.Resize(70, 70, img, resize.Lanczos3)

	// Encode the resized image back to Base64
	var buf bytes.Buffer
	if err = png.Encode(&buf, newImg); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s%s", prefix, base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

func resizeBase64JPEGTo50x50(input, prefix string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(decoded)
	img, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}

	// Resize the image to 50x50
	newImg := resize.Resize(70, 70, img, resize.Lanczos3)

	// Encode the resized image back to Base64
	var buf bytes.Buffer
	if err = jpeg.Encode(&buf, newImg, nil); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s%s", prefix, base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}
