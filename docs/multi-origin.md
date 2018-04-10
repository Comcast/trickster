# Using Multiple-Origins with a single Trickster instance

There are 3 ways to configure multi-origin support.

* HTTP Pathing
* HTTP URL Parameters
* DNS Aliasing

## Basic Usage

To utilize Multiple Origins, you must craft a Trickster configuration file to be read when Trickster starts up. There is a good example in [conf/example.conf](../conf/example.conf). The config file should be placed in `/etc/trickster/trickster.conf` unless you specify a different path when starting Trickster with the `-config` command line argument.

Each origin that your Trickster instance supports must be explicitly enumerated in the configuration file. Trickster does not support open proxying.

Each origin is identified by an origin moniker, provided in the configuration section header for the origin ([origins.MONIKER]). For path and urlparam multi-origin configurations, the Moniker can be simple words. For DNS Aliasing, the origin moniker must match an FQDN that resolves to your Trickster instance.

In all cases, if Trickster cannot identify a valid origin by the client-provided moniker, it will proxy the request to the default origin.

### Path and URL Param Configurations

In these modes, Trickster will use a single FQDN but still map to multiple upstream origins. This is the simplest setup and requires the least amount of work. The client will indicate which origin is desired in the URL Parameter or Path for the request.

Example Path/URLParam Multi-Origin Configuration:
```
[origins]

    # default origin
    [origins.default]
        origin_url = 'http://prometheus.example.com:9090'
        api_path = '/api/v1'
        default_step = 300
        ignore_no_cache_header = false
        max_value_age_secs = 86400

    # "foo" origin
    [origins.foo]
        origin_url = 'http://prometheus-foo.example.com:9090'
        api_path = '/api/v1'
        default_step = 300
        ignore_no_cache_header = false
        max_value_age_secs = 86400

    # "bar" origin
    [origins.bar]
        origin_url = 'http://prometheus-bar.example.com:9090'
        api_path = '/api/v1'
        default_step = 300
        ignore_no_cache_header = false
        max_value_age_secs = 86400
```

#### Using HTTP Path as the Multi-Origin Indicator

The client prefixes the Trickster request path with the origin moniker.

This is the recommended method for integrating multi-origin support into Grafana.

Example Client Request URLs:
* To Request from Origin `foo`: http://trickster.example.com:9090/foo/query?query=xxx

* To Request from Origin `bar`: http://trickster.example.com:9090/bar/query?query=xxx

* To Request from Origin `default` (Method 1, no Moniker): http://trickster.example.com:9090/query?query=xxx

* To Request from Origin `default` (Method 2, with Moniker): http://trickster.example.com:9090/default/query?query=xxx

* Configuring Grafana to request from origin `foo` via Trickster:

<img src="./images/grafana-path-origin.png" width=610 />

#### Using urlparam as the Multi-Origin Indicator

The client provides the origin moniker as an url param.

Example Client Request URLs:

* To Request from Origin `foo`: http://trickster.example.com:9090/query?origin=foo&query=xxx

* To Request from Origin `bar`: http://trickster.example.com:9090/query?origin=bar&query=xxx

### DNS Alias Configuration

In this mode, multiple DNS records point to a single Trickster instance. The FQDN used by the client to reach Trickster represents the Origin Moniker. Therefore, the entire FQDN must be part of the configuration section header. In this mode, the path and url params requested against Trickster are _not_ modified to specify an origin.

Example DNS-based Origin Configuration:
```
[origins]

    # default origin
    [origins.default]
        origin_url = 'http://prometheus.example.com:9090'
        api_path = '/api/v1'
        default_step = 300
        ignore_no_cache_header = false
        max_value_age_secs = 86400

    # "foo" origin
    [origins.trickster-foo.example.com]
        origin_url = 'http://prometheus-foo.example.com:9090'
        api_path = '/api/v1'
        default_step = 300
        ignore_no_cache_header = false
        max_value_age_secs = 86400

    # "bar" origin
    [origins.trickster-bar.example.com]
        origin_url = 'http://prometheus-bar.example.com:9090'
        api_path = '/api/v1'
        default_step = 300
        ignore_no_cache_header = false
        max_value_age_secs = 86400

```

Example Client Request URLs:
*  To Request from Origin `foo`: http://trickster-foo.example.com:9090/query?query=xxx

*  To Request from Origin `bar`: http://trickster-bar.example.com:9090/query?query=xxx

*  To Request from Origin `default`: http://trickster.example.com:9090/query?query=xxx
