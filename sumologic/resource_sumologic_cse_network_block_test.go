package sumologic

import (
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicSCENetworkBlock_create(t *testing.T) {
	SkipCseTest(t)

	var networkBlock CSENetworkBlock
	nAddressBlock := generateRandomCIDRBlock()
	nLabel := "network block test"
	nInternal := true
	nSuppressesSignals := false
	resourceName := "sumologic_cse_network_block.network_block"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSENetworkBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSENetworkBlockConfig(nAddressBlock, nLabel, nInternal, nSuppressesSignals),
				Check: resource.ComposeTestCheckFunc(
					testCheckNetworkBlockExists(resourceName, &networkBlock),
					testCheckNetworkBlockValues(&networkBlock, nAddressBlock, nLabel, nInternal, nSuppressesSignals),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSENetworkBlockDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_network_block" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Network Block destruction check: CSE Network Block ID is not set")
		}

		networkBlockID := rs.Primary.Attributes["id"]

		s, err := client.GetCSENetworkBlock(networkBlockID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("network Block still exists")
		}
	}
	return nil
}

func testCreateCSENetworkBlockConfig(nAddressBlock string, nLabel string, nInternal bool, nSuppressesSignals bool) string {
	return fmt.Sprintf(`
resource "sumologic_cse_network_block" "network_block" {
	address_block = "%s"
	label = "%s"
	internal = "%t"
	suppresses_signals = "%t"
}
`, nAddressBlock, nLabel, nInternal, nSuppressesSignals)
}

func testCheckNetworkBlockExists(n string, networkBlock *CSENetworkBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("network Block ID is not set")
		}

		networkBlockID := rs.Primary.Attributes["id"]

		c := testAccProvider.Meta().(*Client)
		networkBlockResp, err := c.GetCSENetworkBlock(networkBlockID)
		if err != nil {
			return err
		}

		*networkBlock = *networkBlockResp

		return nil
	}
}

func testCheckNetworkBlockValues(networkBlock *CSENetworkBlock, nAddressBlock string, nLabel string, nInternal bool, nSuppressesSignals bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if networkBlock.AddressBlock != nAddressBlock {
			return fmt.Errorf("bad address block, expected \"%s\", got: %#v", nAddressBlock, networkBlock.AddressBlock)
		}
		if networkBlock.Label != nLabel {
			return fmt.Errorf("bad label, expected \"%s\", got: %#v", nLabel, networkBlock.Label)
		}
		if networkBlock.Internal != nInternal {
			return fmt.Errorf("bad internal flag, expected \"%t\", got: %#v", nInternal, networkBlock.Internal)
		}
		if networkBlock.SuppressesSignals != nSuppressesSignals {
			return fmt.Errorf("bad suppressesSignals flag, expected \"%t\", got: %#v", nSuppressesSignals, networkBlock.SuppressesSignals)
		}
		return nil
	}
}

func generateRandomCIDRBlock() string {
	ip := make(net.IP, 4)
	for i := 0; i < 3; i++ {
		ip[i] = byte(rand.Intn(256))
	}
	ip[3] = 0

	return ip.String() + "/26"
}
