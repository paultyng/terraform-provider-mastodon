resource "mastodon_domain_block" "domain_blocks" {
  # this resource is best used with a for_each
  for_each = toset(["nsfw.social", "artalley.social"])
  domain   = each.key
}
