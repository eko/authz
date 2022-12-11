@policy
Feature: policy
  Test policy-related APIs

  Scenario: Create a new policy
    Given I send "POST" request to "/v1/actions" with payload:
      """
      {"name": "create"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"kind": "post", "value": "123"}
      """
    And the response code should be 200
    When I send "POST" request to "/v1/policies" with payload:
      """
      {
        "name": "my-post-123-policy",
        "resources": [
            {"kind": "post", "value": "123"}
        ],
        "actions": ["create"]
      }
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "actions": [
          {
            "id": 1,
            "is_locked": false,
            "name": "create",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "id": 1,
        "name": "my-post-123-policy",
        "resources": [
          {
            "id": 1,
            "is_locked": false,
            "kind": "post",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00",
            "value": "123"
          }
        ],
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Update a policy
    Given I send "POST" request to "/v1/actions" with payload:
      """
      {"name": "create"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/actions" with payload:
      """
      {"name": "update"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"kind": "post", "value": "123"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"kind": "post", "value": "456"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "name": "my-post-123-policy",
        "resources": [
            {"kind": "post", "value": "123"}
        ],
        "actions": ["create"]
      }
      """
    And the response code should be 200
    When I send "PUT" request to "/v1/policies/1" with payload:
      """
      {
        "name": "my-post-456-policy-update",
        "resources": [
            {"kind": "post", "value": "456"}
        ],
        "actions": ["update"]
      }
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "actions": [
          {
            "id": 2,
            "is_locked": false,
            "name": "update",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "id": 1,
        "name": "my-post-456-policy-update",
        "resources": [
          {
            "id": 2,
            "is_locked": false,
            "kind": "post",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00",
            "value": "456"
          }
        ],
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Retrieve a single policy
    Given I send "POST" request to "/v1/actions" with payload:
      """
      {"name": "create"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"kind": "post", "value": "123"}
      """
    And the response code should be 200
    When I send "POST" request to "/v1/policies" with payload:
      """
      {
        "name": "my-post-123-policy",
        "resources": [
            {"kind": "post", "value": "123"}
        ],
        "actions": ["create"]
      }
      """
    And the response code should be 200
    When I send "GET" request to "/v1/policies/1"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "actions": [
          {
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": 1,
            "is_locked": false,
            "name": "create",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "created_at": "2100-01-01T02:00:00+01:00",
        "id": 1,
        "name": "my-post-123-policy",
        "resources": [
          {
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": 1,
            "is_locked": false,
            "kind": "post",
            "updated_at": "2100-01-01T02:00:00+01:00",
            "value": "123"
          }
        ],
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Delete a single policy
    Given I send "POST" request to "/v1/actions" with payload:
      """
      {"name": "create"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"kind": "post", "value": "123"}
      """
    And the response code should be 200
    When I send "POST" request to "/v1/policies" with payload:
      """
      {
        "name": "my-post-123-policy",
        "resources": [
            {"kind": "post", "value": "123"}
        ],
        "actions": ["create"]
      }
      """
    And the response code should be 200
    When I send "DELETE" request to "/v1/policies/1"
    And the response code should be 200
    And the response should match json:
      """
      {
        "success": true
      }
      """
    And I send "GET" request to "/v1/subjects/1"
    And the response code should be 404

  Scenario: Retrieve a list of policies
    Given I send "POST" request to "/v1/actions" with payload:
      """
      {"name": "create"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
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
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "name": "my-post-123-policy-2",
        "resources": [
            {"kind": "post", "value": "123"}
        ],
        "actions": ["create"]
      }
      """
    And the response code should be 200
    When I send "GET" request to "/v1/policies"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "data": [
          {
            "actions": [
              {
                "id": 1,
                "is_locked": false,
                "name": "create",
                "created_at": "2100-01-01T02:00:00+01:00",
                "updated_at": "2100-01-01T02:00:00+01:00"
              }
            ],
            "id": 1,
            "name": "my-post-123-policy-1",
            "resources": [
              {
                "created_at": "2100-01-01T02:00:00+01:00",
                "id": 1,
                "is_locked": false,
                "kind": "post",
                "updated_at": "2100-01-01T02:00:00+01:00",
                "value": "123"
              }
            ],
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "actions": [
              {
                "id": 1,
                "is_locked": false,
                "name": "create",
                "created_at": "2100-01-01T02:00:00+01:00",
                "updated_at": "2100-01-01T02:00:00+01:00"
              }
            ],
            "id": 2,
            "name": "my-post-123-policy-2",
            "resources": [
              {
                "id": 1,
                "is_locked": false,
                "kind": "post",
                "created_at": "2100-01-01T02:00:00+01:00",
                "updated_at": "2100-01-01T02:00:00+01:00",
                "value": "123"
              }
            ],
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "page": 0,
        "size": 100,
        "total": 2
      }
      """
