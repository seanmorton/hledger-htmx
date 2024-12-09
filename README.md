# hledger-htmx
A web server for viewing [hledger](https://hledger.org/) data powered by [htmx](https://htmx.org) and [go templ](https://templ.guide/). It provides similar functionality to [hledger-web](https://hledger.org/1.27/hledger-web.html), with the ability to:
* Visually dive into subaccounts and their transactions
* Easily view data across time periods
* Monitor a monthly budget

Warning: this was written quickly for personal use and contains some hacks! So far this is tested against `hledger v1.27.1` and it should be run on the same machine that has hledger installed with your journal files (the server just calls `hledger` natively and parses its output.)



https://github.com/user-attachments/assets/b9946e5f-b02d-4bbb-9a29-97ae353e5018

