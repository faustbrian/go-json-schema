#!/usr/bin/env python3
"""Validate public Markdown links and prepare Go documentation snippets."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
from urllib.parse import unquote
import argparse
import re


EXCLUDED_PARTS = {"testdata", ".golib-tooling", ".verification"}
REFERENCE_DEFINITION = re.compile(
    r"^\s{0,3}\[([^]]+)\]:\s*(?:<([^>]+)>|(\S+))", re.MULTILINE
)
REFERENCE_USAGE = re.compile(r"(?<!!)\[([^]\n]+)\]\[([^]\n]*)\]")
INLINE_LINK = re.compile(r"(?<!!)\[[^]\n]*\]\(([^)]+)\)")
EXPECTED_OUTPUT = re.compile(r"^\s*// Output:\s*(.*)\s*$", re.MULTILINE)


@dataclass(frozen=True)
class Fence:
    info: str
    source: str


def public_documents(root: Path) -> list[Path]:
    return sorted(
        document
        for document in root.rglob("*.md")
        if not EXCLUDED_PARTS.intersection(document.parts)
    )


def split_markdown(text: str) -> tuple[str, list[Fence]]:
    prose: list[str] = []
    fences: list[Fence] = []
    source: list[str] = []
    marker = ""
    marker_length = 0
    info = ""

    for line in text.splitlines():
        if not marker:
            opening = re.match(r"^\s{0,3}(`{3,}|~{3,})(.*)$", line)
            if opening:
                marker = opening.group(1)[0]
                marker_length = len(opening.group(1))
                info = opening.group(2).strip()
                if marker == "`" and "`" in info:
                    raise ValueError("backtick fence info string contains a backtick")
                source = []
                continue
            prose.append(line)
            continue

        closing = re.match(rf"^\s{{0,3}}{re.escape(marker)}{{{marker_length},}}\s*$", line)
        if closing:
            fences.append(Fence(info=info, source="\n".join(source).rstrip() + "\n"))
            marker = ""
            marker_length = 0
            info = ""
            source = []
            continue
        source.append(line)

    if marker:
        raise ValueError("unterminated fenced code block")

    return "\n".join(prose), fences


def normalized_label(label: str) -> str:
    return " ".join(label.split()).casefold()


def markdown_links(text: str) -> tuple[list[str], list[str]]:
    prose, _ = split_markdown(text)
    prose = re.sub(r"<!--.*?-->", "", prose, flags=re.DOTALL)
    definitions = {
        normalized_label(match.group(1)): match.group(2) or match.group(3)
        for match in REFERENCE_DEFINITION.finditer(prose)
    }
    without_definitions = REFERENCE_DEFINITION.sub("", prose)
    resolved = [match.group(1) for match in INLINE_LINK.finditer(without_definitions)]
    for match in REFERENCE_USAGE.finditer(without_definitions):
        label = normalized_label(match.group(2) or match.group(1))
        target = definitions.get(label)
        if target is None:
            raise ValueError(f"undefined reference-style link label: {label}")
        resolved.append(target)
    return resolved, list(definitions.values())


def heading_anchors(document: Path) -> set[str]:
    prose, _ = split_markdown(document.read_text(encoding="utf-8"))
    seen: dict[str, int] = {}
    result: set[str] = set()
    for line in prose.splitlines():
        match = re.match(r"^#{1,6}\s+(.+?)\s*#*\s*$", line)
        if not match:
            continue
        slug = match.group(1).strip().lower()
        slug = re.sub(r"[^\w\- ]", "", slug)
        slug = re.sub(r"\s+", "-", slug)
        count = seen.get(slug, 0)
        seen[slug] = count + 1
        result.add(slug if count == 0 else f"{slug}-{count}")
    return result


def target_value(raw_target: str) -> str:
    stripped = raw_target.strip()
    if stripped.startswith("<") and ">" in stripped:
        return unquote(stripped[1 : stripped.index(">")])
    return unquote(stripped.split()[0])


def validate_links(root: Path, documents: list[Path]) -> dict[Path, list[str]]:
    anchors = {document.resolve(): heading_anchors(document) for document in documents}
    actionable: dict[Path, list[str]] = {}
    for document in documents:
        try:
            used, definitions = markdown_links(document.read_text(encoding="utf-8"))
        except ValueError as error:
            raise ValueError(f"invalid Markdown in {document.relative_to(root)}: {error}") from error
        actionable[document.resolve()] = [target_value(target) for target in used]
        for raw_target in used + definitions:
            target = target_value(raw_target)
            if target.startswith(("http://", "https://", "mailto:")):
                continue
            relative, _, fragment = target.partition("#")
            linked = (document.parent / relative).resolve() if relative else document.resolve()
            if not linked.exists():
                raise ValueError(
                    f"broken relative link in {document.relative_to(root)}: {raw_target}"
                )
            if fragment and linked.suffix.lower() == ".md":
                if fragment.lower() not in anchors.get(linked, heading_anchors(linked)):
                    raise ValueError(
                        f"broken local anchor in {document.relative_to(root)}: {raw_target}"
                    )
    return actionable


def is_go_info(info: str) -> bool:
    if not info:
        return False
    first = info.split(maxsplit=1)[0].casefold()
    return first == "go" or first == "{.go}"


def go_blocks(root: Path, documents: list[Path]) -> list[tuple[Path, str]]:
    blocks: list[tuple[Path, str]] = []
    for document in documents:
        try:
            _, fences = split_markdown(document.read_text(encoding="utf-8"))
        except ValueError as error:
            raise ValueError(f"invalid Markdown in {document.relative_to(root)}: {error}") from error
        blocks.extend((document, fence.source) for fence in fences if is_go_info(fence.info))
    return blocks


def write_snippets(root: Path, output: Path, blocks: list[tuple[Path, str]]) -> int:
    if not blocks:
        raise ValueError("public documentation contains no Go examples")

    module = "github.com/faustbrian/go-json-schema"
    (output / "go.mod").write_text(
        "module documentation-snippets\n\ngo 1.27.0\n\n"
        f"require {module} v0.0.0\n\nreplace {module} => {root}\n",
        encoding="utf-8",
    )
    fragments: list[tuple[Path, str]] = []
    standalone = 0
    for document, source in blocks:
        if re.match(r"^\s*package\s+main\b", source):
            destination = output / f"standalone-{standalone}"
            destination.mkdir()
            (destination / "main.go").write_text(source, encoding="utf-8")
            expectation = EXPECTED_OUTPUT.search(source)
            if expectation:
                (destination / ".expected-output").write_text(
                    expectation.group(1).rstrip() + "\n", encoding="utf-8"
                )
            standalone += 1
        else:
            fragments.append((document, source))

    generated = '''package snippets

import (
    "context"
    "encoding/json"
    "errors"
    "os"

    jsonschema "github.com/faustbrian/go-json-schema"
)

var (
    ctx = context.Background()
    rawSchema = []byte(`{"type":"not-a-real-type"}`)
    instance = []byte(`42`)
    addressSchema = []byte(`{"type":"string"}`)
    embeddedSchemas = map[string][]byte{
        "https://schemas.example.test/address": addressSchema,
    }
    compiler, _ = jsonschema.NewCompiler()
    schema, _ = compiler.Compile(context.Background(), []byte(`{"type":"string"}`))
    loader, _ = jsonschema.NewMapLoader(embeddedSchemas)
    encoded []byte
    err error
)

var (
    _ = json.Marshal
    _ = errors.Is
    _ = os.OpenRoot
)
'''
    for index, (document, source) in enumerate(fragments):
        generated += f"\n// Source: {document.relative_to(root)}\nfunc snippet{index}() error {{\n{source}\n"
        generated += "\t_ = compiler\n\t_ = schema\n\t_ = loader\n\t_ = encoded\n\t_ = err\n\treturn nil\n}\n"
    (output / "snippets_test.go").write_text(generated, encoding="utf-8")
    return standalone


def validate_routes(root: Path, actionable: dict[Path, list[str]]) -> None:
    required = {
        root / "SUPPORT.md": "https://github.com/faustbrian/go-json-schema/issues",
        root / "SECURITY.md": "https://github.com/faustbrian/go-json-schema/security/advisories/new",
    }
    for document, target in required.items():
        if target not in actionable.get(document.resolve(), []):
            raise ValueError(
                f"canonical actionable route is missing from {document.relative_to(root)}: {target}"
            )


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    arguments = parser.parse_args()
    root = arguments.root.resolve()
    documents = public_documents(root)
    actionable = validate_links(root, documents)
    validate_routes(root, actionable)
    count = write_snippets(root, arguments.output, go_blocks(root, documents))
    (arguments.output / "standalone-count").write_text(f"{count}\n", encoding="utf-8")
    print("documentation links, anchors, references, routes, and Go fences validated")


if __name__ == "__main__":
    main()
