# BBDb

> Free, open-source API for querying banks in Brazil.

## Example

### Request

```
curl https://bbdb.crz.li/?q=itau&compact=yes
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

The `compact` parameter exclude items with a blank `code` field. Optionally you can omit the `q` parameter to get a full list of all banks.

## API

The endpoint is `https://bbdb.crz.li`.

<dl>
  <dt><code>GET /</code></dt>
  <dd>Returns the whole list.</dd>
  <dt><code>GET /?q=...</code></dt>
  <dd>Filter the list by partial match. Normalizes special characters.</dd>
  <dt><code>GET /?compact=on|yes|true|1</code></dt>
  <dd>Exclude banks with blank codes.</dd>
</dl>

## Development

You'll need Go 1.14.2+. Clone and download dependencies.

```shell
$ go mod download
```

See [Makefile](Makefile) for build tasks.

The data is loaded from disk. The expected format is a CSV using `;` as separator.

## Reference

1. [Febraban](https://portal.febraban.org.br/pagina/3164/12/pt-br/associados)

## License

The MIT License © 2014 Arthur Corenzan
