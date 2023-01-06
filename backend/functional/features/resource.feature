@resource
Feature: resource
  Test resource-related APIs

  Scenario: Create a new resource (without value)
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "all-posts", "kind": "post"}
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "id": "all-posts",
        "is_locked": false,
        "kind": "post",
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00",
        "value": "*"
      }
      """

  Scenario: Create a new resource (with value)
    Given I authenticate with username "admin" and password "changeme"
    When I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "custom-post",
        "kind": "post",
        "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681",
        "attributes": [
          {"key": "foo1", "value": "bar1"},
          {"key": "foo2", "value": "bar2"}
        ]
      }
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "attributes": [
          {
            "key": "foo1",
            "value": "bar1"
          },
          {
            "key": "foo2",
            "value": "bar2"
          }
        ],
        "created_at": "2100-01-01T02:00:00+01:00",
        "id": "custom-post",
        "is_locked": false,
        "kind": "post",
        "updated_at": "2100-01-01T02:00:00+01:00",
        "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681"
      }
      """

  Scenario: Update a resource
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "custom-post",
        "kind": "post",
        "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681",
        "attributes": [
          {"key": "foo1", "value": "bar1"},
          {"key": "foo2", "value": "bar2"}
        ]
      }
      """
    And the response code should be 200
    When I send "PUT" request to "/v1/resources/custom-post" with payload:
      """
      {
        "kind": "post",
        "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681",
        "attributes": [
          {"key": "foo3", "value": "bar3"},
          {"key": "foo4", "value": "bar4"}
        ]
      }
      """
    And the response code should be 200
    And the response should match json:
      """
      {
        "attributes": [
          {
            "key": "foo3",
            "value": "bar3"
          },
          {
            "key": "foo4",
            "value": "bar4"
          }
        ],
        "created_at": "2100-01-01T02:00:00+01:00",
        "id": "custom-post",
        "is_locked": false,
        "kind": "post",
        "updated_at": "2100-01-01T02:00:00+01:00",
        "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681"
      }
      """

  Scenario: Retrieve a single resource
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
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
        "is_locked": false,
        "kind": "post",
        "created_at": "2100-01-01T02:00:00+01:00",
        "updated_at": "2100-01-01T02:00:00+01:00",
        "value": "*"
      }
      """

  Scenario: Delete a single resource
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
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
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "all-posts", "kind": "post", "value": "*"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "custom-post", "kind": "post", "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681"}
      """
    And the response code should be 200
    When I send "GET" request to "/v1/resources?filter=kind:contains:post"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "data": [
          {
            "id": "all-posts",
            "is_locked": false,
            "kind": "post",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00",
            "value": "*"
          },
          {
            "id": "custom-post",
            "is_locked": false,
            "kind": "post",
            "created_at": "2100-01-01T02:00:00+01:00",
            "updated_at": "2100-01-01T02:00:00+01:00",
            "value": "97fdb1dc-b1e0-4652-ab82-5d174031a681"
          }
        ],
        "page": 0,
        "size": 100,
        "total": 2
      }
      """
