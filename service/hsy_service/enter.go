package hsy_service

import (
	"encoding/json"
	"flame_clouds/config"
	"flame_clouds/global"
	"flame_clouds/service/sct_service"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SunsetBotReq struct {
	City  string  `json:"city"`
	Aod   float64 `json:"aod"`
	Event string  `json:"event"` // set_1:今天日落, rise_2:明天日出
}

type SunsetBotResponse struct {
	ImgHref     string `json:"img_href"`
	ImgSummary  string `json:"img_summary"`
	PlaceHolder string `json:"place_holder"`
	QueryId     string `json:"query_id"`
	Status      string `json:"status"`
	TbAod       string `json:"tb_aod"`        // 火烧云指标
	TbEventTime string `json:"tb_event_time"` // 事件时间
	TbQuality   string `json:"tb_quality"`
}

// 生成随机查询ID
func generateQueryID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Intn(10000000) + 1)
}

// GetSunsetData 获取日落/日出数据
func GetSunsetData(req SunsetBotReq) (*SunsetBotResponse, error) {
	queryID := generateQueryID()
	baseURL := "https://sunsetbot.top/"

	params := url.Values{}
	params.Add("query_id", queryID)
	params.Add("intend", "select_city")
	params.Add("query_city", req.City)
	params.Add("event_date", "None")
	params.Add("event", req.Event)
	params.Add("times", "None")

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

	var data SunsetBotResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}
	return &data, nil
}

// GetCitySunsetData 获取指定城市的天气数据
func GetCitySunsetData(e config.MonitorEvent) {
	t, err := GetSunsetData(SunsetBotReq{City: global.Config.Monitor.City, Event: e.EventType.Params(), Aod: e.CheckAod})
	if err != nil {
		logrus.Errorf("请求错误 %s", err)
		return
	}
	checkAndNotify(t, e)
}

// 解析火烧云指标
func parseAOD(aodStr string) (float64, error) {
	// 去除HTML标签和额外内容
	cleanStr := strings.Split(aodStr, "<")[0]
	cleanStr = strings.TrimSpace(cleanStr)
	return strconv.ParseFloat(cleanStr, 64)
}

// 检查并处理火烧云指标
func checkAndNotify(data *SunsetBotResponse, e config.MonitorEvent) {
	aod, err := parseAOD(data.TbAod)
	if err != nil {
		fmt.Printf("解析AOD失败: %v\n", err)
		return
	}

	logrus.Infof("城市: %s, 事件: %s, AOD: %.2f", global.Config.Monitor.City, e.EventType.String(), aod)

	// 阈值判断 (0.5是可配置的)
	if aod >= e.CheckAod {
		notification := sct_service.AlertNotification{
			City:      data.TbQuality,
			AOD:       aod,
			EventType: e.EventType.String(),
			EventTime: data.TbEventTime,
			ImageURL:  data.ImgHref,
		}

		if err := sct_service.SendServerNotification(notification); err != nil {
			logrus.Errorf("发送通知失败: %v", err)
		} else {
			logrus.Infof("火烧云预警通知已发送")
		}
	} else {
		logrus.Warnf("火烧云指标未达到阈值")
	}
}
