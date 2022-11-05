provider "mastodon" {
  server = var.mastodon_server

  # these are created when you add an application in the Mastodon web UI
  client_id     = var.mastodon_client_id
  client_secret = var.mastodon_client_secret

  username = var.mastodon_username
  password = var.mastodon_password
}
