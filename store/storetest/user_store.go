// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/store"
)

const (
	DAY_MILLISECONDS   = 24 * 60 * 60 * 1000
	MONTH_MILLISECONDS = 31 * DAY_MILLISECONDS
)

func cleanupStatusStore(t *testing.T, s SqlSupplier) {
	_, execerr := s.GetMaster().ExecNoTimeout(` DELETE FROM Status `)
	require.Nil(t, execerr)
}

func TestUserStore(t *testing.T, ss store.Store, s SqlSupplier) {
	users, err := ss.User().GetAll()
	require.Nil(t, err, "failed cleaning up test users")

	for _, u := range users {
		err := ss.User().PermanentDelete(u.Id)
		require.Nil(t, err, "failed cleaning up test user %s", u.Username)
	}

	t.Run("Count", func(t *testing.T) { testCount(t, ss) })
	t.Run("AnalyticsActiveCount", func(t *testing.T) { testUserStoreAnalyticsActiveCount(t, ss, s) })
	t.Run("AnalyticsGetInactiveUsersCount", func(t *testing.T) { testUserStoreAnalyticsGetInactiveUsersCount(t, ss) })
	t.Run("AnalyticsGetSystemAdminCount", func(t *testing.T) { testUserStoreAnalyticsGetSystemAdminCount(t, ss) })
	t.Run("Save", func(t *testing.T) { testUserStoreSave(t, ss) })
	t.Run("Update", func(t *testing.T) { testUserStoreUpdate(t, ss) })
	t.Run("UpdateUpdateAt", func(t *testing.T) { testUserStoreUpdateUpdateAt(t, ss) })
	t.Run("UpdateFailedPasswordAttempts", func(t *testing.T) { testUserStoreUpdateFailedPasswordAttempts(t, ss) })
	t.Run("Get", func(t *testing.T) { testUserStoreGet(t, ss) })
	t.Run("GetAllUsingAuthService", func(t *testing.T) { testGetAllUsingAuthService(t, ss) })
	t.Run("GetAllProfiles", func(t *testing.T) { testUserStoreGetAllProfiles(t, ss) })
	t.Run("GetProfiles", func(t *testing.T) { testUserStoreGetProfiles(t, ss) })
	t.Run("GetProfilesInChannel", func(t *testing.T) { testUserStoreGetProfilesInChannel(t, ss) })
	t.Run("GetProfilesInChannelByStatus", func(t *testing.T) { testUserStoreGetProfilesInChannelByStatus(t, ss, s) })
	t.Run("GetProfilesWithoutTeam", func(t *testing.T) { testUserStoreGetProfilesWithoutTeam(t, ss) })
	t.Run("GetAllProfilesInChannel", func(t *testing.T) { testUserStoreGetAllProfilesInChannel(t, ss) })
	t.Run("GetProfilesNotInChannel", func(t *testing.T) { testUserStoreGetProfilesNotInChannel(t, ss) })
	t.Run("GetProfilesByIds", func(t *testing.T) { testUserStoreGetProfilesByIds(t, ss) })
	t.Run("GetProfileByGroupChannelIdsForUser", func(t *testing.T) { testUserStoreGetProfileByGroupChannelIdsForUser(t, ss) })
	t.Run("GetProfilesByUsernames", func(t *testing.T) { testUserStoreGetProfilesByUsernames(t, ss) })
	t.Run("GetSystemAdminProfiles", func(t *testing.T) { testUserStoreGetSystemAdminProfiles(t, ss) })
	t.Run("GetByEmail", func(t *testing.T) { testUserStoreGetByEmail(t, ss) })
	t.Run("GetByAuthData", func(t *testing.T) { testUserStoreGetByAuthData(t, ss) })
	t.Run("GetByUsername", func(t *testing.T) { testUserStoreGetByUsername(t, ss) })
	t.Run("GetForLogin", func(t *testing.T) { testUserStoreGetForLogin(t, ss) })
	t.Run("UpdatePassword", func(t *testing.T) { testUserStoreUpdatePassword(t, ss) })
	t.Run("Delete", func(t *testing.T) { testUserStoreDelete(t, ss) })
	t.Run("UpdateAuthData", func(t *testing.T) { testUserStoreUpdateAuthData(t, ss) })
	t.Run("UserUnreadCount", func(t *testing.T) { testUserUnreadCount(t, ss) })
	t.Run("UpdateMfaSecret", func(t *testing.T) { testUserStoreUpdateMfaSecret(t, ss) })
	t.Run("UpdateMfaActive", func(t *testing.T) { testUserStoreUpdateMfaActive(t, ss) })
	t.Run("GetRecentlyActiveUsersForTeam", func(t *testing.T) { testUserStoreGetRecentlyActiveUsersForTeam(t, ss, s) })
	t.Run("GetNewUsersForTeam", func(t *testing.T) { testUserStoreGetNewUsersForTeam(t, ss) })
	t.Run("Search", func(t *testing.T) { testUserStoreSearch(t, ss) })
	t.Run("SearchNotInChannel", func(t *testing.T) { testUserStoreSearchNotInChannel(t, ss) })
	t.Run("SearchInChannel", func(t *testing.T) { testUserStoreSearchInChannel(t, ss) })
	t.Run("SearchNotInTeam", func(t *testing.T) { testUserStoreSearchNotInTeam(t, ss) })
	t.Run("SearchWithoutTeam", func(t *testing.T) { testUserStoreSearchWithoutTeam(t, ss) })
	t.Run("GetProfilesNotInTeam", func(t *testing.T) { testUserStoreGetProfilesNotInTeam(t, ss) })
	t.Run("ClearAllCustomRoleAssignments", func(t *testing.T) { testUserStoreClearAllCustomRoleAssignments(t, ss) })
	t.Run("GetAllAfter", func(t *testing.T) { testUserStoreGetAllAfter(t, ss) })
	t.Run("GetUsersBatchForIndexing", func(t *testing.T) { testUserStoreGetUsersBatchForIndexing(t, ss) })
	t.Run("GetTeamGroupUsers", func(t *testing.T) { testUserStoreGetTeamGroupUsers(t, ss) })
	t.Run("GetChannelGroupUsers", func(t *testing.T) { testUserStoreGetChannelGroupUsers(t, ss) })
	t.Run("PromoteGuestToUser", func(t *testing.T) { testUserStorePromoteGuestToUser(t, ss) })
	t.Run("DemoteUserToGuest", func(t *testing.T) { testUserStoreDemoteUserToGuest(t, ss) })
	t.Run("ResetLastPictureUpdate", func(t *testing.T) { testUserStoreResetLastPictureUpdate(t, ss) })
}

func testUserStoreSave(t *testing.T, ss store.Store) {
	teamId := model.NewId()
	maxUsersPerTeam := 50

	u1 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}

	_, err := ss.User().Save(&u1)
	require.Nil(t, err, "couldn't save user")

	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, maxUsersPerTeam)
	require.Nil(t, err)

	_, err = ss.User().Save(&u1)
	require.NotNil(t, err, "shouldn't be able to update user from save")

	u2 := model.User{
		Email:    u1.Email,
		Username: model.NewId(),
	}
	_, err = ss.User().Save(&u2)
	require.NotNil(t, err, "should be unique email")

	u2.Email = MakeEmail()
	u2.Username = u1.Username
	_, err = ss.User().Save(&u1)
	require.NotNil(t, err, "should be unique username")

	u2.Username = ""
	_, err = ss.User().Save(&u1)
	require.NotNil(t, err, "should be unique username")

	for i := 0; i < 49; i++ {
		u := model.User{
			Email:    MakeEmail(),
			Username: model.NewId(),
		}
		_, err = ss.User().Save(&u)
		require.Nil(t, err, "couldn't save item")

		defer func() { require.Nil(t, ss.User().PermanentDelete(u.Id)) }()

		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u.Id}, maxUsersPerTeam)
		require.Nil(t, err)
	}

	u2.Id = ""
	u2.Email = MakeEmail()
	u2.Username = model.NewId()
	_, err = ss.User().Save(&u2)
	require.Nil(t, err, "couldn't save item")

	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, maxUsersPerTeam)
	require.NotNil(t, err, "should be the limit")
}

func testUserStoreUpdate(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Email: MakeEmail(),
	}
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2 := &model.User{
		Email:       MakeEmail(),
		AuthService: "ldap",
	}
	_, err = ss.User().Save(u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u2.Id}, -1)
	require.Nil(t, err)

	_, err = ss.User().Update(u1, false)
	require.Nil(t, err)

	missing := &model.User{}
	_, err = ss.User().Update(missing, false)
	require.NotNil(t, err, "Update should have failed because of missing key")

	newId := &model.User{
		Id: model.NewId(),
	}
	_, err = ss.User().Update(newId, false)
	require.NotNil(t, err, "Update should have failed because id change")

	u2.Email = MakeEmail()
	_, err = ss.User().Update(u2, false)
	require.NotNil(t, err, "Update should have failed because you can't modify AD/LDAP fields")

	u3 := &model.User{
		Email:       MakeEmail(),
		AuthService: "gitlab",
	}
	oldEmail := u3.Email
	_, err = ss.User().Save(u3)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u3.Id}, -1)
	require.Nil(t, err)

	u3.Email = MakeEmail()
	userUpdate, err := ss.User().Update(u3, false)
	require.Nil(t, err, "Update should not have failed")
	assert.Equal(t, oldEmail, userUpdate.New.Email, "Email should not have been updated as the update is not trusted")

	u3.Email = MakeEmail()
	userUpdate, err = ss.User().Update(u3, true)
	require.Nil(t, err, "Update should not have failed")
	assert.NotEqual(t, oldEmail, userUpdate.New.Email, "Email should have been updated as the update is trusted")

	err = ss.User().UpdateLastPictureUpdate(u1.Id)
	require.Nil(t, err, "Update should not have failed")
}

func testUserStoreUpdateUpdateAt(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	require.Nil(t, err)

	_, err = ss.User().UpdateUpdateAt(u1.Id)
	require.Nil(t, err)

	user, err := ss.User().Get(u1.Id)
	require.Nil(t, err)
	require.Less(t, u1.UpdateAt, user.UpdateAt, "UpdateAt not updated correctly")
}

func testUserStoreUpdateFailedPasswordAttempts(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	require.Nil(t, err)

	err = ss.User().UpdateFailedPasswordAttempts(u1.Id, 3)
	require.Nil(t, err)

	user, err := ss.User().Get(u1.Id)
	require.Nil(t, err)
	require.Equal(t, 3, user.FailedAttempts, "FailedAttempts not updated correctly")
}

func testUserStoreGet(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Email: MakeEmail(),
	}
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	u2, _ := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	})
	_, err = ss.Bot().Save(&model.Bot{
		UserId:      u2.Id,
		Username:    u2.Username,
		Description: "bot description",
		OwnerId:     u1.Id,
	})
	require.Nil(t, err)
	u2.IsBot = true
	u2.BotDescription = "bot description"
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u2.Id)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	require.Nil(t, err)

	t.Run("fetch empty id", func(t *testing.T) {
		_, err := ss.User().Get("")
		require.NotNil(t, err)
	})

	t.Run("fetch user 1", func(t *testing.T) {
		actual, err := ss.User().Get(u1.Id)
		require.Nil(t, err)
		require.Equal(t, u1, actual)
		require.False(t, actual.IsBot)
	})

	t.Run("fetch user 2, also a bot", func(t *testing.T) {
		actual, err := ss.User().Get(u2.Id)
		require.Nil(t, err)
		require.Equal(t, u2, actual)
		require.True(t, actual.IsBot)
		require.Equal(t, "bot description", actual.BotDescription)
	})
}

