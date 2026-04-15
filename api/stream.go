package api

type Stream interface {
	Export(client *XtreamClient, dir string) (int, int, int, error)
}
