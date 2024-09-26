package cmd

import (
	"fmt"

	v2 "github.com/CloudNativeAI/model-spec/specs-go/v2"
	"github.com/CloudNativeAI/model-spec/tools/mdctl/format"
	"github.com/CloudNativeAI/model-spec/tools/mdctl/models"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
)

func BuildModel(commands []format.Command) error {
	manifest := v2.Manifest{MediaType: v2.MediaTypeModelManifest}
	config := v2.Config{}
	weights := v2.Weights{}
	engine := v2.Engine{}

	if len(commands) == 0 {
		return fmt.Errorf("modelfile has no command")
	}

	var modelName string
	if commands[0].Name == format.CREATE {
		modelName = commands[0].Args
		fmt.Println("Create", modelName)
	} else if commands[0].Name == format.FROM {
		modelName = commands[0].Args
		fmt.Println("From ", modelName)
		if err := PullModel(modelName); err != nil {
			return fmt.Errorf("failed to pull base model")
		}
		_, err := FetchManifest(modelName, &manifest, &config)
		if err != nil {
			return fmt.Errorf("failed to get remote manifest")
		}
	} else {
		return fmt.Errorf("first command should be %s or %s", format.CREATE, format.FROM)
	}
	for _, c := range commands {
		switch c.Name {
		case format.CREATE, format.FROM:
			config.Name = c.Args

		case format.NAME:
			config.Name = c.Args

		case format.DESCRIPTION:
			layer, err := models.BuildDescriptor(models.TAR, c.Args, v2.MediaTypeModelDescription, "Description")
			if err != nil {
				return fmt.Errorf("failed to build description layer: %w", err)
			}
			config.Description = append(config.Description, *layer)
			fmt.Printf("Add description [%s]\n", c.Args)

		case format.LICENSE:
			layer, err := models.BuildDescriptor(models.TAR, c.Args, v2.MediaTypeModelLicense, "License")
			if err != nil {
				return fmt.Errorf("failed to build license layer: %w", err)
			}
			config.License = append(config.License, *layer)
			fmt.Printf("Add license [%s]\n", c.Args)

		case format.ARCHITECTURE:
			config.Architecture = c.Args

		case format.FAMILY:
			config.Family = c.Args

		case format.CONFIG:
			layer, err := models.BuildDescriptor(models.TAR, c.Args, v2.MediaTypeModelConfig, "")
			if err != nil {
				return fmt.Errorf("failed to build config layer: %w", err)
			}
			config.Extensions = append(config.Extensions, *layer)
			fmt.Printf("Add config [%s]\n", c.Args)

		case format.PARAM_SIZE:
			engine.Name = c.Args

		case format.FORMAT:
			weights.Format = c.Args

		case format.WEIGHTS:
			layer, err := models.BuildDescriptor(models.TAR, c.Args, v2.MediaTypeModelWeights, "")
			if err != nil {
				return fmt.Errorf("failed to build weights layer: %w", err)
			}
			weights.File = append(weights.File, *layer)
			fmt.Printf("Add weights [%s]\n", c.Args)

		case format.TOKENIZER:
			layer, err := models.BuildDescriptor(models.TAR, c.Args, v2.MediaTypeModelProcessorText, "")
			if err != nil {
				return fmt.Errorf("failed to build tokenizer layer: %w", err)
			}
			manifest.Processor = append(manifest.Processor, *layer)
			fmt.Printf("Add tokenizer [%s]\n", c.Args)

		default:
			fmt.Printf("WARN: [%s] - [%s] not handled\n", c.Name, c.Args)
		}
	}

	manifest.Config = config
	manifest.Weights = weights
	manifest.Engine = engine

	// Commit layers
	_, err := Commit(&manifest)
	if err != nil {
		return fmt.Errorf("failed to commit layers: %w", err)
	}

	// Commit manifest layer
	if err := models.WriteManifest(modelName, &manifest); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	fmt.Println("Build succeed")
	return nil
}

func Commit(m *v2.Manifest) (bool, error) {
	layerGroups := []struct {
		name   string
		layers []oci.Descriptor
	}{
		{"Description", m.Config.Description},
		{"License", m.Config.License},
		{"Extensions", m.Config.Extensions},
		{"Weights", m.Weights.File},
		{"Tokenizer", m.Processor},
	}

	var committed bool
	for _, group := range layerGroups {
		if len(group.layers) == 0 {
			continue // Skip empty layer groups
		}
		groupCommitted, err := commitLayers(group.name, group.layers)
		if err != nil {
			return false, fmt.Errorf("failed to commit %s layers: %w", group.name, err)
		}
		committed = committed || groupCommitted
	}
	return committed, nil
}

func commitLayers(groupName string, layers []oci.Descriptor) (bool, error) {
	var groupCommitted bool
	for _, layer := range layers {
		layerCommitted, err := commitSingleLayer(groupName, layer)
		if err != nil {
			return false, err
		}
		groupCommitted = groupCommitted || layerCommitted
	}
	return groupCommitted, nil
}

func commitSingleLayer(groupName string, layer oci.Descriptor) (bool, error) {
	committed, err := models.Commit(layer)
	if err != nil {
		return false, fmt.Errorf("failed to commit %s layer: %w", groupName, err)
	}

	return committed, nil
}
