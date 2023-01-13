@action
Feature: action
  Test action-related APIs

  Scenario: Retrieve a single action
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "post.123", "kind": "post", "value": "123"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "my-post-123-policy",
        "resources": [
            "post.123"
        ],
        "actions": ["create"]
      }
      """
    And the response code should be 200
    When I send "GET" request to "/v1/actions/create"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "create",
        "created_at": "2100-01-01T01:00:00Z",
        "updated_at": "2100-01-01T01:00:00Z"
      }
      """

  Scenario: Retrieve a list of actions
    Given I authenticate with username "admin" and password "changeme"
    When I send "GET" request to "/v1/actions?sort=id:asc"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "data": [
          {
            "created_at": "2100-01-01T01:00:00Z",
            "id": "create",
            "updated_at": "2100-01-01T01:00:00Z"
          },
          {
            "created_at": "2100-01-01T01:00:00Z",
            "id": "delete",
            "updated_at": "2100-01-01T01:00:00Z"
          },
          {
            "created_at": "2100-01-01T01:00:00Z",
            "id": "get",
            "updated_at": "2100-01-01T01:00:00Z"
          },
          {
            "created_at": "2100-01-01T01:00:00Z",
            "id": "list",
            "updated_at": "2100-01-01T01:00:00Z"
          },
          {
            "created_at": "2100-01-01T01:00:00Z",
            "id": "update",
            "updated_at": "2100-01-01T01:00:00Z"
          }
        ],
        "page": 0,
        "size": 100,
        "total": 5
      }
      """
