host "all" {
    alias = "*"
    config {
        a = 1
    }
}

host "prod" {
    hostname = "prod*"
    config {
        a = 2
    }
}