package logger

// Format 代表日誌輸出的格式類型。
// 不同的格式適用於不同的使用場景：
//   - JSON 格式：適合生產環境，便於日誌解析和分析
//   - Console 格式：適合開發環境，提供更友善的人類閱讀體驗
type Format string

// 日誌格式常數定義
const (
	// JSONFormat JSON 格式輸出
	// 輸出範例：{"level":"info","timestamp":"2024-01-01T12:00:00Z","message":"Hello","user":"alice"}
	// 特點：
	//   - 結構化資料，便於機器解析
	//   - 適合生產環境和日誌分析系統
	//   - 支援複雜的巢狀資料結構
	JSONFormat Format = "json"

	// ConsoleFormat 控制台友善格式輸出
	// 輸出範例：2024-01-01T12:00:00Z  INFO  Hello  user=alice
	// 特點：
	//   - 人類友善的閱讀格式
	//   - 適合開發環境和除錯
	//   - 支援顏色高亮（依 adapter 實作而定）
	ConsoleFormat Format = "console"
)

// ParseFormat 將日誌格式字串解析為對應的 Format 常數。
// 此函數提供容錯處理，無效輸入時返回預設格式。
//
// 支援的字串值：
//   - "json"    -> JSONFormat
//   - "console" -> ConsoleFormat
//
// 參數：
//   - format: 日誌格式的字串表示
//
// 返回值：
//   - Format: 對應的日誌格式常數，無效輸入時返回 ConsoleFormat
//
// 使用範例：
//
//	format := logger.ParseFormat("json")    // 返回 JSONFormat
//	format := logger.ParseFormat("invalid") // 返回 ConsoleFormat (預設值)
//
//	// 在配置中使用
//	cfg := console.Config{
//	    Format: logger.ParseFormat(os.Getenv("LOG_FORMAT")),
//	}
func ParseFormat(format string) Format {
	switch format {
	case "json":
		return JSONFormat
	case "console":
		return ConsoleFormat
	default:
		// 無效輸入時返回 Console 格式作為預設值
		// Console 格式更適合作為預設值，因為它對開發者更友善
		return ConsoleFormat
	}
}
