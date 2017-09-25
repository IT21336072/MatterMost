// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/assert"
)

func TestSqlComplianceStore(t *testing.T) {
	Setup()

	compliance1 := &model.Compliance{Desc: "Audit for federal subpoena case #22443", UserId: model.NewId(), Status: model.COMPLIANCE_STATUS_FAILED, StartAt: model.GetMillis() - 1, EndAt: model.GetMillis() + 1, Type: model.COMPLIANCE_TYPE_ADHOC}
	Must(store.Compliance().Save(compliance1))
	time.Sleep(100 * time.Millisecond)

	compliance2 := &model.Compliance{Desc: "Audit for federal subpoena case #11458", UserId: model.NewId(), Status: model.COMPLIANCE_STATUS_RUNNING, StartAt: model.GetMillis() - 1, EndAt: model.GetMillis() + 1, Type: model.COMPLIANCE_TYPE_ADHOC}
	Must(store.Compliance().Save(compliance2))
	time.Sleep(100 * time.Millisecond)

	c := store.Compliance().GetAll(0, 1000)
	result := <-c
	compliances := result.Data.(model.Compliances)

	if compliances[0].Status != model.COMPLIANCE_STATUS_RUNNING && compliance2.Id != compliances[0].Id {
		t.Fatal()
	}

	compliance2.Status = model.COMPLIANCE_STATUS_FAILED
	Must(store.Compliance().Update(compliance2))

	c = store.Compliance().GetAll(0, 1000)
	result = <-c
	compliances = result.Data.(model.Compliances)

	if compliances[0].Status != model.COMPLIANCE_STATUS_FAILED && compliance2.Id != compliances[0].Id {
		t.Fatal()
	}

	c = store.Compliance().GetAll(0, 1)
	result = <-c
	compliances = result.Data.(model.Compliances)

	if len(compliances) != 1 {
		t.Fatal("should only have returned 1")
	}

	c = store.Compliance().GetAll(1, 1)
	result = <-c
	compliances = result.Data.(model.Compliances)

	if len(compliances) != 1 {
		t.Fatal("should only have returned 1")
	}

	rc2 := (<-store.Compliance().Get(compliance2.Id)).Data.(*model.Compliance)
	if rc2.Status != compliance2.Status {
		t.Fatal()
	}
}

