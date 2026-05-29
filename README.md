# Xtream-API

This is a Xtream API writing in Golang that can be used to fetch IPTV Streams & VOD Content. This
tool will generate `.strm` files for all Livestream, Movie, and Series content offered by the IPTV
service. These files can be used to import into Jellyfin or Plex in order to create your own IPTV
Player.

## Environment Variables

| **Variables**            | **Description**                                                         | **Default** |
| ------------------------ | ----------------------------------------------------------------------- | ----------- |
| `XTREAM_URL`             | The URL used for fetching IPTV streams.                                 | ` `         |
| `XTREAM_USERNAME`        | The Username used to log in to the IPTV service.                        | ` `         |
| `XTREAM_PASSWORD`        | The Password used to log in to the IPTV service.                        | ` `         |
| `XTREAM_LAUNCH`          | The time at which the tool should run, in `hh:mm:ss` format.*           | ` `         |
| `XTREAM_PERIOD`          | The period with which the tool should run, in hours.                    | `24`        |
| `XTREAM_IMAGES`          | Whether the tool should download images alongside stream file.          | `false`     |
| `XTREAM_METADATA`        | Whether the tool should download metadata alongside stream file.        | `false`     |
| `XTREAM_REQUESTS`        | The maximum number of requests per minute.                              | `1000`      |
| `XTREAM_TIMEOUT`         | The maximum time before a request returns a timeout error, in seconds.  | `30`        |
| `XTREAM_DISABLED_LIVE`   | Whether the tool should not download any Livestreams.                   | `false`     |
| `XTREAM_DISABLED_MOVIES` | Whether the tool should not download any Movies.                        | `false`     |
| `XTREAM_DISABLED_SERIES` | Whether the tool should not download any Series.                        | `false`     |
| `XTREAM_BANNED_LIVE`     | A ',' separated list of banned livestream (partial) **category** names. | ` `         |
| `XTREAM_BANNED_MOVIES`   | A ',' separated list of banned movie (partial) **category** names.      | ` `         |
| `XTREAM_BANNED_SERIES`   | A ',' separated list of banned series (partial) **category** names.     | ` `         |

**NOTE:** In case `XTREAM_LAUNCH` is not specified, the tool will run immediately. After that, every
24 hours from that time onwards. If `XTREAM_LAUNCH` is properly specified, it will run at every 24
hours at that specific time.

The `XTREAM_PERIOD` environment variable can be used to run the tool
more often, for example every 12 hours. It is recommended to only use this after the initial few
runs have been performed, and only small additions are being made, as this might cause the tool to
run continuously because the runtime is longer than the period.

## Docker Compose

Here is an example of how to use this tool as a docker container, using a docker compose file:

```yaml
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
      # - XTREAM_LAUNCH="00:00:00"        # Set Launch Time at 00:00:00 (Midnight)
      # - XTREAM_PERIOD=12                # Set Launch Period to every 12 hours

      # - XTREAM_IMAGES=true              # Enable Image Fetching
      # - XTREAM_METADATA=true            # Enable Metadata Fetching
      # - XTREAM_REQUESTS=1000            # Maximum requests per minute
      # - XTREAM_TIMEOUT=30               # Maximum time before request returns timeout error

      # - XTREAM_DISABLED_LIVE=true       # Disable Livestreams
      # - XTREAM_DISABLED_MOVIES=true     # Disable Livestreams
      # - XTREAM_DISABLED_SERIES=true     # Disable Livestreams

      # - XTREAM_BANNED_LIVE=""           # ',' Seperated list of banned livestream categories
      # - XTREAM_BANNED_MOVIES=""         # ',' Seperated list of banned movie categories
      # - XTREAM_BANNED_SERIES=""         # ',' Seperated list of banned series categories
    volumes:
      - /etc/localtime:/etc/localtime:ro
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

To disable any of the three types of content, use the `XTREAM_DISABLED_...` environment variables to
disable them.

The tool will use the `/cache` directory to store any cached content for future runs.

### Images & Metadata

If downloading images is enabled, each piece of content will try to download its `cover.jpg` image.
Image downloading can be enabled using the `XTREAM_IMAGES` environment variable.

If downloading metadata is enabled, each piece of content will try to download its `*.nfo` metadata.
Metadata downloading can be enabled using the `XTREAM_METADATA` environment variable.

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

### Caching

Many IPTV providers have a huge content library. Fetching all of this data every time is going to be
very time and energy consuming. To mitigate this, caching is used. The first time the tool is
launched, it will gather all the data from the IPTV provider. Once it has completed the initial run,
it will cache all the processed entries into a JSON file in the `/cache` directory. All runs from
here on out will use these cached entries to check if any entries don't need updating, and can be
skipped. This will massively save time and requests on successive runs.
