// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

import {intlShape, injectIntl, defineMessages} from 'react-intl';
import EmailToSSO from './email_to_sso.jsx';
import SSOToEmail from './sso_to_email.jsx';

const messages = defineMessages({
    noEmail: {
        id: 'claim.account.noEmail',
        defaultMessage: 'No email specified.'
    }
});

class ClaimAccount extends React.Component {
    constructor(props) {
        super(props);

        this.state = {};
    }
    render() {
        const {formatMessage} = this.props.intl;
        let content;
        if (this.props.email === '') {
            content = <p>{formatMessage(messages.noEmail)}</p>;
        } else if (this.props.currentType === '' && this.props.newType !== '') {
            content = (
                <EmailToSSO
                    email={this.props.email}
                    type={this.props.newType}
                    teamName={this.props.teamName}
                    teamDisplayName={this.props.teamDisplayName}
                />
            );
        } else {
            content = (
                <SSOToEmail
                    email={this.props.email}
                    currentType={this.props.currentType}
                    teamName={this.props.teamName}
                    teamDisplayName={this.props.teamDisplayName}
                />
            );
        }

        return (
            <div>
                {content}
            </div>
        );
    }
}

ClaimAccount.defaultProps = {
};
ClaimAccount.propTypes = {
    currentType: React.PropTypes.string.isRequired,
    newType: React.PropTypes.string.isRequired,
    email: React.PropTypes.string.isRequired,
    teamName: React.PropTypes.string.isRequired,
    teamDisplayName: React.PropTypes.string.isRequired,
    intl: intlShape.isRequired
};

export default injectIntl(ClaimAccount);