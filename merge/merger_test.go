package merge_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Prutswonder/submerger/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMergerRun(t *testing.T) {
	movieExt := []string{".movie-ext"}
	subExt := []string{".sub-ext"}
	logger := &mockLogger{}
	fileWalker := &mockFileWalker{}
	commander := &mockCommander{}
	dirInfo := &mockFileInfo{}
	fileInfo := &mockFileInfo{}
	cmd := &mockCmd{}
	stdOut := ioutil.NopCloser(bytes.NewReader([]byte("")))
	stdErr := ioutil.NopCloser(bytes.NewReader([]byte("")))

	logger.On("Printf", mock.Anything, mock.Anything).Return(nil)
	logger.On("Println", mock.Anything, mock.Anything).Return(nil)
	dirInfo.On("IsDir", mock.Anything).Return(true)
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
			scanFn("/path/to", dirInfo, scanErr)
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

	sut := merge.NewMerger(movieExt, subExt, "_subbed.movie-ext", logger, fileWalker, commander, false)

	// Act
	err := sut.Run("somepath")

	assert.Nil(t, err)
	fileWalker.AssertExpectations(t)
	commander.AssertExpectations(t)
	cmd.AssertExpectations(t)
}

func TestMergerRun_FileWalkerError(t *testing.T) {
	movieExt := []string{}
	subExt := []string{}
	logger := &mockLogger{}
	fileWalker := &mockFileWalker{}
	commander := &mockCommander{}
	expectedErr := errors.New("ewpz")

	logger.On("Printf", mock.Anything, mock.Anything).Return(nil)
	logger.On("Println", mock.Anything, mock.Anything).Return(nil)
	fileWalker.On("Walk", mock.Anything, mock.Anything).Return(expectedErr).Once()

	sut := merge.NewMerger(movieExt, subExt, "_sub", logger, fileWalker, commander, true)

	// Act
	actualErr := sut.Run("")

	assert.Equal(t, expectedErr, actualErr)
	fileWalker.AssertExpectations(t)
}

func TestMergerRun_PipeErrors(t *testing.T) {
	movieExt := []string{".movie-ext"}
	subExt := []string{".sub-ext"}
	logger := &mockLogger{}
	fileWalker := &mockFileWalker{}
	commander := &mockCommander{}
	fileInfo := &mockFileInfo{}
	cmd := &mockCmd{}
	stdOut := ioutil.NopCloser(bytes.NewReader([]byte("")))
	stdErr := ioutil.NopCloser(bytes.NewReader([]byte("")))
	expectedErr := errors.New("Heya")

	logger.On("Printf", mock.Anything, mock.Anything).Return(nil)
	logger.On("Println", mock.Anything, mock.Anything).Return(nil)
	fileInfo.On("IsDir", mock.Anything).Return(false)
	cmd.On("SetEnvironment", mock.Anything)
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
		Return(nil)
	commander.On("Command", "mkvmerge", []string{
		"/path/to/a-movie-with-subs.movie-ext",
		"/path/to/a-movie-with-subs.sub-ext",
		"-o",
		"/path/to/a-movie-with-subs_subbed.movie-ext",
	}).Return(cmd)

	sut := merge.NewMerger(movieExt, subExt, "_subbed.movie-ext", logger, fileWalker, commander, false)

	// Act 1: StdOut
	cmd.On("StdoutPipe").Return(stdOut, expectedErr).Once()
	actualErr := sut.Run("somepath")

	assert.Equal(t, expectedErr, actualErr)
	cmd.AssertExpectations(t)

	// Act 2: StdErr
	cmd.On("StdoutPipe").Return(stdOut, nil).Once()
	cmd.On("StderrPipe").Return(stdErr, expectedErr).Once()
	actualErr = sut.Run("anotherpath")

	assert.Equal(t, expectedErr, actualErr)
	cmd.AssertExpectations(t)
}

func TestRun_CommandErrors(t *testing.T) {
	movieExt := []string{".movie-ext"}
	subExt := []string{".sub-ext"}
	logger := &mockLogger{}
	fileWalker := &mockFileWalker{}
	commander := &mockCommander{}
	fileInfo := &mockFileInfo{}
	cmd := &mockCmd{}
	stdOut := ioutil.NopCloser(bytes.NewReader([]byte("")))
	stdErr := ioutil.NopCloser(bytes.NewReader([]byte("")))
	expectedErr := &exec.ExitError{}

	logger.On("Printf", mock.Anything, mock.Anything).Return(nil)
	logger.On("Println", mock.Anything, mock.Anything).Return(nil)
	fileInfo.On("IsDir", mock.Anything).Return(false)
	cmd.On("SetEnvironment", mock.Anything)
	cmd.On("StdoutPipe").Return(stdOut, nil)
	cmd.On("StderrPipe").Return(stdErr, nil)
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
		Return(nil)
	commander.On("Command", "mkvmerge", []string{
		"/path/to/a-movie-with-subs.movie-ext",
		"/path/to/a-movie-with-subs.sub-ext",
		"-o",
		"/path/to/a-movie-with-subs_subbed.movie-ext",
	}).Return(cmd)

	sut := merge.NewMerger(movieExt, subExt, "_subbed.movie-ext", logger, fileWalker, commander, false)

	// Act 1: Tool exit error
	cmd.On("Start", mock.Anything).Return(nil).Once()
	cmd.On("Wait", mock.Anything).Return(expectedErr).Once()
	actualErr := sut.Run("somepath")

	assert.Equal(t, expectedErr, actualErr)
	cmd.AssertExpectations(t)

	// Act 2: Wait error
	commonErr := errors.New("Nope")
	cmd.On("Start", mock.Anything).Return(nil).Once()
	cmd.On("Wait", mock.Anything).Return(commonErr).Once()
	actualErr = sut.Run("anotherpath")

	assert.Nil(t, actualErr)
	cmd.AssertExpectations(t)

	// Act 3: Start error
	cmd.On("Start", mock.Anything).Return(expectedErr).Once()
	actualErr = sut.Run("yetanotherpath")

	assert.Equal(t, expectedErr, actualErr)
	cmd.AssertExpectations(t)
}
