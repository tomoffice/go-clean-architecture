package logger

// Config 結構（供 factory 使用）
type Config struct {
	UseConsole bool
	UseGCP     bool
	UseSeq     bool

	Level  string // 共用
	Format string // 共用：json or console

	GCP GCPLoggerConfig // 👉 獨立出來，但只在 UseGCP = true 時才用
	Seq SeqLoggerConfig
}

// GCPLoggerConfig
// GCPLoggerConfig 用於設定 GCP Cloud Logging 的參數。
// 此結構會傳入 logger adapter 以建立對應的 GCP logger 實體。
//
// ResourceType（必填）：
//
//	ResourceType 指定 GCP Logging 中的 Monitored Resource 類型。
//	這將決定 log 在 GCP 中的分類方式與附加的 metadata 欄位。
//	常見值如下：
//	  - "k8s_container"：Kubernetes container（建議，最常用）
//	  - "k8s_pod"：Kubernetes Pod（Pod 層級記錄）
//	  - "gce_instance"：GCE instance（虛擬機）
//	  - "global"：通用資源類型（無關容器/機器，可用於本地測試）
//
// ResourceLabels（必填）：
//
//	用來補足 ResourceType 所需欄位，例如：
//	  ResourceType = "k8s_pod" 時，必須提供：
//	    - cluster_name
//	    - namespace_name
//	    - pod_name
//	    - location（zone 或 region）
//
//	範例：
//	  ResourceLabels: map[string]string{
//	      "cluster_name":    "my-cluster",
//	      "namespace_name":  "default",
//	      "pod_name":        "my-app-pod-xyz",
//	      "location":        "asia-east1-a",
//	  }
//
// 可選輔助欄位：
//
//	若你不想手動填 ResourceLabels，可透過下列欄位由程式補齊：
//	  - ClusterName   → 自動補 cluster_name
//	  - NamespaceName → 自動補 namespace_name
//	  - PodName       → 自動補 pod_name
//	  - Location      → 自動補 location
//
// MinSeverity（選填）：
//
//	最小輸出等級，例如 "INFO"、"WARNING"、"ERROR"
//	可以避免 log 雜訊。
type GCPLoggerConfig struct {
	ProjectID string // 必填，用來初始化 GCP Logging Client
	LogName   string // 必填，對應 GCP Log stream 名稱

	ResourceType string // 必填

	ResourceLabels map[string]string // 必填

	// Optional：協助補充 ResourceLabels
	ClusterName   string // 建議自動對應 cluster_name
	Location      string // 建議自動對應 location
	NamespaceName string // 建議自動對應 namespace_name
	PodName       string // 建議自動對應 pod_name

	// Optional：可用來過濾最低嚴重性，減少 log noise
	MinSeverity string // e.g., "INFO", "WARNING", "ERROR"
}
type SeqLoggerConfig struct {
	SeqURL    string
	SeqAPIKey string // optional
}
