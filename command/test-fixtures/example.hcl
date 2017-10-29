alias "service-a" {
  pattern = "service-a[1..2].example.com",
  template = "a%1"
  ssh_config_name = "global"
}

alias "service-b" {
  pattern = "service-b[1..2].example.com",
  template = "b%1"
  ssh_config_name = "global"
}

ssh_config "global" {
  identity_file = "id_rsa.pub"
  port = 22
}
