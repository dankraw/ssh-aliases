host "def" {
  hostname = "servcice-def.example.com"
  alias = "def"
  config = "def_conf"
}

config "def_conf" {
  _import = "root"
}

config "root" {
    _import = "def_conf"
}