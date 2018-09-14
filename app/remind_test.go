// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/utils"

	"time"
	"github.com/mattermost/mattermost-server/mlog"
	"fmt"
)

const (
	UNABLE_TO_SCHEDULE_REMINDER = "unable to schedule reminder"
)

func TestInitReminders(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	th.App.InitReminders()
}

func TestStopReminders(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	th.App.StopReminders()
}

func TestListReminders(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	list := th.App.ListReminders("user_id")
	if list == "" {
		t.Fatal("list should not be empty")
	}
}

func TestDeleteReminders(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	th.App.DeleteReminders("user_id")
}

func TestScheduleReminders(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	user, uErr := th.App.GetUserByUsername("remindbot")
	if uErr != nil { t.Fatal("remindbot doesn't exist") }
	translateFunc := utils.GetUserTranslations(user.Locale)
	request := &model.ReminderRequest{}
	request.UserId = user.Id


	request.Payload = "me foo in 2 seconds"
	response, err := th.App.ScheduleReminder(request)
	if err != nil { t.Fatal(UNABLE_TO_SCHEDULE_REMINDER) }
	var responseParameters = map[string]interface{}{
		"Target":  "You",
		"UseTo":   "",
		"Message": "foo",
		"When":    "in 2 seconds",
	}
	expectedResponse := translateFunc("app.reminder.response", responseParameters)
	if response != expectedResponse {
		t.Fatal("\""+response+"\" doesn't match \""+ expectedResponse+"\"")
	}


	request.Payload = "@bob foo in 2 seconds"
	response, err = th.App.ScheduleReminder(request)
	if err != nil { t.Fatal(UNABLE_TO_SCHEDULE_REMINDER) }
	responseParameters = map[string]interface{}{
		"Target":  "@bob",
		"UseTo":   "",
		"Message": "foo",
		"When":    "in 2 seconds",
	}
	expectedResponse = translateFunc("app.reminder.response", responseParameters)
	if response != expectedResponse {
		t.Fatal("\""+response+"\" doesn't match \""+ expectedResponse+"\"")
	}


	request.Payload = "~off-topic foo in 2 seconds"
	response, err = th.App.ScheduleReminder(request)
	if err != nil { t.Fatal(UNABLE_TO_SCHEDULE_REMINDER) }
	responseParameters = map[string]interface{}{
		"Target":  "~off-topic",
		"UseTo":   "",
		"Message": "foo",
		"When":    "in 2 seconds",
	}
	expectedResponse = translateFunc("app.reminder.response", responseParameters)
	if response != expectedResponse {
		t.Fatal("\""+response+"\" doesn't match \""+ expectedResponse+"\"")
	}


	request.Payload = "me \"foo foo foo\" in 2 seconds"
	response, err = th.App.ScheduleReminder(request)
	if err != nil { t.Fatal(UNABLE_TO_SCHEDULE_REMINDER) }
	responseParameters = map[string]interface{}{
		"Target":  "You",
		"UseTo":   "",
		"Message": "foo foo foo",
		"When":    "in 2 seconds",
	}
	expectedResponse = translateFunc("app.reminder.response", responseParameters)
	if response != expectedResponse {
		t.Fatal("\""+response+"\" doesn't match \""+ expectedResponse+"\"")
	}


	request.Payload = "me foo in 5 seconds"
	response, err = th.App.ScheduleReminder(request)
	if err != nil { t.Fatal(UNABLE_TO_SCHEDULE_REMINDER) }
	responseParameters = map[string]interface{}{
		"Target":  "You",
		"UseTo":   "",
		"Message": "foo",
		"When":    "in 5 seconds",
	}
	expectedResponse = translateFunc("app.reminder.response", responseParameters)
	if response != expectedResponse {
		t.Fatal("\""+response+"\" doesn't match \""+ expectedResponse+"\"")
	}


	request.Payload = "me foo at 2:04 pm"
	response, err = th.App.ScheduleReminder(request)
	if err != nil { t.Fatal(UNABLE_TO_SCHEDULE_REMINDER) }
	responseParameters = map[string]interface{}{
		"Target":  "You",
		"UseTo":   "",
		"Message": "foo",
		"When":    "at 2:04 pm",
	}
	expectedResponse = translateFunc("app.reminder.response", responseParameters)
	if response != expectedResponse {
		t.Fatal("\""+response+"\" doesn't match \""+ expectedResponse+"\"")
	}


	//request.Payload = "me foo on monday at 12:30PM"
	//response, err = th.App.ScheduleReminder(request)
	//if err != nil { t.Fatal(UNABLE_TO_SCHEDULE_REMINDER) }
	//responseParameters = map[string]interface{}{
	//	"Target":  "You",
	//	"UseTo":   "",
	//	"Message": "foo",
	//	"When":    "on monday at 12:30PM",
	//}
	//expectedResponse = translateFunc("app.reminder.response", responseParameters)
	//if response != expectedResponse {
	//	t.Fatal("\""+response+"\" doesn't match \""+ expectedResponse+"\"")
	//}

	//TODO TEST Every

	//TODO TEST Outlier
}


