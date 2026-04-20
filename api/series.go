package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/CaptainKills/xtream-api/utils"
)

type Series struct {
	Cast           string `json:"cast"`
	CategoryId     string `json:"category_id"` // int
	CategoryIds    []int  `json:"category_ids"`
	Cover          string `json:"cover"`
	Director       string `json:"director"`
	EpisodeRunTime string `json:"episode_run_time"` // int
	Genre          string `json:"genre"`
	Id             int    `json:"series_id"`
	LastModified   string `json:"last_modified"` // time.Time
	Name           string `json:"name"`
	Number         int    `json:"num"`
	Plot           string `json:"plot"`
	Rating         string `json:"rating"`        // float64
	Rating5Based   string `json:"rating_5based"` // float64
	ReleaseDate    string `json:"releaseDate"`   // time.Time
	ReleaseDate2   string `json:"release_date"`  // time.Time
	TMDB           string `json:"tmdb"`          // int
	Trailer        string `json:"youtube_trailer"`
}

type SeriesInfo struct {
	Episodes map[string][]Episode `json:"episodes"`
	Info     ExtraSeriesInfo      `json:"info"`
}

type ExtraSeriesInfo struct {
	ActorArray    []Actor  `xml:"actor"`
	Cast          string   `json:"cast" xml:"-"`
	Cover         string   `json:"cover" xml:"-"`
	Director      string   `json:"director" xml:"-"`
	DirectorArray []string `xml:"director"`
	Genre         string   `json:"genre" xml:"-"`
	GenreArray    []string `xml:"genre"`
	Name          string   `json:"name" xml:"title"`
	Plot          string   `json:"plot" xml:"plot"`
	ReleaseDate   string   `json:"releaseDate" xml:"releasedate"` // time.Time
	XMLName       xml.Name `xml:"tvshow"`
}

type Episode struct {
	Added        string   `json:"added" xml:"-"` // time.Time
	CustomSID    string   `json:"custom_sid" xml:"-"`
	DirectSource string   `json:"direct_source" xml:"-"`
	Extension    string   `json:"container_extension" xml:"-"`
	Id           string   `json:"id" xml:"-"` // int
	Number       int      `json:"episode_num" xml:"episode"`
	Season       int      `json:"season" xml:"season"`
	Title        string   `json:"title" xml:"title"`
	XMLName      xml.Name `xml:"episodeinfo"`
}

func (c *XtreamClient) GetSeries() (map[int]Series, error) {
	var series []Series
	series_map := map[int]Series{}

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeries)

	// Fetch Series
	resp, err := c.sendRequest(query)
	if err != nil {
		return map[int]Series{}, err
	}

	// Unmarshal Series
	err = json.Unmarshal(resp, &series)
	if err != nil {
		return map[int]Series{}, err
	}

	// Map Series
	for _, show := range series {
		series_map[show.Id] = show
	}

	return series_map, nil
}

func (c *XtreamClient) GetSeriesInfo(id int) (SeriesInfo, error) {
	var info SeriesInfo
	action := fmt.Sprintf(actionSeriesInfo, id)
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, action)

	resp, err := c.sendRequest(query)
	if err != nil {
		return SeriesInfo{}, err
	}

	err = json.Unmarshal(resp, &info)
	if err != nil {
		return SeriesInfo{}, err
	}

	return info, nil
}

