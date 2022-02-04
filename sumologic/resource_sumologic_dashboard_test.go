package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var firstLabelKey = "service"
var secondLabelKey = "env"
var topologyLabel = TopologyLabel{
	Data: map[string][]string{
		firstLabelKey:  {"collection-proxy"},
		secondLabelKey: {"dev", "prod"},
	},
}

func TestAccSumologicDashboard_basic(t *testing.T) {
	testNameSuffix := acctest.RandString(16)
	title := "terraform_test_dashboard_" + testNameSuffix

	var dashboard Dashboard
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDashboardDestroy(dashboard),
		Steps: []resource.TestStep{
			{
				Config: dashboardImportConfig(title),
			},
			{
				ResourceName:      "sumologic_dashboard.tf_import_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicDashboard_create(t *testing.T) {
	testNameSuffix := acctest.RandString(16)

	// create config
	title := "terraform_test_dashboard_" + testNameSuffix
	description := "Test dashboard description"
	theme := "Dark"
	refreshInterval := 120

	literalRangeName := "today"
	textPanel := TextPanel{
		Key:   "text-panel-001",
		Title: "Text Panel Title",
		Text:  "This is a text panel",
	}
	layout := GridLayout{
		LayoutStructures: []LayoutStructure{
			{
				Key:       "text-panel-001",
				Structure: "{\"height\":15,\"width\":15,\"x\":0,\"y\":0}",
			},
		},
	}
	variable := Variable{
		Name:         "_sourceHost",
		DisplayName:  "Source Host",
		DefaultValue: "host-1",
		SourceDefinition: CsvVariableSourceDefinition{
			Values: "host-1,host-2,host-3",
		},
		AllowMultiSelect: false,
		IncludeAllOption: true,
		HideFromUI:       false,
	}

	var dashboard Dashboard
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDashboardDestroy(dashboard),
		Steps: []resource.TestStep{
			{
				Config: dashboardCreateConfig(title, description, theme, refreshInterval,
					topologyLabel, literalRangeName, textPanel, layout, variable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDashboardExists("sumologic_dashboard.tf_crud_test", &dashboard, t),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"title", title),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"description", description),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"refresh_interval", strconv.FormatInt(int64(refreshInterval), 10)),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"theme", theme),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"topology_label_map.0.data.#", "2"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"time_range.#", "1"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.#", "1"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.0.text_panel.0.key", textPanel.Key),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.0.text_panel.0.text", textPanel.Text),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"layout.0.grid.0.layout_structure.0.key", layout.LayoutStructures[0].Key),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.0.name", variable.Name),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.0.source_definition.0.csv_variable_source_definition.0.values",
						variable.SourceDefinition.(CsvVariableSourceDefinition).Values),
				),
			},
		},
	})
}

