import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { describe, expect, it } from "vitest";
import i18n from "@/app/i18n";
import { references } from "@/lib/references";

const ALLOWED_HOSTS = new Set(["www.rfc-editor.org", "tc39.es", "pkg.go.dev", "github.com"]);
const IXDTF_VERSION = "v0.4.0";

describe("official reference catalog", () => {
  it("uses primary-source hosts and version-pinned Go ixdtf links", () => {
    for (const reference of Object.values(references)) {
      for (const source of reference.sources) {
        const url = new URL(source.href);
        expect(ALLOWED_HOSTS.has(url.hostname), source.href).toBe(true);
        if (source.href.includes("8beeeaaat/ixdtf")) {
          expect(source.href, source.href).toContain(IXDTF_VERSION);
        }
      }
    }
  });

  it("keeps every catalog label available in both locales", () => {
    for (const reference of Object.values(references)) {
      for (const language of ["ja", "en"]) {
        expect(i18n.exists(reference.titleKey, { lng: language }), reference.titleKey).toBe(true);
        expect(i18n.exists(reference.bodyKey, { lng: language }), reference.bodyKey).toBe(true);
        if ("noticeKey" in reference) {
          expect(i18n.exists(reference.noticeKey, { lng: language }), reference.noticeKey).toBe(
            true,
          );
        }
      }
      for (const source of reference.sources) {
        for (const language of ["ja", "en"]) {
          expect(i18n.exists(source.labelKey, { lng: language }), source.labelKey).toBe(true);
        }
      }
    }
  });

  it("matches the Go ixdtf version used by the server", () => {
    const goMod = readFileSync(resolve(process.cwd(), "../server/go.mod"), "utf8");
    expect(goMod).toContain(`github.com/8beeeaaat/ixdtf ${IXDTF_VERSION}`);
  });
});
