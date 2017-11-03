host "service-a" {
  hostname = "service-a[1..2].example.com",
  alias = "a%1"
  config = "global"
}

host "service-b" {
  hostname = "service-b[1..2].example.com",
  alias = "b%1"
  config = "global"
}

config "global" {
  identity_file = "id_rsa.pub"
  port = 22
}
