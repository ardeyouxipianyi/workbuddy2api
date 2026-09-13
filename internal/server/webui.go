package server

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"workbuddy2api/internal/auth"
)

//go:embed dashboard.html
var dashboardHTML []byte

var (
	startTime      = time.Now()
	totalPrompt    atomic.Int64
	totalComplete  atomic.Int64
	totalCached    atomic.Int64
	todayTokens    atomic.Int64
	todayDateStr   string
	todayMu        sync.Mutex
	totalRequests  atomic.Int64
	totalSuccesses atomic.Int64
	totalErrors    atomic.Int64

	modelStatsMu sync.RWMutex
	modelStats   = make(map[string]*ModelUsage)

	metricsMu   sync.RWMutex
	metricsRing = make([]RequestMetric, 0, 100)
	metricIDSeq int64
)

type ModelUsage struct {
	Calls        int64 `json:"calls"`
	PromptToks   int64 `json:"prompt_tokens"`
	CompleteToks int64 `json:"complete_tokens"`
	CachedToks   int64 `json:"cached_tokens"`
}

type RequestMetric struct {
	ID           int64   `json:"id"`
	Time         string  `json:"time"`
	Model        string  `json:"model"`
	Mode         string  `json:"mode"`
	UID          string  `json:"uid"`
	Effort       string  `json:"effort,omitempty"`
	TTFBMs       int64   `json:"ttfb_ms"`
	TotalMs      int64   `json:"total_ms"`
	PromptToks   int     `json:"prompt_tokens"`
	CompleteToks int     `json:"complete_tokens"`
	CachedToks   int     `json:"cached_tokens"`
	CacheHitRate float64 `json:"cache_hit_rate"`
	SpeedTPS     float64 `json:"speed_tps"`
	Status       int     `json:"status"`
}

func RecordRequestMetric(model, mode, uid, effort string, ttfb, total time.Duration, status, promptToks, completeToks, cachedToks int) {
	totalRequests.Add(1)
	if status == 200 {
		totalSuccesses.Add(1)
	} else {
		totalErrors.Add(1)
	}

	if promptToks > 0 {
		totalPrompt.Add(int64(promptToks))
	}
	if completeToks > 0 {
		totalComplete.Add(int64(completeToks))
	}
	if cachedToks > 0 {
		totalCached.Add(int64(cachedToks))
	}

	sumToks := promptToks + completeToks
	if sumToks > 0 {
		todayMu.Lock()
		curDate := time.Now().Format("2006-01-02")
		if todayDateStr != curDate {
			todayDateStr = curDate
			todayTokens.Store(int64(sumToks))
		} else {
			todayTokens.Add(int64(sumToks))
		}
		todayMu.Unlock()

		modelStatsMu.Lock()
		m, ok := modelStats[model]
		if !ok {
			m = &ModelUsage{}
			modelStats[model] = m
		}
		m.Calls++
		m.PromptToks += int64(promptToks)
		m.CompleteToks += int64(completeToks)
		m.CachedToks += int64(cachedToks)
		modelStatsMu.Unlock()
	}

	metricsMu.Lock()
	defer metricsMu.Unlock()
	metricIDSeq++
	totalSec := total.Seconds()
	var tps float64
	if completeToks > 0 && totalSec > 0 {
		tps = float64(completeToks) / totalSec
	}
	ttfbMs := ttfb.Milliseconds()
	if ttfbMs <= 0 && total.Milliseconds() > 0 {
		ttfbMs = total.Milliseconds()
	}

	var hitRate float64
	if promptToks > 0 && cachedToks > 0 {
		hitRate = float64(cachedToks) / float64(promptToks) * 100.0
	}

	item := RequestMetric{
		ID:           metricIDSeq,
		Time:         time.Now().Format("15:04:05"),
		Model:        model,
		Mode:         mode,
		UID:          uid,
		Effort:       effort,
		TTFBMs:       ttfbMs,
		TotalMs:      total.Milliseconds(),
		PromptToks:   promptToks,
		CompleteToks: completeToks,
		CachedToks:   cachedToks,
		CacheHitRate: hitRate,
		SpeedTPS:     tps,
		Status:       status,
	}
	if len(metricsRing) >= 100 {
		metricsRing = metricsRing[1:]
	}
	metricsRing = append(metricsRing, item)
}

