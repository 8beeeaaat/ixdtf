type ReferenceSourceKind = "spec" | "docs" | "code";

interface ReferenceSource {
  kind: ReferenceSourceKind;
  labelKey: string;
  href: string;
}

interface ReferenceDefinition {
  titleKey: string;
  bodyKey: string;
  noticeKey?: string;
  sources: readonly ReferenceSource[];
}

const references = {
  ixdtf: {
    titleKey: "references.topics.ixdtf.title",
    bodyKey: "references.topics.ixdtf.body",
    sources: [
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-3",
      },
      {
        kind: "spec",
        labelKey: "references.sources.rfc3339",
        href: "https://www.rfc-editor.org/rfc/rfc3339.html",
      },
      {
        kind: "docs",
        labelKey: "references.sources.ixdtfPackage",
        href: "https://pkg.go.dev/github.com/8beeeaaat/ixdtf@v0.4.0",
      },
    ],
  },
  temporalApi: {
    titleKey: "references.topics.temporalApi.title",
    bodyKey: "references.topics.temporalApi.body",
    noticeKey: "references.topics.temporalApi.notice",
    sources: [
      {
        kind: "docs",
        labelKey: "references.sources.temporalDocs",
        href: "https://tc39.es/proposal-temporal/docs/",
      },
      {
        kind: "spec",
        labelKey: "references.sources.temporalSpec",
        href: "https://tc39.es/proposal-temporal/",
      },
    ],
  },
  temporalIxdtf: {
    titleKey: "references.topics.temporalIxdtf.title",
    bodyKey: "references.topics.temporalIxdtf.body",
    noticeKey: "references.topics.temporalIxdtf.notice",
    sources: [
      {
        kind: "docs",
        labelKey: "references.sources.temporalStrings",
        href: "https://tc39.es/proposal-temporal/docs/strings.html",
      },
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-3",
      },
    ],
  },
  goParsing: {
    titleKey: "references.topics.goParsing.title",
    bodyKey: "references.topics.goParsing.body",
    sources: [
      {
        kind: "docs",
        labelKey: "references.sources.ixdtfPackage",
        href: "https://pkg.go.dev/github.com/8beeeaaat/ixdtf@v0.4.0",
      },
      {
        kind: "code",
        labelKey: "references.sources.ixdtfParseCode",
        href: "https://github.com/8beeeaaat/ixdtf/blob/v0.4.0/parse.go",
      },
    ],
  },
  temporalParsing: {
    titleKey: "references.topics.temporalParsing.title",
    bodyKey: "references.topics.temporalParsing.body",
    noticeKey: "references.topics.temporalParsing.notice",
    sources: [
      {
        kind: "docs",
        labelKey: "references.sources.temporalZonedDateTime",
        href: "https://tc39.es/proposal-temporal/docs/zoneddatetime.html#static-methods",
      },
      {
        kind: "docs",
        labelKey: "references.sources.temporalStrings",
        href: "https://tc39.es/proposal-temporal/docs/strings.html",
      },
    ],
  },
  offsetConsistency: {
    titleKey: "references.topics.offsetConsistency.title",
    bodyKey: "references.topics.offsetConsistency.body",
    noticeKey: "references.topics.offsetConsistency.notice",
    sources: [
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557Consistency",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-3.4",
      },
      {
        kind: "docs",
        labelKey: "references.sources.temporalZonedDateTime",
        href: "https://tc39.es/proposal-temporal/docs/zoneddatetime.html#static-methods",
      },
      {
        kind: "code",
        labelKey: "references.sources.ixdtfTimezoneCode",
        href: "https://github.com/8beeeaaat/ixdtf/blob/v0.4.0/timezone.go",
      },
    ],
  },
  roundtrip: {
    titleKey: "references.topics.roundtrip.title",
    bodyKey: "references.topics.roundtrip.body",
    noticeKey: "references.topics.roundtrip.notice",
    sources: [
      {
        kind: "docs",
        labelKey: "references.sources.temporalZonedDateTime",
        href: "https://tc39.es/proposal-temporal/docs/zoneddatetime.html",
      },
      {
        kind: "code",
        labelKey: "references.sources.ixdtfFormatCode",
        href: "https://github.com/8beeeaaat/ixdtf/blob/v0.4.0/format.go",
      },
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557Syntax",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-4.1",
      },
    ],
  },
  dstAmbiguity: {
    titleKey: "references.topics.dstAmbiguity.title",
    bodyKey: "references.topics.dstAmbiguity.body",
    noticeKey: "references.topics.dstAmbiguity.notice",
    sources: [
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-3",
      },
      {
        kind: "docs",
        labelKey: "references.sources.temporalAmbiguity",
        href: "https://tc39.es/proposal-temporal/docs/timezone.html#ambiguity-due-to-dst-or-other-time-zone-offset-changes",
      },
      {
        kind: "docs",
        labelKey: "references.sources.temporalZonedDateTime",
        href: "https://tc39.es/proposal-temporal/docs/zoneddatetime.html#static-methods",
      },
    ],
  },
  zoneArithmetic: {
    titleKey: "references.topics.zoneArithmetic.title",
    bodyKey: "references.topics.zoneArithmetic.body",
    noticeKey: "references.topics.zoneArithmetic.notice",
    sources: [
      {
        kind: "docs",
        labelKey: "references.sources.temporalZonedDateTimeAdd",
        href: "https://tc39.es/proposal-temporal/docs/zoneddatetime.html#add",
      },
      {
        kind: "docs",
        labelKey: "references.sources.temporalAmbiguity",
        href: "https://tc39.es/proposal-temporal/docs/timezone.html#ambiguity-due-to-dst-or-other-time-zone-offset-changes",
      },
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-3",
      },
    ],
  },
  calendarAnnotation: {
    titleKey: "references.topics.calendarAnnotation.title",
    bodyKey: "references.topics.calendarAnnotation.body",
    noticeKey: "references.topics.calendarAnnotation.notice",
    sources: [
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557Keys",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-3.2",
      },
      {
        kind: "docs",
        labelKey: "references.sources.temporalCalendars",
        href: "https://tc39.es/proposal-temporal/docs/calendars.html",
      },
      {
        kind: "code",
        labelKey: "references.sources.ixdtfSuffixCode",
        href: "https://github.com/8beeeaaat/ixdtf/blob/v0.4.0/suffix.go",
      },
    ],
  },
  suffixSyntax: {
    titleKey: "references.topics.suffixSyntax.title",
    bodyKey: "references.topics.suffixSyntax.body",
    sources: [
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557Syntax",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-4.1",
      },
      {
        kind: "docs",
        labelKey: "references.sources.temporalStrings",
        href: "https://tc39.es/proposal-temporal/docs/strings.html",
      },
      {
        kind: "code",
        labelKey: "references.sources.ixdtfSuffixCode",
        href: "https://github.com/8beeeaaat/ixdtf/blob/v0.4.0/suffix.go",
      },
    ],
  },
  criticalAnnotations: {
    titleKey: "references.topics.criticalAnnotations.title",
    bodyKey: "references.topics.criticalAnnotations.body",
    sources: [
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557Critical",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-3.3",
      },
      {
        kind: "code",
        labelKey: "references.sources.ixdtfValidateCode",
        href: "https://github.com/8beeeaaat/ixdtf/blob/v0.4.0/validate.go",
      },
    ],
  },
  rejectedExtensions: {
    titleKey: "references.topics.rejectedExtensions.title",
    bodyKey: "references.topics.rejectedExtensions.body",
    sources: [
      {
        kind: "spec",
        labelKey: "references.sources.rfc9557Keys",
        href: "https://www.rfc-editor.org/rfc/rfc9557.html#section-3.2",
      },
      {
        kind: "spec",
        labelKey: "references.sources.rfc6648",
        href: "https://www.rfc-editor.org/rfc/rfc6648.html",
      },
      {
        kind: "code",
        labelKey: "references.sources.ixdtfAbnfCode",
        href: "https://github.com/8beeeaaat/ixdtf/blob/v0.4.0/abnf/abnf.go#L162-L170",
      },
    ],
  },
} as const satisfies Record<string, ReferenceDefinition>;

type ReferenceId = keyof typeof references;

export type { ReferenceDefinition, ReferenceId, ReferenceSource, ReferenceSourceKind };
export { references };
