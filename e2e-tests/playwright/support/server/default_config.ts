// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import merge from 'deepmerge';

import {
    AdminConfig,
    ExperimentalSettings,
    PasswordSettings,
    ServiceSettings,
    TeamSettings,
    PluginSettings,
    ClusterSettings,
    CollapsedThreads,
} from '@mattermost/types/config';
import testConfig from '@e2e-test.config';

export function getOnPremServerConfig(): AdminConfig {
    return merge<AdminConfig>(defaultServerConfig, onPremServerConfig() as AdminConfig);
}

type TestAdminConfig = {
    ClusterSettings: Partial<ClusterSettings>;
    ExperimentalSettings: Partial<ExperimentalSettings>;
    PasswordSettings: Partial<PasswordSettings>;
    PluginSettings: Partial<PluginSettings>;
    ServiceSettings: Partial<ServiceSettings>;
    TeamSettings: Partial<TeamSettings>;
};

// On-prem setting that is different from the default
const onPremServerConfig = (): Partial<TestAdminConfig> => {
    return {
        ClusterSettings: {
            Enable: testConfig.haClusterEnabled,
            ClusterName: testConfig.haClusterName,
        },
        ExperimentalSettings: {
            DisableAppBar: false,
        },
        PasswordSettings: {
            MinimumLength: 5,
            Lowercase: false,
            Number: false,
            Uppercase: false,
            Symbol: false,
            EnableForgotLink: true,
        },
        PluginSettings: {
            EnableUploads: true,
            Plugins: {
                'com.mattermost.calls': {
                    defaultenabled: true,
                },
            },
        },
        ServiceSettings: {
            SiteURL: testConfig.baseURL,
            EnableOnboardingFlow: false,
        },
        TeamSettings: {
            EnableOpenServer: true,
        },
    };
};

