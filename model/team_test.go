// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeamJson(t *testing.T) {
	o := Team{Id: NewId(), DisplayName: NewId()}
	json := o.ToJson()
	ro := TeamFromJson(strings.NewReader(json))

	require.Equal(t, o.Id, ro.Id, "Ids do not match")
}

func TestTeamIsValid(t *testing.T) {
	o := Team{}

	err := o.IsValid()
	require.NotNil(t, err, "should be invalid")

	o.Id = NewId()
	err = o.IsValid()
	require.NotNil(t, err, "should be invalid")

	o.CreateAt = GetMillis()
	err = o.IsValid()
	require.NotNil(t, err, "should be invalid")

	o.UpdateAt = GetMillis()
	err = o.IsValid()
	require.NotNil(t, err, "should be invalid")

	o.Email = strings.Repeat("01234567890", 20)
	err = o.IsValid()
	require.NotNil(t, err, "should be invalid")

	o.Email = "corey+test@hulen.com"
	o.DisplayName = strings.Repeat("01234567890", 20)
	err = o.IsValid()
	require.NotNil(t, err, "should be invalid")

	o.DisplayName = "1234"
	o.Name = "ZZZZZZZ"
	err = o.IsValid()
	require.NotNil(t, err, "should be invalid")

	o.Name = "zzzzz"
	o.Type = TEAM_OPEN
	o.InviteId = NewId()
	err = o.IsValid()
	require.Nil(t, err, err)
}

func TestTeamPreSave(t *testing.T) {
	o := Team{DisplayName: "test"}
	o.PreSave()
	o.Etag()
}

func TestTeamPreUpdate(t *testing.T) {
	o := Team{DisplayName: "test"}
	o.PreUpdate()
}

var domains = []struct {
	value    string
	expected bool
}{
	{"spin-punch", true},
	{"-spin-punch", false},
	{"spin-punch-", false},
	{"spin_punch", false},
	{"a", false},
	{"aa", true},
	{"aaa", true},
	{"aaa-999b", true},
	{"b00b", true},
	{"b)", false},
	{"test", true},
}

func TestValidTeamName(t *testing.T) {
	for _, v := range domains {
		actual := IsValidTeamName(v.value)
		assert.Equal(t, actual, v.expected)
	}
}

var tReservedDomains = []struct {
	value    string
	expected bool
}{
	{"admin", true},
	{"Admin-punch", true},
	{"spin-punch-admin", false},
}

func TestReservedTeamName(t *testing.T) {
	for _, v := range tReservedDomains {
		actual := IsReservedTeamName(v.value)
		assert.Equal(t, actual, v.expected)
	}
}

func TestCleanTeamName(t *testing.T) {
	actual := CleanTeamName("Jimbo's Admin")
	require.Equal(t, actual, "jimbos-admin", "didn't clean name properly")

	actual = CleanTeamName("Admin Really cool")
	require.Equal(t, actual, "really-cool", "didn't clean name properly")

	actual = CleanTeamName("super-duper-guys")
	require.Equal(t, actual, "super-duper-guys", "didn't clean name properly")
}

func TestTeamPatch(t *testing.T) {
	p := &TeamPatch{
		DisplayName:      new(string),
		Description:      new(string),
		CompanyName:      new(string),
		AllowedDomains:   new(string),
		AllowOpenInvite:  new(bool),
		GroupConstrained: new(bool),
	}

	*p.DisplayName = NewId()
	*p.Description = NewId()
	*p.CompanyName = NewId()
	*p.AllowedDomains = NewId()
	*p.AllowOpenInvite = true
	*p.GroupConstrained = true

	o := Team{Id: NewId()}
	o.Patch(p)

	require.Equal(t, *p.DisplayName, o.DisplayName, "DisplayName did not update")
	require.Equal(t, *p.Description, o.Description, "Description did not update")
	require.Equal(t, *p.CompanyName, o.CompanyName, "CompanyName did not update")
	require.Equal(t, *p.AllowedDomains, o.AllowedDomains, "AllowedDomains did not update")
	require.Equal(t, *p.AllowOpenInvite, o.AllowOpenInvite, "AllowOpenInvite did not update")
	require.Equal(t, *p.GroupConstrained, *o.GroupConstrained)
}
