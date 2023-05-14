@user
Feature: user
  Test user-related APIs

  Scenario: Delete its own account
    Given I authenticate with username "admin" and password "changeme"
    And I send "DELETE" request to "/v1/users/admin"
    Then the response code should be 400
    And the response should match json:
      """
      {
        "error": true,
        "message": "a user cannot delete their own account"
      }
      """
