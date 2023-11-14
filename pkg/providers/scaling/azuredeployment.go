// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.
package scaling

import (
	"context"
	"fmt"
	"math"
	"time"

	prototypes "k9s-autoscaler/pkg/proto"
	"k9s-autoscaler/pkg/providers"
	"k9s-autoscaler/pkg/providers/scaling/proto"
	scalingtypes "k9s-autoscaler/pkg/scale/types"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	azcorepolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
	protob "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s.io/klog/v2"
)

type azureDeployment struct {
	creds *azidentity.DefaultAzureCredential
}

type azureDeploymentFactory struct{}

func init() {
	providers.RegisterScalingClient(&proto.AzureDeploymentConfig{}, &proto.AzureDeploymentTargetConfig{}, &azureDeploymentFactory{})
}

func newAzureDeployment(config *proto.AzureDeploymentConfig) (*azureDeployment, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get azure default credentials: %v", err)
	}

	return &azureDeployment{
		creds: cred,
	}, nil
}

func (ad *azureDeployment) SetScaleTarget(ctx context.Context, name, namespace string, scaleTarget *prototypes.AutoscalerTarget, target *prototypes.ScaleSpec) error {
	targetConfig, err := ad.getScaleTargetConfig(scaleTarget)
	if err != nil {
		return err
	}

	resourceID, client, err := ad.getDeploymentClient(targetConfig.ResourceURI)
	if err != nil {
		return err
	}
	resp, err := client.Get(
		ctx,
		resourceID.ResourceGroupName,
		resourceID.Name,
		targetConfig.DeploymentName,
		&armcognitiveservices.DeploymentsClientGetOptions{})
	if err != nil {
		return fmt.Errorf("failed get operation: %v", err)
	}

	targetScale := target.Desired
	if targetConfig.ScaleDenominator != nil {
		if targetScale < *resp.SKU.Capacity {
			targetScale = *targetConfig.ScaleDenominator * int32(math.Floor(float64(targetScale)/float64(*targetConfig.ScaleDenominator)))
		} else {
			targetScale = *targetConfig.ScaleDenominator * int32(math.Ceil(float64(targetScale)/float64(*targetConfig.ScaleDenominator)))
		}
		klog.InfoS("adjusting scale to denominator", "denominator", *targetConfig.ScaleDenominator, "desired", target.Desired, "adjusted", targetScale)
	}

	update := armcognitiveservices.Deployment{
		Properties: resp.Properties,
		SKU:        resp.SKU,
	}
	update.SKU.Capacity = to.Ptr[int32](targetScale)
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		resourceID.ResourceGroupName,
		resourceID.Name,
		targetConfig.DeploymentName,
		update,
		nil)
	if err != nil {
		return fmt.Errorf("failed to initiate scale update request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	if err != nil {
		return fmt.Errorf("failed update operation: %v", err)
	}

	return nil
}

func (ad *azureDeployment) GetScale(ctx context.Context, name, namespace string, scaleTarget *prototypes.AutoscalerTarget) (*prototypes.Scale, error) {
	targetConfig, err := ad.getScaleTargetConfig(scaleTarget)
	if err != nil {
		return nil, err
	}

	resourceID, client, err := ad.getDeploymentClient(targetConfig.ResourceURI)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get(
		ctx,
		resourceID.ResourceGroupName,
		resourceID.Name,
		targetConfig.DeploymentName,
		&armcognitiveservices.DeploymentsClientGetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed get operation: %v", err)
	}

	return &prototypes.Scale{
		Spec: &prototypes.ScaleSpec{
			Desired: *resp.SKU.Capacity,
		},
		Status: &prototypes.ScaleStatus{
			Current: *resp.SKU.Capacity,
		},
	}, nil
}

func (ad *azureDeployment) getScaleTargetConfig(scaleTarget *prototypes.AutoscalerTarget) (*proto.AzureDeploymentTargetConfig, error) {
	config := proto.AzureDeploymentTargetConfig{}
	if err := anypb.UnmarshalTo(scaleTarget.Config, &config, protob.UnmarshalOptions{}); err != nil {
		return nil, err
	}

	return &config, nil
}

func (ad *azureDeployment) getDeploymentClient(resourceURI string) (*arm.ResourceID, *armcognitiveservices.DeploymentsClient, error) {
	resourceID, err := arm.ParseResourceID(resourceURI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse resource ID %s: %v", resourceURI, err)
	}
	clientFactory, err := armcognitiveservices.NewClientFactory(
		resourceID.SubscriptionID,
		ad.creds,
		&armpolicy.ClientOptions{
			ClientOptions: azcorepolicy.ClientOptions{
				APIVersion: "2023-10-01-preview",
			},
		})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client factory: %v", err)
	}

	return resourceID, clientFactory.NewDeploymentsClient(), nil
}

func (f *azureDeploymentFactory) ScalingClient(config *anypb.Any) (scalingtypes.ScalingClient, error) {
	adConfig := proto.AzureDeploymentConfig{}
	if err := anypb.UnmarshalTo(config, &adConfig, protob.UnmarshalOptions{}); err != nil {
		return nil, err
	}

	return newAzureDeployment(&adConfig)
}
