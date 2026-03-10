locals {
  dynamo_tables = {
    forms = {
      hash_key  = "PK"
      range_key = "SK"
    }
    form_submissions = {
      hash_key  = "PK"
      range_key = "SK"
    }
    tags = {
      hash_key  = "PK"
      range_key = "SK"
    }
  }
}

resource "aws_dynamodb_table" "tables" {
  for_each     = local.dynamo_tables
  name         = each.key
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = each.value.hash_key
  range_key    = each.value.range_key

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }

  server_side_encryption {
    enabled = true
  }

  tags = { Name = each.key }
}