func TestComplianceExport(t *testing.T) {
	Setup()

	time.Sleep(100 * time.Millisecond)

	t1 := &model.Team{}
	t1.DisplayName = "DisplayName"
	t1.Name = "zz" + model.NewId() + "b"
	t1.Email = model.NewId() + "@nowhere.com"
	t1.Type = model.TEAM_OPEN
	t1 = Must(store.Team().Save(t1)).(*model.Team)

	u1 := &model.User{}
	u1.Email = model.NewId()
	u1.Username = model.NewId()
	u1 = Must(store.User().Save(u1)).(*model.User)
	Must(store.Team().SaveMember(&model.TeamMember{TeamId: t1.Id, UserId: u1.Id}))

	u2 := &model.User{}
	u2.Email = model.NewId()
	u2.Username = model.NewId()
	u2 = Must(store.User().Save(u2)).(*model.User)
	Must(store.Team().SaveMember(&model.TeamMember{TeamId: t1.Id, UserId: u2.Id}))

	c1 := &model.Channel{}
	c1.TeamId = t1.Id
	c1.DisplayName = "Channel2"
	c1.Name = "zz" + model.NewId() + "b"
	c1.Type = model.CHANNEL_OPEN
	c1 = Must(store.Channel().Save(c1)).(*model.Channel)

	o1 := &model.Post{}
	o1.ChannelId = c1.Id
	o1.UserId = u1.Id
	o1.CreateAt = model.GetMillis()
	o1.Message = "zz" + model.NewId() + "b"
	o1 = Must(store.Post().Save(o1)).(*model.Post)

	o1a := &model.Post{}
	o1a.ChannelId = c1.Id
	o1a.UserId = u1.Id
	o1a.CreateAt = o1.CreateAt + 10
	o1a.Message = "zz" + model.NewId() + "b"
	o1a = Must(store.Post().Save(o1a)).(*model.Post)

	o2 := &model.Post{}
	o2.ChannelId = c1.Id
	o2.UserId = u1.Id
	o2.CreateAt = o1.CreateAt + 20
	o2.Message = "zz" + model.NewId() + "b"
	o2 = Must(store.Post().Save(o2)).(*model.Post)

	o2a := &model.Post{}
	o2a.ChannelId = c1.Id
	o2a.UserId = u2.Id
	o2a.CreateAt = o1.CreateAt + 30
	o2a.Message = "zz" + model.NewId() + "b"
	o2a = Must(store.Post().Save(o2a)).(*model.Post)

	time.Sleep(100 * time.Millisecond)

	cr1 := &model.Compliance{Desc: "test" + model.NewId(), StartAt: o1.CreateAt - 1, EndAt: o2a.CreateAt + 1}
	if r1 := <-store.Compliance().ComplianceExport(cr1); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		cposts := r1.Data.([]*model.CompliancePost)

		if len(cposts) != 4 {
			t.Fatal("return wrong results length")
		}

		if cposts[0].PostId != o1.Id {
			t.Fatal("Wrong sort")
		}

		if cposts[3].PostId != o2a.Id {
			t.Fatal("Wrong sort")
		}
	}

	cr2 := &model.Compliance{Desc: "test" + model.NewId(), StartAt: o1.CreateAt - 1, EndAt: o2a.CreateAt + 1, Emails: u2.Email}
	if r1 := <-store.Compliance().ComplianceExport(cr2); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		cposts := r1.Data.([]*model.CompliancePost)

		if len(cposts) != 1 {
			t.Fatal("return wrong results length")
		}

		if cposts[0].PostId != o2a.Id {
			t.Fatal("Wrong sort")
		}
	}

	cr3 := &model.Compliance{Desc: "test" + model.NewId(), StartAt: o1.CreateAt - 1, EndAt: o2a.CreateAt + 1, Emails: u2.Email + ", " + u1.Email}
	if r1 := <-store.Compliance().ComplianceExport(cr3); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		cposts := r1.Data.([]*model.CompliancePost)

		if len(cposts) != 4 {
			t.Fatal("return wrong results length")
		}

		if cposts[0].PostId != o1.Id {
			t.Fatal("Wrong sort")
		}

		if cposts[3].PostId != o2a.Id {
			t.Fatal("Wrong sort")
		}
	}

	cr4 := &model.Compliance{Desc: "test" + model.NewId(), StartAt: o1.CreateAt - 1, EndAt: o2a.CreateAt + 1, Keywords: o2a.Message}
	if r1 := <-store.Compliance().ComplianceExport(cr4); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		cposts := r1.Data.([]*model.CompliancePost)

		if len(cposts) != 1 {
			t.Fatal("return wrong results length")
		}

		if cposts[0].PostId != o2a.Id {
			t.Fatal("Wrong sort")
		}
	}

	cr5 := &model.Compliance{Desc: "test" + model.NewId(), StartAt: o1.CreateAt - 1, EndAt: o2a.CreateAt + 1, Keywords: o2a.Message + " " + o1.Message}
	if r1 := <-store.Compliance().ComplianceExport(cr5); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		cposts := r1.Data.([]*model.CompliancePost)

		if len(cposts) != 2 {
			t.Fatal("return wrong results length")
		}

		if cposts[0].PostId != o1.Id {
			t.Fatal("Wrong sort")
		}
	}

	cr6 := &model.Compliance{Desc: "test" + model.NewId(), StartAt: o1.CreateAt - 1, EndAt: o2a.CreateAt + 1, Emails: u2.Email + ", " + u1.Email, Keywords: o2a.Message + " " + o1.Message}
	if r1 := <-store.Compliance().ComplianceExport(cr6); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		cposts := r1.Data.([]*model.CompliancePost)

		if len(cposts) != 2 {
			t.Fatal("return wrong results length")
		}

		if cposts[0].PostId != o1.Id {
			t.Fatal("Wrong sort")
		}

		if cposts[1].PostId != o2a.Id {
			t.Fatal("Wrong sort")
		}
	}
}

