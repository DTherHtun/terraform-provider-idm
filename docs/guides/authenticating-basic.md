---
page_title: "Authenticating with REDHAT IDM"
---

```
provider "idm" {
    idm_server = # Redhat IDM server address
    user = # User account for login
    password = # password for login
    insecure = # enable or disable tls (value = true or false)
}
```
