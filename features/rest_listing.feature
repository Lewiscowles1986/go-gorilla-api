Feature: Absent resources return 404s.

  Scenario: A request for a known missing resource.
    When a GET is sent to `/product/ced7f0ac-f6f1-4f82-9890-f1dea2868ed9`
    Then the value of the response status is equal to `404`

  Scenario: A request for a known missing resource.
    When a DELETE is sent to `/product/ced7f0ac-f6f1-4f82-9890-f1dea2868ed9`
    Then the value of the response status is equal to `404`

  Scenario: A request for a known missing resource.
    When the request body is assigned:
      """
      {"name": "cheap guff",
       "price": 11.55}
      """
    And a PUT is sent to `/product/ced7f0ac-f6f1-4f82-9890-f1dea2868ed9`
    Then the value of the response status is equal to `404`
