package bench

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func Test_Run(t *testing.T) {
	benchPath := "github.com/seeleteam/go-seele/core"
	all := Run(benchPath)
	assert.Equal(t, all != "", true)
	fmt.Println("all:\n", all)
}
