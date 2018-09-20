host "service-a" {
  hostname = "service-a[1..5]"
  alias = "a${a}"
  config = "ext"
}

config "ext" {
  user = "${b.c3.d4}"
}

var {
  a = "123"
}