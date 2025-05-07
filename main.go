package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/NEXORA-Studios/Nova.ModDeps/api_client"
	"github.com/NEXORA-Studios/Nova.ModDeps/cache"
	"github.com/NEXORA-Studios/Nova.ModDeps/storage"
)

// 应用程序配置
type AppConfig struct {
	APIClient *api_client.ModrinthClient
	Cache     *cache.Cache
	Storage   *storage.Storage
}

// 获取版本信息，使用三级查询逻辑：缓存 -> 数据库 -> API
func (app *AppConfig) getVersionInfo(ctx context.Context, versionID string) (*api_client.Version, string, error) {
	// 1. 首先检查缓存
	if cached, found := app.Cache.Get(versionID); found {
		return cached.(*api_client.Version), "cache", nil
	}

	// 2. 检查数据库
	dbVersion, err := app.Storage.GetVersion(ctx, versionID)
	if err == nil {
		// 更新缓存
		app.Cache.Set(versionID, dbVersion, 30*time.Minute)
		return dbVersion, "database", nil
	}

	// 3. 从API获取
	version, err := app.APIClient.GetVersion(ctx, versionID)
	if err != nil {
		return nil, "", err
	}

	// 保存到数据库
	if err := app.Storage.SaveVersion(ctx, version); err != nil {
		log.Printf("Failed to save version to database: %v", err)
	}

	// 保存到缓存
	app.Cache.Set(versionID, version, 30*time.Minute)
	return version, "api", nil
}

// HTTP处理函数：获取版本信息
func (app *AppConfig) handleGetVersion(w http.ResponseWriter, r *http.Request) {
	// 从URL路径中提取版本ID
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "无效的URL路径", http.StatusBadRequest)
		return
	}
	versionID := pathParts[3] // /api/version/{id} 中的 {id}
	if versionID == "" {
		http.Error(w, "版本ID不能为空", http.StatusBadRequest)
		return
	}

	version, source, err := app.getVersionInfo(r.Context(), versionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("获取版本信息失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 构建响应
	response := struct {
		Version *api_client.Version `json:"version"`
		Source  string              `json:"source"`
	}{
		Version: version,
		Source:  source,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 健康检查接口
func (app *AppConfig) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	// 初始化各模块
	apiClient := api_client.NewClient()
	cache := cache.NewCache()
	storage, err := storage.NewStorage("moddeps.db")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// 创建应用配置
	app := &AppConfig{
		APIClient: apiClient,
		Cache:     cache,
		Storage:   storage,
	}

	// 创建HTTP处理器
	mux := http.NewServeMux()

	// 注册API路由
	mux.HandleFunc("/api/health", app.handleHealthCheck)
	mux.HandleFunc("/api/version/", app.handleGetVersion)

	// 启动HTTP服务器
	port := ":8080"
	log.Printf("服务器启动在 http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
