package repository_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework/internal/domain"
	"homework/internal/repository"
	"strconv"
	"testing"
)

type RepoSuite struct {
	suite.Suite
	repo *repository.Repo
}

func (suite *RepoSuite) SetupTest() {
	suite.repo = repository.New()
}

func (suite *RepoSuite) TearDownTest() {
	suite.repo = nil
}

func (suite *RepoSuite) TestGetDevice() {
	serialNum := "1"
	device := domain.Device{
		SerialNum: serialNum,
		Model:     "test_model",
		IP:        "0.0.0.0",
	}
	suite.repo.Devices[serialNum] = device

	suite.Run("Existing Device", func() {
		d, err := suite.repo.GetDevice(serialNum)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), device, d)
	})

	suite.Run("Unexisting Device", func() {
		_, err := suite.repo.GetDevice("unexisting_serial")
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, "no such device")
	})
}

func (suite *RepoSuite) TestCreateDevice() {
	serialNum := "1"
	device := domain.Device{
		SerialNum: serialNum,
		Model:     "test_model",
		IP:        "0.0.0.0",
	}

	suite.Run("New Device", func() {
		err := suite.repo.CreateDevice(device)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), device, suite.repo.Devices[serialNum])
	})

	suite.Run("Existing Device", func() {
		err := suite.repo.CreateDevice(device)
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, "device is already in repository")
	})
}

func (suite *RepoSuite) TestDeleteDevice() {
	serialNum := "1"
	device := domain.Device{
		SerialNum: serialNum,
		Model:     "test_model",
		IP:        "0.0.0.0",
	}
	suite.repo.Devices[serialNum] = device

	suite.Run("Existing Device", func() {
		err := suite.repo.DeleteDevice(serialNum)
		assert.NoError(suite.T(), err)
		_, ok := suite.repo.Devices[serialNum]
		assert.False(suite.T(), ok)
	})

	suite.Run("Unexisting Device", func() {
		err := suite.repo.DeleteDevice("unexisting_serial")
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, "no such device")
	})
}

func (suite *RepoSuite) TestUpdateDevice() {
	serialNum := "1"
	device := domain.Device{
		SerialNum: serialNum,
		Model:     "test_model",
		IP:        "0.0.0.0",
	}
	suite.repo.Devices[serialNum] = device

	suite.Run("Existing Device", func() {
		updatedDevice := domain.Device{
			SerialNum: serialNum,
			Model:     "updated_model",
			IP:        "updated_ip",
		}
		err := suite.repo.UpdateDevice(updatedDevice)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), updatedDevice, suite.repo.Devices[serialNum])
	})

	suite.Run("Unexisting Device", func() {
		nonExistingDevice := domain.Device{
			SerialNum: "unexisting_serial",
			Model:     "test_model",
			IP:        "0.0.0.0",
		}
		err := suite.repo.UpdateDevice(nonExistingDevice)
		assert.Error(suite.T(), err)
		assert.EqualError(suite.T(), err, "device not found")
	})
}

func BenchmarkCreateDevice(b *testing.B) {
	repo := repository.New()
	var devices []domain.Device
	for i := 0; i < b.N; i++ {
		d := domain.Device{
			SerialNum: strconv.Itoa(i),
			Model:     "xxx",
			IP:        "0.0.0.0",
		}
		devices = append(devices, d)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := repo.CreateDevice(devices[i])
		if err != nil {
			b.Errorf("unexpected error: %v", err)
		}
	}
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}
