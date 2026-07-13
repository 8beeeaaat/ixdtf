/**
 * IXDTF string tokenizer for highlight display (F-1-2 / F-2-5).
 *
 * DISPLAY ONLY. This never decides validity — correctness always comes from the
 * two real parsers (Go `ixdtf` on the server, native Temporal in the browser).
 * It performs a deliberately loose split so the UI can colour each component;
 * do not grow it into a third interpretation of the standard (CLAUDE.md).
 */
export type IxdtfTokenType = "date" | "time" | "offset" | "timezone" | "extension" | "invalid";

export interface IxdtfToken {
  type: IxdtfTokenType;
  text: string;
  /** Present on `timezone` / `extension` tokens: the `[!...]` critical flag. */
  critical?: boolean;
}

const DATE_RE = /^\d{4}-\d{2}-\d{2}/;
const TIME_RE = /^[Tt ]\d{2}:\d{2}(:\d{2})?(\.\d+)?/;
const OFFSET_RE = /^(Z|z|[+-]\d{2}:?\d{2})/;
const ANNOTATION_RE = /\[([^\]]*)\]/g;

/**
 * Split an IXDTF string into ordered tokens. Concatenating every token's
 * `text` reproduces the original input.
 */
export function tokenize(input: string): IxdtfToken[] {
  const tokens: IxdtfToken[] = [];

  const bracketIndex = input.indexOf("[");
  const core = bracketIndex === -1 ? input : input.slice(0, bracketIndex);
  const annotations = bracketIndex === -1 ? "" : input.slice(bracketIndex);

  // --- core: date / time / offset -----------------------------------------
  let consumed = 0;
  const dateMatch = DATE_RE.exec(core);
  if (dateMatch) {
    tokens.push({ type: "date", text: dateMatch[0] });
    consumed = dateMatch[0].length;

    const timeMatch = TIME_RE.exec(core.slice(consumed));
    if (timeMatch) {
      tokens.push({ type: "time", text: timeMatch[0] });
      consumed += timeMatch[0].length;

      const offsetMatch = OFFSET_RE.exec(core.slice(consumed));
      if (offsetMatch) {
        tokens.push({ type: "offset", text: offsetMatch[0] });
        consumed += offsetMatch[0].length;
      }
    }
  }
  // Anything in the core we could not classify is shown as invalid.
  if (consumed < core.length) {
    tokens.push({ type: "invalid", text: core.slice(consumed) });
  }

  // --- annotations: [timezone] and [key=value] extensions ------------------
  let cursor = 0;
  ANNOTATION_RE.lastIndex = 0;
  let match = ANNOTATION_RE.exec(annotations);
  while (match !== null) {
    // Non-bracket garbage between annotations.
    if (match.index > cursor) {
      tokens.push({ type: "invalid", text: annotations.slice(cursor, match.index) });
    }
    let body = match[1];
    let critical = false;
    if (body.startsWith("!")) {
      critical = true;
      body = body.slice(1);
    }
    tokens.push({
      type: body.includes("=") ? "extension" : "timezone",
      text: match[0],
      critical,
    });
    cursor = ANNOTATION_RE.lastIndex;
    match = ANNOTATION_RE.exec(annotations);
  }
  if (cursor < annotations.length) {
    tokens.push({ type: "invalid", text: annotations.slice(cursor) });
  }

  return tokens;
}
