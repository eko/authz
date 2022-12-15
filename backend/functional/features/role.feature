@role
Feature: role
  Test role-related APIs

  Scenario: Create a new role
    Given I send "POST" request to "/v1/resources" with payload:
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
    When I send "POST" request to "/v1/roles" with payload:
      """
      {
        "id": "my-post-123-role",
        "policies": [
            "my-post-123-policy"
        ]
      }
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "my-post-123-role",
        "policies": [
          {
            "id": "my-post-123-policy",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Update a role
    Given I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "post.123", "kind": "post", "value": "123"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "my-post-policy-create",
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
        "id": "my-post-policy-update",
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
        "id": "my-post-role",
        "policies": [
            "my-post-policy-create"
        ]
      }
      """
    And the response code should be 200
    When I send "PUT" request to "/v1/roles/my-post-role" with payload:
      """
      {
        "policies": [
            "my-post-policy-update"
        ]
      }
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "my-post-role",
        "policies": [
          {
            "id": "my-post-policy-update",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Retrieve a single role
    Given I send "POST" request to "/v1/resources" with payload:
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
    And I send "POST" request to "/v1/roles" with payload:
      """
      {
        "id": "my-post-123-role",
        "policies": [
            "my-post-123-policy"
        ]
      }
      """
    And the response code should be 200
    When I send "GET" request to "/v1/roles/my-post-123-role"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "my-post-123-role",
        "policies": [
          {
            "id": "my-post-123-policy",
            "updated_at": "2100-01-01T02:00:00+01:00",
            "created_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Delete a single role
    Given I send "POST" request to "/v1/resources" with payload:
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
    And I send "POST" request to "/v1/roles" with payload:
      """
      {
        "id": "my-post-123-role",
        "policies": [
            "my-post-123-policy"
        ]
      }
      """
    And the response code should be 200
    When I send "DELETE" request to "/v1/roles/my-post-123-role"
    And the response code should be 200
    And the response should match json:
      """
      {
        "success": true
      }
      """
    And I send "GET" request to "/v1/roles/my-post-123-role"
    And the response code should be 404

  Scenario: Retrieve a list of roles
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
        "id": "my-post-123-role-create",
        "policies": [
            "my-post-123-policy-create"
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/roles" with payload:
      """
      {
        "id": "my-post-123-role-update",
        "policies": [
            "my-post-123-policy-update"
        ]
      }
      """
    And the response code should be 200
    When I send "GET" request to "/v1/roles"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "data": [
          {
            "id": "my-post-123-role-create",
            "policies": [
              {
                "id": "my-post-123-policy-create",
                "created_at": "2100-01-01T02:00:00+01:00",
                "updated_at": "2100-01-01T02:00:00+01:00"
              }
            ],
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "id": "my-post-123-role-update",
            "policies": [
              {
                "id": "my-post-123-policy-update",
                "created_at": "2100-01-01T02:00:00+01:00",
                "updated_at": "2100-01-01T02:00:00+01:00"
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
