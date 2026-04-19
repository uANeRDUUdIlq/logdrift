package contextlines

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew_NegativeValues_ClampedToZero(t *testing.T) {
	b := New(-1, -2)
	assert.Equal(t, 0, b.before)
	assert.Equal(t, 0, b.after)
}

func TestPush_NoContext_OnlyMatchEmitted(t *testing.T) {
	b := New(0, 0)
	out := b.Push("line1", "svc", false)
	assert.Empty(t, out)
	out = b.Push("match", "svc", true)
	assert.Equal(t, []string{"match"}, out)
}

func TestPush_BeforeContext_EmittedOnMatch(t *testing.T) {
	b := New(2, 0)
	b.Push("a", "svc", false)
	b.Push("b", "svc", false)
	b.Push("c", "svc", false)
	out := b.Push("match", "svc", true)
	// ring holds last 2 before lines: b, c
	assert.Contains(t, out, "b")
	assert.Contains(t, out, "c")
	assert.Contains(t, out, "match")
	assert.NotContains(t, out, "a")
}

func TestPush_AfterContext_EmittedAfterMatch(t *testing.T) {
	b := New(0, 2)
	out := b.Push("match", "svc", true)
	assert.Equal(t, []string{"match"}, out)

	out = b.Push("post1", "svc", false)
	assert.Contains(t, out, "post1")

	out = b.Push("post2", "svc", false)
	assert.Contains(t, out, "post2")

	out = b.Push("post3", "svc", false)
	assert.Empty(t, out)
}

func TestPush_NoMatch_NonContextLine_NotEmitted(t *testing.T) {
	b := New(1, 1)
	out := b.Push("irrelevant", "svc", false)
	assert.Empty(t, out)
}

func TestReset_ClearsState(t *testing.T) {
	b := New(2, 2)
	b.Push("a", "svc", false)
	b.Push("match", "svc", true)
	b.Reset()
	out := b.Push("after-reset", "svc", false)
	assert.Empty(t, out)
}
