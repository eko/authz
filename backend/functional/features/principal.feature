@principal
Feature: principal
  Test principal-related APIs

  Scenario: Create a new principal
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/principals" with payload:
      """
      {"id": "f438dfb8-4ae9-4668-9545-f98dba4b2337"}
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "f438dfb8-4ae9-4668-9545-f98dba4b2337",
        "is_locked": false,
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Retrieve a single principal
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/principals" with payload:
      """
      {"id": "f438dfb8-4ae9-4668-9545-f98dba4b2337"}
      """
    And the response code should be 200
    When I send "GET" request to "/v1/principals/f438dfb8-4ae9-4668-9545-f98dba4b2337"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "f438dfb8-4ae9-4668-9545-f98dba4b2337",
        "is_locked": false,
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Delete a single principal
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/principals" with payload:
      """
      {"id": "f438dfb8-4ae9-4668-9545-f98dba4b2337"}
      """
    And the response code should be 200
    When I send "DELETE" request to "/v1/principals/f438dfb8-4ae9-4668-9545-f98dba4b2337"
    And the response code should be 200
    And the response should match json:
      """
      {
        "success": true
      }
      """
    And I send "GET" request to "/v1/principals/f438dfb8-4ae9-4668-9545-f98dba4b2337"
    And the response code should be 404

  Scenario: Retrieve a list of principals
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/principals" with payload:
      """
      {"id": "f438dfb8-4ae9-4668-9545-f98dba4b2337"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {"id": "another.value"}
      """
    And the response code should be 200
    When I send "GET" request to "/v1/principals"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "data": [
          {
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": "authz-admin",
            "is_locked": true,
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": "f438dfb8-4ae9-4668-9545-f98dba4b2337",
            "is_locked": false,
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "created_at": "2100-01-01T02:00:00+01:00",
            "id": "another.value",
            "is_locked": false,
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "page": 0,
        "size": 100,
        "total": 3
      }
      """
