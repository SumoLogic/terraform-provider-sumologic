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


if __name__ == '__main__':
    unittest.main()
