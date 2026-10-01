package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func metricsHandler() http.Handler {
	registry = prometheus.NewRegistry()

	registerAppMetrics(registry)
	registerOptionsMetrics(registry)
	registerTimeMetrics(registry)
	registerStreamMetrics(registry)

	ExportMetrics()

	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}

func ExportMetrics() {
	labels := []string{LABEL_LIVE, LABEL_MOVIES, LABEL_SERIES}

	// App Metrics
	export.app.WithLabelValues(metrics.app.Version, metrics.app.Commit, metrics.app.Branch, metrics.app.Date, metrics.app.Modified).Set(1)

	// Options Metrics
	export.launchTime.Set(float64(metrics.options.LaunchTime.Unix()))
	export.launchPeriod.Set(metrics.options.LaunchPeriod.Hours())

	disabledOptions := []bool{metrics.options.DisabledLive, metrics.options.DisabledMovies, metrics.options.DisabledSeries}
	for index, disabled := range disabledOptions {
		label := labels[index]

		if disabled {
			export.disabled.WithLabelValues(label).Set(1.0)
		} else {
			export.disabled.WithLabelValues(label).Set(0.0)
		}
	}

	if metrics.options.ImagesEnabled {
		export.imagesEnabled.Set(1.0)
	} else {
		export.imagesEnabled.Set(0.0)
	}

	if metrics.options.MetadataEnabled {
		export.metadataEnabled.Set(1.0)
	} else {
		export.metadataEnabled.Set(0.0)
	}

	export.requestPerMinute.Set(float64(metrics.options.RequestPerMinute))
	export.requestTimeout.Set(float64(metrics.options.RequestTimeout))

	// Time Metrics
	export.start.Set(float64(metrics.time.Start.Unix()))
	export.end.Set(float64(metrics.time.End.Unix()))
	export.duration.Set(float64(metrics.time.Duration.Seconds()))
	export.next.Set(float64(metrics.time.Next.Unix()))

	// Stream Metrics
	streams := []StreamMetrics{metrics.live, metrics.movies, metrics.series}
	for index, stream := range streams {
		label := labels[index]

		export.streams.WithLabelValues(label).Set(float64(stream.Streams))
		export.categories.WithLabelValues(label).Set(float64(stream.Categories))
		export.cached.WithLabelValues(label).Set(float64(stream.Cached))
		export.metadata.WithLabelValues(label).Set(float64(stream.Metadata))

		export.filteredStreams.WithLabelValues(label).Set(float64(stream.FilteredStreams))
		export.filteredCategories.WithLabelValues(label).Set(float64(stream.FilteredCategories))

		export.updatedStreams.WithLabelValues(label).Set(float64(stream.UpdatedStreams))
		export.updatedImages.WithLabelValues(label).Set(float64(stream.UpdatedImages))
		export.updatedMetadata.WithLabelValues(label).Set(float64(stream.UpdatedMetadata))

		export.requestRate.WithLabelValues(label).Set(stream.RequestRate)

		export.validatedStreams.WithLabelValues(label).Set(float64(stream.ValidatedStreams))
		export.validatedImages.WithLabelValues(label).Set(float64(stream.ValidatedImages))
		export.validatedMetadata.WithLabelValues(label).Set(float64(stream.ValidatedMetadata))
	}
}
