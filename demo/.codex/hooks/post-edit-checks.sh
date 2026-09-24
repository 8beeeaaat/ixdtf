#!/bin/bash
# PostToolUse: add context for SSOT follow-ups and DESIGN.md guardrails.
set -o pipefail

changed_files=$(
  {
    git diff --name-only 2>/dev/null
    git diff --cached --name-only 2>/dev/null
  } | sort -u
)

ctx=""

if printf '%s\n' "$changed_files" | grep -qx 'api/openapi.yaml'; then
  ctx="${ctx}api/openapi.yaml が変更されています。make generate で orval / oapi-codegen を再実行し、生成物の diff を確認してください (api-contract-sync スキル)。"
fi

if printf '%s\n' "$changed_files" | grep -Eq '^web/src/locales/ja/'; then
  [ -n "$ctx" ] && ctx="${ctx} "
  ctx="${ctx}ja locale が変更されています。en 側にも同一キー構造で反映してください (要件 N-2)。"
fi

if printf '%s\n' "$changed_files" | grep -Eq '^web/src/locales/en/'; then
  [ -n "$ctx" ] && ctx="${ctx} "
  ctx="${ctx}en locale が変更されています。ja 側にも同一キー構造で反映してください (要件 N-2)。"
fi

ui_files=$(printf '%s\n' "$changed_files" | grep -E '\.(tsx|css)$' || true)
if [ -n "$ui_files" ]; then
  hits=$(
    while IFS= read -r file; do
      [ -f "$file" ] || continue
      grep -nE 'from "clsx"|["'"'"' ]glass(-dense|-subtle)?["'"'"' ]|dark:|text-gray-|bg-white|text-black|bg-slate-|text-\[#|bg-\[#' "$file" 2>/dev/null | sed "s#^#$file:#"
    done <<EOF
$ui_files
EOF
  )

  if [ -n "$hits" ]; then
    [ -n "$ctx" ] && ctx="${ctx} "
    ctx="${ctx}DESIGN.md 禁止パターンの疑いを検出: $(printf '%s\n' "$hits" | head -5 | paste -sd '; ' -) — docs/DESIGN.md の禁止パターン節を確認し、必要なら design-system-reviewer エージェントで検査してください。"
  fi
fi

if [ -n "$ctx" ]; then
  jq -n --arg c "$ctx" '{
    hookSpecificOutput: {
      hookEventName: "PostToolUse",
      additionalContext: $c
    }
  }'
fi

exit 0
