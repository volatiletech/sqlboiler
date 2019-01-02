package boil

import (
	"context"
	"testing"
)

func TestSkipHooks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	if HooksAreSkipped(ctx) {
		t.Error("they should not be skipped")
	}

	ctx = SkipHooks(ctx)

	if !HooksAreSkipped(ctx) {
		t.Error("they should be skipped")
	}
}

func TestSkipTimestamps(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	if TimestampsAreSkipped(ctx) {
		t.Error("they should not be skipped")
	}

	ctx = SkipTimestamps(ctx)

	if !TimestampsAreSkipped(ctx) {
		t.Error("they should be skipped")
	}
}
