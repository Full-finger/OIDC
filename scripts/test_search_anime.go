package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Anime struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	EpisodeCount *int   `json:"episode_count"`
	Director     *string `json:"director"`
}

func main() {
	fmt.Println("测试搜索番剧API")
	fmt.Println("================")

	// 测试1: 获取所有番剧
	fmt.Println("1. 获取所有番剧:")
	resp, err := http.Get("http://localhost:8080/api/v1/animes")
	if err != nil {
		log.Fatalf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var animes []Anime
		if err := json.NewDecoder(resp.Body).Decode(&animes); err != nil {
			log.Fatalf("解析响应失败: %v", err)
		}
		
		fmt.Printf("成功获取到 %d 部番剧:\n", len(animes))
		for _, anime := range animes {
			fmt.Printf("- %s (ID: %d)\n", anime.Title, anime.ID)
		}
	} else {
		fmt.Printf("获取番剧列表失败，状态码: %d\n", resp.StatusCode)
	}

	// 测试2: 搜索特定番剧
	fmt.Println("\n2. 搜索番剧 '千与千寻':")
	resp, err = http.Get("http://localhost:8080/api/v1/animes?title=千与千寻")
	if err != nil {
		log.Fatalf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var animes []Anime
		if err := json.NewDecoder(resp.Body).Decode(&animes); err != nil {
			log.Fatalf("解析响应失败: %v", err)
		}
		
		fmt.Printf("搜索到 %d 部匹配的番剧:\n", len(animes))
		for _, anime := range animes {
			fmt.Printf("- %s (ID: %d)\n", anime.Title, anime.ID)
			if anime.EpisodeCount != nil {
				fmt.Printf("  话数: %d\n", *anime.EpisodeCount)
			}
			if anime.Director != nil {
				fmt.Printf("  导演: %s\n", *anime.Director)
			}
		}
	} else {
		fmt.Printf("搜索番剧失败，状态码: %d\n", resp.StatusCode)
	}
}