package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/CloudNativeAI/model-spec/tools/mdctl/format"
	"github.com/CloudNativeAI/model-spec/tools/mdctl/progress"
	"github.com/spf13/cobra"
)

func BuildHandler(cmd *cobra.Command, args []string) error {
	filename, _ := cmd.Flags().GetString("file")
	filename, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	p := progress.NewProgress(os.Stderr)
	// defer p.Stop()
	// bars := make(map[string]*progress.Bar)

	modelFile, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read modelfile: %w", err)
	}

	commands, err := format.Parse(bytes.NewReader(modelFile))
	if err != nil {
		return fmt.Errorf("failed to parse modelfile: %w", err)
	}

	// status := "building"
	// spinner := progress.NewSpinner(status)
	// p.Add(status, spinner)

	if err := BuildModel(commands); err != nil {
		return fmt.Errorf("failed to build model: %w", err)
	}
	p.StopAndClear()

	return nil
}

func RunHandler(cmd *cobra.Command, _ []string) error {
	name, _ := cmd.Flags().GetString("name")
	fmt.Println("Unpack Model: ", name)
	if err := RunModel(name); err != nil {
		return fmt.Errorf("failed to unpack model: %w", err)
	}
	fmt.Println("Unpack succeed")
	return nil
}

func PushHandler(cmd *cobra.Command, _ []string) error {
	name, _ := cmd.Flags().GetString("name")
	fmt.Println("Push Model:", name)
	if err := PushModel(name); err != nil {
		return fmt.Errorf("failed to push model: %w", err)
	}
	return nil
}

func PullHandler(cmd *cobra.Command, _ []string) error {
	name, _ := cmd.Flags().GetString("name")
	fmt.Println("Pull Model:", name)
	if err := PullModel(name); err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	return nil
}

func ListHandler(cmd *cobra.Command, args []string) error {
	return ListModel()
}

func NewCLI() *cobra.Command {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:           "mdctl",
		Short:         "Model management tool",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Print(cmd.UsageString())
		},
	}

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "build models from a Modelfile",
		Args:  cobra.ExactArgs(0),
		RunE:  BuildHandler,
	}
	buildCmd.Flags().StringP("file", "f", "Modelfile", "Path to the Modelfile")

	runCmd := &cobra.Command{
		Use:   "unpack",
		Short: "run a model",
		Args:  cobra.ExactArgs(0),
		RunE:  RunHandler,
	}
	runCmd.Flags().StringP("name", "n", "", "URL of the model")

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "push a model",
		Args:  cobra.ExactArgs(0),
		RunE:  PushHandler,
	}
	pushCmd.Flags().StringP("name", "n", "", "URL of the model")

	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "pull a model",
		Args:  cobra.ExactArgs(0),
		RunE:  PullHandler,
	}
	pullCmd.Flags().StringP("name", "n", "", "URL of the model")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list models",
		Args:  cobra.ExactArgs(0),
		RunE:  ListHandler,
	}

	rootCmd.AddCommand(
		buildCmd,
		listCmd,
		runCmd,
		pushCmd,
		pullCmd,
	)

	return rootCmd
}
