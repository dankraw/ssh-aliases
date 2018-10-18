host "dc1-services" {
  hostname = "instance(\\d+)\\.my\\-service\\-(dev|prod|test)\\..+"
  alias = "host{#1}.{#2}"
  config {
    user = "abc"
    identity_file = "~/.ssh/key.pem"
  }
}