package routes

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

// This function will handle the file uploads and will only work if the AWS environment variables are set:
// AWS_S3_BUCKET, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_REGION

func uploadFiles(ctx *gin.Context) {
	region := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}
	svc := s3.New(sess)

	bucket := os.Getenv("AWS_S3_BUCKET")

	err = ctx.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get all files from the form
	form := ctx.Request.MultipartForm
	files := form.File["files"]

	fmt.Println("Uploading", len(files), "files")
	fmt.Println("Uploading files", files)

	// Iterate over the files and save them
	for _, file := range files {
		// Check if the file is an image, photo, or video
		contentType := getContentType(file)
		if !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "video/") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type", "filename": file.Filename, "content_type": contentType})
			return
		}

		fmt.Println("Uploading file:", file.Filename)
		// Open the file
		src, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		// Upload the file to AWS S3
		err = uploadToAWS(svc, bucket, file.Filename, src)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Save the file to the server

		// err = ctx.SaveUploadedFile(file, "uploads/"+file.Filename)
		// if err != nil {
		// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// 	return
		// }
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Files uploaded successfully"})
}

func getContentType(file *multipart.FileHeader) string {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return ""
	}
	defer src.Close()

	// Get the first 512 bytes to sniff the content type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return ""
	}

	// Reset the file position
	src.Seek(0, 0)

	// Sniff the content type
	mimeType := http.DetectContentType(buffer)
	return mimeType
}

func uploadToAWS(svc *s3.S3, bucket, key string, file multipart.File) error {
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("images/" + key),
		Body:   file,
	})
	if err != nil {
		return err
	}
	return nil
}
