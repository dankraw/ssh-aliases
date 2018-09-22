host "abc" {
  hostname = "servcice-abc.example.com"
  alias = "abc"
  config {
    _extend = "root"
    x = "y"
  }
}

host "def" {
  hostname = "servcice-def.example.com"
  alias = "def"
  config = "def_conf"
}

config "def_conf" {
  some_prop = 123
  this = "never happens"
  _extend = "intermediate"
}

config "intermediate" {
    _extend = ["root", "root_2"]
    this = "happens"
}

config "root" {
    additional = "extension"
    another = "one"
}

config "root_2" {
    additional = "extension 2"
    another = "two"
    _extend = "root" # not a circular dependency in this case
}