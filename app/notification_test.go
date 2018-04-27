// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/utils"
)

func TestSendNotifications(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()

	th.App.AddUserToChannel(th.BasicUser2, th.BasicChannel)

	post1, err := th.App.CreatePostMissingChannel(&model.Post{
		UserId:    th.BasicUser.Id,
		ChannelId: th.BasicChannel.Id,
		Message:   "@" + th.BasicUser2.Username,
		Type:      model.POST_ADD_TO_CHANNEL,
		Props:     map[string]interface{}{model.POST_PROPS_ADDED_USER_ID: "junk"},
	}, true)

	if err != nil {
		t.Fatal(err)
	}

	mentions, err := th.App.SendNotifications(post1, th.BasicTeam, th.BasicChannel, th.BasicUser, nil)
	if err != nil {
		t.Fatal(err)
	} else if mentions == nil {
		t.Log(mentions)
		t.Fatal("user should have been mentioned")
	} else if mentions[0] != th.BasicUser2.Id {
		t.Log(mentions)
		t.Fatal("user should have been mentioned")
	}

	dm, err := th.App.CreateDirectChannel(th.BasicUser.Id, th.BasicUser2.Id)
	if err != nil {
		t.Fatal(err)
	}

	post2, err := th.App.CreatePostMissingChannel(&model.Post{
		UserId:    th.BasicUser.Id,
		ChannelId: dm.Id,
		Message:   "dm message",
	}, true)

	if err != nil {
		t.Fatal(err)
	}

	_, err = th.App.SendNotifications(post2, th.BasicTeam, dm, th.BasicUser, nil)
	if err != nil {
		t.Fatal(err)
	}

	th.App.UpdateActive(th.BasicUser2, false)
	th.App.InvalidateAllCaches()

	post3, err := th.App.CreatePostMissingChannel(&model.Post{
		UserId:    th.BasicUser.Id,
		ChannelId: dm.Id,
		Message:   "dm message",
	}, true)

	if err != nil {
		t.Fatal(err)
	}

	_, err = th.App.SendNotifications(post3, th.BasicTeam, dm, th.BasicUser, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetExplicitMentions(t *testing.T) {
	id1 := model.NewId()
	id2 := model.NewId()
	id3 := model.NewId()

	for name, tc := range map[string]struct {
		Message  string
		Keywords map[string][]string
		Expected *ExplicitMentions
	}{
		"Nobody": {
			Message:  "this is a message",
			Keywords: map[string][]string{},
			Expected: &ExplicitMentions{},
		},
		"NonexistentUser": {
			Message: "this is a message for @user",
			Expected: &ExplicitMentions{
				OtherPotentialMentions: []string{"user"},
			},
		},
		"OnePerson": {
			Message:  "this is a message for @user",
			Keywords: map[string][]string{"@user": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"OnePersonWithPeriodAtEndOfUsername": {
			Message:  "this is a message for @user.name.",
			Keywords: map[string][]string{"@user.name.": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"OnePersonWithPeriodAtEndOfUsernameButNotSimilarName": {
			Message:  "this is a message for @user.name.",
			Keywords: map[string][]string{"@user.name.": {id1}, "@user.name": {id2}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"OnePersonAtEndOfSentence": {
			Message:  "this is a message for @user.",
			Keywords: map[string][]string{"@user": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"OnePersonWithoutAtMention": {
			Message:  "this is a message for @user",
			Keywords: map[string][]string{"this": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
				OtherPotentialMentions: []string{"user"},
			},
		},
		"MultiplePeopleWithOneWord": {
			Message:  "this is a message for @user",
			Keywords: map[string][]string{"@user": {id1, id2}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
					id2: true,
				},
			},
		},
		"OneOfMultiplePeople": {
			Message:  "this is a message for @user",
			Keywords: map[string][]string{"@user": {id1}, "@mention": {id2}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"MultiplePeopleWithMultipleWords": {
			Message:  "this is an @mention for @user",
			Keywords: map[string][]string{"@user": {id1}, "@mention": {id2}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
					id2: true,
				},
			},
		},
		"Channel": {
			Message:  "this is an message for @channel",
			Keywords: map[string][]string{"@channel": {id1, id2}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
					id2: true,
				},
				ChannelMentioned: true,
			},
		},
		"CapitalizedChannel": {
			Message:  "this is an message for @cHaNNeL",
			Keywords: map[string][]string{"@channel": {id1, id2}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
					id2: true,
				},
				ChannelMentioned: true,
			},
		},
		"All": {
			Message:  "this is an message for @all",
			Keywords: map[string][]string{"@all": {id1, id2}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
					id2: true,
				},
				AllMentioned: true,
			},
		},
		"CapitalizedAll": {
			Message:  "this is an message for @ALL",
			Keywords: map[string][]string{"@all": {id1, id2}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
					id2: true,
				},
				AllMentioned: true,
			},
		},
		"UserWithPeriod": {
			Message:  "user.period doesn't complicate things at all by including periods in their username",
			Keywords: map[string][]string{"user.period": {id1}, "user": {id2}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"AtUserWithPeriodAtEndOfSentence": {
			Message:  "this is a message for @user.period.",
			Keywords: map[string][]string{"@user.period": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"UserWithPeriodAtEndOfSentence": {
			Message:  "this is a message for user.period.",
			Keywords: map[string][]string{"user.period": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"PotentialOutOfChannelUser": {
			Message:  "this is an message for @potential and @user",
			Keywords: map[string][]string{"@user": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
				OtherPotentialMentions: []string{"potential"},
			},
		},
		"PotentialOutOfChannelUserWithPeriod": {
			Message: "this is an message for @potential.user",
			Expected: &ExplicitMentions{
				OtherPotentialMentions: []string{"potential.user"},
			},
		},
		"InlineCode": {
			Message:  "`this shouldn't mention @channel at all`",
			Keywords: map[string][]string{},
			Expected: &ExplicitMentions{},
		},
		"FencedCodeBlock": {
			Message:  "```\nthis shouldn't mention @channel at all\n```",
			Keywords: map[string][]string{},
			Expected: &ExplicitMentions{},
		},
		"Emphasis": {
			Message:  "*@aaa @bbb @ccc*",
			Keywords: map[string][]string{"@aaa": {id1}, "@bbb": {id2}, "@ccc": {id3}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
					id2: true,
					id3: true,
				},
			},
		},
		"StrongEmphasis": {
			Message:  "**@aaa @bbb @ccc**",
			Keywords: map[string][]string{"@aaa": {id1}, "@bbb": {id2}, "@ccc": {id3}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
					id2: true,
					id3: true,
				},
			},
		},
		"Strikethrough": {
			Message:  "~~@aaa @bbb @ccc~~",
			Keywords: map[string][]string{"@aaa": {id1}, "@bbb": {id2}, "@ccc": {id3}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
					id2: true,
					id3: true,
				},
			},
		},
		"Heading": {
			Message:  "### @aaa",
			Keywords: map[string][]string{"@aaa": {id1}, "@bbb": {id2}, "@ccc": {id3}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"BlockQuote": {
			Message:  "> @aaa",
			Keywords: map[string][]string{"@aaa": {id1}, "@bbb": {id2}, "@ccc": {id3}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"Emoji": {
			Message:  ":smile:",
			Keywords: map[string][]string{"smile": {id1}, "smiley": {id2}, "smiley_cat": {id3}},
			Expected: &ExplicitMentions{},
		},
		"NotEmoji": {
			Message:  "smile",
			Keywords: map[string][]string{"smile": {id1}, "smiley": {id2}, "smiley_cat": {id3}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"UnclosedEmoji": {
			Message:  ":smile",
			Keywords: map[string][]string{"smile": {id1}, "smiley": {id2}, "smiley_cat": {id3}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"UnopenedEmoji": {
			Message:  "smile:",
			Keywords: map[string][]string{"smile": {id1}, "smiley": {id2}, "smiley_cat": {id3}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"IndentedCodeBlock": {
			Message:  "    this shouldn't mention @channel at all",
			Keywords: map[string][]string{},
			Expected: &ExplicitMentions{},
		},
		"LinkTitle": {
			Message:  `[foo](this "shouldn't mention @channel at all")`,
			Keywords: map[string][]string{},
			Expected: &ExplicitMentions{},
		},
		"MalformedInlineCode": {
			Message:  "`this should mention @channel``",
			Keywords: map[string][]string{},
			Expected: &ExplicitMentions{
				ChannelMentioned: true,
			},
		},

		// The following tests cover cases where the message mentions @user.name, so we shouldn't assume that
		// the user might be intending to mention some @user that isn't in the channel.
		"Don't include potential mention that's part of an actual mention (without trailing period)": {
			Message:  "this is an message for @user.name",
			Keywords: map[string][]string{"@user.name": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"Don't include potential mention that's part of an actual mention (with trailing period)": {
			Message:  "this is an message for @user.name.",
			Keywords: map[string][]string{"@user.name": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"Don't include potential mention that's part of an actual mention (with multiple trailing periods)": {
			Message:  "this is an message for @user.name...",
			Keywords: map[string][]string{"@user.name": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
		"Don't include potential mention that's part of an actual mention (containing and followed by multiple periods)": {
			Message:  "this is an message for @user...name...",
			Keywords: map[string][]string{"@user...name": {id1}},
			Expected: &ExplicitMentions{
				MentionedUserIds: map[string]bool{
					id1: true,
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			m := GetExplicitMentions(tc.Message, tc.Keywords)
			if tc.Expected.MentionedUserIds == nil {
				tc.Expected.MentionedUserIds = make(map[string]bool)
			}
			assert.EqualValues(t, tc.Expected, m)
		})
	}
}

func TestGetExplicitMentionsAtHere(t *testing.T) {
	// test all the boundary cases that we know can break up terms (and those that we know won't)
	cases := map[string]bool{
		"":          false,
		"here":      false,
		"@here":     true,
		" @here ":   true,
		"\n@here\n": true,
		"!@here!":   true,
		"#@here#":   true,
		"$@here$":   true,
		"%@here%":   true,
		"^@here^":   true,
		"&@here&":   true,
		"*@here*":   true,
		"(@here(":   true,
		")@here)":   true,
		"-@here-":   true,
		"_@here_":   false, // This case shouldn't mention since it would be mentioning "@here_"
		"=@here=":   true,
		"+@here+":   true,
		"[@here[":   true,
		"{@here{":   true,
		"]@here]":   true,
		"}@here}":   true,
		"\\@here\\": true,
		"|@here|":   true,
		";@here;":   true,
		":@here:":   false, // This case shouldn't trigger a mention since it follows the format of reactions e.g. :word:
		"'@here'":   true,
		"\"@here\"": true,
		",@here,":   true,
		"<@here<":   true,
		".@here.":   true,
		">@here>":   true,
		"/@here/":   true,
		"?@here?":   true,
		"`@here`":   false, // This case shouldn't mention since it's a code block
		"~@here~":   true,
		"@HERE":     true,
		"@hERe":     true,
	}

	for message, shouldMention := range cases {
		if m := GetExplicitMentions(message, nil); m.HereMentioned && !shouldMention {
			t.Fatalf("shouldn't have mentioned @here with \"%v\"", message)
		} else if !m.HereMentioned && shouldMention {
			t.Fatalf("should've mentioned @here with \"%v\"", message)
		}
	}

	// mentioning @here and someone
	id := model.NewId()
	if m := GetExplicitMentions("@here @user @potential", map[string][]string{"@user": {id}}); !m.HereMentioned {
		t.Fatal("should've mentioned @here with \"@here @user\"")
	} else if len(m.MentionedUserIds) != 1 || !m.MentionedUserIds[id] {
		t.Fatal("should've mentioned @user with \"@here @user\"")
	} else if len(m.OtherPotentialMentions) > 1 {
		t.Fatal("should've potential mentions for @potential")
	}
}

func TestGetMentionKeywords(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	// user with username or custom mentions enabled
	user1 := &model.User{
		Id:        model.NewId(),
		FirstName: "First",
		Username:  "User",
		NotifyProps: map[string]string{
			"mention_keys": "User,@User,MENTION",
		},
	}

	profiles := map[string]*model.User{user1.Id: user1}
	mentions := th.App.GetMentionKeywordsInChannel(profiles, true)
	if len(mentions) != 3 {
		t.Fatal("should've returned three mention keywords")
	} else if ids, ok := mentions["user"]; !ok || ids[0] != user1.Id {
		t.Fatal("should've returned mention key of user")
	} else if ids, ok := mentions["@user"]; !ok || ids[0] != user1.Id {
		t.Fatal("should've returned mention key of @user")
	} else if ids, ok := mentions["mention"]; !ok || ids[0] != user1.Id {
		t.Fatal("should've returned mention key of mention")
	}

	// user with first name mention enabled
	user2 := &model.User{
		Id:        model.NewId(),
		FirstName: "First",
		Username:  "User",
		NotifyProps: map[string]string{
			"first_name": "true",
		},
	}

	profiles = map[string]*model.User{user2.Id: user2}
	mentions = th.App.GetMentionKeywordsInChannel(profiles, true)
	if len(mentions) != 2 {
		t.Fatal("should've returned two mention keyword")
	} else if ids, ok := mentions["First"]; !ok || ids[0] != user2.Id {
		t.Fatal("should've returned mention key of First")
	}

	// user with @channel/@all mentions enabled
	user3 := &model.User{
		Id:        model.NewId(),
		FirstName: "First",
		Username:  "User",
		NotifyProps: map[string]string{
			"channel": "true",
		},
	}

	profiles = map[string]*model.User{user3.Id: user3}
	mentions = th.App.GetMentionKeywordsInChannel(profiles, true)
	if len(mentions) != 3 {
		t.Fatal("should've returned three mention keywords")
	} else if ids, ok := mentions["@channel"]; !ok || ids[0] != user3.Id {
		t.Fatal("should've returned mention key of @channel")
	} else if ids, ok := mentions["@all"]; !ok || ids[0] != user3.Id {
		t.Fatal("should've returned mention key of @all")
	}

	// user with all types of mentions enabled
	user4 := &model.User{
		Id:        model.NewId(),
		FirstName: "First",
		Username:  "User",
		NotifyProps: map[string]string{
			"mention_keys": "User,@User,MENTION",
			"first_name":   "true",
			"channel":      "true",
		},
	}

	profiles = map[string]*model.User{user4.Id: user4}
	mentions = th.App.GetMentionKeywordsInChannel(profiles, true)
	if len(mentions) != 6 {
		t.Fatal("should've returned six mention keywords")
	} else if ids, ok := mentions["user"]; !ok || ids[0] != user4.Id {
		t.Fatal("should've returned mention key of user")
	} else if ids, ok := mentions["@user"]; !ok || ids[0] != user4.Id {
		t.Fatal("should've returned mention key of @user")
	} else if ids, ok := mentions["mention"]; !ok || ids[0] != user4.Id {
		t.Fatal("should've returned mention key of mention")
	} else if ids, ok := mentions["First"]; !ok || ids[0] != user4.Id {
		t.Fatal("should've returned mention key of First")
	} else if ids, ok := mentions["@channel"]; !ok || ids[0] != user4.Id {
		t.Fatal("should've returned mention key of @channel")
	} else if ids, ok := mentions["@all"]; !ok || ids[0] != user4.Id {
		t.Fatal("should've returned mention key of @all")
	}

	dup_count := func(list []string) map[string]int {

		duplicate_frequency := make(map[string]int)

		for _, item := range list {
			// check if the item/element exist in the duplicate_frequency map

			_, exist := duplicate_frequency[item]

			if exist {
				duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
			} else {
				duplicate_frequency[item] = 1 // else start counting from 1
			}
		}
		return duplicate_frequency
	}

	// multiple users but no more than MaxNotificationsPerChannel
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.TeamSettings.MaxNotificationsPerChannel = 4 })
	profiles = map[string]*model.User{
		user1.Id: user1,
		user2.Id: user2,
		user3.Id: user3,
		user4.Id: user4,
	}
	mentions = th.App.GetMentionKeywordsInChannel(profiles, true)
	if len(mentions) != 6 {
		t.Fatal("should've returned six mention keywords")
	} else if ids, ok := mentions["user"]; !ok || len(ids) != 2 || (ids[0] != user1.Id && ids[1] != user1.Id) || (ids[0] != user4.Id && ids[1] != user4.Id) {
		t.Fatal("should've mentioned user1 and user4 with user")
	} else if ids := dup_count(mentions["@user"]); len(ids) != 4 || (ids[user1.Id] != 2) || (ids[user4.Id] != 2) {
		t.Fatal("should've mentioned user1 and user4 with @user")
	} else if ids, ok := mentions["mention"]; !ok || len(ids) != 2 || (ids[0] != user1.Id && ids[1] != user1.Id) || (ids[0] != user4.Id && ids[1] != user4.Id) {
		t.Fatal("should've mentioned user1 and user4 with mention")
	} else if ids, ok := mentions["First"]; !ok || len(ids) != 2 || (ids[0] != user2.Id && ids[1] != user2.Id) || (ids[0] != user4.Id && ids[1] != user4.Id) {
		t.Fatal("should've mentioned user2 and user4 with First")
	} else if ids, ok := mentions["@channel"]; !ok || len(ids) != 2 || (ids[0] != user3.Id && ids[1] != user3.Id) || (ids[0] != user4.Id && ids[1] != user4.Id) {
		t.Fatal("should've mentioned user3 and user4 with @channel")
	} else if ids, ok := mentions["@all"]; !ok || len(ids) != 2 || (ids[0] != user3.Id && ids[1] != user3.Id) || (ids[0] != user4.Id && ids[1] != user4.Id) {
		t.Fatal("should've mentioned user3 and user4 with @all")
	}

	// multiple users and more than MaxNotificationsPerChannel
	th.App.UpdateConfig(func(cfg *model.Config) { *cfg.TeamSettings.MaxNotificationsPerChannel = 3 })
	mentions = th.App.GetMentionKeywordsInChannel(profiles, true)
	if len(mentions) != 4 {
		t.Fatal("should've returned four mention keywords")
	} else if _, ok := mentions["@channel"]; ok {
		t.Fatal("should not have mentioned any user with @channel")
	} else if _, ok := mentions["@all"]; ok {
		t.Fatal("should not have mentioned any user with @all")
	} else if _, ok := mentions["@here"]; ok {
		t.Fatal("should not have mentioned any user with @here")
	}

	// no special mentions
	profiles = map[string]*model.User{
		user1.Id: user1,
	}
	mentions = th.App.GetMentionKeywordsInChannel(profiles, false)
	if len(mentions) != 3 {
		t.Fatal("should've returned three mention keywords")
	} else if ids, ok := mentions["user"]; !ok || len(ids) != 1 || ids[0] != user1.Id {
		t.Fatal("should've mentioned user1 with user")
	} else if ids, ok := mentions["@user"]; !ok || len(ids) != 2 || ids[0] != user1.Id || ids[1] != user1.Id {
		t.Fatal("should've mentioned user1 twice with @user")
	} else if ids, ok := mentions["mention"]; !ok || len(ids) != 1 || ids[0] != user1.Id {
		t.Fatal("should've mentioned user1 with mention")
	} else if _, ok := mentions["First"]; ok {
		t.Fatal("should not have mentioned user1 with First")
	} else if _, ok := mentions["@channel"]; ok {
		t.Fatal("should not have mentioned any user with @channel")
	} else if _, ok := mentions["@all"]; ok {
		t.Fatal("should not have mentioned any user with @all")
	} else if _, ok := mentions["@here"]; ok {
		t.Fatal("should not have mentioned any user with @here")
	}
}

func TestGetNotificationEmailSubject(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	expectedPrefix := "[http://localhost:8065] Notification in team on"
	post := &model.Post{
		CreateAt: 1501804801000,
	}
	translateFunc := utils.GetUserTranslations("en")
	subject := getNotificationEmailSubject(post, translateFunc, "http://localhost:8065", "team")
	if !strings.HasPrefix(subject, expectedPrefix) {
		t.Fatal("Expected subject line prefix '" + expectedPrefix + "', got " + subject)
	}
}

func TestGetNotificationEmailBodyFullNotificationPublicChannel(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	recipient := &model.User{}
	post := &model.Post{
		Message: "This is the message",
	}
	channel := &model.Channel{
		DisplayName: "ChannelName",
		Type:        model.CHANNEL_OPEN,
	}
	senderName := "sender"
	teamName := "team"
	teamURL := "http://localhost:8065/" + teamName
	emailNotificationContentsType := model.EMAIL_NOTIFICATION_CONTENTS_FULL
	translateFunc := utils.GetUserTranslations("en")

	body := th.App.getNotificationEmailBody(recipient, post, channel, senderName, teamName, teamURL, emailNotificationContentsType, translateFunc)
	if !strings.Contains(body, "You have a new notification.") {
		t.Fatal("Expected email text 'You have a new notification. Got " + body)
	}
	if !strings.Contains(body, "CHANNEL: "+channel.DisplayName) {
		t.Fatal("Expected email text 'CHANNEL: " + channel.DisplayName + "'. Got " + body)
	}
	if !strings.Contains(body, senderName+" - ") {
		t.Fatal("Expected email text '" + senderName + " - '. Got " + body)
	}
	if !strings.Contains(body, post.Message) {
		t.Fatal("Expected email text '" + post.Message + "'. Got " + body)
	}
	if !strings.Contains(body, teamURL) {
		t.Fatal("Expected email text '" + teamURL + "'. Got " + body)
	}
}

func TestGetNotificationEmailBodyFullNotificationGroupChannel(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	recipient := &model.User{}
	post := &model.Post{
		Message: "This is the message",
	}
	channel := &model.Channel{
		DisplayName: "ChannelName",
		Type:        model.CHANNEL_GROUP,
	}
	senderName := "sender"
	teamName := "team"
	teamURL := "http://localhost:8065/" + teamName
	emailNotificationContentsType := model.EMAIL_NOTIFICATION_CONTENTS_FULL
	translateFunc := utils.GetUserTranslations("en")

	body := th.App.getNotificationEmailBody(recipient, post, channel, senderName, teamName, teamURL, emailNotificationContentsType, translateFunc)
	if !strings.Contains(body, "You have a new notification.") {
		t.Fatal("Expected email text 'You have a new notification. Got " + body)
	}
	if !strings.Contains(body, "CHANNEL: Group Message") {
		t.Fatal("Expected email text 'CHANNEL: Group Message'. Got " + body)
	}
	if !strings.Contains(body, senderName+" - ") {
		t.Fatal("Expected email text '" + senderName + " - '. Got " + body)
	}
	if !strings.Contains(body, post.Message) {
		t.Fatal("Expected email text '" + post.Message + "'. Got " + body)
	}
	if !strings.Contains(body, teamURL) {
		t.Fatal("Expected email text '" + teamURL + "'. Got " + body)
	}
}

func TestGetNotificationEmailBodyFullNotificationPrivateChannel(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	recipient := &model.User{}
	post := &model.Post{
		Message: "This is the message",
	}
	channel := &model.Channel{
		DisplayName: "ChannelName",
		Type:        model.CHANNEL_PRIVATE,
	}
	senderName := "sender"
	teamName := "team"
	teamURL := "http://localhost:8065/" + teamName
	emailNotificationContentsType := model.EMAIL_NOTIFICATION_CONTENTS_FULL
	translateFunc := utils.GetUserTranslations("en")

	body := th.App.getNotificationEmailBody(recipient, post, channel, senderName, teamName, teamURL, emailNotificationContentsType, translateFunc)
	if !strings.Contains(body, "You have a new notification.") {
		t.Fatal("Expected email text 'You have a new notification. Got " + body)
	}
	if !strings.Contains(body, "CHANNEL: "+channel.DisplayName) {
		t.Fatal("Expected email text 'CHANNEL: " + channel.DisplayName + "'. Got " + body)
	}
	if !strings.Contains(body, senderName+" - ") {
		t.Fatal("Expected email text '" + senderName + " - '. Got " + body)
	}
	if !strings.Contains(body, post.Message) {
		t.Fatal("Expected email text '" + post.Message + "'. Got " + body)
	}
	if !strings.Contains(body, teamURL) {
		t.Fatal("Expected email text '" + teamURL + "'. Got " + body)
	}
}

func TestGetNotificationEmailBodyFullNotificationDirectChannel(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	recipient := &model.User{}
	post := &model.Post{
		Message: "This is the message",
	}
	channel := &model.Channel{
		DisplayName: "ChannelName",
		Type:        model.CHANNEL_DIRECT,
	}
	senderName := "sender"
	teamName := "team"
	teamURL := "http://localhost:8065/" + teamName
	emailNotificationContentsType := model.EMAIL_NOTIFICATION_CONTENTS_FULL
	translateFunc := utils.GetUserTranslations("en")

	body := th.App.getNotificationEmailBody(recipient, post, channel, senderName, teamName, teamURL, emailNotificationContentsType, translateFunc)
	if !strings.Contains(body, "You have a new direct message.") {
		t.Fatal("Expected email text 'You have a new direct message. Got " + body)
	}
	if !strings.Contains(body, senderName+" - ") {
		t.Fatal("Expected email text '" + senderName + " - '. Got " + body)
	}
	if !strings.Contains(body, post.Message) {
		t.Fatal("Expected email text '" + post.Message + "'. Got " + body)
	}
	if !strings.Contains(body, teamURL) {
		t.Fatal("Expected email text '" + teamURL + "'. Got " + body)
	}
}

// from here
func TestGetNotificationEmailBodyGenericNotificationPublicChannel(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	recipient := &model.User{}
	post := &model.Post{
		Message: "This is the message",
	}
	channel := &model.Channel{
		DisplayName: "ChannelName",
		Type:        model.CHANNEL_OPEN,
	}
	senderName := "sender"
	teamName := "team"
	teamURL := "http://localhost:8065/" + teamName
	emailNotificationContentsType := model.EMAIL_NOTIFICATION_CONTENTS_GENERIC
	translateFunc := utils.GetUserTranslations("en")

	body := th.App.getNotificationEmailBody(recipient, post, channel, senderName, teamName, teamURL, emailNotificationContentsType, translateFunc)
	if !strings.Contains(body, "You have a new notification from "+senderName) {
		t.Fatal("Expected email text 'You have a new notification from " + senderName + "'. Got " + body)
	}
	if strings.Contains(body, "CHANNEL: "+channel.DisplayName) {
		t.Fatal("Did not expect email text 'CHANNEL: " + channel.DisplayName + "'. Got " + body)
	}
	if strings.Contains(body, post.Message) {
		t.Fatal("Did not expect email text '" + post.Message + "'. Got " + body)
	}
	if !strings.Contains(body, teamURL) {
		t.Fatal("Expected email text '" + teamURL + "'. Got " + body)
	}
}

func TestGetNotificationEmailBodyGenericNotificationGroupChannel(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	recipient := &model.User{}
	post := &model.Post{
		Message: "This is the message",
	}
	channel := &model.Channel{
		DisplayName: "ChannelName",
		Type:        model.CHANNEL_GROUP,
	}
	senderName := "sender"
	teamName := "team"
	teamURL := "http://localhost:8065/" + teamName
	emailNotificationContentsType := model.EMAIL_NOTIFICATION_CONTENTS_GENERIC
	translateFunc := utils.GetUserTranslations("en")

	body := th.App.getNotificationEmailBody(recipient, post, channel, senderName, teamName, teamURL, emailNotificationContentsType, translateFunc)
	if !strings.Contains(body, "You have a new notification from "+senderName) {
		t.Fatal("Expected email text 'You have a new notification from " + senderName + "'. Got " + body)
	}
	if strings.Contains(body, "CHANNEL: "+channel.DisplayName) {
		t.Fatal("Did not expect email text 'CHANNEL: " + channel.DisplayName + "'. Got " + body)
	}
	if strings.Contains(body, post.Message) {
		t.Fatal("Did not expect email text '" + post.Message + "'. Got " + body)
	}
	if !strings.Contains(body, teamURL) {
		t.Fatal("Expected email text '" + teamURL + "'. Got " + body)
	}
}

func TestGetNotificationEmailBodyGenericNotificationPrivateChannel(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	recipient := &model.User{}
	post := &model.Post{
		Message: "This is the message",
	}
	channel := &model.Channel{
		DisplayName: "ChannelName",
		Type:        model.CHANNEL_PRIVATE,
	}
	senderName := "sender"
	teamName := "team"
	teamURL := "http://localhost:8065/" + teamName
	emailNotificationContentsType := model.EMAIL_NOTIFICATION_CONTENTS_GENERIC
	translateFunc := utils.GetUserTranslations("en")

	body := th.App.getNotificationEmailBody(recipient, post, channel, senderName, teamName, teamURL, emailNotificationContentsType, translateFunc)
	if !strings.Contains(body, "You have a new notification from "+senderName) {
		t.Fatal("Expected email text 'You have a new notification from " + senderName + "'. Got " + body)
	}
	if strings.Contains(body, "CHANNEL: "+channel.DisplayName) {
		t.Fatal("Did not expect email text 'CHANNEL: " + channel.DisplayName + "'. Got " + body)
	}
	if strings.Contains(body, post.Message) {
		t.Fatal("Did not expect email text '" + post.Message + "'. Got " + body)
	}
	if !strings.Contains(body, teamURL) {
		t.Fatal("Expected email text '" + teamURL + "'. Got " + body)
	}
}

func TestGetNotificationEmailBodyGenericNotificationDirectChannel(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	recipient := &model.User{}
	post := &model.Post{
		Message: "This is the message",
	}
	channel := &model.Channel{
		DisplayName: "ChannelName",
		Type:        model.CHANNEL_DIRECT,
	}
	senderName := "sender"
	teamName := "team"
	teamURL := "http://localhost:8065/" + teamName
	emailNotificationContentsType := model.EMAIL_NOTIFICATION_CONTENTS_GENERIC
	translateFunc := utils.GetUserTranslations("en")

	body := th.App.getNotificationEmailBody(recipient, post, channel, senderName, teamName, teamURL, emailNotificationContentsType, translateFunc)
	if !strings.Contains(body, "You have a new direct message from "+senderName) {
		t.Fatal("Expected email text 'You have a new direct message from " + senderName + "'. Got " + body)
	}
	if strings.Contains(body, "CHANNEL: "+channel.DisplayName) {
		t.Fatal("Did not expect email text 'CHANNEL: " + channel.DisplayName + "'. Got " + body)
	}
	if strings.Contains(body, post.Message) {
		t.Fatal("Did not expect email text '" + post.Message + "'. Got " + body)
	}
	if !strings.Contains(body, teamURL) {
		t.Fatal("Expected email text '" + teamURL + "'. Got " + body)
	}
}
