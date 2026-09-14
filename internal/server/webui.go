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
	h.mux.HandleFunc("POST /api/account/toggle", h.withAuth(h.apiAccountToggle))
	h.mux.HandleFunc("GET /api/settings", h.withAuth(h.apiSettingsGet))
	h.mux.HandleFunc("POST /api/settings", h.withAuth(h.apiSettingsSave))
	h.mux.HandleFunc("POST /api/account/revive", h.withAuth(h.apiAccountRevive))
	h.mux.HandleFunc("POST /api/schedule/checkin", h.withAuth(h.apiScheduleCheckin))
	h.mux.HandleFunc("POST /api/schedule/travel", h.withAuth(h.apiScheduleTravel))
	h.mux.HandleFunc("POST /api/schedule/activity", h.withAuth(h.apiScheduleActivity))
	h.mux.HandleFunc("POST /api/schedule/keepalive", h.withAuth(h.apiScheduleKeepalive))
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
	var body struct {
		State string `json:"state"`
		Realm string `json:"realm"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.State == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "state 参数不能为空"})
		return
	}

	base := "https://copilot.tencent.com"
	origin := "https://www.codebuddy.cn"
	if body.Realm == "global" {
		base = "https://www.workbuddy.ai"
		origin = "https://www.workbuddy.ai"
	}

	req, _ := http.NewRequest("GET", base+"/v2/plugin/auth/token?state="+body.State, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", origin)
	req.Header.Set("Referer", origin+"/")
	req.Header.Set("User-Agent", "CLI/2.63.2 CodeBuddy/2.63.2")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"status": "pending", "msg": "网络等待: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var tokEnv struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &tokEnv); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"status": "pending", "msg": "等待中..."})
		return
	}

	if tokEnv.Code != 0 || len(tokEnv.Data) == 0 {
		msg := tokEnv.Msg
		if msg == "" {
			msg = "等待用户在浏览器完成授权..."
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "pending", "msg": msg, "code": tokEnv.Code})
		return
	}

	var tok struct {
		AccessToken       string `json:"accessToken"`
		RefreshToken      string `json:"refreshToken"`
		ExpiresIn         int64  `json:"expiresIn"`
		Domain            string `json:"domain"`
		AccessTokenSnake  string `json:"access_token"`
		RefreshTokenSnake string `json:"refresh_token"`
		ExpiresInSnake    int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(tokEnv.Data, &tok); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "解析 token 失败: " + err.Error()})
		return
	}
	if tok.AccessToken == "" {
		tok.AccessToken = tok.AccessTokenSnake
	}
	if tok.RefreshToken == "" {
		tok.RefreshToken = tok.RefreshTokenSnake
	}
	if tok.ExpiresIn <= 0 {
		tok.ExpiresIn = tok.ExpiresInSnake
	}
	if tok.AccessToken == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "上游未返回有效 accessToken"})
		return
	}

	accReq, _ := http.NewRequest("GET", base+"/v2/plugin/login/account?state="+body.State, nil)
	accReq.Header.Set("Content-Type", "application/json")
	accReq.Header.Set("Accept", "application/json, text/plain, */*")
	accReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	accReq.Header.Set("Origin", origin)
	accReq.Header.Set("Referer", origin+"/")
	accReq.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	accReq.Header.Set("User-Agent", "CLI/2.63.2 CodeBuddy/2.63.2")

	accResp, err := client.Do(accReq)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "获取账号信息失败: " + err.Error()})
		return
	}
	defer accResp.Body.Close()
	accRaw, _ := io.ReadAll(accResp.Body)

	var accEnv struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			UID             string `json:"uid"`
			EnterpriseID    string `json:"enterpriseId"`
			EnterpriseIDAlt string `json:"enterprise_id"`
			Nickname        string `json:"nickname"`
		} `json:"data"`
	}
	if err := json.Unmarshal(accRaw, &accEnv); err != nil || accEnv.Code != 0 || accEnv.Data.UID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "获取 uid 失败: " + accEnv.Msg})
		return
	}

	uid := accEnv.Data.UID
	entID := accEnv.Data.EnterpriseID
	if entID == "" {
		entID = accEnv.Data.EnterpriseIDAlt
	}
	nickname := accEnv.Data.Nickname
	if nickname == "" {
		nickname = "用户-" + uid[:6]
	}

	authDir := "./auths"
	_ = os.MkdirAll(authDir, 0755)

	expiresAt := time.Now().Unix() + tok.ExpiresIn
	if tok.ExpiresIn <= 0 {
		expiresAt = time.Now().Unix() + 2592000
	}

	doc := map[string]any{
		"account": map[string]any{
			"uid":          uid,
			"enterpriseId": entID,
			"nickname":     nickname,
		},
		"auth": map[string]any{
			"accessToken":  tok.AccessToken,
			"refreshToken": tok.RefreshToken,
			"expiresAt":    expiresAt,
			"domain":       tok.Domain,
		},
	}
	targetFile := filepath.Join(authDir, fmt.Sprintf("workbuddy-%s.json", uid))
	content, _ := json.MarshalIndent(doc, "", "  ")
	if err := os.WriteFile(targetFile, content, 0644); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "保存凭据失败: " + err.Error()})
		return
	}

	auths, err := auth.LoadDir(authDir)
	if err == nil {
		h.cfg.Pool.SyncToDir(auths)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "success",
		"uid":      uid,
		"nickname": nickname,
	})
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

func (h *Handler) apiScheduleCheckin(w http.ResponseWriter, r *http.Request) {
	if h.cfg.RunCheckin != nil {
		go h.cfg.RunCheckin()
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "msg": "已触发签到与余额查询解冻（后台执行中）"})
	} else {
		writeJSON(w, http.StatusOK, map[string]any{"status": "disabled", "msg": "未启用签到任务"})
	}
}

func (h *Handler) apiScheduleTravel(w http.ResponseWriter, r *http.Request) {
	if h.cfg.RunTravel != nil {
		go h.cfg.RunTravel()
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "msg": "已触发猫猫旅行巡检（领养/派出/领奖后台执行中）"})
	} else {
		writeJSON(w, http.StatusOK, map[string]any{"status": "disabled", "msg": "未启用猫猫旅行任务"})
	}
}

func (h *Handler) apiScheduleActivity(w http.ResponseWriter, r *http.Request) {
	if h.cfg.RunActivity != nil {
		go h.cfg.RunActivity()
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "msg": "已触发每日活跃上报打卡（连登点亮中）"})
	} else {
		writeJSON(w, http.StatusOK, map[string]any{"status": "disabled", "msg": "未启用活跃上报任务"})
	}
}

func (h *Handler) apiScheduleKeepalive(w http.ResponseWriter, r *http.Request) {
	if h.cfg.RunKeepalive != nil {
		go h.cfg.RunKeepalive()
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "msg": "已触发全池 Token 保活刷新"})
	} else {
		writeJSON(w, http.StatusOK, map[string]any{"status": "disabled", "msg": "未启用 Token 保活任务"})
	}
}

func (h *Handler) apiAccountToggle(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UID     string `json:"uid"`
		Disable bool   `json:"disable"` // true: 开启养号 (暂停聊天挑号), false: 恢复接单
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.UID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "uid 不能为空"})
		return
	}
	if body.Disable {
		h.cfg.Pool.Pause(body.UID)
	} else {
		h.cfg.Pool.Resume(body.UID)
		h.cfg.Pool.ReviveDisabled(body.UID)
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "paused": body.Disable})
}

func (h *Handler) apiSettingsGet(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"api_key":       h.cfg.APIKey,
		"realm_pref":    h.cfg.Pool.RealmPreference(),
		"max_in_flight": h.cfg.Pool.MaxInFlight(),
		"prompt_mode":   h.cfg.PromptMode,
		"session_mode":  h.cfg.Session != nil,
	})
}

func (h *Handler) apiSettingsSave(w http.ResponseWriter, r *http.Request) {
	var body struct {
		APIKey      *string `json:"api_key"`
		RealmPref   *string `json:"realm_pref"`
		MaxInFlight *int    `json:"max_in_flight"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "无效的参数格式"})
		return
	}
	if body.APIKey != nil && *body.APIKey != "" {
		h.cfg.APIKey = *body.APIKey
	}
	if body.RealmPref != nil && *body.RealmPref != "" {
		h.cfg.Pool.SetRealmPreference(*body.RealmPref)
	}
	if body.MaxInFlight != nil && *body.MaxInFlight > 0 {
		h.cfg.Pool.SetMaxInFlight(*body.MaxInFlight)
	}

	saveConfigJSON(body.APIKey, body.MaxInFlight)
	writeJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"msg":        "网关配置与调度模式已成功更新并实时生效！",
		"api_key":    h.cfg.APIKey,
		"realm_pref": h.cfg.Pool.RealmPreference(),
	})
}

func saveConfigJSON(apiKey *string, maxInFlight *int) {
	cfgPaths := []string{"config.json", "/app/config.json"}
	for _, fp := range cfgPaths {
		raw, err := os.ReadFile(fp)
		if err != nil {
			continue
		}
		var m map[string]any
		if json.Unmarshal(raw, &m) != nil {
			continue
		}
		if apiKey != nil && *apiKey != "" {
			m["api_key"] = *apiKey
		}
		if maxInFlight != nil && *maxInFlight > 0 {
			if poolObj, ok := m["pool"].(map[string]any); ok {
				poolObj["max_in_flight"] = *maxInFlight
			}
		}
		out, err := json.MarshalIndent(m, "", "  ")
		if err == nil {
			_ = os.WriteFile(fp, out, 0644)
		}
	}
}
