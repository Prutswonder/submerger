package merge_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Prutswonder/submerger/merge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScannerRun(t *testing.T) {
	movieExt := []string{".movie-ext"}
	subExt := []string{".sub-ext"}
	logger := &mockLogger{}
	fileWalker := &mockFileWalker{}
	dirInfo := &mockFileInfo{}
	fileInfo := &mockFileInfo{}

	logger.On("Printf", mock.Anything, mock.Anything).Return(nil)
	logger.On("Println", mock.Anything, mock.Anything).Return(nil)
	dirInfo.On("IsDir", mock.Anything).Return(true)
	fileInfo.On("IsDir", mock.Anything).Return(false)

	fileWalker.On("Walk", mock.Anything, mock.Anything).
		Run(func(a mock.Arguments) {
			var scanErr error
			scanFn := a.Get(1).(filepath.WalkFunc)
			scanFn("/path/to", dirInfo, scanErr)
			scanFn("/path/to/a-movie-without-subs.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-with-subs.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-with-subs.sub-ext", fileInfo, scanErr)
			scanFn("/path/to/not-a-movie.mkv", fileInfo, scanErr)
			scanFn("/path/to/a-movie-but.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-but_subbed.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-but_subbed.sub-ext", fileInfo, scanErr)
			scanFn("/path/to/another-movie-but_subbed.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/another-movie-but_subbed.sub-ext", fileInfo, scanErr)
			scanFn("/path/to/another-movie-but.movie-ext", fileInfo, scanErr)
		}).
		Return(nil).
		Once()

	sut := merge.NewScanner(movieExt, subExt, "_subbed.movie-ext", logger, fileWalker)

	// Act: With old movies
	err := sut.Run("")

	assert.Nil(t, err)
	fileWalker.AssertExpectations(t)

	// Act: No old movies
	fileWalker.On("Walk", mock.Anything, mock.Anything).
		Run(func(a mock.Arguments) {
			var scanErr error
			scanFn := a.Get(1).(filepath.WalkFunc)
			scanFn("/path/to", dirInfo, scanErr)
			scanFn("/path/to/a-movie-with-subs.sub-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-with-subs.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-without-subs.movie-ext", fileInfo, scanErr)
			scanFn("/path/to/not-a-movie.mkv", fileInfo, scanErr)
			scanFn("/path/to/a-movie-but_subbed.sub-ext", fileInfo, scanErr)
			scanFn("/path/to/a-movie-but_subbed.movie-ext", fileInfo, scanErr)
		}).
		Return(nil).
		Once()
	err = sut.Run(".")

	assert.Nil(t, err)
	fileWalker.AssertExpectations(t)
}

func TestScannerRun_FileWalkerError(t *testing.T) {
	movieExt := []string{}
	subExt := []string{}
	logger := &mockLogger{}
	fileWalker := &mockFileWalker{}
	expectedErr := errors.New("ewpz")

	logger.On("Printf", mock.Anything, mock.Anything).Return(nil)
	logger.On("Println", mock.Anything, mock.Anything).Return(nil)
	fileWalker.On("Walk", mock.Anything, mock.Anything).Return(expectedErr).Once()

	sut := merge.NewScanner(movieExt, subExt, "_sub", logger, fileWalker)

	// Act
	actualErr := sut.Run("")

	assert.Equal(t, expectedErr, actualErr)
	fileWalker.AssertExpectations(t)
}
