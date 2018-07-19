# Brazilian Banks API

> Query names and codes of banks in Brazil.

## Example

#### Request

```
GET /?q=itau HTTP/1.1
Accept: */*
Host: brazilian-banks-api.herokuapp.com
```

#### Response

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Length: 330
Content-Type: application/json
Date: Tue, 12 Sep 2017 23:00:44 GMT

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

Optionally you can omit the `q` parameter to get a list of all banks.

## API

The endpoint is `https://brazilian-banks-api.herokuapp.com`.

<dl>
  <dt><code>GET /</code></dt>
  <dd>Returns the whole list.</dd>
  <dt><code>GET /?q=...</code></dt>
  <dd>Filter the list by partial match, ignoring accents and special characters.</dd>
</dl>

## Developer's Note

```js
document.write(JSON.stringify(Array.from(document.querySelectorAll('td>a[href^="AgenciasRegioes.asp?"]')).map((a) => { const site = a.parentElement.nextElementSibling.querySelector('a:not([href="http://"])');return {code: a.parentElement.previousElementSibling.textContent.trim(), name:a.textContent.trim(), url: (site ? site.href : '')}; })));
```

## Reference

1. [http://www.febraban.org.br/bancos.asp](http://www.febraban.org.br/bancos.asp)

## License

The MIT License &copy; Arthur Corenzan 2017
