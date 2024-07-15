// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';
import type {Dispatch} from 'redux';
import timezones from 'timezones.json';

import {CollapsedThreads} from '@mattermost/types/config';

import {savePreferences} from 'mattermost-redux/actions/preferences';
import {autoUpdateTimezone} from 'mattermost-redux/actions/timezone';
import {patchUser, updateMe} from 'mattermost-redux/actions/users';
import {getConfig, getLicense} from 'mattermost-redux/selectors/entities/general';
import {
    get,
    isCollapsedThreadsAllowed,
    getCollapsedThreadsPreference,
    getFromPreferences,
} from 'mattermost-redux/selectors/entities/preferences';
import {
    generateCurrentTimezoneLabel,
    getCurrentTimezoneFull,
    getCurrentTimezoneLabel,
    getTimezoneForUserProfile,
} from 'mattermost-redux/selectors/entities/timezone';
import {getCurrentUserId, getUser} from 'mattermost-redux/selectors/entities/users';
import {getUserCurrentTimezone} from 'mattermost-redux/utils/timezone_utils';

import {getLanguages, isLanguageAvailable} from 'i18n/i18n';
import {Preferences} from 'utils/constants';

import type {GlobalState} from 'types/store';

import type {OwnProps} from './user_settings_display';
import UserSettingsDisplay from './user_settings_display';

export function makeMapStateToProps() {
    return (state: GlobalState, props: OwnProps) => {
        const config = getConfig(state);
        const currentUserId = getCurrentUserId(state);
        const userTimezone = props.adminMode ? getTimezoneForUserProfile(props.user) : getCurrentTimezoneFull(state);
        const automaticTimezoneNotSet = userTimezone && userTimezone.useAutomaticTimezone && !userTimezone.automaticTimezone;
        const shouldAutoUpdateTimezone = !userTimezone || automaticTimezoneNotSet;
        const timezoneLabel = props.adminMode ? generateCurrentTimezoneLabel(getUserCurrentTimezone(userTimezone)) : getCurrentTimezoneLabel(state);
        const allowCustomThemes = config.AllowCustomThemes === 'true';
        const enableLinkPreviews = config.EnableLinkPreviews === 'true';
        const enableThemeSelection = config.EnableThemeSelection === 'true';
        const lockTeammateNameDisplay = getLicense(state).LockTeammateNameDisplay === 'true' && config.LockTeammateNameDisplay === 'true';
        const configTeammateNameDisplay = config.TeammateNameDisplay as string;
        const emojiPickerEnabled = config.EnableEmojiPicker === 'true';
        const lastActiveTimeEnabled = config.EnableLastActiveTime === 'true';

        let lastActiveDisplay = true;
        const user = props.adminMode ? props.user : getUser(state, currentUserId);
        if (user.props?.show_last_active === 'false') {
            lastActiveDisplay = false;
        }

        let userLocale = props.user.locale;
        if (!isLanguageAvailable(state, userLocale)) {
            userLocale = config.DefaultClientLocale as string;
        }

        let getPreference = (prefCategory: string, prefName: string, defaultValue: string) => get(state, prefCategory, prefName, defaultValue);
        if (props.adminMode && props.userPreferences) {
            const preferences = props.userPreferences;
            getPreference = (prefCategory: string, prefName: string, defaultValue: string) => getFromPreferences(preferences, prefCategory, prefName, defaultValue);
        }

        return {
            lockTeammateNameDisplay,
            allowCustomThemes,
            configTeammateNameDisplay,
            enableLinkPreviews,
            locales: getLanguages(state),
            userLocale,
            enableThemeSelection,
            timezones,
            timezoneLabel,
            userTimezone,
            shouldAutoUpdateTimezone,
            currentUserTimezone: getUserCurrentTimezone(userTimezone) as string,
            availabilityStatusOnPosts: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.AVAILABILITY_STATUS_ON_POSTS, Preferences.AVAILABILITY_STATUS_ON_POSTS_DEFAULT),
            militaryTime: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.USE_MILITARY_TIME, Preferences.USE_MILITARY_TIME_DEFAULT),
            teammateNameDisplay: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.NAME_NAME_FORMAT, configTeammateNameDisplay),
            channelDisplayMode: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.CHANNEL_DISPLAY_MODE, Preferences.CHANNEL_DISPLAY_MODE_DEFAULT),
            messageDisplay: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.MESSAGE_DISPLAY, Preferences.MESSAGE_DISPLAY_DEFAULT),
            colorizeUsernames: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.COLORIZE_USERNAMES, Preferences.COLORIZE_USERNAMES_DEFAULT),
            collapseDisplay: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.COLLAPSE_DISPLAY, Preferences.COLLAPSE_DISPLAY_DEFAULT),
            collapsedReplyThreadsAllowUserPreference: isCollapsedThreadsAllowed(state) && getConfig(state).CollapsedThreads !== CollapsedThreads.ALWAYS_ON,
            collapsedReplyThreads: getCollapsedThreadsPreference(state),
            clickToReply: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.CLICK_TO_REPLY, Preferences.CLICK_TO_REPLY_DEFAULT),
            linkPreviewDisplay: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.LINK_PREVIEW_DISPLAY, Preferences.LINK_PREVIEW_DISPLAY_DEFAULT),
            oneClickReactionsOnPosts: getPreference(Preferences.CATEGORY_DISPLAY_SETTINGS, Preferences.ONE_CLICK_REACTIONS_ENABLED, Preferences.ONE_CLICK_REACTIONS_ENABLED_DEFAULT),
            emojiPickerEnabled,
            lastActiveDisplay,
            lastActiveTimeEnabled,
        };
    };
}

function mapDispatchToProps(dispatch: Dispatch) {
    return {
        actions: bindActionCreators({
            autoUpdateTimezone,
            savePreferences,
            updateMe,
            patchUser,
        }, dispatch),
    };
}

export default connect(makeMapStateToProps, mapDispatchToProps)(UserSettingsDisplay);
