#!/usr/bin/env python3
"""Convert HuggingFace config.json to architecture_config format."""

import json
import sys


REQUIRED_MAPPINGS = {
    "numLayers": "num_hidden_layers",
    "hiddenSize": "hidden_size",
    "numAttentionHeads": "num_attention_heads",
}


def convert_hf_config(hf_config: dict) -> dict:
    """Convert HuggingFace config to architecture_config format."""
    arch_config = {"type": "transformer"}

    for arch_key, hf_key in REQUIRED_MAPPINGS.items():
        if hf_key not in hf_config:
            raise ValueError(f"missing required field: {hf_key}")
        value = hf_config[hf_key]
        if not isinstance(value, int) or isinstance(value, bool):
            raise ValueError(f"field {hf_key} must be an integer, got {type(value).__name__}")
        if value < 1:
            raise ValueError(f"field {hf_key} must be >= 1, got {value}")
        arch_config[arch_key] = value

    return arch_config


def main():
    if len(sys.argv) != 2:
        print(f"usage: {sys.argv[0]} <config.json>", file=sys.stderr)
        sys.exit(1)

    config_path = sys.argv[1]

    try:
        with open(config_path, "r") as f:
            hf_config = json.load(f)
    except FileNotFoundError:
        print(f"error: file not found: {config_path}", file=sys.stderr)
        sys.exit(1)
    except json.JSONDecodeError as e:
        print(f"error: invalid JSON: {e}", file=sys.stderr)
        sys.exit(1)

    try:
        arch_config = convert_hf_config(hf_config)
    except ValueError as e:
        print(f"error: {e}", file=sys.stderr)
        sys.exit(1)

    print(json.dumps(arch_config, indent=2, sort_keys=True))


if __name__ == "__main__":
    main()
