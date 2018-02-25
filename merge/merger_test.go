package merge_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/Prutswonder/submerger/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRun(t *testing.T) {
	movieExt := []string{".movie-ext"}
	subExt := []string{".sub-ext"}
	logger := &mockLogger{}
	fileWalker := &mockFileWalker{}
	commander := &mockCommander{}
	fileInfo := &mockFileInfo{}
	cmd := &mockCmd{}
	stdOut := ioutil.NopCloser(bytes.NewReader([]byte("")))
	stdErr := ioutil.NopCloser(bytes.NewReader([]byte("")))

	logger.On("Printf", mock.Anything, mock.Anything).Return(nil)
	logger.On("Println", mock.Anything, mock.Anything).Return(nil)
	fileInfo.On("IsDir", mock.Anything).Return(false)
	cmd.On("SetEnvironment", mock.Anything).Once()
	cmd.On("StdoutPipe").Return(stdOut, nil).Once()
	cmd.On("StderrPipe").Return(stdErr, nil).Once()
	cmd.On("Start", mock.Anything).Return(nil).Once()
	cmd.On("Wait", mock.Anything).Return(nil).Once()

	fileWalker.On("Walk", mock.Anything, mock.Anything).
		Run(func(a mock.Arguments) {
			var scanErr error
			scanFn := a.Get(1).(filepath.WalkFunc)
			scanFn("/path/to/a-movie-without-subs.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-with-subs.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-with-subs.sub-ext", fileInfo, scanErr)
			scanFn("/path/to/not-a-movie.mkv", fileInfo, scanErr)
			scanFn("/path/to/a-movie-but_subbed.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-but_subbed.sub-ext", fileInfo, scanErr)
		}).
		Return(nil).
		Once()
	commander.On("Command", "mkvmerge", []string{
		"/path/to/a-movie-with-subs.movie-ext",
		"/path/to/a-movie-with-subs.sub-ext",
		"-o",
		"/path/to/a-movie-with-subs_subbed.movie-ext",
	}).Return(cmd).Once()

	sut := merge.NewMerger(movieExt, subExt, "_subbed.movie-ext", logger, fileWalker, commander)

	// Act
	err := sut.Run("somepath")

	assert.Nil(t, err)
	fileWalker.AssertExpectations(t)
	commander.AssertExpectations(t)
	cmd.AssertExpectations(t)
}

func TestRun_FileWalkerError(t *testing.T) {
	movieExt := []string{}
	subExt := []string{}
	logger := &mockLogger{}
	fileWalker := &mockFileWalker{}
	commander := &mockCommander{}
	expectedErr := errors.New("ewpz")

	logger.On("Printf", mock.Anything, mock.Anything).Return(nil)
	logger.On("Println", mock.Anything, mock.Anything).Return(nil)
	fileWalker.On("Walk", mock.Anything, mock.Anything).Return(expectedErr).Once()

	sut := merge.NewMerger(movieExt, subExt, "_sub", logger, fileWalker, commander)

	// Act
	actualErr := sut.Run("somepath")

	assert.Equal(t, expectedErr, actualErr)
	fileWalker.AssertExpectations(t)
}
