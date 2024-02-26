// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import xor from 'lodash.xor';

// *****************************************************************************
// Preferences
// https://api.mattermost.com/#tag/roles
// *****************************************************************************

export const defaultRolesPermissions = {
    channel_admin: 'use_channel_mentions remove_reaction manage_public_channel_members use_group_mentions manage_channel_roles manage_private_channel_members add_reaction read_public_channel_groups create_post read_private_channel_groups',
    channel_guest: 'upload_file edit_post create_post use_channel_mentions read_channel read_channel_content add_reaction remove_reaction',
    channel_user: 'manage_private_channel_members read_public_channel_groups delete_post read_private_channel_groups use_group_mentions manage_private_channel_properties delete_public_channel add_reaction manage_public_channel_properties edit_post upload_file use_channel_mentions get_public_link read_channel read_channel_content delete_private_channel manage_public_channel_members create_post remove_reaction',
    custom_group_user: '',
    playbook_admin: 'playbook_private_manage_properties playbook_public_make_private playbook_public_manage_members playbook_public_manage_roles playbook_public_manage_properties playbook_private_manage_members playbook_private_manage_roles',
    playbook_member: 'playbook_public_view playbook_public_manage_members playbook_public_manage_properties playbook_private_view playbook_private_manage_members playbook_private_manage_properties run_create',
    run_admin: 'run_manage_properties run_manage_members',
    run_member: 'run_view',
    system_admin: 'sysconsole_write_environment_elasticsearch playbook_public_manage_properties sysconsole_write_authentication_ldap run_view manage_jobs manage_roles playbook_public_create manage_public_channel_properties sysconsole_read_plugins delete_post purge_elasticsearch_indexes sysconsole_read_integrations_bot_accounts read_data_retention_job manage_private_channel_members create_elasticsearch_post_indexing_job sysconsole_read_authentication_guest_access create_elasticsearch_post_aggregation_job join_public_teams sysconsole_read_site_public_links add_saml_idp_cert sysconsole_write_site_announcement_banner sysconsole_write_site_notices sysconsole_read_experimental_feature_flags sysconsole_read_site_users_and_teams manage_slash_commands sysconsole_read_authentication_ldap read_channel read_channel_content sysconsole_write_authentication_password list_users_without_team sysconsole_read_authentication_email add_saml_public_cert playbook_private_create promote_guest sysconsole_read_user_management_system_roles manage_public_channel_members create_data_retention_job add_saml_private_cert sysconsole_write_user_management_users sysconsole_read_compliance_compliance_monitoring playbook_public_manage_members sysconsole_write_environment_database sysconsole_write_user_management_teams playbook_private_manage_roles read_public_channel sysconsole_write_plugins sysconsole_read_authentication_openid sysconsole_write_user_management_groups sysconsole_write_site_file_sharing_and_downloads playbook_private_manage_properties sysconsole_read_site_customization join_public_channels add_user_to_team restore_custom_group download_compliance_export_result sysconsole_write_user_management_system_roles sysconsole_write_environment_session_lengths create_custom_group manage_private_channel_properties create_post_public remove_ldap_private_cert sysconsole_write_site_public_links import_team sysconsole_read_environment_developer sysconsole_read_environment_database sysconsole_read_environment_web_server use_channel_mentions view_team remove_others_reactions sysconsole_read_environment_session_lengths sysconsole_write_integrations_bot_accounts playbook_public_view use_group_mentions sysconsole_write_environment_web_server add_ldap_private_cert read_public_channel_groups invite_guest sysconsole_read_environment_smtp create_post sysconsole_read_about_edition_and_license sysconsole_read_authentication_signup sysconsole_read_authentication_saml sysconsole_read_environment_file_storage sysconsole_write_experimental_feature_flags sysconsole_write_site_localization sysconsole_write_environment_rate_limiting sysconsole_read_environment_rate_limiting sysconsole_read_products_boards get_saml_cert_status sysconsole_read_environment_high_availability manage_secure_connections read_compliance_export_job sysconsole_write_compliance_custom_terms_of_service read_user_access_token edit_post sysconsole_write_environment_logging sysconsole_read_environment_push_notification_server sysconsole_write_site_customization read_other_users_teams read_elasticsearch_post_aggregation_job sysconsole_write_compliance_data_retention_policy sysconsole_read_user_management_permissions sysconsole_read_site_emoji sysconsole_read_compliance_data_retention_policy read_license_information sysconsole_read_experimental_features read_deleted_posts sysconsole_read_environment_logging sysconsole_read_reporting_site_statistics test_elasticsearch sysconsole_read_site_posts add_reaction sysconsole_write_authentication_signup manage_outgoing_webhooks create_post_ephemeral sysconsole_read_environment_image_proxy invite_user manage_others_outgoing_webhooks create_user_access_token sysconsole_write_environment_image_proxy sysconsole_write_products_boards read_elasticsearch_post_indexing_job purge_bleve_indexes sysconsole_write_environment_performance_monitoring sysconsole_write_authentication_guest_access sysconsole_read_compliance_custom_terms_of_service edit_others_posts sysconsole_write_billing get_saml_metadata_from_idp sysconsole_write_authentication_saml create_post_bleve_indexes_job invalidate_caches sysconsole_write_experimental_bleve view_members manage_others_bots run_create join_private_teams convert_private_channel_to_public read_audits assign_bot read_jobs remove_user_from_team revoke_user_access_token manage_team sysconsole_read_reporting_server_logs get_public_link manage_others_slash_commands manage_system delete_public_channel read_private_channel_groups sysconsole_read_authentication_mfa delete_emojis list_private_teams create_emojis sysconsole_read_billing sysconsole_write_site_emoji invalidate_email_invite sysconsole_write_environment_file_storage sysconsole_write_compliance_compliance_monitoring remove_saml_public_cert sysconsole_read_compliance_compliance_export sysconsole_read_site_localization manage_team_roles list_public_teams get_logs sysconsole_write_integrations_integration_management sysconsole_read_integrations_cors manage_oauth manage_outgoing_oauth_connections delete_others_emojis sysconsole_write_integrations_gif manage_incoming_webhooks sysconsole_write_authentication_email create_private_channel playbook_private_make_public manage_bots add_ldap_public_cert remove_ldap_public_cert sysconsole_write_site_notifications sysconsole_write_environment_developer playbook_private_manage_members sysconsole_read_user_management_teams edit_custom_group remove_reaction playbook_public_manage_roles sysconsole_write_reporting_server_logs read_others_bots sysconsole_write_site_posts sysconsole_read_site_notifications sysconsole_read_authentication_password playbook_private_view manage_system_wide_oauth get_analytics list_team_channels sysconsole_write_user_management_channels delete_private_channel manage_custom_group_members test_s3 create_ldap_sync_job sysconsole_read_integrations_integration_management test_site_url recycle_database_connections sysconsole_read_site_announcement_banner test_email manage_shared_channels read_bots sysconsole_write_environment_smtp sysconsole_read_experimental_bleve sysconsole_write_environment_push_notification_server sysconsole_write_user_management_permissions sysconsole_read_environment_elasticsearch sysconsole_write_reporting_site_statistics sysconsole_write_site_users_and_teams demote_to_guest create_team test_ldap remove_saml_idp_cert delete_others_posts edit_other_users sysconsole_write_reporting_team_statistics sysconsole_read_integrations_gif sysconsole_read_site_notices sysconsole_write_about_edition_and_license manage_others_incoming_webhooks run_manage_members create_bot sysconsole_write_authentication_mfa sysconsole_read_user_management_users assign_system_admin_role sysconsole_write_experimental_features edit_brand create_group_channel sysconsole_write_authentication_openid create_direct_channel manage_license_information reload_config manage_channel_roles sysconsole_read_user_management_groups create_compliance_export_job read_ldap_sync_job upload_file sysconsole_read_site_file_sharing_and_downloads delete_custom_group sysconsole_read_user_management_channels sysconsole_write_compliance_compliance_export remove_saml_private_cert sysconsole_read_environment_performance_monitoring create_public_channel sysconsole_write_integrations_cors sysconsole_write_environment_high_availability playbook_public_make_private run_manage_properties sysconsole_read_reporting_team_statistics convert_public_channel_to_private',
    system_custom_group_admin: 'create_custom_group edit_custom_group delete_custom_group restore_custom_group manage_custom_group_members',
    system_guest: 'create_group_channel create_direct_channel',
    system_manager: ' sysconsole_read_site_announcement_banner manage_private_channel_properties edit_brand read_private_channel_groups manage_private_channel_members manage_team_roles sysconsole_write_environment_session_lengths sysconsole_read_site_emoji sysconsole_write_environment_developer sysconsole_read_user_management_groups sysconsole_write_user_management_groups sysconsole_write_environment_rate_limiting delete_private_channel sysconsole_read_environment_performance_monitoring sysconsole_read_environment_rate_limiting sysconsole_write_user_management_teams sysconsole_write_integrations_integration_management sysconsole_write_site_public_links sysconsole_read_authentication_ldap sysconsole_write_integrations_cors reload_config sysconsole_write_user_management_channels sysconsole_read_environment_high_availability sysconsole_read_site_users_and_teams sysconsole_read_user_management_teams sysconsole_write_site_users_and_teams sysconsole_read_site_customization sysconsole_write_environment_high_availability sysconsole_read_integrations_bot_accounts sysconsole_read_authentication_guest_access sysconsole_read_site_public_links read_elasticsearch_post_indexing_job sysconsole_read_user_management_channels sysconsole_read_reporting_team_statistics invalidate_caches sysconsole_read_authentication_signup read_elasticsearch_post_aggregation_job sysconsole_write_environment_smtp manage_public_channel_members list_public_teams add_user_to_team sysconsole_read_environment_web_server sysconsole_read_site_localization get_logs sysconsole_write_site_posts sysconsole_write_integrations_bot_accounts sysconsole_write_user_management_permissions sysconsole_read_environment_elasticsearch sysconsole_read_environment_smtp list_private_teams read_public_channel_groups sysconsole_write_environment_file_storage sysconsole_write_integrations_gif manage_public_channel_properties sysconsole_write_environment_performance_monitoring sysconsole_write_site_notifications sysconsole_read_site_notifications sysconsole_read_environment_image_proxy sysconsole_write_site_announcement_banner sysconsole_write_site_emoji test_site_url sysconsole_read_integrations_gif sysconsole_write_environment_logging convert_public_channel_to_private get_analytics sysconsole_read_user_management_permissions sysconsole_write_environment_image_proxy test_elasticsearch recycle_database_connections sysconsole_write_site_localization sysconsole_read_reporting_server_logs create_elasticsearch_post_indexing_job sysconsole_read_reporting_site_statistics test_ldap delete_public_channel sysconsole_write_environment_push_notification_server read_license_information sysconsole_write_products_boards sysconsole_read_about_edition_and_license convert_private_channel_to_public sysconsole_read_integrations_integration_management create_elasticsearch_post_aggregation_job purge_elasticsearch_indexes sysconsole_read_environment_database join_public_teams sysconsole_read_authentication_email sysconsole_read_environment_push_notification_server view_team read_channel sysconsole_read_authentication_password read_ldap_sync_job sysconsole_read_integrations_cors sysconsole_read_environment_logging manage_team sysconsole_read_authentication_openid read_public_channel sysconsole_write_environment_elasticsearch sysconsole_read_plugins manage_channel_roles remove_user_from_team test_email sysconsole_write_site_file_sharing_and_downloads test_s3 sysconsole_read_site_file_sharing_and_downloads sysconsole_read_site_notices sysconsole_read_environment_file_storage join_private_teams sysconsole_read_products_boards sysconsole_read_environment_session_lengths sysconsole_write_environment_database sysconsole_read_authentication_saml sysconsole_read_authentication_mfa sysconsole_write_site_notices sysconsole_write_environment_web_server sysconsole_read_site_posts sysconsole_read_environment_developer sysconsole_write_site_customization manage_outgoing_oauth_connections',
    system_post_all: 'use_group_mentions use_channel_mentions create_post',
    system_post_all_public: 'use_group_mentions use_channel_mentions create_post_public',
    system_read_only_admin: 'sysconsole_read_authentication_guest_access download_compliance_export_result sysconsole_read_compliance_data_retention_policy get_logs sysconsole_read_environment_file_storage read_channel sysconsole_read_integrations_integration_management sysconsole_read_compliance_custom_terms_of_service sysconsole_read_site_notices sysconsole_read_environment_rate_limiting sysconsole_read_about_edition_and_license read_public_channel sysconsole_read_experimental_features test_ldap sysconsole_read_user_management_permissions read_elasticsearch_post_aggregation_job sysconsole_read_environment_image_proxy sysconsole_read_compliance_compliance_export sysconsole_read_integrations_bot_accounts sysconsole_read_authentication_openid sysconsole_read_site_posts sysconsole_read_user_management_users sysconsole_read_experimental_feature_flags sysconsole_read_reporting_team_statistics sysconsole_read_site_localization read_private_channel_groups sysconsole_read_site_file_sharing_and_downloads sysconsole_read_user_management_channels sysconsole_read_authentication_email read_data_retention_job read_audits sysconsole_read_plugins view_team get_analytics sysconsole_read_user_management_groups sysconsole_read_experimental_bleve sysconsole_read_products_boards read_compliance_export_job sysconsole_read_environment_logging sysconsole_read_authentication_signup sysconsole_read_environment_smtp sysconsole_read_environment_session_lengths sysconsole_read_environment_developer sysconsole_read_environment_high_availability read_ldap_sync_job sysconsole_read_environment_performance_monitoring sysconsole_read_authentication_saml read_public_channel_groups sysconsole_read_integrations_gif sysconsole_read_authentication_mfa list_public_teams sysconsole_read_environment_database list_private_teams sysconsole_read_authentication_ldap sysconsole_read_compliance_compliance_monitoring sysconsole_read_site_notifications sysconsole_read_site_announcement_banner read_other_users_teams sysconsole_read_authentication_password sysconsole_read_environment_push_notification_server sysconsole_read_site_users_and_teams sysconsole_read_site_public_links sysconsole_read_site_emoji sysconsole_read_environment_elasticsearch read_license_information sysconsole_read_integrations_cors sysconsole_read_user_management_teams sysconsole_read_reporting_server_logs sysconsole_read_site_customization sysconsole_read_reporting_site_statistics sysconsole_read_environment_web_server read_elasticsearch_post_indexing_job',
    system_user: 'delete_custom_group create_emojis edit_custom_group create_direct_channel view_members join_public_teams restore_custom_group create_custom_group manage_custom_group_members delete_emojis list_public_teams create_team create_group_channel',
    system_user_access_token: 'create_user_access_token read_user_access_token revoke_user_access_token',
    system_user_manager: 'sysconsole_read_authentication_password sysconsole_read_authentication_openid sysconsole_write_user_management_groups list_private_teams sysconsole_read_user_management_groups sysconsole_read_authentication_email manage_public_channel_properties delete_private_channel sysconsole_read_authentication_signup read_private_channel_groups sysconsole_read_user_management_teams test_ldap read_channel view_team manage_team sysconsole_write_user_management_teams manage_channel_roles sysconsole_read_authentication_saml sysconsole_read_authentication_guest_access convert_private_channel_to_public sysconsole_read_user_management_permissions join_public_teams sysconsole_write_user_management_channels read_public_channel_groups sysconsole_read_user_management_channels list_public_teams manage_team_roles join_private_teams manage_public_channel_members convert_public_channel_to_private remove_user_from_team sysconsole_read_authentication_ldap manage_private_channel_properties delete_public_channel manage_private_channel_members read_public_channel add_user_to_team sysconsole_read_authentication_mfa read_ldap_sync_job',
    team_admin: 'manage_others_slash_commands manage_channel_roles manage_others_outgoing_webhooks manage_team_roles use_channel_mentions manage_incoming_webhooks manage_slash_commands manage_public_channel_members convert_private_channel_to_public manage_private_channel_members manage_team convert_public_channel_to_private use_group_mentions delete_post read_public_channel_groups delete_others_posts playbook_private_manage_roles add_reaction remove_reaction remove_user_from_team read_private_channel_groups manage_outgoing_webhooks create_post playbook_public_manage_roles import_team manage_others_incoming_webhooks',
    team_guest: 'view_team',
    team_post_all: 'create_post use_channel_mentions use_group_mentions',
    team_post_all_public: 'create_post_public use_channel_mentions use_group_mentions',
    team_user: 'add_user_to_team view_team playbook_private_create playbook_public_create invite_user join_public_channels list_team_channels read_public_channel create_private_channel create_public_channel',
};

