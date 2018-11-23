/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package config

// config.go
const (
	SeelePath     = "github.com/seeleteam/go-seele"
	CoverFileName = "seele_coverage_detail"
	CoverPackage  = "common\t,core\t,trie\t,p2p\t,seele\t"

	Subject    = "Daily Test Report"
	Sender     = "send@email.com"
	Password   = "password"
	SenderName = "reporter"
	Receivers  = "receiver@email.com"
	CC         = "CC@email.com"
	Host       = "smtp.exmail.qq.com:25"

	StartHour = 04
	StartMin  = 00
	StartSec  = 00

	BenchTopN         = "15"
	BenchReportFormat = "pdf"
)
