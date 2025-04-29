package handler

import (
	"github.com/gin-gonic/gin"
	"krillin-ai/internal/dto"
	"krillin-ai/internal/response"
	"net/http"
	"os"
	"path/filepath"
)

func (h Handler) StartSubtitleTask(c *gin.Context) {
	var req dto.StartVideoSubtitleTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.R(c, response.Response{
			Error: -1,
			Msg:   "参数错误",
			Data:  nil,
		})
		return
	}

	svc := h.Service

	data, err := svc.StartSubtitleTask(req)
	if err != nil {
		response.R(c, response.Response{
			Error: -1,
			Msg:   err.Error(),
			Data:  nil,
		})
		return
	}
	response.R(c, response.Response{
		Error: 0,
		Msg:   "成功",
		Data:  data,
	})
}

func (h Handler) GetSubtitleTask(c *gin.Context) {
	var req dto.GetVideoSubtitleTaskReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.R(c, response.Response{
			Error: -1,
			Msg:   "参数错误",
			Data:  nil,
		})
		return
	}
	svc := h.Service
	data, err := svc.GetTaskStatus(req)
	if err != nil {
		response.R(c, response.Response{
			Error: -1,
			Msg:   err.Error(),
			Data:  nil,
		})
		return
	}
	response.R(c, response.Response{
		Error: 0,
		Msg:   "成功",
		Data:  data,
	})
}

func (h Handler) UploadFile(c *gin.Context) {
	// 获取文件和哈希值
	hash := c.PostForm("hash")
	if hash == "" {
		response.R(c, response.Response{
			Error: -1,
			Msg:   "缺少文件哈希值",
			Data:  nil,
		})
		return
	}

	// 检查文件是否已存在
	savePath := "./uploads/" + hash
	if _, err := os.Stat(savePath); err == nil {
		// 文件已存在
		response.R(c, response.Response{
			Error: 0,
			Msg:   "文件已存在",
			Data:  gin.H{"file_path": "local:" + savePath},
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.R(c, response.Response{
			Error: -1,
			Msg:   "未能获取文件",
			Data:  nil,
		})
		return
	}

	// 保存文件
	if err = c.SaveUploadedFile(file, savePath); err != nil {
		response.R(c, response.Response{
			Error: -1,
			Msg:   "文件保存失败",
			Data:  nil,
		})
		return
	}

	response.R(c, response.Response{
		Error: 0,
		Msg:   "文件上传成功",
		Data:  gin.H{"file_path": "local:" + savePath},
	})

}

func (h Handler) CheckFileExist(c *gin.Context) {
	hash := c.Query("hash")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"exists": false})
		return
	}
	savePath := "./uploads/" + hash
	if _, err := os.Stat(savePath); err == nil {
		c.JSON(http.StatusOK, gin.H{"exists": true, "path": "local:" + savePath})
	} else {
		c.JSON(http.StatusOK, gin.H{"exists": false})
	}
}

func (h Handler) DownloadFile(c *gin.Context) {
	requestedFile := c.Param("filepath")
	if requestedFile == "" {
		response.R(c, response.Response{
			Error: -1,
			Msg:   "文件路径为空",
			Data:  nil,
		})
		return
	}

	localFilePath := filepath.Join(".", requestedFile)
	if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
		response.R(c, response.Response{
			Error: -1,
			Msg:   "文件不存在",
			Data:  nil,
		})
		return
	}
	c.FileAttachment(localFilePath, filepath.Base(localFilePath))
}
