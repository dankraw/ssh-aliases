host "consul" {
  hostname = "consul[1..3].[dc1|dc2].example.com"
  alias = "consul{#1}-{#2}"
  config = {
    identity_file = "some_file.pem"
    user = "ubuntu"
  }
}

host "frontend" {
  hostname = "frontend[1..2].example.com"
  alias = "front{#1}"
  config = "global"
}
