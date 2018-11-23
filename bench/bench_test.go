
/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */
package bench

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func Test_Run(t *testing.T) {
	benchPath := "github.com/seeleteam/go-seele/core/types"
	all := Run(benchPath)
	assert.Equal(t, all != "", true)
	fmt.Println("all:\n", all)
}

func Test_OutputCompressionReport(t *testing.T) {
	output := OutputCompressionReport("reports.zip")
	assert.Equal(t, output, "")
}
