// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import {intlShape, injectIntl, defineMessages} from 'react-intl';
import SettingItemMin from '../setting_item_min.jsx';
import SettingItemMax from '../setting_item_max.jsx';
import SettingPicture from '../setting_picture.jsx';

import UserStore from '../../stores/user_store.jsx';
import ErrorStore from '../../stores/error_store.jsx';

import * as Client from '../../utils/client.jsx';
import Constants from '../../utils/constants.jsx';
import * as AsyncClient from '../../utils/async_client.jsx';
import * as Utils from '../../utils/utils.jsx';

const messages = defineMessages({
    usernameReserved: {
        id: 'user.settings.general.usernameReserved',
        defaultMessage: 'This username is reserved, please choose a new one.'
    },
    usernameRestrictions: {
        id: 'user.settings.general.usernameRestrictions',
        defaultMessage: "Username must begin with a letter, and contain between 3 to 15 lowercase characters made up of numbers, letters, and the symbols '.', '-' and '_'."
    },
    newUsername: {
        id: 'user.settings.general.newUsername',
        defaultMessage: 'You must submit a new username'
    },
    newNickname: {
        id: 'user.settings.general.newNickname',
        defaultMessage: 'You must submit a new nickname'
    },
    newName: {
        id: 'user.settings.general.newName',
        defaultMessage: 'You must submit a new first or last name'
    },
    validEmail: {
        id: 'user.settings.general.validEmail',
        defaultMessage: 'Please enter a valid email address'
    },
    validImage: {
        id: 'user.settings.general.validImage',
        defaultMessage: 'Only JPG or PNG images may be used for profile pictures'
    },
    firstName: {
        id: 'user.settings.general.firstName',
        defaultMessage: 'First Name'
    },
    lastName: {
        id: 'user.settings.general.lastName',
        defaultMessage: 'Last Name'
    },
    notificationsLink: {
        id: 'user.settings.general.notificationsLink',
        defaultMessage: 'Notifications'
    },
    notificationsExtra1: {
        id: 'user.settings.general.notificationsExtra1',
        defaultMessage: 'By default, you will receive mention notifications when someone types your first name. '
    },
    notificationsExtra2: {
        id: 'user.settings.general.notificationsExtra2',
        defaultMessage: 'Go to '
    },
    notificationsExtra3: {
        id: 'user.settings.general.notificationsExtra3',
        defaultMessage: 'settings to change this default.'
    },
    fullName: {
        id: 'user.settings.general.fullName',
        defaultMessage: 'Full Name'
    },
    nickname: {
        id: 'user.settings.general.nickname',
        defaultMessage: 'Nickname'
    },
    nicknameExtra1: {
        id: 'user.settings.general.nicknameExtra1',
        defaultMessage: 'Use Nickname for a name you might be called that is different from your first name and user name.'
    },
    nicknameExtra2: {
        id: 'user.settings.general.nicknameExtra2',
        defaultMessage: 'This is most often used when two or more people have similar sounding names and usernames.'
    },
    username: {
        id: 'user.settings.general.username',
        defaultMessage: 'Username'
    },
    usernameInfo: {
        id: 'user.settings.general.usernameInfo',
        defaultMessage: 'Pick something easy for teammates to recognize and recall.'
    },
    emailHelp1: {
        id: 'user.settings.general.emailHelp1',
        defaultMessage: 'Email is used for sign-in, notifications, and password reset. Email requires verification if changed.'
    },
    emailHelp2: {
        id: 'user.settings.general.emailHelp2',
        defaultMessage: 'Email has been disabled by your system administrator. No notification emails will be sent until it is enabled.'
    },
    emailHelp3: {
        id: 'user.settings.general.emailHelp3',
        defaultMessage: 'Email is used for sign-in, notifications, and password reset.'
    },
    emailHelp4: {
        id: 'user.settings.general.emailHelp4',
        defaultMessage: 'A verification email was sent to '
    },
    primaryEmail: {
        id: 'user.settings.general.primaryEmail',
        defaultMessage: 'Primary Email'
    },
    confirmEmail: {
        id: 'user.settings.general.confirmEmail',
        defaultMessage: 'Confirm Email'
    },
    email: {
        id: 'user.settings.general.email',
        defaultMessage: 'Email'
    },
    profilePicture: {
        id: 'user.settings.general.profilePicture',
        defaultMessage: 'Profile Picture'
    },
    uploadImage: {
        id: 'user.settings.general.uploadImage',
        defaultMessage: 'Click \'Edit\' to upload an image.'
    },
    imageUpdated: {
        id: 'user.settings.general.imageUpdated',
        defaultMessage: 'Image last updated '
    },
    title: {
        id: 'user.settings.general.title',
        defaultMessage: 'General Settings'
    },
    emailMatch: {
        id: 'user.settings.general.emailMatch',
        defaultMessage: 'The new emails you entered do not match.'
    },
    checkEmail1: {
        id: 'user.settings.general.checkEmail1',
        defaultMessage: 'Check your email at '
    },
    checkEmail2: {
        id: 'user.settings.general.checkEmail2',
        defaultMessage: ' to verify the address.'
    },
    emailCantUpdate: {
        id: 'user.settings.general.emailCantUpdate',
        defaultMessage: 'Log in occurs through SSO. Email cannot be updated.'
    },
    newAddress1: {
        id: 'user.settings.general.newAddress1',
        defaultMessage: 'New Address: '
    },
    newAddress2: {
        id: 'user.settings.general.newAddress2',
        defaultMessage: '\nCheck your email to verify the above address.'
    },
    newAddress3: {
        id: 'user.settings.general.newAddress3',
        defaultMessage: 'Check your email to verify your new address'
    },
    newAddress4: {
        id: 'user.settings.general.newAddress4',
        defaultMessage: 'Log in done through SSO'
    },
    close: {
        id: 'user.settings.general.close',
        defaultMessage: 'Close'
    },
    imageTooLarge: {
        id: 'user.settings.general.imageTooLarge',
        defaultMessage: 'Unable to upload profile image. File is too large.'
    }
});

