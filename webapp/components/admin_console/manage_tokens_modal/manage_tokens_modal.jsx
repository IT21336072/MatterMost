// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import LoadingScreen from 'components/loading_screen.jsx';
import RevokeTokenButton from 'components/admin_console/revoke_token_button';

import {Client4} from 'mattermost-redux/client';
import * as UserUtils from 'mattermost-redux/utils/user_utils';

import React from 'react';
import {Modal} from 'react-bootstrap';
import PropTypes from 'prop-types';
import {FormattedMessage, FormattedHTMLMessage} from 'react-intl';

export default class ManageTokensModal extends React.PureComponent {
    static propTypes = {

        /**
         * Set to render the modal
         */
        show: PropTypes.bool.isRequired,

        /**
         * The user the roles are being managed for
         */
        user: PropTypes.object,

        /**
         * The user access tokens for a user, object with token ids as keys
         */
        userAccessTokens: PropTypes.object,

        /**
         * Function called when modal is dismissed
         */
        onModalDismissed: PropTypes.func.isRequired,

        actions: PropTypes.shape({

            /**
             * Function to get a user's access tokens
             */
            getUserAccessTokensForUser: PropTypes.func.isRequired
        }).isRequired
    };

    constructor(props) {
        super(props);
        this.state = {error: null};
    }

    componentWillReceiveProps(nextProps) {
        const userId = this.props.user ? this.props.user.id : null;
        const nextUserId = nextProps.user ? nextProps.user.id : null;
        if (nextUserId && nextUserId !== userId) {
            this.props.actions.getUserAccessTokensForUser(nextUserId, 0, 200);
        }
    }

    handleError = (error) => {
        this.setState({
            error
        });
    }

    renderContents = () => {
        const {user, userAccessTokens} = this.props;

        if (!user) {
            return <LoadingScreen/>;
        }

        let name = UserUtils.getFullName(user);
        if (name) {
            name += ` (@${user.username})`;
        } else {
            name = `@${user.username}`;
        }

        let tokenList;
        if (userAccessTokens) {
            const userAccessTokensList = Object.values(userAccessTokens);

            if (userAccessTokensList.length === 0) {
                tokenList = (
                    <FormattedMessage
                        id='admin.manage_tokens.userAccessTokensNone'
                        defaultMessage='No user access tokens.'
                    />
                );
            } else {
                tokenList = userAccessTokensList.map((token) => {
                    return (
                        <div
                            key={token.id}
                            className='manage-teams__team'
                        >
                            <div className='manage-teams__team-name'>
                                {token.description + ' - ' + token.id}
                            </div>
                            <div className='manage-teams__team-actions'>
                                <RevokeTokenButton
                                    tokenId={token.id}
                                    onError={this.handleError}
                                />
                            </div>
                        </div>
                    );
                });
            }
        } else {
            tokenList = <LoadingScreen/>;
        }

        return (
            <div>
                <div className='manage-teams__user'>
                    <img
                        className='manage-teams__profile-picture'
                        src={Client4.getProfilePictureUrl(user.id, user.last_picture_update)}
                    />
                    <div className='manage-teams__info'>
                        <div className='manage-teams__name'>
                            {name}
                        </div>
                        <div className='manage-teams__email'>
                            {user.email}
                        </div>
                    </div>
                </div>
                <FormattedHTMLMessage
                    id='admin.manage_tokens.userAccessTokensDescription'
                    defaultMessage='Read about <a href="https://about.mattermost.com/default-user-access-tokens" target="_blank">user access tokens</a>.'
                />
                <div className='manage-teams__teams'>
                    {tokenList}
                </div>
            </div>
        );
    }

    render() {
        return (
            <Modal
                show={this.props.show}
                onHide={this.props.onModalDismissed}
                dialogClassName='manage-teams'
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        <FormattedMessage
                            id='admin.manage_tokens.manageTokensTitle'
                            defaultMessage='Manage User Access Tokens'
                        />
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    {this.renderContents()}
                    {this.state.error}
                </Modal.Body>
            </Modal>
        );
    }
}