func testGetAllUsingAuthService(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u1" + model.NewId(),
		AuthService: "service",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u2" + model.NewId(),
		AuthService: "service",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u3" + model.NewId(),
		AuthService: "service2",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()

	t.Run("get by unknown auth service", func(t *testing.T) {
		users, err := ss.User().GetAllUsingAuthService("unknown")
		require.Nil(t, err)
		assert.Equal(t, []*model.User{}, users)
	})

	t.Run("get by auth service", func(t *testing.T) {
		users, err := ss.User().GetAllUsingAuthService("service")
		require.Nil(t, err)
		assert.Equal(t, []*model.User{u1, u2}, users)
	})

	t.Run("get by other auth service", func(t *testing.T) {
		users, err := ss.User().GetAllUsingAuthService("service2")
		require.Nil(t, err)
		assert.Equal(t, []*model.User{u3}, users)
	})
}

func sanitized(user *model.User) *model.User {
	clonedUser := model.UserFromJson(strings.NewReader(user.ToJson()))
	clonedUser.AuthData = new(string)
	*clonedUser.AuthData = ""
	clonedUser.Props = model.StringMap{}

	return clonedUser
}

func testUserStoreGetAllProfiles(t *testing.T, ss store.Store) {
	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()

	u4, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
		Roles:    "system_user some-other-role",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u4.Id)) }()

	u5, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u5" + model.NewId(),
		Roles:    "system_admin",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u5.Id)) }()

	u6, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u6" + model.NewId(),
		DeleteAt: model.GetMillis(),
		Roles:    "system_admin",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u6.Id)) }()

	u7, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u7" + model.NewId(),
		DeleteAt: model.GetMillis(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u7.Id)) }()

	t.Run("get offset 0, limit 100", func(t *testing.T) {
		options := &model.UserGetOptions{Page: 0, PerPage: 100}
		actual, err := ss.User().GetAllProfiles(options)
		require.Nil(t, err)

		require.Equal(t, []*model.User{
			sanitized(u1),
			sanitized(u2),
			sanitized(u3),
			sanitized(u4),
			sanitized(u5),
			sanitized(u6),
			sanitized(u7),
		}, actual)
	})

	t.Run("get offset 0, limit 1", func(t *testing.T) {
		actual, err := ss.User().GetAllProfiles(&model.UserGetOptions{
			Page:    0,
			PerPage: 1,
		})
		require.Nil(t, err)
		require.Equal(t, []*model.User{
			sanitized(u1),
		}, actual)
	})

	t.Run("get all", func(t *testing.T) {
		actual, err := ss.User().GetAll()
		require.Nil(t, err)

		require.Equal(t, []*model.User{
			u1,
			u2,
			u3,
			u4,
			u5,
			u6,
			u7,
		}, actual)
	})

	t.Run("etag changes for all after user creation", func(t *testing.T) {
		etag := ss.User().GetEtagForAllProfiles()

		uNew := &model.User{}
		uNew.Email = MakeEmail()
		_, err := ss.User().Save(uNew)
		require.Nil(t, err)
		defer func() { require.Nil(t, ss.User().PermanentDelete(uNew.Id)) }()

		updatedEtag := ss.User().GetEtagForAllProfiles()
		require.NotEqual(t, etag, updatedEtag)
	})

	t.Run("filter to system_admin role", func(t *testing.T) {
		actual, err := ss.User().GetAllProfiles(&model.UserGetOptions{
			Page:    0,
			PerPage: 10,
			Role:    "system_admin",
		})
		require.Nil(t, err)
		require.Equal(t, []*model.User{
			sanitized(u5),
			sanitized(u6),
		}, actual)
	})

	t.Run("filter to system_admin role, inactive", func(t *testing.T) {
		actual, err := ss.User().GetAllProfiles(&model.UserGetOptions{
			Page:     0,
			PerPage:  10,
			Role:     "system_admin",
			Inactive: true,
		})
		require.Nil(t, err)
		require.Equal(t, []*model.User{
			sanitized(u6),
		}, actual)
	})

	t.Run("filter to inactive", func(t *testing.T) {
		actual, err := ss.User().GetAllProfiles(&model.UserGetOptions{
			Page:     0,
			PerPage:  10,
			Inactive: true,
		})
		require.Nil(t, err)
		require.Equal(t, []*model.User{
			sanitized(u6),
			sanitized(u7),
		}, actual)
	})
}

func testUserStoreGetProfiles(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)

	u4, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
		Roles:    "system_admin",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u4.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u4.Id}, -1)
	require.Nil(t, err)

	u5, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u5" + model.NewId(),
		DeleteAt: model.GetMillis(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u5.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u5.Id}, -1)
	require.Nil(t, err)

	t.Run("get page 0, perPage 100", func(t *testing.T) {
		actual, err := ss.User().GetProfiles(&model.UserGetOptions{
			InTeamId: teamId,
			Page:     0,
			PerPage:  100,
		})
		require.Nil(t, err)

		require.Equal(t, []*model.User{
			sanitized(u1),
			sanitized(u2),
			sanitized(u3),
			sanitized(u4),
			sanitized(u5),
		}, actual)
	})

	t.Run("get page 0, perPage 1", func(t *testing.T) {
		actual, err := ss.User().GetProfiles(&model.UserGetOptions{
			InTeamId: teamId,
			Page:     0,
			PerPage:  1,
		})
		require.Nil(t, err)

		require.Equal(t, []*model.User{sanitized(u1)}, actual)
	})

	t.Run("get unknown team id", func(t *testing.T) {
		actual, err := ss.User().GetProfiles(&model.UserGetOptions{
			InTeamId: "123",
			Page:     0,
			PerPage:  100,
		})
		require.Nil(t, err)

		require.Equal(t, []*model.User{}, actual)
	})

	t.Run("etag changes for all after user creation", func(t *testing.T) {
		etag := ss.User().GetEtagForProfiles(teamId)

		uNew := &model.User{}
		uNew.Email = MakeEmail()
		_, err := ss.User().Save(uNew)
		require.Nil(t, err)
		defer func() { require.Nil(t, ss.User().PermanentDelete(uNew.Id)) }()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: uNew.Id}, -1)
		require.Nil(t, err)

		updatedEtag := ss.User().GetEtagForProfiles(teamId)
		require.NotEqual(t, etag, updatedEtag)
	})

	t.Run("filter to system_admin role", func(t *testing.T) {
		actual, err := ss.User().GetProfiles(&model.UserGetOptions{
			InTeamId: teamId,
			Page:     0,
			PerPage:  10,
			Role:     "system_admin",
		})
		require.Nil(t, err)
		require.Equal(t, []*model.User{
			sanitized(u4),
		}, actual)
	})

	t.Run("filter to inactive", func(t *testing.T) {
		actual, err := ss.User().GetProfiles(&model.UserGetOptions{
			InTeamId: teamId,
			Page:     0,
			PerPage:  10,
			Inactive: true,
		})
		require.Nil(t, err)
		require.Equal(t, []*model.User{
			sanitized(u5),
		}, actual)
	})
}

func testUserStoreGetProfilesInChannel(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	ch1 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	c1, err := ss.Channel().Save(ch1, -1)
	require.Nil(t, err)

	ch2 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}
	c2, err := ss.Channel().Save(ch2, -1)
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)
	t.Run("get in channel 1, offset 0, limit 100", func(t *testing.T) {
		users, err := ss.User().GetProfilesInChannel(c1.Id, 0, 100)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u1), sanitized(u2), sanitized(u3)}, users)
	})

	t.Run("get in channel 1, offset 1, limit 2", func(t *testing.T) {
		users, err := ss.User().GetProfilesInChannel(c1.Id, 1, 2)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u2), sanitized(u3)}, users)
	})

	t.Run("get in channel 2, offset 0, limit 1", func(t *testing.T) {
		users, err := ss.User().GetProfilesInChannel(c2.Id, 0, 1)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u1)}, users)
	})
}

func testUserStoreGetProfilesInChannelByStatus(t *testing.T, ss store.Store, s SqlSupplier) {

	cleanupStatusStore(t, s)

	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	ch1 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	c1, err := ss.Channel().Save(ch1, -1)
	require.Nil(t, err)

	ch2 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}
	c2, err := ss.Channel().Save(ch2, -1)
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{
		UserId: u1.Id,
		Status: model.STATUS_DND,
	}))
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{
		UserId: u2.Id,
		Status: model.STATUS_AWAY,
	}))
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{
		UserId: u3.Id,
		Status: model.STATUS_ONLINE,
	}))

	t.Run("get in channel 1 by status, offset 0, limit 100", func(t *testing.T) {
		users, err := ss.User().GetProfilesInChannelByStatus(c1.Id, 0, 100)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u3), sanitized(u2), sanitized(u1)}, users)
	})

	t.Run("get in channel 2 by status, offset 0, limit 1", func(t *testing.T) {
		users, err := ss.User().GetProfilesInChannelByStatus(c2.Id, 0, 1)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u1)}, users)
	})
}

func testUserStoreGetProfilesWithoutTeam(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
		DeleteAt: 1,
		Roles:    "system_admin",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	t.Run("get, page 0, per_page 100", func(t *testing.T) {
		users, err := ss.User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 0, PerPage: 100})
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u2), sanitized(u3)}, users)
	})

	t.Run("get, page 1, per_page 1", func(t *testing.T) {
		users, err := ss.User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 1, PerPage: 1})
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u3)}, users)
	})

	t.Run("get, page 2, per_page 1", func(t *testing.T) {
		users, err := ss.User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 2, PerPage: 1})
		require.Nil(t, err)
		assert.Equal(t, []*model.User{}, users)
	})

	t.Run("get, page 0, per_page 100, inactive", func(t *testing.T) {
		users, err := ss.User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 0, PerPage: 100, Inactive: true})
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u3)}, users)
	})

	t.Run("get, page 0, per_page 100, role", func(t *testing.T) {
		users, err := ss.User().GetProfilesWithoutTeam(&model.UserGetOptions{Page: 0, PerPage: 100, Role: "system_admin"})
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u3)}, users)
	})
}

func testUserStoreGetAllProfilesInChannel(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	ch1 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	c1, err := ss.Channel().Save(ch1, -1)
	require.Nil(t, err)

	ch2 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}
	c2, err := ss.Channel().Save(ch2, -1)
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	if err != nil {
		panic(err)
	}

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	t.Run("all profiles in channel 1, no caching", func(t *testing.T) {
		var profiles map[string]*model.User
		profiles, err = ss.User().GetAllProfilesInChannel(c1.Id, false)
		require.Nil(t, err)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
			u2.Id: sanitized(u2),
			u3.Id: sanitized(u3),
		}, profiles)
	})

	t.Run("all profiles in channel 2, no caching", func(t *testing.T) {
		var profiles map[string]*model.User
		profiles, err = ss.User().GetAllProfilesInChannel(c2.Id, false)
		require.Nil(t, err)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
		}, profiles)
	})

	t.Run("all profiles in channel 2, caching", func(t *testing.T) {
		var profiles map[string]*model.User
		profiles, err = ss.User().GetAllProfilesInChannel(c2.Id, true)
		require.Nil(t, err)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
		}, profiles)
	})

	t.Run("all profiles in channel 2, caching [repeated]", func(t *testing.T) {
		var profiles map[string]*model.User
		profiles, err = ss.User().GetAllProfilesInChannel(c2.Id, true)
		require.Nil(t, err)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
		}, profiles)
	})

	ss.User().InvalidateProfilesInChannelCacheByUser(u1.Id)
	ss.User().InvalidateProfilesInChannelCache(c2.Id)
}