func TestIn(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	th.App.InitReminders()
	user, uErr := th.App.GetUserByUsername("remindbot")
	if uErr != nil { t.Fatal("remindbot doesn't exist") }


	when := "in one second"
	times, iErr := th.App.in(when, user)
	if iErr != nil { t.Fatal("in one second doesn't parse")}
	var duration time.Duration
	duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	if duration != time.Second {
		t.Fatal("in one second isn't correct")
	}


	when = "in 712 minutes"
	times, iErr = th.App.in(when, user)
	if iErr != nil { t.Fatal("in 712 minutes doesn't parse")}
	duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	if duration != time.Minute * time.Duration(712) {
		t.Fatal("in 712 minutes isn't correct")
	}


	when = "in three hours"
	times, iErr = th.App.in(when, user)
	if iErr != nil { t.Fatal("in three hours doesn't parse")}
	duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	if duration != time.Hour * time.Duration(3) {
		t.Fatal("in three hours isn't correct")
	}


	when = "in 2 days"
	times, iErr = th.App.in(when, user)
	if iErr != nil { t.Fatal("in 2 days doesn't parse")}
	duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	if duration != time.Hour * time.Duration(24) * time.Duration(2) {
		t.Fatal("in 2 days isn't correct")
	}


	when = "in 90 weeks"
	times, iErr = th.App.in(when, user)
	if iErr != nil { t.Fatal("in 90 weeks doesn't parse")}
	duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	if duration != time.Hour * time.Duration(24) * time.Duration(7) * time.Duration(90) {
		t.Fatal("in 90 weeks isn't correct")
	}


	when = "in 4 months"
	times, iErr = th.App.in(when, user)
	if iErr != nil { t.Fatal("in 4 months doesn't parse")}
	duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	if duration != time.Hour * time.Duration(24) * time.Duration(30) * time.Duration(4) {
		t.Fatal("in 4 months isn't correct")
	}


	when = "in one year"
	times, iErr = th.App.in(when, user)
	if iErr != nil { t.Fatal("in one year doesn't parse")}
	duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	if duration != time.Hour * time.Duration(24) * time.Duration(365)  {
		t.Fatal("in one year isn't correct")
	}

}

