package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Body struct {
	// json tag to de-serialize json body
	Name string `json:"name"`
}
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FileName string `json:"filename"`
	FileData string `json:"filedata,omitempty"`
}

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "penis",
		})
	})
	r.POST("/test", func(context *gin.Context) {
		body := Body{}
		// using BindJson method to serialize body with struct
		if err := context.BindJSON(&body); err != nil {
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}
		fmt.Println(body)
		context.JSON(http.StatusAccepted, &body)
	})
	r.POST("/upload", handleFileUpload2)
	r.DELETE("/delete", deleteFile)

	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func deleteFile(c *gin.Context) {
	desktopDir := filepath.Join("C:", "Users", "Dylan", "Desktop", "UploadedFiles")
	filename := c.Query("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'filename' parameter"})
		return
	}

	filePath := filepath.Join(desktopDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	err := os.Remove(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete file: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully", "filename": filename})

}

//	func handleGetFiles(c *gin.Context) {
//		desktopDir := filepath.Join("C:", "Users", "Dylan", "Desktop", "UploadedFiles")
//	}
func handleFileUpload(c *gin.Context) {
	// Process the uploaded text file (if any)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Ensure the uploaded file is a text file
	ext := filepath.Ext(file.Filename)
	if ext != ".txt" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only text files (.txt) are allowed."})
		return
	}

	// Generate a unique filename for the uploaded file
	filename := filepath.Base(file.Filename)

	// Save the uploaded file to a specific location (e.g., the current directory)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	desktopDir := filepath.Join("C:", "Users", "Dylan", "Desktop", "UploadedFiles")
	if err := os.MkdirAll(desktopDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}
	filePath := filepath.Join(desktopDir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Read the content of the text file
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	// Respond with a success message and the file data
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": filename, "filedata": string(fileContent)})
}
func handleFileUpload2(c *gin.Context) {
	// Process the uploaded text file (if any)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Ensure the uploaded file is a text file
	// ext := filepath.Ext(file.Filename)
	// if ext != ".txt" && ext != ".pdf" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only text files (.txt) are allowed."})
	// 	return
	// }

	// Generate a unique filename for the uploaded file
	filename := filepath.Base(file.Filename)

	// Specify the directory path on your desktop where you want to save the file
	desktopDir := "C:\\Users\\Dylan\\Desktop\\UploadedFiles"

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(desktopDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Save the uploaded file to the specified directory
	filePath := filepath.Join(desktopDir, filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	if err := os.MkdirAll(desktopDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory", "details": err.Error()})
		return
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file", "details": err.Error()})
		return
	}

	// Respond with a success message and the file data
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": filename, "filepath": filePath})
}
