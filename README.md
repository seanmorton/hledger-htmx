# hledger-htmx
Web server for viewing hledger data powered by [htmx](htmx.org) and [go templ](https://templ.guide/). It provides similar functionality to the official [hledger-web](https://hledger.org/1.27/hledger-web.html).

So far this is only tested against `hledger v1.27.1`, and it should to be run on the same machine that has hledger installed with your journal files. The server calls `hledger` natively and parses its output for the data.
