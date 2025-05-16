package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware 實現速率限制的 middleware
type RateLimitMiddleware struct {
	limit   int           // 限制數量
	window  time.Duration // 時間窗口
	storage RateLimitStorage
}

// MemoryRateLimitStorage 基於記憶體的速率限制儲存
type MemoryRateLimitStorage struct {
	mu      sync.RWMutex
	clients map[string]*clientRateInfo
}

type clientRateInfo struct {
	count   int
	resetAt time.Time
}

// NewRateLimitMiddleware 建立速率限制 middleware
func NewRateLimitMiddleware(limit int, window time.Duration, storage RateLimitStorage) Middleware {
	if storage == nil {
		storage = NewMemoryRateLimitStorage()
	}

	return &RateLimitMiddleware{
		limit:   limit,
		window:  window,
		storage: storage,
	}
}

// Handle 實現 Middleware 介面
func (m *RateLimitMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 取得客戶端識別 (IP 或 user ID)
		key := m.getClientKey(c)

		// 檢查速率限制
		count, resetAt, err := m.storage.GetAndIncrement(key, m.window)
		if err != nil {
			// 記錄錯誤但繼續處理
			c.Next()
			return
		}

		// 設置速率限制相關的 headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(m.limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(max(0, m.limit-count)))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

		// 檢查是否超過限制
		if count > m.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "速率限制：請求次數過多",
				"retry_after": resetAt.Sub(time.Now()).Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getClientKey 獲取客戶端識別鍵
func (m *RateLimitMiddleware) getClientKey(c *gin.Context) string {
	// 如果有使用者 ID，使用使用者 ID
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return "user:" + id
		}
	}

	// 否則使用 IP 地址
	return "ip:" + c.ClientIP()
}

// max 返回兩個整數中的較大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NewMemoryRateLimitStorage 建立基於記憶體的速率限制儲存
func NewMemoryRateLimitStorage() RateLimitStorage {
	storage := &MemoryRateLimitStorage{
		clients: make(map[string]*clientRateInfo),
	}

	// 啟動清理過期資料的 goroutine
	go storage.cleanup()

	return storage
}

// GetAndIncrement 獲取並增加計數
func (s *MemoryRateLimitStorage) GetAndIncrement(key string, window time.Duration) (int, time.Time, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	client, exists := s.clients[key]

	// 如果不存在或已過期，創建新的
	if !exists || now.After(client.resetAt) {
		resetAt := now.Add(window)
		s.clients[key] = &clientRateInfo{
			count:   1,
			resetAt: resetAt,
		}
		return 1, resetAt, nil
	}

	// 增加計數
	client.count++
	return client.count, client.resetAt, nil
}

// Reset 重置特定鍵的計數
func (s *MemoryRateLimitStorage) Reset(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, key)
	return nil
}

// cleanup 定期清理過期的資料
func (s *MemoryRateLimitStorage) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C

		s.mu.Lock()
		now := time.Now()
		for key, client := range s.clients {
			if now.After(client.resetAt) {
				delete(s.clients, key)
			}
		}
		s.mu.Unlock()
	}
}
