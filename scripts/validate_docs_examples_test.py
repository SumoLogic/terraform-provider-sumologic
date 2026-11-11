#!/usr/bin/env python3
"""Unit tests for the pure helpers in validate_docs_examples.py (no terraform, no network).

Usage: python3 -I scripts/validate_docs_examples_test.py
"""
import importlib.util
import os
import unittest

SRC = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'validate_docs_examples.py')
spec = importlib.util.spec_from_file_location('validate_docs_examples', SRC)
v = importlib.util.module_from_spec(spec)
spec.loader.exec_module(v)

RESOURCE = 'resource "sumologic_field" "f" {\n  field_name = "x"\n}\n'


class StripProviderBlocks(unittest.TestCase):
    def test_removes_terraform_block_with_nested_braces(self):
        body = 'terraform {\n  required_providers {\n    sumologic = { source = "SumoLogic/sumologic" }\n  }\n}\n' + RESOURCE
        self.assertEqual(v.strip_provider_blocks(body).strip(), RESOURCE.strip())

    def test_removes_provider_block(self):
        body = 'provider "sumologic" {\n  access_id = "x"\n}\n' + RESOURCE
        self.assertEqual(v.strip_provider_blocks(body).strip(), RESOURCE.strip())

    def test_removes_one_line_empty_provider_block(self):
        body = 'provider "sumologic" {}\n' + RESOURCE
        self.assertEqual(v.strip_provider_blocks(body).strip(), RESOURCE.strip())

    def test_keeps_resources_and_locals(self):
        body = 'locals {\n  a = { b = 1 }\n}\n' + RESOURCE
        self.assertEqual(v.strip_provider_blocks(body), body)

    def test_ignores_braces_inside_strings(self):
        body = 'terraform {\n  x = "}"\n}\n' + RESOURCE
        self.assertEqual(v.strip_provider_blocks(body).strip(), RESOURCE.strip())

    def test_provider_only_example_is_not_an_example(self):
        body = 'terraform {\n}\n'
        self.assertEqual(v.examples_from_text('```hcl\n' + body + '```\n'), [])


class ForeignProviders(unittest.TestCase):
    def test_flags_example_using_another_provider(self):
        self.assertTrue(v.uses_other_providers('resource "aws_s3_bucket" "b" {}\n' + RESOURCE))
        self.assertTrue(v.uses_other_providers('data "aws_caller_identity" "me" {}\n' + RESOURCE))

    def test_sumologic_only_is_fine(self):
        self.assertFalse(v.uses_other_providers(RESOURCE + 'data "sumologic_personal_folder" "p" {}\n'))


class FileStubs(unittest.TestCase):
    def test_finds_path_module_files(self):
        body = 'x = file("${path.module}/lookup.csv")\ny = jsondecode(file("${path.module}/playbooks.json"))\n'
        self.assertEqual(sorted(v.file_stubs(body)), [('lookup.csv', 'x'), ('playbooks.json', '{}')])

    def test_no_file_calls(self):
        self.assertEqual(v.file_stubs(RESOURCE), [])


class ExampleLines(unittest.TestCase):
    PAGE = 'intro\n\n```hcl\n' + RESOURCE + '```\n\ntext\n```hcl\nterraform {\n}\n```\n\n```hcl\n' + RESOURCE + '```\n'

    def test_line_is_the_fence_line(self):
        self.assertEqual([line for line, _ in v.examples_with_lines(self.PAGE)], [3, 15])

    def test_skips_examples_without_resources(self):
        self.assertEqual(len(v.examples_with_lines(self.PAGE)), 2)

    def test_examples_from_text_matches(self):
        self.assertEqual(v.examples_from_text(self.PAGE), [b for _, b in v.examples_with_lines(self.PAGE)])


class ErrorText(unittest.TestCase):
    OUT = (
        'Error: Reference to undeclared resource\n\n'
        '  on main.tf line %d, in resource "sumologic_slo" "s":\n'
        '  %d:   parent_id = sumologic_slo_folder.f.id\n\n'
        'A managed resource "sumologic_slo_folder" "f" has not been declared in the root module.\n'
    )

    def test_lines_are_relative_to_the_example(self):
        n = v.HEADER_LINES + 4
        out = v.example_relative(self.OUT % (n, n))
        self.assertIn('at line 4 of the example', out)
        self.assertIn('  4:   parent_id', out)
        self.assertNotIn('main.tf', out)

    def test_keeps_the_full_message(self):
        self.assertIn('has not been declared', v.example_relative(self.OUT % (20, 20)))

    def test_first_error(self):
        self.assertEqual(v.first_error(self.OUT % (20, 20)), 'Error: Reference to undeclared resource')
        self.assertEqual(v.first_error('Success! The configuration is valid.'), 'Success! The configuration is valid.')
        self.assertEqual(v.first_error(''), '')


class Warnings(unittest.TestCase):
    OUT = 'Warning: Provider development overrides are in effect\n\nThe following overrides are set.\n\nError: Bad thing\n\nDetails here.\n'

    def test_drops_warning_keeps_error(self):
        out = v.drop_warnings(self.OUT)
        self.assertNotIn('Warning', out)
        self.assertTrue(out.startswith('Error: Bad thing'))
        self.assertIn('Details here.', out)

    def test_success_output_is_kept_when_only_a_warning_is_present(self):
        out = 'Warning: Provider development overrides are in effect\n\nbody\n\nSuccess! The configuration is valid.\n'
        self.assertEqual(v.drop_warnings(out), 'Success! The configuration is valid.')


class Annotations(unittest.TestCase):
    def test_escapes_newlines_and_percent(self):
        self.assertEqual(v.escape_annotation('a\nb%c\r'), 'a%0Ab%25c%0D')


if __name__ == '__main__':
    unittest.main()
