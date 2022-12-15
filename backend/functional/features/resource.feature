@resource
Feature: resource
  Test resource-related APIs

  Scenario: Create a new resource (without value)
    Given I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "all-posts", "kind": "post"}
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "all-posts",
        "kind": "post",
        "value": "*",
        "is_locked": false,
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Create a new resource (with value)
    Given I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "custom-post", "kind": "post", "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681"}
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "custom-post",
        "kind": "post",
        "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681",
        "is_locked": false,
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Retrieve a single resource
    Given I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "all-posts", "kind": "post", "value": "*"}
      """
    And the response code should be 200
    When I send "GET" request to "/v1/resources/all-posts"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "all-posts",
        "kind": "post",
        "value": "*",
        "is_locked": false,
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00"
      }
      """

  Scenario: Delete a single resource
    Given I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "all-posts", "kind": "post", "value": "*"}
      """
    And the response code should be 200
    When I send "DELETE" request to "/v1/resources/all-posts"
    And the response code should be 200
    And the response should match json:
      """
      {
        "success": true
      }
      """
    And I send "GET" request to "/v1/resources/all-posts"
    And the response code should be 404

  Scenario: Retrieve a list of resources
    Given I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "all-posts", "kind": "post", "value": "*"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "custom-post", "kind": "post", "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681"}
      """
    And the response code should be 200
    When I send "GET" request to "/v1/resources"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "data": [
          {
            "id": "all-posts",
            "kind": "post",
            "value": "*",
            "is_locked": false,
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          },
          {
            "id": "custom-post",
            "kind": "post",
            "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681",
            "is_locked": false,
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00"
          }
        ],
        "page": 0,
        "size": 100,
        "total": 2
      }
      """