func testUserStoreGetProfilesNotInChannel(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	ch1 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in channel",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_OPEN,
	}
	c1, err := ss.Channel().Save(ch1, -1)
	require.Nil(t, err)

	ch2 := &model.Channel{
		TeamId:      teamId,
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_PRIVATE,
	}
	c2, err := ss.Channel().Save(ch2, -1)
	require.Nil(t, err)

	t.Run("get team 1, channel 1, offset 0, limit 100", func(t *testing.T) {
		var profiles []*model.User
		profiles, err = ss.User().GetProfilesNotInChannel(teamId, c1.Id, false, 0, 100, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u1),
			sanitized(u2),
			sanitized(u3),
		}, profiles)
	})

	t.Run("get team 1, channel 2, offset 0, limit 100", func(t *testing.T) {
		var profiles []*model.User
		profiles, err = ss.User().GetProfilesNotInChannel(teamId, c2.Id, false, 0, 100, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u1),
			sanitized(u2),
			sanitized(u3),
		}, profiles)
	})

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	if err != nil {
		panic(err)
	}
	t.Run("get team 1, channel 1, offset 0, limit 100, after update", func(t *testing.T) {
		var profiles []*model.User
		profiles, err = ss.User().GetProfilesNotInChannel(teamId, c1.Id, false, 0, 100, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{}, profiles)
	})

	t.Run("get team 1, channel 2, offset 0, limit 100, after update", func(t *testing.T) {
		var profiles []*model.User
		profiles, err = ss.User().GetProfilesNotInChannel(teamId, c2.Id, false, 0, 100, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u2),
			sanitized(u3),
		}, profiles)
	})

	t.Run("get team 1, channel 2, offset 0, limit 0, setting group constrained when it's not", func(t *testing.T) {
		var profiles []*model.User
		profiles, err = ss.User().GetProfilesNotInChannel(teamId, c2.Id, true, 0, 100, nil)
		require.Nil(t, err)
		assert.Empty(t, profiles)
	})

	// create a group
	group, err := ss.Group().Create(&model.Group{
		Name:        "n_" + model.NewId(),
		DisplayName: "dn_" + model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    "ri_" + model.NewId(),
	})
	require.Nil(t, err)

	// add two members to the group
	for _, u := range []*model.User{u1, u2} {
		_, err = ss.Group().UpsertMember(group.Id, u.Id)
		require.Nil(t, err)
	}

	// associate the group with the channel
	_, err = ss.Group().CreateGroupSyncable(&model.GroupSyncable{
		GroupId:    group.Id,
		SyncableId: c2.Id,
		Type:       model.GroupSyncableTypeChannel,
	})
	require.Nil(t, err)

	t.Run("get team 1, channel 2, offset 0, limit 0, setting group constrained", func(t *testing.T) {
		profiles, err := ss.User().GetProfilesNotInChannel(teamId, c2.Id, true, 0, 100, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u2),
		}, profiles)
	})
}

func testUserStoreGetProfilesByIds(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	u4, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u4.Id)) }()

	t.Run("get u1 by id, no caching", func(t *testing.T) {
		users, err := ss.User().GetProfileByIds([]string{u1.Id}, nil, false)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u1)}, users)
	})

	t.Run("get u1 by id, caching", func(t *testing.T) {
		users, err := ss.User().GetProfileByIds([]string{u1.Id}, nil, true)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u1)}, users)
	})

	t.Run("get u1, u2, u3 by id, no caching", func(t *testing.T) {
		users, err := ss.User().GetProfileByIds([]string{u1.Id, u2.Id, u3.Id}, nil, false)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u1), sanitized(u2), sanitized(u3)}, users)
	})

	t.Run("get u1, u2, u3 by id, caching", func(t *testing.T) {
		users, err := ss.User().GetProfileByIds([]string{u1.Id, u2.Id, u3.Id}, nil, true)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{sanitized(u1), sanitized(u2), sanitized(u3)}, users)
	})

	t.Run("get unknown id, caching", func(t *testing.T) {
		users, err := ss.User().GetProfileByIds([]string{"123"}, nil, true)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{}, users)
	})

	t.Run("should only return users with UpdateAt greater than the since time", func(t *testing.T) {
		users, err := ss.User().GetProfileByIds([]string{u1.Id, u2.Id, u3.Id, u4.Id}, &store.UserGetByIdsOpts{
			Since: u2.CreateAt,
		}, true)
		require.Nil(t, err)

		// u3 comes from the cache, and u4 does not
		assert.Equal(t, []*model.User{sanitized(u3), sanitized(u4)}, users)
	})
}

func testUserStoreGetProfileByGroupChannelIdsForUser(t *testing.T, ss store.Store) {
	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()

	u4, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u4.Id)) }()

	gc1, err := ss.Channel().Save(&model.Channel{
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_GROUP,
	}, -1)
	require.Nil(t, err)

	for _, uId := range []string{u1.Id, u2.Id, u3.Id} {
		_, err = ss.Channel().SaveMember(&model.ChannelMember{
			ChannelId:   gc1.Id,
			UserId:      uId,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		})
		require.Nil(t, err)
	}

	gc2, err := ss.Channel().Save(&model.Channel{
		DisplayName: "Profiles in private",
		Name:        "profiles-" + model.NewId(),
		Type:        model.CHANNEL_GROUP,
	}, -1)
	require.Nil(t, err)

	for _, uId := range []string{u1.Id, u3.Id, u4.Id} {
		_, err = ss.Channel().SaveMember(&model.ChannelMember{
			ChannelId:   gc2.Id,
			UserId:      uId,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		})
		require.Nil(t, err)
	}

	testCases := []struct {
		Name                       string
		UserId                     string
		ChannelIds                 []string
		ExpectedUserIdsByChannel   map[string][]string
		EnsureChannelsNotInResults []string
	}{
		{
			Name:       "Get group 1 as user 1",
			UserId:     u1.Id,
			ChannelIds: []string{gc1.Id},
			ExpectedUserIdsByChannel: map[string][]string{
				gc1.Id: {u2.Id, u3.Id},
			},
			EnsureChannelsNotInResults: []string{},
		},
		{
			Name:       "Get groups 1 and 2 as user 1",
			UserId:     u1.Id,
			ChannelIds: []string{gc1.Id, gc2.Id},
			ExpectedUserIdsByChannel: map[string][]string{
				gc1.Id: {u2.Id, u3.Id},
				gc2.Id: {u3.Id, u4.Id},
			},
			EnsureChannelsNotInResults: []string{},
		},
		{
			Name:       "Get groups 1 and 2 as user 2",
			UserId:     u2.Id,
			ChannelIds: []string{gc1.Id, gc2.Id},
			ExpectedUserIdsByChannel: map[string][]string{
				gc1.Id: {u1.Id, u3.Id},
			},
			EnsureChannelsNotInResults: []string{gc2.Id},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			res, err := ss.User().GetProfileByGroupChannelIdsForUser(tc.UserId, tc.ChannelIds)
			require.Nil(t, err)

			for channelId, expectedUsers := range tc.ExpectedUserIdsByChannel {
				users, ok := res[channelId]
				require.True(t, ok)

				var userIds []string
				for _, user := range users {
					userIds = append(userIds, user.Id)
				}
				require.ElementsMatch(t, expectedUsers, userIds)
			}

			for _, channelId := range tc.EnsureChannelsNotInResults {
				_, ok := res[channelId]
				require.False(t, ok)
			}
		})
	}
}

func testUserStoreGetProfilesByUsernames(t *testing.T, ss store.Store) {
	teamId := model.NewId()
	team2Id := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: team2Id, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	t.Run("get by u1 and u2 usernames, team id 1", func(t *testing.T) {
		users, err := ss.User().GetProfilesByUsernames([]string{u1.Username, u2.Username}, &model.ViewUsersRestrictions{Teams: []string{teamId}})
		require.Nil(t, err)
		assert.Equal(t, []*model.User{u1, u2}, users)
	})

	t.Run("get by u1 username, team id 1", func(t *testing.T) {
		users, err := ss.User().GetProfilesByUsernames([]string{u1.Username}, &model.ViewUsersRestrictions{Teams: []string{teamId}})
		require.Nil(t, err)
		assert.Equal(t, []*model.User{u1}, users)
	})

	t.Run("get by u1 and u3 usernames, no team id", func(t *testing.T) {
		users, err := ss.User().GetProfilesByUsernames([]string{u1.Username, u3.Username}, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{u1, u3}, users)
	})

	t.Run("get by u1 and u3 usernames, team id 1", func(t *testing.T) {
		users, err := ss.User().GetProfilesByUsernames([]string{u1.Username, u3.Username}, &model.ViewUsersRestrictions{Teams: []string{teamId}})
		require.Nil(t, err)
		assert.Equal(t, []*model.User{u1}, users)
	})

	t.Run("get by u1 and u3 usernames, team id 2", func(t *testing.T) {
		users, err := ss.User().GetProfilesByUsernames([]string{u1.Username, u3.Username}, &model.ViewUsersRestrictions{Teams: []string{team2Id}})
		require.Nil(t, err)
		assert.Equal(t, []*model.User{u3}, users)
	})
}

func testUserStoreGetSystemAdminProfiles(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Roles:    model.SYSTEM_USER_ROLE_ID + " " + model.SYSTEM_ADMIN_ROLE_ID,
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Roles:    model.SYSTEM_USER_ROLE_ID + " " + model.SYSTEM_ADMIN_ROLE_ID,
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	t.Run("all system admin profiles", func(t *testing.T) {
		result, userError := ss.User().GetSystemAdminProfiles()
		require.Nil(t, userError)
		assert.Equal(t, map[string]*model.User{
			u1.Id: sanitized(u1),
			u3.Id: sanitized(u3),
		}, result)
	})
}

func testUserStoreGetByEmail(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	t.Run("get u1 by email", func(t *testing.T) {
		u, err := ss.User().GetByEmail(u1.Email)
		require.Nil(t, err)
		assert.Equal(t, u1, u)
	})

	t.Run("get u2 by email", func(t *testing.T) {
		u, err := ss.User().GetByEmail(u2.Email)
		require.Nil(t, err)
		assert.Equal(t, u2, u)
	})

	t.Run("get u3 by email", func(t *testing.T) {
		u, err := ss.User().GetByEmail(u3.Email)
		require.Nil(t, err)
		assert.Equal(t, u3, u)
	})

	t.Run("get by empty email", func(t *testing.T) {
		_, err := ss.User().GetByEmail("")
		require.NotNil(t, err)
		require.Equal(t, err.Id, store.MISSING_ACCOUNT_ERROR)
	})

	t.Run("get by unknown", func(t *testing.T) {
		_, err := ss.User().GetByEmail("unknown")
		require.NotNil(t, err)
		require.Equal(t, err.Id, store.MISSING_ACCOUNT_ERROR)
	})
}

func testUserStoreGetByAuthData(t *testing.T, ss store.Store) {
	teamId := model.NewId()
	auth1 := model.NewId()
	auth3 := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u1" + model.NewId(),
		AuthData:    &auth1,
		AuthService: "service",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u3" + model.NewId(),
		AuthData:    &auth3,
		AuthService: "service2",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	t.Run("get by u1 auth", func(t *testing.T) {
		u, err := ss.User().GetByAuth(u1.AuthData, u1.AuthService)
		require.Nil(t, err)
		assert.Equal(t, u1, u)
	})

	t.Run("get by u3 auth", func(t *testing.T) {
		u, err := ss.User().GetByAuth(u3.AuthData, u3.AuthService)
		require.Nil(t, err)
		assert.Equal(t, u3, u)
	})

	t.Run("get by u1 auth, unknown service", func(t *testing.T) {
		_, err := ss.User().GetByAuth(u1.AuthData, "unknown")
		require.NotNil(t, err)
		require.Equal(t, err.Id, store.MISSING_AUTH_ACCOUNT_ERROR)
	})

	t.Run("get by unknown auth, u1 service", func(t *testing.T) {
		unknownAuth := ""
		_, err := ss.User().GetByAuth(&unknownAuth, u1.AuthService)
		require.NotNil(t, err)
		require.Equal(t, err.Id, store.MISSING_AUTH_ACCOUNT_ERROR)
	})

	t.Run("get by unknown auth, unknown service", func(t *testing.T) {
		unknownAuth := ""
		_, err := ss.User().GetByAuth(&unknownAuth, "unknown")
		require.NotNil(t, err)
		require.Equal(t, err.Id, store.MISSING_AUTH_ACCOUNT_ERROR)
	})
}