func TestAt(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	th.App.InitReminders()
	user, uErr := th.App.GetUserByUsername("remindbot")
	if uErr != nil { t.Fatal("remindbot doesn't exist") }


	when := "at noon"
	times, iErr := th.App.at(when, user)
	if iErr != nil { t.Fatal("at noon doesn't parse")}
	if times[0].Hour() != 12 {
		t.Fatal("at noon isn't correct")
	}


	when = "at midnight"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at midnight doesn't parse")}
	if times[0].Hour() != 0 {
		t.Fatal("at midnight isn't correct")
	}


	when = "at two"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at two doesn't parse")}
	if times[0].Hour() != 2 && times[0].Hour() != 14 {
		t.Fatal("at two isn't correct")
	}


	when = "at 7"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 7 doesn't parse")}
	if times[0].Hour() != 7 && times[0].Hour() != 19 {
		t.Fatal("at 7 isn't correct")
	}


	when = "at 12:30pm"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 12:30pm doesn't parse")}
	if times[0].Hour() != 12 && times[0].Minute() != 30 {
		t.Fatal("at 12:30pm isn't correct")
	}


	when = "at 7:12 pm"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 7:12 pm doesn't parse")}
	if times[0].Hour() != 19 && times[0].Minute() != 12 {
		t.Fatal("at 7:12 pm isn't correct")
	}


	when = "at 8:05 PM"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 8:05 PM doesn't parse")}
	if times[0].Hour() != 10 && times[0].Minute() != 5 {
		t.Fatal("at 8:05 PM isn't correct")
	}


	when = "at 9:52 am"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 9:52 am doesn't parse")}
	if times[0].Hour() != 9 && times[0].Minute() != 52 {
		t.Fatal("at 9:52 am isn't correct")
	}


	when = "at 9:12"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 9:12 doesn't parse")}
	if times[0].Hour() != 9 && times[0].Hour() != 21 && times[0].Minute() != 12 {
		t.Fatal("at 9:12 isn't correct")
	}


	when = "at 17:15"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 17:15 doesn't parse") }
	if times[0].Hour() != 17 && times[0].Minute() != 15 {
		t.Fatal("at 17:15 isn't correct")
	}


	when = "at 930am"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 930am doesn't parse") }
	if times[0].Hour() != 9 && times[0].Minute() != 30 {
		t.Fatal("at 930am isn't correct")
	}


	when = "at 1230 am"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 1230 am doesn't parse") }
	if times[0].Hour() != 0 && times[0].Minute() != 30 {
		t.Fatal("at 1230 am isn't correct")
	}


	when = "at 5PM"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 5PM doesn't parse") }
	if times[0].Hour() != 17 && times[0].Minute() != 0 {
		t.Fatal("at 5PM isn't correct")
	}


	when = "at 4 am"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 4 am doesn't parse") }
	if times[0].Hour() != 4 && times[0].Minute() != 0 {
		t.Fatal("at 4 am isn't correct")
	}


	when = "at 1400"
	times, iErr = th.App.at(when, user)
	if iErr != nil { t.Fatal("at 1400 doesn't parse") }
	if times[0].Hour() != 14 && times[0].Minute() != 0 {
		t.Fatal("at 1400 isn't correct")
	}

	//TODO
	/*
	when = "at 11:00 every Thursday";

	when = "at 3pm every day";
	 */
}

