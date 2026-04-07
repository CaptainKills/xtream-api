# Xtream-API

This is an Xtream API writing in Golang that can be used to fetch IPTV Streams & VOD Content.
This tool will generate `.strm` files for all Livestream, Movie, and Series content offered by the IPTV service.
These files can be used to import into Jellyfin or Plex in order to create your own IPTV Player.

## Environment Variables

| **Variables**      | **Description**                                                | **Default** |
| ------------------ | -------------------------------------------------------------- | ----------- |
| `XTREAM_URL`       | The URL used for fetching IPTV streams.                        | ""          |
| `XTREAM_USERNAME`  | The Username used to log in to the IPTV service.               | ""          |
| `XTREAM_PASSWORD`  | The Password used to log in to the IPTV service.               | ""          | 
| `XTREAM_IMAGES`    | Whether the tool should download images alongside stream file. | `true`      |

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
    volumes:
      - /etc/localtime:/etc/localtime:ro
      - /path/to/media:/media             # Change to your desired media directory on Host
    restart: unless-stopped
```

The content fetched by this tool will be downloaded in the `/media` directory:
* `/media/live` directory for Livestream content.
* `/media/movies` for Movie content.
* `/media/series` for Series content.

If downloading images is enabled, each piece of content will try to download its `cover.jpg` image.