func testUserStoreGetByUsername(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	t.Run("get u1 by username", func(t *testing.T) {
		result, err := ss.User().GetByUsername(u1.Username)
		require.Nil(t, err)
		assert.Equal(t, u1, result)
	})

	t.Run("get u2 by username", func(t *testing.T) {
		result, err := ss.User().GetByUsername(u2.Username)
		require.Nil(t, err)
		assert.Equal(t, u2, result)
	})

	t.Run("get u3 by username", func(t *testing.T) {
		result, err := ss.User().GetByUsername(u3.Username)
		require.Nil(t, err)
		assert.Equal(t, u3, result)
	})

	t.Run("get by empty username", func(t *testing.T) {
		_, err := ss.User().GetByUsername("")
		require.NotNil(t, err)
		require.Equal(t, err.Id, "store.sql_user.get_by_username.app_error")
	})

	t.Run("get by unknown", func(t *testing.T) {
		_, err := ss.User().GetByUsername("unknown")
		require.NotNil(t, err)
		require.Equal(t, err.Id, "store.sql_user.get_by_username.app_error")
	})
}

func testUserStoreGetForLogin(t *testing.T, ss store.Store) {
	teamId := model.NewId()
	auth := model.NewId()
	auth2 := model.NewId()
	auth3 := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u1" + model.NewId(),
		AuthService: model.USER_AUTH_SERVICE_GITLAB,
		AuthData:    &auth,
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u2" + model.NewId(),
		AuthService: model.USER_AUTH_SERVICE_LDAP,
		AuthData:    &auth2,
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:       MakeEmail(),
		Username:    "u3" + model.NewId(),
		AuthService: model.USER_AUTH_SERVICE_LDAP,
		AuthData:    &auth3,
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	t.Run("get u1 by username, allow both", func(t *testing.T) {
		user, err := ss.User().GetForLogin(u1.Username, true, true)
		require.Nil(t, err)
		assert.Equal(t, u1, user)
	})

	t.Run("get u1 by username, allow only email", func(t *testing.T) {
		_, err := ss.User().GetForLogin(u1.Username, false, true)
		require.NotNil(t, err)
		require.Equal(t, err.Id, "store.sql_user.get_for_login.app_error")
	})

	t.Run("get u1 by email, allow both", func(t *testing.T) {
		user, err := ss.User().GetForLogin(u1.Email, true, true)
		require.Nil(t, err)
		assert.Equal(t, u1, user)
	})

	t.Run("get u1 by email, allow only username", func(t *testing.T) {
		_, err := ss.User().GetForLogin(u1.Email, true, false)
		require.NotNil(t, err)
		require.Equal(t, err.Id, "store.sql_user.get_for_login.app_error")
	})

	t.Run("get u2 by username, allow both", func(t *testing.T) {
		user, err := ss.User().GetForLogin(u2.Username, true, true)
		require.Nil(t, err)
		assert.Equal(t, u2, user)
	})

	t.Run("get u2 by email, allow both", func(t *testing.T) {
		user, err := ss.User().GetForLogin(u2.Email, true, true)
		require.Nil(t, err)
		assert.Equal(t, u2, user)
	})

	t.Run("get u2 by username, allow neither", func(t *testing.T) {
		_, err := ss.User().GetForLogin(u2.Username, false, false)
		require.NotNil(t, err)
		require.Equal(t, err.Id, "store.sql_user.get_for_login.app_error")
	})
}

func testUserStoreUpdatePassword(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	hashedPassword := model.HashPassword("newpwd")

	err = ss.User().UpdatePassword(u1.Id, hashedPassword)
	require.Nil(t, err)

	user, err := ss.User().GetByEmail(u1.Email)
	require.Nil(t, err)
	require.Equal(t, user.Password, hashedPassword, "Password was not updated correctly")
}

func testUserStoreDelete(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	require.Nil(t, err)

	err = ss.User().PermanentDelete(u1.Id)
	require.Nil(t, err)
}

func testUserStoreUpdateAuthData(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	service := "someservice"
	authData := model.NewId()

	_, err = ss.User().UpdateAuthData(u1.Id, service, &authData, "", true)
	require.Nil(t, err)

	user, err := ss.User().GetByEmail(u1.Email)
	require.Nil(t, err)
	require.Equal(t, service, user.AuthService, "AuthService was not updated correctly")
	require.Equal(t, authData, *user.AuthData, "AuthData was not updated correctly")
	require.Equal(t, "", user.Password, "Password was not cleared properly")
}

func testUserUnreadCount(t *testing.T, ss store.Store) {
	teamId := model.NewId()

	c1 := model.Channel{}
	c1.TeamId = teamId
	c1.DisplayName = "Unread Messages"
	c1.Name = "unread-messages-" + model.NewId()
	c1.Type = model.CHANNEL_OPEN

	c2 := model.Channel{}
	c2.TeamId = teamId
	c2.DisplayName = "Unread Direct"
	c2.Name = "unread-direct-" + model.NewId()
	c2.Type = model.CHANNEL_DIRECT

	u1 := &model.User{}
	u1.Username = "user1" + model.NewId()
	u1.Email = MakeEmail()
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.Username = "user2" + model.NewId()
	_, err = ss.User().Save(u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	_, err = ss.Channel().Save(&c1, -1)
	require.Nil(t, err, "couldn't save item")

	m1 := model.ChannelMember{}
	m1.ChannelId = c1.Id
	m1.UserId = u1.Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := model.ChannelMember{}
	m2.ChannelId = c1.Id
	m2.UserId = u2.Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	_, err = ss.Channel().SaveMember(&m2)
	require.Nil(t, err)

	m1.ChannelId = c2.Id
	m2.ChannelId = c2.Id

	_, err = ss.Channel().SaveDirectChannel(&c2, &m1, &m2)
	require.Nil(t, err, "couldn't save direct channel")

	p1 := model.Post{}
	p1.ChannelId = c1.Id
	p1.UserId = u1.Id
	p1.Message = "this is a message for @" + u2.Username

	// Post one message with mention to open channel
	_, err = ss.Post().Save(&p1)
	require.Nil(t, err)
	err = ss.Channel().IncrementMentionCount(c1.Id, u2.Id)
	require.Nil(t, err)

	// Post 2 messages without mention to direct channel
	p2 := model.Post{}
	p2.ChannelId = c2.Id
	p2.UserId = u1.Id
	p2.Message = "first message"

	_, err = ss.Post().Save(&p2)
	require.Nil(t, err)
	err = ss.Channel().IncrementMentionCount(c2.Id, u2.Id)
	require.Nil(t, err)

	p3 := model.Post{}
	p3.ChannelId = c2.Id
	p3.UserId = u1.Id
	p3.Message = "second message"
	_, err = ss.Post().Save(&p3)
	require.Nil(t, err)

	err = ss.Channel().IncrementMentionCount(c2.Id, u2.Id)
	require.Nil(t, err)

	badge, unreadCountErr := ss.User().GetUnreadCount(u2.Id)
	require.Nil(t, unreadCountErr)
	require.Equal(t, int64(3), badge, "should have 3 unread messages")

	badge, unreadCountErr = ss.User().GetUnreadCountForChannel(u2.Id, c1.Id)
	require.Nil(t, unreadCountErr)
	require.Equal(t, int64(1), badge, "should have 1 unread messages for that channel")

	badge, unreadCountErr = ss.User().GetUnreadCountForChannel(u2.Id, c2.Id)
	require.Nil(t, unreadCountErr)
	require.Equal(t, int64(2), badge, "should have 2 unread messages for that channel")
}

func testUserStoreUpdateMfaSecret(t *testing.T, ss store.Store) {
	u1 := model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(&u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	err = ss.User().UpdateMfaSecret(u1.Id, "12345")
	require.Nil(t, err)

	// should pass, no update will occur though
	err = ss.User().UpdateMfaSecret("junk", "12345")
	require.Nil(t, err)
}

func testUserStoreUpdateMfaActive(t *testing.T, ss store.Store) {
	u1 := model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(&u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	time.Sleep(100 * time.Millisecond)

	err = ss.User().UpdateMfaActive(u1.Id, true)
	require.Nil(t, err)

	err = ss.User().UpdateMfaActive(u1.Id, false)
	require.Nil(t, err)

	// should pass, no update will occur though
	err = ss.User().UpdateMfaActive("junk", true)
	require.Nil(t, err)
}

func testUserStoreGetRecentlyActiveUsersForTeam(t *testing.T, ss store.Store, s SqlSupplier) {

	cleanupStatusStore(t, s)

	teamId := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	millis := model.GetMillis()
	u3.LastActivityAt = millis
	u2.LastActivityAt = millis - 1
	u1.LastActivityAt = millis - 1

	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{UserId: u1.Id, Status: model.STATUS_ONLINE, Manual: false, LastActivityAt: u1.LastActivityAt, ActiveChannel: ""}))
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{UserId: u2.Id, Status: model.STATUS_ONLINE, Manual: false, LastActivityAt: u2.LastActivityAt, ActiveChannel: ""}))
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{UserId: u3.Id, Status: model.STATUS_ONLINE, Manual: false, LastActivityAt: u3.LastActivityAt, ActiveChannel: ""}))

	t.Run("get team 1, offset 0, limit 100", func(t *testing.T) {
		users, err := ss.User().GetRecentlyActiveUsersForTeam(teamId, 0, 100, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u3),
			sanitized(u1),
			sanitized(u2),
		}, users)
	})

	t.Run("get team 1, offset 0, limit 1", func(t *testing.T) {
		users, err := ss.User().GetRecentlyActiveUsersForTeam(teamId, 0, 1, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u3),
		}, users)
	})

	t.Run("get team 1, offset 2, limit 1", func(t *testing.T) {
		users, err := ss.User().GetRecentlyActiveUsersForTeam(teamId, 2, 1, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u2),
		}, users)
	})
}

func testUserStoreGetNewUsersForTeam(t *testing.T, ss store.Store) {
	teamId := model.NewId()
	teamId2 := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	u4, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u4.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId2, UserId: u4.Id}, -1)
	require.Nil(t, err)

	t.Run("get team 1, offset 0, limit 100", func(t *testing.T) {
		result, err := ss.User().GetNewUsersForTeam(teamId, 0, 100, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u3),
			sanitized(u2),
			sanitized(u1),
		}, result)
	})

	t.Run("get team 1, offset 0, limit 1", func(t *testing.T) {
		result, err := ss.User().GetNewUsersForTeam(teamId, 0, 1, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u3),
		}, result)
	})

	t.Run("get team 1, offset 2, limit 1", func(t *testing.T) {
		result, err := ss.User().GetNewUsersForTeam(teamId, 2, 1, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u1),
		}, result)
	})

	t.Run("get team 2, offset 0, limit 100", func(t *testing.T) {
		result, err := ss.User().GetNewUsersForTeam(teamId2, 0, 100, nil)
		require.Nil(t, err)
		assert.Equal(t, []*model.User{
			sanitized(u4),
		}, result)
	})
}

func assertUsers(t *testing.T, expected, actual []*model.User) {
	expectedUsernames := make([]string, 0, len(expected))
	for _, user := range expected {
		expectedUsernames = append(expectedUsernames, user.Username)
	}

	actualUsernames := make([]string, 0, len(actual))
	for _, user := range actual {
		actualUsernames = append(actualUsernames, user.Username)
	}

	if assert.Equal(t, expectedUsernames, actualUsernames) {
		assert.Equal(t, expected, actual)
	}
}

func assertUsersMatchInAnyOrder(t *testing.T, expected, actual []*model.User) {
	expectedUsernames := make([]string, 0, len(expected))
	for _, user := range expected {
		expectedUsernames = append(expectedUsernames, user.Username)
	}

	actualUsernames := make([]string, 0, len(actual))
	for _, user := range actual {
		actualUsernames = append(actualUsernames, user.Username)
	}

	if assert.ElementsMatch(t, expectedUsernames, actualUsernames) {
		assert.ElementsMatch(t, expected, actual)
	}
}

