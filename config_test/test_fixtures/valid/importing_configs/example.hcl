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
  _import = "intermediate"
}

config "intermediate" {
    _import = "root"
    this = "happens"
}

config "root" {
    additional = "extension"
    another = "one"
}