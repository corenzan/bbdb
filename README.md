# BBDb

> Free open-source API for querying banks' numeric code in Brazil.

## About

Monetary transactions in Brazil, such as transfers and deposits, require a numeric code that identifies the financial institution that manages the destinated account. This numeric code is assigned to each member of the STR (Sistema de Transferência de Reservas) by the country's central bank, Banco Central do brasil

BBDb is am open-source service to provide this information up-to-date and free of charge for your applications to consume via an HTTP API.

It's the same API we use when developing our own applications.

## Usage

At a glance, here's an example of filling a `<select>` element with data from the API.

```html
<select></select>

<script>
  const select = document.querySelector("select");

  fetch("https://bbdb.crz.li/?compe=y")
    .then(resp => resp.json())
    .then(data => {
      data.forEach(entry => {
        const opt = document.createElement("option");
        opt.value = entry.code;
        opt.textContent = entry.name;
        select.add(opt);
      }
    });
</script>
```

The `compe` parameter **exclude** items with a blank `code` field. That's usually what you want.

## API

The public endpoint is `https://bbdb.crz.li`.

<dl>
  <dt><code>GET /</code></dt>
  <dd>Get all the records.</dd>
  <dt><code>GET /?q=...</code></dt>
  <dd>Filter the records by partial match. Normalizes special characters.</dd>
  <dt><code>GET /?compe=y|yes|t|true|1</code></dt>
  <dd>Exclude records with a blank <code>code</code> field.</dd>
</dl>

## Development

You'll need Go 1.14.2+. Clone and download dependencies.

```shell
$ go mod download
```

See [Makefile](Makefile) for build tasks.

The data is loaded from disk, and the expected format is CSV.

## Reference

- https://www.bcb.gov.br/estabilidadefinanceira/str

## License

The MIT License © 2014 Arthur Corenzan
