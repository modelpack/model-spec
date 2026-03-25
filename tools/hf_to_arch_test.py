#!/usr/bin/env python3
"""Tests for hf_to_arch.py"""

import json
import subprocess
import sys
import tempfile
import os

SCRIPT_PATH = os.path.join(os.path.dirname(__file__), "hf_to_arch.py")


def run_script(config_content: str) -> tuple:
    """Run hf_to_arch.py with given config content, return (exitcode, stdout, stderr)."""
    with tempfile.NamedTemporaryFile(mode="w", suffix=".json", delete=False) as f:
        f.write(config_content)
        f.flush()
        temp_path = f.name

    try:
        result = subprocess.run(
            [sys.executable, SCRIPT_PATH, temp_path],
            capture_output=True,
            text=True,
        )
        return result.returncode, result.stdout, result.stderr
    finally:
        os.unlink(temp_path)


def test_valid_config():
    """Valid HuggingFace config produces correct output."""
    config = json.dumps({
        "num_hidden_layers": 32,
        "hidden_size": 4096,
        "num_attention_heads": 32,
        "vocab_size": 32000,
    })

    exitcode, stdout, stderr = run_script(config)

    assert exitcode == 0, f"expected exit 0, got {exitcode}: {stderr}"
    output = json.loads(stdout)
    assert output == {
        "type": "transformer",
        "numLayers": 32,
        "hiddenSize": 4096,
        "numAttentionHeads": 32,
    }, f"unexpected output: {output}"
    print("PASS: test_valid_config")


def test_missing_field():
    """Missing required field produces error."""
    config = json.dumps({
        "num_hidden_layers": 32,
        "hidden_size": 4096,
    })

    exitcode, stdout, stderr = run_script(config)

    assert exitcode != 0, "expected non-zero exit for missing field"
    assert "num_attention_heads" in stderr, f"error should mention missing field: {stderr}"
    print("PASS: test_missing_field")


def test_invalid_json():
    """Invalid JSON produces error."""
    exitcode, stdout, stderr = run_script("not valid json {")

    assert exitcode != 0, "expected non-zero exit for invalid JSON"
    assert "invalid JSON" in stderr.lower() or "json" in stderr.lower(), f"error should mention JSON: {stderr}"
    print("PASS: test_invalid_json")


def test_file_not_found():
    """Non-existent file produces error."""
    result = subprocess.run(
        [sys.executable, SCRIPT_PATH, "/nonexistent/path/config.json"],
        capture_output=True,
        text=True,
    )

    assert result.returncode != 0, "expected non-zero exit for missing file"
    assert "not found" in result.stderr.lower(), f"error should mention file not found: {result.stderr}"
    print("PASS: test_file_not_found")


def test_invalid_field_type():
    """Non-integer field produces error."""
    config = json.dumps({
        "num_hidden_layers": "32",
        "hidden_size": 4096,
        "num_attention_heads": 32,
    })

    exitcode, stdout, stderr = run_script(config)

    assert exitcode != 0, "expected non-zero exit for invalid type"
    assert "integer" in stderr.lower(), f"error should mention type: {stderr}"
    print("PASS: test_invalid_field_type")


def test_zero_value():
    """Zero value produces error."""
    config = json.dumps({
        "num_hidden_layers": 0,
        "hidden_size": 4096,
        "num_attention_heads": 32,
    })

    exitcode, stdout, stderr = run_script(config)

    assert exitcode != 0, "expected non-zero exit for zero value"
    assert ">= 1" in stderr, f"error should mention minimum: {stderr}"
    print("PASS: test_zero_value")


def test_bool_value():
    """Boolean value produces error."""
    config = json.dumps({
        "num_hidden_layers": True,
        "hidden_size": 4096,
        "num_attention_heads": 32,
    })

    exitcode, stdout, stderr = run_script(config)

    assert exitcode != 0, "expected non-zero exit for bool value"
    assert "integer" in stderr.lower(), f"error should mention type: {stderr}"
    print("PASS: test_bool_value")


def main():
    test_valid_config()
    test_missing_field()
    test_invalid_json()
    test_file_not_found()
    test_invalid_field_type()
    test_zero_value()
    test_bool_value()
    print("\nAll tests passed.")


if __name__ == "__main__":
    main()
