host "service-a" {
  hostname = "service-a[1..5].example.com",
  alias = "a{#1}"
  config = "service-a"
}

host "service-b" {
  hostname = "service-b[1..2].example.com",
  alias = "b{#1}"
  config = {
    identity_file = "b_id_rsa.pem"
    port = 22
  }
}

config "service-a" {
  identity_file = "a_id_rsa.pem"
  port = 22
}
