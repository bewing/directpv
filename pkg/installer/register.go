// This file is part of MinIO DirectPV
// Copyright (c) 2021, 2022 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package installer

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/fatih/color"
	directpvtypes "github.com/minio/directpv/pkg/apis/directpv.min.io/types"
	"github.com/minio/directpv/pkg/consts"
	"github.com/minio/directpv/pkg/k8s"
	"k8s.io/apiextensions-apiserver/pkg/apihelpers"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

//go:embed directpv.min.io_directpvdrives.yaml
var drivesYAML []byte

//go:embed directpv.min.io_directpvvolumes.yaml
var volumesYAML []byte

func setNoneConversionStrategy(crd *apiextensions.CustomResourceDefinition) {
	crd.Spec.Conversion = &apiextensions.CustomResourceConversion{
		Strategy: apiextensions.NoneConverter,
	}
}

func updateLabels(object metav1.Object, labels map[directpvtypes.LabelKey]directpvtypes.LabelValue) {
	values := object.GetLabels()
	if values == nil {
		values = make(map[string]string)
	}

	for key, value := range labels {
		values[string(key)] = string(value)
	}

	object.SetLabels(values)
}

func getLatestCRDVersionObject(newCRD *apiextensions.CustomResourceDefinition) (crdVersion apiextensions.CustomResourceDefinitionVersion, err error) {
	for i := range newCRD.Spec.Versions {
		if newCRD.Spec.Versions[i].Name == consts.LatestAPIVersion {
			return newCRD.Spec.Versions[i], nil
		}
	}

	return crdVersion, fmt.Errorf("no version %v found crd %v", consts.LatestAPIVersion, newCRD.Name)
}

func syncCRD(ctx context.Context, existingCRD, newCRD *apiextensions.CustomResourceDefinition, c *Config) error {
	existingCRDStorageVersion, err := apihelpers.GetCRDStorageVersion(existingCRD)
	if err != nil {
		return err
	}

	var versionEntryFound bool
	if existingCRDStorageVersion != consts.LatestAPIVersion {
		// Set all the existing versions to false
		for i := range existingCRD.Spec.Versions {
			if existingCRD.Spec.Versions[i].Name == consts.LatestAPIVersion {
				existingCRD.Spec.Versions[i].Storage = true
				versionEntryFound = true
			} else {
				existingCRD.Spec.Versions[i].Storage = false
			}
		}

		if !versionEntryFound {
			latestVersionObject, err := getLatestCRDVersionObject(newCRD)
			if err != nil {
				return err
			}
			existingCRD.Spec.Versions = append(existingCRD.Spec.Versions, latestVersionObject)
		}
	}

	setNoneConversionStrategy(existingCRD)

	if c.DryRun {
		updateLabels(existingCRD, map[directpvtypes.LabelKey]directpvtypes.LabelValue{directpvtypes.VersionLabelKey: consts.LatestAPIVersion})
		existingCRD.TypeMeta = newCRD.TypeMeta
	} else {
		if _, err := k8s.CRDClient().Update(ctx, existingCRD, metav1.UpdateOptions{}); err != nil {
			return err
		}

		fmt.Fprintln(os.Stderr, color.HiYellowString("updated CRD %v to %v", existingCRD.Name, consts.LatestAPIVersion))
	}

	return c.postProc(existingCRD)
}

func registerCRDs(ctx context.Context, c *Config) error {
	register := func(data []byte) error {
		object := map[string]interface{}{}
		if err := yaml.Unmarshal(data, &object); err != nil {
			return err
		}

		var crd apiextensions.CustomResourceDefinition
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object, &crd); err != nil {
			return err
		}

		existingCRD, err := k8s.CRDClient().Get(ctx, crd.Name, metav1.GetOptions{})
		if err != nil {
			if !apierrors.IsNotFound(err) {
				return err
			}

			setNoneConversionStrategy(&crd)

			if c.DryRun {
				updateLabels(&crd, map[directpvtypes.LabelKey]directpvtypes.LabelValue{directpvtypes.VersionLabelKey: consts.LatestAPIVersion})
			} else if _, err = k8s.CRDClient().Create(ctx, &crd, metav1.CreateOptions{}); err != nil {
				return err
			}

			return c.postProc(crd)
		}

		return syncCRD(ctx, existingCRD, &crd, c)
	}

	if err := register(drivesYAML); err != nil {
		return err
	}

	return register(volumesYAML)
}

func unregisterCRDs(ctx context.Context) error {
	if err := k8s.CRDClient().Delete(ctx, consts.DriveResource+"."+consts.GroupName, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	if err := k8s.CRDClient().Delete(ctx, consts.VolumeResource+"."+consts.GroupName, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}
