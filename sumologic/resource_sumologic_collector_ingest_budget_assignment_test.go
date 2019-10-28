package sumologic

// func TestAccSumologicCollectorIngestBudgetAssignment(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSumologicCollectorIngestBudgetAssignmentConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttrSet("sumologic_collector_ingest_budget_assignment.assignment", "id"),
// 					resource.TestCheckResourceAttrSet("sumologic_collector_ingest_budget_assignment.assignment", "collector_id"),
// 					resource.TestCheckResourceAttrSet("sumologic_collector_ingest_budget_assignment.assignment", "ingest_budget_id"),
// 				),
// 			},
// 		}})
// }

var testAccSumologicCollectorIngestBudgetAssignmentConfig = `
resource "sumologic_collector" "test" {
  name = "assignment"
}

resource "sumologic_ingest_budget" "test" {
  name = "assignment"
  field_value = "assignment"
  capacity_bytes = 2
}

resource "sumologic_collector_ingest_budget_assignment" "assignment" {
  collector_id = "${sumologic_collector.test.id}"
  ingest_budget_id = "${sumologic_ingest_budget.test.id}"
}
`
