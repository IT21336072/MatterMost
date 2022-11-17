// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v6/einterfaces/mocks"
	"github.com/mattermost/mattermost-server/v6/model"
)

func TestSelfHostedBootstrap(t *testing.T) {
	t.Run("feature flag off returns not implemented", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		th.Client.Login(th.SystemAdminUser.Email, th.SystemAdminUser.Password)

		os.Setenv("MM_FEATUREFLAGS_SELFHOSTEDFIRSTTIMEPURCHASE", "false")
		defer os.Unsetenv("MM_FEATUREFLAGS_SELFHOSTEDFIRSTTIMEPURCHASE")
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.FeatureFlags.SelfHostedFirstTimePurchase = false })
		th.App.ReloadConfig()

		_, r, err := th.Client.BootstrapSelfHostedSignup(model.BootstrapSelfHostedSignupRequest{Email: th.SystemAdminUser.Email})

		require.Equal(t, http.StatusNotImplemented, r.StatusCode)
		require.Error(t, err)
	})

	t.Run("cloud instances not allowed to bootstrap self-hosted signup", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		th.Client.Login(th.SystemAdminUser.Email, th.SystemAdminUser.Password)

		th.App.Srv().SetLicense(model.NewTestLicense("cloud"))
		os.Setenv("MM_FEATUREFLAGS_SELFHOSTEDFIRSTTIMEPURCHASE", "true")
		defer os.Unsetenv("MM_FEATUREFLAGS_SELFHOSTEDFIRSTTIMEPURCHASE")
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.FeatureFlags.SelfHostedFirstTimePurchase = true })
		th.App.ReloadConfig()

		_, r, err := th.Client.BootstrapSelfHostedSignup(model.BootstrapSelfHostedSignupRequest{Email: th.SystemAdminUser.Email})

		require.Equal(t, http.StatusBadRequest, r.StatusCode)
		require.Error(t, err)
	})

	t.Run("non-admins not allowed to bootstrap self-hosted signup", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		th.Client.Login(th.BasicUser.Email, th.BasicUser.Password)

		os.Setenv("MM_FEATUREFLAGS_SELFHOSTEDFIRSTTIMEPURCHASE", "true")
		defer os.Unsetenv("MM_FEATUREFLAGS_SELFHOSTEDFIRSTTIMEPURCHASE")
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.FeatureFlags.SelfHostedFirstTimePurchase = true })
		th.App.ReloadConfig()

		_, r, err := th.Client.BootstrapSelfHostedSignup(model.BootstrapSelfHostedSignupRequest{Email: th.SystemAdminUser.Email})

		require.Equal(t, http.StatusForbidden, r.StatusCode)
		require.Error(t, err)
	})

	t.Run("self-hosted admins can bootstrap self-hosted signup", func(t *testing.T) {
		th := Setup(t).InitBasic()
		defer th.TearDown()

		th.Client.Login(th.SystemAdminUser.Email, th.SystemAdminUser.Password)

		os.Setenv("MM_FEATUREFLAGS_SELFHOSTEDFIRSTTIMEPURCHASE", "true")
		defer os.Unsetenv("MM_FEATUREFLAGS_SELFHOSTEDFIRSTTIMEPURCHASE")
		th.App.UpdateConfig(func(cfg *model.Config) { cfg.FeatureFlags.SelfHostedFirstTimePurchase = true })
		th.App.ReloadConfig()
		cloud := mocks.CloudInterface{}

		cloud.Mock.On("BootstrapSelfHostedSignup", mock.Anything).Return(&model.BootstrapSelfHostedSignupResponse{Progress: "START"}, nil)

		cloudImpl := th.App.Srv().Cloud
		defer func() {
			th.App.Srv().Cloud = cloudImpl
		}()
		th.App.Srv().Cloud = &cloud

		response, r, err := th.Client.BootstrapSelfHostedSignup(model.BootstrapSelfHostedSignupRequest{Email: th.SystemAdminUser.Email})

		require.Equal(t, http.StatusOK, r.StatusCode)
		require.NoError(t, err)
		require.Equal(t, "START", response.Progress)
	})
}