func TestOn(t *testing.T) {
	th := Setup()
	defer th.TearDown()

	th.App.InitReminders()
	user, uErr := th.App.GetUserByUsername("remindbot")
	if uErr != nil { t.Fatal("remindbot doesn't exist") }


	when := "on Monday"
	times, iErr := th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on Monday doesn't parse")
	}
	if times[0].Weekday().String() != "Monday" {
		t.Fatal("on Monday isn't correct")
	}


	when = "on Tuesday"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on Tuesday doesn't parse")
	}
	if times[0].Weekday().String() != "Tuesday" {
		t.Fatal("on Tuesday isn't correct")
	}


	when = "on Wednesday"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on Wednesday doesn't parse")
	}
	if times[0].Weekday().String() != "Wednesday" {
		t.Fatal("on Wednesday isn't correct")
	}


	when = "on Thursday"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on Thursday doesn't parse")
	}
	if times[0].Weekday().String() != "Thursday" {
		t.Fatal("on Thursday isn't correct")
	}


	when = "on Friday"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on Friday doesn't parse")
	}
	if times[0].Weekday().String() != "Friday" {
		t.Fatal("on Friday isn't correct")
	}


	/*
        when = "on Mondays";

        when = "on Tuesdays at 11:15";

        when = "on Wednesdays";

        when = "on Thursdays";

        when = "on Fridays";

	*/


	when = "on mon"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on mon doesn't parse")
	}
	if times[0].Weekday().String() != "Monday" {
		t.Fatal("on mon isn't correct")
	}


	when = "on wED"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on wED doesn't parse")
	}
	if times[0].Weekday().String() != "Wednesday" {
		t.Fatal("on wED isn't correct")
	}


	when = "on tuesday at noon"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on tuesday at noon doesn't parse")
	}
	if times[0].Weekday().String() != "Tuesday" && times[0].Hour() != 12 {
		t.Fatal("on tuesday at noon isn't correct")
	}


	when = "on sunday at 3:42am"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on sunday at 3:42am doesn't parse")
	}
	if times[0].Weekday().String() != "Sunday" && times[0].Hour() != 3 && times[0].Minute() != 42 {
		t.Fatal("on sunday at 3:42am isn't correct")
	}


	when = "on December 15"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on December 15 doesn't parse")
	}
	if times[0].Month().String() != "December" && times[0].Day() != 15 {
		t.Fatal("on December 15 isn't correct")
	}


	when = "on jan 12"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on jan 12 doesn't parse")
	}
	if times[0].Month().String() != "January" && times[0].Day() != 12 {
		t.Fatal("on jan 12 isn't correct")
	}


	when = "on July 12th"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on July 12th doesn't parse")
	}
	if times[0].Month().String() != "July" && times[0].Day() != 12 {
		t.Fatal("on July 12th isn't correct")
	}


	when = "on March 22"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on March 22 doesn't parse")
	}
	if times[0].Month().String() != "March" && times[0].Day() != 22 {
		t.Fatal("on March 22 isn't correct")
	}


	when = "on March 17 at 5:41pm"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on March 17 at 5:41pm doesn't parse")
	}
	if times[0].Month().String() != "March" && times[0].Day() != 17 && times[0].Hour() != 17 && times[0].Minute() != 41 {
		t.Fatal("on March 17 at 5:41pm isn't correct")
	}


	when = "on September 7th 2020"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on September 7th 2019 doesn't parse")
	}
	if times[0].Month().String() != "September" && times[0].Day() != 7  {
		t.Fatal("on September 7th 2019 isn't correct")
	}


	when = "on April 17 2020"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on April 17 2020 doesn't parse")
	}
	if times[0].Month().String() != "April" && times[0].Day() != 17  {
		t.Fatal("on April 17 2020 isn't correct")
	}


	when = "on April 9 2020 at 11am"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on April 9 2020 at 11am doesn't parse")
	}
	if times[0].Month().String() != "April" && times[0].Day() != 20 && times[0].Hour() != 11  {
		t.Fatal("on April 9 2020 at 11am isn't correct")
	}


	when = "on auguSt tenth 2019"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on auguSt tenth 2019 doesn't parse")
	}
	if times[0].Month().String() != "August" && times[0].Day() != 10 {
		t.Fatal("on auguSt tenth 2019 isn't correct")
	}


	when = "on 7"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on 7 doesn't parse")
	}
	if times[0].Day() != 7 {
		t.Fatal("on 7 isn't correct")
	}


	when = "on 7th"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on 7th doesn't parse")
	}
	if times[0].Day() != 7 {
		t.Fatal("on 7th isn't correct")
	}


	when = "on seven"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on seven doesn't parse")
	}
	if times[0].Day() != 7 {
		t.Fatal("on seven isn't correct")
	}


	when = "on 1/17/20"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on 1/17/20 doesn't parse")
	}
	if times[0].Year() != 2020 && times[0].Month() != 1 && times[0].Day() != 17 {
		t.Fatal("on 1/17/20 isn't correct")
	}


	when = "on 12/17/2020"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on 12/17/2020 doesn't parse")
	}
	if times[0].Year() != 2020 && times[0].Month() != 12 && times[0].Day() != 17 {
		t.Fatal("on 12/17/2020 isn't correct")
	}


	when = "on 12/1"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on 12/1 doesn't parse")
	}
	if times[0].Month() != 12 && times[0].Day() != 1 {
		t.Fatal("on 12/1 isn't correct")
	}


	when = "on 5-17-20"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on 5-17-20 doesn't parse")
	}
	if times[0].Month() != 5 && times[0].Day() != 17 {
		t.Fatal("on 5-17-20 isn't correct")
	}


	when = "on 12-5-2020"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on 12-5-2020 doesn't parse")
	}
	if times[0].Month() != 12 && times[0].Day() != 5 {
		t.Fatal("on 12-5-2020 isn't correct")
	}


	when = "on 12-12"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on 12-12 doesn't parse")
	}
	if times[0].Month() != 12 && times[0].Day() != 12 {
		t.Fatal("on 12-12 isn't correct")
	}


	when = "on 1-1 at midnight"
	times, iErr = th.App.on(when, user)
	if iErr != nil {
		mlog.Error(iErr.Error())
		t.Fatal("on 1-1 at midnight doesn't parse")
	}
	mlog.Info(fmt.Sprintf("%v",times[0]))
	if times[0].Month() != 1 && times[0].Day() != 1 && times[0].Hour() != 0 {
		t.Fatal("on 1-1 at midnight isn't correct")
	}
	
}

