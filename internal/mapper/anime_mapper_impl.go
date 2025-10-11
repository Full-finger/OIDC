package mapper

import (
	"errors"
	"strings"
	"sync"
	"time"
	
	"github.com/Full-finger/OIDC/internal/model"
)

// animeMapper 畠剧映射器实现
type animeMapper struct {
	// 使用内存存储进行测试
	mu     sync.RWMutex
	animes map[uint]*model.Anime
	nextID uint
}

// NewAnimeMapper 创建AnimeMapper实例
func NewAnimeMapper() AnimeMapper {
	mapper := &animeMapper{
		animes: make(map[uint]*model.Anime),
		nextID: 1,
	}
	
	// 添加一些测试数据
	mapper.initializeTestData()
	
	return mapper
}

// Save 保存番剧
func (m *animeMapper) Save(entity interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	anime, ok := entity.(*model.Anime)
	if !ok {
		return errors.New("invalid anime entity")
	}
	
	// 如果是新番剧，分配ID
	if anime.ID == 0 {
		anime.ID = m.nextID
		m.nextID++
	}
	
	// 保存番剧
	m.animes[anime.ID] = anime
	
	return nil
}

// DeleteByID 根据ID删除番剧
func (m *animeMapper) DeleteByID(id interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	animeID, ok := id.(uint)
	if !ok {
		return errors.New("invalid anime id")
	}
	
	delete(m.animes, animeID)
	return nil
}

// GetByID 根据ID获取番剧
func (m *animeMapper) GetByID(id interface{}) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	animeID, ok := id.(uint)
	if !ok {
		return nil, errors.New("invalid anime id")
	}
	
	anime, exists := m.animes[animeID]
	if !exists {
		return nil, errors.New("anime not found")
	}
	
	return anime, nil
}

// GetAll 获取所有番剧
func (m *animeMapper) GetAll() ([]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	animes := make([]interface{}, 0, len(m.animes))
	for _, anime := range m.animes {
		animes = append(animes, anime)
	}
	
	return animes, nil
}

// Update 更新番剧
func (m *animeMapper) Update(entity interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	anime, ok := entity.(*model.Anime)
	if !ok {
		return errors.New("invalid anime entity")
	}
	
	if anime.ID == 0 {
		return errors.New("anime id is required")
	}
	
	// 检查番剧是否存在
	if _, exists := m.animes[anime.ID]; !exists {
		return errors.New("anime not found")
	}
	
	// 更新番剧
	m.animes[anime.ID] = anime
	
	return nil
}

// GetByTitle 根据标题获取番剧
func (m *animeMapper) GetByTitle(title string) (*model.Anime, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for _, anime := range m.animes {
		if anime.Title == title {
			return anime, nil
		}
	}
	
	return nil, errors.New("anime not found")
}

// GetByStatus 根据状态获取番剧列表
func (m *animeMapper) GetByStatus(status string) ([]*model.Anime, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var animes []*model.Anime
	for _, anime := range m.animes {
		if anime.Status == status {
			animes = append(animes, anime)
		}
	}
	
	return animes, nil
}

// Search 搜索番剧
func (m *animeMapper) Search(keyword string) ([]*model.Anime, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var animes []*model.Anime
	for _, anime := range m.animes {
		// 简单的关键词匹配
		if contains(anime.Title, keyword) || contains(anime.Description, keyword) {
			animes = append(animes, anime)
		}
	}
	
	return animes, nil
}

// initializeTestData 初始化测试数据
func (m *animeMapper) initializeTestData() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	testAnimes := []*model.Anime{
		{
			ID:          1,
			Title:       "Test Anime 1",
			Description: "This is a test anime for testing purposes.",
			CoverImage:  "https://example.com/poster1.jpg",
			ReleaseDate: time.Now(),
			Episodes:    12,
			Status:      "airing",
			Rating:      8.5,
		},
		{
			ID:          2,
			Title:       "Test Anime 2",
			Description: "This is another test anime for testing purposes.",
			CoverImage:  "https://example.com/poster2.jpg",
			ReleaseDate: time.Now(),
			Episodes:    24,
			Status:      "finished",
			Rating:      9.0,
		},
	}
	
	for _, anime := range testAnimes {
		m.animes[anime.ID] = anime
	}
	
	m.nextID = 3
}

// contains 检查字符串是否包含子串（不区分大小写）
func contains(s, substr string) bool {
	// 简单实现，实际应用中可能需要更复杂的匹配逻辑
	return len(s) >= len(substr) && 
		(strings.Contains(strings.ToLower(s), strings.ToLower(substr)))
}