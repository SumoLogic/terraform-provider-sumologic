# sumologic_folder
Provides the ability to create, read, delete, update, and manage of folders.

## Example Usage
```hcl
resource "sumologic_folder" "folder" {
  name        = "test-folder"
  description = "Just testing this"
  parent_id   = "<Hexadecimal ID of the parent folder>"
}
```

## Argument reference
The following arguments are supported:
- `name` - (Required) The name of the folder. This is required, and has to be unique.
- `parent_id` - (Required) The ID of the folder in which you want to create the new folder.
- `description` - (Optional) The description of the folder.

## Additional data provided in state
- `created_at` - (Computed) When the folder was created.
- `created_by` - (Computed) Who created the folder.
- `modified_at` - (Computed) When was the folder last modified.
- `modified_by` - (Computed) The ID of the user who modified the folder last.
- `item_type` - (Computed) What the type of the content item is (will obviously be "Folder").
- `permissions` - (Computed) List of permissions the user has on the content item.
- `children` - (Computed) A list of all the content items in the created folder (can be folders or other content items).
