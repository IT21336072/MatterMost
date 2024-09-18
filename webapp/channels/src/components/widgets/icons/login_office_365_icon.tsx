// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {useIntl} from 'react-intl';

export default function EntraIdIcon(props: React.HTMLAttributes<HTMLSpanElement>) {
    const {formatMessage} = useIntl();

    return (
        <span {...props}>
            <svg
                width='16'
                height='16'
                viewBox='0 0 16 16'
                fill='none'
                xmlns='http://www.w3.org/2000/svg'
                aria-label={formatMessage({id: 'generic_icons.login.oneLogin', defaultMessage: 'One Login Icon'})}
            >
                <path
                    d='M3.38 12.473c0.345 0.215 0.918 0.454 1.524 0.454 0.552 0 1.065 -0.16 1.49 -0.433l0.002 -0.001L8 11.491V15.111a1.387 1.387 0 0 1 -0.732 -0.208z'
                    fill='#225086'
                />
                <path
                    d='m6.98 1.34 -6.667 7.52c-0.515 0.581 -0.38 1.46 0.287 1.876l2.779 1.737c0.345 0.215 0.918 0.454 1.524 0.454 0.552 0 1.065 -0.16 1.49 -0.433l0.002 -0.001L8 11.491l-3.879 -2.425 3.88 -4.377V0.889c-0.377 0 -0.753 0.15 -1.02 0.451Z'
                    fill='#6df'
                />
                <path
                    points='4.636 10.199 4.688 10.231 9 12.927 9.001 12.927 9.001 12.927 9.001 5.276 9 5.275 4.636 10.199'
                    fill='#cbf8ff'
                    d='M4.121 9.066L4.167 9.094L8 11.491L8.001 11.491L8.001 11.491L8.001 4.69L8 4.689L4.121 9.066Z'
                />
                <path
                    d='M15.399 10.736c0.668 -0.417 0.802 -1.295 0.287 -1.876l-4.374 -4.934a2.756 2.756 0 0 0 -1.167 -0.259c-0.822 0 -1.557 0.355 -2.046 0.912l-0.097 0.109 3.879 4.377 -3.88 2.425v3.62c0.255 0 0.509 -0.069 0.732 -0.208l6.667 -4.167Z'
                    fill='#074793'
                />
                <path
                    d='M8.001 0.889v3.8l0.097 -0.109a2.711 2.711 0 0 1 2.046 -0.912c0.42 0 0.814 0.095 1.167 0.259l-2.292 -2.586A1.351 1.351 0 0 0 8 0.89Z'
                    fill='#0294e4'
                />
                <path
                    points='13.365 10.199 13.365 10.199 13.365 10.199 9.001 5.276 9.001 12.926 13.365 10.199'
                    fill='#96bcc2'
                    d='M11.88 9.066L11.88 9.066L11.88 9.066L8.001 4.69L8.001 11.49L11.88 9.066Z'
                />
            </svg>
        </span>
    );
}
