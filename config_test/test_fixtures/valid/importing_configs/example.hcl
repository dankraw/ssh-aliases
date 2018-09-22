host "abc" {
  hostname = "servcice-abc.example.com"
  alias = "abc"
  config {
    _import = "root"
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
  _import = "intermediate"
}

config "intermediate" {
    _import = ["root", "root_2"]
    this = "happens"
}

config "root" {
    additional = "extension"
    another = "one"
}

config "root_2" {
    additional = "extension 2"
    another = "two"
    _import = "root" # not a circular dependency in this case
}