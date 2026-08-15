package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// We can't fully test ClubService without a mock repository yet,
// so we will test the constructor and error variables for now.

func TestNewClubService(t *testing.T) {
	// Pass nil as repository since we're just testing the constructor
	svc := NewClubService(nil, nil)
	assert.NotNil(t, svc)
}

func TestClubServiceErrors(t *testing.T) {
	assert.NotNil(t, ErrAlreadyMember)
	assert.Equal(t, "user is already a member of this club", ErrAlreadyMember.Error())

	assert.NotNil(t, ErrClubNotFound)
	assert.Equal(t, "club not found", ErrClubNotFound.Error())
}

