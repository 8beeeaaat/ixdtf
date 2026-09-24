#!/bin/bash
# PreToolUse: block manual edits to codegen output.
payload=$(cat)
command=$(printf '%s' "$payload" | jq -r '.tool_input.command // ""' 2>/dev/null)

if printf '%s\n' "$command" | grep -Eq '(^|\s|\*\*\* (Add|Update|Delete) File: |/)(server/generated|web/src/generated)/'; then
  jq -n \
    --arg reason "codegen 生成物は手動編集禁止です。api/openapi.yaml を編集して make generate を実行してください (api-contract-sync スキル参照)。" \
    '{
      hookSpecificOutput: {
        hookEventName: "PreToolUse",
        permissionDecision: "deny",
        permissionDecisionReason: $reason
      }
    }'
fi

exit 0
