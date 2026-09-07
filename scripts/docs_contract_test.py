#!/usr/bin/env python3

from pathlib import Path
from tempfile import TemporaryDirectory
import unittest

import docs_contract


class DocumentationContractTest(unittest.TestCase):
    def test_reference_usages_require_definitions(self) -> None:
        with self.assertRaisesRegex(ValueError, "undefined reference-style link"):
            docs_contract.markdown_links("[API][missing]\n")
        with self.assertRaisesRegex(ValueError, "undefined reference-style link"):
            docs_contract.markdown_links("[API][]\n")

    def test_reference_usages_resolve_case_insensitive_labels(self) -> None:
        links, _ = docs_contract.markdown_links(
            "[API][Guide]\n[guide]: docs/api.md\n[Support][]\n[support]: SUPPORT.md\n"
        )
        self.assertEqual(links, ["docs/api.md", "SUPPORT.md"])

    def test_go_fences_accept_tildes_and_info_attributes(self) -> None:
        with TemporaryDirectory() as directory:
            root = Path(directory)
            document = root / "README.md"
            document.write_text("~~~go title=example\npackage main\n~~~\n", encoding="utf-8")
            self.assertEqual(
                docs_contract.go_blocks(root, [document]),
                [(document, "package main\n")],
            )

    def test_commented_route_is_not_actionable(self) -> None:
        links, _ = docs_contract.markdown_links(
            "<!-- [Issues](https://github.com/example/issues) -->\n"
        )
        self.assertEqual(links, [])

    def test_root_policy_pages_are_public_documents(self) -> None:
        with TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "SUPPORT.md").write_text("# Support\n", encoding="utf-8")
            (root / "SECURITY.md").write_text("# Security\n", encoding="utf-8")
            self.assertEqual(
                docs_contract.public_documents(root),
                [root / "SECURITY.md", root / "SUPPORT.md"],
            )

    def test_routes_require_resolved_markdown_links(self) -> None:
        with TemporaryDirectory() as directory:
            root = Path(directory)
            support = root / "SUPPORT.md"
            security = root / "SECURITY.md"
            support.write_text("# Support\n", encoding="utf-8")
            security.write_text("# Security\n", encoding="utf-8")
            with self.assertRaisesRegex(ValueError, "SUPPORT.md"):
                docs_contract.validate_routes(
                    root,
                    {
                        support.resolve(): [],
                        security.resolve(): [
                            "https://github.com/faustbrian/go-json-schema/security/advisories/new"
                        ],
                    },
                )

    def test_only_standalone_programs_with_output_oracles_are_designated(self) -> None:
        with TemporaryDirectory() as directory:
            root = Path(directory)
            output = root / "output"
            output.mkdir()
            document = root / "README.md"
            blocks = [
                (document, 'package main\nfunc main() {}\n'),
                (document, 'package main\nfunc main() {}\n// Output: true\n'),
            ]
            self.assertEqual(docs_contract.write_snippets(root, output, blocks), 2)
            self.assertFalse((output / "standalone-0" / ".expected-output").exists())
            self.assertEqual(
                (output / "standalone-1" / ".expected-output").read_text(encoding="utf-8"),
                "true\n",
            )


if __name__ == "__main__":
    unittest.main()
