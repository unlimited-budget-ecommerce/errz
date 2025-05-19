package errorz_test

import (
	"testing"

	"github.com/unlimited-budget-ecommerce/errorz"
	testutils "github.com/unlimited-budget-ecommerce/errorz/utils"
)

func TestGroupErrors_FileNotFound(t *testing.T) {
	_, err := errorz.GroupErrorsByDomain("nonexistent.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestGroupErrors_InvalidJSON(t *testing.T) {
	path := testutils.WriteTempFile(t, "invalid.json", testutils.InvalidJSON)

	_, err := errorz.GroupErrorsByDomain(path)
	if err == nil {
		t.Error("expected parse error, got nil")
	}
}

func TestGroupErrors_EmptyJSON(t *testing.T) {
	content := `{}`
	path := testutils.WriteTempFile(t, "empty.json", content)

	_, err := errorz.GroupErrorsByDomain(path)
	if err == nil {
		t.Error("expected error for empty map, got nil")
	}
}

func TestGroupErrors_ValidJSON(t *testing.T) {
	path := testutils.WriteTempFile(t, "valid.json", testutils.ManyValidJSON)

	groups, err := errorz.GroupErrorsByDomain(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(groups) != 3 {
		t.Errorf("expected 3 domains, got %d", len(groups))
	}

	if len(groups["payment"]) != 2 {
		t.Errorf("expected 2 errors in payment, got %d", len(groups["payment"]))
	}

	if len(groups["auth"]) != 1 {
		t.Errorf("expected 1 error in auth, got %d", len(groups["auth"]))
	}

	if len(groups["order"]) != 1 {
		t.Errorf("expected 1 error in order, got %d", len(groups["order"]))
	}

	// Check only some important fields
	if msg := groups["payment"]["PM0002"].Msg; msg != "timeout" {
		t.Errorf("unexpected message for PM0002: %s", msg)
	}

	if retry := groups["payment"]["PM0002"].IsRetryable; !retry {
		t.Errorf("expected PM0002 to be retryable")
	}
}
