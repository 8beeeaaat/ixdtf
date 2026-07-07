#!/bin/bash
# PreToolUse (Edit|Write): codegen 生成物の手動編集をブロックする (architecture.md「手書きの API 型は禁止」)
file=$(jq -r '.tool_input.file_path // empty' 2>/dev/null)
case "$file" in
  */server/generated/*|*/web/src/generated/*)
    echo "BLOCKED: $file は codegen 生成物です。api/openapi.yaml を編集して 'make generate' を実行してください (api-contract-sync スキル参照)。" >&2
    exit 2
    ;;
esac
exit 0
