# Multi-Notification Terraform Provider - Test Report

## Test Date
July 17, 2026

## Test Environment
- **Backend**: EpD `ashuto2-www.ephemeral.in` (with backend PR #153506 deployed)
- **API Endpoint**: `https://ashuto2-api.ephemeral.in/api/`
- **Terraform Version**: v1.15.8 (darwin_arm64)
- **Go Version**: 1.26.5
- **Provider**: Local dev build from branch `multiNotificationSupport`

## Pre-requisites Completed
1. Installed Go 1.26.5 via Homebrew
2. Installed Terraform v1.15.8 via Homebrew (hashicorp/tap)
3. Built provider binary: `go build -o terraform-provider-sumologic`
4. Configured `~/.terraformrc` with dev_overrides pointing to local build
5. Created test Terraform config with multi-notification scheduled search

## Tests Performed

### Test 1: Static Code Analysis

| Step | Command | Result |
|------|---------|--------|
| Build | `go build ./...` | PASS - No compilation errors |
| Vet | `go vet ./sumologic/` | PASS - No issues found |
| Test registration | `go test ./sumologic/ -run TestAccSumologicLogSearch -count=1` | PASS - All 5 tests found and registered (skipped without TF_ACC) |

**Tests registered:**
- TestAccSumologicLogSearch_basic (SKIP)
- TestAccSumologicLogSearch_create (SKIP)
- TestAccSumologicLogSearch_update (SKIP)
- TestAccSumologicLogSearch_multi_notification (SKIP)
- TestAccSumologicLogSearch_intervalTimeType (SKIP)

---

### Test 2: Terraform Plan (Multi-Notification)

**Config:** Two email notifications on a single scheduled search

```hcl
schedule {
  notifications {
    email_search_notification {
      to_list            = ["team-alerts@example.com"]
      subject_template   = "Alert 1: {{TriggerCondition}} for {{SearchName}}"
      include_query      = true
      include_result_set = true
      include_histogram  = true
      include_csv_attachment = false
    }
  }
  notifications {
    email_search_notification {
      to_list            = ["oncall@example.com"]
      subject_template   = "Alert 2: {{TriggerCondition}} for {{SearchName}}"
      include_query      = false
      include_result_set = true
      include_histogram  = false
      include_csv_attachment = false
    }
  }
}
```

**Command:** `terraform plan`

**Result:** PASS

- Plan correctly shows 2 `notifications` blocks with all attributes
- Plan shows: `1 to add, 0 to change, 0 to destroy`

---

### Test 3: Terraform Apply (Create Resource)

**Command:** `terraform apply -auto-approve`

**Result:** PASS

- Resource created successfully
- Resource ID: `000000000000C355`
- Creation time: ~6 seconds
- Output: `Apply complete! Resources: 1 added, 0 changed, 0 destroyed.`

---

### Test 4: Terraform Plan - No Drift (Read Verification)

**Purpose:** Verify that after creating the resource, reading it back from the API produces no diff (state round-trips correctly).

**Command:** `terraform plan` (after apply)

**Result:** PASS

- Output: `No changes. Your infrastructure matches the configuration.`
- This confirms:
  - The API correctly stores multiple notifications
  - The API correctly returns `notifications` array in the response
  - The provider correctly deserializes the `notifications` array back to Terraform state
  - All notification attributes match what was sent

---

### Test 5: Terraform Destroy (Delete Resource)

**Command:** `terraform destroy -auto-approve`

**Result:** PASS

- Resource `000000000000C355` destroyed successfully
- Destruction time: ~2 seconds
- Output: `Destroy complete! Resources: 1 destroyed.`

---

## Test Summary

| # | Test | Status | Details |
|---|------|--------|---------|
| 1 | Code compilation (`go build`) | PASS | No errors |
| 2 | Code analysis (`go vet`) | PASS | No issues |
| 3 | Test registration | PASS | All 5 tests found |
| 4 | Terraform Plan | PASS | Correct plan output with 2 notifications |
| 5 | Terraform Apply (Create) | PASS | Resource ID: 000000000000C355 |
| 6 | Terraform Plan (No Drift) | PASS | No changes detected after read-back |
| 7 | Terraform Destroy | PASS | Resource deleted successfully |

## Backward Compatibility

The existing singular `notification` field is unchanged and continues to work:
- Changed from `Required` to `Optional` in schema
- `AtLeastOneOf` constraint ensures either `notification` or `notifications` is provided
- Existing Terraform configs with `notification` don't need any changes

## Known Limitation

- **CI Failure**: The `TestAccSumologicLogSearch_multi_notification` test fails in GitHub Actions CI because the CI test environment does not yet have the backend PR (#153506) deployed. The old API rejects requests without `notification` field. Once the backend is deployed to the CI test environment, this test will pass.
- **Error from old API**: `"Invalid field: schedule.notification must not be null"`

## Files Modified

| File | Change |
|------|--------|
| `sumologic/sumologic_log_search.go` | Added `Notifications []interface{}` field, added `omitempty` to `Notification` |
| `sumologic/resource_sumologic_log_search.go` | Schema: made `notification` Optional, added `notifications` field. Conversion: updated read/write paths |
| `sumologic/resource_sumologic_log_search_test.go` | Added `TestAccSumologicLogSearch_multi_notification` test and HCL generator |

## Steps to Reproduce

```bash
# 1. Build the provider
cd /Users/ashutosh.kumbhaj/terraform-provider-sumologic
go build -o terraform-provider-sumologic

# 2. Configure dev override (~/.terraformrc)
cat > ~/.terraformrc <<EOF
provider_installation {
  dev_overrides {
    "SumoLogic/sumologic" = "/Users/ashutosh.kumbhaj/terraform-provider-sumologic"
  }
  direct {}
}
EOF

# 3. Create test config directory with main.tf (see config above)
mkdir test-multi-notif && cd test-multi-notif

# 4. Run terraform
terraform plan
terraform apply -auto-approve
terraform plan          # verify no drift
terraform destroy -auto-approve

# 5. Run acceptance tests (requires TF_ACC and API credentials)
export TF_ACC=1
export SUMOLOGIC_ACCESSID="<id>"
export SUMOLOGIC_ACCESSKEY="<key>"
export SUMOLOGIC_BASE_URL="https://ashuto2-api.ephemeral.in/api/"
go test ./sumologic/ -run TestAccSumologicLogSearch_multi_notification -v -count=1 -timeout 30m
```
