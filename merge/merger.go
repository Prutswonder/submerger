package merge

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type (
	// Merger executes the merge command.
	Merger interface {
		Run(path string) error
	}

	mergerImpl struct {
		Path                        string
		Verbose                     bool
		SupportedMovieExtensions    []string
		SupportedSubtitleExtensions []string
		MergedMovieExtension        string

		movies     map[string]movie
		log        Logger
		fileWalker FileWalker
		commander  Commander
	}

	movie struct {
		moviePath    string
		subtitlePath string
		isSubbed     bool
	}
)

// NewMerger instantiates a new Merger.
func NewMerger(supportedMovieExtensions,
	supportedSubtitleExtensions []string,
	mergedMovieExtension string,
	logger Logger,
	fileWalker FileWalker,
	commander Commander,
	verbose bool) Merger {

	return &mergerImpl{
		log:                         logger,
		fileWalker:                  fileWalker,
		commander:                   commander,
		SupportedMovieExtensions:    supportedMovieExtensions,
		SupportedSubtitleExtensions: supportedSubtitleExtensions,
		MergedMovieExtension:        mergedMovieExtension,
		Verbose:                     verbose,
	}
}

func (c *mergerImpl) Run(path string) error {
	if path == "" {
		path = "."
	}

	c.Path = path

	c.log.Printf("Looking for movies in: %s\n", path)

	c.movies = make(map[string]movie)
	subbedMovies := []string{}
	err := c.fileWalker.Walk(path, c.scan)

	if err != nil {
		c.log.Printf("FileWalker returned %v\n", err)
		return err
	}

	for key := range c.movies {
		movie := c.movies[key]
		if !movie.isSubbed && movie.moviePath != "" && movie.subtitlePath != "" {
			subbedMovie, err := c.merge(key, movie, c.Verbose)
			if err != nil {
				return err
			}
			subbedMovies = append(subbedMovies, subbedMovie)
		}
	}

	if len(subbedMovies) == 0 {
		c.log.Println("\nNo movies found to merge.")
	} else {
		c.log.Println("\nThe following movies were merged:")

		for _, movie := range subbedMovies {
			absMovie, _ := filepath.Abs(movie)
			c.log.Println(absMovie)
		}
	}
	return nil
}

func (c *mergerImpl) scan(path string, fileInfo os.FileInfo, err error) error {
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
	case hasAnySuffix(refPath, c.SupportedSubtitleExtensions):

		if val, exists := c.movies[keyPath]; exists {
			val.subtitlePath = path
			c.movies[keyPath] = val
		} else {
			c.movies[keyPath] = movie{subtitlePath: path}
		}
	}
	return nil
}

func (c *mergerImpl) merge(pathWithoutExtension string, movie movie, verbose bool) (subbedMoviePath string, err error) {
	var stdOut, stdErr bytes.Buffer

	subbedMoviePath = pathWithoutExtension + c.MergedMovieExtension
	var args []string

	if verbose {
		args = append(args, "--verbose")
	}

	args = append(args, movie.moviePath, movie.subtitlePath, "-o", subbedMoviePath)

	cmd := c.commander.Command("mkvmerge", args...)

	// In order to run MKV-Merge properly we need to ensure that the LC_ALL environment
	// variable is set.
	env := replaceOrAppend(os.Environ(), "LC_ALL=", "LC_ALL=C")
	cmd.SetEnvironment(env)

	c.log.Printf("Merging to %s...\n", subbedMoviePath)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd")
		return "", err
	}
	defer outReader.Close()

	outScanner := bufio.NewScanner(outReader)
	outScanner.Split(bufio.ScanRunes)

	go func() {
		for outScanner.Scan() {
			c := outScanner.Text()
			fmt.Print(c)
		}
	}()

	errReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StderrPipe for Cmd")
		return "", err
	}
	defer errReader.Close()

	errScanner := bufio.NewScanner(errReader)
	go func() {
		for errScanner.Scan() {
			fmt.Printf("\t > %s\n", errScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd")
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// MKVMerge sometimes exits with status 1, even though the merge succeeded.
			if exitErr.String() != "exit status 1" {
				fmt.Fprintln(os.Stderr, "Error waiting for Cmd")
				return "", err
			}
		} else {
			c.log.Printf("I/O error: %v\n", err)
		}
	}

	// Only proceed once the process has finished

	if stdErr.Len() > 0 {
		return "", fmt.Errorf(stdErr.String())
	}

	c.log.Println(stdOut.String())
	return subbedMoviePath, nil
}

func hasAnySuffix(str string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}

func replaceOrAppend(list []string, find string, replaceWith string) []string {
	newList := list
	for i, v := range list {
		if strings.HasPrefix(v, find) {
			newList = append(list[:i], list[i+1:]...)
			break
		}
	}
	return append(newList, replaceWith)
}