func TestEvery(t *testing.T) {

	/*

        when = "every Thursday";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.THURSDAY)).atTime(9, 0);
        assertEquals(testDate, checkDate);

        when = "every day";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().plusDays(1).atTime(9, 0);
        assertEquals(testDate, checkDate);

        when = "every 12/18";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDateTime.parse("December 18 " + LocalDateTime.now().getYear() + " 09:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        assertEquals(testDate, checkDate);

        when = "every January 25";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDateTime.parse("January 25 " + LocalDateTime.now().getYear() + " 09:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        assertEquals(testDate, checkDate);

        when = "every other Wednesday";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.WEDNESDAY)).plusWeeks(1).atTime(9, 0);
        assertEquals(testDate, checkDate);

        when = "every day at 11:32am";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().plusDays(1).atTime(11, 32);
        assertEquals(testDate, checkDate);

        when = "every 5/5 at 7";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDateTime.parse("May 5 " + LocalDateTime.now().getYear() + " 07:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        checkDate2 = LocalDateTime.parse("May 5 " + LocalDateTime.now().plusYears(1).getYear() + " 07:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        assertTrue(testDate.equals(checkDate) || testDate.equals(checkDate2));

        when = "every 7/20 at 1100";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDateTime.parse("July 20 " + LocalDateTime.now().getYear() + " 11:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        checkDate2 = LocalDateTime.parse("July 20 " + LocalDateTime.now().plusYears(1).getYear() + " 11:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        assertTrue(testDate.equals(checkDate) || testDate.equals(checkDate2));

        when = "every Monday at 7:32am";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.MONDAY)).atTime(7, 32);
        assertEquals(testDate, checkDate);

        when = "every monday and wednesday";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.MONDAY)).atTime(9, 0);
        assertEquals(testDate, checkDate);
        testDate = occurrence.calculate(when).get(1);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.WEDNESDAY)).atTime(9, 0);
        assertEquals(testDate, checkDate);

        when = "every wednesday, thursday";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.WEDNESDAY)).atTime(9, 0);
        assertEquals(testDate, checkDate);
        testDate = occurrence.calculate(when).get(1);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.THURSDAY)).atTime(9, 0);
        assertEquals(testDate, checkDate);

        when = "every other friday and saturday";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.FRIDAY)).plusWeeks(1).atTime(9, 0);
        assertEquals(testDate, checkDate);
        testDate = occurrence.calculate(when).get(1);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.SATURDAY)).plusWeeks(1).atTime(9, 0);
        assertEquals(testDate, checkDate);

        when = "every monday and wednesday at 1:39am";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.MONDAY)).atTime(1, 39);
        assertEquals(testDate, checkDate);
        testDate = occurrence.calculate(when).get(1);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.WEDNESDAY)).atTime(1, 39);
        assertEquals(testDate, checkDate);

        when = "every monday, tuesday and sunday at 11:00";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.MONDAY)).atTime(11, 0);
        assertEquals(testDate, checkDate);
        testDate = occurrence.calculate(when).get(1);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.TUESDAY)).atTime(11, 0);
        assertEquals(testDate, checkDate);
        testDate = occurrence.calculate(when).get(2);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.SUNDAY)).atTime(11, 0);
        assertEquals(testDate, checkDate);


        when = "every monday, tuesday at 2pm";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.MONDAY)).atTime(14, 0);
        assertEquals(testDate, checkDate);
        testDate = occurrence.calculate(when).get(1);
        checkDate = LocalDate.now().with(TemporalAdjusters.next(DayOfWeek.TUESDAY)).atTime(14, 0);
        assertEquals(testDate, checkDate);

        when = "every 1/30 and 9/30 at noon";
        testDate = occurrence.calculate(when).get(0);
        checkDate = LocalDateTime.parse("January 30 " + LocalDateTime.now().getYear() + " 12:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        checkDate2 = LocalDateTime.parse("January 30 " + LocalDateTime.now().plusYears(1).getYear() + " 12:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        assertTrue(testDate.equals(checkDate) || testDate.equals(checkDate2));
        testDate = occurrence.calculate(when).get(1);
        checkDate = LocalDateTime.parse("September 30 " + LocalDateTime.now().getYear() + " 12:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        checkDate2 = LocalDateTime.parse("September 30 " + LocalDateTime.now().plusYears(1).getYear() + " 12:00", new DateTimeFormatterBuilder()
                .parseCaseInsensitive().appendPattern("MMMM d yyyy HH:mm").toFormatter());
        assertTrue(testDate.equals(checkDate) || testDate.equals(checkDate2));

	 */
}
