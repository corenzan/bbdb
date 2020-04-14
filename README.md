# BBDb

> Free open-source database of banks in Brazil.

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

The endpoint is `https://bbdb.crz.li`.

<dl>
  <dt><code>GET /</code></dt>
  <dd>Returns the whole list.</dd>
  <dt><code>GET /?q=...</code></dt>
  <dd>Filter the list by partial match, ignoring accents and special characters.</dd>
  <dt><code>GET /?compact=on|yes|true|1</code></dt>
  <dd>Exclude banks with blank codes.</dd>
</dl>

## Development

You'll need Go 1.14.2+. Clone and download dependencies.

```shell
$ go mod download
```

See [Makefile](Makefile) for build tasks.

### Generate Database File

The same binary of the server can generate a new `database.json` from the CSV Febraban provides.

```shell
$ ./bbdb -src file.csv
```

## Reference

1. [Febraban](https://portal.febraban.org.br/pagina/3164/12/pt-br/associados)

## License

The MIT License © 2014 Arthur Corenzan
