package main

import (
	"context"
	"fmt"

	"github.com/eko/authz/backend/pkg/authz"
	"github.com/eko/authz/sdk"
	"github.com/eko/authz/sdk/rule"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// Initialize client
	authzClient, err := sdk.NewClient(&sdk.Config{
		ClientID:     "47bb0c4c-9115-11ed-8f3d-acde48001122",
		ClientSecret: "kJVvoHPZQuD7EiwSBbpRqBbzNJHw_WLhj4HRkcC9EBXohavB",
		GrpcAddr:     "127.0.0.1:8081",
	})
	if err != nil {
		panic(fmt.Sprintf("cannot initialize client: %v", err))
	}

	ctx := context.Background()

	// Create a principal named "user-123"
	principalGetResponse, err := authzClient.PrincipalGet(ctx, &authz.PrincipalGetRequest{
		Id: "user-123",
	})
	if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
		principalResponse, err := authzClient.PrincipalCreate(ctx, &authz.PrincipalCreateRequest{
			Id: "user-123",
			Attributes: []*authz.Attribute{
				{Key: "email", Value: "johndoe@acme.tld"},
			},
		})
		if err != nil {
			panic(fmt.Sprintf("cannot create principal: %v", err))
		}

		fmt.Printf("Principal created: %s\n", principalResponse.GetPrincipal().GetId())
	} else {
		fmt.Printf("Principal already exists: %s\n", principalGetResponse.GetPrincipal().GetId())
	}

	// Create a first post resource not matching attribute value
	resourceGetResponse, err := authzClient.ResourceGet(ctx, &authz.ResourceGetRequest{
		Id: "post.123",
	})
	if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
		resourceResponse, err := authzClient.ResourceCreate(ctx, &authz.ResourceCreateRequest{
			Id:    "post.123",
			Kind:  "post",
			Value: "123",
			Attributes: []*authz.Attribute{
				{Key: "owner_email", Value: "someoneelse@acme.tld"},
			},
		})
		if err != nil {
			panic(fmt.Sprintf("cannot create resource: %v", err))
		}

		fmt.Printf("Resource created: %s\n", resourceResponse.GetResource().GetId())
	} else {
		fmt.Printf("Resource already exists: %s\n", resourceGetResponse.GetResource().GetId())
	}

	// Create 2 other post resources matching attributes
	for _, identifier := range []string{"456", "789"} {
		resourceGetResponse, err := authzClient.ResourceGet(ctx, &authz.ResourceGetRequest{
			Id: "post." + identifier,
		})
		if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
			resourceResponse, err := authzClient.ResourceCreate(ctx, &authz.ResourceCreateRequest{
				Id:    "post." + identifier,
				Kind:  "post",
				Value: identifier,
				Attributes: []*authz.Attribute{
					{Key: "owner_email", Value: "johndoe@acme.tld"},
				},
			})
			if err != nil {
				panic(fmt.Sprintf("cannot create resource: %v", err))
			}

			fmt.Printf("Resource created: %s\n", resourceResponse.GetResource().GetId())
		} else {
			fmt.Printf("Resource already exists: %s\n", resourceGetResponse.GetResource().GetId())
		}
	}

	// Create a policy
	policyGetResponse, err := authzClient.PolicyGet(ctx, &authz.PolicyGetRequest{
		Id: "post-owners",
	})
	if status, ok := status.FromError(err); ok && status.Code() == codes.NotFound {
		policyResponse, err := authzClient.PolicyCreate(ctx, &authz.PolicyCreateRequest{
			Id:        "post-owners",
			Resources: []string{"post.*"},
			Actions:   []string{"edit", "delete"},
			AttributeRules: []string{
				rule.AttributeEqual(
					rule.PrincipalResourceAttribute{
						PrincipalAttribute: "email",
						ResourceAttribute:  "owner_email",
					},
				),
			},
		})
		if err != nil {
			panic(fmt.Sprintf("cannot create policy: %v", err))
		}

		fmt.Printf("Policy created: %s\n", policyResponse.GetPolicy().GetId())
	} else {
		fmt.Printf("Policy already exists: %s\n", policyGetResponse.GetPolicy().GetId())
	}

	// Check if principal is allowed
	for _, identifier := range []string{"123", "456", "789"} {
		isAllowed, err := authzClient.IsAllowed(ctx, &authz.Check{
			Principal:     "user-123",
			ResourceKind:  "post",
			ResourceValue: identifier,
			Action:        "edit",
		})
		if err != nil {
			panic(fmt.Sprintf("cannot check if principal is allowed to edit post: %v", err))
		}

		fmt.Printf("Is allowed to edit post #%s? %v\n", identifier, isAllowed)
	}
}
