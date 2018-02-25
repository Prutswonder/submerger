package merge_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/Prutswonder/submerger/merge"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	movieExt := []string{}
	subExt := []string{}
	logger := &mockLogger{}
	fileWalker := &mockFileWalker{}

	logger.On("Printf", mock.Anything, mock.Anything).Return(nil)
	logger.On("Println", mock.Anything, mock.Anything).Return(nil)
	fileWalker.On("Walk", mock.Anything, mock.Anything).Return(nil).Once()

	sut := merge.NewMerger(movieExt, subExt, "_sub", logger, fileWalker)

	// Act
	err := sut.Run("somepath")

	assert.Nil(t, err)
	fileWalker.AssertExpectations(t)
}
