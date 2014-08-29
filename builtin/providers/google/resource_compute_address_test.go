package google

import (
	"fmt"
	"testing"

	"code.google.com/p/google-api-go-client/compute/v1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeAddress_basic(t *testing.T) {
	var addr compute.Address

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeAddress_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeAddressExists(
						"google_compute_address.foobar", &addr),
				),
			},
		},
	})
}

func testAccCheckComputeAddressDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.Resources {
		if rs.Type != "google_compute_address" {
			continue
		}

		_, err := config.clientCompute.Addresses.Get(
			config.Project, config.Region, rs.ID).Do()
		if err == nil {
			return fmt.Errorf("Address still exists")
		}
	}

	return nil
}

func testAccCheckComputeAddressExists(n string, addr *compute.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.Addresses.Get(
			config.Project, config.Region, rs.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.ID {
			return fmt.Errorf("Addr not found")
		}

		*addr = *found

		return nil
	}
}

const testAccComputeAddress_basic = `
resource "google_compute_address" "foobar" {
	name = "terraform-test"
}`