func TestAccSumologicDashboard_update(t *testing.T) {
	testNameSuffix := acctest.RandString(16)

	// create config
	title := "terraform_test_dashboard_" + testNameSuffix
	description := "Test dashboard description"
	theme := "Dark"
	refreshInterval := 120
	literalRangeName := "today"
	textPanel := TextPanel{
		Key:   "text-panel-001",
		Title: "Text Panel Title",
		Text:  "This is a text panel",
	}
	layout := GridLayout{
		LayoutStructures: []LayoutStructure{
			{
				Key:       "text-panel-001",
				Structure: "{\"height\":15,\"width\":15,\"x\":0,\"y\":0}",
			},
		},
	}
	csvVariable := Variable{
		Name:         "_sourceHost",
		DisplayName:  "Source Host",
		DefaultValue: "host-1",
		SourceDefinition: CsvVariableSourceDefinition{
			Values: "host-1,host-2,host-3",
		},
		AllowMultiSelect: false,
		IncludeAllOption: true,
		HideFromUI:       false,
	}

	// updated config
	newTheme := "Light"
	newRefreshInterval := 300
	newFirstLabelValue := "collection-cluster"
	newLiteralRangeName := "week"
	searchPanel := SumoSearchPanel{
		Key:   "search-panel-001",
		Title: "API Errors",
		Queries: []SearchPanelQuery{
			{
				QueryString: "_sourceCategory=api error | timeslice 1h | count by _timeslice",
				QueryType:   "Log",
				QueryKey:    "search-query-01",
			},
		},
		Description: "API errors per hour",
		TimeRange: CompleteLiteralTimeRange{
			RangeName: "today",
		},
	}
	newLayout := GridLayout{
		LayoutStructures: []LayoutStructure{
			{
				Key:       "text-panel-001",
				Structure: "{\"height\":10,\"width\":15,\"x\":0,\"y\":0}",
			},
			{
				Key:       "search-panel-001",
				Structure: "{\"height\":10,\"width\":15,\"x\":0,\"y\":10}",
			},
		},
	}
	newVariables := []Variable{
		csvVariable,
		{
			Name:        "_remoteModule",
			DisplayName: "Remote Module",
			SourceDefinition: LogQueryVariableSourceDefinition{
				Query: "_sourceCategory=api error | parse '[module=*]' as module | count by module",
				Field: "module",
			},
			AllowMultiSelect: true,
			IncludeAllOption: true,
			HideFromUI:       false,
		},
	}

	var dashboard Dashboard
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDashboardDestroy(dashboard),
		Steps: []resource.TestStep{
			{
				Config: dashboardCreateConfig(title, description, theme, refreshInterval,
					topologyLabel, literalRangeName, textPanel, layout, csvVariable),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDashboardExists("sumologic_dashboard.tf_crud_test", &dashboard, t),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"title", title),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"refresh_interval", strconv.FormatInt(int64(refreshInterval), 10)),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"theme", theme),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"topology_label_map.0.data.#", "2"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						literalRangeName),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.#", "1"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.0.text_panel.0.key", textPanel.Key),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"layout.0.grid.0.layout_structure.#", "1"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"layout.0.grid.0.layout_structure.0.key", layout.LayoutStructures[0].Key),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.#", "1"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.0.name", csvVariable.Name),
				),
			},
			{
				Config: dashboardUpdateConfig(title, description, newTheme, newRefreshInterval,
					firstLabelKey, newFirstLabelValue, newLiteralRangeName, textPanel, searchPanel,
					newLayout, newVariables),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDashboardExists("sumologic_dashboard.tf_crud_test", &dashboard, t),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"title", title),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"refresh_interval", strconv.FormatInt(int64(newRefreshInterval), 10)),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"theme", newTheme),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"topology_label_map.0.data.#", "1"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"topology_label_map.0.data.0.label", firstLabelKey),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"topology_label_map.0.data.0.values.0", newFirstLabelValue),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"time_range.0.begin_bounded_time_range.0.from.0.literal_time_range.0.range_name",
						newLiteralRangeName),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.#", "2"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.0.text_panel.0.key", textPanel.Key),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.1.sumo_search_panel.0.key", searchPanel.Key),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.1.sumo_search_panel.0.description", searchPanel.Description),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"panel.1.sumo_search_panel.0.query.0.query_key", searchPanel.Queries[0].QueryKey),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"layout.0.grid.0.layout_structure.#", "2"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"layout.0.grid.0.layout_structure.0.key", newLayout.LayoutStructures[0].Key),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"layout.0.grid.0.layout_structure.1.key", newLayout.LayoutStructures[1].Key),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.#", "2"),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.0.name", newVariables[0].Name),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.0.source_definition.0.csv_variable_source_definition.0.values",
						newVariables[0].SourceDefinition.(CsvVariableSourceDefinition).Values),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.1.name", newVariables[1].Name),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.1.source_definition.0.log_query_variable_source_definition.0.query",
						newVariables[1].SourceDefinition.(LogQueryVariableSourceDefinition).Query),
					resource.TestCheckResourceAttr("sumologic_dashboard.tf_crud_test",
						"variable.1.source_definition.0.log_query_variable_source_definition.0.field",
						newVariables[1].SourceDefinition.(LogQueryVariableSourceDefinition).Field),
				),
			},
		},
	})
}

