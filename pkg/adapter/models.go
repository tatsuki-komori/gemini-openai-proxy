package adapter

import (
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// Geminiモデルの定数定義
const (
	Gemini1Dot5Flash8b = "gemini-1.5-flash-8b"  // Gemini 1.5 Flash 8bモデル
	Gemini1Dot5Pro   = "gemini-1.5-pro-002"     // Gemini 1.5 Proモデル
	Gemini1Dot5Flash = "gemini-1.5-flash-002"    // Gemini 1.5 Flashモデル
	Gemini1Dot5ProV  = "gemini-1.0-pro-vision-latest" // Gemini Vision モデル - struct::ToGenaiModelで他のモデルに変換される
	TextEmbedding004 = "text-embedding-004"      // テキスト埋め込みモデル
)

// モデルマッピングの有効/無効を制御する環境変数
var USE_MODEL_MAPPING bool = os.Getenv("DISABLE_MODEL_MAPPING") != "1"

// モデルの所有者を返す関数
func GetOwner() string {
	if USE_MODEL_MAPPING {
		return "openai"
	} else {
		return "google"
	}
}

// OpenAIモデル名を適切なモデルに変換する関数
func GetModel(openAiModelName string) string {
	if USE_MODEL_MAPPING {
		return openAiModelName
	} else {
		return ConvertModel(openAiModelName)
	}
}

// GeminiモデルをOpenAIモデルにマッピングする関数
func GetMappedModel(geminiModelName string) string {
	if !USE_MODEL_MAPPING {
		return geminiModelName
	}
	switch {
	case geminiModelName == Gemini1Dot5ProV:
		return openai.GPT4VisionPreview
	case geminiModelName == Gemini1Dot5Pro:
		return openai.GPT4TurboPreview
	case geminiModelName == Gemini1Dot5Flash:
		return openai.GPT4
	case geminiModelName == TextEmbedding004:
		return string(openai.AdaEmbeddingV2)
	default:
		return openai.GPT3Dot5Turbo
	}
}

// OpenAIモデルをGeminiモデルに変換する関数
func ConvertModel(openAiModelName string) string {
	switch {
	case openAiModelName == openai.GPT4VisionPreview:
		return Gemini1Dot5ProV
	case openAiModelName == openai.GPT4TurboPreview || openAiModelName == openai.GPT4Turbo1106 || openAiModelName == openai.GPT4Turbo0125:
		return Gemini1Dot5Pro
	case strings.HasPrefix(openAiModelName, openai.GPT4):
		return Gemini1Dot5Flash
	case openAiModelName == string(openai.AdaEmbeddingV2):
		return TextEmbedding004
	default:
		return Gemini1Dot5Flash8b
	}
}

// ChatCompletionRequestのモデルをGeminiモデルに変換するメソッド
func (req *ChatCompletionRequest) ToGenaiModel() string {
	if USE_MODEL_MAPPING {
		return req.ParseModelWithMapping()
	} else {
		return req.ParseModelWithoutMapping()
	}
}

// マッピングなしでモデルを解析するメソッド
func (req *ChatCompletionRequest) ParseModelWithoutMapping() string {
	switch {
	case req.Model == Gemini1Dot5ProV:
		if os.Getenv("GPT_4_VISION_PREVIEW") == Gemini1Dot5Pro {
			return Gemini1Dot5Pro
		}

		return Gemini1Dot5Flash
	default:
		return req.Model
	}
}

// マッピングありでモデルを解析するメソッド
func (req *ChatCompletionRequest) ParseModelWithMapping() string {
	switch {
	case req.Model == openai.GPT4VisionPreview:
		if os.Getenv("GPT_4_VISION_PREVIEW") == Gemini1Dot5Pro {
			return Gemini1Dot5Pro
		}

		return Gemini1Dot5Flash
	default:
		return ConvertModel(req.Model)
	}
}

// EmbeddingRequestのモデルをGeminiモデルに変換するメソッド
func (req *EmbeddingRequest) ToGenaiModel() string {
	if USE_MODEL_MAPPING {
		return ConvertModel(req.Model)
	} else {
		return req.Model
	}
}
