# Xtream-API

This is a Xtream API writing in Golang that can be used to fetch IPTV Streams & VOD Content. This
tool will generate `.strm` files for all Livestream, Movie, and Series content offered by the IPTV
service. These files can be used to import into Jellyfin or Plex in order to create your own IPTV
Player.

## Environment Variables

| **Variables**           | **Description**                                                         | **Default** |
| ----------------------- | ----------------------------------------------------------------------- | ----------- |
| `XTREAM_URL`            | The URL used for fetching IPTV streams.                                 | `""`        |
| `XTREAM_USERNAME`       | The Username used to log in to the IPTV service.                        | `""`        |
| `XTREAM_PASSWORD`       | The Password used to log in to the IPTV service.                        | `""`        |
| `XTREAM_IMAGES`         | Whether the tool should download images alongside stream file.          | `false`     |
| `XTREAM_NFO`            | Whether the tool should download metadata alongside stream file.        | `false`     |
| `XTREAM_REQUESTS`       | The maximum number of requests per minute.                              | `1000`      |
| `XTREAM_TIMEOUT`        | The maximum time before a request returns a timeout error, in seconds.  | `30`        |
| `XTREAM_BANNED_LIVE`    | A ',' separated list of banned livestream (partial) **category** names. | `[]`        |
| `XTREAM_BANNED_MOVIES`  | A ',' separated list of banned movie (partial) **category** names.      | `[]`        |
| `XTREAM_BANNED_SERIES`  | A ',' separated list of banned series (partial) **category** names.     | `[]`        |

## Docker Compose

Here is an example of how to use this tool as a docker container, using a docker compose file:

```docker-compose
---
services:
  xtream-api:
    container_name: xtream-api
    image: captainkills/xtream-api:latest
    environment:
      - XTREAM_URL=...                    # Change to your IPTV URL
      - XTREAM_USERNAME=...               # Change to your IPTV Username
      - XTREAM_PASSWORD=...               # Change to your IPTV Password

      # Optional Environment Variables
      # - XTREAM_IMAGES=true              # Enable Image Fetching
      # - XTREAM_NFO=true                 # Enable Metadata Fetching
      # - XTREAM_REQUESTS=1000            # Maximum requests per minute
      # - XTREAM_TIMEOUT=30               # Maximum time before request returns timeout error
      # - XTREAM_BANNED_LIVE=""           # ',' Seperated list of banned livestream categories
      # - XTREAM_BANNED_MOVIES=""         # ',' Seperated list of banned movie categories
      # - XTREAM_BANNED_SERIES=""         # ',' Seperated list of banned series categories
    volumes:
      - ./media:/media                    # Change to your desired media directory on Host
      - ./cache:/cache                    # Change to your desired cache directory on Host
    restart: unless-stopped
```

Here is an example of how to use this tool as a docker container, using the `docker run` command:

```bash
docker run -d \
  --name xtream-api \
  --restart unless-stopped \
  -e XTREAM_URL="..." \
  -e XTREAM_USERNAME="..." \
  -e XTREAM_PASSWORD="..." \
  -v ./media:/media \
  -v ./cache:/cache \
  captainkills/xtream-api:latest
```

## Functionality

The content fetched by this tool will be downloaded in the `/media` directory:
* `/media/live` directory for Livestream content.
* `/media/movies` for Movie content.
* `/media/series` for Series content.
* `/cache` for caching successfully updated streams for future runs.

### Images & Metadata

If downloading images is enabled, each piece of content will try to download its `cover.jpg` image.
Image downloading can be enabled using the `XTREAM_IMAGES` environment variable.

If downloading metadata is enabled, each piece of content will try to download its `*.nfo` metadata.
Metadata downloading can be enabled using the `XTREAM_NFO` environment variable.

### Rate Limiting

In order to make sure the tool doesn't exceed the rate limit of your IPTV provider, the maximum
number of requests per minute can be set using the `XTREAM_REQUESTS` environment variable.

### Filtering

In case you want to exclude any content from specific categories, you can include the name or
substring of a category in the `XTREAM_BANNED_LIVE`, `XTREAM_BANNED_MOVIES` or
`XTREAM_BANNED_SERIES` environment variables as a ',' separated list. For example, if you want to
ban movies from Germany `[DE] ...`, France `[FR] ...`, and the United Kingdom `[UK] ...`, you can
specify these as follows:

```bash
export XTREAM_BANNED_MOVIES="[DE],[FR],[UK]"
```

Specific categories can be excluded by writing the full category name, and groups of categories can
be excluded by writing a substring that is present in all categories of that group.
