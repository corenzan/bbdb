# BBDb

> Free open-source API for querying banks identification in Brazil.

## About

Monetary transactions in Brazil, such as transfers and deposits, require a numeric code that identifies the financial institution that manages the destined account. This numeric code is assigned to each member of the STR (Sistema de Transferência de Reservas) by the country's central bank, Banco Central do Brasil.

BBDb is an open-source web service that provides this information up-to-date and free of charge.

## Usage

At a glance, here's an example of filling a `<select>` element with options from the API.

```html
<select></select>

<script>
  const select = document.querySelector("select");

  fetch("https://bbdb.crz.li/?compe=y")
    .then((resp) => resp.json())
    .then((data) => {
      data.forEach((entry) => {
        const opt = document.createElement("option");
        opt.value = entry.code;
        opt.textContent = entry.name;
        select.add(opt);
      });
    });
</script>
```

The `compe` parameter **exclude** records with a blank `code` field.

## API

The public endpoint is `https://bbdb.crz.li`.

<dl>
  <dt><code>GET /</code></dt>
  <dd>Get all the records.</dd>
  <dt><code>GET /?q=...</code></dt>
  <dd>Filter the records by partial match. Normalizes special characters.</dd>
  <dt><code>GET /?compe=y|yes|t|true|1</code></dt>
  <dd>Only show banks that participate in the payment system, Compe.</dd>
</dl>

## Development

You'll need Go 1.22+. Clone the repository, then run:

```shell
go mod download
```

To download dependencies. Then:

```shell
go run .
```

To start the server.

The data is read from `data.csv`.

## Reference

- https://www.bcb.gov.br/estabilidadefinanceira/str

## License

The MIT License © 2014 Arthur Corenzan
