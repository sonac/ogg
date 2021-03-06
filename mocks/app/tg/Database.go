// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	models "next-german-words/app/store/models"

	mock "github.com/stretchr/testify/mock"
)

// Database is an autogenerated mock type for the Database type
type Database struct {
	mock.Mock
}

// Connect provides a mock function with given fields:
func (_m *Database) Connect() {
	_m.Called()
}

// FindWord provides a mock function with given fields: german
func (_m *Database) FindWord(german string) (*models.Word, error) {
	ret := _m.Called(german)

	var r0 *models.Word
	if rf, ok := ret.Get(0).(func(string) *models.Word); ok {
		r0 = rf(german)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Word)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(german)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByChatId provides a mock function with given fields: chatId
func (_m *Database) GetUserByChatId(chatId int64) (*models.User, error) {
	ret := _m.Called(chatId)

	var r0 *models.User
	if rf, ok := ret.Get(0).(func(int64) *models.User); ok {
		r0 = rf(chatId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(chatId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWords provides a mock function with given fields:
func (_m *Database) GetWords() ([]*models.Word, error) {
	ret := _m.Called()

	var r0 []*models.Word
	if rf, ok := ret.Get(0).(func() []*models.Word); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Word)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertUser provides a mock function with given fields: user
func (_m *Database) InsertUser(user *models.User) (bool, error) {
	ret := _m.Called(user)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*models.User) bool); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertWord provides a mock function with given fields: word
func (_m *Database) InsertWord(word *models.Word) error {
	ret := _m.Called(word)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.Word) error); ok {
		r0 = rf(word)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUser provides a mock function with given fields: user
func (_m *Database) UpdateUser(user *models.User) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