func testAccCheckDashboardDestroy(dashboard Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetDashboard(dashboard.ID)
		if err == nil {
			return fmt.Errorf("Dashboard (id=%s) still exists", dashboard.ID)
		}
		return nil
	}
}

func testAccCheckDashboardExists(name string, dashboard *Dashboard, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Error = %s. Dashboard not found: %s", strconv.FormatBool(ok), name)
		}

		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("Dashboard ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newDashboard, err := client.GetDashboard(id)
		if err != nil {
			return fmt.Errorf("Dashboard (id=%s) not found", id)
		}
		dashboard = newDashboard
		return nil
	}
}

func dashboardImportConfig(title string) string {
	return fmt.Sprintf(`
		data "sumologic_personal_folder" "personalFolder" {}
		resource "sumologic_dashboard" "tf_import_test" {
			title = "%s"
			description = "Test dashboard description"
			folder_id = data.sumologic_personal_folder.personalFolder.id
			refresh_interval = 120
			theme = "Light"
			time_range {
				begin_bounded_time_range {
					from {
						epoch_time_range {
							epoch_millis = 1612137600
						}
					}
					to {
						epoch_time_range {
							epoch_millis = 1612223999
						}
					}
				}
			}
			panel {
				text_panel {
					key = "tf-text-panel-001"
					title = "What does wsb say?"
					visual_settings = "{\"general\":{\"type\":\"column\"}}"
					keep_visual_settings_consistent_with_parent = true
					text = "Buy AMC!"
				}
			}
			layout {
				grid {
					layout_structure {
						key = "tf-text-panel-001"
						structure = "{\"height\":10,\"width\":15,\"x\":0,\"y\":12}"
					}
				}
			}
			variable {
				name = "idle_cpu"
				display_name = "Idle CPU"
				source_definition {
					metadata_variable_source_definition {
						filter = "_sourceHost=api-* metric=CPU_Idle"
						key = "deployment"
					}
				}
				allow_multi_select = false
				include_all_option = true
				hide_from_ui = false
			}
			coloring_rule {
				scope = "CPU_*"
				single_series_aggregate_function = "Average"
				multiple_series_aggregate_function = "Average"
				color_threshold {
					color = "FFFFFF"
					min = 1
					max = 50
				}
			}
		}`,
		title,
	)
}

func dashboardCreateConfig(title string, description string, theme string, refreshInterval int,
	topologyLabel TopologyLabel, rangeName string, textPanel TextPanel, layout GridLayout,
	variable Variable) string {

	return fmt.Sprintf(`
		data "sumologic_personal_folder" "personalFolder" {}
		resource "sumologic_dashboard" "tf_crud_test" {
			title = "%s"
			description = "%s"
			folder_id = data.sumologic_personal_folder.personalFolder.id
			refresh_interval = %d
			theme = "%s"
			topology_label_map {
				data {
					label = "%s"
					values = ["%s"]
				}
				data {
					label = "%s"
					values = ["%s", "%s"]
				}
			}
			time_range {
				begin_bounded_time_range {
					from {
						literal_time_range {
							range_name = "%s"
						}
					}
				}
			}
			panel {
				text_panel {
					key = "%s"
					title = "%s"
					visual_settings = "{\"general\":{\"type\":\"column\"}}"
					keep_visual_settings_consistent_with_parent = true
					text = "%s"
				}
			}
			layout {
				grid {
					layout_structure {
						key = "%s"
						structure = "{\"height\":10,\"width\":15,\"x\":0,\"y\":0}"
					}
				}
			}
			variable {
				name = "%s"
				display_name = "%s"
				default_value = "%s"
				source_definition {
					csv_variable_source_definition {
						values = "%s"
					}
				}
				allow_multi_select = false
				include_all_option = true
				hide_from_ui = false
			}
		}`,
		title, description, refreshInterval, theme, firstLabelKey, topologyLabel.Data[firstLabelKey][0],
		secondLabelKey, topologyLabel.Data[secondLabelKey][0], topologyLabel.Data[secondLabelKey][1],
		rangeName, textPanel.Key, textPanel.Title, textPanel.Text,
		layout.LayoutStructures[0].Key, variable.Name, variable.DisplayName, variable.DefaultValue,
		variable.SourceDefinition.(CsvVariableSourceDefinition).Values,
	)
}

