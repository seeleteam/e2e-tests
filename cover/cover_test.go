package cover

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/seeleteam/e2e-tests/store"
)

func Test_Run(t *testing.T) {
	coverPath := "github.com/seeleteam/go-seele/common/..."
	all, specified := Run(coverPath)
	assert.Equal(t, all != "", true)
	assert.Equal(t, specified != nil, true)
}

func Test_PrintSpecifiedPkg(t *testing.T) {
	specified, date := make(map[string]string), "20180816"
	specified["common"] = "ok      github.com/seeleteam/go-seele/common    0.105s  coverage: 79.6% of statements"
	specified["core"] = "ok      github.com/seeleteam/go-seele/core      4.745s  coverage: 76.3% of statements"
	bytes, err := json.Marshal(specified)
	assert.Equal(t, err, nil)
	store.Save(date, "", "", bytes)

	result := PrintSpecifiedPkg(date, specified)
	assert.Equal(t, strings.Contains(result, "FAIL"), false)

	specified["p2p"] = "ok      github.com/seeleteam/go-seele/p2p       2.056s  coverage: 21.0% of statements"
	result1 := PrintSpecifiedPkg(date, specified)
	assert.Equal(t, strings.Contains(result1, "FAIL"), false)
}
