package api

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	encoding = "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>"
	regex    = regexp.MustCompile(`[sS][0-9]*[eE][0-9]*\s`)
)

func GenerateMovieNfo(info MovieInfo) string {
	builder := &strings.Builder{}

	builder.WriteString(encoding)
	builder.WriteString("<movie>")

	fmt.Fprintf(builder, "<title>%s</title>", info.Info.Name)
	fmt.Fprintf(builder, "<originaltitle>%s</originaltitle>", info.Info.OriginalName)
	fmt.Fprintf(builder, "<plot>%s</plot>", info.Info.Plot)
	fmt.Fprintf(builder, "<releasedate>%s</releasedate>", info.Info.ReleaseDate)

	genres := strings.SplitSeq(info.Info.Genre, ", ")
	for genre := range genres {
		fmt.Fprintf(builder, "<genre>%s</genre>", genre)
	}

	directors := strings.SplitSeq(info.Info.Director, ", ")
	for director := range directors {
		fmt.Fprintf(builder, "<director>%s</director>", director)
	}

	actors := strings.Split(info.Info.Cast, ", ")
	for index, actor := range actors {
		fmt.Fprintf(builder, "<actor>")
		fmt.Fprintf(builder, "<name>%s</name>", actor)
		fmt.Fprintf(builder, "<order>%d</order>", index)
		fmt.Fprintf(builder, "</actor>")
	}

	builder.WriteString("</movie>")

	return builder.String()
}

func GenerateSeriesNfo(info SeriesInfo) string {
	builder := &strings.Builder{}

	builder.WriteString(encoding)
	builder.WriteString("<tvshow>")

	fmt.Fprintf(builder, "<title>%s</title>", info.Info.Name)
	fmt.Fprintf(builder, "<plot>%s</plot>", info.Info.Plot)
	fmt.Fprintf(builder, "<releasedate>%s</releasedate>", info.Info.ReleaseDate)

	genres := strings.SplitSeq(info.Info.Genre, ", ")
	for genre := range genres {
		fmt.Fprintf(builder, "<genre>%s</genre>", genre)
	}

	directors := strings.SplitSeq(info.Info.Director, ", ")
	for director := range directors {
		fmt.Fprintf(builder, "<director>%s</director>", director)
	}

	actors := strings.Split(info.Info.Cast, ", ")
	for index, actor := range actors {
		fmt.Fprintf(builder, "<actor>")
		fmt.Fprintf(builder, "<name>%s</name>", actor)
		fmt.Fprintf(builder, "<order>%d</order>", index)
		fmt.Fprintf(builder, "</actor>")
	}

	builder.WriteString("</tvshow>")

	return builder.String()
}

func GenerateEpisodeNfo(episode Episode) string {
	builder := &strings.Builder{}
	episode.Title = strings.ReplaceAll(episode.Title, "&", "&amp;")

	var title string
	sections := strings.Split(episode.Title, " - ")
	for index, part := range sections {
		if index == 0 {
			continue
		}
		title += part
		if index != len(sections)-1 {
			title += " "
		}
	}
	title = regex.ReplaceAllString(title, "")

	builder.WriteString(encoding)
	builder.WriteString("<episodedetails>")
	fmt.Fprintf(builder, "<title>%s</title>", title)
	builder.WriteString("</episodedetails>")

	return builder.String()
}