func testUserStoreSearch(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
		Roles:     "system_user system_admin",
	}
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim-bobby" + model.NewId(),
		Email:    MakeEmail(),
		Roles:    "system_user",
	}
	_, err = ss.User().Save(u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
		Roles:    "system_admin",
	}
	_, err = ss.User().Save(u3)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	u5 := &model.User{
		Username:  "yu" + model.NewId(),
		FirstName: "En",
		LastName:  "Yu",
		Nickname:  "enyu",
		Email:     MakeEmail(),
	}
	_, err = ss.User().Save(u5)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u5.Id)) }()

	u6 := &model.User{
		Username:  "underscore" + model.NewId(),
		FirstName: "Du_",
		LastName:  "_DE",
		Nickname:  "lodash",
		Email:     MakeEmail(),
	}
	_, err = ss.User().Save(u6)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u6.Id)) }()

	tid := model.NewId()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u1.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u2.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u3.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u5.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u6.Id}, -1)
	require.Nil(t, err)

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData
	u5.AuthData = nilAuthData
	u6.AuthData = nilAuthData

	testCases := []struct {
		Description string
		TeamId      string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"search jimb",
			tid,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search en",
			tid,
			"en",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u5},
		},
		{
			"search email",
			tid,
			u1.Email,
			&model.UserSearchOptions{
				AllowEmails:    true,
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search maps * to space",
			tid,
			"jimb*",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"should not return spurious matches",
			tid,
			"harol",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"% should be escaped",
			tid,
			"h%",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"_ should be escaped",
			tid,
			"h_",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"_ should be escaped (2)",
			tid,
			"Du_",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u6},
		},
		{
			"_ should be escaped (2)",
			tid,
			"_dE",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u6},
		},
		{
			"search jimb, allowing inactive",
			tid,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1, u3},
		},
		{
			"search jimb, no team id",
			"",
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jim-bobb, no team id",
			"",
			"jim-bobb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2},
		},

		{
			"search harol, search all fields",
			tid,
			"harol",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search Tim, search all fields",
			tid,
			"Tim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search Tim, don't search full names",
			tid,
			"Tim",
			&model.UserSearchOptions{
				AllowFullNames: false,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search Bill, search all fields",
			tid,
			"Bill",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search Rob, search all fields",
			tid,
			"Rob",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowEmails:    true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"leading @ should be ignored",
			tid,
			"@jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jim-bobby with system_user roles",
			tid,
			"jim-bobby",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
				Role:           "system_user",
			},
			[]*model.User{u2},
		},
		{
			"search jim with system_admin roles",
			tid,
			"jim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
				Role:           "system_admin",
			},
			[]*model.User{u1},
		},
		{
			"search ji with system_user roles",
			tid,
			"ji",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
				Role:           "system_user",
			},
			[]*model.User{u1, u2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			users, err := ss.User().Search(testCase.TeamId, testCase.Term, testCase.Options)
			require.Nil(t, err)
			assertUsersMatchInAnyOrder(t, testCase.Expected, users)
		})
	}

	t.Run("search empty string", func(t *testing.T) {
		searchOptions := &model.UserSearchOptions{
			AllowFullNames: true,
			Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
		}

		users, err := ss.User().Search(tid, "", searchOptions)
		require.Nil(t, err)
		assert.Len(t, users, 4)
		// Don't assert contents, since Postgres' default collation order is left up to
		// the operating system, and jimbo1 might sort before or after jim-bo.
		// assertUsers(t, []*model.User{u2, u1, u6, u5}, r1.Data.([]*model.User))
	})

	t.Run("search empty string, limit 2", func(t *testing.T) {
		searchOptions := &model.UserSearchOptions{
			AllowFullNames: true,
			Limit:          2,
		}

		users, err := ss.User().Search(tid, "", searchOptions)
		require.Nil(t, err)
		assert.Len(t, users, 2)
		// Don't assert contents, since Postgres' default collation order is left up to
		// the operating system, and jimbo1 might sort before or after jim-bo.
		// assertUsers(t, []*model.User{u2, u1, u6, u5}, r1.Data.([]*model.User))
	})
}

func testUserStoreSearchNotInChannel(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
	}
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim2-bobby" + model.NewId(),
		Email:    MakeEmail(),
	}
	_, err = ss.User().Save(u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
	}
	_, err = ss.User().Save(u3)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	tid := model.NewId()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u1.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u2.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u3.Id}, -1)
	require.Nil(t, err)

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData

	ch1 := model.Channel{
		TeamId:      tid,
		DisplayName: "NameName",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	c1, err := ss.Channel().Save(&ch1, -1)
	require.Nil(t, err)

	ch2 := model.Channel{
		TeamId:      tid,
		DisplayName: "NameName",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	c2, err := ss.Channel().Save(&ch2, -1)
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	testCases := []struct {
		Description string
		TeamId      string
		ChannelId   string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"search jimb, channel 1",
			tid,
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, allow inactive, channel 1",
			tid,
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, channel 1, no team id",
			"",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, channel 1, junk team id",
			"junk",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, channel 2",
			tid,
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, allow inactive, channel 2",
			tid,
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u3},
		},
		{
			"search jimb, channel 2, no team id",
			"",
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, channel 2, junk team id",
			"junk",
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jim, channel 1",
			tid,
			c1.Id,
			"jim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2, u1},
		},
		{
			"search jim, channel 1, limit 1",
			tid,
			c1.Id,
			"jim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          1,
			},
			[]*model.User{u2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			users, err := ss.User().SearchNotInChannel(
				testCase.TeamId,
				testCase.ChannelId,
				testCase.Term,
				testCase.Options,
			)
			require.Nil(t, err)
			assertUsers(t, testCase.Expected, users)
		})
	}
}

func testUserStoreSearchInChannel(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
	}
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim-bobby" + model.NewId(),
		Email:    MakeEmail(),
	}
	_, err = ss.User().Save(u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
	}
	_, err = ss.User().Save(u3)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	tid := model.NewId()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u1.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u2.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u3.Id}, -1)
	require.Nil(t, err)

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData

	ch1 := model.Channel{
		TeamId:      tid,
		DisplayName: "NameName",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	c1, err := ss.Channel().Save(&ch1, -1)
	require.Nil(t, err)

	ch2 := model.Channel{
		TeamId:      tid,
		DisplayName: "NameName",
		Name:        "zz" + model.NewId() + "b",
		Type:        model.CHANNEL_OPEN,
	}
	c2, err := ss.Channel().Save(&ch2, -1)
	require.Nil(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c2.Id,
		UserId:      u2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   c1.Id,
		UserId:      u3.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	testCases := []struct {
		Description string
		ChannelId   string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"search jimb, channel 1",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, allow inactive, channel 1",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1, u3},
		},
		{
			"search jimb, allow inactive, channel 1, limit 1",
			c1.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          1,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, channel 2",
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, allow inactive, channel 2",
			c2.Id,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			users, err := ss.User().SearchInChannel(
				testCase.ChannelId,
				testCase.Term,
				testCase.Options,
			)
			require.Nil(t, err)
			assertUsers(t, testCase.Expected, users)
		})
	}
}

func testUserStoreSearchNotInTeam(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
	}
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim-bobby" + model.NewId(),
		Email:    MakeEmail(),
	}
	_, err = ss.User().Save(u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
	}
	_, err = ss.User().Save(u3)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	u4 := &model.User{
		Username: "simon" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 0,
	}
	_, err = ss.User().Save(u4)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u4.Id)) }()

	u5 := &model.User{
		Username:  "yu" + model.NewId(),
		FirstName: "En",
		LastName:  "Yu",
		Nickname:  "enyu",
		Email:     MakeEmail(),
	}
	_, err = ss.User().Save(u5)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u5.Id)) }()

	u6 := &model.User{
		Username:  "underscore" + model.NewId(),
		FirstName: "Du_",
		LastName:  "_DE",
		Nickname:  "lodash",
		Email:     MakeEmail(),
	}
	_, err = ss.User().Save(u6)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u6.Id)) }()

	teamId1 := model.NewId()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u1.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u2.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u3.Id}, -1)
	require.Nil(t, err)
	// u4 is not in team 1
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u5.Id}, -1)
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: u6.Id}, -1)
	require.Nil(t, err)

	teamId2 := model.NewId()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId2, UserId: u4.Id}, -1)
	require.Nil(t, err)

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData
	u4.AuthData = nilAuthData
	u5.AuthData = nilAuthData
	u6.AuthData = nilAuthData

	testCases := []struct {
		Description string
		TeamId      string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"search simo, team 1",
			teamId1,
			"simo",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u4},
		},

		{
			"search jimb, team 1",
			teamId1,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, allow inactive, team 1",
			teamId1,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search simo, team 2",
			teamId2,
			"simo",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{},
		},
		{
			"search jimb, team2",
			teamId2,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1},
		},
		{
			"search jimb, allow inactive, team 2",
			teamId2,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u1, u3},
		},
		{
			"search jimb, allow inactive, team 2, limit 1",
			teamId2,
			"jimb",
			&model.UserSearchOptions{
				AllowFullNames: true,
				AllowInactive:  true,
				Limit:          1,
			},
			[]*model.User{u1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			users, err := ss.User().SearchNotInTeam(
				testCase.TeamId,
				testCase.Term,
				testCase.Options,
			)
			require.Nil(t, err)
			assertUsers(t, testCase.Expected, users)
		})
	}
}

func testUserStoreSearchWithoutTeam(t *testing.T, ss store.Store) {
	u1 := &model.User{
		Username:  "jimbo1" + model.NewId(),
		FirstName: "Tim",
		LastName:  "Bill",
		Nickname:  "Rob",
		Email:     "harold" + model.NewId() + "@simulator.amazonses.com",
	}
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	u2 := &model.User{
		Username: "jim2-bobby" + model.NewId(),
		Email:    MakeEmail(),
	}
	_, err = ss.User().Save(u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	u3 := &model.User{
		Username: "jimbo3" + model.NewId(),
		Email:    MakeEmail(),
		DeleteAt: 1,
	}
	_, err = ss.User().Save(u3)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	tid := model.NewId()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: tid, UserId: u3.Id}, -1)
	require.Nil(t, err)

	// The users returned from the database will have AuthData as an empty string.
	nilAuthData := new(string)
	*nilAuthData = ""

	u1.AuthData = nilAuthData
	u2.AuthData = nilAuthData
	u3.AuthData = nilAuthData

	testCases := []struct {
		Description string
		Term        string
		Options     *model.UserSearchOptions
		Expected    []*model.User
	}{
		{
			"empty string",
			"",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2, u1},
		},
		{
			"jim",
			"jim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2, u1},
		},
		{
			"PLT-8354",
			"* ",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          model.USER_SEARCH_DEFAULT_LIMIT,
			},
			[]*model.User{u2, u1},
		},
		{
			"jim, limit 1",
			"jim",
			&model.UserSearchOptions{
				AllowFullNames: true,
				Limit:          1,
			},
			[]*model.User{u2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			users, err := ss.User().SearchWithoutTeam(
				testCase.Term,
				testCase.Options,
			)
			require.Nil(t, err)
			assertUsers(t, testCase.Expected, users)
		})
	}
}

