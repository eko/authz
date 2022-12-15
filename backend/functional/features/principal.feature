@principal
Feature: principal
  Test principal-related APIs

  Scenario: Create a new principal
    Given I send "POST" request to "/v1/principals" with payload:
      """
      {"value": "f438dfb8-4ae9-4668-9545-f98dba4b2337"}
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": 1,
        "is_locked": false,
        "value": "f438dfb8-4ae9-4668-9545-f98dba4b2337",
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Update a principal
    Given I send "POST" request to "/v1/principals" with payload:
      """
      {"value": "f438dfb8-4ae9-4668-9545-f98dba4b2337"}
      """
    And the response code should be 200
    When I send "PUT" request to "/v1/principals/1" with payload:
      """
      {"value": "my-new-value"}
      """
    And the response code should be 200
    And the response should match json:
      """
      {
        "id": 1,
        "is_locked": false,
        "value": "my-new-value",
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Retrieve a single principal
    Given I send "POST" request to "/v1/principals" with payload:
      """
      {"value": "f438dfb8-4ae9-4668-9545-f98dba4b2337"}
      """
    And the response code should be 200
    When I send "GET" request to "/v1/principals/1"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": 1,
        "is_locked": false,
        "value": "f438dfb8-4ae9-4668-9545-f98dba4b2337",
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Delete a single principal
    Given I send "POST" request to "/v1/principals" with payload:
      """
      {"value": "f438dfb8-4ae9-4668-9545-f98dba4b2337"}
      """
    And the response code should be 200
    When I send "DELETE" request to "/v1/principals/1"
    And the response code should be 200
    And the response should match json:
      """
      {
        "success": true
      }
      """
    And I send "GET" request to "/v1/principals/1"
    And the response code should be 404

  Scenario: Retrieve a list of principals
    Given I send "POST" request to "/v1/principals" with payload:
      """
      {"value": "f438dfb8-4ae9-4668-9545-f98dba4b2337"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {"value": "another.value"}
      """
    And the response code should be 200
    When I send "GET" request to "/v1/principals"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "data": [
          {
            "id": 1,
            "is_locked": false,
            "value": "f438dfb8-4ae9-4668-9545-f98dba4b2337",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "id": 2,
            "is_locked": false,
            "value": "another.value",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "page": 0,
        "size": 100,
        "total": 2
      }
      """
