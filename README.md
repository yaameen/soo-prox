SooProx
=======

A simple proxy server

Installation
------------

To install SooProx, run the following command:

goCopy code

`go get github.com/yaamynu/sooprox`

Usage
-----

SooProx can be run in two modes:

1.  CLI mode
2.  Configuration file mode

In CLI mode, SooProx can be run with the following flags:

-   `-c, --config`: Load configuration from `FILE`. Default is `config.yaml`.
-   `-H, --host`: Listen on `HOST`.
-   `-p, --port`: Listen on `PORT`.
-   `-P, --proxies`: Proxy `PREFIX::TARGET`.
-   `-s, --secure`: Use TLS.

In Configuration file mode, SooProx can be run with a YAML configuration file. The file should contain the following fields:

-   `host`: Listen on `HOST`.
-   `port`: Listen on `PORT`.
-   `proxies`: Array of proxies in the format `PREFIX::TARGET`.
-   `secure`: Use TLS.

Trusting the CA certificate
---------------------------

If running in secure mode, the CA certificate must be trusted. To trust the CA certificate, run the following command:

Copy code

`sooprox ca-trust`

Version
-------

`v0.1`

Author
------

Yameen Mohamed (<yaamynu@gmail.com>)