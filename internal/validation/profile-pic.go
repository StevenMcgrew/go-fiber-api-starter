package validation

import (
	"mime/multipart"
	"net/http"
)

func IsProfilePicMimeTypeValid(fileHeader *multipart.FileHeader) (bool, error) {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 512 bytes to detect the MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err
	}

	// Detect the MIME type
	mimeType := http.DetectContentType(buffer)

	// Validate the MIME type (adjust this map to the types you want to allow)
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/bmp":  true,
	}
	if !allowedMimeTypes[mimeType] {
		return false, nil
	}
	return true, nil
}
