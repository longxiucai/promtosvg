package chart

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/common/model"
	"github.com/wcharczuk/go-chart"
)

// QueryResult 结构体用于映射 Prometheus 中的 model.Matrix
type QueryResult struct {
	Status string    `json:"status"`
	Data   QueryData `json:"data"`
}

// QueryData 结构体用于映射 Prometheus 中的 model.Matrix 中的 data 字段
type QueryData struct {
	ResultType model.ValueType `json:"resultType"`
	Result     model.Matrix    `json:"result"`
}

func OutPutSvgFromJson(path string) string {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	return JsonToSvg(jsonData)
}

func JsonToSvg(jsonData []byte) string {
	// 解码 JSON 数据到 QueryResult 结构体
	var queryResult QueryResult
	err := json.Unmarshal(jsonData, &queryResult)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}

	// 提取时间序列数据
	var s []chart.Series
	for _, result := range queryResult.Data.Result {
		var ts []time.Time
		var vs []float64
		// 以当前示例中的数据结构为例，第一个元素是时间戳，第二个元素是值
		for _, value := range result.Values {
			vs = append(vs, float64(value.Value))
			ts = append(ts, value.Timestamp.Time())
		}
		s1 := chart.TimeSeries{
			Name:    fmt.Sprintf("%v", result.Metric["pod"]),
			XValues: ts,
			Style: chart.Style{
				Show: true,
			},
			YValues: vs,
		}
		s = append(s, s1)
	}

	// 创建折线图
	graph := chart.Chart{
		Series: s,
		Width:  3000,
		Height: 800,
	}
	graph.XAxis.Style = chart.Style{
		Show: true,
	}
	graph.YAxis.Style = chart.Style{
		Show: true,
	}
	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph, chart.Style{
			Padding: chart.Box{},
		}),
	}
	hs := md5.New()
	fileName := fmt.Sprintf("%x", hs.Sum(jsonData))
	fileName = fileName[:20]
	svgpath := filepath.Join("tmp", fileName) + ".svg"

	file, _ := os.Create(svgpath)
	// 保存图表为 PNG 文件
	defer file.Close()
	err = graph.Render(chart.SVG, file)
	if err != nil {
		fmt.Println("Error rendering graph:", err)
		return ""
	}
	return svgpath
}
