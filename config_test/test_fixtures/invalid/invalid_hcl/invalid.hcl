host "service-a" {
  hostname = "service-a[1..5].example.com"
  alias = "a{#1}"
  config = {
    user = "joe"
  }
