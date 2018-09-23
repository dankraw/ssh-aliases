host "abc" {
  hostname = "node[1..2].abc.[dev|test].example.com"
  alias = "{#2}.abc{#1}"
  config = "abc-config"
}

config "abc-config" {
  user = "ubuntu"
  identity_file = "${keys.abc}"
  port = 22
}