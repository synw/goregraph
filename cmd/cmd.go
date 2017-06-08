package cmd

import (
	"github.com/synw/microb-goregraph/state"
	"github.com/synw/microb-goregraph/state/mutate"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/terr"
)

func Dispatch(cmd *datatypes.Command) *datatypes.Command {
	com := &datatypes.Command{}
	// TODO: error handling
	if cmd.Name == "start" {
		return Start(cmd)
	} else if cmd.Name == "stop" {
		return Stop(cmd)
	}
	return com
}

func Start(cmd *datatypes.Command) *datatypes.Command {
	tr := mutate.Start()
	if tr != nil {
		cmd.Trace = tr
		cmd.Status = "error"
		terr.Debug("cmd err", tr)
		return cmd
	}
	var resp []interface{}
	resp = append(resp, "Graphql server started at "+state.Conf.Host)
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}

func Stop(cmd *datatypes.Command) *datatypes.Command {
	tr := mutate.Stop()
	if tr != nil {
		cmd.Trace = tr
		cmd.Status = "error"
		return cmd
	}
	var resp []interface{}
	resp = append(resp, "Graphql server stopped")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}
