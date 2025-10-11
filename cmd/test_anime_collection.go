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

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

type Anime struct {
	ID           uint   `json:"id"`
	BangumiID    int    `json:"bangumi_id"`
	Name         string `json:"name"`
	ChineseName  string `json:"chinese_name"`
	Summary      string `json:"summary"`
	AirDate      string `json:"air_date"`
	EpisodesCount int   `json:"episodes_count"`
	PosterURL    string `json:"poster_url"`
	Status       string `json:"status"`
}

type AddToCollectionRequest struct {
	AnimeID uint    `json:"anime_id"`
	Status  string  `json:"status"`
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}

type Collection struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	AnimeID     uint    `json:"anime_id"`
	Status      string  `json:"status"`
	Rating      float64 `json:"rating"`
	Comment     string  `json:"comment"`
	IsFavorite  bool    `json:"is_favorite"`
}

func main() {
	fmt.Println("开始测试番剧和收藏功能...")

	// 1. 用户注册
	fmt.Println("1. 测试用户注册...")
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	testUser := User{
		Username: "anime_user_" + timestamp,
		Email:    "anime_" + timestamp + "@example.com",
		Password: "password123",
		Nickname: "Anime Test User",
	}

	if err := registerUser(testUser); err != nil {
		fmt.Printf("注册失败: %v\n", err)
		return
	}
	fmt.Println("用户注册成功")

	// 2. 用户登录
	fmt.Println("2. 测试用户登录...")
	loginResp, err := loginUser(LoginRequest{
		Username: testUser.Username,
		Password: testUser.Password,
	})
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}
	fmt.Println("用户登录成功")
	fmt.Printf("访问令牌: %s\n", loginResp.AccessToken[:20]+"...")

	// 3. 获取番剧列表
	fmt.Println("3. 测试获取番剧列表...")
	animes, err := listAnimes(loginResp.AccessToken)
	if err != nil {
		fmt.Printf("获取番剧列表失败: %v\n", err)
		// 这可能是由于数据库中还没有番剧数据
		fmt.Println("这可能是由于数据库中还没有番剧数据")
	} else {
		fmt.Printf("获取番剧列表成功，共 %d 个番剧\n", len(animes))
		if len(animes) > 0 {
			// 显示第一个番剧的信息
			fmt.Printf("第一个番剧: ID=%d, 名称=%s\n", animes[0].ID, animes[0].Name)
		}
	}

	// 4. 搜索番剧
	fmt.Println("4. 测试搜索番剧...")
	searchResults, err := searchAnimes("测试", loginResp.AccessToken)
	if err != nil {
		fmt.Printf("搜索番剧失败: %v\n", err)
		// 这可能是由于数据库中还没有匹配的番剧数据
		fmt.Println("这可能是由于数据库中还没有匹配的番剧数据")
	} else {
		fmt.Printf("搜索番剧成功，找到 %d 个结果\n", len(searchResults))
	}

	// 5. 添加到收藏 (使用一个假设存在的番剧ID)
	fmt.Println("5. 测试添加到收藏...")
	collectionReq := AddToCollectionRequest{
		AnimeID: 1, // 假设番剧ID为1存在
		Status:  "Watching",
		Rating:  8.5,
		Comment: "Great anime!",
	}

	collection, err := addToCollection(collectionReq, loginResp.AccessToken)
	if err != nil {
		fmt.Printf("添加到收藏失败: %v\n", err)
		// 这可能是由于番剧ID不存在
		fmt.Println("这可能是由于番剧ID不存在")
	} else {
		fmt.Printf("添加到收藏成功，收藏ID: %d\n", collection.ID)
	}

	// 6. 获取用户收藏
	fmt.Println("6. 测试获取用户收藏...")
	userCollections, err := listUserCollections(loginResp.AccessToken)
	if err != nil {
		fmt.Printf("获取用户收藏失败: %v\n", err)
	} else {
		fmt.Printf("获取用户收藏成功，共 %d 个收藏\n", len(userCollections))
	}

	fmt.Println("\n番剧和收藏功能测试完成！")
}

func registerUser(user User) error {
	jsonData, _ := json.Marshal(user)
	resp, err := http.Post(baseURL+"/api/v1/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("注册失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	return nil
}

func loginUser(loginReq LoginRequest) (*LoginResponse, error) {
	jsonData, _ := json.Marshal(loginReq)
	resp, err := http.Post(baseURL+"/api/v1/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("登录失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, err
	}

	return &loginResp, nil
}

func listAnimes(token string) ([]Anime, error) {
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/anime/list", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取番剧列表失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	var animes []Anime
	if err := json.Unmarshal(body, &animes); err != nil {
		return nil, err
	}

	return animes, nil
}

func searchAnimes(keyword, token string) ([]Anime, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/anime/search?keyword=%s", baseURL, keyword), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
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

func addToCollection(req AddToCollectionRequest, token string) (*Collection, error) {
	jsonData, _ := json.Marshal(req)
	client := &http.Client{}
	httpReq, _ := http.NewRequest("POST", baseURL+"/api/v1/collection/", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("添加到收藏失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	var collection Collection
	if err := json.Unmarshal(body, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

func listUserCollections(token string) ([]Collection, error) {
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/collection/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取用户收藏失败，状态码: %d, 错误信息: %s", resp.StatusCode, string(body))
	}

	var collections []Collection
	if err := json.Unmarshal(body, &collections); err != nil {
		return nil, err
	}

	return collections, nil
}