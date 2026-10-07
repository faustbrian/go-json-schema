# Dialect selection and migration

Select a released dialect explicitly with `WithDialect`. Draft 2020-12 is the
documented constructor default, but applications should still select a dialect
at trust boundaries so a future schema cannot silently change interpretation.

| Constant | Meta-schema URI |
| --- | --- |
| `Draft3` | `http://json-schema.org/draft-03/schema#` |
| `Draft4` | `http://json-schema.org/draft-04/schema#` |
| `Draft6` | `http://json-schema.org/draft-06/schema#` |
| `Draft7` | `http://json-schema.org/draft-07/schema#` |
| `Draft201909` | `https://json-schema.org/draft/2019-09/schema` |
| `Draft202012` | `https://json-schema.org/draft/2020-12/schema` |

The selected dialect controls meta-validation and schemas without a
`$schema`. Embedded resources with released `$schema` declarations use their
declared semantics. Unsupported stable identifiers fail; unreleased dialects
are not exposed through the stable API.

Migration requires semantic review, not keyword renaming. In particular,
review `id`/`$id`, `$ref` siblings, boolean schemas, mathematical integers,
tuple `items`/`additionalItems` versus `prefixItems`/`items`, `dependencies`
versus `dependentRequired` and `dependentSchemas`, exclusive-bound forms,
contains limits, `$recursiveRef` versus `$dynamicRef`, `$defs`, unevaluated
keywords, vocabulary declarations, content, and format assertion policy. Run
both source and target official lanes plus application regressions.

## Regular-expression migration

`pattern` and `patternProperties` use ECMAScript regular expressions in
Unicode mode. The regexp2/v2 v2.8.2 update recognizes standard property
expressions such as `\p{Script=Greek}`. A matching `patternProperties` entry
still validates the property's value and marks that property evaluated for
`unevaluatedProperties`.

Preflight stored schemas before adopting this update. Previously accepted
nonstandard expressions such as `\p{Greek}`, `\p{GCB=RI}`, `^*`, and `(?=a)*`
now fail schema compilation with `ErrInvalidSchema`. Use
`\p{Script=Greek}` for the Greek script instead of the shorthand
`\p{Greek}`. Consult the
[ECMAScript Unicode property grammar](https://tc39.es/ecma262/2025/multipage/text-processing.html#prod-UnicodePropertyValueExpression)
when rewriting other expressions.

The built-in asserted `format: "regex"` checker also recognizes standard
Unicode properties and rejects unsupported property names. Format annotation
does not become assertion implicitly, and a registered custom `regex` format
continues to replace the built-in checker. Existing regex byte, count,
backtracking, and matching-time limits remain caller-configurable.