class UserSettingsGeneralTab extends React.Component {
    constructor(props) {
        super(props);
        this.submitActive = false;

        this.submitUsername = this.submitUsername.bind(this);
        this.submitNickname = this.submitNickname.bind(this);
        this.submitName = this.submitName.bind(this);
        this.submitEmail = this.submitEmail.bind(this);
        this.submitUser = this.submitUser.bind(this);
        this.submitPicture = this.submitPicture.bind(this);

        this.updateUsername = this.updateUsername.bind(this);
        this.updateFirstName = this.updateFirstName.bind(this);
        this.updateLastName = this.updateLastName.bind(this);
        this.updateNickname = this.updateNickname.bind(this);
        this.updateEmail = this.updateEmail.bind(this);
        this.updateConfirmEmail = this.updateConfirmEmail.bind(this);
        this.updatePicture = this.updatePicture.bind(this);
        this.updateSection = this.updateSection.bind(this);

        this.state = this.setupInitialState(props);
    }
    submitUsername(e) {
        e.preventDefault();

        const {formatMessage} = this.props.intl;

        const user = Object.assign({}, this.props.user);
        const username = this.state.username.trim().toLowerCase();

        const usernameError = Utils.isValidUsername(username);
        if (usernameError === 'Cannot use a reserved word as a username.') {
            this.setState({clientError: formatMessage(messages.usernameReserved)});
            return;
        } else if (usernameError) {
            this.setState({clientError: formatMessage(messages.usernameRestrictions)});
            return;
        }

        if (user.username === username) {
            this.updateSection('');
            return;
        }

        user.username = username;

        this.submitUser(user, false);
    }
    submitNickname(e) {
        e.preventDefault();

        const user = Object.assign({}, this.props.user);
        const nickname = this.state.nickname.trim();

        if (user.nickname === nickname) {
            this.updateSection('');
            return;
        }

        user.nickname = nickname;

        this.submitUser(user, false);
    }
    submitName(e) {
        e.preventDefault();

        const user = Object.assign({}, this.props.user);
        const firstName = this.state.firstName.trim();
        const lastName = this.state.lastName.trim();

        if (user.first_name === firstName && user.last_name === lastName) {
            this.updateSection('');
            return;
        }

        user.first_name = firstName;
        user.last_name = lastName;

        this.submitUser(user, false);
    }
    submitEmail(e) {
        e.preventDefault();

        const {formatMessage} = this.props.intl;

        const user = Object.assign({}, this.props.user);
        const email = this.state.email.trim().toLowerCase();
        const confirmEmail = this.state.confirmEmail.trim().toLowerCase();

        if (email === '' || !Utils.isEmail(email)) {
            this.setState({emailError: formatMessage(messages.validEmail), clientError: '', serverError: ''});
            return;
        }

        if (email !== confirmEmail) {
            this.setState({emailError: formatMessage(messages.emailMatch), clientError: '', serverError: ''});
            return;
        }

        if (user.email === email) {
            this.updateSection('');
            return;
        }

        user.email = email;
        this.submitUser(user, true);
    }
    submitUser(user, emailUpdated) {
        const {formatMessage} = this.props.intl;
        Client.updateUser(user,
            () => {
                this.updateSection('');
                AsyncClient.getMe();
                const verificationEnabled = global.window.mm_config.SendEmailNotifications === 'true' && global.window.mm_config.RequireEmailVerification === 'true' && emailUpdated;

                if (verificationEnabled) {
                    ErrorStore.storeLastError({message: formatMessage(messages.checkEmail1) + user.email + formatMessage(messages.checkEmail2)});
                    ErrorStore.emitChange();
                    this.setState({emailChangeInProgress: true});
                }
            },
            (err) => {
                let serverError;
                if (err.message) {
                    serverError = err.message;
                } else {
                    serverError = err;
                }
                this.setState({serverError, emailError: '', clientError: ''});
            }
        );
    }
    submitPicture(e) {
        e.preventDefault();

        const {formatMessage} = this.props.intl;
        if (!this.state.picture) {
            return;
        }

        if (!this.submitActive) {
            return;
        }

        const picture = this.state.picture;

        if (picture.type !== 'image/jpeg' && picture.type !== 'image/png') {
            this.setState({clientError: formatMessage(messages.validImage)});
            return;
        } else if (picture.size > Constants.MAX_FILE_SIZE) {
            this.setState({clientError: formatMessage(messages.imageTooLarge)});
            return;
        }

        var formData = new FormData();
        formData.append('image', picture, picture.name);
        this.setState({loadingPicture: true});

        Client.uploadProfileImage(formData,
            () => {
                this.submitActive = false;
                AsyncClient.getMe();
                window.location.reload();
            },
            (err) => {
                var state = this.setupInitialState(this.props);
                state.serverError = err.message;
                this.setState(state);
            }
        );
    }
    updateUsername(e) {
        this.setState({username: e.target.value});
    }
    updateFirstName(e) {
        this.setState({firstName: e.target.value});
    }
    updateLastName(e) {
        this.setState({lastName: e.target.value});
    }
    updateNickname(e) {
        this.setState({nickname: e.target.value});
    }
    updateEmail(e) {
        this.setState({email: e.target.value});
    }
    updateConfirmEmail(e) {
        this.setState({confirmEmail: e.target.value});
    }
    updatePicture(e) {
        if (e.target.files && e.target.files[0]) {
            this.setState({picture: e.target.files[0]});

            this.submitActive = true;
            this.setState({clientError: null});
        } else {
            this.setState({picture: null});
        }
    }
    updateSection(section) {
        const emailChangeInProgress = this.state.emailChangeInProgress;
        this.setState(Object.assign({}, this.setupInitialState(this.props), {emailChangeInProgress, clientError: '', serverError: '', emailError: ''}));
        this.submitActive = false;
        this.props.updateSection(section);
    }
    setupInitialState(props) {
        const user = props.user;

        return {username: user.username, firstName: user.first_name, lastName: user.last_name, nickname: user.nickname,
                        email: user.email, confirmEmail: '', picture: null, loadingPicture: false, emailChangeInProgress: false, language: user.language};
    }
    render() {
        const {formatMessage, locale} = this.props.intl;
        const user = this.props.user;

        let clientError = null;
        if (this.state.clientError) {
            clientError = this.state.clientError;
        }
        let serverError = null;
        if (this.state.serverError) {
            serverError = this.state.serverError;
        }
        let emailError = null;
        if (this.state.emailError) {
            emailError = this.state.emailError;
        }

        let nameSection;
        const inputs = [];

        if (this.props.activeSection === 'name') {
            inputs.push(
                <div
                    key='firstNameSetting'
                    className='form-group'
                >
                    <label className='col-sm-5 control-label'>{formatMessage(messages.firstName)}</label>
                    <div className='col-sm-7'>
                        <input
                            className='form-control'
                            type='text'
                            onChange={this.updateFirstName}
                            value={this.state.firstName}
                        />
                    </div>
                </div>
            );

            inputs.push(
                <div
                    key='lastNameSetting'
                    className='form-group'
                >
                    <label className='col-sm-5 control-label'>{formatMessage(messages.lastName)}</label>
                    <div className='col-sm-7'>
                        <input
                            className='form-control'
                            type='text'
                            onChange={this.updateLastName}
                            value={this.state.lastName}
                        />
                    </div>
                </div>
            );

            function notifClick(e) {
                e.preventDefault();
                this.updateSection('');
                this.props.updateTab('notifications');
            }

            const notifLink = (
                <a
                    href='#'
                    onClick={notifClick.bind(this)}
                >
                    {formatMessage(messages.notificationsLink)}
                </a>
            );

            const extraInfo = (
                <span>
                    {formatMessage(messages.notificationsExtra1)}
                    {formatMessage(messages.notificationsExtra2)} {notifLink} {formatMessage(messages.notificationsExtra3)}
                </span>
            );

            nameSection = (
                <SettingItemMax
                    title={formatMessage(messages.fullName)}
                    inputs={inputs}
                    submit={this.submitName}
                    server_error={serverError}
                    client_error={clientError}
                    updateSection={(e) => {
                        this.updateSection('');
                        e.preventDefault();
                    }}
                    extraInfo={extraInfo}
                />
            );
        } else {
            let fullName = '';

            if (user.first_name && user.last_name) {
                fullName = user.first_name + ' ' + user.last_name;
            } else if (user.first_name) {
                fullName = user.first_name;
            } else if (user.last_name) {
                fullName = user.last_name;
            }

            nameSection = (
                <SettingItemMin
                    title={formatMessage(messages.fullName)}
                    describe={fullName}
                    updateSection={() => {
                        this.updateSection('name');
                    }}
                />
            );
        }

        let nicknameSection;
        if (this.props.activeSection === 'nickname') {
            let nicknameLabel = formatMessage(messages.nickname);
            if (Utils.isMobile()) {
                nicknameLabel = '';
            }

            inputs.push(
                <div
                    key='nicknameSetting'
                    className='form-group'
                >
                    <label className='col-sm-5 control-label'>{nicknameLabel}</label>
                    <div className='col-sm-7'>
                        <input
                            className='form-control'
                            type='text'
                            onChange={this.updateNickname}
                            value={this.state.nickname}
                        />
                    </div>
                </div>
            );

            const extraInfo = (
                <span>
                    {formatMessage(messages.nicknameExtra1)} {formatMessage(messages.nicknameExtra2)}
                </span>
            );

            nicknameSection = (
                <SettingItemMax
                    title={formatMessage(messages.nickname)}
                    inputs={inputs}
                    submit={this.submitNickname}
                    server_error={serverError}
                    client_error={clientError}
                    updateSection={(e) => {
                        this.updateSection('');
                        e.preventDefault();
                    }}
                    extraInfo={extraInfo}
                />
            );
        } else {
            nicknameSection = (
                <SettingItemMin
                    title={formatMessage(messages.nickname)}
                    describe={UserStore.getCurrentUser().nickname}
                    updateSection={() => {
                        this.updateSection('nickname');
                    }}
                />
            );
        }

        let usernameSection;
        if (this.props.activeSection === 'username') {
            let usernameLabel = formatMessage(messages.username);
            if (Utils.isMobile()) {
                usernameLabel = '';
            }

            inputs.push(
                <div
                    key='usernameSetting'
                    className='form-group'
                >
                    <label className='col-sm-5 control-label'>{usernameLabel}</label>
                    <div className='col-sm-7'>
                        <input
                            className='form-control'
                            type='text'
                            onChange={this.updateUsername}
                            value={this.state.username}
                        />
                    </div>
                </div>
            );

            const extraInfo = (<span>{formatMessage(messages.usernameInfo)}</span>);

            usernameSection = (
                <SettingItemMax
                    title={formatMessage(messages.username)}
                    inputs={inputs}
                    submit={this.submitUsername}
                    server_error={serverError}
                    client_error={clientError}
                    updateSection={(e) => {
                        this.updateSection('');
                        e.preventDefault();
                    }}
                    extraInfo={extraInfo}
                />
            );
        } else {
            usernameSection = (
                <SettingItemMin
                    title={formatMessage(messages.username)}
                    describe={UserStore.getCurrentUser().username}
                    updateSection={() => {
                        this.updateSection('username');
                    }}
                />
            );
        }

        let emailSection;
        if (this.props.activeSection === 'email') {
            const emailEnabled = global.window.mm_config.SendEmailNotifications === 'true';
            const emailVerificationEnabled = global.window.mm_config.RequireEmailVerification === 'true';
            let helpText = formatMessage(messages.emailHelp1);

            if (!emailEnabled) {
                helpText = <div className='setting-list__hint text-danger'>{formatMessage(messages.emailHelp2)}</div>;
            } else if (!emailVerificationEnabled) {
                helpText = formatMessage(messages.emailHelp3);
            } else if (this.state.emailChangeInProgress) {
                const newEmail = UserStore.getCurrentUser().email;
                if (newEmail) {
                    helpText = formatMessage(messages.emailHelp4) + newEmail + '.';
                }
            }

            let submit = null;

            if (this.props.user.auth_service === '') {
                inputs.push(
                    <div key='emailSetting'>
                        <div className='form-group'>
                            <label className='col-sm-5 control-label'>{formatMessage(messages.primaryEmail)}</label>
                            <div className='col-sm-7'>
                                <input
                                    className='form-control'
                                    type='text'
                                    onChange={this.updateEmail}
                                    value={this.state.email}
                                />
                            </div>
                        </div>
                    </div>
                );

                inputs.push(
                    <div key='confirmEmailSetting'>
                        <div className='form-group'>
                            <label className='col-sm-5 control-label'>{formatMessage(messages.confirmEmail)}</label>
                            <div className='col-sm-7'>
                                <input
                                    className='form-control'
                                    type='text'
                                    onChange={this.updateConfirmEmail}
                                    value={this.state.confirmEmail}
                                />
                            </div>
                        </div>
                        {helpText}
                    </div>
                );

                submit = this.submitEmail;
            } else {
                inputs.push(
                    <div
                        key='oauthEmailInfo'
                        className='form-group'
                    >
                        <div className='setting-list__hint'>{formatMessage(messages.emailCantUpdate)}</div>
                        {helpText}
                    </div>
                );
            }

            emailSection = (
                <SettingItemMax
                    title={formatMessage(messages.email)}
                    inputs={inputs}
                    submit={submit}
                    server_error={serverError}
                    client_error={emailError}
                    updateSection={(e) => {
                        this.updateSection('');
                        e.preventDefault();
                    }}
                />
            );
        } else {
            let describe = '';
            if (this.props.user.auth_service === '') {
                if (this.state.emailChangeInProgress) {
                    const newEmail = UserStore.getCurrentUser().email;
                    if (newEmail) {
                        describe = formatMessage(messages.newAddress1) + newEmail + formatMessage(messages.newAddress2);
                    } else {
                        describe = formatMessage(messages.newAddress3);
                    }
                } else {
                    describe = UserStore.getCurrentUser().email;
                }
            } else {
                describe = formatMessage(messages.newAddress4);
            }

            emailSection = (
                <SettingItemMin
                    title={formatMessage(messages.email)}
                    describe={describe}
                    updateSection={() => {
                        this.updateSection('email');
                    }}
                />
            );
        }

        let pictureSection;
        if (this.props.activeSection === 'picture') {
            pictureSection = (
                <SettingPicture
                    title={formatMessage(messages.profilePicture)}
                    submit={this.submitPicture}
                    src={'/api/v1/users/' + user.id + '/image?time=' + user.last_picture_update + '&' + Utils.getSessionIndex()}
                    server_error={serverError}
                    client_error={clientError}
                    updateSection={(e) => {
                        this.updateSection('');
                        e.preventDefault();
                    }}
                    picture={this.state.picture}
                    pictureChange={this.updatePicture}
                    submitActive={this.submitActive}
                    loadingPicture={this.state.loadingPicture}
                />
            );
        } else {
            let minMessage = formatMessage(messages.uploadImage);
            if (user.last_picture_update) {
                minMessage = formatMessage(messages.imageUpdated) + Utils.displayDate(user.last_picture_update, locale);
            }
            pictureSection = (
                <SettingItemMin
                    title={formatMessage(messages.profilePicture)}
                    describe={minMessage}
                    updateSection={() => {
                        this.updateSection('picture');
                    }}
                />
            );
        }

        return (
            <div>
                <div className='modal-header'>
                    <button
                        type='button'
                        className='close'
                        data-dismiss='modal'
                        aria-label={formatMessage(messages.close)}
                        onClick={this.props.closeModal}
                    >
                        <span aria-hidden='true'>{'×'}</span>
                    </button>
                    <h4
                        className='modal-title'
                        ref='title'
                    >
                        <i
                            className='modal-back'
                            onClick={this.props.collapseModal}
                        />
                        {formatMessage(messages.title)}
                    </h4>
                </div>
                <div className='user-settings'>
                    <h3 className='tab-header'>{formatMessage(messages.title)}</h3>
                    <div className='divider-dark first'/>
                    {nameSection}
                    <div className='divider-light'/>
                    {usernameSection}
                    <div className='divider-light'/>
                    {nicknameSection}
                    <div className='divider-light'/>
                    {emailSection}
                    <div className='divider-light'/>
                    {pictureSection}
                    <div className='divider-dark'/>
                </div>
            </div>
        );
    }
}

UserSettingsGeneralTab.propTypes = {
    user: React.PropTypes.object.isRequired,
    updateSection: React.PropTypes.func.isRequired,
    updateTab: React.PropTypes.func.isRequired,
    activeSection: React.PropTypes.string.isRequired,
    closeModal: React.PropTypes.func.isRequired,
    collapseModal: React.PropTypes.func.isRequired,
    intl: intlShape.isRequired
};

export default injectIntl(UserSettingsGeneralTab);