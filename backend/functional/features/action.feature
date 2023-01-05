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
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
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
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": "create",
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": "delete",
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": "get",
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": "list",
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": "update",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "page": 0,
        "size": 100,
        "total": 5
      }
      """
