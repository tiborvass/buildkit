package debug

import (
	"encoding/json"
	"io"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var JSON2LLBCommand = cli.Command{
	Name:      "json2llb",
	Usage:     "convert JSON to LLB. JSON can be also passed via stdin. This command does not require the daemon to be running.",
	ArgsUsage: "<file.json>",
	Action:    json2llb,
}

func json2llb(clicontext *cli.Context) error {
	var r io.Reader
	if jsonFile := clicontext.Args().First(); jsonFile != "" && jsonFile != "-" {
		f, err := os.Open(jsonFile)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}
	def := &llb.Definition{Metadata: make(map[digest.Digest]pb.OpMetadata)}
	dec := json.NewDecoder(r)
	for {
		var op llbOp
		if err := dec.Decode(&op); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		dt, err := op.Op.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal op")
		}
		dgst := digest.FromBytes(dt)
		if dgst != op.Digest {
			//fmt.Printf("digest mismatch: expected %s, got %s from input", dgst, op.Digest)
		}
		def.Metadata[dgst] = op.OpMetadata
		def.Def = append(def.Def, dt)
	}
	return llb.WriteTo(def, os.Stdout)
}
