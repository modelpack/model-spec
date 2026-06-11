#!/usr/bin/env python3
"""Generate Python models from the canonical ModelPack JSON Schema."""

from __future__ import annotations

import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
SCHEMA_PATH = ROOT / "schema" / "config-schema.json"
OUTPUT_PATH = ROOT / "specs-python" / "modelpack" / "v1" / "models.py"


def main() -> int:
    try:
        import datamodel_code_generator  # noqa: F401
    except ModuleNotFoundError:
        print(
            "error: datamodel-code-generator is not installed. "
            "Install it with: pip install datamodel-code-generator",
            file=sys.stderr,
        )
        return 1

    if not SCHEMA_PATH.is_file():
        print(
            f"error: JSON Schema not found at: {SCHEMA_PATH}",
            file=sys.stderr,
        )
        return 1

    OUTPUT_PATH.parent.mkdir(parents=True, exist_ok=True)

    cmd = [
        sys.executable,
        "-m",
        "datamodel_code_generator",
        "--input",
        str(SCHEMA_PATH),
        "--output",
        str(OUTPUT_PATH),
        "--input-file-type",
        "jsonschema",
        "--output-model-type",
        "pydantic_v2.BaseModel",
        "--target-python-version",
        "3.10",
        "--enum-field-as-literal",
        "all",
        "--field-constraints",
        "--disable-timestamp",
    ]

    try:
        subprocess.run(cmd, check=True)
    except subprocess.CalledProcessError as exc:
        print(
            f"error: datamodel-code-generator failed with exit code {exc.returncode}.",
            file=sys.stderr,
        )
        return exc.returncode or 1
    else:
        print(f"Generated: {OUTPUT_PATH}")
        return 0


if __name__ == "__main__":
    raise SystemExit(main())
