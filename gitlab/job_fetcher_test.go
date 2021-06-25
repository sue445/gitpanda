package gitlab

import (
	"github.com/stretchr/testify/assert"
	"github.com/sue445/gitpanda/testutil"
	"testing"
)

func Test_normalizeJobTrace(t *testing.T) {
	raw := testutil.ReadTestData("testdata/job_trace.txt")
	got := normalizeJobTrace(raw)

	expected := testutil.ReadTestData("testdata/job_trace_plain.txt")
	assert.Equal(t, expected, got)
}
