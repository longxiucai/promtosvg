package httpserver

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/longxiucai/promtosvg/pkg/chart"

	"github.com/gin-gonic/gin"
)

const FilePath = "./tmp/"

func RunWeb() error {
	r := gin.Default()

	r.LoadHTMLGlob("html/index.html")
	r.StaticFS("/css", http.Dir("./html/static/css"))
	r.StaticFS("/js", http.Dir("./html/static/js"))
	r.StaticFS("/tmp", http.Dir("./tmp"))

	r.POST("/uploadJsonFile", loadfile)
	r.POST("/uploadJson", loadjson)
	r.GET("/", index)

	err := r.Run(":8081")
	if err != nil {
		return err
	}
	return nil
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

type jsonData struct {
	Prom string `form:"prom" json:"prom" xml:"prom"  binding:"required"`
}

func loadjson(c *gin.Context) {
	var form jsonData
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	svgfile := chart.JsonToSvg([]byte(form.Prom))
	c.JSON(200, gin.H{
		"svgfile": svgfile,
	})
}

func loadfile(c *gin.Context) {
	file, _ := c.FormFile("file")
	hash := fmt.Sprintf("%x", getHash(file.Filename))
	dst := FilePath + hash + ".json"
	c.SaveUploadedFile(file, dst)
	svgfile := chart.OutPutSvgFromJson(dst)

	c.JSON(200, gin.H{
		"svgfile": svgfile,
	})
}

func getHash(s string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(s))
	hash := hasher.Sum(nil)
	return hash
}
