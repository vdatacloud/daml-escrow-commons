package offboarding

import (
	"testing"
	"time"
)

func TestStage(t *testing.T) {
	now := time.Now()
	c := Case{PartySetID: "ps-1"}
	for _, step := range []struct {
		set  func()
		want string
	}{
		{func() {}, "opened"},
		{func() { c.ExportedAt = &now }, "exported"},
		{func() { c.ReceiptRecordedAt = &now }, "receipted"},
		{func() { c.DeletedAt = &now }, "deleted"},
		{func() { c.IdentityErasedAt = &now }, "closed"},
	} {
		step.set()
		if got := c.Stage(); got != step.want {
			t.Errorf("got %s want %s", got, step.want)
		}
	}
	if got := (Case{Institution: "bank", DeletedAt: &now}).Stage(); got != "closed" {
		t.Errorf("institution deleted: %s", got)
	}
}
