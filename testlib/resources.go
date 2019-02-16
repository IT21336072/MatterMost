// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package testlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/utils"
	"github.com/mattermost/mattermost-server/utils/fileutils"
)

const (
	resourceTypeFile = iota
	resourceTypeFolder
)

const (
	actionCopy = iota
	actionSymlink
)

type testResourceDetails struct {
	src     string
	dest    string
	resType int8
	action  int8
}

var testResourcesToSetup = []testResourceDetails{
	{"config/timezones.json", "config/timezones.json", resourceTypeFile, actionCopy},
	{"i18n", "i18n", resourceTypeFolder, actionSymlink},
	{"client", "client", resourceTypeFolder, actionSymlink},
	{"templates", "templates", resourceTypeFolder, actionSymlink},
	{"tests", "tests", resourceTypeFolder, actionSymlink},
	{"utils/policies-roles-mapping.json", "utils/policies-roles-mapping.json", resourceTypeFile, actionSymlink},
}

func resolveResourcesPath() {
	var srcPath string
	var found bool

	// Finding resources and setting full path to source to be used for further processing
	for i, testResource := range testResourcesToSetup {
		if testResource.resType == resourceTypeFile {
			srcPath = fileutils.FindFile(testResource.src)
			if srcPath == "" {
				panic(fmt.Sprintf("Failed to find file %s", testResource.src))
			}

			testResourcesToSetup[i].src = srcPath
		} else if testResource.resType == resourceTypeFolder {
			srcPath, found = fileutils.FindDir(testResource.src)
			if found == false {
				panic(fmt.Sprintf("Failed to find folder %s", testResource.src))
			}

			testResourcesToSetup[i].src = srcPath
		} else {
			panic(fmt.Sprintf("Invalid resource type: %d", testResource.resType))
		}
	}
}

func SetupTestResources() (string, error) {
	resolveResourcesPath()

	tempDir, err := ioutil.TempDir("", "testlib")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temporary directory")
	}

	err = setupConfig(path.Join(tempDir, "config"))
	if err != nil {
		return "", errors.Wrap(err, "failed to setup config")
	}

	var resourceDestInTemp string

	// Setting up test resources in temp.
	// Action in each resource tells whether it needs to be copied or just symlinked

	for _, testResource := range testResourcesToSetup {
		resourceDestInTemp = filepath.Join(tempDir, testResource.dest)

		if testResource.action == actionCopy {
			if testResource.resType == resourceTypeFile {
				err = utils.CopyFile(testResource.src, resourceDestInTemp)
				if err != nil {
					return "", errors.Wrapf(err, "failed to copy file %s to %s", testResource.src, resourceDestInTemp)
				}
			} else if testResource.resType == resourceTypeFolder {
				err = utils.CopyDir(testResource.src, resourceDestInTemp)
				if err != nil {
					return "", errors.Wrapf(err, "failed to copy folder %s to %s", testResource.src, resourceDestInTemp)
				}
			}
		} else if testResource.action == actionSymlink {
			destDir := path.Dir(resourceDestInTemp)
			if destDir != "." {
				err = os.MkdirAll(destDir, os.ModePerm)
				if err != nil {
					return "", errors.Wrapf(err, "failed to make dir %s", destDir)
				}
			}

			err = os.Symlink(testResource.src, resourceDestInTemp)
			if err != nil {
				return "", errors.Wrapf(err, "failed to symlink %s to %s", testResource.src, resourceDestInTemp)
			}
		} else {
			return "", errors.Wrapf(err, "Invalid action: %d", testResource.action)
		}

	}

	return tempDir, nil
}

func setupConfig(configDir string) error {
	var err error
	var config model.Config

	config.SetDefaults()

	err = os.Mkdir(configDir, 0700)
	if err != nil {
		return errors.Wrapf(err, "failed to create config directory %s", configDir)
	}

	defaultJson := path.Join(configDir, "default.json")
	err = ioutil.WriteFile(defaultJson, []byte(config.ToJson()), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write config to %s", defaultJson)
	}

	configJson := path.Join(configDir, "config.json")
	err = utils.CopyFile(defaultJson, configJson)
	if err != nil {
		return errors.Wrapf(err, "failed to copy file %s to %s", defaultJson, configJson)
	}

	return nil
}
