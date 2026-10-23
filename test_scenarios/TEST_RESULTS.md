# Multi-Notification API Test Results

**Date:** 2026-08-13  
**Environment:** epddata.stag.sumologic.net  
**FF State:** BE FF ON (`concierge.enable.multi.notification.scheduled.search`)

## Test Results (FF ON - stag)

### Scenario A: FF ON + Old Terraform (singular notification)
| Step | Result |
|------|--------|
| Create with `notification` (singular) | SUCCESS |
| Read back | `notification`: PRESENT, `notifications`: null |
| Drift check | **PASS** — API preserves singular field |

### Scenario D: FF ON + New Terraform (multiple notifications)
| Step | Result |
|------|--------|
| Create with `notifications` (2 items) | SUCCESS |
| Read back | `notification`: null, `notifications`: PRESENT (2 items) |
| Drift check | **PASS** — API returns same plural field with correct count |
| Update to 3 notifications | SUCCESS |
| Post-update read | `notifications`: 3 items |
| Post-update drift check | **PASS** |

### Scenario B/C: Backward compatibility with FF ON
| Step | Result |
|------|--------|
| Read SS created with singular `notification` (Scenario A) | Returns `notification` (singular), `notifications` null |
| Read SS created BEFORE FF was enabled ("SearchBeforeFF") | Returns `notification` (singular), `notifications` null |
| Conclusion | **No drift** for existing users when FF is turned ON |

## Test Results (FF OFF - ephemeral env)

### Scenario A: FF ON → OFF (singular notification)
| Step | Result |
|------|--------|
| Create with `notification` (singular) | SUCCESS |
| Read back | `notification`: PRESENT, `notifications`: null |
| Drift check | **PASS** |

### Scenario B: FF OFF + plural notifications
| Step | Result |
|------|--------|
| Create with `notifications` (2 items) | SUCCESS (API accepts it) |
| Read back with FF OFF | `notification`: 1st item only, `notifications`: null |
| Drift check | **FAIL** — config has 2 notifications, API returns 1 |
| Data loss? | **NO** — data preserved, reappears when FF ON |

### Scenario C: Created with FF OFF → FF turned ON
| Step | Result |
|------|--------|
| SS created with `notifications` (plural) while FF OFF | Stored correctly |
| Read with FF ON | `notifications`: PRESENT (2 items restored) |
| Drift check | **PASS** (once FF is ON) |

## Summary

| Scenario | FF State | Config uses | API behavior | Drift? |
|----------|----------|-------------|--------------|--------|
| Old TF + FF ON | ON | `notification` | Returns `notification` | NO |
| Old TF + FF OFF | OFF | `notification` | Returns `notification` | NO |
| New TF + FF ON | ON | `notifications` | Returns `notifications` | NO |
| New TF + FF OFF | OFF | `notifications` | Returns `notification` (1st only) | **YES** |
| Pre-FF SS + FF ON | ON | `notification` | Returns `notification` | NO |

## Conclusion

- **Safe to release when FF is permanently ON**: No drift for any scenario
- **Risk if released before FF ON**: Users who use `notifications` (plural) will see perpetual drift
- **Backward compatible**: Existing users with singular `notification` configs are unaffected regardless of FF state
- **No data loss**: Notifications are preserved even with FF toggle — just hidden in GET response when OFF
