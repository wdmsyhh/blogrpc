package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func getOpenAPIFiles() ([]File, error) {
	//pname := fmt.Sprintf("%s/%s", "testgroup", "testprojectname")
	//
	//// GitLab API URL
	//gitlabAPIURL := fmt.Sprintf(`https://git.xxx.cn/api/v4/projects/%s/repository/tree?path=api/proto&recursive=true&page=%d&per_page=100&private_token=%s`, url.QueryEscape(pname), 1, "testtoken")
	//
	//// Create a new request
	//resp, err := http.Get(gitlabAPIURL)
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()
	//
	//var files []File
	//if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
	//	return nil, err
	//}

	// Filter out only .yaml files

	var files = []File{
		{
			Name: "name1.yaml",
			Path: "path1",
		},
		{
			Name: "name2.yaml",
			Path: "path2",
		},
	}
	var openAPIFiles []File
	for _, file := range files {
		if filepath.Ext(file.Name) == ".yaml" {
			openAPIFiles = append(openAPIFiles, file)
		}
	}

	return openAPIFiles, nil
}

func main() {
	r := gin.Default()

	// 提供 OpenAPI 文档
	r.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("templates/openapi.yaml") // 提供 openapi.yaml 文件
	})
	r.GET("/openapi2.yaml", func(c *gin.Context) {
		c.File("templates/openapi2.yaml") // 提供 openapi.yaml 文件
	})

	// 提供 Swagger UI
	r.GET("/swagger", func(c *gin.Context) {
		openAPIFiles, err := getOpenAPIFiles()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error retrieving OpenAPI files: %v", err)
			return
		}

		// 生成 Swagger UI 的 HTML 页面
		c.HTML(http.StatusOK, "swagger.html", gin.H{
			"files": openAPIFiles,
		})
	})

	// 加载 HTML 文件
	r.LoadHTMLGlob("templates/*")

	r.Run(":8080")
}
