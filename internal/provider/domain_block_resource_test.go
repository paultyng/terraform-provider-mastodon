package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDomainBlockResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDomainBlockResourceConfig("nsfw.social"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mastodon_domain_block.test", "domain", "nsfw.social"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "mastodon_domain_block.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "nsfw.social",
			},
			// Update and Read testing (this should recreate, so new id)
			{
				Config: testAccDomainBlockResourceConfig("artalley.social"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mastodon_domain_block.test", "domain", "artalley.social"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDomainBlockResourceConfig(id string) string {
	return fmt.Sprintf(providerConfig+`
resource "mastodon_domain_block" "test" {
	domain = %[1]q
}
`, id)
}
