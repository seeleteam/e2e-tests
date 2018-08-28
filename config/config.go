package config

// config.go
const (
	SeelePath     = "github.com/seeleteam/go-seele"
	CoverFileName = "seele_coverage_detail"
	CoverPackage  = "common\t,core\t,trie\t,p2p\t,seele\t"

	Subject    = "Daily Test Report"
	Sender     = "wangfeifan@zsbatech.com"
	Password   = "Wff19940326..."
	SenderName = "Seele-e2e"
	// Receivers  = "wangfeifan@zsbatech.com"
	// Receivers = "rdc@zsbatech.com"
	Receivers = "dev@seelenet.com"
	Host      = "smtp.exmail.qq.com:25"

	StartHour = 04
	StartMin  = 00
	StartSec  = 00

	BenchTopN         = "15"
	BenchReportFormat = "pdf"
)
