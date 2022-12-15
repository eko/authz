@action
Feature: action
  Test action-related APIs

  Scenario: Retrieve a single action
    Given I send "POST" request to "/v1/resources" with payload:
      """
      {"kind": "post", "value": "123"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "name": "my-post-123-policy-1",
        "resources": [
            {"kind": "post", "value": "123"}
        ],
        "actions": ["create"]
      }
      """
    And the response code should be 200
    When I send "GET" request to "/v1/actions/1"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": 1,
        "is_locked": false,
        "name": "create",
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Retrieve a list of actions
    Given I send "POST" request to "/v1/resources" with payload:
      """
      {"kind": "post", "value": "123"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "name": "my-post-123-policy-1",
        "resources": [
            {"kind": "post", "value": "123"}
        ],
        "actions": ["create", "delete"]
      }
      """
    And the response code should be 200
    When I send "GET" request to "/v1/actions"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "data": [
          {
            "id": 1,
            "is_locked": false,
            "name": "create",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "id": 2,
            "is_locked": false,
            "name": "delete",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "page": 0,
        "size": 100,
        "total": 2
      }
      """