// Should be based only from the generated default config from ./server via "make config-reset"
// Based on v7.10 server
const defaultServerConfig: AdminConfig = {
    ServiceSettings: {
        SiteURL: '',
        WebsocketURL: '',
        LicenseFileLocation: '',
        ListenAddress: ':8065',
        ConnectionSecurity: '',
        TLSCertFile: '',
        TLSKeyFile: '',
        TLSMinVer: '1.2',
        TLSStrictTransport: false,
        TLSStrictTransportMaxAge: 63072000,
        TLSOverwriteCiphers: [],
        UseLetsEncrypt: false,
        LetsEncryptCertificateCacheFile: './config/letsencrypt.cache',
        Forward80To443: false,
        TrustedProxyIPHeader: [],
        ReadTimeout: 300,
        WriteTimeout: 300,
        IdleTimeout: 60,
        MaximumLoginAttempts: 10,
        GoroutineHealthThreshold: -1,
        EnableOAuthServiceProvider: true,
        EnableIncomingWebhooks: true,
        EnableOutgoingWebhooks: true,
        EnableCommands: true,
        EnablePostUsernameOverride: false,
        EnablePostIconOverride: false,
        GoogleDeveloperKey: '',
        EnableLinkPreviews: true,
        EnablePermalinkPreviews: true,
        RestrictLinkPreviews: '',
        EnableTesting: false,
        EnableDeveloper: false,
        DeveloperFlags: '',
        EnableClientPerformanceDebugging: false,
        EnableOpenTracing: false,
        EnableSecurityFixAlert: true,
        EnableInsecureOutgoingConnections: false,
        AllowedUntrustedInternalConnections: '',
        EnableMultifactorAuthentication: false,
        EnforceMultifactorAuthentication: false,
        EnableUserAccessTokens: false,
        AllowCorsFrom: '',
        CorsExposedHeaders: '',
        CorsAllowCredentials: false,
        CorsDebug: false,
        AllowCookiesForSubdomains: false,
        ExtendSessionLengthWithActivity: true,
        SessionLengthWebInDays: 30,
        SessionLengthWebInHours: 720,
        SessionLengthMobileInDays: 30,
        SessionLengthMobileInHours: 720,
        SessionLengthSSOInDays: 30,
        SessionLengthSSOInHours: 720,
        SessionCacheInMinutes: 10,
        SessionIdleTimeoutInMinutes: 43200,
        WebsocketSecurePort: 443,
        WebsocketPort: 80,
        WebserverMode: 'gzip',
        EnableGifPicker: true,
        GfycatAPIKey: '2_KtH_W5',
        GfycatAPISecret: '3wLVZPiswc3DnaiaFoLkDvB4X0IV6CpMkj4tf2inJRsBY6-FnkT08zGmppWFgeof',
        GiphySdkKey: 's0glxvzVg9azvPipKxcPLpXV0q1x1fVP',
        EnableCustomEmoji: true,
        EnableEmojiPicker: true,
        PostEditTimeLimit: -1,
        TimeBetweenUserTypingUpdatesMilliseconds: 5000,
        EnablePostSearch: true,
        EnableFileSearch: true,
        MinimumHashtagLength: 3,
        EnableUserTypingMessages: true,
        EnableChannelViewedMessages: true,
        EnableUserStatuses: true,
        ExperimentalEnableAuthenticationTransfer: true,
        ClusterLogTimeoutMilliseconds: 2000,
        EnablePreviewFeatures: true,
        EnableTutorial: true,
        EnableOnboardingFlow: true,
        ExperimentalEnableDefaultChannelLeaveJoinMessages: true,
        ExperimentalGroupUnreadChannels: 'disabled',
        EnableAPITeamDeletion: false,
        EnableAPITriggerAdminNotifications: false,
        EnableAPIUserDeletion: false,
        ExperimentalEnableHardenedMode: false,
        ExperimentalStrictCSRFEnforcement: false,
        EnableEmailInvitations: false,
        DisableBotsWhenOwnerIsDeactivated: true,
        EnableBotAccountCreation: false,
        EnableSVGs: false,
        EnableLatex: false,
        EnableInlineLatex: true,
        PostPriority: true,
        EnableAPIChannelDeletion: false,
        EnableLocalMode: false,
        LocalModeSocketLocation: '/var/tmp/mattermost_local.socket',
        EnableAWSMetering: false,
        SplitKey: '',
        FeatureFlagSyncIntervalSeconds: 30,
        DebugSplit: false,
        ThreadAutoFollow: true,
        CollapsedThreads: CollapsedThreads.ALWAYS_ON,
        ManagedResourcePaths: '',
        EnableCustomGroups: true,
        SelfHostedPurchase: true,
        AllowSyncedDrafts: true,
        AllowPersistentNotifications: true,
        PersistentNotificationMaxCount: 6,
        PersistentNotificationMaxRecipients: 5,
        PersistentNotificationIntervalMinutes: 5,
        AllowPersistentNotificationsForGuests: false,
    },
    TeamSettings: {
        SiteName: 'Mattermost',
        MaxUsersPerTeam: 50,
        EnableUserCreation: true,
        EnableOpenServer: false,
        EnableUserDeactivation: false,
        RestrictCreationToDomains: '',
        EnableCustomUserStatuses: true,
        EnableCustomBrand: false,
        CustomBrandText: '',
        CustomDescriptionText: '',
        RestrictDirectMessage: 'any',
        EnableLastActiveTime: true,
        UserStatusAwayTimeout: 300,
        MaxChannelsPerTeam: 2000,
        MaxNotificationsPerChannel: 1000,
        EnableConfirmNotificationsToChannel: true,
        TeammateNameDisplay: 'username',
        ExperimentalViewArchivedChannels: true,
        ExperimentalEnableAutomaticReplies: false,
        LockTeammateNameDisplay: false,
        ExperimentalPrimaryTeam: '',
        ExperimentalDefaultChannels: [],
    },
    ClientRequirements: {
        AndroidLatestVersion: '',
        AndroidMinVersion: '',
        IosLatestVersion: '',
        IosMinVersion: '',
    },
    SqlSettings: {
        DriverName: 'postgres',
        DataSource:
            'postgres://mmuser:mostest@localhost/mattermost_test?sslmode=disable\u0026connect_timeout=10\u0026binary_parameters=yes',
        DataSourceReplicas: [],
        DataSourceSearchReplicas: [],
        MaxIdleConns: 20,
        ConnMaxLifetimeMilliseconds: 3600000,
        ConnMaxIdleTimeMilliseconds: 300000,
        MaxOpenConns: 300,
        Trace: false,
        AtRestEncryptKey: '',
        QueryTimeout: 30,
        DisableDatabaseSearch: false,
        MigrationsStatementTimeoutSeconds: 100000,
        ReplicaLagSettings: [],
    },
    LogSettings: {
        EnableConsole: true,
        ConsoleLevel: 'DEBUG',
        ConsoleJson: true,
        EnableColor: false,
        EnableFile: true,
        FileLevel: 'INFO',
        FileJson: true,
        FileLocation: '',
        EnableWebhookDebugging: true,
        EnableDiagnostics: true,
        VerboseDiagnostics: false,
        EnableSentry: true,
        AdvancedLoggingConfig: '',
    },
    ExperimentalAuditSettings: {
        FileEnabled: false,
        FileName: '',
        FileMaxSizeMB: 100,
        FileMaxAgeDays: 0,
        FileMaxBackups: 0,
        FileCompress: false,
        FileMaxQueueSize: 1000,
        AdvancedLoggingConfig: '',
    },
    NotificationLogSettings: {
        EnableConsole: true,
        ConsoleLevel: 'DEBUG',
        ConsoleJson: true,
        EnableColor: false,
        EnableFile: true,
        FileLevel: 'INFO',
        FileJson: true,
        FileLocation: '',
        AdvancedLoggingConfig: '',
    },
    PasswordSettings: {
        MinimumLength: 8,
        Lowercase: false,
        Number: false,
        Uppercase: false,
        Symbol: false,
        EnableForgotLink: true,
    },
    FileSettings: {
        EnableFileAttachments: true,
        EnableMobileUpload: true,
        EnableMobileDownload: true,
        MaxFileSize: 104857600,
        MaxImageResolution: 33177600,
        MaxImageDecoderConcurrency: -1,
        DriverName: 'local',
        Directory: './data/',
        EnablePublicLink: false,
        ExtractContent: true,
        ArchiveRecursion: false,
        PublicLinkSalt: '',
        InitialFont: 'nunito-bold.ttf',
        AmazonS3AccessKeyId: '',
        AmazonS3SecretAccessKey: '',
        AmazonS3Bucket: '',
        AmazonS3PathPrefix: '',
        AmazonS3Region: '',
        AmazonS3Endpoint: 's3.amazonaws.com',
        AmazonS3SSL: true,
        AmazonS3SignV2: false,
        AmazonS3SSE: false,
        AmazonS3Trace: false,
        AmazonS3RequestTimeoutMilliseconds: 30000,
    },
    EmailSettings: {
        EnableSignUpWithEmail: true,
        EnableSignInWithEmail: true,
        EnableSignInWithUsername: true,
        SendEmailNotifications: true,
        UseChannelInEmailNotifications: false,
        RequireEmailVerification: false,
        FeedbackName: '',
        FeedbackEmail: 'test@example.com',
        ReplyToAddress: 'test@example.com',
        FeedbackOrganization: '',
        EnableSMTPAuth: false,
        SMTPUsername: '',
        SMTPPassword: '',
        SMTPServer: 'localhost',
        SMTPPort: '10025',
        SMTPServerTimeout: 10,
        ConnectionSecurity: '',
        SendPushNotifications: true,
        PushNotificationServer: 'https://push-test.mattermost.com',
        PushNotificationServerType: 'custom',
        PushNotificationServerLocation: 'us',
        PushNotificationContents: 'full',
        PushNotificationBuffer: 1000,
        EnableEmailBatching: false,
        EmailBatchingBufferSize: 256,
        EmailBatchingInterval: 30,
        EnablePreviewModeBanner: true,
        SkipServerCertificateVerification: false,
        EmailNotificationContentsType: 'full',
        LoginButtonColor: '#0000',
        LoginButtonBorderColor: '#2389D7',
        LoginButtonTextColor: '#2389D7',
    },
    RateLimitSettings: {
        Enable: false,
        PerSec: 10,
        MaxBurst: 100,
        MemoryStoreSize: 10000,
        VaryByRemoteAddr: true,
        VaryByUser: false,
        VaryByHeader: '',
    },
    PrivacySettings: {
        ShowEmailAddress: true,
        ShowFullName: true,
    },
    SupportSettings: {
        TermsOfServiceLink: 'https://mattermost.com/pl/terms-of-use/',
        PrivacyPolicyLink: 'https://mattermost.com/pl/privacy-policy/',
        AboutLink: 'https://docs.mattermost.com/pl/about-mattermost',
        HelpLink: 'https://mattermost.com/pl/help/',
        ReportAProblemLink: 'https://mattermost.com/pl/report-a-bug',
        ForgotPasswordLink: '',
        SupportEmail: '',
        CustomTermsOfServiceEnabled: false,
        CustomTermsOfServiceReAcceptancePeriod: 365,
        EnableAskCommunityLink: true,
    },
    AnnouncementSettings: {
        EnableBanner: false,
        BannerText: '',
        BannerColor: '#f2a93b',
        BannerTextColor: '#333333',
        AllowBannerDismissal: true,
        AdminNoticesEnabled: true,
        UserNoticesEnabled: true,
        NoticesURL: 'https://notices.mattermost.com/',
        NoticesFetchFrequency: 3600,
        NoticesSkipCache: false,
    },
    ThemeSettings: {
        EnableThemeSelection: true,
        DefaultTheme: 'default',
        AllowCustomThemes: true,
        AllowedThemes: [],
    },
    GitLabSettings: {
        Enable: false,
        Secret: '',
        Id: '',
        Scope: '',
        AuthEndpoint: '',
        TokenEndpoint: '',
        UserAPIEndpoint: '',
        DiscoveryEndpoint: '',
        ButtonText: '',
        ButtonColor: '',
    },
    GoogleSettings: {
        Enable: false,
        Secret: '',
        Id: '',
        Scope: 'profile email',
        AuthEndpoint: 'https://accounts.google.com/o/oauth2/v2/auth',
        TokenEndpoint: 'https://www.googleapis.com/oauth2/v4/token',
        UserAPIEndpoint:
            'https://people.googleapis.com/v1/people/me?personFields=names,emailAddresses,nicknames,metadata',
        DiscoveryEndpoint: '',
        ButtonText: '',
        ButtonColor: '',
    },
    Office365Settings: {
        Enable: false,
        Secret: '',
        Id: '',
        Scope: 'User.Read',
        AuthEndpoint: 'https://login.microsoftonline.com/common/oauth2/v2.0/authorize',
        TokenEndpoint: 'https://login.microsoftonline.com/common/oauth2/v2.0/token',
        UserAPIEndpoint: 'https://graph.microsoft.com/v1.0/me',
        DiscoveryEndpoint: '',
        DirectoryId: '',
    },
    OpenIdSettings: {
        Enable: false,
        Secret: '',
        Id: '',
        Scope: 'profile openid email',
        AuthEndpoint: '',
        TokenEndpoint: '',
        UserAPIEndpoint: '',
        DiscoveryEndpoint: '',
        ButtonText: '',
        ButtonColor: '#145DBF',
    },
    LdapSettings: {
        Enable: false,
        EnableSync: false,
        LdapServer: '',
        LdapPort: 389,
        ConnectionSecurity: '',
        BaseDN: '',
        BindUsername: '',
        BindPassword: '',
        UserFilter: '',
        GroupFilter: '',
        GuestFilter: '',
        EnableAdminFilter: false,
        AdminFilter: '',
        GroupDisplayNameAttribute: '',
        GroupIdAttribute: '',
        FirstNameAttribute: '',
        LastNameAttribute: '',
        EmailAttribute: '',
        UsernameAttribute: '',
        NicknameAttribute: '',
        IdAttribute: '',
        PositionAttribute: '',
        LoginIdAttribute: '',
        PictureAttribute: '',
        SyncIntervalMinutes: 60,
        SkipCertificateVerification: false,
        PublicCertificateFile: '',
        PrivateKeyFile: '',
        QueryTimeout: 60,
        MaxPageSize: 0,
        LoginFieldName: '',
        LoginButtonColor: '#0000',
        LoginButtonBorderColor: '#2389D7',
        LoginButtonTextColor: '#2389D7',
        Trace: false,
    },
    ComplianceSettings: {
        Enable: false,
        Directory: './data/',
        EnableDaily: false,
        BatchSize: 30000,
    },
    LocalizationSettings: {
        DefaultServerLocale: 'en',
        DefaultClientLocale: 'en',
        AvailableLocales: '',
    },
    SamlSettings: {
        Enable: false,
        EnableSyncWithLdap: false,
        EnableSyncWithLdapIncludeAuth: false,
        IgnoreGuestsLdapSync: false,
        Verify: true,
        Encrypt: true,
        SignRequest: false,
        IdpURL: '',
        IdpDescriptorURL: '',
        IdpMetadataURL: '',
        ServiceProviderIdentifier: '',
        AssertionConsumerServiceURL: '',
        SignatureAlgorithm: 'RSAwithSHA1',
        CanonicalAlgorithm: 'Canonical1.0',
        ScopingIDPProviderId: '',
        ScopingIDPName: '',
        IdpCertificateFile: '',
        PublicCertificateFile: '',
        PrivateKeyFile: '',
        IdAttribute: '',
        GuestAttribute: '',
        EnableAdminAttribute: false,
        AdminAttribute: '',
        FirstNameAttribute: '',
        LastNameAttribute: '',
        EmailAttribute: '',
        UsernameAttribute: '',
        NicknameAttribute: '',
        LocaleAttribute: '',
        PositionAttribute: '',
        LoginButtonText: 'SAML',
        LoginButtonColor: '#34a28b',
        LoginButtonBorderColor: '#2389D7',
        LoginButtonTextColor: '#ffffff',
    },
    NativeAppSettings: {
        AppCustomURLSchemes: ['mmauth://', 'mmauthbeta://'],
        AppDownloadLink: 'https://mattermost.com/pl/download-apps',
        AndroidAppDownloadLink: 'https://mattermost.com/pl/android-app/',
        IosAppDownloadLink: 'https://mattermost.com/pl/ios-app/',
    },
    ClusterSettings: {
        Enable: false,
        ClusterName: '',
        OverrideHostname: '',
        NetworkInterface: '',
        BindAddress: '',
        AdvertiseAddress: '',
        UseIPAddress: true,
        EnableGossipCompression: true,
        EnableExperimentalGossipEncryption: false,
        ReadOnlyConfig: true,
        GossipPort: 8074,
        StreamingPort: 8075,
        MaxIdleConns: 100,
        MaxIdleConnsPerHost: 128,
        IdleConnTimeoutMilliseconds: 90000,
    },
    MetricsSettings: {
        Enable: false,
        BlockProfileRate: 0,
        ListenAddress: ':8067',
    },
    ExperimentalSettings: {
        ClientSideCertEnable: false,
        ClientSideCertCheck: 'secondary',
        LinkMetadataTimeoutMilliseconds: 5000,
        RestrictSystemAdmin: false,
        UseNewSAMLLibrary: false,
        EnableSharedChannels: false,
        EnableRemoteClusterService: false,
        DisableAppBar: true,
        DisableRefetchingOnBrowserFocus: false,
        DelayChannelAutocomplete: false,
    },
    AnalyticsSettings: {
        MaxUsersForStatistics: 2500,
    },
    ElasticsearchSettings: {
        ConnectionURL: 'http://localhost:9200',
        Username: 'elastic',
        Password: 'changeme',
        EnableIndexing: false,
        EnableSearching: false,
        EnableAutocomplete: false,
        Sniff: true,
        PostIndexReplicas: 1,
        PostIndexShards: 1,
        ChannelIndexReplicas: 1,
        ChannelIndexShards: 1,
        UserIndexReplicas: 1,
        UserIndexShards: 1,
        AggregatePostsAfterDays: 365,
        PostsAggregatorJobStartTime: '03:00',
        IndexPrefix: '',
        LiveIndexingBatchSize: 1,
        BatchSize: 10000,
        RequestTimeoutSeconds: 30,
        SkipTLSVerification: false,
        CA: '',
        ClientCert: '',
        ClientKey: '',
        Trace: '',
        IgnoredPurgeIndexes: '',
    },
    BleveSettings: {
        IndexDir: '',
        EnableIndexing: false,
        EnableSearching: false,
        EnableAutocomplete: false,
        BatchSize: 10000,
    },
    DataRetentionSettings: {
        EnableMessageDeletion: false,
        EnableFileDeletion: false,
        MessageRetentionDays: 365,
        FileRetentionDays: 365,
        DeletionJobStartTime: '02:00',
        BatchSize: 3000,
    },
    MessageExportSettings: {
        EnableExport: false,
        ExportFormat: 'actiance',
        DailyRunTime: '01:00',
        ExportFromTimestamp: 0,
        BatchSize: 10000,
        DownloadExportResults: false,
        GlobalRelaySettings: {
            CustomerType: 'A9',
            SMTPUsername: '',
            SMTPPassword: '',
            EmailAddress: '',
            SMTPServerTimeout: 1800,
        },
    },
    JobSettings: {
        RunJobs: true,
        RunScheduler: true,
        CleanupJobsThresholdDays: -1,
        CleanupConfigThresholdDays: -1,
    },
    ProductSettings: {},
    PluginSettings: {
        Enable: true,
        EnableUploads: false,
        AllowInsecureDownloadURL: false,
        EnableHealthCheck: true,
        Directory: './plugins',
        ClientDirectory: './client/plugins',
        Plugins: {},
        PluginStates: {
            'com.mattermost.calls': {
                Enable: true,
            },
            'com.mattermost.nps': {
                Enable: true,
            },
        },
        EnableMarketplace: true,
        EnableRemoteMarketplace: true,
        AutomaticPrepackagedPlugins: true,
        RequirePluginSignature: false,
        MarketplaceURL: 'https://api.integrations.mattermost.com',
        SignaturePublicKeyFiles: [],
        ChimeraOAuthProxyURL: '',
    },
    DisplaySettings: {
        CustomURLSchemes: [],
        ExperimentalTimezone: true,
    },
    GuestAccountsSettings: {
        Enable: false,
        HideTags: false,
        AllowEmailAccounts: true,
        EnforceMultifactorAuthentication: false,
        RestrictCreationToDomains: '',
    },
    ImageProxySettings: {
        Enable: false,
        ImageProxyType: 'local',
        RemoteImageProxyURL: '',
        RemoteImageProxyOptions: '',
    },
    CloudSettings: {
        CWSURL: 'https://customers.mattermost.com',
        CWSAPIURL: 'https://portal.internal.prod.cloud.mattermost.com',
    },
    FeatureFlags: {
        TestFeature: 'off',
        TestBoolFeature: false,
        EnableRemoteClusterService: false,
        AppsEnabled: true,
        PluginPlaybooks: '',
        PluginApps: '',
        PluginFocalboard: '',
        PluginCalls: '',
        PermalinkPreviews: true,
        CallsEnabled: true,
        BoardsFeatureFlags: '',
        NormalizeLdapDNs: false,
        GraphQL: false,
        CommandPalette: false,
        SendWelcomePost: true,
        PostPriority: true,
        WysiwygEditor: false,
        ThreadsEverywhere: false,
        OnboardingTourTips: true,
        DeprecateCloudFree: false,
        CloudReverseTrial: false,
        StreamlinedMarketplace: true
    },
    ImportSettings: {
        Directory: './import',
        RetentionDays: 30,
    },
    ExportSettings: {
        Directory: './export',
        RetentionDays: 30,
    },
    WranglerSettings: {
        PermittedWranglerRoles: [],
        AllowedEmailDomain: [],
        MoveThreadMaxCount: 30,
        MoveThreadToAnotherTeamEnable: true,
        MoveThreadFromPrivateChannelEnable: true,
        MoveThreadFromDirectMessageChannelEnable: true,
        MoveThreadFromGroupMessageChannelEnable: true,
    }
};
