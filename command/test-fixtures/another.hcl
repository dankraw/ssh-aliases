host "service-c" {
  hostname = "service-c[1..2].example.com",
  alias = "c%1"
  config = {
    identity_file = "c_id_rsa.pub"
  }
}