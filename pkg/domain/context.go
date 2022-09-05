package domain

import (
	"github.com/pkg/errors"

	"github.com/plumming/dx/pkg/cmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Context defines Kubernetes context.
type Context struct {
	cmd.CommonOptions
	Context string
	Config  *api.Config
}

// NewContext creates a new Context.
func NewContext() *Context {
	c := &Context{}
	return c
}

// Validate input.
func (c *Context) Validate() error {
	k := c.Kuber()
	var err error
	c.Config, err = k.LoadAPIConfig()
	if err != nil {
		return err
	}

	if c.Context == "" {
		c.Context, err = c.selectContext()
		if err != nil {
			return errors.Wrap(err, "failed to select context")
		}
	}

	return nil
}

// Run the cmd.
func (c *Context) Run() error {
	k := c.Kuber()
	var err error
	c.Config, err = k.SetKubeContext(c.Context, c.Config)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) selectContext() (string, error) {
	contexts := c.loadContexts()
	prompter := c.Prompter()
	currentContext := c.Config.CurrentContext
	ctx, err := prompter.SelectFromOptionsWithDefault("Select a context:", currentContext, contexts)
	if err != nil {
		return "", errors.Wrap(err, "failed selecting context from prompter")
	}
	return ctx, nil
}

func (c *Context) loadContexts() []string {
	var contexts []string
	for k := range c.Config.Contexts {
		contexts = append(contexts, k)
	}
	return contexts
}