Cypress.Commands.add('getRoleByName', (name) => {
    return cy.request({
        headers: {'X-Requested-With': 'XMLHttpRequest'},
        url: `/api/v4/roles/name/${name}`,
        method: 'GET',
    }).then((response) => {
        expect(response.status).to.equal(200);
        return cy.wrap({name: response.body});
    });
});

Cypress.Commands.add('apiGetRolesByNames', (names) => {
    return cy.request({
        headers: {'X-Requested-With': 'XMLHttpRequest'},
        url: '/api/v4/roles/names',
        method: 'POST',
        body: names || Object.keys(defaultRolesPermissions),
    }).then((response) => {
        expect(response.status).to.equal(200);
        return cy.wrap({roles: response.body});
    });
});

Cypress.Commands.add('apiPatchRole', (roleID, patch) => {
    return cy.request({
        headers: {'X-Requested-With': 'XMLHttpRequest'},
        url: `/api/v4/roles/${roleID}/patch`,
        method: 'PUT',
        body: patch,
    }).then((response) => {
        expect(response.status).to.equal(200);
        return cy.wrap({role: response.body});
    });
});

Cypress.Commands.add('apiResetRoles', () => {
    cy.apiGetRolesByNames().then(({roles}) => {
        roles.forEach((role) => {
            const defaultPermissions = defaultRolesPermissions[role.name].split(' ');
            const diff = xor(role.permissions, defaultPermissions)?.filter((p) => p?.length);

            if (diff?.length > 0) {
                cy.apiPatchRole(role.id, {permissions: defaultPermissions});
            }
        });
    });
});
