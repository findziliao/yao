package pipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/gou/process"
	"github.com/yaoapp/kun/any"
	"github.com/yaoapp/yao/test"
)

func TestProcessPipes(t *testing.T) {
	prepare(t)
	defer test.Clean()

	p, err := process.Of("pipes.cli.translator", map[string]interface{}{"placeholder": "translate\nhello world"})
	if err != nil {
		t.Fatal(err)
	}

	output, err := p.Exec()
	res := any.Of(output).Map().MapStrAny.Dot()
	assert.True(t, res.Has("global"))
	assert.True(t, res.Has("input"))
	assert.True(t, res.Has("output"))
	assert.True(t, res.Has("sid"))
	assert.True(t, res.Has("switch"))
	assert.Equal(t, "translate\nhello world", res.Get("input[0].placeholder"))
	assert.Len(t, res.Get("switch"), 2)
}

func TestProcessRun(t *testing.T) {
	prepare(t)
	defer test.Clean()

	p, err := process.Of("pipe.Run", "cli.translator", map[string]interface{}{"placeholder": "translate\nhello world"})
	if err != nil {
		t.Fatal(err)
	}

	output, err := p.Exec()
	if err != nil {
		t.Fatal(err)
	}

	res := any.Of(output).Map().MapStrAny.Dot()
	assert.True(t, res.Has("global"))
	assert.True(t, res.Has("input"))
	assert.True(t, res.Has("output"))
	assert.True(t, res.Has("sid"))
	assert.True(t, res.Has("switch"))
	assert.Equal(t, "translate\nhello world", res.Get("input[0].placeholder"))
	assert.Len(t, res.Get("switch"), 2)
}

func TestProcessResume(t *testing.T) {
	prepare(t)
	defer test.Clean()

	p, err := process.Of("pipe.Run", "web.translator", "hello web world")
	if err != nil {
		t.Fatal(err)
	}

	web, err := p.Exec()
	resume := web.(ResumeContext)

	p, err = process.Of("pipe.Resume", resume.ID, "translate", "hello web world")
	output, err := p.Exec()
	if err != nil {
		t.Fatal(err)
	}

	res := any.Of(output).Map().MapStrAny.Dot()
	assert.True(t, res.Has("global"))
	assert.True(t, res.Has("input"))
	assert.True(t, res.Has("output"))
	assert.True(t, res.Has("sid"))
	assert.True(t, res.Has("switch"))
	assert.Equal(t, "hello web world", res.Get("input[0]"))
	assert.Len(t, res.Get("switch"), 2)
}

func TestProcessClose(t *testing.T) {
	prepare(t)
	defer test.Clean()

	p, err := process.Of("pipe.Run", "web.translator", "hello web world")
	if err != nil {
		t.Fatal(err)
	}

	web, err := p.Exec()
	resume := web.(ResumeContext)

	p, err = process.Of("pipe.Close", resume.ID)
	p.Exec()

	p, err = process.Of("pipe.Resume", resume.ID, "translate", "hello web world")
	_, err = p.Exec()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not found")
}
