package service

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/amreo/ercole-services/model"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestProcessHostDataInsertion_SuccessNewHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", p("2019-11-05T14:02:03Z")).Return(model.HostData{}, nil).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, "The server 'superhost1' was added to ercole", alert.Description)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.Date)
	}).Times(1)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(str2oid("5dc3f534db7e81a98b726a52")).Return(model.HostData{}, aerrMock).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", p("2019-11-05T14:02:03Z")).Return(model.HostData{}, aerrMock).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DiffHostError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", p("2019-11-05T14:02:03Z")).Return(model.HostData{}, nil).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestDiffHostDataAndGenerateAlert_SuccessNoDifferences(t *testing.T) {
	as := AlertService{}

	require.Nil(t, as.DiffHostDataAndGenerateAlert(hostData2, hostData1))
}

func TestDiffHostDataAndGenerateAlert_SuccessNewHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewServer,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Nil(t, as.DiffHostDataAndGenerateAlert(model.HostData{}, hostData1))
}

func TestDiffHostDataAndGenerateAlert_SuccessNewDatabase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewDatabase,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
			"dbname":   "acd",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewOption,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
			"dbname":   "acd",
			"features": []string{},
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Nil(t, as.DiffHostDataAndGenerateAlert(hostData1, hostData3))
}

func TestDiffHostDataAndGenerateAlert_SuccessNewEnterpriseLicense(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewOption,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
			"dbname":   "acd",
			"features": []string{"Driving"},
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Nil(t, as.DiffHostDataAndGenerateAlert(hostData3, hostData4))
}

func TestDiffHostDataAndGenerateAlert_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewServer,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataAndGenerateAlert(model.HostData{}, hostData1))
}

func TestDiffHostDataAndGenerateAlert_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewDatabase,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
			"dbname":   "acd",
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataAndGenerateAlert(hostData1, hostData3))
}

func TestDiffHostDataAndGenerateAlert_DatabaseError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataAndGenerateAlert(hostData3, hostData4))
}

func TestDiffHostDataAndGenerateAlert_DatabaseError4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewOption,
		OtherInfo: map[string]interface{}{
			"hostname": "superhost1",
			"dbname":   "acd",
			"features": []string{"Driving"},
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataAndGenerateAlert(hostData3, hostData4))
}
