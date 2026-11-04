terraform {
  required_providers {
    sumologic = { source = "sumologic/sumologic" }
  }
}

provider "sumologic" {
  base_url = "https://stag-api.sumologic.net/api/"
  environment = "stag"   # your deployment
  # access_id / access_key via SUMOLOGIC_ACCESSID / SUMOLOGIC_ACCESSKEY
}

# Top-level folder in the Parsers library root
resource "sumologic_parser_folder" "kevin" {
  name        = "kevin"
}

# Nested folder under it
resource "sumologic_parser_folder" "kevin-sub" {
  name        = "kevin-sub"
  parent_id   = sumologic_parser_folder.kevin.id
}

resource "sumologic_parser" "test" {
  name      = "tf-test-parser"
  parent_id = sumologic_parser_folder.kevin-sub.id
  stanzas   = "[sourcetype:tf-test]\nFORMAT = CSV"
}
