// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"github.com/mattermost/mattermost-server/store"
	"github.com/mattermost/mattermost-server/store/storetest/mocks"
)

type Store struct {
	TeamStore             mocks.TeamStore
	ChannelStore          mocks.ChannelStore
	PostStore             mocks.PostStore
	UserStore             mocks.UserStore
	AuditStore            mocks.AuditStore
	ClusterDiscoveryStore mocks.ClusterDiscoveryStore
	ComplianceStore       mocks.ComplianceStore
	SessionStore          mocks.SessionStore
	OAuthStore            mocks.OAuthStore
	SystemStore           mocks.SystemStore
	WebhookStore          mocks.WebhookStore
	CommandStore          mocks.CommandStore
	CommandWebhookStore   mocks.CommandWebhookStore
	PreferenceStore       mocks.PreferenceStore
	LicenseStore          mocks.LicenseStore
	TokenStore            mocks.TokenStore
	EmojiStore            mocks.EmojiStore
	StatusStore           mocks.StatusStore
	FileInfoStore         mocks.FileInfoStore
	ReactionStore         mocks.ReactionStore
	JobStore              mocks.JobStore
	UserAccessTokenStore  mocks.UserAccessTokenStore
}

func (s *Store) Team() store.TeamStore                         { return &s.TeamStore }
func (s *Store) Channel() store.ChannelStore                   { return &s.ChannelStore }
func (s *Store) Post() store.PostStore                         { return &s.PostStore }
func (s *Store) User() store.UserStore                         { return &s.UserStore }
func (s *Store) Audit() store.AuditStore                       { return &s.AuditStore }
func (s *Store) ClusterDiscovery() store.ClusterDiscoveryStore { return &s.ClusterDiscoveryStore }
func (s *Store) Compliance() store.ComplianceStore             { return &s.ComplianceStore }
func (s *Store) Session() store.SessionStore                   { return &s.SessionStore }
func (s *Store) OAuth() store.OAuthStore                       { return &s.OAuthStore }
func (s *Store) System() store.SystemStore                     { return &s.SystemStore }
func (s *Store) Webhook() store.WebhookStore                   { return &s.WebhookStore }
func (s *Store) Command() store.CommandStore                   { return &s.CommandStore }
func (s *Store) CommandWebhook() store.CommandWebhookStore     { return &s.CommandWebhookStore }
func (s *Store) Preference() store.PreferenceStore             { return &s.PreferenceStore }
func (s *Store) License() store.LicenseStore                   { return &s.LicenseStore }
func (s *Store) Token() store.TokenStore                       { return &s.TokenStore }
func (s *Store) Emoji() store.EmojiStore                       { return &s.EmojiStore }
func (s *Store) Status() store.StatusStore                     { return &s.StatusStore }
func (s *Store) FileInfo() store.FileInfoStore                 { return &s.FileInfoStore }
func (s *Store) Reaction() store.ReactionStore                 { return &s.ReactionStore }
func (s *Store) Job() store.JobStore                           { return &s.JobStore }
func (s *Store) UserAccessToken() store.UserAccessTokenStore   { return &s.UserAccessTokenStore }
func (s *Store) MarkSystemRanUnitTests()                       { /* do nothing */ }
func (s *Store) Close()                                        { /* do nothing */ }
func (s *Store) DropAllTables()                                { /* do nothing */ }
func (s *Store) TotalMasterDbConnections() int                 { return 1 }
func (s *Store) TotalReadDbConnections() int                   { return 1 }
func (s *Store) TotalSearchDbConnections() int                 { return 1 }
