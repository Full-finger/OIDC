package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "http://localhost:8080"
)

type Anime struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverImage  string    `json:"cover_image"`
	ReleaseDate time.Time `json:"release_date"`
	Episodes    int       `json:"episodes"`
	Status      string    `json:"status"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func main() {
	fmt.Println("开始测试番剧增删改查功能...")

	// 1. 创建番剧
	fmt.Println("1. 测试创建番剧...")
	anime := Anime{
		Title:       "Test Anime CRUD",
		Description: "This is a test anime for CRUD operations.",
		CoverImage:  "https://example.com/test-anime-crud.jpg",
		ReleaseDate: time.Now(),
		Episodes:    12,
		Status:      "upcoming",
		Rating:      0.0,
	}

	createdAnime, err := createAnime(anime)
	if err != nil {
		fmt.Printf("创建番剧失败: %v\n", err)
		return
	}
	fmt.Printf("番剧创建成功，ID: %d\n", createdAnime.ID)

	// 2. 获取番剧
	fmt.Println("2. 测试获取番剧...")
	retrievedAnime, err := getAnimeByID(createdAnime.ID)
	if err != nil {
		fmt.Printf("获取番剧失败: %v\n", err)
		return
	}
	fmt.Printf("获取番剧成功，名称: %s\n", retrievedAnime.Title)

	// 3. 更新番剧
	fmt.Println("3. 测试更新番剧...")
	updatedAnime := *createdAnime
	updatedAnime.Status = "airing"
	updatedAnime.Rating = 8.5

	animeAfterUpdate, err := updateAnime(updatedAnime)
	if err != nil {
		fmt.Printf("更新番剧失败: %v\n", err)
		return
	}
	fmt.Printf("番剧更新成功，新状态: %s, 新评分: %.1f\n", animeAfterUpdate.Status, animeAfterUpdate.Rating)

	// 4. 列出所有番剧
	fmt.Println("4. 测试列出所有番剧...")
	animes, err := listAnimes()
	if err != nil {
		fmt.Printf("列出番剧失败: %v\n", err)
		return
	}
	fmt.Printf("列出番剧成功，共 %d 个番剧\n", len(animes))

	// 5. 搜索番剧
	fmt.Println("5. 测试搜索番剧...")
	searchResults, err := searchAnimes("CRUD")
	if err != nil {
		fmt.Printf("搜索番剧失败: %v\n", err)
		return
	}
	fmt.Printf("搜索番剧成功，找到 %d 个结果\n", len(searchResults))

	// 6. 删除番剧
	fmt.Println("6. 测试删除番剧...")
	err = deleteAnime(createdAnime.ID)
	if err != nil {
		fmt.Printf("删除番剧失败: %v\n", err)
		return
	}
	fmt.Println("番剧删除成功")

	// 7. 验证番剧已删除
	fmt.Println("7. 验证番剧已删除...")
	_, err = getAnimeByID(createdAnime.ID)
	if err != nil {
		fmt.Printf("验证番剧已删除: %v\n", err)
	} else {
		fmt.Println("错误：番剧未被删除")
	}

	fmt.Println("\n所有番剧增删改查功能测试完成！")
}

func createAnime(anime Anime) (*Anime, error) {
	jsonData, _ := json.Marshal(anime)
	resp, err := http.Post(baseURL+"/api/v1/anime/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("创建番剧失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	var createdAnime Anime
	if err := json.Unmarshal(body, &createdAnime); err != nil {
		return nil, err
	}

	return &createdAnime, nil
}

func getAnimeByID(id uint) (*Anime, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/anime/%d", baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取番剧失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	var anime Anime
	if err := json.Unmarshal(body, &anime); err != nil {
		return nil, err
	}

	return &anime, nil
}

func updateAnime(anime Anime) (*Anime, error) {
	jsonData, _ := json.Marshal(anime)
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/anime/%d", baseURL, anime.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("更新番剧失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	var updatedAnime Anime
	if err := json.Unmarshal(body, &updatedAnime); err != nil {
		return nil, err
	}

	return &updatedAnime, nil
}

func deleteAnime(id uint) error {
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/anime/%d", baseURL, id), nil)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("删除番剧失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	return nil
}

func listAnimes() ([]Anime, error) {
	resp, err := http.Get(baseURL + "/api/v1/anime/list")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("列出番剧失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	var animes []Anime
	if err := json.Unmarshal(body, &animes); err != nil {
		return nil, err
	}

	return animes, nil
}

func searchAnimes(keyword string) ([]Anime, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/anime/search?keyword=%s", baseURL, keyword))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("搜索番剧失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	var animes []Anime
	if err := json.Unmarshal(body, &animes); err != nil {
		return nil, err
	}

	return animes, nil
}

