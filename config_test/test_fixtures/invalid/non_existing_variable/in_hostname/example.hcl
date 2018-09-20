host "service-a" {
  hostname = "service-a[1..5].${b.c3.d4}"
  alias = "a${a}"
}

var {
  a = "123"
}