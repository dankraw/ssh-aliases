alias "service-c" {
  pattern = "service-c[1..2].example.com",
  template = "c%1"
  ssh_config = {
    identity_file = "c_id_rsa.pub"
  }
}