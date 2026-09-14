// Package headers 构造三类上游请求头（common / chat / billing / refresh）。
// 规则来自 docs/api-reference.md §0/§4/§6。
package upstream

import (
	"net/http"
	"strings"

	"workbuddy2api/internal/auth"
)

const (
	// 国内版 WorkBuddy 桌面端官方默认版本与 CLI 版本
	defaultClientVersionCN = "5.5.6"
	defaultCliVersionCN    = "2.137.1"

	// 国际版 WorkBuddy AI 桌面端官方默认版本与 CLI 版本
	defaultClientVersionGlobal = "5.5.2"
	defaultCliVersionGlobal    = "5.5.2"

	originRefererCN     = "https://www.codebuddy.cn"
	// originRefererGlobal 国际版 Origin/Referer（与 backend 同域）。
	originRefererGlobal = "https://www.workbuddy.ai"
)

func originRefererFor(a *auth.Auth) string {
	if a != nil && a.IsGlobal() {
		return originRefererGlobal
	}
	return originRefererCN
}

func (c *Client) clientVersionCN() string {
	if c != nil && c.ClientVersion != "" {
		return c.ClientVersion
	}
	return defaultClientVersionCN
}

func (c *Client) cliVersionCN() string {
	if c != nil && c.CliVersion != "" {
		return c.CliVersion
	}
	return defaultCliVersionCN
}

func (c *Client) clientVersionGlobal() string {
	if c != nil && c.GlobalClientVersion != "" {
		return c.GlobalClientVersion
	}
	return defaultClientVersionGlobal
}

func (c *Client) cliVersionGlobal() string {
	if c != nil && c.GlobalClientVersion != "" {
		return c.GlobalClientVersion
	}
	return defaultCliVersionGlobal
}

func (c *Client) defaultWorkBuddyUACN() string {
	return "WorkBuddy/" + c.clientVersionCN() + " WorkBuddy/" + c.clientVersionCN() + " CLI/" + c.cliVersionCN()
}

func (c *Client) defaultWorkBuddyUAGlobal() string {
	return "WorkBuddy/" + c.clientVersionGlobal() + " WorkBuddy AI/" + c.clientVersionGlobal() + " CLI/" + c.cliVersionGlobal()
}

func (c *Client) userAgentFor(a *auth.Auth) string {
	if c != nil && c.UserAgent != "" {
		return c.UserAgent
	}
	if a != nil && a.IsGlobal() {
		return c.defaultWorkBuddyUAGlobal()
	}
	return c.defaultWorkBuddyUACN()
}

// CommonHeaders 设置所有 API 共享的请求头。
func (c *Client) CommonHeaders(req *http.Request, a *auth.Auth) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	origin := originRefererFor(a)
	req.Header.Set("Origin", origin)
	req.Header.Set("Referer", origin+"/")
	req.Header.Set("User-Agent", c.userAgentFor(a))
}

func (c *Client) injectAttribution(req *http.Request, a *auth.Auth) {
	req.Header.Set("X-Agent-Purpose", "conversation")
	if a != nil && a.IsGlobal() {
		req.Header.Set("X-IDE-Name", "WorkBuddy AI")
		req.Header.Set("X-IDE-Type", "WorkBuddy")
		req.Header.Set("X-IDE-Product", "WorkBuddy AI")
		req.Header.Set("X-IDE-Version", c.clientVersionGlobal())
		req.Header.Set("X-Product", "SaaS")
		req.Header.Set("X-Domain", "www.workbuddy.ai")
		req.Header.Set("X-No-Enterprise-Id", "1")
	} else {
		req.Header.Set("X-IDE-Name", "WorkBuddy")
		req.Header.Set("X-IDE-Type", "WorkBuddy")
		req.Header.Set("X-IDE-Product", "WorkBuddy")
		req.Header.Set("X-IDE-Version", c.clientVersionCN())
		req.Header.Set("X-Product", "WorkBuddy")
	}
}

func (c *Client) injectClientIP(req *http.Request, clientIP string) {
	if c == nil || !c.PassthroughIP || clientIP == "" {
		return
	}
	req.Header.Set("X-Forwarded-For", clientIP)
	req.Header.Set("X-Real-IP", clientIP)
	req.Header.Set("X-Client-IP", clientIP)
}

func ExtractClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return strings.TrimSpace(xff[:i])
			}
		}
		return strings.TrimSpace(xff)
	}
	if real := strings.TrimSpace(r.Header.Get("X-Real-IP")); real != "" {
		return real
	}
	return ""
}

// ChatHeaders 在 common 之上加 chat 专属的账号头。
// 缺省字段用 X-No-* 约定（与 CodeBuddy 官方 CLI 一致）。
func (c *Client) ChatHeaders(req *http.Request, a *auth.Auth, clientIP string) {
	c.CommonHeaders(req, a)
	if a.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+a.AccessToken)
	} else {
		req.Header.Set("X-No-Authorization", "1")
	}
	if a.UID != "" {
		req.Header.Set("X-User-Id", a.UID)
	} else {
		req.Header.Set("X-No-User-Id", "1")
	}
	if a != nil && a.IsGlobal() {
		req.Header.Set("X-No-Enterprise-Id", "1")
		req.Header.Set("X-Domain", "www.workbuddy.ai")
	} else {
		if a.EnterpriseID != "" {
			req.Header.Set("X-Enterprise-Id", a.EnterpriseID)
		} else {
			req.Header.Set("X-No-Enterprise-Id", "1")
		}
		if a.Domain != "" {
			req.Header.Set("X-Domain", a.Domain)
		} else {
			req.Header.Set("X-No-Department-Info", "1")
		}
	}
	c.injectAttribution(req, a)
	c.injectClientIP(req, clientIP)
}

// BillingHeaders billing 接口请求头。
func (c *Client) BillingHeaders(req *http.Request, a *auth.Auth) {
	req.Header.Set("Authorization", "Bearer "+a.AccessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if c != nil && c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	} else if a != nil && a.IsGlobal() {
		req.Header.Set("User-Agent", "WorkBuddy/"+c.clientVersionGlobal())
	} else {
		req.Header.Set("User-Agent", "WorkBuddy/"+c.clientVersionCN())
	}
	if a.UID != "" {
		req.Header.Set("X-User-Id", a.UID)
	}
	if a != nil && a.IsGlobal() {
		req.Header.Set("X-No-Enterprise-Id", "1")
		req.Header.Set("X-Domain", "www.workbuddy.ai")
	} else {
		if a.EnterpriseID != "" {
			req.Header.Set("X-Enterprise-Id", a.EnterpriseID)
			req.Header.Set("X-Tenant-Id", a.EnterpriseID)
		}
		if a.Domain != "" {
			req.Header.Set("X-Domain", a.Domain)
		}
	}
}

// RefreshHeaders refresh 端点专属头（X-Refresh-Token 只允许出现在这里）。
func (c *Client) RefreshHeaders(req *http.Request, a *auth.Auth) {
	c.CommonHeaders(req, a)
	req.Header.Set("X-Refresh-Token", a.RefreshToken)
	if a.EnterpriseID != "" {
		req.Header.Set("X-Enterprise-Id", a.EnterpriseID)
	}
	req.Header.Set("X-Auth-Refresh-Source", "workbuddy")
}
