# Brazilian Banks API

> Query banks names and codes in Brazil.

## Example

### Request

```
GET /?q=itau&compact=yes HTTP/1.1
Host: brazilian-banks-api.herokuapp.com
Accept: */*
```

### Response

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: : *
Content-Type: application/json; charset=utf-8
Etag: 0ee55f76a1ec068837b939821e96d3de
Content-Length: 382

[
  {
    "name": "Banco Itaú BBA S.A.",
    "code": "184",
    "url": "http://www.itaubba.com.br/"
  },
  { "name": "Banco Itaú Consignado S.A.", "code": "029", "url": "" },
  {
    "name": "Banco ItauBank S.A",
    "code": "479",
    "url": "http://www.itaubank.com.br/"
  },
  {
    "name": "Itaú Unibanco Holding S.A.",
    "code": "652",
    "url": "http://www.itau.com.br/"
  },
  {
    "name": "Itaú Unibanco S.A.",
    "code": "341",
    "url": "http://www.itau.com.br/"
  }
]
```

Optionally you can omit the `q` parameter to get a full list of all banks.

## API

The endpoint is `https://brazilian-banks-api.herokuapp.com`.

<dl>
  <dt><code>GET /</code></dt>
  <dd>Returns the whole list.</dd>
  <dt><code>GET /?q=...</code></dt>
  <dd>Filter the list by partial match, ignoring accents and special characters.</dd>
  <dt><code>GET /?compact=on|yes|true|1</code></dt>
  <dd>Exclude banks with blank codes.</dd>
</dl>

## Reference

1. [http://www.febraban.org.br/bancos.asp](http://www.febraban.org.br/bancos.asp)

## License

The MIT License &copy; Arthur Corenzan 2017