func TestComplianceExportDirectMessages(t *testing.T) {
	Setup()

	time.Sleep(100 * time.Millisecond)

	t1 := &model.Team{}
	t1.DisplayName = "DisplayName"
	t1.Name = "zz" + model.NewId() + "b"
	t1.Email = model.NewId() + "@nowhere.com"
	t1.Type = model.TEAM_OPEN
	t1 = Must(store.Team().Save(t1)).(*model.Team)

	u1 := &model.User{}
	u1.Email = model.NewId()
	u1.Username = model.NewId()
	u1 = Must(store.User().Save(u1)).(*model.User)
	Must(store.Team().SaveMember(&model.TeamMember{TeamId: t1.Id, UserId: u1.Id}))

	u2 := &model.User{}
	u2.Email = model.NewId()
	u2.Username = model.NewId()
	u2 = Must(store.User().Save(u2)).(*model.User)
	Must(store.Team().SaveMember(&model.TeamMember{TeamId: t1.Id, UserId: u2.Id}))

	c1 := &model.Channel{}
	c1.TeamId = t1.Id
	c1.DisplayName = "Channel2"
	c1.Name = "zz" + model.NewId() + "b"
	c1.Type = model.CHANNEL_OPEN
	c1 = Must(store.Channel().Save(c1)).(*model.Channel)

	cDM := Must(store.Channel().CreateDirectChannel(u1.Id, u2.Id)).(*model.Channel)

	o1 := &model.Post{}
	o1.ChannelId = c1.Id
	o1.UserId = u1.Id
	o1.CreateAt = model.GetMillis()
	o1.Message = "zz" + model.NewId() + "b"
	o1 = Must(store.Post().Save(o1)).(*model.Post)

	o1a := &model.Post{}
	o1a.ChannelId = c1.Id
	o1a.UserId = u1.Id
	o1a.CreateAt = o1.CreateAt + 10
	o1a.Message = "zz" + model.NewId() + "b"
	o1a = Must(store.Post().Save(o1a)).(*model.Post)

	o2 := &model.Post{}
	o2.ChannelId = c1.Id
	o2.UserId = u1.Id
	o2.CreateAt = o1.CreateAt + 20
	o2.Message = "zz" + model.NewId() + "b"
	o2 = Must(store.Post().Save(o2)).(*model.Post)

	o2a := &model.Post{}
	o2a.ChannelId = c1.Id
	o2a.UserId = u2.Id
	o2a.CreateAt = o1.CreateAt + 30
	o2a.Message = "zz" + model.NewId() + "b"
	o2a = Must(store.Post().Save(o2a)).(*model.Post)

	o3 := &model.Post{}
	o3.ChannelId = cDM.Id
	o3.UserId = u1.Id
	o3.CreateAt = o1.CreateAt + 40
	o3.Message = "zz" + model.NewId() + "b"
	o3 = Must(store.Post().Save(o3)).(*model.Post)

	time.Sleep(100 * time.Millisecond)

	cr1 := &model.Compliance{Desc: "test" + model.NewId(), StartAt: o1.CreateAt - 1, EndAt: o3.CreateAt + 1, Emails: u1.Email}
	if r1 := <-store.Compliance().ComplianceExport(cr1); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		cposts := r1.Data.([]*model.CompliancePost)

		if len(cposts) != 4 {
			t.Fatal("return wrong results length")
		}

		if cposts[0].PostId != o1.Id {
			t.Fatal("Wrong sort")
		}

		if cposts[len(cposts)-1].PostId != o3.Id {
			t.Fatal("Wrong sort")
		}
	}
}

