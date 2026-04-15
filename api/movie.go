package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/CaptainKills/xtream-api/utils"
)

type Movie struct {
	Added        string `json:"added"`       // time.Time
	CategoryId   string `json:"category_id"` // int
	CategoryIds  []int  `json:"category_ids"`
	CystomSID    string `json:"custom_sid"`
	DirectSource string `json:"direct_source"`
	Extension    string `json:"container_extension"`
	Icon         string `json:"stream_icon"`
	Id           int    `json:"stream_id"`
	IsAdult      int    `json:"is_adult"` // bool
	Name         string `json:"name"`
	Number       int    `json:"num"`
	Rating       string `json:"rating"` // float64
	StreamType   string `json:"stream_type"`
	Trailer      string `json:"trailer"`
}

type MovieInfo struct {
	Info ExtraMovieInfo `json:"info"`
}

type ExtraMovieInfo struct {
	ActorArray    []Actor  `xml:"actor"`
	Actors        string   `json:"actors" xml:"-"`
	Bitrate       int      `json:"bitrate" xml:"-"`
	Cast          string   `json:"cast" xml:"-"`
	Director      string   `json:"director" xml:"-"`
	DirectorArray []string `xml:"director"`
	Genre         string   `json:"genre" xml:"-"`
	GenreArray    []string `xml:"genre"`
	Name          string   `json:"name" xml:"title"`
	OriginalName  string   `json:"o_name" xml:"originaltitle"`
	Plot          string   `json:"plot" xml:"plot"`
	ReleaseDate   string   `json:"releasedate" xml:"releasedate"` // time.Time
	XMLName       xml.Name `xml:"movie"`
}

type Actor struct {
	Name  string `xml:"name"`
	Order int    `xml:"order"`
}

func (c *XtreamClient) GetMovies() (map[int]Movie, error) {
	var movies []Movie
	movie_map := map[int]Movie{}

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovies)

	// Fetch Movies
	resp, err := c.sendRequest(query)
	if err != nil {
		return map[int]Movie{}, err
	}

	// Unmarshal Movies
	err = json.Unmarshal(resp, &movies)
	if err != nil {
		return map[int]Movie{}, err
	}

	// Map Movies
	for _, movie := range movies {
		movie_map[movie.Id] = movie
	}

	return movie_map, nil
}

func (c *XtreamClient) GetMovieInfo(id int) (MovieInfo, error) {
	var info MovieInfo
	action := fmt.Sprintf(actionMovieInfo, id)
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, action)

	resp, err := c.sendRequest(query)
	if err != nil {
		return MovieInfo{}, err
	}

	err = json.Unmarshal(resp, &info)
	if err != nil {
		return MovieInfo{}, err
	}

	return info, nil
}

func (m Movie) Export(c *XtreamClient, dir string) (int, int, int, error) {
	updated_stream := 0
	updated_image := 0
	updated_nfo := 0

	m.Name = strings.ReplaceAll(m.Name, "/", "_")

	pathDirectory := dir + m.Name
	pathStream := pathDirectory + "/" + m.Name + ".strm"
	pathImage := pathDirectory + "/cover" + utils.GetImageExtension(m.Icon)
	pathNfo := pathDirectory + "/movie.nfo"
	url := c.buildURL(m.StreamType, m.Id, m.Extension)

	// Create Subdirectory
	err := os.Mkdir(pathDirectory, 0o755)
	if err != nil && !os.IsExist(err) {
		return updated_stream, updated_image, updated_nfo, err
	}

	// Write Stream to File
	updated_stream, err = utils.WriteFile(pathStream, url)
	if err != nil {
		return updated_stream, updated_image, updated_nfo, err
	}

	// Write Image to File
	if c.Options.ImagesEnabled && !utils.ImageExists(pathImage) && strings.HasPrefix(m.Icon, "http") {
		image, err := c.sendRequest(m.Icon)
		if err != nil {
			// Ignore error for image fetching
			// log.Printf("[WARNING] Failed to fetch Image: %v\n", err)
		} else {
			updated_image, err = utils.WriteImage(pathImage, image)
			if err != nil {
				return updated_stream, updated_image, updated_nfo, err
			}
		}
	}

	// Write NFO to File
	if c.Options.MetadataEnabled {
		info, err := c.GetMovieInfo(m.Id)
		if err != nil {
			// Ignore error for info fetching
			// log.Printf("[WARNING] Failed to fetch Movie Info: %v\n", err)
		} else {
			updated_nfo, err = utils.WriteFile(pathNfo, info.GenerateNfo())
			if err != nil {
				return updated_stream, updated_image, updated_nfo, err
			}
		}
	}

	return updated_stream, updated_image, updated_nfo, nil
}

func (i MovieInfo) GenerateNfo() string {
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
