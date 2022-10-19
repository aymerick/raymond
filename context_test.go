package raymond

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHandlebarsContext(t *testing.T) {
	suite.Run(t, new(HandlebarsContextSuite))
}

type HandlebarsContextSuite struct {
	suite.Suite
	c *handlebarsContext
}

func (s *HandlebarsContextSuite) SetupTest() {
	s.c = newHandlebarsContext()
}

func (s *HandlebarsContextSuite) TestHandlebarsContextAddMemberContext() {
	assert.Equal(s.T(), 0, len(s.c.GetCurrentContext()), "Len expected to be zero.")
	s.c.AddMemberContext("foo", "bar")
	assert.Equal(s.T(), 1, len(s.c.GetCurrentContext()), "Len expected to be one.")
	s.c.AddMemberContext("baz", "bar")
	assert.Equal(s.T(), 2, len(s.c.GetCurrentContext()), "Len expected to be two.")
	s.c.AddMemberContext("bean", "bar")
	assert.Equal(s.T(), 3, len(s.c.GetCurrentContext()), "Len expected to be three.")
	assert.Equal(s.T(), "foo.baz.bean", s.c.GetCurrentContextString(), "Should be all three scopes.")
	assert.Equal(s.T(), 2, len(s.c.GetParentContext(1)), "Len expected to be two.")
	assert.Equal(s.T(), "foo.baz", s.c.GetParentContextString(1), "Should be two scopes.")
	assert.Equal(s.T(), "", s.c.GetParentContextString(3), "Should  be empty string.")
	assert.Equal(s.T(), "foo.baz.bean", s.c.GetParentContextString(4), "Parent context exceeded use default ancestor.")
	assert.Equal(s.T(), 3, len(s.c.GetCurrentContext()), "Len expected to be three.")
	s.c.MoveUpContext()
	assert.Equal(s.T(), 2, len(s.c.GetCurrentContext()), "Len expected to be two.")
	assert.Equal(s.T(), "foo.baz", s.c.GetCurrentContextString(), "Should be two scopes.")
}

func (s *HandlebarsContextSuite) TestHandlebarsContextMappedContextAllTheSameMapping() {
	assert.Equal(s.T(), 0, len(s.c.GetCurrentContext()), "Len expected to be zero.")
	s.c.AddMemberContext("foo", "bar")
	s.c.AddMemberContext("baz", "bar")
	s.c.AddMemberContext("bean", "bar")
	assert.Equal(s.T(), "foo.blah.baz.bing.bean.bong", s.c.GetMappedContextString([]string{"bar", "blah", "bar", "bing", "bar", "bong"}, 0), "Should be all three scopes.")
}

func (s *HandlebarsContextSuite) TestHandlebarsContextMappedContextLongNamesSameMapping() {
	assert.Equal(s.T(), 0, len(s.c.GetCurrentContext()), "Len expected to be zero.")
	s.c.AddMemberContext("foo.foo.foo", "bar")
	s.c.AddMemberContext("baz.baz.baz", "bar")
	s.c.AddMemberContext("bean.bean.bean", "bar")
	assert.Equal(s.T(), "foo.foo.foo.baz.baz.baz.bean.bean.bean", s.c.GetMappedContextString([]string{"bar", "bar", "bar"}, 0), "Should be all three scopes.")
}

func (s *HandlebarsContextSuite) TestHandlebarsContextMappedContextLongNamesSameMappingNoMapping() {
	assert.Equal(s.T(), 0, len(s.c.GetCurrentContext()), "Len expected to be zero.")
	s.c.AddMemberContext("foo.foo.foo", "bar")
	s.c.AddMemberContext("baz.baz.baz", "")
	s.c.AddMemberContext("bean.bean.bean", "bar")
	assert.Equal(s.T(), "foo.foo.foo.baz.baz.baz.bleep.bean.bean.bean.bop", s.c.GetMappedContextString([]string{"bar", "bleep", "bar", "bop"}, 0), "Should be all three scopes.")
}
