#!/bin/bash
# PostToolUse (Edit|Write): SSOT 追従リマインドと DESIGN.md 禁止パターンの早期警告
file=$(jq -r '.tool_input.file_path // empty' 2>/dev/null)
ctx=""
case "$file" in
  */api/openapi.yaml)
    ctx="api/openapi.yaml が変更されました。'make generate' で orval / oapi-codegen を再実行し、生成物の diff を確認してください (api-contract-sync スキル)。"
    ;;
  */web/src/locales/ja/*)
    ctx="ja locale が変更されました。en 側にも同一キー構造で反映してください (要件 N-2)。"
    ;;
  */web/src/locales/en/*)
    ctx="en locale が変更されました。ja 側にも同一キー構造で反映してください (要件 N-2)。"
    ;;
  *.tsx|*.css)
    hits=$(grep -nE 'from "clsx"|[\"'"'"' ]glass(-dense|-subtle)?[\"'"'"' ]|dark:|text-gray-|bg-white|text-black|bg-slate-|text-\[#|bg-\[#' "$file" 2>/dev/null | head -5)
    if [ -n "$hits" ]; then
      ctx="DESIGN.md 禁止パターンの疑いを検出: ${hits} — docs/DESIGN.md の禁止パターン節を確認し、必要なら design-system-reviewer エージェントで検査してください。"
    fi
    ;;
esac
if [ -n "$ctx" ]; then
  jq -n --arg c "$ctx" '{hookSpecificOutput:{hookEventName:"PostToolUse",additionalContext:$c}}'
fi
exit 0
