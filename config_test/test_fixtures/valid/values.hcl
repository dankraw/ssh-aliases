values {
    domain1 = "my.domain1.example.com"
    domain2 = "my.domain2.example.com"
    env {
        dev = "development"
        prod = "production"
        user = "deployment"
        more {
            test = "testing"
            key_name = "secret"
            number = 1001
        }
    }
    threshold = 123
    b_count = 2
    b_alias = "b{#1}"
}
