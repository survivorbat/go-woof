Feature: Parsing a table is a 🍰

  Scenario: Parsing a table to a struct works with default naming
    Given I have a struct type that looks like the following structure:
      """
      {
        "name": "string",
        "age": int,
        "type": "string",
        "vaccinated": bool
      }
      """

    When I use the FromTable function with the table:
      | name | age | type          | vaccinated | nickname |
      | Dex  | 5   | Dachshund     | true       | Dexy     |
      | Bob  | 1   | Berner Sennen | false      | NULL     |

    Then I expect a slice that resembles the following JSON:
      """
      [
        {"name": "Dex", "age": 5, "type": "Dachshund", "vaccinated": true, "nickname": "Dexy"},
        {"name": "Bob", "age": 1, "type": "Berner Sennen", "vaccinated": false, "nickname": null}
      ]
      """

  Scenario: Parsing a table vertically to a struct works with default naming
    Given I have a struct type that looks like the following structure:
      """
      {
        "name": "string",
        "age": int,
        "type": "string",
        "vaccinated": bool
      }
      """
    When I use the FromTable function with the table in vertical mode:
      | name       | Dex       | Bob           |
      | age        | 5         | 1             |
      | type       | Dachshund | Berner Sennen |
      | vaccinated | true      | false         |
      | nickname   | NULL      | Bobby         |

    Then I expect a slice that resembles the following JSON:
      """
      [
        {"name": "Dex", "age": 5, "type": "Dachshund", "vaccinated": true, "nickname": null},
        {"name": "Bob", "age": 1, "type": "Berner Sennen", "vaccinated": false, "nickname": "Bobby"}
      ]
      """
