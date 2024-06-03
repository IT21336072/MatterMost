// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import type {HTMLAttributes} from 'react';
import React, {forwardRef, useRef} from 'react';
import {useDispatch, useSelector} from 'react-redux';
import {Link} from 'react-router-dom';
import styled, {css} from 'styled-components';

import type {ChannelBookmark} from '@mattermost/types/channel_bookmarks';
import type {FileInfo} from '@mattermost/types/files';
import type {Post} from '@mattermost/types/posts';

import {getFile} from 'mattermost-redux/selectors/entities/files';
import {getFileDownloadUrl} from 'mattermost-redux/utils/file_utils';

import {openModal} from 'actions/views/modals';

import ExternalLink from 'components/external_link';
import FilePreviewModal from 'components/file_preview_modal';

import {ModalIdentifiers} from 'utils/constants';
import {getSiteURL, shouldOpenInNewTab} from 'utils/url';

import type {GlobalState} from 'types/store';

import BookmarkItemDotMenu from './bookmark_dot_menu';
import BookmarkIcon from './bookmark_icon';

type Props = {bookmark: ChannelBookmark};
const BookmarkItem = <T extends HTMLAnchorElement>({bookmark}: Props) => {
    const linkRef = useRef<T>(null);
    const dispatch = useDispatch();
    const fileInfo: FileInfo | undefined = useSelector((state: GlobalState) => (bookmark?.file_id && getFile(state, bookmark.file_id)) || undefined);

    const open = () => {
        linkRef.current?.click();
    };

    const handleOpenFile = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        if (fileInfo) {
            dispatch(openModal({
                modalId: ModalIdentifiers.FILE_PREVIEW_MODAL,
                dialogType: FilePreviewModal,
                dialogProps: {
                    post: {} as Post,
                    fileInfos: [fileInfo],
                    startIndex: 0,
                },
            }));
        }
    };

    const icon = (
        <BookmarkIcon
            type={bookmark.type}
            emoji={bookmark.emoji}
            imageUrl={bookmark.image_url}
            fileInfo={fileInfo}
        />
    );
    let link;

    if (bookmark.type === 'link' && bookmark.link_url) {
        link = (
            <DynamicLink
                href={bookmark.link_url}
                ref={linkRef}
                isFile={false}
            >
                {icon}
                <Label>{bookmark.display_name}</Label>
            </DynamicLink>
        );
    } else if (bookmark.type === 'file' && bookmark.file_id) {
        link = (
            <DynamicLink
                href={getFileDownloadUrl(bookmark.file_id)}
                onClick={handleOpenFile}
                ref={linkRef}
                isFile={true}
            >
                {icon}
                <Label>{bookmark.display_name}</Label>
            </DynamicLink>
        );
    }

    return (
        <Chip>
            {link}
            <BookmarkItemDotMenu
                bookmark={bookmark}
                open={open}
            />
        </Chip>
    );
};

const Chip = styled.div`
    position: relative;
    border-radius: 12px;
    overflow: hidden;
    margin: 1px 0;
    flex-shrink: 0;
    min-width: 5rem;
    max-width: 25rem;

    button {
        position: absolute;
        visibility: hidden;
        right: 6px;
        top: 3px;
    }

    &:hover,
    &:focus-within,
    &:has([aria-expanded="true"]) {
        button {
            visibility: visible;
        }
    }

    &:hover,
    &:focus-within {
        a {
            text-decoration: none;
        }
    }

    &:hover,
    &:focus-within,
    &:has([aria-expanded="true"]) {
        a {
            background: rgba(var(--center-channel-color-rgb), 0.08);
            color: rgba(var(--center-channel-color-rgb), 1);
        }
    }

    &:active:not(:has(button:active)),
    &--active,
    &--active:hover {
        a {
            background: rgba(var(--button-bg-rgb), 0.08);
            color: rgb(var(--button-bg-rgb)) !important;

            .icon__text {
                color: rgb(var(--button-bg));
            }

            .icon {
                color: rgb(var(--button-bg));
            }
        }

    }
`;

const Label = styled.span`
    white-space: nowrap;
    padding: 4px 0;
    text-overflow: ellipsis;
    overflow: hidden;
`;

const TARGET_BLANK_URL_PREFIX = '!';

type DynamicLinkProps = {href: string; children: React.ReactNode; isFile: boolean; onClick?: HTMLAttributes<HTMLAnchorElement>['onClick']};
const DynamicLink = forwardRef<HTMLAnchorElement, DynamicLinkProps>(({href, children, isFile, onClick}, ref) => {
    const siteURL = getSiteURL();
    const openInNewTab = shouldOpenInNewTab(href, siteURL);

    const prefixed = href[0] === TARGET_BLANK_URL_PREFIX;

    if (prefixed || openInNewTab) {
        return (
            <StyledExternalLink
                href={prefixed ? href.substring(1) : href}
                rel='noopener noreferrer'
                target='_blank'
                location='channel_bookmarks.item'
                ref={ref}
            >
                {children}
            </StyledExternalLink>
        );
    }

    if (href.startsWith(siteURL) && !isFile) {
        return (
            <StyledLink
                to={href.slice(siteURL.length)}
                ref={ref}
                onClick={onClick}
            >
                {children}
            </StyledLink>
        );
    }

    return (
        <StyledAnchor
            href={href}
            ref={ref}
            onClick={onClick}
        >
            {children}
        </StyledAnchor>
    );
});

const linkStyles = css`
    display: flex;
    padding: 0 12px 0 6px;
    gap: 5px;

    color: rgba(var(--center-channel-color-rgb), 1);
    font-family: Open Sans;
    font-size: 12px;
    font-style: normal;
    font-weight: 600;
    line-height: 16px;
`;

const StyledAnchor = styled.a`
    &&&& {
        ${linkStyles}
    }
`;

const StyledLink = styled(Link)`
    &&&& {
        ${linkStyles}
    }
`;

const StyledExternalLink = styled(ExternalLink)`
    &&&& {
        ${linkStyles}
    }
`;

export default BookmarkItem;
