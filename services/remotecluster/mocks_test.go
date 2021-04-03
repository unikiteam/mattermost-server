// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package remotecluster

import (
	"fmt"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest/mock"
	"github.com/mattermost/mattermost-server/v5/shared/mlog"
	"github.com/mattermost/mattermost-server/v5/store"
	"github.com/mattermost/mattermost-server/v5/store/storetest/mocks"
)

type mockServer struct {
	remotes []*model.RemoteCluster
	logger  *mockLogger
}

func newMockServer(t *testing.T, remotes []*model.RemoteCluster) *mockServer {
	return &mockServer{
		remotes: remotes,
		logger:  &mockLogger{t: t},
	}
}

func (ms *mockServer) Config() *model.Config                                  { return nil }
func (ms *mockServer) GetMetrics() einterfaces.MetricsInterface               { return nil }
func (ms *mockServer) IsLeader() bool                                         { return true }
func (ms *mockServer) AddClusterLeaderChangedListener(listener func()) string { return model.NewId() }
func (ms *mockServer) RemoveClusterLeaderChangedListener(id string)           {}
func (ms *mockServer) GetLogger() mlog.LoggerIFace {
	return ms.logger
}
func (ms *mockServer) GetStore() store.Store {
	anyFilter := mock.MatchedBy(func(filter model.RemoteClusterQueryFilter) bool {
		return true
	})

	remoteClusterStoreMock := &mocks.RemoteClusterStore{}
	remoteClusterStoreMock.On("GetByTopic", "share").Return(ms.remotes, nil)
	remoteClusterStoreMock.On("GetAll", anyFilter).Return(ms.remotes, nil)

	storeMock := &mocks.Store{}
	storeMock.On("RemoteCluster").Return(remoteClusterStoreMock)
	return storeMock
}

type mockLogger struct {
	t *testing.T
}

func (ml *mockLogger) IsLevelEnabled(level mlog.LogLevel) bool {
	return true
}
func (ml *mockLogger) Debug(s string, flds ...mlog.Field) {
	ml.t.Log("debug", s, fieldsToStrings(flds))
}
func (ml *mockLogger) Info(s string, flds ...mlog.Field) {
	ml.t.Log("info", s, fieldsToStrings(flds))
}
func (ml *mockLogger) Warn(s string, flds ...mlog.Field) {
	ml.t.Log("warn", s, fieldsToStrings(flds))
}
func (ml *mockLogger) Error(s string, flds ...mlog.Field) {
	ml.t.Log("error", s, fieldsToStrings(flds))
}
func (ml *mockLogger) Critical(s string, flds ...mlog.Field) {
	ml.t.Log("crit", s, fieldsToStrings(flds))
}
func (ml *mockLogger) Log(level mlog.LogLevel, s string, flds ...mlog.Field) {
	ml.t.Log(level.Name, s, fieldsToStrings(flds))
}
func (ml *mockLogger) LogM(levels []mlog.LogLevel, s string, flds ...mlog.Field) {
	ml.t.Log(levelsToString(levels), s, fieldsToStrings(flds))
}

func levelsToString(levels []mlog.LogLevel) string {
	sb := strings.Builder{}
	for _, l := range levels {
		sb.WriteString(l.Name)
		sb.WriteString(",")
	}
	return sb.String()
}

func fieldsToStrings(fields []mlog.Field) []string {
	encoder := zapcore.NewMapObjectEncoder()
	for _, zapField := range fields {
		zapField.AddTo(encoder)
	}

	var result []string
	for k, v := range encoder.Fields {
		result = append(result, fmt.Sprintf("%s:%v", k, v))
	}
	return result
}