func testCount(t *testing.T, ss store.Store) {
	// Regular
	teamId := model.NewId()
	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	// Deleted
	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.DeleteAt = model.GetMillis()
	_, err = ss.User().Save(u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	// Bot
	u3, err := ss.User().Save(&model.User{
		Email: MakeEmail(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	count, err := ss.User().Count(model.UserCountOptions{
		IncludeBotAccounts: false,
		IncludeDeleted:     false,
		TeamId:             "",
	})
	require.Nil(t, err)
	require.Equal(t, int64(1), count)

	count, err = ss.User().Count(model.UserCountOptions{
		IncludeBotAccounts: true,
		IncludeDeleted:     false,
		TeamId:             "",
	})
	require.Nil(t, err)
	require.Equal(t, int64(2), count)

	count, err = ss.User().Count(model.UserCountOptions{
		IncludeBotAccounts: false,
		IncludeDeleted:     true,
		TeamId:             "",
	})
	require.Nil(t, err)
	require.Equal(t, int64(2), count)

	count, err = ss.User().Count(model.UserCountOptions{
		IncludeBotAccounts: true,
		IncludeDeleted:     true,
		TeamId:             "",
	})
	require.Nil(t, err)
	require.Equal(t, int64(3), count)

	count, err = ss.User().Count(model.UserCountOptions{
		IncludeBotAccounts:  true,
		IncludeDeleted:      true,
		ExcludeRegularUsers: true,
		TeamId:              "",
	})
	require.Nil(t, err)
	require.Equal(t, int64(1), count)

	count, err = ss.User().Count(model.UserCountOptions{
		IncludeBotAccounts: true,
		IncludeDeleted:     true,
		TeamId:             teamId,
	})
	require.Nil(t, err)
	require.Equal(t, int64(1), count)

	count, err = ss.User().Count(model.UserCountOptions{
		IncludeBotAccounts: true,
		IncludeDeleted:     true,
		TeamId:             model.NewId(),
	})
	require.Nil(t, err)
	require.Equal(t, int64(0), count)

	count, err = ss.User().Count(model.UserCountOptions{
		IncludeBotAccounts: true,
		IncludeDeleted:     true,
		TeamId:             teamId,
		ViewRestrictions:   &model.ViewUsersRestrictions{Teams: []string{teamId}},
	})
	require.Nil(t, err)
	require.Equal(t, int64(1), count)

	count, err = ss.User().Count(model.UserCountOptions{
		IncludeBotAccounts: true,
		IncludeDeleted:     true,
		TeamId:             teamId,
		ViewRestrictions:   &model.ViewUsersRestrictions{Teams: []string{model.NewId()}},
	})
	require.Nil(t, err)
	require.Equal(t, int64(0), count)
}

func testUserStoreAnalyticsActiveCount(t *testing.T, ss store.Store, s SqlSupplier) {

	cleanupStatusStore(t, s)

	// Create 5 users statuses u0, u1, u2, u3, u4.
	// u4 is also a bot
	u0Id := model.NewId()
	u1Id := model.NewId()
	u2Id := model.NewId()
	u3Id := model.NewId()

	u4, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u4.Id)) }()

	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u4.Id,
		Username: u4.Username,
		OwnerId:  u1Id,
	})
	require.Nil(t, err)

	millis := model.GetMillis()
	millisTwoDaysAgo := model.GetMillis() - (2 * DAY_MILLISECONDS)
	millisTwoMonthsAgo := model.GetMillis() - (2 * MONTH_MILLISECONDS)

	// u0 last activity status is two months ago.
	// u1 last activity status is two days ago.
	// u2, u3, u4 last activity is within last day
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{UserId: u0Id, Status: model.STATUS_OFFLINE, LastActivityAt: millisTwoMonthsAgo}))
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{UserId: u1Id, Status: model.STATUS_OFFLINE, LastActivityAt: millisTwoDaysAgo}))
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{UserId: u2Id, Status: model.STATUS_OFFLINE, LastActivityAt: millis}))
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{UserId: u3Id, Status: model.STATUS_OFFLINE, LastActivityAt: millis}))
	require.Nil(t, ss.Status().SaveOrUpdate(&model.Status{UserId: u4.Id, Status: model.STATUS_OFFLINE, LastActivityAt: millis}))

	// Daily counts (without bots)
	count, err := ss.User().AnalyticsActiveCount(DAY_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: false})
	require.Nil(t, err)
	assert.Equal(t, int64(2), count)

	// Daily counts (with bots)
	count, err = ss.User().AnalyticsActiveCount(DAY_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: true})
	require.Nil(t, err)
	assert.Equal(t, int64(3), count)

	// Monthly counts (without bots)
	count, err = ss.User().AnalyticsActiveCount(MONTH_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: false})
	require.Nil(t, err)
	assert.Equal(t, int64(3), count)

	// Monthly counts - (with bots)
	count, err = ss.User().AnalyticsActiveCount(MONTH_MILLISECONDS, model.UserCountOptions{IncludeBotAccounts: true})
	require.Nil(t, err)
	assert.Equal(t, int64(4), count)
}

func testUserStoreAnalyticsGetInactiveUsersCount(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	count, err := ss.User().AnalyticsGetInactiveUsersCount()
	require.Nil(t, err)

	u2 := &model.User{}
	u2.Email = MakeEmail()
	u2.DeleteAt = model.GetMillis()
	_, err = ss.User().Save(u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	newCount, err := ss.User().AnalyticsGetInactiveUsersCount()
	require.Nil(t, err)
	require.Equal(t, count, newCount-1, "Expected 1 more inactive users but found otherwise.")
}

func testUserStoreAnalyticsGetSystemAdminCount(t *testing.T, ss store.Store) {
	countBefore, err := ss.User().AnalyticsGetSystemAdminCount()
	require.Nil(t, err)

	u1 := model.User{}
	u1.Email = MakeEmail()
	u1.Username = model.NewId()
	u1.Roles = "system_user system_admin"

	u2 := model.User{}
	u2.Email = MakeEmail()
	u2.Username = model.NewId()

	_, err = ss.User().Save(&u1)
	require.Nil(t, err, "couldn't save user")
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	_, err = ss.User().Save(&u2)
	require.Nil(t, err, "couldn't save user")

	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()

	result, err := ss.User().AnalyticsGetSystemAdminCount()
	require.Nil(t, err)
	require.Equal(t, countBefore+1, result, "Did not get the expected number of system admins.")

}

func testUserStoreGetProfilesNotInTeam(t *testing.T, ss store.Store) {
	team, err := ss.Team().Save(&model.Team{
		DisplayName: "Team",
		Name:        model.NewId(),
		Type:        model.TEAM_OPEN,
	})
	require.Nil(t, err)

	teamId := team.Id
	teamId2 := model.NewId()

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u1" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u1.Id}, -1)
	require.Nil(t, err)

	// Ensure update at timestamp changes
	time.Sleep(time.Millisecond * 10)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId2, UserId: u2.Id}, -1)
	require.Nil(t, err)

	// Ensure update at timestamp changes
	time.Sleep(time.Millisecond * 10)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u3" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u3.Id,
		Username: u3.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u3.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u3.Id)) }()

	var etag1, etag2, etag3 string

	t.Run("etag for profiles not in team 1", func(t *testing.T) {
		etag1 = ss.User().GetEtagForProfilesNotInTeam(teamId)
	})

	t.Run("get not in team 1, offset 0, limit 100000", func(t *testing.T) {
		users, userErr := ss.User().GetProfilesNotInTeam(teamId, false, 0, 100000, nil)
		require.Nil(t, userErr)
		assert.Equal(t, []*model.User{
			sanitized(u2),
			sanitized(u3),
		}, users)
	})

	t.Run("get not in team 1, offset 1, limit 1", func(t *testing.T) {
		users, userErr := ss.User().GetProfilesNotInTeam(teamId, false, 1, 1, nil)
		require.Nil(t, userErr)
		assert.Equal(t, []*model.User{
			sanitized(u3),
		}, users)
	})

	t.Run("get not in team 2, offset 0, limit 100", func(t *testing.T) {
		users, userErr := ss.User().GetProfilesNotInTeam(teamId2, false, 0, 100, nil)
		require.Nil(t, userErr)
		assert.Equal(t, []*model.User{
			sanitized(u1),
			sanitized(u3),
		}, users)
	})

	// Ensure update at timestamp changes
	time.Sleep(time.Millisecond * 10)

	// Add u2 to team 1
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u2.Id}, -1)
	require.Nil(t, err)
	u2.UpdateAt, err = ss.User().UpdateUpdateAt(u2.Id)
	require.Nil(t, err)

	t.Run("etag for profiles not in team 1 after update", func(t *testing.T) {
		etag2 = ss.User().GetEtagForProfilesNotInTeam(teamId)
		require.NotEqual(t, etag2, etag1, "etag should have changed")
	})

	t.Run("get not in team 1, offset 0, limit 100000 after update", func(t *testing.T) {
		users, userErr := ss.User().GetProfilesNotInTeam(teamId, false, 0, 100000, nil)
		require.Nil(t, userErr)
		assert.Equal(t, []*model.User{
			sanitized(u3),
		}, users)
	})

	// Ensure update at timestamp changes
	time.Sleep(time.Millisecond * 10)

	e := ss.Team().RemoveMember(teamId, u1.Id)
	require.Nil(t, e)
	e = ss.Team().RemoveMember(teamId, u2.Id)
	require.Nil(t, e)

	u1.UpdateAt, err = ss.User().UpdateUpdateAt(u1.Id)
	require.Nil(t, err)
	u2.UpdateAt, err = ss.User().UpdateUpdateAt(u2.Id)
	require.Nil(t, err)

	t.Run("etag for profiles not in team 1 after second update", func(t *testing.T) {
		etag3 = ss.User().GetEtagForProfilesNotInTeam(teamId)
		require.NotEqual(t, etag1, etag3, "etag should have changed")
		require.NotEqual(t, etag2, etag3, "etag should have changed")
	})

	t.Run("get not in team 1, offset 0, limit 100000 after second update", func(t *testing.T) {
		users, userErr := ss.User().GetProfilesNotInTeam(teamId, false, 0, 100000, nil)
		require.Nil(t, userErr)
		assert.Equal(t, []*model.User{
			sanitized(u1),
			sanitized(u2),
			sanitized(u3),
		}, users)
	})

	// Ensure update at timestamp changes
	time.Sleep(time.Millisecond * 10)

	u4, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u4" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u4.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: u4.Id}, -1)
	require.Nil(t, err)

	t.Run("etag for profiles not in team 1 after addition to team", func(t *testing.T) {
		etag4 := ss.User().GetEtagForProfilesNotInTeam(teamId)
		require.Equal(t, etag3, etag4, "etag should not have changed")
	})

	// Add u3 to team 2
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId2, UserId: u3.Id}, -1)
	require.Nil(t, err)
	u3.UpdateAt, err = ss.User().UpdateUpdateAt(u3.Id)
	require.Nil(t, err)

	// GetEtagForProfilesNotInTeam produces a new etag every time a member, not
	// in the team, gets a new UpdateAt value. In the case that an older member
	// in the set joins a different team, their UpdateAt value changes, thus
	// creating a new etag (even though the user set doesn't change). A hashing
	// solution, which only uses UserIds, would solve this issue.
	t.Run("etag for profiles not in team 1 after u3 added to team 2", func(t *testing.T) {
		t.Skip()
		etag4 := ss.User().GetEtagForProfilesNotInTeam(teamId)
		require.Equal(t, etag3, etag4, "etag should not have changed")
	})

	t.Run("get not in team 1, offset 0, limit 100000 after second update, setting group constrained when it's not", func(t *testing.T) {
		users, userErr := ss.User().GetProfilesNotInTeam(teamId, true, 0, 100000, nil)
		require.Nil(t, userErr)
		assert.Empty(t, users)
	})

	// create a group
	group, err := ss.Group().Create(&model.Group{
		Name:        "n_" + model.NewId(),
		DisplayName: "dn_" + model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    "ri_" + model.NewId(),
	})
	require.Nil(t, err)

	// add two members to the group
	for _, u := range []*model.User{u1, u2} {
		_, err = ss.Group().UpsertMember(group.Id, u.Id)
		require.Nil(t, err)
	}

	// associate the group with the team
	_, err = ss.Group().CreateGroupSyncable(&model.GroupSyncable{
		GroupId:    group.Id,
		SyncableId: teamId,
		Type:       model.GroupSyncableTypeTeam,
	})
	require.Nil(t, err)

	t.Run("get not in team 1, offset 0, limit 100000 after second update, setting group constrained", func(t *testing.T) {
		users, userErr := ss.User().GetProfilesNotInTeam(teamId, true, 0, 100000, nil)
		require.Nil(t, userErr)
		assert.Equal(t, []*model.User{
			sanitized(u1),
			sanitized(u2),
		}, users)
	})
}

