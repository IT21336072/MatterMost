// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {AnyAction} from 'redux';
import {combineReducers} from 'redux';

import type {ScheduledPost, ScheduledPostsState} from '@mattermost/types/schedule_post';

import {ScheduledPostTypes, UserTypes} from 'mattermost-redux/action_types';

function byTeamId(state: ScheduledPostsState['byTeamId'] = {}, action: AnyAction) {
    switch (action.type) {
    case ScheduledPostTypes.SCHEDULED_POSTS_RECEIVED: {
        const {scheduledPostsByTeamId} = action.data;
        const newState = {...state};

        Object.keys(scheduledPostsByTeamId).forEach((teamId: string) => {
            if (scheduledPostsByTeamId.hasOwnProperty(teamId)) {
                newState[teamId] = scheduledPostsByTeamId[teamId];
            }
        });

        return newState;
    }
    case ScheduledPostTypes.SINGLE_SCHEDULED_POST_RECEIVED: {
        const scheduledPost = action.data.scheduledPost;
        const teamId = action.data.teamId || 'directChannels';

        const newState = {...state};

        if (newState[teamId]) {
            newState[teamId] = [...newState[teamId], scheduledPost];
        } else {
            newState[teamId] = [scheduledPost];
        }

        return newState;
    }
    case UserTypes.LOGOUT_SUCCESS:
        return {};
    default:
        return state;
    }
}

function hasErrorByTeamId(state: ScheduledPostsState['errorsByTeamId'] = {}, action: AnyAction) {
    switch (action.type) {
    case ScheduledPostTypes.SCHEDULED_POSTS_RECEIVED: {
        const {scheduledPostsByTeamId} = action.data;
        const newState = {...state};

        Object.keys(scheduledPostsByTeamId).forEach((teamId: string) => {
            if (scheduledPostsByTeamId.hasOwnProperty(teamId)) {
                const teamScheduledPosts = scheduledPostsByTeamId[teamId] as ScheduledPost[];
                newState[teamId] = teamScheduledPosts.some((scheduledPost) => scheduledPost.error_code);
            }
        });

        return newState;
    }
    case ScheduledPostTypes.SINGLE_SCHEDULED_POST_RECEIVED: {
        const teamId = action.data.teamId || 'directChannels';

        // if team already has an error state, then irrespective of what's the error
        // on new scheduled post, the team would still have an error.
        // So nothing changes so we return the original state as-in.
        if (state[teamId]) {
            return state;
        }

        const scheduledPost = action.data.scheduledPost as ScheduledPost;

        // if team doesn't have any error and neither does the new scheduled post,
        // then nothing changes so we return the original state as-in.
        if (!scheduledPost.error_code && state[teamId]) {
            return state;
        }

        const newState = {...state};
        newState[teamId] = newState[teamId] || Boolean(scheduledPost.error_code);

        return newState;
    }
    case UserTypes.LOGOUT_SUCCESS: {
        return {};
    }
    default:
        return state;
    }
}

export default combineReducers({
    byTeamId,
    hasErrorByTeamId,
});
