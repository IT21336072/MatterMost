// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import PostList from './components/post_list.jsx';
import LoadingScreen from 'components/loading_screen.jsx';

import EmojiStore from 'stores/emoji_store.jsx';
import PostStore from 'stores/post_store.jsx';
import UserStore from 'stores/user_store.jsx';
import ChannelStore from 'stores/channel_store.jsx';
import PreferenceStore from 'stores/preference_store.jsx';

import Constants from 'utils/constants.jsx';
const ScrollTypes = Constants.ScrollTypes;

import React from 'react';

export default class PostFocusView extends React.Component {
    constructor(props) {
        super(props);

        this.onChannelChange = this.onChannelChange.bind(this);
        this.onPostsChange = this.onPostsChange.bind(this);
        this.onUserChange = this.onUserChange.bind(this);
        this.onEmojiChange = this.onEmojiChange.bind(this);
        this.onStatusChange = this.onStatusChange.bind(this);
        this.onPreferenceChange = this.onPreferenceChange.bind(this);
        this.onPostListScroll = this.onPostListScroll.bind(this);

        const focusedPostId = PostStore.getFocusedPostId();

        const channel = ChannelStore.getCurrent();
        let profiles = UserStore.getProfiles();
        if (channel && channel.type === Constants.DM_CHANNEL) {
            profiles = Object.assign({}, profiles, UserStore.getDirectProfiles());
        }

        let statuses;
        if (channel && channel.type !== Constants.DM_CHANNEL) {
            statuses = Object.assign({}, UserStore.getStatuses());
        }

        this.state = {
            postList: PostStore.getVisiblePosts(focusedPostId),
            currentUser: UserStore.getCurrentUser(),
            profiles,
            statuses,
            scrollType: ScrollTypes.POST,
            currentChannel: ChannelStore.getCurrentId().slice(),
            scrollPostId: focusedPostId,
            atTop: PostStore.getVisibilityAtTop(focusedPostId),
            atBottom: PostStore.getVisibilityAtBottom(focusedPostId),
            emojis: EmojiStore.getEmojis(),
            flaggedPosts: PreferenceStore.getCategory(Constants.Preferences.CATEGORY_FLAGGED_POST)
        };
    }

    componentDidMount() {
        ChannelStore.addChangeListener(this.onChannelChange);
        PostStore.addChangeListener(this.onPostsChange);
        UserStore.addChangeListener(this.onUserChange);
        UserStore.addStatusesChangeListener(this.onStatusChange);
        EmojiStore.addChangeListener(this.onEmojiChange);
        PreferenceStore.addChangeListener(this.onPreferenceChange);
    }

    componentWillUnmount() {
        ChannelStore.removeChangeListener(this.onChannelChange);
        PostStore.removeChangeListener(this.onPostsChange);
        UserStore.removeChangeListener(this.onUserChange);
        UserStore.removeStatusesChangeListener(this.onStatusChange);
        EmojiStore.removeChangeListener(this.onEmojiChange);
        PreferenceStore.removeChangeListener(this.onPreferenceChange);
    }

    onChannelChange() {
        const currentChannel = ChannelStore.getCurrentId();
        if (this.state.currentChannel !== currentChannel) {
            this.setState({
                currentChannel: currentChannel.slice(),
                scrollType: ScrollTypes.POST
            });
        }
    }

    onPostsChange() {
        const focusedPostId = PostStore.getFocusedPostId();
        if (focusedPostId == null) {
            return;
        }

        this.setState({
            scrollPostId: focusedPostId,
            postList: PostStore.getVisiblePosts(focusedPostId),
            atTop: PostStore.getVisibilityAtTop(focusedPostId),
            atBottom: PostStore.getVisibilityAtBottom(focusedPostId)
        });
    }

    onUserChange() {
        const channel = ChannelStore.getCurrent();
        let profiles = UserStore.getProfiles();
        if (channel && channel.type === Constants.DM_CHANNEL) {
            profiles = Object.assign({}, profiles, UserStore.getDirectProfiles());
        }
        this.setState({currentUser: UserStore.getCurrentUser(), profiles: JSON.parse(JSON.stringify(profiles))});
    }

    onStatusChange() {
        const channel = ChannelStore.getCurrent();
        let statuses;
        if (channel && channel.type !== Constants.DM_CHANNEL) {
            statuses = Object.assign({}, UserStore.getStatuses());
        }

        this.setState({statuses});
    }

    onEmojiChange() {
        this.setState({
            emojis: EmojiStore.getEmojis()
        });
    }

    onPreferenceChange() {
        this.setState({
            flaggedPosts: PreferenceStore.getCategory(Constants.Preferences.CATEGORY_FLAGGED_POST)
        });
    }

    onPostListScroll() {
        this.setState({scrollType: ScrollTypes.FREE});
    }

    render() {
        const postsToHighlight = {};
        postsToHighlight[this.state.scrollPostId] = true;

        let content;
        if (this.state.postList == null) {
            content = (
                <LoadingScreen
                    position='absolute'
                    key='loading'
                />
            );
        } else {
            content = (
                <PostList
                    postList={this.state.postList}
                    currentUser={this.state.currentUser}
                    profiles={this.state.profiles}
                    scrollType={this.state.scrollType}
                    scrollPostId={this.state.scrollPostId}
                    postListScrolled={this.onPostListScroll}
                    showMoreMessagesTop={!this.state.atTop}
                    showMoreMessagesBottom={!this.state.atBottom}
                    postsToHighlight={postsToHighlight}
                    isFocusPost={true}
                    emojis={this.state.emojis}
                    flaggedPosts={this.state.flaggedPosts}
                    statuses={this.state.statuses}
                />
            );
        }

        return (
            <div id='post-list'>
                {content}
            </div>
        );
    }
}
