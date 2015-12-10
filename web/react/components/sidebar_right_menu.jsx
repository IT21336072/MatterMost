// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import {intlShape, injectIntl, defineMessages} from 'react-intl';

import TeamMembersModal from './team_members_modal.jsx';
import ToggleModalButton from './toggle_modal_button.jsx';
import UserSettingsModal from './user_settings/user_settings_modal.jsx';
import UserStore from '../stores/user_store.jsx';
import * as client from '../utils/client.jsx';
import * as EventHelpers from '../dispatcher/event_helpers.jsx';
import * as utils from '../utils/utils.jsx';

const messages = defineMessages({
    inviteNew: {
        id: 'sidebar_right_menu.inviteNew',
        defaultMessage: 'Invite New Member'
    },
    teamLink: {
        id: 'sidebar_right_menu.teamLink',
        defaultMessage: 'Get Team Invite Link'
    },
    teamSettings: {
        id: 'sidebar_right_menu.teamSettings',
        defaultMessage: 'Team Settings'
    },
    manageTeams: {
        id: 'sidebar_right_menu.manageTeams',
        defaultMessage: 'Manage Members'
    },
    console: {
        id: 'sidebar_right_menu.console',
        defaultMessage: 'System Console'
    },
    accountSettings: {
        id: 'sidebar_right_menu.accountSettings',
        defaultMessage: 'Account Settings'
    },
    logout: {
        id: 'sidebar_right_menu.logout',
        defaultMessage: 'Logout'
    },
    help: {
        id: 'sidebar_right_menu.help',
        defaultMessage: 'Help'
    },
    report: {
        id: 'sidebar_right_menu.report',
        defaultMessage: 'Report a Problem'
    }
});

class SidebarRightMenu extends React.Component {
    componentDidMount() {
        $('.sidebar--left .dropdown-menu').perfectScrollbar();
    }

    constructor(props) {
        super(props);

        this.handleLogoutClick = this.handleLogoutClick.bind(this);

        this.state = {
            showUserSettingsModal: false
        };
    }

    handleLogoutClick(e) {
        e.preventDefault();
        client.logout();
    }

    render() {
        const {formatMessage} = this.props.intl;
        var teamLink = '';
        var inviteLink = '';
        var teamSettingsLink = '';
        var manageLink = '';
        var consoleLink = '';
        var currentUser = UserStore.getCurrentUser();
        var isAdmin = false;
        var isSystemAdmin = false;

        if (currentUser != null) {
            isAdmin = utils.isAdmin(currentUser.roles);
            isSystemAdmin = utils.isSystemAdmin(currentUser.roles);

            inviteLink = (
                <li>
                    <a
                        href='#'
                        onClick={EventHelpers.showInviteMemberModal}
                    >
                        <i className='fa fa-user'></i>{formatMessage(messages.inviteNew)}
                    </a>
                </li>
            );

            if (this.props.teamType === 'O') {
                teamLink = (
                    <li>
                        <a
                            href='#'
                            onClick={EventHelpers.showGetTeamInviteLinkModal}
                        >
                            <i className='glyphicon glyphicon-link'></i>{formatMessage(messages.teamLink)}
                        </a>
                    </li>
                );
            }
        }

        if (isAdmin) {
            teamSettingsLink = (
                <li>
                    <a
                        href='#'
                        data-toggle='modal'
                        data-target='#team_settings'
                    ><i className='fa fa-globe'></i>{formatMessage(messages.teamSettings)}</a>
                </li>
            );
            manageLink = (
                <li>
                    <ToggleModalButton dialogType={TeamMembersModal}>
                        <i className='fa fa-users'></i>{formatMessage(messages.manageTeams)}
                    </ToggleModalButton>
                </li>
            );
        }

        if (isSystemAdmin && !utils.isMobile()) {
            consoleLink = (
                <li>
                    <a
                        href={'/admin_console?' + utils.getSessionIndex()}
                    >
                    <i className='fa fa-wrench'></i>{formatMessage(messages.console)}</a>
                </li>
            );
        }

        var siteName = '';
        if (global.window.mm_config.SiteName != null) {
            siteName = global.window.mm_config.SiteName;
        }
        var teamDisplayName = siteName;
        if (this.props.teamDisplayName) {
            teamDisplayName = this.props.teamDisplayName;
        }

        return (
            <div>
                <div className='team__header theme'>
                    <a
                        className='team__name'
                        href='/channels/general'
                    >{teamDisplayName}</a>
                </div>

                <div className='nav-pills__container'>
                    <ul className='nav nav-pills nav-stacked'>
                        <li>
                            <a
                                href='#'
                                onClick={() => this.setState({showUserSettingsModal: true})}
                            >
                                <i className='fa fa-cog'></i>{formatMessage(messages.accountSettings)}
                            </a>
                        </li>
                        {teamSettingsLink}
                        {inviteLink}
                        {teamLink}
                        {manageLink}
                        {consoleLink}
                        <li>
                            <a
                                href='#'
                                onClick={this.handleLogoutClick}
                            ><i className='fa fa-sign-out'></i>{formatMessage(messages.logout)}</a></li>
                        <li className='divider'></li>
                        <li>
                            <a
                                target='_blank'
                                href='http://ayuda.zboxapp.com/collection/65-chat'
                            ><i className='fa fa-question'></i>{formatMessage(messages.help)}</a></li>
                        <li>
                            <a
                                target='_blank'
                                href='http://ayuda.zboxapp.com/collection/65-chat#contactModal'
                            ><i className='fa fa-phone'></i>{formatMessage(messages.report)}</a></li>
                    </ul>
                </div>
                <UserSettingsModal
                    show={this.state.showUserSettingsModal}
                    onModalDismissed={() => this.setState({showUserSettingsModal: false})}
                />
            </div>
        );
    }
}

SidebarRightMenu.propTypes = {
    intl: intlShape.isRequired,
    teamType: React.PropTypes.string,
    teamDisplayName: React.PropTypes.string
};

export default injectIntl(SidebarRightMenu);