func (h *Handler) registerWebUIRoutes() {
	h.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(dashboardHTML)
	})
	h.mux.HandleFunc("GET /api/metrics", h.withAuth(h.apiMetrics))
	h.mux.HandleFunc("GET /api/gateway", h.withAuth(h.apiGatewayInfo))
	h.mux.HandleFunc("POST /api/auth/start", h.withAuth(h.apiAuthStart))
	h.mux.HandleFunc("POST /api/auth/poll", h.withAuth(h.apiAuthPoll))
	h.mux.HandleFunc("POST /api/account/delete", h.withAuth(h.apiAccountDelete))
	h.mux.HandleFunc("POST /api/account/revive", h.withAuth(h.apiAccountRevive))
}

func (h *Handler) apiMetrics(w http.ResponseWriter, r *http.Request) {
	metricsMu.RLock()
	records := make([]RequestMetric, len(metricsRing))
	for i := range metricsRing {
		records[i] = metricsRing[len(metricsRing)-1-i]
	}
	metricsMu.RUnlock()

	var ttfbList, totalList []int64
	for _, r := range records {
		if r.TTFBMs > 0 { ttfbList = append(ttfbList, r.TTFBMs) }
		if r.TotalMs > 0 { totalList = append(totalList, r.TotalMs) }
	}
	sort.Slice(ttfbList, func(i, j int) bool { return ttfbList[i] < ttfbList[j] })
	sort.Slice(totalList, func(i, j int) bool { return totalList[i] < totalList[j] })

	var p50TTFB, p90TTFB, p50Total, p90Total int64
	if len(ttfbList) > 0 {
		p50TTFB = ttfbList[len(ttfbList)*50/100]
		p90TTFB = ttfbList[len(ttfbList)*90/100]
	}
	if len(totalList) > 0 {
		p50Total = totalList[len(totalList)*50/100]
		p90Total = totalList[len(totalList)*90/100]
	}

	modelStatsMu.RLock()
	mDistribution := make(map[string]ModelUsage)
	for k, v := range modelStats {
		mDistribution[k] = *v
	}
	modelStatsMu.RUnlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	totP := totalPrompt.Load()
	totC := totalComplete.Load()
	totCached := totalCached.Load()
	var globalHitRate float64
	if totP > 0 {
		globalHitRate = float64(totCached) / float64(totP) * 100.0
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"tokens": map[string]any{
			"total":               totP + totC,
			"prompt_tokens":       totP,
			"completion_tokens":   totC,
			"cached_tokens":       totCached,
			"cache_hit_rate":      globalHitRate,
			"today":               todayTokens.Load(),
			"total_reqs":          totalRequests.Load(),
			"success_reqs":        totalSuccesses.Load(),
			"error_reqs":          totalErrors.Load(),
			"model_usage":         mDistribution,
		},
		"latency": map[string]any{
			"p50_ttfb":  p50TTFB,
			"p90_ttfb":  p90TTFB,
			"p50_total": p50Total,
			"p90_total": p90Total,
		},
		"system": map[string]any{
			"uptime_seconds":  int64(time.Since(startTime).Seconds()),
			"memory_alloc_mb": fmt.Sprintf("%.1f", float64(m.Alloc)/(1024*1024)),
			"goroutines":      runtime.NumGoroutine(),
		},
		"records": records,
	})
}

func (h *Handler) apiGatewayInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"api_key":      h.cfg.APIKey,
		"prompt_mode":  h.cfg.PromptMode,
		"max_body_mb":  h.cfg.MaxBodyBytes >> 20,
		"soft_rate":    h.cfg.SoftCooldown.String(),
		"redis_mode":   h.cfg.RedisMode,
		"session_mode": h.cfg.Session != nil,
	})
}

