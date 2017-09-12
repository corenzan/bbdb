# Brazilian Banks API

> Query names and codes of banks in Brazil.

## Example

#### Request

```
GET /?q=itau HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: brazilian-banks-api.herokuapp.com
User-Agent: HTTPie/0.9.9
```

#### Response

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Connection: keep-alive
Content-Length: 330
Content-Type: application/json
Date: Tue, 12 Sep 2017 23:00:44 GMT
Server: Cowboy
Via: 1.1 vegur
X-Response-Time: 39.187ms

[
    {
        "code": "184",
        "name": "Banco Itaú BBA S.A."
    },
    {
        "code": "479",
        "name": "Banco ItaúBank S.A"
    },
    {
        "code": "",
        "name": "Banco Itaucard S.A."
    },
    {
        "code": "M09",
        "name": "Banco Itaucred Financiamentos S.A."
    },
    {
        "code": "",
        "name": "Banco ITAULEASING S.A."
    },
    {
        "code": "652",
        "name": "Itaú Unibanco Holding S.A."
    },
    {
        "code": "341",
        "name": "Itaú Unibanco S.A."
    }
]
```


## API

The endpoint is `https://brazilian-banks-api.herokuapp.com`.

<dl>
  <dt><code>GET /</code></dt>
  <dd>Returns the whole list.</dd>
  <dt><code>GET /?q=...</code></dt>
  <dd>Filter the list by text match, ignoring special characters and accents.</dd>
</dl>

## Reference

1. [http://www.febraban.org.br/bancos.asp](http://www.febraban.org.br/bancos.asp)

## License

The MIT License &copy; Arthur Corenzan 2017
