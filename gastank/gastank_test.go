package gastank

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWireNames(t *testing.T) {
	b, _ := json.Marshal(Response{Networks: []Status{{Coin: "arbeth", Minimum: "0.05", Low: true}}})
	for _, k := range []string{`"networks"`, `"coin":"arbeth"`, `"minimum":"0.05"`, `"low":true`, `"checkedAt"`} {
		if !strings.Contains(string(b), k) {
			t.Errorf("missing %s in %s", k, b)
		}
	}
}
