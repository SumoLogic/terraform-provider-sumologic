# Multi-Notification Terraform Test Scenarios

## Prerequisites

1. Build the provider locally:
   ```bash
   cd /Users/ashutosh.kumbhaj/terraform-provider-sumologic
   go build -o terraform-provider-sumologic
   ```

2. Create a `.terraformrc` to use the local provider build:
   ```bash
   cat > ~/.terraformrc <<'EOF'
   provider_installation {
     dev_overrides {
       "sumologic/sumologic" = "/Users/ashutosh.kumbhaj/terraform-provider-sumologic"
     }
     direct {}
   }
   EOF
   ```

3. Export credentials:
   ```bash
   export SUMOLOGIC_ACCESSID="su7YFH0jPodNk4"
   export SUMOLOGIC_ACCESSKEY="98NX6py5iU6GOy7Dk9Zev2NKuwWRfqyWhzTklSdSX3pF7I0zCGISZ1hZgKinH4P6"
   export SUMOLOGIC_BASE_URL="https://test3ak-www.ephemeral.in/api/"
   ```
   Or for stag:
   ```bash
   export SUMOLOGIC_ACCESSID="suz0eMOJ5tKFsI"
   export SUMOLOGIC_ACCESSKEY="bUDBJFomDRXs6AGSzoRCRQaRGBI3DhQjuYza27wBuUrWkl27eJsVydN39nbes2fq"
   export SUMOLOGIC_BASE_URL="https://epddata.stag.sumologic.net/api/"
   ```

## Test Scenarios

### Scenario A: FF ON + Old Terraform (singular notification)
- **FF state**: BE FF ON
- **Purpose**: Verify old-style singular `notification` still works when FF is enabled
- **Expected**: Create succeeds, plan shows no drift after apply

### Scenario B: FF OFF + New Terraform (plural notifications)
- **FF state**: BE FF OFF
- **Purpose**: Verify what happens when new terraform sends `notifications` but FF is off
- **Expected**: Create may succeed but subsequent plan will show drift (API returns singular)

### Scenario C: FF OFF → create → FF ON → plan
- **FF state**: Start with FF OFF, then turn ON
- **Purpose**: Test migration path — SS created before FF, then FF enabled
- **Steps**:
  1. With FF OFF: `cd scenario_c && terraform apply`
  2. Turn FF ON
  3. Run `terraform plan` — check for drift

### Scenario D: FF ON + New Terraform + Multiple Notifications
- **FF state**: BE FF ON
- **Purpose**: Full happy path — create SS with multiple notifications
- **Expected**: Create succeeds, plan shows no drift, import works

## Running Each Scenario

```bash
cd scenario_X
terraform init
terraform plan
terraform apply
terraform plan  # should show "No changes" if no drift
terraform import sumologic_log_search.test <ID>  # test import
terraform destroy  # cleanup
```
