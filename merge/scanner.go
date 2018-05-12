package merge

import (
	"os"
	"path/filepath"
	"strings"
)

type (
	// Scanner executes the scan command.
	Scanner interface {
		Run(path string) error
	}

	scannerImpl struct {
		Path                        string
		Verbose                     bool
		SupportedMovieExtensions    []string
		SupportedSubtitleExtensions []string
		MergedMovieExtension        string

		movies     map[string]movie
		log        Logger
		fileWalker FileWalker
	}
)

// NewScanner instantiates a new Scanner.
func NewScanner(supportedMovieExtensions,
	supportedSubtitleExtensions []string,
	mergedMovieExtension string,
	logger Logger,
	fileWalker FileWalker) Scanner {

	return &scannerImpl{
		log:                         logger,
		fileWalker:                  fileWalker,
		SupportedMovieExtensions:    supportedMovieExtensions,
		SupportedSubtitleExtensions: supportedSubtitleExtensions,
		MergedMovieExtension:        mergedMovieExtension,
	}
}

func (c *scannerImpl) Run(path string) error {
	if path == "" {
		path = "."
	}

	c.Path = path
	c.movies = make(map[string]movie)
	oldMovies := []string{}
	err := c.fileWalker.Walk(path, c.scan)

	if err != nil {
		c.log.Printf("filepath.Walk(%v) returned %v\n", path, err)
		return err
	}

	for key := range c.movies {
		movie := c.movies[key]
		if movie.isSubbed && movie.moviePath != "" {
			oldMovies = append(oldMovies, movie.moviePath)
		}
	}

	if len(oldMovies) == 0 {
		c.log.Printf("No old movie files found in [%v].\n", path)
	} else {
		c.log.Printf("These old movie files were found in [%v]:\n", path)

		for _, movie := range oldMovies {
			absMovie, _ := filepath.Abs(movie)
			c.log.Println(absMovie)
		}
	}
	return nil
}

func (c *scannerImpl) scan(path string, fileInfo os.FileInfo, err error) error {
	if fileInfo.IsDir() {
		return nil
	}

	refPath := strings.ToLower(path)
	ext := filepath.Ext(path)
	keyPath := strings.TrimSuffix(path, c.MergedMovieExtension)
	keyPath = strings.TrimSuffix(keyPath, ext)
	keyPath = strings.TrimSuffix(keyPath, ".en")

	switch {
	case strings.HasSuffix(refPath, c.MergedMovieExtension):

		if val, exists := c.movies[keyPath]; exists {
			val.isSubbed = true
			c.movies[keyPath] = val
		} else {
			c.movies[keyPath] = movie{isSubbed: true}
		}
	case hasAnySuffix(refPath, c.SupportedMovieExtensions):

		if val, exists := c.movies[keyPath]; exists {
			val.moviePath = path
			c.movies[keyPath] = val
		} else {
			c.movies[keyPath] = movie{moviePath: path}
		}
	}
	return nil
}
