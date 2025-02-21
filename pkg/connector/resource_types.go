package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
)

var (
	userResourceType = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Description: "Pulumi user",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
	}

	teamResourceType = &v2.ResourceType{
		Id:          "team",
		DisplayName: "Team",
		Description: "Pulumi team",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_GROUP,
		},
	}

	orgResourceType = &v2.ResourceType{
		Id:          "organization",
		DisplayName: "Organization",
		Description: "Pulumi organization",
		Traits:      []v2.ResourceType_Trait{},
	}
)
