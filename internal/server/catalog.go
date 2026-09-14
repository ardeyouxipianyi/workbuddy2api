package server

import "strings"

type ModelCatalogEntry struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Realm             string   `json:"realm"` // "cn" | "global"
	RealmLabel        string   `json:"realm_label"`
	Description       string   `json:"description"`
	Credits           string   `json:"credits"`
	MaxInputTokens    int      `json:"max_input_tokens"`
	MaxOutputTokens   int      `json:"max_output_tokens"`
	ContextLength     int      `json:"context_length"`
	SupportsImages    bool     `json:"supports_images"`
	SupportsTools     bool     `json:"supports_tools"`
	SupportsReasoning bool     `json:"supports_reasoning"`
	ReasoningEfforts  []string `json:"reasoning_efforts,omitempty"`
}

// 国内版真实模型清单（共 29 款）
var cnCatalog = map[string]ModelCatalogEntry{
	"hy3": {
		ID: "hy3", Name: "Hy3", Realm: "cn", RealmLabel: "国内版",
		Description: "混元思考模型，具有增强的推理能力", Credits: "0.00x (限时免费)",
		MaxInputTokens: 192000, MaxOutputTokens: 64000, ContextLength: 192000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high"},
	},
	"hy3-x": {
		ID: "hy3-x", Name: "Hy3", Realm: "cn", RealmLabel: "国内版",
		Description: "混元思考模型，具有增强的推理能力", Credits: "0.05x",
		MaxInputTokens: 192000, MaxOutputTokens: 64000, ContextLength: 192000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high"},
	},
	"hy4-preview": {
		ID: "hy4-preview", Name: "Hy4 preview", Realm: "cn", RealmLabel: "国内版",
		Description: "混元思考模型，具有增强的推理能力", Credits: "0.29x (夜间免费)",
		MaxInputTokens: 1000000, MaxOutputTokens: 64000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"high"},
	},
	"hy4-preview-x": {
		ID: "hy4-preview-x", Name: "Hy4 preview", Realm: "cn", RealmLabel: "国内版",
		Description: "混元思考模型，具有增强的推理能力", Credits: "x0.29",
		MaxInputTokens: 1000000, MaxOutputTokens: 64000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"high"},
	},
	"minimax-m2.5": {
		ID: "minimax-m2.5", Name: "MiniMax-M2.5", Realm: "cn", RealmLabel: "国内版",
		Description: "能力均衡，适合日常使用", Credits: "x0.18",
		MaxInputTokens: 200000, MaxOutputTokens: 48000, ContextLength: 200000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"glm-5v-turbo": {
		ID: "glm-5v-turbo", Name: "GLM-5v-Turbo", Realm: "cn", RealmLabel: "国内版",
		Description: "原生多模态模型", Credits: "0.71x",
		MaxInputTokens: 200000, MaxOutputTokens: 64000, ContextLength: 200000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"glm-5.3": {
		ID: "glm-5.3", Name: "GLM-5.3", Realm: "cn", RealmLabel: "国内版",
		Description: "能力均衡，适合日常使用", Credits: "0.79x",
		MaxInputTokens: 1000000, MaxOutputTokens: 48000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high", "max"},
	},
	"glm-5.3-flash": {
		ID: "glm-5.3-flash", Name: "GLM-5.3-Flash", Realm: "cn", RealmLabel: "国内版",
		Description: "原生多模态，擅长处理复杂的长程自主任务。", Credits: "0.06x",
		MaxInputTokens: 1000000, MaxOutputTokens: 32000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high", "max"},
	},
	"glm-5.2": {
		ID: "glm-5.2", Name: "GLM-5.2", Realm: "cn", RealmLabel: "国内版",
		Description: "1M 上下文，擅长长程任务", Credits: "0.79x (夜间折扣)",
		MaxInputTokens: 1000000, MaxOutputTokens: 48000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"high", "xhigh"},
	},
	"glm-5.1": {
		ID: "glm-5.1", Name: "GLM-5.1", Realm: "cn", RealmLabel: "国内版",
		Description: "能力均衡，适合日常使用", Credits: "0.79x",
		MaxInputTokens: 200000, MaxOutputTokens: 48000, ContextLength: 200000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"glm-5.0-turbo": {
		ID: "glm-5.0-turbo", Name: "GLM-5.0-Turbo", Realm: "cn", RealmLabel: "国内版",
		Description: "面向 Agent 场景进行了深度优化", Credits: "x0.95",
		MaxInputTokens: 200000, MaxOutputTokens: 48000, ContextLength: 200000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"glm-4.6v": {
		ID: "glm-4.6v", Name: "GLM-4.6V", Realm: "cn", RealmLabel: "国内版",
		Description: "GLM-4.6V 多模态模型，支持图片输入", Credits: "x0.11",
		MaxInputTokens: 128000, MaxOutputTokens: 32000, ContextLength: 128000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"kimi-k3-1": {
		ID: "kimi-k3-1", Name: "Kimi-K3", Realm: "cn", RealmLabel: "国内版",
		Description: "擅长处理复杂的长程自主任务，前端开发能力突出，同时在知识工作与科研推理上表现出色。", Credits: "x1.62",
		MaxInputTokens: 1000000, MaxOutputTokens: 32000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high", "xhigh"},
	},
	"kimi-k2.8-preview": {
		ID: "kimi-k2.8-preview", Name: "Kimi-K2.8-Preview", Realm: "cn", RealmLabel: "国内版",
		Description: "擅长处理复杂的长程自主任务，前端开发能力突出，同时在知识工作与科研推理上表现出色。", Credits: "x0.77",
		MaxInputTokens: 1000000, MaxOutputTokens: 32000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high", "max"},
	},
	"kimi-k2.7": {
		ID: "kimi-k2.7", Name: "Kimi-K2.7-Code", Realm: "cn", RealmLabel: "国内版",
		Description: "多模态模型，适合日常任务", Credits: "x0.57",
		MaxInputTokens: 256000, MaxOutputTokens: 32000, ContextLength: 256000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"kimi-k2.6": {
		ID: "kimi-k2.6", Name: "Kimi-K2.6", Realm: "cn", RealmLabel: "国内版",
		Description: "多模态模型，适合日常任务", Credits: "x0.52",
		MaxInputTokens: 256000, MaxOutputTokens: 32000, ContextLength: 256000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"kimi-k2.5": {
		ID: "kimi-k2.5", Name: "Kimi-K2.5", Realm: "cn", RealmLabel: "国内版",
		Description: "多模态模型，适合日常任务", Credits: "x0.45",
		MaxInputTokens: 256000, MaxOutputTokens: 32000, ContextLength: 256000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"kimi-k2-thinking": {
		ID: "kimi-k2-thinking", Name: "Kimi-K2-Thinking", Realm: "cn", RealmLabel: "国内版",
		Description: "适合复杂编码任务", Credits: "x0.54",
		MaxInputTokens: 256000, MaxOutputTokens: 32000, ContextLength: 256000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"medium"},
	},
	"minimax-m3": {
		ID: "minimax-m3", Name: "MiniMax-M3", Realm: "cn", RealmLabel: "国内版",
		Description: "原生多模态，擅长代码、智能体任务", Credits: "0.25x",
		MaxInputTokens: 512000, MaxOutputTokens: 128000, ContextLength: 512000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"medium"},
	},
	"minimax-m2.7": {
		ID: "minimax-m2.7", Name: "MiniMax-M2.7", Realm: "cn", RealmLabel: "国内版",
		Description: "能力均衡，适合日常使用", Credits: "x0.26",
		MaxInputTokens: 200000, MaxOutputTokens: 48000, ContextLength: 200000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"glm-4.6": {
		ID: "glm-4.6", Name: "GLM-4.6", Realm: "cn", RealmLabel: "国内版",
		Description: "具有强大推理能力的先进语言模型", Credits: "x0.23",
		MaxInputTokens: 168000, MaxOutputTokens: 32000, ContextLength: 168000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"deepseek-v4-flash": {
		ID: "deepseek-v4-flash", Name: "Deepseek-V4-Flash", Realm: "cn", RealmLabel: "国内版",
		Description: "DeepSeek 旗舰模型，支持 1M 上下文窗口", Credits: "x0.17",
		MaxInputTokens: 1000000, MaxOutputTokens: 50000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh"},
	},
	"deepseek-v4.1-flash": {
		ID: "deepseek-v4.1-flash", Name: "Deepseek-V4.1-Flash", Realm: "cn", RealmLabel: "国内版",
		Description: "DeepSeek 旗舰模型，支持 1M 上下文窗口，原生多模态", Credits: "0.03x (独家优惠)",
		MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"deepseek-v4-pro": {
		ID: "deepseek-v4-pro", Name: "Deepseek-V4-Pro", Realm: "cn", RealmLabel: "国内版",
		Description: "DeepSeek 旗舰模型，支持 1M 上下文窗口", Credits: "x0.51",
		MaxInputTokens: 1000000, MaxOutputTokens: 50000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"deepseek-v3-1": {
		ID: "deepseek-v3-1", Name: "DeepSeek-V3.1", Realm: "cn", RealmLabel: "国内版",
		Description: "DeepSeek 的旗舰模型，适合规划、调试、编码等任务", Credits: "x0.52",
		MaxInputTokens: 96000, MaxOutputTokens: 32000, ContextLength: 128000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: false,
		ReasoningEfforts: nil,
	},
	"deepseek-r1-0528": {
		ID: "deepseek-r1-0528", Name: "deepseek-r1", Realm: "cn", RealmLabel: "国内版",
		Description: "DeepSeek 的开源推理模型，专为逻辑与数学优化", Credits: "",
		MaxInputTokens: 96000, MaxOutputTokens: 8192, ContextLength: 96000,
		SupportsImages: false, SupportsTools: true, SupportsReasoning: false,
		ReasoningEfforts: nil,
	},
	"deepseek-v3-0324": {
		ID: "deepseek-v3-0324", Name: "deepseek-v3", Realm: "cn", RealmLabel: "国内版",
		Description: "DeepSeek 的旗舰模型，适合规划、调试、编码等任务", Credits: "",
		MaxInputTokens: 96000, MaxOutputTokens: 8192, ContextLength: 96000,
		SupportsImages: false, SupportsTools: true, SupportsReasoning: false,
		ReasoningEfforts: nil,
	},
	"hunyuan-2.0-instruct": {
		ID: "hunyuan-2.0-instruct", Name: "Hunyuan-2.0-Instruct", Realm: "cn", RealmLabel: "国内版",
		Description: "国内版专属模型", Credits: "",
		MaxInputTokens: 128000, MaxOutputTokens: 16000, ContextLength: 128000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: nil,
	},
	"hunyuan-chat": {
		ID: "hunyuan-chat", Name: "Hunyuan-Turbos", Realm: "cn", RealmLabel: "国内版",
		Description: "腾讯自研的轻量、快速的通用模型", Credits: "",
		MaxInputTokens: 128000, MaxOutputTokens: 8192, ContextLength: 128000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: false,
		ReasoningEfforts: nil,
	},
}

// 国际版真实模型清单（共 15 款）
var globalCatalog = map[string]ModelCatalogEntry{
	"deepseek-v4.1-flash": {
		ID: "deepseek-v4.1-flash", Name: "Deepseek-V4.1-Flash", Realm: "global", RealmLabel: "国际版",
		Description: "DeepSeek 旗舰模型，支持 1M 上下文窗口，原生多模态", Credits: "x0.00",
		MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"gpt-6-astra": {
		ID: "gpt-6-astra", Name: "GPT-6-Astra", Realm: "global", RealmLabel: "国际版",
		Description: "OpenAI 旗舰模型，擅长复杂推理与长程任务", Credits: "x6.67",
		MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"hy4-preview": {
		ID: "hy4-preview", Name: "Hy4 preview", Realm: "global", RealmLabel: "国际版",
		Description: "混元思考模型，具有增强的推理能力", Credits: "x0.29",
		MaxInputTokens: 1000000, MaxOutputTokens: 64000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"high"},
	},
	"hy3": {
		ID: "hy3", Name: "Hy3", Realm: "global", RealmLabel: "国际版",
		Description: "混元思考模型，具有增强的推理能力", Credits: "x0.00",
		MaxInputTokens: 192000, MaxOutputTokens: 64000, ContextLength: 192000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high"},
	},
	"gpt-5.6-sol": {
		ID: "gpt-5.6-sol", Name: "GPT-5.6-Sol", Realm: "global", RealmLabel: "国际版",
		Description: "OpenAI 旗舰模型，擅长复杂推理与长程任务", Credits: "x3.47",
		MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"gpt-5.6-terra": {
		ID: "gpt-5.6-terra", Name: "GPT-5.6-Terra", Realm: "global", RealmLabel: "国际版",
		Description: "OpenAI 均衡模型，兼顾能力、速度与成本", Credits: "x1.39",
		MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"gpt-5.6-luna": {
		ID: "gpt-5.6-luna", Name: "GPT-5.6-Luna", Realm: "global", RealmLabel: "国际版",
		Description: "OpenAI 轻量模型，响应快速，适合日常任务", Credits: "x0.14",
		MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"gpt-5.5": {
		ID: "gpt-5.5", Name: "GPT-5.5", Realm: "global", RealmLabel: "国际版",
		Description: "OpenAI 旗舰编码模型，擅长长程任务", Credits: "x3.31",
		MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh"},
	},
	"gpt-5.4": {
		ID: "gpt-5.4", Name: "GPT-5.4", Realm: "global", RealmLabel: "国际版",
		Description: "OpenAI 旗舰编码模型，擅长长程任务", Credits: "x1.65",
		MaxInputTokens: 272000, MaxOutputTokens: 72000, ContextLength: 272000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh"},
	},
	"gpt-5.3-codex": {
		ID: "gpt-5.3-codex", Name: "GPT-5.3-Codex", Realm: "global", RealmLabel: "国际版",
		Description: "OpenAI 代码专用模型，非常擅长处理复杂的编码任务", Credits: "x1.25",
		MaxInputTokens: 272000, MaxOutputTokens: 72000, ContextLength: 272000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"medium"},
	},
	"gemini-3.5-flash": {
		ID: "gemini-3.5-flash", Name: "Gemini-3.5-Flash", Realm: "global", RealmLabel: "国际版",
		Description: "能力均衡，适合日常使用", Credits: "x0.99",
		MaxInputTokens: 1000000, MaxOutputTokens: 65536, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "max"},
	},
	"glm-5.3": {
		ID: "glm-5.3", Name: "GLM-5.3", Realm: "global", RealmLabel: "国际版",
		Description: "能力均衡，适合日常使用", Credits: "x0.79",
		MaxInputTokens: 1000000, MaxOutputTokens: 48000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high", "max"},
	},
	"glm-5.2": {
		ID: "glm-5.2", Name: "GLM-5.2", Realm: "global", RealmLabel: "国际版",
		Description: "1M 上下文，擅长长程任务", Credits: "x0.79",
		MaxInputTokens: 1000000, MaxOutputTokens: 48000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"high", "xhigh"},
	},
	"kimi-k3": {
		ID: "kimi-k3", Name: "Kimi-K3", Realm: "global", RealmLabel: "国际版",
		Description: "擅长处理复杂的长程自主任务，前端开发能力突出，同时在知识工作与科研推理上表现出色。", Credits: "x1.62",
		MaxInputTokens: 1000000, MaxOutputTokens: 32000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high", "xhigh"},
	},
	"kimi-k2.6": {
		ID: "kimi-k2.6", Name: "Kimi-K2.6", Realm: "global", RealmLabel: "国际版",
		Description: "多模态模型，适合日常任务", Credits: "x0.52",
		MaxInputTokens: 256000, MaxOutputTokens: 32000, ContextLength: 256000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"medium"},
	},
}

// ModelRealmDetermined 根据模型名称直接精确判断必须由哪个域承接
func ModelRealmDetermined(modelID string) string {
	raw := strings.TrimSpace(modelID)
	// 支持前缀显式指定：cn/xxx 强走国内，global/xxx 强走国际
	if strings.HasPrefix(raw, "cn/") {
		return "cn"
	}
	if strings.HasPrefix(raw, "global/") {
		return "global"
	}

	_, inCN := cnCatalog[raw]
	_, inGlobal := globalCatalog[raw]

	if inGlobal && !inCN {
		return "global" // 仅国际版拥有
	}
	if inCN && !inGlobal {
		return "cn" // 仅国内版拥有
	}
	if inGlobal && inCN {
		return "any" // 两边都有的同名模型，按配置轮转
	}

// 兜底推断规则
	m := strings.ToLower(raw)
	if strings.HasPrefix(m, "gpt-") || strings.HasPrefix(m, "o1") || strings.HasPrefix(m, "o3") || strings.HasPrefix(m, "gemini-") {
		return "global"
	}
	return "cn"
}

// AllDistinctModelsList 返回独立区分的国内外全量模型列表：
// 1. 国内版全部附加 cn/ 前缀（如 cn/deepseek-v4.1-flash, cn/glm-5.3）
// 2. 国际版全部附加 global/ 前缀（如 global/deepseek-v4.1-flash, global/gpt-6-astra）
// 即使同名也完全独立宣告，不带倍率显示，描述中不含任何强制字样。
func AllDistinctModelsList() []map[string]any {
	out := make([]map[string]any, 0, len(cnCatalog)+len(globalCatalog))

	// 1. 国内版模型库 (全部带 cn/ 前缀)
	for _, meta := range cnCatalog {
		supportsReasoning := len(meta.ReasoningEfforts) > 0
		out = append(out, map[string]any{
			"id":                 "cn/" + meta.ID,
			"name":               "[国内] " + meta.Name,
			"object":             "model",
			"created":            1753600000,
			"owned_by":           "workbuddy-cn",
			"realm":              "cn",
			"realm_label":        "国内版",
			"context_length":     meta.ContextLength,
			"max_input_tokens":   meta.MaxInputTokens,
			"max_output_tokens":  meta.MaxOutputTokens,
			"max_tokens":         meta.MaxOutputTokens,
			"description":        meta.Description,
			"supports_vision":    meta.SupportsImages,
			"supports_images":    meta.SupportsImages,
			"supports_tools":     meta.SupportsTools,
			"supports_reasoning": supportsReasoning,
			"reasoning_efforts":  meta.ReasoningEfforts,
			"capabilities": map[string]bool{
				"vision":    meta.SupportsImages,
				"tools":     meta.SupportsTools,
				"reasoning": supportsReasoning,
			},
		})
	}

	// 2. 国际版模型库 (全部带 global/ 前缀)
	for _, meta := range globalCatalog {
		supportsReasoning := len(meta.ReasoningEfforts) > 0
		out = append(out, map[string]any{
			"id":                 "global/" + meta.ID,
			"name":               "[国际] " + meta.Name,
			"object":             "model",
			"created":            1753600000,
			"owned_by":           "workbuddy-global",
			"realm":              "global",
			"realm_label":        "国际版",
			"context_length":     meta.ContextLength,
			"max_input_tokens":   meta.MaxInputTokens,
			"max_output_tokens":  meta.MaxOutputTokens,
			"max_tokens":         meta.MaxOutputTokens,
			"description":        meta.Description,
			"supports_vision":    meta.SupportsImages,
			"supports_images":    meta.SupportsImages,
			"supports_tools":     meta.SupportsTools,
			"supports_reasoning": supportsReasoning,
			"reasoning_efforts":  meta.ReasoningEfforts,
			"capabilities": map[string]bool{
				"vision":    meta.SupportsImages,
				"tools":     meta.SupportsTools,
				"reasoning": supportsReasoning,
			},
		})
	}

	return out
}
