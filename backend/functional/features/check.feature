@check
Feature: check
  Test check-related APIs

  Scenario: Check for access (using RBAC)
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
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
        "id": "my-post-123-role-create-and-update",
        "policies": [
            "my-post-123-policy-create",
            "my-post-123-policy-update"
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "my-principal",
        "roles": [
            "my-post-123-role-create-and-update"
        ]
      }
      """
    And the response code should be 200
    And I wait "500ms"
    When I send "POST" request to "/v1/check" with payload:
      """
      {
        "checks": [
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "create"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "update"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "delete"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "*",
            "action": "update"
          }
        ]
      }
      """
    And the response code should be 200
    And the response should match json:
      """
      {
        "checks": [
          {
            "action": "create",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "is_allowed": true
          },
          {
            "action": "update",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "is_allowed": true
          },
          {
            "action": "delete",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "is_allowed": false
          },
          {
            "action": "update",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "*",
            "is_allowed": false
          }
        ]
      }
      """

  Scenario: Check for access (using ABAC)
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.123",
        "kind": "post",
        "value": "123",
        "attributes": [
          {"key": "owner_id", "value": "owner-123"}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.456",
        "kind": "post",
        "value": "456",
        "attributes": [
          {"key": "owner_id", "value": "owner-456"}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.789",
        "kind": "post",
        "value": "789",
        "attributes": [
          {"key": "owner_id", "value": "owner-123"},
          {"key": "is_editable", "value": true}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "my-principal",
        "attributes": [
          {"key": "owner_id", "value": "owner-123"}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "my-post-123-policy-create",
        "resources": [
            "post.*"
        ],
        "actions": ["create"],
        "attribute_rules": [
          "principal.owner_id == resource.owner_id",
          "resource.is_editable == true"
        ]
      }
      """
    And the response code should be 200
    And I wait "500ms"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.10-updated-after",
        "kind": "post",
        "value": "10-updated-after",
        "attributes": [
          {"key": "owner_id", "value": "owner-123"},
          {"key": "is_editable", "value": true}
        ]
      }
      """
    And the response code should be 200
    And I wait "1s"
    When I send "POST" request to "/v1/check" with payload:
      """
      {
        "checks": [
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "create"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "update"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "456",
            "action": "create"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "789",
            "action": "create"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "10-updated-after",
            "action": "create"
          }
        ]
      }
      """
    And the response code should be 200
    And the response should match json:
      """
      {
        "checks": [
          {
            "action": "create",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "is_allowed": true
          },
          {
            "action": "update",
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "is_allowed": false
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "456",
            "action": "create",
            "is_allowed": false
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "789",
            "action": "create",
            "is_allowed": true
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "10-updated-after",
            "action": "create",
            "is_allowed": true
          }
        ]
      }
      """

  Scenario: Check for access (using ABAC greater operator)
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.123",
        "kind": "post",
        "value": "123",
        "attributes": [
          {"key": "number", "value": 10}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.456",
        "kind": "post",
        "value": "456",
        "attributes": [
          {"key": "number", "value": 20}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.789",
        "kind": "post",
        "value": "789"
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "my-principal"
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "my-post-greater-10",
        "resources": [
            "post.*"
        ],
        "actions": ["create"],
        "attribute_rules": [
          "resource.number > 10"
        ]
      }
      """
    And the response code should be 200
    And I wait "1s"
    When I send "POST" request to "/v1/check" with payload:
      """
      {
        "checks": [
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "create"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "update"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "456",
            "action": "create"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "789",
            "action": "create"
          }
        ]
      }
      """
    And the response code should be 200
    And the response should match json:
      """
      {
        "checks": [
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "create",
            "is_allowed": false
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "update",
            "is_allowed": false
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "456",
            "action": "create",
            "is_allowed": true
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "789",
            "action": "create",
            "is_allowed": false
          }
        ]
      }
      """

  Scenario: Check for access (using ABAC lower operator)
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.123",
        "kind": "post",
        "value": "123",
        "attributes": [
          {"key": "number", "value": 10}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.456",
        "kind": "post",
        "value": "456",
        "attributes": [
          {"key": "number", "value": 20}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/resources" with payload:
      """
      {
        "id": "post.789",
        "kind": "post",
        "value": "789"
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "my-principal"
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "my-post-greater-10",
        "resources": [
            "post.*"
        ],
        "actions": ["create"],
        "attribute_rules": [
          "resource.number < 20"
        ]
      }
      """
    And the response code should be 200
    And I wait "1s"
    When I send "POST" request to "/v1/check" with payload:
      """
      {
        "checks": [
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "create"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "update"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "456",
            "action": "create"
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "789",
            "action": "create"
          }
        ]
      }
      """
    And the response code should be 200
    And the response should match json:
      """
      {
        "checks": [
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "create",
            "is_allowed": true
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "123",
            "action": "update",
            "is_allowed": false
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "456",
            "action": "create",
            "is_allowed": false
          },
          {
            "principal": "my-principal",
            "resource_kind": "post",
            "resource_value": "789",
            "action": "create",
            "is_allowed": false
          }
        ]
      }
      """

  Scenario: Check for access (using ABAC principal rule on wildcard resource)
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "app.123", "kind": "app", "value": "123"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "app-admin",
        "attributes": [
          {"key": "applications", "value": "app1,app2"}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "other-admin",
        "attributes": [
          {"key": "applications", "value": "app10"}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "app1-admins-list",
        "resources": [
            "app.*"
        ],
        "actions": ["list"],
        "attribute_rules": [
          "principal.applications contains app1"
        ]
      }
      """
    And the response code should be 200
    And I wait "1s"
    When I send "POST" request to "/v1/check" with payload:
      """
      {
        "checks": [
          {
            "principal": "app-admin",
            "resource_kind": "app",
            "resource_value": "*",
            "action": "list"
          },
          {
            "principal": "app-admin",
            "resource_kind": "app",
            "resource_value": "123",
            "action": "list"
          },
          {
            "principal": "other-admin",
            "resource_kind": "app",
            "resource_value": "*",
            "action": "list"
          }
        ]
      }
      """
    And the response code should be 200
    And the response should match json:
      """
      {
        "checks": [
          {
            "principal": "app-admin",
            "resource_kind": "app",
            "resource_value": "*",
            "action": "list",
            "is_allowed": true
          },
          {
            "principal": "app-admin",
            "resource_kind": "app",
            "resource_value": "123",
            "action": "list",
            "is_allowed": true
          },
          {
            "principal": "other-admin",
            "resource_kind": "app",
            "resource_value": "*",
            "action": "list",
            "is_allowed": false
          }
        ]
      }
      """

  Scenario: Check for access (ABAC policy attached to a role only grants matching principals)
    Given I authenticate with username "admin" and password "changeme"
    And I send "POST" request to "/v1/resources" with payload:
      """
      {"id": "doc.1", "kind": "doc", "value": "1"}
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "doc-reader",
        "attributes": [
          {"key": "team", "value": "docs"}
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/policies" with payload:
      """
      {
        "id": "docs-team-read",
        "resources": ["doc.*"],
        "actions": ["read"],
        "attribute_rules": [
          "principal.team == docs"
        ]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/roles" with payload:
      """
      {
        "id": "docs-role",
        "policies": ["docs-team-read"]
      }
      """
    And the response code should be 200
    And I send "POST" request to "/v1/principals" with payload:
      """
      {
        "id": "intruder",
        "roles": ["docs-role"],
        "attributes": [
          {"key": "team", "value": "sales"}
        ]
      }
      """
    And the response code should be 200
    And I wait "1s"
    When I send "POST" request to "/v1/check" with payload:
      """
      {
        "checks": [
          {
            "principal": "doc-reader",
            "resource_kind": "doc",
            "resource_value": "1",
            "action": "read"
          },
          {
            "principal": "intruder",
            "resource_kind": "doc",
            "resource_value": "1",
            "action": "read"
          }
        ]
      }
      """
    And the response code should be 200
    And the response should match json:
      """
      {
        "checks": [
          {
            "principal": "doc-reader",
            "resource_kind": "doc",
            "resource_value": "1",
            "action": "read",
            "is_allowed": true
          },
          {
            "principal": "intruder",
            "resource_kind": "doc",
            "resource_value": "1",
            "action": "read",
            "is_allowed": false
          }
        ]
      }
      """
