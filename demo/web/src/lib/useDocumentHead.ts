import { useLocation } from "@tanstack/react-router";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { applyDocumentHead, resolveDocumentHead } from "@/lib/seo";

/**
 * ルート遷移と UI 言語の切り替え (F-0-2) に合わせて title / description / canonical を同期する (N-7)。
 * react-i18next の `t` は言語が変わると別の関数になるため、依存配列の `t` で言語切り替えにも追従する。
 */
export function useDocumentHead(): void {
  const { pathname } = useLocation();
  const { t } = useTranslation();

  useEffect(() => {
    applyDocumentHead(document, resolveDocumentHead(pathname, t));
  }, [pathname, t]);
}
