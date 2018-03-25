host "service-a" {
  hostname = "service-a[1..5].${domain1}"
  alias = "a{#1}"
  config = "service-a"
}

host "service-b" {
  hostname = "service-b[1..${b_count}].example.com"
  alias = "${b_alias}"
  config {
    identity_file = "b_id_${env.more.number}_rsa.pem"
    port = 22
  }
}

config "service-a" {
  identity_file = "a_${env.more.number}_id_${env.more.key_name}_rsa.pem"
  port = 22
  user = "${env.user}"
}
