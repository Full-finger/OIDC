package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type TestAnime struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	EpisodeCount *int   `json:"episode_count"`
	Director     *string `json:"director"`
}

func TestCreateAnime() {
	fmt.Println("测试创建番剧API")
	fmt.Println("================")

	// 番剧数据
	episodeCount := 24
	director := "宫崎骏"
	animeData := map[string]interface{}{
		"title":         "千与千寻",
		"episode_count": episodeCount,
		"director":      director,
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(animeData)
	if err != nil {
		log.Fatalf("序列化JSON失败: %v", err)
	}

	// 发送POST请求
	url := "http://localhost:8080/api/v1/animes"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	
	// 注意：这里需要添加JWT token，由于我们没有实现管理员权限验证，
	// 暂时先不添加认证头，实际使用时需要添加有效的JWT token
	// req.Header.Set("Authorization", "Bearer YOUR_JWT_TOKEN_HERE")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 输出响应状态
	fmt.Printf("响应状态: %s\n", resp.Status)

	// 解析响应
	if resp.StatusCode == http.StatusCreated {
		var anime TestAnime
		if err := json.NewDecoder(resp.Body).Decode(&anime); err != nil {
			log.Fatalf("解析响应失败: %v", err)
		}
		
		fmt.Printf("成功创建番剧:\n")
		fmt.Printf("ID: %d\n", anime.ID)
		fmt.Printf("标题: %s\n", anime.Title)
		if anime.EpisodeCount != nil {
			fmt.Printf("话数: %d\n", *anime.EpisodeCount)
		}
		if anime.Director != nil {
			fmt.Printf("导演: %s\n", *anime.Director)
		}
	} else {
		fmt.Printf("创建番剧失败，状态码: %d\n", resp.StatusCode)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fmt.Printf("响应内容: %s\n", buf.String())
		os.Exit(1)
	}
}