func dashboardUpdateConfig(title string, description string, theme string, refreshInterval int,
	topologyLabel string, topologyLabelValue string, rangeName string, textPanel TextPanel,
	searchPanel SumoSearchPanel, layout GridLayout, variables []Variable) string {

	loqQuerySourceDef := variables[1].SourceDefinition.(LogQueryVariableSourceDefinition)
	csvSourceDef := variables[0].SourceDefinition.(CsvVariableSourceDefinition)

	return fmt.Sprintf(`
		data "sumologic_personal_folder" "personalFolder" {}
		resource "sumologic_dashboard" "tf_crud_test" {
			title = "%s"
			description = "%s"
			folder_id = data.sumologic_personal_folder.personalFolder.id
			refresh_interval = %d
			theme = "%s"
			topology_label_map {
				data {
					label = "%s"
					values = ["%s"]
				}
			}
			time_range {
				begin_bounded_time_range {
					from {
						literal_time_range {
							range_name = "%s"
						}
					}
				}
			}
			panel {
				text_panel {
					key = "%s"
					title = "%s"
					visual_settings = "{\"general\":{\"type\":\"column\"}}"
					keep_visual_settings_consistent_with_parent = true
					text = "%s"
				}
			}

			panel {
				sumo_search_panel {
					key = "%s"
					title = "%s"
					visual_settings = "{\"general\":{\"type\":\"column\"}}"
					keep_visual_settings_consistent_with_parent = true
					description = "%s"
					query {
						query_string = "%s"
						query_type = "Logs"
						query_key = "%s"
					}
					time_range {
						begin_bounded_time_range {
							from {
								relative_time_range {
									relative_time = "-1d"
								}
							}
						}
					}
				}
			}

			layout {
				grid {
					layout_structure {
						key = "%s"
						structure = "{\"height\":10,\"width\":15,\"x\":0,\"y\":0}"
					}
					layout_structure {
						key = "%s"
						structure = "{\"height\":10,\"width\":15,\"x\":0,\"y\":10}"
					}
				}
			}
			variable {
				name = "%s"
				display_name = "%s"
				default_value = "%s"
				source_definition {
					csv_variable_source_definition {
						values = "%s"
					}
				}
				allow_multi_select = false
				include_all_option = true
				hide_from_ui = false
			}
			variable {
				name = "%s"
				display_name = "%s"
				source_definition {
					log_query_variable_source_definition {
						query = "%s"
						field = "%s"
					}
				}
				allow_multi_select = false
				include_all_option = true
				hide_from_ui = false
			}
		}`,
		title, description, refreshInterval, theme, topologyLabel, topologyLabelValue, rangeName,
		textPanel.Key, textPanel.Title, textPanel.Text,
		searchPanel.Key, searchPanel.Title, searchPanel.Description, searchPanel.Queries[0].QueryString,
		searchPanel.Queries[0].QueryKey,
		layout.LayoutStructures[0].Key, layout.LayoutStructures[1].Key,
		variables[0].Name, variables[0].DisplayName, variables[0].DefaultValue, csvSourceDef.Values,
		variables[1].Name, variables[1].DisplayName, loqQuerySourceDef.Query, loqQuerySourceDef.Field,
	)
}