func testUserStoreClearAllCustomRoleAssignments(t *testing.T, ss store.Store) {
	u1 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "system_user system_admin system_post_all",
	}
	u2 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "system_user custom_role system_admin another_custom_role",
	}
	u3 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "system_user",
	}
	u4 := model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "custom_only",
	}

	_, err := ss.User().Save(&u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.User().Save(&u2)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.User().Save(&u3)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u3.Id)) }()
	_, err = ss.User().Save(&u4)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u4.Id)) }()

	require.Nil(t, ss.User().ClearAllCustomRoleAssignments())

	r1, err := ss.User().GetByUsername(u1.Username)
	require.Nil(t, err)
	assert.Equal(t, u1.Roles, r1.Roles)

	r2, err1 := ss.User().GetByUsername(u2.Username)
	require.Nil(t, err1)
	assert.Equal(t, "system_user system_admin", r2.Roles)

	r3, err2 := ss.User().GetByUsername(u3.Username)
	require.Nil(t, err2)
	assert.Equal(t, u3.Roles, r3.Roles)

	r4, err3 := ss.User().GetByUsername(u4.Username)
	require.Nil(t, err3)
	assert.Equal(t, "", r4.Roles)
}

func testUserStoreGetAllAfter(t *testing.T, ss store.Store) {
	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		Roles:    "system_user system_admin system_post_all",
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: "u2" + model.NewId(),
	})
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u2.Id)) }()
	_, err = ss.Bot().Save(&model.Bot{
		UserId:   u2.Id,
		Username: u2.Username,
		OwnerId:  u1.Id,
	})
	require.Nil(t, err)
	u2.IsBot = true
	defer func() { require.Nil(t, ss.Bot().PermanentDelete(u2.Id)) }()

	expected := []*model.User{u1, u2}
	if strings.Compare(u2.Id, u1.Id) < 0 {
		expected = []*model.User{u2, u1}
	}

	t.Run("get after lowest possible id", func(t *testing.T) {
		actual, err := ss.User().GetAllAfter(10000, strings.Repeat("0", 26))
		require.Nil(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("get after first user", func(t *testing.T) {
		actual, err := ss.User().GetAllAfter(10000, expected[0].Id)
		require.Nil(t, err)

		assert.Equal(t, []*model.User{expected[1]}, actual)
	})

	t.Run("get after second user", func(t *testing.T) {
		actual, err := ss.User().GetAllAfter(10000, expected[1].Id)
		require.Nil(t, err)

		assert.Equal(t, []*model.User{}, actual)
	})
}

func testUserStoreGetUsersBatchForIndexing(t *testing.T, ss store.Store) {
	// Set up all the objects needed
	t1, err := ss.Team().Save(&model.Team{
		DisplayName: "Team1",
		Name:        model.NewId(),
		Type:        model.TEAM_OPEN,
	})
	require.Nil(t, err)

	ch1 := &model.Channel{
		Name: model.NewId(),
		Type: model.CHANNEL_OPEN,
	}
	cPub1, err := ss.Channel().Save(ch1, -1)
	require.Nil(t, err)

	ch2 := &model.Channel{
		Name: model.NewId(),
		Type: model.CHANNEL_OPEN,
	}
	cPub2, err := ss.Channel().Save(ch2, -1)
	require.Nil(t, err)

	ch3 := &model.Channel{
		Name: model.NewId(),
		Type: model.CHANNEL_PRIVATE,
	}

	cPriv, err := ss.Channel().Save(ch3, -1)
	require.Nil(t, err)

	u1, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		CreateAt: model.GetMillis(),
	})
	require.Nil(t, err)

	time.Sleep(10 * time.Millisecond)

	u2, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		CreateAt: model.GetMillis(),
	})
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{
		UserId: u2.Id,
		TeamId: t1.Id,
	}, 100)
	require.Nil(t, err)
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		UserId:      u2.Id,
		ChannelId:   cPub1.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		UserId:      u2.Id,
		ChannelId:   cPub2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	startTime := u2.CreateAt
	time.Sleep(10 * time.Millisecond)

	u3, err := ss.User().Save(&model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
		CreateAt: model.GetMillis(),
	})
	require.Nil(t, err)
	_, err = ss.Team().SaveMember(&model.TeamMember{
		UserId:   u3.Id,
		TeamId:   t1.Id,
		DeleteAt: model.GetMillis(),
	}, 100)
	require.Nil(t, err)
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		UserId:      u3.Id,
		ChannelId:   cPub2.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		UserId:      u3.Id,
		ChannelId:   cPriv.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	endTime := u3.CreateAt

	// First and last user should be outside the range
	res1List, err := ss.User().GetUsersBatchForIndexing(startTime, endTime, 100)
	assert.Nil(t, err)

	assert.Len(t, res1List, 1)
	assert.Equal(t, res1List[0].Username, u2.Username)
	assert.ElementsMatch(t, res1List[0].TeamsIds, []string{t1.Id})
	assert.ElementsMatch(t, res1List[0].ChannelsIds, []string{cPub1.Id, cPub2.Id})

	// Update startTime to include first user
	startTime = u1.CreateAt
	res2List, err := ss.User().GetUsersBatchForIndexing(startTime, endTime, 100)
	assert.Nil(t, err)

	assert.Len(t, res2List, 2)
	assert.Equal(t, res2List[0].Username, u1.Username)
	assert.Equal(t, res2List[0].ChannelsIds, []string{})
	assert.Equal(t, res2List[0].TeamsIds, []string{})
	assert.Equal(t, res2List[1].Username, u2.Username)

	// Update endTime to include last user
	endTime = model.GetMillis()
	res3List, err := ss.User().GetUsersBatchForIndexing(startTime, endTime, 100)
	assert.Nil(t, err)

	assert.Len(t, res3List, 3)
	assert.Equal(t, res3List[0].Username, u1.Username)
	assert.Equal(t, res3List[1].Username, u2.Username)
	assert.Equal(t, res3List[2].Username, u3.Username)
	assert.ElementsMatch(t, res3List[2].TeamsIds, []string{})
	assert.ElementsMatch(t, res3List[2].ChannelsIds, []string{cPub2.Id})

	// Testing the limit
	res4List, err := ss.User().GetUsersBatchForIndexing(startTime, endTime, 2)
	assert.Nil(t, err)

	assert.Len(t, res4List, 2)
	assert.Equal(t, res4List[0].Username, u1.Username)
	assert.Equal(t, res4List[1].Username, u2.Username)
}

func testUserStoreGetTeamGroupUsers(t *testing.T, ss store.Store) {
	// create team
	id := model.NewId()
	team, err := ss.Team().Save(&model.Team{
		DisplayName: "dn_" + id,
		Name:        "n-" + id,
		Email:       id + "@test.com",
		Type:        model.TEAM_INVITE,
	})
	require.Nil(t, err)
	require.NotNil(t, team)

	// create users
	var testUsers []*model.User
	for i := 0; i < 3; i++ {
		id = model.NewId()
		user, userErr := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
		})
		require.Nil(t, userErr)
		require.NotNil(t, user)
		testUsers = append(testUsers, user)
	}
	require.Len(t, testUsers, 3, "testUsers length doesn't meet required length")
	userGroupA, userGroupB, userNoGroup := testUsers[0], testUsers[1], testUsers[2]

	// add non-group-member to the team (to prove that the query isn't just returning all members)
	_, err = ss.Team().SaveMember(&model.TeamMember{
		TeamId: team.Id,
		UserId: userNoGroup.Id,
	}, 999)
	require.Nil(t, err)

	// create groups
	var testGroups []*model.Group
	for i := 0; i < 2; i++ {
		id = model.NewId()

		var group *model.Group
		group, err = ss.Group().Create(&model.Group{
			Name:        "n_" + id,
			DisplayName: "dn_" + id,
			Source:      model.GroupSourceLdap,
			RemoteId:    "ri_" + id,
		})
		require.Nil(t, err)
		require.NotNil(t, group)
		testGroups = append(testGroups, group)
	}
	require.Len(t, testGroups, 2, "testGroups length doesn't meet required length")
	groupA, groupB := testGroups[0], testGroups[1]

	// add members to groups
	_, err = ss.Group().UpsertMember(groupA.Id, userGroupA.Id)
	require.Nil(t, err)
	_, err = ss.Group().UpsertMember(groupB.Id, userGroupB.Id)
	require.Nil(t, err)

	// association one group to team
	_, err = ss.Group().CreateGroupSyncable(&model.GroupSyncable{
		GroupId:    groupA.Id,
		SyncableId: team.Id,
		Type:       model.GroupSyncableTypeTeam,
	})
	require.Nil(t, err)

	var users []*model.User

	requireNUsers := func(n int) {
		users, err = ss.User().GetTeamGroupUsers(team.Id)
		require.Nil(t, err)
		require.NotNil(t, users)
		require.Len(t, users, n)
	}

	// team not group constrained returns users
	requireNUsers(1)

	// update team to be group-constrained
	team.GroupConstrained = model.NewBool(true)
	team, err = ss.Team().Update(team)
	require.Nil(t, err)

	// still returns user (being group-constrained has no effect)
	requireNUsers(1)

	// associate other group to team
	_, err = ss.Group().CreateGroupSyncable(&model.GroupSyncable{
		GroupId:    groupB.Id,
		SyncableId: team.Id,
		Type:       model.GroupSyncableTypeTeam,
	})
	require.Nil(t, err)

	// should return users from all groups
	// 2 users now that both groups have been associated to the team
	requireNUsers(2)

	// add team membership of allowed user
	_, err = ss.Team().SaveMember(&model.TeamMember{
		TeamId: team.Id,
		UserId: userGroupA.Id,
	}, 999)
	require.Nil(t, err)

	// ensure allowed member still returned by query
	requireNUsers(2)

	// delete team membership of allowed user
	err = ss.Team().RemoveMember(team.Id, userGroupA.Id)
	require.Nil(t, err)

	// ensure removed allowed member still returned by query
	requireNUsers(2)
}

func testUserStoreGetChannelGroupUsers(t *testing.T, ss store.Store) {
	// create channel
	id := model.NewId()
	channel, err := ss.Channel().Save(&model.Channel{
		DisplayName: "dn_" + id,
		Name:        "n-" + id,
		Type:        model.CHANNEL_PRIVATE,
	}, 999)
	require.Nil(t, err)
	require.NotNil(t, channel)

	// create users
	var testUsers []*model.User
	for i := 0; i < 3; i++ {
		id = model.NewId()
		user, userErr := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
		})
		require.Nil(t, userErr)
		require.NotNil(t, user)
		testUsers = append(testUsers, user)
	}
	require.Len(t, testUsers, 3, "testUsers length doesn't meet required length")
	userGroupA, userGroupB, userNoGroup := testUsers[0], testUsers[1], testUsers[2]

	// add non-group-member to the channel (to prove that the query isn't just returning all members)
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   channel.Id,
		UserId:      userNoGroup.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	// create groups
	var testGroups []*model.Group
	for i := 0; i < 2; i++ {
		id = model.NewId()
		var group *model.Group
		group, err = ss.Group().Create(&model.Group{
			Name:        "n_" + id,
			DisplayName: "dn_" + id,
			Source:      model.GroupSourceLdap,
			RemoteId:    "ri_" + id,
		})
		require.Nil(t, err)
		require.NotNil(t, group)
		testGroups = append(testGroups, group)
	}
	require.Len(t, testGroups, 2, "testGroups length doesn't meet required length")
	groupA, groupB := testGroups[0], testGroups[1]

	// add members to groups
	_, err = ss.Group().UpsertMember(groupA.Id, userGroupA.Id)
	require.Nil(t, err)
	_, err = ss.Group().UpsertMember(groupB.Id, userGroupB.Id)
	require.Nil(t, err)

	// association one group to channel
	_, err = ss.Group().CreateGroupSyncable(&model.GroupSyncable{
		GroupId:    groupA.Id,
		SyncableId: channel.Id,
		Type:       model.GroupSyncableTypeChannel,
	})
	require.Nil(t, err)

	var users []*model.User

	requireNUsers := func(n int) {
		users, err = ss.User().GetChannelGroupUsers(channel.Id)
		require.Nil(t, err)
		require.NotNil(t, users)
		require.Len(t, users, n)
	}

	// channel not group constrained returns users
	requireNUsers(1)

	// update team to be group-constrained
	channel.GroupConstrained = model.NewBool(true)
	_, err = ss.Channel().Update(channel)
	require.Nil(t, err)

	// still returns user (being group-constrained has no effect)
	requireNUsers(1)

	// associate other group to team
	_, err = ss.Group().CreateGroupSyncable(&model.GroupSyncable{
		GroupId:    groupB.Id,
		SyncableId: channel.Id,
		Type:       model.GroupSyncableTypeChannel,
	})
	require.Nil(t, err)

	// should return users from all groups
	// 2 users now that both groups have been associated to the team
	requireNUsers(2)

	// add team membership of allowed user
	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   channel.Id,
		UserId:      userGroupA.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.Nil(t, err)

	// ensure allowed member still returned by query
	requireNUsers(2)

	// delete team membership of allowed user
	err = ss.Channel().RemoveMember(channel.Id, userGroupA.Id)
	require.Nil(t, err)

	// ensure removed allowed member still returned by query
	requireNUsers(2)
}

