package gar

type Source interface {
	Open(name string) (file File, ok bool, err error)
	Files() ([]string, error)
}
