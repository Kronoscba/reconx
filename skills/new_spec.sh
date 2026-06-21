#!/bin/bash
set -euo pipefail

if [ -z "${1:-}" ]; then
    echo "❌ Uso: ./skills/new_spec.sh 'Feature Name'" >&2
    exit 1
fi

FEATURE_NAME="$1"
PROJECT_ROOT="$(pwd)"
[ -d "../docs/specs" ] && PROJECT_ROOT="$(cd .. && pwd)"

SPECS_DIR="$PROJECT_ROOT/docs/specs"
TEMPLATE_FILE="$SPECS_DIR/template.md"
SAFE_NAME=$(echo "$FEATURE_NAME" | tr '[:upper:]' '[:lower:]' | tr ' ' '_' | tr -cd 'a-z0-9_-')
OUTPUT_FILE="$SPECS_DIR/${SAFE_NAME}.md"

mkdir -p "$SPECS_DIR"

if [ ! -f "$TEMPLATE_FILE" ]; then
    echo "❌ Error: No existe template.md" >&2
    exit 1
fi

if [ -e "$OUTPUT_FILE" ]; then
    echo "❌ Error: Ya existe $OUTPUT_FILE" >&2
    exit 1
fi

cp "$TEMPLATE_FILE" "$OUTPUT_FILE"

# Reemplazo seguro (GNU/BSD sed)
if sed --version >/dev/null 2>&1; then
    sed -i "s/\[Name\]/$FEATURE_NAME/g" "$OUTPUT_FILE"
else
    sed -i.bak "s/\[Name\]/$FEATURE_NAME/g" "$OUTPUT_FILE" && rm -f "$OUTPUT_FILE.bak"
fi

echo "✅ Spec creada: $OUTPUT_FILE"
