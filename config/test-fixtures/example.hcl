alias "service-a" {
  pattern = "service-a[1..5].example.com",
  template = "a%1"
  ssh_config_name = "service-a"
}

alias "service-b" {
  pattern = "service-b[1..2].example.com",
  template = "b%1"
  ssh_config = {
    identity_file = "b_id_rsa.pub"
    port = 22
  }
}

ssh_config "service-a" {
  identity_file = "a_id_rsa.pub"
  port = 22
}
