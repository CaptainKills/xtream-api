package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusMetrics struct {
	// app Metrics
	app *prometheus.GaugeVec

	// Option Metrics
	launchTime   prometheus.Gauge
	launchPeriod prometheus.Gauge

	disabled *prometheus.GaugeVec

	imagesEnabled   prometheus.Gauge
	metadataEnabled prometheus.Gauge

	requestPerMinute prometheus.Gauge
	requestTimeout   prometheus.Gauge

	// Time Metrics
	start    prometheus.Gauge
	end      prometheus.Gauge
	duration prometheus.Gauge
	next     prometheus.Gauge

	// Stream Metrics
	streams    *prometheus.GaugeVec
	categories *prometheus.GaugeVec
	cached     *prometheus.GaugeVec
	metadata   *prometheus.GaugeVec

	filteredStreams    *prometheus.GaugeVec
	filteredCategories *prometheus.GaugeVec

	updatedStreams  *prometheus.GaugeVec
	updatedImages   *prometheus.GaugeVec
	updatedMetadata *prometheus.GaugeVec

	failedStreams *prometheus.GaugeVec
	requestRate   *prometheus.GaugeVec

	validatedStreams  *prometheus.GaugeVec
	validatedImages   *prometheus.GaugeVec
	validatedMetadata *prometheus.GaugeVec
}

const (
	NAMESPACE = "xtreamapi"
	TYPE      = "type"

	SUBSYSTEM_APP     = "app"
	SUBSYSTEM_OPTIONS = "options"
	SUBSYSTEM_TIME    = "time"
	SUBSYSTEM_STREAM  = "streams"

	LABEL_LIVE   = "live"
	LABEL_MOVIES = "movies"
	LABEL_SERIES = "series"
)

var (
	registry *prometheus.Registry
	export   PrometheusMetrics
)

func registerAppMetrics(reg prometheus.Registerer) {
	export.app = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_APP,
		Name:      "info",
		Help:      "General information about Xtream-API",
	}, []string{"version", "commit", "branch", "date", "modified"})
}

func registerOptionsMetrics(reg prometheus.Registerer) {
	export.launchTime = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_OPTIONS,
		Name:      "launch_time",
		Help:      "The desired launch time of the application",
	})

	export.launchPeriod = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_OPTIONS,
		Name:      "launch_period",
		Help:      "The desired launch period of the application",
	})

	export.disabled = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_OPTIONS,
		Name:      "disabled",
		Help:      "Which of the stream types is disabled",
	}, []string{TYPE})

	export.imagesEnabled = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_OPTIONS,
		Name:      "images_enabled",
		Help:      "Whether fetching images is enabled",
	})

	export.metadataEnabled = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_OPTIONS,
		Name:      "metadata_enabled",
		Help:      "Whether generating metadata is enabled",
	})

	export.requestPerMinute = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_OPTIONS,
		Name:      "request_per_minute",
		Help:      "The maximum allowed request per minute",
	})

	export.requestTimeout = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_OPTIONS,
		Name:      "request_timeout",
		Help:      "The alloted time before a request times out",
	})
}

func registerTimeMetrics(reg prometheus.Registerer) {
	export.start = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_TIME,
		Name:      "start",
		Help:      "The start time of the most recent run",
	})

	export.end = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_TIME,
		Name:      "end",
		Help:      "The end time of the most recent run",
	})

	export.duration = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_TIME,
		Name:      "duration",
		Help:      "The duration of the most recent run",
	})

	export.next = promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_TIME,
		Name:      "next",
		Help:      "The start time of the next run",
	})
}

func registerStreamMetrics(reg prometheus.Registerer) {
	export.streams = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "total",
		Help:      "The number of streams fetched",
	}, []string{TYPE})

	export.categories = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "categories",
		Help:      "The number of categories fetched",
	}, []string{TYPE})

	export.cached = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "cached",
		Help:      "The number of streams cached",
	}, []string{TYPE})

	export.metadata = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "metadata",
		Help:      "The number of metadata cached",
	}, []string{TYPE})

	export.filteredStreams = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "filtered_streams",
		Help:      "The remaining number of streams after filtering",
	}, []string{TYPE})

	export.filteredCategories = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "filtered_categories",
		Help:      "The remaining number of categories after filtering",
	}, []string{TYPE})

	export.updatedStreams = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "updated_streams",
		Help:      "The number of updated streams during the most recent run",
	}, []string{TYPE})

	export.updatedImages = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "updated_images",
		Help:      "The number of updated images during the most recent run",
	}, []string{TYPE})

	export.updatedMetadata = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "updated_metadata",
		Help:      "The number of updated metadata during the most recent run",
	}, []string{TYPE})

	export.failedStreams = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "failed_streams",
		Help:      "The number of failed stream exports during the most recent run",
	}, []string{TYPE})

	export.requestRate = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "request_rate",
		Help:      "The actual request rate during the most recent run",
	}, []string{TYPE})

	export.validatedStreams = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "validated_streams",
		Help:      "The number of validated streams during the most recent run",
	}, []string{TYPE})

	export.validatedImages = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "validated_images",
		Help:      "The number of validated images during the most recent run",
	}, []string{TYPE})

	export.validatedMetadata = promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Subsystem: SUBSYSTEM_STREAM,
		Name:      "validated_metadata",
		Help:      "The number of validated metadata during the most recent run",
	}, []string{TYPE})
}
