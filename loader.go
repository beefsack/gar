package gar

type Loader struct {
	Sources []Source
}

func NewLoader(sources ...Source) *Loader {
	return &Loader{
		Sources: sources,
	}
}

func (l *Loader) AddSource(source Source) {
	l.Sources = append(l.Sources, source)
}

func (l *Loader) Open(name string) (file File, ok bool, err error) {
	for _, s := range l.Sources {
		file, ok, err = s.Open(name)
		if err != nil || ok {
			return
		}
	}
	return
}

func (l *Loader) Files() ([]string, error) {
	fileMap := map[string]bool{}
	files := []string{}
	for _, s := range l.Sources {
		sourceFiles, err := s.Files()
		if err != nil {
			return nil, err
		}
		for _, f := range sourceFiles {
			if !fileMap[f] {
				fileMap[f] = true
				files = append(files, f)
			}
		}
	}
	return files, nil
}
