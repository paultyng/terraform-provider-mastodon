[![Tests](https://github.com/paultyng/terraform-provider-mastodon/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/paultyng/terraform-provider-mastodon/actions/workflows/test.yml)

# Mastodon Terraform Provider

## Documentation

You can browse documentation on the [Terraform provider registry](https://registry.terraform.io/providers/paultyng/mastodon/latest/docs).

## Contributing

There is a [Docker Compose](./docker-compose.yml) setup provided to run a local Mastodon instance for acceptance testing. This includes some [ruby scripts](./acctest.rb) which modify the instance once it starts up to add users and manage some of the settings.
