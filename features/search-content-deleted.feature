Feature: Search content deleted from elasticsearch

  Scenario: Delete event is processed and query is sent to Elasticsearch
    Given elasticsearch is healthy
    And elasticsearch returns the following response for bulk update
      """
      {
        "took": 13,
        "errors": false,
        "items": []
      }
      """
    When the service starts
    And this delete event is queued, to be consumed
    """
      {
        "URI":"some_deleted_uri"
      }
    """
    Then this delete query is sent to elasticsearch
      """
      {
        "query": {
          "term": {
            "uri":"some_deleted_uri"
          }
        }
      }
      """
