package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFollowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFollowResourceConfig("@acctestadmin"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mastodon_follow.test", "account", "@acctestadmin"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "mastodon_follow.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "@acctestadmin",
			},
			// Update and Read testing (this should recreate, so new id)
			{
				Config: testAccFollowResourceConfig("@acctest2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mastodon_follow.test", "account", "@acctest2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TODO: test another server, @tyng@hachyderm.io

func testAccFollowResourceConfig(id string) string {
	return fmt.Sprintf(providerConfig+`
resource "mastodon_follow" "test" {
	account = %[1]q
}
`, id)
}
