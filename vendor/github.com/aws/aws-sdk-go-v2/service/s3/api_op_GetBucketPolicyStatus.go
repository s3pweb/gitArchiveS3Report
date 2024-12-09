// Code generated by smithy-go-codegen DO NOT EDIT.

package s3

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/aws/aws-sdk-go-v2/service/s3/internal/customizations"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Retrieves the policy status for an Amazon S3 bucket, indicating whether the
// bucket is public. In order to use this operation, you must have the
// s3:GetBucketPolicyStatus permission. For more information about Amazon S3
// permissions, see Specifying Permissions in a Policy
// (https://docs.aws.amazon.com/AmazonS3/latest/dev/using-with-s3-actions.html).
// For more information about when Amazon S3 considers a bucket public, see The
// Meaning of "Public"
// (https://docs.aws.amazon.com/AmazonS3/latest/dev/access-control-block-public-access.html#access-control-block-public-access-policy-status).
// The following operations are related to GetBucketPolicyStatus:
//
// * Using Amazon
// S3 Block Public Access
// (https://docs.aws.amazon.com/AmazonS3/latest/dev/access-control-block-public-access.html)
//
// *
// GetPublicAccessBlock
// (https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetPublicAccessBlock.html)
//
// *
// PutPublicAccessBlock
// (https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutPublicAccessBlock.html)
//
// *
// DeletePublicAccessBlock
// (https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeletePublicAccessBlock.html)
func (c *Client) GetBucketPolicyStatus(ctx context.Context, params *GetBucketPolicyStatusInput, optFns ...func(*Options)) (*GetBucketPolicyStatusOutput, error) {
	if params == nil {
		params = &GetBucketPolicyStatusInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetBucketPolicyStatus", params, optFns, addOperationGetBucketPolicyStatusMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetBucketPolicyStatusOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetBucketPolicyStatusInput struct {

	// The name of the Amazon S3 bucket whose policy status you want to retrieve.
	//
	// This member is required.
	Bucket *string

	// The account id of the expected bucket owner. If the bucket is owned by a
	// different account, the request will fail with an HTTP 403 (Access Denied) error.
	ExpectedBucketOwner *string
}

type GetBucketPolicyStatusOutput struct {

	// The policy status for the specified bucket.
	PolicyStatus *types.PolicyStatus

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func addOperationGetBucketPolicyStatusMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsRestxml_serializeOpGetBucketPolicyStatus{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpGetBucketPolicyStatus{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpGetBucketPolicyStatusValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetBucketPolicyStatus(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addGetBucketPolicyStatusUpdateEndpoint(stack, options); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = v4.AddContentSHA256HeaderMiddleware(stack); err != nil {
		return err
	}
	if err = disableAcceptEncodingGzip(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opGetBucketPolicyStatus(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "s3",
		OperationName: "GetBucketPolicyStatus",
	}
}

// getGetBucketPolicyStatusBucketMember returns a pointer to string denoting a
// provided bucket member valueand a boolean indicating if the input has a modeled
// bucket name,
func getGetBucketPolicyStatusBucketMember(input interface{}) (*string, bool) {
	in := input.(*GetBucketPolicyStatusInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}
func addGetBucketPolicyStatusUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getGetBucketPolicyStatusBucketMember,
		},
		UsePathStyle:            options.UsePathStyle,
		UseAccelerate:           options.UseAccelerate,
		SupportsAccelerate:      true,
		EndpointResolver:        options.EndpointResolver,
		EndpointResolverOptions: options.EndpointOptions,
		UseDualstack:            options.UseDualstack,
		UseARNRegion:            options.UseARNRegion,
	})
}
