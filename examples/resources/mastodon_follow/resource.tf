resource "mastodon_follow" "local_server" {
  account = "@person"
}

resource "mastodon_follow" "remote_server" {
  account = "@person@domain"
}
