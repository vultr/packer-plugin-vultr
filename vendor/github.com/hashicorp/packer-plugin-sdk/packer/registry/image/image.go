// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

/* Package image allows for the management of image metadata that can be stored in a HCP Packer registry.
 */
package image

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/packer-plugin-sdk/packer"
)

// ArtifactStateURI represents the key used by Packer when querying a packersdk.Artifact
// for Image metadata that a particular component would like to have stored on the HCP Packer Registry.
const ArtifactStateURI = "par.artifact.metadata"

// ArtifactOverrideFunc represents a transformation func that can be applied to a
// non-nil *Image. See WithID, WithRegion functions for examples.
type ArtifactOverrideFunc func(*Image) error

// Image represents the metadata for some Artifact in the HCP Packer Registry.
type Image struct {
	// ImageID is a unique reference identifier stored on the HCP Packer registry
	// that can be used to get back the built artifact of a builder or post-processor.
	ImageID string
	// ProviderName represents the name of the top level cloud or service where the built artifact resides.
	// For example "aws, azure, docker, gcp, and vsphere".
	ProviderName string
	// ProviderRegion represents the location of the built artifact.
	// For cloud providers region usually maps to a cloud region or zone, but for things like the file builder,
	// S3 bucket or vsphere cluster region can represent a path on the upstream datastore, or cluster.
	ProviderRegion string
	// Labels represents additional details about an image that a builder or post-processor may with to provide for a given build.
	// Any additional metadata will be made available as build labels within a HCP Packer registry iteration.
	Labels map[string]string
	// SourceImageID is the cloud image id of the image that was used as the
	// source for this image. If set, the HCP Packer registry will be able
	// link the parent and child images for ancestry visualizations and
	// dependency tracking.
	SourceImageID string
}

// Validate checks that the Image i contains a non-empty ImageID and ProviderName.
func (i *Image) Validate() error {
	if i.ImageID == "" {
		return errors.New("error registry image does not contain a valid ImageId")
	}

	if i.ProviderName == "" {
		return errors.New("error registry image does not contain a valid ProviderName")
	}

	return nil
}

// String returns string representation of Image
func (i *Image) String() string {
	return fmt.Sprintf("provider:%s, image:%s, region:%s", i.ProviderName, i.ImageID, i.ProviderRegion)
}

// FromMappedData calls f sequentially for each key and value present in mappedData to create a []*Image
// as the final return value. If there is an error in calling f, FromMappedData will stop processing mapped items
// and result in a nil slice, with the said error.
//
// FromMappedData will make its best attempt to convert the input map into map[interface{}]interface{} before
// calling f(k,v). The func f is responsible for type asserting the expected type for the key and value before
// trying to create an Image from it.
func FromMappedData(mappedData interface{}, f func(key, value interface{}) (*Image, error)) ([]*Image, error) {
	mapValue := reflect.ValueOf(mappedData)
	if mapValue.Kind() != reflect.Map {
		return nil, errors.New("error the incoming mappedData does not appear to be a map; found type to be" + mapValue.Kind().String())
	}

	keys := mapValue.MapKeys()
	var images []*Image
	for _, k := range keys {
		v := mapValue.MapIndex(k)
		i, err := f(k.Interface(), v.Interface())
		if err != nil {
			return nil, err
		}
		images = append(images, i)
	}
	return images, nil
}

// FromArtifact returns an *Image that can be used by Packer core for publishing to the HCP Packer Registry.
// By default FromArtifact will use the a.BuilderID() as the ProviderName, and the a.Id() as the ImageID that
// should be tracked within the HCP Packer Registry. No Region is selected by default as region varies per builder.
// The use of one or more ArtifactOverrideFunc can be used to override any of the defaults used.
func FromArtifact(a packer.Artifact, opts ...ArtifactOverrideFunc) (*Image, error) {
	if a == nil {
		return nil, errors.New("unable to create Image from nil artifact")
	}

	img := Image{
		ProviderName: a.BuilderId(),
		ImageID:      a.Id(),
		Labels:       make(map[string]string),
	}

	for _, opt := range opts {
		err := opt(&img)
		if err != nil {
			return nil, err
		}
	}

	return &img, nil
}

// WithProvider takes a name, and returns a ArtifactOverrideFunc that can be
// used to override the ProviderName for an existing Image.
func WithProvider(name string) func(*Image) error {
	return func(img *Image) error {
		img.ProviderName = name
		return nil
	}
}

// WithID takes a id, and returns a ArtifactOverrideFunc that can be
// used to override the ImageId for an existing Image.
func WithID(id string) func(*Image) error {
	return func(img *Image) error {
		img.ImageID = id
		return nil
	}
}

// WithSourceID takes a id, and returns a ArtifactOverrideFunc that can be
// used to override the SourceImageId for an existing Image.
func WithSourceID(id string) func(*Image) error {
	return func(img *Image) error {
		img.SourceImageID = id
		return nil
	}
}

// WithRegion takes a region, and returns a ArtifactOverrideFunc that can be
// used to override the ProviderRegion for an existing Image.
func WithRegion(region string) func(*Image) error {
	return func(img *Image) error {
		img.ProviderRegion = region
		return nil
	}
}

// SetLabels takes metadata, and returns a ArtifactOverrideFunc that can be
// used to set metadata for an existing Image. The incoming metadata `md`
// will be filtered only for keys whose values are of type string.
// If you wish to override this behavior you may create your own  ArtifactOverrideFunc
// for manipulating and setting Image metadata.
func SetLabels(md map[string]interface{}) func(*Image) error {
	return func(img *Image) error {
		if img.Labels == nil {
			img.Labels = make(map[string]string)
		}

		for k, v := range md {
			v, ok := v.(string)
			if !ok {
				continue
			}
			img.Labels[k] = v
		}

		return nil
	}
}
