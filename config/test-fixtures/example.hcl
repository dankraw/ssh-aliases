# here we have two aliases using the same ssh configuration
alias "hermes-frontend" {
  patterns = ["host[1..5].example.com"]
  regexp = "(host\\d+)"
  template = "%1"
  ssh_config_name = "private"
}

alias "hermes-consumers" {
  patterns = ["host[1..3].example.com"]
  regexp = "(host\\d+)"
  template = "%1"
  ssh_config_name = "private"
}

ssh_config "private" {
  identity_file = "id_rsa.pub"
  port = 22
}