func (h *Handler) apiAuthStart(w http.ResponseWriter, r *http.Request) {
	realm := r.URL.Query().Get("realm")
	base := "https://copilot.tencent.com"
	origin := "https://www.codebuddy.cn"
	if realm == "global" {
		base = "https://www.workbuddy.ai"
		origin = "https://www.workbuddy.ai"
	}
	req, _ := http.NewRequest("POST", base+"/v2/plugin/auth/state?platform=CLI", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", origin)
	req.Header.Set("Referer", origin+"/")
	req.Header.Set("User-Agent", "CLI/2.63.2 CodeBuddy/2.63.2")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var env struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			State   string `json:"state"`
			AuthURL string `json:"authUrl"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &env); err != nil || env.Code != 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": env.Msg})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"state": env.Data.State, "auth_url": env.Data.AuthURL, "realm": realm})
}

func (h *Handler) apiAuthPoll(w http.ResponseWriter, r *http.Request) {
	var body struct { State string `json:"state"`; Realm string `json:"realm"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}
	base := "https://copilot.tencent.com"
	origin := "https://www.codebuddy.cn"
	if body.Realm == "global" {
		base = "https://www.workbuddy.ai"
		origin = "https://www.workbuddy.ai"
	}
	req, _ := http.NewRequest("GET", base+"/v2/plugin/auth/token?state="+body.State, nil)
	req.Header.Set("Origin", origin)
	req.Header.Set("User-Agent", "CLI/2.63.2 CodeBuddy/2.63.2")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var tokEnv struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int64  `json:"expires_in"`
			Domain       string `json:"domain"`
		} `json:"data"`
	}
	if json.Unmarshal(raw, &tokEnv) != nil || tokEnv.Code != 0 {
		writeJSON(w, http.StatusAccepted, map[string]any{"status": "pending", "msg": tokEnv.Msg})
		return
	}
	accReq, _ := http.NewRequest("GET", base+"/v2/plugin/login/account?state="+body.State, nil)
	accReq.Header.Set("Origin", origin)
	accReq.Header.Set("Authorization", "Bearer "+tokEnv.Data.AccessToken)
	accReq.Header.Set("User-Agent", "CLI/2.63.2 CodeBuddy/2.63.2")
	accResp, err := http.DefaultClient.Do(accReq)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	defer accResp.Body.Close()
	accRaw, _ := io.ReadAll(accResp.Body)
	var accEnv struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			UID          string `json:"uid"`
			EnterpriseID string `json:"enterprise_id"`
			Nickname     string `json:"nickname"`
		} `json:"data"`
	}
	if json.Unmarshal(accRaw, &accEnv) != nil || accEnv.Code != 0 || accEnv.Data.UID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "get uid failed: " + accEnv.Msg})
		return
	}
	authDir := "./auths"
	doc := map[string]any{
		"account": map[string]any{ "uid": accEnv.Data.UID, "enterpriseId": accEnv.Data.EnterpriseID, "nickname": accEnv.Data.Nickname },
		"auth": map[string]any{ "accessToken": tokEnv.Data.AccessToken, "refreshToken": tokEnv.Data.RefreshToken, "expiresAt": time.Now().Unix() + tokEnv.Data.ExpiresIn, "domain": tokEnv.Data.Domain },
	}
	targetFile := filepath.Join(authDir, fmt.Sprintf("workbuddy-%s.json", accEnv.Data.UID))
	content, _ := json.MarshalIndent(doc, "", "  ")
	if err := os.WriteFile(targetFile, content, 0644); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "save file failed: " + err.Error()})
		return
	}
	auths, err := auth.LoadDir(authDir)
	if err == nil { h.cfg.Pool.SyncToDir(auths) }
	writeJSON(w, http.StatusOK, map[string]any{"status": "success", "uid": accEnv.Data.UID, "nickname": accEnv.Data.Nickname})
}

func (h *Handler) apiAccountDelete(w http.ResponseWriter, r *http.Request) {
	var body struct { UID string `json:"uid"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.UID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing uid"})
		return
	}
	_ = os.Remove(filepath.Join("./auths", fmt.Sprintf("workbuddy-%s.json", body.UID)))
	auths, _ := auth.LoadDir("./auths")
	h.cfg.Pool.SyncToDir(auths)
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *Handler) apiAccountRevive(w http.ResponseWriter, r *http.Request) {
	var body struct { UID string `json:"uid"` }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.UID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "missing uid"})
		return
	}
	h.cfg.Pool.ReviveDisabled(body.UID)
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}
