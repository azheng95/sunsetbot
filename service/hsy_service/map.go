package hsy_service

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type MapReq struct {
	Region string `json:"region"`
	Event  string `json:"event"`
}

type MapResponse struct {
	MapImgSrc string `json:"map_img_src"`
	Status    string `json:"status"`
}

// 生成随机查询ID
func generateMapQueryID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Intn(1000) + 1)
}

// GetSunsetMapData 获取地图数据
func GetSunsetMapData(req MapReq) (*MapResponse, error) {
	queryID := generateMapQueryID()
	baseURL := "https://sunsetbot.top/map/"

	params := url.Values{}
	params.Add("query_id", queryID)
	params.Add("intend", "select_region")
	params.Add("event", req.Event)
	params.Add("region", req.Region)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回非200状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var data MapResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}
	return &data, nil
}