func testUserStorePromoteGuestToUser(t *testing.T, ss store.Store) {
	// create users
	t.Run("Must do nothing with regular user", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_user",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: true, SchemeUser: false}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)
		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user.Id, SchemeGuest: true, SchemeUser: false, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().PromoteGuestToUser(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_user", updatedUser.Roles)
		require.True(t, user.UpdateAt < updatedUser.UpdateAt)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.False(t, updatedTeamMember.SchemeGuest)
		require.True(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user.Id)
		require.Nil(t, err)
		require.False(t, updatedChannelMember.SchemeGuest)
		require.True(t, updatedChannelMember.SchemeUser)
	})

	t.Run("Must do nothing with admin user", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_user system_admin",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: true, SchemeUser: false}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)
		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user.Id, SchemeGuest: true, SchemeUser: false, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().PromoteGuestToUser(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_user system_admin", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.False(t, updatedTeamMember.SchemeGuest)
		require.True(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user.Id)
		require.Nil(t, err)
		require.False(t, updatedChannelMember.SchemeGuest)
		require.True(t, updatedChannelMember.SchemeUser)
	})

	t.Run("Must work with guest user without teams or channels", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_guest",
		})
		require.Nil(t, err)

		err = ss.User().PromoteGuestToUser(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_user", updatedUser.Roles)
	})

	t.Run("Must work with guest user with teams but no channels", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_guest",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: true, SchemeUser: false}, 999)
		require.Nil(t, err)

		err = ss.User().PromoteGuestToUser(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_user", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.False(t, updatedTeamMember.SchemeGuest)
		require.True(t, updatedTeamMember.SchemeUser)
	})

	t.Run("Must work with guest user with teams and channels", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_guest",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: true, SchemeUser: false}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)
		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user.Id, SchemeGuest: true, SchemeUser: false, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().PromoteGuestToUser(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_user", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.False(t, updatedTeamMember.SchemeGuest)
		require.True(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user.Id)
		require.Nil(t, err)
		require.False(t, updatedChannelMember.SchemeGuest)
		require.True(t, updatedChannelMember.SchemeUser)
	})

	t.Run("Must work with guest user with teams and channels and custom role", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_guest custom_role",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: true, SchemeUser: false}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)
		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user.Id, SchemeGuest: true, SchemeUser: false, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().PromoteGuestToUser(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_user custom_role", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.False(t, updatedTeamMember.SchemeGuest)
		require.True(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user.Id)
		require.Nil(t, err)
		require.False(t, updatedChannelMember.SchemeGuest)
		require.True(t, updatedChannelMember.SchemeUser)
	})

	t.Run("Must no change any other user guest role", func(t *testing.T) {
		id := model.NewId()
		user1, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_guest",
		})
		require.Nil(t, err)

		teamId1 := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: user1.Id, SchemeGuest: true, SchemeUser: false}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId1,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)

		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user1.Id, SchemeGuest: true, SchemeUser: false, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		id = model.NewId()
		user2, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_guest",
		})
		require.Nil(t, err)

		teamId2 := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId2, UserId: user2.Id, SchemeGuest: true, SchemeUser: false}, 999)
		require.Nil(t, err)

		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user2.Id, SchemeGuest: true, SchemeUser: false, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().PromoteGuestToUser(user1.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user1.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_user", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId1, user1.Id)
		require.Nil(t, err)
		require.False(t, updatedTeamMember.SchemeGuest)
		require.True(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user1.Id)
		require.Nil(t, err)
		require.False(t, updatedChannelMember.SchemeGuest)
		require.True(t, updatedChannelMember.SchemeUser)

		notUpdatedUser, err := ss.User().Get(user2.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_guest", notUpdatedUser.Roles)

		notUpdatedTeamMember, err := ss.Team().GetMember(teamId2, user2.Id)
		require.Nil(t, err)
		require.True(t, notUpdatedTeamMember.SchemeGuest)
		require.False(t, notUpdatedTeamMember.SchemeUser)

		notUpdatedChannelMember, err := ss.Channel().GetMember(channel.Id, user2.Id)
		require.Nil(t, err)
		require.True(t, notUpdatedChannelMember.SchemeGuest)
		require.False(t, notUpdatedChannelMember.SchemeUser)
	})
}

func testUserStoreDemoteUserToGuest(t *testing.T, ss store.Store) {
	// create users
	t.Run("Must do nothing with guest", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_guest",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: false, SchemeUser: true}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)
		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user.Id, SchemeGuest: false, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().DemoteUserToGuest(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_guest", updatedUser.Roles)
		require.True(t, user.UpdateAt < updatedUser.UpdateAt)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.True(t, updatedTeamMember.SchemeGuest)
		require.False(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user.Id)
		require.Nil(t, err)
		require.True(t, updatedChannelMember.SchemeGuest)
		require.False(t, updatedChannelMember.SchemeUser)
	})

	t.Run("Must demote properly an admin user", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_user system_admin",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: true, SchemeUser: false}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)
		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user.Id, SchemeGuest: true, SchemeUser: false, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().DemoteUserToGuest(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_guest", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.True(t, updatedTeamMember.SchemeGuest)
		require.False(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user.Id)
		require.Nil(t, err)
		require.True(t, updatedChannelMember.SchemeGuest)
		require.False(t, updatedChannelMember.SchemeUser)
	})

	t.Run("Must work with user without teams or channels", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_user",
		})
		require.Nil(t, err)

		err = ss.User().DemoteUserToGuest(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_guest", updatedUser.Roles)
	})

	t.Run("Must work with user with teams but no channels", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_user",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: false, SchemeUser: true}, 999)
		require.Nil(t, err)

		err = ss.User().DemoteUserToGuest(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_guest", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.True(t, updatedTeamMember.SchemeGuest)
		require.False(t, updatedTeamMember.SchemeUser)
	})

	t.Run("Must work with user with teams and channels", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_user",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: false, SchemeUser: true}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)
		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user.Id, SchemeGuest: false, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().DemoteUserToGuest(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_guest", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.True(t, updatedTeamMember.SchemeGuest)
		require.False(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user.Id)
		assert.Nil(t, err)
		require.True(t, updatedChannelMember.SchemeGuest)
		require.False(t, updatedChannelMember.SchemeUser)
	})

	t.Run("Must work with user with teams and channels and custom role", func(t *testing.T) {
		id := model.NewId()
		user, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_user custom_role",
		})
		require.Nil(t, err)

		teamId := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId, UserId: user.Id, SchemeGuest: false, SchemeUser: true}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)
		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user.Id, SchemeGuest: false, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().DemoteUserToGuest(user.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_guest custom_role", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId, user.Id)
		require.Nil(t, err)
		require.True(t, updatedTeamMember.SchemeGuest)
		require.False(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user.Id)
		assert.Nil(t, err)
		require.True(t, updatedChannelMember.SchemeGuest)
		require.False(t, updatedChannelMember.SchemeUser)
	})

	t.Run("Must no change any other user role", func(t *testing.T) {
		id := model.NewId()
		user1, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_user",
		})
		require.Nil(t, err)

		teamId1 := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId1, UserId: user1.Id, SchemeGuest: false, SchemeUser: true}, 999)
		require.Nil(t, err)

		channel, err := ss.Channel().Save(&model.Channel{
			TeamId:      teamId1,
			DisplayName: "Channel name",
			Name:        "channel-" + model.NewId(),
			Type:        model.CHANNEL_OPEN,
		}, -1)
		require.Nil(t, err)

		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user1.Id, SchemeGuest: false, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		id = model.NewId()
		user2, err := ss.User().Save(&model.User{
			Email:     id + "@test.com",
			Username:  "un_" + id,
			Nickname:  "nn_" + id,
			FirstName: "f_" + id,
			LastName:  "l_" + id,
			Password:  "Password1",
			Roles:     "system_user",
		})
		require.Nil(t, err)

		teamId2 := model.NewId()
		_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: teamId2, UserId: user2.Id, SchemeGuest: false, SchemeUser: true}, 999)
		require.Nil(t, err)

		_, err = ss.Channel().SaveMember(&model.ChannelMember{ChannelId: channel.Id, UserId: user2.Id, SchemeGuest: false, SchemeUser: true, NotifyProps: model.GetDefaultChannelNotifyProps()})
		require.Nil(t, err)

		err = ss.User().DemoteUserToGuest(user1.Id)
		assert.Nil(t, err)
		updatedUser, err := ss.User().Get(user1.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_guest", updatedUser.Roles)

		updatedTeamMember, err := ss.Team().GetMember(teamId1, user1.Id)
		require.Nil(t, err)
		require.True(t, updatedTeamMember.SchemeGuest)
		require.False(t, updatedTeamMember.SchemeUser)

		updatedChannelMember, err := ss.Channel().GetMember(channel.Id, user1.Id)
		assert.Nil(t, err)
		require.True(t, updatedChannelMember.SchemeGuest)
		require.False(t, updatedChannelMember.SchemeUser)

		notUpdatedUser, err := ss.User().Get(user2.Id)
		assert.Nil(t, err)
		require.Equal(t, "system_user", notUpdatedUser.Roles)

		notUpdatedTeamMember, err := ss.Team().GetMember(teamId2, user2.Id)
		require.Nil(t, err)
		require.False(t, notUpdatedTeamMember.SchemeGuest)
		require.True(t, notUpdatedTeamMember.SchemeUser)

		notUpdatedChannelMember, err := ss.Channel().GetMember(channel.Id, user2.Id)
		assert.Nil(t, err)
		require.False(t, notUpdatedChannelMember.SchemeGuest)
		require.True(t, notUpdatedChannelMember.SchemeUser)
	})
}

func testUserStoreResetLastPictureUpdate(t *testing.T, ss store.Store) {
	u1 := &model.User{}
	u1.Email = MakeEmail()
	_, err := ss.User().Save(u1)
	require.Nil(t, err)
	defer func() { require.Nil(t, ss.User().PermanentDelete(u1.Id)) }()
	_, err = ss.Team().SaveMember(&model.TeamMember{TeamId: model.NewId(), UserId: u1.Id}, -1)
	require.Nil(t, err)

	err = ss.User().UpdateLastPictureUpdate(u1.Id)
	require.Nil(t, err)

	user, err := ss.User().Get(u1.Id)
	require.Nil(t, err)

	assert.NotZero(t, user.LastPictureUpdate)
	assert.NotZero(t, user.UpdateAt)

	err = ss.User().ResetLastPictureUpdate(u1.Id)
	require.Nil(t, err)

	user2, err := ss.User().Get(u1.Id)
	require.Nil(t, err)

	assert.True(t, user2.UpdateAt > user.UpdateAt)
	assert.Zero(t, user2.LastPictureUpdate)
}
