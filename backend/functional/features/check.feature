@check
Feature: check
  Test check-related APIs

  Scenario: Check for access
    Given I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "post.123", "kind": "post", "value": "123"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "my-post-123-policy-create",
        "resources": [
            "post.123"
        ],
        "actions": ["create"]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "my-post-123-policy-update",
        "resources": [
            "post.123"
        ],
        "actions": ["update"]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/roles" with payload:
      """
      {
        "id": "my-post-123-role-create-and-update",
        "policies": [
            "my-post-123-policy-create",
            "my-post-123-policy-update"
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "my-principal",
        "roles": [
            "my-post-123-role-create-and-update"
        ]
      }
      """
    And the response code should be 200
    When I send "POST" request to "/v1/check" with payload:
      """
      {
        "checks": [
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "create"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "update"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "delete"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "*",
            "action": "update"
          }
        ]
      }
      """
    And the response code should be 200
    And the response should match json:
      """
      {
        "checks": [
          {
            "action": "create",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "result": "allowed"
          },
          {
            "action": "update",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "result": "allowed"
          },
          {
            "action": "delete",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "result": "denied"
          },
          {
            "action": "update",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "*",
            "result": "denied"
          }
        ]
      }
      """
