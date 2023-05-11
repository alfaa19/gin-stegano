package controller

import (
	"bufio"
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alfaa19/gin-stegano/request"
	"github.com/auyer/steganography"
	"github.com/gin-gonic/gin"
)

func Encode(ctx *gin.Context) {
	var request request.ImageUploadRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	file, err := request.Image.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	outputFileName := request.Image.Filename
	outputFileType := filepath.Ext(outputFileName)
	outputFileName = "storage/" + outputFileName
	var img image.Image

	if outputFileType == ".jpg" {
		img, err = jpeg.Decode(file)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		img, _, err = image.Decode(file)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	w := new(bytes.Buffer) // Encode the message into the image
	if err := steganography.Encode(w, img, []byte(request.Message)); err != nil {
		log.Printf("Error Encoding file %v", err)
		return
	}
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	w.WriteTo(outputFile)
	defer outputFile.Close()
	ctx.JSON(http.StatusOK, gin.H{"success": outputFileName})

}

func Decode(ctx *gin.Context) {
	var request request.ImageUploadRequest
	var img image.Image

	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	file, err := request.Image.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r := bufio.NewReader(file) // buffer reader
	fileType := filepath.Ext(request.Image.Filename)
	if fileType == ".jpg" {
		img, err = jpeg.Decode(r)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		img, _, err = image.Decode(r)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	sizeOfMessage := steganography.GetMessageSizeFromImage(img)
	msg := steganography.Decode(sizeOfMessage, img)
	ctx.JSON(http.StatusOK, gin.H{
		"success": request.Image.Filename,
		"message": string(msg),
	})
}
