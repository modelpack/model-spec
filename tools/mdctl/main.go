package main

import (
	"context"

	"github.com/CloudNativeAI/model-spec/tools/mdctl/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.NewCLI().ExecuteContext(context.Background()))
}
