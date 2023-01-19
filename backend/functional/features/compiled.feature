@compiled
Feature: compiled
  Test compiled-policies APIs

  Scenario: List compiled policies
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "my-principal",
        "attributes": [
          {"key": "email", "value": "johndoe@acme.tld"}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.123",
        "kind": "post",
        "value": "123",
        "attributes": [
          {"key": "owner_email", "value": "johndoe@acme.tld"}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "my-post-policy",
        "resources": [
            "post.*"
        ],
        "actions": ["update", "delete"],
        "attribute_rules": [
            "principal.email == resource.owner_email"
        ]
      }
      """
    And the response code should be 200
    And I wait "500ms"
    When I send "GET" request to "/v1/compiled?filter=policy_id:contains:my-post-policy&sort=action_id:asc"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "data": [
          {
            "action_id": "delete",
            "created_at": "2100-01-01T01:00:00Z",
            "policy_id": "my-post-policy",
            "principal_id": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "updated_at": "2100-01-01T01:00:00Z",
            "version": 4102448400
          },
          {
            "action_id": "update",
            "created_at": "2100-01-01T01:00:00Z",
            "policy_id": "my-post-policy",
            "principal_id": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "updated_at": "2100-01-01T01:00:00Z",
            "version": 4102448400
          }
        ],
        "page": 0,
        "size": 100,
        "total": 2
      }
      """
