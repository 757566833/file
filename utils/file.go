package utils

import (
	"mime/multipart"
	"net/http"
)

func GetFileContentType(out multipart.File) (string, error) {

	// 只需要前 512 个字节就可以了
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// 只识别一部分按规范的 https://github.com/golang/go/issues/47492#issuecomment-891320284
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}