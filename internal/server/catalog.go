package server

// ModelCatalogEntry 完整的模型能力定义
type ModelCatalogEntry struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
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

// 国际版完整模型能力词典 (来自 workbuddy2api-intl 官方快照)
var modelCatalog = map[string]ModelCatalogEntry{
	"deepseek-v4.1-flash": {
		ID: "deepseek-v4.1-flash", Name: "DeepSeek-V4.1-Flash",
		Description: "DeepSeek 旗舰模型，支持 1M 上下文窗口，原生多模态",
		Credits: "x0.00", MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"high"},
	},
	"gpt-6-astra": {
		ID: "gpt-6-astra", Name: "GPT-6-Astra",
		Description: "OpenAI 旗舰模型，擅长复杂推理与长程任务",
		Credits: "x6.67", MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"gpt-5.6-sol": {
		ID: "gpt-5.6-sol", Name: "GPT-5.6-Sol",
		Description: "OpenAI 旗舰模型，擅长复杂推理与长程任务",
		Credits: "x3.47", MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"gpt-5.6-terra": {
		ID: "gpt-5.6-terra", Name: "GPT-5.6-Terra",
		Description: "OpenAI 均衡模型，兼顾能力、速度与成本",
		Credits: "x1.39", MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"gpt-5.6-luna": {
		ID: "gpt-5.6-luna", Name: "GPT-5.6-Luna",
		Description: "OpenAI 轻量模型，响应快速，适合日常任务",
		Credits: "x0.14", MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh", "max"},
	},
	"gpt-5.5": {
		ID: "gpt-5.5", Name: "GPT-5.5",
		Description: "OpenAI 旗舰编码模型，擅长长程任务",
		Credits: "x3.31", MaxInputTokens: 1000000, MaxOutputTokens: 128000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh"},
	},
	"gpt-5.4": {
		ID: "gpt-5.4", Name: "GPT-5.4",
		Description: "OpenAI 旗舰编码模型，擅长长程任务",
		Credits: "x1.65", MaxInputTokens: 272000, MaxOutputTokens: 72000, ContextLength: 272000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "medium", "high", "xhigh"},
	},
	"gpt-5.3-codex": {
		ID: "gpt-5.3-codex", Name: "GPT-5.3-Codex",
		Description: "OpenAI 代码专用模型，非常擅长处理复杂的编码任务",
		Credits: "x1.25", MaxInputTokens: 272000, MaxOutputTokens: 72000, ContextLength: 272000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"medium"},
	},
	"hy4-preview-f": {
		ID: "hy4-preview-f", Name: "Hy4 Preview Free",
		Description: "混元思考模型，具有增强的推理能力",
		Credits: "x0.00", MaxInputTokens: 1000000, MaxOutputTokens: 64000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"high"},
	},
	"hy4-preview": {
		ID: "hy4-preview", Name: "Hy4 Preview",
		Description: "混元思考模型，具有增强的推理能力",
		Credits: "x0.29", MaxInputTokens: 1000000, MaxOutputTokens: 64000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"high"},
	},
	"hy3": {
		ID: "hy3", Name: "Hy3",
		Description: "混元思考模型，具有增强的推理能力",
		Credits: "x0.00", MaxInputTokens: 192000, MaxOutputTokens: 64000, ContextLength: 192000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high"},
	},
	"gemini-3.5-flash": {
		ID: "gemini-3.5-flash", Name: "Gemini-3.5-Flash",
		Description: "Google 全能均衡模型，适合日常及长文本任务",
		Credits: "x0.99", MaxInputTokens: 1000000, MaxOutputTokens: 65536, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"medium"},
	},
	"glm-5.3": {
		ID: "glm-5.3", Name: "GLM-5.3",
		Description: "1M 上下文，能力均衡，适合日常使用",
		Credits: "x0.79", MaxInputTokens: 1000000, MaxOutputTokens: 48000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"low", "high", "max"},
	},
	"glm-5.2": {
		ID: "glm-5.2", Name: "GLM-5.2",
		Description: "1M 上下文，擅长长程任务",
		Credits: "x0.79", MaxInputTokens: 1000000, MaxOutputTokens: 48000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"high", "xhigh"},
	},
	"kimi-k3": {
		ID: "kimi-k3", Name: "Kimi-K3",
		Description: "擅长长程自主任务、前端开发与科研推理",
		Credits: "x1.62", MaxInputTokens: 1000000, MaxOutputTokens: 32000, ContextLength: 1000000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"medium"},
	},
	"kimi-k2.6": {
		ID: "kimi-k2.6", Name: "Kimi-K2.6",
		Description: "多模态模型，适合日常任务",
		Credits: "x0.52", MaxInputTokens: 256000, MaxOutputTokens: 32000, ContextLength: 256000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
		ReasoningEfforts: []string{"medium"},
	},
	"default-model": {
		ID: "default-model", Name: "Auto",
		Description: "优秀编码模型，适合日常开发",
		Credits: "x1.00", MaxInputTokens: 200000, MaxOutputTokens: 24000, ContextLength: 200000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
	},
	"primary-model": {
		ID: "primary-model", Name: "Primary",
		Description: "高质量输出，胜任复杂挑战",
		Credits: "x3.31", MaxInputTokens: 272000, MaxOutputTokens: 72000, ContextLength: 272000,
		SupportsImages: true, SupportsTools: true, SupportsReasoning: true,
	},
}

// enrichModelList 将模型能力元数据挂载到输出列表上
func enrichModelList(list []map[string]any) []map[string]any {
	out := make([]map[string]any, 0, len(list))

	for _, item := range list {
		id, _ := item["id"].(string)
		if id == "" {
			continue
		}

		if meta, ok := modelCatalog[id]; ok {
			item["context_length"] = meta.ContextLength
			item["max_input_tokens"] = meta.MaxInputTokens
			item["max_output_tokens"] = meta.MaxOutputTokens
			item["max_tokens"] = meta.MaxOutputTokens
			item["description"] = meta.Description
			item["credits"] = meta.Credits
			item["supports_vision"] = meta.SupportsImages
			item["supports_images"] = meta.SupportsImages
			item["supports_tools"] = meta.SupportsTools
			item["supports_reasoning"] = meta.SupportsReasoning
			if len(meta.ReasoningEfforts) > 0 {
				item["reasoning_efforts"] = meta.ReasoningEfforts
			}
			item["capabilities"] = map[string]bool{
				"vision":    meta.SupportsImages,
				"tools":     meta.SupportsTools,
				"reasoning": meta.SupportsReasoning,
			}
		} else {
			if item["context_length"] == nil || item["context_length"] == 0 {
				item["context_length"] = 131072
			}
			item["supports_vision"] = true
			item["supports_tools"] = true
			item["capabilities"] = map[string]bool{
				"vision": true,
				"tools":  true,
			}
		}
		out = append(out, item)
	}

	return out
}
