// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// ***************************************************************
// - [#] indicates a test step (e.g. # Go to a page)
// - [*] indicates an assertion (e.g. * Check the title)
// - Use element ID when selecting an element. Create one if none.
// ***************************************************************

// Stage: @prod
// Group: @channels @system_console

import * as TIMEOUTS from '../../../../../fixtures/timeouts';

describe('User Management', () => {
    let testTeam;
    let testChannel;
    let testUser;
    let testUsersForRoles;
    const roleNames = ['system_manager', 'system_user_manager', 'system_read_only_admin'];

    before(() => {
        cy.apiRequireLicense();

        cy.apiInitSetup().then(({team, channel, user}) => {
            testChannel = channel;
            testTeam = team;
            testUser = user;
        });

        testUsersForRoles = {};

        Cypress._.forEach(roleNames, (roleName) => {
            cy.apiCreateUser().then(({user}) => {
                testUsersForRoles[roleName] = user;
                cy.apiAddUserToTeam(testTeam.id, user.id);
            });
        });
    });

    it('Verify manage user\'s settins option is visible in the user configuration page', () => {
        cy.apiAdminLogin();

        cy.visit('/admin_console/user_management/users');
        gotoUserConfigurationPage(testUser);

        verifyManageUserSettingModal(testUser);

        cy.get('#replyNotificationsTitle').should('be.visible').should('have.text', 'Reply notifications').click();
        cy.get('#notificationCommentsNever').should('be.checked');
        cy.get('#notificationCommentsAny').check();
        cy.get('button#saveSetting').last().scrollIntoView().click();

        cy.apiLogout();
        cy.apiLogin(testUser);

        cy.visit(`/${testTeam.name}/channels/${testChannel.name}`);
        cy.get('[aria-label="Settings"]').click();
        cy.get('#replyNotificationsTitle').should('be.visible').should('have.text', 'Reply notifications').click();
        cy.get('#notificationCommentsAny').should('be.checked');
        cy.apiLogout();
    });

    roleNames.forEach((role) => {
        it(`Verify manage user's settings option is visible for role: ${role}`, () => {
            // TODO: remove below if loop after fixing Bug: https://mattermost.atlassian.net/browse/MM-59376
            if (role !== 'system_manager' && role !== 'system_read_only_admin') {
                // # Make the user a System User Manager
                makeUserASystemRole(testUsersForRoles[role].email, role);

                // * Login as the new user and verify the role permissions (ensure they really are a system user manager)
                verifyAccessToUserSettings(testUsersForRoles[role], role);
            }
        });
    });

    function gotoUserConfigurationPage(user) {
        cy.intercept('**api/v4/reports/users?**').as('getUserList');
        cy.get('#input_searchTerm').clear().type(user.id);
        cy.wait('@getUserList');

        cy.get('#systemUsersTable-cell-0_emailColumn').should('have.text', user.email).click();
        cy.url().should('include', `user_management/user/${user.id}`);
    }

    function verifyManageUserSettingModal(user) {
        cy.get('.manageUserSettingsBtn').should('be.visible').should('have.text', 'Manage User Settings').click();
        cy.get('#confirmModalLabel').should('be.visible').should('have.text', `Manage ${user.nickname}'s Settings`);

        cy.get('#cancelModalButton').should('be.visible').should('have.text', 'Cancel');
        cy.get('#confirmModalButton').should('be.visible').should('have.text', 'Manage User Settings').click();
        cy.get('h1#accountSettingsModalLabel').should('be.visible').should('have.text', `Manage ${user.nickname}'s Settings`);
        cy.get('.adminModeBadge').should('be.visible').should('have.text', 'Admin Mode');
    }

    function makeUserASystemRole(userEmail, role) {
        // # Login as each new role.
        cy.apiAdminLogin();

        // # Go the system console.
        cy.visit('/admin_console/user_management/system_roles');

        cy.get('.admin-console__header').within(() => {
            cy.findByText('System Roles', {timeout: TIMEOUTS.ONE_MIN}).should('exist').and('be.visible');
        });

        // # Click on edit for the role
        cy.findByTestId(`${role}_edit`).click();

        cy.get('button#systemRolePermissionDropdownuser_management_users').click();
        cy.get('div.PermissionSectionDropdownOptions_label').first().should('have.text', 'Can edit').click();

        // # Click Add People button
        cy.findByRole('button', {name: 'Add People'}).click().wait(TIMEOUTS.HALF_SEC);

        // # Type in user name
        cy.findByRole('textbox', {name: 'Search for people'}).typeWithForce(`${userEmail}`);

        // # Find the user and click on him
        cy.get('#multiSelectList').should('be.visible').children().first().click({force: true});

        // # Click add button
        cy.findByRole('button', {name: 'Add'}).click().wait(TIMEOUTS.HALF_SEC);

        // # Click save button
        cy.findByRole('button', {name: 'Save'}).click().wait(TIMEOUTS.HALF_SEC);
        cy.apiLogout();
    }

    function verifyAccessToUserSettings(user) {
        // # Login as each new role.
        cy.apiLogin(user);

        // # Go the system console.
        cy.visit('/admin_console/user_management/users');

        gotoUserConfigurationPage(user);
        verifyManageUserSettingModal(user);
        cy.apiLogout();
    }
});
