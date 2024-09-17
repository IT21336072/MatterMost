// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import type {IntlShape} from 'react-intl';
import {FormattedMessage, defineMessages} from 'react-intl';

import type {PasswordConfig} from 'mattermost-redux/selectors/entities/general';

import Constants from 'utils/constants';

export function isValidPassword(password: string, passwordConfig: PasswordConfig, intl?: IntlShape) {
    const telemetryErrorIds = [];
    let valid = true;
    const minimumLength = passwordConfig.minimumLength || Constants.MIN_PASSWORD_LENGTH;

    if (password.length < minimumLength || password.length > Constants.MAX_PASSWORD_LENGTH) {
        valid = false;
        telemetryErrorIds.push({field: 'password', rule: 'error_length'});
    }

    if (passwordConfig.requireLowercase) {
        if (!password.match(/[a-z]/)) {
            valid = false;
        }

        telemetryErrorIds.push({field: 'password', rule: 'lowercase'});
    }

    if (passwordConfig.requireUppercase) {
        if (!password.match(/[A-Z]/)) {
            valid = false;
        }

        telemetryErrorIds.push({field: 'password', rule: 'uppercase'});
    }

    if (passwordConfig.requireNumber) {
        if (!password.match(/[0-9]/)) {
            valid = false;
        }

        telemetryErrorIds.push({field: 'password', rule: 'number'});
    }

    if (passwordConfig.requireSymbol) {
        if (!password.match(/[ !"\\#$%&'()*+,-./:;<=>?@[\]^_`|~]/)) {
            valid = false;
        }

        telemetryErrorIds.push({field: 'password', rule: 'symbol'});
    }

    let error;
    if (!valid) {
        const errorMessage = getPasswordErrorMessage(
            passwordConfig.requireLowercase,
            passwordConfig.requireUppercase,
            passwordConfig.requireNumber,
            passwordConfig.requireSymbol,
        );

        error = intl ? (
            intl.formatMessage(
                errorMessage,
                {
                    min: minimumLength,
                    max: Constants.MAX_PASSWORD_LENGTH,
                },
            )
        ) : (
            <FormattedMessage
                {...errorMessage}
                values={{
                    min: minimumLength,
                    max: Constants.MAX_PASSWORD_LENGTH,
                }}
            />
        );
    }

    return {valid, error, telemetryErrorIds};
}

export function getPasswordErrorMessage(needsLowercase: boolean, needsUppercase: boolean, needsNumber: boolean, needsSymbol: boolean) {
    /* eslint-disable no-else-return, no-lonely-if */
    if (needsLowercase) {
        if (needsUppercase) {
            if (needsNumber) {
                if (needsSymbol) {
                    return passwordErrors.passwordErrorLowercaseUppercaseNumberSymbol;
                } else {
                    return passwordErrors.passwordErrorLowercaseUppercaseNumber;
                }
            } else {
                if (needsSymbol) {
                    return passwordErrors.passwordErrorLowercaseUppercaseSymbol;
                } else {
                    return passwordErrors.passwordErrorLowercaseUppercase;
                }
            }
        } else {
            if (needsNumber) {
                if (needsSymbol) {
                    return passwordErrors.passwordErrorLowercaseNumberSymbol;
                } else {
                    return passwordErrors.passwordErrorLowercaseNumber;
                }
            } else {
                if (needsSymbol) {
                    return passwordErrors.passwordErrorLowercaseSymbol;
                } else {
                    return passwordErrors.passwordErrorLowercase;
                }
            }
        }
    } else {
        if (needsUppercase) {
            if (needsNumber) {
                if (needsSymbol) {
                    return passwordErrors.passwordErrorUppercaseNumberSymbol;
                } else {
                    return passwordErrors.passwordErrorUppercaseNumber;
                }
            } else {
                if (needsSymbol) {
                    return passwordErrors.passwordErrorUppercaseSymbol;
                } else {
                    return passwordErrors.passwordErrorUppercase;
                }
            }
        } else {
            if (needsNumber) {
                if (needsSymbol) {
                    return passwordErrors.passwordErrorNumberSymbol;
                } else {
                    return passwordErrors.passwordErrorNumber;
                }
            } else {
                if (needsSymbol) {
                    return passwordErrors.passwordErrorSymbol;
                } else {
                    return passwordErrors.passwordError;
                }
            }
        }
    }
    /* eslint-enable no-else-return, no-lonely-if */
}

export const passwordErrors = defineMessages({
    passwordError: {id: 'user.settings.security.passwordError', defaultMessage: 'Must be {min}-{max} characters long.'},
    passwordErrorLowercase: {id: 'user.settings.security.passwordErrorLowercase', defaultMessage: 'Must be {min}-{max} characters long and include lowercase letters.'},
    passwordErrorLowercaseNumber: {id: 'user.settings.security.passwordErrorLowercaseNumber', defaultMessage: 'Must be {min}-{max} characters long and include lowercase letters and numbers.'},
    passwordErrorLowercaseNumberSymbol: {id: 'user.settings.security.passwordErrorLowercaseNumberSymbol', defaultMessage: 'Must be {min}-{max} characters long and include lowercase letters, numbers, and special characters.'},
    passwordErrorLowercaseSymbol: {id: 'user.settings.security.passwordErrorLowercaseSymbol', defaultMessage: 'Must be {min}-{max} characters long and include lowercase letters and special characters.'},
    passwordErrorLowercaseUppercase: {id: 'user.settings.security.passwordErrorLowercaseUppercase', defaultMessage: 'Must be {min}-{max} characters long and include both lowercase and uppercase letters.'},
    passwordErrorLowercaseUppercaseNumber: {id: 'user.settings.security.passwordErrorLowercaseUppercaseNumber', defaultMessage: 'Must be {min}-{max} characters long and include both lowercase and uppercase letters, and numbers.'},
    passwordErrorLowercaseUppercaseNumberSymbol: {id: 'user.settings.security.passwordErrorLowercaseUppercaseNumberSymbol', defaultMessage: 'Must be {min}-{max} characters long and include both lowercase and uppercase letters, numbers, and special characters.'},
    passwordErrorLowercaseUppercaseSymbol: {id: 'user.settings.security.passwordErrorLowercaseUppercaseSymbol', defaultMessage: 'Must be {min}-{max} characters long and include both lowercase and uppercase letters, and special characters.'},
    passwordErrorNumber: {id: 'user.settings.security.passwordErrorNumber', defaultMessage: 'Must be {min}-{max} characters long and include numbers.'},
    passwordErrorNumberSymbol: {id: 'user.settings.security.passwordErrorNumberSymbol', defaultMessage: 'Must be {min}-{max} characters long and include numbers and special characters.'},
    passwordErrorSymbol: {id: 'user.settings.security.passwordErrorSymbol', defaultMessage: 'Must be {min}-{max} characters long and include special characters.'},
    passwordErrorUppercase: {id: 'user.settings.security.passwordErrorUppercase', defaultMessage: 'Must be {min}-{max} characters long and include uppercase letters.'},
    passwordErrorUppercaseNumber: {id: 'user.settings.security.passwordErrorUppercaseNumber', defaultMessage: 'Must be {min}-{max} characters long and include uppercase letters, and numbers.'},
    passwordErrorUppercaseNumberSymbol: {id: 'user.settings.security.passwordErrorUppercaseNumberSymbol', defaultMessage: 'Must be {min}-{max} characters long and include uppercase letters, numbers, and special characters.'},
    passwordErrorUppercaseSymbol: {id: 'user.settings.security.passwordErrorUppercaseSymbol', defaultMessage: 'Must be {min}-{max} characters long and include uppercase letters, and special characters.'},
});
