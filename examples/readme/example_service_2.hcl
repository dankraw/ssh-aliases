host "other" {
  hostname = "other[1..2].example.com"
  alias = "other{#1}"
  config {
    user = "lurker"
    identity_file = "${keys.other}"
    port = 22
  }
}