func TestMessageExport(t *testing.T) {
	Setup()

	// get the starting number of message export entries
	startTime := model.GetMillis()
	var numMessageExports = 0
	if r1 := <-store.Compliance().MessageExport(startTime - 10, 10); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		messages := r1.Data.([]*model.MessageExport)
		numMessageExports = len(messages)
	}

	// need a team
	team := &model.Team{
		DisplayName: "DisplayName",
		Name: "zz" + model.NewId() + "b",
		Email: model.NewId() + "@nowhere.com",
		Type: model.TEAM_OPEN,
	}
	team = Must(store.Team().Save(team)).(*model.Team)

	// and two users that are a part of that team
	user1 := &model.User{
		Email:    model.NewId(),
		Username: model.NewId(),
	}
	user1 = Must(store.User().Save(user1)).(*model.User)
	Must(store.Team().SaveMember(&model.TeamMember{
		TeamId: team.Id,
		UserId: user1.Id,
	}))

	user2 := &model.User{
		Email: model.NewId(),
		Username: model.NewId(),
	}
	user2 = Must(store.User().Save(user2)).(*model.User)
	Must(store.Team().SaveMember(&model.TeamMember{
		TeamId: team.Id,
		UserId: user2.Id,
	}))

	// need a public channel as well as a DM channel between the two users
	channel := &model.Channel{
		TeamId: team.Id,
		DisplayName: "Channel2",
		Name: "zz" + model.NewId() + "b",
		Type: model.CHANNEL_OPEN,
	}
	channel = Must(store.Channel().Save(channel)).(*model.Channel)
	directMessageChannel := Must(store.Channel().CreateDirectChannel(user1.Id, user2.Id)).(*model.Channel)

	// user1 posts twice in the public channel
	post1 := &model.Post{
		ChannelId: channel.Id,
		UserId: user1.Id,
		CreateAt: startTime,
		Message: "zz" + model.NewId() + "a",
	}
	post1 = Must(store.Post().Save(post1)).(*model.Post)

	post2 := &model.Post{
		ChannelId: channel.Id,
		UserId: user1.Id,
		CreateAt: startTime + 10,
		Message: "zz" + model.NewId() + "b",
	}
	post2 = Must(store.Post().Save(post2)).(*model.Post)

	// they also send a DM to user2
	post3 := &model.Post{
		ChannelId: directMessageChannel.Id,
		UserId: user1.Id,
		CreateAt: startTime + 20,
		Message: "zz" + model.NewId() + "c",
	}
	post3 = Must(store.Post().Save(post3)).(*model.Post)

	// user2 has seen all messages in the public channel
	channelMember1 := &model.ChannelMember{
		ChannelId: channel.Id,
		UserId: user2.Id,
		LastViewedAt: startTime + 30,
		LastUpdateAt: startTime + 30,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	channelMember1 = Must(store.Channel().SaveMember(channelMember1)).(*model.ChannelMember)

	// a ChannelMember record is implicitly created for all users in a DM, so we need to update the existing record for
	// user2 to make it look like they've read the message that user1 sent
	channelMember2 := &model.ChannelMember{
		ChannelId: directMessageChannel.Id,
		UserId: user2.Id,
		LastViewedAt: startTime + 30,
		LastUpdateAt: startTime + 30,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	}
	channelMember2 = Must(store.Channel().UpdateMember(channelMember2)).(*model.ChannelMember)

	// fetch the message exports for all three posts that user1 sent
	messageExportMap := map[string]model.MessageExport{}
	if r1 := <-store.Compliance().MessageExport(startTime - 10, 10); r1.Err != nil {
		t.Fatal(r1.Err)
	} else {
		messages := r1.Data.([]*model.MessageExport)
		assert.Equal(t, numMessageExports + 3, len(messages))

		for _, v := range messages {
			messageExportMap[*v.PostId] = *v
		}
	}

	// post1 was made by user1 in channel1 and team1, but has no channel member because user1 hasn't viewed the channel
	assert.Equal(t, post1.Message, *messageExportMap[post1.Id].PostMessage)
	assert.Equal(t, channel.Id, *messageExportMap[post1.Id].ChannelId)
	assert.Equal(t, user1.Id, *messageExportMap[post1.Id].UserId)
	assert.Equal(t, team.Id, *messageExportMap[post1.Id].TeamId)
	assert.Empty(t, messageExportMap[post1.Id].ChannelMemberLastViewedAt)

	// post2 was made by user1 in channel1 and team1, but has no channel member because user1 hasn't viewed the channel
	assert.Equal(t, post2.Message, *messageExportMap[post2.Id].PostMessage)
	assert.Equal(t, channel.Id, *messageExportMap[post2.Id].ChannelId)
	assert.Equal(t, user1.Id, *messageExportMap[post2.Id].UserId)
	assert.Equal(t, team.Id, *messageExportMap[post2.Id].TeamId)
	assert.Empty(t, messageExportMap[post2.Id].ChannelMemberLastViewedAt)

	// post3 is a DM, so it has no team info, and channel member records were implicitly created for both users
	assert.Equal(t, post3.Message, *messageExportMap[post3.Id].PostMessage)
	assert.Equal(t, directMessageChannel.Id, *messageExportMap[post3.Id].ChannelId)
	assert.Equal(t, user1.Id, *messageExportMap[post3.Id].UserId)
	assert.Empty(t, messageExportMap[post3.Id].TeamId)
	assert.Equal(t, int64(0), *messageExportMap[post3.Id].ChannelMemberLastViewedAt)
}
