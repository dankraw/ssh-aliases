host "my-service" {
  hostname = "instance[1..2].my-service.example.com"
  alias = "myservice{#1}"
  config = {
    identity_file = "my_service.pem"
  }
}