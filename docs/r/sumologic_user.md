
# sumologic_user
Provides a [Sumologic User][1].

## Example Usage
```hcl
resource "sumologic_role" "example_role" {
  name        = "TestRole123"
  description = "Testing resource sumologic_role"

  lifecycle {
    ignore_changes = ["users"]
  }
}

resource "sumologic_user" "example_user1" {
  first_name = "Jon"
  last_name  = "Doe"
  email      = "jon.doe@gmail.com"
  active     = false
  role_ids   = ["${sumologic_role.example_role.id}"]
}

resource "sumologic_user" "example_user2" {
  first_name = "Jane"
  last_name  = "Smith"
  email      = "jane.smith@gmail.com"
  role_ids   = ["${sumologic_role.example_role.id}"]
}
```

## Argument reference


## Attributes reference


[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Manage/Users-and-Roles/Manage-Users
