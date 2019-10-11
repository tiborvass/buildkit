// +build !dfrunnocache

package dockerfile2llb

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

func dispatchRunNoCache(c *instructions.RunCommand) (llb.RunOption, error) {
	return nil, nil
}
