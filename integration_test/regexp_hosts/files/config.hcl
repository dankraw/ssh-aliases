host "dc1-services" {
    hostname = "ab-([a-z]+\\d+)\\.([a-z-]+)\\-(prod|test).my.dc1.com"
    alias = "{#3}.{#2}-{#1}.dc1"
}

host "dc2-services" {
    hostname = "ab-([a-z]+\\d+)\\.([a-z-]+)\\-(prod|test).my.dc2.net"
    alias = "{#3}.{#2}-{#1}.dc2"
}