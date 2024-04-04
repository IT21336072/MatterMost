// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

import {
    noBillingHistory,
    freeTrial,
} from './billing_summary';

import './billing_summary.scss';

type BillingSummaryProps = {
    isFreeTrial: boolean;
    daysLeftOnTrial: number;
}

const BillingSummary = ({isFreeTrial, daysLeftOnTrial}: BillingSummaryProps) => {
    let body = noBillingHistory;

    if (isFreeTrial) {
        body = freeTrial(daysLeftOnTrial);
    }
    return (
        <div className='BillingSummary'>
            {body}
        </div>
    );
};

export default BillingSummary;
