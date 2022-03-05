# elasticsearch-init

An init container meant to be spun up alongside a new Elasticsearch instance to set up the required items for Honeystats.

Sets up:
- ingest pipelines
  - geoip
- dashboards
  - one for each service
  - general dashboard for source IPs etc.