func (s Series) Export(c *XtreamClient, dir string) (int, int, int, error) {
	updated_streams := 0
	updated_images := 0
	updated_nfos := 0

	s.Name = strings.ReplaceAll(s.Name, "/", "_")

	pathDirectory := dir + s.Name
	pathImage := pathDirectory + "/cover" + utils.GetImageExtension(s.Cover)
	pathNfo := pathDirectory + "/tvshow.nfo"

	// Create Subdirectory
	err := os.Mkdir(pathDirectory, 0o755)
	if err != nil && !os.IsExist(err) {
		return updated_streams, updated_images, updated_nfos, err
	}

	// Fetch Series Info
	info, err := c.GetSeriesInfo(s.Id)
	if err != nil {
		return updated_streams, updated_images, updated_nfos, err
	}

	// Write Episodes to File
	for season, episodes := range info.Episodes {
		pathSeason := pathDirectory + "/Season " + season

		// Create Season Subdirectory
		err := os.Mkdir(pathSeason, 0o755)
		if err != nil && !os.IsExist(err) {
			return updated_streams, updated_images, updated_nfos, err
		}

		for _, episode := range episodes {
			updated_stream, updated_nfo, err := episode.Export(c, pathSeason)
			if err != nil {
				return updated_streams, updated_images, updated_nfos, err
			}

			updated_streams += updated_stream
			updated_nfos += updated_nfo
		}
	}

	// Write Image to File
	if c.Options.ImagesEnabled && !utils.ImageExists(pathImage) && strings.HasPrefix(s.Cover, "http") {
		image, err := c.sendRequest(s.Cover)
		if err != nil {
			// Ignore error for image fetching
			// log.Printf("[WARNING] Failed to fetch Image: %v\n", err)
		} else {
			updated_image, err := utils.WriteImage(pathImage, image)
			if err != nil {
				return updated_streams, updated_images, updated_nfos, err
			}

			updated_images += updated_image
		}
	}

	// Write NFO to File
	if c.Options.MetadataEnabled {
		updated_nfo, err := utils.WriteFile(pathNfo, info.GenerateNfo())
		if err != nil {
			return updated_streams, updated_images, updated_nfos, err
		}

		updated_nfos += updated_nfo
	}

	return updated_streams, updated_images, updated_nfos, nil
}

func (i SeriesInfo) GenerateNfo() string {
	genres := strings.SplitSeq(i.Info.Genre, ",")
	for genre := range genres {
		i.Info.GenreArray = append(i.Info.GenreArray, strings.Trim(genre, " "))
	}

	directors := strings.SplitSeq(i.Info.Director, ",")
	for director := range directors {
		i.Info.DirectorArray = append(i.Info.DirectorArray, strings.Trim(director, " "))
	}

	actors := strings.Split(i.Info.Cast, ",")
	for index, actor := range actors {
		i.Info.ActorArray = append(i.Info.ActorArray, Actor{Name: strings.Trim(actor, " "), Order: index})
	}

	nfo, err := xml.MarshalIndent(i.Info, "", "  ")
	if err != nil {
		log.Println(err)
	}

	return xmlHeader + string(nfo)
}

func (e Episode) Export(c *XtreamClient, dir string) (int, int, error) {
	updated_stream := 0
	updated_nfo := 0

	e.Title = strings.ReplaceAll(e.Title, "/", "_")

	pathStream := dir + "/" + e.Title + ".strm"
	pathNfo := dir + "/" + e.Title + ".nfo"

	id, err := strconv.Atoi(e.Id)
	if err != nil {
		return updated_stream, updated_nfo, err
	}
	url := c.buildURL("series", id, e.Extension)

	// Write Stream to File
	updated_stream, err = utils.WriteFile(pathStream, url)
	if err != nil {
		return updated_stream, updated_nfo, err
	}

	// Write NFO to File
	if c.Options.MetadataEnabled {
		updated_nfo, err = utils.WriteFile(pathNfo, e.GenerateNfo())
		if err != nil {
			return updated_stream, updated_nfo, err
		}
	}

	return updated_stream, updated_nfo, nil
}

func (e Episode) GenerateNfo() string {
	var title strings.Builder

	sections := strings.Split(e.Title, " - ")
	for index, part := range sections {
		if index == 0 && len(sections) > 1 {
			continue
		}
		title.WriteString(part)
		if index != len(sections)-1 {
			title.WriteString(" ")
		}
	}

	re := regexp.MustCompile(`[sS][0-9]*[eE][0-9]*\s`)
	e.Title = re.ReplaceAllString(title.String(), "")
	e.Title = strings.Trim(e.Title, " ")

	nfo, err := xml.MarshalIndent(e, "", "  ")
	if err != nil {
		log.Println(err)
	}

	return xmlHeader + string(nfo)
}
